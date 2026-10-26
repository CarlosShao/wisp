# 票 347 落地腿 `347-r1` — 形**一**（＝`AC#1` **甲**）：POSIX 侧 C26 的 `realpath` ＋ `lstat` 拒 symlink **真落地**

- 锚＝`HEAD ba0f4597d`（2026-10-26 19:4x +08）。起手核：`git diff --cached --name-status` **空**；`git status --porcelain -- internal/risk internal/tools` **空**（`cmd/wisp` 有 `wire-r1` 在飞件 6 枚＝他腿，本腿**零触碰**，见 §3 末）。
- 本腿射程＝编排者裁的「一」，**零契约字面改动**（`C26`＝`PLAN.md:1382`、`SPEC-06 §4`＝`:47-58`、`D7`＝`:184-199` 一字未动；`PLAN.md:2385` 那句"**macOS 侧**"**未改**）。
- 读数全部是**我自己现跑**的：包级 `go test -v ./internal/risk/`（Windows 本机 ＋ Linux `golang:1.27` 容器）、门禁级 `bash scripts/portable-tests.sh --scope={core,windows}`、`d22scan`、`gofmt`/`gofumpt`。日志全在 `D:\tmp\wisp\20261026\347r1\*.log`（只建不删）。

## 1. 票面八格逐条（原文 → 判定 → 读数）

