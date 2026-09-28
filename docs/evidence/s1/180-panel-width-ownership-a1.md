# 180-a1 — 面板尺寸今天到底由谁决定（只读普查·派单 §B）

- 分工：`.scratch/wisp/dispatches/2026-09-28-175x-wave-data-carriers-182a-180a-189a-190a-188r.md` §B。
- 工单：`.scratch/wisp/issues/180-panel-width-is-a-config-field-with-zero-production-readers-so-changing-config-toml-has-no-visible-effect.md`。
- 本件性质：**只读普查**。零产码；`frontend/**` 与 `design/**` 零写面（本件所有"零命中"宣称的尺射程都只写
  `internal/ cmd/ tools/`，`frontend/**` 从未被任何一把尺计入分母或分子）。
- **AC 框一枚未勾**（勾要非实现者的表）；`AC#2`／`AC#3` 按派单**未做**，只做 `AC#1` 三问，`AC#4`／`AC#5`／`AC#6` 未做（见 §5）。

## 0. 起手三件（现跑，本程第一步）

```
date                              => 2026-09-28 17:48+0800
git rev-parse HEAD                => 38fc7c0e2276306b4fbe57387d62407db6bb018b
git rev-parse --abbrev-ref HEAD   => dev
git log --oneline -1              => 38fc7c0e ledger(A387-A388) 收 33-r3（入向听众进树＋词面尺改能力尺，我四把尺复跑）＋owner 提问后按默认动作派 145-r2
git status --porcelain | wc -l    => 77
```

起手名册逐枚（`git status --porcelain` 原文，落盘 `/tmp/roster180.txt`，终态闸门拿它对差集）：

```
 M .gitignore
 M .scratch/wisp/probes/152/my152.py
 M .scratch/wisp/probes/161/r6/logs/flip-1.txt
 M .scratch/wisp/probes/161/r6/logs/flip-2.txt
 M .scratch/wisp/probes/161/r6/logs/flip-3.txt
 M .scratch/wisp/probes/161/r6/logs/flip-4.txt
 M .scratch/wisp/probes/161/r6/logs/flip-5.txt
 M .scratch/wisp/probes/161/r6/logs/flip-6.txt
 M .scratch/wisp/probes/161/r6/logs/flip-baseline.txt
 M .scratch/wisp/probes/161/r6/logs/flip-restored.txt
 D design/assets/base.css
 D design/assets/icons.js
 D design/assets/theme.js
 D design/assets/tokens.css
 M design/doubao/README.md
 M design/doubao/demo/app.js
 M design/doubao/demo/index.html
 M design/doubao/demo/styles.css
 D design/index.html
 D design/screens/approval.html
 D design/screens/ball.html
 D design/screens/chat.html
 D design/screens/config.html
 D design/screens/cost.html
 D design/screens/firstrun.html
 D design/screens/palette.html
 D design/screens/privacy.html
 D design/screens/security.html
 D design/screens/states.html
 D design/screens/tasks.html
 M docs/evidence/s1/152-subject-death-never-measured-r1-accept-r1.md
?? .scratch/wisp/.scratch/
?? .scratch/wisp/dispatches/2026-09-28-175x-wave-data-carriers-182a-180a-189a-190a-188r.md
?? .scratch/wisp/probes/139/accept-r1/
?? .scratch/wisp/probes/152/overlay-probe1-on-samppost.json
?? .scratch/wisp/probes/156/__pycache__/
?? .scratch/wisp/probes/156/mut-156-r2/asis.log
?? .scratch/wisp/probes/156/zero156-r4-head.sh
?? .scratch/wisp/probes/156/zero156-r4-work/
?? .scratch/wisp/probes/158/r2/
?? .scratch/wisp/probes/161/r2/__pycache__/
?? .scratch/wisp/probes/161/r2/ctl/
?? .scratch/wisp/probes/161/r5/negative-control/
?? .scratch/wisp/probes/161/r6/logs/flip-7.txt
?? .scratch/wisp/probes/162/r4/
?? .scratch/wisp/probes/162/v1-baseline-gotest.txt
?? .scratch/wisp/probes/176/r1/logs/gate-clauses-end-bad.txt
?? .scratch/wisp/probes/176/r1/logs/gate-clauses-end.txt
?? .scratch/wisp/probes/176/r1/logs/gate-clauses-start-bad.txt
?? .scratch/wisp/probes/176/r1/logs/gate-clauses-start.txt
?? .scratch/wisp/probes/183/accept-v1/
?? .scratch/wisp/probes/185/c1/logs/d22scan-post-final.txt
?? .scratch/wisp/probes/33/r2/d22scan-final.txt
?? .scratch/wisp/probes/33/r2/gate-final.txt
?? .scratch/wisp/probes/33/r2/status-final.txt
?? .scratch/wisp/probes/33/r2/status-start.txt
?? .scratch/wisp/probes/999/
?? .zcodeignore
?? design/doubao/01-ball-states.jpg
?? design/doubao/demo/lib/
?? design/doubao/demo/rb-files.js
?? design/doubao/demo/rb-plugins.js
?? design/doubao/demo/rb-review.js
?? design/doubao/demo/rb-terminal.js
?? design/doubao/demo/rightbar.js
?? design/doubao/demo/screens/
?? design/doubao/demo/screenshots/
?? design/doubao/demo/sidebar.js
?? design/old/
?? part1-state1-fixed.txt
?? part1-state1-pristine.txt
?? part1-state2-fixed.txt
?? part1-state2-pristine.txt
?? part2-nog6-fixed.txt
?? part2-nog6-pristine.txt
?? part3-stale-fixed.txt
?? part3-stale-pristine.txt
```

