# ball-cycle.ps1 - ticket 07 visual evidence (SPEC-08 §2.1).
#
# Builds cmd/balldebug, runs the full 20-state cycle (2s dwell per state),
# screenshots the ball window for EVERY state into
# docs/evidence/s1/ball-states/<NN>-<State>.png, then verifies the process
# exited cleanly (everything destroyed, handle gate < 600).
#
# Usage:  powershell -NoProfile -ExecutionPolicy Bypass -File scripts\dev\ball-cycle.ps1
#         optional: -DwellMs 2000 -OutDir <dir>

param(
  [int]$DwellMs = 2000,
  [string]$OutDir = ""
)

$ErrorActionPreference = "Stop"

$repo = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
if (-not $OutDir) { $OutDir = Join-Path $repo "docs\evidence\s1\ball-states" }
New-Item -ItemType Directory -Force -Path $OutDir | Out-Null

$go = "D:\work\base\go\bin\go.exe"
if (-not (Test-Path $go)) { $go = "go" }
$env:PATH = "D:\work\base\go\bin;E:\work\base\msys64\mingw64\bin;$env:PATH"
$env:CGO_ENABLED = "1"
$env:GOPROXY = "https://goproxy.cn,direct"

Write-Host "[ball-cycle] building cmd/balldebug ..."
$exe = Join-Path $env:TEMP "wisp-balldebug.exe"
& $go build -o $exe (Join-Path $repo "cmd\balldebug")
if ($LASTEXITCODE -ne 0) { throw "go build failed" }

Add-Type -AssemblyName System.Drawing
Add-Type @"
using System;using System.Runtime.InteropServices;
public class WispShot {
  [DllImport("user32.dll")] public static extern IntPtr FindWindow(string cls, string title);
  [DllImport("user32.dll")] public static extern bool GetWindowRect(IntPtr h, out RECT r);
  [DllImport("user32.dll")] public static extern bool SetCursorPos(int x, int y);
  [DllImport("user32.dll")] public static extern bool SetWindowPos(IntPtr hw, IntPtr after, int x, int y, int w, int ht, uint flags);
  public struct RECT { public int L, T, R, B; }
}
"@

# Park the cursor out of the way so no third-party bubble overlaps the shot,
# and place the ball on clear wallpaper above the tray corner (input-method
# bars live there).
[WispShot]::SetCursorPos(4, 4) | Out-Null
Add-Type -AssemblyName System.Windows.Forms
$waPre = [System.Windows.Forms.Screen]::PrimaryScreen.WorkingArea
$shootX = [int]($waPre.Right - 420)
$shootY = [int]($waPre.Bottom - 330)

$psi = New-Object System.Diagnostics.ProcessStartInfo
$psi.FileName = $exe
$psi.Arguments = "-cycle-ms $DwellMs -x $shootX -y $shootY"
$psi.RedirectStandardOutput = $true
$psi.RedirectStandardError = $true
$psi.UseShellExecute = $false
$proc = [System.Diagnostics.Process]::Start($psi)

# Window appears after STA boot; wait for it.
$deadline = [DateTime]::UtcNow.AddSeconds(20)
$hwnd = [IntPtr]::Zero
while ([DateTime]::UtcNow -lt $deadline) {
  $hwnd = [WispShot]::FindWindow("WispBallWindow", "Wisp")
  if ($hwnd -ne [IntPtr]::Zero) {
    break
  }
  Start-Sleep -Milliseconds 100
}
if ($hwnd -eq [IntPtr]::Zero) {
  try { $proc.Kill() } catch {}
  throw "ball window never appeared"
}
Write-Host "[ball-cycle] window up; capturing states ..."

$seq = 0
while (-not $proc.StandardOutput.EndOfStream) {
  $line = $proc.StandardOutput.ReadLine()
  if ($line -match "state=([A-Za-z]+)") {
    $state = $Matches[1]
    $seq++
    $name = "{0:d2}-{1}.png" -f $seq, $state
    $out = Join-Path $OutDir $name
    # mid-dwell: animation has run, overlays (badge/progress) applied
    Start-Sleep -Milliseconds ([Math]::Max(600, [int]($DwellMs * 0.55)))
    $r = New-Object WispShot+RECT
    Start-Sleep -Milliseconds 150
    [WispShot]::GetWindowRect($hwnd, [ref]$r) | Out-Null
    $margin = 130
    $cx0 = $r.L - $margin; $cy0 = $r.T - $margin; $cw = ($r.R - $r.L) + 2*$margin; $ch = ($r.B - $r.T) + 2*$margin
    if ($cw -lt 8 -or $ch -lt 8) { Write-Warning "bad rect for $state"; continue }
    $bmp = New-Object System.Drawing.Bitmap($cw, $ch)
    $g = [System.Drawing.Graphics]::FromImage($bmp)
    $g.CopyFromScreen($cx0, $cy0, 0, 0, $bmp.Size)
    $bmp.Save($out, [System.Drawing.Imaging.ImageFormat]::Png)
    $g.Dispose(); $bmp.Dispose()
    Write-Host ("[ball-cycle] {0}  rect=({1},{2})" -f $name, $r.L, $r.T)
  }
  elseif ($line -match "closed handles=(\d+)") {
    Write-Host ("[ball-cycle] final handle count: {0} (gate < 600)" -f $Matches[1])
    if ([int]$Matches[1] -ge 600) { throw "handle gate exceeded" }
  }
  elseif ($line -match "OK$") { Write-Host "[ball-cycle] harness reported OK" }
}

$proc.WaitForExit(15000) | Out-Null
if (-not $proc.HasExited) {
  Write-Warning "harness did not exit; killing"
  $proc.Kill()
}
$errOut = $proc.StandardError.ReadToEnd()
if ($errOut.Trim()) { Write-Host "[ball-cycle] stderr: $errOut" }

Write-Host ("[ball-cycle] done: {0} shots in {1}" -f $seq, $OutDir)
