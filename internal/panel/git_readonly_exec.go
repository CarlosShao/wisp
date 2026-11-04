package panel

// Ticket 189: the two halves of the review section that files cannot answer.
//
// WHAT THIS FILE IS. The one place in this package that starts a process, and the
// only implementation of the two dimensions ticket 189 adds under composer.git:
// `uncommitted` (how many changes are not committed yet, split into staged and
// unstaged, with the untracked half in its own column) and `diff` (the per-file
// diff with a byte budget, a truncation verdict and the command that recovers the
// whole thing). Both are read out of `git` and nothing else.
//
// WHY A PROCESS AT ALL, WHEN TICKET 181 NEEDED NONE. Ticket 181 answers three
// questions - is this a tree, which branch, which worktrees - and all three are in
// file names and file contents (.git/HEAD, refs/heads, .git/worktrees), so it needs
// no binary and git.go says so. These two dimensions are different in kind:
//
//   - the STAGED column of the count is "index versus HEAD", and HEAD is a tree
//     object inside a packfile. Reading it in Go means re-implementing pack idx v2,
//     offset and tree deltas, zlib inflate and the text conversion rules; the census
//     (docs/evidence/s1/189-uncommitted-count-design-a1.md section 3.2) measured this
//     repository's object store at 2,016 loose + 17,653 packed objects and 4 packs,
//     and go.mod is frozen, so no git library is available either.
//   - the UNTRACKED column is "which files are not ignored", and the ignore
//     language is gitignore(5). The census measured the naive walk at 8,684
//     candidates against git's authoritative 46 - a 188x over-count that arrives
//     with no error. tools/d22scan/gitignore.go is 817 lines of that language and
//     still asks git for the answer.
//   - the DIFF is a content comparison against HEAD, so it inherits the object store
//     problem in full.
//
// So this file rents correctness from git, once, and pays for it with a dependency
// the product states out loud: no git binary means these two dimensions report the
// third state (GitReadNoGit) and never a zero. git.go's header went from "NO
// EXTERNAL PROCESS" to the honest sentence this file makes true.
//
// WHAT THIS FILE IS NOT. It is not a tool: nothing here is registered in
// internal/tools, D34's table carries no git row, and no model call reaches it -
// TestGitDimensionHasNoModelCallableTool keeps that door shut. It is not a write
// path in any mode: the command shapes are fixed literals, every one of them is a
// query (`status --porcelain`, `diff --numstat`, `diff -- <path>`), and
// git_readonly_boundary_189_test.go reads this file's AST to prove there is no
// second shape. It is not a display scrub: paths and diff text it returns are
// somebody else's strings and reach a human through the packet unwashed, which
// leaves ticket 216's roster (display_scrub_216.go) as the one place that has to
// grow for them - named here rather than half-done here.
//
// TIMEOUTS. One command, one context deadline (context.WithTimeout). That is the
// monotonic clock D42#9 and d22scan ban #4 ask for; there is no time.Now().Sub
// anywhere in this file, and the diff cost is bounded by the byte budget rather
// than by a stopwatch.
//
// THE PATH IS NOT TAKEN FROM THE CALLER'S STRING. Every command runs with `-C
// <tree>`, where <tree> is C26's own answer for the workspace root (risk.Resolve,
// the resolver AGENTS.md 1.2 names and d22scan ban #2 enforces). A path C26 refuses
// - an unexcepted reparse traversal, for instance - is reported as unreadable and
// is never handed to a subprocess.

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/CarlosShao/wisp/internal/risk"
)

// ------------------------------------------------------------- view types
//
// The four types below are the packet's own shape for the two dimensions, and
// they are declared here rather than in git.go for the reason the header gives:
// the reader and the shape it answers with are one object, and splitting them
// across two files is how a field quietly stops being produced. Both views are
// VALUE types on purpose (design §5/F1): no handle, no writer, no pointer a
// caller could use to push bytes back through - writing would require changing
// the shape, which is a visible edit instead of a quiet one.

// GitChangedFile is one row of the uncommitted list.
//
// Staged and Unstaged are the two columns of `git status` XY carried separately
// rather than collapsed into Status, because a path can be BOTH (stage a change,
// then edit the file again: "MM") and a renderer that only received one word
// would have to guess which half it was looking at. Untracked is a third
// category, not a third column: `??` has no HEAD side at all, which is why it
// gets its own column in the count and its own sentence in the diff.
type GitChangedFile struct {
	// Path is the repository-relative path exactly as git printed it (byte for
	// byte; a path is not required to be UTF-8 and is not scrubbed here).
	Path string `json:"path"`
	// Status is one of the gitStatus* words below.
	Status string `json:"status"`
	// Staged says the index differs from HEAD for this path.
	Staged bool `json:"staged"`
	// Unstaged says the worktree differs from the index for this path.
	Unstaged bool `json:"unstaged"`
	// Untracked says this path is not in the index at all ("??").
	Untracked bool `json:"untracked"`
}

