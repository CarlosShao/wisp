# 票 89 复验（独立对抗验收 · 第二任）—— `acceptor-ticket89b`

**被验对象**：`.scratch/wisp/issues/89-0600-is-decorative-on-windows-acl-for-private-data.md`
（`Status: fix-complete-awaiting-reacceptance`）· 退回单四条 = 本文件同目录
`89-adversarial-acceptance.md`（上一轮 `acceptor-ticket89` 的面，本代理**不写它**）的"编排者·验收退回单"。
**被验的三枚 commit**：`c8d5c94`（退回单 #1 #2 判据+接线）/ `01e7007`（#3 #4 + 前三次变异读数）/
`f801d34`（#4 变异读数 + 全量门禁，纯 docs）。本轮 HEAD = `63fc82a`，
`git diff --name-only f801d34 HEAD -- internal/winsec internal/secret internal/memory internal/agent` **空**
⇒ 复验在 `f801d34`/HEAD 的码上取数即等价于被验三枚。

**环境**：Windows 11 Pro，`go version go1.27.1 windows/amd64`，当前账户 `DESKTOP-LVS7839\swq`
（SID 见下），`icacls` = `C:\Windows\System32\icacls.exe`。
**修复者的任何数字都不作为本代理的证据**：下面每一条都是本代理自己跑出来的，红名自己贴。

**快照纪律**：全部在仓外 `%TEMP%`（Git Bash 的 `/tmp`），仓内**未建 worktree**、工作树未被本代理改动。

| 快照 | SHA | 用途 |
| --- | --- | --- |
| `/tmp/wisp89bacc-89b2` | `c8d5c94` | 复现修复者变异 #1（删 `propagatePrivate`） |
| `/tmp/wisp89bacc-hd` | `63fc82a`(HEAD) | 四包 `-count=2` 回归 + 本代理在 HEAD 上重放变异 |

---

## 第 1 条：AC#2 覆盖面判据搬进包内 —— **复现成立（红名唯一、其余全绿），但票面 AC#2 的表述要降一档**

### 1.1 我自己重做那发变异（`c8d5c94`，与修复者同一 SHA、同一锚点）

- 锚点：`internal/winsec/winsec_windows.go:293` 的 `return propagatePrivate(path)` → `return nil // MUTATION-89BACC-WALK`
 （`applyDescriptor` 一行未动，walk 整体留着但没人调它）。
- **同链证明**：`grep -n "MUTATION-89BACC-WALK"` ⇒ `293:	return nil // MUTATION-89BACC-WALK`（就是锚点那一行）。
- **先证可编译**（编译失败不算变异）：变异前 `go build ./...` **rc=0**，变异后 `go build ./...` **rc=0**。
- 跑：`go test -count=1 -v ./internal/winsec/` ⇒ **TEST_RC=1**。

| 读数 | 本代理实测 |
| --- | --- |
| `=== RUN` | **25** |
| 顶层 `--- PASS` | **14** |
| 顶层 `--- FAIL` | **1** |
| 缩进 `    --- PASS`（子测试） | 10 |
| `SKIP`（全量输出 grep -c） | **0** |

**唯一红名**：`--- FAIL: TestAC2SealDirNarrowsChildrenThatCarryTheirOwnExplicitACEs (0.45s)`
⇒ 与修复者自称的"rc=1、`=== RUN` 25、14 绿 / 1 红、0 SKIP、红名就是这一条"**逐字对上**。

红因（OS 自己的原文，本代理从日志逐字取）：

```
acl_windows_test.go:454: icacls text of operator-widened-file.txt still names the foreign principal (Everyone):
    …\TestAC2SealDirNarrowsChildrenThatCarryTheirOwnExplicitACEs…\data\operator-widened-file.txt Everyone:(RX)
        NT AUTHORITY\SYSTEM:(I)(F)
        BUILTIN\Administrators:(I)(F)
        DESKTOP-LVS7839\swq:(I)(F)
acl_windows_test.go:455: NOT PRIVATE operator-widened-file.txt: foreign SID(s) S-1-1-0
    (principals: [Everyone NT AUTHORITY\SYSTEM BUILTIN\Administrators DESKTOP-LVS7839\swq])
```

`operator-widened-dir` 与 `operator-widened-dir/inside.txt` 同形红（三条都带显式 ACE）。

**还原自证**：把快照的该文件用 `git show c8d5c94:…` 覆盖回去后
`grep -rn "MUTATION-89BACC" --include=*.go .` **0 命中**、`diff` 与 `c8d5c94` 逐字相同、
仓内 `git diff --quiet -- internal/winsec/winsec_windows.go` **rc=0**（工作树未做变异）。

### 1.2 变异在 HEAD 上打到哪几格（本代理补测，修复者没报这一格）

同一锚点在 HEAD 快照 `/tmp/wisp89bacc-hd`（先证 `go build ./...` **rc=0**，变异后仍 **rc=0**，
`grep -n` 命中 `382:	return nil // MUTATION-89BACC-WALK`）⇒ `go test -count=1 -v ./internal/winsec/`
**TEST_RC=1、`=== RUN` 29、17 顶层 PASS / 2 顶层 FAIL / 10 子 PASS、0 SKIP**。
**HEAD 上红名是两条不是壹条**（修复者的"唯一红名"是在 `c8d5c94` 上取的，那时第 4 条的用例还不存在）：

