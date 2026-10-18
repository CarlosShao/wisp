# 322-r2 本机 internal/winsec 取证（原生 PowerShell，⛔ MSYS 路径形，避免台面故障）。
# 机主 2026-10-18 10:1x 新铁律：盘符根零污染，唯一允许落点 D:\tmp\wisp\20261018\322r2\*。
# ⛔ 在共享工作树跑 go test —— 全部在仓外导出树 $Tree 里。
# 成对两枚临时根台面：gtD（D: 卷·预建·有 8.3 别名）⇄ NONE（默认 TMP＝C 卷）。
$ErrorActionPreference = 'Continue'
$Repo = 'D:\work\workspace\projects plans\Wisp'
$Root = 'D:\tmp\wisp\20261018\322r2'
$Tree = (Get-Content "$Root\tree-head.txt" -Raw).Trim()
$Gt   = "$Root\gtD"
$Logs = "$Repo\.scratch\wisp\probes\322\r2\logs"
New-Item -ItemType Directory -Force -Path $Logs | Out-Null
New-Item -ItemType Directory -Force -Path $Gt | Out-Null
$Head = (git -C $Repo rev-parse HEAD)
$Sum = "$Logs\host-winsec-summary.txt"
Set-Content -Path $Sum -Value '' -Encoding ascii
function say([string]$m) { $m | Tee-Object -FilePath $Sum -Append }

say "### 322-r2 host winsec started $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss zzz') HEAD=$Head"
say "TREE=$Tree"
say "GOTMP-D=$Gt  (D: 卷·预建·8.3 别名)"
# cmd/wisp 才需要的两串 dll；对本包无害，照配方带上（派单★色尺）
$env:PATH = "$Tree\third_party\sherpa-onnx;$Tree\build;$env:PATH"
say "PATH-HEAD=$($env:PATH.Substring(0,[Math]::Min(170,$env:PATH.Length)))"

function Four([string]$b) {
  $run  = (Select-String -Path $b -Pattern '^=== RUN'   -AllMatches).Count
  $pass = (Select-String -Path $b -Pattern '^--- PASS'  -AllMatches).Count
  $fail = (Select-String -Path $b -Pattern '^--- FAIL'  -AllMatches).Count
  $skip = (Select-String -Path $b -Pattern '^--- SKIP'  -AllMatches).Count
  return "RUN=$run PASS_TOP=$pass FAIL_TOP=$fail SKIP=$skip"
}
function ToolWords([string]$b) {
  $t = @('walk','build failed','no tests to run','setup failed','0xc0000135')
  $out = foreach ($w in $t) { "$($w -replace ' ','')=$((Select-String -Path $b -Pattern ([regex]::Escape($w)) -AllMatches).Count)" }
  return ($out -join ' ')
}
function HostRun([string]$label, [string]$tmp, [string]$filter) {
  $body = Join-Path $Logs "$label.log"
  Push-Location $Tree
  if ($tmp -ne 'NONE') { $env:GOTMPDIR = $tmp; New-Item -ItemType Directory -Force -Path $tmp | Out-Null }
  else { Remove-Item Env:\GOTMPDIR -ErrorAction SilentlyContinue }
  $args = @('test','./internal/winsec/','-count=1','-timeout','30m','-v')
  if ($filter -ne '') { $args += @('-run', $filter) }
  & go @args > $body 2>&1
  $rc = $LASTEXITCODE
  Pop-Location
  say "=== $label rc=$rc GOTMPDIR=$(if($tmp -eq 'NONE'){'unset (默认 TMP=' + $env:TMP + ' = C 卷)'}else{$tmp}) filter='$filter'"
  say (Four $body)
  say (ToolWords $body)
  $fl = (Select-String -Path $body -Pattern '^--- FAIL' | ForEach-Object { $_.Line })
  if ($fl) { say '--- top FAIL names:'; $fl | ForEach-Object { say "  $_" } } else { say '--- top FAIL names: (none)' }
  $sk = (Select-String -Path $body -Pattern '^--- SKIP' | ForEach-Object { $_.Line })
  if ($sk) { say '--- top SKIP names:'; $sk | ForEach-Object { say "  $_" } } else { say '--- top SKIP names: (none)' }
  $wt = (Select-String -Path $body -Pattern 'witnessTruncated' | ForEach-Object { $_.Line })
  if ($wt) { say '--- witnessTruncated rows:'; $wt | ForEach-Object { say "  $_" } } else { say '--- witnessTruncated rows: (none)' }
}

$PatEstate = 'TestGateJudgesThePrivateSetByResolvedSID|TestSealNarrowsAndNamesThePrincipalItRemovedBySID'
$step = if ($args.Count -gt 0) { $args[0] } else { 'all' }
switch ($step) {
  'a' { HostRun 'wh-D' $Gt '' }
  'b' { HostRun 'wh-C' 'NONE' '' }
  'c' { HostRun 'wd-D' $Gt $PatEstate }
  'd' { HostRun 'wd-C' 'NONE' $PatEstate }
  'e' { HostRun 'wt-D' $Gt 'TestAC3Windows' }
  default {
    HostRun 'wh-D' $Gt ''
    HostRun 'wh-C' 'NONE' ''
    HostRun 'wd-D' $Gt $PatEstate
    HostRun 'wd-C' 'NONE' $PatEstate
    HostRun 'wt-D' $Gt 'TestAC3Windows'
  }
}
say "### done $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss zzz')"
