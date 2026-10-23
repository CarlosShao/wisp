# 票 132 终裁腿 `132-v1` —— 五格逐格判：`132-r1` 的产品面读数**复核成立**；两格缺件（本件补齐）＋三格本件直接翻

- **被裁**：`132-r1` 四笔 `b22225513`／`6c381ab9c`／`191dfdf0b`／`ce3773d2a`（件 `132-r1-impl.md` 152 行＋`-part2.md` 74 行）。
- **起手锚**：`c68b4c35c`（本件期间共享 HEAD 推进到 `66e16de00`＝票 130 的三笔；见 §⑨ 写面现场）。两枚真二进制分别由 `git archive` 导出**仓外**树构建；新根全在 `D:\tmp\wisp\20261023\132v1\**`，⛔ 未碰 owner 真实数据目录。
- **写面**：本件＋仓外；**仓内零 `.go` 字节**。仓外探针一枚见 §④ 末与 §⑩。

## ① 九项逐项判语与读数

### 1) 票面前提新鲜度（对 HEAD blob 复量）
| 票面说 | 今天 | 尺 |
|---|---|---|
| `SealDir` 非测试调用点 **0** | **0**：非测试命中 1 行＝**定义行** `internal/winsec/winsec.go:301` | `grep -rn 'SealDir(' --include='*.go' cmd internal \| grep -v _test.go` |
| 测试引用 **16** | **16** | 同上再 `grep -c _test.go` |
| 定义在 `winsec.go:222` | **`:301`**（漂 79 行；追加裁决写的 `:301` 是对的） | 现读 |
| 洞在 `observe/logging.go:72`／`:244` | **改前 blob 逐字命中 `:72`（MkdirAll）／`:244`（OpenFile 0644）**；改后 `:99`／`:280`／封在 `:299` | `git show b22225513^:internal/observe/logging.go \| grep -n` |
| `logsink.go:71 const logDirName` | **`:76`**（`InitLog` 在 `:149`） | 现读 |
⇒ **前提新鲜**：标题那把尺今天仍读 0；行号漂移已在票面/追加裁决具名。