// GitUncommittedCounts is the three-column count. Every field counts FILES, not
// hunks and not lines.
//
// The columns deliberately overlap: Staged + Unstaged + Untracked is NOT
// necessarily Total, because a path can be both staged and unstaged (see
// GitChangedFile). Total is the number of rows `git status` produced, and the
// sentence uncommittedSentence states all of them so a reader never has to add
// them up and guess which arithmetic was intended.
type GitUncommittedCounts struct {
	// Total is the number of changed paths: staged, unstaged and untracked.
	Total int `json:"total"`
	// Staged is the number of paths whose index entry differs from HEAD.
	Staged int `json:"staged"`
	// Unstaged is the number of paths whose worktree bytes differ from the index.
	Unstaged int `json:"unstaged"`
	// TrackedModified is the subset of the two columns whose word is "modified"
	// (M or T on either side).
	TrackedModified int `json:"trackedModified"`
	// TrackedDeleted is the subset whose word is "deleted" (D on either side).
	TrackedDeleted int `json:"trackedDeleted"`
	// Untracked is the `??` column. It is NOT part of TrackedModified or
	// TrackedDeleted and has no HEAD-side content for the diff to compare.
	Untracked int `json:"untracked"`
}

// GitUncommittedView is the `composer.git.uncommitted` sub-section.
//
// Counts is a POINTER and it is omitted (not zeroed) when Kind is not
// GitReadCounted. That is the whole rule of this ticket's count half: "0
// uncommitted" is the sentence "this tree is clean", and a dimension nobody read
// must never be able to wear it. A renderer that sees no counts object knows to
// show Reason instead of a number.
type GitUncommittedView struct {
	// Kind is GitReadCounted when Counts is present, and a GitKind* value or
	// GitReadNoGit when nobody counted. The vocabulary is the git section's own
	// four states plus the two read states, not a second enum invented here.
	Kind string `json:"kind"`
	// Reason is the human sentence. It is mandatory and never empty; for
	// GitReadCounted it states all three columns (see uncommittedSentence).
	Reason string `json:"reason"`
	// Counts is nil exactly when Kind is not GitReadCounted.
	Counts *GitUncommittedCounts `json:"counts,omitempty"`
	// Files is always an array, never nil, and capped at gitChangedFileCap.
	Files []GitChangedFile `json:"files"`
	// FilesTotal is the exact number of rows, whether or not Files was cut.
	FilesTotal int `json:"filesTotal"`
	// FilesTruncated says Files holds only the first gitChangedFileCap rows.
	FilesTruncated bool `json:"filesTruncated"`
}

// GitDiffFile is one `git diff HEAD` row plus at most gitDiffPerFileByteCap bytes
// of its text.
//
// The truncated / bytesTotal / recoverPointer triple is the ticket-164 lesson
// made structural: `Truncated=true` alone would be "done but not working", so a
// truncated row carries the exact size it WOULD have had (BytesTotal, only when
// that file's diff was actually run) and the command that retrieves the whole
// thing (RecoverPointer). Truncated=false with BytesTotal absent means the row
// was never materialized at all (past gitDiffMaxFiles), which the row still says
// honestly via Truncated.
type GitDiffFile struct {
	// Path is the destination path git printed.
	Path string `json:"path"`
	// Status is the same gitStatus* word the count section published for this
	// path when it had one (statusWordFor), so the two sections cannot label the
	// same file differently.
	Status string `json:"status"`
	// Additions and Deletions are `--numstat` line counts. A binary file
	// reports 0/0 with Binary=true instead of a made-up line count.
	Additions int `json:"additions"`
	Deletions int `json:"deletions"`
	// Binary says numstat reported "-" for a line count; no text was requested.
	Binary bool `json:"binary"`
	// Text is the diff body, at most gitDiffPerFileByteCap bytes and cut on a
	// rune boundary. Empty is possible for a binary row and for a row that was
	// not materialized; Truncated and Binary say which.
	Text string `json:"text"`
	// BytesTotal is the exact byte count of this file's whole diff. It is
	// ABSENT (nil) when the file's diff was never run, because 0 would be a lie
	// about its size, and it is present-but-zero is impossible from git.
	BytesTotal *int64 `json:"bytesTotal,omitempty"`
	// BytesIncluded is how many of those bytes Text holds (0 when not run).
	BytesIncluded int64 `json:"bytesIncluded"`
	// Truncated says Text is shorter than BytesTotal, or that the row was not
	// materialized inside this snapshot's budget.
	Truncated bool `json:"truncated"`
	// RecoverPointer is the copy-pasteable command that returns this file's
	// whole diff, e.g. `git --no-optional-locks -C "<tree>" diff HEAD -- <path>`.
	RecoverPointer string `json:"recoverPointer"`
}

