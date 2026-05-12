package binlayout

import (
	"encoding/binary"
	"errors"
	"io"
	"reflect"
)

func Read(r io.Reader, data any) error {
	switch rv := reflect.ValueOf(data); rv.Kind() {
	case reflect.Pointer:
		return read(r, rv.Elem(), nil)
	case reflect.Slice:
		return fillSlice(r, rv)
	default:
		return errors.New("data must be a slice or a pointer")
	}
}

func read(r io.Reader, rv reflect.Value, byteOrder binary.ByteOrder) error {
	if rv.Kind() == reflect.Struct {
		return readStruct(r, rv, byteOrder)
	}

	if byteOrder == nil {
		return errors.New("byte order is not set")
	}

	switch rv.Kind() {
	case reflect.Int8:
	case reflect.Int16:
	case reflect.Int32:
	case reflect.Int64:
	case reflect.Uint8:
	case reflect.Uint16:
	case reflect.Uint32:
	case reflect.Uint64:
	}
}

func readStruct(r io.Reader, rv reflect.Value, byteOrder binary.ByteOrder) error {
	return nil
}

func fillSlice(r io.Reader, rv reflect.Value) error {
	return nil
}

const (
	tagEndian = "bl-endian"
	tagMagic  = "bvalid-magic-hex"
	tagValues = "bvalid-values"
	tagBytes  = "bl-bytes"
)
