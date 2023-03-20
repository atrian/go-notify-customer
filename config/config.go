package config

import (
	"flag"
	"github.com/caarlos0/env/v6"

	"github.com/atrian/go-notify-customer/internal/interfaces"
)

var (
	_ senderConfig   = (*Config)(nil)
	_ grpcConfig     = (*Config)(nil)
	_ webConfig      = (*Config)(nil)
	_ securityConfig = (*Config)(nil)
)

type webConfig interface {
	GetHttpServerAddress() string
}

type securityConfig interface {
	GetTrustedSubnetAddress() string
}

type grpcConfig interface {
	GetGRPCAddress() string
}

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
	httpAddress             string `env:"NC_HTTP_ADDRESS"`
	httpTrustedSubnet       string `env:"NC_TRUSTED_SUBNET"`
	grpcVaultAddress        string `env:"NC_GRPC_VAULT_ADDRESS"`
	ampqDSN                 string `env:"NC_AMPQDSN"`
	notificationQueue       string `env:"NC_DISPATCH_QUEUE"`
	failedWorksQueue        string `env:"NC_FAILED_QUEUE"`
	mailSenderAddress       string `env:"NC_MAIL_SENDER_ADDRESS"`
	mailSMTPHost            string `env:"NC_MAIL_SMTP_HOST"`
	mailLogin               string `env:"NC_MAIL_LOGIN"`
	mailPassword            string `env:"NC_MAIL_PASSWORD"`
	mailDefaultMessageTheme string `env:"NC_MAIL_DEFAULT_MESSAGE_THEME"`
	twilioAccountSid        string `env:"NC_TWILIO_ACCOUNT_ID"`
	twilioAuthToken         string `env:"NC_TWILIO_ACCOUNT_ID"`
	twilioSenderPhone       string `env:"NC_TWILIO_SENDER_PHONE"`
	log                     interfaces.Logger
}

func (config *Config) GetTrustedSubnetAddress() string {
	return config.httpTrustedSubnet
}

func NewConfig(logger interfaces.Logger) Config {
	conf := Config{
		log: logger,
	}

	conf.loadDefaults()
	conf.loadEnv()
	conf.loadFlags()

	return conf
}

func (config *Config) GetAmpqDSN() string {
	return config.ampqDSN
}

func (config *Config) GetHttpServerAddress() string {
	return config.httpAddress
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
	return config.mailSenderAddress
}

func (config *Config) GetMailSMTPHost() string {
	return config.mailSMTPHost
}

func (config *Config) GetMailLogin() string {
	return config.mailLogin
}

func (config *Config) GetMailPassword() string {
	return config.mailPassword
}

func (config *Config) GetMailMessageTheme() string {
	return config.mailDefaultMessageTheme
}

func (config *Config) GetTwilioAccountSid() string {
	return config.twilioAccountSid
}

func (config *Config) GetTwilioAuthToken() string {
	return config.twilioAuthToken
}

func (config *Config) GetTwilioSenderPhone() string {
	return config.twilioSenderPhone
}

func (config *Config) loadDefaults() {
	config.notificationQueue = "planned_notifications"
	config.failedWorksQueue = "failed_notifications"
}

// loadFlags загрузка в конфигурацию флагов запуска приложения
func (config *Config) loadFlags() {
	httpAddress := flag.String("a", "127.0.0.1:8080", "Address and port used for GO-notify-customer app webserver.")
	grpcVaultAddress := flag.String("g", "127.0.0.1:50051", "Address and port used GRPC vault connection.")
	ampqDsn := flag.String("ad", "amqp://guest:guest@localhost:5672/", "DSN for AMPQ server.")

	flag.Parse()

	config.httpAddress = *httpAddress
	config.grpcVaultAddress = *grpcVaultAddress
	config.ampqDSN = *ampqDsn
}

// loadEnv загрузка в конфигурацию данных из переменных окружения
func (config *Config) loadEnv() {
	err := env.Parse(config)
	if err != nil {
		config.log.Error("Config parse error", err)
	}
}
