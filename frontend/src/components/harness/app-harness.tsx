/* ============================================================================
   AppHarness - the harness as a WHOLE application (owner 2026-09-26: "用户刚
   进来的时候，不是都进入一个新任务的窗口吗？左边侧边栏会有新建任务等…右侧
   边栏等等…一个 harness 前端正常应该有什么页面，你也得弄全了").
   ----------------------------------------------------------------------------
   This is the ?harness=1 page: a full-application walkthrough in the shape
   mainstream agent harnesses use - sidebar (new task, search, history) on the
   left, main column (new-task entry / session view) in the middle, a
   collapsible context rail (approvals / tasks / usage) on the right, and the
   overlay family (command palette, settings, workspace picker, first-run,
   toasts) above everything. main.tsx routes ?harness=1 here and ?harness=2 to
   the component showcase.

   What is real: every component. Sidebar, main, right rail and overlays are
   props-driven files the panel reuses or mirrors; the L2 card, config screen,
   firstrun screen and palette are the SAME files production mounts. What is
   demo: every value - fixtures/harness-app.ts and the RB_* consts carry the
   "demo 数据" burden, and the banner from main.tsx says so.

   State policy: this page is a walkthrough, so it may hold view state (which
   session, which overlay) the way showcase.tsx does - none of it persists,
   none of it reaches the bridge. The product panel (App) stays stateless.
   ============================================================================ */

import { useEffect, useState } from "react";
import { Command, Moon, PanelRight, PanelRightClose, Sun } from "lucide-react";
import { EntityChip } from "@/components/ai-native/entity-chip";
import { PaletteScreen } from "@/components/palette-screen";
import { FirstrunOverlay, RB_FIRSTRUN } from "@/components/harness/firstrun-overlay";
import { HarnessMain } from "@/components/harness/main";
import { HarnessSidebar } from "@/components/harness/sidebar";
import { RightRail, RB_RIGHT_RAIL } from "@/components/harness/right-rail";
import { SettingsModal } from "@/components/harness/settings-modal";
import { ToastStack, RB_TOASTS } from "@/components/harness/toast";
import { WorkspacePicker, RB_WORKSPACES } from "@/components/harness/workspace-picker";
import {
  AUTOMATIONS,
  GREETING,
  EMPTY_HINT,
  HISTORY_GROUPS,
  PROJECT,
  QUICK_COMMANDS,
  RECENT_TASKS,
  SESSIONS,
  WEEK_USAGE,
  WORKSPACE_NOW,
} from "@/fixtures/harness-app";
import { SHOWCASE_APPROVALS, SHOWCASE_PALETTE_GROUPS } from "@/fixtures/harness";

