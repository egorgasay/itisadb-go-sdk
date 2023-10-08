package itisadb

import (
	"context"
	"google.golang.org/grpc/metadata"
)

const token = "token"

var authMetadata = metadata.New(map[string]string{token: ""})

func withAuth(ctx context.Context) context.Context {
	return metadata.NewOutgoingContext(ctx, authMetadata)
}
