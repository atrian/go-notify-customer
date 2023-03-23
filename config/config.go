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
	IsMailTLSRequired() bool
}

type twilioConfig interface {
	GetTwilioAccountSid() string
	GetTwilioAuthToken() string
	GetTwilioSenderPhone() string
}

type Config struct {
	data Params
	log  interfaces.Logger
}

type Params struct {
	HttpAddress             string `env:"NC_HTTP_ADDRESS"`
	HttpTrustedSubnet       string `env:"NC_TRUSTED_SUBNET"`
	GrpcVaultAddress        string `env:"NC_GRPC_VAULT_ADDRESS"`
	AmpqDSN                 string `env:"NC_AMPQDSN"`
	NotificationQueue       string `env:"NC_DISPATCH_QUEUE" envDefault:"planned_notifications"`
	FailedWorksQueue        string `env:"NC_FAILED_QUEUE" envDefault:"failed_notifications"`
	MailSenderAddress       string `env:"NC_MAIL_SENDER_ADDRESS"`
	MailSMTPHost            string `env:"NC_MAIL_SMTP_HOST"`
	MailTLSRequired         bool   `env:"NC_MAIL_TLS_REQUIRED"`
	MailLogin               string `env:"NC_MAIL_LOGIN"`
	MailPassword            string `env:"NC_MAIL_PASSWORD"`
	MailDefaultMessageTheme string `env:"NC_MAIL_DEFAULT_MESSAGE_THEME"`
	TwilioAccountSid        string `env:"NC_TWILIO_ACCOUNT_ID"`
	TwilioAuthToken         string `env:"NC_TWILIO_ACCOUNT_ID"`
	TwilioSenderPhone       string `env:"NC_TWILIO_SENDER_PHONE"`
}

func (config *Config) GetDefaultResponseContentType() string {
	return "application/json"
}

func (config *Config) GetTrustedSubnetAddress() string {
	return config.data.HttpTrustedSubnet
}

func NewConfig(logger interfaces.Logger) Config {
	conf := Config{
		log: logger,
	}

	conf.loadEnv()
	conf.loadFlags()

	return conf
}

func (config *Config) GetAmpqDSN() string {
	return config.data.AmpqDSN
}

func (config *Config) GetHttpServerAddress() string {
	return config.data.HttpAddress
}

func (config *Config) GetGRPCAddress() string {
	return config.data.GrpcVaultAddress
}

func (config *Config) GetNotificationQueue() string {
	return config.data.NotificationQueue
}

func (config *Config) GetFailedWorksQueue() string {
	return config.data.FailedWorksQueue
}

func (config *Config) GetMailSenderAddress() string {
	return config.data.MailSenderAddress
}

func (config *Config) GetMailSMTPHost() string {
	return config.data.MailSMTPHost
}

func (config *Config) IsMailTLSRequired() bool {
	return config.data.MailTLSRequired
}

func (config *Config) GetMailLogin() string {
	return config.data.MailLogin
}

func (config *Config) GetMailPassword() string {
	return config.data.MailPassword
}

func (config *Config) GetMailMessageTheme() string {
	return config.data.MailDefaultMessageTheme
}

func (config *Config) GetTwilioAccountSid() string {
	return config.data.TwilioAccountSid
}

func (config *Config) GetTwilioAuthToken() string {
	return config.data.TwilioAuthToken
}

func (config *Config) GetTwilioSenderPhone() string {
	return config.data.TwilioSenderPhone
}

// loadFlags загрузка в конфигурацию флагов запуска приложения
func (config *Config) loadFlags() {
	httpAddress := flag.String("a", "127.0.0.1:8080", "Address and port used for GO-notify-customer app webserver.")
	grpcVaultAddress := flag.String("g", "127.0.0.1:50051", "Address and port used GRPC vault connection.")
	ampqDsn := flag.String("ad", "amqp://guest:guest@localhost:5672/", "DSN for AMPQ server.")

	flag.Parse()

	config.data.HttpAddress = *httpAddress
	config.data.GrpcVaultAddress = *grpcVaultAddress
	config.data.AmpqDSN = *ampqDsn
	config.log.Debug("Flag params processed")
}

// loadEnv загрузка в конфигурацию данных из переменных окружения
func (config *Config) loadEnv() {
	params := Params{}
	err := env.Parse(&params)

	if err != nil {
		config.log.Error("Config parse error", err)
	}

	config.data = params
	config.log.Debug("Env params processed")
}
