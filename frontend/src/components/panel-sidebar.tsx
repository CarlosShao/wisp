/* ============================================================================
   Panel sidebar (product path, rebuilt 2026-09-28 after the first attempt
   translated the old 9-row icon rail into 9 rows of text and called it a
   design).
   ----------------------------------------------------------------------------
   What the first attempt got wrong, recorded so it cannot come back:

     1. 设置 appeared TWICE - once as a row in the section list, once at the
        bottom. The same door in two places is not a choice, it is a bug.
     2. The section list carried all NINE legacy rail rows. Five of them
        (命令 / 设置 / 安全 / 隐私 / 成本) are not views a person navigates
        to: 命令 is the Ctrl+K palette, and the other four are settings
        CATEGORIES, which live on the settings page. Putting them in the
        sidebar made the sidebar a transcription of a retired icon strip
        instead of a shape with a job.
     3. No search field, even though the shape owner approved has one.
     4. 会话历史 was a sentence pretending to be a section.

   The shape now, top to bottom, is the approved one:

     workspace chip -> search -> 新建对话 -> 面板 section (the four places
     the main column can go) -> 会话历史 section (time groups) -> a hairline
     -> 设置 / 主题.

   Colours: token utilities only, zero literals; no emoji.
   ============================================================================ */

import { Fragment, useMemo, useState } from "react";
import { ChevronDown, Moon, Plus, Search, Settings, Sun, X } from "lucide-react";
import { GlideMenu } from "@/components/ai-native/glide-menu";
import { cn } from "@/lib/cn";

/** The four places the main column can go. Each is a SECTION of the sidebar,
    not a page of its own: clicking one swaps the main column, it does not
    open a window. */
export type PanelSectionId = "chat" | "approval" | "tasks" | "ball";

/** One 会话历史 row. Empty until Go pushes a session list (ticket 166), which
    is why the list is a prop and not a constant in this file. */
export interface PanelSessionRow {
  id: string;
  title: string;
  /** Muted right-edge text (time or state), verbatim. */
  meta?: string;
}

export interface PanelSessionGroup {
  /** 今天 / 近 7 天 / 更早. */
  label: string;
  rows: readonly PanelSessionRow[];
}

export interface PanelSidebarProps {
  /** Which of the four the main column is showing. */
  active: PanelSectionId;
  onSelect: (id: PanelSectionId) => void;
  onOpenSettings: () => void;
  onToggleTheme: () => void;
  theme: "light" | "dark";
  /** True while the settings page is up, so the bottom row lights as the
      place you actually are. */
  settingsOpen: boolean;
  /** The workspace row's two lines. Both come from the snapshot's composer
      section; when Go has not filled it the caller passes the reason. */
  workspaceName: string;
  workspacePath: string;
  /** Native approval queue depth, or 0. Drives the 审批 badge. */
  pendingCount: number;
  /** Time-grouped session history. Empty array is the honest empty state. */
  sessionGroups: readonly PanelSessionGroup[];
  activeSessionId: string | null;
  onSelectSession: (id: string) => void;
}

/** 面板 section 的四行。写死在这里而不是从 panel-views 过滤，是有意的：
    那张表是「屏的词汇表」（门禁按它查图标名），不是「侧栏该显示几行」。
    两者混过一次，代价是设置出现在两个地方。 */
const SECTION_ROWS: readonly { id: PanelSectionId; label: string; hint: string }[] = [
  { id: "chat", label: "对话", hint: "当前这次对话" },
  { id: "approval", label: "审批", hint: "等你点头的操作" },
  { id: "tasks", label: "任务", hint: "正在跑和排队的活" },
  { id: "ball", label: "球状态", hint: "悬浮球现在在哪个态" },
];

function SectionHeader({ label, action }: { label: string; action?: React.ReactNode }) {
  return (
    <div className="flex items-center gap-1 px-2 pb-1 pt-2 text-[10.5px] font-medium tracking-[0.08em] text-ink-3">
      <span className="min-w-0 flex-1 truncate">{label}</span>
      {action}
    </div>
  );
}

