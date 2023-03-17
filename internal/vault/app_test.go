package vault

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/atrian/go-notify-customer/proto"
)

type VaultTestSuite struct {
	suite.Suite
	contacts []*pb.Contact
	app      *App
	port     int
}

func (suite *VaultTestSuite) SetupSuite() {
	listener, err := net.Listen("tcp", ":0")
	assert.NoError(suite.T(), err)

	suite.port = listener.Addr().(*net.TCPAddr).Port

	ctx := context.Background()

	suite.app = New(ctx)
	go suite.app.SetCustomListener(listener).Run()
}

func (suite *VaultTestSuite) TearDownSuite() {
	suite.app.Stop()
}

func (suite *VaultTestSuite) Test_RunWithCustomListener() {
	ctx := context.Background()
	personUUID := uuid.New()

	// устанавливаем соединение с сервером
	conn, err := grpc.Dial(fmt.Sprintf(":%v", suite.port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()

	assert.NoError(suite.T(), err)

	client := pb.NewVaultClient(conn)
	resp, err := client.GetContacts(ctx, &pb.GetContactsRequest{PersonUUID: personUUID.String()})

	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 2, len(resp.Contacts))
	assert.Equal(suite.T(), personUUID.String(), resp.Contacts[0].GetPersonUuid())

	resp, err = client.GetContacts(ctx, &pb.GetContactsRequest{PersonUUID: "RandomData"})
	assert.Error(suite.T(), err)
}

// Для запуска через Go test
func TestVaultSuite(t *testing.T) {
	suite.Run(t, new(VaultTestSuite))
}
