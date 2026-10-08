# 33-r12（票 33 那一格）名册＝第 1 笔交付（先取证、未动字）

- 腿：`33-r12`（文本归位腿，零代码／零测试／零票面框）｜工单＝票 33 `.scratch/wisp/issues/33-panel-host-c27.md`
- 起手 HEAD 现量：`bb871191`（派单给的是 `23ade0da`／台账笔 `bb871191`；我这一把 `git log --oneline -1` 真身＝`bb871191`，⛔ 不采派单转述）
- 起手 `probes/33/` 现量：`a1 a2 a3 h1 n1 p1 r1 r2 r3 r4 r5 r6 r7 r8 r8b r9 r10 r11 v2 v4 winc1`＝21 枚，**`r12` 目录当时不存在**（＝空，与派单预期一致）
- 起手写面核查（`git status --porcelain`）：满树脏项属共享工作树别人格的活，与本腿无关；本腿只写 `.md`。
- 全量尺（第 1 发取证）：`grep -rn 'firstRoundTripLocked\|serveNotBuiltNoticeLocked' --include='*.md' .`
  - rc=0，读数落 `grep-raw.txt`＝**102 行／39,234 字节／27 枚文件**（逐枚计数＝`grep-byfile.txt`），最长行 1,344 字符（`awk` 量）。
  - 名册逐行都按**整行**读（`Read` 读 `grep-raw.txt` 全 102 行，非 `cut -c1-150`）。
- 代码侧现状（只读尺，未 build／未 test）：`grep -rl --include='*.go'` 旧名只剩 6 枚文件＝`.scratch/wisp/probes/33/{p1/q3,p1/receipt,r11/mutations×3}` 的不可编译拷贝 ＋ `cmd/wisp/panel_locked_naming_33r11_windows_test.go:360` 的**正控夹具** `roundTripHost.firstRoundTripLocked`（该名字是那枚钉的输入，改了钉就哑）。产码与既有测试调用点＝新名（`cmd/wisp/panel_host_windows.go` 内 `firstRoundTrip` 命中 5）。⇒ 与 33-r11 报的"代码里已全改完"相符。
- 判据（唯一标准）：这句在描述"代码现在长什么样"（现在时）还是"当时发生了什么"（历史）。历史的绝不抹名。
- 写面（派单具名放开）＝`.scratch/wisp/issues/**` ＋ `docs/evidence/**` ＋ 本探针目录；`docs/reports/**`（含台账 `pending-and-issues.md`）与别人格的 `probes/**`＝禁区。

处置图例：`改`＝本腿第 2 笔换字；`历史`＝不许动；`禁区`＝写面外或判据文字；`不可纯换字`＝换名会让句子本身变成新的假话。

---

## 1. 写面内：`.scratch/wisp/issues/33-panel-host-c27.md`（5 行）

| 行 | 那句在说什么 | 判 | 处置 |
|---|---|---|---|
| 39 | AC#13 框正文：10-01 11:1x 的现读时序（`:221` 供页→`:250`→`:227` 调 `firstRoundTripLocked`→`:372-375` 覆盖），尾部 10-08 11:1x 就地收窄注自己写着"那四个行号是 10-01 的现读、已漂" | 历史（带时刻的现读＋同格已注过期）＋判据文字（AC 框） | 不改（双重理由：硬约束禁改判据文字；且它是被保留的读数） |
| 321 | 裁定 P2：判据形状就此定死，"⇒ `firstRoundTripLocked` 的射程写死＝页→Go 到达" | 历史（裁定原文），且与台账 `pending-and-issues.md:10292` 逐字同源 | 不改（台账不可动，单方面改票面会造出票面↔台账的假分歧） |
| 348 | 33-r10 进度追加 ①："现量（尺＝grep…，本文件**当时** 861 行）：`:425` 先调 `firstRoundTripLocked` → 它在 `:757-758` 发探测页…" | 历史（那句自带的尺与"当时"，行号今天已漂到 `:448`） | 不改（换名留旧行号＝半新半旧的假读数） |
| 358 | 同一节 ② 里对新台件第一问的机制描述："探测页由 `firstRoundTripLocked` 对着一次性接收器当场跑出来再取，不是抄进台件的字符串"——断言的是取法（运行时取），名字只是指针，无日期、无行号 | 现在时描述代码 | **改**：`firstRoundTripLocked` → `firstRoundTrip` |
| 385 | "这笔账是什么"：`serveNotBuiltNoticeLocked`（`:460`）与 `firstRoundTripLocked`（`:840` 一带）**带 `Locked` 后缀却自己取锁** | 现在时描述代码，但**不可纯换字**：谓词就是"这枚名字带 `Locked` 后缀"，换成如实名后句子自相矛盾（"带 `Locked` 后缀"变假）；如实更正只能加一行时间戳注记＝本腿禁加行，且该段措辞编排者已在台账 `A713` 自留（"这枚票面措辞要我改与①同批"） | 不改（顶回派单，见 §5） |

