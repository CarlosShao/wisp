# 342-a1 census（票 342 `AC#0` 三形代价表｜只读普查腿｜⛔ 选形、⛔ 判语、⛔ 落任何一形）

腿＝`342-a1`（票 342 第一枚腿）；锚＝起手现量 `HEAD=c596df55`（本腿所有行号＝`git show HEAD:<path>`，⛔ 在飞工作树；`internal/winsec` 全部 40 枚 blob 在 `98dbf20a`⇄`c596df55`⇄现 HEAD **逐枚同串**⇒本票产码/测试在数小时内未动，行号仍可用）。
台面＝`git archive c596df55`（八件套＋`docs/`）导出树 `D:\tmp\wisp\20261018\342a1\{treeA,treeB,mut-*}`；POSIX 档＝`golang:1.27` 容器（`CGO_ENABLED=0`、`GOPROXY=off`、宿主 module 缓存只读挂载、`MSYS_NO_PATHCONV=1`）；Windows 档＝本机 `go1.27.1 windows/amd64`。逐字读数在 `logs/`，尺在 `scripts/`＋`probe/`。

## 0 起手
- `git log -1`＝`98dbf20a`（★起手时刻；现 HEAD 已进到 `062dbe44`＝别腿三笔，`internal/winsec` 零动）。
- `git status --porcelain -- internal/winsec cmd/wisp internal`＝**空（0 行）**。
- `git log --all --oneline --grep='342-a1'`＝**1 笔＝`98dbf20a`**，那是编排者的**派单存档**（正文里写着腿名，`git show --name-only` 只含 `probes/orch/20261018/txt/*`）⇒ 尺按「正文也算命中」＝1；按「腿的提交」＝**0**。第二把尺：`git log --all --oneline -- .scratch/wisp/probes/342`＝**空**、`.scratch/wisp/probes/342` 起手不存在 ⇒ 自证＝本腿是第一枚 342 腿。
- 票框现量：**未勾 7／已勾 0**（尺见 `logs/00-h0-open.txt` 一族）。

## 1 票面「现量」六行逐枚重打（全部从具名提交树取）
| 行 | 票面原读数 | 本腿现读（`c596df55`） | 判定 |
|---|---|---|---|
| 那次 open 的现场 | `winsec.go:149` 逐字 `f, err := os.OpenFile(resolved, flags\|oNoFollowPlatform, perm)`；`:146`＝`preOpenLeafGuard(resolved)` | **同、行号未漂**：`:146`/`:149` | 同 |
| 产码注释自陈段 | `winsec.go:118-138` | **同、未漂**：起句 `:118`，`ANCESTOR` 句 `:128`，`This window is not claimed closed anywhere.` `:133`，Windows 句 `:134`/`:138`；逐字全部对得上（`logs/01-h1-rows123.txt`） | 同 |
| 敞着的调用者 | `:87`＝覆盖写旗标；`:97`＝`O_EXCL` | **行号漂了 −9**：覆盖写＝`winsec.go:96`，独占＝`winsec.go:106`（`:87`/`:97` 是旧锚） | **异⇒具名更正** |
| 生产调用者枚数 | `PrivateFile(` 5＋`PrivateFileExclusive(` 2 | **枚数同＝5＋2**，逐枚 `file:line` 见 H2 | 同（枚数）；本腿逐枚点名 |
| Windows 更难 | `322-a1`：`FILE_FLAG_OPEN_REPARSE_POINT` 只对叶子有效 | **本腿自己复现了**（不是引腿报）：祖先＝junction 时 `anc+REPARSE` 报 `reparse=true`（index 不同），同一祖先下的 `anc/blob.bin+REPARSE` 报 `reparse=false`＝句柄命名**目标文件**（`logs/12-host-probe-windows-rows-run2.log`） | 同，且已复现 |
| 后果形状 | 祖先被换成链接 ⇒ 那次 `O_TRUNC` 跟随它清空树外真文件 | **两平台都现量到**：POSIX 真入口 1500 发 `damage=106`（`flips=358`）；Windows 本机 BASE-W `damage=897/1500`（`flips=132`）。★另抓到**第二枚损伤形状**（票面没写）：拒绝路径自己的 `_ = os.Remove(resolved)`（`winsec.go:160`）沿祖先链接**删掉树外对象**（A3 件：`unlink err=<nil>` 且 witness `no such file`） | 同＋一枚新形状上报 |

