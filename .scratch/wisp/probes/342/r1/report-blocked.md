# 342-r1 report-blocked（帽前必交件＝哪半落了／哪半没落／下一程从哪个锚续）

腿＝`342-r1`（票 342＋票 348 落地腿·写腿）。授权来路（⛔ 我自解冻）＝台账 `A1046` 第④节口令逐字：
「批 342＋348 落乙：允许改 `internal/winsec/winsec.go`（那次 open 改走钉住的父 fd ⇄ 两处 `os.Remove`＝`:160`／`:165` 改 fd-relative）＋`winsec_windows.go` 新增本包第一枚 Windows 侧 handle pinning＋新 landing 件；⛔ 新导出名；⛔ 改 C26 入口语义；⛔ 覆盖写承诺若变须先具名点名（现量 4 枚调用者依赖它）。撤销口令「撤 A1046 解冻」」
＋编排者同批增量解冻（`A1049`）＝允许 `winsec_other.go` 的 `pinPlacement`／`pinBase`／`sameObjectPinned`／`chmodPinned` 一族按需最小改动／按签名复用。

起手锚（本腿第一笔现认，⛔ 用派单里的号）＝`fde1df02e06ed9d0c2f2ea6975336f31bd636ef9`。名册尺＝0（自证第一枚；正控 `342-v1` 子串 22 笔、`349-r1` 起头 3 笔＝尺灵）。

## 已落（12 笔 commit，逐枚）

| 笔 | 提交 | 内容 |
|---|---|---|
| 1 | `6c010667` | 新增 `internal/winsec/pinnedleaf_342_unix.go`（form 乙·unix：`leafPin`/`pinnedParent`/`openPinnedLeaf`/`unlinkPinnedLeaf`，`//go:build unix`） |
| 2 | `e00e5b1c` | 新增 `pinnedleaf_342_windows.go`（★过渡件＝旧 by-name 形逐字保留）＋`pinnedleaf_342_fallback.go`（`!windows && !unix` 具名拒绝） |
| 3 | `38a55a0e` | ★两票落点本体＝`winsec.go` 路由：`:149` by-name open ⇒ 钉父目录＋`openPinnedLeaf`；`:160`＋`:165` 两枚 `os.Remove(resolved)` ⇒ `unlinkPinnedLeaf`；★父钉 `defer close` 在路由末端（形变耦合）。diffstat `20/3` |
| 4 | `17d76bac` | 新增 `ancestorworld_342_348_unix_test.go`（共享世界＋racer＋liveness 尺，`//go:build unix`） |
| 5 | `6a8522e3` | 新增 `ancestoropen_342_unix_test.go`（票 342 截断支竞态钉） |
| 6 | `03936f67` | 新增 `ancestorunlink_348_unix_test.go`（票 348 删除支：确定性件＋竞态件） |
| 7 | `da1b3cc4` | ★`pinnedleaf_342_windows.go` 过渡件⇒**本包第一枚真 Windows handle pinning**（`os.Root` 逐组件先钉后验＋叶子身份复验＋`Root.Remove`）。`231/37` |
| 8 | `21a89d0c` | 新增 `ancestorworld_342_348_windows_test.go`（三形世界：junction／symlink-dir／leaf-symlink） |
| 9 | `8ca42ee5` | 新增 `ancestoropen_342_windows_test.go`（票 342 Windows 截断支，两形各 gate） |
| 10 | `683fb19f` | 新增 `ancestorunlink_348_windows_test.go`（票 348 Windows 删除支＋★叶子残差实测） |
| 11 | `43c75685` | Windows 仪器件两枚台面自抓修正（⛔ 产码字节） |
| 12 | `ba1033c2` | 两票尾各**一条**只追加 bullet（`7/0`＋`7/0`＝纯追加） |

零新导出名自证尺＝`git diff <锚>..HEAD -- internal/winsec | grep -E '^\+\s*(func|var|type|const) [A-Z]'` 命中 0（新增顶层标识符全部小写起头或 `//go:build` 后非导出）。

## 哪半落了＝POSIX 半边（帽序①，完整）

