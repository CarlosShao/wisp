# 票 337 · 只读普查件（腿 337-a1，2026-10-18 第一波）

零写码：`.go` 命中 0 枚（尺＝逐笔 `git show --name-only` 里 `grep -c '\.go$'`）。⛔ `go test`／⛔ `go build`／⛔ `go vet`（只读腿在共享工作树跑它会真编译 internal/risk＝本波 336-r1 写面）。
起手锚自认＝`c11163c896b09e7b6863d347dadbcb65702585cc`（`2026-10-18 10:02:56 +0800`）。所有读数取具名提交树 `c11163c8`（`git show <rev>:<path>`／`git grep <pat> <rev>`），⛔ 读工作树在飞字节。
尺的公共射程＝全仓 `*.go` 并 **剔 `.scratch/**` 影子件**（旧腿在 `.scratch` 里存了大量 `bridge.go` 快照，不剔＝假命中堆；两枚未剔的原始全量件作对照件入库并具名）。逐跳 logs＝`logs/10–23,90`。
机主新铁律（lead 10-1x 转达）已写进起手件：盘符根目录零污染，本腿零导出树、零工作区。

## 1. AC#0 链逐跳读数（票面「现量」那一串的独立复跑）

| 跳 | 读数（★都带尺） | 与票面 |
|---|---|---|
| a 生产者 | A 档 case **9** 枚；走 `aRuleReason` **7** 枚（`blacklist.go:221/226/228/230/232/236/238`），走裸字面串 **2** 枚（`:223 .git/config`、`:234 browser credential store`）；拼装点唯一 `:329` | **差异①**：判语件 `80-chain-reach.log:2` 写「九枚 A case 走 aRuleReason」→ 现跑 9 枚中 7 枚（logs/16）。★成因也具名：`blacklist.go:214` 那行注释自陈「every case below returns through aRuleReason」，而同一个 switch 里 `:223`／`:234` 两支返回裸字面串 ⇒ **注释过期**（非判级差异，档位一字未变）；本波 `internal/risk/**` 写腿＝336-r1，⛔ 本腿动它，只上报 |
| b Gate 调用点 | `risk.Gate(` 生产调用点**恰 1 枚**＝`internal/tools/mode.go:98`；`internal/**` 测试对该串 0 枚 | 一致（行号未漂） |
| c `Reason` 读者 | 判定型 **0**／搬运型 **1**／打印型 **0**（尺＝收紧形 `(note\|n\|Blacklist)\.Reason[^A-Za-z]`，唯一命中 `mode.go:111`＝写侧赋值） | 一致（细节见必答④） |
| d 切片被读 | `pending_read.go:73/74/75` 三枚（Absolute／Unlockable／AlreadyUnlocked） | 一致 |
| e `blacklistedSensitive()` | 命中 **2** 行＝注释 `:117`＋定义 `:120` ⇒ 调用点 **0** | 一致 |
| f 分类器那一跳 | `internal/tools/paths.go:380` 返回 `risk.Tier`（无字符串）；`rules_gateway.go:77` 只 switch tier | 一致 |
| g `risk.Logf` | 生产注入点 **0**（9 枚赋值命中全在 `_test.go`）；sink＝包级 `var`，注入只能取 `risk.Logf =` 形 | 一致 |
| h 到模型那句 | `rules_gateway.go:82` → 融合 `assessor.go:316-352`（`:345 join`）→ `bridge.go:301` → `:439` Deny 支 → `reject(:1090)` 返回 `ToolOutcome{Text: why}`（`Text:` 在 `:1106`）→ `loop.go:745` | **差异②**：票面/判语件引的 `ToolOutcomeText(:1105)` 全仓 **0 命中**＝不存在的标识符（尺＝`git grep -n ToolOutcomeText <rev>` 剔 .scratch） |

