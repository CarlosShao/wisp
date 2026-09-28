package panel

// The git dimension of the snapshot (ticket 181: the read-only half).
//
// WHAT THIS FILE IS. A host-side read surface that answers three questions by
// reading files and nothing else: is this workspace a git tree, which branch (or
// which bare commit) is checked out in it, and which worktrees does it have. The
// answers ride the snapshot push that already exists (internal/panel/pump.go ->
// cmd/wisp/panel_pump.go), which is the outbound half ticket 181's census
// (docs/evidence/s1/181-186-git-detection-census-c1.md §⑤) measured as really
// running today.
//
// WHAT THIS FILE IS NOT. It is not a tool. D34's built-in tool table
// (docs/PLAN.md §D34) carries no git row, and giving the model a `git.*` tool
// would be a contract change, not a display feature - so nothing here is
// registered in internal/tools and nothing here is reachable from a model call.
// TestGitDimensionHasNoModelCallableTool nails that. It is also not the inbound
// half: "switch branch / switch tree" needs the page -> Go hop (a postMessage
// receiver plus a router), and that hop does not exist in this tree, which is why
// SwitchBlocked below is a constant sentence rather than a computed verdict.
//
// NO EXTERNAL PROCESS. Everything is os.Stat / os.ReadFile / os.ReadDir. The
// census §② took this repository's own three shapes as live data and concluded
// that recognising a repo, the current branch and detached HEAD needs no `git`
// binary; spawning one would additionally drag in D38's goroutine-owner rules and
// d22scan ban #1 for no gain.
//
// THE THREE SHAPES THIS HAS TO EAT (census §②, all measured on real data):
//
//	ordinary repo   <dir>/.git is a DIRECTORY; HEAD lives inside it.
//	linked tree <dir>/.git is a FILE whose content is "gitdir: <path>" -
//	                absolute or relative; the path is the per-tree private
//	                dir, and the shared dir (refs, packed-refs, worktrees/) comes
//	                from the "commondir" file sitting next to HEAD.
//	detached HEAD   the HEAD file is a 40-hex sha with no "ref: " prefix.
//
// A "gitdir:" that points at a directory which is gone is UNREADABLE, never
// "not a repo": the .git file is right there, so the workspace was a tree and
// the tree it named was deleted underneath it. Collapsing those two would be the
// same lie ticket 92/145 refuse in every other section of this packet.

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// ONE MORE THING THIS FILE DOES NOT DO: it makes no path AUTHORISATION decision.
// The probe walks up from the directory the operator already pointed the host at,
// and the two pointer files it follows (.git's "gitdir:" and the per-tree
// "commondir") are read as the metadata they are - the strings inside them come
// from the repository, so they are used only to locate more metadata to read,
// never as a scope, never written to, never executed. C26's PathResolver keeps the
// only vote on what any TOOL may touch (AGENTS §1.2 names it, and d22scan ban #2
// enforces it); this section reports a fact about the host's own tree and reaches
// the user through the snapshot, which is the other half of why it is not a Tool.

// The four states of GitView.kind. An enum, not a boolean: "this is not a repo",
// "I am not allowed to look" and "I could not read it" are three different facts,
// and the panel must not render the second and third as the first (census §⑦).
const (
	GitKindRepo             = "repo"
	GitKindNotARepo         = "not_a_repo"
	GitKindPermissionDenied = "permission_denied"
	GitKindUnreadable       = "unreadable"
)

// GitSwitchBlockedReason is the snapshot's standing statement about the two
// actions this ticket does NOT implement. It is not a policy call and it is not
// a refusal to be worked around: the inbound channel that a switch would arrive
// on (the WebView2 host of ticket 33, its postMessage receiver, and a router for
// ComposerRequest) has zero production callers in this tree - measured, with the
// grep lines, in census §⑤. A panel that drew a working switcher on top of that
// would be a button wired to nothing.
//
// The day the inbound hop lands, this constant stops being true and ticket 186
// owns replacing it - not this file.
const GitSwitchBlockedReason = "切换分支／切换工作树今天不可用：网页到宿主的那一跳还没有落地（无 postMessage 接收器、无 router、无 WebView2 宿主），面板只有快照这一条出向通道"