同节附带：`:391` 那句"现存那 2 枚进具名豁免名册直到 ⓐ 落地"不含旧方法名，本腿射程外（名字尺抓不到），但它同样是"与 ⓐ 同批执行＝没有可落瞬间"的措辞，编排者已自留——本腿一字未碰，具名告知。

## 2. 写面内：`docs/evidence/**`（6 行／4 枚文件）

| 文件:行 | 那句在说什么 | 判 | 处置 |
|---|---|---|---|
| `33-panel-host-c27-r4.md:251` | 33-r4 仪器腿的"本程现量"：`bringUp` 的 `:221 serveEntry()` 之后 `:227 firstRoundTripLocked(...)`，后者 `:374-375` 又 SetHtml | 历史（10-01 该腿的读数＋当时行号） | 不改 |
| `33-panel-host-c27-r5.md:56` | 表行"改前（行号现取于 `7a02b121`）"：那串时序 ⇒ 每次冷启动最终显示探测页 | 历史（引用某枚 commit 的内容＋锚） | 不改 |
| `33-panel-host-c27-r5.md:57` | 表行"改后"：33-r5 那一笔另加 `serveNotBuiltNoticeLocked()`、`firstRoundTripLocked` 的射程在注释里写死 | 历史（记那枚 commit 造了什么） | 不改 |
| `33-panel-host-c27-v1.md:83` | v1 验收腿"顺路核 AGENTS §1.2"的读数：超时那一形＝0，附 `firstRoundTripLocked` 的 `t0.Add(5*time.Second)` 与当时注释行号 `:168-173` | 历史（该腿那把尺的原读数） | 不改 |
| `33-panel-host-c27-v1.md:89` | §A#28 产码形状：`:221`→`:227`→`:372-375`，"`firstRoundTripLocked` 的调用点全仓唯一＝`:227`（grep 现量）" | 历史（带时刻、带当时行号的验收读数） | 不改 |
| `33-panel-host-c27-v2.md:34` | AC#13 判语行的证据半："次序已重排——`panel_host_windows.go:308` 先 `firstRoundTripLocked`、`:310` 后 `serveEntry`"，起手锚 `3216ba6d`（10-01 16:47） | 历史（那一波终裁的读数，且行号属当时那个树） | 不改 |

⇒ **`docs/evidence/**` 六行全判历史，一枚未动。** 这六行是被后面那些引用（台账、票面 `:39` 收窄注）逐字引着的验收记录，抹名＝毁掉可核性。

## 3. 禁区：`docs/reports/pending-and-issues.md`（6 行）——既有行一字不许动（连追加都不必）

