param(
    [Parameter(Mandatory = $true)][string]$Tag
)

$ErrorActionPreference = 'Continue'
$repo = 'D:\work\workspace\projects plans\Wisp'
$logs = Join-Path $repo '.scratch\wisp\probes\307\r1\logs'
$aliasLong = 'C:\Users\swq\tmp\riskattrib1\RunnerTempAreaForAttribution1'
$aliasShort = 'C:\Users\swq\tmp\RISKAT~1\RUNNER~1'

Set-Location $repo

function Invoke-GoTest {
    param([string]$Name, [string]$RunFilter, [string]$Table)
    $file = Join-Path $logs ("{0}-{1}-{2}.txt" -f $Name, $Table, $Tag)
    if ($Table -eq 'aliased') {
        $env:TMPDIR = $aliasShort
        $env:TEMP = $aliasShort
        $env:TMP = $aliasShort
    }
    else {
        $env:TMPDIR = ''
        $env:TEMP = $env:ORIGINAL_TEMP
        $env:TMP = $env:ORIGINAL_TMP
    }
    $envFile = Join-Path $logs ("env-{0}-{1}-{2}.txt" -f $Name, $Table, $Tag)
    "table=$Table TMPDIR=[$env:TMPDIR] TEMP=[$env:TEMP] TMP=[$env:TMP]" | Set-Content -Encoding utf8 $envFile
    $sw = [System.Diagnostics.Stopwatch]::StartNew()
    if ($RunFilter -eq '') {
        & go test ./internal/risk/ -count=1 -v 2>&1 | Out-File -Encoding utf8 $file
    }
    else {
        & go test ./internal/risk/ -count=1 -v -run $RunFilter 2>&1 | Out-File -Encoding utf8 $file
    }
    $rc = $LASTEXITCODE
    $sw.Stop()
    Add-Content -Encoding utf8 $file ("---- rc={0} total_seconds={1:N2} table={2} name={3} tag={4}" -f $rc, $sw.Elapsed.TotalSeconds, $Table, $Name, $Tag)
    $p = (Select-String -Path $file -Pattern '^\s*--- PASS' -AllMatches | Measure-Object).Count
    $f = (Select-String -Path $file -Pattern '^\s*--- FAIL' -AllMatches | Measure-Object).Count
    $k = (Select-String -Path $file -Pattern '^\s*--- SKIP' -AllMatches | Measure-Object).Count
    $r = (Select-String -Path $file -Pattern '^\s*=== RUN' -AllMatches | Measure-Object).Count
    $line = ("COUNTS name={0} table={1} tag={2} rc={3} RUN={4} PASS={5} FAIL={6} SKIP={7} seconds={8:N2}" -f $Name, $Table, $Tag, $rc, $r, $p, $f, $k, $sw.Elapsed.TotalSeconds)
    Add-Content -Encoding utf8 (Join-Path $logs ('counts-summary.txt')) $line
    Write-Output $line
}

$env:ORIGINAL_TEMP = $env:TEMP
$env:ORIGINAL_TMP = $env:TMP

Add-Content -Encoding utf8 (Join-Path $logs 'counts-summary.txt') ("==== begin tag={0} at {1}" -f $Tag, (Get-Date -Format 'yyyy-MM-dd HH:mm:ss'))

Invoke-GoTest -Name 'full' -RunFilter '' -Table 'default'
Invoke-GoTest -Name 'full2' -RunFilter '' -Table 'default'
Invoke-GoTest -Name 'full' -RunFilter '' -Table 'aliased'
Invoke-GoTest -Name 'full2' -RunFilter '' -Table 'aliased'
Invoke-GoTest -Name 'syncrun' -RunFilter 'TestSync' -Table 'default'
Invoke-GoTest -Name 'syncrun' -RunFilter 'TestSync' -Table 'aliased'

Add-Content -Encoding utf8 (Join-Path $logs 'counts-summary.txt') ("==== end tag={0} at {1}" -f $Tag, (Get-Date -Format 'yyyy-MM-dd HH:mm:ss'))
