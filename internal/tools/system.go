package tools

// ---------------------------------------------------------------------------
// The machine-control tool family (ticket 23): the C4 face over SystemBridge.
//
// system_bridge.go froze the vocabulary (Action), the one entry point
// (SystemBridge.Do) and the refusals (ErrNoSystemBackend / ErrElevatedTarget /
// ErrNotImplementedHere); system_backend_windows.go is the OS half. This file is
// the third layer the D34 table actually names: eleven builtin tools whose
// declared risk floor, capability face and R8 facts are what C19 judges BEFORE
// Execute runs.
//
// Why one generic systemTool rather than eleven types: every verb in the family
// goes through the SAME seam with the same argument map, so a per-verb type
// would be eleven copies of the same Execute with a different Action literal -
// and the thing worth reviewing is the (tool, action) -> Action table plus the
// declaration each row carries, both of which are here in one place.
//
// The declaration is where this layer's security lives, so the split follows
// D34's own numbers rather than one number per tool:
//
//	system.get      L0   read-only; the focused-window title is a C25 source
//	system.set      L1   reversible mixer/brightness/IME changes
//	system.power    L2   R8 power class; restart/shutdown need the typed word
//	window.list     L0   read-only
//	window.focus    L1   reversible (the user can focus back)
//	window.manage   L1   ... EXCEPT close, which is L2 + R8 close-window
//	clipboard.read  L0   C25 source (the read is marked, not blocked)
//	clipboard.write L1   R8 send class: F4's exfil channel
//	media.control   L1   reversible transport keys
//	screen.capture  L1   pixels stay in memory; see the C7 note below
//	input.type      L2   injection into an arbitrary window (PLAN.md:2557)
//
// R1 makes Declared a LOWER BOUND, not a verdict: C19 fuses it with the rules,
// which is why the L1 rows still carry their R8 class through Decl.Facts. That
// hook is what puts 「关闭窗口」/「发送类外发操作」 on the card instead of letting
// a declared L1 read as reversible.
//
// AC COVERAGE IN THIS LAYER (read before claiming this ticket is closed):
//
//	AC1 risk matrix      covered here and by system_backend's own make: every
//	                     Decl().Declared and every Facts class is asserted.
//	AC4 clipboard taint  the WRITE column (Facts "send" -> C19 L2, and the
//	                     provenance channel naming clipboard.write) is asserted
//	                     here; the READ->mark->R4 propagation is the C25 loop's
//	                     wiring, which no tool can assert about itself.
//	AC5 power word       covered here: restart/shutdown with a missing or
//	                     different confirm_word never reaches the backend.
//
//	AC2 C7 ImagePart     NOT COVERED, and deliberately not faked. C7 has no
//	                     image part in this tree's llm.Content, and Result has no
//	                     field that carries bytes, so screen.capture renders
//	                     {mime, bytes_len, drm} as TEXT and the PNG never leaves
//	                     process memory. The model therefore cannot SEE the
//	                     screenshot until the C7 ImagePart -> LLM seam exists
//	                     (owner-approved contract work, PLAN.md:2962). Until it
//	                     does, screen.capture is a tool that reports what it
//	                     captured, not one that shows it.
//	AC3 confirm-card UI  NOT COVERED here: the three D45 options for input.type
//	                     (仅本次 / 本会话内允许向此应用 / 所有应用 with the
//	                     danger styling) are rendered by the native card, not by
//	                     a tool. This layer supplies what that card reads - L2 on
//	                     input.type, Facts on the rest, ErrElevatedTarget's exact
//	                     wording - and the rendering is the front-end's.
//
// What is NOT wired into this layer either: the C25 mark for system.get /
// clipboard.read happens in the bridge's provenance path (risk.SrcSystemGet,
// risk.SrcClipboardRead are already in the frozen source set), and the D45
// session-grant downgrade for input.type (L2 -> L1 bound to the target process
// name, grant lifecycle ticket 49) needs a grant, not a declaration: no row here
// mints or trusts one.
// ---------------------------------------------------------------------------

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/CarlosShao/wisp/internal/risk"
)

