/* ============================================================================
   Panel skeleton (ticket 77) - the frame R19 scoped the first version to:
   a glass shell on C21's ambient background, four tabs (command / result /
   history / settings) and the slot the L2 card occupies above them.

   It is a layout, not a router: which tab shows what is decided by the
   snapshot the Go side pushes, so a reload loses nothing (PLAN.md:1044).
   The full-screen WebGL ambience R19 lists for phase two is deliberately
   absent - the CSS radial-gradient ambient below is C21's own four lights,
   which cost no frame loop at all.
   ============================================================================ */

import { Shimmer } from "@/components/ai-native/shimmer";
import type { ReactNode } from "react";

export const PANEL_TABS = ["命令", "结果", "历史", "配置"] as const;
export type PanelTab = (typeof PANEL_TABS)[number];

export function PanelSkeleton({
  active,
  children,
  waitingLabel,
}: {
  active: PanelTab;
  children: ReactNode;
  waitingLabel?: string;
}) {
  return (
    <div className="wisp-shell flex min-h-screen w-full items-start justify-center p-4">
      <div
        className="glass-raised flex w-full max-w-[var(--panel-w)] flex-col overflow-hidden rounded-xl border-0"
        style={{ animation: "fade-up 350ms cubic-bezier(0.23,1,0.32,1) both" }}
      >
        <header className="flex items-center gap-3 border-b border-line px-4 py-2.5">
          <span className="text-[12.5px] font-medium text-ink">Wisp</span>
          <nav className="flex items-center gap-1" aria-label="面板分区">
            {PANEL_TABS.map((tab) => (
              <span
                aria-current={tab === active ? "page" : undefined}
                className={
                  "rounded-chip px-2 py-0.5 text-[11.5px] " +
                  (tab === active ? "bg-overlay text-ink" : "text-ink-3")
                }
                key={tab}
              >
                {tab}
              </span>
            ))}
          </nav>
          {waitingLabel && (
            <Shimmer className="ml-auto text-[11.5px]">{waitingLabel}</Shimmer>
          )}
        </header>
        <main className="flex flex-col gap-3 p-4">{children}</main>
      </div>
    </div>
  );
}
