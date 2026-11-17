/* ============================================================================
   History browser (ticket 36 AC#3's page half, readface wave 2 2026-11-17).
   ----------------------------------------------------------------------------
   WHAT THIS IS. The history & transcript browser the ticket asked for, fed
   by the REAL paged read face: every window comes off lib/panel.ts's
   requestHistoryPage (method "history.page", internal/panel/bridge.go's
   MethodHistoryPage, answered by cmd/wisp/panel_history_wave2.go over
   internal/memory/dao_history_page_36.go's paged read of the two D35 tables
   task_log + tool_call). Rows render {ts,kind,tool,outcome,summary} verbatim.

   HOW 1000+ ROWS SCROLL. Two independent mechanisms, both exact:
     - DATA paging: windows of HISTORY_PAGE_WINDOW (the store's own ceiling,
       500) are fetched as the list scrolls to its bottom and merged by
       lib/history-view.ts's mergeHistoryPage - consecutive windows tile the
       one total order, and a gap or an id seen twice is a NAMED defect line,
       never a silent absorb.
     - RENDER virtualization: only the viewport band plus overscan renders,
       via the pure visibleSlice over a fixed row height. The instrument pins
       the slice math at every boundary, so "不重不漏" is a property of the
       list, not a claim about it.

   THE EMPTY AND THE REFUSED. total=0 renders 暂无记录; a refusal renders the
   door's own sentence verbatim; a receipt without its payload line or a
   payload that cannot hand over its rows renders 未判 - never a fabricated
   list, never a fabricated count.

   DELETE / PURGE / EXPORT RIDE THE EXISTING PRIVACY FACE, and only it
   (禁止新造第二套): the destructive verbs of this same data keep their only
   route - lib/panel.ts's three privacy requests over internal/memory's
   privacy layer, landed by the wave-1 privacy leg - and this page reuses the
   wave-1 page half wholesale: usePrivacyRequest + PrivacyReplyBlock from
   privacy-settings/privacy-view, the purge confirm card included. Nothing
   here opens a second destructive door. The two domain words this strip
   offers are NOT a copied PrivacyDomains() list: they are the read face's
   own two table names (dao_history_page_36.go pages exactly task_log +
   tool_call), and the door still judges every word against PrivacyDomains()
   - its refusal rides verbatim.

   Colours: token utilities only, zero literals; no emoji.
   ============================================================================ */

import { useCallback, useEffect, useRef, useState, type UIEvent } from "react";
import { RotateCcw, ScrollText, TriangleAlert } from "lucide-react";
import { Button } from "@/components/ai-native/button";
import { PrivacyReplyBlock } from "@/components/privacy-settings";
import { cn } from "@/lib/cn";
import { privacyRunners, usePrivacyRequest } from "@/lib/privacy-view";
import {
  EMPTY_HISTORY_MERGE,
  HISTORY_PAGE_WINDOW,
  HISTORY_ROW_HEIGHT,
  formatHistoryTs,
  historyHasMore,
  mergeHistoryPage,
  parseHistoryPageReply,
  visibleSlice,
  type HistoryLoader,
  type HistoryMergeState,
  type HistoryPageRow,
} from "@/lib/history-view";

/** The browser viewport's fixed height in px - the second constant the slice
    arithmetic needs (the first is HISTORY_ROW_HEIGHT). */
export const HISTORY_VIEWPORT_HEIGHT = 380;

/** One history row's domain word, for the delete-one privacy door. The
    mapping is the read face's own: kind "task" IS task_log, kind "tool" IS
    tool_call (dao_history_page_36.go's own comment). A kind word the face
    does not define maps to "" - the delete button then renders disabled
    instead of sending a word nobody owns. */
export function historyDomainOf(kind: string): string {
  if (kind === "task") return "task_log";
  if (kind === "tool") return "tool_call";
  return "";
}

