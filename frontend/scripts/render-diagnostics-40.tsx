/* ============================================================================
   Ticket 40 AC#5's behavioural tooth (render evidence), 2026-11-18 - the
   diagnostics page's preview / export leg.
   ----------------------------------------------------------------------------
   WHAT THIS FILE PROVES, in the shape render-privacy-40 established: every
   check below feeds REAL-SHAPED data through the PRODUCTION component and
   asserts the SAME VALUE comes back out in the markup - "喂真 payload ⇒
   行内真字段逐字上屏".

   The fixtures are the wire's own documented shapes, copied from the Go leg
   that answers them (cmd/wisp/panel_diagnostics_40.go, the doors bridge.go's
   MethodDiagnosticsPreview / MethodDiagnosticsExport name, composer_dispatch.go
   routes, panel_inbound.go wires) and from the engine that owns the payload
   (internal/observe's BundlePreview JSON tags - the SAME object the builder
   ships, compact-marshalled behind the door's marker):
     - preview reply: one prose sentence plus the ONE payload line behind
       wisp-diagnostics-preview:, whose JSON carries categories (each with the
       engine's own detail and items name/bytes/head/included/reason), the
       exact entries list, and total_bytes;
     - the ZERO-ENTRY preview (tooth ③'s green half): entries [] with the
       excluded categories' own statements - the page must render 条目 0 个
       OUT LOUD with those statements verbatim and draw NO entry-list chrome;
       an empty list that reads as a normal reading is exactly what this tooth
       refuses;
     - export receipt: HandleDiagnosticsExport's sentence naming the written
       path, plus the same preview JSON whose out_path is the real target;
     - a refusal: RefusedEnvelopeForUser's wrapper sentence, including the
       switch-off one that NAMES diagnostics_opt_in - shown verbatim, never an
       empty list that reads as "nothing to export".

   Where a reading cannot be had, the check pins the 未判 rule: no payload
   line ⇒ the page says 未判 and NEVER draws a fabricated entry or path; a
   payload whose entries are damaged ⇒ skipped count stated, zeros never
   folded in (the cost-view rule).

   Source-level teeth:
     - the two route literals live in panel.ts (the ONE door file), spelled
       exactly as bridge.go's constants carry them;
     - the page renders the payload's own exclusion statements and never
       spells the engine's sentences itself (a page-side exclusion sentence
       would be a fabricated reading - mutation ③'s target);
     - the export receipt re-reads the preview (no local bundle copy to
       patch);
     - zero emoji.
   ============================================================================ */

import { readFileSync } from "node:fs";
import { renderToStaticMarkup } from "react-dom/server";
import { DiagnosticsExportBlock, DiagnosticsPreviewBlock } from "@/components/diagnostics-settings";
import {
  parseDiagnosticsExportReply,
  parseDiagnosticsPreviewReply,
  type DiagnosticsPreviewState,
} from "@/lib/diagnostics-view";

const failures: string[] = [];
function check(label: string, ok: boolean, detail = ""): void {
  if (!ok) failures.push(`${label}${detail ? ` [${detail}]` : ""}`);
}

// The emoji band the repo scanner walks (ban #8, tools/d22scan's own set).
const EMOJI_RE =
  /[\u{1F000}-\u{1FAFF}\u{2200}-\u22FF\u{2600}-\u{27BF}\u{2B00}-\u{2BFF}\u{1F1E6}-\u{1F1FF}]|\uFE0F/u;

/** Same escaping React applies to a text node - so a match here is a match on
    screen, not a match on an unescaped substring that never renders. */
function esc(text: string): string {
  return text
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#x27;");
}

/* ---------------------------------------------------------------------------
   1. The preview reply parses: prose unparsed, the engine's own fields
      verbatim
   --------------------------------------------------------------------------- */

// HandleDiagnosticsPreview's prose shape, then the marker line whose JSON is
// the engine's own BundlePreview (the field set internal/observe's json tags
// carry - categories with detail verbatim, items with the redacted head).
const PREVIEW_PROSE =
  "诊断包预览（条目 4 个，合计 2048 字节）。预览与导出走同一份文件清单来源：同一个 BundleOptions 组装、同一对 Preview/Build 函数（票 45 引擎），所以这份清单就是导出 zip 的条目清单。排除类（转写/watchdog/音频）按引擎自报的默认排除原因逐类列出。";
