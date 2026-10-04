param(
    [Parameter(Mandatory)][string]$Exe,
    [string]$ShellExecute = 'false'
)
# Prototype of the helper that goes into scripts/slo-check.ps1 (ticket 263):
# start an external program, WAIT for it, and take its real exit code without
# ever reading the automatic exit-code variable. Probes both UseShellExecute
# shapes and whether the child's own lines still reach this caller's streams.
$ErrorActionPreference = 'Stop'
Set-StrictMode -Version 2.0

Write-Host ("helper-probe: launching {0} (UseShellExecute={1})" -f $Exe, $ShellExecute)
$psi = New-Object System.Diagnostics.ProcessStartInfo
$psi.FileName = $Exe
$psi.Arguments = '2 3'
$psi.UseShellExecute = ($ShellExecute -eq 'true')
$psi.WorkingDirectory = Split-Path -Parent $Exe
$proc = [System.Diagnostics.Process]::Start($psi)
$proc.WaitForExit()
$code = [int]$proc.ExitCode
$proc.Dispose()
Write-Host ("helper-probe: waited, exit={0}" -f $code)
