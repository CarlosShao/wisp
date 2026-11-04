# 404 裁决 — `screenCaptureSchema` 坏 JSON 修复（机械可验，源起票378-R4）

## AC 判语

- **AC#1 修**：`internal/tools/system.go:470` 在 `"target":{...一块区域"}` 之后补 `,`（`...一块区域"},`），使 `"target"` 对象与 `windowTargetProperties`（以 `"hwnd":{...}` 起）之间合法分隔。**已改**。
- **AC#2 牙（修前红/修后绿，机械可验）**：
  - 修前现跑：`go test ./cmd/wisp -run Test233HotTierAgentReachesTheNextTask` → **FAIL**，报 `json: error calling MarshalJSON ... invalid character '"' after object key:value pair`（截屏工具入 tool-list 序列化即炸）。
  - 修后现跑：同命令 → **ok 2.375s**；新增 `internal/tools/screencapture_schema_404_test.go::Test404ScreenCaptureSchemaIsValidJSON`（`json.Valid(screenCaptureSchema)`）→ **ok 0.045s**。
  - ⇒ 修前红、修后绿＋专项回归，牙成立。
- **AC#3 门禁**：`go build ./internal/tools` rc=0 · `go test ./internal/tools -run Test404` ok · `go test ./cmd/wisp -run Test233…` ok · `gofumpt -l` 空 · `d22scan` clean。
- **AC#4 纪律**：仅 `system.go` 一行＋一枚测试；⛔ 碰 thresholds/golden/frontend/design；⛔ 改 C/D 契约文字（工具参数 schema 非契约）；实现者＝编排者、本判语建立在机械 before/after 读数上（非主观裁量）。

## 意义

- 这是**票 378 R4 的前置真 bug**：它使「任何真 wisp run 跑一发」的测试今天不可能绿。
- 修复后，**票 378 的 AC#2 真 spawn+fkc 验收测试变得可行**（run 级不再被此封死）。

## 双名核
- 关票后 `git ls-tree -r HEAD --name-only | grep -c 'issues/404-'` 须 = 1（见关闭提交）。
