/* ============================================================================
   Panel root (ticket 77; third cut, 2026-09-28).
   ----------------------------------------------------------------------------
   State policy: NONE (PLAN.md:1044). Everything rendered here comes from the
   snapshot the Go side pushes; the only local state is view preference (which
   rail is open, which theme, which panel is up), and it dies with the
   document. A WebView restart must land in the state Go reports, never in
   whatever the page remembered.

   What this cut fixes, from owner's 2026-09-28 review
   --------------------------------------------------
   The two prior cuts both put a MENU in the left sidebar. It is not a menu.
   It is a tree of the user's own work - the workspace and the sessions under
   it - and three things that were sitting in it do not belong there at all:

     审批   → a card that appears when a gated call needs you, plus a count on
              the header. A permanent row for something empty 99% of the time
              trains you to ignore it.
     任务   → the 任务监控 popover in the top bar, which is where a live view
              of running work belongs.
     球状态 → gone. The ball is a native Direct2D surface (D29); its settings
              are on the settings page.

   Also fixed here, all of them found by actually clicking rather than by
   reading the code:
     - the command palette button was drawn with no handler. A control that
       does nothing is worse than a missing one: it teaches people to stop
       believing the panel.
     - popovers closed only on their own trigger. They now dismiss on an
       outside pointerdown and on Escape (lib/use-dismiss.ts).
     - the panel did not survive a 640px window, which is the DEFAULT
       `[panel] width` in internal/config/schema.go. The sidebar collapses to
       icons and the rail overlays instead of stealing a third of the column.

   Colours: token utilities only, zero literals; no emoji.
   ============================================================================ */

import { useCallback, useEffect, useRef, useState } from "react";
import { Activity, Command, Moon, PanelLeft, PanelRight, PanelRightClose, Sun } from "lucide-react";
import { L2ApprovalCard } from "@/components/l2-approval-card";
import { ContextMeter } from "@/components/context-meter";
import { GitBranchChip } from "@/components/git-branch";
import { MonitorPopover } from "@/components/harness/monitor-popover";
import { PaletteScreen } from "@/components/palette-screen";
import { PromptBar } from "@/components/ai-native/prompt-bar";
import { LineSidebar } from "@/components/react-bits/line-sidebar";
import { RightRail } from "@/components/harness/right-rail";
import { SettingsPage } from "@/components/harness/settings-page";
import { PanelSidebar } from "@/components/panel-sidebar";
import { useDismiss } from "@/lib/use-dismiss";
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

/** Native mode names -> the Chinese the tier row shows. Spelled once, in the
    component that draws them, so a fourth tier cannot be invented here. */
const TIER_NAMES: Record<string, string> = {
  ask_every_step: "每步都问",
  ask_high_risk: "高风险才问",
  auto_approve: "全自动",
  unknown: "档位未知",
};

/** Narrower than this and the sidebar stops fitting a session title, so it
    collapses to its icon rail. 640px is the DEFAULT panel width, not an edge
    case - see the header. */
const SIDEBAR_COLLAPSE_BELOW = 860;

