#Requires -Version 5.1
<#
.SYNOPSIS
    One-shot build of wisp.exe on Windows (SPEC-11 §2.2). From a clean clone,
    this is the only command needed besides installing the toolchain in
    docs/BUILD.md.

.DESCRIPTION
    1. scripts/fetch-deps.ps1        (verify/complete third_party/, cached runs finish fast)
    2. frontend: built here - npm ci (when node_modules is absent) + npm run
       build in frontend/, then frontend/dist is verified to hold a real page.
       No node/npm on this machine means this script FAILS; the step never skips
       (ticket 274: '//go:embed all:dist' in frontend/embed.go is what carries
       the page bytes into wisp.exe, so a skipped step ships a pageless exe).
    3. go build with CGO_ENABLED=1 + mingw-w64 gcc, real link of the
       sherpa-onnx C API via the official Go bindings
    4. output: build\wisp.exe + onnxruntime.dll + sherpa-onnx DLLs colocated
       in the same directory (SPEC-11 §7.1)
    5. build\SHA256SUMS
    6. smoke test: runs build\wisp.exe doctor

.USAGE
    scripts/build.ps1 [-Env dev|prod]
#>
[CmdletBinding()]
param(
    [ValidateSet('dev', 'prod')]
    [string]$Env = 'dev'
)

$ErrorActionPreference = 'Stop'
Set-StrictMode -Version 2.0
$ProgressPreference = 'SilentlyContinue'

$RepoRoot = Split-Path -Parent $PSScriptRoot
$BuildInfoPkg = 'github.com/CarlosShao/wisp/internal/buildinfo'

function Fail([string]$Message) {
    Write-Host "build.ps1: FATAL: $Message"
    exit 1
}

Write-Host "build.ps1: repo root $RepoRoot (Env=$Env)"

# --- 0. toolchain discovery -------------------------------------------------
$go = Get-Command 'go.exe' -ErrorAction SilentlyContinue
if ($null -eq $go) {
    Fail 'go.exe not found on PATH. Install the toolchain pinned in docs/BUILD.md.'
}
$goVersionLine = (& $go.Source version)

$ccCandidates = @()
if ($env:CC) { $ccCandidates += $env:CC }
$ccCandidates += 'gcc.exe'
if (Get-Item 'env:MINGW64_ROOT' -ErrorAction SilentlyContinue) { $ccCandidates += (Join-Path $env:MINGW64_ROOT 'gcc.exe') }
$ccCandidates += 'C:\msys64\mingw64\bin\gcc.exe'
$cc = $null
foreach ($cand in $ccCandidates) {
    $resolved = Get-Command $cand -ErrorAction SilentlyContinue
    if ($null -ne $resolved) { $cc = $resolved.Source; break }
}
if ($null -eq $cc) {
    Fail 'mingw-w64 gcc not found. Add mingw64\bin to PATH or set MINGW64_ROOT. See docs/BUILD.md.'
}
# Pitfall (see docs/BUILD.md): mingw gcc fails silently when its bin directory
# is not on PATH, even when invoked by full path. Prepend it so both the gcc
# driver and cgo's compiler invocations resolve its helper DLLs and tools.
$ccDir = Split-Path -Parent $cc
if (-not ($env:PATH -like "*$ccDir*")) { $env:PATH = "$ccDir;$env:PATH" }
$ccVersionLine = (& $cc --version | Select-Object -First 1)
Write-Host "build.ps1: toolchain: $($goVersionLine); CC=$cc ($ccVersionLine)"

# --- 1. fetch-deps ----------------------------------------------------------
& (Join-Path $PSScriptRoot 'fetch-deps.ps1') -RepoRoot $RepoRoot
if ($LASTEXITCODE -ne 0) { Fail 'fetch-deps.ps1 failed (see output above).' }

