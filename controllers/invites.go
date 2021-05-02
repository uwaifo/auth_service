package controllers

import (
	"app-auth/auth"
	"app-auth/db"
	"app-auth/types"
	"app-auth/utils"
	"fmt"
	"os"
	"time"

	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mongodb/mongo-go-driver/bson"
)

func ConfirmMemberInvite(ctx echo.Context) error {
	mailId := ctx.Param("confirmationId")

	// we get the user object from redis.
	res, err := db.GetObject(mailId)
	if err != nil {
		log.Println(err)
		return ctx.JSON(http.StatusBadRequest, bson.M{"status": "Verification link has expired or is not available"})
	}

	if res == nil || len(res) < 1 {
		log.Println(err)
		return ctx.JSON(http.StatusBadRequest, bson.M{"status": "Verification link already used"})
	}

	app := res[0].(string)
	newUser := res[2].(string)
	userId := res[1].(string)
	teamId := res[3].(string)
	userEmail := res[4].(string)
	organisationId := res[5].(string)
	signupUrl := res[6].(string)
	appRedirectUrl := res[7].(string)

	userMemberScopeUtil := types.UserMemberScope{
		Id:             userId,
		OrganisationId: organisationId,
		TeamId:         teamId,
		UserId:         userId,
		App:            app,
	}

	userUtils := utils.UserMemberUtil{UserId: &userId, UserEmail: &userEmail, App: app}

	// if the user member scope is nil, the user does not exist
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

	// update the approved state of the
	userMemberScopeUpdateRes := userMemberScopeUtil.UpdateUserMemberScope(bson.D{{"state", "Approved"}})
	if userMemberScopeUpdateRes == nil {
		return ctx.JSON(http.StatusInternalServerError, types.TeamOperation{
			Message: fmt.Sprintf("Updating team user member %s state failed", userId),
			Status:  500,
			Error:   true,
			Type:    "ConfirmMemberInvite",
			State:   "Unsuccessful",
		})
	}

	// if the user scope exists, then the user exists
	userUpdateRes := userUtils.UpdateUser(bson.D{{"active", true}, {"verified", true}, {"lastloggedinscope", teamId}})
	if userUpdateRes == nil {
		return ctx.JSON(http.StatusInternalServerError, types.TeamOperation{
			Message: fmt.Sprintf("Updating team user member %s state failed", userId),
			Status:  500,
			Error:   true,
			Type:    "ConfirmMemberInvite",
			State:   "Unsuccessful",
		})
	}

	expires := time.Now().Add(15 * time.Minute)
	token := auth.CreateRSA256SignedToken(utils.ReadFile("./jwtRS256.key"), *userUpdateRes, userMemberScope, expires, app)

	// after all is done, remove the link
	// db.Del(mailId)

	// get the signup url
	if signupUrl == "" && os.Getenv("REDIRECT_URL") != "" {
		signupUrl = os.Getenv("REDIRECT_URL")
	} else if signupUrl == "" {
		signupUrl = "http://app-article-server:8081"
	}

	// set the redirect url from the signup url and the app redirect url
	// if the invited user already has an account, we take them to the app, else we redirect them to enter a one-time password
	// and then we take them to the app...
	redirectUrl := ""
	if newUser == "1" {
		redirectUrl = fmt.Sprintf(`%s?redirectUrl=%s&confirmEmail=%s&mailId=%s`, signupUrl, appRedirectUrl, userEmail, mailId)
	} else {
		redirectUrl = fmt.Sprintf(`%s?token=%s&expires=%s&scope=%s`, appRedirectUrl, token, expires.String(), userMemberScope.TeamId)
	}

	// redirect the user to where it's needed
	return ctx.Redirect(http.StatusTemporaryRedirect, redirectUrl)
}
