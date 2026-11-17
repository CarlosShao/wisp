package main

// Ticket 38 AC#5's Go half (readface wave 3, 2026-11-17): the force-rule
// resolver - the ONE object both consumers of D11③ read, so the face and the
// prompt cannot drift.
//
// WHAT D11③ IS, verbatim from its sources (read-only for this leg): D11③
// (PLAN.md:299-301) demotes the draft's 6.5 rule system to a USER PREFERENCE
// injection point - "用户显式偏好覆盖"; SPEC-05 §2 names the mechanism:
// "注入为指令而非前置分类" (injected as an INSTRUCTION, never a
// pre-classification layer). The D11③ 欠账行 in SPEC-12 §5 (the row that
// matches the marker at internal/agent/control.go:15) states the landing
// condition for the full routing-side ticket ("规则表在 config.toml 落地且
// 路由生效") and records that today no ticket owns it. THIS leg does not land
// that row: internal/agent/control.go's欠账标记 is untouched, internal/config's
// schema is read-only here, and config.toml's
// strict loader (DisallowUnknownFields) refuses an unknown key by design - so
// a force key written into config.toml [agent] could never coexist with a
// loading config. The rule table therefore lives in its OWN file beside
// config.toml: force-rules.toml. That is a composition-root decision, in
// cmd/wisp's write face, and it changes no frozen surface; the day the
// routing-side ticket lands the key in config.toml, this file's loader is the
// one place that moves.
//
// THE TWO SOURCES, and the precedence the teeth pin:
//
//	override - the palette's per-task chip (ticket 38 AC#5), staged through
//	           panel.message.send's own payload (ComposerRequest.ForceMode)
//	           and consumed at admission; one slot, the 453 pattern;
//	toml     - [[force_rule]] rows in force-rules.toml, matched by substring
//	           against the task text (D11③'s rule-table shape: a pattern
//	           steering the mode for matching tasks);
//	default  - none of the above: no force instruction, the prompt is exactly
//	           what it was before this file existed (byte-identical).
//
// The vocabulary is D11③'s own two words - force_tool (必须调工具) and
// force_chat (不得调工具) - never a third one minted here.

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	toml "github.com/pelletier/go-toml/v2"

	"github.com/CarlosShao/wisp/internal/agent"
)

// ForceRulesFileName is the rule table's file name, beside config.toml.
const ForceRulesFileName = "force-rules.toml"

// The D11③ mode vocabulary. Two words, both from PLAN.md:299.
const (
	ForceModeTool = "force_tool"
	ForceModeChat = "force_chat"
)

// forceInstruction blocks are the D11③ instructions the assembler injects -
// imperative, one paragraph, no rule engine prose. They ride the ticket-200
// InstructionSource seam (the guidance/authority split ticket 200 AC#4
// pinned), NOT one of the seven frozen D39 sections and NOT the safety
// section: the section table is frozen by D39 and a per-task block does not
// belong in it.
const (
	forceToolInstruction = "任务模式强制（用户显式偏好，优先于你的默认判断）：本任务必须通过调用工具完成，不要只给文字回答。"
	forceChatInstruction = "任务模式强制（用户显式偏好，优先于你的默认判断）：本任务直接对话回答，不要调用任何工具。"
)

// forceRule is one TOML row.
type forceRule struct {
	// Pattern is the substring the task text is matched against. Empty
	// matches nothing: a rule that names no pattern is dead, not wildcard.
	Pattern string `toml:"pattern"`
	// Mode is force_tool | force_chat; anything else is refused at load.
	Mode string `toml:"mode"`
}

type forceRulesFile struct {
	Rules []forceRule `toml:"force_rule"`
}

// forceStore is the runtime's one holder: the TOML rules (re-read per
// resolve, so an edit lands without a restart) and the override slots. All
// reads go through resolveForce; nothing here writes a prompt directly.
type forceStore struct {
	// path is the rules file's absolute path; "" means no table was located,
	// which is the honest default (no rules, no instruction).
	path string
	mu   sync.Mutex
	// pendingOverride is the palette chip staged by the message door for the
	// NEXT admitted task (one slot, the 453 pattern).
	pendingOverride string
	// active maps task id -> the override that rode its admission, so the
	// face can name which task runs forced.
	active map[string]string
}

func newForceStore(path string) *forceStore {
	return &forceStore{path: path, active: map[string]string{}}
}

// loadRules reads the TOML table. A missing file is the default state (no
// rules), NOT an error; a malformed file or an unknown mode word refuses the
// WHOLE table (the ticket-83 lying-key rule: a half-loaded rule table would
// let one bad row silently steer nothing while the face reports the rest).
func (s *forceStore) loadRules() ([]forceRule, error) {
	if s == nil || s.path == "" {
		return nil, nil
	}
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("force-rules.toml 读取失败：%w", err)
	}
	var f forceRulesFile
	if err := toml.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("force-rules.toml 不是合法的规则表：%w", err)
	}
	for i, r := range f.Rules {
		if r.Mode != ForceModeTool && r.Mode != ForceModeChat {
			return nil, fmt.Errorf("force-rules.toml 第 %d 条规则的 mode=%q 不在 force_tool|force_chat 里，整表拒绝（不许猜）", i+1, r.Mode)
		}
	}
	return f.Rules, nil
}

