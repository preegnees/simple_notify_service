package kafka

import (
	models "notifier/pkg/models"
	"testing"

	"github.com/segmentio/kafka-go/protocol"
)

func TestConsume(t *testing.T) {
	k := kafka_{}

	messageBefore := models.Message{
		Destination: "Email",
		Email: "Email@Email",
		Username: "radmir",
		MessageSubject: "ist test",
		Message: "hello world",
	}

	hs := []protocol.Header{
		{Key: "Destination", Value: []byte(messageBefore.Destination)},
		{Key: "Email", Value: []byte(messageBefore.Email)},
		{Key: "Username", Value: []byte(messageBefore.Username)},
		{Key: "MessageSubject", Value: []byte(messageBefore.MessageSubject)},
		{Key: "Message", Value: []byte(messageBefore.Message)},

	}
	messageAfter, err := k.convertToMessage(hs)
	if err != nil {
		panic(err)
	}
	if messageAfter != messageBefore {
		panic("messageAfter != messageBefore")
	}
}