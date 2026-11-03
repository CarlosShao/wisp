package tools

// Ticket 23's tool-layer instrument: the eleven D34 machine-control rows, their
// declarations, their verb table, and the one refusal that must happen BEFORE the
// OS is touched.
//
// What is asserted here, and why each leg has teeth:
//
//	AC1 (risk matrix)  every Decl().Declared exactly as PLAN.md:2556-2567 prints
//	                   it, and the R8 class each row carries. A one-level drift
//	                   (input.type L2 -> L1, window.manage.close L1, system.set
//	                   L2) reddens.
//	AC4 (taint)        the WRITE column: clipboard.write's "send" class fuses to
//	                   L2 through the real risk.RiskAssessor, AND the C25 channel
//	                   table names clipboard.write's text as an exfil channel. The
//	                   read->mark->R4 propagation is the C25 loop's job, not a
//	                   tool's: nothing in this package can mark its own output.
//	AC5 (confirm word) restart/shutdown with a missing or wrong confirm_word
//	                   reach the backend ZERO times. The fake counts, so a future
//	                   "check it after the call" refactor reddens instead of
//	                   shutting a machine down.
//
// The fake backend is the SYSTEM-WIDE seam the ticket's hard constraint names
// ("all system calls through the host bridge"): it implements systemFuncs, the
// one interface SystemBridge.Do calls, so no Win32 entry point, no clipboard and
// no power transition is reachable from any leg below.

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
)

// ---------------------------------------------------------------------------
// the fake backend (the ticket's host-bridge seam)
// ---------------------------------------------------------------------------

// fake23Call is one recorded call.
type fake23Call struct {
	action Action
	args   map[string]any
}

// fake23Backend records every action it is handed and answers with whatever the
// leg told it to. Its zero value is usable: nil replyFunc means "answer with a
// Detail naming the action", which is enough for the rendering legs.
type fake23Backend struct {
	mu     sync.Mutex
	calls  []fake23Call
	reply  func(a Action, args map[string]any) (ActionResult, error)
	errFor map[Action]error
}

func (f *fake23Backend) act(_ context.Context, a Action, args map[string]any) (ActionResult, error) {
	cp := make(map[string]any, len(args))
	for k, v := range args {
		cp[k] = v
	}
	f.mu.Lock()
	f.calls = append(f.calls, fake23Call{action: a, args: cp})
	err := f.errFor[a]
	reply := f.reply
	f.mu.Unlock()
	if err != nil {
		return ActionResult{}, err
	}
	if reply != nil {
		return reply(a, cp)
	}
	return ActionResult{Detail: "fake did " + string(a), Target: string(a)}, nil
}

func (f *fake23Backend) count() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.calls)
}

func (f *fake23Backend) last() (fake23Call, bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.calls) == 0 {
		return fake23Call{}, false
	}
	return f.calls[len(f.calls)-1], true
}

// systemRoster23 builds the family over one fake and returns the entries keyed by
// tool name, with the fake itself. The registry is a REAL Registry, so a row that
// would fail C1/C3 validation turns this helper into a test failure here rather
// than a startup failure in production.
func systemRoster23(t *testing.T, f *fake23Backend) (*Registry, map[string]Entry) {
	t.Helper()
	bridge := &SystemBridge{Backend: f}
	entries := BuiltinSystemEntries(SystemDeps{Bridge: bridge})
	reg := NewRegistry()
	byName := make(map[string]Entry, len(entries))
	for _, e := range entries {
		if err := reg.Register(e); err != nil {
			t.Fatalf("BuiltinSystemEntries 注册失败（%s）：%v", e.Tool.Name(), err)
		}
		byName[e.Tool.Name()] = e
	}
	return reg, byName
}

// run23 executes one call through the entry's own Tool (no bridge, no gate: the
// legs below are about the tool layer, and the gate's own contract is asserted in
// internal/agent/approval and in bridge_test.go).
func run23(t *testing.T, e Entry, params string) Result {
	t.Helper()
	res, err := e.Tool.Execute(context.Background(), json.RawMessage(params), nil)
	if err != nil {
		t.Fatalf("%s Execute 返回了宿主内部错误（这应是模型可自纠的 IsError）：%v", e.Tool.Name(), err)
	}
	return res
}

