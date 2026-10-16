# 120 v2 终裁表（裁决腿 120-v2，非实现者，只读）

票＝`.scratch/wisp/issues/120-posix-seal-is-check-then-act-154-of-4000-retries-hit-the-window.md`
实现方＝`120-r1`/`120-r2`/`120-r3`；前任裁决腿＝`120-v1`（`docs/evidence/s1/120-v1.md`）。
本件**不产码、不翻 AC 框、不动 Status/票名/-done、不 push**。判语只用：成立／不成立／不可判〔缺哪一发读数〕。

起手锚现读：`git rev-parse HEAD`＝**232f309e8092b15ede2d6f0455d71775e5a5de96**、时刻 **2026-10-16 11:14 +0800**、
`git status --porcelain -- internal/winsec`＝**空**（本腿全程零字节触碰工作树）。
票框尺＝`grep -c '^\- \[ \]'` 与 `grep -c '^\- \[x\]'` 打在票 120 一枚件上：现量 **5 枚未勾／0 枚已勾**（与派单抬头一致，我自己数的）。

## 0. 一句话结论

**`sealHandle` 那一半：`120-v1` 判的"缺世界"在今天这枚面上已经不成立——两形各自有确定性台件，且各自单独去牙都会精确转红；
本格判「成立」（裁活）。** 五格：AC#1 成立、AC#2 成立（附一句可执行答复）、AC#3 成立（含一枚具名"113 那批对 swap-back 失明"的口径落差）、
AC#4 成立、AC#5 成立。另交回三枚交付缺陷（过期行号／过期"票 322 未修"文字／Windows 半格仍空），⛔ 本腿动一字。

台面一律＝`git archive <锚> go.mod go.sum internal cmd tools scripts | tar -x -C /d/tmp/probe120v2/{head,prebase}` 的**仓外导出树**
（尺＝`find <树> -name '*.go' | wc -l`：head **742**、prebase **729**，两枚都非空；全仓尺另有一棵 `full`/`fullpre`，`.go`＝1098/1072）；
容器＝`golang:1.27`（docker **29.6.2**，`-e CGO_ENABLED=0 -e GOPROXY=off` 只读挂宿主 mod 缓存，实测 `getuid=0`）；
突变＝**物理改写仓外树的字节＋sha256 还原**，⛔ `-overlay`。所有数字都是我自己现跑的：⛔ 抄过 `120-v1`/`120-r2`/`120-r3` 任一枚读数。

HEAD 那棵树的 `internal/winsec/winsec_other.go` 现量 sha256＝**19f3c211055011cc699a8d8027ef49ff543ede0da5c863f58aaaad21a640d3e2**（456 行，量的）。
⚠ 这一枚与 `120-r3` 全件引用的 clean 基准 `137442cd…` **不同**：票 322 的两枚写腿（`8baec646`、`1fbc9067`）在 `120-r3` 交件之后改过这枚文件（＋42/−1 与拆件）。
⇒ 本腿所有读数都落在**今天交付的字节**上，比腿报那一版更强；下面凡引用腿的数处，都给了我的同尺现量。

## 1. 先复现前任（派单第 0 步第 3 条）

| 我打算当依据的前任说法 | 我的同尺现量 | 复现？ |
|---|---|---|
| `120-v1` §1：改前码面 hit-foreign>0（它 43/6000；`120-r2` 26/6000） | 改前树＝`git archive 58de465a`（其 `winsec_other.go` 现量 sha `34cc6bd8ca0d885f658e62373b062b48d2b737095ff29b19f7cd98be47f00ec0`＝腿报那枚"改前"串）＋本轮仪器，`WISP_120_BASELINE=1`：三形 **9 / 8 / 17**，合 **34/6000>0** | 方向复现（按 `A991` 口径三枚数**不可相减**，只作存在性独立复认） |
| `120-v1` §1 旁证：改前 `refused-at-act=0`（闸门不存在） | 同发逐形 `of-which-refused-at-act=0 / 0 / 0`；改后 HEAD 树同尺＝19/9/20 非零 | 复现 |
| `120-v1` §5 名册顶层 20 枚、`120-r3` 报 27 枚 | 我今天同一把尺现量 **28 枚**（见 §5） | 复现其增长方向；绝对值以我这枚为准 |
| `120-v1` §2：摘掉 `sealHandle` 两支的牙 ⇒ 整包 `55/55/0/0` 零变色 | **不复现，也⛔ 复现不了**：那正是"缺世界"的成因读数，世界补上之后同一条牙必然叫红（§4 两发各红 2 枚）。这一枚是本腿要裁的题，不是依据 | 具名：不适用 |

