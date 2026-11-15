# Receipt-leg readings for ticket 219 AC#6 (leg B of the 2026-11-15 dispatch).
# Read-only on production code: every mutation below is applied to a COPY and
# reaches the compiler only through `go test -overlay`; the frontend mutations
# run in .scratch/wisp/probes/219/mut/fe (a copy of frontend/{src,scripts} with
# node_modules junctioned to the real one), never in frontend/.
# Run from the repo root:  pwsh -File .scratch/wisp/probes/219/readings/00-run-all.ps1
$ErrorActionPreference = 'Continue'
$root = (Get-Location).Path
$mut = Join-Path $root '.scratch\wisp\probes\219\mut'
$fe  = Join-Path $mut 'fe'
$roster = Join-Path $root 'internal\panel\subagent_roster_197.go'

function Section($t) { Write-Output ""; Write-Output ("=" * 78); Write-Output "== $t"; Write-Output ("=" * 78) }

# One filtered, truncated view of a go test run: verdict lines, the exact red
# sentences, and any overlay/parse error - never the multi-KB packet dump.
function Show-GoTest($out) {
  $out | ForEach-Object {
    if ($_ -match '^(--- (PASS|FAIL)|FAIL|ok |PASS)' -or
        $_ -match 'subagent_(blocked_220|roster_197)_test\.go:\d+:' -or
        $_ -match 'overlay|error|cannot|invalid') {
      $_.Substring(0, [Math]::Min(300, $_.Length))
    }
  } | Select-Object -First 18
}

function Write-Overlay($name, $mutated) {
  $obj = @{ Replace = @{ $roster = $mutated } }
  $json = ConvertTo-Json -InputObject $obj -Depth 5
  Set-Content -LiteralPath (Join-Path $mut $name) -Value $json
  return (Join-Path $mut $name)
}

Section 'A. The two tables: declaration, write path, read door (grep, non-test)'
foreach ($f in @('internal\agent\approval\gate.go','internal\agent\approval\window_read.go','internal\agent\approval\queue.go','internal\agent\approval\pending_read.go')) {
  Write-Output "--- $f"
  Select-String -Path $f -Pattern 'windows  map|windows:  map|g\.windows\[|delete\(g\.windows|range g\.windows|func \(g \*Gate\) LiveL1Windows|pending \[\]\*qitem|len\(q\.pending\) >= q\.maxPend|q\.pending = append|range q\.pending|func \(q \*Queue\) LiveApprovals' |
    ForEach-Object { "{0}:{1}: {2}" -f $f, $_.LineNumber, $_.Line.Trim() }
}

Section 'B. Production call sites of each door (non-test .go under internal/ cmd/)'
$hits = Get-ChildItem -Recurse -Filter *.go -Path internal,cmd |
  Where-Object { $_.Name -notlike '*_test.go' } |
  Select-String -Pattern '\.LiveL1Windows\(\)|\.LiveApprovals\(\)'
$hits | ForEach-Object { "{0}:{1}: {2}" -f $_.Path.Substring($root.Length + 1), $_.LineNumber, $_.Line.Trim() }
Write-Output ("PRODUCTION CALLSITES: LiveL1Windows={0}  LiveApprovals={1}" -f `
  (@($hits | Where-Object { $_.Line -match 'LiveL1Windows' }).Count), `
  (@($hits | Where-Object { $_.Line -match 'LiveApprovals' }).Count))

Section 'C. The hop: pump -> roster join (L1 and L2 stay two sources into one boolean)'
foreach ($f in @('internal\panel\pump.go','internal\panel\subagent_roster_197.go','cmd\wisp\panel_pump.go','cmd\wisp\run.go')) {
  Write-Output "--- $f"
  Select-String -Path $f -Pattern 'L1WindowWait|L1Windows|Verdicts|p\.src\.Verdicts|l1Windows|indexWaitingKey|waiting\[row\.TaskID\]|LiveApprovals|taskRosterSectionFrom' |
    ForEach-Object { "{0}:{1}: {2}" -f $f, $_.LineNumber, $_.Line.Trim() }
}

Section 'D. Go baseline: the L1 join cells and the L1/L2 separation cells'
Push-Location $root
go test ./internal/panel -run 'TestAL1WindowInWaitingLightsItsRow|TestASecondCardForTheSameTaskStillLightsItsRow|TestTheClashFormIsRecognisedOnlyAsTheQueueWritesIt|TestABlockedRowKeepsItsRunningDimension|TestTheRosterReaderPutsSubagentsOnTheWire|TestTruncationFactsRideThePacket|TestDroppedStreamsAreNamedOnTheWire' -count=1 -v 2>&1 |
  Select-String -Pattern '^(=== RUN|--- (PASS|FAIL)|PASS|FAIL|ok)'
