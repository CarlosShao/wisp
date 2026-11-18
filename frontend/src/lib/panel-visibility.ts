/* ============================================================================
   Panel visibility pause (ticket 36 AC#5).
   ----------------------------------------------------------------------------
   WHAT THIS FILE IS. The panel's single answer to the ticket's frozen rule
   "All loops respect the CPU rule: no animations when panel hidden
   (visibilitychange pause)" (36-result-history-panel.md:32) - one mechanism
   for every JS-driven animation work source the panel has, subscribed to the
   document once per page.

   THE MAINSTREAM SHAPE, AND WHY. A browser already parks what it can see as
   presentation on a hidden document: zero requestAnimationFrame callbacks
   fire and CSS animations stop advancing. The CPU-rule gap is the JS work
   those primitives cannot see as animation:
     - setTimeout-driven reveal ladders (RevealText's 46ms segment chain)
       keep waking a hidden document's event loop - throttled, but waking;
     - any frame-driven meter or settle loop would keep scheduling frames
       from JS if nobody stopped it.
   So the mechanism, exactly the mainstream shape the dispatch names
   (document.visibilityState 变化 => 停 rAF 循环、恢复再启动):

     - document.visibilityState is the ONE fact. Read directly through
       panelWorkAllowed - never inferred from a timestamp, never guessed
       from a frame count (the wall-clock ban).
     - the `visibilitychange` event is the ONE signal. installPanelVisibilityPause
       subscribes once (panel-root.tsx, at the mount point); on hidden every
       registered loop stops (its pending rAF frame is CANCELLED - the frozen
       state the instrument asserts); on visible again the loops that were
       running when the panel hid restart from where they stopped - no catch-up
       burst, no lost frame. The resume half is asserted by counters, not by
       the clock.
     - RevealText's ladder gates its scheduling on the same one fact and
       resumes on the same one signal (usePanelVisibilityTick): while hidden
       no timer is scheduled at all, so the wake source is gone, and the
       segment counter freezes mid-stream, then continues on show.

   WHY A rAF LOOP FACTORY AT ALL. The counters the AC demands ride rAF: the
   instrument (scripts/render-hiddenpause36.tsx) feeds this file's own
   createPanelAnimationLoop through a parser-level browser-class engine
   (schedule -> run -> schedule until cancelled) and judges by counter
   deltas: visible => grows, hidden => the same tick sampled twice is equal
   and further pumps add nothing, visible again => grows again. Production
   code under test, mutation teeth both red (no pause => grows while hidden;
   no resume => frozen after visible), zero wall-clock judgement anywhere.

   SCOPE HONESTY (the browser's own share is not faked). The CSS-driven
   animations (the chunk fade-up, the stream caret, the ball rings) pause on
   a hidden page because the compositor stops - that half needs no code and
   this file does not pretend to add it. What this file owns is exactly the
   JS half the compositor cannot see.
   ============================================================================ */

import { useEffect, useState } from "react";

/** Whether the panel's animation work may currently run. The one fact, read
    off the document - never derived from a timestamp (no wall-clock judging).
    A null document (no DOM, e.g. a static render pass) allows no work: the
    honest answer for "nowhere to read the fact from" is the paused one. */
export function panelWorkAllowed(doc: Document | null = typeof document === "undefined" ? null : document): boolean {
  if (doc === null) return false;
  return doc.visibilityState === "visible";
}

/** One pausable animation work source. Both shapes the panel has - an rAF
    loop and RevealText's segment ladder - present as this: stop() takes the
    work off the event loop (the pending frame/timer is cancelled), start()
    resumes it, running says whether work is currently scheduled or mid-step.
    start() after start() is a no-op, stop() after stop() is a no-op. */
export interface PausableWork {
  start(): void;
  stop(): void;
  readonly running: boolean;
}