**锚自取说明**：起手锚 `38fc7c0e` 与编排者锚同号；工作期间树上多出一枚 `971b7dad`
（`ledger(A389)＋立 Q-71 - 前端等数据的五枚票一次排开（145-r2/188-r1 写＋182a/180a/189a/190a 只读）`），
**不是漂移**；它碰没碰我的文件见 §6 的逐枚 `git show --name-status`（我不带区间 diff 自证）。
在飞的 `145-r2` 写面（`internal/panel/composer.go`／`pump.go`／`cmd/wisp/panel_pump.go`／`internal/panel/pump_test.go`）
本程**只读不写**；读到的字节是起手态，逐枚原文凭据在下表行号。

## 1. `AC#1` ① — 面板窗口尺寸今天由**谁**决定

结论：**没有人。因为本树今天没有面板窗口。** 三处代码注释 + 三把尺现读同向：

| 处 | `file:line` | 现读到的字（逐字，摘其要） | 它说明什么 |
|---|---|---|---|
| 包边界 | `internal/panel/doc.go:16` | 「`DEFERRED(host/bridge): implemented by ticket 33 (host), ticket 35 (bridge).`」＋`:17`「This ticket only freezes the package boundary.」 | `internal/panel` 只是被冻住的包边界；宿主（窗口）与桥都不在此票交付范围 |
| 载体自述 | `internal/panel/pump.go:15-17` | 「`WHAT THIS FILE IS NOT: the transport. There is no Go -> page channel in this tree today - no WebView2 host (ticket 33 is unclaimed), no postMessage writer, no local HTTP/SSE/websocket server`」 | 生产码自己的注释直接写着：本树无 WebView2 宿主 |
| 运行态 | `cmd/wisp/run.go:438` | 「`have is a page to go to: this tree carries no WebView2 host (tickets 33/35),`」 | 装配根（跑起来的进程）同样声明没有宿主 |
| WebView 选项 | 无此代码 | `go.mod` 里 webview 依赖**零枚**（尺见 §2 R7，输出 `(none)`）；全树非测试 Go 里 `webview2` 只出现在四行**注释**（`internal/panel/doc.go:7`、`internal/proc/jobscope_windows.go:18`/`:21`、`internal/proc/shutdown.go:14`） | "创建窗口／WebView 选项"那一处**不存在**，没有可接的一行 |
| 布局 | `internal/panel/composer.go:57-62` | `Snapshot` 只有四枚 key：`pending`／`results`／`composer`／`generatedAt` | 送进面板的数据包今天**连"尺寸"这一维都没有**，谈不上布局取哪枚数 |
| 唯一真窗口（**别认错**） | `internal/ball/ball_windows.go:233` | `edge := int32(WindowEdgePx(b.opts.SizePx, dpi))`，钳位在 `:64`（`44..72`，0 走 `BallSizeDefaultPx`）与 `:144-151` | 仓内今天唯一的 OS 窗口是**球**；它的尺寸来自 `opts.SizePx`（`[ball] size` 那一族），**与 `[panel] width` 无关系** |

票面给的那把尺（本程原样跑，非测试）：