// GitDiffView is the `composer.git.diff` sub-section.
//
// Files is always an array and never nil. FilesTotal counts the files
// `git diff HEAD --numstat` named (tracked changes only - untracked files have
// no HEAD-side baseline and are called out in Reason instead of being listed
// with invented content).
type GitDiffView struct {
	// Kind is GitReadCounted when this dimension was read, otherwise a GitKind*
	// value or GitReadNoGit.
	Kind string `json:"kind"`
	// Reason is the mandatory human sentence (see diffSentence).
	Reason string `json:"reason"`
	// Files is the per-file list, always an array, never nil.
	Files []GitDiffFile `json:"files"`
	// FilesTotal is how many tracked files differ from HEAD.
	FilesTotal int `json:"filesTotal"`
	// FilesIncluded is how many of those rows were fully answered: the diff text
	// was materialized, or the file is binary and the verdict needs no text.
	FilesIncluded int `json:"filesIncluded"`
	// BytesIncluded is the text bytes actually sent for this snapshot.
	BytesIncluded int64 `json:"bytesIncluded"`
	// ByteCap is the budget BytesIncluded was held inside.
	ByteCap int `json:"byteCap"`
	// Truncated says something was left out (files, text, or both).
	Truncated bool `json:"truncated"`
	// RecoverPointer is the whole-tree recovery command
	// (`git --no-optional-locks -C "<tree>" diff HEAD --`).
	RecoverPointer string `json:"recoverPointer"`
}

// The two read states these dimensions add to the git section's four GitKind*
// values. GitReadCounted is the only state that carries numbers; it is spelled
// separately from "repo" because a repo can be on a clean tree (0 counts) or on
// a tree nobody could read (no counts), and collapsing those two is the exact
// lie ticket 145 refuses in every other section.
const (
	// GitReadCounted says the dimension was read and the numbers are real
	// (which may be zero).
	GitReadCounted = "counted"

	// GitReadNoGit is the third state the census demanded: the machine has no
	// git binary, so the dimension is unread and says so instead of reporting 0.
	GitReadNoGit = "no_git"
)

// gitNoBinaryReason is the sentence for GitReadNoGit. It names the missing
// thing and states the consequence, the same shape
// tools/d22scan/gitignore.go:49-62 uses for the same absence.
const gitNoBinaryReason = "这台机器上没有 git 可执行文件，所以这一维没有读：未提交枚数与逐文件差异都报读取失败，不报 0（0 枚的意思是这棵树干净，是另一件事）"

// The per-path status words. They are a closed vocabulary the count section and
// the diff section share, which is why they are constants rather than strings
// spelled at each site. Values are lower-case words so the packet's JSON stays
// uniform.
const (
	gitStatusModified   = "modified"
	gitStatusAdded      = "added"
	gitStatusDeleted    = "deleted"
	gitStatusRenamed    = "renamed"
	gitStatusCopied     = "copied"
	gitStatusTypechange = "typechange"
	gitStatusConflicted = "conflicted"
	gitStatusUntracked  = "untracked"
)

