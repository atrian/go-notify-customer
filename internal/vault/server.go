package vault

import (
	"context"
	"errors"
	"fmt"

	pb "github.com/atrian/go-notify-customer/proto"
	"github.com/google/uuid"
)

var BadRequest = errors.New("bad request")

type ContactServer struct {
	pb.UnimplementedVaultServer
}

func (s ContactServer) GetContacts(ctx context.Context, in *pb.GetContactsRequest) (*pb.GetContactsResponse, error) {
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
		return nil, BadRequest
	}

	response.Contacts = []*pb.Contact{
		&phone,
		&email,
	}

	return &response, nil
}
