param(
    [Parameter(Mandatory = $true)][string]$Tree,
    [Parameter(Mandatory = $true)][string]$StubBin,
    [Parameter(Mandatory = $true)][string]$Log,
    [string]$Mode = 'nothing'
)
# 274-r3 case driver. Runs the MIRROR TREE's scripts\build.ps1 with:
#   * git.exe removed from PATH (the tree has no .git; under $ErrorActionPreference='Stop'
#     a stderr-writing git rev-parse turns into a terminating error - v1/v2 did the same),
#   * the stub npm dir first on PATH (STUB_MODE picks whether it writes one space-named file).
$ErrorActionPreference = 'Continue'

$gitless = @($env:PATH -split ';' | Where-Object {
        $_ -ne '' -and -not (Test-Path -LiteralPath (Join-Path $_ 'git.exe') -ErrorAction SilentlyContinue)
    })
$env:PATH = (@(@($StubBin) + $gitless) -join ';')
$env:STUB_MODE = $Mode

Remove-Item -LiteralPath $Log -ErrorAction SilentlyContinue
"=== CASE tree=$Tree mode=$Mode ===" | Out-File -FilePath $Log -Append -Encoding utf8
"NPM_RESOLVES_TO: $((Get-Command npm.cmd -ErrorAction SilentlyContinue).Source)" | Out-File -FilePath $Log -Append -Encoding utf8
"GIT_ON_PATH: $((Get-Command git.exe -ErrorAction SilentlyContinue).Source)" | Out-File -FilePath $Log -Append -Encoding utf8
$build = Join-Path $Tree 'scripts\build.ps1'
"BUILD_SCRIPT: $build" | Out-File -FilePath $Log -Append -Encoding utf8
& powershell.exe -NoProfile -ExecutionPolicy Bypass -File $build -Env dev 2>&1 | Out-File -FilePath $Log -Append -Encoding utf8
"PROCESS_EXITCODE=$LASTEXITCODE" | Out-File -FilePath $Log -Append -Encoding utf8
Write-Output "CASE_DONE $Tree mode=$Mode rc=$LASTEXITCODE log=$Log"
