//go:build windows

package tools

// Ticket 20 frame 7 (c)③ — the security core: a junction / 8.3 rejection must
// NEVER be presentable as a [fs] allowed_dirs first-use ask. Real fixtures
// only: a real NTFS junction (mklink /J) and a real 8.3 name
// (GetShortPathNameW), both taken from the shared harness in
// bridge_junction_windows_test.go — zero t.Skip, a missing prerequisite fails
// loudly (this repository's standing rule for these fixtures).
//
// WHY THE CALL IS REFUSED BY A GATE ANSWER OF REJECT in the junction case: with
// an allow-all gate the junction call is refused LATER, by the tool's own second
// C26 parse, and that refusal is a tool-class error — not a bridge refusal, and
// therefore not something an ask could ever be derived from. The refusal shape
// (the one the frame names: 「被拒时」) is produced by the L2 card being
// answered no, and that is what these cases drive.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
)

// Test20r2JunctionRefusalIsNeverAskable: the frame's own attack narrative — an
// injected document names a path that traverses a junction, and the refusal it
// produces must never read as "just out of scope, shall we authorize the
// directory?". The value is RefusalRedTeamDenied, and the ask predicate is
// false.
func Test20r2JunctionRefusalIsNeverAskable(t *testing.T) {
	f := newReparseFixture(t) // positive control inside: the junction really reads f.secret
	g := &gateSpy{approveAns: AnswerReject, approveWhy: "用户否决了这次执行"}
	b, _, dl := bridgeOver(t, f.allowed, false, nil, g)

	out, err := b.Execute(t.Context(), req("fs.read", pathArgsOf(t, f.viaJunction("victim.txt"))))
	if err != nil {
		t.Fatalf("拒绝不是宿主故障: %v", err)
	}
	if !out.IsError || strings.Contains(out.Text, f.secret) {
		t.Fatalf("junction 对面的字节进了结果: %+v", out)
	}
	d := dl.last()
	if d.Level != risk.L2 || !contains(d.RulesHit, risk.R2) {
		t.Fatalf("判定 = %+v, want L2 via R2（解析器一拒，R2 fail-closed）", d)
	}
	if got := string(d.Refusal); got != "redteam-denied" {
		t.Fatalf("junction 拒绝的机器可辨值 = %q, want %q", got, "redteam-denied")
	}
	if askable20r2(d) {
		t.Fatal("junction 拒绝带上了可问类值：注入文档就此能驱动「要不要授权这个目录」的询问 (c)③")
	}
	if d.Refusal == RefusalNone {
		t.Fatal("junction 拒绝的机器可辨值是空的：两类拒绝又回到只靠措辞区分（(a) 那一格）")
	}
}

