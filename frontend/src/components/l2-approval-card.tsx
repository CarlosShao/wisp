/* ============================================================================
   L2 strong-confirmation card (ticket 77 AC#3; demo anatomy pass 2026-09-25)
   ----------------------------------------------------------------------------
   Re-cut onto the high-fidelity demo's approval-card anatomy
   (design/doubao/demo/screens/approval.js renderDetail): a card-spotlight
   surface with a level bar across the top edge, an icon slot beside the
   "需要你的确认" head, a mono tool name with the level badge in the title,
   verbatim params in a muted inset, the reason and the C19 rule ids that
   fired, the call chain, and a footer that keeps the only two buttons
   PLAN.md:3481 allows. The props contract is untouched: ApprovalCardView's
   ten fields, pinned on the Go side by TestApprovalCardViewJSONKeysMatchFrontendTypes.

   Two demo ornaments are deliberately NOT drawn, because the snapshot has no
   field behind them and a decoration that implies state is a lie:
     - the three-dot pager: demo pages through three question previews;
       production shows params, reason and call chain at once and has no
       page field (ticket 145's census decides if one ever lands);
     - the countdown / recommendation blocks: same absence, same rule.

   The pieces theme.css freezes from PLAN.md:3481 stay exactly as they were:
   the top bar, the level badge in the title, the standing hint with the 20px
   ball glyph, and the two-button footer. scripts/render-l2.tsx counts all of
   them, so their per-card multiplicities are part of the contract.

   One spelling note: the top bar paints with bg-stop, not bg-destructive.
   Both resolve to hsl(var(--destructive)); --stop is the C21 danger token the
   frozen row names, and the render harness asserts that spelling.
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

/** Demo headerBadge (approval.js:133-144): L1 rides the primary capsule, the
    strong levels the destructive one. A level with no demo twin (Deny) shows
    its raw string rather than an invented label. The demo look comes from
    theme.css's unlayered .badge classes; the 10.5px size rides the same inline
    style the demo uses, because unlayered CSS beats Tailwind's utilities. */
function LevelBadge({ level }: { level: string }) {
  if (level === "L1") {
    return (
      <Badge className="badge badge-primary" style={{ fontSize: "10.5px" }}>
        <ShieldCheck aria-hidden="true" />
        <span>L1 轻确认</span>
      </Badge>
    );
  }
  if (level === "L2") {
    return (
      <Badge className="badge badge-destructive" style={{ fontSize: "10.5px" }}>
        <ShieldAlert aria-hidden="true" />
        <span className="uppercase tracking-wide">L2 强确认</span>
      </Badge>
    );
  }
  return (
    <Badge
      className="badge badge-destructive font-mono uppercase"
      style={{ fontSize: "10.5px" }}
    >
      {level}
    </Badge>
  );
}

/** demo's .section-title (styles.css:501), transcribed to utilities: 12px,
    semibold, uppercase, 0.06em tracking, muted. */
function SectionTitle({ children }: { children: string }) {
  return (
    <p className="mb-2 text-[12px] font-semibold uppercase tracking-[0.06em] text-muted-foreground">
      {children}
    </p>
  );
}

/**
 * The argv, verbatim - SPEC-06 §7 / PLAN.md:3481 want 完整参数, never a
 * summary, so nothing here is elided. Collapsed: one line per element, tool
 * first, the demo's own pre shape. Expanded: the same lines with an index and
 * JSON quoting, because join() cannot show an empty argument, a leading space
 * or a trailing newline - and those are exactly the shapes that change what a
 * command deletes. The expand adds precision; it never unlocks withheld text.
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
    <div className="overflow-x-auto rounded-lg border border-border bg-muted/40 px-3 py-2.5">
      {expanded ? (
        <ul className="m-0 flex list-none flex-col gap-0.5 p-0">
          {[tool, ...args].map((element, i) => (
            <li
              className="whitespace-pre font-mono text-[11.5px] leading-[var(--lh-mono)] text-foreground"
              key={i}
            >
              <span className="text-muted-foreground">
                {String(i).padStart(2, "0")}{" "}
              </span>
              {JSON.stringify(element)}
            </li>
          ))}
        </ul>
      ) : (
        <pre className="m-0 whitespace-pre font-mono text-[11.5px] leading-[var(--lh-mono)] text-foreground">
          {[tool, ...args].join("\n")}
        </pre>
      )}
    </div>
  );
}

/** Rules that fired, as C19 ids - the operator reads these against SPEC-06 §3.
    The title repeats the id itself; inventing rule prose here would put words
    in the assessor's mouth. */
