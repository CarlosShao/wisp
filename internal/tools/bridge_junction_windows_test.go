//go:build windows

package tools

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/risk"
)

// Ticket 20 box 5 (AC#5): the BRIDGE layer's reparse-point / 8.3-short-name
// denial.
//
// Why this file exists even though ticket 18 already denies both artifacts:
// `internal/risk/pathresolver_junction_windows_test.go` proves the LOWER layer
// says no. It does not prove the choke point cannot swallow that no and carry
// on. The difference is exactly the A33② family ("the implementation exists,
// nobody upstream uses it") and the fake-green shape this project keeps
// registering, so these cases drive REAL calls through `Bridge.Execute` and
// assert that zero bytes landed on the far side of a REAL junction.
//
// No string mocks. Every artifact below is a real NTFS junction
// (`mklink /J`, which needs no administrator token, unlike a symlink) and a
// real 8.3 short name taken from `GetShortPathNameW`. If the machine cannot
// produce one, the case fails LOUDLY with its prerequisite spelled out; it
// never skips, because a skipped security assertion is the empty proof this
// repository registered against itself today (Go counts SKIP as ok).
//
// The adversary in every case is a gate that ANSWERS ALLOW on both routes: a
// human who waved the card through must still not be able to reach the far
// side of a junction, because the card cannot show what the resolver refused
// to resolve.

// ---------------------------------------------------------------------------
// real-artifact fixtures
// ---------------------------------------------------------------------------

// mkRealJunction creates a REAL NTFS junction. mklink is an internal cmd
// command, so cmd.exe has to be the one to run it; /J is the junction form,
// the only reparse point creatable without SeCreateSymbolicLinkPrivilege.
func mkRealJunction(t *testing.T, link, target string) {
	t.Helper()
	out, err := exec.Command("cmd", "/c", "mklink", "/J", link, target).CombinedOutput()
	if err != nil {
		t.Fatalf("前置条件缺失：造不出真 NTFS junction（mklink /J %s %s: %v: %s）。\n"+
			"需要：%s 所在卷支持 reparse point（NTFS）且当前用户能在临时目录建 junction。\n"+
			"本用例不许改成用字符串假装 junction。", link, target, err, out,
			filepath.VolumeName(link))
	}
	if _, err := os.Lstat(link); err != nil {
		t.Fatalf("mklink 报告成功但 %s 不存在: %v", link, err)
	}
	// Go's Lstat renders a junction as an irregular mode on this toolchain, so
	// "it is really a junction" is proven the only way that matters: the
	// fixture reads through it (see newReparseFixture's positive control).
}

// shortNameOf returns the REAL 8.3 spelling of an existing path.
func shortNameOf(t *testing.T, long string) string {
	t.Helper()
	p16, err := syscall.UTF16PtrFromString(long)
	if err != nil {
		t.Fatalf("UTF16PtrFromString(%q): %v", long, err)
	}
	n, err := syscall.GetShortPathName(p16, nil, 0)
	vol := filepath.VolumeName(long)
	if err != nil || n == 0 {
		t.Fatalf("前置条件缺失：%s 所在卷 %s 没有给出 8.3 短名（GetShortPathNameW: %v）。\n"+
			"需要该卷启用短名生成（查：%s 应为 \"8dot3name 设置为 0（启用）\"，"+
			"改它需要管理员，所以这条用例在关闭了 8.3 的卷上必须红并写明缺什么，不许 skip）。",
			long, vol, err, "fsutil 8dot3name query "+vol)
	}
	buf := make([]uint16, n)
	if n, err = syscall.GetShortPathName(p16, &buf[0], n); err != nil || n == 0 {
		t.Fatalf("GetShortPathNameW(%q) 第二遍: %v", long, err)
	}
	short := syscall.UTF16ToString(buf[:n])
	if strings.EqualFold(short, long) || !strings.Contains(strings.ToLower(filepath.Base(short)), "~") {
		t.Fatalf("前置条件缺失：%s 的短名回显是 %q，与长名无法区分（该文件名太短或卷未生成 8.3）。\n"+
			"这条用例的输入必须是与长名不同的真短名拼法，否则它什么也没测。", long, short)
	}
	return short
}

