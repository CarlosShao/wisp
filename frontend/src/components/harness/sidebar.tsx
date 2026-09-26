/* ============================================================================
   HarnessSidebar - the library's Sidebar Nav, adapted to props, body now in
   Qoder's three-section form (owner 2026-09-26, second form ruling from the
   Qoder screenshot: 分组段头 自动化 / 工作区 / 最近任务, muted small text).
   ----------------------------------------------------------------------------
   Derived from TurboKach/ai-native-react-components components/sidebar-nav.tsx
   (MIT, Copyright (c) 2026 Turbo, upstream commit 05dab2d2).

   Local changes, each load-bearing:
     - props-driven: workspace / automations / project / history groups /
       usage / callbacks all arrive from the parent; upstream's hardcoded
       Creamery demo data is gone.
     - PROPS CONTRACT lives in src/fixtures/harness-app.ts (HarnessSidebarProps
       + HarnessWorkspaceInfo); this file imports it instead of re-declaring -
       same one-contract-page rule main.tsx follows. HarnessWorkspaceInfo is
       re-exported here so the old import path keeps working.
     - Qoder three-section body (2026-09-26): below the preserved top anatomy,
       the nav is three sections with muted small headers:
         自动化   - status-dot rows (accent dot = running, ink-3 dot = idle)
                    with a hover-revealed gear slot on the right edge;
         工作区   - a project row (chevron + Folder icon + name, collapsible)
                    over nested session children (indented, leading status
                    icon: red CircleAlert = failed, accent dot = streaming,
                    nothing = done);
         最近任务 - the previous history groups (今天 / 近 7 天 / 更早),
                    collapsible chevrons and session rows unchanged.
     - the MEASURED gliding hover highlight (one bg-hover box easing between
       rows, useLayoutEffect over hovered ?? active) now measures EVERY row
       kind - automation rows, the project row, children, history rows - via
       keyed refs ("auto:*" / "project" / "child:*" / plain session ids).
     - what survives verbatim from the upstream anatomy: the workspace
       monogram row, the quick-search field with the "/" kbd chip, the accent
       action's plus-badge circle, the gliding highlight mechanics, the
       count-badge pop-in curve, the hover-revealed affordance pattern, and
       the bottom usage/settings/theme block under a hairline.
     - icons are lucide-react; upstream's inline SVG icon map is dropped.
     - TS6 strict interfaces; runtime statements otherwise kept.

   Colours: token utilities only, zero literals; no emoji.
   ============================================================================ */

import { useLayoutEffect, useRef, useState } from "react";
import { ChevronDown, CircleAlert, Folder, Moon, Plus, Search, Settings, Sun } from "lucide-react";
import { cn } from "@/lib/cn";
import { ValuePill } from "@/components/ai-native/value-pill";
import type {
  HarnessProjectChildStatus,
  HarnessSidebarProps,
  HarnessWorkspaceInfo,
} from "@/fixtures/harness-app";

export type { HarnessWorkspaceInfo };

/** 工作区段嵌套会话行的行首状态图标：红圈叹号=failed（已拒绝）、蓝点=
    streaming（进行中）、无=done（已完成，占位空槽保持对齐）。 */
function ChildStatusIcon({ status }: { status: HarnessProjectChildStatus }) {
  return (
    <span className="flex w-3.5 shrink-0 justify-center">
      {status === "failed" ? (
        <CircleAlert aria-hidden="true" className="shrink-0 text-red" size={12} strokeWidth={2} />
      ) : status === "streaming" ? (
        <span aria-hidden="true" className="block size-1.5 rounded-full bg-accent" />
      ) : null}
    </span>
  );
}

