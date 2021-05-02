package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"app-auth/iam"
	"app-auth/types"

	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"

	"app-auth/auth"
	"app-auth/db"
	"app-auth/session"
	"app-auth/utils"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/markbates/goth/gothic"
	"github.com/mongodb/mongo-go-driver/bson"
)

var _ = godotenv.Load()
var Store = session.NewCookieStore()
var _ = auth.Goth()

func PingHandler(req echo.Context) error {

	sess, err := Store.Get(req.Request(), "user-session")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(sess.Values)
	fmt.Println(req.Cookies())

	return req.String(http.StatusOK, "[PONG] AUTH SERVER IS UP\n")
}

func OAuthCallback(req echo.Context) error {

	responseBody := CallResponseJSON{"Authentication success", "ok"}

	req.Set("Content-Type", "application/json")

	return req.JSON(http.StatusOK, responseBody)
}

func LoginPost(req echo.Context) error {

	loginDetails := new(types.UserLoginPOSTPayload)

	err := req.Bind(loginDetails)
	if err != nil {
		fmt.Println(err)
		res := types.LoginResponse{Token: "", UserId: "", Expires: time.Now(), Message: "Email / User Invalid"}
		return req.JSON(http.StatusNotFound, res)
	}

	userUtil := utils.UserMemberUtil{UserId: nil, UserEmail: &loginDetails.Email}

	// if there's no such user
	resultContainer := userUtil.GetUserByEmail()
	if resultContainer == nil {
		res := types.LoginResponse{Token: "", UserId: "", Expires: time.Now(), Message: "Email / User Invalid"}
		return req.JSON(http.StatusNotFound, res)
	}

	userUtil = utils.UserMemberUtil{UserId: &resultContainer.Id, UserEmail: &loginDetails.Email, App: loginDetails.App}

	if resultContainer.Verified == false {
		res := types.LoginResponse{Token: "", Expires: time.Now(), Message: "Please verify your email"}
		return req.JSON(http.StatusNotFound, res)
	}

	if utils.CheckPasswordHash(loginDetails.Password, resultContainer.Password) {
		expires := time.Now().Add(15 * time.Minute)
		refresh := time.Now().Add(14 * 24 * time.Hour)
		scopes := userUtil.GetUserScopeWithID(resultContainer.LastLoggedInScope)
		token := auth.CreateRSA256SignedToken(utils.ReadFile("./jwtRS256.key"), *resultContainer, scopes, expires, loginDetails.App)
		refreshToken := auth.CreateRSA256RefreshToken(utils.ReadFile("./jwtRS256.key"), *resultContainer, refresh, loginDetails.App)
		teams := *userUtil.GetUserActiveTeams()
		res := types.LoginResponse{Token: token, UserId: resultContainer.Id, Expires: expires, Message: "User Login Successful", Refresh: refreshToken, Scope: resultContainer.LastLoggedInScope, Teams: teams}
		return req.JSON(http.StatusOK, res)
	} else {
		res := types.LoginResponse{Token: "", UserId: "", Expires: time.Now(), Message: "User Email / Password Invalid", Scope: ""}
		return req.JSON(http.StatusNotFound, res)
	}
}

func AssignDefaultTeam(req echo.Context) error {

	user := req.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	admin := claims["id"].(string)
	app := claims["app"].(string)
	email := claims["email"].(string)
	id := req.Param("id")

	userUtil := utils.UserMemberUtil{UserId: &admin, UserEmail: &email, App: app}

	// if there's no such user
	resultContainer := userUtil.GetUserById()
	if resultContainer == nil {
		res := types.LoginResponse{Token: "", UserId: "", Expires: time.Now(), Message: "Email / User Invalid"}
		return req.JSON(http.StatusNotFound, res)
	}

	upData := bson.D{{"lastloggedinscope", id}}
	updatedUser := userUtil.UpdateUser(upData)

	log.Print(updatedUser)

	if resultContainer.Verified == false {
		res := types.LoginResponse{Token: "", Expires: time.Now(), Message: "Please verify your email"}
		return req.JSON(http.StatusNotFound, res)
	}

	expires := time.Now().Add(15 * time.Minute)
	refresh := time.Now().Add(14 * 24 * time.Hour)
	scopes := userUtil.GetUserScopeWithID(id)
	log.Print(scopes)
	token := auth.CreateRSA256SignedToken(utils.ReadFile("./jwtRS256.key"), *resultContainer, scopes, expires, app)
	refreshToken := auth.CreateRSA256RefreshToken(utils.ReadFile("./jwtRS256.key"), *resultContainer, refresh, app)
	res := types.LoginResponse{Token: token, UserId: resultContainer.Id, Expires: expires, Message: "User Login Successful", Refresh: refreshToken, Scope: resultContainer.LastLoggedInScope, Teams: nil}
	return req.JSON(http.StatusOK, res)
}

