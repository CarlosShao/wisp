/* ============================================================================
   Panel views (Q1 = 甲, owner 2026-09-25; frontend session)

   Q1 settled the panel's navigation as a vertical icon rail on the LEFT,
   copied from design/doubao/demo: icons only, the Chinese name on hover,
   one click to change screen. The nine rows below are that rail's whole set,
   taken entry-for-entry from the demo's own list (design/doubao/demo/app.js:18-26).

   The previous shape was four text tabs across the header (PANEL_TABS in
   src/components/panel-skeleton.tsx) - that was 乙, the branch owner deprecated.

   ---------------------------------------------------------------------------
   Icons: PLAN.md §17.4 is a FROZEN list of about 55 names, and adding to it is
   a human-approval change (that branch is 乙). So only names from that list may
   be imported here, even where the demo's own icon is a better fit.
   Measured 2026-09-25 (script D:/tmp/wisp-fe-icon-setdiff.mjs, reading the list
   out of PLAN.md and the rail out of app.js rather than from this comment):
   6 of the demo's 9 icon names are not in the frozen list. Four of those six
   have a near-relative that IS in the list (shield / list / circle-dot /
   settings-2) and read correctly. Two do not, and are marked `interim` below:
     对话  message-square-text -> layers  (a stack of turns; nothing chat-shaped
           exists in the list - chat / message / conversation are zero-hit)
     成本  chart-column        -> cpu     (compute, not money - coin / dollar /
           receipt / wallet / chart / gauge are all zero-hit in the list)
   The marker lives in the data, not only in a document, so
   scripts/render-nav.tsx can count it: if either row is ever given a proper
   name, that row's `interim` must be deleted or the harness goes red.
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
  /** A name from PLAN.md §17.4's frozen list. Nothing else may appear here;
      scripts/render-nav.tsx checks every value against the list itself. */
  icon: string;
  /** The icon the demo actually drew, kept so the substitution is auditable. */
  demoIcon: string;
  /** Required whenever `icon` differs from `demoIcon`: what the frozen list had
      instead of what the demo wanted. scripts/render-nav.tsx asserts that every
      such row carries one, and that no row which drew the demo's own name does. */
  note?: string;
  /** The narrower claim: this row's substitute is not merely different, it is
      not apt, because §17.4 has nothing in that shape at all. Exactly two rows
      qualify and owner's ruling counted them; adding a third means the frozen
      list has to change, which is branch 乙 and a human-approval change. */
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

const NOT_APT = "INTERIM(图标不贴切，Q1=甲 2026-09-25)";

export const PANEL_VIEWS: readonly PanelView[] = [
  {
    id: "chat",
    label: "对话",
    icon: "layers",
    demoIcon: "message-square-text",
    note: "§17.4 里 chat / message / conversation 零命中，退而取表示一轮轮叠起来的 layers",
    interim: NOT_APT,
    fed: true,
  },
  {
    id: "approval",
    label: "审批",
    icon: "shield",
    demoIcon: "shield-check",
    note: "清单里有 shield 与 shield-alert，没有 shield-check；少了那个勾",
    fed: true,
  },
  { id: "palette", label: "命令", icon: "command", demoIcon: "command", fed: false },
  {
    id: "tasks",
    label: "任务",
    icon: "list",
    demoIcon: "list-checks",
    note: "清单里有 list，没有 list-checks；同样只少勾",
    fed: false,
  },
  {
    id: "ball",
    label: "球状态",
    icon: "circle-dot",
    demoIcon: "orbit",
    note: "清单里没有 orbit，circle-dot 是同一件事（一颗带点的圆）的最直白写法",
    fed: false,
  },
  {
    id: "config",
    label: "设置",
    icon: "settings-2",
    demoIcon: "settings",
    note: "清单里只有 settings-2 这一枚齿轮，没有裸的 settings",
    fed: false,
    selfFed: true,
  },
  { id: "security", label: "安全", icon: "shield-alert", demoIcon: "shield-alert", fed: false },
  { id: "privacy", label: "隐私", icon: "lock", demoIcon: "lock", fed: false },
  {
    id: "cost",
    label: "成本",
    icon: "cpu",
    demoIcon: "chart-column",
    note: "§17.4 里 coin / dollar / receipt / wallet / chart / gauge 全部零命中，退而取表示算力的 cpu",
    interim: NOT_APT,
    fed: false,
  },
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
