# 77 — scaffold 第二批裁决表素材：AC#3（L2 卡真数据）＋ AC#11（深浅色渲染与对比度）

- 交付者：`agent-93`（2026-11-16 19:27-20:0x，frontend 唯一写腿，按派单
  `.scratch/wisp/probes/orch/20261116/txt/77-scaffold-batch2-dispatch.txt` 只打 AC#11 与 AC#3）
- AC 框一律不勾，勾选与裁决归非实现者。
- 仪器/读数件：`frontend/scripts/render-shell-77.tsx`（新，AC#11）、
  `frontend/fixtures/l2-card-unjudged-reason.json`（新，AC#3 未判负控）、
  `77-scaffold-batch2-contrast.txt`（本轮最终对比度逐对读数，`npm --prefix frontend run render:shell77` 全输出）。
- 既有仪器复用：`frontend/scripts/render-l2.tsx`（`npm run render:l2`，CI 第 6 步已在用）。

## AC#3 — L2 确认卡渲染真数据

### 判据①：`internal/risk` 真决策对象 → 面板 wire → 卡，整条 file:line

| 跳 | 证据（file:line） | 内容 |
|---|---|---|
| 1. 真决策对象 | `internal/risk/assessor.go:230` | `RiskAssessor.Assess(tool, params, facts) Decision`；`Decision` 定义 `:92`；`Level.String()` `:56`；`NewRiskAssessor()` `:190` |
| 2. 决策→卡视图（纯映射） | `internal/panel/approval.go:64-67` | `NewApprovalCardView` 跑真评估器；`:72-99` `CardViewFromDecision` 把 `risk.Decision` 映射为 10 键 `ApprovalCardView`（`ReasonKnown = reason 非空`，`:84`） |
| 3. 生产装配根 | `cmd/wisp/run.go:1185-1186` | `rt.pump = panel.NewSnapshotPump(panel.PumpSources{Verdicts: rt.liveVerdicts, ...})`；`rt.liveVerdicts` = `cmd/wisp/panel_pump.go:50` |
| 4. 驻留泵同一来源 | `cmd/wisp/panel_resident_pump.go:94-99` | `Verdicts: func() []panel.NativeVerdict { return nativeVerdictsFromGate(gate) }`——与 run.go 共用同一映射函数，无第二张卡 |
| 5. verdict→卡（不二次评估） | `internal/panel/pump.go:85-99` | `NativeVerdict` 是「internal/risk 已判定、internal/agent/approval 已放行」的副本；`:106-121` `CardView()` 走同一个 `CardViewFromDecision`（approval.go:72）⇒ 推送卡与 CLI 卡对同一 verdict 不可能分歧 |
| 6. 快照装配 | `internal/panel/pump.go:349-360` | `Snapshot()` 把 verdicts 变 cards；`internal/panel/composer.go:58` `Pending []ApprovalCardView` json:"pending" |
| 7. 推送 | `cmd/wisp/panel_snapshot_deliver_329.go:57` | `window.wispOnPanelSnapshot(<json>)`（每次 publish 整包推送） |
| 8. 页接收与类型 | `frontend/src/panel-root.tsx:37` 装 `frontend/src/lib/panel.ts:568-578` `installPanelSnapshotReceiver`；类型 `panel.ts:24-47`（`ApprovalCardView`）与 `panel.ts:470-471`（`PanelSnapshot.pending`） | 类型被 `TestApprovalCardViewJSONKeysMatchFrontendTypes` 双向钉（Go↔TS 键集） |
| 9. 卡渲染 | `frontend/src/App.tsx:351-353` | `snapshot.pending.map(v => <L2ApprovalCard view={v} key={v.correlationId}/>)`；字段绘制 `frontend/src/components/l2-approval-card.tsx`：工具名 `:292`、风险级徽标 `:93-120`、参数 `:139-168`/`:305`、原因 `:206-218`/`:310`、命中规则 `:175-201`/`:312`、调用链 `:316-319`、判定来源 `:338` |

生产调用者（能力类判据）：`run.go:1185` 的装配根——每次 panel publish 由
`panel_snapshot_deliver_329.go:57` 整包推送 `pending`；今日另有命令级调用者
`wisp panel-assets -l2`（`cmd/wisp/panel_assets.go:68`，同一 `NewApprovalCardView`）。

### 判据②：行为牙（两枚，均 rc=0）

- **牙 1 `render:l2` 真决策逐字段牙**（`frontend/scripts/render-l2.tsx`）：
  输入 = **今日现跑** `build\wisp.exe panel-assets -l2 fs.delete -irreversible delete -- "C:/Users/swq/Documents/notes.md"`
  （rc=0，`go build -o build/wisp.exe ./cmd/wisp` rc=0），其输出与
  `frontend/fixtures/l2-card-fs-delete.json` **归一化后逐字节一致**（仅行尾差异）。
  `npm run render:l2 -- fixtures/l2-card-fs-delete.json` rc=0，1 卡 21822 B HTML。
  断言全部由输入派生（`render-l2.tsx:71-102` `expectationsFor`，无被断言的字面量）⇒
  组件若忽略 props 即红。实测渲染树含：`fs.delete` ×3、路径 ×1、`R1` ×3、`R8` ×3、
  原因原文 `R1: 工具声明为下界（L1）; R8: 不可逆操作（永久删除）` ×1、
  调用链 `cli &gt; panel-assets &gt; fs.delete` ×1、`判定来自 原生风险评估` ×1、
  卡标题 `需要你的确认` ×1、`信息不足` ×0——同一份决策对象逐字段上卡。
