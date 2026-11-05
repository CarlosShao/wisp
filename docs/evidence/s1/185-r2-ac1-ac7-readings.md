# `185-r2-ac1-ac7-readings.md` — 票 185：AC#1 缺的那半（**未修码红**，两包）＋ AC#7 缺的那半（**牙现量**）

- 腿代号：`185-r2`｜性质：**纯测试面**（零产码改动）｜派单口径：写面只 `internal/risk/**`＋`internal/tools/**` 的 `_test.go` 与 `docs/evidence/s1/185-*.md`
- 起手 HEAD（本程 `git log -1`）：`b88f3ef94`（`ae62f761a` 是它的父提交＝派单写的锚）
- 终态读数取于本腿**最后一枚 commit 之后**（读数时 HEAD 见 §5 末行；期间别的写腿在推进主线，属别人的活，本程不碰、不提交、不还原）
- ⚠ **票面那句"缺陷在主线上"在本锚点上不成立**：修复本体已由 `ee1e2118c`（`185-r1`，A381 第四支，owner 2026-09-28 16:5x 批「好的，那超长输出，我批了，要做的」）落进主线。本腿做的是票 185 两格**欠下的读数**，不是重新实现一遍修法。
- ⚠ **AC 框一枚没勾**（票面文件不在本腿写面；勾/改名归编排者）

---

## ① 起手状态（逐枚现跑）

```
$ git log -1 --format='%h %ad %s' --date=iso
b88f3ef94 ... 收wave21:103 AC#3③PASS/102已结案/retriage-189;A1557

$ git status --porcelain -- internal/ cmd/
?? internal/tools/ticket23_image_channel_test.go      <- 另一枚写腿的在飞件，本程未碰
```

修复在位的尺（现跑）：

```
$ grep -rn "hostMintedPathFor" internal/risk/*.go | grep -v _test.go
internal/risk/provenance.go:606:  if v := p.hostMintedPathFor(scopeID, origin); v != "" {
internal/risk/hostpath_185.go:66: func (p *Provenance) hostMintedPathFor(scopeID, origin string) string {

$ grep -rn "MarkWithHostPath(" --include=*.go internal/ cmd/ | grep -v _test.go
internal/risk/provenance.go:484   (Mark -> MarkWithHostPath, hostPath="")
internal/risk/provenance.go:563   (定义)
internal/tools/bridge.go:752      (唯一生产盖章点)
```

⇒ 第二条尺今天与 185-c1 量的**同值**：生产盖章口只有桥一处、带声明的工具仍只有 `task.output` 一支（`internal/tools/task.go:558-559` 的 `hostPathBoxFromCtx`／`box.set`）。本票**没有**把 `fs.read` 也加进声明（那是 A381 排除的"参数侧按值放行"）。

## ② AC#1 的**红前**读数（判据要成对；两包分取）

方法：`-overlay` 把 `internal/risk/provenance.go` 换成一枚**只摘掉名册那一跳**的副本（`provenance.go:606` 的 `if v := p.hostMintedPathFor(...)` 换成 `if v := ""`，其余字节逐字＝HEAD），即票 185 立案时那棵树的行为；**跟踪件一字节未动**。

### ②.1 包级（`./internal/risk/`，26 枚：`pointer_183` 8＋`pointer_185` 10＋`shape_a` 8）

**4 红 / 22 绿**，红句逐字：

