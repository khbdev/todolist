package rabbitmq

import (
	
	"html/template"
	"os"

	"gopkg.in/gomail.v2"
)

func SendEmail(job EmailJob) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := 587
	smtpEmail := os.Getenv("SMTP_EMAIL")
	smtpPass := os.Getenv("SMTP_PASSWORD")

	tmpl := `<h1>Welcome to TodoList, {{.Name}}</h1>
<p>We’re excited to have you on board!</p>`

	t, err := template.New("welcome").Parse(tmpl)
	if err != nil {
		return err
	}

	var bodyStr string
	buf := &bodyStrBuilder{&bodyStr}
	t.Execute(buf, job)

	m := gomail.NewMessage()
	m.SetHeader("From", smtpEmail)
	m.SetHeader("To", job.Email)
	m.SetHeader("Subject", "Welcome to TodoList")
	m.SetBody("text/html", bodyStr)

	d := gomail.NewDialer(smtpHost, smtpPort, smtpEmail, smtpPass)
	return d.DialAndSend(m)
}

// Helper template uchun
type bodyStrBuilder struct{ s *string }
func (b *bodyStrBuilder) Write(p []byte) (int, error) {
	*b.s += string(p)
	return len(p), nil
}
