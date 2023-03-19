package notificationDispatcher

import (
	"context"
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/atrian/go-notify-customer/internal/dto"
	"github.com/atrian/go-notify-customer/pkg/logger"
	pb "github.com/atrian/go-notify-customer/proto"
)

type ContactVaultTestSuite struct {
	suite.Suite
	vaultClient contactVault
	grpcServ    *grpc.Server
	config      grpcConfig
	ctx         context.Context
	cancel      context.CancelFunc
	port        int
}

func (suite *ContactVaultTestSuite) SetupSuite() {
	suite.ctx, suite.cancel = context.WithCancel(context.Background())

	listener, _ := net.Listen("tcp", ":0")
	port := listener.Addr().(*net.TCPAddr).Port

	suite.config = grpcConfigMock{port: port}

	// создаём gRPC-сервер без зарегистрированной службы
	suite.grpcServ = grpc.NewServer()

	// регистрируем сервис
	pb.RegisterVaultServer(suite.grpcServ, &contactServerMock{})

	// получаем запрос gRPC
	go func() {
		for {
			select {
			case <-suite.ctx.Done():
				_ = suite.vaultClient.Stop()
				suite.grpcServ.Stop()
			default:
				_ = suite.grpcServ.Serve(listener)
			}
		}
	}()

	// формируем клиент
	suite.vaultClient = NewContactVaultClient(suite.config, logger.NewZapLogger())
}

func (suite *ContactVaultTestSuite) TearDownSuite() {
	suite.cancel()
}

func (suite *ContactVaultTestSuite) Test_GrpcContactVault_FindByPersonUUID() {
	pUUID := uuid.New()
	res, err := suite.vaultClient.FindByPersonUUID(context.Background(), pUUID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), pUUID, res.PersonUUID)
	assert.Equal(suite.T(), []dto.Contact{
		{
			Channel:     "sms",
			Destination: "+79876543210",
		}, {
			Channel:     "email",
			Destination: "dummy@mail.com",
		},
	}, res.Contacts)
}

func (suite *ContactVaultTestSuite) Test_GrpcContactVault_Stop() {
	err := suite.vaultClient.Stop()
	assert.NoError(suite.T(), err)

	_, err = suite.vaultClient.FindByPersonUUID(context.Background(), uuid.New())
	stat := status.Code(err)
	assert.Equal(suite.T(), codes.Canceled, stat)
}

var _ grpcConfig = (*grpcConfigMock)(nil)

type grpcConfigMock struct {
	port int
}

func (g grpcConfigMock) GetGRPCAddress() string {
	return fmt.Sprintf(":%v", g.port)
}

type contactServerMock struct {
	pb.UnimplementedVaultServer
}

func (c *contactServerMock) GetContacts(ctx context.Context, in *pb.GetContactsRequest) (*pb.GetContactsResponse, error) {
	var response pb.GetContactsResponse

	personUUID, err := uuid.Parse(in.GetPersonUUID())

	phone := pb.Contact{
		PersonUuid:  personUUID.String(),
		Channel:     "sms",
		Destination: "+79876543210",
	}

	email := pb.Contact{
		PersonUuid:  personUUID.String(),
		Channel:     "email",
		Destination: "dummy@mail.com",
	}

	if err != nil {
		response.Error = fmt.Sprintf("Bad request")
		return nil, errors.New("bad request")
	}

	response.Contacts = []*pb.Contact{
		&phone,
		&email,
	}

	return &response, nil
}

// Для запуска через Go test
func TestVaultSuite(t *testing.T) {
	suite.Run(t, new(ContactVaultTestSuite))
}
