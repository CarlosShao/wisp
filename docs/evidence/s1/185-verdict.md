# `185-verdict.md` — 票 185 **非实现者终裁表**（185-v1｜每格 fresh 读数，全部本程自跑）

- 腿代号：`185-v1`｜性质：**终裁（非实现者）**｜写面＝本表一枚；产码/票面/台账零触碰、零 push
- 读数时点：2026-11-06 12:27–12:52（全部本程现跑；探针/replica/日志在 `.scratch/wisp/probes/185/v1/**`，未提交）
- 前置阅读：票面 185 全文＋`docs/evidence/s1/185-r3-head-criterion.md`＋`internal/tools/pointer_185_reread_chain_test.go`（均已读）

---

## ① 锚点与起手状态

**锚点更正（先测，不是偏好）**：派单写 HEAD `abf2a35a6`。本程起手实为 **`8879c3d7b`**（＝`abf2a35a6` 之上一枚台账收账 commit `A1567`，只动 `docs/reports/pending-and-issues.md` +7 行）；本程进行中主线又推进两枚（`55f3d4dea` 票98证据件、`5e344d27c` 票177-r2证据件）。**代码面零漂移的尺**：`git diff --name-only abf2a35a6..HEAD -- internal/ cmd/ tools/` ＝ **空** ⇒ 本程在 `8879c3d7b` 上取的全部读数对 `5e344d27c` 逐字有效。

**上下文**：台账 `A1567` 记「V185 null 阵亡（无产出，wave31 重派）」——本表即 wave31 重派的交付。

**起手 `git status --porcelain -- internal/ cmd/`**（排程护栏，逐枚具名，全部未碰）：

```
?? cmd/wisp/panel_resident_323_windows_test.go      （票323 在飞）
?? cmd/wisp/panel_resident_assembly_323_test.go     （票323 在飞）
?? cmd/wisp/ticket123_replyleg_other_test.go        （W123g 在飞）
?? cmd/wisp/ticket123_replyleg_test.go              （W123g 在飞）
?? cmd/wisp/ticket123_replyleg_windows_test.go      （W123g 在飞）
?? internal/tools/ticket23_image_channel_test.go    （票23 在飞）
```

**第21处 自证**：本程四门与全部判据读数**没有一条红可归因到未跟踪在飞件**（唯一两处红——未修复刻与两枚变异——都是我自己的 `-overlay` replica 造的，且经 `diff` 逐行核过）；仍按派单机制补做了 `git clone --shared` 纯态导出树复核（§⑦）。

---

## ② 门禁电池（派单四门＋补强，本程自跑）

| 门 | 命令 | 读数 |
|---|---|---|
| G1 | `go build ./cmd/... ./internal/... ./tools/...` | **rc=0**（输出空） |
| G2 | `go test ./internal/tools ./internal/risk -count=1 -timeout 420s` | **rc=0**：`ok internal/tools 26.822s`＋`ok internal/risk 13.293s` |
| G3 | `gofumpt -l internal/tools/ internal/risk/`（及宽扫 `internal/ cmd/ tools/`） | **双空**（rc=0） |
| G4 | `sh scripts/d22scan.sh` | **rc=0**：阳性对照 `PASS=35 FAIL=0 SKIP=0`；实扫 `clean - no D22 ban violations`；examined **421** production Go files；ban#8 internal/=**907**／cmd/=**216**（r3 记 906/213，差值＝别腿在飞件，非本票） |
| 补 | `go test ./internal/agent -count=1` | `ok 7.361s` |
| 补 | risk 单跑／tools 单跑 | `ok 13.956s`／`ok 30.931s` |
| 补 | **纯态**（`git clone --shared` 导出树 `/tmp/wisp185v1-export-8879`，checkout `8879c3d7b`） | build rc=0；`internal/risk ok 14.727s`；`internal/tools` **仅一枚环境形红**＝`TestTicket252R1AlignmentAddsNoAuthorization`（见 §⑦尾）；**票185 族纯态 risk 10/10＋tools 8/8 全 PASS** |

（PATH 前缀按在册形 `"$PWD/build:$PWD/third_party/sherpa-onnx"`；`GOTMPDIR/TMPDIR/TEMP/TMP` 已 unset；未设 `CGO_ENABLED=0`。）

---