## 2 生产调用者名册（7 枚逐枚；「依赖覆盖写吗」＝AC#2 的前置料）
| # | `file:line`（现跑） | 形状 | 依赖覆盖写吗 | 凭据（调用点上下文） |
|---|---|---|---|---|
| 1 | `internal/secret/store.go:79` | `PrivateFile` | **是**（唯一一枚"重写自己已有的路径"） | `Store.Store` 文档行 `:59-63` 逐字 "creating or **overwriting** it"；`322-a1 4` 同判 |
| 2 | `internal/config/migrate.go:93` | `PrivateFile` | **是** | 注释 `:89-92` 逐字 "because this call **rewrites an existing backup** rather than skipping it, PrivateFile is also the **repair path** for a wide backup left by an older build" |
| 3 | `internal/secret/migrate.go:169` | `PrivateFile` | 否（`os.Stat` 前置 `:168`，存在即走 `SealFile` 修复支 `:174`） | 只在 `fs.ErrNotExist` 分支调用 |
| 4 | `internal/secret/migrate.go:189` | `PrivateFile` | **是** | pre-rename temp（`:184` 固定名 `tmpPath`），崩溃残留下次必须被覆盖 |
| 5 | `internal/agent/spill.go:166` | `PrivateFile` | **是** | `.retry` temp＋rename（D31 last-writer-wins，`:165-172`） |
| 6 | `internal/agent/spill.go:343` | `PrivateFileExclusive` | 否（要的就是 `fs.ErrExist`） | `writeFileExclusive` 文档 `:331-337` 逐字 "FAILING with fs.ErrExist" |
| 7 | `internal/memory/artifacts.go:75` | `PrivateFileExclusive` | 否 | 文档 `:67-70` 逐字 "A collision is `fs.ErrExist`, **not an overwrite**" |
⇒ **4 枚依赖覆盖写**（1/2/4/5）。★票 322 已把这四枚钉成 `openleaf_overwrite_322_test.go:36`（两平台同件；`:70` 红句＝"refused an existing name"、`:79` mode 收窄）⇒ **AC#2 今天有钉**（322-a1 那句"没有任何一枚现有钉钉住"已被 322-r1 的这批件推翻，具名更正）。

## 3 三形代价表（甲／乙／丙；每形＝改动面｜POSIX 实测｜Windows 实测｜顶红名册）
> 甲/乙的读数＝**假想落点**：在导出树里新增一枚 landing 件＋把 `winsec.go:149` 那一行改调 landing 函数（`1/1` 行改动＋`41`/`34`/`53`/`47` 行新件，尺＝`logs/31-landing-diffstats.txt`；落刀前后锚命中各验一次＝旧 1→0、新 0→1，`scripts/`＋`logs/`）。⛔ 跟踪文件被动、⛔ `-overlay`。

