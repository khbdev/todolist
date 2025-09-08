package rabbitmq

import (

	"os"

	"gopkg.in/gomail.v2"
)

func SendEmail(job EmailJob) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := 587
	smtpEmail := os.Getenv("SMTP_EMAIL")
	smtpPass := os.Getenv("SMTP_PASSWORD")





	var bodyStr string = "Salom Todolist App Xush Kelibsiz"


	m := gomail.NewMessage()
	m.SetHeader("From", smtpEmail)
	m.SetHeader("To", job.Email)
	m.SetHeader("Subject", "Welcome to TodoList")
	m.SetBody("text/html", bodyStr)

	d := gomail.NewDialer(smtpHost, smtpPort, smtpEmail, smtpPass)
	return d.DialAndSend(m)
}

