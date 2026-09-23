# 124 AC#2b 批次 3 后一半 — `internal/winsec` 的 17 枚：先裁该不该转，再动手

## 0. 锚点与环境

- 开工首读 `git rev-parse --short HEAD` = **`4ea0db2`**（分支 `dev`）。开工时 `git status --porcelain` 只有一枚
  来源未明的未跟踪文件 `docs/reports/2026-09-23-gap-analysis-vs-oss-harnesses.md` —— 按派单：**未提交、未改、未删、未据它开票或改判据**。
- 纯净快照：`git archive 4ea0db2 | tar -x -C /d/tmp/wisp124-2b3b`（461 枚 `.go`），**未在仓内建 worktree、未 checkout**。
  快照内 `md5sum internal/winsec/resolve.go` = **`b6876a5efe759f6e17434d1b50a129c3`**，与 AC#1／AC#2a／2b-1／2b-2 逐字同字
  ⇒ 两批前置动码（`fafe2b4`、`dfa3dc4`）**没有碰到 `internal/winsec`**，本锚的 winsec 码就是它们量过的同一版本。
- 容器：`golang:1.27`（镜像 `3680233e3204`），容器内 `go version` = **`go1.27.1 linux/amd64`**。
  挂载一律 `/d/...` + `MSYS_NO_PATHCONV=1`（派单警告的 `docker run -v "C:\…"` 静默空挂＝假绿），
  每发进容器第一件事 `ls -l /src/go.mod` = `-rwxrwxrwx 1 root root 883`（证明文件真在，非空挂）。
- 形状：`ln -s /realpriv /varlink` + `TMPDIR=/varlink/w124tmp`（软链形）；`TMPDIR=/plainroot/w124tmp`（普通形）。
  硬断言沿用 2b-1 的 `exit 97`（空挂）／`exit 98`（形状没建成）／`exit 99`（普通形被 `/varlink` 污染），本批每一发均未命中。
- 跑法脚本：`/d/tmp/wisp124-2b3b-run.sh`（CRLF 原件）+ `/d/tmp/wisp124-2b3b-run-lf.sh`（LF 副本，容器实际执行的那枚）。
- 时间戳：本件所有 +8 时间由 `date -u` 换算（开工 12:52 UTC / 20:52 +8）。

## 1. 17 枚名单，与派单给的三条实据的**自核**原文

### 1.1 名单（`./internal/winsec/`，只在软链形红）

AC#1 在 `7b4c36a` 上量出 17 枚；本批在 `4ea0db2` 的**改前基线**（容器软链形 `-count=1 -v`，
`/d/tmp/wisp124-2b3b-logs/BASE-L.log`）复算，**逐名零差 = 13 枚顶层 FAIL + 4 枚子测试 FAIL**：

| # | 用例名 | 文件 | 构建标签 | 宿主(Windows)编不编译 |
|---|---|---|---|---|
| 1 | `TestAC5FailedSealRefusesTheWrite` | `private_fail_test.go:26` | 无标签 | **编译** |
| 2 | `TestAC5FailedSealRefusesTheWrite/exclusive_artifact` | `private_fail_test.go:30` | 无标签 | **编译** |
| 3 | `TestAC5FailedSealRefusesTheWrite/replacement_write` | `private_fail_test.go:43` | 无标签 | **编译** |
| 4 | `TestAC5FailedSealRefusesTheWrite/directory_chain` | `private_fail_test.go:56` | 无标签 | **编译** |
| 5 | `TestAC5FailedSealRefusesTheWrite/seal_file_and_dir_direct` | `private_fail_test.go:66` | 无标签 | **编译** |
| 6 | `TestAC5FailureIsNotSwallowedByTheHappyPath` | `private_fail_test.go:99` | 无标签 | **编译** |
| 7 | `TestPOSIXPrivateFileIsReally0600` | `private_other_test.go:20` | `!windows` | 不编译 |
| 8 | `TestPOSIXPrivateDirIsReally0700` | `private_other_test.go:37` | `!windows` | 不编译 |
| 9 | `TestPOSIXSymlinkAtArtifactPositionIsNotRecursed` | `private_other_test.go:60` | `!windows` | 不编译 |
| 10 | `TestPOSIXMissingFileIsNotAnError` | `private_other_test.go:97` | `!windows` | 不编译 |
| 11 | `TestAC2POSIXDoesNotFoldABackslashIntoASeparator` | `ancestor_separator_108_other_test.go:94` | `!windows` | 不编译 |
| 12 | `TestAC4POSIXFloorAnswersInsideTheNamedTree` | `ancestor_separator_108_other_test.go:152` | `!windows` | 不编译 |
| 13 | `TestAC118POSIXSealFileRefusesALinkStandingWhereTheFileWasNamed` | `placement_leaf_118_other_test.go:75` | `!windows` | 不编译 |
| 14 | `TestAC118POSIXPrivateFileRefusesALinkStandingWhereTheFileWasNamed` | `placement_leaf_118_other_test.go:96` | `!windows` | 不编译 |
| 15 | `TestAC3POSIXSealFileStillNarrowsAPlainFileInsideTheNamedTree` | `placement_symlink_113_other_test.go:253` | `!windows` | 不编译 |
| 16 | `TestAC3POSIXSealStillWorksNextToAndThroughRealDirectoriesAndLinks` | `placement_symlink_113_other_test.go:287` | `!windows` | 不编译 |
| 17 | `TestAC3POSIXSealDoesNotFoldABackslashIntoASeparator` | `placement_symlink_113_other_test.go:335` | `!windows` | 不编译 |

宿主侧不是推断，是量出来的：`GOOS=windows go test ./internal/winsec/ -list '.*'`（在 `4ea0db2` 的纯净快照里跑）
⇒ 上表 13 枚顶层名逐名 `-cx` 计数：第 1、6 枚命中 1，第 7–17 枚命中 **0**。
⇒ **派单实据③成立且可量化：17 枚里 11 枚在宿主上连编译都不参与，只有 6 枚（#1–#6，`private_fail_test.go` 无标签）宿主有分母。**
本件后面**每一个数都标"在哪台、哪个构建标签"**。

### 1.2 三条实据的自核结论（我的转述也是断言，逐条重走）

**① 票 119 面 `R-119-9` 的禁令原文 —— 对得上，逐字为真。**
`.scratch/wisp/issues/119-posix-link-leg-refuses-legitimate-symlinked-data-roots.md:73`：

> - **（2026-09-22 返修 `agent-ticket119b` 新增，来源 `R-119-9`）不许拿被测函数算 fixture。**
>       一枚用例的期望值只能来自文件系统或调用方自己声明的字面值，不能来自被测函数（含其幂等组合）：
>       `SealableRoot(base)` 当期望值 ⇒ 把"声明的树原样返回"这条纪律抹掉也照样绿（验收方变异 R 实测 56/56 全绿）。
>       用例落笔前先问一句"哪一发生产变异能把它打红"，答不出就是恒真，别当覆盖交出去。