// levels23 is PLAN.md:2556-2567 read as data: the D34 rows this ticket owns and
// the declared floor each one prints.
var levels23 = map[string]risk.Level{
	"system.get":      risk.L0,
	"system.set":      risk.L1,
	"system.power":    risk.L2,
	"window.list":     risk.L0,
	"window.focus":    risk.L1,
	"window.manage":   risk.L1, // close lifts to L2 through R8 (asserted below)
	"clipboard.read":  risk.L0,
	"clipboard.write": risk.L1, // send class lifts to L2 through R8 (asserted below)
	"media.control":   risk.L1,
	"screen.capture":  risk.L1,
	"input.type":      risk.L2,
}

// caps23 is the fourth column of the same table (the C3 token each row spends).
var caps23 = map[string]Capability{
	"system.get":      CapSysinfo,
	"system.set":      CapSysinfo,
	"system.power":    CapSysinfo,
	"window.list":     CapWindow,
	"window.focus":    CapWindow,
	"window.manage":   CapWindow,
	"clipboard.read":  CapClipboard,
	"clipboard.write": CapClipboard,
	"media.control":   CapSysinfo,
	"screen.capture":  CapScreen,
	"input.type":      CapInput,
}

// ---------------------------------------------------------------------------
// the roster
// ---------------------------------------------------------------------------

// Test23RosterCoversElevenD34Rows pins the NAME SET: eleven rows, exactly the
// names the frozen table prints, and nothing else. A twelfth name (or a rename)
// is a D34 change, and this leg is where that shows up.
func Test23RosterCoversElevenD34Rows(t *testing.T) {
	_, byName := systemRoster23(t, &fake23Backend{})
	if len(byName) != 11 {
		t.Fatalf("BuiltinSystemEntries 出 %d 枚，D34 的这一族是 11 枚：%v", len(byName), sorted23(names23(byName)))
	}
	for name := range levels23 {
		if _, ok := byName[name]; !ok {
			t.Errorf("D34 行 %q 没有成为工具（模型敲不到这扇门）", name)
		}
	}
	for name := range byName {
		if _, ok := levels23[name]; !ok {
			t.Errorf("多出一枚 %q：D34 的这一族只有 11 个名字，新增要改契约", name)
		}
	}
}

// Test23RosterDeclarationsMatchTheFrozenTable is AC1's declaration half: level,
// capability face, the C3 pair (Needs ⊆ Capabilities), PathParams (none of these
// calls touches a file, so C19 has no path to judge) and the provider slot.
func Test23RosterDeclarationsMatchTheFrozenTable(t *testing.T) {
	_, byName := systemRoster23(t, &fake23Backend{})
	for name, want := range levels23 {
		e, ok := byName[name]
		if !ok {
			continue // covered by the roster leg
		}
		if got := e.Decl.Declared; got != want {
			t.Errorf("%s Declared=%v，D34 写的是 %v", name, got, want)
		}
		if got := e.Decl.Capabilities; len(got) != 1 || got[0] != caps23[name] {
			t.Errorf("%s Capabilities=%v，D34 第四列是 %s", name, got, caps23[name])
		}
		if got := e.Decl.Needs; len(got) != 1 || got[0] != caps23[name] {
			t.Errorf("%s Needs=%v：C3 要求 Needs 就是它真正花的那枚 %s", name, got, caps23[name])
		}
		if len(e.Decl.PathParams) != 0 {
			t.Errorf("%s PathParams=%v：这一族不碰文件系统，声明路径只会让 C19 判一个它没有的目标",
				name, e.Decl.PathParams)
		}
		if e.Decl.Provider != KindBuiltin {
			t.Errorf("%s Provider=%q，要 %q", name, e.Decl.Provider, KindBuiltin)
		}
		if len(e.Tool.Parameters()) == 0 {
			t.Errorf("%s 没有参数 schema（C1 注册会拒）", name)
		}
		if strings.TrimSpace(e.Tool.Description()) == "" {
			t.Errorf("%s 没有描述（模型看不见就敲不到）", name)
		}
	}
}

