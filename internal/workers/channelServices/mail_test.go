package channelServices

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var _ configMail = (*mailConfigMock)(nil)

func TestMail_SendMessage(t *testing.T) {
	// TODO восстановить тест воркера
	assert.NoError(t, nil)
}

type mailConfigMock struct{}

func (m mailConfigMock) IsMailTLSRequired() bool {
	return true
}

func (m mailConfigMock) GetMailSenderAddress() string {
	return "test@sender.ru"
}

func (m mailConfigMock) GetMailSMTPHost() string {
	return "sandbox.smtp.mailtrap.io:465"
}

func (m mailConfigMock) GetMailLogin() string {
	return "88aab5f5167a58"
}

func (m mailConfigMock) GetMailPassword() string {
	return "4b1cf7212d8f34"
}

func (m mailConfigMock) GetMailMessageTheme() string {
	return "Mail theme"
}
