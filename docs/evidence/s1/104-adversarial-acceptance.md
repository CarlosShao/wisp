# 104 — 对抗验收：`SealFile` 对"继承来的"带外授权出声（两桶 + `kind=`）

**验收方**：`acceptor-ticket104`（只读角色，一行生产码未改）
**日期**：2026-09-21 20:2x–21:0x（本机 CST）
**被验版本（锚定 sha）**：
- `4693feb517e0c7f59143a44696bb20706a8d4e31` = 修前红（AC#1 用例落地，行为仍 HEAD 原样）
- `4d4344798ab708bc69b0ea4756d481ebb580d522` = 修复（只动 `internal/winsec/winsec_windows.go` + 票面）

**当前 HEAD = `a505607`（票 112 已进），工作树另有票 112/113/111/92b 在飞 ⇒ 本票所有读数**不取**当前工作树。**

**快照口径（A38④：仓内不建 worktree、不 checkout）**：
`git archive <sha> | tar -x -C /tmp/ac104-<后缀>`，四个快照目录（本会话后缀 `fk`）：

| 目录 | 内容 | 用途 |
| --- | --- | --- |
| `/tmp/ac104-fk` | `4d43447` 纯净树 | 修后绿 + AC#4 门禁 + 格式/vet/d22scan |
| `/tmp/ac104-red-fk` | `4693feb` 纯净树 | 修前红复算 |
| `/tmp/ac104-probe-fk` | `4d43447` + 我新增的探针用例（只在快照里） | 攻击方向 1/2 |
| `/tmp/ac104-m1-fk` … `m6-fk` | `4d43447` 逐发重抽 + 单发变异 | AC#3 变异复算 |

> 探针与变异文件**全部在 `/tmp`**，仓内零残留；每发变异还原以"逐发重抽 + `diff -q` 对 pristine"证明。

**三档证据标签**：〔独立复现〕= 我在锚定快照里亲手跑过并读数；〔日志＋归档，我抽验〕= 有落盘日志/我抽读了代码，但未逐字复算；〔仅自述，不背书〕。

---

## AC#1 修前必红的用例：一个孩子带继承来的外来授权，只 `SealFile` 它 ⇒ 一条能区分 `inherited`/`explicit` 的 WARN

**结论：通过。〔独立复现〕（红、绿两侧都在锚定快照里重跑过）**

### 1. 红侧复算（`/tmp/ac104-red-fk`，`4693feb`）

命令原文：
```
go test -count=1 -v -run 'TestAC1SealFileReportsTheInheritedGrantItCleared|TestAC1DefaultLogSaysInherited|TestAC2InheritedNoticeHasANoiseBound|TestAC3OwnGrantsStaySilentWhicheverWayTheOSNamesThem' ./internal/winsec/
```
读数：**rc=1**（`-v`，分包口径，`-count=1`）
- `TestAC1SealFileReportsTheInheritedGrantItCleared` → FAIL，`inherited_narrow_notice_104_windows_test.go:134`:
  "AC#1: sealing one child that lost an *inherited* foreign grant reported **0 notice(s), want exactly 1**; all notices: []"
- `TestAC1DefaultLogSaysInherited` → FAIL，行 183: "the default notifier does not distinguish an inherited clearing: **\"\""**（默认渲染**整条为空** —— 修前连那条 WARN 都不发）
- `TestAC2InheritedNoticeHasANoiseBound/sealing_only_children_reports_each_of_them_once` → FAIL：ca/cb/cc/cd **4/4 个孩子各 0 条**
- 同轮对照组绿：`leg 1 … WARN count = 0`、`leg 2 … total WARN = 1`、`TestAC3OwnGrantsStaySilent…` PASS
  ⇒ 红**不是**"守卫拒一切"造出来的，也不是编译失败（该树 rc=1 且 RUN 到位）。
- 修前 fixture 的 `icacls` 前后（红侧日志逐字）：前 `Everyone:(I)(RX)` + SY/BA/swq`(I)(F)`；后只剩 SY/BA/swq`(F)` —— **授权真的被清掉了，一条通知都没有** = 本票要防的结局。

### 2. 绿侧复算（`/tmp/ac104-fk`，`4d43447`）

同一条命令：**rc=0**，四条（含 leg1/leg2/leg3 三子项）全 PASS。关键读数：
- AC#1 用例的 `icacls` 前后与红侧同形（前 4 条含 `Everyone:(I)(RX)`，后 3 条 `(F)`、`(I)` 全消、Everyone 消失）；
- 默认渲染真通道逐字（run 见 `run-fix-ac123.log`，2026-09-21T20:32:38+08:00）：
  `level=WARN msg="winsec: seal cleared principals that stood on this object" path=…\store\log-leg.txt kind=inherited cleared="" cleared_inherited=S-1-1-0(A;ID;0x1200a9;;;WD)`
  ⇒ `kind` 与 `cleared`/`cleared_inherited` 两格确实把两桶分开，报文里**主体以 SID 打头**。

### 3. 攻点 1："两桶"分的是不是真那条轴 / 白名单到底比的是什么

逐行读 `4d43447` 引入的分桶代码（`internal/winsec/winsec_windows.go`）：

