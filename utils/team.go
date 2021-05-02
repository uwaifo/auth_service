package utils

import (
	"app-auth/db"
	"app-auth/iam"
	"app-auth/types"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"

	"context"
	"log"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

type TeamUtil struct {
	TeamName    string `json:"team_name"`
	Description string `json:"description"`
	TeamId      string `json:"team_id"`
	UserId      string `json:"user_id"`
	UserEmail   string `json:"email"`
	App         string `json:"app"`

	team         types.Team
	Organisation OrganisationUtil
}

// toString() representation of the
func (teamUtil TeamUtil) String() string {
	return "<TeamUtil team_id='" + teamUtil.TeamId + "' team_name='" + teamUtil.TeamName + "' organisation_id='" + teamUtil.Organisation.Id + "'/>"
}

func (teamUtil TeamUtil) getTeamFilter() bson.D {
	return bson.D{{"id", teamUtil.TeamId}, {"app", teamUtil.App}, {"organisationid", teamUtil.Organisation.Id}}
}

// GetTeam() returns the team data
func (teamUtil TeamUtil) GetTeam() *types.Team {
	// if there's no user init id
	if teamUtil.TeamId == "" {
		return nil
	}

	teamData := types.Team{}
	filter := teamUtil.getTeamFilter()
	fmt.Println(filter)
	err := db.TeamCollection.FindOne(context.Background(), filter).Decode(&teamData)
	if err != nil {
		log.Println(err)
		return nil
	}

	teamUtil.team = teamData
	return &teamUtil.team
}

func (teamUtil TeamUtil) DestroyTeamRecord() (bool, error) {
	// if there's no user init id
	foundRecord := false
	if teamUtil.TeamId == "" {
		return foundRecord, nil
	}

	//teams := organisationUtil.GetOrganisationTeams()

	teamData := types.Team{}
	filter := teamUtil.getTeamFilter()

	//filter := teamUtil.GetTeam()
	err := db.TeamCollection.FindOne(context.Background(), filter).Decode(&teamData)
	if err != nil {
		fmt.Println("didnt find record")

		log.Println(err)
		return foundRecord, err
	} else {
		fmt.Println("found record")
	}
	/////
	//deletableTeams , err := db.TeamCollection.DeleteMany(context.Background())

	ok, err := db.TeamCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Println(err)
		return foundRecord, nil
	}
	fmt.Println(ok)
	foundRecord = true
	return foundRecord, nil

}

// UpdateTeam() returns an update of the team data
func (teamUtil TeamUtil) UpdateTeam(teamUpdateInfo bson.D) *types.Team {
	// if there's no such team, return nil
	var team = teamUtil.GetTeam()
	if team == nil {
		return nil
	}

	teamData := types.Team{}
	teamUpdateData := bson.D{{"$set", teamUpdateInfo}}
	filter := teamUtil.getTeamFilter()

	tr := true
	rd := options.After
	opts := &options.FindOneAndUpdateOptions{Upsert: &tr, BypassDocumentValidation: &tr, ReturnDocument: &rd}

	err := db.TeamCollection.FindOneAndUpdate(context.Background(), filter, teamUpdateData, opts).Decode(&teamData)
	if err != nil {
		log.Println(err)
		return nil
	}

	filter = bson.D{{"teamid", teamData.Id}, {"app", teamUtil.App}}
	teamAndScopeUpdateData := bson.D{{"$set", bson.D{{"teamname", teamData.Name}}}}

	// update the scopes if you update the organisation
	_, err = db.UserScopesCollection.UpdateMany(context.Background(), filter, teamAndScopeUpdateData)
	if err != nil {
		log.Println(err)
	}

	teamUtil.team = teamData
	teamUtil.TeamName = teamData.Name
	teamUtil.TeamId = teamData.Id
	return &teamUtil.team
}

