# 票 333 · 裁决腿 333-v1 · 九格判语

- 裁决者≠实现者（`SPEC-12` 第 4.3 节第 1/3 条）。实现腿＝`333-r1`，起手锚＝`ae710c49`，交件锚＝`859026d2`。
- 本腿⛔ 翻一枚框、⛔ 改判据文字／票名／`Status:`、⛔ 动 `internal/**` 一字、⛔ push；产出＝判语＋读数。
- 台面：三枚**仓外**导出台树（`mktemp -d` on D:，`dir /x D:` 证 8.3 别名在册）——
  `A=D:/W333V1-SuqCO8`（`git archive 859026d2`，`.go`＝757）／`P=D:/W333V2-hLzM9f`（`ae710c49`，753）／
  `M=D:/W333V3-b3pHjd`（859026d2 的第二棵，只当突变面）。凡"交件侧"读 A、"改前侧"读 P，⛔ 取共享工作树。
- ★`internal/panel/display_scrub_216_test.go:444` 那枚旧钉子当前状态：10-17 11:4x 该文件正被 `335-r1` 改写，
  **本腿读数取自 `859026d2` 树**（`git status --porcelain -- internal/panel` 本腿起手与收尾各量一次＝皆空）。
- 全部原始读数：`.scratch/wisp/probes/333/v1/logs/`（`readings-final.log` 是决定性那一份，脚本 `build-tree.sh`／
  `ac2-ruler.sh`／`ac1-ac7-rulers.sh`／`ac0-ac8-rulers.sh`／`colors.sh`／`mut333v1.sh` 可逐字重跑）。
- ★本腿自纠一枚数：本件起笔自测 15480 字节，随后把标题行那枚节号改成"第 4.3 节"，字节数随之变（现值逐字见
  `logs/final-check.log` 的 ruler 2，⛔ 在本件里钉死字节数＝那行自己也在分母里）。票尾那条 bullet 写的 15480 是改前旧值，
  两数⛔ 相减、本条具名留下更正；票面一字没动过，框一枚没翻（现跑未勾仍 9＝第 19–27 行）。

## AC#0 交件三把尺 ＝ **成立**
六笔逐枚 `git show --name-only` 本腿独立再走一遍：`3f3dde46`／`dbaad267`／`5f4f7945`／`9574e794`／`04f2e847`／`859026d2`
＝35 枚唯一路径，全部落在四枚允许形状里，**越格 0**（尺＝逐枚路径 `grep -v -E` 四条前缀，命中计数 0）。
`internal/panel/display_scrub_216_test.go` 在 `ae710c49..859026d2` 的 `--stat` **不打字**＝零字节，与腿自述一致。
票尾 bullet `numstat`＝`2 0`（一条 bullet ＋ 它自己那行前的空行，⛔ 第二条 bullet）。
五词空尺本腿自跑：三枚代码件 0／证据件 1（第 73 行＝写下那把尺的那一句自己，腿已具名披露）。件非 0 字节（78 行／14238 字节）。

## AC#1 起手三把尺自己现跑 ＝ **成立**（三把都带了正控）
① 名册枚数＝**15**（本腿尺＝`awk` 只取函数体＋`grep -o '0x..' | sort -u | wc -l`，树 P；跨行 `case` 陷阱复现＝整文件
   `grep -c 'case 0x'` 读 **1**，正控＝同一 awk 尺打在合成三枚行读 3，腿也带了这条）。② `unicode.Cf`＝**170**，
   `GO_UNICODE_VERSION=17.0.0`——本腿自己另建仓外纯 stdlib 探针（`D:/wisp-333v1-tmp/cfcheck`，`go run` rc=0、stderr 空），
   ⛔ 抄票面也⛔ 抄 `332-a1`；正控 `CC_COUNT=65`、`PC_U0041_cf=false`、八枚逐枚 `cf=true/zs=false/cc=false`、
   旧十五枚 `FIFTEEN_ALL_IN_CF=15 of 15`（⇒ 按范畴只可能扩）。③ 被洗字段＝2、生产调用点＝1（`composer.go:145` 非测试唯一
   调用，树 A 现量 `display_scrub_216.go:218`/`:224`；票面 `:228`/`:234` 本腿在树 P 逐行 sed 复认＝**锚点侧行号**，腿的据实更正成立）。

