package tools

// ticket 162 AC#3 (line endings and BOM, both directions) and AC#4 (atomicity,
// as a negative experiment). The AC the ticket froze for these two cells:
//
//	:34 AC#3 换行/BOM 正反两向（CRLF 文件改完还是 CRLF；带 BOM 的改完 BOM 还在；
//	           不许把 LF 文件的行尾改掉）
//	:35 AC#4 原子性：改到一半被取消不许留损坏文件（并答"摘掉临时文件那一味，是否
//	           存在一发损坏从此看不见"）
//
// What the code under test actually is (internal/tools/fs_edit.go:136, re-read
// this cell): `content := string(raw)` - the whole file's bytes, with the BOM's
// three bytes and every CR of a CRLF pair inside it. Locating is strings.Count /
// strings.Index over those bytes, and the bytes that go out are
// content[:start] + new + content[stop:]. So the tool never converts anything:
// the only way a line ending or a BOM can change is if it sits INSIDE the
// matched old, or inside the new the caller spelled.
//
// That is a real guarantee and these cases pin it in both directions. It is NOT
// the guarantee 票面 :27 describes ("先剥 BOM、按首次出现判定行尾、内部统一成
// LF、写完还原"), which is a WIDENING of matching: it exists to make an
// LF-spelled old hit a CRLF file. Per this dispatch, AC#3 forensics that
// widening's cost below (TestFSEditLineEndingForensicsRefusesRealSpellings) and
// does not implement it - whether to build it is the orchestrator's call.

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// utf8BOM is the byte-order mark as it exists ON DISK in a UTF-8 file: the three
// bytes EF BB BF, the UTF-8 encoding of U+FEFF. Spelled byte-by-byte because a
// literal BOM inside this source line would be invisible to whoever reads it,
// which is the exact class of thing this cell is about.
const utf8BOM = "\xef\xbb\xbf"

// lfView is the normalization 票面 :27 would install, used here ONLY as the
// measuring stick that says "a normalizing locator would have found this": a
// refusal is only worth registering as a cost if the same old would have hit
// once endings are unified. Both sides go through it, because that proposal's
// own wording is "旧文本与新文本同一把尺".
func lfView(s string) string {
	s = strings.TrimPrefix(s, utf8BOM)
	return strings.ReplaceAll(strings.ReplaceAll(s, "\r\n", "\n"), "\r", "\n")
}

// endings returns how one blob spells its line breaks: (crlf, bare-cr, bare-lf).
// A pure-CRLF blob has bareLF 0; a pure-LF blob has bareCR 0; any other pair is
// a mixed-ending file, which is what a silently rewritten ending turns into.
func endings(s string) (crlf, bareCR, bareLF int) {
	crlf = strings.Count(s, "\r\n")
	bareCR = strings.Count(s, "\r") - crlf
	bareLF = strings.Count(s, "\n") - crlf
	return crlf, bareCR, bareLF
}

// ---------------------------------------------------------------------------
// AC#3 forward: BOM
// ---------------------------------------------------------------------------