| 格 | 票面原文（摘） | 判定 | 读数 |
|---|---|---|---|
| **AC#0** 前置 | 「编排者把一／二／三摆给机主…之后才许开工」＋起手「五格逐枚重跑」「`DEFERRED(macOS/Linux)` ⇄ `SPEC-12 §5` 现读对撞」 | **成立**（编排者已裁「一」，本腿不自己选形） | ① 五格重跑：桩本体 `pathresolver_other.go:10/:12` **两枚**桩；`f.root` 唯一赋值点 `pathresolver.go:354-357`；守卫 `blacklist.go:363`（票面写 `:359-361`＝**偏 4 行，票面错**）；`portable-tests.sh` 档位＝`core`(scope `:347`)/`windows`(`:359`)、逐名 `core_pin :291`/`win_pin :308`（票面写 `:266/:282`＝**两处都错**）。② 标记对撞：删掉 `DEFERRED(macOS/Linux)` 后 `internal/risk` 仍有三枚 `DEFERRED(...)` 标记（`doc.go:15`、`provenance.go:58`、`syncdirs_other.go:21` `DEFERRED(P12-macos)`）⇒ `SPEC-12 §5:66`「`DEFERRED｜macOS 平台层`」那行**仍有 1:1 对面**，**不必新增登记行**。 |
| **AC#1** 三形代价表 | 「**⛔ 未出表不许动码**」；甲＝`resolveHandle` 落 `realpath`＋`lstat` 拒 symlink | **成立**：表已由 `347-a1` 出（`docs/evidence/s1/347-a1-census.md`），本腿＝甲，并按普查 §6 把三条射程补进落地（§2） | 甲落地后的包级实测＝§4 的 shot B（4 新红，逐名对得上普查那 5 枚）＋ shot C（全绿） |
| **AC#2** 单独摘守卫 | 「⛔ 只把 `blacklist.go:359-361` 那枚首行守卫摘掉…⇒ `isUnder("","")` 恒真会把 POSIX 反向全量 fail-closed 成 A」；「去牙那一发必须点名这一形」 | **成立，但落地后这一形结构上被退役**（本腿**具名**它，⛔ 掩盖） | 突变 **M2**（摘掉 `uncertainAnchorMiss` 首行 `cand.root == "" \|\| anchor.root == ""`）＝**rc=0 / 0 FAIL**。原因现量：`realpath` 落地后 `formsOf` 的 root **最短到 `/`**（`/x` 的父链一路 Lstat 到 `/`），而 `isUnder(x,"/")` 要前缀 `"//"` ⇒ 两侧 root **恒非空且不构成祖先关系** ⇒ "摘守卫＝全 A"**不再成立**。反过来说：`AC#2` 要防的"缺世界退化成全 A"**没有发生**，证据＝§5 的 M1/M3/M4 都是**定向红**、不是全量红。 |
| **AC#3** 搬牙不删牙 | 「两枚 336 仪器＋339 绊线＝**设计好的"搬我"信号**⇒ 新台面下要改其射程，⛔ 删除、⛔ 让它静默绿」；「合件时必须**保名** `candidateOffAnchorChain334`」 | **成立**：336 的 8 枚牙**整份 roster 搬到 untagged**、两平台各带自己的 partial 制造器；339 的两枚绊线**换成 reason/class 两枚真牙**；`candidateOffAnchorChain334` 保名（搬到 `anchor_world_334_test.go`） | `TestAnchorFallbackTeethPerAnchor336` 在 Linux **首次跑起来**（8 行＋8 个 offchain 子测试全绿）；`TestAnchorMissReason…Posix339`／`…ClassPosix339` 两枚真牙。**⛔ 未用"搬牙即绿"交差**：搬之前先量了 shot B（5 红），搬之后才绿（§4）。 |
| **AC#4** 颜色与门禁 | 两平台成对；四数逐枚；红名册**逐名作差**；「门禁四数只在共享工作树、★跑前 `unset GOTMPDIR TMPDIR`」；⛔ `./cmd/wisp` 整包、⛔ `go test ./...` | **成立**（两平台都给了；门禁级另给**基线对照**，见 §4） | 包级：Windows `147/1/1→147/1/1`、Linux `122/1/1→127/1/0`。门禁级：Linux `core` 基线 `1388/10/1` ⇄ 本腿 `1392/11/1`，**集合差＝+1 枚 config**（§4.3）；Windows `windows` 档 `525/2/0`（两枚均基线既有，§4.4）。`cmd/wisp` **零字节**（§8）。 |
| **AC#5** 零新导出名／零契约移动 | ⛔ 新导出名／新 JSON 键／`panel.*`；⛔ 动 `C26` 入口语义、档位、豁免键、`FILE_ATTRIBUTE_REPARSE_POINT` 半；冻结件与 `thresholds.go`/golden/`go.mod`/`go.sum`/D/C/R 表一字不动 | **成立** | `git diff` 里**零新导出标识符**（新增全部是 `!windows` 文件内的小写函数与 test 内符号）；`pathresolver_windows.go` **零字节**；`C26` 的输入/输出类型与档位不变（`Result` 结构零改动）。**唯一语义变化**见 §2-④（不是契约移动，见那一条）。 |
| **AC#6** 纪律格 | ⛔ 翻框／改票名／`Status:`／`-done`；写面＝`internal/risk/**`＋`.scratch/wisp/probes/347/r1/**`＋本票末尾**一条** bullet；`add`/`commit` 两处 pathspec 逐枚手写；⛔ push | **成立，一处按本腿派单收窄** | 写面实际＝`internal/risk/**`＋`scripts/portable-tests.sh`＋`docs/evidence/s1/347-r1-impl.md`（派单逐字给的三处）；**⛔ 未动票面任何一字节**（含那条 bullet——派单禁改票面，二者冲突时以派单的"⛔ 票面"为准，请编排者决定是否要一条 appended bullet）。`add`/`commit` 逐枚字面点名（§10）；**⛔ 未 push**。 |
| **禁区** | ⛔ 与票 55 同波动 `pathresolver*.go`；⛔ 与票 342 混记；⛔ 回头改票 339 判据文字 | **成立** | 票 55 未动；票 342（`internal/winsec`）零字节；票 339 票面未动。 |

## 2. 三条射程补充：逐条落地