| 行 | 那句在说什么 | 判 | 处置 |
|---|---|---|---|
| 10204 | 10-01"本轮新开的一格（票 33 AC#13）"：那串行号＋`firstRoundTripLocked` | 历史＋禁区 | 不改 |
| 10292 | 裁定 P2 原文（票面 `:321` 的同源那份） | 历史＋禁区 | 不改 |
| 12883 | 逐字引用 `probes/msg-a653.md:15` 的读数（含 `serveNotBuiltNoticeLocked()`） | 历史（引用他人原话）＋禁区 | 不改 |
| 13737 | 收 `33-r10` 交件面：三笔 commit＋"先探测 `firstRoundTripLocked`、最后 `serveEntry`、失败支 `serveNotBuiltNoticeLocked`"（读 diff 的记录） | 历史（引用 commit 内容）＋禁区 | 不改 |
| 13964 | 收 `33-a3`：两枚旧名的定义行/调用行现量（`:840`/`:448`/`:460`/`:451`） | 历史＋禁区 | 不改 |
| 13991 | 编排者采 ⓐ＋ⓒ 那笔：旧名→新名的对应表（改名账目本身） | 历史＋禁区 | 不改 |

## 4. 禁区：别人格的 `probes/**`（21 枚文件／85 行）——不在本程写面，逐枚仍给判

