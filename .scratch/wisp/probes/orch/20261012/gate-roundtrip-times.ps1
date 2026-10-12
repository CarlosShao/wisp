$names = Get-ChildItem -Path $env:TEMP -Directory -ErrorAction SilentlyContinue | Where-Object { $_.Name -like "TestRoundTrip*" }
Write-Output ("rows=" + @($names).Count)
foreach ($d in $names) {
  $wrote = ""
  $f = Get-ChildItem -Path $d.FullName -File -ErrorAction SilentlyContinue | Select-Object -First 1
  if ($f) { $wrote = $f.LastWriteTime.ToString("yyyy-MM-dd HH:mm:ss") }
  Write-Output ("name=" + $d.Name + " created=" + $d.CreationTime.ToString("yyyy-MM-dd HH:mm:ss") + " firstfile_write=" + $wrote + " files=" + (@(Get-ChildItem -Path $d.FullName -ErrorAction SilentlyContinue).Count))
}