## AC#2 范畴判定落成 ＝ **成立**
尺＝`git diff ae710c49 859026d2 -- internal/panel/display_scrub_216.go` 的新增行：本腿读 **added=31 / deleted=41**
（腿报 32/42＝它把 `+++`/`---` 头部各算了一行，两对⛔ 相减、只作口径具名）；新增侧 `case 0x`＝0、`0x`＝0、`U+`＝0、
十六进制字面量枚数＝0；**正控＝同一把尺打在被删侧 1／5／15**（尺不哑）。
函数体逐字＝`return unicode.Is(unicode.Cf, r)` 一枚。第二枚来源查法＝非测试产码全量 `grep -rn 'unicode.Cf'` →
只有 `:164` 判定与 `:141` 一句散文；写入侧 `:124`、快路 `:204` 都问同一枚符号，⛔ 表、⛔ 常量、⛔ 白名单。

## AC#3 关系钉的恒真性 ＝ **成立**（本票最高危那格，实测四发）
必答题 (i)：**如果实际侧直接抄期望侧、不再经过产品函数，这枚钉子不还响。** 本腿自己构造了那一发（树 M、内容锚、逐发还原）：
- `M0` pristine 未突变＝`RUN=6 PASS_TOP=5 FAIL_TOP=1 SKIP=0`，五枚 333 钉全绿（那 1 枚红＝`:444` 设计外上报件，见 AC#6）。
- `M1` 产码谓词逐字换回锚点那 15 枚 `switch`（`MUT 8c12e7223cf148d5`，⛔ 动测试）⇒ 关系钉**红**：
  `display_scrub_327r2_test.go:66: nail 1: product-washed=15, unicode.Cf=170` ＋ `:72: a Cf-table rune is not washed by the
  product: [...] (product=15 category=170)` ＋ `:76: product roster cardinality=15, want the category's 170`。
- `M3` **同一枚产码突变**＋只把测试实际侧换成期望侧那枚表达式（`washSetOver(isDisplayWashCf)` → `washSetOver(func(r rune) bool
  { return unicode.Is(unicode.Cf, r) })`，`MUT 8c12e7223cf148d5 / 59fb54bc4e2b118f`）⇒ 关系钉**绿**：
  `:66: nail 1: product-washed=170, unicode.Cf=170` ＋ `--- PASS`。`M4` 反形单发（产码正确）也绿。
⇒ 反形与正形在 M1/M3 这一对上**读数不同**（一红一绿），按本仓定式"判据换成反形仍全绿＝不敏感"，本钉**不是**恒真式：
它的牙由 `:53` 那一行"实际侧读产品"承载，本腿逐字读过并把它打红过一次。
本腿同时把话说清：产码体与期望侧今天同值 ⇒ 这枚钉**今天**为真是定义式的，它买的是两件事——产码漂移（M1 已证会红）与
Go／Unicode 表漂移（`:66` 那行把基数钉在仪器读数上），它⛔ 检验"该不该按范畴"这半句（那属票面裁定，属契约不属测试）。
排除钉翻向逐枚具名（八枚在册，`unicode.Is(unicode.Cf,r)` 那枚守卫留在原地＝表一动它自己先红）；两枚旧名映射由本腿用
颜色名行序列独立证明是**改名⛔ 删用例**：pre 侧 `...CardinalityIsExactlyTheRuledFifteen`／`...ExcludesTheEightNamedRunes`
→ 交件侧 `...IsExactlyTheCfCategory`／`...WashesTheEightRunesTheOldRulingExcluded`，皆 PASS。

