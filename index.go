package itisadb

import (
	"context"
	"errors"
	"fmt"
	"github.com/egorgasay/itisadb-go-sdk/api/balancer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Object struct {
	name string
	cl   balancer.BalancerClient
}

var ErrObjectNotFound = errors.New("object not found")

// Set sets the value for the key in the specified object.
func (i *Object) Set(ctx context.Context, key, value string, uniques bool) error {
	_, err := i.cl.SetToObject(ctx, &balancer.BalancerSetToObjectRequest{
		Key:     key,
		Value:   value,
		Object:  i.name,
		Uniques: uniques,
	})

	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return err
		}

		if st.Code() == codes.ResourceExhausted {
			return ErrObjectNotFound
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

// Get gets the value for the key from the specified object.
func (i *Object) Get(ctx context.Context, key string) (string, error) {
	res, err := i.cl.GetFromObject(ctx, &balancer.BalancerGetFromObjectRequest{
		Key:    key,
		Object: i.name,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return "", err
		}

		if st.Code() == codes.ResourceExhausted {
			return "", ErrObjectNotFound
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

// Object returns a new or an existing object.
func (i *Object) Object(ctx context.Context, name string) (*Object, error) {
	name = fmt.Sprint(i.name, ".", name)
	_, err := i.cl.Object(ctx, &balancer.BalancerObjectRequest{
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

	return &Object{
		name: name,
		cl:   i.cl,
	}, nil
}

// GetName returns the name of the object.
func (i *Object) GetName() string {
	return i.name
}

// JSON returns the object in JSON.
func (i *Object) JSON(ctx context.Context) (string, error) {
	r, err := i.cl.ObjectToJSON(ctx, &balancer.BalancerObjectToJSONRequest{
		Name: i.name,
	})

	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return "", err
		}
		if st.Code() == codes.NotFound {
			return "", ErrObjectNotFound
		}
		if st.Code() == codes.Unavailable {
			return "", ErrUnavailable
		}
		return "", err
	}

	return r.GetObject(), nil
}

// Size returns  the size of the object.
func (i *Object) Size(ctx context.Context) (uint64, error) {
	r, err := i.cl.Size(ctx, &balancer.BalancerObjectSizeRequest{
		Name: i.name,
	})

	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return 0, err
		}
		if st.Code() == codes.NotFound {
			return 0, ErrObjectNotFound
		}
		if st.Code() == codes.Unavailable {
			return 0, ErrUnavailable
		}
		return 0, err
	}

	return r.GetSize(), nil
}

// DeleteObject deletes the object.
func (i *Object) DeleteObject(ctx context.Context) error {
	_, err := i.cl.DeleteObject(ctx, &balancer.BalancerDeleteObjectRequest{
		Object: i.name,
	})

	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return err
		}
		if st.Code() == codes.NotFound {
			return ErrObjectNotFound
		}
		if st.Code() == codes.Unavailable {
			return ErrUnavailable
		}
		return err
	}

	return nil
}

// Attach attaches the object to another object.
func (i *Object) Attach(ctx context.Context, name string) error {
	_, err := i.cl.AttachToObject(ctx, &balancer.BalancerAttachToObjectRequest{
		Dst: i.name,
		Src: name,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return err
		}
		if st.Code() == codes.NotFound {
			return ErrObjectNotFound
		}
		if st.Code() == codes.Unavailable {
			return ErrUnavailable
		}
		return err
	}
	return nil
}

// DeleteAttr deletes the attribute from the object.
func (i *Object) DeleteAttr(ctx context.Context, key string) error {
	_, err := i.cl.DeleteAttr(ctx, &balancer.BalancerDeleteAttrRequest{
		Object: i.name,
		Key:    key,
	})

	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return err
		}
		if st.Code() == codes.ResourceExhausted {
			return ErrObjectNotFound
		}
		if st.Code() == codes.NotFound {
			return ErrNotFound
		}
		if st.Code() == codes.Unavailable {
			return ErrUnavailable
		}
		return err
	}

	return nil
}