// CreateTeam() returns an update of the team data
func (teamUtil TeamUtil) CreateTeam(imageUrl string) *types.Team {
	// if there's such team, return nil
	var team = teamUtil.GetTeam()
	if team != nil {
		return nil
	}
	var organisationName = teamUtil.Organisation.Name

	teamData := types.Team{
		Name:             teamUtil.TeamName,
		Id:               teamUtil.TeamId,
		ImageUrl:         imageUrl,
		OrganisationId:   teamUtil.Organisation.Id,
		OrganisationName: organisationName,
		CreatedBy:        teamUtil.UserId,
		App:              teamUtil.App,
		CreatedAt:        time.Now(),
		UpdateAt:         time.Now(),
		Members:          []string{teamUtil.UserId},
		Deleted:          false,
	}

	// add team to the database
	_, err := db.TeamCollection.InsertOne(context.Background(), teamData)
	if err != nil {
		log.Println(err)
		return nil
	}

	// create a new scope for the user
	userMemberScope := types.UserMemberScope{
		Id:               teamUtil.UserId,
		OrganisationId:   teamUtil.Organisation.Id,
		OrganisationName: organisationName,
		TeamId:           teamUtil.TeamId,
		TeamName:         teamUtil.TeamName,
		UserId:           teamUtil.UserId,
		UserEmail:        teamUtil.UserEmail,
		App:              teamUtil.App,
		State:            "creator",
		Scopes:           []string{"Root"},
	}

	res := userMemberScope.SaveUserMemberScope()
	if res == nil {
		return nil
	}

	// let us update the organisation then
	data := teamUtil.Organisation.AddOrganisationTeam(teamUtil.TeamId)
	if data == nil || data == &AndFalse {
		return nil
	}

	team = teamUtil.GetTeam()
	if team == nil {
		return nil
	}
	// now that a team is saved, getting the team should return some data
	teamUtil.team = *team
	teamUtil.TeamName = teamUtil.team.Name
	teamUtil.TeamId = teamUtil.team.Id

	// then return the team
	return &teamUtil.team
}

// RemoveTeam() returns an update of the team data
func (teamUtil TeamUtil) RemoveTeam() *bool {
	// if there's no such team, return nil
	var team = teamUtil.GetTeam()
	if team == nil {
		return &AndFalse
	}

	filter := bson.D{{
		"id", bson.D{{"$in", team.Members}},
	}, {
		"teamid", team.Id,
	}}

	_, err := db.UserScopesCollection.DeleteMany(context.Background(), filter)
	if err != nil {
		log.Println(err)
		return &AndFalse
	}

	_, err = db.TeamCollection.DeleteOne(context.Background(), teamUtil.getTeamFilter())
	if err != nil {
		log.Println(err)
		return &AndFalse // return false if the team is not deleted
	}
	// let us update the organisation then
	data := teamUtil.Organisation.RemoveOrganisationTeam(teamUtil.TeamId)
	if data == nil || data == &AndFalse {
		return nil
	}

	// then return true if the team is removed
	return &AndTrue
}

func (teamUtil TeamUtil) RestoreDeletedTeam() *types.Team {
	var team = teamUtil.GetTeam()
	if team == nil {
		return nil
	}

	teamData := types.Team{}

	// instruction to set deleted field to false
	restoreStatment := bson.M{"$set": bson.M{
		"deleted": false,
	}}
	filter := teamUtil.getTeamFilter()

	tr := true
	rd := options.After
	opts := &options.FindOneAndUpdateOptions{Upsert: &tr, BypassDocumentValidation: &tr, ReturnDocument: &rd}

	err := db.TeamCollection.FindOneAndUpdate(context.Background(), filter, restoreStatment, opts).Decode(&teamData)
	if err != nil {
		log.Println(err)
		return nil
	}

	//return &teamData
	teamUtil.team = teamData
	teamUtil.TeamName = teamData.Name
	teamUtil.TeamId = teamData.Id
	fmt.Println("Name : ", &teamUtil.team.Name)
	return &teamUtil.team

}
func (teamUtil TeamUtil) RemoveTeamTemp() *types.Team {

	var team = teamUtil.GetTeam()
	if team == nil {
		return nil
	}

	teamData := types.Team{}
	//teamUpdateData := bson.D{{"$set", bson.D{{"deleted", true}, {"deleted_by", "somebody"}, {"deleted_at", time.Now()}}}}
	teamUpdateData := bson.M{"$set": bson.M{
		"deleted":   true,
		"deletedat": time.Now(),
		//"deleted_at": time.Now(),
		"deletedby": teamUtil.UserId,
		//"deleted_by": teamUtil.UserId,
	}}
	filter := teamUtil.getTeamFilter()

	tr := true
	rd := options.After
	opts := &options.FindOneAndUpdateOptions{Upsert: &tr, BypassDocumentValidation: &tr, ReturnDocument: &rd}

	err := db.TeamCollection.FindOneAndUpdate(context.Background(), filter, teamUpdateData, opts).Decode(&teamData)
	if err != nil {
		log.Println(err)
		return nil
	}

	// Call function to schedule deleteton in 30 days
	//Schedule Permanent Deletion
	scheduleId := uuid.NewV4().String()
	scheduleObj := ScheduleUtil{Id: scheduleId}

	scheduleItem, err := scheduleObj.ScheduleTask("TEAM", "DestroyTeamRecord", team.Id)
	if err != nil {
		log.Println(err)

		return nil
	}
	fmt.Println(scheduleItem)

	//scheduleTask()

	filter = bson.D{{"teamid", teamData.Id}, {"app", teamUtil.App}}
	//teamAndScopeUpdateData := bson.D{{"$set", bson.D{{"deleted", true}, {"deleted_by", "somebody"}, {"deleted_at", time.Now()}}}}
	_, err = db.UserScopesCollection.UpdateMany(context.Background(), filter, teamUpdateData)
	if err != nil {
		log.Println(err)
	}
	teamUtil.team = teamData
	teamUtil.TeamName = teamData.Name
	teamUtil.TeamId = teamData.Id
	return &teamUtil.team

}