**甲＝POSIX 逐组件 `openat(O_NOFOLLOW)` 下探到父目录、叶子在钉住的 fd 上开**
- 改动面：新件 `internal/winsec/<landing>_linux.go`（41 行，`syscall.Open("/")+Openat` 链）；`winsec.go:149`（1/1）；`pinBase` 同款卷根口径（`winsec_other.go:232`）；**Windows 侧必须另写一枚**（下条）。
- POSIX 实测：同形假想落点的定向竞态 `JIA damage=0 / flips=362 / flipErrs=142`（基线同尺 `BASE damage=242`、`A2 真入口 damage=106`）⇒ 截断那一支被关掉。
- ★**但生产路由上甲没有把伤害关掉**：`CENSUS A2`（驱动产码 `PrivateFile`，⛔ 只换 open）落点后 `damage=66/1500`、`flips=312`，逐枚归类＝**损伤全部来自拒绝路径的 `os.Remove(resolved)`（`winsec.go:160`）**，不是截断（`leg=nil` 276、floor 960、gate 3、seal-refusal 96＋descent-noise 165）。⇒ 甲只关 AC 的一半，**第二枚损伤形状要另防**（`firstLinkAncestor` 或把 unlink 改成 fd-relative）。
- Windows 实测：甲＝`FILE_FLAG_OPEN_REPARSE_POINT` 逐前缀复探＋叶子按名开（`winsec_windows.go:374` 同一枚调用 widened）。JIA-W `damage=0 / flips=72`（⛔ `PROD-W damage=3` 的残余是 unlink 形，见上）⇒ **判据能关，仍是 check-then-act**：前缀复探与 `os.OpenFile` 之间可换手，平台**没有**逐组件 open-at-time 的写形（W2 行＋`winsec_windows.go:361-369` 自陈）。成本＝每次 open 多 `2x` 枚 `CreateFile`（x＝组件数）＋`GetFileInformationByHandle`。
- 顶红名册：**实测 0 红**（`./internal/winsec/` 全档：容器基线 `RUN=71 PASS=43 FAIL=0 SKIP=0` ⇄ 甲落点**逐字同号**；本机基线 `RUN=110 PASS=61 FAIL=0 SKIP=0` ⇄ 甲落点**逐字同号**）。★误形正控：`yinofb`（钉住父目录但叶子带 `O_EXCL`）⇒ **4 枚真红**：`openleaf_322_other.go:200`（AC#1 RED damaged）／`openleaf_overwrite_322:70`（refused an existing name）／`sealhandle_privateroute_120:…` 两枚 ⇒ 尺敏感，且证明覆盖写那枚钉今天咬得住。
- 可达性代价（现量，⛔ 推断）：本包产码今天 `go build` 在 `openbsd/netbsd/solaris/plan9` **已 rc=1**（`fstypeOf` 用 `syscall.Statfs*`；`winsec_other.go:309-315`），`darwin/freebsd/dragonfly` rc=0；`syscall.Openat` 只在 linux ⇒ **甲不需要为它们新增断口，但也不能声称覆盖 unix 全族**。

**乙＝复用本包为密封做过的 handle pinning 把父目录钉住**
- 改动面：POSIX 侧 `winsec.go:149`（1/1）＋**零枚新导出名**：`pinPlacement(string)(*os.File,error)`（`winsec_other.go:149`）／`sameObjectPinned`（`:214`）／`chmodPinned`（`:255`）／`pinBase`（`:232`）／`sealHandle`（`:287`）**同包非导出、签名对得上**（`f.Fd()` ⇄ `Openat`）。★但"复用"只在 !windows 成立：`winsec_windows.go` 里**没有这一族**（`sealHandle`＝`applyDescriptor(f.Name(), false)`，`:413`；`pinPlacement` 带 `//go:build !windows`）⇒ 乙在 Windows＝**新写**（唯一 act-time 形＝`os.Root` 链）。
- POSIX 实测：假想落点 `YIa damage=0 / flips=506`（同尺基线 242）；生产路由 `damage=59`、`flips=365`＝**同一枚 unlink 残余**（形状与甲一致）。落点色＝容器 `71/43/0/0` 逐字同基线；顶红 0。
- Windows 实测：`os.Root` 链（`OpenRoot(卷)`→逐组件 `OpenRoot`→叶在钉住的 child 上开）＝`YI-W damage=0 / flips=89 / flipErrs=1`，逐发归类：`path escapes from parent` 1312／`not found` 134／`too many levels` 3／`nil` 45／`Access is denied` 6。D 行：Windows `os.Root` 对**junction 子件直接拒**（绝对目标＋树内相对目标都拒＝`D2/D3/D4`），但对**树内相对 FILE 符号链接照旧跟随**（`W-ROOTSYM followed=true`）⇒ Windows 的乙**钉祖先、不钉叶子**，叶子那一支仍要 `sameObjectPinned` 同款复验（今天 Windows ⛔ 没有它）。落点色＝本机 `110/61/0/0` 逐字同基线；顶红 0。
- 常驻仪器（有）：`posix_seal_window_120_other_test.go`（三枚，⚠ 只驱动 `SealFile`）、`sealhandle_identity/privateroute/window_120_other_test.go`（120-r3 那批，驱动 `sealHandle`/`PrivateFile` 的**叶子**身份）、`placement_symlink_113`。⇒ **这族今天有钉，但一枚都量不到祖先那一支**（见必答 (vi)）。

