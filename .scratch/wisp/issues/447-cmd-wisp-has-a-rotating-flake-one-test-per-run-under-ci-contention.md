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

- [x] **AC#1 量出"谁在红"而不是"红了几次"**：把最近 **20 发** `test-windows` 的 step10 日志全取下来，逐发抽出 `--- FAIL: <测试名>`，产出「测试名 → 红过几次／在哪几发」表。⛔ 只报总数。
- [x] **AC#2 逐枚定性**：每一枚红过的测试，判它是**哪一种**：
  (a) 读快照读早了（同 445 的形状）；
  (b) 真实硬件/设备前提缺失（该**具名弃权并给数字**，⛔ 硬跑）；
  (c) 真产码缺陷（⇒ 立刻停手上报，票升格）；
  (d) 墙钟上界（⛔ 不许拿 `took<5s` 当判据，改用可判事实）。
- [x] **AC#3 修或弃权，逐枚给读数**：(a)/(d) 按 441/445 已确立的形修（等落定＋具名弃权行＋数字＋可 grep token）；(b) 具名弃权；(c) 单独立票。
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
---

## 编排者收件（2026-11-14 19:4x，非实现者）：AC#1–#3 翻；AC#4 记裁定；AC#5 等五发绿；AC#6 腿已给牙

- **尺1**：`5dee200c0` 在。`git show --name-only` = **4 枚 `cmd/wisp/*_test.go`**（`panel_host_windows_budget_winlive_test.go` 新增、`panel_host_windows_test.go`、`resident_voice_interrupt_390r1_windows_test.go`、`resident_voice_wire_r1_windows_test.go`）＋证据＋probes＋本票 ⇒ **产码 0 字节**。**`thresholds.go` 在本 commit 里被改 = 0 枚**（我自己核）。
- **尺2**：`docs/evidence/s1/447-rotating-flake.md` **23445 B**；`probes/447` 里 **20 发 step10 日志 ＋ 3 份红窗原文 ＋ 6 份牙日志**，收集脚本 `collect-ci-fails.ps1` 可逐字重跑。
- **尺3 我自己复跑**（走载具）：`-test.run '^(TestWireR1AnnounceEndToEnd|TestPanelHostRealWindowHopAndLifecycle)$'` ⇒ **RC=0，RUN=2 PASS=2 FAIL=0 SKIP=0**，`ok … cmd/wisp 19.00s`；`PanelHost…` **PASS 1.70s**、`WireR1…` **PASS 17.35s**。
- **AC#1 翻**：20 发里**红过 3 枚**（`390r1…Walks27` ×2、`PanelHostRealWindowHopAndLifecycle` ×1、`WireR1AnnounceEndToEnd` ×1），另 2 发红在 `lint` step7（path-length，票 262）**不在本包**；16 发绿的四数恒 `846/604/0/19`。
- **AC#2 翻**：①`390r1` ×2 = **(a) 读早了**，两发都在 445 修复（`d1dfc6fd7`）**之前**，之后 **11 发零复发**；②`PanelHost` = **(d) 墙钟上界**，红句 `4130.721 ms exceeds D32 panel cold budget 1500 ms`，同码本机 964.3/1098.4/1112.8 ms = **4x**；③`WireR1` = **(a) 读早了**，`finish()` 在 `:1396` 交转写、`:1408` 派 #11、**`:1418` 才 `enterWarm`**，测试在 `:319` 读到合法的 `Awake`（本机修前 `-count=20` ＝ 20 发 4 红，其中 2 发逐字同 CI）。
- **AC#3 翻**：两处读点换 `waitStage447`（谓词＝断言条件、1ms 轮询、5s 上界、超时行带六个数），原断言句子保留；`PanelHost` 的冷预算判决搬进 winlive 档。
- **AC#4 —— 我的裁定（⚠ 标放宽，成立，但要机主知情）**：腿把 **D32 面板冷启动预算的判决**从默认档搬进 `winlive` 档，并在证据 §5 逐条列了「放宽了什么／为什么／代价／撤销口令」。我判**成立**：数目字（1500）**一字节没动**，默认档仍**记录读数**并在 `coldMs<=0` 时红，真机档仍判 ⇒ 这与票 111 wave 10 对 hot 与 P50/P95 的处理同形，是**既有先例的形状**。
  ⚠ **但我必须说清射程变了**：D32 属冻结契约，**CI 默认档从此不再对冷启动 SLO 出判决** —— 这是覆盖率下降，不是文案改动。⇒ 我会把这条**单独告知机主**并附撤销口令（撤销＝把判决搬回默认档，代价＝CI 上该包继续 4x，`PanelHost` 会继续偶发红）。**撤**：「447 冷预算搬回去」。
