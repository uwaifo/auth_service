package utils

import (
	"fmt"
	"net/http"

	"app-auth/db"
	"app-auth/iam"
	"app-auth/types"

	"context"
	"log"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	uuid "github.com/satori/go.uuid"
)

type OrganisationUtil struct {
	Name        string `json:"organisation_name"`
	Description string `json:"description"`
	Id          string `json:"organisation_id"`

	App    string `json:"app"`
	UserId string `json:"user_id"`

	organisation types.Organisation
}

// String() result of logging the Organisation util
func (organisationUtil OrganisationUtil) String() string {
	return "<OrganisationUtil organisation_id='" + organisationUtil.Id + "' organisation_name='" + organisationUtil.Name + "'/>"
}

// GetOrganisation() return a db organisation
func (organisationUtil OrganisationUtil) GetOrganisation() *types.Organisation {
	// if there's no user init id
	if organisationUtil.Id == "" {
		return nil
	}

	organisationData := types.Organisation{}
	filter := organisationUtil.getOrganisationFilter()
	err := db.OrganisationCollection.FindOne(context.Background(), filter).Decode(&organisationData)
	if err != nil {
		log.Println(err)
		return nil
	}

	organisationUtil.organisation = organisationData
	organisationUtil.Name = organisationData.Name
	return &organisationUtil.organisation
}

// DestroyOrganisationRecord() return a an empty object
func (organisationUtil OrganisationUtil) DestroyOrganisationRecord() (bool, error) {
	foundRecord := false

	if organisationUtil.Id == "" {
		return foundRecord, nil
	}

	//teams := organisationUtil.GetOrganisationTeams()

	organisationData := types.Organisation{}
	filter := organisationUtil.getOrganisationFilter()
	err := db.OrganisationCollection.FindOne(context.Background(), filter).Decode(&organisationData)
	if err != nil {
		log.Println(err)
		return foundRecord, err
	}
	/////
	//deletableTeams , err := db.TeamCollection.DeleteMany(context.Background())

	ok, err := db.OrganisationCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Println(err)
		return foundRecord, nil
	}

	fmt.Println(ok)
	foundRecord = true
	return foundRecord, nil

}

// returns the filter for the organisation
func (organisationUtil OrganisationUtil) getOrganisationFilter() bson.D {
	return bson.D{{"id", organisationUtil.Id}, {"app", organisationUtil.App}}
}

// UpdateOrganisation() update organisation and return the new update
func (organisationUtil OrganisationUtil) UpdateOrganisation(organisationUpdateInfo bson.D) *types.Organisation {
	// if there's no such team, return nil
	var organisation = organisationUtil.GetOrganisation()
	if organisation == nil {
		return nil
	}

	organisationData := types.Organisation{}
	organisationUpdateData := bson.D{{"$set", organisationUpdateInfo}}
	filter := organisationUtil.getOrganisationFilter()

	tr := true
	rd := options.After
	opts := &options.FindOneAndUpdateOptions{Upsert: &tr, BypassDocumentValidation: &tr, ReturnDocument: &rd}

	err := db.OrganisationCollection.FindOneAndUpdate(context.Background(), filter, organisationUpdateData, opts).Decode(&organisationData)
	if err != nil {
		log.Println(err)
		return nil
	}

	filter = bson.D{{"organisationid", organisationData.Id}, {"app", organisationUtil.App}}
	teamAndScopeUpdateData := bson.D{{"$set", bson.D{{"organisationname", organisationData.Name}}}}

	// update the scopes if you update the organisation
	_, err = db.UserScopesCollection.UpdateMany(context.Background(), filter, teamAndScopeUpdateData)
	if err != nil {
		log.Println(err)
	}

	// update all teams in necessary when you update the organisation...
	_, err = db.TeamCollection.UpdateMany(context.Background(), filter, teamAndScopeUpdateData)
	if err != nil {
		log.Println(err)
	}

	organisationUtil.organisation = organisationData
	organisationUtil.Name = organisationData.Name
	organisationUtil.Id = organisationData.Id
	return &organisationUtil.organisation
}

