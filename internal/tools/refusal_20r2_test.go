package tools

// Ticket 20 frame 7 — the [fs] allowed_dirs first-use ask flow's minimal
// closed subset. This file carries items (a), (b) and (c) on every platform;
// the junction/8.3 fixtures live in refusal_20r2_windows_test.go.
//
// The three assertions (frame text):
//
//	① (a) 的两类拒绝原因可区分  — Test20r2RefusalKindsAreDistinct +
//	   Test20r2DenyIsRedTeamNeverAskable
//	② 批准 → 下一次同路径调用成功，拒绝 → 仍被拒 —
//	   Test20r2FirstUseAskFlowGoesOnlyThroughConfig +
//	   Test20r2DirectionEngineAdjudicatesFileLevelAdditions
//	③ junction / 8.3 拒绝绝不产生询问 — the windows file (real fixtures) plus
//	   the askable20r2 predicate every case in both files asserts with
//
// AND WHAT "下一次" MEANS HERE, so nobody reads more into it than the repo
// built: same-run widening is deliberately NOT live (cmd/wisp/config_reload.go
// names it: "swapping rt.paths mid-run would need an exported reset on
// internal/tools, which this leg may not mint (no new exported names)"), and
// ticket 201's card text promises 「下一次启动生效」. So the approval leg
// asserts BOTH halves honestly: the bridge that already refused keeps refusing
// (no silent in-run widening — i.e. no 「这次先放行」), and the NEXT assembly —
// the shape a restart produces — reads the persisted entry and lets the same
// path through.

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/risk"
)

// askable20r2 is the ask rule as a predicate: a first-use ask may ONLY be
// derived from RefusalOutOfAllowlist (refusal.go's contract). Every case in
// both files states its expectation through this function, so the security
// core (c)③ is asserted as a decision, not as a comment.
func askable20r2(d Decision) bool { return d.Refusal == RefusalOutOfAllowlist }

// r20Rig is one composed bridge plus the two sinks a refusal is read from: the
// verdict the host receives (OnDecision) and the audit line (Logf).
type r20Rig struct {
	b    *Bridge
	dec  *decLog
	logs []string
}

func r20Bridge(t *testing.T, allowed []string, g Gate) *r20Rig {
	t.Helper()
	if g == nil {
		g = NoGate{}
	}
	rig := &r20Rig{dec: &decLog{}}
	paths := NewPathCanonicalizer(allowed, nil)
	reg := NewRegistry()
	for _, e := range BuiltinFSEntries(FSDeps{Paths: paths}) {
		if err := reg.Register(e); err != nil {
			t.Fatalf("Register(%s): %v", e.Tool.Name(), err)
		}
	}
	rig.b = New(Options{
		Registry:   reg,
		Paths:      paths,
		Gate:       g,
		Provenance: risk.NewProvenance(risk.ProvOptions{NoProbe: true}),
		OnDecision: rig.dec.add,
		Logf: func(format string, args ...any) {
			rig.logs = append(rig.logs, fmt.Sprintf(format, args...))
		},
	})
	return rig
}

func (r *r20Rig) audit() string { return strings.Join(r.logs, "\n") }

// read is one fs.read through this rig: allowed plus the refusal verdict.
func (r *r20Rig) read(t *testing.T, target string) (bool, Decision) {
	t.Helper()
	out, err := r.b.Execute(t.Context(), req("fs.read", pathArgsOf(t, target)))
	if err != nil {
		t.Fatalf("拒绝不是宿主故障: %v", err)
	}
	return !out.IsError, r.dec.last()
}

