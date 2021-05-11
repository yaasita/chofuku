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
	c, err := New(target_dir)
	if err != nil {
		t.Fatal(err)
	}
	c.UpdateHead100k()
	duplicates, err := c.GetDuplicates()
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range duplicates {
		t.Logf("size = %d, 100k = %s \n%+v", v.Size, v.Head100kHash, v.Names)
	}
}
