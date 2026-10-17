# 票 330 — 门铃仪器（票 165 `AC#6`）的三枚具名缺口：镜像⇄真名册零对接／`Resident` 是装饰／`THE JUDGMENT` 标错了门

- Status: **已立，未派**（2026-10-16 16:3x 编排者立；来源＝非实现者终裁腿 `165-v2` 的 `G1`/`G2`/`G3`，判语件 `docs/evidence/s1/165-v2-verdict.md`＝118 行／13890 字节；台账 `A1009`）
- 关联：票 165 `AC#6`（本格已由 `165-v2` 裁"成立、带三枚具名缺口"⇒ 编排者已翻勾；缺口归本票，⛔ 记在 165 头上）／票 329（同面 `cmd/wisp`，**必须串行**）

## 这是什么（人话）

165 那格新装的仪器能听见"模型点了这枚工具的名"。但它听的是**自己那份抄下来的名册**，不是**真名册**。今天把抄本里的名字整对换掉，仪器一声不响——也就是它没在管"这枚工具到底存不存在"。另外两枚小账：抄本里"这枚工具常驻不常驻"那个开关，拨了仪器也不看；还有一句注释把"判据"标在了一扇本来不管这事的门上。

## 三枚缺口的凭据（逐字取自 `165-v2`，⚠ 引用前先重跑突变）

| 号 | 那一发突变 | 实测 | 缺口 |
|---|---|---|---|
| **G1**（重） | 把用例里的期望名与夹具里的名字**一起**改成 `plan.presentX`（`44f51fc8` 导出树） | **绿、rc=0** | "镜像名册 ⇄ 真名册"今天**零枚用例对接**：`planPresentName165r3` 一枚常量喂四处（`:38/:67/:136/:219`），与夹具同步改名即自洽。机理＝导入环（票 165 `AC#6` 原文硬约束：`internal/agent` 只能包内形状、⛔ import `internal/tools`）⇒ 真名册那一半只在 `cmd/wisp` 够得着 |
| **G2**（中） | 镜像条目 `Resident: true → false` | **绿、rc=0** | 档位开关在这枚用例里**不起可见作用**：`internal/agent/inject.go:44-49` 切驻留／检索两支，`bm25Select`（`:172-224`）`score>0` 即入选，而本用例目录**只有 1 枚候选**、查询"帮我把这个项目收拾一下"与描述共用 unigram `我/一/下` ⇒ 非驻留照样被注入 ⇒ `:219` 看不见档位换了。★按派单口径这一格＝**缺世界**（需要"多候选"那个世界才谈得上有牙），⛔ 写判据成立 |
| **G3**（轻·描述面） | 只把 sink 计数喂成常量 `n := 1`、牙完好 | **绿 rc=0** | `:131` 把 `THE JUDGMENT` 标在对牙不敏感的那扇门上（牙实际在 `:143`：`rowOfTool165r3` 独立遍历行、不复用计数）⇒ ⛔ 恒真，算一枚冗余门＋一句会误导下一位读者的注释（`:78-82` 作者已自陈，故⛔ 记谎、记误导） |

★`165-v2` 同时**收窄了票面上两句**（⛔ 改原文，只登记）：那句"两枚字段都⛔ 红"要说清＝`Declared` **有牙**（M-b2 `L2→L1` 唯一红在 `:157`；M-g `L2→L0` 五枚红 `:157/:184/:191/:199/:208`），`Resident` **装饰**（G2）。

## 现量（编排者 10-16 16:3x 自己跑，行号只当读数）

- 被裁字节基线 sha256（`165-v2` 自报并逐位回基线）：`internal/agent/testdata/golden/plan-present.sse` `7d26fa889b59fc3a…`（28 行）／`internal/agent/plan_present_doorbell_165r3_test.go` `00c32049994c441c…`（221 行）。
- 真名册那枚条目出处＝`internal/tools/plan_present.go:72`（`Name()`）与 `:123-132`（参数面）；装配注册＝`cmd/wisp/run.go`（票 270 的 `270-v1` 现量：生产链上装配点 1 枚）。

