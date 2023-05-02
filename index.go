package itisadb

import (
	"context"
	"errors"
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

var ErrIndexNotFound = errors.New("index not found")

// Set sets the value for the key in the specified index.
func (i *Index) Set(ctx context.Context, key, value string, uniques bool) error {
	_, err := i.cl.SetToIndex(ctx, &balancer.BalancerSetToIndexRequest{
		Key:     key,
		Value:   value,
		Index:   i.name,
		Uniques: uniques,
	})

	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			return ErrUniqueConstraint
		}
		return fmt.Errorf("an unknown error occurred while setting the value in the storage: %w", err)
	}

	return nil
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
	_, err := i.cl.Index(ctx, &balancer.BalancerIndexRequest{
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

// GetName returns the name of the index.
func (i *Index) GetName() string {
	return i.name
}

// GetIndex returns the index.
func (i *Index) GetIndex(ctx context.Context) (map[string]string, error) {
	r, err := i.cl.GetIndex(ctx, &balancer.BalancerGetIndexRequest{
		Name: i.name,
	})

	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			return nil, ErrIndexNotFound
		}
		return nil, fmt.Errorf("an unknown error occurred while deleting the index: %w", err)
	}

	return r.GetIndex(), nil
}

// Size returns  the size of the index.
func (i *Index) Size(ctx context.Context) (uint64, error) {
	r, err := i.cl.Size(ctx, &balancer.BalancerIndexSizeRequest{
		Name: i.name,
	})

	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			return 0, ErrIndexNotFound
		}
		return 0, fmt.Errorf("an unknown error occurred while deleting the index: %w", err)
	}

	return r.GetSize(), nil
}

// Delete deletes the index.
func (i *Index) Delete(ctx context.Context) error {
	_, err := i.cl.DeleteIndex(ctx, &balancer.BalancerDeleteIndexRequest{
		Name: i.name,
	})

	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			return ErrIndexNotFound
		}
		return fmt.Errorf("an unknown error occurred while deleting the index: %w", err)
	}

	return nil
}
