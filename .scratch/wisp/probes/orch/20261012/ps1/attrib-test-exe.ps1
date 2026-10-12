$ErrorActionPreference = 'SilentlyContinue'
$rows = @(Get-CimInstance Win32_Process -Filter "Name like '%test.exe%'")
"test_exe_count=$($rows.Count)"
foreach ($r in $rows) {
  $cl = $r.CommandLine
  if ($null -eq $cl) { $cl = '(unreadable)' }
  $root = 'n/a'
  if ($cl -match 'wisp-test-gotmp[\\/](?<rn>[^\\/ "\]+)') { $root = $Matches['rn'] }
  $pkg = 'n/a'
  if ($r.Name -match '^(?<p>[^.]+)\.test\.exe$') { $pkg = $Matches['p'] }
  "pid=$($r.ProcessId) name=$($r.Name) pkg_hint=$pkg run_root=$root"
}
$parents = Get-CimInstance Win32_Process -Filter "Name = 'go.exe'"
"go_exe_count=$($parents.Count)"
