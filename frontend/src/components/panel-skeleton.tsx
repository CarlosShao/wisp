/* ============================================================================
   Panel skeleton - the demo's window, transcribed from
   design/doubao/demo/index.html:85-115 and styles.css .main-window/.title-bar/
   .window-body/.content-area (:112-235, v4 质感层 :733-741).

   Structure: title bar (ball-dot + 「一缕」 on the left; the button cluster on
   the right only under harness chrome), then the body - the 60px NavRail on the
   left and the content area on the right. Each screen change remounts the
   content area (key={active}), which is what replays the demo's screen-enter
   transition (rbScreenIn, duration read from --dur-screen rather than copied).

   `chrome` is the harness switch: main.tsx passes it only under ?harness=1, and
   it buys the floating-window dressing a browser has to fake (desktop backdrop,
   fog spheres, 720px window with radius + hairline + layered shadow). In
   production the WebView2 window IS the panel, so the shell renders bare - and
   the title bar's right side stays button-free, because a min/close button that
   does nothing would be a fake control in a window where Win32 owns those. What
   production does get there is status: the shimmer label while a screen
   transition runs, and the pending-approval pill when the queue is not empty.
   ============================================================================ */

import { useEffect, useRef, useState } from "react";
import { Minus, Sun, X } from "lucide-react";
import { NavRail } from "@/components/nav-rail";
import { Shimmer } from "@/components/ai-native/shimmer";
import { StatusPill } from "@/components/ai-native/status-pill";
import type { PanelViewId } from "@/lib/panel-views";
import type { ReactNode } from "react";

/** demo 屏切换时长。运行时读 token（--dur-screen）而不是把数字抄进代码 - 抄了
    就是第二份真相源；没有 document 的渲染环境（render-nav 的静态标记）回落到
    demo 原值 220ms。 */
function screenDurationMs(): number {
  if (typeof document === "undefined") return 220;
  const raw = getComputedStyle(document.documentElement).getPropertyValue("--dur-screen").trim();
  const n = Number.parseFloat(raw);
  return Number.isFinite(n) && n > 0 ? n : 220;
}

export function PanelSkeleton({
  active,
  children,
  pendingCount,
  onSelect,
  chrome = false,
}: {
  active: PanelViewId;
  children: ReactNode;
  pendingCount?: number;
  onSelect?: (id: PanelViewId) => void;
  chrome?: boolean;
}) {
  // 屏切换的 shimmer 加载态：从 active 变化那一刻起，到过渡结束为止（demo 各
  // 屏切换的 shimmer 过渡）。首挂不算切换 - 首屏谈不上「切换中」。
  const [entering, setEntering] = useState(false);
  const prevRef = useRef<PanelViewId>(active);
  useEffect(() => {
    if (prevRef.current === active) return;
    prevRef.current = active;
    setEntering(true);
    const t = window.setTimeout(() => setEntering(false), screenDurationMs());
    return () => window.clearTimeout(t);
  }, [active]);

  const frame = (
    <>
      {/* 标题栏：demo .title-bar（40px、左 16 右 12 的两段 padding、发丝底边）。
          底色与边线是 demo 自己的 token（--titlebar-bg / --window-border），没进
          Tailwind 命名空间，走内联 var() - main.tsx 的横幅同款做法。 */}
      <header
        className="flex h-[var(--titlebar-h)] shrink-0 select-none items-center justify-between gap-2 pl-4 pr-3"
        style={{ background: "var(--titlebar-bg)", borderBottom: "1px solid var(--window-border)" }}
      >
        <div className="flex items-center gap-2.5">
          <div className="ball-dot size-4 rounded-full" />
          <span className="text-[13px] font-medium text-foreground/80">一缕</span>
        </div>
        <div className="flex items-center gap-1">
          {entering && <Shimmer className="text-[11px]">切换中</Shimmer>}
          {pendingCount ? (
            <StatusPill tone="orange">待审批 {pendingCount}</StatusPill>
          ) : null}
          {chrome && (
            <div className="flex items-center gap-1">
              {/* demo 标题栏按钮位（styles.css:148-168：28px、6px 圆角、120ms
                  变色；关闭位悬停 destructive）。harness 里的造型位，不带行为 -
                  真正的窗口按钮归 Win32。 */}
              <button
                aria-label="切换主题"
                className="flex size-7 items-center justify-center rounded-sm text-muted-foreground transition-colors duration-[var(--dur-fast)] hover:bg-muted hover:text-foreground"
                title="切换主题"
                type="button"
              >
                <Sun size={16} />
              </button>
              <button
                aria-label="最小化"
                className="flex size-7 items-center justify-center rounded-sm text-muted-foreground transition-colors duration-[var(--dur-fast)] hover:bg-muted hover:text-foreground"
                title="最小化"
                type="button"
              >
                <Minus size={16} />
              </button>
              <button
                aria-label="关闭"
                className="flex size-7 items-center justify-center rounded-sm text-muted-foreground transition-colors duration-[var(--dur-fast)] hover:bg-destructive hover:text-destructive-foreground"
                title="关闭"
                type="button"
              >
                <X size={16} />
              </button>
            </div>
          )}
        </div>
      </header>
      <div className="flex min-h-0 flex-1 overflow-hidden">
        <NavRail active={active} onSelect={onSelect} pendingCount={pendingCount} />
        {/* key={active} 是屏切换动画的触发器：换屏即整块内容重挂，.screen-enter
            （rbScreenIn）随之重放。滚动条照 demo .content-area：6px、圆头、
            border 色拇指、悬停加深（styles.css:229-235）。 */}
        <main
          className="screen-enter flex min-h-full min-w-0 flex-1 flex-col overflow-x-hidden overflow-y-auto px-7 pb-8 pt-6 [&::-webkit-scrollbar]:w-[6px] [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:bg-border [&::-webkit-scrollbar-thumb:hover]:bg-muted-foreground [&::-webkit-scrollbar-track]:bg-transparent"
          key={active}
        >
          {children}
        </main>
      </div>
    </>
  );

  if (chrome) {
    return (
      <div className="harness-desktop fixed inset-0 flex items-center justify-center overflow-hidden">
        <div aria-hidden="true" className="fog-sphere fog-sphere-1" />
        <div aria-hidden="true" className="fog-sphere fog-sphere-2" />
        {/* demo .main-window：720x780、12px 圆角、发丝边、三层投影（几何与
            投影都来自 token：--panel-w/--panel-h/--r-xl/--shadow-window-v4，
            由 .wisp-window 携带）。z 压过两枚 fog-sphere。 */}
        <div
          className="wisp-shell wisp-window relative z-50 flex flex-col"
          style={{ width: "var(--panel-w)", height: "var(--panel-h)" }}
        >
          {frame}
        </div>
      </div>
    );
  }
  return <div className="wisp-shell flex h-screen w-full flex-col overflow-hidden">{frame}</div>;
}