func LoginSignUpPost(req echo.Context) error {
	mailId := req.Param("confirmationId")
	loginDetails := new(types.UserLoginSignupPOSTPayload)

	err := req.Bind(loginDetails)
	if err != nil {
		fmt.Println(err)
		return err
	}

	userUtil := utils.UserMemberUtil{UserId: nil, UserEmail: &loginDetails.Email, App: loginDetails.App}

	// check if this user exists, if not return an error
	resultContainer := userUtil.GetUserByEmail()
	if resultContainer == nil {
		res := types.LoginResponse{Token: "", UserId: "", Expires: time.Now(), Message: "User Does not exist"}
		return req.JSON(http.StatusBadRequest, res)
	}

	// assign the user Id to the util
	userUtil = utils.UserMemberUtil{UserId: &resultContainer.Id, UserEmail: &loginDetails.Email, App: loginDetails.App}
	// update the user password
	userUtil.UpdateUser(bson.D{{"password", utils.HashPassword(loginDetails.Password)}})

	// then generate the refresh and user tokens
	expires := time.Now().Add(15 * time.Minute)
	refresh := time.Now().Add(14 * 24 * time.Hour)
	scopes := userUtil.GetUserScopeWithID(resultContainer.LastLoggedInScope)
	token := auth.CreateRSA256SignedToken(utils.ReadFile("./jwtRS256.key"), *resultContainer, scopes, expires, loginDetails.App)
	refreshToken := auth.CreateRSA256RefreshToken(utils.ReadFile("./jwtRS256.key"), *resultContainer, refresh, loginDetails.App)
	teams := *userUtil.GetUserActiveTeams()

	db.Del(mailId)

	// send this as the response
	res := types.LoginResponse{Token: token, UserId: resultContainer.Id, Expires: expires, Message: "User Login Successful", Refresh: refreshToken, Scope: resultContainer.LastLoggedInScope, Teams: teams}
	return req.JSON(http.StatusCreated, res)
}

func PasswordResetRequest(req echo.Context) error {
	requestDetails := new(types.UserPasswordResetRequestPOSTPayload)

	err := req.Bind(requestDetails)
	if err != nil {
		fmt.Println(err)
		res := types.GenericResponse{Status: http.StatusBadRequest, Message: "Unable to parse post request"}
		return req.JSON(http.StatusBadRequest, res)
	}

	userUtil := utils.UserMemberUtil{UserId: nil, UserEmail: &requestDetails.Email, App: requestDetails.App}

	// check if this user exists, if not return an error
	resultContainer := userUtil.GetUserByEmail()
	if resultContainer == nil {
		res := types.GenericResponse{Status: http.StatusNotFound, Error: "User not found", Message: "User Does not exist"}
		return req.JSON(http.StatusNotFound, res)
	}

	// generate a new uuid as a redis id for password reset
	redirectId := uuid.NewV4().String()
	err = db.Set(redirectId, resultContainer.Id)
	if err != nil {
		fmt.Print(err)
		res := types.GenericResponseWithError{Status: http.StatusFailedDependency, Error: err, Message: "Unable to set reset ID to redis"}
		return req.JSON(http.StatusFailedDependency, res)
	}

	// then create a new url and send to the user email
	changeRequestURL := fmt.Sprintf("%s?resetId=%s&email=%s", requestDetails.RedirectUrl, redirectId, resultContainer.Email)
	err = utils.SendPasswordResetEmail(resultContainer.Email, changeRequestURL)
	if err != nil {
		fmt.Print(err)
		res := types.GenericResponseWithError{Status: http.StatusFailedDependency, Error: err, Message: "Unable to send email!"}
		return req.JSON(http.StatusFailedDependency, res)
	}

	res := types.GenericResponseWithError{Message: "Password reset request successful", Error: nil, Status: http.StatusOK}
	return req.JSON(http.StatusOK, res)
}

func PasswordResetChange(req echo.Context) error {
	changeDetails := new(types.UserPasswordResetChangePOSTPayload)

	err := req.Bind(changeDetails)
	if err != nil {
		fmt.Println(err)
		res := types.GenericResponseWithError{Status: http.StatusFailedDependency, Message: "Unable to parse post request", Error: err}
		return req.JSON(http.StatusBadRequest, res)
	}

	userId := db.Get(changeDetails.ResetId)
	if userId == "" {
		return req.JSON(http.StatusNotFound, types.GenericResponseWithError{Status: http.StatusFailedDependency, Message: "Password Reset verification link expired or unavailable", Error: nil})
	}

	userUtil := utils.UserMemberUtil{UserId: &userId, UserEmail: &changeDetails.Email, App: changeDetails.App}

	// check if this user exists, if not return an error
	resultContainer := userUtil.GetUserById()
	if resultContainer == nil {
		return req.JSON(http.StatusNotFound, types.GenericResponseWithError{Status: http.StatusFailedDependency, Message: "User Does not exist", Error: nil, Time: time.Now().String()})
	}

	// update the user password
	userUtil.UpdateUser(bson.D{{"password", utils.HashPassword(changeDetails.Password)}})

	// remove the reset ID generated
	db.Del(changeDetails.ResetId)

	// send this as the response
	res := types.GenericResponseWithError{Status: http.StatusCreated, Message: "Password successfully updated!", Error: nil, Time: time.Now().String()}
	return req.JSON(http.StatusCreated, res)
}