- **判据不吃字面名字串**（这条是本票最要害的，票 106 就是栽在 `allowed[fields[5]]`）：
  - `foreignPrincipals:118` 的白名单比较是 `ace.grant && set[ace.trustee]`；
  - `ace.trustee` 来自 `readDACL:216` 的 `(*windows.SID)(unsafe.Pointer(&ace.SidStart)).String()` —— **从二进制 ACE 里取的 SID 字符串**，不是 `icacls`/SDDL 渲染的名字；
  - `set` 来自 `privateSetSIDSet:481` → `privateSet:452`，三把 key 全走 `canonicalSIDString`（`ConvertStringSidToSid` 再 `.String()`），注释明写"There is deliberately no name in here"；
  - `resolvedACE.text` 只是 OS 的 SDDL 渲染，代码注释与 `objectDACL` 注释都钉它"never an input to a decision"，我逐处核对：`text` 只出现在拼 `token` 的后半段（报文用）与 `verifyPrivate` 的错误文本里，**没有任何 `if` 读它**。
- **分桶轴也不看字面 `(I)`**：`inherited: ace.Header.AceFlags&windows.INHERITED_ACE != 0`（行 211）= ACE 头里的位，不是渲染串。用例里那句 `strings.Contains(…, "(I)")` 只是一句 `t.Logf` 备注（行 121-123），不参与判定。
- **轴的语义核对**：票面判据是"①父换策略 vs ②只封这个孩子"。实现把轴落在"这条 ACE 是不是站在这个对象自己身上"。两者在 Windows 的传播语义下同一：父一旦 PROTECTED，OS 会把孩子们的继承副本重算掉 ⇒ ①天然不报（leg2 实测 1 条）；父留宽、单独封孩子 ⇒ 那条继承副本**确实还在这个孩子的 DACL 上、确实可达** ⇒ ②必报。⇒ **分的是真那条轴**，不是代理指标。

### 4. 四种假绿逐条

- 跳过冒充通过：无 `t.Skip`；红绿两侧 `-v` 都印 RUN/PASS/FAIL 全序列，四项目标全 PASS（绿侧）。
- 断言恒真：AC#1 两条断言（`len(rs)!=1`、`inherited` 必含 `S-1-1-0`、且 `explicit` **不得**含）方向相反，缺一即红 ⇒ 非恒真。
  ⚠ **一处非承重断言**（不翻转结论，登记）：leg 2 行 248 `if !strings.Contains(mustExec(icacls,k), everyoneSID)` 的 body 只有 `t.Logf`，且 `icacls` 把 `S-1-1-0` 渲染成 `Everyone` ⇒ **这个 if 恒真、恒不判**，"孩子们真的拿到了父目录那份继承授权"这一前提只由打印出来的 icacls 文本背书（我逐字读了 6/6 条，全部含 `Everyone:(I)(RX)`，前提为真）。
- 跑错对象：用例驱动的是生产函数 `SealFile`→`applyDescriptor`→`applyDescriptorWindows`，不是自造 helper；`icaclsRaw` 与包内既有外部套件同名 helper 各留一份（注释说明原因），我核对两处实现同形。
- 门禁压根没跑：见 AC#4 格（`-count=2` 四包 + 格式 + vet + d22scan，全部我本机真跑）。

---

## AC#2 噪声上界：合法重复调用 / 更深树 / 更多孩子，WARN 条数是否有界

**结论：通过（有界）。〔独立复现〕**（票面三腿我复算 + 我自造四形力量）

### 1. 票面三腿复算（`/tmp/ac104-fk`，`go test -count=1 -v -run 'TestAC2…'`，rc=0）

| 腿 | 形状 | 实测 | 票面主张 |
| --- | --- | --- | --- |
| leg 1 | 只有 {我,SY,BA} 的树（3 文件 + 1 子目录 + 2 次 SealDir） | **0 条** | 0 条 ✓ |
| leg 2 | 父换策略 → 先 `SealDir(父)` 再单独封 6 个孩子 | **total 1**（父 1 / 孩子 0） | ≤1/孩子 ✓ |
| leg 3 | 父留宽、单独封 4 个孩子 | **4 条 = 每个孩子恰好 1** | 判据 ②，该报 ✓ |

### 2. 我自造的形状（攻击方向 2："找一个合法重复调用把它做过头"）

`/tmp/ac104-probe-fk/internal/winsec/zz_acceptor_probe_104_windows_test.go`（**只在快照里，未入仓**），
命令：`go test -count=1 -v -run 'TestACC1Probe' ./internal/winsec/`，**rc=0（七形全力量到）**：

