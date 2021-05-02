package controllers

import (
	"app-auth/iam"
	"app-auth/types"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/mongodb/mongo-go-driver/bson"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func GetIAMScopes(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	app := claims["app"].(string)

	iamScopes := iam.ScopeActions { Id: "", App: app, Name: "", Permission: []string{} }

	allScopes := iamScopes.GetScopes(); if allScopes == nil {
		return ctx.JSON(http.StatusNotFound, types.ScopeNotFound{
			Message: "Scopes Not Found",
			Status: 404,
			Error: true,
			Type: "GetIAMScopes",
		})
	}

	return ctx.JSON(http.StatusOK, allScopes)
}

func GetIAMScope(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	app := claims["app"].(string)
	id := ctx.Param("id")

	iamScopes := iam.ScopeActions { Id: id, App: app, Name: "", Permission: []string{} }

	iamScope := iamScopes.GetScope(); if iamScope == nil {
		return ctx.JSON(http.StatusNotFound, types.ScopeNotFound{
			Message: "Scope Not Found",
			Status: 404,
			Error: true,
			Type: "GetIAMScope",
		})
	}

	return ctx.JSON(http.StatusOK, iamScopes)
}

func PutIAMScope(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	app := claims["app"].(string)
	id := ctx.Param("id")

	iamScopePutDetails := new(types.IAMScopePUTPayload)
	err := ctx.Bind(iamScopePutDetails); if err != nil {
		fmt.Println(err)
		return err
	}

	iamScopes := iam.ScopeActions { Id: id, App: app, Name: "", Permission: []string{}}

	updateData := bson.D{{"name", iamScopePutDetails.Scope}}
	iamScope := iamScopes.UpdateScope(updateData); if iamScope == nil {
		return ctx.JSON(http.StatusInternalServerError, types.ScopeOperation{
			Message: "Scope Not Updated",
			Status: 500,
			Error: true,
			Type: "PutIAMScope",
			State: "Unsuccessful",
		})
	}

	return ctx.JSON(http.StatusOK, iamScope)
}

func PostIAMScope(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	app := claims["app"].(string)

	iamScopePostDetails := new(types.IAMScopePOSTPayload)
	err := ctx.Bind(iamScopePostDetails); if err != nil {
		fmt.Println(err)
		return err
	}

	iamScopeId := uuid.NewV4().String()
	iamScopes := iam.ScopeActions { Id: iamScopeId, App: app, Name: iamScopePostDetails.Scope, Permission: iamScopePostDetails.Permissions }

	iamScope := iamScopes.CreateScope(); if iamScope == nil {
		return ctx.JSON(http.StatusInternalServerError, types.ScopeOperation{
			Message: "Scope Not Created",
			Status: 500,
			Error: true,
			Type: "PostIAMScope",
			State: "Unsuccessful",
		})
	}

	return ctx.JSON(http.StatusOK, iamScope)
}

func DeleteIAMScope(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	app := claims["app"].(string)
	id := ctx.Param("id")

	iamScopes := iam.ScopeActions { Id: id, App: app, Name: "", Permission: []string{}}

	iamScope := iamScopes.RemoveScope(); if iamScope == nil {
		return ctx.JSON(http.StatusInternalServerError, types.ScopeOperation{
			Message: "Scope Not Updated",
			Status: 500,
			Error: true,
			Type: "PutIAMScope",
			State: "Unsuccessful",
		})
	}

	return ctx.JSON(http.StatusOK, iamScope)
}

func PostIAMScopesPermission(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	app := claims["app"].(string)
	id := ctx.Param("id")

	iamScopePostPermissionDetails := new(types.IAMScopePOSTPermissionPayload)
	err := ctx.Bind(iamScopePostPermissionDetails); if err != nil {
		fmt.Println(err)
		return err
	}

	iamScopes := iam.ScopeActions { Id: id, App: app, Name: "", Permission: []string{}}

	iamScope := iamScopes.AddScopePermission(iamScopePostPermissionDetails.Permission); if iamScope == nil {
		return ctx.JSON(http.StatusInternalServerError, types.ScopeOperation{
			Message: "Scope Permission Not Added",
			Status: 500,
			Error: true,
			Type: "PostIAMScopesPermission",
			State: "Unsuccessful",
		})
	}

	return ctx.JSON(http.StatusOK, iamScope)
}

func PostIAMScopesPermissions(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	app := claims["app"].(string)
	id := ctx.Param("id")

	iamScopePostPermissionsDetails := new(types.IAMScopePOSTPermissionsPayload)
	err := ctx.Bind(iamScopePostPermissionsDetails); if err != nil {
		fmt.Println(err)
		return err
	}

	iamScopes := iam.ScopeActions { Id: id, App: app, Name: "", Permission: []string{}}

	iamScope := iamScopes.AddScopePermissions(iamScopePostPermissionsDetails.Permissions); if iamScope == nil {
		return ctx.JSON(http.StatusInternalServerError, types.ScopeOperation {
			Message: "Scope Permissions Not Added",
			Status: 500,
			Error: true,
			Type: "PostIAMScopesPermissions",
			State: "Unsuccessful",
		})
	}

	return ctx.JSON(http.StatusOK, iamScope)
}

func DeleteIAMScopesPermission(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	app := claims["app"].(string)
	id := ctx.Param("id")

	iamScopeDeletePermissionDetails := new(types.IAMScopeDELETEPermissionPayload)
	err := ctx.Bind(iamScopeDeletePermissionDetails); if err != nil {
		fmt.Println(err)
		return err
	}

	iamScopes := iam.ScopeActions { Id: id, App: app, Name: "", Permission: []string{}}

	iamScope := iamScopes.RemoveScopePermission(iamScopeDeletePermissionDetails.Permission); if iamScope == nil {
		return ctx.JSON(http.StatusInternalServerError, types.ScopeOperation {
			Message: "Scope Permission Not Removed",
			Status: 500,
			Error: true,
			Type: "DeleteIAMScopesPermission",
			State: "Unsuccessful",
		})
	}

	return ctx.JSON(http.StatusOK, iamScope)
}
