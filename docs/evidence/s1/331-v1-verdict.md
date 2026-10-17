# 331-v1 verdict (non-implementer) — ticket 331, MessageKey comment in internal/observe/errors.go

- Leg `331-v1`, adjudicator, read-only on `internal/observe/**`. Material: `docs/evidence/s1/331-r1-evidence.md`
  (114 lines / 8082 bytes at `d4a1da8b`, re-measured here) + `.scratch/wisp/probes/331/r1/**` + commits
  `e09187e1 5cbec52d adcbe8f3 7ab51937 47f7a6a3 d4a1da8b`.
- Every number below was re-run by this leg, off named commit trees (`72b65c61` start, `e09187e1`
  landed), never off in-flight working bytes. The shared worktree file hash `ef90aa58…` equals the
  `e09187e1` blob, so gates run in the worktree do judge the bytes under review.
- My artifacts: `.scratch/wisp/probes/331/v1/**` — 台件 `build-trees.sh`, `ac2-replay.sh`, `ac2-blindness.sh`;
  读数 `logs/ac0-rerun`, `ac34-rerun`, `ac2-replay`, `ac2-blindness`, `ac2-scope-boundaries`, `ac4-gates`,
  两发 `go test` 日志, 两枚名册 + `roster-diff.txt`, plus `impl-leg-zero-change-check.copy.sh` (the ruler under judgment, archived verbatim so every reading here can be replayed).

## Verdicts

AC#0 成立 · AC#1 成立 · AC#2 成立（钉的是产码字节那一维，票面与尺射程相符）· AC#3 成立 · AC#4 成立 · AC#5 成立。三条残余在文末，其中一枚归票 225、不记本票。

## AC#0 — startup scales really re-run, ticket face checked cell by cell: 成立

| scale (all on `72b65c61`) | my reading | ticket face | leg |
|---|---|---|---|
| `git grep -n MessageKey -- internal cmd tools ':(exclude)*_test.go'` rc=0 | 2 lines = `:160` comment + `:163` decl, embedders **0** | same | same |
| same ruler without exclude rc=0 | 8 lines, real call sites 4 (`errors_test.go:191,192`, `matrix_14_2_test.go:175,249`) | same 4 | same 4 |
| `sed -n '160,162p'` rc=0 | three lines verbatim incl. "UI and logs embed the key" | same | same |
| family control `grep -nE 'strings.(Contains\|ContainsAny\|HasPrefix\|HasSuffix\|EqualFold\|Index)' -- internal/risk ':(exclude)*_test.go'` | **37** | 37 | **not re-run** |

The fourth row is the leg's only gap: it never re-ran the `internal/risk` 37 positive control. I ran
it, it matches the ticket, so the ticket face is not stale and nothing is unreadable. Extra: the ticket's
older anchor `0e944d93` is an ancestor of `72b65c61` and gives the same comment lines and the same
2-line production reading, so the leg's "行号未漂" claim holds at both anchors.

## AC#1 — the rewrite tells the truth: 成立

Before `72b65c61:errors.go:160-162` / after `e09187e1:errors.go:160-166`, verbatim:

    -// MessageKey returns the stable user-visible copy key for the class. The copy
    -// text itself is deferred (D23 error-copy work, S6); UI and logs embed the key
    -// so copy can be added later without re-touching error sites.
    +// MessageKey returns the stable user-visible copy key reserved for the
    +// deferred error-copy layer (D23 error-copy work, S6; registered DEFERRED in
    +// docs/specs/SPEC-12-roadmap-governance.md section 5). As read on 2026-10-17
    +// (ticket 331) no production code embeds this key: its only callers are the
    +// tests internal/observe/errors_test.go (TestMessageKeys) and
    +// internal/llm/matrix_14_2_test.go. Whether UI and logs embed the key is
    +// that deferred layer's decision, not a promise made here.

Asked per sentence, "is this supported by the tree now":
1. "returns the stable user-visible copy key" — body is `return "error." + string(c)` (`:168`),
   stability pinned by a real test. Supported.
2. "registered DEFERRED in docs/specs/SPEC-12-roadmap-governance.md section 5" — file exists, section
   5 is 「推迟项登记表」, line 77 verbatim `| DEFERRED | 完整错误文案体系 | D23 | … | S6 | … |`.
   Supported; not a phantom citation — `TestMessageKeys` at `errors_test.go:188`, both named test
   files on disk, every referenced artifact swept present (`logs/ac0-rerun.log` tail).
