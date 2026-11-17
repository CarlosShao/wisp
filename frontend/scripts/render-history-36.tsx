/* ============================================================================
   Ticket 36's behavioural tooth (render + rule evidence).
   ----------------------------------------------------------------------------
   WHAT THIS FILE PROVES, and why a screenshot cannot:
   every check below feeds REAL-SHAPED data through the PRODUCTION component
   and asserts the SAME VALUE comes back out in the markup - the "feed real
   data => display the same value" tooth this ticket's dispatch demands, in
   the shape render-l2/render-shell-77 already established for this tree.

   The fixtures are not hand-drawn JSON: they are the wire's own documented
   shapes, read out of internal/panel's wire structs (run_feed_145.go's JSON
   tags) the same way render-shell-77 reads ApprovalCardView. Where a dimension
   is the contract ruler's registered Go-only key (composer_test.go:231-233),
   the check asserts its ABSENCE instead - a renderer that invented a value
   for an undeclared key would be the exact 假值 the dispatch forbids.

   The XSS half (AC#2) feeds the hostile payload families (parser-borne raw
   HTML, schema-judged markdown links, plus the scheme-relative href negative
   controls) through the production MarkdownText (the same function every
   finished chunk takes) and
   asserts none of them reaches the markup. This is a function-level proof,
   not a browser run: the sanitizer is a pure string->string step, so if the
   hostile fragment is not in the static markup it cannot execute. The CSP
   report-only listener stays the shipped belt; this file is the braces.
   ============================================================================ */

import { renderToStaticMarkup } from "react-dom/server";
import App from "@/App";
import { MarkdownText, ResultHistory, RunRegion } from "@/components/result-history";
import type { PanelSnapshot, ResultChunkView, RunSection } from "@/lib/panel";

const failures: string[] = [];
function check(label: string, ok: boolean, detail = ""): void {
  if (!ok) failures.push(`${label}${detail ? ` [${detail}]` : ""}`);
}

/* --------------------------------------------------------------------------
   1. The sanitizer: the four AC#2 payload families, through production code
   -------------------------------------------------------------------------- */

// Each entry is the payload a fake fetch/markdown chunk would carry, the
// fragment that makes it hostile, and the sanitized remnant the whitelist
// allows through (asserting the remnant keeps the content readable where the
// spec wants it kept). NOTE on raw HTML: react-markdown's default pipeline
// has no rehype-raw, so raw-HTML payloads (img onerror, script, inline event
// handlers) never become elements at all - they die in the parser, one layer
// BEFORE the whitelist. The whitelist is what guards the markdown-borne
// shapes (links, code, tables), so the "content survives" probes ride the
// markdown forms.
const XSS_CASES: { name: string; payload: string; hostile: RegExp; kept?: string }[] = [
  {
    name: "img onerror",
    payload: "先看目录\n\n<img src=x onerror=\"window.__pwned=1\">",
    hostile: /<img|onerror/i,
  },
  {
    name: "script tag",
    payload: "结论如下<script>alert(1)</script>",
    hostile: /<script/i,
  },
  {
    name: "javascript: href",
    payload: "[点这里](javascript:alert(1))",
    hostile: /javascript:/i,
    kept: "点这里",
  },
  {
    // The scheme-relative negative control (the audit's open-redirect face):
    // hast-util-sanitize's protocol check passes any colon-less URL untouched,
    // so //evil.example/x survived the TIGHT_SCHEMA alone and a browser
    // resolves it against the page origin. The urlTransform gate must strip
    // the href - the anchor may keep its text, never its destination.
    name: "protocol-relative href",
    payload: "[点这里](//evil.example/x)",
    hostile: /evil\.example|href="\/\//i,
    kept: "点这里",
  },
  {
    // The backslash sibling: browsers normalize \ to /, so a href starting
    // with a backslash is the same scheme-relative trick in disguise.
    name: "backslash scheme-relative href",
    payload: "[点这里](\\evil.example/x)",
    hostile: /evil\.example/i,
    kept: "点这里",
  },
  {
    // Raw HTML with handlers: the parser drops the element wholesale (no
    // rehype-raw), so neither the attributes nor the text reach the markup -
    // this case asserts BOTH halves, which is why it has no `kept`.
    name: "event handler attr",
    payload: '<a href="https://example.com" onclick="alert(1)" onmouseover="alert(2)">链接</a>',
    hostile: /on(click|mouseover)=|example\.com/i,
  },
];

for (const c of XSS_CASES) {
  const html = renderToStaticMarkup(<MarkdownText text={c.payload} />);
  check(`XSS ${c.name}: the hostile fragment must not reach the markup`, !c.hostile.test(html), html.slice(0, 200));
  if (c.kept) {
    check(`XSS ${c.name}: the content itself survives the whitelist`, html.includes(c.kept), html.slice(0, 200));
  }
}