// SystemDeps is what the machine-control family needs from the host. One field,
// because the seam IS the host face: there is no roster, no canonicalizer and no
// file handle in this family (nothing here touches a path - C19 has no path of
// any of these calls to judge, which is why every Decl.PathParams is nil).
type SystemDeps struct {
	// Bridge is the host bridge. A nil Bridge is SAFE rather than fatal: Do
	// answers ErrNoSystemBackend, the tools turn that into IsError text, and
	// the task continues (D37). The zero value of SystemBridge is also usable
	// (Backend nil = the platform default backend), so a composition root that
	// only wants the real machine writes &SystemBridge{}.
	Bridge *SystemBridge
}

// systemConfirmWord is the typed confirmation word PLAN.md:2562 requires for
// restart and shutdown: 「关机/重启必须二次输入确认词」. It is compared verbatim,
// so a model that writes "confirm" or "确定" gets the refusal text back rather
// than a power transition - the check exists to be un-guessable by accident, and
// an over-eager assistant is exactly what it is guarding against.
const systemConfirmWord = "确认"

// systemConfirmParam is the parameter name the word arrives in.
const systemConfirmParam = "confirm_word"

// ---------------------------------------------------------------------------
// the tool
// ---------------------------------------------------------------------------

// systemTool is one D34 row. verbs maps the parameter action ("" when the tool
// has no action parameter) onto the bridge Action, so the table in
// BuiltinSystemEntries reads as the dispatch it is.
type systemTool struct {
	d     SystemDeps
	name  string
	verbs map[string]Action
}

// systemActionParam is the one parameter that selects a verb inside the tools
// that have several (system.set, system.power, window.manage, media.control).
const systemActionParam = "action"

// Execute implements C1 Tool.
//
// Every failure the MODEL can act on comes back as Result.IsError (D37: a tool
// refusal is booked class "tool", the model self-corrects and the task lives on);
// only a genuine host fault is returned as an error. The three refusals the
// bridge vocabulary defines are all model-actable - "no backend here", "not
// implemented here" and the UIPI answer - so they are carried in the text, with
// ErrElevatedTarget's PLAN.md:3234 wording intact because that sentence is what
// the user is supposed to read.
func (t systemTool) Execute(ctx context.Context, params json.RawMessage, onUpdate func(string)) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	var args map[string]any
	if err := json.Unmarshal(params, &args); err != nil {
		return Result{Text: "参数解析失败：" + err.Error(), IsError: true}, nil
	}
	action, err := t.resolve(args)
	if err != nil {
		return Result{Text: err.Error(), IsError: true}, nil
	}
	// AC5: the typed word is checked BEFORE the bridge is called, so a missing
	// or wrong word costs zero backend calls. The check lives here rather than
	// in the backend because it is a property of the D34 row (only restart and
	// shutdown carry it), not of the OS call underneath.
	if needWord := action == ActionPowerRestart || action == ActionPowerShutdown; needWord {
		if got := argString(args, systemConfirmParam); got != systemConfirmWord {
			return Result{
				Text:    "关机/重启需在参数 confirm_word 键入确认词「" + systemConfirmWord + "」，本次未执行",
				IsError: true,
			}, nil
		}
	}
	// confirm_word is not an argument of any OS call: it is dropped before the
	// bridge sees the map, so the audit line and the backend both describe the
	// action rather than the ritual around it.
	call := make(map[string]any, len(args))
	for k, v := range args {
		if k == systemConfirmParam {
			continue
		}
		call[k] = v
	}
	res, err := t.d.bridge().Do(ctx, action, call)
	if err != nil {
		if isSystemRefusal(err) {
			return Result{IsError: true, Text: err.Error()}, nil
		}
		return Result{}, err
	}
	if onUpdate != nil {
		onUpdate(t.name + " " + string(action))
	}
	return renderSystemResult(t.name, res), nil
}

// resolve maps (tool, params[action]) onto one bridge Action. An empty action is
// a single-verb tool (system.get, window.list, window.focus, clipboard.*,
// screen.capture, input.type); anything else must name a verb from the table -
// including the empty string, which is refused rather than defaulted. Defaulting
// would turn "I forgot the action" into "close the window" for window.manage.
func (t systemTool) resolve(args map[string]any) (Action, error) {
	if len(t.verbs) == 1 {
		for k, a := range t.verbs {
			if k == "" {
				return a, nil
			}
		}
	}
	verb := strings.ToLower(strings.TrimSpace(argString(args, systemActionParam)))
	a, ok := t.verbs[verb]
	if !ok {
		return "", fmt.Errorf("%s 的 action 必须是以下之一：%s（收到 %q，未执行）",
			t.name, systemVerbNames(t.verbs), verb)
	}
	return a, nil
}

