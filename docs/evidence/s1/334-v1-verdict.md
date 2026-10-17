# 票 334 · 非实现者终裁判语（裁腿 `334-v1`，2026-10-17）

我没写这枚票的任何产码。下面每一格都写清**哪一层是我自己伸手真跑出来的**、**哪一层是读码推的**。
读数与日志的逐字出处＝`.scratch/wisp/probes/334/v1/readings-and-mutations.md` ＋ `logs/**`（逐枚已入库）。
台面＝三棵仓外导出树（`D:\W334V1TreeBefore`＝`03f7e7a6`／`…After`／`…Mut`＝`3a118226`），
⛔ 共享工作树跑过任何一发 `go test`／`go build`（门禁四数只在共享树，那是它要求的地方）。

## 起手认锚（⛔ 抄派单给我的号）

`date=12:35 +08`、HEAD=`58806772`、跟踪件零脏（`git status --porcelain --untracked-files=no` 空）。
派单说腿的起手锚是 `6e4abaf1`、腿自报现认 `03f7e7a6`：我把那把尺自己跑了一遍
（`git diff --stat 6e4abaf1 03f7e7a6 -- internal/risk/` ＝ 空输出 rc=0）⇒ **两枚号下 `internal/risk/**` 字节全等，成立**。
票面七格现跑：未勾 **7**、已勾 **0**；★票面 `AC#6` 正文那句"本票六格"是编排者笔误，我按七枚裁，已具名。

## 逐格判语

**`AC#0` 起手现量＝成立。** 六行我逐枚重跑（读数在读数件第 5 节）。两枚具名差异：
③′ 腿不在机主 profile 里建目录、改用例自造世界复现同一读数——我也用我自己的 solo 世界复现了
（`chain:case5-aws` 点名 `.aws` ⇄ `chain:case6-kube` 点名 `.kube`），**认这座桥**；
④"LOCALAPPDATA 挪走仍 A"那一发我没重造，**我不替它背书**，但同一行里"链内 A ⇄ 链外 B"我有九对独立读数。
本票加的那枚分母我数得出：**链内 13 枚读 A、同层干净文件 10 枚读 none、链外 9 枚读 B**（⛔ 全称句）。

**`AC#1` reason 说实话＝成立（带一枚具名缺陷）。** 真跑：M-V4（reason 拼装那处 → `if true`）、
M-V5（锚展开那行 → 丢掉 root 的裸形）、**外加一发腿没发的 M-V6**＝只把那句实话换成**另一种假话**
（仍带 `fail-closed` 字样、不提候选也不提缺位锚）⇒ 三发都只把 ⓐⓑ 两枚 reason 用例咬红、闸门格仍绿，
落点两两不同、缺任一发那把尺就失明。读码：拼装点唯一（`blacklist.go:329`），九枚 case 的条件表达式我逐枚对过 diff＝未动。
★缺陷（不翻判据、要落账）：**票面说"人（和模型）看到的那句话"这一半在盘上不成立**——见必答①，
那句 reason 今天**没有任何生产读者**，所以 AC#1 买到的是"返回值与一枚还没接上的日志 sink 里说实话"。

