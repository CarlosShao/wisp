function Get-V1ExitCode {
    param([string]$Exe, [string]$State, [double]$Seconds, [string]$Out)
    & $Exe slo -state $State -seconds $Seconds -interval-ms 250 -out $Out
    return $LASTEXITCODE
}
