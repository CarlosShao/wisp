# 309-v3b anchor (non-implementer adjudication, cells AC#2 + AC#3 only)

- Role: judge for ticket 309. No production code, no assertion edits, no box flips, no push.
- Leg: 309-v3b, continuation of 309-v3 (died on daily usage limit; only left anchor cf627625 + logs, zero verdicts).
- My anchor: `git rev-parse --short HEAD` = 3b8066a2 (branch dev), taken 2026-10-11 19:0x +08 by me. staged = 0 lines (`git diff --cached --name-only` empty at take time).
- Worktree carries other legs' edits (do not touch); full filter `internal cmd tools scripts` in logs/anchor-worktree.txt.
- Ticket state at read: 4 boxes [x] / 2 boxes [ ] = AC#2 and AC#3 unchecked; box count re-verified by me before any verdict.
- Briefs read verbatim: 309v3b-dispatch.txt (continuation), 309v3-dispatch.txt (original, still in force).
- Predecessor leftovers (probes/309/v3/**; readings citable with attribution, any conclusion of mine must come from my own re-run):
  - 00-anchor.md, logs/anchor-state.txt, logs/ancestry.txt, logs/gitlog.txt
  - logs/r-labels.txt (shape-label ruler raw output); logs/labels.stderr = 0 bytes, I re-shoot this ruler with stderr landed
  - logs/rig-shapes.txt (committed by orchestrator as 3b8066a2)
- Commits under review: df8d517d (A: matcher widened to new/make/conversion), a1f84d39 (B: bridge.go comment fix), 6fed247e (r3 evidence). Orchestrator control shot: probes/orch/20261011/verify-304r1-309r3.md sec.3 (new(ToolRequest) bare in package agent => red; mv aside => ok; pristine => ok). That is a reference, not my proof.
- Harness-fault guard: any red whose text contains walk / GetFileAttributesEx / build failed / no test files is a void, not teeth; every red must carry the exact line that produced it.
- Rigs: out-of-repo only, one fresh mktemp -d per tree, export = git archive HEAD go.mod go.sum internal cmd tools scripts (six pieces; missing scripts => fake walk red), non-empty proof by .go count first.
- Deliverables planned: probes/309/v3b/{10-teeth-pairs.md, 20-cells.md, 90-final.md} + logs; docs/evidence/s1/309-v3b-verdict.md; one append-only tail bullet in the ticket.