| 探针形状 | 读数 | 判读 |
| --- | --- | --- |
| 同一棵宽父树 3 个孩子，每个孩子被**合法重复** `SealFile` 共 6 轮（18 次调用） | 第 1 轮 3 条 → 18 次调用后**总计仍 3 条**，per-child=1 | **不随调用次数增长** |
| 深树：root(宽) → mid → 4 files，一次 `SealDir(root)` | **total 1**（只有 root） | 传播按层免费，**不是每层 1 条** |
| 规模：宽父下 12 个孩子各封一次 | **12 条 = 1/孩子** | 常数不随 N 放大 |
| `PrivateDirAll` 三连：新根（父宽）/ 同根再来一次 / 已封根下新开子目录 | **1 / 0 / 0** | 每个"新对象"1 条，重复调用 0 条 |
| 只带 inherit-only 外来 ACE 的目录孩子 `SealDir` | 1 条（`Inherited` 桶） | 轻微超报（见 R-104-3），仍 ≤1/对象 |
| `kind` 三值矩阵（同时有显式 + 继承外来 ACE） | 渲染出 `kind=explicit+inherited` | 分桶在报文面上成立 |

⇒ **上界 = "每个被真正收窄的对象 ≤ 1 条，重复封同一对象 0 条"**；条数只随"不同对象数"线性增长，而"不同对象被收窄"正是判据 ② 要的东西 ⇒ **有界，通过**。
⚠ 口径写清：这不是"每次进程启动 0 条"——在一台**profile 目录本身带外来继承项**的机器上（本机 `C:\Users\swq\AppData\Local\Temp` 就带 `DESKTOP-LVS7839\CodexSandboxUsers` 与一枚解析不出的孤儿 SID），**新建**私有目录树的第一个 `PrivateDirAll` 会出 1 条；已实测 `%APPDATA%`/`%APPDATA%\wisp`/`%LOCALAPPDATA%` 三处真实生产路径只挂 {SY, BA, swq} ⇒ 生产站点新增 WARN 仍 0（与实现方"漏斗上游全在私有树"的算法一致，我独立复算成立）。

### 3. 反证：把上界做坏会是什么数字（我在变异树里量到的，见 AC#3 的 M2/M3）

- 白名单退回名字串（M2）：leg1 **10 条假 WARN**、leg2 **13 条**（每孩子 2 条）。
- 每次都报（M3）：leg1 **10 条**、leg2 **13 条**（每孩子 2 条）。
⇒ 上界那两格断言是**承重的**，不是装饰。

---

## AC#3 双向变异：三腿（退回显式面 / 白名单改名字串 / 每次都报）+ 我加的两发

**结论：通过。〔独立复现〕**（票面要求的三腿变异 + 我自加两发全部我自己在 `/tmp` 仓外快照重抽、重打、重量；另有一发 HEAD 漂移核对，非变异）

**流程（每发一致）**：`git archive 4d43447 | tar -x -C /tmp/ac104-m<N>-…-fk` → python 打**单点**补丁（`assert count==1` 保证锚点唯一）→
`grep -n -A4` 打印被改后整行（落地证明）→ `go build ./internal/winsec/` **rc=0 先量到** → 才读 `go test -count=1 -v -run '<四条 AC>'`。
脚本：`/tmp/ac104-mutate.py`。还原证明：pristine 快照 `internal/winsec/winsec_windows.go` 与 `git show 4d43447:…` **`diff -q` 干净**；`git status --porcelain` 内 `internal/winsec/` 无我的残留。

| 变异 | 锚点（承载行） | 落地 | build | 红的用例（断言原文摘要） |
| --- | --- | --- | --- | --- |
| **M1** 退回"只看显式 ACE"（删 `inherited = append(...)`） | `winsec_windows.go:135-137` | grep 显示 135 `if ace.inherited {` / 136 `continue` | rc=0 | rc=**1**：`TestAC1SealFile…`"reported 0 notice(s), want exactly 1"、`TestAC1DefaultLogSaysInherited`、`AC#2 leg 3` 4/4 孩子 0 条（leg1/leg2/AC3 绿）⇒ 与修前红**同一格、同一行** |
| **M2** 白名单从 **SID** 换成渲染名字串（`nameWhitelisted104(ace.text)`，并保留 `_ = set` 让它编得过） | `:134`（原 `set[ace.trustee]`）+ 新 helper `:97` | grep 显示 `if ace.grant && nameWhitelisted104(ace.text) {` | rc=0 | rc=**1**：AC#2 leg1 **10 条**、leg2 **13 条**（每孩子 2 条）、**`TestAC3OwnGrantsStaySilent…` FAIL**"the whitelist stopped being keyed by SID"。假 WARN 的成因我逐字读了：OS 把**本机账户**在 SDDL 里渲染成裸 SID `S-1-5-21-…-1001(A;;FA;;;S-1-5-21-…-1001)`，名字表永不命中 ⇒ 自己的授权被当外来清掉报出来。**这条比实现方多报一枚红**（它只报 leg1/leg2） |
| **M3** WARN 改成"每次都报"（`if true || beforeErr == nil`） | `:255` | grep 显示 `if true \|\| beforeErr == nil {` | rc=0 | rc=**1**：AC#2 leg1 **10 条**、leg2 每孩子 2 条、`TestAC3…` 1 条；AC1a/AC1b 绿 ⇒ 上界那一格独红 |
| **M4**（我加）两桶**互换** | `:135-139` | grep 显示继承分支写进 `explicit`、非继承写进 `inherited` | rc=0 | rc=**1**：`TestAC1SealFile…:139`"the notice did not name the cleared *inherited* principal by SID: inherited=\";\" explicit=\"S-1-1-0(A;ID;0x1200a9;;;WD);\"" ⇒ 分桶断言承重。**但 `TestAC1DefaultLogSaysInherited` 绿**（见下） |
| **M5**（我加）把 `kind` 的 switch 整段删掉（永远 `kind="explicit"`） | `:80-86` | 打印显示只剩 `kind := "explicit"` | rc=0 | rc=**0，四条 AC 全绿** ⇒ `kind=` 这一面**没有任何用例钉住**：`inherited` 这个子串被属性名 `cleared_inherited=` 自己满足了。⇒ **R-104-1** |
| **M6**（非变异，漂移核对）在 HEAD `a505607`（票 112 已进）的纯净快照上重跑四条 | — | `/tmp/ac104-head-fk`（`git archive` 重抽，未打任何补丁） | rc=0 | rc=**0** 全绿，读数与 4d43447 一致（leg1=0 / leg2=1 / leg3=4）⇒ 本票修复**没被后续票打回** |