1. **`reparseComponents` 也落**（普查 §6-1）→ **落了**：`internal/risk/pathresolver_other.go:88` `reparseComponents` 自 `/` 逐组件 `os.Lstat`，`ModeSymlink` 命中即收；首枚缺失/不可 Lstat 即停（与 Windows 版 `:74-86` 同形，返回**拼写照**的组件＝例外表与审计行用的那个拼写）。⇒ 只落 `realpath` 会造成"软链被跟随"，本腿没走那条路。**证明它载荷**＝突变 **M4**（把收集行去掉）：两条 symlink 用例立刻红，mutant 读数逐字 `Resolved=true … Canonical=…/outside/secret.txt`（**跟随到界外**＝正是 fail-open 形）。
2. **`normalizeLocalUNC` 的 POSIX 唤醒**（普查 §6-2）→ **落了编译期早退**：`pathresolver.go:251-263` `if filepath.Separator != '\\' { return p }`。**普查那半句我现核后修正**：经 `Resolve` 的 `//` 起头路径**今天已经被 `lexCanonical`(`Clean`) 折掉**，所以走管道**不会**被换成 `\`；能唤醒它的是**直接调用**。两条都测了：`TestNormalizeLocalUNC…347` 直接调（`//localhost/c$/x`、`//server/share/y`、`\\server\share\z` 全部原样）；同文件再跑管道形（`Resolve("//localhost/c$/x")` 不带 `\`，`TestDoubleSlashSpellingIsNotUNC` 仍绿）。**突变 M5**（摘掉早退）⇒ 红，mutant 读数 `normalizeLocalUNC("//localhost/c$/x") = "C:\x"`（比普查写的"整条换 `\`"更糟：它命中了 localhost/`c$` 那支，变成**Windows 盘符路径**）。
3. **"搬牙"实为"重造牙"**（普查 §6-3）→ **换掉了 partial 制造器**：POSIX 上的 `U+001B`／`U+202E` 是**合法文件名字符**，所以 `hostile334` 在 POSIX 只会是"普通缺失名"⇒ `partial=false`、clause 不响（shot B 日志逐字：`class=B … clause=false`）。新的制造器＝**dangling symlink**（`anchor_fallback_teeth_336_other_test.go:57`、`anchor_miss_reason_339_other_test.go:46`）：`Lstat` 成功（**存在**，正是 `pathForms.partial` 说的那半）＋`Stat/EvalSymlinks` 打不开（目标不在）⇒ 与 Windows "名字被拒解析"给出**同一对事实**，且**不依赖 root 权限**（容器里 000 目录会被 root 绕过，所以没选它）。制造器自己带断言（`Stat` 必须失败、`Lstat` 必须成功），⛔ 假设。

**④ 本腿在射程补充之外**顺带测出并修掉的一枚同类缺陷（⛔ 当"顺手改"）：`classifyForms` 原先把 `normPath` 的**比较折叠形**（小写化）喂给 `anchorForms`→`formsOf`，而后者要 **Lstat** 它。Windows 上无所谓（大小写不敏感），POSIX 上 `$HOME=/home/Alice` 折成 `/home/alice` ＝**另一棵树**，锚侧 root 停在错误祖先（实测：`anchor.root="/tmp"` ⇒ off-chain 对照被误判成 A）。修法＝锚侧改用**环境原拼写**（`blacklist.go:203-217` + 六个 `anchorPath(homeRaw,…)`），折叠加在 `formsOf` 内部（`raw` 仍是折叠形，比较语义不变）。这正是该文件 `normPath` 自己的注释所要求的"折叠形**不得**交给 OS"。Windows 读数证明零回归（§4.1 三数全同）。

## 3. 产码面逐文件（"搬牙与产码**同一提交**"的划分）

**产码（5 件，2 件有可执行字节＋3 件注释；`git diff --numstat` 现量）**
| 文件 | 改动 | +/− |
|---|---|---|
| `internal/risk/pathresolver_other.go` | 删两枚桩＋`DEFERRED(macOS/Linux)`；落 `resolveHandle`(`EvalSymlinks`＋空/错守卫)＋`reparseComponents`(lstat 走) | **+87/−8** |
| `internal/risk/pathresolver.go` | `normalizeLocalUNC` 编译期早退；`pathForms.real` 注释 | **+16/−2** |
| `internal/risk/blacklist.go` | 锚侧原拼写（上面的 ④）＋两处陈旧注释 | **+39/−24** |
| `internal/risk/syncdirs.go` | 注释（`pathHandleVerify` 残留原因） | **+4/−2** |
| `internal/risk/syncdirs_other.go` | 注释（同上，`pathHandleVerify` 仍 `false`） | **+9/−5** |