## ③ AC#1 —— **PASS**（红前/绿后判据已落地、成对、有牙；本程自铸 replica 复量）

**落地件**（全部 tracked 在 HEAD）：`internal/risk/pointer_185_test.go`（**10 腿**：裁判正反两半＋跨 scope＋名册随 mark 生死＋同目录兄弟＋一 rune 之差＋**同名不同目录**＋代价面钉）＋ `internal/tools/pointer_185_cli_seam_test.go`（T1–T6，真桥真 Spiller 真 20000 字节产物）＋ `internal/tools/pointer_185_reread_chain_test.go`（r3 新增 2 腿，`git ls-files` 证 tracked）。

**红前（未修复刻，本程自铸）**：`-overlay` 把 `internal/risk/provenance.go` 换成只摘名册那一跳的副本（`:606` 的 `if v := p.hostMintedPathFor(scopeID, origin)` → `if v := ""`；**副本 diff 恰 1 行**，其余逐字＝HEAD）：

- `./internal/risk -run TestPointer185`：**4 红**——`SecondRereadOfAHostArtifactIsClean`／`ForeignContentCarryingTheHostPathStillHits`／`RosterDiesWithTheMarksItWasDerivedFrom`／`CostFaceOfTheRosterIsPinnedAndNarrowerThanTheSpan`（与 185-r2 §②.1 名单逐枚同名）。
- `./internal/tools -run TestPointer185`：**8 红**——seam 6 枚＋**链双腿**。链腿红句逐字（`pointer_185_reread_chain_test.go:77`／`:106`）：
  `第 2 发续读判成 "L2"（is_error=true）：一次成功的续读又亲手造出了下一次的阻断者（票 185 本体在深度 2 复活）`
  ⚠ r2/r3 记 6 红，本程 8 红＝**链判据比 r2 时代多了两枚真牙**：新判据自己在未修复刻上真红，不是挂在旧牙上的转述。

**绿后（HEAD 原样）**：同过滤全 PASS（链双腿 `-v` 单跑 PASS；纯态克隆里 185 族 risk 10/10＋tools 8/8 PASS）。

**成对性**：正半（同指针第二发/12+1 发链全 L0 逐字节）＋反半 W-2（外来正文携带同路径**仍命中**且归因落外来 mark：risk `ForeignContent…`＋seam `T3`＋链 `ChainDoesNotLaunderAForeignCarrier`）——两半都在场，缺一形的风险不存在。

## ④ AC#2 —— **PASS**（「声明该谁给」有裁定读数，非偏好；本程双变异复证）

- **裁定**（r3 §⑤ 记录，本程认可）：C25 声明的唯一生产入口保持＝桥的盖章步、只替 `task.output` 一枚；`fs.read` 读回正文的豁免由 risk 层从 scope 自家 mark **派生**（`hostpath_185.go` 的 `hostMintedPathFor`，A381 第四支），不是新声明。**两边都不给**。
- **候选 (a) 桥侧逐发声明＝落了就红（本程复量）**：`-overlay` 把 `internal/tools/fs.go` 换成成功返回前插 `box.set(canon)` 三行的副本（**diff 恰 3 行**）⇒ `./internal/tools -run TestPointer185` **8 红**（seam 6＋链 2）。⚠ 185-c1 在它的锚上量过同形变异＝**零牙**；今天的 8 枚牙全是 185 三腿落地的——「这支今天没人看着」的局面已被本票自己堵上。
- **候选 (b) 落盘层＝仍在禁令面后（本程复量）**：摘 `fs.read` 的戳会撞的三枚承重判据 `TestFSReadReturnsTheFileAndTaintsIt`／`TestFSReadTaintFeedsR4`／`TestEveryRegisteredToolIsClassifiedForMarking175r2` 在 HEAD **全 PASS**（`ok 0.083s`）＝票 175「成功结果必须有来源标记」仍被现役判据承着。
- **停手上报被遵守**：c1 普查四支全部撞边界即停（未擅自落地）；owner 随后 `A381` 批「好的，那超长输出，我批了，要做的」并要求改掉 taintmatch 那句 "NOT approved"——r1 照批改（ sanction 出处＝A381 本身）。

## ⑤ AC#3 —— **PASS**（覆盖面名册本程现跑，非抄票面）

