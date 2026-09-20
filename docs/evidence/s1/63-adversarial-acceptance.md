# Ticket 63 — adversarial acceptance (credential-entry CLI)

**Acceptor:** orchestrator (not the implementer). **Date:** 2026-09-20.
**Commits under review:** `722f819`, `cb1bf21`, `4477f56`, `6ab2a3d`.
**Files:** `cmd/wisp/secret.go`, `cmd/wisp/secret_test.go`, `cmd/wisp/secret_argv_windows_test.go`, `internal/secret/*`.

A dispatched acceptance agent was killed by a platform connection error after 26 tool
uses, mid-experiment. It had already staged its own mutation (M2 below) in
`cmd/wisp/secret.go`. I preserved that patch outside the repo
(`ticket63-mutation-M2.patch`), ran its experiment to completion, then reverted the
production file and confirmed the tree clean. This report is the orchestrator's own
verdict, with the implementer's claims treated as unverified until re-derived here.

## Verdict table (1:1 with the ticket's AC boxes, per README rule 6)

| AC | Verdict | Evidence I produced myself |
|---|---|---|
| 1 — hidden input, two-entry confirm, DPAPI blob, masked `get` | **PASS** | `TestSecretSetGetListRoundTrip` (`secret_test.go:173`) asserts the blob bytes contain no plaintext (`:190`) and that `get` returns only `RedactSecret`'s masked form (`:203`); `TestSecretSetConfirmationMismatchStoresNothing` (`:331`) asserts a mismatch stores zero blobs. Re-run: `go test -count=2 ./cmd/wisp/` → `ok 11.317s`. |
| 2 — no plaintext in argv / logs / errors / temp files, proven not asserted | **PASS** | (a) **My own mutation**: added `"value", value` to the audit `slog.Info` at `secret.go:326` → **5 assertions went red** (`set logged the plaintext secret`, `--show must not log the plaintext`, `FailurePathsLogAndPrintNoPlaintext`, `UnsetRefusesWhileReferenced`, `EndToEnd…`), then reverted. (b) The argv claim is **not vacuous**: the positive control (`secret_argv_windows_test.go:452`) plants `fakeKey` as a real positional token in a live child's argv and demands the **same** `pollCommandLine()` reader used by the negative test (`:253` vs `:463`) report it. (c) The temp-file half has its own planted-decoy control (`secret_test.go:556`) and asserts the data dir holds **exactly one** file after `set` (`:572`) — an extra scratch file fails on name, not just on content. |
| 3 — `list` shows id+time only; `unset` refuses while referenced; `--force` leaves an audit line | **PASS with MINOR-1** | `TestSecretUnsetRefusesWhileReferenced` names the referencing field and asserts `forced=true` appears in the forced delete's own record (`:822`). The narrowing that `4477f56` introduced is exploitable — see MINOR-1. |
| 4 — dev/test/prod are three independent blobs | **PASS** | `TestSecretSameNameUnderThreeEnvsIsThreeBlobs` (`:886`) is content-level, not path-level: it stores a **different** key per env, asserts 3 distinct ciphertexts (`len(seen) != 3` → Fatal, `:939`), reads each env back and asserts it equals its own key **and** equals no other env's (`:944-958`), and asserts the prod listing never reaches into the dev dir (`:970`). |
| 5 — portable and non-portable both covered via ticket 06's seam | **PASS** | `TestSecretPortableModeUsesTicket06Seam` (`:977`) proves the marker flips `layout.Portable` both ways, the portable data dir is exe-relative, the round-trip still works, and the portable blob holds no plaintext (`:1024`). The non-portable branch is exercised by every other probe (`portable=false`) with the same no-plaintext-byte assertion, so a silent plaintext downgrade fails on content, not on a flag. |
| 6 — provider resolves `api_key_ref` and issues the request | **TRANSFERRED, not failed** | The test proves ref → `config.LoadFile` → `ProviderKeys` → `Authorization` header, but the request is sent by the test's own http client, so the provider segment is caller-self-certified — which README rule 6 forbids as proof. Handed to **ticket 12** with a mutation-based completion criterion; registered as **A8** in `docs/reports/pending-and-issues.md`. The implementer left this box unticked and I confirm that was the correct call. |
| 7 — adversarial acceptance by a non-implementer, 1:1 table | **PASS** | This document. |

## Defects

**[MINOR-1] `auditRecordFor` narrowing opened a hole the code it replaced did not have.**
The old assertion grepped the whole log buffer for `forced=true`; `4477f56` scoped it to the
single record containing the needle. That fixed a real cross-subtest pollution bug, but it
lost a detection the whole-buffer grep had. Reproduction — the dead agent's own mutation,
which logs `forced=true` into a **separate** record that does not contain the needle, on
**both** branches (forced and non-forced):

```
$ git apply ticket63-mutation-M2.patch   # two added slog.Info lines at secret.go:471/473
$ go test -count=1 -run 'TestSecretUnset' ./cmd/wisp/
ok  	github.com/CarlosShao/wisp/cmd/wisp	0.064s        # FALSE GREEN
```

The non-forced delete still passes while the process logs a `forced=true` record that no
assertion accounts for. Under the pre-`4477f56` assertion this failed.

Severity MINOR, not MAJOR: what is weakened is an **audit-trail** assertion, not a leak
assertion — the plaintext-detection path is independently proven to bite (AC#2(a) above).

Prescribed fix (do not re-scope; count instead): keep `auditRecordFor` for the positive
half, and add a whole-buffer invariant that the number of records containing `forced=true`
**equals the number of forced deletes the test actually performed**. That closes M2 without
reopening the pollution bug, because the count is asserted per test rather than per buffer.
The mutation above must then go red.

**[NOTE, not a defect] The temp-file claim's bound is honestly narrow.** It covers the
injected data dir (fully walked, and asserted to contain only the blob) plus files **newly
created** in the OS temp dir (name-set diff before/after). An in-place append to a file that
already existed in `%TEMP%` would not be flagged. Given the store's only write path is
`blobPath` under the injected data dir, I am not opening a ticket for this; recording it so
the claim is never quoted more broadly than "no new file anywhere carries the value".

## Forbidden-pattern sweep (my own run over ticket 63's files)

`go func(` / `filepath.Clean` / `filepath.Abs` / `time.Now().Sub` timeouts / emoji —
zero hits across `cmd/wisp/secret.go`, both test files, and `internal/secret/*.go`.
`gofmt -l` empty. `go vet ./cmd/wisp/ ./internal/secret/` RC 0.

## Gates I re-ran myself (not the implementer's numbers)

```
gofmt -l cmd/wisp internal/secret                         → empty
go vet ./cmd/wisp/ ./internal/secret/                     → RC 0
go test -count=2 ./internal/secret/                       → ok 0.195s
go test -count=2 ./cmd/wisp/                              → ok 11.317s
go test -race -count=1 ./cmd/wisp/ ./internal/secret/     → ok 7.875s / 1.225s
```

`cmd/wisp`'s test binary links sherpa-onnx: `third_party/sherpa-onnx` must be on `PATH` or
child processes die with `0xc0000135`. Pre-existing property of the package, not introduced
here, and not documented in BUILD.md — worth a line there.

## Disposition

**Clean to close, with MINOR-1 registered for a follow-up.** AC#6 is a legitimate transfer
rather than a shortfall of this ticket. No BLOCKER, no MAJOR, one MINOR with a prescribed
fix and a reproduction that currently goes green.
