package itisadb

import (
	"github.com/egorgasay/gost"
	"github.com/egorgasay/itisadb-go-sdk/api"
	"golang.org/x/net/context"
)

func (c *Client) del(ctx context.Context, key string, opts DeleteOptions) (res gost.Result[gost.Nothing]) {
	_, err := c.cl.Delete(withAuth(ctx), &api.DeleteRequest{
		Key:     key,
		Options: &api.DeleteRequest_Options{Server: opts.Server},
	})

	if err != nil {
		return res.Err(errFromGRPCError(err))
	}

	return res.Ok(gost.Nothing{})
}

func (c *Client) DelOne(ctx context.Context, key string, opts ...DeleteOptions) gost.Result[gost.Nothing] {
	c.mu.RLock()
	defer c.mu.RUnlock()

	s, ok := c.keysAndServers[key]
	delete(c.keysAndServers, key)

	opt := DeleteOptions{}

	if len(opts) > 0 {
		opt = opts[0]
	}

	if ok {
		opt.Server = &s
	}

	return c.del(ctx, key, opt)
}
