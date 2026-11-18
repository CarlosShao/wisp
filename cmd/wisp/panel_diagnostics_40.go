package main

// Readface wave 5 (2026-11-18), ticket 40 AC#5's Go half: the real leg behind
// internal/panel's diagnostics.preview / diagnostics.export doors - the panel
// face of the ticket-45 diagnostics bundle engine.
//
// WHY THIS DOOR EXISTS. The engine has answered since ticket 45
// (internal/observe's PreviewDiagnosticsBundle / BuildDiagnosticsBundle pair,
// which share the one BundleOptions source - the preview's entry list is the
// builder's zip entry list by construction), the CLI face since deferred
// disposition 20261110 #5 (diagnostics_cmd.go: `wisp diagnostics preview` /
// `export` under the [privacy] diagnostics_opt_in total switch), and the
// panel had neither. SPEC-08 §5.2 names the diagnostics view's flow - preview
// checklist the user reviews BEFORE export - and ticket 40 AC#5 asks that the
// preview list exactly the bundle contents and the exported bundle exclude
// audio / keys / transcript-full. Ticket 194's roster census recorded
// diagnostics.export as "engine in, delivery hop missing"; this file is that
// hop.
//
// ONE ENUMERATE SOURCE, NOT TWO (the AC's hardest clause). The door owns no
// file list of its own: both handlers assemble the SAME BundleOptions the CLI
// assembles - diagnosticsBundleOptions180P(h.dir, c, outPath, "") - and hand
// it to the SAME Preview/Build pair internal/observe owns. The panel preview,
// the panel zip and the CLI's cannot drift, because there is only one options
// assembler and one enumerate pair in the process. A preview that spelled its
// own list would be exactly the lie the ticket-45 preview function was built
// as a sibling of the builder to refuse.
//
// THE GATE IS THE CLI'S. [privacy] diagnostics_opt_in (config.toml, schema
// default false) is the total switch: the panel door consults the SAME key in
// the SAME file the CLI gate reads, and false refuses both verbs by name
// naming the key and the file - a diagnostics face that exists while the user
// said "no diagnostics" is fail-rude, and a panel face answering where the
// CLI refuses would be two privacy postures wearing one name. The door reads
// nothing else out of the config beyond what diagnosticsBundleOptions180P
// consumes (redact_paths).
//
// PREVIEW FIRST, ALWAYS. The export handler runs PreviewDiagnosticsBundle
// before BuildDiagnosticsBundle on the same options, exactly like the CLI's
// export: a preview failure aborts the export before anything is written (the
// builder is fail-closed and would leave no zip, but the preview is what
// makes "nothing written" a checked claim rather than a hope), and the
// payload the page renders IS that preview - so what the user saw before
// clicking export is the same object the builder ships.
//
// THE PAYLOAD IS THE ENGINE'S OWN. observe.BundlePreview carries JSON tags
// (created_at / out_path / categories[] / entries / total_bytes); the wire
// half re-marshals it compact and rides the settings-route shape: one prose
// sentence plus ONE payload line behind the door's marker, the page finds the
// marker and parses the rest of the line. The prose is never parsed. The
// excluded categories (transcripts / watchdog / audio) ride the payload with
// the engine's own reasons verbatim - "the bundle excludes audio / keys /
// transcript-full" is a fact the engine states, not a sentence this file
// invents.

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/panel"
)

// DiagnosticsPreviewPayloadMarker / DiagnosticsExportPayloadMarker start the
// two doors' payload lines. Same naming properties as every marker in the
// family (wisp-privacy-export:, wisp-session-rules:, wisp-cost-summary:): no
// `panel.` prefix and no dot - neither is a method name nor a route.
const (
	DiagnosticsPreviewPayloadMarker = "wisp-diagnostics-preview:"
	DiagnosticsExportPayloadMarker  = "wisp-diagnostics-export:"
)

