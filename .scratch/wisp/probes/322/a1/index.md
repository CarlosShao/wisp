# 322-a1 三形代价表（只读普查腿 agent=322-a1，2026-10-15 21:3x +08）

锚点（起手现读）`HEAD=c673f67b`；仓内产码**零改动**（`git status --porcelain -- internal cmd docs .scratch/wisp/issues` 起手非 0 的那两枚是别家腿的面，见 §6）。
容器＝`golang:1.27` `CGO_ENABLED=0` uid 0，树放 **D:**（`/d/tmp/probe322/`，导出后 `find -name '*.go' | wc -l` = **1072** 证非空）；Windows 档＝本机 `go1.27.1 windows/amd64`。
探针源副本（本目录三枚 `.go.txt`）＝实跑源的**剥注释＋gofmt 归一副本**（**逐字节相同**，`cmp` 四对已证 SAME；本表引用的竞态读数即由这三枚原字节跑出，逐批对照见 `logs.md §3`）。
每一格标〔已测〕／〔未测〕；**未测一律具名留给下一轮，⛔ 用"按 POSIX 应该会"补位**。

## 0. 一句话（不是结论，是读数摘要）

叶子在窗口里被换成链接时，**`:97`（`O_CREATE|O_EXCL`）那一支两平台都关着**（实测 EEXIST、不跟随、见证逐字节不变）；**`:87`（`O_TRUNC`）那一支两平台都开着**，且 Windows 比 POSIX 更宽（真入口 3000 发里 `damaged=780`（三批 571/759/780）、`ok=1216`＝其中一批样本行 `YES … err=<nil>` 同色（截断了外来文件**还报成功**））。甲形在 POSIX 把该窗口关到 0 且现有钉一枚不红；**甲/乙两形在 Windows 都等于没修**（甲＝常量 0，乙的第二跳退化回跟随）；Windows 唯一量到的对应物是 `FILE_FLAG_OPEN_REPARSE_POINT` 那枚"不跟随的开"，它给判据、**不给写路径**。

## 1. 三形代价表（AC#0 主表）