**搬牙／重造牙（10 件改动 ＋ 3 件新文件 ＋ 1 件脚本，**与产码同一提交**）**
`anchor_world_334_test.go`(新 107 行,untagged 世界唯一家) · `anchor_fallback_teeth_336_roster_test.go`(新 141 行,8 枚 roster＋断言 untagged) · `anchor_fallback_teeth_336_other_test.go`(+60/−75,posix 制造器) · `anchor_fallback_teeth_336_windows_test.go`(+14/−126,缩到 windows 制造器) · `anchor_miss_reason_339_other_test.go`(+120/−123,reason/class 两枚真牙) · `anchor_miss_reason_334_test.go`(+8/−78,删搬走的夹具块) · `pathresolver_symlink_other_test.go`(新 139 行,§6 正向＋UNC 直接调用) · `syncdirs_other_test.go`(+47/−77,绊线→grade 因果＋前提翻转) · `syncdirs_test.go`(+65/−11,grade 用例搬入 untagged＋弱 grade 的 canonical 前提) · `syncdirs_windows_test.go`(+21/−48,donor 移出＋头注) · 两处注释件(`provenance_syncdirs_{other,windows}_test.go`、`syncdirs_ancestor_actable_leg_116_test.go`) · `scripts/portable-tests.sh`(+1/−2,账本：删一行、重写一行)。
★**比普查估的"7 件≈300-460 行"多**：实为 **18 枚文件**（387 行新文件＋约 500 行改动）。多出来的量全在"重造牙"（普查 §6-3 已预警）与两平台夹具合流上，⛔ 是跑偏。

**为什么必须同一提交**：①只回退产码 ⇒ 搬过来的牙**全红**（突变 M3 就是那一发：`resolveHandle` 回桩 ⇒ 336＋339 三枚红）；②只回退牙 ⇒ 两枚绊线**立刻红**（shot B：`119/5/0`）。两个方向都红 ⇒ 单独回退任何一半都不成立，故"可单独回退"的前提＝两者同提交。`cmd/wisp` 6 枚在飞件属 `wire-r1`，本腿 `git diff` 里**零命中**（§8）。

## 4. 两平台分母（包级 ＋ **门禁级**）

### 4.1 包级 `./internal/risk`（同一 `-count=1`，逐枚点名，⛔ 相减）
| 平台 | 基线 HEAD | shot B（只落产码） | shot C（终态） |
|---|---|---|---|
| Windows（本机） | **147 PASS / 1 FAIL / 1 SKIP**（149） | —（`!windows` 改动，不适用） | **147 / 1 / 1**（149）**一数不变**；唯一改名 `…DisarmFallbackWindows→…DisarmFallback` |
| Linux（容器） | **122 / 1 / 1**（124） | **119 / 5 / 0**（124） | **127 / 1 / 0**（128） |

Linux 名册差（`124−3+7`）：删 `TestFallbackBranchHasNoWorldOnThisPlatform336`／`TestAnchorMissHasNoWorldOnThisPlatform339`／`TestSyncNoGradeIsConfirmedOnPosix`；加 `TestAnchorFallbackTeethPerAnchor336`／`TestAnchorMissReasonNamesCandidateAndMissingAnchorPosix339`／`TestAnchorMissKeepsNonOverridableClassPosix339`／`TestResolveDeniesSymlinkTraversalByDefault347`／`TestResolveFollowsExemptedSymlinkOntoRealPath347`／`TestNormalizeLocalUNCLeavesPosixSpellingsAlone347`＋搬入的 `TestSyncConfirmedGradesDisarmFallback`。Windows 的 `149` 里含**同一枚改名**，所以两平台名册不能直接相减。