// diagnosticsDoor40 is the wave-5 socket of the panel chain: the data root
// the gate reads and the bundle assembles over, the audit sink, the actor.
type diagnosticsDoor40 struct {
	dir   string
	audit panel.AuditFunc
	actor string
	now   func() time.Time // injectable for tests; nil = time.Now
}

// newDiagnosticsDoor40 builds the leg over the data root. No file is read
// here - every read happens on the request hop, gated.
func newDiagnosticsDoor40(dir string, auditf panel.AuditFunc, actor string) *diagnosticsDoor40 {
	return &diagnosticsDoor40{dir: dir, audit: auditf, actor: actor}
}

// gate40 consults [privacy] diagnostics_opt_in - the CLI's total switch, the
// same key in the same file diagnosticsGate180P reads. Three refusal shapes,
// three different fixes, the CLI gate's own discipline: no file, unreadable
// file, switch off. The returned config is what diagnosticsBundleOptions180P
// consumes (redact_paths).
func (h *diagnosticsDoor40) gate40() (*config.Config, error) {
	cfgPath := filepath.Join(h.dir, configFileName)
	if _, statErr := os.Stat(cfgPath); statErr != nil {
		return nil, fmt.Errorf(
			"总闸读不到 %s（%v）：[privacy] diagnostics_opt_in 取 schema 默认值 false，按关闸处理（fail-closed）；先创建配置并把 [privacy] 段的 diagnostics_opt_in 设为 true",
			cfgPath, statErr)
	}
	c, _, lerr := config.LoadFile(cfgPath, nil)
	if lerr != nil || c == nil {
		return nil, fmt.Errorf(
			"%s 读不进来（%v）：总闸读不到就是关着，fail-closed；先修好配置文件", cfgPath, lerr)
	}
	if !c.Privacy.DiagnosticsOptIn {
		return nil, fmt.Errorf(
			"总闸 [privacy] diagnostics_opt_in = false（来源 %s）。面板诊断面与 wisp diagnostics CLI 是同一把总闸：关着时不做预览、不读任何日志；要导出诊断包，请把 config.toml 里 [privacy] 段的 diagnostics_opt_in 改为 true 后重试。",
			cfgPath)
	}
	return c, nil
}

// diagnosticsOutPath40 is the export default, the CLI's own shape: one stamped
// zip inside the data root. The page echoes the path verbatim; the zip is a
// file the user opens from disk, not bytes over the bridge.
func (h *diagnosticsDoor40) diagnosticsOutPath40() string {
	stamp := h.nowStamp()
	return filepath.Join(h.dir, "wisp-diagnostics-"+stamp+".zip")
}

func (h *diagnosticsDoor40) nowStamp() string {
	if h.now != nil {
		return h.now().UTC().Format("20060102-150405Z")
	}
	return observe.NowWallUTC().Format("20060102-150405Z")
}

// HandleDiagnosticsPreview answers diagnostics.preview: gate, the CLI's own
// options assembler, the engine's own preview, the engine's own JSON.
// Nothing is written anywhere - the builder is not called.
func (h *diagnosticsDoor40) HandleDiagnosticsPreview(ctx context.Context, req panel.ComposerRequest) (string, error) {
	c, err := h.gate40()
	if err != nil {
		h.record(req, "preview", "refused", err)
		return "", fmt.Errorf("panel: diagnostics.preview 被拒绝（requestId=%q）：%w", req.RequestID, err)
	}
	o, aerr := diagnosticsBundleOptions180P(h.dir, c, "(preview only: nothing is written)", "")
	if aerr != nil {
		h.record(req, "preview", "options-failed", aerr)
		return "", fmt.Errorf("panel: diagnostics.preview 被拒绝（requestId=%q）：%w", req.RequestID, aerr)
	}
	p, perr := observe.PreviewDiagnosticsBundle(o)
	if perr != nil {
		// The preview failing means the bundle really could not be built -
		// the engine fails on the same inputs the builder fails on.
		h.record(req, "preview", "preview-failed", perr)
		return "", fmt.Errorf("panel: diagnostics.preview 被拒绝（requestId=%q）：%w", req.RequestID, perr)
	}
	b, merr := json.Marshal(p)
	if merr != nil {
		h.record(req, "preview", "marshal-failed", merr)
		return "", fmt.Errorf("panel: diagnostics.preview 载荷编不出来（requestId=%q）：%w", req.RequestID, merr)
	}
	h.record(req, "preview", "answered", nil)
	return fmt.Sprintf(
		"诊断包预览（条目 %d 个，合计 %d 字节）。预览与导出走同一份文件清单来源：同一个 BundleOptions 组装、同一对 Preview/Build 函数（票 45 引擎），所以这份清单就是导出 zip 的条目清单。排除类（转写/watchdog/音频）按引擎自报的默认排除原因逐类列出。",
		len(p.Entries), p.TotalBytes) + "\n" + DiagnosticsPreviewPayloadMarker + string(b), nil
}