// stageOverride arms the palette chip for the next admitted task. The
// vocabulary is judged by the caller (the message door refuses before it
// stages); this slot takes the word verbatim.
func (rt *agentRuntime) stageForceOverride(mode string) {
	if rt == nil || rt.force == nil {
		return
	}
	rt.force.mu.Lock()
	rt.force.pendingOverride = mode
	rt.force.mu.Unlock()
}

// takeOverride consumes the staged chip at admission (admitTask, the first
// hook that receives the loop-minted task id - the 453 pattern exactly).
func (rt *agentRuntime) takeForceOverride(taskID string) string {
	if rt == nil || rt.force == nil {
		return ""
	}
	rt.force.mu.Lock()
	mode := rt.force.pendingOverride
	rt.force.pendingOverride = ""
	if mode != "" && taskID != "" {
		if rt.force.active == nil {
			rt.force.active = map[string]string{}
		}
		rt.force.active[taskID] = mode
	}
	rt.force.mu.Unlock()
	return mode
}

// clearForceOverride is the refused-submission cleanup: a chip staged for a
// task that never started must not leak into the next one (the 453 rule).
func (rt *agentRuntime) clearForceOverride() {
	if rt == nil || rt.force == nil {
		return
	}
	rt.force.mu.Lock()
	rt.force.pendingOverride = ""
	rt.force.mu.Unlock()
}

// ForceDecision is resolveForce's answer, and the face's payload: which mode
// the next prompt will instruct, WHO said so, the exact instruction block the
// assembler will inject ("" = none), and the winning rule's pattern.
type ForceDecision struct {
	Mode        string // "force_tool" | "force_chat" | "" (default: no force)
	Source      string // "override" | "toml" | "default"
	Pattern     string // the toml rule's pattern ("" for override/default)
	Instruction string // the prompt block ("" = the prompt is unchanged)
}

// resolveForce is the ONE evaluation both the face and the prompt injection
// call. Precedence: override > toml > default. The TOML table is re-read per
// call, so an edit to the file takes effect on the next task without a
// restart; a broken table resolves to default WITH the error surfaced to the
// caller (never a silent half-table).
func (rt *agentRuntime) resolveForce(taskText string) (ForceDecision, error) {
	dec := ForceDecision{Source: "default"}
	if rt == nil || rt.force == nil {
		return dec, nil
	}
	rt.force.mu.Lock()
	ov := rt.force.pendingOverride
	rt.force.mu.Unlock()
	if ov == ForceModeTool || ov == ForceModeChat {
		dec.Mode, dec.Source = ov, "override"
		dec.Instruction = forceInstructionFor(ov)
		return dec, nil
	}
	rules, err := rt.force.loadRules()
	if err != nil {
		return dec, err
	}
	for _, r := range rules {
		if r.Pattern != "" && strings.Contains(taskText, r.Pattern) {
			dec.Mode, dec.Source, dec.Pattern = r.Mode, "toml", r.Pattern
			dec.Instruction = forceInstructionFor(r.Mode)
			return dec, nil
		}
	}
	return dec, nil
}

// activeOverrides is the face's override rows: the tasks this process
// admitted while a chip was staged.
func (rt *agentRuntime) activeOverrides() map[string]string {
	if rt == nil || rt.force == nil {
		return nil
	}
	rt.force.mu.Lock()
	defer rt.force.mu.Unlock()
	out := make(map[string]string, len(rt.force.active))
	for id, mode := range rt.force.active {
		out[id] = mode
	}
	return out
}

func forceInstructionFor(mode string) string {
	switch mode {
	case ForceModeTool:
		return forceToolInstruction
	case ForceModeChat:
		return forceChatInstruction
	default:
		return ""
	}
}

// forceRulesPathFor derives the table's path from the config file's own
// directory (the same dir layout the ticket-200 global tier uses), so the
// rules live where the operator already edits config.toml.
func forceRulesPathFor(cfgPath string) string {
	if strings.TrimSpace(cfgPath) == "" {
		return ""
	}
	return filepath.Join(filepath.Dir(cfgPath), ForceRulesFileName)
}

// forceInstructionSource wraps the ticket-200 loader: the force block rides
// the SAME InstructionSource seam (guidance, not a D39 section), prepended so
// the steering sentence is the first thing the turn's guidance says. An empty
// block renders exactly what the inner source rendered - the default path is
// byte-identical to before this file existed.
type forceInstructionSource struct {
	inner agent.InstructionSource
	block string
}

func (f *forceInstructionSource) BlockForTurn() string {
	var b strings.Builder
	if f != nil && f.block != "" {
		b.WriteString(f.block)
		b.WriteString("\n")
	}
	if f != nil && f.inner != nil {
		b.WriteString(f.inner.BlockForTurn())
	}
	return b.String()
}

func (f *forceInstructionSource) ManifestForTurn() []string {
	var out []string
	if f != nil && f.block != "" {
		out = append(out, "force-mode: "+forceModeWordFor(f.block))
	}
	if f != nil && f.inner != nil {
		out = append(out, f.inner.ManifestForTurn()...)
	}
	return out
}

func forceModeWordFor(block string) string {
	if strings.Contains(block, "不要调用任何工具") {
		return ForceModeChat
	}
	return ForceModeTool
}
