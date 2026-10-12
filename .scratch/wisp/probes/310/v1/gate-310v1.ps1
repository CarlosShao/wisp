# 310-v1 起手闸门名册（正典＝.scratch/wisp/probes/orch/gate-roster-304.md §1）
# 规矩：每条 0 读数同行带正控（必非零那枚）；计数全部在 PS Inside 里算好再打印；
#       本文件用 -File 跑，绝不内联 -Command（本仓 10-11 那枚假 0 就是这么造的）。
$ErrorActionPreference = "Continue"

function TaskCount([string]$name) {
  $lines = @(tasklist //FI "IMAGENAME eq $name" 2>$null)
  $hits = @($lines | Where-Object { $_ -match [regex]::Escape($name) })
  return $hits.Count
}

function CimCount([string]$pattern) {
  $all = @(Get-CimInstance Win32_Process | Where-Object { $_.Name -like $pattern })
  return $all.Count
}

$svchost = CimCount "svchost.exe"
$allexe  = CimCount "*.exe"

"== 1) four exact tasklist names (each shot carries its own positive control) =="
foreach ($n in @("wisp.exe", "balldebug.exe", "mockllm.exe", "wisp.test.exe")) {
  $c = TaskCount $n
  "exact_name=$n count=$c positive_control_same_ruler=svchost.exe:$(TaskCount 'svchost.exe')"
}

"== 2) family ruler, verbatim canonical form (count computed PS-Inside) =="
$hits = @(Get-CimInstance Win32_Process | Where-Object { $_.Name -like "*.test.exe" })
"canonical_count=$($hits.Count)"
foreach ($h in $hits) { "  pid=$($h.ProcessId) ppid=$($h.ParentProcessId) name=$($h.Name) cmd=$($h.CommandLine)" }
"family_positive_control_svchost=$svchost family_positive_control_all_exe=$allexe"

"== 3) %TEMP% Test* census (the ticket-310 AC#4 roster line) =="
$env:TEMP
$troot = [System.IO.Path]::GetTempPath()
$testdirs = @((Get-ChildItem -Path $troot -Directory -ErrorAction SilentlyContinue) | Where-Object { $_.Name -like "Test*" })
$alltempdirs = @((Get-ChildItem -Path $troot -Directory -ErrorAction SilentlyContinue))
"temp_root=$troot test_count=$($testdirs.Count) positive_control_all_temp_dirs=$($alltempdirs.Count)"
$testdirs | ForEach-Object { "  name=$($_.Name) created=$($_.CreationTime.ToString('yyyy-MM-dd HH:mm:ss'))" } | Select-Object -First 80

"== 4) the two dedicated GOTMPDIR roots (present on disk, I do NOT touch them) =="
foreach ($r in @("$env:LOCALAPPDATA\wisp-test-gotmp", "D:\tmp\wisp-test-gotmp")) {
  if (Test-Path $r) {
    $kids = @((Get-ChildItem -Path $r -Directory -ErrorAction SilentlyContinue))
    "root=$r exists=1 child_dirs=$($kids.Count)"
    foreach ($k in $kids) {
      $e = @((Get-ChildItem -Path $k.FullName -ErrorAction SilentlyContinue))
      "    child=$($k.Name) entries=$($e.Count) created=$($k.CreationTime.ToString('yyyy-MM-dd HH:mm:ss'))"
    }
  } else {
    "root=$r exists=0 child_dirs=n/a positive_control_all_temp_dirs=$($alltempdirs.Count)"
  }
}
"ps_done=1"
