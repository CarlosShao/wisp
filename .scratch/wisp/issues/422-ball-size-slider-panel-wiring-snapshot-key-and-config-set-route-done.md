# 422 — `[ball] size` 滑块的**面板接线**：快照键 `ballSize` 与 `config.set ball_size` 路由未接，页面拿不到读口/写被门拒 (DONE ✅)

- **Status**：done（2026-11-10 19:4x 编排者裁决关闭；落地 commit `587e0e22e`；裁决者≠实现者＝编排者）
- 来路：票 413 落地（`c8dbd2288` 喂线＋`baefa152b` 写接口/往返/契约证据）的**诚实边界**——413 是后端 + 契约证据（AC#6 明说"交外部 agent 接线"），但**面板两条路由**在 `internal/panel/**`（413 写面外）未接。
- 关联：票 413（契约件 `docs/evidence/s1/413-ball-size-slider-contract.md`）· 票 65（ball_style 的 composer/pump/carrier + config.set roster 链，照它复刻）· 票 253（inbound roster）· 票 255（config_readers roster）

## 现量（⛔ 引用前先重跑）

- 413 已建：读口 `BallSizeOptions`/`BallSizeView`、写口 `SetBallSize`（`internal/config/settings_413.go`）。
- **未接**：`grep -rn "ballSize\|ball_size" internal/panel/ cmd/wisp/ | grep -v _test` 应见零（除 413 写面）；composer Snapshot 无 `ballSize` 键；`config.set` 名册（`internal/panel/config_write_65.go:39-46`）只有 `ball_style`，无 `ball_size` case ⇒ 页面写 ball.size 被 `unlisted-field` 门拒。
- 照 65 链：`ball_style` 的 composer/pump/carrier + config.set roster 四件（读 65 那几张票/码定位）。

## AC

- [x] **AC#1 快照读键**：`Snapshot` 增 `ballSize` 键（值＝`BallSizeOptions().Size`），composer/pump/carrier 三处照 `ballStyle` 同形扩；判据＝快照 JSON 含 `ballSize` 且值==config 现值。
- [x] **AC#2 写路由**：`config.set ball_size` 进 `config_write_65.go` 名册 + `panel_config_store.go` ApplySetting 加 `ball.size` case（调 `SetBallSize`，越界拒绝透出 validate 句）；判据＝经 config.set 写 90 读到 90、写 43/161 拒绝贴原文。
- [x] **AC#3 反向判据**：删 ballSize 快照键 ⇒ 读判据红；删 ball_size 写 case ⇒ 写判据红。两读数贴。
- [x] **AC#4 门禁**：gofumpt 绝对空／go vet ./internal/panel ./cmd/wisp rc0／go build 同包 rc0／go test ./internal/panel -count=1 绿 ＋ go test ./cmd/wisp -run '<新用例>' -count=1 -timeout 30m 绿／d22scan clean／path-length GREEN。
- [x] **AC#5 禁碰**：⛔ frontend/**（外部 agent）；⛔ thresholds/golden/allowlist；⛔ 改 413 契约件原文。

## 撤销口令
- 「**422 撤回**」。

## Progress log（append-only）

- [2026-11-10 19:2x] 编排者立票（从 413 诚实边界反推；5 AC）。
- [2026-11-10 19:3x] 落地腿交付 `587e0e22e`（12 件 +483/−7）：读路由＝新 carrier `internal/panel/ball_size_413.go`＋`composer.go`/`pump.go`/`run.go` 三处照 ballStyle 同形扩（reader-gated）；写路由＝`config_write_65.go` 名册 1→2（`ball_size`）＋`panel_config_store.go` ApplySetting case 调 `SetBallSize` 档位静态答 restart；`run.go` 装配根用同一 Manager 的 `BallSizeOptions()` 注入。**两处自判**：①动 `cmd/wisp` 生产码（AC#2 明令要 ApplySetting case，属票面要求非越界）②把 415 缺口台账延伸到顶层对账登记 `PanelSnapshot:{ballSize}`（panel.ts 未声明该键；登记带双向腐烂判定＝UI 腿一旦声明即自红逼删，加牙非静音）。

## 编排者裁决关闭（2026-11-10 19:4x，判决者≠实现者）

### 1:1 裁决表（AC ↔ 判决 ↔ 证据）

| AC | 判决 | 证据（编排者现跑/引用） |
|---|---|---|
| **AC#1** 快照读键 | **PASS** | `internal/panel/composer.go:160 BallSize *BallSizeSection json:"ballSize,omitempty"`＋`ball_size_413.go` carrier＋run.go 注入；反向①证明无 reader 不发键 |
| **AC#2** 写路由 | **PASS** | `config_write_65.go` 名册含 `ball_size`（fields65+takesValue）＋`panel_config_store.go` ApplySetting case 调 `SetBallSize`；**编排者现跑** `go test ./cmd/wisp -run 'Test422\|Test413\|Test65'` **ok 4.397s**（写90→读90、43/161 拒贴原文、拒绝后字节不动） |
| **AC#3** 反向判据 | **PASS** | 腿：①摘 pump 填充块 ⇒ `TestTheBallSizeKeyTravelsOnlyFromItsReader422` **FAIL**（sent no ballSize key）②摘 ApplySetting case ⇒ `Test422ApplySettingBallSizeWritesAndRefuses` **FAIL**；均还原绿 |
| **AC#4** 门禁 | **PASS** | **编排者现跑**：`go test ./internal/panel` **ok 18.978s**；`go vet ./internal/panel ./cmd/wisp` rc0；gofumpt 空；`d22scan` **clean**（腿）；path-length **over-budget=62 covered=62 not in roster=0**（VERDICT GREEN） |
| **AC#5** 禁碰 | **PASS** | 只 commit internal/panel＋cmd/wisp 接线；`git status -- frontend design thresholds golden` 零改动；413 契约件原文一字未动 |

**结论：5/5 AC PASS ⇒ 置 done、改名 `-done`。** ⚠ **诚实边界（不属本票 AC）**：`frontend/src/lib/panel.ts` 未声明 `ballSize`、`config-screen.tsx:105` 范围仍旧 44-72 —— 归外部 UI agent（Q-51），本票台账条目即留给它的接口。撤销口令：「**422 撤回**」。
