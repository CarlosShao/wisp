# Stops ONLY processes that this machine's forensics attributed to Wisp test runs:
# name must be mockllm.exe AND command line must reference this repo's golden dir.
# Guarded: re-verifies both conditions at kill time; anything else is left alone.
$ErrorActionPreference = 'Stop'
$out = Join-Path $PSScriptRoot 'kill-mine-1.txt'
$L = @()
$L += 'ruler = Get-CimInstance (re-verify) -> Stop-Process -Id (named pids only)'
$L += 'captured_local = ' + (Get-Date -Format 'yyyy-MM-dd HH:mm:ss zzz')
$L += ''

$needle = 'projects plans\Wisp\internal\llm\testdata\golden'
$cands = @(Get-CimInstance Win32_Process | Where-Object {
    $_.Name -eq 'mockllm.exe' -and $null -ne $_.CommandLine -and $_.CommandLine.Contains($needle)
})
$L += 'candidates matching BOTH rules (name=mockllm.exe AND this repo golden dir) = ' + $cands.Count

$killed = 0
foreach ($p in $cands) {
    $ppid = [int]$p.ParentProcessId
    $parent = @(Get-CimInstance Win32_Process -Filter "ProcessId=$ppid" -ErrorAction SilentlyContinue)
    if ($parent.Count -gt 0) {
        $L += ('  SKIP pid=' + $p.ProcessId + ' : parent pid=' + $ppid + ' is still ALIVE (' + $parent[0].Name + ') - not treating it as litter')
        continue
    }
    $L += ('  KILL pid=' + $p.ProcessId + ' started=' + ([datetime]$p.CreationDate).ToString('yyyy-MM-dd HH:mm:ss') + ' parent=' + $ppid + ' GONE, ppid re-verified absent')
    Stop-Process -Id $p.ProcessId -Force
    $killed++
}
$L += 'killed = ' + $killed
Start-Sleep -Milliseconds 800
$after = @(Get-CimInstance Win32_Process | Where-Object { $_.Name -eq 'mockllm.exe' })
$L += 'mockllm.exe count after = ' + $after.Count
$wv = @(Get-CimInstance Win32_Process | Where-Object { $_.Name -eq 'msedgewebview2.exe' })
$L += 'msedgewebview2.exe count after (must stay 24, NOT mine) = ' + $wv.Count
$L += ''
$L += 'end_of_kill_mine'
$L | Out-File -FilePath $out -Encoding utf8
Write-Output ('wrote ' + $out)
