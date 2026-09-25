package tst

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

var goldDir = filepath.Join("testdata", "golden")

// RecordGolden records a golden file and makes the test fail after printing its value.
func (u *Unit) RecordGolden(name string, v any) {
	u.t.Helper()
	name = u.t.Name() + "_" + name
	s, err := os.Stat(goldDir)
	switch {
	case err == nil && !s.IsDir():
		u.t.Fatalf("%s is not a directory", goldDir)
	case errors.Is(err, os.ErrNotExist):
		u.No(os.MkdirAll(goldDir, 0o750))
	case err != nil:
		u.No(err)
	}
	buf := u.Do(json.MarshalIndent(v, "", "\t"))
	u.No(os.WriteFile(filepath.Join(goldDir, name), buf, 0o600))
	u.t.Fatalf("Gold file successfully recorded.\nOriginal data:\n%#v\n\nGolden file:\n%s", v, buf)
}

// LoadGolden loads a previously recorded golden file.
func (u *Unit) LoadGolden[V any](name string) V {
	u.t.Helper()
	var v V
	name = u.t.Name() + "_" + name
	buf := u.Do(os.ReadFile(filepath.Join(goldDir, name)))
	u.No(json.Unmarshal(buf, &v))
	return v
}

// With replaces the value pointed by val with temp, and resets it after the
// test is done running.
func (u *Unit) With[T any](val *T, temp T) {
	u.t.Helper()
	bak := *val
	*val = temp
	u.t.Cleanup(func() {
		*val = bak
	})
}
