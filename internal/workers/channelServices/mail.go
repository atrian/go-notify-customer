package channelServices

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"strings"
	
	"github.com/atrian/go-notify-customer/internal/interfaces"
)

type Mail struct {
	conf   configMail
	ctx    context.Context
	logger interfaces.Logger
}

type configMail interface {
	GetMailSenderAddress() string
	GetMailSMTPHost() string
	IsMailTLSRequired() bool
	GetMailLogin() string
	GetMailPassword() string
	GetMailMessageTheme() string
}

func NewMail(ctx context.Context, conf configMail, logger interfaces.Logger) *Mail {
	return &Mail{
		conf:   conf,
		ctx:    ctx,
		logger: logger,
	}
}

func (s *Mail) SendMessage(message string, destination string) error {
	headers := make(map[string]string)
	headers["From"] = s.conf.GetMailSenderAddress()
	headers["To"] = destination
	headers["Subject"] = s.conf.GetMailMessageTheme()

	mail := strings.Builder{}
	mail.Grow(len(headers) * 2)

	for k, v := range headers {
		mail.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	mail.WriteString("\r\n")
	mail.WriteString(message)

	// Connect to the SMTP Server
	servername := s.conf.GetMailSMTPHost()

	host, _, _ := net.SplitHostPort(servername)
	auth := smtp.PlainAuth("", s.conf.GetMailLogin(), s.conf.GetMailPassword(), host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		log.Panic(err)
	}

	c, err := smtp.NewClient(conn, host)
	defer func(c *smtp.Client) {
		cErr := c.Close()
		if cErr != nil {
			s.logger.Error("smtp client close err", err)
		}
	}(c)
	if err != nil {
		return err
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		return err
	}

	// To && From
	if err = c.Mail(s.conf.GetMailSenderAddress()); err != nil {
		return err
	}

	if err = c.Rcpt(destination); err != nil {
		return err
	}

	// Data
	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		s.logger.Error("Writer close err", err)
	}

	return nil
}