const PREVIEW_PAYLOAD = `{"created_at":"2026-11-18T07:40:00Z","out_path":"(preview only: nothing is written)","categories":[` +
  `{"category":"logs","default_included":true,"collected":true,"detail":"1 file(s) would be copied, 0 skipped (oversized); rolling JSONL logs; redacted at rest and re-passed through the redactor","items":[{"name":"logs/wisp-20261118-001.jsonl","bytes":900,"head":"{\\"level\\":\\"INFO\\",\\"msg\\":\\"provider call\\"}","included":true}]},` +
  `{"category":"config","default_included":true,"collected":true,"detail":"config.toml re-passed through the non-disableable redactor; a surviving key shape aborts the build","items":[{"name":"config.redacted.toml","bytes":420,"head":"[privacy]\\ndiagnostics_opt_in = true","included":true}]},` +
  `{"category":"transcripts","default_included":false,"collected":false,"detail":"ASR transcripts - excluded by default; each item must be named after it was previewed (none opted in: nothing under the transcript dir is read)"},` +
  `{"category":"audio","default_included":false,"collected":false,"detail":"never collected: ASR audio is not read by the bundle builder at all (D16 privacy boundary)"}` +
  `],"entries":["logs/wisp-20261118-001.jsonl","config.redacted.toml","version.json","manifest.json"],"total_bytes":2048}`;
const PREVIEW_REPLY = `${PREVIEW_PROSE}\nwisp-diagnostics-preview:${PREVIEW_PAYLOAD}`;

const previewParsed = parseDiagnosticsPreviewReply(PREVIEW_REPLY);
check("the preview reply must parse as the parsed state", previewParsed.kind === "parsed");
if (previewParsed.kind === "parsed") {
  check("the preview prose must arrive verbatim, unparsed", previewParsed.prose === PREVIEW_PROSE, previewParsed.prose);
  check("the engine's entries list must ride verbatim", JSON.stringify(previewParsed.envelope.entries) === JSON.stringify([
    "logs/wisp-20261118-001.jsonl", "config.redacted.toml", "version.json", "manifest.json",
  ]), JSON.stringify(previewParsed.envelope.entries));
  check("the engine's total_bytes must ride verbatim", previewParsed.envelope.totalBytes === 2048);
  const logs = previewParsed.envelope.categories.find((c) => c.category === "logs");
  check("the logs category must ride with the engine's own detail", logs !== undefined && logs.detail.includes("redactor"), logs?.detail ?? "(missing)");
  const audio = previewParsed.envelope.categories.find((c) => c.category === "audio");
  check(
    "the audio category must ride with the engine's own never-collected statement verbatim",
    audio !== undefined && audio.detail === "never collected: ASR audio is not read by the bundle builder at all (D16 privacy boundary)",
    audio?.detail ?? "(missing)",
  );
  check("no damage in a clean payload means skipped=0", previewParsed.skipped === 0);
}

/* ---------------------------------------------------------------------------
   2. The preview renders field by field: entries verbatim, head under the
      脱敏 label, the audio statement verbatim
   --------------------------------------------------------------------------- */

const previewState: DiagnosticsPreviewState =
  previewParsed.kind === "parsed"
    ? { phase: "ready", prose: previewParsed.prose, envelope: previewParsed.envelope, skipped: previewParsed.skipped, raw: PREVIEW_REPLY }
    : ({ phase: "loading" } as DiagnosticsPreviewState);
const previewHtml = renderToStaticMarkup(<DiagnosticsPreviewBlock state={previewState} />);
check("the entry count must state the payload's own entries length", previewHtml.includes(esc("zip 条目清单（4 条")));
check("the log entry must render verbatim", previewHtml.includes(esc("logs/wisp-20261118-001.jsonl")));
check("the log item must render its real bytes verbatim", previewHtml.includes(esc("900 字节")));
check("the log item's redacted head must render under the 脱敏 label", previewHtml.includes(esc("脱敏")) && previewHtml.includes(esc('"level":"INFO"')));
check("the audio category's never-collected statement must render verbatim", previewHtml.includes(esc("never collected: ASR audio is not read by the bundle builder at all (D16 privacy boundary)")));
check("the transcripts default-excluded statement must render verbatim", previewHtml.includes(esc("excluded by default")));
check("the same-source claim must be on screen (the preview IS the export list)", previewHtml.includes(esc("同一份文件清单来源")));
check("the preview block must hold zero emoji", !EMOJI_RE.test(previewHtml));

