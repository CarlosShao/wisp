$ErrorActionPreference = 'Continue'
$leaves = @(
 'D:\tmp\wisp-test-gotmp\311a1-B-20261012-133040-1594\TestSealNarrowsAndNamesThePrincipalItRemovedBySID1250499198\002\la',
 'D:\tmp\wisp-test-gotmp\311a1-B-20261012-133040-1594\TestGateJudgesThePrivateSetByResolvedSID3972594964\001\la',
 'D:\tmp\wisp-test-gotmp\311a1-B-20261012-133040-1594\TestGateJudgesThePrivateSetByResolvedSIDrunner_descriptor_is_ju1386323932\001\runner'
)
foreach ($p in $leaves) {
  if (Test-Path -LiteralPath $p) {
    try {
      $a = Get-Acl -LiteralPath $p
      Write-Output ("leaf_present=1 owner=" + $a.Owner + " protected=" + $a.AreAccessRulesProtected + " sddl=" + $a.Sddl)
      Write-Output ("  path=" + $p)
    } catch { Write-Output ("  getacl_failed=" + $_.Exception.Message) }
    try {
      $r = [System.IO.Directory]::GetFiles($p)
      Write-Output ("  can_list_inside=1 n=" + $r.Count)
    } catch { Write-Output ("  can_list_inside=0 err=" + $_.Exception.Message) }
  } else {
    Write-Output ("leaf_present=0 path=" + $p)
  }
}
$g = cmd /c "whoami /groups /fo csv /nh"
Write-Output ("ba_row=" + (($g | Select-String -SimpleMatch 'S-1-5-32-544') -join ' '))
Write-Output ("sy_row_present=" + [bool]($g | Select-String -SimpleMatch 'S-1-5-18'))
foreach ($root in @('C:\Users\swq\AppData\Local\wisp-test-gotmp', 'D:\tmp\wisp-test-gotmp')) {
  $rs = @(Get-ChildItem -LiteralPath $root -Directory)
  $sum = 0
  foreach ($x in $rs) { $sum += @(Get-ChildItem -LiteralPath $x.FullName -Directory SilentlyContinue).Count }
  Write-Output ("root=" + $root + " run_dirs=" + $rs.Count + " entry_dirs_sum=" + $sum)
}
$tops = @(Get-ChildItem -Path $env:TEMP -Directory | Where-Object { $_.Name -like 'Test*' })
Write-Output ("temp_test_star=" + $tops.Count)
Write-Output ("temp_all_dirs_positive_control=" + @(Get-ChildItem -Path $env:TEMP -Directory).Count)
Write-Output "leafcheck_rc=0"