### 2) AC#1（现状表）—— **能翻**（条件：采信本件 §④；`132-r1` 件里缺 5/6）
- ★**它交的表只有 1/6 的分母**：两份件 `grep -n 'memory\.db\|config\.toml\|DPAPI\|backup\|artifacts'` ⇒ **0 命中**；表里只有 数据根／`logs\`／文件 三行。
- ★★**本件新发现（票面、腿都没写）——同一棵树，读数取决于走哪条腿**：`wisp slo`／`secret set` 不开 memory ⇒ 改前 `logs` **永久宽**；`wisp run` 在 `memory.Open`（`internal/memory/open.go:176` `PrivateDirAll(abs)`）**事后把整棵树收窄**（日志里 `12:04:06 WARN seal cleared principals path=<root>` 就是它）。⇒ 腿 §3 形状 B「根已存在且宽 ⇒ 保持宽」**只对 observe 那一次调用成立，对整个进程不成立**。
- ⇒ 缺什么：`config.toml`／DPAPI blob 与备份／`memory.db`／`artifacts`／`backup\` 与四枚副本的逐枚 `icacls`＋「已存在未锁／新建即锁」两形状；§④ 已补齐（DPAPI **备份**站点本机无 legacy 输入 ⇒ 未触发；静态读 `internal/secret/migrate.go:169` `PrivateFile`＋`:174` `SealFile` 两半都已封）。

### 3) AC#2（接线）—— **能翻**（条件：本件三档读数并入；另有一句假话要改）
- 三跳逐字命中：`:99` `winsec.PrivateDirAll(cfg.Dir, 0o700)`／`:280` `OpenFile(..., 0o600)`／`:299` `if err := sealLogFile(name); err != nil { f.Close(); return Wrap(...) }`（create 之后、**本进程首个内容字节之前**）。
- **① 同账户零影响：成立**（`secret set` 续写 `rc=0`、**1012→1562 字节／4→7 条**、ACL 仍窄、无侧车；`doctor` `[PASS] data dir writable (test)`）。
- **② 三档审计照写：腿的件里零读数**（`grep 'R20\|审计\|三档'` 只命中一个标题行）。**本件现跑**：`ask_every_step`／`ask_high_risk`／`auto_approve` 各写 **8 条** `audit:` 记录进**已封**文件（含 `audit: perm: MODE-READ … mode=<档>`）⇒ 性质成立。
- ⇒ 缺什么：② 的三档读数（本件 §④）＋ §⑤ 那句「封条在层出现之前就位」的改正。

### 4) AC#3（危害证明）—— **能翻**（腿只到 DACL 差，本件补了"别人真来写"）
- 腿只有 `icacls /findsid`（`SID Found`→`Not found`），**没有一次"别的主体来写"的尝试**（它自己明说）。
- 本件的"别的主体"：当前用户 token **把自身 SID 置 deny-only**、保留 `Authenticated Users`／`Users`／`Everyone` ⇒ 等价于"同机另一个已认证标准账户"，然后**真发 `CreateFile`**。
- 读数：宽文件 **APPEND_ALLOWED**（1023→1024）／READ_ALLOWED／LIST_ALLOWED；窄文件 **APPEND／READ／LIST 全 DENIED**；把窄描述符搬到**宽目录**里的副本仍 **APPEND／READ DENIED**（目录 LIST_ALLOWED）⇒ 文件那一半独立生效，不靠目录代偿。

### 5) AC#4（变异自证）—— **能翻**
- 本件在仓外导出树自造 M0–M3（**⛔ 未参考它的红名**）：红集合 **6／4／3／1**，与它逐名相同（§⑥）。
- **两半不代偿**：M1 后 **PASS=2**（文件两枚仍绿）、M2 后 **PASS=3**（目录三枚仍绿）。
- 枚举门：AST、跟随别名、**只认调用**；钉 `internal/observe` 必须同时有 dir-sealer 与 file-sealer 两类**调用**（能力口径）。⇒ 与 §③ 裁定一致；但它**不钉"`SealDir` 这个导出名零调用者"**（那枚今天仍是死导出，另行登记）。

### 6) AC#5（门禁四数）—— 见 §⑦（rc=1；四数在册；新增红 0 见 §⑧）
### 7) 门禁 cp —— 全绿：`build 0`／`d22scan 0`／长度门 0／格式甲乙两尺都 0（§⑦）
### 8) 新增红 0／写面纪律 —— §⑧／§⑨
### 9) 另报两件 —— §⑩

## ② 五格逐格

| 格 | 判 | 缺什么／条件 |
|---|---|---|
| **AC#1** 现状表 | **能翻** | `132-r1` 件缺 5/6 具名项读数＋按腿分栏；本件 §④ 的表即该格交付物，采信即可翻 |
| **AC#2** 接线 | **能翻** | ② 三档读数在腿件里没有（本件 §④ 补）；腿 §2「封条在层出现之前就位」是假话，须改成 §⑤ 的写法 |
| **AC#3** 危害 | **能翻** | 腿件只有 SID 差；本件 §④ 的"别的主体真写"补上 |
| **AC#4** 变异 | **能翻** | 无（枚举门是能力口径，`SealDir` 死导出另行登记） |
| **AC#5** 门禁 | **能翻** | 无（四数在册＋各门 rc=0＋新增红 0；已知红 cmd 7 名／winsec 2 名，逐名同对照树） |

## ③ 编排者「按能力判、⛔ 按函数名判」——**结论成立，理由有一处不成立**

- ★**结论对**：`PrivateDirAll` 封**它创建的每一层**并在建下一层前封好上一层；且 Windows 的 `sealDir`＝`applyDescriptor`＋`propagatePrivate`（`winsec_windows.go:673-678`），把**已在里面的孩子**一并收窄。⇒ 容量确实交付，`SealDir` 计数只是代理指标；为让 grep 非 0 再加一枚冗余调用是"给尺表演"，该禁。
- ⛔**理由错**：它引的 `privateDirAll` 文档（`winsec.go:257-259`）逐字讲的是 **C26**（"no path reaching os.Mkdir or **SealDir** here can have skipped C26"），**不是顺序**；代码是 `os.Mkdir`（`:283`）**然后** `sealDir`（`:292`）——**"在每层出现之前封"字面为假**。真正的强保证是导出文档 `:223-229` 那句"sealing **every level it creates** before descending"＋`propagatePrivate`，与"父先封、再建子（孩子出生即窄）"。
- ⇒ **票面标题那把尺不必非 0**（AC#2 不因此 ⛔）；但票面应把这把尺改写成能力句，否则下一位读者还会按名字重开一次。`SealDir` 导出今天无人调，另行登记。

## ④ 我自己造的产品面 ACL 读数（仓外新根、两枚真二进制、`icacls` 逐枚）

| 对象 | 改前（`wisp slo`／`secret set` 腿） | 改后 |
|---|---|---|
| 0 字节日志文件（首启） | `Administrators:(I)(F)`／`SYSTEM:(I)(F)`／**`Authenticated Users:(I)(M)`**／**`Users:(I)(RX)`** | `SYSTEM:(F)`／`Administrators:(F)`／`swq:(F)`（**全显式**） |
| `logs\` | 全 `(I)` 宽（同上四条） | 六条显式（含 `(OI)(CI)(IO)`） |
| 数据根（本调用创建） | 全 `(I)` 宽 | 六条显式 |
| `/findsid *S-1-5-11` | **`SID Found`** | **`No files with a matching SID was found`** |
| `/findsid *S-1-5-32-545` | **`SID Found`** | **`No files with a matching SID was found`** |

**AC#1 的 12 个落点（`wisp run "hello"` 全路径，改后树，逐枚 `icacls`）**：数据根／`logs\`／`wisp-…-001.jsonl`／`config.toml`／`secrets\`／DPAPI blob／`wisp.db`／`artifacts\`／`backup\`／`bak-0-1`…`bak-3-4` —— **全部窄**。机制分三种：**显式**（根、logs、日志文件、config、blob、artifacts、backup 目录，及后续 `PrivateDirAll` 的 propagate 过的 wisp.db 与备份）／**继承 `(I)`**（首跑的 `wisp.db` 自已封根继承；末枚备份只剩 `(I)`）。**改前同一路径**：`wisp run` 终态也窄，但那是在 `memory.Open` 事后收窄（日志 `12:04:06`），宽窗口里早期记录已落盘；`wisp slo`／`secret set` 腿则**永久宽**。

**形状 B（已存在且宽 ⇒ 修）产品面**：预建宽根＋`icacls /grant:r *S-1-5-11:(OI)(CI)(M)` 的 `logs\` ⇒ 跑一发后 `logs\` **变成六条显式窄**、文件三条显式、**数据根仍宽**（observe 不碰 `path` 以上祖先）——与腿 §3 形状 B 一致（限"不走 memory"的腿）。

**滚动文件**（腿没测"每个文件"）：仓外探针用包内 `newRollingWriterClock`（1 MB 阈值）造出 `001`＋`002` 两枚，目录是**未封的 `t.TempDir()`** ⇒ 改后两枚都带**显式** ACE；同一探针在 HEAD-去本笔的对照树上 **两枚都 `explicit=[]`**（窄全靠继承＝"测空气"的形）⇒ "滚动出来的每个文件都封"成立。

## ⑤ 「创建时」字面 —— **不是"层出现之前"，是"建完立刻、且在进下一层之前"**

- 代码事实：`privateDirAll` 对每层 `os.Mkdir`（`winsec.go:283`）**后**才 `sealDir`（`:292`）；最后再封一次 `abs`（`:296`）。⇒ **mkdir 在前、封在后**，⛔ 不是"出生即封"。
- 但**效果上不留宽孩子**：层 k 建在"已封的 k-1"下面 ⇒ 出生即窄；层 k 的 mkdir→seal 窗口里被别人塞进来的孩子，由 `sealDir` 的 `propagatePrivate` 收窄；seal 之后建的孩子继承窄（实测：在已封 `logs\` 里新建一枚文件 ⇒ `SYSTEM:(I)(F)`／`Administrators:(I)(F)`／`swq:(I)(F)`）。
- ⇒ AC#2 的字面「创建时**就** SealDir」**成立**（创建与封同一次调用，且在任何日志文件打开之前）；⛔ "在每层出现之前"的说法要改。若要真正的"出生即封"，得给目录也用 `CreateDirectoryW`＋`SECURITY_ATTRIBUTES`（今天是先建后 `SetNamedSecurityInfo`）——**那不在本票字面要求内**。

## ⑥ 反向判据：本件自造变异（仓外导出树，`sha256=13d8180e…` 为基）

| 突变 | 摘掉哪跳 | 结果 | 红名（逐名） |
|---|---|---|---|
| M0 | 两半全摘 | build 0；**6 红**（PASS=0） | 全部六枚 |
| M1 | 只摘目录封（`PrivateDirAll`→`MkdirAll`） | **4 红，PASS=2** | `…LogDirIsSealedAtCreation`／`…ExistingWideLogTreeIsRepaired`／`…UnsealableLogDirIsRefusedLoudly`／`…SealEntryPointsHaveProductionCallers`（**文件两枚仍绿**） |
| M2 | 只摘文件封（删 `:299` 整块） | **3 红，PASS=3** | `…LogFileIsSealedBeforeTheFirstContentByte`／`…SealFailureRefusesTheWrite`／`…SealEntryPointsHaveProductionCallers`（**目录三枚仍绿**） |
| M3 | 留调用吞错误（`_ = sealLogFile(name)`） | **1 红** | `…SealFailureRefusesTheWrite` |

⇒ **两半互不代偿**（M1/M2 的 PASS 集合恰好互补）；每次突变后 `build rc=0`、还原 `sha256=13d8180e…` 后原样重跑 **6 PASS**。

## ⑦ 门禁 rc（现跑）

| 门 | rc | 读数 |
|---|---|---|
| `go build ./...` | **0** | — |
| `go vet ./internal/observe/ ./internal/winsec/` | **0** | — |
| `sh scripts/d22scan.sh` | **0** | `clean`；`bans #1-5 internal/=248 cmd/=42`；`ban #8 internal/=631 cmd/=148`（＝腿报的 290 production files 对得上） |
| `sh scripts/check-path-length-budget.sh` | **0** | `VERDICT GREEN`；`over=61／roster=61`；bands 57/4/0 |
| `gofmt -l`（5 枚） | **0** | 空 |
| `gofumpt -l`（5 枚） | **0** | 空；正控（导出树放一枚未格式化）**点名**且输出是**反斜杠**形，`[\\/]` 归属命中 **1** |
| CRLF（`tr -cd "\r" \| wc -c`） | — | 5 枚全 **0** |
| `go test -count=2 -v ./cmd/wisp/ ./internal/winsec/` | **1** | 干净导出树 `head-now`（`66e16de00`，树内自带三枚 DLL）：**RUN=755／PASS_ALL=737（TOP 500）／FAIL_ALL=12（TOP 10）／SKIP=5**；cmd/wisp `RUN=509／PASS=496／FAIL=7／SKIP=5`（整包 **600.264s 撞 10 分钟上限**，goroutine dump，无终判）；winsec `RUN=246／PASS=241／FAIL=5（TOP 3）／SKIP=0` |

