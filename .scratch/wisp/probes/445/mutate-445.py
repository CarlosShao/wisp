#!/usr/bin/env python3
"""Ticket 445 AC#3 -- the two mutations that prove the repaired instrument still
has teeth. RUNS ON AN OUT-OF-REPO EXPORT TREE ONLY (D:/tmp/wisp/<date>/445-tree):
the main repository is never mutated by this file.

  m1  !TriggerInterrupt  the click's gesture still runs and its executor is still
                         attached, but `trigger()` is removed from
                         resident_ball_windows.go's clickInterruptGesture: the
                         ball click REALLY does not call TriggerInterrupt, i.e.
                         the "real hole" world AC#1 asks about. The repaired
                         instrument must go red and name the MUTANT line.
  m2  !edge-recorded      the walk still happens (the state machine moves to
                         Listening) but rv.dispatch no longer appends the edge to
                         rv.timeline. That is the world where "wait for the edge
                         to land" can never succeed, so the repaired instrument's
                         WAIT must go red with its own named line and numbers -
                         this is what separates "wait for the target" from "wait
                         longer": a longer wait would not help here at all.
  restore  put both files back from their `.445pristine` copies (created on the
           first mutation, never deleted) and print the sha256 of each, so
           "restored to the same bytes" is a reading rather than a claim.

Usage:
  python .scratch/wisp/probes/445/mutate-445.py <tree> m1|m2|restore
"""

import hashlib
import pathlib
import shutil
import sys

BALL = "cmd/wisp/resident_ball_windows.go"
VOICE = "cmd/wisp/resident_voice_windows.go"

M1_OLD = """\ttrigger()
\tslog.Info("ball: click interrupt dispatched", "executed", true)
"""
M1_NEW = """\t// MUTATION 445-M1: TriggerInterrupt is never called. The executor stays
\t// attached (so the gesture is still "hosted" and the ball's own nil branch does
\t// not fire) and the log still claims execution - this is the honest-defect
\t// world: the ball click and the voice leg are not wired.
\tslog.Info("ball: click interrupt dispatched", "executed", true)
"""

M2_OLD = """\trv.mu.Lock()
\trv.timeline = append(rv.timeline, voiceTransition{
\t\tFrom: from, To: to, Event: ev, AtSample: rv.fed.Load(), Say: say,
\t})
\trv.mu.Unlock()
"""
M2_NEW = """\t// MUTATION 445-M2: the edge is NOT recorded. The state walk above still
\t// happened, so every machine.State() marker still fires - only the timeline
\t// is silent. Nothing but a wait-for-the-edge instrument can catch this.
"""


def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def mutate(tree, rel, old, new):
    path = tree / rel
    pristine = path.with_suffix(path.suffix + ".445pristine")
    if not pristine.exists():
        shutil.copy2(path, pristine)
    text = path.read_text(encoding="utf-8")
    if text.count(old) != 1:
        sys.exit("REFUSED: %s: the anchor occurs %d time(s), want exactly 1" % (rel, text.count(old)))
    before, after = sha(path), None
    path.write_text(text.replace(old, new), encoding="utf-8", newline="\n")
    after = sha(path)
    print("mutated   %-46s %s -> %s" % (rel, before[:16], after[:16]))
    print("pristine  %-46s %s" % (pristine.name, sha(pristine)[:16]))


def restore(tree):
    for rel in (BALL, VOICE):
        path = tree / rel
        pristine = path.with_suffix(path.suffix + ".445pristine")
        if not pristine.exists():
            print("no pristine copy for %s (never mutated in this tree)" % rel)
            continue
        shutil.copy2(pristine, path)
        same = sha(path) == sha(pristine)
        print("restored  %-46s %s pristine_same=%s" % (rel, sha(path), same))
        if not same:
            sys.exit("RESTORE FAILED: %s" % rel)


def main():
    if len(sys.argv) != 3:
        sys.exit(__doc__)
    tree = pathlib.Path(sys.argv[1])
    if not (tree / "go.mod").exists():
        sys.exit("REFUSED: %s does not look like a Wisp tree (no go.mod)" % tree)
    mode = sys.argv[2]
    if mode == "m1":
        mutate(tree, BALL, M1_OLD, M1_NEW)
    elif mode == "m2":
        mutate(tree, VOICE, M2_OLD, M2_NEW)
    elif mode == "restore":
        restore(tree)
    else:
        sys.exit(__doc__)


if __name__ == "__main__":
    main()