| 形 | 改动面（文件∶行，现跑） | POSIX 实测读数（容器 uid 0） | Windows 实测读数（本机） | 顶红（实测＋读断言） |
| --- | --- | --- | --- | --- |
| **甲**：open 上 `O_NOFOLLOW`（平台拆常量） | `internal/winsec/winsec.go:113` `os.OpenFile(resolved, flags, perm)` → `flags|oNoFollowPlatform`；常量落 `winsec_other.go`（已 import `syscall`）＋ `winsec_windows.go` 各 1 行 | 〔已测〕by-name 全部**叶子**链接形（rel/abs/up/absout/dangling）→ `ELOOP`「too many levels of symbolic links」，见证 `size=70 mtime` 逐位不变；真文件 → `OPENED` 正常。**竞态真入口**：`PrivateFile` 3000 发 `damaged=0`（五批全 0；同法基线 `damaged` 五批 `108/130/124/89/111`），另有 107–139 发**新**拒绝原因是 `ELOOP` ⇒ 拦在了 act 上（活性非零） | 〔已测〕`syscall.O_NOFOLLOW` 在 `GOOS=windows` **undefined** ⇒ 该形 Windows＝`0`＝逐字没改；`GOOS=windows go vet` 在四份副本上 rc=1，**死因＝容器无模块网**（`x/sys` 下载被拒），不是产码问题 | 〔已测〕`go test -count=1 -v ./internal/winsec/` = **55 RUN / 33 PASS / 0 FAIL / 0 SKIP**，与基线**逐枚同名**；消费者 `risk/secret/memory/config` 全 `ok` |
| **乙**：`:87` 改独占开＋第二跳复核 | `winsec.go:87` 旗标 → `O_CREATE\|O_EXCL`；`privateFile`（`:108`，签名要加一枚"是否允许覆盖"）插第二跳；调用者两枚（`:87`/`:97`）逐枚点名 | 〔已测〕第二跳形状 `O_WRONLY\|O_NOFOLLOW`（带不带 `O_TRUNC` 都量了）：叶子链接 → `ELOOP` 拒；真文件 → `OPENED`。**竞态真入口**：`damaged=0`、`ok=8`（与甲同色）。**覆盖写承诺**：`PrivateFile` 连写两次 → `nil/nil`，内容 `SECOND`；预存 0644 文件 → 覆盖并收窄成 0600 | 〔已测〕乙的第二跳用的就是 `oNoFollowPlatform`，Windows 上＝0 ⇒ 纯 `O_WRONLY\|O_TRUNC`，**照旧跟随**：真入口 3000 发 `damaged=1022`（基线 780，同量级、没关） | 〔已测〕同上：**0 红**。**⚠ 关键代价**：把第二跳摘掉（纯"已存在即拒"）后，`./internal/winsec/` 与 `risk/secret/memory/config/agent` **仍然全绿**，而 `PrivateFile` 第二次写返回 `file exists`、内容停在 `FIRST`、预存的 0644 文件**停在 -rw-r--r-- 没被收窄**（四批 `refused-exist` 366/378/401/462，`damaged=0`、`ok=0`） ⇒ 覆盖写承诺**今天没有钉**（见 §4） |
| **丙**：不动 open，截断后同一对象复验再报错 | `winsec.go:113` 之后插 `f.Stat()` vs `os.Lstat(resolved)`（POSIX 侧这已是 `winsec_other.go:261-274` `sealHandle` 里**现存**的两条判据 ⇒ 丙≈不新增东西） | 〔已测〕信号**看得见叶子换手**：`SameFile=false linkNow=true`；**看不见祖先换手**（`anc-sib`/`anc-out` → `SameFile=true linkNow=false`）。真入口 `damaged=142`（四批 165/135/148/142；基线同批 108/130/124/111 ⇒ 同量级）＝**报告了损伤、没防住损伤** | 〔已测〕同一对象尺在 Windows 也**只对叶子有效**（W2 `SameFile=false linkNow=true`／W4 祖先 junction `SameFile=true`）；且 Windows 的 `sealHandle`（`winsec_windows.go:328`）＝ `applyDescriptor(f.Name(), false)`，**根本没有这道闸门**（120-r2 那半是 `//go:build !windows`）⇒ 本机基线真入口 `damaged=780`／`ok=1216`（三批 571/759/780；样本行 `YES … err=<nil>` n=367 与 n=113） | 〔已测〕**0 红**（丙形不改任何现有承诺） |

`O_EXCL` 那一支的稳态读数（票面 ①，容器）：**EEXIST、不跟随、见证 size/mtime 逐字节不变**，叶子链接的五种拼法（相对／绝对／爬出根／绝对但目标在根外／悬空）**全部一致**；`errors.Is(err, fs.ErrExist)=true`（`PrivateFileExclusive` 承诺的哨兵在链接态也成立）。⚠ 一枚附带读数：`O_TRUNC|O_CREATE` 遇**悬空**链接会在树外**创建**对象（`size=0`）＝另一种损伤形状，票面没写。

## 2. Windows 那一半（AC#3 具名答案）

