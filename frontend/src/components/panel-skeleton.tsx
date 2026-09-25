/* ============================================================================
   Panel skeleton (ticket 77; navigation re-cut for Q1 = 甲 on 2026-09-25)
   ----------------------------------------------------------------------------
   The frame R19 scoped the first version to: a glass shell on C21's ambient
   background, the slot the L2 card occupies, and - since Q1 - a vertical icon
   rail down the LEFT edge instead of four text tabs across the top.

   The old header tabs (PANEL_TABS: 命令 / 结果 / 历史 / 配置) were 乙, the branch
   owner deprecated: text labels, no icons, and only four of the nine screens
   the demo has. They are gone rather than kept alongside, because two ways of
   changing screen in one panel is two answers to "which screen is this".

   It is still a layout, not a router: which view shows is decided by the
   snapshot the Go side pushes (Q2: no component holds it), so a reload loses
   nothing (PLAN.md:1044). The full-screen WebGL ambience R19 lists for phase
   two is deliberately absent - the CSS radial-gradient ambient is C21's own
   four lights, which cost no frame loop at all.
   ============================================================================ */

import { NavRail, railTitle } from "@/components/nav-rail";
import { Shimmer } from "@/components/ai-native/shimmer";
import type { PanelViewId } from "@/lib/panel-views";
import type { ReactNode } from "react";

export function PanelSkeleton({
  active,
  children,
  pendingCount,
  waitingLabel,
  onSelect,
}: {
  active: PanelViewId;
  children: ReactNode;
  pendingCount?: number;
  waitingLabel?: string;
  onSelect?: (id: PanelViewId) => void;
}) {
  return (
    <div className="wisp-shell flex min-h-screen w-full items-start justify-center p-4">
      <div
        className="glass-raised flex w-full max-w-[var(--panel-w)] overflow-hidden rounded-xl border border-border-hair"
        style={{ animation: "fade-up 350ms cubic-bezier(0.23,1,0.32,1) both" }}
      >
        <NavRail active={active} onSelect={onSelect} pendingCount={pendingCount} />
        <div className="flex min-w-0 flex-1 flex-col">
          <header className="flex items-center gap-3 border-b border-line px-4 py-2.5">
            <span className="text-[12.5px] font-medium text-ink">Wisp</span>
            <span className="text-[11.5px] text-ink-3">{railTitle(active)}</span>
            {/* Q1 = 甲 says one click changes the screen, and the panel is allowed
                to know which of its own screens it is showing: PLAN.md:1043-1044
                forbids REMEMBERING across a WebView restart, not having a current
                screen. So the pick lives in App and dies with the document, and a
                snapshot that names a view outranks it outright. What is still not
                wired is the other direction - the host is never told, which is
                ticket 35's push and not this file's business. */}
            {waitingLabel && (
              <Shimmer className="ml-auto text-[11.5px]">{waitingLabel}</Shimmer>
            )}
          </header>
          <main className="flex flex-col gap-3 p-4">{children}</main>
        </div>
      </div>
    </div>
  );
}
