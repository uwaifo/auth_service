package controllers

import (
	"app-auth/types"
	"app-auth/utils"
	"log"
	"time"

	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/mongodb/mongo-go-driver/bson"
	uuid "github.com/satori/go.uuid"
)

// this route controller gets teams in an organisation
func GetTeams(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	admin := claims["id"].(string)
	app := claims["app"].(string)
	id := ctx.Param("id")

	organisationUtil := utils.OrganisationUtil{Id: id, Name: "", UserId: admin, App: app}

	teams := organisationUtil.GetOrganisationTeams()
	if teams == nil {
		return ctx.JSON(http.StatusNotFound, types.TeamNotFound{
			Message: "Organisation Teams Not Found",
			Status:  404,
			Error:   true,
			Type:    "GetTeams",
		})
	}

	return ctx.JSON(http.StatusOK, teams)
}

//
func GetAllDeletedTeams(ctx echo.Context) error {

	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	admin := claims["id"].(string)
	app := claims["app"].(string)
	// id := ctx.Param("id")

	teamUtil := utils.TeamUtil{UserId: admin, App: app}
	//organisationUtil := utils.OrganisationUtil{Id: id, Name: "", UserId: admin, App: app}

	//teams := organisationUtil.GetOrganisationDeletedTeams()
	teams := teamUtil.GetAllDeletedTeams()
	if teams == nil {
		return ctx.JSON(http.StatusNotFound, types.TeamNotFound{
			Message: "Organisation Teams Not Found",
			Status:  404,
			Error:   true,
			Type:    "GetTeams",
		})
	}

	return ctx.JSON(http.StatusOK, teams)

}

// this route controller returns organisation deleted teams
func GetDeletedTeams(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	admin := claims["id"].(string)
	app := claims["app"].(string)
	id := ctx.Param("id")

	organisationUtil := utils.OrganisationUtil{Id: id, Name: "", UserId: admin, App: app}

	teams := organisationUtil.GetOrganisationDeletedTeams()
	if teams == nil {
		return ctx.JSON(http.StatusNotFound, types.TeamNotFound{
			Message: "Organisation Teams Not Found",
			Status:  404,
			Error:   true,
			Type:    "GetTeams",
		})
	}

	return ctx.JSON(http.StatusOK, teams)

}

func GetDeletedTeamsCount(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	admin := claims["id"].(string)
	app := claims["app"].(string)

	organisationUtil := utils.OrganisationUtil{Id: "", Name: "", UserId: admin, App: app}
	deletedTeamCount := organisationUtil.GetOrganisationDeletedTeamsCount()
	if deletedTeamCount == nil {
		return ctx.JSON(http.StatusNotFound, types.TeamNotFound{
			Message: "Organisation Teams Not Found",
			Status:  404,
			Error:   true,
			Type:    "GetTeams",
		})
	}

	return ctx.JSON(http.StatusOK, deletedTeamCount)
}

// this route controller returns all active user teams
func UserTeams(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["id"].(string)
	userEmail := claims["email"].(string)
	app := claims["app"].(string)

	userUtil := utils.UserMemberUtil{UserId: &userId, UserEmail: &userEmail, App: app}

	teams := userUtil.GetUserActiveTeams()
	if teams == nil {
		return ctx.JSON(http.StatusNotFound, types.TeamNotFound{
			Message: "User Teams Not Found",
			Status:  404,
			Error:   true,
			Type:    "GetUserTeams",
		})
	}

	return ctx.JSON(http.StatusOK, teams)
}

// returns all the teams deleted
func UserTeamsDeleted(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["id"].(string)
	userEmail := claims["email"].(string)
	app := claims["app"].(string)

	userUtil := utils.UserMemberUtil{UserId: &userId, UserEmail: &userEmail, App: app}

	teams := userUtil.GetUserDeletedTeams()
	if teams == nil {
		return ctx.JSON(http.StatusNotFound, types.TeamNotFound{
			Message: "User Teams Not Found",
			Status:  404,
			Error:   true,
			Type:    "GetUserTeams",
		})
	}

	return ctx.JSON(http.StatusOK, teams)
}

