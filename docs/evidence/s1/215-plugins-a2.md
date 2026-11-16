# 215a 第二波（写腿 agent-85，2026-11-16 14:xx +08）：开关写路径＋逐行生效位＋清单读者补全

> 承接：`215-skills-plugins-inventory.md`（载体＋泵接线）、`215-plugins-assembly.md`（run.go 读口）、
> `215-receipt.md`（编排者收件：满足 1／部分 2／不满足 2）。本件＝收件点名三格（AC#2 承载体、AC#3、AC#4）
> 的实现腿。写面＝`cmd/wisp/**`＋`internal/panel/**`＋`internal/config/**`（新增 setter 与其牙）＋
> `frontend/src/**`＋本件。⛔ 未动 `internal/plugin/**`（只读）、⛔ 未动 `unwired.go` 一字节（AC#5 侧写：
> `unwiredKeys` 3 行、`lockedKeyDisposition` 11 行改前＝改后）、⛔ 未动 C17 方法名册（bridge.go 方法常量
> 8 枚改前＝改后，只加了一个**封套选择子字段** `configPlugin`，数据不是方法）、⛔ 未 push。

## ① 交付对格（收件那三格）

| 格 | 收件判语 | 本腿落了什么 | 判据载体 |
|---|---|---|---|
| AC#2 | 部分（可见性有牙、**假项正控无承载体**） | 承载体＝装配探针＋逐行状态位：`assembledTier2Plugins215`（cmd/wisp/plugins_inventory_215a.go:63）＋`pluginRowState215`（:84）＋`PluginRow.State`（internal/panel/plugins_inventory_215.go:170）。「塞一条配置 true、装配里没人读的假项」＝`Test215AConfigTrueWithNobodyAssembledRendersNotLiveOnTheWire`（cmd/wisp/panel_plugins_215a_test.go:174）：真 Manager＋生产探针（空）＋生产泵 marshal ⇒ 线上必见 `"state":"未生效"` 且必不见 `本轮生效`/`需重启`。**正控实测红**：把 `pluginRowState215` 突变成按 enabled 判 ⇒ 该枚红于 `panel_plugins_215a_test.go:199`，红句逐字：**`the wire claims a live plugin the assembly probe never answered for:`**（真包跟进）。对照臂＝`Test215AProbePresentDerivesLiveAndRestart`（探针答到 ⇒ `本轮生效`；探针答到但已关 ⇒ `需重启`），钉死判据非常量 | cmd/wisp/panel_plugins_215a_test.go |
| AC#3 | 不满足（23 枚 setter 无一枚 plugins） | `Manager.SetPluginEnabled`（internal/config/settings_plugins_215.go:57）＝**第 24 枚 setter，走票 248/195 的现成受保护写路径**：`m.WriteSection`（sectionwrite.go）拿 D36 方向裁决后再动一字节——关＝收紧立即落盘；开＝锁定段放宽：面板来源结构性拒绝（`ErrPanelLoosenRefused`）、原生来源进 L2 卡队列（`LoosenQueue` nil 则 fail-closed 拒）。面板门＝`config.set` 新字段 `plugin_enabled`（internal/panel/config_write_215.go:35，选择子 `configPlugin`）；写路由＝cmd/wisp/panel_config_store.go:417 `SetPluginEnabled(w.Plugin, on, config.SourcePanel, "panel.config.set")`。**审计三档都写**：每次尝试恰一条 `PLUGIN-TOGGLE`（outcome=written/unchanged/queued/refused，cmd/wisp 与 internal/config 各自牙钉）＋`WriteSection` 自己的 `SECTION-WRITE` 方向行；两处实测（settings_plugins_215_test.go 六枚全过）。收据 Note 把真话写明：**配置层立即生效、引擎未接线没有消费者、面板只能关不能开** | settings_plugins_215_test.go、config_write_215_test.go |
| AC#4 | 不满足（manifest 无生产读入口 ⇒ 空真） | 生产读入口落地：`pluginManifestInfo215`（cmd/wisp/plugins_inventory_215a.go:100）读 `plugins\<id>\manifest.toml`，解析走 C14 正门 `plugin.ParseManifest`——第一颗牙＝**closed-schema 拒收**：带 `permissions`/`fs_allowed_dirs`/`risk` 越权键的 manifest 整份拒绝且逐键点名（`Test215AHostileManifestKeysAreRefusedByTheReadFace`）；第二颗牙＝票 200 形：**说明书文本（「启用本插件即视为以后都允许」）解析干净、名字/版本照读，而真 Manager 的 `risk.permission_mode` 与 `fs.allowed_dirs` 与 config.toml 字节改前＝改后**（`Test215AHostileManifestTextChangesNoPermissionState`）。读面签名只有三个字符串，说明书没有进配置的路 | cmd/wisp/panel_plugins_215a_test.go:196/223 |