## 选形（编排者裁定，⛔ 腿自己改形）

**甲**：在 `cmd/wisp/**` 加一枚常驻判据——**真名册里那枚工具名 ⇄ `plan-present.sse` 线上那枚名字必须逐字相等**；★正控＝把夹具名与镜像名一起改而**真名不动**，那一发**必须红**（这正好是把 G1 那一发从"绿"变成"红"）。
- ⛔ 乙（在 `internal/agent` 里再抄一份"真名"常量）：那只是把镜像挪了个地方，G1 原样存在。
- ⛔ 丙（导入环硬拆：让 `internal/agent` 引 `internal/tools`）＝碰 `bridge.go:13` 那条已钉死的环，属契约级、要人工批准。
- G2 那一格**⛔ 塞进甲里顺手做**：要装牙必须先造"多候选"那个世界（同目录再造一枚**不驻留**的工具条目＋查询与它共享词），属**第二枚落地腿**，落 `internal/agent/**`（与本票写面不同面 ⇒ 可并行，但要⛔ 撞车）。
- G3 只改注释一句（`:131` 那处标成"冗余非零门，牙在 `:143`"），⛔ 动产码字节，随甲同批。

撤销口令＝**「330 改形乙」**。

## 验收

- [ ] **AC#0 交件三把尺**（名册＝`cmd/wisp/**`／`probes/330/r1/**`／`docs/evidence/s1/**` 一枚短名件／本票末尾一条 bullet）。
- [ ] **AC#1 G1 对接钉落地＋正控两发**：正向＝真名与线上名相等 ⇒ 绿；★**G1 那一发重放**＝夹具名与镜像名一起改成 `plan.presentX` 而真名册不动 ⇒ **必须红**（这一发是 `165-v2` 报的绿，本票把它变成红＝本票唯一判据）；当场还原并报还原 sha256 回基线。
- [ ] **AC#2 ⛔ 新名册假对接**：判据必须**从装配出来的真名册读**（⛔ 从第二份硬编码字面量读）；★尺＝`git diff` 里新增的名字字面量枚数，除真名册那一处外必须为 0，并具名写"这枚字符串是从哪一行读出来的"。
- [ ] **AC#3 G3 那句注释说实话**：`internal/agent/plan_present_doorbell_165r3_test.go:131` 附近改成"冗余非零门，牙在 `:143`"；★产码字节零改动那把尺必带（滤整行注释后 `diff` 空、rc=0，**并先跑正控证明这把尺有牙**）。
- [ ] **AC#4 G2 归口具名**：本票⛔ 关掉 G2，必须在票尾一条 bullet 里写"缺世界＝多候选夹具未造"，并留第二枚落地腿的编号位（编排者排程）。
- [ ] **AC#5 台面与颜色**：`cmd/wisp` 整包那格＝★`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp/ -count=1 -timeout 30m`，四数点名、终态与起手红名册作差＝0；⛔ `./...`；起手那一发⛔ 共享工作树 ⇒ 仓外导出树八件套＋`docs`（`paths_shortname_252_*` 那族要仓根 `docs/`）＋手工三件＋`D:` 卷带 8.3 短名。
- [ ] **AC#6 门禁四数只在工作树**：`gofumpt`／`d22scan`／长度门／`go vet ./cmd/wisp/`，逐枚 `rc=`。
- [ ] **AC#7 ⛔ 撞 `329-r1`**：本票**必须排在票 329 之后**（同一枚共享编译面每波只许一枚写腿）；派单里必须具名写"起手锚＝329 落地后的那一枚提交"。
- **330-r2 交件（2026-10-17 09:3x +08·只追加这一条、⛔ 翻任何 AC 框、⛔ 改上面一字）**：G2 那格的"缺世界"**已造出**——新建 `internal/agent/plan_present_crowded_world_330r2_test.go`（272 行，⛔ 产码字节）把真镜像条目（组合 `planPresentDirectory165r3()`、⛔ 第二份抄本）与**一枚**不驻留的长尾候选放进同一目录，检索名额按 `internal/agent/inject.go:232-243` 的 300 token 预算削：长尾那条索引行占 291、镜像那条占 20，⇒ 镜像条目一旦从驻留支拨到检索支就抢不过名额、整枚不上线（档位这时才决定谁被注入）。★成对两发现跑（只打仓外导出树 `D:\tmp\w330pst` 的字节、逐发当场还原核回提交内 sha256）：去牙发＝长尾候选 `Resident: false→true` ⇒ 唯一红在 `plan_present_crowded_world_330r2_test.go:247`（`tier leak: resident arm = [plan.present ledger.sweep330r2], want "ledger.sweep330r2" kept out of it …`）；反控发＝查询换成与长尾无共同字 `"frobnicate widget quietly"` ⇒ 红在另一行 `:254`（`retrieval arm = [], want "ledger.sweep330r2" in it …`）——红句与行号都不同 ⇒ 这枚用例响的是档位、⛔ 恒响。G2 那一发重放（`plan_present_doorbell_165r3_test.go:69` 拨 `false`）⇒ 红在 `:188`（第一轮请求体里那扇门整个不见了）＋`:240`（驻留臂空），而 165 那枚用例本身同发仍 `--- PASS`＝它的盲区属 G1／330-r1 的 `cmd/wisp` 面。★必答题：把新增那枚候选删掉（只删本腿自己加的夹具）再重放档位突变，`TestCrowdedWorldStillOffersTheBrake330r2` **变回绿＝答"是"**（牙确实挂在多候选世界上）；另一枚 `TestResidencyTierDecidesTheInjectionArm330r2` 仍红，接住它的是**枚数前提门 `:221`**（那扇门按名字提名了被删的候选、与档位无关），具名交裁、⛔ 算本腿的 G2 牙。台面与颜色：仓外树八件套＋仓根 `docs`＋手工三件＋`GOTMPDIR` 指 `D:` 卷，改前（锚 `feab45e3`）`rc=0／RUN 110／PASS_TOP 86／FAIL_TOP 0／SKIP 0／PANIC 0` ⇄ 改后（含本腿两枚文件）`rc=0／RUN 112／PASS_TOP 88／FAIL_TOP 0／SKIP 0／PANIC 0`，红名册作差＝0、同包⛔ 一枚别枚用例被打红（故无需改任何既有期望值）。★零产码尺：`git diff --name-only feab45e3 HEAD -- internal/agent | grep -vE '(_test\.go|^internal/agent/testdata/)'` ＝空、rc=1；同一把尺的正控跑在 `git clone --local` 的副本里（把 `internal/agent/inject.go` 一处非注释字面量 `MaxRetrievedTools` 改成 `MaxRetrievedTools+1`）＝非空、rc=0 ⇒ 这把尺不失明。门禁四数只在共享工作树：`gofumpt -l internal/agent/` 列出 `approval/pending_read.go`／`tools.go` 两枚＝本锚点既存漂移（`git status` 对它们为空、⛔ 本腿碰过，本腿两枚文件不在名单）；`sh scripts/d22scan.sh` rc=1，唯一 finding＝`cmd/wisp/panel_snapshot_deliver_329.go:14 [phantom-citation]`（并行腿 329-r1 的在飞字节，本腿新文件被提及 0 次，具名交编排者）；长度门 rc=0 VERDICT GREEN；`go vet ./internal/agent/` rc=0 输出 0 字节。G3 随批＝`plan_present_doorbell_165r3_test.go:131` 那句注释改成"冗余非零门，牙在 :145 的 `rowOfTool165r3` 独立遍历"（1 删 3 加、221→223 行、新串 `6542773d3d2edcc6…`；⛔ 删那扇门；"把 sink 计数喂成常量"那一发本腿未复测，⛔ 记本腿读数）。基线件一字未动：`internal/agent/testdata/golden/plan-present.sse` 仍 `7d26fa889b59fc3a…`。证据件＝`docs/evidence/s1/330-r2-g2-teeth.md`（102 行／9647 字节），五发原始读数＝`.scratch/wisp/probes/330/r2/logs/m1..m5.log`。AC#1/#2 仍空（等 330-r1 在 `cmd/wisp/**` 落对接钉）；AC#0/#3/#4/#5/#6/#7 本腿⛔ 自翻、交裁决腿。
- **330-v1 终裁（2026-10-17 10:3x +08·非实现者裁决腿·只追加这一条、⛔ 翻任何 AC 框、⛔ 改上面一字）**：料＝`330-r2-g2-teeth.md`＋`probes/330/r2/**`＋上条 bullet；读数取具名提交树 `b9069dfa`（自建仓外树 `D:\tmp\w330v1`、两枚目标文件 sha 逐位等基线 `6542773d…`/`891fbf5c…`、`b9069dfa` 后 `internal/agent` 零枚别腿字节）。★编排者没复跑 `330-r2` 那对成对突变，本腿**自己重放**（⛔ 抄它读数），判语件＝`docs/evidence/s1/330-v1-verdict.md`（109 行／9316 字节、分词形未填标记自尺 0 命中）。逐格：**AC#1 不适用**（对接钉今天不存在＝`grep planPresentName165r3 cmd/wisp` 零命中，去重放暴露"钉不存在"，属 `330-r1` 的 `cmd/wisp` 面）；**AC#2 不适用**（同因，且 `330-r2` 夹具走组合 `planPresentDirectory165r3()`＋头注 `:66-70` 明言⛔ 声称对接真名册，没造假对接）；**AC#3 成立**（⛔ 让注释改口背书：本腿独立把 `sinkCount165r3` 喂常量 `n:=1` 单跑 165 用例＝rc=0 仍绿，证 `:131` 非零门确冗余、牙真在 `:145` 的 `rowOfTool165r3` 独立遍历 `:95-105`；票面"牙在 :143"是改注释前旧号、现量真实行 `:145`；零产码尺空 rc=1＋正控改 `inject.go` 字面量→尺输出 `inject.go` rc=0＝不失明）；**AC#4 成立**（真把两臂世界造出：M1 镜像 `Resident true→false`⇒红 `:188`＋`:240`、165 用例同发仍 PASS；★成对两发落点不同——去牙 M2⇒唯一红 `:247`、反控 M3⇒唯一红 `:254`＝档位门⇄检索门各响、⛔ 恒响）；**AC#5**＝internal/agent 半格成立（本腿独跑整包 forward rc=0／RUN 112／PASS_TOP 88／FAIL_TOP 0／SKIP 0／PANIC 0、`.go`=750 证非空、⛔ 与 330-r2/编排者的数相减＝三把台面不同）、cmd/wisp 半格不适用；**AC#6 成立**（只算 internal/agent：`gofumpt -l <两枚交付件>`＝空 rc0、d22scan 现量 rc0、长度门 rc0 GREEN、`go vet ./internal/agent/` rc0）；**AC#7 不适用**（量排程纪律非字节：`330-r2` 只写 `internal/agent`、与 329 的 `cmd/wisp` ⛔ 同枚共享编译面，且今天 329-r1 已死、cmd/wisp 写面交 228-r3，判语⛔ 编"并发写同面"）。**必答题**：G2 钉住的机制今天被产码读到＝`internal/agent/loop.go:573` 调 `Inject`、臂切分 `inject.go:47`、检索削名额 `inject.go:55/:238`（逐枚 file:line、⛔ 靠用例绿或 `go build`/`vet` rc 背书接线）；G1 对接钉不存在⇒真名册 `internal/tools/plan_present.go:72` ⇄ 镜像 `plan_present_doorbell_165r3_test.go:38` 今天零枚用例对接。AC#1/#2/#7 归 `330-r1` 的 `cmd/wisp` 面，本裁对象⛔ 交付；六框翻否归编排者。
