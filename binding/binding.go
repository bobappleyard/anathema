package binding

import (
	"reflect"
	"sort"
	"unsafe"
)

// Binding associates names with struct offsets.
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

// Constructor describes how to create a binding for structs.
type Constructor func(reflect.StructField) (string, bool)

// ForStruct will create a binding for a struct type.
func (c Constructor) ForStruct(t reflect.Type) Binding {
	res := Binding{typ: t}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		n, ok := c(f)
		if !ok {
			continue
		}
		res.names = append(res.names, n)
		res.fields = append(res.fields, field{
			m:      mechanismFor(f.Type),
			offset: f.Offset,
		})
		sort.Sort((*sortBindings)(&res))
	}
	return res
}

// Fields returns a constructor that uses the names of the fields to create
// bindings.
func Fields() Constructor {
	return func(f reflect.StructField) (string, bool) {
		if f.PkgPath != "" {
			return "", false
		}
		return f.Name, true
	}
}

// Tag returns a constructor that uses the named tag as the basis for binding
// names.
func Tag(name string) Constructor {
	return func(f reflect.StructField) (string, bool) {
		if f.PkgPath != "" {
			return "", false
		}
		return f.Tag.Lookup(name)
	}
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

// Slice returns a binding that contains the provided names. Any names that are
// not present in the original binding appear as undefined fields.
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

// ToStrings returns a slice containing the string representation of fields on
// the provided struct value.
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

// FromStrings parses the slice of strings into the provided struct value.
func (b Binding) FromStrings(s []string, v reflect.Value) error {
	p := unsafe.Pointer(v.Elem().UnsafeAddr())
	for i, f := range b.fields {
		err := f.write(p, s[i])
		if err != nil {
			return err
		}
	}
	return nil
}

// FromFunc maps the names from the binding through the func. This provides
// values to be parsed into the provided struct value.
func (b Binding) FromFunc(fn func(string) (string, bool), v reflect.Value) error {
	p := unsafe.Pointer(v.Elem().UnsafeAddr())
	for i, n := range b.names {
		f := b.fields[i]
		v, ok := fn(n)
		if !ok {
			continue
		}
		err := f.write(p, v)
		if err != nil {
			return err
		}
	}
	return nil
}

// Defined will return false if any of the fields are undefined.
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
