// Tideland Common Go Library - Assert
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
	"bytes"
	"fmt"
	"path"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

//--------------------
// CONST
//--------------------

const RELEASE = "Tideland Common Go Library - Assert - Release 2012-02-06"

//--------------------
// TEST
//--------------------

// Test represents the test inside an assert.
type Test uint

const (
	Invalid Test = iota
	True
	False
	Nil
	NotNil
	Equal
	Different
	Matches
	ErrorMatches
	Implements
	Assignable
	Unassignable
)

var testNames = []string{
	Invalid:      "invalid",
	True:         "true",
	False:        "false",
	Nil:          "nil",
	NotNil:       "not nil",
	Equal:        "equal",
	Different:    "different",
	Matches:      "matches",
	ErrorMatches: "error matches",
	Implements:   "implements",
	Assignable:   "assignable",
	Unassignable: "unassignable",
}

func (t Test) String() string {
	if int(t) < len(testNames) {
		return testNames[t]
	}
	return "invalid"
}

//--------------------
// FAIL FUNC
//--------------------

// FailFunc is a user defined function that will be call by an assert if
// a test fails.
type FailFunc func(test Test, obtained, expected interface{}, msg string) bool

// panicFailFunc just panics if an assert fails.
func panicFailFunc(test Test, obtained, expected interface{}, msg string) bool {
	var obex string
	switch test {
	case True, False, Nil, NotNil:
		obex = fmt.Sprintf("'%v'", obtained)
	case Implements, Assignable, Unassignable:
		obex = fmt.Sprintf("'%v' <> '%v'", ValueDescription(obtained), ValueDescription(expected))
	default:
		obex = fmt.Sprintf("'%v' <> '%v'", obtained, expected)
	}
	panic(fmt.Sprintf("assert '%s' failed: %s (%s)", test, obex, msg))
	return false
}

// generateTestingFailFunc creates a fail func bound to a testing.T.
func generateTestingFailFunc(t *testing.T, fail bool) FailFunc {
	return func(test Test, obtained, expected interface{}, msg string) bool {
		pc, file, line, _ := runtime.Caller(2)
		_, fileName := path.Split(file)
		funcNameParts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
		funcNamePartsIdx := len(funcNameParts) - 1
		funcName := funcNameParts[funcNamePartsIdx]
		buffer := &bytes.Buffer{}
		fmt.Fprintf(buffer, "--------------------------------------------------------------------------------\n")
		fmt.Fprintf(buffer, "Assert '%s' failed!\n\n", test)
		fmt.Fprintf(buffer, "Filename: %s\n", fileName)
		fmt.Fprintf(buffer, "Function: %s()\n", funcName)
		fmt.Fprintf(buffer, "Line    : %d\n", line)
		switch test {
		case True, False, Nil, NotNil:
			fmt.Fprintf(buffer, "Obtained: %v\n", obtained)
		case Implements, Assignable, Unassignable:
			fmt.Fprintf(buffer, "Obtained: %v\n", ValueDescription(obtained))
			fmt.Fprintf(buffer, "Expected: %v\n", ValueDescription(expected))
		default:
			fmt.Fprintf(buffer, "Obtained: %v\n", obtained)
			fmt.Fprintf(buffer, "Expected: %v\n", expected)
		}
		fmt.Fprintf(buffer, "Message : %s\n", msg)
		fmt.Fprintf(buffer, "--------------------------------------------------------------------------------\n")
		fmt.Print(buffer)
		if fail {
			t.Fail()
		}
		return false
	}
}

//--------------------
// ASSERT
//--------------------

// Assert instances provide the test methods.
type Assert struct {
	failFunc FailFunc
}

// NewAssert creates a new assert.
func NewAssert(ff FailFunc) *Assert {
	return &Assert{ff}
}

// NewPanicAssert creates a new assert which panics if an assert fails.
func NewPanicAssert() *Assert {
	return NewAssert(panicFailFunc)
}

// NewTestingAssert creates a new assert for use with the testing package.
func NewTestingAssert(t *testing.T, fail bool) *Assert {
	return NewAssert(generateTestingFailFunc(t, fail))
}

// True tests if obtained is true.
func (a Assert) True(obtained bool, msg string) bool {
	if obtained == false {
		return a.failFunc(True, obtained, true, msg)
	}
	return true
}

// False tests if obtained is false.
func (a Assert) False(obtained bool, msg string) bool {
	if obtained == true {
		return a.failFunc(False, obtained, false, msg)
	}
	return true
}

