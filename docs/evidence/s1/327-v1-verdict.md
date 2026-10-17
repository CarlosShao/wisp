# 327-v1 终裁：票 327 `AC#1`–`AC#4` 四格有牙没有

裁决腿 `327-v1`（非实现者·只读，独占 `internal/panel` 那包 `go test`）｜票＝`.scratch/wisp/issues/327-display-boundary-washes-cc-only-so-cf-invisible-chars-still-pass-through.md`
被裁料＝`c2f9310bd`＋`0eaf5a6a2`＋`60c328686`（第四笔票尾 bullet＝编排者代提 `0438b84d`）｜锚＝具名提交 `60c328686`，起手对照＝`7c9d7ebe`
逐字读数＝`.scratch/wisp/probes/327/v1/logs/readings.log`（台面、四数、七发突变、还原哈希、调用点尺）｜台件＝同目录 `v1-mutate.sh`、`zz_v1_exhaustive_test.go`
判语只许四形：成立／不成立／不适用／不可判。`AC#0`（普查）与 `AC#5`（门禁四数＋颜色）⛔ 本腿射程＝**不适用**。
★本腿⛔ 翻任何框、⛔ 改票名／`Status:`／`-done`、⛔ 动 `internal/**`·`cmd/**` 任何跟踪字节（七发突变全打在仓外导出台副本，每发 `RESTORE_OK`）。

## 1. 四格判语

| 格 | 判语 | 凭据（全部本人现跑） |
|---|---|---|
| `AC#1` 射程成文 | **成立** | 名册 15 枚＝裁定枚举逐枚相同；全标量普查 170 枚 `Cf` 里边界只洗 15 枚，`extra` 空、`missing` 空；⛔ `unicode.Mn`（`git grep rc=1`）；判级轴那张表零字节 |
| `AC#2` 两把尺不许分家 | **成立** | 成对两发本人重放＝两发各红、红名册作差恰为两枚对方靶子；另加一发"两把尺同时失明"作反证；再加一把全标量属性尺（1,112,063 枚）零分歧 |
| `AC#3` 洗＝转义⛔ 删 | **成立** | 反形（删而不逃）那发多红 4 枚、逐字红句 `:411` 15 次；三处装饰断言逐枚具名见第 3 节 |
| `AC#4` 判定不许被名字改变 | **成立** | 正向半格常驻绿（真 `risk.NewRiskAssessor()`）；★**票点名要的那一发"攻它的"＝本人自己跑了**（改导出台副本的 `internal/risk/assessor.go`）⇒ 两枚用例红、逐字红句在读数件 `[C]` |

## 2. `AC#2` 那对成对突变的逐字重放（派单最要紧的一件）

台面＝`D:/tmp/w327v1/bench`（`git archive 60c328686` 八件套＋`docs/`＋手工三件；`.go` 747 枚；被裁字节与工作树逐字同哈希 `8050f22d…`）。
还原写在 `trap … EXIT`（＝finally），七发每发当场还原并自报 `RESTORE_OK`；还原后哈希逐字＝`8050f22d14203bba006cbfabcfb2a748081286d492691109c3a35a7a435ff8fd`。

| 发 | 摘掉的东西 | 突变件 sha256 | 四数（RUN／PASS_TOP／FAIL_TOP／SKIP） | 比终态多出的红 |
|---|---|---|---|---|
| 一 | 写入侧 `case isDisplayWashCf(r)` 那 7 行支（判据留着） | `e3bb5a1c…4598db` | 252／151／11／0 | `TestTheWriterBranchWashesCfEvenWhenTheFastPathAlreadyFired`（**唯一**）＋`TestEveryNamedCfRuneIsEscapedAndNeverDeleted`＋`TestABidiNameWalksTheRealBoundaryIntoAWashedRow`＋`TestTheBidiClaimInANameChangesNoVerdict` |
| 二 | 判据侧 `|| isDisplayWashCf(r)` 一枚 term（写入侧留着） | `579e02d5…58e5e7` | 252／151／11／0 | `TestNeedsDisplayWashNamesEveryWashedCfRune`（**唯一**）＋同名册那三枚 |