func GetTeam(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	admin := claims["id"].(string)
	app := claims["app"].(string)
	id := ctx.Param("id")
	teamParam := ctx.Param("team")

	organisationUtil := utils.OrganisationUtil{Id: id, Name: "", UserId: admin, App: app}
	organisationUtil.GetOrganisation()

	teamUtil := utils.TeamUtil{TeamId: teamParam, TeamName: "", UserId: admin, App: app, Organisation: organisationUtil}

	team := teamUtil.GetTeam()
	if team == nil {
		return ctx.JSON(http.StatusNotFound, types.TeamNotFound{
			Message: fmt.Sprintf("Team %s Not Found", teamParam),
			Status:  404,
			Error:   true,
			Type:    "GetTeam",
		})
	}

	return ctx.JSON(http.StatusOK, team)
}

func DestroyTeam(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	admin := claims["id"].(string)
	app := claims["app"].(string)
	id := ctx.Param("id")
	teamParam := ctx.Param("team")

	organisationUtil := utils.OrganisationUtil{Id: id, Name: "", UserId: admin, App: app}
	organisationUtil.GetOrganisation()

	teamUtil := utils.TeamUtil{TeamId: teamParam, TeamName: "", UserId: admin, App: app, Organisation: organisationUtil}

	destroyTeam, err := teamUtil.DestroyTeamRecord()
	if !destroyTeam && err != nil {
		return ctx.JSON(http.StatusNotFound, types.TeamNotFound{
			Message: fmt.Sprintf("Team %s Not Found", teamParam),
			Status:  404,
			Error:   true,
			Type:    "GetTeam",
		})
	}

	deletedResponse := types.GenericResponse{
		Message: "Successfully Deleted",
		Status:  202,
		Error:   "false",
		Time:    time.Now().String(),
	}

	return ctx.JSON(http.StatusAccepted, deletedResponse)
}

func PostTeam(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	admin := claims["id"].(string)
	email := claims["email"].(string)
	app := claims["app"].(string)
	id := ctx.Param("id")
	teamUuid := ""

	if claims["scope"] != nil {
		scopes := claims["scope"].(map[string]interface{})
		teamUuid = scopes["team_id"].(string)
	}

	teamPostDetails := new(types.TeamPOSTPayload)
	err := ctx.Bind(teamPostDetails)
	if err != nil {
		fmt.Println(err)
		return err
	}

	organisationUtil := utils.OrganisationUtil{Id: id, Name: "", UserId: admin, App: app}
	organisationUtil.GetOrganisation()

	isValidOrganization := organisationUtil.GetOrganisation()

	if isValidOrganization.Deleted {
		return ctx.JSON(http.StatusBadRequest, types.TeamOperation{
			Message: fmt.Sprintf("Team Not Created. The parent organization %s has been deleted.", isValidOrganization.Name),
			Status:  400,
			Error:   true,
			Type:    "PostTeam",
			State:   "Unsuccessful",
		})
	}

	teamId := uuid.NewV4().String()
	teamUtil := utils.TeamUtil{TeamId: teamId, TeamName: teamPostDetails.Name, UserId: admin, UserEmail: email, App: app, Organisation: organisationUtil}

	team := teamUtil.CreateTeam(teamPostDetails.ImageUrl)
	if team == nil {
		return ctx.JSON(http.StatusInternalServerError, types.TeamOperation{
			Message: fmt.Sprintf("Team %s Not Created", teamId),
			Status:  500,
			Error:   true,
			Type:    "PostTeam",
			State:   "Unsuccessful",
		})
	}

	// the user has no available team at the moment
	if teamUuid == "" {
		userUtils := utils.UserMemberUtil{UserId: &admin, UserEmail: &email, App: app}
		// update the user logged in scope now then
		userUpdateRes := userUtils.UpdateUser(bson.D{{"active", true}, {"verified", true}, {"lastloggedinscope", teamId}})
		if userUpdateRes == nil {
			log.Println(fmt.Sprintf("Updating team user %s lastloggedinscope failed", admin))
		}
	}

	// return the created team
	return ctx.JSON(http.StatusCreated, team)
}

