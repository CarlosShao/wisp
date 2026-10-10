# `305-a1` · `90` 现态＋我⛔ 做到的＋越界自查＋与派单/票面冲突处（具名）

## A. 现态（交件名册，逐件有正文，⛔ 占位、⛔ 0 字节）
`.scratch/wisp/probes/305/a1/` 下：`00-anchor.md` · `10-table1-live-document-roster.md` · `20-table2-attribution.md` ·
`30-table3-push-dimension.md` · `40-table4-cost-and-pump-answer.md` · `90-final.md`（本件）·
`logs/l1-anchor.txt` · `l2-call-sites.txt` · `l3-library-hops.txt` · `l4-earlier-green-readings.txt` ·
`l5-suspect-70b00885.txt` · `l6-entry-html-surface.txt` · `l7-anchor-verify-head-blob.txt` · `l8-document-literals.txt` ·
`l9-dependency-surface.txt` · `l10-oob-selfaudit.txt`。
**每把尺的件各自落一行 `rc=`**（`l1`-`l10` 末节逐字可见；`l2` 里 `EVAL_PROD_HITS=0` 那行是**自证计数行**，`l2` 的"导航调用"一节是**空集＝结论**）。

四张表各一行结论（带尺名）：
1. **表①**（尺＝逐枚调用点 `文件:行`＋调用序，HEAD blob，`logs/l7` 28 枚锚全 MATCH）：冷启动序上**最后一次写文档的是 `:486` 入口（1044B）或 `:467` 公告（181B），⛔ `:892` 探针（136B）**；而两枚症状页面自报的签名（`"0"` ＋ `title=""`）在三份候选里**唯一匹配探针页** ⇒ **最后一次 `SetHtml` 的那份文档⛔ 是 `Eval` 实际落到的那份**。
2. **表②**（尺＝逐字回读编排者已有五发＋`docs/evidence/s1/` 既有读数尺）：**nail2 ⛔ 由 `fb2fb802` 造成**（clone 同台面 `r3`↔`r4` 逐字同句）★**但"从来没绿过"⛔ 成立**——`33-panel-host-c27-r7.md:201/:207`（HEAD `416d9d56`，带 bundle 台面，整包 240-0-0）里 AC13 与 nail2 **两枚逐字 PASS**；嫌疑笔 `70b00885` 现量＝**纯 extract-function（同语句、同顺序）⛔ 改任何交接**，只把过期诊断文字搬了家。
3. **表③**（尺＝逐跳调用点＋依赖 blob 的导出面）：push 那一维缺**两截**——库侧 `pkg/edge/chromium.go:144` 的 `ExecuteScript(script, 0)`（**NULL callback ⇒ ⛔ HRESULT／⛔ 脚本结果／⛔ 完成信号**）＋与 AC13 共享的"最后一份文档⛔ 落地"；**分层假设⛔ 作废**（凭据＝clone 上最后文档是 181B 公告、母仓上是 1044B 入口，nail2 四发一句没变 ⇒ 它的红⛔ 依赖入口交接成功），而 nail2 独有那一截 AC13 判据覆盖⛔ 了。
4. **表④**（尺 1＝`grep -rn -E "\.Eval\(" --include=*.go cmd internal | grep -v _test.go` ⇒ **生产调用者 0 枚**；尺 2＝`logs/l9` 控制面调用计数）：修②的"页面先自报就绪再推"一形**⛔ 撞 `A502` P1**（`Eval` 本来就在泵线程上跑）；"要回程"一形**撞**，撞处＝依赖面（`common.go:26` 那个导出 interface ⛔ callback／⛔ 订阅点）＋泵面（`A798` 甲形逐字"泵形状要重来一遍"，锚＝`A502` P1），且与〔待人拍板〕`Q-84` 同格 ⇒ **具名报回，⛔ 本腿绕**。

