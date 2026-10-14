# 票 114 验收裁决表 — 腿 `114-v1`（非实现者；只读＋overlay／仓外种形突变＋重跑）

- 原告：`114-a5` 产码 `d8dd4b38`（`internal/panel/composer_test.go`＋`bridge_test.go`）＋交付 `37eb6992`（拆分表 `probes/114/a5/00-split-and-landing.md`）。
- 锚：起手 HEAD `9b51859b`（2026-10-14 16:0x +08）；窗口内他腿 `20-r2`（internal/tools+config）在飞、`frontend/` 被界面侧持续写（109→113 枚）。
- 方法：仓内零突变（两测试件 sha256 前后同值 `d2b2e848…`／`6ed0d1a4…`）；"改前"＝`go test -overlay` 换编译输入，种形＝仓外 `git archive HEAD` 树真种盘；logs 在 `probes/114/v1/logs/01…05`。⛔ 框未翻、未动 `Status:`/票名、零 push。

## 1 四格判语（结论先行）

| 格 | 判语 | 一句依据 | 翻框建议 |
|---|---|---|---|
| **AC#4** | **成立** | 自造"单删 `.json`"变异 ⇒ 新对齐尺红且逐字点名；全量回退旧 7 类 ⇒ 点名 9 枚；旧集下 `.json` 种形同时逃过结构钉（unknown 2≠3）与门二（少点名）＝"改前逃过两枚仪器"复读成立 | 可翻；留档句① |
| **AC#7** | **成立（抽核口径）** | panel 整包 windows `223/143/5/0` 与 88c7e752 对照件 `220/141/5/0` 红集**逐名同集**；vet rc=0、gofumpt -l 零行、`bash -n` 两脚本 rc=0 亲跑；**未复跑**：d22scan 全跑／长度门／POSIX（凭 a5 logs/05、06c） | 可翻，附具名未复跑项 |
| **AC#8** | **成立** | 仓外树真种盘页侧漂移（`config.set`⇒`config.sets`）⇒ 新门 FAIL 逐字点名 `config.set`；同一棵树上 253 尺（6/6/6）与 N7 **PASS**（盲） | 可翻；留档句② |
| **AC#9** | **成立（钉半支）** | 仓外树真种 `frontend/public/*.js` 门：旧根 overlay ⇒ **绿**（`89 files scanned, 2 host call sites`）＝"绿与事实无关"现场复现；同一棵树新根 ⇒ FAIL 且逐字点名 `public/ac92-v1-door.js` | 可翻；留档句③（精确加载集未答） |

## 2 逐格证据（红/绿句逐字；行号＝本窗口现读）

### 2.1 AC#4 `成立`
- 基线绿（16:05）：`door 2 reads all 15 text types the scanner walks: .css .go .html .js .json .jsx .md .svg .toml .ts .tsx .txt .vue .yaml .yml`（composer_test.go:709）。
- 我自造 m1（单删 `.json`）红逐字：`tools/d22scan's text whitelist admits ".json" and this package's door-2 filter does not read it: a route spelled in that type walks past every composer gate (ticket 114 AC#4, R-92-1), and the fix is to widen this filter, never to narrow the scanner`（:706）；同发结构钉红 `3 unknown literals were expected; got [2…]`（:662）＋`door 2 must name ac92-plant-d.json, got [2 hits]`（:674）＝旧集下 `.json` 逃过两枚仪器。
- m2（全量回退旧集）红 9 枚逐字：`.go .json .md .svg .toml .txt .vue .yaml .yml`（:705）。
- 留档句①：门二文本类集＝d22scan `isTextFile`（main.go:1312，15 类）的超集；`.mjs` 是刻意不对称，⛔ 不许为对称缩门（R16#4）。

### 2.2 AC#7 `成立（抽核口径）`
- live＝`RUN=223 PASS=143 FAIL=5 SKIP=0 rc=1`；对照（overlay 回退两件到 `88c7e752`）＝`220/141/5/0`；FAIL 名集 diff 空。
- 抽两枚逐字（两跑同句）：`composer_test.go:78: Go Snapshot emits [instructions tasks] that interface PanelSnapshot does not declare`；`backpressure_bound_35r8_test.go:309: … pinned text 4800000 runes > the cap 36864` ⇒ 5 红均非本波（3 枚前端耦合＋2 枚 35-r8）。
- 门禁 rc 亲跑：vet 0；gofumpt -l 零行 0；`bash -n d22scan.sh` 0；`bash -n check-path-length-budget.sh` 0。

### 2.3 AC#8 `成立`
- 我的种形（页侧少发一枚）：`frontend/src/lib/panel.ts never spells config.set as a literal: adding a route on one side only is the drift ticket 114 AC#8 turns into a gate - land the page literal in the same change or take the Go name out` ⇒ FAIL。
- 同一棵树：253 尺 `full inbound roster held: 6 names declared in bridge.go, 6 guard case labels, 6 router case labels…` ⇒ PASS；N7 ⇒ PASS（只认手抄四枚 `Method*`）。
- 253 不覆盖此方向的现读引句（`inbound_roster_253_test.go:19-20`）：`and it reads nothing outside internal/panel's production sources - cmd/wisp's inbound surface belongs to that package's legs.`
- 留档句②：新增名从 `bridge.go` 经 `parseFileInbound` 现取、页侧字面量为必需；⛔ 不许把清单抄回测试件（那正是 AC#8 立案的形状）。

