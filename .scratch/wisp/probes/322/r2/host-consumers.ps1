# 322-r2 本机（Windows 档）消费者五包读数：internal/risk/secret/config/agent/memory。
# 目的＝AC#4 消费者维的"本机 ⇄ Linux 档"成对读数的本机那一半。
# ★⛔ cmd/wisp（本波唯一持跑者＝228-r3c，等编排者放行）。
$ErrorActionPreference = 'Continue'
$Repo = 'D:\work\workspace\projects plans\Wisp'
$R = 'D:\tmp\wisp\20261018\322r2'
$Tree = (Get-Content "$R\tree-head.txt" -Raw).Trim()
$Gt = "$R\gtD"
$Logs = "$Repo\.scratch\wisp\probes\322\r2\logs"
New-Item -ItemType Directory -Force -Path $Gt | Out-Null
$Sum = "$Logs\host-consumers-summary.txt"
Set-Content -Path $Sum -Value '' -Encoding ascii
function say([string]$m) { $m | Tee-Object -FilePath $Sum -Append }
$env:GOTMPDIR = $Gt
$env:PATH = "$Tree\third_party\sherpa-onnx;$Tree\build;$env:PATH"
say "### 322-r2 host consumers started $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss zzz')"
say "TREE=$Tree (pristine 字节，本腿那一枚 .go 改前) GOTMPDIR=$Gt"
Push-Location $Tree
# 【第一版故障具名】原来用 $p.Trim('/').Replace('internal/','') 取 tag：'./internal/risk/' 的
# Trim('/') 被开头的 '.' 挡住 ⇒ tag 变成 './risk.' ⇒ 正文写进不存在的 'hc-\risk.log'、
# 五包全部零读数＝**台面故障、⛔ 计入判据**（作废的 summary 只建不删）。本版显式配对。
$pairs = @(
  @{ pkg = './internal/risk/';    tag = 'risk' },
  @{ pkg = './internal/secret/';  tag = 'secret' },
  @{ pkg = './internal/config/';  tag = 'config' },
  @{ pkg = './internal/agent/';   tag = 'agent' },
  @{ pkg = './internal/memory/';  tag = 'memory' }
)
foreach ($one in $pairs) {
  $p = $one.pkg; $tag = $one.tag
  $body = Join-Path $Logs "hc-$tag.log"
  & go test $p -count=1 -timeout 30m -v > $body 2>&1
  $rc = $LASTEXITCODE
  $run = (Select-String -Path $body -Pattern '^=== RUN' -AllMatches).Count
  $pa = (Select-String -Path $body -Pattern '^--- PASS' -AllMatches).Count
  $fa = (Select-String -Path $body -Pattern '^--- FAIL' -AllMatches).Count
  $sk = (Select-String -Path $body -Pattern '^--- SKIP' -AllMatches).Count
  $pkgline = (Select-String -Path $body -Pattern '^(ok|FAIL|\?)\s' | ForEach-Object { $_.Line }) -join ' ;; '
  say "PKG $p rc=$rc RUN=$run PASS_TOP=$pa FAIL_TOP=$fa SKIP=$sk :: $pkgline"
  $fl = (Select-String -Path $body -Pattern '^--- FAIL' | ForEach-Object { $_.Line })
  if ($fl) { say '  FAIL names:'; $fl | ForEach-Object { say "    $_" } }
  say ('  tooling-words: walk=' + (Select-String -Path $body -Pattern 'walk' -AllMatches).Count +
      ' buildfailed=' + (Select-String -Path $body -Pattern 'build failed' -AllMatches).Count +
      ' noteststorun=' + (Select-String -Path $body -Pattern 'no tests to run' -AllMatches).Count +
      ' setupfailed=' + (Select-String -Path $body -Pattern 'setup failed' -AllMatches).Count +
      ' 0xc0000135=' + (Select-String -Path $body -Pattern '0xc0000135' -AllMatches).Count)
}
Pop-Location
say "### done $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss zzz')"
