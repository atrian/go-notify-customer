package channelServices

type Twilio struct {
}

func NewTwilio() *Twilio {
	return &Twilio{}
}

func (s *Twilio) SendMessage(message string, destination string) error {
	return nil
}