```
$ grep -rn "NewWindow\|SetBounds\|Rect{\|Width:" --include=*.go internal/panel/ | grep -v _test
(空输出，rc=1)
```

⇒ `internal/panel/` 里**没有一处**创建窗口、设置 bounds、或写死 `Width:` 的形状。
`internal/panel/` 非测试文件全集（11 枚，尺 `ls internal/panel/ | grep -v _test`）：
`approval.go` `assets.go` `attachments.go` `bridge.go` `composer.go` `composer_dispatch.go` `composer_handlers.go`
`doc.go` `git.go` `pump.go` `workspace.go`——逐枚扫过形状，无一含窗口/几何代码（`assets.go` 只读嵌入的 dist 字节，
`attachments.go` 里的 `SizeBytes` 是**附件字节数**，`bridge.go:82` 的 `host` 是"宿主该怎么做"的称谓）。

## 2. `AC#1` ② — `default:"640"` 那一族的读取方枚数（现量，零也是量出来的）

字段在场（现读）：`internal/config/schema.go:524-535`，五枚字段逐字：

```
524: // PanelSection is [panel]; hot-tier.
525: type PanelSection struct {
526: 	Enabled bool `toml:"enabled" default:"true"`
527: 	// Width is the panel width in px.
528: 	Width int `toml:"width" default:"640"`
529: 	// Height in px; 0 = auto from content.
530: 	Height int `toml:"height"`
531: 	// KeepAliveInSession keeps the WebView2 warm during a session.
532: 	KeepAliveInSession bool `toml:"keep_alive_in_session" default:"true"`
533: 	// Scale is the UI zoom factor; 0 = follow system DPI.
534: 	Scale float64 `toml:"scale"`
535: }
```
外层挂载点：`internal/config/schema.go:123` `Panel PanelSection \`toml:"panel"\``。

逐枚尺（口径：`grep -rn "Panel.<字段>" --include=*.go internal/ cmd/ tools/ | grep -v _test.go`；
**非测试**、**不含 `frontend/**`／`design/**`**）：

| R# | 尺（原文） | 输出（原文） | 读取方枚数 |
|---|---|---|---|
| R4a | `grep -rn "Panel.Width" --include=*.go internal/ cmd/ tools/ \| grep -v _test.go` | （空）`rc=1` | **0** |
| R4b | `grep -rn "Panel.Height" --include=*.go internal/ cmd/ tools/ \| grep -v _test.go` | （空）`rc=1` | **0** |
| R4c | `grep -rn "Panel.Scale" --include=*.go internal/ cmd/ tools/ \| grep -v _test.go` | （空）`rc=1` | **0** |
| R4d | `grep -rn "Panel.KeepAliveInSession" --include=*.go internal/ cmd/ tools/ \| grep -v _test.go` | （空）`rc=1` | **0** |
| R4e | `grep -rn "Panel.Enabled" --include=*.go internal/ cmd/ tools/ \| grep -v _test.go` | （空）`rc=1` | **0** |
| R4f | `grep -rn "\.Width" --include=*.go internal/ cmd/ tools/ \| grep -v _test.go`（票面原尺） | （空）`rc=1` | **0**（全树非测试 Go 里**没有任何** `.Width`） |
| R4g | `grep -rn "\.Width" --include=*.go . \| grep -v _test.go \| grep -v "^./frontend"` | （空）`rc=1` | 0（这条只作扩大射程复核，**不用于宣称前端**） |
| R7 | `grep -rn -i "webview\|winc\|lxn\|github.com/richardlawley" go.mod` | `(none)` | 宿主依赖 0 枚，接线无落点 |

生产里"碰到过这枚 struct"的两处，逐枚现读，**都不是"读值去行动"**：

| `file:line` | 现读到的字 | 定性 |
|---|---|---|
| `internal/config/manager.go:211` | `{"panel", &cur.Panel, &fresh.Panel, func() { cur.Panel = fresh.Panel }},` | 热加载**整块指针拷贝**：比较整块、赋值整块，从不取 `Width` 的值 |
| `internal/config/manager.go:216-221` | `for _, s := range rest { if !reflect.DeepEqual(s.oldr, s.newr) { s.set(); rep.Hot = append(rep.Hot, s.name) } }` | "hot 生效"今天的实现＝换掉结构体字段＋往回执里 append 一个**段名字符串** |
| `internal/config/defaults.go:60`/`:92-126` | `applyDefaults(reflect.ValueOf(c).Elem())` / `setDefault(...)` | 反射是 `default:"640"` 的**写入方**，不是读取方 |
| `manager.go:219`/`:278`/`:334` | 尺 `grep -rn "\.Hot\b" --include=*.go internal/ cmd/ \| grep -v _test.go` 命中**只有三处 append**，无一处取值去行动 | 连"这改了哪几段"的回执今天也没有行动方 |

