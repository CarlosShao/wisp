# 322-r2 容器（golang:1.27 linux）取证：internal/winsec 整包 -count=2 -v（容器四数）
# ＋ AC#4 消费者面 internal/risk/secret/config/agent/memory 的 Linux 档成对读数。
# ⛔ 在共享工作树跑 test。GOMODCACHE＝宿主 /d/work/base/gopath/pkg/mod 只读（⛔ 命名卷 ro 坑，见 r1b §6 自纠）。
# 从 PowerShell 直接调 docker ⇒ 挂载路径用 D:/ 形（MSYS 形 /d/... 是给 sh 用的，r1b/v1 那两版台件在 sh 下有效）。
# 机主 2026-10-18 10:1x 新铁律：落点唯一允许 D:\tmp\wisp\20261018\322r2\*，盘符根零污染。
$ErrorActionPreference = 'Continue'
$Repo = 'D:\work\workspace\projects plans\Wisp'
$Root = 'D:\tmp\wisp\20261018\322r2'
$Tree = (Get-Content "$Root\tree-head.txt" -Raw).Trim()
$TreeWin = $Tree -replace '\\', '/'
$Logs = "$Repo\.scratch\wisp\probes\322\r2\logs"
New-Item -ItemType Directory -Force -Path $Logs | Out-Null
$ModCacheWin = 'D:/work/base/gopath/pkg/mod'
$Head = (git -C $Repo rev-parse HEAD)
$Sum = "$Logs\container-summary.txt"
Set-Content -Path $Sum -Value '' -Encoding ascii
function say([string]$m) { $m | Tee-Object -FilePath $Sum -Append }

say "### 322-r2 container started $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss zzz') HEAD=$Head"
say "TREE=$TreeWin  MODCACHE=$ModCacheWin"
say "GO-VERSION-IMAGE: $(docker run --rm golang:1.27 go version)"

function Four([string]$b) {
  $run  = (Select-String -Path $b -Pattern '^=== RUN'   -AllMatches).Count
  $pass = (Select-String -Path $b -Pattern '^--- PASS'  -AllMatches).Count
  $fail = (Select-String -Path $b -Pattern '^--- FAIL'  -AllMatches).Count
  $skip = (Select-String -Path $b -Pattern '^--- SKIP'  -AllMatches).Count
  return "RUN=$run PASS_TOP=$pass FAIL_TOP=$fail SKIP=$skip"
}

function CRun([string]$label, [string]$pkgs, [string]$runfilter) {
  $inner = "export GOMODCACHE=/gomod GOTMPDIR=/tmp/gt-$label; mkdir -p `"`$GOTMPDIR`"; go test $pkgs -count=2 -timeout 30m -v"
  if ($runfilter -ne '') { $inner += " -run `"$runfilter`"" }
  $inner += " > /work/$label.log 2>&1; rc=`$?; echo `"go-test-rc=`$rc`" >> /work/$label.log; exit `$rc"
  $rc = 0
  $out = & docker run --rm -e CGO_ENABLED=0 -e GOPROXY=off -e GOFLAGS=-mod=mod `
      -v "${TreeWin}:/work" -v probe322r2build:/root/.cache/go-build `
      -v "${ModCacheWin}:/gomod:ro" -w /work golang:1.27 `
      sh -c $inner 2>&1
  $rc = $LASTEXITCODE
  $body = Join-Path $Tree "$label.log"
  say "=== $label pkgs='$pkgs' runfilter='$runfilter' host-rc=$rc"
  if ($out) { say '--- docker stderr/stdout (first 8 lines):'; ($out | Select-Object -First 8) | ForEach-Object { say "  $_" } }
  if (Test-Path $body) {
    Copy-Item $body "$Logs\$label.body.log" -Force
    $b = "$Logs\$label.body.log"
    say (Four $b)
    say '--- package result lines:'
    (Select-String -Path $b -Pattern '^(ok|\?|FAIL)\s' | ForEach-Object { $_.Line }) | ForEach-Object { say "  $_" }
    say '--- top FAIL names:'
    $fl = (Select-String -Path $b -Pattern '^--- FAIL' | ForEach-Object { $_.Line })
    if ($fl) { $fl | ForEach-Object { say "  $_" } } else { say '  (none)' }
    say '--- RED lines verbatim:'
    $red = (Select-String -Path $b -Pattern 'RED|paired reading' | ForEach-Object { $_.Line })
    if ($red) { $red | ForEach-Object { say "  $_" } } else { say '  (no RED line)' }
    say ('tooling-words: walk=' + (Select-String -Path $b -Pattern 'walk' -AllMatches).Count +
         ' buildfailed=' + (Select-String -Path $b -Pattern 'build failed' -AllMatches).Count +
         ' noteststorun=' + (Select-String -Path $b -Pattern 'no tests to run' -AllMatches).Count +
         ' setupfailed=' + (Select-String -Path $b -Pattern 'setup failed' -AllMatches).Count)
  } else {
    say "  NO BODY at $body ⇒ 台面故障，⛔ 计入判据"
  }
}

$step = if ($args.Count -gt 0) { $args[0] } else { 'all' }
switch ($step) {
  'w' { CRun 'cw-winsec' './internal/winsec/' '' }
  'c' { CRun 'cc-consumers' './internal/risk/ ./internal/secret/ ./internal/config/ ./internal/agent/ ./internal/memory/' '' }
  default {
    CRun 'cw-winsec' './internal/winsec/' ''
    CRun 'cc-consumers' './internal/risk/ ./internal/secret/ ./internal/config/ ./internal/agent/ ./internal/memory/' ''
  }
}
say "### done $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss zzz')"
