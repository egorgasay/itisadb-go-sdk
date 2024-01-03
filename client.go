package itisadb

import (
	"context"
	"fmt"
	"github.com/egorgasay/gost"
	api "github.com/egorgasay/itisadb-shared-proto/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	token string

	cl api.ItisaDBClient
}

type Options struct {
	Server int32
	Unique bool
}

type ValueSpec struct {
	Value   string
	Options SetOptions
}

type KeySpec struct {
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

func New(ctx context.Context, balancerIP string, conf ...Config) (res gost.Result[*Client]) {
	var conn *grpc.ClientConn
	var err error

	config := defaultConfig
	if len(conf) > 0 {
		config = conf[0]
	}

	conn, err = grpc.Dial(balancerIP, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return res.Err(gost.NewError(0, 0, err.Error()))
	}

	client := api.NewItisaDBClient(conn)

	resp, err := client.Authenticate(ctx, &api.AuthRequest{
		Login:    config.Credentials.Login,
		Password: config.Credentials.Password,
	})

	if err != nil {
		return res.Err(gost.NewError(0, 0, fmt.Sprintf("failed to authenticate: %s", err)))
	}

	authMetadata.Set(token, resp.Token)

	return res.Ok(&Client{
		cl:    client,
		token: resp.Token,
	})
}

// Object creates a new object.
func (c *Client) Object(ctx context.Context, name string, opts ...ObjectOptions) (res gost.Result[*Object]) {
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

	if err != nil {
		return res.Err(errFromGRPCError(err))
	}

	return res.Ok(&Object{
		cl:   c.cl,
		name: name,
	})
}

// IsObject checks if it is an object or not.
func (c *Client) IsObject(ctx context.Context, name string) (res gost.Result[bool]) {
	r, err := c.cl.IsObject(withAuth(ctx), &api.IsObjectRequest{
		Name: name,
	})

	if err != nil {
		return res.Err(errFromGRPCError(err))
	}

	return res.Ok(r.Ok)
}

func ToServerNumber(x int) *int32 {
	var y = int32(x)
	return &y
}
