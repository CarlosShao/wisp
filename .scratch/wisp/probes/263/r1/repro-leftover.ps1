$ErrorActionPreference = 'Stop'
Set-StrictMode -Version 2.0
# Does the ticket-263 shape leave the child running behind it? Uses the GUI
# subsystem probe with a 6s lifetime, the broken `& exe` + immediate read, and
# then asks the OS whether the process is still alive at the moment the script
# is already dead.
$exe = Join-Path $PSScriptRoot 'bin\probe-gui.exe'
try {
    & $exe 6 3
    $code = $LASTEXITCODE
    Write-Host ("leftover-probe: no throw, read exit={0}" -f $code)
} catch {
    Write-Host ('leftover-probe: THROW - ' + $_.Exception.Message)
}
# The process name is the EXE base name (probe-gui), not the Go module name;
# the first run of this probe filtered on 'probe263' and reported 0 alive,
# which was a broken ruler, not a fact. Fixed 2026-10-04 by $exe's own base name.
$procName = [IO.Path]::GetFileNameWithoutExtension($exe)
$alive = @(Get-Process -Name $procName -ErrorAction SilentlyContinue)
Write-Host ("leftover-probe: {0} processes still alive right after the death = {1}" -f $procName, $alive.Count)
foreach ($p in $alive) { Write-Host ("leftover-probe: alive pid={0}" -f $p.Id) }
foreach ($p in $alive) { $p.WaitForExit() }
Write-Host ("leftover-probe: alive after waiting them out = {0}" -f @(Get-Process -Name $procName -ErrorAction SilentlyContinue).Count)
