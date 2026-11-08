/* ============================================================================
   Panel root: the Go -> page snapshot receiver (ticket 329 AC#5).
   ----------------------------------------------------------------------------
   Why this file exists: main.tsx used to render <App /> with no snapshot and
   App held no state for one, so every push the host sent died on an undefined
   window.wispOnPanelSnapshot and the panel sat on EMPTY forever. The receiver
   lives HERE, at the mount point, and not inside App on purpose: App is
   pinned by its render-nav ruler to exactly ONE state object (the view pick
   plus its view preferences), and "the latest package the host pushed" is
   not a view preference - it is the document-lifetime feed the whole panel
   renders from (PLAN.md:1044: the page holds no data of its own; this state
   is the host's data arriving, not a second store - nothing here caches,
   merges, or survives a reload).

   Contract semantics (.scratch/wisp/probes/329/r1/page-contract.md):
     - one call per publish; the argument is already the snapshot OBJECT
       (the host pasted the JSON verbatim into the call), so the state holds
       what the host sent - never a string to parse, never a partial patch;
     - WHOLE replacement: setSnapshot(s) with the new package, no diff, no
       merge - calling it any number of times just makes the latest one win;
     - cold start: this global is defined only once React has mounted. A push
       that arrives before that is dropped by WebView2 and the page does NOT
       invent a re-push or a catch-up read - cold-package semantics belong to
       tickets 33/35. Until the first push lands, <App /> renders its own
       EMPTY frame, which says so out loud.
   ============================================================================ */

import { useEffect, useState } from "react";
import App from "@/App";
import { installPanelSnapshotReceiver, type PanelSnapshot } from "@/lib/panel";

export function PanelRoot() {
  const [snapshot, setSnapshot] = useState<PanelSnapshot | null>(null);
  // Install once per mount; the uninstaller covers StrictMode's extra
  // mount/unmount cycle and a real unmount (the global dies with the page
  // anyway, but leaving it pointing at a dead setter would be a lie).
  useEffect(() => installPanelSnapshotReceiver(window, setSnapshot), []);
  return <App snapshot={snapshot ?? undefined} />;
}
