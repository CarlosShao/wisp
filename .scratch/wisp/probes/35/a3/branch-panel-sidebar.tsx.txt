/* ============================================================================
   Panel sidebar - the WORKSPACE SESSION TREE (2026-09-28, third cut).
   ----------------------------------------------------------------------------
   Two cuts of this file were wrong, and both are recorded because the third
   one is only correct if you know why the second one was not.

   Cut 1 (the icon rail's successor) put nine rows in a section. Settings
   appeared twice, and five of the nine were not navigation at all.

   Cut 2 fixed that by printing the same nine rows with nicer type. It was
   still wrong, and owner's words were: 「这不就是工作区会话的骨架树吗，怎么出来
   什么审批，任务，球状态」 - a sidebar's job is to show you YOUR WORK, not to
   be a menu of subsystems. So:

     审批  left the sidebar. It is a CARD that appears when a gated tool call
           needs you, plus a count on the header. A permanent row for
           something that is empty 99% of the time teaches you to ignore it.
     任务  left the sidebar for the top bar's 任务监控 popover, which is where
           a live view of running work belongs.
     球状态 left the panel entirely. The ball is a native Direct2D surface
           (D29) and has no business in the panel's navigation; its SETTINGS
           live on the settings page.

   What is left is the tree, and it is the shape the ?harness=1 walkthrough
   already established: a workspace row, the sessions under it, and recent
   sessions by time. Nothing else.

   Colours: token utilities only, zero literals; no emoji.
   ============================================================================ */

import { Fragment, useMemo, useState } from "react";
import {
  ChevronDown,
  ChevronRight,
  CircleAlert,
  Folder,
  Moon,
  Plus,
  Search,
  Settings,
  Sun,
  X,
} from "lucide-react";
import { GlideMenu } from "@/components/ai-native/glide-menu";
import { cn } from "@/lib/cn";

/** A session row's state. The icon vocabulary is the harness's own: red alert
    = refused, accent dot = streaming, nothing = done. */
export type SessionState = "done" | "failed" | "streaming";

export interface TreeSession {
  id: string;
  title: string;
  state: SessionState;
  /** Muted right-edge text (time or status), verbatim. */
  meta?: string;
}

export interface PanelSidebarProps {
  /** Which session the main column is showing, or null. */
  activeSessionId: string | null;
  onSelectSession: (id: string) => void;
  onNewTask: () => void;
  onOpenSettings: () => void;
  onToggleTheme: () => void;
  theme: "light" | "dark";
  settingsOpen: boolean;
  /** The workspace line. Path comes from the snapshot's C26-resolved field. */
  workspaceName: string;
  workspacePath: string;
  /** Sessions under the workspace row. Empty today (ticket 166 owns the feed). */
  workspaceSessions: readonly TreeSession[];
  /** Recent sessions by time bucket (今天 / 近 7 天 / 更早). */
  historyGroups: readonly { label: string; rows: readonly TreeSession[] }[];
}

function SectionHeader({ label }: { label: string }) {
  return (
    <div className="px-2 pb-1 pt-2 text-[10.5px] font-medium tracking-[0.08em] text-ink-3">
      {label}
    </div>
  );
}

/** The leading state glyph on a session row. A leaf slot is kept for `done` so
    every row's text starts at the same x - a column that jitters when a
    session changes state is unreadable. */
function StateGlyph({ state }: { state: SessionState }) {
  return (
    <span className="flex w-3.5 shrink-0 justify-center">
      {state === "failed" ? (
        <CircleAlert aria-hidden="true" className="shrink-0 text-red" size={12} strokeWidth={2} />
      ) : state === "streaming" ? (
        <span aria-hidden="true" className="block size-1.5 rounded-full bg-accent" />
      ) : null}
    </span>
  );
}

