package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

type Config struct {
	address           string `env:"NC_ADDRESS"`
	grpcVaultAddress  string `env:"NC_GRPC_VAULT_ADDRESS"`
	ampqDSN           string `env:"NC_AMPQDSN"`
	notificationQueue string `env:"NC_DISPATCH_QUEUE"`
	failedWorksQueue  string `env:"NC_FAILED_QUEUE"`
}

func (config *Config) GetWebServerAddress() string {
	return config.address
}

func (config *Config) GetGRPCAddress() string {
	return config.grpcVaultAddress
}

func (config *Config) GetNotificationQueue() string {
	return config.notificationQueue
}

func (config *Config) GetFailedWorksQueue() string {
	return config.failedWorksQueue
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
	ampqDsn := flag.String("ad", "amqp://guest:guest@localhost:5672/", "DSN for AMPQ server")

	flag.Parse()

	config.address = *address
	config.ampqDSN = *ampqDsn
}

// loadServerFlags загрузка в конфигурацию данных из переменных окружения
func (config *Config) loadServerEnvConfiguration() {
	err := env.Parse(&config.notificationQueue)
	if err != nil {
		log.Fatal("TODO - доработать конфиги")
	}
}