### 4.2 中间那 5 红逐名 ＋ 守卫行（shot B，`D:\tmp\…\linux-shotB.log`）
| # | 红名 | 守卫行 → 报错行 | 性质 |
|---|---|---|---|
| 1 | `TestFallbackBranchHasNoWorldOnThisPlatform336` | `anchor_fallback_teeth_336_other_test.go:75`（`if f.root != ""`）→ `:76`，4 个 relCand 各一条 | 绊线（设计即"搬我"） |
| 2 | `TestAnchorMissHasNoWorldOnThisPlatform339` | `…339_other_test.go:139`（`if f.root != ""`）→`:140`；`:142`（`if a.root != ""`）→`:143` | 绊线 |
| 3 | `TestSyncNoGradeIsConfirmedOnPosix` | `syncdirs_other_test.go:63`（`Resolve(root,nil)…|| res.Resolved`）→`:64`；另 `:94`→`:95`、`:99`→`:100` | 绊线（文件头 `:37-38` 自己写"delete"） |
| 4 | `TestSyncMembershipDecidesOffProfilePosix` | `syncdirs_other_test.go:94`（`SyncDetectionComplete()`）→`:95`；`:99`（under-profile 必须 suspect）→`:100` | 前提随世界翻转 |
| 5 | `TestC26RewriteAccountIsConsumedAtEverySecurityLeg` | `pathresolver_rewrite_account_test.go:200`（报 `internal/panel/composer.go:127`＋`display_scrub_216.go:26/37/218` 直读 `Result.Canonical`） | **存量 `A380`**，HEAD 两平台都红（§7） |

⇒ 前 4 枚已由搬牙**全转绿**（shot C）；第 5 枚**不是本腿的红**。

### 4.3 **门禁级**（`scripts/portable-tests.sh --scope=core`，同一容器、同一 git-less 导出树、同 `unset GOTMPDIR TMPDIR`）
| 树 | 四数 | rc | 红集 |
|---|---|---|---|
| **基线**（`HEAD` 逐文件 `git show` 还原，16 枚 `hash-object` 全等） | `PASS=1388 FAIL=10 SKIP=1`（`=== RUN=2030`） | **1** | memory 3（`NothingWritten…`／`HandDeleted…`／`Dsn196r3…`）＋risk 1（第 5 枚）＋panel 5（`ApprovalCard…`／`StreamLogFlood…`／`StreamLogDropped…`／`ComposerContract…`／`ComposerRender…`／`ReadGit…` 中的 5 枚） |
| **本腿** | `PASS=1392 FAIL=11 SKIP=1`（`=== RUN=2052`） | **1** | 上面 10 枚**逐名全同**＋**新增 1 枚** |

**集合差＝+1 枚，且只有 1 枚**：`TestAC2POSIXSeamInstallSurvivesASymlinkSpelledTemp125`（`internal/config/c26_seam_posix_125_test.go:117`）。机制现量：它把子进程 `TMPDIR` 设成穿过 symlink 的 temp 根，钉子调 `winsec.PathResolverInstalled().Resolve(base+"/../wisp-125-probe")` 并要求**非错**；`reparseComponents` 落地后该路径祖先 `varlink125` 是 symlink ⇒ 正确返回 `ErrReparseDenied`。**我判这是契约正确的后果**（`SPEC-06 §4` 默认拒；Windows 那半同样无条件拒），而**该文件在我的禁写面**（`internal/config/**`）⇒ 已具名上报编排者（附三条修法）。`SKIP=1` 两发同枚（`TestAppendSessionTurnCrashChild`）＝既有。**PASS 差闭合**：`1388+5(risk 净增)+(−1)(config)=1392`。

### 4.4 Windows 门禁级（`--scope=windows`，共享工作树，`unset GOTMPDIR`）
- **基线**（`export-gate-base` 同档）：`PASS=526 FAIL=1 SKIP=0` rc=1，红＝`TestC26RewriteAccountIsConsumedAtEverySecurityLeg`。
- **本腿**：`PASS=526 FAIL=1 SKIP=0` rc=1，红＝**同一枚** ⇒ **Windows 门禁级零变化**。
- 第一发曾读到 `525/2/0`（多一枚 `TestResolvePerCallBudget`）：那是**与 Linux `core` 档并发跑**的产物；单独重跑第二发回到 `526/1/0`。该用例是**时序**（预算 1.000 ms/op），`0.593 ms/op` 孤立绿，且本腿改动在 `!windows` 文件里、**不编译进它的路径** ⇒ ⛔ 当新红。

## 5. 反向判据：突变三读数 ＋ 还原核（全在**仓外导出树**，工作树零突变窗口）