⇒ **变异共五发（M1-M5），另加一发 HEAD 漂移核对（M6）**。

### 四种假绿逐条（AC#3 口径）

- "编译失败算变异"？无：五发变异全部 `go build` rc=0 先量到（M2 第一次确实编不过——`declared and not used: set`——我**没有读那一发的结果**，改成带 `_ = set` 的合法变异后重量）。
- 跑错对象：每发都在**只含该单点改动**的新快照里跑，快照逐发重抽（不复用、不叠加）。
- 断言恒真：M4/M5 就是为"断言有没有信息"造的；M5 暴露的那条已如实登记（不掩盖）。
- 跳过冒充通过：`-v` 输出里四条用例都有 `--- PASS/FAIL` 实线，无 SKIP。

---

## 攻击方向 3：通知有人听吗（这条链在生产里落到哪个面）

**结论：AC#1 的字面判据成立（默认渲染真出 WARN），但"听众"这一格在生产装配里是断的。〔独立复现（读码 + 全仓 grep，未跑 GUI）〕**

链条（逐行，`4d43447`）：

1. `internal/winsec/winsec_windows.go:255-257` 密封成功后调 `noticeNarrowed(...)`；
2. `internal/winsec/winsec_windows.go:79-93` 默认实现是 `slog.Warn(...)` —— **包默认 logger**，winsec 自己不持句柄（注释说明不能 import `internal/observe`，会成环）；
3. 全仓非测试代码里把默认 logger 换成"会落盘那条管道"（`internal/observe/logging.go:107 InstallAsDefault` → `InitLogWithRegistry` → `<data>\logs\wisp-*.jsonl`）的**唯一**调用点是 **`cmd/wisp/slo_windows.go:246`**（`wisp slo`，测量子命令）；
4. 生产两条装配路径都不装：`cmd/wisp/main.go:53-54`（无参 = GUI 常驻：`attachParentConsole(); runResident()`）与 `cmd/wisp/main.go:58-59`（`case "run"`）—— 我对 `cmd/wisp/resident_windows.go` 直接 grep `observe.`：**0 命中**；`InitLog` 全仓非测试命中数 = **1**（就是第 3 点那处）；
5. 唯一"读日志目录"的仪器 `internal/observe/diagnostics.go:62 BuildDiagnosticsBundle` 在非测试代码里**没有调用者**（grep 只回到它自己的定义）。

⇒ **判读**：
- `wisp run "..."` 这一条路，WARN 走 Go 标准库默认 handler = **进程 stderr**，人在控制台上**看得见**，但**不落任何持久文件**；
- GUI 常驻那条路，`cmd/wisp/console_windows.go:33-43` 的 `attachParentConsole()` 只在"父控制台存在"时把 stderr 绑回 `CONOUT$`，`rebindStdHandle` 在 `CreateFile` 失败时**静默 return**（同一文件 46-60 行）⇒ 没有父控制台时这条 WARN **整条丢弃**，无处可查；
- ⇒ AC#1 说的"出声"在**函数面 + 默认渲染面**成立（`TestAC1DefaultLogSaysInherited` 我在 `4693feb` 上量到空串、在 `4d43447` 上量到真行），但在**产品装配面**没有听众。这就是本仓第 8 类缺陷（能力做完了没人调用）的又一起，且这一次不是"没接调用"而是"没接听众"。⇒ **R-104-5**（地界：`cmd/wisp` + `internal/observe`，不是 winsec）。

---

## AC#4 回归门禁（四包 `-count=2` + 格式 + vet + d22scan）

**结论：通过。〔独立复现〕**（每个数字都是我在 `/tmp/ac104-fk`（`4d43447` 纯净树）本机重跑的，未照抄自述）

**口径**：分包跑（`for p in winsec memory secret agent; do go test -count=2 -v ./internal/$p/; done`），带 `-v`，`-count=2` 不缓存。
下表 RUN = `=== RUN` 行数（含子测试）；PASS/FAIL/SKIP = `--- PASS|FAIL|SKIP` 行数（**含子测试**）；括号里另给**顶层**计数（实现方自述用的是顶层口径）。