const (
	// gitBinaryName is the single binary this package starts. It is a literal on
	// purpose: the boundary nail refuses a second `exec.Command*` call, so widening
	// this to a variable would be the first step of a second command shape.
	gitBinaryName = "git"

	// gitExecTimeout bounds ONE command. Every command here is a read of a tree the
	// host already holds, and the census measured this repository's whole `git diff`
	// at 267 ms, so five seconds is a stuck-subprocess bound rather than a budget.
	gitExecTimeout = 5 * time.Second

	// gitStatusByteCap and gitNumstatByteCap are the retention caps for the two
	// listing commands. Both are past anything a review pane can read (this
	// repository's status at the census anchor was 66 KB for 1,322 rows) and both are
	// COUNTING sinks: if one is cut, the counts are reported as unavailable instead of
	// being computed from a partial stream.
	gitStatusByteCap  = 4 << 20
	gitNumstatByteCap = 4 << 20

	// gitDiffByteCap is the whole-snapshot text budget for the diff dimension, and
	// gitDiffPerFileByteCap is the per-file slice of it. They exist because the
	// census measured this repository's diff at 1,521,885 bytes / 17,374 lines: the
	// snapshot is pushed on every state change, so an unbounded diff would turn the
	// panel's pipe into a freight route.
	gitDiffByteCap        = 256 << 10
	gitDiffPerFileByteCap = 32 << 10

	// gitDiffMaxFiles bounds how many per-file commands one snapshot runs. The rest
	// of the files still travel as rows (path, additions, deletions, the recover
	// pointer) with truncated=true, because "this file was not sent" and "this file
	// does not exist" are different facts and a row is the cheap way to say the first.
	gitDiffMaxFiles = 12

	// gitChangedFileCap bounds the file list in the count section. The exact total
	// travels in filesTotal whether or not the list was cut.
	gitChangedFileCap = 200
)

// gitAnswer is one read-only command's answer. text holds at most the retention cap;
// total counts every byte the process wrote, so a caller can report "how big was it"
// even when it did not keep the bytes. kind is "" when the command answered, and one
// of the GitRead*/GitKind* values when it did not.
type gitAnswer struct {
	text   []byte
	total  int64
	cut    bool
	kind   string
	reason string
}

func (a gitAnswer) ok() bool { return a.kind == "" }

// cappedSink is the retention sink: it keeps the first `keep` bytes and counts the
// rest without holding them. Two things depend on that and both are the ticket's
// point - a pathological diff cannot become a pathological allocation, and bytesTotal
// is exact even when bytesIncluded was cut.
type cappedSink struct {
	keep  int64
	buf   []byte
	total int64
	cut   bool
}

func (s *cappedSink) Write(p []byte) (int, error) {
	s.total += int64(len(p))
	if s.keep > 0 {
		n := int64(len(p))
		if n > s.keep {
			n = s.keep
		}
		s.buf = append(s.buf, p[:n]...)
		s.keep -= n
		if n < int64(len(p)) {
			s.cut = true
		}
	} else if len(p) > 0 {
		s.cut = true
	}
	return len(p), nil
}

// runGitReadOnly runs one `git --no-optional-locks -C <tree> <args...>` with a
// deadline and a counting sink. It never returns a partial answer as a success: an
// exit status, a deadline or a missing binary all come back as a kind and a sentence.
//
// `--no-optional-locks` is not decoration. The census measured this repository's
// .git/index unchanged across a `git status` (md5 identical before and after), but
// the existence of the flag is git's own admission that the default path takes a
// lock, and a lock is taken in order to write. Asking for no optional locks removes
// that write from the road instead of trusting it not to happen.
func runGitReadOnly(tree string, keep int64, args ...string) gitAnswer {
	ctx, cancel := context.WithTimeout(context.Background(), gitExecTimeout)
	defer cancel()

	argv := make([]string, 0, len(args)+3)
	argv = append(argv, "--no-optional-locks", "-C", tree)
	argv = append(argv, args...)

	cmd := exec.CommandContext(ctx, gitBinaryName, argv...)
	// A sink, not Output(): Output() would hold the whole stream before anyone can
	// bound it, and a diff of a generated file is exactly the stream that must not
	// be held.
	sink := &cappedSink{keep: keep}
	var stderr bytes.Buffer
	cmd.Stdout = sink
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err == nil {
		return gitAnswer{text: sink.buf, total: sink.total, cut: sink.cut}
	}
	switch {
	case errors.Is(err, exec.ErrNotFound):
		return gitAnswer{kind: GitReadNoGit, reason: gitNoBinaryReason}
	case ctx.Err() != nil:
		// The deadline is what killed it, so the sentence says that rather than
		// quoting "signal: killed" at a reader.
		return gitAnswer{
			kind: GitKindUnreadable,
			reason: "git 在 " + gitExecTimeout.String() + " 内没有应答（" +
				strings.Join(args, " ") + "）：这一维报读取失败，不报 0",
		}
	default:
		reason := firstLine(stderr.String())
		if reason == "" {
			reason = err.Error()
		}
		return gitAnswer{
			kind:   GitKindUnreadable,
			reason: "git 读这一维失败（" + strings.Join(args, " ") + "）：" + reason,
		}
	}
}

