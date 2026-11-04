# 404 — `screenCaptureSchema` 是坏 JSON：任何含截屏工具的 tool-list 序列化即炸（封死所有 run 级测试）

- Status: **done**
- 立案：2026-11-04 18:5x +08 编排者立（源自票 **378** 的 R4 退回发现，见 `docs/evidence/s1/378-ac-verdict.md`）
- 类型：**前置真 bug**（非 378 造）

## 这是什么（人话）

`internal/tools/system.go` 里 `screenCaptureSchema` 把 `"target":{...}` 之后**直接拼** `windowTargetProperties`（后者以 `"hwnd":{...}` 起），**漏了一个逗号**：

```
`"target":{...region 一块区域"}` +     ← 这里 `}` 后应有一个 `,`
windowTargetProperties + `,` +          ← 以 `"hwnd":{...}` 起
```
⇒ 组出来的 schema 是 `..."region"}"hwnd":{...}` —— **非法 JSON**。

## 后果

- 任何把截屏工具纳入 tool-list 的请求，序列化即报：
  `json: error calling MarshalJSON for type *jsontext.Value: invalid character '"' after object key:value pair`。
- ⇒ **今天任何「真 wisp run 跑一发」的测试都不可能绿**（如 `Test233HotTierAgentReachesTheNextTask` 现红即此因）。这正是票 378 的 AC#2 想要 run 级验收却做不出的 R4 前置红。

## 结案判据

- [x] AC#1 修：`system.go:470` `"target":{...}` 后补 `,`（`...一块区域"},`）。
- [x] AC#2 牙：新增 `internal/tools/screencapture_schema_404_test.go::Test404ScreenCaptureSchemaIsValidJSON`（`json.Valid(screenCaptureSchema)`）；修前红、修后绿。
- [x] AC#3 门禁：`go build ./internal/tools`、`go test ./internal/tools -run Test404`、`go test ./cmd/wisp -run Test233…`（转绿＝R4 解除）、`gofumpt`、`d22scan`。
- [x] AC#4 纪律：仅动 `system.go` 一行＋一枚测试；⛔ 碰 thresholds/golden/frontend/design；⛔ 改任何 C/D 契约文字（tool 参数 schema 非契约）。

## Progress log

- [2026-11-04 18:5x +08] agent=`编排者` did=修 `system.go:470` 补逗号＋加 `Test404ScreenCaptureSchemaIsValidJSON` 回归；门禁跑尺；源起票378-R4。
