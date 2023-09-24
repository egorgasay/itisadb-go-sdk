package itisadb

import (
	"context"
	"google.golang.org/grpc/metadata"
)

const token = "token"

func (c *Client) withAuth(ctx context.Context) context.Context {
	md := metadata.New(map[string]string{token: c.token})
	return metadata.NewOutgoingContext(ctx, md)
}
