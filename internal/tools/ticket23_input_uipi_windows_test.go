//go:build windows

package tools

// Ticket 23 AC#3's backend-residue legs, narrowed to what a no-GUI gate run
// can honestly claim.
//
// What the ticket log (2026-11-03) still lists open on AC#3 has three parts:
// a real fixture app, the real-machine UIPI refusal, and the three-option
// confirm card. The card is rendered by the native front-end (frontend/** is
// not this team's write face), and a fixture app requires starting and
// focusing a window - which this run is forbidden to do. What this file pins
// instead is the predicate and the binding key those missing pieces stand on,
// exercised against targets that need no window of their own:
//
//   - the UIPI gate (inputTarget -> ErrElevatedTarget BEFORE any byte is
//     sent, never a bypass attempt), driven through the seam the backend
//     already owns so the elevated answer is deterministic on any runner;
//   - that same gate letting a live, non-elevated target through with the
//     resolved hwnd and the resolved process name. The target is this test
//     process's own console window (GetConsoleWindow, kernel-owned, same
//     integrity level as us) when the run has one, and a fake hwnd that
//     resolves to our own pid when it does not - either way the pid under the
//     predicate is OURS, so the non-elevated answer is guaranteed and no
//     ambient focus can flip it;
//   - realProcessName producing exactly the D45 input.type binding key
//     (approval_grant.pattern, memory/schema.go:105): the image name, lower
//     case, no directory - the shape ticket 49's grant rows will be written
//     against.
//
// The PLAN.md:3234 sentence itself is asserted verbatim in
// Test23RefusalsBecomeModelActableErrors; these legs assert the machinery
// that fires it.

import (
	"errors"
	"strings"
	"testing"

	"golang.org/x/sys/windows"
)

// selfHwnd23 returns a live top-level window owned by this process without
// creating or focusing anything: the console window a `go test` run always has
// on a gate runner. A headless run (no console, rc=0 from Find) answers 0 and
// the caller falls back to the pid-stub shape.
var procGetConsoleWindow = windows.NewLazySystemDLL("user32.dll").NewProc("GetConsoleWindow")

func selfHwnd23(t *testing.T) uintptr {
	t.Helper()
	if procGetConsoleWindow.Find() != nil {
		return 0
	}
	h := callUintptr(procGetConsoleWindow)
	if h == 0 {
		t.Log("本进程没有控制台窗口（headless run）：live 目标走 pid 桩形")
	}
	return h
}

// Test23UIPIPredicateRefusesElevatedTargetBeforeAnyByte pins the AC#3 refusal
// order: the elevated answer comes from inputTarget, which runs before
// sendText/sendKeys are ever reached. The seam is b.processElevated, so the
// elevated answer is deterministic instead of depending on what happens to be
// running on the gate runner.
func Test23UIPIPredicateRefusesElevatedTargetBeforeAnyByte(t *testing.T) {
	sentText, sentKeys := false, false
	b := &systemBackend{
		sendText: func(uintptr, string) error { sentText = true; return nil },
		sendKeys: func(uintptr, string) error { sentKeys = true; return nil },
		// Every target is "elevated" here: the one no-bypass answer.
		processElevated: func(uint32) bool { return true },
	}
	// No hwnd, no title: the target is the foreground window. That is THIS
	// process on the gate runner - no window is created, focused or closed to
	// make it so.
	_, _, err := b.inputTarget(map[string]any{})
	if err == nil {
		t.Fatalf("提权目标没有在 inputTarget 被拒：注入会静默丢弃（UIPI 的老陷阱）")
	}
	if !errors.Is(err, ErrElevatedTarget) {
		t.Fatalf("拒绝错误是 %v，要 ErrElevatedTarget（PLAN.md:3234 的那一句）", err)
	}
	if sentText || sentKeys {
		t.Fatalf("拒绝之前就已经发了输入：sendText=%v sendKeys=%v（先注入后拒绝等于没拦）", sentText, sentKeys)
	}
	// And the tool layer turns it into model-actable IsError text carrying the
	// frozen sentence, with no bypass hint in it.
	d := SystemDeps{Bridge: &SystemBridge{Backend: b}}
	entries := BuiltinSystemEntries(d)
	var input Entry
	for _, e := range entries {
		if e.Tool.Name() == "input.type" {
			input = e
			break
		}
	}
	if input.Tool == nil {
		t.Fatalf("input.type 不在名册里")
	}
	res, err := input.Tool.Execute(t.Context(), []byte(`{"text":"hi"}`), nil)
	if err != nil {
		t.Fatalf("拒绝变成了宿主错误 %v：这应是模型可自纠的 IsError", err)
	}
	if !res.IsError || !strings.Contains(res.Text, ErrElevatedTarget.Error()) {
		t.Fatalf("IsError=%v Text=%q，要 IsError + ErrElevatedTarget 原句", res.IsError, res.Text)
	}
	// The refusal must CARRY the no-bypass commitment ("本工具不会尝试提权绕过"
	// is part of the frozen answer), and must not hand the model a workaround
	// (no runas/bypass/ShellExecute instructions). Naming UIPI or 提权 in a
	// negation is the honest answer; suggesting one is not.
	if !strings.Contains(res.Text, "不会尝试提权绕过") {
		t.Errorf("拒绝文本丢了「不尝试绕过」的承诺：%q", res.Text)
	}
	for _, banned := range []string{"runas", "bypass", "shellexecute"} {
		if strings.Contains(strings.ToLower(res.Text), banned) {
			t.Errorf("拒绝文本出现了疑似绕过指引 %q：%q", banned, res.Text)
		}
	}
}

