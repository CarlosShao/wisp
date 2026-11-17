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
import { Activity, Command, Moon, PanelLeft, PanelRight, PanelRightClose, ScrollText, Sun } from "lucide-react";
import { ApprovalQueueSection } from "@/components/approval-queue";
import { L2ApprovalCard } from "@/components/l2-approval-card";
import { ComposerRow } from "@/components/composer-row";
import { GitBranchChip, GitBranchUnread } from "@/components/git-branch";
import { HistoryBrowser } from "@/components/history-page";
import { MonitorPopover, type MonitorAgentRow } from "@/components/harness/monitor-popover";
import {
  PaletteScreen,
  type PaletteDiscoveryPhase,
  type PaletteForceChip,
  type PaletteForceView,
  type PaletteGroup,
  type PaletteItem,
  type PaletteNote,
  type PaletteRosterRow,
  type PaletteToolRow,
} from "@/components/palette-screen";
import { LineSidebar } from "@/components/react-bits/line-sidebar";
import { ResultHistory } from "@/components/result-history";
import { RightRail } from "@/components/harness/right-rail";
import { SettingsPage } from "@/components/harness/settings-page";
import { PanelSidebar } from "@/components/panel-sidebar";
import { useDismiss } from "@/lib/use-dismiss";
import { currentView, type PanelViewId } from "@/lib/panel-views";
import {
  parseRulesForceReply,
  parseToolsListReply,
} from "@/lib/tools-discovery-view";
import {
  requestHistoryPage,
  requestRosterCancel,
  requestRulesForce,
  requestToolLaunch,
  requestToolsList,
  sendMessage,
} from "@/lib/panel";
import { TIER_NAMES as libTierNames } from "@/lib/tier-names";
import type { ComposerState, GitView, PanelSnapshot, RunSection, TaskRosterSection } from "@/lib/panel";

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

/** Native mode names -> the Chinese the tier rows show. Moved to
    lib/tier-names.ts (ticket 40) so the settings page's 安全 category shows
    the SAME words from the SAME map instead of growing a second copy that
    could drift from this one - a fourth tier cannot be invented twice
    differently. */
const TIER_NAMES = libTierNames;

/** The task monitor's roster rows, straight off the wire (ticket 197/220: the
    fields are copied verbatim - blockedOnApproval is the host's own join,
    status is a D43 name this file neither judges nor fills in, and the row's
    stream triple rides as-is (the popover draws it as a suffix, never a
    status); the display words live in the popover that draws them). Until
    this derivation existed the monitor was fed a bare empty array and every
    pushed row died unseen. An absent tasks key stays distinguishable from an
    empty roster: the popover is told which one it got (lib/panel.ts's
    optionality rule).

    Readface wave 2 (2026-11-17): the three task-card keys (ticket 38 AC#4)
    ride verbatim too - startedAt / costToDate / currentStep are omitempty on
    the wire, so an unstamped row stays ABSENT here and the popover renders
    its 未判, never an epoch, a zero cost, or a fake step. */
export function monitorAgents(tasks: TaskRosterSection | undefined): MonitorAgentRow[] {
  return (tasks?.rows ?? []).map((row) => ({
    id: row.taskId,
    name: row.label || row.taskId,
    blockedOnApproval: row.blockedOnApproval,
    status: row.status,
    statusKnown: row.statusKnown,
    statusReason: row.statusReason,
    streamTruncated: row.streamTruncated,
    streamElidedRunes: row.streamElidedRunes,
    streamDropped: row.streamDropped,
    startedAt: row.startedAt,
    costToDate: row.costToDate,
    currentStep: row.currentStep,
  }));
}

/** Narrower than this and the sidebar stops fitting a session title, so it
    collapses to its icon rail. 640px is the DEFAULT panel width, not an edge
    case - see the header. */
const SIDEBAR_COLLAPSE_BELOW = 860;

