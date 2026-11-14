# 314-deshare2 A/B for the winlive 35v2 transport case and the cold/hot budget
# pair, alternating builds so machine-load drift hits both. Same cmd.exe reason as
# ab-exit2s.ps1: PowerShell 7.6 mangles "-test.run=...".
param(
  [int]$Rounds = 2,
  [string]$BinDir = 'D:\tmp\wisp\20261114\314-deshare2\bin'
)
$ErrorActionPreference = 'Stop'
$root = (Resolve-Path (Join-Path $PSScriptRoot '..\..\..\..')).Path
$logDir = Join-Path $PSScriptRoot 'logs'
New-Item -ItemType Directory -Force -Path $logDir | Out-Null
$results = Join-Path $PSScriptRoot 'results-winlive-ab.jsonl'
$env:PATH = (Join-Path $root 'third_party\sherpa-onnx') + ';' + $env:PATH

$cases = @(
  @{ name = 'transport'; regex = '^TestLive35v2PagePostMessageEnvelopeReachesTheGoDoor$'; pattern = '35v2 webview processes after teardown: .*|35v2 left .* tree processes behind' },
  @{ name = 'budget'; regex = '^(TestPanelHostRealWindowHopAndLifecycle_WinLive|TestPanelHostLatencyPercentilesAC2_WinLive)$'; pattern = 'winlive hot re-show measured.*|sample  1: .*|cold P95 .* exceeds.*' }
)
foreach ($round in 1..$Rounds) {
  foreach ($case in $cases) {
    foreach ($tag in @('deshare', 'before')) {
      $exe = Join-Path $BinDir "wisp-$tag.test.exe"
      $stem = "$($case.name)-$tag-r$round"
      $out = Join-Path $logDir "$stem.out.txt"
      $sw = [Diagnostics.Stopwatch]::StartNew()
      $cmd = '"{0}" "-test.run={1}" -test.v -test.timeout=5m > "{2}" 2>&1' -f $exe, $case.regex, $out
      & $env:ComSpec /c $cmd
      $rc = $LASTEXITCODE
      $sw.Stop()
      $text = Get-Content $out -Raw
      $verdicts = ([regex]::Matches($text, '(?m)^--- (PASS|FAIL|SKIP): \S+') | ForEach-Object { $_.Value }) -join ' | '
      $reading = ([regex]::Matches($text, "(?m)$($case.pattern)") | ForEach-Object { $_.Value }) -join ' | '
      $line = [ordered]@{
        round = $round
        case = $case.name
        build = $tag
        rc = $rc
        elapsed_s = [math]::Round($sw.Elapsed.TotalSeconds, 2)
        verdicts = $verdicts
        reading = $reading
      }
      $json = ($line | ConvertTo-Json -Compress)
      Add-Content -Path $results -Value $json -Encoding utf8
      $json
    }
  }
}
