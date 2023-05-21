package itisadb

import (
	"github.com/egorgasay/itisadb-go-sdk/api/balancer"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Client) del(ctx context.Context, key string, server int32) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	_, err := c.cl.Delete(ctx, &balancer.BalancerDeleteRequest{
		Key:    key,
		Server: server,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return err
		}

		if st.Code() == codes.NotFound {
			return ErrNotFound
		}

		if st.Code() == codes.Unavailable {
			return ErrUnavailable
		}

		return err
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
