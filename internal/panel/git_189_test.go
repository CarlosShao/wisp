package panel

// Ticket 189 AC#2 / AC#3 / AC#4 / AC#5: the resident criteria for the two
// review dimensions (uncommitted counts and per-file diff).
//
// WHY THE FIXTURE IS A REAL REPOSITORY. AC#2 asks for a value that comes from a
// real repo, and the failure this ticket exists to prevent is a
// plausible-looking number computed from something else: the census
// (docs/evidence/s1/189-uncommitted-count-design-a1.md §3.2) built a pure-Go
// index reader whose field offsets were wrong, and it reported 3,316 dirty
// files out of 3,332 without an error. So every expectation below is written to
// disk by THIS file, and then compared against TWO rulers for the same tree:
// the reader under test, and git's own `status --porcelain` / `diff --numstat`
// output. A reader that agrees with itself but not with git is red here.
//
// The tree contains one file per category, because the whole point of the count
// is that the categories do not collapse: a.txt is STAGED (index differs from
// HEAD), b.txt is UNSTAGED (worktree differs from index), gone.txt was deleted
// without staging it, and new.txt is UNTRACKED (no HEAD side at all).
//
// THE READ-ONLY BOUNDARY (AC#3) IS TWO TESTS, AND NEITHER IS DECORATION.
//   - a measurement: .git/index hashes the same before and after the reader
//     runs, and HEAD / the worktree bytes are unchanged too. This is only "this
//     run did not write" - the census measured the same thing for one `git
//     status` and refused to call it a general law, because the existence of
//     --no-optional-locks is git's own admission that the default path takes a
//     lock, and a lock is taken in order to write.
//   - a structural scan: the two production files have no write sink at all.
//     No os.WriteFile/Create/OpenFile/Remove/Rename/Mkdir/Chmod calls, no
//     import outside a read-only roster, exactly one process constructor
//     (exec.CommandContext, in git_readonly_exec.go only), no exported function
//     in the file that starts the process, and not one mutating git subcommand
//     spelled as a string literal. A positive control plants a write and a
//     `git commit` into a copy and asserts the scan fires, so "quiet" cannot
//     mean "blind".
//
// The dispatch for this leg spelled the index check as md5; this file uses
// sha256, which is the hash the repo already uses (internal/llm tests), and the
// guarantee measured is the same one: the bytes did not change.

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// ------------------------------------------------------------- fixture

// git189CLI runs one git command for the fixture or for the authority reading.
// The two readings must run under the SAME config, or "the reader disagrees
// with git" would only be measuring a config difference.
func git189CLI(t *testing.T, dir string, args ...string) string {
	t.Helper()
	full := append([]string{
		"-C", dir,
		"-c", "user.name=wisp-189",
		"-c", "user.email=wisp-189@example.invalid",
		"-c", "commit.gpgsign=false",
		"-c", "core.autocrlf=false",
	}, args...)
	cmd := exec.Command("git", full...)
	// The system config can carry core.autocrlf (Git for Windows installers
	// write it there); the empty HOME set by git189Fixture removes the global
	// one. Both rulers then read one deterministic config.
	cmd.Env = append(os.Environ(), "GIT_CONFIG_NOSYSTEM=1", "GIT_TERMINAL_PROMPT=0")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s in %s: %v\n%s", strings.Join(args, " "), dir, err, out)
	}
	return string(out)
}