〔已测〕本机可创建 file symlink＋dir symlink＋junction（`W1 …=<nil>`）⇒ 三形在 Windows **都能量**，不是"无从下手"。
1. `:87` 形 over symlink leaf → `OPENED`＋见证 `size 70→0`；over **祖先 junction** → 同样 `DAMAGED`。**洞在本机是活的**。
2. `:97` 形 over symlink leaf → `The file exists.` `is(ErrExist)=true`＋见证不变；over 祖先 junction → 同样拒。**独占那一支两平台都关着。**
3. `os.Root` 在 Windows 不兜：相对目标叶链接 `O_TRUNC` → `OPENED`＋`DAMAGED`；`Root.Create` → `OPENED`＋`DAMAGED`；绝对目标 → `path escapes from parent`；`O_EXCL` → 拒。⇒ 与 POSIX 同结论：**别把 Root 的链接规则当防线**（`120-r1`/`120-r2` 那句在 Windows 也成立）。
4. 对应物＝`FILE_FLAG_OPEN_REPARSE_POINT`：`OPEN_EXISTING`＋该旗标 → 句柄命名**链接本体**（`attrs=0x420 reparse=true size=0`，file index 与目标不同：`…591` vs 目标 `…590`）；此后 `Write(20 字节)=nil`、`Truncate(0)=nil` 都"成功"但**见证仍 size=70、链接未坏**（`readlink` 可读）⇒ 它能**拒绝**、**不能当写路径**（写进去的字节到不了调用方以为的对象，而且返回成功）。
5. ⚠ 陷阱一枚：同一枚链接改用 `CREATE_ALWAYS`＋该旗标 → 句柄报 `attrs=0x20 reparse=false`、**链接条目被换成一枚空普通文件**（见证未损、链接没了）⇒ ⛔ 任何 Windows 落地把 `CREATE_ALWAYS`/`O_CREATE` 当"不跟随的覆盖开"。
6. 因此本腿给编排者的具名回答是：**甲/乙两形 Windows 无对应物（甲＝0、乙的第二跳照旧跟随，均已实测）；Windows 侧唯一可用的形状是"reparse 旗标开 → 句柄报 reparse 即拒"，成本＝一次额外 `CreateFile`，它关的是"判据"那一半，写路径需要另一跳。** 〔仅本机可量〕CI/容器读数覆盖不到 Windows：本轮**未**取 Windows runner 的符号链接权限档（下一轮）。
7. 另外具名：`winsec_windows.go:328` 的 `sealHandle` 没有 120-r2 的同一对象闸门 ⇒ Windows 上"私密字节落到外来对象"也仍是开的（本机读数 `ok=868`，其中 damaged&nil 480 发）——**这一条不属于票 322 的射程，属票 120 的 Windows 半格**，本腿只登记、不动。

## 3. 顶红名册（逐名读断言 ＋ 实测色；行号＝本轮现跑）

实测色先说：**四形（甲／乙／丙／乙无第二跳）× `./internal/winsec/` = 0 FAIL / 0 SKIP，与基线逐枚同名**（`placement_symlink_113` 8 顶层＋4 子、`ancestor_separator_108` 4、`dataroot_symlink_119` 5、`seam_probe_root_125` 3、`posix_seal_window_120` 3、`placement_leaf_118` 2、`private_fail` / `private_other` 全在 33 枚 PASS 里）。**票面点名的五枚一枚都不会因为这三形变红**，原因写在下面"射程"一栏。

