package itisadb

import (
	"context"
	"errors"
	"fmt"
	"github.com/egorgasay/itisadb-go-sdk/api/balancer"
	gcredentials "google.golang.org/grpc/credentials"
	"reflect"
	"strconv"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	keysAndServers map[string]int32
	mu             sync.RWMutex

	cl balancer.BalancerClient
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

var ErrUnavailable = errors.New("storage is unavailable")

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

	client := balancer.NewBalancerClient(conn)

	return &Client{
		keysAndServers: make(map[string]int32, 100),
		cl:             client,
	}, nil
}

// Index creates a new index.
func (c *Client) Index(ctx context.Context, name string) (*Index, error) {
	_, err := c.cl.Index(ctx, &balancer.BalancerIndexRequest{
		Name: name,
	})

	if err != nil {
		return nil, err
	}

	return &Index{
		name: name,
		cl:   c.cl,
	}, nil
}

// IsIndex checks if it is an index or not.
func (c *Client) IsIndex(ctx context.Context, name string) (bool, error) {
	r, err := c.cl.IsIndex(ctx, &balancer.BalancerIsIndexRequest{
		Name: name,
	})

	if err != nil {
		return false, err
	}

	return r.Ok, nil
}

// StructToIndex creates an index from a struct.
func (c *Client) StructToIndex(ctx context.Context, name string, structure any) (*Index, error) {
	// checks if it is a struct
	structureValue := reflect.ValueOf(structure)
	switch structureValue.Type().Kind() {
	case reflect.Struct:
		return c.structToIndex(ctx, name, structure, nil)
	default:
		return nil, fmt.Errorf("invalid type: %s", structureValue.Type().Kind().String())
	}
}

func (c *Client) structToIndex(ctx context.Context, name string, structure any, parent *Index) (index *Index, err error) {
	if parent != nil {
		index, err = parent.Index(ctx, name)
	} else {
		index, err = c.Index(ctx, name)
	}

	if err != nil {
		return nil, err
	}

	// reflection is used to iterate over the struct
	// and create the index

	structureValue := reflect.ValueOf(structure)
	switch structureValue.Type().Kind() {
	case reflect.Pointer:
		structureValue = structureValue.Elem()
	}

	for i := 0; i < structureValue.NumField(); i++ {
		field := structureValue.Field(i)

		fieldType := structureValue.Type().Field(i)
		key := fieldType.Name
		if !field.CanInterface() {
			continue
		}

		value := field.Interface()

		switch fieldType.Type.Kind() {
		case reflect.String:
			err = index.Set(ctx, key, value.(string), false)
		case reflect.Struct:
			_, err = c.structToIndex(ctx, key, value, index)
		case reflect.Pointer:
			for field.Type().Kind() == reflect.Pointer {
				field = field.Elem() // handle nil
			}

			switch field.Type().Kind() {
			case reflect.Struct:
				_, err = c.structToIndex(ctx, key, field.Interface(), index)
			default:
				err = index.Set(ctx, key, fmt.Sprint(field.Interface()), false)
			}
		default:
			err = index.Set(ctx, key, fmt.Sprint(value), false)
		}

		if err != nil {
			index.DeleteIndex(ctx)
			return nil, err
		}
	}
	return index, nil
}

// IndexToStruct creates a struct from an index.
// Supported field types: Strings, Ints, Uints, Pointers to structs, Structs, Booleans, Floats.
// TODO: ADD SUPPORT OF STRUCT COMMENTS
func (c *Client) IndexToStruct(ctx context.Context, name string, obj any) error {
	// checks if it is a struct
	objValue := reflect.ValueOf(obj)

	objType := objValue.Type().Kind()

	if objType != reflect.Pointer {
		return fmt.Errorf("invalid type: %s, should be pointer", objType.String())
	}

	return c.indexToStruct(ctx, name, objValue, nil)
}

func (c *Client) indexToStruct(ctx context.Context, name string, obj reflect.Value, parent *Index) (err error) {
	var index *Index
	if parent != nil {
		index, err = parent.Index(ctx, name)
	} else {
		index, err = c.Index(ctx, name)
	}

	if err != nil {
		return err
	}

	structureValue := obj.Elem()

	val := ""

	for i := 0; i < structureValue.NumField(); i++ {
		field := structureValue.Field(i)
		fieldType := structureValue.Type().Field(i)
		key := fieldType.Name

		switch fieldType.Type.Kind() {
		case reflect.String:
			val, err = index.Get(ctx, key)
			if err != nil {
				return err
			}
			structureValue.FieldByName(key).SetString(val)
		case reflect.Ptr: // TODO: handle *str etc
			switch field.Type().Elem().Kind() {
			case reflect.Struct:
				field.Set(reflect.New(field.Type().Elem()))
				err = c.indexToStruct(ctx, key, field, index)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val, err = index.Get(ctx, key)
			if err != nil {
				return err
			}

			num, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return err
			}

			field.SetInt(num)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			val, err = index.Get(ctx, key)
			if err != nil {
				return err
			}

			num, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return err
			}

			field.SetUint(uint64(num))
		case reflect.Struct:
			field.Set(reflect.New(field.Type()))
			err = c.indexToStruct(ctx, key, field, index)
		case reflect.Slice:
			// TODO: handle slices
		case reflect.Map:
			// TODO: handle maps
		case reflect.Bool:
			val, err = index.Get(ctx, key)
			if err != nil {
				return err
			}

			boolean, err := strconv.ParseBool(val)
			if err != nil {
				return err
			}

			field.SetBool(boolean)
		case reflect.Float32, reflect.Float64:
			val, err = index.Get(ctx, key)
			if err != nil {
				return err
			}

			float, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return err
			}

			field.SetFloat(float)
		}

		if err != nil {
			index.DeleteIndex(ctx)
			return err
		}
	}
	return nil
}

func GetCmp[V comparable](ctx context.Context, key V) (V, error) {
	v := reflect.ValueOf(key)
	var out = "123"
	fmt.Println(v.Type().Kind())
	a := func() any {
		switch v.Type().Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			num, err := strconv.ParseInt(out, 10, 64)
			if err != nil {
				panic("not implemented")
			}
			return num
		case reflect.String:
			return out
		}
		panic("not implemented")
	}()

	return a.(V), nil
}