// git189Fixture writes a real repository with exactly one path per category and
// returns its root. The directory name carries a space on purpose: this
// repository's own path has one, and a `-C` that quotes badly would fail here.
func git189Fixture(t *testing.T) string {
	t.Helper()
	// An empty HOME (and USERPROFILE, which is what git for Windows falls back
	// to) keeps a machine-global ignore file from hiding new.txt and keeps a
	// machine-global core.autocrlf from moving the line counts. Both the
	// fixture commands and the production reader inherit this environment.
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")

	root := filepath.Join(t.TempDir(), "review 189")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatalf("create fixture root: %v", err)
	}

	git189CLI(t, root, "init", "-q")
	git189Write(t, root, "a.txt", "alpha\n")
	git189Write(t, root, "b.txt", "beta\n")
	git189Write(t, root, "gone.txt", "gamma\n")
	// A binary file is part of the fixture because AC#4 gives it its own
	// category: no line counts, no text, and no truncation verdict - the row is
	// complete once it says "binary".
	git189WriteBytes(t, root, "bin.dat", []byte{0x00, 0x01, 0x02, 0x00, 'w', 'i', 's', 'p', '\n'})
	git189CLI(t, root, "add", "--", "a.txt", "b.txt", "gone.txt", "bin.dat")
	git189CLI(t, root, "commit", "-q", "-m", "baseline")

	// Staged: the index differs from HEAD.
	git189Write(t, root, "a.txt", "alpha\nstaged change\n")
	git189CLI(t, root, "add", "--", "a.txt")
	// Unstaged: the worktree differs from the index.
	git189Write(t, root, "b.txt", "beta\nunstaged change\n")
	// Deleted, not staged: an unstaged deletion.
	if err := os.Remove(filepath.Join(root, "gone.txt")); err != nil {
		t.Fatalf("remove gone.txt: %v", err)
	}
	// Binary and modified, unstaged: the text budget must not be spent on it.
	git189WriteBytes(t, root, "bin.dat", []byte{0x00, 0x01, 0x02, 0x00, 0x7f, '\n'})
	// Untracked: no HEAD side at all, so the diff dimension must not list it.
	git189Write(t, root, "new.txt", "brand new\n")

	return root
}

func git189Write(t *testing.T, root, rel, content string) {
	t.Helper()
	git189WriteBytes(t, root, rel, []byte(content))
}

func git189WriteBytes(t *testing.T, root, rel string, content []byte) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir for %s: %v", rel, err)
	}
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("write %s: %v", rel, err)
	}
}

// ------------------------------------------------- AC#2: counts and diff come
// from the real repository

