# `185-r3-head-criterion.md` — 票 185：常驻续读链判据（HEAD-first red 现量）＋ C25 声明归属裁定读数

- 腿代号：`185-r3`｜性质：**纯测试面**（写面＝`internal/tools/**` 的 `_test.go` ＋ 本表；`internal/risk/**` 零触碰）
- 起手＝终态锚：HEAD **`84b0a7c4d`**（本程 `git log -1` 现跑；`git status --porcelain -- internal/ cmd/` 起手即在飞件三枚，逐枚见 §①，本程未碰、未提交、未还原）
- ⚠ **对派单前提的一处更正（先测的结论，不是偏好）**：派单写「缺陷在主线上」——**在 `84b0a7c4d` 上这句话是假的**。修法本体 `internal/risk/hostpath_185.go`（A381 第四支，`ee1e2118c` 落地）在场，且生产盖章/名册路径自 `185-r2` 锚 `b88f3ef94` 以来**零 commit**（尺：`git log --oneline b88f3ef94..HEAD -- internal/risk/provenance.go internal/risk/hostpath_185.go internal/risk/taintmatch.go internal/tools/bridge.go internal/tools/fs.go internal/tools/task.go internal/tools/task_backfill.go internal/agent/spill.go` ＝ 空，唯一 tracked 差异是 `pointer_185_test.go` 的注释 24/−5）。本程按「补常驻判据＋裁定读数」执行，不是重新修一遍。

---

## ① 起手状态（现跑）

```
?? cmd/wisp/panel_resident_323_windows_test.go      <- 票 323 在飞件，未碰
?? cmd/wisp/panel_resident_assembly_323_test.go     <- 票 323 在飞件，未碰
?? internal/tools/ticket23_image_channel_test.go    <- 票 23 在飞件，未碰；整包可编译（G2 全绿为证）
```

## ② RED-before（未修码复刻，本程在 `84b0a7c4d` 上现量，非抄 185-r2）

复刻法：`-overlay` 把 `internal/risk/provenance.go` 换成只摘掉名册那一跳的副本（`:606` 的 `p.hostMintedPathFor(scopeID, origin)` → `""`，其余逐字＝HEAD；副本与 overlay 在 `.scratch/wisp/probes/185/r3/mut-unfix/`、`ov-unfix.json`）。

- `./internal/risk/`（`-run TestPointer185`）：**4 红**——`SecondRereadOfAHostArtifactIsClean`／`ForeignContentCarryingTheHostPathStillHits`／`RosterDiesWithTheMarksItWasDerivedFrom`／`CostFaceOfTheRosterIsPinnedAndNarrowerThanTheSpan`（与 185-r2 §②.1 四枚逐枚同名单）。
- `./internal/tools/`（`-run TestPointer185`）：**6 红**——seam 全家（与 185-r2 §②.2 同名单），红句含 `第二次照同一条宿主指针续读判成 "L2"`。
- ⇒ 派单那句话（task.output→fs.read 成功→同路径再发命中 R4）**在未修树上真、在 HEAD 上假**；入库判据的牙今天仍活（两包合计 10 红，逐枚可复）。

## ③ GREEN-after（HEAD 原样）

`-run TestPointer185` 两包全绿：risk `ok`、tools `ok`（逐字在 `probes/185/r3/logs/green-risk.txt`／`green-tools.txt`）。

## ④ 本腿落地的常驻判据：`internal/tools/pointer_185_reread_chain_test.go`（新件，两腿成对）

| 腿 | 钉什么 | HEAD | 未修复刻 | 候选(a)变异 |
|---|---|---|---|---|
| `TestPointer185LongRereadChainNeverRemintsTheBlocker` | task.output→**12 发**混合整份/有界续读全部 L0 且逐字节相等，链条后第 13 发仍 L0；mark 枚数＝1 桩＋13 读回＝14 **分毫不差**（豁免不许吞掉 175 的盖章义务） | **PASS** | **FAIL**（第 2 发即死） | **FAIL**（第 2 发即死） |
| `TestPointer185ChainDoesNotLaunderAForeignCarrier` | 长链条之后，外来正文逐字拼出宿主指针 ⇒ 续读**仍命中**且归因落**外来 mark**（票 177 W-2 在链条深度） | **PASS** | **FAIL**（precondition 死） | **FAIL**（precondition 死） |

红句逐字（两种变异同形，`pointer_185_reread_chain_test.go:77`）：
`第 2 发续读判成 "L2"（is_error=true）：一次成功的续读又亲手造出了下一次的阻断者（票 185 本体在深度 2 复活）`

⇒ 票题那个量词（「**每一次**成功的续读都造出下一次的阻断者」）从此有常驻反判据：旧判据只钉到深度 3（T1/T2），本判据钉到 12＋1，且对两种复活路径都红。

## ⑤ AC#2 裁定读数：「C25 声明该谁给——桥侧还是落盘层」

现量（本程逐条现跑，名册不从票面抄）：

