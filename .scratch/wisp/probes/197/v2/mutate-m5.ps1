# M5 re-run: same mutant as mutate.ps1's M5, with the one compile fix the first attempt
# needed (watchRoot must stay used). Overlay only; no repository byte is modified.
$ErrorActionPreference = 'Stop'
$repo = 'D:\work\workspace\projects plans\Wisp'
$out  = Join-Path $repo '.scratch\wisp\probes\197\v2\M5-watcher-bare-goroutine'
$src  = Join-Path $repo 'internal\tools\subagent_197.go'
$text = Get-Content -Raw -LiteralPath $src
$old1 = "`tobserve.Default.Spawn(`"agent-task-`"+bg.ID, `"tools`", watchRoot, func(c context.Context) {"
$old2 = "`t`tdone <- res`n`t})`n`n`tselect {"
foreach ($o in @($old1, $old2)) {
    $n = ([regex]::Matches($text, [regex]::Escape($o))).Count
    if ($n -ne 1) { throw "ANCHOR hit $n times (want 1): $o" }
}
$text = $text.Replace($old1, "`t_ = watchRoot`n`tgo func(c context.Context) {")
$text = $text.Replace($old2, "`t`tdone <- res`n`t}(context.WithoutCancel(ctx))`n`n`tselect {")
New-Item -ItemType Directory -Force -Path $out | Out-Null
$mutant = Join-Path $out 'subagent_197.go'
Set-Content -LiteralPath $mutant -Value $text -NoNewline
@{ Replace = @{ $src = $mutant } } | ConvertTo-Json -Depth 4 | Set-Content -LiteralPath (Join-Path $out 'overlay.json') -NoNewline
Push-Location $repo
try { $t = & go test "-overlay=$(Join-Path $out 'overlay.json')" ./internal/tools -run 'Test229FinishWatcherIsRosterNamedWhileChildRuns' -count 1 -v 2>&1 | Out-String }
finally { Pop-Location }
Set-Content -LiteralPath (Join-Path $out 'run.txt') -Value $t -NoNewline
$t