// TestFSEditKeepsABOMItWasNotAskedToTouch is AC#3's "带 BOM 的文件改完 BOM 还在"
// in the direction that matters: the bytes before the edit and the bytes after
// it, checked at byte level, not through a string comparison that could hide a
// re-encoding.
//
// 未修码上响不响：这一格没有"响不响"的问题 - 它是取证＋钉住一条不许退让的性质。
// 今天的通路（纯字节匹配 + 只换命中段）本来就不会碰 BOM，所以这条断言在未修码上
// 同样成立，修码只是让它变成**有人守着的**性质而不是巧合。
func TestFSEditKeepsABOMItWasNotAskedToTouch(t *testing.T) {
	root := sealableTempDir124(t)
	original := utf8BOM + "title: 报表\nstatus: draft\n"
	target := filepath.Join(root, "report.md")
	if err := os.WriteFile(target, []byte(original), 0o600); err != nil {
		t.Fatal(err)
	}
	if got := []byte(original)[:3]; string(got) != "\xef\xbb\xbf" {
		t.Fatalf("the fixture itself is wrong, it does not start with EF BB BF: % x", got)
	}
	b, _ := fsEditBridge(t, mustCanonical(t, root))

	out, err := b.Execute(t.Context(), req("fs.edit", editArgs(t, target,
		map[string]any{"old": "status: draft", "new": "status: final"})))
	if err != nil || out.IsError {
		t.Fatalf("out=%+v err=%v", out, err)
	}
	after := readString(t, target)

	// Byte level, as the AC asks it: the leading three bytes are still the BOM.
	if len(after) < 3 || string([]byte(after)[:3]) != "\xef\xbb\xbf" {
		t.Fatalf("AC#3 violated: the BOM is gone, file now starts % x", []byte(after)[:min(3, len(after))])
	}
	if after != utf8BOM+"title: 报表\nstatus: final\n" {
		t.Fatalf("the edit changed more than the matched bytes:\n got=%q\nwant=%q", after, utf8BOM+"title: 报表\nstatus: final\n")
	}
	if n := strings.Count(after, utf8BOM); n != 1 {
		t.Errorf("exactly one BOM byte sequence may exist, got %d (a re-encode would duplicate or move it)", n)
	}
	noStagingFilesLeft(t, root)
	t.Logf("AC#3 BOM 正向读数：%d 字节 → %d 字节，前三字节仍为 EF BB BF，BOM 出现 %d 次，文案=%q",
		len(original), len(after), strings.Count(after, utf8BOM), out.Text)
}

// ---------------------------------------------------------------------------
// AC#3 reverse: CRLF file, LF-spelled old = 0 hits (and WHY it refuses)
// ---------------------------------------------------------------------------

// TestFSEditOnACRLFFileRefusesAnLFSpelledOld is AC#3's second named case. The
// dispatch asks for more than "it fails loudly": it asks whether the failure is
// a CORRECT REFUSAL or the file simply not being readable.
//
// The two are told apart here by two readings on the SAME bytes:
//  1. the refusal names the file's byte length ("文件 23 字节"), which the code
//     can only print after a successful read (fs_edit.go:171-176 counts inside
//     `content`), and the read-failure wording ("读取目标失败") must be absent;
//  2. immediately after, a line-break-free old on the same file HITS and lands,
//     which is impossible if the file were unreadable.
func TestFSEditOnACRLFFileRefusesAnLFSpelledOld(t *testing.T) {
	root := sealableTempDir124(t)
	const original = "alpha\r\nbravo\r\ncharlie\r\n"
	target := filepath.Join(root, "crlf.txt")
	if err := os.WriteFile(target, []byte(original), 0o600); err != nil {
		t.Fatal(err)
	}
	crlf, cr, lf := endings(original)
	if crlf != 3 || cr != 0 || lf != 0 {
		t.Fatalf("fixture is not a pure CRLF file: crlf=%d bareCR=%d bareLF=%d", crlf, cr, lf)
	}
	b, _ := fsEditBridge(t, mustCanonical(t, root))

	// (1) the LF spelling of the same two lines: 0 hits.
	out, err := b.Execute(t.Context(), req("fs.edit", editArgs(t, target,
		map[string]any{"old": "alpha\nbravo", "new": "ALPHA\nBRAVO"})))
	if err != nil {
		t.Fatal(err)
	}
	assertRefusal(t, out, "0 命中", fmt.Sprintf("文件 %d 字节", len(original)))
	if strings.Contains(out.Text, "读取目标失败") {
		t.Errorf("this must be the exact-match refusal, not a read failure: %q", out.Text)
	}
	assertUntouched(t, target, original)
	// A normalizing locator WOULD have found this: the cost is real, not a typo.
	if n := strings.Count(lfView(original), "alpha\nbravo"); n != 1 {
		t.Fatalf("expected 1 hit under LF unification to justify registering this shape, got %d", n)
	}

	// (2) the control that settles 正确拒绝 vs 根本读不到: same file, one call
	// later, a line-break-free old hits and lands.
	out2, err := b.Execute(t.Context(), req("fs.edit", editArgs(t, target,
		map[string]any{"old": "bravo", "new": "BRAVO"})))
	if err != nil || out2.IsError {
		t.Fatalf("the control edit must land on these same bytes: out=%+v err=%v", out2, err)
	}
	const want = "alpha\r\nBRAVO\r\ncharlie\r\n"
	if got := readString(t, target); got != want {
		t.Fatalf("control edit changed the wrong bytes:\n got=%q\nwant=%q", got, want)
	}
	t.Logf("AC#3 CRLF 反向读数：%d 字节的纯 CRLF 文件里 LF 写法的 old 被响亮拒绝（文案含文件长度 %d，"+
		"证明读到了整文件），同一文件上不含换行的 old 立即命中并落盘 ⇒ 这是正确拒绝，不是读不到",
		len(original), len(original))
}

