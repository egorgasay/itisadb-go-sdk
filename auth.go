package itisadb

import (
	"context"

	"github.com/egorgasay/gost"
	api "github.com/egorgasay/itisadb-shared-proto/go"
	"google.golang.org/grpc/metadata"
)

const token = "token"

const (
	DefaultLevel Level = iota
	RestrictedLevel
	SecretLevel
)

// TODO: refactor
var authMetadata = metadata.New(map[string]string{token: ""})

func withAuth(ctx context.Context) context.Context {
	if ctx.Value(token) != nil {
		return ctx
	}

	return metadata.NewOutgoingContext(ctx, authMetadata)
}

func (c *Client) CreateUser(ctx context.Context, login, password string, opts ...CreateUserOptions) (res gost.Result[gost.Nothing]) {
	opt := CreateUserOptions{}

	if len(opts) > 0 {
		opt = opts[0]
	}

	_, err := c.cl.CreateUser(withAuth(ctx), &api.CreateUserRequest{
		User: &api.User{Login: login, Password: password, Level: uint32(opt.Level)},
	})

	if err != nil {
		err := errFromGRPCError(err)
		if err == ErrUniqueConstraint {
			return res
		}

		return res.Err(err)
	}

	return res.Ok(gost.Nothing{})
}

func (c *Client) DeleteUser(ctx context.Context, login string) (res gost.Result[bool]) {
	_, err := c.cl.DeleteUser(withAuth(ctx), &api.DeleteUserRequest{
		Login: login,
	})

	if err != nil {
		err := errFromGRPCError(err)
		if err == ErrNotFound {
			return res
		}

		return res.Err(err)
	}

	return res.Ok(true)
}

func (c *Client) ChangePassword(ctx context.Context, login, newPassword string) error {
	_, err := c.cl.ChangePassword(withAuth(ctx), &api.ChangePasswordRequest{
		Login:       login,
		NewPassword: newPassword,
	})

	if err != nil {
		return errFromGRPCError(err)
	}

	return nil
}

func (c *Client) ChangeLevel(ctx context.Context, login string, level Level) error {
	_, err := c.cl.ChangeLevel(withAuth(ctx), &api.ChangeLevelRequest{
		Login: login,
		Level: int32(level),
	})

	if err != nil {
		return errFromGRPCError(err)
	}

	return nil
}
