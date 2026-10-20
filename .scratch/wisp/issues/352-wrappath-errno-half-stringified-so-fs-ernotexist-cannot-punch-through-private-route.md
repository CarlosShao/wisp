# 票 352 — `privateFile` 路由把"父目录不存在"吞成 `ErrNotSealable`：`wrapPath` 的 errno 那半 `%v` 令 `errors.Is(err, fs.ErrNotExist)` 失效

**立票**：2026-10-20 15:1x 编排者（来路＝实现腿 `342-r1` 交件时上报的调用者可见行为差 → 终裁腿 `342-v2` 必答(v) 独立签"修回 errno 可穿"→ 机主口令「批」＝台账 `A1056`④/`A1058`；★本票＝把那道已批准的一步落成票面，⛔ 新批准）
**性质**：★错误分类可穿性缺陷（D37 错误模型⇄C26 边界语义＝契约邻接，故走人工批准，已批）。
**后果（对用户怎么说）**：调用方问"这路径不存在吗"，今天拿不到真话——`errors.Is(err, fs.ErrNotExist)`＝false，错误读起来像"这个对象密封不了"。今天全仓无真依赖它的调用者（终裁腿 38 个判定点逐枚现读＝零枚），所以无现红；但 `errors.Is` 是仓内唯一可编程通道，报反分类＝埋雷。

## 现量（每条带尺；⛔ 引用前先重跑）
- 真身（`git show HEAD:internal/winsec/winsec.go`）＝`:430 func wrapPath(path string, err error) error {` ⇒ 体内 `fmt.Errorf("%w: %s: %v", ErrNotSealable, path, err)`＝`%w` 只包 `ErrNotSealable`、errno 走 `%v` 字符串化 ⇒ 原始 errno 不可穿（尺＝`git grep -n 'func wrapPath'`＋现读函数体；★它在 winsec.go、⛔ 在 `winsec_other.go`——实现腿当时报错了位置，终裁腿复核过）。
- 调用点数＝`git cat-file blob HEAD:internal/winsec/winsec.go | grep -c 'wrapPath('`＝**2**（含定义行；实际调用点少，落地腿逐枚现数并点名）。
- 实测差（`342-v2` 交件⑥，仓外探针同时打进两棵树）：POSIX"父目录不存在"那一跳 基线 `isErrNotExist=true/isNotSealable=false` ⇒ 落地 `isErrNotExist=false/isNotSealable=true`；Windows 同跳 落地另带 `isUnresolvedPath=true`。★其余七跳两侧逐字同（覆盖写 err=nil、独占 `isErrExist=true`、叶子/祖先链接 `isUnresolvedPath`、空串/相对、`SealFile` 缺失两侧本就都 `ErrNotSealable`＝那跳从来⛔ 可穿）。
- ★**射程只到 `wrapPath` 这一处**：`pinnedParent` 的 Windows 支经 `placementRefusalWindows`→`preOpenRefusal` 已自带双 sentinel（源码注释逐字"pinPlacement carries both sentinels itself, so nothing here may strip them"）＝⛔ 第二处、⛔ 顺手扩。

## 要建什么
- [ ] **AC#0 现状钉死（改前读数）**：两平台各自量出"父目录不存在"那一跳今天 `errors.Is` 的两个真值（基线语义＝`fs.ErrNotExist` true／现语义＝`ErrNotExist` false＋`ErrNotSealable` true），逐字红/绿句入库＝修后要有对照物。
- [ ] **AC#1 落地（形状＝已批准的那一行）**：`wrapPath` 的 errno 那半 `%v`→`%w`（＝`fmt.Errorf("%w: %s: %w", ErrNotSealable, path, err)`，双 `%w` 需 Go≥1.20，本仓 go1.27）。★⛔ 动任何其他跳的分类、⛔ 改 `sealError` 幂等判定、⛔ 新增导出名、⛔ 改函数签名。判据＝POSIX⇄Windows 两形各量到 `errors.Is(err, fs.ErrNotExist)`＝**true** 且 `errors.Is(err, ErrNotSealable)` **照旧 true**（两分类并存＝这才叫"可穿"⛔"改判成另一种不可穿"）。
- [ ] **AC#2 消费者面不成对⛔ 结案**：`internal/secret`⇄`config`⇄`memory`⇄`agent`⇄`panel` 五包各自跑色（终裁腿那 38 枚判定点的归属包）＋★`cmd/wisp` 整包⛔ 本趟跑（归编排者那窗，防互抢）；落地⇄基线成对逐名作差＝新增红 0 枚。★既知连带具名：`private_fail_test.go:37/50/61/72/75` 五枚 `ErrNotSealable` 判定须仍绿；读 `ErrUnresolvedPath` 的 `openleaf_322`／`placement_symlink_113` ⛔ 受影响。
- [ ] **AC#3 纪律格**：⛔ 翻任何框（本票四格由非实现者裁）、⛔ 改票名／`Status:`／`-done`、⛔ push、名册尺锚定行首＋两把锚、显式 pathspec、新件先 `git add`、提交后对回。

## 禁区
⛔ 借本票顺手碰 `pinnedleaf_342_windows.go:11-20`/`:158-170` 的机制归因注释（＝`342-v3` 具名的〔丙〕＝另批）；⛔ 碰 `:118-140` 那段"only half closed"自陈；⛔ 新增任何产码接缝（`342-v3` 已签⛔ 该批＝本票只此一行）。

**Status:** **未开工**（机主已批、腿待派）。⛔ 零翻框、零 push。
## Progress log
（追加式；每条 `- [日期 时刻 +08] agent=… did=… next=…`，钟点由 `date` 插值，不手打。）
