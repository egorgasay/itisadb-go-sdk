package itisadb

import (
	"context"
	"errors"
	"fmt"
	"github.com/egorgasay/itisadb-go-sdk/api"
)

const (
	_ = -iota
	getFromDisk
)

var ErrNotFound = errors.New("not found")

func (c *Client) GetFrom(ctx context.Context, key string, server int32) (res Result[string]) {
	r, err := c.cl.Get(withAuth(ctx), &api.GetRequest{
		Key:    key,
		Server: &server,
	})

	if err != nil {
		res.err = convertGRPCError(err)
	} else {
		res.value = r.Value
	}

	return res
}

// GetOne gets the value by the key from gRPCis.
func (c *Client) GetOne(ctx context.Context, key string) (res Result[string]) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.GetFrom(ctx, key, c.keysAndServers[key])
}

// GetMany gets a lot of values from gRPCis.
func (c *Client) GetMany(ctx context.Context, keys []string) (map[string]string, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	var keyValue = make(map[string]string, 10)
	var err error

	for _, key := range keys {
		keyValue[key], err = c.GetOne(ctx, key).ValueAndErr()
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
		keyValue[key.Key], err = c.GetFrom(ctx, key.Key, key.Opts.Server).ValueAndErr()
		if err != nil {
			return nil, err
		}
	}

	return keyValue, nil
}