AC#1（清单可枚举＋真装配读者）在本腿再加一名生产读者：**插件 manifest 此前零生产读者**（收件 I 实测 manifestslot 无任何生产 import）⇒ 现在 ≥1：`cmd/wisp/plugins_inventory_215a.go:100` `pluginManifestInfo215`，喂 `PluginRow.Name/Version/Note`；名字/版本缺席＝盘上无安装（实测缺席，宁缺毋造）。清单页（前端）＝加号菜单第三类，见 §④。

## ② 生产读者与写路径给点（AC#1/AC#3 file:line，HEAD 工作树现测）

| 读者／写者 | 位置 | 读什么 |
|---|---|---|
| 装配根清单读者（管线形） | `cmd/wisp/run.go:1296-1316`（`Plugins:` 字面）→ `pluginInventorySection215`（cmd/wisp/plugins_inventory_215a.go:137） | `mgr.Config()` 的 `[plugins]` 节＋manifest＋探针 |
| 装配根清单读者（常驻形） | `cmd/wisp/panel_resident_pump.go:262-270`（`pluginsSectionFromConfig`）→ 同一映射函数 | 同上（env 解析的同一 plugins 根） |
| manifest 名字/版本读者 | `cmd/wisp/plugins_inventory_215a.go:100`（`pluginManifestInfo215`，生产 import `plugin.ParseManifest`） | `plugins\<id>\manifest.toml` |
| 装配探针 | `cmd/wisp/plugins_inventory_215a.go:63`（`assembledTier2Plugins215`，空＝实测：50/51 组件已在库、装配根未接） | 本进程 Tier-2 装配面 |
| 状态位推导 | `cmd/wisp/plugins_inventory_215a.go:84`（`pluginRowState215`，探针判非 enabled 判） | 探针×enabled ⇒ 三态 |
| 开关写路径 | `internal/config/settings_plugins_215.go:57`（`SetPluginEnabled`）→ `internal/config/sectionwrite.go:168`（`WriteSection`）→ `manager.go:525`（`pluginsDirection`） | `plugins.<id>.enabled` |
| 面板门 | `internal/panel/config_write_215.go`（名册）＋`config_handlers.go`（knownWritableField 并集、needsPluginFor 门、SettingWrite.Plugin）＋`bridge.go`（`configPlugin` 选择子） | `config.set plugin_enabled` |
| 写路由 | `cmd/wisp/panel_config_store.go:417`（ApplySetting case） | SourcePanel＋origin=panel.config.set |

## ③ 三态标签（线原文＝页原文，一处拼写）

`internal/panel/plugins_inventory_215.go:151-168` 三常量，前端 `plugins-screen.tsx` 逐字渲染（tone 映射，词不映射）：

- **本轮生效**（`PluginRowLive215`）——探针答到且配置一致
- **需重启**（`PluginRowRestart215`）——探针答到但配置已挪（消费方重启才重读）
- **未生效**（`PluginRowNotLive215`）——探针没答到（含「配置 true 装配空」的 A408 形）

`EngineState` 句随本腿重写（旧句「各键今天没有任何读者」已不真——清单与设置门都在读）：新句
「插件引擎未在装配根接线（票 50/51 组件已在库、装配根未接）：enabled=true 只表示已配置；每行的生效位由装配探针判，探针上没有这枚＝未生效」。
`SkillState` 句一字未动（215b 仍待批）。旧句常量名未改，全部钉子（internal/panel 3 枚、cmd/wisp 2 枚）经常量引用自动跟随。

