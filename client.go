package itisadb

import (
	"context"
	"errors"
	"fmt"
	"github.com/egorgasay/itisadb-go-sdk/api"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	keysAndServers map[string]int32
	mu             sync.RWMutex

	token string

	cl api.ItisaDBClient
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

	client := api.NewItisaDBClient(conn)

	resp, err := client.Authenticate(ctx, &api.AuthRequest{
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
func (c *Client) Object(ctx context.Context, name string) Result[*Object] {
	_, err := c.cl.Object(withAuth(ctx), &api.ObjectRequest{
		Name: name,
	})

	r := Result[*Object]{}

	if err != nil {
		r.err = convertGRPCError(err)
	} else {
		r.value = &Object{
			cl:   c.cl,
			name: name,
		}
	}

	return r
}

// IsObject checks if it is an object or not.
func (c *Client) IsObject(ctx context.Context, name string) (res Result[bool]) {
	r, err := c.cl.IsObject(withAuth(ctx), &api.IsObjectRequest{
		Name: name,
	})

	if err != nil {
		res.err = convertGRPCError(err)
	} else {
		res.value = r.Ok
	}

	return res
}
