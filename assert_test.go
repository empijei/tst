package tst_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/empijei/tst"
)

const (
	fatal = "fatal"
	errr  = "error"
)

type call struct {
	typ string
	msg string
}

type stubTest struct {
	helperCalled int
	calls        []call
	clean        []func()
}

func newStub(t *testing.T) *stubTest {
	t.Helper()
	s := &stubTest{}
	t.Cleanup(func() {
		t.Helper()
		if len(s.calls) != 0 {
			t.Errorf("unchecked assertion")
		}
		if s.helperCalled == 0 {
			t.Errorf("helper not called")
		}
	})
	return s
}

func (*stubTest) Name() string {
	return "StubTest"
}

func (t *stubTest) Helper() {
	t.helperCalled++
}

func (t *stubTest) Fatalf(f string, args ...any) {
	t.calls = append(t.calls, call{fatal, fmt.Sprintf(f, args...)})
}

func (t *stubTest) Errorf(f string, args ...any) {
	t.calls = append(t.calls, call{errr, fmt.Sprintf(f, args...)})
}

func (t *stubTest) Failed() bool {
	return len(t.calls) > 0
}

func (t *stubTest) Cleanup(f func()) {
	t.clean = append(t.clean, f)
}

func (t *stubTest) pop() call {
	c := t.calls[0]
	t.calls = t.calls[1:]
	return c
}

func TestDo(t *testing.T) {
	t.Parallel()
	t.Run("ok", func(t *testing.T) {
		st := newStub(t)
		f := func() (int, error) { return 1, nil }
		v := tst.Do(f())(st)
		if v != 1 {
			t.Fatalf("bad value forwarded: want 1 got %v", v)
		}
	})
	t.Run("not ok", func(t *testing.T) {
		st := newStub(t)
		f := func() (int, error) { return 0, errors.New("argh") }
		tst.Do(f())(st)
		want := call{fatal, `❌ Do: got unexpected error: "argh".`}
		if got := st.pop(); got != want {
			t.Errorf("\nwant\n%#v\ngot\n%#v", got, want)
		}
	})
}

func TestDo2(t *testing.T) {
	t.Parallel()
	t.Run("ok", func(t *testing.T) {
		st := newStub(t)
		f := func() (int, int, error) { return 1, 2, nil }
		v1, v2 := tst.Do2(f())(st)
		if v1 != 1 || v2 != 2 {
			t.Fatalf("bad value forwarded: want (1,2) got (%v,%v)", v1, v2)
		}
	})
	t.Run("not ok", func(t *testing.T) {
		st := newStub(t)
		f := func() (int, int, error) { return 0, 0, errors.New("argh") }
		tst.Do2(f())(st)
		want := call{fatal, `❌ Do2: got unexpected error: "argh".`}
		if got := st.pop(); got != want {
			t.Errorf("\nwant\n%#v\ngot\n%#v", got, want)
		}
	})
}

func TestNo(t *testing.T) {
	t.Parallel()
	t.Run("ok", func(t *testing.T) {
		st := newStub(t)
		tst.No(nil, st)
	})
	t.Run("not ok", func(t *testing.T) {
		st := newStub(t)
		tst.No(errors.New("argh"), st)
		want := call{fatal, `❌ No: got unexpected error: "argh".`}
		if got := st.pop(); got != want {
			t.Errorf("\nwant\n%#v\ngot\n%#v", got, want)
		}
	})
}

func TestErr(t *testing.T) {
	t.Parallel()
	t.Run("ok", func(t *testing.T) {
		st := newStub(t)
		tst.Err("argh", errors.New("something argh happened"), st)
	})
	t.Run("nil error", func(t *testing.T) {
		st := newStub(t)
		tst.Err("argh", nil, st)
		want := call{fatal, `❌ Err: expected error, got <nil>`}
		if got := st.pop(); got != want {
			t.Errorf("\nwant\n%#v\ngot\n%#v", got, want)
		}
	})
	t.Run("mismatch", func(t *testing.T) {
		st := newStub(t)
		tst.Err("foo", errors.New("argh"), st)
		want := call{errr, `⚠️ Err: want error message to contain "foo", got "argh"`}
		if got := st.pop(); got != want {
			t.Errorf("\nwant\n%#v\ngot\n%#v", got, want)
		}
	})
}

func TestIs(t *testing.T) {
	t.Parallel()
	t.Run("ok", func(t *testing.T) {
		st := newStub(t)
		tst.Is(1, 1, st)
	})
	t.Run("mismatch", func(t *testing.T) {
		st := newStub(t)
		tst.Is(1, 2, st)
		got := st.pop()
		if got.typ != errr {
			t.Errorf("want err, got %v", got.typ)
		}
	})
}

func TestKo(t *testing.T) {
	t.Parallel()
	t.Run("ok", func(t *testing.T) {
		st := newStub(t)
		tst.Ko(st)
	})
	t.Run("not ok", func(t *testing.T) {
		st := newStub(t)
		st.Errorf("previous failure")
		tst.Ko(st)
		_ = st.pop()
		got := st.pop()
		want := call{fatal, "⛔ Ko: Test aborted due to previous failures."}
		if got != want {
			t.Errorf("call[1]:\nwant\n%#v\ngot\n%#v", got, want)
		}
	})
}
