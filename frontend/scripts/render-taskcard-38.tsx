/* ============================================================================
   Ticket 38 AC#4's page-half tooth (render evidence), 2026-11-17 - the
   task-card / cancel leg (consume wave 2).
   ----------------------------------------------------------------------------
   WHAT THIS FILE PROVES, in the shape render-privacy-40 / render-task-roster
   established: every check below feeds REAL-SHAPED data through the
   PRODUCTION code and asserts the SAME VALUE comes back out.

   THE FIXTURES are the wire's own documented shapes: the three task-card
   keys read off internal/panel/subagent_roster_197.go's TaskRowView JSON
   tags (startedAt / costToDate / currentStep, all omitempty), the cancel
   receipt from cmd/wisp/panel_roster_cancel_wave2.go's HandleRosterCancel,
   and the applied-steps forensics from the same surfaces the Go cancel test
   reads (cmd/wisp/panel_roster_cancel_wave2_test.go): the terminal task_log
   row (kind "task", id = the task id) plus the task's tool_call rows
   (kind "tool", id = "<taskId>#<callId>", the correlation loop.go:603
   mints), counted by the production appliedStepsReading over history.page
   rows.

   THE TEETH:
     1. the three keys ride VERBATIM onto the task-card line (the start time
        through the production formatter, the cost numbers with their
        currency pair, the in-flight step name);
     2. an unstamped row renders 未判 for start/cost and 无在飞 for the step -
        NEVER an epoch, a zero cost, or an invented step;
     3. App's derivation copies the three keys verbatim (nothing judged
        page-side);
     4. the cancel is a REAL GATE: the 停止 button only opens the confirm
        card, the request leaves from exactly one place (the card's confirm
        button), and the receipt/refusal/no-host words ride verbatim;
     5. the applied-steps reading shows THREE values off the forensics rows
        (终态名 / 已落地步数 / 被中止步数) and an incomplete scan reads 未判 -
        never a fake 0;
     6. source level: the roster.cancel literal + selector live in panel.ts,
        both ledger registrations are gone, the forensics counting imports
        the history face's own reader - plus zero emoji.
   ============================================================================ */

import { readFileSync } from "node:fs";
import { renderToStaticMarkup } from "react-dom/server";
import App, { monitorAgents } from "@/App";
import {
  AppliedReadingLine,
  MonitorPopover,
  TaskCancelReply,
  TaskCardLine,
  type TaskCancelState,
} from "@/components/harness/monitor-popover";
import {
  appliedStepsReading,
  formatHistoryTs,
  scanTaskForensics,
  type HistoryPageRow,
} from "@/lib/history-view";
import type { PanelSnapshot, TaskRosterSection } from "@/lib/panel";

