// Package stat Сервис статистики отправки сообщений. Общается с внутренними сервисами посредством канала
// statChan chan dto.Stat
//
// Структура передачи данных между слоями приложения dto.Stat:
//
//	 type Stat struct {
//		StatUUID         uuid.UUID  `json:"stat_uuid"`         // StatUUID id записи статистики
//		PersonUUID       uuid.UUID  `json:"person_uuid"`       // PersonUUID связь отправленного уведомления с клиентом
//		NotificationUUID uuid.UUID  `json:"notification_uuid"` // NotificationUUID связь с уведомлением
//		CreatedAt        string     `json:"created_at"`        // CreatedAt дата и время отправки
//		Status           StatStatus `json:"status"`            // Status статус отправки
//	}
//
// Возможные статусы dto.Stat
//
//	Sent       StatStatus = iota + 1 // Уведомление отправлено
//	Failed                           // Ошибка отправки
//	BadChannel                       // Канал отправки не поддерживается
package stat

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/atrian/go-notify-customer/internal/dto"
	"github.com/atrian/go-notify-customer/internal/interfaces"
)

var _ interfaces.StatService = (*Service)(nil)

// Service структура содержит канал для получения статистикт отправок
// хранилище и логгер с интерфейсом interfaces.Logger
type Service struct {
	statChan <-chan dto.Stat
	storage  Storager
	logger   interfaces.Logger
}

func New(statChan chan dto.Stat, logger interfaces.Logger) *Service {
	s := Service{
		statChan: statChan,
		storage:  NewMemoryStorage(),
		logger:   logger,
	}

	return &s
}

// Start стартовые процедуры для сервиса
func (s Service) Start(ctx context.Context) {
	// слушаем канал statChan в который другие сервисы передают данные об отправках
	go func(ctx context.Context, statChan <-chan dto.Stat) {
		s.logger.Info("Stat listener is UP")
		for {
			select {
			case stat := <-statChan:
				_ = s.Store(ctx, stat) // сохраняем полученную статистику

			case <-ctx.Done(): // завершение по контексту
				return

			default: // обеспечиваем простой пока ждем данных в канале
			}
		}
	}(ctx, s.statChan)
}

// Stop корректное завершение работы
func (s Service) Stop() {
	s.logger.Info("Stat service stopped")
}

// All возвращает все хранящиеся шаблоны
func (s Service) All(ctx context.Context) []dto.Stat {
	res, err := s.storage.All(ctx)
	if err != nil {
		s.logger.Error("Stat service storage.All err", err)
	}

	if res == nil {
		return []dto.Stat{}
	}

	return res
}

// Store сохранение шаблона в харнилище
func (s Service) Store(ctx context.Context, stat dto.Stat) error {
	stat.StatUUID = uuid.New()
	stat.CreatedAt = time.Now().Format(dateTimeFormat) // сохраняем время записи

	return s.storage.Store(ctx, stat)
}

// FindByPersonUUID возвращает статистику по получателю уведомления
func (s Service) FindByPersonUUID(ctx context.Context, personUUID uuid.UUID) ([]dto.Stat, error) {
	return s.storage.GetByPersonId(ctx, personUUID)
}

// FindByNotificationId возвращает статистику по конкретному уведомлению
func (s Service) FindByNotificationId(ctx context.Context, notificationUUID uuid.UUID) ([]dto.Stat, error) {
	return s.storage.GetByNotificationId(ctx, notificationUUID)
}
