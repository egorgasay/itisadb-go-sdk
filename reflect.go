package itisadb

import (
	"context"
	"fmt"
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
	var res Result[*Object]

	if parent != nil {
		res = parent.Object(ctx, name)
	} else {
		res = c.Object(ctx, name)
	}

	if res.Err() != nil {
		return nil, res.Err()
	}
	object = res.Val()

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
			err = object.Set(ctx, key, value.(string)).Err()
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
				err = object.Set(ctx, key, fmt.Sprint(field.Interface())).Err()
			}
		default:
			err = object.Set(ctx, key, fmt.Sprint(value)).Err()
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
func (c *Client) ObjectToStruct(ctx context.Context, name string, obj any) error {
	// checks if it is a struct
	objValue := reflect.ValueOf(obj)

	objType := objValue.Type().Kind()

	if objType != reflect.Pointer {
		return fmt.Errorf("invalid type: %s, should be pointer", objType.String())
	}

	return c.objectToStruct(ctx, name, objValue, nil)
}

func (c *Client) objectToStruct(ctx context.Context, name string, obj reflect.Value, parent *Object) (err error) {
	var res Result[*Object]
	if parent != nil {
		res = parent.Object(ctx, name)
	} else {
		res = c.Object(ctx, name)
	}

	if res.Err() != nil {
		return res.Err()
	}

	object := res.Val()

	if obj.Type().Kind() == reflect.Pointer {
		obj = obj.Elem()
	}

	val := ""

	for i := 0; i < obj.NumField(); i++ {
		field := obj.Field(i)
		fieldType := obj.Type().Field(i)
		key := fieldType.Name

	define:

		switch field.Type().Kind() {
		case reflect.String:
			val, err = object.Get(ctx, key).ValueAndErr()
			if err != nil {
				return err
			}
			field.SetString(val)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val, err = object.Get(ctx, key).ValueAndErr()
			if err != nil {
				return err
			}

			num, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return err
			}

			field.SetInt(num)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			val, err = object.Get(ctx, key).ValueAndErr()
			if err != nil {
				return err
			}

			num, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return err
			}

			field.SetUint(uint64(num))
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
		case reflect.Slice:
			// TODO: handle slices
		case reflect.Map:
			// TODO: handle maps
		case reflect.Bool:
			val, err = object.Get(ctx, key).ValueAndErr()
			if err != nil {
				return err
			}

			boolean, err := strconv.ParseBool(val)
			if err != nil {
				return err
			}

			field.SetBool(boolean)
		case reflect.Float32, reflect.Float64:
			val, err = object.Get(ctx, key).ValueAndErr()
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