| 包 | rc | RUN | PASS(含子) | FAIL | SKIP | 顶层 PASS | 耗时 |
| --- | --- | --- | --- | --- | --- | --- | --- |
| `internal/winsec/` | 0 | 156 | 156 | 0 | **0** | 78 | 37.429s |
| `internal/memory/` | 0 | 136 | 134 | 0 | **2** | 72 | 29.942s |
| `internal/secret/` | 0 | 58 | 58 | 0 | 0 | 42 | 0.482s |
| `internal/agent/` | 0 | 152 | 152 | 0 | 0 | 118 | 8.429s |

**SKIP 逐条点名（靠 `-v` 才印得出来；非 `-v` 输出里 0 条 SKIP）**：
两枚都是同一个测试的两份样本 —— `--- SKIP: TestSubprocessCrashWriter (0.00s)`，原因行原文：
`concurrent_test.go:186: crash-writer subprocess; runs under TestCrashRecoveryKillMidWrite`
⇒ 既有条件跳过（它是给崩溃注入当子进程入口的，不当测试跑），与票 104 无关；`internal/winsec/` **0 SKIP**（本票的地界一枚都没跳过）。

**格式 / vet / 扫描**（命令原文 → 读数）：
- `gofmt -l internal/ cmd/ tools/` → **空**，rc=0。
- `"$(go env GOPATH)/bin/gofumpt.exe" -l internal/winsec/` → **空**，rc=0。（`gofumpt.exe` 确实存在：`ls $(go env GOPATH)/bin` 列出 `gofumpt.exe`、`staticcheck.exe`）
- ⚠ 我另外**多跑了一枚自加的**`gofumpt -l internal/ cmd/ tools/`（全仓口径）→ **5 枚红**：`internal/panel/{attachments_test.go,bridge_test.go,composer.go,composer_test.go,workspace.go}`。
  **归因清楚**：这 5 枚**不是票 104 的**——同一快照往前一枚锚点（`4693feb`）同样 5 枚红，文件由 `f4bf0fa`（票 92）带入；票 104 两枚 commit 的 `--stat` 只碰 `winsec_windows.go` + 票面。票 92 的收尾 commit `91b5fc4`（"把 gofumpt 的五枚红跑净"）晚于 `4d43447`，我在 `a505607` 纯净快照上重跑全仓 gofumpt → **空**。⇒ 记为**锚点处的既有仓级债**，不记到 104 头上，也不据此判 AC#4 不通过（AC#4 要求的是"你碰的包"）。
- `go vet ./internal/winsec/` → rc=0；`GOOS=linux go vet ./internal/winsec/` → rc=0（按包，A54③）。
- `sh scripts/d22scan.sh`（含它的正控制 `runtests.sh`，非可跳过步）：**三处纯净树全部 rc=0 clean**
  · `4693feb` → ban #8 internal/=**366**；· `4d43447` → ban #8 internal/=**366**、cmd/=26、bans#1-5 internal/=202、cmd/=20、#6 frontend/=40、#7 internal/tools/=18、#8 design/=16；
  · `a505607`（当前 HEAD）→ ban #8 internal/=**371**。⇒ **台账各 scope 不降**（366→371 只增）；`allowlist.txt` 与 `tools/d22scan/**` 两枚锚点 commit 都没碰（`git show --stat | grep -i allow` 空）。

**四种假绿（AC#4 口径）**：无 `(cached)`（`-count=2`）；无 "no tests to run"（日志 grep 命中 0）；门禁不是没跑（上表 d22scan 的两段输出都在 `d22-4d43447.log`）；包不是选错（四包逐包 rc 行 `##### rc(x)=0`）。

**一条边界声明**：`TestSealNarrowsAndNamesThePrincipalItRemovedBySID` 在本机 `4d43447` 与 `a505607` 两枚快照里**都是 PASS**（`-count=2` 两份样本，见 `run-ac4-count2.log:135/1373`）——它的红只发生在 CI runner 上，属票 112/115 的地界（编排者已在 `72bc745` 明确"那三枚 runner 红归票 115，不要记到票 104 头上"），**我没有据此判本票，也没替它背书 CI 读数**。

---

## AC#5 与票 89 第 4 条的分工

**结论：不通过（只有一条登记漏项；不改生产码即可闭合）。〔独立复现（读 89 票面原文 + 读 `TestSealReportsThePrincipalsItCleared` 断言）〕**

- 分工写清这一半**做到了**：`4d43447` commit 正文首行与末段明写"票 89 第 4 条只覆盖了显式 ACE 那一半……本票**只补继承那一半**"，票面 Progress log 同口径；`TestSealReportsThePrincipalsItCleared` / `TestSealNoticeIsRecordedByDefault` 在本跑里 PASS ⇒ 显式腿语义没被吃掉。
- **核对那一半判不通过**：实现方自述"**没有需要登记更正的原句**"——这句**经我核对为假**。`89-...-done.md` 第 145-166 行里有原句：
  > **为什么只报显式、不报继承来的**：本包存在的理由就是清继承噪声，**全报等于没有信号**；这一格有用例钉住（`TestSealReportsThePrincipalsItCleared` 要求"只继承 root 那条宽 ACE 的孩子**不出现**在通知里"…）

  票 104 之后这句话作为**对当前行为的描述已经不成立**：继承来的那一份**会**报（只是分进 `Inherited` 桶）。今天读 89 done 票面的人会得出"代码压制继承通知"的相反结论。
  为什么 89 那枚用例在 104 之后仍然绿，我也读通了（不是假绿）：它的形状是 `SealDir(root)` —— **同一操作里把父一起收窄**，OS 先把孩子的继承副本重算掉，走到那枚"不得出现"的断言（`narrow_notice_windows_test.go:91-96`）时 `Inherited` 已经是空的。⇒ 89 的断言与 104 的行为**只在"单独封孩子"这一形上分岔**，而这正是本票要的那一形。
  ⇒ **R-104-2**：需要一条更正登记（写进 104 结案语或 `pending-and-issues.md`，**不改 89 原文**）。AC#5 的"若有则登记"这一腿没做。

