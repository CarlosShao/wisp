# 票 447 — `cmd/wisp` 在 CI 争用下**每发红一枚不同的测试**：四数永远一样，红的那枚在换

**状态**：待派（**系统性**，不是单枚票）
**为什么单独立票**：445 修好了 `Test390r1BallClick...`（它现在本地 17/17 绿），但**红的换人了**。单枚修等于追着影子打。

## 现量（三次 CI 失败，逐条来自 `--log-failed`，不是转述）

| run | head | 红的那一枚 | 四数 |
|---|---|---|---|
| `43575891409` | `d2cfa96e8` | `Test390r1BallClickPlayingInterruptAbortsAndWalks27`（2.47s） | `RUN=846 PASS=603 FAIL=1 SKIP=19` |
| `43583105791` | `b7e85ab75` | `TestPanelHostRealWindowHopAndLifecycle`（4.30s） | `RUN=846 PASS=603 FAIL=1 SKIP=19` |
| 期间另有 5 发 success | — | 无 | 同样 `RUN=846` |

⇒ **同一套 846 枚用例、同一套 SKIP 19 枚，每发恰好红 1 枚，但红的是不同的一枚**。这不是某枚测试坏了，是**这一包在 CI 的争用/时序下有多个敏感点**。

## 完成判据（AC）

- [ ] **AC#1 量出"谁在红"而不是"红了几次"**：把最近 **20 发** `test-windows` 的 step10 日志全取下来，逐发抽出 `--- FAIL: <测试名>`，产出「测试名 → 红过几次／在哪几发」表。⛔ 只报总数。
- [ ] **AC#2 逐枚定性**：每一枚红过的测试，判它是**哪一种**：
  (a) 读快照读早了（同 445 的形状）；
  (b) 真实硬件/设备前提缺失（该**具名弃权并给数字**，⛔ 硬跑）；
  (c) 真产码缺陷（⇒ 立刻停手上报，票升格）；
  (d) 墙钟上界（⛔ 不许拿 `took<5s` 当判据，改用可判事实）。
- [ ] **AC#3 修或弃权，逐枚给读数**：(a)/(d) 按 441/445 已确立的形修（等落定＋具名弃权行＋数字＋可 grep token）；(b) 具名弃权；(c) 单独立票。
- [ ] **AC#4 不许钝化**：⛔ 不许删/放宽任何既有断言来让它绿；若要放宽，必须在证据里逐条列出放宽了什么、为什么，并标 `⚠ 放宽`。
- [ ] **AC#5 收尾判据**：修后**连续 5 发** `test-windows` 的 step10 里 `FAIL` 那一位恒为 0（给 5 个 run 号）；且本地 `./cmd/wisp` 走载具跑一次给四数。
- [ ] **AC#6 牙**：至少对修过的其中 **2 枚**各造一发「把护栏拆掉」的突变（仓外导出树）⇒ 必须真红并点名；还原 SHA256 同串、再绿。

## 边界
写面：`cmd/wisp/**` 的 `*_test.go`（**本波只有你一枚写腿碰 `cmd/wisp`**）、`docs/evidence/s1/447-*`、`.scratch/wisp/probes/447/**`、本票 Progress log（只追加）。
⛔ 产码（除非 AC#2 判出 (c) 真缺陷 —— 那种情况**停手上报**）；⛔ 阈值/golden/契约/规格；⛔ 阈值与断言放宽；⛔ `internal/panel/**`、`internal/tools/**`、`tools/d22scan/**`、`internal/llm/**`；⛔ add -A/--amend/reset/rebase/stash/checkout .；⛔push。
⚠ `go test ./cmd/wisp` **必须走载具** `scripts/wisp-cli-staged-tests.sh`（裸跑因 DLL 出假红 `0xc0000135`），⛔ 同时跑两发；整包一次要几分钟，⛔ 拿定向子集冒充整包四数。

## 撤销口令
「票 447 撤」＝删本票＋回退其写腿 commit（新 commit 反改，不 reset/amend）。

## Progress log

- **2026-11-14 18:38–19:5x +08，写腿（委派，起手 HEAD `c605f3ecb`）** — AC#1 量完、AC#2 三枚全部定性、AC#3 修两处（全是 `*_test.go`，产码 0 字节）、AC#6 两枚牙在仓外导出树里真红并点名（`verdict=ALL-TEETH-BIT`）。证据全文 `docs/evidence/s1/447-rotating-flake.md`，台件 `.scratch/wisp/probes/447/**`。
  - **谁在红（20 发，结论非空的最近 20 发）**：只有 **3 枚**——`Test390r1BallClickPlayingInterruptAbortsAndWalks27`（2 发：`43575891409`/`43577965575`，**都早于 445 的 `d1dcf6fd7`（09:23Z）**，445 之后 11 发零复发）、`TestPanelHostRealWindowHopAndLifecycle`（1 发：`43583105791`）、`TestWireR1AnnounceEndToEnd`（1 发：`43575529807`）。另 2 发（`43575854808`/`43575830504`）红在 `lint` 的路径长度预算，不在本包。
  - **逐枚定性**：390r1 = (a) 已由 445 修；PanelHost = (d) hosted runner 的冷启动 4130.7ms 越 D32 1500ms（同码本机 964–1113ms），判决该由 winlive 档下；WireR1 = (a) 在「转写到文本路」这个标记上读 `stage()`，而 `enterWarm` 在那之后——本机修前 `-count=20` **4 红**，其中 2 发逐字就是 CI 那句。
  - **修**：WireR1 两处读点换 `waitStage447`（谓词＝断言条件，1ms 轮询，5s 上界，超时行带 `engines_hot/warm_turn_end/loads_after_warm/refs/finals/handed`）；D32 冷预算 1500ms 的判决从默认档搬到 winlive 档（**边界一字未改**，默认档保留读数与「没有读数＝红」，并留可 grep token `447-cold-budget-verdict-in-winlive`）。⚠ 一条**放宽**已在证据 §5 逐条列出并给了撤销口令。
  - **牙**：m1 去掉 `enterWarm("turn-end")` ⇒ 红并点名（`engines_hot=2, warm_turn_end=0, stage="Awake"`）；m2 冷启动真慢 2s ⇒ winlive 档红 `2991.5 ms exceeds D32 panel cold budget 1500 ms`；两次还原 sha256 与仓同值、再绿。
  - **AC#5 的一半做不了**：CI 连续 5 发要等编排者推送后用 `collect-ci-fails.ps1 -Limit 12 -Need 5` 复读；本机整包四数见证据 §7。
  - **报而不改**：`-test.count=20`（同进程重复）会多出第三种红（`wire_r1:370` #26 走未落，`state=Sleeping`，2/20；`-count=1` 不见）——像同进程里上一发的收尾把机器按回初始态，属第五个候选，建议单独立票。

