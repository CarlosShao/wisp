# 314-deshare2 A/B: the winlive 2s exit clause on the de-share build vs the
# pre-change (fixed-name) build, alternating so machine load drifts hit both.
# Runs the compiled test binaries through cmd.exe because PowerShell 7.6's
# native argument passing mangles "-test.run=..." (observed: the child got
# "-test" alone and exited rc=2).
param(
  [int]$Rounds = 3,
  [string]$BinDir = 'D:\tmp\wisp\20261114\314-deshare2\bin'
)
$ErrorActionPreference = 'Stop'
$root = (Resolve-Path (Join-Path $PSScriptRoot '..\..\..\..')).Path
$logDir = Join-Path $PSScriptRoot 'logs'
New-Item -ItemType Directory -Force -Path $logDir | Out-Null
$results = Join-Path $PSScriptRoot 'results-exit2s.jsonl'
$env:PATH = (Join-Path $root 'third_party\sherpa-onnx') + ';' + $env:PATH

foreach ($round in 1..$Rounds) {
  foreach ($tag in @('deshare', 'before')) {
    $exe = Join-Path $BinDir "wisp-$tag.test.exe"
    $stem = "exit2s-$tag-r$round"
    $out = Join-Path $logDir "$stem.out.txt"
    $sw = [Diagnostics.Stopwatch]::StartNew()
    $cmd = '"{0}" -test.run=^TestPanelHostWebViewChildrenExitWithinTwoSeconds_WinLive$ -test.v -test.timeout=5m > "{1}" 2>&1' -f $exe, $out
    & $env:ComSpec /c $cmd
    $rc = $LASTEXITCODE
    $sw.Stop()
    $text = Get-Content $out -Raw
    $reading = [regex]::Match($text, 'winlive exit clause: .*')
    $verdict = [regex]::Match($text, '(?m)^--- (PASS|FAIL): \S+')
    $line = [ordered]@{
      round = $round
      build = $tag
      rc = $rc
      elapsed_s = [math]::Round($sw.Elapsed.TotalSeconds, 2)
      verdict = $verdict.Value
      reading = $reading.Value
    }
    $json = ($line | ConvertTo-Json -Compress)
    Add-Content -Path $results -Value $json -Encoding utf8
    $json
  }
}
