// FourCoreLabs fork patch; see FORK.md.

package interp

import "testing"

// TestAbsPathVFS pins POSIX resolution under VFSPaths. Without it, on Windows
// filepath.IsAbs("/work") is false and a virtual absolute path was joined
// onto the current directory.
func TestAbsPathVFS(t *testing.T) {
	tests := []struct{ dir, path, want string }{
		{"/work", "a/b", "/work/a/b"},
		{"/work", "/shared/x", "/shared/x"},
		{"/work/sub", "../x", "/work/x"},
		{"/work", "./a//b/", "/work/a/b"},
		{"/", "..", "/"},
		{"/work", "", ""},
	}
	for _, tt := range tests {
		if got := absPath(tt.dir, tt.path, true); got != tt.want {
			t.Errorf("absPath(%q, %q, vfs) = %q, want %q", tt.dir, tt.path, got, tt.want)
		}
	}
}

// TestVFSPathsSurvivesSubshellAndReset: the option must reach subshells
// (pipelines, $(...)) and survive Reset, or those paths fall back to host
// semantics.
func TestVFSPathsSurvivesSubshellAndReset(t *testing.T) {
	r, err := New(VFSPaths())
	if err != nil {
		t.Fatal(err)
	}
	if !r.subshell(false).vfsPaths {
		t.Error("subshell lost vfsPaths")
	}
	r.Reset()
	if !r.vfsPaths {
		t.Error("Reset lost vfsPaths")
	}
}
