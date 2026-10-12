param([string]$Root,[string]$TempPath)
"ROOT=" + $Root
"ROOT_EXISTS_DIR=" + (Test-Path -Path $Root -PathType Container)
$full = (Get-Item -Force $Root).FullName
"ROOT_FULL=" + $full
$t = (Get-Item -Force $TempPath).FullName.TrimEnd('\')
$underTemp = $full.StartsWith($t + '\')
"ROOT_UNDER_TEMP=" + $underTemp
"ROOT_TEMPNAME_IN_TEMP_ROSTER=" + [bool](Get-ChildItem -Path $env:TEMP -Directory | Where-Object { $_.Name -eq (Split-Path $Root -Leaf) })
