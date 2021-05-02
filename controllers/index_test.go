package controllers

import (
	"app-auth/db"
	"app-auth/types"
	"app-auth/utils"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var (
	_ = godotenv.Load()
	sampleUsers map[string]*json.RawMessage
	fileData = utils.ReadFile("../data/data.json")
	data = types.UserData{}
	users = []string{"user_one", "user_two", "user_three", "user_four", "user_five"}
)

func parseDataFromJSONAndSaveToDatabase(t *testing.T, user string) {

	// this is a sample user data
	err := json.Unmarshal(*sampleUsers[user], &data)
	assert.Nil(t, err)

	//assert.True(t, data == config.UserData{})
	data.Id = uuid.NewV4().String()
	data.Password = utils.HashPassword(data.Password)

	filter := bson.M{"email": data.Email}
	resultContainer := types.UserData{}

	err = db.FindOne(collection, utils.CreateTimeoutContext(5), filter).Decode(&resultContainer); if err != nil {
		// there is no such document in the database.
		response := db.InsertIntoDB(collection, utils.CreateTimeoutContext(5), utils.UserBsonM(data))

		// assert that response is of kind mongo.InsertOneResult.
		assert.False(t, reflect.TypeOf(response) == reflect.TypeOf(mongo.InsertOneResult{}))
	} else {
		// this user already exists, there's no need to add him(gender-agnostic) again.
		assert.True(t, reflect.TypeOf(resultContainer) == reflect.TypeOf(types.UserData{}))
	}
}


func TestSampleUserSave(t *testing.T) {
	err := json.Unmarshal([]byte(fileData), &sampleUsers)
	// always assert that the error is nil and the json is correctly parsed
	assert.Nil(t, err)

	for _, value := range users {
		parseDataFromJSONAndSaveToDatabase(t, value)
	}
}


func TestLoginPost(t *testing.T) {
	// user login details sample.
	loginDetails := fmt.Sprintf(`{"email": "%s", "password": "@testUser"}`, data.Email)

	t.Log(os.Getwd())

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(loginDetails))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	body := types.LoginResponse{}

	// make sure calling the LoginPost returns no errors.
	if assert.NoError(t, LoginPost(c)) {
		// check is the status of the post request is OK.
		assert.Equal(t, http.StatusOK, rec.Code)

		err := json.Unmarshal(rec.Body.Bytes(), &body); if err != nil {
			t.Log(err)
		}

		// assert that the body message and the user login are successful.
		assert.Equal(t, body.Message, "User Login Successful")
	}
}


func TestGolangPaths(t *testing.T) {
	t.Log(os.Getwd())
}
