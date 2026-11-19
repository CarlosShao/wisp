# 407 — K5 插件面板：后端有读面、界面未接、无规格（182 普查 K5，无票认领）

**Status:** **open**（2026-11-05 编排者按机主 2026-11-05「立新票」决定建；来路＝182 普查 K5）。
**Type:** 功能缺口（插件功能缺一块界面 + 无规格）＋无票认领。
**Packages:** 前端 frontend/**（插件面板 UI，**本编队全权负责**——机主 2026-11-16 铁律 4；旧"归机主外部 agent"作废）；可能牵 internal/plugin 只读 bridge。
**Blocked by:** 无（前端 face 现归本编队，不再等外部）；Go 侧 bridge 若有缺口可另立子格。

## 现场
182 普查（docs/evidence/s1/182-census.md）K5：插件面板后台能读（internal/plugin 库有读面），但前端从未接、也无界面规格；普查时无任何 open 票负责。

## 建议
前端 UI 现由本编队做；先把「插件面板该显示什么」的规格写清（可用 C17/面板既有形态），再接线。Go 侧若缺 bridge 读口，另开一条本编队可做的子票。

## 规格（2026-11-19 编排者补：机主 2026-11-16 铁律 4 后前端归本编队，本票可做）

**显示什么**（数据源只许用现成读面，⛔ 自造）：
- 来源＝`internal/panel/composer.go:223` 的 `Plugins *PluginInventorySection`（215a 载体，`plugins_inventory_215.go` 定形状）；写面＝`config.set` 的 `plugin_enabled`（`config_write_215.go`，须带 `configPlugin=<id>`，`config_handlers.go:575-579` 缺 id 即拒）。
- 每行：`<id>`、启用/停用、来源/锁定说明、以及**读面给了才显示**的其余字段（没给的字段**显示"未知"**，⛔ 编造）。
- 空清单：显式说"**还没有插件**"（零表格），⛔ 画 0 行假表。

## AC

- [x] **AC#1 规格落地**：插件页渲染上面那张表；每枚字段**逐字来自 `PluginInventorySection`**（⛔ 页面自造字段名/默认值）；空清单显示"还没有插件"，零表格。
- [x] **AC#2 开关走现成门**：启用/停用走 `config.set`（`plugin_enabled`＋`configPlugin`），⛔ 自造第二条写路；缺 `configPlugin` 的拒绝句原话上屏（`config_handlers.go:575-579`）。
- [x] **AC#3 缺失即未知**：读面未给某字段 ⇒ 该格显示"未知"，**两发变异**：①把缺字段当默认 true/false ⇒ 红 ②空清单画表 ⇒ 红。
- [x] **AC#4 仪器**：`render:pluginpage407` rc=0，出**两张**静态 markup（空态＋有插件态），零 emoji（`TestFrontendHasNoEmoji` 绿，仪器文件也在射程内）。
- [x] **AC#5 牙**：①页面引用 `PluginInventorySection` 里不存在的字段 ⇒ 红（源级/形状尺）②开关未走 `config.set`（自造写路）⇒ 红 ③写面回执把"需要重启"吞掉 ⇒ 红。
## 编排者收件（2026-11-19 18:2x）—— `1d08513ec`（实现）＋`aed0809df`（审计后补）—— 票结

**尺1** 两提交在。**尺3 我自己复跑**：仪器 `render-pluginpage407` **rc=0**（末行 OK…）,`npm run typecheck` **rc=0**、`go test ./internal/panel -run '^TestFrontendHasNoEmoji$'` ok、`d22scan` clean。
**对抗验收（`aud407`，只读，裁决者≠实现者）＝ 7 通过 / 1 可疑 / 0 不通过**，逐条读数：
- AC#1 字段逐字来自读面：页面 9 枚 row 键＋4 枚 section 键 **== Go json tag 全等**，反向 0；空态 markup `table/tr/td/th/ul/li` **全 0**（有插件态 table=1/tr=6/td=40 对照）。
- AC#2 写面：`git grep plugin_enabled` 真写入点 **1 处**（`config-actions.ts:151`→`config.set`）；Go 唯一字节改动 `internal/config/settings_plugins_215.go:138`；`SourceNative` 无插件写调用者、CLI 名册只开 `ball.style` ⇒ **无第二条路**。
- AC#3 未知规则：`enabled` 缺省＝"未知"且不画开关（`:78-82`/`:116-118`/`:170`）；**审计指出行级 `rows` 无守卫**（段级 `:261` 已合规）⇒ `aed0809df` 补 `rowsGiven`（`Array.isArray`）：缺/null/非数组 ⇒「未知（读面没有给清单）」零表格不抛；`[]` ⇒「还没有插件」。
- AC#4 仪器：两张 static markup 在 `docs/evidence/s1/`（空态 1606B／有插件态 8445B），零 emoji。
- AC#5 牙：五发变异**审计亲跑 5/5 全红**且红句逐字一致；补件后新增三发（M1 守卫拆⇒崩红／M2 手抄折断⇒红／M3 下游折"保存失败"⇒**3 处 FAIL** 含 `[refusal text head: 保存失败]`）。fixture 已**不再手抄**：仪器启动现读 `internal/panel/config_handlers.go`＋`bridge.go` 组装。