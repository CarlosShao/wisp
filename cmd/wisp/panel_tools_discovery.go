package main

// Readface wave 3 (2026-11-17), ticket 38 AC#3's discovery/launch doors and
// AC#5's force-rules read: the real legs behind internal/panel's tools.list /
// tool.launch / rules.force (the C17 roster additions' reasons live on
// bridge.go's Method* constants; here are the bodies).
//
// THE LAUNCH'S PATH, stated once: a tool.launch request is refused unless the
// registry holds the named tool; the accepted request is composed into a task
// book and handed to the ONE task slot the composition root owns, exactly the
// panel-message machinery (453). The loop then dispatches the tool call the
// way every model-driven call travels - internal/agent/loop.go's single
// agent.ToolRequest construction point (ticket 309's frozen nail pins the
// count; a second construction site in this package would trip it, and the
// dispatch's "same host bridge, no bypass" is satisfied a fortiori: the
// launch rides THE path, not a copy of it). An L2 verdict raises the REAL
// approval card; the receipt does NOT wait for it (the roster-cancel
// precedent: a dispatch goroutine parked on a 300s card stalls every other
// panel request).
//
// THE ROSTER'S PILL IS THE TRUTH: riskPill is rendered from the entry's own
// declared level (internal/tools' Decl.Declared, the R1 floor the D34 table
// pins) through internal/risk's Level.String() - the same vocabulary the
// confirmation card shows. No level word is spelled in this file.
//
// THE FORCE FACE HAS ONE RESOLVER: the view below is resolveForce's answer
// (force_rules.go), the same function execute() consults before attaching the
// instruction block - the face and the next prompt cannot disagree.

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/panel"
)

// payload markers, the history-page family's naming properties: no "panel."
// prefix, no dot.
const (
	toolsListPayloadMarker   = "wisp-tools-list:"
	rulesForcePayloadMarker  = "wisp-rules-force:"
	toolLaunchReceiptMarker  = "wisp-tool-launch:"
	maxToolLaunchArgsRunes   = 2048
	forceModeVocabularyW3    = "force_tool|force_chat"
	toolLaunchTaskPrefixW3   = "工具启动（面板）：请调用工具 "
	toolLaunchTaskArgsNoteW3 = "，参数用上面给出的 JSON。只做这一次调用，然后把结果回报给我。"
)

// ---------------------------------------------------------------------------
// tools.list
// ---------------------------------------------------------------------------

// toolLaunchSeam is the one thing the launch door needs from the task
// pipeline: hand it the tool name and the JSON args the page echoed, and it
// submits the composed task book through the resident task slot. A nil seam
// is the CLI / refused-pipeline state and refuses by name.
type toolLaunchSeam func(ctx context.Context, tool, args string) error

// toolsListSeam returns the roster rows and the source that produced them.
type toolsListSeam func() []panel.ToolRosterRow

// rulesForceSeam resolves the face's view for one task text.
type rulesForceSeam func(taskText string) (panel.ForceRulesView, error)

// discoveryDoors is the wave-3 socket of the panel chain: three doors, one
// struct, three late-bound seams (the roster-cancel/453 pattern - the chain
// is assembled before the task pipeline and the runtime exist).
type discoveryDoors struct {
	audit panel.AuditFunc
	actor string

	mu         sync.RWMutex
	toolsSeam  toolsListSeam
	launchSeam toolLaunchSeam
	rulesSeam  rulesForceSeam
}

func newDiscoveryDoors(audit panel.AuditFunc, actor string) *discoveryDoors {
	return &discoveryDoors{audit: audit, actor: actor}
}

func (d *discoveryDoors) bind(tools toolsListSeam, launch toolLaunchSeam, rules rulesForceSeam) {
	if d == nil {
		return
	}
	d.mu.Lock()
	d.toolsSeam, d.launchSeam, d.rulesSeam = tools, launch, rules
	d.mu.Unlock()
}

// HandleToolsList answers one parsed roster read: the seam's rows (the
// registry's own truth) behind the payload marker. Read-only end to end.
func (d *discoveryDoors) HandleToolsList(ctx context.Context, req panel.ComposerRequest) (string, error) {
	if d == nil {
		return "", fmt.Errorf("%w: 工具名册门不存在", panel.ErrNoHandlerAttached)
	}
	d.mu.RLock()
	seam := d.toolsSeam
	d.mu.RUnlock()
	if seam == nil {
		err := fmt.Errorf("%w: 工具名册未接入（本进程没有任务管线）：名册只读不执行，这一问没有答案，也不会被静默丢掉",
			panel.ErrNoHandlerAttached)
		d.record(req, "tools.list", "refused", err)
		return "", err
	}
	rows := seam()
	d.record(req, "tools.list", "answered", nil)
	var b strings.Builder
	fmt.Fprintf(&b, "已安装工具共 %d 只（riskPill 为 C19 风险级真值，只读名册；启动走 tool.launch，照样过审批门）。", len(rows))
	payload, merr := json.Marshal(map[string]any{"rows": rows})
	if merr != nil {
		b.WriteString("\n" + toolsListPayloadMarker + fmt.Sprintf(`{"payloadError":%q}`, merr.Error()))
		return b.String(), nil
	}
	b.WriteString("\n" + toolsListPayloadMarker + string(payload))
	return b.String(), nil
}