### 2.4 AC#9 `成立（钉半支）`
- 我自己种门（仓外树 `frontend/public/ac92-v1-door.js`）：旧根 overlay＝`renderer door held: 89 files scanned, 2 host call sites…`（**PASS**，绿与事实无关复现）；新根同一棵树＝FAIL：`the renderer reaches the host from outside src/lib/panel.ts…`＋`public/ac92-v1-door.js:1: window.chrome.webview.postMessage(… "panel.mode.set" …)`＋`names "panel.mode.set"`。
- 真树改后态（16:2x）：`renderer door held: 113 files scanned, 2 host call sites (all in src/lib/panel.ts)`；真树种形自检：`red as required: 3 outside call sites, 1 assembled routes, 3 unknown literals, door 2 named 3 line(s)`。
- 两平台登记：`go doc io/fs.WalkDir`＝`The files are walked in lexical order`＋`WalkDir does not follow symbolic links found in directories`；`frontend` 零 tracked symlink（`git ls-files -s` 120000 零行）⇒ junction 侧**只登记、不演练**（按派单）。
- 留档句③：根＝整棵 `frontend`＝任何加载集的超集（不漏 loadable 文件）；"渲染器到底加载哪些文件"的精确答案未答，归界面侧（普查 J3）。

## 3 两停手格复核（判"停得对否"）

### 3.1 AC#3 — **停得对**
- `composer_handlers.go:133: if modeIsWidening(from, to) && h.Confirm == nil {`（拒在任何 `Modes.Set` 前）。
- `cmd/wisp/panel_inbound.go:254`／`:264` 两枚 `Confirm: nil,`（同一 `newComposerDispatchChain` 体；CLI seam 与常驻腿 `newResidentComposerDispatch`（panel_resident_windows.go:222/244）共用）⇒ 两条产码腿的变宽请求都被 `ErrNoL2Confirm` 拒，**走不到 Set 的 C18 卡**。
- `cmd/wisp/run.go:340` 字段＋`:677` 赋值（带真 confirm 转发）＝全仓**零读者**死码（git grep 仅注释/声明/赋值 3 处；余命中在 probes 非编译件）。
- J1"常驻腿拿什么举 C18 卡"在 spec/票面无一处具名（普查 J1）⇒ 卡片半支生产不可达＋真机对照无从构造 ⇒ 未定义即停正确。

### 3.2 AC#11 — **停得对（措辞更正一处）**
- `ci.yml:1026 run: npm run render:composer`；其后无 `git diff` 步骤；`grep -c "exit-code" .github/workflows/ci.yml` = **0**（防腐坏半支确缺）。
- 票 111 `open`、Packages 含 `ci.yml`（末条 log 自陈"本票另有 AC#11 等未勾格…"）；票 85 `open`（排队，Packages 含 `ci.yml`）；两票零处"让渡" ⇒ 票面自设前置（"先看 111/85a 是否已让出该文件"）不满足。
- ⚠ 更正：`.github` **不在** AGENTS.md §1.1／本票禁改清单的明列冻结里 ⇒ 准确根因＝**他票占用未让渡＋CI 属治理面**；"契约邻接"宜降级为上述表述（框：仍不翻，等 Q）。

## 4 拆分表抽核（随机 2 格）＋**停手上报 1 项**

- 抽样命令 `shuf -e "AC#5" "AC#6" "AC#10" -n 2` ⇒ **AC#10、AC#5**。
- **AC#5 判对**：`internal/agent/approval/ui.go:167-171` PanelAPI＝{Reject,Head,View} 无 `Allow`（`:164-166` 逐字 `It has no Allow method`）；`internal/panel/config_handlers.go:108` `lockedFieldFamilies` 拒 **7 族＋2 叶名**（a5 写"risk./fs."＝低报，方向对）；U+2190–U+2BFF 与 141 定案不是同一把——该提示与 AGENTS.md §1.2 现读一致。
- **AC#10 ②面判不成立 ⇒ 停手上报**：a5 写 `面=frontend(owner-delegated)`，但（i）a2 普查 `:198` 逐字列 AC#10 为**仪器格**（"5 枚是仪器格（AC#4/AC#8/AC#9/AC#10/AC#11）…另有 1 枚（AC#6）永不归本编队"）；（ii）a5 本行 ① 两枚落点（`cmd/wisp/panel_host_windows.go`／`internal/panel/composer_test.go`）全在本队写面；（iii）现读真门＝`panel_host_windows.go:80` const `"wispDispatch"`＋`:842 w.Bind(panelDispatchBinding,…)`（a5 引的 `:82/:319` 系普查旧锚、已腐烂；`:892` 另有探针一次性绑定 `wispProbeRT`），而钉正则 `composer_test.go:441` 只认 `.postMessage(`，页面真调用＝`panel.ts:350 window.wispDispatch` ⇒ 该格含**本队仪器写面**，"frontend 面"名不副实。③"本波不落"站得住（普查 §5-丙自陈"该上报而不是自选"＋J3 精确加载集归界面侧）⇒ 请编排者按已挂的"AC#10 去向"重裁面签。

## 5 纪律与没做的（具名）

- 只读仓内；全部突变在仓外（overlay＋`D:\tmp\wisp114v1\{ac8tree,ac9tree,mut}`）；仓内两测试件 sha256 前后同值；`GOTMPDIR=C:\tmp\wisp114v1gotmp`（反斜杠形）；只 commit、⛔ push；显式 pathspec。
- 没做：d22scan.sh 全跑／长度门／POSIX 三合包**未复跑**（凭 a5 logs，窗口非纯净）；junction 演练（派单：只登记）；AC#6/#10 的前端面判读（∉本腿派单）；真机窗口（无）。