**红名（逐名）**。cmd 7 枚：`TestCleanCheckoutBuilds_AC11`／`TestAC14GoSideEvalPushReachesThePage`／`Test258BridgeMutationNoSrcKeepsOldBinding`／`Test258OccupiedCombinationNamesTheNewValue`／`Test258V1ProbeSummonEditRebindsLiveBall`／`TestAC1ResidentLegInstallsItsLogListenerOnDisk`／`TestAC246DevLegIgnoresTheTestTaskInjection`。winsec 2 名：`TestAC1WindowsPrivateFileThroughASwappedAncestorNeverTruncatesAForeignFile342r1`（＋`/junction`／`/symlink-dir`；`-count=2` 跑了两轮）、`TestAC2WindowsLeafIdentityReCheckRefusesAForeignNameAndNeverDeletesIt342r1`。

**逐名归因（都不是本笔）**：`…CleanCheckoutBuilds_AC11`＝`git rev-parse --show-toplevel: exit status 128`（`git archive` 树无 `.git`）；`…AC246Dev…` 与 `…ResidentLegInstallsItsLogListenerOnDisk`＝**本机有一条 dev Wisp 常驻在跑**（`another instance is running in this session`／`the child never reached the event loop`，单实例把子进程"激活并退出"；该用例仍量到 `1 file(s), 7 record(s)` 落在已封 `logs` 里）；`Test258*`／`AC14GoSideEval`＝面板/WebView2 现场；winsec 两枚＝342r1 的 400 次竞态仪器（且 `internal/winsec` 不依赖 observe）。**同一命令三发三组四数**（首发主检出 `RUN=467／PASS=462／FAIL=4／SKIP=0`，cmd 零具名红但在 `TestAC13BringUpSurvivesAReusedThreadQuit` 处截断）⇒ 四数只能当形状、判据只按名。

