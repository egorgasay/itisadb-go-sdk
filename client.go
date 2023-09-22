package itisadb

import (
	"context"
	"errors"
	"github.com/egorgasay/itisadb-go-sdk/api/balancer"
	gcredentials "google.golang.org/grpc/credentials"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	keysAndServers map[string]int32
	mu             sync.RWMutex

	cl balancer.BalancerClient
}

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

var ErrUnavailable = errors.New("storage is unavailable")

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

	client := balancer.NewBalancerClient(conn)

	return &Client{
		keysAndServers: make(map[string]int32, 100),
		cl:             client,
	}, nil
}

// Object creates a new object.
func (c *Client) Object(ctx context.Context, name string) (*Object, error) {
	_, err := c.cl.Object(ctx, &balancer.BalancerObjectRequest{
		Name: name,
	})

	if err != nil {
		return nil, err
	}

	return &Object{
		name: name,
		cl:   c.cl,
	}, nil
}

// IsObject checks if it is an object or not.
func (c *Client) IsObject(ctx context.Context, name string) (bool, error) {
	r, err := c.cl.IsObject(ctx, &balancer.BalancerIsObjectRequest{
		Name: name,
	})

	if err != nil {
		return false, err
	}

	return r.Ok, nil
}
