# 票 314 — `cmd/wisp` 测试二进制在真窗测试区**整进程死亡**：4/4 连发复现（16:54–17:17 窗口）；同日早些时候同形 4/4 全活（16:22–16:50）——未定性，归因候选已列

**立票**：2026-10-13 17:1x 编排者（来路＝亲跑票 311 `Q-311-r1b-1` 的 q1＋顺手复跑 q2/q3/q4；票 311 那格只判"守卫一声不响＋四数行在场"，死亡**出那格射程**、另立本票）
**性质**：★**事故票（测试台面稳性）**。与票 310（死掉的二进制留下临时目录）同因不同面：310 管死后残留，本票管**死本身**。⚠ **未定性为 299 回归**——归因候选 4 枚见 §现量末；⛔ 任何结论先过 `AC#0/#1`。

**为什么现在立**：① `ci.yml:621`（`run: bash scripts/wisp-cli-tests.sh`）那一跳就是默认形 ⇒ **下一次推送，CI 会看见这个死亡**；② 4/4 连发复现，且 q1 与 q4 死在**同一枚测试**（`TestAC14GoSideEvalPushReachesThePage`）。

## 现量（逐发；原始日志根 `C:\Users\swq\AppData\Local\wisp-test-gotmp\orch311-q-20261013\logs\`；⛔ 引用前重跑）

| 发 | 调用形 | 读数 | 死点（最后一行 → ⛔ 无下一行） |
|---|---|---|---|
| q1 | `wisp-cli-tests.sh` 默认形（单包＋CI 那份 `-skip`） | 死于 **185.422s**；`runtests` 四数＝`=== RUN=176 --- PASS=126 --- FAIL=0 --- SKIP=0`（**无任何 `--- FAIL` 行**） | `=== RUN TestAC14GoSideEvalPushReachesThePage` → 窗建起（`MARK299 after CreateWindowExW hwnd=0x17950e52`＋数条 wndproc）→ 死 |
| q2 | 同上复跑 | 死于 **187.924s**；四数 `RUN=172 PASS=120 FAIL=2 SKIP=0`（两枚具名红＝`TestTicket223ModeLooseningChangesTheRunningModeAfterAllow` 2.71s／`TestAC4FocusReturnToPriorWindowGap33r5` 0.35s） | `TestAC13BringUpSurvivesAReusedThreadQuit` → `after CreateWindowExW hwnd=0x2cd0ed4`＋`msg=0x86` → 死 |
| q3 | **直发**（`go test -count=1 -v ./cmd/wisp/ -skip …`；⛔ 走 wrapper） | 死于 **211.592s**；我 grep＝`RUN=173 PASS=170`＋1 枚具名红（`TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card` 41.24s） | `TestAC13BringUpSurvivesAReusedThreadQuit` → `after CreateWindowExW hwnd=0xeeb0382` → 死（更早一帧） |
| q4 | **两包**（`./cmd/wisp/ ./internal/audio/`，无 `-skip`＝与 16:22 活的那 4 发**同形**） | 死于 **173.961s**（cmd/wisp）；`ok internal/audio 16.162s` 照常 | `TestAC14GoSideEvalPushReachesThePage` →（同 q1） |
| **对照**：16:22–16:50 的 4 发 | 同 q4 形 | **4/4 全活**（cmd/wisp 各 397–407s；恰两枚真窗红＝软形态） | — |
| q5 | 形②（`WISP_TEST_GOTMP_ROOT` 指同卷换目录）＋**user-data 夹先挪走**（fresh profile，跑中重生成） | 死于 **178.028s**；守卫**静**（GUARD=0／INCONCLUSIVE=0）、**换目录被认领**（`GOTMPDIR=…/orch311-q-20261013/root2/cli-20261013-171826-1438` 逐字）、四数行在（`RUN=172 PASS=121 FAIL=1 SKIP=0`；1 枚具名红＝`TestTicket223ModeLoosening…` 2.27s） | `TestAC13BringUpSurvivesAReusedThreadQuit` →（同 q2/q3） |

- ★**调用形与共享夹的嫌疑全排除**：q3 干掉 wrapper（直发照样死）、q4 干掉"单包/`-skip`"差异（同形照样死）、**q5 干掉"共享 user-data 夹状态"主嫌**（夹先挪走、跑中重生成，照样死）⇒ 当前候选（按嫌疑排序）：
  (a) ★**桌面/会话状态**〔现最强〕：**边界证据＝存活窗口 16:22–16:50 与机主在旁/交互重叠；死亡窗口 5/5 起于他 16:50 最后一条消息之后（此后无交互）**——候选机制＝锁屏/空闲桌面对真窗测试的影响（⛔ 未验；下一发便宜实验＝他人在时复跑同形一发）；
  (b) 299 宿主（`panel_control_windows.go`）／新前端页面字节在真窗区的负载（⛔ 未验；死点只落在 2 枚测试＝更像"这两枚附近"而非全测试面）；
  (c) 随机性（5 死 4 活若 p=0.5 是 1/512 ⇒ ⛔ 当纯巧合，但也未到定性线）。