// systemVerbNames lists a tool's verbs for the refusal text, sorted so the text
// a user sees does not reshuffle between builds.
func systemVerbNames(verbs map[string]Action) string {
	names := make([]string, 0, len(verbs))
	for k := range verbs {
		if k == "" {
			continue
		}
		names = append(names, k)
	}
	sort.Strings(names)
	return strings.Join(names, " / ")
}

// bridge returns the host bridge to call. A nil deps.Bridge becomes the zero
// SystemBridge, which answers ErrNoSystemBackend for every action: a tool
// assembled without a bridge refuses, it does not panic.
func (d SystemDeps) bridge() *SystemBridge {
	if d.Bridge != nil {
		return d.Bridge
	}
	return &SystemBridge{}
}

// isSystemRefusal reports whether err is one of the three refusals the family
// defines in its vocabulary. They are the model-actable half: an unreachable
// backend, an unimplemented verb on this platform, and the UIPI answer.
func isSystemRefusal(err error) bool {
	return errors.Is(err, ErrNoSystemBackend) ||
		errors.Is(err, ErrNotImplementedHere) ||
		errors.Is(err, ErrElevatedTarget)
}

// renderSystemResult turns an ActionResult into the text the model reads. Three
// shapes, one per kind of verb:
//
//   - Detail: the write verbs' one-liner ("音量已设为 40%").
//   - Report: the read verbs' structured payload, as compact JSON so a test
//     asserts keys rather than prose.
//   - Bytes: screen.capture. There is NO image channel in Result and none in
//     C7 in this tree, so what goes back is the metadata ({mime, bytes_len,
//     drm}) and not one byte of the PNG. The wording says so out loud: a text
//     render that stayed silent would let the model claim it "saw" the screen.
func renderSystemResult(name string, res ActionResult) Result {
	out := Result{Origin: res.Target}
	switch {
	case len(res.Bytes) > 0:
		rep := map[string]any{
			"mime":      orDefaultString(res.Mime, "application/octet-stream"),
			"bytes_len": len(res.Bytes),
			"drm":       res.DRM,
		}
		if res.Target != "" {
			rep["target"] = res.Target
		}
		if res.Detail != "" {
			rep["detail"] = res.Detail
		}
		out.Text = renderSystemReport(rep) + "\n" +
			"（截图像素只在内存里，没有落盘；C7 ImagePart 到 LLM 的接缝本仓当前不存在，" +
			"所以模型看不到画面，只能拿到上面这份元数据）"
		return out
	case len(res.Report) > 0:
		out.Text = renderSystemReport(res.Report)
		if res.Detail != "" {
			out.Text = res.Detail + "\n" + out.Text
		}
		return out
	case res.Detail != "":
		out.Text = res.Detail
		return out
	}
	// Nothing at all came back. Saying "done" would be a claim about the
	// machine that no verb here made; the honest sentence is that it ran.
	out.Text = name + " 已完成，但后端没有返回任何内容"
	return out
}

// renderSystemReport marshals a report as compact, key-sorted JSON. json.Marshal
// already sorts map keys, so the text is stable across runs - which is what lets
// a test compare it byte for byte.
func renderSystemReport(rep map[string]any) string {
	b, err := json.Marshal(rep)
	if err != nil {
		return "（后端返回值无法序列化：" + err.Error() + "）"
	}
	return string(b)
}

func orDefaultString(s, def string) string {
	if strings.TrimSpace(s) == "" {
		return def
	}
	return s
}

// ---------------------------------------------------------------------------
// facts: the R8 input C19 judges the call with
// ---------------------------------------------------------------------------

// powerFacts is system.power's R8 class. It is declared for EVERY verb in the
// row, not only restart/shutdown: lock/sleep/logoff also end the session (and
// with it every unsaved editor), and the frozen class set has one word for the
// whole row (SPEC-06 §3 R8: "power"). Narrowing it here would be this file
// inventing a rule.
func powerFacts(context.Context, map[string]any, []string) risk.Facts {
	return risk.Facts{Irreversible: []string{"power"}}
}