// gitNoWorkspaceReason is the honest answer when nothing was narrowed. It is a
// separate sentence from "not a repo" because there was no directory to look at.
const gitNoWorkspaceReason = "本轮未选择工作区，git 这一维没有可探测的目录"

// GitWorktree is one entry of the tree list. The list ALWAYS includes the
// main tree: git keeps the main tree's own HEAD at <common>/HEAD and lists
// only the attached ones under <common>/worktrees/, so an implementation that
// enumerates the directory alone drops the tree the user is standing in (census
// §③ measured this repository: ls .git/worktrees = 2, git's own listing = 3).
type GitWorktree struct {
	// Path is that tree's directory ("" when its gitdir pointer cannot be
	// read - a gap, never a guess).
	Path string `json:"path"`
	// Branch is the checked-out branch ("" when detached).
	Branch string `json:"branch"`
	// Sha is the bare HEAD when detached ("" otherwise).
	Sha string `json:"sha"`
	// Detached says Branch is empty because HEAD is a raw commit.
	Detached bool `json:"detached"`
	// Main marks the primary tree (the one whose .git is a directory, or the
	// one the shared git dir belongs to).
	Main bool `json:"main"`
}

// GitView is the git section the composer reads. Every field is produced by
// reading files; the two list fields are never nil (an absent array and an empty
// array mean different things to a renderer, and this packet's rule is that
// "nothing" has to be stated rather than left blank - ticket 145).
//
// remoteBranches is deliberately NOT a field. refs/remotes/* is a cache of the
// last fetch, so shipping it as if it were current state would put a number on
// screen that decays with no indication; census §④ recommends against it and the
// orchestrator's dispatch kept that call ("默认不送").
type GitView struct {
	// Kind is one of the four GitKind* values above.
	Kind string `json:"kind"`
	// Reason is mandatory and never empty: it is the sentence the panel shows
	// when the state is anything but "here is a repo on branch X".
	Reason string `json:"reason"`
	// Branch is the symref name (short: "dev", not "refs/heads/dev").
	Branch string `json:"branch"`
	// DetachedSha is the 40-hex HEAD when there is no branch.
	DetachedSha string `json:"detachedSha"`
	// IsDetached is HEAD's shape, decided by its first five bytes ("ref: ").
	IsDetached bool `json:"isDetached"`
	// RepoRoot is the working tree that owns the shared git dir.
	RepoRoot string `json:"repoRoot"`
	// CurrentWorktree is the directory this run is scoped to. It is the SAME
	// value WorkspaceView.Canonical carries (cmd/wisp/panel_pump.go), so the two
	// sections of one packet cannot name different trees.
	CurrentWorktree string `json:"currentWorktree"`
	// Worktrees is main-first, then attached, each sorted by path.
	Worktrees []GitWorktree `json:"worktrees"`
	// Branches is the full local set: refs/heads walked RECURSIVELY (a branch
	// named "feat/x" is a directory "feat" holding a file "x", so one-level
	// listing drops it - census §④ measured 3 vs the real 4 here) merged with
	// packed-refs, loose winning on a name clash.
	Branches []string `json:"branches"`
	// SwitchBlocked is the status field, not git data: empty would mean
	// "switchable", which is false in this tree, so it carries the reason.
	SwitchBlocked string `json:"switchBlocked"`
}

// GitViewNotProbed is what the pump sends when the assembly attached no reader
// at all. Same rule as the mode section: a value nobody read is reported as
// unread, never as a plausible default.
func GitViewNotProbed() GitView {
	return newGitView(GitKindUnreadable, gitNoWorkspaceReason)
}

