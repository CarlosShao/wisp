# 232 — 票 223 那枚"重启档要告诉用户为什么不生效"的用例**只钉了"有没有那句"，没钉"那句里有没有为什么"**：把解释三段删掉，它照样绿

- Status: **待派，测试面加固**（与票 231 同撞 `cmd/wisp/config_reload*.go` ⇒ 同批或串行）。来源＝非实现者验收腿 `223-v2` 第②节（它**造出了那一发突变并跑出绿**，不是我推的）＋台账 `A444`。⚠ **形状不是 r2 引入的**（`223-v2` 具名：归 r1 那发形状），但今天还在。
- ⚠ **这不是安全漏洞，照三行读**：①现象＝**我们自己的用例**在一种产品文案下不报警；②**没有本机被入侵的证据**；③最坏后果＝**产品那句提示变空壳而测试仍绿**（用户看到"本次不生效"却看不到为什么，而门不会响），**不是**任何权限被放宽。

## 现量（起手逐条复算）

| 事实 | 读数 | 尺 |
|---|---|---|
| 内容断言在哪 | `cmd/wisp/config_reload_223_test.go:482-486`：三枚内容 needle 在 **`why+out` 拼接串**上找 | 腿现读＋写腿复核 |
| ⛔ **瘦句不响** | 突变 B＝**保留那枚 needle、把"原因／涉及哪些段／两件事都没发生"三段删掉** ⇒ 用例 **PASS 2.24s**（红名＝无） | 腿实测（它自己造的那发） |
| 方向说反**有**牙 | `:481`（等不到那句就 `t.Fatalf`）＋`:495-497` 的**负向断言** ⇒ "把重启档说成立即生效"这一形会红 | 腿现读 |
| ⚠ 一枚**前向**风险（今天不红） | needle 是全文匹配、**不限"在种下改动之后"** ⇒ 若启动横幅里也含同样字样，用例会 **0ms 假绿**。今天产码里那字样的唯一来源是 `config_reload.go:289` | 腿 `grep` 现跑＋具名标"前向" |
| 助手本身干净 | `awaitStdout`（`:140-158`）与既有 `awaitAudit`（`:109-127`）**同构**：同 `reloadCaseBudget`=40s、同 Timer/20ms Ticker、同 `Fatalf` 出口、**无 `Skip`** ⇒ 期限到是**红**、不会被伪装成绿或被吞成 skip | 腿逐行对拉 |

## 后果

1. **"告诉用户为什么"这一半，今天只由审计日志钉着，操作员面上那一半是空的**——产品文案可以被删薄到只剩一句"本次不生效"，测试仍全绿。这正是票 223 AC#7 字面要的"**为什么**"，所以它**不能算已钉死**（编排者据此把 AC#7 留着不勾）。
2. 前向那一枚更贵：横幅文案一改（那是常改的），这一整格判据会**静默变成恒真**，而没人会去看。

## 判据（每格都要现跑；⛔ 不许用"我把断言写长了"充当修好）

- [x] **AC#1 复现瘦句不响**：把 `223-v2` 那发突变 B 自己重做一遍（删三段、留 needle）⇒ **必须仍绿**，抄原始读数。⛔ 不许跳过这一格直接改断言——**不先复现就不知道自己在修什么**。
- [x] **AC#2 内容断言改到 stdout＋窗口化**：三枚内容 needle **只在 stdout 上找**（别在拼接串上找），且**只搜"种下改动之后"那一段输出**（带起点标记）。⚠ 期限仍用现成的 monotonic 形状，⛔ 不许写减法、不许裸起协程（`d22scan` ban #4／#1）。
- [x] **AC#3 反向正控要真响**：改完后重做突变 B ⇒ **必须红**；并补一发新突变（把"原因"那段换成一句无关话）⇒ 也必须红。红句原文逐枚抄进表。**若两发都仍绿＝你这发断言还是装饰**，具名写出来。
- [ ] **AC#4 横幅假绿那一枚也要钉**：种一发"启动横幅里先含同样字样、之后才种改动"的文件／参数，断用例**不许 0ms 通过**（起点窗口化之后它自然会红或变慢——读数写清是哪种）。
- [x] **AC#5 零放宽**：⛔ 一条既有断言都不许删、不许改成"或满足任一即可"；本票动完之后 `TestTicket223RestartTierSaysItWillNotApply` 与 `TestTicket223R2FailureSentenceRouting` **必须仍全绿**（逐名抄回）。
- [x] **AC#6 整包终态**：`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp ./internal/... -count=1`＋逐名比红名册（历史 5 枚属别人地界；`internal/risk` 那枚争用型假红安静复量为准）。

