package tst

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

var goldDir = filepath.Join("testdata", "golden")

// RecordGolden records a golden file and makes the test fail after printing its value.
func (a *Assertions) RecordGolden(name string, v any) {
	a.t.Helper()
	name = a.t.Name() + "_" + name
	s, err := os.Stat(goldDir)
	switch {
	case err == nil && !s.IsDir():
		a.t.Fatalf("%s is not a directory", goldDir)
	case errors.Is(err, os.ErrNotExist):
		a.No(os.MkdirAll(goldDir, 0o750))
	case err != nil:
		a.No(err)
	}
	buf := a.Do(json.MarshalIndent(v, "", "\t"))
	a.No(os.WriteFile(filepath.Join(goldDir, name), buf, 0o600))
	a.t.Fatalf("Gold file successfully recorded.\nOriginal data:\n%#v\n\nGolden file:\n%s", v, buf)
}

// LoadGolden loads a previously recorded golden file.
func (a *Assertions) LoadGolden[V any](name string) V {
	a.t.Helper()
	var v V
	name = a.t.Name() + "_" + name
	buf := a.Do(os.ReadFile(filepath.Join(goldDir, name)))
	a.No(json.Unmarshal(buf, &v))
	return v
}

// With replaces the value pointed by val with temp, and resets it after the
// test is done running.
func (a *Assertions) With[T any](val *T, temp T) {
	a.t.Helper()
	bak := *val
	*val = temp
	a.t.Cleanup(func() {
		*val = bak
	})
}
