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

   Four shapes added 2026-09-25, all named by the frozen row PLAN.md:3481
   ("L2 确认卡（面板内）") and none of them needing a ruling:
     顶部 2px --danger 横条  |  标题 = 工具名 + L2 徽标  |  底部常驻一行「点击悬浮
     球以批准」+ 20px 悬浮球图示（环 2s 呼吸）  |  按钮区补「查看完整参数」（ghost）。
   The row's 明确不用 column names exactly one forbidden thing: a panel-side
   "allow" button (D33/F2, SPEC-06:19, owner's P9 red line #4). It is not here,
   and the hint text deliberately stops at PLAN's own words -- the demo's longer
   line (design/doubao/demo/screens/approval.js:307) claims an F2 hotkey and
   that "panel-side allows are structurally refused server-side", and neither is
   true on disk today (the Go naming test those claims lean on has zero
   definitions, ledger A222#7). Putting them on screen would be a promise the
   native side does not yet keep.
   ============================================================================ */

import { useState } from "react";
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

/** D23: icons, never emoji - every glyph here comes from lucide-react.
    PLAN.md:3478 draws every risk-level badge as a 1px outline capsule with no
    fill; the stricter reading is taken, so the /15 tints that used to sit here
    are gone (a filled amber chip on a card whose whole job is "this may be
    irreversible" reads as a status pill, not as a warning). */
function LevelBadge({ level }: { level: string }) {
  const deny = level === "Deny";
  return (
    <Badge
      className={cn(
        "gap-1 border text-[10.5px] font-semibold uppercase",
        deny ? "border-stop text-stop" : "border-caution text-caution",
      )}
      variant="outline"
    >
      {deny ? <ShieldAlert size={12} /> : <AlertTriangle size={12} />}
      {level}
    </Badge>
  );
}

/**
 * The standing row PLAN.md:3481 puts at the foot of the card: where the allow
 * actually happens, drawn rather than asserted. The circle is the ball; the
 * ring around it breathes on the 2s period that row's 动效 cell names, so the
 * eye is led to it without the card ever growing a button that could approve.
 * Loop animation is allowed here and only here: PLAN.md:3491 lists 审批等待 among
 * the non-idle states, and this card does not exist while the panel is idle.
 */
function BallApproveHint() {
  return (
    <p className="flex items-center gap-1.5 text-[11px] text-ink-3">
      <svg
        aria-hidden="true"
        className="shrink-0 text-ink-2"
        height="20"
        viewBox="0 0 20 20"
        width="20"
      >
        <circle className="l2-ball-ring" cx="10" cy="10" fill="none" r="8" stroke="currentColor" strokeWidth="1" />
        <circle cx="10" cy="10" fill="currentColor" r="5" />
      </svg>
      点击悬浮球以批准
    </p>
  );
}

/**
 * The argv, verbatim. Collapsed it is one line per element (what the card has
 * always shown); expanded it is one line per element WITH its index and its
 * quoting, because join() cannot show an empty argument, a leading space or a
 * trailing newline - and those are exactly the shapes an argument can have that
 * change what a command deletes.
 *
 * SPEC-06 §7 / PLAN.md:3481 want 完整参数 rather than a summary, so nothing here
 * is ever elided: the expand adds precision, it does not unlock withheld text.
 * The control that drives it is a prop, not a local state, because PLAN.md:3481
 * puts that button in the 按钮区 at the foot of the card, away from this block.
 */
function ParamsBlock({ args, expanded, tool }: { args: string[]; expanded: boolean; tool: string }) {
  return (
    <div className="glass-inset overflow-x-auto rounded-control p-2.5">
      {expanded ? (
        <ul className="flex flex-col gap-0.5">
          {[tool, ...args].map((element, i) => (
            <li className="font-mono text-[11.5px] leading-[var(--lh-mono)] whitespace-pre text-ink" key={i}>
              <span className="text-ink-3">{String(i).padStart(2, "0")} </span>
              {JSON.stringify(element)}
            </li>
          ))}
        </ul>
      ) : (
        <pre className="font-mono text-[11.5px] leading-[var(--lh-mono)] whitespace-pre text-ink">
          {[tool, ...args].join("\n")}
        </pre>
      )}
    </div>
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
  const [expanded, setExpanded] = useState(false);
  const send = (outcome: ApprovalOutcome) => {
    // Native side decides (D33/F2); this only reports what was asked for.
    requestApprovalResolution(view.correlationId, outcome);
    onIntent?.(view.correlationId, outcome);
  };

  return (
    <Card className="glass-raised w-full max-w-[var(--panel-w)] gap-0 overflow-hidden rounded-card border-0 py-0">
      {/* PLAN.md:3481's "视觉重量 for an irreversible operation": a 2px danger
          bar across the top edge, drawn from the C21 token, never a literal. */}
      <div aria-hidden="true" className="h-[2px] w-full shrink-0 bg-stop" />
      <CardHeader className="gap-1 px-4 pt-3 pb-2">
        <p className="text-[11px] text-ink-3">需要你的确认</p>
        <div className="flex items-center justify-between gap-3">
          {/* 标题 = 工具名 + L2 徽标 (PLAN.md:3481) */}
          <CardTitle className="flex min-w-0 items-center gap-2 text-[13px] font-medium text-ink">
            <span className="truncate font-mono">{view.tool}</span>
            <LevelBadge level={view.level} />
          </CardTitle>
          <Button
            aria-label="关闭"
            className="-mr-1 shrink-0 size-6 text-ink-3 hover:bg-hover hover:text-ink"
            onClick={() => send("refuse")}
            size="icon"
            variant="ghost"
          >
            <X size={14} />
          </Button>
        </div>
      </CardHeader>

      <CardContent className="flex flex-col gap-2.5 px-4 pb-3">
        {/* the command, verbatim - one line per argv element so a quoted
            argument cannot be misread as two */}
        <ParamsBlock args={view.args} expanded={expanded} tool={view.tool} />

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

      <CardFooter className="flex-col items-stretch gap-2 border-t border-line px-4 py-2.5">
        {/* 常驻一行 (PLAN.md:3481): it is there whether or not anything else is */}
        <BallApproveHint />
        <div className="flex items-center gap-2">
          <span className="mr-auto text-[11px] text-ink-3">
            判定来自 {view.decidedBy === "native" ? "原生风险评估" : view.decidedBy}
          </span>
          <Button
            className="rounded-control text-ink-2 hover:bg-hover hover:text-ink"
            onClick={() => setExpanded((e) => !e)}
            size="sm"
            variant="ghost"
          >
            {expanded ? "收起参数" : "查看完整参数"}
          </Button>
          <Button
            className="rounded-control text-ink"
            onClick={() => send("refuse")}
            size="sm"
            variant="secondary"
          >
            拒绝
          </Button>
        </div>
      </CardFooter>
    </Card>
  );
}
