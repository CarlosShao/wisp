# 322-r2 头号活：`cmd/wisp` 那一发（★只在仓外导出树跑，⛔ 共享工作树）
# ★车道纪律：本波 `go test ./cmd/wisp` 同一时刻只许一枚腿；放行只由编排者发消息（派单 §9＋裁定⑤）。
#   本腿第一／二版⛔ 自判 228-r3c 进度——那一版因函数名撞车根本没跑成工具链，见下面【故障具名】。
# ★PATH 必须带两串 dll（third_party/sherpa-onnx ⇄ build），否则 0xc0000135＋零枚 --- FAIL
#   ＝那包根本没跑＝台面故障、⛔ 计入判据（派单颜色尺原文）。
# ★编排者 11:0x 裁定②⒜＝AC#4 那一格**以定向 -run 的 estate 色为准、整包色只作台面终态记录**
#   （凭据＝228-r3b/228-r3c 两程都实测"真窗区整进程死亡、resident_* 在整包从不执行"）
#   ⇒ 本台件三档定向＝estate-resident(28 枚)／estate-live(12 枚)／estate-deathpoints(那两枚具名死亡点)，
#   整包两档（D 台 ⇄ C 台）照派单"两把临时根台面成对"留作终态记录。
# ----------------------------------------------------------------------------
# 【故障件具名·前两版读数全部作废·只建不删】cwisp-full-D/C、estate-resident、estate-live、
#   estate-deathpoints 五发＝RUN=0 的零判据读数＝台面故障、⛔ 计入任何格的判据。
#   病根＝函数名 `Go` 撞上外部命令 `go`：PowerShell 函数名**大小写不敏感**，函数体里那行
#   `& go @a > $body` 递归调用了**本函数自己**（现场＝错位行 $label='test'／$filter='./cmd/wisp/'／
#   $tmp='-count=1'），根本没碰工具链 ⇒ 五发两秒"跑完"。
#   而 `0xc0000135=2` 是嵌套 summary 文本被自己数到的**自指产物**，⛔ 读成"dll 缺失"
#   （现量：build/ 与 third_party/sherpa-onnx/ 各 3 枚 dll 都在树里＝台面配方没问题）。
#   对照凭据＝同目录 host-winsec.ps1 用不撞名的 `HostRun`，那一族五发都出了真数（RUN=110 等）
#   ⇒ 病在函数名，不在台面配方、⛔ 在 cmd/wisp 的码。本版改名 `Run-Shot`。
# ----------------------------------------------------------------------------
# 用法：pwsh -File cmd-wisp.ps1 full | estate | both
$ErrorActionPreference = 'Continue'
$Repo = 'D:\work\workspace\projects plans\Wisp'
$R = 'D:\tmp\wisp\20261018\322r2'
$Tree = (Get-Content "$R\tree-head.txt" -Raw).Trim()   # pristine 树：本腿对 cmd/wisp 零写
$Gt = "$R\gtD"
$Logs = "$Repo\.scratch\wisp\probes\322\r2\logs"
New-Item -ItemType Directory -Force -Path $Logs, $Gt | Out-Null
if (-not (Test-Path $Tree)) { Write-Error "TREE MISSING: $Tree - ABORT"; exit 91 }
$Sum = "$Logs\cmdwisp-summary.txt"
Set-Content -Path $Sum -Value '' -Encoding utf8
function say([string]$m) { $m | Tee-Object -FilePath $Sum -Append }
function Count4([string]$b) {
  "RUN=" + (Select-String -Path $b -Pattern '^=== RUN' -AllMatches).Count +
  " PASS_TOP=" + (Select-String -Path $b -Pattern '^--- PASS' -AllMatches).Count +
  " FAIL_TOP=" + (Select-String -Path $b -Pattern '^--- FAIL' -AllMatches).Count +
  " SKIP=" + (Select-String -Path $b -Pattern '^--- SKIP' -AllMatches).Count
}
function ToolWords([string]$b) {
  '0xc0000135=' + (Select-String -Path $b -Pattern '0xc0000135' -AllMatches).Count +
  ' buildfailed=' + (Select-String -Path $b -Pattern 'build failed' -AllMatches).Count +
  ' noteststorun=' + (Select-String -Path $b -Pattern 'no tests to run' -AllMatches).Count +
  ' setupfailed=' + (Select-String -Path $b -Pattern 'setup failed' -AllMatches).Count +
  ' dead=' + (Select-String -Path $b -Pattern 'signal: killed|exit status|panic:' -AllMatches).Count
}
function Run-Shot([string]$label, [string]$filter, [string]$tmp) {
  $body = Join-Path $Logs "$label.log"
  Push-Location $Tree
  if ($tmp -ne 'NONE') { $env:GOTMPDIR = $tmp; New-Item -ItemType Directory -Force -Path $tmp | Out-Null }
  else { Remove-Item Env:\GOTMPDIR -ErrorAction SilentlyContinue }
  $env:PATH = "$Tree\third_party\sherpa-onnx;$Tree\build;$env:PATH"
  $ga = @('test', './cmd/wisp/', '-count=1', '-timeout', '30m', '-v')
  if ($filter -ne '') { $ga += @('-run', $filter) }
  # ★外部命令显式到底，⛔ 再让任何同名函数有机会截走这一行
  & 'go' @ga > $body 2>&1
  $rc = $LASTEXITCODE
  Pop-Location
  $tg = if ($tmp -eq 'NONE') { 'unset (default TMP=' + $env:TMP + ' = C drive)' } else { $tmp }
  say "=== $label rc=$rc filter='$filter' GOTMPDIR=$tg tree=$Tree"
  say (Count4 $body)
  say (ToolWords $body)
  $pk = (Select-String -Path $body -Pattern '^(ok|FAIL|\?)\s' | ForEach-Object { $_.Line }) -join ' ;; '
  say "  package-line: $pk"
  say '  --- top FAIL names (unique):'
  $fl = (Select-String -Path $body -Pattern '^--- FAIL' | ForEach-Object { ($_.Line -replace '\s+\(.*\)\s*$', '') } | Sort-Object -Unique)
  if ($fl) { $fl | ForEach-Object { say "    $_" } } else { say '    (none)' }
  say '  --- top SKIP names (unique):'
  $sk = (Select-String -Path $body -Pattern '^--- SKIP' | ForEach-Object { ($_.Line -replace '\s+\(.*\)\s*$', '') } | Sort-Object -Unique)
  if ($sk) { $sk | ForEach-Object { say "    $_" } } else { say '    (none)' }
  say '  --- last 8 lines verbatim (整进程死亡时这里是唯一现场):'
  (Get-Content $body -Tail 8) | ForEach-Object { say "    $_" }
}
say "### 322-r2 cmd/wisp started $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss zzz') HEAD=$(git -C $Repo rev-parse HEAD)"
$mode = if ($args.Count -gt 0) { $args[0] } else { 'both' }
switch ($mode) {
  'full' {
    Run-Shot 'cwisp-full-D' '' $Gt
    Run-Shot 'cwisp-full-C' '' 'NONE'
  }
  'estate' {
    Run-Shot 'cwisp-estate-resident' 'Resident' $Gt
    Run-Shot 'cwisp-estate-live' 'TestLive' $Gt
    Run-Shot 'cwisp-estate-deathpoints' 'TestPanelThreadIsSTAAndExitsCleanly|TestAC13BringUpSurvivesAReusedThreadQuit' $Gt
  }
  default {
    Run-Shot 'cwisp-full-D' '' $Gt
    Run-Shot 'cwisp-full-C' '' 'NONE'
    Run-Shot 'cwisp-estate-resident' 'Resident' $Gt
    Run-Shot 'cwisp-estate-live' 'TestLive' $Gt
    Run-Shot 'cwisp-estate-deathpoints' 'TestPanelThreadIsSTAAndExitsCleanly|TestAC13BringUpSurvivesAReusedThreadQuit' $Gt
  }
}
say "### done $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss zzz')"
