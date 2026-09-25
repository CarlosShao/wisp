/* ============================================================================
   Approval screen (third-generation rebuild, 2026-09-25)
   ----------------------------------------------------------------------------
   The 审批中心 view, cut to what PanelSnapshot.pending carries. Layout follows
   the demo's approval screen (design/doubao/demo/screens/approval.js): page
   title and subtitle, the queue depth beside the section title, the fail-closed
   notice, then the queue rows.

   What the demo has that the snapshot does not, and what happens to each:
     - "9 / 10" capacity and 接近容量上限: no capacity field, so the depth badge
       shows the plain count;
     - selected-row state feeding the detail card: the panel has no second
       route, and the demo's selected-card job is done by the L2 card App
       renders globally above this screen - so rows highlight on hover only;
     - batch checkboxes / 批量拒绝: C17's whitelist has no such route, so no
       checkbox and no batch button is drawn;
     - countdowns, ages, risk-rec cards: no fields. Not faked.

   The warn notice uses the C21 warn tokens (var(--warn-line) / var(--warn-soft))
   rather than Tailwind's amber palette - the palette here is the token table.
   ============================================================================ */

import { TriangleAlert } from "lucide-react";
import type { ApprovalCardView, PanelSnapshot } from "@/lib/panel";

/** demo's rowBadge (approval.js:147-153) minus the risk half: the snapshot
    carries the level, not the risk verdict, so the badge shows only what the
    field actually says. */
function RowLevelBadge({ level }: { level: string }) {
  if (level === "L1") {
    return (
      <span className="badge badge-primary" style={{ fontSize: "10px", height: "18px" }}>
        L1
      </span>
    );
  }
  if (level === "L0") {
    return (
      <span className="badge" style={{ fontSize: "10px", height: "18px" }}>
        L0
      </span>
    );
  }
  return (
    <span className="badge badge-destructive" style={{ fontSize: "10px", height: "18px" }}>
      {level}
    </span>
  );
}

function QueueRow({ view }: { view: ApprovalCardView }) {
  return (
    <li className="flex items-center gap-2 rounded-lg border border-transparent px-2.5 py-2 transition-colors hover:border-border hover:bg-muted/50">
      <RowLevelBadge level={view.level} />
      <span className="min-w-0 flex-1 truncate text-[13px] font-medium text-foreground">
        {view.tool}
      </span>
      <span className="shrink-0 font-mono text-[11px] text-muted-foreground">
        {view.correlationId}
      </span>
      <span className="shrink-0 font-mono text-[10.5px] text-muted-foreground">
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
        <h1 className="text-[20px] font-semibold text-foreground">审批中心</h1>
        <p className="mt-1 text-[12px] text-muted-foreground">
          不可逆与高危操作在原生侧批准，面板仅可拒绝与查看完整参数
        </p>
      </div>

      <div className="flex flex-wrap items-center gap-2">
        {/* demo's .section-title (styles.css:501), transcribed to utilities. */}
        <span className="text-[12px] font-semibold uppercase tracking-[0.06em] text-muted-foreground">
          待审批队列
        </span>
        <span className="badge badge-primary">{pending.length}</span>
      </div>

      <div className="flex items-start gap-2 rounded-lg border border-[var(--warn-line)] bg-[var(--warn-soft)] px-3 py-2 text-xs leading-[var(--lh-body)] text-warn">
        <TriangleAlert aria-hidden="true" className="mt-px size-3.5 shrink-0" />
        <span>
          队列容量上限 10 条：满额后新的待确认请求按 fail-closed 自动拒绝，绝不默认放行；宿主不可达或队列出错同样判拒绝。
        </span>
      </div>

      {pending.length === 0 ? (
        <p className="text-[12px] text-muted-foreground">
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
