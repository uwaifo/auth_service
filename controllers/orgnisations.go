package controllers

import (
	"app-auth/iam"
	"app-auth/types"
	"app-auth/utils"
	"time"

	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/mongodb/mongo-go-driver/bson"
	uuid "github.com/satori/go.uuid"
)

func GetOrganisations(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	admin := claims["id"].(string)
	email := claims["email"].(string)
	app := claims["app"].(string)

	userUtil := utils.UserMemberUtil{UserId: &admin, UserEmail: &email, App: app}

	//userOrganisations := userUtil.GetUserOrganisations()
	userOrganisations := userUtil.GetUserActiveOrganisations()

	if userOrganisations == nil {
		return ctx.JSON(http.StatusNotFound, types.OrganisationNotFound{
			Message: "Organisations Not Found",
			Status:  404,
			Error:   true,
			Type:    "GetOrganisations",
		})
	}

	return ctx.JSON(http.StatusOK, *userOrganisations)
}

func GetDeletedOrganisations(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	admin := claims["id"].(string)
	email := claims["email"].(string)
	app := claims["app"].(string)

	userUtil := utils.UserMemberUtil{UserId: &admin, UserEmail: &email, App: app}
	userOrganisations := userUtil.GetUserInactiveOrganisations()

	if userOrganisations == nil {
		return ctx.JSON(http.StatusNotFound, types.OrganisationNotFound{
			Message: "Organisations Not Found",
			Status:  404,
			Error:   true,
			Type:    "GetOrganisations",
		})
	}

	return ctx.JSON(http.StatusOK, *userOrganisations)
}
func GetDeletedOrganisationsCount(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	admin := claims["id"].(string)
	email := claims["email"].(string)
	app := claims["app"].(string)

	userUtil := utils.UserMemberUtil{UserId: &admin, UserEmail: &email, App: app}
	userDeletedOrganizationsCount := userUtil.GetUserInactiveOrganizationsCount()

	if userDeletedOrganizationsCount == nil {
		return ctx.JSON(http.StatusNotFound, types.OrganisationNotFound{
			Message: "Organisations Not Found",
			Status:  404,
			Error:   true,
			Type:    "GetOrganisations",
		})
	}

	return ctx.JSON(http.StatusOK, *userDeletedOrganizationsCount)
}
func GetOrganisation(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	admin := claims["id"].(string)
	app := claims["app"].(string)
	id := ctx.Param("id")
	organisationUtil := utils.OrganisationUtil{Id: id, Name: "", UserId: admin, App: app}

	organisation := organisationUtil.GetOrganisation()
	if organisation == nil {
		return ctx.JSON(http.StatusNotFound, types.OrganisationNotFound{
			Message: fmt.Sprintf("Organisation %s Not Found", id),
			Status:  404,
			Error:   true,
			Type:    "GetOrganisation",
		})
	}

	return ctx.JSON(http.StatusOK, organisation)
}

func PostOrganisation(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	admin := claims["id"].(string)
	app := claims["app"].(string)

	organisationPostDetails := new(types.OrganisationPOSTPayload)
	err := ctx.Bind(organisationPostDetails)
	if err != nil {
		fmt.Println(err)
		return err
	}

	organisationId := uuid.NewV4().String()
	organisationUtil := utils.OrganisationUtil{Id: uuid.NewV4().String(), Name: organisationPostDetails.Name, UserId: admin, App: app}

	organisation := organisationUtil.CreateOrganisation(organisationPostDetails.ImageUrl)
	if organisation == nil {
		return ctx.JSON(http.StatusInternalServerError, types.OrganisationOperation{
			Message: fmt.Sprintf("Organisation %s Not Created", organisationId),
			Status:  500,
			Error:   true,
			Type:    "PostOrganisation",
			State:   "Unsuccessful",
		})
	}

	return ctx.JSON(http.StatusCreated, organisation)
}

