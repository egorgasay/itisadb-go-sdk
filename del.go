package itisadb

import (
	"github.com/egorgasay/itisadb-go-sdk/api"
	"golang.org/x/net/context"
)

func (c *Client) del(ctx context.Context, key string, server int32) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	_, err := c.cl.Delete(withAuth(ctx), &api.DeleteRequest{
		Key:    key,
		Server: &server,
	})
	if err != nil {
		return convertGRPCError(err)
	}
	return nil
}

func (c *Client) Del(ctx context.Context, key string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.del(ctx, key, c.keysAndServers[key])
}