# --- 2. frontend ------------------------------------------------------------
# This step BUILDS the page bundle the exe is supposed to carry. frontend/dist
# is produced here ('npm ci' when node_modules is absent, then 'npm run build'
# = tsc -b && vite build), and frontend/embed.go's '//go:embed all:dist' is the
# only channel that moves those bytes into build\wisp.exe. Without this step the
# build still succeeds - go:embed happily embeds the single tracked anchor file
# frontend/dist/.gitkeep - and the shipped exe carries zero page bytes. That is
# the defect ticket 274 was filed about.
# So node and npm are HARD requirements of this build line, and there is
# deliberately no skip path: unresolvable node or npm, a non-zero 'npm ci' or
# 'npm run build', or a dist that still looks like a clean checkout all make
# this script stop with a named cause and exit code 1.
$FrontendDir = Join-Path $RepoRoot 'frontend'
$DistDir = Join-Path $FrontendDir 'dist'
if (-not (Test-Path -LiteralPath (Join-Path $FrontendDir 'package.json'))) {
    Fail "frontend step: no package.json under $FrontendDir - the page source tree is missing, so nothing can produce the frontend/dist bytes go:embed puts into wisp.exe."
}

# Tool resolution. PowerShell maps a bare 'npm' lookup to its own npm.ps1 shim,
# so npm.cmd - the plain batch wrapper - is preferred: under the
# $ErrorActionPreference = 'Stop' at the top of this file the only failure this
# step should ever surface is npm's own exit code, not a wrapper's exception
# behaviour. Each lookup falls back to the extension-less name so a Node install
# that ships only 'node'/'npm' still resolves, and a miss is a hard failure
# rather than a skip.
function Resolve-BuildTool([string]$Preferred, [string]$Fallback) {
    $found = Get-Command $Preferred -ErrorAction SilentlyContinue
    if ($null -eq $found) { $found = Get-Command $Fallback -ErrorAction SilentlyContinue }
    return $found
}

$node = Resolve-BuildTool 'node.exe' 'node'
if ($null -eq $node) {
    Fail 'frontend step: node could not be resolved on PATH (tried node.exe, then node). The page bundle is a build input of wisp.exe, so this script fails instead of skipping it. Install the toolchain pinned in docs/BUILD.md, or run the build on a machine that has node and npm on PATH.'
}
$npm = Resolve-BuildTool 'npm.cmd' 'npm'
if ($null -eq $npm) {
    Fail 'frontend step: npm could not be resolved on PATH (tried npm.cmd, then npm). Without npm there is no frontend/dist, and an exe built from that tree carries no page at all - which is exactly what this step exists to prevent, so it fails instead of skipping.'
}
$nodeVersionLine = (& $node.Source --version)
$npmVersionLine = (& $npm.Source --version)
Write-Host "build.ps1: frontend: node = $($node.Source) ($nodeVersionLine); npm = $($npm.Source) ($npmVersionLine)"

if (Test-Path -LiteralPath (Join-Path $FrontendDir 'node_modules')) {
    Write-Host 'build.ps1: frontend: node_modules EXISTS -> taking the npm-ci-SKIPPED branch (npm run build only). This branch does not verify the installed tree against package-lock.json; delete frontend/node_modules to force a clean npm ci.'
} else {
    Write-Host 'build.ps1: frontend: node_modules ABSENT -> running npm ci (clean install pinned by package-lock.json).'
    Push-Location $FrontendDir
    try {
        & $npm.Source ci
        if ($LASTEXITCODE -ne 0) { Fail "frontend step: 'npm ci' failed with exit code $LASTEXITCODE in $FrontendDir (npm = $($npm.Source)). See the npm output above; the page bundle cannot be built without its dependencies." }
    } finally {
        Pop-Location
    }
}

Write-Host 'build.ps1: frontend: running npm run build (tsc -b && vite build, output goes to frontend/dist).'
Push-Location $FrontendDir
try {
    & $npm.Source run build
    if ($LASTEXITCODE -ne 0) { Fail "frontend step: 'npm run build' failed with exit code $LASTEXITCODE in $FrontendDir (npm = $($npm.Source)). See the npm output above." }
} finally {
    Pop-Location
}

