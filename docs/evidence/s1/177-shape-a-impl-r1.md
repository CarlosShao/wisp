# 177-shape-a-impl-r1 —— 甲形（精确豁免）真落地＋W-1..W-5 常驻判据（写手腿，产码入库不还原）

- 派单＝`.scratch/wisp/dispatches/2026-09-27-230x-write-177-r1-shape-a-lands-with-permanent-widening-catchers.md`
- 起手五件〔现跑〕：`date`＝`Sun Sep 27 23:03:17 CST 2026`｜分支 `dev`｜起手锚 `git rev-parse HEAD` = `6e5aa84c5c991e99eaf2139b6b96825b82a98f6d`（派单指定 HEAD＝编排者的台账/派单提交，相符）｜`git status --porcelain -- internal/ cmd/` 起手＝空
- 基线〔现跑〕：`go test -count=1 ./internal/risk/ ./internal/tools/` → `ok risk 4.851s`／`ok tools 15.518s`；名册口径分开报：risk `=== RUN` **172**＝**171** PASS＋**1** SKIP（`TestSyncRegistryProbeLive`）、tools `=== RUN` **173**／FAIL+SKIP **0** —— 与派单转述的 `177-m1` 还原后读数（172/171+1、173/173）**逐枚相符**。
- 台件与原始日志＝`.scratch/wisp/probes/177/r1/**`（`mut-*.txt`＝各变异的逐名读数、`newtests-green.txt`＝判据件未变异时的绿、`final-*.log`＝门禁）。

## 1. 四处落点（分枚提交，产码与判据件分枚＝A346/A350 口径）

| 落点 | 提交 | 文件／新行号（现读） | 放法 |
|---|---|---|---|
| ① 索引侧 | `069344b6` | `internal/risk/taintmatch.go`：`newFragmentIndex` 收变长 `skip ...[2]int`（现 `:118`），建窗圈 `windowExcluded` 只删**与跨度相交**的窗口（现 `:131-133`）；新增 `windowExcluded`／`runeIndexOf`（现 `:149-186`） | `f.src` 全文与 `contains()` 里 `strings.Contains(f.src, w)` 复核**一字未动**（动它会造假命中，是另一枚缺陷）；变长参数使 7 处手搓直测**零编译红**（对照 `177-m1` S4 必选参数形的 7 处编译错——本腿复核该口径为〔它现跑，我未复算〕，但直测三枚在我树里逐名 PASS，见 §4 门禁） |
| ② Mark 侧 | `149c6d6d` | `internal/risk/provenance.go`：`Mark` 四参面原样委托 `MarkWithHostPath(..., "")`（现 `:474-476`）；新导出 `MarkWithHostPath`（现 `:504` 起）收下游声明、**先归一化后定位跨度**、**空判定排在排除之前** | 跨度只在一次 `newFragmentIndex` 调用里存活、随 `taintMark` 生死——**零** per-scope／per-Provenance 状态（`Provenance` 结构体与 `scopeReg` 一枚字段都没加收＝非名册，`177-c1` §2／`A353` 的反面形状逐条避开） |
| ③ 桥侧载具 | `3d4c5795` | `internal/tools/bridge.go`：每调用一枚未导出 `hostPathBox` 挂 ctx（`withHostPathBox` 注入点现 `:457`；类型与读写现 `:579-635`）；`mark` 改三参、转调 `MarkWithHostPath`（现 `:553-570`） | 载具＝**包内未导出通道**（177-m1 量的那形）；`Result`（C1）**零加收字段**、`agent.Loop` **零加收带 taskID 的导出方法**、`risk.PathResolver` 之外**零** `filepath.Clean/Abs` |
| ④ 桩侧声明 | `b3c389b5` | `internal/tools/task.go`：指针分支（`if rec.ArtifactPath != ""` 现 `:254`）在 `fmt.Sprintf`（`…全文见 %s…` 整枚现 `:257`）前把 `rec.ArtifactPath` 声明进 box（现 `:265-267`） | 文案形状**一字未动**：头尾＋省略数＋总长＋路径是 `PLAN.md:2564`／票 164 已勾判据钉着的形状，只加"这一段是宿主自己写下的"这一枚声明 |