func PutTeam(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	admin := claims["id"].(string)
	app := claims["app"].(string)
	id := ctx.Param("id")
	teamParam := ctx.Param("team")

	teamPutDetails := new(types.TeamPUTPayload)
	err := ctx.Bind(teamPutDetails)
	if err != nil {
		fmt.Println(err)
		return err
	}

	organisationUtil := utils.OrganisationUtil{Id: id, Name: "", UserId: admin, App: app}
	organisationUtil.GetOrganisation()

	teamUtil := utils.TeamUtil{TeamId: teamParam, TeamName: "", UserId: admin, App: app, Organisation: organisationUtil}

	updateData := bson.D{{"name", teamPutDetails.Name}, {"imageurl", teamPutDetails.ImageUrl}}

	team := teamUtil.UpdateTeam(updateData)
	if team == nil {
		return ctx.JSON(http.StatusInternalServerError, types.TeamOperation{
			Message: fmt.Sprintf("Team %s Not Updated", teamParam),
			Status:  500,
			Error:   true,
			Type:    "PutTeam",
			State:   "Unsuccessful",
		})
	}

	return ctx.JSON(http.StatusOK, team)
}

func DeleteTeam(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	admin := claims["id"].(string)
	email := claims["email"].(string)
	app := claims["app"].(string)
	id := ctx.Param("id")
	teamParam := ctx.Param("team")

	organisationUtil := utils.OrganisationUtil{Id: id, Name: "", UserId: admin, App: app}
	organisationUtil.GetOrganisation()

	teamUtil := utils.TeamUtil{TeamId: teamParam, TeamName: "", UserId: admin, UserEmail: email, App: app, Organisation: organisationUtil}

	removedTeam := teamUtil.RemoveTeamTemp()
	if removedTeam == nil {
		return ctx.JSON(http.StatusInternalServerError, types.TeamOperation{
			Message: fmt.Sprintf("Team %s Not Removed", teamParam),
			Status:  500,
			Error:   true,
			Type:    "DeleteTeam",
			State:   "Unsuccessful",
		})
	}

	fmt.Println("user_id: ", admin)
	fmt.Println("time is : ", time.Now())
	return ctx.JSON(http.StatusCreated, removedTeam)

	/*
		return ctx.JSON(http.StatusOK, types.TeamOperation{
			Message: fmt.Sprintf("Team %s Removed", id),
			Status:  200,
			Error:   false,
			Type:    "DeleteTeam",
			State:   "Success",
		})
	*/
}

func RestoreTeam(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	admin := claims["id"].(string)
	app := claims["app"].(string)
	id := ctx.Param("id")
	teamParam := ctx.Param("team")

	//	userUtil := utils.UserMemberUtil{UserId: &userId, UserEmail: &userEmail, App: app}
	organisationUtil := utils.OrganisationUtil{Id: id, Name: "", UserId: admin, App: app}
	organisation := organisationUtil.GetOrganisation()

	// Check if the organization exits

	if organisation == nil {
		return ctx.JSON(http.StatusNotFound, types.OrganisationNotFound{
			Message: "Organisations Not Found",
			Status:  404,
			Error:   true,
			Type:    "GetOrganisations",
		})
	}

	if organisation.Deleted {
		return ctx.JSON(http.StatusBadRequest, types.TeamOperation{
			Message: "Restoration Error: Team's Organisation is deleted",
			Status:  400,
			Error:   true,
			Type:    "RestoreTeam",
			State:   "Unsuccessful",
		})

	}

	teamUtil := utils.TeamUtil{TeamId: teamParam, TeamName: "", UserId: admin, App: app, Organisation: organisationUtil}
	restoreTeam := teamUtil.RestoreDeletedTeam()
	if restoreTeam == nil {
		return ctx.JSON(http.StatusBadRequest, types.TeamOperation{
			Message: fmt.Sprintf("Team %s Not Restored", teamParam),
			Status:  400,
			Error:   true,
			Type:    "RestoreTeam",
			State:   "Unsuccessful",
		})
	}

	return ctx.JSON(http.StatusAccepted, restoreTeam)

}