| 文件 | 行 | 那批句子在说什么 | 判 |
|---|---|---|---|
| `33/a3/verdict.md` | 8, 9 | 行号漂移记录（`:751`→`:840`、`:536`→`:625`、`:460` 未漂）＋"凡引用行号请连本段一起引" | 历史 |
| `33/a3/verdict.md` | 50, 51 | 名册表两行：定义处／注释／"是（`:461 m.mu.Lock()`）"／调用点枚数 | 历史（该腿现量） |
| `33/a3/verdict.md` | 79 | grep 枚数读数：测试调用者 3 处含 `firstRoundTripLocked`×1 | 历史（尺读数） |
| `33/a3/verdict.md` | 93, 101 | 小标题级读数：`*PanelManager.firstRoundTripLocked`（`:840`，第一句 `m.mu.Lock()`）／`serveNotBuiltNoticeLocked`（`:460`） | 历史（带行号） |
| `33/a3/verdict.md` | 128 | "`firstRoundTripLocked` 里就有一枚 5 秒 `deadline` 循环＋`pnlPumpOnce()`"——用它举例说明死锁表现为超时红 | 现在时描述代码（写面外，未动） |
| `33/a3/verdict.md` | 145, 167 | 尺读数（现量 5 命中／3 命中）；主体是最多 5 秒的消息泵（`:868`） | 历史（尺读数＋带行号） |
| `33/a3/verdict.md` | 153 | "票面会变旧：issues/33… 里 `firstRoundTripLocked` 4 处（`:39`/`:321`/`:348`/`:358`）"——正是本程处置清单的来路 | 历史（预言性记录，抹掉就没法核本程） |
| `33/a3/verdict.md` | 223, 236 | 归口表：`AC#13` 框"今天就写着 `firstRoundTripLocked` 的名字"（引票面原文） | 历史（引用他件原文） |
| `33/a3/00-anchor.md` | 24 | "涉及载体：`cmd/wisp/panel_host_windows.go` 的 `firstRoundTripLocked` / `serveNotBuiltNoticeLocked`" | 现在时描述代码（写面外，未动） |
| `33/v4/verdict.md` | 33, 34, 61, 62 | 搬动前后那 5 行逐字节相同（源码原文）＋改前调用行的 blame 是 `13acad460`/`697b4faeb` | 历史（源码原文／blame） |
| `33/v4/verdict.md` | 50 | "`rtMs` 以值过函数边界…；`firstRoundTripLocked` 的四条 `-1` 出口" | 现在时描述代码（写面外，未动） |
| `33/v4/verdict.md` | 57, 58 | 与票面 `:385` 同形："`Locked` 后缀意味着调用点应持锁。`firstRoundTripLocked`（`:751`）与 `serveNotBuiltNoticeLocked`（`:460`）自己取锁" | 现在时但**不可纯换字**（谓词是后缀本身）＋写面外 |
| `33/v4/verdict.md` | 88, 92 | 突变 `m2`/`m6` 的操作串（把那几行换成什么） | 历史（引用当时源码） |
| `33/v4/verdict.md` | 162, 165, 171 | 行号序列与 `-S` 会漏那把尺的说明（`:221`→`:227`／`:268`→`:270`） | 历史 |
| `33/r10/impl.md` | 25, 34 | 现量时序：`:425` 调、`:757-758` 发探测页；失败支 `:427-429`→`:448` | 历史（行号读数） |
| `33/r10/impl.md` | 77, 82, 85, 98, 116, 119, 124 | 逐字引用的 diff 上下文行（`-`/`+` 与注释原文） | 历史（commit 内容） |
| `33/r10/impl.md` | 141 | "`:217` `captureProbeDoc33r10`＝当场把探测页跑出来取（对一次性 sink 调 `firstRoundTripLocked`）" | 历史（带行号） |
| `33/r10/impl.md` | 154 | "探测页不是字符串常量而是运行时从 `firstRoundTripLocked` 那里取来的" | 现在时描述代码（写面外，未动） |
| `33/r10/impl.md` | 234 | "不许读成探测被削弱：`firstRoundTripLocked` 一字未动" | 历史（记那一笔没动什么） |
| `33/r10/00-anchor.md` | 20 | AC#13 框逐字抄本 | 历史（引用他件原文） |
| `33/r10/00-anchor.md` | 31, 33, 40, 41, 43 | 带行号的源码 dump（注释行 `:420`／调用 `:425`／`// firstRoundTripLocked drives…` `:716`／声明 `:732`） | 历史（源码原文） |
| `33/r10/00-anchor.md` | 52, 63, 65 | 现量叙述：`:425` 调、泵那一段 `:425`→`:427`、失败走 `serveNotBuiltNoticeLocked()`（`:448`） | 历史（行号读数） |
| `33/r11/impl.md` | 18, 19 | 改名对照表两行：旧名→新名＋声明／调用／注释枚数＋"体第 841 行就是 `m.mu.Lock()`"（这格账目本身） | 历史（旧名是新名这条记录的被引体，抹掉就没人知道改的是哪枚） |
| `33/r11/impl.md` | 49, 93 | 正控夹具形状（`roundTripHost.firstRoundTripLocked` 喂给尺必须报）＋突变 M1"把旧名原样放回"的报句逐字 | 历史（仪器输入／突变读数） |
| `33/r11/impl.md` | 160 | "旧名在跟踪文档里的残留：命中约 20 枚 `.md`…本腿禁区一字未改，这些引用从今天起是历史名"＝本程存在的来路 | 历史（该腿的具名移交） |
| `33/p1/probe.md` | 24, 65, 180 | ⑤-3 复认（`:361-368` 的 `done`）、R17 整文件全貌 dump、"射程就此定死" | 历史（尺读数／裁定） |
| `33/p1/probe.md` | 169 | "那正是今天 `firstRoundTripLocked` 会绿的那一枚" | 现在时描述代码（写面外，未动） |
| `33/a2/census.md` | 62, 251 | `done` 关闭点（`:361-368`）；R16 411 行文件全貌（`:221`/`:227`/`:349-394`） | 历史 |
| `35/a1/outbound-bridge.md` | 47, 48 | 非导出方法名册带行号（`:442`、`:656`） | 历史（dump） |
| `35/a2/transport-cost.md` | 146, 592 | 行号读数（`:449`、`:468`、`:681`）＋dump 里的 `func … serveNotBuiltNoticeLocked()` 原文 | 历史 |
| `35/census-resident-panel/summary.md` | 18, 20, 24 | 调用链带行号（`firstRoundTripLocked(:425)`）；"not built" 分支 `:441`/`:449`；引票面并注旧行号反过来 | 历史 |
| `35/r1/impl.md` | 34 | "`installPanelTransport` 在 `bringUp` 里跑于 `firstRoundTripLocked`/`serveEntry` 之前" | 现在时描述代码（写面外，未动） |
| `35/v1/verdict.md` | 43 | 现量：`:408`/`:425`/`:427`/`:467`/`serveNotBuiltNoticeLocked` 的 SetHtml `:448` | 历史 |
| `253/r3/precheck.md` | 62, 97 | 现成样板位置（`:570-615`）与 `wispProbeRT`（`:582`）读数 | 历史 |
| `274/a2/costs.md` | 165 | 纠派单过期行号的现量：`:451` 调 `m.serveNotBuiltNoticeLocked()`、`:460` 定义、`:468` HTML | 历史（尺读数） |
| `111/c2/logs/05-minus-one-exits.md` | 14, 15 | "腿报：`firstRoundTripLocked` 定义在 `:840`、`t0` 来自 `bringUp` 第 319 行"＋返回值去处（`:447`/`:448`/`:424`） | 历史（派单明确列为不许改的那形） |
| `111/c2/logs/06-ci-and-carrier-facts.md` | 28 | 产码侧对应：`:472-488 serveEntry`、`:450-452` 改走 `serveNotBuiltNoticeLocked` | 历史 |
| `bundle/1/census.md` | 75, 76, 79, 80, 83 | 逐字源码行、`:442-451` 定义位置、"调用点现量＝全仓 1 次 `:429`"、现读时序 `:426`/`:428`、`:428-429` 报错即告示页 | 历史（尺读数／源码原文） |
| `e2e-panel-1/readiness.md` | 42, 113 | H5 行的现量时序（`:340`/`:342`/`:343`/`:356-365`）＋码侧事实 | 历史 |
| `msg-a653.md` | 15 | 台账 12883 逐字引的那句读数 | 历史 |
| `panel-resident/1/hops.md` | 76, 308, 312, 315 | 干净检出走 `serveNotBuiltNoticeLocked`（`:442`）；调用/定义行号读数；产码自泵 `:656`/`:668`/`:681` | 历史 |

