/* ============================================================================
   HarnessSidebar - the harness APP page's left column (showcase demo data,
   ?harness=1 only; the product path never mounts this).
   ----------------------------------------------------------------------------
   Owner's shape ruling (2026-09-26): the app's left skeleton is a ~240px
   inset column - app identity on top, a primary "new task" button, a search
   field, grouped session history, and at the bottom the settings entry, the
   weekly usage subtotal and the theme toggle.

   PROPS CONTRACT (the assembler wires every one of these; no data lives in
   the component - all rows, pills and labels arrive via props):

     groups: HarnessSidebarGroup[]
         Session history, already grouped (今天 / 近 7 天 / 更早) and already
         flattened to rows by src/fixtures/harness-app.ts (toSidebarRow).
         Rows render top-down in array order; empty groups still print their
         label, an all-empty list renders the "没有匹配的会话" empty state.
     activeSessionId: string | null
         The highlighted row gets bg-hover plus a 2px accent bar on the left
         edge. null highlights nothing.
     onSelectSession: (id: string) => void
         Row click. Fires with the row's id.
     onNewTask: () => void
         The primary accent button under the identity row ("新建任务").
     searchValue: string
     onSearchChange: (value: string) => void
         Controlled search field (rounded-control, bg-field, Search icon).
     searchPlaceholder?: string
         Defaults to "搜索会话…".
     usage: HarnessUsagePill[]
         Bottom subtotal - each entry renders "{label} <ValuePill>{value}</ValuePill>".
         Two entries in the fixture (本周 tokens / 本周花费); any count renders.
     onOpenSettings: () => void
         The settings row (Settings icon + text).
     theme: "light" | "dark"
     onToggleTheme: () => void
         Theme row: shows the icon of the theme you would switch TO
         (dark -> Sun "浅色", light -> Moon "暗色"), showcase's convention.
     className?: string
         Appended to the root aside.

   Interactions are callbacks only - no storage API, no bridge reads. All
   colours are token utilities; zero literals.
   ============================================================================ */

import { Moon, Search, Settings, Sun, Plus } from "lucide-react";
import { cn } from "@/lib/cn";
import { ValuePill } from "@/components/ai-native/value-pill";
import type {
  HarnessBadgeTone,
  HarnessSidebarGroup,
  HarnessSidebarProps,
} from "@/fixtures/harness-app";

/* 次行状态徽标的缩微色表（tint 底 + 语义前景，token 工具类）。 */
const BADGE_TONES: Record<HarnessBadgeTone, string> = {
  green: "bg-green-tint text-green",
  accent: "bg-accent-tint text-accent-ink",
  red: "bg-red-tint text-red",
  neutral: "bg-field text-ink-2",
};