// Nil tests if obtained is nil.
func (a Assert) Nil(obtained interface{}, msg string) bool {
	if !isNil(obtained) {
		return a.failFunc(Nil, obtained, nil, msg)
	}
	return true
}

// NotNil tests if obtained is not nil.
func (a Assert) NotNil(obtained interface{}, msg string) bool {
	if isNil(obtained) {
		return a.failFunc(NotNil, obtained, nil, msg)
	}
	return true
}

// Equal tests if expected and obtained are equal.
func (a Assert) Equal(obtained, expected interface{}, msg string) bool {
	if !reflect.DeepEqual(obtained, expected) {
		return a.failFunc(Equal, obtained, expected, msg)
	}
	return true
}

// Different tests if expected and obtained are different.
func (a Assert) Different(obtained, expected interface{}, msg string) bool {
	if reflect.DeepEqual(obtained, expected) {
		return a.failFunc(Different, obtained, expected, msg)
	}
	return true
}

// Matches tests if the obtained string matches a regular expression.
func (a Assert) Matches(obtained, regex, msg string) bool {
	matches, err := regexp.MatchString("^"+regex+"$", obtained)
	if err != nil {
		return a.failFunc(Matches, obtained, regex, "can't compile regex: "+err.Error())
	}
	if !matches {
		return a.failFunc(Matches, obtained, regex, msg)
	}
	return true
}

// ErrorMatches tests if the obtained error as string matches a regular expression.
func (a Assert) ErrorMatches(obtained error, regex, msg string) bool {
	matches, err := regexp.MatchString("^"+regex+"$", obtained.Error())
	if err != nil {
		return a.failFunc(ErrorMatches, obtained, regex, "can't compile regex: "+err.Error())
	}
	if !matches {
		return a.failFunc(ErrorMatches, obtained, regex, msg)
	}
	return true
}

// Implements tests if obtained implements the expected interface variable pointer.
func (a Assert) Implements(obtained, expected interface{}, msg string) bool {
	obtainedValue := reflect.ValueOf(obtained)
	expectedValue := reflect.ValueOf(expected)
	if !obtainedValue.IsValid() {
		return a.failFunc(Implements, obtained, expected, "obtained value is invalid")
	}
	if !expectedValue.IsValid() || expectedValue.Kind() != reflect.Ptr || expectedValue.Elem().Kind() != reflect.Interface {
		return a.failFunc(Implements, obtained, expected, "expected value is no interface variable pointer")
	}
	if !obtainedValue.Type().Implements(expectedValue.Elem().Type()) {
		return a.failFunc(Implements, obtained, expected, msg)
	}
	return true
}

// Assignable tests if the types of expected and obtained are assignable.
func (a Assert) Assignable(obtained, expected interface{}, msg string) bool {
	obtainedValue := reflect.ValueOf(obtained)
	expectedValue := reflect.ValueOf(expected)
	if !obtainedValue.Type().AssignableTo(expectedValue.Type()) {
		return a.failFunc(Assignable, obtained, expected, msg)
	}
	return true
}

// Unassignable tests if the types of expected and obtained are not assignable.
func (a Assert) Unassignable(obtained, expected interface{}, msg string) bool {
	obtainedValue := reflect.ValueOf(obtained)
	expectedValue := reflect.ValueOf(expected)
	if obtainedValue.Type().AssignableTo(expectedValue.Type()) {
		return a.failFunc(Unassignable, obtained, expected, msg)
	}
	return true
}

//--------------------
// HELPER
//--------------------

// ValueDescription returns a description of a value as string.
func ValueDescription(value interface{}) string {
	rvalue := reflect.ValueOf(value)
	kind := rvalue.Kind()
	switch kind {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		return kind.String() + " of " + rvalue.Type().Elem().String()
	case reflect.Func:
		return kind.String() + " " + rvalue.Type().Name() + "()"
	case reflect.Interface, reflect.Struct:
		return kind.String() + " " + rvalue.Type().Name()
	case reflect.Ptr:
		return kind.String() + " to " + rvalue.Type().Elem().String()
	}
	// Default.
	return kind.String()
}

// isNil is a safer way to test if a value is nil.
func isNil(value interface{}) bool {
	if value == nil {
		// Standard test.
		return true
	} else {
		// Some types have to be tested via reflection.
		rvalue := reflect.ValueOf(value)
		kind := rvalue.Kind()
		switch kind {
		case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
			return rvalue.IsNil()
		}
	}
	return false
}

// EOF