func SignupPost(req echo.Context) error {
	signDetails := new(types.UserSignUpPOSTPayload)

	err := req.Bind(signDetails)
	if err != nil {
		fmt.Println(err)
		return err
	}

	userUtil := utils.UserMemberUtil{UserId: nil, UserEmail: &signDetails.Email, App: signDetails.App}

	// there's no such user
	resultContainer := userUtil.GetUserByEmail()
	if resultContainer == nil {
		userData := types.UserData{
			Id:                uuid.NewV4().String(),
			Password:          utils.HashPassword(signDetails.Password),
			Email:             signDetails.Email,
			Username:          signDetails.Username,
			Firstname:         signDetails.FirstName,
			Lastname:          signDetails.LastName,
			Picture:           signDetails.Picture,
			Age:               0,
			Active:            true,
			Verified:          false,
			Provider:          "local",
			LastLoggedInScope: "",
		}

		res := userUtil.CreateNewUser(userData)
		if res == nil {
			res := types.LoginResponse{Token: "", UserId: "", Expires: time.Now(), Message: "User Creation Unsuccessful"}
			return req.JSON(http.StatusBadRequest, res)
		}

		// send the user an email
		err = utils.SendSignupEmail(userData.Email, userData.Id, signDetails.App, signDetails.AppRedirectUrl)
		if err != nil {
			log.Print(err)
			res := types.LoginResponse{Token: "", UserId: "", Expires: time.Now(), Message: "Error Sending Signup Email. Try Again Later"}
			return req.JSON(http.StatusBadRequest, res)
		}

		response := types.SignupResponse{Time: time.Now(), Status: "200", Message: "User Sign up successful. Please confirm your email."}

		return req.JSON(http.StatusCreated, response)
	}

	res := types.SignupResponse{Time: time.Now(), Status: "404", Message: "User email is already registered"}
	return req.JSON(http.StatusBadRequest, res)
}

func SignupTmPost(req echo.Context) error {
	signDetails := new(types.UserSignUpPOSTPayload)

	err := req.Bind(signDetails)
	if err != nil {
		fmt.Println(err)
		return err
	}

	userUtil := utils.UserMemberUtil{UserId: nil, UserEmail: &signDetails.Email, App: signDetails.App}

	// there's no such user
	resultContainer := userUtil.GetUserByEmailTm()
	if resultContainer == nil {
		userData := types.UserData{
			Id:                uuid.NewV4().String(),
			Password:          utils.HashPassword(signDetails.Password),
			Email:             signDetails.Email,
			Username:          signDetails.Username,
			Firstname:         signDetails.FirstName,
			Lastname:          signDetails.LastName,
			Picture:           signDetails.Picture,
			Age:               0,
			Active:            false,
			Verified:          false,
			Provider:          "local",
			LastLoggedInScope: "",
		}

		res := userUtil.CreateNewUserTm(userData)
		if res == nil {
			res := types.LoginResponse{Token: "", UserId: "", Expires: time.Now(), Message: "User Creation Unsuccessful"}
			return req.JSON(http.StatusBadRequest, res)
		}

		// generate a redisId
		redisAdminConfirmId := uuid.NewV4().String()
		// send the URL of the ID to the admin account
		redisObject := map[string]interface{}{
			"userId":         userData.Id,
			"userEmail":      userData.Email,
			"appRedirectUrl": signDetails.AppRedirectUrl,
			"app":            signDetails.App,
		}

		if err := db.SetObject(redisAdminConfirmId, redisObject); err != nil {
			res := types.LoginResponse{Token: "", UserId: "", Expires: time.Now(), Message: "Error caching user data for admin confirmation"}
			return req.JSON(http.StatusBadRequest, res)
		}

		authUrl := os.Getenv("ENDPOINT")
		if authUrl == "" {
			authUrl = "http://app-auth/user"
		}

		// send the user an email
		err = utils.SendAdminConfirmationMail("dannymcwaves@icloud.com", authUrl+"/tm/admin-confirm/"+redisAdminConfirmId, userData.Email, userData.Username)
		if err != nil {
			log.Print(err)
			res := types.LoginResponse{Token: "", UserId: "", Expires: time.Now(), Message: "Error Sending Admin Confirmation Email. Try Again Later"}
			return req.JSON(http.StatusBadRequest, res)
		}

		response := types.SignupResponse{Time: time.Now(), Status: "200", Message: "User Sign up successful. Waiting for admin to confirm Signup"}

		return req.JSON(http.StatusCreated, response)
	}

	message := "User email is already registered"
	// user is not yet accepted by admin
	if !resultContainer.Verified {
		message = "Email already registered. Please verify you email"
	}
	// user is not yet accepted by admin
	if !resultContainer.Active {
		message = "Email is registered. Waiting for admin to approve signup"
	}

	res := types.SignupResponse{Time: time.Now(), Status: "402", Message: message}
	return req.JSON(http.StatusBadRequest, res)
}