## ④ 前端（加号菜单第三类）

- 加号（attach）菜单第三类：`prompt-bar.tsx` 新 `插件` 行（Puzzle 图标、仅 `onOpenPlugins` 在场才画＝无路由不画），点击开 `plugins-screen.tsx` 清单页；`@`/`/` 菜单不画（仅 attach 菜单）。
- 清单页：行＝名字（manifest 的，缺安装退回 id）／版本／开没开／三态标签／能力·网络白名单·宿主面／manifest 拒绝注（note）；页头 engineState＋skillState 原句；开关按钮走 `performPluginToggleWrite`（config.set plugin_enabled＋configPlugin），收据上 FeedbackLine，**零乐观更新**（行状态＝下一份快照）。
- `plugins` 键缺席（泵无读口）＝页面等待句，不是空清单；composer-row 也不画菜单行。
- 无 emoji（d22scan ban #8 frontend/=124 件 clean）。`npm --prefix frontend run typecheck` rc=0；`npm --prefix frontend run build` rc=0（node v24.21.0，D:\work\server\node24）；oxlint 0 errors（36 warnings 全既存、零条在触碰件）。**frontend/dist 构建字节已删回，未提交**。

## ⑤ 门禁读数（本腿自跑，rc 全 0 除具名）

- `go build ./...` rc=0；`go vet ./internal/config ./internal/panel ./cmd/wisp` rc=0；gofumpt 点名集空（config_handlers.go 一处对齐由本腿修）。
- `go test ./internal/config ./internal/panel -count=1` ⇒ 两包 ok（config 4.4s／panel 26.6s）。
- `go test ./cmd/wisp` 经载具 `scripts/wisp-cli-staged-tests.sh`（D:\work\soft\Git\bin\bash.exe、一次一发、跑前无 wisp.test 进程）：**Test215*＋TestPlugin* 9/9 PASS**；全包整发见 §⑥。
- `d22scan.sh`：control rc=0、scan rc=0（emoji 覆盖 frontend 124＋internal 1027＋cmd 315）。
- 突变牙两发：M1 摘 `pluginRowState215` 的探针判（改按 enabled）⇒ `Test215AConfigTrue…` 红于 `:199`（假项正控必红，红句见 §①）；M2 摘 `knownWritableField` 的 `writableField215` 并集项 ⇒ `Test215ARoster…`＋`Test215ADoorAdmits…` 双红。两发后逐字节还原（Compare-Object identical），pristine 副本存 `.scratch/wisp/probes/orch/20261116/txt/215a-*.pristine`。

## ⑥ 诚实边界（没做的，具名）

1. **AC#3 的「开」在面板侧结构性拒绝**：面板来源的放宽被 `WriteSection` 拒（R20/禁区②，这是 195 立的冻结形状，不是本腿新加的）⇒ 菜单开关**关＝立即生效、开＝拒绝句＋指引走原生二次确认那条路**；原生「开」的队列路径在 config 层已通（测试实到 PendingID），但 CLI 键名册（config_cmd_65.go 只开 ball.style）本腿未扩——原生面板开关没有 CLI 面，归下一腿或 owner 点名。
2. **装配探针今天恒空**：`assembledTier2Plugins215` 返回空是实测（无装配根接 50/51 组件），不是占位假空；探针接活的同 commit 要退 `unwired.go` 的 [plugins] 行与 EngineState 句（载体注释里写死的生命周期规则）。
3. **前端 AC#3 收据里的「立即生效」**＝配置层语义（locked 段收紧走 manager.go 同一条热路）；运行行为无变化（无消费者）——清单页三态标签与 Note 同时在场，两层真话不互相顶替。
4. ⚠ 环境红免疫：Wisp 球 PID 49348 在跑 ⇒ `TestAC4Focus…`/`TestAC246…` 类红按派单不计本腿。