func TestGit189UncommittedAndDiffComeFromTheRealRepo(t *testing.T) {
	root := git189Fixture(t)

	view := ReadGit(root)
	if view.Kind != GitKindRepo {
		t.Fatalf("ReadGit(%q).Kind = %q, want %q (reason %q)",
			root, view.Kind, GitKindRepo, view.Reason)
	}
	if view.CurrentWorktree != root {
		t.Errorf("CurrentWorktree = %q, want the scoped directory %q", view.CurrentWorktree, root)
	}

	// ---- the count section ------------------------------------------------
	un := view.Uncommitted
	if un.Kind != GitReadCounted {
		t.Fatalf("uncommitted.Kind = %q (%s), want %q - a real repo with real changes must be counted",
			un.Kind, un.Reason, GitReadCounted)
	}
	if un.Counts == nil {
		t.Fatal("uncommitted.Counts is nil for a counted view: the counts are the dimension")
	}
	want := GitUncommittedCounts{
		Total:           5,
		Staged:          1,
		Unstaged:        3,
		TrackedModified: 3,
		TrackedDeleted:  1,
		Untracked:       1,
	}
	if *un.Counts != want {
		t.Errorf("uncommitted.Counts = %+v, want %+v\n(reason: %s)", *un.Counts, want, un.Reason)
	}
	if un.FilesTotal != 5 || len(un.Files) != 5 {
		t.Errorf("uncommitted files: total=%d listed=%d, want 5/5", un.FilesTotal, len(un.Files))
	}
	if un.FilesTruncated {
		t.Error("uncommitted.FilesTruncated = true on a five-file tree")
	}

	// Ruler 2: derive the same three columns from git's raw porcelain stream.
	status := git189CLI(t, root, "status", "--porcelain=v1")
	total, staged, unstaged, untracked := git189CountsFromPorcelain(t, status)
	if total != un.Counts.Total || staged != un.Counts.Staged ||
		unstaged != un.Counts.Unstaged || untracked != un.Counts.Untracked {
		t.Errorf("count section disagrees with git status: reader total/staged/unstaged/untracked = %d/%d/%d/%d, "+
			"git says %d/%d/%d/%d\nraw status:\n%s",
			un.Counts.Total, un.Counts.Staged, un.Counts.Unstaged, un.Counts.Untracked,
			total, staged, unstaged, untracked, status)
	}

	// The per-path rows: the two columns are kept, not collapsed.
	rows := map[string]GitChangedFile{}
	for _, f := range un.Files {
		rows[f.Path] = f
	}
	git189AssertRow(t, rows, "a.txt", gitStatusModified, true, false, false)
	git189AssertRow(t, rows, "b.txt", gitStatusModified, false, true, false)
	git189AssertRow(t, rows, "gone.txt", gitStatusDeleted, false, true, false)
	git189AssertRow(t, rows, "bin.dat", gitStatusModified, false, true, false)
	git189AssertRow(t, rows, "new.txt", gitStatusUntracked, false, false, true)

	// ---- the diff section -------------------------------------------------
	df := view.Diff
	if df.Kind != GitReadCounted {
		t.Fatalf("diff.Kind = %q (%s), want %q", df.Kind, df.Reason, GitReadCounted)
	}
	if df.FilesTotal != 4 || len(df.Files) != 4 {
		t.Errorf("diff files: total=%d listed=%d, want 4/4 (tracked changes only)", df.FilesTotal, len(df.Files))
	}
	if df.FilesIncluded != 4 {
		t.Errorf("diff.FilesIncluded = %d, want 4 (three texts plus one complete binary verdict)", df.FilesIncluded)
	}
	if df.Truncated {
		t.Errorf("diff.Truncated = true on a four-file tree with nothing cut: %s", df.Reason)
	}
	if df.ByteCap != gitDiffByteCap {
		t.Errorf("diff.ByteCap = %d, want %d", df.ByteCap, gitDiffByteCap)
	}
	if !strings.Contains(df.RecoverPointer, "diff HEAD --") {
		t.Errorf("whole-tree RecoverPointer = %q, want the command that recovers everything", df.RecoverPointer)
	}
	if !strings.Contains(df.Reason, "未跟踪") || !strings.Contains(df.Reason, "1 枚") {
		t.Errorf("diff.Reason does not call out the untracked file it structurally cannot cover: %q", df.Reason)
	}

	// Ruler 2: git's own --numstat for the same tree.
	numstat := git189Numstat(t, git189CLI(t, root, "diff", "HEAD", "--numstat"))
	if len(numstat) != len(df.Files) {
		t.Fatalf("numstat named %d files, the diff view listed %d\nraw numstat:\n%v",
			len(numstat), len(df.Files), numstat)
	}
	diffRows := map[string]GitDiffFile{}
	var textBytes int64
	for _, f := range df.Files {
		diffRows[f.Path] = f
		textBytes += f.BytesIncluded
	}
	for path, wantRow := range numstat {
		row, ok := diffRows[path]
		if !ok {
			t.Errorf("git diff HEAD --numstat names %q, the diff view does not list it", path)
			continue
		}
		if row.Binary != wantRow.binary {
			t.Errorf("%s: Binary = %t, git's numstat says %t", path, row.Binary, wantRow.binary)
		}
		if row.Binary {
			// A binary row is complete without text: no bytes were budgeted, no
			// BytesTotal exists to be cut, and the untruncated verdict is the
			// whole truth about it.
			if row.Text != "" || row.BytesTotal != nil || row.BytesIncluded != 0 {
				t.Errorf("%s: binary row carries text=%q bytesTotal=%v bytesIncluded=%d, want no text and no byte claim",
					path, row.Text, row.BytesTotal, row.BytesIncluded)
			}
			if row.Truncated {
				t.Errorf("%s: a binary row has nothing to truncate, but Truncated = true", path)
			}
			if row.Additions != 0 || row.Deletions != 0 {
				t.Errorf("%s: binary row reports %d/%d line counts, want 0/0", path, row.Additions, row.Deletions)
			}
			continue
		}
		if row.Additions != wantRow.add || row.Deletions != wantRow.del {
			t.Errorf("%s: additions/deletions = %d/%d, git says %d/%d",
				path, row.Additions, row.Deletions, wantRow.add, wantRow.del)
		}
		if row.BytesTotal == nil {
			t.Errorf("%s: BytesTotal is absent although this file's diff was run", path)
			continue
		}
		if *row.BytesTotal != int64(len(row.Text)) || row.BytesIncluded != int64(len(row.Text)) {
			t.Errorf("%s: BytesTotal=%d BytesIncluded=%d len(Text)=%d - a small file was cut without saying so",
				path, *row.BytesTotal, row.BytesIncluded, len(row.Text))
		}
		if row.Truncated {
			t.Errorf("%s: Truncated = true although the whole file fits", path)
		}
		if !strings.Contains(row.RecoverPointer, path) || !strings.Contains(row.RecoverPointer, "diff HEAD --") {
			t.Errorf("%s: RecoverPointer = %q, want the copy-pasteable command naming the path", path, row.RecoverPointer)
		}
	}
	if textBytes != df.BytesIncluded {
		t.Errorf("diff.BytesIncluded = %d, sum of per-file BytesIncluded = %d", df.BytesIncluded, textBytes)
	}
	if strings.Contains(strings.Join(git189DiffPaths(df.Files), ","), "new.txt") {
		t.Error("new.txt (untracked) appears in the diff file list: git diff HEAD has no HEAD side for it")
	}
	// The texts are the real content, not a placeholder.
	if !strings.Contains(diffRows["a.txt"].Text, "+staged change") {
		t.Errorf("a.txt text does not carry the staged line: %q", diffRows["a.txt"].Text)
	}
	if !strings.Contains(diffRows["b.txt"].Text, "+unstaged change") {
		t.Errorf("b.txt text does not carry the unstaged line: %q", diffRows["b.txt"].Text)
	}
	if !strings.Contains(diffRows["gone.txt"].Text, "-gamma") {
		t.Errorf("gone.txt text does not carry the deleted line: %q", diffRows["gone.txt"].Text)
	}
}

