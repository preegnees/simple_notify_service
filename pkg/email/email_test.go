package email

import (
	"testing"
	models "notifier/pkg/models"
)

func TestSendEmai(t *testing.T) {
	emailService := New("", "")
	m := models.Message{
	Destination: "email",
	Email: "gobox.v1@gmail.com",
	Username: "",
	MessageSubject: "hello, v1",
	Message: "your code 40444",
	}
	emailService.Send(m)
}