param(
    [string]$WispExe = '',
    [string]$Subset = 'smoke',
    [int]$SecondsPerState = 1,
    [string]$Only = ''
)
# Ticket 263 AC#3/AC#4 harness: build mutated COPIES of the tracked
# scripts/slo-check.ps1, run each, and print what colour it comes out.
# Nothing here edits a tracked file - the mutations live only under
# mutations/, because this is a shared working tree and a transient red inside
# the real gate could be read by another leg as a regression.
#
# Every mutation asserts its anchor was found before writing, so a mutation
# that silently no-ops cannot report a fake reading.
$ErrorActionPreference = 'Stop'
Set-StrictMode -Version 2.0

$here = $PSScriptRoot
$root = $here
for ($i = 0; $i -lt 5; $i++) { $root = Split-Path -Parent $root }
$src = Join-Path $root 'scripts\slo-check.ps1'
if (-not (Test-Path -LiteralPath $src)) { throw "cannot locate the tracked script at $src" }
if (-not $WispExe) { $WispExe = Join-Path $here 'bin\probe-gui.exe' }
$mutDir = Join-Path $here 'mutations'
New-Item -ItemType Directory -Force -Path $mutDir | Out-Null

$text = Get-Content -LiteralPath $src -Raw
Write-Host ("harness: source={0} bytes={1}" -f $src, $text.Length)

function Make-Mutation {
    param([string]$Name, [string]$Anchor, [string]$Replacement)
    if (-not $text.Contains($Anchor)) {
        throw ("mutation {0}: anchor not found - the mutation did NOT land, refusing to report a reading" -f $Name)
    }
    $out = $text.Replace($Anchor, $Replacement)
    if ($out -eq $text) {
        throw ("mutation {0}: replacement produced no change - refusing to report a reading" -f $Name)
    }
    $path = Join-Path $mutDir ("slo-check-{0}.ps1" -f $Name)
    Set-Content -LiteralPath $path -Value $out -Encoding UTF8
    Write-Host ("mutation {0}: landed, copy written to {1}" -f $Name, $path)
    return $path
}

# Anchors inside the tracked script (verified against it before writing this).
$loopAnchor = '    $run = Invoke-ExternalProgram -FilePath $WispExe -ArgumentList @('
$waitAnchor = '        $proc.WaitForExit()'
$codeAnchor = "        return [pscustomobject]@{ ok = `$true; code = [int]`$proc.ExitCode; why = '' }"
$passAnchor = '            $pass = [bool]$reportJson.pass'
$takeAnchor = '    $code = $run.code'

$cases = [ordered]@{
    'clean-copy'         = @{ anchor = $loopAnchor; repl = $loopAnchor }  # no-op copy, the control
    'plant-bare'         = @{ anchor = $loopAnchor; repl = "    `$stale = `$LASTEXITCODE`n" + $loopAnchor }
    'plant-braced'       = @{ anchor = $loopAnchor; repl = "    `$stale = `${LASTEXITCODE}`n" + $loopAnchor }
    'plant-global'        = @{ anchor = $loopAnchor; repl = "    `$stale = `$global:LASTEXITCODE`n" + $loopAnchor }
    'plant-global-used'  = @{ anchor = $takeAnchor; repl = '    $code = $global:LASTEXITCODE' }
    # The remaining bypass after the nail learned the scope-qualified spelling:
    # a runtime lookup of the same variable is invisible to a text scan.
    'plant-getvariable'  = @{ anchor = $loopAnchor; repl = "    `$stale = (Get-Variable -Name 'LASTEXITCODE' -ValueOnly)`n" + $loopAnchor }
    'drop-wait'          = @{ anchor = $waitAnchor; repl = '        # drop-wait: the wait was removed by mutations/drop-wait' }
    'hardcode-zero'      = @{ anchor = $codeAnchor; repl = "        return [pscustomobject]@{ ok = `$true; code = 0; why = '' }" }
    # AC#3's other half: the gate must STILL redden when a state is judged
    # against a must-fail shape. This mutates a JUDGMENT line, not the
    # instrument, and the clean fixture (which otherwise passes) stays in place.
    'force-state-fail'   = @{ anchor = $passAnchor; repl = '            $pass = $false' }
}

foreach ($name in $cases.Keys) {
    if ($Only -and $Only -ne $name) { continue }
    $case = $cases[$name]
    if ($name -eq 'clean-copy') {
        $path = Join-Path $mutDir 'slo-check-clean-copy.ps1'
        Set-Content -LiteralPath $path -Value $text -Encoding UTF8
        Write-Host ("mutation {0}: byte-identical copy written to {1}" -f $name, $path)
    } else {
        $path = Make-Mutation -Name $name -Anchor $case.anchor -Replacement $case.repl
    }
    $outDir = Join-Path $mutDir ("out-{0}" -f $name)
    Write-Host ("=== RUN {0} (subset={1} secondsPerState={2}) ===" -f $name, $Subset, $SecondsPerState)
    $prev = $ErrorActionPreference
    $ErrorActionPreference = 'Continue'
    # Pre-assigned on purpose: this harness may run its FIRST copy in a session
    # where no exit code exists yet, and reading the automatic variable then is
    # exactly the ticket-263 death. Measured: with -Only plant-global-used the
    # harness's own read threw VariableIsUndefined on the first (and only) call.
    $rc = -1
    & $path -Subset $Subset -SecondsPerState $SecondsPerState -WispExe $WispExe -OutDir $outDir 2>&1 | ForEach-Object { Write-Host ("  | {0}" -f $_) }
    if ($null -ne $LASTEXITCODE) { $rc = $LASTEXITCODE }
    $ErrorActionPreference = $prev
    Write-Host ("=== RUN {0} rc={1} ===" -f $name, $rc)
}
