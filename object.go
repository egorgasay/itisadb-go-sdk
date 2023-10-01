package itisadb

import (
	"context"
	"errors"
	"fmt"
	"github.com/egorgasay/itisadb-go-sdk/api"
)

type Object struct {
	name string
	cl   api.ItisaDBClient
}

var ErrObjectNotFound = errors.New("object not found")

// Set sets the value for the key in the specified object.
func (i *Object) Set(ctx context.Context, key, value string, uniques bool) (res Result[bool]) {
	_, err := i.cl.SetToObject(withAuth(ctx), &api.SetToObjectRequest{
		Key:     key,
		Value:   value,
		Object:  i.name,
		Uniques: uniques,
	})

	if err != nil {
		res.err = convertGRPCError(err)
	} else {
		res.value = true
	}

	return res
}

// Get gets the value for the key from the specified object.
func (i *Object) Get(ctx context.Context, key string) (res Result[string]) {
	r, err := i.cl.GetFromObject(withAuth(ctx), &api.GetFromObjectRequest{
		Key:    key,
		Object: i.name,
	})

	if err != nil {
		res.err = convertGRPCError(err)
	} else {
		res.value = r.Value
	}

	return res
}

// Object returns a new or an existing object.
func (i *Object) Object(ctx context.Context, name string) (res Result[*Object]) {
	name = fmt.Sprint(i.name, ".", name)
	_, err := i.cl.Object(withAuth(ctx), &api.ObjectRequest{
		Name: name,
	})

	if err != nil {
		res.err = convertGRPCError(err)
	} else {
		res.value = &Object{
			name: name,
			cl:   i.cl,
		}
	}

	return res
}

// GetName returns the name of the object.
func (i *Object) GetName() string {
	return i.name
}

// JSON returns the object in JSON.
func (i *Object) JSON(ctx context.Context) (res Result[string]) {
	r, err := i.cl.ObjectToJSON(withAuth(ctx), &api.ObjectToJSONRequest{
		Name: i.name,
	})

	if err != nil {
		res.err = convertGRPCError(err)
	} else {
		res.value = r.Object
	}

	return res
}

// Size returns  the size of the object.
func (i *Object) Size(ctx context.Context) (res Result[uint64]) {
	r, err := i.cl.Size(withAuth(ctx), &api.ObjectSizeRequest{
		Name: i.name,
	})

	if err != nil {
		res.err = convertGRPCError(err)
	} else {
		res.value = r.Size
	}

	return
}

// DeleteObject deletes the object.
func (i *Object) DeleteObject(ctx context.Context) (res Result[bool]) {
	_, err := i.cl.DeleteObject(withAuth(ctx), &api.DeleteObjectRequest{
		Object: i.name,
	})

	if err != nil {
		res.err = convertGRPCError(err)
	} else {
		res.value = true
	}

	return res
}

// Attach attaches the object to another object.
func (i *Object) Attach(ctx context.Context, name string) (res Result[bool]) {
	_, err := i.cl.AttachToObject(withAuth(ctx), &api.AttachToObjectRequest{
		Dst: i.name,
		Src: name,
	})

	if err != nil {
		res.err = convertGRPCError(err)
	} else {
		res.value = true
	}

	return res
}

// DeleteAttr deletes the attribute from the object.
func (i *Object) DeleteAttr(ctx context.Context, key string) (res Result[bool]) {
	_, err := i.cl.DeleteAttr(withAuth(ctx), &api.DeleteAttrRequest{
		Object: i.name,
		Key:    key,
	})

	if err != nil {
		res.err = convertGRPCError(err)
	} else {
		res.value = true
	}

	return res
}
