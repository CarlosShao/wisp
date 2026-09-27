/* ============================================================================
   Panel root (ticket 77; rebuilt 2026-09-27 to the harness shape owner
   approved: sidebar + main conversation + collapsed right rail, not a bare
   icon rail).
   ----------------------------------------------------------------------------

   State policy: NONE (PLAN.md:1044). The only data this component renders is
   the snapshot the Go side hands it, and the only way to change what is on
   screen is a new snapshot - including which of the views is showing
   (Q2: `view` arrives through currentView(), no component holds it). In
   production the pump is ticket 35's bridge push; with no host attached the
   component tree still renders an empty frame rather than reaching for a
   cached copy, because a WebView restart must land in the same state the Go
   side reports, not in whatever the page remembered.

   Two things deliberately do NOT move with the view: the pending L2 cards and
   the composer row. Hiding a strong-confirmation card because the user clicked
   to another screen would be a safety change smuggled into a navigation change;
   the sidebar's 审批 badge is what carries the queue depth instead.

   The visual acceptance surface is the library-style showcase (main.tsx renders
   it under ?harness=1); this tree is the production panel - the WebView2 window
   IS the panel, there is no fake desktop or floating window to draw.

   Shape (owner 2026-09-27): the old bare icon rail is gone. The panel is now
   the harness's own skeleton: a left sidebar (workspace chip, search, 新建
   对话, conversation history, sections for 审批/任务/球状态, bottom settings),
   a main conversation area (new-task entry or session view), and a right rail
   that stays collapsed until asked. The 9 views become sidebar sections plus
   a full-page settings; the rail icons are retired as navigation.
   ============================================================================ */

import { useState } from "react";
import { Moon, PanelRight, PanelRightClose, Sun } from "lucide-react";
import { ConfigScreen } from "@/components/config-screen";
import { L2ApprovalCard } from "@/components/l2-approval-card";
import { PanelSidebar, type PanelSidebarSection } from "@/components/panel-sidebar";
import { RightRail } from "@/components/harness/right-rail";
import { SettingsPage } from "@/components/harness/settings-page";
import { Composer } from "@/components/composer";
import { currentView, type PanelViewId } from "@/lib/panel-views";
import type { ComposerState, PanelSnapshot } from "@/lib/panel";

const EMPTY_COMPOSER: ComposerState = {
  mode: { current: "unknown", names: [], l2ConfirmNames: [] },
  workspace: {
    set: false,
    spelling: "",
    canonical: "",
    reparse: false,
    rewritten: false,
    reason: "尚未收到原生侧的状态快照",
  },
  attachments: [],
  acceptedAttachmentMimes: [],
  maxAttachmentBytes: 0,
  attachmentError: "",
};

const EMPTY: PanelSnapshot = {
  pending: [],
  results: [],
  composer: EMPTY_COMPOSER,
  generatedAt: "",
};

/** A screen Go has no data for yet says so, and says which ticket would fill
    it. Inventing rows here is what owner's P9 red line forbids ("不得造假数据
    当真实字段"), and an empty <div> would read as a bug rather than as a gap. */
function UnfedScreen({ id }: { id: PanelViewId }) {
  return (
    <p className="text-[12px] text-muted-foreground">
      {id} 屏还没有接到数据。面板不许自己造内容，所以这一屏只说明它缺
      C17 的读口与票 35 的推送。
    </p>
  );
}

