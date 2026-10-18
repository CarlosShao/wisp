# 339-r1 AC#2 去牙反控（Windows 台面＝牙齿新家）。两发串行，每发前后验 pristine 哈希。
# 正控＝未突变 blacklist.go ⇒ 334 牙齿绿；去牙＝摘掉 aRelates 尾行回落支 ⇒ 334 牙齿点名红。
# ⛔ 动跟踪文件：突变只打在仓外导出树的物理字节上；逐发还原并报还原哈希。
$ErrorActionPreference = 'Continue'
$REPO = 'D:\work\workspace\projects plans\Wisp'
$TREE = 'D:\tmp\wisp\20261018\339r1\win\after'
$BL = Join-Path $TREE 'internal\risk\blacklist.go'
$GT = 'D:\tmp\wisp\20261018\339r1\gt'
$PRISTINE = 'c1799cd24afce19e6accef23db23a7d7ee746dec'
$RUN = 'TestAnchorMissReasonNamesCandidateAndMissingAnchor|TestAnchorMissKeepsNonOverridableClassAndNeedL2|TestAnchorFallbackTeethPerAnchor336'
$env:TEMP = $GT; $env:TMP = $GT; $env:GOTMPDIR = $GT
$env:GOFLAGS = '-mod=mod'; $env:GOPROXY = 'off'; $env:CGO_ENABLED = '0'
$env:GOCACHE = 'D:\tmp\wisp\20261018\339r1\gocache'
New-Item -ItemType Directory -Force -Path $GT | Out-Null
Set-Location $TREE

function Hash-File($p) { (Get-FileHash -Algorithm SHA1 -Path $p).Hash.ToLower().Replace('','') }
function Git-Hash($p) { (& git -C $REPO hash-object $p).Trim() }
function Run-334($label) {
  $log = Join-Path $REPO ".scratch\wisp\probes\339\r1\logs\$label.log"
  go test ./internal/risk/ -run $RUN -count=1 -timeout 30m -v *> $log
  $lines = Get-Content $log -Encoding utf8
  "[$label] RC=$LASTEXITCODE RUN=$(($lines|?{$_ -match '^=== RUN'}).Count) PASS=$(($lines|?{$_ -match '^--- PASS'}).Count) FAIL=$(($lines|?{$_ -match '^--- FAIL'}).Count) SKIP=$(($lines|?{$_ -match '^--- SKIP'}).Count)"
  $lines | ?{$_ -match '^--- FAIL'}
}

"=== PRISTINE check (git hash-object of exported blacklist.go) ==="
$pre = Git-Hash $BL
"pristine_now=$pre expected=$PRISTINE match=$(if($pre -eq $PRISTINE){'YES'}else{'NO'})"

'=== 正控：未突变 ⇒ 期望绿 ==='
Run-334 'win-ac2-positive-unmutated'

'=== 去牙一发：aRelates 尾行回落支 return uncertainAnchorMiss -> return false（内容锚 :281，⛔ 行号当判据）==='
$txt = Get-Content $BL -Raw
$anchor = "`treturn uncertainAnchorMiss(cand, anchor)"
$cnt = ([regex]::Matches($txt, [regex]::Escape($anchor))).Count
"anchor_hits=$cnt (must be 1)"
$muted = $txt.Replace($anchor, "`treturn false")
Set-Content -Path $BL -Value $muted -Encoding utf8 -NoNewline
"mutated_hash=$(Git-Hash $BL)"
Run-334 'win-ac2-deteeth-fallback'

'=== 当场还原 blacklist.go（⛔ 信 echo，重新拉物理字节）==='
git -C $REPO show "e6fd0f5051812a23c72db444d12d76b4d27285a7:internal/risk/blacklist.go" | Out-File -FilePath $BL -Encoding utf8 -NoNewline
$post = Git-Hash $BL
"restored_hash=$post expected=$PRISTINE match=$(if($post -eq $PRISTINE){'YES'}else{'NO'})"

'=== 成对正控二跑：还原后必须回到绿（⛔ 恒红）==='
Run-334 'win-ac2-positive-after-restore'
$post2 = Git-Hash $BL
"final_hash=$post2 match=$(if($post2 -eq $PRISTINE){'YES'}else{'NO'})"
