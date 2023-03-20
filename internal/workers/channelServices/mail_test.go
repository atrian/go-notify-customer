package channelServices

import (
	"context"
	"testing"
	"time"

	external "github.com/nikoksr/notify"
	"github.com/stretchr/testify/assert"
)

var _ configMail = (*mailConfigMock)(nil)

func TestMail_SendMessage(t *testing.T) {
	conf := mailConfigMock{}

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	chanService := NewMail(ctx, conf)

	err := chanService.SendMessage("Test", "destination@unknown.ru")
	assert.ErrorIs(t, err, external.ErrSendNotification)
}

type mailConfigMock struct{}

func (m mailConfigMock) GetMailSenderAddress() string {
	return "test@sender.ru"
}

func (m mailConfigMock) GetMailSMTPHost() string {
	return "localhost:465"
}

func (m mailConfigMock) GetMailLogin() string {
	return "mail@login.ru"
}

func (m mailConfigMock) GetMailPassword() string {
	return "mail@login.ru"
}

func (m mailConfigMock) GetMailMessageTheme() string {
	return "mail theme"
}