**丙＝只诚实报告（`f.Stat()` ⇄ 逐组件 `Lstat` 复验，截断已发生后报错）**
- 改动面：`winsec.go:149` 之后插复验（≈3-8 行）＋逐组件 `Lstat` 尺；POSIX 侧这条**今天已在跑**（`sealHandle` 的 `:300-305` 两道身份判据就是它），Windows 侧今天⛔ 没有（`winsec_windows.go:413`）。
- POSIX 实测：`YIb damage=589 / flips=748`（同尺 BASE 242／甲 0／乙 0）⇒ **报告了损伤、没防住损伤**＝⛔ 防住形。Windows：`PROD-W damage=3/1500、flips=99`，残余全来自 unlink 形；截断报告今天已部分在场（`preOpenLeafGuard` 的"未答的问题即拒绝"），祖先一支⛔ 无人复验。
- 顶红名册：**0 红**（丙不改任何现有承诺；实测＝`A2`/`PROD-W` 这些件在落点上照旧绿，红只出现在"窗口被打开"那侧）。票面与注释若选丙 ⇒ 一律⛔ 出现"窗口已关"（票面 `AC#5`）。

## 4 六枚在册钉的断言名册（尺＝逐枚 grep HEAD blob 快照；0 命中处都带正控）
| 件 | 承重断言（现跑行号） | 它调 `privateFile` 吗 | 哪一形会打红它＋红句形状 |
|---|---|---|---|
| `placement_symlink_113_other_test.go` | `:277` 枚经 `:286` 调 `PrivateFile`；`assertRefused113` `:189`，硬判据 `:197`（`ErrUnresolvedPath`）＋`:201`（`refusalCreditsLink137` `:160`）；反半 `:353`/`:376-377`（plain path 必须写成） | **是（2 处）** | 甲/乙把叶子位置的裸 `*PathError` 直送 ⇒ `:197` 红句 "refused with …, which does not name ErrUnresolvedPath"＋`:201` "…does not credit the link…"；"逢名就拒" ⇒ `:377` "refused a plain path inside the named tree"。实测落点 0 红。 |
| `ancestor_separator_108_other_test.go` | 三枚 `RemoveUnlinked`（`:91`/`:132`/`:158`，红句 `:103`/`:169`）＋`:181`（`PrivateDirAll:187`／`ResolvePath:190`/`:206`） | 否（命中 0；正控＝同尺数别处调用 **12** 枚） | 三形都碰不到 ⇒ 0 红；会红的改法＝动 `pathPieces`/`firstLinkAncestor`（⛔ 本票射程）。 |
| `dataroot_symlink_119_other_test.go` | `:139`/`:165`/`:224`（`PrivateDirAll`，红句 `:240`）＋`:259` 枚（`SealFile:287`、`PrivateFile:300`）＋`:336` | **是（1 处）** | 过度拒绝形（逢祖先链接就拒）⇒ `:300` 或 `:267-268` "PrivateDirAll refused its own tree" 红。实测落点 0 红。⚠ 它拦不住乙的承诺退化（只写一次）。 |
| `seam_probe_root_125_other_test.go` | `:119`/`:222`/`:275`（守卫件，`:58` 夹具） | 否（命中 0；正控＝`func Test`/夹具锚在场：`:58`/`:119`/`:222`/`:275`） | 三形不碰守卫 ⇒ 0 红；会红＝改 `resolveString` 的拼写口径（`:650`）。 |
| `posix_seal_window_120_other_test.go` | `:321`（驱动 `:364` `SealFile`）／`:378`／`:404`（驱动 `:479`）；红句 `:206` "AC#1 RED sample …"、`:233` "RED instrument is dead …"、`:225` flip 红 | 否（命中 0；正控＝别处调用计数 2） | ★三形**一枚都不会红，也一枚都量不到本票窗口**（只驱动 `SealFile`）⇒ 甲/乙落地必须自带驱动 `PrivateFile` 的新仪器（本腿 `probe/census342a1_probe_test.go` 就是这个形状的雏形）。 |
| 票 322 的 `openleaf_*_322_*_test.go`（四枚） | `openleaf_322_other.go:75` 枚：`:199` damaged 红、`:205` 仪器死红、`:210` landed/narrowed 红（调用 `:159`）；`openleaf_flag:35`（`:53` by-name open、`:60` 拒不得删链接）；`openleaf_guard_322_windows:91`/`:158`（`:119` 真文件不许拒／`:123` 缺席名不许拒／`:240` 仪器红）；`openleaf_overwrite_322:36`（`:70`/`:79`/`:107`/`:124`） | **是**（`1`/`0`/`1`/`5` 处） | 覆盖写承诺退化 ⇒ `openleaf_overwrite:70` "AC#2 RED: PrivateFile refused an existing name"（`yinofb` 正控实测打红它）；探针旗标丢 ⇒ `openleaf_flag:56`/`:57`；守卫逢名就拒 ⇒ `openleaf_guard:119`/`:123`。⚠ `openleaf_flag` ⛔ 经 `privateFile`（头注 `:10-17` 自陈"那两发去牙留它绿"）。 |

