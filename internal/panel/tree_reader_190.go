package panel

import (
	"context"
	"errors"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const treeGitTimeout190 = 8 * time.Second

func splitNUL190(raw []byte) []string {
	var out []string
	for _, part := range strings.Split(string(raw), "\x00") {
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

func gitTreeRun190(ctx context.Context, root string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "git", append([]string{"-C", root}, args...)...)
	return cmd.Output()
}

// readGitTree190 returns the file tree under root using git, or a third state.
func readGitTree190(root string) TreeView {
	if strings.TrimSpace(root) == "" {
		return TreeView{Kind: TreeKindNotRepo, Reason: "no workspace to read", Entries: []TreeEntry{}}
	}
	ctx, cancel := context.WithTimeout(context.Background(), treeGitTimeout190)
	defer cancel()
	out, err := gitTreeRun190(ctx, root, "rev-parse", "--is-inside-work-tree")
	if err != nil {
		var exErr *exec.Error
		if errors.As(err, &exErr) {
			return TreeView{Kind: TreeKindNoGit, Reason: "no usable git on this machine", Entries: []TreeEntry{}}
		}
		return TreeView{Kind: TreeKindNotRepo, Reason: "not a git work tree", Entries: []TreeEntry{}}
	}
	if strings.TrimSpace(string(out)) != "true" {
		return TreeView{Kind: TreeKindNotRepo, Reason: "not a git work tree", Entries: []TreeEntry{}}
	}
	files, err := gitTreeRun190(ctx, root, "ls-files", "-z")
	if err != nil {
		return TreeView{Kind: TreeKindTooLarge, Reason: "listing files failed", Entries: []TreeEntry{}}
	}
	more, _ := gitTreeRun190(ctx, root, "ls-files", "-o", "-i", "--exclude-standard", "-z")
	paths := append(splitNUL190(files), splitNUL190(more)...)
	const cap190 = 2000
	truncated := false
	total := 0
	entries := make([]TreeEntry, 0, len(paths))
	for _, p := range paths {
		if p == ".git" || strings.HasPrefix(p, ".git/") {
			continue
		}
		total++
		if len(entries) >= cap190 {
			truncated = true
			continue
		}
		entries = append(entries, TreeEntry{Path: p, IsDir: false})
	}
	reason := "total " + strconv.Itoa(total) + " files"
	return TreeView{Kind: TreeKindCounted, Reason: reason, Entries: entries, FilesTotal: total, FilesTruncated: truncated}
}