**范围要点（决定本批怎么用这条）**：这条钉的是**期望值**的来路，不是**输入根**的来路。本批 17 枚的改动全部落在
"测试自己把哪份根递给底线"这一侧（输入），没有任何一枚的期望值改成由被测函数算出来。§2 逐枚复核过这一点。

**② `internal/proc/envfork.go` 头部明写 `nothing in internal/winsec calls it` —— 对得上，逐字为真。**
`internal/proc/envfork.go:148`（`SealableRoot` 的文档块内）：

> // cleans, never absolutizes, never case-folds, and nothing in internal/winsec
> // calls it - the only C26 entry point stays internal/risk/pathresolver.go.

并且这条边界**在 winsec 这一侧被第二次写下并解释了理由**——`internal/winsec/resolve.go:224-228`（票 125 的注释）：

> // (internal/proc's SealableRoot) - applied to the only OS read this package makes
> // about its *own* fixture. It is deliberately not a call into internal/proc:
> // envfork.go's boundary note ("nothing in internal/winsec calls it") is the
> // package doc ticket 113 AC#6 wrote, and the floor must not start depending on
> // the layer above it to keep that boundary checkable.

⇒ **批次 2 那句"winsec 不许照抄本批的委托形状"逐字成立**，而且理由比"幂等所以恒真"更硬：这是一条**双向写下的地界**，
票 125 已经为同一道题给出过答案——winsec 要解析自己的 fixture 根，就在包内自己走（`resolveProbeRoot`，`resolve.go:253`），
**不伸到上层去买**。本批照这个先例走，不接 `proc.SealableRoot`。

**③ winsec 有一枚 `*_other_test.go`（非 windows 标签）宿主根本不编译 —— 对得上，且比转述更强。**
转述只说"有一枚"；实测是**6 枚文件带 `!windows`**（`ancestor_separator_108_other`、`dataroot_symlink_119_other`、
`placement_leaf_118_other`、`placement_symlink_113_other`、`private_other`、`seam_probe_root_125_other`），
覆盖 17 枚里的 **11 枚**；另有 `private_fail_test.go` **无标签**，覆盖 6 枚。
⇒ 派单要求的"读数必须说明在哪台、哪个构建标签"本批落成硬格式：§3/§7 每个数后面带 `〔容器/!windows〕` 或 `〔容器/无标签〕`。

### 1.3 派单没给、但我撞上的两条（都要报，不当已验证用）

**(a) `proc.SealableRoot` 在 winsec 里**不是**编译不通**。批次 1 的 `next=` 让下游怀疑接不上；我实测把它接通了：
在 `/d/tmp/wisp124-2b3b-cycle`（快照副本，**仓外**）放一枚 `package winsec` + `//go:build linux` + `import "…/internal/proc"` 的探针
`zz_cycleprobe_124_test.go`，容器原生 `go vet ./internal/winsec/` = **rc=0**、`go test -run NONE` = **ok**、
`go test -list 'CycleProbe.*'` 列得出 `TestCycleProbe124`。根因：`go list -f '{{join .Imports " "}}' ./internal/proc`
= `context errors fmt github.com/CarlosShao/wisp/internal/buildinfo` —— **`internal/proc` 并不 import `internal/winsec`**，两包无环。
⇒ **"不许照抄委托形状"是一条纪律／地界约束，不是一条编译约束**。这点必须纠正，否则会下游会以为 winsec 是"转不了"。
（我一开始按 `go list -deps ./internal/proc` 的输出判过"有环"，那枚输出把 winsec/secret/observe 都列进来了，与
`.Imports`、`grep -r` 与**实际编译**三条都矛盾；以三条一致的为准，`-deps` 那一读登记为来源未明、不复用。）

**(b) 软链形下 winsec 的 113/108 族拒绝腿是"绿得没有理由"的**（不在 17 枚里，是本批改动面必然牵到的邻居）。
同一支仪器在两形下的逐字对比（`placement_symlink_113_other_test.go:153` 的 `t.Logf`）：

- 软链形 `BASE-L.log:111`：`AC#1 SealFile("/varlink/.../002/root/link/keep-me.txt") -> err=winsec: refusing to seal …
  **reaches it through the link at /varlink**, which is not the tree this call names`
- 普通形 `BASE-P.log:128`：`AC#1 SealFile("/plainroot/.../002/root/link/keep-me.txt") -> err=winsec: refusing to seal …
  **reaches it through the link at /plainroot/.../002/root/link**, which is not the tree this call names`

⇒ 普通形拒的是**用例自己种的** `root/link`；软链形拒的是**宿主的** `/varlink`，而 `assertRefused113`
（`placement_symlink_113_other_test.go:133-138`）只要求 `errors.Is(err, winsec.ErrUnresolvedPath)` —— 两根链接给出**同一个 sentinel**。
⇒ 这几枚在软链形**过，但没过在它测的那条腿上**。逐枚点名（软链形 PASS、且日志里拒因写着 `/varlink`）：
`TestAC1POSIXSealFileThroughASymlinkRefusesAndLeavesTheForeignTreeAlone`、
`TestAC1POSIXSealFileRefusesALinkAncestorAtEveryDepth` 的 4 枚子测试、
`TestAC1POSIXSealDirThroughASymlinkRefuses`、`TestAC1POSIXPrivateFileThroughASymlinkRefusesAndWritesNothing`、
`TestAC1POSIXSealFileThroughABackslashNamedLink`、`TestAC2POSIXAncestorGuardRefusesASpellingThroughASymlink`、
`TestAC2POSIXABackslashInALinkNameIsStillALinkAncestor` = **8 枚顶层 + 4 枚子测试 = 12 枚**。
这 12 枚**不在 17 枚的分母里**（它们是 PASS，不是 FAIL），但它们的 fixture 根由 `newForeign113`／`foreignTree108`
两枚共享 helper 供给，**转 17 枚必然同时转它们**（§8 逐枚登记）。**这条是本批交回的最重要发现，请 2b-4 与验收方复算。**

## 2. 甲／乙／丙 三分类，逐枚给理由（本件的主产物）

### 2.0 裁的口径

三档按派单定义。判"甲类"要同时过两条：(i) 它红是因为**测试根的来路不对**（与前两批同一枚病）；
(ii) 换掉这枚根之后，**它被测的东西还在**。判"乙类"要拿出**断言原文**证明它要的就是"未解析那一形"被拒。
判"丙类"要给出形状／依赖／生产侧的**可复核**理由。

先把两条**通用事实**立住，后面逐枚只写这枚自己的部分：

- **F1（红因全是宿主的链接，没有一枚是用例种的链接）**：17 枚的软链形失败串逐字回看（`BASE-L.log`，全文见 §3.1），
  只有三种文案，且**每一种点名的链接都是 `/varlink`**（宿主自己那枚），不是用例在树里种的 link：
  (a) `winsec: path is not provably resolved, refusing to seal …reaches it through the link at /varlink, which is not the tree this call names`；
  (b) `winsec: entry is a link to something else, not private data …the spelling reaches it through the link at /varlink, and whatever lives behind that link is not this tree's data to delete`；
  (c) `the leaf link does not answer with its own target: "/realpriv/…" vs "/varlink/…" (err=<nil>)`。
