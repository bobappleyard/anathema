package binding

import (
	"errors"
	"reflect"
	"sort"
	"unsafe"
)

type Binding struct {
	typ    reflect.Type
	names  []string
	fields []field
}

type field struct {
	m      mechanism
	offset uintptr
}

type sortBindings Binding

var ErrNotFound = errors.New("not found")

func ForStruct(t reflect.Type) Binding {
	res := Binding{typ: t}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.PkgPath != "" {
			continue
		}
		res.names = append(res.names, f.Name)
		res.fields = append(res.fields, field{
			m:      mechanismFor(f.Type),
			offset: f.Offset,
		})
		sort.Sort((*sortBindings)(&res))
	}
	return res
}

func (b *sortBindings) Len() int {
	return len(b.fields)
}

func (b *sortBindings) Less(i, j int) bool {
	return b.names[i] < b.names[j]
}

func (b *sortBindings) Swap(i, j int) {
	b.names[i], b.names[j] = b.names[j], b.names[i]
	b.fields[i], b.fields[j] = b.fields[j], b.fields[i]
}

func (b Binding) Slice(names []string) Binding {
	res := Binding{typ: b.typ}
	res.names = names
	for _, n := range names {
		f := field{m: undefined{}}
		idx := sort.Search(len(b.fields), func(i int) bool {
			return b.names[i] >= n
		})
		if idx < len(b.fields) && b.names[idx] == n {
			f = b.fields[idx]
		}
		res.fields = append(res.fields, f)
	}
	return res
}

func (b Binding) ToStrings(v reflect.Value) ([]string, error) {
	res := make([]string, len(b.fields))
	p := unsafe.Pointer(v.UnsafeAddr())
	for i, f := range b.fields {
		x, err := f.read(p)
		if err != nil {
			return nil, err
		}
		res[i] = x
	}
	return res, nil
}

func (b Binding) FromStrings(s []string) (reflect.Value, error) {
	v := reflect.New(b.typ).Elem()
	p := unsafe.Pointer(v.UnsafeAddr())
	for i, f := range b.fields {
		err := f.write(p, s[i])
		if err != nil {
			return reflect.Value{}, err
		}
	}
	return v, nil
}

func (b Binding) FromFunc(fn func(string) (string, bool)) (reflect.Value, error) {
	v := reflect.New(b.typ).Elem()
	p := unsafe.Pointer(v.UnsafeAddr())
	for i, n := range b.names {
		f := b.fields[i]
		v, ok := fn(n)
		if !ok {
			return reflect.Value{}, ErrNotFound
		}
		err := f.write(p, v)
		if err != nil {
			return reflect.Value{}, err
		}
	}
	return v, nil
}

func (b Binding) Defined() bool {
	for _, f := range b.fields {
		if !f.m.defined() {
			return false
		}
	}
	return true
}

func (f field) read(p unsafe.Pointer) (string, error) {
	return f.m.read(f.shift(p))
}

func (f field) write(p unsafe.Pointer, s string) error {
	return f.m.write(f.shift(p), s)
}

func (f field) shift(p unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) + f.offset)
}