// AddMember() Add a new member to the team using the user email
func (teamUtil TeamUtil) AddMember(memberEmail string, signupAuthUrl string, appRedirecturl string, app string) *types.UserMemberScope {

	// the team we're trying to add the user to doesn't exist
	team := teamUtil.GetTeam()
	if team == nil {
		return nil
	}

	// we check if the user exists, if s/he does, cool and if not, we create a new user
	userUtil := UserMemberUtil{UserEmail: &memberEmail, UserId: nil, App: app}
	newUser := false
	newUserId := uuid.NewV4().String()

	user := userUtil.GetUserByEmail()
	if user == nil {
		userData := types.UserData{
			Id:                newUserId,
			Password:          "",
			Email:             memberEmail,
			Username:          "",
			Firstname:         "",
			Lastname:          "",
			Picture:           "",
			Age:               0,
			Active:            true,
			Verified:          false,
			Provider:          "local",
			LastLoggedInScope: "",
		}

		newUser = true
		// for some reason the user was not created and so we just return nil
		// and an error report to the request
		user = userUtil.CreateNewUser(userData)
		if user == nil {
			log.Println(user)
			return nil
		}
	}

	userUtil.UserId = &user.Id

	// send out an invitation email to the user
	mailId := uuid.NewV4().String()

	// now we create a user scope
	scope := types.UserMemberScope{
		Id:               user.Id,
		OrganisationName: team.OrganisationName,
		OrganisationId:   team.OrganisationId,
		TeamName:         team.Name,
		TeamId:           team.Id,
		UserEmail:        user.Email,
		UserId:           user.Id,
		App:              app,       // this state would correspond to a user invite
		State:            "Invited", // this state would correspond to a user invite
		Scopes:           []string{},
	}

	// check if this scopeMember exists
	scopeMember := scope.GetUserMemberScope()
	if scopeMember != nil {
		// this particular user already exists
		return nil
	}

	// when there's an error sending mail to the user
	// delete the user and return nil
	err := SendTeamEmail(user.Email, mailId)
	if err != nil {
		if newUser == true {
			_ = userUtil.RemoveUser()
		}
		return nil
	}

	type inviteObjectType = map[string]interface{}
	var inviteObject = inviteObjectType{
		"app":              app,
		"user_id":          user.Id,
		"new_user":         newUser,
		"team_id":          team.Id,
		"user_email":       user.Email,
		"signup_url":       signupAuthUrl,
		"app_redirect_url": appRedirecturl,
		"organisation_id":  team.OrganisationId,
	}

	// raises an error during redis cache
	// remove the saved user and return nil
	err = db.SetObject(mailId, inviteObject)
	if err != nil {
		if newUser == true {
			_ = userUtil.RemoveUser()
		}
		return nil
	}

	// if saving the scope return nil
	// delete the redis object and delete the saved user
	res := scope.SaveUserMemberScope()
	if res == nil {
		if newUser == true {
			_ = userUtil.RemoveUser()
		}
		_ = db.Del(mailId)
		return nil
	}

	// the team has not been update
	updateRes := teamUtil.UpdateTeam(bson.D{{"members", append(team.Members, scope.Id)}})
	if updateRes == nil {
		// we can check in the confirmation if the scope exists in the user
	}

	return &scope
}

