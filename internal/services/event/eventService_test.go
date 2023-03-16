package event

import (
	"context"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/atrian/go-notify-customer/internal/services/event/entity"
)

type TestSuite struct {
	suite.Suite
	service Service
	events  []entity.Event
}

func (suite *TestSuite) SetupSuite() {
	suite.events = []entity.Event{
		{
			EventUUID:            uuid.New(),
			Title:                "New appointment",
			Description:          "Test description",
			DefaultPriority:      100,
			NotificationChannels: []string{"sms"},
		}, {
			EventUUID:            uuid.New(),
			Title:                "Appointment canceled",
			Description:          "Test description 2",
			DefaultPriority:      999,
			NotificationChannels: []string{"sms", "email"},
		},
	}
}

func (suite *TestSuite) SetupTest() {
	suite.service.storage = NewMemoryStorage()
	_, _ = suite.service.StoreBatch(context.TODO(), suite.events)
}

func (suite *TestSuite) TestService_All() {
	result := suite.service.All(context.TODO())

	if !reflect.DeepEqual(result, suite.events) {
		suite.T().Errorf("Expected events %v, got %v", suite.events, result)
	}
}

func (suite *TestSuite) TestService_Store() {
	newEvent := entity.Event{
		Title:                "New event",
		Description:          "For test",
		DefaultPriority:      1,
		NotificationChannels: nil,
	}

	storeResult, err := suite.service.Store(context.TODO(), newEvent)
	assert.NoError(suite.T(), err)

	// При сохроанении событию выдается UUID, сбрасываем его для сравнения в тесте
	storeResult.EventUUID = uuid.UUID{}
	assert.Equal(suite.T(), storeResult, newEvent)
}

func (suite *TestSuite) TestService_StoreBatch() {
	newEvents := []entity.Event{
		{
			Title:                "New event 1",
			Description:          "For test 1",
			DefaultPriority:      15,
			NotificationChannels: []string{"sms"},
		},
		{
			Title:                "New event 2",
			Description:          "For test 2",
			DefaultPriority:      19,
			NotificationChannels: []string{"email"},
		},
	}

	storeResult, err := suite.service.StoreBatch(context.TODO(), newEvents)
	assert.NoError(suite.T(), err)

	// При сохроанении событию выдается UUID, сбрасываем его для сравнения в тесте
	storeResult[0].EventUUID = uuid.UUID{}
	storeResult[1].EventUUID = uuid.UUID{}
	assert.Equal(suite.T(), storeResult, newEvents)
}

func (suite *TestSuite) TestService_Update() {
	eventForUpdate := suite.events[0]
	eventForUpdate.Title = "Updated Title"

	result, err := suite.service.Update(context.TODO(), eventForUpdate)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), result, eventForUpdate)
}

func (suite *TestSuite) TestService_DeleteById() {
	eventUUID := suite.events[0].EventUUID

	err := suite.service.DeleteById(context.TODO(), eventUUID)
	assert.NoError(suite.T(), err)

	err = suite.service.DeleteById(context.TODO(), eventUUID)
	assert.ErrorIs(suite.T(), err, NotFound)
}

func (suite *TestSuite) TestService_FindById() {
	eventUUID := suite.events[0].EventUUID

	result, err := suite.service.FindById(context.TODO(), eventUUID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), result, suite.events[0])

	// Поиск несуществующего события
	err = suite.service.DeleteById(context.TODO(), uuid.New())
	assert.ErrorIs(suite.T(), err, NotFound)
}

// Для запуска через Go test
func TestEventServiceTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
