package tst_test

import (
	"testing"

	"github.com/empijei/tst"
)

func TestGolden(t *testing.T) {
	dir := t.TempDir()
	tst.WithGoldDir(t, dir)
	st := newStub(t)
	type v struct {
		Val int
	}
	tst.RecordGolden("foo", v{3}, st)
	tst.Is(fatal, st.pop().typ, t)
	got := tst.LoadGolden[v]("foo", st)
	tst.Is(v{3}, got, t)
}