## 2. AC#1 重试形状做成可重跑用例 —— 判 **成立**

① **用例与多样本机制**：`internal/winsec/posix_seal_window_120_other_test.go` 的三枚具名
`TestPOSIX120SealWindowLeafSwapIsClosed`（`:321`）／`…AncestorAbsoluteLinkIsClosed`（`:378`）／`…AncestorRelativeLinkIsClosed`（`:404`）。
多样本＝常量 `sealWindowSamples120 = 4`（`:100`）× `sealWindowAttempts120 = 500`（`:101`）⇒ **每形分母 2000、每形 4 枚样本各判各的**（尺＝读那两枚常量＋`judgeRun120` 的样本循环 `:196-212`）。
只跑一发交不了差：`-run TestPOSIX120` 打印 3×4 枚样本行＋3 枚 total 行。

② **成对读数（同尺、两台仓外树，缺一不成凭据）**：

| 形 | 改前（`WISP_120_BASELINE=1`）hit-foreign/2000 | 改前活性 refused-at-act | HEAD 默认严格档 hit-foreign/2000（4 枚样本逐枚） | HEAD 活性（total） |
|---|---|---|---|---|
| leaf-swap | **9**（of which reported-nil 9） | 0 | 0/0/0/0 | check=942 act=19 sealed-own=943 |
| ancestor-swap-absolute-link | **8**（nil=7） | 0 | 0/0/0/0 | check=82 act=9 sealed-own=121 |
| ancestor-swap-relative-sibling-link | **17**（nil=17） | 0 | 0/0/0/0 | check=133 act=20 sealed-own=137 |
| 合 | **34/6000 > 0** | 皆 0 | **0/6000** | 皆非零（"命中 0"与"仪器死了"没互抵） |

**"抓到窗口的次数"是量出来的**：打在 `posix_seal_window_120_other_test.go:171-172` —— 每次 attempt 之后 `witnessMode()` 读回外站证人的**盘上 mode**，
等于 `0o600` 才 `hitForeign++`（⛔ 从错误文案推、⛔ 从返回值推）；单位＝**次数**（分母＝attempt 数，同一行 `:138` 打印 `attempts=%d hit-foreign=%d`）。

③ **反形之问（会不会红？还是只会报 0？）**：三发独立证据，都会红 ⇒ "只会报 0"这一支排除。
(a) **未修码面用默认严格档跑**＝`v2-prebase-strict.log`：三枚具名全红 `RUN=3 TOPPASS=0 TOPFAIL=3 TOPSKIP=0 rc=1`，
红句逐字 `AC#1 RED sample 1/4 shape=leaf-swap: the seal landed on the foreign tree 1 time(s), 1 of them while returning nil (window open on THIS sample; every sample is judgement, no averaging)`
—— **一发样本里只命中 1 次也叫红**，敏感性有读数。
(b) **在已修码面上把窗口放大**＝物理摘 `sameObjectPinned` 的同一对象那一支（`winsec_other.go:222` 一支，⛔ 两支一起摘）⇒
`TestPOSIX120SealWindowAncestorRelativeLinkIsClosed` 转红（`v2-tooth-walk-samefile-face.log`，FAIL=2 枚）。
(c) **种一次真竞态形状（确定性、非概率）**＝`120-r3` 的两枚台件在摘对应牙时精确转红（§4）。

④ ⛔ 用平均抹尾部、⛔ 只报最好那次：判据在**样本循环内**（`:206` 逐样本 `t.Errorf`），total 行（`:219`）只是额外打印；
两份日志里 12 枚样本行逐枚全列。实现件没有"取最好"的写法 ⇒ 成立。

## 3. AC#2 前置问题（openat / O_NOFOLLOW / 按已验证 fd 密封，在本仓现有 API 形状下做得到吗）—— 判 **成立**

**可执行答复一句**：**"按已验证的描述符密封"这一形做得到，而且就是今天落地的形状；"逐级 `openat`+`O_NOFOLLOW` 的完全体"做不到 ⇒
本票该做机制（也做了），但只关到"密封那一跳"，⛔ 任何人把它读成"这一族从根上关掉了"。**

凭据（每一枚都是本腿现跑，标尺与射程）：

