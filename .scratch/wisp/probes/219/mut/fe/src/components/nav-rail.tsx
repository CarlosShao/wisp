/* ============================================================================
   Nav rail (Q1 = 甲; icons per owner's 2026-09-25 ruling) - the library's
   vocabulary: bg-surface column, hover:bg-hover rows, active row lit by the
   accent-tint glide block with an accent-ink icon.
   ---------------------------------------------------------------------------

   Icons only, always; the Chinese name appears on hover AND on keyboard focus,
   in pure CSS (.nav-rail-label in src/styles/theme.css). One click changes the
   screen via onSelect UP to App - this component holds no view state of its
   own, so scripts/render-nav.tsx can assert the nav layer is stateless. The
   panel never tells the host which screen it moved to (Q-50 = 甲).
   ============================================================================ */

import {
  ChartColumn,
  Command,
  ListChecks,
  Lock,
  MessageSquareText,
  Orbit,
  Settings,
  ShieldAlert,
  ShieldCheck,
} from "lucide-react";
import { cn } from "@/lib/cn";
import { PANEL_VIEWS, type PanelViewId } from "@/lib/panel-views";

/** Keys are the names panel-views.ts draws (frozen §17.4 + the approved seven). */
const ICONS = {
  "message-square-text": MessageSquareText,
  "shield-check": ShieldCheck,
  command: Command,
  "list-checks": ListChecks,
  orbit: Orbit,
  settings: Settings,
  "shield-alert": ShieldAlert,
  lock: Lock,
  "chart-column": ChartColumn,
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
    <nav
      aria-label="面板视图"
      className="relative flex w-[60px] min-w-[60px] shrink-0 flex-col items-center gap-0.5 border-r border-line bg-surface py-3"
    >
      {/* 滑翔指示条：一枚 accent-tint 圆角块，280ms 滑到 active 行。偏移量 =
          容器顶距(py-3=12px) + 行距×序号；行距 = 40px 行 + 2px 间距 = 42px。
          三个数都由本文件的布局类决定（py-3 / size-10 / gap-0.5），改布局就
          同步改这两个常量——owner 2026-09-26 指出 glide 整体偏上一行距的 12px
          就是漏加了容器顶距。 */}
      <div
        aria-hidden="true"
        className="nav-glide"
        style={{ transform: `translateY(${12 + PANEL_VIEWS.findIndex((v) => v.id === active) * 42}px)` }}
      />
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
              "nav-rail-item flex size-10 shrink-0 items-center justify-center rounded-control transition-colors duration-150",
              on ? "text-accent-ink" : "text-ink-3 hover:bg-hover hover:text-ink",
            )}
            key={v.id}
            onClick={() => onSelect?.(v.id)}
            title={v.label}
            type="button"
          >
            <Icon size={20} strokeWidth={1.75} />
            {v.id === "approval" && pendingCount ? (
              <span className="absolute right-1 top-1 flex size-4 items-center justify-center rounded-full bg-accent text-[9px] font-medium leading-none text-white">
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