- **F2（`SealableRoot` 幂等 ⇒ 普通形逐字不动）**：普通形下 `t.TempDir()` 的整条前缀里没有链接，
  解析＝原样返回 ⇒ 改动**在普通形是字面 no-op**。这是判据②"普通形一枚都不许多红"的机理，§4 再用量去核。

**结论先给：甲类 17 / 乙类 0 / 丙类 0。** 派单预留的"绝大多数属乙类、基本不该动码"这一支**没有成立**——
理由是 §2.1–§2.3 逐枚的断言原文：17 枚里**没有任何一枚**的断言是"就是要拿未解析的根去试底线拒绝"。
下面逐枚，并按派单要求把"如果哪天有人把它转绿会丢掉什么"这一栏**反向**填成"它现在红着会挡住什么"（乙类为 0，所以那一栏本批全空）。

### 2.1 `private_fail_test.go`（无标签 ⇒ 宿主有分母）—— 6 枚，全甲类

| # | 用例 | 断言原文（被测对象） | 类 | 本枚理由 |
|---|---|---|---|---|
| 1 | `TestAC5FailedSealRefusesTheWrite` | `private_fail_test.go:35` `t.Fatal("the write succeeded although the seal was refused")`、`:38` `t.Fatalf("error does not name the refusal: %v", err)`（配 `!errors.Is(err, ErrNotSealable)`） | **甲** | 它要的是**注入的** `applyDescriptor` 失败（`withInjectedSealFailure`，`:14-19`）。软链形下底线自己先拒（F1-a），sentinel 不是 `ErrNotSealable` ⇒ 拒因被换掉了。被测对象是"密封失败时写不写得下去"，与根的来路无关：换根后注入腿重新是操作性拒因。 |
| 2–5 | 上表的 4 枚子测试 `exclusive_artifact`／`replacement_write`／`directory_chain`／`seal_file_and_dir_direct` | `:40`／`:53` `assertNoBytesOnDisk(t, p)`（`:88` `t.Errorf("refused write left %d bytes on disk at %s", …)`）；`:61` `t.Fatal("PrivateDirAll reported a private directory it could not seal")`；`:72-77` `SealFile`/`SealDir` 要 `ErrNotSealable` | **甲** | 同上：量的是"拒绝之后不留敏感字节、不留半成品目录"。`assertNoBytesOnDisk` 的期望来自**文件系统**（`os.Stat` 的 `fs.ErrNotExist`），不是被测函数 ⇒ 换根不触 R-119-9。 |
| 6 | `TestAC5FailureIsNotSwallowedByTheHappyPath` | `:103` `t.Fatalf("PrivateFileExclusive: %v", err)`、`:109` `exclusive create must keep failing on an occupied name` | **甲** | 它是 #1 的反向半（不注入时必须成）。软链形红在 `:103` 的 `t.Fatalf` —— **正向那半根本没跑**，所以这一对目前是"一枚红遮一枚没测"。换根后两半才真的互为对照。 |

⚠ 这 6 枚是 17 枚里**唯一宿主有分母**的（无 `!windows` 标签）。⇒ 本批给它们换根必须**在 Windows 上是 no-op**，
否则会把宿主/CI 的 windows 腿一起动了。做法见 §8（平台分文件的 no-op 那一支）。

### 2.2 `private_other_test.go`（`!windows` ⇒ 宿主零分母）—— 4 枚，全甲类

| # | 用例 | 断言原文 | 类 | 本枚理由 |
|---|---|---|---|---|
| 7 | `TestPOSIXPrivateFileIsReally0600` | `:26` `t.Fatalf("PrivateFile: %v", err)`；`:33` `t.Errorf("artifact mode is %v, want -rw-------", got)` | **甲** | 被测对象是 **POSIX 模式位**（`0o644` 进、`0600` 出，靠 `applyDescriptorPOSIX` 回读）。软链形红在 `:26` 的 setup 拒 ⇒ 模式那半根本没读。换根不碰模式断言，也不产生任何"期望值来自被测函数"。 |
| 8 | `TestPOSIXPrivateDirIsReally0700` | `:40` `t.Fatalf("PrivateDirAll: %v", err)`；`:48` `%s mode is %v, want drwx------` | **甲** | 同 #7，方向是目录 `0755` 进 `0700` 出。 |
| 9 | `TestPOSIXSymlinkAtArtifactPositionIsNotRecursed` | `:73` `t.Fatalf("RemoveUnlinked: %v", err)`；`:76` `the link survived`；`:79` `TARGET DELETED - removal followed the symlink` | **甲**（但**最接近乙**，逐字核过） | 它**自己种**一枚 link 在 artifact 位置（`:77 os.Symlink(outside, link)`），要的是"unlink 作用于链接本身、不递归进别人的树"。软链形红在 `:73`——被**宿主的** `/varlink`（F1-a）拒了，`:77` 那枚被测链接**还没种出来**。⇒ 换根后它测的仍是自己种的链接；**而且这枚必须换两枚根**（`outside` 与 `root` 是两个 `t.TempDir()`），只换一枚会连 setup 都过不了。 |
| 10 | `TestPOSIXMissingFileIsNotAnError` | `:99` `t.Errorf("RemoveUnlinked of a non-existent entry: %v", err)` | **甲** | 要的是"不存在的条目幂等、不算错"。软链形它拿到 F1-b（`entry is a link to something else`），因为宿主的链接在它的不存在路径的祖先上。⇒ 这枚的红恰恰**证明底线在链接上是拒的**：换根后"未解析形被拒"这条**仍然被别的用例钉着**（`seam_probe_root_125_other_test.go` 与票 113 族），不是被本批抹掉。 |

### 2.3 `ancestor_separator_108` / `placement_leaf_118` / `placement_symlink_113`（均 `!windows`）—— 7 枚，全甲类

