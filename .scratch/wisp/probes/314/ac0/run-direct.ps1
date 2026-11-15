# 票 314 AC#0 — direct run of the compiled cmd/wisp test binary (winlive tag).
# No wrapper, no `go test`; capture the true process exit code + full stderr.
# Usage: pwsh -NoProfile -File run-direct.ps1
$ErrorActionPreference = 'Continue'

$root = 'D:\work\workspace\projects plans\Wisp'
$dir  = 'D:\tmp\wisp\20261115\w314'
$exe  = Join-Path $dir 'wisp-winlive.test.exe'
$out  = Join-Path $dir 'stdout.log'
$err  = Join-Path $dir 'stderr.log'
$rcf  = Join-Path $dir 'rc.txt'
$pidf = Join-Path $dir 'run.pid'

# DLLs: missing them yields 0xC0000135 (environment, not ticket 314).
$env:PATH = (Join-Path $root 'third_party\sherpa-onnx') + ';' +
            (Join-Path $root 'build') + ';' + $root + ';' + $env:PATH

foreach ($f in @($out, $err, $rcf, $pidf)) { if (Test-Path $f) { Remove-Item $f -Force } }

$start = Get-Date
$cpuBefore = (Get-Process msedgewebview2 -ErrorAction SilentlyContinue | Measure-Object).Count

$p = Start-Process -FilePath $exe `
     -ArgumentList '-test.v -test.timeout=60m' `
     -RedirectStandardOutput $out -RedirectStandardError $err `
     -PassThru -NoNewWindow
Set-Content -Path $pidf -Value $p.Id -Encoding ascii
$pid0 = $p.Id

Wait-Process -Id $pid0
$p.Refresh()
$rc  = $p.ExitCode
$end = Get-Date

# signed int32 -> unsigned 32-bit hex, and the two readings PowerShell gives
$u32 = ([int64]$rc) -band 0xFFFFFFFFL
$hex = '0x' + ([uint32]$u32).ToString('X8')
$dec = $rc
$seconds = [math]::Round(($end - $start).TotalSeconds, 1)

$cpuAfter = (Get-Process msedgewebview2 -ErrorAction SilentlyContinue | Measure-Object).Count

$summary = @(
  "EXE=$exe"
  "ARGS=-test.v -test.timeout=60m"
  "PID=$pid0"
  "START=$($start.ToString('yyyy-MM-dd HH:mm:ss'))"
  "END=$($end.ToString('yyyy-MM-dd HH:mm:ss'))"
  "ELAPSED_S=$seconds"
  "EXIT_CODE_DEC=$dec"
  "EXIT_CODE_HEX=$hex"
  "EXIT_CODE_U32=$u32"
  "STDOUT_BYTES=$((Get-Item $out).Length)"
  "STDERR_BYTES=$((Get-Item $err).Length)"
  "WEBVIEW2_BEFORE=$cpuBefore"
  "WEBVIEW2_AFTER=$cpuAfter"
)
$summary | ForEach-Object { $_ }
$summary | Set-Content -Path $rcf -Encoding utf8
