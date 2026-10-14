$ErrorActionPreference = "Continue"
$w = @{LogName="Application"; StartTime=(Get-Date "2026-10-14 15:20:00"); EndTime=(Get-Date "2026-10-14 16:00:00")}
Get-WinEvent -FilterHashtable $w | Select-Object TimeCreated,ProviderName,Id,LevelDisplayName | Sort-Object TimeCreated | Format-Table -AutoSize | Out-String -Width 200