| 件∶枚（行号） | 它钉什么 | 三形为什么顶不红它／**什么改法会红** |
| --- | --- | --- |
| `placement_symlink_113_other_test.go:277` `TestAC1POSIXPrivateFileThroughASymlinkRefusesAndWritesNothing`（调用在 `:286`） | `PrivateFile(link/secret.txt)` 必须**返回 `ErrUnresolvedPath` 且错误里点到这枚链接**（`assertRefused113` `:189-205`，`refusalCreditsLink137` `:160`），并断言树外没多出 `secret.txt`、名字本身没被创建 | 稳态（叶子**本来就是**链接）在 open 之前就被 `resolveString` 拒 ⇒ 三形都轮不到。**会红的改法**：① 把拒绝原因改成裸 `*PathError`（`ELOOP`/`EEXIST`）而不再包 `ErrUnresolvedPath`／不再点名 planted link ⇒ `:197` 与 `:201` 两条一起红；② 任何"先开后判"的形状若把这条路径变成 `nil` 返回 ⇒ `:193` 红 |
| `placement_symlink_113_other_test.go:353` `TestAC3POSIXSealStillWorksNextToAndThroughRealDirectoriesAndLinks`（`PrivateFile` 在 `:376`，perm 0644，断言 `:379` 落到 0600）／`:319` 同族 | 反半："一切皆拒"不许绿：树内**普通路径**必须写成、且收窄到 0600 | **会红的改法**：把甲形写成"open 前先 `Lstat`，是链接就拒"却没走 `O_NOFOLLOW`／或乙形第二跳丢了 `O_TRUNC` 之后不收窄 ⇒ `:377` `t.Fatalf`；perm→mode 的收窄丢了 ⇒ `:380` |
| `placement_leaf_118_other_test.go:96` `TestAC118POSIXPrivateFileRefusesALinkStandingWhereTheFileWasNamed`（调用 `:106`，断言 `:113` 内容未变／`:117` mode 未变／`:121` 链接还在） | 叶位置（链接**就是**实参）：`PrivateFile` 不许换掉受害者的字节、不许改它的 mode、不许删那枚链接 | 三形（稳态）都在 open 前被等住 ⇒ 不红。**会红的改法**：① 若把甲写成"只加 `O_EXCL`、丢 `O_NOFOLLOW`"，ELOOP 变成 EEXIST 后仍走第二跳跟随 ⇒ `:114`；② 若"拒绝"顺手 `os.Remove(resolved)`（现有失败路径 `winsec.go:119` 就删名字！）⇒ `:121` 红——**这一条是本表里离三形最近的一枚**，落地腿必须检查新加的拒绝支有没有继承那次 unlink |
| `ancestor_separator_108_other_test.go:80`/`:115`/`:146`（全是 `RemoveUnlinked`）＋`:181` `TestAC4POSIXFloorAnswersInsideTheNamedTree`（`PrivateDirAll`/`ResolvePath`，**不调 `privateFile`**） | 祖先分隔符口径与"底线不许改写拼写" | 与 `:87`/`:97` 无交集 ⇒ 三形都不红。**会红的改法**：动 `pathPieces`/`firstLinkAncestor`（不在本票射程）。⚠ 只在有真实 C26 pipeline 的二进制里 `:182` 会 `Skipf` ⇒ 数它时先证它跑过（`-count=1` 容器里 `SKIP=0` 是真跑） |
| `dataroot_symlink_119_other_test.go:259` `TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119`（`PrivateFile` 在 `:300`：解析过的根内**普通路径**必须写成）＋`:139`/`:165`/`:224`/`:336` | 根内普通路径必须**仍可写**（防"过度拒绝"），同时根内链接必须拒 | **会红的改法**：甲若写成"祖先里出现任何链接都拒"（O_NOFOLLOW 不管祖先，但有人会用 `firstLinkAncestor` 去补）→ `:267` `PrivateDirAll` 或 `:300` 红；乙若无第二跳 → 本枚仍绿（它只写一次）⇒ **这枚拦不住乙的承诺退化** |
| `seam_probe_root_125_other_test.go:119`/`:222`/`:275`（只走 `resolverProbeShapes()`/`resolverConformanceFailure()`/`builtinVerifier`，**不调 `privateFile`**） | 安装 C26 解析器的接缝守卫 | 三形都不碰守卫 ⇒ 不红。**会红的改法**：把 `resolveString` 换成"自己 open 再问平台"或改动探针根拼写 ⇒ `:134`/`:253`/`:295` |
| `posix_seal_window_120_other_test.go:321`/`:378`/`:404`（**只驱动 `SealFile`**：调用点 `:364`/`:479`） | 密封（chmod）那族窗口必须关到 0 且活性非零（默认严格档） | 三形都只改 `privateFile` 的 open ⇒ **一枚都不红，也一枚都量不到本票的窗口**。⚠ **本票 AC#1 不能拿这三枚冒充**：需要一枚同款但驱动 `PrivateFile` 的仪器（本腿的 `probe-racer.go.txt` 就是这个形状的雏形，⛔ 未入库） |
| `private_fail_test.go:26` `TestAC5FailedSealRefusesTheWrite`（含 `:46` `PrivateFile` 注入失败支）／`:99` `TestAC5FailureIsNotSwallowedByTheHappyPath`（`:108` 断言独占创建**第二次必须 `fs.ErrExist`**） | 失败方向只许收紧；`ErrExist` 哨兵不许丢 | 四形全绿。**会红的改法**：① 乙第二跳把 `ErrExist` 吃掉又改成 `nil` 返回 → `:108` 红（`PrivateFileExclusive` 不许走第二跳！`winsec.go:97` 那枚必须原样哨兵）；② 丙/甲的拒绝支若绕过 `sealError` 包装 → `:37`/`:51` `ErrNotSealable` 断言红 |
| `private_other_test.go:25`（`PrivateFile` 0644 → 期望落 0600）／`acl_windows_test.go:307`/`:310`／`reparse_windows_test.go:313`（Windows 档，本轮容器**跑不到**） | 模式/描述符落点承诺 | POSIX 已覆盖那枚；**Windows 三枚＝〔未测〕**（见 §5） |