// resolveGitTree is C26's vote, taken once per read. The canonical answer is what
// `-C` receives; a refusal (ErrReparseDenied, or an empty answer) is reported as an
// unreadable dimension, never as a tree the reader picked for itself.
//
// The rewrite account is deliberately NOT re-judged here: GitViewRewritten already
// reports the expansion (candidate 1 of ticket 181 AC#6, ruled by the orchestrator),
// and ReadGit has already read that account before it hands the tree over.
func resolveGitTree(root string) (string, string, string) {
	if strings.TrimSpace(root) == "" {
		return "", GitKindUnreadable, gitNoWorkspaceReason
	}
	res, err := risk.Resolve(root, nil)
	if err != nil {
		return "", GitKindUnreadable,
			"C26 拒绝这条工作区路径，git 这一维不越过它去取数：" + err.Error()
	}
	if strings.TrimSpace(res.Canonical) == "" {
		return "", GitKindUnreadable, "C26 没有给出可用的树，git 这一维没有可探测的目录：" + root
	}
	return res.Canonical, "", ""
}

// readGitWorkingState answers both ticket-189 dimensions for one tree, in the order
// the packet's own consistency needs: one `git status` decides what "not committed"
// means, and the diff is then taken with the SAME definition (staged plus unstaged,
// i.e. `git diff HEAD`), so a packet cannot say "12 uncommitted" in one section and
// list 9 files in the other.
func readGitWorkingState(currentWorktree string) (GitUncommittedView, GitDiffView) {
	tree, kind, reason := resolveGitTree(currentWorktree)
	if kind != "" {
		return gitUncommittedUnavailable(kind, reason), gitDiffUnavailable(kind, reason)
	}

	status := runGitReadOnly(tree, gitStatusByteCap, "status", "--porcelain=v1", "-z")
	if !status.ok() {
		return gitUncommittedUnavailable(status.kind, status.reason),
			gitDiffUnavailable(status.kind, status.reason)
	}
	if status.cut {
		// A cut stream means the counts would be computed from the first N bytes of
		// it. That number would look completely reasonable and be wrong, which is the
		// failure mode the census named as this ticket's most expensive one.
		unreadable := "git status 的输出超过 " + strconv.Itoa(gitStatusByteCap) +
			" 字节上限被截断，枚数不可信，所以不报数"
		return gitUncommittedUnavailable(GitKindUnreadable, unreadable),
			gitDiffUnavailable(GitKindUnreadable, unreadable)
	}

	counts, rows, err := parseGitStatusZ(status.text)
	if err != nil {
		return gitUncommittedUnavailable(GitKindUnreadable, err.Error()),
			gitDiffUnavailable(GitKindUnreadable, err.Error())
	}
	return gitUncommittedCounted(rows, counts), readGitDiff(tree, rows, counts.Untracked)
}

// gitUncommittedCounted builds the count section from the parsed rows. The list is
// capped; the total is not.
func gitUncommittedCounted(rows []GitChangedFile, counts GitUncommittedCounts) GitUncommittedView {
	view := GitUncommittedView{
		Kind:           GitReadCounted,
		Counts:         &counts,
		Files:          make([]GitChangedFile, 0, len(rows)),
		FilesTotal:     len(rows),
		FilesTruncated: len(rows) > gitChangedFileCap,
	}
	if len(rows) > gitChangedFileCap {
		rows = rows[:gitChangedFileCap]
	}
	view.Files = append(view.Files, rows...)
	view.Reason = uncommittedSentence(counts, view.FilesTruncated)
	return view
}

// gitUncommittedUnavailable is every state in which nobody counted. Counts stays nil
// on purpose: the packet must not be able to render "0 changes" for a dimension that
// was never read, because 0 is the sentence "this tree is clean".
func gitUncommittedUnavailable(kind, reason string) GitUncommittedView {
	return GitUncommittedView{Kind: kind, Reason: reason, Files: []GitChangedFile{}}
}

// gitDiffUnavailable is the same rule for the diff: an empty file list with an
// unreadable kind, never an empty file list that reads as "nothing changed".
func gitDiffUnavailable(kind, reason string) GitDiffView {
	return GitDiffView{
		Kind:           kind,
		Reason:         reason,
		Files:          []GitDiffFile{},
		ByteCap:        gitDiffByteCap,
		RecoverPointer: "",
	}
}

