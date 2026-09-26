/* ============================================================================
   Harness right rail - the TABBED context column (?harness=1 demo).
   ----------------------------------------------------------------------------
   Owner 2026-09-26: the rail must mirror real harness tools, and ZCode's right
   panel is a TABBED container (子智能体目录 + a live canvas/preview tab) - not
   a dump of demo components. So this rail is three tabs:

     任务运行  the run directory (running/finished groups, status icons,
               durations, failure lines, expandable rows) - ZCode's
               子智能体目录 shape.
     审批      the compact approval queue (level tag + tool + correlationId,
               查看全部 slot). The FULL L2 card deliberately does NOT live
               here any more: it is a main-flow component (production renders
               it above the view switch), and stuffing it into the rail was
               the exact thing owner rejected.
     画布      a live canvas tab in ZCode's sense - an embedded preview of the
               component showcase (?harness=2), with refresh and open-in-new
               affordances. ZCode's own right panel does exactly this with its
               "Wisp" tab.

   Props are wired by app-harness.tsx; demo data ships as RB_RIGHT_RAIL.
   Showing/hiding the rail is the PARENT's conditional render.
   ============================================================================ */

import { useState } from "react";
import { Showcase } from "@/components/showcase";
import { StatusPill } from "@/components/ai-native/status-pill";
import { cn } from "@/lib/cn";

/** One compact queue line: what tool, what level, which correlation. */
export interface RightRailQueueRow {
  /** "L0" | "L1" | "L2" | "Deny" - the risk level, painted like the L2
      card's outline capsule (L2/Deny red, L1 accent, rest neutral). */
  level: string;
  /** The gated tool name, set in mono. */
  tool: string;
  /** The C17 correlation id, mono, truncated, full text on the title attr. */
  correlationId: string;
}

/** One run row of the 任务运行 directory (owner 2026-09-26: right rail should
    mirror real harness tools - ZCode's 子智能体目录 shape: status groups,
    durations, failure lines, expandable rows). */
export interface RightRailRunRow {
  id: string;
  name: string;
  status: "running" | "done" | "failed";
  /** Wall-clock label like "2 分" / "38 分" - display copy, not a countdown. */
  duration?: string;
  /** Expanded detail line (a run's outcome summary). */
  summary?: string;
  /** Failed runs carry the reason in red. */
  error?: string;
}

export interface RightRailProps {
  pendingCount: number;
  queueRows: readonly RightRailQueueRow[];
  onViewAllApprovals?: () => void;
  runs: readonly RightRailRunRow[];
}

/** How many queue lines the 审批 tab shows before the 查看全部 slot. */
const MAX_QUEUE_ROWS = 4;

/** The canvas content: the component showcase, mounted DIRECTLY (no iframe).
    The panel CSP is default-src 'none' (ticket 77) - framing even a same-origin
    page is refused, and loosening a security invariant for a demo is not on the
    table. Mounting the component costs nothing: it is already in the bundle. */

type RailTab = "runs" | "approvals" | "canvas";

const TABS: readonly { key: RailTab; label: string }[] = [
  { key: "runs", label: "任务运行" },
  { key: "approvals", label: "审批" },
  { key: "canvas", label: "画布" },
];

function RailTabBar({
  active,
  onSelect,
  pendingCount,
}: {
  active: RailTab;
  onSelect: (tab: RailTab) => void;
  pendingCount: number;
}) {
  return (
    <div
      className="flex shrink-0 items-stretch gap-0.5 border-b border-line bg-inset p-1"
      role="tablist"
      aria-label="右侧栏"
    >
      {TABS.map((tab) => {
        const on = tab.key === active;
        return (
          <button
            aria-selected={on}
            className={cn(
              "flex flex-1 items-center justify-center gap-1.5 rounded-control px-2 py-1.5 text-[12px] font-medium transition-colors duration-150",
              on ? "bg-surface text-ink shadow-btn" : "text-ink-2 hover:bg-hover hover:text-ink",
            )}
            key={tab.key}
            onClick={() => onSelect(tab.key)}
            role="tab"
            type="button"
          >
            {tab.label}
            {tab.key === "approvals" && pendingCount > 0 ? (
              <span className="flex h-4 min-w-4 items-center justify-center rounded-full bg-accent-tint px-1 text-[9.5px] font-semibold tabular-nums text-accent-ink">
                {pendingCount}
              </span>
            ) : null}
          </button>
        );
      })}
    </div>
  );
}

