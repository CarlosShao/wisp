/* ============================================================================
   L2 strong-confirmation card (ticket 77 AC#3; second-round rebuild 2026-09-25)
   ----------------------------------------------------------------------------
   The owner tore up the demo-faithful route; the visual system is now the
   beautiful-ui component library speaking its own vocabulary (bg-surface /
   border-line / text-ink / shadow-card), with the shadcn Card/Badge/Button
   primitives kept ONLY as the data-slot carriers the render instruments read.
   The shadcn colour names no longer resolve under the fourth-generation token
   table, so every surface/edge/ink class this card needs is spelled out here
   in library vocabulary - the primitives contribute structure, not paint.

   Blueprint for the anatomy: TurboKach's approval-card (the vendored demo at
   src/components/ai-native/approval-card.tsx). The props contract is
   untouched: ApprovalCardView's ten fields, pinned on the Go side by
   TestApprovalCardViewJSONKeysMatchFrontendTypes.

   Two demo ornaments are deliberately NOT drawn, because the snapshot has no
   field behind them and a decoration that implies state is a lie:
     - the three-dot pager: the demo pages through question previews;
       production shows params, reason and call chain at once;
     - the countdown / recommendation blocks: same absence, same rule.

   The pieces the frozen PLAN.md row (section 17.5, the L2 confirmation card
   row) demands stay exactly as the instruments count them:
     - the 2px bar across the top edge, painted with bg-stop for the strong
       levels (the danger token spelling scripts/render-l2.tsx asserts);
     - the title = tool name + level badge, both inside the one card-title
       element, the badge an outline capsule (never a fill);
     - verbatim params, the reason, the C19 rule ids, the call chain;
     - the standing hint with the 20px ball glyph whose ring breathes 2s
       (.l2-ball-ring in theme.css; loop permitted by the PLAN.md:3491 row,
       because this card does not exist while the panel is idle);
     - the two-button footer with no approval-shaped control of any kind.
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
import {
  FilePen,
  FileSearch,
  FolderInput,
  Info,
  PencilLine,
  ShieldAlert,
  ShieldCheck,
  Terminal,
  Trash,
  type LucideIcon,
} from "lucide-react";
import { cn } from "@/lib/cn";
import {
  requestApprovalResolution,
  type ApprovalCardView,
  type ApprovalOutcome,
} from "@/lib/panel";

/** The header icon slot. A table keyed on the real gated tool names, falling
    back to the name's first segment, and to FileSearch past that - a guessed
    glyph per unknown tool would decorate the card with meaning nobody chose. */
const TOOL_ICONS: Record<string, LucideIcon> = {
  "fs.delete": Trash,
  "fs.rename": PencilLine,
  "fs.write": FilePen,
  "fs.mv": FolderInput,
  "shell.exec": Terminal,
};

const TOOL_PREFIX_ICONS: Record<string, LucideIcon> = {
  fs: FileSearch,
  shell: Terminal,
};

function toolIcon(tool: string): LucideIcon {
  const exact = TOOL_ICONS[tool];
  if (exact) return exact;
  return TOOL_PREFIX_ICONS[tool.split(".")[0] ?? ""] ?? FileSearch;
}

/**
 * The level badge: PLAN's outline capsule - a 1px stroke in the level's hue,
 * never a fill (the harness fails the card for the two fill spellings it
 * knows). A level with no twin (Deny) takes the strong stroke and shows its
 * raw string rather than an invented label.
 */
function LevelBadge({ level }: { level: string }) {
  const strong = level === "L2" || level === "Deny";
  if (level === "L1") {
    return (
      <Badge
        className="rounded-full border-accent/40 px-1.5 text-accent-ink"
        variant="outline"
      >
        <ShieldCheck aria-hidden="true" className="size-3" />
        <span>L1 轻确认</span>
      </Badge>
    );
  }
  return (
    <Badge
      className={cn(
        "rounded-full px-1.5",
        strong ? "border-red/40 text-red" : "border-line text-ink-2",
      )}
      variant="outline"
    >
      {strong && <ShieldAlert aria-hidden="true" className="size-3" />}
      <span className={strong ? "uppercase tracking-wide" : ""}>
        {strong ? (level === "L2" ? "L2 强确认" : level) : level}
      </span>
    </Badge>
  );
}

