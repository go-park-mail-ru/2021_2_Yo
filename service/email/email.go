package email

import (
	log "github.com/sirupsen/logrus"
	"net/smtp"
	"os"
)

func SendEmail(theme, message string, recievers []string) {
	from := os.Getenv("EMAIL_ADDR")
	password := os.Getenv("EMAIL_PASSWORD")

	host := "smtp.gmail.com"
	port := "587"
	address := host + ":" + port

	theme += "\n"
	emailMessage := []byte(theme + message)

	auth := smtp.PlainAuth("", from, password, host)

	err := smtp.SendMail(address, auth, from, recievers, emailMessage)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("Check email")
}
