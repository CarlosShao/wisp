# Retry the "sample genuinely over budget" reading on the TRACKED bytes until the
# sampling-validity precheck accepts the machine (ticket 263 AC#3).
#
# Why a retry loop is needed and is NOT a weakening of anything: this machine is
# the self-hosted runner, another leg is compiling, and ticket 134 AC#4's precheck
# correctly refuses to sample over 50% cpu. A refusal is the gate doing its job,
# so the harness waits for a quiet window instead of pretending the reading exists.
$ErrorActionPreference = 'Continue'
$here = $PSScriptRoot
$tracked = Join-Path (Split-Path -Parent (Split-Path -Parent (Split-Path -Parent (Split-Path -Parent (Split-Path -Parent $here))))) 'scripts\slo-check.ps1'
$exe = Join-Path $here 'bin\probe-gui.exe'
$log = Join-Path $here 'tracked-fail-retry.txt'
Set-Content -LiteralPath $log -Value "retry harness: tracked=$tracked" -Encoding UTF8
$env:PROBE263_VERDICT = 'fail'
$env:WISP_ENV = 'test'
for ($attempt = 1; $attempt -le 6; $attempt++) {
    $outDir = Join-Path $here ("out-retry-{0}" -f $attempt)
    # ORACLE, after a first version that got this wrong: the gate prints with
    # Write-Host, which does NOT enter the pipeline, so the first version of
    # this loop looked for "precheck ok" in a captured string that was always
    # empty and reported reached-sampling=False on a run that had in fact
    # sampled and reddened (attempts 2 and 4 in retry-wrapper.log). The report
    # file is the honest signal: it exists only when a conclusion was produced.
    $console = Join-Path $here ("retry-console-{0}.txt" -f $attempt)
    & $tracked -Subset smoke -SecondsPerState 1 -WispExe $exe -OutDir $outDir 6>&1 2>&1 | Set-Content -LiteralPath $console -Encoding UTF8
    $rc = 0
    if ($null -ne $LASTEXITCODE) { $rc = $LASTEXITCODE }
    $report = Join-Path $outDir 'slo-report.json'
    $reached = Test-Path -LiteralPath $report
    Add-Content -LiteralPath $log -Value ("attempt {0}: rc={1} report-written(=reached-sampling)={2} console={3}" -f $attempt, $rc, $reached, $console) -Encoding UTF8
    if ($reached) { break }
    Start-Sleep -Seconds 15
}
Write-Host ("retry harness done, log at {0}" -f $log)
