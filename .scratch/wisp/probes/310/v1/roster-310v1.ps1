# 310-v1 %TEMP% Test* roster ruler (canonical shape: names + CreationTime, count computed PS-Inside)
# usage: powershell -NoProfile -ExecutionPolicy Bypass -File this.ps1 -Out <name-file> [-Round <label>]
param(
  [Parameter(Mandatory = $true)][string]$Out,
  [string]$Round = "n/a",
  [string]$Filter = "Test*"
)
$ErrorActionPreference = "Continue"
$troot = [System.IO.Path]::GetTempPath()
$all = @((Get-ChildItem -Path $troot -Directory -ErrorAction SilentlyContinue) | Where-Object { $_.Name -like $Filter })
$all | ForEach-Object { $_.Name } | Sort-Object | ForEach-Object { $_ } | Out-File -FilePath $Out -Encoding ascii
$rt = @($all | Where-Object { $_.Name -like "TestRoundTrip*" })
"round=$Round temp_root=$troot filter=$Filter count=$($all.Count) positive_control_all_temp_dirs=$((Get-ChildItem -Path $troot -Directory -ErrorAction SilentlyContinue).Count) stderr_channel=separate_file"
foreach ($d in ($rt | Sort-Object CreationTime)) {
  "  roundtrip_name=$($d.Name) created=$($d.CreationTime.ToString('yyyy-MM-dd HH:mm:ss.fff'))"
}
"roster_done=1"