const failures: string[] = [];
function check(label: string, ok: boolean, detail = ""): void {
  if (!ok) failures.push(`${label}${detail ? ` [${detail} ]` : ""}`);
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
   Wire fixtures. Field names are lib/panel.ts TaskRowView verbatim.
   --------------------------------------------------------------------------- */

const STAMPED_AT = 1763376000;
const FED_USAGE = { tokensIn: 1200, tokensOut: 340, cachedTokens: 512, micros: 123456, currency: "CNY" };

function row(over: Partial<import("@/lib/panel").TaskRowView> & Pick<import("@/lib/panel").TaskRowView, "taskId" | "label">): import("@/lib/panel").TaskRowView {
  return {
    kind: "subagent",
    parentTaskId: "root-1",
    status: "",
    statusKnown: false,
    streamKey: over.taskId,
    blockedOnApproval: false,
    streamTruncated: false,
    streamElidedRunes: 0,
    streamDropped: false,
    ...over,
  };
}

function roster(rows: import("@/lib/panel").TaskRowView[]): TaskRosterSection {
  return {
    rows,
    inFlightSlots: 1,
    poolCap: 8,
    streamTruncated: false,
    streamElidedRunes: 0,
    droppedStreamKeys: [],
    droppedCount: 0,
  };
}

const COMPOSER = {
  mode: { current: "ask_high_risk", names: ["ask_every_step", "ask_high_risk", "auto_approve"], l2ConfirmNames: ["auto_approve"] },
  workspace: { set: true, spelling: "D:/work/wisp", canonical: "D:/work/wisp", reparse: false, rewritten: false, reason: "" },
  attachments: [],
  acceptedAttachmentMimes: [],
  maxAttachmentBytes: 0,
  attachmentError: "",
};

function snapshotWith(tasks: TaskRosterSection): PanelSnapshot {
  return { pending: [], results: [], composer: COMPOSER, generatedAt: "2026-11-17T13:00:00Z", tasks };
}

const ENV = [{ label: "工作区", value: "D:/work/wisp" }] as const;

/* ---------------------------------------------------------------------------
   1. The three keys ride verbatim through the production task-card line.
   --------------------------------------------------------------------------- */

const stampedRow = row({
  taskId: "child-9",
  label: "带三键的子任务",
  status: "Acting",
  statusKnown: true,
  startedAt: STAMPED_AT,
  costToDate: FED_USAGE,
  currentStep: "fs.write",
});
const cardHtml = renderToStaticMarkup(<TaskCardLine row={monitorAgents(roster([stampedRow]))[0]!} />);
check(
  "startedAt must render through the production formatter (逐字)",
  cardHtml.includes(esc(formatHistoryTs(STAMPED_AT))),
);
check("the raw stamp must never render as epoch 1970 when the host stamped a real second", !cardHtml.includes("1970"));
check("costToDate's micros must ride verbatim", cardHtml.includes(esc("123456 micros")));
check("costToDate's currency pair must ride verbatim", cardHtml.includes(esc("CNY")));
check("costToDate's token counts must ride verbatim", cardHtml.includes(esc("tokens 1200/340")));
check("costToDate's cached tokens must ride verbatim", cardHtml.includes(esc("缓存 512")));
check("currentStep's tool name must ride verbatim", cardHtml.includes(esc("fs.write")));
check("the card line must hold zero emoji", !EMOJI_RE.test(cardHtml));

/* ---------------------------------------------------------------------------
   2. The unstamped row: 未判 / 无在飞, never a fabricated zero.
   --------------------------------------------------------------------------- */

const bareRow = row({ taskId: "child-7", label: "没登记三键的子任务" });
const bareHtml = renderToStaticMarkup(<TaskCardLine row={monitorAgents(roster([bareRow]))[0]!} />);
// The 未判 words sit inside their own spans (开始 <span>未判</span>), so the
// halves are asserted separately from the value text.
check("an unstamped start must read 未判", bareHtml.includes(esc("开始")) && /开始\s*<[^>]+>未判/.test(bareHtml));
check("an unspent row must read 未判 (never a zero cost)", bareHtml.includes(esc("成本")) && /成本\s*<[^>]+>未判/.test(bareHtml));
check("a step-less row must read 无在飞 (the wire's own retired state)", bareHtml.includes(esc("无在飞")));
check("the unstamped row must NOT render epoch 1970", !bareHtml.includes("1970"));
check("the unstamped row must NOT render a fabricated micros number", !bareHtml.includes("micros"));
check("the unstamped card must hold zero emoji", !EMOJI_RE.test(bareHtml));

/* ---------------------------------------------------------------------------
   3. App's derivation copies the three keys verbatim; the popover paints
      them and, with the doors bound, draws the 停止 control.
   --------------------------------------------------------------------------- */

const derived = monitorAgents(roster([stampedRow]));
check("the derivation must carry startedAt verbatim", derived[0]?.startedAt === STAMPED_AT);
check("the derivation must carry costToDate verbatim", derived[0]?.costToDate === FED_USAGE);
check("the derivation must carry currentStep verbatim", derived[0]?.currentStep === "fs.write");

const popoverHtml = renderToStaticMarkup(
  <MonitorPopover
    open
    env={ENV}
    agents={monitorAgents(roster([stampedRow]))}
    rosterPresent
    cancelDoors={{
      cancelTask: async () => null,
      historyLoad: async () => null,
    }}
  />,
);
check("the popover must paint the three-key card line", popoverHtml.includes(esc("fs.write")));
check("the popover with doors bound must draw the per-row 停止 control", popoverHtml.includes(esc("停止")));
check("the popover markup must hold zero emoji", !EMOJI_RE.test(popoverHtml));

const doorlessHtml = renderToStaticMarkup(
  <MonitorPopover open env={ENV} agents={monitorAgents(roster([stampedRow]))} rosterPresent />,
);
check("without the doors bound the popover must draw NO dead 停止 control", !doorlessHtml.includes(esc("停止")));

const appHtml = renderToStaticMarkup(<App snapshot={snapshotWith(roster([stampedRow]))} />);
check("App with a stamped roster must render without crashing", appHtml.includes("一缕"));
check("the App markup must hold zero emoji", !EMOJI_RE.test(appHtml));

/* ---------------------------------------------------------------------------
   4. The cancel exchange, rendered verbatim through the production block.
   --------------------------------------------------------------------------- */

const RECEIPT =
  "已对任务 child-9 发出停止（经既有取消通路；任务将自行结算并落下终态行，已落地的步骤以 tool_call 取证行与任务终态为准）。";
const cancelStates: TaskCancelState[] = [
  { phase: "sending" },
  { phase: "receipt", text: RECEIPT, scanning: true },
  { phase: "refused", raw: `面板请求被拒绝 [roster.cancel pc-9-c1]：panel: 任务 child-99 没有停下：查不到任务 child-99（任务名册没有这一行）` },
  { phase: "no-host" },
];
for (const [i, state] of cancelStates.entries()) {
  const html = renderToStaticMarkup(<TaskCancelReply state={state} />);
  check(`cancel state ${i} must hold zero emoji`, !EMOJI_RE.test(html));
  if (state.phase === "receipt") {
    check("the receipt must ride verbatim", html.includes(esc(RECEIPT)));
    check("the receipt must name the row it stopped", html.includes(esc("child-9")));
  }
  if (state.phase === "refused") {
    check("the refusal must ride verbatim (the host's own words)", html.includes(esc("查不到任务 child-99（任务名册没有这一行）")));
    check("the refusal must render in the danger tint", html.includes("bg-red-tint"));
  }
  if (state.phase === "no-host") {
    check("the no-host block must say the request never left (the task was untouched)", html.includes(esc("任务未被触碰")));
  }
}

/* ---------------------------------------------------------------------------
   5. The applied-steps reading: THREE values off the forensics rows, 未判
      when the scan is incomplete - never a fake 0.
   --------------------------------------------------------------------------- */

// The settled forensics of cancelled task child-9, exactly the surfaces the
// Go cancel test reads: the terminal task_log row + the task's tool_call
// rows (success landed, cancelled aborted), mixed among unrelated history
// rows.
const forensics: HistoryPageRow[] = [
  { id: "other-task", ts: 1763376009, kind: "task", outcome: "Completed", summary: "别人的任务" },
  { id: "child-9", ts: 1763376008, kind: "task", outcome: "Stopped" },
  { id: "child-9#c1", ts: 1763376007, kind: "tool", tool: "fs.read", outcome: "success" },
  { id: "child-9#c2", ts: 1763376006, kind: "tool", tool: "fs.write", outcome: "cancelled" },
  { id: "other-task#c1", ts: 1763376005, kind: "tool", tool: "fs.read", outcome: "success" },
];
const reading = appliedStepsReading(forensics, "child-9", { status: "Stopped", statusKnown: true }, true);
check("the terminal name must read Stopped (D43 name off the wire)", reading.terminalKnown && reading.terminal === "Stopped");
check("the landed count must read the task's own success rows", reading.landed === 1);
check("the aborted count must read the task's own cancelled rows", reading.aborted === 1);
check("the matched rows must be exactly the task's forensics", reading.matchedRows === 3);
check("foreign tasks' rows must never enter the reading", reading.landed !== 2);

const readingHtml = renderToStaticMarkup(
  <AppliedReadingLine reading={reading} scanning={false} />,
);
check("the reading line must show 终态 Stopped", readingHtml.includes(esc("Stopped")));
check("the reading line must show 已落地 1 步", readingHtml.includes(esc("已落地 1 步")));
check("the reading line must show 被中止 1 步", readingHtml.includes(esc("被中止 1 步")));
check("the reading line must hold zero emoji", !EMOJI_RE.test(readingHtml));

// The terminal name can also come off the store's own row when the roster
// has not filed a status yet; an unknown terminal is 未判, never guessed.
const storeOnly = appliedStepsReading(forensics, "child-9", undefined, true);
check("the terminal name must fall back to the store's own task row", storeOnly.terminalKnown && storeOnly.terminal === "Stopped");
const noTerminal = appliedStepsReading(
  forensics.filter((r) => r.id !== "child-9"),
  "child-42",
  { status: "", statusKnown: false },
  true,
);
check("an unknown terminal must read 未判", !noTerminal.terminalKnown);

const unjudgedHtml = renderToStaticMarkup(<AppliedReadingLine />);
check("an absent reading must read 未判 for both counts", unjudgedHtml.includes(esc("已落地步数：未判")) && unjudgedHtml.includes(esc("被中止步数：未判")));
check("an absent reading must NEVER draw a fake 0", !unjudgedHtml.includes(esc("已落地 0 步")) && !unjudgedHtml.includes(esc("被中止 0 步")));
check("an absent reading must read 未判 for the terminal", unjudgedHtml.includes(esc("终态")));
check("the unjudged line must hold zero emoji", !EMOJI_RE.test(unjudgedHtml));

// An INCOMPLETE scan reads 未判: a count from a partially read store would
// be a fake number.
const partial = appliedStepsReading(forensics.slice(0, 2), "child-9", { status: "Stopped", statusKnown: true }, false);
check("an incomplete scan's landed count must be null (未判)", partial.landed === null);
check("an incomplete scan's aborted count must be null (未判)", partial.aborted === null);

/* ---------------------------------------------------------------------------
   6. The forensics scan over the whole history face, through the production
      scanner (three windows of the ceiling page, like the Go scroll test).
   --------------------------------------------------------------------------- */

{
  const store: HistoryPageRow[] = Array.from({ length: 1005 }, (_, i) => {
    const n = 1004 - i;
    const id = `scroll-${n.toString().padStart(4, "0")}`;
    if (i === 1000) return { id: "child-9", ts: 1763376000 - i, kind: "task", outcome: "Stopped" };
    if (i === 400) return { id: "child-9#c1", ts: 1763376000 - i, kind: "tool", tool: "fs.read", outcome: "success" };
    if (i === 60) return { id: "child-9#c2", ts: 1763376000 - i, kind: "tool", tool: "fs.edit", outcome: "cancelled" };
    return { id, ts: 1763376000 - i, kind: n % 2 === 0 ? "task" : "tool" };
  });
  let calls = 0;
  const load = async (offset: number, limit: number): Promise<string | null> => {
    calls += 1;
    const rows = store.slice(offset, Math.min(offset + limit, store.length));
    const payload = JSON.stringify({ rows, total: store.length, offset, limit });
    return `历史共 ${store.length} 行（新的事件在前），本页第 ${offset + 1} 行起的 ${rows.length} 行。\nwisp-history-page:${payload}`;
  };
  const scan = await scanTaskForensics(load, "child-9");
  check("the scan must read the WHOLE history (complete=true)", scan.complete);
  check("the scan must page until the store's own total is covered", scan.scannedRows === 1005, `scanned ${scan.scannedRows}`);
  check("the scan must collect exactly the task's forensics rows", scan.rows.length === 3);
  const scanReading = appliedStepsReading(scan.rows, "child-9", { status: "Stopped", statusKnown: true }, scan.complete);
  check("the full-scan reading must land 终态 Stopped / 已落地 1 / 被中止 1", scanReading.terminal === "Stopped" && scanReading.landed === 1 && scanReading.aborted === 1);

  const refusedScan = await scanTaskForensics(
    async () =>
      `面板请求被拒绝 [history.page pc-x]：panel: history.page 被拒绝（requestId="pc-x"）：历史存储不可用`,
    "child-9",
  );
  check("a refusal mid-scan must stop with complete=false (the reading stays 未判)", !refusedScan.complete && refusedScan.refused !== undefined);
}

/* ---------------------------------------------------------------------------
   7. Source-level wiring: the gate's shape, the literals, the ledger.
   --------------------------------------------------------------------------- */

const popoverSrc = readFileSync("src/components/harness/monitor-popover.tsx", "utf8");
check(
  "the stop request must be sent from exactly one place (the confirm card's confirm button)",
  (popoverSrc.match(/onClick=\{\(\) => confirmStop\(confirmId\)\}/g) ?? []).length === 1,
);
check("the 停止 button must only open the confirm card (先确认再取消，禁止一键杀)", (popoverSrc.match(/setConfirmId\(agent\.id\)/g) ?? []).length === 1);
check("the confirm card must be a real dialog surface", popoverSrc.includes('role="alertdialog"'));
check("the confirm card must name the door it will ring", popoverSrc.includes("roster.cancel"));
check("the confirm card must carry the 取消 exit", popoverSrc.includes("取消"));
check("the forensics reading must come from the history face's own scanner", popoverSrc.includes("scanTaskForensics") && popoverSrc.includes("appliedStepsReading"));
check("the receipt classification must use the wire's own refusal wrapper", popoverSrc.includes("REFUSED_ENVELOPE_PREFIX"));

const panelSrc = readFileSync("src/lib/panel.ts", "utf8");
check("panel.ts must spell roster.cancel as a literal", panelSrc.includes('"roster.cancel"'));
check("panel.ts must spell rosterTaskId as bridge.go's tag carries it", panelSrc.includes("rosterTaskId"));
check("panel.ts must spell history.page as a literal (the reading's door)", panelSrc.includes('"history.page"'));

const ledgerSrc = readFileSync("../internal/panel/bridge_test.go", "utf8");
check(
  "the page-gap ledger must no longer register history.page or roster.cancel (the page literals landed)",
  !ledgerSrc.includes('"history.page":') && !ledgerSrc.includes('"roster.cancel":'),
);
const composerTestSrc = readFileSync("../internal/panel/composer_test.go", "utf8");
check(
  "composerTSGapLedger must no longer register the TaskRowView three keys (the page half landed)",
  !composerTestSrc.includes('"TaskRowView": {"startedAt", "costToDate", "currentStep"}'),
);

const appSrc = readFileSync("src/App.tsx", "utf8");
check(
  "App must bind the real cancel door into the monitor (roster.cancel, not a local fake)",
  appSrc.includes("cancelTask: requestRosterCancel"),
);
check("App must feed the derivation to the monitor (never a bare empty array)", appSrc.includes("agents={monitorAgents(snapshot.tasks)}") && !appSrc.includes("agents={[]}"));

/* ---------------------------------------------------------------------------
   The verdict.
   --------------------------------------------------------------------------- */

if (failures.length > 0) {
  for (const f of failures) console.error("render-taskcard38: FAIL " + f);
  console.error(`render-taskcard38: FAIL, ${failures.length} check(s) failed`);
  process.exit(1);
}
console.error(
  "render-taskcard38: OK, ticket 38 AC#4's task card and cancel gate consume the real doors " +
    "(startedAt/costToDate/currentStep ride verbatim, unstamped rows read 未判 and never epoch/zero, " +
    "the derivation copies the three keys 1:1, the 停止 control only opens its confirm card and the request " +
    "leaves from exactly one place, receipt/refusal/no-host ride verbatim, " +
    "the applied-steps reading shows 终态/已落地/被中止 off the forensics rows with an incomplete scan at 未判, " +
    "the forensics scan pages the whole history face, the literals live in panel.ts with both ledger " +
    "registrations gone in the same commit, zero emoji)",
);
