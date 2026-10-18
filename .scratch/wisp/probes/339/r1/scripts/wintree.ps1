# 339-r1 交件侧 Windows 尺：指定树整包（原生 PowerShell，⛔ Git Bash 通道改写 TMP）。
param([string]$Tree, [string]$Label)
$ErrorActionPreference = 'Continue'
$REPO = 'D:\work\workspace\projects plans\Wisp'
$LOG = Join-Path $REPO ".scratch\wisp\probes\339\r1\logs\$Label.log"
$env:TEMP = 'D:\tmp\wisp\20261018\339r1\gt'
$env:TMP = $env:TEMP
$env:GOTMPDIR = $env:TEMP
$env:GOFLAGS = '-mod=mod'; $env:GOPROXY = 'off'; $env:CGO_ENABLED = '0'
$env:GOCACHE = 'D:\tmp\wisp\20261018\339r1\gocache'
New-Item -ItemType Directory -Force -Path $env:TEMP | Out-Null
Set-Location $Tree
go test ./internal/risk/ -count=1 -timeout 30m -v *> $LOG
"GO_TEST_RC=$LASTEXITCODE"
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
$lines | Select-String -Pattern 'build failed|no test files|setup failed|creating work dir|0xc0000135' | ForEach-Object { $_.Line } | Select-Object -First 5
