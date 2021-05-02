package controllers

import (
	"app-auth/types"

	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

func AddTeamMemberScope(ctx echo.Context) error  {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	app := claims["app"].(string)
	id := ctx.Param("id")
	teamParam := ctx.Param("team")
	userId := ctx.Param("userId")

	memberScopePostDetails := new(types.MemberScopePUTPayload)
	err := ctx.Bind(memberScopePostDetails); if err != nil {
		fmt.Println(err)
		return err
	}

	userMemberScope := types.UserMemberScope {
		Id: userId,
		OrganisationId: id,
		TeamId: teamParam,
		UserId: userId,
		App: app,
	}

	memberScope := userMemberScope.AddUserMemberScopeScopes(memberScopePostDetails.Scope); if memberScope == nil {
		return ctx.JSON(http.StatusInternalServerError, types.TeamOperation {
			Message: fmt.Sprintf("Member Scope %s Not Added", userId),
			Status: 400,
			Error: true,
			Type: "AddTeamMemberScope",
			State: "Unsuccessful",
		})
	}

	return ctx.JSON(http.StatusCreated, memberScope)
}

func RemoveTeamMemberScope(ctx echo.Context) error  {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	app := claims["app"].(string)
	id := ctx.Param("id")
	teamParam := ctx.Param("team")
	userId := ctx.Param("userId")

	memberScopePostDetails := new(types.MemberScopePUTPayload)
	err := ctx.Bind(memberScopePostDetails); if err != nil {
		fmt.Println(err)
		return err
	}

	userMemberScope := types.UserMemberScope {
		Id: userId,
		OrganisationId: id,
		TeamId: teamParam,
		UserId: userId,
		App: app,
	}

	memberScope := userMemberScope.RemoveUserMemberScopeScopes(memberScopePostDetails.Scope); if memberScope == nil {
		return ctx.JSON(http.StatusInternalServerError, types.TeamOperation {
			Message: fmt.Sprintf("Member Scope %s Not Removed", userId),
			Status: 400,
			Error: true,
			Type: "RemoveTeamMemberScope",
			State: "Unsuccessful",
		})
	}

	return ctx.JSON(http.StatusOK, memberScope)
}
