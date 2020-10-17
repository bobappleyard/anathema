package di

import (
	"reflect"
	"testing"
)

type testInterface interface {
	method()
}

type testImplementationA struct{}
type testImplementationB struct{}

func (testImplementationA) method() {}
func (testImplementationB) method() {}

type testRule struct {
	t  reflect.Type
	fs []Furnisher
}

func (r *testRule) Apply(start *Scope, t reflect.Type, results []Furnisher) []Furnisher {
	if r.t.AssignableTo(t) {
		return append(results, r.fs...)
	}
	return results
}

func TestFindFurnishers(t *testing.T) {
	s := &Scope{}
	s.rules = append(s.rules, &testRule{
		reflect.TypeOf(testImplementationA{}),
		[]Furnisher{&instanceFurnisher{reflect.ValueOf(testImplementationA{})}},
	})
	s.rules = append(s.rules, &testRule{
		reflect.TypeOf(testImplementationB{}),
		[]Furnisher{&instanceFurnisher{reflect.ValueOf(testImplementationB{})}},
	})

	fs := s.findFurnishers(reflect.TypeOf(new(testInterface)).Elem())
	expect := []Furnisher{
		&instanceFurnisher{reflect.ValueOf(testImplementationA{})},
		&instanceFurnisher{reflect.ValueOf(testImplementationB{})},
	}
	if !reflect.DeepEqual(fs, expect) {
		t.Fatalf("got %v, expected %v", fs, expect)
	}
}

func TestFurnish(t *testing.T) {
	s := &Scope{}
	f := &instanceFurnisher{reflect.ValueOf(testImplementationA{})}
	s.rules = append(s.rules, &testRule{
		reflect.TypeOf(testImplementationA{}),
		[]Furnisher{f},
	})

	var got testInterface
	err := s.Furnish(reflect.ValueOf(&got))
	if err != nil {
		t.Fatalf("got error %v", err)
	}
	if got != (testImplementationA{}) {
		t.Fatalf("got %v, expecting %v", got, testImplementationA{})
	}

	f.value = reflect.Value{}
	err = s.Furnish(reflect.ValueOf(&got))
	if err != nil {
		t.Fatalf("got error %v", err)
	}
	if got != (testImplementationA{}) {
		t.Fatalf("got %v, expecting %v", got, testImplementationA{})
	}

}

type testProvider struct {
}

func (*testProvider) Provide() testImplementationA {
	return testImplementationA{}
}

func TestMethodFurnisher(t *testing.T) {
	s := &Scope{}
	m, _ := reflect.ValueOf(new(testProvider)).Type().MethodByName("Provide")
	s.rules = append(s.rules, &testRule{
		reflect.TypeOf(new(testProvider)),
		[]Furnisher{&instanceFurnisher{reflect.ValueOf(new(testProvider))}},
	})
	s.rules = append(s.rules, &testRule{
		reflect.TypeOf(testImplementationA{}),
		[]Furnisher{&methodFurnisher{m}},
	})

	var got testInterface
	err := s.Furnish(reflect.ValueOf(&got))
	if err != nil {
		t.Fatalf("got error %v", err)
	}
	if got != (testImplementationA{}) {
		t.Fatalf("got %v, expecting %v", got, testImplementationA{})
	}
}
