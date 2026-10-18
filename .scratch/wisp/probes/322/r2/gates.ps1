# 322-r2 门禁（第二版·修 & 调用形）。
# 分工照派单：格式化/d22scan/长度门＝**共享工作树**；go build/go vet（含跨档）＝**仓外导出台**。
# ★两枚台：head 树（pristine，⛔ 我的字节）⇄ mine 树（＝head＋我的新字节）——差色只可能来自我那一枚文件。
$ErrorActionPreference = 'Continue'
$Repo = 'D:\work\workspace\projects plans\Wisp'
$R    = 'D:\tmp\wisp\20261018\322r2'
$Tree  = (Get-Content "$R\tree-head.txt" -Raw).Trim()
$Mine  = (Get-Content "$R\tree-mine.txt" -Raw).Trim()
# 上一版这里用 `if ($Mine -match '^/')` 把 MSYS 形改成 'D:'+'\...'，结果造出 D:\d\tmp\... 这种
# 不存在的路径 ⇒ 那一版 6 枚 *-mine 发全部 Push-Location 失败、rc=0 是**无效读数**（已具名作废）。
# 本版显式两形都认，并加存在性硬闸：路径不在 ⇒ 当场 ABORT，⛔ 让假绿进判据。
if ($Mine -match '^/d/') { $Mine = 'D:\' + ($Mine.Substring(3) -replace '/', '\') }
if ($Mine -match '^D:/') { $Mine = $Mine -replace '/', '\' }
if (-not (Test-Path $Mine)) { Write-Error "MINE TREE MISSING: $Mine - ABORT"; exit 91 }
if (-not (Test-Path $Tree)) { Write-Error "HEAD TREE MISSING: $Tree - ABORT"; exit 91 }
$Logs = "$Repo\.scratch\wisp\probes\322\r2\logs"
$Sum  = "$Logs\gates-summary.txt"
Set-Content -Path $Sum -Value '' -Encoding ascii
function say([string]$m) { $m | Tee-Object -FilePath $Sum -Append }

function Gate([string]$label, [scriptblock]$body, [string]$wd, [string]$env1) {
  Push-Location $wd
  if ($env1) { $k, $v = $env1 -split '=', 2; Set-Item -Path ("env:" + $k) -Value $v }
  $out = & $body 2>&1
  $rc = $LASTEXITCODE
  if ($env1) { $k = ($env1 -split '=', 2)[0]; Remove-Item ("env:" + $k) -ErrorAction SilentlyContinue }
  Pop-Location
  $joined = ($out | Out-String)
  Set-Content -Path "$Logs\gate-$label.txt" -Value $joined -Encoding utf8
  $tail = (($out | Where-Object { "$_".Trim() -ne '' } | Select-Object -Last 4) -join ' | ')
  say "GATE $label rc=$rc wd=$wd :: $tail"
}

Set-Location $Repo
say "### 322-r2 gates v2 started $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss zzz')"
say "MY BYTE = internal/winsec/sealhandle_privateroute_120_other_test.go"
say "TREE(head,pristine)=$Tree"
say "TREE(mine=head+mybyte)=$Mine"

# --- 格式化尺：只在共享工作树，只对我动的件 ---
Push-Location $Repo
$gf = & 'D:\work\base\gopath\bin\gofumpt.exe' -l 'internal/winsec/sealhandle_privateroute_120_other_test.go' 2>&1
$grc = $LASTEXITCODE
Pop-Location
say "GATE gofumpt-mine rc=$grc listed=$(if([string]::IsNullOrWhiteSpace(($gf|Out-String))){0}else{@($gf).Count}) :: $(($gf) -join ';')"
Push-Location $Repo
$gm = & gofmt -l 'internal/winsec/sealhandle_privateroute_120_other_test.go' 2>&1
$gmrc = $LASTEXITCODE
Pop-Location
say "GATE gofmt-mine rc=$gmrc listed=$(if([string]::IsNullOrWhiteSpace(($gm|Out-String))){0}else{@($gm).Count})"

# --- d22scan / 长度门：只在共享工作树（需要 design/ 与 .git） ---
Gate 'd22scan'      { sh scripts/d22scan.sh } $Repo $null
Gate 'lengthgate'   { sh scripts/check-path-length-budget.sh --with-self-test } $Repo $null

# --- build / vet：仓外导出台，两枚台成对 ---
Gate 'build-head'   { go build ./... } $Tree $null
Gate 'build-mine'   { go build ./... } $Mine $null
Gate 'vet-winsec-head' { go vet ./internal/winsec/ } $Tree $null
Gate 'vet-winsec-mine' { go vet ./internal/winsec/ } $Mine $null
Gate 'vet-cmdwisp-head' { go vet ./cmd/wisp/ } $Tree $null
Gate 'vet-cmdwisp-mine' { go vet ./cmd/wisp/ } $Mine $null
foreach ($os in 'linux', 'darwin', 'windows') {
  Gate "vet-$os-winsec-head" { go vet ./internal/winsec/ } $Tree "GOOS=$os"
  Gate "vet-$os-winsec-mine" { go vet ./internal/winsec/ } $Mine "GOOS=$os"
}
say "### done $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss zzz')"