// TestGit189NoGitAndNoRepoReportAThirdStateNeverZero is the "0 is a sentence"
// nail. The third state must be separable from "the tree is clean" in the
// packet itself, not only in a comment: a nil Counts pointer and an absent
// "counts" JSON key are what stop a renderer from printing "0 changes" for a
// dimension nobody read.
func TestGit189NoGitAndNoRepoReportAThirdStateNeverZero(t *testing.T) {
	notRepo := filepath.Join(t.TempDir(), "plain dir")
	if err := os.MkdirAll(notRepo, 0o755); err != nil {
		t.Fatalf("create %s: %v", notRepo, err)
	}
	un, df := readGitWorkingState(notRepo)
	if un.Kind == GitReadCounted || un.Counts != nil {
		t.Errorf("not-a-repo: uncommitted kind=%q counts=%v, want an unread state with nil counts",
			un.Kind, un.Counts)
	}
	if un.Reason == "" || df.Reason == "" {
		t.Errorf("not-a-repo: the mandatory reason is empty (uncommitted %q, diff %q)", un.Reason, df.Reason)
	}
	if un.Files == nil || df.Files == nil {
		t.Error("not-a-repo: files is nil; the packet's rule is an empty array, never null")
	}

	// GitViewNotProbed (the pump attached no reader) takes the same road.
	probed := GitViewNotProbed()
	if probed.Uncommitted.Kind == GitReadCounted || probed.Diff.Kind == GitReadCounted {
		t.Errorf("GitViewNotProbed reports a counted dimension: uncommitted=%q diff=%q",
			probed.Uncommitted.Kind, probed.Diff.Kind)
	}

	raw, err := json.Marshal(probed)
	if err != nil {
		t.Fatalf("marshal the not-probed view: %v", err)
	}
	var wire struct {
		Uncommitted map[string]json.RawMessage `json:"uncommitted"`
		Diff        map[string]json.RawMessage `json:"diff"`
	}
	if err := json.Unmarshal(raw, &wire); err != nil {
		t.Fatalf("unmarshal the not-probed view: %v", err)
	}
	if _, ok := wire.Uncommitted["counts"]; ok {
		t.Errorf("not-probed JSON carries a counts key: %s", raw)
	}
	if got := string(wire.Uncommitted["files"]); got != "[]" {
		t.Errorf("not-probed JSON files = %s, want [] (an empty array states 'nothing listed')", got)
	}
	if got := string(wire.Diff["files"]); got != "[]" {
		t.Errorf("not-probed diff files = %s, want []", got)
	}
	if _, ok := wire.Uncommitted["kind"]; !ok {
		t.Errorf("not-probed JSON has no kind: %s", raw)
	}
	if _, ok := wire.Diff["recoverPointer"]; !ok {
		t.Errorf("diff JSON has no recoverPointer: %s", raw)
	}
}

