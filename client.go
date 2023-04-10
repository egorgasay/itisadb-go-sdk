package grpcisclient

import (
	"context"
	"fmt"
	gcredentials "google.golang.org/grpc/credentials"
	"grpcis-client/api/balancer"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	client         balancer.BalancerClient
	keysAndServers map[string]int32
	mu             sync.RWMutex
}

func New(balancerIP string, credentials ...gcredentials.TransportCredentials) (*Client, error) {
	var conn *grpc.ClientConn
	var err error

	if credentials == nil || len(credentials) == 0 {
		conn, err = grpc.Dial(balancerIP, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		conn, err = grpc.Dial(balancerIP, grpc.WithTransportCredentials(credentials[0]))
	}

	if err != nil {
		return nil, err
	}

	return &Client{
		client:         balancer.NewBalancerClient(conn),
		keysAndServers: make(map[string]int32, 100),
	}, nil
}

func (c *Client) GetOnce(ctx context.Context, key string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	res, err := c.client.Get(ctx, &balancer.BalancerGetRequest{
		Key:    key,
		Server: c.keysAndServers[key],
	})

	if err != nil {
		return "", fmt.Errorf("an error occurred while getting the value in the storage: %w", err)
	}

	return res.Value, nil
}

func (c *Client) SetOnce(ctx context.Context, key, value string) error {
	res, err := c.client.Set(ctx, &balancer.BalancerSetRequest{
		Key:   key,
		Value: value,
	})

	if err != nil {
		return fmt.Errorf("an error occurred while setting the value in the storage: %w", err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.keysAndServers[key] = res.SavedTo
	return nil
}

func (c *Client) SetTo(ctx context.Context, key, value string, server int32) {}
func (c *Client) GetFrom(ctx context.Context, key, server int32)             {}

func (c *Client) SetMany(ctx context.Context, keyValue map[string]string) {}
func (c *Client) GetMany(ctx context.Context, keys []string)              {}

type Opts struct {
	Server int32
}

type Value struct {
	Value string
	Opts  Opts
}

type Key struct {
	Key  string
	Opts Opts
}

func (c *Client) SetManyOpts(ctx context.Context, keyValue map[string]Value) {}
func (c *Client) GetManyOpts(ctx context.Context, keys []Key)                {}
