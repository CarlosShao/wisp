//go:build windows

package tools

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"testing"
)

// Ticket 174, AC#3 (i) + (ii), landed by 174-r4. This file adds NO production
// behaviour: it turns the two claims the fix rests on into criteria that go red
// when someone breaks them, which is what "自证" means here.
//
// AC#3(ii) - "判这条路径在不在授权根里，必须走 C26/PathCanonicalizer，不许在它
// 之外用 filepath.Clean|Abs 做文件系统决策". The prose in task.go:818-820 says
// that today, and CI's d22scan only greps the two TOKENS: a hand-rolled
// `strings.HasPrefix(raw, root)` decision (or handing InAllowlist the RAW
// string instead of the canonical one) passes d22scan clean. So a token scan is
// not evidence for this AC; the verdict has to be shown to move when - and only
// when - C26's answer moves, over paths where the two spellings DISAGREE:
//
//	TestPointerAuthorityFollowsC26NotTheSpelling174r4
//	    C26 says INSIDE, lexical prefix on the filed string says OUTSIDE
//	    (an 8.3 short name, taken from the OS, never typed by hand).
//	TestPointerAuthorityNarrowsWithWorkspaceNotLexicalRoot174r4
//	    C26 says OUTSIDE, lexical prefix on the filed string says INSIDE
//	    (ticket 92 AC#3's workspace narrowing, over the SAME real file).
//
// Each arm carries its own control, because "it was quiet" and "it spoke" are
// both vacuous unless the same input flips when the authority answer flips:
// an empty allowlist must make the short name speak, and clearing the workspace
// must make the narrowed file quiet again. The silence arms additionally check
// what fs.read really returns through the SAME bridge, so a notice that merely
// sounds honest while fs.read disagrees goes red too (that drift is the exact
// failure this ticket exists to prevent).
//
// AC#3(i) - AGENTS §1.2 bans turning the host's artifacts into a gated Tool; the
// neighbouring risk for this fix is the existence probe becoming a write. The
// doc comment at task.go:821-824 promises "Read-only: this function never
// creates, moves, truncates or removes anything" and no criterion anywhere
// checked it.
//
//	TestPointerCheckWritesNothing174r4
//
// Injection surface (AGENTS §1.3): the real bridge, the real C26 canonicalizer,
// real fs.read, real files, real OS short-name queries. Nothing is mocked,
// nothing is skipped: where a prerequisite is missing (no 8.3 alias on this
// volume) the test FAILS LOUDLY with the prerequisite spelled out, which is
// this repository's standing rule for short-name fixtures.

// shortName174r4 asks the OS for the 8.3 spelling of an EXISTING path. A reply
// identical to the long spelling means this volume gives this name no usable
// alias, and then an arm that needs two spellings measures nothing - so it
// fails instead of passing vacuously.
func shortName174r4(t *testing.T, long string) string {
	t.Helper()
	p16, err := syscall.UTF16PtrFromString(long)
	if err != nil {
		t.Fatalf("UTF16PtrFromString(%q): %v", long, err)
	}
	n, err := syscall.GetShortPathName(p16, nil, 0)
	if err != nil || n == 0 {
		t.Fatalf("前置条件缺失：%s 拿不到 8.3 短名（GetShortPathNameW: %v, n=%d）。"+
			"需要该卷开启短名生成（fsutil 8dot3name query %s），改它要管理员权限。", long, err, n, long)
	}
	buf := make([]uint16, n)
	if n, err = syscall.GetShortPathName(p16, &buf[0], n); err != nil || n == 0 {
		t.Fatalf("GetShortPathNameW(%q) 第二遍: %v (n=%d)", long, err, n)
	}
	short := syscall.UTF16ToString(buf[:n])
	if short == long {
		t.Fatalf("前置条件缺失：%s 的 8.3 短名与长名逐字相同（卷未给这枚名字生成别名），"+
			"两形对拼在这一发上量不到任何东西。", long)
	}
	return short
}

