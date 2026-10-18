# 322-r2 容器成对复跑：internal/winsec 整包 -count=2 -v，两枚台成对——
#   cw2-head  = 起树锚 pristine 字节（sealhandle_privateroute_120_other_test.go = fb6c54ed…）
#   cw2-mine  = 同一棵树＋我那一枚新字节（f87035bd…）
# 目的＝证明「只改注释行＋一枚永不执行的 t.Logf 字符串」不改变任何颜色。
# ★这一枚文件是 `//go:build !windows` ⇒ 本机 Windows 档根本不编译它，
#   所以它的射程只在 POSIX/容器这一侧；本机那些发是 pristine 字节下的读数（具名）。
$ErrorActionPreference = 'Continue'
$Repo = 'D:\work\workspace\projects plans\Wisp'
$R = 'D:\tmp\wisp\20261018\322r2'
$Logs = "$Repo\.scratch\wisp\probes\322\r2\logs"
$TreeHead = (Get-Content "$R\tree-head.txt" -Raw).Trim()
$TreeMine = (Get-Content "$R\tree-mine.txt" -Raw).Trim()
# 同 gates.ps1 的坑：MSYS 形/反斜杠形都要认，并且**硬验存在性**，路径不在就 ABORT，
# ⛔ 让"容器里没跑成"冒充成四数全 0 的绿。
if ($TreeMine -match '^/d/') { $TreeMine = 'D:\' + ($TreeMine.Substring(3) -replace '/', '\') }
if ($TreeMine -match '^D:/') { $TreeMine = $TreeMine -replace '/', '\' }
if (-not (Test-Path $TreeMine)) { Write-Error "MINE TREE MISSING: $TreeMine - ABORT"; exit 91 }
if (-not (Test-Path $TreeHead)) { Write-Error "HEAD TREE MISSING: $TreeHead - ABORT"; exit 91 }
$ModCacheWin = 'D:/work/base/gopath/pkg/mod'
$Sum = "$Logs\container-mine-summary.txt"
Set-Content -Path $Sum -Value '' -Encoding ascii
function say([string]$m) { $m | Tee-Object -FilePath $Sum -Append }

function CRun([string]$label, [string]$treeWin) {
  $inner = 'export GOMODCACHE=/gomod GOTMPDIR=/tmp/gt-' + $label + '; mkdir -p "$GOTMPDIR"; ' +
    'go test ./internal/winsec/ -count=2 -timeout 30m -v > /work/' + $label + '.log 2>&1; rc=$?; echo "go-test-rc=$rc" >> /work/' + $label + '.log; exit $rc'
  $tw = $treeWin -replace '\\', '/'
  $out = & docker run --rm -e CGO_ENABLED=0 -e GOPROXY=off -e GOFLAGS=-mod=mod `
      -v "${tw}:/work" -v probe322r2build:/root/.cache/go-build `
      -v "${ModCacheWin}:/gomod:ro" -w /work golang:1.27 sh -c $inner 2>&1
  $rc = $LASTEXITCODE
  $body = Join-Path $treeWin "$label.log"
  if (Test-Path $body) {
    Copy-Item $body "$Logs\$label.body.log" -Force
    $b = "$Logs\$label.body.log"
    say "=== $label tree=$tw rc=$rc"
    say ("RUN=" + (Select-String -Path $b -Pattern '^=== RUN' -AllMatches).Count +
        " PASS_TOP=" + (Select-String -Path $b -Pattern '^--- PASS' -AllMatches).Count +
        " FAIL_TOP=" + (Select-String -Path $b -Pattern '^--- FAIL' -AllMatches).Count +
        " SKIP=" + (Select-String -Path $b -Pattern '^--- SKIP' -AllMatches).Count)
    say '--- package result lines:'
    (Select-String -Path $b -Pattern '^(ok|\?|FAIL)\s' | ForEach-Object { $_.Line }) | ForEach-Object { say "  $_" }
    say '--- top FAIL names:'
    $fl = (Select-String -Path $b -Pattern '^--- FAIL' | ForEach-Object { $_.Line })
    if ($fl) { $fl | ForEach-Object { say "  $_" } } else { say '  (none)' }
    say '--- AC#1 paired readings (openleaf 322 race rows):'
    (Select-String -Path $b -Pattern 'AC#1 paired reading on one screen' | ForEach-Object { $_.Line.Trim() }) | ForEach-Object { say "  $_" }
    say ('--- branch尺: holds-its-own=' + (Select-String -Path $b -Pattern 'witness holds its own bytes').Count +
        ' R-322-SHAPE=' + (Select-String -Path $b -Pattern "R-322 SHAPE OBSERVED").Count +
        ' payload-IN=' + (Select-String -Path $b -Pattern 'raced payload IN the foreign').Count)
  } else {
    say "=== $label NO BODY at $body (rc=$rc) ⇒ 台面故障"
    if ($out) { ($out | Select-Object -First 6) | ForEach-Object { say "  $_" } }
  }
}

say "### 322-r2 container paired started $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss zzz')"
$step = if ($args.Count -gt 0) { $args[0] } else { 'all' }
switch ($step) {
  'h' { CRun 'cw2-head' $TreeHead }
  'm' { CRun 'cw2-mine' $TreeMine }
  default { CRun 'cw2-head' $TreeHead; CRun 'cw2-mine' $TreeMine }
}
say "### done $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss zzz')"