func newGitView(kind, reason string) GitView {
	return GitView{
		Kind:          kind,
		Reason:        reason,
		Worktrees:     []GitWorktree{},
		Branches:      []string{},
		SwitchBlocked: GitSwitchBlockedReason,
	}
}

// ReadGit answers the git dimension for one workspace directory.
//
// The probe starts at the workspace root and walks UP, because the folder a run
// is scoped to is very often a subdirectory of the repository (census §②: "从
// 工作区路径起逐级往上"). Walking up stops at the volume root. Nothing here
// resolves a path through C26 and nothing here decides permissions - it reads
// metadata about a tree the operator already pointed the host at, and reports
// what the filesystem says, including "I was refused".
func ReadGit(workspaceRoot string) GitView {
	if strings.TrimSpace(workspaceRoot) == "" {
		return GitViewNotProbed()
	}
	view := newGitView(GitKindUnreadable, "")
	view.CurrentWorktree = workspaceRoot

	loc, kind, reason := locateGitDirs(workspaceRoot)
	if kind != GitKindRepo {
		view.Kind = kind
		view.Reason = reason
		return view
	}

	branch, sha, detached, headKind, headReason := readHead(loc.worktreeGitDir)
	view.RepoRoot = loc.repoRoot
	switch headKind {
	case GitKindRepo:
		view.Kind = GitKindRepo
		view.Branch = branch
		view.DetachedSha = sha
		view.IsDetached = detached
		if detached {
			view.Reason = "这棵工作树是分离头指针：HEAD 直接指向提交 " + sha + "，不属于任何分支"
		} else {
			view.Reason = "这棵工作树当前检出分支 " + branch
		}
	default:
		view.Kind = headKind
		view.Reason = headReason
	}

	view.Branches = listLocalBranches(loc.commonGitDir)

	// The main tree is not in <common>/worktrees/, so it is added by hand
	// from the shared dir's own HEAD, and only then does the enumeration start.
	mainBranch, mainSha, mainDetached, _, _ := readHead(loc.commonGitDir)
	view.Worktrees = append(view.Worktrees, GitWorktree{
		Path: loc.repoRoot, Branch: mainBranch, Sha: mainSha,
		Detached: mainDetached, Main: true,
	})
	view.Worktrees = append(view.Worktrees, listLinkedWorktrees(loc.commonGitDir)...)
	sort.SliceStable(view.Worktrees, func(i, j int) bool {
		if view.Worktrees[i].Main != view.Worktrees[j].Main {
			return view.Worktrees[i].Main
		}
		return view.Worktrees[i].Path < view.Worktrees[j].Path
	})
	if view.Kind == GitKindRepo {
		view.Reason += "；本地分支 " + strconv.Itoa(len(view.Branches)) + " 枚，工作树 " +
			strconv.Itoa(len(view.Worktrees)) + " 枚（含主工作树）"
	}
	return view
}

// ReadGitForWorkspace is the account-aware entry to the git dimension (ticket
// 181's AC#6 debt, ticket 102's rule).
//
// ReadGit takes a bare string and therefore cannot know whether the coordinate
// it was handed is the one the operator named. A WorkspaceView carries C26's
// rewrite account on its own fields (composer.go:189-191 — Rewritten, plus
// Spelling at :184), which is exactly what its header promises: the account
// "must be visible as such in the panel instead of being presented as 'your
// folder'". So the leg that renders the git dimension consumes the account here
// instead of reaching past it for Canonical.
//
// A view that was never rewritten is the ordinary road and behaves exactly as
// before, byte for byte.
func ReadGitForWorkspace(ws WorkspaceView) GitView {
	if ws.Rewritten {
		return GitViewRewritten(ws)
	}
	return ReadGit(ws.Canonical)
}

