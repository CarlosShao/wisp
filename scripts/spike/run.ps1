# run.ps1 - S0 spike orchestrator (ticket 02).
# Builds the measurement programs (cgo + pure-Go flavors) and runs them all,
# writing JSON evidence into docs/evidence/s0/data/.
# Run from anywhere:  powershell -ExecutionPolicy Bypass -File scripts\spike\run.ps1
# Optional: -Stage baselines|verdict|goja|webview2|models|all   (default all)
param(
    [string]$Stage = "all",
    [int]$WebView2Children = 12
)

$ErrorActionPreference = "Stop"

# ---- resolve repo root / dirs
$SpikeDir = $PSScriptRoot
$RepoRoot = (Resolve-Path (Join-Path $SpikeDir "..\..")).Path
$BinDir   = Join-Path $SpikeDir "bin"
$DataDir  = Join-Path $RepoRoot "docs\evidence\s0\json"
$ModelsDir = Join-Path $RepoRoot "third_party\spike-models"
$DLLDir   = Join-Path $RepoRoot "third_party\sherpa-onnx"
New-Item -ItemType Directory -Force -Path $BinDir | Out-Null
New-Item -ItemType Directory -Force -Path $DataDir | Out-Null

# ---- resolve toolchain (same rules as build.ps1)
if (-not $env:GOPROXY) { $env:GOPROXY = "https://goproxy.cn,direct" }
$go = $null
foreach ($c in @("go", "$env:USERPROFILE\scoop\apps\go\current\bin\go.exe", "D:\work\base\go\bin\go.exe")) {
    try { $v = & $c version 2>$null; if ($LASTEXITCODE -eq 0) { $go = $c; break } } catch {}
}
if (-not $go) { throw "go not found" }
$cc = $null
if ($env:CC) { $cc = $env:CC }
foreach ($c in @("gcc.exe", "$env:MINGW64_ROOT\gcc.exe", "E:\work\base\msys64\mingw64\bin\gcc.exe", "C:\msys64\mingw64\bin\gcc.exe")) {
    if (-not $cc) { try { & $c --version 2>$null | Out-Null; if ($LASTEXITCODE -eq 0) { $cc = $c; break } } catch {} }
}
if (-not $cc) { throw "gcc not found (needed for cgo builds)" }
$ccDir = Split-Path $cc -Parent
$env:PATH = "$ccDir;$env:PATH"
Write-Host "go=$go cc=$cc"

function Build-Spike {
    param([string]$Pkg, [string]$Out, [bool]$Cgo, [string]$Tags)
    $exe = Join-Path $BinDir $Out
    Push-Location $SpikeDir
    try {
        $env:CGO_ENABLED = if ($Cgo) { "1" } else { "0" }
        $args = @("build", "-o", $exe)
        if ($Tags) { $args += @("-tags", $Tags) }
        $args += "./$Pkg"
        if ($Cgo) { $env:CC = $cc }
        & $go @args
        if ($LASTEXITCODE -ne 0) { throw "build failed: $Pkg" }
    } finally {
        Pop-Location
    }
    Write-Host "built $exe"
}

function Copy-DLLs {
    foreach ($d in @("onnxruntime.dll", "sherpa-onnx-c-api.dll", "sherpa-onnx-cxx-api.dll")) {
        $src = Join-Path $DLLDir $d
        if (Test-Path $src) { Copy-Item $src $BinDir -Force }
    }
}

function Run-Json {
    param([string]$Exe, [string[]]$Args, [string]$Out)
    $outPath = Join-Path $DataDir $Out
    $all = $Args + @("-out", $outPath)
    Write-Host ">> $Exe $($all -join ' ')"
    & $Exe @all
    if ($LASTEXITCODE -ne 0) { Write-Warning "$Exe exited $LASTEXITCODE" }
    return $outPath
}

# ---- builds
Build-Spike "shell-baseline"    "shell-baseline.exe"     $false ""
Build-Spike "webview2-latency"  "webview2-latency.exe"   $false ""
Build-Spike "goja-caps"         "goja-caps.exe"          $false ""
Build-Spike "xy-verdict"        "xy-verdict-nocgo.exe"   $false ""
Build-Spike "speech-baseline"   "speech-baseline.exe"    $true  "cgo_sherpa"
Build-Spike "xy-verdict"        "xy-verdict.exe"         $true  "cgo_sherpa"
Build-Spike "model-residency"   "model-residency.exe"    $true  "cgo_sherpa"
Copy-DLLs