- `grep -rn "MarkWithHostPath\|hostPathBox" --include=*.go internal/ cmd/`（非测试）：生产盖章点唯一＝`bridge.go:752`（`provenance.go:484` 是 `Mark` 转发 wrapper、`:563` 是定义）；`box.set` 唯一命中＝`task.go:559` ⇒ **带声明的工具今天仍只有 `task.output` 一枚**（与 c1/r2/r3 同值，本程独立复得）。
- 盖章名册 `sensitiveSourceTools`＝9 枚（fs.read／search.content／clipboard.read／system.get／web.fetch／doc.read／screen.capture／transcript／task.output）。豁免按 mark **自己的 origin** 判（不看工具名、不看候选参数），所以任何 sensitive-source 工具对宿主落盘产物**自己的读回**都被覆盖；只在回执里**携带**路径的非读回（如 `fs.list` 列目录）仍算证据——墙没有放宽。
- `fs.list` 不在盖章名册（c1 的反向面）原样在册：少一张证据、不是阻断者，登记给下一轮，本票不裁。

## ⑥ AC#4 —— **PARTIAL（半格）**：退让面已具名＋量化钉在场；**3/10 尺今日树未复量**

- **已交付**：①退让与收益同批具名——`hostpath_185.go` 文档块 + 台账 `A385` 逐字（「往后凡宿主亲手落盘、又被 fs.read 整份读回的那枚 mark，其正文里恰好连续拼出该产物路径 ≥8 rune 的片段…不再算续读证据；同一段文字若来自 web.fetch 或模型自己提名，仍算」）；②量化常驻钉 `TestPointer185CostFaceOfTheRosterIsPinnedAndNarrowerThanTheSpan` **本程 PASS**（rostered mark 的证据窗集永远 ⊇ 同体声明 mark，宽不过 183）。
- **欠的一格**：3/10 真机频率尺在**修复后树上**的复量。r1 证据表 §⑥ 至今「（待填）」，r2 §⑦.2／r3 §⑧.2 均明记未复量。本程重放原探针（`probes/183/r1/overlay-freq-head.json`，9-28 未修复码副本）⇒ **build failed**：那批副本缺 HEAD 包现已定义的 `declaredPath`（taintmatch）与 `markObs`（394 seam）字段——**时代不匹配，原样不可重放**，读数留档 `p10-freq-replay.txt`。
- **裁定**：AC 的实质（不许只报收益不报退让）已满足；缺的是一个**语料数字**，不是产码。**框：带注可翻，或持开**（见 §⑩）。

## ⑦ AC#5 —— **PASS**（四门全绿 fresh；CLI 那一格本程补测两形 PASS；两条残留如实记）

- 四门＋agent 包＋单包读数见 §②，全绿。
- **CLI 真机那一格（r1 §⑦/r2 §⑦.1/r3 §⑧.1 三腿都欠的那格）本程补上**：
  - 方法：把 9-28 的 `zz183r1_e2e_test.go` **移植成自己的副本**（`.scratch/wisp/probes/185/v1/zz185v1_e2e_test.go`，三处机械修：`+context` import；`rt.execute` 按现行签名补 ctx 参 `run.go:1539`；config 模板按现行 schema 给 m1 补 `enabled = true`；sink 改道我自己 `probes/185/v1/logs/`）。**被引读的 183 日志零触碰**：`probes/183/r1/logs/e2e-readings.txt` 跑前跑后 md5 同为 `c2ac391f80c5d072a1f7f393f51416ac`（跑前备份在 `e2e-readings.BEFORE-20261106.txt`）。
  - 读数（**两形都贴**）：带在册 PATH 前缀形 **PASS**；不带前缀形 **PASS**。两形一致：run A 退出 0／run B 退出 0／确认卡 0→0；整份读回**逐字节相等＝true**（20000 字节）；`task.output` kind=success L0；**第一发 fs.read kind=success risk=L0 rules_hit=[]**；**第二发（有界）fs.read kind=success risk=L0 rules_hit=[]**；分岔计数 `R4 console=0 rig-turns=0 tool-rows=0`（没有任何一发 refused）；C25 收尾 **dropped=3**（1 桩＋2 读回——修复前同台件是 refused＋dropped=2）。⇒ **票 185 本体在真 CLI 接缝上消失，两形都量到了**。
  - ⚠ 条件如实记：该读数是**工作树条件**（在飞 W123g 未跟踪件在场）——纯态 `go test ./cmd/wisp` 的编译缺口已由票 98 腿登记（abf2a35a6 引用未提交定义，归票 123/W123g，第21处镜像形状），非本票可修。
