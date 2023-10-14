package itisadb

import (
	"context"
	"github.com/egorgasay/itisadb-go-sdk/api"
)

var (
	setDefault int32 = 0
	setToAll   int32 = -1
)

func (c *Client) set(ctx context.Context, key, val string, opt SetOptions) (res Result[int32]) {
	r, err := c.cl.Set(withAuth(ctx), &api.SetRequest{
		Key:   key,
		Value: val,
		Options: &api.SetRequest_Options{
			Server:   opt.Server,
			Uniques:  opt.Uniques,
			ReadOnly: opt.ReadOnly,
		},
	})

	if err != nil {
		res.err = convertGRPCError(err)
	} else {
		res.val = r.SavedTo
	}

	return res
}

// SetOne sets the val for the key to gRPCis.
func (c *Client) SetOne(ctx context.Context, key, val string, opts ...SetOptions) Result[bool] {
	opt := SetOptions{}
	if len(opts) > 0 {
		opt = opts[0]
	}

	server, err := c.set(ctx, key, val, opt).ValueAndErr()
	if err != nil {
		return Result[bool]{err: err}
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.keysAndServers[key] = server
	return Result[bool]{val: true}
}

// SetToAll sets the val for the key on all servers.
func (c *Client) SetToAll(ctx context.Context, key, val string, opts ...SetOptions) Result[bool] {
	opt := SetOptions{
		Server: &setToAll,
	}

	if len(opts) > 0 {
		opt = opts[0]
	}

	r := c.set(ctx, key, val, opt)
	if r.Err() != nil {
		return Result[bool]{err: convertGRPCError(r.Err())}
	}

	return Result[bool]{val: true}
}

// SetMany sets a set of vals for gRPCis.
func (c *Client) SetMany(ctx context.Context, kv map[string]string, opts ...SetOptions) Result[bool] {
	opt := SetOptions{}

	if len(opts) > 0 {
		opt = opts[0]
	}

	for key, val := range kv {
		err := c.set(ctx, key, val, opt).Err()
		if err != nil {
			return Result[bool]{err: convertGRPCError(err)}
		}
	}
	return Result[bool]{val: true}
}

// SetManyOpts gets a lot of vals from gRPCis with opts.
func (c *Client) SetManyOpts(ctx context.Context, keyValue map[string]Value) Result[bool] {
	for key, val := range keyValue {
		err := c.set(ctx, key, val.Value, val.Options).Err()
		if err != nil {
			return Result[bool]{err: convertGRPCError(err)}
		}
	}
	return Result[bool]{val: true}
}
