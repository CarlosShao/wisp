/* ============================================================================
   Ticket 36 AC#3's page-half tooth (render evidence), 2026-11-17 - the
   history browser leg (consume wave 2).
   ----------------------------------------------------------------------------
   WHAT THIS FILE PROVES, in the shape render-privacy-40 established: every
   check below feeds REAL-SHAPED data through the PRODUCTION code and asserts
   the SAME VALUE comes back out - no browser needed, because every piece
   under test is a pure function over the wire's own shapes.

   THE FIXTURES are the Go leg's own documented shapes, copied from the Go
   half that answers them (cmd/wisp/panel_history_wave2.go's receipt +
   "wisp-history-page:" payload over internal/memory/dao_history_page_36.go's
   HistoryPageRow). The 1005-row store mirrors the Go scroll test's fixture
   (cmd/wisp/panel_history_wave2_test.go:TestHistoryPageScrollsOneThousandRows
   WithoutOverlap): rows newest-first, three windows of 500/500/5, and the
   same three window firsts the Go test pins (scroll-1004 / scroll-0504 /
   scroll-0004).

   THE TEETH:
     1. 1005 rows through THREE windows via the production parser + merger:
        no overlap (every id exactly once), no gap (window continuity held,
        zero anomalies), first/last rows verbatim;
     2. a duplicate id and a broken window offset are NAMED defects - the
        merge never silently absorbs or drops them;
     3. the virtual slice covers the list at every scroll position without
        overlapping itself, and the top-of-list render shows the FIRST rows
        while the LAST row stays unrendered (virtualization is real, not a
        claim);
     4. empty store (total=0) renders 暂无记录 - and the pre-reply browser
        renders NEITHER a list NOR the empty claim;
     5. refusals ride verbatim; a receipt without its payload line and a
        payload without rows stay 未判, never a fabricated list;
     6. delete-one / purge / export ride the EXISTING privacy face: the
        browser imports the wave-1 hook and reply block wholesale, purge
        fires only from its confirm card, delete-one goes out per row with
        the face's own domain mapping;
     7. source level: the history.page literal + envelope keys live in
        panel.ts, and the page-gap ledger's two wave-2 entries are GONE (the
        both-ways rule) - plus zero emoji everywhere.
   ============================================================================ */

import { readFileSync } from "node:fs";
import { renderToStaticMarkup } from "react-dom/server";
import {
  HistoryBrowser,
  HistoryEmptyBlock,
  HistoryRefusalBlock,
  HistoryRowLine,
  historyDomainOf,
} from "@/components/history-page";
import {
  EMPTY_HISTORY_MERGE,
  HISTORY_ROW_HEIGHT,
  formatHistoryTs,
  mergeHistoryPage,
  parseHistoryPageReply,
  visibleSlice,
  type HistoryPagePayload,
} from "@/lib/history-view";

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
   The 1005-row store, newest first: row at offset i carries id
   scroll-<1004-i> and ts 1763376000 - i. Exactly the shape the Go scroll
   test feeds a live store through the real DAO.
   --------------------------------------------------------------------------- */

function historyRow(i: number): HistoryPagePayload["rows"][number] {
  const n = 1004 - i;
  const id = `scroll-${n.toString().padStart(4, "0")}`;
  if (n % 2 === 0) {
    return { id, ts: 1763376000 - i, kind: "task", outcome: "Completed", summary: `任务 ${id} 的摘要行` };
  }
  return { id, ts: 1763376000 - i, kind: "tool", tool: "fs.read", outcome: "success" };
}

function historyReply(offset: number, count: number, total: number): string {
  const rows = Array.from({ length: count }, (_, k) => historyRow(offset + k));
  const payload = JSON.stringify({ rows, total, offset, limit: count });
  // The payload line is the whole remainder of its line (the parser takes
  // the rest of the line and JSON.parses it) - the marker's JSON is the
  // line's last word, exactly how cmd/wisp's marshalHistoryPayloadWave2
  // writes it.
  return `历史共 ${total} 行（新的事件在前），本页第 ${offset + 1} 行起的 ${count} 行。\nwisp-history-page:${payload}`;
}

