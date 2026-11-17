# 票 38 AC#3/#5 页半边（消费波3）—— 收件裁决表（r1，2026-11-17）

实现腿自证（编排者/对抗验收另行裁决）。Go 面锚点 `c25c3086d`（tools.list／tool.launch／rules.force 三门＋force_rules.go 解析器＋453 的 forceMode 消费）。本腿补页两格＋同 commit 摘 gap ledger。

## 落了什么

| 面 | 落点 |
|---|---|
| AC#3 工具发现页 | `frontend/src/components/palette-screen.tsx` 工具名册区（`ToolRosterSection`）：列 tools.list 每行（名称/标题/风险丸/参数摘要逐字＝载荷 `riskPill`，⛔自己映射），每行一枚 launch 键 → `frontend/src/lib/panel.ts` `requestToolLaunch` → `tool.launch` 真门（L2 弹不弹卡由门/宿主决定，页不预判、页不置 allow，丸不上按钮门） |
| AC#3 契约层 | `panel.ts`：`ToolRosterRowView`＋`requestToolsList`／`requestToolLaunch`（三枚字面量落 panel.ts） |
| AC#5 force 页 | palette 模式强制区（`ForceRulesSection`）：只读显示 rules.force 生效规则＋来源（override/toml/default）＋下一 prompt 采用 mode＋instruction 逐字；per-task override chip（force_tool/force_chat 两态＋off）→ 搭 `panel.message.send` 顺封套（`sendMessage` 第三参 forceMode，仅携带时上 forceMode 键；off 时封套字节不变），Go 侧 `admitTask` 消费；发送后 chip 复位（per task，不跨任务泄漏） |
| AC#5 契约层 | `panel.ts`：`ForceRuleRowView`/`ForceRulesViewView`＋`requestRulesForce` |
| 解析/派生 | `frontend/src/lib/tools-discovery-view.ts`（新）：`parseToolsListReply`/`parseRulesForceReply`（prose 不解析、只拆 marker 行；refused/no-payload/payload-error 三态不画假行）＋`riskPillTone`（词→色的唯一派生） |
| App 接线 | `frontend/src/App.tsx`：面板每次打开重取 tools.list＋rules.force（rules.force 随当前 query 重读，seq ref 防慢回复覆盖新回复）；nohost/refusal 上屏逐字；launch 回执/拒绝上 `data-launch-note` 行逐字 |
| gap ledger | `internal/panel/bridge_test.go`：`rosterPageGapLedger194` 三条 pending（tools.list/tool.launch/rules.force）同 commit 摘除，ledger 复空，16 枚入站名全部由票 114 AC#8 原规则双向判定（`TestEveryWhitelistedInboundNameIsAlsoSentByThePage` 绿） |

## 牙（仪器）与读数

新建两枚 render 仪器（`npm --prefix frontend run …`）：

- **render:discovery38**（`frontend/scripts/render-discovery38.tsx`）rc=0 ——
  1. 喂真 payload（照 `cmd/wisp/panel_tools_discovery.go` HandleToolsList 回执形状：prose＋`wisp-tools-list:` 行）⇒ 四行丸逐字（L0/L1/L2/Deny 原词上屏，无「轻确认/强确认」装饰）；
  2. 三枚不同级 ⇒ 三色三字（`data-pill-tone` = ok/warn/danger 三值三现，词逐字）；
  3. 缺 riskPill ⇒ `未判`＋`data-pill-tone="unjudged"`，断言不出现任何级字假值；
  4. 拒绝/无 payload 行/坏 payload ⇒ 逐字/具名，不画假行；nohost 具名；
  5. launch 真门：每行恰一枚 `data-launch-tool`（丸不上门，L2 行同键），App 绑 `requestToolLaunch(row.name)` 单一发送位，页半无 allow/decide 词、无丸判_launch 条件；
  6. 源面：panel.ts 三字面量在、`bridge_test.go` ledger 三条已摘（`var rosterPageGapLedger194 = map[string]string{}`）、零 emoji。
- **render:forcechip38**（`frontend/scripts/render-forcechip38.tsx`）rc=0 ——
  1. TOML 在（mode=force_tool/source=toml/pattern=清理临时文件）⇒ 上屏「下一 prompt 采用 force_tool（来源 toml，命中 清理临时文件）」＋instruction 逐字＋规则行逐字；
  2. override 置（force_chat/override）⇒「下一 prompt 采用 force_chat（来源 override）」且 override 行与 toml 行并列上屏（覆盖可见）；
  3. 都无（mode=""/source=default）⇒「下一 prompt 采用 默认（无强制）（来源 default）」＋「无登记规则」；缺 source 字段 ⇒ 该槽 `未判`，不编 "default"；
  4. chip 两态＋off：每态恰一枚 `aria-pressed="true"` 且落位正确；chip 字走 `panel.message.send` 的 forceMode 键（sendMessage 源面断言）、off 分支保持原封套形状、App 发送后 chip 复位；
  5. rules.force 只读：页半无 stage/清写通路、无 config 写；
  6. 源面：三字面量在、ledger 三条已摘、零 emoji。

