package grpcclient

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"osbourne.local/frontend/gen/profile"
)

type ProfileClient struct {
	conn   *grpc.ClientConn
	Client profile.ProfileServiceClient
}

func NewProfileClient(addr string) (*ProfileClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &ProfileClient{
		conn:   conn,
		Client: profile.NewProfileServiceClient(conn),
	}, nil
}

func (c *ProfileClient) Close() {
	_ = c.conn.Close()
}
