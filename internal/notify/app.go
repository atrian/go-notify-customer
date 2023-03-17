package notify

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/atrian/go-notify-customer/config"
	"github.com/atrian/go-notify-customer/internal/dto"
	"github.com/atrian/go-notify-customer/internal/interfaces"
	"github.com/atrian/go-notify-customer/internal/services/event"
	"github.com/atrian/go-notify-customer/internal/services/notify"
	"github.com/atrian/go-notify-customer/internal/services/stat"
	"github.com/atrian/go-notify-customer/internal/services/template"
	"github.com/atrian/go-notify-customer/pkg/logger"
)

type App struct {
	services services
	config   config.Config
	logger   interfaces.Logger
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
	// общий конфиг приложения
	appConf := config.NewConfig()
	// логгер приложения
	appLogger := logger.NewZapLogger()

	// канал для передачи уведомлений
	notificationChan := make(chan dto.Notification)

	// канал для передачи статистики отправки
	statChan := make(chan dto.Stat)

	return App{
		config: appConf,
		services: services{
			notificationService:    notify.New(notificationChan),
			notificationDispatcher: nil,
			eventService:           event.New(),
			templateService:        template.New(),
			statisticService:       stat.New(ctx, statChan),
		},
		logger: appLogger,
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

	startMessage := fmt.Sprintf("Server started @ %v", a.config.Address)
	a.logger.Info(startMessage)

	// запуск веб сервера, по умолчанию с адресом localhost, порт 8080
	// TODO прокинуть роутер в http сервер
	log.Fatal(http.ListenAndServe(a.config.Address, nil))
}

func (a App) Stop() {
	// корректно завершаем сервисы
	//a.services.[nService].Stop()
}

func (a App) StartWorkers() {
	var channelWorker interfaces.Worker
	_ = channelWorker
	// конфигурация и запуск воркеров отправки конкретных каналов
}
