package email

import (
	"fmt"
	"net/smtp"
	models "notifier/pkg/models"
	sender "notifier/pkg/sender"
)

type email struct {
	apiKey    string
	mainEmail string

	// kafkaErrCli
}

func New(apiKey, mainEmail string) *email {

	if apiKey == "" || mainEmail == "" {
		apiKey = "dvjbyrbgeyuvrjgt"
		mainEmail = "gobox.v1@gmail.com"
	}
	return &email{
		apiKey:    apiKey,
		mainEmail: mainEmail,
	}
}

var _ sender.IEmail = (*email)(nil)

func (e email) Send(m models.Message) {
	to := []string{m.Email}

	host := "smtp.gmail.com"
	port := "587"
	address := host + ":" + port

	subject := fmt.Sprintf("Subject: %s\n", m.MessageSubject)
	body := fmt.Sprintf("%s", m.Message)
	message := []byte(subject + body)

	auth := smtp.PlainAuth("", e.mainEmail, e.apiKey, host)

	err := smtp.SendMail(address, auth, e.mainEmail, to, message)
	if err != nil {
		panic(err) // тут в кафку отправлять
	}
}