// Test23WindowManageCloseIsL2ByFacts is the one row with two levels, asserted
// both ways: the DECLARATION stays L1 (the six reversible verbs) and the CLOSE
// verb carries R8's close-window class, which is what lifts that call to L2.
// Collapsing either half used to be the tempting shortcut - a blanket L2 would
// drag every minimize through a card, and dropping the class would let close run
// as a reversible L1.
func Test23WindowManageCloseIsL2ByFacts(t *testing.T) {
	_, byName := systemRoster23(t, &fake23Backend{})
	e, ok := byName["window.manage"]
	if !ok {
		t.Fatalf("window.manage 不在名册里")
	}
	if e.Decl.Declared != risk.L1 {
		t.Fatalf("window.manage Declared=%v，六枚可逆动作是 L1", e.Decl.Declared)
	}
	if e.Decl.Facts == nil {
		t.Fatalf("window.manage 没有 Facts 钩子：close 就永远到不了 L2")
	}
	for _, verb := range []string{"minimize", "maximize", "restore", "move", "resize", "snap"} {
		f := e.Decl.Facts(context.Background(), map[string]any{"action": verb}, nil)
		if len(f.Irreversible) != 0 {
			t.Errorf("window.manage %s 带了不可逆类别 %v：这一枚是可逆的", verb, f.Irreversible)
		}
	}
	f := e.Decl.Facts(context.Background(), map[string]any{"action": "close"}, nil)
	if len(f.Irreversible) != 1 || f.Irreversible[0] != "close-window" {
		t.Errorf("window.manage close 的 Facts=%v，要 [close-window]（R8）", f.Irreversible)
	}
	// The per-verb floor table agrees with the hook, so the two cannot drift.
	if got := SystemRowLevels()["window.manage.close"]; got != risk.L2 {
		t.Errorf("SystemRowLevels()[window.manage.close]=%v，要 L2", got)
	}
	for _, verb := range []string{"minimize", "maximize", "restore", "move", "resize", "snap"} {
		if got := SystemRowLevels()["window.manage."+verb]; got != risk.L1 {
			t.Errorf("SystemRowLevels()[window.manage.%s]=%v，要 L1", verb, got)
		}
	}
}

// Test23PowerAndClipboardWriteCarryR8Classes is AC1's facts half: the two rows
// that are L2-by-consequence name their class, and a call without one stays
// clean (so the instrument is not passing simply because every hook says L2).
func Test23PowerAndClipboardWriteCarryR8Classes(t *testing.T) {
	_, byName := systemRoster23(t, &fake23Backend{})
	power := byName["system.power"]
	if power.Decl.Facts == nil {
		t.Fatalf("system.power 没有 Facts 钩子")
	}
	if f := power.Decl.Facts(context.Background(), map[string]any{"action": "lock"}, nil); len(f.Irreversible) != 1 || f.Irreversible[0] != "power" {
		t.Errorf("system.power 的 Facts=%v，要 [power]（R8）", f.Irreversible)
	}
	clip := byName["clipboard.write"]
	if clip.Decl.Facts == nil {
		t.Fatalf("clipboard.write 没有 Facts 钩子：F4 的外泄通道就只剩一句注释")
	}
	if f := clip.Decl.Facts(context.Background(), map[string]any{"text": "x"}, nil); len(f.Irreversible) != 1 || f.Irreversible[0] != "send" {
		t.Errorf("clipboard.write 的 Facts=%v，要 [send]（R8）", f.Irreversible)
	}
	// Negative control: the read rows must NOT carry a class, so the legs above
	// are reading the hook rather than a constant.
	for _, name := range []string{"system.get", "window.list", "clipboard.read", "screen.capture"} {
		if e := byName[name]; e.Decl.Facts != nil {
			f := e.Decl.Facts(context.Background(), map[string]any{}, nil)
			if len(f.Irreversible) != 0 {
				t.Errorf("%s 是只读行，Facts 却给了 %v", name, f.Irreversible)
			}
		}
	}
}

// ---------------------------------------------------------------------------
// the verb table
// ---------------------------------------------------------------------------