// HandleDiagnosticsExport answers diagnostics.export: gate, the CLI's own
// options assembler with a real OutPath, the preview FIRST (a failure aborts
// before anything is written), then the builder. The receipt names the path;
// the payload is the same preview JSON whose out_path is now the real target.
func (h *diagnosticsDoor40) HandleDiagnosticsExport(ctx context.Context, req panel.ComposerRequest) (string, error) {
	c, err := h.gate40()
	if err != nil {
		h.record(req, "export", "refused", err)
		return "", fmt.Errorf("panel: diagnostics.export 被拒绝（requestId=%q）：%w", req.RequestID, err)
	}
	outPath := h.diagnosticsOutPath40()
	o, aerr := diagnosticsBundleOptions180P(h.dir, c, outPath, "")
	if aerr != nil {
		h.record(req, "export", "options-failed", aerr)
		return "", fmt.Errorf("panel: diagnostics.export 被拒绝（requestId=%q）：%w", req.RequestID, aerr)
	}
	p, perr := observe.PreviewDiagnosticsBundle(o)
	if perr != nil {
		h.record(req, "export", "preview-failed", perr)
		return "", fmt.Errorf("panel: diagnostics.export 被拒绝（requestId=%q，导出中止，什么都没写）：%w", req.RequestID, perr)
	}
	written, berr := observe.BuildDiagnosticsBundle(o)
	if berr != nil {
		// The builder is fail-closed: an aborted bundle leaves no zip behind.
		h.record(req, "export", "build-failed", berr)
		return "", fmt.Errorf("panel: diagnostics.export 被拒绝（requestId=%q，zip 已被构建器回收，盘上没有半成品）：%w", req.RequestID, berr)
	}
	if written != outPath {
		werr := fmt.Errorf("构建器回答的路径与请求的不一致（%q != %q）", written, outPath)
		h.record(req, "export", "path-mismatch", werr)
		return "", fmt.Errorf("panel: diagnostics.export 被拒绝（requestId=%q）：%w", req.RequestID, werr)
	}
	b, merr := json.Marshal(p)
	if merr != nil {
		h.record(req, "export", "marshal-failed", merr)
		return "", fmt.Errorf("panel: diagnostics.export 载荷编不出来（requestId=%q）：%w", req.RequestID, merr)
	}
	h.record(req, "export", "answered", nil)
	return fmt.Sprintf(
		"已导出诊断包：%s（条目 %d 个；内容以预览与包内 manifest.json 为准——manifest 逐条记录了每个条目与每个未收集项的原因）。音频、密钥、完整转写按引擎默认排除；转写与 watchdog 快照默认排除、逐项同意。",
		written, len(p.Entries)) + "\n" + DiagnosticsExportPayloadMarker + string(b), nil
}

func (h *diagnosticsDoor40) record(req panel.ComposerRequest, verb, outcome string, err error) {
	if h == nil || h.audit == nil {
		return
	}
	line := "ok"
	if err != nil {
		line = err.Error()
	}
	h.audit("panel: DIAGNOSTICS request=%q method=%q verb=%q actor=%q outcome=%q err=%s",
		req.RequestID, req.Method, verb, h.actor, outcome, line)
}
