// Package tst provides a collection of small, focused helpers designed to make Go
// tests leaner and more readable. It aims for a minimal learning curve by
// providing intuitive functions for common testing patterns like error handling,
// value unwrapping, and deep equality checks.
package tst

import (
	"context"
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
	Fatalf(string, ...any)
	Errorf(string, ...any)
	Failed() bool
	Cleanup(func())
}

// DoB unwraps a result and stops the test immediately (t.Fatalf) if ok is false.
//
// Example:
//
//	val := tst.DoB(syncMap.Load("foo"))(t)
func DoB[V any](v V, ok bool) func(t Test) V {
	return func(t Test) V {
		t.Helper()

		if !ok {
			t.Fatalf(fatalEmoji + "DoB: got ok==false")
		}
		return v
	}
}

// Be stops the test immediately (t.Fatalf) if ok is false.
//
// Example:
//
//	tst.Be(len(list) > 0, t)
func Be(ok bool, t Test) {
	t.Helper()
	if !ok {
		t.Fatalf(fatalEmoji + "Be: !ok")
	}
}

// Do unwraps a result and stops the test immediately (t.Fatalf) if an error
// occurred.
//
// Example:
//
//	f := tst.Do(os.Open("file.txt"))(t)
//	defer f.Close()
func Do[V any](v V, err error) func(t Test) V {
	return func(t Test) V {
		t.Helper()

		if err != nil {
			t.Fatalf(fatalEmoji+"Do: got unexpected error: %q.", err)
		}
		return v
	}
}

// Do2 is like [Do], but for functions that return two values and an error.
//
// Example:
//
//	v1, v2 := tst.Do2(returnsTwoValuesAndError())(t)
func Do2[V1, V2 any](v1 V1, v2 V2, err error) func(t Test) (V1, V2) {
	return func(t Test) (V1, V2) {
		t.Helper()

		if err != nil {
			t.Fatalf(fatalEmoji+"Do2: got unexpected error: %q.", err)
		}
		return v1, v2
	}
}

// No stops the test immediately (t.Fatalf) if the provided error is not nil.
//
// Example:
//
//	tst.No(err, t)
func No(err error, t Test) {
	t.Helper()
	if err != nil {
		t.Fatalf(fatalEmoji+"No: got unexpected error: %q.", err)
	}
}

// Ko stops the test immediately (t.Fatalf) if the test has already failed.
// This is useful to prevent cascading errors if a previous check failed.
//
// Example:
//
//	tst.Is(want, got, t)
//	tst.Ko(t) // Stop here if Is failed.
func Ko(t Test) {
	t.Helper()
	if t.Failed() {
		t.Fatalf(stopEmoji + "Ko: Test aborted due to previous failures.")
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
//	tst.Is(want, got, t)
func Is[T any](want, got T, t Test, opts ...cmp.Option) {
	t.Helper()
	opts = append(opts, cmpopts.EquateErrors())
	diff := cmp.Diff(want, got, opts...)
	if diff == "" {
		return
	}
	t.Errorf(errorEmoji+"Is: mismatch:\n\nwant:\n%#v\n\ngot:\n%#v\n\ndiff:\n%s\n", want, got, diff)
}

// Err checks if the provided error is not nil and contains an optional message.
//
// It calls t.Fatalf if err is nil, and t.Errorf if the message doesn't match.
// If no message is passed, it just checks that an error occurred.
//
// Example:
//
//	tst.Err("permission denied", err, t)
func Err(errorSubMessage string, err error, t Test) {
	t.Helper()
	if err == nil {
		t.Fatalf(fatalEmoji+"Err: expected error, got %v", err)
		return
	}
	if strings.Contains(err.Error(), errorSubMessage) {
		return
	}
	t.Errorf(errorEmoji+"Err: want error message to contain %q, got %q", errorSubMessage, err.Error())
}

// PTest is an abstraction over [*testing.T] that includes Parallel and Context.
type PTest interface {
	Test
	Context() context.Context
	Parallel()
}

var _ PTest = &testing.T{}

// Go is a shorthand for t.Parallel() and returns the test context.
//
// Example:
//
//	ctx := tst.Go(t)
func Go(t PTest) context.Context {
	t.Helper()
	t.Parallel()
	return t.Context()
}
