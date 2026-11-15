# 197-v2 adjudication leg — overlay mutants for the A1875 "would this test go red" half.
# Overlay ONLY: no repository source byte is modified. Each mutant is a full copy of
# one production file with one exact anchor replaced; the overlay JSON points the
# compiler at the copy. Anchors are asserted to hit exactly once before any run.
$ErrorActionPreference = 'Stop'
$repo = 'D:\work\workspace\projects plans\Wisp'
$out  = Join-Path $repo '.scratch\wisp\probes\197\v2'

function New-Mutant {
    param([string]$Name, [string]$Rel, [string[]]$Old, [string[]]$New)
    $src = Join-Path $repo $Rel
    $text = Get-Content -Raw -LiteralPath $src
    for ($i = 0; $i -lt $Old.Count; $i++) {
        $count = ([regex]::Matches($text, [regex]::Escape($Old[$i]))).Count
        if ($count -ne 1) { throw "ANCHOR $Name#$i hit $count times (want exactly 1) in $Rel" }
        $text = $text.Replace($Old[$i], $New[$i])
    }
    $dir = Join-Path $out $Name
    New-Item -ItemType Directory -Force -Path $dir | Out-Null
    $mutant = Join-Path $dir (Split-Path $Rel -Leaf)
    Set-Content -LiteralPath $mutant -Value $text -NoNewline
    $overlay = Join-Path $dir 'overlay.json'
    $json = @{ Replace = @{ $src = $mutant } } | ConvertTo-Json -Depth 4
    Set-Content -LiteralPath $overlay -Value $json -NoNewline
    return $overlay
}

function Run-Mutant {
    param([string]$Name, [string]$Pkg, [string]$Run, [int]$Count)
    $overlay = Get-ChildItem -LiteralPath (Join-Path $out $Name) -Filter 'overlay.json' | Select-Object -First 1
    $log = Join-Path $out $Name 'run.txt'
    Push-Location $repo
    try {
        $args = @('test', "-overlay=$($overlay.FullName)", $Pkg, '-run', $Run, '-count', "$Count", '-v')
        $outText = & go @args 2>&1 | Out-String
        $rc = $LASTEXITCODE
    } finally { Pop-Location }
    Set-Content -LiteralPath $log -Value $outText -NoNewline
    $fails = [regex]::Matches($outText, '(?m)^\s*--- FAIL: (\S+)') | ForEach-Object { $_.Groups[1].Value }
    $passes = [regex]::Matches($outText, '(?m)^\s*--- PASS: (\S+)') | ForEach-Object { $_.Groups[1].Value }
    [pscustomobject]@{
        Mutant = $Name; Pkg = $Pkg; Rc = $rc
        Red = ($fails -join ', '); Green = ($passes -join ', ')
    }
}

$M1old = @("`tr.recordLocked(taskID, TaskOutput{`n`t`tParentTaskID: parentTaskID, Label: label, Kind: TaskKindSubagent,`n`t`tState: statemachine.State(subagentStateRunning),`n`t})`n`tif r.cancel == nil {")
$M1new = @("`tr.recordLocked(taskID, TaskOutput{`n`t`tParentTaskID: parentTaskID, Label: label, Kind: TaskKindSubagent,`n`t`tState: `"`",`n`t})`n`tif r.cancel == nil {")
New-Mutant -Name 'M1-inflight-state-blank' -Rel 'internal\tools\task.go' -Old $M1old -New $M1new | Out-Null

$M2old = @("`tstate := subagentStateSettled`n`tswitch res.Status {`n`tcase agent.StatusCancelled:`n`t`tstate = subagentStateStopped`n`tcase agent.StatusFailed:`n`t`tstate = subagentStateFailed`n`t}`n`ttext := res.Text")
$M2new = @("`tstate := statemachine.State(`"`")`n`tswitch res.Status {`n`tcase agent.StatusCancelled:`n`t`tstate = subagentStateStopped`n`tcase agent.StatusFailed:`n`t`tstate = subagentStateFailed`n`t}`n`ttext := res.Text")
New-Mutant -Name 'M2-settled-state-blank' -Rel 'internal\tools\subagent_197.go' -Old $M2old -New $M2new | Out-Null

$M3old = @("func (s *StreamLog) Append(key, text string) {`n`tif s == nil || text == `"`" {`n`t`treturn`n`t}")
$M3new = @("func (s *StreamLog) Append(key, text string) {`n`tif s == nil || text == `"`" {`n`t`treturn`n`t}`n`tkey = `"subagent:merged`"")
New-Mutant -Name 'M3-stream-keys-merge' -Rel 'internal\panel\pump.go' -Old $M3old -New $M3new | Out-Null

New-Mutant -Name 'M4-parent-cancel-cascades' -Rel 'internal\tools\subagent_197.go' `
    -Old @('context.WithCancel(context.WithoutCancel(ctx))') -New @('context.WithCancel(ctx)') | Out-Null

New-Mutant -Name 'M5-watcher-bare-goroutine' -Rel 'internal\tools\subagent_197.go' `
    -Old @("`tobserve.Default.Spawn(`"agent-task-`"+bg.ID, `"tools`", watchRoot, func(c context.Context) {", "`t`tdone <- res`n`t})`n`n`tselect {") `
    -New @("`tgo func(c context.Context) {", "`t`tdone <- res`n`t}(context.WithoutCancel(ctx))`n`n`tselect {") | Out-Null

$results = @()
$results += Run-Mutant -Name 'M1-inflight-state-blank' -Pkg './internal/tools' -Count 1 `
    -Run 'Test197RowExistsBeforeFirstChildModelCall|Test197CancelIsPerRowAndNeverCascades|Test221ParentCancellationStillDoesNotCascade|Test197SpawnPublishesIdentityRow'
$results += Run-Mutant -Name 'M2-settled-state-blank' -Pkg './internal/tools' -Count 1 `
    -Run 'Test197SpawnPublishesIdentityRow|Test197RowExistsBeforeFirstChildModelCall'
$results += Run-Mutant -Name 'M3-stream-keys-merge' -Pkg './internal/panel' -Count 1 `
    -Run 'TestNPlusOneSubagentsEachGetTheirOwnStream|TestTwoSubagentsWritingAtOnceDoNotPolluteEachOther'
$results += Run-Mutant -Name 'M4-parent-cancel-cascades' -Pkg './internal/tools' -Count 3 `
    -Run 'Test197CancelIsPerRowAndNeverCascades|Test221ParentCancellationStillDoesNotCascade'
$results += Run-Mutant -Name 'M5-watcher-bare-goroutine' -Pkg './internal/tools' -Count 1 `
    -Run 'Test229FinishWatcherIsRosterNamedWhileChildRuns'

$results | Format-Table -AutoSize -Wrap | Out-String -Width 200
