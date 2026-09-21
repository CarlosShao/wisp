# 79 — `artifactName` folds distinct tool-call ids onto one disk name, and `writeFileExclusive` isn't exclusive

**Status:** ready-for-agent
**Type:** correctness/data-integrity defect (artifact clobbering) + a storage-hygiene gap
**Blocks:** nothing · **Blocked by:** nothing (`internal/agent`/`internal/memory` are free once ticket 76 landed)
**Packages:** `internal/agent/spill.go`, `internal/memory/artifacts.go` + tests. Do **not** touch
`internal/risk` (ticket 75 in flight), `internal/ball`, `cmd/wisp`, `internal/proc` (ticket 78 in flight).
**Found by:** ticket 76's agent, reported and **deliberately not fixed** (correct call — out of that
ticket's scope). Spec refs: C25 provenance/taint, D31 atomic-rename, D6 single source of truth.

## Defect 1 (the real one): two different artifacts can silently overwrite each other

Measured by ticket 76 and re-read by me at `internal/agent/spill.go:132`:
- `artifactName(callID, seq)` keeps only `[A-Za-z0-9_-]`, so **`p/q`, `p\q` and `pq` all collapse to
  `tool-output-pq.txt`**. Those are *different* model-supplied tool-call ids.
- The doc comment admits only "same call id retried"; the sanitizer's collision is wider than that.
- Worse, the function named **`writeFileExclusive` does not use `O_EXCL`** — it's `os.WriteFile`
  (`O_CREATE|O_TRUNC`), so the name is a lie and the second writer silently wins.
- Consequence: the earlier `Spill.Path` (handed to the model / usable via `fs.read`) can point at
  bytes that were replaced by a different call's output. **That is a provenance break, not cosmetics**:
  a tool result the model believes it can re-read may be another tool's output.

### What to build
- Make the on-disk name **injective** with respect to the logical id (hash or escape rather than strip —
  stripping is what merges them), or detect the collision and fail loudly. Say which you chose and why.
- Make `writeFileExclusive` actually exclusive (`O_CREATE|O_EXCL`), or rename it to what it does.
  ⚠ If you make it truly exclusive, decide and **test** what happens on retry of the same id
  (the current documented behavior is that retry overwrites) — don't let that become an error path
  that breaks the caller.
- ⚠ **Do not "fix" the collision by allowing caller-controlled paths through** (that's the invariant
  ticket 76 just pinned: no caller-controlled destination path; see its AC#1/AC#2 tests).

## Defect 2: nobody enforces that `artifacts/` is flat

`listArtifactsDir` does `if e.IsDir() { continue }`, so a stray subdirectory is **invisible** to
`ListArtifacts`, `PurgeArtifacts`, and the 500 MB LRU quota. Ticket 76's fixture proved it: a `nested\`
canary directory **survived a purge** (existing behavior, recorded honestly, not caused by them).

### What to build
Pick one and write the reason in the commit message: (a) the quota/purge walks recursively so a stray
dir can't hide bytes from the cap; or (b) the store treats a non-conforming entry as an error it reports
rather than silently ignoring. **"Silently ignore" is the status quo and it is the option that must not win
by default** — an entry invisible to the quota is how a 500 MB cap becomes 2 GB.

## AC (1:1 verdict table, one row per box)
- [ ] **AC#1** A test proving two *different* ids that previously collided now produce *different* on-disk
  names (or a loud failure), with the collision pair from the report (`p/q` vs `pq`) as named subtests.
- [ ] **AC#2** `writeFileExclusive` is genuinely exclusive (or renamed + behavior documented), **plus** a test
  pinning what a same-id retry does now (both outcomes are acceptable; *unspecified* is not).
- [ ] **AC#3** Mutation: revert the injectivity change ⇒ AC#1 must go red; **grep-prove the mutation landed
  before running and prove the revert after** (this repo has produced several false greens from
  no-op mutations today — anchor on the line that carries the thing you're changing, not on a symbol name).
- [ ] **AC#4** Defect 2 resolved one of the two ways above, with a test where a stray subdirectory either
  counts against the quota or is reported — and that test must go red under the current `continue`.
- [ ] **AC#5** Gates: `gofmt -l` on touched pkgs empty, `go vet` scoped rc=0, `go test -count=2` scoped, and
  the RUN-count invariant (`=== RUN` == 2 × distinct names) with **the 2 known pre-existing `SKIP`s in
  `internal/memory/concurrent_test.go:186` named explicitly** rather than filtered out of the report.

## Rules
`git commit -q -F - -- <explicit paths> <<'MSGEOF'` (quoted heredoc, A31); no `git add -A`; check
`git diff --cached --name-only`; no `--amend`/`reset`/`rebase`/`stash` (A34); no whole-repo gates
(shared tree); no worktree inside the repo (A38④); don't touch ticket 76's or ticket 20's faces
(those are mine). First checkpoint commit within 15 tool calls; sync Status + boxes + `next=` every commit.
