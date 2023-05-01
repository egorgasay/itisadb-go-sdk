package itisadb

import (
	"context"
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

// Index creates a new area.
func (c *Client) Index(ctx context.Context, name string) (*Index, error) {
	_, err := c.cl.Index(ctx, &balancer.BalancerIndexRequest{
		Name: name,
	})

	if err != nil {
		return nil, err
	}

	return &Index{
		name: name,
		cl:   c.cl,
	}, nil
}

// IsIndex checks if it is an index or not.
func (c *Client) IsIndex(ctx context.Context, name string) (bool, error) {
	r, err := c.cl.IsIndex(ctx, &balancer.BalancerIsIndexRequest{
		Name: name,
	})

	if err != nil {
		return false, err
	}

	return r.Ok, nil
}
