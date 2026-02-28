// Package tst is a simple package to write lean and readable tests that has a very
// low learning curve.
package tst

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

const (
	fatalEmoji = "❌ "
	stopEmoji  = "⛔ "
	errorEmoji = "⚠️ "
	bugEmoji   = "🚨 "
)

var (
	_ Test = &testing.T{}
	_ Test = &testing.F{}
	_ Test = &testing.B{}
)

// Test is a test case.
type Test interface {
	Helper()
	Fatalf(string, ...any)
	Errorf(string, ...any)
	Failed() bool
	Cleanup(func())
}

func checkFormat(t Test, msgArgs []any) string {
	t.Helper()

	if len(msgArgs) == 0 {
		return ""
	}

	format, ok := msgArgs[0].(string)
	if !ok {
		t.Fatalf(bugEmoji+"Bug in tests detected:\nFirst element of msgArgs must be format string, but type is %T.", msgArgs[0])
	}
	return format
}

func format(t Test, msgArgs []any) string {
	t.Helper()

	if len(msgArgs) == 0 {
		return ""
	}

	format := checkFormat(t, msgArgs)
	return fmt.Sprintf(format, msgArgs[1:]...) + ": "
}

// Do unwraps a result and stops the test if an error occurred.
func Do[V any](v V, err error) func(t Test, msgArgs ...any) V {
	// TODO
	return func(t Test, msgArgs ...any) V {
		t.Helper()
		checkFormat(t, msgArgs)

		if err != nil {
			t.Fatalf(fatalEmoji+"Do: %sgot unexpected error: %q.", format(t, msgArgs), err)
		}
		return v
	}
}

// Do2 is like [Do], but for functions that return 2 values and an error.
func Do2[V1, V2 any](v1 V1, v2 V2, err error) func(t Test, msgArgs ...any) (V1, V2) {
	// TODO
	return func(t Test, msgArgs ...any) (V1, V2) {
		t.Helper()
		checkFormat(t, msgArgs)

		if err != nil {
			t.Fatalf(fatalEmoji+"Do2: %sgot unexpected error: %q.", format(t, msgArgs), err)
		}
		return v1, v2
	}
}

// No stops the test if an error occurred.
func No(err error) func(t Test, msgArgs ...any) {
	// TODO
	return func(t Test, msgArgs ...any) {
		t.Helper()
		checkFormat(t, msgArgs)
		if err != nil {
			t.Fatalf(fatalEmoji+"No: %sgot unexpected error: %q.", format(t, msgArgs), err)
		}
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
func Is[T any](want, got T, opts ...cmp.Option) func(t Test, msgArgs ...any) {
	// TODO
	return func(t Test, msgArgs ...any) {
		t.Helper()
		checkFormat(t, msgArgs)
		opts = append(opts, cmpopts.EquateErrors())
		diff := cmp.Diff(want, got, opts...)
		if diff == "" {
			return
		}
		usrMsg := format(t, msgArgs)
		if usrMsg == "" {
			usrMsg = "mismatch:"
		}
		t.Errorf(errorEmoji+"Is: %s\n\nwant:\n%#v\n\ngot:\n%#v\n\ndiff:\n%s\n", usrMsg, want, got, diff)
	}
}

// Err checks if the provided error is not nil and contains an optional message.
func Err(errorSubMessage string, err error) func(t Test, msgArgs ...any) {
	// TODO
	return func(t Test, msgArgs ...any) {
		t.Helper()
		checkFormat(t, msgArgs)
		if err == nil {
			t.Fatalf(fatalEmoji+"Err: expected error, got %v", err)
			return
		}
		if strings.Contains(err.Error(), errorSubMessage) {
			return
		}
		t.Errorf(errorEmoji+"Err: %swant error message to contain %q, got %q", format(t, msgArgs), errorSubMessage, err.Error())
	}
}
