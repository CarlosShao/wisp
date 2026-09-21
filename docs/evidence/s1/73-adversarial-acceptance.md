# 票 73 对抗验收 —— 孤儿暂存文件清扫器（A18 判据②的可执行化）

**验收人**：编排者（**非实现者**；实现方 `agent-ticket73`，5 个 commit）
**验收时间**：2026-09-21 10:04
**被验收的 commit**：`a0072b0`（清扫器 + 翻转 A18 断言）、`4445911`（AC#2/3/4 用例）、
`7e8332e`（门禁数字）、`d49b0e6`（代理自我撤回一条假观察）
**基线**：`5849435`
**所有验收动作都在 `git archive HEAD` 解到 `/tmp` 的纯净树里做**——共享工作树里有四个代理在写码，
任何在仓内跑的变异都可能被邻居吞进他下一次 commit（今天已经为此付过一次 A31、一次 A34）。

## 裁决表（与票面 5 个 AC 框 1:1）

| AC | 裁决 | 我跑的判据 | 档位 |
|---|---|---|---|
| **AC#1** 翻转 A18：真 kill 之后下一次写盘必须回收残留 | **PASS** | 我的独立变异：把 `sweepStagingOrphans` 首行改成 `return stagingSweep{}`（＝"清扫器不存在"）→ **`--- FAIL: TestA18RealTaskkill…TheNextWriteReclaimsTheStagingFile`**，失败报文**点名残留全路径**。还原后 `grep -c MUTATION-73`=0、复跑 `ok`。 | 〔独立复现〕 |
| **AC#2** 归属是真的，不是"前缀匹配就删" | **PASS** | 我的独立变异：`stagingAttributable` 恒真 → **`--- FAIL: TestSweepNamingSchemeRejectsNearMisses`**（这条正是"差一点像我的名字"的守卫）。代理自己那次的变异点是 `fs_staging.go:224`（前缀归属），也红过。 | 〔独立复现〕（代理那一发为〔日志＋归档，我抽验〕） |
| **AC#3** 清扫**永不**穿过真 junction 删东西 | **PASS，但我要替它记一条更准的说法** | 见下面「变异矩阵」：单独关掉 reparse 门 **不红**，单独关掉 `IsRegular` 门也**不红**；**两门全下 + 前缀归属**（=现实里的"天真清扫器"）→ **三条安全用例全红**，含 `--- FAIL: TestSweepNeverDeletesThroughARealJunction (0.30s)`。 | 〔独立复现〕 |
| **AC#4** 活句柄不被删、写盘也不失败 | **PASS** | `TestSweepSparesATempFileHeldByAnotherProcess` 在我的每轮矩阵里都**先绿后红**：天真清扫器下红（0.01s），其余配置绿（0.76s）⇒ 用例既咬得住也跑得动，不是空转。 | 〔独立复现〕 |
| **AC#5** 门禁 | **PASS（部分为抽验）** | 我在纯净树跑的：`go test ./internal/tools/ -run 'A18\|Staging\|Orphan\|Sweep\|Junction\|Tmp'` → **19 RUN / 19 PASS / 0 FAIL+SKIP**、`ok 2.674s`。代理自报 `-count=2 ok 27.617s`、`-race -count=2 ok 36.283s`、`gofmt -l` 空、`vet` rc=0、`sh scripts/d22scan.sh` clean——**我没有重跑 race 与 -count=2**，按第二档记。 | 〔独立复现〕+〔日志＋归档，我抽验〕 |

## 变异矩阵（本票真正的交付物，全部我在同一棵纯净树里手动做的）

| 变异 | 落在哪 | 结果 | 读法 |
|---|---|---|---|
| M1 `sweepStagingOrphans` → 直接 return | 清扫器整体不存在 | **A18 红** | 护栏为真：没有清扫器就没人回收 |
| M2 `stagingAttributable` → 恒真 | 归属判定 | **命名方案用例红**（其余绿） | 名字过滤器在第一线，M2 没走到删除 |
| M3 `stagingIsReparse` → 恒假 | reparse 门 | **全绿** | ⚠ 见下 |
| M4 `stagingDeletable` 的 `!IsRegular` 去掉 | 常规文件门 | **全绿** | ⚠ 见下 |
| M5 = 前缀归属 + `stagingDeletable` 恒真 | "天真清扫器" | **三条安全用例全红** | 用例有效，保护是**分层冗余** |

**M3/M4 全绿该怎么定性（这条判断比结果本身值钱）**：它**不是**"测试无效"，而是
**同一件事有三层保证（名字切分 / `IsRegular` / reparse 属性），删任一层结果不变**。
M5 证明这三层一旦全撤，测试立刻三条齐红 ⇒ **测试咬得住的是"结果"，不是某一行实现**。
这与我记忆里那条判据对上了：**"两条互为冗余的保证，删任一条测试都不红"**——
此时必须分清"测试无效"与"冗余防御"，前者是缺陷，后者不是。

**但 M3 顺带抓到一处必须登记的东西**：`internal/tools/staging_live_windows.go:49` 的注释写着
"stagingIsReparse is **C26's rule applied to the sweeper's own traversal**"，
把 junction 安全**归因给 reparse 门**；实测真 junction 是**非 regular**，
所以先挡下它的是 `IsRegular()`。⚠ **错归因的危险是反的**：后人读注释会以为"reparse 门在兜"，
于是敢删 `IsRegular` 那行——而那行才是先起作用的。**登记，不改**（改注释属下一张碰 `internal/tools` 的票顺手做）。

## 结论

**PASS，5/5，归档 `-done`。** 本票的价值不止是"不再漏垃圾"，而是它把 A18 从
**一条在真实故障下不可观测的断言**（原测试用进程内 `Hooks.Kill`，返回 error 一定跑 Go 清理）
变成**真 `taskkill /F` 之下可证伪的不变式**，并且给了三层冗余的删除防护。

## 代理交回的六条，我的处置（不留口头账）

1. **旧命名方案的孤儿永久不可归属**（改名后 `.wisp-tmp-<digits>` 形状再也匹配不上）⇒ **真缺口**，
   需一个显式"legacy 名字清扫"决定 ⇒ 登记为 **Q-18**（owner：要不要为一次性历史残留写一段带截止的兼容清扫）。
2. **A18 可见残留计数 2 → 1**（每次写盘都扫）⇒ 引用过"两次 kill 留两个"的文档**全部过期**；
   `docs/evidence/s1/20-…` 作为历史证据**不动**，但我在此登记，避免下次被当成矛盾。
3. **A18 测试函数已改名** ⇒ 票 20 证据文档里的旧名字**不再解析** ⇒ 与 ②同批处理。
4. **测试在修好之前先发现了一个真 bug**：`pid == "0"` 曾被当成合法创建者
   （`OpenProcess(0)` 打开的是 System Idle Process ⇒ 那个名字**永远"活着"**、永远扫不掉）。现在匹配器拒了。
   ⇒ 这条是本票质量最高的信号，**写进 A42 表扬并保留**。
5. 共享树噪声非其造成（`tools/d22scan` 一度 `undefined: io`），且它**主动撤回了自己一次 `script-rc=0` 的误读**
   （`d49b0e6`）⇒ 认可：这正是 A15"没报错≠跑过了"的自查。
6. A36① 那两枚 Unicode 引号目录在它开工前就已被票 70 的代理删掉 ⇒ 已在 **A40⑤** 记过，不重复计。
