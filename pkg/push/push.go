package push

import (
	dto "notifier/pkg/dto"
)

type IConnection interface {
	SendMessage(dto.DTOMessagePush) (err error)
}