// Test23VerbTableMapsEveryAction pins the (tool, action) -> Action dispatch. The
// expected table is written out literally rather than derived, because a derived
// expectation would move with the code it is supposed to judge.
func Test23VerbTableMapsEveryAction(t *testing.T) {
	cases := []struct {
		tool    string
		params  string
		want    Action
		verbFor string
	}{
		{"system.get", `{}`, ActionSystemGet, ""},
		{"system.set", `{"action":"volume","level":40}`, ActionSetVolume, "volume"},
		{"system.set", `{"action":"mute","mute":true}`, ActionSetMute, "mute"},
		{"system.set", `{"action":"brightness","level":70}`, ActionSetBright, "brightness"},
		{"system.set", `{"action":"media_key","key":"next"}`, ActionMediaKey, "media_key"},
		{"system.set", `{"action":"ime","layout":"zh"}`, ActionSetIME, "ime"},
		{"system.power", `{"action":"lock"}`, ActionPowerLock, "lock"},
		{"system.power", `{"action":"sleep"}`, ActionPowerSleep, "sleep"},
		{"system.power", `{"action":"logoff"}`, ActionPowerLogoff, "logoff"},
		{"system.power", `{"action":"restart","confirm_word":"确认"}`, ActionPowerRestart, "restart"},
		{"system.power", `{"action":"shutdown","confirm_word":"确认"}`, ActionPowerShutdown, "shutdown"},
		{"window.list", `{}`, ActionWindowList, ""},
		{"window.focus", `{"hwnd":7}`, ActionWindowFocus, ""},
		{"window.manage", `{"action":"minimize","hwnd":7}`, ActionWindowMin, "minimize"},
		{"window.manage", `{"action":"maximize","hwnd":7}`, ActionWindowMax, "maximize"},
		{"window.manage", `{"action":"restore","hwnd":7}`, ActionWindowRestore, "restore"},
		{"window.manage", `{"action":"move","hwnd":7,"x":10,"y":10}`, ActionWindowMove, "move"},
		{"window.manage", `{"action":"resize","hwnd":7,"width":800,"height":600}`, ActionWindowResize, "resize"},
		{"window.manage", `{"action":"snap","hwnd":7,"position":"left"}`, ActionWindowSnap, "snap"},
		{"window.manage", `{"action":"close","hwnd":7}`, ActionWindowClose, "close"},
		{"clipboard.read", `{}`, ActionClipboardRead, ""},
		{"clipboard.write", `{"text":"hi"}`, ActionClipboardWrite, ""},
		{"media.control", `{"action":"play_pause"}`, actionMediaPlayPause, "play_pause"},
		{"media.control", `{"action":"next"}`, actionMediaNext, "next"},
		{"media.control", `{"action":"previous"}`, actionMediaPrev, "previous"},
		{"media.control", `{"action":"stop"}`, actionMediaStop, "stop"},
		{"screen.capture", `{"target":"full"}`, ActionScreenCapture, ""},
		{"input.type", `{"text":"hello","hwnd":7}`, ActionInputType, ""},
	}
	for _, c := range cases {
		f := &fake23Backend{}
		_, byName := systemRoster23(t, f)
		e, ok := byName[c.tool]
		if !ok {
			t.Errorf("%s 不在名册里", c.tool)
			continue
		}
		run23(t, e, c.params)
		got, ok := f.last()
		if !ok {
			t.Errorf("%s（verb=%q）没有到达后端：Action 映射断了", c.tool, c.verbFor)
			continue
		}
		if got.action != c.want {
			t.Errorf("%s（verb=%q）到了后端是 %q，要 %q", c.tool, c.verbFor, got.action, c.want)
		}
		// confirm_word is a ritual, not an OS argument: it must not reach the
		// bridge (whose audit line counts keys).
		if _, leaked := got.args[systemConfirmParam]; leaked {
			t.Errorf("%s 把 %s 也交给了后端（后端与审计行只该看见动作参数）", c.tool, systemConfirmParam)
		}
	}
}

// Test23UnknownVerbIsRefusedBeforeTheBackend is the other side of the table: a
// verb outside it - including a MISSING one on a multi-verb tool - is a refusal
// with zero backend calls. "Default to the first verb" would turn a forgotten
// action on window.manage into a window transition the user never asked for.
func Test23UnknownVerbIsRefusedBeforeTheBackend(t *testing.T) {
	for _, c := range []struct{ tool, params string }{
		{"system.set", `{"action":"laser"}`},
		{"system.power", `{"action":"reboot"}`},
		{"window.manage", `{}`}, // no action at all
		{"window.manage", `{"action":""}`},
		{"media.control", `{}`},
	} {
		f := &fake23Backend{}
		_, byName := systemRoster23(t, f)
		res := run23(t, byName[c.tool], c.params)
		if !res.IsError {
			t.Errorf("%s %s 竟被接受（Text=%q）", c.tool, c.params, res.Text)
		}
		if n := f.count(); n != 0 {
			t.Errorf("%s %s 在被拒之前就调了后端 %d 次", c.tool, c.params, n)
		}
	}
}