export function AppHarness() {
  // view state: which main column mode, which session, which overlays.
  const [sessionId, setSessionId] = useState<string | null>(null);
  const [search, setSearch] = useState("");
  const [railOpen, setRailOpen] = useState(true);
  const [paletteOpen, setPaletteOpen] = useState(false);
  const [paletteQuery, setPaletteQuery] = useState("");
  const [pickerOpen, setPickerOpen] = useState(false);
  const [pickerQuery, setPickerQuery] = useState("");
  const [settingsOpen, setSettingsOpen] = useState(false);
  const [firstrunOpen, setFirstrunOpen] = useState(true);
  const [dark, setDark] = useState(false);

  function toggleTheme() {
    const next = !dark;
    setDark(next);
    document.documentElement.dataset.theme = next ? "dark" : "light";
  }

  // Ctrl+K opens the command palette - the one keyboard verb every harness
  // shares. Esc closes whichever overlay is up, topmost first.
  useEffect(() => {
    function onKey(e: KeyboardEvent) {
      if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === "k") {
        e.preventDefault();
        setPaletteOpen((v) => !v);
        return;
      }
      if (e.key === "Escape") {
        if (paletteOpen) setPaletteOpen(false);
        else if (pickerOpen) setPickerOpen(false);
        else if (settingsOpen) setSettingsOpen(false);
      }
    }
    document.addEventListener("keydown", onKey);
    return () => document.removeEventListener("keydown", onKey);
  }, [paletteOpen, pickerOpen, settingsOpen]);

  const activeSession = sessionId ? (SESSIONS.find((s) => s.id === sessionId) ?? null) : null;
  const queueRows = SHOWCASE_APPROVALS.map((a) => ({
    level: a.level,
    tool: a.tool,
    correlationId: a.correlationId,
  }));

  function newTask() {
    setSessionId(null);
    setFirstrunOpen(false);
  }

  return (
    <div className="flex h-screen w-full overflow-hidden bg-page text-ink">
      <HarnessSidebar
        activeSessionId={sessionId}
        automations={AUTOMATIONS}
        groups={HISTORY_GROUPS}
        onNewTask={newTask}
        onOpenSettings={() => setSettingsOpen(true)}
        onSelectSession={(id) => setSessionId(id)}
        onToggleTheme={toggleTheme}
        project={PROJECT}
        searchValue={search}
        onSearchChange={setSearch}
        theme={dark ? "dark" : "light"}
        usage={WEEK_USAGE}
        workspace={{ name: "Wisp", path: WORKSPACE_NOW }}
      />

      <div className="flex min-w-0 flex-1 flex-col">
        {/* 顶栏：工作区（点击换）+ 命令面板 + 右栏开关 —— 应用骨架的第三块 */}
        <header className="flex h-11 shrink-0 items-center justify-between gap-3 border-b border-line bg-surface px-4">
          <button
            className="flex min-w-0 items-center gap-2 rounded-control px-1.5 py-1 transition-colors duration-150 hover:bg-hover"
            onClick={() => setPickerOpen(true)}
            title="选择工作区"
            type="button"
          >
            <EntityChip color="var(--accent)" name={WORKSPACE_NOW.slice(0, 2)} />
            <span className="truncate font-mono text-[11.5px] text-ink-2">{WORKSPACE_NOW}</span>
          </button>
          <div className="flex shrink-0 items-center gap-1.5">
            <button
              className="flex items-center gap-1.5 rounded-control px-2 py-1 text-[11.5px] text-ink-2 transition-colors duration-150 hover:bg-hover hover:text-ink"
              onClick={() => setPaletteOpen(true)}
              type="button"
            >
              <Command aria-hidden="true" size={13} />
              <span>命令面板</span>
              <kbd className="kbd">Ctrl K</kbd>
            </button>
            <button
              aria-label={railOpen ? "收起右侧栏" : "展开右侧栏"}
              className="flex size-7 items-center justify-center rounded-control text-ink-2 transition-colors duration-150 hover:bg-hover hover:text-ink"
              onClick={() => setRailOpen((v) => !v)}
              title={railOpen ? "收起右侧栏" : "展开右侧栏"}
              type="button"
            >
              {railOpen ? <PanelRightClose aria-hidden="true" size={15} /> : <PanelRight aria-hidden="true" size={15} />}
            </button>
            <button
              aria-label="切换主题"
              className="flex size-7 items-center justify-center rounded-control text-ink-2 transition-colors duration-150 hover:bg-hover hover:text-ink"
              onClick={toggleTheme}
              title="切换主题"
              type="button"
            >
              {dark ? <Sun aria-hidden="true" size={15} /> : <Moon aria-hidden="true" size={15} />}
            </button>
          </div>
        </header>

        <div className="flex min-h-0 flex-1">
          <HarnessMain
            emptyHint={EMPTY_HINT}
            greeting={GREETING}
            onNewTask={newTask}
            onPickCommand={() => setPaletteOpen(true)}
            onPickTask={() => setSessionId(SESSIONS[0]?.id ?? null)}
            onPickWorkspace={() => setPickerOpen(true)}
            quickCommands={QUICK_COMMANDS}
            recentTasks={RECENT_TASKS}
            session={activeSession}
            view={activeSession ? "session" : "new"}
            workspace={WORKSPACE_NOW}
          />
          {railOpen && (
            <RightRail
              onViewAllApprovals={() => undefined}
              pendingCount={SHOWCASE_APPROVALS.length}
              queueRows={queueRows}
              runs={RB_RIGHT_RAIL.runs}
            />
          )}
        </div>
      </div>

      {/* 浮层族：命令面板 / 工作区选择 / 设置 / 首启 / Toast */}
      {paletteOpen && (
        <div
          className="fixed inset-0 z-40 flex items-start justify-center bg-black/35 pt-[12vh]"
          onClick={(e) => {
            if (e.target === e.currentTarget) setPaletteOpen(false);
          }}
          role="presentation"
        >
          <div className="w-[560px] overflow-hidden rounded-card border border-line bg-surface shadow-overlay" style={{ animation: "fade-up 200ms var(--ease-out-strong) both" }}>
            <PaletteScreen
              groups={SHOWCASE_PALETTE_GROUPS}
              onQueryChange={setPaletteQuery}
              query={paletteQuery}
            />
          </div>
        </div>
      )}
      {pickerOpen && (
        <WorkspacePicker
          current={RB_WORKSPACES.current}
          onBrowse={() => undefined}
          onClose={() => setPickerOpen(false)}
          onPick={() => setPickerOpen(false)}
          query={pickerQuery}
          recent={RB_WORKSPACES.recent}
          onQueryChange={setPickerQuery}
        />
      )}
      <SettingsModal onClose={() => setSettingsOpen(false)} open={settingsOpen} />
      {firstrunOpen && <FirstrunOverlay {...RB_FIRSTRUN} onClose={() => setFirstrunOpen(false)} />}
      <ToastStack items={RB_TOASTS} />
    </div>
  );
}
