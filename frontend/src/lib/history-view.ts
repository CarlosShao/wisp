/* ============================================================================
   History view - the history.page read face's reply, its paging merge, the
   render slice, and the applied-steps forensics reading (readface wave 2,
   ticket 36 AC#3 page half + ticket 38 AC#4's cancel reading, 2026-11-17).
   ----------------------------------------------------------------------------
   WHAT ARRIVES. history.page answers over the same settings door
   (window.wispDispatch) with ONE string, exactly like cost.summary and the
   three privacy doors:

     - receipt (cmd/wisp/panel_history_wave2.go HandleHistoryPage): prose
       ("历史共 <n> 行（新的事件在前），本页第 <k> 行起的 <m> 行。") followed
       by ONE line starting with the door's own marker
       ("wisp-history-page:") whose remainder is the page JSON:
       {rows:[{id,ts,kind,tool,outcome,summary}],total,offset,limit}. The
       prose is for humans and is NEVER parsed (the same rule
       lib/privacy-view.ts and lib/cost-view.ts state); the page's facts are
       read off the payload only.

     - a refusal: RefusedEnvelopeForUser's one sentence - a negative window,
       a page over the store's ceiling, or an unopenable store refuses by
       name. The page shows it verbatim.

   THE MERGE IS EXACT, NOT ABSORBING. Windows come back over one total order
   (ts DESC, kind DESC, id DESC - internal/memory/dao_history_page_36.go's
   own comment), so consecutive windows tile the history: window N's offset
   must equal the rows already merged, and an id already on screen means the
   store reordered under the reader. Both are NAMED defects (never silently
   absorbed, never silently dropped rows): a merge that can hide a gap or a
   duplicate is how "滚动不重不漏" becomes a slogan instead of a property.

   THE SLICE. The merged list is virtualized: only the rows around the
   viewport render. The slice math is a pure exported function so the render
   instrument can pin "不重不漏" at every scroll position without a browser.

   THE APPLIED-STEPS READING (ticket 38 AC#4). After a roster.cancel receipt
   the task settles on its own and books the evidence rows itself: the
   terminal task_log row (kind "task", id = the task id) and the task's
   tool_call forensics rows (kind "tool", id = "<taskId>#<callId>" - the
   correlation the loop mints, internal/agent/loop.go:603). The page counts
   them off the history read face (已落地 = success rows, 被中止 = cancelled
   rows); it never invents a count, and an incomplete scan reads 未判 - a
   count read from a partially read store would be a fake number.
   ============================================================================ */

import { REFUSED_ENVELOPE_PREFIX } from "@/lib/privacy-view";

/** The wire marker, spelled once here; cmd/wisp's historyPagePayloadMarker
    ("wisp-history-page:") is the other side of the same string. */
export const HISTORY_PAGE_MARKER = "wisp-history-page:";

/** The door's largest served page (internal/memory's HistoryPageLimitCeiling).
    The scroller asks windows of this size; scrolling past it is what offset
    is for. Spelled as a literal - it is the face's own constant, and the
    refusal that names it is the ask, kept alive. */
export const HISTORY_PAGE_WINDOW = 500;

/** One history row as the door's payload carries it
    (internal/memory/dao_history_page_36.go HistoryPageRow's JSON tags,
    verbatim). kind is the face's own two-table vocabulary ("task" | "tool"):
    a task row's outcome IS task_log.state, a tool row's outcome IS
    tool_call.outcome, and a tool row carries no summary (args_json is the
    §5.1-sanitized forensics column and never travels under a second name). */
export interface HistoryPageRow {
  id: string;
  ts: number;
  kind: string;
  tool?: string;
  outcome?: string;
  summary?: string;
}

/** One page as the door's payload carries it (cmd/wisp/panel_history_wave2.go
    historyPagePayload). total is the count the page was read against, so the
    scroller bounds itself without a second request per page; offset/limit
    are the request's own answers back, clamped to what was served. */
export interface HistoryPagePayload {
  rows: HistoryPageRow[];
  total: number;
  offset: number;
  limit: number;
}

