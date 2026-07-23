package grpcclient

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"osbourne.local/frontend/gen/student"
)

type StudentClient struct {
	conn   *grpc.ClientConn
	Client student.StudentServiceClient
}

func NewStudentClient(addr string) (*StudentClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &StudentClient{
		conn:   conn,
		Client: student.NewStudentServiceClient(conn),
	}, nil
}

func (c *StudentClient) Close() {
	_ = c.conn.Close()
}
