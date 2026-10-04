param(
    [Parameter(Mandatory)][string]$Child
)
# Question behind scripts/slo-check.ps1 line 105-106 (ticket 263, latent site):
# does calling a .ps1 with & leave an exit code readable in THIS scope?
# Three child shapes are probed by repro-child-*.ps1:
#   1. child that only calls `exit 3`            (build.ps1's Fail path)
#   2. child that only runs a native command      (build.ps1's go build path)
#   3. child that runs a native command then exits 0
$ErrorActionPreference = 'Stop'
Set-StrictMode -Version 2.0

Write-Host ("repro-parent: calling {0}" -f $Child)
& $Child
try {
    $code = $LASTEXITCODE
    Write-Host ("repro-parent: read exit={0}" -f $code)
} catch {
    Write-Host ("repro-parent: THROW - {0}" -f $_.Exception.Message)
}
