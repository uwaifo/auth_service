package mail


import (
	"fmt"

	"app-auth/config"
	
	"gopkg.in/gomail.v2"
)


func Mailer(to string, subject string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "no-reply@scaratec.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer("smtp.scaratec.com", 587, config.DefaultMail, config.DefaultMailPassword)

	err := d.DialAndSend(m); if err != nil { 
		fmt.Println(err)
	}

	return err
}