/* ============================================================================
   Panel root (ticket 77; rebuilt 2026-09-28).
   ----------------------------------------------------------------------------
   State policy: NONE (PLAN.md:1044). The only data this component renders is
   the snapshot the Go side hands it, and the only way to change what is on
   screen is a new snapshot. In production the pump is ticket 35's bridge
   push; with no host attached the tree still renders an empty frame rather
   than reaching for a cached copy, because a WebView restart must land in the
   state the Go side reports, not in whatever the page remembered.

   What changed in this cut, and why
   --------------------------------
   The 2026-09-27 cut rebuilt the shell and got the bones right, but the first
   pass at the main column left three things that made it read as unfinished
   rather than as a panel waiting for data:

     - a second, redundant right-rail toggle sitting in a content-area header
       whose only other content was the current screen's name. One control per
       action is the rule; there was one control, twice.
     - an empty main column that was two sentences adrift in a void. Owner's
       words on a panel with no data: 「哪怕数据没有，但是你该有的样子得合理吧」
       - the SHAPE has to be right even when the content is not. The empty
       column now states what will appear there, using the real section names
       the panel is built from (流式回复 / 工具调用 / 思考过程), which is a
       statement about this panel, not invented content.
     - no context/usage readout above the composer. ContextMeter fills that,
       and says it has not been told when it has not.

   The 5 rows that are not navigable views (命令 / 设置 / 安全 / 隐私 / 成本)
   do NOT appear here: 命令 is the Ctrl+K palette and the other four are
   settings categories, both reachable from their own places. The sidebar is a
   shape with a job, not a transcription of the retired icon rail.

   Colours: token utilities only, zero literals; no emoji.
   ============================================================================ */

import { useState } from "react";
import { Command, Moon, PanelRight, PanelRightClose, Sun } from "lucide-react";
import { L2ApprovalCard } from "@/components/l2-approval-card";
import { Composer } from "@/components/composer";
import { ContextMeter } from "@/components/context-meter";
import { PanelSidebar, type PanelSectionId } from "@/components/panel-sidebar";
import { RightRail } from "@/components/harness/right-rail";
import { SettingsPage } from "@/components/harness/settings-page";
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

/** The main column's four destinations, and the honest sentence for the three
    that Go has not fed. The sentence names the missing feed on purpose: an
    empty <div> reads as a bug, and a row of dashes reads as data. */
const SECTION_COPY: Record<
  Exclude<PanelSectionId, "chat">,
  { title: string; lede: string; body: string }
> = {
  approval: {
    title: "审批",
    lede: "不可逆与高危操作在原生侧批准，面板只递请求、不做决定。",
    body: "队列为空。出现需要你点头的调用时，它会在这里排队，同时在面板顶部出卡。",
  },
  tasks: {
    title: "任务",
    lede: "正在跑的活、排队的活、和被路径锁挡住的活。",
    body: "还没有收到任务数据。面板不自己造这一屏，所以它只说明缺 C17 的读口与票 35 的推送。",
  },
  ball: {
    title: "球状态",
    lede: "悬浮球现在在哪个态，以及它是被什么带过去的。",
    body: "还没有收到球状态数据。面板读不到原生那颗球（D29 是 Direct2D 画的），这里等原生侧推送。",
  },
};

