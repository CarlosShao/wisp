param(
  [Parameter(Mandatory = $true)][string]$Exe,
  [Parameter(Mandatory = $true)][string]$WorkDir,
  [Parameter(Mandatory = $true)][string]$WatchRoot,
  [Parameter(Mandatory = $true)][string]$Label
)
# 310-r1 shape (c): external TerminateProcess of a compiled Go test binary, watched in THIS
# process (bash <-> powershell round trips are 0.5-1 s and the in-flight iteration lives
# 1-2 s, so an outside watcher cannot see the moment of death).
#
# Derived from probes/310/a1/logs/20-shape-c-watch.ps1 (census leg 310-a1, tracked file, NOT
# edited here). Two deliberate changes, both named in 20-paired-readings.md:
#   1. the watched root is a parameter, so the same台件 measures the OS temp root AND a
#      redirected GOTMPDIR root;
#   2. the leak is decided by a NAME DIFF between a snapshot taken BEFORE the process starts
#      and a snapshot taken after the kill. The census leg's own report (20-shapes.md 3.3)
#      says its first two shots read NO-LEAK because they judged the dir the loop happened to
#      catch, while the dir that actually leaks is the one whose iteration was in flight when
#      the kill landed. The name diff is the shape of ruler that caught it there, too.
# ASCII only on purpose: this file is tracked and the emoji scanner does not exempt strings.
$ErrorActionPreference = 'Stop'

function Snap([string]$root) {
  @(Get-ChildItem -Path $root -Directory -ErrorAction SilentlyContinue |
    Where-Object { $_.Name -like 'TestRoundTripLoadMarshalLoad*' } | ForEach-Object { $_.Name })
}

$start = Get-Date
$before = Snap $WatchRoot
$p = Start-Process -FilePath $Exe `
  -ArgumentList @('-test.run', '^TestRoundTripLoadMarshalLoad$', '-test.count', '4000') `
  -WorkingDirectory $WorkDir -PassThru -NoNewWindow `
  -RedirectStandardOutput "$WatchRoot-$Label.run.txt" -RedirectStandardError "$WatchRoot-$Label.run.txt.err"
Write-Output "$Label watch_root=$WatchRoot pid=$($p.Id) dirs_before=$($before.Count)"

$found = $null
$deadline = $start.AddSeconds(60)
while ($null -eq $found -and (Get-Date) -lt $deadline) {
  $fresh = @(Snap $WatchRoot) | Where-Object { $before -notcontains $_ }
  if ($fresh.Count -gt 0) { $found = $fresh[0] } else { Start-Sleep -Milliseconds 80 }
}
if ($null -eq $found) {
  Write-Output "$Label world=NOT-CONSTRUCTED reading=VOID"
  Stop-Process -Id $p.Id -Force
  exit 2
}
Write-Output "$Label world=CONSTRUCTED dir=$found"
$b = @(Get-CimInstance Win32_Process -Filter "ProcessId=$($p.Id)")
if ($b.Count -eq 0) { Write-Output "$Label target_vanished_before_kill reading=VOID"; exit 2 }
$g1 = ($b[0].Name -like '*.test.exe')
$g2 = ($b[0].CommandLine -like "*$WorkDir*")
$g3 = ($b[0].ParentProcessId -eq $PID)
Write-Output "$Label guard_name_is_test_exe=$g1 guard_cmdline_has_workdir=$g2 guard_parent_is_this_leg=$g3"
if (-not ($g1 -and $g2 -and $g3)) {
  Write-Output "$Label guards_failed reading=VOID"
  Stop-Process -Id $p.Id -Force
  exit 3
}
Stop-Process -Id $p.Id -Force
Write-Output "$Label external_kill_issued at=$(Get-Date -Format 'HH:mm:ss.fff')"
Start-Sleep -Milliseconds 1500
$alive = @(Get-Process -Id $p.Id -ErrorAction SilentlyContinue).Count
$after = @(Snap $WatchRoot)
$leaked = @($after | Where-Object { $before -notcontains $_ })
$topNow = @(Get-ChildItem -Path $env:TEMP -Directory | Where-Object { $_.Name -like 'Test*' }).Count
Write-Output "$Label after_kill target_pid_still_present=$alive dirs_before=$($before.Count) dirs_after=$($after.Count) leaked_count=$($leaked.Count) leaked_names=[$($leaked -join ',')]"
Write-Output "$Label os_temp_test_star_now=$topNow"
if ($leaked.Count -gt 0) { Write-Output "$Label verdict=DIR-LEAKED-in-watch-root shape_c_reproduced"; exit 0 }
Write-Output "$Label verdict=NO-LEAK-in-watch-root"; exit 1
