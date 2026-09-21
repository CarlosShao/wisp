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

## 第 2 条：明文密钥备份 —— **接线本代理自己复现（独立仪器、同一对象前后 icacls 对照），另登记两条残留**

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