// windowManageFacts is window.manage's R8 class: "close-window" for close,
// nothing for the six reversible verbs. It is a hook rather than two
// declarations because the level and the class both depend on the SAME action
// parameter - one switch is the only way the two can never disagree.
//
// An unknown/missing action returns NO class here: the refusal in resolve() is
// what answers it, and inflating it to close-window would put a card in front of
// the user for a call that never runs.
func windowManageFacts(_ context.Context, params map[string]any, _ []string) risk.Facts {
	verb := strings.ToLower(strings.TrimSpace(argString(params, systemActionParam)))
	if verb == "close" {
		return risk.Facts{Irreversible: []string{"close-window"}}
	}
	return risk.Facts{}
}

// clipboardWriteFacts is clipboard.write's R8 class. D34's note owns the reason
// ("算外泄通道", F4): the clipboard is where content leaves the assistant and
// becomes available to any program the user later pastes into, and the frozen
// class set spells that "send". It is L2 by consequence, above the row's own L1
// floor - which is exactly what R1's lower-bound rule is for.
func clipboardWriteFacts(context.Context, map[string]any, []string) risk.Facts {
	return risk.Facts{Irreversible: []string{"send"}}
}

// ---------------------------------------------------------------------------
// declarations
// ---------------------------------------------------------------------------

// systemDecl builds one row's declaration. Capabilities and Needs are the same
// token: this family reads/acts on the machine and nothing else, and a Needs
// outside Capabilities is a hard C3 reject, so the pair is written once.
func systemDecl(level risk.Level, cap Capability, facts func(context.Context, map[string]any, []string) risk.Facts) Decl {
	return Decl{
		Capabilities: []Capability{cap},
		Needs:        []Capability{cap},
		Declared:     level,
		PathParams:   nil,
		Facts:        facts,
		Resident:     true,
		Provider:     KindBuiltin,
	}
}

// Decls for the family. Each is a function rather than a package variable for
// the same reason PlanPresentDecl is: the roster a test reads must be the roster
// the registry froze, never a second copy that can drift.
func SystemGetDecl() Decl      { return systemDecl(risk.L0, CapSysinfo, nil) }
func SystemSetDecl() Decl      { return systemDecl(risk.L1, CapSysinfo, nil) }
func SystemPowerDecl() Decl    { return systemDecl(risk.L2, CapSysinfo, powerFacts) }
func WindowListDecl() Decl     { return systemDecl(risk.L0, CapWindow, nil) }
func WindowFocusDecl() Decl    { return systemDecl(risk.L1, CapWindow, nil) }
func WindowManageDecl() Decl   { return windowManageDecl() }
func ClipboardReadDecl() Decl  { return systemDecl(risk.L0, CapClipboard, nil) }
func ClipboardWriteDecl() Decl { return systemDecl(risk.L1, CapClipboard, clipboardWriteFacts) }
func MediaControlDecl() Decl   { return systemDecl(risk.L1, CapSysinfo, nil) }
func ScreenCaptureDecl() Decl  { return systemDecl(risk.L1, CapScreen, nil) }
func InputTypeDecl() Decl      { return systemDecl(risk.L2, CapInput, nil) }

// windowManageDecl is WindowManageDecl's body. The row's floor is L1 (the six
// reversible verbs); close reaches L2 through R8 in the Facts hook, so the
// declared number stays honest for the common case instead of dragging every
// minimize through a confirmation card.
func windowManageDecl() Decl {
	return systemDecl(risk.L1, CapWindow, windowManageFacts)
}

// SystemRowLevels is the per-verb effective floor of the one row whose level is
// not a single number: window.manage. It exists so a caller (a test, a future
// roster view) can ask "what does this verb declare" without re-deriving it from
// prose, and it is asserted against the Facts hook by the ticket-23 test.
func SystemRowLevels() map[string]risk.Level {
	return map[string]risk.Level{
		"window.manage.minimize": risk.L1,
		"window.manage.maximize": risk.L1,
		"window.manage.restore":  risk.L1,
		"window.manage.move":     risk.L1,
		"window.manage.resize":   risk.L1,
		"window.manage.snap":     risk.L1,
		"window.manage.close":    risk.L2,
	}
}

