package channelServices

type Mail struct {
}

func NewMail() *Mail {
	return &Mail{}
}

func (s *Mail) SendMessage(message string, destination string) error {
	return nil
}
