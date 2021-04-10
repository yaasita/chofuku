package chofuku

import (
	"os"
	"testing"
)

func TestCompareHash(t *testing.T) {
	h1, _ := read100k("../tmp/a.dat")
	h2, _ := read_all("../tmp/a.dat")
	t.Logf("read100k = %s", h1)
	t.Logf("read_all = %s", h2)
	if h1 != h2 {
		t.Error("hash mismatch")
	}
}
func TestNew(t *testing.T) {
	target_dir := os.Getenv("TARGET_DIR")
	c, _ := New(target_dir)
	c.UpdateHead100k()
	duplicates, _ := c.GetDuplicates()
	for _, v := range duplicates {
		t.Logf("size = %d", v.Size)
	}
}
