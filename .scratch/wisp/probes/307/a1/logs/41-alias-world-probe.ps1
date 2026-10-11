$ErrorActionPreference = 'Continue'
$fs = New-Object -ComObject Scripting.FileSystemObject

function ShortOf($p) {
  try {
    if (Test-Path -LiteralPath $p -PathType Container) { return $fs.GetFolder($p).ShortPath }
    if (Test-Path -LiteralPath $p -PathType Leaf) { return $fs.GetFile($p).ShortPath }
    return 'PATH-MISSING'
  } catch {
    return 'ERROR'
  }
}

$base = Join-Path $env:TEMP 'wisp307a1AliasWorldProbe'
$leaf = Join-Path $base 'LongSegmentNameBeyondEight'
New-Item -ItemType Directory -Force -Path $base   | Out-Null
New-Item -ItemType Directory -Force -Path $leaf   | Out-Null
Write-Output ('PROBE-BASE long=' + $base)
Write-Output ('PROBE-BASE ' + (ShortOf $base))
Write-Output ('PROBE-LEAF long=' + $leaf)
Write-Output ('PROBE-LEAF ' + (ShortOf $leaf))

Write-Output '--- fsutil 8dot3name query (machine + C: volume state, read-only) ---'
& fsutil 8dot3name query
Write-Output ('fsutil_machine_rc=' + $LASTEXITCODE)
& fsutil 8dot3name query C:
Write-Output ('fsutil_volume_rc=' + $LASTEXITCODE)

Write-Output '--- what Go would see: TMP/TEMP/TMPDIR as inherited by a child process ---'
Write-Output ('TEMP=' + $env:TEMP)
Write-Output ('TMP=' + $env:TMP)
Write-Output ('TMPDIR=' + $env:TMPDIR)
Write-Output ('TEMP-short=' + (ShortOf $env:TEMP))