1. `internal/winsec/winsec_other.go:29` `const oNoFollowPlatform = syscall.O_NOFOLLOW`；**今天的数值＝`0x20000`（131072）**
   （尺＝仓外独立模块 `/d/tmp/probe120v2/probe/main.go` 在容器里 `go run .` 打印，`PROBE_rc=0`，件 `logs/v2-ac2-probe-readout.txt`）。
2. **`O_NOFOLLOW` 逐 GOOS 编译矩阵**（尺＝在导出树里放一枚只引用该常量的件 `probenof/`，逐 `GOOS/GOARCH` 交叉 `go build`，只看 rc；11 平台×2 枚＝22 发）：
   COMPILES＝linux/amd64、darwin/arm64、freebsd/amd64、openbsd/amd64、netbsd/amd64、dragonfly/amd64、solaris/amd64、aix/ppc64；
   rc≠0＝**windows/amd64、js/wasm、plan9/amd64**（三枚都 rc=1；我只抓了 rc 与包名首行，⛔ 抓全错误串 ⇒ 说"取不到该常量"有据、说逐字⛔ 有据）。
3. **`syscall.Openat` 同把尺**：只有 **linux/amd64 与 aix/ppc64** COMPILES，其余 9 平台 rc≠0；darwin 那一发的逐字第二行＝
   `probeat\probeat.go:5:19: undefined: syscall.Openat`。⇒ ★**这一枚更正前任口径**：`120-r1`/`120-v1` 写的"只有 linux"少了 aix（⛔ 改任何已定案文字，只具名记现量）。
4. **`os.Root` 令 `O_NOFOLLOW` 失效**（现跑，⛔ 引用注释）：同一枚探针里 `Root.OpenFile("leaf.bin", O_RDONLY|O_NOFOLLOW)` 对**相对兄弟链接**
   读数＝`OPENED, and the descriptor stands on the TARGET object=true`（`os.SameFile(fd, target)`＝真）⇒ 叶子链接照样被跟，这条⛔ 当防线。
   同发 `Root.Create` 两读：对空名 `OPENED`（不宣称独占）、对带链接的名也 `OPENED`。
5. **Windows 那枚常量今天的值＝`0`**：`internal/winsec/winsec_windows.go:335` `const oNoFollowPlatform = 0`；
   OS 手段＝`winsec_windows.go:374 func preOpenLeafGuard`（`:379` 以 `FILE_FLAG_OPEN_REPARSE_POINT|FILE_FLAG_BACKUP_SEMANTICS` 开句柄问"是不是链接"），
   那是**判据不是不跟随的覆盖开** ⇒ windows 侧没有对应物可关这一族。
6. **做得到且已落地的那一半**（读码定位，行号＝我今天对 `19f3c211…` 那枚文件现量）：`pinPlacement :149`（钉→复验→动，叶子 `:189` 用 `O_RDONLY|syscall.O_NONBLOCK`）、
   `chmodPinned :255`（`f.Chmod`＋同一枚 fd 上 `f.Stat` 读回）、`sealHandle :287`（持 fd 入口自带的两次 inline 复验）。
7. **还开着的那一腿**（产码自己写着，⛔ 本腿新造）：`winsec_other.go:103-113` ⇒ 叶子腿已由票 322 用 `O_NOFOLLOW` 关上（接线尺见 §7），
   **祖先腿仍开**（`O_NOFOLLOW` 只说最后一枚成分）。

## 4. ★`sealHandle` 那一半（`120-v1` 缺口①、`120-r3` 声称补齐）—— 判 **成立**（"缺世界"作废；裁死→裁活）

台面＝HEAD 导出树（`19f3c211…`），脚本 `.scratch/wisp/probes/120/v2/tooth.sh` 与 `tooth2.sh`，逐发驱动日志
`logs/v2-tooth-driver.log`、`logs/v2-tooth2-driver.log`。整包基线色（未去牙）＝`RUN=71 TOPPASS=43 TOPFAIL=0 TOPSKIP=0 rc=0`
（尺＝`grep -c '^=== RUN'` 含子测／`grep -c '^--- PASS'` 只数顶层，两把⛔混用）。

