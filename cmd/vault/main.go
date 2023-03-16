package main

import (
	"context"

	"github.com/atrian/go-notify-customer/internal/vault"
)

func main() {
	ctx := context.Background()

	application := vault.New(ctx)
	application.Run()
}