必答两问：
- **同源与否**＝**部分同源，⛔ 一枚**：共享上游（最后一份 `SetHtml` 的文档没落地成 `Eval` 的目标），nail2 另多一枚只有它才暴露的下游（第一发的写⛔ 进到第二发读的 DOM）。裁"同源/⛔ 同源"的唯一一发＝票面 `AC#3` 那枚跨维度负控。
- **撞不撞 STA 泵**＝**分形状**：甲形⛔ 撞；乙形（要回程）**撞，且撞两处已裁**，交编排者裁。

## B. 我⛔ 做到的（照直写，⛔ 粉）
1. ⛔ 任何**新颜色**：本腿⛔ `go test`／⛔ `go build`／⛔ `go vet`／⛔ 任何 `cmd/wisp` 真跑（只跑过 `go list -m -f '{{.Dir}}'`，`logs/l3` 逐字有 `rc_golist`）。⇒ 表② 里"同一枚台面的绿↔红对照"**仍欠一发**（R2），我在 `20-` §C 末已具名，⛔ 拿跨台面读数冒充。
2. ⛔ 判出 nail2 缺跳 (i) 的**机制**：`脚本被丢` / `跑在另一份文档实例` / `文档在两发之间被重建` 三种形状，静态读码⛔ 分得开（都⛔ 走同一枚无回程通道）。⇒ 欠 R1＋R4。
3. ⛔ 一枚 `AC` 框都没翻；⛔ 对 `AC#1`/`AC#2`/`AC#3` 下任何判语；⛔ 建议改判据（`AC#1`/`AC#2` 那句"⛔ 放宽"我原样当硬约束）。
4. ⛔ 取得 10-01 那份 `frontend/dist` 字节：那枚文件在 HEAD ⛔ 受跟踪（`git ls-tree` 只回 `.gitkeep`）⇒ **台面第二轴⛔ 复现**（`20-` §E 的 R3 已登记为"做不到"，⛔ 假装能跑）。
5. ⛔ 判 `NewWithOptions` 那一跳给⛔ 给一份初始文档（表① §A 第 5 行写成"未定性"，⛔ 猜）。
6. ⛔ 任何一发起 winlive 档：`-tags winlive` 那一发今天⛔ 人跑过（凭据 `panel_host_windows_live_test.go:18` 逐字 `winlive has NO CI job`），我只把它**列为缺的那一发**。
7. ⛔ 动 `frontend/src/**`、⛔ 动 `SLO`/`golden`/`thresholds.go`/D 表 C 表、⛔ 动别人的脏面（`l10` 逐字可见那 12 行 `M` 原样在盘）、⛔ push、⛔ worktree/checkout。

## C. 越界自查（尺＝`git show --name-only --format=` **逐笔**，⛔ 区间尺；件 `logs/l10`）
- 本腿三笔（`5f891575`／`53c854e3`／本笔）名册**全部落在 `.scratch/wisp/probes/305/a1/**`**，越出＝**0 格**。
- 新建 `.go` 枚数：**逐笔 `=0`／`=0`**（`l10` 那两行 `go_files=0`）⇒ `gofumpt`/`d22scan` 分母**零移动**。
- 证据件命名：⛔ 一枚 `.out`（全 `.md`／`.txt`）。
- commit 全部**显式 pathspec 写在 `$( … )` 之外**（`git add -- <三枚具体路径>`）；⛔ `add -A`／`.`、⛔ `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`。
- 闸门自计：真·读码/尺调用约 20 发（⛔ 超 25 那枚目标调用闸），第 45 发硬 gate 之前已完成三笔 commit；大原文（CI 日志 1.4MB）⛔ 重新抓远端，逐字两句用派单已落的 `r9` 件。

