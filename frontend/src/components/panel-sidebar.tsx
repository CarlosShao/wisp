/* ============================================================================
   Panel sidebar (product path only, 2026-09-27).
   ----------------------------------------------------------------------------
   The harness skeleton's left column, adapted for the product panel: no
   fixture data, honest empty states, and the four view sections (chat / task /
   ball / settings) rendered as sidebar buttons instead of a separate icon rail.

   The shape is Qoder's three-section form (automations / workspace / history),
   but the product path has no data for any of them today, so every section
   shows the honest empty state rather than invented rows.
   ============================================================================ */

import { Moon, Plus, Settings, Sun } from "lucide-react";
import { GlideMenu } from "@/components/ai-native/glide-menu";
import { cn } from "@/lib/cn";

export type PanelSidebarSection =
  | "chat"
  | "approval"
  | "palette"
  | "tasks"
  | "ball"
  | "config"
  | "security"
  | "privacy"
  | "cost"
  | "settings";

export interface PanelSidebarProps {
  active: PanelSidebarSection;
  onSelect: (section: PanelSidebarSection) => void;
  onNewTask: () => void;
  onOpenSettings: () => void;
  onToggleTheme: () => void;
  theme: "light" | "dark";
  pendingCount?: number;
}

function SectionHeader({ label }: { label: string }) {
  return (
    <div className="px-2 pb-1 pt-1.5 text-[10.5px] font-medium tracking-[0.08em] text-ink-3">
      {label}
    </div>
  );
}

export function PanelSidebar({
  active,
  onSelect,
  onNewTask,
  onOpenSettings,
  onToggleTheme,
  theme,
  pendingCount,
}: PanelSidebarProps) {
  return (
    <aside className="flex h-full w-60 min-w-60 shrink-0 flex-col border-r border-line bg-surface">
      {/* The aria-label is the gate's needle: render-nav counts rows inside
          this landmark to prove the sidebar renders every view. */}
      <nav aria-label="面板视图" className="flex min-h-0 flex-1 flex-col overflow-y-auto p-2">
        {/* workspace row */}
        <button
          type="button"
          className="mb-2 flex w-full items-center gap-2.5 rounded-control p-1.5 text-left transition-[background-color,transform] duration-100 hover:bg-hover active:scale-[0.96]"
        >
          <span className="flex size-8 shrink-0 items-center justify-center rounded-[8px] bg-ink text-[13px] font-semibold text-surface">
            W
          </span>
          <span className="min-w-0 flex-1">
            <span className="block truncate text-[13px] font-medium leading-tight text-ink">
              一缕
            </span>
            <span className="block truncate text-[11px] leading-tight text-ink-3">
              产品工作区
            </span>
          </span>
        </button>

        {/* accent action */}
        <button
          type="button"
          onClick={onNewTask}
          className="mb-2 flex w-full items-center gap-2 rounded-control px-2 py-1.5 text-[13px] font-medium text-accent transition-[background-color,transform] duration-100 hover:bg-accent-tint active:scale-[0.96]"
        >
          <span className="min-w-0 flex-1 truncate text-left">新建对话</span>
          <span className="flex size-4 shrink-0 items-center justify-center rounded-full bg-accent text-white">
            <Plus aria-hidden="true" size={9} strokeWidth={3} />
          </span>
        </button>

        {/* navigation sections */}
        <GlideMenu className="flex flex-col gap-2" highlightClassName="inset-x-0 rounded-[7px] bg-hover">
          <div>
            <SectionHeader label="视图" />
            <div className="flex flex-col gap-px">
              {(
                [
                  { key: "chat", label: "对话" },
                  { key: "approval", label: `审批${pendingCount ? ` (${pendingCount})` : ""}` },
                  { key: "palette", label: "命令" },
                  { key: "tasks", label: "任务" },
                  { key: "ball", label: "球状态" },
                  { key: "config", label: "设置" },
                  { key: "security", label: "安全" },
                  { key: "privacy", label: "隐私" },
                  { key: "cost", label: "成本" },
                ] as { key: PanelSidebarSection; label: string }[]
              ).map((item) => {
                const on = item.key === active;
                return (
                  <button
                    key={item.key}
                    type="button"
                    aria-current={on ? "page" : undefined}
                    aria-label={item.label}
                    data-menu-row
                    onClick={() => onSelect(item.key)}
                    className={cn(
                      "relative z-10 flex w-full items-center gap-2 rounded-[7px] px-2 py-1.5 text-left transition-[color,transform] duration-150 active:scale-[0.96]",
                      on ? "font-medium text-ink" : "text-ink-2",
                    )}
                  >
                    {on ? (
                      <span
                        aria-hidden="true"
                        className="absolute bottom-1.5 left-0.5 top-1.5 w-[2px] rounded-full bg-accent"
                      />
                    ) : null}
                    <span className="min-w-0 flex-1 truncate text-[13px]">{item.label}</span>
                  </button>
                );
              })}
            </div>
          </div>

          <div>
            <SectionHeader label="会话历史" />
            <p className="px-2 py-3 text-[12px] text-ink-3">
              还没有会话。面板不许自己造内容，Go 侧推送后会出现。
            </p>
          </div>
        </GlideMenu>
      </nav>

      {/* bottom block */}
      <GlideMenu
        className="flex flex-col gap-1 border-t border-line p-2"
        highlightClassName="inset-x-2 rounded-[7px] bg-hover"
      >
        <button
          type="button"
          data-menu-row
          onClick={onOpenSettings}
          className="flex w-full items-center gap-2 rounded-[7px] px-2 py-1.5 text-left text-[12.5px] text-ink-2 transition-colors duration-150 hover:text-ink"
        >
          <Settings aria-hidden="true" size={13} strokeWidth={1.8} />
          设置
        </button>
        <button
          type="button"
          data-menu-row
          onClick={onToggleTheme}
          className="flex w-full items-center gap-2 rounded-[7px] px-2 py-1.5 text-left text-[12.5px] text-ink-2 transition-colors duration-150 hover:text-ink"
        >
          {theme === "dark" ? (
            <Sun aria-hidden="true" size={13} strokeWidth={1.8} />
          ) : (
            <Moon aria-hidden="true" size={13} strokeWidth={1.8} />
          )}
          {theme === "dark" ? "浅色" : "暗色"}
        </button>
      </GlideMenu>
    </aside>
  );
}
