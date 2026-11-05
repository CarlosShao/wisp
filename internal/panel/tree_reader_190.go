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

// ReadFileTreeForWorkspace190 is the account-aware entry to ticket 190's file
// list: the same shape the git dimension got in ReadGitForWorkspace (git.go:298),
// for the same rule (ticket 102, ledger A380).
//
// WHY THIS EXISTS AT ALL. A file list is read off a coordinate, and a coordinate
// is only "the folder the user named" if C26's account says so. ReadFileTree190
// takes a bare string and cannot know whether expansion moved that string onto a
// different tree, so a leg that reaches past the account for the path renders a
// listing of a tree nobody named as if it were the operator's own folder - the
// exact shape the account exists to prevent. Taking the whole WorkspaceView is
// what makes the account readable at this leg instead of a field the caller
// silently skips.
//
// THE DISPOSITION IS THE GIT DIMENSION'S, deliberately. GitViewRewritten
// (git.go:335) reports candidate ① of ticket 181 AC#6 as ruled by the
// orchestrator: state the expanded truth AND say out loud that the coordinate
// was rewritten. The tree rides the same ruling so one packet cannot narrate the
// same workspace two different ways - git's section saying "rewritten" while the
// files section quietly says "your folder". Refusing the whole dimension
// (candidate ②) is not implemented here either, for the same reason it is not
// implemented there: that ruling is the orchestrator's, and this leg does not
// take it by itself.
//
// A view that was never rewritten is the ordinary road and behaves exactly as
// before, byte for byte - the same contract ReadGitForWorkspace carries.
func ReadFileTreeForWorkspace190(ws WorkspaceView) TreeView {
	if ws.Rewritten {
		return fileTreeRewritten190(ws)
	}
	return ReadFileTree190(ws.Canonical)
}

// fileTreeRewritten190 is ticket 190's file list for a coordinate C26's expansion
// step moved onto another tree. It reads the expanded tree (the account-bearing
// reading) and puts the account into the sentence the page renders, in the same
// words the git section uses: both the spelling the user gave and the tree it
// moved to, plus the literal `rewritten=true` a reader can grep for.
func fileTreeRewritten190(ws WorkspaceView) TreeView {
	view := ReadFileTree190(ws.Canonical)
	view.Reason = "注意：这条工作区路径被 C26 的展开步改写过（账户 rewritten=true），" +
		"下面这份文件清单来自改写后的树 " + ws.Canonical + "，不是原拼法 " + ws.Spelling + "。" + view.Reason
	return view
}