export function HarnessSidebar({
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
  const totalRows = groups.reduce((sum, g) => sum + g.rows.length, 0);

  return (
    <aside
      className={cn("flex h-full w-60 shrink-0 flex-col border-r border-line bg-inset", className)}
    >
      {/* 应用标识行：青点 + 「一缕」 */}
      <div className="flex items-center gap-2 px-4 pb-3 pt-4">
        <span aria-hidden="true" className="size-3 shrink-0 rounded-full bg-accent" />
        <span className="text-[13px] font-semibold text-ink">一缕</span>
      </div>

      {/* 新建任务主按钮 */}
      <div className="px-3">
        <button
          className="flex w-full items-center justify-center gap-1.5 rounded-control bg-accent
            px-3 py-2 text-[13px] font-medium text-white shadow-btn
            transition-[background-color,transform] duration-150 hover:bg-accent-ink active:scale-[0.98]"
          onClick={onNewTask}
          style={{ animation: "fade-up 300ms var(--ease-out-strong) both" }}
          type="button"
        >
          <Plus aria-hidden="true" className="size-4" strokeWidth={2.2} />
          新建任务
        </button>
      </div>

      {/* 搜索框（受控） */}
      <div className="px-3 pt-2.5">
        <div
          className="flex h-8 items-center gap-1.5 rounded-control border border-transparent
            bg-field px-2.5 transition-[border-color] duration-150 focus-within:border-line-strong"
        >
          <Search aria-hidden="true" className="size-3.5 shrink-0 text-ink-3" strokeWidth={2} />
          <input
            aria-label="搜索会话"
            className="min-w-0 flex-1 bg-transparent text-[12.5px] text-ink outline-none placeholder:text-ink-3"
            onChange={(e) => onSearchChange(e.target.value)}
            placeholder={searchPlaceholder}
            value={searchValue}
          />
        </div>
      </div>

      {/* 会话历史列表 */}
      <nav aria-label="会话历史" className="mt-3 min-h-0 flex-1 overflow-y-auto px-3 pb-2">
        {totalRows === 0 ? (
          <p className="px-2 py-6 text-center text-[11.5px] leading-relaxed text-ink-3">
            没有匹配的会话
          </p>
        ) : (
          groups.map((group) => (
            <SidebarGroupBlock
              activeSessionId={activeSessionId}
              group={group}
              key={group.label}
              onSelectSession={onSelectSession}
            />
          ))
        )}
      </nav>

      {/* 底部：设置入口 / 用量小计 / 主题切换 */}
      <div className="flex flex-col gap-1 border-t border-line px-3 py-3">
        <button
          className="flex items-center gap-2 rounded-control px-2 py-1.5 text-[12.5px] text-ink-2
            transition-colors duration-100 hover:bg-hover hover:text-ink"
          onClick={onOpenSettings}
          type="button"
        >
          <Settings aria-hidden="true" className="size-4 shrink-0" strokeWidth={2} />
          设置
        </button>

        <div className="flex flex-wrap items-center gap-x-1 gap-y-1 px-2 py-1">
          {usage.map((pill) => (
            <span className="inline-flex items-center gap-1" key={pill.label}>
              <span className="text-[10.5px] text-ink-3">{pill.label}</span>
              <ValuePill tone={pill.tone}>{pill.value}</ValuePill>
            </span>
          ))}
        </div>

        <button
          className="flex items-center gap-2 rounded-control px-2 py-1.5 text-[12.5px] text-ink-2
            transition-colors duration-100 hover:bg-hover hover:text-ink"
          onClick={onToggleTheme}
          type="button"
        >
          {theme === "dark" ? (
            <Sun aria-hidden="true" className="size-4 shrink-0" strokeWidth={2} />
          ) : (
            <Moon aria-hidden="true" className="size-4 shrink-0" strokeWidth={2} />
          )}
          {theme === "dark" ? "浅色" : "暗色"}
        </button>
      </div>
    </aside>
  );
}

/* One labelled group of history rows. */
function SidebarGroupBlock({
  group,
  activeSessionId,
  onSelectSession,
}: {
  group: HarnessSidebarGroup;
  activeSessionId: string | null;
  onSelectSession: (id: string) => void;
}) {
  return (
    <div className="mb-3">
      <p className="px-1.5 pb-1 text-[11px] uppercase tracking-wide text-ink-3">{group.label}</p>
      {group.rows.map((row) => {
        const active = row.id === activeSessionId;
        return (
          <button
            className={cn(
              "relative mb-0.5 flex w-full items-start gap-2 rounded-control px-2.5 py-2 text-left",
              "transition-colors duration-100",
              active ? "bg-hover" : "hover:bg-hover",
            )}
            key={row.id}
            onClick={() => onSelectSession(row.id)}
            type="button"
          >
            {active ? (
              <span
                aria-hidden="true"
                className="absolute bottom-1.5 left-0 top-1.5 w-0.5 rounded-full bg-accent"
              />
            ) : null}
            <span className="min-w-0 flex-1">
              <span
                className={cn(
                  "block truncate text-[12.5px] leading-5",
                  active ? "font-medium text-ink" : "text-ink",
                )}
              >
                {row.title}
              </span>
              <span className="mt-0.5 flex items-center gap-1.5">
                {row.badge ? (
                  <span
                    className={cn(
                      "inline-flex items-center rounded-chip px-1 text-[10px] font-medium leading-[16px]",
                      BADGE_TONES[row.badge.tone],
                    )}
                  >
                    {row.badge.label}
                  </span>
                ) : null}
                {row.time ? (
                  <span className="truncate text-[10.5px] text-ink-3">{row.time}</span>
                ) : null}
              </span>
            </span>
          </button>
        );
      })}
    </div>
  );
}