function SessionRow({
  session,
  active,
  onSelect,
  indent,
}: {
  session: TreeSession;
  active: boolean;
  onSelect: () => void;
  indent: number;
}) {
  return (
    <button
      type="button"
      aria-current={active ? "page" : undefined}
      aria-label={session.title}
      data-menu-row
      onClick={onSelect}
      style={{ paddingLeft: 8 + indent * 12 }}
      className={cn(
        // relative z-10: the GlideMenu highlight is an absolutely-positioned
        // sibling and paints OVER any row that does not raise itself. Omitting
        // it hid the 设置 label behind a grey bar on 2026-09-28.
        "relative z-10 flex w-full items-center gap-1.5 rounded-[7px] py-1.5 pr-2 text-left transition-[color,transform] duration-150 active:scale-[0.96]",
        active ? "font-medium text-ink" : "text-ink-2",
      )}
    >
      {/* Where you are is a FIXED accent bar, never the hover block: the
          highlight glides, and a marker that glides is a marker you lose. */}
      {active ? (
        <span
          aria-hidden="true"
          className="absolute bottom-1.5 left-1 top-1.5 w-[2px] rounded-full bg-accent"
        />
      ) : null}
      <StateGlyph state={session.state} />
      <span className="min-w-0 flex-1 truncate text-[13px]">{session.title}</span>
      {session.meta ? (
        <span className="shrink-0 text-[10.5px] tabular-nums text-ink-3">{session.meta}</span>
      ) : null}
    </button>
  );
}

