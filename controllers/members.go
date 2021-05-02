package controllers

import (
	"app-auth/types"
	"app-auth/utils"

	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

func GetTeamMembers(ctx echo.Context) error  {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	admin := claims["id"].(string)
	app := claims["app"].(string)
	id := ctx.Param("id")
	teamParam := ctx.Param("team")

	organisationUtil := utils.OrganisationUtil{Id: id, Name: "", UserId: admin, App: app,}
	organisationUtil.GetOrganisation()

	teamUtil := utils.TeamUtil{TeamId: teamParam, TeamName: "", UserId: admin, App: app, Organisation: organisationUtil}

	members := teamUtil.GetMembers(); if members == nil {
		return ctx.JSON(http.StatusNotFound, types.TeamNotFound {
			Message: fmt.Sprintf("No Members Found"),
			Status: 404,
			Error: true,
			Type: "GetTeamMembers",
		})
	}

	return ctx.JSON(http.StatusOK, members)
}

func GetTeamMember(ctx echo.Context) error  {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	admin := claims["id"].(string)
	app := claims["app"].(string)
	id := ctx.Param("id")
	teamParam := ctx.Param("team")
	userParam := ctx.Param("userId")

	organisationUtil := utils.OrganisationUtil{Id: id, Name: "", UserId: admin, App: app,}
	organisationUtil.GetOrganisation()

	teamUtil := utils.TeamUtil{TeamId: teamParam, TeamName: "", UserId: admin, App: app, Organisation: organisationUtil}

	member := teamUtil.GetMember(userParam); if member == nil {
		return ctx.JSON(http.StatusNotFound, types.TeamNotFound {
			Message: fmt.Sprintf("Member %s Not Found", userParam),
			Status: 404,
			Error: true,
			Type: "GetTeamMember",
		})
	}

	return ctx.JSON(http.StatusOK, member)
}

func AddTeamMember(ctx echo.Context) error  {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	admin := claims["id"].(string)
	app := claims["app"].(string)
	id := ctx.Param("id")
	teamParam := ctx.Param("team")

	memberPostDetails := new(types.TeamMemberInvitePOSTPayload)
	err := ctx.Bind(memberPostDetails); if err != nil {
		fmt.Println(err)
		return err
	}

	organisationUtil := utils.OrganisationUtil{Id: id, Name: "", UserId: admin, App: app,}
	organisationUtil.GetOrganisation()

	teamUtil := utils.TeamUtil{TeamId: teamParam, TeamName: "", UserId: admin, App: app, Organisation: organisationUtil}

	member := teamUtil.AddMember(memberPostDetails.Email, memberPostDetails.SignupUrl, memberPostDetails.AppRedirectUrl, app); if member == nil {
		return ctx.JSON(http.StatusInternalServerError, types.TeamOperation {
			Message: fmt.Sprintf("Team Member %s Not Created", memberPostDetails.Email),
			Status: 500,
			Error: true,
			Type: "PostTeamMember",
			State: "Unsuccessful",
		})
	}

	return ctx.JSON(http.StatusCreated, member)
}

func RemoveTeamMember(ctx echo.Context) error  {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	admin := claims["id"].(string)
	app := claims["app"].(string)
	id := ctx.Param("id")
	teamParam := ctx.Param("team")
	userParam := ctx.Param("userId")

	organisationUtil := utils.OrganisationUtil{Id: id, Name: "", UserId: admin, App: app,}
	organisationUtil.GetOrganisation()

	teamUtil := utils.TeamUtil{TeamId: teamParam, TeamName: "", UserId: admin, App: app, Organisation: organisationUtil}

	member := teamUtil.RemoveMember(userParam); if member == &utils.AndFalse {
		return ctx.JSON(http.StatusInternalServerError, types.TeamOperation {
			Message: fmt.Sprintf("Team Member Not Removed"),
			Status: 500,
			Error: true,
			Type: "RemoveTeamMember",
			State: "Unsuccessful",
		})
	}

	return ctx.JSON(http.StatusOK, types.TeamOperation {
		Message: fmt.Sprintf("Team Member Not Removed"),
		Status: 200,
		Error: true,
		Type: "RemoveTeamMember",
		State: "Success",
	})

}