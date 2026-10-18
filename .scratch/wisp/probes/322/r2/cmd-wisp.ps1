# 322-r2 头号活：`cmd/wisp` 整包那一发（★只在仓外导出台跑，⛔ 共享工作树）
# ★车道纪律：本波 `go test ./cmd/wisp` 同一时刻只许一枚腿。唯一持跑者＝228-r3c。
#   本腿⛔ 自判进度起跑——只在编排者亲发放行消息之后跑（派单 §9＋编排者 10:0x 回复①）。
# ★PATH 必须带两串 dll（third_party/sherpa-onnx ⇄ build），否则 0xc0000135＋零枚 --- FAIL
#   ＝那包根本没跑＝台面故障、⛔ 计入判据（派单颜色尺原文）。
# 已知本机一格：cmd/wisp 整包有"真窗区整进程死亡"（228-r3b 具名：c 面死在
#   TestPanelThreadIsSTAAndExitsCleanly、d 面死在 TestAC13BringUpSurvivesAReusedThreadQuit，
#   两台面不同枚＝非确定性）⇒ estate 色取自定向 -run，⛔ 拿整包 rc=1 当"322 把包打红"、⛔ 当全绿。
# 用法：pwsh -File cmd-wisp.ps1 full   |  estate   |  both
$ErrorActionPreference = 'Continue'
$Repo = 'D:\work\workspace\projects plans\Wisp'
$R = 'D:\tmp\wisp\20261018\322r2'
$Tree = (Get-Content "$R\tree-head.txt" -Raw).Trim()   # pristine 树：本腿对 cmd/wisp 零写、⛔ 用 mine 台
$Gt = "$R\gtD"
$Logs = "$Repo\.scratch\wisp\probes\322\r2\logs"
New-Item -ItemType Directory -Force -Path $Logs, $Gt | Out-Null
if (-not (Test-Path $Tree)) { Write-Error "TREE MISSING: $Tree - ABORT"; exit 91 }
$Sum = "$Logs\cmdwisp-summary.txt"
if (Test-Path $Sum) { Add-Content -Path $Sum -Value '' } else { Set-Content -Path $Sum -Value '' -Encoding ascii }
function say([string]$m) { $m | Tee-Object -FilePath $Sum -Append }
function Count4([string]$b) {
  "RUN=" + (Select-String -Path $b -Pattern '^=== RUN'  -AllMatches).Count +
  " PASS_TOP=" + (Select-String -Path $b -Pattern '^--- PASS' -AllMatches).Count +
  " FAIL_TOP=" + (Select-String -Path $b -Pattern '^--- FAIL' -AllMatches).Count +
  " SKIP=" + (Select-String -Path $b -Pattern '^--- SKIP' -AllMatches).Count
}
function ToolWords([string]$b) {
  '0xc0000135=' + (Select-String -Path $b -Pattern '0xc0000135' -AllMatches).Count +
  ' buildfailed=' + (Select-String -Path $b -Pattern 'build failed' -AllMatches).Count +
  ' noteststorun=' + (Select-String -Path $b -Pattern 'no tests to run' -AllMatches).Count +
  ' setupfailed=' + (Select-String -Path $b -Pattern 'setup failed' -AllMatches).Count +
  ' fatal-error=' + (Select-String -Path $b -Pattern '^fatal:|signal: killed|exit status' -AllMatches).Count
}
function Go([string]$label, [string]$filter, [string]$tmp) {
  $body = Join-Path $Logs "$label.log"
  Push-Location $Tree
  if ($tmp -ne 'NONE') { $env:GOTMPDIR = $tmp; New-Item -ItemType Directory -Force -Path $tmp | Out-Null }
  else { Remove-Item Env:\GOTMPDIR -ErrorAction SilentlyContinue }
  # ★两串 dll 前置（⛔ 漏＝0xc0000135＋零枚 --- FAIL＝那包根本没跑＝台面故障）
  $env:PATH = "$Tree\third_party\sherpa-onnx;$Tree\build;$env:PATH"
  $a = @('test', './cmd/wisp/', '-count=1', '-timeout', '30m', '-v')
  if ($filter -ne '') { $a += @('-run', $filter) }
  & go @a > $body 2>&1
  $rc = $LASTEXITCODE
  Pop-Location
  say "=== $label rc=$rc filter='$filter' GOTMPDIR=$(if($tmp -eq 'NONE'){'unset (默认 TMP=' + $env:TMP + ' = C 卷)'}else{$tmp}) tree=$Tree"
  say (Count4 $body)
  say (ToolWords $body)
  $pk = (Select-String -Path $body -Pattern '^(ok|FAIL|\?)\s' | ForEach-Object { $_.Line }) -join ' ;; '
  say "  package-line: $pk"
  say '  --- top FAIL names (unique):'
  $fl = (Select-String -Path $body -Pattern '^--- FAIL' | ForEach-Object { ($_.Line -replace '\s+\(.*\)\s*$', '') } | Sort-Object -Unique)
  if ($fl) { $fl | ForEach-Object { say "    $_" } } else { say '    (none)' }
  say '  --- last 6 lines verbatim (进程死亡时这里是唯一现场):'
  (Get-Content $body -Tail 6) | ForEach-Object { say "    $_" }
  return $rc
}
say "### 322-r2 cmd/wisp started $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss zzz') HEAD=$(git -C $Repo rev-parse HEAD)"
$mode = if ($args.Count -gt 0) { $args[0] } else { 'both' }
# ★-run 串逐枚现量过实名册（⛔ 指不存在的名字当凭据）：
#    Resident -> 28 枚顶层（含 228/246/256/258/260r3/r4/265/268/290/293/294/296 那一族 resident_*）
#    TestLive  -> 12 枚顶层（真窗区 live 一族）
#    TestPanelThreadIsSTAAndExitsCleanly -> 1（228-r3b 具名 c 面死亡点）
#    TestAC13BringUpSurvivesAReusedThreadQuit -> 1（228-r3b 具名 d 面死亡点）
# 尺＝git grep -h '^func Test' -- cmd/wisp | Where-Object { $_ -match <pat> } | Measure-Object
# ★成对两枚临时根台面（派单颜色尺："两把临时根台面成对（默认 C 卷 ⇄ D: 预建 GOTMPDIR）"）——
#   228-r3b 具名的整进程死亡在两台落在**不同枚**＝非确定性 ⇒ 整包那一发两台各跑一发，
#   死亡枚归位才有对照；⛔ 拿单台整包 rc=1 当"322 把包打红"、⛔ 当全绿。
switch ($mode) {
  'full' {
    Go 'cwisp-full-D' '' $Gt
    Go 'cwisp-full-C' '' 'NONE'
  }
  'estate' {
    Go 'cwisp-estate-resident' 'Resident' $Gt
    Go 'cwisp-estate-live' 'TestLive' $Gt
    Go 'cwisp-estate-deathpoints' 'TestPanelThreadIsSTAAndExitsCleanly|TestAC13BringUpSurvivesAReusedThreadQuit' $Gt
  }
  default {
    Go 'cwisp-full-D' '' $Gt
    Go 'cwisp-full-C' '' 'NONE'
    Go 'cwisp-estate-resident' 'Resident' $Gt
    Go 'cwisp-estate-live' 'TestLive' $Gt
    Go 'cwisp-estate-deathpoints' 'TestPanelThreadIsSTAAndExitsCleanly|TestAC13BringUpSurvivesAReusedThreadQuit' $Gt
  }
}
say "### done $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss zzz')"
