package grpcis

import (
	"context"
	"fmt"
	"github.com/egorgasay/grpcis-go-sdk/api/balancer"
)

const (
	setDefault = -iota
	setToDB
	setToAll
	setToAllAndToDB
)

func (c *Client) set(ctx context.Context, key, value string, server int32, uniques bool) (int32, error) {
	res, err := c.cl.Set(ctx, &balancer.BalancerSetRequest{
		Key:     key,
		Value:   value,
		Server:  server,
		Uniques: uniques,
	})

	if err != nil {
		return 0, fmt.Errorf("an error occurred while setting the value in the storage: %w", err)
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

// SetToDB sets the value for the key in the physical database.
func (c *Client) SetToDB(ctx context.Context, key, value string, uniques bool) error {
	_, err := c.set(ctx, key, value, setToDB, uniques)
	return err
}

// SetToAll sets the value for the key on all servers.
func (c *Client) SetToAll(ctx context.Context, key, value string, uniques bool) error {
	_, err := c.set(ctx, key, value, setToAll, uniques)
	return err
}

// SetToAllAndToDB sets the value for the key on all servers and in th physical database.
func (c *Client) SetToAllAndToDB(ctx context.Context, key, value string, uniques bool) error {
	_, err := c.set(ctx, key, value, setToAllAndToDB, uniques)
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
