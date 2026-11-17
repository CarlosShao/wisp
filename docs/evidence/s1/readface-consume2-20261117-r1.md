# readface-consume2-20261117-r1 — 前端消费波 2（写腿 agent-113 · 裁决表）

派单：`.scratch/wisp/probes/orch/20261117/txt/frontend-consume-wave2-dispatch.txt`（读后照办）。
前置：`3a62bed97`（票 36 AC#3＋38 AC#4 的 Go 半）已落两面；票 36/38 尾「编排者追加」节 file:line 现读核对。
写面：`frontend/src/**`、`frontend/scripts/render-*`、`internal/panel/composer_test.go`（ledger 行）、
`docs/evidence/s1/readface-consume2-*`。⛔ 契约/PLAN/specs/thresholds/golden/allowlist/D43/design 零字节；
⛔ 放宽断言（palette38 仪一处**随真相再裁**，见 §4③，钉得更紧不是放松）；⛔ push；⛔ dist 字节
（`git ls-files frontend/dist` ＝ 1，仍只有 `.gitkeep`，构建产物不入库）；⛔ UI emoji（修复过程见 §4②）；
⛔ 造数据（所有仪器喂的都是 Go 侧 wire 形自己的 fixture）。

## 0. 如实上报的三处偏差/裁量

1. **两枚页字面 ledger 注册物理上在 `bridge_test.go`，不在派单写的 `composer_test.go`。**
   派单第 3 条写「删 internal/panel/composer_test.go 里这两条 ledger 注册」，但现读盘面：
   history.page / roster.cancel 的**待页字面**注册在 `internal/panel/bridge_test.go` 的
   `rosterPageGapLedger194`（Go 半落地时的注册，`bridge.go:172-174/:209-210` 注释即指向那里）；
   `composer_test.go` 里挂号的只有 **TaskRowView 三键**（`composerTSGapLedger:244`）。两处都按
   ledger 双向规则**同 commit 摘除**：`bridge_test.go` 两条（map 归空）＋ `composer_test.go` 一条。
   顺带把 `inbound_roster_253_test.go:82-83` 一句**注释**（"registered pending"）更正为新真相——
   纯注释、零断言改动。
2. **`RunUsageRow` 翻译随走线新增进 `composerTSNameMap`（composer_test.go）。** TaskRowView 在
   panel.ts 声明 `costToDate` 后，契约走线（`nestedContractViolations`）第一次下潜进这个嵌套结构，
   红句自己点名「add the translation and the panel.ts interface in the same move」⇒ 同 commit 加
   `"RunUsageRow": "RunUsageRow"` ＋ panel.ts 的 `RunUsageRow` 接口（五键 tokensIn/tokensOut/
   cachedTokens/micros/currency，与 `run_feed_145.go:115-121` JSON tag 逐一对照）。这不是放宽——
   是走线**新到达**一格的必经动作；`RunSection.usage` 三键（running/usage/done）的注册照旧在册未动。
3. **`frontend/package.json` 加了两行仪器脚本**（render:historypage36 / render:taskcard38）。派单写面
   列了 `frontend/scripts/render-*` 而 package.json 不在其中；不加则编排者无法按 `npm run <名>` 复跑
   仪器（本波之前所有仪器都是这个形状）。两行都只挂仪器命令，零依赖变化。

## 1. 票 36 AC#3：history.page 画成 1000 行滚动列表

- **面板声明（lib/panel.ts）**：`requestHistoryPage(offset, limit)` 走既有 settings 封套
  （`window.wispDispatch`，method `history.page`，选择子键 `historyOffset/historyLimit`＝
  bridge.go JSON tag 逐字）；同 commit 摘 `bridge_test.go` 的待页注册（stale-red）。
- **回复解析（lib/history-view.ts）**：`parseHistoryPageReply` 按 privacy-view 同一纪律——
  拒绝句（`REFUSED_ENVELOPE_PREFIX`）逐字、prose 不解析、只认 `wisp-history-page:` 载荷行、
  rows/total/offset/limit 缺或型不对 ⇒ payload-error（未判，不画假行）。行形
  `{id,ts,kind,tool?,outcome?,summary?}`＝`dao_history_page_36.go HistoryPageRow` 逐字。
- **1000 行滚动＝两套精确机制**：
  - **数据分窗**：`mergeHistoryPage` 只接受**连续窗**（page.offset 必须＝已合并行数），重复 id 与
    窗口断档都记成**具名缺陷**（duplicates/anomalies）上屏告警——「不重不漏」是合并器的性质，
    不是口号；吸收或静默丢行都进不去。
  - **渲染虚拟化**：`visibleSlice(count, scrollTop, viewport, overscan)` 纯函数（行高固定
    `HISTORY_ROW_HEIGHT=26`），列表容器只渲染可视带＋overscan。