// RemoveOrganisation() deletes an organisation
func (organisationUtil OrganisationUtil) RemoveOrganisation(deleteData bson.D) *types.Organisation {
	var organisation = organisationUtil.GetOrganisation()
	if organisation == nil {
		return nil
	}

	organisationData := types.Organisation{}
	//organisationUpdateData := bson.D{{"$set", bson.D{{"deleted", true}, {"deleted_at", time.Now()}}}}
	//organisationUpdateData := bson.D{{"$set", deleteData}}
	organisationUpdateData := bson.M{"$set": bson.M{
		"deleted":    true,
		"deletedat":  time.Now(),
		"deleted_at": time.Now(),
		"deletedby":  organisationUtil.UserId,
		"deleted_by": organisationUtil.UserId,
	}}

	filter := organisationUtil.getOrganisationFilter()

	tr := true
	rd := options.After
	opts := &options.FindOneAndUpdateOptions{Upsert: &tr, BypassDocumentValidation: &tr, ReturnDocument: &rd}

	err := db.OrganisationCollection.FindOneAndUpdate(context.Background(), filter, organisationUpdateData, opts).Decode(&organisationData)
	if err != nil {
		log.Println(err)
		return nil
	}
	//Schedule Permanent Deletion
	scheduleId := uuid.NewV4().String()
	scheduleObj := ScheduleUtil{Id: scheduleId}

	scheduleItem, err := scheduleObj.ScheduleTask("ORG", "DestroyOrganisationRecord", organisationData.Id)
	if err != nil {
		log.Println(err)

		return nil
	}
	fmt.Println(scheduleItem)

	/*

		 	scheduleId := uuid.NewV4().String()
			scheduleData := types.Schedule{
				Id:        scheduleId,
				ItemId:    organisationData.Id,
				ItemType:  "ORG",
				Action:    "DestroyOrganisationRecord",
				ExecuteOn: time.Now().Add(2 * time.Minute),
				CreatedAt: time.Now(),
			}

			scheduleResponse, err := db.AuthScheduleCollection.InsertOne(context.Background(), scheduleData)
			if err != nil {
				log.Println(err)
				return nil
			}
			fmt.Println(scheduleResponse)
	*/

	filter = bson.D{{"organisationid", organisationData.Id}, {"app", organisationUtil.App}}
	teamAndScopeUpdateData := bson.D{{"$set", bson.D{
		{"deleted", organisationData.Deleted},
		{"deleted_at", time.Now()},
		{"deletedat", time.Now()},
		{"deletedby", organisationUtil.UserId},
		{"deleted_by", organisationUtil.UserId},
	}}}

	//fmt.Println(teamAndScopeUpdateData)

	// update the scopes if you update the organisation
	_, err = db.UserScopesCollection.UpdateMany(context.Background(), filter, teamAndScopeUpdateData)
	if err != nil {
		log.Println(err)
	}

	// update all teams in necessary when you update the organisation...
	_, err = db.TeamCollection.UpdateMany(context.Background(), filter, teamAndScopeUpdateData)
	if err != nil {
		log.Println(err)
	}

	organisationUtil.organisation = organisationData
	organisationUtil.Name = organisationData.Name
	organisationUtil.Id = organisationData.Id

	//Call to cascade delete nested teams
	/*for _, v := range organisationUtil.organisation.Teams {
		fmt.Println(v)

	}*/

	return &organisationUtil.organisation

}

// RemoveOrganisation() deletes an organisation
func (organisationUtil OrganisationUtil) RestoreDeletedOrganisation(deleteData bson.D) *types.Organisation {
	var organisation = organisationUtil.GetOrganisation()
	if organisation == nil {
		return nil
	}

	organisationData := types.Organisation{}
	//organisationUpdateData := bson.D{{"$set", bson.D{{"deleted", true}, {"deleted_at", time.Now()}}}}
	//organisationUpdateData := bson.D{{"$set", deleteData}}
	organisationUpdateData := bson.M{"$set": bson.M{
		"deleted":    false,
		"deletedat":  time.Now(),
		"deleted_at": time.Now(),
		"deletedby":  "",
		"deleted_by": "",
	}}

	filter := organisationUtil.getOrganisationFilter()

	tr := true
	rd := options.After
	opts := &options.FindOneAndUpdateOptions{Upsert: &tr, BypassDocumentValidation: &tr, ReturnDocument: &rd}

	err := db.OrganisationCollection.FindOneAndUpdate(context.Background(), filter, organisationUpdateData, opts).Decode(&organisationData)
	if err != nil {
		log.Println(err)
		return nil
	}

	filter = bson.D{{"organisationid", organisationData.Id}, {"app", organisationUtil.App}}
	teamAndScopeUpdateData := bson.D{{"$set", bson.D{
		{"deleted", organisationData.Deleted},
		{"deletedat", time.Now()},
		{"deletedby", organisationUtil.UserId},
	}}}

	//fmt.Println(teamAndScopeUpdateData)

	// update the scopes if you update the organisation
	_, err = db.UserScopesCollection.UpdateMany(context.Background(), filter, teamAndScopeUpdateData)
	if err != nil {
		log.Println(err)
	}

	// update all teams in necessary when you update the organisation...
	_, err = db.TeamCollection.UpdateMany(context.Background(), filter, teamAndScopeUpdateData)
	if err != nil {
		log.Println(err)
	}

	organisationUtil.organisation = organisationData
	organisationUtil.Name = organisationData.Name
	organisationUtil.Id = organisationData.Id

	return &organisationUtil.organisation

}

