/* ============================================================================
   Panel skeleton - fourth generation, the library's own chrome.
   ----------------------------------------------------------------------------
   Minimal window: a 44px header (ball dot + 一缕 + status), the 60px icon rail
   on the left, the content column on the right. Everything reads the
   beautiful-ui vocabulary (surface / ink / line) directly - there is no
   floating-window chrome any more: in production the WebView2 window IS the
   panel, and the library-style showcase (not a fake desktop) is where the
   visuals get accepted (main.tsx ?harness=1).

   Screen changes remount the content column (key={active}) so the screen-enter
   transition replays. The composer row is rendered by App AFTER children and
   docked with mt-auto - this main is a flex column for exactly that.
   ============================================================================ */

import { useEffect, useRef, useState } from "react";
import { NavRail } from "@/components/nav-rail";
import { Shimmer } from "@/components/ai-native/shimmer";
import { StatusPill } from "@/components/ai-native/status-pill";
import type { PanelViewId } from "@/lib/panel-views";
import type { ReactNode } from "react";

export function PanelSkeleton({
  active,
  children,
  pendingCount,
  onSelect,
}: {
  active: PanelViewId;
  children: ReactNode;
  pendingCount?: number;
  onSelect?: (id: PanelViewId) => void;
}) {
  // 屏切换的 shimmer 加载态：换屏那一刻亮起到过渡结束。首挂不算切换。
  const [entering, setEntering] = useState(false);
  const prevRef = useRef<PanelViewId>(active);
  useEffect(() => {
    if (prevRef.current === active) return;
    prevRef.current = active;
    setEntering(true);
    const t = window.setTimeout(() => setEntering(false), 300);
    return () => window.clearTimeout(t);
  }, [active]);

  return (
    <div className="flex h-screen w-full flex-col overflow-hidden bg-page text-ink">
      {/* 标题栏：一缕 + 状态。生产里窗口按钮归 Win32，这里永远不放假按钮。 */}
      <header className="flex h-11 shrink-0 select-none items-center justify-between gap-2 border-b border-line bg-surface pl-4 pr-3">
        <div className="flex items-center gap-2.5">
          <div className="size-3.5 rounded-full bg-accent" />
          <span className="text-[13px] font-medium text-ink">一缕</span>
        </div>
        <div className="flex items-center gap-2">
          {entering && <Shimmer className="text-[11px]">切换中</Shimmer>}
          {pendingCount ? (
            <StatusPill tone="orange">待审批 {pendingCount}</StatusPill>
          ) : null}
        </div>
      </header>
      <div className="flex min-h-0 flex-1 overflow-hidden">
        <NavRail active={active} onSelect={onSelect} pendingCount={pendingCount} />
        {/* key={active} 触发屏切换动画；flex 列让 App 的 composer 包能钉底。 */}
        <main
          className="screen-enter flex min-h-full min-w-0 flex-1 flex-col overflow-x-hidden overflow-y-auto bg-page px-6 pb-6 pt-5"
          key={active}
        >
          {children}
        </main>
      </div>
    </div>
  );
}