export default function App({ snapshot = EMPTY }: { snapshot?: PanelSnapshot }) {
  // A host that has not filled the composer section falls back to the empty
  // view, which says so in words; inventing a mode would be the panel guessing
  // at a permission it was never told.
  const composer = snapshot.composer ?? EMPTY_COMPOSER;

  // Which screen is showing. A snapshot that names one outranks anything local
  // (Q2): the host is the truth source. The local pick applies only while no
  // snapshot names a view - today that is every snapshot, since ticket 35's
  // field does not exist yet - and it dies with the document, which is what
  // PLAN.md:1043-1044 asks for: nothing may be REMEMBERED across a WebView
  // restart, but "which of my own screens am I on" is not a memory.
  //
  // ONE state object, not four. The gate pins the count so this cannot quietly
  // grow into a store; rail/settings/theme are view preferences of the same
  // class as the alpha knob in the settings screen, and they die with it too.
  const [local, setLocal] = useState<{
    view: PanelViewId | null;
    railOpen: boolean;
    settingsOpen: boolean;
    settingsSection: string;
    dark: boolean;
  }>({
    view: null,
    railOpen: false,
    settingsOpen: false,
    settingsSection: "通用",
    dark: false,
  });

  const hostView = (snapshot as { view?: string }).view;
  const view = hostView !== undefined ? currentView(snapshot) : (local.view ?? currentView(snapshot));

  /** The sidebar's four sections are the four the main column can show. A
      snapshot naming one of the five retired rows lands on 对话 rather than
      guessing - see currentView for the same rule one layer down. */
  const section: PanelSectionId =
    view === "approval" || view === "tasks" || view === "ball" ? view : "chat";

  function selectSection(next: PanelSectionId) {
    setLocal((current) => ({ ...current, view: next, settingsOpen: false }));
  }

  function toggleTheme() {
    const next = !local.dark;
    setLocal((current) => ({ ...current, dark: next }));
    document.documentElement.dataset.theme = next ? "dark" : "light";
  }

  return (
    <div className="flex h-screen w-full flex-col overflow-hidden bg-page text-ink">
      <header className="flex h-11 shrink-0 select-none items-center justify-between gap-2 border-b border-line bg-surface pl-4 pr-2.5">
        <div className="flex min-w-0 items-center gap-2.5">
          <span aria-hidden="true" className="size-2.5 shrink-0 rounded-full bg-accent" />
          <span className="shrink-0 text-[13px] font-medium text-ink">一缕</span>
          {snapshot.pending.length > 0 ? (
            <span className="shrink-0 rounded-full bg-orange-tint px-2 py-0.5 text-[10.5px] font-medium tabular-nums text-orange">
              待审批 {snapshot.pending.length}
            </span>
          ) : null}
        </div>
        <div className="flex shrink-0 items-center gap-0.5">
          <kbd className="kbd mr-1 hidden sm:flex">Ctrl K</kbd>
          <button
            aria-label="命令面板"
            className="flex size-7 items-center justify-center rounded-control text-ink-2 transition-colors duration-150 hover:bg-hover hover:text-ink"
            title="命令面板"
            type="button"
          >
            <Command aria-hidden="true" size={14} strokeWidth={1.9} />
          </button>
          <button
            aria-label={local.railOpen ? "收起右侧栏" : "展开右侧栏"}
            className="flex size-7 items-center justify-center rounded-control text-ink-2 transition-colors duration-150 hover:bg-hover hover:text-ink"
            onClick={() => setLocal((current) => ({ ...current, railOpen: !current.railOpen }))}
            title={local.railOpen ? "收起右侧栏" : "展开右侧栏"}
            type="button"
          >
            {local.railOpen ? (
              <PanelRightClose aria-hidden="true" size={15} />
            ) : (
              <PanelRight aria-hidden="true" size={15} />
            )}
          </button>
          <button
            aria-label="切换主题"
            className="flex size-7 items-center justify-center rounded-control text-ink-2 transition-colors duration-150 hover:bg-hover hover:text-ink"
            onClick={toggleTheme}
            title="切换主题"
            type="button"
          >
            {local.dark ? <Sun aria-hidden="true" size={15} /> : <Moon aria-hidden="true" size={15} />}
          </button>
        </div>
      </header>

      <div className="flex min-h-0 flex-1 overflow-hidden">
        <PanelSidebar
          active={section}
          activeSessionId={null}
          onOpenSettings={() => setLocal((current) => ({ ...current, settingsOpen: true }))}
          onSelect={selectSection}
          onSelectSession={() => undefined}
          onToggleTheme={toggleTheme}
          pendingCount={snapshot.pending.length}
          sessionGroups={[]}
          settingsOpen={local.settingsOpen}
          theme={local.dark ? "dark" : "light"}
          workspaceName="一缕"
          workspacePath={
            composer.workspace.set ? composer.workspace.canonical : "尚未收到工作区"
          }
        />

        <div className="flex min-w-0 flex-1 flex-col">
          {local.settingsOpen ? (
            <SettingsPage
              activeSection={local.settingsSection}
              onClose={() => setLocal((current) => ({ ...current, settingsOpen: false }))}
              onSelectSection={(next) =>
                setLocal((current) => ({ ...current, settingsSection: next }))
              }
            />
          ) : (
            <>
              {/* Pending cards sit ABOVE the section switch and on every
                  section: hiding a strong-confirmation card because the user
                  navigated would be a safety change smuggled into a
                  navigation change. The sidebar's 审批 badge carries depth. */}
              {snapshot.pending.map((v) => (
                <L2ApprovalCard key={v.correlationId} view={v} />
              ))}

              <div className="flex min-h-0 flex-1 flex-col">
                {section === "chat" ? (
                  <>
                    <div className="flex min-h-0 flex-1 flex-col overflow-y-auto">
                      {/* ONE measuring column for the whole main area. The
                          results, the placeholder and the composer all sit on
                          the same axis, so the text above the input lines up
                          with the text in it. Two different widths read as two
                          different panels. */}
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
                    <div className="shrink-0 border-t border-line bg-surface px-6 pb-4 pt-2.5">
                      <div className="mx-auto w-full max-w-[720px]">
                        <ContextMeter usage={{}} className="mb-1.5" />
                        <Composer state={composer} />
                      </div>
                    </div>
                  </>
                ) : (
                  <SectionPlaceholder copy={SECTION_COPY[section]} />
                )}
              </div>
            </>
          )}
        </div>

        {/* Right rail: collapsed until asked, and then it lands on a LAUNCHER
            of things you can open rather than on tabs that are already up. */}
        {local.railOpen ? (
          <RightRail reviewFiles={[]} fileTree={[]} terminalLines={[]} browserUrl="" />
        ) : null}
      </div>
    </div>
  );
}

