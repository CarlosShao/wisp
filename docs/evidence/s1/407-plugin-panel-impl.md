# 票 407 实现记录 — 插件面板（清单表 + 四枚牙 + 仪器）

**日期**：2026-11-19 17:2x+08。**来路**：派单 `.scratch/wisp/probes/orch/20261119/txt/407-plugin-panel-dispatch.txt`（本编队前端写腿）。
**性质**：实现腿自录证据（裁决/对抗验收由非实现者另出）。

## 1. 交付面（与派单写面逐项对齐）

| 件 | 变更 |
|---|---|
| `frontend/src/components/plugins-screen.tsx` | 共享清单体 `PluginInventoryList` 重切为票 407 的**清单表**（`<table>`，列＝id／启用／生效／名称／版本／来源/锁定说明／能力／网络白名单／宿主面）；空清单显式「还没有插件」零表格；缺键格＝「未知」（cell 助手 `enabledWord/stateWord/textCell/noteCell/listCell/idCell`，全部 unknown-typed 守卫）；缺 enabled 的行**不画开关**（猜方向＝猜写入）。`PluginsScreen` 外壳不动（缺键等待句、tier2 芯片、反馈线）。 |
| `frontend/src/lib/config-feedback.ts` | 新增 `translatePluginToggleReply`：插件开关的回执**原样上屏**（refusal 原话、已接受回执的诚实注记、`需要重启` 逐字保留），机器行（`wisp-config-reload:`，复用 `CONFIG_RELOAD_MARKER` 常量）剥掉不进屏；共享的 `translateTieredWriteReply`/`REFUSALS` 一字未动（其余写路零行为变化）。 |
| `frontend/src/lib/config-actions.ts` | `performPluginToggleWrite` 换用上述翻译器；写路本身一字未动：仍 `requestConfigWrite` → `method:"config.set"`＋`configField:"plugin_enabled"`＋`configPlugin:<id>`（panel.ts 现成信封）。 |
| `frontend/scripts/render-pluginpage407.tsx` | 新仪器（npm script `render:pluginpage407`）：两张静态 markup（空态＋有插件态）写盘＋四枚牙。 |
| `frontend/package.json` | 加 `render:pluginpage407` script（其余不动）。 |
| `docs/evidence/s1/407-plugin-panel-markup-empty.html` / `-with-plugins.html` | 仪器产出的两张静态 markup（最终代码产出，1606/8445 字节）。 |

**没动的**：`cmd/wisp/**`、`internal/panel/**`、`internal/plugin/**`（读面/写面/门在票 215a 已齐，零缺口）、`frontend/fixtures/composer-states.html`（他人在飞件）、panel.ts（声明已与载体对齐，仪器双向尺钉住）。

## 2. 显示规则（AC#1/#3 的落点，组件头注释同文）

- **只许用现成读面**：每格＝`PluginRowView`/`PluginsSection` 的 wire 值原样（name/version/note/enabled/state/capabilities/netAllowlist/hostApi/tier2Enabled/engineState/skillState），两态句（engineState/skillState）逐字。
- **未知规则**：键没上线（`typeof`/`Array.isArray` 守卫不成立）⇒ 该格「未知」，⛔ 编默认。name/version 是非 omitempty——`""` 是**给了的实测值**（盘上无安装），照本文件既有空档词「（无）」显示，**不是**「未知」；note 是 omitempty——缺席是常态形状且页侧无法区分"干净"与"没读到"，按票面规则答「未知」。
- **空清单**：`rows.length===0` ⇒ 「还没有插件」＋原实测说明句，零 `<table>/<ul>/<li>`（变异 M5 证明这枚尺咬得住）。
- **缺键行无开关**：enabled 缺键 ⇒ 启用格「未知」且不画 开启/关闭 按钮（按钮方向来自当前态，缺态画钮＝编默认）。

## 3. 写面（AC#2/#5-③ 的落点）

- 开关走现成门：`performPluginToggleWrite` → `requestConfigWrite`（`config.set`，`plugin_enabled`＋`configPlugin`）——⛔ 自造第二条写路；页/门源码零 `wispDispatch`/`requestHostAnswer` 直调。
- **拒绝句原话上屏**：`translatePluginToggleReply` 对 `被拒绝` 回执逐字上屏（含 `config_handlers.go:575-579` 的「必须指名插件 id（configPlugin＝[plugins.<id>] 的 <id>…）」与锁定段放宽拒绝）；已接受回执的诚实注记（「配置层已生效；引擎未接线…面板侧只能关不能开」）与 `需要重启` 逐字保留；机器 reload 事件行剥除。修掉的旧形状：原 `translateTieredWriteReply` 把这类拒绝折成通用「保存失败」（吞句）。

## 4. 仪器（AC#4）与四枚牙（AC#5）

`npm run render:pluginpage407` ＝ vite SSR build + node 检查；空态/有插件态两张 markup 写 `docs/evidence/s1/407-plugin-panel-*.html` 并逐项断言（逐字回传、未知格、零 emoji、零表格空态）。四枚牙均可红：

