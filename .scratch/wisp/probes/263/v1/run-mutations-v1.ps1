# 263-v1 突变台件：从 tracked 字节逐字节生成副本，锚不命中即 throw；⛔ 不在跟踪文件上做瞬时破坏
# 用法：powershell -NoProfile -ExecutionPolicy Bypass -File run-mutations-v1.ps1
$ErrorActionPreference = 'Stop'
$here = Split-Path -Parent $PSCommandPath
$base = Join-Path $here 'mutations\base.ps1'          # 由 git cat-file blob HEAD: 抽出的字节件
$tracked = 'D:\work\workspace\projects plans\Wisp\scripts\slo-check.ps1'
$md5Before = (Get-FileHash -Algorithm MD5 -LiteralPath $tracked).Hash

function New-Mutation([string]$name, [array]$edits) {
    # $edits: 每枚 @{Kind='text'|'regex'; Anchor=...; Replacement=...}
    $text = [IO.File]::ReadAllText($base)
    foreach ($e in $edits) {
        if ($e.Kind -eq 'text') {
            if (-not $text.Contains($e.Anchor)) { throw "anchor miss ($name): $($e.Anchor)" }
            $text = $text.Replace($e.Anchor, $e.Replacement)
        } else {
            if (-not [regex]::IsMatch($text, $e.Anchor)) { throw "regex anchor miss ($name): $($e.Anchor)" }
            $text = [regex]::Replace($text, $e.Anchor, $e.Replacement)
        }
    }
    $out = Join-Path $here "mutations\slo-check-$name.ps1"
    [IO.File]::WriteAllText($out, $text, (New-Object System.Text.UTF8Encoding($false)))
    return $out
}

function Invoke-Mutation([string]$name, [string]$file, [string]$env) {
    $outDir = Join-Path $here "mutations\out-$name"
    if (-not (Test-Path $outDir)) { New-Item -ItemType Directory -Path $outDir | Out-Null }
    $exe = Join-Path $here 'bin\probe263v1-gui.exe'
    $log = Join-Path $here "mutations\log-$name.txt"
    $raw = Join-Path $here "mutations\raw-$name.txt"
    $rawErr = Join-Path $here "mutations\raw-$name.err.txt"
    if ($env -eq 'fail') { $env:V1_VERDICT = 'fail' } else { Remove-Item Env:V1_VERDICT -ErrorAction SilentlyContinue }
    # 子进程 stderr 一律重定向到文件：PS 5.1 里把外部程序的 stderr 并进管道（2>&1）会在
    # `$ErrorActionPreference='Stop'` 下变成 NativeCommandError 把台件自己打死（本腿踩过一次）；
    # 退码用 Start-Process -PassThru -Wait 取，不依赖 $LASTEXITCODE（本腿也撞上它不落）。
    $argLine = '-NoProfile -ExecutionPolicy Bypass -File "' + $file + '" -Subset smoke -SecondsPerState 1 -WispExe "' + $exe + '" -OutDir "' + $outDir + '"'
    # 起-等-取三件事自己做（与门内 Invoke-ExternalProgram 同形）：Start-Process 与 1>/2> 两把尺
    # 都在本腿上取不到退码，已具名记在证据件 §⑨。
    $psi = New-Object System.Diagnostics.ProcessStartInfo
    $psi.FileName = (Get-Command powershell.exe -CommandType Application | Select-Object -First 1).Source
    $psi.Arguments = $argLine
    $psi.UseShellExecute = $false
    $psi.RedirectStandardOutput = $true
    $psi.RedirectStandardError = $true
    $child = [System.Diagnostics.Process]::Start($psi)
    $so = $child.StandardOutput.ReadToEnd()
    $se = $child.StandardError.ReadToEnd()
    $child.WaitForExit()
    $rc = $child.ExitCode
    $child.Dispose()
    Set-Content -LiteralPath $raw -Value $so -Encoding UTF8
    Set-Content -LiteralPath $rawErr -Value $se -Encoding UTF8
    Remove-Item Env:V1_VERDICT -ErrorAction SilentlyContinue
    $lines = @(Get-Content -LiteralPath $raw -ErrorAction SilentlyContinue) + @(Get-Content -LiteralPath $rawErr -ErrorAction SilentlyContinue)
    @("=== MUTATION $name rc=$rc ===") + @($lines | ForEach-Object { [string]$_ }) | Set-Content -Path $log -Encoding UTF8
    return [pscustomobject]@{ name = $name; rc = $rc; log = $log; lines = @($lines) }
}