★正控（派单 AC#0 明文要的自带正控，⛔ 交两种形状都放行的尺）：
- c 跳同一把尺换 `Decision.Reason` 读形＝非测试 **6 枚**（`bridge.go:301/439/1093/1094/1165`、`panel_pump.go:75`）；再放宽一把（含测试）＝**72 枚**（logs/13）。
- e 跳同一把尺换同文件确有调用点的 `b.permissionMode()` ＝ **1 枚**。
- g 跳同一把尺换 `opt.Sink`／`opt.Logger` 族＝非测试 **24 枚**。⚠ 我第一发正控（`risk.NowFunc` 形）落空＝那把尺不可信，旧发原样留 `logs/15`，决定性读数以 `logs/17` 为准。
- b 跳同一把尺换 `risk.Resolve(`＝非测试 **1** 枚（`paths.go:116`）、含测试 7 枚。⚠★本腿先把预期写成「确有 ≥2 枚」＝错，已当场具名更正（logs/11 末节）。**四跳的 0 都不是尺坏。**
★本腿自抓两枚尺坑并具名：① `n\.Reason` 子串命中 `turn.Reasoning`（llm 的字段，与 Note 无关）⇒ 收紧 `[^A-Za-z]`；② `ls-tree` 未关 `core.quotePath` 时中文名以八进制转义呈现，`grep 264` 会假命中 ⇒ 重跑用 `-c core.quotePath=false`（logs/21）。

## 2. AC#1 三形代价表（先量再选；★本腿⛔ 选形、⛔ 推荐、⛔ 替机主决定）

**甲＝只接进本机日志（`risk.Logf` 补注入点）**
- 出向面变化枚数：出向**零**变化（日志非出向面）。要动的字节＝一个注入点；今天已存在的同类形状＝`cmd/wisp` 里 `Logf:` **6 枚**生产注入（`panel_inbound.go:255`、`resident_approval_windows.go:375`、`run.go:494/630/670/787`），但那些接的是**另一枚同名 sink** `tools.Options.Logf`（`bridge.go:73`）；两枚 sink 之间今天的接线＝**0 枚**（logs/22）。落点行的现长度尺＝logs/23。
- 触碰 `frontend/**`：**假**。要不要新导出名：**假**（`risk.Logf` 已是导出 var，`blacklist.go:51`）。
- 绝对路径＋控制字符进模型上下文：**假**（不进模型）。
- 不做的代价／没防住的形状：那句实话今天已按 `rule=` 参数在 sink 上（`blacklist.go:83`），sink 为 nil ⇒ 整条丢 ⇒ 机主与审计日志都看不到拒绝的真实理由，票 334 买到的只活在返回值里。

**乙＝接进面板那条拒绝文案（走 `C17` 现有方法面，⛔ 新方法名）**
- 出向面变化枚数：面板字节面今天**完全不带** `BlacklistNote`（`internal/panel` 非测试对 `Blacklist` **0 命中**，logs/22）。既有承载体只有三处：`panel/approval.go:48 json:"reason"`（注自陈 verbatim）、`:83 Reason: decision.Reason`、`cmd/wisp/panel_pump.go:75`。⇒ 乙形若走"复用现有键"＝改的是**卡片 reason 这句话本身**；若走"新键"＝票面 AC#5 禁区。
- 触碰 `frontend/**`：**假（就出向字节而言）**——前端 `l2-approval-card.tsx:207/:216` 已在渲染 `view.reason`（不改前端也会显示）；⚠ 但 `frontend/**` 归机主外部 agent（本编队零写面），任何显示侧改动都需另排。
- 新导出名：**取决于走哪条路**（本腿不选）。★真缺口：面板那条 reason 今天来自 `risk.Decision` 融合句（R3 中文短句），与 Gate 那句英文长话**不同串**（logs/22）。
- 绝对路径＋控制字符进模型：**假**（进人眼不进模型）。⚠ 但 216 显示底座的清洗名册今天**不含**卡片 Reason：`display_scrub_216.go:58` 自陈「does not touch ApprovalCardView.Tool/Args[]/Reason」，`scrubForDisplay` 生产调用点仅 `composer.go:145` 一处 ⇒ 把带真路径的真串交给卡片，**没有 216 那层洗**（读码推的；尺＝上述两处调用点枚数）。
- 不做的代价／没防住的形状：机主屏幕上仍是那句不含理由的拒绝（票面"人看到的"那一半继续不成立）。