export function PanelSidebar({
  activeSessionId,
  onSelectSession,
  onNewTask,
  onOpenSettings,
  onToggleTheme,
  theme,
  settingsOpen,
  workspaceName,
  workspacePath,
  workspaceSessions,
  historyGroups,
}: PanelSidebarProps) {
  const [query, setQuery] = useState("");
  const [workspaceOpen, setWorkspaceOpen] = useState(true);
  const [collapsed, setCollapsed] = useState<ReadonlySet<string>>(new Set());

  // The search box filters the tree under it. Memoised so the arrays keep their
  // identity: a projection built inline gets a fresh one every render, and any
  // effect keyed on it spins (this tree already hit a Maximum-update-depth loop
  // that way - see harness/sidebar.tsx).
  const q = query.trim().toLowerCase();
  const shownWorkspace = useMemo(
    () => (q ? workspaceSessions.filter((s) => s.title.toLowerCase().includes(q)) : workspaceSessions),
    [workspaceSessions, q],
  );
  const shownHistory = useMemo(
    () =>
      historyGroups
        .map((g) => ({
          ...g,
          rows: q ? g.rows.filter((s) => s.title.toLowerCase().includes(q)) : g.rows,
        }))
        .filter((g) => g.rows.length > 0),
    [historyGroups, q],
  );

  function toggleGroup(label: string) {
    setCollapsed((current) => {
      const next = new Set(current);
      if (next.has(label)) next.delete(label);
      else next.add(label);
      return next;
    });
  }

  const nothingToShow = shownWorkspace.length === 0 && shownHistory.length === 0;

  return (
    <aside className="flex h-full w-60 min-w-60 shrink-0 flex-col border-r border-line bg-surface">
      <nav aria-label="工作区会话" className="flex min-h-0 flex-1 flex-col overflow-y-auto p-2">
        {/* workspace row: the root of the tree */}
        <button
          type="button"
          title={workspacePath}
          className="mb-1.5 flex w-full items-center gap-2.5 rounded-control p-1.5 text-left transition-[background-color,transform] duration-100 hover:bg-hover active:scale-[0.96]"
        >
          <span className="flex size-8 shrink-0 items-center justify-center rounded-[8px] bg-ink text-[13px] font-semibold text-surface">
            {/* W, not 一: a lone horizontal stroke in a dark tile reads as a
                minus sign, which is the one thing a workspace chip must not
                look like. */}
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
          <ChevronDown aria-hidden="true" className="shrink-0 rotate-180 text-ink-3" size={12} strokeWidth={2} />
        </button>

        {/* search: filters the tree. The clear button replaces the "/" chip
            while there is something to clear - a key hint that never gives way
            teaches people to stop reading the box. */}
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

        <button
          type="button"
          onClick={onNewTask}
          className="mb-1 flex w-full items-center gap-2 rounded-control px-2 py-1.5 text-[13px] font-medium text-accent transition-[background-color,transform] duration-100 hover:bg-accent-tint active:scale-[0.96]"
        >
          <span className="min-w-0 flex-1 truncate text-left">新建任务</span>
          <span className="flex size-4 shrink-0 items-center justify-center rounded-full bg-accent text-white">
            <Plus aria-hidden="true" size={9} strokeWidth={3} />
          </span>
        </button>

        <GlideMenu className="flex flex-col gap-0.5" highlightClassName="inset-x-0 rounded-[7px] bg-hover">
          {/* ── 工作区：项目行 + 挂在它下面的会话 ────────────────────────── */}
          <div>
            <SectionHeader label="工作区" />
            <button
              type="button"
              aria-expanded={workspaceOpen}
              aria-label={workspaceName}
              data-menu-row
              onClick={() => setWorkspaceOpen((v) => !v)}
              className="relative z-10 flex w-full items-center gap-1.5 rounded-[7px] px-2 py-1.5 text-left transition-[color,transform] duration-150 active:scale-[0.96]"
            >
              <ChevronRight
                aria-hidden="true"
                className={cn("shrink-0 text-ink-3 transition-transform duration-150", workspaceOpen && "rotate-90")}
                size={10}
                strokeWidth={2.2}
              />
              <Folder aria-hidden="true" className="shrink-0 text-ink-2" size={13} strokeWidth={1.8} />
              <span className="min-w-0 flex-1 truncate text-[13px] text-ink-2">{workspaceName}</span>
              {shownWorkspace.length > 0 ? (
                <span className="shrink-0 text-[10.5px] tabular-nums text-ink-3">{shownWorkspace.length}</span>
              ) : null}
            </button>
            {workspaceOpen && shownWorkspace.length > 0 ? (
              <div className="flex flex-col gap-px">
                {shownWorkspace.map((session) => (
                  <SessionRow
                    active={session.id === activeSessionId}
                    indent={1}
                    key={session.id}
                    onSelect={() => onSelectSession(session.id)}
                    session={session}
                  />
                ))}
              </div>
            ) : null}
            {/* The empty line lives UNDER the project row, where a file tree
                puts it. It was floating between the row and the next section,
                which read as a status message about the sidebar rather than as
                "this folder has nothing in it yet". */}
            {workspaceOpen && shownWorkspace.length === 0 ? (
              <p className="py-1 pl-6 pr-2 text-[11.5px] leading-[1.6] text-ink-3">
                {q ? "没有匹配的会话。" : "这个工作区下还没有会话。"}
              </p>
            ) : null}
          </div>

          {/* ── 最近任务：按时间分组，可折叠 ─────────────────────────────── */}
          {shownHistory.map((group) => {
            const shut = collapsed.has(group.label);
            return (
              <Fragment key={group.label}>
                <SectionHeader label={group.label} />
                <button
                  type="button"
                  aria-expanded={!shut}
                  onClick={() => toggleGroup(group.label)}
                  className="flex w-full items-center gap-1 px-2 pb-0.5 text-[10.5px] font-medium tracking-[0.06em] text-ink-3"
                >
                  <ChevronDown
                    aria-hidden="true"
                    className={cn("shrink-0 transition-transform duration-150", shut && "-rotate-90")}
                    size={10}
                    strokeWidth={2.2}
                  />
                  <span className="min-w-0 flex-1 truncate text-left">{group.label}</span>
                  <span className="shrink-0 tabular-nums">{group.rows.length}</span>
                </button>
                {!shut && (
                  <div className="flex flex-col gap-px">
                    {group.rows.map((session) => (
                      <SessionRow
                        active={session.id === activeSessionId}
                        indent={1}
                        key={session.id}
                        onSelect={() => onSelectSession(session.id)}
                        session={session}
                      />
                    ))}
                  </div>
                )}
              </Fragment>
            );
          })}
        </GlideMenu>

        {nothingToShow && q ? (
          <p className="px-3 py-2 text-[12px] leading-[1.6] text-ink-3">
            没有匹配「{query.trim()}」的会话。
          </p>
        ) : null}
      </nav>

      {/* bottom: settings and theme. 设置 appears exactly ONCE in this whole
          sidebar - the duplicate shipped on 2026-09-27 and the gate now pins
          the count. */}
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
            <span aria-hidden="true" className="absolute bottom-1.5 left-0.5 top-1.5 w-[2px] rounded-full bg-accent" />
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