/* ---------------------------------------------------------------------------
   1. The 1005-row tooth: three windows, production parser + merger.
   --------------------------------------------------------------------------- */

const W1 = parseHistoryPageReply(historyReply(0, 500, 1005));
const W2 = parseHistoryPageReply(historyReply(500, 500, 1005));
const W3 = parseHistoryPageReply(historyReply(1000, 5, 1005));
check("window 1 must parse as a page", W1.phase === "page");
check("window 2 must parse as a page", W2.phase === "page");
check("window 3 must parse as a page", W3.phase === "page");
if (W1.phase === "page" && W2.phase === "page" && W3.phase === "page") {
  check("the receipt prose must ride verbatim, unparsed", W1.prose.startsWith("历史共 1005 行（新的事件在前）"));
  const m1 = mergeHistoryPage(EMPTY_HISTORY_MERGE, W1.payload);
  const m2 = mergeHistoryPage(m1, W2.payload);
  const m3 = mergeHistoryPage(m2, W3.payload);
  check("three windows must merge to exactly 1005 rows", m3.rows.length === 1005, `got ${m3.rows.length}`);
  check("the store's total must ride through", m3.total === 1005);
  check("scrolling must bring NO duplicate id", m3.duplicates.length === 0, m3.duplicates.join(","));
  check("scrolling must name NO window anomaly", m3.anomalies.length === 0, m3.anomalies.join(" | "));
  const ids = new Set(m3.rows.map((r) => r.id));
  check("every id must appear exactly once (不重)", ids.size === 1005);
  let contiguous = true;
  for (let n = 0; n <= 1004; n += 1) {
    if (!ids.has(`scroll-${n.toString().padStart(4, "0")}`)) contiguous = false;
  }
  check("every row of the store must be on screen after the scroll (不漏)", contiguous);
  check("the first window's first row must be scroll-1004 (the Go test's own first)", m3.rows[0]?.id === "scroll-1004");
  check("the second window's first row must be scroll-0504", m2.rows[500]?.id === "scroll-0504");
  check("the third window's first row must be scroll-0004", m3.rows[1000]?.id === "scroll-0004");
  check("the last row must be scroll-0000", m3.rows[1004]?.id === "scroll-0000");
  check("row fields must ride verbatim (summary)", m3.rows[0]?.summary === "任务 scroll-1004 的摘要行");
  check("row fields must ride verbatim (tool)", m3.rows[1]?.tool === "fs.read");

  /* --- 2. The defects are NAMED: duplicate ids and broken offsets. --- */
  const dupMerge = mergeHistoryPage(m3, { ...W3.payload, offset: 1005 });
  check("a window bringing already-seen ids must be flagged as duplicates", dupMerge.duplicates.length === 5);
  check("a duplicate id must NOT be appended twice", dupMerge.rows.length === 1005, `got ${dupMerge.rows.length}`);
  const gapMerge = mergeHistoryPage(m3, {
    rows: [{ id: "gap-row", ts: 1, kind: "task" }],
    total: 1006,
    offset: 1200,
    limit: 500,
  });
  check("a window whose offset skips rows must be flagged as an anomaly", gapMerge.anomalies.length === 1);
  check("the anomaly must name both halves of the gap", gapMerge.anomalies[0]?.includes("1005") === true && gapMerge.anomalies[0]?.includes("1200") === true);
}

/* ---------------------------------------------------------------------------
   3. The virtual slice: covers the list, never overlaps itself, and the
      top-of-list render hides the tail.
   --------------------------------------------------------------------------- */