3. "As read on 2026-10-17 … no production code embeds this key" — re-run at `e09187e1` and at today's
   HEAD; broadened to `grep -niE 'message_?key'` over `frontend internal cmd tools` (non-test), which
   leaves only the declaration. Supported.
4. "its only callers are the tests …" — exactly those four call sites. Supported.
5. "not a promise made here" — no forward promise. Supported.

Direction: it states the missing wiring as fact and hands the decision off; the promise sentence is
gone; nothing is laundered into "no problem"; no sibling "UI already uses the key" phrasing appears.
Dropping "so copy can be added later without re-touching error sites" is right — with zero consumers
there is no error site to re-touch, so that clause described a no-op.

## AC#2 — mandatory answer: what the zero-change ruler pins: 成立

Ruler = drop every `^[[:space:]]*//` line on both sides, then `sha256sum` + `diff`.
Main scale reproduced: both sides `f51f944c8aaf1d7a…c408efd5539`, diff empty, rc=0. Positive control
reproduced: non-comment literal `error.` -> `errorX.` -> `990b32975783…`, rc=1. Not tautological.
Blind band mapped shape by shape (mutations only on out-of-repo copies, `logs/ac2-blindness.log` +
`logs/ac2-scope-boundaries.log`):

| green under the ruler | red under the ruler |
|---|---|
| B1 anchor vs itself — a leg that changed nothing at all | A2 non-comment literal changed (`errorX.`) |
| B2 the whole seven-line honest comment block deleted | B4 `+ ""` appended on a code line |
| B3 a U+2713 planted inside a whole-line comment | C1 trailing comment added on a code line |
| C3 a whole-line comment far away edited (`errors.go:82` RetryNever) | C2 one blank line added |
| A3b truth file vs reverse-written lie file (bytes differ, comment opposite) | |
| A4 comment rewritten toward a different true sentence | |

Its exact reach in one sentence: "no byte of errors.go that is not a whole-line comment may move."
C1/C2 show the filter drops only lines that are wholly comments, so trailing comments and blank lines
stay in reach; the blind band is precisely whole-line comment text. Three consequences, which are the
answer: (a) a no-op leg takes the main scale green (B1); (b) blindness is file-wide, so green does not
even prove the edits stayed inside the MessageKey block (C3); (c) green cannot separate truth from
lie (A3b). What pins "the comment really was rewritten" is a different ruler — the raw `git diff` line
audit: changed 10 (3 removed, 7 added), non-whole-comment changed 0 (re-run as A5), single file,
single hunk `@@ -157,9 +157,13 @@`.

Ruling: 成立, pinned dimension = production bytes unchanged outside whole-line comments. Ticket face
and ruler reach DO match: the cell is titled 「AC#2（产码字节零改动那把尺）」, its first clause reads
「滤掉整行注释后」 (declares the filter out loud), its mandatory positive control is placed on a
NON-comment line, and its closing 「其中非整行注释枚数必为 0」 is the same dimension. So AC#2 is not
「不成立/不可判」, and I have no sentence to name as over-claiming — no criteria-text change is
needed. The real hazard is a reader citing AC#2's rc=0 as proof of honesty: that is an AC#1 claim,
carriable only by a non-implementer (this file). Criteria wording belongs to the orchestrator and
adjudicators either way; a writing leg must not touch it either.

## AC#3 — behavior unchanged: 成立

Trees built by me (`build-trees.sh`): archive of the eight-piece set + root `docs/` + manual `build/*.dll`,
`third_party/sherpa-onnx/*.dll`, `frontend/dist`; both trees 750 `.go` (non-empty proof), `diff -rq`
shows exactly one differing file; trees under `D:\V331V1` where ShortPath equals the full path (8.3
verified by script); `GOTMPDIR` created first. `go test ./internal/observe/ -count=1 -timeout 30m -v`:
before rc=0, RUN=71 PASS_TOP=71 FAIL_TOP=0 SKIP=0 panic=0 (`ok … 3.566s`, stderr empty); after identical
71/71/0/0 (`ok … 3.561s`). Red rosters both empty; sorted name lists 142 lines each, `roster-diff.txt`
empty (rc=0). Sanity on the count shape: `grep -c '^func Test'` = 71 and indented subtest lines = 0, so
`=== RUN` (with subtests) legitimately equals `--- PASS` here — not a mixed-scale artifact. `-v` only.

