package notificationDispatcher

import (
	"context"
	"regexp"

	"github.com/google/uuid"

	"github.com/atrian/go-notify-customer/internal/dto"
	"github.com/atrian/go-notify-customer/internal/interfaces"
)

var (
	_  serviceGateway = (*ServiceFacade)(nil)
	re                = regexp.MustCompile(`(?m)\[([a-zA-Z]+)]`)
)

type contactVault interface {
	FindByPersonUUID(ctx context.Context, personUUID uuid.UUID) (dto.PersonContacts, error)
	Stop() error
}

type ServiceFacade struct {
	contact  contactVault
	template interface {
		FindByEventId(ctx context.Context, eventUUID uuid.UUID) ([]dto.Template, error)
	}
	event interface {
		FindById(ctx context.Context, eventUUID uuid.UUID) (dto.Event, error)
	}
}

func NewDispatcherServiceFacade(template interfaces.TemplateService, event interfaces.EventService) *ServiceFacade {
	f := ServiceFacade{
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
	// TODO замена подстановок в тексте сообщения
	return ""
}
