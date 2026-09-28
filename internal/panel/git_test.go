package panel

// Ticket 181 AC#2/AC#3/AC#4: the git read surface's resident nails.
//
// WHY THESE TESTS PLANT THEIR OWN TREES. AC#2 says a positive criterion's value
// has to come from a real file, and the one version of that which is worth
// nothing is a test that asserts this repository's current branch name - it goes
// green on the machine that was used to write it and red on every other one. So
// every expected value below is a value THIS FILE wrote to disk a few lines
// earlier, and the assertions read back what the reader found.
//
// The three shapes come from the census (docs/evidence/s1/181-186-git-detection-census-c1.md
// §②③④), which measured all three on live data in this repository: the directory
// form, the worktree form where .git is a POINTER FILE (absolute and relative),
// and detached HEAD as a bare 40-hex line. Two of the rulers the dispatch carried
// were overturned by that census and both are nailed here: a one-level
// refs/heads listing (3 vs the real 4), and a worktree list that enumerates
// .git/worktrees only (2 vs the real 3, because the main tree is not in there).

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// ------------------------------------------------------------- fixtures

// mkTree writes files under a fresh temp dir and returns the dir. The sub-name
// carries a space on purpose: this repository's own path has one, and the census
// measured the live gitdir pointer as an absolute path containing a space.
func mkTree(t *testing.T, sub string, files map[string]string) string {
	t.Helper()
	root := filepath.Join(t.TempDir(), sub)
	for rel, content := range files {
		p := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatalf("plant %s: %v", rel, err)
		}
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			t.Fatalf("plant %s: %v", rel, err)
		}
	}
	return root
}

// gitRefFile content is irrelevant to the reader (git stores the sha there); the
// file's EXISTENCE and path are what name a branch.
const gitRefFile = "0123456789abcdef0123456789abcdef0123456789\n"

// ------------------------------------------------------------ AC#2 positive

// TestReadGitSyntheticSymrefRepo is the resident positive criterion: a tree this
// test built, read back through the surface, answering with the values that are
// written in those files.
//
// Mutation: replace listLocalBranches' recursion in walkRefs with a single
// os.ReadDir of refs/heads (the ruler the census overturned) and the branches
// assertion below goes red on the nested "dsh/feat/frontend-p0" entry. Remove
// the hand-added main worktree in ReadGit and the worktrees assertion goes red
// (only .git/worktrees/ is listed, which is the other overturned ruler).
func TestReadGitSyntheticSymrefRepo(t *testing.T) {
	worktree := "feature/x"
	branches := []string{"dev", "master", worktree, "dsh/feat/frontend-p0"}

	files := map[string]string{
		".git/HEAD":                            "ref: refs/heads/" + worktree + "\n",
		".git/refs/heads/dev":                  gitRefFile,
		".git/refs/heads/master":               gitRefFile,
		".git/refs/heads/dsh/feat/frontend-p0": gitRefFile,
		// refs/heads/feature/ is a directory whose entries are branches too.
		".git/refs/heads/feature/x": gitRefFile,
	}
	root := mkTree(t, "my repo", files)

	// The expectation is the string this test wrote, not a name read off the
	// machine running it.
	want := append([]string{}, branches...)
	sort.Strings(want)

	view := ReadGit(root)

	if view.Kind != GitKindRepo {
		t.Fatalf("kind = %q (reason %q), want %q for a tree with a .git directory",
			view.Kind, view.Reason, GitKindRepo)
	}
	if view.Branch != worktree {
		t.Errorf("branch = %q, want %q read from .git/HEAD", view.Branch, worktree)
	}
	if view.IsDetached || view.DetachedSha != "" {
		t.Errorf("a symref HEAD reported detached=%v sha=%q", view.IsDetached, view.DetachedSha)
	}
	if view.RepoRoot != root {
		t.Errorf("repoRoot = %q, want the directory holding .git (%q)", view.RepoRoot, root)
	}
	if view.CurrentWorktree != root {
		t.Errorf("currentWorktree = %q, want the probed directory", view.CurrentWorktree)
	}
	if got := view.Branches; !equalStrings(got, want) {
		t.Errorf("branches = %v, want every file under refs/heads walked recursively: %v", got, want)
	}
	if len(view.Branches) == 3 {
		// The exact shape of the overturned ruler: one-level ls counts
		// dev/master/feature(direct
		// ory) and misses the nested branch.
		t.Errorf("branch count 3 is the one-level listing bug; nested refs/heads/dsh/feat/* was dropped")
	}
	if len(view.Worktrees) != 1 || !view.Worktrees[0].Main || view.Worktrees[0].Path != root {
		t.Errorf("worktrees = %+v, want exactly the main worktree %q for a repo with no attached trees",
			view.Worktrees, root)
	}
	if strings.TrimSpace(view.Reason) == "" {
		t.Error("reason is empty - this packet's rule is that every state says why it is what it is")
	}
	if !strings.Contains(view.Reason, worktree) {
		t.Errorf("reason %q does not name the branch it read from the file", view.Reason)
	}
	if view.SwitchBlocked == "" {
		t.Error("switchBlocked empty would claim the panel can switch trees, which this tree cannot")
	}
}

