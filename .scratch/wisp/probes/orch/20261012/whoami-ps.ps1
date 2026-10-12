$ErrorActionPreference = 'Continue'
$u = [System.Security.Principal.WindowsIdentity]::GetCurrent()
Write-Output ("token_user=" + $u.Name)
Write-Output ("token_user_sid=" + $u.User.Value)
Write-Output ("token_groups_count=" + $u.Groups.Count)
foreach ($g in $u.Groups) {
  if ($g.Value -in @('S-1-5-32-544','S-1-5-18','S-1-5-11','S-1-5-32-546')) {
    Write-Output ("group_sid=" + $g.Value)
  }
}
$b = New-Object System.Security.Principal.NTAccount('BUILTIN','Administrators')
try {
  $t = $b.Translate([System.Security.Principal.SecurityIdentifier])
  Write-Output ("BA_resolves=" + $t.Value)
} catch { Write-Output ("BA_translate_failed=" + $_.Exception.Message) }
$la = New-Object System.Security.Principal.NTAccount('BUILTIN','Administrators')
Write-Output ("deny_only_check: BA present in token above? " + ($u.Groups | Where-Object { $_.Value -eq $t.Value } | Measure-Object).Count)
Write-Output "whoami_ps_rc=0"