/** Whether the run feed (ticket 145) carries anything renderable: ended tool
    rows, reasoning text, the classified error or the stuck sentence. The
    running/usage/done dimensions are the contract ruler's registered Go-only
    keys (internal/panel/composer_test.go:233) - they are NOT consulted here;
    naming the gap in the report is this wave's honest answer, not a fake
    empty. */
function runHasContent(run?: RunSection): boolean {
  if (run === undefined) return false;
  return run.tools.length > 0 || Boolean(run.reasoning) || run.error != null || run.stuck != null;
}

/** lib/panel's GitView is the declared wire shape (tickets 181/189/190;
    internal/panel/git.go ReadGitForWorkspace fills it on the pipeline leg).
    This projection is the chip's feed: kind/branch/worktrees/branches
    verbatim, no second mapping. Before this existed the chip was fed a bare
    empty object and every pushed git row died unseen (the same defect the
    monitor feed's comment records for tasks). The chip's two write callbacks
    stay no-ops on purpose: a checkout/new-worktree WRITE has no route - what
    ticket 186 landed is panel.workspace.request, a scope narrowing, not a git
    operation - and a button wired to a door that does not exist is the lie
    the render-nav gate names. */
function gitChipView(git: GitView): {
  view: { current: string; isRepo: boolean; local: readonly string[]; worktrees: readonly { branch: string; path: string }[] };
} {
  return {
    view: {
      current: git.isDetached ? git.detachedSha : git.branch,
      isRepo: git.kind === "repo",
      local: git.branches,
      worktrees: git.worktrees.map((w) => ({ branch: w.detached ? w.sha : w.branch, path: w.path })),
    },
  };
}

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
    /** Ticket 38: the palette's free text. An ephemeral view preference of
        the same class as monitorOpen - it dies with the document and the
        only thing it ever feeds is the send row's own submit. */
    paletteQuery: string;
    /** The send row's failure line (lib/panel's sendRequest throws when no
        host is attached). Rendered inside the palette modal, never
        swallowed - a request that never left would otherwise look like a
        palette that worked. */
    paletteError: string;
    /** Readface wave 2 (2026-11-17): the history browser's overlay. The
        same ephemeral view-preference class as monitorOpen/paletteOpen -
        the LIST's data lives in the browser's own state fed by
        history.page, never cached here. */
    historyOpen: boolean;
    /** Readface wave 3 (2026-11-17, ticket 38 AC#3/AC#5): the palette's two
        read faces and the per-task chip. All ephemeral view preference -
        the tools.list / rules.force data is refetched every open, never
        cached across one; the chip's word dies with the document. */
    paletteToolsPhase: PaletteDiscoveryPhase;
    paletteToolsRows: PaletteToolRow[];
    paletteToolsProse: string;
    paletteToolsNote: string;
    paletteLaunchNote: string;
    paletteForcePhase: PaletteDiscoveryPhase;
    paletteForceView?: PaletteForceView;
    paletteForceNote: string;
    paletteForceChip: PaletteForceChip;
    sessionId: string | null;
    activeTurn: string | null;
  }>({
    view: null,
    railOpen: false,
    settingsOpen: false,
    settingsSection: "通用",
    dark: false,
    monitorOpen: false,
    paletteOpen: false,
    paletteQuery: "",
    paletteError: "",
    historyOpen: false,
    paletteToolsPhase: "pending",
    paletteToolsRows: [],
    paletteToolsProse: "",
    paletteToolsNote: "",
    paletteLaunchNote: "",
    paletteForcePhase: "pending",
    paletteForceNote: "",
    paletteForceChip: "",
    sessionId: null,
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
  //
  // Ticket 38 AC#6 completes the vocabulary:
  //  - Tab, while the palette is open, is TRAPPED inside the dialog: the
  //    focusables of [role="dialog"] are walked in DOM order (search box,
  //    then the rows as buttons), the ends wrap instead of escaping, and a
  //    focus that somehow landed outside is pulled back into the first
  //    focusable. Shift+Tab walks the same ring backwards.
  //  - Whatever held focus when the palette OPENED (the header button, or
  //    whatever Ctrl+K fired from) gets it BACK when the palette closes -
  //    a keyboard-only user is not dropped at the top of the document.
  //  - While an IME composition is open every key belongs to it (AC#2), so
  //    the guard returns before any palette verb - Ctrl+K, Tab, Escape -
  //    can fire mid-composition.
  const paletteReturnFocusRef = useRef<HTMLElement | null>(null);
  function rememberPaletteReturnFocus() {
    paletteReturnFocusRef.current =
      document.activeElement instanceof HTMLElement ? document.activeElement : null;
  }
  useEffect(() => {
    function onKey(event: KeyboardEvent) {
      if (event.isComposing) return;
      if ((event.ctrlKey || event.metaKey) && event.key.toLowerCase() === "k") {
        event.preventDefault();
        if (!local.paletteOpen) rememberPaletteReturnFocus();
        setLocal((c) => ({
          ...c,
          paletteOpen: !c.paletteOpen,
          // Readface wave 3: reopening refetches both read faces - reset the
          // phases here (an event handler, not a cascading effect), so the
          // sections say 正在读取 instead of showing the previous open's
          // data for a tick.
          paletteToolsPhase: "pending",
          paletteForcePhase: "pending",
          paletteLaunchNote: "",
        }));
        return;
      }
      if (event.key === "Tab" && local.paletteOpen) {
        const dialog = document.querySelector<HTMLElement>('[role="dialog"][aria-label="命令面板"]');
        if (dialog) {
          const focusables = Array.from(
            dialog.querySelectorAll<HTMLElement>("input, button:not([disabled])"),
          );
          if (focusables.length > 0) {
            const first = focusables[0];
            const last = focusables[focusables.length - 1];
            const current = document.activeElement;
            if (!dialog.contains(current)) {
              // Stray focus cannot linger outside an open modal: pull it in.
              event.preventDefault();
              first.focus();
            } else if (event.shiftKey && current === first) {
              event.preventDefault();
              last.focus();
            } else if (!event.shiftKey && current === last) {
              event.preventDefault();
              first.focus();
            }
          }
        }
        return;
      }
      if (event.key === "Escape") {
        setLocal((c) => ({ ...c, paletteOpen: false, monitorOpen: false, historyOpen: false }));
        if (narrow) setSidebarOpen(false);
      }
    }
    document.addEventListener("keydown", onKey);
    return () => document.removeEventListener("keydown", onKey);
  }, [narrow, local.paletteOpen]);

  // Closing the palette (any path: Escape, Ctrl+K again, backdrop click, a
  // dispatched row) hands focus back to the element that opened it.
  useEffect(() => {
    if (local.paletteOpen) return;
    const back = paletteReturnFocusRef.current;
    if (!back) return;
    paletteReturnFocusRef.current = null;
    back.focus();
  }, [local.paletteOpen]);

  /* ---------------------------------------------------------------------------
     Readface wave 3 (ticket 38 AC#3/AC#5): the palette's two read faces.
     Every open refetches both doors through lib/panel's own wire functions -
     nothing is cached across an open, and a null reply (no host attached)
     lands as the section's own nohost line, never an empty roster dressed up
     as an answer. The force read rides the CURRENT query: rules.force's
     "next prompt" reading is computed against the task text the user is
     about to send, so it re-reads as the query moves (a seq ref keeps a
     slow reply from overwriting a fresh one).
     --------------------------------------------------------------------------- */
  const toolsReqRef = useRef(0);
  const forceReqRef = useRef(0);
  useEffect(() => {
    if (!local.paletteOpen) return;
    const seq = ++toolsReqRef.current;
    requestToolsList().then((reply) => {
      if (seq !== toolsReqRef.current) return;
      if (reply === null) {
        setLocal((c) => ({ ...c, paletteToolsPhase: "nohost", paletteToolsRows: [], paletteToolsProse: "", paletteToolsNote: "" }));
        return;
      }
      const state = parseToolsListReply(reply);
      if (state.phase === "roster") {
        setLocal((c) => ({
          ...c,
          paletteToolsPhase: "roster",
          paletteToolsRows: state.rows,
          paletteToolsProse: state.prose,
          paletteToolsNote: "",
        }));
      } else if (state.phase === "payload-error") {
        setLocal((c) => ({ ...c, paletteToolsPhase: "payload-error", paletteToolsRows: [], paletteToolsProse: "", paletteToolsNote: state.message }));
      } else {
        setLocal((c) => ({ ...c, paletteToolsPhase: state.phase, paletteToolsRows: [], paletteToolsProse: "", paletteToolsNote: state.raw }));
      }
    });
  }, [local.paletteOpen]);
  useEffect(() => {
    if (!local.paletteOpen) return;
    const seq = ++forceReqRef.current;
    requestRulesForce(local.paletteQuery).then((reply) => {
      if (seq !== forceReqRef.current) return;
      if (reply === null) {
        setLocal((c) => ({ ...c, paletteForcePhase: "nohost", paletteForceView: undefined, paletteForceNote: "" }));
        return;
      }
      const state = parseRulesForceReply(reply);
      if (state.phase === "view") {
        setLocal((c) => ({ ...c, paletteForcePhase: "view", paletteForceView: state.view, paletteForceNote: "" }));
      } else if (state.phase === "payload-error") {
        setLocal((c) => ({ ...c, paletteForcePhase: "payload-error", paletteForceView: undefined, paletteForceNote: state.message }));
      } else {
        setLocal((c) => ({ ...c, paletteForcePhase: state.phase, paletteForceView: undefined, paletteForceNote: state.raw }));
      }
    });
  }, [local.paletteOpen, local.paletteQuery]);

  function toggleTheme() {
    const next = !local.dark;
    setLocal((c) => ({ ...c, dark: next }));
    document.documentElement.dataset.theme = next ? "dark" : "light";
  }

  // The permission tiers come from the snapshot's own vocabulary. The panel
  // asks; it never records the answer, so nothing here is optimistic.
  const tiers = composer.mode.names.length > 0 ? composer.mode.names : ["unknown"];

  /* ---------------------------------------------------------------------------
     Ticket 38: the palette's content, all of it fed, none of it remembered.
     --------------------------------------------------------------------------- */
  const paletteQueryTrim = local.paletteQuery.trim();

  // The free-text send row (D9's noise-environment fallback): the query, one
  // row, always first when there is text. Its select handler dispatches the
  // SAME envelope the composer's send button uses (panel.message.send via
  // lib/panel's sendMessage), so palette text reaches the agent exactly the
  // way composer text does - no second path, no bypass.
  const paletteGroups: PaletteGroup[] = [];
  if (paletteQueryTrim !== "") {
    paletteGroups.push({
      name: "发送任务",
      items: [
        { label: `以任务发送：${paletteQueryTrim}`, icon: "terminal", hint: "Enter", action: "send" },
      ],
    });
  }
  paletteGroups.push({
    name: "面板",
    items: [
      { label: "打开设置", icon: "settings", action: "settings" },
      { label: local.dark ? "切到浅色主题" : "切到深色主题", icon: local.dark ? "sun" : "moon", action: "theme" },
      { label: "任务监控", icon: "activity", action: "monitor" },
      { label: "新建会话", icon: "message-square-plus", action: "new-task" },
      { label: local.railOpen ? "收起右侧栏" : "展开右侧栏", icon: "panel-right", action: "rail" },
    ],
  });

  // The command roster: the snapshot's commands section (ticket 213 交付 3),
  // read verbatim - until this wiring existed every pushed roster row died
  // unseen (lib/panel.ts declared the shape; nothing on the page consumed
  // it). The three states lib/panel.ts's optionality rule separates stay
  // separate: key absent (no reader), reason set (the roster table refused
  // to assemble), rows present (the real menu).
  const commands = snapshot.commands;
  const paletteRoster: readonly PaletteRosterRow[] = commands?.rows ?? [];
  const paletteNotes: PaletteNote[] = [];
  if (commands === undefined) {
    paletteNotes.push({
      word: "命令名册 未接线",
      reason:
        "快照没有 commands 键：面板装配没接名册读取器（PumpSources.Commands，internal/panel/pump.go:221）",
    });
  } else if ((commands.reason ?? "") !== "") {
    paletteNotes.push({ word: "命令名册 装配拒绝", reason: commands.reason ?? "" });
  }
  paletteNotes.push(
    {
      word: "命令执行 未接线",
      reason:
        "名册只承载不执行（internal/panel/commands_roster_213.go:21-26），C17 名册（internal/panel/bridge.go:41-107）没有执行通路，名册行不可点",
    },
    {
      // Readface wave 2 (2026-11-17): roster.cancel landed on the C17 roster,
      // so this is no longer dormant - the stop button lives on the task
      // monitor's rows (monitor-popover.tsx), behind its own confirm card.
      word: "取消任务 已接线",
      reason:
        "走 roster.cancel（internal/panel/bridge.go 的取消门，票 38 AC#4）：任务监控浮层的每一行有停止键，先确认再取消；回执逐字上屏，已落地步骤读数走 history.page 的取证行。",
    },
    {
      // Readface wave 3 (2026-11-17): the force doors landed on the C17
      // roster (internal/panel/bridge.go's MethodRulesForce, ticket 38
      // AC#5's Go half), so this is no longer dormant - the read face
      // reports the effective rules and the chip rides the send envelope.
      word: "模式强制 已接线",
      reason:
        "走 rules.force（internal/panel/bridge.go 的 force 规则门，票 38 AC#5）：本面板只读显示生效规则＋来源（override/toml/default）＋下一 prompt 采用 mode；per-task 覆盖搭 panel.message.send 的 forceMode 顺封套，由 Go 侧 admitTask 消费——本页不预判档位、不置 allow。",
    },
    {
      // Readface wave 3 (2026-11-17): the discovery face landed too - the
      // palette's 工具名册 section reads tools.list and its launch keys go
      // through tool.launch, the real gate. The Tier-2 plugin inventory
      // keeps its own composer plus-menu face (snapshot.plugins).
      word: "第三方工具清单",
      reason:
        "内置工具走本面板的工具名册（tools.list 只读名册＋tool.launch 真门，票 38 AC#3）；Tier-2 插件清单走 composer 加号菜单（snapshot.plugins，internal/panel/plugins_inventory_215.go 的读取面，composer-row.tsx 已接）",
    },
  );

  // The palette's one execution face. Everything except "send" is a view
  // preference this document owns; "send" is the single bridge request, and
  // its failure lands on the palette's error line instead of vanishing.
  // Readface wave 3 (ticket 38 AC#5): the per-task chip's word rides the
  // SAME panel.message.send envelope as the forceMode payload key when one
  // is selected; off sends the envelope byte-identically to its pre-chip
  // shape. The Go side judges the word and consumes it at admission - this
  // page never stages, clears or reads back an override.
  function runPaletteItem(item: PaletteItem) {
    switch (item.action) {
      case "send": {
        const body = paletteQueryTrim;
        if (body === "") return;
        try {
          if (local.paletteForceChip === "") {
            sendMessage(body, []);
          } else {
            sendMessage(body, [], local.paletteForceChip);
          }
        } catch (err) {
          const detail = err instanceof Error ? err.message : String(err);
          setLocal((c) => ({ ...c, paletteError: `消息没有发出去：${detail}` }));
          return;
        }
        // The chip is per task: it applied to THIS send (the Go side's
        // admitTask consumed it there), so it resets instead of leaking
        // into whatever the user types next.
        setLocal((c) => ({ ...c, paletteOpen: false, paletteQuery: "", paletteError: "", paletteForceChip: "" }));
        return;
      }
      case "settings":
        setLocal((c) => ({ ...c, paletteOpen: false, settingsOpen: true }));
        return;
      // Theme stays open on purpose: the palette overlays the page, so the
      // change is visible the moment it lands.
      case "theme":
        toggleTheme();
        return;
      case "monitor":
        setLocal((c) => ({ ...c, paletteOpen: false, monitorOpen: true }));
        return;
      case "new-task":
        setLocal((c) => ({ ...c, paletteOpen: false, sessionId: null, view: "chat" }));
        return;
      case "rail":
        setLocal((c) => ({ ...c, railOpen: !c.railOpen }));
        return;
      default:
        return;
    }
  }

  // Readface wave 3 (ticket 38 AC#3): the tool row's launch key. The page
  // sends tool.launch through lib/panel's ONE wire function and shows the
  // receipt or refusal verbatim; whether an L2 confirmation card appears is
  // the gate's decision on the host's own path - the pill never gates this
  // button, no verdict is predicted, and no allow vocabulary rides along.
  function launchTool(row: PaletteToolRow) {
    requestToolLaunch(row.name).then((reply) => {
      setLocal((c) => ({ ...c, paletteLaunchNote: reply ?? "没接到宿主：tool.launch 没有送达，工具未被触碰" }));
    });
  }

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
          {/* git branch. Ticket 181/190's composer.git section is the feed:
              a repo shows its real branch; a non-repo shows nothing (the
              chip's own rule); a snapshot with no composer.git (the
              no-pipeline resident shape) shows the unread chip instead of a
              fake empty - three states, never a guessed branch name. */}
          {composer.git === undefined ? (
            <GitBranchUnread />
          ) : (
            <GitBranchChip
              view={gitChipView(composer.git).view}
              onCheckout={() => undefined}
              onNewWorktree={() => undefined}
            />
          )}
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
                  agents={monitorAgents(snapshot.tasks)}
                  rosterPresent={snapshot.tasks !== undefined}
                  poolInFlight={snapshot.tasks?.inFlightSlots}
                  poolCap={snapshot.tasks?.poolCap}
                  streamTruncated={snapshot.tasks?.streamTruncated}
                  streamElidedRunes={snapshot.tasks?.streamElidedRunes}
                  streamDroppedCount={snapshot.tasks?.droppedCount}
                  /* Readface wave 2 (ticket 38 AC#4): the cancel flow's two
                     doors, the same wispDispatch pipe the settings pages
                     use. A reply of null (no host) lands on the row's own
                     未判 line inside the popover - never a silent no-op. */
                  cancelDoors={{ cancelTask: requestRosterCancel, historyLoad: requestHistoryPage }}
                  onPin={() => undefined}
                />
              </div>
            ) : null}
          </div>

          {/* 历史 (ticket 36 AC#3, readface wave 2): the paged history
              browser over history.page. The read face landed on its own
              roster name (Q-50's panel.view.request stays deleted); the
              sidebar's history tree below is a SESSION-tree shape the read
              face does not answer, so it stays an honest empty - this
              button opens the paged list the face does answer. */}
          <button
            type="button"
            aria-label="历史记录"
            aria-expanded={local.historyOpen}
            className="flex size-7 items-center justify-center rounded-control text-ink-2 transition-colors duration-150 hover:bg-hover hover:text-ink"
            onClick={() => setLocal((c) => ({ ...c, historyOpen: !c.historyOpen }))}
            title="历史记录"
          >
            <ScrollText aria-hidden="true" size={14} strokeWidth={1.9} />
          </button>

          <kbd className="kbd hidden sm:flex">Ctrl K</kbd>
          <button
            type="button"
            aria-label="命令面板"
            className="flex size-7 items-center justify-center rounded-control text-ink-2 transition-colors duration-150 hover:bg-hover hover:text-ink"
            onClick={() => {
              // The button is the second open path; whoever used it gets
              // focus back on close (AC#6), same as the Ctrl+K path.
              if (!local.paletteOpen) rememberPaletteReturnFocus();
              setLocal((c) => ({
                ...c,
                paletteOpen: !c.paletteOpen,
                // Readface wave 3: reopening refetches both read faces -
                // reset the phases here (an event handler, not a cascading
                // effect), so the sections say 正在读取 instead of showing
                // the previous open's data for a tick.
                paletteToolsPhase: "pending",
                paletteForcePhase: "pending",
                paletteLaunchNote: "",
              }));
            }}
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
              /* Dormancy, narrowed (ticket 36 AC#3, readface wave 2): the
                 paged HISTORY list the read face answers now has its page
                 (the header's 历史 button -> HistoryBrowser over
                 history.page). What the sidebar tree feeds is a different
                 shape - a session tree, not the flat paged rows - and the
                 read face declares no session grouping, so the groups stay
                 an honest empty rather than rows reshaped to look like
                 sessions. */
              historyGroups={[]}
              onNewTask={() => setLocal((c) => ({ ...c, sessionId: null, view: "chat" }))}
              onOpenSettings={() => setLocal((c) => ({ ...c, settingsOpen: true }))}
              onSelectSession={(id) => setLocal((c) => ({ ...c, sessionId: id, view: "chat" }))}
              onToggleTheme={toggleTheme}
              settingsOpen={local.settingsOpen}
              theme={local.dark ? "dark" : "light"}
              workspaceName="一缕"
              workspacePath={composer.workspace.set ? composer.workspace.canonical : "尚未收到工作区"}
              /* Same engine gap as historyGroups: no session-list reader
                 exists on the snapshot (ticket 166 owns the feed per
                 panel-sidebar.tsx's own prop doc). */
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
              ballStyle={snapshot.ballStyle ?? null}
              ballSize={snapshot.ballSize ?? null}
              ballOpacityIdle={snapshot.ballOpacityIdle ?? null}
              ballClickThrough={snapshot.ballClickThrough ?? null}
              panelWidth={snapshot.panelWidth ?? null}
              panelFontSize={snapshot.panelFontSize ?? null}
              /* Ticket 40: the 安全 category reads the tier off the snapshot's
                 composer section, the 技能 category reads the plugin
                 inventory - both pushed every packet, both rendered here and
                 nowhere cached. */
              composer={composer}
              plugins={snapshot.plugins ?? null}
            />
          ) : (
            <>
              {/* The L2 cards live here, ABOVE everything and on every view. A
                  strong-confirmation card that vanished because you clicked
                  elsewhere would be a safety change smuggled into navigation. */}
              {/* Ticket 37 (D31): the queue depth + head ride WITH the cards on
                  every view, off the same pushed pending list - a queue you
                  cannot see is a queue you cannot answer. Zero cards = no
                  strip at all (a permanent zero widget trains you to ignore
                  it, the same argument the header count makes). */}
              {snapshot.pending.length > 0 ? (
                <ApprovalQueueSection snapshot={snapshot} />
              ) : null}
              {snapshot.pending.map((v) => (
                <L2ApprovalCard key={v.correlationId} view={v} />
              ))}

              <div className="flex min-h-0 flex-1 flex-col">
                {view === "approval" || view === "tasks" || view === "ball" ? (
                  <RetiredView id={view} />
                ) : (
                  <Conversation
                    composer={composer}
                    onJump={(id) => setLocal((c) => ({ ...c, activeTurn: id }))}
                    permission={composer.mode.current}
                    snapshot={snapshot}
                    tiers={tiers}
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
            {local.paletteError !== "" ? (
              <p className="border-b border-line bg-red-tint px-3.5 py-2 text-[12px] leading-[1.6] text-red" role="alert">
                {local.paletteError}
              </p>
            ) : null}
            <PaletteScreen
              forceChip={local.paletteForceChip}
              forceNote={local.paletteForceNote}
              forcePhase={local.paletteForcePhase}
              forceView={local.paletteForceView}
              groups={paletteGroups}
              launchNote={local.paletteLaunchNote}
              notes={paletteNotes}
              onForceChipChange={(next) => setLocal((c) => ({ ...c, paletteForceChip: next }))}
              onLaunchTool={launchTool}
              onQueryChange={(next) => setLocal((c) => ({ ...c, paletteQuery: next, paletteError: "" }))}
              onSelect={runPaletteItem}
              query={local.paletteQuery}
              roster={paletteRoster}
              toolsNote={local.paletteToolsNote}
              toolsPhase={local.paletteToolsPhase}
              toolsProse={local.paletteToolsProse}
              toolsRows={local.paletteToolsRows}
            />
          </div>
        </div>
      ) : null}

      {/* 历史 (ticket 36 AC#3): the paged browser as an overlay dialog, the
          palette's shape. The door's receipt, refusals, empty state and
          privacy strip live inside HistoryBrowser; the window owns only the
          backdrop and the dismiss paths (backdrop click, Escape). */}
      {local.historyOpen ? (
        <div
          aria-label="历史记录"
          aria-modal="true"
          className="fixed inset-0 z-40 flex items-start justify-center bg-black/35 pt-[10vh]"
          onClick={(e) => {
            if (e.target === e.currentTarget) setLocal((c) => ({ ...c, historyOpen: false }));
          }}
          role="dialog"
        >
          <div
            className="w-[680px] max-w-[92vw] overflow-hidden rounded-card border border-line bg-page shadow-overlay"
            style={{ animation: "fade-up 200ms var(--ease-out-strong) both" }}
          >
            <div className="max-h-[80vh] overflow-y-auto p-4">
              <HistoryBrowser load={requestHistoryPage} />
            </div>
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
  composer,
  onJump,
}: {
  snapshot: PanelSnapshot;
  tiers: string[];
  permission: string;
  /** The snapshot's composer section: the write targets (the attachments a
      send carries, the workspace in force) ride with it into ComposerRow. */
  composer: ComposerState;
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
          <div className="mx-auto flex w-full max-w-[720px] flex-1 flex-col px-6 py-5">
            {/* Ticket 36's region. The placeholder stays the empty state of the
                conversation column (the render-nav gate pins its sentences);
                ResultHistory takes over the moment anything arrives - chunks
                (sanitized markdown / caret) or run-feed rows (tools, error,
                stuck, reasoning). runHasContent distinguishes "run key present
                but nothing recorded" from content, so a fresh pipeline with
                zero events still shows the placeholder, not a bare region. */}
            {snapshot.results.length === 0 && !runHasContent(snapshot.run) ? (
              <ConversationPlaceholder />
            ) : (
              <ResultHistory results={snapshot.results} run={snapshot.run} />
            )}
          </div>
        </div>
      </div>
      <div className="shrink-0 border-t border-line bg-surface px-6 pb-4 pt-2.5">
        <div className="mx-auto w-full max-w-[720px]">
          {/* Ticket 187's UI half lives in ComposerRow: the model and
              thinking-tier pickers ride the snapshot's modelSwitch section,
              and their writes go out through the settings door. The four
              composer callbacks - permission tier, workspace, attachment,
              send - are wired to lib/panel's routes inside ComposerRow, which
              is also where the write receipts and the errors are shown. */}
          <ComposerRow
            composer={composer}
            modelSwitch={snapshot.modelSwitch}
            permission={permission}
            permissions={tiers.map((t) => ({
              key: t,
              name: TIER_NAMES[t] ?? t,
              gated: composer.mode.l2ConfirmNames.includes(t),
            }))}
            plugins={snapshot.plugins}
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
