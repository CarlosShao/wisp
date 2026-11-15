# run-probe.ps1 -- ticket 455 recon: which onnxruntime.dll does the process bind?
#
# Answers with a PATH and a SIZE, never with the version string alone.
# Arms:
#   A  carrier shape  : probe.exe + the three pinned DLLs in ONE dir, cwd=cmd/wisp
#   B  raw go-test shape: probe.exe in a DLL-free dir, third_party/sherpa-onnx on PATH, cwd=cmd/wisp
#   C1 explicit LoadLibraryEx(other dir\sherpa-onnx-c-api.dll), flags=0
#   C2 same, flags=0x8  (LOAD_WITH_ALTERED_SEARCH_PATH -> the DLL's own dir becomes a search root)
#   C3 same, flags=0x100 (LOAD_LIBRARY_SEARCH_DLL_LOAD_DIR)
#
# Usage: powershell -ExecutionPolicy Bypass -File .scratch\wisp\probes\455\run-probe.ps1
$ErrorActionPreference = 'Stop'

$here   = Split-Path -Parent $MyInvocation.MyCommand.Path
$root   = (Resolve-Path (Join-Path $here '..\..\..\..')).Path
$stage  = Join-Path $here 'stage'
$tmpRaw = Join-Path $here 'tmp-raw'
New-Item -ItemType Directory -Force -Path $stage  | Out-Null
New-Item -ItemType Directory -Force -Path $tmpRaw | Out-Null

Write-Host "repo root: $root"

Write-Host "`n=== candidate files (size + file version) ==="
$cands = @(
    (Join-Path $env:SystemRoot 'System32\onnxruntime.dll'),
    (Join-Path $root 'third_party\sherpa-onnx\onnxruntime.dll'),
    (Join-Path $root '.scratch\wisp\stage-cli\onnxruntime.dll')
)
foreach ($c in $cands) {
    if (Test-Path $c) {
        $i = Get-Item $c
        Write-Host ("  {0}`n      {1} B   FileVersion={2}" -f $c, $i.Length, $i.VersionInfo.FileVersion)
    } else {
        Write-Host "  MISSING: $c"
    }
}

Write-Host "`n=== build ==="
Push-Location $here
try {
    & go build -o (Join-Path $stage 'probe.exe') .
    if ($LASTEXITCODE -ne 0) { throw "go build probe failed: $LASTEXITCODE" }
    & go build -o (Join-Path $stage 'loadlib_probe.exe') .\loadlib_probe
    if ($LASTEXITCODE -ne 0) { throw "go build loadlib_probe failed: $LASTEXITCODE" }
} finally { Pop-Location }

# Stage the pinned DLLs next to probe.exe; the raw arm gets only the exes.
Copy-Item (Join-Path $root 'third_party\sherpa-onnx\*.dll') $stage -Force
Copy-Item (Join-Path $stage 'probe.exe')         $tmpRaw -Force
Copy-Item (Join-Path $stage 'loadlib_probe.exe') $tmpRaw -Force

$cmdwisp   = Join-Path $root 'cmd\wisp'
$sherpaDll = Join-Path $stage 'sherpa-onnx-c-api.dll'

Write-Host "`n=== ARM A: carrier shape (DLLs staged next to the exe), cwd=cmd/wisp ==="
Push-Location $cmdwisp
try { & (Join-Path $stage 'probe.exe') } finally { Pop-Location }

Write-Host "`n=== ARM B: raw go-test shape (exe in a DLL-free dir, sherpa on PATH), cwd=cmd/wisp ==="
$oldPath = $env:PATH
$env:PATH = (Join-Path $root 'third_party\sherpa-onnx') + ';' + $oldPath
Push-Location $cmdwisp
try { & (Join-Path $tmpRaw 'probe.exe') } finally { Pop-Location; $env:PATH = $oldPath }

Write-Host "`n=== ARM C1: explicit LoadLibraryEx(absolute other-dir sherpa), flags=0 ==="
Push-Location $cmdwisp
try { & (Join-Path $tmpRaw 'loadlib_probe.exe') $sherpaDll 0 } finally { Pop-Location }

Write-Host "`n=== ARM C2: explicit LoadLibraryEx(absolute other-dir sherpa), flags=8 (LOAD_WITH_ALTERED_SEARCH_PATH) ==="
Push-Location $cmdwisp
try { & (Join-Path $tmpRaw 'loadlib_probe.exe') $sherpaDll 8 } finally { Pop-Location }

Write-Host "`n=== ARM C3: explicit LoadLibraryEx(absolute other-dir sherpa), flags=100 (LOAD_LIBRARY_SEARCH_DLL_LOAD_DIR) ==="
Push-Location $cmdwisp
try { & (Join-Path $tmpRaw 'loadlib_probe.exe') $sherpaDll 100 } finally { Pop-Location }

Write-Host "`n=== done ==="
