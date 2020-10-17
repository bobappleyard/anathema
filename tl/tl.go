package tl

import (
	"reflect"
	"unsafe"
)

func ListTypes() []reflect.Type {
	var res []reflect.Type
	sections, offsets := typelinks()
	for i, base := range sections {
		for _, offset := range offsets[i] {
			typeAddr := unsafe.Pointer(uintptr(base) + uintptr(offset))
			typ := reflect.TypeOf(*(*interface{})(unsafe.Pointer(&typeAddr)))
			res = append(res, typ)
		}
	}
	return res
}

func typelinks() (sections []unsafe.Pointer, offset [][]int32)