## AC#4 两侧同集不许破 ＝ **成立**
"整族逐枚"是真整族：本腿独立取全集枚数（⛔ 抄票面 170、⛔ 抄 `332-a1`）＝`CF_COUNT=170`；`Cc` 那一支另跑一枚探针
（`D:/wisp-333v1-tmp/ctlexp`）＝`CC_CATEGORY=65`、`ISCONTROL_TOTAL=65`、`ISCONTROL_AND_NOT_CF=65`、`ISCONTROL_BUT_NEITHER_Cc_NOR_Cf=0`
⇒ go1.27.1 的 `unicode.IsControl` 恰等于 `Cc`，测试里 `full Cc=65` 那支计的就是 Cc，没夹带 Zl／Zp（`U+2028/2029 iscontrol=false`）。
测试自己那两行读数（本腿在 M0 复跑）＝`display_scrub_333_333r1_test.go:82: full Cf=170, full Cc=65, boundaries=22, failures=0`、
`:120: examined=170, product washes=170, failures=0` ⇒ 基数与产品洗集对撞＝整族是钉出来的。
`U+180E` **在**那 22 枚具名边界里（`boundary333r1` 逐枚数＝22，`0x180E` 在册），而边界那支的 `want` 取自
`unicode.Is(unicode.Cf,r) || unicode.IsControl(r)`＝**跟着 Go 的表走**，与票面射程声明第 3 条同形；本格凭据就写成这一形，
⛔ 替票面加强成"Unicode 官方类别"（Go 表今天对 `U+180E` 回 `cf=true/zs=false`，本腿探针逐枚复认）。

## AC#5 成对两发（牙）＝ **成立**（本腿自己重放，⛔ 抄腿的哈希）
落点不同函数／不同断言行：M1＝`isDisplayWashCf` 体内 token（→ 关系钉 `:72`/`:76`）；M2＝`func needsDisplayWash(...) {` 之后插
`return false`（`MUT 7626aa28a2c5646e`）→ 同集钉红：`:87: Cc U+0000: needsDisplayWash stayed silent, so scrubForDisplay never runs`、
`:82: ... failures=488`，另 `:125: U+00AD: "photo\u00ad报告.png" came back untouched`。
红名册 `comm` 作差＝**M1 独有** `TestDisplayWashRosterIsExactlyTheCfCategory`／**M2 独有**
`TestDisplayWashJudgeAndWriterAgreeOnBothFamiliesAndBoundaries`／三枚共有（两枚 AC#4 钉＋排除钉）⇒ 两发不互替。
pristine 那一发（M0）绿，"红"不是台面坏。逐发当场还原：四发各自 `RESTORE` 行与 PRE 逐字同串
（`dc86d19ce9fdf52b`／`ddf652a42a0e5f70`），末了 `diff -r A/internal/panel M/internal/panel` **零输出**。
★本腿抓到一条腿没写的形状：M2 下 nail 3 的同集臂读 `disagreements=0`（两臂一起哑＝同集不红），红的是它自己的正控 `:158`；
真正把短路打红的是 333 新文件的方向臂 `:87` ⇒ AC#4 的"同集不许破"由**新钉**承载，⛔ 把 nail 3 当那把尺。