// ---------------------------------------------------------------------------
// AC#3 forward: CRLF file, CRLF-spelled old = must hit, and stay CRLF
// ---------------------------------------------------------------------------

// TestFSEditOnACRLFFileAppliesACRLFSpelledOldAndStaysCRLF is the other half of
// the ticket's "CRLF 文件改完还是 CRLF": the edit must land, and no line ending
// outside the matched region may be rewritten. The inserted line is spelled
// CRLF on purpose - that is what "the model matched the file's convention"
// looks like, and the file must come out pure.
func TestFSEditOnACRLFFileAppliesACRLFSpelledOldAndStaysCRLF(t *testing.T) {
	root := sealableTempDir124(t)
	const original = "alpha\r\nbravo\r\ncharlie\r\n"
	target := filepath.Join(root, "crlf-hit.txt")
	if err := os.WriteFile(target, []byte(original), 0o600); err != nil {
		t.Fatal(err)
	}
	b, _ := fsEditBridge(t, mustCanonical(t, root))

	out, err := b.Execute(t.Context(), req("fs.edit", editArgs(t, target,
		map[string]any{"old": "alpha\r\nbravo", "new": "ALPHA\r\nALIAS\r\nBRAVO"})))
	if err != nil || out.IsError {
		t.Fatalf("a CRLF-spelled old on a CRLF file must hit: out=%+v err=%v", out, err)
	}
	const want = "ALPHA\r\nALIAS\r\nBRAVO\r\ncharlie\r\n"
	after := readString(t, target)
	if after != want {
		t.Fatalf("the CRLF edit did not land as expected:\n got=%q\nwant=%q", after, want)
	}
	crlf, cr, lf := endings(after)
	if cr != 0 || lf != 0 {
		t.Fatalf("AC#3 violated: the file is no longer pure CRLF (crlf=%d bareCR=%d bareLF=%d)", crlf, cr, lf)
	}
	if crlf != 4 {
		t.Errorf("expected 4 CRLF pairs after adding one line, got %d", crlf)
	}
	noStagingFilesLeft(t, root)
	t.Logf("AC#3 CRLF 正向读数：%d 字节 → %d 字节，CRLF 对 %d→%d，裸 CR/裸 LF 均为 0（行尾未被改写），文案=%q",
		len(original), len(after), 3, crlf, out.Text)
}

// ---------------------------------------------------------------------------
// AC#3 clause 3: an LF file's endings must not be touched, either
// ---------------------------------------------------------------------------