## D. 与派单／票面**冲突**处（具名，逐条给尺；⛔ 我据此改任何规矩）
1. ★**票面标题那句「`Eval` 推不进文档那一枚早在回归之前就红、从来没绿过」** ⇄ 盘上：
   尺＝`grep -rn -E "PASS: TestAC13ColdStart|PASS: TestAC14GoSideEvalPush" docs/evidence/s1/` ⇒
   `33-panel-host-c27-r7.md:201` AC13 **PASS(0.97s)**、`:241`、`:277` 亦 PASS；`:207` **`TestAC14GoSideEvalPushReachesThePage` PASS(0.50s)**；
   `33-panel-host-c27-r6.md:36` AC13 **PASS(1.47s)** 单跑。身份（`r7.md:10`／`:23`）＝HEAD **`416d9d56`**、整包 **rc=0／240-0-0**。
   ⇒ **"从来没绿过"⛔ 成立**；票面 判语① 只否掉了 `fb2fb802` 那一支。**这直接决定 `AC#2` 该⛔ 该按"缺功能"降级＝编排者裁**，本腿⛔ 替。
2. **票面第 18 行「AC13 在母仓⛔ 已知'改前绿'过」**：在票 303 那一对改前/改后的范围内**真**（改前那发红句是"回执为零"，通道挡在前面）；⛔ 覆盖到 10-01 那批发。⇒ 两句话射程⛔ 同，引用时必须带"哪一对"。
3. **派单第 4 节「票 33 的 `A475` 那批裁过'两维⛔ 混一枚'」**：现量**⛔ 对不上**——`grep -n "A475" docs/reports/pending-and-issues.md` ⇒ `9890` 那格是**票 228-r1／球住进会跑任务的那条腿／新立票 245**，⛔ 面板两维。
   "两维⛔ 混一枚"实际在 **10-01 13:12 那道 P2**（凭据逐字：`panel_resident_windows_test.go:11-16` 用例头 ＋ `panel_host_windows.go:860-862`；台账那一批＝`A502`，`:10287`）。⇒ 建议派单口径改指 `P2`／`A502`，**⛔ 我改票面一字**。
4. **`:660`/`:662`/`:665` 那族行号属 `cmd/wisp/panel_host_windows_test.go`（该用例从 `:605` 起），⛔ 属 `panel_resident_windows_test.go`**；本腿按后者去核 `:660`，那一条⛔ 命中（逐字打印在 `logs/l7` 末节）。⇒ 这是**引用面写错文件**那一族，⛔ 行号腐烂那一族，但会让下一腿白找。
5. **票面 `AC#0` ①「每一次 `SetHtml`／导航调用」**：现量生产面**⛔ 任何导航调用点**（尺＝`logs/l2` 的"导航调用"一节＝**空集**；库内唯一的 `browser.Navigate` 在 `webview.go:390`，本仓⛔ 调用者）⇒ 活文档只由**三枚 `SetHtml`** 决定。
6. **票面"台面规则"只覆盖了一根轴**：它写的是 母仓（带旧 `frontend/dist`）↔干净 clone（⛔）；盘上还有第二轴＝**同一枚母仓里那枚未跟踪 dist 的内容自己漂过**（`logs/l6`：mtime `2026-10-10 08:51`、1044 字节、`:13-16` 带冻结 CSP）。10-01 那批绿读数用的是**旧字节** ⇒ 跨时间的"母仓↔母仓"对照⛔ 天然同台面。⇒ 具名欠账，⛔ 本腿补规矩。
7. **判据串里那句现在时诊断**（`panel_resident_windows_test.go:330` 逐字 `the round-trip probe page is the last document shown … bringUp must hand the entry over AFTER the probe`）在 `a81c2980` 的源码序上**⛔ 成立**（`logs/l5` 的 `70b00885` diff＝同语句同顺序）。⇒ 它现在是**待验断言**；本腿⛔ 动那句一字（`AC#1` 硬约束⛔ 放宽判据），只报"文字与序⛔ 一致"。
8. **注释比读数乐观一处**：`panel_resident_windows_test.go:14` 逐字 `(ii) a Go-side Eval push, **which arrives even without it**` ⇄ 今天 nail2 四发都红。⇒ 引用它当既成事实＝把注释读得比代码乐观（表③ §D）。
9. **`70b00885` 是〔仅腿报〕**那一格，票面第 19/48 行要求⛔ 当归因结论——本腿现量后**主动否掉**它（纯搬函数），⛔ 把它写成结论。