// readGitDiff materializes the per-file diff inside the byte budget.
//
// The status rows are handed in for two reasons: the file's status word comes from
// the one reading the count section already published (so the two sections cannot
// disagree about whether a file was added, renamed or deleted), and the untracked
// count lets the sentence say out loud what `git diff` structurally does not cover.
func readGitDiff(tree string, rows []GitChangedFile, untracked int) GitDiffView {
	numstat := runGitReadOnly(tree, gitNumstatByteCap, "diff", "HEAD", "--numstat", "-z")
	if !numstat.ok() {
		return gitDiffUnavailable(numstat.kind, numstat.reason)
	}
	if numstat.cut {
		return gitDiffUnavailable(GitKindUnreadable,
			"git diff HEAD --numstat 的输出超过 "+strconv.Itoa(gitNumstatByteCap)+
				" 字节上限被截断，文件清单不可信，所以不报")
	}
	entries, err := parseGitNumstatZ(numstat.text)
	if err != nil {
		return gitDiffUnavailable(GitKindUnreadable, err.Error())
	}

	statusByPath := make(map[string]string, len(rows))
	for _, r := range rows {
		if !r.Untracked {
			statusByPath[r.Path] = r.Status
		}
	}

	view := GitDiffView{
		Kind:           GitReadCounted,
		Files:          make([]GitDiffFile, 0, len(entries)),
		FilesTotal:     len(entries),
		ByteCap:        gitDiffByteCap,
		RecoverPointer: gitRecoverCommand(tree, ""),
	}
	var used int64
	for i, e := range entries {
		row := GitDiffFile{
			Path:           e.path,
			Status:         statusWordFor(e, statusByPath),
			Additions:      e.additions,
			Deletions:      e.deletions,
			Binary:         e.binary,
			RecoverPointer: gitRecoverCommand(tree, e.path),
		}
		if e.binary {
			// A binary row is not text and does not enter the text budget. Saying so
			// is the honest version of sending "Binary files differ" three times,
			// and the row is COMPLETE as it stands - there is nothing further to
			// fetch - so it counts as included. Not counting it here would make the
			// closing FilesIncluded < FilesTotal check call an untruncated binary
			// snapshot truncated.
			row.Truncated = false
			view.Files = append(view.Files, row)
			view.FilesIncluded++
			continue
		}
		keep := int64(gitDiffPerFileByteCap)
		if rest := gitDiffByteCap - used; rest < keep {
			keep = rest
		}
		if i >= gitDiffMaxFiles || keep <= 0 {
			// Not materialized: the row still names the file, its line counts and the
			// command that retrieves it in full. bytesTotal stays absent because this
			// snapshot never ran that file's diff - 0 would be a lie about its size.
			row.Truncated = true
			view.Truncated = true
			view.Files = append(view.Files, row)
			continue
		}
		ans := runGitReadOnly(tree, keep, "diff", "HEAD", "--no-color", "--unified=3", "--", e.path)
		if !ans.ok() {
			row.Truncated = true
			view.Truncated = true
			view.Files = append(view.Files, row)
			continue
		}
		exact := ans.total
		row.BytesTotal = &exact
		row.BytesIncluded = int64(len(ans.text))
		row.Text = string(trimToRuneBoundary(ans.text))
		row.Truncated = ans.cut
		view.BytesIncluded += int64(len(ans.text))
		view.FilesIncluded++
		used += int64(len(ans.text))
		if row.Truncated {
			view.Truncated = true
		}
		view.Files = append(view.Files, row)
	}
	if view.FilesIncluded < view.FilesTotal {
		view.Truncated = true
	}
	view.Reason = diffSentence(view, untracked)
	return view
}

// ------------------------------------------------------------------ parsing

