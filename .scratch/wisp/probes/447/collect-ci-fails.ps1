# collect-ci-fails.ps1 -- ticket 447 AC#1 instrument.
#
# Reads the LAST N test-windows step-10 sections straight from GitHub and prints
# one machine-readable line per run:
#
#   RUN=<id> SHA=<sha> CONCLUSION=<x> STEP10=<success|failure|missing>
#     FOURNUM=RUN=846/PASS=603/FAIL=1/SKIP=19
#     FAIL=<name>;FAIL=<name>            (empty when green)
#     SKIP=<name>;SKIP=<name>            (the named skips, in log order)
#
# step 10 of the test-windows job is "cmd/wisp CLI tests (needs the sherpa DLLs
# staged above, ticket 111 AC#4)" -- confirmed from the job JSON, not guessed
# from a position.
#
# The raw job log is NOT kept (3-4 MB per run): the four numbers and the FAIL /
# SKIP name lists are the measurement, and every one of them is re-derivable by
# re-running this script. Usage: pwsh -File collect-ci-fails.ps1 [-Limit 30] [-Need 20]
param(
    [int]$Limit = 30,
    [int]$Need = 20
)
$ErrorActionPreference = 'Stop'
$gh = 'C:\Program Files\GitHub CLI\gh.exe'
$root = 'D:\work\workspace\projects plans\Wisp'
Set-Location $root
$outDir = Join-Path $root '.scratch\wisp\probes\447\logs'
New-Item -ItemType Directory -Force $outDir | Out-Null

$runs = & $gh run list --workflow ci --limit $Limit --json databaseId,headSha,conclusion,createdAt | ConvertFrom-Json
$done = 0
foreach ($r in $runs) {
    if ($done -ge $Need) { break }
    $id = $r.databaseId
    $jobs = (& $gh run view $id --json jobs) | ConvertFrom-Json
    $job = $jobs.jobs | Where-Object { $_.name -eq 'test-windows' } | Select-Object -First 1
    if (-not $job) {
        "RUN=$id SHA=$($r.headSha.Substring(0,9)) CONCLUSION=$($r.conclusion) STEP10=no-test-windows-job"
        continue
    }
    $step = $job.steps | Where-Object { $_.number -eq 10 } | Select-Object -First 1
    $stepName = if ($step) { $step.name } else { '' }
    if (-not $step -or $step.status -ne 'completed') {
        "RUN=$id SHA=$($r.headSha.Substring(0,9)) CONCLUSION=$($r.conclusion) STEP10=incomplete($($step.status))"
        continue
    }
    # Stream the job log once; keep only the step-10 lines that carry a verdict.
    $lines = & $gh run view $id --job $job.databaseId --log 2>&1
    $fails = New-Object System.Collections.Generic.List[string]
    $skips = New-Object System.Collections.Generic.List[string]
    $four = ''
    foreach ($l in $lines) {
        # gh's job log is  <job>\t<step>\t<timestamp> <content>
        # MEASURED: only TWO tabs per line (7,072 of 7,073 step-10 lines of run
        # 43584365596 carry 3 fields, one carries 4) - the timestamp and the
        # content share the last field, unlike the 4-field shape one assumes.
        # Reading field 3 as the content made every verdict line look empty
        # (first draft of this script read FOURNUM=UNPARSED for green runs).
        if ($l -isnot [string]) { continue }
        # 3 fields max: <job> \t <step> \t <timestamp> <content>; the content may
        # itself carry a tab (one such line measured in 7,073), so take the whole
        # tail instead of a single field.
        $p = $l -split "`t", 3
        if ($p.Count -lt 3) { continue }
        # 'UNKNOWN STEP' is what older job logs say in the step column for the
        # SAME content (measured: run 43578228598's cmd/wisp four-numbers line
        # arrives as `test-windows\tUNKNOWN STEP\t...`). Refusing that token is
        # how the first draft reported FOURNUM=UNPARSED on runs that were green.
        # BUT the very same bucket carries the LATER "Portable windows tests"
        # step too (its own four numbers are RUN=1275/PASS=833 and it prints the
        # same --- SKIP names), and reading it wholesale reported the wrong
        # package's four numbers for the five older runs. The step's own
        # [startedAt, completedAt] window is the discriminator: a line stamped
        # inside step 10's window is step 10's line, in either spelling.
        $inStep = ($p[1] -eq $stepName)
        if (-not $inStep -and $p[1] -ne 'UNKNOWN STEP') { continue }
        $c = $p[2]
        # strip the optional leading BOM and capture the "<RFC3339> " stamp
        $c = $c -replace '^\uFEFF', ''
        if ($c -match '^(\d{4}-\d\d-\d\dT[\d:.]+Z)\s') {
            $ts = [datetime]::Parse($Matches[1]).ToUniversalTime()
            $c = $c.Substring($Matches[0].Length)
            if (-not $inStep) {
                $inStep = ($ts -ge $step.startedAt.AddSeconds(-30) -and $ts -le $step.completedAt.AddSeconds(30))
            }
        }
        if (-not $inStep) { continue }
        # 10, not 9: "--- FAIL: " is ten characters and index 9 is the SPACE, so
        # splitting at 9 yielded the empty string as the name (first draft
        # printed FAIL= with no name on runs that really had one).
        if ($c.StartsWith('--- FAIL: ')) {
            $n = ($c.Substring(10) -split ' ')[0]
            if (-not $fails.Contains($n)) { $fails.Add($n) }
        }
        elseif ($c.StartsWith('--- SKIP: ')) {
            $n = ($c.Substring(10) -split ' ')[0]
            if (-not $skips.Contains($n)) { $skips.Add($n) }
        }
        elseif ($c -like 'portable-tests.sh: four numbers*') { $four = $c }
    }
    $fourBits = ''
    if ($four -match 'RUN=(\d+)\s+--- PASS=(\d+)\s+--- FAIL=(\d+)\s+--- SKIP=(\d+)') {
        $fourBits = "RUN=$($Matches[1])/PASS=$($Matches[2])/FAIL=$($Matches[3])/SKIP=$($Matches[4])"
    } else { $fourBits = 'UNPARSED' }
    $failBits = ($fails | ForEach-Object { "FAIL=$_" }) -join ';'
    $skipBits = ($skips | ForEach-Object { "SKIP=$_" }) -join ';'
    $line = "RUN=$id SHA=$($r.headSha.Substring(0,9)) CONCLUSION=$($r.conclusion) STEP10=$($step.conclusion) FOURNUM=$fourBits`n  $failBits`n  $skipBits"
    $line
    Set-Content -Encoding utf8 (Join-Path $outDir "$id.step10.txt") $line
    $done++
}