/* ---------------------------------------------------------------------------
   3. THE ZERO-ENTRY PREVIEW: said out loud with the payload's own statements,
      never a bare empty list (tooth ③'s green half)
   --------------------------------------------------------------------------- */

const ZERO_PROSE = "诊断包预览（条目 0 个，合计 0 字节）。预览与导出走同一份文件清单来源。";
const AUDIO_DETAIL = "never collected: ASR audio is not read by the bundle builder at all (D16 privacy boundary)";
const TRANS_DETAIL = "ASR transcripts - excluded by default; each item must be named after it was previewed (none opted in: nothing under the transcript dir is read)";
const ZERO_PAYLOAD = `{"created_at":"2026-11-18T07:41:00Z","out_path":"(preview only: nothing is written)","categories":[` +
  `{"category":"logs","default_included":true,"collected":false,"detail":"not collected: log dir not provided"},` +
  `{"category":"transcripts","default_included":false,"collected":false,"detail":"${TRANS_DETAIL}"},` +
  `{"category":"audio","default_included":false,"collected":false,"detail":"${AUDIO_DETAIL}"}` +
  `],"entries":[],"total_bytes":0}`;
const zeroParsed = parseDiagnosticsPreviewReply(`${ZERO_PROSE}\nwisp-diagnostics-preview:${ZERO_PAYLOAD}`);
check("a zero-entry preview must parse as the parsed state", zeroParsed.kind === "parsed");
if (zeroParsed.kind === "parsed") {
  check("the zero-entry envelope must carry the payload's empty entries verbatim", zeroParsed.envelope.entries.length === 0);
  const zeroHtml = renderToStaticMarkup(
    <DiagnosticsPreviewBlock state={{ phase: "ready", prose: zeroParsed.prose, envelope: zeroParsed.envelope, skipped: zeroParsed.skipped, raw: "" }} />,
  );
  check("the zero-entry preview must say 条目 0 个 out loud", zeroHtml.includes(esc("条目 0 个")));
  check("the zero-entry preview must refuse to read as a normal empty list", zeroHtml.includes(esc("不是正常的空读数")));
  check("the zero-entry preview must carry the audio statement verbatim", zeroHtml.includes(esc(AUDIO_DETAIL)));
  check("the zero-entry preview must carry the transcripts statement verbatim", zeroHtml.includes(esc(TRANS_DETAIL)));
  check("the zero-entry preview must draw NO entry-list chrome (no zip 条目清单 table)", !zeroHtml.includes(esc("zip 条目清单（")));
  check("the zero-entry block must hold zero emoji", !EMOJI_RE.test(zeroHtml));
}

/* ---------------------------------------------------------------------------
   4. Refusals ride verbatim; damaged payloads are named, never drawn
   --------------------------------------------------------------------------- */

// RefusedEnvelopeForUser's wrapper around the switch-off refusal: the cause
// NAMES the total switch. Fed to the page, the words must come back out.
const SWITCH_OFF_REFUSAL =
  `面板请求被拒绝 [diagnostics.preview pc-5-abc]：panel: diagnostics.preview 被拒绝（requestId="pc-5-abc"）：` +
  `总闸 [privacy] diagnostics_opt_in = false（来源 C:/data/config.toml）。面板诊断面与 wisp diagnostics CLI 是同一把总闸：关着时不做预览、不读任何日志；要导出诊断包，请把 config.toml 里 [privacy] 段的 diagnostics_opt_in 改为 true 后重试。`;
const refusalHtml = renderToStaticMarkup(<DiagnosticsPreviewBlock state={{ phase: "unreadable", raw: SWITCH_OFF_REFUSAL }} />);
check("the switch-off refusal must ride verbatim", refusalHtml.includes(esc(SWITCH_OFF_REFUSAL)));
check("the refusal must name the total switch (diagnostics_opt_in)", refusalHtml.includes(esc("diagnostics_opt_in")));
check("the refusal must render in the danger tint", refusalHtml.includes("bg-red-tint"));
check("the refusal block must hold zero emoji", !EMOJI_RE.test(refusalHtml));