// HandleToolLaunch answers one parsed launch: validate the selectors, hand
// the invocation to the seam, receipt without waiting. The receipt names the
// tool; the card (if the verdict is L2) is the gate's own.
func (d *discoveryDoors) HandleToolLaunch(ctx context.Context, req panel.ComposerRequest) (string, error) {
	if d == nil {
		return "", fmt.Errorf("%w: 工具启动门不存在", panel.ErrNoHandlerAttached)
	}
	tool := strings.TrimSpace(req.ToolName)
	if tool == "" {
		err := fmt.Errorf("%w: toolName 是空的：launch 启动的是名册里的一只工具，不给名字就什么都不启动", panel.ErrComposerRequest)
		d.record(req, tool, "refused", err)
		return "", err
	}
	if strings.TrimSpace(req.ToolArgs) != "" && !json.Valid([]byte(req.ToolArgs)) {
		err := fmt.Errorf("%w: toolArgs 不是合法的 JSON：%s", panel.ErrComposerRequest, truncateRunesW3(req.ToolArgs, 64))
		d.record(req, tool, "refused", err)
		return "", err
	}
	d.mu.RLock()
	seam := d.launchSeam
	d.mu.RUnlock()
	if seam == nil {
		err := fmt.Errorf("%w: 启动通路未接入（本进程没有任务管线）：工具 %s 未被触碰，也不会被静默丢掉",
			panel.ErrNoHandlerAttached, tool)
		d.record(req, tool, "refused", err)
		return "", err
	}
	if err := seam(ctx, tool, req.ToolArgs); err != nil {
		err = fmt.Errorf("panel: 工具 %s 没有启动：%w", tool, err)
		d.record(req, tool, "refused", err)
		return "", err
	}
	d.record(req, tool, "launched", nil)
	// The receipt does NOT wait for the card or the result (see the file
	// header): the launch rides the task slot, the gate's card and the tool's
	// outcome surface through the channels that already exist.
	return fmt.Sprintf("已发出工具 %s 的启动（经既有任务管线；若风险判定为 L2 会弹真实审批卡，结果与取证按既有通道呈现）。", tool), nil
}

// HandleRulesForce answers one parsed force read: the resolver's view - the
// effective rules with their sources, plus the mode the NEXT prompt will
// adopt for the requested task text. Read-only: nothing here stages or
// writes.
func (d *discoveryDoors) HandleRulesForce(ctx context.Context, req panel.ComposerRequest) (string, error) {
	if d == nil {
		return "", fmt.Errorf("%w: force 规则门不存在", panel.ErrNoHandlerAttached)
	}
	d.mu.RLock()
	seam := d.rulesSeam
	d.mu.RUnlock()
	if seam == nil {
		err := fmt.Errorf("%w: force 规则未接入（本进程没有任务管线）：这一问没有答案，也不会被静默丢掉",
			panel.ErrNoHandlerAttached)
		d.record(req, "rules.force", "refused", err)
		return "", err
	}
	view, err := seam(req.ForceTask)
	if err != nil {
		err = fmt.Errorf("panel: rules.force 读取失败：%w", err)
		d.record(req, "rules.force", "refused", err)
		return "", err
	}
	d.record(req, "rules.force", "answered", nil)
	var b strings.Builder
	if view.Mode == "" {
		b.WriteString("当前没有生效的强制规则，下一 prompt 采用默认模式（不注入 force 指令）。")
	} else {
		fmt.Fprintf(&b, "下一 prompt 采用 %s（来源 %s）：%s", view.Mode, view.Source, view.Instruction)
	}
	payload, merr := json.Marshal(view)
	if merr != nil {
		b.WriteString("\n" + rulesForcePayloadMarker + fmt.Sprintf(`{"payloadError":%q}`, merr.Error()))
		return b.String(), nil
	}
	b.WriteString("\n" + rulesForcePayloadMarker + string(payload))
	return b.String(), nil
}

func (d *discoveryDoors) record(req panel.ComposerRequest, what, outcome string, err error) {
	if d == nil || d.audit == nil {
		return
	}
	line := "ok"
	if err != nil {
		line = err.Error()
	}
	d.audit("panel: DISCOVERY-FORCE request=%q what=%q actor=%q outcome=%q err=%s",
		req.RequestID, what, d.actor, outcome, line)
}

func truncateRunesW3(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + "…"
}

// ---------------------------------------------------------------------------
// the seams' production implementations, on the composition root
// ---------------------------------------------------------------------------

