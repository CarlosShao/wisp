# 票 314 AC#0 — canonical direct run: the binary is driven EXACTLY the way
# scripts/wisp-cli-staged-tests.sh:492 drives it, only the wrapper is gone:
#   (cd cmd/wisp && <binary> -test.count=1 -test.timeout=90m -test.v)
# plus the carrier's staging: the pinned sherpa/onnxruntime DLLs sit NEXT TO the
# binary (System32 otherwise shadows onnxruntime -- stage-sherpa-dlls.sh:8-21),
# and the five voice fixtures + WISP_43_REPO_ROOT are exported.
# Run A (repo-root cwd, no staging) is kept only as the "not the same test"
# control: 121 reds, 62 of them purely cwd, ~40 purely environment.
param([string]$Suffix = 'c')
$ErrorActionPreference = 'Continue'

$root = 'D:\work\workspace\projects plans\Wisp'
$dir  = 'D:\tmp\wisp\20261115\w314'
$exe  = Join-Path $dir 'wisp-winlive.test.exe'
$out  = Join-Path $dir "stdout-$Suffix.log"
$err  = Join-Path $dir "stderr-$Suffix.log"
$rcf  = Join-Path $dir "rc-$Suffix.txt"
$pidf = Join-Path $dir "run-$Suffix.pid"

$env:PATH = (Join-Path $root 'third_party\sherpa-onnx') + ';' +
            (Join-Path $root 'build') + ';' + $root + ';' + $env:PATH

$env:WISP_SPIKE_MODELS = Join-Path $root 'third_party\spike-models'
$env:WISP_SILERO_MODEL = Join-Path $root 'third_party\spike-models\silero_vad.onnx'
$env:WISP_TTS_DIR      = Join-Path $root 'third_party\spike-models\matcha-icefall-zh-baker'
$env:WISP_KWS_DIR      = Join-Path $root 'third_party\spike-models\sherpa-onnx-kws-zipformer-wenetspeech-3.3M-2024-01-01'
$env:WISP_PUNC_MODEL   = Join-Path $root 'third_party\model-fetch\x-punc\sherpa-onnx-punct-ct-transformer-zh-en-vocab272727-2024-04-12-int8\model.int8.onnx'
$env:WISP_43_REPO_ROOT = $root

foreach ($f in @($out, $err, $rcf, $pidf)) { if (Test-Path $f) { Remove-Item $f -Force } }

Set-Location (Join-Path $root 'cmd\wisp')
$cwd = (Get-Location).Path
$start = Get-Date

$p = Start-Process -FilePath $exe `
     -ArgumentList '-test.count=1 -test.timeout=90m -test.v' `
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
  "ARGS=-test.count=1 -test.timeout=90m -test.v"
  "PID=$pid0"
  "START=$($start.ToString('yyyy-MM-dd HH:mm:ss'))"
  "END=$($end.ToString('yyyy-MM-dd HH:mm:ss'))"
  "ELAPSED_S=$seconds"
  "EXIT_CODE_DEC=$rc"
  "EXIT_CODE_HEX=$hex"
  "EXIT_CODE_U32=$u32"
  "STDOUT_BYTES=$((Get-Item $out).Length)"
  "STDERR_BYTES=$((Get-Item $err).Length)"
)
$summary | ForEach-Object { $_ }
$summary | Set-Content -Path $rcf -Encoding utf8
