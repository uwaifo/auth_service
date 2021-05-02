package types

import (
	"app-auth/db"

	"context"
	"fmt"
	"log"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

type TeamMapType = map[string]interface{}
type OrganisationMapType = map[string]interface{}
type ScopeObjectMapping = map[string]OrganisationMapType

type Scopes struct {
	App         string
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

type ScopeAlreadyExists struct {
	Name    string `json:"name"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type ScopeNotFound struct {
	Type    string `json:"type"`
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   bool   `json:"error"`
}

type ScopeOperation struct {
	Type    string `json:"type"`
	Status  int    `json:"status"`
	Message string `json:"message"`
	State   string `json:"state"`
	Error   bool   `json:"error"`
}

type ScopeTeam struct {
	Name  string `json:"name"`
	Id    string `json:"id"`
	Scope string `json:"scope"`
}

type ScopeOrganisation struct {
	Name  string      `json:"name"`
	Id    string      `json:"id"`
	Admin bool        `json:"admin"`
	Team  []ScopeTeam `json:"team"`
}

// experimental searches for list fails with longer search times...
type ScopeObjectList struct {
	Scopes []ScopeOrganisation
}

type ScopeObject struct {
	Scopes ScopeObjectMapping
}

type UserMemberScope struct {
	Id               string   `json:"id"`
	OrganisationName string   `json:"organisation_name"`
	OrganisationId   string   `json:"organisation_id"`
	TeamName         string   `json:"team_name"`
	TeamId           string   `json:"team_id"`
	UserEmail        string   `json:"user_email"`
	UserId           string   `json:"user_id"`
	App              string   `json:"app"`   // this state would correspond to a user invite
	State            string   `json:"state"` // this state would correspond to a user invite
	Scopes           []string `json:"scopes"`
}

func (userMemberScopes UserMemberScope) String() string {
	return fmt.Sprintf(
		`<UserMemberScopes email="%s" id="%s" team_id="%s" organisation_id="%" />`,
		userMemberScopes.UserEmail,
		userMemberScopes.Id,
		userMemberScopes.TeamId,
		userMemberScopes.OrganisationId,
	)
}

func (userMemberScopes UserMemberScope) getScopeFilter() bson.D {
	return bson.D{
		{"id", userMemberScopes.Id},
		{"app", userMemberScopes.App},
		{"userid", userMemberScopes.UserId},
		{"teamid", userMemberScopes.TeamId},
		{"organisationid", userMemberScopes.OrganisationId},
	}
}

func (userMemberScopes UserMemberScope) GetUserMemberScope() *UserMemberScope {
	userMemberScopesData := UserMemberScope{}

	err := db.UserScopesCollection.FindOne(context.Background(), userMemberScopes.getScopeFilter()).Decode(&userMemberScopesData)
	if err != nil {
		log.Println(err)
		return nil
	}

	return &userMemberScopesData
}

func (userMemberScopes UserMemberScope) SaveUserMemberScope() *UserMemberScope {
	userScopes := userMemberScopes.GetUserMemberScope()
	if userScopes != nil {
		return nil
	}

	_, err := db.UserScopesCollection.InsertOne(context.Background(), userMemberScopes)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &userMemberScopes
}

func (userMemberScopes UserMemberScope) RemoveUserMemberScope() *bool {
	userScopes := userMemberScopes.GetUserMemberScope()
	if userScopes == nil {
		return nil
	}

	_, err := db.UserScopesCollection.DeleteOne(context.Background(), userMemberScopes.getScopeFilter())
	if err != nil {
		log.Println(err)
		return &AndFalse
	}
	return &AndTrue
}

func (userMemberScopes UserMemberScope) RemoveUserMembersScope(userScopeIds []string) *bool {
	filter := bson.D{{
		"id", bson.D{{"$in", userScopeIds}},
	}}

	_, err := db.UserScopesCollection.DeleteMany(context.Background(), filter)
	if err != nil {
		log.Println(err)
		return &AndFalse
	}
	return &AndTrue
}

func (userMemberScopes UserMemberScope) UpdateUserMemberScope(userScopeUpdateInfo bson.D) *UserMemberScope {

	userScopes := userMemberScopes.GetUserMemberScope()
	if userScopes == nil {
		return nil
	}

	userScopesData := UserMemberScope{}
	userScopesUpdateData := bson.D{{"$set", userScopeUpdateInfo}}

	tr := true
	rd := options.After
	opts := &options.FindOneAndUpdateOptions{Upsert: &tr, BypassDocumentValidation: &tr, ReturnDocument: &rd}

	err := db.UserScopesCollection.FindOneAndUpdate(
		context.Background(),
		userMemberScopes.getScopeFilter(),
		userScopesUpdateData,
		opts,
	).Decode(&userScopesData)
	if err != nil {
		log.Println(err)
		return nil
	}

	return &userScopesData
}

func (userMemberScopes UserMemberScope) AddUserMemberScopeScopes(userScopeData string) *UserMemberScope {
	userScope := userMemberScopes.GetUserMemberScope()
	if userScope == nil {
		return nil
	}
	return userMemberScopes.UpdateUserMemberScope(bson.D{{"scopes", append(userScope.Scopes, userScopeData)}})
}

func (userMemberScopes UserMemberScope) RemoveUserMemberScopeScopes(userScopeData string) *UserMemberScope {
	// this scope doesn't exists, return nil
	userScope := userMemberScopes.GetUserMemberScope()
	if userScope == nil {
		return nil
	}

	// remove the permission from the list of permission
	index := FindIndex(userScope.Scopes, userScopeData)
	newUserScopes := RemoveIndex(userScope.Scopes, index)

	return userMemberScopes.UpdateUserMemberScope(bson.D{{"permissions", newUserScopes}})
}

var data = `
{
  "StandardClaims": {
    "exp": 1585989679,
    "iat": 1585988779,
    "iss": "Clipsynphony"
  },
  "age": 21,
  "email": "user_two.test@clipsynphony.com",
  "expires": "2020-04-04T08:41:19.520330457Z",
  "firstname": "Test",
  "id": "987ff34b-e68b-4ff2-8f2c-804a03209315",
  "lastname": "User",
  "provider": "local",
  "username": "test_user",
  "verified": true
  "scopes": {
    "987ff34b-e68b-4ff2-8f2c-804a03209315": {
		 "id": "987ff34b-e68b-4ff2-8f2c-804a03209315",
		 "name": "CNN"
		 "987ff34b-e68b-4ff2-8f2c-804a03209315": {
			"name": "CNN AFRICA"
			"id": "987ff34b-e68b-4ff2-8f2c-804a03209315"
			"scope": ["EDITOR"]
		  }
   },
  }


   Example:
   "987ff34b-e68b-4ff2-8f2c-804a03209315": {
	 "id": "987ff34b-e68b-4ff2-8f2c-804a03209315",
     "name": "CNN"
	 "987ff34b-e68b-4ff2-8f2c-804a03209315": {
        "name": "CNN AFRICA"
      	"id": "987ff34b-e68b-4ff2-8f2c-804a03209315"
		"scope": ["EDITOR"]
      }
   }
  },
}
`
