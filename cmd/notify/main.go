package notify

import "github.com/atrian/go-notify-customer/internal/notify"

func main() {
	application := notify.New()
	application.Run()
}
