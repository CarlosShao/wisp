/* ============================================================================
   Panel root (ticket 77; navigation per Q1 = 甲).

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
   ============================================================================ */

import { Composer } from "@/components/composer";
import { L2ApprovalCard } from "@/components/l2-approval-card";
import { PanelSkeleton } from "@/components/panel-skeleton";
import { ResultStream } from "@/components/result-stream";
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
    <p className="text-[12px] text-ink-3">
      {v.label}屏还没有接到数据。面板不许自己造内容，所以这一屏只说明它缺
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
  const view = currentView(snapshot);
  return (
    <PanelSkeleton
      active={view}
      pendingCount={snapshot.pending.length}
      waitingLabel={snapshot.pending.length > 0 ? "等待确认" : undefined}
    >
      {snapshot.pending.map((v) => (
        <L2ApprovalCard key={v.correlationId} view={v} />
      ))}
      {view === "chat" && <ResultStream chunks={snapshot.results} />}
      {view !== "chat" && view !== "approval" && <UnfedScreen id={view} />}
      <Composer state={composer} />
    </PanelSkeleton>
  );
}
