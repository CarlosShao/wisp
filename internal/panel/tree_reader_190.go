package panel

import (
	"context"
	"errors"
	"os"
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
		return TreeView{Kind: TreeKindNotChosen, Reason: "no workspace was narrowed; the whole disk is not a fallback", Entries: []TreeEntry{}}
	}
	info, statErr := os.Lstat(root)
	if statErr != nil {
		if os.IsNotExist(statErr) {
			return TreeView{Kind: TreeKindUnreadable, Reason: "the root does not exist", Entries: []TreeEntry{}}
		}
		if os.IsPermission(statErr) {
			return TreeView{Kind: TreeKindPermissionDenied, Reason: "the filesystem refused this process", Entries: []TreeEntry{}}
		}
		return TreeView{Kind: TreeKindUnreadable, Reason: "the root could not be examined", Entries: []TreeEntry{}}
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return TreeView{Kind: TreeKindLinkOutside, Reason: "the root is a link; C26 default-deny refuses to follow it", Entries: []TreeEntry{}}
	}
	if !info.IsDir() {
		return TreeView{Kind: TreeKindUnreadable, Reason: "the root is not a directory", Entries: []TreeEntry{}}
	}
	if d, openErr := os.Open(root); openErr != nil {
		if os.IsPermission(openErr) {
			return TreeView{Kind: TreeKindPermissionDenied, Reason: "the filesystem refused this process", Entries: []TreeEntry{}}
		}
		return TreeView{Kind: TreeKindUnreadable, Reason: "the root could not be opened", Entries: []TreeEntry{}}
	} else {
		_ = d.Close()
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
		return TreeView{Kind: TreeKindUnreadable, Reason: "listing files failed", Entries: []TreeEntry{}}
	}
	more, _ := gitTreeRun190(ctx, root, "ls-files", "-o", "--exclude-standard", "-z")
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

// ReadFileTree190 returns the working-tree file list for root (ticket 190 host-side source).
func ReadFileTree190(root string) TreeView { return readGitTree190(root) }
