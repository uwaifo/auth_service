package utils

import (
	"app-auth/db"
	"app-auth/mail"

	"fmt"
	"log"
	"os"

	uuid "github.com/satori/go.uuid"
)

func SendTeamEmail(rcpt string, emailId string) error {
	url := os.Getenv("ENDPOINT")
	if url == "" {
		url = "http://app-auth"
	}
	body := mail.TeamInviteTemplate(fmt.Sprintf(`%s/user-invite-confirm/%s`, url, emailId))
	subject := "Confirm Team Invite !"

	err := mail.Mailer(rcpt, subject, body)
	if err != nil {
		return err
	}
	return err
}

func SendPasswordResetEmail(rcpt string, redirectUrl string) error {
	body := mail.TeamInviteTemplate(redirectUrl)
	subject := "Password Change Request!"
	err := mail.Mailer(rcpt, subject, body)
	if err != nil {
		return err
	}
	return err
}

func SendAdminConfirmationMail(rcpt string, redirectUrl string, email string, username string) error {
	body := mail.AdminConfirmTemplate(redirectUrl, email, username)
	subject := "Confirm User Signup!"

	err := mail.Mailer(rcpt, subject, body)
	if err != nil {
		return err
	}
	return err

}

func SendResetEmail(rcpt string, redirectUrl string) error {
	body := mail.EmailResetTemplate(redirectUrl)
	subject := "Email Change Request!"

	err := mail.Mailer(rcpt, subject, body)
	if err != nil {
		return err
	}
	return err
}

func SendSignupEmail(rcpt string, userId string, app string, appRedirect string) error {

	confirmId := uuid.NewV4().String()
	redirectId := uuid.NewV4().String()

	body := mail.SignupTemplate(confirmId, redirectId, app)
	subject := "Email Verification!"

	err := mail.Mailer(rcpt, subject, body)
	if err != nil {
		log.Panic(err)
	}

	err1 := db.Set(confirmId, userId)
	err2 := db.Set(redirectId, appRedirect)

	log.Print(err1)
	log.Print(err2)

	return err
}