① **世界真被造出来**（证据与读数同一发，`logs/r1c-own7.log`，`-run '120r3$'`：RUN=7 TOPPASS=7 TOPFAIL=0 TOPSKIP=0 rc=0）：
真对象枚 `sealhandle_identity_120_other_test.go:44: AC#1' world: fd and name = two different objects; name kind=file mode=-rw-rw-rw-`
＋同发 `:48 … leg that spoke = sealHandle inline SameFile face`；链接枚 `:86 … leg that spoke = sealHandle inline link face`
＋`:91 read-back: name kind=link, caller's object=-rw-------, witness=-rw-rw-rw-`。夹具的活性是**硬**的：`identity:41-43` 在调用前 `t.Fatalf` 问"名与 fd 是不是两枚对象"，
世界没造出来这一发就红 ⇒ "绿＝没世界"这一形被仪器自己堵死。

② **逐支单独去牙**（每支一发；`false && ` 前插、只动一行）：

| 摘的牙（现量行号） | 突变 sha256 | rc / 四数 | 变红的用例（逐名） |
|---|---|---|---|
| `winsec_other.go:300`（sealHandle 链接那一支） | `989af68717ee5b05296c8fc094be253dd2e0588712afdb00fef4a08ab1d3f88a` | **1**／RUN=71 TOPPASS=41 **TOPFAIL=2** TOPSKIP=0 | `TestAC1POSIXSealHandleNameBecomesLinkAfterPinnedSeal120r3`、`TestAC1POSIXPrivateFileRefusesWhenNameBecomesLinkAfterPinnedSeal120r3` |
| `winsec_other.go:303`（sealHandle 真对象那一支） | `76e49806edae651d03f2fcaf3f58ef009f780e820309031494bfecb56124dd54` | **1**／RUN=71 TOPPASS=41 **TOPFAIL=2** TOPSKIP=0 | `TestAC1POSIXSealHandleNameCarriesAnotherRealObject120r3`、`TestAC1POSIXPrivateFileRefusesWhenNameCarriesAnotherRealObject120r3` |

去牙红句（逐字一行，两支各一枚）：
`:300` 摘后＝`sealhandle_identity_120_other_test.go:86: AC#1 RED name-becomes-a-sibling-link-after-the-pinned-seal: refusal came from sealHandle inline SameFile face (winsec_other.go:272), want the leg carrying "the name carries a link, so this descriptor is not the object the caller named"; verbatim: …`
——⇒ **仍然被拒、但拒因换腿**，所以判据必须问"哪一腿说的"，⛔ 只问"有没有拒"（这条判据真的在分辨：它红了）。
`:303` 摘后＝`sealhandle_identity_120_other_test.go:48: AC#1 RED name-carries-another-real-object: the call returned nil, so the leg carrying "this descriptor is not the object the name carries now, so the bytes would not land in what was just sealed" never spoke - this is the shape of that leg with its teeth taken out`
＋盘上事实那一行 `:60: AC#1 RED: the caller's own object was narrowed to -rw------- while the refusal says the bytes would not land in what was just sealed`（没牙时 `chmodPinned` 真落在 fd 的对象上）。

**盲跑那枚（枚 7）今天没跟着红** ⇒ 与 `120-r3` §5① 自写的"红 3 枚是概率性额外证据、⛔ 写成保证"**一致**，本腿⛔ 据此加判据、也⛔ 据此减判据。
自然 landing 我读了两发（⛔ 升硬，按 120-v1 §2 裁③）：整包那发 `sealHandle inline SameFile=10 / link=9`、`-run '120r3$'` 那发 `9 / 1`
（同发另有 `returned nil: the flip missed every check = 26 / 50`、`other failure = 41 / 21`，枚 7 的两发**证人 mode 与内容判据全绿**：`witness holds its own bytes after 2000 raced attempts`）。