## AC#6 两把台面 ＝ **成立**（射程分开写清）
本腿四发（同一条 `go test -v ./internal/panel/ -count=1 -timeout 30m`，只挪临时根所在卷，⛔ 混台面、⛔ 相减）：
改前 C 卷 `255/158/7/0` ⇄ D 卷 `255/158/7/0`；交件 C 卷 `257/159/8/0` ⇄ D 卷 `257/159/8/0`（rc 皆 1，包本就带红）。
同卷成对＝改前⇄改后各按同一配方、同一棵树集。名行序列剥耗时后每枚锚的 C⇄D **只差汇总行那列耗时** ⇒
★"无差"这一条的射程＝**`internal/panel` ＋ 本配方**；`A1013` 那族换卷差一枚红是在**别的面**（`internal/risk`、`cmd/wisp`）量的，
两件事⛔ 混写。本腿四数比腿的五枚红多两枚，具名＝归档树工件、⛔ 计入判据：`git_test.go:600 kind="not_a_repo"`（树无 `.git`）、
`composer_test.go:873`＋`:902 "0 painted states"`（`git archive` 出的 `frontend/fixtures/composer-states.html` 与工作树副本字节不同
＝16649⇄16641；本腿⛔ 把机制归到底，但它两枚锚同现 ⇒ 作差里抵消）。
红名册作差＝**只多一枚、一枚也没少**：多的那枚逐字＝`--- FAIL: TestTheRosterIsClosedAgainstTheRunesTheRulingLeftOut`，
红句 `display_scrub_216_test.go:448: U+2061 is not on the ruled list but needsDisplayWash fired on it` ／
`:451: U+2061 is not on the ruled list but the wash changed it: "photo\u2061.png" -> "photo\\u2061.png"`（七枚 ×2 行）。
本格⛔ 因它打红：**本腿确认那枚红不是 333 产码的缺陷**——产码按甲-全裁形执行（AC#2 已证），而那枚钉子断言的前提是
"这七枚⛔ 在裁定内、应原样穿过"，那前提属票 327 的旧名册，已被票 333 的裁定取代；它住在 `display_scrub_216_test.go`
＝票 335 的写面里。归口＝票 335／编排者一纸裁定（本腿⛔ 处方）。

## AC#7 生产调用者必答 ＝ **成立**，链由本腿逐枚 sed 复认并补两跳
11 枚 `file:line` 本腿在树 A 逐枚读回（`run.go:714`／`panel_pump.go:451`／`pump.go:354`／`pump.go:309`／`composer.go:106`／
`composer.go:145`／`display_scrub_216.go:218`／`:224`／`pump.go:359`／`composer.go:231 json:"canonical"`／`attachments.go:86 json:"name"`）
全部逐字对上。**接线的尺＝调用点枚数**（⛔ 拿 `go build`/`vet` rc=0 背书）：`scrubComposerForDisplay` 非测试命中＝3
＝调用点 1 枚（`composer.go:145`）＋定义＋散文各一枚。本腿另补腿那两跳之间缺的绑定与驱动点：
`panel_pump.go:447 func (rt *agentRuntime) publishPanelSnapshot()`（:451 `Publish()`、:463 `pushPanelSnapshotToHost()`）被以方法值绑在
`run.go:754` 与 `run.go:1075`，实际调用在 `run.go:1375`／`:1464`／`:1493`——腿那句 `run.go:714 → panel_pump.go:451` 之间的箭头⛔ 是一枚调用边。
"到人的眼睛那一跳仍不成立"＝成立，但票面那句**字面不准**：本腿尺＝`grep -rn -i canonical frontend/src`＝8 命中（正控 `composer`＝81），
渲染器在（`App.tsx:243`/`:303`、`components/composer.tsx:150`/`:152`）；真正 0 命中的是推包入口（`onmessage`＝0、`chrome.webview`＝0、
`addEventListener("message"`＝0；`addEventListener` 14 枚全为 keydown/resize/pointer），`main.tsx:97` 挂 `<App />` 无入参 ⇒
走 `App.tsx:91` 默认 `EMPTY`。字节已出门、页面没人接＝票 329 那一半，本票⛔ 写成屏幕上已改。