- ⚠ **q5 副产物（正向）**：该发同时是票 311 形②的闭档读数——守卫静＋换目录被认领＋四数行在 ⇒ 形② 判据成立。
- 死法共同特征：**无 Go panic 输出、无 `--- FAIL`（q1/q4 一枚都⛔ 有）、无 `exit status` 行**——整进程消失；死点全落在 `CreateWindowExW` 返回到首几条 wndproc 之间。
- 前例与对照：`A773`（10-09）**同脚本跑完**（`PASS=260/FAIL=4`）；`pending-and-issues.md:14768` 那族"真窗族不是确定性基线"是**软形态**（红绿漂），本票是**硬形态**（整进程走）。
- 关联票：票 310（死的后果＝临时目录残留；本票是死本身）；票 303（页面→Go 回执断口，已 `-done`）。

## 要建什么（草案；落地腿定最终形）

- [ ] **AC#0 死法取证**：`go test -c` 编出 cmd/wisp 测试二进制 → **直跑它**（`-test.v -test.skip='…'`），拿**进程退出码**（native 异常码 vs `os.Exit(n)` vs `-test.timeout` panic）＋stderr 全量；⛔ 只看 wrapper 汇总。
- [ ] **AC#1 最小重现＋归因 (a)**：q5 A/B 式（fresh profile ⇄ 旧 profile 各 ≥2 发）；判"(a) 共享夹状态"是否充分；若成立 ⇒ 修法候选＝夹具改用每发独立 profile（先查它为什么是固定名：`:61` 处注释逐字）。
- [ ] **AC#2 归因 (b)/(d)**：安静桌面（⛔ 人不用机器／⛔ 挂窗工具退场）下 1-2 发；以及 `main` 对照树 1-2 发（⛔ 动共享工作树；导出树跑）。⚠ 涉及"机主别用机器"的窗口＝要编排者约，⛔ 腿自约。
- [ ] **AC#3 CI 暴露与推送前门**：现读 `ci.yml:601-621` 一带确认哪一跳跑它＋**推送前**按本票形状重跑并把颜色写进推送前检查（⛔ 没跑过不许推）。⚠ 首推会让 CI 见生死——先在本机把死法定性再推。
- [ ] **AC#4 处置**：修（谁）／登记（多长）；⛔ 为变绿放宽任何断言、⛔ 动阈值／golden。

## 排程与边界

- 本票的"跑"（`AC#1/#2`）要 `cmd/wisp` 整包色＋静桌面 ⇒ ⛔ 与 `294-r1`·任何在飞整包腿同波。
- 写面（未来落地腿）：`probes/314/**`＋必要时 `cmd/wisp/*_test.go`（夹具 profile 名一处）；⛔ 本票⛔ 动 `.github/**` 一字。
- **不阻塞机主那条线**（他的产品面正常；死的是测试台面）⇒ ⛔ 拿它去打扰他；推送计划受影响（见 `AC#3`）。

## Progress log

（追加式；每条 `- [YYYY-MM-DD HH:MM:SS +08] agent=… did=… next=…`，钟点由 `date` 的 stdout 插值，不手打。）

- [2026-10-13 17:22:56 +08] agent=编排者 did=立票即连发 5 枚（q1–q5）**全部整进程死亡**（174.0–211.6s；死点仅 2 枚测试：3×`TestAC13BringUpSurvivesAReusedThreadQuit`＋2×`TestAC14GoSideEvalPushReachesThePage`）；排除面＝wrapper（q3 直发照死）／单包与 `-skip`（q4 两包无 skip 照死）／共享 user-data 夹（q5 挪夹重生成照死）；对照＝16:22–16:50 同形 4/4 全活 ⇒ 边界＝机主 16:50 最后交互前后（候选 (a) 桌面会话状态，现最强、未验）next=机主在场时复跑同形一发（便宜实验）；`AC#0`（直跑测试二进制取退出码）随下一枚能占 `cmd/wisp` 色的腿走；原始日志 `orch311-q-20261013/logs/{qshot-default,qshot-default-2,qshot-direct,q4-two-package,qshot-shape2}.txt`（⛔ 入仓，同根另存）