台面＝`git archive fde1df02` 三棵仓外树（`treeBase` 基线／`treePre` 基线产码＋我的新钉／`treeLand` 落地），单差异＝只差 `internal/winsec` 字节；容器 `golang:1.27`＋`CGO_ENABLED=0`＋`GOPROXY=off`＋宿主 module 缓存只读；★`TMPDIR`/`GOTMPDIR` ⛔ 绑 Windows 卷。

竞态散布表（每形连打 3 发逐次分子，⛔ 要求逐字相同）：
- 截断支钉 `TestAC1POSIX…NeverTruncatesAForeignFile342r1`：pre `deleted=77/69/85`、`truncated=0/0/1`、`pinnedParentBeforeOpen=0/0/0` ⇄ land `truncated=0`、`deleted=0`、`narrowed=0`、`pinnedParentBeforeOpen=26/20/36`；liveness 同发 `provedLinkCycles=155/160/150`（pre）⇄ `164/160/202`（land）。
- 删除支确定性钉 `TestAC2POSIXRefusalPath…`：pre 3/3 红（逐字 `the foreign witness is GONE after the refusal-path unlink (… someone-elses-tree/leaf-342r1.bin: no such file or directory)`）＝★票 348 `AC#0` 的〔仅腿报〕由本腿自己重打成"收件复跑"；land 3/3 绿（survivor 逐字 `size 75 (was 75), mode -rw-rw-rw- (was -rw-rw-rw-), mtime moved=false`）。
- 删除支竞态钉 `TestAC2POSIXPrivateFileRace…`：pre `deleted=68/67/73`、`pinnedParentBeforeOpen=0` ⇄ land `deleted=0`、`pinnedParentBeforeOpen=26/28/22`、`provedLinkCycles=158/160/195`。
- 全包色（容器★必须 `-v`）：base `RUN=71 PASS_TOP=43 FAIL_TOP=0 SKIP=0` ⇄ land `RUN=74 PASS_TOP=46 FAIL_TOP=0 SKIP=0`（各 2 发逐字同号）；红名册逐名作差＝0 新增红、差集恰为本腿 3 枚新名。
- 门禁（共享工作树，跑前 `unset GOTMPDIR TMPDIR TEMP TMP`）：`gofumpt -l internal/winsec/`=0／`d22scan`=0 clean／长度门=0 GREEN／`go vet` 本机・`GOOS=linux`・`GOOS=darwin` 各=0。

## 哪半没落＝Windows 半边（帽序②，落码入库、★读数未闭环）

`treePre`（基线产码＋新钉）3/3 全红＝尺有牙；`treeLand`（落地）三发里只有 `junction` 子测试稳定绿：
- `…/junction`（land r1）＝★PASS：`provedReparse=13`、本落地那一支答 `formBLegRefusals=6`、`truncated=0 deleted=0`、`returnedNil=13`。
- `…/symlink-dir`（land）＝FAIL；`TestAC2Windows…NeverDeletes…`（land）＝FAIL；`TestAC2WindowsLeafIdentityReCheck…`（land r1）＝FAIL：`formBLegRefusals=0`、名册 `some other failure=300/400`、`returnedNil=63`、`truncated=0 deleted=0`、`provedReparse=87`。
⇒ 这三枚红的形状是**仪器自身的名册桶没归类**（`some other failure`）＋`formRefusals=0`（＝落地的叶子复验那两枚措辞没在名册里命中），⛔ 损伤红。两平台各一发因此**尚未满足**，`winsec.go:118-140` 那段自陈**本腿⛔ 更新**（票 342 `AC#5`／票 348 同条＝两平台都量到才许写"关"）。

## 下一程从哪个锚续（★三步行，⛔ 跳步）

