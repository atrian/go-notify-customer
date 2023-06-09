package notify

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/atrian/go-notify-customer/config"
	"github.com/atrian/go-notify-customer/internal/dto"
	"github.com/atrian/go-notify-customer/internal/interfaces"
	"github.com/atrian/go-notify-customer/internal/notify/handlers"
	"github.com/atrian/go-notify-customer/internal/notify/router"
	"github.com/atrian/go-notify-customer/internal/services/event"
	"github.com/atrian/go-notify-customer/internal/services/notificationDispatcher"
	"github.com/atrian/go-notify-customer/internal/services/notify"
	"github.com/atrian/go-notify-customer/internal/services/stat"
	"github.com/atrian/go-notify-customer/internal/services/template"
	"github.com/atrian/go-notify-customer/internal/workers"
	"github.com/atrian/go-notify-customer/pkg/ampq"
	"github.com/atrian/go-notify-customer/pkg/logger"
)

type App struct {
	services         services
	config           config.Config
	notificationChan chan dto.Notification
	statChan         chan dto.Stat
	logger           interfaces.Logger
}

// services - регистр всех доступных сервисов
type services struct {
	notificationService    interfaces.NotificationService         // notificationService приоритезация и органичение уведомлений
	notificationDispatcher interfaces.NotificationDispatchService // notificationDispatcher отправка уведомлений в BUS
	eventService           interfaces.EventService                // eventService CRUD сервис для бизнес событий
	templateService        interfaces.TemplateService             // templateService CRUD сервис для шаблонов событий
	statisticService       interfaces.StatService                 // statisticService сервис статистики отправки
}

func New() App {
	// логгер приложения
	appLogger := logger.NewZapLogger()

	// общий конфиг приложения
	appConf := config.NewConfig(appLogger)

	// канал для передачи уведомлений
	notificationChan := make(chan dto.Notification)

	// канал для передачи статистики отправки
	statChan := make(chan dto.Stat)

	// Подготовка зависимостей сервисов
	ampqClient := ampq.New("", appLogger)
	notificationService := notify.New(notificationChan, appLogger)
	eventService := event.New(appLogger)
	templateService := template.New(appLogger)
	statisticService := stat.New(statChan, appLogger)

	contactVault := notificationDispatcher.NewContactVaultClient(&appConf, appLogger)
	serviceFacade := notificationDispatcher.NewDispatcherServiceFacade(contactVault, templateService, eventService)
	dispatcherService := notificationDispatcher.New(notificationChan, &appConf, serviceFacade, ampqClient, appLogger)

	return App{
		config: appConf,
		services: services{
			notificationService:    notificationService,
			notificationDispatcher: dispatcherService,
			eventService:           eventService,
			templateService:        templateService,
			statisticService:       statisticService,
		},
		notificationChan: notificationChan,
		statChan:         statChan,
		logger:           appLogger,
	}
}

// Run инициализация зависимостей, запуск стартовых методов сервисов, запуск воркеров, запуск веб сервера
func (a App) Run(ctx context.Context) {
	// операции корректного завершения работы
	defer a.Stop()

	// Предварительная готовность сервисов
	a.services.notificationDispatcher.Start(ctx)
	a.services.notificationService.Start(ctx)
	a.services.eventService.Start(ctx)
	a.services.templateService.Start(ctx)
	a.services.statisticService.Start(ctx)

	// запуск фоновых воркеров
	a.StartWorkers(ctx)

	// подготовка роутера для http сервера, передаем хендлерам сервисы
	// и логгер
	h := handlers.New(
		&a.config,
		a.services.eventService,
		a.services.notificationService,
		a.services.statisticService,
		a.services.templateService,
		a.logger)

	routes := router.New(h, &a.config)

	startMessage := fmt.Sprintf("Server started @ %v", a.config.GetHttpServerAddress())
	a.logger.Info(startMessage)

	// запуск веб сервера, по умолчанию с адресом localhost, порт 8080
	log.Fatal(http.ListenAndServe(a.config.GetHttpServerAddress(), routes))
}

func (a App) Stop() {
	a.services.notificationService.Stop()
	a.services.eventService.Stop()
	a.services.templateService.Stop()
	a.services.statisticService.Stop()
	a.services.notificationDispatcher.Stop()
	a.logger.Info("All services stopped")
}

// StartWorkers запуск фоновых воркеров непосредственной отправки сообщений
func (a App) StartWorkers(ctx context.Context) {
	var (
		ampqClient    interfaces.AmpqClient
		channelWorker interfaces.Worker
	)

	ampqClient = ampq.NewWithConnection(a.config.GetAmpqDSN(), a.logger)
	channelWorker = workers.NewChannelWorker(ctx, &a.config, ampqClient, a.statChan, a.logger)

	go func() {
		channelWorker.Start(ctx, a.config.GetNotificationQueue(), "", a.config.GetFailedWorksQueue())
		defer channelWorker.Stop()
	}()
}