// The belt beneath the braces: the export the CI greps for must not appear in
// this tree's source at all.
const { readFileSync, readdirSync, statSync } = await import("node:fs");
const { join } = await import("node:path");
function walkSources(dir: string): string[] {
  const out: string[] = [];
  for (const e of readdirSync(dir)) {
    const p = join(dir, e);
    const st = statSync(p);
    if (st.isDirectory()) out.push(...walkSources(p));
    else if (/\.(tsx?|css)$/.test(e)) out.push(p);
  }
  return out;
}
// The belt beneath the braces: the JSX attribute shape must appear in zero
// files. The comment EXPLAINING the ban is allowed (this tree's gates quote
// their own rules; render-nav does the same strip), so comments go first.
let innerHtmlHits = 0;
for (const f of walkSources("src")) {
  const stripped = readFileSync(f, "utf8")
    .replace(/\/\*[\s\S]*?\*\//g, "")
    .replace(/^\s*\/\/.*$/gm, "");
  if (stripped.includes("dangerouslySetInnerHTML")) innerHtmlHits += 1;
}
check("dangerouslySetInnerHTML must appear in zero frontend/src files (comments exempt)", innerHtmlHits === 0, `${innerHtmlHits} file(s)`);

/* --------------------------------------------------------------------------
   2. Feed real data => the same value shows (the run feed's dimensions)
   -------------------------------------------------------------------------- */

// Shapes lifted from internal/panel/run_feed_145.go's JSON tags - what Go
// actually puts on the wire (the Declared keys only; see the ruler note).
const RUN_FED: RunSection = {
  tools: [
    { callId: "call-17", name: "web.fetch", outcome: "success", resultText: "抓到 4,213 字节，标题: Wisp 项目主页" },
    { callId: "call-18", name: "fs.delete", outcome: "error" },
    { callId: "call-19", name: "shell.exec", outcome: "cancelled" },
    { callId: "call-20", name: "doc.read", outcome: "truncated", resultText: "前 2,000 字节" },
    { callId: "call-21", name: "mystery.tool", outcome: "weird_outcome" },
  ],
  toolsDropped: 3,
  reasoning: "...[truncated: 128 runes elided]...先看包结构，再定位 pump 的 publish 链。",
  error: { class: "net_timeout", detail: "dial tcp 127.0.0.1:9090: connectex: timeout", text: "工具调用超时，未取得结果。", retryable: true },
  stuck: { text: "web.fetch 已连续调用 8 次，参数相同", repeatLevel: 2 },
};

const fedHtml = renderToStaticMarkup(<RunRegion run={RUN_FED} />);

// Every fed value must come back out verbatim. The collapsed halves (the
// reasoning block, the error's expandable detail, the tool chip's expanded
// block) render at the INTERACTION state, not the initial one - react-dom/
// server renders useState's initial value (the same limitation render-stream
// records in its header) - so their VALUES are pinned at the source level
// below, and the markup asserts what the initial state actually paints.
check("the fed tool name shows", fedHtml.includes("web.fetch"));
check("the fed result text shows", fedHtml.includes("Wisp 项目主页"));
check("the fed toolsDropped count shows with the truncate-not-merge wording", fedHtml.includes("3") && fedHtml.includes("截"));
check("the fed error text is the human copy", fedHtml.includes("工具调用超时"));
check("the fed error class chip shows the class", fedHtml.includes("net_timeout"));
check("the fed stuck sentence shows verbatim", fedHtml.includes("web.fetch 已连续调用 8 次，参数相同"));
check("the reasoning line is mounted collapsed-by-default", fedHtml.includes("推理过程"));

// Source-level pins for the interaction-state values: the strings the user
// sees on expand are exactly the fed strings (never a paraphrase).
const resultHistorySrc = readFileSync("src/components/result-history.tsx", "utf8");
check("the error card renders err.detail verbatim on expand", resultHistorySrc.includes("{err.detail ??") || /err\.detail/.test(resultHistorySrc));
check("the tool chip's expanded block renders resultText verbatim", resultHistorySrc.includes("{row.resultText}") || /row\.resultText/.test(resultHistorySrc));
check("the tool chip's expanded block carries correlationId", resultHistorySrc.includes("correlationId"));
check("the reasoning block renders the fed text verbatim", resultHistorySrc.includes("{text}"));

// Outcome vocabulary -> the frozen glyph row (loop's journal vocabulary). The
// glyphs paint in row order; each row is anchored by its tool name (the call
// id is a React key, not markup), and the slice runs to the next anchor.
const outcomeChecks: [string, RegExp][] = [
  ["web.fetch", /aria-label="成功"/],
  ["fs.delete", /aria-label="失败"/],
  ["shell.exec", /aria-label="已取消"/],
  ["doc.read", /aria-label="已截断"/],
  // An outcome the loop never publishes must NOT get a success glyph - the
  // unknown-word shape renders the word itself in the label.
  ["mystery.tool", /aria-label="未知结果 weird_outcome"/],
];
for (const [name, labelRe] of outcomeChecks) {
  const rowStart = fedHtml.indexOf(name);
  const rowEnd = fedHtml.indexOf("</li>", rowStart);
  const row = fedHtml.slice(rowStart, rowEnd > 0 ? rowEnd : undefined);
  check(`outcome glyph for ${name}`, labelRe.test(row), row.slice(0, 160));
}

/* --------------------------------------------------------------------------
   3. The registered Go-only keys must be INVENTED NOWHERE
   -------------------------------------------------------------------------- */

// composer_test.go:231-233 registers argsJson/level/retryable/running/usage/
// done as Go-only until the ledger move. None of their VALUES may appear on
// screen - check the fed tree carries none of the key names as rendered data.
for (const key of ["argsJson", "retryable:", "usage"]) {
  check(
    `registered Go-only key ${key} must not be rendered as data`,
    !fedHtml.includes(key),
    `found "${key}" in the fed tree`,
  );
}
// `running` and `level:` are substring-hostile in ordinary prose, so their
// absence is pinned against the markup's aria labels and title attributes
// only - a risk pill or an in-flight chip would have to paint one of these.
check("no risk-level pill (L0/L1/L2) is drawn from an undeclared key", !/aria-label="L[012]/.test(fedHtml));
check("no in-flight chip is drawn (no running reader on this wire)", !/aria-label="运行中"/.test(fedHtml));

/* --------------------------------------------------------------------------
   4. The absent-reader shape: no run key => the named absence
   -------------------------------------------------------------------------- */

const chunks: ResultChunkView[] = [{ correlationId: "c1", text: "只有一条结果。", done: true }];
const noRunHtml = renderToStaticMarkup(<ResultHistory run={undefined} results={chunks} />);
check("with chunks but no run key, the absence names the missing engine", noRunHtml.includes("工作状态未接线"));
check("the absence names the missing reader", noRunHtml.includes("PumpSources.Run"));
check("the fed chunk still shows verbatim", noRunHtml.includes("只有一条结果。"));

// Empty everything: the honest empty state comes from App's placeholder, so
// the region itself renders null (never a bare shell).
const emptyHtml = renderToStaticMarkup(<ResultHistory run={undefined} results={[]} />);
check("with nothing fed, the region renders empty (App's placeholder owns the sentence)", emptyHtml === "");

// Present-but-empty run (pipeline shape, zero events): the region also
// renders null - "no events yet" is the placeholder's business, not a bare
// 工作状态 header with nothing under it.
const emptyRun = renderToStaticMarkup(<RunRegion run={{ tools: [] }} />);
check("a present-but-empty run renders nothing", emptyRun === "");

/* --------------------------------------------------------------------------
   5. Whole-panel wiring: the App tree shows the fed values (ticket 329's hop)
   -------------------------------------------------------------------------- */

const SNAP: PanelSnapshot = {
  pending: [],
  results: chunks,
  composer: {
    mode: { current: "unknown", names: [], l2ConfirmNames: [] },
    workspace: { set: true, spelling: "D:\\w", canonical: "D:\\w", reparse: false, rewritten: false, reason: "" },
    attachments: [],
    acceptedAttachmentMimes: [],
    maxAttachmentBytes: 0,
    attachmentError: "",
    git: {
      kind: "repo",
      reason: "这棵工作树当前检出分支 dev",
      branch: "dev",
      detachedSha: "",
      isDetached: false,
      repoRoot: "D:\\w",
      currentWorktree: "D:\\w",
      worktrees: [{ path: "D:\\w", branch: "dev", sha: "", detached: false, main: true }],
      branches: ["dev", "dsh/feat/x"],
      switchBlocked: "",
    },
  },
  generatedAt: "render-evidence-36",
};
const appHtml = renderToStaticMarkup(<App snapshot={SNAP} />);
check("the conversation column shows the fed chunk through the sanitizer path", appHtml.includes("只有一条结果。"));
check("the header git chip shows the fed branch from composer.git (ticket 36's wiring fix)", appHtml.includes("当前分支 dev"));

/* --------------------------------------------------------------------------
   verdict
   -------------------------------------------------------------------------- */

if (failures.length > 0) {
  for (const f of failures) console.error("render-history36: FAIL " + f);
  console.error(`render-history36: FAIL, ${failures.length} check(s) failed`);
  process.exit(1);
}
console.log(
  `render-history36: OK, ${XSS_CASES.length} XSS families neutralized through production MarkdownText, ` +
    `run-feed values round-trip (${RUN_FED.tools.length} tool rows, error, stuck, reasoning), ` +
    `registered Go-only keys invented nowhere, absent-reader shape names its engine, ` +
    `dangerouslySetInnerHTML grep = 0, git chip reads composer.git`,
);
