package runs

import (
	"app-auth/db"

	"log"

	uuid "github.com/satori/go.uuid"
)

func main() {
	log.Println("this is it")

	// send out an invitation email to the user
	mailId := uuid.NewV4().String()

	type inviteObjectType = map[string]interface{}
	var inviteObject = inviteObjectType{
		"app": "app",
		"user_id": "newUserId",
		"new_user": "newUser",
		"team_id": "team.Id",
		"user_email": "user.Email",
		"organisation_id": "team.OrganisationId",
		"signup_url": "signupAuthUrl",
	}

	// raises an error during redis cache
	// remove the saved user and return nil
	err := db.SetObject(mailId, inviteObject); if err != nil {
		log.Println(err)
	}

	res, err := db.GetObject(mailId); if err != nil {
		log.Println(err)
	} else {
		app := res[0]
		userId := res[1]
		newUser := res[2]
		teamId := res[3]
		userEmail := res[4]
		organisationId := res[5]
		signupUrl := res[6]

		log.Println(app)
		log.Println(userId)
		log.Println(newUser)
		log.Println(teamId)
		log.Println(userEmail)
		log.Println(organisationId)
		log.Println(signupUrl)
	}
}
