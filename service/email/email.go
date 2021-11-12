package email

import (
	"net/smtp"
	"os"
	log "github.com/sirupsen/logrus"
)

func SendEmail(theme, message string) {
	from := os.Getenv("EMAIL_ADDR")
	password := os.Getenv("EMAIL_PASSWORD")

	toEmail := "artyomsh01@yandex.ru"
	to := []string{toEmail}

	host := "smtp.gmail.com"
	port := "587"
	address := host + ":" + port


	emailMessage := []byte(theme + message)

	auth := smtp.PlainAuth("",from,password,host)

	err := smtp.SendMail(address, auth, from, to, emailMessage)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("Check email")
}