**`AC#2` 判级方向一字不许动＝半成立。** 站得住的那半（全是我自己跑的）：
改前⇄改后同一批具名输入的 `class／allow／needL2／点名规则` 在**我这把 42 行的尺**上逐字相等（`diff`＝0 行；
腿的那把是 14 行）；`rules_gateway` 与 `internal/tools/` 对它零字节（numstat 空）；豁免键买不到放行那一发我复跑了。
牙齿也在：M-V1（只动判级支 `aRelates` 尾行＝回落支整支，reason 一字未动）当场把闸门格咬红，
红句三面镜子（`:193`／`:197`／`:201`）是同一次降级的三次观测，不是两件事——这把尺**认得"整支被摘"**。
★不成立的那半＝判据写的是"**任何**试图把这一形降成 B 档／可豁免的改动必须被咬红"，我实测到两类反例：
**M-V2a**（只摘 `%LOCALAPPDATA%\Microsoft\Credentials` 一枚 case 的回落支）与
**M-V2b**（只摘 `~/.ssh/**` 一枚 case 的回落支）——两发都让同一形从**不可豁免的 A** 落到
**可单文件豁免的 B**（我的尺逐字看到 `chain:case9 A→B`、`chain:case4b A→B`），
而腿的三枚仪器**全绿**（`10／3／0／0 rc=0`）。腿的仪器真正钉住的锚只有 `.aws`（闸门格＋ⓐ）与 `.kube`（ⓑ，
M-V2c 那发 ⓑ 当场红就是证明）；其余五枚 aEq/aUnder 锚上的这一形**没人钉**。
另记一发方向性的：M-V3（把回落支的祖先关系界去掉＝**变严**，整台机器都算 A）闸门格**恒不红**，
红的是 AC#1 那两枚与票 216 那枚——那属于偶然接住，不是这格设计的敏感方向。
⇒ 判 `AC#2` **半成立**，缺陷具名如上；补法（我⛔ 动手）：闸门格按锚逐枚各发一两枚用例，
或把我那 42 行对撞搬成仓内断言（世界已现成，见 `logs/probe-334v1_test.go.txt`）。

**`AC#3` 给仪器补世界＝成立。** (a) 成对反控我自己伸手跑（⛔ 抄腿的 `.log`）：改前字节 `24／5／1／0 rc=1`
⇄ 交件字节 `24／6／0／0 rc=0` ⇄ **把夹具改回 `t.TempDir()` 那一发回到 `24／5／1／0 rc=1`，同枚名册同枚用例同三发子用例**
＋夹具未动的正控 `24／6／0／0`；改的那 1 行是取材，四对断言与 `want` 值一字未动（我逐字 diff 过）。
⚠ 我第一发反控把 `GOTMPDIR` 放到了 `D:` 卷＝**台面缺陷**（那一发红本身活在"临时根落在真主目录链内"上），已用 CI 形重跑。
(b) 世界造对了：尺＝候选 root 与七枚锚展开 root 的**公共前缀段数**（全是 9＝只到 temp 根）＋
`uncertainAnchorMiss` 实值逐枚 `false`＋`offchain:*` 九行全读 `B` ⇒ 候选真在 A 锚祖先链**之外**；两臂不同哑。

**`AC#4` 字面串消费者普查＝成立。** 我的尺带射程（`3a118226` 全仓跟踪件，剔 `.scratch/**`），先跑正控
（`A-list sensitive path` 命中 2 枚文件＝生产者＋仪器自己）再报数：`fail-closed, NOT a match` 只命中生产者 1 枚；
各锚字面串在测试里只做**子用例名／注释**（`pathshape_portable_test.go:134` 那张表比的是 `Classify()==ClassA`，⛔ 比 reason）；
`design/doubao/demo/screens/security.js`＝机主外部 mock；台账里那句＝散文 ⇒ **判定型消费者 0 枚**，腿"无需停手"这一句认。
（我自抓一发尺坑：`Microsoft\Protect` 用单反斜杠 `-F` 跑会假 0，真实串是 Go 源码里的双反斜杠形。）

**`AC#5` 门禁四数＋颜色＝成立。** 四数我在共享工作树逐枚重跑：`gofumpt.exe -l` rc=0（空输出）／
`d22scan.sh` rc=0 `clean`（射程 internal=586 cmd=143 design=39）／长度门 rc=0 `VERDICT GREEN`／`go vet ./internal/risk/` rc=0。
分母 14818→**15035** 是这段时间别腿入库的差，⛔ 读成本腿的红。整包颜色我自己在仓外成对跑了
（`logs/90/91-full-*.log`，两棵树里都装了我那枚探针，故 `RUN` 比腿的 225／230 各多 22 行）：
改前 `247／140／2／1 rc=1` ⇄ 改后 `252／143／1／1 rc=1`；红名册作差逐名归位＝
转绿那枚是 `TestFileBaseNameClaimChangesNoClassOrReason`（＋三发子用例），残红一枚
`TestC26RewriteAccountIsConsumedAtEverySecurityLeg`，SKIP 两发同一枚 `TestSyncRegistryProbeLive`＝本机注册表现实。
★归位更正（这条比腿的写法更硬）：那枚残红的四行逐字指 `internal\panel\composer.go:127` 与
`display_scrub_216.go:26/37/218`，而**我这两棵树里的 panel 字节是两枚锚的已提交字节**
（`composer.go` 与 `display_scrub_216.go` 两发哈希逐枚相等，`0071996d…`／`230d3443…`）⇒ 它红在
**已入库的 panel 面**、⛔ 别腿在飞字节、⛔ 本票射程，改前改后四行逐字相同。本腿⛔ 碰那面。

