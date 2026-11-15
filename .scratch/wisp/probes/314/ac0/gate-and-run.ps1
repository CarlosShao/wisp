# Ticket 314 AC#0 — exclusivity gate, then the canonical direct run.
# The fleet rule (ticket 314 §3, ledger A1860): only ONE full `cmd/wisp` suite may
# run at a time. A run started while another is in flight is not a reading.
# This gate waits until no wisp.test.exe / wisp-winlive.test.exe is alive for a
# full 60s, then starts the canonical staged direct run.
param([string]$Suffix = 'c')
$ErrorActionPreference = 'Continue'
$dir = 'D:\tmp\wisp\20261115\w314'
$log = Join-Path $dir "gate-$Suffix.log"
"GATE_START=$(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')" | Set-Content -Path $log -Encoding utf8

function Test-Competitor {
  $p = Get-CimInstance Win32_Process -Filter "Name='wisp.test.exe' OR Name='wisp-winlive.test.exe'" -ErrorAction SilentlyContinue
  if ($p) { return ($p | ForEach-Object { "pid=$($_.ProcessId) cl=$($_.CommandLine)" }) -join ' | ' }
  return ''
}

$freeStreak = 0
while ($true) {
  $c = Test-Competitor
  if ($c -eq '') {
    $freeStreak++
  } else {
    $freeStreak = 0
    "GATE_BUSY ts=$(Get-Date -Format 'HH:mm:ss') $c" | Add-Content -Path $log -Encoding utf8
  }
  if ($freeStreak -ge 3) { break }
  Start-Sleep -Seconds 20
}
"GATE_FREE ts=$(Get-Date -Format 'yyyy-MM-dd HH:mm:ss') (60s with no test binary alive)" | Add-Content -Path $log -Encoding utf8
Start-Sleep -Seconds 5
"GATE_LAUNCH ts=$(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')" | Add-Content -Path $log -Encoding utf8

pwsh -NoProfile -File 'D:\work\workspace\projects plans\Wisp\.scratch\wisp\probes\314\ac0\run-direct-canonical.ps1' -Suffix $Suffix |
  Tee-Object -FilePath (Join-Path $dir "gate-run-$Suffix.out")
"GATE_DONE ts=$(Get-Date -Format 'yyyy-MM-dd HH:mm:ss') rc=$LASTEXITCODE" | Add-Content -Path $log -Encoding utf8