// CreateOrganisation() create a new organisation
func (organisationUtil OrganisationUtil) CreateOrganisation(imageUrl string) *types.Organisation {

	// if there's such an organisation, return nil
	var organisation = organisationUtil.GetOrganisation()
	if organisation != nil {
		return nil
	}

	organisationData := types.Organisation{
		Id:          organisationUtil.Id,
		App:         organisationUtil.App,
		ImageUrl:    imageUrl,
		Name:        organisationUtil.Name,
		Description: organisationUtil.Description,
		CreatedBy:   organisationUtil.UserId,
		CreatedAt:   time.Now(),
		UpdateAt:    time.Now(),
		Teams:       []string{},
		Deleted:     false,
	}

	// if there's a problem with login return error
	_, err := db.OrganisationCollection.InsertOne(context.Background(), organisationData)
	if err != nil {
		log.Println(err)
		return nil
	}

	organisationUtil.organisation = *organisationUtil.GetOrganisation()
	organisationUtil.Name = organisationData.Name
	organisationUtil.Id = organisationData.Id
	return &organisationUtil.organisation
}

// AddOrganisationTeam() Add a newly created team to the organisation
func (organisationUtil OrganisationUtil) AddOrganisationTeam(teamData string) *bool {
	// if there's no such team, return nil
	var organisation = organisationUtil.GetOrganisation()
	if organisation == nil {
		return &AndFalse
	}

	// add the new team id to the organisation teams
	updateData := bson.D{{"teams", append(organisation.Teams, teamData)}}

	data := organisationUtil.UpdateOrganisation(updateData)
	if data == nil {
		return &AndFalse
	}

	return &AndTrue
}

// RemoveOrganisationTeam() remove a deleted team from the organisation
func (organisationUtil OrganisationUtil) RemoveOrganisationTeam(teamData string) *bool {
	// if there's no such team, return nil
	var organisation = organisationUtil.GetOrganisation()
	if organisation == nil {
		return &AndFalse
	}
	var teams = organisation.Teams

	// remove the team from the list of teams
	index := FindIndex(teams, teamData)
	newTeams := RemoveIndex(teams, index)

	// remove the deleted team id to the organisation teams
	updateData := bson.D{{"teams", newTeams}}

	data := organisationUtil.UpdateOrganisation(updateData)
	if data == nil {
		return &AndFalse
	}

	return &AndTrue
}

// GetOrganisationTeams() get all the teams within an organisation
func (organisationUtil OrganisationUtil) GetOrganisationTeams() *[]types.Team {
	// if there's no such organisation, return nil
	var organisation = organisationUtil.GetOrganisation()
	if organisation == nil {
		return nil
	}

	teamsFilter := bson.D{
		{"id", bson.D{{"$in", organisation.Teams}}},
		{"deleted", false},
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

		// append the teams data to the various array variables
		teams = append(teams, team)
	}

	// return all the teams
	return &teams
}

// GetOrganisationTeams() get all the teams within an organisation
func (organisationUtil OrganisationUtil) GetOrganisationDeletedTeams() *[]types.Team {
	// if there's no such organisation, return nil
	var organisation = organisationUtil.GetOrganisation()
	if organisation == nil {
		return nil
	}

	teamsFilter := bson.D{
		{"id", bson.D{{"$in", organisation.Teams}}},
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

func (organisationUtil OrganisationUtil) GetOrganisationDeletedTeamsCount() *types.TeamsCount {
	// Always add the scope so we don't count teams from different oprganisations or scopes and apps
	teamsFilter := bson.D{
		{"deleted", true},
		/*{"app", organisationUtil.App},
		{"$or", bson.A{
			// return any team you created
			bson.D{{"createdby", organisationUtil.UserId}},
			// or are a member of
			bson.D{
				{"members", bson.D{
					{"$in", []string{organisationUtil.UserId}},
				}},
			},
		}},*/
	}

	cur, err := db.TeamCollection.CountDocuments(context.Background(), teamsFilter)
	if err != nil {
		log.Println(err)
		return nil
	}

	var teamsCount = types.TeamsCount{
		Type:   "",
		Count:  int(cur),
		Status: http.StatusOK,
		Error:  false,
	}

	return &teamsCount
}

// HasOrganisationPermission() checks if a user has permission to perform a task on the organisation
func (organisationUtil OrganisationUtil) HasOrganisationPermission(userOrganisationScope iam.FindUserScopes, permission string) int {
	// if there's no such organisation, return 404
	var organisation = organisationUtil.GetOrganisation()
	if organisation == nil {
		return 404
	}

	// if the organisation found was created by the current user, then bypass the permission checks
	if organisation.CreatedBy == userOrganisationScope.Id {
		return 200
	}

	// get scopes pertaining to the organisation level
	scopes := userOrganisationScope.FindAndReturnUserScopesBasedOnFilter(userOrganisationScope.GetOrganisationsScopesLevelFilter())

	// we initialize a scope checker, which would check the permissions
	scopesChecker := iam.Scope{Scopes: scopes, App: userOrganisationScope.App}

	// check if the user has the required permission to perform this task
	if !scopesChecker.HasPermission(permission) {
		return 401
	}

	return 200
}