- **盖章口**：全仓非测试 `MarkWithHostPath(` 生产调用唯一一处＝`internal/tools/bridge.go:752`；`box.set(` 唯一命中＝`internal/tools/task.go:559`（task.output 落盘回执）⇒ **今天带声明的工具仍只有 task.output 一枚**（与 185-c1 同值，两年未变）。
- **候选 (a) 桥侧逐发声明**（185-c1 定的形状：`fs.read` 成功时把读的那条路径填 `hostPathBox`）——本程用 `-overlay` 在 `fs.go` 成功返回前插 `box.set(canon)` 三行复量（`probes/185/r3/mut-claima/`）：**seam 6 枚全红**，含 T1 `第二次照同一条宿主指针续读判成 "L2"`。⇒ 候选 (a) 在 185-c1 的锚上「零牙」，**今天的牙＝6 枚红，且红法是机械性的**：`MarkWithHostPath` 里 `hostPath != ""` 直接关掉名册那一跳（`provenance.go:604` `if declaredNorm == "" && hostPath == ""`），而声明路径的全文拼写不在产物正文里 ⇒ 读回 mark 颗粒无收、**阻断者原地复活**。⇒ 桥侧逐发声明＝不但越契约（183 AC#5① 按值放行），连这道门自己都过不去。
- **候选 (b) 落盘层（D15）**：通道今天不存在（185-c1 尺 `os.ReadFile|os.Open` 在 `internal/agent` 非测试＝空，本程未复量通道、只复量承重面）：摘 `fs.read` 的戳会撞的三枚承重判据 `TestFSReadReturnsTheFileAndTaintsIt`／`TestFSReadTaintFeedsR4`／`TestEveryRegisteredToolIsClassifiedForMarking175r2` 在 HEAD **全绿**（`ok 0.095s`）＝票 175「成功结果必须有来源标记」今天仍被现役判据承着，落盘层任何「不产生第二枚 mark」的落法都要先拆它们＝禁令面。
- **裁定（读数支持，非偏好）**：**两边都不给。** C25 声明的唯一生产入口保持＝**桥的盖章步、只替 task.output 一枚**（宿主亲手落盘的那一处）；`fs.read` 读回正文的豁免由 **risk 层从 scope 自家 mark 派生**（`hostMintedPathFor`，A381 批准的第四支），不是新声明。桥侧逐发声明被本表 §④/⑤ 双变异判据钉死为「落了就红」；落盘层被票 175 三枚承重判据钉死为「拆了才通」。若日后要翻案，翻案者必须先让 `pointer_185_reread_chain_test.go` 两腿＋T4＋那三枚承重判据换形，属改契约文字，须人工批准。

## ⑥ 门禁（终态读数，取于 §⑦ 提交之前的工作树＝新判据在场＋在飞件在场）

| 门 | 读数 |
|---|---|
| `go build ./cmd/... ./internal/... ./tools/...` | `rc=0` |
| `go test ./internal/tools ./internal/risk -count=1 -timeout 420s` | `ok … internal/tools 22.304s`＋`ok … internal/risk 12.487s`（rc=0） |
| `gofumpt -l internal/tools/ internal/risk/` | **空**（rc=0） |
| `sh scripts/d22scan.sh` | `rc=0`，`clean - no D22 ban violations`（internal/ 906、cmd/ 213、frontend/ 116、design/ 39） |

（PATH 前缀按在册形 `$PWD/build:$PWD/third_party/sherpa-onnx`；`GOTMPDIR/TMPDIR/TEMP/TMP` 已 unset；未设 `CGO_ENABLED=0`。）

## ⑦ 本腿 tracked 改动（显式 pathspec 提交，零 push）

| 文件 | 改动 |
|---|---|
| `internal/tools/pointer_185_reread_chain_test.go` | 新增（两腿＋链条驱动器，零产码行） |
| `docs/evidence/s1/185-r3-head-criterion.md` | 新增（本表） |

探针与变异全在 `.scratch/wisp/probes/185/r3/**`（mut-unfix／mut-claima／ov-*.json／logs），不提交、不删。

## ⑧ 本程没测的（逐名，接着 185-r2 §⑦ 的账）

1. **AC#6 真机 CLI 那一发**（`probes/183/r1` 族的 `zz*` 台件）：写面不含 `cmd/wisp/**`，且其 sink 会原地重写别家票逐行引用的日志 ⇒ 仍未做，归面不变。
2. **AC#4 的 3/10 代价尺复量**（真实正文命中率）：未复量，本表不引新数；代价面的常驻判据（risk `CostFaceOfTheRosterIsPinnedAndNarrowerThanTheSpan`）在 G2 全包里绿。
3. `185-r1` 证据表 §②–⑫ 叙述格：仍空（别人的件，本腿不改写）。
4. 候选 (b)「宿主替模型读」新通道的落地实验：只量了承重面，未做通道原型（那要先拆 175，禁令面）。
