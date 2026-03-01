// Package tst is a simple package to write lean and readable tests that has a very
// low learning curve.
//
// THIS API IS CLUNKY: we need to wait for generic methods to properly do this.
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

// Test is an abstraction over *testing.(T|B|F).
type Test interface {
	Helper()
	Fatalf(string, ...any)
	Errorf(string, ...any)
	Failed() bool
	Cleanup(func())
}

// Do unwraps a result and stops the test if an error occurred.
func Do[V any](v V, err error) func(t Test) V {
	return func(t Test) V {
		t.Helper()

		if err != nil {
			t.Fatalf(fatalEmoji+"Do: got unexpected error: %q.", err)
		}
		return v
	}
}

// Do2 is like [Do], but for functions that return 2 values and an error.
func Do2[V1, V2 any](v1 V1, v2 V2, err error) func(t Test) (V1, V2) {
	return func(t Test) (V1, V2) {
		t.Helper()

		if err != nil {
			t.Fatalf(fatalEmoji+"Do2: got unexpected error: %q.", err)
		}
		return v1, v2
	}
}

// No stops the test if an error occurred.
func No(err error, t Test) {
	t.Helper()
	if err != nil {
		t.Fatalf(fatalEmoji+"No: got unexpected error: %q.", err)
	}
}

// Ko stops the test if it has already failed.
func Ko(t Test) {
	t.Helper()
	if t.Failed() {
		t.Fatalf(stopEmoji + "Ko: Test aborted due to previous failures.")
	}
}

// Is checks that want matches got via [cmp.Diff] using the options provided.
// Errors are compared with [errors.Is] by default.
//
// Options can be found both in [cmp] and [cmpopts].
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

// PTest is an abstraction over [*testing.T].
type PTest interface {
	Test
	Context() context.Context
	Parallel()
}

var _ PTest = &testing.T{}

// Go is a shorthand for t.Parallel(); t.Context().
func Go(t PTest) context.Context {
	t.Helper()
	return t.Context()
}
