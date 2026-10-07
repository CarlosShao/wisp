@echo off
rem 274-v3 stub npm for AC#11 case 4: on "npm run build" writes exactly ONE new
rem artifact whose name contains spaces, into tree-c's dist; exits 0 otherwise.
if /I "%1"=="run" (
  echo page bytes written by this run's stub npm> "D:\tmp\wisp274v3\tree-c\frontend\dist\stale spaced page.html"
  echo 274-v3 stub-write npm ran build, wrote the spaced artifact
  exit /b 0
)
echo 274-v3 stub-write npm (%*)
exit /b 0
