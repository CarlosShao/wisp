# 36 (leg 1, agent-95) — Result & history panel：接线读数、行为牙、门禁 rc

- 时刻：2026-11-16 22:1x +08；起手锚 `3d3ed8c08`；写面＝`frontend/**`（铁律 4）＋本件＋票面 Progress log。
- commit：见回话第一行（本件与票面追记、frontend 件同 commit，显式 pathspec，未 push）。

## 0. 一句话

票 36 的**工作状态面**真接上了：快照 `run` 键（票 145 的生产读面）第一次有 UI 读者；结果流第一次过
**F2 白名单消毒**（react-markdown + remark-gfm + rehype-sanitize，收紧 schema）；git 分支 chip 从
**喂死空对象**改成**真读 composer.git**。合同尺注册的 6 枚 Go-only 键**没有在前端造值**，按派单留给
Go 腿（见 §4）。

## 1. 逐格 AC 处置（票面 :37-47）

| AC | 处置 | 牙 |
|---|---|---|
| AC#0(31) 14-state 视觉清单 + 截图签收 | **留白**。本腿按 PLAN §17.5 表把**本 wire 实际携带的态**画了（工具 chip 五态 outcome 字形、错误卡 danger 竖条、Stuck 卡 warn 竖条、推理折叠行、SSE 光标、L2 卡已有 render-l2 钉着）。**等用户签收的截图格不由实现腿勾**（票面 09-26 追记同规则）；且 14 态里的 审批等待队列深度角标 / L1 阻止窗口 / L2 原生降级卡 / 注入检出 / 已取消分隔线 / 思考中光带+计时 / token 计数 的 wire 字段不存在或属他票（37/145 的 pending[].position 等） | 本件 §2 缺字段清单 |
| AC#1 XSS 套件 | **做了**（功能层）：四族载荷 `<img onerror>` / `<script>` / `javascript:` href / 事件处理器属性 走产码 `MarkdownText`，标记无宿主碎片；`dangerouslySetInnerHTML`＝0（去注释扫描）。**CSP report-only listener 的浏览器运行监听留白**——那需要 WebView 宿主起真窗（票 33 的场），本腿无宿主 | `scripts/render-history-36.tsx` §1 |
| AC#2 sanitizer 白名单落盘 | **做了**：TIGHT_SCHEMA＝default schema 减 img/div/input、href 协议收窄 http/https/mailto、class 全剥；白名单即产码本体（组件里，不是副本） | `src/components/result-history.tsx:86-106` |
| AC#3 dangerouslySetInnerHTML grep=0 | **做了**：牙里去注释扫 `frontend/src`＝0 命中 | `render-history-36.tsx` §1 末 |
| AC#4 History 1000 行 60fps + delete/purge/export 对 DB | **留白（具名缺引擎）**：快照无任务历史 section（internal/panel 无 history reader；SQLite task/transcript 无 panel 读面；`panel.view.request` 被 Q-50=甲 封死）。**⛔ 不在前端造 1000 行假历史**。虚拟化、删除/清空/导出三钮都无从喂真数据 | App.tsx:359-380 注释具名 |
| AC#5 200-chunk golden 回放 token 计数 | **留白（半格做了）**：流式渲染面真接上了（done=false 走 RevealText+光标；结果流有界截断不合并＝票 197 R2 裁定，牙里断言截断句），但 **token 计数字段未上线**（ResultChunkView 无 token 字段；usage 维度在合同尺注册清单里） | result-history.tsx:393-413 |
| AC#6 隐藏面板动画暂停 (rAF) | **留白**：本腿新增的循环动画为零（光标呼吸/推理 chevron 全是既有 PLAN 冻结行，归 render-stream 仪器管）；rAF 计数器牙需要真宿主窗口 | — |

## 2. 本腿接上的 dormancy（有生产读者）

| dormancy | 之前 | 现在 |
|---|---|---|
| `snapshot.run`（票 145 工作状态） | **零 UI 读者**（grep `run\?\.` 只有 panel.ts 声明） | `ResultHistory`/`RunRegion` 渲染 tools/reasoning/error/stuck；读者＝`src/components/result-history.tsx`，喂面＝`cmd/wisp/run.go:1184/1281`（RunFeed 装配）＋`run.go:2729-2792`（consoleSink 四支路喂值） |
| git 分支 chip | `App.tsx` 喂死空对象（`view={{current:"",isRepo:false,...}}`），**生产推送的 git 维度全灭** | 真读 `composer.git`（票 181/190 读面，`internal/panel/git.go ReadGitForWorkspace` 填）；无 git 键时显式 `GitBranchUnread` |
| 结果流消毒 | App 直接 `chunk.text` 裸文本节点（流式安全但完成态无 markdown/无消毒） | 完成态过白名单；流式态保持文本节点+冻结光标 |

## 3. 保持「未接线」并具名缺哪台引擎（⛔ 不改成即时）