## 5 `AC#1..AC#6` 在选形落地前能否判
- `AC#1`（成对读数＋liveness、两平台各一发）：**⛔ 判不了**——选形未定 ⇒ 没有"修后"那一侧；缺的仪器＝驱动 `PrivateFile` 的祖先竞态钉（本腿有雏形，未入库）。
- `AC#2`：料已齐（H2 那 4 枚＋`openleaf_overwrite_322_test.go` 这枚钉今天在场）⇒ **可判形、不可判终态**；仍缺"要不要为祖先守卫补第二枚钉（unlink 形）"的裁定。
- `AC#3`（去牙反控两发）：⛔ 判（本腿⛔ 去牙）；但尺已敏感＝`yinofb` 一发打出 4 枚真红，可作落地腿的反控起点。
- `AC#4`（门禁＋消费者）：⛔ 判（落地的字节今天不存在）；⚠ 本腿自抓一枚越界，见 H6 末。
- `AC#5`（`winsec.go:118-140` 自陈随终态更新）：⛔ 判；★丙形下⛔ 出现"窗口已关"。
- `AC#6`（纪律格）：⛔ 判（要等落地腿的提交面）。

## 6 落地腿需要的人工批准＋本腿自抓
- 甲/乙**都动 `internal/winsec` 产码**＝不是仪器活 ⇒ 落地腿必须是**写腿**，且需要人工批准口令（形状＝「批 342 落<甲|乙>：允许改 `internal/winsec/winsec.go`＋新 landing 件；`winsec_windows.go` 是否同批要单独点名」）。选形权＝编排者＋非实现者，⛔ 本腿。
- 两枚要新形状就⛔ 自行加：① 甲/乙**都不防**的第二枚损伤＝拒绝路径 `os.Remove(resolved)`（`winsec.go:160`）沿祖先链接删树外对象（A3 现量）；建议具名上报给编排者立新票或并入本票判据，⛔ 我立票。② Windows 若选乙，需要本包**第一枚** Windows 侧 handle pinning（今天⛔ 该族）＝新的平台代码面。
- ★本腿自抓（入账归编排者，⛔ 藏）：`29-container-consumers-*` 与 `33-*` 两发里我跑了 `./internal/risk/`⇄`./internal/config/`⇄`./internal/agent/` 的色彩——`internal/risk` 是 `339-v1` 的写面/跑面，与"你的色只取 `./internal/winsec/` 定向与 Linux 容器档"冲突。读数照实入库但**⛔ 当判据**：`TestC26RewriteAccountIsConsumedAtEverySecurityLeg` 在**基线树上 3/3 红**（`/src`⇄`/jia`⇄`/yi` 同号）＝与落点无关、且很可能与别腿同刻跑同包互相污染；`TestNothingWrittenPromiseOnlyHoldsWhenSchemaIsCurrent` 基线 1/3 红⇄甲 3/3⇄乙 2/3＝**摆动枚，本腿不归因**。另：首发 host 档我带 `GOTMPDIR=D:\tmp` 跑 `./internal/winsec/` ⇒ `FAIL=2`（`TestGateJudgesThePrivateSetByResolvedSID`/`TestSealNarrowsAndNamesThePrincipalItRemovedBySID`，逐字 `TempDir RemoveAll cleanup: openfdat … Access is denied`）＝票 311 已具名的**卷形台面既有红**，⛔ 记成祖先那一支新造的红；此后 host 档一律 `unset GOTMPDIR`（`11-*` 起 `110/61/0/0`）。

