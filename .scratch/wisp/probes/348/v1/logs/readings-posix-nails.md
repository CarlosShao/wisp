# 342-v2 读数 1/6 · POSIX 两枚新常驻钉独立重打（票 342 AC#1 ⇄ 票 348 AC#0/AC#2）

台面＝**自己** `git archive` 两棵仓外导出树 `D:\tmp\wisp\20261019\342v2\{treeB,treeL}`（⛔ 复用 `342-r1`／`342-v1` 任何一棵；⛔ `-overlay`）。
锚＝本腿第一笔 `git log -1` 现认 `c3c0ac80`；单差异两端＝基线 `fde1df02`（`winsec.go` blob `67783ae0`）⇄ 落地 `9581a48c`（blob `07300bff`）。
pristine 尺＝两棵树 `internal/winsec/*.go` 逐枚 sha256 并排差＝**5 行**（落地 4 枚新产码件＋`winsec.go` 换串）＝台面对称（六枚新钉两棵**同串**、⛔ 引用落地专属符号＝见 `logs-04-buildtrees.txt`）。
容器＝`golang:1.27`＋`CGO_ENABLED=0`＋`GOPROXY=off`＋宿主 module 缓存只读；★`TMPDIR`/`GOTMPDIR` **留容器内默认 `/tmp`**（⛔ 绑 Windows 卷＝`342-v1` 定式：`os.SameFile` 在 bind mount 失真＝身份门被台面摘掉）。`GOCACHE`＝D 卷宿主目录（跨发复用，⛔ 测试临时根）。
尺＝`go test ./internal/winsec/ -v -count=1 -timeout 40m -run '<三枚名>'`；四数＝`RUN/PASS_TOP/FAIL_TOP/SKIP`，⛔ 混用⛔ 相减。原始逐字＝`logs/i01-tree{B,L}-shot{1..4}.raw.txt`。★竞态读数⛔ 要求逐字相同＝复现量级与形状。

## 表一 · 票 342 那枚截断支钉 `TestAC1POSIX…NeverTruncatesAForeignFile342r1`（attempts=1500，target=rel）

| 侧 | 发 | cycles | provedLinkCycles | pinnedParentBeforeOpen | refusedAfterOpen | truncated | **deleted** | untouched | enoentGaps | flipErrs | 色 |
|---|---|---|---|---|---|---|---|---|---|---|---|
| 基线 treeB | 1 | 150 | 150 | **0** | 22 | 0 | **72** | 1428 | 131 | 17 | FAIL |
| 基线 treeB | 2 | 147 | 147 | **0** | 22 | 0 | **65** | 1435 | 116 | 21 | FAIL |
| 基线 treeB | 3 | 152 | 152 | **0** | 21 | 0 | **66** | 1434 | 116 | 19 | FAIL |
| 基线 treeB | 4 | 137 | 137 | **0** | 22 | 0 | **79** | 1421 | 119 | 18 | FAIL |
| 落地 treeL | 1 | 171 | 171 | **25** | 17 | 0 | **0** | 1500 | 108 | 31 | PASS |
| 落地 treeL | 2 | 158 | 158 | **22** | 9 | 0 | **0** | 1500 | 105 | 25 | PASS |
| 落地 treeL | 3 | 160 | 160 | **14** | 13 | 0 | **0** | 1500 | 124 | 22 | PASS |
| 落地 treeL | 4 | 169 | 169 | **29** | 8 | 0 | **0** | 1500 | 101 | 17 | PASS |

★基线侧那枚 `AC#1 RED instrument` 逐字（四发同形，shot1）＝`the racer proved 150 link-state cycles and the parent pin refused 0 attempts BEFORE any open (walk-after-open answers: 22, attempts: 1500, flipErrs=17, …)`＝**父钉那支在基线产码上根本不存在**（`pinnedParentBeforeOpen=0`）⇒ 那台的绿⛔ 可得、只有红可得。
★落地侧 liveness 钥匙**逐发在场**＝`pinnedParentBeforeOpen = 25/22/14/29` 与 `provedLinkCycles = 171/158/160/169` 同发非零 ⇒ 归零不是"没人到过窗口"（票 342 AC#1 逐字"修后 0 且 liveness 非零"＝**满足**）。

## 表二 · 票 348 竞态侧 `TestAC2POSIX…Race…NeverDeletesAForeignFile342r1`（同世界、`deleted` 是本票分子）

