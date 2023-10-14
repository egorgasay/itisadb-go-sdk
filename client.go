package itisadb

import (
	"context"
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

type Options struct {
	Server int32
	Unique bool
}

type Value struct {
	Value   string
	Options SetOptions
}

type Key struct {
	Key     string
	Options GetOptions
}

type Credentials struct {
	Login    string
	Password string
}

type Config struct {
	Credentials Credentials
}

const (
	DefaultUser     = "itisadb"
	DefaultPassword = "itisadb"
)

var defaultConfig = Config{
	Credentials: Credentials{
		Login:    DefaultUser,
		Password: DefaultPassword,
	},
}

func New(ctx context.Context, balancerIP string, conf ...Config) (*Client, error) {
	var conn *grpc.ClientConn
	var err error

	config := defaultConfig
	if len(conf) > 0 {
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
func (c *Client) Object(ctx context.Context, name string, opts ...ObjectOptions) Result[*Object] {
	opt := ObjectOptions{
		Level: Level(api.Level_DEFAULT),
	}

	if len(opts) > 0 {
		opt = opts[0]
	}

	_, err := c.cl.Object(withAuth(ctx), &api.ObjectRequest{
		Name: name,
		Options: &api.ObjectRequest_Options{
			Server: opt.Server,
			Level:  api.Level(opt.Level),
		},
	})

	r := Result[*Object]{}

	if err != nil {
		r.err = convertGRPCError(err)
	} else {
		r.val = &Object{
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
		res.val = r.Ok
	}

	return res
}

func ToServerNumber(x int) *int32 {
	var y = int32(x)
	return &y
}