- **牙（新仪器 `npm run render:historypage36`，rc=0）**：1005 行喂三窗（500/500/5，窗首
  scroll-1004/scroll-0504/scroll-0004＝Go 测试自己的三窗首行）过**生产** parse＋merge ⇒ 恰 1005 行、
  id 全集各一次、零重复零断档；重复窗 ⇒ duplicates=5 且不重复上屏；跳窗 ⇒ anomaly 点名两半；
  虚拟切片自顶滚到底**并集恰全覆盖 1005 行**、顶位渲染藏尾行（虚拟化是真的）；total=0 ⇒
  「暂无记录」而**首答前不画空态**；拒绝/无载荷/坏载荷 ⇒ 逐字/未判，零假行零假计数。
- **行显示**：`{ts,kind,tool,outcome,summary}` 逐字（ts 走 `formatHistoryTs`，0＝存储自己的
  「无时刻」而非 1970；outcome 原词上屏不判色）。
- **delete-one/purge/export 走既有 privacy 面（⛔新造第二套）**：浏览器内嵌的隐私条**整块复用
  wave-1 页半**——`usePrivacyRequest`＋`privacyRunners`＋`PrivacyReplyBlock`（privacy-view /
  privacy-settings 原件），purge 只从确认卡确认键发出（一处 `purge.submit(`）、delete-one 逐行
  发（域词映射＝读面自己的两表名 task_log/tool_call，`historyDomainOf`；门仍按 PrivacyDomains()
  判，拒绝逐字上屏）。页面上**没有任何第二套破坏性方法名**（仪器源级钉）。
- **接线**：App 顶栏新增「历史记录」钮（lucide `ScrollText`；`History` 图标本仓 lucide-react 无此
  导出，现查）开 `HistoryBrowser` 对话框（palette 同形：backdrop 点击/Escape 关）；`load` 只喂
  `requestHistoryPage`。侧栏 historyGroups 保持诚实空——那是**会话树形状**，读面不声明会话分组，
  注释已改写说明（历史平铺列表在顶栏浏览器里）。

## 2. 票 38 AC#4 后半：TaskRowView 三键＋roster.cancel 真门＋applied 读数

- **三键上屏**：panel.ts `TaskRowView` 补 `startedAt?: number / costToDate?: RunUsageRow /
  currentStep?: string`（Go 侧全 omitempty ⇒ 页侧全 optional；缺键＝宿主没盖戳 ⇒ 显示**未判**，
  ⛔ 假 0/假 epoch）；`App.monitorAgents` 三键 1:1 拷贝；`MonitorPopover` 每行加任务卡行
  （`TaskCardLine`）：开始（生产 formatter）／成本（micros+currency 成对逐字，tokens 数逐字）／
  当前步（缺＝线自己的「无在飞」语义，非未判）。
- **roster.cancel 走真门（先确认再取消 ⛔一键杀）**：`MonitorPopover` 新 prop `cancelDoors`
  （App 绑 `requestRosterCancel`＋`requestHistoryPage`）；每行「停止」键**只开确认卡**
  （alertdialog：点名行 id、roster.cancel、不可从面板撤回、只停一行不级联），请求**只从卡上确认键**
  发出（源级钉：`confirmStop(` 只有一个调用位）。回执/拒绝/无宿主三态逐字上屏
  （`TaskCancelReply`）；未绑门（showcase/演示形）不画停止键——不画死按钮。
- **取消后 applied-steps 读数（终态名/已落地步数/被中止步数）**：读数从**取证面**来，不造——
  Go 半的裁决表写明结算面＝tool_call 取证行＋终态 task_log 行＋名册行；页侧经 history.page 读：
  - `scanTaskForensics`：按存储自报 total **整史扫**（500/窗）捞该任务的终态行（kind=task 且
    id=taskId）与工具取证行（kind=tool 且 id 前缀 `<taskId>#`——loop.go:603 自己的 correlation
    拼法，与名册 blockedOnApproval join 同族）；
  - `appliedStepsReading`：终态名＝名册行 status（statusKnown）否则终态 task_log 行 outcome；
    已落地＝success 行数；被中止＝cancelled 行数；**扫描不完整 ⇒ null ⇒ 未判**，⛔ 半截计数假报。
  - 读数行（`AppliedReadingLine`）三槽各自「有值显值、缺显未判」。
  - 牙（新仪器 `npm run render:taskcard38`，rc=0）：三键喂值⇒逐字（开始＝formatter 输出、
    123456 micros CNY、tokens 1200/340、缓存 512、fs.write）；缺键⇒未判/无在飞、无 1970 无假
    micros；取证行喂 child-9（1 success＋1 cancelled＋终态 Stopped）⇒ 终态 Stopped · 已落地 1 步 ·
    被中止 1 步，外人行不串数；1005 行整史扫描（三窗）⇒ complete 且读数同值；中途拒绝 ⇒
    complete=false ⇒ 未判；确认门源级钉＋回执逐字＋零 emoji。
