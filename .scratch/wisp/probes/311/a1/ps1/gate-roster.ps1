# Probe 311-a1: gate roster + token identity (read-only, no ACL writes)
$ErrorActionPreference = 'Continue'

Write-Output ("now=" + (Get-Date -Format 'yyyy-MM-dd HH:mm:ss'))
Write-Output ("temp_env=" + $env:TEMP)

# ruler A: top-level directories under %TEMP% whose name starts with "Test"
$all = @(Get-ChildItem -Path $env:TEMP -Directory -ErrorAction SilentlyContinue)
$tests = @($all | Where-Object { $_.Name -like 'Test*' })
Write-Output ("temp_top_all_count=" + $all.Count)
Write-Output ("temp_top_TestRoster_count=" + $tests.Count)

# ruler B: the *.test.exe family (process level), count computed inside PS
$procs = @(Get-CimInstance Win32_Process -ErrorAction SilentlyContinue)
$family = @($procs | Where-Object { $_.Name -like '*.test.exe' })
Write-Output ("testexe_family_count=" + $family.Count)
foreach ($p in $family) {
    Write-Output ("FAMILY pid=" + $p.ProcessId + " ppid=" + $p.ParentProcessId + " name=" + $p.Name + " cmd=" + $p.CommandLine)
}

# positive controls on the same pipeline (must be non-zero)
$svchost = @($procs | Where-Object { $_.Name -eq 'svchost.exe' })
Write-Output ("positive_control_svchost_count=" + $svchost.Count)
$allexe = @($procs | Where-Object { $_.Name -like '*.exe' })
Write-Output ("positive_control_all_exe_count=" + $allexe.Count)

# token identity: which SIDs does this process actually hold
$id = [Security.Principal.WindowsIdentity]::GetCurrent()
Write-Output ("token_name=" + $id.Name)
Write-Output ("token_user_sid=" + $id.User.Value)
foreach ($g in $id.Groups) { Write-Output ("token_group_sid=" + $g.Value) }

# local account table: which name holds RID 500 (the account SDDL spells "LA")
$acct = @(Get-CimInstance Win32_UserAccount -ErrorAction SilentlyContinue | Where-Object { $_.LocalAccount -eq $true })
Write-Output ("localaccount_count=" + $acct.Count)
foreach ($a in $acct) {
    $rid = ($a.SID -split '-')[-1]
    Write-Output ("LOCALACCOUNT name=" + $a.Name + " sid=" + $a.SID + " rid=" + $rid + " disabled=" + $a.Disabled)
}
Write-Output ("ps_failed=" + ($Error.Count))
