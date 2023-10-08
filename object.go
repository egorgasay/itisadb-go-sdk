package itisadb

import (
	"context"
	"fmt"
	"github.com/egorgasay/itisadb-go-sdk/api"
)

type Object struct {
	name string
	cl   api.ItisaDBClient
}

// Set sets the value for the key in the specified object.
func (i *Object) Set(ctx context.Context, key, value string, opts ...SetToObjectOptions) (res Result[bool]) {
	opt := SetToObjectOptions{}

	if len(opts) > 0 {
		opt = opts[0]
	}

	_, err := i.cl.SetToObject(withAuth(ctx), &api.SetToObjectRequest{
		Key:    key,
		Value:  value,
		Object: i.name,
		Options: &api.SetToObjectRequest_Options{
			Server:   opt.Server,
			Uniques:  opt.Uniques,
			ReadOnly: opt.ReadOnly,
		},
	})

	if err != nil {
		res.err = convertGRPCError(err)
	} else {
		res.value = true
	}

	return res
}

// Get gets the value for the key from the specified object.
func (i *Object) Get(ctx context.Context, key string, opts ...GetFromObjectOptions) (res Result[string]) {
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
		res.err = convertGRPCError(err)
	} else {
		res.value = r.Value
	}

	return res
}

// Object returns a new or an existing object.
func (i *Object) Object(ctx context.Context, name string, opts ...ObjectOptions) (res Result[*Object]) {
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
func (i *Object) JSON(ctx context.Context, opts ...ObjectToJSONOptions) (res Result[string]) {
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
		res.err = convertGRPCError(err)
	} else {
		res.value = r.Object
	}

	return res
}

// Size returns  the size of the object.
func (i *Object) Size(ctx context.Context, opts ...SizeOptions) (res Result[uint64]) {
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
		res.err = convertGRPCError(err)
	} else {
		res.value = r.Size
	}

	return
}

// DeleteObject deletes the object.
func (i *Object) DeleteObject(ctx context.Context, opts ...DeleteObjectOptions) (res Result[bool]) {
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
		res.err = convertGRPCError(err)
	} else {
		res.value = true
	}

	return res
}

// Attach attaches the object to another object.
func (i *Object) Attach(ctx context.Context, name string, opts ...AttachToObjectOptions) (res Result[bool]) {
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
		res.err = convertGRPCError(err)
	} else {
		res.value = true
	}

	return res
}

// DeleteKey deletes the attribute from the object.
func (i *Object) DeleteKey(ctx context.Context, key string, opts ...DeleteKeyOptions) (res Result[bool]) {
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
		res.err = convertGRPCError(err)
	} else {
		res.value = true
	}

	return res
}