- **gate-clauses.sh**：**不是仓内脚本**（git 历史零命中；只有 `.scratch/wisp/probes/154/` 探针件）。本程复跑：仪器自身腿数断言过（名册=14 声明=14 记账=14 缺腿=0），但 10 枚腿的 git grep 在本程嵌套 harness 里 rc=126（手工等价命令 rc=0＝harness 伪象，非仓态信号）；**名册最新鲜的有效读数仍是 185-c1**（BAD=1=G6neg、与在册差集空）。
- 纯态克隆复核（§② 末行）：185 族全绿；唯一红＝票 252 的环境形探针 `TestTicket252R1AlignmentAddsNoAuthorization`——克隆根在 %TEMP% 下取不到 8.3 短名别名，前置自失败（该腿明写「不许 skip」故响）——**与票 185 无关、与未跟踪件无关**，归票 252 的机器依赖面。

## ⑧ AC#6 —— **PASS**（冻结文字内容锚逐字等值；185 五枚 commit 写面零越界）

- **内容锚（本程自跑）**：`provenance.go` 从 `// Mark records one sensitive-source read` 到 `never a wider index.` 整段——`ee1e2118c` 与 HEAD **逐字节相同**（该段同时覆盖 Mark 契约块与票 183 声明块 25 行；`diff` 输出 `MARK BLOCK IDENTICAL`，md5 `0f7633e659657b0b7dca8c548241eb43` 两侧同值）。185 的更正段是**追加**在该段之后（A385 语义更正），符合「要更正只许追加」。
- `internal/risk/taintmatch.go:11-15`（逐 token 追踪 REJECTED）：md5 **`5680ddd18e2d2ec2a85e485b54f4c12e`** ＝A363/A364 基线原值（该窗行号未漂）。
- **185 五枚 commit 的文件清单**（ee1e2118c／7f6d7a26a／567c7875a／f606fe4be／b7157ac19）：只有 `internal/risk/*`（hostpath_185.go、pointer_185_test.go、provenance.go、taintmatch.go）、`internal/tools/*`（两枚测试件）、`docs/evidence/s1/185-*.md`、`.scratch/wisp/probes/**`。**零触碰**：`docs/PLAN.md`／`docs/specs/**`／thresholds／golden／审批超时常量／`allowlist.txt`。taintmatch 那句名册改写＝A381 批准本身要求的，非越界。
- 沿袭 r2 §⑥ 的提醒（本程再证）：A363/A364 的**行号**基线（provenance.go:468-473/482-506）已漂移失效，内容锚才是可复跑的尺——台账基线换锚一事仍悬，归编排者。

## ⑨ AC#7 —— **PASS，牙活（本程自铸变异现量，独红一枚）**

- 守卫绿：`TestPointer185SameBasenameOtherDirectoryStillHits`（落点＝`internal/risk/pointer_185_test.go` 第 9 枚/共 10——票面点名「第 9 枚」那支，rail 尺 `grep -c "^func Test" internal/risk/pointer_183_test.go`＝**8** 未变）在全部全包/族跑里 PASS。
- 牙（本程自写变异，非抄 r2 件）：`-overlay` 把 `internal/risk/taintmatch.go` 换成「价值规则从窗口级放宽到 **basename 级**」的副本（`fragmentIndex.contains` 里：候选自己的 basename 出现在该 mark 的 declaredPath ⇒ 整个候选不算证据；**副本 diff 恰 16 行插入，跑前核过**）：
  - `-run TestPointer185`：**独红 1 枚**＝`TestPointer185SameBasenameOtherDirectoryStillHits`，其余全绿；
  - **整包对照**（全 `./internal/risk`，同一变异）：仍**只有这 1 红**（pointer_183 8＋shape_a 8＋pointer_185 其余 9 全绿）⇒ **牙是它独有的、活的**——「豁免改成只比 basename」那一发今天会被入库判据当场抓住。
  - 纪律自报：本变异第一发探针跑在**未真正施加变异的副本**上（0 红＝no-op，作废）；重做后先 `diff` 核实 16 行在位再跑，红读数才取。两次日志都在（`p4-basename-risk.txt` 作废／`p4r` 为准），不删。