// Test20r2RefusalKindsAreDistinct is assertion ① / item (a): a refusal that is
// merely out of scope and an absolute denial are two MACHINE-READABLE values
// rather than two wordings, and the out-of-scope one is the askable class.
func Test20r2RefusalKindsAreDistinct(t *testing.T) {
	root := tempCanonical(t)
	outside := tempCanonical(t)
	target := writeUnder(t, outside, "wanted.txt", "the bytes at stake")
	if strings.HasPrefix(strings.ToLower(target), strings.ToLower(root)+`\`) {
		t.Fatalf("fixture broken: %s is under %s; the out-of-scope leg measures nothing", target, root)
	}

	g := &gateSpy{approveAns: AnswerReject, approveWhy: "用户否决了这次执行"}
	rig := r20Bridge(t, []string{root}, g)
	allowed, d := rig.read(t, target)
	if allowed {
		t.Fatalf("越界读被放行: allowed=%v", allowed)
	}

	// The nail (a) asks for: distinct machine values, not shared words. Both
	// halves are asserted, because either edit that merges them has to ring.
	if string(RefusalOutOfAllowlist) == string(RefusalRedTeamDenied) {
		t.Fatalf("两类拒绝共用同一个字符串 %q —— 被拒原因不再可机器区分（票 20 第 7 框 (a)）",
			string(RefusalOutOfAllowlist))
	}
	if got := string(d.Refusal); got != "out-of-allowlist" {
		t.Fatalf("越界拒绝的机器可辨值 = %q, want %q（措辞不参与判定）", got, "out-of-allowlist")
	}
	if !askable20r2(d) {
		t.Fatalf("越界拒绝必须是可问类（「要不要把这个目录加进来」唯一允许的来源），got %q", d.Refusal)
	}
	// The audit line is the second host-facing sink; it must carry the value,
	// not only the verdict the OnDecision callback saw.
	if !strings.Contains(rig.audit(), "refusal=out-of-allowlist") {
		t.Fatalf("审计行没带上机器可辨值:\n%s", rig.audit())
	}
	if !strings.Contains(d.Reason, "R2") {
		t.Fatalf("reason = %q, want the R2 verdict wording (unchanged by this ticket)", d.Reason)
	}
}

// Test20r2DenyIsRedTeamNeverAskable is assertion ①'s other half plus (c)③'s
// principle on a cross-platform fixture: an absolute Deny (R3 tier A) must
// NEVER be presentable as "authorize its directory". The 8.3/UNC/\\?\ spellings
// of A-list files land HERE (they resolve to the long path and are denied by
// classification — ticket 18 Cases 2-4), so this case is also the
// cross-platform anchor for the 「8.3 拒绝绝不产生询问」 limb.
func Test20r2DenyIsRedTeamNeverAskable(t *testing.T) {
	home := tempCanonical(t)
	t.Setenv("USERPROFILE", home)
	t.Setenv("HOME", home)
	if err := os.MkdirAll(filepath.Join(home, ".aws"), 0o755); err != nil {
		t.Fatal(err)
	}
	// Inside the allowed root on purpose: R2 is satisfied, so only R3 can
	// catch it — the shape an A-list dodge wants.
	target := writeUnder(t, home, filepath.Join(".aws", "credentials"), "AKIAIOSFODNN7EXAMPLE")

	g := &gateSpy{windowAns: AnswerAllow, approveAns: AnswerAllow}
	rig := r20Bridge(t, []string{home}, g)
	out, err := rig.b.Execute(t.Context(), req("fs.read", pathArgsOf(t, target)))
	if err != nil {
		t.Fatalf("Deny 是拒绝不是故障: %v", err)
	}
	if !out.IsError || out.ErrorClass != "permission_denied" {
		t.Fatalf("A 档读必须被拒: %+v", out)
	}
	if w, a := g.counts(); w != 0 || a != 0 {
		t.Fatalf("A 档绝不可走到等人批准/被问的卡：window=%d approval=%d", w, a)
	}
	d := rig.dec.last()
	if d.Level != risk.Deny || !contains(d.RulesHit, risk.R3) {
		t.Fatalf("判定 = %+v, want Deny via R3", d)
	}
	if got := string(d.Refusal); got != "redteam-denied" {
		t.Fatalf("A 档拒绝的机器可辨值 = %q, want %q（Deny 永远不可问）", got, "redteam-denied")
	}
	if askable20r2(d) {
		t.Fatal("A 档拒绝带上了可问类值：等于让红队门自己开口要社工点击 (c)③")
	}
}

// Test20r2FirstUseAskFlowGoesOnlyThroughConfig is (c)② on the existing 🔒
// channel — item (b) being "there is no other door". The approve action is
// Manager.AddAllowedDir, the one programmatic landing ticket 201 built and
// cmd/wisp/approval_always.go already drives; nothing here touches the
// canonicalizer.
func Test20r2FirstUseAskFlowGoesOnlyThroughConfig(t *testing.T) {
	cfgPath := filepath.Join(sealableTempDir124(t), "config.toml")
	base := "schema_version = 2\n\n[fs]\nallowed_dirs = []\n"
	if err := os.WriteFile(cfgPath, []byte(base), 0o600); err != nil {
		t.Fatal(err)
	}
	mgr, err := config.NewManager(cfgPath, nil)
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}

	outside := sealableTempCanonical124(t)
	target := writeUnder(t, outside, "wanted.txt", "the bytes the ask is about")
	offer := filepath.Dir(target) // the line an ask would offer: [fs] allowed_dirs += offer

	// build assembles the NEXT run's bridge the way cmd/wisp does: roots come
	// from config, never from a mutated canonicalizer.
	build := func() *r20Rig {
		return r20Bridge(t, mgr.Config().FS.AllowedDirs,
			&gateSpy{approveAns: AnswerReject, approveWhy: "用户否决了这次执行"})
	}
	var refusedRun *r20Rig

	t.Run("refused: the refusal is the askable class", func(t *testing.T) {
		refusedRun = build()
		allowed, d := refusedRun.read(t, target)
		if allowed || !askable20r2(d) {
			t.Fatalf("allowed=%v refusal=%q：空名单下这一发必须是可问类拒绝", allowed, d.Refusal)
		}
	})

	t.Run("denied: nothing is appended and the next call is still refused", func(t *testing.T) {
		if got := mgr.Config().FS.AllowedDirs; len(got) != 0 {
			t.Fatalf("没有人批准，名单却变了: %v", got)
		}
		allowed, d := build().read(t, target)
		if allowed || d.Refusal != RefusalOutOfAllowlist {
			t.Fatalf("拒绝后必须仍被拒且仍是可问类: allowed=%v refusal=%q", allowed, d.Refusal)
		}
	})

	t.Run("approved: the append lands via Manager.AddAllowedDir and the next assembly passes", func(t *testing.T) {
		if err := mgr.AddAllowedDir(offer); err != nil {
			t.Fatalf("AddAllowedDir（既有锁定通道）拒绝了这次追加: %v", err)
		}
		// Persisted, not memory-only: a fresh Manager over the same file — the
		// shape a restart produces — must read it back.
		raw, err := os.ReadFile(cfgPath)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(raw), offer) {
			t.Fatalf("批准后 config.toml 里没有这一行:\n%s", raw)
		}
		next, err := config.NewManager(cfgPath, nil)
		if err != nil {
			t.Fatalf("next-start NewManager: %v", err)
		}
		if !slices.Contains(next.Config().FS.AllowedDirs, offer) {
			t.Fatalf("下一次启动读不回来: %v", next.Config().FS.AllowedDirs)
		}
		// No silent in-run widening: the bridge that already refused still
		// refuses, and it must still be the askable class (the repo's named
		// residual, not a defect).
		if allowed, d := refusedRun.read(t, target); allowed || d.Refusal != RefusalOutOfAllowlist {
			t.Fatalf("同一枚桥的复核必须仍是被拒且可问（allowed=%v refusal=%q）——"+
				"名单一写盘就跑通，那是「这次先放行下次再说」，被 (b) 明令禁止", allowed, d.Refusal)
		}
		// And the next assembly lets the SAME path through, L0, no card.
		run := build()
		allowed, d := run.read(t, target)
		if !allowed {
			t.Fatalf("批准后下一次装配必须放行: refusal=%q", d.Refusal)
		}
		if d.Level != risk.L0 {
			t.Fatalf("放行了却不是 L0 直读: level=%v", d.Level)
		}
	})

	t.Run("scoped: the widening authorized one directory, not the machine", func(t *testing.T) {
		sibling := tempCanonical(t)
		sib := writeUnder(t, sibling, "other.txt", "越界内容")
		allowed, d := build().read(t, sib)
		if allowed || d.Refusal != RefusalOutOfAllowlist {
			t.Fatalf("被授权的是 %s，兄弟目录不该跟着放行: allowed=%v refusal=%q", offer, allowed, d.Refusal)
		}
	})
}

// Test20r2DirectionEngineAdjudicatesFileLevelAdditions is the reload half of
// item (b), in the precedent's shape (manager.go's fsDirection row +
// manager_test.go:162-225): a hand edit that loosens [fs] allowed_dirs is a
// locked-section change, and it takes effect only if ConfirmLocked approves.
// This test only ever changes the FILE — no in-process widen call exists.
func Test20r2DirectionEngineAdjudicatesFileLevelAdditions(t *testing.T) {
	cfgPath := filepath.Join(sealableTempDir124(t), "config.toml")
	base := "schema_version = 2\n\n[fs]\nallowed_dirs = []\n"
	if err := os.WriteFile(cfgPath, []byte(base), 0o600); err != nil {
		t.Fatal(err)
	}
	mgr, err := config.NewManager(cfgPath, nil)
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	outside := sealableTempCanonical124(t)
	target := writeUnder(t, outside, "wanted.txt", "reload-leg bytes")
	loosened := strings.Replace(base, "allowed_dirs = []",
		"allowed_dirs = ["+strconv.Quote(outside)+"]", 1)
	build := func() *r20Rig {
		return r20Bridge(t, mgr.Config().FS.AllowedDirs,
			&gateSpy{approveAns: AnswerReject, approveWhy: "用户否决了这次执行"})
	}

	// Rejected by the hook: still empty, same path still refused.
	var askedSection string
	var askedKeys []string
	mgr.ConfirmLocked = func(section string, keys []string) bool {
		askedSection, askedKeys = section, keys
		return false
	}
	if err := os.WriteFile(cfgPath, []byte(loosened), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := mgr.CheckAndReload(); err != nil {
		t.Fatalf("reload: %v", err)
	}
	if askedSection != "fs" || len(askedKeys) != 1 || askedKeys[0] != "fs.allowed_dirs" {
		t.Fatalf("confirm hook got section=%q keys=%v, want fs [fs.allowed_dirs]", askedSection, askedKeys)
	}
	if got := mgr.Config().FS.AllowedDirs; len(got) != 0 {
		t.Fatalf("拒绝的放宽必须保持旧值（空），got %v", got)
	}
	if allowed, d := build().read(t, target); allowed || !askable20r2(d) {
		t.Fatalf("拒绝后同路径必须仍被拒: allowed=%v refusal=%q", allowed, d.Refusal)
	}

	// Approved (a later edit of the SAME loosening): now it takes effect. The
	// 10ms gap is the reload tracker's own requirement — the poll is keyed on
	// mtime+size, and two writes inside one timestamp tick would read as "no
	// change" (manager_test.go's mutate helper sleeps for the same reason).
	time.Sleep(10 * time.Millisecond)
	mgr.ConfirmLocked = func(section string, keys []string) bool { return true }
	if err := os.WriteFile(cfgPath, []byte(loosened), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := mgr.CheckAndReload(); err != nil {
		t.Fatalf("reload: %v", err)
	}
	if got := mgr.Config().FS.AllowedDirs; len(got) != 1 {
		t.Fatalf("批准后才生效，got %v", got)
	}
	if allowed, d := build().read(t, target); !allowed {
		t.Fatalf("批准后同路径必须成功: refusal=%q", d.Refusal)
	}
}

// Test20r2NoSecondAuthorizationPath is item (b)'s structural half — the
// forbidden shapes, pinned so they cannot arrive quietly: nothing on the C26
// seam can widen what InAllowlist authorizes (construction is the only door),
// and this package never reaches into config to write an allowlist itself.
func Test20r2NoSecondAuthorizationPath(t *testing.T) {
	// 1. The exported method set is EXACTLY the read/judge surface plus the
	//    workspace NARROWING. A widening method (AddRoot/SetRoots/Reload...)
	//    fails this roster; a new read must extend it in the same commit.
	want := map[string]bool{
		"Canonicalize": true, "InAllowlist": true, "Roots": true,
		"UnusableRoots": true, "RewrittenRoots": true,
		"WorkspaceRoot": true, "ResolveWorkspace": true,
		"SetWorkspaceRoot": true, "ClearWorkspace": true,
	}
	typ := reflect.TypeOf((*PathCanonicalizer)(nil))
	got := map[string]bool{}
	for i := 0; i < typ.NumMethod(); i++ {
		got[typ.Method(i).Name] = true
	}
	for name := range want {
		if !got[name] {
			t.Errorf("PathCanonicalizer lost %s: the seam's roster changed", name)
		}
		delete(got, name)
	}
	for name := range got {
		t.Errorf("PathCanonicalizer grew %s. If it WIDENS what InAllowlist authorizes outside "+
			"NewPathCanonicalizer, that is the second authorization path ticket 20 frame 7 (b) forbids.", name)
	}

	// 2. No production file in this package imports internal/config: the
	//    widening lands through the host's Manager, so a tools-side write
	//    would be the birth of exactly the path (b) bans.
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	const importPath = `"github.com/CarlosShao/wisp/internal/config"`
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") || strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}
		b, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(b), importPath) {
			t.Errorf("internal/tools/%s imports internal/config: the allowlist write must stay on the "+
				"host's side of the bridge (ticket 20 frame 7 (b))", e.Name())
		}
	}
}
