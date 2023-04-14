package grpcis

import (
	"github.com/egorgasay/grpcis-go-sdk/api/balancer"
	gcredentials "google.golang.org/grpc/credentials"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	client         balancer.BalancerClient
	keysAndServers map[string]int32
	mu             sync.RWMutex
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

	return &Client{
		client:         balancer.NewBalancerClient(conn),
		keysAndServers: make(map[string]int32, 100),
	}, nil
}
