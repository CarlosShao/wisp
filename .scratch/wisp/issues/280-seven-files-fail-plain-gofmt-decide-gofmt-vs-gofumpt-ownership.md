# 票 280 — 7 枚文件不过裸 `gofmt`：逐枚判归属（机器 `gofmt` vs CI 那步 `gofumpt`）

**立票**：2026-10-08 19:3x 编排者（来路＝台账 `A736` 收工门禁快照：`gofmt -l cmd internal` ＝ 7 命中）
**性质**：读数为先；**不预设要改**（改它们＝"顺手改"，且可能与 CI 那步不是同一把尺）。

## 现量（引用前重跑；⚠ 数字是快照）

- 7 枚：`cmd/wisp/models.go`／`panel_inbound_guards_35r3_test.go`／`panel_transport_35r2_test.go`（10-08 今天那枚）／`internal/agent/approval/pending_read.go`／`internal/agent/tools.go`／`internal/risk/provenance.go`／`internal/tools/bridge.go`。
- 来历（`A736` 逐枚 `git log -1`）：6 枚既有（10-03～10-07），1 枚今天已提交、树净。
- ⚠ CI 的 lint 步量的是 **gofumpt**，与裸 `gofmt` **不是同一把尺**（`A721` 一族）。

## 要建什么

- [ ] **AC#1 逐枚判**：7 枚各给"差异属于 `gofmt` 风格还是 `gofumpt` 加严"（本机若有 gofumpt 二进制就跑，没有就具名写"判不动＋缺什么"）。
- [ ] **AC#2 与 CI 的关系**：逐枚答"它在 CI 那步会不会被点红"，⛔ 不许用"应该会"填空。
- [ ] **AC#3 处置建议只写形状**：改／不改／留给下一枚真正动该文件的票；⛔ 本票不落地任何格式改动。

## 禁区

⛔ 不许 `-w`；⛔ 不许批量重排；⛔ 三枚冻结件一字不动；⛔ 零 push。
