# 331-r1 evidence (implementer, comment-only leg) — errors.go MessageKey comment tells the truth

- Leg: `331-r1` (implementer, whole-line comment surface only). Anchor recorded at start:
  `72b65c614af0a8441f0c2d31ddbe4569c42d4d20` (dev). Deliverable commits: `e09187e1` (errors.go),
  `5cbec52d` (probes). Working-tree target file verified byte-identical to anchor before editing.
- Note vs ticket face: ticket's anchor `0e944d93` is older; all readings below re-run live on
  `72b65c61`. Line numbers `:160` / `:163` and the four test call sites are unchanged at this
  anchor, so no correction needed.

## 1. The boundary this ticket exists for (dispatch requires it spelled out)
What is deferred is *building the error-copy layer* (D23 / S6, `DEFERRED` row in
`docs/specs/SPEC-12-roadmap-governance.md` section 5). Deferred work must not buy itself a
fictitious claim in a code comment: the old text "UI and logs embed the key ... so copy can be
added later without re-touching error sites" asserted a consumer base that measures **zero** in
production. This leg corrected the sentence only; it did NOT add any production caller (that is
the deferred layer's job, not this ticket's).

## 2. Startup scales (ruler named with each number; logs under .scratch/wisp/probes/331/r1/logs)
- R1 `git grep -n "MessageKey" 72b65c61 -- internal cmd tools ':(exclude)*_test.go'`, rc=0,
  stderr empty: exactly 2 hits = `errors.go:160` (own comment) + `errors.go:163` (declaration)
  => production embedders = 0 (count of call sites inside production scope).
- R2 same ruler without `:(exclude)` (positive control proving R1 is not blind), rc=0: real
  call sites present = `internal/observe/errors_test.go:191`, `:192`,
  `internal/llm/matrix_14_2_test.go:175`, `:249` (plus comment/name matches `errors_test.go:186`
  /`:188`).
- R3 `git show 72b65c61:internal/observe/errors.go | sed -n '158,166p'`, rc=0 — before text,
  verbatim (archived in `r0-scale3-comment-before.txt`):
  `// MessageKey returns the stable user-visible copy key for the class. The copy`
  `// text itself is deferred (D23 error-copy work, S6); UI and logs embed the key`
  `// so copy can be added later without re-touching error sites.`

## 3. The replacement (AC#1: three facts, no forward promise)
New comment states: (1) the key is *reserved for* the deferred copy layer; (2) as read on
2026-10-17 no production code embeds it, only tests call it (both test files named;
`TestMessageKeys` verified to exist before naming it); (3) the copy layer stays DEFERRED under
D23/S6; final line "Whether UI and logs embed the key is that deferred layer's decision, not a
promise made here" deliberately avoids any "we will always ..." guarantee (AC#5 discipline).

## 4. AC#2 zero-change ruler, its positive control, and the reverse-write experiment
- Ruler `zero-change-check.sh`: strips whole-line comments (`^[[:space:]]*//`) from anchor blob
  and candidate, then `sha256sum` + `diff`. Scope declared honestly: it measures *production
  bytes excluding whole-line comments only*.
- Main run (candidate = committed worktree file): both sides
  `f51f944c8aaf1d7a7c708a50dcabcf7408ce2cc82c470ccd4591bc408efd5539`, diff empty, rc=0.
- Line audit (`git diff 72b65c61..e09187e1`): changed lines total = 10 (3 removed + 7 added),
  non-whole-comment changed lines = 0 (`grep -vcE '^[+-][[:space:]]*//'` rc=1, zero matches).
- Positive control (mandatory, else the "empty" is not evidence): on an out-of-tree copy, change
  one literal in a NON-comment line (`return "error."` -> `"errorX."`) => same ruler returns
  non-empty diff, differing sha (`990b32975783...`), **rc=1**. Ruler is not tautological.
  (First attempt at the control was void: sed anchored on the stale pre-edit line number 164
  silently matched nothing; detected via mutation-evidence grep, redone pattern-anchored.)
- Reverse-write experiment (orchestrator's mandatory question): copy outside the repo, rewrite
  the comment to the lie "production UI and logs already embed this key" => ruler returns
  **rc=0, green, identical sha** (`ac2-reverse-write-experiment.log`). Restored on the spot;
  restore hash `ef90aa588a40c6e359b15c7f887755b8bfa10293186aca784dc684ea365488aa` equals the
  committed file. Conclusion for ruling: the ruler *should stay green* on a reversed comment --
  that is precisely its designed blindness to whole-line comments. It tests "production bytes
  unchanged", NOT "the comment was edited" and NOT "the comment is truthful". Comment truth
  content is adjudicated by AC#1 (non-implementer), and a *detecting* ruler for "comment lines
  changed at all" is the raw `git diff` in section 4 line-audit, not the zero-change ruler.

## 5. AC#3 behavior unchanged (both runs in out-of-repo export trees, never the shared worktree)
- Tables: `git archive 72b65c61 go.mod go.sum internal cmd tools scripts frontend models` +
  root `docs/` + manual `build/*.dll`, `third_party/sherpa-onnx/*.dll`, `frontend/dist`;
  T1 = anchor as-is (before), T2 = T1 + this leg's errors.go only (`diff -rq` pair shows that
  single file differing). Trees on `D:\W331R1` (8.3 short alias == full name, no spaces;
  verified via COM ShortPath). 750 `.go` files each (non-empty proof). `GOTMPDIR=D:\W331R1\tmp`
  created first. Command: `go test ./internal/observe/ -count=1 -timeout 30m -v`.
- Before (rc=0): `^=== RUN`=71 (incl. subtests), `^--- PASS`=71 (top-level only),
  `^--- FAIL`=0, `^panic`=0; log ends `ok ... 3.562s`.
- After (rc=0): 71 / 71 / 0 / 0; ends `ok ... 3.603s`. Red rosters both empty, diff rc=0;
  full run/pass/fail name-line lists byte-diff empty (rc=0).

## 6. AC#4 gates, all run in the shared working tree (not a partial export)
- `gofumpt.exe -l internal/observe/errors.go` (real binary /d/work/base/gopath/bin): rc=0, empty.
- `sh scripts/d22scan.sh`: rc=0, "clean - no D22 ban violations" (scope includes design/,
  frontend/, internal/, cmd/; my comment is pure ASCII).
- `sh scripts/check-path-length-budget.sh`: rc=0, "VERDICT GREEN" (needs .git; ran in shared tree).
- `go vet ./internal/observe/`: rc=0, empty.
- No red line containing walk / build failed / no test files / setup failed / 0xc0000135 seen.

## 7. Deliverable commit self-audit (rulers named)
- `git show --name-only e09187e1` => exactly one file: `internal/observe/errors.go`; every `+`/`-`
  diff line is a whole-line `//` comment (paste archived in `ac2-raw-diff.txt`).
- Write surface touched: errors.go comment lines, this evidence file,
  `.scratch/wisp/probes/331/r1/**`, one appended bullet on ticket 331. No AC box flipped, no
  Status/ticket-name change, no push, no amend/reset/rebase/stash/clean.
- Deliverable placeholder scan: the pattern and its command are written verbatim in
  `.scratch/wisp/probes/331/r1/logs/ac-selfcheck.log` (tokens kept out of this sentence on
  purpose, so the file cannot self-trigger its own scan); outcome 0 hits, grep rc=1.
