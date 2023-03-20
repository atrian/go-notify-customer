package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

var _ senderConfig = (*Config)(nil)

type senderConfig interface {
	GetAmpqDSN() string
	GetNotificationQueue() string
	GetFailedWorksQueue() string
	mailConfig
	twilioConfig
}

type mailConfig interface {
	GetMailSenderAddress() string
	GetMailSMTPHost() string
	GetMailLogin() string
	GetMailPassword() string
	GetMailMessageTheme() string
}

type twilioConfig interface {
	GetTwilioAccountSid() string
	GetTwilioAuthToken() string
	GetTwilioSenderPhone() string
}

type Config struct {
	address           string `env:"NC_ADDRESS"`
	grpcVaultAddress  string `env:"NC_GRPC_VAULT_ADDRESS"`
	ampqDSN           string `env:"NC_AMPQDSN"`
	notificationQueue string `env:"NC_DISPATCH_QUEUE"`
	failedWorksQueue  string `env:"NC_FAILED_QUEUE"`
}

func (config *Config) GetAmpqDSN() string {
	return config.ampqDSN
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

func (config *Config) GetMailSenderAddress() string {
	//TODO implement me
	panic("implement me")
}

func (config *Config) GetMailSMTPHost() string {
	//TODO implement me
	panic("implement me")
}

func (config *Config) GetMailLogin() string {
	//TODO implement me
	panic("implement me")
}

func (config *Config) GetMailPassword() string {
	//TODO implement me
	panic("implement me")
}

func (config *Config) GetMailMessageTheme() string {
	//TODO implement me
	panic("implement me")
}

func (config *Config) GetTwilioAccountSid() string {
	//TODO implement me
	panic("implement me")
}

func (config *Config) GetTwilioAuthToken() string {
	//TODO implement me
	panic("implement me")
}

func (config *Config) GetTwilioSenderPhone() string {
	//TODO implement me
	panic("implement me")
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