// dirNames174r4 snapshots what a directory holds, so "the check wrote nothing"
// is measured instead of asserted.
func dirNames174r4(t *testing.T, dir string) []string {
	t.Helper()
	des, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir(%s): %v", dir, err)
	}
	names := make([]string, 0, len(des))
	for _, de := range des {
		names = append(names, de.Name())
	}
	sort.Strings(names)
	return names
}

func sameNames(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// TestPointerAuthorityFollowsC26NotTheSpelling174r4 is AC#3(ii)'s first arm: the
// copy file is a real regular file under an authorized root, and the string the
// host filed is its 8.3 alias. C26 expands the alias and answers INSIDE, so the
// reply must stay silent - while any decision made by comparing the filed
// string against the root lexically would answer OUTSIDE and spray a false
// "读不到" over a pointer that works.
func TestPointerAuthorityFollowsC26NotTheSpelling174r4(t *testing.T) {
	dir := tempCanonical(t)
	body := asciiRun(20000)
	good := mustWriteFile(t, filepath.Join(dir, "tool-output-174r4-alias.txt"), body)
	short := shortName174r4(t, good)
	t.Logf("two spellings of one file: long=%q short=%q", good, short)

	roster := NewTaskRoster()
	roster.Record("bg-alias", TaskOutput{Text: body, ArtifactPath: short})
	b := judgedTaskBridge(t, roster, NewPathCanonicalizer([]string{dir}, nil))

	out := callTaskOutput(t, b, "bg-alias")
	t.Logf("alias arm verbatim: IsError=%v Truncated=%v Text=%q",
		out.IsError, out.Truncated, announceOf(out.Text))
	if out.IsError || !out.Truncated {
		t.Fatalf("a recorded long output must still be answered as a truncation: %+v", out)
	}
	for _, banned := range []string{noticeLead, "读不到", "不在你被授权的目录范围内", "未接线", "不可找回"} {
		if strings.Contains(out.Text, banned) {
			t.Errorf("a pointer C26 authorizes must stay silent whatever spelling the host filed;"+
				" found %q in %q - that is a decision made on the string, not on C26's answer",
				banned, announceOf(out.Text))
		}
	}
	if !strings.Contains(out.Text, short) || pointerRe.FindStringSubmatch(out.Text) == nil {
		t.Fatalf("the pointer must still name the filed path, got %q", announceOf(out.Text))
	}
	// The silence has to be TRUE: following the pointer as printed works.
	rd := readThroughBridge(t, b, short)
	if rd.IsError || rd.RiskLevel != "L0" || rd.Text != body {
		t.Fatalf("control leg: fs.read on the alias must return the whole %d bytes at L0,"+
			" got isError=%v level=%v len=%d", len(body), rd.IsError, rd.RiskLevel, len(rd.Text))
	}

	// Reverse control: the SAME alias, authority removed (empty allowlist is
	// the factory default, SPEC-03 allowed_dirs[]=[]) - now it must speak. This
	// is what keeps the silence above from being a check that never fires.
	roster2 := NewTaskRoster()
	roster2.Record("bg-alias-noauth", TaskOutput{Text: body, ArtifactPath: short})
	b2 := judgedTaskBridge(t, roster2, NewPathCanonicalizer(nil, nil))
	out2 := callTaskOutput(t, b2, "bg-alias-noauth")
	t.Logf("alias arm, no root verbatim: IsError=%v Truncated=%v Text=%q",
		out2.IsError, out2.Truncated, announceOf(out2.Text))
	if !strings.Contains(out2.Text, noticeLead) || !strings.Contains(out2.Text, "不在你被授权的目录范围内") {
		t.Fatalf("with nothing authorized the alias arm must say the pointer cannot be followed, got %q",
			announceOf(out2.Text))
	}
	if rd2 := readThroughBridge(t, b2, short); !rd2.IsError || rd2.RiskLevel != "L2" {
		t.Fatalf("the notice says fs.read is refused; same bridge says %+v", rd2)
	}
}

// TestPointerAuthorityNarrowsWithWorkspaceNotLexicalRoot174r4 is AC#3(ii)'s
// second arm and the dangerous direction: a real file under an authorized root,
// with a workspace chosen elsewhere in that same root (ticket 92 AC#3 narrows
// scope). Lexically the filed string IS inside the configured root, so a
// hand-rolled containment answer would stay silent - while C26 answers OUTSIDE
// and fs.read really is refused. The reply must follow C26.
func TestPointerAuthorityNarrowsWithWorkspaceNotLexicalRoot174r4(t *testing.T) {
	dir := tempCanonical(t)
	body := asciiRun(20000)
	proj := filepath.Join(dir, "proj174r4")
	if err := os.MkdirAll(proj, 0o700); err != nil {
		t.Fatal(err)
	}
	spill := filepath.Join(dir, "artifacts174r4")
	if err := os.MkdirAll(spill, 0o700); err != nil {
		t.Fatal(err)
	}
	good := mustWriteFile(t, filepath.Join(spill, "tool-output-174r4-narrowed.txt"), body)
	// The filed string sits inside the configured root by simple prefix:
	if !strings.HasPrefix(good, dir) {
		t.Fatalf("fixture broken: %q must sit lexically under %q for this arm to mean anything", good, dir)
	}

	paths := NewPathCanonicalizer([]string{dir}, nil)
	res, err := paths.ResolveWorkspace(proj)
	if err != nil {
		t.Fatalf("ResolveWorkspace(%s): %v", proj, err)
	}
	if err := paths.SetWorkspaceRoot(res.Canonical, res); err != nil {
		t.Fatalf("SetWorkspaceRoot(%s): %v", res.Canonical, err)
	}

	roster := NewTaskRoster()
	roster.Record("bg-narrowed", TaskOutput{Text: body, ArtifactPath: mustCanonical(t, good)})
	b := judgedTaskBridge(t, roster, paths)

	out := callTaskOutput(t, b, "bg-narrowed")
	t.Logf("workspace-narrowed arm verbatim: IsError=%v Truncated=%v Text=%q",
		out.IsError, out.Truncated, announceOf(out.Text))
	if out.IsError || !out.Truncated {
		t.Fatalf("the tool still answers with the truncation stub: %+v", out)
	}
	if !strings.Contains(out.Text, noticeLead) || !strings.Contains(out.Text, "不在你被授权的目录范围内") {
		t.Fatalf("C26 says this file is outside the narrowed scope, so the pointer must say so; got %q",
			announceOf(out.Text))
	}
	if rd := readThroughBridge(t, b, good); !rd.IsError || rd.RiskLevel != "L2" {
		t.Fatalf("the notice must agree with what fs.read really does under the same narrowing,"+
			" got %+v", rd)
	}

	// Reverse control: same file, same root, workspace cleared - C26 now answers
	// INSIDE, and the reply must go back to being silent with a working read.
	paths.ClearWorkspace()
	roster.Record("bg-widened", TaskOutput{Text: body, ArtifactPath: mustCanonical(t, good)})
	out2 := callTaskOutput(t, b, "bg-widened")
	t.Logf("workspace cleared verbatim: IsError=%v Truncated=%v Text=%q",
		out2.IsError, out2.Truncated, announceOf(out2.Text))
	for _, banned := range []string{noticeLead, "读不到", "未接线", "不可找回"} {
		if strings.Contains(out2.Text, banned) {
			t.Errorf("clearing the narrowing must silence the notice; found %q in %q",
				banned, announceOf(out2.Text))
		}
	}
	if rd2 := readThroughBridge(t, b, good); rd2.IsError || rd2.RiskLevel != "L0" || rd2.Text != body {
		t.Fatalf("control leg: fs.read must return the whole %d bytes at L0 once the workspace"+
			" is cleared, got isError=%v level=%v len=%d", len(body), rd2.IsError, rd2.RiskLevel, len(rd2.Text))
	}
}

// TestPointerCheckWritesNothing174r4 is AC#3(i)'s shape for this fix: the reply
// may probe whether the copy file exists, but the probe must stay read-only -
// it may not materialize, create, truncate or move anything in the artifacts
// tree. AGENTS §1.2 bans making the host's artifacts a gated Tool; the creep
// this pins against is subtler and lives on the same two lines: turning
// os.Stat(raw) into os.OpenFile(raw, O_CREATE) would make a "can I still read
// this?" check into a write to a directory no root authorizes.
func TestPointerCheckWritesNothing174r4(t *testing.T) {
	dir := tempCanonical(t)
	body := asciiRun(20000)
	ghost := filepath.Join(dir, "tool-output-174r4-ghost.txt")
	if _, err := os.Lstat(ghost); !os.IsNotExist(err) {
		t.Fatalf("the ghost fixture must be absent for this leg to mean anything: %v", err)
	}
	good := mustWriteFile(t, filepath.Join(dir, "tool-output-174r4-intact.txt"), body)
	before := dirNames174r4(t, dir)
	goodInfo, err := os.Lstat(good)
	if err != nil {
		t.Fatal(err)
	}

	roster := NewTaskRoster()
	roster.Record("bg-ghost", TaskOutput{Text: body, ArtifactPath: mustCanonical(t, ghost)})
	roster.Record("bg-good", TaskOutput{Text: body, ArtifactPath: mustCanonical(t, good)})
	b := judgedTaskBridge(t, roster, NewPathCanonicalizer([]string{dir}, nil))

	ghostOut := callTaskOutput(t, b, "bg-ghost")
	t.Logf("ghost arm verbatim: IsError=%v Truncated=%v Text=%q",
		ghostOut.IsError, ghostOut.Truncated, announceOf(ghostOut.Text))
	// Errorf, not Fatalf: "说了没有" and "写没写" are two questions, and stopping
	// at the first must not blind the second - a fix that keeps the wording while
	// creating the file has to be caught by the filesystem facts below.
	if !strings.Contains(ghostOut.Text, "并不存在") || !strings.Contains(ghostOut.Text, noticeLead) {
		t.Errorf("a pointer at a missing copy file must say so, got %q", announceOf(ghostOut.Text))
	}
	if _, err := os.Lstat(ghost); !os.IsNotExist(err) {
		t.Errorf("answering a pointer at an absent copy file must NOT create it: Lstat=%v", err)
	}

	goodOut := callTaskOutput(t, b, "bg-good")
	t.Logf("intact arm verbatim: IsError=%v Truncated=%v Text=%q",
		goodOut.IsError, goodOut.Truncated, announceOf(goodOut.Text))
	for _, banned := range []string{noticeLead, "读不到", "未接线", "不可找回"} {
		if strings.Contains(goodOut.Text, banned) {
			t.Errorf("a healthy pointer must stay silent; found %q in %q", banned, announceOf(goodOut.Text))
		}
	}

	after := dirNames174r4(t, dir)
	if !sameNames(before, after) {
		t.Errorf("the pointer check wrote into the artifacts tree: before=%v after=%v", before, after)
	}
	goodInfoAfter, err := os.Lstat(good)
	if err != nil {
		t.Fatalf("the intact copy file must still be there: %v", err)
	}
	if !goodInfoAfter.ModTime().Equal(goodInfo.ModTime()) || goodInfoAfter.Size() != goodInfo.Size() {
		t.Errorf("answering task.output must not touch the copy file it only points at:"+
			" mtime %v -> %v, size %d -> %d",
			goodInfo.ModTime(), goodInfoAfter.ModTime(), goodInfo.Size(), goodInfoAfter.Size())
	}
	if data, err := os.ReadFile(good); err != nil || string(data) != body {
		t.Errorf("the copy file's bytes changed (err=%v len=%d)", err, len(data))
	}
	// Pointers are still owed in both arms - a "read-only fix" that dropped them
	// would pass the two assertions above for the wrong reason.
	if pointerRe.FindStringSubmatch(ghostOut.Text) == nil || pointerRe.FindStringSubmatch(goodOut.Text) == nil {
		t.Errorf("both replies must still carry the pointer, got %q / %q",
			announceOf(ghostOut.Text), announceOf(goodOut.Text))
	}
}
