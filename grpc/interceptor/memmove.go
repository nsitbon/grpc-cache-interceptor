package interceptor

import (
	"fmt"
	"reflect"
	"unsafe"
)

func MemMove(dest, src interface{}) error {
	if err := checkParamsType(dest, src); err != nil {
		return err
	}

	copy(interfaceToSlice(dest), interfaceToSlice(src))
	return nil
}

func checkParamsType(dest interface{}, src interface{}) error {
	if dest == nil {
		return fmt.Errorf("MemMove requires arguments to be non nil: dest (nil)")
	} else if src == nil {
		return fmt.Errorf("MemMove requires arguments to be non nil: src (nil)")
	} else if !hasPointerType(dest) {
		return fmt.Errorf("MemMove requires arguments to be of pointer type: dest (%T)", dest)
	} else if !hasPointerType(src) {
		return fmt.Errorf("MemMove requires arguments to be of pointer type: src (%T)", src)
	} else if !haveSameType(src, dest) {
		return fmt.Errorf("MemMove requires arguments of the same type: src (%T), dest (%T)", src, dest)
	}
	return nil
}

func hasPointerType(dest interface{}) bool {
	return reflect.ValueOf(dest).Kind() == reflect.Ptr
}

func haveSameType(src interface{}, dest interface{}) bool {
	return reflect.ValueOf(src).Type() == reflect.ValueOf(dest).Type()
}

func interfaceToSlice(req interface{}) []byte {
	return buildByteSlice(getPointerValue(req), getPointedTypeSize(req))
}

/*
	reflect.SliceHeader is the internal structure used by Go to represent Slice.
*/
func buildByteSlice(data uintptr, size int) []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{Data: data, Len: size, Cap: size}))
}

func getPointedTypeSize(req interface{}) int {
	return int(reflect.Indirect(reflect.ValueOf(req)).Type().Size())
}

func getPointerValue(req interface{}) uintptr {
	return reflect.ValueOf(req).Pointer()
}
