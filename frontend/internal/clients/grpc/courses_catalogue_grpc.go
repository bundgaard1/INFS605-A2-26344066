package grpcclient

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	coursecatalogue "osbourne.local/frontend/gen/course-catalogue"
)

type CourseCatalogueClient struct {
	conn   *grpc.ClientConn
	Client coursecatalogue.CourseCatalogueServiceClient
}

func NewCourseCatalogueClient(addr string) (*CourseCatalogueClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &CourseCatalogueClient{
		conn:   conn,
		Client: coursecatalogue.NewCourseCatalogueServiceClient(conn),
	}, nil
}

func (c *CourseCatalogueClient) Close() {
	_ = c.conn.Close()
}