// GitViewRewritten is the git dimension for a workspace coordinate that C26's
// expansion step moved onto another tree.
//
// DISPOSITION = AC#β, and the orchestrator rules on it; this is the single line
// of the packet where that ruling lands. What is implemented here is candidate
// ① — report the expanded truth AND say out loud that the path was rewritten,
// the shape internal/tools/paths.go:82-93 already uses for config roots
// (RewrittenRoots is recorded and reported, never silently swallowed).
// Candidate ② would be `return newGitView(GitKindUnreadable, <reason>)`,
// dropping the whole dimension; candidate ③ — show Canonical with no word of
// the account — is the disease this function exists to close, and is not
// implemented anywhere.
//
// MEASURED TODAY: this branch is unreachable on the current tree. The only
// producer that ever sets WorkspaceView.Rewritten is RequestWorkspaceSwitch
// (workspace.go:104-111), and its success path runs after res.Actable() at
// workspace.go:85, which errors on any rewritten spelling — so a view that
// reaches the packet has Rewritten == false. The branch is defence in depth for
// the day the upstream loosens, which is the same reason workspace.go:83-85
// re-reads the account "on the native side of the boundary".
func GitViewRewritten(ws WorkspaceView) GitView {
	view := ReadGit(ws.Canonical)
	view.Reason = "注意：这条工作区路径被 C26 的展开步改写过（账户 rewritten=true），" +
		"下面报的是改写后的树 " + ws.Canonical + "，不是原拼法 " + ws.Spelling + "。" + view.Reason
	return view
}

// gitDirs is the pair of locations a probe resolves: the per-tree git dir
// (HEAD, index, commondir) and the shared git dir (refs, packed-refs,
// worktrees/), plus the working tree that owns the shared one.
type gitDirs struct {
	worktreeGitDir string
	commonGitDir   string
	repoRoot       string
}

// locateGitDirs walks up from start looking for a .git, in either shape.
func locateGitDirs(start string) (gitDirs, string, string) {
	for dir := start; ; {
		dot := filepath.Join(dir, ".git")
		st, err := os.Stat(dot)
		switch {
		case err == nil && st.IsDir():
			// Ordinary repository: the .git directory is both the tree's own
			// git dir and the shared one.
			return gitDirs{worktreeGitDir: dot, commonGitDir: dot, repoRoot: dir},
				GitKindRepo, ""
		case err == nil:
			// Linked tree: .git is a pointer file.
			return resolvePointerFile(dot, dir)
		case errors.Is(err, fs.ErrNotExist):
			parent := filepath.Dir(dir)
			if parent == dir {
				return gitDirs{}, GitKindNotARepo,
					"这棵工作区目录及其以上都没有 .git，所以没有分支可显示"
			}
			dir = parent
		case errors.Is(err, fs.ErrPermission):
			return gitDirs{}, GitKindPermissionDenied,
				"探测路径的权限被拒绝：能看到这里可能是一棵 git 树，但宿主读不到，面板不会把它说成不是仓库"
		default:
			return gitDirs{}, GitKindUnreadable, ioReason(err)
		}
	}
}

// resolvePointerFile eats the tree-shaped .git: a file whose one line is
// "gitdir: <path>", absolute or relative to the directory holding the file.
func resolvePointerFile(dot, dir string) (gitDirs, string, string) {
	raw, err := os.ReadFile(dot)
	if err != nil {
		if errors.Is(err, fs.ErrPermission) {
			return gitDirs{}, GitKindPermissionDenied, ioReason(err)
		}
		return gitDirs{}, GitKindUnreadable, ioReason(err)
	}
	pointer := parseGitdirFile(string(raw))
	if pointer == "" {
		return gitDirs{}, GitKindUnreadable, ".git 文件存在，但里面没有可解析的 gitdir: 指针"
	}
	if !filepath.IsAbs(pointer) {
		// Relative gitdir pointers are written relative to the directory that
		// holds the .git file, which is `dir` above.
		pointer = filepath.Join(dir, pointer)
	}
	st, err := os.Stat(pointer)
	switch {
	case err == nil && st.IsDir():
		common := pointer
		if cd, ok := readCommondir(pointer); ok {
			common = cd
		}
		repoRoot := filepath.Dir(common)
		return gitDirs{worktreeGitDir: pointer, commonGitDir: common, repoRoot: repoRoot},
			GitKindRepo, ""
	case err == nil:
		return gitDirs{}, GitKindUnreadable,
			".git 指向的 git 目录存在但不是目录：" + pointer
	case errors.Is(err, fs.ErrPermission):
		return gitDirs{}, GitKindPermissionDenied, ioReason(err)
	default:
		// THE point of this branch: the pointer file is present, so this was a
		// tree; the tree it names is gone. "not_a_repo" would be a lie about
		// what the user is standing in.
		return gitDirs{}, GitKindUnreadable,
			".git 指向的 git 目录已不存在（上一棵树可能被删了）：" + pointer
	}
}

