# 票 89 独立对抗验收（acceptor-ticket89）

**被验收对象**：`.scratch/wisp/issues/89-0600-is-decorative-on-windows-acl-for-private-data.md`（`Status: ready-for-review`，六框自勾）
**代码链**：`b994a2c`(骨架) → `3054007`(基线 icacls) → `de15a6b`(真 SD) → `57bdbb2`(接线 + AC#4/AC#5) → `0a3a445`(AC#6 门禁)
**验收环境**：本机真机，`go version go1.27.1 windows/amd64`，Windows 11，当前账户 `DESKTOP-LVS7839\swq`，
icacls = `C:\Windows\System32\icacls.exe`。本文件由**独立对抗验收代理**撰写，实现者的任何数字都不作为本代理的证据。

**快照纪律（全部在仓外 `/tmp`，仓内无 worktree、工作树未被本代理改动）**：

| 快照目录 | SHA | 用途 |
| --- | --- | --- |
| `/tmp/wisp89acc-89` | `0a3a445`（票 89 终态，`internal/winsec` 到 HEAD 无漂移，已用 `git diff --stat 0a3a445 HEAD -- internal/winsec` 核对为空） | 判据重跑 + 各项变异 |
| `/tmp/wisp89acc89-pre` | `3054007`（实现仍是 `os.Chmod` 占位） | 第 3 项"差点假绿"反向复现 |
| `/tmp/wisp89acc89-base` | `942ab5a` = `b994a2c^`（票 89 之前的生产码） | 第 1 项基线（本代理自造的仪器） |

---

## 进度总览（先落盘，边做边追加）

| 项 | 内容 | 状态 |
| --- | --- | --- |
| 2 | 纯净快照 `go test -count=2 -v ./internal/winsec/` 四数 | **已完成 → PASS** |
| 1 | 基线为假（本代理自造 icacls 仪器，票前生产码） | 进行中 |
| 3 | "差点假绿"独立证实（解析器退回跳过首行） | 未开始 |
| 4 | `PROTECTED_DACL` 变异 | 未开始 |
| 5 | AC#5 fail-closed + `sealError()` 归一化拿掉 | 未开始 |
| 6 | 覆盖面主张（`SealDir` 传播 + 出生那一刻继承） | 未开始 |
| 7 | `os.Remove` 删 junction 与票 79 A51② 的冲突 | 未开始 |
| 8 | `verifyPrivate` 白名单口径裁决 | 未开始 |

---

## 第 2 项：纯净快照判据重跑 —— **PASS（票面数字被独立复现）**

仓外纯净快照 `git archive 0a3a445 | tar -x -C /tmp/wisp89acc-89`，本代理自己的读数：

```
$ go test -count=2 -v ./internal/winsec/        # 在 /tmp/wisp89acc-89
ok  github.com/CarlosShao/wisp/internal/winsec  14.791s
RC=0

=== RUN   : 34        (全量 grep '^=== RUN' 计数)
--- PASS  : 20        (顶层；另有 14 条缩进的子测试 PASS 行，合计 34)
--- FAIL  : 0
--- SKIP  : 0
不同名测试数 : 17      (10 顶层 + 7 子测试名)
34 == 2 × 17  ✔  （每个名字恰好出现两次，无一例外，见下）
```

`=== RUN` 逐名计数（`sort | uniq -c`，全部恰好 `2`）：
`TestAC1BaselineProductionPaths` / `TestAC1BaselineForeignACEPropagatesIntoModeOnlyWrites` /
`TestAC2SealedDirCoversFilesItNeverTouched` / `TestAC3SealedWritesCarryNoForeignSID`（+ `artifact-exclusive`、
`secret-blob`、`staging-temp` 三子）/ `TestAC3ProductionDataRootIsPrivateEndToEnd` /
`TestAC4JunctionAtArtifactPositionIsNotRecursed` / `TestAC4UnlinkableLinkGivesNamedError` /
`TestAC4SealedWalkSkipsLinks` / `TestAC5FailedSealRefusesTheWrite`（+ `exclusive_artifact`、`replacement_write`、
`directory_chain`、`seal_file_and_dir_direct` 四子）/ `TestAC5FailureIsNotSwallowedByTheHappyPath`。
⇒ 票面"0 SKIP、17 条具名 PASS、无重名跑两遍冒充条数"**独立复现**。
结构上还有一条支撑：`grep -rn "t.Skip\|SkipNow\|testing.Short" internal/winsec/` = **0 命中**，本包没有任何 skip 通道。

**四种假绿逐条点名（本项）**：
1. `--- SKIP`：0 条，且包内无 `t.Skip` 调用 ⇒ 不是"skip 当 ok"。
2. `-run` 过滤：本项**没用** `-run`，全量输出计数。
3. 步骤被静默跳过：`assertPrivateACL`/`requireForeignModify` 走 `run(t,…)`，`icacls` 或 `powershell` 非零退出即 `t.Fatalf`；
   `TestAC1BaselineProductionPaths` 里 `wisp.db-wal` 缺失只会 `t.Logf` 不会 skip（本代理在第 1 项自造仪器里加了
   "一条都没测到就判失败"的 `if probed == 0 { t.Fatal }` 钉子，见下）。
4. 计数不达标 / 重名：34 = 2×17 且逐名 uniq -c 全为 2 ⇒ 票 79 那种"同一测试名跑两遍"未重现。

**残留观察（不算 FAIL，记给编排者）**：`TestAC1BaselineProductionPaths` 这个名字里的"基线"在 `0a3a445` 上
已经**不再是纯基线** —— `secret.NewStore`/`memory.Open` 在 `57bdbb2` 接了 winsec，所以那两类的 icacls 原文
测的是**修好之后**的状态（票面 AC#1 那段的原文是 `3054007` 时代录的，仍然成立）。⇒ 本代理不采信它的基线读数，
第 1 项改用票前树 `942ab5a` 上自己写的仪器重测。

---

## 第 1 项：基线为假 —— **PASS（本代理自造仪器，原文读数）**

**仪器**：`/tmp/wisp89acc89-base`（= `git archive b994a2c^` = `942ab5a`，票 89 之前的生产码）里
**本代理自己写的** `internal/acceptor89/baseline_windows_test.go`（不 import `internal/winsec`，不用它的解析器/断言）。
生产入口用票前代码：`os.MkdirAll(0o755)` + `os.OpenFile(O_CREATE|O_EXCL,0o600)`（`agent/spill.go:106/:245` 原形）、
`secret.NewStore`+`Store`（票前 `store.go:44/:74`）、`memory.Open`（票前 `open.go:172`，DSN 里 `journal_mode(WAL)`）。
**全部落在 `t.TempDir()` 下，未碰用户真实数据目录。**

`go test -count=1 -v ./internal/acceptor89/` → **RC=0，5 条 PASS，0 FAIL，0 SKIP**。判据是
"必须存在一个既不是我、也不是 SYSTEM/Administrators 的主体，且它手里的 ACE 能读"——**测不到就 `t.Errorf`**，
所以绿=基线确实为假。逐类原文（本代理的解析行，非票面）：

```
artifact.txt         foreign [DESKTOP-LVS7839\CodexSandboxUsers (I)(M,DC)
                              S-1-5-21-3623186960-731165060-4091685855-1717338598 (I)(M,DC)]
   all aces += NT AUTHORITY\SYSTEM (I)(F) / BUILTIN\Administrators (I)(F) / DESKTOP-LVS7839\swq (I)(F)
acceptor(DPAPI blob) foreign [同一两条 (I)(M,DC)]
wisp.db              foreign [同一两条 (I)(M,DC)]
wisp.db-wal          foreign [同一两条 (I)(M,DC)]
wisp.db-shm          foreign [同一两条 (I)(M,DC)]
model.bin.part       foreign [同一两条 (I)(M,DC)]   ← staging（downloader 的 0o644 同形 idiom）
```

⇒ **四类（artifact / DPAPI blob / db+wal+shm / staging）在票前生产码上全部被外来主体持
`(M,DC)` = MODIFY+DELETE_CHILD，且全部 `(I)` 继承而来**。基线主张成立。
同一次运行里 `os.Stat` 给的 `FileInfo.Mode()` 原文 = **`-rw-rw-rw-`** ⇒ "A51① 那套假位"也被独立复现。

**"这两个外来主体到底是不是真账户"（替代判据的关键一环）——本代理独立测得**：

```
net localgroup CodexSandboxUsers -> CodexSandboxOffline, CodexSandboxOnline
net user -> Administrator,CarlosShao,CodexSandboxOffline,CodexSandboxOnline,DefaultAccount,Guest,shaowq,swq,WDAGUtilityAccount,WsiAccount
CodexSandboxUsers 的 SID = S-1-5-21-1228170099-895614386-1166154857-1005（与我 -1001 同机器前缀 ⇒ 本机局部组）
```

⇒ 持有私有数据 MODIFY 的是**两个真实存在、与 `swq` 不同的本机账户**，不是空组名；这台机上还有
`CarlosShao`、`shaowq`、`Administrator` 等别的本地账户 ⇒ "同机其他账户"不是假想敌。
**替代判据在"有没有第二主体"这一环上不需要口令就能站住。**

---

## 第 3 项：它自曝的"差点假绿" —— **PASS（独立复现，反向钉子有效）**

在**未修复的实现**上做（`/tmp/wisp89acc89-pre` = `3054007`，`winsec_windows.go` 通体只有 `os.Chmod`）。

**(a) 先证"真红"**（解析器原样）：`-count=1 -v -run 'AC3|AC2|AC1BaselineForeign'` → **RC=1**：
`--- FAIL: TestAC3SealedWritesCarryNoForeignSID` + 三子 `{artifact-exclusive,secret-blob,staging-temp}` 全红、
`--- FAIL: TestAC2SealedDirCoversFilesItNeverTouched`、`--- PASS: TestAC1BaselineForeignACEPropagatesIntoModeOnlyWrites`。
红因原文：`NOT PRIVATE artifact-exclusive: foreign SID(s) S-1-1-0 (principals: [Everyone BUILTIN\Administrators NT AUTHORITY\SYSTEM DESKTOP-LVS7839\swq])`
（AC#2 那条的红因里还额外出现本机外来主体 `S-1-5-21-…-1005`(CodexSandboxUsers) 与 `S-1-5-21-3623…`）。

**(b) 把解析器退回"跳过首行"那个形状**（`lines = lines[1:] // MUTATION-89ACC-3`，可编译、只改判据侧，实现一行未动）：

```
--- PASS: TestAC3SealedWritesCarryNoForeignSID        3/3 子测试全绿   ← 什么都没修，判据假绿
--- FAIL: TestAC1BaselineForeignACEPropagatesIntoModeOnlyWrites
    "the baseline is not leaking: the foreign ACE never reached …artifact.txt
     (principals [BUILTIN\Administrators NT AUTHORITY\SYSTEM DESKTOP-LVS7839\swq])"
--- FAIL: TestAC2SealedDirCoversFilesItNeverTouched   （仍红，见下）
```

⇒ **票面自曝的那次假绿是真的**：`icacls` 把第一条 ACE 和路径印在同一行（我的原文读数：
`…\data\artifact-exclusive Everyone:(I)(RX)`，第二条起才换行），跳过首行**恰好丢掉 `Everyone:(I)(RX)`**，
AC#3 三子在 chmod-only 实现上集体变绿。**这条反向钉子不是装饰**——它在解析器坏掉的同一刻变红，
红因写明了"本该漏的东西没漏"。
AC#2 在这个变异下**没有**跟着变绿，原因是本机 data 根的继承链上有**两个**外来主体，
丢掉首行还剩第二个 ⇒ 这条红是本机状态给的余量，不是判据给的；记下来免得下次以为 AC#2 能兜住解析器。

**(c) 那条钉子"永不 skip"**：它在 (a)(b) 两次运行里分别 PASS/FAIL，从未 SKIP；`-count=2` 全量 **0 SKIP**；
`grep -rn "t.Skip\|SkipNow\|testing.Short" internal/winsec/` = **0 命中**（结构上没有 skip 通道）；
`wideParent` 走 `run(t,"icacls",…)`，`icacls`/`powershell` 非零退出直接 `t.Fatalf`；
断言体是 `if privateACLError(…) == "" { t.Fatalf("the baseline is not leaking…") }` ⇒ **测不到漏就红**。

---

## 第 4 项：`PROTECTED_DACL` 变异 —— **PASS（红因是权限真的变宽），票面计数偏低**

`/tmp/wisp89acc89-m4`：`DACL_SECURITY_INFORMATION|PROTECTED_DACL_SECURITY_INFORMATION` → 只留前者（可编译）。
全量 `go test -count=1 -v ./internal/winsec/` **RC=1**，本代理读数：**8 条顶层 FAIL + 3 条子测试 FAIL = 11 红，2 顶层 PASS**
（`TestAC5FailedSealRefusesTheWrite`、`TestAC1BaselineForeignACEPropagatesIntoModeOnlyWrites`）。红名逐条：

```
FAIL TestAC5FailureIsNotSwallowedByTheHappyPath   FAIL TestAC4JunctionAtArtifactPositionIsNotRecursed
FAIL TestAC4UnlinkableLinkGivesNamedError         FAIL TestAC4SealedWalkSkipsLinks
FAIL TestAC1BaselineProductionPaths               FAIL TestAC3SealedWritesCarryNoForeignSID (+3 子)
FAIL TestAC2SealedDirCoversFilesItNeverTouched    FAIL TestAC3ProductionDataRootIsPrivateEndToEnd
```

⇒ 票面自称的"AC#3 三子 + AC#2 四红"**方向对、数报少了**（没提 AC#4 三条、AC#5 happy-path、生产端到端、AC#1 一起塌）。
记为**票面计数偏差**，不是假绿。

**红因是"权限真的变宽"而不是断言字符串没匹配**——OS 自己读回的 SDDL 原文（我日志里 4 处 `;;;WD`）：

```
D:AI(A;;FA;;;SY)(A;;FA;;;BA)(A;;FA;;;S-1-5-21-1228170099-895614386-1166154857-1001)
    (A;ID;0x1301ff;;;S-1-5-21-1228170099-895614386-1166154857-1005)      ← CodexSandboxUsers 全权继承进来
    (A;ID;0x1301ff;;;S-1-5-21-3623186960-731165060-4091685855-1717338598)
    (A;ID;FA;;;SY)(A;ID;FA;;;BA)(A;ID;FA;;;…-1001)
生产端到端那条还多一条： (A;OICIID;0x1200a9;;;WD)                          ← Everyone，票面点名的正是它
```

`D:AI` 里**没有 `P`**，且进来的不只是合成宽父的 `WD`——**本机真实外来主体 `CodexSandboxUsers`
带 `0x1301ff` 一起长进了"我们已经密封过"的文件**。运行时 `verifyPrivate` 当场拒绝落盘
（`DACL is missing or unprotected, so inherited grants still apply`），`memory.Open` 直接返回
`create data dir: … cannot apply a private security descriptor` ⇒ 失败方向 = 收紧，与 AC#5 一致。

---

## 第 5 项：AC#5 fail-closed —— **PASS（两条腿各自都被独立钉住）**

**(a) 去掉 `sealError()` 归一化**（`/tmp/wisp89acc89-m5`：`privateFile` 的 seal 失败分支 + `SealFile`/`SealDir`
三个公共边界全部直接返回平台错误；可编译）→ **RC=1，`9 PASS / 1 FAIL`**，唯一红的是

```
--- FAIL: TestAC5FailedSealRefusesTheWrite
    /exclusive_artifact /replacement_write /directory_chain /seal_file_and_dir_direct   ← 4/4 子全红
    红因：private_fail_test.go:38: error does not name the refusal: injected: descriptor could not be applied
```

⇒ **"契约不能靠每个实现记得 wrap"这句确实有用例兜住**（注入的错误故意**不**包 `ErrNotSealable`，
正是"忘了 wrap 的实现"的形状）。票面只报了 `sealHandle` 那枚变异（2 红），这条归一化变异它没跑；
本代理跑出来是 **4 红**，比票面更强。

**(b) 把"零字节残留"那半单独钉一次**（`/tmp/wisp89acc89-m7`：删掉拒绝分支里的 `_ = os.Remove(path)`，
错误照旧返回）→ **RC=1**：

```
--- FAIL: TestAC5FailedSealRefusesTheWrite/exclusive_artifact
      refused write left 0 bytes on disk at artifact.txt
--- FAIL: TestAC5FailedSealRefusesTheWrite/replacement_write
      refused write left 0 bytes on disk at blob.bin
    /directory_chain、/seal_file_and_dir_direct 仍 PASS（这两条不产生文件），happy-path 仍 PASS
```

⇒ `assertNoBytesOnDisk` 是承重的。**一条口径修正给票面**：残留物实测大小就是 **0 字节**
（封权限发生在写内容**之前**，`privateFile` 的顺序决定），所以"零字节残留"字面成立；这条腿真正防的是
"留一个宽 ACL 的空条目在磁盘上"（`O_TRUNC` 那条腿还会顺手清空旧内容），**不是**防敏感字节外泄——
防外泄靠的是顺序本身。别让后人把这条读成内容泄露钉子。

---

## 第 6 项：覆盖面主张 —— **一半 PASS、一半 FAIL（传播那一半没有任何用例，实测全绿）**

**(a) 那条"winsec 从没碰过的文件也变私有"的判据：找到了，断言本体是真测量。**
`internal/winsec/acl_windows_test.go:349-355`（`TestAC2SealedDirCoversFilesItNeverTouched` 内）：

```go
untouched := filepath.Join(root, "sqlite-like-sidecar.tmp")
if err := os.WriteFile(untouched, []byte("x"), 0o644); err != nil { … }   // 裸 os，winsec 不知情
assertPrivateACL(t, untouched)                                            // ← 断言本体
```

`assertPrivateACL` → `aclSIDs` → `run(t,"icacls",path)` → `privateACLError`（SID 白名单）。
⇒ 断言的确实是"这个文件在 OS 眼里只剩 {我,SY,BA}"，不是换个名字自夸。**这一半 PASS。**

**(b) "`SealDir` 对已存在子项传播"这条：没有任何用例钉住。** 变异 `/tmp/wisp89acc89-m6`：
`sealDir` 里删掉 `propagatePrivate(path)`（保留 `applyDescriptor`，可编译）⇒
**全量 `go test -count=1 -v ./internal/winsec/` RC=0，34 条 `=== RUN` 全绿，0 红**。
⇒ 票面 AC#2 那段"没有 walk 就会留下一整棵宽权限的老树"**在判据层面是空的**。

为什么它绿？本代理自造的第二台仪器（`internal/accprobe`，同一份代码分别放进 HEAD 快照与 m6）直接问 OS：
在**已密封的根**里造四个孩子，其中三个各带**自己的显式** `Everyone` ACE（`icacls /grant`，非继承），
然后只 `SealDir(root)`：

| build | `a-inherited-only.txt`（纯继承） | `b-own-explicit-ace.txt` | `c-dir`（带显式 ACE 的目录） | `c-dir/inside.txt` |
| --- | --- | --- | --- | --- |
| **HEAD（有 walk）** | 私有 | **私有** | **私有** | **私有** |
| **m6（无 walk）** | 私有（OS 自动重算） | **仍宽 `Everyone (RX)`** | **仍宽 `Everyone (OI)(CI)(RX)`** | **仍宽 `Everyone (I)(RX)`** |

⇒ **"先建目录后设权限没用"这条物理事实的真实边界是"子项 DACL 里有没有它自己的显式 ACE"**：
纯继承的孩子 OS 会替你重算（所以 AC#2 里的 `stray` 恰好被 OS 救了，判据测了个空）；
**带显式 ACE 的孩子只有 walk 能救**。实现是对的（HEAD 四个全私有），**判据漏的正是这一格**。

**(c) 还有一条历史账**：`3054007` 版 AC#2 里**本来有**钉这个物理事实的负向腿——
`if …== "" { t.Errorf("expected the pre-existing child to still carry the foreign ACE (proof that seal-before-write ordering matters)…") }`
（`git diff 3054007 0a3a445 -- internal/winsec/acl_windows_test.go` 里整段被删）。
`de15a6b` 加上 walk 之后那条腿**必然**变红，实现者的选择是**删判据**而不是换成"带显式 ACE 的孩子"这一格
⇒ 今天 suite 里连"出生宽、后来封"的正反面都不剩。

**缺的那一次具体测量（一句话可执行）**：AC#2 里补一个子项——先
`icacls <child-file> /grant *S-1-1-0:(RX)`（显式、非继承），再 `PrivateDirAll(root)`，
断言该子项 icacls 原文里**不再出现 `S-1-1-0`**；本代理已实测这条判据在删掉 `propagatePrivate` 的 build 上必红、
在 HEAD 上必绿。

---

## 第 7 项：junction 与 A51② —— **票面结论复现（自己的读数），但票面另一句被推翻**

自造仪器 `/tmp/wisp89acc89-junction`（独立 module，零仓依赖，全部 `t.TempDir()` 下，测完自清）：

```
mklink /J  -> Junction created for …\artifact-link <<===>> …\target   exec err: <nil>
Lstat(link) before remove : "?rw-rw-rw-"        ← Go 不把 junction 标成 ModeSymlink
target 里文件数 before=3 after=3                 (sub/keep-me.txt + sub/deeper/also.txt)
os.Remove(junction-to-non-empty-dir) error = <nil>
link 已消失=true ; target 仍在=true ; keep-me.txt 仍在=true
VERDICT-89ACC: os.Remove DID unlink the junction directly; A51(2) does not reproduce for junctions on this box
（空目标那条：os.Remove err=<nil>、linkGone=true、target 未被删）
```

⇒ **票 89 推翻票 79/A51② 的这条，本代理独立测得同一读数：能删，目标内容毫发无损。**
可以进 registry 更正（A51② 对 **junction** 不成立）。**清理**：两条 junction/符号链接都在自己的
`t.TempDir()` 里，测试内已拆；未在任何真实数据目录留下链接。

**但票面另一句被本代理推翻**：它说"**目录符号链接这台机普通权限建不出来**（需
`SeCreateSymbolicLinkPrivilege`/开发者模式）"。实测：

```
os.Symlink(dir → 非空目标)  error = <nil>       Lstat = "Lrw-rw-rw- (ModeSymlink set)"
os.Remove(该目录符号链接)    error = <nil>，目标内容仍在
IsInRole(Administrator) = False
HKLM…AppModelUnlock\AllowDevelopmentWithoutDevLicense = 1   ← 开发者模式是开着的
```

⇒ **这台机上目录符号链接本来就免特权可造**，AC#4 的"Windows 侧符号链接"那一格是**可构造而没构造**，不是不可构造。
缺的测量同样一句话：`reparse_windows_test.go` 里把现有 junction 用例的 `mklink /J` 换成 `os.Symlink`
（同一组断言：目标内容存活 + 只拆链接 + 失败具名）再跑一遍；并把"开发者模式下符号链接免特权"写进票面，
别再以"造不出来"推给 Linux 侧（换一台没开开发者模式的机器，这条前提会静默反过来）。
**顺带记一条实现正确性**：Go 的 `Lstat` 对 junction **不**给 `ModeSymlink`，所以 `isReparsePoint`
用 `FILE_ATTRIBUTE_REPARSE_POINT` 而非 `os.ModeSymlink` 是**必须的**——用后者密封游走会**漏掉所有 junction**；
这条他们做对了，但判据里没有区分两种链接的读数。