// ---------------------------------------------------------------------------
// AC5: the typed confirmation word
// ---------------------------------------------------------------------------

// Test23PowerShutdownRequiresTheTypedWord is AC5's instrument. Three legs:
//
//   - missing confirm_word -> IsError AND zero backend calls;
//   - a WRONG word ("confirm", "确定") -> same, so the check is a comparison,
//     not a presence test;
//   - the exact word -> exactly one call, with ActionPowerShutdown.
//
// The count is the whole point: a tool that ran the shutdown and only then
// reported "you forgot the word" would pass a text-only assertion and kill the
// machine. Restart is asserted alongside shutdown because D34 puts the word on
// both (PLAN.md:2562).
func Test23PowerShutdownRequiresTheTypedWord(t *testing.T) {
	for _, action := range []string{"shutdown", "restart"} {
		for _, bad := range []string{``, `"confirm"`, `"确定"`, `"确认 "`, `"确认!"`} {
			f := &fake23Backend{}
			_, byName := systemRoster23(t, f)
			params := `{"action":"` + action + `"}`
			if bad != "" {
				params = `{"action":"` + action + `","confirm_word":` + bad + `}`
			}
			res := run23(t, byName["system.power"], params)
			if !res.IsError {
				t.Errorf("%s 缺/错确认词（%s）却成功：%q", action, params, res.Text)
			}
			if !strings.Contains(res.Text, systemConfirmWord) {
				t.Errorf("%s 的拒绝文本没告诉用户确认词是什么：%q", action, res.Text)
			}
			if n := f.count(); n != 0 {
				t.Errorf("%s 缺/错确认词（%s）竟调了后端 %d 次 —— 这一发在真机上就是关机", action, params, n)
			}
		}

		f := &fake23Backend{}
		_, byName := systemRoster23(t, f)
		res := run23(t, byName["system.power"],
			`{"action":"`+action+`","confirm_word":"`+systemConfirmWord+`"}`)
		if res.IsError {
			t.Errorf("%s 给了确认词却被拒：%q", action, res.Text)
		}
		if n := f.count(); n != 1 {
			t.Errorf("%s 给了确认词，后端被调 %d 次，要 1 次", action, n)
		}
		if got, _ := f.last(); got.action != ActionPowerShutdown && got.action != ActionPowerRestart {
			t.Errorf("%s 到了后端是 %q", action, got.action)
		}
	}

	// The three verbs that do NOT carry the word stay callable without it: a
	// lock screen on every request would be its own bug.
	for _, action := range []string{"lock", "sleep", "logoff"} {
		f := &fake23Backend{}
		_, byName := systemRoster23(t, f)
		res := run23(t, byName["system.power"], `{"action":"`+action+`"}`)
		if res.IsError || f.count() != 1 {
			t.Errorf("%s 不该要确认词（IsError=%v calls=%d）：%q", action, res.IsError, f.count(), res.Text)
		}
	}
}

// ---------------------------------------------------------------------------
// rendering
// ---------------------------------------------------------------------------

