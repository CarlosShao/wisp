# 311-r1 starting/handover roster ruler (ticket 311, leg 311-r1)
# Reads, all printed as key=value lines:
#   temp_test_star   = top-level Test* directories under $env:TEMP (the ticket-310/311 registered residue roster)
#   temp_all_dirs    = positive control: every top-level directory under $env:TEMP (a plain 0 here means a broken ruler)
#   family_test_exe  = live processes whose name ends in .test (Go test binaries) -- must be 0 at handover
#   svchost / all_exe = positive controls so an empty process list is never read as "no processes"
# Get-ChildItem carries -ErrorAction SilentlyContinue (without it a denied path reads as a false 0).
$ErrorActionPreference = 'Continue'
$t = Get-ChildItem $env:TEMP -Directory -ErrorAction SilentlyContinue
$testStar = @($t | Where-Object { $_.Name -like 'Test*' })
$allDirs = @($t)
$procs = @(Get-Process -ErrorAction SilentlyContinue)
$family = @($procs | Where-Object { $_.ProcessName -like '*.test' })
$svchost = @($procs | Where-Object { $_.ProcessName -eq 'svchost' })
$allExe = @($procs)
"temp_test_star=$($testStar.Count)"
"temp_all_dirs_positive_control=$($allDirs.Count)"
"family_test_exe=$($family.Count)"
"svchost_positive_control=$($svchost.Count)"
"all_processes_positive_control=$($allExe.Count)"
if ($family.Count -gt 0) {
    foreach ($p in $family) { "family_pid=$($p.Id) name=$($p.ProcessName) started=$($p.StartTime.ToString('yyyy-MM-dd HH:mm:ss'))" }
}