## 7 必答（逐条）
- **(i)** `oNoFollowPlatform`：`winsec_other.go:29`＝`syscall.O_NOFOLLOW`；`winsec_windows.go:335`＝`0`；★**没有第三枚平台件**（`internal/winsec` 里 `_(windows|other).go` 只有这两枚，`winsec_other.go` 带 `//go:build !windows` ⇒ plan9/js 也走它，而 `syscall.O_NOFOLLOW` 在 plan9/js 未定义＝今天 `GOOS=plan9 go build ./internal/winsec/` **rc=1**，`winsec_other.go:29:35: undefined`）。盘上机制锚：`winsec.go:149` 那一枚 by-name `os.OpenFile` 里，`O_NOFOLLOW` 只约束最后一次成分解析（内核语义＋`winsec_other.go:22-28` 自陈"speaks about the LAST component only"），祖先成分在同一次 `openi()` 里被正常跟随 ⇒ 祖先换手那一跳没有任何在场仪器拦它。
- **(ii)** `sealHandle`/`applyDescriptorPOSIX` 一族：POSIX 侧**今天就能被 `PrivateFile` 复用**（同包非导出；`pinPlacement`→`*os.File`，`f.Fd()` 喂 `Openat`，签名对得上；`sealHandle(*os.File) error` 已在 `winsec.go:158` 被这条路由调用）⇒ **不需要新导出名**。★但"本包为密封做过的 handle pinning"字面上**只在 !windows 存在**（Windows 的 `sealHandle` 走名字、`winsec_windows.go:407-413` 自陈原因＝`WRITE_DAC`），乙在 Windows 是**新写**、不是复用。该族自己的常驻仪器＝`posix_seal_window_120`（`SealFile` 驱动）＋`sealhandle_{identity,privateroute,window}_120`（`sealHandle`/`PrivateFile` 叶子身份）＋`placement_symlink_113`。
- **(iii)** Windows 逐组件下探需要什么：本腿实测两条形都对祖先有效——(a) 逐前缀 `CreateFile(…, OPEN_EXISTING, FILE_FLAG_OPEN_REPARSE_POINT|FILE_FLAG_BACKUP_SEMANTICS)`＋`GetFileInformationByHandle` 读 `FILE_ATTRIBUTE_REPARSE_POINT`（＝今天 `preOpenLeafGuard` 那枚调用 widened；`winsec_windows.go:374-395`），⛔ 判据、叶子写路径仍要按名开 ⇒ check-then-act；(b) `os.Root` 链（`OpenRoot(卷)`→逐组件 `OpenRoot`→叶在 child Root 上开）＝act-time 拒（`path escapes from parent` 1312/1500）。`FILE_FLAG_OPEN_REPARSE_POINT` 只对叶子有效这一句：**本腿复现了**（祖先 junction 自身 `reparse=true` 且 index 换；同一 junction 下的 `anc/blob.bin` 带该旗标 `reparse=false`＝命名目标），⛔ 引腿报。★附带读数：Windows `os.Root` 对树内**相对 FILE 符号链接照旧跟随**（`W-ROOTSYM followed=true`）＝120-r1 那句 POSIX 结论在 Windows 同形成立，Root 不是叶子防线。
- **(iv)** 造不出可断世界的形：**乙在 Windows**（今天没有 pinning 家族可复用 ⇒ 缺 platform code 面：`winsec_windows.go` 新增 `os.Root` 链或 `ByHandleFileInformation` 链），与**甲在非 linux 的 unix**（`syscall.Openat` 不存在；`darwin/freebsd/dragonfly` 只能走 `os.Root` 链，而本包今天 `openbsd/netbsd/solaris/plan9` **已编不过**＝`fstypeOf`）。注入接缝两平台都**不缺**：`SealableTempDirForTest124`＋`os.Symlink`/`mkJunction` 在本机与容器都真造得出（`flips>0` 即证）。新导出名：甲/乙在 POSIX **都不需要**；只有把守卫暴露给消费者才需要（⛔ 本腿不建、不裁）。
- **(v)** 依赖"允许覆盖写"的调用者：**4 枚**＝`store.go:79`、`config/migrate.go:93`、`secret/migrate.go:189`、`agent/spill.go:166`（凭据逐枚引了调用点文档行，⛔ 猜）；`secret/migrate.go:169` 有 `os.Stat` 前置所以不依赖；两枚 `PrivateFileExclusive` 要的正是 `fs.ErrExist`。⇒ 甲/乙若把它改成"已存在即拒"，红会落在 `openleaf_overwrite_322_test.go:70`/`:79`/`:107`（本腿 `yinofb` 正控已把这枚钉打红过）。
- **(vi)** 两枚在册摆动枚的颜色：容器基线一发 `--- PASS` ×2；★**安静复量 3 发**（`logs/24-container-swing-nails-3x.log`）＝**3/3 全 PASS、0 红**；本机档它们带 `!windows` 标签、不参与。两枚红句自陈我**独立复跑了语义**并认其结论：`TestPOSIX120SealWindowLeafSwapIsClosed` 只驱动 `SealFile`（调用点 `:364`/`:479`，本腿尺命中 `PrivateFile`＝**0**）⇒ "a zero window reading here proves nothing" 对**祖先那一支同样成立**；`TestAC1POSIX…322r1` 的腿表（`:46-54`）里唯一算 liveness 的钥匙是**叶子** ELOOP（`:64`）＋祖先换手在它射程外（叶子名自己永远是真文件）⇒ 它那两枚自陈句我认：**今天的仪器没有任何一枚看得见本票要补的那一维**（"祖先那一支"），两枚的绿⛔ 读成"叶子那一支已钉住"的凭据延伸，⛔ 用"票 322 已结案"糊过去；本票新造的可见性就是 H3 里那枚 `CENSUS A2`/`BASE-W`/`PROD-W` 尺。