func PutOrganisation(ctx echo.Context) error {
	permission := "update.organisation"
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	admin := claims["id"].(string)
	email := claims["email"].(string)
	app := claims["app"].(string)
	id := ctx.Param("id")

	// before anything else, we have to make sure the user has the correct permissions to perform this action
	// use the organisation util to find the organisation
	// and update it
	organisationUtil := utils.OrganisationUtil{Id: id, Name: "", UserId: admin, App: app}

	// we get the user scopes within this organisation
	userOrganisationScopes := iam.FindUserScopes{Id: admin, App: app, Email: email, OrganisationId: id, TeamId: ""}

	// check if the user has permissions to perform this action
	hasPermission := organisationUtil.HasOrganisationPermission(userOrganisationScopes, permission)

	// if the team is not found
	if hasPermission == 404 {
		return ctx.JSON(http.StatusNotFound, types.OrganisationOperation{
			Message: fmt.Sprintf("Organisation %s Not Found", id),
			Status:  404,
			Error:   true,
			Type:    "PutOrganisation",
			State:   "Unsuccessful",
		})
	}

	// the user does not have a permission to perform this task
	if hasPermission == 401 {
		return ctx.JSON(401, types.OrganisationOperation{
			Message: fmt.Sprintf("Need a %s permission to perform this action", permission),
			Status:  401,
			Error:   false,
			Type:    "PutOrganisation",
			State:   "Unsuccessful",
		})
	}

	organisationPutDetails := new(types.OrganisationPUTPayload)
	err := ctx.Bind(organisationPutDetails)
	if err != nil {
		fmt.Println(err)
		return err
	}

	updateData := bson.D{{"name", organisationPutDetails.Name}, {"imageurl", organisationPutDetails.ImageUrl}}

	organisation := organisationUtil.UpdateOrganisation(updateData)
	if organisation == nil {
		return ctx.JSON(http.StatusInternalServerError, types.OrganisationOperation{
			Message: fmt.Sprintf("Organisation %s Not Updated", id),
			Status:  500,
			Error:   true,
			Type:    "PutOrganisation",
			State:   "Unsuccessful",
		})
	}

	return ctx.JSON(http.StatusOK, organisation)
}

func DeleteOrganisation(ctx echo.Context) error {
	permission := "update.organisation"
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	admin := claims["id"].(string)
	email := claims["email"].(string)
	app := claims["app"].(string)
	id := ctx.Param("id")

	// before anything else, we have to make sure the user has the correct permissions to perform this action
	// use the organisation util to find the organisation
	// and update it
	organisationUtil := utils.OrganisationUtil{Id: id, Name: "", UserId: admin, App: app}

	// we get the user scopes within this organisation
	userOrganisationScopes := iam.FindUserScopes{Id: admin, App: app, Email: email, OrganisationId: id, TeamId: ""}

	// check if the user has permissions to perform this action
	hasPermission := organisationUtil.HasOrganisationPermission(userOrganisationScopes, permission)

	// if the team is not found
	if hasPermission == 404 {
		return ctx.JSON(http.StatusNotFound, types.OrganisationOperation{
			Message: fmt.Sprintf("Organisation %s Not Found", id),
			Status:  404,
			Error:   true,
			Type:    "PutOrganisation",
			State:   "Unsuccessful",
		})
	}

	// the user does not have a permission to perform this task
	if hasPermission == 401 {
		return ctx.JSON(401, types.OrganisationOperation{
			Message: fmt.Sprintf("Need a %s permission to perform this action", permission),
			Status:  401,
			Error:   false,
			Type:    "PutOrganisation",
			State:   "Unsuccessful",
		})
	}

	organisation := organisationUtil.GetOrganisation()
	if organisation.Deleted {
		return ctx.JSON(http.StatusNotFound, types.OrganisationNotFound{
			Message: fmt.Sprintf("Organisation %s with ID %s has been deleted already.", organisation.Name, id),
			Status:  404,
			Error:   true,
			Type:    "GetOrganisation",
		})
	}

	deleteData := bson.D{{"deleted", true} /*, {"deleted_at", organisation.DeletedAt}*/}

	deleteOrganisation := organisationUtil.RemoveOrganisation(deleteData)
	if deleteOrganisation == nil {
		return ctx.JSON(http.StatusInternalServerError, types.OrganisationOperation{
			Message: fmt.Sprintf("Organisation %s Not Updated", id),
			Status:  500,
			Error:   true,
			Type:    "PutOrganisation",
			State:   "Unsuccessful",
		})
	}
	return ctx.JSON(http.StatusAccepted, types.OrganisationOperation{
		Message: fmt.Sprintf("Organisation %s Removed", id),
		Status:  202,
		Error:   false,
		Type:    "DeleteOrganisation",
		State:   "Success",
	})

}

