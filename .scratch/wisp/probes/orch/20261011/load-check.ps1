$os = Get-CimInstance Win32_OperatingSystem
$cpu = (Get-CimInstance Win32_Processor | Measure-Object -Property LoadPercentage -Average).Average
$memPct = [math]::Round((($os.TotalVisibleMemorySize - $os.FreePhysicalMemory) * 100 / $os.TotalVisibleMemorySize), 1)
Write-Output ("CPU_loadpct=" + $cpu + " MEM_usedpct=" + $memPct)
foreach ($n in @('go','gopls','gotest','compile','link','wisp','node','runner','git')) {
  $c = (Get-Process -Name ($n + '*') -ErrorAction SilentlyContinue | Measure-Object).Count
  Write-Output ($n + "*=" + $c)
}