function QueueLevelTag({ level }: { level: string }) {
  const strong = level === "L2" || level === "Deny";
  return (
    <span
      className={cn(
        "inline-flex h-[18px] shrink-0 items-center rounded-full border px-1.5 font-mono text-[10.5px] leading-none",
        strong
          ? "border-red/40 text-red"
          : level === "L1"
            ? "border-accent/40 text-accent-ink"
            : "border-line text-ink-2",
      )}
    >
      {level}
    </span>
  );
}

export function RightRail({
  pendingCount,
  queueRows,
  onViewAllApprovals,
  runs,
}: RightRailProps) {
  const [tab, setTab] = useState<RailTab>("runs");
  const [expandedRun, setExpandedRun] = useState<string | null>(null);
  const runningRuns = runs.filter((r) => r.status === "running");
  const finishedRuns = runs.filter((r) => r.status !== "running");
  const shown = queueRows.slice(0, MAX_QUEUE_ROWS);

  return (
    <aside
      aria-label="上下文栏"
      className={cn(
        "flex shrink-0 flex-col overflow-hidden border-l border-line bg-surface text-ink transition-[width] duration-200",
        tab === "canvas" ? "w-[460px]" : "w-[320px]",
      )}
    >
      <RailTabBar active={tab} onSelect={setTab} pendingCount={pendingCount} />

      {/* 任务运行 ---------------------------------------------------------- */}
      {tab === "runs" && (
        <section aria-label="任务运行" className="flex flex-col gap-2 overflow-y-auto px-3 py-3">
          {runs.length === 0 ? (
            <p className="text-[11.5px] text-ink-3">没有任务运行记录</p>
          ) : (
            <div className="flex flex-col gap-1">
              {runningRuns.map((run) => (
                <button
                  className="flex w-full items-center gap-2 rounded-control px-1.5 py-1.5 text-left transition-colors duration-100 hover:bg-hover"
                  key={run.id}
                  onClick={() => setExpandedRun((current) => (current === run.id ? null : run.id))}
                  type="button"
                >
                  <span
                    aria-hidden="true"
                    className="shrink-0 text-accent-ink"
                    style={{ animation: "spin 1.1s linear infinite" }}
                  >
                    <svg height="14" viewBox="0 0 24 24" width="14" fill="none" stroke="currentColor" strokeLinecap="round" strokeWidth="2.2">
                      <path d="M21 12a9 9 0 1 1-9-9" />
                    </svg>
                  </span>
                  <span className="min-w-0 flex-1 truncate font-mono text-[11.5px] text-ink" title={run.name}>
                    {run.name}
                  </span>
                  <span className="shrink-0 text-[10.5px] tabular-nums text-ink-3">{run.duration}</span>
                </button>
              ))}
              {finishedRuns.map((run) => {
                const failed = run.status === "failed";
                const expanded = expandedRun === run.id;
                return (
                  <button
                    className="flex w-full flex-col items-stretch gap-0.5 rounded-control px-1.5 py-1.5 text-left transition-colors duration-100 hover:bg-hover"
                    key={run.id}
                    onClick={() => setExpandedRun((current) => (current === run.id ? null : run.id))}
                    type="button"
                  >
                    <span className="flex items-center gap-2">
                      <span
                        aria-hidden="true"
                        className={cn(
                          "flex size-3.5 shrink-0 items-center justify-center rounded-full",
                          failed ? "bg-red-tint text-red" : "bg-green-tint text-green",
                        )}
                      >
                        {failed ? (
                          <svg height="9" viewBox="0 0 24 24" width="9" fill="none" stroke="currentColor" strokeLinecap="round" strokeWidth="3">
                            <path d="M18 6L6 18M6 6l12 12" />
                          </svg>
                        ) : (
                          <svg height="9" viewBox="0 0 24 24" width="9" fill="none" stroke="currentColor" strokeLinecap="round" strokeWidth="3">
                            <path d="M20 6L9 17l-5-5" />
                          </svg>
                        )}
                      </span>
                      <span
                        className={cn(
                          "min-w-0 flex-1 truncate font-mono text-[11.5px]",
                          failed ? "text-red" : "text-ink",
                        )}
                        title={run.name}
                      >
                        {run.name}
                      </span>
                      <span className="shrink-0 text-[10.5px] tabular-nums text-ink-3">{run.duration}</span>
                    </span>
                    {failed && run.error ? (
                      <span className="truncate pl-5.5 text-[10.5px] leading-relaxed text-red" title={run.error}>
                        {run.error}
                      </span>
                    ) : null}
                    {expanded && run.summary ? (
                      <span className="pl-5.5 text-[10.5px] leading-relaxed text-ink-3">{run.summary}</span>
                    ) : null}
                  </button>
                );
              })}
            </div>
          )}
        </section>
      )}

      {/* 审批 ---------------------------------------------------------------- */}
      {tab === "approvals" && (
        <section aria-label="审批" className="flex flex-col gap-2 overflow-y-auto px-3 py-3">
          <div className="flex items-center justify-between gap-2">
            <p className="text-[11.5px] text-ink-2">待确认请求按 fail-closed 处理：超时一律判拒绝。</p>
            <StatusPill tone={pendingCount > 0 ? "orange" : "neutral"}>待审批 {pendingCount}</StatusPill>
          </div>
          {shown.length === 0 ? (
            <p className="text-[11.5px] text-ink-3">队列为空</p>
          ) : (
            <ul className="m-0 flex list-none flex-col gap-0.5 p-0">
              {shown.map((row) => (
                <li
                  className="flex items-center gap-2 rounded-control px-1.5 py-1.5 transition-colors duration-100 hover:bg-hover"
                  key={row.correlationId}
                >
                  <QueueLevelTag level={row.level} />
                  <span className="min-w-0 flex-1 truncate font-mono text-[11.5px] text-ink" title={row.tool}>
                    {row.tool}
                  </span>
                  <span className="shrink-0 font-mono text-[10.5px] text-ink-3" title={row.correlationId}>
                    {row.correlationId}
                  </span>
                </li>
              ))}
            </ul>
          )}
          {queueRows.length > MAX_QUEUE_ROWS && (
            <button
              className="self-start rounded-full px-1.5 py-0.5 text-[11px] text-accent-ink transition-colors duration-100 hover:bg-hover hover:underline"
              onClick={onViewAllApprovals}
              type="button"
            >
              查看全部
            </button>
          )}
          <p className="text-[10.5px] leading-relaxed text-ink-3">
            完整的 L2 确认卡在组件陈列室（画布标签可预览）与生产面板的会话流里，这里只放队列。
          </p>
        </section>
      )}

      {/* 画布 ------------------------------------------------------------------
          ZCode 的右侧栏有画布/预览标签（其「Wisp」标签即内嵌实时页面）。这里
          把组件陈列室直接挂载进画布（不 iframe：面板 CSP 是 default-src 'none'，
          安全底线不为演示松口），标签激活时右栏自动加宽。 */}
      {tab === "canvas" && (
        <section aria-label="画布" className="flex min-h-0 flex-1 flex-col overflow-y-auto">
          <Showcase />
        </section>
      )}
    </aside>
  );
}

