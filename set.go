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

func (c *Client) set(ctx context.Context, key, value string, server int32) (int32, error) {
	res, err := c.cl.Set(ctx, &balancer.BalancerSetRequest{
		Key:    key,
		Value:  value,
		Server: server,
	})

	if err != nil {
		return 0, fmt.Errorf("an error occurred while setting the value in the storage: %w", err)
	}

	return res.SavedTo, nil
}

// SetOne sets the value for the key to gRPCis.
func (c *Client) SetOne(ctx context.Context, key, value string) error {
	server, err := c.set(ctx, key, value, setDefault)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.keysAndServers[key] = server
	return nil
}

// SetTo sets the value for the key on the specified server.
func (c *Client) SetTo(ctx context.Context, key, value string, server int32) error {
	_, err := c.set(ctx, key, value, server)
	return err
}

// SetToDB sets the value for the key in the physical database.
func (c *Client) SetToDB(ctx context.Context, key, value string) error {
	_, err := c.set(ctx, key, value, setToDB)
	return err
}

// SetToAll sets the value for the key on all servers.
func (c *Client) SetToAll(ctx context.Context, key, value string) error {
	_, err := c.set(ctx, key, value, setToAll)
	return err
}

// SetToAllAndToDB sets the value for the key on all servers and in th physical database.
func (c *Client) SetToAllAndToDB(ctx context.Context, key, value string) error {
	_, err := c.set(ctx, key, value, setToAllAndToDB)
	return err
}

// SetMany sets a set of values for gRPCis.
func (c *Client) SetMany(ctx context.Context, keyValue map[string]string) error {
	for key, value := range keyValue {
		_, err := c.set(ctx, key, value, setDefault)
		if err != nil {
			return err
		}
	}
	return nil
}

// SetManyOpts gets a lot of values from gRPCis with opts.
func (c *Client) SetManyOpts(ctx context.Context, keyValue map[string]Value) error {
	for key, value := range keyValue {
		_, err := c.set(ctx, key, value.Value, value.Opts.Server)
		if err != nil {
			return err
		}
	}
	return nil
}
