//go:build windows

package main

// 票 30 收口残余③（30-v1 裁决 §9.3 的具名交回）：present.New(present.Channels{...})
// 一次装配，把票 30 落地的 internal/present 库（D10 结果呈现路由）接进常驻进程。
// 本文件只做装配与宿主适配器；internal/present、internal/agent、internal/risk、
// C17 桥、D43 表一个字节都不动。
//
// 逐通道（与 30-r1-impl.md §8① 的交回清单逐字对齐）：
//
//   - Speak → src.run.voiceAnnounce 的同款缝（wire-r1 在 resident_windows.go 里
//     接的 rv.announce）。ChTTS 出纳闸在 announce 路径上（26-r3），这里不重复设闸
//     ——端口文档（internal/present/dispatch.go）自己声明了这一点。
//   - Notify → rt.notify（runSpec 的通知缝，resident 装的是 postSystemNotification）。
//   - Badge → ball.SetBadgeText（STA 线程 post，任意 goroutine 可调）。
//   - TaskList → nil。面板任务表是快照派生的（PumpSources.Tasks 读名册），没有
//     单条入口 API（票 36 的面），具名跳过——present 对 nil 端口记 skipped 不炸。
//   - Panel → C17 桥结果展示：结果区（rt.stream）+ publishPanelSnapshot，桥一字未动。
//   - Clipboard → 真写适配器。端口文档的义务在这里兑现：写真写之前先
//     risk.Provenance.CheckText(scope, risk.ChClipboard, text)。
//   - Artifacts → *memory.Store 原样即插（PutArtifact/ArtifactsDir 方法集逐字吻合，零适配器）。
//
// 调用点在 execute()（cmd/wisp/run.go）：loop.Run 返回后对 completed 结果
// Dispatch(ctx, present.Input{TaskID, Text})。

import (
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"runtime"
	"unsafe"

	"github.com/CarlosShao/wisp/internal/present"
	"github.com/CarlosShao/wisp/internal/risk"
	"golang.org/x/sys/windows"
)

// presentClipboardScope is the C25 scope the clipboard gate reads, opened once
// at assembly on the process's own engine (the one the tool bridge marks into).
// It is the same shape the announce gate's voiceAnnounceScope documents
// (resident_voice_engines_windows.go): a process-lifetime scope, NOT a shadow of
// a task's scope, because admitTask's revocation closes a task's scope before
// loop.Run returns and CheckText fail-closes on an unbound id - a gate pointed
// at the closed task scope would refuse every write for the wrong reason. The
// named residual is the same one the announce path carries: a task's own taints
// are dropped when its scope closes, so carrying them to presentation time is a
// change in internal/risk's scope lifetime and outside this leg's write surface.
// What the gate does guarantee here is the port's obligation itself: no bytes
// reach the real clipboard without CheckText having run first, and a refused
// text is named, not swallowed.
const presentClipboardScope = "present-clipboard"

// presentClipboard is tier-4's real clipboard writer (D10 row 4). The gate call
// is unconditional and sits in front of the write half: write is the real Win32
// body in production and a recorder in the wiring tests, and neither can be
// reached while the gate refuses.
type presentClipboard struct {
	prov  *risk.Provenance
	scope string
	write func(string) error
	logf  func(string, ...any)
}

// newPresentClipboard opens the gate's scope and returns nil only when the gate
// itself cannot exist: no write half was handed in, or the scope would not open.
// Both are fail-closed shapes - a nil port is recorded as "not wired" by the
// router, and an unchecked write never happens.
func newPresentClipboard(prov *risk.Provenance, write func(string) error, logf func(string, ...any)) *presentClipboard {
	if write == nil {
		return nil
	}
	if prov == nil {
		// The same posture attachProvenance takes: an own engine keeps the gate
		// real (fail-closed) where a nil engine would have made CheckText a
		// nil dereference or the check a silent skip.
		prov = risk.NewProvenance(risk.ProvOptions{NoProbe: true})
	}
	if prov.OpenScope(presentClipboardScope) == nil {
		return nil
	}
	return &presentClipboard{prov: prov, scope: presentClipboardScope, write: write, logf: logf}
}

// WriteClipboard is present.ClipboardWriter. The port's documented obligation
// lives on the first two lines: CheckText against this host's scope and the
// clipboard channel, and a hit refuses the write with its source named.
func (c *presentClipboard) WriteClipboard(text string) error {
	if c == nil || c.prov == nil || c.write == nil {
		return errors.New("present: clipboard gate not assembled; refusing to write unchecked")
	}
	hit, refused := c.prov.CheckText(c.scope, risk.ChClipboard, text)
	if refused {
		if c.logf != nil {
			c.logf("present: 剪贴板写入被 ChClipboard 闸拦下（来源 %s）：%s", hit.Source(), hit.Origin)
		}
		return fmt.Errorf("present: clipboard refused by the %s gate (source %s)", risk.ChClipboard, hit.Source())
	}
	return c.write(text)
}

// ---- the real clipboard write (Win32) --------------------------------------

const (
	// clipGMEMFixed makes GlobalAlloc return the block's own address as the
	// handle: SetClipboardData takes it as the handle AND the bytes are copied
	// to it as the destination, so this package never needs GlobalLock - and
	// never needs the one conversion go vet's unsafeptr check refuses, a
	// uintptr turning back into a Pointer.
	clipGMEMFixed     = uintptr(0x0000)
	clipCFUnicodeText = uintptr(13)
)