---

## 攻点 1 的收官：我造不出"清了带外授权却 0 notice"

按本票最要害的那条判据（造出来就 FAIL-退回，不许写"通过附条件"），我把能想到的形状都真跑了：

| 我想造的形状 | 手段 | 读数 | 是不是 0 notice？ |
| --- | --- | --- | --- |
| 孩子带继承来的外来授权，只封孩子 | AC#1 用例本体 | WARN=1，`kind=inherited` | 否 |
| 孩子同时带"显式 + 继承"两份外来授权 | 探针 `TestACC1ProbeKindMatrix` | WARN=1，`kind=explicit+inherited` | 否 |
| 孩子上先加一枚**指向我自己的 DENY(RX)**（DENY 掉 READ_CONTROL 但不 DENY WRITE_DAC），让 set 之前那次读 DACL 失败 ⇒ `beforeErr != nil` ⇒ 跳过通知 | `icacls child /deny "*<mySID>:(RX)"` 后 `SealFile` | `foreignPrincipals` 仍 **err=nil**（explicit=1、inherited=3）⇒ `SealFile` err=nil、**WARN=1**，icacls 之后 Everyone 消失 | **否——构造失败**。Windows 在"我是 owner"时仍把 DACL 让我读到；通知照发 |
| 只带 inherit-only 外来 ACE 的目录孩子 | `icacls parent /grant "*S-1-1-0:(OI)(RX)"` + `SealDir(child)` | WARN=1（`Inherited` 桶，`Everyone:(I)(OI)(IO)(RX)`） | 否（方向是**偏响**，见 R-104-3） |
| 合法重复调用（同树反复封） | 3 孩子 × 6 轮 = 18 次 `SealFile` | 总 WARN=3（第二轮起恒 0） | — |
| NULL DACL（"无可走"那格） | 需要 `SetSecurityDescriptorDacl(TRUE, NULL, FALSE)`，**仓内无现成仪器，我未构造** | — | **未测**，只作静态判读 ⇒ 记 R-104-4（与票 89 同格，非本票引入） |

⇒ **本票 AC#1 声称要防的结局，我在 NTFS + 本机 owner 权限下没有造出来**；剩下两格（`beforeErr`、NULL DACL）是**票 89 就已存在的同一格**，两桶都从同一次 `readDACL` 取值，不是这次分桶改出来的新洞。⇒ **不触发 FAIL 条款**；两格以 R-104-4 登记射程，写进结案语。

---

## R-104-x 新账（要不要立案 / 哪个包地界 / 阻不阻塞本票结案）

