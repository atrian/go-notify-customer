package notificationDispatcher

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/atrian/go-notify-customer/internal/dto"
	"github.com/atrian/go-notify-customer/pkg/logger"
	pb "github.com/atrian/go-notify-customer/proto"
)

var (
	_                contactVault = (*GrpcContactVault)(nil)
	ErrNilConnection              = errors.New("nil grpc connection")
)

// grpcConfig требования к конфигу для GrpcContactVault
type grpcConfig interface {
	GetGRPCAddress() string
}

// GrpcContactVault клиент для полуения контактов из хранилища по grpc
type GrpcContactVault struct {
	config grpcConfig
	logger logger.Logger
	conn   *grpc.ClientConn
}

func NewContactVaultClient(config grpcConfig, logger logger.Logger) *GrpcContactVault {
	cv := GrpcContactVault{
		config: config,
		logger: logger,
	}

	// Устанавливаем соединение с GRPC сервером
	conn, err := grpc.Dial(config.GetGRPCAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		cv.logger.Error("NewContactVaultClient grpc.Dial err", err)
	}
	cv.conn = conn

	return &cv
}

func (g GrpcContactVault) FindByPersonUUID(ctx context.Context, personUUID uuid.UUID) (dto.PersonContacts, error) {
	// Запрашиваем контактные данные у внешнего хранилища
	client := pb.NewVaultClient(g.conn)
	resp, err := client.GetContacts(ctx, &pb.GetContactsRequest{PersonUUID: personUUID.String()})

	if err != nil {
		return dto.PersonContacts{}, err
	}

	// Формируем слайс dto.Contact с контактами и отдаем ответ в обертке dto.PersonContacts
	var contacts []dto.Contact
	for _, contact := range resp.Contacts {
		contacts = append(contacts, dto.Contact{
			Channel:     contact.GetChannel(),
			Destination: contact.GetDestination(),
		})
	}

	return dto.PersonContacts{
		PersonUUID: personUUID,
		Contacts:   contacts,
	}, nil
}

func (g GrpcContactVault) Stop() error {
	if g.conn == nil {
		return ErrNilConnection
	}

	return g.conn.Close()
}
