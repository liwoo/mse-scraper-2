package conf

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

func Parse(v interface{}) error {
	ptrRef := reflect.ValueOf(v)

	if ptrRef.Kind() != reflect.Ptr {
		return errors.New("NOT A VALID POINTER")
	}

	ref := ptrRef.Elem()
	if ref.Kind() != reflect.Struct {
		return errors.New("NOT A STRUCT")
	}

	return parseField(ref)
}

func parseField(ref reflect.Value) error {
	refType := ref.Type()
	for i := 0; i < refType.NumField(); i++ {
		refField := ref.Field(i)
		refTypeField := refType.Field(i)

		if !refField.CanSet() {
			continue
		}
		tag := refTypeField.Tag.Get("env")
		var key = refTypeField.Name
		if len(tag) > 0 {
			key = tag
		}
		value := os.Getenv(key)
		err := setField(refField, refTypeField, value)

		if err != nil {
			fmt.Println(err)
			continue
		}
	}
	return nil
}

func setField(ref reflect.Value, refTypeField reflect.StructField, value string) error {

	typee := ref.Type()
	fieldee := ref
	if typee.Kind() == reflect.Ptr {
		typee = typee.Elem()
		fieldee = ref.Elem()
	}
	switch typee.Kind() {
	case reflect.String:
		fieldee.Set(reflect.ValueOf(value))
		break
	case reflect.Bool:
		i, err := strconv.ParseBool(value)
		if err != nil {
			return nil
		}
		fieldee.Set(reflect.ValueOf(i))
		break
	case reflect.Int:
		i, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return nil
		}
		fieldee.Set(reflect.ValueOf(int(i)))
		break
	case reflect.Int64:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil
		}
		fieldee.Set(reflect.ValueOf(i))
		break
	case reflect.Float32:
		i, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return nil
		}
		fieldee.Set(reflect.ValueOf(float32(i)))
		break
	case reflect.Float64:
		i, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil
		}
		fieldee.Set(reflect.ValueOf(i))
		break
	default:
		return errors.New("TYPE NOT SUPPORTED")
	}
	return nil
}
