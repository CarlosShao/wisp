/* ============================================================================
   Harness right rail - the TABBED context column (?harness=1 demo).
   ----------------------------------------------------------------------------
   Owner 2026-09-26 (second pass, with ZCode screenshots annotated): the tab
   container must carry the WHOLE ZCode interaction model, not just tabs -

     - every tab has a hover × that CLOSES it;
     - a "+" control opens the most recently closed tab;
     - the ⌄ control opens a DROPDOWN listing 打开的标签页 (with per-tab ×)
       and 最近关闭的标签页 (click to reopen), with a 搜索标签页 filter;
     - closing every tab shows the Qoder-style empty state (big action cards
       over a radial glow) whose cards reopen the tabs.

   Three tabs: 任务运行 (run directory), 审批 (compact queue), 画布 (live
   mounted showcase; widens the rail). The FULL L2 card deliberately does not
   live here - it is a main-flow component (owner rejected it in the rail).
   ============================================================================ */

import { useState } from "react";
import {
  Activity,
  PanelRight,
  Plus,
  ShieldCheck,
  X,
} from "lucide-react";
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

/** One run row of the 任务运行 directory (ZCode 子智能体目录 shape). */
export interface RightRailRunRow {
  id: string;
  name: string;
  status: "running" | "done" | "failed";
  /** Wall-clock label like "2 分" - display copy, not a countdown. */
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

const MAX_QUEUE_ROWS = 4;

type RailTabKey = "runs" | "approvals" | "canvas";

interface RailTabDef {
  key: RailTabKey;
  label: string;
  icon: typeof Activity;
  /** Closed-time label for the 最近关闭 list (demo copy). */
  closedAgo: string;
}

const TAB_DEFS: readonly RailTabDef[] = [
  { key: "runs", label: "任务运行", icon: Activity, closedAgo: "刚刚" },
  { key: "approvals", label: "审批", icon: ShieldCheck, closedAgo: "3 小时" },
  { key: "canvas", label: "画布", icon: PanelRight, closedAgo: "5 小时" },
];

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
  // ZCode 的标签容器模型：打开集合 + 激活 + 最近关闭栈。
  const [openTabs, setOpenTabs] = useState<RailTabKey[]>(["runs", "approvals", "canvas"]);
  const [activeTab, setActiveTab] = useState<RailTabKey | null>("runs");
  const [closedTabs, setClosedTabs] = useState<RailTabKey[]>([]);
  const [menuOpen, setMenuOpen] = useState(false);
  const [menuQuery, setMenuQuery] = useState("");
  const [expandedRun, setExpandedRun] = useState<string | null>(null);
  const runningRuns = runs.filter((r) => r.status === "running");
  const finishedRuns = runs.filter((r) => r.status !== "running");
  const shown = queueRows.slice(0, MAX_QUEUE_ROWS);

  function closeTab(key: RailTabKey) {
    setOpenTabs((current) => {
      const next = current.filter((k) => k !== key);
      if (activeTab === key) {
        setActiveTab(next[next.length - 1] ?? null);
      }
      return next;
    });
    setClosedTabs((current) => [key, ...current.filter((k) => k !== key)].slice(0, 3));
    setMenuOpen(false);
  }

  function reopenTab(key: RailTabKey) {
    setOpenTabs((current) => (current.includes(key) ? current : [...current, key]));
    setClosedTabs((current) => current.filter((k) => k !== key));
    setActiveTab(key);
    setMenuOpen(false);
  }

  const active = activeTab && openTabs.includes(activeTab) ? activeTab : null;
  const onCanvas = active === "canvas";
  const menuMatches = (def: RailTabDef) =>
    def.label.toLowerCase().includes(menuQuery.toLowerCase());