/** Build the panel's rAF loop over a target document: the callback runs every
    frame while running. Each start() schedules a fresh chain of
    requestAnimationFrame frames; each stop() cancels the pending frame and
    the chain dies there. This is the loop the counters count. */
export function createPanelAnimationLoop(
  doc: Document,
  onFrame: (frame: number) => void,
): PausableWork {
  let id: number | null = null;
  let isRunning = false;
  const view = doc.defaultView;
  if (view === null) {
    throw new Error("createPanelAnimationLoop: the document has no view; no frame scheduler exists to drive a loop");
  }
  const tick = (frame: number): void => {
    if (!isRunning) return;
    onFrame(frame);
    id = view.requestAnimationFrame(tick);
  };
  return {
    start(): void {
      if (isRunning) return;
      isRunning = true;
      id = view.requestAnimationFrame(tick);
    },
    stop(): void {
      if (!isRunning) return;
      isRunning = false;
      if (id !== null) {
        view.cancelAnimationFrame(id);
        id = null;
      }
    },
    get running(): boolean {
      return isRunning;
    },
  };
}

/** The page-level pause subscription. On hidden: every loop that is running
    stops (frozen) and is remembered in the suspended set. On visible again:
    exactly the suspended set restarts - a loop that was already stopped when
    the panel hid is NOT started by the resume (the pause never invents work),
    and a loop that was running is never left stopped (the resume half the
    instrument's MUT-B reddens). The uninstaller rides the returned handle
    (`handle.uninstall()`), the shape panel-root.tsx's cleanup returns;
    `paused` is the instrument-visible state the listener last drove. */
export function installPanelVisibilityPause(
  doc: Document,
  loops: PausableWork[],
): { uninstall(): void; readonly paused: boolean } {
  let paused = false;
  const suspended = new Set<PausableWork>();
  const onVisibility = (): void => {
    if (panelWorkAllowed(doc)) {
      // MUTATION-36B (the "stopped but never resumed" tooth): delete exactly
      // the start loop below, keep the clear - hidden still freezes, visible
      // never restarts, and the instrument's post-resume counter delta reddens.
      for (const loop of suspended) loop.start();
      suspended.clear();
      paused = false;
      return;
    }
    // MUTATION-36A (the "hidden does not pause" tooth): delete exactly the
    // stop loop below - hidden flips `paused` but no loop ever stops, and the
    // instrument's hidden-state counter delta reddens.
    for (const loop of loops) {
      if (loop.running) {
        loop.stop();
        suspended.add(loop);
      }
    }
    paused = true;
  };
  doc.addEventListener("visibilitychange", onVisibility);
  return {
    uninstall(): void {
      doc.removeEventListener("visibilitychange", onVisibility);
    },
    get paused(): boolean {
      return paused;
    },
  };
}

/** The panel page's one registry. The mount point (panel-root.tsx) subscribes
    THIS array to the document once; frame-driven sources on the panel path
    register their loop here instead of owning a visibility listener of their
    own - one subscription per page, N loops on it. */
export const panelAnimationLoops: PausableWork[] = [];

/** A counter that bumps on every visibilitychange of the document. The
    resume signal for scheduling effects (RevealText's ladder): a hidden
    bump re-runs the effect, the gate inside declines to schedule; the
    visible bump re-runs it, the gate schedules again. Effects that only
    read the fact at schedule time stay honest across the transition without
    owning a listener. The document defaults to the page's own, guarded the
    same way panelWorkAllowed guards: a no-DOM render pass (SSR, the render
    instruments) gets no listener and a stable tick. */
export function usePanelVisibilityTick(
  doc: Document | undefined = typeof document === "undefined" ? undefined : document,
): number {
  const [tick, setTick] = useState(0);
  useEffect(() => {
    if (doc === undefined) return;
    const on = (): void => setTick((t) => t + 1);
    doc.addEventListener("visibilitychange", on);
    return () => doc.removeEventListener("visibilitychange", on);
  }, [doc]);
  return tick;
}