/* ============================================================================
   RB_RIGHT_RAIL - the harness fixture for <RightRail />, defined here (not in
   src/fixtures/harness.ts) so the rail's props shape and its demo data ship
   together. Demo only: the product path reads the C17 snapshot, never this.
   ============================================================================ */

export const RB_RIGHT_RAIL: RightRailProps = {
  pendingCount: 3,
  queueRows: [
    { level: "L2", tool: "fs.delete", correlationId: "harness-0001" },
    { level: "L2", tool: "shell.exec", correlationId: "harness-0004" },
    { level: "L1", tool: "fs.write", correlationId: "harness-0007" },
    { level: "L0", tool: "fs.listdir", correlationId: "harness-0009" },
    { level: "L2", tool: "web.fetch", correlationId: "harness-0012" },
  ],
  runs: [
    { id: "rb-run-1", name: "归档桌面截图", status: "running" as const, duration: "2 分" },
    { id: "rb-run-2", name: "下载模型文件", status: "running" as const, duration: "1 分" },
    {
      id: "rb-run-3",
      name: "会议纪要待办",
      status: "done" as const,
      duration: "5 分",
      summary: "抽出待办 3 项，已按人汇总写入台账",
    },
    { id: "rb-run-4", name: "周报生成", status: "done" as const, duration: "12 分", summary: "周报已落 Desktop\\周报" },
    {
      id: "rb-run-5",
      name: "磁盘清理脚本",
      status: "failed" as const,
      duration: "3 分",
      error: "R6 命中：argv 含管道与重定向，需逐条确认",
      summary: "已改为移入「待确认」目录",
    },
  ],
};