/* ---------------------------------------------------------------------------
   The conversation placeholder. This is the piece the earlier cut got wrong by
   omission: it drew nothing, so a panel with no data looked like a panel that
   had failed to load. It states the shape of what arrives instead - the three
   blocks the result stream is actually built from - and says plainly that
   nothing has arrived yet. Every word here is about THIS panel's own
   construction, so it cannot be mistaken for a fabricated result.
   --------------------------------------------------------------------------- */

const ARRIVES = ["流式回复", "工具调用", "思考过程"] as const;

function ConversationPlaceholder() {
  return (
    /* Not dead centre. A panel window is short, and dead-centring in it drops
       the sentence into the middle of a void with nothing above it - which is
       how the 2026-09-27 build looked like it had failed to load. Sitting at
       roughly a third puts the eye at the top third of the column, which is
       where an opening line belongs. */
    <div className="flex flex-1 items-start justify-center pt-[18vh]">
      <div className="w-full max-w-[440px] text-center" style={{ animation: "fade-up 320ms var(--ease-out-strong) both" }}>
        <p className="text-[15px] font-medium text-ink">一缕在这儿</p>
        <p className="mx-auto mt-1.5 max-w-[360px] text-[12.5px] leading-[1.7] text-ink-3">
          面板不自己造内容。在下面说一句话，或者从悬浮球开口，
          原生侧会把结果一条条推上来。
        </p>
        <div className="mt-4 flex items-center justify-center gap-1.5">
          {ARRIVES.map((label) => (
            <span
              key={label}
              className="rounded-full border border-line bg-surface px-2.5 py-1 text-[11px] text-ink-3"
            >
              {label}
            </span>
          ))}
        </div>
        <p className="mt-3.5 text-[10.5px] text-ink-3">还没收到任何一条</p>
      </div>
    </div>
  );
}

/* ---------------------------------------------------------------------------
   A section with no feed yet. The lede says what the section is FOR (so the
   shape is legible), the body says the feed is missing and names it. Two
   sentences doing two different jobs, which is the difference between an
   empty state and a shrug.
   --------------------------------------------------------------------------- */

function SectionPlaceholder({
  copy,
}: {
  copy: { title: string; lede: string; body: string };
}) {
  return (
    <div className="flex min-h-0 flex-1 flex-col overflow-y-auto">
      <div className="mx-auto flex w-full max-w-[720px] flex-1 flex-col justify-center px-6">
        <div className="w-full max-w-[520px]" style={{ animation: "fade-up 320ms var(--ease-out-strong) both" }}>
          <h1 className="text-[15px] font-semibold text-ink">{copy.title}</h1>
          <p className="mt-1 text-[12.5px] leading-[1.7] text-ink-2">{copy.lede}</p>
          <div className="mt-4 rounded-card border border-line bg-surface px-3.5 py-3">
            <p className="text-[12px] leading-[1.7] text-ink-3">{copy.body}</p>
          </div>
        </div>
      </div>
    </div>
  );
}
