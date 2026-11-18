/* ============================================================================
   Ticket 36 AC#5's behavioural tooth (render + rule evidence), 2026-11-18.
   ----------------------------------------------------------------------------
   WHAT THIS FILE PROVES. The AC's sentence is "Hidden-panel animation pause
   asserted (rAF counters)" (36-result-history-panel.md:47), and the dispatch
   spells the judgement: 可见态 => 计数增长；切 hidden => 计数停住（同一 tick
   多次采样相等）；切回 => 恢复增长。⛔ 墙钟猜测，用计数器差值判. So every
   reading below is a COUNTER DELTA:

     - the engine: a parser-level, browser-class rAF implementation (this
       tree's instrument convention - the tree ships no jsdom), driving
       createPanelAnimationLoop from lib/panel-visibility.ts: schedule ->
       run -> schedule until cancelled, requestAnimationFrame/
       cancelAnimationFrame on the injected view, one pump = one frame's
       work. Because a real hidden document never fires rAF, the engine
       models the compositor's stop - but the PAUSE itself is judged on the
       production loop's own state: while hidden, stop() must have taken the
       loop off the scheduler, so pumps add zero ticks (the counter delta is
       the verdict, not the engine's model).
     - the resume: the same one fact drives installPanelVisibilityPause's
       visible branch, which restarts exactly the suspended set.

   THE TEETH (the dispatch's two, both mutation-proven):
     ① 隐藏不暂停 (MUT-A: the hidden branch's stop loop deleted) => the
       hidden-state counter delta is > 0 => RED;
     ② 只停不恢复 (MUT-B: the visible branch's restart loop deleted) => the
       post-resume counter delta is 0 => RED.
   Both ran against an exported copy of this tree (git archive export, the
   same off-repo mutation bench ticket 36 leg 2 used - never the real repo);
   readings in docs/evidence/s1/36-hiddenpause-deflake39-20261118.md.

   The reveal-ladder half (RevealText riding the same fact) is asserted at
   source level below: renderToStaticMarkup runs no effects, so the
   schedule-time gate cannot be exercised in markup - its WIRING (the gate
   in the effect, the visibility tick in the deps) is pinned instead, and
   the rAF loop above is the exercised proof of the shared mechanism.
   panel-root's subscription and App's tracker span are pinned the same way.
   No wall clock is read anywhere in this file.
   ============================================================================ */

import { readFileSync } from "node:fs";
import { renderToStaticMarkup } from "react-dom/server";
import { PanelHiddenPauseProbe } from "@/App";
import {
  createPanelAnimationLoop,
  installPanelVisibilityPause,
  panelWorkAllowed,
  type PausableWork,
} from "@/lib/panel-visibility";

const failures: string[] = [];
function check(label: string, ok: boolean, detail = ""): void {
  if (!ok) failures.push(`${label}${detail ? ` [${detail}]` : ""}`);
}

/* ---------------------------------------------------------------------------
   The parser-level browser-class document/view: the mainstream fake the
   panel-root subscription and the instrument share - an EventTarget that
   models visibilitychange, and a view whose rAF is a real handle registry
   (cancel actually takes the pending frame off; pumps are the frames; no
   wall clock anywhere in the judgement).
   --------------------------------------------------------------------------- */

class FakeView {
  private scheduled = new Map<number, (frame: number) => void>();
  private seq = 0;
  requestAnimationFrame(cb: (frame: number) => void): number {
    const handle = ++this.seq;
    this.scheduled.set(handle, cb);
    return handle;
  }
  cancelAnimationFrame(handle: number): void {
    this.scheduled.delete(handle);
  }
  /** One compositor pump: run every scheduled callback once. A pump is one
      frame's worth of work; the callbacks reschedule themselves. */
  pump(): void {
    const next = [...this.scheduled.entries()];
    this.scheduled.clear();
    let frame = 0;
    for (const [, cb] of next) cb(++frame);
  }
  get pending(): number {
    return this.scheduled.size;
  }
}

class FakeDocument {
  visibilityState: "visible" | "hidden" = "visible";
  private listeners = new Map<string, ((ev: { type: string }) => void)[]>();
  readonly defaultView = new FakeView();
  addEventListener(type: string, cb: (ev: { type: string }) => void): void {
    const list = this.listeners.get(type) ?? [];
    list.push(cb);
    this.listeners.set(type, list);
  }
  removeEventListener(type: string, cb: (ev: { type: string }) => void): void {
    const list = this.listeners.get(type);
    if (!list) return;
    const i = list.indexOf(cb);
    if (i >= 0) list.splice(i, 1);
  }
  /** The one signal the mechanism subscribes to, driven exactly like the
      browser drives it: flip the fact, then fire visibilitychange. */
  setVisibility(state: "visible" | "hidden"): void {
    this.visibilityState = state;
    for (const cb of this.listeners.get("visibilitychange") ?? []) cb({ type: "visibilitychange" });
  }
}

