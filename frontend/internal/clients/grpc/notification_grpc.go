package grpcclient

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"osbourne.local/frontend/gen/notification"
)

type NotificationClient struct {
	conn   *grpc.ClientConn
	Client notification.NotificationServiceClient
}

func NewNotificationClient(addr string) (*NotificationClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &NotificationClient{
		conn:   conn,
		Client: notification.NewNotificationServiceClient(conn),
	}, nil
}

func (c *NotificationClient) Close() {
	_ = c.conn.Close()
}
