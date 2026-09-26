/* ============================================================================
   Harness right rail - the context column (?harness=1 demo, 2026-09-26)
   ----------------------------------------------------------------------------
   The right-hand ~280px column of the harness workspace: approval queue,
   running tasks, today's usage, quick actions. Visual vocabulary is the
   beautiful-ui theme verbatim (bg-surface / border-line / text-ink-3 /
   shadow-card), no new paint.

   PROPS CONTRACT - the assembler wires every one of these; this file holds
   no state and no behaviour. Showing/hiding the rail (折叠) is the PARENT's
   conditional render; the component itself always renders at its fixed width.

     pendingCount: number
         The number inside the orange pill on the 审批队列 head.
         0 renders the same pill in the neutral tone.
     queueRows: readonly RightRailQueueRow[]
         Compact queue lines: level badge + tool (mono) + correlationId.
         At most MAX_QUEUE_ROWS (3) render; when the array is longer the
         查看全部 slot appears under them.
     onViewAllApprovals?: () => void
         The 查看全部 slot's callback. The slot renders only when the queue
         is longer than 3 rows; the callback is an affordance, the rail
         routes nothing itself.
     approvalCard?: ApprovalCardView
         When present, the REAL L2ApprovalCard (src/components/
         l2-approval-card.tsx, untouched) renders under the queue lines,
         fed verbatim. Absent = nothing renders, no placeholder lie.
     onApprovalIntent?: (correlationId: string, outcome: ApprovalOutcome) => void
         Passthrough for L2ApprovalCard's onIntent (local UI feedback only;
         the verdict itself is the host's, never this callback's).
     tasks: readonly RightRailTaskRow[]
         Running-task rows: spinner arc + name (mono) + progress 0..1
         rendered as a tabular percentage on the right. Empty array renders
         the honest 「没有进行中的任务」 line.
     usage: RightRailUsage
         tokens / spend as two ValuePills, plus the monthly-budget bar:
         budgetPct >= 100 paints track bg-red-tint with a bg-red fill (the
         117% case), >= 80 fills bg-warn, below that bg-accent on bg-field.
     onNewTask? / onOpenPalette? / onOpenSettings? / onExportLogs?
         The four quick-action icon slots (新任务 / 命令面板 / 设置 /
         导出日志). Callbacks only; the rail wires no behaviour.

   DEMO DATA: RB_RIGHT_RAIL at the bottom of this file is the ?harness=1
   fixture for this component. Spread it: <RightRail {...RB_RIGHT_RAIL} />.
   It is demo-only, reachable through the harness page, never the product
   path (the product panel is fed by the C17 snapshot push).
   ============================================================================ */

import {
  Check,
  Command,
  FileDown,
  LoaderCircle,
  Settings,
  SquarePlus,
  X,
} from "lucide-react";
import { useState } from "react";
import { L2ApprovalCard } from "@/components/l2-approval-card";
import { StatusPill } from "@/components/ai-native/status-pill";
import { ValuePill } from "@/components/ai-native/value-pill";
import { cn } from "@/lib/cn";
import type { ApprovalCardView, ApprovalOutcome } from "@/lib/panel";

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

