# 339-r1 编译两侧尺：GOOS=linux（含 _other 件）与 GOOS=windows（含 _334 件）都 vet 得过。
# 仓外导出树上打，⛔ 共享工作树。
$ErrorActionPreference = 'Continue'
$TREE = 'D:\tmp\wisp\20261018\339r1\buildcheck\tree'
$REPO = 'D:\work\workspace\projects plans\Wisp'
New-Item -ItemType Directory -Force -Path $TREE | Out-Null
Set-Location $REPO
git archive e6fd0f5051812a23c72db444d12d76b4d27285a7 go.mod go.sum internal cmd tools scripts frontend models docs | tar -xf - -C $TREE
# 打上我的三件
Copy-Item "$REPO\internal\risk\anchor_miss_reason_334_test.go" "$TREE\internal\risk\anchor_miss_reason_334_test.go" -Force
Copy-Item "$REPO\internal\risk\anchor_miss_reason_339_other_test.go" "$TREE\internal\risk\anchor_miss_reason_339_other_test.go" -Force
"go=$(Get-ChildItem $TREE -Recurse -Filter *.go).Count"
Set-Location $TREE
$env:GOFLAGS = '-mod=mod'; $env:GOPROXY = 'off'; $env:CGO_ENABLED = '0'
$env:GOTMPDIR = 'D:\tmp\wisp\20261018\339r1\gt'; $env:GOCACHE = 'D:\tmp\wisp\20261018\339r1\gocache'
New-Item -ItemType Directory -Force -Path $env:GOTMPDIR | Out-Null
foreach ($os in @('linux', 'windows')) {
  "--- vet GOOS=$os ---"
  $out = & { $env:GOOS = $os; go vet ./internal/risk/ } 2>&1
  "VET_RC=$LASTEXITCODE"
  $out | Select-Object -First 12
  "--- list GOOS=$os ---"
  $env:GOOS = $os
  $l = & go test -run '^$' -count=1 ./internal/risk/ 2>&1
  "BUILD_RC=$LASTEXITCODE"
  $l | Select-String -Pattern 'FAIL|ok |cannot|undefined' | Select-Object -First 5
}
