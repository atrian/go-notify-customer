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
	ctx      context.Context
	logger   interfaces.Logger
	listener net.Listener
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
	if a.listener == nil {
		a.SetDefaultListener()
	}

	// создаём gRPC-сервер без зарегистрированной службы
	s := grpc.NewServer()

	// регистрируем сервис
	pb.RegisterVaultServer(s, &ContactServer{})

	a.logger.Info("Vault gRPC server started")

	// получаем запрос gRPC
	if err := s.Serve(a.listener); err != nil {
		log.Fatal(err)
	}
}

func (a *App) SetCustomListener(listener net.Listener) *App {
	a.listener = listener
	return a
}

func (a *App) SetDefaultListener() *App {
	// определяем порт для сервера
	l, err := net.Listen("tcp", ":3200")
	if err != nil {
		a.logger.Fatal("net.Listen error", err)
	}

	a.listener = l
	return a
}

func (a *App) Stop() {
	// TODO shutdown
}