var (
	clipUser32 = windows.NewLazySystemDLL("user32.dll")
	clipOpen   = clipUser32.NewProc("OpenClipboard")
	clipClose  = clipUser32.NewProc("CloseClipboard")
	clipEmpty  = clipUser32.NewProc("EmptyClipboard")
	clipSet    = clipUser32.NewProc("SetClipboardData")

	clipKernel32 = windows.NewLazySystemDLL("kernel32.dll")
	clipAlloc    = clipKernel32.NewProc("GlobalAlloc")
	clipFree     = clipKernel32.NewProc("GlobalFree")
	clipMove     = clipKernel32.NewProc("RtlMoveMemory")
)

// writeClipboardWindows puts one CF_UNICODETEXT blob on the system clipboard.
// On success the system owns the block; every failure path releases what it
// allocated and closes what it opened, so the clipboard is never left locked.
func writeClipboardWindows(text string) error {
	u16, err := windows.UTF16FromString(text)
	if err != nil {
		return fmt.Errorf("present: clipboard text is not UTF-16 encodable: %w", err)
	}
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if r, _, callErr := clipOpen.Call(0); r == 0 {
		return fmt.Errorf("present: OpenClipboard: %v", callErr)
	}
	defer clipClose.Call()
	if r, _, callErr := clipEmpty.Call(); r == 0 {
		return fmt.Errorf("present: EmptyClipboard: %v", callErr)
	}
	size := uintptr(len(u16)) * unsafe.Sizeof(u16[0])
	h, _, callErr := clipAlloc.Call(clipGMEMFixed, size)
	if h == 0 {
		return fmt.Errorf("present: GlobalAlloc: %v", callErr)
	}
	clipMove.Call(h, uintptr(unsafe.Pointer(&u16[0])), size)
	if r, _, callErr := clipSet.Call(clipCFUnicodeText, h); r == 0 {
		clipFree.Call(h)
		return fmt.Errorf("present: SetClipboardData: %v", callErr)
	}
	return nil
}

// ---- the channel adapters ---------------------------------------------------
// One thin type per port, so the closures the assembly hands present.Channels
// carry the real objects' own methods rather than re-typed copies of them.

type presenterNotifier func(title, body string) error

func (f presenterNotifier) Notify(title, body string) error { return f(title, body) }

type presenterSpeaker func(text string) bool

func (f presenterSpeaker) Speak(text string) bool { return f(text) }

type presenterBadger func(text string)

func (f presenterBadger) Badge(text string) { f(text) }

type presenterPanelSink func(text, artifact string)

func (f presenterPanelSink) ShowResult(text, artifact string) { f(text, artifact) }

// showPresentResult is the PanelSink adapter: the finished result reaches the
// panel's results section, the one C17-approved outbound face this process has
// (snapshot packets only - the four-method panel.* whitelist is untouched).
//
// The mapping this host makes, named: the result's own text is already in the
// results section by the time the router fires - the loop's sink appended every
// delta under the task id and closed the row at EvDone - so this adapter adds
// the one thing the streamed row cannot carry, the landed artifact's location
// (D10 rows 1-2: 落文件必须在面板可见), and then publishes. The row key is the
// artifact's own bare file name: present.ArtifactName minted it,
// memory.ValidArtifactName guards it, so no second stream-key spelling is
// invented here and a re-read of the same result stays the same row.
func (rt *agentRuntime) showPresentResult(text, artifact string) {
	if rt == nil || rt.stream == nil {
		return
	}
	if artifact != "" {
		rt.stream.Append(filepath.Base(artifact), "全文已存："+artifact)
	}
	rt.publishPanelSnapshot()
}

// assembleResidentPresenter is the ONE assembly of present.Channels for this
// process. Every nil a degraded boot left behind becomes a nil port, which the
// router records as skipped - the wiring never invents a stand-in. The returned
// line names what was really wired, in the boot report's own honesty convention.
func assembleResidentPresenter(src *residentTaskSource, rb *residentBall, rv *residentVoice,
	clipWrite func(string) error,
) (*present.Dispatcher, string) {
	if src == nil || src.run == nil {
		return nil, ""
	}
	run := src.run
	ch := present.Channels{Logger: slog.Default()}
	line := "present 装配："
	if run.notify != nil {
		ch.Notify = presenterNotifier(run.notify)
		line += "Notify "
	}
	if rv != nil && run.voiceAnnounce != nil {
		ch.Speak = presenterSpeaker(run.voiceAnnounce)
		line += "Speak "
	}
	if rb != nil && rb.b != nil {
		ch.Badge = presenterBadger(rb.b.SetBadgeText)
		line += "Badge "
	}
	ch.Panel = presenterPanelSink(run.showPresentResult)
	line += "Panel "
	if run.store != nil {
		ch.Artifacts = run.store
		line += "Artifacts "
	}
	if clip := newPresentClipboard(run.c25, clipWrite, run.auditf); clip != nil {
		ch.Clipboard = clip
		line += "Clipboard "
	}
	// TaskList stays nil on purpose: the panel task table is snapshot-derived
	// (PumpSources.Tasks) and has no per-entry entry point to hand over - the
	// named skip the handoff allows.
	line += "TaskList=nil"
	return present.New(ch), line
}
