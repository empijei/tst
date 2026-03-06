package tst

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

var goldDir = filepath.Join("testdata", "golden")

// RecordGolden records a golden file and makes the test fail after printing its value.
func RecordGolden(name string, v any, t Test) {
	t.Helper()
	name = t.Name() + "_" + name
	s, err := os.Stat(goldDir)
	switch {
	case err == nil && !s.IsDir():
		t.Fatalf("%s is not a directory", goldDir)
	case errors.Is(err, os.ErrNotExist):
		No(os.MkdirAll(goldDir, 0o750), t)
	case err != nil:
		No(err, t)
	}
	buf := Do(json.MarshalIndent(v, "", "\t"))(t)
	No(os.WriteFile(filepath.Join(goldDir, name), buf, 0o600), t)
	t.Fatalf("Gold file successfully recorded.\nOriginal data:\n%#v\n\nGolden file:\n%s", v, buf)
}

// LoadGolden loads a previously recorded golden file.
func LoadGolden[V any](name string, t Test) V {
	t.Helper()
	var v V
	name = t.Name() + "_" + name
	buf := Do(os.ReadFile(filepath.Join(goldDir, name)))(t)
	No(json.Unmarshal(buf, &v), t)
	return v
}

// With replaces the value pointed by val with temp, and resets it after the
// test is done running.
func With[T any](val *T, temp T, t Test) {
	t.Helper()
	bak := *val
	*val = temp
	t.Cleanup(func() {
		*val = bak
	})
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