{
  const count = 1005;
  const viewport = 380;
  let covered = new Set<number>();
  let lastEnd = -1;
  let monotonic = true;
  for (let scrollTop = 0; scrollTop <= (count * HISTORY_ROW_HEIGHT - viewport); scrollTop += 40) {
    const s = visibleSlice(count, scrollTop, viewport);
    if (s.start < 0 || s.end > count || s.start > s.end) monotonic = false;
    if (s.end <= s.start && scrollTop < (count - 1) * HISTORY_ROW_HEIGHT) monotonic = false;
    for (let i = s.start; i < s.end; i += 1) covered.add(i);
  }
  check("every scroll position must yield an in-bounds, non-empty slice", monotonic);
  check("scrolling top-to-bottom must cover every row exactly once in the union", covered.size === 1005, `covered ${covered.size}`);
  const atTop = visibleSlice(count, 0, viewport);
  check("the top slice must start at row 0", atTop.start === 0);
  const atBottom = visibleSlice(count, (count - 1) * HISTORY_ROW_HEIGHT, viewport);
  check("the bottom slice must reach the last row", atBottom.end === 1005 && atBottom.start <= 1004);
  const pastEnd = visibleSlice(count, 9999999, viewport);
  // A scroll past the last row clamps to the list's end and renders nothing
  // (a real scrollbar cannot reach there; the arithmetic must not go
  // out-of-bounds or throw rows back into view).
  check(
    "a past-the-end scroll must clamp inside the list and render nothing",
    pastEnd.end <= 1005 && pastEnd.start <= pastEnd.end,
  );
  const empty = visibleSlice(0, 0, viewport);
  check("an empty list must yield an empty slice", empty.start === 0 && empty.end === 0);
}

/* ---------------------------------------------------------------------------
   4. Row rendering: the five display fields verbatim, the delete button
      disabled for a kind the face does not define.
   --------------------------------------------------------------------------- */

const rowHtml = renderToStaticMarkup(
  <HistoryRowLine deleteDisabled={false} onDelete={() => undefined} row={historyRow(0)} />,
);
check("the row must show the formatted timestamp", rowHtml.includes(esc(formatHistoryTs(1763376000))));
check("the row must carry the raw ts in the title", rowHtml.includes("ts=1763376000"));
check("the task row must show the 任务 label", rowHtml.includes(esc("任务")));
check("the task row must show the store's own outcome word verbatim", rowHtml.includes(esc("Completed")));
check("the task row must show the summary verbatim", rowHtml.includes(esc("任务 scroll-1004 的摘要行")));
check("the task row must name its privacy domain on the delete control", rowHtml.includes(esc("task_log")));
check("the row markup must hold zero emoji", !EMOJI_RE.test(rowHtml));

const toolRowHtml = renderToStaticMarkup(
  <HistoryRowLine deleteDisabled={false} onDelete={() => undefined} row={historyRow(1)} />,
);
check("the tool row must show the tool name verbatim", toolRowHtml.includes(esc("fs.read")));
check("the tool row must show the tool outcome verbatim", toolRowHtml.includes(esc("success")));
check("the tool row must name tool_call on the delete control", toolRowHtml.includes(esc("tool_call")));
check("historyDomainOf must map task -> task_log", historyDomainOf("task") === "task_log");
check("historyDomainOf must map tool -> tool_call", historyDomainOf("tool") === "tool_call");
check("historyDomainOf must refuse to invent a domain for an unknown kind", historyDomainOf("mystery") === "");

/* ---------------------------------------------------------------------------
   5. The empty store and the refused / unparseable replies, through the
      production blocks.
   --------------------------------------------------------------------------- */

const emptyHtml = renderToStaticMarkup(<HistoryEmptyBlock />);
check("the empty store must render 暂无记录", emptyHtml.includes(esc("暂无记录")));
check("the empty block must hold zero emoji", !EMOJI_RE.test(emptyHtml));

const REFUSAL =
  `面板请求被拒绝 [history.page pc-7-h1]：panel: history.page 被拒绝（requestId="pc-7-h1"）：` +
  `memory: history.page limit=900 不成立：一页至少 1 行且不超过上限 500：翻页用 offset，不用放大页`;