## ⑧ 新增红 0

- **winsec**：`go list -deps ./internal/winsec` 对 `internal/observe` 命中 **0** ⇒ 本笔改不到它；两枚 342r1 红是**竞态仪器**（红名逐发漂，腿也这么说）。
- **cmd/wisp**：本发 **7 枚具名红**（§⑦ 逐名），**逐名与对照树相同** ⇒ 不是本笔新增；整包 FAIL 的另一半是 **600s 超时**（无终判），与"测试判红"是两件事。
- **对照尺**：`HEAD`（`66e16de00`，含票 130 三笔）导出仓外，**只把本笔的三处改回去**（`PrivateDirAll`→`MkdirAll`、删 `:299` 封块、`0o600`→`0o644`）并去掉四枚本笔测试；两棵树各带 `third_party/sherpa-onnx` 的三枚 DLL。第一次的对照树（把 `logging.go` 还原成 `b22225513^` 的 blob）同时倒退了票 130，量出 19 枚 `no native DLLs` 红 ⇒ **作废重做**，本节是重做的那一发。
- **读数**：`head-m0`（`-count=1`）cmd `RUN=479／PASS=467／FAIL=7／SKIP=5`、winsec `RUN=123／PASS=120／FAIL=4（TOP 2）／SKIP=0`；**cmd 的 7 名与 `head-now` 逐名相同、winsec 顶层两名相同**（`-count=1` 只跑一轮，故 342r1 少一枚顶层）。
- **结论**：**新增红 0**（逐名差集为空；两棵树的多名红都是本机环境：无 `.git` 的导出树／常驻 dev Wisp／WebView2 面板／342r1 竞态）。

