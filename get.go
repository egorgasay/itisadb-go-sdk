package itisadb

import (
	"context"
	"errors"
	"github.com/egorgasay/itisadb-go-sdk/api/balancer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	_ = -iota
	getFromDB
)

var ErrNotFound = errors.New("not found")

func (c *Client) get(ctx context.Context, key string, server int32) (string, error) {
	res, err := c.cl.Get(ctx, &balancer.BalancerGetRequest{
		Key:    key,
		Server: server,
	})

	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return "", err
		}

		if st.Code() == codes.NotFound {
			return "", ErrNotFound
		}

		if st.Code() == codes.Unavailable {
			return "", ErrUnavailable
		}

		return "", err
	}

	return res.Value, nil
}

// GetOne gets the value by the key from gRPCis.
func (c *Client) GetOne(ctx context.Context, key string) (string, error) {
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

// GetFromDB gets the value by key from the physical database.
func (c *Client) GetFromDB(ctx context.Context, key string) (string, error) {
	return c.get(ctx, key, getFromDB)
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
