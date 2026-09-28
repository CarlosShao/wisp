# 派单：`185-r1`（写码腿）——分页续读：豁免改绑"宿主亲手落盘过的那一枚文件"（owner 09-28 16:5x 点名批准）

- 派单时刻：`2026-09-28 16:5x`｜编排者锚点 HEAD＝**起手自取**（`git log -1 --format='%h'`，派单时盘上是 `8ca2291f`；⚠ 别照抄，漂了就登记）
- 工单：`.scratch/wisp/issues/185-every-successful-reread-mints-the-blocker-for-the-next-one-…md`（票 185）
- 性质：**写码腿（含产码改动）**｜⚠ **AC 框一枚不许勾**
- ⚠ **这一程动的是安全匹配侧（`internal/risk`）⇒ 授权已落账：owner 原话「**那超长输出，我批了，要做的**」，逐字在台账 `A381`。它批的形状＝下面 §0 那枚"第四支"，⚠ **不是**"放宽 R4"、**不是**"参数侧按值放行"。判"必须动那两支之一才修得好"＝**停手上报**。

## 0. 要修什么（一句话）＋已批形状

现状：模型照宿主指针把一枚超长产物**整份读回一次**已经能成（票 183 已结案）；**第二次再读同一条路径被自家 R4 拒**——阻断者是"第一次读回来的那 20000 字节自己盖的那枚**没有声明的** C25 mark"（票 185 本体）。
`185-c1` 把四支候选逐支现量（表 `docs/evidence/s1/185-…-c1.md`，台账 `A370`）：**(a) 桥侧给 `fs.read` 也声明**＝声明值就是模型递进来的那个字符串＝票 183 AC#5① 明禁的"参数侧按值放行"，且它的变异下 `./internal/tools/` 全包 `ok 17.264s`、`--- FAIL` **0**＝今天落地零牙；**(b) 走再落盘层不产生第二枚 mark**＝要拆票 175"成功结果必须有来源标记"那枚承重声明；**(c) 不修、写成"同一路径只许读一次"**＝这句话今天不真（频率尺 3/10、真实语料零读数）；⇒ **只剩第四支＝本程要落的**：**豁免不再问"模型报了哪条路径"，改问"这条路径是不是本 scope 里宿主亲手落盘过的那一枚"**（宿主确实登记过：`internal/tools/task_backfill.go:112` 那一名册）。
⚠ **这一支的本体是 per-scope 路径名册＝票 177 当时裁"未批"的那一形**（`internal/risk/taintmatch.go` 逐字写着 "a per-scope path roster is explicitly NOT approved"，出处 `A353`）⇒ **owner 这句批准就是翻那一枚**，撤销口令「撤分页」。**你落地时要把那句注释改掉并具名指向 `A381`**（那是文档级更正，不是放宽）。

## 1. 写面（只这些，逐枚具名）

✅ `internal/risk/`（名册类型＋匹配侧查询；⚠ `provenance.go:468-473` 那句 `Mark` 契约文字与 `taintmatch.go:11-15` 那句"逐 token 追踪已被否决"**保持不动**，基线 md5 `858e45116383caa3e7c1dd4b0924fad1`／`5680ddd18e2d2ec2a85e485b54f4c12e`；`taintmatch.go:100-109` 是**仅参考**基线、允许随改动变）｜✅ `internal/tools/`（把名册喂给盖章那一步；`bridge.go` 的 `mark(...)` 附近）｜✅ 新增跟踪判据件（命名带 `185`）｜✅ 新建 `docs/evidence/s1/185-paged-reread-r1.md`｜✅ 新建 `.scratch/wisp/probes/185/r1/**`｜✅ 追加票 185 一段 Progress log（只追加、不改原句）。
⛔ `docs/PLAN.md`（含 `:1531`／`:1532` 两行 DEFERRED）、`docs/specs/**`、`go.mod`／`go.sum`、`allowlist.txt`、`thresholds.go`／golden／审批超时常量 **一字不动**；⛔ **票 183 的文档块原 25 行**（`provenance.go:482-506`，md5 `f89e891e5eee3f3ea2b4f89d921072c4`）不许改写，要更正只许追加。
⛔ `internal/panel/**`、`cmd/wisp/**`（⚠ 另一枚写腿 `181-r2` 正在改 `cmd/wisp/panel_pump.go` 与 `internal/panel/git.go` ⇒ 撞车即停手上报，别硬合）；⛔ `frontend/**`／`design/**` 不读不写不引不转述。
⛔ 别家脏件（`.gitignore`、`probes/152/**`、`probes/161/r6/logs/flip-*`、`docs/evidence/s1/152-*.md`、`design/**` 那批未提交删除）**不动不提交不评论不还原**。
⛔ 零删除命令；不许跑 `probes/161/r6/flip-declaration.sh`；不许改 `scripts/d22scan.sh`／`probes/154/gate-clauses.sh`；只 commit **不 push**；显式 pathspec；禁 `add -A`／`.`／`-a`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`restore`／`clean`／`worktree`／`switch`／`merge`。
⚠ 写面闸门：起手把 `git status --porcelain -- internal/ cmd/` 名册原样贴进表 §1；**终态必须等于起手名册（逐枚具名差集为空）**，不是"必须为空"。
⚠ 变异载具自证（`33-r1` 用超顶 4 枚换来的）：**先证明变异真落到了那个字节**（比对前后哈希，**两枚都必须非空**），再谈"红了"；`grep` 未勾框要写 `'^- \[ \]'`（`[ ]` 是字符类）。