/** The face's own display word for a row kind. A word the face does not
    define renders as itself - never as a guessed label. */
function kindLabel(kind: string): string {
  if (kind === "task") return "任务";
  if (kind === "tool") return "工具";
  return kind;
}

/** One merged history row: the wire's five display fields, nothing derived.
    ts renders through formatHistoryTs (0 = the store's own "no stamp", shown
    as 无时刻, never as 1970); outcome rides verbatim whatever word it is -
    the store's own vocabularies, not a page-side judgement. Exported pure
    for the render instrument. */
export function HistoryRowLine({
  row,
  onDelete,
  deleteDisabled,
}: {
  row: HistoryPageRow;
  onDelete: (row: HistoryPageRow) => void;
  deleteDisabled: boolean;
}) {
  const domain = historyDomainOf(row.kind);
  return (
    <div
      className="flex items-center gap-2 border-b border-line px-3 text-[11.5px]"
      style={{ height: HISTORY_ROW_HEIGHT }}
      role="listitem"
    >
      <span className="w-[118px] shrink-0 font-mono tabular-nums text-ink-3" title={`ts=${row.ts}`}>
        {formatHistoryTs(row.ts)}
      </span>
      <span className="w-[38px] shrink-0 text-ink-2">{kindLabel(row.kind)}</span>
      <span className="w-[110px] shrink-0 truncate font-mono text-ink" title={row.tool ?? ""}>
        {row.tool ?? ""}
      </span>
      <span className="w-[76px] shrink-0 truncate font-mono text-ink-2" title={row.outcome ?? ""}>
        {row.outcome ?? ""}
      </span>
      <span className="min-w-0 flex-1 truncate text-ink-2" title={row.summary ?? ""}>
        {row.summary ?? ""}
      </span>
      <button
        type="button"
        aria-label={`删除这条记录（${domain} ${row.id}）`}
        disabled={deleteDisabled || domain === ""}
        onClick={() => onDelete(row)}
        className="shrink-0 rounded-chip px-1.5 py-0.5 text-[10.5px] text-ink-3 transition-colors duration-150 hover:bg-hover hover:text-red disabled:pointer-events-none disabled:opacity-40"
      >
        删除
      </button>
    </div>
  );
}

/** The browser's empty state, its own exported component so the render
    instrument can pin the sentence: 暂无记录 is ONLY the store's own
    total=0 answer, never the pre-reply state of a browser that has not
    asked yet. */
export function HistoryEmptyBlock() {
  return (
    <div className="rounded-card border border-line bg-surface px-3.5 py-6">
      <p className="text-center text-[13px] text-ink">暂无记录</p>
      <p className="mt-1 text-center text-[11px] leading-[1.6] text-ink-3">
        存储自报 total=0：任务历史与工具取证两张表都还没有行。
      </p>
    </div>
  );
}

/** The refusal / unparseable / no-host blocks, exported pure: whatever the
    door answered comes back out VERBATIM, and a reading the door did not
    hand over says 未判 and never draws a fabricated list. */
export function HistoryRefusalBlock({ phase, raw, message }: { phase: ListPhase; raw: string; message: string }) {
  if (phase === "no-host") {
    return (
      <div className="rounded-card border border-line bg-surface px-3.5 py-3">
        <p className="text-[12px] leading-[1.7] text-ink-2">
          没接到本机宿主（这一页没有 wispDispatch 门）：请求没有发出去，历史未判——这里不画任何假列表。
        </p>
      </div>
    );
  }
  if (phase === "refused") {
    return (
      <div className="rounded-card border border-red/40 bg-red-tint px-3.5 py-3">
        <p className="text-[12px] leading-[1.7] text-red">门拒绝了这个请求，原话逐字如下：</p>
        <p className="mt-1.5 whitespace-pre-wrap break-all font-mono text-[11px] leading-[1.6] text-ink">{raw}</p>
      </div>
    );
  }
  return (
    <div className="rounded-card border border-orange/40 bg-orange-tint px-3.5 py-3">
      <p className="text-[12px] leading-[1.7] text-orange">
        {phase === "no-payload"
          ? "回执上没有解析到历史载荷行（wisp-history-page: 行不在或不是 JSON）：列表未判，不画假行。原话逐字如下："
          : `${message} 原话逐字如下：`}
      </p>
      <p className="mt-1.5 whitespace-pre-wrap break-all font-mono text-[11px] leading-[1.6] text-ink">{raw}</p>
    </div>
  );
}

