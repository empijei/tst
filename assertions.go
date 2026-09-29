// Package tst provides a collection of small, focused helpers designed to make Go
// tests leaner and more readable. It aims for a minimal learning curve by
// providing intuitive functions for common testing patterns like error handling,
// value unwrapping, and deep equality checks.
package tst

import (
	"strings"
	"sync/atomic"
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

type Assertions struct {
	t    Test
	used atomic.Bool
}

// Go is a shorthand for t.Parallel() and returns the unit to use assertions.
//
// Example:
//
//	a := tst.Go(t)
//	a.Is("42", "42")
func Go(t Test) *Assertions {
	t.Helper()
	if pt, ok := t.(interface{ Parallel() }); ok {
		pt.Parallel()
	}
	a := &Assertions{t: t}
	t.Cleanup(func() {
		t.Helper()
		if !a.used.Load() {
			t.Fatalf(fatalEmoji + "Assertion was created but not used.")
		}
	})
	return a
}

// Sync is like Go, but doesn't call Parallel().
func Sync(t Test) *Assertions {
	a := &Assertions{t: t}
	t.Cleanup(func() {
		t.Helper()
		if !a.used.Load() {
			t.Fatalf(fatalEmoji + "Assertion was created but not used.")
		}
	})
	return a
}

// DoB unwraps a result and stops the test immediately (t.Fatalf) if ok is false.
//
// Example:
//
//	val := a.DoB(syncMap.Load("foo"))
func (a *Assertions) DoB[V any](v V, ok bool) V {
	a.t.Helper()
	a.used.Store(true)
	if !ok {
		a.t.Fatalf(fatalEmoji + "DoB: got ok==false")
	}
	return v
}

// Be stops the test immediately (t.Fatalf) if ok is false.
//
// Example:
//
//	a.Be(len(list) > 0)
func (a *Assertions) Be(ok bool) {
	a.t.Helper()
	a.used.Store(true)
	if !ok {
		a.t.Fatalf(fatalEmoji + "Be: !ok")
	}
}

// Do unwraps a result and stops the test immediately (t.Fatalf) if an error
// occurred.
//
// Example:
//
//	f := a.Do(os.Open("file.txt"))
//	defer f.Close()
func (a *Assertions) Do[V any](v V, err error) V {
	a.t.Helper()
	a.used.Store(true)
	if err != nil {
		a.t.Fatalf(fatalEmoji+"Do: got unexpected error: %q.", err)
	}
	return v
}

// Do2 is like [Do], but for functions that return two values and an error.
//
// Example:
//
//	v1, v2 := a.Do2(returnsTwoValuesAndError())
func (a *Assertions) Do2[V1, V2 any](v1 V1, v2 V2, err error) (V1, V2) {
	a.t.Helper()
	a.used.Store(true)
	if err != nil {
		a.t.Fatalf(fatalEmoji+"Do2: got unexpected error: %q.", err)
	}
	return v1, v2
}

// No stops the test immediately (t.Fatalf) if the provided error is not nil.
//
// Example:
//
//	a.No(err)
func (a *Assertions) No(err error) {
	a.t.Helper()
	a.used.Store(true)
	if err != nil {
		a.t.Fatalf(fatalEmoji+"No: got unexpected error: %q.", err)
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
//	a.Is(want, got, t)
func (a *Assertions) Is[T any](want, got T, opts ...cmp.Option) {
	a.t.Helper()
	a.used.Store(true)
	opts = append(opts, cmpopts.EquateErrors())
	diff := cmp.Diff(want, got, opts...)
	if diff == "" {
		return
	}
	a.t.Errorf(errorEmoji+"Is: mismatch:\n\nwant:\n%#v\n\ngot:\n%#v\n\ndiff:\n%s\n", want, got, diff)
}

// IsSubString is a specialized version of Is to check that want is a substring of got.
func (a *Assertions) IsSubString(want, got string) {
	a.t.Helper()
	a.used.Store(true)
	if strings.Contains(got, want) {
		return
	}
	a.t.Errorf(errorEmoji+"IsSubString: wanted %q to be substring of %q", want, got)
}

// Err checks if the provided error is not nil and contains an optional message.
//
// It calls t.Fatalf if err is nil, and t.Errorf if the message doesn't match.
// If no message is passed, it just checks that an error occurred.
//
// Example:
//
//	a.Err("permission denied", err)
func (a *Assertions) Err(errorSubMessage string, err error) {
	a.t.Helper()
	a.used.Store(true)
	if err == nil {
		a.t.Fatalf(fatalEmoji+"Err: expected error, got %v", err)
		return
	}
	if strings.Contains(err.Error(), errorSubMessage) {
		return
	}
	a.t.Errorf(errorEmoji+"Err: want error message to contain %q, got %q", errorSubMessage, err.Error())
}

// Ko stops the test immediately (t.Fatalf) if the test has already failed.
// This is useful to prevent cascading errors if a previous check failed.
//
// Example:
//
//	a.Is(want, got)
//	a.Ko(t) // Stop here if Is failed.
func (a *Assertions) Ko() {
	a.t.Helper()
	if a.t.Failed() {
		a.t.Fatalf(stopEmoji + "Ko: Test aborted due to previous failures.")
	}
}
