/*
 * The main job of this functionality is to:

 * Find all the teams a user belongs in
 * Find all the organisations that have those teams
 * Get the scopes of the user within the teams, using the team's data
 */

package utils

import (
	"context"
	"log"
	"net/http"

	"app-auth/db"
	"app-auth/types"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

type UserMemberUtil struct {
	UserId    *string `json:"user_id"`
	UserEmail *string `json:"user_email"`
	App       string  `json:"app"`

	user   *types.UserData
	active bool

	teams   []types.Team
	teamIds []string

	organisations   []types.Organisation
	organisationIds []string

	scopes      []types.UserMemberScope
	scopeObject types.ScopeObjectMapping
}

// toString() representation of the
func (userUtil UserMemberUtil) String() string {
	return "<UserUtil user_id='" + *userUtil.UserId + "' user_email='" + *userUtil.UserEmail + "'/>"
}

// GetUserByEmail() returns the found user or nil
func (userUtil UserMemberUtil) GetUserByEmail() *types.UserData {

	// if there's no user init email
	if userUtil.UserEmail == nil || *userUtil.UserEmail == "" {
		return nil
	}

	userData := types.UserData{}
	filter := bson.D{{"email", *userUtil.UserEmail}, {"active", true}}
	err := db.UserCollection.FindOne(context.Background(), filter).Decode(&userData)
	if err != nil {
		log.Println(err)
		return nil
	}

	userUtil.UserId = &userData.Id
	userUtil.user = &userData
	return userUtil.user
}

// GetUserById() returns the found user or nil
func (userUtil UserMemberUtil) GetUserById() *types.UserData {
	// if there's no user init id
	if userUtil.UserId == nil || *userUtil.UserId == "" {
		return nil
	}

	log.Println(*userUtil.UserId)

	userData := types.UserData{}
	filter := bson.D{{"id", *userUtil.UserId}, {"active", true}}
	err := db.UserCollection.FindOne(context.Background(), filter).Decode(&userData)
	if err != nil {
		log.Println(err)
		return nil
	}

	userUtil.UserEmail = &userData.Email
	userUtil.user = &userData
	return userUtil.user
}

// CreateNewUser() returns a new user created
func (userUtil UserMemberUtil) CreateNewUser(userData types.UserData) *types.UserData {

	userUtil.UserEmail = &userData.Email
	userUtil.UserId = &userData.Id

	// there's no user data here, then we create and return the new user data
	data := userUtil.GetUserByEmail()
	if data == nil {
		res, err := db.UserCollection.InsertOne(context.Background(), userData)
		log.Println(err)
		log.Println(res)
		if err != nil {
			return nil
		}

		//fmt.Println(res)
		userUtil.user = userUtil.GetUserById()

		return userUtil.user
	}

	// else we just return the user data
	return data
}

// GetUserByEmail() returns the found user or nil
func (userUtil UserMemberUtil) GetUserByEmailTm() *types.UserData {

	// if there's no user init email
	if userUtil.UserEmail == nil || *userUtil.UserEmail == "" {
		return nil
	}

	userData := types.UserData{}
	filter := bson.D{{"email", *userUtil.UserEmail}}
	err := db.UserCollection.FindOne(context.Background(), filter).Decode(&userData)
	if err != nil {
		log.Println(err)
		return nil
	}

	userUtil.UserId = &userData.Id
	userUtil.user = &userData
	return userUtil.user
}

// GetUserById() returns the found user or nil
func (userUtil UserMemberUtil) GetUserByIdTm() *types.UserData {
	// if there's no user init id
	if userUtil.UserId == nil || *userUtil.UserId == "" {
		return nil
	}

	log.Println(*userUtil.UserId)

	userData := types.UserData{}
	filter := bson.D{{"id", *userUtil.UserId}}
	err := db.UserCollection.FindOne(context.Background(), filter).Decode(&userData)
	if err != nil {
		log.Println(err)
		return nil
	}

	userUtil.UserEmail = &userData.Email
	userUtil.user = &userData
	return userUtil.user
}

// CreateNewUserTm() returns a new user created for tiermedizin
func (userUtil UserMemberUtil) CreateNewUserTm(userData types.UserData) *types.UserData {

	userUtil.UserEmail = &userData.Email
	userUtil.UserId = &userData.Id

	// there's no user data here, then we create and return the new user data
	data := userUtil.GetUserByEmailTm()
	if data == nil {
		res, err := db.UserCollection.InsertOne(context.Background(), userData)
		log.Println(err)
		log.Println(res)
		if err != nil {
			return nil
		}

		//fmt.Println(res)
		userUtil.user = userUtil.GetUserByIdTm()

		return userUtil.user
	}

	// else we just return the user data
	return data
}

// UpdateUser() returns an update user in the db
func (userUtil UserMemberUtil) UpdateUser(userData bson.D) *types.UserData {

	// there's no user data here we return nil
	data := userUtil.GetUserById()
	if data == nil {
		return nil
	}

	returnUserData := types.UserData{}
	userUpdateData := bson.D{{"$set", userData}}
	filter := bson.D{{"id", userUtil.UserId}, {"active", true}}

	tr := true
	rd := options.After
	opts := &options.FindOneAndUpdateOptions{Upsert: &tr, BypassDocumentValidation: &tr, ReturnDocument: &rd}

	// else we find and update the user data
	err := db.UserCollection.FindOneAndUpdate(context.Background(), filter, userUpdateData, opts).Decode(&returnUserData)
	if err != nil {
		return nil
	}
	userUtil.user = &returnUserData
	return userUtil.user
}

// UpdateUser() returns an update user in the db
func (userUtil UserMemberUtil) UpdateUserTm(userData bson.D) *types.UserData {

	// there's no user data here we return nil
	data := userUtil.GetUserByIdTm()
	if data == nil {
		return nil
	}

	returnUserData := types.UserData{}
	userUpdateData := bson.D{{"$set", userData}}
	filter := bson.D{{"id", userUtil.UserId}}

	tr := true
	rd := options.After
	opts := &options.FindOneAndUpdateOptions{Upsert: &tr, BypassDocumentValidation: &tr, ReturnDocument: &rd}

	// else we find and update the user data
	err := db.UserCollection.FindOneAndUpdate(context.Background(), filter, userUpdateData, opts).Decode(&returnUserData)
	if err != nil {
		return nil
	}
	userUtil.user = &returnUserData
	return userUtil.user
}

// RemoveUser() returns a new user created
func (userUtil UserMemberUtil) RemoveUser() bool {

	userUpdateData := bson.D{{"active", false}}
	data := userUtil.UpdateUser(userUpdateData)
	if data == nil {
		return false
	}

	userUtil.active = data.Active
	return true
}

// GetTeams() returns all the teams that has a specific userId as a member
func (userUtil UserMemberUtil) GetUserTeamsAll() *[]types.Team {
	// if there's no userId in the struct
	if userUtil.UserId == nil || *userUtil.UserId == "" {
		return nil
	}

	// get all the teams where the creator matches the current userId
	// or where the current userId is a part of the members
	teamsFilter := bson.D{
		{"app", userUtil.App},
		{"$or", bson.A{
			// return any team you created
			bson.D{{"createdby", *userUtil.UserId}},
			// or are a member of
			bson.D{
				{"members", bson.D{
					{"$in", []string{*userUtil.UserId}},
				}},
			},
		}},
	}

	cur, err := db.TeamCollection.Find(context.Background(), teamsFilter)
	if err != nil {
		log.Println(err)
		return nil
	}

	defer cur.Close(context.TODO())

	var teams []types.Team
	var teamIds []string
	var organisationIds []string

	for cur.Next(context.TODO()) {
		var team types.Team
		err := cur.Decode(&team)
		if err != nil {
			log.Print(err)
		}

		// append the teams data to the various array variables
		teams = append(teams, team)
		teamIds = append(teamIds, team.Id)
		organisationIds = append(organisationIds, team.OrganisationId)
	}

	userUtil.teams = teams
	userUtil.teamIds = teamIds
	userUtil.organisationIds = organisationIds

	return &userUtil.teams
}

// GetTeams() returns all the teams that has a specific userId as a member
func (userUtil UserMemberUtil) GetUserTeams(deleted bool) *[]types.Team {
	// if there's no userId in the struct
	if userUtil.UserId == nil || *userUtil.UserId == "" {
		return nil
	}

	// get all the teams where the creator matches the current userId
	// or where the current userId is a part of the members
	teamsFilter := bson.D{
		{"app", userUtil.App},
		{"deleted", deleted},
		{"$or", bson.A{
			// return any team you created
			bson.D{{"createdby", *userUtil.UserId}},
			// or are a member of
			bson.D{
				{"members", bson.D{
					{"$in", []string{*userUtil.UserId}},
				}},
			},
		}},
	}

	cur, err := db.TeamCollection.Find(context.Background(), teamsFilter)
	if err != nil {
		log.Println(err)
		return nil
	}

	defer cur.Close(context.TODO())

	var teams []types.Team
	var teamIds []string
	var organisationIds []string

	for cur.Next(context.TODO()) {
		var team types.Team
		err := cur.Decode(&team)
		if err != nil {
			log.Print(err)
		}

		// append the teams data to the various array variables
		teams = append(teams, team)
		teamIds = append(teamIds, team.Id)
		organisationIds = append(organisationIds, team.OrganisationId)
	}

	userUtil.teams = teams
	userUtil.teamIds = teamIds
	userUtil.organisationIds = organisationIds

	return &userUtil.teams
}

func (userUtil UserMemberUtil) GetUserDeletedTeams() *[]types.Team {
	return userUtil.GetUserTeams(true)
}

func (userUtil UserMemberUtil) GetUserActiveTeams() *[]types.Team {
	return userUtil.GetUserTeams(false)
}

// GetListUserOrganisations returns a list of the users non-deleted organisations
func (userUtil UserMemberUtil) GetUserActiveOrganisations() *[]types.Organisation {
	// return all the active teams
	return userUtil.GetUserOrganisations(false)
}

// GetListUserOrganisations returns a list of the users non-deleted organisations
func (userUtil UserMemberUtil) GetUserInactiveOrganisations() *[]types.Organisation {
	// pass true to return a list of deleted teams
	return userUtil.GetUserOrganisations(true)
}

func (userUtil UserMemberUtil) GetUserInactiveOrganizationsCount() *types.DeletedOrganisationsCount {
	var teams = userUtil.GetUserTeamsAll()
	var organisatioFilter bson.D

	var organisationIds []string
	for _, team := range *teams {
		organisationIds = append(organisationIds, team.OrganisationId)
	}

	// filter all teams that have an ID in the organisationIds array
	organisatioFilter = bson.D{
		{"$or", bson.A{
			bson.D{{"id", bson.D{{"$in", organisationIds}}}},
			bson.D{{"createdby", *userUtil.UserId}},
		}},
		{"deleted", true},
	}

	cur, err := db.OrganisationCollection.CountDocuments(context.Background(), organisatioFilter)
	if err != nil {
		log.Println(err)
		return nil
	}

	//defer cur.Close(context.TODO())

	var organizationCount = types.DeletedOrganisationsCount{
		Type:   "",
		Count:  int(cur),
		Status: http.StatusOK,
		// Message: "",
		Error: false,
	}

	return &organizationCount

}

// GetOrganisations() returns all the organisation th user created or is a part of
func (userUtil UserMemberUtil) GetUserOrganisations(deleted bool) *[]types.Organisation {
	var teams = userUtil.GetUserTeamsAll()

	var organisatioFilter bson.D

	if teams == nil || len(*teams) == 0 {
		// get all the organisations created by this user
		organisatioFilter = bson.D{{"createdby", *userUtil.UserId}, {"deleted", deleted}}
	} else {
		var organisationIds []string
		for _, team := range *teams {
			organisationIds = append(organisationIds, team.OrganisationId)
		}

		// filter all teams that have an ID in the organisationIds array
		organisatioFilter = bson.D{
			{"$or", bson.A{
				bson.D{{"id", bson.D{{"$in", organisationIds}}}},
				bson.D{{"createdby", *userUtil.UserId}},
			}},
			{"deleted", deleted},
		}
	}

	cur, err := db.OrganisationCollection.Find(context.Background(), organisatioFilter)
	if err != nil {
		log.Println(err)
		return nil
	}

	defer cur.Close(context.TODO())

	var organisations []types.Organisation

	for cur.Next(context.TODO()) {
		var organisation types.Organisation
		err := cur.Decode(&organisation)
		if err != nil {
			log.Print(err)
		}

		// append the teams data to the various array variables
		organisations = append(organisations, organisation)
	}

	userUtil.organisations = organisations
	return &userUtil.organisations

}

// GetUserScopes() returns all the scopes of the current user.
func (userUtil UserMemberUtil) GetUserScopeWithID(teamid string) *types.UserMemberScope {

	// if there's no userId and there's not user email, then return nil
	if (userUtil.UserId == nil || *userUtil.UserId == "") && (userUtil.UserEmail == nil || *userUtil.UserEmail == "") {
		return nil
	}

	var userMemberScopesFilter = bson.D{}
	if *userUtil.UserId != "" {
		userMemberScopesFilter = bson.D{
			{"id", *userUtil.UserId},
			{"teamid", teamid},
		}
	} else {
		userMemberScopesFilter = bson.D{
			{"useremail", *userUtil.UserEmail},
			{"teamid", teamid},
		}
	}

	var userMemberScopeData = new(types.UserMemberScope)

	err := db.UserScopesCollection.FindOne(context.TODO(), userMemberScopesFilter).Decode(userMemberScopeData)
	if err != nil {
		log.Println(err)
		return nil
	}

	return userMemberScopeData
}

// GetUserScopes() returns all the scopes of the current user.
func (userUtil UserMemberUtil) GetUserScopes() *[]types.UserMemberScope {

	// if there's no userId and there's not user email, then return nil
	if (userUtil.UserId == nil || *userUtil.UserId == "") && (userUtil.UserEmail == nil || *userUtil.UserEmail == "") {
		return nil
	}

	var userMemberScopesFilter = bson.D{}
	if *userUtil.UserId != "" {
		userMemberScopesFilter = bson.D{
			{"id", *userUtil.UserId},
		}
	} else {
		userMemberScopesFilter = bson.D{
			{"useremail", *userUtil.UserEmail},
		}
	}

	cur, err := db.UserScopesCollection.Find(context.TODO(), userMemberScopesFilter)
	if err != nil {
		log.Println(err)
		return nil
	}

	defer cur.Close(context.TODO())

	var userMemberScopes []types.UserMemberScope

	for cur.Next(context.TODO()) {
		var userMemberScope types.UserMemberScope
		err := cur.Decode(&userMemberScope)
		if err != nil {
			log.Print(err)
		}

		// append the teams data to the various array variables
		userMemberScopes = append(userMemberScopes, userMemberScope)
	}

	userUtil.scopes = userMemberScopes
	return &userUtil.scopes
}

// GetUserScopeObject construct a mapping for the scope to used as a part for the jwt
func (userUtil UserMemberUtil) GetUserScopesObject() types.ScopeObjectMapping {
	var scopes = userUtil.GetUserScopes()
	if scopes == nil {
		return types.ScopeObjectMapping{}
	}
	var scopeObject = types.ScopeObjectMapping{}

	for _, scope := range *scopes {
		teamMap := types.TeamMapType{
			"id":     scope.TeamId,
			"name":   scope.TeamName,
			"scopes": scope.Scopes,
		}

		// check if the organisation is already created as a prop on the topLevel Object
		organisationMap := scopeObject[scope.OrganisationId]
		if organisationMap != nil {
			// if it exists, get the organisation and add a new team
			organisationMap[scope.TeamId] = teamMap
		} else {
			// if not, then create the organisation props on the top level object
			scopeObject[scope.OrganisationId] = types.OrganisationMapType{
				scope.TeamId: teamMap,
				"id":         scope.OrganisationId,
				"name":       scope.OrganisationName,
			}
		}
	}

	userUtil.scopeObject = scopeObject
	return userUtil.scopeObject
}
