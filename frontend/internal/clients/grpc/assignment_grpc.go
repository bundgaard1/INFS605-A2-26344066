package grpcclient

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	assignment "osbourne.local/frontend/gen/assignment"
)

type AssignmentClient struct {
	conn   *grpc.ClientConn
	Client assignment.AssignmentServiceClient
}

func NewAssignmentClient(addr string) (*AssignmentClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &AssignmentClient{
		conn:   conn,
		Client: assignment.NewAssignmentServiceClient(conn),
	}, nil
}

func (c *AssignmentClient) Close() {
	_ = c.conn.Close()
}