export default function App({ snapshot = EMPTY }: { snapshot?: PanelSnapshot }) {
  const composer = snapshot.composer ?? EMPTY_COMPOSER;

  // ONE state object. The gate pins the count, and the reason is in the
  // render-nav header: a gate derived from what it checks cannot catch that
  // thing being wrong.
  const [local, setLocal] = useState<{
    view: PanelViewId | null;
    railOpen: boolean;
    settingsOpen: boolean;
    settingsSection: string;
    dark: boolean;
    monitorOpen: boolean;
    paletteOpen: boolean;
    sessionId: string | null;
    model: string;
    effort: string;
    activeTurn: string | null;
  }>({
    view: null,
    railOpen: false,
    settingsOpen: false,
    settingsSection: "通用",
    dark: false,
    monitorOpen: false,
    paletteOpen: false,
    sessionId: null,
    model: "",
    effort: "",
    activeTurn: null,
  });

  const hostView = (snapshot as { view?: string }).view;
  const view = hostView !== undefined ? currentView(snapshot) : (local.view ?? currentView(snapshot));
  // Four of the nine legacy screens were never navigation - they are settings
  // categories. A snapshot that names one lands on the settings page AT that
  // category, which is where it belongs. Sending it to the conversation column
  // would put a setting request in a chat window.
  const CATEGORY_VIEW: Partial<Record<PanelViewId, string>> = {
    config: "外观",
    security: "安全",
    privacy: "隐私",
    cost: "成本",
  };
  const categoryOf = CATEGORY_VIEW[view];
  const narrow = useNarrow();
  // At the DEFAULT 640px panel width a 240px sidebar leaves 400px for the
  // conversation, which is not a conversation. So below the threshold the
  // sidebar becomes an overlay DRAWER, closed by default: a panel that opens
  // with a tree laid over it is a panel you cannot read. The handle in the
  // header opens it; the conversation is what you came for.
  const [sidebarOpen, setSidebarOpen] = useState(false);
  useEffect(() => {
    // Crossing the threshold resets the drawer. Otherwise shrinking the window
    // leaves a tree covering the conversation, and growing it back leaves a
    // tree nobody asked for.
    setSidebarOpen(false);
  }, [narrow]);
  const sidebarVisible = !narrow || sidebarOpen;
  const monitorRef = useRef<HTMLDivElement>(null);
  useDismiss(monitorRef, local.monitorOpen, useCallback(() => setLocal((c) => ({ ...c, monitorOpen: false })), []));

  // Ctrl+K is the one keyboard verb every harness shares, and the chip in the
  // header advertises it. A chip promising a key with no listener is worse
  // than no chip. Escape closes the palette from the same handler, because a
  // modal you can only leave with the mouse is not a modal you can use.
  useEffect(() => {
    function onKey(event: KeyboardEvent) {
      if ((event.ctrlKey || event.metaKey) && event.key.toLowerCase() === "k") {
        event.preventDefault();
        setLocal((c) => ({ ...c, paletteOpen: !c.paletteOpen }));
        return;
      }
      if (event.key === "Escape") {
        setLocal((c) => ({ ...c, paletteOpen: false, monitorOpen: false }));
        if (narrow) setSidebarOpen(false);
      }
    }
    document.addEventListener("keydown", onKey);
    return () => document.removeEventListener("keydown", onKey);
  }, [narrow]);

  function toggleTheme() {
    const next = !local.dark;
    setLocal((c) => ({ ...c, dark: next }));
    document.documentElement.dataset.theme = next ? "dark" : "light";
  }

  // The permission tiers come from the snapshot's own vocabulary. The panel
  // asks; it never records the answer, so nothing here is optimistic.
  const tiers = composer.mode.names.length > 0 ? composer.mode.names : ["unknown"];

  return (
    <div className="flex h-screen w-full flex-col overflow-hidden bg-page text-ink">
      <header className="flex h-11 shrink-0 select-none items-center gap-2 border-b border-line bg-surface pl-4 pr-2.5">
        <div className="flex min-w-0 shrink-0 items-center gap-2.5">
          {/* At the default panel width the tree is a drawer, so it needs a
              handle. Above the threshold this button is not rendered at all -
              a control for a mode you are not in is clutter. */}
          {narrow ? (
            <button
              type="button"
              aria-label={sidebarOpen ? "收起会话树" : "展开会话树"}
              aria-expanded={sidebarOpen}
              onClick={() => setSidebarOpen((v) => !v)}
              className="flex size-7 shrink-0 items-center justify-center rounded-control text-ink-2 transition-colors duration-150 hover:bg-hover hover:text-ink"
            >
              <PanelLeft aria-hidden="true" size={15} />
            </button>
          ) : null}
          <span aria-hidden="true" className="size-2.5 rounded-full bg-accent" />
          <span className="shrink-0 text-[13px] font-medium text-ink">一缕</span>
          {/* 审批 is a COUNT here, not a place. It lights only when something
              is actually waiting, which is the whole point of moving it out
              of the sidebar. */}
          {snapshot.pending.length > 0 ? (
            <span className="shrink-0 rounded-full bg-warn-soft px-2 py-0.5 text-[10.5px] font-medium tabular-nums text-warn">
              待审批 {snapshot.pending.length}
            </span>
          ) : null}
        </div>

        <div className="flex min-w-0 flex-1 items-center gap-1.5">
          {/* git branch. Reads its strings from the snapshot; a workspace that
              is not a repo renders nothing at all rather than a chip. */}
          <GitBranchChip
            view={{ current: "", isRepo: false, local: [], worktrees: [] }}
            onCheckout={() => undefined}
            onNewWorktree={() => undefined}
          />
        </div>

        <div className="flex shrink-0 items-center gap-0.5">
          {/* 任务监控 - a live view of running work, in a popover that gets out
              of the way. lib/use-dismiss closes it on an outside click. */}
          <div className="relative" ref={monitorRef}>
            <button
              type="button"
              aria-expanded={local.monitorOpen}
              aria-label="任务监控"
              onClick={() => setLocal((c) => ({ ...c, monitorOpen: !c.monitorOpen }))}
              className={`flex h-7 size-7 shrink-0 items-center justify-center rounded-control transition-colors duration-150 ${
                local.monitorOpen ? "bg-hover text-ink" : "text-ink-2 hover:bg-hover hover:text-ink"
              }`}
            >
              <Activity aria-hidden="true" size={14} />
            </button>
            {local.monitorOpen ? (
              <div className="absolute right-0 top-full z-40 mt-1.5">
                <MonitorPopover
                  open
                  env={[
                    { label: "工作区", value: composer.workspace.set ? composer.workspace.canonical : "尚未收到" },
                    { label: "分支", value: "尚未读到" },
                  ]}
                  agents={[]}
                  onPin={() => undefined}
                />
              </div>
            ) : null}
          </div>

          <kbd className="kbd hidden sm:flex">Ctrl K</kbd>
          <button
            type="button"
            aria-label="命令面板"
            className="flex size-7 items-center justify-center rounded-control text-ink-2 transition-colors duration-150 hover:bg-hover hover:text-ink"
            onClick={() => setLocal((c) => ({ ...c, paletteOpen: !c.paletteOpen }))}
            title="命令面板"
          >
            <Command aria-hidden="true" size={14} strokeWidth={1.9} />
          </button>
          <button
            type="button"
            aria-label={local.railOpen ? "收起右侧栏" : "展开右侧栏"}
            className="flex size-7 items-center justify-center rounded-control text-ink-2 transition-colors duration-150 hover:bg-hover hover:text-ink"
            onClick={() => setLocal((c) => ({ ...c, railOpen: !c.railOpen }))}
            title={local.railOpen ? "收起右侧栏" : "展开右侧栏"}
          >
            {local.railOpen ? <PanelRightClose aria-hidden="true" size={15} /> : <PanelRight aria-hidden="true" size={15} />}
          </button>
          <button
            type="button"
            aria-label="切换主题"
            className="flex size-7 items-center justify-center rounded-control text-ink-2 transition-colors duration-150 hover:bg-hover hover:text-ink"
            onClick={toggleTheme}
            title="切换主题"
          >
            {local.dark ? <Sun aria-hidden="true" size={15} /> : <Moon aria-hidden="true" size={15} />}
          </button>
        </div>
      </header>

      <div className="relative flex min-h-0 flex-1 overflow-hidden">
        {sidebarVisible ? (
          <div
            className={
              narrow
                ? "absolute inset-y-0 left-0 z-30 shadow-overlay"
                : "contents"
            }
          >
            <PanelSidebar
              activeSessionId={local.sessionId}
              historyGroups={[]}
              onNewTask={() => setLocal((c) => ({ ...c, sessionId: null, view: "chat" }))}
              onOpenSettings={() => setLocal((c) => ({ ...c, settingsOpen: true }))}
              onSelectSession={(id) => setLocal((c) => ({ ...c, sessionId: id, view: "chat" }))}
              onToggleTheme={toggleTheme}
              settingsOpen={local.settingsOpen}
              theme={local.dark ? "dark" : "light"}
              workspaceName="一缕"
              workspacePath={composer.workspace.set ? composer.workspace.canonical : "尚未收到工作区"}
              workspaceSessions={[]}
            />
          </div>
        ) : null}

        <div className="flex min-w-0 flex-1 flex-col">
          {local.settingsOpen || categoryOf ? (
            <SettingsPage
              activeSection={categoryOf ?? local.settingsSection}
              onClose={() => setLocal((c) => ({ ...c, settingsOpen: false }))}
              onSelectSection={(next) => setLocal((c) => ({ ...c, settingsSection: next }))}
            />
          ) : (
            <>
              {/* The L2 cards live here, ABOVE everything and on every view. A
                  strong-confirmation card that vanished because you clicked
                  elsewhere would be a safety change smuggled into navigation. */}
              {snapshot.pending.map((v) => (
                <L2ApprovalCard key={v.correlationId} view={v} />
              ))}

              <div className="flex min-h-0 flex-1 flex-col">
                {view === "approval" || view === "tasks" || view === "ball" ? (
                  <RetiredView id={view} />
                ) : (
                  <Conversation
                    effort={local.effort}
                    onEffortChange={(effort) => setLocal((c) => ({ ...c, effort }))}
                    onPermissionChange={() => undefined}
                    onSend={() => undefined}
                    onJump={(id) => setLocal((c) => ({ ...c, activeTurn: id }))}
                    permission={composer.mode.current}
                    snapshot={snapshot}
                    tiers={tiers}
                    model={local.model}
                    onModelChange={(model) => setLocal((c) => ({ ...c, model }))}
                  />
                )}
              </div>
            </>
          )}
        </div>

        {local.railOpen ? (
          <div className={narrow ? "absolute inset-y-0 right-0 z-30 shadow-overlay" : ""}>
            <RightRail reviewFiles={[]} fileTree={[]} terminalLines={[]} browserUrl="" />
          </div>
        ) : null}
      </div>

      {local.paletteOpen ? (
        <div
          aria-label="命令面板"
          aria-modal="true"
          className="fixed inset-0 z-40 flex items-start justify-center bg-black/35 pt-[12vh]"
          onClick={(e) => {
            if (e.target === e.currentTarget) setLocal((c) => ({ ...c, paletteOpen: false }));
          }}
          role="dialog"
        >
          <div
            className="w-[560px] overflow-hidden rounded-card border border-line bg-surface shadow-overlay"
            style={{ animation: "fade-up 200ms var(--ease-out-strong) both" }}
          >
            <PaletteScreen
              groups={[]}
              onQueryChange={() => undefined}
              onSelect={() => setLocal((c) => ({ ...c, paletteOpen: false }))}
              query=""
            />
          </div>
        </div>
      ) : null}
    </div>
  );
}

