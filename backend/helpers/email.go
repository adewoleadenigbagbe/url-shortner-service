package helpers

import (
	"os"
	"strconv"

	"github.com/go-mail/mail"
)

func SendEmail(tos []string) error {
	from := os.Getenv("EMAIL_FROM")
	pass := os.Getenv("EMAIL_PASSWORD")

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))

	m := mail.NewMessage()

	m.SetHeader("From", from)

	m.SetHeader("To", tos...)

	m.SetHeader("Subject", "Complete Onboarding Process")

	//TODO: change the link later after FE impl
	message := `
		<h2>
			Hello there. Welcome to Url Shortner Organization. Please click below link to get started
		</h2>
		<p>
		 	<a href="https://www.w3schools.com">Register user</a>
		</p>
	`

	m.SetBody("text/html", message)

	d := mail.NewDialer(smtpHost, smtpPort, from, pass)

	return d.DialAndSend(m)
}