Write-Output ("[rc=$LASTEXITCODE]")
go test ./internal/agent/approval -run 'TestLiveL1WindowsSeesAWindowInTheLoop|TestReadingL1WindowsMovesNoGateState|TestL1WindowIsPlainDataWithNoRouteToAnAnswer|TestStartedCallsAreNotReportedAsWaiting' -count=1 -v 2>&1 |
  Select-String -Pattern '^(=== RUN|--- (PASS|FAIL)|PASS|FAIL|ok)'
Write-Output ("[rc=$LASTEXITCODE]")
Pop-Location

Section 'E. M1: L1 arm removed from the roster join (go test -overlay only)'
$p = Join-Path $mut 'm1-no-L1-arm\subagent_roster_197.go'
Copy-Item $roster $p -Force
$t = Get-Content -Raw -LiteralPath $p
$old = "`tfor _, w := range l1Windows {`n`t`t// Both halves of one window: an L1 window is keyed by its correlation, and`n`t`t// the row it blocks is named by its task.`n`t`tindexWaitingKey(waiting, known, w.CorrelationID)`n`t`tindexWaitingKey(waiting, known, w.TaskID)`n`t}`n"
if (-not $t.Contains($old)) { Write-Output 'MUT-ANCHOR NOT FOUND (m1)'; exit 1 }
Set-Content -LiteralPath $p -Value $t.Replace($old, "`t// MUTATION M1: the L1 arm is removed - the pre-220 shape.`n") -NoNewline
$ov = Write-Overlay 'm1-overlay.json' $p
$out = go test ./internal/panel -run 'TestAL1WindowInWaitingLightsItsRow|TestASecondCardForTheSameTaskStillLightsItsRow|TestTheClashFormIsRecognisedOnlyAsTheQueueWritesIt|TestABlockedRowKeepsItsRunningDimension' -count=1 -overlay $ov 2>&1
$rc = $LASTEXITCODE
Show-GoTest $out
Write-Output ("[M1 rc=$rc]")

Section 'F. M2: L2 card arm removed from the roster join (go test -overlay only)'
$p = Join-Path $mut 'm2-no-L2-arm\subagent_roster_197.go'
Copy-Item $roster $p -Force
$t = Get-Content -Raw -LiteralPath $p
$old = "`tfor _, card := range pending {`n`t`tindexWaitingKey(waiting, known, card.CorrelationID)`n`t}`n"
if (-not $t.Contains($old)) { Write-Output 'MUT-ANCHOR NOT FOUND (m2)'; exit 1 }
Set-Content -LiteralPath $p -Value $t.Replace($old, "`t// MUTATION M2: the L2 card arm is removed.`n") -NoNewline
$ov = Write-Overlay 'm2-overlay.json' $p
$out = go test ./internal/panel -run 'TestAL1WindowInWaitingLightsItsRow|TestASecondCardForTheSameTaskStillLightsItsRow|TestTheClashFormIsRecognisedOnlyAsTheQueueWritesIt|TestABlockedRowKeepsItsRunningDimension' -count=1 -overlay $ov 2>&1
$rc = $LASTEXITCODE
Show-GoTest $out
Write-Output ("[M2 rc=$rc]")

Section 'G. M3: the three numbers zeroed on the carrier (go test -overlay only)'
$p = Join-Path $mut 'm3-no-numbers\subagent_roster_197.go'
Copy-Item $roster $p -Force
$t = Get-Content -Raw -LiteralPath $p
$pairs = @{
  "`t`tInFlightSlots:     state.InFlightSlots," = "`t`tInFlightSlots:     0, // MUTATION M3"
  "`t`tPoolCap:           state.PoolCap," = "`t`tPoolCap:           0, // MUTATION M3"
  "`t`tStreamElidedRunes: state.StreamElidedRunes," = "`t`tStreamElidedRunes: 0, // MUTATION M3"
  "`t`t`tStreamTruncated:   row.StreamTruncated," = "`t`t`tStreamTruncated:   false, // MUTATION M3"
  "`t`t`tStreamElidedRunes: row.StreamElidedRunes," = "`t`t`tStreamElidedRunes: 0, // MUTATION M3"
}
foreach ($k in $pairs.Keys) { if (-not $t.Contains($k)) { Write-Output "MUT-ANCHOR NOT FOUND (m3): $k"; exit 1 }; $t = $t.Replace($k, $pairs[$k]) }
Set-Content -LiteralPath $p -Value $t -NoNewline
$ov = Write-Overlay 'm3-overlay.json' $p
$out = go test ./internal/panel -run 'TestTheRosterReaderPutsSubagentsOnTheWire|TestTruncationFactsRideThePacket|TestDroppedStreamsAreNamedOnTheWire' -count=1 -overlay $ov 2>&1
$rc = $LASTEXITCODE
Show-GoTest $out
Write-Output ("[M3 rc=$rc]")

