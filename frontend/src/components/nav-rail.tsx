/* ============================================================================
   Nav rail (Q1 = 甲) - the left vertical icon strip that replaces the old
   header text tabs.

   Icons only, always; the Chinese name appears on hover AND on keyboard focus,
   in pure CSS (:hover / :focus-visible on the button, see .nav-rail-label in
   src/styles/theme.css). No local state holds which row is lit, because
   `active` is a prop and the only way to change it is a new snapshot (Q2).

   Clicking asks the native side; it does not switch the screen here. That is
   not shyness about the click - PLAN.md:1044 and SPEC-08:150 want the panel to
   remember nothing, and a rail that flipped its own row would be the first
   piece of state the panel had ever kept.
   ============================================================================ */

import {
  CircleDot,
  Command,
  Cpu,
  Layers,
  List,
  Lock,
  Settings2,
  Shield,
  ShieldAlert,
} from "lucide-react";
import { cn } from "@/lib/cn";
import { PANEL_VIEWS, viewOf, type PanelViewId } from "@/lib/panel-views";

/** Keys are PLAN.md §17.4 names; the harness asserts every rail row is in here. */
const ICONS = {
  "circle-dot": CircleDot,
  command: Command,
  cpu: Cpu,
  layers: Layers,
  list: List,
  lock: Lock,
  "settings-2": Settings2,
  shield: Shield,
  "shield-alert": ShieldAlert,
} as const;

export function NavRail({
  active,
  pendingCount,
}: {
  active: PanelViewId;
  pendingCount?: number;
}) {
  return (
    <nav aria-label="面板视图" className="flex w-[44px] shrink-0 flex-col items-center gap-1 border-r border-line py-3">
      {PANEL_VIEWS.map((v) => {
        const Icon = ICONS[v.icon as keyof typeof ICONS];
        if (!Icon) {
          throw new Error(`nav-rail: ${v.icon} is not a drawn icon`);
        }
        const on = v.id === active;
        return (
          <button
            aria-current={on ? "page" : undefined}
            aria-disabled
            aria-label={v.label}
            className={cn(
              "nav-rail-item relative flex size-8 shrink-0 items-center justify-center rounded-control",
              on ? "bg-overlay text-ink" : "text-ink-3 hover:bg-hover hover:text-ink",
            )}
            disabled
            key={v.id}
            title={`${v.label}（换屏待接线）`}
            type="button"
          >
            <Icon size={16} />
            {v.id === "approval" && pendingCount ? (
              <span className="absolute -right-0.5 -top-0.5 rounded-full bg-inset px-1 text-[9px] leading-4 text-ink">
                {pendingCount}
              </span>
            ) : null}
            <span className="nav-rail-label">{v.label}</span>
          </button>
        );
      })}
    </nav>
  );
}

/** The name of the row the panel is on, for the header line. */
export function railTitle(active: PanelViewId): string {
  return viewOf(active).label;
}