func SignupTmAdminConfirm(req echo.Context) error {
	redisId := req.Param("redisUserData")

	data, err := db.GetObjectByKey(redisId)
	if err != nil {
		res := types.SignupResponse{Time: time.Now(), Status: "404", Message: "User Signup Not Available!"}
		return req.JSON(http.StatusNotFound, res)
	}

	app := data["app"]
	userId := data["userId"]
	userEmail := data["userEmail"]
	appRedirectUrl := data["appRedirectUrl"]

	userUtil := utils.UserMemberUtil{UserId: &userId, UserEmail: &userEmail, App: app}

	// there's no such user
	resultContainer := userUtil.GetUserByIdTm()
	if resultContainer != nil {
		user := userUtil.UpdateUserTm(bson.D{{"active", true}})
		if user == nil {
			res := types.LoginResponse{Token: "", UserId: "", Expires: time.Now(), Message: "User confirmation Update Unsuccessful"}
			return req.JSON(http.StatusBadRequest, res)
		}

		// send the user an email
		err = utils.SendSignupEmail(userEmail, userId, app, appRedirectUrl)
		if err != nil {
			res := types.LoginResponse{Token: "", UserId: "", Expires: time.Now(), Message: "Error Sending Verification Email to user. Try Again Later"}
			return req.JSON(http.StatusBadRequest, res)
		}

		_ = db.Del(redisId)

		response := types.SignupResponse{Time: time.Now(), Status: "200", Message: "User Verification Email sent!"}
		return req.JSON(http.StatusCreated, response)
	}

	res := types.SignupResponse{Time: time.Now(), Status: "404", Message: "User Data Not Found!"}
	return req.JSON(http.StatusNotFound, res)
}

func Logout(req echo.Context) error {
	cookie := utils.CreateHttpCookie("", time.Now(), true)
	req.SetCookie(cookie)
	return req.JSON(http.StatusOK, types.LogoutResponse{Status: "200", Message: "Logout successful"})
}

func SocialCallback(req echo.Context) error {

	user, err := gothic.CompleteUserAuth(req.Response(), req.Request())
	if err != nil {
		log.Print(err)
		return req.JSON(http.StatusBadRequest, bson.M{"error": err, "status": 400, "message": "Authentication failed"})
	}

	app := req.QueryParam("app")
	if app == "" {
		app = "clipsynphony"
	}

	log.Print(user)

	userData := utils.GenerateUserFromGoogleAuth(user)
	log.Print(userData)

	userUtil := utils.UserMemberUtil{UserId: &userData.Id, UserEmail: &userData.Email, App: app}

	// there's no such user
	resultContainer := userUtil.GetUserByEmail()
	if resultContainer == nil {
		res := userUtil.CreateNewUser(userData)
		if res == nil {
			response := types.LoginResponse{Token: "", UserId: "", Expires: time.Now(), Message: "User Creation Unsuccessful"}
			return req.JSON(http.StatusBadRequest, response)
		}

		// then assign the created user data to the resultContainer
		resultContainer = res
	}

	scopes := userUtil.GetUserScopeWithID(resultContainer.LastLoggedInScope)
	expires := time.Now().Add(15 * time.Minute)
	token := auth.CreateRSA256SignedToken(utils.ReadFile("./jwtRS256.key"), userData, scopes, expires, app)

	RedirectUrl := os.Getenv("REDIRECT_URL")
	if RedirectUrl == "" {
		RedirectUrl = "http://localhost:8000"
	}

	return req.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s?token=%s&expires=%s", RedirectUrl, token, expires))
}

func AuthenticateWithSocial(req echo.Context) error {

	log.Print(req.Request().URL.Query().Get("state"))

	req.Request().URL.Query().Set("provider", req.Param("provider"))

	fmt.Println(req.Request().URL.Query())

	// in this case we should find the cookies set by the goth lib.
	fmt.Println(req.Cookies())

	if user, err := gothic.CompleteUserAuth(req.Response(), req.Request()); err == nil {
		return req.JSON(http.StatusOK, user)
	} else {
		log.Print(err)
		gothic.BeginAuthHandler(req.Response(), req.Request())
	}

	return req.String(http.StatusOK, "3u45uby8934nvw8t92vw89")
}