// ---------------------------------------------------------------------------
// schemas
// ---------------------------------------------------------------------------

// The parameter documents. Shape rules for the whole family: every object is
// closed (additionalProperties:false), every enumerated verb names its
// alternatives in the description as well as the enum, and no schema offers a
// field a caller could use to mark something approved - confirm_word is a
// RITUAL the user's own words supply, not an approval channel (the L2 card is).
var (
	// systemEmptySchema is the closed empty object. It is never reached by a
	// name in BuiltinSystemEntries (the switch covers all eleven), but a
	// C1 Parameters must never be nil: a tool with no document is a
	// registration failure, and returning an empty one is the honest shape
	// for a row nobody wrote a schema for yet.
	systemEmptySchema = JSONSchema(`{"type":"object","properties":{},"additionalProperties":false}`)

	systemGetSchema = JSONSchema(`{"type":"object","properties":{` +
		`"detail":{"type":"boolean","description":"保留位：true 时后端应返回更详细的读数（当前后端总是返回全量）"}` +
		`},"additionalProperties":false}`)

	systemSetSchema = JSONSchema(`{"type":"object","properties":{` +
		`"action":{"type":"string","enum":["volume","mute","brightness","media_key","ime"],` +
		`"description":"要改的设置：volume 音量 / mute 静音 / brightness 亮度 / media_key 媒体键 / ime 输入法"},` +
		`"level":{"type":"integer","minimum":0,"maximum":100,"description":"action=volume 或 brightness 时的目标百分比 0..100"},` +
		`"mute":{"type":"boolean","description":"action=mute 时 true 静音、false 取消静音"},` +
		`"key":{"type":"string","enum":["play_pause","next","previous","stop"],"description":"action=media_key 时按哪个媒体键"},` +
		`"layout":{"type":"string","enum":["en","zh"],"description":"action=ime 时切到哪种输入法"}` +
		`},"required":["action"],"additionalProperties":false}`)

	systemPowerSchema = JSONSchema(`{"type":"object","properties":{` +
		`"action":{"type":"string","enum":["lock","sleep","restart","shutdown","logoff"],` +
		`"description":"lock 锁屏 / sleep 睡眠 / restart 重启 / shutdown 关机 / logoff 注销"},` +
		`"confirm_word":{"type":"string",` +
		`"description":"restart 与 shutdown 必须把确认词原文「确认」交给用户核对后填在这里；缺失或不符则不执行"}` +
		`},"required":["action"],"additionalProperties":false}`)

	windowListSchema = JSONSchema(`{"type":"object","properties":{` +
		`"max_results":{"type":"integer","minimum":1,"description":"最多列多少个窗口（默认 50）"}` +
		`},"additionalProperties":false}`)

	windowTargetProperties = `"hwnd":{"type":"integer","minimum":1,` +
		`"description":"window.list 给出的窗口句柄数字（首选）"},` +
		`"title":{"type":"string","description":"窗口标题子串（大小写不敏感；匹配到多个会被拒绝，请改用 hwnd）"}`

	windowFocusSchema = JSONSchema(`{"type":"object","properties":{` + windowTargetProperties +
		`},"additionalProperties":false}`)

	windowManageSchema = JSONSchema(`{"type":"object","properties":{` +
		`"action":{"type":"string","enum":["minimize","maximize","restore","move","resize","snap","close"],` +
		`"description":"minimize 最小化 / maximize 最大化 / restore 还原 / move 移动 / resize 缩放 /` +
		` snap 贴靠 / close 关闭（可能丢失未保存内容，走 L2 确认）"},` +
		windowTargetProperties + `,` +
		`"x":{"type":"integer","description":"action=move 时的左上角 x（屏幕像素）"},` +
		`"y":{"type":"integer","description":"action=move 时的左上角 y（屏幕像素）"},` +
		`"width":{"type":"integer","minimum":1,"description":"action=resize 时的宽（像素）"},` +
		`"height":{"type":"integer","minimum":1,"description":"action=resize 时的高（像素）"},` +
		`"position":{"type":"string","enum":["left","right","top","bottom","top_left","top_right","bottom_left","bottom_right"],` +
		`"description":"action=snap 时贴到屏幕哪一半"}` +
		`},"required":["action"],"additionalProperties":false}`)

	clipboardReadSchema = JSONSchema(`{"type":"object","properties":{` +
		`"max_chars":{"type":"integer","minimum":1,"description":"最多读回多少字符（后端按上限截断）"}` +
		`},"additionalProperties":false}`)

	clipboardWriteSchema = JSONSchema(`{"type":"object","properties":{` +
		`"text":{"type":"string","minLength":1,"description":"要写进剪贴板的正文"}` +
		`},"required":["text"],"additionalProperties":false}`)

	mediaControlSchema = JSONSchema(`{"type":"object","properties":{` +
		`"action":{"type":"string","enum":["play_pause","next","previous","stop"],` +
		`"description":"play_pause 播放/暂停 / next 下一首 / previous 上一首 / stop 停止"}` +
		`},"required":["action"],"additionalProperties":false}`)

	screenCaptureSchema = JSONSchema(`{"type":"object","properties":{` +
		`"target":{"type":"string","enum":["full","window","region"],"description":"截什么：full 全屏 / window 某个窗口 / region 一块区域"},` +
		windowTargetProperties + `,` +
		`"x":{"type":"integer","description":"target=region 时区域的左上角 x"},` +
		`"y":{"type":"integer","description":"target=region 时区域的左上角 y"},` +
		`"width":{"type":"integer","minimum":1,"description":"target=region 时的宽"},` +
		`"height":{"type":"integer","minimum":1,"description":"target=region 时的高"}` +
		`},"additionalProperties":false}`)

	inputTypeSchema = JSONSchema(`{"type":"object","properties":{` +
		`"text":{"type":"string","description":"要注入的文本（与 keys 至少给一个）"},` +
		`"keys":{"type":"string","description":"要按下的键序列，例如 ^s 或 {ENTER}（与 text 至少给一个）"},` +
		windowTargetProperties +
		`},"additionalProperties":false}`)
)

