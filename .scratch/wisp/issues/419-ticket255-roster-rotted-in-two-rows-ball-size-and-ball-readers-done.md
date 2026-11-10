# 419 — 票 255 的名册**两行腐烂** ⇒ 它的守护测试在 CI 上红（其中一行是**我自己**改出来的） (DONE ✅)

**Status**: done（2026-11-10 19:5x 编排者裁决关闭（2026-11-10 09:2x 编排者建，从 CI 真红反推；票 255 本体已 `-done`）

## 现量（⛔ 引用前先重跑）

**CI 真红**（`test-windows` job `129250564875`，run `42643523161`，headSha `d797a855`）：
```
2026-11-09T15:28:30.4293776Z config_receipt_255_test.go:240: row "ball" cites internal/ball/ball_windows.go:64 for
  "SizePx   int // configured orb size 44..72" and the token is now nowhere in that file
2026-11-09T15:28:30.4309177Z --- FAIL: TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim (0.00s)

2026-11-09T15:28:30.4314596Z config_receipt_255_test.go:303: SECTION 4/14 section="ball" field=Ball
  roster=[] scan=[cmd/wisp/ball_style_65_windows.go cmd/wisp/config_cmd_65.go] match=false
2026-11-09T15:28:30.4316213Z config_receipt_255_test.go:305: section "ball" (field Ball) is read by
  [cmd/wisp/ball_style_65_windows.go cmd/wisp/config_cmd_65.go] in production, the roster names []: something...
--- FAIL: TestTicket255RosterStillMatchesTheActualReadSites
```

## 两处根因（逐条现跑核过）

### ① 尺寸注释 `44..72` → 盘上已是 `44..160`（**我造成的**）
| 项 | 读数 |
|---|---|
| 盘上 `internal/ball/ball_windows.go:64` | `SizePx   int // configured orb size 44..160 (0 = default 56)` |
| 名册 `cmd/wisp/config_readers_255.go:167` | `... ball_windows.go:64 [SizePx   int // configured orb size 44..72] ...` |
| 谁改的 | **`931077cfa`**「ball: internal/ball 尺寸上限跟随 72 -> 160（A1636 批 B）」 |
| 那次改了哪些文件 | **只有 `internal/ball/ball_windows.go`（1 文件 1 行）** ⇒ **名册没跟** |

⇒ **`931077cfa` 是我（编排者）在 `A1636` 批准上限 72→160 后落的批 B。** 我改了产码注释，**没同步票 255 名册里逐字引用那一行的 token** ⇒ 名册按 token 查不到 ⇒ 红。

### ② 名册说球"零读者"，实际有两处真读者
| 项 | 读数 |
|---|---|
| 名册 `:167` 逐字 | `hotClaimNoReader + "... the ball's own option carries the size, and **nothing outside internal/config reads cfg.Ball**"` |
| 真读者 ① | `cmd/wisp/ball_style_65_windows.go:50` → `return c.Ball.Style, true` |
| 真读者 ② | `cmd/wisp/config_cmd_65.go:235` → `return c.Ball.Style, true, nil` |
| 读者 ③ | `cmd/wisp/config_cmd_65.go:150` → `config.NewDefaults().Ball.Style` |

⇒ **名册断言"零读者"是 2026-10-04 立票时刻的事实，已腐烂**（与票 180 `A1667` **同一味**）。

## ★为什么这两行都该"改名册"而不是"改测试"

- 票 255 的守护测试**问的是对的**：名册里逐字引用的 token 必须还在那个文件里（防"名册说的是老事实"）。
- **腐烂的是名册那两行**，不是测试。⇒ 修法＝**把名册那两行改成今天的事实**。
- ⛔ **不许**把测试改成"忽略 token 找不到"——那是**把闸拆了**（票 255 的全部价值就在这道闸上）。

## AC

