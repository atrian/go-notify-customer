package template

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/atrian/go-notify-customer/internal/dto"
)

type MemoryStorageTestSuite struct {
	suite.Suite
	storage   Storager
	templates []dto.Template
}

func (suite *MemoryStorageTestSuite) SetupSuite() {
	persistentEventUUID := uuid.New()

	suite.templates = []dto.Template{
		{
			TemplateUUID: uuid.New(),
			EventUUID:    persistentEventUUID,
			Title:        "Test 3",
			Description:  "Description 1",
			Body:         "Body [param] 1",
			ChannelType:  "sms",
		}, {
			TemplateUUID: uuid.New(),
			EventUUID:    persistentEventUUID,
			Title:        "Test 4",
			Description:  "Description 1",
			Body:         "Body [param] 1",
			ChannelType:  "sms",
		}, {
			TemplateUUID: uuid.New(),
			EventUUID:    uuid.New(),
			Title:        "Test 1",
			Description:  "Description 1",
			Body:         "Body [param] 1",
			ChannelType:  "sms",
		}, {
			TemplateUUID: uuid.New(),
			EventUUID:    uuid.New(),
			Title:        "Test 2",
			Description:  "Description 2",
			Body:         "Body [param] 2",
			ChannelType:  "email",
		},
	}
}

func (suite *MemoryStorageTestSuite) SetupTest() {
	suite.storage = NewMemoryStorage()
	for i := 0; i < len(suite.templates); i++ {
		_ = suite.storage.Store(context.TODO(), suite.templates[i])
	}
}

func (suite *MemoryStorageTestSuite) Test_GetById() {
	// Запрос несуществующего объекта
	_, err := suite.storage.GetById(context.TODO(), uuid.New())
	assert.ErrorIs(suite.T(), err, NotFound)

	// Запрос существующего объекта
	result, err := suite.storage.GetById(context.TODO(), suite.templates[0].TemplateUUID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), result, suite.templates[0])
}

func (suite *MemoryStorageTestSuite) Test_GetByEventId() {
	// Запрос несуществующего объекта
	_, err := suite.storage.GetByEventId(context.TODO(), uuid.New())
	assert.ErrorIs(suite.T(), err, NotFound)

	// Запрос существующего объекта
	result, err := suite.storage.GetByEventId(context.TODO(), suite.templates[0].EventUUID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2, len(result)) // По условию теста 2 результата с нужным EventUUID
	assert.Equal(suite.T(), suite.templates[0].EventUUID, result[0].EventUUID)
	assert.Equal(suite.T(), suite.templates[0].EventUUID, result[1].EventUUID)
}

func (suite *MemoryStorageTestSuite) Test_All() {
	result, err := suite.storage.All(context.TODO())
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), len(result), len(suite.templates))
}

func (suite *MemoryStorageTestSuite) Test_DeleteById() {
	err := suite.storage.DeleteById(context.TODO(), suite.templates[0].TemplateUUID)
	assert.NoError(suite.T(), err)

	err = suite.storage.DeleteById(context.TODO(), suite.templates[0].TemplateUUID)
	assert.ErrorIs(suite.T(), err, NotFound)
}

func (suite *MemoryStorageTestSuite) Test_Update() {
	template := suite.templates[0]
	template.Title = "Updated field"

	err := suite.storage.Update(context.TODO(), template)
	assert.NoError(suite.T(), err)

	updated, err := suite.storage.GetById(context.TODO(), template.TemplateUUID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), template, updated)
}

// Для запуска через Go test
func TestMemoryStorageSuite(t *testing.T) {
	suite.Run(t, new(MemoryStorageTestSuite))
}