func RestoreOrganisation(ctx echo.Context) error {
	permission := "update.organisation"
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	admin := claims["id"].(string)
	email := claims["email"].(string)
	app := claims["app"].(string)
	id := ctx.Param("id")

	// before anything else, we have to make sure the user has the correct permissions to perform this action
	// use the organisation util to find the organisation
	// and update it
	organisationUtil := utils.OrganisationUtil{Id: id, Name: "", UserId: admin, App: app}

	// we get the user scopes within this organisation
	userOrganisationScopes := iam.FindUserScopes{Id: admin, App: app, Email: email, OrganisationId: id, TeamId: ""}

	// check if the user has permissions to perform this action
	hasPermission := organisationUtil.HasOrganisationPermission(userOrganisationScopes, permission)

	// if the team is not found
	if hasPermission == 404 {
		return ctx.JSON(http.StatusNotFound, types.OrganisationOperation{
			Message: fmt.Sprintf("Organisation %s Not Found", id),
			Status:  404,
			Error:   true,
			Type:    "PutOrganisation",
			State:   "Unsuccessful",
		})
	}

	// the user does not have a permission to perform this task
	if hasPermission == 401 {
		return ctx.JSON(401, types.OrganisationOperation{
			Message: fmt.Sprintf("Need a %s permission to perform this action", permission),
			Status:  401,
			Error:   false,
			Type:    "PutOrganisation",
			State:   "Unsuccessful",
		})
	}

	organisation := organisationUtil.GetOrganisation()
	if !organisation.Deleted {
		return ctx.JSON(http.StatusNotFound, types.OrganisationNotFound{
			Message: fmt.Sprintf("Organisation %s with ID %s is already active.", organisation.Name, id),
			Status:  404,
			Error:   true,
			Type:    "GetOrganisation",
		})
	}

	deleteData := bson.D{{"deleted", false} /*, {"deleted_at", organisation.DeletedAt}*/}
	//fmt.Println(deleteData)

	deleteOrganisation := organisationUtil.RestoreDeletedOrganisation(deleteData)
	if deleteOrganisation == nil {
		return ctx.JSON(http.StatusInternalServerError, types.OrganisationOperation{
			Message: fmt.Sprintf("Organisation %s Not Updated", id),
			Status:  500,
			Error:   true,
			Type:    "PutOrganisation",
			State:   "Unsuccessful",
		})
	}
	return ctx.JSON(http.StatusAccepted, types.OrganisationOperation{
		Message: fmt.Sprintf("Organisation %s Restored", id),
		Status:  202,
		Error:   false,
		Type:    "DeleteOrganisation",
		State:   "Success",
	})

}

func DestroyOrganisation(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	admin := claims["id"].(string)
	app := claims["app"].(string)
	id := ctx.Param("id")
	organisationUtil := utils.OrganisationUtil{Id: id, Name: "", UserId: admin, App: app}

	destroyedOrganisation, err := organisationUtil.DestroyOrganisationRecord()
	if !destroyedOrganisation && err != nil {
		return ctx.JSON(http.StatusNotFound, types.OrganisationNotFound{
			Message: fmt.Sprintf("Organisation %s Not Found", id),
			Status:  404,
			Error:   true,
			Type:    "GetOrganisation",
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

/*
func DeleteOrganisation(ctx echo.Context) error {
	permission := "delete.organisation"
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	admin := claims["id"].(string)
	email := claims["email"].(string)
	app := claims["app"].(string)
	id := ctx.Param("id")

	organisationUtil := utils.OrganisationUtil{Id: id, Name: "", UserId: admin, App: app}

	// we get the user scopes within this organisation
	userOrganisationScopes := iam.FindUserScopes{Id: admin, App: app, Email: email, OrganisationId: id, TeamId: ""}

	// check if the user has permissions to perform this action
	hasPermission := organisationUtil.HasOrganisationPermission(userOrganisationScopes, permission)

	// if the team is not found
	if hasPermission == 404 {
		return ctx.JSON(http.StatusNotFound, types.OrganisationOperation{
			Message: fmt.Sprintf("Organisation %s Not Found", id),
			Status:  404,
			Error:   true,
			Type:    "PutOrganisation",
			State:   "Unsuccessful",
		})
	}

	// the user does not have a permission to perform this task
	if hasPermission == 401 {
		return ctx.JSON(401, types.OrganisationOperation{
			Message: fmt.Sprintf("Need a %s permission to perform this action", permission),
			Status:  401,
			Error:   false,
			Type:    "PutOrganisation",
			State:   "Unsuccessful",
		})
	}

	removedOrganisation := organisationUtil.RemoveOrganisation()
	if *removedOrganisation == false {
		return ctx.JSON(http.StatusInternalServerError, types.OrganisationOperation{
			Message: fmt.Sprintf("Organisation %s Not Removed", id),
			Status:  500,
			Error:   true,
			Type:    "DeleteOrganisation",
			State:   "Unsuccessful",
		})
	}

	return ctx.JSON(http.StatusOK, types.OrganisationOperation{
		Message: fmt.Sprintf("Organisation %s Removed", id),
		Status:  200,
		Error:   false,
		Type:    "DeleteOrganisation",
		State:   "Success",
	})
}
*/
