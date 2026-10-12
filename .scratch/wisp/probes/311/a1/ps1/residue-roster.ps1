# Probe 311-a1: roster the residue inside a dedicated test-temp root (read-only).
# Usage: powershell -NoProfile -File residue-roster.ps1 <root> [<root> ...]
$ErrorActionPreference = 'Continue'
Write-Output ('now=' + (Get-Date -Format 'yyyy-MM-dd HH:mm:ss'))
foreach ($r in $args) {
    Write-Output ('===== ROOT=' + $r)
    if (-not (Test-Path -LiteralPath $r)) { Write-Output 'root_exists=NO'; continue }
    $dirs = @(Get-ChildItem -LiteralPath $r -Directory -Recurse -Force -ErrorAction SilentlyContinue)
    $top = @(Get-ChildItem -LiteralPath $r -Directory -Force -ErrorAction SilentlyContinue)
    $allEntries = @(Get-ChildItem -LiteralPath $r -Force -ErrorAction SilentlyContinue)
    Write-Output ('count_top_Directory=' + $top.Count)
    Write-Output ('count_recursive_Directory=' + $dirs.Count)
    Write-Output ('count_top_all_entries=' + $allEntries.Count)
    foreach ($d in $dirs) {
        Write-Output ('DIR|depth=' + (($d.FullName -split '\\').Count - ($r -split '[\\/]').Count) + '|creation=' + $d.CreationTime.ToString('yyyy-MM-dd HH:mm:ss.fff') + '|' + $d.FullName)
    }
}
Write-Output ('ps_errors=' + ($Error.Count))