func ConfirmEmail(ctx echo.Context) error {
	confirmationId := ctx.Param("confirmationId")
	redirectId := ctx.Param("redirectId")

	userId := db.Get(confirmationId)
	if userId == "" {
		return ctx.JSON(http.StatusNotFound, bson.M{"status": "Verification link has expired or is not available"})
	}
	redirectUrl := db.Get(redirectId)

	app := ctx.QueryParam("app")

	userUtil := utils.UserMemberUtil{UserId: &userId, UserEmail: nil, App: app}

	// there's no such user
	resultContainer := userUtil.GetUserById()
	if resultContainer == nil {
		// user was not found
		loginPage := os.Getenv("SIGNUP_URL")
		if loginPage == "" {
			loginPage = "http://app-server:8080"
		}
		return ctx.Redirect(http.StatusTemporaryRedirect, loginPage)
	} else {
		// if the user is found, the update the verified status of the user account
		updata := bson.D{{"verified", true}}
		user := userUtil.UpdateUser(updata)
		if user == nil {
			res := types.LoginResponse{Token: "", UserId: "", Expires: time.Now(), Message: "User verification Update Unsuccessful"}
			return ctx.JSON(http.StatusBadRequest, res)
		}

		userUtil.UserId = &resultContainer.Id

		expires := time.Now().Add(15 * time.Minute)
		scopes := userUtil.GetUserScopeWithID(resultContainer.LastLoggedInScope)
		token := auth.CreateRSA256SignedToken(utils.ReadFile("./jwtRS256.key"), *user, scopes, expires, app)

		// remove the redis id that's already confirmed
		_ = db.Del(confirmationId)
		_ = db.Del(redirectId)

		if redirectUrl == "" && os.Getenv("REDIRECT_URL") != "" {
			redirectUrl = os.Getenv("REDIRECT_URL")
		} else if redirectUrl == "" {
			redirectUrl = "http://app-article-server:8081"
		}

		return ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf(`%s?token=%s&expires=%s&scope=%s`, redirectUrl, token, expires, resultContainer.LastLoggedInScope))
	}
}

func RefreshToken(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["id"].(string)
	userEmail := claims["email"].(string)
	app := claims["app"].(string)

	refresh := new(types.RefreshToken)
	err := ctx.Bind(refresh)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if refresh.Token == "" {
		return ctx.JSON(http.StatusFailedDependency, types.InvalidRequest{Status: http.StatusFailedDependency, Message: "Refresh token missing"})
	}

	userUtil := utils.UserMemberUtil{UserId: &userId, UserEmail: &userEmail}

	// there's no such user
	resultContainer := userUtil.GetUserByEmail()
	if resultContainer == nil {
		res := types.LoginResponse{Token: "", UserId: "", Expires: time.Now(), Message: "Email / User Invalid"}
		return ctx.JSON(http.StatusNotFound, res)
	}

	expires := time.Now().Add(15 * time.Minute)
	refreshTime := time.Now().Add(14 * 24 * time.Hour)
	scopes := userUtil.GetUserScopeWithID(resultContainer.LastLoggedInScope)
	token := auth.CreateRSA256SignedToken(utils.ReadFile("./jwtRS256.key"), *resultContainer, scopes, expires, app)
	refreshToken := auth.CreateRSA256RefreshToken(utils.ReadFile("./jwtRS256.key"), *resultContainer, refreshTime, app)

	res := types.LoginResponse{Token: token, Expires: expires, Message: "User Login Successful", Refresh: refreshToken, UserId: resultContainer.Id, Scope: resultContainer.LastLoggedInScope}

	return ctx.JSON(http.StatusOK, res)
}

func GetRefreshToken(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["id"].(string)
	userEmail := claims["email"].(string)
	app := claims["app"].(string)

	userUtil := utils.UserMemberUtil{UserId: &userId, UserEmail: &userEmail}

	// there's no such user
	resultContainer := userUtil.GetUserById()
	if resultContainer == nil {
		res := types.LoginResponse{Token: "", UserId: "", Expires: time.Now(), Message: "Email / User Invalid"}
		return ctx.JSON(http.StatusNotFound, res)
	}

	refreshTime := time.Now().Add(14 * 24 * time.Hour)
	refreshToken := auth.CreateRSA256RefreshToken(utils.ReadFile("./jwtRS256.key"), *resultContainer, refreshTime, app)

	res := types.LoginResponse{Message: "Refresh Token Generated", Refresh: refreshToken, UserId: resultContainer.Id, Scope: resultContainer.LastLoggedInScope}

	return ctx.JSON(http.StatusOK, res)
}