// TestFSEditDoesNotRewriteAnLFFilesLineEndings is the ticket's third clause
// ("不许把 LF 文件的行尾改掉") in both directions: an LF file edited in LF stays
// free of CR, and a CRLF-spelled old on an LF file is the same loud 0-hit
// refusal - the asymmetry cannot run one way only.
func TestFSEditDoesNotRewriteAnLFFilesLineEndings(t *testing.T) {
	root := sealableTempDir124(t)
	const original = "alpha\nbravo\ncharlie\n"
	target := filepath.Join(root, "lf.txt")
	if err := os.WriteFile(target, []byte(original), 0o600); err != nil {
		t.Fatal(err)
	}
	b, _ := fsEditBridge(t, mustCanonical(t, root))

	// CRLF-spelled old on an LF file: refused, untouched, and the refusal is the
	// exact-match one.
	out, err := b.Execute(t.Context(), req("fs.edit", editArgs(t, target,
		map[string]any{"old": "alpha\r\nbravo", "new": "ALPHA\r\nBRAVO"})))
	if err != nil {
		t.Fatal(err)
	}
	assertRefusal(t, out, "0 命中")
	assertUntouched(t, target, original)

	// LF-spelled batch on the LF file: lands, and not one CR appears.
	out2, err := b.Execute(t.Context(), req("fs.edit", editArgs(t, target,
		map[string]any{"old": "bravo", "new": "bravo\nbravo2"})))
	if err != nil || out2.IsError {
		t.Fatalf("out=%+v err=%v", out2, err)
	}
	after := readString(t, target)
	const want = "alpha\nbravo\nbravo2\ncharlie\n"
	if after != want {
		t.Fatalf("the LF edit changed bytes outside the match:\n got=%q\nwant=%q", after, want)
	}
	if strings.ContainsRune(after, '\r') {
		t.Fatalf("AC#3 violated: CR appeared in an LF file: %q", after)
	}
	if crlf, cr, lf := endings(after); crlf != 0 || cr != 0 || lf != 4 {
		t.Errorf("expected 4 bare LF and no CR/CRLF, got crlf=%d bareCR=%d bareLF=%d", crlf, cr, lf)
	}
	t.Logf("AC#3 LF 双向读数：%d 字节 → %d 字节，裸 LF %d 条、CR 0 枚（LF 文件的行尾没被改）",
		len(original), len(after), strings.Count(after, "\n"))
}

// ---------------------------------------------------------------------------
// AC#3 forensics: what NOT implementing 票面 :27 costs, in readings
// ---------------------------------------------------------------------------

// TestFSEditLineEndingForensicsRefusesRealSpellings is the取证 half of AC#3, and
// it is取证 only: no normalization was implemented for it. Each row is a call a
// model actually makes once the file it read came from a Windows checkout - the
// old is byte-correct in LF, the file is CRLF (or CR, or mixed).
//
// Every row asserts two things so the number means something:
//   - today: a loud 0-hit refusal and an untouched file (the safe side of the
//     failure - it costs a round trip, never data);
//   - under 票面 :27's unification: the same old would have matched EXACTLY once,
//     i.e. it was a good call that is being refused. A row that would not match
//     after unification is not a cost, so it fails the case.
//
// The blast radius is therefore "N refusals of edits that were right", and it
// stays reversible. That is the argument FOR not building :27 today, and it is
// registered for the orchestrator to decide - not implemented here.
func TestFSEditLineEndingForensicsRefusesRealSpellings(t *testing.T) {
	cases := []struct {
		name     string
		file     string
		old      string
		wouldHit int
	}{
		{
			name:     "crlf_file_lf_spelled_old",
			file:     "alpha\r\nbravo\r\ncharlie\r\n",
			old:      "alpha\nbravo",
			wouldHit: 1,
		},
		{
			name:     "cr_only_file_lf_spelled_old",
			file:     "alpha\rbravo\rcharlie\r",
			old:      "alpha\nbravo",
			wouldHit: 1,
		},
		{
			name:     "mixed_file_old_over_the_crlf_line_lf_spelled",
			file:     "alpha\r\nbravo\ncharlie\r\n",
			old:      "alpha\nbravo",
			wouldHit: 1,
		},
		{
			name:     "mixed_file_old_over_the_lf_line_crlf_spelled",
			file:     "alpha\r\nbravo\ncharlie\r\n",
			old:      "bravo\r\ncharlie",
			wouldHit: 1,
		},
		{
			name:     "lf_file_crlf_spelled_old",
			file:     "alpha\nbravo\ncharlie\n",
			old:      "alpha\r\nbravo",
			wouldHit: 1,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := sealableTempDir124(t)
			target := filepath.Join(root, "shape.txt")
			if err := os.WriteFile(target, []byte(tc.file), 0o600); err != nil {
				t.Fatal(err)
			}
			if n := strings.Count(lfView(tc.file), lfView(tc.old)); n != tc.wouldHit {
				t.Fatalf("this row is not a cost: under LF unification old hits %d times, want %d", n, tc.wouldHit)
			}
			b, _ := fsEditBridge(t, mustCanonical(t, root))

			out, err := b.Execute(t.Context(), req("fs.edit", editArgs(t, target,
				map[string]any{"old": tc.old, "new": "REPLACED"})))
			if err != nil {
				t.Fatal(err)
			}
			assertRefusal(t, out, "0 命中", "must match exactly including all whitespace and newlines")
			assertUntouched(t, target, tc.file)
			t.Logf("归一取证：%d 字节文件（行尾形 %s）被 0 命中拒掉一枚归一后可命中 %d 次的 old；"+
				"落盘前后一字未变", len(tc.file), describeEndings(tc.file), tc.wouldHit)
		})
	}
}