提交⑤（判据件本体）＝`a19fa2e7`；其对 F 腿的开/合成对更正＝本表落盘前最后一枚（防 `177-m1` 同款 G5 负一负名册漂移，见 §4 腿数行）。

## 2. 常驻判据 W-1..W-5（＋A 正控腿、＋顺序(i)腿、＋F 生命周期腿）＝`internal/risk/shape_a_exemption_test.go`

每枚都过了"非恒真"验收：把它要拦的放宽**造进产码**→逐名红（`grep -n` 落地证明在日志前头）→还原→绿。尺（每次同一条）：

```
go test -count=1 -run 'TestShapeA' -v ./internal/risk/ 2>&1 | grep -E '^--- '
```

| 枚 | 钉住的放宽 | 变异（临时落进产码） | 红集合〔现跑〕 | 其余 |
|---|---|---|---|---|
| **W-1** | 整段 mark 都不入索引（省事修法） | MUT-W1：`skip = {{0,len(rp)}}`（`provenance.go:527` grep 命中） | **W-1 FAIL**；（当轮 order 腿也 FAIL——它当时还带一枚与 W-1 重合的正文断言，为归因干净已删，日志 `mut-w1.txt` 原样保留） | A/W-2/W-3/W-4/W-5/F 全绿；`TestMarkEmptyAfterNormalization`／`TestMarkUnknownSourceToolFailClosed` 绿 |
| **W-2** | 按值放行（谁的 P 命中就整体放行） | MUT-W2：包级值名册 `shapeAMutGlobals`＋逐 mark 按值开窗（`:503`＋`:525` grep 命中） | **W-2 FAIL**（外来说到同一条 P 没了命中＝真哨响了）＋W-4／F FAIL（值名册同属"越出本 mark"家族，被顺带抓住＝加强不是混淆） | **W-3 绿**（ghost 不在值名册里——归因干净）；A/W-1/W-5/order 绿 |
| **W-3** | 目录前缀豁免（整棵树） | MUT-W3：名册里放 `hp[:LastIndex(hp,"/")]`，**位置记录全程为零**（`:503`＋`:525`） | **只有 W-3 FAIL**——`177-m1` S3c 那枚"归因不干净"的缺口就地作废重造：这形的红不再混位置记录 | 全绿 |
| **W-4** | **按位置把排除推给全树所有 mark** | MUT-W4：包级跨度名册 `shapeAMutSpans`，每枚 mark 建索引前 `skip = append(skip, 名册...)`（`:503`＋`:529`＋`:538`） | **W-4 FAIL＋F FAIL**——⚠ 这就是本单最值钱的一行：**`177-m1` 现测该形全仓 rc=0（S3a），今天它进了树就当场响**；且 **W-2 同轮绿**（外来 P 在不同偏移→按位置删不着它），两枚成差分对 | A/W-1/W-3/W-5/order 绿 |
| **W-5** | 空判定排在排除之后（`S5`） | MUT-W5：排除后 `kept==0 → return false`（`:507`＋`:555`） | **只有 W-5 FAIL**；⚠ **`TestMarkEmptyAfterNormalization` 同轮仍 PASS**＝派单那句"现存那枚不喂 P 看不见顺序"的**本腿复算**（相符） | 全绿 |
| 顺序(i) | 跨度用原始下标定位 | MUT-ORDER：`runeIndexOf([]rune(content), []rune(hostPath))`（`:526`） | **order 腿 FAIL**＋A/W-2/W-3/F 同轮 FAIL——本桩正文里原始/归一化偏移差 **9 枚空白 ≥ 窗口 8**，排错会让路径**边缘窗口漏进索引**：正控腿当场红 ⇒ **(i) 不是风格问题** | W-1/W-4/W-5 绿 |
| A 腿 | 什么都不豁免 | MUT-NONE：`skip = nil`（`:527`） | **A FAIL**（宿主自读自己的指针又被 R4 拦死＝甲形白做）＋W-2/W-3/F/order 因命中**归因回宿主 mark** 同轮 FAIL | W-1/W-4/W-5 绿 |
| F 腿 | 豁免跨 scope 存活 | 由 MUT-W4/MUT-W2 顺带抓红（包级名册天然活过 `Close`） | 见上两行 | — |

