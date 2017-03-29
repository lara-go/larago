package cache

import (
	"bytes"
	"encoding/gob"
	"reflect"
)

// Set value to the target
func setValue(target interface{}, value interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ErrorTypeMissmatch
		}
	}()

	v := reflect.ValueOf(target)
	if v.Type().Kind() != reflect.Ptr {
		return ErrorPtr
	}

	v.Elem().Set(reflect.ValueOf(value))

	return nil
}

// go binary encoder
func serializeValue(value interface{}) ([]byte, error) {
	buffer := bytes.Buffer{}

	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(value); err != nil {
		return buffer.Bytes(), ErrorSerialize
	}

	return buffer.Bytes(), nil
}

// go binary decoder
func unserializeValue(by []byte, target interface{}) error {
	buff := bytes.Buffer{}
	buff.Write(by)

	decoder := gob.NewDecoder(&buff)
	if err := decoder.Decode(target); err != nil {
		return err
	}

	return nil
}
