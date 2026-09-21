/* ============================================================================
   Panel root (ticket 77).

   State policy: NONE (PLAN.md:1044). The only data this component renders is
   the snapshot the Go side hands it, and the only way to change what is on
   screen is a new snapshot. In production the pump is ticket 35's bridge push;
   with no host attached the component tree still renders an empty frame rather
   than reaching for a cached copy, because a WebView restart must land in the
   same state the Go side reports, not in whatever the page remembered.
   ============================================================================ */

import { Composer } from "@/components/composer";
import { L2ApprovalCard } from "@/components/l2-approval-card";
import { PanelSkeleton, type PanelTab } from "@/components/panel-skeleton";
import { ResultStream } from "@/components/result-stream";
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

export default function App({
  snapshot = EMPTY,
  tab = "结果",
}: {
  snapshot?: PanelSnapshot;
  tab?: PanelTab;
}) {
  // A snapshot from a host that has not filled the composer section yet falls
  // back to the empty view, which says so in words; inventing a mode here would
  // be the panel guessing about a permission it was never told.
  const composer = snapshot.composer ?? EMPTY_COMPOSER;
  return (
    <PanelSkeleton active={tab} waitingLabel={snapshot.pending.length > 0 ? "等待确认" : undefined}>
      {snapshot.pending.map((view) => (
        <L2ApprovalCard key={view.correlationId} view={view} />
      ))}
      <ResultStream chunks={snapshot.results} />
      <Composer state={composer} />
    </PanelSkeleton>
  );
}
