# 304-a1 — 起手锚点（票 304 `AC#0` 只读普查，腿 `304-a1`）

派单件＝`.scratch/wisp/probes/orch/20261011/txt/304a1-dispatch.txt`；票＝`.scratch/wisp/issues/304-external-killed-test-binaries-orphan-mockllm-exe.md`（**已 `sed -n '1,40p'`＋整件读完**，判读⛔ 基于转述）。

## 1. 锚点三把尺（现跑，落笔即读数）

| 尺（逐字） | 派单期望 | 本腿现跑 | 判定 |
|---|---|---|---|
| `git rev-parse --short HEAD` | `d75db9ee` | **`d75db9ee`** | 相符，开工 |
| `git status --porcelain -- internal cmd tools \| wc -l` | 0 行（派单时刻） | **0 行** | 相符；Go 编译面仍按"有人占着"处理（`309-r2` 产码笔已落 `d75db9ee`＝`internal/tools/bridge.go`＋`toolrequest_nail_309_test.go`，证据件笔可能仍在交） |
| `git rev-parse --abbrev-ref HEAD` | `dev` | **`dev`** | 相符 |

`date` 现跑＝**2026-10-11 15:23 +0800**（本件全部读数的时刻基准；进程名册那一发在 15:1x，同一小时窗内）。

## 2. 本腿射程与硬边界（派单 §1/§0.3 逐字守）

- **只做 `AC#0` 的表①＋表③（含必答两问）＋残留名册重量。⛔ 做表②**（三形实测＝要真起子进程＋`go test -timeout` 那一形）。
- **⛔ 任何 Go 编译面／运行面命令**：本腿全程没跑过 `go build`／`go vet`／`go test`／任何 `.exe`／`scripts/d22scan.sh`／`scripts/portable-tests.sh`；**`go env`、`go list` 也一枚没用**（⇒ 需要它们的读数一律登记成 `Q-304-a1-*`，见 `90-final.md` 末节，⛔ 冒充实测）。
- 写面＝**只新建** `.scratch/wisp/probes/304/a1/**`（`.md`＋`logs/*.txt`）＋票尾**纯 ＋1 行 −0 行**一枚 Progress log bullet。⛔ 动 `internal/**`／`cmd/**`／`tools/**`／`scripts/**`／`.github/**`／`docs/reports/**`／`docs/evidence/s1/**`／`AGENTS.md`／`docs/PLAN.md`／`docs/specs/**`；⛔ 碰 `tools/d22scan/allowlist.txt` 一字节；⛔ 在 `probes/**` 下建 `.go`；⛔ 删任何件；⛔ 杀任何进程。

## 3. 件名册

| 件 | 内容 |
|---|---|
| `10-residual.md` | (A) 当前残留名册重量（四族进程＋归属三件＋`%TEMP%` 目录尺） |
| `20-roster.md` | 表① 全名册（尺名册＋逐枚回收形＋⛔ 看不见的形状具名） |
| `30-costs.md` | 表③ 三形代价表＋必答 (I)(II) |
| `90-final.md` | 结论一页＋欠读数 `Q-304-a1-*` |
| `logs/*.txt` | 每把尺的原始输出（`proc-*`＝PowerShell/CIM，`grep-m*`＝词面尺） |
