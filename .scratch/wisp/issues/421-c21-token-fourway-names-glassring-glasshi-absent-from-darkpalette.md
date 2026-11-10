# 421 — `TestC21DesignTokensFourWayAgree` 真红：证据表点名 `Palette.GlassRing/GlassHi`，但 `internal/ball/tokens.go` 的 `DarkPalette()` 无此字面量

- **Status**：ready-for-agent（编排者 2026-11-10 17:2x 立，从 CI `test-windows` 真红反推；415 腿复跑时也撞到，确认非 415 造成）
- 来路：CI run `42643523161` `test-windows` job；以及本机 `go test ./internal/panel -run TestC21DesignTokensFourWayAgree` 现跑复现。

## 现量（⛔ 引用前先重跑）

本机现跑读数：
```
--- FAIL: TestC21DesignTokensFourWayAgree (0.00s)
tokens_fourway_test.go:638: parties: frontend/scripts/gen-tokens.mjs=43 tokens,
  docs/evidence/s1/c21-native-tokens.md=78 colour rows, internal/ball/tokens.go=37/35 palette fields,
  frontend/src/styles/tokens.generated.css=43 light + 43 dark declarations, 11 mapped + 23 waived
tokens_fourway_test.go:726: docs/evidence/s1/c21-native-tokens.md:58: the table names Palette.GlassRing,
  but internal/ball/tokens.go holds no such colour literal in DarkPalette()
tokens_fourway_test.go:726: docs/evidence/s1/c21-native-tokens.md:59: the table names Palette.GlassHi,
  but internal/ball/tokens.go holds no such colour literal in DarkPalette()
```
**归属**：415 只改 `internal/panel/composer_test.go`，未碰 tokens/css/evidence ⇒ **非 415 造成，预存在红**（可核 `git show --stat 551d2a16d`）。

## 这是什么（人话）

四方一致测试要求四份东西互相对得上：①前端 `gen-tokens.mjs`（43 token）②证据表 `c21-native-tokens.md`（78 色行）③Go `internal/ball/tokens.go` 的 palette 字段（37/35）④生成出来的 `tokens.generated.css`（43 light + 43 dark）。现在证据表列了 `Palette.GlassRing`、`Palette.GlassHi` 两个名字，而 Go 的 `DarkPalette()` 里没有对应字面量 ⇒ 四方对不上。

## AC

- [ ] **AC#1 先判后改（⛔ 直接选最小动作）**：判这两枚名字**该删（已废弃）还是该补（漏了）**。判据＝读 `internal/ball/tokens.go` 的 palette 结构、`gen-tokens.mjs`、`tokens.generated.css` 三处看这两个名字在别处是否存在；按「主流 harness／最有价值」定，写明理由。
- [ ] **AC#2 对齐（改哪一份由 AC#1 定）**：
  - 若两名字**已废弃** ⇒ 改 `docs/evidence/s1/c21-native-tokens.md` 删/改那两行（证据表是派生物）；
  - 若两名字**该在** ⇒ 补进 `internal/ball/tokens.go` 的 `DarkPalette()`（并同步 light 侧／css 若需要）；
  - ⚠ ⛔ 只为了让测试变绿去弱化断言——那条四方尺是契约护栏（`23 waived` 已是官方豁免数，⛔ 加大 waived）。
- [ ] **AC#3 正向判据**：`go test ./internal/panel -run 'TestC21DesignTokensFourWayAgree' -count=1` **PASS**，贴读数。
- [ ] **AC#4 反向判据**：在 AC#1 判定的"正确"侧故意造一处漂移 ⇒ 测试必红；还原 ⇒ 绿。两读数都贴。
- [ ] **AC#5 门禁**：gofumpt（绝对路径）空／`go vet ./internal/panel ./internal/ball` rc0／`go build ./internal/panel ./internal/ball` rc0／`go test ./internal/ball ./internal/panel` 绿／`d22scan` clean／`scripts/check-path-length-budget.sh` GREEN。
- [ ] **AC#6 禁碰**：⛔ `frontend/**`（若判定要改 css/gen-tokens.mjs，那是机主外部 agent 的 —— **停手报回来**，本票只做 Go 侧或证据表侧）；⛔ thresholds/golden。

## 撤销口令
- 「**421 撤回**」⇒ 本票作废。