## 禁区

不动 `PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt`；三枚冻结件一字不动；⛔ **不改产品文案来让测试变绿**（那正是本票要防的"文案瘦成空壳"）；`frontend/**`／`design/**` 零读零写零转述；不新增导出名；⚠ 与票 **231** 同文件面（`cmd/wisp/config_reload*.go`）⇒ 同批或串行；与票 222／221／224／228 同撞 `cmd/wisp` ⇒ 串行。
## Progress log (append-only, newest last)

- [2026-10-06T13:06+08] agent=232-r2 did=AC#1-AC#6 交读数：起手锚 HEAD `590fb855`、三把尺原文（第 3/4 行零 WITHDRAWN/撤/作废、scoped porcelain=0 行、票面 `- [ ]`=6 全不归本腿）；**AC#1 复现成立**＝-overlay 种 B（三段删薄）＋本腿自补 C／E／F 共 4 发在 **HEAD 未修态用例**上逐枚 PASS（2.28/2.29/2.50/2.58s）；**AC#2 落笔**＝`config_reload_223_test.go` 取 pre-plant mark ＋ 文件内小写辅助 `windowSince232`/`stdoutSince232`/`awaitStdoutSince232`/`awaitAuditSince232`（零导出、零产码改动），三枚 needle 拆到具名流＋窗口，并实测（logs/P0）出 `重启进程后生效` **只在审计句**、操作员句里没有 ⇒ 那枚只能钉审计流，钉 stdout 需改产品文案＝本票禁区，具名上报；**AC#3 六发全红**（B 5 枚 omits／C 2 枚／E 1 枚／F 1 枚／J `never carried AFTER the plant` 42s／K 同句 43s）；**AC#4** 横幅形 M 在 HEAD 用例绿（4.19s）→ 终态用例红（43.8s），起作用机制＝起点窗口不是变慢，通知完好的正例 G 仍绿＝不误伤；**AC#5 零放宽**＝10 枚 `TestTicket223*` 全 PASS、`-count=5` 五发全绿、负向断言故意保持全流读；**AC#6** 整包 9 枚红逐名比＝ball 1／panel 4 在册别人地界、`TestResolvePerCallBudget` 安静 -count=3 转 `ok`、`TestAC1ResidentLeg…` 安静单跑转 `ok`（带载偶发在册），新见 2 枚 `Test258Occupied…`/`Test258V1ProbeSummon…` 用"把本腿写面退回 HEAD 的整包对照发"证明**同样红**（本机 `Ctrl+Alt+P` 被别的程序占着，`Hot key is already registered.`）⇒ 零新增红。件 `.scratch/wisp/probes/232/r2/evidence.md`；⛔ 六枚 AC 框／`Status:`／票名一字未动（append-only 小节由本腿新增，此前票面无此节）；`md5sum cmd/wisp/config_reload.go` 13 发逐次对拉＝`5ce441ca…`＝HEAD 值，零 push next=交编排者
- [2026-10-06T19:3x+08] agent=232-v1 did=非实现者终裁（⛔ 本腿不勾任何 `- [ ]`、不写"完成"、`docs/**`／台账一字未动）：起手四把尺原文取到（HEAD `d6243757`、scoped porcelain=0 行、`config_reload.go` md5＝`5ce441ca…`＝`git show HEAD:` 同值、基线三枚具名 0 枚 `--- FAIL`）。**全部 23 发读数走 `go test -overlay`，共享工作树零原地编辑**（旧形那一把尺由 git blob `ae9f346b`＝`a16d1ff7:` 抽出重建，件里已注明"不是盘上历史"）。判语＝**AC#1 成立**（旧尺上 B/C/E/F 四形逐枚 PASS 3.16/3.08/3.07/3.67s＝瘦句不响复现，正控 K 旧尺 FAIL 43.17s）·**AC#2 成立带具名让步**（六枚 needle 只在 stdout 窗口那一路，`why+out` 全文件计数旧形 1→终态 0；`P`/`Q` 两发证明具名流真的分了家：把五字搬到 stdout、旧形同发绿 4.21s 而终态红 3.63s）·**AC#3 成立**（终态 B/C/E/F/N/K/P/Q/SA/MK **10 发全红**，红句逐字在案；`MIRROR` 逐字节副本正控 PASS 2.97s 把"走 overlay 这个动作"洗清）·**AC#4 部分成立**（拦住＝`SA` 旧形 PASS 2.29s→终态 FAIL 41.45s、`MK` PASS 4.38s→FAIL 42.71s；**没拦住＝`M`（横幅先带同样字样、重启句照打）终态仍 PASS 2.62s**，根因＝窗口下界取在 plant 之前，用例没有任何一条要求"窗口内该句只出现一次/必须晚于 applied 行"；票面那句"不许 0ms"在全部 23 发里从不出现 0ms）·**AC#5 成立**（断言计数 6/2/5 ≥ 旧 5/2/4、删除列只落在被替换的 `why+out` 块、无"或满足任一"、helper 无 `Skip`；本腿具名复跑 10 枚 `TestTicket223*` 全 PASS 含 `TestTicket223R2FailureSentenceRouting` 0.62s）·**AC#6 成立零新增持久红**（整包 rc=1／8 枚逐名：ball 1＋panel 4 在册别人地界；`TestAC1AlwaysBranch…` 安静单跑 PASS 1.19s、`TestTenOpsInOneToolCallGetOneConfirm` 安静单跑 PASS 0.07s（★这枚不在派单在册名里，待你归账）、`TestResolvePerCallBudget` -count=3 `ok 4.063s`；`Test258*` 今日未红）。门禁：`sh scripts/d22scan.sh` rc=0 且正控先绿 `top-level: PASS=35 FAIL=0 SKIP=0`；`gofmt -l`/`gofumpt -l` 只喂 `.go`（本票两枚文件空＝干净；另 5 枚系本机 checkout CRLF，`i/lf w/crlf`＋`git diff --quiet HEAD` 逐枚核过与本腿无关，未触碰）；CI 形 tracked 尺 `attrib.sh --tracked-only` rc=2，根因逐字定位＝tracked 故意坏样本 `.scratch/wisp/probes/185/c1/mut/fs_broken.go:4:1`（早于本票，未动）；`go vet ./cmd/wisp/` rc=0；终态 scoped porcelain＝起手 0 行、`config_reload.go` md5 复量同值。件 `.scratch/wisp/probes/232/v1/verdict.md`（§0–§5 六节全实，零"待填"）＋`logs/**` 24 枚；待你裁五枚已具名写进 §5（AC#2 半格／AC#4"0ms"读法／★`M` 形真缺口建议立新票／`TestTenOps…` 归账／格式仪器 bench 形状）；⛔ 六枚 AC 框／`Status:`／票名一字未动（本腿只做 append-only 一行）；零 push next=交编排者