// parseGitdirFile strips git's "gitdir: " header. It returns "" for anything
// else-shaped, which the caller reports instead of guessing at.
func parseGitdirFile(content string) string {
	for _, l := range strings.Split(content, "\n") {
		l = strings.TrimSpace(l)
		if v, ok := strings.CutPrefix(l, "gitdir:"); ok {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

// plainGitdirPath reads the reverse pointer that sits at
// <common>/worktrees/<name>/gitdir. git writes it as a BARE path; the "gitdir: "
// header belongs to the other file (the one inside a linked tree). Both spellings
// are accepted here, because dropping every attached path to an empty string over
// a prefix is the failure this function was written after being measured live.
func plainGitdirPath(content string) string {
	line := firstLine(content)
	if v, ok := strings.CutPrefix(line, "gitdir:"); ok {
		return strings.TrimSpace(v)
	}
	return line
}

// readCommondir follows the commondir file, which points from one linked tree's
// private git dir to the shared one (this repository: "../..").
func readCommondir(worktreeGitDir string) (string, bool) {
	raw, err := os.ReadFile(filepath.Join(worktreeGitDir, "commondir"))
	if err != nil {
		return "", false
	}
	rel := strings.TrimSpace(string(raw))
	if rel == "" {
		return "", false
	}
	return filepath.Join(worktreeGitDir, rel), true
}

// readHead classifies a HEAD file: symref, detached sha, or a failure with its
// own kind. kind is GitKindRepo when a usable answer came back.
func readHead(gitDir string) (branch, sha string, detached bool, kind, reason string) {
	raw, err := os.ReadFile(filepath.Join(gitDir, "HEAD"))
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return "", "", false, GitKindUnreadable, "git 目录里没有 HEAD：" + gitDir
	case errors.Is(err, fs.ErrPermission):
		return "", "", false, GitKindPermissionDenied, ioReason(err)
	case err != nil:
		return "", "", false, GitKindUnreadable, ioReason(err)
	}
	line := firstLine(string(raw))
	// The shape git writes is "ref: refs/heads/<name>"; the prefix test is on the
	// five bytes "ref: ", exactly as census §② names the detached test.
	if ref, ok := strings.CutPrefix(line, "ref: "); ok {
		ref = strings.TrimSpace(ref)
		if ref == "" {
			return "", "", false, GitKindUnreadable, "HEAD 的 ref: 行没有指向任何引用"
		}
		short := ref
		if trimmed, ok := strings.CutPrefix(ref, "refs/heads/"); ok {
			short = trimmed
		}
		return short, "", false, GitKindRepo, ""
	}
	if isHex(line) {
		return "", line, true, GitKindRepo, ""
	}
	return "", "", false, GitKindUnreadable, "HEAD 既不是 ref: 行也不是裸提交号：" + line
}

// listLocalBranches is refs/heads walked recursively, merged with packed-refs.
//
// The recursion is the fix for the ruler this ticket's dispatch carried and the
// census then overturned: `ls .git/refs/heads | wc -l` said 3 while
// `git branch --list` says 4, because "dsh/feat/frontend-p0" is a directory
// "dsh" inside a directory "feat" inside a directory "heads".
func listLocalBranches(commonGitDir string) []string {
	loose := map[string]bool{}
	walkRefs(filepath.Join(commonGitDir, "refs", "heads"), "", loose)
	out := make([]string, 0, len(loose))
	for name := range loose {
		out = append(out, name)
	}
	// packed-refs only fills names loose refs do not already carry; a name in
	// both places is the loose one's, which is what git itself does.
	for name := range readPackedRefs(commonGitDir) {
		if !loose[name] {
			loose[name] = true
			out = append(out, name)
		}
	}
	sort.Strings(out)
	return out
}

// walkRefs collects every branch name under dir, recursing through the nested
// directories that slash-named branches become.
func walkRefs(dir, prefix string, set map[string]bool) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return // no loose refs at all is a legal repo state, not a failure here
	}
	for _, e := range entries {
		name := prefix + e.Name()
		if e.IsDir() {
			walkRefs(filepath.Join(dir, e.Name()), name+"/", set)
			continue
		}
		set[name] = true
	}
}