// TestFSEditSilentShapeChangesLandToday is the other side of the same forensic
// cell, and the side that was NOT safe.
//
// IT NO LONGER PINS WHAT ITS NAME SAYS. The function name is r2's and the test
// roster is a stability key (dispatch r3 门禁: "缺那一侧必须为 0"), so the name
// stays; the assertions flipped in AC#3b (r3). r2 registered these two shapes as
// readings of defects, with tripwires in place - "the answer now signals the
// encoding change; update this case's wording" - and r3 is the cell those
// tripwires were waiting for. Both subtests now pin the OPPOSITE: each shape is a
// hard refusal, nothing lands, and the receipt says so. The forward direction of
// each shape (an old copied after the BOM; a new spelled in the file's own
// endings) is exercised by the control subtests of the same shape in
// fs_edit_ac3b_test.go, so this cell cannot be satisfied by a tool that refuses
// everything.
func TestFSEditSilentShapeChangesLandToday(t *testing.T) {
	t.Run("an_old_that_carries_the_bom_is_refused_r3", func(t *testing.T) {
		root := sealableTempDir124(t)
		original := utf8BOM + "alpha\nbeta\n"
		target := filepath.Join(root, "bom-zero.txt")
		if err := os.WriteFile(target, []byte(original), 0o600); err != nil {
			t.Fatal(err)
		}
		b, _ := fsEditBridge(t, mustCanonical(t, root))

		// The same call r2 measured landing at 14 bytes -> 11 with zero signal: a
		// model that copies the first line INCLUDING the BOM's three bytes matches
		// at offset 0, and the marker sits inside the replaced region.
		out, err := b.Execute(t.Context(), req("fs.edit", editArgs(t, target,
			map[string]any{"old": utf8BOM + "alpha", "new": "ALPHA"})))
		if err != nil {
			t.Fatal(err)
		}
		assertRefusal(t, out, "UTF-8 BOM")
		if after := readString(t, target); after != original {
			t.Fatalf("AC#3b① violated: the BOM-swallowing edit moved bytes:\n before=%q\n  after=%q", original, after)
		}
		if !strings.HasPrefix(readString(t, target), utf8BOM) {
			t.Fatal("the BOM is gone - AC#3b① regressed")
		}
		t.Logf("AC#3b① 翻转读数（r2 的缺陷读数①在此变红）：%d 字节带 BOM 的文件，old 从偏移 0 起并含 BOM ⇒ "+
			"现在是硬拒、一字未落，前三字节仍 EF BB BF，IsError=%v ErrorClass=%q",
			len(original), out.IsError, out.ErrorClass)
	})

	t.Run("an_lf_spelled_new_in_a_crlf_file_is_refused_r3", func(t *testing.T) {
		root := sealableTempDir124(t)
		const original = "alpha\r\nbravo\r\ncharlie\r\n"
		target := filepath.Join(root, "crlf-mixed.txt")
		if err := os.WriteFile(target, []byte(original), 0o600); err != nil {
			t.Fatal(err)
		}
		b, _ := fsEditBridge(t, mustCanonical(t, root))

		// r2's defect reading ②: the old spans no break, the new brings bare LFs,
		// and the file came out 3 CRLF + 1 bare LF with an answer that mentioned
		// nothing. That is now the refusal's own wording.
		out, err := b.Execute(t.Context(), req("fs.edit", editArgs(t, target,
			map[string]any{"old": "bravo", "new": "bravo\nbravo2"})))
		if err != nil {
			t.Fatal(err)
		}
		assertRefusal(t, out, "行尾约定", "裸 LF 1 条")
		if got := readString(t, target); got != original {
			t.Fatalf("AC#3b② violated: the mixed-ending edit moved bytes:\n before=%q\n  after=%q", original, got)
		}
		if crlf, cr, lf := endings(readString(t, target)); crlf != 3 || cr != 0 || lf != 0 {
			t.Fatalf("AC#3b② violated: the file is no longer pure CRLF: %d/%d/%d", crlf, cr, lf)
		}
		t.Logf("AC#3b② 翻转读数（r2 的缺陷读数②在此变红）：纯 CRLF 文件里一枚 LF 写法的 new ⇒ 硬拒，"+
			"落盘前后仍 3 组 CRLF / 0 裸 LF，IsError=%v ErrorClass=%q", out.IsError, out.ErrorClass)
	})
}