```
--- FAIL: TestSealReportsThePrincipalsItCleared (0.07s)
--- FAIL: TestAC2SealDirNarrowsChildrenThatCarryTheirOwnExplicitACEs (0.45s)
```

第一条的红因（本代理从自己日志逐字取）：
`seal cleared a grant on shared-with-a-service-account.txt without reporting it` +
`shared-with-a-service-account.txt still is not private after the reported narrowing: … D:AI(A;;0x1200a9;;;WD)…`
⇒ walk 在 HEAD 上**多了一格承重**（通知不能报 walk 没碰过的东西），覆盖面判据比 `c8d5c94` 更强不是更弱。
还原后 `grep -rn MUTATION-89BACC` **0 命中**、`diff` 与 HEAD 逐字相同。

### 1.3 判"注释里那句'纯继承的孩子不区分任何事'"—— **是诚实，不是覆盖面自我缩窄**

本代理不采信注释，自己量了一次：在 §1.1 的**无 walk 变异**日志里，`operator-widened-file.txt` 前后两次 icacls 原文是

```
（前置断言处，SealDir 之前）      Everyone:(RX)  +  Everyone:(I)(RX)  + SY/BA/我 (I)(F)
（SealDir 之后，无 walk 的 build） Everyone:(RX)                        + SY/BA/我 (I)(F)
```

⇒ **继承来的那一条 `Everyone:(I)(RX)` 被 OS 自己重算掉了，自带显式的 `Everyone:(RX)` 留着** ——
注释那句"纯继承的孩子在断言上不区分任何事"是**被本代理独立复现的物理事实**，不是给自己找的台阶。
判据的形状也因此是对的：它挑的三条子项（file / dir / dir 里的 file）**全部自带显式 ACE**，
对照项 `inherits-only.txt` 只 `t.Logf` 不当判据。

**票面 AC#2 的表述要不要降一档：要，降半档（措辞级，不是结论级）。** 票面 AC#2 第 (2) 句写
"'先建目录、后设权限'**没用**：子项继承的是创建那一刻的 ACL"。本代理读数给出的准确边界是：

| 既存子项的形态 | 只重封父目录（不 walk）会发生什么 |
| --- | --- |
| 纯继承（ACE 带 `ID`） | OS 重算 ⇒ **会变私有**（"后设权限"在这一格其实管用，代价是中间那段可读窗口） |
| 自带显式 ACE | **仍然宽**，只有 walk 能救（上面那两行原文就是这一格） |

⇒ "边建边封"的真实理由应写成"**不留可读窗口** + **修自带显式 ACE 的老孩子**"，而不是"后设权限没用"。
这是**票面措辞**要改；`AC#2` 这一格本代理判 **PASS**（退回单第 1 条的两条要求逐字复现）。

---

## 第 2 条：明文密钥备份 —— **接线本代理自己复现（独立仪器、同一对象前后 icacls 对照），另登记四条残留与更正**

本代理**不 import `internal/winsec` 的测试helper**（它们在 `winsec_test` 包里也 import 不到），
自己写了 `internal/acc89b/mig89b_windows_test.go`（只在仓外快照里，仓内没有这个目录）：
自己的 icacls 调用、自己的 SID 白名单（名字→SID 走 powershell 翻译，不做字符串猜名）、
自己的宽父目录构造（`icacls /inheritance:r /grant:r …*S-1-1-0:(OI)(CI)(RX)`）。
**构造全在 `t.TempDir()` 下，没有往 `%APPDATA%` 的真实数据目录写一个字节。**

### 2.1 本代理自己的 icacls BEFORE/AFTER 对照（HEAD 干净码，走生产入口 `secret.MigratePlaintext`）

```
BEFORE-migration  …\data\config.toml                 Everyone:(I)(RX) + Administrators/SYSTEM/swq (I)(F)
AFTER-backup      …\data\config.toml.bak-plaintext   NT AUTHORITY\SYSTEM:(F)
                                                      BUILTIN\Administrators:(F)
                                                      DESKTOP-LVS7839\swq:(F)
AFTER-config      …\data\config.toml                 同上三条，Everyone 不见了
```

⇒ 备份里那份明文（本代理断言过它确实还带着 `sk-PLAINTEXT-89b-DO-NOT-LEAK`）落盘即私有；
迁移后的 `config.toml` 也私有（tmp 封了 + rename 保留 descriptor 这一句被读数支持）。

既存宽备份那一腿（repair leg），**同一条路径前后对照**：

```
BEFORE-repair …\config.toml.bak-plaintext   Everyone:(RX)  +  Everyone:(I)(RX) + SY/BA/我
AFTER-repair  …\config.toml.bak-plaintext   NT AUTHORITY\SYSTEM:(F) BUILTIN\Administrators:(F) DESKTOP-LVS7839\swq:(F)
```

同一次运行里还打出本代理**没要求**但白捡的一条第 4 条证据（默认通知器真的在真生产路径上响了）：

```
WARN winsec: seal cleared principals that were placed on this object explicitly
  path=…\config.toml.bak-plaintext cleared=WD(A;;0x1200a9;;;WD)
  policy="winsec owns the grants on this tree; out-of-band ACEs are removed at the next seal"
```

### 2.2 两发变异本代理自己重做（编译都先证 rc=0，锚点在同一条 grep 链里）