**丙＝接进模型正文（`ToolOutcomeText` 那一跳）**
- 出向面变化枚数：那一跳今天已有出向句——`bridge.go:439` 拼 `'该目标被绝对禁止访问（R3 A 档，任何授权都不可豁免）：' + dec.Reason`，`dec.Reason`＝R3 那句**已含 `canonical` 绝对路径**（`rules_gateway.go:82` 的 `%s`）⇒ 丙形不是"新增一句"，是**换／加长已有那句**。落点行现长度尺＝logs/23。
- 触碰 `frontend/**`：**假**。新导出名：**假**（改的是既有拼装行的字符串）。
- 绝对路径＋控制字符进模型上下文：**真**（且今天模型侧无拦网：`sanitizeArgs`（`bridge.go:1290`）只作用 `ArgsJSON` 的 forensics 列（`:1201`），不作用进模型的 `Text`＝读码推的，尺＝该函数调用点枚数）。对照先例逐枚点名（都是跟踪件真名）：**票 264**＝`error-text-carries-absolute-paths-into-model-visible-text`（失败句把绝对路径逐字发给云端服务商、这一侧一枚仪器都没有）——丙形正是同一族形状；**票 216**＝控制字符显示注入底座，且自陈不覆盖卡片 Reason——丙形把带控制字符的真串送进上下文，两侧都没有这道洗；**票 174**＝指针诚实性族（`docs/evidence/s1/174-pointer-honesty-r1.md` 等）——"给模型一个它读不到的路径"这一形，丙形给真候选路径时同样撞上（默认 `[fs] allowed_dirs` 空）。
- 不做的代价／没防住的形状：模型永远只看到"被绝对禁止"而看不到"为什么"，它无法据此改计划；但代价同时是泄露面（上述三族）。

★★硬门（票面 AC#2／派单）：未经台账 `A##` 具名批准，⛔ 任何人把这句接进模型正文或面板出向内容；默认只落本机日志。乙／丙属功能级，由编排者摆给机主拍板——**本腿不选形、不推荐**。

## 3. AC#3 前置料（只量现状，⛔ 提改法）

门拦下文案＝`bridge.go:439` 中文短句 ＋ R3 融合句；权限判定文案＝`rules_gateway.go:49` `'R2: 目标路径在授权目录之外: %s'`。两支今天**不同形**（不同规则号、不同语言壳）。★恒真判据那一形的既有钉子＝`bridge_junction_windows_test.go:309-310`（逐字："红队拒绝不得与'不在授权目录内'共用措辞（否则未来的询问流会被注入文档驱动）"）＋`:462` 那条反向钉（解析失败不得含"A 档"）。票面点名的 242／244 一族真名＝`internal/agent/approval/ticket242_binding_test.go`／`ticket242_panelface_test.go`／`internal/agent/corr_percall_242_test.go`／`docs/evidence/s1/244-black-console-r1.md`。⇒ 若落地时把 reason 支与权限支写成同串，这几枚是现成的咬合点（本腿⛔ 跑它们）。

## 4. AC#4 前置料（成对台面点名列，⛔ 跑）

