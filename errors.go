package itisadb

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrNotFound         = errors.New("not found")
	ErrUnavailable      = errors.New("storage is unavailable")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrObjectNotFound   = errors.New("object not found")
	ErrUniqueConstraint = errors.New("unique constraint failed")
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