- **牙 2 `render:l2` 未判负控（喂不到 ⇒ 显式未判）**：新 fixture
  `frontend/fixtures/l2-card-unjudged-reason.json`（`reasonKnown=false`、`reason=""`——
  这是 wire 自己文档化的缺原因形状，`internal/panel/approval.go:49-52`，不是假数据）。
  `npm run render:l2 -- fixtures/l2-card-unjudged-reason.json` rc=0：
  `信息不足` ×1、真原因文本 ×0（真没画）、`>L2<` 徽标仍在——不降级成看似安全的空位（Q-23）。
- 负向对照（`render-l2.tsx:135-144`）：上游 demo 问卷串 0 命中（不是组件库假 props 渲染）；卡数==pending 数。

### AC#3 诚实边界（留非实现者裁）

渲染证据是 `react-dom/server` 的 HTML 串（React 自己画的树），**不是真 WebView2 里的
layout/paint DOM**——真窗口归票 33，真事件流推送归票 35 的 host 侧（接收端
`panel-root.tsx` 已就位）。AC 框照旧不勾。

## AC#11 — 深色＋浅色渲染外壳＋样例屏，对比度逐对 ≥AA

### 仪器

`frontend/scripts/render-shell-77.tsx`（`npm run render:shell77`，照
`render-model-settings.tsx`/`render-l2.tsx` 的 vite --ssr 形状，零新依赖）：

1. 渲**生产根** `src/App.tsx`（外壳＝顶栏/侧栏/会话列/composer）＋一枚样例屏（L2 强确认卡，
   吃真决策 fixture），在 `data-theme="dark"` 与浅色（:root）各渲一遍，未判 fixture 也各渲一遍
   （4 棵树）；
2. 走渲染树标签栈提取**实际绘制的** text/bg 配对，逐对要求「被渲染出来」——重构悄悄删掉一对即红；
3. 配对颜色从 `tokens.generated.css` 解析（浅=`:root`、深=`[data-theme="dark"]`，别名跟链、
   8 位 hex alpha 按所在面合成），`--panel-alpha` 断言 100%；
4. WCAG 2 相对亮度实算比值，逐对判 AA（本壳无大字，全部按 4.5:1；本壳无 ≥24px/≥18.66px bold 文本，
   3:1 档无行）；禁用钮按 WCAG 1.4.3 豁免但仍合成实测并打印；
5. 判定规则：低于 AA 且未挂账 ⇒ rc=1 点名配对与比值；挂账配对一旦变绿 ⇒ rc=1 要求删账
   （账不能藏修复也不能静默累积）。

### token 小改（改生成器真相源 `frontend/scripts/gen-tokens.mjs`，重生成，`tokens:check` rc=0）

| token | 前 | 后 | 依据 |
|---|---|---|---|
| `--ink-3` 亮 | `#9a9da3` | `#6d727a` | 亮色对 surface/page 实测 2.72/2.61:1（<AA） |
| `--ink-3` 深 | `#6c6f75` | `#8b8f96` | 深色对 surface/page 实测 3.08/3.53:1（<AA） |
| `--red` 亮 | `#e3474c` | `#d13d43` | 白字对红钮与红徽标对白面实测 3.98:1（<AA） |
| `--red` 深 | `#ee5c61` | 不改 | 深色 `--red` 兼作徽标文字（4.68:1 过 AA），压深即击穿文字红——单 token 无解，见断跳 |

### 对比度逐对数值（最终态，`77-scaffold-batch2-contrast.txt` 全量 50 次测量）

亮色（浅）：

| 类别 | 配对 | 值 | 判（4.5:1） |
|---|---|---|---|
| 正文 | ink/surface `#1f2124 on #ffffff` | 16.14:1 | PASS |
| 正文 | ink/page `#1f2124 on #fafafb` | 15.47:1 | PASS |
| 正文 | ink/inset `#1f2124 on #f7f8f9` | 15.18:1 | PASS |
| 次要 | ink-2/surface `#62656b on #ffffff` | 5.84:1 | PASS |
| 标签 | ink-3/surface（改后 `#6d727a`） | 4.84:1 | PASS（改前 2.72:1 FAIL） |
| 标签 | ink-3/page（改后 `#6d727a`） | 4.64:1 | PASS（改前 2.61:1 FAIL） |
| 标签 | ink-2/field（Ctrl K 键帽） | 5.22:1 | PASS |
| 按钮 | white/red（拒绝钮，改后 `#d13d43`） | 4.71:1 | PASS（改前 3.98:1 FAIL） |
| 按钮 | ink/surface（outline 钮） | 16.14:1 | PASS |
| 标签 | warn/warn-soft（待审批胶囊） | 2.67:1 | **FAIL(<AA) 挂账断跳** |
| 标签(未判) | warn/surface（信息不足前缀） | 2.97:1 | **FAIL(<AA) 挂账断跳** |
| 标签 | red/surface（L2 徽标，改后 `#d13d43`） | 4.71:1 | PASS（改前 3.98:1 FAIL） |
| 禁用态 | ink-3@0.5/surface | 1.97:1 | 实测留档，WCAG 1.4.3 豁免 |