| 变异 | 锚点 | 落地证明 | 构建 | 本代理仪器 | 修复者的用例 |
| --- | --- | --- | --- | --- | --- |
| `MUTATION-89BACC-MIG` 两处 `winsec.PrivateFile` 退回 `os.WriteFile(…,0o600)` | `migrate.go:169` + `:189` | `grep -c` = **2** | `go build ./internal/secret/` **rc=0** | `TestProbeFreshMigrationBackupIsPrivate` **红**：`NOT PRIVATE config.toml.bak-plaintext: [Everyone [S-1-1-0] (I)(RX)]`，且 `config.toml` 也红 | `--- FAIL: TestAC2MigrationBackupIsPrivate`，红因 `NOT PRIVATE config.toml.bak-plaintext: foreign SID(s) S-1-1-0`（与修复者贴的逐字同形） |
| `MUTATION-89BACC-REPAIR` 把 `else if serr := winsec.SealFile(backupPath); serr != nil` 换成永不触发的等价式（**不调用** SealFile） | `migrate.go:174` | `grep -n` 命中 174 行 | build **rc=0** | `TestProbePreExistingWideBackupIsRepaired` **红**：AFTER 仍是 `Everyone:(RX) + Everyone:(I)(RX)` | `--- FAIL: TestAC2PreExistingMigrationBackupIsRepaired`，红因 `foreign SID(s) S-1-1-0, S-1-1-0 (principals: [Everyone Everyone …])` |

⇒ 两条腿各自承重，`TestAC2SealedDirCoversFilesItNeverTouched` 与 `TestAC2MigrationBackupIsPrivate`
在后一发变异下保持绿 ⇒ 变异只打到该打的那一格。两发都还原后 `grep -rn MUTATION-89BACC --include=*.go` **0 命中**、
`diff` 与 HEAD 逐字相同。

**本代理自己补的一条断言**：变异 #MIG 下 `config.toml`（迁移后的主配置）**也**是宽的 ⇒
说明"`:189` 那一行"不是凑数，rename 前不封就等于 rename 后仍宽；修复者两条都接是对的。

### 2.3 残留登记（**只登记，不自己修**）

- **(R-89b-1) 明文在别处不留副本：这一格本代理量出来是干净的。** `MigratePlaintext` 跑完后
  在本代理那个宽根下**全树扫字节**（含 `secrets\` 里的 DPAPI blob、`artifacts`、临时件），
  除 `config.toml.bak-plaintext` 之外**命中 0 份** `sk-PLAINTEXT` 明文（日志原文
  `PROBE plaintext copies under … besides the backup: 0`）；`config.toml.migrate-tmp` 在 rename 后**不存在**
  （`GetFileAttributesEx … The system cannot find the file specified.`），失败路径上 `migrate.go:193`
  也 `os.Remove(tmpPath)`。⇒ tmp 那一格同样密封了，没有中间态残留。
- **(R-89b-2) 同类明文站点仍在仓里，且不属于这三枚 commit**：`internal/config/migrate.go:83` 的
  `os.WriteFile(backup, raw, 0o600)`（`config.toml.bak-<from>`，schema 迁移前的原始配置）与
  `internal/config/parse.go:213` 的 `os.Chmod(tmpName, 0o600)`（`atomicWrite` 的随机名 tmp）
  仍是装饰性模式位。上一轮已裁"`internal/config` 归票 95/90"，本代理**不重开地界**，
  但它限制了"对 owner 那句话"的范围（见文末直答）：**真正跑在生产装配根里的那份配置备份还没有封**。
- **(R-89b-3) 被这次修复接上线的函数目前没有任何生产调用方**：
  `grep -rn "MigratePlaintext" --include=*.go .` 排除自身与 `_test.go` 只剩
  `internal/secret/configrefs.go:16` 的一句**注释** ⇒ D33 明文迁移在装配根里**还没接线**。
  ⇒ "迁移备份已密封"这句话今天只能说"**它被产出的那条路径上**密封了"，不能说"owner 机器上会出现这种备份，它们是私有的"。
  这不是缺陷（判据与接线都对、两发变异都红），是**口径的射程**，必须写进结论。
- **(R-89b-4) 行为变化（方向是紧的，登记不判红）**：`PrivateFile`/`SealFile` 走 C26 铸造口，
  相对路径与含 `..` 的拼法**直接拒**。本代理实测 `secret.MigratePlaintext("config.toml")`（chdir 到该目录）
  现在**失败**：`secret: migrate: secret: create secrets: winsec: refusing to seal secrets: … secrets is not absolute`。
  **归因取了三档读数**（这是本代理不写"大概是谁的锅"的那一格）：

  | 快照 | SHA | 同一条相对路径探针 |
  | --- | --- | --- |
  | `0a3a445`（票 89 上一轮终态，票 94 未落） | 相对迁移 **PASS**（`PROBE relative migration ok: backup=config.toml.bak-plaintext`） |
  | `c8d5c94^` = `468dd27`（票 94 的 `resolve.go` 已在，本次两行接线未落） | 相对迁移 **FAIL**，错误与下面逐字同形 |
  | HEAD（三枚全在） | 相对迁移 **FAIL** |

  ⇒ **不是 `c8d5c94` 这两行接出来的**，是票 94 的解析口带来的语义（`secret.NewStore`→`PrivateDirAll` 先拒）；
  同一份 `config.toml.bak-plaintext` 在 `0a3a445` 上还是**宽的**（本代理在 0a3a445 上跑同一仪器：
  `NOT PRIVATE config.toml.bak-plaintext: [Everyone [S-1-1-0] (I)(RX)]` + `NOT PRIVATE config.toml after the rename` 两条红）
  ⇒ 这条同时把上一轮验收的"最坏产物没封"那笔账**用本代理自己的 icacls 原文**复现了一遍。

---

## 第 3 条：目录符号链接那一格本代理自己复测 —— **修复者说的是真的；A51② 的旧登记要更正，判据给在下面**

本代理自己写仪器（`internal/acc89b/probe34_windows_test.go`，全在 `t.TempDir()` 里，未碰任何真实数据目录），
不 import 被测包的 helper。原始读数（本代理的 `t.Logf` 行，逐字）：

```
PRIV IsInRole(Administrator) = False
PRIV dev-mode = HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\AppModelUnlock
        AllowDevelopmentWithoutDevLicense    REG_DWORD    0x1