## AC#8 ⛔ 五件事 ＝ **成立**
`git diff ae710c49 859026d2 -- internal/panel/`：新增 `t.Skip`＝**0**；删除 `func Test`＝2、新增＝4（那 2 枚＝改名，映射见 AC#3）；
断言行 删 8／增 15。★本票最硬一问（改名＋期望侧翻向是更严还是放宽）＝**更严**，判据两条：
(1) 旧钉断言"这八枚留在原地不被洗"（`isDisplayWashCf=false`、快路沉默、`scrubForDisplay` 逐字返回），新钉要求
`isDisplayWashCf=true`、快路必开、**且**结果必须等于 `escapeOf(r)` 那枚精确形状、**且** `ContainsRune` 证⛔ 删除——对产码的
要求只增不减；(2) 方向轴在本仓是"洗＝转义显示、不洗＝放行"，把 15→170 是让更多字节变成可见转义，⛔ 一步放行。
判级轴另在 `internal/risk/taintmatch.go`，本票对它**零字节**（`git diff --name-only ae710c49 859026d2 -- internal/risk`＝0 文件），
`thresholds.go`／golden／`allowlist.txt`／`go.mod`／`go.sum`／`frontend/**`／`design/**`／`docs/specs/**`／`PLAN.md` 同尺＝0 文件。
Cc 那一支（`unicode.IsControl` 两处，`:121`/`:204`）一字未动；`display_scrub_216_test.go:400 len(cfRoster) != 15` 两枚锚都在原位（⛔ 动）。
转义⛔ 删除的逐字依据＝树 A `:74` 那句 VISIBLE ESCAPE, not deletion。
三处过期散文（树 A `:90`/`:92`/`:193`）留原地＝**⛔ 算 333 的债**，本腿**驳回不了、并加分一条理由**：`:92` 那一行本身就写着
`U+202E` 与 `U+00AD`，改它＝把码位名字送进 AC#2 那把尺的**新增行**、与票面那把尺正面撞（尺会读出 `U+`≥1 而判"没按裁形做"）。
归口＝票 335 的整行注释面，与 R4 同形；要修得先修尺（见下面更正建议第 1 条）。

## 强制回答三句
**(a) 九格里哪几格不成立、缺哪一发读数**：**没有一格不成立**——`AC#0`…`AC#8` 九格全部**成立**，且每格都由本腿自己的
读数支撑（⛔ 一条抄腿）：AC#3 的反形攻击四发、AC#5 的成对两发＋pristine、AC#6 的四发四数、AC#1/AC#4 的两枚仓外探针。
两处**缺东西**但⛔ 打红本格，具名上交：① AC#6 那枚"多的红"的处方归票 335（本腿只归因、⛔ 判它）；② AC#7 票面那句
"页面接收器 0 命中"字面不准，更正建议见 (c)，框照旧不勾由你裁。
**(b) 本腿为收件动过产码字节吗**：**零**。本腿只写 `docs/evidence/s1/333-v1-verdict.md`、`.scratch/wisp/probes/333/v1/**`
与票 333 尾一条 bullet；`internal/**` 的任何字节都只在仓外树里被读过或（在树 M 里）被临时突变并逐发还原（末次
`diff -r` 与 PRE／A 三方同串）。`git status --porcelain -- internal/panel` 由**编排者**在共享工作树现量盖章，本腿⛔ 自证。
**(c) 判据文字更正建议（交你裁，本腿⛔ 动票面一字）**：
1. AC#2 那把尺建议限定射程为"⛔ 函数体内出现字面码位或 `case 0x`"，并写明"注释行⛔ 计入"——现形（新增行 `U+` 命中 0）会把
   修散文与守裁形变成互斥，腿为了守尺只能留下三处过期话。
2. AC#7 末句"页面接收器 0 命中"建议改成"页面**推包入口** 0 命中（`onmessage`／`chrome.webview`／`addEventListener("message"`），
   渲染器 4 枚在位"——现文字面会被读成"页面里什么都没有"，与 `frontend/src` 的实测形状相反。
3. AC#0 的"分词尺 0 命中"建议补一句"写下那把尺的自己那一行除外"，否则交件只能靠披露过关（腿正是这么办的）。
4. AC#6 建议加一句"改前⇄交件的四发必须在**同一棵树集**上跑"——腿那四发跑在共享工作树、本腿那四发跑在归档树，
   两把数天然差两枚红（`.git` 与归档后的 html 夹具）；不写清，后来的腿会把口径差读成判据差。