func UserStats(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["id"].(string)
	userEmail := claims["email"].(string)
	username := claims["username"].(string)
	verified := claims["verified"].(bool)
	active := claims["active"].(bool)
	app := claims["app"].(string)

	userStats := types.UserStats{}
	userStats.Username = username
	userStats.Verified = verified
	userStats.Active = active
	userStats.ReceivedInvites = 0
	userStats.SentInvites = 0

	userUtil := utils.UserMemberUtil{UserId: &userId, UserEmail: &userEmail, App: app}

	teams := userUtil.GetUserActiveTeams()
	if teams == nil {
		teams = &[]types.Team{}
	}

	organisations := userUtil.GetUserActiveOrganisations()
	if organisations == nil {
		organisations = &[]types.Organisation{}
	}

	userStats.Organisations = len(*organisations)
	userStats.Teams = len(*teams)

	return ctx.JSON(http.StatusOK, userStats)
}

func GetMe(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["id"].(string)
	userEmail := claims["email"].(string)
	app := claims["app"].(string)

	userUtil := utils.UserMemberUtil{UserId: &userId, UserEmail: &userEmail, App: app}

	userData := userUtil.GetUserById()
	if userData == nil {
		userData = &types.UserData{}
	}

	return ctx.JSON(http.StatusOK, userData)
}

func EditUserInfo(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["id"].(string)
	userEmail := claims["email"].(string)
	app := claims["app"].(string)

	userUtil := utils.UserMemberUtil{UserId: &userId, UserEmail: &userEmail, App: app}
	userData := userUtil.GetUserById()
	if userData == nil {
		return ctx.JSON(
			http.StatusNotFound,
			types.GenericResponse{
				Message: "User not Found",
				Status:  http.StatusNotFound,
				Time:    time.Now().String(),
				Error:   "User not Found by email or Id. cannot update user details",
			})
	}

	// get the data from the request
	userUpdateDetails := new(types.UserEditPUTPayload)
	err := ctx.Bind(userUpdateDetails)
	if err != nil {
		fmt.Println(err)
		return ctx.JSON(
			http.StatusOK,
			types.GenericResponse{
				Message: "Error parsing request data",
				Status:  http.StatusOK,
				Time:    time.Now().String(),
				Error:   "Error parsing data in request",
			})
	}

	// update the user information
	updatedUser := userUtil.UpdateUser(bson.D{
		{"email", userUpdateDetails.Email},
		{"lastname", userUpdateDetails.LastName},
		{"firstname", userUpdateDetails.FirstName},
		{"username", userUpdateDetails.UserName},
		{"picture", userUpdateDetails.Picture},
	})

	// if there's an error updating the user information
	// then return that an error to the requesting user
	if updatedUser == nil {
		return ctx.JSON(
			http.StatusInternalServerError,
			types.GenericResponse{
				Message: "Error updating user data",
				Status:  http.StatusInternalServerError,
				Time:    time.Now().String(),
				Error:   "Error updating user information",
			})
	}

	return ctx.JSON(http.StatusOK, updatedUser)
}

func EditUserPassword(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["id"].(string)
	userEmail := claims["email"].(string)
	app := claims["app"].(string)

	// init a user util module
	userUtil := utils.UserMemberUtil{UserId: &userId, UserEmail: &userEmail, App: app}

	// get the user information
	userData := userUtil.GetUserById()
	if userData == nil {
		return ctx.JSON(
			http.StatusOK,
			types.GenericResponse{
				Message: "User not Found",
				Status:  http.StatusNotFound,
				Time:    time.Now().String(),
				Error:   "User not Found by email or Id. cannot update user details",
			})
	}

	// get the data from the request
	userPasswordUpdateDetails := new(types.UserUpdatePasswordPUTPayload)
	err := ctx.Bind(userPasswordUpdateDetails)
	if err != nil {
		fmt.Println(err)
		return ctx.JSON(
			http.StatusOK,
			types.GenericResponse{
				Message: "Error parsing request data",
				Status:  http.StatusOK,
				Time:    time.Now().String(),
				Error:   "Error parsing data in request",
			})
	}

	// if the old password does not match the existing password
	if !utils.CheckPasswordHash(userPasswordUpdateDetails.OldPassword, userData.Password) {
		return ctx.JSON(
			http.StatusOK,
			types.GenericResponse{
				Message: "Old Password Incorrect",
				Status:  http.StatusOK,
				Time:    time.Now().String(),
				Error:   "Old password provided does not match currently existing password",
			})
	}

	// update the user password
	updatedUser := userUtil.UpdateUser(bson.D{{"password", utils.HashPassword(userPasswordUpdateDetails.Password)}})

	// if there was an error updating the user information
	if updatedUser == nil {
		return ctx.JSON(
			http.StatusInternalServerError,
			types.GenericResponse{
				Message: "Error updating user data",
				Status:  http.StatusInternalServerError,
				Time:    time.Now().String(),
				Error:   "Error updating user information",
			})
	}

	// return an update response
	return ctx.JSON(
		http.StatusOK,
		types.GenericResponse{
			Message: "Password successfully updated",
			Status:  http.StatusOK,
			Time:    time.Now().String(),
			Error:   "",
		})
}