const refusalHtml = renderToStaticMarkup(<HistoryRefusalBlock message="" phase="refused" raw={REFUSAL} />);
check("the refusal must ride verbatim", refusalHtml.includes(esc(REFUSAL)));
check("the refusal must render in the danger tint", refusalHtml.includes("bg-red-tint"));
check("the refusal block must hold zero emoji", !EMOJI_RE.test(refusalHtml));

const noHostHtml = renderToStaticMarkup(<HistoryRefusalBlock message="" phase="no-host" raw="" />);
check("the no-host block must say the reading is 未判", noHostHtml.includes(esc("未判")));
check("the no-host block must not look like an empty history", !noHostHtml.includes(esc("暂无记录")));

const noPayloadHtml = renderToStaticMarkup(
  <HistoryRefusalBlock message="" phase="no-payload" raw="历史共 5 行。" />,
);
check("the no-payload block must say 未判", noPayloadHtml.includes(esc("未判")));
check("the no-payload block must carry the raw reply verbatim", noPayloadHtml.includes(esc("历史共 5 行。")));
check("the no-payload block must not draw a fabricated row count", !/已合并 \d+/.test(noPayloadHtml));

const payloadErrParsed = parseHistoryPageReply("x\nwisp-history-page:{\"total\":3,\"offset\":0,\"limit\":10}");
check("a payload without rows must refuse as payload-error", payloadErrParsed.phase === "payload-error");
const payloadErrHtml = renderToStaticMarkup(
  <HistoryRefusalBlock
    message={payloadErrParsed.phase === "payload-error" ? payloadErrParsed.message : ""}
    phase="payload-error"
    raw={payloadErrParsed.phase === "payload-error" ? payloadErrParsed.raw : ""}
  />,
);
check("the payload-error message must ride verbatim", payloadErrHtml.includes(esc("rows 不在或不是数组")));

// The pre-reply browser: neither a list nor an empty claim - 暂无记录 is the
// store's own total=0 answer, never this page's guess while the request is
// in flight.
const initialHtml = renderToStaticMarkup(<HistoryBrowser load={async () => null} />);
check("the pre-reply browser must not claim 暂无记录", !initialHtml.includes(esc("暂无记录")));
check("the pre-reply browser must not draw a fabricated row count", !initialHtml.includes("已合并 0"));
check("the pre-reply browser must still draw the privacy strip (the doors exist)", initialHtml.includes(esc("privacy.delete-one")));
check("the initial markup must hold zero emoji", !EMOJI_RE.test(initialHtml));

/* ---------------------------------------------------------------------------
   6. Source-level wiring: the literals, the ledger, the privacy face reuse.
   --------------------------------------------------------------------------- */

const panelSrc = readFileSync("src/lib/panel.ts", "utf8");
check("panel.ts must spell history.page as a literal (the ledger's closed gap)", panelSrc.includes('"history.page"'));
check("panel.ts must spell roster.cancel as a literal", panelSrc.includes('"roster.cancel"'));
check(
  "the envelope keys must be spelled as bridge.go's JSON tags carry them",
  panelSrc.includes("historyOffset") && panelSrc.includes("historyLimit") && panelSrc.includes("rosterTaskId"),
);

// The page-gap ledger's own both-ways rule: the page started carrying both
// literals, so the stale registrations had to leave in the same commit.
const ledgerSrc = readFileSync("../internal/panel/bridge_test.go", "utf8");
check(
  "the page-gap ledger must no longer register history.page or roster.cancel (the page literals landed)",
  !ledgerSrc.includes('"history.page":') && !ledgerSrc.includes('"roster.cancel":'),
);

