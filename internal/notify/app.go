package notify

import (
	"context"
	"fmt"
	"github.com/atrian/go-notify-customer/internal/workers"
	"log"
	"net/http"

	"github.com/atrian/go-notify-customer/config"
	"github.com/atrian/go-notify-customer/internal/dto"
	"github.com/atrian/go-notify-customer/internal/interfaces"
	"github.com/atrian/go-notify-customer/internal/services/event"
	"github.com/atrian/go-notify-customer/internal/services/notificationDispatcher"
	"github.com/atrian/go-notify-customer/internal/services/notify"
	"github.com/atrian/go-notify-customer/internal/services/stat"
	"github.com/atrian/go-notify-customer/internal/services/template"
	"github.com/atrian/go-notify-customer/pkg/ampq"
	"github.com/atrian/go-notify-customer/pkg/logger"
)

type App struct {
	ctx              context.Context
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

func New(ctx context.Context) App {
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
	notificationService := notify.New(notificationChan)
	eventService := event.New()
	templateService := template.New()
	statisticService := stat.New(ctx, statChan)

	contactVault := notificationDispatcher.NewContactVaultClient(&appConf, appLogger)
	serviceFacade := notificationDispatcher.NewDispatcherServiceFacade(contactVault, templateService, eventService)
	dispatcherService := notificationDispatcher.New(ctx, notificationChan, &appConf, serviceFacade, ampqClient, appLogger)

	return App{
		ctx:    ctx,
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
func (a App) Run() {
	// операции корректного завершения работы
	defer a.Stop()

	// Предварительная готовность сервисов
	a.services.notificationService.Start()
	a.services.eventService.Start()
	a.services.templateService.Start()
	a.services.statisticService.Start()

	// запуск фоновых воркеров
	a.StartWorkers()

	// подготовка роутера для http сервера, передаем хендлерам сервисы
	// и логгер
	/*routes := router.New(handlers.New(
	a.services.eventService,
	a.services.templateService,
	a.logger))*/

	startMessage := fmt.Sprintf("Server started @ %v", a.config.GetHttpServerAddress())
	a.logger.Info(startMessage)

	// запуск веб сервера, по умолчанию с адресом localhost, порт 8080
	// TODO прокинуть роутер в http сервер
	log.Fatal(http.ListenAndServe(a.config.GetHttpServerAddress(), nil))
}

func (a App) Stop() {
	// корректно завершаем сервисы
	//a.services.[nService].Stop()
}

func (a App) StartWorkers() {
	var (
		ampqClient    interfaces.AmpqClient
		channelWorker interfaces.Worker
	)

	ampqClient = ampq.NewWithConnection(a.config.GetAmpqDSN(), logger.NewZapLogger())
	channelWorker = workers.NewChannelWorker(a.ctx, &a.config, ampqClient, a.statChan, a.logger)

	go func() {
		channelWorker.Start(a.config.GetNotificationQueue(), "", a.config.GetFailedWorksQueue())
		defer channelWorker.Stop()
	}()
}
