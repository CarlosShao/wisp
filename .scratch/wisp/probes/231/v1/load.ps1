$cpu = (Get-Counter '\Processor(_Total)\% Processor Time' -SampleInterval 1 -MaxSamples 1).CounterSamples[0].CookedValue
$os = Get-CimInstance Win32_OperatingSystem
$mem = (1 - $os.FreePhysicalMemory / $os.TotalVisibleMemorySize) * 100
Write-Output ("CPU={0:N1} MEM={1:N1}" -f $cpu, $mem)