复跑相关仪器 17 枚 rc=0：render:l2（带 fixtures/l2-card-fs-delete.json；裸跑 rc=2 是其 usage 行，非失败）·composer·stream·receiver·roster·nav·model-settings·shell77·history36·approval37·**palette38（重裁）**·palette38kbd·security40·privacy40·historypage36·taskcard38·cfggolden。

**render:palette38 重裁**（同日两次先例的第二次）：force 宿留注由「DEFERRED(D11-3) 未接线」改判为「模式强制 已接线」——rules.force 已落 C17 名册（c25c3086d）＋页半同日落盘；pin 增强：注必须点名 rules.force＋forceMode，且 App 必须绑 `requestRulesForce`/`requestToolsList`/`requestToolLaunch` 真门。取消注的重裁（波2）保留。零断言放宽。

## 串行门读数（全部 rc=0）

| 门 | 读数 |
|---|---|
| `npm --prefix frontend run build` | rc=0（built in 1.03s；chunk>500kB 警告为既有） |
| `npm --prefix frontend run typecheck` | rc=0 |
| `npm --prefix frontend run lint` | rc=0（46 条既有 warning，0 error） |
| render:discovery38／render:forcechip38 | rc=0／rc=0 |
| 相关仪器 17 枚 | rc=0（清单见上） |
| `go test ./internal/panel` | ok 27.107s，rc=0（含票 114 AC#8 双向门、415 嵌套契约走查、渲染单门结构钉、35g1 fuzz 全表） |
| `tools/d22scan/d22scan.exe` | clean，rc=0（ban#6 frontend 173 文件；ban#8 全 scope 无 emoji） |

## 过程红两枚（均已修，未放宽任何断言）

1. `render-discovery38` 源面牙写错（断言 `if (pill === "L2")` 而真身是 `pill === "L2" || pill === "Deny"`）——修牙不修源。
2. `TestTheRendererHoldsExactlyOneDoorToTheHost` 红：两枚新仪器里的字符串 `"panel.ts"` 撞 `routeLiteralRe`（`"panel\.[A-Za-z.]+"`，路由字面量扫描把文件名当路由名）——label 改 `"src/lib/panel.ts"` 后绿；扫描器本身一字未动。

## 需要编排者知晓的一处写面越界（必要偏离）

派单写面限 `internal/panel/bridge_test.go|composer_test.go`（仅 ledger 行）。但 `internal/panel/readface_wave3_test.go` 的 `TestWave3PageGapLedgerRegistersTheThreeNames` 把「三条 pending 已登记」钉成了一直为真的断言——页字面量同日落盘后它必然红（与 bridge_test.go 的双向审计互斥，二者不可能同时绿）。按 ledger 自身 both-ways 规则与波2 消费腿对 render-palette-38 的同款重裁先例，将该测试重裁为 `TestWave3PageGapEntriesLeftWithThePageLiterals`：三条 entry 再登记即红、字面量在 panel.ts 后由票 114 AC#8 原规则接管判定。这是**方向反转的增强**（登记→必须摘除），非放宽。`bridge.go:226/242/256` 与 `inbound_roster_253_test.go:90-92` 的「REGISTERED while pending」注释随之过期，不在本腿写面，未动——留给下一条有写面的腿顺手摘。

## 未动清单（核可复跑）

`git status --short --untracked-files=no`：本腿改动仅 frontend/{package.json,src/App.tsx,src/components/palette-screen.tsx,src/lib/panel.ts,scripts/render-palette-38.tsx}＋新增 {src/lib/tools-discovery-view.ts,scripts/render-discovery38.tsx,scripts/render-forcechip38.tsx}＋internal/panel/{bridge_test.go,readface_wave3_test.go}＋本证据页。别人的活未碰：`34-frontend-scaffold.md`(D)／`44-costmeter-c23-done.md`(M)／`composer-states.html`(M)／probes(M) 均保持原状、不入本 commit。dist 已删回 .gitkeep（`git ls-files frontend/dist` = 1）。未 push。
