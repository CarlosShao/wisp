#Requires -Version 5.1
<#
.SYNOPSIS
    One-shot build of wisp.exe on Windows (SPEC-11 §2.2). From a clean clone,
    this is the only command needed besides installing the toolchain in
    docs/BUILD.md.

.DESCRIPTION
    1. scripts/fetch-deps.ps1        (verify/complete third_party/, cached runs finish fast)
    2. frontend: skipped (none yet, lands in S5)
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
Write-Host 'build.ps1: frontend step skipped (no frontend yet; embed lands in S5 per SPEC-11 §2.2)'

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
$ldflags = (
    "-X $BuildInfoPkg.Version=0.0.0-dev",
    "-X $BuildInfoPkg.Commit=$commit",
    "-X $BuildInfoPkg.BuildDate=$buildDate",
    "-X $BuildInfoPkg.DefaultEnv=$Env",
    "-X $BuildInfoPkg.SherpaOnnxVersion=$sherpaVersion",
    "-X $BuildInfoPkg.OnnxRuntimeVersion=$ortVersion"
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
