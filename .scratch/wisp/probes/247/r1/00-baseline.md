# 247 r1 — evidence leg (write-code leg)

Anchor at start: `b02fda30` on `dev`, shared worktree. Clock of this leg's start: see the
Progress-log line appended to the ticket (interpolated from `date` stdout).

## 0. Pre-change baselines (measured by this leg, before any of its own edits)

Ruler 1 — AC#1 non-test importers of `internal/audio` (ticket 现量 1, verbatim command):

```
grep -rl "CarlosShao/wisp/internal/audio" --include=*.go . | grep -v "/internal/audio/" | grep -vc _test
```

Reading: `0`

Ruler 2 — AC#2 `SetAudioLevel` non-test producers (`grep -rn "SetAudioLevel" --include=*.go internal/ cmd/`):

Production-code hits (non-test):

- `internal/ball/liquid_windows.go:42` — the seam itself (`func (b *Ball) SetAudioLevel(level float32)`)
- `cmd/balldebug/main.go:418` — synthetic envelope feed (flag text says `synthetic audio envelope`)
- `cmd/balldebug/main.go:477` — feeds `0`

=> real-microphone producer count before this leg: `0`. Ticket 现量 2 confirmed, not stale.

## 1. Gates (AC#8). One line per command with `rc=N`, appended as each is run.

(to be filled by this leg; each gate writes its own section below, never an empty file)