function RuleList({ rulesHit }: { rulesHit: string[] }) {
  if (rulesHit.length === 0) {
    return (
      <span className="text-[11.5px] text-muted-foreground">
        no rule id recorded
      </span>
    );
  }
  return (
    <ul aria-label="matched risk rules" className="m-0 flex list-none flex-wrap gap-1 p-0">
      {rulesHit.map((rule) => (
        <li key={rule}>
          <Badge
            className="badge"
            title={rule}
            variant="outline"
            style={{
              fontFamily: "var(--font-mono-stack)",
              fontSize: "10.5px",
              height: "20px",
            }}
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
      <p className="text-[12.5px] leading-[var(--lh-body)] text-muted-foreground">
        <span className="font-medium text-warn">信息不足：</span>
        风险评估未给出原因文本，本卡不能推断为"无风险"。
      </p>
    );
  }
  return (
    <p className="text-[12.5px] leading-[var(--lh-body)] text-muted-foreground">
      {view.reason}
    </p>
  );
}

/**
 * The standing row PLAN.md:3481 puts at the foot of the card, frozen against
 * the demo alignment: where the allow actually happens, drawn with the 20px
 * ball glyph whose ring breathes on the 2s period that row's 动效 cell names.
 * Loop animation is allowed here and only here: PLAN.md:3491 lists 审批等待
 * among the non-idle states, and this card does not exist while the panel
 * is idle.
 */
function BallApproveHint() {
  return (
    <p className="flex items-center gap-1.5 text-[11px] text-muted-foreground">
      <svg
        aria-hidden="true"
        className="shrink-0 text-muted-foreground"
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
    <Card
      className="card-spotlight relative w-full max-w-[var(--panel-w)] gap-0 overflow-hidden rounded-xl border-border py-0 shadow-none"
      onMouseMove={(e) => {
        // The demo's spotlight: the gradient follows the pointer through the
        // two custom properties .card-spotlight::after reads. A style write on
        // the element itself, so no state and nothing remembered.
        const rect = e.currentTarget.getBoundingClientRect();
        e.currentTarget.style.setProperty("--mx", `${e.clientX - rect.left}px`);
        e.currentTarget.style.setProperty("--my", `${e.clientY - rect.top}px`);
      }}
    >
      {/* PLAN.md:3481's "视觉重量 for an irreversible operation": the bar
          across the top edge, drawn from the C21 danger token. bg-stop is the
          spelling the render harness counts; it resolves to the destructive
          hue for L2/Deny and the primary hue sits here only for the lighter
          levels. Exactly one per card. */}
      <div
        aria-hidden="true"
        className={cn("h-[2px] w-full shrink-0", strong ? "bg-stop" : "bg-primary")}
      />
      <CardHeader className="block px-4 pb-2 pt-3">
        <div className="flex items-start justify-between gap-3">
          <div className="flex min-w-0 items-center gap-2.5">
            <div className="flex size-8 shrink-0 items-center justify-center rounded-lg bg-muted/60 text-muted-foreground">
              <ToolIcon aria-hidden="true" className="size-4" />
            </div>
            <div className="min-w-0">
              <p className="text-[13px] font-medium text-foreground">
                需要你的确认
              </p>
              {/* 标题 = 工具名 + 徽标 (PLAN.md:3481); the harness reads both
                  out of this element, so they stay its direct children. */}
              <CardTitle className="mt-0.5 flex min-w-0 items-center gap-2 text-[11.5px] font-normal leading-snug text-muted-foreground">
                <span className="truncate font-mono">{view.tool}</span>
                <LevelBadge level={view.level} />
              </CardTitle>
            </div>
          </div>
          {/* demo's pagination-dots slot lives here; hidden on purpose - see
              the file header. */}
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
          <p className="font-mono text-[11.5px] text-muted-foreground">
            {`调用链：${view.callChain.join(" > ")}`}
          </p>
        )}
        {view.sessionOverrideBlocked && (
          <p className="text-[11.5px] font-medium text-stop">
            本调用带 C25 污染标记，会话级授权对它无效。
          </p>
        )}

        <div className="flex items-start gap-1.5 text-[11px] leading-[var(--lh-body)] text-muted-foreground">
          <Info aria-hidden="true" className="mt-px size-3.5 shrink-0" />
          <span>批准权在原生侧 · 面板不提供允许按钮</span>
        </div>
      </CardContent>

      <CardFooter className="flex-col items-stretch gap-2 border-t border-border px-4 py-2.5">
        {/* 常驻一行 (PLAN.md:3481): it is there whether or not anything else is. */}
        <BallApproveHint />
        <div className="flex items-center gap-2">
          <span className="mr-auto text-[11px] text-muted-foreground">
            {`判定来自 ${view.decidedBy === "native" ? "原生风险评估" : view.decidedBy}`}
          </span>
          <Button onClick={() => send("refuse")} size="sm" variant="destructive">
            拒绝
          </Button>
          <Button
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