/** The one history request this page knows how to send - lib/panel.ts's
    requestHistoryPage bound into a loader so the components and the
    instruments share the same arity. */
export type HistoryLoader = (offset: number, limit: number) => Promise<string | null>;

/** What one history reply turned out to be, classified on the wire's own
    words (the refusal wrapper prefix, then the marker line). */
export type HistoryPageState =
  | { phase: "page"; prose: string; payload: HistoryPagePayload }
  | { phase: "refused"; raw: string }
  /** A receipt without its payload line: 未判, nothing drawn. */
  | { phase: "no-payload"; raw: string }
  /** A payload line that cannot hand over its rows: named, never partially
      drawn. */
  | { phase: "payload-error"; message: string; raw: string };

/** Classify one history.page reply. The prose rides verbatim; only the
    marker line is parsed; a payload whose rows/total are missing or
    mistyped is refused whole - the page never guesses fields. */
export function parseHistoryPageReply(reply: string): HistoryPageState {
  if (reply.startsWith(REFUSED_ENVELOPE_PREFIX)) {
    return { phase: "refused", raw: reply };
  }
  const at = reply.indexOf(HISTORY_PAGE_MARKER);
  if (at < 0) {
    return { phase: "no-payload", raw: reply };
  }
  const afterMarker = reply.slice(at + HISTORY_PAGE_MARKER.length);
  const line = afterMarker.split("\n", 1)[0] ?? "";
  let data: unknown;
  try {
    data = JSON.parse(line);
  } catch {
    return { phase: "no-payload", raw: reply };
  }
  if (typeof data !== "object" || data === null) {
    return { phase: "no-payload", raw: reply };
  }
  const raw = data as { rows?: unknown; total?: unknown; offset?: unknown; limit?: unknown };
  if (!Array.isArray(raw.rows)) {
    return {
      phase: "payload-error",
      message: "历史载荷上 rows 不在或不是数组——列表未判，不画假行",
      raw: reply,
    };
  }
  if (typeof raw.total !== "number" || typeof raw.offset !== "number" || typeof raw.limit !== "number") {
    return {
      phase: "payload-error",
      message: "历史载荷上 total/offset/limit 不在或不是数——分窗未判，不猜窗口",
      raw: reply,
    };
  }
  const rows: HistoryPageRow[] = [];
  for (const item of raw.rows) {
    if (typeof item !== "object" || item === null) {
      return { phase: "payload-error", message: "历史载荷有一行不是对象——列表未判", raw: reply };
    }
    const r = item as { id?: unknown; ts?: unknown; kind?: unknown };
    if (typeof r.id !== "string" || r.id === "" || typeof r.ts !== "number" || typeof r.kind !== "string" || r.kind === "") {
      return {
        phase: "payload-error",
        message: "历史载荷有一行缺 id/ts/kind（页面拒绝猜行）——列表未判",
        raw: reply,
      };
    }
    const rowOut: HistoryPageRow = { id: r.id, ts: r.ts, kind: r.kind };
    if (typeof (item as { tool?: unknown }).tool === "string") rowOut.tool = (item as { tool: string }).tool;
    if (typeof (item as { outcome?: unknown }).outcome === "string") {
      rowOut.outcome = (item as { outcome: string }).outcome;
    }
    if (typeof (item as { summary?: unknown }).summary === "string") {
      rowOut.summary = (item as { summary: string }).summary;
    }
    rows.push(rowOut);
  }
  return {
    phase: "page",
    prose: reply.slice(0, at).trimEnd(),
    payload: { rows, total: raw.total, offset: raw.offset, limit: raw.limit },
  };
}

/* ---------------------------------------------------------------------------
   The merge: consecutive windows tile the history, and anything else is a
   NAMED defect - never a silent absorb, never a silent drop.
   --------------------------------------------------------------------------- */

/** The scroller's merged state. duplicates carries the ids a later window
    brought back (the store reordered under the reader); anomalies carries
    window-continuity breaks (a page whose offset is not the row count the
    merge already holds). Both render as named defect lines - a merge that
    hid them would turn "不重不漏" into an unverified claim. */