## 编排者翻勾记录（2026-10-06 19:4x +08）

凭据＝**非实现者终裁腿 `232-v1`**：件 `.scratch/wisp/probes/232/v1/verdict.md`（**195 行／22,571 字节**，六节全实、占位符尺 0 命中）＋`logs/` **24 枚**原文读数；两笔提交 `fda6a5a5`／`75691f1a`（编排者逐枚 `git log -1` 验存在；交件笔逐枚删除列＝0，票面那笔＝**1/0**）。★该腿**一枚框都没翻**（我翻前现量未勾 6／已勾 0，与它 §0 起手逐字一致），实现者是 `232-r2`（`3d9b8374`），判语全部建在**它自己造的 23 发 overlay 读数**上（12 发具名红），⛔ 未采信写腿日志。

- **AC#1 → 翻勾**（成立）。它用 git blob `ae9f346b` **反向重建出旧形**（件里注明"非盘上历史"，这跟我让 `236-v1` 补"改之前那一形"时用的是同一招），旧形上 B/C/E/F 四形逐枚**仍绿**（3.16／3.08／3.07／3.67s）＋正控 K 红 43.17s ⇒ "瘦句不响"这一形**被量到**，不是被声明。
- **AC#2 → 翻勾**（成立·带具名让步）。`why+out` 拼接计数旧形 1 ⇒ 终态 **0**，六枚 needle 只读 stdout 的种后窗口。★**让步具名**：`重启进程后生效` 那五字**只在审计句**（产品真源 `cmd/wisp/config_reload.go:317-319`，操作员那句 `:321-326` 里没有），要把它挪进 stdout＝**改产品文案＝本票禁区**。⇒ 我裁"成立"不裁"半格"：判据原文钉的是"内容断言不许在拼接串上找"，这一点实测有牙（`P` 2.61s 红／`Q` 3.63s 红而同发旧形 4.21s 绿）；"操作员那句也必须含这五字"**不是票面要求**，⛔ 我不替它加要求。
- **AC#3 → 翻勾**（成立）。改后终态 **10 发红**（B/C/E/F/N/K/P/Q/MK/SA），`MIRROR` 逐字节副本正控绿 2.97s ⇒ 红不是"载体动了"造的。
- **AC#4 → ⛔ 不翻，保持未勾**。腿量到**真缺口**：拦住的是 `MK`（旧 PASS 4.38s ⇒ 终态 FAIL 42.71s）与 `SA`（2.29s ⇒ FAIL 41.45s），但 **`M` 形（横幅先带同样字样、之后种改动、重启句照打）终态仍 PASS 2.62s**（件 `logs/AC4c-M-final.txt` 逐字）——根因＝窗口下界取在 `plant` **之前**，用例里既没有"窗口内出现次数＝1"也没有"必须晚于 `state=applied` 行"。⇒ 这一形**具名立成票 272**，⛔ 我不把它读成"部分成立就算完成"。★顺带更正我自己写在 AC#4 里的那句"不许 **0ms** 通过"：23 发读数里**从不出现 0ms**（最快 2.29s）⇒ 那三个字是我造的一个不存在的读数形状，本格按**"不许假绿"**读；原句不抹（就在上面那行里）。
- **AC#5 → 翻勾**（成立）。断言枚数 6/2/5 ≥ 旧 5/2/4 ⇒ 零放宽有数；10 枚 `TestTicket223*` 全 PASS（含 `R2FailureSentenceRouting` 0.62s）。
- **AC#6 → 翻勾**（成立）。终态整包 rc=1／逐名 **8 枚**，`ball` 1＋`panel` 4 逐名对上历史在册；★零枚新增持久红。