③ **同形两对之问（⛔ 两枚腿互冒充）**：这把尺我自己现跑（`grep -n` 整文件）——
`ModeSymlink != 0` 命中 **3 行**＝`:219`（`sameObjectPinned` 内）、`:300`（`sealHandle` 内）、`:455`（`ancestorIsLink` 的 return，⛔ 是判据）；
`!os.SameFile(pinned, now)` 命中 **2 行**＝`:222`、`:303`。⇒ 同一文件里确实有**两对同形判据**，按行号下牙会打错。
本腿的改法＝**由 `func sealHandle` 的体范围（现量 287..307）里找行**，落牙后硬校 `false &&` 全文件**恰好 1 处**并打印 `LANDS` 前后行逐字（驱动日志里每发都在）。
★再补一发正面区分：把**另一对**（`sameObjectPinned` 两支一起摘，突变 sha `3dc077484d9b4fd1e6e6cd87f73e49d7b9282958c0caa647cf6a5e095e9dddbf`）
⇒ 红的是**完全不同的一批**：`TestAC1POSIXOpenThroughASwappedLeafLinkNeverTruncatesAForeignFile322r1`、`TestAC1POSIXSealHandleNotReachedByAPrePlantedLink120r3`、
`TestAC1POSIXPrivateFileWindowIsRacedWithoutTheSeam120r3`、`TestPOSIX120SealWindowAncestorRelativeLinkIsClosed`（TOPFAIL=4），
而我上面两发的 4 枚链接/真对象形**全绿** ⇒ 两对判据可分辨、本腿下的牙在对该的那对上。拆开单摘（`tooth2.sh`）：只摘 `:219` ⇒ 只红成因枚 3（1 枚）；只摘 `:222` ⇒ 红 322 枚＋仪器形 3（2 枚）。

④ **还原**：三发全部 `RESTORE MATCHES`，还原 sha256＝**19f3c211055011cc699a8d8027ef49ff543ede0da5c863f58aaaad21a640d3e2**（＝起手 HEAD 现量，逐字符同串），
还原后整包复跑回到基线色：`v2-restored-reconfirm.log`＝`RUN=71 TOPPASS=43 TOPFAIL=0 TOPSKIP=0 rc=0`。
底版留 `/d/tmp/probe120v2/head/internal/winsec/winsec_other.go.bak-120v2-teeth-on`（只建不删）。

⇒ **判语**：`sealHandle` 的两枚 inline 复验今天**各有能走到它的世界、各有单独的牙**；`120-v1` 的"缺世界"是对当时那枚字节（`137442cd…`）的正确读数，
对今天的交付字节**不再成立**。本格判**成立**。

## 5. AC#3 反半边（票 113 那批容器用例：一枚不变红、一枚不被 Skip）—— 判 **成立**（含一枚具名口径落差）

* 枚数现量：尺＝容器 `-run 'TestAC[1-4]POSIX' -count=1 -v` 于 HEAD 导出树，`grep '^--- PASS' | awk '{print $3}' | sort -u | wc -l`（射程＝`internal/winsec`，量的）
  ＝ **顶层 28 枚**；同发四数 `RUN=39 / TOPPASS=28 / TOPFAIL=0 / TOPSKIP=0 rc=0`（39−28＝11 枚是子测试的 `=== RUN` 行）。
  ⚠ 票面写"16"、`120-v1` 现量 20、`120-r3` 报 27、本腿 28 ⇒ 只具名记现量与尺口径，⛔ 改票面原文（措写要人批）。
  第 28 枚是票 322 落的 `TestAC1POSIXOpenThroughASwappedLeafLinkNeverTruncatesAForeignFile322r1`。逐枚名册＝`logs/r1b-roster.log`。
* 覆盖"链接在祖先各位"的（点名）：`TestAC1POSIXSealFileRefusesALinkAncestorAtEveryDepth`、`TestAC2POSIXAncestorGuardRefusesASpellingThroughASymlink`、
  `TestAC2POSIXABackslashInALinkNameIsStillALinkAncestor`、`TestAC3POSIXSealStillWorksNextToAndThroughRealDirectoriesAndLinks`、
  `TestAC1POSIXPrivateFileThroughASymlinkRefusesAndWritesNothing`、`TestAC1POSIXSealDirThroughASymlinkRefuses`、
  `TestAC1POSIXSealFileThroughASymlinkRefusesAndLeavesTheForeignTreeAlone`、`TestAC1POSIXSealFileThroughABackslashNamedLink`、
  `TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119`、119 那三枚 route 枚、125 那三枚 seam 枚。
* 覆盖**反半边"合法路径照旧成功"**的（点名）：`TestAC3POSIXSealFileStillNarrowsAPlainFileInsideTheNamedTree`、
  `TestAC4POSIXFloorAnswersInsideTheNamedTree`、`TestAC2POSIXSeamAcceptsTheHonestPOSIXAnswer125`、
  `TestAC1POSIXSymlinkedConfigDirRouteBecomesSealable119`、`TestAC1POSIXSymlinkedTempDirRouteBecomesSealable119`、
  `TestAC2POSIXDoesNotFoldABackslashIntoASeparator`、`TestAC3POSIXSealDoesNotFoldABackslashIntoASeparator`，
  以及 120 家新落的 `TestAC1POSIXSealHandlePlainObjectStillSeals120r3`（诚实路径必须仍 seal）。
