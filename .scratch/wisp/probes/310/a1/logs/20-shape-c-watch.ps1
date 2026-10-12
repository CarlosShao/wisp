param(
  [Parameter(Mandatory=$true)][string]$Exe,
  [Parameter(Mandatory=$true)][string]$WorkDir,
  [Parameter(Mandatory=$true)][string]$LogFile,
  [string]$Label = 'shapec'
)
# 310-a1 shape (c): start a compiled Go test binary from THIS process, watch %TEMP% for the
# directory its in-flight t.TempDir() creates, then TerminateProcess it from the SAME process
# (no bash<->powershell roundtrip between watch and kill), then re-check the directory.
# ASCII only on purpose: this file is tracked, and the emoji scanner does not exempt strings.
$ErrorActionPreference = 'Stop'
$start = Get-Date
$stdout = "$LogFile.run.txt"
$p = Start-Process -FilePath $Exe `
  -ArgumentList @('-test.run', '^TestRoundTripLoadMarshalLoad$', '-test.count', '4000') `
  -WorkingDirectory $WorkDir -PassThru -NoNewWindow `
  -RedirectStandardOutput $stdout -RedirectStandardError "$stdout.err"
Write-Output "$Label started pid=$($p.Id) name=$($p.Name) at=$($start.ToString('HH:mm:ss.fff'))"
$found = $null
$deadline = $start.AddSeconds(60)
while ($null -eq $found -and (Get-Date) -lt $deadline) {
  $found = @(Get-ChildItem -Path $env:TEMP -Directory |
    Where-Object { $_.Name -like 'TestRoundTripLoadMarshalLoad*' -and $_.CreationTime -ge $start }) |
    Select-Object -First 1
  if ($null -eq $found) { Start-Sleep -Milliseconds 80 }
}
if ($null -eq $found) {
  Write-Output "$Label world=NOT-CONSTRUCTED reading=VOID killing_target_only"
  Stop-Process -Id $p.Id -Force
  exit 2
}
Write-Output "$Label world=CONSTRUCTED dir=$($found.Name) dir_creation=$($found.CreationTime.ToString('HH:mm:ss.fff'))"
$b = @(Get-CimInstance Win32_Process -Filter "ProcessId=$($p.Id)")
if ($b.Count -eq 0) { Write-Output "$Label target_vanished_before_kill reading=VOID"; exit 2 }
$g1 = ($b[0].Name -like '*.test.exe')
$g2 = ($b[0].CommandLine -like "*$WorkDir*")
$g3 = ($b[0].ParentProcessId -eq $PID)
Write-Output "$Label target_name=$($b[0].Name) pid=$($p.Id) ppid=$($b[0].ParentProcessId) this_ps_pid=$PID"
Write-Output "$Label cmdline=$($b[0].CommandLine)"
Write-Output "$Label guard_name_is_test_exe=$g1 guard_cmdline_has_export_tree=$g2 guard_parent_is_this_leg=$g3"
if (-not ($g1 -and $g2 -and $g3)) {
  Write-Output "$Label guards_failed refusing_to_kill reading=VOID"
  Stop-Process -Id $p.Id -Force
  exit 3
}
Stop-Process -Id $p.Id -Force
Write-Output "$Label external_kill_issued at=$(Get-Date -Format 'HH:mm:ss.fff')"
# The watched dir is NOT the one that necessarily leaks: iterations are ~10 ms apart, so the
# first dir the loop saw may already have been cleaned by its own finished iteration when the
# kill lands. Snapshot every TestRoundTripLoadMarshalLoad* dir at kill time and diff after.
$atKill = @(Get-ChildItem -Path $env:TEMP -Directory |
  Where-Object { $_.Name -like 'TestRoundTripLoadMarshalLoad*' } | ForEach-Object { $_.Name })
Write-Output "$Label at_kill_mine_dirs=$($atKill.Count) at_kill_names=$($atKill -join ',')"
Start-Sleep -Milliseconds 1500
$alive = @(Get-Process -Id $p.Id -ErrorAction SilentlyContinue).Count
$afterKill = @(Get-ChildItem -Path $env:TEMP -Directory |
  Where-Object { $_.Name -like 'TestRoundTripLoadMarshalLoad*' } | ForEach-Object { $_.Name })
$leaked = @($afterKill | Where-Object { $atKill -notcontains $_ })
$allNow = @(Get-ChildItem -Path $env:TEMP -Directory | Where-Object { $_.Name -like 'Test*' }).Count
Write-Output "$Label after_kill target_pid_still_present=$alive mine_dirs_after=$($afterKill.Count) leaked_names=[$($leaked -join ',')] leaked_count=$($leaked.Count) test_star_dirs_now=$allNow"
Write-Output "$Label watched_dir_now_exists=$(@(Get-ChildItem -Path $env:TEMP -Directory | Where-Object { $_.Name -eq $found.Name }).Count)"
if ($leaked.Count -gt 0) { Write-Output "$Label verdict=DIR-LEAKED shape_c_reproduced"; exit 0 }
Write-Output "$Label verdict=NO-LEAK"; exit 1
