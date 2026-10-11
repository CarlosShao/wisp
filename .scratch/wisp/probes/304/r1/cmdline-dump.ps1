# 304-r1 instrument: dump pid/ppid/full command line of every mockllm.exe whose
# command line contains THIS repo's golden dir.
# Purpose: the ticket's own attribution ruler (issue 304 lines 11-14) - the first
# path segment of the child's t.TempDir() is the TEST NAME, so a residue is
# assigned to a named spawn point by its command line, not by a count or by
# adjacency in time.
$ErrorActionPreference = 'Stop'
$goldenPat = '*Wisp\internal\llm\testdata\golden*'
$rows = @(Get-CimInstance Win32_Process -Filter "Name='mockllm.exe'" |
    Where-Object { [string]$_.CommandLine -like $goldenPat })
foreach ($r in $rows) {
  $parentAlive = [bool](Get-Process -Id $r.ParentProcessId -ErrorAction SilentlyContinue)
  Write-Output ("pid=" + $r.ProcessId + " ppid=" + $r.ParentProcessId + " parent_alive=" + $parentAlive +
    " tempdir_testname=" + (($r.CommandLine -split '\\') | Where-Object { $_ -like 'Test*' } | Select-Object -First 1) +
    " cmdline=" + $r.CommandLine)
}
Write-Output ("count_repo_golden=" + $rows.Count)