## AC#4 — four gates: 成立

Shared worktree only. gofumpt (`/d/work/base/gopath/bin/gofumpt.exe -l internal/observe/errors.go`)
rc=0 empty; `sh scripts/d22scan.sh` rc=0 "clean - no D22 ban violations" with `ban #8 design/ examined
39 text files` in scope (no missing-`design/` exit-2 GUARD case); `sh
scripts/check-path-length-budget.sh` rc=0 "VERDICT GREEN" (14579 tracked paths, 60 over-budget, 60
rostered, 0 unrostered); `go vet ./internal/observe/` rc=0 empty. Red-word sweep on my own log matched
7 `walk` lines — all d22scan self-test case names or PASS lines
(`TestDeclaredEmojiScopeCannotWalkZeroFiles`, `TestLedgerCountsMatchAnIndependentWalk`,
`TestWalksSkipGitIgnoredPaths`) or an expected-stderr sample quoted inside a passing selftest.
`grep -cE '^( *)--- FAIL|^FAIL'` on that log = 0. No table failure, nothing excluded.

## AC#5 — discipline: 成立

Per-commit `git show --name-only`: `e09187e1` = one file, `internal/observe/errors.go` (10 lines, all
whole-line comments); `5cbec52d` = 27 files all under `.scratch/wisp/probes/331/r1/`; `adcbe8f3` /
`7ab51937` / `d4a1da8b` = that probe dir (commit-msg drafts, one log) + the single evidence file;
`47f7a6a3` = ticket 331 (numstat 2 added / 0 deleted: blank + one bullet) + a probe draft. Four kinds
of surface, nothing else; evidence filename is short and carries no section symbol. Ticket at HEAD: six
AC boxes still `- [ ]`, zero `- [x]`, `Status:` still 待派, no `-done`; no leg commit is an ancestor of
`origin/dev` (nothing pushed); evidence corrections were appended as three separate commits
(88 -> 90 -> 114 lines), history never rewritten.

## Three residuals (for the orchestrator; box-flipping and criteria wording stay with you)

- R1 (this ticket, wording): line 1 still reads "the stable user-visible copy key" while the body says
  nothing embeds it. Honest as written — "user-visible" names the key's designated role — but it only
  becomes literally true once the copy layer exists. Not disqualifying for AC#1.
- R2 (this ticket, perishability): "As read on 2026-10-17" is a dated snapshot with no standing guard. A
  later wiring ticket does turn the AC#2-style ruler red for that file (A2/B4/C1 show it has teeth on
  real bytes), but no test pins the date clause. Whoever opens the wiring ticket owns this sentence.
- R3 (归票 225, not this ticket's bill): production `DEFERRED(` count is 26 at `72b65c61`, `e09187e1`
  and today's HEAD, and `errors.go` still holds 0 markers — that registration gap is 225's audit. New
  detail for 225: this rewrite adds a bare-word `DEFERRED` hit in `errors.go` (production bare-word
  count 46 -> 47). 225's ruler needs the parenthesized form, so its 26 stands, but a future looser grep
  on bare `DEFERRED` would pick this comment up.

## Self-audit (three rulers, not "done")

1. My commits, each with `git show --name-only`, are listed in my ticket bullet; production-file bytes
   touched by this leg = zero (`internal/observe/**` read-only for me).
2. This file: 160 lines or fewer, non-empty. The fill-in-token ruler and its readings are spelled out in
   `.scratch/wisp/probes/331/v1/logs/ac-selfcheck.log` (tokens kept out of this sentence so the file
   cannot self-trigger its own scan): strict 0 hits rc=1, angle-bracket form 0 hits rc=1.
3. Decisive readings quoted verbatim above and in `logs/`: stripped-string match `f51f944c…` rc=0;
   control `990b3297…` rc=1; A3b truth-vs-lie same stripped string with the two files proven
   byte-different; B1 no-op green; C3 file-wide blindness; A2/B4/C1/C2 red; four counts 71/71/0/0
   twice; gates rc=0 x4; restore check `git status --porcelain -- internal/observe/errors.go` empty
   with the blob hash `ef90aa58…` unchanged.
Mandatory answer: AC#2 成立 on the production-byte dimension; ticket wording matches the ruler's reach;
no criteria-text change requested, and none may be tasked to a writing leg.
