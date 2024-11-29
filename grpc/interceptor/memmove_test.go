package interceptor

import (
	"fmt"
	"reflect"
	"testing"
)

type FakeType1 struct {
	Foo int
}

type FakeType2 struct {
	Foo string
}

type FakeType3 struct {
	FakeType1
	Bar bool
}

func TestMemMove(t *testing.T) {
	t.Parallel() // marks TLog as capable of running in parallel with other tests
	tests := []struct {
		src      interface{}
		dest     interface{}
		hasError bool
	}{
		{nil, nil, true},
		{1, nil, true},
		{nil, 2, true},
		{1, 2, true},
		{1, 2.0, true},
		{FakeType1{Foo: 1}, FakeType1{Foo: 2}, true},
		{&FakeType1{Foo: 1}, &FakeType1{Foo: 2}, false},
		{FakeType1{Foo: 1}, FakeType2{Foo: "2"}, true},
		{&FakeType1{Foo: 1}, &FakeType2{Foo: "2"}, true},

		{FakeType1{Foo: 1}, FakeType3{FakeType1: FakeType1{Foo: 2}, Bar: true}, true},
		{&FakeType1{Foo: 1}, &FakeType3{FakeType1: FakeType1{Foo: 2}, Bar: true}, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("src(%#v)/dest(%#v)", tt.src, tt.dest), func(t *testing.T) {
			t.Parallel()

			err := MemMove(tt.dest, tt.src)

			if err != nil != tt.hasError {
				expected := "not"

				if tt.hasError {
					expected = ""
				}

				t.Errorf(fmt.Sprintf("src (%#v) / dest (%#v) error was %s expected: err = %v", tt.src, tt.dest, expected, err))
			}

			if !tt.hasError && !reflect.DeepEqual(tt.dest, tt.src) {
				t.Errorf("dest and src were supposed to be equal : dest (%v) / src (%v)", tt.dest, tt.src)
			}
		})
	}
}