PRIV whoami = swq
SYMLINK os.Symlink(non-empty dir) error = <nil>
LSTAT mode=Lrw-rw-rw- ModeSymlink=true reparse=true FILE_ATTRIBUTE_REPARSE_POINT(0x400) set=true
A51-RECHECK os.Remove(dir symlink to non-empty dir) err = <nil> ; target still there = true
TARGET-INTACT after RemoveUnlinked
WALK-DID-NOT-CROSS target still wide: [Everyone [S-1-1-0] (OI)(CI)(RX)
    DESKTOP-LVS7839\CodexSandboxUsers [S-1-5-21-…-1005] (I)(OI)(CI)(M,DC)
    S-1-5-21-3623186960-…-1717338598 (I)(OI)(CI)(M,DC)]
JUNCTION mklink /J -> Junction created for …\junction <<===>> …\target
JUNCTION Lstat mode=?rw-rw-rw- ModeSymlink=false reparse-bit=true
```

⇒ 四条都成立，**且是本代理自己测的**：① 这台机**未提权 + 开发者模式=1 ⇒ `os.Symlink` 到非空目录直接成功**；
② `Lstat` 给 `Lrw-rw-rw-`、`ModeSymlink` 与 reparse 属性**同时**置位；
③ `os.Remove` 对**目录符号链接**也**能**直接拆、目标内容毫发无损；
④ `SealDir`/`PrivateDirAll` 的传播游走**不穿**符号链接（链接目标那条宽 DACL 原样留着，
连本机那两个真实外来主体的 `(M,DC)` 都没被动过）。
另外 `JUNCTION` 那一发独立复现了"Go 的 `Lstat` 对 junction **不**给 `ModeSymlink`、只给 reparse 位"
⇒ `isReparsePoint` 用属性而不是 `ModeSymlink` 是**必须的**（这条实现做对了，本代理支持它继续留在判据里）。

**仪器自曝一条（不拿它当结论）**：本代理想补第三档证据 `whoami /priv`，在这台机上取到**空输出**
（`cmd /c whoami /priv` 从 Go 与从 MSYS 两个方向都是 1 行空表头）。
按本仓"空输出先怀疑仪器"的规矩，这一档**不计**：结论建立在 `IsInRole=False` + 注册表 `0x1` +
`os.Symlink` 返回 `<nil>` 三条实测上，这三条本代理都拿到了原文。

### A51② 那条旧登记要不要更正：**要，改两处，判据与原文读数如下**

`docs/reports/pending-and-issues.md` 的 A51② 现文（本代理读到的口径）说的是
"`os.Remove` 删不掉指向目录的符号链接 ⇒ 游离子树回收清不掉"。本轮读数：

1. **junction**：上一轮验收与本代理本轮**各自独立**测到 `os.Remove` 返回 `<nil>`、目标内容全在
   （本轮原文见上 `JUNCTION …` 两行 + 上一轮报告的 `VERDICT-89ACC` 行）；
2. **目录符号链接**：本代理本轮原文 `A51-RECHECK os.Remove(dir symlink to non-empty dir) err = <nil> ; target still there = true`。

⇒ 落账判据（**编排者的面，本代理不写 registry**）：A51② 应改成
**"在本机 Go 1.27.1 / Win11 上不复现（junction 与目录符号链接两类对象都实测 `os.Remove` 能直接拆、目标内容存活）；
保留 `RemoveUnlinked` 的理由从'删不掉'改成'只可能删到链接本身、绝不把删除半径交给别人'"**。
第二句是必须的：如果只把"删不掉"划掉而不换理由，后人会顺手把 `RemoveUnlinked` 当成多余代码删了。
**registry 现文的本代理读数**（`docs/reports/pending-and-issues.md:1865`，逐字）：
"**A51② `os.Remove` 删不掉"指向目录的符号链接"** ⇒ 一个裸名 artifact 若被替换成 symlink-to-dir，
票 79 的游离子树回收路径会清不掉它"——**这句在本机为假**；`A64③(1)`（同一文件 `:1453`）只更正了
**junction 一半**（"能删指向非空目录的 junction"），**目录符号链接那一半还挂着**，
且两处都没换掉"`RemoveUnlinked` 存在的理由"。⇒ 落账动作 = 改 `:1865` 本体 + 补 `:1453` 的第二类对象 + 换理由句。

---

## 第 4 条：本代理攻"报出来"这条 —— **白名单按 SID 成立、SYSTEM/Administrators 不噪；但"继承来的那份"确实绕过检测（登记为射程限定，不判 FAIL）**

### 4.1 攻"挂在父目录上继承下来（SDDL 带 `ID`）会不会绕过检测"：**会，实测确认**

本代理的仪器（`TestProbeNoticeMissesInheritedGrants`）：`PrivateDirAll(root)` → 操作员在 **root** 上
`icacls /grant *S-1-1-0:(OI)(CI)(RX)` → 在 root 里建 `artifact.txt`（它**存着**一份
`Everyone:(I)(RX)`，本代理的前置断言验过）→ **只**调 `winsec.SealFile(child)`（生产里
"某个工件被重新落一次盘"就是这个形状，`memory.Open` 不会被触发）⇒

```
NOTICE-LINES for a SealFile that removed an inherited copy of an operator grant: 0 ([])
AUDIT GAP CONFIRMED: the child lost read access for Everyone and nothing named it
```

同一发里 `foreignPrincipals(child)` 从"宽"变成"空" ⇒ **子项上那份带外授权确实被清掉了，一条日志都没有**。

**判：这一格不是 FAIL，是这条通知的射程比它的名字窄。** 三条理由：
① 机密性方向是**紧**的（对象只会变窄，不会变宽）；② 那份 ACE 的**来源**（root 上那条显式的）
在 root 下一次被密封时**会**报（§4.2 每一发报的就是它）；
③ 代码注释与 `TestSealReportsThePrincipalsItCleared` 明写并且**钉住**了"纯继承的孩子不出现在通知里"
⇒ 这是**已经写在纸面上的取舍**，不是偷偷漏。
但**必须落一句账**（R-89b-5）：票面/registry 里"带外授权会在下一次密封时被清除，**并且打一条 WARN**"
这句要改成——"**记在该对象自己 DACL 上的（SDDL 不带 `ID`）授权被清时才打 WARN；
只继承自父目录的那一份不打**（它随父目录被密封时由父目录那一条承担），
所以'某个孩子悄悄读不到了'这一格今天没有日志"。

### 4.2 攻"SYSTEM/Administrators 不该报（避免噪声）"与"白名单按 SID 而不是按名字字符串"：**两条都过**

九发对照（每发独立 `t.TempDir()` 根，`icacls /grant <拼法>:(OI)(CI)(RX)` 后 `SealDir(root)`，
数 `seal cleared principals` 那一条 WARN）。**本代理一律把 grantee 先翻成 SID 再落地**（`sidOfName` 走
`NTAccount→SecurityIdentifier`），这样构造本身就不依赖 icacls 会不会解析某个显示名：

| grantee 拼法 | 期望 | 实测 WARN 条数 | 通知里的 `cleared=` |
| --- | --- | --- | --- |
| `*S-1-5-18`（SYSTEM，SID） | 不报 | **0** | — |
| `*S-1-5-32-544`（Administrators，SID） | 不报 | **0** | — |
| `BUILTIN\Administrators`（**名字**） | 不报 | **0** | — |
| `DESKTOP-LVS7839\swq`（**我，名字**） | 不报 | **0** | — |
| `*S-1-1-0`（Everyone，SID） | 报 | 1 | `WD(A;OICI;0x1200a9;;;WD)` |
| `Everyone`（**名字**） | 报 | 1 | `WD(A;OICI;0x1200a9;;;WD)` |
| `*S-1-5-6`（SERVICE） | 报 | 1 | `SU(A;OICI;0x1200a9;;;SU)` |
| `*S-1-5-11`（Authenticated Users） | 报 | 1 | `AU(A;OICI;0x1200a9;;;AU)` |
| `DESKTOP-LVS7839\CodexSandboxUsers`（本机真账户组） | 报 | 1 | `S-1-5-21-…-1005(A;OICI;0x1200a9;;;S-1-5-21-…-1005)` |

⇒ **白名单确实按 SID 判**：`allowedSIDStrings()` 的键是 `SY`/`S-1-5-18`/`BA`/`S-1-5-32-544`/`user.Current().Uid`，
**一个账户显示名都不在里面**；SDDL 由 OS 产出（缩写或裸 SID，**从不出显示名**）
⇒ "本地化机器上名字会变"打不到这条判据。
最硬的一发是 `Everyone`（名字形式）与 `*S-1-1-0`（SID 形式）**报出同一个 `WD`** ⇒ 检测不看操作员怎么拼，
看 OS 把 trustee 落成什么。`CodexSandboxUsers` 那一发同时证明：能被 icacls 解析出来的本机组，
通知里仍以**裸机器 SID** 出现（与 AC#1 那条"解析不出名字的别名 SID"同源，两条腿都覆盖到了）。

### 4.3 本代理另外两发攻击（都干净，写出来免得下一任再猜）

- **DENY ACE**：`icacls root /deny *S-1-1-0:(OI)(CI)(RX)` 后 `SealDir(root)` = **`<nil>`**，
  并且**报了**：`WARN winsec: … cleared=WD(D;OICI;0x1200a9;;;WD)`。
  ⇒ 本代理原本怀疑"`verifyPrivate` 把 `fields[0] != "A"` 当外来主体 ⇒ 带 DENY 的对象会被拒绝密封、
  运维可以让我们自我锁死"——**在本机不成立**：整个 DACL 被换掉，读回的是换后的那份，没有 DENY 可挑。
- **备份位置被人放了一根符号链接**（攻第 2 条新接线的分支判定：`os.Stat(backupPath)` 走**未解析**拼法，
  两条腿的动作都走**已解析**拼法，问它们不一致时会怎样）：
  `BRANCH-A`（`config.toml.bak-plaintext` 是指向目录的链接）与 `BRANCH-B`（指向别人家的文件）
  **两次都拒绝**：`winsec: refusing to seal … traverses a reparse point …, which is not the tree this call names`，
  并且 `VICTIM-INTACT the foreign file behind the link is untouched`、链接目标那条宽 DACL 原样保留。
  ⇒ **没有写穿链接**，票 94 的牙没被这两行接线卸掉。代价登记为可用性（R-89b-6）：
  备份位置上有一根既存链接 ⇒ 迁移**永久失败**（错误具名、带路径，方向是紧的）。

---

## 回归：纯净快照四包 `-count=2`（本代理自己的四数）

快照 `/tmp/wisp89bacc-hd` = `63fc82a`（HEAD），先证 `go build ./...` **rc=0**，再跑
`go test -count=2 -v ./internal/winsec/ ./internal/secret/ ./internal/memory/ ./internal/agent/` ⇒ **RC=0**。

| 包 | 结果 | 用时 |
| --- | --- | --- |
| `internal/winsec` | `ok` | 10.735s |
| `internal/secret` | `ok` | 0.361s |
| `internal/memory` | `ok` | 28.000s |
| `internal/agent` | `ok` | 4.268s |

| 计数（全量 `-v` 输出） | 数 |
| --- | --- |
| `=== RUN` | **402** |
| 顶层 `--- PASS` | **268** |
| 缩进 `    --- PASS` | **132** |
| `--- FAIL`（顶层+子） | **0** |
| `--- SKIP`（顶层） | **2** |

**`=== RUN` == 2 × 不同名**：distinct = **201**，`uniq -c` 后**没有任何一个名字的次数 ≠ 2**
（`grep '^=== RUN' | sed 's/^=== RUN   //' | sort | uniq -c | awk '$1!=2'` 输出为空）。

**SKIP/FAIL 逐条点名**：FAIL **0 条**；SKIP 只有**一个名字出现两行**——
`--- SKIP: TestSubprocessCrashWriter (0.00s)` × 2，
原因 `concurrent_test.go:186: crash-writer subprocess; runs under TestCrashRecoveryKillMidWrite`
（`WISP_CRASH_CHILD != "1"` 时子进程入口在父进程里跳过，真腿由 `TestCrashRecoveryKillMidWrite` 带起来跑；
`git log -1 -- internal/memory/concurrent_test.go` = `41d707c`，**票前既有**，与这三枚 commit 无关）。

⚠ **(R-89b-7) 对修复者一句读数的更正**：票面收尾枚写"全量输出里 `--- SKIP` **0 条**（用 `grep -E "SKIP|FAIL"` 扫全量输出计数，不是 `head`）"。
本代理实测：**不加 `-v` 时 `go test` 根本不印 `--- SKIP`**——
`go test -count=2 ./internal/memory/` 的全部输出就是 `ok … 25.731s` 一行（`grep -c SKIP` = **0**），
同一包加 `-v` 就有 **2** 行。⇒ "0 SKIP" 若是从非 `-v` 输出扫出来的，那是**仪器看不见**，不是**没有**。
真数是 **1 个具名 helper skip × 2 轮**。这一格不判红（skip 本身良性且既有），
但 AC#6 那句"skip 要逐条点名，不许当 ok"要求票面把它改成点名读数。

工具门（本代理复核，只当门不当判据）：`gofmt -l internal/winsec internal/secret` 空、
`$GOPATH/bin/gofumpt -l` 同两包空、`go vet` 与 `GOOS=linux go vet`（**按包作用域**，A54③）两包 rc=0。

---

## 票 94 的牙：修复有没有绕开 `ResolvedPath` 铸造口 —— **没有，判 PASS（并附一发攻击读数）**

- **没有新增公开入口**：`git diff 01e7007^..01e7007 -- internal/winsec/winsec.go internal/winsec/resolve.go` 为空，
  第 4 条只动 `applyDescriptorWindows`；第 2 条用的是既有公开 API `winsec.PrivateFile`/`winsec.SealFile`。
- **密封路径上每一个会碰 OS 的入口都过了铸造口**（本代理逐个数调用点）：
  `PrivateFile`/`PrivateFileExclusive`→`privateFile`→`resolveString`；`SealFile`/`SealDir`→`resolveString`；
  `PrivateDirAll`→`ResolvePath` 后走 `privateDirAll(dir ResolvedPath)`；
  `applyDescriptor` 的三个调用点分别拿 `f.Name()`（已用 resolved 拼法打开）、`sealDir(resolved)`、
  `propagatePrivate` 里 `filepath.Join(resolved, ReadDir 的名字)`。新增的那一发
  `windows.GetNamedSecurityInfo(path,…)` 就挂在 `applyDescriptorWindows` 里，拿的是同一个已过铸造口的 `path`。
  唯一故意不解析的仍是 `RemoveUnlinked`（票 94 自认弱处，本三枚**未碰**）。
- `diff 0a3a445..f801d34` 生产码里**没有**新增 `filepath.Abs` / `os.Chmod` / 裸 `os.Create`（`grep` 命中全在测试与注释里）。
- **本代理的攻击读数**（§4.3 第二发）：把备份位置换成符号链接，让"分支判定用的未解析 `os.Stat`"与
  "动作用的解析后拼法"不一致 ⇒ 两条腿**都拒绝**并给出 `traverses a reparse point …, which is not the tree this call names`，
  链接后面那个别人的文件**内容未动、DACL 未动**。⇒ 铸造口在这条新接线上仍然说得过去。
- 游走不穿链接也独立复现（§第 3 条 `WALK-DID-NOT-CROSS`）。
- **残留（不判红）**：`migrate.go:168` 用**未解析**的 `os.Stat(backupPath)` 选分支这件事本身是"词法决定、动作解析"的错位；
  本代理测到它目前只会走向拒绝（紧），但它是一句该写进注释的弱点，而不是该留白的巧合。
- 另一处口径：通知把**完整路径**写进 `slog`，而日志文件今天还是 `0o644`（票 95 的账）。路径不是密钥，
  但"哪台机器的哪个目录里有你的私有树"这件事在日志里是可 grep 的，归票 95 一起算。

---

## 第 4 条的变异本代理也重做了（证明那条通知是承重的）

`/tmp/wisp89bacc-hd`：`applyDescriptorWindows` 里 `noticeNarrowed(narrowNotice{Path: path, Principals: before})`
→ `_ = before // MUTATION-89BACC-SILENT`（`grep -n` 命中 144 行 = 锚点本行；变异前后 `go build ./internal/winsec/` **rc=0**）⇒

