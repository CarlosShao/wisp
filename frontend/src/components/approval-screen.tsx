/* ============================================================================
   Approval screen (second-round rebuild, 2026-09-25)
   ----------------------------------------------------------------------------
   The 审批中心 view, cut to what PanelSnapshot.pending carries, on the
   beautiful-ui vocabulary (surface / ink / line / tint) instead of the demo's
   class names.

   What the old demo had that the snapshot does not, and what happens to each:
     - "9 / 10" capacity and 接近容量上限: no capacity field, so the depth
       badge shows the plain count and nothing else;
     - selected-row state feeding a detail card: the panel has no second
       route, and the demo's selected-card job is done by the L2 card App
       renders globally above this screen - so rows highlight on hover only;
     - batch checkboxes / 批量拒绝: C17's whitelist has no such route, so no
       checkbox and no batch button is drawn;
     - countdowns, ages, risk-rec cards: no fields. Not faked.

   The warn notice uses the token-backed utilities (border-warn-line /
   bg-warn-soft / text-warn) rather than Tailwind's amber palette - the
   palette here is the generated token table, where warn rides orange.
   ============================================================================ */

import { TriangleAlert } from "lucide-react";
import { cn } from "@/lib/cn";
import type { ApprovalCardView, PanelSnapshot } from "@/lib/panel";

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

function QueueRow({ view }: { view: ApprovalCardView }) {
  return (
    <li className="flex items-center gap-2 rounded-control px-2.5 py-2 transition-colors duration-100 hover:bg-hover">
      <RowLevelBadge level={view.level} />
      <span className="min-w-0 flex-1 truncate font-mono text-[13px] font-medium text-ink">
        {view.tool}
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
      </div>

      <div className="flex items-start gap-2 rounded-control border border-warn-line bg-warn-soft px-3 py-2 text-xs leading-[1.6] text-warn">
        <TriangleAlert aria-hidden="true" className="mt-px size-3.5 shrink-0" />
        <span>
          队列容量上限 10 条：满额后新的待确认请求按 fail-closed 自动拒绝，绝不默认放行；宿主不可达或队列出错同样判拒绝。
        </span>
      </div>

      {pending.length === 0 ? (
        <p className="text-[12px] leading-[1.6] text-ink-2">
          队列当前为空。待确认请求由原生侧 Go 推送：出现需要人工确认的调用时，会在这里排队并在面板顶部出卡。
        </p>
      ) : (
        <ul className="m-0 flex list-none flex-col gap-1 p-0">
          {pending.map((view) => (
            <QueueRow key={view.correlationId} view={view} />
          ))}
        </ul>
      )}
    </section>
  );
}
