package main

import (
	"app-auth/auth"
	"app-auth/config"
	"app-auth/controllers"
	"app-auth/schedule"

	"log"
	"os"
	"regexp"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	// instantiate the server
	e := echo.New()
	e.HideBanner = true

	// configure a logger middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, error=${error}\n",
	}))

	// allow development URL in CORS setup
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:8081",
			"http://localhost:8000",
			"http://localhost:10101",
			"http://127.0.0.1:8081",
			"http://localhost:5000",
			"http://app-auth:5000",
			"http://app-auth",
			"http://app-server:8000",
			"http://app-server:8080",
			"http://app-server",
			"http://app-test:2000",
			"http://app-test",
			"http://app-user-server:3005",
			"http://app-user-server:8082",
			"http://localhost:8082",
			"http://app-user-server",
			"http://localhost:3000",
			"http://app-user:3000",
			"http://app-user",
			"http://app-article-server:8081",
			"http://app-article-server",
			"http://app-article:6001",
			"http://app-article",
			"https://app.staging.clipsynphony.com",
			"https://www.staging.clipsynphony.com",
			"https://www.testing.clipsynphony.com",
			"https://www.dev.clipsynphony.com",
			"https://www.clipsynphony.com",
			"https://www.clipsymphony.com",
			"http://*.35.246.198.190.nip.io",
			"https://*.35.246.198.190.nip.io",
			"http://*.staging.id.scaratec.com",
			"https://*.staging.id.scaratec.com",
			"http://www.pr-{{NUMBER}}.staging.id.scaratec.com",
			"https://www.pr-{{NUMBER}}.staging.id.scaratec.com",
			"http://www.staging.id.scaratec.com",
			"https://www.staging.id.scaratec.com",
			"http://*.prod.staging.clipsynphony.com",
			"http://*.staging.clipsynphony.com",
			"https://*.prod.staging.clipsynphony.com",
			"https://*.staging.clipsynphony.com",
			"http://www.prod.staging.clipsynphony.com",
			"http://www.staging.clipsynphony.com",
			"https://www.prod.staging.clipsynphony.com",
			"https://www.staging.clipsynphony.com",
			"https://www.staging.tiermedizin-forum.de",
			"https://www.tiermedizin-forum.de",
		},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))

	// setup jwt authentication middleware for processing requests...
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    auth.GetSecret(),
		TokenLookup:   "header:" + echo.HeaderAuthorization,
		AuthScheme:    "Bearer",
		SigningMethod: "RS256",
		Skipper: func(context echo.Context) bool {
			// regex for organization
			orgRegex, err := regexp.Compile(`(^\/ping$|^\/user.+)`)
			if err != nil {
				log.Println(err)
			}

			return orgRegex.MatchString(context.Request().URL.String())
		},
	}))

	// use the cookie middleware
	e.Use(session.Middleware(controllers.Store))

	e.GET("/ping", controllers.PingHandler)

	e.GET("/user/callback", controllers.OAuthCallback)
	// authentication with social login
	e.GET("/user/auth/:provider", controllers.AuthenticateWithSocial)
	// authentication callback using social login
	e.GET("/user/auth/:provider/callback", controllers.SocialCallback)

	// user login. posting user data
	e.POST("/user/login", controllers.LoginPost)

	// user password reset
	e.POST("/user/password-reset/request", controllers.PasswordResetRequest)
	e.POST("/user/password-reset/request/", controllers.PasswordResetRequest)

	// user login. posting user data
	e.POST("/user/password-reset/change", controllers.PasswordResetChange)
	e.POST("/user/password-reset/change/", controllers.PasswordResetChange)

	//
	//e.POST("/email/change/:redirect/:emailchangeid", controllers.)
	//e.POST("/email/change/:redirect/:emailchangeid", controllers.)

	// when the user is invited and the user isn't logged in or doesn't have an account on id.scaratec.com
	// this handles the little popup screen where the user logs in and gets authentication
	// [Will be deprecated for a finer confirmation handling]
	e.POST("/user/login-signup-confirm/:confirmationId", controllers.LoginSignUpPost)
	e.POST("/user/login-signup-confirm/:confirmationId/", controllers.LoginSignUpPost)
	// handling user signup
	e.POST("/user/signup", controllers.SignupPost)
	e.POST("/user/tm/signup", controllers.SignupTmPost)
	e.GET("/user/tm/admin-confirm/:redisUserData", controllers.SignupTmAdminConfirm)
	e.GET("/user/tm/admin-confirm/:redisUserData/", controllers.SignupTmAdminConfirm)
	// handling user logout
	e.GET("/user/logout", controllers.Logout)
	// handling user signup email confirmation
	e.GET("/user/confirm/:confirmationId/redirect/:redirectId", controllers.ConfirmEmail)
	// handling user invite confirmations
	e.GET("/user/user-invite-confirm/:confirmationId", controllers.ConfirmMemberInvite)

	// get stats for the currently logged in user
	e.GET("/stats", controllers.UserStats)
	// get stats for the currently logged in user
	e.GET("/me", controllers.GetMe)
	// Edit user information
	e.PUT("/me/edit", controllers.EditUserInfo)
	e.PUT("/me/edit/", controllers.EditUserInfo)
	// update user password
	e.PUT("/me/password", controllers.EditUserPassword)
	e.PUT("/me/password/", controllers.EditUserPassword)
	e.PUT("/me/password/update", controllers.EditUserPassword)
	e.PUT("/me/password/update/", controllers.EditUserPassword)

	// user email reset request
	e.POST("/me/email-reset/request", controllers.EmailResetRequest)
	e.POST("/me/email-reset/request/", controllers.EmailResetRequest)

	//process email reset verification link
	//http://app-auth:8080/user/email/change/?emailChangeId=0fce79ab-ce61-4cfd-9b92-63b3dc004370
	e.GET("user/email/change", controllers.ProcessEmailReset)
	e.GET("user/email/change/", controllers.ProcessEmailReset)

	// get teams for the currently logged in user
	e.GET("/teams", controllers.UserTeams)
	e.GET("/teams/deleted", controllers.UserTeamsDeleted)
	e.GET("/teams/deleted/", controllers.UserTeamsDeleted)
	// assign default team to the user
	e.GET("/assign-default-team/:id", controllers.AssignDefaultTeam)
	// get a count of all the deleted teams within a specific organisation
	e.GET("/stats/deleted-teams/count", controllers.GetDeletedTeamsCount)
	// get a count of deleted organizations
	e.GET("/stats/deleted-organisations/count", controllers.GetDeletedOrganisationsCount)
	// get all deleted teams accross organisations
	e.GET("/stats/deleted-teams", controllers.GetAllDeletedTeams)

	// get a list of all the organisation belonging to a user that are not deleted
	e.GET("/organisations", controllers.GetOrganisations)
	// get a list of all the deleted organisation belonging to a user
	e.GET("/organisations/deleted", controllers.GetDeletedOrganisations)
	e.GET("/organisations/deleted/", controllers.GetDeletedOrganisations)
	// get a specific organisation by id

	e.GET("/organisation/:id", controllers.GetOrganisation)
	// post an organisation
	e.POST("/organisation/", controllers.PostOrganisation)
	e.POST("/organisation", controllers.PostOrganisation)
	// remove a specific organisation by id
	//e.DELETE("/organisation/:id", controllers.DeleteOrganisation)
	e.DELETE("/organisation/:id", controllers.DeleteOrganisation)

	// permamently detroy an organisation redord
	e.DELETE("/organisation/destroy/:id", controllers.DestroyOrganisation)
	e.DELETE("/organisation/destroy/:id/", controllers.DestroyOrganisation)

	// permamently detroy a team redord
	e.DELETE("/organisation/:id/destroyteam/:team", controllers.DestroyTeam)
	e.DELETE("/organisation/:id/destroyteam/:team/", controllers.DestroyTeam)

	// restore a deleted organisation
	e.POST("/organisation/restore/:id", controllers.RestoreOrganisation)

	// update a specific organisation by id
	e.PUT("/organisation/:id", controllers.PutOrganisation)

	// get all the teams within a specific organisation
	e.GET("/organisation/:id/teams", controllers.GetTeams)
	// get all the deleted teams within a specific organisation
	e.GET("/organisation/:id/deletedteams", controllers.GetDeletedTeams)
	// get a specific team within a speicific organisation by Id
	e.GET("/organisation/:id/team/:team", controllers.GetTeam)
	// add a team to a specific organisation
	e.POST("/organisation/:id/team/", controllers.PostTeam)
	e.POST("/organisation/:id/team", controllers.PostTeam)
	// remove a team within an organisation by ID
	e.DELETE("/organisation/:id/team/:team", controllers.DeleteTeam)
	e.DELETE("/organisation/:id/team/:team/", controllers.DeleteTeam)

	// restore a team within an organisation by ID
	e.POST("/organisation/:id/team/:team/restore", controllers.RestoreTeam)
	// update a team within a specific organisation by Id
	e.PUT("/organisation/:id/team/:team", controllers.PutTeam)

	// add/invite a new member to a team
	e.POST("/organisation/:id/team/:team/user", controllers.AddTeamMember)
	e.POST("/organisation/:id/team/:team/user/", controllers.AddTeamMember)
	// remove a member within a team
	e.DELETE("/organisation/:id/team/:team/user/:userId", controllers.RemoveTeamMember)
	// get a member of a team by Id
	e.GET("/organisation/:id/team/:team/user/:userId", controllers.GetTeamMember)
	// get all members of a team
	e.GET("/organisation/:id/team/:team/users", controllers.GetTeamMembers)
	e.GET("/organisation/:id/team/:team/users/", controllers.GetTeamMembers)

	// update the scope of a user within a team
	e.POST("/organisation/:id/team/:team/user/:userId/scope", controllers.AddTeamMemberScope)
	// delete an assign scope to user
	e.DELETE("/organisation/:id/team/:team/user/:userId/scope", controllers.RemoveTeamMemberScope)

	// for posting new IAM scopes
	// get all iam scopes created within for an application
	e.GET("/iam", controllers.GetIAMScopes)
	// get user scopes by email
	e.GET("/me/scopes", controllers.GetScopeByUserEmail)
	// get a specific IAM scope
	e.GET("/iam/:id", controllers.GetIAMScope)
	// add an IAM scope within an application
	e.POST("/iam/", controllers.PostIAMScope)
	e.POST("/iam", controllers.PostIAMScope)
	// remove an IAM scope
	e.DELETE("/iam/:id", controllers.DeleteIAMScope)
	// update an IAM scope
	e.PUT("/iam/:id", controllers.PutIAMScope)

	// update the permissions
	// Add new permissions to an IAM scope
	e.POST("/iam/:id/permission", controllers.PostIAMScopesPermission)
	e.POST("/iam/:id/permissions", controllers.PostIAMScopesPermissions)
	// remove permission within an IAM scope
	e.DELETE("/iam/:id/permission", controllers.DeleteIAMScopesPermission)

	// get a new access token using a refresh token
	e.GET("/token/refresh", controllers.GetRefreshToken)
	e.POST("/token/refresh", controllers.RefreshToken)

	// access denied for unsupported methods
	e.GET("/", controllers.AccessDenied)

	// get the port from the env files
	port := os.Getenv("PORT")
	if port == "" {
		port = config.HttpPort
	}

	// Start all the pending jobs
	//gocron.Every(30).Second().Do(task)
	//go gocron.Start()
	schedule.Starter()

	// start the server with the port
	e.Logger.Fatal(e.Start(port))
}
