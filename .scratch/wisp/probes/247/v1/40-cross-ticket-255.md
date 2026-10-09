# 247 v1 — question 2: does the ticket-255 roster edit loosen that ticket's judgement?

Commit under review: `ada5c563` (`cmd/wisp/config_readers_255.go`, +24/−7).
This leg answers three sub-questions independently; all readings are mine.

**Verdict: 不 loosening. The 255 guard did not go blind, and no dumb key was washed into
"explained" in any machine-checked sense.** One wording change is real but strengthens the
instrument rather than weakening it; details and rulers below. One residual honest gap is named at
the end (the tier choice, not the tier's strength).

## (a) Does `hotClaimNoReader` → `hotClaimOtherProcess` for `audio` launder the key?

No, because the only thing the tiers *do* is decide whether a section may be named in the
「已立即生效」 sentence, and that gate is keyed on a **different** prefix:

`cmd/wisp/config_readers_255.go:289-300`, verbatim head of the decision:

```go
		consumed := true
		for _, row := range rows {
			verdict, ok := hotRowClaims[row]
			if !ok || !strings.HasPrefix(verdict, hotClaimConsumed) {
				consumed = false
			}
		}
		if consumed {
			s.claimable = append(s.claimable, name)
```

`hotClaimConsumed` is `"consumed: "` (`config_readers_255.go:79`); `hotClaimOtherProcess` is
`"other-process: "` (`:82`). So `[audio]` lands in `hotSectionSplit.quiet` exactly as it did before,
and the receipt still may not claim it took effect immediately. The claim-strength change is zero.

Precedent shape check (the same file, rows this leg read verbatim):

- `:118` `"models": hotClaimOtherProcess + "cmd/wisp/models.go:163 [cfg.Models.Mirror] - only the
  wisp models subcommands read it, in their own process, off a fresh LoadFile",`
- `:129` `"hotkey": hotClaimOtherProcess + "cmd/wisp/resident_ball_windows.go:276 [Hotkeys:  cfg,] - …
  but the ball lives in the resident process, not in this `wisp run` host",`
- new `:150` `"audio":   hotClaimOtherProcess + "cmd/wisp/resident_audio_windows.go:196
  [c.Audio.MicMutedDefault] - the resident capture leg is [audio]'s first production reader and it
  reads the file once at boot, so a hot reload of this section waits for a restart instead of acting
  in this run",`

Same tier, same reason (a reader that lives in the *resident* process, re-read per boot, not a live
re-reader), so the row is not a new class of claim. "First production reader" is also true by this
leg's own negative ruler — everything outside `internal/config` and outside tests that touches the
section:

```
grep -rn "\.Audio\." --include=*.go cmd/ internal/ | grep -v "_test.go" | grep -v "^internal/config/"
  => cmd/wisp/config_readers_255.go:149 (comment), :150 (this roster row),
     cmd/wisp/resident_audio_windows.go:32 (comment), :193, :196, :213   (the leg itself)
```

Nothing was deleted by the commit (the `−7` are the seven rewritten lines, and this leg diffed it:
`"Ball"/"Session"/"Privacy"/"Memory"/"Cost"` rows untouched, `sectionReadSites` keys all still
present, only `Audio`/`Voice` values went from `nil` to a one-element list).

## (b) Are the new `file:line` anchors verbatim correct?

Yes — this leg pulled the lines itself rather than trusting the row text:

```
sed -n '196p' cmd/wisp/resident_audio_windows.go
	ra.mutedAtBoot = c.Audio.MicMutedDefault
sed -n '197p' cmd/wisp/resident_audio_windows.go
	ra.voiceEnabled = c.Voice.Enabled
```

`196` carries the cited token `c.Audio.MicMutedDefault`; `197` carries `c.Voice.Enabled`. And those
anchors are not decoration: `cmd/wisp/config_receipt_255_test.go:197-211` re-reads every cited file
from disk (`raw, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))`), bounds-checks
the line number (`:203`) and fails when the line no longer carries the token (`:208`
`			if !strings.Contains(got, token) {`). So a future move of that read site reddens 255 — teeth
intact, and now pointed at **five more** lines than before (`+4` voice rows that gained a cite plus
the audio row).

The four `[voice.*]` rows: tier constant unchanged (`hotClaimNoReader` still the prefix at
`config_readers_255.go:199-202`), the `扫描零命中` marker still present, and only the *subject* of
the sentence narrowed from "reads cfg.Voice" (section-wide) to "reads this [voice] key". The narrow
claim is true per this leg's own ruler:

```
grep -rlE "TTS\.Speed|\.Punctuation|WakeWord\.Thresholds|WakeWord\.VetoWords" --include=*.go cmd internal \
  | grep -v "_test.go" | grep -v "^internal/config/" | wc -l
0
```

Had the wording stayed section-wide it would have become **false** the moment `[voice] gained a
production reader — so the edit is a correction, not a softening.

## (c) Did the 255 guard go blind? Final state, run by this leg

```
PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp/ -count=1 -run 'TestTicket255' -v
rc_255_run=0
--- PASS: TestTicket255SplitOnlyClaimsSectionsWithALiveReader (0.00s)
--- PASS: TestTicket255HotRowRosterCoversTheRegistry (0.00s)
--- PASS: TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim (0.00s)
--- PASS: TestTicket255RosterStillMatchesTheActualReadSites (0.46s)
--- PASS: TestTicket255ReceiptOmitsPanelFromTheImmediateSentence (2.07s)
--- PASS: TestTicket255ReceiptOmitsASectionWithNoReaderAnywhere (2.02s)
--- PASS: TestTicket255ReceiptStillNamesTheLiveReadSection (2.00s)
--- PASS: TestTicket255ReceiptSentenceAssemblyIsFiltered (1.02s)
   (+12 more PASS in the family; 20 top-level, 0 FAIL)
ok  	github.com/CarlosShao/wisp/cmd/wisp	7.764s
```

Why the guard is *not* blind, stated as a mechanism this leg checked rather than as a hope:

1. `TestTicket255RosterStillMatchesTheActualReadSites` compares the live scan against
   `sectionReadSites` by **set equality** (`config_receipt_255_test.go:244`
   `		if !slicesEqualSorted(got, want) {`). A passing equality with `want =
   {"cmd/wisp/resident_audio_windows.go"}` proves the scan really finds that reader — had the batch
   not registered it, the same test reddens (that is precisely the 12:4x red the implementer
   reports; the direction of the failure is symmetric, so it cannot be gamed by leaving the roster
   stale *or* by padding it).
2. The reader-less sample that 255 uses to prove the receipt can say "no reader" is `[session]`, not
   `[audio]` — verbatim at `config_receipt_255_test.go:462-466` ("[session] takes the sample over:
   three keys … still zero production readers in every process") and its assertion plants
   `[session]` (`:471`). And `[llm]` remains the live-reader negative control (`:504`). So the two
   load-bearing samples of 255's judgement are untouched by this batch. Ruler for "no test quietly
   became about audio": `grep -rn "section=audio" --include=*.go cmd/wisp/` => 0 hits.
3. `citedRowsFloor = 8` (`:220`) is a floor on checked cites, so more cites cannot satisfy it worse.

## Residual (for the orchestrator, not a failure)

- Whether `[audio]` should be `other-process` or `snapshot-only` is a naming argument I do not
  settle: neither prefix reaches `claimable`, and the row's own text says "reads the file once at
  boot, so a hot reload waits for a restart" — which is `snapshot-only`'s meaning too. `models`/
  `hotkey` precedent picks `other-process` for exactly the "resident process, fresh LoadFile" shape,
  so the batch is consistent; if the orchestrator wants the *timing* fact on the tier rather than in
  prose, that is a 255-side decision, and per `AGENTS.md §0.2` it is human-approved, so this leg
  leaves it alone.
- The genuine new fact this edit records honestly: `mic_muted_default` is read **once at boot**, so
  changing it needs a restart. That is the same shape ticket 255 already labels `hotkey`, and it is
  why ticket 247's AC#2 missing shot ("put the mic in a state where the device actually opens") also
  requires a restart or a hot-tier reader that ticket 228/180 owns.