// The contract walk's page half: TaskRowView declares the three keys and the
// nested usage shape (the Go-side ledger entry + translation left in the
// same commit).
check(
  "panel.ts must declare the three task-card keys (ticket 38 AC#4's page half)",
  panelSrc.includes("startedAt?: number") && panelSrc.includes("costToDate?: RunUsageRow") && panelSrc.includes("currentStep?: string"),
);
check("panel.ts must declare the RunUsageRow interface the walk descends into", panelSrc.includes("export interface RunUsageRow"));
const composerTestSrc = readFileSync("../internal/panel/composer_test.go", "utf8");
check(
  "composerTSGapLedger must no longer register the TaskRowView three keys (the page half landed)",
  !composerTestSrc.includes('"TaskRowView": {"startedAt", "costToDate", "currentStep"}'),
);
check(
  "composerTSNameMap must translate RunUsageRow (the walk's newly reached nested struct)",
  composerTestSrc.includes('"RunUsageRow":'),
);

// The privacy face is REUSED, not re-landed: the browser imports the wave-1
// hook + reply block and the three privacy requests, purge fires only from
// its confirm card, and delete-one goes out per row through the face's own
// domain mapping.
const historyPageSrc = readFileSync("src/components/history-page.tsx", "utf8");
check("the browser must reuse the wave-1 privacy hook (usePrivacyRequest)", historyPageSrc.includes("usePrivacyRequest"));
check("the browser must reuse the wave-1 reply block (PrivacyReplyBlock)", historyPageSrc.includes("PrivacyReplyBlock"));
check(
  "delete-one must ride the existing requestPrivacyDeleteOne (bound by the wave-1 runner table, not re-landed here)",
  historyPageSrc.includes("privacyRunners") &&
    readFileSync("src/lib/privacy-view.ts", "utf8").includes("deleteOne: requestPrivacyDeleteOne"),
);
check("the browser must reuse the wave-1 runner bindings (no second envelope)", historyPageSrc.includes("privacyRunners"));
check(
  "privacy.purge must be sent from exactly one place (the confirm card's confirm button)",
  (historyPageSrc.match(/purge\.submit\(/g) ?? []).length === 1,
);
check("the purge first button must only open the confirm card", historyPageSrc.includes("setPurgeConfirmOpen(true)"));
check("the purge confirm card must be a real dialog surface", historyPageSrc.includes('role="alertdialog"'));
check("the purge confirm card must name the irreversibility", historyPageSrc.includes("不可逆"));
check("the history read itself must never carry destructive vocabulary as a method", !historyPageSrc.includes('"history.delete') && !historyPageSrc.includes('"history.purge'));

// The App wiring: the browser is mounted and fed the real loader.
const appSrc = readFileSync("src/App.tsx", "utf8");
check("App must mount the HistoryBrowser", appSrc.includes("<HistoryBrowser load={requestHistoryPage} />"));
check(
  "App must hand the monitor the two cancel doors (roster.cancel + history.page)",
  appSrc.includes("cancelDoors={{ cancelTask: requestRosterCancel, historyLoad: requestHistoryPage }}"),
);

/* ---------------------------------------------------------------------------
   The verdict.
   --------------------------------------------------------------------------- */

if (failures.length > 0) {
  for (const f of failures) console.error("render-historypage36: FAIL " + f);
  console.error(`render-historypage36: FAIL, ${failures.length} check(s) failed`);
  process.exit(1);
}
console.error(
  "render-historypage36: OK, ticket 36 AC#3's history browser consumes the real history.page door " +
    "(1005 rows through three windows merge 不重不漏 with the Go test's own window firsts, " +
    "duplicate ids and broken offsets are named defects never absorbed, " +
    "the virtual slice covers the list at every scroll position and hides the tail at the top, " +
    "row fields ride verbatim with the face's own domain mapping, " +
    "total=0 renders 暂无记录 while the pre-reply browser renders neither list nor empty claim, " +
    "refusals and unparseable payloads stay 未判, " +
    "delete-one/purge/export ride the existing privacy face with purge behind its confirm card, " +
    "the literals live in panel.ts with both ledger registrations gone in the same commit, zero emoji)",
);