| # | 用例 | 断言原文 | 类 | 本枚理由 |
|---|---|---|---|---|
| 11 | `TestAC2POSIXDoesNotFoldABackslashIntoASeparator`（108:94） | `:112` `t.Errorf("AC#2 RED: POSIX folded the backslash in %q into a separator and refused the tree's own file: %v", …)`；`:115` `the stray inside root/%s should have been removed`；`:117` `assertStillThere108(…)` | **甲** | 被测对象是 `pathPieces` 在非 Windows 上**只切 `/`**。红串（F1-b）逐字写着"…through the link at **/varlink**"，而它种的链接是 `root/a` ⇒ 拒因不是 `a\b` 那条例外。⚠ 这枚两半都红了（`:115` 说明文件**确实没被删**）⇒ 换根后"删得掉、且不误删别人树"两半才同时有分母。 |
| 12 | `TestAC4POSIXFloorAnswersInsideTheNamedTree`（108:152） | `:159` `t.Fatalf("PrivateDirAll(%s): %v", dir, err)`；`:169` `AC#4 RED: the floor answered %q for %q, i.e. it rewrote the spelling…`；`:172` `AC#4 RED: the answer %q is not inside the named tree %q` | **甲**（R-119-9 专项核过） | ⚠ 本批最需要盯 R-119-9 的一枚：它的期望值就是**声明的那个字符串本身**（`got.String() != dir`）。⇒ 只要 `dir` 仍来自"调用方自己声明的局部变量"，纪律就还在；**红线是**不能写成 `winsec.ResolvePath(winsec.SealableRoot(...))` 之类拿被测函数再算一遍。本批按前一种改（只换 `t.TempDir()` 那一步的来路），并把它记进 §9 的复核点。另外它开头有 `:153` `PathResolverInstalled() != nil ⇒ t.Skipf`，是既有的形状无关守卫，不动。 |
| 13 | `TestAC118POSIXSealFileRefusesALinkStandingWhereTheFileWasNamed`（118:75） | `:81` `assertRefused113(t, "SealFile", link, winsec.SealFile(link))`；`:84` `the mode change followed the leaf` | **甲**（红因是**守卫救回来的**） | 红在 helper `leafLinkTo118:50-52` 的前置（F1-c）：`EvalSymlinks(link)`=`/realpriv/…` ≠ `target`=`/varlink/…`。⇒ **这枚前置正是防 §1.3(b) 那种假绿的闸**：它没被放宽、也没被绕过，它把一枚本来会"因宿主的链接而绿"的用例**改成了响亮地红**。换根后前置成立、被测的 `SealFile(link)` 才真的跑。 |
| 14 | `TestAC118POSIXPrivateFileRefusesALinkStandingWhereTheFileWasNamed`（118:96） | `:107` `assertRefused113(t, "PrivateFile", link, err)`；`:114` `PrivateFile(%q) replaced the victim's bytes` | **甲** | 同 #13。 |
| 15 | `TestAC3POSIXSealFileStillNarrowsAPlainFileInsideTheNamedTree`（113:253） | `:269` `AC#3 RED: a plain file inside the named tree was refused`；`:281` `AC#3 RED: the seal did not narrow anything, mode is %o` | **甲** | AC#3 的反向半（"全都拒"的假修法是它打红）。软链形它被 `/varlink` 拒 ⇒ 这条**反向腿在软链形完全没有分母**。 |
| 16 | `TestAC3POSIXSealStillWorksNextToAndThroughRealDirectoriesAndLinks`（113:287） | `:298` `AC#3 RED: PrivateDirAll refused its own tree`；`:311` `PrivateFile refused a plain path inside the named tree`；`:317` `SealDir refused its own directory`；`:326` `…a directory that merely has a link among its children` | **甲** | 它种的链接是 `root/not-an-ancestor` 与 `deep/dangling`（**兄弟/子节点**位置），红在 `:298` 宿主的链接。⇒ 换根后"旁边有链接也要能密封"这条仍然被测，且**只有**换根才测得到。 |
| 17 | `TestAC3POSIXSealDoesNotFoldABackslashIntoASeparator`（113:335） | `:355` `AC#3 RED: POSIX folded the backslash in %q into a separator and refused the tree's own file`；`:357` `%s is %v (%v), want -rw-------` | **甲** | 同 #11 的另一半（`root/a` 是链接 + 真目录 `a\b`）。红串又是 F1-a 的 `/varlink`。 |

### 2.4 乙类＝0：这一档为什么空，以及它意味着什么

派单要求乙类登记"保留红是有意的"并写清"转绿会丢掉什么检测力"。**本批没有这一档**，理由不是没找，是逐枚断言原文都对不上：

- 票面 AC#2 那句"或改成显式钉住『就是要拿未解析的根试底线拒绝』"在 winsec 里**已经有专门的用例在做**，
  而且**不是这 17 枚**：`seam_probe_root_125_other_test.go`（3 枚 `TestAC2POSIXSeam…125`，软链形逐名 `t.Skipf`，见 §5.4 对账）
  与票 113 族 `placement_symlink_113_other_test.go` 的 `assertRefused113` 拒绝腿。
  ⇒ **要保住的"未解析那一形"另有其人，本批 0 枚动它们**（批次 2 保留 `internal/config/c26_seam_posix_125_test.go` 是先例，本批同处理）。
- 17 枚里名字带"Refuses"的只有 #1/#6/#13/#14/#16，逐字回看软链形**实际拿到**的字符串（§1.3 F1）：
  没有一枚拿到"它自己要的那句拒"——#1/#6 要 `ErrNotSealable` 拿到 `not provably resolved`；#13/#14 根本没跑到断言。
  ⇒ 与 AC#2a 账上"拒绝腿 0"一致（那本账是**断言方向**口径，本批是**改动面＋纪律**口径，两口径在此对上）。

⇒ **裁决：17 枚全部该转（甲类），本批要动码。** 但**形状不许照抄批次 1/2**：不接 `proc.SealableRoot`，
改按票 125 先例在包内自走。理由与实现见 §8。

## 3. 判据① 逐枚转绿且点名（四数只从 `-v` 量）

四台＝改前/改后 × 软链/普通，全部 `golang:1.27` 容器内 `go1.27.1 linux/amd64`、`./internal/winsec/` `-count=1 -v`
（`BASE-L/BASE-P/POST-L/POST-P.log`；每台进容器先 `ls -l /src/go.mod`=883 + `md5sum resolve.go`=`b6876a5e…`）。
**这 17 枚里 11 枚在宿主 Windows 上没有分母**（`!windows` 标签），所以本表**只有容器读数**，宿主那一台见 §7.2 的 windows 腿。

| # | 用例名 | PRE-L 软链 | POST-L 软链 | PRE-P 普通 | POST-P 普通 |
|---|---|---|---|---|---|
| 1 | `TestAC118POSIXPrivateFileRefusesALinkStandingWhereTheFileWasNamed` | FAIL | PASS | PASS | PASS |
| 2 | `TestAC118POSIXSealFileRefusesALinkStandingWhereTheFileWasNamed` | FAIL | PASS | PASS | PASS |
| 3 | `TestAC2POSIXDoesNotFoldABackslashIntoASeparator` | FAIL | PASS | PASS | PASS |
| 4 | `TestAC3POSIXSealDoesNotFoldABackslashIntoASeparator` | FAIL | PASS | PASS | PASS |
| 5 | `TestAC3POSIXSealFileStillNarrowsAPlainFileInsideTheNamedTree` | FAIL | PASS | PASS | PASS |
| 6 | `TestAC3POSIXSealStillWorksNextToAndThroughRealDirectoriesAndLinks` | FAIL | PASS | PASS | PASS |
| 7 | `TestAC4POSIXFloorAnswersInsideTheNamedTree` | FAIL | PASS | PASS | PASS |
| 8 | `TestAC5FailedSealRefusesTheWrite` | FAIL | PASS | PASS | PASS |
| 9 | `TestAC5FailureIsNotSwallowedByTheHappyPath` | FAIL | PASS | PASS | PASS |
| 10 | `TestPOSIXMissingFileIsNotAnError` | FAIL | PASS | PASS | PASS |
| 11 | `TestPOSIXPrivateDirIsReally0700` | FAIL | PASS | PASS | PASS |
| 12 | `TestPOSIXPrivateFileIsReally0600` | FAIL | PASS | PASS | PASS |
| 13 | `TestPOSIXSymlinkAtArtifactPositionIsNotRecursed` | FAIL | PASS | PASS | PASS |
| 14 | `TestAC5FailedSealRefusesTheWrite/directory_chain` | FAIL | PASS | PASS | PASS |
| 15 | `TestAC5FailedSealRefusesTheWrite/exclusive_artifact` | FAIL | PASS | PASS | PASS |
| 16 | `TestAC5FailedSealRefusesTheWrite/replacement_write` | FAIL | PASS | PASS | PASS |
| 17 | `TestAC5FailedSealRefusesTheWrite/seal_file_and_dir_direct` | FAIL | PASS | PASS | PASS |

