param([Parameter(Mandatory)][string]$Exe, [string]$OutFile)
# 263-v1 最小复现件：把改前 scripts/slo-check.ps1 的 :324/:325/:326 三条语句逐字搬出来
#   Write-Host 'sampling state ...'   -> & $WispExe slo ...   -> $code = $LASTEXITCODE
# 目的：亲手量"PowerShell 的 & 对 GUI 子系统进程等不等"，两形对照（CUI vs GUI）。
$ErrorActionPreference = 'Stop'
Set-StrictMode -Version 2.0
$label = Split-Path -Leaf $Exe
Write-Host "v1-repro: sampling state Sleeping for 2s ($label)"
$sw = [System.Diagnostics.Stopwatch]::StartNew()
if ($OutFile) {
    & $Exe slo -state Sleeping -seconds 2 -interval-ms 250 -out $OutFile
} else {
    & $Exe slo -state Sleeping -seconds 2 -interval-ms 250 -out (Join-Path $env:TEMP 'v1-repro-state.json')
}
$elapsed = $sw.ElapsedMilliseconds
Write-Host "v1-repro: call returned after ${elapsed}ms"
$code = $LASTEXITCODE
Write-Host "v1-repro: state exit=$code elapsed-ms=$elapsed"
