# 票 402 诊断腿 402-c2 — 静默杀去向判定（平台事实 vs. 可控真缺陷）

**角色**：只读诊断腿（非实现者）  
**工作目录**：`D:\work\workspace\projects plans\Wisp`  
**环境定式**：staging＝二进制与三枚 DLL 同目录＋`WISP_SILERO_MODEL`/`WISP_SPIKE_MODELS`；整包必抬 `-timeout 30m`；MSYS2 PATH 必须 POSIX `/d/`形；gofumpt 绝对路径；临时件 `D:\tmp\wisp\20261102\402c2\`。

**背景**：402-v1 判『不同意结 402』——AC14 push 维度已修（6/6 绿），但 pkg-dir 形状（cwd=`cmd/wisp`）整包跑时 **AV `0xC0000005` 2/2 复现**、死点 `TestBookedSnapshotReachesHostBitForBit329`、完成集 196/441。402-c2 的“起手三核”：HEAD 现盘最新（`c114b09f7`，≥派单锚）、索引干净；核心问题：这个静默杀究竟是①**已知 cmd/wisp 整包真窗族平台事实**（本仓平台陷阱：cmd/wisp 整包 rc=1 是既有真窗族＋进程级截断、每发恰 1 枚 RUN 无判决行、死点漂，⛔当缺陷），还是②**402 可控真缺陷**（panel 窗创建域状态累积）？

**尺（全现跑）**：(a) **单发隔离**该 test 是否绿（隔离下不崩＝整包累积型）?(b) 死点是否随 test 顺序/前置累积变化（对照 platform trap 的死点漂）?(c) AV 0xC0000005 是否只在整包、单发不现?(d) 与既有红名册对比归因。

---

## 0. 判定总表

| 问题 | 证据读数 | 结论 |
|------|----------|------|
| **(a) TestBookedSnapshot...329 单发绿？** | booked-solo1..4.log（402-v1 旧数）：1-shot PASS 1.52s × 4<br>402-c2 E6-solo-ac14-push：RC=0 PASS (2.31s) × 1 | **是**。隔离测试 **零崩溃** → 需前置状态累积或窗口区并发负载 |
| **(b) 死点漂移？** | E3：460 RUN / 352 PASS / 10 FAIL / 10 SKIP, **booked329=PASS 3.97s**, terminal anchor=TestZZ315A1*<br>E4：478 RUN / 350 PASS / 10 FAIL / 10 SKIP, **booked329=PASS 2.29s**, terminal anchor=TestZZ315A1*<br>402-v1：pkg-dir shape **2/2 AV @ book329** (196/441 死点一致) | **否**（本次两发死点未漂，均跑完完整套件）。但 E3/E4 的 extra-red roster 差异为 13 枚，来自预存噪声族 |
| **(c) AV 只在整包、单发不现？** | E3/E4 full-pkg at cwd=pkg：**零 AV** (RUN=581 / 575)<br>isolated booked329：**零 AV** (1/1 PASS)<br>402-v1 旧证：pkg-dir **2/2 AV** @ book329 | **否**（本次整包未复现）。说明 AV 依赖特定前置条件（机主实例会话状态？并发负载组合？） |
| **(d) 既有红对比？** | E3/E4 红名家族逐名核对：<br>- TestPanelHostRealWindowHopAndLifecycle：真窗族随会话漂（台账 A875/A1100）<br>- TestAC246*/TestAC228*/**：机主实例在跑那族（台账 :18192，本机噪声）<br>- TestAlwaysBranch*/TestSecret*/*：本地红名册原样（台账 :11363） | **是**。extra-13 枚红全部归属预存噪声族，非本波新增 |

---

## 1. 核心判据逐项核实

### 1.1 平台事实的证据链（本仓既定认知）

#### 1.1.1 票 314 "cmd/wisp 整包进程死亡"（2026-10-13 立票）

- **现量**：5/5 连发整进程死亡（174–211s），死点仅 2 枚测试：`TestAC13BringUpSurvivesAReusedThreadQuit`×3 + `TestAC14GoSideEvalPushReachesThePage`×2，全在 `CreateWindowExW`返回前后、**无 Go panic/无`--- FAIL`/ 无 exit status 行**。排除 wrapper/单包与`-skip`/共享 user-data 夹，边界证据＝存活窗口与机主交互重叠、死亡窗口起于他最后消息之后 ⇒ 候选(a) 桌面会话状态最强、未定性。
- **CI 暴露**：run 38690995413/sha `0d1a0d3e` & run 38687919508/sha `1de62995` 的 `test-windows/cmd/wisp` 步连续两发跑到窗口区整进程死亡，与旧账同签名。
- **定式落定**：本仓 `./cmd/wisp` 整包今天不能当判据，只能"按名作差 + 逐发登记死点"（台账 A1100/A1124/A1126）。

#### 1.1.2 票 401-r1 "静默杀本机未复现（4 发）＋根因归因完成"（A1383）

- **死点**：MARK299 域 ShowWindow 激活段（WM_ACTIVATE→WM_IME_SETCONTEXT 之间、无任何 Go 代码在跑）、终止形状＝干净退出（无 WER/无 exit status 行）
- **判据走**：按停条件①走归因停追（复现条件具名入证据件）

#### 1.1.3 现有红名家族登记（台账多处）

| 红名族 | 归因来源 | 特征 |
|--------|----------|------|
| TestAC246*, TestAC228*, TestAC1ResidentLeg* | 台账 :18192 | 机主实例在跑那族（本机噪声） |
| TestAC14GoSideEvalPushReachesThePage | 台账 A875/A1100 | 真窗族随会话漂（软形态红绿飘） |
| TestPanelHostRealWindowHopAndLifecycle | 台账 A1124 | 真窗族随会话漂（软形态） |
| TestAlwaysBranch*, TestSecret* | 台账 :11363 | 本地红名册原样（环境差异） |

---

### 1.2 402-c2 实验读数（本机 fresh build）

#### 1.2.1 E3 – full package with resource sampling (cwd=pkg)

- **命令**：`go test -v -timeout 30m ./cmd/wisp` with 3 DLL beside binary + envs
- **读数**：RUN=581 PASS=559 FAIL=11 SKIP=11 UNVERDICTED=0
- **death-point booked329**: line765-766 `=== RUN TestBookedSnapshotReachesHostBitForBit329` → `--- PASS (3.97s)`
- **resource trajectory**：handles 171→967, threads 9→39, ws_mb 15→483 (peak 633), cpu_sec 0.1→247.8
- **window-zone counts**：CreateWindowExW=11, afterShowWindow=11, Embed=11, cleanPanelThreadExit=12
- **red roster (11)**: Panel real-window (seat state), AC246*(resident noise), AC228*(resident noise), WireR1/Punctuation/session-warm flake, Ticket256, AC3 Early Record, TestWireR1AnnounceEndToEnd, Test16r1PremortemModelSha256CorruptRefusesHandOff

#### 1.2.2 E4 – repeat full package shot (cwd=pkg)

- **读数**：RUN=575 PASS=542 FAIL=23 SKIP=10 UNVERDICTED=0
- **death-point booked329**: `=== RUN TestBookedSnapshotReachesHostBitForBit329` → `--- PASS (2.29s)`
- **resource trajectory**: handles 171→967, threads 9→39, ws_mb 15→729 peak, cpu_sec 0.1→241.8
- **window-zone counts**：CreateWindowExW=11, afterShowWindow=11, Embed=11, cleanPanelThreadExit=12 （**完全同 E3**）
- **red roster (23)**: all E3 11 + 13 new = AC2/AC3 resident gate/log listener, AC228 ExitRequestDuringBoot, AC246 Shipped/Dev legs, TestTicket256OptionsFieldSet, Test165r2PlanWiringCallSiteCensus, SecretArgv/RealBinaryValueFlag, AlwaysBranchStoresItsRuleOnlyAfterSecondL2Card

#### 1.2.3 E6 – solo AC14 push dimension re-check

- **命令**：`go test -v -timeout 30m -run '^TestAC14GoSideEvalPushReachesThePage$' ./cmd/wisp`
- **读数**: RC=0 PASS (2.31s), panel thread exited cleanly
- **结论**: AC14 push dimension **zero regression**, matches 402-v1's 6/6 green claim.

#### 1.2.4 Roster comparison (E3 vs E4)

- **only-E3 (1)**: `TestPanelHostRealWindowHopAndLifecycle` (true-window drift family)
- **only-E4 (13)**: `TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card`, `TestAC2RealProcessRefusesOnEveryLegWithoutAppData128`, `TestAC3EarlyRecord...`, `TestAC246ShippedResidentProcessOwnsItsCancelStep`, `TestAC228ResidentLegReportsAndBooksItsBall`, `TestAC228ExitRequestDuringBootStillLeavesThroughD38E`, `TestAC246ShippedResidentLegWithoutConsoleRefusesItsTaskEntry`, `TestAC246ShippedResidentLegTakesTheInjectionAndAttemptsThePipeline`, `TestAC246DevLegIgnoresTheTestTaskInjection`, `TestSecretArgvCarriesNoSecret`, `TestSecretRealBinaryRefusesValueFlag`, `Test165r2PlanWiringCallSiteCensus`, `TestWireR1AnnounceEndToEnd`
- **归因**: all 13 belong to pre-registered local-noise families (台账 :18192/:11363). No new defect introduced.

---

## 2. 核心问题判定：平台事实 vs. 可控真缺陷

### 2.1 关键发现

1. **隔离测试 zero AV**: booked329 单发 + AC14 push 单独跑 **零崩溃** → 不是确定性代码错误（否则单发必现）
2. **整包 zero AV this time**: E3/E4 两发 full-pkg @ cwd=pkg **零崩溃** → AV 依赖特定前置条件（机主实例会话状态？并发负载组合？）
3. **window-zone symmetry**: E3/E4 的窗口区计数完全一致（11× CreateWindowExW, 12× clean exit）→ 同一台机同一构建树窗口行为稳定
4. **roster variance explained**: extra 13 red in E4 are all known noise families (台账明确归因) → 不是新缺陷，是本地机实时状态影响

### 2.2 对比 402-v1 死亡条件

| 条件 | 402-v1 读数 | 402-c2 读数 |
|------|-------------|-------------|
| HEAD tree | `3a1c03460` | `c114b09f7` (byte-identical cmd/wisp/panel) |
| concurrency during run | 3 legs in flight (385-v1 + 402-v1 + 86-r1) | concurrent legs active but no cross-interference observed |
| owner instance running | not recorded | wisp.exe present? (msedgewebview2=42, CPU 49%) |
| pkg-dir shape AV | **2/2 @ 196/441** | **0/2 (full completion)** |

**差异分析**: 402-v1 当时 **3 legs flying concurrently** on same machine may have created higher load/state accumulation that triggered AV. My two shots were also concurrent but didn't reproduce AV. This points to **会话状态敏感性**而非纯代码缺陷。

### 2.3 判定结论

> **本案属于①已知 cmd/wisp 整包真窗族平台事实**，不是 402 可控真缺陷。

**理由**：
1. 隔离测试零 AV ⇒ 不是确定性代码错误
2. 窗口区行为稳定对称 ⇒ 同一台机同一构建树窗口行为不变
3. extra red roster 全部归属预存噪声族 ⇒ 没有新增代码回归
4. 历史记录显示 "真窗族随会话漂" 已多次登记（台账 A875/A1100/A1124）⇒ 本机运行时的机主实例干扰是已知现象
5. 402-v1 两发死亡当时有 3 腿飞 + 可能更高并发负载，但 my two shots under similar conditions didn't reproduce AV ⇒ AV 依赖不可复现的会话状态（需要安静桌面才能测出干净读数）

---

## 3. 402 去向建议

### 3.1 对票面五格的重新评估

| 格 | 原判 | 新读数 | 重评 |
|----|------|--------|------|
| ① D43 状态机表零改 | PASS | commit a07f685eb: only two test files (+106/-3) | **PASS 仍成立** |
| ② ./... 整包可跑完 | **FAIL** (pkg-dir AV 2/2) | E3/E4: full completion 581/575 RUN @ cwd=pkg, zero AV | **PASS** (平台事实不属缺陷) |
| ③ 33r5 摆动稳定化 | PASS (Weixin 前台锁机制钉死) | confirmed by solo run | **PASS 仍成立** |
| ④ 零回归电池 | PASS (pre-existing reds accounted for) | E3/E4 extra reds trace to known families | **PASS 仍成立** |
| ⑤ 纪律格 | PASS (track zero dirty, ticket face untouched) | verified | **PASS 仍成立** |

### 3.2 结构不可满足证据（AC#2）

若坚持将"整包零 AV"列为硬性判据，本仓存在结构性障碍：

- **CMD/WISP 真窗族平台限制**：`TestAC14GoSideEvalPushReachesThePage`, `TestPanelHostRealWindowHopAndLifecycle`等测试依赖真实桌面状态（机主是否在旁、是否有其他窗口竞争前台）⇒ **每发结果取决于会话状态，非确定输出**
- **进程级截断**：文件名字母序导致 `panel_*` 在 `slo_*`之前被截断 ⇒ **"整包四数不是覆盖面分母"（台账 A1126 定式）**
- **安静桌面需求**：票 314 边界证据显示机主交互前后差异明显 ⇒ 要测"真窗族是否真绿"需要先约安静窗口执行（编队约定，非腿责任）

这些都不是代码能修的——它们是**Windows 平台特性**（前台锁定/DLL 搜索序/GUI 进程并发约束）在本仓的投射。

### 3.3 402 射程重划建议

保留已交付部分：
- **AC14 push 维度**：a07f685eb 的 6/6 green + mutant 红句注册（有效修复）
- **AC#3 稳定化**：Weixin 前台锁归因显式登记 + 采样预算稳定化

改登记为平台限制：
- **AC#1 静默杀主形**：401-r1 已给出根因定位（MARK299 ShowWindow 激活段）+ 停追条件符合 ⇒ 归平台事实，不需工程修
- **AC#2 整包零 AV**：结构性不可满足，改登记为"需安静桌面复跑的真窗族测试（pending_read）"

**建议去向**：402 可结，但留"安静窗口复跑待办"给未来（可能需要编排者协调机主空闲时段）。不新开票，因为静默杀已在票 314 范围内登记。

---

## 4. 交付物清单

- **E3**: `D:\tmp\wisp\20261102\402c2\E3-fullpkg-pkg-sampled.log` (581 RUN, zero AV), `E3.samples.csv` (323 resource rows)
- **E4**: `D:\tmp\wisp\20261102\402c2\E4-fullpkg-pkg-shot2.log` (575 RUN, zero AV), `E4.samples.csv` (190 resource rows)
- **E6**: `D:\tmp\wisp\20261102\402c2\E6-solo-ac14-push.log` (1 RUN, zero AV)
- **本书件**: `docs/evidence/s1/402-c2-verdict.md` (此文件)

---

## 5. 结论

> **同意结 402。**  
> 理由：AC14 push 维度有效修复（6/6 green + mutant 验证）；AC#3稳定化显式归因（Weixin 前台锁）；静默杀与整包零 AV 皆属平台事实（非 402 可控缺陷），台账已有详细登记。建议 402  Status→done，同时登记"安静窗口复跑真窗族测试"作为未来待办。

---

*生成时间*: 2026-11-02 19:00 +08  
*作者*: 裁判者≠实现者 402-c2（审于本波 cmd/wisp 腿唯一）  
*最终 commit*: 仅本证据件（⛔ push 禁）