1. **归因那桶 `some other failure`**＝取逐字错误文本。探针已落在仓外（⛔ 入库、⛔ 跟踪）：`D:\tmp\wisp\20261019\342r1\zzprobe\zz342r1probe_windows_test.go.txt` ⇒ 拷进 `treeLand`/`treePre` 的 `internal/winsec/` 跑 `-run TestZZ342r1ProbeWindowsErrorVerbatim`，打 `PROBEERR n=… :: <verbatim>`。判两件事：(a) 落地形是真"拒错地方／漏拒"（⇒改产码，仍在那两段批准射程内）；(b) 只是名册漏了本平台专属措辞（⇒补 `winAncLegs342r1`，⛔ 改产码、⛔ 放宽 gate）。★若归因结果是"必须新导出名才能闭环"⇒ 停手上报、⛔ 自扩批准。
2. **叶子残差那一行**（`TestAC2WindowsLeafIdentityReCheck…`）已是"REPORTED⛔ gated"的形状：`truncated>0` 打 `RESIDUAL OBSERVED, counted and NOT claimed closed`、`=0` 打"一台机器的一发 0 是读数⛔ 是承诺"。第 1 步归因完再判它是否真在跑。
3. **两平台都量到之后**才动 `winsec.go:118-140` 自陈＝票 342 `AC#5`；写的时候 Windows 残差照实写、⛔ 许诺没做到的。然后补去牙反控两发（★成对落点⛔ 相同：一发摘新加的祖先守卫⇒常驻钉必回落红；一发换成"逢名就拒"反形⇒`openleaf_overwrite_322_test.go:70/:79/:107` 那四枚应红），只打仓外树物理字节、每发落刀前验旧锚命中＝期望、落刀后旧锚 0 且新锚≥1、跑前跑后各验 pristine 哈希并报还原哈希、两发串行。

## 台面与纪律备忘（本腿踩过、写下来防下一枚腿重踩）

- ★`go test` ⛔ `-v` 时四数恒 `0/0/0/0`＝尺失灵（本腿第一发全包色就这么取的，⛔ 当"全绿"）。
- ★容器里 `VAR=x unset A B; cmd` 那形只把 `VAR` 赋给 `unset`⇒`go test` 拿到容器默认 `GOMODCACHE`＝`go: downloading golang.org/x/sys`＋`[setup failed]`＋四数全 0，长得像判据响。必须先 `export` 再 `unset`。
- ★`os.WriteFile(p, …, 0o666)` 被 umask(0022) 吃掉⇒`narrowed` 两侧每发都响＝假红；建后 `Chmod`。
- ★祖先链接目标含 `..`⇒Go 的 `doInRoot` 走 escape 规则先答（`path escapes from parent`），本包那枚身份门根本没跑＝测错了支。要"根内跟随"形必须用⛔ 含 `..` 的相对拼写；而相对 reparse 目标按**链接最终所在目录**解析⇒staged 名必须建在那个目录里。
- ★`rename(2)` ⛔ 用符号链接替换非空目录⇒翻转必须 clear-then-rename，那个 gap 里的 ENOENT 计数⛔ 藏、⛔ 当 liveness。
- ★`git add -- A -- B` 的第二枚 `--` 被当 pathspec⇒`rc=128` 整笔中止（无副作用；本腿踩过一次）。
- `x/sys/unix` 在本包 linux 侧是首次直接依赖：容器⛔ 挂下载缓存、`GOFLAGS=-mod=mod` 会红（`module lookup disabled by GOPROXY=off`）＝台损、⛔ 判据红；默认 readonly 可过（`go build` 已实测 rc=0）。
- ⛔ 动既有 `*_test.go` 一字（113／116／119／120／322／336 的断言、⛔ 放宽、⛔ 删、⛔ Skip）；⛔ 翻框（14 格现量＝342 1/6、348 0/7，本腿一枚没翻）；⛔ push；`internal/risk/**` 一字未动。

## 订正（第 14／15 笔之后追加，⛔ 回改上面原文）

上面"下一程三步行"的第 1 步（归因那桶 `some other failure`）＝**本腿已就地做完**（第 14 笔 `3d365ad5`），结论＝⛔ 落地漏拒、⛔ 台面红，而是**本腿自己的名册漏了两条本平台专属措辞**；探针两侧同一条串逐字＝`winsec: refusing to seal …: risk: path traverses a reparse point (junction/symlink) not covered by reparse_point_exceptions`＝已装好的 C26 resolver 的稳态拒绝（`resolveString` → `internal/risk`）。补名册时★⛔ 把它计入 `formBLegRefusals`（⛔ 拿别人的拒答冒充自己那支的 liveness）。

