$path = 'D:\work\workspace\projects plans\Wisp\build\wisp.exe'
$b = [IO.File]::ReadAllBytes($path)
$pe = [BitConverter]::ToInt32($b, 60)
$sub = [BitConverter]::ToInt16($b, $pe + 92)
Write-Output ("pe_offset=" + $pe + " subsystem=" + $sub + " (2=GUI, 3=CUI)")
