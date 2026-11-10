# 421 — `TestC21DesignTokensFourWayAgree` 真红：证据表点名 `Palette.GlassRing/GlassHi`，但 `internal/ball/tokens.go` 的 `DarkPalette()` 无此字面量 (DONE ✅)

- **Status**：done（2026-11-10 17:5x 编排者裁决关闭；两腿：`a89101612` 落地 panel 四方尺翻转；裁决者≠实现者＝编排者）
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

- [x] **AC#1 先判后改（⛔ 直接选最小动作）**：判这两枚名字**该删（已废弃）还是该补（漏了）**。判据＝读 `internal/ball/tokens.go` 的 palette 结构、`gen-tokens.mjs`、`tokens.generated.css` 三处看这两个名字在别处是否存在；按「主流 harness／最有价值」定，写明理由。
- [x] **AC#2 对齐（改哪一份由 AC#1 定）**：
  - 若两名字**已废弃** ⇒ 改 `docs/evidence/s1/c21-native-tokens.md` 删/改那两行（证据表是派生物）；
  - 若两名字**该在** ⇒ 补进 `internal/ball/tokens.go` 的 `DarkPalette()`（并同步 light 侧／css 若需要）；
  - ⚠ ⛔ 只为了让测试变绿去弱化断言——那条四方尺是契约护栏（`23 waived` 已是官方豁免数，⛔ 加大 waived）。
- [x] **AC#3 正向判据**：`go test ./internal/panel -run 'TestC21DesignTokensFourWayAgree' -count=1` **PASS**，贴读数。
- [x] **AC#4 反向判据**：在 AC#1 判定的"正确"侧故意造一处漂移 ⇒ 测试必红；还原 ⇒ 绿。两读数都贴。
- [x] **AC#5 门禁**：gofumpt（绝对路径）空／`go vet ./internal/panel ./internal/ball` rc0／`go build ./internal/panel ./internal/ball` rc0／`go test ./internal/ball ./internal/panel` 绿／`d22scan` clean／`scripts/check-path-length-budget.sh` GREEN。
- [x] **AC#6 禁碰**：⛔ `frontend/**`（若判定要改 css/gen-tokens.mjs，那是机主外部 agent 的 —— **停手报回来**，本票只做 Go 侧或证据表侧）；⛔ thresholds/golden。

## 撤销口令
- 「**421 撤回**」⇒ 本票作废。

## Progress log（append-only）

- [2026-11-10 17:2x] 编排者立票（从 CI `42643523161` 真红反推；6 AC）。
- [2026-11-10 17:3x] **腿一（判据腿）零写入停手上报**：AC#1 判定＝**废弃**——`GlassRing/GlassHi/GlassHiSoft` 是机主 344 票（A1395/`ebf33fd77`）已裁死色（`gen-tokens.mjs`/`tokens.generated.css`/`theme.css` 均 `grep -ci glass`=0；球侧 `c21V2RetiredGlass` 已反向钉死字段必须不再持有）。两条 AC#2 路（删表行／补字段）均被机主批过的机器钉死（票 380「78 行一字节不删」＋球侧测试）。**真根因＝344 漏改 `internal/panel/tokens_fourway_test.go` 的 P2→P3 遍历**（344 只碰 5 文件，internal/panel 零命中）⇒ 红 8 天。诚实修法在 421 写面外 ⇒ 停手，请扩写面。
- [2026-11-10 17:4x] 编排者**扩写面到 `internal/panel/tokens_fourway_test.go`**，派 421-panel。
- [2026-11-10 17:45 +08] **腿二（421-panel）交付 `a89101612`**：只在 `tokens_fourway_test.go`（+30/−2）新增包级 `c21V2RetiredGlass` map（照抄球侧 `tokens_table_test.go:678-682`）＋ P2→P3 循环镜像翻转（命中 retired 表仍持有即红，然后 continue）。**加牙不是减牙**：命题从"必须有字面量"翻成"必须没有"，78 行证据表一字节未动。判据：修后 `TestC21DesignTokensFourWayAgree` PASS（`72 agreed＋6 枚翻转行＝78 全有牙`）＋ `TestC21TableColourRowsMatchTokensCSS` PASS；反向（往 tokens.go 临时塞回 GlassRing）⇒ 四方红＋球侧红，还原绿（`git status -- tokens.go` 空）；门禁 gofumpt 空/vet rc0/build rc0/两包全量 ok/d22scan clean/path-length GREEN；只 commit 一文件未 push。

## 编排者裁决关闭（2026-11-10 17:5x，判决者≠实现者）

### 1:1 裁决表（AC ↔ 判决 ↔ 证据）

| AC | 判决 | 证据（编排者现跑/引用） |
|---|---|---|
| **AC#1** 先判后改 | **PASS** | 腿一证据链：三处 `grep -ci glass`=0＋球侧 `c21V2RetiredGlass` 反向钉死＋台账 A1395（机主 344 甲批准删死玻璃色）⇒ **废弃非漏补**，判定正确 |
| **AC#2** 对齐 | **PASS** | 选形状＝**panel 四方尺镜像退役翻转**（非删表行/非补字段，两条均被机主钉子堵死）；78 行证据表**一字节未动**，waived 仍 23 ⛔ 未加大；命题翻成"必须没有"＝加牙 |
| **AC#3** 正向 | **PASS** | **编排者现跑** `TestC21DesignTokensFourWayAgree` **ok 0.064s**；全 `./internal/panel` **ok 17.134s**；`./internal/ball` **ok 1.058s** |
| **AC#4** 反向 | **PASS** | 腿二：往 `tokens.go` 临时塞回 `GlassRing` ⇒ 四方 `tokens_fourway_test.go:748` 红 ＋ 球侧 `tokens_table_test.go:1619/1654` 红；还原 ⇒ 两包绿（`git status -- tokens.go` 空） |
| **AC#5** 门禁 | **PASS** | **编排者现跑**：gofumpt 空／`go vet ./internal/panel ./internal/ball` rc0／`go build` 同两包 rc0／两包全量 ok／`d22scan` **clean（exit 0）**／path-length **over-budget=62 covered=62 not in roster=0**（VERDICT GREEN） |
| **AC#6** 禁碰 | **PASS** | 腿二只 commit `internal/panel/tokens_fourway_test.go` 一枚；未碰 `frontend/**`；未动 tokens.go（只临时变异后还原）；未碰 thresholds/golden |

**结论：6/6 AC PASS ⇒ 置 done、改名 `-done`。** 撤销口令：「**421 撤回**」。