  return (
    <aside
      aria-label="上下文栏"
      className={cn(
        "flex shrink-0 flex-col overflow-hidden border-l border-line bg-surface text-ink transition-[width] duration-200",
        onCanvas ? "w-[460px]" : "w-[320px]",
      )}
    >
      {/* ── 标签条：⌄ 下拉 + 可关标签 + ＋ 重开 ─────────────────────────── */}
      <div className="relative flex shrink-0 items-stretch gap-0.5 border-b border-line bg-inset p-1">
        <button
          aria-expanded={menuOpen}
          aria-label="标签页列表"
          className={cn(
            "flex w-7 shrink-0 items-center justify-center rounded-control text-ink-2 transition-colors duration-150 hover:bg-hover hover:text-ink",
            menuOpen && "bg-hover text-ink",
          )}
          onClick={() => setMenuOpen((v) => !v)}
          title="标签页列表"
          type="button"
        >
          <ChevronGlyph />
        </button>
        {openTabs.map((key) => {
          const def = TAB_DEFS.find((t) => t.key === key);
          if (!def) return null;
          const on = key === active;
          return (
            <button
              aria-selected={on}
              className={cn(
                "group flex min-w-0 flex-1 items-center justify-center gap-1 rounded-control px-2 py-1.5 text-[12px] font-medium transition-colors duration-150",
                on ? "bg-surface text-ink shadow-btn" : "text-ink-2 hover:bg-hover hover:text-ink",
              )}
              key={key}
              onClick={() => setActiveTab(key)}
              role="tab"
              type="button"
            >
              <span className="truncate">{def.label}</span>
              {key === "approvals" && pendingCount > 0 ? (
                <span className="flex h-4 min-w-4 shrink-0 items-center justify-center rounded-full bg-accent-tint px-1 text-[9.5px] font-semibold tabular-nums text-accent-ink">
                  {pendingCount}
                </span>
              ) : null}
              <span
                aria-label={`关闭 ${def.label}`}
                className={cn(
                  "flex size-3.5 shrink-0 items-center justify-center rounded-[4px] text-ink-3 opacity-0 transition-[opacity,background-color,color] duration-100 hover:bg-line/70 hover:text-ink group-hover:opacity-100",
                  on && "opacity-100",
                )}
                onClick={(event) => {
                  event.stopPropagation();
                  closeTab(key);
                }}
                role="button"
                tabIndex={-1}
              >
                <X aria-hidden="true" size={9} strokeWidth={2.6} />
              </span>
            </button>
          );
        })}
        <button
          aria-label="打开最近关闭的标签页"
          className="flex w-7 shrink-0 items-center justify-center rounded-control text-ink-2 transition-colors duration-150 hover:bg-hover hover:text-ink disabled:pointer-events-none disabled:opacity-40"
          disabled={closedTabs.length === 0}
          onClick={() => reopenTab(closedTabs[0])}
          title={closedTabs.length ? `打开 ${TAB_DEFS.find((t) => t.key === closedTabs[0])?.label}` : "没有最近关闭的标签页"}
          type="button"
        >
          <Plus aria-hidden="true" size={13} strokeWidth={2.2} />
        </button>

        {/* ── 下拉：搜索 + 打开的标签页 + 最近关闭的标签页 ─────────────── */}
        {menuOpen && (
          <div
            className="absolute left-1 right-1 top-full z-40 mt-1 rounded-card border border-line bg-surface p-1.5 shadow-overlay"
            style={{ animation: "pop-in 160ms var(--ease-out-strong) both" }}
          >
            <label className="mb-1 flex h-8 items-center gap-2 rounded-control bg-inset px-2.5 shadow-hairline">
              <input
                onChange={(event) => setMenuQuery(event.target.value)}
                placeholder="搜索标签页…"
                value={menuQuery}
                className="min-w-0 flex-1 bg-transparent text-[12.5px] text-ink outline-none placeholder:text-ink-3"
              />
            </label>
            <p className="px-1.5 pb-0.5 pt-1.5 text-[10.5px] font-medium uppercase tracking-[0.08em] text-ink-3">
              打开的标签页
            </p>
            {openTabs.filter((k) => menuMatches(TAB_DEFS.find((t) => t.key === k)!)).map((key) => {
              const def = TAB_DEFS.find((t) => t.key === key)!;
              return (
                <div className="flex items-center gap-2 rounded-chip px-1.5 py-1.5 hover:bg-hover" key={key}>
                  <def.icon aria-hidden="true" className="shrink-0 text-ink-2" size={13} strokeWidth={1.8} />
                  <button
                    className="min-w-0 flex-1 truncate text-left text-[12px] text-ink"
                    onClick={() => {
                      setActiveTab(key);
                      setMenuOpen(false);
                    }}
                    type="button"
                  >
                    {def.label}
                  </button>
                  <span
                    aria-label={`关闭 ${def.label}`}
                    className="flex size-4 shrink-0 items-center justify-center rounded-[4px] text-ink-3 hover:bg-line/70 hover:text-ink"
                    onClick={() => closeTab(key)}
                    role="button"
                    tabIndex={-1}
                  >
                    <X aria-hidden="true" size={9} strokeWidth={2.6} />
                  </span>
                </div>
              );
            })}
            {closedTabs.filter((k) => menuMatches(TAB_DEFS.find((t) => t.key === k)!)).length > 0 && (
              <p className="px-1.5 pb-0.5 pt-2 text-[10.5px] font-medium uppercase tracking-[0.08em] text-ink-3">
                最近关闭的标签页
              </p>
            )}
            {closedTabs.filter((k) => menuMatches(TAB_DEFS.find((t) => t.key === k)!)).map((key) => {
              const def = TAB_DEFS.find((t) => t.key === key)!;
              return (
                <button
                  className="flex w-full items-center gap-2 rounded-chip px-1.5 py-1.5 text-left hover:bg-hover"
                  key={key}
                  onClick={() => reopenTab(key)}
                  type="button"
                >
                  <def.icon aria-hidden="true" className="shrink-0 text-ink-3" size={13} strokeWidth={1.8} />
                  <span className="min-w-0 flex-1 truncate text-[12px] text-ink-2">{def.label}</span>
                  <span className="shrink-0 text-[10.5px] text-ink-3">{def.closedAgo}</span>
                </button>
              );
            })}
          </div>
        )}
      </div>

      {/* ── 全关后的空状态（Qoder 大卡片组） ───────────────────────────── */}
      {openTabs.length === 0 ? (
        <div
          className="flex min-h-0 flex-1 flex-col items-center justify-center gap-3 px-4 py-6"
          style={{
            backgroundImage:
              "radial-gradient(circle at 50% 55%, var(--accent-tint) 0%, transparent 62%)",
          }}
        >
          <p className="text-[13px] font-semibold text-ink">打开面板</p>
          <p className="mb-1 text-[11.5px] text-ink-3">选择要在侧边面板中查看的内容。</p>
          <div className="flex w-full flex-col gap-2.5">
            {TAB_DEFS.map((def) => (
              <button
                className="flex w-full items-center gap-3 rounded-card border border-line bg-surface p-4 text-left shadow-card transition-colors duration-150 hover:border-accent/40"
                key={def.key}
                onClick={() => reopenTab(def.key)}
                type="button"
              >
                <span className="flex size-8 shrink-0 items-center justify-center rounded-control bg-field text-ink-2">
                  <def.icon aria-hidden="true" size={15} strokeWidth={1.8} />
                </span>
                <span className="min-w-0 flex-1 truncate text-[12.5px] font-medium text-ink">
                  {def.label}
                </span>
              </button>
            ))}
          </div>
        </div>
      ) : (
        <>
          {/* ── 任务运行 ---------------------------------------------------- */}
          {active === "runs" && (
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
                        <svg fill="none" height="14" stroke="currentColor" strokeLinecap="round" strokeWidth="2.2" viewBox="0 0 24 24" width="14">
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
                              <svg fill="none" height="9" stroke="currentColor" strokeLinecap="round" strokeWidth="3" viewBox="0 0 24 24" width="9">
                                <path d="M18 6L6 18M6 6l12 12" />
                              </svg>
                            ) : (
                              <svg fill="none" height="9" stroke="currentColor" strokeLinecap="round" strokeWidth="3" viewBox="0 0 24 24" width="9">
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

          {/* ── 审批 ---------------------------------------------------------- */}
          {active === "approvals" && (
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

          {/* ── 画布 ----------------------------------------------------------
              ZCode 的画布/预览标签。面板 CSP 是 default-src 'none'（票 77），
              不走 iframe——直接挂载陈列室，标签激活时右栏加宽。 */}
          {active === "canvas" && (
            <section aria-label="画布" className="flex min-h-0 flex-1 flex-col overflow-y-auto">
              <Showcase />
            </section>
          )}
        </>
      )}
    </aside>
  );
}

/* ── 内联小图标（chevron 上下双箭头，ZCode 下拉钮同款） ─────────────────── */

function ChevronGlyph() {
  return (
    <svg
      fill="none"
      height="12"
      stroke="currentColor"
      strokeLinecap="round"
      strokeLinejoin="round"
      strokeWidth="2"
      viewBox="0 0 24 24"
      width="12"
      aria-hidden="true"
    >
      <path d="M7 15l5 5 5-5M7 9l5-5 5 5" />
    </svg>
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
