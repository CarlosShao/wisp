/* ============================================================================
   Panel root (ticket 77; navigation per Q1 = 甲, icons per owner 2026-09-25).

   State policy: NONE (PLAN.md:1044). The only data this component renders is
   the snapshot the Go side hands it, and the only way to change what is on
   screen is a new snapshot - including which of the nine views is showing
   (Q2: `view` arrives through currentView(), no component holds it). In
   production the pump is ticket 35's bridge push; with no host attached the
   component tree still renders an empty frame rather than reaching for a cached
   copy, because a WebView restart must land in the same state the Go side
   reports, not in whatever the page remembered.

   Two things deliberately do NOT move with the view: the pending L2 cards and
   the composer row. Hiding a strong-confirmation card because the user clicked
   to another screen would be a safety change smuggled into a navigation change;
   the rail's 审批 badge is what carries the queue depth instead.

   `chrome` is the harness switch: main.tsx passes it only under ?harness=1,
   and it buys the floating-window dressing a browser needs (desktop backdrop,
   fog, window radius and shadow). In production the WebView2 window IS the
   panel, so the shell renders bare - the chrome a real window gets comes from
   Win32, not from CSS.
   ============================================================================ */

import { useState } from "react";
import { ApprovalScreen } from "@/components/approval-screen";
import { ChatScreen } from "@/components/chat-screen";
import { Composer } from "@/components/composer";
import { ConfigScreen } from "@/components/config-screen";
import { L2ApprovalCard } from "@/components/l2-approval-card";
import { PanelSkeleton } from "@/components/panel-skeleton";
import { currentView, viewOf, type PanelViewId } from "@/lib/panel-views";
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

/** A screen Go has no data for yet says so, and says which ticket would fill it.
    Inventing rows here is what owner's P9 red line forbids ("不得造假数据当真实
    字段"), and an empty <div> would read as a bug rather than as a gap. */
function UnfedScreen({ id }: { id: PanelViewId }) {
  const v = viewOf(id);
  return (
    <p className="text-[12px] text-muted-foreground">
      {v.label}屏还没有接到数据。面板不许自己造内容，所以这一屏只说明它缺
      C17 的读口与票 35 的推送。
    </p>
  );
}

export default function App({
  snapshot = EMPTY,
  chrome = false,
}: {
  snapshot?: PanelSnapshot;
  chrome?: boolean;
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
  const [picked, setPicked] = useState<PanelViewId | null>(null);
  const hostView = (snapshot as { view?: string }).view;
  const view = hostView !== undefined ? currentView(snapshot) : (picked ?? currentView(snapshot));
  return (
    <PanelSkeleton
      active={view}
      chrome={chrome}
      onSelect={setPicked}
      pendingCount={snapshot.pending.length}
    >
      {/* Pending cards stay ABOVE the view switch and on every view: hiding a
          strong-confirmation card because the user clicked to another screen
          would be a safety change smuggled into a navigation change. */}
      {snapshot.pending.map((v) => (
        <L2ApprovalCard key={v.correlationId} view={v} />
      ))}
      {view === "chat" ? (
        <ChatScreen snapshot={snapshot} chrome={chrome} />
      ) : view === "approval" ? (
        <ApprovalScreen snapshot={snapshot} />
      ) : view === "config" ? (
        <ConfigScreen />
      ) : (
        <UnfedScreen id={view} />
      )}
      {/* The composer row is the demo's bottom-docked bar: mt-auto pins it to
          the window floor while the content is shorter than one screen and
          lets it scroll normally once the stream outgrows the viewport
          (panel-skeleton's main is a flex column for exactly this). */}
      <div className="mt-auto">
        <Composer state={composer} />
      </div>
    </PanelSkeleton>
  );
}
