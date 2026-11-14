$exe = Join-Path $PWD 'build/wisp.exe'
if (-not (Test-Path -LiteralPath $exe)) {
  Write-Host "the build step above reported success but $exe does not exist - nothing was measured, so this step does not pass"
  exit 1
}
# `& $exe panel-assets -check` is NOT the shape used here. build\wisp.exe
# links with -H=windowsgui (scripts/build.ps1:244-257), and PowerShell does
# not wait for a GUI-subsystem child: measured on this bench 2026-11-14,
# `& $exe panel-assets -check *> file` returned after 20 ms with an EMPTY
# file and $LASTEXITCODE still reading 99 from the statement before it.
# The measured shape is the one build.ps1's own smoke test uses
# (scripts/build.ps1:307-361, ticket 244-r2): Process.Start with
# UseShellExecute=$false, so the child's stdout still reaches this log,
# then WaitForExit and the real ExitCode.
$psi = New-Object System.Diagnostics.ProcessStartInfo
$psi.FileName = $exe
$psi.Arguments = 'panel-assets -check'
$psi.UseShellExecute = $false
$proc = $null
try {
  $proc = [System.Diagnostics.Process]::Start($psi)
} catch {
  Write-Host "could not start ${exe}: $($_.Exception.Message) - no verdict was produced, so this step does not pass"
  exit 1
}
if ($null -eq $proc) {
  Write-Host "no process object came back for ${exe} - no verdict was produced, so this step does not pass"
  exit 1
}
$code = $null
try {
  $proc.WaitForExit()
  $code = [int]$proc.ExitCode
} catch {
  Write-Host "${exe} started but its exit code could not be read: $($_.Exception.Message) - no verdict was produced, so this step does not pass"
  exit 1
} finally {
  $proc.Dispose()
}
Write-Host "wisp.exe panel-assets -check -> exit code $code"
if ($code -ne 0) {
  Write-Host "the shipping executable cannot serve its own embedded page bundle (panel-assets -check exit code $code); the lines above name the missing piece"
  exit 1
}