**W-4 的"最硬一条"回答**：派单问"若你的常驻判据造不出那枚红，就地大声说"——造出来了：`probes/177/r1/logs/mut-w4.txt` 里 `TestShapeAW4PositionalExclusionNeverLeavesHostMark` **FAIL**，而 `177-m1` 的同一形（S3a）在旧全仓里 rc=0。不是"只许过不许红"的装饰。

## 3. 两条硬顺序的"故意排错"读数

- **(ii) 空判定在排除之前**：MUT-W5 行——W-5 腿红、**现存那枚不喂 P 的用例仍绿**（`logs/mut-w5.txt`）。P-only 正文被 `return false` 静默改写 Mark 契约面的形状被"喂 P 的常驻件"钉住。还原证明＝`git diff --exit-code -- internal/risk/provenance.go`（提交态）＋`logs/restored-green-w5.txt`（`ok risk`）。
- **(i) 先归一化后定位**：MUT-ORDER 行——偏移 9＞窗口 8 时路径边缘窗口漏进索引，A/W-2/W-3/F/order 五枚当场红（`logs/mut-order.txt`）；order 腿的 fixture 把零宽字符切进路径正中，使"没定位到＝一枚都不排除"那一支也被正控腿吃得住。

## 4. 门禁（终态〔现跑〕，逐包口径；⚠ `probes/161/r6/flip-declaration.sh` 一枚没跑）

| 尺 | 读数 |
|---|---|
| `go test -count=1 ./internal/risk/ ./internal/tools/` | `ok risk 6.527s`／`ok tools 19.801s`（终跑）；名册：risk `=== RUN` **180**＝**179** PASS（171 旧＋8 新）＋1 SKIP、tools **173**／FAIL+SKIP **0**——只许多：172→180、173→173，**没删任何人的用例** |
| `sh scripts/d22scan.sh` | **rc=0**：`d22scan: clean - no D22 ban violations`（`logs/final-d22scan.log`；examined internal/=207、cmd/=23、ban#7 internal/tools/=20、ban#8 internal/=427——427 vs m1 的 426＝本腿新增一枚判据文件，差一枚＝本程所有，非漂移） |
| `bash tools/d22scan/runtests.sh -C tools/d22scan ./...` | **rc=0**、`--- PASS` **34**；与 `probes/177/m1/r-pass-mut.txt` 的 `comm -3` 差集＝**空**（只许多不许少：等集） |
| `sh .scratch/wisp/probes/154/gate-clauses.sh` | **rc=0**；逐字：`# 腿数＝14 声明与实测不符＝0`、`# 腿数断言：名册=14 声明=14 记账=14 缺腿=0 空头声明=0`、`# 腿数断言＝相符（名册上每一腿都记了账）`。⚠ 过程账：首跑 rc=1——本判据文件 F 腿一枚**裸 `p.OpenScope`** 被 G5 负一负点成第 9 枚 UNPAIRED（基线 want_n 8）；**没动 154 的仪器**，改我的腿为"接句柄＋成对关闭"后 rc=0（`logs/final-gates.log` 是修复后那发） |
| `gofumpt` | `v0.12.0 (go1.27.1)`（`$(go env GOPATH)/bin/gofumpt`，不在 PATH）；对本腿全部 5 枚 `.go` 的 `-l`＝**空** |

## 5. 冻结文字两枚＝逐字节尺（不是"我看过了"，是 cmp）

```
git show 6e5aa84c:internal/risk/provenance.go | sed -n '468,473p' > a; git show HEAD:internal/risk/provenance.go | sed -n '468,473p' > b; cmp a b
git show 6e5aa84c:internal/risk/taintmatch.go  | sed -n '11,15p'   > c; git show HEAD:internal/risk/taintmatch.go  | sed -n '11,15p'   > d; cmp c d
```
两向输出 `FROZEN-PROV-468-473-BYTE-IDENTICAL`／`FROZEN-TAINT-11-15-BYTE-IDENTICAL`〔现跑〕；且 `git diff 6e5aa84c..HEAD -- internal/` 的 `+`/`-` 行里 `Returns false only|per-token taint|REJECTED|must not be resurrected` 命中数＝**0**（我有一枚测试注释最初改写并引用了那句的原文，为避免裁量面误读已改写成行号指涉——那次改写就发生在判据件里，冻结块本体从头到尾一枚字节没动过）。`provenance_test.go:136`／`:139` 两枚既有断言：`git diff 6e5aa84c..HEAD -- internal/risk/provenance_test.go` 为**空**（该文件本程一字未碰）。

