package mail

import (
	"app-auth/config"
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"strings"
)


type Mail struct {
	senderId string
	toIds    []string
	subject  string
	body     string
}

type SmtpServer struct {
	host string
	port string
}

func (s *SmtpServer) ServerName() string {
	return s.host + ":" + s.port
}

func (mail *Mail) BuildMessage() string {
	message := ""
	message += fmt.Sprintf("From: %s\r\n", mail.senderId)

	if len(mail.toIds) > 0 { message += fmt.Sprintf("To: %s\r\n", strings.Join(mail.toIds, ";")) }

	message += fmt.Sprintf("Subject: %s\r\n", mail.subject)
	message += "\r\n" + mail.body

	return message
}

func SendMail(to string, subject string, body string) error {
	mail := Mail{}
	mail.senderId = config.DefaultMail
	mail.toIds = []string{to}
	mail.subject = subject
	mail.body = body

	messageBody := mail.BuildMessage()

	smtpServer := SmtpServer{host: "smtp.gmail.com", port: "465"}

	// log the host name.
	log.Println(smtpServer.host)

	//build an auth
	auth := smtp.PlainAuth("", mail.senderId, config.DefaultMailPassword, smtpServer.host)

	// Gmail will reject connection if it's not secure
	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName: smtpServer.host,
	}

	// create a tls security connection
	conn, err := tls.Dial("tcp", smtpServer.ServerName(), tlsconfig)
	if err != nil { log.Panic(err) }

	// create a new smtp client with the host using the secure tls connection protocol.
	client, err := smtp.NewClient(conn, smtpServer.host)
	if err != nil { log.Panic(err) }

	// step 1: Use Auth
	if err = client.Auth(auth); err != nil { log.Panic(err) }

	// step 2: add all from and to
	if err = client.Mail(mail.senderId); err != nil { log.Panic(err) }

	// iterate through the recipients and add them to client receipts.
	for _, k := range mail.toIds {
		if err = client.Rcpt(k); err != nil { log.Panic(err) }
	}

	// Get Data object from the client
	w, err := client.Data()
	if err != nil { log.Panic(err) }

	// write the data on the mail body to the smtp Data object
	_, err = w.Write([]byte(messageBody))
	if err != nil { log.Panic(err) }

	// close the data object
	err = w.Close()
	if err != nil { log.Panic(err) }

	err = client.Quit(); if err != nil { log.Panic(err) }

	log.Println("Mail sent successfully")

	return err

}