| 变异 | 打哪枚守卫 | 读数 |
|---|---|---|
| **M1** | `aRelates` 去掉 `uncertainAnchorMiss`（＝`AC#2`/336 观察的那一支） | **红 3**：`TestAnchorFallbackTeethPerAnchor336`（报 `…roster_test.go:114/117/120/123`）＋339 两枚；mutant 逐字 `class: got class=B … reason="B-list … id_*"`＝**A 降 B 且可覆盖**，正是牙要抓的那格 |
| **M2** | `uncertainAnchorMiss` 首行 root 守卫（`AC#2` 点名的那枚） | **rc=0 / 0 FAIL**（负读数，见 §1 AC#2：落地后该守卫在 POSIX 结构上死掉） |
| **M3** | `resolveHandle` 回桩 `("",false)` | **红 4**：336＋339 两枚＋`TestResolveFollowsExemptedSymlinkOntoRealPath347` ⇒ 牙**载荷于解析器**，不是空转 |
| **M4** | `reparseComponents` 停止收 symlink | **红 2**：两枚 symlink 用例；mutant 逐字 `Resolved=true … Canonical=…/outside/secret.txt`（跟随到界外） |
| **M5** | `normalizeLocalUNC` 的 POSIX 早退 | **红 1**：`TestNormalizeLocalUNC…347`；mutant 逐字 `= "C:\x"` |

**三读数**：变异红＝上表 4 发（M1/M3/M4/M5，逐枚红句已具名）／**还原绿**＝`mut-restored.log` rc=0、0 FAIL（同一 `-run` 联合筛，含上表全部用例）／**原样绿**＝**从未变异**的导出树上同一联合筛 `linux-pristine-union.log` **rc=0、0 FAIL**。
**还原核（时机＝最后一发变异 M5 之后）**：5 枚产码件逐枚 `SHA256` 对工作树 **SAME**（`blacklist.go`/`pathresolver.go`/`pathresolver_other.go`/`syncdirs.go`/`syncdirs_other.go`）；`export-mut` 的 `git status` 回到与工作树同一集合。工作树**全程未发生变异**（突变只落在 `D:\tmp\wisp\20261026\347r1\export-mut`）。

## 6. A 档在 POSIX 上**第一次可达**（正向，真 symlink 用例）

- **指向界外 ⇒ 必须拒绝**：`<base>/profile/shortcut → <base>/outside/secret.txt`，`Resolve(link,nil)` ＝ `ErrReparseDenied`、`Resolved=false`；`reparseComponents` 恰返回 `[link]`（拼写照）；例外写成**目的地**（`target`）**不授权**（表里放的是组件，不是"同一棵树"）。
- **例外后 ⇒ 必须解析**（真 `realpath`，⛔ 原样传回）：`<home>/work/innocent.txt → <home>/.ssh/id_ed25519`，不带例外 ⇒ 仍拒（**界内也默认拒**：契约拒的是**遍历**）；把该组件列进 `[fs] reparse_point_exceptions` ⇒ `Resolved=true`、`Reparse=true`、`Canonical` **＝目标真路径**（≠ link 拼写），`Classify(Canonical)==A`（`~/.ssh/**`）。
- **A 档的两条来源在 POSIX 上都活了**：①**名字命中**（上面这条，靠 `realpath` 把 link 折成目标）；②**fail-closed 支**（`uncertainAnchorMiss`）——`TestAnchorFallbackTeethPerAnchor336` 八行＋`…Posix339` 两枚就是它的世界（Windows 那 8 枚牙直到 347 在这里**没有世界**）。

## 7. `internal/risk` 既存红（⛔ 当新红）

1. `TestC26RewriteAccountIsConsumedAtEverySecurityLeg`：**HEAD 复现**，两平台都红（`win-base-risk.log` 与 `linux-shotB.log` 的基线发），根因在 `internal/panel` 三处直读 `Result.Canonical`（禁写面），账 `A380`。本腿**未碰**。
2. `TestResolvePerCallBudget`：**时序**。包级两平台都绿；Windows `windows` 档一发红、**孤立跑绿**（`0.593 ms/op ≤ 1.000 ms`）。本腿改动在 `!windows` 文件里，**不在其路径上**。

