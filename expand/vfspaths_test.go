// FourCoreLabs fork patch; see FORK.md.

package expand_test

import (
	"io/fs"
	"slices"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/FourCoreLabs/sh/expand"
	"github.com/FourCoreLabs/sh/syntax"
)

// TestGlobVFSPaths globs over a virtual filesystem whose paths are POSIX.
// Without VFSPaths, on Windows the base "/work" is not absolute and is joined
// with "\", so ReadDir2 is asked for "/work/\work" and nothing matches.
func TestGlobVFSPaths(t *testing.T) {
	vfs := fstest.MapFS{
		"work/a.txt":     {},
		"work/b.txt":     {},
		"work/c.log":     {},
		"work/sub/d.txt": {},
		"shared/e.txt":   {},
	}
	var asked []string
	readDir := func(p string) ([]fs.DirEntry, error) {
		asked = append(asked, p)
		if strings.Contains(p, `\`) || !strings.HasPrefix(p, "/") {
			t.Errorf("ReadDir2(%q): want an absolute POSIX path", p)
		}
		name := strings.TrimPrefix(p, "/")
		if name == "" {
			name = "."
		}
		return fs.ReadDir(vfs, name)
	}
	tests := []struct {
		pat  string
		want []string
	}{
		{"*.txt", []string{"a.txt", "b.txt"}},
		{"[ab].txt", []string{"a.txt", "b.txt"}},
		{"sub/*.txt", []string{"sub/d.txt"}},
		{"/work/*.txt", []string{"/work/a.txt", "/work/b.txt"}},
		{"/shared/*", []string{"/shared/e.txt"}},
		{"../shared/*.txt", []string{"../shared/e.txt"}},
		{"**/*.txt", []string{"a.txt", "b.txt", "sub/d.txt"}},
	}
	for _, tt := range tests {
		t.Run(tt.pat, func(t *testing.T) {
			asked = asked[:0]
			cfg := &expand.Config{
				Env:      expand.ListEnviron("PWD=/work"),
				ReadDir2: readDir,
				VFSPaths: true,
				GlobStar: true,
			}
			word, err := syntax.NewParser().Document(strings.NewReader(tt.pat))
			if err != nil {
				t.Fatal(err)
			}
			got, err := expand.Fields(cfg, word)
			if err != nil {
				t.Fatalf("Fields: %v (ReadDir2 calls %q)", err, asked)
			}
			slices.Sort(got)
			if !slices.Equal(got, tt.want) {
				t.Errorf("glob %q = %q, want %q (ReadDir2 calls %q)", tt.pat, got, tt.want, asked)
			}
		})
	}
}
