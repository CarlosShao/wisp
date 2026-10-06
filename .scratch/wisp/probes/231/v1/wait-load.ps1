param([int]$Threshold = 70, [int]$MaxWaits = 12)
for ($i = 0; $i -lt $MaxWaits; $i++) {
    $cpu = (Get-Counter '\Processor(_Total)\% Processor Time' -SampleInterval 1 -MaxSamples 1).CounterSamples[0].CookedValue
    $os = Get-CimInstance Win32_OperatingSystem
    $mem = (1 - $os.FreePhysicalMemory / $os.TotalVisibleMemorySize) * 100
    $line = ("try={0} CPU={1:N1} MEM={2:N1}" -f $i, $cpu, $mem)
    Write-Output $line
    if ($cpu -lt $Threshold) {
        Write-Output ("CPU_UNDER_THRESHOLD after {0} waits" -f $i)
        exit 0
    }
    Start-Sleep -Seconds 20
}
Write-Output "CPU_STILL_HIGH - proceeding with named caveat"
exit 1
