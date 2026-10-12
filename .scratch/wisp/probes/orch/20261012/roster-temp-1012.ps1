$ErrorActionPreference = 'Continue'
$tops = @(Get-ChildItem -Path $env:TEMP -Directory | Where-Object { $_.Name -like 'Test*' })
Write-Output ("temp_test_star_count=" + $tops.Count)
Write-Output ("temp_all_dirs_positive_control=" + @(Get-ChildItem -Path $env:TEMP -Directory).Count)
Write-Output "--- roster (name|CreationTime|lastWrite) ---"
foreach ($d in ($tops | Sort-Object CreationTime)) {
  Write-Output ($d.Name + "|" + $d.CreationTime.ToString("yyyy-MM-dd HH:mm:ss") + "|" + $d.LastWriteTime.ToString("yyyy-MM-dd HH:mm:ss"))
}
foreach ($root in @("$env:LOCALAPPDATA\wisp-test-gotmp", "D:\tmp\wisp-test-gotmp")) {
  if (Test-Path $root) {
    $rs = @(Get-ChildItem -Path $root -Directory)
    Write-Output ("root=" + $root + " run_dirs=" + $rs.Count)
    $sum = 0
    foreach ($r in $rs) {
      $e = @(Get-ChildItem -Path $r.FullName -Directory).Count
      $f = @(Get-ChildItem -Path $r.FullName -File).Count
      $sum = $sum + $e
      Write-Output ("  " + $r.Name + "|entry_dirs=" + $e + "|entry_files=" + $f + "|created=" + $r.CreationTime.ToString("yyyy-MM-dd HH:mm:ss"))
    }
    Write-Output ("  root_total_entry_dirs=" + $sum)
  } else {
    Write-Output ("root=" + $root + " ABSENT")
  }
}
Write-Output "roster_rc=0"