* **"同一个对象被换掉"这一形真被造出来了吗？** 造出来了，但**不在票 113 那批里**——凭牙（§4 ③ 的拆分）：只摘 `sameObjectPinned` 的同一对象支（`:222`）⇒
  变红的是 `TestPOSIX120SealWindowAncestorRelativeLinkIsClosed`（120 自己的仪器形 3）与票 322 那枚；**28 枚 113 名册一枚都不红**。
  ⇒ 具名口径落差（⛔ 判成 120 的缺陷，票面 AC#3 只要求那批不变红）：**113 那批钉的是"链接站着时的拒绝"，对"换过去又换回另一枚真对象"这一形是失明的**；
  这一形今天的凭据只落在 120 自己的仪器＋`120-r3` 的台件上。下一枚读票面的人不许把 113 那批读成"祖先被换掉这一族都钉住了"。
* 容器是 root（实测 `getuid=0`）⇒ 权限类断言在本台面失真（票 113 的常注）；这一族夹具头自己写明不碰权限面（`sealhandle_window_120_other_test.go:52-55`）。

## 6. AC#4 若 AC#2 的答案要改 `winsec.go` 那段边界话 ⇒ 先登记交回编排者 —— 判 **成立**

尺＝`git diff 58de465a..232f309e -- internal/winsec/winsec.go`：现量 **2 个 hunk、+42/−1**（`@@ -83,6 +83,15 @@` 与 `@@ -105,13 +114,45 @@`；被删的那一行是 `f, err := os.OpenFile(resolved, flags, perm)`）。
`git log 58de465a..HEAD -- internal/winsec/winsec.go` 现量**只有 1 枚 commit**＝`8baec646`（收票 322-r1 的盘上遗产，编排者代提）。
⇒ **票 120 的三枚腿一枚都没碰 `winsec.go`**（`120-r2` 的写面＝`winsec_other.go`＋一枚仪器；`120-r3` 全程零产码，腿自报且与我现量 sha 相符）。
唯一改动那段话的来源是**票 322**，而它**已登记**：台账 `docs/reports/pending-and-issues.md` 的 **A1001**（第 7 条）逐字接线尺＝
`winsec.go:146 preOpenLeafGuard(resolved)`＋`:149 flags|oNoFollowPlatform`、`winsec_other.go:29`、`winsec_windows.go:335/:374`——
这四枚落点本腿逐枚 `grep -n` 复认存在（§3 第 1、5 条即同两枚）。⇒ 判成立：无须本腿新登记，也没有"改了却没登记"。

## 7. 能力类必答：这半格在生产调用者上通不通 —— 通（调用点数出来，⛔ 拿 `go build` rc=0 冒充）

尺＝`git grep -n sealHandle -- 'internal/**'`，按文件分栏（量的）：
* **产码**（非 `_test.go`）3 枚件、**调用点 1 枚**＝`internal/winsec/winsec.go:158 if err := sealHandle(f); err != nil {`，位于 `privateFile` 体内；
  入口两条都经它：`PrivateFile`（`winsec.go:95`→`privateFile(…, O_WRONLY|O_CREATE|O_TRUNC, …)`）与 `PrivateFileExclusive`（`winsec.go:105`→`…, O_EXCL, 0o600`）。
  定义 2 枚＝POSIX `winsec_other.go:287`、Windows `winsec_windows.go:413`。
* **测试**5 枚件、29 处命中＝`sealhandle_identity_120_other_test.go`(12)、`sealhandle_window_120_other_test.go`(8)、
  `sealhandle_privateroute_120_other_test.go`(6)、`openleaf_322_other_test.go`(2)、`openleaf_guard_322_windows_test.go`(1)。
⇒ 不是"只有测试在调"：真入口一枚、两条公开 API 都走它，且枚 5/6 走的就是 `PrivateFile` 这条**产码路由**（§4 读数）。
⚠ 同一把尺露出的半格：Windows 那枚 `sealHandle` ＝ `applyDescriptor(f.Name(), false)`，**没有 inline 同一对象闸门**（见 §9）。

## 8. AC#5 门禁 —— 判 **成立**（逐把带 rc=，全为本腿现跑）

