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

## 进度总览（**全部 8 项 + AC#6 的 POSIX 腿已做完；结论在文末裁决表**）

| 项 | 内容 | 状态 |
| --- | --- | --- |
| 1 | 基线为假（本代理自造 icacls 仪器，票前生产码） | **已完成 → PASS** |
| 2 | 纯净快照 `go test -count=2 -v ./internal/winsec/` 四数 | **已完成 → PASS** |
| 3 | "差点假绿"独立证实（解析器退回跳过首行） | **已完成 → PASS（假绿复现，钉子有效）** |
| 4 | `PROTECTED_DACL` 变异 | **已完成 → PASS，票面计数偏低** |
| 5 | AC#5 fail-closed + `sealError()` 归一化拿掉 | **已完成 → PASS（两腿各测一次）** |
| 6 | 覆盖面主张（`SealDir` 传播 + 出生那一刻继承） | **已完成 → 一半 FAIL（传播无用例）** |
| 7 | `os.Remove` 删 junction 与票 79 A51② 的冲突 | **已完成 → 票面结论复现，另推翻票面一句** |
| 8 | `verifyPrivate` 白名单口径裁决 | **已完成 → 正确的严格（描述错位已纠）** |

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

---

---

## 补：AC#6 的 POSIX 腿独立复现（真 Linux，本代理自己跑）

票面自称"非 Windows 平台不是跳过，是测另一套机制（docker alpine:3.20，10 条具名结果、0 SKIP）"。
本代理不采信转述，自己交叉编译 + 真跑一次：

```
$ GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go test -c -o winsec.test ./internal/winsec/   # 在 0a3a445 快照
$ docker run -d --name wisp89acc89-posix -v wisp89acc89-89vol:/d alpine:3.20 sleep 900
$ docker exec … /d/winsec.test -test.v -test.count=1
Linux 944184f65f48 6.6.114.1-microsoft-standard-WSL2 x86_64 GNU/Linux
=== RUN  10   |   --- PASS 10（6 顶层 + 4 子）   |   --- FAIL 0   |   --- SKIP 0   |   RC=0
TestAC5FailedSealRefusesTheWrite(+exclusive_artifact/replacement_write/directory_chain/seal_file_and_dir_direct)
TestAC5FailureIsNotSwallowedByTheHappyPath / TestPOSIXPrivateFileIsReally0600 / TestPOSIXPrivateDirIsReally0700
TestPOSIXSymlinkAtArtifactPositionIsNotRecursed / TestPOSIXMissingFileIsNotAnError
```

断言本体也读了（`internal/winsec/private_other_test.go`），不是"跑过就算"：
`if got := info.Mode().Perm(); got != 0o600 { t.Errorf(...) }`（**传进去的是 0o644，要求出来是 0600**）、
`!= 0o700` 对 `PrivateDirAll` 建的**两级**目录都查、符号链接那条断言 `RemoveUnlinked` 后
`os.Stat(innocent)` **必须还在**（"TARGET DELETED - removal followed the symlink"）。
⇒ **"显式不适用 vs 静默通过"这条：PASS**，POSIX 侧是真测量、真断言，且没有 skip 分支。
容器与卷已清理（`0 container(s) / 0 volume(s) left`），名字带本代理会话后缀。

**顺带一条被本代理第 7 项推翻的注释**：`private_other_test.go:56-59` 与 `:77` 写着
"Windows needs a privilege for（构造符号链接）… the construction is stronger here than on Windows"。
这台机开发者模式=1、进程未提升，`os.Symlink(dir)` 直接成功 ⇒ **这句前提在本机为假**，
"Linux 侧更强"因此不成立；改判据时应把 Windows 侧的符号链接腿补齐（第 7 项给的方法），
而不是继续引用这句 privilege 说法。

---

## 第 8 项：`verifyPrivate` 只认 `{我, SY, BA}` —— **判：正确的严格；但票面对它后果的描述错位**

