param(
  [string]$Pattern = "*.test.exe",
  [switch]$Quick
)
# 编排者 10-12 复跑：① 家族尺逐字（正典 gate-roster-304.md §1）② 正控（同管线必命中）③ AC#4b 归属三问在非零台面上咬不咬
$ErrorActionPreference = "Stop"

function FamilyCount([string]$p) {
  $all = @(Get-CimInstance Win32_Process | Where-Object { $_.Name -like $p })
  $ids = ($all | ForEach-Object { $_.ProcessId } | Sort-Object) -join ","
  $names = ($all | ForEach-Object { $_.Name } | Sort-Object -Unique) -join ","
  "pattern=$p count=$($all.Count) ids=[$ids] uniq_names=[$names]"
}

"Witness: family ruler (verbatim canonical form, no inline -join in the PS command line)"
$hits = @(Get-CimInstance Win32_Process | Where-Object { $_.Name -like $Pattern })
"canonical_count=$($hits.Count)"
foreach ($h in $hits) { "  pid=$($h.ProcessId) ppid=$($h.ParentProcessId) name=$($h.Name) cmd=$($h.CommandLine)" }
""
"Positive controls (same pipeline, a pattern that MUST be non-zero today):"
FamilyCount -p "svchost.exe"
FamilyCount -p "explorer.exe"
FamilyCount -p "*.exe" | ForEach-Object { $_.Substring(0, [Math]::Min(120, $_.Length)) }
""
"AC#4b bite test: is there a NON-ZERO surface where 'point to a tree or write unknown' actually bites?"
$wv = @(Get-CimInstance Win32_Process | Where-Object { $_.Name -like "*msedgewebview2.exe" })
"webview_total=$($wv.Count)"
$noarg = @($wv | Where-Object { $_.CommandLine -notmatch "--webview-exe-name=" })
"webview_without_webview_exe_name=$($noarg.Count)"
foreach ($x in $noarg) { "  pid=$($x.ProcessId) ppid=$($x.ParentProcessId) cmdlen=$([string]::Length($x.CommandLine)) cmd=$($x.CommandLine)" }
$parentless = @($wv | Where-Object { -not (Get-CimInstance Win32_Process -Filter "ProcessId=$($_.ParentProcessId)" -ErrorAction SilentlyContinue) })
"webview_parent_already_gone=$($parentless.Count)"
""
"Roots of the webview trees (parent chain one hop up, by name):"
$wv | ForEach-Object { $_.ParentProcessId } | Sort-Object -Unique | ForEach-Object {
  $p = Get-CimInstance Win32_Process -Filter "ProcessId=$_" -ErrorAction SilentlyContinue
  if ($p) { "  root_pid=$_ name=$($p.Name) cmd=$($p.CommandLine)" } else { "  root_pid=$_ name=<GONE>" }
}
if ($Quick) { exit 0 }
