param(
    [Parameter(Mandatory = $true)][string]$Mutation,
    [string]$Run = 'TestSync',
    [string]$Table = 'aliased'
)

$ErrorActionPreference = 'Stop'
$tree = 'C:\Users\swq\tmp\wisp-307-r1-tree-main'
$pristine = 'C:\Users\swq\tmp\wisp-307-r1-pristine-syncdirs.go'
$target = Join-Path $tree 'internal\risk\syncdirs.go'
$logs = 'D:\work\workspace\projects plans\Wisp\.scratch\wisp\probes\307\r1\logs'
$aliasShort = 'C:\Users\swq\tmp\RISKAT~1\RUNNER~1'

# literal anchors (exact substrings of the committed file)
$lineRoot = 'if aUnder(candForms, anchorForms(e.canon)) {'
$lineHome = 'if !s.complete && s.home != "" && aUnder(candForms, anchorForms(s.home)) {'
$lineForms = "`tcandForms := formsOf(cand)"

$src = [System.IO.File]::ReadAllText($pristine)

switch ($Mutation) {
    'pristine' {
        $new = $src
    }
    'a' {
        $new = $src.Replace($lineHome, 'if !s.complete && s.home != "" && isUnder(cand, s.home) {')
    }
    'c' {
        $new = $src.Replace($lineRoot, 'if isUnder(cand, e.canon) {')
    }
    'fullrevert' {
        $new = $src.Replace($lineRoot, 'if isUnder(cand, e.canon) {').
            Replace($lineHome, 'if !s.complete && s.home != "" && isUnder(cand, s.home) {').
            Replace($lineForms, "`t_ = formsOf(cand)")
    }
    'd' {
        $new = $src.Replace($lineRoot, 'if aUnder(candForms, anchorForms(e.canon)) || true {').
            Replace($lineHome, 'if !s.complete && s.home != "" && aUnder(candForms, anchorForms(s.home)) || true {')
    }
    default {
        throw "unknown mutation $Mutation"
    }
}

$changed = ($new -ne $src)
if ($Mutation -ne 'pristine' -and -not $changed) {
    throw "mutation $Mutation changed nothing - the literal anchors did not match"
}

[System.IO.File]::WriteAllText($target, $new)

$outFile = Join-Path $logs ("mut-{0}-{1}.txt" -f $Mutation, $Table)
$head = ("MUTATION={0} TABLE={1} literal-anchor-roundtrip-changed={2}`r`n" -f $Mutation, $Table, $changed)
Set-Content -NoNewline -Encoding utf8 $outFile $head

$origTemp = $env:TEMP
$origTmp = $env:TMP
if ($Table -eq 'aliased') {
    $env:TMPDIR = $aliasShort
    $env:TEMP = $aliasShort
    $env:TMP = $aliasShort
}
else {
    $env:TMPDIR = ''
    $env:TEMP = $origTemp
    $env:TMP = $origTmp
}
"TABLE=$Table TMPDIR=[$env:TMPDIR] TEMP=[$env:TEMP] TMP=[$env:TMP]" | Add-Content -Encoding utf8 $outFile

Set-Location $tree
& go test ./internal/risk/ -count=1 -v -run $Run 2>&1 | Out-File -Append -Encoding utf8 $outFile
$rc = $LASTEXITCODE

Copy-Item -Force $pristine $target

$shaTree = (Get-FileHash -Algorithm SHA256 $target).Hash
$shaPris = (Get-FileHash -Algorithm SHA256 $pristine).Hash
$p = (Select-String -Path $outFile -Pattern '^\s*--- PASS' | Measure-Object).Count
$f = (Select-String -Path $outFile -Pattern '^\s*--- FAIL' | Measure-Object).Count
$k = (Select-String -Path $outFile -Pattern '^\s*--- SKIP' | Measure-Object).Count
$tail = ("---- mutation={0} rc={1} PASS={2} FAIL={3} SKIP={4} run={5} table=aliased`r sha256_tree_after_restore={6}`r sha256_pristine={7}`r" -f $Mutation, $rc, $p, $f, $k, $Run, $shaTree, $shaPris)
Add-Content -NoNewline -Encoding utf8 $outFile $tail
Write-Output ("MUTATION name={0} rc={1} PASS={2} FAIL={3} SKIP={4}" -f $Mutation, $rc, $p, $f, $k)
Write-Output ("RESTORED sha_tree={0} sha_pristine={1} equal={2}" -f $shaTree, $shaPris, ($shaTree -eq $shaPris))