// parseGitStatusZ reads `git status --porcelain=v1 -z`. Each record is "XY path"
// terminated by NUL; a rename or copy is followed by a second record holding the
// source path (measured against git 2.54 on Windows: `R  c.txt\0b.txt\0`).
func parseGitStatusZ(raw []byte) (GitUncommittedCounts, []GitChangedFile, error) {
	var counts GitUncommittedCounts
	records := splitNUL(raw)
	rows := make([]GitChangedFile, 0, len(records))
	for i := 0; i < len(records); i++ {
		rec := records[i]
		if len(rec) == 0 {
			continue
		}
		if len(rec) < 4 || rec[2] != ' ' {
			return counts, nil, errors.New(
				"git status 的记录形状不认识（不是 XY + 空格 + 路径）：" + string(rec))
		}
		x, y := rec[0], rec[1]
		path := string(rec[3:])
		row := GitChangedFile{Path: path}
		switch {
		case x == '?' && y == '?':
			row.Status = gitStatusUntracked
			row.Untracked = true
			counts.Untracked++
		case x == '!':
			// Ignored files only appear with --ignored, which this file never asks
			// for. If one shows up the stream is not the stream that was requested.
			continue
		default:
			row.Status = statusWord(x, y)
			row.Staged = x != ' ' && x != '?'
			row.Unstaged = y != ' ' && y != '?'
			if row.Staged {
				counts.Staged++
			}
			if row.Unstaged {
				counts.Unstaged++
			}
			if x == 'M' || y == 'M' || x == 'T' || y == 'T' {
				counts.TrackedModified++
			}
			if x == 'D' || y == 'D' {
				counts.TrackedDeleted++
			}
			if x == 'R' || x == 'C' {
				// The next record is the rename/copy source. It is consumed here so it
				// cannot be mistaken for a row of its own (that mistake would add a
				// phantom file and inflate the count by one per rename).
				if i+1 >= len(records) {
					return counts, nil, errors.New("git status 报了一次改名但没有给出源路径")
				}
				i++
			}
		}
		counts.Total++
		rows = append(rows, row)
	}
	return counts, rows, nil
}

// gitNumstatEntry is one `git diff HEAD --numstat -z` record.
type gitNumstatEntry struct {
	path      string
	oldPath   string
	additions int
	deletions int
	binary    bool
}

// parseGitNumstatZ reads `git diff HEAD --numstat -z`: "added\tdeleted\tpath" per
// record, and for a rename the path field is EMPTY and the next two records are the
// source and the destination (measured: `0\t0\t\0b.txt\0c.txt\0`).
func parseGitNumstatZ(raw []byte) ([]gitNumstatEntry, error) {
	records := splitNUL(raw)
	out := make([]gitNumstatEntry, 0, len(records))
	for i := 0; i < len(records); i++ {
		rec := records[i]
		if len(rec) == 0 {
			continue
		}
		// The record is bytes because a path is not required to be UTF-8; only
		// the three tab-separated fields are read as text.
		fields := strings.Split(string(rec), "\t")
		if len(fields) != 3 {
			return nil, errors.New("git diff --numstat 的记录形状不认识：" + string(rec))
		}
		e := gitNumstatEntry{}
		if fields[0] == "-" || fields[1] == "-" {
			e.binary = true
		} else {
			add, err := strconv.Atoi(fields[0])
			if err != nil {
				return nil, errors.New("git diff --numstat 的增行数不是数字：" + string(rec))
			}
			del, err := strconv.Atoi(fields[1])
			if err != nil {
				return nil, errors.New("git diff --numstat 的删行数不是数字：" + string(rec))
			}
			e.additions, e.deletions = add, del
		}
		if fields[2] == "" {
			if i+2 >= len(records) {
				return nil, errors.New("git diff --numstat 报了一次改名但没有给出两个路径")
			}
			e.oldPath = string(records[i+1])
			e.path = string(records[i+2])
			i += 2
		} else {
			e.path = fields[2]
		}
		out = append(out, e)
	}
	return out, nil
}

// splitNUL splits a NUL-terminated record stream, dropping the empty tail git leaves
// behind. Byte-level on purpose: a path is not required to be UTF-8 and a "-z"
// stream is the only shape that hands the raw bytes over.
func splitNUL(raw []byte) [][]byte {
	var out [][]byte
	start := 0
	for i := 0; i < len(raw); i++ {
		if raw[i] != 0 {
			continue
		}
		out = append(out, raw[start:i])
		start = i + 1
	}
	if start < len(raw) {
		// A stream that did not end on NUL still carries its last record; dropping it
		// would silently lose a file from the list.
		out = append(out, raw[start:])
	}
	return out
}

// statusWord turns git's two-column status into one word a pane can show. The staged
// column wins when a path is both staged and unstaged, because the row also carries
// both booleans for a renderer that wants to show "MM" instead.
func statusWord(x, y byte) string {
	switch {
	case x == 'U' || y == 'U' || (x == 'A' && y == 'A') || (x == 'D' && y == 'D'):
		return gitStatusConflicted
	case x == 'R' || y == 'R':
		return gitStatusRenamed
	case x == 'C' || y == 'C':
		return gitStatusCopied
	case x == 'A' || y == 'A':
		return gitStatusAdded
	case x == 'D' || y == 'D':
		return gitStatusDeleted
	case x == 'T' || y == 'T':
		return gitStatusTypechange
	default:
		return gitStatusModified
	}
}

