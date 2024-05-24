package itisadb

import (
	"fmt"

	"github.com/egorgasay/gost"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrNotFound         = gost.NewErrX(0, "not found")
	ErrUnavailable      = gost.NewErrX(0, "storage is unavailable")
	ErrUnauthorized     = gost.NewErrX(0, "unauthorized")
	ErrObjectNotFound   = gost.NewErrX(0, "object not found")
	ErrUniqueConstraint = gost.NewErrX(0, "unique constraint failed")
	ErrPermissionDenied = gost.NewErrX(0, "permission denied")
	ErrContextCanceled  = gost.NewErrX(0, "context canceled")
)

func errFromGRPCError(err error) *gost.ErrX {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return gost.NewErrX(
			0,
			"unknown error",
		).ExtendMsg(fmt.Sprint(err.Error()))
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
		return gost.NewErrX(
			0,
			"unknown error",
		).ExtendMsg(fmt.Sprint(err.Error()))
	}
}