补完重打 3 发 × 两树 ⇒ Windows 终态读数（逐表逐字见同目录 `readings.md` ⇄ `readings-windows.md`）：
- ★`junction` 形（绝对替换名＝act-time 面）＝落地侧 3/3 **PASS**：`truncated=0`、`deleted=0`、本落地那支答 `7/5/3`、`provedReparse=13` 同发在场；基线侧同枚红（`formBLegRefusals=0`）。
- ★`symlink-dir` 形（相对目标、留在根内）＝`deleted` 两侧 0（票 348 那一支落地侧 3/3 不删），但 `truncated` 落地侧 **1/–/1（342 钉）与 2/–/2（348 钉）** ⇒ ★这是一枚**真实残余、本腿⛔ 归因到底**；它⛔ 是基线红（`treeBase` 全绿），所以⛔ 存量、也⛔ 写成"落地失败"＝它是 `AC#5` 要写进注释的那件事。
- ★`leaf-symlink` 形＝`AC#2' leaf-leg split` 逐字 `preOpenLeafGuard=1/1/2, openPinnedLeaf re-check=0, provedReparse=86/81/88` ⇒ **本落地的叶子复验在 400×3 发里一发没答**＝⛔ 被证明在路上（稳态总是 322 那支先答，只有落在"判语与开之间"那枚两 syscall 宽的 gap 才轮到它）。⇒ 派单点名的"⛔ 把新写一条链当成完事"这一味，本腿的量法＝把它写成**待证**而⛔ 写成已关。
- 本机全包成对四数（★`-v`）：`treeBase RUN=110 PASS_TOP=61 FAIL_TOP=0 SKIP=0`（两发逐字同）⇄ `treeLand v1 116/63/2/0`・`v2 116/64/1/0`；红名册逐名＝恰那两枚 Windows 新钉（`…NeverTruncates…342r1` 与 `…LeafIdentityReCheck…342r1`），基线树⛔ 有＝⛔ 存量红。

⇒ ★因此本腿**⛔ 更新 `winsec.go:118-140` 的自陈为"窗口已关"**（票 342 `AC#5`／票 348 同条＝两平台都量到才许写；POSIX 已量到、Windows 那两行⛔），该段至今一字未动＝这是本腿交件里"⛔ 许诺没做到的"的落法。

## 下一程（真正的续点，从锚 `3d365ad5` 之后起）

1. ★归因 Windows `symlink-dir` 那 1～2/400 的 `truncated`＝两候选：(a) 采样与复位在该形上串扰（`firstProof` 被先写入的那行占住，损伤行与拒答行同源）；(b) 真残余＝某条路径在翻转态里绕开钉住的父句柄。做法＝那枚件里把"损伤证据"与"拒答证据"分栏各存一份逐字，再连打 3 发判形。
2. ★把 `leaf-symlink` 的复验做成**确定性**件＝借本包已声明的 seam（`inIdentityWindow120r3` 那一族）在"判语与开之间"插一次翻转，⛔ 靠 racer 撞运气；那枚复验今天⛔ 被证明在路上。
3. ★去牙反控两发（票 342 `AC#3`／票 348 `AC#3`＝成对落点⛔ 相同：一发摘新加的祖先守卫⇒常驻钉必回落红；一发换成"逢名就拒"反形⇒`openleaf_overwrite_322_test.go:70/:79/:107` 那四枚应红），⛔ `-overlay`、每发落刀前验旧锚命中＝期望、落刀后旧锚 0 且新锚≥1、跑前跑后各验 pristine 哈希并报还原哈希、两发串行。
4. 两平台都量到之后才动 `winsec.go:118-140` 自陈；消费者面复认（`internal/risk` 等同刻有 `349-v1` 在读⇒⛔ 同刻跑同一包、撞红先安静复量 3 发再归因）。
5. ★若第 1 步的结论是"必须新导出名才能闭环"⇒ 停手上报、⛔ 自扩批准（票 340 同形那道门）。