func AccessDenied(ctx echo.Context) error {
	return ctx.String(http.StatusForbidden, "Access Denied")
}

func EmailResetRequest(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["id"].(string)
	userEmail := claims["email"].(string)
	//app := claims["app"].(string)
	// init a user util module
	userUtil := utils.UserMemberUtil{UserId: &userId, UserEmail: &userEmail}
	// find the user by email and ID using the JWT token details
	userData := userUtil.GetUserById()
	// if the user does not exist, return 404

	if userData == nil {
		return ctx.JSON(
			http.StatusOK,
			types.GenericResponse{
				Message: "User not Found",
				Status:  http.StatusNotFound,
				Time:    time.Now().String(),
				Error:   "User not Found by email or Id. cannot update user details",
			})
	}
	// get the data from the request
	resetEmailRequest := new(types.UserEmailChangeRequestPostPayload)
	err := ctx.Bind(resetEmailRequest)
	if err != nil {
		fmt.Println(err)
		return ctx.JSON(
			http.StatusOK,
			types.GenericResponse{
				Message: "Error parsing request data",
				Status:  http.StatusOK,
				Time:    time.Now().String(),
				Error:   "Error parsing data in request",
			})
	}
	// check if the user password from the request object is valid, if it’s invalid return error
	if !utils.CheckPasswordHash(resetEmailRequest.Password, userData.Password) {
		return ctx.JSON(
			http.StatusOK,
			types.GenericResponse{
				Message: "Password Incorrect",
				Status:  http.StatusOK,
				Time:    time.Now().String(),
				Error:   "Password provided does not match currently existing password",
			})
	}

	// check if the new email exists… if it does, return an appropriate error notifying email existence.
	newUserUtil := utils.UserMemberUtil{UserId: nil, UserEmail: &resetEmailRequest.NewEmail}

	validNewEmail := newUserUtil.GetUserByEmail()
	if validNewEmail != nil {
		res := types.GenericResponseWithError{Status: http.StatusBadRequest, Error: err, Message: "New email already in use"}

		return ctx.JSON(http.StatusNotFound, res)
	}

	redisEmailChangeId := uuid.NewV4().String()
	err = db.SetObject(redisEmailChangeId, map[string]interface{}{"userId": userData.Id, "userNewEmail": resetEmailRequest.NewEmail})
	if err != nil {
		fmt.Println(err)
		res := types.GenericResponseWithError{
			Status:  http.StatusFailedDependency,
			Error:   err,
			Message: "Unable to set the ID in redis",
		}
		return ctx.JSON(http.StatusFailedDependency, res)
	}

	//Create a redirectUrl

	redirectUrl := os.Getenv("ENVIROMENT")
	if redirectUrl == "" {
		redirectUrl = "http://app-auth:8080"
	}
	//emailSentURL := fmt.Sprintf("%s?redirectUrl=%s&emailChangeId=%s", authEndpoint, resetEmailRequest.RedirectUrl, redisEmailChangeId)

	emailSentURL := fmt.Sprintf("%s/email/change?redirectUrl=%s&emailChangeId=%s", redirectUrl, resetEmailRequest.RedirectUrl, redisEmailChangeId)
	fmt.Println(emailSentURL)

	// send email to the new user email, if there’s an error sending the email, return email sending error

	err = utils.SendResetEmail(resetEmailRequest.NewEmail, emailSentURL)
	if err != nil {
		fmt.Print(err)
		res := types.GenericResponseWithError{Status: http.StatusFailedDependency, Error: err, Message: "Unable to send email!"}
		return ctx.JSON(http.StatusFailedDependency, res)
		//fmt.Println(res)
	}

	res := types.EmailResetSuccess{Message: "Password reset request successful", Error: nil, Status: http.StatusOK, EmailChangeId: redisEmailChangeId}

	// If all actions pass, return a success message object

	return ctx.JSON(http.StatusOK, res)

}