- `go test -count=1 -v -run 'TestSeal' ./internal/winsec/`：**`=== RUN` 2、两条全红**，
  `--- FAIL: TestSealReportsThePrincipalsItCleared`（`seal cleared a grant on data without reporting it; notices: []` +
  同一发对 `shared-with-a-service-account.txt` 的同类红）、
  `--- FAIL: TestSealNoticeIsRecordedByDefault`（`the default notifier wrote nothing auditable: ""`）
  ⇒ 与修复者自称"两条同时红"**逐字对上**。
- 全量 `-v` 跑同一 build：**RUN=29 / 17 顶层 PASS / 2 顶层 FAIL / 0 SKIP** ⇒ 变异**只**打到那两格，
  其余 17 条（含 §1 的覆盖面判据、AC#4/AC#5、生产端到端）保持绿。
- 还原：`grep -rn MUTATION-89BACC --include=*.go internal/winsec internal/secret` **0 命中**，
  `winsec_windows.go` 与 `migrate.go` 两个被改过的文件 `diff` 对 HEAD **IDENTICAL**，还原后 `go test ./internal/winsec/` **ok 4.024s**。

---

## 裁决表（退回单四条 1:1 + 两个顺带项）

| 项 | 修复者的主张 | 本代理裁决 | 本代理自己的依据 |
| --- | --- | --- | --- |
| **#1** AC#2 覆盖面判据搬进包内 | `c8d5c94` 上删 `propagatePrivate` ⇒ rc=1、唯一红名这一条、其余 14 绿 | **PASS** | §1.1 逐字复现（25 RUN/14 绿/1 红/0 SKIP + 红因原文）；§1.2 HEAD 上红两条（更强的承重证据） |
| **#1 附判** "纯继承的孩子不区分任何事" | 注释明写它不当判据 | **诚实**，但**票面 AC#2 第 (2) 句要降半档措辞** | §1.3：无 walk 的 build 里 `Everyone:(I)(RX)` 被 OS 重算掉、`Everyone:(RX)` 留着（本代理从变异日志取的原文） |
| **#2** 明文密钥备份接线 + repair 腿 | 同一 `.bak-plaintext` BEFORE 含 `S-1-1-0` → AFTER 只剩 SY/BA/swq；两条用例变异必红 | **PASS** | §2.1 本代理独立仪器的 BEFORE/AFTER 原文；§2.2 两发变异各红一条腿（他们的两条用例本代理也各自看到红）；`0a3a445` 上同一探针**红** ⇒ 上一轮那笔账也被本代理复现 |
| **#2 附判** 明文是否在别处留副本 | 未主张 | **没有残留（本代理扫过）**；但登记 R-89b-2/3/4 | §2.3：全树字节扫描除备份外 **0 命中**、`.migrate-tmp` 不存在；`internal/config/migrate.go:83`+`parse.go:213` 同类站点仍敞；`MigratePlaintext` **无生产调用方** |
| **#3** 目录符号链接真测掉 + A51② 再扩一格 | 未提权 + 开发者模式=1 ⇒ `os.Symlink` 成功；`Lstat` `Lrw-rw-rw-` + reparse 位；`os.Remove` 也能拆 | **PASS** | §第 3 条本代理自己的八行读数（`whoami /priv` 那一档本代理仪器空输出，**不计**，已在报告里点名） |
| **#4** 静默清除改成会报出来 | 只报显式、白名单不开口；两条用例必红 | **PASS**，附**射程限定 R-89b-5** | §4.1 攻出"继承来的那份不报"（0 条 WARN，实测）；§4.2 九发对照证白名单按 SID、SY/BA/我 全静默；§"第 4 条的变异"两条同时红、其余 17 条绿；§4.3 DENY ACE 与符号链接备份位两发都干净 |
| **回归** 四包 `-count=2` | "四包全 ok、全量输出 0 SKIP 0 FAIL" | **PASS**，读数更正一条 | 本代理：RC=0、402 RUN=2×201、FAIL 0、**SKIP 2 行 = 1 个具名 helper skip**；并测到"非 `-v` 输出根本看不见 SKIP"（R-89b-7） |
| **票 94 的牙** | "未新增任何绕开 `ResolvedPath` 的入口" | **PASS（没有绕开）** | 调用点逐个核 + `diff` 无新增 `filepath.Abs`/裸原语 + 本代理一发"未解析 Stat vs 已解析动作"的攻击被拒 |

