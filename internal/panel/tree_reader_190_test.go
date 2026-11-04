package panel

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadGitTree190NotARepo(t *testing.T) {
	dir := t.TempDir()
	v := readGitTree190(dir)
	if v.Kind != TreeKindNotRepo {
		t.Fatalf("plain temp dir: kind=%q want not_repo", v.Kind)
	}
	if v.Entries == nil {
		t.Fatalf("entries must be non-nil empty array, got nil")
	}
	b, _ := json.Marshal(v)
	if !strings.Contains(string(b), "\"entries\":[]") {
		t.Fatalf("entries must serialize as [], got %s", b)
	}
}

func TestReadGitTree190EmptyRoot(t *testing.T) {
	v := readGitTree190("")
	if v.Kind != TreeKindNotChosen {
		t.Fatalf("empty root kind=%q want not_chosen (the whole disk is not a fallback)", v.Kind)
	}
}

func TestReadGitTree190RealRepo(t *testing.T) {
	dir := t.TempDir()
	if out, err := exec.Command("git", "-C", dir, "init").CombinedOutput(); err != nil {
		t.Fatalf("git init: %v %s", err, out)
	}
	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	v := readGitTree190(dir)
	if v.Kind != TreeKindCounted {
		t.Fatalf("real repo kind=%q reason=%q", v.Kind, v.Reason)
	}
	if v.FilesTotal < 1 {
		t.Fatalf("want >=1 file, got %d", v.FilesTotal)
	}
	found := false
	for _, e := range v.Entries {
		if e.Path == "a.txt" {
			found = true
		}
	}
	if !found {
		t.Fatalf("a.txt not in entries: %+v", v.Entries)
	}
}
