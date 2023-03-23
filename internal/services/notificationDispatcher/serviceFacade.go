package notificationDispatcher

import (
	"context"
	"regexp"

	"github.com/google/uuid"

	"github.com/atrian/go-notify-customer/internal/dto"
)

var (
	_ serviceGateway = (*ServiceFacade)(nil)
	// regexp для поиска параметров [param1] [param2] в теле сообщения
	templateRe = regexp.MustCompile(`(?m)\[([a-zA-Z]+\d*)]`)
	// regexp для замены множественных пробелов на один
	spaceRe = regexp.MustCompile(`\s+`)
)

// contactVault интерфейс клиента хранилища контакных данных
type contactVault interface {
	// FindByPersonUUID поиск по uuid
	FindByPersonUUID(ctx context.Context, personUUID uuid.UUID) (dto.PersonContacts, error)
	// Stop остановка клиента
	Stop() error
}

// templateService контракт на сервис шаблонов. Доступен поиск по бизнес событию
type templateService interface {
	FindByEventId(ctx context.Context, eventUUID uuid.UUID) ([]dto.Template, error)
}

// eventService контракт на сервис бизнес событий.
type eventService interface {
	FindById(ctx context.Context, eventUUID uuid.UUID) (dto.Event, error)
}

// ServiceFacade сервисный фасад для нужд notificationDispatcher
type ServiceFacade struct {
	contact  contactVault
	template templateService
	event    eventService
}

func NewDispatcherServiceFacade(contact contactVault, template templateService, event eventService) *ServiceFacade {
	f := ServiceFacade{
		contact:  contact,
		event:    event,
		template: template,
	}

	return &f
}

func (f *ServiceFacade) getContacts(ctx context.Context, personUUIDs []uuid.UUID) ([]dto.PersonContacts, error) {
	contacts := make([]dto.PersonContacts, 0, len(personUUIDs))

	for _, personUuid := range personUUIDs {
		contact, _ := f.contact.FindByPersonUUID(ctx, personUuid)
		contacts = append(contacts, contact)
	}

	return contacts, nil
}

func (f *ServiceFacade) getTemplates(ctx context.Context, eventUuid uuid.UUID) ([]dto.Template, error) {
	return f.template.FindByEventId(ctx, eventUuid)
}

func (f *ServiceFacade) getEvent(ctx context.Context, eventUuid uuid.UUID) (dto.Event, error) {
	return f.event.FindById(ctx, eventUuid)
}

func (f *ServiceFacade) prepareTemplate(template string, replaces []dto.MessageParam) string {
	// Собираем таблицу замен
	replaceDict := make(map[string]string, len(replaces))
	for _, el := range replaces {
		replaceDict[el.Key] = el.Value
	}

	// Заменяем в строке все key1 в формате [key1] на значение replaceDict[key1]
	result := templateRe.ReplaceAllFunc([]byte(template), func(bytes []byte) []byte {
		key := string(bytes[1 : len(bytes)-1])
		if val, ok := replaceDict[key]; ok {
			return []byte(val)
		}
		return []byte{}
	})

	return spaceRe.ReplaceAllString(string(result), " ")
}