// TestGit189NoBinaryIsTheThirdStateNotZero nails the other third state. The
// census named "portable bundle without git" as a real configuration, and the
// rule is that these two dimensions then read as unread - never as "0
// uncommitted". The fixture is a real repository, so nothing but the missing
// binary can explain the state.
func TestGit189NoBinaryIsTheThirdStateNotZero(t *testing.T) {
	root := git189Fixture(t)

	// The same coordinate must resolve without any process; the two dimensions
	// alone should report GitReadNoGit underneath an otherwise intact git view.
	// An existing but empty directory is a legal PATH with no executables in it,
	// which is the portable-bundle shape the census named.
	emptyPath := t.TempDir()
	t.Setenv("PATH", emptyPath)

	view := ReadGit(root)
	if view.Kind != GitKindRepo {
		t.Fatalf("ReadGit(...).Kind = %q with git missing; the file-based half must still answer %q (reason %q)",
			view.Kind, GitKindRepo, view.Reason)
	}
	if view.Uncommitted.Kind != GitReadNoGit || view.Diff.Kind != GitReadNoGit {
		t.Errorf("with no git binary on PATH: uncommitted=%q diff=%q, want %q for both",
			view.Uncommitted.Kind, view.Diff.Kind, GitReadNoGit)
	}
	if view.Uncommitted.Counts != nil || view.Diff.FilesTotal != 0 {
		t.Errorf("the no-git state carries a number: counts=%v diffTotal=%d",
			view.Uncommitted.Counts, view.Diff.FilesTotal)
	}
	if view.Uncommitted.Reason != gitNoBinaryReason || view.Diff.Reason != gitNoBinaryReason {
		t.Errorf("no-git reasons = %q / %q, want the one sentence %q",
			view.Uncommitted.Reason, view.Diff.Reason, gitNoBinaryReason)
	}
}

// TestGit189UncommittedJSONShape pins the sub-section names and the one pointer
// field this ticket adds, because those are the keys the panel renders.
func TestGit189UncommittedJSONShape(t *testing.T) {
	root := git189Fixture(t)
	view := ReadGit(root)

	raw, err := json.Marshal(view)
	if err != nil {
		t.Fatalf("marshal GitView: %v", err)
	}
	var top map[string]json.RawMessage
	if err := json.Unmarshal(raw, &top); err != nil {
		t.Fatalf("unmarshal GitView: %v", err)
	}
	for _, key := range []string{"uncommitted", "diff"} {
		if _, ok := top[key]; !ok {
			t.Errorf("GitView JSON has no %q sub-section: %s", key, raw)
		}
	}
	// The two new dimensions are NESTED under composer.git; they are not new
	// top-level snapshot keys (that would be a frontend-contract change, Q-51).
	var un GitUncommittedView
	if err := json.Unmarshal(top["uncommitted"], &un); err != nil {
		t.Fatalf("unmarshal uncommitted: %v", err)
	}
	if un.Counts == nil || un.Counts.Total != 5 {
		t.Errorf("wire uncommitted counts = %+v, want total 4", un.Counts)
	}
	var df GitDiffView
	if err := json.Unmarshal(top["diff"], &df); err != nil {
		t.Fatalf("unmarshal diff: %v", err)
	}
	if df.FilesTotal != 4 || df.RecoverPointer == "" {
		t.Errorf("wire diff = total %d recover %q, want total 4 and a recover pointer", df.FilesTotal, df.RecoverPointer)
	}
}

// ------------------------------------------------- AC#3: the read-only boundary

// TestGit189ReaderDoesNotTouchTheRepository measures the one thing a test can
// measure: this run left the repository's own bookkeeping byte-identical.
func TestGit189ReaderDoesNotTouchTheRepository(t *testing.T) {
	root := git189Fixture(t)
	gitDir := filepath.Join(root, ".git")

	indexBefore := git189HashFile(t, filepath.Join(gitDir, "index"))
	headBefore := git189HashFile(t, filepath.Join(gitDir, "HEAD"))
	worktreeBefore := map[string]string{}
	for _, rel := range []string{"a.txt", "b.txt", "new.txt"} {
		worktreeBefore[rel] = git189HashFile(t, filepath.Join(root, rel))
	}

	for i := 0; i < 2; i++ {
		view := ReadGit(root)
		if view.Uncommitted.Kind != GitReadCounted || view.Diff.Kind != GitReadCounted {
			t.Fatalf("round %d did not actually read (uncommitted %q, diff %q): the no-write reading below "+
				"would then be measuring a reader that never ran", i+1, view.Uncommitted.Kind, view.Diff.Kind)
		}
	}

	if after := git189HashFile(t, filepath.Join(gitDir, "index")); after != indexBefore {
		t.Errorf(".git/index changed across two reads: %s -> %s. --no-optional-locks exists because the "+
			"default `git status` path takes a lock, and a lock is taken in order to write", indexBefore, after)
	}
	if after := git189HashFile(t, filepath.Join(gitDir, "HEAD")); after != headBefore {
		t.Errorf(".git/HEAD changed across two reads: %s -> %s", headBefore, after)
	}
	for rel, before := range worktreeBefore {
		if after := git189HashFile(t, filepath.Join(root, rel)); after != before {
			t.Errorf("worktree file %s changed across the read: %s -> %s", rel, before, after)
		}
	}
}

