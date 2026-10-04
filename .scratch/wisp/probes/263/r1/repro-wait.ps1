param(
    [Parameter(Mandatory)][string]$Exe,
    [int]$SleepSeconds = 2,
    [int]$ExitCode = 3
)
# Minimal stand-in for scripts/slo-check.ps1 lines 324-326/337 (ticket 263).
# Same three statements in the same order, same host settings:
#   announce the state -> launch the external program with & -> read the exit code.
$ErrorActionPreference = 'Stop'
Set-StrictMode -Version 2.0

$t0 = Get-Date
Write-Host ("repro: sampling state Sleeping for {0}s" -f $SleepSeconds)
& $Exe $SleepSeconds $ExitCode
$code = $LASTEXITCODE
$elapsedMs = [int]((Get-Date) - $t0).TotalMilliseconds
Write-Host ("repro: state exit={0} elapsed-ms={1}" -f $code, $elapsedMs)