**`AC#6` 纪律格＝成立。** 三枚产码 commit 的 `git show --name-only` 我逐枚核过＝只 `internal/risk/**` 三枚文件；
框没翻（现跑未勾仍 7）；票名／`Status:`／`-done` 未动；证据落 `probes/334/r1/`，⛔ 提前占 `docs/evidence/s1/`；
判级那一支一字未动（我也没看到"为了绿做的等价事"——M-V2 那两类放松**没有**出现在交件字节里）。
★票面"未定案那一支"（要不要把这一形降回可豁免的 B）本腿**没动也没替它铺路**；我这轮多出一枚事实请一并入账：
误伤分母不止腿数的那几枚——**七枚** aEq/aUnder 锚都能各自成为这一形的落点（我只实测了其中 4 枚的 A 形读数）。

**事故那一格＝只报盘上事实（纪律结论归编排者，我⛔ 处置）。**
`git show --stat 8823549f` 里那行 `…census-done.md} | 0`；`--name-status`＝`R100 旧名 → 新名`；
`--numstat`＝`0 0` ⇒ **行丢失枚数 0**；`dc367a2a` 的 parent＝`8823549f`；旧名的消失笔＝`8823549f`（`--diff-filter=D`）；
`dc367a2a` 在同一枚件上翻了 **6** 枚框（diff 里 `- [ ]`×6 ⇄ `+ [x]`×6，`14 insertions/7 deletions`），
该件 `dc367a2a` 的 blob＝HEAD 的 blob（`1d6b147e…` 同串）；HEAD 树里 332 那枚**只命中 `-done` 一个名字**；
盘上现值未勾 0／已勾 6。⇒ 内容零损失，损害只有**归属**：那枚防重领键的产生记录写在腿的提交号下。

## 派单末尾两枚必答题

**① 那句 fail-closed 长话今天到人眼的射程＝哪几跳、有没有落进模型正文或快照字节。**
链（每一跳 `file:line` 都在 `3a118226` 上现跑过 `git grep -n`，逐字见 `logs/80-chain-reach.log`）：
`classifyForms` 第二返回值（`blacklist.go:221…236`）→ 唯一拼装点 `:329` → `Gate` 装进 `PathDecision.Reason`（`:77`/`:84`）
→ **`risk.Gate` 全仓只有一枚生产调用点**＝`internal/tools/mode.go:98` → 写进 `BlacklistNote.Reason`（`mode.go:111`）
→ 而这枚字段**读者为零**：`Blacklist` 被读到的只有 `pending_read.go:73-75` 克隆的那三枚路径切片，
`mode.go:120` 那句 `blacklistedSensitive()` 连它自己都没有调用点（`git grep` 只命中定义处 2 行＝定义＋注释），
`sensitiveClassifier.Classify`（`paths.go:380`）只取 tier、把规则串丢掉；
另一出口 `risk.Logf`（sink 变量在 `blacklist.go:51`）在生产码里**没有注入点**（`git grep 'Logf = '` 剔测试后 0 命中）。
到模型的那句**不是**这一句：它是 `rules_gateway.go:82` 的 `R3: 命中 A 档敏感路径，拒绝（不可豁免）: <canonical>`
→ assessor 融合（`assessor.go:345` 一带）→ `bridge.go:301` → `route` 的 Deny 支 `bridge.go:439`
→ `reject` 的 `ToolOutcome.Text`（`bridge.go:1105`，函数注释自陈"returns what the model sees"）
→ `loop.go:745` `toolResultMessage(c.ID, log.Text, …)` 进对话正文。面板那两跳（`panel_pump.go:75`、`run.go:785`）
渲染的也是 `Decision.Reason`＝R3 那句，⛔ Gate 那句；而 A 档在 `route` 里直接 return，根本不生成卡片。
⇒ **结论：这句长话今天不进模型正文、不进面板快照字节；票 334 的修法没有新引入出向泄露面。**
同时给出另一面（这是票面的前提要更正的地方）：它今天**也不到人眼**——零读者。
所以"理由指错东西"今天实害的是**返回值消费者与审计日志**（一旦有人接上 `risk.Logf` 或给 `BlacklistNote.Reason` 找读者，
那句里新添的真实候选路径与控制字符、绝对 profile 路径就一起进那一面）⇒ 接线那枚票落地时必须带上"这句会带真路径"这一条。

