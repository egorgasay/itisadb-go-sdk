package itisadb

import (
	"context"
	"fmt"
	"github.com/egorgasay/gost"
	"reflect"
	"strconv"
)

// StructToObject creates an object from a struct.
//
// DO NOT USE CYCLIC STRUCTURES.
func (c *Client) StructToObject(ctx context.Context, name string, structure any) (*Object, error) {
	// checks if it is a struct
	structureValue := reflect.ValueOf(structure)
	switch structureValue.Type().Kind() {
	case reflect.Struct:
		return c.structToObject(ctx, name, structure, nil)
	default:
		return nil, fmt.Errorf("invalid type: %s", structureValue.Type().Kind().String())
	}
}

func (c *Client) structToObject(ctx context.Context, name string, structure any, parent *Object) (object *Object, err error) {
	var res gost.Result[*Object]

	if parent != nil {
		res = parent.Object(ctx, name)
	} else {
		res = c.Object(ctx, name)
	}

	if res.IsErr() {
		return nil, res.Error()
	}
	object = res.Unwrap()

	// reflection is used to iterate over the struct
	// and create the object

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
			err = object.Set(ctx, key, value.(string)).Error()
		case reflect.Struct:
			_, err = c.structToObject(ctx, key, value, object)
		case reflect.Pointer:
			for field.Type().Kind() == reflect.Pointer {
				field = field.Elem() // handle nil
			}

			switch field.Type().Kind() {
			case reflect.Struct:
				_, err = c.structToObject(ctx, key, field.Interface(), object)
			default:
				err = object.Set(ctx, key, fmt.Sprint(field.Interface())).Error()
			}
		default:
			err = object.Set(ctx, key, fmt.Sprint(value)).Error()
		}

		if err != nil {
			object.DeleteObject(ctx)
			return nil, err
		}
	}
	return object, nil
}

// ObjectToStruct creates a struct from an object.
// Supported field types: Strings, Ints, Uints, Pointers to structs, Structs, Booleans, Floats.
//
// DO NOT USE CYCLIC STRUCTURES.
func (c *Client) ObjectToStruct(ctx context.Context, name string, data any) error {
	// checks if it is a struct
	objValue := reflect.ValueOf(data)

	objType := objValue.Type().Kind()

	if objType != reflect.Pointer {
		return fmt.Errorf("invalid type: %s, should be pointer", objType.String())
	}

	return c.objectToStruct(ctx, name, objValue, nil)
}

func (c *Client) objectToStruct(ctx context.Context, name string, obj reflect.Value, parent *Object) (err error) {
	var res gost.Result[*Object]

	if parent != nil {
		res = parent.Object(ctx, name)
	} else {
		res = c.Object(ctx, name)
	}

	if res.IsErr() {
		return res.Error()
	}

	object := res.Unwrap()

	if obj.Type().Kind() == reflect.Pointer {
		obj = obj.Elem()
	}

	for i := 0; i < obj.NumField(); i++ {
		field := obj.Field(i)
		fieldType := obj.Type().Field(i)
		key := fieldType.Name

	define:

		switch kind := field.Type().Kind(); kind {
		case reflect.Ptr:
			field.Set(reflect.New(field.Type().Elem()))
			if field.Type().Elem().Kind() == reflect.Struct {
				err = c.objectToStruct(ctx, key, field, object)
			} else {
				field = field.Elem()
				goto define
			}
		case reflect.Struct:
			field.Set(reflect.Zero(field.Type()))
			err = c.objectToStruct(ctx, key, field, object)
		case reflect.Slice, reflect.Map: // TODO: support slices, maps
			return fmt.Errorf("unsupported type: %s", kind.String())
		default:
			r := object.Get(ctx, key)
			if r.IsErr() {
				return err
			}
			val := r.Unwrap()

			switch kind {
			case reflect.String:
				field.SetString(val)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				num, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					return err
				}

				field.SetInt(num)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				num, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					return err
				}

				field.SetUint(uint64(num))
			case reflect.Float32, reflect.Float64:
				float, err := strconv.ParseFloat(val, 64)
				if err != nil {
					return err
				}

				field.SetFloat(float)
			case reflect.Bool:
				boolean, err := strconv.ParseBool(val)
				if err != nil {
					return err
				}

				field.SetBool(boolean)
			default:
				return fmt.Errorf("unsupported type: %s", kind.String())
			}
		}

		if err != nil {
			object.DeleteObject(ctx)
			return err
		}
	}
	return nil
}

type getter interface {
	Get(ctx context.Context, key string) (string, error)
}

var ErrWrongTypeParameter = fmt.Errorf("wrong type parameter")

func GetCmp[V comparable](ctx context.Context, from getter, key string) (val V, err error) {
	strVal, err := from.Get(ctx, key)
	if err != nil {
		return val, nil
	}

	v := reflect.ValueOf(val)
	a := func() any {
		switch v.Type().Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			num, err := strconv.ParseInt(strVal, 10, 64)
			if err != nil {
				return ErrWrongTypeParameter
			}
			return num
		case reflect.String:
			return strVal
		case reflect.Bool:
			b, err := strconv.ParseBool(strVal)
			if err != nil {
				return ErrWrongTypeParameter
			}
			return b
		case reflect.Float32, reflect.Float64:
			num, err := strconv.ParseFloat(strVal, 64)
			if err != nil {
				return ErrWrongTypeParameter
			}
			return num
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			num, err := strconv.ParseInt(strVal, 10, 64)
			if err != nil {
				return ErrWrongTypeParameter
			}
			if num < 0 {
				return ErrWrongTypeParameter
			}
		}
		return ErrWrongTypeParameter
	}()

	return a.(V), nil
}
