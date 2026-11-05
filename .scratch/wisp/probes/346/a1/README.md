# 346 AC#3 - CI red roster attribution ruler

Ticket: `.scratch/wisp/issues/346-two-race-window-instruments-red-on-abstention-under-ci-contention-blob-identical-flip.md`, AC#3.

The ticket asks for a re-runnable bench piece that answers one question per red name:
**which ticket owns this?** It is not a CI step (`.github/**` needs human approval) and not a
standing assertion. It writes only inside a temp dir.

## Files

| file | what it is |
|---|---|
| `ci-red-attribution.sh` | the ruler: run ids (via `gh`) or pre-fetched `--log` files in, per-name attribution out |
| `selftest.sh` | offline, fixture-driven selftest (34 checks); no `gh`, no network, no repo pool |
| `fixtures/gh-prefixed.log` | gh-shaped log: `job<TAB>step<TAB>timestamp --- FAIL: Name` |
| `fixtures/all-green.log` | a log with no top-level FAIL row (must be REFUSED, not reported as 0 red) |
| `fixtures/pool/` | mini issue pool: criteria owner, file owner, prose-only name |
| `fixtures/ledger.md` | mini ledger with one `A####` reading |
| `fixtures/funcs.tsv` | test seam for the file fallback (`Name<TAB>path`); default is `git grep` |
| `fixtures/bin/gh` | fixture `gh` that always fails, for the exit-3 check |
| `.gitattributes` | pins this subtree to LF, so a Windows checkout cannot turn the fixture log into CRLF (the ruler anchors `^---`, and `bin/gh` is a shebang script) |

## Usage

From the repository root:

```bash
# two CI run ids (gh must be logged in)
bash .scratch/wisp/probes/346/a1/ci-red-attribution.sh 39124122121 39101060753

# already-fetched logs (offline)
bash .scratch/wisp/probes/346/a1/ci-red-attribution.sh \
  --log /d/tmp/wisp/20261105/346r1/run-39124122121.log \
  --log /d/tmp/wisp/20261105/346r1/run-39101060753.log

# selftest
bash .scratch/wisp/probes/346/a1/selftest.sh
```

Exit codes: `0` ok / `1` positive control / `2` usage / `3` gh fetch / `4` unreadable log /
`5` empty roster (refused on purpose).

## Why it is shaped this way

1. **Prefix first, anchor second.** gh prefixes every line with `job<TAB>step<TAB>timestamp `,
   so the ticket face's `^--- FAIL` ruler reads 0 hits on a real red log. The ruler strips the
   prefix, then anchors `^---`; indented subtests and mid-line prose stay out of the denominator.
2. **A zero-row roster is refused (exit 5).** A ruler that measured nothing may not hand back
   `unique red names=0` as if that were a green.
3. **Positive control before attribution (exit 1).** A name that is known to be registered must
   resolve; otherwise the ruler is blind and says so instead of attributing.
4. **File fallback.** Ticket 346 names its instruments by **file**, not by Go func name, so the
   resolver falls back from `func <Name>` to the file basename and then to the tickets naming it.
   `fixtures/funcs.tsv` is the hermetic seam for that step; the default is `git grep`.
5. **Prose alone is not an owner.** `散文提及` is a reading, not a registration; only a criteria
   line, a named file, or a ledger id makes a row `REGISTERED`.

## Known limits (do not read this as a green)

- It answers **"is this red name owned by a ticket?"**, not **"is this red new in this leg?"**.
  The latter needs the two-run diff plus the blob-identity ruler (`git rev-parse <rev>:<path>`).
- `--log` only sees the failed-step log that `gh run view --log-failed` returns, so step counts
  are a lower bound, not the full job roster.
- It reads the issue pool as text. If a ticket names a red only in prose, the row is
  `UNREGISTERED` by design; add a criteria line or a ledger id if the ticket really owns it.

Provenance: the baseline ruler was delivered inline in `docs/evidence/s1/346-c1-census.md`
section 3; this directory is that piece, hardened with the guards above and a selftest that
fails on the baseline (see `docs/evidence/s1/346-r1-impl.md`).
