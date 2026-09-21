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