表由日志程序化生成、非手抄：对 17 枚逐名在四份 `-v` 日志中各查一次 `^ *--- (PASS|FAIL|SKIP): <名> ` 命中。
**17 行 × 4 列无一格变差**（没有任何一枚从 PASS 变成 FAIL/SKIP）。

包级四数（`〔容器/软链形〕`，`-count=1 -v`）：

| 台 | RUN | 顶层 PASS | 顶层 FAIL | 顶层 SKIP | 子测 PASS | 子测 FAIL | rc |
|---|---|---|---|---|---|---|---|
| PRE-L | 45 | 14 | **13** | 3 | 11 | **4** | 1 |
| POST-L | 45 | 27 | **0** | 3 | 15 | **0** | **0** |

⇒ 本批软链形红名 **17 → 0**（顶层 13→0、子测试 4→0），且顶层 PASS 恰 +13、子测 PASS 恰 +4
⇒ **变绿的是那 17 枚本身，不是分母移动**。`-v` 日志里点名宿主的 `refusing to seal` / `not provably resolved` 行数：
PRE-L **20** 行 → POST-L **0** 行。

## 4. 判据② 普通形一枚都不许多红，且判定分支一字未动

**普通形八个数改前改后逐数相同**（`〔容器/普通形〕`，`-count=1 -v`）：
PRE-P `RUN=52 顶层 30/0/0 子测 22/0` ＝ POST-P `RUN=52 顶层 30/0/0 子测 22/0`，`rc` 两枚都 0
⇒ **普通形一枚都没多红，也没有一枚变绿**（本来就全绿）。

**`git diff` 的删除侧全文**（`git diff 8ced405^ 8ced405 -- internal/winsec/` 里**全部** `-` 行，15 行，一字不差）：

```
-	dir = filepath.Join(t.TempDir(), name)
-	root := filepath.Join(t.TempDir(), "root")
-	root := filepath.Join(t.TempDir(), "data")
-	root := filepath.Join(t.TempDir(), "root")
-	f.dir = filepath.Join(t.TempDir(), name)
-	root := filepath.Join(t.TempDir(), "data")
-	root := filepath.Join(t.TempDir(), "data")
-	root := filepath.Join(t.TempDir(), "root")
-	dir := t.TempDir()
-	dir := t.TempDir()
-	dir := t.TempDir()
-	root := filepath.Join(t.TempDir(), "data", "artifacts")
-	outside := filepath.Join(t.TempDir(), "someone-elses-tree")
-	root := filepath.Join(t.TempDir(), "data")
-	if err := RemoveUnlinked(filepath.Join(t.TempDir(), "gone")); err != nil {
```

⇒ **15 行删除全部是"哪份根递给底线"这一个表达式**；`t.Errorf` / `t.Fatalf` / `assertRefused113` / `assertStillThere108` /
`errors.Is(…)` / 模式位比较 / 阈值 **0 命中**。新增侧除两枚新文件外只有 15 处同名替换（`git diff --numstat` ＝ 15 增 15 删，逐文件 3/1/4/2/5）。
⇒ **票面 AC#3 成立**：底线码与判定分支都没碰，`internal/winsec/resolve.go` 零 hunk、`internal/proc/**` 零 hunk、
`internal/risk/**` 零 hunk、前两批已交文件零回退。

## 5. 判据③ 两形各一枚，RUN/SKIP 差逐包解释

受影响包只有一枚 `./internal/winsec/`（改动面见 §8，可盘上核）。两形各一枚 `-count=1 -v` 的差：

| 项 | 软链形 | 普通形 | 差 | 逐名解释 |
|---|---|---|---|---|
| RUN | 45 | 52 | **7** | 全部是票 125 三枚自拒探针的 7 枚子测试（父用例 `t.Skipf` ⇒ 子测试从未被创建），逐名见下 |
| 顶层 SKIP | 3 | 0 | 3 | 上述三枚父用例，软链形 `t.Skipf`（`seam_probe_root_125_other_test.go:120/223/276`），普通形跑 |
| 顶层合计 | 27P+0F+3S=30 | 30P+0F+0S=30 | **0** | 顶层分母两形**一模一样** |
| 子测合计 | 15P+0F | 22P+0F | 7 | 同 RUN 差 |

那 7 枚子测试逐名：`TestAC2POSIXSeamProbeShapesAreBuiltOnAResolvedRoot125/{control_plain_temp,measured_symlink_spelled_temp}`、
`TestAC2POSIXSeamGuardStillRefusesEveryHostileShape125/{control_plain_temp,measured_symlink_spelled_temp}`、
`TestAC2POSIXSeamAcceptsTheHonestPOSIXAnswer125/{control_plain_temp,measured_symlink_spelled_temp,control_unresolved_root_still_refused_by_the_floor_itself}`。

⚠ **派单点名的那条（"软链形自己会缩小分母，'不再红'要分清变绿还是被跳过"）逐枚排掉**：

- **SKIP 名册改前改后逐名相同**：`diff` PRE-L 与 POST-L 的 `--- SKIP` 名单 ＝ 空（就是那 3 枚 125 探针）
  ⇒ 本批**没有新增任何一枚 SKIP**；
- §3 表里 17 枚旁写的是 `PASS` 不是 `SKIP`，且 POST-L 的用例名集合与 PRE-P 逐名相等（只差那 7 枚 125 子测试）；
- 那 3 枚探针 + 7 枚子测试**一枚未动**（既不在 17 枚分母里，也不在 §8 改动面里）；
  批次 2 留下的跑法账（**每形一枚新容器 + `-timeout` 够长**）本批照做，两形分别用了独立容器。

⇒ 形状自带的缩小（7 枚子测试）改前改后**同数同因**，没有被算进"本批清掉的红"。

## 6. 判据④ "另一种拒"逐枚实拿字符串

探针台 `PROBE-L`（`/d/tmp/wisp124-2b3b-probe/`，仓外）：只在 `private_fail_test.go` 加 **4 行 `t.Logf`**
（`grep -c PROBE124` ＝ 4），不改任何断言 ⇒ 四数与 POST-L **逐数相同**（`RUN=45 顶层 27/0/3 子测 15/0`、rc=0），
先证明探针本身没移动分母，再拿字符串。

- **#1–#5（AC5 注入族，它自己要的是 `ErrNotSealable`）**——四枚都拿回**注入的那一句**，且 `IsErrNotSealable=true`：
  `winsec: cannot apply a private security descriptor: /realpriv/w124tmp/TestAC5FailedSealRefusesTheWrite…/001/{artifact.txt,blob.bin,a,plain.txt}: injected: descriptor could not be applied`
  ⇒ "未解析根"那句从这五枚的拒因里**消失**，注入腿重新成为操作性拒因（转前它拿到的是 `not provably resolved`，见 §1.3/F1-a）。