// toolsRosterRows is the roster seam: the registry's own entries, one row
// each, the pill rendered by internal/risk's vocabulary (never spelled here).
func (rt *agentRuntime) toolsRosterRows() []panel.ToolRosterRow {
	if rt == nil || rt.bridge == nil {
		return nil
	}
	entries := rt.bridge.Registry().List()
	rows := make([]panel.ToolRosterRow, 0, len(entries))
	for _, e := range entries {
		rows = append(rows, panel.ToolRosterRow{
			Name:      e.Tool.Name(),
			Title:     toolTitleW3(e.Tool.Description()),
			RiskPill:  e.Decl.Declared.String(),
			ArgSummary: toolArgSummaryW3(e.Tool.Parameters()),
		})
	}
	return rows
}

// toolTitleW3 keeps the description's first line - the sentence the model
// itself reads to decide whether to call the tool.
func toolTitleW3(desc string) string {
	line := strings.TrimSpace(desc)
	if i := strings.IndexAny(line, "\r\n"); i >= 0 {
		line = strings.TrimSpace(line[:i])
	}
	return line
}

// toolArgSummaryW3 flattens the tool's JSON Schema top-level properties into
// one comma-separated selector list. A schema without named properties is
// "无参数", never an invented summary.
func toolArgSummaryW3(schema json.RawMessage) string {
	if len(schema) == 0 {
		return "无参数"
	}
	var parsed struct {
		Properties map[string]json.RawMessage `json:"properties"`
	}
	if err := json.Unmarshal(schema, &parsed); err != nil || len(parsed.Properties) == 0 {
		return "无参数"
	}
	names := make([]string, 0, len(parsed.Properties))
	for k := range parsed.Properties {
		names = append(names, k)
	}
	sort.Strings(names)
	return strings.Join(names, ",")
}

// composeToolLaunch turns an accepted launch into the task book text, judged
// against the registry the run actually assembled: a name the registry does
// not hold is refused HERE (fail-closed), never silently submitted to a model
// that may improvise one. The args travel inside the book verbatim, capped at
// the forensics cap's own order of magnitude.
func (rt *agentRuntime) composeToolLaunch(tool, args string) (string, error) {
	if rt == nil || rt.bridge == nil {
		return "", fmt.Errorf("本进程没有任务管线")
	}
	name := strings.TrimSpace(tool)
	if _, ok := rt.bridge.Registry().Lookup(name); !ok {
		return "", fmt.Errorf("工具 %s 不在本机名册里（tools.list 是唯一的启动菜单）", name)
	}
	if len([]rune(args)) > maxToolLaunchArgsRunes {
		return "", fmt.Errorf("toolArgs 超长（%d 字符 > %d），面板启动不受理", len([]rune(args)), maxToolLaunchArgsRunes)
	}
	argText := strings.TrimSpace(args)
	if argText == "" {
		argText = "{}"
	}
	return toolLaunchTaskPrefixW3 + name + toolLaunchTaskArgsNoteW3 + "\n参数 JSON：" + argText, nil
}

// SubmitPanelToolLaunch lives in panel_discovery_force_windows.go with the
// resident source it submits through (the windows-only task pipeline, the
// 453/bind pattern).

// forceRulesView resolves the face's view. The mode is the one resolveForce
// computed for the requested task text - the SAME decision execute() applies
// to the next prompt.
func (rt *agentRuntime) forceRulesView(taskText string) (panel.ForceRulesView, error) {
	if rt == nil {
		return panel.ForceRulesView{Source: "default"}, fmt.Errorf("本进程没有任务管线")
	}
	dec, err := rt.resolveForce(taskText)
	if err != nil {
		return panel.ForceRulesView{}, err
	}
	view := panel.ForceRulesView{
		Mode:        dec.Mode,
		Source:      dec.Source,
		Pattern:     dec.Pattern,
		Instruction: dec.Instruction,
	}
	if rules, rerr := rt.force.loadRules(); rerr != nil {
		return view, rerr
	} else {
		for _, r := range rules {
			view.Rules = append(view.Rules, panel.ForceRuleRow{Source: "toml", Pattern: r.Pattern, Mode: r.Mode})
		}
	}
	for id, mode := range rt.activeOverrides() {
		view.Rules = append(view.Rules, panel.ForceRuleRow{Source: "override", Pattern: id, Mode: mode})
	}
	sort.SliceStable(view.Rules, func(i, j int) bool {
		if view.Rules[i].Source != view.Rules[j].Source {
			return view.Rules[i].Source < view.Rules[j].Source
		}
		return view.Rules[i].Pattern < view.Rules[j].Pattern
	})
	return view, nil
}

// compile-time proofs: the doors satisfy the roster shapes internal/panel
// declared, and the instruction source satisfies the ticket-200 seam.
var (
	_ panel.ToolsListHandler  = (*discoveryDoors)(nil)
	_ panel.ToolLaunchHandler = (*discoveryDoors)(nil)
	_ panel.RulesForceHandler = (*discoveryDoors)(nil)
	_ agent.InstructionSource = (*forceInstructionSource)(nil)
)
