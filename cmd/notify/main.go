package main

import (
	"context"

	"github.com/atrian/go-notify-customer/internal/notify"
)

func main() {
	ctx := context.Background()

	application := notify.New(ctx)
	application.Run()
}