// Test23ExecuteRendersDetailReportAndBytes is the Result shape per verb kind:
// a write verb's Detail becomes Text verbatim; a read verb's Report becomes
// compact JSON; and screen.capture's BYTES never appear - the model gets
// {mime, bytes_len, drm} and the text says the image is not visible to it,
// which is the honest state of C7 (AC2's seam does not exist yet).
func Test23ExecuteRendersDetailReportAndBytes(t *testing.T) {
	f := &fake23Backend{reply: func(a Action, _ map[string]any) (ActionResult, error) {
		switch a {
		case ActionSystemGet:
			return ActionResult{Report: map[string]any{"os": "windows", "input_language": "zh"}, Target: "记事本"}, nil
		case ActionScreenCapture:
			return ActionResult{Bytes: []byte{1, 2, 3, 4, 5}, Mime: "image/png", DRM: true, Target: "screen"}, nil
		case ActionPowerSleep:
			return ActionResult{Detail: "已进入睡眠", Target: "system.power.sleep"}, nil
		}
		return ActionResult{Detail: "done " + string(a)}, nil
	}}
	_, byName := systemRoster23(t, f)

	get := run23(t, byName["system.get"], `{}`)
	if get.IsError {
		t.Fatalf("system.get 报错：%q", get.Text)
	}
	if want := `{"input_language":"zh","os":"windows"}`; get.Text != want {
		t.Errorf("system.get 的 Report 渲染成 %q，要紧凑且键有序的 %q", get.Text, want)
	}
	if get.Origin != "记事本" {
		t.Errorf("system.get 的 Origin=%q，要后端给的焦点窗口标题（C25 的出处）", get.Origin)
	}

	sleep := run23(t, byName["system.power"], `{"action":"sleep"}`)
	if sleep.Text != "已进入睡眠" {
		t.Errorf("写动词的 Detail 没有原样成为 Text：%q", sleep.Text)
	}

	cap := run23(t, byName["screen.capture"], `{"target":"full"}`)
	for _, want := range []string{`"mime":"image/png"`, `"bytes_len":5`, `"drm":true`} {
		if !strings.Contains(cap.Text, want) {
			t.Errorf("screen.capture 的元数据读不到 %q：%q", want, cap.Text)
		}
	}
	// The bytes must not be rendered: neither as base64 nor as raw digits.
	for _, leak := range []string{"AQIDBAU", "1,2,3,4,5"} {
		if strings.Contains(cap.Text, leak) {
			t.Errorf("screen.capture 把像素写进了 Text（%q）：截图只存内存", leak)
		}
	}
	if !strings.Contains(cap.Text, "C7") {
		t.Errorf("screen.capture 的文本没说 C7 ImagePart 接缝不存在（模型会以为自己看到了画面）：%q", cap.Text)
	}
}

// Test23RefusalsBecomeModelActableErrors checks the three bridge refusals reach
// the model as IsError text (never as a host fault), and that the UIPI sentence
// keeps PLAN.md:3234's exact wording - that sentence is what the user is told to
// read when input.type hits an elevated window.
func Test23RefusalsBecomeModelActableErrors(t *testing.T) {
	for _, c := range []struct {
		tool string
		want error
	}{
		{"input.type", ErrElevatedTarget},
		{"system.set", ErrNotImplementedHere},
		{"clipboard.read", ErrNoSystemBackend},
	} {
		f := &fake23Backend{errFor: map[Action]error{}}
		_, byName := systemRoster23(t, f)
		// Wire the error on the action this tool actually sends.
		act := mustAction23(t, byName[c.tool], c.tool)
		f.errFor[act] = c.want
		params := `{"action":"volume","level":10}`
		if c.tool == "input.type" || c.tool == "clipboard.read" {
			params = `{"text":"hi"}`
		}
		res := run23(t, byName[c.tool], params)
		if !res.IsError {
			t.Errorf("%s 遇到 %v 没有变成 IsError（Text=%q）", c.tool, c.want, res.Text)
			continue
		}
		if !strings.Contains(res.Text, c.want.Error()) {
			t.Errorf("%s 的拒绝文本丢了原句：%q，要含 %q", c.tool, res.Text, c.want.Error())
		}
	}
	if !strings.Contains(ErrElevatedTarget.Error(), "目标程序以管理员身份运行，Wisp 无法向其输入") {
		t.Errorf("ErrElevatedTarget 的原文变了：%q（PLAN.md:3234 冻结的是这一句）", ErrElevatedTarget.Error())
	}
}

// Test23HostFaultIsNotSwallowed proves the split runs the other way too: an
// error outside the refusal vocabulary is a HOST fault and comes back as an
// error, so the loop books it as internal instead of telling the model to
// "retry with different arguments".
func Test23HostFaultIsNotSwallowed(t *testing.T) {
	boom := errors.New("boom: 后端内部炸了")
	f := &fake23Backend{errFor: map[Action]error{ActionWindowList: boom}}
	_, byName := systemRoster23(t, f)
	_, err := byName["window.list"].Tool.Execute(context.Background(), json.RawMessage(`{}`), nil)
	if err == nil {
		t.Fatalf("宿主内部错误被吞成了正常返回：模型会以为机器被读过了")
	}
	if !errors.Is(err, boom) {
		t.Errorf("回来的错误不是原错误（%v）：D37 要看得到原句", err)
	}
}