// ---------------------------------------------------------------------------
// descriptions
// ---------------------------------------------------------------------------

// Each Description states what the tool does AND what it cannot or will not do,
// because that second half is what stops the model from promising the user
// something the row refuses: the backend is not wired on every platform, screens
// cannot be seen yet, and an elevated window is a hard Windows limit rather than
// a permission Wisp can ask for.
func (t systemTool) Description() string {
	switch t.name {
	case "system.get":
		return "读本机状态：当前焦点窗口标题与进程、时间、输入法，以及平台能读到的电量/音量/亮度（读不到的值不编）"
	case "system.set":
		return "改本机设置：音量、静音、亮度、媒体键、输入法（可逆；音量绝对电平与亮度依赖后端，未接入时会如实报未实现）"
	case "system.power":
		return "锁屏 / 睡眠 / 重启 / 关机 / 注销（L2 必确认；重启与关机还要用户把确认词「确认」交回 confirm_word，缺失即不执行）"
	case "window.list":
		return "列出可见且有标题的顶层窗口（hwnd、标题、进程、是否最小化/最大化），供 window.* 与 input.type 指定目标"
	case "window.focus":
		return "把指定窗口切到前台（L1；系统前台锁可能拒绝，拒绝时会说明）"
	case "window.manage":
		return "管理指定窗口：最小化/最大化/还原/移动/缩放/贴靠是 L1；关闭是 L2（可能丢失未保存内容）"
	case "clipboard.read":
		return "读剪贴板文本（L0；内容会被标记为外部来源，之后写往外部的调用可能因此升级确认）"
	case "clipboard.write":
		return "把文本写进剪贴板（L1；剪贴板是外泄通道，内容按 R8 send 类判 L2 确认）"
	case "media.control":
		return "控制正在播放的媒体：播放/暂停、下一首、上一首、停止（发媒体键，可逆）"
	case "screen.capture":
		return "截屏（全屏/窗口/区域）。当前只能回报 {mime, bytes_len, drm} 元数据，像素留在内存里、模型看不到画面"
	case "input.type":
		return "向目标窗口注入文本或键序列（L2；目标程序以管理员身份运行时为 Windows UIPI 硬限制，会明确拒绝且不提权绕过）"
	}
	return ""
}

