package itisadb

import (
	"context"
	"errors"
	"github.com/egorgasay/itisadb-go-sdk/api"
)

const (
	setDefault = -iota
	setToDisk
	setToAll
	setToAllAndToDisk
)

var ErrUniqueConstraint = errors.New("unique constraint failed")

func (c *Client) set(ctx context.Context, key, value string, server int32, uniques bool) (int32, error) {
	res, err := c.cl.Set(withAuth(ctx), &api.SetRequest{
		Key:     key,
		Value:   value,
		Server:  &server,
		Uniques: uniques,
	})

	if err != nil {
		return 0, convertGRPCError(err)
	}

	return res.SavedTo, nil
}

// SetOne sets the value for the key to gRPCis.
func (c *Client) SetOne(ctx context.Context, key, value string, uniques bool) error {
	server, err := c.set(ctx, key, value, setDefault, uniques)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.keysAndServers[key] = server
	return nil
}

// SetTo sets the value for the key on the specified server.
func (c *Client) SetTo(ctx context.Context, key, value string, server int32, uniques bool) error {
	_, err := c.set(ctx, key, value, server, uniques)
	return err
}

// SetToDisk sets the value for the key in the physical database.
func (c *Client) SetToDisk(ctx context.Context, key, value string, uniques bool) error {
	_, err := c.set(ctx, key, value, setToDisk, uniques)
	return err
}

// SetToAll sets the value for the key on all servers.
func (c *Client) SetToAll(ctx context.Context, key, value string, uniques bool) error {
	_, err := c.set(ctx, key, value, setToAll, uniques)
	return err
}

// SetToAllAndToDisk sets the value for the key on all servers and in th physical database.
func (c *Client) SetToAllAndToDisk(ctx context.Context, key, value string, uniques bool) error {
	_, err := c.set(ctx, key, value, setToAllAndToDisk, uniques)
	return err
}

// SetMany sets a set of values for gRPCis.
func (c *Client) SetMany(ctx context.Context, keyValue map[string]string, uniques bool) error {
	for key, value := range keyValue {
		_, err := c.set(ctx, key, value, setDefault, uniques)
		if err != nil {
			return err
		}
	}
	return nil
}

// SetManyOpts gets a lot of values from gRPCis with opts.
func (c *Client) SetManyOpts(ctx context.Context, keyValue map[string]Value, uniques bool) error {
	for key, value := range keyValue {
		_, err := c.set(ctx, key, value.Value, value.Opts.Server, uniques)
		if err != nil {
			return err
		}
	}
	return nil
}