/** The classify phase of the last window reply, kept beside the merged list
    so a refusal never looks like an empty list. */
export type ListPhase = "idle" | "no-host" | "refused" | "no-payload" | "payload-error";

export function HistoryBrowser({ load }: { load: HistoryLoader }) {
  const [merge, setMerge] = useState<HistoryMergeState>(EMPTY_HISTORY_MERGE);
  const [phase, setPhase] = useState<ListPhase>("idle");
  const [raw, setRaw] = useState("");
  const [message, setMessage] = useState("");
  const [prose, setProse] = useState("");
  const [loading, setLoading] = useState(false);
  const [scrollTop, setScrollTop] = useState(0);
  /** Whether at least one window REQUEST has gone out. 暂无记录 needs this:
      before the first reply the browser shows neither a list nor an empty
      claim - "the store has no rows" is the store's own answer, not this
      page's guess while the request is still in flight. */
  const [started, setStarted] = useState(false);

  const seqRef = useRef(0);
  const loadingRef = useRef(false);
  const mergeRef = useRef(merge);
  mergeRef.current = merge;

  /** Ask for one window and merge it. `base` is the merge the window joins:
      EMPTY for the first window / a reload, the live merge for a scroll-in
      window. The seq guard retires stale replies (StrictMode mounts, rapid
      reloads) exactly like lib/privacy-view.ts's hook. */
  const fetchWindow = useCallback(
    (offset: number, base: HistoryMergeState) => {
      const mine = ++seqRef.current;
      loadingRef.current = true;
      setStarted(true);
      setLoading(true);
      void load(offset, HISTORY_PAGE_WINDOW).then((reply) => {
        if (mine !== seqRef.current) return;
        loadingRef.current = false;
        setLoading(false);
        if (reply === null) {
          setPhase("no-host");
          return;
        }
        const parsed = parseHistoryPageReply(reply);
        if (parsed.phase === "refused") {
          setPhase("refused");
          setRaw(parsed.raw);
          return;
        }
        if (parsed.phase === "no-payload") {
          setPhase("no-payload");
          setRaw(parsed.raw);
          return;
        }
        if (parsed.phase === "payload-error") {
          setPhase("payload-error");
          setMessage(parsed.message);
          setRaw(parsed.raw);
          return;
        }
        setProse(parsed.prose);
        setPhase("idle");
        setMerge(mergeHistoryPage(base, parsed.payload));
      });
    },
    [load],
  );

  /** First window on mount - a browser that opens empty would look like an
      empty history, which is the 暂无记录 state's job, not this one's. */
  useEffect(() => {
    fetchWindow(0, EMPTY_HISTORY_MERGE);
  }, [fetchWindow]);

  function reload() {
    seqRef.current += 1;
    loadingRef.current = false;
    setLoading(false);
    setPhase("idle");
    setRaw("");
    setMessage("");
    setProse("");
    setMerge(EMPTY_HISTORY_MERGE);
    setScrollTop(0);
    fetchWindow(0, EMPTY_HISTORY_MERGE);
  }

  function onScroll(event: UIEvent<HTMLDivElement>) {
    const el = event.currentTarget;
    setScrollTop(el.scrollTop);
    const nearBottom = el.scrollTop + el.clientHeight >= el.scrollHeight - 3 * HISTORY_ROW_HEIGHT;
    if (nearBottom && !loadingRef.current && phase === "idle" && historyHasMore(mergeRef.current)) {
      fetchWindow(mergeRef.current.rows.length, mergeRef.current);
    }
  }

  /* ------------------------- the existing privacy face ---------------------- */

  /** The two domain words this strip can address, in a fixed order. These
      are the read face's own tables, not a copied D35 vocabulary - the
      caption says so, and the door's refusal is still the judge of any
      word. */
  const HISTORY_DOMAINS = ["task_log", "tool_call"] as const;
  type HistoryDomain = (typeof HISTORY_DOMAINS)[number];

  const [domain, setDomain] = useState<HistoryDomain>("task_log");
  const [purgeConfirmOpen, setPurgeConfirmOpen] = useState(false);

  // The wave-1 hook and reply block, reused wholesale: one receipt per verb,
  // refusals verbatim, the export payload's own records counted - nothing
  // here re-implements any of it (禁止新造第二套).
  const del = usePrivacyRequest("plain", privacyRunners.deleteOne);
  const purge = usePrivacyRequest("plain", privacyRunners.purge);
  const exp = usePrivacyRequest("export", privacyRunners.export);

  function deleteRow(row: HistoryPageRow) {
    del.submit(historyDomainOf(row.kind), row.id);
  }

  /* --------------------------------- render -------------------------------- */

  const slice = visibleSlice(merge.rows.length, scrollTop, HISTORY_VIEWPORT_HEIGHT);
  const emptyHistory =
    merge.rows.length === 0 && merge.total === 0 && phase === "idle" && !loading && started;

  return (
    <div className="flex min-w-0 flex-col gap-3">
      <div className="flex items-center gap-2">
        <ScrollText aria-hidden="true" className="size-4 shrink-0 text-ink-2" />
        <span className="min-w-0 flex-1 text-[12.5px] font-semibold text-ink">历史记录</span>
        <button
          type="button"
          aria-label="重新载入历史"
          onClick={reload}
          className="flex items-center gap-1 rounded-chip px-1.5 py-0.5 text-[11px] text-ink-3 transition-colors duration-150 hover:bg-hover hover:text-ink"
        >
          <RotateCcw aria-hidden="true" size={11} strokeWidth={2} />
          重新载入
        </button>
      </div>

      {/* The door's own prose sentence (历史共 N 行…), verbatim, while the
          list reads clean. A refusal replaces it with the refusal. */}
      {phase === "idle" && prose !== "" ? (
        <p className="text-[11.5px] leading-[1.6] text-ink-3">{prose}</p>
      ) : null}

      {phase !== "idle" ? <HistoryRefusalBlock message={message} phase={phase} raw={raw} /> : null}

      {emptyHistory ? <HistoryEmptyBlock /> : null}

      {/* The named defects: a gap or a duplicate is the store reordering
          under the reader, and the list says so instead of hiding it. */}
      {merge.anomalies.map((a) => (
        <p className="text-[11.5px] leading-[1.6] text-warn" key={a} role="alert">
          <TriangleAlert aria-hidden="true" className="mr-1 inline size-3 align-[-1px]" />
          {a}
        </p>
      ))}
      {merge.duplicates.length > 0 ? (
        <p className="text-[11.5px] leading-[1.6] text-warn" role="alert">
          <TriangleAlert aria-hidden="true" className="mr-1 inline size-3 align-[-1px]" />
          这一页带回 {merge.duplicates.length} 行已在屏上的记录（存储在读取间隙发生了变动，请重新载入）：重复行未重复显示。
        </p>
      ) : null}

      {merge.rows.length > 0 ? (
        <div
          className="overflow-y-auto rounded-card border border-line bg-surface"
          style={{ height: HISTORY_VIEWPORT_HEIGHT }}
          onScroll={onScroll}
          data-testid="history-list"
        >
          <div className="relative" style={{ height: merge.rows.length * HISTORY_ROW_HEIGHT }}>
            <div className="absolute inset-x-0" style={{ top: slice.start * HISTORY_ROW_HEIGHT }} role="list" aria-label="历史行">
              {merge.rows.slice(slice.start, slice.end).map((row) => (
                <HistoryRowLine deleteDisabled={false} key={row.id} onDelete={deleteRow} row={row} />
              ))}
            </div>
          </div>
        </div>
      ) : null}

      {/* The bound: 合并 X / 共 Y, straight off the wire's own numbers - a
          window request in flight says so, never a silent pause. */}
      {merge.rows.length > 0 || merge.total > 0 ? (
        <p className="text-[11px] tabular-nums text-ink-3">
          已合并 {merge.rows.length} / 共 {merge.total} 行
          {loading ? "　正在取下一窗…" : ""}
          {!loading && historyHasMore(merge) ? "　（滚到底部继续取）" : ""}
        </p>
      ) : null}

      {/* ------------------- delete / purge / export: the existing privacy face ------------------- */}
      <div className="rounded-card border border-line bg-surface px-3.5 py-3">
        <p className="text-[12px] font-semibold text-ink">删除 / 清空 / 导出（走既有 privacy 面）</p>
        <p className="mt-1 text-[11px] leading-[1.6] text-ink-3">
          这份数据的破坏性动词只有一条路由：privacy.delete-one / privacy.purge / privacy.export
          （internal/memory 的 privacy 层，readface wave 1 落的那扇门），本页不造第二套。下面两枚域词是本读面自己的两张表名
          （task_log＝任务历史，tool_call＝工具取证），门仍按 PrivacyDomains() 判，拒绝时原话逐字上屏。
        </p>
        <div className="mt-2 flex items-center gap-1.5">
          {HISTORY_DOMAINS.map((d) => (
            <button
              type="button"
              key={d}
              aria-pressed={domain === d}
              onClick={() => setDomain(d)}
              className={cn(
                "rounded-chip px-2 py-0.5 font-mono text-[11px] transition-colors duration-150",
                domain === d ? "bg-ink text-canvas" : "bg-inset text-ink-2 hover:bg-hover",
              )}
            >
              {d}
            </button>
          ))}
        </div>
        <div className="mt-2 flex flex-wrap items-center gap-2">
          <Button onClick={() => exp.submit(domain)} size="xs" variant="secondary">
            导出 JSON
          </Button>
          <Button onClick={() => setPurgeConfirmOpen(true)} size="xs" variant="secondary">
            清空这一域
          </Button>
        </div>
        {purgeConfirmOpen ? (
          <div className="mt-2.5 rounded-card border border-red/40 bg-red-tint px-3 py-2.5" role="alertdialog" aria-label="确认清空">
            <p className="text-[12px] font-semibold text-red">确认清空</p>
            <p className="mt-1 text-[11.5px] leading-[1.7] text-ink">
              privacy.purge 不可逆：整域清空后无法恢复。即将清空的域：{domain}
              。其他域不会动；清空的条数以门回执自报的计数为准。
            </p>
            <div className="mt-2 flex items-center gap-2">
              <Button onClick={() => setPurgeConfirmOpen(false)} size="xs" variant="secondary">
                取消
              </Button>
              <Button
                aria-label="确认清空"
                className="bg-red text-white hover:bg-red/90"
                onClick={() => {
                  setPurgeConfirmOpen(false);
                  purge.submit(domain);
                }}
                size="xs"
                variant="primary"
              >
                确认清空
              </Button>
            </div>
          </div>
        ) : null}
        <div className="mt-2 flex flex-col gap-2">
          <PrivacyReplyBlock state={del.state} />
          <PrivacyReplyBlock state={purge.state} />
          <PrivacyReplyBlock state={exp.state} />
        </div>
      </div>
    </div>
  );
}
