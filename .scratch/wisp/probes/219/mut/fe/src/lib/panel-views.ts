/* ============================================================================
   Panel views (owner ruling 2026-09-25, second pass; frontend owner-agent)

   The rail is a vertical icon strip on the LEFT, copied entry-for-entry from
   the demo (design/doubao/demo/app.js:18-26). Q1 = 甲 settled the shape
   (icons only, the Chinese name on hover, one click to change screen).

   ---------------------------------------------------------------------------
   Icons, second pass. PLAN.md §17.4 is a frozen list of about 55 names and
   adding to it is a human-approval change (branch 乙). On 2026-09-25 owner
   granted that approval in person: the frontend was told the demo wins where
   it collides with the frozen spec, and the rail was rebuilt to draw the
   demo's OWN icons. That ruling covers exactly the seven names the demo uses
   that §17.4 does not list (message-square-text, shield-check, list-checks,
   orbit, settings, chart-column, life-buoy); scripts/render-nav.tsx checks
   every row against frozen + approved and fails on anything else.

   The Q1 = 甲 substitution ledger (layers -> 对话, cpu -> 成本) is therefore
   RETIRED: every row now draws its demoIcon, and no row may carry a `note` or
   an `interim` marker any more. If a future row ever has to substitute again,
   the ledger shape comes back with it.
   ============================================================================ */

import type { PanelSnapshot } from "@/lib/panel";

export type PanelViewId =
  | "chat"
  | "approval"
  | "palette"
  | "tasks"
  | "ball"
  | "config"
  | "security"
  | "privacy"
  | "cost";

export interface PanelView {
  id: PanelViewId;
  /** The Chinese label the rail reveals on hover and focus. */
  label: string;
  /** The icon the demo drew. Since owner's 2026-09-25 ruling this is always a
      legal name: it is either in PLAN.md §17.4's frozen list or in the
      seven-name approved extension, and render-nav.tsx checks both sets. */
  icon: string;
  /** Kept for auditability: the demo's own name for this row. A row that does
      not draw it (icon !== demoIcon) must say why in `note`. */
  demoIcon: string;
  /** Required whenever `icon` differs from `demoIcon`. Empty today - owner's
      ruling retired the substitutions. */
  note?: string;
  /** The Q1 = 甲 interim marker. Retired by the same ruling; render-nav.tsx
      fails if any row carries one again. */
  interim?: string;
  /** Whether Go has anything to push for this screen yet. False renders a
      named empty state - never invented content (owner's P9 red line). */
  fed: boolean;
  /** Whether the panel can put something real on this screen WITHOUT Go. Only
      设置 qualifies today: its one live control moves a style on this page and
      reads nothing from anywhere. `fed: false, selfFed: true` is therefore the
      honest shape for a screen that is neither empty nor fed - the gap notice
      would be a lie, and claiming `fed` would be a lie the other way. */
  selfFed?: boolean;
}

export const PANEL_VIEWS: readonly PanelView[] = [
  { id: "chat", label: "对话", icon: "message-square-text", demoIcon: "message-square-text", fed: true },
  { id: "approval", label: "审批", icon: "shield-check", demoIcon: "shield-check", fed: true },
  { id: "palette", label: "命令", icon: "command", demoIcon: "command", fed: false },
  { id: "tasks", label: "任务", icon: "list-checks", demoIcon: "list-checks", fed: false },
  { id: "ball", label: "球状态", icon: "orbit", demoIcon: "orbit", fed: false },
  { id: "config", label: "设置", icon: "settings", demoIcon: "settings", fed: false, selfFed: true },
  { id: "security", label: "安全", icon: "shield-alert", demoIcon: "shield-alert", fed: false },
  { id: "privacy", label: "隐私", icon: "lock", demoIcon: "lock", fed: false },
  { id: "cost", label: "成本", icon: "chart-column", demoIcon: "chart-column", fed: false },
];

/** The screen the panel opens on when the snapshot names none. Chosen because
    it is the only row with anything behind it today; it is a fallback for a
    missing field, not a memory of where the user last was. */
export const DEFAULT_VIEW: PanelViewId = "chat";

export function viewOf(id: PanelViewId): PanelView {
  const found = PANEL_VIEWS.find((v) => v.id === id);
  if (!found) throw new Error(`panel-views: no view ${JSON.stringify(id)}`);
  return found;
}

/**
 * The one honest answer to "which screen is this". The field does not exist on
 * PanelSnapshot yet - ticket 35's pump owns it, and `panel.resync` (the push
 * that would re-establish it on every show) is still zero-hit repo-wide. Until
 * then a screen change is asked for and refused, so what is on screen stays
 * whatever the caller passed.
 *
 * Reading it through this one function is what keeps the swap a one-line change
 * later: no component holds a view, and none derives one either.
 */
export function currentView(snapshot: PanelSnapshot & { view?: string }): PanelViewId {
  const named = snapshot.view;
  if (named === undefined) return DEFAULT_VIEW;
  if (PANEL_VIEWS.some((v) => v.id === named)) return named as PanelViewId;
  // A snapshot naming a screen the frontend has no row for is a contract drift,
  // not a reason to guess: say so, and show the default.
  console.error(`panel-views: snapshot names unknown view ${JSON.stringify(named)}`);
  return DEFAULT_VIEW;
}
