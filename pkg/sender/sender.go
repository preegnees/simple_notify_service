package sender

import (
	"context"

	models "notifier/pkg/models"
)

const EMAIL = "email"
const PUSH = "push"

type ISender interface {
	Send(models.Message)
}

type IEmail interface {
	ISender
}

type IPush interface {
	ISender
}

type sender struct {
	ctx          context.Context
	emailService IEmail
	pushService  IPush
	source       <-chan models.Message
}

func New(ctx context.Context, emailService IEmail, pushService IPush, source <-chan models.Message) *sender {

	return &sender{
		ctx:          ctx,
		emailService: emailService,
		pushService:  pushService,
		source:       source,
	}
}

func (s *sender) Run() error {

	select {
	case <-s.ctx.Done():
		return nil
	case m, ok := <-s.source:
		if !ok {
			return nil
		}
		if m.Destination == EMAIL {
			s.emailService.Send(m)
		} else if m.Destination == PUSH {
			s.pushService.Send(m)
		}
	}
	return nil
}