/** 12px semibold uppercase section head, the library's muted ink. */
function SectionTitle({ children }: { children: string }) {
  return (
    <p className="mb-2 text-[12px] font-semibold uppercase tracking-[0.06em] text-ink-3">
      {children}
    </p>
  );
}

/**
 * The argv, verbatim - SPEC-06 §7 / the frozen row want 完整参数, never a
 * summary, so nothing here is elided. Collapsed: one line per element, tool
 * first. Expanded: the same lines with an index and JSON quoting, because
 * join() cannot show an empty argument, a leading space or a trailing
 * newline - and those are exactly the shapes that change what a command
 * deletes. The expand adds precision; it never unlocks withheld text.
 */
function ParamsBlock({
  args,
  expanded,
  tool,
}: {
  args: string[];
  expanded: boolean;
  tool: string;
}) {
  return (
    <div className="overflow-x-auto rounded-control border border-line bg-inset px-3 py-2.5 shadow-hairline">
      {expanded ? (
        <ul className="m-0 flex list-none flex-col gap-0.5 p-0">
          {[tool, ...args].map((element, i) => (
            <li
              className="whitespace-pre font-mono text-[11.5px] leading-[1.45] text-ink"
              key={i}
            >
              <span className="text-ink-3">{String(i).padStart(2, "0")} </span>
              {JSON.stringify(element)}
            </li>
          ))}
        </ul>
      ) : (
        <pre className="m-0 whitespace-pre font-mono text-[11.5px] leading-[1.45] text-ink">
          {[tool, ...args].join("\n")}
        </pre>
      )}
    </div>
  );
}

/** Rules that fired, as C19 ids - the operator reads these against SPEC-06 §3.
    The title repeats the id itself; inventing rule prose here would put words
    in the assessor's mouth. Outline capsule, mono 10.5px, the id in the
    title attribute for the audit hover. */
function RuleList({ rulesHit }: { rulesHit: string[] }) {
  if (rulesHit.length === 0) {
    return (
      <span className="text-[11.5px] text-ink-3">no rule id recorded</span>
    );
  }
  return (
    <ul aria-label="matched risk rules" className="m-0 flex list-none flex-wrap gap-1 p-0">
      {rulesHit.map((rule) => (
        <li key={rule}>
          <Badge
            className="rounded-full border-line px-1.5 text-ink-2"
            style={{
              fontFamily: "var(--font-mono-stack)",
              fontSize: "10.5px",
              height: "20px",
            }}
            title={rule}
            variant="outline"
          >
            {rule}
          </Badge>
        </li>
      ))}
    </ul>
  );
}

/** The reason line, with the fail-closed wording when there is none: an absent
    reason must never render as an empty slot beside a card that looks safe
    (Q-23). Warn colour, not Tailwind amber - the token is the palette here. */
function ReasonLine({ view }: { view: ApprovalCardView }) {
  if (!view.reasonKnown || view.reason.trim() === "") {
    return (
      <p className="text-[12.5px] leading-[1.6] text-ink-2">
        <span className="font-medium text-warn">信息不足：</span>
        风险评估未给出原因文本，本卡不能推断为"无风险"。
      </p>
    );
  }
  return (
    <p className="text-[12.5px] leading-[1.6] text-ink-2">{view.reason}</p>
  );
}

/**
 * The standing row the frozen L2 row puts at the foot of the card: where the
 * allow actually happens, drawn with the 20px ball glyph whose ring breathes
 * on the 2s period the row's motion cell names. Loop animation is allowed
 * here and only here: the non-idle state list names the approval wait, and
 * this card does not exist while the panel is idle.
 */
