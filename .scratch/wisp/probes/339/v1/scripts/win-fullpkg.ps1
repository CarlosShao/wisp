# 339-v1 台件⑦：Windows 本机档整包（★PowerShell 原生，⛔ Git Bash ⇒ MSYS 改写 TMP 形）。
# 用法: pwsh -File win-fullpkg.ps1 <tree-win-path> <label>
# 台面＝D: 卷仓外导出树；TMP/GOTMPDIR＝Windows 反斜杠形、先 New-Item；四数逐枚点名⛔ 混用⛔ 相减。
param([string]$TreeWin = 'D:\tmp\wisp\20261018\339v1\winpost\tree', [string]$Label = 'v1-win-fullpkg', [string]$RunPat = '.')
$ErrorActionPreference = 'Continue'
$REPO = 'D:\work\workspace\projects plans\Wisp'
$GT = 'D:\tmp\wisp\20261018\339v1\gt-win'
$LOGDIR = Join-Path $REPO '.scratch\wisp\probes\339\v1\logs'
$LOG = Join-Path $LOGDIR "$Label.log"
New-Item -ItemType Directory -Force -Path $GT | Out-Null
New-Item -ItemType Directory -Force -Path $LOGDIR | Out-Null
$env:TMPDIR = $GT; $env:TEMP = $GT; $env:TMP = $GT; $env:GOTMPDIR = $GT
$env:GOFLAGS = '-mod=mod'; $env:GOPROXY = 'off'; $env:CGO_ENABLED = '0'
Set-Location $TreeWin
"label=$Label tree=$PWD run=$RunPat"
go test ./internal/risk/ -run $RunPat -count=1 -timeout 30m -v *> $LOG
"GO_TEST_RC=$LASTEXITCODE"
"bytes=$((Get-Item $LOG).Length)"
$lines = Get-Content $LOG -Encoding utf8
"RUN=$((($lines | Where-Object { $_ -match '^=== RUN' }).Count))"
"PASS_TOP=$((($lines | Where-Object { $_ -match '^--- PASS' }).Count))"
"FAIL_TOP=$((($lines | Where-Object { $_ -match '^--- FAIL' }).Count))"
"SKIP=$((($lines | Where-Object { $_ -match '^--- SKIP' }).Count))"
'=== red roster ==='
$lines | Where-Object { $_ -match '^--- FAIL' }
'=== skip roster ==='
$lines | Where-Object { $_ -match '^--- SKIP' }
'=== 台面故障词 ==='
($lines | Select-String -Pattern 'build failed|no test files|setup failed|creating work dir|0xc0000135|mklink' | ForEach-Object { $_.Line }) | Select-Object -First 8
