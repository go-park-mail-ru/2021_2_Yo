package email

import (
	"backend/internal/models"
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	log "backend/pkg/logger"
)

type Mail struct {
	Sender  string
	Subject string
	Body    bytes.Buffer
}

func BuildMessage(mail Mail) string {
	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg += fmt.Sprintf("From: %s\r\n", mail.Sender)
	msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	msg += fmt.Sprintf("\r\n%s\r\n", mail.Body.String())
	return msg
}

func SendEmail(theme, htmlTemplate string, info []*models.Info) {
	from := os.Getenv("EMAIL_ADDR")
	password := os.Getenv("EMAIL_PASSWORD")

	host := "smtp.gmail.com"
	port := "587"
	address := host + ":" + port

	ts, err := template.ParseFiles(htmlTemplate)
	if err != nil {
		fmt.Println(err)
	}

	for _, reciever := range info {

		var body bytes.Buffer
		err = ts.Execute(&body, reciever)
		if err != nil {
			log.Error(err)
		}

		request := Mail{
			Sender:  from,
			Subject: theme,
			Body:    body,
		}
		msg := BuildMessage(request)
		auth := smtp.PlainAuth("", from, password, host)

		err := smtp.SendMail(address, auth, from, []string{reciever.Mail}, []byte(msg))
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
