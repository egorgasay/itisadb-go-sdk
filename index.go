package itisadb

import (
	"context"
	"fmt"
	"github.com/egorgasay/itisadb-go-sdk/api/balancer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Index struct {
	next *Index
	name string
	cl   balancer.BalancerClient
}

// Set sets the value for the key in the specified index.
func (i *Index) Set(ctx context.Context, key, value string, uniques bool) (int32, error) {
	res, err := i.cl.SetToIndex(ctx, &balancer.BalancerSetToIndexRequest{
		Key:     key,
		Value:   value,
		Index:   i.name,
		Uniques: uniques,
	})

	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			return 0, ErrUniqueConstraint
		}
		return 0, fmt.Errorf("an unknown error occurred while setting the value in the storage: %w", err)
	}

	return res.SavedTo, nil
}

// Get gets the value for the key from the specified index.
func (i *Index) Get(ctx context.Context, key string) (string, error) {
	res, err := i.cl.GetFromIndex(ctx, &balancer.BalancerGetFromIndexRequest{
		Key:   key,
		Index: i.name,
	})
	if err != nil { // TODO: add more info
		return "", fmt.Errorf("an unknown error occurred while getting the value from the area: %w", err)
	}
	return res.Value, nil
}

// Index returns a new or an existing index.
func (i *Index) Index(ctx context.Context, name string) (*Index, error) {
	name = i.name + "/" + name
	_, err := i.cl.Index(ctx, &balancer.IndexRequest{
		Name: name,
	})

	if err != nil {
		return nil, err
	}

	return &Index{
		name: name,
		cl:   i.cl,
	}, nil
}
