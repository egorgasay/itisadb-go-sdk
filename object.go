package itisadb

import (
	"context"
	"fmt"
	"github.com/egorgasay/gost"
	"github.com/egorgasay/itisadb-go-sdk/api"
)

type Object struct {
	name string
	cl   api.ItisaDBClient
}

// Set sets the value for the key in the specified object.
func (i *Object) Set(ctx context.Context, key, value string, opts ...SetToObjectOptions) (res gost.Result[int32]) {
	opt := SetToObjectOptions{}

	if len(opts) > 0 {
		opt = opts[0]
	}

	r, err := i.cl.SetToObject(withAuth(ctx), &api.SetToObjectRequest{
		Key:    key,
		Value:  value,
		Object: i.name,
		Options: &api.SetToObjectRequest_Options{
			Server:  opt.Server,
			Uniques: opt.Uniques,
		},
	})

	if err != nil {
		return res.Err(errFromGRPCError(err))
	}

	return res.Ok(r.SavedTo)
}

// Get gets the value for the key from the specified object.
func (i *Object) Get(ctx context.Context, key string, opts ...GetFromObjectOptions) (res gost.Result[string]) {
	opt := GetFromObjectOptions{}

	if len(opts) > 0 {
		opt = opts[0]
	}

	r, err := i.cl.GetFromObject(withAuth(ctx), &api.GetFromObjectRequest{
		Key:    key,
		Object: i.name,
		Options: &api.GetFromObjectRequest_Options{
			Server: opt.Server,
		},
	})

	if err != nil {
		return res.Err(errFromGRPCError(err))
	}

	return res.Ok(r.Value)
}

// Object returns a new or an existing object.
func (i *Object) Object(ctx context.Context, name string, opts ...ObjectOptions) (res gost.Result[*Object]) {
	opt := ObjectOptions{
		Level: Level(api.Level_DEFAULT),
	}

	if len(opts) > 0 {
		opt = opts[0]
	}

	name = fmt.Sprint(i.name, ".", name)
	_, err := i.cl.Object(withAuth(ctx), &api.ObjectRequest{
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
		name: name,
		cl:   i.cl,
	})
}

// Name returns the name of the object.
func (i *Object) Name() string {
	return i.name
}

// JSON returns the object in JSON.
func (i *Object) JSON(ctx context.Context, opts ...ObjectToJSONOptions) (res gost.Result[string]) {
	opt := ObjectToJSONOptions{}

	if len(opts) > 0 {
		opt = opts[0]
	}

	r, err := i.cl.ObjectToJSON(withAuth(ctx), &api.ObjectToJSONRequest{
		Name: i.name,
		Options: &api.ObjectToJSONRequest_Options{
			Server: opt.Server,
		},
	})

	if err != nil {
		return res.Err(errFromGRPCError(err))
	}

	return res.Ok(r.Object)
}

// Size returns  the size of the object.
func (i *Object) Size(ctx context.Context, opts ...SizeOptions) (res gost.Result[uint64]) {
	opt := SizeOptions{}

	if len(opts) > 0 {
		opt = opts[0]
	}

	r, err := i.cl.Size(withAuth(ctx), &api.ObjectSizeRequest{
		Name: i.name,
		Options: &api.ObjectSizeRequest_Options{
			Server: opt.Server,
		},
	})

	if err != nil {
		return res.Err(errFromGRPCError(err))
	}

	return res.Ok(r.Size)
}

// DeleteObject deletes the object.
func (i *Object) DeleteObject(ctx context.Context, opts ...DeleteObjectOptions) (res gost.Result[gost.Nothing]) {
	opt := DeleteObjectOptions{}

	if len(opts) > 0 {
		opt = opts[0]
	}

	_, err := i.cl.DeleteObject(withAuth(ctx), &api.DeleteObjectRequest{
		Object: i.name,
		Options: &api.DeleteObjectRequest_Options{
			Server: opt.Server,
		},
	})

	if err != nil {
		return res.Err(errFromGRPCError(err))
	}

	return res.Ok(gost.Nothing{})
}

// Attach attaches the object to another object.
func (i *Object) Attach(ctx context.Context, name string, opts ...AttachToObjectOptions) (res gost.Result[gost.Nothing]) {
	opt := AttachToObjectOptions{}

	if len(opts) > 0 {
		opt = opts[0]
	}

	_, err := i.cl.AttachToObject(withAuth(ctx), &api.AttachToObjectRequest{
		Dst: i.name,
		Src: name,
		Options: &api.AttachToObjectRequest_Options{
			Server: opt.Server,
		},
	})

	if err != nil {
		return res.Err(errFromGRPCError(err))
	}

	return res.Ok(gost.Nothing{})
}

// DeleteKey deletes the attribute from the object.
func (i *Object) DeleteKey(ctx context.Context, key string, opts ...DeleteKeyOptions) (res gost.Result[gost.Nothing]) {
	opt := DeleteKeyOptions{}

	if len(opts) > 0 {
		opt = opts[0]
	}

	_, err := i.cl.DeleteAttr(withAuth(ctx), &api.DeleteAttrRequest{
		Object: i.name,
		Key:    key,
		Options: &api.DeleteAttrRequest_Options{
			Server: opt.Server,
		},
	})

	if err != nil {
		return res.Err(errFromGRPCError(err))
	}

	return res.Ok(gost.Nothing{})
}
