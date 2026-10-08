# 02 · 跨发红名册作差（8 发 `ci`，`cmd/wisp` 那一步）

## 取数（`gh` 全部一次通过，零重试）

| run# | run id | headSha | 触发 | conclusion | 本腿取的件 | 字节 |
|---|---|---|---|---|---|---|
| 305 | 37703959747 | `cc31526165` | schedule | failure | `D:/tmp/wisp111c2/run305.log`（复用 `111-c1` 那份，**本腿自己重切的步**） | 2 013 255 |
| 304 | 37545246395 | `cc31526165` | schedule | failure | `run304.log` | 2 014 234 |
| 303 | 37406757402 | `cc31526165` | push | failure | `run303.log` | 2 206 960 |
| 302 | 37406422380 | `b948bcb84e` | push | failure | `run302.log` | 1 916 415 |
| 301 | 37405698188 | `ec84cff16c` | push | failure | `run301.log` | 1 915 934 |
| 300 | 37396530365 | `c6cf66e648` | schedule | failure | `run300.log`（复用 `111-c1`，本腿自切） | 1 922 082 |
| 299 | 37257547572 | `c6cf66e648` | push | failure | `run299.log` | 2 088 835 |
| 294 | 37076437856 | `941805d0e1` | schedule | failure | `run294.log` | 1 838 679 |

命令：`gh run list --repo CarlosShao/wisp --workflow ci --limit 20 --json ...` → `ci20.json`（3 778 B，`wc -c`＝**3778**）。
下载：`gh run view --repo CarlosShao/wisp --log-failed <id> > run<n>.log 2> run<n>.err`，逐发 `rc=0`、`errsz=0`（六发连跑，无一次 `EOF`）。
步切法与去前缀见 `01-*.md`；切片行数：305=2190／304=2197／303=2201／302=2196／301=2196／300=2125／299=2129／294=1525。

## 步末四数（逐字，`portable-tests.sh: four numbers`）

| run | `=== RUN` | PASS | FAIL | SKIP |
|---|---|---|---|---|
| 305 | 335 | 230 | **6** | 2 |
| 304 | 335 | 230 | **7** | 1 |
| 303 | 335 | 229 | **8** | 1 |
| 302 | 335 | 231 | **6** | 1 |
| 301 | 335 | 231 | **6** | 1 |
| 300 | 323 | 220 | **6** | 1 |
| 299 | 323 | 219 | **7** | 1 |
| 294 | 261 | 167 | **12** | 1 |

⚠ 数的是**顶层** `--- FAIL`（`grep -E '^--- FAIL: '`）；作差前已按第 128 条剥掉汇总/计数行——`runtests.sh:` 与 `portable-tests.sh:` 那两行含 `FAIL=` 但形状不是 `--- FAIL: <名>`，不在名册里。

## 逐名色矩阵（F＝FAIL／P＝PASS／S＝SKIP／—＝该发不存在此名）

| 用例 | 305 | 304 | 303 | 302 | 301 | 300 | 299 | 294 |
|---|---|---|---|---|---|---|---|---|
| `TestRunPacketCarriesTheLoadedInstructionFiles` | F | F | F | F | F | F | F | F |
| `TestPanelHostRealWindowHopAndLifecycle` | F | F | F | F | F | F | F | F |
| `TestTicket223PermissionDeniedSitsInItsOwnSentence` | F | F | F | F | F | F | F | F |
| `TestAC1AlwaysBranchDoesNotRevertAHandEditedKey` | F | F | F | F | F | F | F | **P** |
| `TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card` | F | F | F | F | F | F | F | **P** |
| `TestPanelHostLatencyPercentilesAC2` | **S** | F | F | F | F | F | F | F |
| `TestAC14GoSideEvalPushReachesThePage` | F | F | F | **P** | **P** | **P** | **P** | F |
| `TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe` | S | S | S | S | S | S | S | S |
| `TestTicket223HandEditedFsLooseningCostsAnL2Card` | P | P | **F** | P | P | P | **F** | P |

## 同名不同坏法：`TestPanelHostRealWindowHopAndLifecycle` 的两种红句（逐字，本腿现切）

| run | `:660` 那一栏 | 红在哪一行 | 红句类别 |
|---|---|---|---|
| 305 | `-1.000 ms` | `:662`（`t.Fatalf`） | 哨兵值／通道没起，用例终止 |
| 304 | `3414.339 ms` | `:665`（`t.Errorf`） | 冷启 **3414.3 ms 超 D32 面板冷预算 1500 ms** |
| 303 | `3428.032 ms` | `:665` | 同上 |
| 302 | `3191.700 ms` | `:665` | 同上 |
| 301 | `3941.081 ms` | `:665` | 同上 |
| 300 | `3743.192 ms` | `:665` | 同上 |
| 299 | `3883.234 ms` | `:665` | 同上 |
| 294 | `3126.549 ms` | 旧行号 `:549`（该发树里行号不同） | 同上 |

⇒ **8 发里只有 1 发（305）是 `-1` 那形；7 发是"量到了、超预算"那形**。`111-c1` 的件（`ci-colours.md` Q5、票 111 追加段第 392 行）把 run 300 也写成"两发都是 `-1.000 ms` → `:662`"，**与原文冲突，本腿具名顶回**（见 `attribution.md` 顶回第 1 条）。

## SKIP↔FAIL 互换的机理（不是随机漂移）

`TestPanelHostLatencyPercentilesAC2` 的色由同一发里 lifecycle 那一枚**怎么红**决定：
- lifecycle 红在 `:665`（`t.Errorf`，非终止）⇒ 样本记上了 ⇒ 那枚跑分位并**红在预算**：
  `:993 sample 1: cold=3414.339 ms hot=35.152 ms` / `:999 AC#2 percentiles over 1 runs (nearest-rank): cold P50=3414.339 P95=3414.339 (budget 1500) | hot P50=35.152 P95=35.152 (budget 200)` / `:1005 cold P95 ... exceeds the D32 panel cold budget 1500 ms`（run 300/302/303/304 同形，hot 侧 25.996–41.580 ms **一律在 200 ms 预算内**）。
- lifecycle 红在 `:662`（`t.Fatalf`，终止）⇒ 样本没记上 ⇒ 那枚走 `:985` 具名 **SKIP**（run 305）。

⇒ 这**两枚是同一个因的两个果**，作差时不许当两枚独立信息量看。

## 另一族只出现在更早那发的红（run 294，12 枚里的 6 枚在后续发消失）

run 294 顶层红名（逐字顺序）：`TestL1VetoNeedsAChannelTheHostReallyWired (22.35s)`、`TestTicket223PermissionDeniedSitsInItsOwnSentence`、`TestRunPacketCarriesTheLoadedInstructionFiles`、`TestPanelHostRealWindowHopAndLifecycle`、`TestPanelHostLatencyPercentilesAC2`、`TestAC14GoSideEvalPushReachesThePage`、`TestAC246ResidentPipelineAsksThroughTheOneGate`、`TestTicket101ManualSwitchSurvivesRestart`、`TestTicket101UntouchedConfigRestartsAtDefault`、`TestTicket101SessionGrantDoesNotCrossRestart`、`TestComposedGateBlocksAWriteForTwoSeconds (301.19s)`、`TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking (41.23s)`。
⇒ 那一发的树比 `=== RUN=261` 还小、且带 301 s 这种量级的用例 ⇒ **名册的规模本身在漂**，不是只有尾巴在漂。

rc=0