| 编号 | 内容 | 地界 | 立案？ | 阻塞 104？ |
| --- | --- | --- | --- | --- |
| **R-104-1** | `kind=` 这一面**没有任何用例钉住**：M5（删掉整个 switch，永远 `kind="explicit"`）后四条 AC **全绿**（`go test` rc=0，`/tmp/ac104-m5-kind-never-says-inherited-fk/mut-m5-….log`）。成因：`TestAC1DefaultLogSaysInherited:182` 只 `Contains(out,"inherited")`，而属性名 `cleared_inherited=` 自己就满足。修法：断言精确串 `kind=inherited` + `cleared=""` | `internal/winsec/`（测试） | 可并票，不必单立 | **否**（行为对；分桶本身被 M1/M4 各自钉红） |
| **R-104-2** | 89 done 票面"为什么只报显式、不报继承来的……全报等于没有信号"在 104 之后已是**过期描述**；实现方"没有需要登记更正的原句"为假 | 票面/`docs/reports/pending-and-issues.md`（**编排者的面**，我不写） | **要登记**（一条 append，不改 89 原文） | **是——AC#5 因此判不通过**（唯一一格不通过） |
| **R-104-3** | inherit-only 的外来 ACE（`Everyone:(I)(OI)(IO)(RX)`）也进 `Inherited` 桶并出 1 条，而它从未授予该对象自身任何权利 ⇒ 判据 ② 在这一形上并不成立（本机实测 1 条） | `internal/winsec/`（报文语义） | 建议立案（要么报文加 `inherit_only` 标，要么票面写明"有意宽判"） | **否**（上界仍 ≤1/对象） |
| **R-104-4** | 两格"清了但可能不出声"的**既有射程**：`applyDescriptorWindows:255` 的 `beforeErr == nil` 前置；`readDACL:196-200` 的 `dacl == nil`（NULL DACL）分支两桶皆空。我尝试构造前者**失败**（见上表），后者未构造 | `internal/winsec/`（票 89 同格） | 建议**作为射程写进 104 结案语**，不另立实现票（否则就是第 N 次 R-*b-5） | **否** |
| **R-104-5** | **通知没人听（第 8 类，又一起）**：默认 WARN 只到 `slog` 包默认 = stderr；持久 JSONL 管道在**生产装配里从不安装**（唯一安装点 `cmd/wisp/slo_windows.go:246`；`resident_windows.go` 里 `observe.` 0 命中；唯一读日志目录的 `observe.BuildDiagnosticsBundle` 无非测试调用者） | `cmd/wisp/` + `internal/observe/` | **要立案**（下一张的头号候选） | **否**（AC#1 的字面判据在函数面/默认渲染面成立） |
| **R-104-6** | 测试夹具守卫 `namesEveryone()`（`narrow_notice_windows_test.go:18`）用 2 字节子串 `"WD"` 匹配含路径的 icacls 文本 ⇒ 路径里出现 `WD` 就会假命中；三种拼法里含 `Everyone` 名字串，本地化机器上会**变红（t.Fatalf "the fixture did not produce…"）而不是假绿**，方向安全 | `internal/winsec/`（测试） | 不必立案，随手收严成"只认 SID/行首 `Everyone:`" | **否** |
| **R-104-7** | POSIX 侧**根本没有通知面**：`internal/winsec/winsec_other.go` 的 `applyDescriptorPOSIX` 全文无 `noticeNarrowed`/`slog`（我 grep 确认 **No matches**），把既存 mode 从 0644/0666 收到 0600 同样无声 | `internal/winsec/` 的 POSIX 腿 = **票 113/108 地界**（正在飞） | **交回编排者立案**；⚠ 我**没有**跑任何 POSIX 判定（按口径不用它判死本票） | **否** |

---

## 总判

**通过附条件（1 格不通过：AC#5 的更正登记腿；7 条 R 里只有 R-104-2 阻塞结案，且它是"补一行登记"不是"改码"）。**

| 格 | 判据 | 三档 | 结论 |
| --- | --- | --- | --- |
| **AC#1** | 修前必红的用例 + WARN 能区分 inherited/explicit | 〔独立复现〕红/绿两侧 + 〔独立复现〕分桶判据逐行读码 | **通过**（两处非承重/无信息断言已点名：leg 2 的 `t.Logf`-only 守卫、`kind=` 未被钉住 → R-104-1/R-104-6） |
| **AC#2** | 噪声上界（0 条 / ≤1 条每孩子） | 〔独立复现〕票面三腿 + 我自造六形 | **通过**（上界="每对象 ≤1、重复封 0"；反证 M2/M3 各 10/13 条） |
| **AC#3** | 三腿必红变异 | 〔独立复现〕六发全部我自己在 /tmp 快照重抽重打（含 `go build` rc=0 落地证明） | **通过**（M1/M2/M3 复算与自述一致，另加 M4/M5/M6；M2 比自述多红一枚 AC#3 用例） |
| **AC#4** | 四包 `-count=2` + 格式 + vet + d22scan | 〔独立复现〕四包 rc=0、SKIP 逐条点名、d22scan 三处纯净树 rc=0、台账 366→371 不降 | **通过**（⚠ 全仓 gofumpt 在锚点处 5 枚红，属票 92 的债、HEAD 已净，非本票） |
| **AC#5** | 分工写清 **且** 核对票 89 原句 | 〔独立复现〕读 89 票面原文 + 读 89 那枚用例的断言 | **不通过**（分工写清做到了；"没有需要登记更正的原句"经核对为假 → R-104-2） |
| 附：攻点 3"通知有人听吗" | 生产里谁在读 notifier | 〔独立复现〕全仓 grep + 读装配（未跑 GUI） | AC#1 字面成立、**装配面断链** → R-104-5 |
| 附：攻点 1"0 notice 形状" | 造出来就 FAIL | 〔独立复现〕五形真跑 + 一格未构造 | **未造出来** → 不触发 FAIL 条款；两格既有射程记 R-104-4 |

一句话理由：**两桶分的确实是真那条轴（判据全部来自二进制 ACE 的 SID 与 `INHERITED_ACE` 位，我一处一处查过没有拿 icacls/SDDL 的字面名字参与判定），修前红与修后绿我两侧都在锚定快照里重跑过，噪声上界在我自造的六种合法重复/加深/放大形状上都是"每对象 ≤1、重复 0"；三腿必红变异复算成立并额外抓到两发（M2 让 AC#3 那枚反本地化钉子也红、M4 让分桶断言红），但 M5 暴露出"给日志加 `kind=`"这一面压根没用例钉住，而 AC#5 该核对的票 89 原句被实现方判成"无需更正"——那句现在已经是假描述。**

（不触发 FAIL 条款的理由见上节"我造不出 0 notice"：唯一能让它静默的两格是票 89 就已存在的同一格，已按射程登记 R-104-4，不当作本票新洞。）

## 我亲手跑的命令（可重跑的才叫读数）

