# roster.ps1 - 311-r1b: process + temp-root roster with positive controls on the same line.
# Rulers are deliberately paired with a control so a "0" can never be read as "empty output".
$ErrorActionPreference = 'SilentlyContinue'
$all = Get-Process
$family = @($all | Where-Object { $_.ProcessName -like '*.test' })
$svchost = @($all | Where-Object { $_.ProcessName -eq 'svchost' })
$wispexe = @($all | Where-Object { $_.ProcessName -in @('wisp', 'balldebug') })
Write-Output ("all_processes_positive_control=" + $all.Count)
Write-Output ("svchost_positive_control=" + $svchost.Count)
Write-Output ("family_test_exe=" + $family.Count)
Write-Output ("wisp_or_balldebug_exe=" + $wispexe.Count)
foreach ($p in $family) {
    $ct = 'unknown'
    try { $ct = $p.StartTime.ToString('yyyy-MM-dd HH:mm:ss') } catch { $ct = 'unreadable' }
    $cl = 'unreadable'
    try { $cl = (Get-CimInstance Win32_Process -Filter "ProcessId=$($p.Id)").CommandLine } catch {}
    Write-Output ("family_pid=" + $p.Id + " name=" + $p.ProcessName + " started=" + $ct + " cmdline=" + $cl)
}
foreach ($p in $wispexe) { Write-Output ("wispish_pid=" + $p.Id + " name=" + $p.ProcessName) }
$tempRoot = [System.IO.Path]::GetFullPath($env:TEMP)
$dirs = @(Get-ChildItem -LiteralPath $tempRoot -Directory)
$testStar = @($dirs | Where-Object { $_.Name -like 'Test*' })
Write-Output ("temp_root=" + $tempRoot)
Write-Output ("temp_test_star=" + $testStar.Count)
Write-Output ("temp_all_dirs_positive_control=" + $dirs.Count)
foreach ($d in $testStar) { Write-Output ("temp_test_dir name=" + $d.Name + " created=" + $d.CreationTime.ToString('yyyy-MM-dd HH:mm:ss')) }
