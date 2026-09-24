/* ============================================================================
   L2 strong-confirmation card (ticket 77 AC#3)
   ----------------------------------------------------------------------------
   Adapted from the vendored Beautiful UI ApprovalCard
   (src/components/ai-native/approval-card.tsx, MIT, upstream commit
   05dab2d2b5f1) - its visual grammar is kept: glass card with rounded-card /
   surface / shadow-card, an icon-button dismiss, a ring-dot pager and filled
   circular action buttons. What is NOT kept is upstream's demo questionnaire:
   this card renders the real fields SPEC-06 §7 requires of an L2 prompt, and
   its props type is the Go view model panel.ApprovalCardView, checked against
   internal/risk by TestApprovalCardViewFieldsCoverRiskDecision.

   Q-17 / Q-23 landing: the card shows the complete command and arguments, the
   risk level, the call chain and WHY the call needs a human at all. When the
   assessor gave no reason (RiskReasonUnstated upstream), the card says so in
   words - "unknown reason" must never render as an empty slot, because an
   empty slot next to a green-looking card reads as "no risk".
   ============================================================================ */

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { AlertTriangle, ShieldAlert, X } from "lucide-react";
import { cn } from "@/lib/cn";
import {
  requestApprovalResolution,
  type ApprovalCardView,
  type ApprovalOutcome,
} from "@/lib/panel";

/** D23: icons, never emoji - every glyph here comes from lucide-react. */
function LevelBadge({ level }: { level: string }) {
  const deny = level === "Deny";
  return (
    <Badge
      className={cn(
        "gap-1 border text-[10.5px] font-semibold uppercase",
        deny
          ? "border-stop/40 bg-stop/15 text-stop"
          : "border-caution/40 bg-caution/15 text-caution",
      )}
      variant="outline"
    >
      {deny ? <ShieldAlert size={12} /> : <AlertTriangle size={12} />}
      {level}
    </Badge>
  );
}

/** Rules that fired, as C19 ids - the operator reads these against SPEC-06 §3. */
function RuleList({ rulesHit }: { rulesHit: string[] }) {
  if (rulesHit.length === 0) {
    return <span className="text-ink-3">no rule id recorded</span>;
  }
  return (
    <ul className="flex flex-wrap gap-1" aria-label="matched risk rules">
      {rulesHit.map((rule) => (
        <li key={rule}>
          <Badge className="border-line-strong bg-inset text-ink-2" variant="outline">
            {rule}
          </Badge>
        </li>
      ))}
    </ul>
  );
}

/** The reason line, with the fail-closed wording when there is none. */
function ReasonLine({ view }: { view: ApprovalCardView }) {
  if (!view.reasonKnown || view.reason.trim() === "") {
    return (
      <p className="text-ink-2">
        <span className="font-medium text-caution">信息不足：</span>
        风险评估未给出原因文本，本卡不能推断为"无风险"。
      </p>
    );
  }
  return <p className="text-ink-2">{view.reason}</p>;
}

export function L2ApprovalCard({
  view,
  onIntent,
}: {
  view: ApprovalCardView;
  /** Called after the intent is handed to the host, for local UI feedback
      only. The verdict itself never comes back through this callback. */
  onIntent?: (correlationId: string, outcome: ApprovalOutcome) => void;
}) {
  const send = (outcome: ApprovalOutcome) => {
    // Native side decides (D33/F2); this only reports what was asked for.
    requestApprovalResolution(view.correlationId, outcome);
    onIntent?.(view.correlationId, outcome);
  };

  return (
    <Card className="glass-raised w-full max-w-[var(--panel-w)] gap-0 overflow-hidden rounded-card border-0 py-0">
      <CardHeader className="gap-1 px-4 pt-3.5 pb-2">
        <div className="flex items-start justify-between gap-3">
          <CardTitle className="text-[13px] font-medium text-ink">
            需要你的确认
          </CardTitle>
          <Button
            aria-label="关闭"
            className="-mr-1 size-6 text-ink-3 hover:bg-hover hover:text-ink"
            onClick={() => send("refuse")}
            size="icon"
            variant="ghost"
          >
            <X size={14} />
          </Button>
        </div>
        <div className="flex items-center gap-2">
          <LevelBadge level={view.level} />
          <span className="font-mono text-[11.5px] text-ink-3">{view.tool}</span>
        </div>
      </CardHeader>

      <CardContent className="flex flex-col gap-2.5 px-4 pb-3">
        {/* the command, verbatim - one line per argv element so a quoted
            argument cannot be misread as two */}
        <div className="glass-inset overflow-x-auto rounded-control p-2.5">
          <pre className="font-mono text-[11.5px] leading-[var(--lh-mono)] whitespace-pre text-ink">
            {[view.tool, ...view.args].join("\n")}
          </pre>
        </div>

        <ReasonLine view={view} />
        <RuleList rulesHit={view.rulesHit} />

        {view.callChain.length > 0 && (
          <p className="text-[11.5px] text-ink-3">
            调用链：{view.callChain.join(" > ")}
          </p>
        )}
        {view.sessionOverrideBlocked && (
          <p className="text-[11.5px] font-medium text-stop">
            本调用带 C25 污染标记，会话级授权对它无效。
          </p>
        )}
      </CardContent>

      <CardFooter className="items-center justify-end gap-2 border-t border-line px-4 py-2.5">
        <span className="mr-auto text-[11px] text-ink-3">
          判定来自 {view.decidedBy === "native" ? "原生风险评估" : view.decidedBy}
        </span>
        <Button
          className="rounded-control text-ink-2 hover:bg-hover hover:text-ink"
          onClick={() => send("refuse")}
          size="sm"
          variant="ghost"
        >
          拒绝
        </Button>
      </CardFooter>
    </Card>
  );
}