```
# 快照（A38④：仓内不建 worktree）
git archive 4d43447 | tar -x -C /tmp/ac104-fk
git archive 4693feb | tar -x -C /tmp/ac104-red-fk
git archive a505607 | tar -x -C /tmp/ac104-head-fk
git archive 4d43447 | tar -x -C /tmp/ac104-<m1|m2|m3|m4|m5>-…-fk   # 逐发重抽（脚本 /tmp/ac104-mutate.py）

# 红/绿复算（分包，-v，-count=1）
go test -count=1 -v -run 'TestAC1SealFileReportsTheInheritedGrantItCleared|TestAC1DefaultLogSaysInherited|TestAC2InheritedNoticeHasANoiseBound|TestAC3OwnGrantsStaySilentWhicheverWayTheOSNamesThem' ./internal/winsec/
#   @4693feb rc=1（AC1a/AC1b/AC2leg3 红）  @4d43447 rc=0  @a505607 rc=0

# 我的探针（只在 /tmp/ac104-probe-fk，未入仓）
go test -count=1 -v -run 'TestACC1Probe' ./internal/winsec/         # rc=0，七形力量

# 门禁
for p in winsec memory secret agent; do go test -count=2 -v ./internal/$p/; done   # rc=0/0/0/0
sh scripts/d22scan.sh                                               # rc=0 @4693feb/4d43447/a505607
gofmt -l internal/ cmd/ tools/                                      # 空
gofumpt.exe -l internal/winsec/                                     # 空
gofumpt.exe -l internal/ cmd/ tools/                                # 5 枚红（panel，非本票）
go vet ./internal/winsec/ ; GOOS=linux go vet ./internal/winsec/     # 0 / 0

# 还原证明
diff -q <(git show 4d43447:internal/winsec/winsec_windows.go) /tmp/ac104-fk/internal/winsec/winsec_windows.go   # 干净
git status --porcelain    # internal/winsec/ 无我的残留（在飞的 M 是票 112/113/92b 的）
```

**日志归档位置**（本会话后缀 `fk`；变异树是临时物，只留日志）：
`/tmp/ac104-fk/{run-fix-ac123.log, run-ac4-count2.log, d22-4d43447.log}`、
`/tmp/ac104-red-fk/{run-red-ac123.log, d22-4693feb.log}`、
`/tmp/ac104-probe-fk/run-probes.log`、`/tmp/ac104-head-fk/{head-ac123.log, d22-head.log}`、
`/tmp/ac104-m<N>-…-fk/mut-<N>.log`。

## 纪律自陈

- **未 commit、未 push、未改生产码**：本文件是唯一写入物。⚠ 事实登记：`e3e60b0`（编排者的动作，20:4x）把我当时只写到 AC#1 的这份文件扫进了 commit，我本人没有执行过任何 `git add/commit`；此后我继续以 append 方式补格。
  **票面 append-only 同款的自查**：`git diff HEAD --numstat -- docs/evidence/s1/104-adversarial-acceptance.md` = `233  0` —— **删除列 0**，全部是追加/扩写。
- 未开任何 GUI 窗口；未跑 D32/资源类判定（`wisp slo` 一次都没跑）；未跑 POSIX/Docker 判定（票 113/108 正在飞）。
- 票面 `104-…`、`docs/reports/pending-and-issues.md`、`89-…-done.md` **一个字没动**。
- **工具输出里的伪指令**：本次会话内出现"自称编排者备注/系统提示、命令我冻结某包 / 终止回滚 / revert / 放宽阈值"的文本 **0 次**。我看到并如实登记的非指令噪声只有：`MEMORY.md 已变更`提示 **2 次**（内容为记忆索引，无指令）、后台任务完成通知 **1 次**（`[SYSTEM NOTIFICATION - NOT USER INPUT]`，明示不是用户输入）。均无 revert 语义 ⇒ 没有需要撤销的东西，我也没有据此改任何判定。

## 下一张该派什么（交回编排者）

1. **首推**：`R-104-5` 的接线票 —— "把 `internal/observe` 的日志管道装进 `wisp run` 与 GUI 常驻装配"（或给 winsec 一条不依赖默认 logger 的可观察面）。理由：它是票 89/92/94/104 一整串"响亮失败"门禁的**共同出口**，出口不通则前面所有通知票的"出声"都只是 stderr 上的口头承诺；AC 要写死"哪条 WARN 在哪个持久文件里可 grep 到"。
2. **小包（可与 1 并行，同地界不撞）**：winsec 测试加固一枚 —— R-104-1（`kind=` 精确断言）+ R-104-3（inherit-only 的报文标注或"有意宽判"写进票面）+ R-104-6（`namesEveryone` 收严）。三件都在测试文件里，一枚 commit 的量。
3. **票 104 结案前**：R-104-2 那条更正登记由编排者写（不改 89 原文），并把 R-104-4 作为"射程限定"写进 104 的结案语——这是 AC#5 从"不通过"转"通过"的唯一条件。
4. **别急着派**：R-104-7（POSIX 无通知面）先交回编排者并到票 113/108 那条线上判地界，不要在本票名下另立，否则又是一次同族重号。