---

## ⑩ 终裁表＋可结性

| AC | 裁决 | 附条件/残留 |
|---|---|---|
| AC#1 红前/绿后判据落地（成对） | **PASS** | 无——本程自铸 replica：risk 4 红＋tools 8 红（含链双腿），绿后全过 |
| AC#2 声明归属裁定（逐支现量） | **PASS** | 无——候选(a)变异 8 红、175 三承重绿，均本程复量 |
| AC#3 覆盖面名册现跑 | **PASS** | 无——盖章点/声明工具与 c1 同值，独立复得 |
| AC#4 代价面复量＋具名 | **PARTIAL** | 具名＋量化钉已交付；**3/10 尺修复后树未复量**（原探针时代不匹配，重放 build failed） |
| AC#5 门禁（含 CLI 两形） | **PASS** | CLI 两形本程补测 PASS（移植探针，183 引读日志 md5 未动）；gate-clauses.sh 非仓内脚本＋本程复跑 harness 阻断（126），名册有效读数仍 c1 |
| AC#6 契约轴 | **PASS** | 冻结文字内容锚等值；五枚 commit 写面零越界；行号基线过期提醒沿袭 |
| AC#7 同名不同目录守卫（牙活） | **PASS** | 无——basename 变异独红 1 枚，全包对照仍独红 |

**closeable？＝YES（可结票）。**

- **可安全翻的框：AC#1、AC#2、AC#3、AC#5、AC#6、AC#7（6 枚），无附条件。**
- **AC#4：带注翻或持开，二选一**——①翻，框内注明「3/10 复量欠：原探针时代不匹配（缺 `declaredPath`/`markObs`），修复后树语料数字不存在」；②持开，把「3/10 尺重写探针复量」立为纯测量小票（零产码）。本终裁推荐 **①**（实质已满足，欠的是一个可后补的语料数字；持开的代价是整票为一句注释挂着）。
- 票 185 的本体缺陷（每一次成功续读造出下一发的阻断者）已被 A381 第四支修掉、判据落地且有牙（红前 4+8、变异独红、CLI 两形 clean）；撤销口令沿袭 A381 的「**撤分页**」。

## ⑪ 本程没做的（逐名）

1. gate-clauses 名册在本 harness 的复推（rc=126 伪象）——未强解，c1 读数有效。
2. 3/10 语料尺的修复后复量——AC#4 的欠格，未重写探针（重写属新测量腿，非终裁范围）。
3. `./cmd/wisp` 全包测试——不在派单四门；纯态编译缺口已由票 98 腿登记归 W123g。
4. frontend／design／台账（`pending-and-issues.md`）——零触碰。
5. 票面 7 个框一枚未勾、票面状态行未改——翻框归编排者。

## ⑫ 机读读数索引

- 门禁：`.scratch/wisp/probes/185/v1/logs/`（g1-build／g2-full／g3-scoped／g3-wide／g4-d22／gates-summary；p6-risk-alone／p7-tools-alone）
- AC#1：p1-red-risk（4 红）／p2-red-tools（8 红）／p5-chain／p9-named-tools；replica `mut-unfix/`＋`ov-unfix.json`（diff 1 行）
- AC#2：p3-claima-tools（8 红）／p8r-trio-tools（3 承重 PASS）；replica `mut-claima/`＋`ov-claima.json`（diff 3 行）
- AC#7：p4-basename-risk（**作废 no-op**）／p4r 重跑（独红 1）／p4b-basename-fullrisk（全包仍独红 1）；replica `mut-basename/`＋`ov-basename.json`（diff 16 行）
- AC#5 CLI：e2e-withprefix／e2e-noprefix（rc=0 两形 PASS）＋e2e-readings-withprefix／-noprefix（审计行逐字）；探针副本 `zz185v1_e2e_test.go`＋`ov-e2e.json`；备份 `e2e-readings.BEFORE-20261106.txt`（md5 c2ac391f…）
- AC#4：p10-freq-replay（build failed，时代不匹配留档）
- gate-clauses：gateclauses-fresh／fresh2（腿数断言过、10 腿 git-grep-rc-126 harness 伪象留档）
- 纯态克隆：`/tmp/wisp185v1-export-8879`（build rc=0；risk ok；tools 唯一红＝票252 环境形；185 族 10/10＋8/8）
