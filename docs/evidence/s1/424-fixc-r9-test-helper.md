# 票 424 — 第三批：`r9_mode_card_test.go` 跨平台 test 泄漏修复（编排者直做·零逻辑改动）

- 时刻：2026-11-12 11:1x +08
- 锚：`3f15ede4d` 之后（424-fix-b 已提交）
- 病灶（overlay 现测）：`cmd/wisp/r9_mode_card_test.go:146/162: undefined: modeRequest114r2`
  - `modeRequest114r2` 定义在 `panel_mode_confirm_114_windows_test.go`（`//go:build windows` test），却被**无 build tag 的** `r9_mode_card_test.go` 引用 ⇒ `GOOS=linux` 下 test 变体编不过 ⇒ 会红 `go vet ./...`（`go vet` 含 test 文件）。

## 改了什么（零逻辑改动）

1. **新建 `cmd/wisp/mode_request_114r2_helper_test.go`**（无 build tag，`package main`）——把 `modeRequest114r2` 原函数体**逐字**搬入；它只依赖跨平台 test 常量 `inboundPanelSource`（`panel_inbound_33_test.go:39`，无 tag），故跨平台安全。
2. **`cmd/wisp/panel_mode_confirm_114_windows_test.go`** —— 删掉被搬走的 `modeRequest114r2`（7 行：3 注释 + 函数签名 + 3 体行）。
- 函数体一字未改；无断言/阈值/契约改动。

## 判据（我自己复跑）

| 项 | 读数 |
|---|---|
| **Linux `go vet -overlay(go.mod→sherpa桩) ./cmd/wisp`（含 test）** | **rc=0，undefined=0** |
| **Linux `go test -overlay -gcflags=all=-e -c ./cmd/wisp`（编译 test 变体，全量）** | **rc=0，undefined=0** |
| Windows `go build ./cmd/wisp` | **rc=0** |
| Windows `go vet ./cmd/wisp` | **rc=0** |
| Windows 定向 `go test -run 'TestTicket114r2|TestR9|TestTicket255|Test324|Test180K|Test423|Test65'` | **RUN=50 PASS=49 FAIL=0 SKIP=1，rc=0** |

## 越界自证

`git status --porcelain -- cmd/wisp/`：本件只 2 个落点——` M panel_mode_confirm_114_windows_test.go`、`?? mode_request_114r2_helper_test.go`＋本证据件。其余 `??`（`panel_resident_323_windows_test.go` 等）是**别的腿既有未跟踪件，未动未提交**。

## 一句话

把纯字符串构造 `modeRequest114r2` 从 windows-only test 文件搬进无 tag 的 test helper，`GOOS=linux` 下 cmd/wisp 的 build 与 test 变体都 type-check 通过（undefined=0），Windows build/vet/定向测试不回归；这是 cmd/wisp 在 Linux 上编过的最后一块。
