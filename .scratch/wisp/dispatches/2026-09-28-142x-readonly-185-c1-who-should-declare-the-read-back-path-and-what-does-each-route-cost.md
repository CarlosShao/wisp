# 派单：`185-c1`（只读普查·零产码）——"照宿主指针读回来的那一发，第二次为什么还被拒"这枚**声明该由谁给**，三支候选逐支给现量与代价

- 派单时刻：`2026-09-28 14:2x`（`date` 现量 14:15 之后）
- 编排者锚点：HEAD＝**`78ea1d19`**（票 183 已结线 `-done`，十格全勾）
- 工单：`.scratch/wisp/issues/185-every-successful-reread-mints-the-blocker-for-the-next-one-because-the-fs-read-result-becomes-an-undeclared-c25-mark.md`
- 性质：**只读普查**。⇒ **AC 框一枚不许勾**（勾由编排者翻），**产码零字节不许动**。

---

## 0. 先把你最可能踩的那枚坑写在最前面：**"未修码"就是今天的 HEAD**

上一枚同族派单（`183-r2`）里我写了一句"这一枚今天就能红"，**那句话错了**，错因是**我把"未修码"当成了 HEAD**：票 183 的修法早在 `0662a35a` 就落地，HEAD 已经是"修完之后"。⇒ 那枚判据在 HEAD 上是**绿**的，而格面要的"未修码红／落地后绿"两向我只能**自加一把尺**（把 `internal/risk/{provenance.go,taintmatch.go}` 钉回 `b277e1e1` 走 `-overlay`）才取回来。

**本单的更正定式**：票 185 的缺陷**今天仍然存在**，所以——

> **本票的"未修码"＝你起手那一枚 HEAD（`78ea1d19` 或它之后的继任）。凡你说"这枚判据今天应该红"，请同时具名你拿的是哪一枚 commit 当基线；不许把"HEAD"和"未修码"当同义词用而不检查。**

⚠ 反向也一样：**不许拿"今天已经绿"当"这件事做完了"**——票 183 的教训是判据可以全绿而承诺没兑现。本票的兑现判据是**同一条路径的第二次读回**，那一发今天在真机上是红的（读数见 §2）。

## 1. 现象（编排者本程在 HEAD 上亲取，不是我转述）

尺：`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test -count=1 -v -overlay=.scratch/wisp/probes/183/r1/overlay-cli.json -run 'TestRereadHostPointerOnCLISeam183a1' ./cmd/wisp/`
读数：`.scratch/wisp/probes/183/orch/logs/cli-with-prefix.txt`（60.10s，`--- FAIL` 红在 `zz183r1_e2e_test.go:373 run B（模型续读）退出 1`）

三行审计逐字（同一份文件里）：

```
kind=success … tool=task.output risk=L0 decision=allow rules_hit=[]
kind=success … tool=fs.read     risk=L0 decision=allow rules_hit=[]      ← 第一次照指针续读：通了（票 183 的功劳）
kind=refused … tool=fs.read     risk=L2 decision=reject rules_hit=[R4] reason="R4: 包含来自 fs.read C:\Users\s…"  ← 第二次：被自己读回来的内容咬住
```

⇒ **阻断者的来源那一枚换成了 `fs.read` 自己**：第一次读回把那 20000 字节盖成了一枚**没有声明的** C25 mark，第二次续读的证据就是这枚 mark。按票 183 AC#3／票 177 W-2 的边界，**没声明的正文命中是应当的**——所以这不是回归，是**下一次覆盖**（`183-v2` 判、编排者采纳，故立本票）。

不带 PATH 前缀那一形（同一枚台件）＝`exit status 0xc0000135`＝票 98 那枚加载期坑，**两形都要贴**、不许只贴绿的那形。

## 2. 本单只答三格：AC#2（主菜）＋ AC#3（分母）＋ AC#7（那一形的现状）

### AC#2 —— "第二次读回该由谁给这枚声明"，三支候选**逐支给现量与代价**，不许凭偏好选