本代理不猜，直接在 HEAD 快照上做了一次测量（`internal/accprobe/probe2_windows_test.go`）：

```
IDENTITY u.Name="Carlos Shao" u.Uid="S-1-5-21-…-1001" u.Username="DESKTOP-LVS7839\swq"   ← 顺带纠我自己的仪器：u.Name 是显示名，不是 icacls 打印的账户名
操作员在密封根上主动加第三主体：icacls <root> /grant *S-1-5-6:(OI)(CI)(RX)
  AFTER-OPERATOR-GRANT root foreign=[NT AUTHORITY\SERVICE (OI)(CI)(RX)]
随后 winsec.PrivateFile(<root>/artifact.txt, …) = <nil>        ← 写成功，没有"在写入点炸"
再 winsec.SealDir(<root>) = nil，之后 root foreign=[]           ← 操作员那条有意授权被静默清除
```

⇒ 票面留的口径问题（"将来若有意的第三主体会在写入点直接失败"）**描述错位**：
`verifyPrivate` 读回的是**我们刚写的那个对象自己**的 DACL，而 PROTECTED DACL 让子项根本不继承父目录上的第三主体，
所以写入点**不会**失败；真正发生的是**下一次 `memory.Open`/`PrivateDirAll` 把运维有意留下的授权清掉**（上面实测第三行）。
**裁决：正确的严格，不改成"允许白名单外主体"。** 三条理由，不给和稀泥版：
1. 这个包唯一的承诺是"只有我"。任何形式的"可配置额外主体"都会把它变成"一个默认很严的 ACL 生成器"，
   而默认值一定会被放宽——票 89 的整票证据就是"默认继承"这三个字坑了我们多少个落点。
2. 严格性不是纸面姿态，它被 m4 变异证明承重：去掉 `PROTECTED` 之后本机真实主体
   `CodexSandboxUsers` 带着 `0x1301ff` 长进"已密封"的文件，`verifyPrivate` 是唯一当场拦住它的东西。
   白名单一开口，这条拦截就同时失效。
3. 未来的自己撞墙点**不在读取端，在共享端**：真要给别人读（服务账户读 artifacts 做报表、运维跨账户排障），
   现在的行为是**静默撤销**而不是**报错**。静默撤销比拒写更坏——运维会以为自己加了权限没用，然后去关 winsec。

⇒ 由此得出**本代理唯一要求的口径动作**（不是要求放宽）：把"winsec 拥有 data 树的唯一授权权、
任何带外授权会在下一次密封时被清除"写进 `docs/reports/pending-and-issues.md` 一条 A/`R` 记录，
并且**补一条用例声明这个语义**（我这次是探到的，suite 里没有；我给的探针现成可搬：加带外 ACE → `SealDir` → 断言它没了）。
如果哪天真要支持服务账户，要求是**加一个显式 API**（`PrivateFileSharedWith(sids …)`）+ 自己的判据，
而不是往现在这个白名单里塞例外。

---

## 覆盖面裁决：未接的同族落点，哪些算票 89 的账

本代理先把接线事实核了（`grep -rn "winsec\." --include=*.go | 排除测试与本包`）：
`agent/spill.go:111/:125/:254`、`memory/open.go:176/:188/:503`、`memory/artifacts.go:171/:181`、
`secret/store.go:49/:79` ⇒ **票面点名的四个家真的接上了**，不是只写了 API。
（另有一条我没在票面看到的：`internal/risk/winsec_c26.go:21` 已 `winsec.SetPathResolver(...)`，那是票 94 工作树，见下节。）

