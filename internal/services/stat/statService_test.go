package stat

import (
	"context"
	"testing"

	"github.com/atrian/go-notify-customer/internal/dto"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type StatTestSuite struct {
	suite.Suite
	stats   []dto.Stat
	service *Service
	cancel  context.CancelFunc
}

func (suite *StatTestSuite) SetupSuite() {
	personOneUUID := uuid.New()
	personTwoUUID := uuid.New()

	suite.stats = []dto.Stat{
		{
			NotificationUUID: uuid.New(),
			PersonUUID:       personOneUUID,
			Status:           dto.Sent,
		}, {
			NotificationUUID: uuid.New(),
			PersonUUID:       personOneUUID,
			Status:           dto.Failed,
		}, {
			NotificationUUID: uuid.New(),
			PersonUUID:       personTwoUUID,
			Status:           dto.Sent,
		}, {
			NotificationUUID: uuid.New(),
			PersonUUID:       personTwoUUID,
			Status:           dto.Sent,
		},
	}
}

func (suite *StatTestSuite) SetupTest() {
	ctx, cancel := context.WithCancel(context.TODO())
	suite.cancel = cancel
	statChan := make(chan dto.Stat)

	suite.service = New(ctx, statChan)

	suite.service.Start()

	for i := 0; i < len(suite.stats); i++ {
		statChan <- suite.stats[i]
	}
}

func (suite *StatTestSuite) Test_All() {
	result := suite.service.All(context.TODO())

	assert.Equal(suite.T(), len(suite.stats), len(result))
}

func (suite *StatTestSuite) Test_Store() {
	newStat := dto.Stat{
		PersonUUID:       uuid.New(),
		NotificationUUID: uuid.New(),
		Status:           dto.Sent,
	}

	err := suite.service.Store(newStat)
	assert.NoError(suite.T(), err)

	// Запрос существующего объекта
	result, err := suite.service.FindByEventUUID(context.TODO(), newStat.NotificationUUID)
	result.CreatedAt = "" // опускаем временную метку для сравнения
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), newStat, result)
}

func (suite *StatTestSuite) Test_FindByPersonUUID() {
	// Запрос несуществующего объекта
	_, err := suite.service.FindByPersonUUID(context.TODO(), uuid.New())
	assert.ErrorIs(suite.T(), err, NotFound)

	// Запрос существующего объекта
	result, err := suite.service.FindByPersonUUID(context.TODO(), suite.stats[0].PersonUUID)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2, len(result)) // по условиям теста у нас 2 события на человека
}

func (suite *StatTestSuite) Test_FindByEventUUID() {
	// Запрос несуществующего объекта
	_, err := suite.service.FindByEventUUID(context.TODO(), uuid.New())
	assert.ErrorIs(suite.T(), err, NotFound)

	// Запрос существующего объекта
	result, err := suite.service.FindByEventUUID(context.TODO(), suite.stats[0].NotificationUUID)
	result.CreatedAt = "" // опускаем временную метку для сравнения
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), suite.stats[0], result)
}

// Для запуска через Go test
func TestStatServiceTestSuite(t *testing.T) {
	suite.Run(t, new(StatTestSuite))
}