// --------------------------------------------------------- the other shapes

// TestReadGitDetachedSyntheticTree nails the third shape: HEAD as a bare 40-hex
// line, which the census read off this repository's own worktree
// (D:/tmp/wisp136instr-r1/wt-head) but asserts here from a planted file.
func TestReadGitDetachedSyntheticTree(t *testing.T) {
	sha := strings.Repeat("ab", 20)
	root := mkTree(t, "detached", map[string]string{
		".git/HEAD":           sha + "\n",
		".git/refs/heads/dev": gitRefFile,
	})

	view := ReadGit(root)
	if view.Kind != GitKindRepo {
		t.Fatalf("kind = %q (reason %q), want repo", view.Kind, view.Reason)
	}
	if !view.IsDetached {
		t.Error("isDetached false for a HEAD with no ref: prefix")
	}
	if view.DetachedSha != sha {
		t.Errorf("detachedSha = %q, want the 40-hex line this test wrote (%q)", view.DetachedSha, sha)
	}
	if view.Branch != "" {
		t.Errorf("branch = %q, want empty for detached HEAD", view.Branch)
	}
	if view.Reason == "" || !strings.Contains(view.Reason, sha) {
		t.Errorf("reason = %q, want it to name the commit it read", view.Reason)
	}
	if len(view.Worktrees) != 1 || !view.Worktrees[0].Detached || view.Worktrees[0].Sha != sha {
		t.Errorf("main worktree entry = %+v, want the detached sha", view.Worktrees)
	}
}

