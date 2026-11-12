# 431 — 两枚 `TestAC1ResidentLeg*` 红在同一件事：**启动预算按固定秒数写死**（10s 窗口 vs 本机 3.4–5.6s 启动；1s 轮询 vs 2.9–5.6s 启动）

- Status: **done**（写面＝`cmd/wisp/resident_sink_nail_127_windows_test.go` 与 `cmd/wisp/resident_log_listener_127_windows_test.go` 的**预算常量/等待方式**；⚠ 先只读取证，判清是"预算该加大"还是"该改成事件驱动"，再定改法）
- 来源：round 88 整包复跑（`D:\tmp\wisp\20261112\fullclean\cmdwisp.log`，`PASS=593 / FAIL=14`）
- 关联：票 127（resident sink nail，两枚测试宿主）· 票 43（十步回放）· 票 425（同族第三枚，已修）· 台账 `A1801`／`A1805`

## 这是什么（人话）

有两枚测试都在测**常驻进程（resident leg）启动时**的行为，但都**只给了一个固定的短时间窗口**，
而本机（这台机器）启动一次要 **3.4–5.6 秒**——窗口不够，于是它们**在机器还没启动完就判红**。

**这⛔ 是产品坏了，是"尺子刻错了"。** 但⛔ 也不能简单把秒数加大——那只是把偶发变少见。

## 两枚的现量（逐字）

### ① `TestAC1ResidentLegInstallsItsLogListenerOnDisk`
**逐字红句**：预算 10s 窗口不够（本机启动 3.4–5.6s，其中 `speech.NewDetector` 一枚就 1.83s，见 `cmd/wisp/resident_voice_windows.go:713`）。
⇒ **稳定红**（本会话多轮整包均红）。

### ② `TestAC1ResidentLegOutlivesItsOwnLogFailure`
**逐字红句**：`pollUntil127(20, …)` 给 **1 秒**预算（20 × 50ms），而启动要 2.9–5.6s。
⇒ **⚠ 偶发**：腿在 round 88 测过 5 次 ⇒ **3 PASS / 2 FAIL**。

**关键差别**：①是**稳定红**，②是**偶发**。⛔ 用同一个修法糊过去——先判清两者是不是同一个根。

## 机制链（现读，供你起步；⛔ 别照抄结论）

- `pollUntil127(n, cond)`：轮询 `n` 次、每次间隔约 50ms ⇒ **`n=20` 只有 1 秒**。
- 启动成本大头：`speech.NewDetector` 在 `cmd/wisp/resident_voice_windows.go:713`，**1.83s**。
- 同族第三枚 `TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink` 已由**票 425** 修好（`e4c611d26`），它的红因不同（断言落后于产码，⛔ 预算）。**⛔ 把这三枚混为一谈。**

## ⚠ 必须自己判清的核心问题（⛔ 不许我替你定）

**"把秒数加大" vs "改成事件驱动/等真实信号"——哪个是主流 harness 的做法？**

- **加大秒数**：简单，但**仍然是一个魔数**，换台更慢的机器还会红；而且会让真故障时**多等那么久**才报。
- **等真实信号**：例如等 `slog` 里那条"启动完成"的**逐字行**出现（`pollUntil127` 本来就支持条件函数），
  或等某个 `ready` 通道关闭。**这才是主流做法**：等**事实**而不是等**时间**。
- **中间形**：把预算按"启动实测成本 × 安全系数"算出来，并在注释里写明依据（⛔ 拍脑袋）。

**你必须先量出**：这两枚各自**真正在等什么信号**、那个信号**最早何时可观测**。
然后选形并**写明理由**（⛔ "因为简单"不是理由）。

## 要做的（按序）

- [x] **AC#0 只读取证（⛔ 本格不许改任何文件）**：
  1. 两枚各**复跑 ≥3 次**定性（确认①稳定红、②偶发，并给出每次耗时）：
     ```
     D:\work\soft\Git\bin\bash.exe -c "scripts/wisp-cli-staged-tests.sh '-test.run' '^TestAC1ResidentLegInstallsItsLogListenerOnDisk$'"
     D:\work\soft\Git\bin\bash.exe -c "scripts/wisp-cli-staged-tests.sh '-test.run' '^TestAC1ResidentLegOutlivesItsOwnLogFailure$'"
     ```
     ⚠ **必须用 `D:\work\soft\Git\bin\bash.exe`**（本机 `bash` 是 WSL 转发器，一枚测试都不跑）。
  2. **逐字抄出两枚的红句**（日志里被截断的要回源码抄）。
  3. 读 `pollUntil127` 的实现与两枚各自的**等待条件**，写清它们**真正在等什么**。
  4. **量出**那个信号最早何时可观测（打时间戳），以及 `speech.NewDetector` 的实际耗时。
  5. 判清①与②**是否同一个根**，给一句话结论＋`文件:行`。

