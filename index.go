package itisadb

import (
	"context"
	"errors"
	"github.com/egorgasay/itisadb-go-sdk/api/balancer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Index struct {
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
		if !ok {
			return err
		}
		if st.Code() == codes.AlreadyExists {
			return ErrUniqueConstraint
		}
		if st.Code() == codes.Unavailable {
			return ErrUnavailable
		}
		return err
	}

	return nil
}

// Get gets the value for the key from the specified index.
func (i *Index) Get(ctx context.Context, key string) (string, error) {
	res, err := i.cl.GetFromIndex(ctx, &balancer.BalancerGetFromIndexRequest{
		Key:   key,
		Index: i.name,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return "", err
		}

		if st.Code() == codes.ResourceExhausted {
			return "", ErrIndexNotFound
		}

		if st.Code() == codes.NotFound {
			return "", ErrNotFound
		}

		if st.Code() == codes.Unavailable {
			return "", ErrUnavailable
		}

		return "", err
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
		st, ok := status.FromError(err)
		if !ok {
			return nil, err
		}
		if st.Code() == codes.NotFound {
			return nil, ErrNotFound
		}
		if st.Code() == codes.Unavailable {
			return nil, ErrUnavailable
		}
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
		st, ok := status.FromError(err)
		if !ok {
			return nil, err
		}
		if st.Code() == codes.NotFound {
			return nil, ErrIndexNotFound
		}
		if st.Code() == codes.Unavailable {
			return nil, ErrUnavailable
		}
		return nil, err
	}

	return r.GetIndex(), nil
}

// Size returns  the size of the index.
func (i *Index) Size(ctx context.Context) (uint64, error) {
	r, err := i.cl.Size(ctx, &balancer.BalancerIndexSizeRequest{
		Name: i.name,
	})

	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return 0, err
		}
		if st.Code() == codes.NotFound {
			return 0, ErrIndexNotFound
		}
		if st.Code() == codes.Unavailable {
			return 0, ErrUnavailable
		}
		return 0, err
	}

	return r.GetSize(), nil
}

// DeleteIndex deletes the index.
func (i *Index) DeleteIndex(ctx context.Context) error {
	_, err := i.cl.DeleteIndex(ctx, &balancer.BalancerDeleteIndexRequest{
		Index: i.name,
	})

	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return err
		}
		if st.Code() == codes.NotFound {
			return ErrIndexNotFound
		}
		if st.Code() == codes.Unavailable {
			return ErrUnavailable
		}
		return err
	}

	return nil
}

// Attach attaches the index to another index.
func (i *Index) Attach(ctx context.Context, name string) error {
	_, err := i.cl.AttachToIndex(ctx, &balancer.BalancerAttachToIndexRequest{
		Dst: i.name,
		Src: name,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return err
		}
		if st.Code() == codes.NotFound {
			return ErrIndexNotFound
		}
		if st.Code() == codes.Unavailable {
			return ErrUnavailable
		}
		return err
	}
	return nil
}

// DeleteAttr deletes the attribute from the index.
func (i *Index) DeleteAttr(ctx context.Context, key string) error {
	_, err := i.cl.DeleteAttr(ctx, &balancer.BalancerDeleteAttrRequest{
		Index: i.name,
		Key:   key,
	})

	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return err
		}
		if st.Code() == codes.ResourceExhausted {
			return ErrNotFound
		}
		if st.Code() == codes.NotFound {
			return ErrIndexNotFound
		}
		if st.Code() == codes.Unavailable {
			return ErrUnavailable
		}
		return err
	}

	return nil
}
