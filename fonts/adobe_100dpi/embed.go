package adobe_100dpi

import (
	"embed"
	"strings"
)

//go:embed *.bdf
var fontFS embed.FS

// lines returns the contents of an embedded BDF file split into lines.
func lines(name string) []string {
	b, err := fontFS.ReadFile(name)
	if err != nil {
		return nil
	}
	return strings.Split(strings.TrimRight(string(b), "\n"), "\n")
}
