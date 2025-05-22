package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

func NewSendGrid(apiKey, fromEmail string) *SendGridMailer {
	client := sendgrid.NewSendClient(apiKey)

	return &SendGridMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    client,
	}
}

func (m *SendGridMailer) Send(templateFile, username, email string, data any, isSandbox bool) (int, error) {
	from := mail.NewEmail(FromName, m.fromEmail)
	to := mail.NewEmail(username, email)

	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return -1, err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return -1, err
	}
	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return -1, err
	}

	message := mail.NewSingleEmail(from, subject.String(), to, "", body.String())

	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSandbox,
		},
	})
	var retryErr error
	for i := 0; i < MaxRetries; i++ {
		response, retryErr := m.client.Send(message)
		if retryErr != nil {
			log.Printf("Failed to send email to %v, attempt %d of %d", email, i+1, MaxRetries)
			log.Printf("Error: %v", retryErr.Error())

			// exponential backoff
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}

		log.Printf("Email sent with status code %v", email, response.StatusCode)
		return response.StatusCode, nil
	}

	return -1, fmt.Errorf("failed to send email to %v after %d attempts, error: %v", email, MaxRetries, retryErr)
}
