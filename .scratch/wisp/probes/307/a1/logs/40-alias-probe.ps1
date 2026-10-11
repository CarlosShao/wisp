$ErrorActionPreference = 'Continue'
$fs = New-Object -ComObject Scripting.FileSystemObject

function ShortOf($p) {
  try {
    if (Test-Path -LiteralPath $p -PathType Container) {
      return $fs.GetFolder($p).ShortPath
    }
    if (Test-Path -LiteralPath $p -PathType Leaf) {
      return $fs.GetFile($p).ShortPath
    }
    return 'PATH-MISSING'
  } catch {
    return 'ERROR:' + $_.Exception.Message
  }
}

Write-Output ('COMPUTERNAME=' + $env:COMPUTERNAME)
Write-Output ('USERNAME=' + $env:USERNAME)
Write-Output ('USERPROFILE=' + $env:USERPROFILE)
Write-Output ('TEMP=' + $env:TEMP)
Write-Output ('TMP=' + $env:TMP)
Write-Output ('TMPDIR=' + $env:TMPDIR)
Write-Output ('LOCALAPPDATA=' + $env:LOCALAPPDATA)

$targets = @(
  $env:USERPROFILE,
  'C:\Users',
  'C:\Users\swq',
  $env:LOCALAPPDATA,
  (Join-Path $env:LOCALAPPDATA 'Temp'),
  'C:\Users\swq\tmp',
  'C:\Users\swq\tmp\riskattrib1',
  'C:\Users\swq\tmp\riskattrib1\RunnerTempAreaForAttribution1'
)

foreach ($t in $targets) {
  if ([string]::IsNullOrEmpty($t)) { continue }
  $s = ShortOf $t
  $same = if ($s -eq $t) { 'SHORT==LONG' } else { 'SHORT!=LONG' }
  Write-Output ('ALIAS ' + $same + ' long=' + $t + ' short=' + $s)
}

Write-Output '--- segment walk of LOCALAPPDATA\Temp (each prefix, short vs long) ---'
$la = Join-Path $env:LOCALAPPDATA 'Temp'
$cur = ''
foreach ($seg in ($la -split '\\')) {
  if ($seg -eq '') { continue }
  if ($cur -eq '') { $cur = $seg } else { $cur = $cur + '\' + $seg }
  $s = ShortOf $cur
  $same = if ($s -eq $cur) { 'SHORT==LONG' } else { 'SHORT!=LONG' }
  Write-Output ('SEG ' + $same + ' prefix=' + $cur + ' short=' + $s)
}

Write-Output '--- dir listing of C:\Users (does any account name have a short alias?) ---'
try {
  foreach ($d in (Get-ChildItem -LiteralPath 'C:\Users' -Directory -Force)) {
    $s = ShortOf $d.FullName
    $same = if ($s -eq $d.FullName) { 'SHORT==LONG' } else { 'SHORT!=LONG' }
    Write-Output ('USERS ' + $same + ' long=' + $d.FullName + ' short=' + $s)
  }
} catch {
  Write-Output ('USERS-ERROR ' + $_.Exception.Message)
}

Write-Output '--- ntfs 8dot3 name creation behaviour (read-only query) ---'
& fsutil behavior query Usn8dot3namecreation
Write-Output ('fsutil_rc=' + $LASTEXITCODE)
