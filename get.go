package itisadb

import (
	"context"

	"github.com/egorgasay/gost"
	api "github.com/egorgasay/itisadb-shared-proto/go"
)

func (c *Client) get(ctx context.Context, key string, opts GetOptions) (res gost.Result[Value]) {
	r, err := c.cl.Get(withAuth(ctx), &api.GetRequest{
		Key: key,
		Options: &api.GetRequest_Options{
			Server: opts.Server,
		},
	})

	if err != nil {
		return res.Err(errFromGRPCError(err))
	}

	var level Level

	switch gost.SafeDeref(r.Level.Enum()) {
	case api.Level_DEFAULT:
		level = DefaultLevel
	case api.Level_RESTRICTED:
		level = RestrictedLevel
	case api.Level_SECRET:
		level = SecretLevel
	}

	return res.Ok(Value{
		Value:    r.Value,
		ReadOnly: r.ReadOnly,
		Level:    level,
	})
}

// GetOne gets the value by the key from gRPCis.
func (c *Client) GetOne(ctx context.Context, key string, opts ...GetOptions) (res gost.Result[Value]) {
	opt := GetOptions{}

	if len(opts) > 0 {
		opt = opts[0]
	}

	return c.get(ctx, key, opt)
}

// GetMany gets a lot of values from gRPCis.
func (c *Client) GetMany(ctx context.Context, keys []string, opts ...GetOptions) (res gost.Result[map[string]Value]) {
	if ctx.Err() != nil {
		return res.Err(gost.NewError(0, 0, ctx.Err().Error()))
	}

	var keyValue = make(map[string]Value, 10)

	opt := GetOptions{}

	if len(opts) > 0 {
		opt = opts[0]
	}

	for _, key := range keys {
		switch r := c.GetOne(ctx, key, opt); r.Switch() {
		case gost.IsOk:
			keyValue[key] = r.Unwrap()
		case gost.IsErr:
			return res.Err(r.Error())
		}
	}
	return res.Ok(keyValue)
}

// GetManyOpts gets a lot of values from gRPCis with opts.
func (c *Client) GetManyOpts(ctx context.Context, keys []KeySpec) (res gost.Result[map[string]Value]) {
	if ctx.Err() != nil {
		return res.Err(gost.NewError(0, 0, ctx.Err().Error()))
	}

	var keyValue = make(map[string]Value, 10)

	for _, key := range keys {
		switch r := c.get(ctx, key.Key, key.Options); r.Switch() {
		case gost.IsOk:
			keyValue[key.Key] = r.Unwrap()
		case gost.IsErr:
			return res.Err(r.Error())
		}
	}

	return res.Ok(keyValue)
}