| 落点 | 算不算 89 的账 | 理由 / 代价 |
| --- | --- | --- |
| `internal/secret/migrate.go:154`（`os.WriteFile(backupPath, raw, 0o600)`） | **算（本代理新账，最重的一条）** | 它写的是**迁移前的原始 config**，里面按 `secret/migrate.go` 自己的注释还是"planned fields 未剥离"的明文密钥；票 89 已拥有并改过 `internal/secret/`（`store.go` 两处），同包同威胁模型（本机外来账户 `(M,DC)`）漏掉一半，"范围外"这个说法站不住。代价 = 一行替换 + 我第 1 项那台仪器加一类。**这条若编排者认账，就是 89 唯一"内容级"欠账。** |
| `internal/secret/migrate.go:168`（迁移后 tmp + rename） | **算，轻** | 同上（同函数），且它是 rename 前的临时件——正是票 89 AC#2 论证里"必须靠目录链兜住"的那一类形状。 |
| `models/downloader.go:224/:287/:566` + `archive.go:41/73/78` | **不算 89 的账；也不该按"私有数据"接** | 裁决见下"日志/模型缓存要/不要"。 |
| `config/migrate.go:83` + `config/parse.go:213` | **不算 89 的账** | 那是 `internal/config` 的写路径，**票 90 正在改这个包**；89 去接会撞车。转票 95/90 协调。 |
| `observe/logging.go:72/:244` | **不算 89 的账** | 日志从未进过 89 任何一条 AC 判据，票面也没 claim 过它。 |
| `ball/position.go:73/:77` | **不算** | 游戏位置状态，不是私有数据；`(M,DC)` 的后果是被人改分，归完整性票。 |
| `cmd/wisp/doctor.go:248`、`cmd/wisp/slo_windows.go:559` | **不算** | 票面自己的禁改区，票 77 活着。 |

**给编排者的一句话裁决（"日志/模型缓存算不算私有数据"）**：
- **日志：要**（但记在票 95，不是 89）。理由：日志里是 prompt 与 tool output，含用户粘过的东西，
  与 artifacts 同一威胁模型、同一台机、同一批 `(M,DC)` 主体；机密性主张成立。
  **代价与坑要点名**：`observe/logging.go:244` 是**追加**语义，而 `winsec.PrivateFile` 是
  `O_WRONLY|O_CREATE|O_TRUNC` ⇒ 直接替换会**把日志清空**。正确接法是"照旧 `OpenFile(O_APPEND)` 打开写完后
  `winsec.SealFile(path)`"（`SealFile` 已存在、已封已验证），代价 = 两处一行 + 一条 icacls 判据 + 轮转文件名跟着做。
- **模型缓存：不要按"私有数据"接**。理由：权重是公开可下载物（`internal/models/doc.go:11` 明说 hash 只来自签名 manifest，
  "mirrors serve bytes only"），机密性对它没有意义；而**完整性**已经有机制：`downloader.go:134` 的 cache 分支
  要求"每个已装文件 hash 为真"，`:264` 在离开 staging 前再验一遍 ⇒ 被同机账户偷改的字节在 cache-hit 时会被拒。
  残余风险只有"别人能看出你下载过哪些模型"（元数据），那是另一票。
  ⇒ **结论：不要，代价为零**；除非 owner 明确把"下载史"当隐私，那时接法也只是 `<DataDir>/staging` 一处 `PrivateDirAll`。

---

## `filepath.Abs`（D22 / 票 94）：一句裁决，不重复报

**不挡住本验收结论**。本代理全部读数取自信任 SHA `0a3a445`（`internal/winsec` 到 `0a3a445..HEAD` 无漂移，
`git diff --stat` 已核），`PrivateDirAll` 用 `filepath.Abs` 决定封哪棵树这件事**不改变**我测到的任何一条：
基线为假、密封后只剩 {SY,BA,我}、变异红名与红因。
并且本代理看到 94 的代理**已经在工作树里做了这件事**（未提交）：`git status` 显示
`?? internal/winsec/…`→`M internal/winsec/winsec.go`，且有新的 `internal/winsec/resolve.go`（`SetPathResolver`）
与 `internal/risk/winsec_c26.go` 把 `filepath.Abs` 换成 C26 解析。⇒ 这条**归票 94 交账**，89 不背。

---

## 替代判据到底够不够（编排者的头号问题，直答）

