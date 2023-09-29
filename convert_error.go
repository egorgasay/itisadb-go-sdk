package itisadb

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func convertGRPCError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.NotFound:
		return ErrNotFound
	case codes.Unavailable:
		return ErrUnavailable
	case codes.ResourceExhausted:
		return ErrObjectNotFound
	case codes.AlreadyExists:
		return ErrUniqueConstraint
	case codes.Unauthenticated:
		return ErrUnauthorized
	default:
		return err
	}
}