## 2. 必答六格（每格都要答"这一发在**未修码**上响不响"；先测→再写→再提交）

**AC#1 成对判据（缺一形＝把墙拆掉或把豁免做成装饰）**：① 同一任务里 `task.output`→`fs.read` 成功→对**同一条路径**再发一发 ⇒ **不命中 R4**；② **外来正文携带同一条路径**、模型随后去读它 ⇒ **仍然命中**（票 177 W-2 原句）。⚠ 判据必须长在**CLI 接缝形状**上（含回填产生的那条 `ArtifactPath`），不许只是桥级复用——票 164／183 两次都是"包级全绿≠端到端通"。
**AC#2 名册的边界必须被钉，不是被声明**：逐枚答"名册活多久、跨不跨 scope、跨不跨任务"——⚠ 现量起点（我本轮复跑）：`Scope.Close` 删 mark 那一路今天就把 mark 带走了，名册**必须同样随 scope 失效**；补一发"另一枚任务借不到这一枚的名册"（对照 `183-v2` 的 L3 归因形状）。
**AC#3 覆盖面要量，不要宣布**：逐枚答"哪些工具回执里会带宿主路径、它们各自会不会造出同款阻断者"（尺要现跑：`grep -rn "MarkWithHostPath\|hostPathBox" --include=*.go internal/ cmd/`），**名册与枚数不许从票面抄**。
**AC#4 第三形不许坏**：同目录兄弟产物／改一个 rune 的近邻路径／**同名不同目录**（票 185 AC#7，我复跑过 basename 放宽变异 ⇒ `internal/risk` 已入库那批一枚都没红＝今天零牙，**必须新写**）——⚠ 这三形今天应当是**绿的**，它们是守卫、不许充当本票新增的牙。
**AC#5 代价面与收益同批发（第 46 条）**：复量第 4 把尺那一面（真实正文命中率 **3/10** 那把尺，语料组成要具名：`183-v2` 已裁出它是"8 枚仓内文件＋2 枚作者自造夹具"，**不是用户文档命中率**），并逐句答"这次往前推一步，让哪一类外来正文不再算续读证据"。**只报收益不报退让＝本票判不通过。**
**AC#6 真机那一发要到手**：CLI 端到端跑两发续读（第一发＋第二发），逐字贴 `kind=`／`rules_hit=`／`bytes=` 三行；⚠ 带在册 PATH 前缀 `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`，**两形都贴**（不带＝`0xc0000135`＝票 98）；⚠ 别拿包级绿充当端到端绿。

## 3. 门禁（终态读数取在你自己最后一枚 commit 之后）

`go test -count=1 ./internal/risk/`（**必须单包**，四包并发假红＝`A359`）＋`./internal/tools/`＋`./internal/agent/`＋`./internal/panel/`（⚠ 该包今天有 **3 枚在册红**：`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourWayAgree`，照实记、不许当绿、不许修）；⚠ **另加一条我自己欠的**（`A380` 那格账）：**收尾必跑一次全仓不变式面**——`sh scripts/d22scan.sh`（基线 `ban #8 internal/`＝**438**，你新增件会涨，具名登记）＋`bash .scratch/wisp/probes/154/gate-clauses.sh` **比红腿名册不比退码**（在册只 `G6neg`）＋`gofumpt -l`（在 `"$(go env GOPATH)/bin/gofumpt"`）；`git status --porcelain -- internal/ cmd/` 终态等于起手名册。

## 4. 表骨架（第 ≤5 枚调用内落盘并 commit 第一枚）

① 起手锚＋起手脏件名册＋授权出处（引 `A381` 那两句原话） ② 未修码红：AC#1 两发＋AC#4 三形逐字读数 ③ 落地形状：名册存在哪、谁写、谁查、活多久 ④ AC#2 跨 scope／跨任务那两发 ⑤ AC#3 覆盖面逐枚名册（现跑尺） ⑥ **AC#5 代价面**（哪一类外来正文不再算证据＋3/10 那把尺的复量与语料组成） ⑦ AC#6 真机两发逐字 ⑧ 变异自证（含"两枚哈希非空"） ⑨ 门禁全套 ⑩ 本程没测的（逐名） ⑪ 被拒调用＋零删除自证＋工具调用终值 ⑫ 对编排者的不服（若有）

## 5. 硬顶与计数

工具调用**硬顶 40**（这一程要动产码＋端到端，我给足，但**每答完一格自报累计**）；**第 30 枚起不许开新探索**；判"必须动 R4 本身／必须把 `fs.read` 摘出名册／必须设 `PassThroughUnclassifiedRisk`"＝**停手上报**（票 183 AC#5 那三支毒修法一枚都不许用）。

## 6. 结束消息回我七项

① 名册的存储与生命周期（哪一枚结构、谁填、什么时候失效）＋AC#2 两发读数；② AC#1 成对两发的**未修码红**与**落地后绿**逐字；③ AC#3 覆盖面逐枚（带尺）；④ AC#4 三形今天是否仍绿（它们坏了就是回归）；⑤ **AC#5 代价面**（一句"让哪类外来正文不再算证据"＋3/10 复量）；⑥ AC#6 真机两发逐字（含 PATH 两形）；⑦ 门禁全套＋起手/终态名册差集＋`taintmatch.go` 那句"NOT approved"注释你改成什么了（逐字）＋被拒调用＋零删除自证＋工具调用终值＋没测的逐名。