| 尺 | 口径／射程 | rc | 读数 |
|---|---|---|---|
| 容器整包 | `go test -count=2 -v ./internal/winsec/`，HEAD 导出树 | **0** | `RUN=142 / TOPPASS=86 / TOPFAIL=0 / TOPSKIP=0`；口径＝不同顶层名 **43** 枚 ×2＝86 条 `--- PASS`，`=== RUN` 含子测（43 顶层＋28 子测＝71）×2＝142 ⇒ 与票面"=== RUN == 不同测试名 ×2"一致 |
| 容器名册（AC#3） | `-run 'TestAC[1-4]POSIX'` -count=1 -v | **0** | 39/28/0/0 |
| 120r3 本族 | `-run '120r3$'` | **0** | 7/7/0/0 |
| 仪器默认严格档 | `-run TestPOSIX120` | **0** | 3/3/0/0，六样本 hit-foreign 全 0 |
| `gofmt -l` | 五枚（`winsec_other.go`＋四枚 120 件） | **0** | 列出 **0 枚**（退码恒 0 ⇒ 读列件） |
| `gofumpt -l` | `/d/work/base/gopath/bin/gofumpt.exe` 同五枚 | **0** | 列出 0 行／**0 字节** |
| `go vet ./internal/winsec/` | 宿主 windows/amd64（⛔ 不编 `winsec_other.go`） | **0** | — |
| `GOOS=linux GOARCH=amd64 go vet` | 真读 `!windows` 那一半（含四枚测试件） | **0** | — |
| `GOOS=darwin GOARCH=arm64 go vet` | 交叉复认 | **0** | — |
| `sh scripts/d22scan.sh` | **纯净快照**＝`git archive 232f309e` 全仓导出树 | **0** | `clean - no D22 ban violations` |
| 同把尺成对分母 | 58de465a 全仓导出树 ⇄ HEAD | 0／0 | `ban#8 internal/=563→573`、`cmd/=136→139`；`#1-5 internal/=239`、`cmd/=41`、`#6/#8 frontend/=116`、`#7 tools/=25`、`examined 280` 持平 ⇒ **无一 scope 下降** |

⚠ 口径注两枚：(1) `design/` 我这把尺读 **30** 枚、前任报 39——⛔ 是 scope 掉了，是他们跑在**带未跟踪件的工作树**上、我跑在 `git archive` 导出树上
（`git ls-tree -r design/` 对 `58de465a` 与 `232f309e` 两枚锚都＝30 枚，现量）。(2) 导出树里 `git rev-parse` 取不到 ⇒ d22scan 自己印了
`gitignore rules NOT APPLIED`（A207 的已知机器相关口径），方向是**扫得更宽**、⛔ 放行。
`RUN=0`／`no tests to run`／`0xc0000135`／`build failed`／`walk` 这一族在我所有发里**都没出现**（正控＝同镜像同台的 `VET_rc=0` 与 8 发 rc 全读到位）。

## 9. 台面故障具名（⛔ 计入判据）＋交付缺陷（⛔ 本腿动一字）＋没裁完的格

**本腿自己撞坏的两把尺（实录，随批入库）**：
1. 我把读者档与严格档两发塞进**同一枚** `docker run -e WISP_120_BASELINE=1` ⇒ 第二发也带着 env，`STRICT_rc=0` 是**假读数**；
   作废件留底 `logs/v2-prebase-strict-DUPLICATE-of-reader.log`（只建不删），分开两枚容器重跑后 `rc=1` 才是 §2③(a) 那一发。
2. 我第一次用 `grep 'live scope work' | head -1` 抓 d22scan 分母 ⇒ 抓到仪器自测里那枚临时小树的摘要（`internal/=14`），坏尺；改 `tail -1` 取最终行为 §8 的数。

**交付缺陷三枚（交回编排者处置；都⛔ 改判据，本腿一枚没翻）**：
1. **过期行号**：`spokeLeg120r3` 的名册（`sealhandle_window_120_other_test.go:88-92`）把四枚腿印成 `winsec_other.go:188/191/269/272`，
   今天在 **219/222/300/303**（票 322 把那枚文件从 425 行抬到 456 行）。判据本身按**文案子串**匹配，四串在产码逐字仍在 ⇒ 红绿与归属不受影响，
   但打印出来的行号会误导下一枚腿。应改句子＝把行号列与函数体范围现算，或只印函数名。
