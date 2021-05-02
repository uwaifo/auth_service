package db

import (
	"os"

	"github.com/mongodb/mongo-go-driver/mongo"
)

func createCollection(collectionName string) *mongo.Collection {
	// get environment variable...
	databaseName := os.Getenv("MONGO_DB_NAME")
	if databaseName == "" {
		databaseName = "testing"
	}
	return Collection(Client, databaseName, collectionName)
}

var UserCollection = createCollection("users")
var OrganisationCollection = createCollection("organisation")
var TeamCollection = createCollection("teams")
var UserScopesCollection = createCollection("user_scopes")
var ScopesCollection = createCollection("scopes")
var AuthScheduleCollection = createCollection("auth_schedules")
