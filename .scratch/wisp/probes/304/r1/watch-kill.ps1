param(
  [Parameter(Mandatory = $true)][int]$TestPid,
  [int]$MaxSeconds = 25
)
# 304-r1 instrument: atomically OBSERVE a live mockllm.exe child of $TestPid and
# TerminateProcess THAT pid from inside one PowerShell process.
# Why: observing with one call and killing with another costs ~2-3 s of process
# startup on this box, which is long enough for a short test to reach its own
# t.Cleanup - the reverse shot then "measures" a world that never happened
# (that is exactly how the first c-SP1 reverse shot died: Stop-Process reported
# "Cannot find a process with the process identifier 22456").
# Scope: only the pid handed in (started by paired.sh) is ever stopped; children
# are only matched by name 'mockllm.exe' AND a command line containing this
# repo's golden dir AND ppid == $TestPid.
$ErrorActionPreference = 'Stop'
$goldenPat = '*Wisp\internal\llm\testdata\golden*'
$deadline = (Get-Date).AddSeconds($MaxSeconds)
while ((Get-Date) -lt $deadline) {
  $kids = @(Get-CimInstance Win32_Process -Filter "Name='mockllm.exe'" |
      Where-Object { (([string]$_.CommandLine -like $goldenPat) -and ($_.ParentProcessId -eq $TestPid)) })
  if ($kids.Count -gt 0) {
    foreach ($k in $kids) {
      Write-Output ("world_constructed: child pid=" + $k.ProcessId + " ppid=" + $k.ParentProcessId +
        " start=" + $k.CreationDate.ToString('yyyy-MM-dd HH:mm:ss') + " cmdline_has_repo_golden=YES")
    }
    Write-Output ("parent_present_before_stop=" + [bool](Get-Process -Id $TestPid -ErrorAction SilentlyContinue))
    try {
      Stop-Process -Id $TestPid -Force
      Write-Output ("terminate_issued pid=" + $TestPid + " (the only pid this script stops)")
    } catch {
      Write-Output ("terminate_FAILED pid=" + $TestPid + " err=" + $_.Exception.Message)
      exit 4
    }
    Start-Sleep -Milliseconds 400
    Write-Output ("parent_present_after_stop=" + [bool](Get-Process -Id $TestPid -ErrorAction SilentlyContinue))
    exit 0
  }
  Start-Sleep -Milliseconds 200
}
Write-Output ("NOT-CONSTRUCTED within " + $MaxSeconds + "s: no mockllm.exe child of pid " + $TestPid + " was ever observed")
exit 3