// TestGit189ReviewPathHasNoWriteSink is the structural half: the read path does
// not merely refrain from writing, it has nothing to write with. The positive
// control at the bottom is the reason a quiet reading above can be trusted.
func TestGit189ReviewPathHasNoWriteSink(t *testing.T) {
	dir := filepath.Join(panelRepoRoot(t), "internal", "panel")
	prod := []string{"git.go", "git_readonly_exec.go"}
	for _, finding := range scanGit189WriteSinks(t, dir, prod) {
		t.Errorf("the git review path grew a write sink: %s", finding)
	}

	// Functional shape (design §5/F1): the entries take a coordinate and return
	// a value. No []byte, no io.Writer, no *os.File, no command name, no
	// context-write channel - a write would first have to change the shape.
	entry := reflect.TypeOf(ReadGit)
	if entry.NumIn() != 1 || entry.In(0).Kind() != reflect.String {
		t.Errorf("ReadGit signature changed: %v; it must take one workspace coordinate", entry)
	}
	if entry.NumOut() != 1 || entry.Out(0) != reflect.TypeOf(GitView{}) {
		t.Errorf("ReadGit returns %v; it must return the GitView value", entry)
	}
	state := reflect.TypeOf(readGitWorkingState)
	if state.NumIn() != 1 || state.In(0).Kind() != reflect.String {
		t.Errorf("readGitWorkingState signature changed: %v; it must take one tree coordinate", state)
	}
	if state.NumOut() != 2 || state.Out(0) != reflect.TypeOf(GitUncommittedView{}) ||
		state.Out(1) != reflect.TypeOf(GitDiffView{}) {
		t.Errorf("readGitWorkingState returns %v; it must return the two value views", state)
	}

	// Positive control: plant a write and an action subcommand into a COPY of
	// the production files and require the scan to name all of them. A nail that
	// cannot go red is decoration (ticket 190's S3 lesson).
	planted := t.TempDir()
	for _, name := range prod {
		src, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			t.Fatalf("read %s for the planted copy: %v", name, err)
		}
		body := string(src)
		if name == "git_readonly_exec.go" {
			body += "\nfunc Planted189Write(p string) { _ = os.WriteFile(p, nil, 0o644); _ = exec.Command(\"git\", \"commit\", \"-m\", \"x\") }\n"
		}
		if err := os.WriteFile(filepath.Join(planted, name), []byte(body), 0o600); err != nil {
			t.Fatalf("write the planted copy of %s: %v", name, err)
		}
	}
	findings := scanGit189WriteSinks(t, planted, prod)
	joined := make([]string, 0, len(findings))
	for _, f := range findings {
		joined = append(joined, f.String())
	}
	all := strings.Join(joined, "\n")
	if !strings.Contains(all, "os.WriteFile") {
		t.Errorf("the scan stayed silent on a planted os.WriteFile - it is blind\nfindings:\n%s", all)
	}
	if !strings.Contains(all, `"commit"`) {
		t.Errorf("the scan stayed silent on a planted `git commit` string literal - it is blind\nfindings:\n%s", all)
	}
	if !strings.Contains(all, "exec.Command") {
		t.Errorf("the scan stayed silent on a planted second process constructor\nfindings:\n%s", all)
	}
	if !strings.Contains(all, "Planted189Write") {
		t.Errorf("the scan stayed silent on a new exported function in the process-starting file\nfindings:\n%s", all)
	}
}

