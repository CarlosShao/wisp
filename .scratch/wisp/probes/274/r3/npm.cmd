@echo off
echo stub-npm 274-r3: STUB_MODE=%STUB_MODE% args=%*
if "%STUB_MODE%"=="fresh" echo fresh-page-bytes-from-this-run > "dist\index with a space.html"
exit /b 0