- [x] **AC#1 修 ①**：`cmd/wisp/config_readers_255.go:167` 里的 `44..72` 改成 **`44..160`**（**逐字对齐盘上 `ball_windows.go:64` 的当前内容**）。改完现跑 `grep -n "44\.\.160" cmd/wisp/config_readers_255.go` 必须有 1 命中。
- [x] **AC#2 修 ②**：`:167` 那行的 `hotClaimNoReader` 判定与「nothing outside internal/config reads cfg.Ball」那句**必须重判**：
  - 现量有 **3 处** `c.Ball.Style` 读者（上表）⇒ 若它们算"生产读者"，则该行**不该是 `hotClaimNoReader`**，须改成具名读者形（照同文件其它有读者行的形状）；
  - ⚠ **若判"`.Style` 算读者但 `.Size` 不算"**（票 413 正是要喂 `.Size`），则须在行内**分开写清两个字段各自的读者状态**，⛔ 用一句话盖住两个字段。
  - **腿先判后改，写明理由。**
- [x] **AC#3 正向判据**：`go test ./cmd/wisp -run 'TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim|TestTicket255RosterStillMatchesTheActualReadSites' -count=1 -timeout 30m` **两枚都绿**，贴输出。
- [x] **AC#4 反向判据**：把 `:167` 的 `44..160` 故意改回 `44..72` ⇒ 第一枚测试**必须变红**；还原 ⇒ 变绿。两个读数都贴。
- [x] **AC#5 别只修这两行**：那 15 条 `drifted, the token` 是 **`t.Logf`（设计允许，不是红）**，⛔ 当缺陷修；但**顺手**用它们刷新名册行号（票 255 AC#6 要求"下次重裁可刷新"）。
- [x] **AC#6 门禁**：`gofumpt`（绝对路径）空／`go vet ./cmd/wisp` rc0／`go build ./cmd/wisp` rc0／`d22scan` clean／path-length GREEN。
- [x] **AC#7 CI 侧验证**：推一发后 `gh api repos/CarlosShao/wisp/actions/jobs/<id>/logs` 拉 `test-windows`，**这两枚测试必须 PASS**，贴那几行。

## 禁区
- ⛔ 放宽或删除那两枚守护测试（**那是把闸拆了**）
- ⛔ 动 `internal/ball/ball_windows.go` 的 `44..160`（**那是机主批过的数字**，`A1636`）
- ⛔ 动 thresholds / golden / allowlist.txt
- ⛔ 碰 `frontend/**`、`design/**`

## 撤销口令
- 「**419 撤回**」⇒ 本票作废。
## 编排者裁决关闭（2026-11-10 19:5x，判决者≠实现者）

### 1:1 裁决表（AC ↔ 判决 ↔ 证据）

| AC | 判决 | 证据（编排者现跑/引用） |
|---|---|---|
| **AC#1** 44..72→44..160 | **PASS** | **编排者现跑** `grep -c "44\.\.160" cmd/wisp/config_readers_255.go` = **1** |
| **AC#2** 重判 no-reader | **PASS** | 名册 ball 行现分两段具名读者：`.Size` 半＝`resident_ball_windows.go [SizePx: hs.sizePx]`（413 喂线），`.Style` 半＝`ball_style_65_windows.go:50`（名册 181–189 行注释）——⛔ 未用一句盖两字段 |
| **AC#3** 两测绿 | **PASS** | **编排者现跑** `go test ./cmd/wisp -run TestTicket255Roster` **ok 1.179s**；CI 同（见 AC#7） |
| **AC#4** 反向 | **PASS** | 腿：`44..160`→`44..72` ⇒ 第一枚红；还原 ⇒ 绿（两读数） |
| **AC#5** 刷新行号 | **PASS** | 名册行注释按 413 落地后读点刷新（181–189 引 `SizePx: hs.sizePx`）；15 条 `drifted` t.Logf 未当缺陷修（正确） |
| **AC#6** 门禁 | **PASS** | 腿：gofumpt 空／`go vet ./cmd/wisp` rc0／build rc0／d22scan clean／path-length GREEN |
| **AC#7** CI 侧验证 | **PASS** | **实 CI 日志** run `42804517323` / test-windows job `129790147091`：`TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim` **PASS**、`TestTicket255RosterStillMatchesTheActualReadSites` **PASS 0.68s**、`TestTicket255RosterReJudgedAfterTickets44And61` **PASS** |

**结论：7/7 AC PASS ⇒ 置 done、改名 `-done`。** 撤销口令：「**419 撤回**」。