## ⑨ 提交号＋逐笔越界（`git show --name-status` 逐笔）

| 提交 | 文件 | 越界 |
|---|---|---|
| `b22225513` | `M internal/observe/logging.go`；`A` 四枚 `internal/observe/*132*_test.go` | **0**（全在 `internal/observe/**`；无 `cmd/**`、无 docs、无 scripts） |
| `6c381ab9c` | `A docs/evidence/s1/132-r1-impl{,-part2}.md` | **0** |
| `191dfdf0b` | `M` 同上两枚 | **0** |
| `ce3773d2a` | `M docs/evidence/s1/132-r1-impl.md` | **0** |
| 正控 `6e19f4a43`（别人一笔） | 同尺列出 **7 枚 `.go`** ⇒ 尺看得见 `.go`；`96f1af9c9`＝**0 枚**（空正控，腿自报并换掉，本件复量同） | — |

- `internal/winsec` **一字未改**（`git diff b22225513^ b22225513 -- internal/winsec` 空）。本笔自己的提交号见件尾。
- ★**写面现场**：本件期间共享 HEAD 从 `c68b4c35c` 推进到 `66e16de00`（票 130 三笔，改的也是 `internal/observe/logging.go`）。⇒ 本件第一发 AC#5 门禁（`11:55:45–12:01:00`）跑在**正在被 130 改写的工作树**上；为准，§⑦ 的数字改用**导出树 `head-now`（`git archive HEAD`＋DLL）**重跑，并在 §⑧ 用同一棵树做对照。

## ⑩ 它另报的两件（本件判）

1. **`cmd/wisp/logsink.go:44-46` 的 `INTERIM (ticket 117 AC#3)` 已过期 —— 成立**。它写 "nothing here seals or un-seals the log file"，又写"日志算不算私有数据是 owner 没答的那一问"；owner 09-23 已批 `Q-31`「算、要锁」，且管线现在真的封。该文件在 `cmd/**` ＝它禁写，它**如实拒绝动**，对。⇒ 排队给下一波 `cmd/wisp` 写腿。
2. **运行中途 roll 封失败 ⇒ `slog` 吞错 —— 成立**。代码：`:299` 失败即 `f.Close()`＋返回错；`redactHandler.Handle` 原样 `return h.next.Handle(...)`（`:176-182`），`slog.Logger` 不向调用方上报 Handler 错误 ⇒ 记录**消失**、产品侧**没有一行"seal 失败"**；启动时那一半仍响（`run`／`resident` 打印 `持久日志未启用（…）…不会落盘`，`cmd/wisp/run.go:222`／`resident_windows.go:114`）。"不许没锁但能写"这条是关的，"失败有一行名字"这条**没有**。⇒ 编排者把它放"残余、另开票"的处置**对**；建议单开一枚"Handler 失败必须有一行"的票。