func ProcessEmailReset(ctx echo.Context) error {

	// we get the redisEmailChangeId from the request params
	redisEmailChangeId := ctx.QueryParam("emailChangeId")
	//Get the redirectUrl Param from the request object
	//redirectUrl := ctx.QueryParam("redirectUrl")

	// we look for the corresponding object in redis, if the redis object  saved cannot be found, then return an error.
	reqObj, err := db.GetObjectByKey(redisEmailChangeId)
	if err != nil {
		fmt.Println(err)
		return ctx.JSON(http.StatusBadRequest, bson.M{"status": "The email reset link can not be found"})
	}

	// check the redis object to make sure it's not empty or the length is not less than 2, if it is then return an error saying the link is already used.

	if reqObj == nil || len(reqObj) < 1 {
		fmt.Println(err)
		return ctx.JSON(http.StatusBadRequest, bson.M{"status": "Unable to get a valid object associated with that emailChangeId. It may have expired"})
	}

	// Read the data from the redis object
	userId := reqObj["userId"]
	userNewEmail := reqObj["userNewEmail"]

	updateUserUtil := types.ConfirmationData{
		Userid: userId,
		Email:  userNewEmail,
	}
	userUtils := utils.UserMemberUtil{UserId: &userId, UserEmail: &userNewEmail}

	findUser := userUtils.GetUserByEmail()
	if findUser == nil {
		fmt.Println("e no dey")

	}

	fmt.Println(updateUserUtil.Email)

	// Find the user by Id
	userUtil := utils.UserMemberUtil{UserId: &userId, UserEmail: nil}
	userData := userUtil.GetUserById()
	if userData == nil {
		return ctx.JSON(
			http.StatusNotFound,
			types.GenericResponse{
				Message: "User not Found",
				Status:  http.StatusNotFound,
				Time:    time.Now().String(),
				Error:   "User not Found by Id. cannot update user details",
			})

	}

	// 	When the user exists, then you update the user email to the new email
	updateData := bson.D{{"email", updateUserUtil.Email}}
	updatedUser := userUtil.UpdateUser(updateData)
	if updatedUser == nil {
		return ctx.JSON(
			http.StatusInternalServerError,
			types.GenericResponse{
				Message: "Error updating user data",
				Status:  http.StatusInternalServerError,
				Time:    time.Now().String(),
				Error:   "Error updating user information",
			})
	}

	// Find all the user member scopes with the userId and update the email.

	userMemberScopeUtil := types.UserMemberScope{
		Id: userData.Id,
		//OrganisationId: organisationId,
		//TeamId:         teamId,
		UserId: userId,
		//App:            app,
	}

	userMemberScope := userMemberScopeUtil.GetUserMemberScope()
	if userMemberScope == nil {
		return ctx.JSON(http.StatusNotFound, types.TeamOperation{
			Message: fmt.Sprintf("Team user member %s not found", userId),
			Status:  404,
			Error:   true,
			Type:    "ConfirmMemberInvite",
			State:   "Unsuccessful",
		})
	}

	// TODO
	userMemberScopeUpdateRes := userMemberScopeUtil.UpdateUserMemberScope(bson.D{{"useremail", updateUserUtil.Email}})
	if userMemberScopeUpdateRes == nil {
		return ctx.JSON(http.StatusInternalServerError, types.TeamOperation{
			Message: fmt.Sprintf("Updating team user member %s state failed", userId),
			Status:  500,
			Error:   true,
			Type:    "ConfirmMemberInvite",
			State:   "Unsuccessful",
		})
	}
	// 	Generate a new user token
	expires := time.Now().Add(15 * time.Minute)
	// TODO ASK danny about the lasp app parameter that I have omited in the bellow token
	token := auth.CreateRSA256SignedToken(utils.ReadFile("./jwtRS256.key"), *updatedUser, userMemberScope, expires, "")

	// redirect the user to the redirect url and pass the generated token as a token param
	redirectUrl := os.Getenv("ENVIROMENT")
	if redirectUrl == "" {
		redirectUrl = "http://app-auth:8080"
	}

	redirectUrl = fmt.Sprintf(`%s?token=%s&expires=%s&scope=%s`, redirectUrl, token, expires.String(), userData.LastLoggedInScope)

	//Redirect the user
	return ctx.Redirect(http.StatusTemporaryRedirect, redirectUrl)

}

func GetScopeByUserEmail(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["id"].(string)
	userEmail := claims["email"].(string)
	//app := claims["app"].(string)
	// init a user util module
	userUtil := utils.UserMemberUtil{UserId: &userId, UserEmail: &userEmail}
	// find the user by email and ID using the JWT token details
	userData := userUtil.GetUserById()
	// if the user does not exist, return 404

	if userData == nil {
		return ctx.JSON(
			http.StatusOK,
			types.GenericResponse{
				Message: "User not Found",
				Status:  http.StatusNotFound,
				Time:    time.Now().String(),
				Error:   "User not Found by email or Id. cannot update user details",
			})
	}
	iamScopes := iam.ScopeActions{Id: "", App: "", Name: "", Permission: []string{}}
	allUserScopes := iamScopes.GetScopesByEmail(userData.Email)
	if allUserScopes == nil {
		return ctx.JSON(http.StatusNotFound, types.ScopeNotFound{
			Message: "Scopes Not Found",
			Status:  404,
			Error:   true,
			Type:    "GetIAMScopes",
		})

	}

	return ctx.JSON(http.StatusOK, allUserScopes)
}
