$os = Get-CimInstance Win32_OperatingSystem
Write-Output ("TotalGB={0:N1} FreeGB={1:N1}" -f ($os.TotalVisibleMemorySize / 1MB), ($os.FreePhysicalMemory / 1MB))
Get-Process | Sort-Object WorkingSet64 -Descending | Select-Object -First 8 Name, @{n = 'MB'; e = { [int]($_.WorkingSet64 / 1MB) } } | Format-Table -AutoSize
$g = Get-Process go, wisp, test -ErrorAction SilentlyContinue
if ($g) { Write-Output ("INFLIGHT_GO_OR_WISP=" + $g.Count) } else { Write-Output "INFLIGHT_GO_OR_WISP=0" }