- **命名**：照 Go 半编排者裁定落地名 `roster.cancel`（票面 `task.cancel` 被 internal/tools 票 188
  名轴仪钉死，撤销口令在 `bridge.go MethodRosterCancel` 注释），页侧字面与选择子键 `rosterTaskId`
  同名同源。

## 3. 读数（收件三把尺，全部现跑）

| 尺 | 读数 |
|---|---|
| `npm --prefix frontend run build` | **rc=0**（1.08s，chunk>500kB 警告为既有） |
| `npm --prefix frontend run typecheck` | **rc=0** |
| `npm --prefix frontend run lint` | **rc=0**（45 warnings 全为既有件，0 errors） |
| 新仪器 `render:historypage36` / `render:taskcard38` | **rc=0 / rc=0**（OK 句见各仪器尾行） |
| 相关仪器回归 | roster / history36 / privacy40 / palette38 / palette38kbd / nav / shell77 / approval37 / stream / composer / receiver / model-settings / security40 全 **rc=0** |
| `go test ./internal/panel` | **ok 23.4s rc=0**（含 frontend_hygiene 零 emoji 门、契约双向走线、名册 13/9 闭合尺、页字面镜像门） |
| `go build ./...` | **rc=0**；`go vet ./cmd/wisp` rc=0 |
| `sh scripts/d22scan.sh` | **OK - control rc=0, scan rc=0**（末行），0 findings；ban #8 frontend/ 149 件 examined |
| 附加（载具） | `scripts/wisp-cli-staged-tests.sh` 三发定向：TestHistoryPage* ／TestTaskCardFields*／TestRosterCancelDoor* 全 **PASS**（含真子代理取消 2.86s）。裸 `go test ./cmd/wisp` 本机 0xc0000135＝载具头注写明的已知形状（语音 DLL 只随 exe 目录绑定，裸跑拿不到），非本腿产物，走载具即绿 |

## 4. 过程中修的三处红（无一放宽）

1. **TS6133**：history-page.tsx 曾 import `requestPrivacyDeleteOne` 但实际调用走
   `privacyRunners.deleteOne`（wave-1 绑定表）⇒ 摘 import，仪器改为源级钉「绑定表里的
   deleteOne＝requestPrivacyDeleteOne」——同一事实，钉点更准。
2. **零 emoji 门红 4 处**：我初稿在注释/断言里用了 `⛔`（U+26D4，d22scan/字号仪 2600–27BF 带）
   ——本仓 frontend 面**注释也不豁免**（internal/panel 的 frontend_hygiene_test.go 现跑抓出）⇒
   四处改纯文字（「禁止一键杀」「禁止新造第二套」）。改后门绿、d22scan clean。
3. **render:palette38 一枚旧钉随真相再裁**：该仪源级钉「取消任务 未接线＋点名 composer_dispatch.go/
   task.cancel」——roster.cancel 落地后这钉钉的是**旧假话**（与票 38 第二批 render:security40 同形）。
   再裁（未松）：钉改「取消任务 已接线＋点名 roster.cancel＋App 必须绑
   `cancelDoors={{ cancelTask: requestRosterCancel, historyLoad: requestHistoryPage }}`」——
   取消门本体逐齿钉在 render:taskcard38。改后 palette38 rc=0。
   注：`render:config-screen` 脚本名在本仓不存在（我误列），battery 里那枚 rc=1 是 Missing script，
   非仪器红。

## 5. 诚实缺口（具名，不勾框）

- **AC 勾框一律留给非实现者**（票 36 AC#3、票 38 AC#4 的格子本腿不勾）。
- **1000 行 60fps**：仪器钉的是虚拟化数学与数据不重不漏（SSR/纯函数可证的全部）；真实 60fps 需要
  真宿主窗滚动采样，本腿无桌面独占环境，归桌面验收。
- **子代理行 costToDate 恒未判**：Go 半写明子代理 spend 在其 finalize 落账、EvUsage 不过本 sink
  ——这是线自己的语义，页面照显未判，不是缺陷。
- **applied-steps 整史扫描的代价**：一次取消后按 total 全量分页读历史（500/窗）；自用库规模可接受，
  若日后 history.face 长出逐任务过滤参数，此扫描可一行收窄。
- **取消读数的终态名依赖快照推送节奏**：回执即回（停止已发出），终态/取证行由任务自行结算后经
  下一次 snapshot push ＋取证扫描可见——页侧不等待、不猜。
- **palette 命令执行/模式强制**仍未接线（D11③ DEFERRED、名册无执行通路），与本腿无关，钉句照旧。

## 6. 撞车账

本腿写面独占 `frontend/src/**`＋`internal/panel` 本波写额；起跑时共享树有其他腿的未提交件
（`.scratch/**`、`fixtures/composer-states.html` 等，均不在本腿 pathspec）。本 commit 显式 pathspec
只含上列 13 件，未碰他人未提交活。
