// ported from https://github.com/blocto/solana-go-sdk/blob/v1.30.0/pkg/bincode/serialize.go

package zkm_runtime

import (
	"encoding/binary"
	"fmt"
	"reflect"
)

func SerializeData(data any) ([]byte, error) {
	return serializeData(reflect.ValueOf(data))
}

func MustSerializeData(data interface{}) []byte {
	serializedData, err := serializeData(reflect.ValueOf(data))
	if err != nil {
		panic(err)
	}
	return serializedData
}

func serializeData(v reflect.Value) ([]byte, error) {
	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			return []byte{1}, nil
		}
		return []byte{0}, nil
	case reflect.Int8:
		return []byte{uint8(v.Int())}, nil
	case reflect.Uint8:
		return []byte{uint8(v.Uint())}, nil
	case reflect.Int16:
		b := make([]byte, 2)
		binary.LittleEndian.PutUint16(b, uint16(v.Int()))
		return b, nil
	case reflect.Uint16:
		b := make([]byte, 2)
		binary.LittleEndian.PutUint16(b, uint16(v.Uint()))
		return b, nil
	case reflect.Int32:
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(v.Int()))
		return b, nil
	case reflect.Uint32:
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(v.Uint()))
		return b, nil
	case reflect.Int64:
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(v.Int()))
		return b, nil
	case reflect.Uint64:
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(v.Uint()))
		return b, nil
	case reflect.Slice:
		switch v.Type().Elem().Kind() {
		case reflect.Uint8:
			output := make([]byte, 8)
			binary.LittleEndian.PutUint64(output, uint64(v.Len()))

			for i := 0; i < v.Len(); i++ {
				d, err := serializeData(v.Index(i))
				if err != nil {
					return nil, err
				}
				output = append(output, d...)
			}
			return output, nil
		}
		return nil, fmt.Errorf("unsupport type: %v, elem: %v", v.Kind(), v.Elem().Kind())
	case reflect.Array:
		switch v.Type().Elem().Kind() {
		case reflect.Uint8:
			b := make([]byte, 0, v.Len())
			for i := 0; i < v.Len(); i++ {
				b = append(b, byte(v.Index(i).Uint()))
			}
			return b, nil
		}
		return nil, fmt.Errorf("unsupport type: %v, elem: %v", v.Kind(), v.Elem().Kind())
	case reflect.String:
		b := make([]byte, 8+len(v.String()))
		binary.LittleEndian.PutUint64(b, uint64(len(v.String())))
		copy(b[8:], []byte(v.String()))
		return b, nil
	case reflect.Ptr:
		if v.IsNil() {
			return []byte{0}, nil
		}
		d, err := serializeData(v.Elem())
		if err != nil {
			return nil, err
		}
		b := make([]byte, 1+len(d))
		b[0] = 1
		copy(b[1:], d[:])
		return b, nil
	case reflect.Struct:
		data := make([]byte, 0, 1024)
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			d, err := serializeData(field)
			if err != nil {
				return nil, err
			}
			data = append(data, d...)
		}
		return data, nil
	}
	return nil, fmt.Errorf("unsupport type: %v", v.Kind())
}