// A reply WITHOUT its payload line: the reading stays 未判 and no fabricated
// entry may be drawn (喂不到 ⇒ 显式未判, 禁止假清单).
const noPayloadHtml = renderToStaticMarkup(
  <DiagnosticsPreviewBlock state={{ phase: "unreadable", raw: "诊断包预览被拒绝，载荷行不在。" }} />,
);
check("a marker-less reply must stay an explicit named state", noPayloadHtml.includes(esc("诊断包预览被拒绝")));
check("the no-payload preview block must hold zero emoji", !EMOJI_RE.test(noPayloadHtml));

// A payload whose entries are damaged: payload-error named, zeros never
// folded in (the cost-view rule).
const damagedParsed = parseDiagnosticsPreviewReply(
  `x\nwisp-diagnostics-preview:{"created_at":"2026-11-18T07:42:00Z","out_path":"p","categories":[],"entries":["ok.jsonl",""],"total_bytes":10}`,
);
check("a damaged entries payload must refuse to parse as ready", damagedParsed.kind === "payload-error");
if (damagedParsed.kind === "payload-error") {
  const damagedHtml = renderToStaticMarkup(<DiagnosticsPreviewBlock state={{ phase: "payload-error", message: damagedParsed.message, raw: damagedParsed.raw }} />);
  check("the damaged-entries block must name the 未判 rule", damagedHtml.includes(esc("未判")) && damagedHtml.includes(esc("不画 0 条")));
  check("the damaged-entries block must hold zero emoji", !EMOJI_RE.test(damagedHtml));
}

// The no-host state (plain browser / hand-assembled harness): explicit, and
// nothing that reads like a result.
const noHostHtml = renderToStaticMarkup(<DiagnosticsPreviewBlock state={{ phase: "no-host" }} />);
check("the no-host block must say the reading is 未判", noHostHtml.includes(esc("未判")));
check("the no-host block must hold zero emoji", !EMOJI_RE.test(noHostHtml));

/* ---------------------------------------------------------------------------
   5. The export receipt rides verbatim and names the path it wrote
   --------------------------------------------------------------------------- */

const EXPORT_PROSE =
  "已导出诊断包：C:/data/wisp-diagnostics-20261118-074300Z.zip（条目 4 个；内容以预览与包内 manifest.json 为准——manifest 逐条记录了每个条目与每个未收集项的原因）。音频、密钥、完整转写按引擎默认排除；转写与 watchdog 快照默认排除、逐项同意。";
const EXPORT_REPLY = `${EXPORT_PROSE}\nwisp-diagnostics-export:{"created_at":"2026-11-18T07:43:00Z","out_path":"C:/data/wisp-diagnostics-20261118-074300Z.zip","categories":[],"entries":["logs/wisp-20261118-001.jsonl","config.redacted.toml","version.json","manifest.json"],"total_bytes":2048}`;
const exportParsed = parseDiagnosticsExportReply(EXPORT_REPLY);
check("the export reply must parse as the parsed state", exportParsed.kind === "parsed");
if (exportParsed.kind === "parsed") {
  check("the export prose must arrive verbatim", exportParsed.prose === EXPORT_PROSE, exportParsed.prose);
  const exportHtml = renderToStaticMarkup(
    <DiagnosticsExportBlock state={{ phase: "ready", prose: exportParsed.prose, envelope: exportParsed.envelope, skipped: exportParsed.skipped }} />,
  );
  check("the export receipt must ride verbatim", exportHtml.includes(esc(EXPORT_PROSE)));
  check("the written path must ride verbatim (payload out_path)", exportHtml.includes(esc("C:/data/wisp-diagnostics-20261118-074300Z.zip")));
  check("the export block must state the exclusions point back at the payload's own statements", exportHtml.includes(esc("按引擎默认排除")));
  check("the export block must hold zero emoji", !EMOJI_RE.test(exportHtml));
}

const exportRefusal =
  `面板请求被拒绝 [diagnostics.export pc-7-def]：panel: diagnostics.export 被拒绝（requestId="pc-7-def"）：总闸 [privacy] diagnostics_opt_in = false（来源 C:/data/config.toml）。面板诊断面与 wisp diagnostics CLI 是同一把总闸。`;
const exportRefusalHtml = renderToStaticMarkup(<DiagnosticsExportBlock state={{ phase: "refused", raw: exportRefusal }} />);
check("the export refusal must ride verbatim", exportRefusalHtml.includes(esc(exportRefusal)));
check("the export refusal must render in the danger tint", exportRefusalHtml.includes("bg-red-tint"));
check("the export refusal block must hold zero emoji", !EMOJI_RE.test(exportRefusalHtml));