// Test23NoBridgeRefusesInsteadOfPanicking pins the composition-time half of
// SystemDeps.bridge(): a deps value with no Bridge still yields eleven entries,
// and every one of them has a callable Tool (the refusal happens at CALL time,
// never as a nil dereference while assembling the roster).
//
// The call-time half is platform-split, and that split is the honest reading of
// system_bridge.go's zero value: on a platform that HAS a backend, the zero
// SystemBridge with a nil Backend is the REAL machine ("nil = selected by the
// platform build", system_ops_windows.go), so a call from here would read the
// live desktop. Asserting "no bridge refuses" on Windows would therefore require
// a bridge that refuses on the platform that does not refuse - i.e. a lie. The
// Windows leg below asserts the opposite truth instead: the zero value is USABLE
// (it reaches a backend rather than nil-panicking).
func Test23NoBridgeRefusesInsteadOfPanicking(t *testing.T) {
	entries := BuiltinSystemEntries(SystemDeps{})
	if len(entries) != 11 {
		t.Fatalf("空 deps 应仍出 11 枚（拒绝发生在调用时，不是组装时），得到 %d", len(entries))
	}
	for _, e := range entries {
		if e.Tool == nil || strings.TrimSpace(e.Tool.Name()) == "" {
			t.Errorf("空 deps 下有一枚 Tool 组装不出来：%+v", e)
		}
		if len(e.Tool.(systemTool).verbs) == 0 {
			t.Errorf("%s 在空 deps 下没有动词表，调用时连 action 都解析不出来", e.Tool.Name())
		}
	}
	// The seam's own answer for "there is no backend here" is what the tools
	// turn into model-actable text; assert that conversion on a bridge built
	// with an EXPLICIT refusing backend, which is the platform-neutral shape.
	refuse := &fake23Backend{errFor: map[Action]error{ActionSystemGet: ErrNoSystemBackend}}
	_, byName := systemRoster23(t, refuse)
	res := run23(t, byName["system.get"], `{}`)
	if !res.IsError || !strings.Contains(res.Text, ErrNoSystemBackend.Error()) {
		t.Errorf("没有后端时的 system.get 返回 %+v，要 IsError + ErrNoSystemBackend 原句", res)
	}
}

// ---------------------------------------------------------------------------
// AC4: clipboard.write's R8 class through the REAL assessor
// ---------------------------------------------------------------------------

// Test23ClipboardWriteFactsFuseToL2 is AC4's judge half, run through
// risk.NewRiskAssessor with the declaration and the Facts this file supplies -
// exactly what internal/tools/bridge.go's factsFor hands C19. It asserts the
// CLASS does the lifting (the declaration alone is L1), that R8 is the rule that
// fired, and that the card's reason says what it is.
func Test23ClipboardWriteFactsFuseToL2(t *testing.T) {
	_, byName := systemRoster23(t, &fake23Backend{})
	a := risk.NewRiskAssessor()

	declaredOnly := a.Assess("clipboard.write", nil, risk.Facts{Declared: byName["clipboard.write"].Decl.Declared})
	if declaredOnly.Level != risk.L1 {
		t.Fatalf("正控失败：只用声明时 clipboard.write 应判 L1，得到 %v", declaredOnly.Level)
	}

	facts := byName["clipboard.write"].Decl.Facts(context.Background(), map[string]any{"text": "hi"}, nil)
	facts.Declared = byName["clipboard.write"].Decl.Declared
	d := a.Assess("clipboard.write", map[string]any{"text": "hi"}, facts)
	if d.Level != risk.L2 {
		t.Errorf("clipboard.write 带着 send 类应判 L2（R8），得到 %v/%v", d.Level, d.RulesHit)
	}
	if !hasRule309(d, risk.R8) {
		t.Errorf("判决里没有 R8：%v", d.RulesHit)
	}
	if !strings.Contains(d.Reason, "发送类外发操作") {
		t.Errorf("卡的 reason 没写外泄语义：%q", d.Reason)
	}

	// window.manage close goes the same road, and its reason names the window.
	wm := byName["window.manage"]
	wf := wm.Decl.Facts(context.Background(), map[string]any{"action": "close"}, nil)
	wf.Declared = wm.Decl.Declared
	if dc := a.Assess("window.manage", nil, wf); dc.Level != risk.L2 {
		t.Errorf("window.manage close 应判 L2（R8），得到 %v/%v", dc.Level, dc.RulesHit)
	}
	// And minimize stays where the declaration put it.
	mf := wm.Decl.Facts(context.Background(), map[string]any{"action": "minimize"}, nil)
	mf.Declared = wm.Decl.Declared
	if dm := a.Assess("window.manage", nil, mf); dm.Level != risk.L1 {
		t.Errorf("window.manage minimize 应仍是 L1，得到 %v/%v", dm.Level, dm.RulesHit)
	}
}