```
--- FAIL: TestPointer185SecondRereadOfAHostArtifactIsClean (0.01s)
    pointer_185_test.go:142: 185 referee: the host's own artifact must stay re-readABLE over and over,
      got hit from origin="C:\\Users\\swq\\AppData\\Roaming\\wisp\\artifacts\\tool-output-agent-task-c66c0634.txt"
      fragment="c:\\users"
--- FAIL: TestPointer185ForeignContentCarryingTheHostPathStillHits (0.01s)
    pointer_185_test.go:211: W-2: the hit must be attributed to the foreign mark, got "fs.read"
      (origin="C:\\Users\\swq\\...\\tool-output-agent-task-c66c0634.txt")
--- FAIL: TestPointer185RosterDiesWithTheMarksItWasDerivedFrom (0.01s)
    pointer_185_test.go:258: precondition broken: the referee's own scope must be the clean one before Close means anything
--- FAIL: TestPointer185CostFaceOfTheRosterIsPinnedAndNarrowerThanTheSpan (0.01s)
    pointer_185_test.go:381: cost leg broken: the rostered mark carries no value rule at all, so this leg measures nothing
```

⚠ 第二枚红是**归因那半**（不是"没拦住"，是"拦住却记在宿主自己那枚上"）：这就是 183-v1 的"归因本身也是判据"，也说明**成对判据的"反"那一半在未修码上红的是归因、不是颜色**——与票 183 的记录同形。

### ②.2 接缝（`./internal/tools/`，真桥具：真 `TaskRoster`＋真 `agent.Spiller`＋20000 字节真产物）

**6 红 / 0 绿**（整包口径：这 6 枚是本票的接缝判据全集），红句逐字：

```
--- FAIL: TestPointer185SecondRereadOfTheSameArtifactStaysClean
    pointer_185_cli_seam_test.go:212: 第二次照同一条宿主指针续读判成 "L2", want L0
      （每一次成功续读又造出一枚没声明的 mark＝票 185 本体）
--- FAIL: TestPointer185PagedRereadsKeepWorking
    pointer_185_cli_seam_test.go:232: 第 2 发分页续读判成 "L2"（is_error=true）：分页读今天还是走不通
--- FAIL: TestPointer185ForeignFileCarryingTheHostPathStillBlocksTheReread
    pointer_185_cli_seam_test.go:279: 归因错：这一发的证据判成来自 "<本测试临时根>\001\artifacts\tool-output-agent-task-c66c0634185r1a1b.txt",
      want 外来正文 "<...>\001\foreign-note.txt"
--- FAIL: TestPointer185AModelNominatedFileIsNotRostered
    pointer_185_cli_seam_test.go:302: precondition 对照失效：宿主自己落盘的那一枚第 2 发被拒（"L2"），
      本腿测的就不是名册边界而是全拒
--- FAIL: TestPointer185AnotherTaskCannotBorrowTheRoster
    pointer_185_cli_seam_test.go:340: precondition: 声明方的第二次没通（"L2"），T1 那一格也一起红了
--- FAIL: TestPointer185SiblingHostArtifactIsNotLaundered
    pointer_185_cli_seam_test.go:390: 声明方第 2 发被拒（"L2"）：T1/T2 那一格也一起红了
```

⇒ **票 185 本体在接缝上可复现**（第 212／232 两行就是"第二发 L2"），不是只有包级推理。

## ③ AC#1 的**绿后**读数（同一套判据、同一棵树）

```
$ go test ./internal/risk/ -run 'TestPointer18|TestShapeA' -count=1 -v
26 PASS / 0 FAIL        ok  github.com/CarlosShao/wisp/internal/risk

$ go test ./internal/tools/ -run 'TestPointer185' -count=1 -v
6 PASS / 0 FAIL         ok  github.com/CarlosShao/wisp/internal/tools
      pointer_185_cli_seam_test.go:242: 读数：三发分页续读后 scope=185r1-paged 的 mark 枚数=4
```

成对两半都在：正（同一路径第二发／第三发不命中 R4）＋反（外来正文携带同一条宿主路径**仍**命中，且归因落在外来那一枚）。⇒ AC#1 的判据**已入库且成对**，缺的只是这两组读数——本表补上。

## ④ AC#7 的**牙现量**（票面原话："它的牙必须现量"）