## 4. `:87` 那一支的全仓调用者与承诺（AC#2；grep 有界，全带文件∶行）

生产调用者（5 处 4 文件）＋各自今天拿到的承诺，**本轮实测**（`P3` 三行＝`logs.md`）：

| 调用者 | 形状 | 依赖的承诺 | 乙（带第二跳） | 乙（纯拒绝） |
| --- | --- | --- | --- | --- |
| `internal/secret/store.go:79` `Store.Set` → `s.blobPath(id)` | `:87` | **同 id 再存必须覆盖**（唯一一枚"重写自己已有的路径"的生产调用者） | 仍成立（实测 `nil/nil`＋内容 `SECOND`） | **破**（实测第二次 `file exists`、内容留 `FIRST`） |
| `internal/config/migrate.go:93` | `:87` | 注释逐字：'because this call **rewrites an existing backup** rather than skipping it, PrivateFile is also the **repair path** for a wide backup left by an older build' | 仍成立 | **破**（且预存 0644 文件停在 `-rw-r--r--` 没收窄 ⇒ 旧构建留下的宽备份再也修不好） |
| `internal/secret/migrate.go:169`（`backupPath`） | `:87` | 有 `os.Stat` 前置：只在**不存在**时写，存在就走 `SealFile` 修复支 | 无所谓 | 不受影响 |
| `internal/secret/migrate.go:189`（`tmpPath`，pre-rename temp） | `:87` | 崩溃残留的 tmp 下次必须被覆盖写 | 仍成立 | **破**（残留 tmp 会让迁移永久失败循环） |
| `internal/agent/spill.go:166`（`path+".retry"` temp＋rename，D31 last-writer-wins） | `:87` | 残留 `.retry` 必须能被覆盖 | 仍成立 | **破** |
| `:97` 两枚：`internal/agent/spill.go:343`、`internal/memory/artifacts.go:75` | `O_EXCL` | 要的就是 `fs.ErrExist`（`artifacts.go:66-69` 注释逐字：'A collision is `fs.ErrExist`, **not an overwrite**'） | 不许被卷进第二跳 | 原样 |

测试侧调用者（点名，别漏）：`internal/winsec/acl_windows_test.go:307,310`、`private_fail_test.go:46`、`private_other_test.go:25`、`placement_symlink_113_other_test.go:286,376`、`placement_leaf_118_other_test.go:106`、`dataroot_symlink_119_other_test.go:300`、`cmd/wisp/secret_dataroot_119b_test.go:253`、`reparse_windows_test.go:313`。
**★读断言＋实测双证的结论**：全仓**没有一枚**现有用例对同一枚路径调 `PrivateFile` 两次 ⇒ 覆盖写承诺是"活的、无钉的"。若选乙，落地腿**必须自带新钉**（同路径两次 `PrivateFile` → 第二次 `nil` 且内容替换），否则承诺丢了没人报红。