// TestReadGitWorktreePointerFile is the linked-worktree shape: .git is a FILE,
// the shared dir comes from commondir, and the repo root is its parent.
func TestReadGitWorktreePointerFile(t *testing.T) {
	branch := "attached/branch"
	mainRoot := mkTree(t, "main tree", map[string]string{
		".git/HEAD":                       "ref: refs/heads/dev\n",
		".git/refs/heads/dev":             gitRefFile,
		".git/refs/heads/attached/branch": gitRefFile,
		// One attached worktree, and - the point of the fixture - it is NOT
		// mirrored in the main tree's own HEAD.
		".git/worktrees/attached/HEAD":      "ref: refs/heads/" + branch + "\n",
		".git/worktrees/attached/commondir": "../..\n",
	})
	attachedDir := filepath.Join(mainRoot, "attached")
	if err := os.MkdirAll(attachedDir, 0o755); err != nil {
		t.Fatal(err)
	}
	// git keeps the reverse pointer in the shared dir: <common>/worktrees/<name>/
	// gitdir names the .git FILE inside that tree, and the reader derives the
	// tree's root as its parent. The content is written here in the shape git
	// itself writes - a BARE path with no "gitdir: " header, which is what the
	// live probe on this repository caught the first version of the reader
	// getting wrong (every attached path came back empty).
	if err := os.WriteFile(filepath.Join(mainRoot, ".git", "worktrees", "attached", "gitdir"),
		[]byte(filepath.Join(attachedDir, ".git")+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// A second attached tree written with the prefixed spelling, so the reader is
	// held to accepting both rather than to whichever one the fixture happened to
	// imitate.
	prefixed := filepath.Join(mainRoot, "attached-prefixed")
	if err := os.MkdirAll(prefixed, 0o755); err != nil {
		t.Fatal(err)
	}
	prefixedGit := filepath.Join(mainRoot, ".git", "worktrees", "prefixed")
	if err := os.MkdirAll(prefixedGit, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(prefixedGit, "HEAD"),
		[]byte("ref: refs/heads/prefixed/pointer\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(prefixedGit, "gitdir"),
		[]byte("gitdir: "+filepath.Join(prefixed, ".git")+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(prefixedGit, "commondir"),
		[]byte("../..\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Absolute pointer, exactly the shape census §② measured on D:/wt/fe.
	pointer := "gitdir: " + filepath.Join(mainRoot, ".git", "worktrees", "attached") + "\n"
	if err := os.WriteFile(filepath.Join(attachedDir, ".git"), []byte(pointer), 0o644); err != nil {
		t.Fatal(err)
	}

	view := ReadGit(attachedDir)
	if view.Kind != GitKindRepo {
		t.Fatalf("kind = %q (reason %q), want repo for a worktree pointer file", view.Kind, view.Reason)
	}
	if view.Branch != branch {
		t.Errorf("branch = %q, want %q from the per-worktree HEAD", view.Branch, branch)
	}
	if view.RepoRoot != mainRoot {
		t.Errorf("repoRoot = %q, want the shared dir's parent (%q) resolved through commondir",
			view.RepoRoot, mainRoot)
	}
	if view.CurrentWorktree != attachedDir {
		t.Errorf("currentWorktree = %q, want the attached tree this probe started in", view.CurrentWorktree)
	}
	// Every tree must be listed: main (which is NOT inside worktrees/) plus both
	// attached ones.
	if len(view.Worktrees) != 3 {
		t.Fatalf("worktrees = %+v, want 3 entries (main + 2 attached)", view.Worktrees)
	}
	if !view.Worktrees[0].Main {
		t.Errorf("worktrees[0] = %+v, want the main tree first", view.Worktrees[0])
	}
	var attached *GitWorktree
	for i := range view.Worktrees {
		if view.Worktrees[i].Path == attachedDir {
			attached = &view.Worktrees[i]
		}
	}
	if attached == nil {
		t.Fatalf("the attached worktree %q is missing from %+v", attachedDir, view.Worktrees)
	}
	if attached.Branch != branch || attached.Main {
		t.Errorf("attached entry = %+v, want branch %q and Main=false", attached, branch)
	}
	var found *GitWorktree
	for i := range view.Worktrees {
		if view.Worktrees[i].Path == prefixed {
			found = &view.Worktrees[i]
		}
	}
	if found == nil || found.Branch != "prefixed/pointer" {
		t.Errorf("the prefixed-spelling tree %q did not resolve: %+v", prefixed, view.Worktrees)
	}
	for i, e := range view.Worktrees {
		if e.Path == "" {
			t.Errorf("worktrees[%d] = %+v, want a resolved path for every entry", i, e)
		}
	}
	if !equalStrings(view.Branches, []string{"attached/branch", "dev"}) {
		t.Errorf("branches = %v, want the shared refs/heads walked recursively", view.Branches)
	}
}

// TestReadGitRelativePointerFile covers the second spelling git writes - census
// §⑨ item 3 lists it as the shape live data could not supply.
func TestReadGitRelativePointerFile(t *testing.T) {
	mainRoot := mkTree(t, "rel main", map[string]string{
		".git/HEAD":                    "ref: refs/heads/dev\n",
		".git/refs/heads/dev":          gitRefFile,
		".git/worktrees/sub/HEAD":      "ref: refs/heads/topic\n",
		".git/worktrees/sub/commondir": "../..\n",
	})
	sub := filepath.Join(mainRoot, "wt", "sub")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	rel := filepath.Join("..", "..", ".git", "worktrees", "sub")
	if err := os.WriteFile(filepath.Join(sub, ".git"), []byte("gitdir: "+rel+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	view := ReadGit(sub)
	if view.Kind != GitKindRepo {
		t.Fatalf("kind = %q (reason %q), want repo for a relative gitdir:", view.Kind, view.Reason)
	}
	if view.Branch != "topic" {
		t.Errorf("branch = %q, want topic", view.Branch)
	}
	if view.RepoRoot != mainRoot {
		t.Errorf("repoRoot = %q, want %q", view.RepoRoot, mainRoot)
	}
}

// TestReadGitVanishedPointerIsUnreadable: a .git FILE whose gitdir: names a
// directory that is gone is UNREADABLE, never not_a_repo.
func TestReadGitVanishedPointerIsUnreadable(t *testing.T) {
	root := mkTree(t, "orphan", map[string]string{
		".git": "gitdir: " + filepath.Join(t.TempDir(), "no-such-dir") + "\n",
	})
	view := ReadGit(root)
	if view.Kind != GitKindUnreadable {
		t.Errorf("kind = %q, want unreadable - the pointer file exists, so this WAS a worktree", view.Kind)
	}
	if view.Reason == "" {
		t.Error("unreadable without a reason")
	}
}

// ------------------------------------------------------- AC#4: not a repo

// TestReadGitNonRepoSaysSoOutLoud is AC#4's real criterion.
func TestReadGitNonRepoSaysSoOutLoud(t *testing.T) {
	root := mkTree(t, "plain folder", map[string]string{"notes.txt": "no git here\n"})

	view := ReadGit(root)
	if view.Kind != GitKindNotARepo {
		t.Fatalf("kind = %q (reason %q), want not_a_repo", view.Kind, view.Reason)
	}
	if strings.TrimSpace(view.Reason) == "" {
		t.Error("reason empty: a workspace that is not a repo must say so, not render a blank")
	}
	if view.Branch != "" || view.DetachedSha != "" || len(view.Branches) != 0 || len(view.Worktrees) != 0 {
		t.Errorf("a non-repo invented git state: %+v", view)
	}
	if view.Branches == nil || view.Worktrees == nil {
		t.Error("lists are nil: the packet owes an empty array, not a null")
	}
	if view.SwitchBlocked == "" {
		t.Error("switchBlocked must still carry its sentence")
	}
}

// TestReadGitWithoutAWorkspaceDoesNotProbe: no narrowed workspace is a different
// fact from "not a repo".
func TestReadGitWithoutAWorkspaceDoesNotProbe(t *testing.T) {
	for _, in := range []string{"", "   "} {
		view := ReadGit(in)
		if view.Kind != GitKindUnreadable {
			t.Errorf("ReadGit(%q).kind = %q, want unreadable", in, view.Kind)
		}
		if view.Reason != gitNoWorkspaceReason {
			t.Errorf("reason = %q, want the no-workspace sentence", view.Reason)
		}
	}
}

// ------------------------------------------------------- AC#3: the reverse

// TestGitDimensionHasNoModelCallableTool is the anti-drift nail. The census drew
// the line this protects: showing a branch goes through the host read surface and
// the snapshot; handing the model a git tool would go through D34's tool table,
// which has no git row, and would turn a display feature into a contract change.
//
// Three doors, all read off the tree rather than asserted from memory:
//  1. D34's authoritative table in docs/PLAN.md declares no git.* row;
//  2. nothing in internal/tools names a tool whose name starts with "git.";
//  3. the C17 whitelist in bridge.go is still exactly the four panel.* methods.
func TestGitDimensionHasNoModelCallableTool(t *testing.T) {
	root := panelRepoRoot(t)

	plan, err := os.ReadFile(filepath.Join(root, "docs", "PLAN.md"))
	if err != nil {
		t.Fatalf("read docs/PLAN.md: %v", err)
	}
	if row := d34ToolRowMentioningGit(string(plan)); row != "" {
		t.Errorf("D34's built-in tool table now has a git row (%s) - the git dimension was turned "+
			"into a model-callable tool, which is a contract change this ticket does not make", row)
	}

	offenders := gitToolNamesUnder(filepath.Join(root, "internal", "tools"))
	if len(offenders) > 0 {
		t.Errorf("a git.* tool name appears in internal/tools production code (%s) - the read surface "+
			"must reach the panel through the snapshot, never through a tool the model can call",
			strings.Join(offenders, ", "))
	}

	wantMethods := []string{
		MethodModeRequest, MethodWorkspaceRequest, MethodAttachmentAdd, MethodMessageSend,
	}
	sort.Strings(wantMethods)
	if got := whitelistMethodsFromSource(t, filepath.Join(root, "internal", "panel", "bridge.go")); !equalStrings(got, wantMethods) {
		t.Errorf("the C17 panel.* whitelist in bridge.go = %v, want the four methods this ticket may not "+
			"extend - read-only display rides the existing snapshot push, so a fifth method means the "+
			"read half got tied back to the action half", got)
	}
}

var (
	gitToolNameRe = regexp.MustCompile(`"(git\.[a-z0-9_.-]+)"`)
	panelMethodRe = regexp.MustCompile(`"(panel\.[a-z0-9_.-]+)"`)
	// d34ToolNameRe matches a tool name the way D34's table writes it - in
	// BACKTICKS, not double quotes. The positive control caught the first version
	// of this door looking for the quoted spelling, which no markdown row can
	// produce: that door was green because it was blind, not because the table
	// was clean.
	d34ToolNameRe = regexp.MustCompile("`(git\\.[a-z0-9_.-]+)`")
)

// whitelistMethodsFromSource reads the method names off the file that declares
// them, because a check that re-lists the Go constants proves nothing: it stays
// green no matter what gets added. Same ruler the census used
// (grep -n '"panel\.' internal/panel/bridge.go).
func whitelistMethodsFromSource(t *testing.T, path string) []string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	seen := map[string]bool{}
	var out []string
	for _, m := range panelMethodRe.FindAllStringSubmatch(string(data), -1) {
		if !seen[m[1]] {
			seen[m[1]] = true
			out = append(out, m[1])
		}
	}
	sort.Strings(out)
	return out
}

// d34ToolRowMentioningGit looks for a markdown table row whose first cell is a
// backticked tool name starting with git. - inside the D34 section.
func d34ToolRowMentioningGit(plan string) string {
	lines := strings.Split(strings.ReplaceAll(plan, "\r\n", "\n"), "\n")
	inD34 := false
	for _, l := range lines {
		if strings.HasPrefix(l, "## ") || strings.HasPrefix(l, "### ") {
			inD34 = strings.Contains(l, "D34")
			continue
		}
		if !inD34 || !strings.HasPrefix(strings.TrimSpace(l), "|") {
			continue
		}
		cell := strings.TrimSpace(strings.Split(strings.TrimSpace(l), "|")[1])
		if m := d34ToolNameRe.FindString(cell); m != "" {
			return strings.TrimSpace(l)
		}
	}
	return ""
}

// composerMethodWhitelist is the whitelist as ONE list, so a fifth method cannot
// slip in by being declared somewhere else in the package.
func composerMethodWhitelist() []string {
	return []string{
		MethodModeRequest, MethodWorkspaceRequest, MethodAttachmentAdd, MethodMessageSend,
	}
}

// gitToolNamesUnder returns every `"git.*"` tool-name literal in a tree's
// production Go files. Split out of the criterion above so that criterion can be
// pointed at a knowingly wrong tree the way TestPlantedComposerModeWriteGoesRed
// does for the mode door: a negative criterion that has never been seen to fire
// is a sentence, not an instrument.
func gitToolNamesUnder(root string) []string {
	var offenders []string
	// Walk errors are ignored deliberately - this reports what it could read, and
	// the caller asserts on the list.
	_ = filepath.Walk(root, func(p string, st os.FileInfo, err error) error {
		if err != nil || st.IsDir() || !strings.HasSuffix(p, ".go") || strings.HasSuffix(p, "_test.go") {
			return nil
		}
		data, rErr := os.ReadFile(p)
		if rErr != nil {
			return nil
		}
		for _, m := range gitToolNameRe.FindAllStringSubmatch(string(data), -1) {
			offenders = append(offenders, filepath.Base(p)+": \""+m[1]+"\"")
		}
		return nil
	})
	sort.Strings(offenders)
	return offenders
}

// TestPlantedGitToolShapesGoRed is the positive control for AC#3's three doors.
// Every plant is the exact drift the door exists to catch, read by the same
// helper the resident criterion runs - and the real tree has to stay quiet under
// it.
func TestPlantedGitToolShapesGoRed(t *testing.T) {
	dir := t.TempDir()
	plant := "package tools\n\n// A branch switcher offered to the model.\n" +
		"func (gitSwitch) Name() string { return \"git.switch\" }\n"
	if err := os.WriteFile(filepath.Join(dir, "git_switch.go"), []byte(plant), 0o644); err != nil {
		t.Fatal(err)
	}
	// A _test.go file is outside the door's file set on purpose and must not fire.
	if err := os.WriteFile(filepath.Join(dir, "noise_test.go"), []byte(plant), 0o644); err != nil {
		t.Fatal(err)
	}
	got := gitToolNamesUnder(dir)
	if len(got) != 1 || !strings.Contains(got[0], "git.switch") {
		t.Errorf("planted a git.switch tool and the scan returned %v, want exactly that one hit", got)
	}
	if clean := gitToolNamesUnder(filepath.Join(panelRepoRoot(t), "internal", "tools")); len(clean) != 0 {
		t.Errorf("the tool door fires on the real tree too: %v", clean)
	}

	// Door 3: a fifth panel.* method declared in the whitelist file.
	five := "package panel\n\nconst (\n" +
		"\tMethodModeRequest      = \"panel.mode.request\"\n" +
		"\tMethodWorkspaceRequest = \"panel.workspace.request\"\n" +
		"\tMethodAttachmentAdd    = \"panel.attachment.add\"\n" +
		"\tMethodMessageSend      = \"panel.message.send\"\n" +
		"\tMethodWorktreeSwitch   = \"panel.worktree.switch\"\n)\n"
	fivePath := filepath.Join(dir, "bridge_five.go")
	if err := os.WriteFile(fivePath, []byte(five), 0o644); err != nil {
		t.Fatal(err)
	}
	if n := len(whitelistMethodsFromSource(t, fivePath)); n != 5 {
		t.Errorf("planted a fifth panel.* method and the whitelist scan counted %d, want 5", n)
	}
	if real := whitelistMethodsFromSource(t, filepath.Join(panelRepoRoot(t), "internal", "panel", "bridge.go")); len(real) != 4 {
		t.Errorf("the whitelist door reads %d methods on the real bridge.go, want 4", len(real))
	}

	// Door 1: a git row inside the D34 section, and the same row outside it.
	if d34ToolRowMentioningGit("### D34 内置工具权威表\n\n| Tool | Level |\n|---|---|\n| `git.branch` | L0 |\n") == "" {
		t.Error("planted a git row inside the D34 section and the table scan did not see it")
	}
	if d34ToolRowMentioningGit("## Something else\n\n| `git.branch` | L0 |\n") != "" {
		t.Error("the D34 door fires outside the D34 section")
	}
}

// ------------------------------------------------------- snapshot carriage

// TestGitSectionTravelsInTheSnapshotPacket is the field-family nail: the git view
// has to reach the panel through the pump, not just exist as a struct.
func TestGitSectionTravelsInTheSnapshotPacket(t *testing.T) {
	root := mkTree(t, "packet tree", map[string]string{
		".git/HEAD":                 "ref: refs/heads/release/9\n",
		".git/refs/heads/release/9": gitRefFile,
	})
	pump := NewSnapshotPump(PumpSources{
		Workspace: func() WorkspaceView { return WorkspaceViewFromRoot(root) },
		Git:       func() GitView { return ReadGit(root) },
	})

	snap := pump.Snapshot()
	if snap.Composer.Git.Branch != "release/9" {
		t.Errorf("composer.git.branch = %q, want release/9 from the planted HEAD", snap.Composer.Git.Branch)
	}
	if snap.Composer.Git.CurrentWorktree != snap.Composer.Workspace.Canonical {
		t.Errorf("git.currentWorktree %q and workspace.canonical %q disagree - one packet, two trees",
			snap.Composer.Git.CurrentWorktree, snap.Composer.Workspace.Canonical)
	}
	data, err := pump.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	var wire struct {
		Composer struct {
			Git struct {
				Kind    string `json:"kind"`
				Reason  string `json:"reason"`
				Branch  string `json:"branch"`
				Main    bool   `json:"main"`
				Blocked string `json:"switchBlocked"`
			} `json:"git"`
		} `json:"composer"`
	}
	if err := json.Unmarshal(data, &wire); err != nil {
		t.Fatal(err)
	}
	if wire.Composer.Git.Kind != GitKindRepo || wire.Composer.Git.Branch != "release/9" {
		t.Errorf("wire git = %+v, want kind %q branch release/9", wire.Composer.Git, GitKindRepo)
	}
	if wire.Composer.Git.Reason == "" || wire.Composer.Git.Blocked == "" {
		t.Errorf("wire git lost its two mandatory sentences: %+v", wire.Composer.Git)
	}
}

// TestPumpWithoutGitReaderSaysUnreadable: a section nobody attached a reader to
// reports that, the same way the mode section refuses to invent a default.
func TestPumpWithoutGitReaderSaysUnreadable(t *testing.T) {
	snap := NewSnapshotPump(PumpSources{}).Snapshot()
	if snap.Composer.Git.Kind != GitKindUnreadable || snap.Composer.Git.Reason == "" {
		t.Errorf("git section = %+v, want an unreadable state with a reason", snap.Composer.Git)
	}
}

// ------------------------------------------------------------- real tree

// TestReadGitOnThisRepositoryIsSelfConsistent runs the surface over the tree the
// tests live in. No expectation here names a branch, so the criterion holds on
// any checkout: it asserts the surface's own answers agree with the files.
func TestReadGitOnThisRepositoryIsSelfConsistent(t *testing.T) {
	root := panelRepoRoot(t)
	view := ReadGit(root)
	if view.Kind != GitKindRepo {
		t.Fatalf("kind = %q reason = %q: this repository has a .git directory", view.Kind, view.Reason)
	}
	head, err := os.ReadFile(filepath.Join(root, ".git", "HEAD"))
	if err != nil {
		t.Fatalf("read .git/HEAD: %v", err)
	}
	line := firstLine(string(head))
	if strings.HasPrefix(line, "ref: ") {
		want := strings.TrimSpace(strings.TrimPrefix(line, "ref: "))
		want = strings.TrimPrefix(want, "refs/heads/")
		if view.Branch != want {
			t.Errorf("branch = %q, want %q as .git/HEAD names it", view.Branch, want)
		}
		if view.IsDetached {
			t.Error("detached for a symref HEAD")
		}
	} else if isHex(line) {
		if !view.IsDetached || view.DetachedSha != line {
			t.Errorf("detached state = %+v, want sha %q", view, line)
		}
	} else {
		t.Fatalf("unrecognised HEAD shape on this machine: %q", line)
	}

	// The recursive walk is what the flat ruler got wrong: count the files under
	// refs/heads here and require the surface to agree, file for file.
	var files []string
	if err := filepath.Walk(filepath.Join(root, ".git", "refs", "heads"), func(p string, st os.FileInfo, err error) error {
		if err != nil || st.IsDir() {
			return nil
		}
		rel, rErr := filepath.Rel(filepath.Join(root, ".git", "refs", "heads"), p)
		if rErr != nil {
			return nil
		}
		files = append(files, filepath.ToSlash(rel))
		return nil
	}); err != nil {
		t.Fatalf("walk refs/heads: %v", err)
	}
	sort.Strings(files)
	if !equalStrings(view.Branches, files) && !hasPackedRefs(filepath.Join(root, ".git", "packed-refs")) {
		t.Errorf("branches = %v, want the files under refs/heads: %v", view.Branches, files)
	}
	if len(view.Branches) != len(files) && !hasPackedRefs(filepath.Join(root, ".git", "packed-refs")) {
		t.Errorf("branch count = %d, want the recursive file count %d (the flat ls counted directories)",
			len(view.Branches), len(files))
	}

	// The main tree must be in the list even though git keeps it out of
	// .git/worktrees/.
	entries, _ := os.ReadDir(filepath.Join(root, ".git", "worktrees"))
	if len(view.Worktrees) != len(entries)+1 {
		t.Errorf("worktrees = %d entries, want the %d attached ones plus the main tree",
			len(view.Worktrees), len(entries))
	}
	if !view.Worktrees[0].Main || view.Worktrees[0].Path != view.RepoRoot {
		t.Errorf("first worktree entry = %+v, want this repository's main tree (%q)",
			view.Worktrees[0], view.RepoRoot)
	}
	// Every ATTACHED entry has to carry a resolved path. This repository's own
	// attached trees are the live data that caught the first version of the
	// reader returning "" for all of them, because git writes the reverse
	// pointer as a bare path and the reader demanded the "gitdir: " header.
	for i, e := range view.Worktrees[1:] {
		if e.Path == "" {
			t.Errorf("attached entry %d = %+v, want the path its gitdir pointer names", i+1, e)
			continue
		}
		if st, err := os.Stat(e.Path); err != nil {
			t.Logf("attached entry %d path not present on this machine (%v) - prunable shape, reported not hidden", i+1, err)
		} else if !st.IsDir() {
			t.Errorf("attached entry %d path %q is not a directory", i+1, e.Path)
		}
	}
}

func hasPackedRefs(p string) bool {
	st, err := os.Stat(p)
	return err == nil && !st.IsDir()
}

// equalStrings is already declared in this package's test files
// (l2_grant_boundary_test.go:1919) with the same signature, so it is used here
// rather than shadowed.

// TestReadGitForWorkspaceConsumesRewriteAccount is ticket 181's AC#6, the
// judgement piece for ticket 102's rule on the git dimension (181-r2): the leg
// that renders git must read the rewrite account the WorkspaceView carries, not
// reach past it for the coordinate alone.
func TestReadGitForWorkspaceConsumesRewriteAccount(t *testing.T) {
	dir := t.TempDir()

	t.Run("not_rewritten_is_the_ordinary_road", func(t *testing.T) {
		ws := WorkspaceView{Set: true, Spelling: dir, Canonical: dir}
		got := ReadGitForWorkspace(ws)
		want := ReadGit(dir)
		if got.Kind != want.Kind || got.Reason != want.Reason || got.CurrentWorktree != want.CurrentWorktree {
			t.Errorf("account-aware entry disagreed with the plain read: got kind=%q reason=%q tree=%q, want kind=%q reason=%q tree=%q",
				got.Kind, got.Reason, got.CurrentWorktree, want.Kind, want.Reason, want.CurrentWorktree)
		}
		if strings.Contains(got.Reason, "c26") {
			t.Errorf("a view with rewritten=false was narrated as rewritten: %q", got.Reason)
		}
	})

	t.Run("rewritten_must_say_so_out_loud", func(t *testing.T) {
		ws := WorkspaceView{Set: true, Spelling: `%WISP181R2%\named`, Canonical: dir, Rewritten: true}
		got := ReadGitForWorkspace(ws)
		if !strings.Contains(got.Reason, "rewritten=true") {
			t.Errorf("reason %q does not name the rewrite account: a moved coordinate would read as the operator's own folder", got.Reason)
		}
		if !strings.Contains(got.Reason, ws.Spelling) || !strings.Contains(got.Reason, ws.Canonical) {
			t.Errorf("reason %q must carry both the spelling the user named and the tree C26 moved it to", got.Reason)
		}
		if got.CurrentWorktree != ws.Canonical {
			t.Errorf("currentWorktree = %q, want the expanded coordinate %q (candidate 2 would drop the dimension instead; that ruling is open)", got.CurrentWorktree, ws.Canonical)
		}
	})

	t.Run("unset_view_probes_nothing", func(t *testing.T) {
		if got := ReadGitForWorkspace(UnsetWorkspaceView()); got.Kind != GitKindUnreadable {
			t.Errorf("kind = %q for an unset workspace, want %q", got.Kind, GitKindUnreadable)
		}
	})
}
