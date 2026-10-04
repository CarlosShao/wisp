param(
    [string]$Exe = 'bin\probe-gui.exe'
)
# Is the exit-code probe slow enough to detect "the wait was removed"?
# `cmd /c exit 7` finishes so fast that a helper which forgot to wait can still
# read a real ExitCode, i.e. it cannot see that shape. This probe checks the
# alternative: a command that lives ~1s before exiting 7.
$ErrorActionPreference = 'Stop'
Set-StrictMode -Version 2.0

function Get-NativeArgumentLine {
    param([string[]]$Values)
    $quoted = @(foreach ($v in $Values) {
        $s = [string]$v
        if ($s -match '[\s"]') { '"' + ($s -replace '"', '\"') + '"' } else { $s }
    })
    return ($quoted -join ' ')
}

$cases = @(
    @{ name = 'cmd-fast';  args = @('/c', 'exit 7') },
    @{ name = 'cmd-slow';  args = @('/c', 'ping -n 2 127.0.0.1 >nul & exit 7') }
)
foreach ($case in $cases) {
    $psi = New-Object System.Diagnostics.ProcessStartInfo
    $psi.FileName = 'C:\Windows\System32\cmd.exe'
    $psi.Arguments = (Get-NativeArgumentLine -Values $case.args)
    $psi.UseShellExecute = $false
    Write-Host ("probe-case {0}: command-line = {1} | {2}" -f $case.name, $psi.FileName, $psi.Arguments)
    $proc = [System.Diagnostics.Process]::Start($psi)
    $proc.WaitForExit()
    Write-Host ("probe-case {0}: waited exit={1}" -f $case.name, [int]$proc.ExitCode)
    $proc.Dispose()

    $psi2 = New-Object System.Diagnostics.ProcessStartInfo
    $psi2.FileName = 'C:\Windows\System32\cmd.exe'
    $psi2.Arguments = (Get-NativeArgumentLine -Values $case.args)
    $psi2.UseShellExecute = $false
    $p2 = [System.Diagnostics.Process]::Start($psi2)
    # The mutated shape: no WaitForExit() at all.
    try {
        $code2 = [int]$p2.ExitCode
        Write-Host ("probe-case {0}: NO-WAIT read succeeded exit={1} (hasExited={2}) -> the probe cannot see this shape" -f $case.name, $code2, $p2.HasExited)
    } catch {
        Write-Host ("probe-case {0}: NO-WAIT read threw: {1} -> the probe DOES see this shape" -f $case.name, $_.Exception.Message)
    }
    $p2.WaitForExit()
    $p2.Dispose()
}