export interface HistoryMergeState {
  rows: HistoryPageRow[];
  total: number;
  duplicates: string[];
  anomalies: string[];
}

export const EMPTY_HISTORY_MERGE: HistoryMergeState = { rows: [], total: 0, duplicates: [], anomalies: [] };

/** Merge one parsed page into the scroller's state. Rows append in payload
    order (the store's own total order); an id already merged is recorded as
    a duplicate and NOT appended twice, and a page whose offset is not the
    merged row count is recorded as a continuity anomaly - both keep the
    reading honest instead of quietly reshaping the list. */
export function mergeHistoryPage(state: HistoryMergeState, page: HistoryPagePayload): HistoryMergeState {
  const rows = state.rows.slice();
  const duplicates = state.duplicates.slice();
  const anomalies = state.anomalies.slice();
  if (page.offset !== rows.length) {
    anomalies.push(
      `窗口断档：已合并 ${rows.length} 行，但这一页从 offset=${page.offset} 起（存储在读取间隙发生了变动，请重新载入）`,
    );
  }
  const seen = new Set(state.rows.map((r) => r.id));
  for (const row of page.rows) {
    if (seen.has(row.id)) {
      duplicates.push(row.id);
      continue;
    }
    seen.add(row.id);
    rows.push(row);
  }
  return { rows, total: page.total, duplicates, anomalies };
}

/** Whether another window is worth asking for: the store reported more rows
    than the merge holds, AND the merge actually holds rows (a first page
    that came back empty against a positive total is the anomaly line's
    business, not a reason to loop forever). */
export function historyHasMore(state: HistoryMergeState): boolean {
  return state.rows.length > 0 && state.rows.length < state.total;
}

/* ---------------------------------------------------------------------------
   The render slice: pure virtualization math, exported for the tooth.
   --------------------------------------------------------------------------- */

/** One row's fixed render height in px. A fixed height is what makes the
    slice arithmetic exact - variable-height virtualization is a scroll bar
    that lies about where rows are. */
export const HISTORY_ROW_HEIGHT = 26;

/** The rows to render for a scroll position: the visible band plus overscan
    on both sides, clamped to the list. Pure so the instrument can pin every
    boundary (top, bottom, oversized viewport, past-the-end scroll) without a
    browser. start/end are half-open: rows[start .. end-1], never overlapping
    the band twice, never covering more than the list holds. */
export function visibleSlice(
  count: number,
  scrollTop: number,
  viewportHeight: number,
  overscan = 8,
): { start: number; end: number } {
  if (count <= 0) return { start: 0, end: 0 };
  const first = Math.floor(Math.max(0, scrollTop) / HISTORY_ROW_HEIGHT) - overscan;
  const visible = Math.ceil(Math.max(0, viewportHeight) / HISTORY_ROW_HEIGHT) + 2 * overscan;
  const start = Math.max(0, Math.min(first, count));
  const end = Math.max(start, Math.min(count, start + visible));
  return { start, end };
}

/** Display timestamp for one row. ts is the store's own unix second
    (COALESCEd: a tool row with no start AND no end stamp travels as 0 and
    sorts last - it renders as 无时刻, never as 1970). */
export function formatHistoryTs(ts: number): string {
  if (ts <= 0) return "无时刻";
  const d = new Date(ts * 1000);
  const p = (n: number) => n.toString().padStart(2, "0");
  return `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())} ${p(d.getHours())}:${p(d.getMinutes())}:${p(d.getSeconds())}`;
}

/* ---------------------------------------------------------------------------
   The applied-steps reading (ticket 38 AC#4): what the cancelled task
   already landed, read off the history face's own rows.
   --------------------------------------------------------------------------- */

/** The correlation prefix the loop mints for every tool call of a task
    ("<taskID>#<callID>", internal/agent/loop.go:603, the same join shape the
    roster's blockedOnApproval join uses). A tool row whose id carries this
    prefix belongs to that task's forensics; a task row whose id EQUALS the
    task id is its terminal row. */