export function PanelSidebar({
  active,
  onSelect,
  onOpenSettings,
  onToggleTheme,
  theme,
  settingsOpen,
  workspaceName,
  workspacePath,
  pendingCount,
  sessionGroups,
  activeSessionId,
  onSelectSession,
}: PanelSidebarProps) {
  const [query, setQuery] = useState("");
  const [collapsed, setCollapsed] = useState<ReadonlySet<string>>(new Set());

  // The search box filters the history below it. Memoised so the array keeps
  // its identity across renders: a projection built inline gets a fresh
  // identity every render, and any effect keyed on it spins (this tree hit a
  // Maximum-update-depth loop that way once already - see sidebar.tsx).
  const q = query.trim().toLowerCase();
  const shownGroups = useMemo(
    () =>
      sessionGroups
        .map((group) => ({
          ...group,
          rows: q ? group.rows.filter((row) => row.title.toLowerCase().includes(q)) : group.rows,
        }))
        .filter((group) => group.rows.length > 0),
    [sessionGroups, q],
  );

  function toggleGroup(label: string) {
    setCollapsed((current) => {
      const next = new Set(current);
      if (next.has(label)) next.delete(label);
      else next.add(label);
      return next;
    });
  }

  const historyEmpty = shownGroups.length === 0;

  return (
    <aside className="flex h-full w-60 min-w-60 shrink-0 flex-col border-r border-line bg-surface">
      <nav
        aria-label="面板导航"
        className="flex min-h-0 flex-1 flex-col overflow-y-auto p-2"
      >
        {/* 工作区行：名字 + 路径。路径是快照里 C26 授权后的那条，不是用户
            输入框里还没提交的那条。 */}
        <button
          title={workspacePath}
          type="button"
          className="mb-2 flex w-full items-center gap-2.5 rounded-control p-1.5 text-left transition-[background-color,transform] duration-100 hover:bg-hover active:scale-[0.96]"
        >
          <span className="flex size-8 shrink-0 items-center justify-center rounded-[8px] bg-ink text-[13px] font-semibold text-surface">
            {/* The monogram is W, not the first character of 一缕: a lone
                horizontal stroke in a dark tile reads as a minus sign, which
                is the one thing a workspace chip must not look like. */}
            W
          </span>
          <span className="min-w-0 flex-1">
            <span className="block truncate text-[13px] font-medium leading-tight text-ink">
              {workspaceName}
            </span>
            <span className="block truncate font-mono text-[10.5px] leading-tight text-ink-3">
              {workspacePath}
            </span>
          </span>
        </button>

        {/* 搜索：过滤下面的会话历史。有内容时让位给清除键，否则给 "/" 角标
            —— 一个角标承诺的键，必须真的有监听器接它。 */}
        <label className="mb-1 flex h-8 items-center gap-2 rounded-control bg-inset px-2.5 shadow-hairline">
          <Search aria-hidden="true" className="shrink-0 text-ink-3" size={12} strokeWidth={2} />
          <input
            aria-label="搜索会话"
            onChange={(event) => setQuery(event.target.value)}
            placeholder="搜索会话"
            value={query}
            className="min-w-0 flex-1 bg-transparent text-[12.5px] text-ink outline-none placeholder:text-ink-3"
          />
          {query ? (
            <button
              type="button"
              aria-label="清空搜索"
              onClick={() => setQuery("")}
              className="flex size-4.5 shrink-0 items-center justify-center rounded-[5px] text-ink-3 transition-colors duration-100 hover:bg-hover hover:text-ink"
            >
              <X aria-hidden="true" size={11} strokeWidth={2.4} />
            </button>
          ) : (
            <kbd className="flex size-4.5 shrink-0 items-center justify-center rounded-[5px] bg-surface text-[10px] text-ink-3 shadow-hairline">
              /
            </kbd>
          )}
        </label>

        {/* 主行动：新建对话。整页 harness 的 accent 行 + 加号徽章圆。 */}
        <button
          type="button"
          onClick={() => onSelect("chat")}
          className="mb-1 flex w-full items-center gap-2 rounded-control px-2 py-1.5 text-[13px] font-medium text-accent transition-[background-color,transform] duration-100 hover:bg-accent-tint active:scale-[0.96]"
        >
          <span className="min-w-0 flex-1 truncate text-left">新建对话</span>
          <span className="flex size-4 shrink-0 items-center justify-center rounded-full bg-accent text-white">
            <Plus aria-hidden="true" size={9} strokeWidth={3} />
          </span>
        </button>

        <GlideMenu
          className="flex flex-col gap-0.5"
          highlightClassName="inset-x-0 rounded-[7px] bg-hover"
        >
          {/* ── 面板区段：主栏能去的四个地方 ─────────────────────────────── */}
          <div>
            <SectionHeader label="面板" />
            <div className="flex flex-col gap-px">
              {SECTION_ROWS.map((row) => {
                const on = row.id === active;
                return (
                  <button
                    key={row.id}
                    type="button"
                    aria-current={on ? "page" : undefined}
                    aria-label={row.label}
                    title={row.hint}
                    data-menu-row
                    onClick={() => onSelect(row.id)}
                    className={cn(
                      "relative z-10 flex w-full items-center gap-2 rounded-[7px] py-1.5 pr-2 text-left transition-[color,transform] duration-150 active:scale-[0.96]",
                      on ? "pl-3.5 font-medium text-ink" : "pl-3.5 text-ink-2",
                    )}
                  >
                    {/* 你在这一行：一根不动的 accent 左标。hover 滑块会滑走，
                        「我在这」不能跟着滑走。 */}
                    {on ? (
                      <span
                        aria-hidden="true"
                        className="absolute bottom-1.5 left-1 top-1.5 w-[2px] rounded-full bg-accent"
                      />
                    ) : null}
                    <span className="min-w-0 flex-1 truncate text-[13px]">{row.label}</span>
                    {row.id === "approval" && pendingCount > 0 ? (
                      <span className="shrink-0 rounded-full bg-accent-tint px-1.5 py-px text-[10.5px] font-medium tabular-nums text-accent-ink">
                        {pendingCount}
                      </span>
                    ) : null}
                  </button>
                );
              })}
            </div>
          </div>

          {/* ── 会话历史区段：按时间分组，可折叠；空时是一句实话 ─────────── */}
          <div>
            <SectionHeader label="会话历史" />
            {historyEmpty ? (
              <p className="px-3 py-2 text-[12px] leading-[1.6] text-ink-3">
                {q ? `没有匹配「${query.trim()}」的会话` : "还没有会话"}
                {!q ? "。面板不自己造会话，原生侧推送后这里才会有内容。" : "。"}
              </p>
            ) : (
              shownGroups.map((group) => {
                const shut = collapsed.has(group.label);
                return (
                  <Fragment key={group.label}>
                    <button
                      type="button"
                      aria-expanded={!shut}
                      onClick={() => toggleGroup(group.label)}
                      className="flex w-full items-center gap-1 px-2 pb-0.5 pt-1 text-[10.5px] font-medium tracking-[0.06em] text-ink-3"
                    >
                      <ChevronDown
                        aria-hidden="true"
                        className={cn(
                          "shrink-0 transition-transform duration-150",
                          shut && "-rotate-90",
                        )}
                        size={10}
                        strokeWidth={2.2}
                      />
                      <span className="min-w-0 flex-1 truncate text-left">{group.label}</span>
                    </button>
                    {!shut && (
                      <div className="flex flex-col gap-px">
                        {group.rows.map((row) => {
                          const on = row.id === activeSessionId;
                          return (
                            <button
                              key={row.id}
                              type="button"
                              aria-current={on ? "page" : undefined}
                              data-menu-row
                              onClick={() => onSelectSession(row.id)}
                              className={cn(
                                "relative z-10 flex w-full items-center gap-2 rounded-[7px] py-1.5 pr-2 text-left transition-[color,transform] duration-150 active:scale-[0.96]",
                                on ? "pl-3.5 font-medium text-ink" : "pl-3.5 text-ink-2",
                              )}
                            >
                              {on ? (
                                <span
                                  aria-hidden="true"
                                  className="absolute bottom-1.5 left-1 top-1.5 w-[2px] rounded-full bg-accent"
                                />
                              ) : null}
                              <span className="min-w-0 flex-1 truncate text-[13px]">{row.title}</span>
                              {row.meta ? (
                                <span className="shrink-0 text-[10.5px] tabular-nums text-ink-3">
                                  {row.meta}
                                </span>
                              ) : null}
                            </button>
                          );
                        })}
                      </div>
                    )}
                  </Fragment>
                );
              })
            )}
          </div>
        </GlideMenu>
      </nav>

      {/* ── 底部：设置与主题。设置在这张侧栏里只出现这一次 ────────────────
          The rows carry relative z-10 for the same reason the nav rows do:
          the GlideMenu highlight block is an absolutely-positioned sibling, and
          without the z-index it paints OVER the label. That shipped as a grey
          bar where 设置 should have been - the button was still clickable,
          which is exactly why nobody noticed. */}
      <GlideMenu
        className="flex flex-col gap-1 border-t border-line p-2"
        highlightClassName="inset-x-2 rounded-[7px] bg-hover"
      >
        <button
          type="button"
          aria-label="设置"
          aria-current={settingsOpen ? "page" : undefined}
          data-menu-row
          onClick={onOpenSettings}
          className={cn(
            "relative z-10 flex w-full items-center gap-2 rounded-[7px] px-2 py-1.5 text-left text-[12.5px] transition-colors duration-150",
            settingsOpen ? "font-medium text-ink" : "text-ink-2 hover:text-ink",
          )}
        >
          {settingsOpen ? (
            <span
              aria-hidden="true"
              className="absolute bottom-1.5 left-0.5 top-1.5 w-[2px] rounded-full bg-accent"
            />
          ) : null}
          <Settings aria-hidden="true" size={13} strokeWidth={1.8} />
          设置
        </button>
        <button
          type="button"
          data-menu-row
          onClick={onToggleTheme}
          className="relative z-10 flex w-full items-center gap-2 rounded-[7px] px-2 py-1.5 text-left text-[12.5px] text-ink-2 transition-colors duration-150 hover:text-ink"
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