export default function App({
  snapshot = EMPTY,
}: {
  snapshot?: PanelSnapshot;
}) {
  // A snapshot from a host that has not filled the composer section yet falls
  // back to the empty view, which says so in words; inventing a mode here would
  // be the panel guessing about a permission it was never told.
  const composer = snapshot.composer ?? EMPTY_COMPOSER;

  // Which screen is showing. A snapshot that names one outranks anything local,
  // because the host is the truth source (Q2). The pick below only applies while
  // no snapshot names a view - today that is every snapshot, since ticket 35's
  // field does not exist yet - and it dies with the document, which is what
  // PLAN.md:1043-1044 asks for: nothing may be REMEMBERED across a WebView
  // restart, but "which of my own screens am I on" is not a memory.
  // railOpen / settingsOpen / settingsSection / dark live inside the same
  // state object so the gate's "exactly one view pick" rule stays true while
  // the rail toggle and the settings page still work.
  const [local, setLocal] = useState<{
    view: PanelViewId | null;
    railOpen: boolean;
    settingsOpen: boolean;
    settingsSection: string;
    dark: boolean;
  }>({ view: null, railOpen: false, settingsOpen: false, settingsSection: "通用", dark: false });
  const hostView = (snapshot as { view?: string }).view;
  const view = hostView !== undefined ? currentView(snapshot) : (local.view ?? currentView(snapshot));

  // Right rail: collapsed until asked.
  const railOpen = local.railOpen;
  function setRailOpen(open: boolean | ((v: boolean) => boolean)) {
    setLocal((current) => ({ ...current, railOpen: typeof open === "function" ? open(current.railOpen) : open }));
  }

  // Settings is a full page, not a modal. It owns the left category sidebar.
  const settingsOpen = local.settingsOpen;
  const settingsSection = local.settingsSection;
  function setSettingsOpen(open: boolean) {
    setLocal((current) => ({ ...current, settingsOpen: open }));
  }
  function setSettingsSection(section: string) {
    setLocal((current) => ({ ...current, settingsSection: section }));
  }

  const dark = local.dark;

  function toggleTheme() {
    const next = !dark;
    setLocal((current) => ({ ...current, dark: next }));
    document.documentElement.dataset.theme = next ? "dark" : "light";
  }

  // The sidebar's sections map to the old views. Clicking one selects it; the
  // main column shows the matching screen. The settings section opens the
  // full settings page instead of a sidebar selection.
  function handleSidebarSelect(id: PanelSidebarSection) {
    if (id === "settings") {
      setSettingsOpen(true);
      return;
    }
    setLocal((current) => ({ ...current, view: id as PanelViewId }));
    setSettingsOpen(false);
  }

  return (
    <div className="flex h-screen w-full flex-col overflow-hidden bg-page text-ink">
      {/* 标题栏：一缕 + 状态。生产里窗口按钮归 Win32，这里永远不放假按钮。 */}
      <header className="flex h-11 shrink-0 select-none items-center justify-between gap-2 border-b border-line bg-surface pl-4 pr-3">
        <div className="flex items-center gap-2.5">
          <div className="size-3.5 rounded-full bg-accent" />
          <span className="text-[13px] font-medium text-ink">一缕</span>
        </div>
        <div className="flex items-center gap-1.5">
          {snapshot.pending.length > 0 && (
            <span className="inline-flex h-5 items-center rounded-full bg-orange-tint px-2 text-[10.5px] font-medium text-orange">
              待审批 {snapshot.pending.length}
            </span>
          )}
          <button
            className="flex items-center gap-1.5 rounded-control px-2 py-1 text-[11.5px] text-ink-2 transition-colors duration-150 hover:bg-hover hover:text-ink"
            onClick={() => setRailOpen((v) => !v)}
            title={railOpen ? "收起右侧栏" : "展开右侧栏"}
            type="button"
          >
            {railOpen ? <PanelRightClose aria-hidden="true" size={15} /> : <PanelRight aria-hidden="true" size={15} />}
            <span className="hidden sm:inline">{railOpen ? "收起" : "展开"}</span>
          </button>
          <button
            className="flex size-7 items-center justify-center rounded-control text-ink-2 transition-colors duration-150 hover:bg-hover hover:text-ink"
            onClick={toggleTheme}
            title="切换主题"
            type="button"
          >
            {dark ? <Sun aria-hidden="true" size={15} /> : <Moon aria-hidden="true" size={15} />}
          </button>
        </div>
      </header>

      <div className="flex min-h-0 flex-1 overflow-hidden">
        {/* 左栏：区段式导航。视图 + 会话历史都是 section，设置打开整页。 */}
        <PanelSidebar
          active={view}
          onNewTask={() => setLocal((current) => ({ ...current, view: "chat" }))}
          onOpenSettings={() => setSettingsOpen(true)}
          onSelect={handleSidebarSelect}
          onToggleTheme={toggleTheme}
          pendingCount={snapshot.pending.length}
          theme={dark ? "dark" : "light"}
        />

        {/* 主区：设置整页 或 对话区 */}
        <div className="flex min-w-0 flex-1 flex-col">
          {settingsOpen ? (
            <SettingsPage
              activeSection={settingsSection}
              onSelectSection={setSettingsSection}
              onClose={() => setSettingsOpen(false)}
            />
          ) : (
            <>
              {/* 顶栏：当前视图指示 + 右栏开关 */}
              <div className="flex h-10 shrink-0 items-center justify-between border-b border-line bg-surface px-4">
                <span className="text-[12px] font-medium text-ink-2">
                  {view === "chat" && "对话"}
                  {view === "approval" && "审批中心"}
                  {view === "palette" && "命令"}
                  {view === "tasks" && "任务"}
                  {view === "ball" && "球状态"}
                  {view === "config" && "设置"}
                  {view === "security" && "安全"}
                  {view === "privacy" && "隐私"}
                  {view === "cost" && "成本"}
                </span>
                <button
                  className="flex items-center gap-1 rounded-control px-2 py-1 text-[11.5px] text-ink-2 transition-colors duration-150 hover:bg-hover hover:text-ink"
                  onClick={() => setRailOpen((v) => !v)}
                  type="button"
                >
                  {railOpen ? <PanelRightClose aria-hidden="true" size={13} /> : <PanelRight aria-hidden="true" size={13} />}
                  <span>{railOpen ? "收起右栏" : "展开右栏"}</span>
                </button>
              </div>

              {/* 内容区 */}
              <div className="flex min-h-0 flex-1 flex-col overflow-y-auto">
                {snapshot.pending.map((v) => (
                  <L2ApprovalCard key={v.correlationId} view={v} />
                ))}

                {view === "chat" ? (
                  <div className="flex min-h-0 flex-1 flex-col">
                    <div className="flex min-h-0 flex-1 flex-col overflow-y-auto px-6 py-5">
                      {snapshot.results.length === 0 ? (
                        <div className="flex flex-1 flex-col items-center justify-center gap-3">
                          <p className="text-[13px] text-ink-3">
                            还没有会话。面板不许自己造内容，Go 侧推送的流式结果会出现在这里。
                          </p>
                          <p className="text-[11.5px] text-ink-3">
                            点击左侧「新建对话」开始，或等待原生侧推送。
                          </p>
                        </div>
                      ) : (
                        <div className="mx-auto flex w-full max-w-[720px] flex-col gap-5">
                          {snapshot.results.map((chunk) => (
                            <div key={chunk.correlationId} className="text-[13px] text-ink">
                              {chunk.text}
                            </div>
                          ))}
                        </div>
                      )}
                    </div>
                    {/* composer 钉底 */}
                    <div className="shrink-0 border-t border-line bg-surface px-6 py-3">
                      <Composer state={composer} />
                    </div>
                  </div>
                ) : view === "approval" ? (
                  <div className="flex min-h-0 flex-1 flex-col gap-3 px-6 py-5">
                    <h1 className="text-[15px] font-semibold text-ink">审批中心</h1>
                    <p className="text-[12.5px] text-ink-2">
                      不可逆与高危操作在原生侧批准，面板只递请求、不做决定。
                    </p>
                    {snapshot.pending.length === 0 ? (
                      <p className="text-[12px] leading-[1.6] text-ink-2">
                        队列当前为空。待确认请求由原生侧 Go 推送：出现需要人工确认的调用时，会在这里排队并在面板顶部出卡。
                      </p>
                    ) : (
                      <p className="text-[12px] text-ink-3">
                        有 {snapshot.pending.length} 条待审批，卡片在上方常驻显示。
                      </p>
                    )}
                  </div>
                ) : view === "tasks" ? (
                  <div className="flex min-h-0 flex-1 flex-col gap-3 px-6 py-5">
                    <h1 className="text-[15px] font-semibold text-ink">任务</h1>
                    <p className="text-[12.5px] text-ink-2">
                      还没有接到数据。面板不许自己造内容，这里只说明它缺 C17 的读口与票 35 的推送。
                    </p>
                  </div>
                ) : view === "ball" ? (
                  <div className="flex min-h-0 flex-1 flex-col gap-3 px-6 py-5">
                    <h1 className="text-[15px] font-semibold text-ink">球状态</h1>
                    <p className="text-[12.5px] text-ink-2">
                      还没有接到数据。面板不许自己造内容，这里只说明它缺 C17 的读口与票 35 的推送。
                    </p>
                  </div>
                ) : view === "config" ? (
                  /* config 在旧骨架里是「设置」视图；新骨架里设置是一整页，从
                     侧栏底部进。快照若仍指到 config，渲染外观一节（selfFed 的
                     ConfigScreen），而不是一句「还没接到数据」——那样对设置
                     才是说谎。 */
                  <div className="px-6 py-5">
                    <ConfigScreen />
                  </div>
                ) : (
                  <UnfedScreen id={view} />
                )}
              </div>
            </>
          )}
        </div>

        {/* 右栏：默认收起，打开后显示 launcher（卡片 + 快捷键 chip），不是
            直接进 tabs。tabs 按需打开或 agent 请求。 */}
        {railOpen && (
          <RightRail
            reviewFiles={[]}
            fileTree={[]}
            terminalLines={[]}
            browserUrl=""
          />
        )}
      </div>
    </div>
  );
}