$reports = @()

# 正控甲：删掉等待 —— 能力钉应当当场红
$f = New-Mutation 'drop-wait' @(@{Kind='text'; Anchor='        $proc.WaitForExit()'; Replacement='        # WAIT-REMOVED-BY-263-V1'})
$reports += Invoke-Mutation 'drop-wait' $f ''

# 正控乙：种一枚裸读 —— 词面钉应当当场红
$f = New-Mutation 'plant-bare' @(@{Kind='text'; Anchor='    $code = $run.code'; Replacement='    $code = $LASTEXITCODE'})
$reports += Invoke-Mutation 'plant-bare' $f ''

# 绕过形一：运行期取名（Get-Variable -ValueOnly）+ 真的改用 &
$plant = "    & `$WispExe slo -state `$state -seconds `$SecondsPerState -interval-ms 250 -out `$outFile`n    `$code = (Get-Variable -Name 'LASTEXITCODE' -ValueOnly)"
$f = New-Mutation 'plant-getvariable' @(@{Kind='text'; Anchor='    $code = $run.code'; Replacement=$plant})
$reports += Invoke-Mutation 'plant-getvariable' $f ''

# 绕过形二：把两枚钉的调用删掉（其余一字不动）
$f = New-Mutation 'delete-nails' @(
    @{Kind='regex'; Anchor='(?m)^Test-ScriptShapeNail$'; Replacement='# NAIL-CALL-REMOVED-BY-263-V1'},
    @{Kind='regex'; Anchor='(?m)^Test-ExitCodeInstrument$'; Replacement='# NAIL-CALL-REMOVED-BY-263-V1'}
)
$reports += Invoke-Mutation 'delete-nails' $f ''

# 绕过形三：读取挪到别的文件（本目录另有一枚 helper，由本腿单独写好、⛔ 不在这里用 here-string 生成）
$helper = Join-Path $here 'mutations\helper-v1.ps1'
if (-not (Test-Path -LiteralPath $helper)) { throw "helper missing: $helper" }
$plant3 = "    . (Join-Path `$PSScriptRoot 'helper-v1.ps1')`n    `$code = Get-V1ExitCode -Exe `$WispExe -State `$state -Seconds `$SecondsPerState -Out `$outFile"
$f = New-Mutation 'read-in-other-file' @(@{Kind='text'; Anchor='    $code = $run.code'; Replacement=$plant3})
$reports += Invoke-Mutation 'read-in-other-file' $f ''

# 判据那一行的活路（AC#3 ④）：把某档判据换成必红形状，样本本身仍达标
$f = New-Mutation 'force-state-fail' @(@{Kind='text'; Anchor='            $pass = [bool]$reportJson.pass'; Replacement='            $pass = $false # FORCE-RED-BY-263-V1'})
$reports += Invoke-Mutation 'force-state-fail' $f ''

# 副本机制自身的对照：同一套生成路径、零改动 —— 应当与盘上交付字节同色（绿）
$noop = Invoke-Mutation 'copy-noop' @(@{Kind='text'; Anchor='    $code = $run.code'; Replacement='    $code = $run.code'})
$reports += $noop

Write-Host '================ 汇总 ================'
foreach ($r in $reports) { Write-Host ("{0,-24} rc={1}" -f $r.name, $r.rc) }
$md5After = (Get-FileHash -Algorithm MD5 -LiteralPath $tracked).Hash
Write-Host ("tracked md5 before={0}" -f $md5Before)
Write-Host ("tracked md5 after ={0}" -f $md5After)
Write-Host ("tracked UNCHANGED = {0}" -f ($md5Before -eq $md5After))
