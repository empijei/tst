// Package tst provides a collection of small, focused helpers designed to make Go
// tests leaner and more readable. It aims for a minimal learning curve by
// providing intuitive functions for common testing patterns like error handling,
// value unwrapping, and deep equality checks.
package tst

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

const (
	fatalEmoji = "❌ "
	stopEmoji  = "⛔ "
	errorEmoji = "⚠️ "
)

var (
	_ Test = &testing.T{}
	_ Test = &testing.F{}
	_ Test = &testing.B{}
)

// Test is an abstraction over *testing.(T|B|F). It allows tst helpers to work
// with different testing types.
type Test interface {
	Helper()
	Name() string
	Fatalf(string, ...any)
	Errorf(string, ...any)
	Failed() bool
	Cleanup(func())
}

// PTest is an abstraction over [*testing.T] that includes Parallel and Context.
type PTest interface {
	Test
	Parallel()
}

var _ PTest = &testing.T{}

// Sync allows to call Go without triggering a parallel execution.
//
// Example:
//
//	u := tst.Go(tst.Sync(t))
func Sync(pt PTest) (t Test) { return t }

type Unit struct {
	t Test
}

// Go is a shorthand for t.Parallel() and returns the unit to use assertions.
//
// Example:
//
//	u := tst.Go(t)
//	u.Is("42", "42")
func Go(t Test) *Unit {
	t.Helper()
	if pt, ok := t.(PTest); ok {
		pt.Parallel()
	}
	return &Unit{t}
}

// DoB unwraps a result and stops the test immediately (t.Fatalf) if ok is false.
//
// Example:
//
//	val := u.DoB(syncMap.Load("foo"))
func (u *Unit) DoB[V any](v V, ok bool) V {
	u.t.Helper()
	if !ok {
		u.t.Fatalf(fatalEmoji + "DoB: got ok==false")
	}
	return v
}

// Be stops the test immediately (t.Fatalf) if ok is false.
//
// Example:
//
//	u.Be(len(list) > 0)
func (u *Unit) Be(ok bool) {
	u.t.Helper()
	if !ok {
		u.t.Fatalf(fatalEmoji + "Be: !ok")
	}
}

// Do unwraps a result and stops the test immediately (t.Fatalf) if an error
// occurred.
//
// Example:
//
//	f := u.Do(os.Open("file.txt"))
//	defer f.Close()
func (u *Unit) Do[V any](v V, err error) V {
	u.t.Helper()
	if err != nil {
		u.t.Fatalf(fatalEmoji+"Do: got unexpected error: %q.", err)
	}
	return v
}

// Do2 is like [Do], but for functions that return two values and an error.
//
// Example:
//
//	v1, v2 := u.Do2(returnsTwoValuesAndError())
func (u *Unit) Do2[V1, V2 any](v1 V1, v2 V2, err error) (V1, V2) {
	u.t.Helper()
	if err != nil {
		u.t.Fatalf(fatalEmoji+"Do2: got unexpected error: %q.", err)
	}
	return v1, v2
}

// No stops the test immediately (t.Fatalf) if the provided error is not nil.
//
// Example:
//
//	u.No(err)
func (u *Unit) No(err error) {
	u.t.Helper()
	if err != nil {
		u.t.Fatalf(fatalEmoji+"No: got unexpected error: %q.", err)
	}
}

// Is checks that want matches got via [cmp.Diff] using the options provided.
// It calls t.Errorf if there's a mismatch.
// Errors are compared with [cmpopts.EquateErrors] by default.
//
// Options can be found in both [cmp] and [cmpopts] packages.
//
// Example:
//
//	u.Is(want, got, t)
func (u *Unit) Is[T any](want, got T, opts ...cmp.Option) {
	u.t.Helper()
	opts = append(opts, cmpopts.EquateErrors())
	diff := cmp.Diff(want, got, opts...)
	if diff == "" {
		return
	}
	u.t.Errorf(errorEmoji+"Is: mismatch:\n\nwant:\n%#v\n\ngot:\n%#v\n\ndiff:\n%s\n", want, got, diff)
}

// IsSubString is a specialized version of Is to check that want is a substring of got.
func (u *Unit) IsSubString(want, got string) {
	u.t.Helper()
	if strings.Contains(got, want) {
		return
	}
	u.t.Errorf(errorEmoji+"IsSubString: wanted %q to be substring of %q", want, got)
}

// Err checks if the provided error is not nil and contains an optional message.
//
// It calls t.Fatalf if err is nil, and t.Errorf if the message doesn't match.
// If no message is passed, it just checks that an error occurred.
//
// Example:
//
//	u.Err("permission denied", err)
func (u *Unit) Err(errorSubMessage string, err error) {
	u.t.Helper()
	if err == nil {
		u.t.Fatalf(fatalEmoji+"Err: expected error, got %v", err)
		return
	}
	if strings.Contains(err.Error(), errorSubMessage) {
		return
	}
	u.t.Errorf(errorEmoji+"Err: want error message to contain %q, got %q", errorSubMessage, err.Error())
}

// Ko stops the test immediately (t.Fatalf) if the test has already failed.
// This is useful to prevent cascading errors if a previous check failed.
//
// Example:
//
//	u.Is(want, got)
//	u.Ko(t) // Stop here if Is failed.
func (u *Unit) Ko() {
	u.t.Helper()
	if u.t.Failed() {
		u.t.Fatalf(stopEmoji + "Ko: Test aborted due to previous failures.")
	}
}