**(a) 桥侧**：`fs.read` 成功时也填 `hostPathBox`（它读的就是那条路径，宿主当场知道）。
必答：**现量取号**（`grep -n "hostPathBox\|withHostPathBox" internal/tools/bridge.go internal/tools/task.go`，号会漂、以你现跑为准）；这一支要落在哪一枚函数、那一枚函数今天有没有别的调用方；**票 177 已裁过的三条边界里它会碰到哪条**（per-scope 名册／参数侧按值放行／`allowlist`）——⚠ **若你判它等价于"按参数值放行"，那正是票 183 AC#5①明禁的形状，具名停手上报，不要自己放行**。

**(b) D15 再落盘层**：续读走"宿主替模型读"的通道，根本不产生第二枚 mark。
必答：那一层今天在哪（`grep -n "func " internal/agent/spill.go | head -30`）、"不产生 mark"这句话在码上等价于改哪一处、改了之后**哪些既有判据会红**（现跑：`grep -rln "C25\|Mark(" internal/agent/*_test.go internal/risk/*_test.go internal/tools/*_test.go`，逐枚点名、别只报枚数）、以及它会不会把票 175 的"成功结果必须有来源标记"那枚承重声明拆掉（**那一条是禁令面**）。

**(c) 不修**，把"同一条路径只许读一次"写成明说的事实与文案。
必答：**文案落在哪一枚文件、那句话是不是冻结件**（⚠ `task.output` 的桩文案由 `PLAN.md:2564` 冻结，尺：`grep -n "全文见" internal/tools/task.go docs/PLAN.md`）⇒ 若落点在冻结件里，这一支**要人工批准**，请具名报回、不要自己改一字。

**再补一枚第四支**（如果码上真有第五条路就报，没有就明说"没有"）：有没有哪一支能**只作用于这一发**（同一条产物路径的续读）而不放宽任何一类外来正文？判据形状是什么？

### AC#3 —— 覆盖面**要量、不要宣布**

尺（现跑，名册与枚数**不许从票面抄**）：`grep -rn "MarkWithHostPath\|hostPathBox" --include=*.go internal/ cmd/ | grep -v _test.go`
必答：逐枚答"哪些工具的回执里会带宿主路径、它们各自会不会造出同款阻断者"（已知候选：`fs.read`／`task.output`／`clip.*`／落盘回执），**每一枚带一条 grep 或一条读数**；今天带声明的只有哪几枚、不带声明的有哪几枚。

### AC#7 —— "同名不同目录"那一形（票 183 AC#4 归口过来的）现状

必答一把尺：今天的豁免会不会作用到"文件名相同但目录不同"的那条路径？（读 `spellsDeclaredPath` 与 `attachDeclaredPath` 的比对面，尺：`grep -n "spellsDeclaredPath\|attachDeclaredPath" internal/risk/*.go | grep -v _test.go`）
⚠ **不许拿"看着到不了"当凭据**——那正是票 183 整张票的教训。要么造一枚 `-overlay` 变异（把豁免改成只比 basename，**别动跟踪件**）看它红，要么明确写"本程未验、这一形仍无人钉"。

## 3. 写面（只这些，别的都算越界）

- ✅ 新建：`docs/evidence/s1/185-reread-owner-census-c1.md`（裁决表，**骨架先落盘再逐格填**）
- ✅ 新建：`.scratch/wisp/probes/185/c1/**`（读数／台件副本／overlay 名册）
- ✅ 追加：票 185 的 Progress log 一段（**只追加、不改原句、不勾任何 AC 框**）
- ⛔ **零产码**：`internal/**`／`cmd/**` 一字不动（起手与终态各跑一次 `git status --porcelain -- internal/ cmd/`，都必须是空；不为空就把你要依赖的跟踪件用 `git show HEAD:<path> > <你的探针目录>` 钉进自己的 `-overlay`，并在表里具名）
- ⛔ **不许动冻结面**：`docs/PLAN.md`、`docs/specs/**`、`thresholds.go`／golden／审批超时常量／`allowlist.txt`、`internal/risk/provenance.go:468-473`、`internal/risk/taintmatch.go:11-15`、`provenance.go` 的 `MarkWithHostPath` 文档块原 25 行。三枚受保护文字交件前逐枚复量 md5（基线：`858e4511…`／`5680ddd1…`／`f89e891e…`）；⚠ **`taintmatch.go:100-109` 不在受保护名单里**（那是"仅参考"的一格，会随行号推移而变，别据此报警）
- ⛔ **不许派落地腿**：本程只交读数与代价，**选哪一支归编排者裁**
- ⛔ `frontend/**`／`design/**` 属别家地界：不读、不写、不引、不转述
- ⛔ Git 纪律：只 commit 不 push；commit 必须带**显式 pathspec**；禁 `git add -A`／`git add .`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`restore`／`clean`／`worktree`；**临时件只建不删**（不许跑任何删除命令）；**不许跑 `probes/161/r6/flip-declaration.sh`**
- ⛔ 凭据：任何 API 密钥／DPAPI 明文都不许进读数文件，路径与哈希值可以

