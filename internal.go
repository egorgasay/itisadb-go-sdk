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