/* ---------------------------------------------------------------------------
   The instrument's counter loop: a PausableWork over the engine, whose
   onFrame bumps the counter. The judgement is always deltas between
   samples.
   --------------------------------------------------------------------------- */

interface CounterLoop {
  work: PausableWork;
  reads(): number;
}

function counterLoop(doc: FakeDocument): CounterLoop {
  let count = 0;
  const work = createPanelAnimationLoop(doc as unknown as Document, () => {
    count += 1;
  });
  return { work, reads: () => count };
}

/* ---------------------------------------------------------------------------
   1. The exercised tooth: visible grows, hidden freezes, visible resumes
   --------------------------------------------------------------------------- */

const doc = new FakeDocument();
const counter = counterLoop(doc);
// The page-level subscription is what drives the loop on visibilitychange -
// the same wiring panel-root.tsx mounts (one listener, the loops array).
const sub1 = installPanelVisibilityPause(doc as unknown as Document, [counter.work]);

check("the fact starts visible", doc.visibilityState === "visible" && panelWorkAllowed(doc as unknown as Document));
counter.work.start();
check("the loop starts running", counter.work.running);
doc.defaultView.pump();
const visibleBefore = counter.reads();
doc.defaultView.pump();
const visibleAfter = counter.reads();
check(
  `visible: the counter grows across pumps (${visibleBefore} -> ${visibleAfter})`,
  visibleBefore > 0 && visibleAfter > visibleBefore,
  `${visibleBefore}/${visibleAfter}`,
);
check("the engine schedules the next frame while running", doc.defaultView.pending > 0);

// --- hidden: the freeze (tooth ①'s green control) --------------------------
doc.setVisibility("hidden");
check("hidden: the pause stopped the loop", !counter.work.running);
check("hidden: the pending frame is off the scheduler", doc.defaultView.pending === 0);
doc.defaultView.pump();
const hiddenSample1 = counter.reads();
doc.defaultView.pump();
const hiddenSample2 = counter.reads();
doc.defaultView.pump();
const hiddenSample3 = counter.reads();
// The dispatch's own shape: the SAME tick sampled repeatedly is equal.
check(
  `hidden: the counter is frozen (pumps 1/2/3 equal: ${hiddenSample1}/${hiddenSample2}/${hiddenSample3})`,
  hiddenSample1 === visibleAfter && hiddenSample2 === hiddenSample1 && hiddenSample3 === hiddenSample2,
  `${visibleAfter}/${hiddenSample1}/${hiddenSample2}/${hiddenSample3}`,
);

// --- visible again: the resume (tooth ②'s green control) -------------------
doc.setVisibility("visible");
check("visible again: the suspended loop restarted", counter.work.running);
check("visible again: a frame is scheduled again", doc.defaultView.pending > 0);
doc.defaultView.pump();
const resumedAfterOne = counter.reads();
check(`resume: the counter grows again (${hiddenSample3} -> ${resumedAfterOne})`, resumedAfterOne > hiddenSample3);
doc.defaultView.pump();
const resumedAfterTwo = counter.reads();
check(`resume: the growth continues (${resumedAfterOne} -> ${resumedAfterTwo})`, resumedAfterTwo > resumedAfterOne);
sub1.uninstall();

/* ---------------------------------------------------------------------------
   2. The resume membership rule: exactly the suspended set restarts
   --------------------------------------------------------------------------- */

const doc2 = new FakeDocument();
const a = counterLoop(doc2);
const b = counterLoop(doc2);
a.work.start();
b.work.start();
const sub2 = installPanelVisibilityPause(doc2 as unknown as Document, [a.work, b.work]);
doc2.setVisibility("hidden");
check("both running loops are suspended on hidden", !a.work.running && !b.work.running);
doc2.setVisibility("visible");
check("both suspended loops resume on visible", a.work.running && b.work.running);

// A loop that was stopped BEFORE the hide must not be started by the resume:
// the pause never invents work.
const c = counterLoop(doc2);
check("a never-started loop is not running", !c.work.running);
doc2.setVisibility("hidden");
doc2.setVisibility("visible");
check("the resume does not invent work for a loop that was already stopped", !c.work.running);
sub2.uninstall();

/* ---------------------------------------------------------------------------
   3. The subscription wiring: one visibilitychange listener drives the set
   --------------------------------------------------------------------------- */