## 4. 门禁（只读程也照样取数，但**按 `A363` 薄规矩不充当 AC 结案凭据**，只作现状记录）

- `sh scripts/d22scan.sh`（基线：rc＝0；`ban #8 internal/` **examined=433**——多一枚少一枚都要具名解释，不许默默漂）
- `bash .scratch/wisp/probes/154/gate-clauses.sh` **比红腿名册、不比退码**（今天唯一在册红腿＝`G6neg`＝票 178；名册差集非空就停下报我）
- 逐包 `go test -count=1 ./internal/risk/`（**必须单跑**，四包并发会假红＝`A359`）＋`./internal/tools/`＋`./internal/agent/`
- CLI 那一面走 `-overlay` **且带在册 PATH 前缀**（`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`，出处 `scripts/wisp-cli-tests.sh:20`；不带就是 `0xc0000135`），**两形都要贴**；⚠ 真机那一发单次约 **60 s**，**最多跑两发**（一发足够就别跑第二发），跑之前先确认真的有必要
- 终态读数取在**你自己最后一枚 commit 之后**

## 5. 表要长的骨架（先落这一枚再开始取数，之后增量填）

`185-reread-owner-census-c1.md` 十二节：① 起手锚／写面闸门／三枚受保护文字 md5 起手向 ② AC#2(a) 现量与代价 ③ AC#2(b) 现量与代价 ④ AC#2(c) 现量与代价＋冻结判定 ⑤ 有没有第四支 ⑥ AC#3 覆盖面名册（逐枚带尺） ⑦ AC#7 那一形的现状（验了还是没验，具名） ⑧ **本程没测什么**（逐名，别写"其余都覆盖了"） ⑨ 门禁终态 ⑩ 被拒／没成功的调用（逐条） ⑪ 有没有跑过删除命令＋工具调用终值自报 ⑫ next＝**落地腿派之前还缺什么**（含"要不要摆人工批准项"）

⚠ 恒真判据老规矩：**今天会绿的守卫不许充当新功能的牙**；要报"有牙"必须给变异。

## 6. 硬顶与增量交付（上一枚同职能死于 52／50 枚，别再犯）

- 工具调用硬顶 **40 枚**；**到 28 枚就停止新探索**，只余量用于填表、跑门禁、commit。
- **骨架先落盘**（第 ≤6 枚调用内把 §5 那十二节标题＋待填写进去并 commit 第一枚），之后**每两格一交**，每交一枚 commit。表没写完不许继续取数。
- 交件顺序：表 → 票面 Progress log 追加 → 门禁终态 → commit → 结束消息。

## 7. 结束消息必须回给我的（缺一我就只能退回）

① 三支候选**各自**的一句话结论＋最贵的那枚代价；② **有没有一支会碰到票 177 已裁过的边界**（碰到＝具名，我自己都不敢替你判）；③ AC#3 那份名册（工具→回执带不带宿主路径→会不会造同款阻断者）；④ `internal/**` 的写面自证（起手／终态两条 `git status --porcelain -- internal/ cmd/` 都为空）；⑤ 三枚受保护文字 md5 终态复量；⑥ 门禁四数＋红腿名册；⑦ 被拒的调用逐条；⑧ 跑没跑过删除命令；⑨ 工具调用终值枚数；⑩ next＝落地腿派之前还缺什么。

**要不要摆人工批准项，请只报"是哪一格、为什么"，不要替编排者拍板。**
