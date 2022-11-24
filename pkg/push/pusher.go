package push

import (
	models "notifier/pkg/models"
)

type IConnection interface {
	SendMessage(models.Message) (err error)
}