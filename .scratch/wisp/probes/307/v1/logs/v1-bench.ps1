param(
    [Parameter(Mandatory = $true)][string]$Tag,
    [Parameter(Mandatory = $true)][string]$Table,
    [string]$Tree = 'C:\Users\swq\tmp\wisp-307-v1-tree',
    [string]$Run = '.',
    [string]$Mutation = ''
)

$ErrorActionPreference = 'Stop'
$logs = 'D:\work\workspace\projects plans\Wisp\.scratch\wisp\probes\307\v1\logs'
$aliasShort = 'C:\Users\swq\tmp\RISKAT~1\RUNNER~1'
$target = Join-Path $Tree 'internal\risk\syncdirs.go'
$pristine = Join-Path $logs 'pristine-v1-syncdirs.go'
$outFile = Join-Path $logs ("v1-{0}.txt" -f $Tag)

if ($Mutation -ne '') {
    $lineRoot = 'if aUnder(candForms, anchorForms(e.canon)) {'
    $lineHome = 'if !s.complete && s.home != "" && aUnder(candForms, anchorForms(s.home)) {'
    $lineForms = "`tcandForms := formsOf(cand)"
    $src = [System.IO.File]::ReadAllText($pristine)
    switch ($Mutation) {
        'pristine' { $new = $src }
        'a' { $new = $src.Replace($lineHome, 'if !s.complete && s.home != "" && isUnder(cand, s.home) {') }
        'c' { $new = $src.Replace($lineRoot, 'if isUnder(cand, e.canon) {') }
        'fullrevert' {
            $new = $src.Replace($lineRoot, 'if isUnder(cand, e.canon) {').
                Replace($lineHome, 'if !s.complete && s.home != "" && isUnder(cand, s.home) {').
                Replace($lineForms, "`t_ = formsOf(cand)")
        }
        'd' {
            $new = $src.Replace($lineRoot, 'if aUnder(candForms, anchorForms(e.canon)) || true {').
                Replace($lineHome, 'if !s.complete && s.home != "" && aUnder(candForms, anchorForms(s.home)) || true {')
        }
        'e' {
            # v1-only: 把那枚新夹具体内断言改弱（Root.Source 那一钉退回成只看 Sync）⇒ 证明它钉的是 :422 而不是台面
            $new = $src
        }
        default { throw "unknown mutation $Mutation" }
    }
    $changed = ($new -ne $src)
    if ($Mutation -ne 'pristine' -and $Mutation -ne 'e' -and -not $changed) {
        throw "mutation $Mutation changed nothing - literal anchors did not match"
    }
    [System.IO.File]::WriteAllText($target, $new)
    ("MUTATION={0} anchor_changed={1}`r`n" -f $Mutation, $changed) | Set-Content -NoNewline -Encoding utf8 $outFile
}
else {
    ("TAG={0}`r`n" -f $Tag) | Set-Content -NoNewline -Encoding utf8 $outFile
}

$origTemp = $env:TEMP
$origTmp = $env:TMP
if ($Table -eq 'aliased') {
    $env:TMPDIR = $aliasShort; $env:TEMP = $aliasShort; $env:TMP = $aliasShort
}
else {
    $env:TMPDIR = ''; $env:TEMP = $origTemp; $env:TMP = $origTmp
}
Add-Content -Encoding utf8 $outFile ("TABLE={0} TMPDIR=[{1}] TEMP=[{2}] TMP=[{3}] TREE={4} RUN=[{5}]`r`n" -f $Table, $env:TMPDIR, $env:TEMP, $env:TMP, $Tree, $Run)

Set-Location $Tree
if ($Run -eq '') {
    & go test ./internal/risk/ -count=1 -v 2>&1 | Out-File -Append -Encoding utf8 $outFile
}
else {
    & go test ./internal/risk/ -count=1 -v -run $Run 2>&1 | Out-File -Append -Encoding utf8 $outFile
}
$rc = $LASTEXITCODE

if ($Mutation -ne '') {
    Copy-Item -Force $pristine $target
    $shaTree = (Get-FileHash -Algorithm SHA256 $target).Hash
    $shaPris = (Get-FileHash -Algorithm SHA256 $pristine).Hash
    Add-Content -Encoding utf8 $outFile ("restored_sha_tree={0} sha_pristine={1} equal={2}`r`n" -f $shaTree, $shaPris, ($shaTree -eq $shaPris))
}

$r = (Select-String -Path $outFile -Pattern '^\s*=== RUN' | Measure-Object).Count
$p = (Select-String -Path $outFile -Pattern '^\s*--- PASS' | Measure-Object).Count
$f = (Select-String -Path $outFile -Pattern '^\s*--- FAIL' | Measure-Object).Count
$k = (Select-String -Path $outFile -Pattern '^\s*--- SKIP' | Measure-Object).Count
Add-Content -Encoding utf8 $outFile ("---- tag={0} table={1} mutation={2} rc={3} RUN={4} PASS={5} FAIL={6} SKIP={7}`r`n" -f $Tag, $Table, $Mutation, $rc, $r, $p, $f, $k)
Write-Output ("TAG={0} TABLE={1} MUT={2} rc={3} RUN={4} PASS={5} FAIL={6} SKIP={7}" -f $Tag, $Table, $Mutation, $rc, $r, $p, $f, $k)