两发红名册 `comm` 作差＝**一多一少、不同枚用例、不同枚断言行**（发一红在 `:352`，发二红在 `:368`，逐字见读数件 `[B]`）。
★**"两把尺其实响在同一处"这一条我找过、找不到**：发一里判据靶子绿、发二里写入侧靶子绿，互不含对方。
再加两把：
- **发[m3]＝整枚谓词失明**（`isDisplayWashCf`=>`false`，两把尺同时去牙）⇒ 名册＝发一∪发二全部 5 枚，**没有一枚该响没响**（`TestTheRosterIsClosed…` 按设计仍绿，它只防加宽）。
- **补尺（常驻候选）**＝全标量属性尺：每枚标量问三件事——判据响不响、干净串里写入侧改不改、**判据被 ESC 强行打开后**写入侧逃不逃。终态 `rule-disagreements=0`；快路径早退会**天然掩盖**写入侧失明（干净串两侧都说"没变"），所以第三问不可省——省了就等于把 `AC#2` 读回"两边都加就完事了"。

四数与 327-r1 自报差 1 枚（它 12／我 11，PASS 150／151）＝**台面配方差、⛔ 判据差**：它那份导出台⛔ 带 `docs/`，`TestGitDimensionHasNoModelCallableTool` 读 `docs/PLAN.md`（`internal/panel/git_test.go:367`）必红；本人按现行配方（八件套＋`docs/`）跑，该枚绿 ⇒ 起手名册 7 枚。作差比的是同名册⇄同名册，结论⛔ 受影响。

## 3. 六枚新增用例逐枚点名（常驻⇄有牙⇄装饰）

常驻性：六枚**全常驻**——文件零 `//go:build`、零 `t.Skip`、零 `testing.Short`（尺＝读数件 `[G]` 同批），且六枚都在不带 `-v` 的普通 `go test` 里执行、终态逐枚 `--- PASS`。