// Test20r2ShortNameOfAnAListFileIsNeverAskable: the other shape ticket 18
// calls red-team. The 8.3 spelling resolves to the long A-list path, so R2 is
// satisfied (the file sits INSIDE an allowed root) and R3 denies it; the gate
// is allow-all on both routes, which is what makes the denial the judge's and
// not the human's. Its refusal class must be the red-team one.
func Test20r2ShortNameOfAnAListFileIsNeverAskable(t *testing.T) {
	homeRaw := tempRaw(t)
	t.Setenv("USERPROFILE", homeRaw)
	t.Setenv("HOME", homeRaw)
	if err := os.MkdirAll(filepath.Join(homeRaw, ".ssh"), 0o700); err != nil {
		t.Fatal(err)
	}
	secret := "id_ed25519 AAAC3NzaC1lZDI1NTE5AAAA 20r2 askable probe"
	long := filepath.Join(homeRaw, ".ssh", "id_ed25519_20r2_askable_probe")
	if err := os.WriteFile(long, []byte(secret), 0o600); err != nil {
		t.Fatal(err)
	}
	short := shortNameOf(t, long)
	home := mustCanonical(t, homeRaw)
	if mustCanonical(t, short) != mustCanonical(t, long) {
		t.Fatalf("短名 %q 与长名 %q 解析到了不同目标，用例前提不成立", short, mustCanonical(t, long))
	}

	g := &gateSpy{windowAns: AnswerAllow, approveAns: AnswerAllow}
	b, _, dl := bridgeOver(t, home, false, nil, g)
	out, err := b.Execute(t.Context(), req("fs.read", pathArgsOf(t, short)))
	if err != nil {
		t.Fatalf("Deny 是拒绝不是故障: %v", err)
	}
	if !out.IsError || strings.Contains(out.Text, secret) || out.ErrorClass != "permission_denied" {
		t.Fatalf("8.3 拼法绕过了 A 档拒绝或错分类: %+v", out)
	}
	if w, a := g.counts(); w != 0 || a != 0 {
		t.Fatalf("A 档绝不可走到等人批准/被问的卡：window=%d approval=%d", w, a)
	}
	d := dl.last()
	if d.Level != risk.Deny || !contains(d.RulesHit, risk.R3) {
		t.Fatalf("判定 = %+v, want Deny via R3", d)
	}
	if got := string(d.Refusal); got != "redteam-denied" {
		t.Fatalf("8.3 拒绝的机器可辨值 = %q, want %q", got, "redteam-denied")
	}
	if askable20r2(d) {
		t.Fatal("8.3 拒绝带上了可问类值 (c)③")
	}
}

// Test20r2ShortNameOfAPlainFileIsTheOrdinaryOutOfScopeClass registers the
// boundary this leg deliberately drew, with the ruling it follows rather than
// a preference: 8.3 is a SPELLING, not a refusal class. The committed
// TestShortNameSpellingGetsTheSameVerdictAsTheLongOne (same package) pins
// verbatim that "8.3 is deliberately NOT a resolver-refused class" and that an
// approved out-of-scope short-name read must read — so the ask class follows
// the TARGET: a plain file outside every root is out-of-scope whichever
// spelling named it, and a blacklisted target is red-team whichever spelling
// named it (the case above). If the owner wants the stricter reading — "any
// 8.3-spelled path is red-team" — it is a one-line change in refusal.go, and
// THIS case is the one that will ring.
func Test20r2ShortNameOfAPlainFileIsTheOrdinaryOutOfScopeClass(t *testing.T) {
	root := tempCanonical(t)
	outside := tempRaw(t)
	long := filepath.Join(outside, "outside_secret_data.txt")
	if err := os.WriteFile(long, []byte("plain bytes, not blacklisted"), 0o600); err != nil {
		t.Fatal(err)
	}
	short := shortNameOf(t, long)
	longCanon := mustCanonical(t, long)
	if mustCanonical(t, short) != longCanon {
		t.Fatalf("短名 %q 与长名 %q 解析到了不同目标，用例前提不成立", short, longCanon)
	}

	var classes []RefusalKind
	for _, spelling := range []string{short, longCanon} {
		g := &gateSpy{approveAns: AnswerReject, approveWhy: "用户否决了这次执行"}
		b, _, dl := bridgeOver(t, root, false, nil, g)
		out, err := b.Execute(t.Context(), req("fs.read", pathArgsOf(t, spelling)))
		if err != nil {
			t.Fatalf("拒绝不是宿主故障（拼法 %q）: %v", spelling, err)
		}
		if !out.IsError {
			t.Fatalf("越界读被放行（拼法 %q）: %+v", spelling, out)
		}
		classes = append(classes, dl.last().Refusal)
	}
	if classes[0] != classes[1] {
		t.Fatalf("两种拼法必须落进同一类拒绝：short=%q long=%q", classes[0], classes[1])
	}
	if got := string(classes[0]); got != "out-of-allowlist" {
		t.Fatalf("非黑名单越界文件的 8.3 拼法类 = %q, want %q（8.3 是拼法不是拒绝类，"+
			"见 TestShortNameSpellingGetsTheSameVerdictAsTheLongOne 的已裁定注释）", got, "out-of-allowlist")
	}
}