被测判据：`TestPointer185SameBasenameOtherDirectoryStillHits`（同名不同目录：`D:\wisp\other\artifacts\tool-output-agent-task-c66c0634.txt` vs 声明的 `C:\Users\swq\AppData\Roaming\wisp\artifacts\tool-output-agent-task-c66c0634.txt`，宿主正文里逐字拼出兄弟路径）。尺：`grep -rn "ptr185SameBasename" internal/risk/pointer_185_test.go`。

**变异甲（真牙：价值规则放宽到 basename 级，判据在 `contains` 里）**——`-overlay` 把 `internal/risk/taintmatch.go` 换成一枚副本：当**候选**自己的 basename 出现在本 mark 的 `declaredPath` 里，整个候选不再算证据（`grep -n "func (f \*fragmentIndex) contains" internal/risk/taintmatch.go`）：

```
1 红 / 25 绿（26 枚全集）
--- FAIL: TestPointer185SameBasenameOtherDirectoryStillHits
其余 25 枚（pointer_183 8＋pointer_185 其余 9＋shape_a 8）全 PASS
```

⇒ 这一格**有牙**，且牙是它独有的：`pointer_183` 8 枚＋`shape_a` 8 枚在该变异下**一枚都不红**（与 185-c1 的复量一致）。

**变异乙（`185-r1` 表 §8 自称的 M1：名册键放宽成"最后 24 个 rune"）**——`-overlay` 换 `internal/risk/hostpath_185.go`（`0 红 / 26 绿`）：

```
26 PASS / 0 FAIL
```

⇒ **M1 不是这一格的牙**（`185-r1` 的注释"its tooth is table §8's M1"与现量相反）。放松名册键动的是"哪一枚 mark 被戴上价值规则"，动不到这一格检查的那条候选。已就地更正该注释（本腿唯一 tracked 改动，只改注释）。

**变异丙（对照：这一格在"未修码"上是绿的）**——用 §② 那枚摘名册副本跑：

```
--- PASS: TestPointer185SameBasenameOtherDirectoryStillHits
```

⇒ 与 AC#7 票面"今天应当是绿的＝不许弄坏的守卫，不许充当本票新增的牙"**逐字吻合**；它的牙只能靠变异甲现量（本表就是那枚读数）。

## ⑤ 门禁（本腿最后一枚 commit 之后复跑）

| 门 | 命令 | 读数 |
|---|---|---|
| build | `go build ./cmd/... ./internal/... ./tools/...` | `rc=0` |
| risk 全包 | `go test ./internal/risk/ -count=1 -timeout 420s` | `ok … 14.700s`（rc=0） |
| tools 全包 | `go test ./internal/tools/ -count=1 -timeout 420s` | `ok … 29.638s`（rc=0） |
| gofumpt | `/d/work/base/gopath/bin/gofumpt.exe -l internal/risk/ internal/tools/` | **空**（rc=0） |
| d22scan | `sh scripts/d22scan.sh` | `d22rc=0`：`clean - no D22 ban violations`；`examined 421 production Go files`；`ban #8 internal/ 905`／`cmd/ 211`／`frontend/ 116`／`design/ 39` |

（PATH 前缀按在册形：`export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"`；`GOTMPDIR/TMPDIR/TEMP/TMP` 已 unset。）

## ⑥ 契约轴（票 185 AC#6）：**冻结文字内容未变，但"按行号取 md5"那两把尺在本锚点已过期**

`185-r1` 记的基线用**行号窗口**取；本锚点行号已被别家写腿顶移，逐字复跑**对不上**：

```
$ sed -n '468,473p' internal/risk/provenance.go | md5sum
447a92b2c856b114f0e2cf2aef0f2333        (基线 858e45116383caa3e7c1dd4b0924fad1)
$ sed -n '11,15p'   internal/risk/taintmatch.go | md5sum
5680ddd18e2d2ec2a85e485b54f4c12e        (基线同值，命中)
$ sed -n '482,506p' internal/risk/provenance.go | md5sum
a0c379a2c52083a26f04e5db431f8325        (基线 f89e891e5eee3f3ea2b4f89d921072c4)
```