# Post-build verification - the nail this ticket was really asking for.
# 'npm run build' exiting 0 is NOT evidence that page bytes exist: a stubbed or
# no-op npm leaves frontend/dist looking exactly like a clean checkout, one
# tracked anchor file (frontend/dist/.gitkeep) and no index.html, and go:embed
# embeds that anchor without complaint. Every branch below names its own cause.
$DistEntry = Join-Path $DistDir 'index.html'
if (-not (Test-Path -LiteralPath $DistDir)) {
    Fail "frontend step: npm run build reported success but $DistDir does not exist. Named cause: dist directory missing, so wisp.exe would embed no page."
}
$distAll = @(Get-ChildItem -LiteralPath $DistDir -Recurse -File -Force)
# The anchor-only shape is .gitkeep alone; anything real has to add at least one
# file beside it. Counting non-anchor files rather than all files keeps a build
# that wiped .gitkeep but did produce a page from failing for the wrong reason.
$distReal = @($distAll | Where-Object { $_.Name -ne '.gitkeep' })
if ($distReal.Count -le 0) {
    Fail "frontend step: npm run build reported success but frontend/dist holds only the anchor file(s) [$($distAll.Count) file(s), none beside frontend/dist/.gitkeep]. That is the clean-checkout shape, so the exe would ship with zero page bytes. Named cause: build produced no page artifacts."
}
if (-not (Test-Path -LiteralPath $DistEntry)) {
    Fail "frontend step: frontend/dist has $($distReal.Count) artifact file(s) but no index.html - the entry file internal/panel resolves (panel.EntryFile) is absent, so the bundle can never be served. Named cause: missing entry file."
}
$distEntryBytes = (Get-Item -LiteralPath $DistEntry).Length
if ($distEntryBytes -le 0) {
    Fail "frontend step: $DistEntry is 0 bytes - a present-but-empty entry file is not a page. Named cause: empty entry file."
}
$distRoster = ($distAll | Sort-Object FullName | ForEach-Object { "$($_.FullName.Substring($DistDir.Length + 1))=$($_.Length)" }) -join ' '
Write-Host "build.ps1: frontend ok: $($distAll.Count) file(s) in frontend/dist (entry index.html is $distEntryBytes bytes): $distRoster"

# --- 3. go build ------------------------------------------------------------
function Get-DepsValue([string]$Section, [string]$Key) {
    $inSection = $false
    foreach ($raw in (Get-Content -LiteralPath (Join-Path $RepoRoot 'deps.toml'))) {
        $line = $raw.Trim()
        if ($line.StartsWith('[')) {
            $inSection = ($line.Trim('[', ']').Trim() -eq $Section)
            continue
        }
        if (-not $inSection) { continue }
        if ($line -match ('^\s*' + [regex]::Escape($Key) + '\s*=\s*"([^"]*)"')) { return $Matches[1] }
    }
    return $null
}
$sherpaVersion = Get-DepsValue 'sherpa-onnx' 'version'
$ortVersion = Get-DepsValue 'onnxruntime' 'version'
if (-not $sherpaVersion -or -not $ortVersion) {
    Fail 'could not read sherpa-onnx/onnxruntime versions from deps.toml'
}

$commit = 'nogit'
$git = Get-Command 'git.exe' -ErrorAction SilentlyContinue
if ($null -ne $git) {
    $short = (& git rev-parse --short HEAD 2>$null)
    if ($LASTEXITCODE -eq 0 -and $short) { $commit = $short }
}
$buildDate = (Get-Date).ToUniversalTime().ToString('yyyy-MM-ddTHH:mm:ssZ')
# -H=windowsgui (SPEC-11 §2.2: no-arg launch = GUI): the win32 desktop binary must
# not pop a black console window when double-clicked. CLI legs (`wisp run` etc.)
# keep their stdout because attachParentConsole (cmd/wisp/console_windows.go:33-43)
# returns early when a usable stdout already exists (redirected or parent console)
# and otherwise attaches ATTACH_PARENT_PROCESS and rebinds the std handles.
$ldflags = (
    "-X $BuildInfoPkg.Version=0.0.0-dev",
    "-X $BuildInfoPkg.Commit=$commit",
    "-X $BuildInfoPkg.BuildDate=$buildDate",
    "-X $BuildInfoPkg.DefaultEnv=$Env",
    "-X $BuildInfoPkg.SherpaOnnxVersion=$sherpaVersion",
    "-X $BuildInfoPkg.OnnxRuntimeVersion=$ortVersion",
    "-H=windowsgui"
) -join ' '

