package itisadb

import (
	"context"
	"fmt"

	"github.com/egorgasay/gost"
	api "github.com/egorgasay/itisadb-shared-proto/go"
)

type Object struct {
	name string
	opt  ObjectOptions

	cl api.ItisaDBClient
}

func (c *Client) Object(name string, opts ...ObjectOptions) *Object {
	opt := ObjectOptions{}

	if len(opts) > 0 {
		opt = opts[0]
	}

	return &Object{
		name: name,
		opt:  opt,
		cl:   c.cl,
	}
}

func (o *Object) Create(ctx context.Context) (res gost.Result[*Object]) {
	_, err := o.cl.Object(withAuth(ctx), &api.ObjectRequest{
		Name: o.name,
		Options: &api.ObjectRequest_Options{
			Server: o.opt.Server,
			Level:  api.Level(o.opt.Level),
		},
	})

	if err != nil {
		return res.Err(errFromGRPCError(err))
	}

	return res.Ok(o)
}

// Server returns the server ID of the object.
func (o *Object) Server() int32 {
	return o.opt.Server
}

// Set sets the value for the key in the specified object.
func (o *Object) Set(ctx context.Context, key, value string, opts ...SetToObjectOptions) (res gost.Result[int32]) {
	opt := SetToObjectOptions{Server: o.opt.Server}

	if len(opts) > 0 {
		opt = opts[0]
	}

	r, err := o.cl.SetToObject(withAuth(ctx), &api.SetToObjectRequest{
		Key:    key,
		Value:  value,
		Object: o.name,
		Options: &api.SetToObjectRequest_Options{
			Server:   opt.Server,
			ReadOnly: opt.ReadOnly,
		},
	})

	if err != nil {
		return res.Err(errFromGRPCError(err))
	}

	return res.Ok(r.SavedTo)
}

// Get gets the value for the key from the specified object.
func (o *Object) Get(ctx context.Context, key string, opts ...GetFromObjectOptions) (res gost.Result[string]) {
	opt := GetFromObjectOptions{Server: o.opt.Server}

	if len(opts) > 0 {
		opt = opts[0]
	}

	r, err := o.cl.GetFromObject(withAuth(ctx), &api.GetFromObjectRequest{
		Key:    key,
		Object: o.name,
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
func (o *Object) Object(name string, opts ...ObjectOptions) *Object {
	opt := ObjectOptions{}

	if len(opts) > 0 {
		opt = opts[0]
	}

	return &Object{
		name: fmt.Sprint(o.name, ObjectSeparator, name),
		cl:   o.cl,
		opt:  opt,
	}
}

// Name returns the name of the object.
func (o *Object) Name() string {
	return o.name
}

// JSON returns the object in JSON.
func (o *Object) JSON(ctx context.Context, opts ...ObjectToJSONOptions) (res gost.Result[string]) {
	opt := ObjectToJSONOptions{Server: o.opt.Server}

	if len(opts) > 0 {
		opt = opts[0]
	}

	r, err := o.cl.ObjectToJSON(withAuth(ctx), &api.ObjectToJSONRequest{
		Name: o.name,
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
func (o *Object) Size(ctx context.Context, opts ...SizeOptions) (res gost.Result[uint64]) {
	opt := SizeOptions{Server: o.opt.Server}

	if len(opts) > 0 {
		opt = opts[0]
	}

	r, err := o.cl.Size(withAuth(ctx), &api.ObjectSizeRequest{
		Name: o.name,
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
func (o *Object) DeleteObject(ctx context.Context, opts ...DeleteObjectOptions) (res gost.Result[gost.Nothing]) {
	opt := DeleteObjectOptions{Server: o.opt.Server}

	if len(opts) > 0 {
		opt = opts[0]
	}

	_, err := o.cl.DeleteObject(withAuth(ctx), &api.DeleteObjectRequest{
		Object: o.name,
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
func (o *Object) Attach(ctx context.Context, name string, opts ...AttachToObjectOptions) (res gost.Result[gost.Nothing]) {
	opt := AttachToObjectOptions{Server: o.opt.Server}

	if len(opts) > 0 {
		opt = opts[0]
	}

	_, err := o.cl.AttachToObject(withAuth(ctx), &api.AttachToObjectRequest{
		Dst: o.name,
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
func (o *Object) DeleteKey(ctx context.Context, key string, opts ...DeleteKeyOptions) (res gost.Result[gost.Nothing]) {
	opt := DeleteKeyOptions{Server: o.opt.Server}

	if len(opts) > 0 {
		opt = opts[0]
	}

	_, err := o.cl.DeleteAttr(withAuth(ctx), &api.DeleteAttrRequest{
		Object: o.name,
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

// Is checks if it is an object or not.
func (o *Object) Is(ctx context.Context) (res gost.Result[bool]) {
	r, err := o.cl.IsObject(withAuth(ctx), &api.IsObjectRequest{
		Name: o.name,
	})

	if err != nil {
		return res.Err(errFromGRPCError(err))
	}

	return res.Ok(r.Ok)
}
