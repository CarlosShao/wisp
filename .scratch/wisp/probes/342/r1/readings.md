# 342-r1 读数汇总（POSIX 半边完整＋Windows 半边成对）

⛔ 本文件只抄日志里已有的读数，⛔ 复述、⛔ 圆场。逐字全文＝仓外 `D:\tmp\wisp\20261019\342r1\logs\`（POSIX：`10-pair-*`／`30-fullpkg-*`；Windows：`60-winpair-*`／`70-probe-*`／`80-host-full-*`；门禁：`20-*`／`21-*`／`22-*`～`24-*`）。Windows 那三张表与残余的展开见同目录 `readings-windows.md`。

## 台面（三棵仓外导出树，⛔ 复用别人的树、⛔ `-overlay`、⛔ 改跟踪文件）

- 锚＝`fde1df02e06ed9d0c2f2ea6975336f31bd636ef9`＝本腿起手的 HEAD（⛔ 用派单里的号）；`git archive <锚> go.mod go.sum internal cmd tools scripts frontend models docs` 八件套＋仓根 `docs`＋手工三件目录（`build/`、`third_party/sherpa-onnx/`、`frontend/dist`）。
- `treeBase`＝锚本身（基线色）／`treePre`＝锚＋只拷本腿新增的 6 枚 `*_test.go`（＝常驻钉打在旧产码上）／`treeLand`＝锚＋本腿 `internal/winsec` 全量字节。
- ★单差异性质现量＝`winsec.go` 哈希 base⇄pre 逐字同（`41ba9a9b8fcf7050`）、land 侧 `9086f2bbbf531823`；`GO_COUNT` 766／772／775。
- POSIX 档＝`golang:1.27`＋`CGO_ENABLED=0`＋`GOPROXY=off`＋宿主 module 缓存只读；★`TMPDIR`/`GOTMPDIR` ⛔ 绑 Windows 卷（＝容器自 `/tmp`，overlayfs 存 mode bits）。Windows 档＝本机 `go1.27.1 windows/amd64`。

## 成对四数（★必须带 `-v`；⛔ 四数混用、⛔ 相减）

POSIX（容器，`go test -v -count=1 -timeout 30m ./internal/winsec/`）：
- `treeBase`＝`RUN=71 PASS_TOP=43 FAIL_TOP=0 SKIP=0`（两发逐字同号）
- `treeLand`＝`RUN=74 PASS_TOP=46 FAIL_TOP=0 SKIP=0`（两发逐字同号）
- 红名册逐名作差＝**0 枚新增红**；差集恰＝本腿 3 枚 POSIX 新钉，全 PASS。

Windows（本机）：
- `treeBase`＝`RUN=110 PASS_TOP=61 FAIL_TOP=0 SKIP=0`（两发逐字同号＝基线全绿）
- `treeLand`＝v1 `RUN=116 PASS_TOP=63 FAIL_TOP=2 SKIP=0`／v2 `RUN=116 PASS_TOP=64 FAIL_TOP=1 SKIP=0`
- 红名册逐名（v1）＝`TestAC1WindowsPrivateFileThroughASwappedAncestorNeverTruncatesAForeignFile342r1`（＝`symlink-dir` 那一形量到 `truncated=1`／400）＋`TestAC2WindowsLeafIdentityReCheckRefusesAForeignNameAndNeverDeletesIt342r1`（＝instrument 红：两支叶子腿那发都没答，逐字 `no leaf leg answered (guard=0, landing's re-check=0) while reparseAt proved 83 flipped leaf cycles`）。
- ★基线树⛔ 有这两枚红＝这两枚红只在落地树上出现⇒⛔ "存量红"、⛔ "台面红"；但也⛔ 读成"我的字节把安全dim 打红＝落地失败"＝两枚红各自说的是**残余尚未被证明关掉**（损伤 1/400 与复验 0 次），正是 `AC#5` 要求写进注释而不是写成"已关"的那两件事。

## 竞态散布表（每形连打 3 发、逐次分子；竞态读数⛔ 要求逐字相同＝复现量级与形状）

POSIX · 票 342 截断支钉（相对祖先链接形，attempts=1500）：
| 侧 | cycles | provedLinkCycles | pinnedParentBeforeOpen | truncated | deleted | narrowed |
|---|---|---|---|---|---|---|
| pre r1/r2/r3 | 154/160/150 | 155/160/150 | 0/0/0 | 0/0/1 | 77/69/85 | 0 |
| land r1/r2/r3 | 164/160/201 | 164/160/202 | 26/20/36 | 0 | 0 | 0 |

POSIX · 票 348 删除支竞态钉（attempts=1500）：
| 侧 | provedLinkCycles | pinnedParentBeforeOpen | deleted | truncated |
|---|---|---|---|---|
| pre r1/r2/r3 | 131/134/158 | 0/0/0 | 68/67/73 | 0 |
| land r1/r2/r3 | 158/160/195 | 26/28/22 | 0 | 0 |

POSIX · 票 348 确定性钉（＝342-a1 的 A3 那一发由本腿自己重打＝把〔仅腿报〕升成"收件复跑"）：
- pre 3/3 红，红句逐字：`AC#2 RED: the foreign witness is GONE after the refusal-path unlink (lstat …\001\data\someone-elses-tree\leaf-342r1.bin: no such file or directory)`
- land 3/3 绿，同发逐字读回：`AC#2' survivor read-back: size 75 (was 75), mode -rw-rw-rw- (was -rw-rw-rw-), mtime moved=false`

Windows 三形（junction／symlink-dir／leaf-symlink，attempts=400×3 发）逐表逐字在 `readings-windows.md`；★要点＝`deleted` 落地侧 0（票 348 那一支）、`truncated` 在 `junction` 形 0、在 `symlink-dir` 形 1～2/400（**残余、⛔ 归因到底**）、`leaf-symlink` 形本落地复验 0 次作答（＝⛔ 被证明在路上）。

liveness 钥匙（两票 AC 逐字＝修后 0 必须与"窗口真到达过"同发）＝`provedLinkCycles`／`provedReparse` 与 `pinnedParentBeforeOpen`／`formBLegRefusals` 同时非 0，每行同发在场；⛔ 拿 0 命中当凭据、⛔ 一发代两发。

## 覆盖写四枚调用者结论（必答 (iv)）

逐枚现读调用点（取 `HEAD` blob）＝`internal/secret/store.go:79`（重写同 id 的 blob）／`internal/config/migrate.go:93`（旧构建留下的宽备份的修复路径）／`internal/secret/migrate.go:189`（崩溃遗留的临时文件被覆盖）／`internal/agent/spill.go:166`（`.retry` 临时文件被覆盖）＝4 枚依赖"允许覆盖写"；另 2 枚 `PrivateFileExclusive`（`agent/spill.go:343`、`memory/artifacts.go:75`）要的正是 `fs.ErrExist`。
⇒ 结论＝**本落形⛔ 动这一语义**：`winsec.go:96`/`:106` 两枚旗标一字未动，`PrivateFile` 的"or replaces"承诺与"名字已存在⛔ 构成拒绝理由"那段注释一字未动；新增的拒绝只来自"本包能证明这枚名字不安全"（祖先钉不住的换手），＝票面已有边界（322 的叶子支同形）。
凭据＝`openleaf_overwrite_322_test.go:70/:79/:107` 所在的 `TestAC2OverwritePromiseHoldsForEveryProductionCaller322r1`（含 exclusive 那一支）在 POSIX 容器两发与 Windows 本机两发里**全绿**（名册逐名见上＝无新增红）。

## 零新导出名自证（必答 (v)）

★先纠一把尺：本腿第一次用 PowerShell 的 `-match` 跑这条尺，命中 45 枚＝⛔ 能用、⛔ 报那个数（PowerShell 默认⛔ 区分大小写，`[A-Z]` 把 `func leafPin` 也算进去了）。改正后的尺＝大小写敏感（`-cmatch`／等价 `grep -E`）：
`git diff fde1df02..HEAD -- internal/winsec/` 里匹配 `^\+\s*(func|var|type|const) [A-Z]` ⇒ 命中 **6**，逐枚点名＝6 枚全是 `Test…` 用例入口（Go 约定测试入口必须大写起头、只在 `*_test.go`、⛔ 包外 API 面）。
⇒ **非测试的新顶层导出标识符＝0 枚**。新增的产码形全部非导出：`leafPin`／`pinnedParent`／`openPinnedLeaf`／`unlinkPinnedLeaf`／`pinVolumeRoot`／`sameObjectPinnedWindows`／`placementRefusalWindows`／`pinRefusal`／`entryIsLink`／`underlyingErrorOf`／`underlyingErrorMessage`／`openModeBits`／`errNoActTimePrimitive`。⇒ ⛔ 触到"必须新导出名才能闭环"那道停手上报门（票 340 同形）。

## 父目录句柄生命周期（必答 (ii)）＝怎么证明它活过两处删除那一跳

拿到＝`privateFile` 在 `preOpenLeafGuard` 之后、open 之前（`winsec.go:158` `parent, err := pinnedParent(resolved)`）；释放＝★路由末端唯一一处 `defer func() { _ = parent.close() }()`（`:163`），两处删除（`:174` 密封失败支、`:179` 写失败支）都在该 defer 之前发生；POSIX 侧每次 act 之后另有 `runtime.KeepAlive`（`openModeBits`/`unlinkPinnedLeaf`）把"⛔ 提前被 finalize 掉描述符"写成显式依赖。
仪器侧的证＝两票竞态钉的 `pinnedParentBeforeOpen` 那一列（land 侧 POSIX 26/20/36 与 26/28/22＝同一发里既证明窗口到达过、又证明是父钉在 open 之前答的）＋票 348 的确定性件那一行 survivor 读回。

## 本腿自抓的台损（5 枚，⛔ 计入判据，逐枚已具名修）

1. `go test` ⛔ 带 `-v` ⇒ 四数恒 `0/0/0/0`（尺失灵，⛔ 读成"全绿"）。
2. 容器 `sh -c "GOMODCACHE=/gomod … unset …; go test"` ⇒ 赋值只作用于 `unset`、`go test` 拿容器默认 `GOMODCACHE`⇒`go: downloading golang.org/x/sys`＋`[setup failed]`＋四数全 0（长得像判据响）；改 `export …; unset …`。
3. `os.WriteFile(witness, …, 0o666)` 被容器 umask(0022) 吃掉 ⇒ `narrowed` 两侧每发都响（首发 `narrowed=1423`）；改"建后 `Chmod`"。
4. Windows 相对链接的 staged 名建错层（相对 reparse 目标按**链接最终所在目录**解析）＋`toReal` ⛔ 分形（叶子形去移 `sub`⇒ENOTEMPTY）⇒ 首发两枚 `settle … no such file or directory`。
5. Windows 名册漏 `risk: path traverses a reparse point…`（＝已装好的 C26 resolver 的稳态拒绝）＋两支叶子腿没分开 ⇒ 300/400、373/400 落 `some other failure`，看着像"落地漏拒"、其实是仪器自己失真；归因＝探针两侧同一条串（见 `70-probe-*.log`），补名册时★⛔ 把它计入 `formBLegRefusals`（⛔ 拿别人的拒答冒充自己那支的 liveness）。

## ⛔ 尚未做（⛔ 把没做的说成做了）

1. ★去牙反控两发（票 342 `AC#3`／票 348 `AC#3`＝成对落点⛔ 相同：一发摘新加的祖先守卫⇒常驻钉必回落红；一发换成"逢名就拒"反形⇒覆盖写那四枚应红），只打仓外树物理字节、每发落刀前后验锚、跑前跑后验 pristine 哈希、报还原哈希、两发串行。
2. ★`winsec.go:118-140` 终态自陈（票 342 `AC#5`＝两平台都量到才许写"关"；Windows 的 `symlink-dir` 残余与叶子复验未证在路上⇒本腿⛔ 更新那段、只追加"Windows 残差"的写法在第 15 笔之后待做）。
3. ★Windows `symlink-dir` 那 1～2/400 的归因（两种候选见 `readings-windows.md`），以及把 `leaf-symlink` 的复验做成**确定性**件（借本包已声明的 seam 在判语与开之间插一次翻转，形同 `inIdentityWindow120r3`）。
4. 消费者面复认（`internal/risk`／`secret`／`config`／`agent`／`memory` 的 Linux 档与本机全 ok 成对）＝⛔ 跑（★同波裁腿 `349-v1` 在读 `internal/risk`；本腿⛔ 写它、⛔ 与它同刻跑同一包，撞红也只安静复量）。
