/* ============================================================================
   Approval screen - the 审批中心 view of the LIVE approval queue.
   ----------------------------------------------------------------------------
   Rebuilt 2026-11-16 for ticket 37. The old cut was a STATIC demo listing
   (GlideMenu rows of whatever pending carried) plus a capacity sentence that
   hard-coded a number the wire never sent. This cut draws what the snapshot
   actually carries, with every value fed in displayed verbatim:

     - FIFO order: rows are numbered 1..n in wire order (LiveApprovals walks
       q.pending head-first, internal/agent/approval/pending_read.go:162-180),
       and the first row carries the 队头 mark - what the user answers through
       the ball is the HEAD, not the list;
     - the honest capacity sentence: the queue's own bound is
       DefaultMaxPending = 8 (internal/agent/approval/queue.go:128-130), and
       the fail-closed behaviour on overflow is the queue's own sentence;
     - 等待于 per row: the task dimension a row is waiting on. The wire's
       ApprovalCardView carries correlationId but no taskId (the panel shape
       is the six-fields-plus whitelist of internal/panel/approval.go:38-59),
       so this column derives from TaskRosterSection's blockedOnApproval join
       (internal/panel/subagent_roster_197.go:186-247) when the snapshot has a
       roster: the row whose task's card is THIS correlationId names the task
       it waits on, verbatim. Without a roster there is no reader for the
       task id, and the row says so instead of guessing;
     - what moved out of this file: the batch view and the 50+ escalation
       bench. The D45-1 aggregation (internal/agent/approval/batch.go:35-74)
       is produced inside the gate and consumed by the native card face
       (internal/agent/approval/cardrender.go:84 renderBatch); the panel wire
       (ApprovalCardView) carries no BatchView, and a summary drawn here would
       be a value no producer sent. The panel never aggregates L2 and the
       gate refuses to - R7's 50+ batch escalates to L2 at
       internal/agent/approval/gate.go:317-324 and stays unaggregated;
     - timeout UX: the 270s warning / 300s auto-reject ladder is the NATIVE
       card's job (internal/agent/approval/gate.go deadline + EventWarning,
       consumed by the ball's card UI), not this page's: the wire carries no
       deadline field, and drawing one here would invent a number.

   The level capsule keeps the old rule: the snapshot carries the level the
   queue admitted, the badge shows only what the field says.
   ============================================================================ */

import { TriangleAlert } from "lucide-react";
import { cn } from "@/lib/cn";
import type { ApprovalCardView, PanelSnapshot } from "@/lib/panel";

/** The queue's capacity, as the queue itself fixes it:
    DefaultMaxPending = 8 (internal/agent/approval/queue.go:128-130). A named
    constant here rather than a number buried in prose, so the render
    instrument can pin sentence == constant and the sentence cannot drift
    from the bound it describes. */
export const QUEUE_CAPACITY = 8;

/** The queue row's level capsule: strong levels ride red tint, L1 the brand
    accent, L0 the neutral field. The snapshot carries the level, not the risk
    verdict, so the badge shows only what the field actually says. */
function RowLevelBadge({ level }: { level: string }) {
  const strong = level === "L2" || level === "Deny";
  return (
    <span
      className={cn(
        "inline-flex h-[18px] shrink-0 items-center rounded-full px-1.5 text-[10px] font-medium leading-none",
        strong && "bg-red-tint text-red",
        level === "L1" && "bg-accent-tint text-accent-ink",
        (level === "L0" || (!strong && level !== "L1")) && "bg-field text-ink-2",
      )}
    >
      {level}
    </span>
  );
}

/** One FIFO row: index + level + tool + correlationId + rules count +
    等待于. Every value on the row is either verbatim off the wire or the
    row's own position - nothing is derived twice. */