**按内容锚点重取（这才是可复跑的尺）**：在 `ee1e2118c`（`185-r1` 落地那枚）的 blob 上，`468-473`／`482-506` 两窗**逐字命中**基线；同一段文字在本锚点上是：

```
Mark 契约块（6 行，`grep -n "^// Mark records one sensitive-source read"` 起）
  ee1e2118c：fcb9105f8222c0a2e28a2806576f0cde     HEAD：fcb9105f8222c0a2e28a2806576f0cde   <- 同一段
183 声明块（`grep -n "^// MarkWithHostPath is Mark plus the ticket 177 shape-A declaration"` 到
           `grep -n "never a wider index\."`）
  ee1e2118c：bc76bbe6c39f798aaf9848ee2b8a9e17     HEAD：bc76bbe6c39f798aaf9848ee2b8a9e17   <- 同一段
```

⇒ 结论：**冻结两句一字未动**，票面要求的"要更正只许追加"也守住了（185 的更正段在 183 那段**之后**追加，见 `provenance.go` "Delivery semantic correction, ticket 185 (ledger A381...)"）。⚠ 但请编排者把 `A363`／`A364` 那两把**行号**基线改成**内容锚点**基线，否则任何人下次复跑都会看到"两枚对不上"的假红。

## ⑦ 本程**没测**的（逐名，不许当已做）

1. **AC#6 真机 CLI 那一发**（`probes/183/r1` 那族的台件 `zz183r1_e2e_test.go`）：本腿写面不含 `cmd/wisp/**`，且那枚台件的 sink 会**原地重写** `probes/183/r1/logs/e2e-readings.txt`（票 185 现量段与 `A364` 逐行引第 55／60／64 行）⇒ 按 `A364` 的定式，重跑前必须先把它复制到自己的路径，属 `185-r2` 原排程的活，本腿不顶替。
2. **AC#4／§⑥ 的 `3/10` 代价尺复量**（真实正文命中率，语料＝8 枚仓内文件＋2 枚自造）：未复量，本表不引它的数。
3. **`185-r1` 证据件 `docs/evidence/s1/185-paged-reread-r1.md` 的 §②-§⑫ 叙述格**：仍空。本腿**不改写别人的件**（该件是 `185-r1` 的交付表），把读数落在这枚新件里；要合并由编排者裁。
4. **面板／前端**：零触碰（不是本腿写面）。

## ⑧ 对票面与对编排者的**更正**（本腿唯一不服）

1. **票面状态过期**：`Status: 已立、按住不派` 与 `Progress log` 停在 09-28 14:4x，而主线已在 09-28 17:0x 由 `ee1e2118c` 落地修法（`A385` 收下、明写"票 185 不结案"）。⇒ **"缺陷在主线上"这句话对今天的树是假的**，本腿按"补读数"而不是"重新修一遍"来做；若下一步还要"实现第一格"，唯一还缺的是 §⑦ 那三件（真机／代价尺／别人那张表），**都不是产码**。
2. **`185-r1` 注释里的牙归属是错的**（M1 → 实为 basename 级价值规则变异）：已在本腿唯一 tracked 改动里就地更正，并写了可复跑的尺。这一条属"注释说的动作没有判据钉着"同族，是票 183／185 反复付学费的那一类。
3. **未勾任何 AC 框、未改票面、未 push、零删除命令**。

## ⑨ 本腿 tracked 改动（逐枚）

| 文件 | 改动 | 说明 |
|---|---|---|
| `internal/risk/pointer_185_test.go` | 注释 +24/−5（**零代码行**） | 把 AC#7 那格的"牙"从错的 M1 更正为实测的 basename 级价值规则变异，并指向本表 §4 |
| `docs/evidence/s1/185-r2-ac1-ac7-readings.md` | 新增（本件） | 上面全部读数 |