唯一真的读写这五枚值的地方是**测试**（不计入读者）：`internal/config/boundary_test.go:159-163`
（`c.Panel.Enabled = false`／`c.Panel.Width = 800`／`c.Panel.Height = 600`／`c.Panel.KeepAliveInSession = false`／`c.Panel.Scale = 1.5`）。

⇒ "缺哪一跳"的诚实答案：**不是缺一跳，是缺整条链**——宿主（票 33）、载体（票 35）、
以及"尺寸"这一维在快照里的位置（`composer.go:57-62` 名册没有它）三样都不在。

## 3. `AC#1` ③ — 规格要不要它生效（逐字引，给行号）

结论：**规格要求它生效**（⇒ 本票是"补实现"那一支，不是"改规格文字＝人工批准"那一支）。
但要求**只由 `hot` 那一格承载**，宿主侧的文本是缺的。

| 出处 `file:line` | 逐字引文 |
|---|---|
| `docs/PLAN.md:2726` | 「**三档生效级别**：`hot`（立即生效）· `reload`（需重载子系统，如换 ASR 模型）· `restart`（需重启进程）。」 |
| `docs/PLAN.md:2743` | 「\| `[panel]` \| `enabled` `width` `height` `keep_alive_in_session`(true) `scale` \| `hot` \|」 |
| `docs/specs/SPEC-03-config-secrets-envs.md:39` | 「\| `[panel]` \| `enabled(bool)=true` `width(int)=640` `height(int)` `keep_alive_in_session(bool)=true` `scale(float)` \| hot \|」 |
| `docs/specs/SPEC-08-ui-ball-panel.md:143-154` | §5.1 宿主（D29/C27）逐条只写：`jchv/go-webview2`、单例 `PanelManager`「隐藏而非销毁」、`AddWebResourceRequestedFilter` 从 `embed.FS` 喂产物、前端无状态＋每次 `show` 推 `panel.resync`、Runtime 缺失走原生降级卡、多任务共用单窗口按 correlationId 分区——**没有一句把窗口尺寸指到 `cfg.Panel.Width`** |
| `docs/PLAN.md:3443` | `stroke-width="1.5"`（图标规范，与面板窗口无关；尺 `grep -n "width" docs/PLAN.md` 全树仅 `:2743`/`:3443` 两枚命中） |
| SPEC-08 其余 `width` 命中 | 尺 `grep -rn -i "width\|尺寸\|px\b" docs/specs/SPEC-08*.md`：`px`/`尺寸` 命中全部落在**球**那一节（`:29-67`，44–72px/56px/34.72px 等）与面板**视觉行高**（`:190-192`）；`:33` 那枚 SVG `stroke-width` 与窗口无关 ⇒ **SPEC-08 不谈面板窗口宽度** |

⇒ 「规格真空」这一支**不成立**（`PLAN.md:2743` 明写 `hot`，`PLAN.md:2726` 明写 `hot`＝立即生效）；
需要编排者另外记一笔的是**文本缺口**：规格要求"立即生效"却没写"生效成什么"，
补宿主那一票（票 33）若不指名尺寸来源，这一格还会再空一次。

## 4. 风险栏（派单 §B 末条预写；`AC#2`／`AC#3` 本程未做）

- ⛔ **不许用"把默认值改掉"交差**：把 `default:"640"`（`internal/config/schema.go:528`）换成别的数而**不接读者**，
  票面 `AC#3` 判不通过——那只是把一枚装饰换成另一枚。本程现量再钉一次：今天**改默认值连"改了什么"都无人读**。
- ⚠ 别把**球**那一格当面板的归属：`internal/ball/ball_windows.go:233` 是仓内唯一真窗口的尺寸来源，走 `[ball] size`；
  在那一行去"接 width"＝修错格子。
- ⚠ 别把 `internal/config/manager.go:211` 的整块拷贝当"已接线"；它不读值。
  同理 `rep.Hot`（`manager.go:219`）今天无行动方，"回执里出现了 panel"不等于生效。