**机密性这一格：够。** 三条独立复现撑着：① 基线不是"没验证"，是当场为假（我自己的仪器，四类原文 `(M,DC)`）；
② 外来主体不是空组名——`CodexSandboxUsers` 里有 `CodexSandboxOffline`/`CodexSandboxOnline` 两个活本地账户，
这台机上还有 `CarlosShao`/`shaowq`/`Administrator`；③ 密封后 icacls 只剩 {SY, BA, 我}，且 `verifyPrivate`
每次写完都用 OS 自己的 SDDL 读回拒绝外来主体（m4 变异证明这条运行时腿承重）。
"第二账户真的 open 失败"能多证的是**内核 access-check 与 icacls 文本一致**，而它不改变主张的内容。

**但本代理要指出替代判据没覆盖的那一格，而且它不是"第二账户"**：
所有证据都在**新建文件**上取。真正没有测量的是"**既存宽权限树**被 winsec 接手之后是否真的变窄"——
票面 AC#2 的传播主张恰恰就是这一格，而它**无用例**（m6 全绿，第 6 项）。
⇒ **缺的那一次具体测量**（按价值排序）：
1. 既存子项**自带显式**外来 ACE（不是继承来的）→ `PrivateDirAll(root)` → 断言它没了。
   本代理的 `accprobe` 已证明这条区分 HEAD 与 m6，直接可搬进 suite。
2. Windows 侧**目录符号链接**那一格（第 7 项：这台机 `os.Symlink` 免特权可造，票面说不可造）。
3. `secret/migrate.go` 那两类既存产物（明文密钥备份）的 icacls 读数——第 1 项那台仪器加一类即可。
4. **不该再向本箱要**的：`runas` 那条。它在这台机器上不可闭合（未提升 + 无口令），
   票面 next=(4) 已如实写明；要求是把它从"可选升级"改记为**本箱不可闭合的边界 + 一条 CI/第二机待办**，
   别让下一任以为漏跑。

---

## 最终裁决表（AC 1:1 + 八项 1:1）

### 票面六个 AC

