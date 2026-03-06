package tst

func WithGoldDir(t Test, path string) {
	t.Helper()
	With(&goldDir, path, t)
}
