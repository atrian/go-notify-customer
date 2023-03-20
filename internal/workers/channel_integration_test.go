//go:build integration
// +build integration

package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/atrian/go-notify-customer/internal/dto"
	"github.com/dhui/dktest"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/atrian/go-notify-customer/pkg/ampq"
	"github.com/atrian/go-notify-customer/pkg/logger"
)

var (
	dsn  = "amqp://guest:guest@localhost:5672/"
	opts = dktest.Options{
		ReadyTimeout: 15 * time.Second,
		Hostname:     "localhost",
		PortBindings: nat.PortMap{
			nat.Port("5672/tcp"): []nat.PortBinding{{
				HostIP:   "0.0.0.0",
				HostPort: "5672",
			}},
		},
		ReadyFunc: isReady,
	}

	image = "rabbitmq:3-management-alpine"

	_ config         = (*configMock)(nil)
	_ channelService = (*smsServiceMock)(nil)
)

func TestRabbitAlive(t *testing.T) {
	dktest.Run(t, image, opts, func(t *testing.T, c dktest.ContainerInfo) {
		client := ampq.NewWithConnection(dsn, logger.NewZapLogger())
		client.MigrateDurableQueues("test")

		err := client.Publish("test", []byte("message"))
		assert.NoError(t, err)

		_, err = client.Consume("test")
		assert.NoError(t, err)
	})
}

func TestChannelWorker_Start(t *testing.T) {
	dktest.Run(t, image, opts, func(t *testing.T, c dktest.ContainerInfo) {
		client := ampq.NewWithConnection(dsn, logger.NewZapLogger())
		conf := configMock{}
		client.MigrateDurableQueues(conf.GetNotificationQueue(), conf.GetFailedWorksQueue())
		zapLogger := logger.NewZapLogger()
		sendStatChan := make(chan dto.Stat)
		outputChan := make(chan string)
		done := make(chan struct{})

		// Создаем нового воркера
		worker := NewChannelWorker(conf, client, sendStatChan, context.Background(), zapLogger)
		worker.ReloadService("sms", newSmsMock(outputChan))

		// Run worker
		go func() {
			worker.Start(conf.GetNotificationQueue(), "", "")
			defer worker.Stop()
		}()

		// Result checker
		go func() {
			// Сравниваем результат 1
			result := <-outputChan
			assert.Equal(t, "Test rabbit:+79876543210", result)

			// Сравниваем результат 2
			result = <-outputChan
			assert.Equal(t, "Test rabbit second:+123456", result)

			done <- struct{}{}
		}()

		// Stat receiver stub
		go func() {
			for {
				select {
				case <-sendStatChan:
					// do something with stat data
				default:
					// do nothing
				}
			}
		}()

		// отправляем в очередь тестовое сообщение
		err := client.Publish("test", []byte("message"))
		notificationUUID := uuid.New()
		personUUID := uuid.New()

		message := dto.Message{
			NotificationUUID:   notificationUUID,
			PersonUUID:         personUUID,
			Text:               "Test rabbit",
			Channel:            "sms",
			DestinationAddress: "+79876543210",
		}
		message2 := dto.Message{
			NotificationUUID:   notificationUUID,
			PersonUUID:         personUUID,
			Text:               "Test rabbit second",
			Channel:            "sms",
			DestinationAddress: "+123456",
		}

		jsonMessage, err := json.Marshal(message)
		assert.NoError(t, err)
		jsonMessage2, err := json.Marshal(message2)
		assert.NoError(t, err)

		// отправляем первое сообщение
		err = client.Publish(conf.GetNotificationQueue(), jsonMessage)
		assert.NoError(t, err)
		// отправляем второе сообщение
		err = client.Publish(conf.GetNotificationQueue(), jsonMessage2)
		assert.NoError(t, err)

		<-done
	})
}

func isReady(ctx context.Context, c dktest.ContainerInfo) bool {
	client := ampq.New(dsn, logger.NewFatalZapLogger())
	err := client.Connect(dsn)

	if err != nil {
		return false
	}

	return true
}

type configMock struct{}

func (c configMock) GetAmpqDSN() string {
	return dsn
}

func (c configMock) GetNotificationQueue() string {
	return "notifications"
}

func (c configMock) GetFailedWorksQueue() string {
	return "failed"
}

func (c configMock) GetMailConfig() {
	//TODO implement me
	panic("implement me")
}

func (c configMock) GetTwilioConfig() {
	//TODO implement me
	panic("implement me")
}

type smsServiceMock struct {
	output chan<- string
}

func newSmsMock(output chan string) *smsServiceMock {
	s := smsServiceMock{output: output}
	return &s
}

func (s *smsServiceMock) SendMessage(message string, destination string) error {
	s.output <- fmt.Sprintf("%v:%v", message, destination)
	return nil
}
