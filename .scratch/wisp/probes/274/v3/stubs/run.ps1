# 274-v3 battery runner (ASCII only for PS 5.1).
$ErrorActionPreference = 'Continue'
$repo = 'D:\work\workspace\projects plans\Wisp'
$V3 = Join-Path $repo '.scratch\wisp\probes\274\v3'
$LOGS = Join-Path $V3 'logs'
$base = 'D:\tmp\wisp274v3'
$inner = Join-Path $base 'inner'
$psExe = (Get-Command powershell.exe).Source

function Repo-Md5 {
    $files = @('frontend\dist\.gitkeep', 'frontend\dist\index.html',
        'frontend\dist\assets\index-BRKj5OIJ.css', 'frontend\dist\assets\index-BVKlegVD.js',
        'build\wisp.exe', 'scripts\build.ps1')
    foreach ($f in $files) {
        $p = Join-Path $repo $f
        (Get-FileHash -Algorithm MD5 -LiteralPath $p).Hash.ToLower() + '  ' + $f
    }
}

# ---- 05: does the :238-242 shape (git rev-parse 2>$null) throw under EAP=Stop outside a repo?
$testSrc = @'
$ErrorActionPreference = 'Stop'
Set-StrictMode -Version 2.0
$g = Get-Command 'git.exe' -ErrorAction SilentlyContinue
if ($null -eq $g) { 'NO_GIT_ON_PATH'; exit }
try {
    $s = (& git rev-parse --short HEAD 2>$null)
    "NO_THROW git=$($g.Source) LASTEXITCODE=$LASTEXITCODE short=[$s]"
} catch {
    "THREW git=$($g.Source) EXC=$($_.Exception.GetType().FullName) MSG=$($_.Exception.Message)"
}
'@
$testPs = Join-Path $inner 'gitthrow-test.ps1'
Set-Content -LiteralPath $testPs -Value $testSrc -Encoding ASCII
$throwLines = & $psExe -NoProfile -ExecutionPolicy Bypass -Command "Set-Location -LiteralPath '$inner'; & '$testPs'" 2>&1 | ForEach-Object { "$_" }
$throwText = $throwLines -join "`n"
$gitThrows = ($throwText -match 'THREW')
$log05 = @(
    '=== 05 git-throw probe (replicates build.ps1 :238-242 shape in a dir without .git, EAP=Stop) ===',
    "command: & $psExe -NoProfile -ExecutionPolicy Bypass -Command ""Set-Location -LiteralPath '$inner'; & '$testPs'""",
    "date: $(Get-Date -Format o)",
    $throwText,
    "gitThrows=$gitThrows  (if true: per-tree fresh 'git init' + empty commit below neutralizes :240; repo .git is never mirrored)"
)
Set-Content -LiteralPath (Join-Path $LOGS '05-git-throw-test.txt') -Value $log05 -Encoding UTF8
Write-Host ("274-v3 run.ps1: gitThrows=" + $gitThrows)

$cases = @(
    @{ name = '01-spaced-stub-nowrite';  file = '10-case1-spaced-stub-nowrite.txt';  tree = 'tree-a'; stub = 'stub-nowrite'; expect = 'rc=1, Named cause: unprovenanced page bytes (spaced name intact in BOTH array rosters)' },
    @{ name = '02-reverted-join-split';  file = '11-case2-reverted-join-split.txt';  tree = 'tree-b'; stub = 'stub-nowrite'; expect = 'rc=0 false green again (tautology check: teeth live exactly at :209-210)' },
    @{ name = '03-nospace-stub-nowrite'; file = '12-case3-nospace-stub-nowrite.txt'; tree = 'tree-c'; stub = 'stub-nowrite'; expect = 'rc=1 same named cause (teeth not dulled)' },
    @{ name = '04-spaced-written-green'; file = '13-case4-spaced-written-green.txt'; tree = 'tree-c'; stub = 'stub-write';    expect = 'rc=0 and intact spaced name printed in fresh artifact list' }
)

$summary = @()
foreach ($c in $cases) {
    $treeDir = Join-Path $base $c.tree
    $bp = Join-Path $treeDir 'scripts\build.ps1'
    $stubDir = Join-Path $base $c.stub
    if ($gitThrows -and -not (Test-Path -LiteralPath (Join-Path $treeDir '.git'))) {
        & git init -q $treeDir | Out-Null
        & git -C $treeDir -c user.email=probe@local -c user.name=probe commit -q --allow-empty -m '274-v3 bench anchor' | Out-Null
    }
    $env:PATH = $stubDir + ';' + $env:PATH
    $pre = @(
        ('=== CASE ' + $c.name + ' ==='),
        ('command: & ' + $psExe + ' -NoProfile -ExecutionPolicy Bypass -File "' + $bp + '" -Env dev'),
        ('expect: ' + $c.expect),
        ('date: ' + (Get-Date -Format o)),
        ('tree: ' + $treeDir + '  stub: ' + $stubDir + '  treeHasFreshGitRepo=' + (Test-Path -LiteralPath (Join-Path $treeDir '.git'))),
        ('npm resolved -> ' + (Get-Command npm.cmd -ErrorAction SilentlyContinue).Source),
        ('node resolved -> ' + (Get-Command node.exe -ErrorAction SilentlyContinue).Source),
        ('go resolved -> ' + (Get-Command go.exe -ErrorAction SilentlyContinue).Source),
        ('gcc resolved -> ' + (Get-Command gcc.exe -ErrorAction SilentlyContinue).Source),
        ('git resolved -> ' + (Get-Command git.exe -ErrorAction SilentlyContinue).Source),
        ('tree dist BEFORE: ' + ((Get-ChildItem -LiteralPath (Join-Path $treeDir 'frontend\dist') -Recurse -File -Force | ForEach-Object { $_.Name }) -join ' | '))
    )
    $out = & $psExe -NoProfile -ExecutionPolicy Bypass -File $bp -Env dev 2>&1 | ForEach-Object { "$_" }
    $rc = $LASTEXITCODE
    $post = @(
        ('CASE_' + $c.name + '_PROCESS_EXITCODE=' + $rc),
        ('tree dist AFTER: ' + ((Get-ChildItem -LiteralPath (Join-Path $treeDir 'frontend\dist') -Recurse -File -Force | ForEach-Object { $_.Name }) -join ' | ')),
        '--- repo (in-repo) md5 AFTER case (must be untouched all through) ---'
    ) + (Repo-Md5)
    Set-Content -LiteralPath (Join-Path $LOGS $c.file) -Value ($pre + $out + $post) -Encoding UTF8
    Write-Host ('274-v3 run.ps1: case ' + $c.name + ' rc=' + $rc)
    $summary += ($c.name + '=' + $rc)
    $env:PATH = (($env:PATH -split ';') | Where-Object { $_ -notlike ($base + '\stub-*') }) -join ';'
}
Write-Host ('274-v3 run.ps1: SUMMARY ' + ($summary -join ' '))