const doc3 = new FakeDocument();
const w3 = counterLoop(doc3);
w3.work.start();
const sub3 = installPanelVisibilityPause(doc3 as unknown as Document, [w3.work]);
check("the subscription starts in the running state (the page is visible)", !sub3.paused);
doc3.setVisibility("hidden");
check("the subscription reports paused after the hide", sub3.paused);
doc3.setVisibility("visible");
check("the subscription reports running after the show", !sub3.paused);
check("the loop restarted through the subscription", w3.work.running);
// Uninstall: the listener dies, so a hide no longer pauses (the page-level
// contract panel-root.tsx rides: install once per mount, uninstall on unmount).
sub3.uninstall();
w3.work.start();
doc3.setVisibility("hidden");
check("after uninstall the subscription no longer drives the loop", w3.work.running);
w3.work.stop();

/* ---------------------------------------------------------------------------
   4. The App-level wiring: the tracker span rides the same one fact
   --------------------------------------------------------------------------- */

const trackerDoc = new FakeDocument();
trackerDoc.visibilityState = "hidden";
const hiddenRender = renderToStaticMarkup(<PanelHiddenPauseProbe doc={trackerDoc as unknown as Document} />);
check(
  "hidden document => the tracker reports paused",
  hiddenRender.includes(`data-panel-hidden-pause="paused"`),
  hiddenRender.slice(0, 160),
);
trackerDoc.visibilityState = "visible";
const visibleRender = renderToStaticMarkup(<PanelHiddenPauseProbe doc={trackerDoc as unknown as Document} />);
check(
  "visible document => the tracker reports running",
  visibleRender.includes(`data-panel-hidden-pause="running"`),
  visibleRender.slice(0, 160),
);

const appSrc = readFileSync("src/App.tsx", "utf8");
check(
  "App exports the PanelHiddenPauseProbe tracker that reports the fact",
  appSrc.includes("export function PanelHiddenPauseProbe") && appSrc.includes("data-panel-hidden-pause"),
);
check(
  "the tracker judges on panelWorkAllowed - the one fact, never a clock",
  appSrc.includes(`panelWorkAllowed(doc) ? "running" : "paused"`),
);
check(
  "the tracker's document is injectable (the instrument feeds the same fake the subscription rides)",
  appSrc.includes(`doc = typeof document === "undefined" ? null : document`),
);

/* ---------------------------------------------------------------------------
   5. The panel-root subscription is mounted (the wiring the loops ride)
   --------------------------------------------------------------------------- */

const rootSrc = readFileSync("src/panel-root.tsx", "utf8");
check(
  "panel-root installs the visibility pause once per mount",
  rootSrc.includes("installPanelVisibilityPause(document, panelAnimationLoops)"),
);
check(
  "panel-root's cleanup covers StrictMode's double mount (uninstall on teardown)",
  rootSrc.includes("return () => pause.uninstall();"),
);

/* ---------------------------------------------------------------------------
   6. The reveal ladder rides the same fact (the setTimeout wake source)
   --------------------------------------------------------------------------- */

const revealSrc = readFileSync("src/components/reveal-text.tsx", "utf8");
check(
  "RevealText gates its schedule on panelWorkAllowed (no timer exists while hidden)",
  revealSrc.includes("if (!panelWorkAllowed()) return;"),
);
check(
  "RevealText's effect re-runs on the visibility tick (the resume path)",
  revealSrc.includes("[count, done, visTick]"),
);
check(
  "RevealText does not own a private visibility listener (one subscription per page)",
  !revealSrc.includes("addEventListener"),
);

/* ---------------------------------------------------------------------------
   7. The production file's mutation anchors are present (the teeth's
      targets exist where the bench cut them)
   --------------------------------------------------------------------------- */

const visSrc = readFileSync("src/lib/panel-visibility.ts", "utf8");
check(
  "the MUT-A anchor is on disk (the hidden branch's stop loop)",
  visSrc.includes("MUTATION-36A") && visSrc.includes("loop.stop()"),
);
check(
  "the MUT-B anchor is on disk (the visible branch's restart loop)",
  visSrc.includes("MUTATION-36B") && visSrc.includes("loop.start()"),
);

/* ---------------------------------------------------------------------------
   The verdict.
   --------------------------------------------------------------------------- */

if (failures.length > 0) {
  for (const f of failures) console.error("render-hiddenpause36: FAIL " + f);
  console.error(`render-hiddenpause36: FAIL, ${failures.length} check(s) failed`);
  process.exit(1);
}
console.log(
  `render-hiddenpause36: OK, rAF counters judge the pause by deltas ` +
    `(visible ${visibleBefore}->${visibleAfter} grows, hidden frozen ${hiddenSample1}/${hiddenSample2}/${hiddenSample3} equal, ` +
    `resume ${hiddenSample3}->${resumedAfterTwo} grows), ` +
    `resume restarts only what it suspended, uninstall detaches, ` +
    `App tracker + panel-root subscription + reveal-ladder gate wired, mutation anchors present`,
);
