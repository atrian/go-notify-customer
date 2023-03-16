package template

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/atrian/go-notify-customer/internal/services/template/entity"
)

type TemplateTestSuite struct {
	suite.Suite
	service   Service
	templates []entity.Template
}

func (suite *TemplateTestSuite) SetupSuite() {
	suite.templates = []entity.Template{
		{
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

func (suite *TemplateTestSuite) SetupTest() {
	suite.service.storage = NewMemoryStorage()
	_, _ = suite.service.StoreBatch(context.TODO(), suite.templates)
}

func (suite *TemplateTestSuite) TestService_All() {
	result := suite.service.All(context.TODO())

	assert.Equal(suite.T(), len(suite.templates), len(result))
}

func (suite *TemplateTestSuite) TestService_Store() {
	newEvent := entity.Template{
		EventUUID:   uuid.UUID{},
		Title:       "Test Title",
		Description: "Description",
		Body:        "Body",
		ChannelType: "sms",
	}

	storeResult, err := suite.service.Store(context.TODO(), newEvent)
	assert.NoError(suite.T(), err)

	// При сохроанении событию выдается UUID, сбрасываем его для сравнения в тесте
	storeResult.EventUUID = uuid.UUID{}
	assert.Equal(suite.T(), storeResult, newEvent)
}

func (suite *TemplateTestSuite) TestService_StoreBatch() {
	newTemplates := []entity.Template{
		{
			EventUUID:   uuid.New(),
			Title:       "Test 1",
			Description: "Description",
			Body:        "Body",
			ChannelType: "ChannelType",
		}, {
			EventUUID:   uuid.New(),
			Title:       "Test 2",
			Description: "Description",
			Body:        "Body",
			ChannelType: "ChannelType",
		},
	}

	storeResult, err := suite.service.StoreBatch(context.TODO(), newTemplates)
	assert.NoError(suite.T(), err)

	// При сохроанении событию выдается UUID, сбрасываем его для сравнения в тесте
	storeResult[0].TemplateUUID = uuid.UUID{}
	storeResult[1].TemplateUUID = uuid.UUID{}
	assert.Equal(suite.T(), storeResult, newTemplates)
}

func (suite *TemplateTestSuite) TestService_Update() {
	itemForUpdate := suite.templates[0]
	itemForUpdate.Title = "Updated Title"

	result, err := suite.service.Update(context.TODO(), itemForUpdate)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), result, itemForUpdate)
}

func (suite *TemplateTestSuite) TestService_DeleteById() {
	templateUUID := suite.templates[0].TemplateUUID

	err := suite.service.DeleteById(context.TODO(), templateUUID)
	assert.NoError(suite.T(), err)

	err = suite.service.DeleteById(context.TODO(), templateUUID)
	assert.ErrorIs(suite.T(), err, NotFound)
}

func (suite *TemplateTestSuite) TestService_FindById() {
	templateUUID := suite.templates[0].TemplateUUID

	result, err := suite.service.FindById(context.TODO(), templateUUID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), result, suite.templates[0])

	// Поиск несуществующего шаблона
	_, err = suite.service.FindById(context.TODO(), uuid.New())
	assert.ErrorIs(suite.T(), err, NotFound)
}

func (suite *TemplateTestSuite) TestService_FindByEventID() {
	eventUUID := suite.templates[0].EventUUID

	result, err := suite.service.FindByEventId(context.TODO(), eventUUID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), result, suite.templates[0])

	// Поиск несуществующего события
	_, err = suite.service.FindByEventId(context.TODO(), uuid.New())
	assert.ErrorIs(suite.T(), err, NotFound)
}

// Для запуска через Go test
func TestTemplateServiceTestSuite(t *testing.T) {
	suite.Run(t, new(TemplateTestSuite))
}
