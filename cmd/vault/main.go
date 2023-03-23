package main

import (
	"context"

	"github.com/atrian/go-notify-customer/config"
	"github.com/atrian/go-notify-customer/internal/vault"
	"github.com/atrian/go-notify-customer/pkg/logger"
)

func main() {
	ctx := context.Background()
	conf := config.NewConfig(logger.NewZapLogger())

	application := vault.New(&conf)
	application.Run(ctx)
}