// ---------------------------------------------------------------------------
// AC#4: kill the write mid-flight
// ---------------------------------------------------------------------------

// TestFSEditKilledMidWriteLeavesTheTargetByteIdentical is AC#4 (a) and (b): the
// same seam TestAtomicWriteKillsMidWrite uses for fs.write (Hooks.Kill at an
// exact byte boundary), now mounted on fs.edit through the existing bridge
// builder, on a batch that has ALREADY passed every check - locating, counting,
// overlap, size cap. So the only thing left between here and the disk is the
// writer itself.
func TestFSEditKilledMidWriteLeavesTheTargetByteIdentical(t *testing.T) {
	const original = "alpha\nbravo\ncharlie\ndelta\n" // 26 bytes
	var killed []string
	root := sealableTempDir124(t)
	target := filepath.Join(root, "keep.txt")
	if err := os.WriteFile(target, []byte(original), 0o600); err != nil {
		t.Fatal(err)
	}
	b, _ := fsEditBridgeDeps(t, mustCanonical(t, root), func(d *FSDeps) {
		d.WriteChunk = 8
		d.Hooks = Hooks{
			AtStep: func(step string) { killed = append(killed, step) },
			Kill: func(step string) error {
				if step == "write:16" {
					return errors.New("模拟进程在此刻被杀")
				}
				return nil
			},
		}
	})

	out, err := b.Execute(t.Context(), req("fs.edit", editArgs(t, target,
		map[string]any{"old": "alpha", "new": "ALPHA-LONGER"},
		map[string]any{"old": "charlie", "new": "CHARLIE-LONGER-THAN-BEFORE"})))
	if err != nil {
		t.Fatal(err)
	}
	// (a) THE HEADLINE, FIRST. 162-v1 measured the old order: under its M-4 (the
	// staging call taken out of fs.edit) this case went red at the boundary probe
	// below, and the byte-level assertion underneath it never ran - one
	// assertion's failure swallowed another's reading. The load claim of AC#4 is
	// "the target holds the pre-edit bytes", so that is the first thing read; the
	// probes after it only say HOW GOOD the reading is (the seam fired where we
	// named it), and they are all still here, not one weaker. Order matters for
	// the counterfactual above too: with the byte check first, a rig that breaks
	// atomicity reports damage, not a missing step name.
	if got := readString(t, target); got != original {
		t.Fatalf("AC#4 violated: after a mid-write kill the target holds %d bytes %q, want the original %q",
			len(got), got, original)
	}
	// (0) the seam actually fired at the byte boundary we named, and only there.
	if len(killed) == 0 || killed[0] != "create-temp" {
		t.Fatalf("expected the write to reach the staging boundaries, got %q", killed)
	}
	var sawWrite16 bool
	for _, s := range killed {
		if s == "write:16" {
			sawWrite16 = true
		}
	}
	if !sawWrite16 {
		t.Fatalf("the kill boundary was never reached, the case proves nothing: %q", killed)
	}

	// (b) today's actual shape is BOTH halves of the "或": no .wisp-tmp-* file
	// survives AND the ledger names the one that was removed.
	noStagingFilesLeft(t, root)
	if !out.IsError {
		t.Fatalf("a killed edit must report a failure: %+v", out)
	}
	if len(out.AppliedSteps) == 0 {
		t.Fatalf("the D31 ledger must say what landed before the stop: %+v", out)
	}
	var namedTemp, namedRemoval bool
	for _, s := range out.AppliedSteps {
		if strings.Contains(s, "创建临时文件") {
			namedTemp = true
		}
		if strings.Contains(s, "删除临时文件") && strings.Contains(s, "目标从头到尾未被改动") {
			namedRemoval = true
		}
	}
	if !namedTemp || !namedRemoval {
		t.Errorf("the ledger must name both the staging file and its removal (temp=%v removal=%v): %q",
			namedTemp, namedRemoval, out.AppliedSteps)
	}
	t.Logf("AC#4(a)(b) 读数：改前 %d 字节 → 杀点 write:16（边界在下一块落笔前触发，暂存件里此刻 8 字节）"+
		"→ 目标仍是 %d 字节原字节（字节断言现在跑在边界探针之前）；"+
		"目录内 .wisp-tmp-* 残件 0 枚，台账 %q", len(original), len(readString(t, target)), out.AppliedSteps)
}