2. **过期"未修"文字**：`sealhandle_privateroute_120_other_test.go:233-238`（含 `winsec_other.go:78-82` 的指路）与
   `sealhandle_window_120_other_test.go:7/:21-23` 仍写"票 322 unfixed at this anchor"；今天 322 的**叶子腿已落地**（`winsec.go:146/:149`＋`oNoFollowPlatform`），
   仍开的只有**祖先腿**（`winsec_other.go:103-113` 逐字）。本腿两发盲跑里那一形**都没出现**（`witness holds its own bytes after 2000 raced attempts`）＝与已落地一致。
3. **票 120 的 Windows 半格仍空**：`winsec_windows.go:413` 的 `sealHandle` 只有 `applyDescriptor(f.Name(), false)`，缺同一对象闸门；
   票面 `Packages:` 段逐字**禁改** `winsec_windows.go`（票 115 在飞）⇒ ⛔ 本腿动、也⛔ 判。

**没裁完的格（具名＋缺哪一发读数）**：
* **Windows 侧**：本格**不可判**〔缺的读数＝windows/amd64 上"名字被换成另一枚真对象"这一形的台件与去牙反控；今天仓里**没有**这一枚用例，
  且本腿禁碰那枚件、容器只读 linux〕。出处＝`120-r3` §5④ 那句"`winsec_windows.go:328` 缺同一对象闸门"（⛔ 票面文字），
  该指路行号今天也已过期：现量 `:328` 落在票 322 的 `preOpenLeafGuard` 注释段里，`sealHandle` 定义在 **`winsec_windows.go:413`**（＝缺陷 1 同一族的第二例，⛔ 本腿动字）。
  票面 `Packages:` 段逐字**禁改** `winsec_windows.go`（票 115 在飞）。
* **射程外消费者变色**（`internal/risk`、`cmd/wisp`、`internal/secret` 在 Linux CI 上会不会因 r 位收紧变绿变红）：**不可判**
  〔缺的读数＝`go test ./cmd/wisp` 整包（票面 314 死亡区、本腿⛔ 跑）与 `internal/tools` 整包（同批只读腿 `270-v1` 在用，本腿⛔ 跑）〕；
  本腿只交回 `d22scan` 与 `go vet` 两把静态尺同色。
* **自然 landing 的分布**：只两发（10/9 与 9/1），按 120-v1 §2 裁③**只上报⛔ 升硬**；要多发归编排者（⛔ 阈值）。

## 10. 复跑（可逐字重放）

```sh
cd 'D:/work/workspace/projects plans/Wisp'
V=/d/tmp/probe120v2
mkdir -p "$V/head" "$V/prebase"
git archive 232f309e go.mod go.sum internal cmd tools scripts | tar -x -C "$V/head"
git archive 58de465a go.mod go.sum internal cmd tools scripts | tar -x -C "$V/prebase"
cp internal/winsec/posix_seal_window_120_other_test.go "$V/prebase/internal/winsec/"
sh .scratch/wisp/probes/120/v2/tooth.sh    # 基线＋两枚 sealHandle 牙＋另一对两支一起摘＋还原复跑
sh .scratch/wisp/probes/120/v2/tooth2.sh   # 只摘 sameObjectPinned 的一支（两发）
MSYS_NO_PATHCONV=1 docker run --rm -e CGO_ENABLED=0 -e GOPROXY=off \
  -v /d/work/base/gopath/pkg/mod:/gomod -e GOMODCACHE=/gomod -e WISP_120_BASELINE=1 \
  -v "$V/prebase":/work -v probe120v2build:/root/.cache/go-build -w /work golang:1.27 \
  sh -c 'go test -count=1 -v -run TestPOSIX120 ./internal/winsec/'      # 改前读者档
# ⛔ 把上面那枚 -e 复用给"严格档"那一发（那是我撞过的台面故障，见 §9）
sha256sum internal/winsec/winsec_other.go   # 19f3c211055011cc699a8d8027ef49ff543ede0da5c863f58aaaad21a640d3e2
```
读数件（⛔ 删）：`.scratch/wisp/probes/120/v2/logs/`（23 枚）＋两枚驱动脚本；仓外：`/d/tmp/probe120v2/{head,prebase,probe,full,fullpre,v2-*.log}`、
去牙底 `head/internal/winsec/winsec_other.go.bak-120v2-teeth-on`。仓内工作树全程未碰：`git status --porcelain -- internal/winsec` 空。
