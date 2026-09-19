#Requires -Version 5.1
<#
.SYNOPSIS
    C29 model-manifest signing (ticket 14, D33/F3). Signs models/manifest.json
    with the minisign-format models key and verifies it back through the
    runtime verifier.

.DESCRIPTION
    The SECRET key never enters the repository: it is injected by path, from
    (in order) -KeyPath, $env:WISP_MINISIGN_KEY, or the default dev location
    E:\work\base\wisp-minisign\wisp-models.key. The dev key is an unencrypted
    DEV-ONLY envelope; the production key ceremony is S8 and MUST rotate.

    Key regeneration (loses the old key => re-sign everything + new buildinfo
    pubkey): scripts\sign-models.ps1 -Keygen -KeyDir E:\work\base\wisp-minisign

.USAGE
    scripts\sign-models.ps1                      # sign models/manifest.json
    scripts\sign-models.ps1 -Verify              # verify only (no key needed)
    scripts\sign-models.ps1 -Keygen              # generate a new dev keypair
#>
[CmdletBinding()]
param(
    [string]$KeyPath,
    [string]$KeyDir = 'E:\work\base\wisp-minisign',
    [switch]$Verify,
    [switch]$Keygen,
    [string]$Manifest
)

$ErrorActionPreference = 'Stop'
Set-StrictMode -Version 2.0

$RepoRoot = Split-Path -Parent $PSScriptRoot
if (-not $Manifest) { $Manifest = Join-Path $RepoRoot 'models\manifest.json' }
$go = Get-Command 'go.exe' -ErrorAction SilentlyContinue
if ($null -eq $go) { $go = Get-Item 'D:\work\base\go\bin\go.exe' -ErrorAction SilentlyContinue }
if ($null -eq $go) { Write-Host 'sign-models: FATAL: go.exe not found'; exit 1 }
$env:GOPROXY = 'https://goproxy.cn,direct'

Push-Location $RepoRoot
try {
    if ($Keygen) {
        & $go.Source run ./tools/signmodels keygen -dir $KeyDir
        if ($LASTEXITCODE -ne 0) { Write-Host 'sign-models: FATAL: keygen failed'; exit 1 }
        Write-Host "sign-models: NEXT STEP: paste the second line of $KeyDir\wisp-models.pub into internal/buildinfo.MinisignPublicKey, then re-sign."
        return
    }

    $sigPath = "$Manifest.minisig"
    if ($Verify) {
        & $go.Source run ./tools/signmodels verify -pub "$KeyDir\wisp-models.pub" -in $Manifest -sig $sigPath
        if ($LASTEXITCODE -ne 0) { Write-Host 'sign-models: FATAL: verification failed'; exit 1 }
        return
    }

    if (-not $KeyPath) { $KeyPath = $env:WISP_MINISIGN_KEY }
    if (-not $KeyPath) { $KeyPath = Join-Path $KeyDir 'wisp-models.key' }
    if (-not (Test-Path $KeyPath)) {
        Write-Host "sign-models: FATAL: secret key not found at $KeyPath (never commit it; generate with -Keygen)"
        exit 1
    }

    & $go.Source run ./tools/signmodels sign -key $KeyPath -in $Manifest -out $sigPath
    if ($LASTEXITCODE -ne 0) { Write-Host 'sign-models: FATAL: signing failed'; exit 1 }
    Write-Host "sign-models: commit BOTH $Manifest and $sigPath together."
} finally {
    Pop-Location
}