// TestFSEditWithoutTheStagedWriterWouldCorruptTheTarget answers AC#4's load-bit:
// 摘掉 temp＋rename 这一味，是否存在一发损坏从此看不见？
//
// The experiment is a counterfactual, not an opinion: take the bytes fs.edit
// itself produced (control run, unmodified code, no hooks), and send them at the
// target DIRECTLY - the upstream call shape 票面 :28 refuses to copy - using
// production's OWN chunked writer and the SAME kill seam. The only variable is
// the missing medicine.
//
// Reading: yes, and it is invisible twice over. The truncation is not a
// half-written file the reader can spot as half-written (nothing marks it), and
// the tool that did it is no longer there to complain, so the next instrument
// that touches the path - fs.read - reports success on 8 bytes of a file that
// used to hold 26. That is AC#1's "no signal" class arriving from a different
// direction, and temp+rename is what closes it.
func TestFSEditWithoutTheStagedWriterWouldCorruptTheTarget(t *testing.T) {
	const original = "alpha\nbravo\ncharlie\ndelta\n" // 26 bytes
	root := sealableTempDir124(t)
	target := filepath.Join(root, "counterfactual.txt")
	if err := os.WriteFile(target, []byte(original), 0o600); err != nil {
		t.Fatal(err)
	}

	// Step 1 - control: the real fs.edit, no hooks, same batch. Whatever lands
	// here is the byte string the counterfactual will write directly.
	b, _ := fsEditBridge(t, mustCanonical(t, root))
	out, err := b.Execute(t.Context(), req("fs.edit", editArgs(t, target,
		map[string]any{"old": "alpha", "new": "ALPHA-LONGER"},
		map[string]any{"old": "charlie", "new": "CHARLIE-LONGER-THAN-BEFORE"})))
	if err != nil || out.IsError {
		t.Fatalf("the control edit must land: out=%+v err=%v", out, err)
	}
	intended := readString(t, target)
	if intended == original || len(intended) <= 16 {
		t.Fatalf("control reading unusable: %d bytes", len(intended))
	}

	// Step 2 - the counterfactual: same target, same bytes, same kill boundary,
	// no staging file. O_TRUNC is the whole point: the original is gone the
	// instant the writer opens, one byte of new content at a time.
	if err := os.WriteFile(target, []byte(original), 0o600); err != nil {
		t.Fatal(err)
	}
	d := FSDeps{
		WriteChunk: 8,
		Hooks: Hooks{Kill: func(step string) error {
			if step == "write:16" {
				return errors.New("模拟进程在此刻被杀")
			}
			return nil
		}},
	}
	se := d.ledger()
	f, oerr := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if oerr != nil {
		t.Fatal(oerr)
	}
	written, werr := se.writeAll(t.Context(), f, strings.NewReader(intended), nil, d.writeCap())
	closeErr := f.Close()
	if closeErr != nil {
		t.Errorf("close after the kill: %v", closeErr)
	}
	if werr == nil {
		t.Fatal("the counterfactual writer ran to completion - the seam never fired, so this case proves nothing")
	}

	// (c)-1: the damage is real and it is at the target, not in a temp file.
	// writeAll checks the boundary BEFORE writing the chunk whose arrival it
	// names, so a kill at "write:16" leaves exactly one 8-byte chunk on disk -
	// the same arithmetic the staged case above runs on.
	if written != 8 {
		t.Errorf("expected 8 bytes written before the stop, got %d", written)
	}
	got := readString(t, target)
	if got != intended[:8] {
		t.Fatalf("expected a truncated prefix of the new content, got %q", got)
	}
	if got == original {
		t.Fatal("the counterfactual left the target intact - temp+rename would then be worthless, re-read this case")
	}
	if len(got) >= len(original) {
		t.Fatalf("expected a target SMALLER than the original (silent truncation), got %d vs %d", len(got), len(original))
	}

	// (c)-2: and it is INVISIBLE to the next instrument. fs.read on the corrupted
	// file is a success carrying the corrupted bytes, no error class, no
	// truncation flag - the reader has nothing to compare against.
	b2, _ := fsEditBridge(t, mustCanonical(t, root))
	rd, rerr := b2.Execute(t.Context(), req("fs.read", pathArgsOf(t, target)))
	if rerr != nil || rd.IsError || rd.ErrorClass != "" {
		t.Fatalf("the claim under test is that a later read says nothing, got: out=%+v err=%v", rd, rerr)
	}
	if rd.Truncated {
		t.Errorf("Truncated=true would be a signal; this case pins today's silence: %+v", rd)
	}
	if rd.Text != got {
		t.Fatalf("fs.read must hand back exactly the corrupted bytes for the reading to mean anything: got %q want %q", rd.Text, got)
	}
	t.Logf("AC#4(c) 反面实验读数：同一批字节、同一道 write:16 杀点，摘掉 temp+rename 之后目标从 %d 字节变成 "+
		"%d 字节（%q），既不是原文也不是要落的内容；事后 fs.read 返回 IsError=false ErrorClass=\"\" "+
		"Truncated=false，把损坏原文照读照回 ⇒ 这发损坏从此看不见。装回 temp+rename（上一枚用例）则目标一字未动。",
		len(original), len(got), got)
}

// describeEndings names a blob's line-ending shape for a log line.
func describeEndings(s string) string {
	crlf, cr, lf := endings(s)
	switch {
	case crlf == 0 && cr == 0 && lf == 0:
		return "无行尾"
	case crlf > 0 && cr == 0 && lf == 0:
		return "纯 CRLF"
	case lf > 0 && crlf == 0 && cr == 0:
		return "纯 LF"
	case cr > 0 && crlf == 0 && lf == 0:
		return "纯 CR"
	default:
		return fmt.Sprintf("混合(CRLF=%d CR=%d LF=%d)", crlf, cr, lf)
	}
}
