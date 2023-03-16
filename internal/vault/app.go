package vault

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/atrian/go-notify-customer/internal/interfaces"
	"github.com/atrian/go-notify-customer/pkg/logger"
	pb "github.com/atrian/go-notify-customer/proto"
)

type App struct {
	ctx    context.Context
	logger interfaces.Logger
}

func New(ctx context.Context) *App {
	l := logger.NewZapLogger()

	a := App{
		ctx:    ctx,
		logger: l,
	}

	return &a
}

func (a *App) Run() {
	// определяем порт для сервера
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		a.logger.Fatal("net.Listen error", err)
	}

	// создаём gRPC-сервер без зарегистрированной службы
	s := grpc.NewServer()

	// регистрируем сервис
	pb.RegisterVaultServer(s, &ContactServer{})

	a.logger.Info("Vault gRPC server started")

	// получаем запрос gRPC
	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}
}
