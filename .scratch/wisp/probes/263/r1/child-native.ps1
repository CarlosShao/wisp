$ErrorActionPreference = 'Stop'
Set-StrictMode -Version 2.0
Write-Host 'repro-child-native: running cmd.exe /c exit 5, then returning without exit'
& cmd.exe /c 'exit 5'
Write-Host 'repro-child-native: returned'
