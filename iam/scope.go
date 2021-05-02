package iam

import (
	"app-auth/db"
	"app-auth/types"
	"context"
	"github.com/mongodb/mongo-go-driver/bson"
	"log"
	"strings"
)

type Scope struct {
	Scopes []string `json:"scopes"`
	App string `json:"app"`

	scopes []types.Scopes
	permissions []string
}

// toString() representation of the
func (scope Scope) String() string  {
	return "<Scope " + strings.Join(scope.Scopes, " ") + " />"
}

// ToApp() returns the name of
func (scope Scope) GetApp() string  {
	return strings.ToUpper(scope.App)
}

func (scope Scope) getScopesFromDB() []types.Scopes {

	// ["EDITOR", "READER"]

	scopesFilter := bson.D{
		{"name", bson.D{{"$in", scope.Scopes}}},
		{"$or", bson.A{
			bson.D{{"app", scope.App}},
			bson.D{{"app", "id.scaratec.com"}},
		}},
	}

	cur, err := db.ScopesCollection.Find(context.TODO(), scopesFilter);  if err != nil { log.Println(err) }

	var scopes []types.Scopes

	defer cur.Close(context.TODO())

	// loop through all the
	for cur.Next(context.Background()) {
		var elem types.Scopes
		err := cur.Decode(&elem)
		if err != nil { log.Print(err) }
		scopes = append(scopes, elem)
	}

	// try avoiding multiple database calls by setting
	// this value as a struct variable
	scope.scopes = scopes

	return scopes
}

func (scope Scope) getPermissionList() []string {
	var scopes = scope.getScopesFromDB()
	var permissions []string

	for _, scp := range scopes {
		permissions = append(permissions, scp.Permissions...)
	}

	scope.permissions = permissions

	return permissions
}

func (scope Scope) permission() Permission {
	// this function call gets all the required scopes from the database
	var _ = scope.getPermissionList()
	return Permission{
		Scopes: scope.Scopes,
		App: scope.App,
		Permissions: scope.permissions,
	}
}

// when given a specific, return if the Scope has this specific permission
func (scope Scope) HasPermission(permit string) bool {
	return scope.permission().HasPermission(permit)
}




// this util function is for finding user scopes
type FindUserScopes struct {
	Id string `json:"id"`
	App string `json:"app"`
	Email string `json:"email"`
	OrganisationId string `json:"organisation_id"`
	TeamId string `json:"team_id"`
}

// filter for returning scopes at the organisation level
func (findUserScopes FindUserScopes) GetOrganisationsScopesLevelFilter() bson.D {
	return bson.D {
		{"organisationid", findUserScopes.OrganisationId},
		{"id", findUserScopes.Id},
		{"useremail", findUserScopes.Email},
		{"app", findUserScopes.App},
	}
}


// filter for returing scopes at the team level
func (findUserScopes FindUserScopes) GetTeamsScopesLevelFilter() bson.D {
	return bson.D {
		{"organisationid", findUserScopes.OrganisationId},
		{"teamid", findUserScopes.TeamId},
		{"id", findUserScopes.Id},
		{"useremail", findUserScopes.Email},
		{"app", findUserScopes.App},
	}
}


// get app level scopes. for setting and updating IAMS
// filter for returing scopes at the team level
func (findUserScopes FindUserScopes) GetAppScopesLevelFilter() bson.D {
	return bson.D {
		{"id", findUserScopes.Id},
		{"useremail", findUserScopes.Email},
		{"app", findUserScopes.App},
	}
}


func (findUserScopes FindUserScopes) FindAndReturnUserScopesBasedOnFilter(filter bson.D) []string {
	// get a list of user scopes based filter
	cur, err := db.UserScopesCollection.Find(context.TODO(), filter); if err != nil {
		log.Println(err)
		return []string{}
	}

	defer cur.Close(context.TODO())

	var userMemberScopes []string

	for cur.Next(context.TODO()) {
		var userMemberScope types.UserMemberScope
		err := cur.Decode(&userMemberScope)
		if err != nil { log.Print(err) }

		// append the scopes of the user to the upper string
		userMemberScopes = append(userMemberScopes, userMemberScope.Scopes...)
	}

	return userMemberScopes
}
