package itisadb

import (
	"context"

	"github.com/egorgasay/gost"
	api "github.com/egorgasay/itisadb-shared-proto/go"
)

type RAM struct {
	Total     uint64
	Available uint64
}

type internal struct{}

var Internal internal

func (i *internal) GetRAM(ctx context.Context, c *Client) (res gost.Result[RAM]) {
	r, err := c.cl.GetRam(withAuth(ctx), &api.GetRamRequest{})
	if err != nil {
		return res.Err(errFromGRPCError(err))
	}

	ram := r.GetRam()

	return res.Ok(RAM{Total: ram.Total, Available: ram.Available})
}

func (i *internal) GetLastUserChangeID(ctx context.Context, c *Client) (res gost.Result[uint64]) {
	r, err := c.cl.GetLastUserChangeID(withAuth(ctx), &api.GetLastUserChangeIDRequest{})
	if err != nil {
		return res.Err(errFromGRPCError(err))
	}

	return res.Ok(r.LastChangeID)
}

func (i *internal) Sync(ctx context.Context, c *Client, users []Internal_User) (res gost.ResultN) {
	_, err := c.cl.Sync(withAuth(ctx), &api.SyncData{
		Users: apiUsersFromInternalUsers(users),
	})

	if err != nil {
		return res.Err(errFromGRPCError(err))
	}

	return res.Ok()
}

func apiUsersFromInternalUsers(users []Internal_User) []*api.User {
	var out []*api.User
	for _, u := range users {
		out = append(out, &api.User{
			Id:       uint64(u.ID),
			Login:    u.Login,
			Password: u.Password,
			Level:    uint32(u.Level),
			Active:   u.Active,
		})
	}
	return out
}
