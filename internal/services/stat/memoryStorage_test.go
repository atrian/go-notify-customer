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
	result.CreatedAt = "" // сбрасываем данные о времени записи для сравнения
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), result, suite.stats[0])
}

func (suite *MemoryStorageTestSuite) Test_GetByPersonId() {
	// Запрос несуществующего объекта
	_, err := suite.storage.GetByPersonId(context.TODO(), uuid.New())
	assert.ErrorIs(suite.T(), err, NotFound)

	// Запрос существующего объекта
	result, err := suite.storage.GetByPersonId(context.TODO(), suite.stats[0].PersonUUID)
	assert.NoError(suite.T(), err)
	// Исходя из тестовых данных ожидаем 2 элемента на 1 идентификатор клиента
	assert.Equal(suite.T(), len(result), 2)
}

func (suite *MemoryStorageTestSuite) Test_All() {
	result, err := suite.storage.All(context.TODO())
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), len(result), len(suite.stats))
}

// Для запуска через Go test
func TestMemoryStorageSuite(t *testing.T) {
	suite.Run(t, new(MemoryStorageTestSuite))
}
