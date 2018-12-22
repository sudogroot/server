package email

import (
	"runtime"

	"gopkg.in/gomail.v2"
)

// SendEmailRequest struct
type SendEmailRequest struct {
	From         string
	To           string
	Subject      string
	Body         string
	SMTPHostname string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
}

// SendAsyncEmail ...
func SendAsyncEmail(sendEmailRequest SendEmailRequest) {
	// Set home many CPU cores this function wants to use
	runtime.GOMAXPROCS(runtime.NumCPU())

	m := gomail.NewMessage()
	m.SetHeader("From", sendEmailRequest.From)
	m.SetHeader("To", sendEmailRequest.To)
	m.SetHeader("Subject", sendEmailRequest.Subject)
	m.SetBody("text/html", sendEmailRequest.Body)

	d := gomail.NewDialer(sendEmailRequest.SMTPHostname, sendEmailRequest.SMTPPort, sendEmailRequest.SMTPUsername, sendEmailRequest.SMTPPassword)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

// SendSyncEmail ...
func SendSyncEmail(sendEmailRequest SendEmailRequest) bool {
	m := gomail.NewMessage()
	m.SetHeader("From", sendEmailRequest.From)
	m.SetHeader("To", sendEmailRequest.To)
	m.SetHeader("Subject", sendEmailRequest.Subject)
	m.SetBody("text/html", sendEmailRequest.Body)

	d := gomail.NewDialer(sendEmailRequest.SMTPHostname, sendEmailRequest.SMTPPort, sendEmailRequest.SMTPUsername, sendEmailRequest.SMTPPassword)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return false
	}

	return true
}