| 牙 | 变异（临时改源→跑仪器→还原，盘上字节已复原） | 读数 |
|---|---|---|
| ①引用不存在字段⇒红 | 页面 id 格插 `{String(row.author)}`（载体无此字段） | `FAIL tooth 1: every row.<field> ... [author]`，rc=1 |
| ②开关没走 config.set⇒红 | 写字段改自造名 `"plugin_on"` | `FAIL tooth 2: the write is plugin_enabled addressed by the plugin id selector`，rc=1 |
| ③拒绝句被吞⇒红 | `performPluginToggleWrite` 换回 `translateTieredWriteReply`（通用摘要器） | `FAIL tooth 3 (source): ... its OWN verbatim translator`，rc=1 |
| ④缺字段当默认⇒红 | `enabledWord` 改成 `true?已启用:已停用`（缺键落 false 默认） | `FAIL tooth 4 (AC#3 mutation 1): an ABSENT enabled must render 未知, not the false-default 已停用`，rc=1 |
| AC#3-②空清单画表⇒红 | 表格分支条件改恒真（空清单也画表） | `FAIL the empty inventory must say 还没有插件`＋`ZERO table markup`＋`screen empty`，rc=1（3 检查） |

每次变异后均已还原并复跑仪器回绿（rc=0）。还原核验：`git diff --stat` 仅含本节第 1 表的四枚产品件。

牙①的形状尺是**活的**：仪器运行时读 `internal/panel/plugins_inventory_215.go`，正则抽 `PluginRow`/`PluginInventorySection` 的 json tag 集，双向核对 ①页面所有 `row.<f>`/`section.<f>` 引用 ⊆ tag 集 ②panel.ts 两个 interface 的键集 == tag 集（抓到过我自己正则漏数字的 bug：`tier2Enabled` 的 `2` 使首版 `[A-Za-z]+` 漏配——尺先红了仪器作者）。

## 5. 现跑读数（2026-11-19 17:1x–17:2x+08，全部本腿现跑）

| 闸 | 命令 | 读数 |
|---|---|---|
| 仪器 | `npm run render:pluginpage407` | **rc=0**，OK 句见仪器末行；两张 markup 已写盘 |
| 前端 typecheck | `npm run typecheck`（tsc -b） | **rc=0** |
| 前端 lint | `npm run lint`（oxlint） | 0 errors（50 warnings＝存量，非本腿引入） |
| 前端 build | `npm run build` | **rc=0**（1.04s；chunk>500kB 警告＝存量） |
| 共享面回归 | `npm run render:security40`（挂 SkillsSettings→本清单） | **rc=0**（其插件 pin 全存活） |
| 共享面回归 | `npm run render:model-settings`（读 config-actions/panel.ts 源尺） | **rc=0** |
| Go 门/hygiene | `go test ./internal/panel -count=1 -run 'Test215A\|TestFrontendHasNoEmoji\|TestPanelFrontendIsStateless\|TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme'` | **ok ... rc=0**（含写面三针 Test215A×3＋零 emoji＋无状态＋色字面量） |
| 载具 | `D:\work\soft\Git\bin\bash.exe scripts/wisp-cli-staged-tests.sh '-test.run' '^Test215'`（跑前 `Get-Process wisp.test`＝0） | **rc=0，PASS=9 SKIP=0 FAIL=0**，`ok github.com/CarlosShao/wisp/cmd/wisp 1.00s`，log＝`.scratch/wisp/probes/orch/20261119/407/vehicle-test215.log`（载具形：DLL staged sha256 ok，fixtures present 无跳过） |
| 全树编译 | `go build ./...` | **rc=0** |
| 扫描 | `sh scripts/d22scan.sh` | **clean - no D22 ban violations ... d22scan.sh: OK - control rc=0, scan rc=0**（末行照录：`d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=422, bans #1-5 cmd/=107, ban #6 frontend/=188, ban #7 internal/tools/=52, ban #10 internal/risk/=21, ban #8 design/=39, ban #8 frontend/=188, ban #8 internal/=1066, ban #8 cmd/=359`） |

**记过一处**：首版 `plugins-screen.tsx` 注释里带了 `⛔`（U+26D4 在 d22scan/卫生测试的扫描带内、注释不豁免），`TestFrontendHasNoEmoji` 当场红（`plugins-screen.tsx:97`）——已改词复跑回绿。仪器/页面零 emoji 的断言自此坐实。

## 6. 载具操作记录（串行纪律）

- 跑前 `Get-Process wisp.test`＝0 枚（17:19 现量），go.exe 0 枚。
- 第一发 rc=2＝**旗标被拆**：pwsh 把裸 `-test.run` 拆成 `-test`＋`.run` 再交 bash（探针 `.scratch/wisp/probes/orch/20261119/407/argprobe.sh` 坐实；带引号 `'-test.run'` 即完整）。属调用方拼写错误，非载具缺陷；第二发带引号 rc=0。裸跑禁令未犯（两发均经载具）。

## 7. 未定项/移交

- 「来源/锁定说明」列的内容＝note（读面给的逐字说明行），缺席答「未知」——若编排者希望该列在无 note 时改答别的词，属契约文字，待批。
- 仪器产出的两张 markup 为**证据快照**（非 golden）：无任何逐字节比对闸门引用它们，重跑仪器即覆盖更新。