// git189ScanFinding is one structural violation, addressed well enough to fix.
type git189ScanFinding struct {
	file string
	line int
	what string
}

func (f git189ScanFinding) String() string {
	return fmt.Sprintf("%s:%d: %s", f.file, f.line, f.what)
}

// git189ReadOnlyImports is the roster of packages the read path may depend on.
// Not a style rule: os/exec is here because git_readonly_exec.go must start git,
// and every other entry is a read or a value manipulation. A new import has to
// be added here on purpose, which is the point.
var git189ReadOnlyImports = map[string]bool{
	"bytes":         true,
	"context":       true,
	"errors":        true,
	"io/fs":         true,
	"os":            true,
	"os/exec":       true,
	"path/filepath": true,
	"sort":          true,
	"strconv":       true,
	"strings":       true,
	"time":          true,
	"unicode/utf8":  true,
	"github.com/CarlosShao/wisp/internal/risk": true,
}

// git189ReadOnlyOSCalls is the whitelist from the design (189-a1 §5/F2): the
// read path may ask the filesystem questions, never make statements.
var git189ReadOnlyOSCalls = map[string]bool{
	"Stat":     true,
	"Lstat":    true,
	"ReadFile": true,
	"ReadDir":  true,
	"Open":     true,
}

// git189MutatingGitVerbs is the action roster. None of these may appear as a
// string literal anywhere on the read path, because the argv is assembled from
// literals: if the verb is not spelled, it cannot be passed to git.
var git189MutatingGitVerbs = map[string]bool{
	"add": true, "am": true, "apply": true, "checkout": true, "clean": true,
	"commit": true, "fetch": true, "gc": true, "merge": true, "mv": true,
	"prune": true, "pull": true, "push": true, "rebase": true, "reset": true,
	"restore": true, "rm": true, "stash": true, "switch": true, "tag": true,
	"update-index": true,
}

// git189ExecConstructors is the single process constructor the package may use.
var git189ExecConstructors = map[string]bool{"CommandContext": true}

// scanGit189WriteSinks parses the named production files and returns every
// shape that could write or that could name an action. It is shared with the
// positive control on purpose: the control must run the shipped ruler, not a
// copy of it.
func scanGit189WriteSinks(t *testing.T, dir string, names []string) []git189ScanFinding {
	t.Helper()
	fset := token.NewFileSet()
	var out []git189ScanFinding
	execCalls := 0

	for _, name := range names {
		path := filepath.Join(dir, name)
		src, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		file, err := parser.ParseFile(fset, path, src, parser.ParseComments)
		if err != nil {
			t.Fatalf("parse %s: %v", path, err)
		}

		for _, imp := range file.Imports {
			pkg, err := strconv.Unquote(imp.Path.Value)
			if err != nil {
				continue
			}
			if !git189ReadOnlyImports[pkg] {
				out = append(out, git189ScanFinding{
					name, fset.Position(imp.Pos()).Line,
					"import " + strconv.Quote(pkg) + " is not on the read-only roster",
				})
			}
		}

		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			pkgIdent, ok := sel.X.(*ast.Ident)
			if !ok {
				return true
			}
			switch pkgIdent.Name {
			case "os":
				if !git189ReadOnlyOSCalls[sel.Sel.Name] {
					out = append(out, git189ScanFinding{
						name, fset.Position(call.Pos()).Line,
						"os." + sel.Sel.Name + " is not a read-only filesystem call",
					})
				}
			case "exec":
				execCalls++
				if !git189ExecConstructors[sel.Sel.Name] {
					out = append(out, git189ScanFinding{
						name, fset.Position(call.Pos()).Line,
						"exec." + sel.Sel.Name + " is a second process shape; only exec.CommandContext is allowed",
					})
				}
				if name != "git_readonly_exec.go" {
					out = append(out, git189ScanFinding{
						name, fset.Position(call.Pos()).Line,
						"this file starts a process; the package's one named spawn site is git_readonly_exec.go",
					})
				}
			}
			return true
		})

		ast.Inspect(file, func(n ast.Node) bool {
			lit, ok := n.(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				return true
			}
			val, err := strconv.Unquote(lit.Value)
			if err != nil {
				return true
			}
			if git189MutatingGitVerbs[val] {
				out = append(out, git189ScanFinding{
					name, fset.Position(lit.Pos()).Line,
					"mutating git subcommand " + strconv.Quote(val) + " is spelled as a literal; argv is assembled from literals, so this is the first half of an action",
				})
			}
			return true
		})

		if name != "git_readonly_exec.go" {
			continue
		}
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			// Methods are skipped on purpose: cappedSink.Write is the io.Writer
			// the child process's stdout is copied INTO (an in-memory buffer),
			// not a file sink. What matters is that the package exposes no free
			// function that could be handed a destination.
			if !ok || fn.Recv != nil {
				continue
			}
			if ast.IsExported(fn.Name.Name) {
				out = append(out, git189ScanFinding{
					name, fset.Position(fn.Pos()).Line,
					"exported function " + fn.Name.Name + " turns the internal reader into a callable surface",
				})
			}
		}
	}

	if execCalls != 1 {
		out = append(out, git189ScanFinding{
			"git_readonly_exec.go", 0,
			fmt.Sprintf("the package starts %d processes through the exec package; exactly one read-only exec.CommandContext is the contract", execCalls),
		})
	}
	return out
}

