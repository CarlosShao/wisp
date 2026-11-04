package panel

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestTree190RootDoesNotExistUnreadable pins AC#4(a): a root that is not there is
// reported as unreadable, not as an empty tree.
func TestTree190RootDoesNotExistUnreadable(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "does-not-exist")
	v := readGitTree190(missing)
	if v.Kind != TreeKindUnreadable {
		t.Fatalf("missing root kind=%q want unreadable", v.Kind)
	}
	if v.Entries == nil {
		t.Fatalf("entries must be a non-nil empty array, got nil")
	}
}

// TestTree190LinkOutsideRefused pins AC#4(c): a root that is itself a symlink is
// refused (C26 default-deny) rather than followed and listed.
func TestTree190LinkOutsideRefused(t *testing.T) {
	base := t.TempDir()
	target := filepath.Join(base, "target")
	if err := os.Mkdir(target, 0o755); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(base, "link")
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("cannot create symlink on this host: %v", err)
	}
	v := readGitTree190(link)
	if v.Kind != TreeKindLinkOutside {
		t.Fatalf("symlink root kind=%q want link_outside", v.Kind)
	}
}

// TestTree190ReadSurfaceHasNoWriteCalls pins AC#3: this dimension is read only.
// It scans the reader's own source for any write-shaped os call, so a later leg
// that slips a Remove/Mkdir/WriteFile into the read surface goes red here.
func TestTree190ReadSurfaceHasNoWriteCalls(t *testing.T) {
	src, err := os.ReadFile("tree_reader_190.go")
	if err != nil {
		t.Fatalf("read own source: %v", err)
	}
	forbidden := []string{
		"os.Remove", "os.Rename", "os.Mkdir", "os.WriteFile",
		"os.Create", "os.RemoveAll", "os.Chmod", "os.Symlink",
	}
	for _, f := range forbidden {
		if strings.Contains(string(src), f) {
			t.Fatalf("read-only surface calls a write: %q present in tree_reader_190.go", f)
		}
	}
}
