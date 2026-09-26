/* ============================================================================
   HarnessSidebar - the library's Sidebar Nav, adapted to props (owner
   2026-09-26: "左边侧边栏人家也有对应的组件，你非要自己搞" - this replaces the
   hand-rolled sidebar).
   ----------------------------------------------------------------------------
   Derived from TurboKach/ai-native-react-components components/sidebar-nav.tsx
   (MIT, Copyright (c) 2026 Turbo, upstream commit 05dab2d2).

   Local changes, each load-bearing:
     - props-driven: workspace / history groups / usage / callbacks all arrive
       from the parent; upstream's hardcoded Creamery demo data is gone.
     - the flat "Workspace/Objects" sections become our collapsible history
       groups (今天 / 近 7 天 / 更早) - the demo README's "Collapsible ...
       chat navigation" row, with rotating chevrons.
     - items are session rows (label + time meta, no leading icon) and the
       accent action is 新建任务; the settings/usage/theme block joins the
       bottom under a hairline.
     - icons are lucide-react; upstream's inline SVG icon map is dropped.
     - TS6 strict interfaces; runtime statements otherwise kept.

   What survives verbatim: the workspace monogram row, the quick-search field
   with the "/" kbd chip, the accent action's plus-badge circle, the MEASURED
   gliding hover highlight (one bg-hover box that eases between rows, driven
   by useLayoutEffect over hovered ?? active), the count-badge pop-in curve
   for badges, and the hover-revealed plus affordance pattern.
   ============================================================================ */

import { useLayoutEffect, useRef, useState } from "react";
import { ChevronDown, Moon, Plus, Search, Settings, Sun } from "lucide-react";
import { cn } from "@/lib/cn";
import { ValuePill } from "@/components/ai-native/value-pill";
import type { HarnessSidebarGroup, HarnessUsagePill } from "@/fixtures/harness-app";

export interface HarnessWorkspaceInfo {
  /** 短名（monogram 与标题行）。 */
  name: string;
  /** 完整路径（副标题行）。 */
  path: string;
}

export interface HarnessSidebarProps {
  workspace: HarnessWorkspaceInfo;
  groups: readonly HarnessSidebarGroup[];
  activeSessionId: string | null;
  onSelectSession: (id: string) => void;
  onNewTask: () => void;
  searchValue: string;
  onSearchChange: (value: string) => void;
  searchPlaceholder?: string;
  usage: readonly HarnessUsagePill[];
  onOpenSettings: () => void;
  theme: "light" | "dark";
  onToggleTheme: () => void;
  className?: string;
}

