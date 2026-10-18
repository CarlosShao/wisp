# 339-r1 AC#0/AC#4 台件③b：Windows 本机档（PowerShell 原生，⛔ Git Bash ⇒ MSYS 改写 TMP）
# 用法: pwsh -File win.ps1 <run-regex> <label>
param([string]$RunPat = '.', [string]$Label = 'win')
$ErrorActionPreference = 'Continue'
$REPO = 'D:\work\workspace\projects plans\Wisp'
$TREE = 'D:\tmp\wisp\20261018\339r1\win\tree'
$GT = 'D:\tmp\wisp\20261018\339r1\gt'
$LOG = Join-Path $REPO ".scratch\wisp\probes\339\r1\logs\$Label.log"
New-Item -ItemType Directory -Force -Path $GT | Out-Null
New-Item -ItemType Directory -Force -Path (Split-Path $LOG) | Out-Null
$env:TMPDIR = $GT; $env:TEMP = $GT; $env:TMP = $GT; $env:GOTMPDIR = $GT
$env:GOFLAGS = '-mod=mod'; $env:GOPROXY = 'off'; $env:CGO_ENABLED = '0'
$env:GOCACHE = 'D:\tmp\wisp\20261018\339r1\gocache'
Set-Location $TREE
"label=$Label tree=$PWD run=$RunPat"
go test ./internal/risk/ -run $RunPat -count=1 -timeout 30m -v *> $LOG
$RC = $LASTEXITCODE
"GO_TEST_RC=$RC"
"size=$((Get-Item $LOG).Length)"
$lines = Get-Content $LOG -Encoding utf8
"RUN=$(($lines | Where-Object { $_ -match '^=== RUN' }).Count)"
"PASS_TOP=$(($lines | Where-Object { $_ -match '^--- PASS' }).Count)"
"FAIL_TOP=$(($lines | Where-Object { $_ -match '^--- FAIL' }).Count)"
"SKIP=$(($lines | Where-Object { $_ -match '^--- SKIP' }).Count)"
'=== red roster ==='
$lines | Where-Object { $_ -match '^--- FAIL' }
'=== skip roster ==='
$lines | Where-Object { $_ -match '^--- SKIP' }
'=== fault words ==='
($lines | Select-String -Pattern 'build failed|no test files|setup failed|creating work dir|0xc0000135' | ForEach-Object { $_.Line }) | Select-Object -First 5