// Test23UIPIPredicateLetsLiveNonElevatedTargetThrough is the other half of the
// same gate, against the REAL predicate (realProcessElevated: integrity
// levels, no seam stub): the live target's pid is our own process, so the
// non-elevated answer is guaranteed, and inputTarget must hand back the hwnd
// plus a resolved process name. If the predicate ever fails closed on a
// legitimate same-integrity target, every input.type call dies here.
func Test23UIPIPredicateLetsLiveNonElevatedTargetThrough(t *testing.T) {
	b := &systemBackend{processElevated: realProcessElevated}
	h := selfHwnd23(t)
	if h != 0 {
		got, name, err := b.inputTarget(map[string]any{"hwnd": int(h)})
		if err != nil {
			t.Fatalf("本进程自己的窗口被判提权：%v（谓词对同完整性目标必须放行）", err)
		}
		if got != h {
			t.Fatalf("inputTarget 返回的 hwnd=%#x，传入 %#x：目标被换走了", got, h)
		}
		if strings.TrimSpace(name) == "" {
			t.Fatalf("inputTarget 返回空进程名：D45 绑定键没有着落")
		}
		return
	}
	// Headless shape: hwnd is unknown to the test, so pidOf is stubbed to our
	// own pid - the window the predicate is asked about is still THIS process,
	// which is the one thing every runner can promise.
	b.windowPID = func(uintptr) uint32 { return uint32(windows.GetCurrentProcessId()) }
	got, name, err := b.inputTarget(map[string]any{"hwnd": 1234567890})
	if err != nil {
		t.Fatalf("自己的 pid 被判提权：%v", err)
	}
	if got != 1234567890 || strings.TrimSpace(name) == "" {
		t.Fatalf("放行腿读数异常：hwnd=%d name=%q", got, name)
	}
}

// Test23ProcessNameIsTheD45BindingKey pins the binding key the confirm card's
// 「本会话内允许向此应用」 option writes into approval_grant.pattern (ticket
// 49's lifecycle, fed here): realProcessName must be the image name, lower
// case, no directory - the canonical shape the grant matches against. The pid
// is the test process itself, taken from the same live window the leg above
// resolved (or from our own pid in the headless shape).
func Test23ProcessNameIsTheD45BindingKey(t *testing.T) {
	pid := uint32(windows.GetCurrentProcessId())
	if h := selfHwnd23(t); h != 0 {
		if viaWindow := realWindowPID(h); viaWindow != pid {
			t.Fatalf("前置失败：控制台窗口 %#x 的 pid=%d 不是本进程 %d", h, viaWindow, pid)
		}
	}
	got := realProcessName(pid)
	if got == "" {
		t.Fatalf("realProcessName(pid=%d) 回了空串：绑定键无从谈起", pid)
	}
	if got != strings.ToLower(got) {
		t.Errorf("绑定键 %q 不是小写（D45 pattern 是小写镜像名）", got)
	}
	if strings.ContainsAny(got, `\/`) {
		t.Errorf("绑定键 %q 带目录分隔符：pattern 只要镜像名", got)
	}
	if strings.HasSuffix(got, ".exe") == false && strings.Contains(got, ".") == false {
		t.Logf("绑定键 %q 没有扩展名（镜像名本身没有时也成立，只记录不判红）", got)
	}
}