// RemoveMember()
func (teamUtil TeamUtil) RemoveMember(memberId string) *bool {
	// if there's no such team, return nil
	var team = teamUtil.GetTeam()
	if team == nil {
		return nil
	}

	// member does not exist in the team
	index := FindIndex(team.Members, memberId)
	if index == -1 {
		return &AndFalse
	}

	// this filter is getting all the user scope objects within a team
	filter := bson.D{{"id", memberId}, {"teamid", team.Id}}

	// find a specific user, delete it and if an error occurs, return nil
	_, err := db.UserScopesCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Println(err)
		return &AndFalse
	}

	// remove the team from the list of teams
	newMembers := RemoveIndex(team.Members, index)

	// remove the user id from the list of members
	team = teamUtil.UpdateTeam(bson.D{{"members", newMembers}})
	if team == nil {
		// update of the team has failed...
		return nil
	}

	return &AndTrue
}

// GetMembers()
func (teamUtil TeamUtil) GetMembers() *[]types.UserMemberScope {
	// if there's no such team, return nil
	var team = teamUtil.GetTeam()
	if team == nil {
		return nil
	}

	// this filter is getting all the user scope objects within a team
	filter := bson.D{{
		"id", bson.D{{"$in", team.Members}},
	}, {
		"teamid", team.Id,
	}}

	// find all the user scope object and return nil is an error occurs during the search
	cur, err := db.UserScopesCollection.Find(context.Background(), filter)
	if err != nil {
		log.Println(err)
		return nil
	}

	defer cur.Close(context.TODO())
	var userMemberScopes []types.UserMemberScope

	// loop through the results till and append the items to the userMemberScope object
	for cur.Next(context.TODO()) {
		var userMemberScope types.UserMemberScope
		err := cur.Decode(&userMemberScope)
		if err != nil {
			log.Print(err)
		}

		// append the user member scope objects data to the various array variables
		userMemberScopes = append(userMemberScopes, userMemberScope)
	}

	return &userMemberScopes
}

// GetMember()
func (teamUtil TeamUtil) GetMember(memberId string) *types.UserMemberScope {
	// if there's no such team, return nil
	var team = teamUtil.GetTeam()
	if team == nil {
		return nil
	}

	// this filter is getting all the user scope objects within a team
	filter := bson.D{{"id", memberId}, {"teamid", team.Id}}

	// get a new variable of the user member scope type
	userMemberScope := types.UserMemberScope{}

	// find all the user scope object and return nil is an error occurs during the search
	err := db.UserScopesCollection.FindOne(context.Background(), filter).Decode(&userMemberScope)
	if err != nil {
		log.Println(err)
		return nil
	}

	// return this specific member scope
	return &userMemberScope
}

// HasOrganisationTeamPermission() checks if a user has a permission to perform a specific task on the team
func (teamUtil TeamUtil) HasOrganisationTeamPermission(userOrganisationScope iam.FindUserScopes, permission string) int {
	// get scopes pertaining to the team level
	scopes := userOrganisationScope.FindAndReturnUserScopesBasedOnFilter(userOrganisationScope.GetTeamsScopesLevelFilter())

	// we initialize a scope checker, which would check the permissions
	scopesChecker := iam.Scope{Scopes: scopes, App: userOrganisationScope.App}

	// check if the user has the required permission to perform this task
	if !scopesChecker.HasPermission(permission) {
		return 401
	}

	return 200
}

func (teamUtil TeamUtil) GetAllDeletedTeams() *[]types.Team {

	teamsFilter := bson.D{
		//{"id", bson.D{{"$in", organisation.Teams}}},
		{"deleted", true},
	}

	cur, err := db.TeamCollection.Find(context.Background(), teamsFilter)
	if err != nil {
		log.Println(err)
		return nil
	}

	defer cur.Close(context.TODO())
	var teams []types.Team

	for cur.Next(context.TODO()) {
		var team types.Team
		err := cur.Decode(&team)
		if err != nil {
			log.Print(err)
		}

		// append the teams data to the array
		teams = append(teams, team)
	}

	// return all the teams
	return &teams

}
