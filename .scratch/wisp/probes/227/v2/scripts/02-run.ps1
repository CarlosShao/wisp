# AC#5 整包终态那一发（票 227 票面命令）。
# 三条台规：1) 整条 PATH 一律 Windows 形＋分号（pwsh 原生＝无 POSIX 冒号形混入）；
# 2) GOTMPDIR＝短、C 卷（与 %LOCALAPPDATA% 同卷＝票 311 认可形），先建再指；
# 3) 退码＝go 命令裸跑重定向到文件后立刻取 $LASTEXITCODE（无管道、无 tail）。
param([string]$Gotmp = 'C:\Users\swq\AppData\Local\wisp-test-gotmp\227v2-a',
       [string]$Log = 'D:\tmp\wisp\20261019\227v2\full.log',
       [string]$Sel = './cmd/wisp ./internal/...')
$ErrorActionPreference = 'Continue'
$tree = 'D:\tmp\wisp\20261019\227v2\tree'
New-Item -ItemType Directory -Force -Path $Gotmp | Out-Null
$env:GOTMPDIR = $Gotmp
Remove-Item Env:TMPDIR -ErrorAction SilentlyContinue
Set-Location $tree
$env:PATH = "$tree\third_party\sherpa-onnx;$tree\build;$env:PATH"
"START=$(Get-Date -Format s) GOTMPDIR=$env:GOTMPDIR" | Out-File -FilePath $Log -Encoding utf8
"PATH_HEAD=$($env:PATH.Split(';')[0..1] -join ';')" | Out-File -FilePath $Log -Append -Encoding utf8
$pkgs = $Sel.Split(' ')
& 'D:\work\base\go\bin\go.exe' test $pkgs -count=1 -timeout 30m -v *> $Log
$rc = $LASTEXITCODE
"RC=$rc" | Out-File -FilePath $Log -Append -Encoding utf8
"END=$(Get-Date -Format s)" | Out-File -FilePath $Log -Append -Encoding utf8
Write-Output "RC=$rc LOG=$Log"
