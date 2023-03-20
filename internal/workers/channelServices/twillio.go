package channelServices

import (
	"context"
	external "github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/twilio"
)

type Twilio struct {
	ctx context.Context
	cfg configTwilio
}

type configTwilio interface {
	GetTwilioAccountSid() string
	GetTwilioAuthToken() string
	GetTwilioSenderPhone() string
}

func NewTwilio(ctx context.Context, cfg configTwilio) *Twilio {
	return &Twilio{
		ctx: ctx,
		cfg: cfg,
	}
}

func (s *Twilio) SendMessage(message string, destination string) error {
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
