package itisadb

import (
	"github.com/egorgasay/gost"
	api "github.com/egorgasay/itisadb-shared-proto/go"
	"golang.org/x/net/context"
)

func (c *Client) del(ctx context.Context, key string, opts DeleteOptions) (res gost.ResultN) {
	_, err := c.cl.Delete(withAuth(ctx), &api.DeleteRequest{
		Key:     key,
		Options: &api.DeleteRequest_Options{Server: opts.Server},
	})

	if err != nil {
		return res.Err(errFromGRPCError(err))
	}

	return res.Ok()
}

func (c *Client) DelOne(ctx context.Context, key string, opts ...DeleteOptions) gost.ResultN {
	opt := DeleteOptions{}

	if len(opts) > 0 {
		opt = opts[0]
	}

	return c.del(ctx, key, opt)
}