/** Track the window width. A resize listener, not a media query: the threshold
    is about whether a session title fits beside the icons, not about a device
    class. 640px is the DEFAULT panel width (internal/config/schema.go), so
    this path is the common one, not the edge case. */
function useNarrow(): boolean {
  const read = useCallback(
    () => typeof window !== "undefined" && window.innerWidth < SIDEBAR_COLLAPSE_BELOW,
    [],
  );
  const [narrow, setNarrow] = useState(read);
  useEffect(() => {
    const onResize = () => setNarrow(read());
    window.addEventListener("resize", onResize);
    // Re-read once on mount: the window can have been resized between the
    // useState initialiser and this effect running.
    onResize();
    return () => window.removeEventListener("resize", onResize);
  }, [read]);
  return narrow;
}

/* ---------------------------------------------------------------------------
   The conversation column: the turn rail, the stream, and the composer.
   --------------------------------------------------------------------------- */

function Conversation({
  snapshot,
  tiers,
  permission,
  onPermissionChange,
  model,
  onModelChange,
  effort,
  onEffortChange,
  onSend,
  onJump,
}: {
  snapshot: PanelSnapshot;
  tiers: string[];
  permission: string;
  onPermissionChange: (key: string) => void;
  model: string;
  onModelChange: (key: string) => void;
  effort: string;
  onEffortChange: (key: string) => void;
  onSend: (text: string) => void;
  onJump: (id: string) => void;
}) {
  // One turn per result chunk, keyed by correlationId. The rail needs at least
  // two marks to be worth drawing, and TurnRail hides itself below that.
  const turns = snapshot.results.map((chunk) => ({
    id: chunk.correlationId,
    label: chunk.text.slice(0, 18) || "（空）",
  }));

  return (
    <>
      <div className="flex min-h-0 flex-1 overflow-hidden">
        {/* react-bits LineSidebar: the cursor-proximity rail for long
            conversations. Mounted only past one turn - a rail with a single
            mark is furniture pretending to be a control. */}
        {turns.length > 1 ? (
          <div className="flex w-44 shrink-0 items-center border-r border-line px-2">
            <LineSidebar
              fontSize={0.8}
              itemGap={10}
              items={turns.map((t) => t.label)}
              markerLength={26}
              maxShift={14}
              onItemClick={(index) => onJump(turns[index]!.id)}
              showIndex
            />
          </div>
        ) : null}
        <div className="flex min-h-0 flex-1 flex-col overflow-y-auto">
          <div className="mx-auto flex w-full max-w-[720px] flex-1 flex-col px-6">
            {snapshot.results.length === 0 ? (
              <ConversationPlaceholder />
            ) : (
              <div className="flex flex-col gap-5 py-5">
                {snapshot.results.map((chunk) => (
                  <p className="text-[13px] leading-relaxed text-ink" key={chunk.correlationId}>
                    {chunk.text}
                  </p>
                ))}
              </div>
            )}
          </div>
        </div>
      </div>
      <div className="shrink-0 border-t border-line bg-surface px-6 pb-4 pt-2.5">
        <div className="mx-auto w-full max-w-[720px]">
          <ContextMeter usage={{}} className="mb-1.5" />
          <PromptBar
            effort={effort}
            efforts={[]}
            model={model}
            models={[]}
            onEffortChange={onEffortChange}
            onModelChange={onModelChange}
            onPermissionChange={onPermissionChange}
            onSend={onSend}
            permission={permission}
            permissions={tiers.map((t) => ({ key: t, name: TIER_NAMES[t] ?? t }))}
          />
        </div>
      </div>
    </>
  );
}

