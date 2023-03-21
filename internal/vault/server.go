package vault

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/atrian/go-notify-customer/internal/interfaces"
	pb "github.com/atrian/go-notify-customer/proto"
)

var BadRequest = errors.New("bad request")

type ContactServer struct {
	pb.UnimplementedVaultServer
	logger interfaces.Logger
}

func NewContactServer(logger interfaces.Logger) *ContactServer {
	s := ContactServer{
		logger: logger,
	}

	return &s
}

func (s ContactServer) GetContacts(ctx context.Context, in *pb.GetContactsRequest) (*pb.GetContactsResponse, error) {
	var response pb.GetContactsResponse
	s.logger.Debug("Contact request for UUID: ", in.GetPersonUUID())

	personUUID, err := uuid.Parse(in.GetPersonUUID())
	if err != nil {
		s.logger.Error("GetContacts uuid.Parse err", err)
	}

	phone := pb.Contact{
		PersonUuid:  personUUID.String(),
		Channel:     "sms",
		Destination: "+79876543210",
	}

	email := pb.Contact{
		PersonUuid:  personUUID.String(),
		Channel:     "mail",
		Destination: "dummy@mail.com",
	}

	if err != nil {
		response.Error = fmt.Sprintf("Bad request")
		return nil, BadRequest
	}

	response.Contacts = []*pb.Contact{
		&phone,
		&email,
	}

	return &response, nil
}
