package stat

import (
	"context"
	"testing"

	"github.com/atrian/go-notify-customer/internal/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/google/uuid"
)

type MemoryStorageTestSuite struct {
	suite.Suite
	storage Storager
	stats   []dto.Stat
}

func (suite *MemoryStorageTestSuite) SetupSuite() {
	personOneUUID := uuid.New()
	personTwoUUID := uuid.New()

	suite.stats = []dto.Stat{
		{
			StatUUID:         uuid.New(),
			NotificationUUID: uuid.New(),
			PersonUUID:       personOneUUID,
			Status:           dto.Sent,
		}, {
			StatUUID:         uuid.New(),
			NotificationUUID: uuid.New(),
			PersonUUID:       personOneUUID,
			Status:           dto.Failed,
		}, {
			StatUUID:         uuid.New(),
			NotificationUUID: uuid.New(),
			PersonUUID:       personTwoUUID,
			Status:           dto.Sent,
		}, {
			StatUUID:         uuid.New(),
			NotificationUUID: uuid.New(),
			PersonUUID:       personTwoUUID,
			Status:           dto.Sent,
		},
	}
}

func (suite *MemoryStorageTestSuite) SetupTest() {
	suite.storage = NewMemoryStorage()
	for i := 0; i < len(suite.stats); i++ {
		_ = suite.storage.Store(context.TODO(), suite.stats[i])
	}
}

func (suite *MemoryStorageTestSuite) Test_GetByNotificationId() {
	// Запрос несуществующего объекта
	_, err := suite.storage.GetByNotificationId(context.TODO(), uuid.New())
	assert.ErrorIs(suite.T(), err, NotFound)

	// Запрос существующего объекта
	result, err := suite.storage.GetByNotificationId(context.TODO(), suite.stats[0].NotificationUUID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, len(result))
}

func (suite *MemoryStorageTestSuite) Test_GetByPersonId() {
	// Запрос несуществующего объекта
	_, err := suite.storage.GetByPersonId(context.TODO(), uuid.New())
	assert.ErrorIs(suite.T(), err, NotFound)

	// Запрос существующего объекта
	result, err := suite.storage.GetByPersonId(context.TODO(), suite.stats[0].PersonUUID)
	assert.NoError(suite.T(), err)
	// Исходя из тестовых данных ожидаем 2 элемента на 1 идентификатор клиента
	assert.Equal(suite.T(), 2, len(result))
}

func (suite *MemoryStorageTestSuite) Test_All() {
	result, err := suite.storage.All(context.TODO())
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), len(suite.stats), len(result))
}

// Для запуска через Go test
func TestMemoryStorageSuite(t *testing.T) {
	suite.Run(t, new(MemoryStorageTestSuite))
}
