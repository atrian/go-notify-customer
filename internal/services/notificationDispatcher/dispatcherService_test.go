package notificationDispatcher

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/atrian/go-notify-customer/internal/dto"
	"github.com/atrian/go-notify-customer/internal/interfaces"
	"github.com/atrian/go-notify-customer/pkg/logger"
)

var (
	_ dispatcherConfig      = (*configMock)(nil)
	_ eventService          = (*eventMock)(nil)
	_ templateService       = (*templateMock)(nil)
	_ contactVault          = (*contactMock)(nil)
	_ interfaces.AmpqClient = (*ampqMock)(nil)
)

type DispatcherServiceTestSuite struct {
	suite.Suite
	config     dispatcherConfig
	dispatcher *Dispatcher
	ampq       interfaces.AmpqClient
	inputCh    chan dto.Notification
	outChan    chan string
}

func (suite *DispatcherServiceTestSuite) SetupSuite() {
	inputCh := make(chan dto.Notification)
	suite.inputCh = inputCh
	outChan := make(chan string)
	suite.ampq = newAmpqMock(outChan)
	suite.outChan = outChan

	suite.config = &configMock{}
	suite.dispatcher = New(context.Background(),
		inputCh,
		suite.config,
		NewDispatcherServiceFacade(&contactMock{}, &templateMock{}, &eventMock{}),
		suite.ampq,
		logger.NewZapLogger())

	suite.dispatcher.Start()
}

func (suite *DispatcherServiceTestSuite) Test_dispatch() {
	personUUID := uuid.New()
	notificationUUID := uuid.New()

	suite.inputCh <- dto.Notification{
		Index:            111,
		NotificationUUID: notificationUUID,
		EventUUID:        uuid.UUID{},
		PersonUUIDs: []uuid.UUID{
			personUUID,
		},
		MessageParams: nil,
		Priority:      222,
	}

	message := <-suite.outChan
	expected := fmt.Sprintf("queue: message_queue, message: {\"notification_uuid\":\"%v\",\"person_uuid\":\"%v\",\"text\":\"test message\",\"channel\":\"sms\",\"destination_address\":\"888\"}", notificationUUID, personUUID)
	assert.Equal(suite.T(), expected, message)
}

type configMock struct{}

func (c *configMock) GetNotificationQueue() string {
	return "message_queue"
}

type eventMock struct{}

func (e *eventMock) FindById(ctx context.Context, eventUUID uuid.UUID) (dto.Event, error) {
	return dto.Event{
		EventUUID:            uuid.New(),
		Title:                "test event",
		Description:          "dd",
		DefaultPriority:      0,
		NotificationChannels: []string{"sms"},
	}, nil
}

type contactMock struct{}

func (c *contactMock) FindByPersonUUID(ctx context.Context, personUUID uuid.UUID) (dto.PersonContacts, error) {
	return dto.PersonContacts{
		PersonUUID: personUUID,
		Contacts: []dto.Contact{
			{
				Channel:     "sms",
				Destination: "888",
			},
		},
	}, nil
}

func (c *contactMock) Stop() error {
	return nil
}

type templateMock struct{}

func (t *templateMock) FindByEventId(ctx context.Context, eventUUID uuid.UUID) ([]dto.Template, error) {
	return []dto.Template{
		dto.Template{
			TemplateUUID: uuid.New(),
			EventUUID:    uuid.New(),
			Title:        "test",
			Description:  "none",
			Body:         "test message",
			ChannelType:  "sms",
		},
	}, nil
}

type ampqMock struct {
	outputChan chan string
}

func newAmpqMock(outputChan chan string) *ampqMock {
	am := ampqMock{outputChan: outputChan}
	return &am
}

func (a *ampqMock) Connect(dsn string) error {
	return nil
}

func (a *ampqMock) Reconnect() error {
	return nil
}

func (a *ampqMock) MigrateDurableQueues(queues ...string) {
	return
}

func (a *ampqMock) Channel() *amqp.Channel {
	return nil
}

func (a *ampqMock) Consume(queue string) (<-chan amqp.Delivery, error) {
	return nil, nil
}

func (a *ampqMock) Publish(queue string, msgBody []byte) error {
	res := fmt.Sprintf("queue: %v, message: %v", queue, string(msgBody))
	a.outputChan <- res
	return nil
}

func (a *ampqMock) Stop() {
	return
}

// Для запуска через Go test
func TestDispatcherServiceTestSuite(t *testing.T) {
	suite.Run(t, new(DispatcherServiceTestSuite))
}