## 5. 汇总与本腿的裁定

- 全量命中：**102 行／27 枚文件**（写面内 11 行＝票面 5＋证据 6；台账 6 行；别人格探针 85 行／21 枚）。判**现在时描述代码**＝**10 行**；判**历史叙述**＝**92 行**。
- 10 行里：**改 1 行**（`issues/33-panel-host-c27.md:358`，纯换字、不增删行）；**改不了 3 行**（票面 `:385`＋`probes/33/v4/verdict.md:57`、`:58` 同一形：谓词断的是"`Locked` 后缀"，换名后句子从"过期"变成"荒谬"，如实更正必须加一行带时刻的注记＝本腿禁加行，且该段编排者已在台账 `A713` 自留）；**写面外 6 行**（`33/a3/verdict.md:128`、`33/a3/00-anchor.md:24`、`33/v4/verdict.md:50`、`33/r10/impl.md:154`、`33/p1/probe.md:169`、`35/r1/impl.md:34`）——本腿不扩写面，列名册交编排者处置。
- 顶回派单 1 条：§1 里"该改"的示例句（票面"两枚带 `Locked` 后缀却自己取锁"那一形＝`:385`）**不是可换字的现在时描述**，理由如上；同一形在 `probes/33/v4/verdict.md:57`、`:58` 各一枚，本腿也照"不改"判（那两枚在写面外，即便在写面仍不改）。本腿照死规矩一枚没动。
- 票面框尺（尺＝`grep -cE '^[[:space:]]*- \[ \]'`）：起手现量 `issues/33-panel-host-c27.md`＝**13 未勾**；四枚证据件＝0／0／0／0（读数＝`rulers-before.txt`）。第 2 笔改完必须复尺同数。
- 本程零代码、零测试、零 build、零 vet；零 push；两笔 commit 各带显式 pathspec。