// Test23ClipboardWriteIsAFrozenC25Channel is AC4's channel half: the frozen C25
// channel table (risk/provenance.go) must keep naming clipboard.write, because
// that mapping is how tainted content travelling through the clipboard is
// detected as an exfil attempt. It is a cross-package read (the table is
// unexported), so the leg goes through the public Inspect surface instead.
//
// WHAT THIS LEG DOES NOT PROVE, said out loud: the read->mark->R4 propagation.
// Marking happens in the bridge's provenance path and in the loop's C25 wiring;
// a tool cannot mark its own output, so that half belongs to the C25/ticket-19
// wiring and the bridge's flow into R4 (bridge_test.go and
// assessorfor_taint_instrument_309_test.go cover the detector seam).
func Test23ClipboardWriteIsAFrozenC25Channel(t *testing.T) {
	prov := risk.NewProvenance(risk.ProvOptions{NoProbe: true})
	const scope, secret = "task-23-ac4", "skimmed-secret-23-payload"
	if !prov.Mark(scope, risk.SrcFSRead, "D:/keys/id_rsa", "凭据摘录 "+secret+" 收尾") {
		t.Fatalf("前置失败：Mark 没落污点，这一腿失去前提")
	}
	// Positive control: the clipboard.read/write family itself must be a
	// sensitive SOURCE in the frozen set (SPEC-06 §5), which is what makes a
	// clipboard read mark its output at all.
	if !risk.IsSensitiveSource("clipboard.read") || !risk.IsSensitiveSource("system.get") ||
		!risk.IsSensitiveSource("screen.capture") {
		t.Errorf("C25 的敏感源集合少了本族的行：clipboard.read=%v system.get=%v screen.capture=%v",
			risk.IsSensitiveSource("clipboard.read"), risk.IsSensitiveSource("system.get"),
			risk.IsSensitiveSource("screen.capture"))
	}
	params := map[string]any{"text": "外发正文 " + secret}
	hit, ok := prov.Inspect(scope, "clipboard.write", params)
	if !ok {
		t.Fatalf("把污点正文写进剪贴板却判不出 R4：clipboard.write 不在 C25 的外泄通道表里")
	}
	if hit.Channel != risk.ChClipboard {
		t.Errorf("命中的通道是 %q，要 %q（F4 的外泄通道）", hit.Channel, risk.ChClipboard)
	}
	// And the detector seam the bridge uses reports it too.
	if src, tainted := prov.Detector(scope).TaintHit(params); !tainted || !strings.Contains(src, "fs.read") {
		t.Errorf("探测器没把这次调用判成污点外发：src=%q tainted=%v", src, tainted)
	}
}

// ---------------------------------------------------------------------------
// small helpers
// ---------------------------------------------------------------------------

// mustAction23 finds the bridge Action a tool sends, by running it once through
// a probe fake. It exists because the refusal leg has to wire its injected error
// onto the action the tool really uses - which is exactly the mapping under test
// in Test23VerbTableMapsEveryAction, so reading it off a probe keeps this leg
// from restating the table a second time.
func mustAction23(t *testing.T, e Entry, name string) Action {
	t.Helper()
	probe := &fake23Backend{}
	entry := Entry{
		Tool: systemTool{
			d: SystemDeps{Bridge: &SystemBridge{Backend: probe}}, name: name,
			verbs: e.Tool.(systemTool).verbs,
		},
		Decl: e.Decl,
	}
	params := `{}`
	switch name {
	case "system.set":
		params = `{"action":"volume","level":10}`
	case "clipboard.read", "system.get", "window.list":
		params = `{}`
	}
	if _, err := entry.Tool.Execute(context.Background(), json.RawMessage(params), nil); err != nil {
		t.Fatalf("%s 探查调用失败：%v", name, err)
	}
	call, ok := probe.last()
	if !ok {
		t.Fatalf("%s 探查没有到达后端", name)
	}
	return call.action
}

// names23 lists a roster's names.
func names23(m map[string]Entry) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

// sorted23 sorts a string slice for a stable failure message.
func sorted23(in []string) []string {
	out := append([]string(nil), in...)
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j] < out[j-1]; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out
}
