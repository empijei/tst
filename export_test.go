package tst

func (u *Assertions) WithGoldDir(t Test, path string) {
	t.Helper()
	u.With(&goldDir, path)
}