- **#13/#14（AC118 leaf-link，它自己要的是"我种的叶子链接被拒"）**——转后逐字点名**自己的叶子**：
  `AC#1 SealFile("/realpriv/…/002/root/artifact.txt") -> err=winsec: refusing to seal … reaches it through the link at /realpriv/…/002/root/artifact.txt, which is not the tree this call names`
  （`AC#1 PrivateFile(…)` 同形）⇒ 链接名＝用例自己 `os.Symlink` 出来的那枚 `artifact.txt`，不再是宿主的 `/varlink`；
  并且 `leafLinkTo118:50-52` 的 `filepath.EvalSymlinks` 前置**没被放宽也没被绕过**，它现在通过是因为形状真的对了。
- **#3/#4/#5/#6/#7/#10/#11/#12/#13(模式族)/#15/#16/#17（反向腿，要的是"别拒"）**——转后拿到的是**成功**
  （`SealFile`/`PrivateFile`/`PrivateDirAll`/`SealDir`/`RemoveUnlinked` 返回 nil，模式位回读 0600/0700 成立）。
- ⚠ **反扫**：`POST-L.log` 与 `GATE-L.log` 里 `refusing to seal`、`not provably resolved`、`/varlink` 三串命中 **0 行**
  ⇒ 没有一枚靠"换个理由拒"或"换个理由跳"过断言。

## 7. 判据⑤ 变异自证三态 ＋ 票面 AC#5 门禁

### 7.1 变异三态（票面 AC#4）

`MUTATION-124-2B3B`：把 `tempdir_resolved_124_other_test.go` 的函数体退回原样交 `t.TempDir()`（＝拆掉这一层解析＝旧实现）。
**先证落地、再读数**：

- `grep -n` 落地：`…/wisp124-2b3b-mut/internal/winsec/tempdir_resolved_124_other_test.go:35` 命中 `MUTATION-124-2B3B`、
  `:36` 命中 `return t.TempDir()`，**`return resolveProbeRoot` 命中 0**；
- 容器原生 `go build ./internal/winsec/` **rc=0** ＋ `go vet ./internal/winsec/` **rc=0**（`MUT-BUILD.log`）**之后**才读红名。

| 态 | 形状 | RUN | 顶层 PASS | 顶层 FAIL | 顶层 SKIP | 子测 PASS | 子测 FAIL | rc |
|---|---|---|---|---|---|---|---|---|
| PRE（旧实现） | 软链 | 45 | 14 | 13 | 3 | 11 | 4 | 1 |
| POST（转后） | 软链 | 45 | 27 | **0** | 3 | 15 | **0** | **0** |
| MUT（退回旧实现） | 软链 | 45 | 14 | **13** | 3 | 11 | **4** | **1** |
| PRE / POST / MUT | 普通 | 52 | 30 | 0 | 0 | 22 | 0 | 0 |

⇒ **三态齐**：拆掉解析 ⇒ 软链形**重新 17 枚红**，且 MUT-L 的红名与 PRE-L 的红名 **`diff` 逐名完全相同**（17 枚）；
同一发变异在普通形 **0 红**、八个数与 PRE-P/POST-P **逐数相同**
⇒ 这层解析**只在软链形起作用**，不是把普通形一起喂绿。

### 7.2 门禁（票面 AC#5，本批受影响包＝`./internal/winsec/`）

`-count=2 -v` 两形各一枚**新**容器（`〔容器/linux/amd64〕`，`-timeout 30m`）：

| 形 | RUN | 顶层 PASS | 顶层 FAIL | 顶层 SKIP | 子测 PASS | 子测 FAIL | rc |
|---|---|---|---|---|---|---|---|
| GATE-L 软链 | 90 | 54 | **0** | 6 | 30 | **0** | 0 |
| GATE-P 普通 | 104 | 60 | **0** | 0 | 44 | **0** | 0 |

⇒ 每枚数都是 §3 那发 `-count=1` 的**正好 2 倍**（无缓存、无 flake、与原读数无分歧）；GATE-L 的 6 枚 SKIP＝那 3 枚 125 探针 ×2、
名册逐名相同；GATE-P 的 SKIP 名册为空；两形机制字串 `not provably resolved` / `refusing to seal` / `/varlink` 命中 **0**。

**宿主那半边（票面 AC#5 的 windows 腿，`〔宿主/windows〕`）**：`private_fail_test.go` 的 6 枚是本批唯一在宿主有分母的用例，
它们在 Windows 上走 `tempdir_resolved_124_windows_test.go` 的字面 no-op ⇒ 改动面对宿主是零 hunk 的根表达式替换。
宿主 `go vet ./internal/winsec/` **rc=0**（类型读数，见 §7.3）。

静态与仪器：

- `gofmt -l internal/winsec`（**整包**）宿主 **0 行**；`gofumpt -l internal/winsec`（**整包**）**0 行**，
  版本 **`v0.7.0 (go1.27.1)`** —— CI 没钉版本，这里钉明我用的这枚；真跑，非"未跑"。
- `d22scan` 纯净快照（正确调用式 `cd tools/d22scan && go run . -root "<仓绝对路径>"`，两发都 rc=0 clean）：
  - PRE 快照（`4ea0db2` → `D:\tmp\wisp124-2b3b`）八 scope `203/22/40/18/16/40/397/38`；
  - POST 快照（`8ced405` → `D:\tmp\wisp124-2b3b-post`）正向对照 `examined 225 production Go files`，八 scope `203/22/40/18/16/40/`**`400`**`/38`。
  ⇒ **一枚 scope 都不降**；唯一变化 `ban #8 internal/` 397→400（+3）**逐名有出处**：本批 2 枚
     （`tempdir_resolved_124_other_test.go`、`tempdir_resolved_124_windows_test.go`）
     ＋ 兄弟 `worker-ticket124-ac2b-3a` 的 `62dda11`（`internal/agent` 26 枚，1 枚新文件）＝3。
     其余七枚与批次 2 交回的 `203/22/40/18/16/40/397/38` **逐数相同**。

### 7.3 `go vet` 双 GOOS

| 读法 | 台 | 范围 | rc | 逐行归因 |
|---|---|---|---|---|
| `go vet ./internal/winsec/` | **容器原生 linux/amd64**（`CGO_ENABLED` 默认，`8ced405` 快照） | 本批改动面 | **0** | **这才是那 11 枚 `!windows` 用例＋`tempdir_resolved_124_other_test.go` 的真类型读数**——宿主结构上到不了它们 |
| `go vet ./internal/winsec/` | 宿主 windows | 本批改动面 | **0** | 覆盖 windows 那一支 no-op helper 与无标签的 `private_fail_test.go` |
| `go vet ./...` | 宿主 **GOOS=windows** 全树 | 30 枚包 | **0** | 无输出 |
| `GOOS=linux go vet ./...` | 宿主 **交叉** 全树 | 30 枚包 | **1** | 1 个错误块，逐字归因到 `package github.com/CarlosShao/wisp/cmd/wisp → imports github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx → imports github.com/k2-fsa/sherpa-onnx-go-linux: build constraints exclude all Go files in …sherpa-onnx-go-linux@v1.13.8` |
| `go vet ./...` | 容器原生 全树 | 30 枚包 | **未取得** | 首发为整树拉模块缓存，卡在 `go: downloading …/sherpa-onnx-go-linux v1.13.8` 未收敛；按批次 1/2 的同一账，**包级 rc=0 已覆盖本批改动面**，全树那一发留给 2b-4 的终判据复算 |

