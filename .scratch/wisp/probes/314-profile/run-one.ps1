# 314 profile A/B harness - ONE test process at a time, always.
# usage: pwsh -File run-one.ps1 -Bin <path to .test.exe> -Tag <A1|A2|B> -Run 1
param(
  [Parameter(Mandatory=$true)][string]$Bin,
  [Parameter(Mandatory=$true)][string]$Tag,
  [Parameter(Mandatory=$true)][int]$Run
)
$ErrorActionPreference = 'Stop'
$root = (Resolve-Path (Join-Path $PSScriptRoot '..\..\..\..')).Path
$probe = $PSScriptRoot
$logDir = Join-Path $probe 'logs'
$tagRun = "$Tag-r$Run"
$out = Join-Path $logDir "$tagRun.out.txt"
$err = Join-Path $logDir "$tagRun.err.txt"
if (Test-Path $out) { Remove-Item $out -Force }
if (Test-Path $err) { Remove-Item $err -Force }
$re = '^(TestAC14GoSideEvalPushReachesThePage|TestAC13BringUpSurvivesAReusedThreadQuit|TestPanelThreadIsSTAAndExitsCleanly|TestAC13BringUpRefusesAThreadWithAQueuedClose)$'
$fixedLeafs = @('wisp-33r5-panel-profile','wisp-33r6-reused-thread-quit','wisp-33r7-queued-close')

function Snap-Profiles {
  $res = @()
  foreach ($leaf in $fixedLeafs) {
    $p = Join-Path $env:TEMP $leaf
    if (Test-Path $p) {
      $f = Get-ChildItem $p -Recurse -File -ErrorAction SilentlyContinue
      $res += [pscustomobject]@{ leaf = $leaf; exists = $true; files = $f.Count; bytes = [int64](($f | Measure-Object Length -Sum).Sum); lastWrite = (Get-Item $p).LastWriteTime.ToString('yyyy-MM-dd HH:mm:ss') }
    } else {
      $res += [pscustomobject]@{ leaf = $leaf; exists = $false; files = 0; bytes = 0; lastWrite = '' }
    }
    # per-run unique leaves (B shape): prefix-<pid>-<nanotime>
    $u = @(Get-ChildItem $env:TEMP -Directory -Filter "$leaf-*" -ErrorAction SilentlyContinue)
    if ($u.Count -gt 0) {
      $res += [pscustomobject]@{ leaf = "$leaf-*"; exists = $true; files = $u.Count; bytes = 0; lastWrite = ($u | Sort-Object LastWriteTime | Select-Object -Last 1).LastWriteTime.ToString('yyyy-MM-dd HH:mm:ss') }
    }
  }
  return $res
}

function Snap-Browsers {
  $all = @(Get-Process msedgewebview2 -ErrorAction SilentlyContinue)
  $profs = @()
  foreach ($p in (Get-CimInstance Win32_Process -Filter "Name='msedgewebview2.exe'" -ErrorAction SilentlyContinue)) {
    $m = [regex]::Match($p.CommandLine, '--user-data-dir="?([^" ]+)"?')
    if ($m.Success -and ($m.Groups[1].Value -match 'wisp-33r[567]-')) { $profs += "$($p.ProcessId)=$($m.Groups[1].Value)" }
  }
  return [pscustomobject]@{ total = $all.Count; onTestProfiles = $profs }
}

$compBefore = @(Get-Process wisp.test -ErrorAction SilentlyContinue).Count
$profBefore = Snap-Profiles
$brBefore = Snap-Browsers

$sw = [Diagnostics.Stopwatch]::StartNew()
$p = Start-Process -FilePath $Bin -ArgumentList @('-test.run', $re, '-test.v', '-test.timeout', '10m') `
  -WorkingDirectory $root -RedirectStandardOutput $out -RedirectStandardError $err -PassThru -NoNewWindow
$hung = -not $p.WaitForExit(300000)
if ($hung) { Stop-Process -Id $p.Id -Force -ErrorAction SilentlyContinue; Start-Sleep -Seconds 3 }
$sw.Stop()
$rc = if ($hung) { 'KILLED-300s' } else { $p.ExitCode }

$profAfter = Snap-Profiles
$brAfter = Snap-Browsers
$text = if (Test-Path $out) { Get-Content $out -Raw } else { '' }
$errText = if (Test-Path $err) { Get-Content $err -Raw } else { '' }
$verdicts = @([regex]::Matches($text, '(?m)^\s*--- (PASS|FAIL|SKIP): (\S+)') | ForEach-Object { "$($_.Groups[1].Value):$($_.Groups[2].Value)" })
$rc32 = $null
if ($rc -is [int]) { $rc32 = '0x{0:X8}' -f ([uint32]($rc -band 0xFFFFFFFF)) }

$row = [ordered]@{
  tag = $Tag; run = $Run; binary = (Split-Path $Bin -Leaf)
  rc = $rc; rcHex = $rc32
  elapsed_s = [math]::Round($sw.Elapsed.TotalSeconds, 2)
  verdicts = $verdicts
  runLines = ([regex]::Matches($text, '(?m)^=== RUN')).Count
  passLines = ([regex]::Matches($text, '(?m)^--- PASS:')).Count
  failLines = ([regex]::Matches($text, '(?m)^--- FAIL:')).Count
  skipLines = ([regex]::Matches($text, '(?m)^--- SKIP:')).Count
  panic = ([regex]::Matches($text + $errText, '(?m)^panic:')).Count
  bootFailed = ([regex]::Matches($text + $errText, 'webview2: .*boot failed')).Count
  embedFailed = ([regex]::Matches($text + $errText, 'WebView2 Embed failed')).Count
  createWindow = ([regex]::Matches($text + $errText, 'MARK299 after CreateWindowExW')).Count
  markBefore = ([regex]::Matches($text + $errText, 'MARK299 before Embed')).Count
  markAfter = ([regex]::Matches($text + $errText, 'MARK299 after Embed')).Count
  accessViolation = ([regex]::Matches($text + $errText, '0xc0000005|access violation|exception')).Count
  errBytes = $errText.Length
  competitorsBefore = $compBefore
  browsersBefore = $brBefore.total
  browsersOnTestProfilesBefore = ($brBefore.onTestProfiles -join ';')
  browsersAfter = $brAfter.total
  browsersOnTestProfilesAfter = ($brAfter.onTestProfiles -join ';')
  profilesBefore = ($profBefore | ForEach-Object { "{0}:exists={1},files={2},bytes={3},lw={4}" -f $_.leaf,$_.exists,$_.files,$_.bytes,$_.lastWrite }) -join ' | '
  profilesAfter = ($profAfter | ForEach-Object { "{0}:exists={1},files={2},bytes={3},lw={4}" -f $_.leaf,$_.exists,$_.files,$_.bytes,$_.lastWrite }) -join ' | '
  outFile = (Split-Path $out -Leaf)
}
$line = ($row | ConvertTo-Json -Compress -Depth 6)
Add-Content -Path (Join-Path $probe 'results.jsonl') -Value $line -Encoding utf8
$line