## 8 节奏与帽（自抓，公共块口径）
- ★本腿在 150 轮帽前收到交件面：31 枚 logs＋11 枚 probe 不可能逐枚一笔入库 ⇒ 按 14:3x 裁定退让为**四笔批量**（census／logs／scripts／probe 各一笔；枚数帽＝参考值⛔判据，"每落一枚就 commit"与笔数相冲时按公共块走＝本句具名入账）。仓外原件（`treeA`／`treeB`／`mut-*`／`gotmp`／`gobuild`＋全部原始 `.log`）⛔ 删、只建、留盘可复跑；入库件凡 `keyrows` 抽取都逐字带 RAW 路径＋**同尺在 RAW 上重算的四数**，⛔ 用抽取冒充原发。
- ★⛔ 选形句自查：全文 **0 处推荐性选形句**（尺＝`grep -nE "推荐|应当选|不妨选"`＝0，正控＝同一把尺在本件能命中"⛔ 选形"字样 5 处 ⇒ 尺不空转；选形权⛔ 本腿）。
- ★票面 H1 第三行（`:87`/`:97`）与 H2 末行（"覆盖写⛔ 任何一枚现有钉钉住"）＝两处**具名更正**，均只追加、⛔ 回改票 322 判据文字、⛔ 翻任何框、⛔ 改 `Status:`。
- ★探针缺陷三枚自抓（读数以修复后的发为准，缺陷发原样入库 `logs/21-…run1.log`）：① A1 第一发把见证名写成 `someone-elses-blob.bin`⇄拼写里的叶子名 `blob.bin` 不符 ⇒ 祖先被跟随后开的是**另一枚新文件**、"ambiguous row"是假读数；② A2 第一发 `seed real` 死于 re-stage 用了已被 rename 消费的目录；③ YI-W 第一发 1500/1500 全红死于 `pathComponents` 在 Windows **不含 volume**、我的下探循环却从 `comps[1:]` 起 ⇒ 少钉一级（＝322-a1 `logs.md 6.4` 那族"探针自身缺陷"，同形复发，具名入账）。
- ★要新形状＝具名上报编排者（H6 那两枚：拒绝路径 unlink 沿祖先链接删树外对象、Windows 侧⛔ 密封 pinning 家族），⛔ 本腿立票。