- **AC#5 未翻**：要**连续 5 发** `test-windows` step10 的 `FAIL` 恒 0（给 5 个 run 号），现在还没跑够。
- **AC#6 未翻**：腿给了两发牙（m1 去掉 `enterWarm` ⇒ 红且点名 `engines_hot=2, warm_turn_end=0`；m2 冷启动真慢 2s ⇒ winlive 档红 `2991.5 ms`；还原 sha256 与仓同值再绿），但 AC#6 要的是「**至少 2 枚**各造一发」，它给的是 2 发 —— 我判**内容满足**，但落点是 `cmd/wisp` 的测试面（回归闸）还是只在 probes，等我下一条审计腿核完再翻（今天已经被这条打过一次）。
- **腿的诚实缺口（照记）**：本机整包四数 `RUN=847 PASS=605 FAIL=2 SKIP=17`，两枚红不是它的用例、且都自述非产码：`TestAC4Focus…`（前台被另一扇真窗口占着）、`TestAC246DevLeg…`「a dev Wisp is already running in this session」。★**第二条是我造成的** —— 我 16:23 起的那个驻留球就在跑，所以本机跑整包时那枚必红。⇒ 记：**本机跑 `cmd/wisp` 整包前，先把我起的球停掉**。

**新读数（2026-11-14 20:5x，编排者自己复跑，来自票 192 接线腿的整包满载跑）**：`Test219AC2SecondRevokeOfTheSameRuleIsRefused` 在**整包满载**时红，红句 `waited on w.reply 46s`（单枚整包 55.78s）；我**低位**单独跑 **`-count=3` ⇒ 3/3 PASS**（2.57／2.12／2.10s，`ok … 7.00s`）⇒ 判 **(d) 满载时序**，不是产码。⇒ 本票 flake 集合从 3 枚增至 **4 枚**（`390r1…Walks27`／`PanelHostRealWindowHopAndLifecycle`／`WireR1AnnounceEndToEnd`／`219AC2SecondRevoke`）。AC#1 那张 20 发表只覆盖 CI；**这一枚是"整包满载"才出现的，CI 上未观测到** —— 口径记明。

---

## 射程更正（2026-11-14 21:1x，编排者现查 CI 日志后自己改的）

**我上面那张 20 发表只看 `test-windows` step10，漏了两件事** ⇒ 本票射程必须扩大到**两个 job**：

**① 含 447 修复的 6 发里，3 发红、3 发绿**（不是我以为的"修完就绿"）：

| run | head | 红的那一枚 | 用时 | job |
|---|---|---|---|---|
| `ff0edb4b0` | `ff0edb4b0` | `TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe` | 15.64s | test-windows step10 |
| `1984618a4` | `1984618a4` | `Test26R3AnnounceClosesTheCaptureGate` | 9.87s | test-windows step10 |
| `35faf8064` | `35faf8064` | `Test285LoopDispatchesOneCorrPerCall` | **0.00s** | **test-core step7（linux）** |

绿的 3 发：`445ea5b98`、`468cd2894`、`6aa244b67`。

**② 落点（我自己 grep 的）**：
- `cmd/wisp/panel_resident_windows_test.go:461`
- `cmd/wisp/resident_voice_selfecho_26r3_windows_test.go:104`
- `internal/agent/ticket285_corr_distinct_rulers_test.go:101`（同包还有 `corr_roster_join_328_test.go` 也提到它）

**③ 由此更正本票的计数**：flake 集合**不是 3 枚、也不是 4 枚，现在是 7 枚**（`390r1…Walks27`／`PanelHostRealWindowHopAndLifecycle`／`WireR1AnnounceEndToEnd`／`219AC2SecondRevoke`（本地满载）／`AC13ColdStart`／`26R3Announce`／`285LoopDispatches`），**跨 `cmd/wisp` 与 `internal/agent` 两个包、`test-windows` 与 `test-core` 两个 job**。

**④ `285Loop` 那枚最可疑**：`0.00s` 就红 ⇒ 不是"读早了"（那类要等），形状像**真缺陷或顺序依赖**。⇒ 派腿时写死：这一枚若判成 (c) 真缺陷，**立刻停手上报**，不许用改测试的方式抹平。