const exportNoPayloadHtml = renderToStaticMarkup(
  <DiagnosticsExportBlock state={{ phase: "no-payload", raw: "导出回执没有载荷行。" }} />,
);
check("the export no-payload block must refuse to fake a success (写盘结果未判)", exportNoPayloadHtml.includes(esc("写盘结果未判")) && exportNoPayloadHtml.includes(esc("不画假成功")));
check("the export no-payload block must hold zero emoji", !EMOJI_RE.test(exportNoPayloadHtml));

const exportNoHostHtml = renderToStaticMarkup(<DiagnosticsExportBlock state={{ phase: "no-host" }} />);
check("the export no-host block must refuse to fake a success", exportNoHostHtml.includes(esc("不画假成功")));
check("the export no-host block must hold zero emoji", !EMOJI_RE.test(exportNoHostHtml));

/* ---------------------------------------------------------------------------
   6. Source-level wiring: the literals, the engine's own words, the re-read
   --------------------------------------------------------------------------- */

const panelSrc = readFileSync("src/lib/panel.ts", "utf8");
check("panel.ts must spell diagnostics.preview as a literal (the ledger's closed gap)", panelSrc.includes('"diagnostics.preview"'));
check("panel.ts must spell diagnostics.export as a literal", panelSrc.includes('"diagnostics.export"'));
// The exclusion statements are the ENGINE's own words (they ride the
// payload's category detail). The page must render state-enveloped detail
// verbatim and never spell the Go sentences itself - a page-side exclusion
// sentence would be a fabricated reading (mutation ③'s target).
const viewSrc = readFileSync("src/lib/diagnostics-view.ts", "utf8");
const cardSrc = readFileSync("src/components/diagnostics-settings.tsx", "utf8");
check("the page must render the payload's category detail, not a page-side exclusion sentence", cardSrc.includes("category.detail") && !cardSrc.includes("never collected"));
check("the lib's markers must be spelled and match the Go constants", viewSrc.includes('"wisp-diagnostics-preview:"') && viewSrc.includes('"wisp-diagnostics-export:"'));
check("the page's requests must come through the ONE door file (lib/panel.ts)", viewSrc.includes('from "@/lib/panel"'));
check("the preview must be re-read after an export receipt (no local bundle copy to patch)", cardSrc.includes("useDiagnosticsExport(preview.reload)"));
check("the zero-entry branch must draw no entry-list chrome", cardSrc.includes("entries.length === 0"));
check("the settings page's diagnostics cell must mount the real card", readFileSync("src/components/harness/settings-page.tsx", "utf8").includes("DiagnosticsSettings"));

// The whole card in its SSR state: loading is explicit, no emoji anywhere.
const loadingHtml = renderToStaticMarkup(<DiagnosticsPreviewBlock state={{ phase: "loading" }} />);
check("the loading state must say what is in flight", loadingHtml.includes(esc("正在读取诊断包预览")));
check("the loading card must hold zero emoji", !EMOJI_RE.test(loadingHtml));
check("the card source must hold zero emoji", !EMOJI_RE.test(cardSrc) && !EMOJI_RE.test(viewSrc) && !EMOJI_RE.test(panelSrc));

/* ---------------------------------------------------------------------------
   The verdict.
   --------------------------------------------------------------------------- */

if (failures.length > 0) {
  for (const f of failures) console.error("render-diagnostics40: FAIL " + f);
  console.error(`render-diagnostics40: FAIL, ${failures.length} check(s) failed`);
  process.exit(1);
}
console.error(
  "render-diagnostics40: OK, ticket 40 AC#5's diagnostics page consumes the real doors " +
    "(the engine's own BundlePreview rides verbatim - entries list, item name/bytes, the redacted head under the 脱敏 label, " +
    "the excluded categories' own statements verbatim including audio's never-collected line, " +
    "the zero-entry preview says 条目 0 个 out loud with those statements and never reads as a normal empty list, " +
    "the export receipt names the written path verbatim and the exclusions point back at the payload, " +
    "refusals ride verbatim including the sentence that names diagnostics_opt_in, " +
    "a marker-less reply stays 未判 and never draws a fabricated entry or path, damaged entries are named never folded into zeros, " +
    "the two literals live in panel.ts, the page renders the payload's own words never the engine's sentences, the card re-reads after an export receipt, " +
    "zero emoji)",
);