---

## 三句直答

**(1) 票 89 现在能不能挂 `-done`？—— 能，但有三个"改名同时要做"的动作，都不需要再动生产码。**
四条退回全部由本代理独立复现（三条完全干净，第 1 条附带票面措辞降半档）。挂 `-done` 的同时请编排者落三笔：
① 票面 AC#2 第 (2) 句改成"边建边封的理由 = 不留可读窗口 + 修自带显式 ACE 的老孩子"（§1.3）；
② R-89b-5 那条通知射程限定 + R-89b-7 那条 SKIP 读数更正（票面自己的面，不是码的面）；
③ A51② 的 registry 更正（`:1865` 本体 + `:1453` 第二半 + 换掉 `RemoveUnlinked` 的理由句）。
`R-89b-2/3`（`internal/config` 同类明文站点、`MigratePlaintext` 无生产调用方）**不算 89 的欠账**，
但必须由编排者转成票 95/90 与装配根那张票的条目——否则会一直挂在"对 owner 的那句话"里说不满。

**(2) "私有数据只对我可读"今天能对 owner 说到哪一步（更新后的那一句）：**
"**这台机上，走 winsec 的那几条路——artifacts、DPAPI blob、`wisp.db`+`-wal`+`-shm`、staging、
迁移产出的 `config.toml.bak-plaintext`（含既存那份的 repair）与迁移临时件——落盘即只有
 `SYSTEM / Administrators / 我`；而且一处既存的老树，只要它自己 DACL 上带过显式的外来授权，
 下一次密封会被清掉并留一条 WARN**"。