⇒ 交叉那一发按派单口径**既不算破口也不算清白**：它停在外部模块的 build constraints、**到不了类型检查**（错误块里没有任何本仓
`file:line`，且 `internal/winsec` 不在归因路径上）。同一发在批次 1、批次 2 逐字相同 ⇒ 与本批无关。
本批改动面的类型读数取上面前两行的 **rc=0**。

## 8. 改动面

### 8.1 交件 commit

| commit | 内容 |
|---|---|
| `857a5fe` | `docs/evidence/s1/124-ac2b-3b-conversion.md` §0–§2（三分类，本件主产物） |
| `8ced405` | `internal/winsec/**` 动码（2 枚新文件 + 5 枚改 15 处） |
| `09edf02` | 证据 §3–§7（判据①-⑤ 与门禁读数） |

（本节与 §9/§10 是第 4 枚 commit。**AC#2 本格与 AC#5 都不翻**，见 §9。）

### 8.2 文件级

**新增（2 枚，都是 test-only，都不进生产码）**
- `internal/winsec/tempdir_resolved_124_other_test.go` —— `//go:build !windows`，`package winsec`，
  `SealableTempDirForTest124` ＝ `resolveProbeRoot(t.TempDir())`。**复用包内既有的那一次走法（`resolve.go:253`），
  没写第四份重复解析**（票 125 `R-125-2` 那本副本账不增行），**也没 import `internal/proc`**。
- `internal/winsec/tempdir_resolved_124_windows_test.go` —— `//go:build windows`，字面 no-op（宿主/CI 的 windows 腿一字不动，理由写在该文件注释里）。

**修改（5 枚，共 15 处递根点，逐枚 file:line）**

| 文件 | 处数 | 改后位置 | 覆盖的用例 |
|---|---|---|---|
| `private_fail_test.go`（无标签，`package winsec`） | 2 | `:28`、`:100` | #1–#6 |
| `private_other_test.go`（`!windows`，`package winsec`） | 5 | `:21`、`:38`、`:61`、`:71`、`:98` | #7–#10 |
| `ancestor_separator_108_other_test.go`（`!windows`，`package winsec_test`） | 3 | `:38`(helper `foreignTree108`)、`:96`、`:156` | #11、#12 |
| `placement_leaf_118_other_test.go`（`!windows`，`package winsec_test`） | 1 | `:60`(helper `ownTree118`) | #13、#14 |
| `placement_symlink_113_other_test.go`（`!windows`，`package winsec_test`） | 4 | `:53`(helper `newForeign113`)、`:254`、`:289`、`:337` | #15、#16、#17 |

⇒ 三枚 helper 里只有 **2 枚**（`newForeign113`、`foreignTree108`）被换，因为它们**同时**供给了本批用例；
`leafLinkTo118` 的 `filepath.EvalSymlinks` 前置**一字未动**（它是 §6 里那两枚的防假绿闸）。

### 8.3 禁改列与兄弟地界（逐条核过，不是声明）

- `docs/PLAN.md`、`docs/specs/**`、`internal/risk/**`、`internal/risk/pathresolver*.go`、`rules_gateway.go`、
  `tools/d22scan/**`、`allowlist.txt`、任何阈值／golden／`thresholds.go`、`frontend/**`：**零 hunk**
  （`git diff --name-only 4ea0db2 8ced405` 只出现上面这 7 枚路径，已在 §4 后逐枚核）。
- `internal/proc/**` 的 `SealableRoot` 本体：**零 hunk，且根本没 import 它**。
- `internal/winsec/resolve.go`（票面 AC#6 点名的禁改列）：**零 hunk**，只**读**了它的 `resolveProbeRoot`；
  锚点与改后 `md5sum` 同为 `b6876a5efe759f6e17434d1b50a129c3`（每台容器都自证过）。
- `internal/agent/**`：**零 hunk**（兄弟 3a 在那儿；它的 `62dda11` 在我之前落库，我没碰它的文件）。
- 前两批已交文件（memory/config/tools/llm/perm/agent/Approval 的 `tempdir_resolved_124*_test.go`）：**零回退**。
- `internal/winsec/dataroot_symlink_119_other_test.go`、`seam_probe_root_125_other_test.go`：**一字未动**
  （前者是 `R-119-9` 的守卫文件、后者是那 3 枚自拒探针）。
- 未跟踪的 `docs/reports/2026-09-23-gap-analysis-vs-oss-harnesses.md`：**未提交、未改、未删、未据它开票或改判据**
  （全程 `git status --porcelain` 里它一直单独挂着）。

### 8.4 宿主（`〔宿主/windows〕`）读数补全：本批唯一有宿主分母的那 6 枚

派单要求"每个数标在哪台、哪个构建标签"，所以宿主半边单独量（`go test -count=1 -v -run 'TestAC5…'`，无标签文件，Windows 真跑）：

| 用例 | 改前（`4ea0db2` 快照） | 改后（`8ced405` 工作树） |
|---|---|---|
| `TestAC5FailedSealRefusesTheWrite` | PASS | PASS |
| `…/exclusive_artifact`、`…/replacement_write`、`…/directory_chain`、`…/seal_file_and_dir_direct` | PASS ×4 | PASS ×4 |
| `TestAC5FailureIsNotSwallowedByTheHappyPath` | PASS | PASS |

⇒ **6/6 两态同 PASS** ＝ windows no-op 那一支按设计生效（宿主半边没被解析层碰过）。
另外 11 枚（`!windows`）宿主**连编译都不参与**（`GOOS=windows go test -list` 逐名命中 0，见 §1.1），
所以它们在宿主**没有"改前改后"可言**，只有 §3 的容器读数。

## 9. 未验证项与 `next=`

### 9.1 未验证 / 已登记不裁

1. **AC#2 本格不翻、AC#2b 不翻、AC#5 不翻**（派单明写）。票面 16:33 的终判据"软链形红名数＝0"
   要 2b-4 交回后在全树一次性复算，本批只把 `./internal/winsec/` 这一包从 17 做到 0。
2. **§1.3(b) 那 12 枚邻居的"vacuous green"是字符串对比推出来的，没做变异自证**。
   我量到的是：软链形它们 PASS 且拒因写 `/varlink`、普通形同枚 PASS 且拒因写自己种的 link，
   而 `assertRefused113` 两形都只要求 `ErrUnresolvedPath` ⇒ 同 sentinel、无法区分。
   **没有做**的是"把叶子/祖先腿拆掉（票 113 验收方的 MUT-B）看它们在软链形会不会照样绿"那一发 ⇒ 属推断，交验收方。