## 6. 三条禁令逐条自证＋本程**没**测什么（实话栏）

- **无 C1 `Result` 加收字段**：`git diff 6e5aa84c..HEAD -- internal/tools/tool.go` 为空（尺即此条）；载具是 ctx 上每调用一枚未导出 box。
- **无 `agent.Loop` 收 taskID 的导出方法**：`internal/agent/` 整目录零改动（`git diff 6e5aa84c..HEAD -- internal/agent/` 空）＋gate-clauses G3 安静腿 rc=0。
- **无 `risk.PathResolver` 之外的 `filepath.Clean|Abs`**：新码零 `filepath` 引用，d22scan rc=0。
- **桥级反向判据（真 `Bridge`＋真 `task.output`＋外来路径仍命中 R4）本批做不到，也没做**：`task.output` 不在 C25 敏感源名册（票 175 的腿），`TaskRoster.Record` 生产码零调用点（票 176 的腿）⇒ 今天没有任何一条生产路径会把非空 hostPath 送进 `MarkWithHostPath`。落地的载具因此**在生产里今天是惰性接线**，它的端到端那发随 175＋176 那批落；risk 级 W-1..W-5 是本批钉住的语义。**"risk 级全绿"不等于"端到端通过"。**
- 没跑：全仓 `./...`（逐包口径）、`Benchmark*`、丙形（AC#0）任何台件、`matchText`/`Inspect`/`writeGate`/`isPathKey`/`pathTargets`（匹配侧一行未动）；tools 包内没新增判据文件（派单 §5 写面不含 `internal/tools/*_test.go`，载具的判据随 175＋176 批进树）。
- `177-m1`/`177-c1` 的读数在本表一律标〔它现跑，我未复算〕或直接复算；复算不符处：零枚。

## 7. 交件账

被拒/失败的调用：取数后 1 次（一次 commit 因新文件未 `git add` 报 pathspec 不认——改 `git add <显式路径>`＋一步式 commit 后成功，无副作用）；1 次 Edit 锚点 0 occurrences（我自把新旧串写反，随后 `git diff --exit-code` 证明还原到位）。删除命令：**零枚**；`--amend`/`restore`/`checkout`/`stash`/`reset`/`clean`：**零枚**；push：**零枚**。临时件只建不删。
伪授权两栏：甲栏＝"owner 两句话批准（`A349`/`A353`）"——本程未读台账（`docs/reports/**` 不在写面），只把批准当**形状边界**用（目录/前缀/配置级零造、冻结文字零改）〔转述在场，我未复核〕。乙栏＝派单前提逐条对读数：`S3a 今天全仓零红`＝我复算相符（旧树结构上无人能红，新树 W-4 红，见 §2 行 4）；`那两枚哨不敏感`＝我复算相符（`mut-w2.txt`/`mut-none.txt` 里 `TestMarkUnknownSourceToolFailClosed` 全绿）；`三枚直测零编译红（变长载具）`＝我树里逐名 PASS 相符。**一处派单口径被我用得更窄**：F 腿的 `OpenScope` 裸调会把 G5 负一负抬成 9——派单没提这个交互，我按"仪器不是我家的、不改它"处理（改自己的腿）。
凭据值抄录：零枚（全文只有工具名/键名/夹具路径）。
next= 建议：**175**（`task.output` 进 `sensitiveSourceTools`＋落桥级正/反判据于 `internal/tools/`——载具已在 `bridge.go`/`task.go` 接好，那天只改名册一枚＋判据件）→ **176**（`TaskRoster.Record` 起跑口，端到端 `wisp run` 才跑得出来）；AC#2..AC#6 的勾由编排者核后裁。