/** One running task line of the 进行中任务 block. */
export interface RightRailTaskRow {
  id: string;
  /** The task name, set in mono. */
  name: string;
  /** 0..1; rendered as a rounded percentage, tabular. */
  progress: number;
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

/** The 今日用量 numbers. Strings stay verbatim (they are display copy). */
export interface RightRailUsage {
  /** Token count of today, e.g. "12,480". */
  tokens: string;
  /** Spend of today, e.g. "¥0.31". */
  spend: string;
  /** Monthly budget consumption in percent; 117 means over budget. */
  budgetPct: number;
}

export interface RightRailProps {
  pendingCount: number;
  queueRows: readonly RightRailQueueRow[];
  onViewAllApprovals?: () => void;
  approvalCard?: ApprovalCardView;
  onApprovalIntent?: (correlationId: string, outcome: ApprovalOutcome) => void;
  /** Deprecated demo field kept for the assembler's old wiring; the 任务运行
      directory below is runs-driven. */
  tasks?: readonly RightRailTaskRow[];
  runs: readonly RightRailRunRow[];
  usage: RightRailUsage;
  onNewTask?: () => void;
  onOpenPalette?: () => void;
  onOpenSettings?: () => void;
  onExportLogs?: () => void;
}

/** How many queue lines the rail shows before the 查看全部 slot. */
const MAX_QUEUE_ROWS = 3;

/** 11px muted section head - the L2 card's SectionTitle at rail scale. */
function RailSectionHead({
  title,
  trailing,
}: {
  title: string;
  trailing?: React.ReactNode;
}) {
  return (
    <div className="flex items-center justify-between gap-2">
      <p className="text-[11px] font-medium uppercase tracking-[0.06em] text-ink-3">
        {title}
      </p>
      {trailing}
    </div>
  );
}

/** The queue line's level badge: the L2 card's outline-capsule vocabulary
    (1px stroke in the level's hue, never a fill), mono 10.5px. */
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

function QuickActionButton({
  label,
  icon: Icon,
  onClick,
}: {
  label: string;
  icon: React.ComponentType<{ className?: string; "aria-hidden"?: boolean }>;
  onClick?: () => void;
}) {
  return (
    <button
      aria-label={label}
      className="flex size-7 items-center justify-center rounded-control text-ink-3 transition-colors duration-150 hover:bg-hover hover:text-ink"
      onClick={onClick}
      title={label}
      type="button"
    >
      <Icon aria-hidden className="size-4" />
    </button>
  );
}

export function RightRail({
  pendingCount,
  queueRows,
  onViewAllApprovals,
  approvalCard,
  onApprovalIntent,
  runs,
  usage,
  onNewTask,
  onOpenPalette,
  onOpenSettings,
  onExportLogs,
}: RightRailProps) {
  const shown = queueRows.slice(0, MAX_QUEUE_ROWS);
  const overBudget = usage.budgetPct >= 100;
  const nearBudget = usage.budgetPct >= 80;
  const fillCls = overBudget ? "bg-red" : nearBudget ? "bg-warn" : "bg-accent";
  // 任务运行目录：点行展开详情（ZCode 子智能体目录的形态）。
  const [expandedRun, setExpandedRun] = useState<string | null>(null);
  const runningRuns = runs.filter((r) => r.status === "running");
  const finishedRuns = runs.filter((r) => r.status !== "running");

  return (
    <aside
      aria-label="上下文栏"
      className="flex w-[280px] shrink-0 flex-col overflow-y-auto border-l border-line bg-surface text-ink"
    >
      {/* 审批队列 ------------------------------------------------------------ */}
      <section className="flex flex-col gap-2 border-b border-line px-4 py-4" aria-label="审批队列">
        <RailSectionHead
          title="审批队列"
          trailing={
            <StatusPill tone={pendingCount > 0 ? "orange" : "neutral"}>
              待审批 {pendingCount}
            </StatusPill>
          }
        />
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
        {approvalCard && (
          <L2ApprovalCard
            onIntent={onApprovalIntent}
            view={approvalCard}
          />
        )}
      </section>

      {/* 任务运行 ------------------------------------------------------------
          形态对照真实 harness（ZCode 子智能体目录）：运行中/已结束两组、状态
          图标、时长右对齐、失败行红字原因、点行展开详情。 */}
      <section className="flex flex-col gap-2 border-b border-line px-4 py-4" aria-label="任务运行">
        <RailSectionHead
          title="任务运行"
          trailing={
            <span className="text-[11px] tabular-nums text-ink-3">
              运行中 {runningRuns.length} · 已结束 {finishedRuns.length}
            </span>
          }
        />
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
                <LoaderCircle
                  aria-hidden="true"
                  className="size-3.5 shrink-0 text-accent-ink"
                  style={{ animation: "spin 1.1s linear infinite" }}
                />
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
                      {failed ? <X size={9} strokeWidth={3} /> : <Check size={9} strokeWidth={3} />}
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

      {/* 今日用量 ------------------------------------------------------------ */}
      <section className="flex flex-col gap-2 border-b border-line px-4 py-4" aria-label="今日用量">
        <RailSectionHead title="今日用量" />
        <div className="flex flex-wrap items-center gap-x-3 gap-y-1 text-[12px] text-ink-2">
          <span>
            Tokens <ValuePill>{usage.tokens}</ValuePill>
          </span>
          <span>
            花费 <ValuePill>{usage.spend}</ValuePill>
          </span>
        </div>
        <div className="flex items-center gap-2">
          <div
            aria-hidden="true"
            className={cn("h-1 flex-1 overflow-hidden rounded-full", overBudget ? "bg-red-tint" : "bg-field")}
          >
            <div
              className={cn("h-full rounded-full", fillCls)}
              style={{ width: `${Math.min(100, Math.max(0, usage.budgetPct))}%` }}
            />
          </div>
          <span
            className={cn(
              "shrink-0 text-[11px] tabular-nums",
              overBudget ? "text-red" : "text-ink-3",
            )}
          >
            月预算 {Math.round(usage.budgetPct)}%
          </span>
        </div>
      </section>

      {/* 快捷操作 ------------------------------------------------------------ */}
      <section className="flex flex-col gap-2 px-4 py-4" aria-label="快捷操作">
        <RailSectionHead title="快捷操作" />
        <div className="flex items-center gap-1">
          <QuickActionButton icon={SquarePlus} label="新任务" onClick={onNewTask} />
          <QuickActionButton icon={Command} label="命令面板" onClick={onOpenPalette} />
          <QuickActionButton icon={Settings} label="设置" onClick={onOpenSettings} />
          <QuickActionButton icon={FileDown} label="导出日志" onClick={onExportLogs} />
        </div>
      </section>
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
  ],
  approvalCard: {
    correlationId: "harness-0001",
    tool: "fs.delete",
    args: ["--paths", "C:\\Users\\swq\\Desktop\\shot-01.png", "--permanent"],
    level: "L2",
    rulesHit: ["R3", "R8"],
    reason: "该操作将永久删除文件，不进回收站，无法撤销。",
    reasonKnown: true,
    sessionOverrideBlocked: true,
    callChain: ["wisp.run", "agent.turn", "tool.fs.delete"],
    decidedBy: "native",
  },
  tasks: [
    { id: "rb-task-1", name: "归档桌面截图", progress: 0.66 },
    { id: "rb-task-2", name: "下载模型文件", progress: 0.25 },
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
  usage: { tokens: "12,480", spend: "¥0.31", budgetPct: 117 },
};
