package authservice

import (
	"github.com/231031/authservice/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.AuthServiceClient
}

func NewClient(url string) (*Client, error) {
	c, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	s := pb.NewAuthServiceClient(c)
	return &Client{conn: c, service: s}, nil
}