深色：

| 类别 | 配对 | 值 | 判（4.5:1） |
|---|---|---|---|
| 正文 | ink/surface `#f2f3f4 on #232427` | 13.97:1 | PASS |
| 正文 | ink/page `#f2f3f4 on #17181a` | 15.99:1 | PASS |
| 正文 | ink/inset `#f2f3f4 on #1f2022` | 14.68:1 | PASS |
| 次要 | ink-2/surface `#a5a8ad on #232427` | 6.51:1 | PASS |
| 标签 | ink-3/surface（改后 `#8b8f96`） | 4.78:1 | PASS（改前 3.08:1 FAIL） |
| 标签 | ink-3/page（改后 `#8b8f96`） | 5.47:1 | PASS（改前 3.53:1 FAIL） |
| 标签 | ink-2/field | 5.85:1 | PASS |
| 按钮 | white/red `#ffffff on #ee5c61` | 3.31:1 | **FAIL(<AA) 挂账断跳** |
| 按钮 | ink/surface | 13.97:1 | PASS |
| 标签 | warn/warn-soft（深色合成底 #41332a） | 5.15:1 | PASS |
| 标签(未判) | warn/surface | 6.59:1 | PASS |
| 标签 | red/surface `#ee5c61 on #232427` | 4.68:1 | PASS |
| 禁用态 | ink-3@0.5/surface | 2.24:1 | 实测留档，WCAG 1.4.3 豁免 |

### 任一对 <AA 的具名断跳（⛔ 不写整体达标；终态 47/50 过，3 对挂账）

1. **深色 拒绝钮 white on red = 3.31:1**：修复需按钮面专用 token（如 `--red-strong`）。
   新增彩色 token 会让四方对账（`internal/panel/tokens_fourway_test.go` check 2）未覆盖，
   须同步 `docs/evidence/s1/c21-native-tokens.md` 第四代未覆盖表＋该测试的钉名表
   `c21FourthGenUncovered`——两者均在本波写面（`frontend/**`＋`docs/evidence/s1/77-*`）之外 ⇒ 停，报断跳。
2. **亮色 待审批胶囊 warn on warn-soft = 2.67:1**：橙字在橙底浅 tint 上无解，
   需文本专用 token（如 `--warn-ink`，库已有 `--accent-ink` 先例）；同上断跳。
3. **亮色 信息不足前缀 warn on surface = 2.97:1**：同上，需 `--warn-ink`；断跳同上。

### AC#11 诚实边界

- 测的是**静态渲染树实际绘制的配对**（renderToStaticMarkup 的 class + 生成 token 值），
  不是浏览器合成像素；hover/focus 态不在静态树里，本轮未测（深色 `--red` 的 hover 半透明
  组合与 ink-3 on hover-2 属同类，留待真浏览器一轮）。
- 对比度仪器是本轮新牙，**未接 CI**（ci.yml 本波不在写面）；本地命令可复跑
  `npm --prefix frontend run render:shell77`。

## 门禁复读（本波全部改动的最终读数）

| 门 | 读数 |
|---|---|
| `npm --prefix frontend run typecheck` | rc=0 |
| `npm --prefix frontend run lint` | rc=0（36 warnings 0 errors，109 files；基线 36/108 ⇒ 新文件 0 新增） |
| `npm --prefix frontend run build` | rc=0（首屏 JS gzip 147.70 kB，与 AC#12 记录同水位；dist 已删回 `.gitkeep`，`git status -- frontend/dist` 空） |
| `npm --prefix frontend run tokens:check` | rc=0（97 keys，43 tokens，生成件=生成器） |
| `npm --prefix frontend run render:shell77` | rc=0（50 次测量，3 对挂账断跳具名） |
| `npm run render:l2`（真决策/未判两 fixture） | rc=0 / rc=0 |
| `go build -o build/wisp.exe ./cmd/wisp` | rc=0 |
| `go test ./internal/panel -run 'FourWay|FrontendHasNoEmoji|FrontendIsStateless|NeverNamesAnApprovalDecision|ColourLiterals' -count=1` | rc=0（四方对账在 token 改后仍绿） |
| `sh scripts/d22scan.sh` | rc=0（clean；ban#8 frontend/ 128 text files 含两枚新文件；ban#6 frontend/ 128） |

硬约束遵守：未动 thresholds/golden/allowlist/PLAN/specs/D43/design；未放宽任何既有断言
（render-l2 的断言一字未改，只加新仪器与新 fixture）；未 push；UI 零 emoji（TestFrontendHasNoEmoji
＋d22scan 双读数绿）；dist 零字节入库。