$outDir = Join-Path $RepoRoot 'build'
New-Item -ItemType Directory -Path $outDir -Force | Out-Null

$env:CGO_ENABLED = '1'
$env:CC = $cc
$env:GOOS = 'windows'
$env:GOARCH = 'amd64'
if (-not $env:GOPROXY) { $env:GOPROXY = 'https://goproxy.cn,direct' }

Push-Location $RepoRoot
try {
    & $go.Source build -trimpath -ldflags $ldflags -o (Join-Path $outDir 'wisp.exe') ./cmd/wisp
    if ($LASTEXITCODE -ne 0) { Fail "go build failed with exit code $LASTEXITCODE" }
} finally {
    Pop-Location
}
Write-Host 'build.ps1: go build ok (cgo linked against sherpa-onnx C API)'

# --- 4. colocate DLLs with the exe (SPEC-11 §7.1) ---------------------------
$dllSource = Join-Path $RepoRoot 'third_party\sherpa-onnx'
Copy-Item -Path (Join-Path $dllSource '*.dll') -Destination $outDir -Force
Write-Host "build.ps1: DLLs colocated into $outDir"

# --- 4b. colocate the signed C29 model manifest (ticket 14) -----------------
# The runtime resolves <exe dir>/models/manifest.json (ResolveManifestPath);
# the manifest + its minisign signature ship with the binary and must stay in
# lockstep with the commit.
$manifestDir = Join-Path $outDir 'models'
New-Item -ItemType Directory -Force -Path $manifestDir | Out-Null
Copy-Item -Path (Join-Path $RepoRoot 'models\manifest.json')      -Destination $manifestDir -Force
Copy-Item -Path (Join-Path $RepoRoot 'models\manifest.json.minisig') -Destination $manifestDir -Force
Write-Host "build.ps1: signed model manifest colocated into $manifestDir"

# --- 5. SHA256SUMS ----------------------------------------------------------
$artifacts = @('wisp.exe', 'onnxruntime.dll', 'sherpa-onnx-c-api.dll', 'sherpa-onnx-cxx-api.dll')
$sums = @()
foreach ($f in $artifacts) {
    $p = Join-Path $outDir $f
    if (-not (Test-Path -LiteralPath $p)) { Fail "expected artifact missing: $f" }
    $h = (Get-FileHash -Algorithm SHA256 -LiteralPath $p).Hash.ToLowerInvariant()
    $sums += "$h  $f"
}
$sumsPath = Join-Path $outDir 'SHA256SUMS'
# LF line endings and UTF-8 without BOM so `sha256sum -c` (GNU) validates the file.
$utf8NoBom = New-Object System.Text.UTF8Encoding($false)
[System.IO.File]::WriteAllText($sumsPath, (($sums -join "`n") + "`n"), $utf8NoBom)
Write-Host "build.ps1: wrote $sumsPath"

# --- 6. smoke test ----------------------------------------------------------
Write-Host 'build.ps1: smoke test - running wisp.exe doctor'
& (Join-Path $outDir 'wisp.exe') doctor
if ($LASTEXITCODE -ne 0) { Fail 'wisp doctor reported FAIL.' }

Write-Host 'build.ps1: done. Artifacts in build\: wisp.exe, onnxruntime.dll, sherpa-onnx-c-api.dll, sherpa-onnx-cxx-api.dll, SHA256SUMS'
exit 0