// ------------------------------------------------------------- test helpers

func git189HashFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

// git189CountsFromPorcelain derives the three columns straight from git's own
// output - ruler 2 for the count assertions above.
func git189CountsFromPorcelain(t *testing.T, out string) (total, staged, unstaged, untracked int) {
	t.Helper()
	for _, line := range strings.Split(strings.TrimRight(out, "\r\n"), "\n") {
		if line == "" {
			continue
		}
		if len(line) < 4 || line[2] != ' ' {
			t.Fatalf("unexpected porcelain record %q", line)
		}
		x, y := line[0], line[1]
		total++
		if x == '?' && y == '?' {
			untracked++
			continue
		}
		if x != ' ' && x != '?' {
			staged++
		}
		if y != ' ' && y != '?' {
			unstaged++
		}
	}
	return total, staged, unstaged, untracked
}

// git189NumstatRow is one parsed `git diff HEAD --numstat` record. A binary
// record carries "-" in both count columns, which is why `binary` is a field
// rather than an Atoi error.
type git189NumstatRow struct {
	add    int
	del    int
	binary bool
}

// git189Numstat parses `git diff HEAD --numstat` (non -z; the fixture paths
// contain no tab or newline), ruler 2 for the diff assertions above.
func git189Numstat(t *testing.T, out string) map[string]git189NumstatRow {
	t.Helper()
	rows := map[string]git189NumstatRow{}
	for _, line := range strings.Split(strings.TrimRight(out, "\n"), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Split(line, "\t")
		if len(fields) != 3 {
			t.Fatalf("unexpected numstat record %q", line)
		}
		row := git189NumstatRow{}
		if fields[0] == "-" || fields[1] == "-" {
			row.binary = true
		} else {
			add, err := strconv.Atoi(fields[0])
			if err != nil {
				t.Fatalf("numstat additions %q: %v", fields[0], err)
			}
			del, err := strconv.Atoi(fields[1])
			if err != nil {
				t.Fatalf("numstat deletions %q: %v", fields[1], err)
			}
			row.add, row.del = add, del
		}
		rows[fields[2]] = row
	}
	return rows
}

func git189AssertRow(t *testing.T, rows map[string]GitChangedFile, path, status string, staged, unstaged, untracked bool) {
	t.Helper()
	row, ok := rows[path]
	if !ok {
		t.Errorf("%s is missing from the uncommitted rows: %v", path, git189RowPaths(rows))
		return
	}
	if row.Status != status || row.Staged != staged || row.Unstaged != unstaged || row.Untracked != untracked {
		t.Errorf("%s = {status:%q staged:%t unstaged:%t untracked:%t}, want {status:%q staged:%t unstaged:%t untracked:%t}",
			path, row.Status, row.Staged, row.Unstaged, row.Untracked,
			status, staged, unstaged, untracked)
	}
}

func git189RowPaths(rows map[string]GitChangedFile) []string {
	out := make([]string, 0, len(rows))
	for p := range rows {
		out = append(out, p)
	}
	sort.Strings(out)
	return out
}

func git189DiffPaths(files []GitDiffFile) []string {
	out := make([]string, 0, len(files))
	for _, f := range files {
		out = append(out, f.Path)
	}
	sort.Strings(out)
	return out
}