## 8. 四数 ＋ 写面零越界

| 对象 | 改前 | 改后 |
|---|---|---|
| `internal/risk`（包级） | Linux `122/1/1`、Windows `147/1/1` | Linux `127/1/0`、Windows `147/1/1` |
| `internal/tools`（普查说会"开跑"的那枚） | Linux 门禁：`ok`（`TestWorkspaceSwitchRefusesAJunctionToOutside` **账本 skip**） | **同**。★**普查这里错**：`:198` 的 skip 是 `runtime.GOOS != "windows"` **平台守卫**，不是能力探针 ⇒ 落 `reparseComponents` 后它**仍然在 `:198` 跳**（现跑逐字），**零开跑**；账本那一行留着（只把"detection 是 Windows 实现"这句改准） |
| `cmd/wisp` | — | **零字节**（`git diff --name-only -- cmd/` 空；在飞 6 枚属 `wire-r1`） |

## 9. 格式甲乙两尺 ＋ `d22scan`

- **甲 `gofmt -l`**／**乙 `gofumpt -l`**（绝对路径 `/d/work/base/gopath/bin/gofumpt.exe`，归属判定用 `[\\/]`）：本腿 **22 枚改动/新增文件逐枚不在名单里**；两条尺都只报那 **3 条既存**（`cmd/wisp/models.go`、`internal/agent/approval/pending_read.go`、`internal/agent/tools.go`）＋`.scratch/**` 下的探针件（非本腿责任）。
- **`d22scan`**：在**仓外导出树**（`export-gate`，HEAD＋本腿文件、无 `wire-r1` 在飞件）`rc=0 clean`；**正控**＝同一工具扫同一棵树、只塞一枚 `U+2264` 到**字符串**里 ⇒ `rc=1` 且逐字点名 `internal/risk/zz_d22_positive_control_347.go:4: [emoji] ban #8 … string literals included`（⇒ 尺是活的，`clean` 不是说谎）。共享工作树那发**也 `rc=0`**（这次 `wire-r1` 的在飞件没触规），仍按派单以导出树读数为主。

## 10. 提交与门禁 rc

- 门禁 rc：Linux `core` **rc=1**（基线也 rc=1；差集只有上面那 1 枚 config）；Windows `windows` **rc=1**（红集与基线**逐名相同**，零差）。包级 Linux `rc=1`（只剩存量 panel 那枚）、Windows `rc=1`（同）。
- 提交：`0c730cac2`（20 枚文件，`+1064/−608`）。`add`/`commit` 两处**逐枚字面 pathspec**，主题逐字以「票 347 落解析器腿 347-r1」开头；`git show --name-only HEAD | grep '^cmd/'` **＝0**；**⛔ push**。**提交后复跑**（按 `git archive HEAD` 导出的树，⛔ 只信工作树）：Linux 包级 **`127/1/0` 同数**、红同枚。
- **遗留（具名，⛔ 藏）**：①`internal/config` 那 1 枚门禁新红（禁写面，已上报）；②`pathHandleVerify` 仍 `false`（`syncdirs_other.go`，只改注释）——落地后它只剩"Lstat 能过而 realpath 过不了"的残留（EACCES 等），收紧它属**独立决定＋独立分母**；③`syncdirs.go` 的 `s.home`／`roots[].canon` 仍是折叠形喂 `formsOf`（同类折叠缺陷，只影响精度不影响颜色，未动）；④macOS 的 `t.TempDir()` 走 `/var→/private/var`，**会被这枚默认拒命中**（本腿在 Linux 容器里复现了同形：§4.3 的 config 红）⇒ 票 55 落地前必须先解决例外注入。

## 11. 我没核（⛔ 当已核）

⛔ 未跑 macOS 实机／未交叉编译 darwin；⛔ 未跑 `go test ./...`（派单禁）；⛔ 未动 `internal/tools`／`internal/config`；⛔ 未验 `pathHandleVerify=true` 若收紧会红几枚；⛔ 未核 `A380` 台账原文（"存量"二字来自票面＋HEAD 复现）；⛔ 未验 GA CI（本机容器 ≠ ubuntu-latest 内核/权限）。
