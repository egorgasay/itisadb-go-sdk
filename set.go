package itisadb

import (
	"context"
	"errors"
	"github.com/egorgasay/itisadb-go-sdk/api"
)

const (
	setDefault = -iota
	setToAll
)

var ErrUniqueConstraint = errors.New("unique constraint failed")

func (c *Client) SetTo(ctx context.Context, key, value string, server int32, uniques bool) (res Result[int32]) {
	r, err := c.cl.Set(withAuth(ctx), &api.SetRequest{
		Key:     key,
		Value:   value,
		Server:  &server,
		Uniques: uniques,
	})

	if err != nil {
		res.err = convertGRPCError(err)
	} else {
		res.value = r.SavedTo
	}

	return res
}

// SetOne sets the value for the key to gRPCis.
func (c *Client) SetOne(ctx context.Context, key, value string, uniques bool) Result[bool] {
	server, err := c.SetTo(ctx, key, value, setDefault, uniques).ValueAndErr()
	if err != nil {
		return Result[bool]{err: err}
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.keysAndServers[key] = server
	return Result[bool]{value: true}
}

// SetToAll sets the value for the key on all servers.
func (c *Client) SetToAll(ctx context.Context, key, value string, uniques bool) Result[bool] {
	r := c.SetTo(ctx, key, value, setToAll, uniques)
	if r.Err() != nil {
		return Result[bool]{err: convertGRPCError(r.Err())}
	}

	return Result[bool]{value: true}
}

// SetMany sets a set of values for gRPCis.
func (c *Client) SetMany(ctx context.Context, keyValue map[string]string, uniques bool) Result[bool] {
	for key, value := range keyValue {
		err := c.SetTo(ctx, key, value, setDefault, uniques).Err()
		if err != nil {
			return Result[bool]{err: convertGRPCError(err)}
		}
	}
	return Result[bool]{value: true}
}

// SetManyOpts gets a lot of values from gRPCis with opts.
func (c *Client) SetManyOpts(ctx context.Context, keyValue map[string]Value, uniques bool) Result[bool] {
	for key, value := range keyValue {
		err := c.SetTo(ctx, key, value.Value, value.Opts.Server, uniques).Err()
		if err != nil {
			return Result[bool]{err: convertGRPCError(err)}
		}
	}
	return Result[bool]{value: true}
}