// Parameters implements C1 Tool.
func (t systemTool) Parameters() JSONSchema {
	switch t.name {
	case "system.get":
		return systemGetSchema
	case "system.set":
		return systemSetSchema
	case "system.power":
		return systemPowerSchema
	case "window.list":
		return windowListSchema
	case "window.focus":
		return windowFocusSchema
	case "window.manage":
		return windowManageSchema
	case "clipboard.read":
		return clipboardReadSchema
	case "clipboard.write":
		return clipboardWriteSchema
	case "media.control":
		return mediaControlSchema
	case "screen.capture":
		return screenCaptureSchema
	case "input.type":
		return inputTypeSchema
	}
	return systemEmptySchema
}

// Name implements C1 Tool.
func (t systemTool) Name() string { return t.name }

// ---------------------------------------------------------------------------
// the roster
// ---------------------------------------------------------------------------

// BuiltinSystemEntries returns the eleven D34 rows of this family, in the order
// the frozen table prints them. The caller (cmd/wisp/run.go) registers them into
// the same Registry as every other builtin, so the bridge's C3 check, C19
// verdict, C18 card and C22 deadline all apply without a second code path.
//
// Verb tables are part of the declaration review: a verb missing here is
// unreachable (the model gets the list back and can retry), and a verb mapped to
// the wrong Action is what the ticket-23 test's table nails down.
func BuiltinSystemEntries(d SystemDeps) []Entry {
	return []Entry{
		{
			Tool: systemTool{d: d, name: "system.get", verbs: map[string]Action{"": ActionSystemGet}},
			Decl: SystemGetDecl(),
		},
		{
			Tool: systemTool{d: d, name: "system.set", verbs: map[string]Action{
				"volume":     ActionSetVolume,
				"mute":       ActionSetMute,
				"brightness": ActionSetBright,
				"media_key":  ActionMediaKey,
				"ime":        ActionSetIME,
			}},
			Decl: SystemSetDecl(),
		},
		{
			Tool: systemTool{d: d, name: "system.power", verbs: map[string]Action{
				"lock":     ActionPowerLock,
				"sleep":    ActionPowerSleep,
				"restart":  ActionPowerRestart,
				"shutdown": ActionPowerShutdown,
				"logoff":   ActionPowerLogoff,
			}},
			Decl: SystemPowerDecl(),
		},
		{
			Tool: systemTool{d: d, name: "window.list", verbs: map[string]Action{"": ActionWindowList}},
			Decl: WindowListDecl(),
		},
		{
			Tool: systemTool{d: d, name: "window.focus", verbs: map[string]Action{"": ActionWindowFocus}},
			Decl: WindowFocusDecl(),
		},
		{
			Tool: systemTool{d: d, name: "window.manage", verbs: map[string]Action{
				"minimize": ActionWindowMin,
				"maximize": ActionWindowMax,
				"restore":  ActionWindowRestore,
				"move":     ActionWindowMove,
				"resize":   ActionWindowResize,
				"snap":     ActionWindowSnap,
				"close":    ActionWindowClose,
			}},
			Decl: WindowManageDecl(),
		},
		{
			Tool: systemTool{d: d, name: "clipboard.read", verbs: map[string]Action{"": ActionClipboardRead}},
			Decl: ClipboardReadDecl(),
		},
		{
			Tool: systemTool{d: d, name: "clipboard.write", verbs: map[string]Action{"": ActionClipboardWrite}},
			Decl: ClipboardWriteDecl(),
		},
		{
			Tool: systemTool{d: d, name: "media.control", verbs: map[string]Action{
				"play_pause": actionMediaPlayPause,
				"next":       actionMediaNext,
				"previous":   actionMediaPrev,
				"stop":       actionMediaStop,
			}},
			Decl: MediaControlDecl(),
		},
		{
			Tool: systemTool{d: d, name: "screen.capture", verbs: map[string]Action{"": ActionScreenCapture}},
			Decl: ScreenCaptureDecl(),
		},
		{
			Tool: systemTool{d: d, name: "input.type", verbs: map[string]Action{"": ActionInputType}},
			Decl: InputTypeDecl(),
		},
	}
}
