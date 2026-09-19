#Requires -Version 5.1
<#
.SYNOPSIS
    Downloads the native dependencies pinned in deps.toml into third_party/
    (git-ignored). SPEC-11 §2.1: build-time download + SHA256 verify + cache.

.DESCRIPTION
    - Reads deps.toml from the repo root.
    - Cache fast path: if every pinned DLL already sits in third_party/sherpa-onnx/
      and its SHA256 matches the pin, exits in about a second (files are re-hashed
      on every run, so a corrupted cache is re-fetched, never trusted).
    - Otherwise downloads the pinned sherpa-onnx Windows shared release (direct
      GitHub first, then the mirror prefix), verifies the archive SHA256,
      extracts the pinned DLLs, verifies each DLL SHA256, and installs them flat
      into third_party/sherpa-onnx/.
    - Any hash mismatch = loud failure, bad file deleted. Never bypass this.

.USAGE
    scripts/fetch-deps.ps1 [-RepoRoot <path>]
#>
[CmdletBinding()]
param(
    [string]$RepoRoot
)

$ErrorActionPreference = 'Stop'
Set-StrictMode -Version 2.0
$ProgressPreference = 'SilentlyContinue'
try { [Net.ServicePointManager]::SecurityProtocol = [Net.ServicePointManager]::SecurityProtocol -bor [Net.SecurityProtocolType]::Tls12 } catch {}

if (-not $RepoRoot) {
    # $PSScriptRoot is empty in param defaults under some invocation styles
    # (see docs/BUILD.md); resolve here in the body instead.
    if ($PSScriptRoot) { $RepoRoot = Split-Path -Parent $PSScriptRoot } else { $RepoRoot = (Get-Location).Path }
}

function Split-TomlKey([string]$Value) {
    # Splits a dotted TOML key on dots outside double quotes; strips the quotes.
    $parts = New-Object System.Collections.Generic.List[string]
    $cur = New-Object System.Text.StringBuilder
    $inQuote = $false
    foreach ($ch in $Value.ToCharArray()) {
        if ($ch -eq '"') { $inQuote = -not $inQuote; continue }
        if ($ch -eq '.' -and -not $inQuote) { $parts.Add($cur.ToString()); [void]$cur.Clear(); continue }
        [void]$cur.Append($ch)
    }
    $parts.Add($cur.ToString())
    # The leading comma prevents PowerShell from unwinding a single-element
    # collection into a scalar (then [0] would return the first character).
    return ,$parts
}

function Read-DepsToml([string]$Path) {
    # Minimal TOML reader for deps.toml's schema: [sections], string values,
    # inline # comments. Returns a flat map: "section" and "section.key" -> value.
    $result = @{}
    $section = ''
    foreach ($raw in (Get-Content -LiteralPath $Path)) {
        $line = $raw.Trim()
        if ($line.Length -eq 0 -or $line.StartsWith('#')) { continue }
        if ($line.StartsWith('[') -and $line.EndsWith(']')) {
            $section = (Split-TomlKey $line.Substring(1, $line.Length - 2).Trim()) -join '.'
            $result[$section] = @{}
            continue
        }
        $eq = $line.IndexOf('=')
        if ($eq -lt 0 -or -not $result.ContainsKey($section)) { continue }
        $key = (Split-TomlKey $line.Substring(0, $eq).Trim())[0]
        $val = $line.Substring($eq + 1).Trim()
        if ($val -match '^"([^"]*)"') { $val = $Matches[1] }
        else { $val = ($val -replace '#.*$', '').Trim() }
        $result[$section][$key] = $val
    }
    return $result
}

function Get-Sha256([string]$Path) {
    return (Get-FileHash -Algorithm SHA256 -LiteralPath $Path).Hash.ToLowerInvariant()
}

function Invoke-DownloadWithFallback([string]$Url, [string]$MirrorPrefix, [string]$OutFile) {
    $urls = @($Url)
    if ($MirrorPrefix) { $urls += (($MirrorPrefix.TrimEnd('/')) + '/' + $Url) }
    foreach ($round in 1..2) {
        foreach ($u in $urls) {
            Write-Host "fetch-deps: downloading (attempt $round): $u"
            try {
                Invoke-WebRequest -Uri $u -OutFile $OutFile -UseBasicParsing -TimeoutSec 900
                if ((Get-Item -LiteralPath $OutFile).Length -gt 0) { return $true }
                Write-Warning "fetch-deps: downloaded file is empty: $u"
            } catch {
                Write-Warning ("fetch-deps: download failed: " + $_.Exception.Message)
            }
        }
        Start-Sleep -Seconds 2
    }
    return $false
}

$depsPath = Join-Path $RepoRoot 'deps.toml'
if (-not (Test-Path -LiteralPath $depsPath)) {
    Write-Host "fetch-deps: FATAL: deps.toml not found at $depsPath"
    exit 1
}
$deps = Read-DepsToml $depsPath
$sherpa = $deps['sherpa-onnx']
$targetDir = Join-Path $RepoRoot 'third_party\sherpa-onnx'
$dllSections = @($deps.Keys | Where-Object { $_ -like 'sherpa-onnx.dll.*' } | Sort-Object)
if ($dllSections.Count -eq 0) {
    Write-Host "fetch-deps: FATAL: no [sherpa-onnx.dll.*] entries in deps.toml"
    exit 1
}

