# 303-r1 起手锚 (anchor)

现量时刻锚定：HEAD 与台面读数同一次命令取得。

## HEAD
`git log -1 --format='%h %ad' --date=iso-strict`
4e00357f 2026-10-10T17:34:47+08:00

## 台面干净度
`git status --porcelain -- cmd internal scripts .github docs | wc -l`
0

## 进程计数 (tasklist)
- wisp.exe: 0
- balldebug.exe: 0
- msedgewebview2.exe: 24

## 我的尺清单 (本次将用)
- 定向两发受害用例: `go test ./cmd/wisp/ -run 'TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe|TestAC14AwaitedBindingReplyReachesThePage|TestAC14GoSideEvalPushReachesThePage' -v`
  - 前置: 母仓跑 ./cmd/wisp/ 必须铺 sherpa (shell 形 PATH)
    `export PATH="/d/work/workspace/projects plans/Wisp/third_party/sherpa-onnx:/d/work/workspace/projects plans/Wisp/build:$PATH"`
- 整包两发取交集: `go test ./cmd/wisp/ -v` (改前一发 / 改后一发, 逐字落文件)
- 已知负载敏感红 (不算新增红, 也必报): `TestResolvePerCallBudget`
- 已知两形 (预算红 / 哨兵 `-1.000` 红): `TestPanelHostRealWindowHopAndLifecycle`
- 门禁四把: `go vet ./cmd/wisp/` / `go vet -tags winlive ./cmd/wisp/ ./internal/ball/` / `sh scripts/d22scan.sh` / `gofmt -l cmd/wisp` (HEAD blob 名册并排)
- 归因承重复跑: `git show fb2fb802 -- cmd/wisp/panel_host_windows.go` (那 57 行 = 第一手现场)

## 写面范围
- cmd/wisp/** (受约束产码/夹具)
- .scratch/wisp/probes/303/r1/** (证据件, 每件自落一行 rc=N, 非 .out, 非 0 字节)