仍然宽、不许混进这句话的：**`internal/config` 自己的写路径与它的 `config.toml.bak-<schema>` 备份**（票 95/90）、
**日志**（`0o644`，票 95）、**`models` 的 staging/缓存**（裁过：不按私有数据接）、`ball/position.go`、
`cmd/wisp/*`（禁改区）；以及两条结构性边界——**只继承自父目录的那份带外授权被清时不报**（R-89b-5）、
**`MigratePlaintext` 还没有生产调用方**（R-89b-3，所以那句话里的"迁移产出"今天只能读作"它一旦被产出"）。

**(3) A51② 要不要改：要，改成"本机不复现"并换掉存在理由。**
原文读数（`docs/reports/pending-and-issues.md:1865`）：
"**A51② `os.Remove` 删不掉"指向目录的符号链接"** ⇒ 一个裸名 artifact 若被替换成 symlink-to-dir，
票 79 的游离子树回收路径会清不掉它（它只用 `os.Remove`，绝不 `RemoveAll`，这个选择本身是对的…）"。
本代理与上一轮验收**各自独立**测到它对 **junction** 与 **目录符号链接** 两类对象**都不复现**
（本轮原文：`A51-RECHECK os.Remove(dir symlink to non-empty dir) err = <nil> ; target still there = true`；
`JUNCTION … os.Remove` 同样 `<nil>`、目标内容全在）。
⇒ 更正句：**"A51② 在本机 Go 1.27.1 / Win11 上不复现；`RemoveUnlinked` 继续保留，
理由从'删不掉'换成'只可能删到链接本身、绝不把删除半径交给别人'"**——
`A64③(1)`（`:1453`）只补了 junction 一半，另半边与理由句都还欠着。