function BallApproveHint() {
  return (
    <p className="flex items-center gap-1.5 text-[11px] text-ink-2">
      <svg
        aria-hidden="true"
        className="shrink-0 text-ink-2"
        height="20"
        viewBox="0 0 20 20"
        width="20"
      >
        <circle
          className="l2-ball-ring"
          cx="10"
          cy="10"
          fill="none"
          r="8"
          stroke="currentColor"
          strokeWidth="1"
        />
        <circle cx="10" cy="10" fill="currentColor" r="5" />
      </svg>
      点击悬浮球以批准
    </p>
  );
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
  const ToolIcon = toolIcon(view.tool);
  const strong = view.level === "L2" || view.level === "Deny";

  return (
    <Card className="relative w-full gap-0 overflow-hidden rounded-card border-line bg-surface py-0 text-ink shadow-card">
      {/* The frozen row's "visual weight for an irreversible operation": the
          bar across the top edge. bg-stop is the spelling the render harness
          counts; it resolves to the danger hue for L2/Deny, and the brand
          accent sits here only for the lighter levels. Exactly one per card. */}
      <div
        aria-hidden="true"
        className={cn("h-[2px] w-full shrink-0", strong ? "bg-stop" : "bg-accent")}
      />
      <CardHeader className="block px-4 pb-2 pt-3">
        <div className="flex items-start justify-between gap-3">
          <div className="flex min-w-0 items-center gap-2.5">
            <div className="flex size-8 shrink-0 items-center justify-center rounded-control bg-inset text-ink-2 shadow-hairline">
              <ToolIcon aria-hidden="true" className="size-4" />
            </div>
            <div className="min-w-0">
              <p className="text-[13px] font-medium text-ink">需要你的确认</p>
              {/* 标题 = 工具名 + 徽标 (the frozen row); the harness reads both
                  out of this element, so they stay its direct children. */}
              <CardTitle className="mt-0.5 flex min-w-0 items-center gap-2 text-[11.5px] font-normal leading-snug text-ink-2">
                <span className="truncate font-mono">{view.tool}</span>
                <LevelBadge level={view.level} />
              </CardTitle>
            </div>
          </div>
          {/* The demo's pagination-dots slot lives here; hidden on purpose -
              see the file header. */}
        </div>
      </CardHeader>

      <CardContent className="flex flex-col gap-2.5 px-4 pb-3">
        <div>
          <SectionTitle>完整参数</SectionTitle>
          <ParamsBlock args={view.args} expanded={expanded} tool={view.tool} />
        </div>

        <div>
          <SectionTitle>为什么需要确认</SectionTitle>
          <ReasonLine view={view} />
          <div className="mt-2">
            <RuleList rulesHit={view.rulesHit} />
          </div>
        </div>

        {view.callChain.length > 0 && (
          <p className="font-mono text-[11.5px] text-ink-2">
            {`调用链：${view.callChain.join(" > ")}`}
          </p>
        )}
        {view.sessionOverrideBlocked && (
          <p className="text-[11.5px] font-medium text-red">
            本调用带 C25 污染标记，会话级授权对它无效。
          </p>
        )}

        <div className="flex items-start gap-1.5 text-[11px] leading-[1.6] text-ink-3">
          <Info aria-hidden="true" className="mt-px size-3.5 shrink-0" />
          <span>批准权在原生侧 · 面板不提供允许按钮</span>
        </div>
      </CardContent>

      <CardFooter className="flex-col items-stretch gap-2 border-t border-line px-4 py-2.5">
        {/* 常驻一行 (the frozen row): it is there whether or not anything else is. */}
        <BallApproveHint />
        <div className="flex items-center gap-2">
          <span className="mr-auto text-[11px] text-ink-3">
            {`判定来自 ${view.decidedBy === "native" ? "原生风险评估" : view.decidedBy}`}
          </span>
          <Button
            className="rounded-control bg-red text-white hover:bg-red/90"
            onClick={() => send("refuse")}
            size="sm"
            variant="destructive"
          >
            拒绝
          </Button>
          <Button
            className="rounded-control border-line bg-surface text-ink hover:bg-hover"
            onClick={() => setExpanded((e) => !e)}
            size="sm"
            variant="outline"
          >
            {expanded ? "收起参数" : "查看完整参数"}
          </Button>
        </div>
      </CardFooter>
    </Card>
  );
}
