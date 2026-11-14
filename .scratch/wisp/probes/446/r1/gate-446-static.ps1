# gate-446-static.ps1 -- the static half of AC#3.
#
# AC#3 has two mutation directions and each needs a check that can go RED:
#   M1  the JUDGEMENT goes back to comparing the message text  -> behaviour probe
#       (probe446_errno_presence_windows_test.go) catches it on the wrapped
#       errno=0 rows.
#   M2  one CALL SITE goes back to an inline comparison and stops using the
#       judgement -> THIS script catches it, because the behaviour probe would
#       still be green (it keeps calling a correct errnoPresent that the call
#       site no longer uses).
#
# Reads only internal/tools/system_backend_windows.go. Exit 0 = GREEN, 1 = RED.
# The sentence is assembled from pieces so that this script does not itself
# become a hit for the census it enforces.

$ErrorActionPreference = 'Stop'
$file = 'internal/tools/system_backend_windows.go'
if (-not (Test-Path $file)) { Write-Output "RED: $file missing"; exit 1 }
$src = Get-Content $file

$sentence = 'operation ' + 'completed ' + 'successfully'
$red = @()

$sentenceHits = @($src | Select-String -SimpleMatch $sentence)
if ($sentenceHits.Count -ne 0) {
  $red += ("sentence still decided on at line(s): " + (($sentenceHits | ForEach-Object { $_.LineNumber }) -join ','))
}

$textCompare = @($src | Select-String -Pattern '\.Error\(\)\s*(==|!=)\s*"')
if ($textCompare.Count -ne 0) {
  $red += ("error TEXT compared to a literal at line(s): " + (($textCompare | ForEach-Object { $_.LineNumber }) -join ','))
}

$defs = @($src | Select-String -Pattern '^func errnoPresent\(e error\) bool \{')
if ($defs.Count -ne 1) { $red += "expected exactly 1 definition of errnoPresent, found $($defs.Count)" }

$calls = @($src | Select-String -Pattern 'errnoPresent\(e\)')
if ($calls.Count -ne 2) {
  $red += ("expected exactly 2 call sites of errnoPresent(e), found {0} at line(s) {1}" -f ($calls.Count), (($calls | ForEach-Object { $_.LineNumber }) -join ','))
}

$codeJudge = @($src | Select-String -Pattern 'errors\.Is\(e, syscall\.Errno\(0\)\)')
if ($codeJudge.Count -ne 1) { $red += "expected exactly 1 code-based judgement errors.Is(e, syscall.Errno(0)), found $($codeJudge.Count)" }

Write-Output ("readings: sentence-hits={0} text-compares={1} errnoPresent-defs={2} errnoPresent-calls={3} code-judgements={4}" -f `
  $sentenceHits.Count, $textCompare.Count, $defs.Count, $calls.Count, $codeJudge.Count)
$calls | ForEach-Object { Write-Output ("  call site line {0}: {1}" -f $_.LineNumber, $_.Line.Trim()) }

if ($red.Count -ne 0) {
  $red | ForEach-Object { Write-Output "RED: $_" }
  exit 1
}
Write-Output 'GREEN: no sentence decides anything in this file; both call sites use the code judgement.'
exit 0
