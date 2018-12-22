/*
	Copyright (C) 2018 Nirmal Almara
	
    This file is part of Joyread.

    Joyread is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    Joyread is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
	along with Joyread.  If not, see <https://www.gnu.org/licenses/>.
*/

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