/* The three views that were in the sidebar and are not views at all. Kept as
   one honest sentence each so that a snapshot still naming one lands somewhere
   that explains itself, instead of on a blank screen. */
function RetiredView({ id }: { id: PanelViewId }) {
  const copy: Record<string, string> = {
    approval: "审批不是一屏，是一个卡片：真有高危请求时它会自己出现在面板顶部，右上角的计数告诉你有几个在等你。",
    tasks: "任务不是一个地方，是一个浮层：点顶栏的「任务监控」看正在跑的活。",
    ball: "悬浮球是原生画的（D29），不在面板里。球的设置在设置页。",
  };
  return (
    <div className="flex min-h-0 flex-1 items-start justify-center px-6 py-10">
      <div className="w-full max-w-[520px]" style={{ animation: "fade-up 320ms var(--ease-out-strong) both" }}>
        <p className="text-[13px] leading-[1.8] text-ink-2">{copy[id] ?? "这一屏已经不存在了。"}</p>
      </div>
    </div>
  );
}

const ARRIVES = ["流式回复", "工具调用", "思考过程"] as const;

function ConversationPlaceholder() {
  return (
    <div className="flex flex-1 items-start justify-center pt-[13vh]">
      <div className="w-full max-w-[440px] text-center" style={{ animation: "fade-up 320ms var(--ease-out-strong) both" }}>
        <p className="text-[15px] font-medium text-ink">一缕在这儿</p>
        <p className="mx-auto mt-1.5 max-w-[360px] text-[12.5px] leading-[1.7] text-ink-3">
          面板不自己造内容。在下面说一句话，或者从悬浮球开口，原生侧会把结果一条条推上来。
        </p>
        <div className="mt-4 flex items-center justify-center gap-1.5">
          {ARRIVES.map((label) => (
            <span key={label} className="rounded-full border border-line bg-surface px-2.5 py-1 text-[11px] text-ink-3">
              {label}
            </span>
          ))}
        </div>
        <p className="mt-3.5 text-[10.5px] text-ink-3">还没收到任何一条</p>
      </div>
    </div>
  );
}