3. **邻居自己的 7 处 `t.TempDir()` 未换**（`ancestor_108:69,127`、`placement_symlink_113:146,169,197,213,238`），
   所以那 12 枚在软链形**仍然 vacuous**。换它们＝动票 113/108 的用例面，不属本批 17 枚分母 ⇒ 有意不做，不属漏做。
4. 交叉 `GOOS=linux go vet ./...`（宿主）按派单**既不算破口也不算清白**：停在外部模块 cgo 包的 build constraints，
   到不了类型检查。真正的 linux 类型读数是**容器原生**那一发（§9.4）。
5. `d22scan` 跑在含码的纯净快照（`8ced405`）上；**证据文件本身不在台账 scope 内**（台账只看 `internal/`、`cmd/`、
   `frontend/`、`design/`），所以 §9/§10 追加不会再动那八枚数。
6. **全树 30 枚包的分母本批未复跑**（那是 2b-4 的终判据）。本批改动面只有一枚包，且判据② 已证普通形逐数不变。
7. 变异自证只发了**一票**（拆解析层）。票面 AC#4 要求的最低量是"至少一发把新加的那层解析拆掉"⇒ 已满足；
   没有发"把 `resolveProbeRoot` 本身拆掉"那一发（它属票 125 的地界）。

### 9.2 `next=`（交编排者 / 2b-4）

1. **2b-4＝`cmd/wisp` 5 枚**（`TestProvidersProbeRecordsMeasuredThinkingFalse/True`、`TestProvidersDiscoverListsWhatTheServerServes`、
   `TestProvidersProbeUnconfiguredRefIsNotSilentlyKeyless` 四枚顶层走同一枚 `newProvidersFixture` ⇒ 一处递根覆盖 4 枚；
   加 `secret_test.go:679` 的 `TestSecretFailurePathsLogAndPrintNoPlaintext/unset_name_(no_such_blob)`）。
   ⚠ 同包另有 **28 枚两形都红**（`DPAPI` 族 + 命令面腿，`R-119-7` 那本账）⇒ **一枚都不许顺手修绿**，
   也不许因包级 `FAIL≠0` 把自己那 5 枚判成没转绿：**判据是逐枚点名、不是包级 rc**。
   ⚠ 地界按**文件级**核：`cmd/wisp/leg_dispatch_gate_133_test.go` 系票 133 在动。
2. ⚠ **形状别照抄**：`cmd/wisp` 与 2b-1/2 同为"调用方解析 OS 给的答案"，可直接用 `proc.SealableRoot`；
   **只有 `internal/winsec` 是例外**（底线自己的包，见 §2.4 与 §8.2），2b-4 不需要再走 3b 这套。
3. **全树终判据复算**（2b-4 交回后一次做，须含 3a＋3b 合并态）：纯净快照 + 两形各一枚**新**容器 +
   逐包 `-count=1 -v` 跑 AC#1 那本 30 枚包分母 ⇒ 软链形红名与账上 131 枚"可转"做差集**期望为空**。
   豁免名单逐名钉死为：票 123 那族 **三枚** 300 s 腿（`TestL1Write`、`TestLateVeto…`、`TestFSReadOnlyNeverOpensACard`）
   ＋ 两形都红 29 枚（`panel` 1 ＋ `cmd/wisp` 28）＋ 形状自带 SKIP 4 枚（`winsec` 3 ＋ `config` 1，既非红亦非绿）。
   ⇒ 本批把这 4 枚 winsec/config 侧的数字钉成了**改后仍 3 枚 SKIP、名册逐名相同**，复算时按名核即可。
4. **请一并裁 §1.3(b)/§9.1(2)(3) 那本邻居账**（12 枚 winsec POSIX 拒绝腿在软链形是"绿得没有理由"）。
   可选修法就一行形状：把那 7 处邻居根也换成 `SealableTempDirForTest124(t)`。
   代价：动的是票 113/108 的用例面、且必须补一发 MUT-B 型变异才能证明它们真的变强 ⇒ 本批按地界没做。
   不修的后果：软链形（＝macOS 真实形状）下这 12 枚**零检测力**，下一位改 `winsec` 叶子腿的人不会收到红。
5. ⚠ **本批沿用的两条跑法账值得写进每份派单**：`exit 99`（普通形不许存在 `/varlink`，批次 1 加的）
   与"每形一枚**新**容器"（批次 2 交的）；本批有一次**自己撞出来的**教训：容器内 `sh -c` 忘 `cd /src`
   ⇒ `go.mod not found`、`-count=2` 两发第一轮全部 rc=1 空跑（**没被当成绿**，因为四数是 0）。登记在此以免下游重踩。

### 9.4 容器原生 `go vet ./...`（真类型读数）

**未取得，且不是破口**：那一发（`VET-CONTAINER.log`）在首次为整树拉模块缓存时卡在
`go: downloading github.com/k2-fsa/sherpa-onnx-go-linux v1.13.8` 未收敛，本批不等它、改取**包级**那一发
——`go vet ./internal/winsec/` 容器原生 **rc=0**（读数与逐行归因在 §7.3）。
本批改动面只有 `internal/winsec` 一枚包，包级 rc=0 就是它的类型结论；**全树那一发留给 2b-4 的终判据复算顺手补**。

## 10. 临时件（只建不删）

| 路径 | 用途 |
|---|---|
| `/d/tmp/wisp124-2b3b/` | 锚点 `4ea0db2` 纯净快照（PRE 台；宿主改前读数也在这台上量） |
| `/d/tmp/wisp124-2b3b-post/` | 改件 `8ced405` 纯净快照（POST 台 + `-count=2` 门禁台 + d22scan POST） |
| `/d/tmp/wisp124-2b3b-mut/` | `MUTATION-124-2B3B` 台（判据⑤三态） |
| `/d/tmp/wisp124-2b3b-probe/` | 探针台（判据④"另一种拒"实拿，只加 4 行 `t.Logf`） |
| `/d/tmp/wisp124-2b3b-cycle/` | 环引用探针（§1.3(a)：证 `package winsec` import `proc` 无环、vet rc=0）——**在仓外** |
| `/d/tmp/wisp124-2b3b-gocache/` | 容器共享 `GOMODCACHE`（被上面各台复用） |
| `/d/tmp/wisp124-2b3b-run.sh`／`-run-lf.sh` | 跑法脚本（CRLF 原件／LF 副本，容器执行 LF 那枚） |
| `/d/tmp/wisp124-2b3b-logs/` | 全部 `-v` 日志：`BASE-L/P`、`POST-L/P`、`MUT-L/P`、`PROBE-L`、`GATE-L/P`、`MUT-BUILD`、`VET-{CONTAINER,HOST-WIN,HOST-LINUXX}`、`HOST-LIST`、`PRE-L-red.txt`、`MUT-L-red.txt` |
| `/d/tmp/wisp124-2b3b-frag-3-7.md` | 本件 §3–§7 的写作件（`cat >>` 用） |

被本件引用的目录谁都别删——它们是这份裁决表可复算的凭据。
