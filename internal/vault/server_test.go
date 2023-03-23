package vault

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	"github.com/atrian/go-notify-customer/pkg/logger"
	pb "github.com/atrian/go-notify-customer/proto"
)

const bufSize = 1024 * 1024

type ServerTestSuite struct {
	suite.Suite
	contacts []*pb.Contact
	listener *bufconn.Listener
}

func (suite *ServerTestSuite) SetupSuite() {
	suite.listener = bufconn.Listen(bufSize)

	s := grpc.NewServer()
	pb.RegisterVaultServer(s, NewContactServer(logger.NewZapLogger()))

	go func() {
		if err := s.Serve(suite.listener); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func (suite *ServerTestSuite) Test_GetContacts() {
	ctx := context.Background()
	personUUID := uuid.New()

	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(suite.buffDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()

	assert.NoError(suite.T(), err)

	client := pb.NewVaultClient(conn)
	resp, err := client.GetContacts(ctx, &pb.GetContactsRequest{PersonUUID: personUUID.String()})

	// Нет ошибок, получили 2 тестовые записи
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 2, len(resp.Contacts))
	assert.Equal(suite.T(), personUUID.String(), resp.Contacts[0].GetPersonUuid())

	// На запрос рандомных данных получили ошибку
	resp, err = client.GetContacts(ctx, &pb.GetContactsRequest{PersonUUID: "RandomData"})
	assert.Error(suite.T(), err)
}

func (suite *ServerTestSuite) buffDialer(context.Context, string) (net.Conn, error) {
	return suite.listener.Dial()
}

// Для запуска через Go test
func TestStatServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}