| AC | 票面声称 | 本代理裁决 | 依据（我自己的读数） |
| --- | --- | --- | --- |
| **AC#1** 基线 icacls 原文，四类各走生产入口，不许用 `FileInfo.Mode()` | 四类首条 ACE 都是外来 `CodexSandboxUsers (M,DC)` | **PASS（允许残留一条措辞）** | 第 1 项：票前树 `942ab5a` 上自造仪器，六条读数全中（含 db/wal/shm）；`Mode()` 原文 `-rw-rw-rw-` 也复现。残留：HEAD 里那个 `TestAC1BaselineProductionPaths` 因 `57bdbb2` 接线**已不再是纯基线**，名字要改或注释要写 |
| **AC#2** 方向（目录链先封 + 文件兜底），"建目录/建文件先后必须有用例钉住"，"继承不能放大要显式断言" | 逐层建+逐层封、walk 传播、PROTECTED | **FAIL（只在"用例"这一半；实现正确）** | "继承不能放大"这半 **PASS**：m4 变异 11 红、SDDL 里 `0x1301ff;;;S-1-5-21-…-1005` + `;;;WD` 真进来了。"先后必须有用例钉住"这半 **FAIL**：删掉 `propagatePrivate` 后 winsec 全套 34 条全绿（m6，RC=0）；且 `3054007` 原本钉它的负向腿在 `de15a6b` 被删而没换格。⇒ 补一条（本代理 accprobe 现成）即可转 PASS |
| **AC#3** 真红过的判据；替代判据 = SID 白名单 + icacls 原文 | 红→绿、0 SKIP、17 条具名 PASS | **PASS** | 第 2 项：纯净快照 `0a3a445` 自己跑 `-count=2 -v` → `=== RUN 34 / PASS 20(+14 子) / FAIL 0 / SKIP 0`，34=2×17 逐名核对。第 3 项：先红（RC=1，3 子 + AC#2）后绿，且假绿路径被独立复现、反向钉子当场变红。组里两个活账户 ⇒ 替代判据的红因落在真主体上 |
| **AC#4** A51② 单独钉：链接不被递归、失败具名；造不出来要如实写 | junction 可造、三条全绿；符号链接"本机造不出来"故放 Linux | **PASS（结论）/ FAIL（那句"造不出来"）** | 第 7 项：`os.Remove` 删指向非空目录的 junction **返回 nil、目标 3 文件全在** ⇒ 票面 A51② 更正成立（我自己的读数）；但同一台机、同一个未提升进程 `os.Symlink(dir)` **成功**（开发者模式=1）⇒ Windows 侧符号链接那一格是**可测而未测**。票面 AC#4 自身三条用例在我的 `-count=2` 里全绿，未递归删除由 suite 覆盖 |
| **AC#5** 失败方向只能收紧；注入一次失败 ⇒ 必须是错误 | 4/4 腿绿 + `sealError()` 归一化；票面变异 2 红 | **PASS（比票面更强）** | 第 5 项：m5 拿掉归一化 ⇒ **4/4 红**（红因 `error does not name the refusal`）；m7 拿掉 `os.Remove` ⇒ 2 红（`refused write left 0 bytes on disk`）⇒ 两条腿各自承重。口径修正：残留物本就是 0 字节，这条防的是"留下宽 ACL 空条目"，不是防内容外泄 |
| **AC#6** 门禁：两平台、四假绿逐条点名、非 Windows 不许静默通过 | gofmt/gofumpt/vet/`GOOS=linux` vet/test 全清，POSIX 侧真测 | **PASS（Linux 侧本代理自己跑过）** | Windows 侧：0 SKIP、0 FAIL、逐名 2×；`t.Skip` 在包内 0 命中（结构证据）。**Linux 侧见下"补：AC#6 的 POSIX 腿独立复现"**。gofmt/gofumpt/vet 三项是工具门不是判据，本代理未复核、不盖章 |

### 编排者点名的八项

| # | 项 | 裁决 |
| --- | --- | --- |
| 1 | 基线为假（自取 icacls） | **PASS** —— 四类 + wal/shm 全中 `(I)(M,DC)`，仪器与实现零共享 |
| 2 | AC#3 红→绿是真的（纯净快照四数） | **PASS** —— 34/20+14/0/0，17 名各 2 次；票面数字被复现而非照抄 |
| 3 | 自曝的"差点假绿"独立证实 | **PASS** —— 退回"跳过首行"⇒ AC#3 三子在未修实现上全绿；同一刻反向钉子变红；钉子无 skip 通道 |
| 4 | `PROTECTED_DACL` 变异 | **PASS（红因=权限真变宽）**，但**票面计数偏低**（自称 4 红，实测 8 顶层 + 3 子 = 11 红） |
| 5 | AC#5 fail-closed + 拿掉 `sealError()` | **PASS** —— 4/4 红；另测 `os.Remove` 腿 2 红 |
| 6 | 覆盖面主张（传播 + 出生那一刻） | **一半 FAIL** —— "winsec 从没碰过的文件也私有"那条判据在（`acl_windows_test.go:349-355`，断言本体是 icacls→SID 白名单）；但**传播无用例**（m6 全绿），真实边界是"子项自带显式 ACE 时只有 walk 能救"（accprobe 实测 HEAD 四条全收 / m6 三条仍宽） |
| 7 | 推翻 A51②（`os.Remove` 删 junction） | **PASS，可进 registry 更正**（我自己的读数）；**附带推翻票面另一句**：本机目录符号链接免特权可造 |
| 8 | `verifyPrivate` 白名单口径 | **正确的严格**（不放宽）；但票面对后果的描述错位——实测第三主体**不会**在写入点失败，而是**下一次密封时被静默清除** ⇒ 要求登记语义 + 补一条声明用例 |

### 两种"该不该算账"