---

## 本代理的破坏性检查（自证清白）

* 全部变异与仪器在 `%TEMP%` 下四个快照（`wisp89bacc-89b2` / `-hd` / `-pre` / `-c1`），**仓内未建 worktree**；
  工作树里 `internal/winsec`、`internal/secret`、`internal/memory`、`internal/agent`、`internal/risk`、`cmd/wisp`、`internal/perm`
  与本代理**零接触**（`git diff --stat -- internal/winsec internal/secret` 空；`git status --porcelain` 里
  只有别人的 `internal/risk/*`、`internal/tools/*` 与本报告）。
* 未 `git add -A`、未 `--amend`/`reset`/`rebase`/`stash`/`checkout .`、**未 push**；只 commit 了
  `docs/evidence/s1/89b-reacceptance.md` 这一个文件，且没往上一轮的 `89-adversarial-acceptance.md` 写一个字。
* 未在 `%APPDATA%` 的真实数据目录写过任何文件；符号链接/junction/DENY ACE/宽 ACL 全造在 `t.TempDir()` 内，
  测完由 `t.TempDir` 自清；快照目录本代理跑完即删。
* 变异 token 只存在于快照：`MUTATION-89BACC-WALK/-MIG/-REPAIR/-SILENT`；
  仓内 `grep -rn "MUTATION-89BACC" --include=*.go .` **0 命中**（终态复核）。