# --- cache fast path: re-hash cached files against the pins -----------------
# The cache is only valid when (a) the manifest records the exact archive pin
# from deps.toml, and (b) every cached DLL still hashes to its pin. This keeps
# a tampered deps.toml pin from being silently satisfied by stale cache.
$cachedOk = $true
$manifestPath = Join-Path $targetDir '.cache-manifest.json'
if (-not (Test-Path -LiteralPath $manifestPath)) {
    $cachedOk = $false
} else {
    try {
        $manifest = Get-Content -LiteralPath $manifestPath -Raw | ConvertFrom-Json
        if ($manifest.sha256 -ne $sherpa['sha256'] -or $manifest.version -ne $sherpa['version']) { $cachedOk = $false }
    } catch { $cachedOk = $false }
}
if ($cachedOk) {
    foreach ($sec in $dllSections) {
        $dest = Join-Path $targetDir ($sec -replace '^sherpa-onnx\.dll\.', '')
        if (-not (Test-Path -LiteralPath $dest)) { $cachedOk = $false; break }
        if ((Get-Sha256 $dest) -ne $deps[$sec]['sha256']) { $cachedOk = $false; break }
    }
}
if ($cachedOk) {
    Write-Host "fetch-deps: cache hit - third_party/sherpa-onnx matches deps.toml (sherpa-onnx $($sherpa['version']))"
    exit 0
}

# --- download, verify, extract ----------------------------------------------
$tmp = Join-Path ([System.IO.Path]::GetTempPath()) ("wisp-fetch-deps-" + [guid]::NewGuid().ToString('N').Substring(0, 8))
New-Item -ItemType Directory -Path $tmp -Force | Out-Null
$archiveName = [System.IO.Path]::GetFileName($sherpa['source'])
$archive = Join-Path $tmp $archiveName
try {
    if (-not (Invoke-DownloadWithFallback -Url $sherpa['source'] -MirrorPrefix $sherpa['mirror_prefix'] -OutFile $archive)) {
        Write-Host "fetch-deps: FATAL: could not download $($sherpa['source']) (direct and mirror both failed)"
        exit 1
    }
    $actual = Get-Sha256 $archive
    if ($actual -ne $sherpa['sha256']) {
        Remove-Item -LiteralPath $archive -Force -ErrorAction SilentlyContinue
        Write-Host "fetch-deps: FATAL: SHA256 MISMATCH for the downloaded archive."
        Write-Host "  expected: $($sherpa['sha256'])"
        Write-Host "  actual:   $actual"
        Write-Host "  The bad file has been deleted. Do not edit the pin to make this pass;"
        Write-Host "  investigate the source instead (deps.toml is the trust anchor, SPEC-11 §2.1)."
        exit 1
    }
    Write-Host "fetch-deps: archive SHA256 verified ($archiveName)"

    & "$env:SystemRoot\System32\tar.exe" -xf $archive -C $tmp
    if ($LASTEXITCODE -ne 0) {
        Write-Host "fetch-deps: FATAL: tar extraction failed with exit code $LASTEXITCODE"
        exit 1
    }

    if (Test-Path -LiteralPath $targetDir) { Remove-Item -LiteralPath $targetDir -Recurse -Force }
    New-Item -ItemType Directory -Path $targetDir -Force | Out-Null

    foreach ($sec in $dllSections) {
        $pin = $deps[$sec]
        $destName = $sec -replace '^sherpa-onnx\.dll\.', ''
        $rel = ($pin['archive_path'] -replace '/', '\')
        $src = Get-ChildItem -Path $tmp -Recurse -File | Where-Object { $_.FullName.EndsWith($rel) } | Select-Object -First 1
        if ($null -eq $src) {
            Write-Host "fetch-deps: FATAL: archive does not contain $($pin['archive_path'])"
            exit 1
        }
        $h = Get-Sha256 $src.FullName
        if ($h -ne $pin['sha256']) {
            Write-Host "fetch-deps: FATAL: SHA256 MISMATCH for $($pin['archive_path'])."
            Write-Host "  expected: $($pin['sha256'])"
            Write-Host "  actual:   $h"
            exit 1
        }
        Copy-Item -LiteralPath $src.FullName -Destination (Join-Path $targetDir $destName) -Force
        Write-Host "fetch-deps: extracted $($pin['archive_path']) -> third_party/sherpa-onnx/$destName (sha256 ok)"
    }

    @{
        version    = $sherpa['version']
        url        = $sherpa['source']
        sha256     = $sherpa['sha256']
        fetched_at = (Get-Date).ToUniversalTime().ToString('o')
    } | ConvertTo-Json | Set-Content -LiteralPath (Join-Path $targetDir '.cache-manifest.json') -Encoding ascii

    Write-Host "fetch-deps: done (sherpa-onnx $($sherpa['version']))"
    exit 0
} finally {
    Remove-Item -LiteralPath $tmp -Recurse -Force -ErrorAction SilentlyContinue
}