export function HarnessSidebar({
  workspace,
  groups,
  activeSessionId,
  onSelectSession,
  onNewTask,
  searchValue,
  onSearchChange,
  searchPlaceholder = "搜索会话…",
  usage,
  onOpenSettings,
  theme,
  onToggleTheme,
  className,
}: HarnessSidebarProps) {
  const [collapsed, setCollapsed] = useState<ReadonlySet<string>>(new Set());
  const [hovered, setHovered] = useState<string | null>(null);
  const [box, setBox] = useState<{ top: number; height: number } | null>(null);
  const navRef = useRef<HTMLDivElement>(null);
  const itemRefs = useRef<Record<string, HTMLButtonElement | null>>({});

  // 滑翔高亮：一枚 bg-hover 方块量测着落到 hovered ?? active 行（上游逐字）。
  useLayoutEffect(() => {
    const container = navRef.current;
    const target = itemRefs.current[hovered ?? activeSessionId ?? ""];
    if (!container || !target) {
      setBox(null);
      return;
    }
    const containerRect = container.getBoundingClientRect();
    const targetRect = target.getBoundingClientRect();
    setBox({ top: targetRect.top - containerRect.top, height: targetRect.height });
  }, [hovered, activeSessionId, groups, collapsed]);

  function toggleGroup(title: string) {
    setCollapsed((current) => {
      const next = new Set(current);
      if (next.has(title)) next.delete(title);
      else next.add(title);
      return next;
    });
  }

  return (
    <aside
      className={cn(
        "flex h-full w-60 min-w-60 shrink-0 flex-col border-r border-line bg-surface",
        className,
      )}
    >
      <div className="flex min-h-0 flex-1 flex-col overflow-y-auto p-2">
        {/* workspace row（上游解剖，内容换 Wisp） */}
        <button
          type="button"
          className="mb-2 flex w-full items-center gap-2.5 rounded-control p-1.5 text-left transition-[background-color,transform] duration-100 hover:bg-hover active:scale-[0.96]"
          title={workspace.path}
        >
          <span className="flex size-8 shrink-0 items-center justify-center rounded-[8px] bg-ink text-[13px] font-semibold text-surface">
            {workspace.name.slice(0, 1)}
          </span>
          <span className="min-w-0 flex-1">
            <span className="block truncate text-[13px] font-medium leading-tight text-ink">
              {workspace.name}
            </span>
            <span className="block truncate text-[11px] leading-tight text-ink-3">
              个人工作区
            </span>
          </span>
          <ChevronDown aria-hidden="true" className="shrink-0 rotate-180 text-ink-3" size={12} strokeWidth={2} />
        </button>

        {/* quick search（上游逐字：inset 底 + "/" kbd 角标） */}
        <label className="mb-1 flex h-8 items-center gap-2 rounded-control bg-inset px-2.5 shadow-hairline">
          <Search aria-hidden="true" className="shrink-0 text-ink-3" size={12} strokeWidth={2} />
          <input
            value={searchValue}
            onChange={(event) => onSearchChange(event.target.value)}
            placeholder={searchPlaceholder}
            className="min-w-0 flex-1 bg-transparent text-[12.5px] text-ink outline-none placeholder:text-ink-3"
          />
          <kbd className="flex size-4.5 shrink-0 items-center justify-center rounded-[5px] bg-surface text-[10px] text-ink-3 shadow-hairline">
            /
          </kbd>
        </label>

        {/* accent action：新建任务（上游逐字的 accent 行 + 加号徽章圆） */}
        <button
          type="button"
          onClick={onNewTask}
          className="mb-2 flex w-full items-center gap-2 rounded-control px-2 py-1.5 text-[13px] font-medium text-accent transition-[background-color,transform] duration-100 hover:bg-accent-tint active:scale-[0.96]"
        >
          <span className="min-w-0 flex-1 truncate text-left">新建任务</span>
          <span className="flex size-4 shrink-0 items-center justify-center rounded-full bg-accent text-white">
            <Plus aria-hidden="true" size={9} strokeWidth={3} />
          </span>
        </button>

        {/* items：滑翔高亮 + 可折叠分组 + 会话行 */}
        <div ref={navRef} className="relative flex flex-col gap-2" onMouseLeave={() => setHovered(null)}>
          <span
            aria-hidden="true"
            className="pointer-events-none absolute inset-x-0 rounded-[7px] bg-hover"
            style={{
              top: box?.top ?? 0,
              height: box?.height ?? 0,
              opacity: box ? 1 : 0,
              transition:
                "top 220ms cubic-bezier(0.23,1,0.32,1), height 220ms cubic-bezier(0.23,1,0.32,1), opacity 150ms ease",
            }}
          />
          {groups.map((group) => {
            const isCollapsed = collapsed.has(group.label);
            return (
              <div key={group.label}>
                <button
                  type="button"
                  onClick={() => toggleGroup(group.label)}
                  className="group flex w-full items-center gap-1 px-2 pb-1 pt-1 text-[10.5px] font-medium uppercase tracking-[0.08em] text-ink-3"
                >
                  <ChevronDown
                    aria-hidden="true"
                    className={cn("transition-transform duration-150", isCollapsed && "-rotate-90")}
                    size={10}
                    strokeWidth={2.2}
                  />
                  <span className="min-w-0 flex-1 text-left">{group.label}</span>
                </button>
                {!isCollapsed && (
                  <div className="flex flex-col gap-px">
                    {group.rows.map((row) => {
                      const isActive = row.id === activeSessionId;
                      return (
                        <button
                          key={row.id}
                          ref={(el) => {
                            itemRefs.current[row.id] = el;
                          }}
                          type="button"
                          aria-current={isActive ? "page" : undefined}
                          onMouseEnter={() => setHovered(row.id)}
                          onFocus={() => setHovered(row.id)}
                          onBlur={() => setHovered(null)}
                          onClick={() => onSelectSession(row.id)}
                          className={cn(
                            "relative z-10 flex w-full items-center gap-2 rounded-[7px] px-2 py-1.5 text-left transition-[color,transform] duration-150 active:scale-[0.96]",
                            isActive ? "font-medium text-ink" : "text-ink-2",
                          )}
                        >
                          <span className="min-w-0 flex-1 truncate text-[13px]">{row.title}</span>
                          {row.badge ? (
                            <span className="shrink-0 text-[10.5px] text-ink-3">{row.badge.label}</span>
                          ) : null}
                          <span className="shrink-0 text-[10.5px] tabular-nums text-ink-3">{row.time}</span>
                        </button>
                      );
                    })}
                  </div>
                )}
              </div>
            );
          })}
          {groups.every((g) => g.rows.length === 0) && (
            <div className="px-2 py-3 text-[12px] text-ink-3">没有匹配的会话</div>
          )}
        </div>
      </div>

      {/* 底部：用量 + 设置 + 主题（此块不进滑翔区，发丝线分隔） */}
      <div className="flex flex-col gap-1 border-t border-line p-2">
        <div className="flex items-center gap-1.5 px-2 py-1">
          {usage.map((u) => (
            <span className="flex min-w-0 items-center gap-1" key={u.label}>
              <span className="text-[10.5px] text-ink-3">{u.label}</span>
              <ValuePill>{u.value}</ValuePill>
            </span>
          ))}
        </div>
        <button
          type="button"
          onClick={onOpenSettings}
          className="flex w-full items-center gap-2 rounded-[7px] px-2 py-1.5 text-left text-[12.5px] text-ink-2 transition-colors duration-150 hover:bg-hover hover:text-ink"
        >
          <Settings aria-hidden="true" size={13} strokeWidth={1.8} />
          设置
        </button>
        <button
          type="button"
          onClick={onToggleTheme}
          className="flex w-full items-center gap-2 rounded-[7px] px-2 py-1.5 text-left text-[12.5px] text-ink-2 transition-colors duration-150 hover:bg-hover hover:text-ink"
        >
          {theme === "dark" ? (
            <Sun aria-hidden="true" size={13} strokeWidth={1.8} />
          ) : (
            <Moon aria-hidden="true" size={13} strokeWidth={1.8} />
          )}
          {theme === "dark" ? "浅色" : "暗色"}
        </button>
      </div>
    </aside>
  );
}
