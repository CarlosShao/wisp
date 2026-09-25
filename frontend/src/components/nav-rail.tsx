/* ============================================================================
   Nav rail (Q1 = 甲) - the demo's own 60px icon strip, transcribed from
   design/doubao/demo/styles.css: .side-nav/.nav-item (:178-213) and the v4
   .nav-glide block (:663-678).

   Icons only, always; the Chinese name appears on hover AND on keyboard focus,
   in pure CSS (:hover / :focus on the button, see .nav-rail-label in
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
   scripts/render-nav.tsx keeps that ban nailed from this side too.
   ============================================================================ */

import { useEffect, useRef } from "react";
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

/** Keys are the demo's own icon names (panel-views.ts `icon` field, owner's
    2026-09-25 demo ruling). A miss is contract drift, so it throws instead of
    rendering a silent empty row. */
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
  const navRef = useRef<HTMLElement | null>(null);
  const glideRef = useRef<HTMLDivElement | null>(null);

  // nav-glide 滑翔指示条：demo 的招牌动效（styles.css:664-678，transform 走
  // --dur-glide 280ms）。量 active 行的 offsetTop 而不是拿「12px padding +
  // 行高 40 + 间距 2」做算术：布局尺寸一改，算术会悄悄错位，量出来永远是对的。
  // SSR（render-nav 的静态标记）不跑 effect，glide 停在 CSS 的 top:0 起点；
  // 默认屏是第一行，静态标记的几何仍然成立。
  useEffect(() => {
    const nav = navRef.current;
    const glide = glideRef.current;
    if (!nav || !glide) return;
    const row = nav.querySelector<HTMLButtonElement>(`[data-view="${active}"]`);
    if (row) glide.style.transform = `translateY(${row.offsetTop}px)`;
  }, [active]);

  return (
    <nav
      aria-label="面板视图"
      className="relative flex w-[var(--r-nav)] min-w-[var(--r-nav)] shrink-0 flex-col items-center gap-0.5 py-3"
      ref={navRef}
      style={{ background: "var(--nav-bg)", borderRight: "1px solid var(--window-border)" }}
    >
      <div aria-hidden="true" className="nav-glide" ref={glideRef} />
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
              "nav-rail-item flex size-10 shrink-0 items-center justify-center rounded-lg transition-colors duration-[var(--dur-fast)]",
              // active 行不带底色：那一格颜色属于 nav-glide 滑翔条，按钮自己
              // 保持透明（demo .nav-item.active 就是这么叠的）。
              on ? "text-primary" : "text-muted-foreground hover:bg-muted hover:text-foreground",
            )}
            data-view={v.id}
            key={v.id}
            onClick={() => onSelect?.(v.id)}
            type="button"
          >
            <Icon size={20} strokeWidth={1.75} />
            {v.id === "approval" && pendingCount ? (
              <span className="absolute -right-0.5 -top-0.5 flex size-4 items-center justify-center rounded-full border border-border bg-popover text-[10px] font-medium leading-none text-foreground">
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
