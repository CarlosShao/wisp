# 票 314 AC#0 — direct run, canonical launch shape.
# `go test ./cmd/wisp` runs the test binary with the PACKAGE directory as cwd
# (see cmd/wisp/panel_profile_leaf_shape_windows_test.go:569 filepath.Glob("*.go")); a
# direct run from the repo root is not the same test (run A, 121 bogus 0.00s reds).
$ErrorActionPreference = 'Continue'

$root = 'D:\work\workspace\projects plans\Wisp'
$dir  = 'D:\tmp\wisp\20261115\w314'
$exe  = Join-Path $dir 'wisp-winlive.test.exe'
$out  = Join-Path $dir 'stdout-cwd.log'
$err  = Join-Path $dir 'stderr-cwd.log'
$rcf  = Join-Path $dir 'rc-cwd.txt'
$pidf = Join-Path $dir 'run-cwd.pid'

$env:PATH = (Join-Path $root 'third_party\sherpa-onnx') + ';' +
            (Join-Path $root 'build') + ';' + $root + ';' + $env:PATH

foreach ($f in @($out, $err, $rcf, $pidf)) { if (Test-Path $f) { Remove-Item $f -Force } }

Set-Location (Join-Path $root 'cmd\wisp')
$cwd = (Get-Location).Path
$start = Get-Date

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

$u32 = ([int64]$rc) -band 0xFFFFFFFFL
$hex = '0x' + ([uint32]$u32).ToString('X8')
$seconds = [math]::Round(($end - $start).TotalSeconds, 1)

$summary = @(
  "EXE=$exe"
  "CWD=$cwd"
  "ARGS=-test.v -test.timeout=60m"
  "PID=$pid0"
  "START=$($start.ToString('yyyy-MM-dd HH:mm:ss'))"
  "END=$($end.ToString('yyyy-MM-dd HH:mm:ss'))"
  "ELAPSED_S=$seconds"
  "EXIT_CODE_DEC=$rc"
  "EXIT_CODE_HEX=$hex"
  "STDOUT_BYTES=$((Get-Item $out).Length)"
  "STDERR_BYTES=$((Get-Item $err).Length)"
)
$summary | ForEach-Object { $_ }
$summary | Set-Content -Path $rcf -Encoding utf8
