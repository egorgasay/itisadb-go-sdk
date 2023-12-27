package itisadb

import (
	"fmt"
	"github.com/egorgasay/gost"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrNotFound         = gost.NewError(0, 0, "not found")
	ErrUnavailable      = gost.NewError(0, 0, "storage is unavailable")
	ErrUnauthorized     = gost.NewError(0, 0, "unauthorized")
	ErrObjectNotFound   = gost.NewError(0, 0, "object not found")
	ErrUniqueConstraint = gost.NewError(0, 0, "unique constraint failed")
	ErrPermissionDenied = gost.NewError(0, 0, "permission denied")
)

func errFromGRPCError(err error) *gost.Error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return gost.NewError(
			0, 0,
			fmt.Sprintf("unknown error: %s\n", err.Error()),
		)
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
	case codes.PermissionDenied:
		return ErrPermissionDenied
	default:
		return gost.NewError(
			0, 0,
			fmt.Sprintf("unknown error: %s\n", st.Message()),
		)
	}
}