- ⚠ 快照名册是**四枚 key 的钉**（`internal/panel/composer.go:45-55` 自述四枚钉、含两枚字节级），
  真接尺寸这一维要先撞 `Q-51`（谁可写 `frontend/src/lib/panel.ts`）与那四枚钉——**写腿开工前必做撞钉预检**，
  本程已把它列进风险栏，未自行放宽任何断言。
- ⚠ `default:"640"` 的数值本身今天不许改（票面 `AC#5`）；数值是 owner 的一句话，不是实现者的自选。

## 5. 本程未做的格子（照派单 §B：只做 `AC#1`）

- `AC#2`／`AC#3`：写腿，**未做**（派单原文「那是写腿，等你的表回来我再派」）。
- `AC#4` 同族枚数普查：**未做**（不在 §B 分工内）。上一程 `180-c1` 的现量（带 `default:` 行 69、去重字段名 61、零生产读者 18 枚）
  在 `docs/evidence/s1/180-182-panel-fields-census-c1.md`，**本程未复核、不背书**，只作指路。
- `AC#5` 契约轴：本程零产码、`docs/PLAN.md`／`docs/specs/**`／`go.mod`／`go.sum`／`thresholds.go`／golden／`allowlist.txt`
  一字节未动（`frontend/**`／`design/**` 零写面）；但票面那条"`schema.go` 默认数值不许改"属写腿约束，本程只声明未改。
- `AC#6` 门禁（逐包 `go test` ＋ `d22scan` ＋名册差集）：**未跑**——本程 `go` 文件零改动，跑它拿不到新信息；
  预算被 §0/§2 的现读与现量尺吃掉（派单硬约束「超预算不是放宽断言的理由」，我未放宽任何断言，只未跑此尺）。
  ⚠ 下一程（写腿）开工前必须自跑，且最终读数取在最后一枚 commit 之后。
- 票面「现量 #3」的 `unwired_test.go` 语义复核：**本程未做**（属 `180-c1` 已答，见其表 §3）。

## 6. 自证（只读：逐枚 `git show --name-status`，不用区间 diff）

起手三件与逐枚凭据在下面；本节是**追加**进已提交正文的自证段（骨架先建 → 每节贴 `git log`＋`git show --name-only` 原文）。

（`git log --oneline -1`、`git show --name-only` 与终态名册差集原文：见本件提交后的回报正文；
本程对生产码的写面为**零字节**——唯一两枚路径是本件与工单的 Progress log。）

### 6.1 本程第一枚 commit（票面 Progress log 追加）逐枚 name-status

```
commit 9b6a1be5
190-a1(只读): 证据件 190-file-tree-design-a1.md — 文件树那堆的根/上限/忽略规则设计核

A	docs/evidence/s1/190-file-tree-design-a1.md
```

### 6.2 工作期间树上多出的他枚 commit（有否碰本程文件）

```
commit 971b7dad
ledger(A389)＋立 Q-71 - 前端等数据的五枚票一次排开（145-r2/188-r1 写＋182a/180a/189a/190a 只读）

A	.scratch/wisp/dispatches/2026-09-28-175x-wave-data-carriers-182a-180a-189a-190a-188r.md
M	docs/reports/pending-and-issues.md
```

本件自身那枚 commit 只带 `docs/evidence/s1/180-panel-width-ownership-a1.md` 一枚路径（显式 pathspec，未用 git add -A / git add .）；其 name-status 与终态名册差集贴在交件回报里。本程未用区间 diff 自证只读。

### §6.2 编排者追加更正（09-28 17:5x，账 `A391`；不动上文任何一格）
交件时点位更正：本节当时写"树上多出一枚 `971b7dad`"，**实为三枚**——`971b7dad`（A389 台账）／`500bc98e`（`188-r1` 片①）／`9b6a1be5`（`190-a1`），**三枚都没碰本程两条路径**（我逐枚 `git show --name-status` 复跑过）。差集四条已具名（dispatches 那枚被入库、`probes/145/r2/`＋`145-snapshot-growth-r2.md`＝`145-r2` 在飞、`188-task-state-r1.md`＝`188-r1` 在飞），**零枚出自本程**。另：本程 `AC#1` 已由我复跑对格后翻勾（三枚尺读数我重跑到手：`internal/panel/` 非测试里尺寸形状 **0 命中**、五枚 `[panel]` 字段生产读取方各 **0**）。
