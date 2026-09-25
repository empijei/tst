package tst_test

import (
	"testing"

	"github.com/empijei/tst"
)

func TestGolden(t *testing.T) {
	st := newStub(t)
	u := tst.Go(st)
	dir := t.TempDir()
	u.WithGoldDir(t, dir)
	type v struct {
		Val int
	}
	u.RecordGolden("foo", v{3})
	u.Is(fatal, st.pop().typ)
	got := u.LoadGolden[v]("foo")
	u.Is(v{3}, got)
}
