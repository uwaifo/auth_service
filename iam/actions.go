package iam

import (
	"app-auth/db"
	"app-auth/types"
	"context"
	"log"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

type ScopeActions struct {
	Id         string   `json:"id"`
	App        string   `json:"app"`
	Name       string   `json:"name"`
	Permission []string `json:"permission"`

	permission []string
}

func (scopeActions ScopeActions) String() string {
	return "<Scope id='" + scopeActions.Id + "' name='" + scopeActions.Name + "'/>"
}

func (scopeActions ScopeActions) getScopeFilter() bson.D {
	return bson.D{{"id", scopeActions.Id}, {"app", scopeActions.App}}
}

func (scopeActions ScopeActions) GetScope() *types.Scopes {
	if scopeActions.Id == "" || scopeActions.Name == "" {
		return nil
	}

	scopesData := types.Scopes{}
	err := db.ScopesCollection.FindOne(context.Background(), scopeActions.getScopeFilter()).Decode(&scopesData)
	if err != nil {
		log.Println(err)
		return nil
	}

	scopeActions.permission = scopesData.Permissions

	return &scopesData
}

func (scopeActions ScopeActions) GetScopes() *[]types.Scopes {

	filter := bson.D{{"$or", bson.A{
		bson.D{{"app", scopeActions.App}},
		bson.D{{"app", "id.scaratec.com"}},
	}}}

	cur, err := db.ScopesCollection.Find(context.Background(), filter)
	if err != nil {
		log.Println(err)
		return nil
	}

	defer cur.Close(context.Background())
	var scopes []types.Scopes

	for cur.Next(context.TODO()) {
		var scope types.Scopes
		err := cur.Decode(&scope)
		if err != nil {
			log.Print(err)
		}

		// append the teams data to the various array variables
		scopes = append(scopes, scope)
	}

	return &scopes
}

func (scopeActions ScopeActions) GetScopesByEmail(email string) *[]types.UserMemberScope {
	filter := bson.D{
		{"useremail", email},
	}

	cur, err := db.UserScopesCollection.Find(context.Background(), filter)
	if err != nil {
		log.Println(err)
		return nil
	}

	defer cur.Close(context.Background())
	var scopes []types.UserMemberScope

	for cur.Next(context.TODO()) {
		var scope types.UserMemberScope
		err := cur.Decode(&scope)
		if err != nil {
			log.Print(err)
		}

		// append the teams data to the various array variables
		scopes = append(scopes, scope)
	}

	return &scopes
}

func (scopeActions ScopeActions) CreateScope() *types.Scopes {
	// this scopes already exists
	scope := scopeActions.GetScope()
	if scope != nil {
		return nil
	}

	data := types.Scopes{
		Id:          scopeActions.Id,
		App:         scopeActions.App,
		Name:        scopeActions.Name,
		Permissions: scopeActions.Permission,
	}

	_, err := db.ScopesCollection.InsertOne(context.Background(), data)
	if err != nil {
		log.Println(err)
		return nil
	}

	return scopeActions.GetScope()
}

func (scopeActions ScopeActions) UpdateScope(data bson.D) *types.Scopes {
	// this scope doesn't exists, return nil
	scope := scopeActions.GetScope()
	if scope == nil {
		return nil
	}

	scopeData := types.Scopes{}
	scopeUpdateData := bson.D{{"$set", data}}

	tr := true
	rd := options.After
	opts := &options.FindOneAndUpdateOptions{Upsert: &tr, BypassDocumentValidation: &tr, ReturnDocument: &rd}

	err := db.ScopesCollection.FindOneAndUpdate(context.Background(), scopeActions.getScopeFilter(), scopeUpdateData, opts).Decode(&scopeData)
	if err != nil {
		log.Println(err)
		return nil
	}

	return &scopeData
}

func (scopeActions ScopeActions) AddScopePermission(permission string) *types.Scopes {
	// this scope doesn't exists, return nil
	scope := scopeActions.GetScope()
	if scope == nil {
		return nil
	}
	permissions := append(scope.Permissions, permission)
	return scopeActions.UpdateScope(bson.D{{"permissions", permissions}})
}

func (scopeActions ScopeActions) AddScopePermissions(permission []string) *types.Scopes {
	// this scope doesn't exists, return nil
	scope := scopeActions.GetScope()
	if scope == nil {
		return nil
	}
	permissions := append(scope.Permissions, permission...)
	return scopeActions.UpdateScope(bson.D{{"permissions", permissions}})
}

func (scopeActions ScopeActions) RemoveScopePermission(permission string) *types.Scopes {
	// this scope doesn't exists, return nil
	scope := scopeActions.GetScope()
	if scope == nil {
		return nil
	}

	// remove the permission from the list of permission
	index := types.FindIndex(scope.Permissions, permission)
	newPermissions := types.RemoveIndex(scope.Permissions, index)

	return scopeActions.UpdateScope(bson.D{{"permissions", newPermissions}})
}

func (scopeActions ScopeActions) RemoveScope() *bool {
	// this scope doesn't exists, return nil
	scope := scopeActions.GetScope()
	if scope == nil {
		return nil
	}

	// then delete the scope
	_, err := db.ScopesCollection.DeleteOne(context.Background(), scopeActions.getScopeFilter())
	if err != nil {
		log.Println(err)
		return &types.AndFalse // return false if the organisation is not deleted
	}

	// return true when the
	return &types.AndTrue
}

// HasIAMScopePermissions() checks if a user has a permission to perform a specific task on iam scopes
func (scopeActions ScopeActions) HasIAMScopePermissions(userOrganisationScope FindUserScopes, permission string) int {

	// get scopes pertaining to the team level
	scopes := userOrganisationScope.FindAndReturnUserScopesBasedOnFilter(userOrganisationScope.GetAppScopesLevelFilter())

	// we initialize a scope checker, which would check the permissions
	scopesChecker := Scope{Scopes: scopes, App: userOrganisationScope.App}

	// check if the user has the required permission to perform this task
	if !scopesChecker.HasPermission(permission) {
		return 401
	}

	return 200
}