// readPackedRefs parses .git/packed-refs: "<sha> refs/heads/<name>", "#" header
// lines and "^<sha>" peeled lines skipped.
func readPackedRefs(commonGitDir string) map[string]bool {
	out := map[string]bool{}
	raw, err := os.ReadFile(filepath.Join(commonGitDir, "packed-refs"))
	if err != nil {
		return out
	}
	for _, l := range strings.Split(strings.ReplaceAll(string(raw), "\r\n", "\n"), "\n") {
		l = strings.TrimSpace(l)
		if l == "" || strings.HasPrefix(l, "#") || strings.HasPrefix(l, "^") {
			continue
		}
		_, ref, ok := strings.Cut(l, " ")
		if !ok {
			continue
		}
		ref = strings.TrimSpace(ref)
		if name, ok := strings.CutPrefix(ref, "refs/heads/"); ok && name != "" {
			out[name] = true
		}
	}
	return out
}

// listLinkedWorktrees enumerates <common>/worktrees/* - the ATTACHED trees only.
// Each one's own root comes from its gitdir file, which names the ".git" pointer
// file inside that tree, so the root is its parent (census §③).
func listLinkedWorktrees(commonGitDir string) []GitWorktree {
	entries, err := os.ReadDir(filepath.Join(commonGitDir, "worktrees"))
	if err != nil {
		return nil
	}
	out := make([]GitWorktree, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		dir := filepath.Join(commonGitDir, "worktrees", e.Name())
		branch, sha, detached, _, _ := readHead(dir)
		path := ""
		if raw, err := os.ReadFile(filepath.Join(dir, "gitdir")); err == nil {
			// This file is the OTHER spelling: git writes the reverse pointer as
			// a bare path, with no "gitdir: " header (measured live on this
			// repository's two attached trees - requiring the prefix here reads
			// every attached path as empty).
			if p := plainGitdirPath(string(raw)); p != "" {
				path = filepath.Dir(p)
			}
		}
		// An unresolvable pointer keeps an empty path: this list is data the user
		// reads to decide which tree to work in, and a guessed path is worse than
		// a visible gap.
		out = append(out, GitWorktree{Path: path, Branch: branch, Sha: sha, Detached: detached})
	}
	return out
}

// firstLine returns the first non-empty line of a git metadata file, with the
// line endings git writes on either OS normalised away.
func firstLine(content string) string {
	for _, l := range strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n") {
		if l = strings.TrimSpace(l); l != "" {
			return l
		}
	}
	return ""
}

// isHex recognises a bare commit id.
func isHex(s string) bool {
	if len(s) != 40 {
		return false
	}
	for _, r := range s {
		switch {
		case r >= '0' && r <= '9':
		case r >= 'a' && r <= 'f':
		case r >= 'A' && r <= 'F':
		default:
			return false
		}
	}
	return true
}

// ioReason turns an IO failure into the sentence the panel shows. The error text
// is kept: "读不到" without a name is the "未知" ticket 147 refused.
func ioReason(err error) string {
	return "git 这一维读取失败：" + err.Error()
}