export function taskForensicsPrefix(taskId: string): string {
  return `${taskId}#`;
}

/** The three readings the cancel tooth asks for, off the store's own rows:
    终态名 (the terminal task_log row's state, or the live roster row's D43
    name when the host filed one), 已落地步数 (the task's tool_call rows with
    outcome "success") and 被中止步数 (outcome "cancelled"). scanComplete
    says the caller read the WHOLE history before counting: a count from a
    partially read store would be a fake number, so an incomplete scan
    carries null - which the page renders as 未判, never as 0. */
export interface AppliedStepsReading {
  terminalKnown: boolean;
  terminal: string;
  landed: number | null;
  aborted: number | null;
  /** The forensics rows that matched the task (terminal row + tool rows). */
  matchedRows: number;
}

export function appliedStepsReading(
  historyRows: HistoryPageRow[],
  taskId: string,
  rowStatus?: { status: string; statusKnown: boolean },
  scanComplete = true,
): AppliedStepsReading {
  const prefix = taskForensicsPrefix(taskId);
  let terminalFromStore = "";
  let landed = 0;
  let aborted = 0;
  let matchedRows = 0;
  for (const row of historyRows) {
    if (row.kind === "task" && row.id === taskId) {
      matchedRows += 1;
      if (typeof row.outcome === "string" && row.outcome !== "") terminalFromStore = row.outcome;
      continue;
    }
    if (row.kind === "tool" && row.id.startsWith(prefix)) {
      matchedRows += 1;
      if (row.outcome === "success") landed += 1;
      if (row.outcome === "cancelled") aborted += 1;
    }
  }
  const terminalKnown = (rowStatus?.statusKnown === true && rowStatus.status !== "") || terminalFromStore !== "";
  const terminal = rowStatus?.statusKnown === true && rowStatus.status !== "" ? rowStatus.status : terminalFromStore;
  return {
    terminalKnown,
    terminal,
    landed: scanComplete ? landed : null,
    aborted: scanComplete ? aborted : null,
    matchedRows,
  };
}

/** What a whole-history forensics scan turned out to be. complete=false
    means the reading stays 未判; the failure halves carry the door's own
    words (verbatim) so the page can say WHY nothing was counted. */
export interface ForensicsScanOutcome {
  complete: boolean;
  rows: HistoryPageRow[];
  scannedRows: number;
  total: number;
  refused?: string;
  payloadError?: string;
  noPayload?: boolean;
  noHost?: boolean;
}

/** Page through the WHOLE history (windows of HISTORY_PAGE_WINDOW) and keep
    the rows that belong to one task: its terminal task row and its tool
    forensics rows. The store's own total bounds the loop; a refusal, an
    unparseable payload or a no-host door stops the scan with complete=false
    - the reading is then 未判, never a partial count dressed up as data. */
export async function scanTaskForensics(load: HistoryLoader, taskId: string): Promise<ForensicsScanOutcome> {
  const prefix = taskForensicsPrefix(taskId);
  const out: ForensicsScanOutcome = { complete: false, rows: [], scannedRows: 0, total: 0 };
  let offset = 0;
  for (;;) {
    const reply = await load(offset, HISTORY_PAGE_WINDOW);
    if (reply === null) {
      out.noHost = true;
      return out;
    }
    const state = parseHistoryPageReply(reply);
    if (state.phase === "refused") {
      out.refused = state.raw;
      return out;
    }
    if (state.phase === "no-payload") {
      out.noPayload = true;
      return out;
    }
    if (state.phase === "payload-error") {
      out.payloadError = state.message;
      return out;
    }
    out.total = state.payload.total;
    out.scannedRows += state.payload.rows.length;
    for (const row of state.payload.rows) {
      if ((row.kind === "task" && row.id === taskId) || (row.kind === "tool" && row.id.startsWith(prefix))) {
        out.rows.push(row);
      }
    }
    const served = state.payload.offset + state.payload.rows.length;
    if (state.payload.rows.length === 0 || served >= state.payload.total) {
      out.complete = true;
      return out;
    }
    offset = served;
  }
}
