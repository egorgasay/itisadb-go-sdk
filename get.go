package itisadb

import (
	"context"
	"errors"
	"github.com/egorgasay/itisadb-go-sdk/api"
)

const (
	_ = -iota
	getFromDisk
)

var ErrNotFound = errors.New("not found")

func (c *Client) get(ctx context.Context, key string, server int32) (string, error) {
	res, err := c.cl.Get(withAuth(ctx), &api.GetRequest{
		Key:    key,
		Server: &server,
	})

	if err != nil {
		return "", convertGRPCError(err)
	}

	return res.Value, nil
}

// Get gets the value by the key from gRPCis.
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.get(ctx, key, c.keysAndServers[key])
}

// GetFrom gets the value by key from the specified server.
func (c *Client) GetFrom(ctx context.Context, key string, server int32) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	return c.get(ctx, key, server)
}

// GetFromDisk gets the value by key from the physical database.
func (c *Client) GetFromDisk(ctx context.Context, key string) (string, error) {
	return c.get(ctx, key, getFromDisk)
}

// GetMany gets a lot of values from gRPCis.
func (c *Client) GetMany(ctx context.Context, keys []string) (map[string]string, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	var keyValue = make(map[string]string, 10)
	var err error

	for _, key := range keys {
		keyValue[key], err = c.get(ctx, key, 0)
		if err != nil {
			return nil, err
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
		keyValue[key.Key], err = c.get(ctx, key.Key, key.Opts.Server)
		if err != nil {
			return nil, err
		}
	}

	return keyValue, nil
}
