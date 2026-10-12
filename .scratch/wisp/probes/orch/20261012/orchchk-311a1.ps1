$ErrorActionPreference = 'Continue'
$ts = Get-Date -Format 'HHmmss'
$rootA = "C:\Users\swq\AppData\Local\wisp-test-gotmp\orchchk-A-$ts"
$rootB = "D:\tmp\wisp-test-gotmp\orchchk-B-$ts"
New-Item -ItemType Directory -Path $rootA | Out-Null
New-Item -ItemType Directory -Path $rootB | Out-Null
Write-Output ("whoami_user_sid=" + (whoami /user /fo csv | ConvertFrom-Csv).Sid)
$groups = whoami /groups /fo csv | ConvertFrom-Csv
Write-Output ("group_BA_line=" + (($groups | Where-Object { $_.SID -eq 'S-1-5-32-544' }) -join ' '))
Write-Output ("group_SY_present=" + [bool]($groups | Where-Object { $_.SID -eq 'S-1-5-18' }))
$lu = Get-LocalUser | Where-Object { $_.SID.Value -like '*-500' }
Write-Output ("localadmin_500 name=" + $lu.Name + " enabled=" + $lu.Enabled + " sid=" + $lu.SID.Value)
foreach ($pair in @(@('A', $rootA), @('B', $rootB))) {
  $tag = $pair[0]; $r = $pair[1]
  foreach ($p in @($r, (Split-Path $r -Parent), (Split-Path (Split-Path $r -Parent) -Parent))) {
    if (-not $p) { continue }
    $txt = (icacls $p 2>&1 | Out-String)
    $mine = ($txt -split "\r?\n") | Where-Object { $_ -match 'swq|Authenticated Users|BUILTIN\\Users' }
    Write-Output ("chain_" + $tag + " dir=" + $p + " Dc_present=" + ($txt -match '\(Dc\)'))
    foreach ($m in $mine) { Write-Output ("   ace: " + $m.Trim()) }
  }
  $acl = Get-Acl -Path $r
  Write-Output ("root_" + $tag + " owner=" + $acl.Owner + " protected=" + $acl.AreAccessRulesProtected)
}
Write-Output "orchchk_rc=0"
