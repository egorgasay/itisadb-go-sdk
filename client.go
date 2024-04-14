package itisadb

import (
	"context"

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

const ObjectSeparator = "."

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

	conn, err = grpc.Dial(balancerIP,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return res.Err(gost.NewErrX(0, "grpc dial failed").Extend(0, err.Error()))
	}

	client := api.NewItisaDBClient(conn)

	resp, err := client.Authenticate(ctx, &api.AuthRequest{
		Login:    config.Credentials.Login,
		Password: config.Credentials.Password,
	})

	if err != nil {
		return res.Err(errFromGRPCError(err))
	}

	authMetadata.Set(token, resp.Token)

	return res.Ok(&Client{
		cl:    client,
		token: resp.Token,
	})
}

//type Error struct {
//	*gost.ErrX
//}
//
//const (
//	NotFound = iota
//)
//
//func (e *Error) IsNotFound() bool {
//	return e.BaseCode() == NotFound
//}