- [x] **AC#1 按 AC#0 的判定修（⛔ 不许只把秒数乘个系数了事）**：
  - 优先改成**等真实信号**（事件/日志行/ready 通道）；若确实只能靠时间，必须**写明依据**（实测成本 × 系数）。
  - ⛔ **不许放宽断言本身**（要等的那个事实⛔ 改）。
  - **必须配突变控制真跑**：把产码里那个"启动完成"的标记删掉 ⇒ 新等待**必须超时红**（证明它等的是真事实，⛔ 空转）。

- [x] **AC#2 越界检查**：`git diff --name-only` 逐行贴。出现 `docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**` ⇒ **直接退回**。

- [x] **AC#3 门禁**：`go build ./cmd/wisp/` · `go vet ./cmd/wisp/` · `gofumpt -l` 零输出 · **`scripts/d22scan.sh`**（⚠ ⛔ 用 `go run ./tools/d22scan`，独立 module）。四条读数全贴。
  ⚠ `build/mut_rvw.go` 是未跟踪残留件，会让 `go build ./...` 红。**⛔ 不许删**。

- [x] **AC#4 与票 425 的关系写清**：本票只处理**预算/等待方式**这两枚，⛔ 不翻票 127／425 任何框；⛔ 不宣称 `TestAC1*` 全族绿。

## 本票**不**解决

- 不碰 `Test26P*`（票 428 已修）／`Test384*`（票 429 已修）／`Test365*`（票 430 已修）。
- 不改 `thresholds.go`／SLO 阈值（冻结件）。
- 不改 `docs/PLAN.md` 与 `docs/specs/**` 一字。

## 派单约束

