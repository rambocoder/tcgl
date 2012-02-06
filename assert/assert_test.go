// Tideland Common Go Library - Assert - Unit Test
//
// Copyright (C) 2012 Frank Mueller / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed 
// by the new BSD license.

package assert

//--------------------
// IMPORTS
//--------------------

import (
	"errors"
	"io"
	"testing"
)

//--------------------
// FAIL FUNCS
//--------------------

// createValueAssert returns an assert with a value logging fail func.
func createValueAssert(t *testing.T) *Assert {
	return NewAssert(func(test Test, obtained, expected interface{}, msg string) bool {
		t.Logf("testing assert '%s' failed: '%v' <> '%v' (%s)", test, obtained, expected, msg)
		return false
	})
}

// createTypeAssert returns an assert with a value description (type) logging fail func.
func createTypeAssert(t *testing.T) *Assert {
	return NewAssert(func(test Test, obtained, expected interface{}, msg string) bool {
		t.Logf("testing assert '%s' failed: '%v' <> '%v' (%s)",
			test, ValueDescription(obtained), ValueDescription(expected), msg)
		return false
	})
}

//--------------------
// TESTS
//--------------------

// Test the isNil() helper.
func TestIsNilHelper(t *testing.T) {
	if !isNil(nil) {
		t.Errorf("nil is not nil?")
	}
	if isNil("nil") {
		t.Errorf("'nil' is nil?")
	}
	var c chan int
	if !isNil(c) {
		t.Errorf("channel is not nil?")
	}
	var f func()
	if !isNil(f) {
		t.Errorf("func is not nil?")
	}
	var i interface{}
	if !isNil(i) {
		t.Errorf("interface is not nil?")
	}
	var m map[string]string
	if !isNil(m) {
		t.Errorf("map is not nil?")
	}
	var p *bool
	if !isNil(p) {
		t.Errorf("pointer is not nil?")
	}
	var s []string
	if !isNil(s) {
		t.Errorf("slice is not nil?")
	}
}

// Test the True() assertion.
func TestAssertTrue(t *testing.T) {
	a := createValueAssert(t)

	a.True(true, "should not fail")
	if a.True(false, "should fail and be logged") {
		t.Errorf("True() returned true")
	}
}

// Test the False() assertion.
func TestAssertFalse(t *testing.T) {
	a := createValueAssert(t)

	a.False(false, "should not fail")
	if a.False(true, "should fail and be logged") {
		t.Errorf("False() returned true")
	}
}

// Test the Nil() assertion.
func TestAssertNil(t *testing.T) {
	a := createValueAssert(t)

	a.Nil(nil, "should not fail")
	if a.Nil("not nil", "should fail and be logged") {
		t.Errorf("Nil() returned true")
	}
}

// Test the NotNil() assertion.
func TestAssertNotNil(t *testing.T) {
	a := createValueAssert(t)

	a.NotNil("not nil", "should not fail")
	if a.NotNil(nil, "should fail and be logged") {
		t.Errorf("NotNil() returned true")
	}
}

// Test the Equal() assertion.
func TestAssertEqual(t *testing.T) {
	a := createValueAssert(t)
	m := map[string]int{"one": 1, "two": 2, "three": 3}

	a.Equal(nil, nil, "should not fail")
	a.Equal(true, true, "should not fail")
	a.Equal(1, 1, "should not fail")
	a.Equal("foo", "foo", "should not fail")
	a.Equal(map[string]int{"one": 1, "three": 3, "two": 2}, m, "should not fail")
	if a.Equal("one", 1, "should fail and be logged") {
		t.Errorf("Equal() returned true")
	}
	if a.Equal("two", "2", "should fail and be logged") {
		t.Errorf("Equal() returned true")
	}
}

// Test the Different() assertion.
func TestAssertDifferent(t *testing.T) {
	a := createValueAssert(t)
	m := map[string]int{"one": 1, "two": 2, "three": 3}

	a.Different(nil, "nil", "should not fail")
	a.Different("true", true, "should not fail")
	a.Different(1, 2, "should not fail")
	a.Different("foo", "bar", "should not fail")
	a.Different(map[string]int{"three": 3, "two": 2}, m, "should not fail")
	if a.Different("one", "one", "should fail and be logged") {
		t.Errorf("Different() returned true")
	}
	if a.Different(2, 2, "should fail and be logged") {
		t.Errorf("Different() returned true")
	}
}

// Test the Matches() assertion.
func TestAssertMatches(t *testing.T) {
	a := createValueAssert(t)

	a.Matches("this is a test", "this.*test", "should not fail")
	a.Matches("this is 1 test", "this is [0-9] test", "should not fail")
	if a.Matches("this is a test", "foo", "should fail and be logged") {
		t.Errorf("Matches() returned true")
	}
	if a.Matches("this is a test", "this*test", "should fail and be logged") {
		t.Errorf("Matches() returned true")
	}
}

// Test the ErrorMatches() assertion.
func TestAssertErrorMatches(t *testing.T) {
	a := createValueAssert(t)
	err := errors.New("oops, an error")

	a.ErrorMatches(err, "oops, an error", "should not fail")
	a.ErrorMatches(err, "oops,.*", "should not fail")
	if a.ErrorMatches(err, "foo", "should fail and be logged") {
		t.Errorf("ErrorMatches() returned true")
	}
}

// Test the Implements() assertion.
func TestAssertImplements(t *testing.T) {
	a := createTypeAssert(t)

	var err error
	var w io.Writer

	a.Implements(errors.New("error test"), &err, "should not fail")
	if a.Implements("string test", &err, "should fail and be logged") {
		t.Errorf("Implements() returned true")
	}
	if a.Implements(errors.New("error test"), &w, "should fail and be logged") {
		t.Errorf("Implements() returned true")
	}
}

// Test the Assignable() assertion.
func TestAssertAssignable(t *testing.T) {
	a := createTypeAssert(t)

	a.Assignable(1, 5, "should not fail")
	if a.Assignable("one", 5, "should fail and be logged") {
		t.Errorf("Assignable() returned true")
	}
}

// Test the Unassignable() assertion.
func TestAssertUnassignable(t *testing.T) {
	a := createTypeAssert(t)

	a.Unassignable("one", 5, "should not fail")
	if a.Unassignable(1, 5, "should fail and be logged") {
		t.Errorf("Unassignable() returned true")
	}
}

// Test if the panic assert panics when failing.
func TestPanicAssert(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Logf("panic worked: '%v'", err)
		}
	}()

	a := NewPanicAssert()
	foo := func() {}

	a.Assignable(47, 11, "should not fail")
	a.Assignable(47, foo, "should fail")

	t.Errorf("should not be reached")
}

// Test the testing assert.
func TestTestingAssert(t *testing.T) {
	a := NewTestingAssert(t, false)
	foo := func() {}
	bar := 4711

	a.Assignable(47, 11, "should not fail")
	a.Assignable(foo, bar, "should fail")
}

// EOF