1. **30 天历史树/会话列表**（`App.tsx:361-380` 注释）：缺 SQLite 历史 reader（internal/panel 无读面）＋
   `panel.view.request` 被 Q-50=甲 封死（重开 C17 白名单要人工批）。
2. **ContextMeter 用量行**（`composer-row.tsx:293-301` 注释）：`RunSection.usage` 在合同尺注册清单
   （`internal/panel/composer_test.go:233`），Go wire 有值但 panel.ts ⛔ 声明（声明＝stale 红，除非
   Go 腿同 commit 摘注册）。仪表保持「还没收到数据」真话句。
3. **expanded 工具块的 argsJson 着色 JSON / L0-L2 胶囊**（`result-history.tsx` ToolChip 注释）：
   同上，注册键 `composer_test.go:231`。
4. **error 卡 retryable**（牙里断言未判）：注册键，未判句显式落场，不猜布尔。
5. **stuck 三钮（继续/换个方式/放弃）**：C17 入站面（Q-76），wire 无路由——不画死钮。
6. **usage/stop 终止词/done 维度**：注册键，未渲染，具名在 result-history.tsx 头注。

## 4. 需要 Go 侧做的事（先停下写进回话，⛔ 硬造）

- **Go 腿一波内做**（internal/panel 写面归他腿）：`composerTSGapLedger`（`internal/panel/composer_test.go:229-235`）
  摘 `ToolRunRowView{argsJson,level}`、`RunErrorRowView{retryable}`、`RunSection{running,usage,done}` 三行，
  同 commit 由前端腿补 `panel.ts` 三接口（`ToolRunningRowView/RunUsageRowView/RunDoneRowView`，形状照
  `docs/evidence/s1/415-r145-2-nested-contract-ruler-impl-v1.md` §AC#5 清单）。摘完前端即可渲染
  取证参数/风险胶囊/用量行/终止词——**不用改任何 Go 结构体**，只是解封页侧声明。
- **常驻免管道装配若要 run 键**（cmd/wisp/panel_resident_pump.go:39-46 现状＝无 RunFeed）：那是新读面，
  属 cmd/wisp 波次调度，本腿未动、未要求。

## 5. 新牙

- **`npm run render:history36`**（`scripts/render-history-36.tsx`）：
  - §1 四族 XSS 载荷过产码消毒 → 标记零宿主碎片；`dangerouslySetInnerHTML` 去注释扫＝0；
  - §2 喂真形状数据（按 `run_feed_145.go` JSON tag 抄形）⇒ 同值回显：5 行工具（含 unknown outcome
    的第四态）、toolsDropped 截断句、错误三件套、stuck 句、推理行；outcome→字形四态映射逐行断言；
  - §3 注册键**零造值**（argsJson/retryable/usage 不上屏；无 L0-L2 胶囊、无运行中 chip）；
  - §4 缺 reader 形态具名（「工作状态未接线：… PumpSources.Run 读取器未接」）；空 run／空 results 双显式；
  - §5 全面板走线：喂 composer.git ⇒ header chip 显真分支（`aria-label="当前分支 dev"`）。
- 读数：**OK，0 失败**（详见回话第一行 rc）。

## 6. 门禁读数（全部现跑）

| 门 | rc | 备注 |
|---|---|---|
| `npm run build` | 0 | 首屏 chunk 500.87KB→678.99KB（+54KB gzip）＝PLAN §17.6 冻结栈（react-markdown+remark-gfm+rehype-sanitize）进主 chunk 的代价；不做懒加载——每会话第一条完成消息就要走消毒路径 |
| `npm run typecheck` | 0 | |
| `npm run lint` | 0 | 0 errors / 37 warnings（全部 ai-native vendored 存量，本腿零新增） |
| `npm run tokens:check` | 0 | 97 keys / 43 tokens 未动 |
| `render:nav / render:stream / render:shell77 / render:receiver / render:roster / render:composer / render:model-settings` | 全 0 | App 形状尺、光标尺、对比度尺（47 对）全绿；KNOWN_BREAKS 3 枚照旧注册 |
| `render:l2` + fixture | 0 | 裸跑 rc=2 是它自己的 usage 行，非失败 |
| `go test ./internal/panel -run 'Composer/Approval/PlantedNested/前端卫生五连'` | ok | 合同尺 0 红＝panel.ts 零越权声明被尺确认 |
| `sh scripts/d22scan.sh` | OK | frontend/ 130 件 ban#8 零命中 |

## 7. 诚实边界

- 14 态**没有全画**（wire 没有的字段不造）；本腿画的是 §2 表列出的、`run` 键真实携带的维度。
- AC#1 的 CSP listener 与 AC#4 的 1000 行 DB 对账需要真宿主/真 reader，**留白不是跳过**。
- 首屏 +54KB gzip 是真代价，写在 §6；不掩盖。
- dist 构建产物已按派单删回 `.gitkeep`（`frontend/.gitignore:12-13` 本来就 ignore，git status 干净）。