- 先只读（AC#0 零写），判定后再定写面。
- ⛔ `git add -A`／`add .`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；`add`／`commit` 两处显式 pathspec。
- 只 commit，⛔ push。
- ⛔ 翻任何 AC 框。
- 临时件只建不删，落 `D:\tmp\wisp\<YYYYMMDD>\431-<腿码>\`。
- **仪器**：运行时行为检查 ⇒ `go test -overlay`（拷贝放**仓外**）有效。

## Progress log

- **2026-11-12 22:1x +08 — 写腿（431-ac1budget）**：AC#0 只读取证完成，**票面一处假设被实测推翻**：①的红**不是** 10s 窗口不够——它的盘上记录轮询 **611ms 就成功**（预算从未绑住），红在 `resident_sink_nail_127_windows_test.go:505`，因为用例在 ~0.6s 就 `leg.stop()` 了子进程，而它要读的 `resident event loop running` 要到 **2.8–4.3s** 才出现（16 发实测 2.77–4.32s，中位 ~3.4s；`speech.NewDetector` 实测 **2.4–2.6s**，非票面写的 1.83s）。②才是真预算病：`pollUntil127(20,…)`=1s vs 启动 2.84–2.99s。⇒ **两枚不同根，同一修法：等事实不等时间。**
  选形＝**乙（等真实信号）**，⛔ 加大秒数、⛔ 算系数——因为甲/丙**买不到①的绿**（事实在 `stop()` 之后永不到来），且只会让真故障多等。新增 `bootFactCeiling127 = 60s`（**上限**，非等待，注释写明 16 发实测依据）＋ `waitForBootFact127`（`observe.NewTimeout` 单调钟；事实一到立刻返回；**子进程一死立刻报**，⛔ 坐满上限）。断言事实**逐字未改**，无 `t.Skip`，⛔ 碰 `thresholds.go`/golden。
  **突变控制真跑**（产码 `resident_windows.go:851` 那句标记删除，经仓外 `-overlay` 映射，仓内零字节）：① `FAIL (66.69s)` `the child never reached the event loop within 1m0s`、② `FAIL (66.90s)` `the refusal stopped the app: … within 1m0s` ⇒ **超时红，不是空转绿**；同发里**不相关的那个事实照常到达**（盘上记录 609ms / 拒绝句 151ms），证明各自钉的是自己那句话。
  改后：① 3/3 PASS（10.65 / 10.12 / 10.32s）、② 3/3 PASS（10.02 / 9.32 / 9.67s）；同族三枚同跑全绿（含票 425 那枚）。
  门禁：`go build ./cmd/wisp/` rc=0 · `go vet ./cmd/wisp/` rc=0 · `gofumpt -l cmd internal` 零输出 · `scripts/d22scan.sh` rc=0（`PASS=45 FAIL=0 SKIP=0`，`clean - no D22 ban violations`）。
  证据：`docs/evidence/s1/431-ac1-boot-budget.md`。⛔ 未翻任何 AC 框；⛔ 未 push。
- [2026-11-12 22:1x +08] **编排者：收件＋三把尺全过＋AC#0–AC#4 五格翻勾＋`-done`**。腿 `ticket-431-ac1-boot-budget` 交件（`fa39931bd`，3 files；证据件 `docs/evidence/s1/431-ac1-boot-budget.md` 16815B）。**三把尺**：①`git cat-file -t` = commit；②件 16815B 在盘；③**我自己复跑**两枚 ⇒ **PASS 10.88s / PASS 7.94s，FAIL=0**。**越界**：只改 `cmd/wisp/resident_sink_nail_127_windows_test.go` 一枚。
- **★★腿推翻了我票面的因果故事：① 的 10 秒预算根本没生效**。腿的实测：那 10s 轮询**在 611ms 就成功了**，预算从未绑住。真正的红在 **`:505`——下一条断言**，它在 `leg.stop()` **已经把子进程杀掉之后**才去读 `resident event loop running` 那一行。**那条事实在 `stop()` 之后永远不会出现** ⇒ **加多少秒都救不了**。腿用仓外 overlay 证明：把等待移到 `stop()` 之前、断言一字不动 ⇒ **PASS 8.20s**。
- **★两个不同的根，⛔ 一个**：① **读得太早**（顺序错），② **预算真的不够**（boot 2.84–2.99s vs 1s 窗口）。腿的定性：① **稳定红 5/5**（但⛔ 因为预算）；② **偶发**——round 88 读 3 PASS/2 FAIL、本会话读 5 FAIL/5，**两者都为真且一致**：它是对 2.84–2.99s 启动的**硬币**。
- **★各事实首次可观测时刻（腿实测）**：① 盘上安装记录 **611ms**；① 真正的红点 `resident event loop running` **2.36–2.59s**；② 拒绝行 **100ms**；② 启动完成行 **2.84–2.99s**。**启动完成行跨 16 次独立启动：2.77–4.32s，中位 ~3.4s**。`speech.NewDetector` **2.4–2.6s**（**我票面写的 1.83s 偏低**）；去掉语音夹具后启动降到 1.75–1.84s。
- **★选形＝等事实**（带单调上限作为唯一的时间量），理由⛔ "因为简单"：**两个更便宜的形根本买不到 ① 的绿**——那条事实在 `stop()` 之后永不出现，**乘秒数什么都不修**；对 ② 系数只让偶发变罕见。腿**加了第二个出口**：上限**只被"活着但永不到达"的进程**花掉；启动期就死掉的子进程**在 reaped 那一刻立刻报**——这正面回答了我票面自己提的"加大数字会让真故障多等"那条反对意见。**事实逐字节未改，无 `t.Skip`，断言未放宽。**
- **突变控制两发真跑**（仓外 overlay，产码零改动）：删掉产码里那个标记 ⇒ ① **FAIL 66.69s**（`the child never reached the event loop within 1m0s`）／② **FAIL 66.90s**（`the refusal stopped the app: … within 1m0s`）。**两发都是超时红，⛔ 空转绿**。同一次运行里各自等的**另一个**无关事实仍然到达（609ms 盘上记录、151ms 拒绝行）⇒ **证明每枚钉住的是自己那句话**。
- **★我票面的第 2 处错（第 7 次同根）**：我写「本机启动 3.4–5.6s，`speech.NewDetector` 1.83s」——腿实测 **2.77–4.32s（中位 3.4s）／2.4–2.6s**。我的数字偏高/偏低都不准，是从旧台账印象里填的。