// dirSnapshot hashes every regular file directly under dir. "No byte landed in
// the target" is asserted as before == after over the whole directory, not as
// a hopeful stat of the one file the test happened to name.
func dirSnapshot(t *testing.T, dir string) map[string]string {
	t.Helper()
	des, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir(%s): %v", dir, err)
	}
	out := map[string]string{}
	for _, e := range des {
		if !e.Type().IsRegular() {
			continue
		}
		b, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			t.Fatalf("ReadFile(%s): %v", filepath.Join(dir, e.Name()), err)
		}
		sum := sha256.Sum256(b)
		out[e.Name()] = fmt.Sprintf("%d:%s", len(b), hex.EncodeToString(sum[:]))
	}
	return out
}

func sameSnapshot(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

// reparseFixture is the shape that has bitten before: the SOURCE sits inside an
// authorized root, the TARGET sits outside it. A prefix check on the spelled
// path waves this through; only refusing reparse traversal catches it
// (SPEC-06 §4).
type reparseFixture struct {
	allowed  string            // canonical [fs] allowed_dirs root
	outside  string            // canonical: NOT in the allowlist
	junction string            // <allowed>\jn -> <outside>, a real junction
	victim   string            // the real file on the far side
	secret   string            // its bytes
	before   map[string]string // snapshot of outside, taken after the fixture is live
}

func newReparseFixture(t *testing.T) *reparseFixture {
	t.Helper()
	allowed := tempCanonical(t)
	outsideRaw := tempRaw(t)
	f := &reparseFixture{
		allowed: allowed,
		outside: mustCanonical(t, outsideRaw),
		secret:  "只有穿过 junction 才读得到的内容：" + strings.Repeat("s3cr3t", 8),
	}
	f.victim = filepath.Join(outsideRaw, "victim.txt")
	if err := os.WriteFile(f.victim, []byte(f.secret), 0o600); err != nil {
		t.Fatal(err)
	}
	f.junction = filepath.Join(allowed, "jn")
	mkRealJunction(t, f.junction, outsideRaw)

	// POSITIVE CONTROL: the junction must really reach the bytes. Without this,
	// a "denied" verdict could just mean the fixture was broken.
	got, err := os.ReadFile(filepath.Join(f.junction, "victim.txt"))
	if err != nil || string(got) != f.secret {
		t.Fatalf("junction 没有真的通向目标（读 %s: %v），后面的拒绝就什么也没证明",
			filepath.Join(f.junction, "victim.txt"), err)
	}
	// And the LONG spelling of the same outside file is the plain D34 out-of-
	// allowlist row (L2), i.e. the harness is not rejecting everything.
	f.before = dirSnapshot(t, outsideRaw)
	return f
}

func (f *reparseFixture) viaJunction(name string) string {
	return filepath.Join(f.junction, name)
}

func (f *reparseFixture) outsideSnapshot(t *testing.T) map[string]string {
	t.Helper()
	return dirSnapshot(t, filepath.Dir(f.victim))
}

// bridgeOver wires the whole fs roster over a real one-root allowlist and the
// real C19 assessor, with the gate and (optionally) the real journal the case
// asks for. It is a local constructor rather than fsDepsBridge because the
// shared one books no rows, and half of what this file proves is forensics.
func bridgeOver(t *testing.T, allowed string, deleteOn bool, j agent.Journal,
	g Gate,
) (*Bridge, *gateSpy, *decLog) {
	t.Helper()
	spy, ok := g.(*gateSpy)
	if !ok {
		t.Fatal("these cases need a gateSpy so the routed decision can be read back")
	}
	dl := &decLog{}
	paths := NewPathCanonicalizer([]string{allowed}, nil)
	reg := NewRegistry()
	for _, e := range BuiltinFSEntries(FSDeps{Paths: paths, DeleteEnabled: deleteOn}) {
		if err := reg.Register(e); err != nil {
			t.Fatalf("Register(%s): %v", e.Tool.Name(), err)
		}
	}
	return New(Options{
		Registry: reg, Paths: paths, Gate: spy, Journal: j,
		Provenance: risk.NewProvenance(risk.ProvOptions{NoProbe: true}),
		OnDecision: dl.add, Logf: func(string, ...any) {},
	}), spy, dl
}

// allowAllBridge is the composition under test: real C26 canonicalizer over a
// real one-root allowlist, the real C19 assessor, the whole fs roster, and a
// gate that waves everything through.
func allowAllBridge(t *testing.T, allowed string, deleteOn bool) (*Bridge, *gateSpy, *decLog) {
	t.Helper()
	return bridgeOver(t, allowed, deleteOn, nil, &gateSpy{
		windowAns: AnswerAllow, windowWhy: "test gate allowed it",
		approveAns: AnswerAllow, approveWhy: "test gate approved it",
	})
}

// assertJunctionRefused is the shared teeth: refused, booked as an approved
// call that still did nothing, and not one byte changed on the far side.
func assertJunctionRefused(t *testing.T, g *gateSpy, log *decLog, f *reparseFixture,
	out agent.ToolOutcome, err error,
) {
	t.Helper()
	if err != nil {
		t.Fatalf("拒绝必须是判定不是宿主故障: err=%v", err)
	}
	if !out.IsError {
		t.Fatalf("桥穿过了一个真 junction（%d 字节返回）: text=%q", len(out.Text), out.Text)
	}
	if out.RiskLevel != "L2" {
		t.Errorf("risk_level = %q, want L2（R2 对无法规范化的路径 fail-closed）", out.RiskLevel)
	}
	if _, a := g.counts(); a != 1 {
		t.Errorf("want exactly one L2 approval route: approval=%d", a)
	}
	// The card WAS answered allow (that is what route() acts on), so the only
	// thing that stopped this call is the bridge's own second resolution inside
	// the tool. "无法解析" is the TOOL's wording, not the gate's: if the lower
	// layer's denial had been swallowed upstream, this call would either have
	// succeeded (bytes moved) or been reported as a gate rejection.
	if d := g.approvalDecision(); d.Level != risk.L2 {
		t.Fatalf("人看到的那张卡 = %+v, want L2（红队拒绝得先成为判定，才会到人面前）", d)
	}
	if !strings.Contains(out.Text, "无法解析") {
		t.Fatalf("拒绝不是来自工具层的第二次 C26 解析: %+v", out)
	}
	if d := log.last(); !contains(d.RulesHit, risk.R2) {
		t.Errorf("rules_hit = %v, want R2 in it", d.RulesHit)
	}
	if after := f.outsideSnapshot(t); !sameSnapshot(f.before, after) {
		t.Errorf("junction 目标目录的内容被改动了：before=%v after=%v", f.before, after)
	}
	noStagingFilesLeft(t, f.allowed)
	noStagingFilesLeft(t, filepath.Dir(f.victim))
}

// ---------------------------------------------------------------------------
// 1. the read route
// ---------------------------------------------------------------------------

// TestBridgeRefusesARealJunctionOnTheReadRoute is box 5's read half: fs.read
// and fs.list driven through Bridge.Execute with a path that sits inside an
// authorized root and lands outside it through a real junction. The gate
// answers ALLOW, so a pass here means the refusal is the bridge's own, not a
// side effect of nobody approving the call.
func TestBridgeRefusesARealJunctionOnTheReadRoute(t *testing.T) {
	f := newReparseFixture(t)

	for _, tc := range []struct {
		name string
		tool string
		args func(t *testing.T) string
	}{
		{"fs.read_through_a_real_junction", "fs.read", func(t *testing.T) string {
			return pathArgsOf(t, f.viaJunction("victim.txt"))
		}},
		{"fs.list_through_a_real_junction", "fs.list", func(t *testing.T) string {
			return pathArgsOf(t, f.junction)
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			b, g, log := allowAllBridge(t, f.allowed, false)
			out, err := b.Execute(t.Context(), req(tc.tool, tc.args(t)))
			assertJunctionRefused(t, g, log, f, out, err)
			// The content must not leak, not even inside an error string.
			if strings.Contains(out.Text, f.secret) {
				t.Fatalf("拒绝的返回文本里出现了目标内容: %q", out.Text)
			}
			if strings.Contains(out.Text, "victim.txt") && tc.tool == "fs.list" {
				t.Fatalf("fs.list 把 junction 对面的条目列出来了: %q", out.Text)
			}
		})
	}

	// Observation for box 7(a) (NOT its acceptance, and nothing here ticks it):
	// the red-team refusal does not share the wording of the plain
	// out-of-allowlist refusal, and it travels as R2's fail-closed branch.
	// Today the distinction lives in the REASON STRING only - there is no typed
	// reason value on Decision, which box 7(a) asks for.
	b, _, log := allowAllBridge(t, f.allowed, false)
	if _, err := b.Execute(t.Context(), req("fs.read", pathArgsOf(t, f.viaJunction("victim.txt")))); err != nil {
		t.Fatal(err)
	}
	d := log.last()
	if !strings.Contains(d.Reason, "无法规范化") {
		t.Fatalf("junction 拒绝的 reason 必须说明是解析失败: %q", d.Reason)
	}
	if strings.Contains(d.Reason, "目标路径在授权目录之外") {
		t.Fatalf("红队拒绝不得与'不在授权目录内'共用措辞（否则未来的询问流会被注入文档驱动）: %q", d.Reason)
	}
	// Contrast: the same bridge on the long spelling of an out-of-allowlist file
	// does produce that other wording, so the distinction above is real and not
	// two ways of saying the same thing.
	g2 := &gateSpy{}
	b2, _ := fsBridgeWith(t, nil, g2, f.allowed)
	if _, err := b2.Execute(t.Context(), req("fs.read", pathArgsOf(t, f.victim))); err != nil {
		t.Fatal(err)
	}
	if r := g2.approvalDecision().Reason; !strings.Contains(r, "目标路径在授权目录之外") {
		t.Errorf("越界读的 reason = %q, want 越界措辞（对照项）", r)
	}
}

// TestBridgeWritesNothingThroughARealJunction is box 5's write half: every
// fs entry that can touch the disk, called with a junction target and a gate
// that approves. The target directory is snapshotted by content hash, so
// "no byte landed" covers the overwrite case, the create case and the
// staging-file case at once.
func TestBridgeWritesNothingThroughARealJunction(t *testing.T) {
	cases := []struct {
		name     string
		tool     string
		deleteOn bool
		args     func(t *testing.T, f *reparseFixture) string
		extra    func(t *testing.T, f *reparseFixture)
	}{
		{
			name: "fs.write_over_the_junctioned_target", tool: "fs.write",
			args: func(t *testing.T, f *reparseFixture) string {
				return writeArgsOf(t, f.viaJunction("victim.txt"), "被 junction 换掉的内容")
			},
		},
		{
			name: "fs.write_a_new_file_through_the_junction", tool: "fs.write",
			args: func(t *testing.T, f *reparseFixture) string {
				return writeArgsOf(t, f.viaJunction("brand-new.txt"), "x")
			},
			extra: func(t *testing.T, f *reparseFixture) {
				if existsFile(filepath.Join(filepath.Dir(f.victim), "brand-new.txt")) {
					t.Error("junction 对面多出了一个新文件")
				}
			},
		},
		{
			name: "fs.trash_through_the_junction", tool: "fs.trash",
			args: func(t *testing.T, f *reparseFixture) string {
				return pathArgsOf(t, f.viaJunction("victim.txt"))
			},
			extra: func(t *testing.T, f *reparseFixture) {
				if !existsFile(f.victim) {
					t.Error("fs.trash 把 junction 对面的真文件收走了")
				}
			},
		},
		{
			name: "fs.move_through_the_junction", tool: "fs.move",
			args: func(t *testing.T, f *reparseFixture) string {
				return moveArgsOf(t, f.viaJunction("victim.txt"), filepath.Join(f.allowed, "moved.txt"))
			},
			extra: func(t *testing.T, f *reparseFixture) {
				if existsFile(filepath.Join(f.allowed, "moved.txt")) {
					t.Error("fs.move 把 junction 对面的文件搬进了授权根")
				}
				if !existsFile(f.victim) {
					t.Error("fs.move 把 junction 对面的源搬走了")
				}
			},
		},
		{
			name: "fs.delete_through_the_junction", tool: "fs.delete", deleteOn: true,
			args: func(t *testing.T, f *reparseFixture) string {
				return pathArgsOf(t, f.viaJunction("victim.txt"))
			},
			extra: func(t *testing.T, f *reparseFixture) {
				if !existsFile(f.victim) {
					t.Error("fs.delete 永久删掉了 junction 对面的文件")
				}
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := newReparseFixture(t)
			b, g, log := allowAllBridge(t, f.allowed, tc.deleteOn)
			out, err := b.Execute(t.Context(), req(tc.tool, tc.args(t, f)))
			assertJunctionRefused(t, g, log, f, out, err)
			if readString(t, f.victim) != f.secret {
				t.Fatal("目标内容被改写过")
			}
			if tc.extra != nil {
				tc.extra(t, f)
			}
		})
	}
}

// TestJunctionInsideAnAllowedRootCannotReachAnAListFile is the direction the
// blacklist actually has to worry about: a junction planted in an authorized
// directory whose target is an A-tier tree. The bridge must never hand the
// bytes over, and - because A tier is absolute - must not even open a card.
func TestJunctionInsideAnAllowedRootCannotReachAnAListFile(t *testing.T) {
	homeRaw := tempRaw(t)
	t.Setenv("USERPROFILE", homeRaw)
	t.Setenv("HOME", homeRaw)
	ssh := filepath.Join(homeRaw, ".ssh")
	if err := os.MkdirAll(ssh, 0o700); err != nil {
		t.Fatal(err)
	}
	secret := "OPENSSH-KEY-MATERIAL-DO-NOT-LEAK"
	key := filepath.Join(ssh, "id_ed25519_import_key")
	if err := os.WriteFile(key, []byte(secret), 0o600); err != nil {
		t.Fatal(err)
	}

	// The authorized root is a DIFFERENT directory, so only the junction can
	// reach the A-list tree from inside it.
	allowed := tempCanonical(t)
	link := filepath.Join(allowed, "into-profile")
	mkRealJunction(t, link, ssh)
	if got, err := os.ReadFile(filepath.Join(link, "id_ed25519_import_key")); err != nil ||
		string(got) != secret {
		t.Fatalf("阳性对照失败：junction 没通向 A 档文件（%v）", err)
	}

	b, g, log := allowAllBridge(t, allowed, true)
	out, err := b.Execute(t.Context(), req("fs.read", pathArgsOf(t, filepath.Join(link, "id_ed25519_import_key"))))
	if err != nil {
		t.Fatalf("拒绝不是宿主故障: %v", err)
	}
	if !out.IsError || strings.Contains(out.Text, secret) {
		t.Fatalf("穿过 junction 读到了 A 档文件: %+v", out)
	}
	// NOTE(d22 ban #8): 本用例发现的真实软处（原样钉住，不顺手改）：junction 挡在 A 档目标前面时，
	// 判定层给出的**不是** Deny，而是一张人可以点批准的 L2 卡——因为解析器拒绝
	// 规范化，R3 根本没拿到长路径，只能 fail-closed。卡面上目标只显示成
	// "<原样路径> (无法规范化: …)"，批准它的人看不到自己批准的是什么。
	// 真正拦住字节的仍然只是工具层的第二次解析（上面的 out.IsError）。
	// ⇒ 属安全判定（"红队拒绝该不该在判定阶段就不可批准"），本项目硬规矩要求
	// owner 先判放哪层，这里只登记不修。
	w, a := g.counts()
	if w != 0 || a != 1 {
		t.Fatalf("今日真实路由是 L2 审批（不是 Deny、也不是 L1）: window=%d approval=%d", w, a)
	}
	d := log.last()
	if d.Level != risk.L2 || !contains(d.RulesHit, risk.R2) || !contains(d.RulesHit, risk.R3) {
		t.Fatalf("判定 = %+v, want L2 via R2+R3（双双 fail-closed）", d)
	}
	if len(d.Paths) != 1 || !strings.Contains(d.Paths[0], "无法规范化") {
		t.Errorf("卡面上的目标 = %v, want 带'无法规范化'标记（今天的形状）", d.Paths)
	}
	if strings.Contains(d.Reason, "A 档") {
		t.Errorf("reason = %q：解析器一拒，A 档匹配根本没发生（这正是上面那段登记的成因）", d.Reason)
	}
}

// ---------------------------------------------------------------------------
// 2. the real 8.3 short name
// ---------------------------------------------------------------------------

// TestBridgeRefusesTheRealShortNameOfAnAListFile: the 8.3 spelling is the
// classic way to dodge a string-prefix blacklist. Here the A-list file sits
// INSIDE the allowed root, so R2 is satisfied and only R3 can catch it - and
// the input is a real GetShortPathNameW spelling, not a hand-typed "~1".
func TestBridgeRefusesTheRealShortNameOfAnAListFile(t *testing.T) {
	homeRaw := tempRaw(t)
	t.Setenv("USERPROFILE", homeRaw)
	t.Setenv("HOME", homeRaw)
	if err := os.MkdirAll(filepath.Join(homeRaw, ".ssh"), 0o700); err != nil {
		t.Fatal(err)
	}
	secret := "id_ed25519 AAAC3NzaC1lZDI1NTE5AAAA short-name-probe"
	long := filepath.Join(homeRaw, ".ssh", "id_ed25519_short_name_probe")
	if err := os.WriteFile(long, []byte(secret), 0o600); err != nil {
		t.Fatal(err)
	}
	short := shortNameOf(t, long)
	home := mustCanonical(t, homeRaw)

	// The short spelling must be a genuinely different input that C26 folds
	// onto the same canonical target; otherwise this test measures nothing.
	if mustCanonical(t, short) != mustCanonical(t, long) {
		t.Fatalf("短名 %q 与长名 %q 解析到了不同目标，用例前提不成立", short, mustCanonical(t, long))
	}

	b, g, log := allowAllBridge(t, home, true)
	out, err := b.Execute(t.Context(), req("fs.read", pathArgsOf(t, short)))
	if err != nil {
		t.Fatalf("Deny 是拒绝不是故障: %v", err)
	}
	if !out.IsError || strings.Contains(out.Text, secret) {
		t.Fatalf("8.3 拼法绕过了 A 档拒绝: %+v", out)
	}
	if out.ErrorClass != "permission_denied" {
		t.Errorf("error_class = %q, want permission_denied（A 档是判定，不是工具出错）", out.ErrorClass)
	}
	w, a := g.counts()
	if w != 0 || a != 0 {
		t.Fatalf("A 档绝不可走到等人批准的卡：window=%d approval=%d", w, a)
	}
	d := log.last()
	if d.Level != risk.Deny || !contains(d.RulesHit, risk.R3) {
		t.Fatalf("判定 = %+v, want Deny via R3", d)
	}
	// What a human WOULD have seen on a card must name the long path: a
	// shortened spelling cannot be allowed to hide the real target.
	if len(d.Paths) != 1 || !strings.EqualFold(d.Paths[0], mustCanonical(t, long)) {
		t.Errorf("决策里显示的路径 = %v, want 解析后的长路径 %s", d.Paths, mustCanonical(t, long))
	}
}

// TestShortNameSpellingGetsTheSameVerdictAsTheLongOne covers the other 8.3
// shape: a real short name of a file OUTSIDE every allowed root. The two
// spellings must land on the SAME L2/R2 row (a shortened string must not read
// as in-scope, and must not become a red-team class either), the card must
// name the long path, and the bytes must come from the file the card named.
func TestShortNameSpellingGetsTheSameVerdictAsTheLongOne(t *testing.T) {
	outsideRaw := tempRaw(t)
	allowed := tempCanonical(t)
	body := "越界内容：短名拼法也必须折叠到同一个真目标上"
	long := filepath.Join(outsideRaw, "outside_secret_data.txt")
	if err := os.WriteFile(long, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	short := shortNameOf(t, long)
	wantCanon := mustCanonical(t, long)

	// An unanswered gate refuses, so this half also proves nothing was read.
	// Both spellings run the identical expectation: that IS the "same verdict"
	// claim, asserted per spelling rather than compared after the fact.
	for _, spelling := range []string{short, long} {
		g := &gateSpy{} // empty answers -> AnswerReject on both routes
		b, _, _ := bridgeOver(t, allowed, false, nil, g)
		out, err := b.Execute(t.Context(), req("fs.read", pathArgsOf(t, spelling)))
		if err != nil {
			t.Fatalf("拒绝不是宿主故障（拼法 %q）: %v", spelling, err)
		}
		if !out.IsError || strings.Contains(out.Text, body) {
			t.Fatalf("越界读被放行（拼法 %q）: %+v", spelling, out)
		}
		d := g.approvalDecision()
		if d.Level != risk.L2 || !contains(d.RulesHit, risk.R2) {
			t.Fatalf("拼法 %q 判定 = %+v, want L2 via R2", spelling, d)
		}
		if len(d.Paths) != 1 || !strings.EqualFold(d.Paths[0], wantCanon) {
			t.Errorf("卡面上的路径 = %v, want 长路径 %s（短名不许替真目标打掩护）", d.Paths, wantCanon)
		}
		if !strings.Contains(d.Reason, "目标路径在授权目录之外") {
			t.Errorf("拼法 %q 的 reason = %q, want 越界措辞（8.3 不是红队门）", spelling, d.Reason)
		}
	}

	// 8.3 is deliberately NOT a resolver-refused class: once a human approves an
	// out-of-scope read it must read. What must never happen is the card naming
	// one file while the bytes come from another.
	t.Run("an_approved_out_of_scope_short_name_reads_the_file_the_card_named", func(t *testing.T) {
		b, g, _ := allowAllBridge(t, allowed, false)
		out, err := b.Execute(t.Context(), req("fs.read", pathArgsOf(t, short)))
		if err != nil {
			t.Fatal(err)
		}
		if out.IsError || out.Text != body {
			t.Fatalf("短名没有折叠到长名指向的那个文件: out=%+v", out)
		}
		d := g.approvalDecision()
		if d.Level != risk.L2 || len(d.Paths) != 1 ||
			!strings.EqualFold(d.Paths[0], wantCanon) {
			t.Fatalf("人批准的卡必须说清它批准了哪个文件: %+v", d)
		}
	})
}

// TestWavedThroughJunctionWriteBooksAnApprovedRowThatWroteNothing is box 5's
// forensic half: what the audit trail says when a human approves a call whose
// target the resolver refused. The bytes are safe (asserted), but the row
// reads decision=allow / outcome=error / error_class=tool - "the tool said no",
// not "the policy said no". Recorded as today's real shape, with the SQLite row
// as the witness rather than a comment: whether a red-team refusal deserves its
// own disposition (observe.ClassPermissionDenied already exists) is a security
// judgement and belongs to the owner, not to this file.
func TestWavedThroughJunctionWriteBooksAnApprovedRowThatWroteNothing(t *testing.T) {
	f := newReparseFixture(t)
	store := openStore(t, filepath.Join(tempRaw(t), "data"))
	mustStartTask(t, store, "task-1")
	b, g, _ := bridgeOver(t, f.allowed, false, store,
		&gateSpy{windowAns: AnswerAllow, approveAns: AnswerAllow})

	out, err := b.Execute(t.Context(), req("fs.write",
		writeArgsOf(t, f.viaJunction("victim.txt"), "被批准了但仍不许落盘的内容")))
	if err != nil {
		t.Fatalf("拒绝不是宿主故障: %v", err)
	}
	if !out.IsError {
		t.Fatalf("桥在人批准之后把 junction 对面的文件写了: %+v", out)
	}
	if after := f.outsideSnapshot(t); !sameSnapshot(f.before, after) {
		t.Fatalf("junction 目标目录被改动: before=%v after=%v", f.before, after)
	}
	if _, a := g.counts(); a != 1 {
		t.Fatalf("want one approved L2 route: approval=%d", a)
	}

	rows, err := store.ListToolCallsByTask(t.Context(), "task-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("booked %d rows, want 1: %+v", len(rows), rows)
	}
	r := rows[0]
	if r.RiskLevel != "L2" || r.Decision != agent.DecisionAllow ||
		r.Outcome != agent.OutcomeError || r.ErrorClass != "tool" {
		t.Fatalf("审计行的真实形状 = risk=%q decision=%q outcome=%q class=%q, %+v",
			r.RiskLevel, r.Decision, r.Outcome, r.ErrorClass, r)
	}
	if !strings.Contains(r.ArgsJSON, "jn") {
		t.Errorf("args_json 里必须留着那次调用的原样路径: %s", r.ArgsJSON)
	}
}