| 侧 | 发 | cycles | provedLinkCycles | pinnedWalkRefusals | pinnedParentBeforeOpen | **deleted** | truncated | narrowed | untouched | 色 |
|---|---|---|---|---|---|---|---|---|---|---|
| 基线 treeB | 1 | 144 | 144 | 32 | 0 | **86** | 0 | 0 | 1414 | FAIL |
| 基线 treeB | 2 | 131 | 131 | 25 | 0 | **71** | 0 | 0 | 1429 | FAIL |
| 基线 treeB | 3 | 135 | 135 | 22 | 0 | **73** | 0 | 0 | 1427 | FAIL |
| 基线 treeB | 4 | 138 | 139 | 20 | 0 | **73** | 0 | 0 | 1427 | FAIL |
| 落地 treeL | 1 | 159 | 159 | 28 | 22 | **0** | 0 | 0 | 1500 | PASS |
| 落地 treeL | 2 | 152 | 152 | 32 | 26 | **0** | 0 | 0 | 1500 | PASS |
| 落地 treeL | 3 | 146 | 147 | 43 | 25 | **0** | 0 | 0 | 1500 | PASS |
| 落地 treeL | 4 | 156 | 156 | 26 | 15 | **0** | 0 | 0 | 1500 | PASS |

基线侧分子＝**65~86/1500**（腿报 67/68/73~85＝量级与形状都复现）；落地侧 `deleted=0` ×4 且同发 liveness 非零。`narrowed=0` 两侧都 0（票 120 那一支在此形没被重新打开）。

## 表三 · 票 348 确定性侧 `TestAC2POSIXRefusalPath…NeverDeletesAForeignFile342r1`（腿报 pre 3/3 红 ⇄ land 3/3 绿；本腿重打 4 发）

- 基线 treeB 4/4 **红**，红句逐字（shot1）＝`AC#2 RED: the foreign witness is GONE after the refusal-path unlink (lstat …/001/data/someone-elses-tree/leaf-342r1.bin: no such file or directory), so the removal re-resolved the caller's spelling through the swapped ancestor and deleted a name in somebody else's tree - ticket 348's shape, the second damage of the pair` ＝与腿报名句**逐字同**（路径前缀是本腿台面的，⛔ 抄腿的树）。
- 落地 treeL 4/4 **绿**，绿侧逐字（shot1）＝`AC#2' survivor read-back: size 75 (was 75), mode -rw-rw-rw- (was -rw-rw-rw-), mtime moved=false`；同发拒绝逐字＝`AC#2' refusal verbatim (the leg that answered): winsec: path is not provably resolved, refusing to seal: …/data/sub/leaf-342r1.bin: pinned object at …/data/sub/leaf-342r1.bin is not what the name carries: this descriptor is not the object the name carries now, so the bytes would not land in what was just sealed (winsec: cannot apply a private security descriptor)`。
- ★这一发的"修前分子>0"在 POSIX **确定性成立**（4/4 删、⛔ 竞态运气），这就是票 348 `AC#0` 要的"把〔仅腿报〕升成本腿收件复跑"。

## 逐发四数（尺＝同一把，⛔ 相减）

`i01-treeB-shot{1..4}`＝`RUN=3 PASS_TOP=0 FAIL_TOP=3 SKIP=0` ×4 逐字同号；`i01-treeL-shot{1..4}`＝`RUN=3 PASS_TOP=3 FAIL_TOP=0 SKIP=0` ×4 逐字同号（`logs/four-counts-posix.tsv`）。

## 本腿自己看到的两枚仪器事实（⛔ 计入判据、供下一程）

1. `ancTargetAbs`（绝对目标那一形）在这三枚钉里**没被驱动**（两枚 POSIX 钉都取 `ancTargetRel`）⇒ 本腿的读数只覆盖 rel 形；世界件注释自陈 abs 形＝`Root` containment 先答、专门给票 348 的确定性形状供血。⛔ 外推成"两形都量到"。
2. `flipErrs` 基线侧 17~21／落地侧 14~31，首错逐字＝`remove …/data/sub: directory not empty`＝racer 的 clear-then-rename 在叶子刚被写入时扑空、随即自旋重试成功（`provedLinkCycles`≈`cycles`＝证到）。⛔ 把它当"窗口没到"。
