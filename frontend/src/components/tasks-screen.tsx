/* ============================================================================
   Tasks screen (showcase-only, 2026-09-25)
   ----------------------------------------------------------------------------
   Cut from the library's task-rows + filter-table blueprints: a filter chip
   row over a flat list, one row per task with a six-state capsule and the
   running row's spinner arc. Rows enter staggered (80ms apart, the
   blueprint's own rhythm); the stagger replays once on mount, it does not
   loop (the D32 idle rule).

   STATE RULE: the only state here is which filter chip is active - a view
   preference of the same class as the settings screen's alpha knob: it dies
   with the document and feeds nothing back to the host. The tasks themselves
   are pure props; when Go grows a tasks feed, the snapshot becomes the
   source and this file's props shape is what the view model gets.

   The six states are the panel's honest task vocabulary - NOT invented
   states: 运行中 / 等待确认 (D31 queue) / 被阻塞 (C20 path conflict) /
   排队中 / 失败 (D37 classes) / 已完成.
   ============================================================================ */

import { useState } from "react";
import { cn } from "@/lib/cn";

export type TaskStatus = "running" | "approval" | "blocked" | "queued" | "failed" | "done";

export interface TaskRow {
  id: string;
  title: string;
  sub: string;
  status: TaskStatus;
  /** A finished/running duration label, verbatim; no timer runs here. */
  elapsed?: string;
}

const STATUS_LABELS: Record<TaskStatus, string> = {
  running: "运行中",
  approval: "等待确认",
  blocked: "被阻塞",
  queued: "排队中",
  failed: "失败",
  done: "已完成",
};

/** The six-state colour code: green runs, the brand accent waits for a human,
    amber blocks, red fails, and the two resting states share the neutral
    field. 琥珀是等待、红是错误 - the frozen rows' semantic, on library tint
    classes. */
const STATUS_PILLS: Record<TaskStatus, string> = {
  running: "bg-green-tint text-green",
  approval: "bg-accent-tint text-accent-ink",
  blocked: "bg-orange-tint text-orange",
  queued: "bg-field text-ink-2",
  done: "bg-field text-ink-2",
  failed: "bg-red-tint text-red",
};

/** The running row's spinner arc: a track ring plus a 28%-length arc that
    spins on the `spin` keyframe theme.css carries (loop allowed: a running
    task is a non-idle state). Non-running rows show a plain status dot. */
function StatusGlyph({ status }: { status: TaskStatus }) {
  if (status !== "running") {
    return <span aria-hidden="true" className="size-2 shrink-0 rounded-full bg-line-strong" />;
  }
  const size = 18;
  const stroke = 2;
  const r = (size - stroke) / 2;
  const c = 2 * Math.PI * r;
  return (
    <span className="relative inline-flex size-[18px] shrink-0 items-center justify-center">
      <svg
        aria-hidden="true"
        className="absolute inset-0"
        height={size}
        style={{ animation: "spin 1.1s linear infinite" }}
        viewBox={`0 0 ${size} ${size}`}
        width={size}
      >
        <circle cx={size / 2} cy={size / 2} fill="none" r={r} stroke="var(--line)" strokeWidth={stroke} />
        <circle
          cx={size / 2}
          cy={size / 2}
          fill="none"
          r={r}
          stroke="var(--accent)"
          strokeLinecap="round"
          strokeDasharray={`${c * 0.28} ${c * 0.72}`}
          strokeWidth={stroke}
        />
      </svg>
    </span>
  );
}

export function TasksScreen({ tasks }: { tasks: readonly TaskRow[] }) {
  const [filter, setFilter] = useState<"all" | TaskStatus>("all");
  const counts = new Map<TaskStatus, number>();
  for (const task of tasks) {
    counts.set(task.status, (counts.get(task.status) ?? 0) + 1);
  }
  const shown = filter === "all" ? tasks : tasks.filter((task) => task.status === filter);

  const chips: { key: "all" | TaskStatus; label: string; count: number }[] = [
    { key: "all", label: "全部", count: tasks.length },
    ...(Object.keys(STATUS_LABELS) as TaskStatus[])
      .filter((status) => counts.has(status))
      .map((status) => ({
        key: status,
        label: STATUS_LABELS[status],
        count: counts.get(status) ?? 0,
      })),
  ];

  return (
    <div className="flex w-full flex-col gap-2">
      {/* filter chips */}
      <div className="-mx-1 flex items-center gap-1 overflow-x-auto px-1 py-1" style={{ scrollbarWidth: "none" }}>
        {chips.map((chip) => {
          const active = filter === chip.key;
          return (
            <button
              aria-pressed={active}
              className={cn(
                "flex h-6.5 shrink-0 items-center gap-1.5 rounded-full px-2.5 text-[12px] font-medium transition-[background-color,box-shadow,color] duration-200",
                active ? "bg-accent-tint text-accent-ink" : "text-ink-2 hover:bg-hover",
              )}
              key={chip.key}
              onClick={() => setFilter(chip.key)}
              type="button"
            >
              {chip.label}
              <span
                className={cn(
                  "rounded-[4px] px-1 text-[10.5px] tabular-nums",
                  active ? "bg-field text-ink-2" : "text-ink-3",
                )}
              >
                {chip.count}
              </span>
            </button>
          );
        })}
      </div>

      {/* the list */}
      <div className="overflow-hidden rounded-card border border-line bg-surface shadow-card">
        {shown.length === 0 ? (
          <p className="px-3 py-6 text-center text-[12px] text-ink-3">
            这一档下没有任务。
          </p>
        ) : (
          shown.map((task, i) => (
            <div
              className="flex h-11 items-center gap-2.5 border-b border-line px-3 transition-colors duration-100 last:border-b-0 hover:bg-hover"
              key={task.id}
              style={{ animation: `fade-up 450ms cubic-bezier(0.23,1,0.32,1) ${i * 80}ms both` }}
            >
              <StatusGlyph status={task.status} />
              <span className="min-w-0 flex-1 truncate text-[13px] font-medium text-ink">
                {task.title}
                <span className="ml-2 text-[12px] font-normal text-ink-3">{task.sub}</span>
              </span>
              {task.elapsed && (
                <span className="shrink-0 font-mono text-[11px] tabular-nums text-ink-3">
                  {task.elapsed}
                </span>
              )}
              <span
                className={cn(
                  "inline-flex h-5.5 shrink-0 items-center rounded-full px-2 text-[11.5px] font-medium leading-none",
                  STATUS_PILLS[task.status],
                )}
              >
                {STATUS_LABELS[task.status]}
              </span>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