**我另外自己复跑的两把尺（⛔ 没把腿的自述当凭据）**：
1. 那枚不在我给的名册里的 `TestTenOpsInOneToolCallGetOneConfirm`（真身 `internal/agent/approval/batch_test.go:142`）：整包红 5.24s（`logs/G5-fullpackage.txt` 逐字）vs **我安静单跑 PASS 0.12s／rc=0** ⇒ 与 `TestAC1*`／`TestResolvePerCallBudget` 同族＝**带载型计时红**。⇒ 具名入在册（台账 `A647`），⛔ 任何腿不许为它改断言。
2. 腿报的 tracked 尺 rc=2：那把尺真身在 `.scratch/wisp/probes/161/r5/attrib.sh`（⚠ 我先在 `scripts/` 找它、找不到 ⇒ 差点下"仓里没有这把尺"的错结论），`ci.yml:204` 在跑它。★我 19:4x 现量 CI 那发：**step 8 `gofmt (gofumpt)` [failure] ⇒ step 9 `…tracked set is the denominator` [skipped]**——**step 9 没带 `if: ${{ !cancelled() }}`** ⇒ 那把 tracked 尺**今天在 CI 里从未执行**，腿的 rc=2 只是本机读数、没有 CI 后果。⛔ 但这**不是好消息**：这正是票 111 AC#5 那条"新步排到会失败的步骤之后就必须自带守卫，否则这道门从未执行"**在我自己仓里还没堵完的一枚实例**，且 step 8 常红的根因＝票 269 那枚入库夹具 `probes/185/c1/mut/fs_broken.go`（＝票 236 AC#6 那格，三支未预选）。⇒ 三条读数同批记 `A647`，⛔ 本条不动 `ci.yml` 一个字。

**剩余未勾一枚＝AC#4**，缺口具名归口票 272。本票**不加 `-done`**（六枚里五枚成立＋一枚因判据缺口未闭合）。