/** 三段的分组段头：muted 小字（Qoder 形态），不折叠。 */
function SectionHeader({ label }: { label: string }) {
  return (
    <div className="px-2 pb-1 pt-1.5 text-[10.5px] font-medium tracking-[0.08em] text-ink-3">
      {label}
    </div>
  );
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
  automations,
  onOpenAutomation,
  project,
  className,
}: HarnessSidebarProps) {
  const [collapsed, setCollapsed] = useState<ReadonlySet<string>>(new Set());
  const [projectOpen, setProjectOpen] = useState(true);
  const [hovered, setHovered] = useState<string | null>(null);
  const [box, setBox] = useState<{ top: number; height: number } | null>(null);
  const navRef = useRef<HTMLDivElement>(null);
  const itemRefs = useRef<Record<string, HTMLButtonElement | HTMLDivElement | null>>({});

  const setItemRef = (key: string) => (el: HTMLButtonElement | HTMLDivElement | null) => {
    itemRefs.current[key] = el;
  };

  // 滑翔高亮：一枚 bg-hover 方块量测着落到 hovered ?? active 行（上游逐字，
  // 键域扩到三段全部行种）。行被折叠隐藏时量测落空，方块淡出。
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
  }, [hovered, activeSessionId, groups, collapsed, projectOpen, automations, project]);

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

        {/* 三段式主体：滑翔高亮覆盖全部行种 */}
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

          {/* ── 自动化段：状态圆点行 + 右缘齿轮位（hover 显现） ── */}
          <div>
            <SectionHeader label="自动化" />
            <div className="flex flex-col gap-px">
              {automations.map((automation) => (
                <div
                  key={automation.id}
                  ref={setItemRef(`auto:${automation.id}`)}
                  onMouseEnter={() => setHovered(`auto:${automation.id}`)}
                  className="group relative z-10 flex w-full items-center gap-2 rounded-[7px] px-2 py-1.5"
                >
                  <span
                    aria-hidden="true"
                    className={cn(
                      "size-1.5 shrink-0 rounded-full",
                      automation.state === "running" ? "bg-accent" : "bg-ink-3",
                    )}
                  />
                  <span className="min-w-0 flex-1 truncate text-[13px] text-ink-2">
                    {automation.name}
                  </span>
                  <button
                    type="button"
                    aria-label={`打开 ${automation.name} 设置`}
                    className="flex size-5 shrink-0 items-center justify-center rounded-chip text-ink-3 opacity-0 transition-[opacity,background-color,color] duration-150 hover:bg-hover hover:text-ink focus-visible:opacity-100 group-hover:opacity-100"
                    onClick={() => onOpenAutomation?.(automation.id)}
                  >
                    <Settings aria-hidden="true" size={12} strokeWidth={1.8} />
                  </button>
                </div>
              ))}
            </div>
          </div>

          {/* ── 工作区段：项目行（可折叠）+ 嵌套会话行（缩进 + 行首状态图标） ── */}
          <div>
            <SectionHeader label="工作区" />
            <button
              ref={setItemRef("project")}
              type="button"
              aria-expanded={projectOpen}
              onMouseEnter={() => setHovered("project")}
              onClick={() => setProjectOpen((current) => !current)}
              className="relative z-10 flex w-full items-center gap-1.5 rounded-[7px] px-2 py-1.5 text-left transition-[color,transform] duration-150 active:scale-[0.96]"
            >
              <ChevronDown
                aria-hidden="true"
                className={cn(
                  "shrink-0 text-ink-3 transition-transform duration-150",
                  !projectOpen && "-rotate-90",
                )}
                size={10}
                strokeWidth={2.2}
              />
              <Folder aria-hidden="true" className="shrink-0 text-ink-2" size={13} strokeWidth={1.8} />
              <span className="min-w-0 flex-1 truncate text-[13px] text-ink-2">{project.name}</span>
            </button>
            {projectOpen && (
              <div className="flex flex-col gap-px">
                {project.children.map((child) => {
                  const isActive = child.id === activeSessionId;
                  return (
                    <button
                      key={child.id}
                      ref={setItemRef(`child:${child.id}`)}
                      type="button"
                      aria-current={isActive ? "page" : undefined}
                      onMouseEnter={() => setHovered(`child:${child.id}`)}
                      onFocus={() => setHovered(`child:${child.id}`)}
                      onBlur={() => setHovered(null)}
                      onClick={() => onSelectSession(child.id)}
                      className={cn(
                        "relative z-10 flex w-full items-center gap-2 rounded-[7px] py-1.5 pl-6 pr-2 text-left transition-[color,transform] duration-150 active:scale-[0.96]",
                        isActive ? "font-medium text-ink" : "text-ink-2",
                      )}
                    >
                      <ChildStatusIcon status={child.status} />
                      <span className="min-w-0 flex-1 truncate text-[13px]">{child.title}</span>
                    </button>
                  );
                })}
              </div>
            )}
          </div>

          {/* ── 最近任务段：现会话历史（分组折叠 / 会话行 / 滑翔解剖不变） ── */}
          <div>
            <SectionHeader label="最近任务" />
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
                            ref={setItemRef(row.id)}
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