甲＝`internal/risk`＋`cmd/wisp`；乙＝`internal/tools`＋`internal/panel`＋`cmd/wisp`（＋`frontend/**` 归机主外部 agent）；丙＝`internal/tools`＋`internal/agent`。共同约束（派单编队）：`cmd/wisp/**` 与 `internal/panel/**` 本波零写腿；`internal/risk/**` 本波唯一写腿＝336-r1；`go test ./cmd/wisp` 本波唯一持跑者＝228-r3c；台面须仓外导出树且目录名产生 8.3 短名、`GOTMPDIR` 先建——★落点按新铁律只许 `D:\tmp\wisp\20261018\<腿名>\`，⛔ 盘符根目录。红名册里 `internal/panel` 那两枚前端契约键族与 `TestC21…` 属别面存量（票面 AC#4 明文），⛔ 计入本票欠账。

## 5. 门禁与 LF（本腿零 .go，只跑适用两枚；logs/90 逐枚落 rc）

`sh scripts/d22scan.sh` rc＝0（`go run` 打的是它自己的扫描模块＝guard 本身）；`sh scripts/check-path-length-budget.sh` rc＝0。`gofumpt -l`／`go vet` 对零 `.go` 的普查腿＝不适用，具名不跑。★新字节 LF：工作树侧 CR＝0（bash 通道落笔）；本会话 `core.autocrlf=true` ⇒ 早期 PowerShell `Out-File` 那批已就地转 LF 并重落，入库 blob CR＝0（尺＝`git cat-file blob` 数 `\r`）。存量 6 枚脏＝派单具名，⛔ 顺手重排别腿文件。

## 6. 必答（派单末四问）

**(i) "消费者 0 枚"那把尺的正控逐字读数**：同一把尺（`git grep -nE '(dec|it\.Decision|decision)\.Reason[^A-Za-z]' c11163c8 -- "*.go"` 剔 `.scratch`、剔 `_test.go`）命中 **6** 枚——`bridge.go:301/439/1093/1094/1165`、`panel_pump.go:75`；放宽一把（含测试）**72** 枚。逐字见 `logs/12-hop-c-blacklistnote-landing.log` 尺三。
**(ii) 逐跳读数与票面差异清单**：**2 枚**＝① hop a 分母：票面正文未写分母，参照件 `80-chain-reach.log:2`（＝派单指定逐跳对照的判语件）写「九枚 A case 走 aRuleReason」→ 现跑 9 枚中 7 枚（logs/16）；② hop h 的 `ToolOutcomeText` 全仓 0 命中＝不存在的标识符（票面第 13 行与判语件都引它；真身 `ToolOutcome.Text`，`bridge.go:1106`）。其余 a—h 六跳**零差异**（尺见 §1 表与 logs/10–23）。★另具名一枚非链读数的发现：`blacklist.go:214` 注释自陈过期（见 §1 差异①成因行），属产码注释、⛔ 本腿动。
**(iii) "会动 `frontend/**`"逐形点名**：甲＝**假**；乙＝**假**（就出向字节而言：前端 `l2-approval-card.tsx:207/:216` 已渲染 `view.reason`；⚠ 但 `frontend/**` 是机主外部 agent 的写面，任何显示侧改动需另排）；丙＝**假**。
**(iv) 这句实话今天最远流到哪一枚字节就停了**：最远＝`internal/agent/approval/pending_read.go:67` 的 `out := d`（`cloneDecision` 整枚 struct 值拷贝，把 `BlacklistNote` 连 `Reason` 一起搬进 `LiveApproval.Decision`）＝**搬运型接触，不指名读该字段**；再往前一步是写侧 `internal/tools/mode.go:111`。之后判定型／打印型／面板出向／模型正文**全部 0 位读者**（`internal/panel` 对 `Blacklist` 非测试 0 命中）。⇒ 这是对票面"读者 0 枚"的一处**细化**（搬运型 1 枚），不推翻它。

## 7. 本腿未定案／上报（⛔ 自行假设）

1. 票面第 13 行把 `ToolOutcomeText(:1105)` 当既有标识符引用——请编排者裁定是否只追加更正注记（⛔ 本腿动票面原文）。
2. 必答④那枚"搬运型 1 枚"要不要计入票面"读者 0 枚"的定义——归裁决者，⛔ 本腿盖章。
3. `internal/panel/display_scrub_216.go:58` 那句自陈（卡片 Reason 不进 216 清洗名册）对乙形是**既有事实**还是**待立缺口**——属票 216／乙形裁量，本腿只报盘上字节。