**② 什么样的产码改动能让 `AC#2` 闸门不红、而 `M-A`／`M-B` 全绿——它的失明方向。**
实测到两枚（不同落点、同一机理）：**把单枚 A 档锚 case 的判级从 `aUnder/aEq` 换成不带回落支的 `anchorMatches`**。
M-V2a 摘 `%LOCALAPPDATA%\Microsoft\Credentials`、M-V2b 摘 `~/.ssh/**` ⇒ 腿的三枚仪器 `10／3／0／0 rc=0` 全绿、
AC#1 那两发的基线行为不受影响，而我那把 42 行的尺当场看到 `A → B`、`needL2=false → true`（＝降到可单文件豁免）。
读码解释（这层是推的）：switch 按序取**第一枚缺位锚**，腿的三枚仪器只把世界造到 `.aws` 与 `.kube` 缺位，
其余五枚锚上的这一形从未被任何断言碰过。⇒ **这枚闸门买的到底是哪一件事**的尺面答案是：
"`.aws`（闸门格＋ⓐ）与 `.kube`（ⓑ）这两枚锚上的这一形不许降档"，⛔ "这一形在任何一枚 A 档锚上都不许降档"。
反向也具名：M-V3 那发**变严**（去掉祖先关系界）闸门恒不红。

## 交件自证（三把尺）

1. 逐枚提交号＋逐枚 `git show --name-only` 全集见 `logs/95-final-selfcheck.log`；写面只三处＝
   `.scratch/wisp/probes/334/v1/**` ＋ 本件 ＋ 票 334 末尾一条只追加 bullet；
   **本腿 `.go` 命中 0 枚**（尺＝逐笔 `git show --name-only` 里 `grep -c '\.go$'`，⛔ 子串匹配主题）。
2. 件：本件＋`readings-and-mutations.md`＋`logs/**` 逐枚 `wc -c` 非 0，见同上日志；
   占位符尺（`TODO|FIXME|placeholder|占位`）在本件与读数件 0 命中，正控＝同尺在票 334 正文命中≥1 枚。
3. 决定性尺逐字（起手侧⇄交件侧）：M-V2a／M-V2b 的"落点行差＋四数＋我这把尺的行差"逐字在
   `logs/60-v2-runner.log`；还原哈希每发＝`c1799cd24afce19e6accef23db23a7d7ee746dec`（＝交件 blob），末态三枚件哈希逐字同串；
   门禁四数 rc＝0/0/0/0（`logs/81-gates-shared-worktree.log`）；起手与收尾各跑一次
   `git status --porcelain --untracked-files=no`，收尾新出现的那枚 ` M .scratch/wisp/issues/335-…md`＝**别腿在飞字节，本腿⛔ 碰过**。

⛔ push；⛔ 翻框；⛔ 改判据文字／票名／`Status:`；⛔ 动 `internal/**`／`cmd/**` 一字。