function Run-All {
    Run-Json (Join-Path $BinDir "shell-baseline.exe") @() "01-shell-baselines.json" | Out-Null
    Run-Json (Join-Path $BinDir "speech-baseline.exe") @("-models", $ModelsDir) "02-speech-baselines.json" | Out-Null
    Run-Json (Join-Path $BinDir "xy-verdict.exe") @("-role", "idle-y") "03-idle-y.json" | Out-Null
    Run-Json (Join-Path $BinDir "xy-verdict-nocgo.exe") @("-role", "idle-x") "03-idle-x.json" | Out-Null
    Run-Json (Join-Path $BinDir "xy-verdict.exe") @("-role", "unload-test", "-which", "asr", "-models", $ModelsDir) "04-unload-asr.json" | Out-Null
    Run-Json (Join-Path $BinDir "goja-caps.exe") @() "05-goja.json" | Out-Null
    Run-Json (Join-Path $BinDir "webview2-latency.exe") @("-mode", "driver", "-children", "$WebView2Children") "06-webview2.json" | Out-Null
    Run-Json (Join-Path $BinDir "model-residency.exe") @("-which", "kws", "-models", $ModelsDir) "07-residency-kws.json" | Out-Null
    Run-Json (Join-Path $BinDir "model-residency.exe") @("-which", "vad", "-models", $ModelsDir) "07-residency-vad.json" | Out-Null
    Run-Json (Join-Path $BinDir "model-residency.exe") @("-which", "asr", "-models", $ModelsDir) "07-residency-asr.json" | Out-Null
    Run-Json (Join-Path $BinDir "model-residency.exe") @("-which", "tts", "-models", $ModelsDir) "07-residency-tts.json" | Out-Null
    Run-Json (Join-Path $BinDir "model-residency.exe") @("-which", "switch", "-models", $ModelsDir) "08-switch.json" | Out-Null
    Write-Host ""
    Write-Host "=== data files ==="
    Get-ChildItem $DataDir | Format-Table Name, Length -AutoSize
}

switch ($Stage) {
    "baselines" {
        Run-Json (Join-Path $BinDir "shell-baseline.exe") @() "01-shell-baselines.json" | Out-Null
        Run-Json (Join-Path $BinDir "speech-baseline.exe") @("-models", $ModelsDir) "02-speech-baselines.json" | Out-Null
    }
    "verdict" {
        Run-Json (Join-Path $BinDir "xy-verdict.exe") @("-role", "idle-y") "03-idle-y.json" | Out-Null
        Run-Json (Join-Path $BinDir "xy-verdict-nocgo.exe") @("-role", "idle-x") "03-idle-x.json" | Out-Null
        Run-Json (Join-Path $BinDir "xy-verdict.exe") @("-role", "unload-test", "-which", "asr", "-models", $ModelsDir) "04-unload-asr.json" | Out-Null
    }
    "goja"     { Run-Json (Join-Path $BinDir "goja-caps.exe") @() "05-goja.json" | Out-Null }
    "webview2" { Run-Json (Join-Path $BinDir "webview2-latency.exe") @("-mode", "driver", "-children", "$WebView2Children") "06-webview2.json" | Out-Null }
    "models" {
        Run-Json (Join-Path $BinDir "model-residency.exe") @("-which", "kws", "-models", $ModelsDir) "07-residency-kws.json" | Out-Null
        Run-Json (Join-Path $BinDir "model-residency.exe") @("-which", "vad", "-models", $ModelsDir) "07-residency-vad.json" | Out-Null
        Run-Json (Join-Path $BinDir "model-residency.exe") @("-which", "asr", "-models", $ModelsDir) "07-residency-asr.json" | Out-Null
        Run-Json (Join-Path $BinDir "model-residency.exe") @("-which", "tts", "-models", $ModelsDir) "07-residency-tts.json" | Out-Null
        Run-Json (Join-Path $BinDir "model-residency.exe") @("-which", "switch", "-models", $ModelsDir) "08-switch.json" | Out-Null
    }
    "all"      { Run-All }
    default    { throw "unknown stage $Stage" }
}
