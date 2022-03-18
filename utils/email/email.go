package email

import (
	"log"
	"net/smtp"
)

// Send an e-mail using the configuration to obtain login details, recipient, etc.
// @param subject - The subject for the e-mail
// @param body - The body of the e-mail
// @param emailUser - The user to use for login to the SMTP server
// @param emailPassword - The password to use for login to the SMTP server
// @param emailServer - The SMTP server to use for sending e-mail
// @param emailPort - The port for the SMTP server
// @param emailTo - The recipient
// @param emailFrom - The originator of the e-mail
func SendEmail(subject string, body string, emailUser string, emailPassword string, emailServer string,
	emailPort string, emailTo string, emailFrom string) {
	if len(emailServer) == 0 { // Do not proceed if an email server was not provided.
		return
	}
	// Create the authentication object
	auth := smtp.PlainAuth("", emailUser, emailPassword, emailServer)

	// Prepare the e-mail for transmission.
	to := []string{emailTo}
	msg := []byte("To: " + emailTo + "\r\n" +
		"From: " + emailFrom + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body + "\r\n")

	// Transmit.
	err := smtp.SendMail(emailServer+":"+emailPort, auth, emailFrom, to, msg)
	if err != nil { // Do not fatally error in case the e-mail server is just temporarily down.
		log.Printf("Problem sending e-mail notification: \n", err)
	}
	log.Println("Email sent: " + subject)
}
