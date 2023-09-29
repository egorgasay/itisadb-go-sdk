package itisadb

import (
	"context"
	"errors"
	"fmt"
	"github.com/egorgasay/itisadb-go-sdk/api/balancer"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	keysAndServers map[string]int32
	mu             sync.RWMutex

	token string

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

type Credentials struct {
	Login    string
	Password string
}

type Config struct {
	Credentials Credentials
}

var defaultConfig = Config{
	Credentials: Credentials{
		Login:    "itisadb",
		Password: "itisadb",
	},
}

func New(ctx context.Context, balancerIP string, conf ...Config) (*Client, error) {
	var conn *grpc.ClientConn
	var err error

	config := defaultConfig
	if len(conf) > 1 {
		config = conf[0]
	}

	conn, err = grpc.Dial(balancerIP, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := balancer.NewBalancerClient(conn)

	resp, err := client.Authenticate(ctx, &balancer.BalancerAuthRequest{
		Login:    config.Credentials.Login,
		Password: config.Credentials.Password,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to authenticate: %w", err)
	}

	authMetadata.Set(token, resp.Token)

	return &Client{
		keysAndServers: make(map[string]int32, 100),
		cl:             client,
		token:          resp.Token,
	}, nil
}

// Object creates a new object.
func (c *Client) Object(ctx context.Context, name string) (*Object, error) {
	_, err := c.cl.Object(ctx, &balancer.BalancerObjectRequest{
		Name: name,
	})

	if err != nil {
		return nil, convertGRPCError(err)
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
		return false, convertGRPCError(err)
	}

	return r.Ok, nil
}
