$all = @(Get-CimInstance Win32_Process -Filter "Name='msedgewebview2.exe'")
$byPid = @{}
foreach ($p in $all) { $byPid[[string]$p.ProcessId] = $p }
function RootOf($p) { $cur = $p; $g = 0
  while ($true) { $par = $byPid[[string]$cur.ParentProcessId]; if ($null -eq $par) { return $cur }; $cur = $par; $g++; if ($g -gt 20) { return $cur } } }
$groups = @{}
foreach ($p in $all) { $r = RootOf $p; $k = [string]$r.ProcessId
  if (-not $groups.ContainsKey($k)) { $groups[$k] = New-Object System.Collections.ArrayList }; [void]$groups[$k].Add($p) }
Write-Output ("captured=" + (Get-Date -Format 'yyyy-MM-dd HH:mm:ss'))
Write-Output ("total=" + $all.Count + " trees=" + $groups.Count)
foreach ($k in ($groups.Keys | Sort-Object { [int]$_ })) {
  $root = $byPid[$k]; $m = $groups[$k]
  $nm = [regex]::Match($root.CommandLine, '--webview-exe-name=([^\s"]+)')
  $hostNm = $nm.Groups[1].Value
  $pp = Get-CimInstance Win32_Process -Filter "ProcessId=$($root.ParentProcessId)" -ErrorAction SilentlyContinue
  $ppn = if ($pp) { $pp.Name } else { "GONE(ppid=" + $root.ParentProcessId + ")" }
  $start = $root.CreationDate.ToString('MM-dd HH:mm:ss')
  $wsMB = [math]::Round((($m | Measure-Object -Property WorkingSetSize -Sum).Sum/1MB),1)
  Write-Output ("  tree root=" + $k + " host=" + $hostNm + " parent=" + $ppn + " started=" + $start + " size=" + $m.Count + " wsMB=" + $wsMB)
}