Section 'H. Frontend: the real tree, unmodified (the gate this leg cites)'
Push-Location (Join-Path $root 'frontend')
npm run render:roster 2>&1 | Select-Object -Last 3
Write-Output ("[render:roster rc=$LASTEXITCODE]")
Pop-Location

Section 'I. Frontend mutations, run in the copy tree only (frontend/** bears zero writes)'
function Reset-Fe {
  Copy-Item (Join-Path $root 'frontend\src\App.tsx') (Join-Path $fe 'src\App.tsx') -Force
  Copy-Item (Join-Path $root 'frontend\src\components\harness\monitor-popover.tsx') (Join-Path $fe 'src\components\harness\monitor-popover.tsx') -Force
}
function Build-Run($tag) {
  Push-Location $fe
  node node_modules/vite/bin/vite.js build --ssr scripts/render-task-roster.tsx --outDir (".tmp/" + $tag) --logLevel warn 2>&1 | Select-String -Pattern 'error during' | Select-Object -First 2
  node (".tmp/" + $tag + "/render-task-roster.js") 2>&1 | Select-Object -Last 10
  Write-Output ("[$tag rc=$LASTEXITCODE]")
  Pop-Location
}
Reset-Fe
Write-Output '--- I0 copy baseline (must be OK)'
Build-Run 'i0'

Write-Output '--- I1 popover: pool span + row suffix + section line removed'
$p = Join-Path $fe 'src\components\harness\monitor-popover.tsx'
$t = Get-Content -Raw -LiteralPath $p
$t = $t.Replace('const boundLine = sectionBoundLine(streamTruncated ?? false, streamElidedRunes, streamDroppedCount);', 'const boundLine = null as string | null; // MUTATION I1')
$t = $t.Replace('const bound = rowBoundLabel(agent);', 'const bound = null as string | null; // MUTATION I1')
$t = $t.Replace('{poolInFlight !== undefined && poolCap !== undefined ? (', '{false ? ( // MUTATION I1')
Set-Content -LiteralPath $p -Value $t -NoNewline
Build-Run 'i1'
Reset-Fe

Write-Output '--- I2 App: the five feed lines + the derivation stream triple removed'
$p = Join-Path $fe 'src\App.tsx'
$t = Get-Content -Raw -LiteralPath $p
$map = @{
 '    streamTruncated: row.streamTruncated,' = '    // MUT-I2'
 '    streamElidedRunes: row.streamElidedRunes,' = ''
 '    streamDropped: row.streamDropped,' = ''
 '                  poolInFlight={snapshot.tasks?.inFlightSlots}' = '                  /* MUT-I2 */'
 '                  poolCap={snapshot.tasks?.poolCap}' = ''
 '                  streamTruncated={snapshot.tasks?.streamTruncated}' = '                  /* MUT-I2 */'
 '                  streamElidedRunes={snapshot.tasks?.streamElidedRunes}' = ''
 '                  streamDroppedCount={snapshot.tasks?.droppedCount}' = ''
}
foreach ($k in $map.Keys) { if (-not $t.Contains($k)) { Write-Output "MUT-ANCHOR NOT FOUND (i2): $k"; exit 1 }; $t = $t.Replace($k, $map[$k]) }
Set-Content -LiteralPath $p -Value $t -NoNewline
Build-Run 'i2'
Reset-Fe

Write-Output '--- I3 popover: the 等批准 tier removed'
$p = Join-Path $fe 'src\components\harness\monitor-popover.tsx'
$t = Get-Content -Raw -LiteralPath $p
$t = $t.Replace('  if (row.blockedOnApproval) return { label: "等批准", dot: "blocked" };', '  // MUTATION I3')
Set-Content -LiteralPath $p -Value $t -NoNewline
Build-Run 'i3'
Reset-Fe

Write-Output '--- I4 App: the blockedOnApproval derivation removed'
$p = Join-Path $fe 'src\App.tsx'
$t = Get-Content -Raw -LiteralPath $p
$t = $t.Replace('    blockedOnApproval: row.blockedOnApproval,', '    // MUT-I4')
Set-Content -LiteralPath $p -Value $t -NoNewline
Build-Run 'i4'
Reset-Fe

Section 'J. Pristine check: zero writes to production code by this leg'
Push-Location $root
git status --porcelain=v1 -- frontend/src internal cmd
Write-Output '(empty above = zero writes)'
Pop-Location