| 用例 | 打的是哪格 | 有牙的证据（哪一发让它红） | 具名装饰部分 |
|---|---|---|---|
| `TestTheWriterBranchWashesCfEvenWhenTheFastPathAlreadyFired` | `AC#2` 写入侧靶子 | 只在发[m1]（＋[m3][m7]）红 | 末半句（`:357`–`:360`）的裁判是 `isDisplayWashCf` 自己＝**圈回**：发[m3] 里它 0 次命中而 `:352` 响 ⇒ 它抓不到"名单变窄"，只抓"输出还留着原始字形" |
| `TestNeedsDisplayWashNamesEveryWashedCfRune` | `AC#2` 判据侧靶子 | 只在发[m2]（＋[m3]）红，`:368` 逐枚 15 行 | 名单是测试件手抄的第二份（⇄ 产码那份）＝**这是它有效的原因**，⛔ 缺陷；但它只管"少认"，不管"多认" |
| `TestEveryNamedCfRuneIsEscapedAndNeverDeleted` | `AC#3` 主格 | 发[m1][m2][m3][m7] 全红，`:411` 各 15 次 | ★`:416`–`:418`"两半不许拼成第三枚名字"＝**恒不响**：八次运行 0 次命中，包括[m7] 那发正是它自称防的形——夹具 `photo`＋r＋`报告.png` 删掉 r 后前缀后缀仍在。真正在把关的是 `:411` 的逐枚 `want` 相等 |
| `TestTheRosterIsClosedAgainstTheRunesTheRulingLeftOut` | `AC#1` 排除钉 | 只在发[m4]（改成 `unicode.Cf` 整类）红：`:448` 逐字 `U+2061 is not on the ruled list but needsDisplayWash fired on it` | 单边钉（⛔ 抓"少洗"）；且⛔ 抓"加一枚没点名的"＝见第 5 节残余 R1 |
| `TestABidiNameWalksTheRealBoundaryIntoAWashedRow` | `AC#3`＋`AC#0` 边界真件 | 发[m1][m2][m3][m7] 红：`:477`/`:480`/`:486`/`:490` | 真 broker＋真两个构造器＝非装饰；`Spelling` 留干净、`U+202E` 放在用户没打的 `Canonical` 上＝与名册一致 |
| `TestTheBidiClaimInANameChangesNoVerdict` | `AC#4` | 发[m1]（`:534`）/[m2]/[m3]/[m7]/**[m5]** 全红 | 真 `risk.NewRiskAssessor()`、三名同判级；它同时把显示半格与判级半格钉在一枚用例里＝`AC#4` 要的"洗过的也要同判" |

"换反形会不会也红"逐枚问过：`escapeOf`（测试件自带的词汇抄本）与 `want` 相等那组**会**红；上面三处具名装饰⛔ 红。⛔ 一枚用例的存在本身被本人当成牙。

## 4. 能力必答：那 15 枚今天在生产路径上真被读到吗

**被产码读到＝是；被人眼读到＝未闭合（归口票 329）。** 尺＝`git grep` 非 `_test` 口径，逐枚具名（读数件 `[E]`）：
- `scrubForDisplay` 的产码调用者 **2 枚**＝`internal/panel/display_scrub_216.go:228`（`Workspace.Canonical`）与 `:234`（`Attachments[].Name`），两枚都在 `scrubComposerForDisplay` 体内。
- `scrubComposerForDisplay` 的产码调用者 **1 枚**＝`internal/panel/composer.go:145`（坐在 `NewSnapshot`＝同件 `:106` 体内）。
- `NewSnapshot` 的产码调用者 **2 枚**＝`internal/panel/pump.go:238`、`pump.go:309`；`NewComposerState` **1 枚**＝`pump.go:275`（第三实参＝`nil`）。
- 泵的产码装配 **4 枚**＝`cmd/wisp/run.go:702`／`:732`（`Out=rt.bookPanelSnapshot`）／`:742`／`:1063`；`publishPanelSnapshot` 体＝`rt.pump.Publish()`（`cmd/wisp/panel_pump.go:447`）⇒ **每次 publish 都过这道边界**，`U+202E` 那支 `case` 是真会执行的字节。
- 出口的远端＝`bookPanelSnapshot`（`panel_pump.go:360`）只存字段＋记一条摘要日志；唯一读它的 `lastPanelSnapshot`（`:424`）**产码读者 0 枚**（尺只剩定义行）；入向读快照那枚方法常量**不存在**（`MethodSnapshot` 尺 rc=1）；出向 `Eval` 产码 4 枚（`panel_control_windows.go:737/744/748/752`）逐枚都是 `window._rpc` 回执、⛔ 一处承载快照。
- 两枚被洗字段的**源**：`Workspace.Canonical` 源在（`run.go:705 Workspace: rt.workspaceView`，读数腿 `panel_pump.go:124`）；`Attachments[].Name` 源 0（`pump.go:275` 传 `nil`）＝与普查 ② 一致，本票⛔ 欠这半格（票 216 残余 ③ 同形）。
- 结论：`AC#1`–`AC#4` 四格的本体是**出向字节边界**，那半闭合；"走到屏幕上那一行"今天仍断在 `panel_pump.go:360`→`:424` 之间＝**票 329 在飞的活**，⛔ 记在 327 的欠账上。⛔ 拿 `go build` 的 rc=0 当"接上了"——上面那串调用点尺就是为这一条跑的。

## 5. 残余（具名，附归口）

- **R1｜名册基数⛔ 有一枚钉在产码上**：发[m6] 把**没被裁定的 `U+206B` 加进名单** ⇒ `RUN 252／PASS 155／FAIL 7`＝与终态**逐枚同名、多红 0 枚**，全部用例照绿。同一把补尺打在同一个副本上立刻红并点出 `U+206B`。⇒ 建议把 `zz_v1_exhaustive_test.go` 那两枚落常（需要一枚 `internal/panel/**` 写腿，本裁决腿⛔ 落它）。
- **R2｜`AC#1` 那句"具名裁要不要"还欠一纸**：`U+2061`–`U+2064`、`U+061C`、`U+180E`、`U+206A` 这 7 枚今天**照旧穿过**，且被钉成"排除"。产码与用例都自陈这是临时（注释原文：stay raw on purpose until a ruling names them）。裁定归编排者／机主，⛔ 本腿代裁。
- **R3｜`AC#4` 的归口更正**：327-r1 交回 ① 说那一发"射程＝`internal/risk/**`／`internal/perm/**`＝本腿禁区，请派有那面写权的腿"——**前提错了**：那一发是**仓外副本里的定向突变**，⛔ 需要任何跟踪字节的写权（本人刚跑，risk 件还原哈希逐字回 `9e73889b…`）。⇒ 建议：⛔ 为这一格另派写腿，改由本件结案；若仍想动 `internal/risk` 的**产码**（另件事），才需要新票。
- **R4｜装饰三处**（第 3 节具名）：`:416`–`:417` 恒不响、`:357`–`:360` 圈回、`:400`–`:401` 钉的是测试件自己那份抄本。⛔ 判据受损，⛔ 要改（改它们＝改既有断言，先停手上报那条规矩在这里生效）。

## 6. 本腿自证（派单"交件自证"三把尺）

1. **逐枚提交号＋逐枚 `git show --name-only`**：本腿落 4 笔（见下），每笔名册只含四处写面：`docs/evidence/s1/327-v1-verdict.md`、`.scratch/wisp/probes/327/v1/**`、票 327 末尾**一条** bullet、commit-msg 源件（也在 `probes/327/v1/`）。被裁料名册另见读数件 `[G]`：三枚提交越界 0 枚、`git diff --name-only 7c9d7ebe 60c328686`＝7 条路径全在写面内。
2. **判语件**＝`docs/evidence/s1/327-v1-verdict.md`；尺＝`wc -c`／`wc -l`／五枚空壳记号尺（尺面见派单"输出预算"那一行）＝0 命中（数值见交件报告）。
3. **成对两发**＝第 2 节逐字红句＋红名册＋两发不同处；**还原哈希**＝七发逐发 `restored_sha256=8050f22d…`（m5 另有 `restored_risk_sha256=9e73889b…`）；**起手⇄终态四数**＝起手名册 7 枚（本人台面无突变）⇄ 七发各 11／11／12／8／9／7／11；**调用点尺逐字输出**＝读数件 `[E]` 八把。
4. **本腿⛔ 跑过**：整包 `./cmd/wisp`、`go test ./...`、push、翻框、动 `internal|cmd` 跟踪字节、杀 `wisp.exe` PID 21892。

## 7. 我在票面／派单／证据件里查到的不实或过期（⛔ 替谁圆）

1. **"22 枚"**（`A1008` ⑤1 头一句、`327r1-dispatch.txt`、本派单同引）＝加法错，枚举清单逐枚求和＝**15**（方向 5＋4＝9，零宽/连接 4，`FEFF` 1，`00AD` 1）。编排者已在 `A1009` ③2 自纠；本人**独立复算**同样＝15，与裁定枚枚对齐（补尺 `extra/missing` 双空）。
2. **327-r1 证据件第 5 节第 1 条的归口前提**＝不实（见 R3）。
3. **327-r1 证据件第 4 节的名册枚数**（"8 枚起手红""两发各 12 枚红"）＝**台面配方过期**而非造假：现行配方（八件套＋`docs/`）下该包起手红名册＝7 枚、两发＝11 枚；差值恰为 `TestGitDimensionHasNoModelCallableTool` 那枚 `docs/` 依赖。它的"名册作差＝0"在两把尺下都成立。
4. **票面现量表第 4 行**"全树认得双向控制符＝命中 0 枚文件"＝今天重跑命中 2 枚文件（产码 3 行＋用例 14 行）＝**按设计过期**（那行本来就写"起手逐条复算，⛔ 信这里的号"）。
5. **派单引用的三枚行号**（`:123` 写入侧、`:214` 判据侧、`:164` 函数行）**今天逐枚对得上**；票面 AC 文字与实现无冲突。派单头"未勾 6／已勾 0"＝对（尺复跑）。
6. 本人**⛔ 发现**任何一处把"两把尺"写成同一处的假凭据；也**⛔ 发现**任何放宽既有断言的痕迹（测试件 `numstat`＝204/0，零删除行；产码件被删 13 行＝12 行注释＋1 行被替换的旧判据）。