- **未接同族**：只有 `secret/migrate.go:154/:168` **算票 89 的账**（同包、写的是迁移前**明文密钥**、票已 own 该包写路径）；
  `models/*`、`config/*`、`observe/*`、`ball/*`、`cmd/wisp/*` **不算**（分属票 95/90/77）。
- **日志/模型缓存要不要**：**日志要**（记票 95，接法必须用 `SealFile` 而不是 `PrivateFile`，否则清空日志）；
  **模型缓存不要**（公开字节 + manifest hash + cache-hit 复验已覆盖完整性，代价为零）。
- **票 94 的 `filepath.Abs`**：**不挡住验收结论**（94 工作树里已有 `resolve.go`+`SetPathResolver` 在处理），89 不背。

---

## 两句直答

**(1) 票 89 能不能改名 `-done`？——今天不能，但**不是**因为票 94。**
真实挡路的是本代理测出的**三条**未闭项，都在 89 自己 claim 的范围内：
① AC#2 的传播判据缺失（m6 全绿，票面把它写成"这条先后关系有用例钉"）；
② `secret/migrate.go:154/:168` 的明文密钥备份没接（票 own 该包）；
③ AC#4 "本机符号链接造不出来"这句**是错的**，那句写在票面上的话会直接被 registry 抄走。
⇒ 建议：①③ 一次 commit 就能补完（我给的 accprobe 现成），②要么现在接（一行级×2）要么编排者明文把它转给票 95
并**在票面划掉"范围外"这个说法**。这三条清完即可改名。票 94 未闭只影响 `internal/winsec` 的所有权交接顺序
（89 改名后 94 继续改这个包没冲突——`-done` 只是防重领键），**不构成理由**。

**(2)"我们的私有数据现在真的只对我可读吗"——今天这句**不能对 owner 无条件说**。**
可以说的是（每条都有本代理自己的 icacls 读数背书）：
"**artifacts、DPAPI blob、`wisp.db`+`-wal`+`-shm`、spill 的 `.retry` 这四条主路，在这台机上确实只剩
 `NT AUTHORITY\SYSTEM`、`BUILTIN\Administrators`、`DESKTOP-LVS7839\swq`；改之前它们对
 `CodexSandboxOffline`/`CodexSandboxOnline` 这类同机账户是可读可写的**"。
不可以说的是"所有私有数据"：`config.toml.bak`（迁移前含明文密钥）仍走装饰性的 `0o600`；日志仍 `0o644`；
**任何在 winsec 之前既存、且自带显式外来 ACE 的老树，只有跑到 `SealDir`/`PrivateDirAll` 才会被清**
（实现会清，但我没找到用例保证它一直清）。⇒ 对 owner 的诚实句式应带**范围状语**，不带"全部/从此"。

---

## 本代理的破坏性检查（自证清白）

* 仓内**未建 worktree**；所有变异与仪器都在 `/tmp/wisp89acc*` 快照里；`git status --short` 里
  除我自己的 `docs/evidence/s1/89-adversarial-acceptance.md` 之外**没有任何一条是本代理产生的**
  （`internal/config/`、`internal/panel/`、`cmd/wisp/`、`internal/winsec/`、`frontend/`、别人的票面一律未动）。
* 未 `git add -A`、未 `--amend`/`reset`/`rebase`/`stash`/`checkout .`、**未 push**。
* 未在用户真实数据目录写任何文件；junction / 符号链接全部造在 `t.TempDir()` 内并在测试内拆除。
* 四枚变异 token 只存在于快照：`MUTATION-89ACC-3/-4/-5/-7`；仓内 `grep -rn "MUTATION-89ACC" --include=*.go .` **= 0 命中（已复核）**，
  且 `git diff --quiet 0a3a445 -- internal/winsec/*_test.go` 干净（工作树里的 winsec 判据与本票终态逐字相同）。
* docker：容器 `wisp89acc89-posix` 与卷 `wisp89acc89-89vol` 跑完即删，复核读数 `0 container(s) / 0 volume(s) left`。
