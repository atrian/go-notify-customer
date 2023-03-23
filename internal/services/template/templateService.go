// Package template Сервис работы с шаблонами сообщений для нужд
// REST CRUD и формирования сообщений
// при отправке уведомления во входящих данных указывается EventUUID
// по которому находится шаблон сообщения для разных каналов доставки
//
// Формат передачи между слоями приложения dto.Template
// Формат входящего json см. dto.IncomingTemplate
//
//	type Template struct {
//		TemplateUUID uuid.UUID `json:"template_uuid"`         // TemplateUUID - id шаблона
//		EventUUID    uuid.UUID `json:"event_uuid"`            // EventUUID связь с UUID бизнес события
//		Title        string    `json:"title"`                 // Title название шаблона
//		Description  string    `json:"description,omitempty"` // Description описание шаблона
//		Body         string    `json:"body"`                  // Body тело шаблона
//		ChannelType  string    `json:"channel_type"`          // ChannelType связь с каналом отправки
//	}
//
// В поле Body (тело шаблона) можно указывать места для подстановки.
// Пример "Ваша запись на [date] подтверждена. [company]"
// Плейсхолдер должен начинаться с квадратной скобки [ и заканчиваться закрывающейся квадратной скобкой ]
// Внутри допустимы латинские буквы в нижнем и верхнем регистре.
// Можно указывать цифры ПОСЛЕ обозначения буквенного ключа. Другие символы запрещены.
//
// Правильно: [param1], [param2]
// Не правильно: [1param], [param 1], [param_1], [para1m], [$3fds] и тд.
//
// Входяшее уведомление приходит в JSON отражающим следующую структуру
// Плейсхолдеры будут заменены при подготовке сообщения данными из []MessageParams
//
//	type IncomingNotification struct {
//			EventUUID     uuid.UUID      `json:"event_uuid"`               // uuid бизнес события
//			PersonUUIDs   []uuid.UUID    `json:"person_uuids"`             //
//			MessageParams []MessageParam `json:"message_params,omitempty"` //
//			Priority      uint           `json:"priority,omitempty"`       //
//	}
//
// Подробнее см. dto.IncomingNotification
package template

import (
	"context"

	"github.com/google/uuid"

	"github.com/atrian/go-notify-customer/internal/dto"
	"github.com/atrian/go-notify-customer/internal/interfaces"
)

// Service содержит хранилище данных и логгер удовлетворяющий интерфейсу interfaces.Logger
type Service struct {
	storage Storager
	logger  interfaces.Logger
}

func New(logger interfaces.Logger) *Service {
	s := Service{
		storage: NewMemoryStorage(),
		logger:  logger,
	}
	return &s
}

// Start стартовые процедуры сервиса
func (s Service) Start(ctx context.Context) {
	s.logger.Info("Template service started")
}

// Stop корректное завершение работы
func (s Service) Stop() {
	s.logger.Info("Template service stopped")
}

// All возвращает все хранящиеся шаблоны
func (s Service) All(ctx context.Context) []dto.Template {
	templates, err := s.storage.All(ctx)

	if err != nil {
		s.logger.Error("Template service storage.All err", err)
	}

	if templates == nil {
		return []dto.Template{}
	}

	return templates
}

// Store сохранение шаблона в харнилище
func (s Service) Store(ctx context.Context, template dto.Template) (dto.Template, error) {
	template.TemplateUUID = uuid.New()

	err := s.storage.Store(ctx, template)
	if err != nil {
		s.logger.Error("Template service storage.Store err", err)
		return dto.Template{}, err
	}

	return template, nil
}

// StoreBatch массовое сохранение шаблонов.
// В данной версии не используется хендлерами, задел на будущее.
func (s Service) StoreBatch(ctx context.Context, templates []dto.Template) ([]dto.Template, error) {
	for i := 0; i < len(templates); i++ {
		templates[i].TemplateUUID = uuid.New()
		err := s.storage.Store(ctx, templates[i])
		if err != nil {
			s.logger.Error("Template service storage.Store err", err)
		}
	}

	return templates, nil
}

// Update обновление шаблона
func (s Service) Update(ctx context.Context, template dto.Template) (dto.Template, error) {
	err := s.storage.Store(ctx, template)
	if err != nil {
		s.logger.Error("Template service storage.Store err (Update)", err)
		return dto.Template{}, err
	}

	return template, nil
}

// FindById поиск шаблона по его uuid
func (s Service) FindById(ctx context.Context, templateUUID uuid.UUID) (dto.Template, error) {
	return s.storage.GetById(ctx, templateUUID)
}

// FindByEventId поиск шаблонов по uuid бизнес события
func (s Service) FindByEventId(ctx context.Context, eventUUID uuid.UUID) ([]dto.Template, error) {
	return s.storage.GetByEventId(ctx, eventUUID)
}

// DeleteById удаление шаблона из хранилища. Hard delete!
func (s Service) DeleteById(ctx context.Context, templateUUID uuid.UUID) error {
	return s.storage.DeleteById(ctx, templateUUID)
}
