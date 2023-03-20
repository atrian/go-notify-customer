package channelServices

import (
	"context"

	external "github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/mail"
)

type Mail struct {
	conf configMail
	ctx  context.Context
}

type configMail interface {
	GetMailSenderAddress() string
	GetMailSMTPHost() string
	GetMailLogin() string
	GetMailPassword() string
	GetMailMessageTheme() string
}

func NewMail(ctx context.Context, conf configMail) *Mail {
	return &Mail{
		conf: conf,
		ctx:  ctx,
	}
}

func (s *Mail) SendMessage(message string, destination string) error {
	mailService := mail.New(s.conf.GetMailSenderAddress(), s.conf.GetMailSMTPHost())
	mailService.AddReceivers(destination)
	mailService.AuthenticateSMTP("", s.conf.GetMailLogin(), s.conf.GetMailPassword(), s.conf.GetMailSMTPHost())

	external.UseServices(mailService)

	return external.Send(
		context.TODO(),
		s.conf.GetMailMessageTheme(),
		message,
	)
}
