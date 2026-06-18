package bh_lucidatypewriter_100dpi

import (
	"io/fs"
	"testing"

	"github.com/craigk5n/ilibgo"
)

// TestEmbeddedFontsParse loads every embedded .bdf and confirms it parses with
// a positive font size, guarding the go:embed migration.
func TestEmbeddedFontsParse(t *testing.T) {
	const wantCount = 14

	entries, err := fs.ReadDir(fontFS, ".")
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, e := range entries {
		name := e.Name()
		data, err := fontFS.ReadFile(name)
		if err != nil {
			t.Errorf("ReadFile(%s): %v", name, err)
			continue
		}
		f, err := ilibgo.LoadFontFromBytes(name, data)
		if err != nil {
			t.Errorf("LoadFontFromBytes(%s): %v", name, err)
			continue
		}
		if sz, err := ilibgo.GetFontSize(f); err != nil || sz <= 0 {
			t.Errorf("%s: GetFontSize = %d, err %v; want > 0", name, sz, err)
		}
		count++
	}
	if count != wantCount {
		t.Errorf("embedded font count = %d, want %d", count, wantCount)
	}
}
