# 333-r1 · 显示边界改按 `Cf` 范畴判定 · 产码腿交件料

- 票＝`.scratch/wisp/issues/333-display-boundary-wash-cf-category-not-hand-roster.md`（AC#0–AC#8 九格，起手现跑 `grep -c '^[[:space:]]*- \[ \]'`＝**9**，⛔ 抄派单）
- 起手锚＝`ae710c49`（立票那一笔）；本腿提交＝`3f3dde46`(AC#2 产码) · `dbaad267`(改前双色) · `5f4f7945`(AC#3 换形) · `9574e794`(AC#4 新钉) · 本笔(证据＋牙)
- 写面逐笔 `git show --name-only` 只含：`internal/panel/display_scrub_216.go`、`internal/panel/display_scrub_327r2_test.go`、`internal/panel/display_scrub_333_333r1_test.go`(新建)、`.scratch/wisp/probes/333/r1/**`、本件、票尾一条 bullet。
  `internal/panel/display_scrub_216_test.go`＝**零字节**（尺＝`git diff --name-only ae710c49 HEAD` 名册里它不出现），三枚冻结件／`cmd/wisp/**`／`internal/risk/**`／`frontend/**`／`design/**` 零字节。
- 件：`wc -l`＝229 / 160 / 126（帽 160 内），`gofumpt -l`＝空，`go vet ./internal/panel/`＝rc 0。

## AC#1 起手三把尺（逐枚现跑，⛔ 抄本票，各带正控）
1. 名册枚数＝**15**。尺＝`sed -n '166,170p' internal/panel/display_scrub_216.go | grep -o '0x[0-9A-Fa-f]\{2,4\}' | wc -l`＝15；同尺对合成三枚行读 3（正控，尺不哑）；`grep -c 'case 0x'` 整文件读 **1**（＝派单点名的跨行 `case` 陷阱，实测复现）。
2. `unicode.Cf` 全集＝**170**，`GO_UNICODE_VERSION=17.0.0`。尺＝仓外纯 stdlib 探针 `D:\wisp-333r1-tmp\cfcheck\main.go`（⛔ 本仓文件，`go run` rc=0），逐标量 `unicode.Is(unicode.Cf,r)`；正控＝`CC_COUNT=65` 对上算式 `32+33`、`CTRL_CF_U202B=true`、`CTRL_CF_U0041=false`。⛔ 引 `332-a1` 的数。
   同发逐枚读数：排除名册八枚 **全 `cf=true`／`zs=false`**（含 `U+180E`＝本票射程声明第 3 条那枚表分歧，判据锚 Go 的表）；旧十五枚 **全在 `Cf` 内**（`FIFTEEN_ALL_IN_CF=15`）⇒ 按范畴只可能扩、不可能丢。
3. 被洗字段＝**2 枚**、生产调用点＝**1 处**，行号自量：`internal/panel/composer.go:145`（`Composer: scrubComposerForDisplay(composer)`，非测试唯一命中）→ `display_scrub_216.go:218`（`c.Workspace.Canonical`）／`:224`（`a.Name`）。正控＝`git grep -n scrubComposerForDisplay` 全仓另两命中是本文件自己的注释与函数定义。
   （注：`display_scrub_216.go` 被我改写后行号位移，`:218`/`:224` 是**交件侧**现量，票面 `:228`/`:234` 是锚点侧读数。）

## AC#2 范畴判定落成 ＋ 名册清零
尺＝`git diff ae710c49 -- internal/panel/display_scrub_216.go` 的**新增行**（日志 `logs/ac2-ruler.log`）：`added_lines=32`、`case 0x`＝0、`0x`＝0、`U+`＝0。正控＝被删侧 42 行里 `case 0x`＝1、`U+`＝15（同一把尺反向不哑）。
函数体现形＝`return unicode.Is(unicode.Cf, r)` 一枚；Cc 那一支（`unicode.IsControl`）在写入侧与判侧各一枚，**一字未动**（AC#8）。
★我刻意⛔ 改的三处过期措辞（改了会把 Cc/转义词汇的码位界尺带进新增行、与 AC#2 的尺对撞）：`:90` "(the ticket 327 Cf roster)"、`:92` "the invisible roster"、`:193` "the invisible formatting roster in isDisplayWashCf"。三处都是**散文**，不影响判定；账归编排者另裁（与票 327 残余 R4 同形）。

## AC#3 两枚既有钉：换形⛔ 打红、逐枚具名
| 锚点名（`ae710c49`） | 交件点名 | 变了什么 |
|---|---|---|
| `TestDisplayWashRosterCardinalityIsExactlyTheRuledFifteen` | `TestDisplayWashRosterIsExactlyTheCfCategory` | 期望侧从 15 枚字面清单（`ruledRoster327r2`，**已删、⛔ 留作常量**）换成按范畴收集；实际侧仍读产码谓词 `isDisplayWashCf` 本身；两臂＝extra（`:68`）／missing（`:72`）两枚**不同断言行**，双向差集皆空；名里⛔ 枚数词 |
| `TestDisplayWashRosterExcludesTheEightNamedRunes` | `TestDisplayWashRosterWashesTheEightRunesTheOldRulingExcluded` | 方向整批翻成"被洗＝true"：谓词、快路、精确转义形状、`ContainsRune` 点名"⛔ 删除" |
两枚都是**改名⛔ 删用例**——旧名在册、映射在上表，⛔ 悄悄消失。`excludedRoster327r2` 名册本体留在原位（`:22`），八枚逐枚读数＝`U+2061 FUNCTION APPLICATION`／`U+2062 INVISIBLE TIMES`／`U+2063 INVISIBLE SEPARATOR`／`U+2064 INVISIBLE PLUS`／`U+061C ARABIC LETTER MARK`／`U+180E MONGOLIAN VOWEL SEPARATOR`／`U+206A ACTOR WITHOUT`／`U+206B SPEAKER WITHOUT`，八枚现全走 `\uXXXX` 小写形（钉的是同包 `display_scrub_216_test.go:425` 那枚 `escapeOf` 复用了同一套词汇，界面仍一套转义格式）。
`nail 3`（同集钉）与票 327 残余 R4 三处装饰断言**未动**：尺＝`logs/ac3-nails-reshaped.log` 第 2 条，`^// Nail 3` 起至文件末与锚点逐字 `IDENTICAL`。

## AC#4 两侧同集：整族逐枚，⛔ 抽样冒充整族
新建 `internal/panel/display_scrub_333_333r1_test.go` 两枚钉（`logs/ac4-new-nails.log`，`-run TestDisplayWash`）：
- `TestDisplayWashFastPathAndWriterBothAnswerWashAcrossBothFamilies`：读数 `full Cf=170, full Cc=65, boundaries=22, failures=0`（`:82`）。三臂三枚不同断言行（`:64` 判据必须响／`:68` 写入侧必须真转义／`:72` 两臂不许分歧），另加一枚空射程正控（`Cf==0 || Cc==0` 直接红）。
- `TestDisplayWashEscapesAndNeverDeletesAcrossTheWholeCategory`：读数 `examined=170, product washes=170, failures=0`（`:120`）——`examined` 与产品洗集基数对撞，"这一臂跑满了整族"是钉出来的、不是我说出来的；同时钉住原位形状、两半可见字节不许接成新名字（AC#8 的"转义⛔ 删除"）、幂等。

## AC#5 成对两发（牙，落点不同函数／不同断言行）
台面＝仓外导出树 `D:\W333R1MUT`（`git archive 9574e794` 八件套＋`docs/`＋手工三件；`.go` 取出即数＝**757**）。字节同源＝`sha256(树内件) = sha256(git show 9574e794:…216.go) = sha256(工作树件) = dc86d19ce9fdf5…d267d552d0`（三向逐字同串）。
| 发 | 落点（内容锚，⛔ 绝对行号） | 突变字节数 | PRE | MUT | RESTORE | 还原核对 |
|---|---|---|---|---|---|---|
| M1 | `isDisplayWashCf` 体内 token `return unicode.Is(unicode.Cf, r)` → 锚点那 10 行名册 `switch` | diff 11 行 | `dc86d19c…` | `8c12e722…` | `dc86d19c…` | YES |
| M2 | `needsDisplayWash` 签名行后插 `return false`（快路短路） | diff 1 行 | `dc86d19c…` | `d71c1ab2…` | `dc86d19c…` | YES |
同一条 `-run` 名册（5 枚）三发：`M0-control 9/8/1/0`、`M1 9/5/4/0`、`M2 9/2/7/0`（`RUN/PASS_TOP/FAIL_TOP/SKIP`）。
红名册作差（`logs/ac5-pair-summary.log`）：M1 独有＝`TestDisplayWashRosterIsExactlyTheCfCategory`；M2 独有＝`TestTheWriterBranchWashesCfEvenWhenTheFastPathAlreadyFired`、`TestNeedsDisplayWashNamesEveryWashedCfRune`、`TestEveryNamedCfRuneIsEscapedAndNeverDeleted`、`TestDisplayWashJudgeAndWriterAgreeOnBothFamiliesAndBoundaries`；两发共有＝`TestDisplayWashRosterWashesTheEightRunesTheOldRulingExcluded`、AC#4 两枚新钉。**⇒ 两发红在不同用例、不同断言行，少任一发另一发的绿都不算凭据。**
必红的那两枚，逐字带文案：
- M1 → `display_scrub_327r2_test.go:72: a Cf-table rune is not washed by the product: [U+0600 … U+FFF9 U+FFFA U+FFFB] (product=15 category=170)` ＋ `:76: product roster cardinality=15, want the category's 170`（关系钉真读产码：把产码退回名册它就红＝327-r2 摘除对照的形状在 333 仍在）。
- M2 → `display_scrub_333_333r1_test.go:87: Cc U+0000: needsDisplayWash stayed silent, so scrubForDisplay never runs`（首枚；全族 488 条同形失败，截前 20 条打印，`:82` 读数行点名 `failures=488`）。
两发各自当场还原（脚本 `mut333r1.sh` 里 restore 在 test 之后无条件执行，⛔ 突变留仓；工作树与树的字节同源见上）。

## AC#6 台面与颜色（两把临时根台面成对，四数逐枚点名）
一条命令逐字不变＝`go test -v ./internal/panel/ -count=1 -timeout 30m`，只挪 `TMP`/`TEMP`/`GOTMPDIR`；台面＝共享工作树（`.git`/`docs`/`design` 齐，门禁四数同面跑）。
| 发 | 临时根 | RUN | PASS_TOP | FAIL_TOP | SKIP |
|---|---|---|---|---|---|
| pre-anchor-c | `C:\Users\swq\AppData\Local\Temp`（GOTMPDIR 未设） | 255 | 160 | 5 | 0 |
| pre-anchor-d | `D:\wisp-333r1-tmp\tmp` ＋ GOTMPDIR `…\gotmp`（先 `mkdir -p`） | 255 | 160 | 5 | 0 |
| post-c | C 卷同前 | 257 | 161 | 6 | 0 |
| post-d | `D:` 卷同前 | 257 | 161 | 6 | 0 |
四数未混台面、未相减。★成对结果＝本包**没有** A1013 那一族差一枚红：两枚台面的名行序列剥耗时后逐行全等（`420` 名行 ⇄ `420`，差异名行 `0`，尺＝`logs/pre-anchor-surface-diff.log`；两枚 .log 与 .sha256 已入库）。
红名册作差（改前 5 ⇄ 交件 6）＝**只多一枚**：`TestTheRosterIsClosedAgainstTheRunesTheRulingLeftOut`。改前五枚（`TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestStreamLogFloodBelowKeyBoundIsNotBounded35r8`／`TestStreamLogDroppedNamingLedgerIsNotBounded35r8`／`TestComposerContractTypesMatchFrontend`／`TestC21DesignTokensFourWayAgree`）交件逐字同在，⛔ 我碰过；算式（⛔ 混台面相减，只做同一枚台面内的名数对账）＝顶层名数 165→167（我新增两枚用例，两枚都绿：见 AC#4 读数），PASS_TOP 160→161＝+2 新绿 −1 枚由绿转红，FAIL_TOP 5→6＝正是那枚转红的上报件。⛔ 整包 `./cmd/wisp`、⛔ `go test ./...`、⛔ 动 PID 21892。

## ★上报：一枚设计外红（未定义即停，本腿⛔ 自行放宽）
`internal/panel/display_scrub_216_test.go:444 TestTheRosterIsClosedAgainstTheRunesTheRulingLeftOut`＝327-r1 当年按"⛔ 在裁定内"把七枚原样穿过钉成的用例（票 327 尾 bullet 自己写了"纳入与否请拍一纸"）。扩集后它红 14 句（`:448`＋`:451` 各 7 枚：`U+2061 U+2062 U+2063 U+2064 U+061C U+180E U+206A`）。
它**不在**票 333 AC#3 与派单划定的"设计内两枚"里，又**在本腿写面之外**（写面只有 216.go／327r2_test.go／新建 333 件／probes／证据／票尾一条）⇒ 按 D22 闸门③ 停手上报，等一纸裁定：它应与 327r2 那枚排除钉同形翻向（七枚翻"被洗＝true"），还是随 333 的裁定作废。本腿**未**放宽、**未**删、**未** `t.Skip`，交件终态这一枚就是红的（AC#6 那 6 里有它）。
`U+206B` 不在它那七枚里（327-r1 那批漏了它，正是 R1 的由来），所以红的枚数差是 7 不是 8——这句是为了让裁决者⛔ 误读成"少了一枚牙"。

## AC#7 生产调用者必答（能力类，⛔ 用例绿背书接线真）
扩集后那两枚字段真被产码读到，全链逐枚 `file:line`（交件侧现量）：
`cmd/wisp/run.go:714`（`rt.pump = panel.NewSnapshotPump(...)`）→ `cmd/wisp/panel_pump.go:451`（`rt.pump.Publish()`）→ `internal/panel/pump.go:354`（`Publish`）→ `internal/panel/pump.go:309`（`snap := NewSnapshot(cards, results, composer, now())`）→ `internal/panel/composer.go:106`（`func NewSnapshot`）→ `internal/panel/composer.go:145`（`Composer: scrubComposerForDisplay(composer)`，全仓非测试唯一调用点）→ `internal/panel/display_scrub_216.go:218`（`c.Workspace.Canonical = scrubForDisplay(...)`）／`:224`（`a.Name = scrubForDisplay(...)`）→ `internal/panel/pump.go:359`（`json.Marshal(snap)`）→ 字节里 `canonical`＝`composer.go:231` 的 `json:"canonical"`、`name`＝`attachments.go:86` 的 `json:"name"`。
★到人的眼睛那一跳**仍不成立**，但我自己的尺读到的形状与票面那句话不同，据实写：
- 页面**渲染器存在**（⛔ 0 命中）：`frontend/src/App.tsx:243`／`:303`、`frontend/src/components/composer.tsx:150`／`:152` 逐枚读 `composer.workspace.canonical`（尺＝`git grep -n -i 'canonical' -- 'frontend/src/**'`＝8 命中，正控＝同尺 `composer` 命中 81）。
- 页面**接收器 0 命中**：`onmessage`＝0、`chrome.webview`＝0、`hostObject`＝1（只有 `lib/panel.ts:146` 那句注释说它是"Installed by WebView2's AddHostObjectToScript / postMessage pipe"）、`addEventListener` 全部命中是 keydown/resize/pointer，没有一枚是 `"message"`；`frontend/src/main.tsx:97` 挂的是 `<App />`（⛔ `snapshot` 入参）⇒ 走 `App.tsx:91` 的默认 `EMPTY`。
⇒ 结论＝包已经带洗好的字节出门（票 329 甲形把字节推给宿主，`panel_pump.go:451` 后一行的 `pushPanelSnapshotToHost()`），页面那一侧还**没有人把推来的包接进 `<App>`**；屏幕上今天没有改，本票⛔ 写成已改（AC#7 那一半仍归票 329）。

## AC#8 ⛔ 五件事自查（尺在括号里）
⛔ 动 R4 三处装饰断言（`:416-417`/`:357-360`/`:400-401`：红句名册只出现 `:448`/`:451`，`logs/ac3-nails-reshaped.log` 第 6 条）；⛔ 放宽既有断言／⛔ 删用例／`t.Skip`＝0（尺＝`git diff ae710c49 HEAD -- internal/panel/ | grep -c '^+.*t.Skip'`＝0；`len(cfRoster) != 15` 那枚照原样在，账归编排者）；⛔ 碰 `Cc` 那一支（`unicode.IsControl` 两处一字未动）；⛔ 发明删除形（依据＝`display_scrub_216.go:74` 那句 VISIBLE ESCAPE, not deletion，交件侧现量行号；AC#4 第二枚钉整族逐枚验过不接成新产品名）；⛔ 指名一枚不存在的件（本件点名的日志／脚本／用例名全部在盘上，`find .scratch/wisp/probes/333 -type f` 可数）。
占位符尺（分词形四类词，逐件跑，日志 `logs/placeholders.log`）：三枚代码件命中 **0**／本件命中 **1**，那一枚命中就是写下这把尺的这一行自己（自指，log 里逐行点名行号），⛔ 别的占位待补。

## 本腿没做的事（⛔ 冒充欠账清零）
- 没翻任何 AC 框、没改票名／`Status:`、没加 `-done`、没 push。
- 没裁那枚设计外红，只上报（AC#6 那格与这里同一枚）。
- `docs/evidence/s1/333-r1-display-cf-category.md` 之外的 R-333-1（`Mn`／`Zs` 另族）与 R-333-2（票 327 的 AC#0／AC#7）⛔ 本票欠账，一字未动。
