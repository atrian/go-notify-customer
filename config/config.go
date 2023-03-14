package config

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Address           string `env:"NC_ADDRESS"`
	AMPQDSN           string `env:"NC_AMPQDSN"`
	DBDSN             string `env:"NC_DBDSN"`
	NotificationQueue QueueConfig
	FailedWorksQueue  QueueConfig
}

type QueueConfig struct {
	TaskFlushInterval time.Duration `env:"NC_TASK_FLUSH_INTERVAL"` // периодичность скидывания записей из буфера в хранение
	ListenQueue       string        `env:"NC_LISTEN_QUEUE"`        // входящая очередь сообщений
	DispatchQueue     string        `env:"NC_DISPATCH_QUEUE"`      // исходящая очередь сообщений
}

func NewConfig() Config {
	conf := Config{}

	conf.loadServerEnvConfiguration()
	conf.loadServerFlags()

	return conf
}

// loadServerFlags загрузка в конфигурацию флагов запуска приложения
func (config *Config) loadServerFlags() {
	address := flag.String("a", "127.0.0.1:8080", "Address and port used for GO notify customer app.")
	dsn := flag.String("d", "postgresql://user:password@localhost:5432/postgres?sslmode=disable", "DSN for SQL server")
	ampqDsn := flag.String("ad", "amqp://guest:guest@localhost:5672/", "DSN for AMPQ server")

	flag.Parse()

	config.Address = *address
	config.DBDSN = *dsn
	config.AMPQDSN = *ampqDsn
}

// loadServerFlags загрузка в конфигурацию данных из переменных окружения
func (config *Config) loadServerEnvConfiguration() {
	err := env.Parse(&config.NotificationQueue)
	if err != nil {
		log.Fatal("TODO - доработать конфиги")
	}
}
