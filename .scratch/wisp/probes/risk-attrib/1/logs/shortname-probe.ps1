$ErrorActionPreference = 'Stop'
$long = 'C:\Users\swq\tmp\riskattrib1\RunnerTempAreaForAttribution1'
$fso = New-Object -ComObject Scripting.FileSystemObject
$folder = $fso.GetFolder($long)
Write-Output ('LONG   = ' + $long)
Write-Output ('SHORT  = ' + $folder.ShortPath)
Write-Output ('8DOT3  = ' + (cmd /c fsutil 8dot3name query C: | Select-Object -Last 1))
