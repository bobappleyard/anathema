package di

import (
	"errors"
	"github.com/bobappleyard/anathema/assert"
	"reflect"
	"strconv"
	"testing"
)

type testProcess struct {
	values map[reflect.Type]reflect.Value
}

func (s *testProcess) withInstance(x interface{}) *testProcess {
	if s.values == nil {
		s.values = map[reflect.Type]reflect.Value{}
	}
	s.values[reflect.TypeOf(x)] = reflect.ValueOf(x)
	return s
}

func (s *testProcess) Furnish(p interface{}) error {
	return furnish(s, p)
}

func (s *testProcess) RequireValue(t reflect.Type) (reflect.Value, error) {
	v, ok := s.values[t]
	if !ok {
		return v, ErrInjectionFailed
	}
	return v, nil
}

func (s *testProcess) Apply(b Builder, t reflect.Type) {
	for _, v := range s.values {
		b.Constructor(&instanceRule{v})
	}
}

func TestApplyConstructors(t *testing.T) {
	must := func(r Rule, err error) Rule {
		if err != nil {
			panic(err)
		}
		return r
	}

	type testStruct struct {
		A int
	}

	for _, test := range []struct {
		name    string
		rules   []Rule
		request reflect.Type
		match   bool
	}{
		{
			"instance rule",
			[]Rule{
				Instance(123),
				Instance("hello"),
			},
			reflect.TypeOf(123),
			true,
		},
		{
			"factory rule",
			[]Rule{
				must(Factory(func() int { return 123 })),
				must(Factory(func() (string, error) { return "abc", nil })),
			},
			reflect.TypeOf(123),
			true,
		},
		{
			"slice rule",
			[]Rule{
				new(sliceRule),
			},
			reflect.SliceOf(reflect.TypeOf(new(interface{})).Elem()),
			true,
		},
		{
			"struct rule",
			[]Rule{
				new(structRule),
			},
			reflect.TypeOf(testStruct{}),
			true,
		},
		{
			"pointer rule",
			[]Rule{
				new(ptrRule),
			},
			reflect.TypeOf(&testStruct{}),
			true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			var b builder
			for _, r := range test.rules {
				r.Apply(&b, test.request)
			}
			cons := b.constructors
			if len(cons) == 0 {
				cons = b.fallback
			}
			ts := make([]reflect.Type, len(cons))
			for i, c := range cons {
				ts[i] = c.WillCreate()
			}
			if test.match {
				assert.Equal(t, ts, []reflect.Type{test.request})
			}
		})
	}
}

type testInjected struct{}

func (*testInjected) Inject() {}

func TestApplyMutators(t *testing.T) {
	must := func(r Rule, err error) Rule {
		if err != nil {
			panic(err)
		}
		return r
	}

	for _, test := range []struct {
		name    string
		rules   []Rule
		request reflect.Type
		ok      bool ``
	}{
		{
			"instance no mutators",
			[]Rule{
				Instance(123),
			},
			reflect.TypeOf(123),
			false,
		},
		{
			"factory no mutators",
			[]Rule{
				must(Factory(func() int { return 123 })),
			},
			reflect.TypeOf(123),
			false,
		},
		{
			"pointer no mutators",
			[]Rule{
				new(ptrRule),
			},
			reflect.PtrTo(reflect.TypeOf(123)),
			false,
		},
		{
			"injected",
			[]Rule{
				new(injectedRule),
			},
			reflect.TypeOf(new(testInjected)),
			true,
		},
		{
			"struct",
			[]Rule{
				new(structRule),
			},
			reflect.TypeOf(testInjected{}),
			true,
		},
		{
			"slice",
			[]Rule{
				new(sliceRule),
			},
			reflect.TypeOf([]interface{}{}),
			false,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			var b builder
			for _, r := range test.rules {
				r.Apply(&b, test.request)
			}
			if test.ok != (len(b.mutators) != 0) {
				t.Fail()
			}
		})
	}
}

func TestFactoryConstructor(t *testing.T) {
	for _, test := range []struct {
		name string
		fn   interface{}
		ok   bool
	}{
		{
			"no args ok",
			func() int { return 0 },
			true,
		},
		{
			"args ok",
			func(s string, t int) int { return 0 },
			true,
		},
		{
			"error ok",
			func() (int, error) { return 0, nil },
			true,
		},
		{
			"non func",
			0,
			false,
		},
		{
			"non error second arg",
			func() (int, int) { return 0, 0 },
			false,
		},
		{
			"no returns",
			func() {},
			false,
		},
		{
			"too many returns",
			func() (int, int, error) { return 0, 0, nil },
			false,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, err := Factory(test.fn)
			if test.ok != (err == nil) {
				t.Fail()
			}
		})
	}
}

func TestFactoryCreate(t *testing.T) {
	injector := &factoryRule{reflect.ValueOf(strconv.Atoi)}
	v, err := injector.Create(new(testProcess).withInstance("100"))
	assert.NoError(t, err)
	assert.Equal(t, v.Interface(), 100)
}

type injectedStruct struct {
	x int
}

func (s *injectedStruct) Inject(x int) {
	s.x = x
}

func TestInjected(t *testing.T) {
	var s injectedStruct
	sv := reflect.ValueOf(&s)
	p := new(testProcess).withInstance(10)
	err := new(injectedRule).Update(p, sv)
	assert.NoError(t, err)
	assert.Equal(t, s.x, 10)
}

var injectionError = errors.New("error calling Inject()")

type injectedStructErr struct {
	x int
}

func (s *injectedStructErr) Inject(x int) error {
	s.x = x
	return injectionError
}

func TestInjectedErr(t *testing.T) {
	var s injectedStructErr
	sv := reflect.ValueOf(&s)
	p := new(testProcess).withInstance(10)
	err := new(injectedRule).Update(p, sv)
	assert.Equal(t, err, injectionError)
	assert.Equal(t, s.x, 10)
}

func TestPtrCreate(t *testing.T) {
	p := new(testProcess).withInstance(10)
	injector := &ptrInjector{reflect.PtrTo(reflect.TypeOf(10))}
	v, err := injector.Create(p)
	assert.NoError(t, err)
	assert.Equal(t, *(v.Interface().(*int)), 10)
}

func TestSlice(t *testing.T) {
	p := new(testProcess).withInstance("hello").withInstance(20)
	injector := &sliceInjector{reflect.SliceOf(reflect.TypeOf(new(interface{})).Elem())}

	v, err := injector.Create(p)
	assert.NoError(t, err)

	contains := func(search interface{}) bool {
		xs := v.Interface().([]interface{})
		for _, x := range xs {
			if x == search {
				return true
			}
		}
		return false
	}

	if !contains("hello") {
		t.Fail()
	}
	if !contains(20) {
		t.Fail()
	}
}

type testStruct struct {
	Injected   int
	unexported int
	Tagged     int `tag:"value"`
	Inner      struct {
		Value int
	}
}

func TestStruct(t *testing.T) {
	process := new(testProcess).withInstance(10)
	injector := &structInjector{reflect.TypeOf(testStruct{})}

	v, err := injector.Create(process)
	assert.NoError(t, err)
	err = injector.Update(process, v)
	assert.NoError(t, err)

	s := v.Interface().(testStruct)
	assert.Equal(t, s.Injected, 10)
	assert.Equal(t, s.unexported, 0)
	assert.Equal(t, s.Tagged, 0)
	assert.Equal(t, s.Inner.Value, 10)
}
