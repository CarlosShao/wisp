/* ============================================================================
   ApprovalQueueSection (ticket 37) - the queue's always-on panel presence.
   ----------------------------------------------------------------------------
   D31 says the queue depth must ALWAYS be visible with the card, and the
   tasks page must show blocked tasks WITH "waiting on what". The ball's
   badge (native, approvalDepthBadge) is the primary channel and the header
   chip is the panel's count, but both of those go dark the moment the user
   is on the tasks page doing something else - and a queue you cannot see is
   a queue you cannot answer.

   What this section draws, all fed verbatim off the pushed snapshot:

     - depth: pending.length. It counts the packet's pending section, which
       nativeVerdictsFromGate copies straight off Queue.LiveApprovals()
       (cmd/wisp/panel_pump.go:78-96) - the queue's own head-first FIFO.
     - 队头: the FIRST card in the same list is what the ball's buttons
       answer (the gate stamps Prompt.Depth from the head's own position, and
       Head() is what the panel face serves, internal/agent/approval/gate.go:747
       /queue.go head()). Drawn from the list, never counted twice.
     - the queue rows with 等待于, sharing one derivation with the approval
       screen (components/approval-screen.tsx waitingOnFor): the join is the
       snapshot's own TaskRosterSection blockedOnApproval, and a snapshot
       without a roster says so instead of inventing a task name.
     - two honest absences, never faked:
         * tasks key absent  -> "任务名册未送达（宿主没带 tasks 这一节）" and
           no task names anywhere;
         * no reader at all  -> the empty sentence names the producer that
           exists (Queue.LiveApprovals -> NativeVerdict.CardView) and states
           that nothing is waiting, which is a MEASURED empty, not a shrug.

   Placement: App mounts this above the page body on every view while
   something is pending - the queue is not a place, it is a condition. With
   nothing pending it renders nothing at all: a permanent zero-depth widget
   is alarm fatigue (the same argument App's header makes for the approval
   chip), and the visibility mandate is about WHILE cards wait.
   ============================================================================ */

import { waitingOnFor } from "@/components/approval-screen";

/** The queue strip the L2 cards sit above. Zero cards = zero markup. */
export function ApprovalQueueSection({ snapshot }: { snapshot: Parameters<typeof waitingOnFor>[0] }) {
  const pending = snapshot.pending;
  if (pending.length === 0) {
    return null;
  }
  return (
    <div
      aria-label="审批队列"
      className="flex flex-col gap-1.5 rounded-card border border-warn-line bg-warn-soft px-3.5 py-2.5"
    >
      <div className="flex items-center gap-2">
        <span className="text-[12px] font-semibold text-warn">审批队列</span>
        <span className="inline-flex h-[18px] items-center rounded-full bg-warn-soft px-1.5 text-[10.5px] font-medium tabular-nums text-warn">
          {`深度 ${pending.length}`}
        </span>
        <span className="text-[10.5px] text-ink-3">悬浮球只答队头</span>
      </div>
      <ul className="m-0 flex list-none flex-col gap-0.5 p-0">
        {pending.map((view, i) => (
          <li
            className="flex items-center gap-2 text-[11.5px]"
            key={view.correlationId}
          >
            <span className="shrink-0 font-mono tabular-nums text-ink-3">
              {i === 0 ? "队头" : String(i + 1).padStart(2, "0")}
            </span>
            <span className="shrink-0 font-mono text-ink">{view.tool}</span>
            <span className="shrink-0 font-mono text-[10.5px] text-ink-3">
              {view.correlationId}
            </span>
            <span className="shrink-0 text-[10.5px] text-ink-3">
              {waitingOnFor(snapshot, view.correlationId)}
            </span>
          </li>
        ))}
      </ul>
    </div>
  );
}
