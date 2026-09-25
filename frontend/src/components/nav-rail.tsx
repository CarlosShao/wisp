/* ============================================================================
   Nav rail (Q1 = 甲) - the left vertical icon strip that replaces the old
   header text tabs.

   Icons only, always; the Chinese name appears on hover AND on keyboard focus,
   in pure CSS (:hover / :focus-visible on the button, see .nav-rail-label in
   src/styles/theme.css).

   One click changes the screen, which is what 甲 promised. The click goes UP to
   App as onSelect - this component holds no view state of its own, so
   scripts/render-nav.tsx can still assert that the nav layer is stateless. App
   keeps the pick for the life of the document and nothing longer:
   PLAN.md:1043-1044 forbids remembering across a WebView restart, not having a
   current screen, and a snapshot that names a view outranks the local pick
   outright (see App.tsx).

   What is still NOT wired is the other direction: the panel never tells the host
   which screen it moved to. Q-50 = 甲 bans a fifth outbound route, and
   scripts/render-nav.tsx:221 keeps that ban nailed from this side too.
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
  onSelect,
}: {
  active: PanelViewId;
  pendingCount?: number;
  onSelect?: (id: PanelViewId) => void;
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
            aria-label={v.label}
            className={cn(
              "nav-rail-item relative flex size-8 shrink-0 items-center justify-center rounded-control",
              on ? "bg-nav-active text-primary" : "text-ink-3 hover:bg-hover hover:text-ink",
            )}
            key={v.id}
            onClick={() => onSelect?.(v.id)}
            title={v.label}
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
