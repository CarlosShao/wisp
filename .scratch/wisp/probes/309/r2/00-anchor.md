# 309-r2 · 00 锚件（写腿，起手声明写面）

腿＝`309-r2`（写腿，非裁决者）。派单＝票 309 末节「编排者复跑＋裁定（2026-10-11 14:1x +08）」§4 那五件。
起手 HEAD＝`202ed4a5`（与派单给的锚点一致，尺＝`git rev-parse --short HEAD`，时刻 2026-10-11 14:21 +08）。

## 写面（三件套，逐字）
1. `internal/tools/**` —— 只两枚测试件与本票那枚钉（`toolrequest_nail_309_test.go`）＋`bridge.go` 的**注释行**。
2. `.scratch/wisp/probes/309/r2/**` —— 证据件（`.md`）＋`logs/`（原始读数、Rig 脚本 `.txt`）。
3. 本票末尾**一条** append-only bullet（1 加／0 删）。

## ⛔ 零（本腿自认的硬边界）
- ⛔ 翻任何框、⛔ 改任何 `Status:` 行、⛔ 声明任何格"成立"（判语归 `309-v2`）。
- ⛔ 动 `bridge.go` 的布尔式 `if b.injected || b.prov == nil || taskID == "" {`、⛔ 删那一肢、
  ⛔ 动 `internal/risk/**` 产品码、⛔ 改 `lit236*` 两枚常数、⛔ 给 `agent.ToolRequest` 加字段。
  ⇒ 行为必须逐字节不变，交付里带尺（`10-behbyte.md`）。
- ⛔ 为变绿放宽断言、⛔ 造 `--- SKIP`、⛔ 把空名册当"扫过了"（分母照打印）。
- ⛔ 在字符串字面量里写 `⛔ ≈ ✓ ≤`（d22scan ban #8 射程，注释豁免、字符串不豁免）。
- ⛔ push／`git add -A`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`／在仓内建 worktree。
- ⛔ 在 `probes/**` 里新建 `.go`（证据里的程序写 `.go.txt` 或嵌进 `.md`）；证据件⛔ 叫 `.out`。
- ⛔ 为 `b.prov == nil` 那一肢新造夹具（编排者 §4 第 4 件裁＝只补注释一行）。
- ⛔ 为 `bridge.go:935` 那句"其余闸门"造第二枚仪器（裁＝改射程句）。

## 台面备忘（⛔ 我造的，落档免得下一位误读）
- 工作树里 `go.mod` 处于别腿在飞态（`git diff -- go.mod`＝把 `github.com/jchv/go-webview2` 从 indirect 移到 direct），
  `internal/tools`／`internal/risk`／`cmd` 三处起手 `git status --porcelain`＝**0 行**（尺见 `logs/start-state.log`）。
  ⇒ 整包名册的**改前/改后两发都走仓外导出树**（`git archive <sha>`），同一台面可比；工作树那发只作旁证并具名它带着别腿的 go.mod 在飞态。
- 起手 `sh scripts/check-path-length-budget.sh --with-self-test`＝**rc=0**（VERDICT GREEN，读数 `logs/pathlen-start.log`）。
