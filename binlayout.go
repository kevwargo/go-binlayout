package binlayout

import (
	"encoding/binary"
	"errors"
	"io"
	"reflect"
)

var ErrUndefinedByteOrder = errors.New("byte order is not defined")

func Read(r io.Reader, data any) error {
	s := state{
		r:   r,
		val: reflect.ValueOf(data),
	}
	s.kind = s.val.Kind()
	s.typ = s.val.Type()

	switch s.kind {
	case reflect.Pointer:
		s.ptr = s.val
		s.val = s.val.Elem()
		s.kind = s.val.Kind()
		s.typ = s.val.Type()

		return s.read()
	default:
		return errors.New("data must be a pointer")
	}
}

type state struct {
	r    io.Reader
	val  reflect.Value
	ptr  reflect.Value
	typ  reflect.Type
	kind reflect.Kind
	bo   binary.ByteOrder
}

func (s state) read() error {
	if s.kind == reflect.Struct {
		return s.readStruct()
	}

	if s.bo == nil {
		return ErrUndefinedByteOrder
	}

	if s.isInt() {
		return s.readInt(nil)
	}

	switch s.kind {
	case reflect.Slice:
		// if
	}

	return nil
}

func (s state) readStruct() error {
	var size int64
	sizes := make(map[string]int64)

	for i := range s.typ.NumField() {
		f := s.typ.Field(i)

		if i == 0 && s.setByteOrder(f) {
			continue
		}

		fs := s.withVal(s.val.Field(i))

		if !f.IsExported() {
			if err := fs.discard(); err != nil {
				return err
			}
		}

		if fs.isInt() {
			if err := fs.readInt(&size); err != nil {
				return err
			}

			sizes[f.Name] = size
		}
		// if fv.(int) { decode; storeIntVal(f.name, fv.(int)) }
	}

	return nil
}

func (s state) isInt() bool {
	switch s.kind {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	}

	return false
}

func (s state) readInt(res *int64) error {
	if err := binary.Read(s.r, s.bo, s.ptr.Interface()); err != nil {
		return err
	}

	if res != nil {
		if s.val.CanInt() {
			*res = s.val.Int()
		} else if s.val.CanUint() {
			*res = int64(s.val.Uint())
		}
	}

	return nil
}

func (s *state) setByteOrder(f reflect.StructField) bool {
	if f.Type != typEmptyStruct || f.Name != "_" {
		return false
	}

	val, ok := f.Tag.Lookup(tagEndian)
	if !ok {
		return false
	}

	switch val {
	case "big":
		s.bo = binary.BigEndian
	case "little":
		s.bo = binary.LittleEndian
	default:
		return false
	}

	return true
}

func (s state) withVal(val reflect.Value) state {
	s.val = val
	s.ptr = val.Addr()
	s.kind = val.Kind()
	s.typ = val.Type()

	return s
}

func (s state) withBO(bo binary.ByteOrder) state {
	s.bo = bo

	return s
}

var typEmptyStruct = reflect.TypeFor[struct{}]()

const (
	tagEndian = "bl-endian"
	tagMagic  = "bvalid-magic-hex"
	tagValues = "bvalid-values"
	tagBytes  = "bl-bytes"

	maxFieldSize int64 = 1 << 34
)
