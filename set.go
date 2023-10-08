package itisadb

import (
	"context"
	"github.com/egorgasay/itisadb-go-sdk/api"
)

var (
	setDefault int32 = 0
	setToAll   int32 = -1
)

func (c *Client) set(ctx context.Context, key, value string, opt SetOptions) (res Result[int32]) {
	r, err := c.cl.Set(withAuth(ctx), &api.SetRequest{
		Key:   key,
		Value: value,
		Options: &api.SetRequest_Options{
			Server:   opt.Server,
			Uniques:  opt.Uniques,
			ReadOnly: opt.ReadOnly,
		},
	})

	if err != nil {
		res.err = convertGRPCError(err)
	} else {
		res.value = r.SavedTo
	}

	return res
}

// SetOne sets the value for the key to gRPCis.
func (c *Client) SetOne(ctx context.Context, key, value string, opts ...SetOptions) Result[bool] {
	opt := SetOptions{}
	if len(opts) > 0 {
		opt = opts[0]
	}

	server, err := c.set(ctx, key, value, opt).ValueAndErr()
	if err != nil {
		return Result[bool]{err: err}
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.keysAndServers[key] = server
	return Result[bool]{value: true}
}

// SetToAll sets the value for the key on all servers.
func (c *Client) SetToAll(ctx context.Context, key, value string, opts ...SetOptions) Result[bool] {
	opt := SetOptions{
		Server: &setToAll,
	}

	if len(opts) > 0 {
		opt = opts[0]
	}

	r := c.set(ctx, key, value, opt)
	if r.Err() != nil {
		return Result[bool]{err: convertGRPCError(r.Err())}
	}

	return Result[bool]{value: true}
}

// SetMany sets a set of values for gRPCis.
func (c *Client) SetMany(ctx context.Context, kv map[string]string, opts ...SetOptions) Result[bool] {
	opt := SetOptions{}

	if len(opts) > 0 {
		opt = opts[0]
	}

	for key, value := range kv {
		err := c.set(ctx, key, value, opt).Err()
		if err != nil {
			return Result[bool]{err: convertGRPCError(err)}
		}
	}
	return Result[bool]{value: true}
}

// SetManyOpts gets a lot of values from gRPCis with opts.
func (c *Client) SetManyOpts(ctx context.Context, keyValue map[string]Value) Result[bool] {
	for key, value := range keyValue {
		err := c.set(ctx, key, value.Value, value.Options).Err()
		if err != nil {
			return Result[bool]{err: convertGRPCError(err)}
		}
	}
	return Result[bool]{value: true}
}