function QueueRow({
  view,
  position,
  waitingOn,
}: {
  view: ApprovalCardView;
  /** 1-based place in the FIFO, in wire order. */
  position: number;
  /** The 等待于 sentence: the task name when a roster row's card is this
      corr, the host's own missing-reader sentence otherwise. */
  waitingOn: string;
}) {
  return (
    <li
      aria-label={`审批队列第 ${position} 项`}
      className="flex items-center gap-2 rounded-control px-2.5 py-2"
      data-menu-row
    >
      <span
        className="shrink-0 font-mono text-[11px] tabular-nums text-ink-3"
        title="FIFO 顺序（LiveApprovals 队头在先）"
      >
        {String(position).padStart(2, "0")}
      </span>
      <RowLevelBadge level={view.level} />
      <span className="min-w-0 flex-1 truncate font-mono text-[13px] font-medium text-ink">
        {view.tool}
      </span>
      <span
        className="shrink-0 text-[10.5px] text-ink-3"
        title="这一项正在等谁（读自快照的任务名册 join；名册不在时报缺）"
      >
        {waitingOn}
      </span>
      <span className="shrink-0 font-mono text-[11px] text-ink-3">
        {view.correlationId}
      </span>
      <span className="shrink-0 text-[10.5px] text-ink-3">
        {view.rulesHit.length} 条规则
      </span>
    </li>
  );
}

/** The task a queued card waits on, read off the SAME snapshot's task roster:
    a roster row with blockedOnApproval:true whose join covers this
    correlationId names its own task. The Go join (indexWaitingKey) matches
    both id families - the corr as the queue filed it and the task id a corr
    was built on - so a row here means the HOST matched it, not this page.
    No tasks key on the snapshot = no reader for the task dimension at all,
    and the honest sentence says exactly that instead of a guess. Shared with
    ApprovalQueueSection (components/approval-queue.tsx), which renders the
    same sentences on every view. */
export function waitingOnFor(snapshot: PanelSnapshot, corr: string): string {
  const roster = snapshot.tasks;
  if (roster === undefined) {
    return "等待于：快照未带任务名册";
  }
  const match = roster.rows.find((row) => row.blockedOnApproval && row.taskId === corr) ??
    roster.rows.find((row) => row.blockedOnApproval && corr.startsWith(`${row.taskId}#`));
  if (match) {
    const label = match.label.trim() !== "" ? match.label : match.taskId;
    return `等待于：${label}`;
  }
  return "等待于：名册里没有认领这张卡的任务";
}

export function ApprovalScreen({ snapshot }: { snapshot: PanelSnapshot }) {
  const pending = snapshot.pending;
  return (
    <section aria-label="审批中心" className="flex min-w-0 flex-col gap-3">
      <div>
        <h1 className="text-[15px] font-semibold text-ink">审批中心</h1>
        <p className="mt-1 text-[12.5px] text-ink-2">
          不可逆与高危操作在原生侧批准，面板只递请求、不做决定。
        </p>
      </div>

      <div className="flex flex-wrap items-center gap-2">
        <span className="text-[12px] font-semibold uppercase tracking-[0.06em] text-ink-3">
          待审批队列
        </span>
        <span className="inline-flex h-6 items-center rounded-full bg-accent-tint px-2.5 text-[11px] font-medium text-accent-ink">
          {pending.length}
        </span>
        {pending.length > 0 ? (
          <span className="text-[11px] text-ink-3">队头在先（FIFO），悬浮球只答队头</span>
        ) : null}
      </div>

      <div className="flex items-start gap-2 rounded-control border border-warn-line bg-warn-soft px-3 py-2 text-xs leading-[1.6] text-warn">
        <TriangleAlert aria-hidden="true" className="mt-px size-3.5 shrink-0" />
        <span>
          队列容量上限 {QUEUE_CAPACITY} 条：满额后新的待确认请求按 fail-closed 自动拒绝，绝不默认放行；宿主不可达或队列出错同样判拒绝。
        </span>
      </div>

      {pending.length === 0 ? (
        <p className="text-[12px] leading-[1.6] text-ink-2">
          队列当前为空。待确认请求由原生侧 Go 推送：出现需要人工确认的调用时，会在这里排队并在面板顶部出卡。
        </p>
      ) : (
        <ul className="m-0 flex list-none flex-col gap-1 p-0">
          {pending.map((view, i) => (
            <QueueRow
              key={view.correlationId}
              position={i + 1}
              view={view}
              waitingOn={waitingOnFor(snapshot, view.correlationId)}
            />
          ))}
        </ul>
      )}
    </section>
  );
}
