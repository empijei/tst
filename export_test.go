package tst

func (u *Unit) WithGoldDir(t Test, path string) {
	t.Helper()
	u.With(&goldDir, path)
}