## 5. 未测格（逐枚具名，留给下一轮）

1. darwin/freebsd/openbsd/netbsd/dragonfly/solaris/aix 上 `O_NOFOLLOW` 的**行为**读数——本轮只量到**编译可用**（见 `logs.md` 矩阵：`COMPILES` 逐枚）；行为只有 linux 档。
2. `ios/*`（`go build` 死于"requires external (cgo) linking"，未定案）、`android/*` 两枚未进矩阵。
3. **Windows CI runner 能否创建 symlink/junction**（AC#3 的"CI 覆盖不到"这一格只有本机读数）；`acl_windows_test.go`/`reparse_windows_test.go` 那两枚 Windows 档在本轮任何一档都**未跑过**。
4. **uid 1000 档**（容器只跑 uid 0；120 的先例：权限面在 uid 0 失真）——三形是否新增 EACCES 形状未量。
5. 竞态命中的**方差/区间**：同形四次抓取 base `108`/`130`/`124`/`111`、bing `165`/`135`/`148`/`142`、yi-nofb 的 `refused-exist` `366`/`378`/`401`/`462`（各 3000 发），未算分布、未换 CPU 压力。
6. Windows racer 的仪器缺陷：`flip-errs` base-win `580`/`808`/`768`、yi-win `253`/`570`/`412`（三批）——`os.Rename` 抢占会失败 ⇒ Windows 分母失真，两枚数只能当"窗口确实到达过链接态"的凭据，不能当命中率。
7. 甲/乙在 **FIFO、硬链接、UNC、`\\?\`、8.3 短名、尾点/空格** 拼写下的形状（`resolve.go` 的射程，本票未碰）。
8. `O_NOFOLLOW` **经 `os.Root`** 在 Windows 档：本轮只量了不加旗标的跟随（W5a/W5c）。
9. `go test ./cmd/wisp`、`go test ./internal/tools`（派单明令本腿⛔跑，两面上有腿在飞）⇒ AC#4 的"消费者不推红"只覆盖 `risk/secret/memory/config/agent`。
10. `:97` 独占那一支的**窗口**形态：本轮量的是"叶子已是链接"的稳态（两平台 EEXIST）；"检查后被换成链接"需要一枚 `PrivateFileExclusive` 版竞态仪器，未造。
11. 甲形对 `privateFile` **目录**侧（`PrivateDirAll`）与 `SealDir` 的影响面：未碰、未量（本票只管文件那一跳）。

## 6. 纪律与自抓（不藏）

- 写面只有 `.scratch/wisp/probes/322/a1/**` ＋票 322 尾 1 条 bullet；`internal/**`/`cmd/**`/票面正文/AC 框/Status/台账/`docs/**` **零改动**（提交前后各一枚 `git status` 自证，见 `logs.md` 末行）。
- 起手检查非 0 的别家面**具名**：`internal/tools/subagent_197.go`(M)＋`internal/tools/nilbase_270_test.go`(??)=`270-r1` 在飞；本腿没动过它们，也没跑它们的包。
- ⛔ 仓内产码突变：三形的"顶红"是在**仓外副本**（`/d/tmp/probe322/mut/{base,jia,yi,bing,yi-nofb}`）上物理改写后跑的，仓内字节全程未变；副本留底在盘（只建不删）。
- 自抓一枚：`mv` 移走过一份探针源（`exp3-windows/main.go` → `mut/base/probe322win_tmp`，**未删**，随后 `cp` 复位；两处现在都在盘上）。入账归编排者。
- 台面故障新增（详见 `logs.md` §末）：① 后台 `docker run` 被 CLI 判"completed rc=0"而容器仍在跑 ⇒ 尾段没落盘；② 容器无外网会让需要模块的包以 `connection refused`+`[setup failed]` **冒充测试红**，定式＝只读挂宿主模块缓存。
