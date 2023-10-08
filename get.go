package itisadb

import (
	"context"
	"fmt"
	"github.com/egorgasay/itisadb-go-sdk/api"
)

const (
	_ = -iota
	getFromDisk
)

func (c *Client) get(ctx context.Context, key string, opts GetOptions) (res Result[string]) {
	r, err := c.cl.Get(withAuth(ctx), &api.GetRequest{
		Key: key,
		Options: &api.GetRequest_Options{
			Server: opts.Server,
		},
	})

	if err != nil {
		res.err = convertGRPCError(err)
	} else {
		res.value = r.Value
	}

	return res
}

// GetOne gets the value by the key from gRPCis.
func (c *Client) GetOne(ctx context.Context, key string, opts ...GetOptions) (res Result[string]) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	opt := GetOptions{}

	if len(opts) > 0 {
		opt = opts[0]
	}

	if opt.Server != nil {
		if s, ok := c.keysAndServers[key]; ok {
			opt.Server = &s
		}
	}

	return c.get(ctx, key, opt)
}

// GetMany gets a lot of values from gRPCis.
func (c *Client) GetMany(ctx context.Context, keys []string, opts ...GetOptions) (map[string]string, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	var keyValue = make(map[string]string, 10)
	var err error

	opt := GetOptions{}

	if len(opts) > 0 {
		opt = opts[0]
	}

	for _, key := range keys {
		keyValue[key], err = c.GetOne(ctx, key, opt).ValueAndErr()
		if err != nil {
			return nil, fmt.Errorf("get %s: %w", key, err)
		}
	}
	return keyValue, nil
}

// GetManyOpts gets a lot of values from gRPCis with opts.
func (c *Client) GetManyOpts(ctx context.Context, keys []Key) (map[string]string, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	var keyValue = make(map[string]string, 10)
	var err error

	for _, key := range keys {
		keyValue[key.Key], err = c.get(ctx, key.Key, key.Options).ValueAndErr()
		if err != nil {
			return nil, err
		}
	}

	return keyValue, nil
}