// statusWordFor prefers the word the count section published for that path: one
// snapshot, one reading, so the two sections cannot label the same file differently.
func statusWordFor(e gitNumstatEntry, byPath map[string]string) string {
	if w, ok := byPath[e.path]; ok && w != "" {
		return w
	}
	switch {
	case e.oldPath != "":
		return gitStatusRenamed
	case e.binary:
		return gitStatusModified
	case e.additions == 0 && e.deletions > 0:
		return gitStatusDeleted
	case e.deletions == 0 && e.additions > 0:
		return gitStatusModified
	default:
		return gitStatusModified
	}
}

// ---------------------------------------------------------------- sentences

// uncommittedSentence is the mandatory human sentence for the count. It states all
// three columns because "uncommitted" here means staged AND unstaged AND untracked
// (the orchestrator's ruling on ticket 189: staged changes are still uncommitted,
// because the other reading teaches a user that `git add` made them safe).
func uncommittedSentence(c GitUncommittedCounts, filesCut bool) string {
	sentence := "未提交改动 " + strconv.Itoa(c.Total) + " 枚：已暂存 " + strconv.Itoa(c.Staged) +
		" 枚（含已跟踪的修改 " + strconv.Itoa(c.TrackedModified) + " 枚、删除 " +
		strconv.Itoa(c.TrackedDeleted) + " 枚），未暂存 " + strconv.Itoa(c.Unstaged) +
		" 枚，未跟踪 " + strconv.Itoa(c.Untracked) + " 枚；暂存与未暂存是同一份未提交改动里的两档，都算未提交"
	if filesCut {
		sentence += "；文件清单只列前 " + strconv.Itoa(gitChangedFileCap) + " 枚，总数见 filesTotal"
	}
	return sentence
}

// diffSentence is the diff's mandatory sentence: what was sent, what was not, and
// where the whole of it lives.
func diffSentence(v GitDiffView, untracked int) string {
	sentence := "逐文件差异共 " + strconv.Itoa(v.FilesTotal) + " 枚已跟踪改动的文件，本次送出 " +
		strconv.Itoa(v.FilesIncluded) + " 枚 / " + strconv.FormatInt(v.BytesIncluded, 10) +
		" 字节（单枚上限 " + strconv.Itoa(gitDiffPerFileByteCap) + " 字节、整包上限 " +
		strconv.Itoa(gitDiffByteCap) + " 字节、最多 " + strconv.Itoa(gitDiffMaxFiles) +
		" 枚逐枚取）"
	if v.Truncated {
		sentence += "；有内容被截断，全量按每枚文件自己的 recoverPointer 取回"
	}
	if untracked > 0 {
		sentence += "；另有 " + strconv.Itoa(untracked) +
			" 枚未跟踪文件不在这张 diff 里（git diff 只比对已跟踪文件，未跟踪件没有 HEAD 侧基准）"
	}
	if v.FilesTotal == 0 {
		sentence = "已跟踪的改动 0 枚：git diff HEAD 没有任何输出（这不包括未跟踪文件）"
		if untracked > 0 {
			sentence += "；未跟踪 " + strconv.Itoa(untracked) + " 枚在这一维没有内容可显示"
		}
	}
	return sentence
}

// gitRecoverCommand is the pointer AC#4 requires: the exact command a reader can run
// to get the whole thing back. Without it, truncated=true is "done but not working",
// which is the failure ticket 164's overflow lesson names.
func gitRecoverCommand(tree, path string) string {
	cmd := gitBinaryName + " --no-optional-locks -C " + quoteForDisplay(tree) + " diff HEAD --"
	if path != "" {
		cmd += " " + quoteForDisplay(path)
	}
	return cmd
}

// quoteForDisplay wraps an argument so the pointer is copy-pasteable. It is display
// text: nothing in this file executes it, and the escape is for a human shell, not
// for git's own parser.
func quoteForDisplay(s string) string {
	return `"` + strings.ReplaceAll(s, `"`, `\"`) + `"`
}

// trimToRuneBoundary drops a partial UTF-8 sequence left at the cut point. A byte cut
// in the middle of a rune renders as a replacement character, which a reader would
// misread as "this file contains a strange glyph" rather than "the text was cut here".
func trimToRuneBoundary(b []byte) []byte {
	for i := 0; i < utf8.UTFMax && len(b) > 0; i++ {
		if r, size := utf8.DecodeLastRune(b); r != utf8.RuneError || size > 1 {
			return b
		}
		b = b[:len(b)-1]
	}
	return b
}
