package data

import (
	"app-auth/types"
	"app-auth/utils"

	"github.com/satori/go.uuid"
)

var SampleUser = types.UserData{
	Id: uuid.NewV4().String(),
	Password: utils.HashPassword("@testDeveloper"),
	Email: "developer.test@clipsynphony.com",
	Username: "test_developer",
	Firstname: "Test",
	Lastname: "Developer",
	Picture: "",
	Age: 22,
	Verified: true,
	Provider: "local",
}
