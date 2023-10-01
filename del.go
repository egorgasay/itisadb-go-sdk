package itisadb

import (
	"github.com/egorgasay/itisadb-go-sdk/api"
	"golang.org/x/net/context"
)

func (c *Client) del(ctx context.Context, key string, server int32) (res Result[bool]) {
	_, err := c.cl.Delete(withAuth(ctx), &api.DeleteRequest{
		Key:    key,
		Server: &server,
	})

	if err != nil {
		res.err = convertGRPCError(err)
	} else {
		res.value = true
	}

	return res
}

func (c *Client) Del(ctx context.Context, key string) Result[bool] {
	c.mu.RLock()
	defer c.mu.RUnlock()

	s := c.keysAndServers[key]
	delete(c.keysAndServers, key)

	return c.del(ctx, key, s)
}
