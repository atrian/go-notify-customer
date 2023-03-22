package channelServices

import (
	"context"
	external "github.com/nikoksr/notify"

	"github.com/atrian/go-notify-customer/internal/interfaces"
	"github.com/nikoksr/notify/service/twilio"
)

type Twilio struct {
	cfg    configTwilio
	logger interfaces.Logger
}

type configTwilio interface {
	GetTwilioAccountSid() string
	GetTwilioAuthToken() string
	GetTwilioSenderPhone() string
}

func NewTwilio(cfg configTwilio, logger interfaces.Logger) *Twilio {
	return &Twilio{
		cfg:    cfg,
		logger: logger,
	}
}

func (s *Twilio) SendMessage(ctx context.Context, message string, destination string) error {
	twilioSvc, err := twilio.New(
		s.cfg.GetTwilioAccountSid(),
		s.cfg.GetTwilioAuthToken(),
		s.cfg.GetTwilioSenderPhone(),
	)

	if err != nil {
		return err
	}

	twilioSvc.AddReceivers(destination)

	notifier := external.New()
	notifier.UseServices(twilioSvc)

	err = notifier.Send(context.Background(), "", message)
	if err != nil {
		return err
	}

	return nil
}
