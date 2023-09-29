package itisadb

import (
	"context"
	"errors"
	"google.golang.org/grpc/metadata"
)

const token = "token"

var authMetadata = metadata.New(map[string]string{token: ""})

var ErrUnauthorized = errors.New("unauthorized")

func withAuth(ctx context.Context) context.Context {
	return metadata.NewOutgoingContext(ctx, authMetadata)
}
