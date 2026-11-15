# 票 215 收件裁决 —— 加号菜单第三类（技能／插件可枚举）

- **裁决者**：收件腿 `agent-54`（**非实现者**：只读＋裁决；写面只有本件与 `.scratch/wisp/probes/215/**`）
- **锚点 commit**：`9ed201f5fc6dad9b976ade01434de126f318b0cd`（2026-11-15 13:41:57 +08，`git rev-parse HEAD` 现跑）
- **盘上读数时刻**：2026-11-15 13:42–13:58 +08（所有读数均现跑，⛔ 转述）
- **判语汇总**：**满足 1 ／ 部分 2 ／ 不满足 2（总 5）** ⇒ AC#5 满足；AC#1、AC#2 部分；AC#3、AC#4 不满足
- **「摘接线后哪枚红」一句话**：**`go test -overlay` 摘掉 `cmd/wisp/run.go:1307` 那格 ⇒ 0 枚红**（3/3 PASS）；那格只有**盘上字面钉** `Test215AssemblyRootCarriesThePluginsReader`，我用「沙箱 cwd＋mutant run.go」把同一枚红复现出来（§④.5）；**另一生产赋值点 `cmd/wisp/panel_resident_pump.go:144` 摘掉 ⇒ `Test65ResidentPumpCarriesConfigTruth` 红**（实测，§④.6）
- **票面与盘上不一致处**：`I1`–`I7`（§⑤），其中要命的是 **`I1`：票面全篇只记 `run.go` 一处赋值点，盘上是两处**；**`I4`：`grep -c Plugin internal/config/schema.go` 的「8」在立票提交 `c1fa2e1d4` 上就是 18 行、HEAD 19 行，「8」从未复现**；`I2/I3/I5/I6/I7` 是行号与时刻读数腐烂（逐条点名，不据以翻案）
- 本件按 `AGENTS.md` 与派单口径：**票 215 的 AC 框、Status、Progress log 一字未动**；`cmd/wisp/**` 只读；**Git：本腿不 commit、不 push**（收件件留 untracked 交编排者，理由见 §⑥.1）

---

## ① 逐格五列判语

| AC 原文逐字（`.scratch/wisp/issues/215-composer-plus-menu-skills-plugins-inventory.md:32-36`） | 判语 | 依据 `file:line` ＋原文 | 可核命令 | 缺什么 |
|---|---|---|---|---|
| **AC#1** 清单可枚举，**每一条都有真装配读者**（现跑 grep 点数：改前"某键零读者"⇒ 改后 ≥1）。 | **部分** | 插件半**成立**：载体 `internal/panel/plugins_inventory_215.go:106-107`（`type PluginInventorySection struct`）＋读口 `internal/panel/pump.go:296`（`Plugins func() PluginInventorySection`）＋两个**生产**赋值点 `cmd/wisp/run.go:1307`、`cmd/wisp/panel_resident_pump.go:144`；读者 `540-563` 拷贝排序、nil→`[]`。技能半**不存在**：`grep -c Skill internal/config/schema.go`＝0（现跑），线上只有登记句 `internal/panel/plugins_inventory_215.go:79`（`SkillsNotInConfigModel215`）。**「改前零读者⇒改后≥1」实测成立**：`5341ead7a^` 上 `internal/config` **包外**消费者＝0（包内只有自己的 marshal 回写 `internal/config/parse.go:172` 与方向审计 `internal/config/manager.go:525`），HEAD＝两处装配根 | `git grep -n 'c.Plugins\|cfg.Plugins' -- 'cmd/*.go' 'internal/**/*.go' ':!*_test.go'`（输出见 §③ 第 5 条） | 技能那一份清单（215b，补配置键＝动 D36，待人工批准）；`run.go` 那格**没有编译视图的牙**（只有盘上字面钉，§④）；界面渲染（`A410` 前置） |
| **AC#2** **"配置开着但没生效"必须能被看出来**——正控：塞一条"配置 true、装配里没人读"的假项 ⇒ 判据必红（这一形就是 `A408` 点名的"文案在、控件不在"的镜像）。 | **部分** | **可见性成立且可验**：常驻泵映射 `cmd/wisp/panel_resident_pump.go:252`（`EngineState: panel.PluginEngineNotWired215`）＋`253`（`SkillState`）；句子本体 `internal/panel/plugins_inventory_215.go:74`（`"插件引擎未装配（票 50/51）：[plugins] 各键今天没有任何读者…"`）。**我实测的牙**：把 `252` 那行摘掉 ⇒ `Test65ResidentPumpCarriesConfigTruth` 红于 `cmd/wisp/panel_resident_pump_windows_test.go:141`（`plugins.engineState = "", want the carrier's honest not-wired sentence constant`）。**但 AC 字面的正控（"塞假项 ⇒ 判据红"）没有承载体**：今天**没有任何一行可能是「真生效」**（引擎 50/51 未建），假项与真项同形，判据无法靠"假项"区分；已实现的是更强的整节声明（每行都是"配了但没跑"）＋摘句即红 | `go test ./cmd/wisp -run 'Test65ResidentPumpCarriesConfigTruth' -count=1 -overlay=<…>`（§④.7，实测红）；`go test ./internal/panel -run 'TestPlugin\|TestThePluginsKeyTravelsOnlyFromItsReader215' -count=1` ⇒ 3/3 PASS | 一处"每一行这一轮到底生效没有"的**逐行**位（今天只有整节句）；`run.go:1311` 那侧的产句**无牙**（改坏没人红）；AC 字面正控要等引擎（50/51）落地才有"真生效"可对照 |
| **AC#3** 开关动作**跨档位一致**：默认最严档里关一个技能要能立刻生效、不需要重启；做不到就**在清单上写明"这一枚要重启"**，不许静默假装生效。 | **不满足** | **开关本身不存在**：`internal/config` 的 `Manager` 全部写口 23 枚里**没有 plugins 写口**（`git grep -n 'func (m \*Manager) Set' -- 'internal/config/*.go'` ⇒ 无 `SetPlugin*`，输出见 §③ 第 7 条）；载体自己的注释把这一格登记为未建：`internal/panel/plugins_inventory_215.go:41-48`（`"THE FLIP HALF IS REGISTERED, NOT BUILT."`）。也没有"这一枚要重启"的**逐枚**标注（只有整节"引擎未装配"句），也没有开关审计 | 同上 + `git grep -n -i 'SetPlugin\|PluginSet' -- 'internal/config/*.go' 'cmd/wisp/*.go'` ⇒ 0 命中 | 交付 2 **整个**：读/写路、三档审计、跨档位行为、逐枚"要重启／即时"标注；这一条不是"差一点"，是"没开始" |
| **AC#4** 外部内容带进来的指令**改不动权限**（复用票 200 那枚 hostile-file 判据的形状：投一枚写着"以后都允许"的插件说明 ⇒ 档位与 `allowed_dirs` 一字不变）。 | **不满足** | **生产侧零读入口**：唯一会读插件清单文本的 `internal/plugin/manifestslot`（`manifestslot.go:84`、`install.go:207`、`lifecycle.go:305/429`）在 `internal/plugin/**` 之外**没有任何生产 import**（`git grep -n 'internal/plugin/manifestslot' -- '*.go' ':!internal/plugin/manifestslot/**' ':!*_test.go'` ⇒ 只有 `internal/plugin/command.go:16` 一条注释）；`c.Plugins` 的两个读者只读配置字段，不读任何外部文本（§③ 第 5 条）。⇒ 今天"改不动权限"是**空真**：不是做到了，是没有可被投毒的入口；也**没有一枚判据**（hostile-plugin 形状的用例不存在） | `git grep -n 'internal/plugin/manifestslot' -- '*.go' ':!internal/plugin/manifestslot/**' ':!*_test.go'`（输出见 §③ 第 8 条） | 插件文本的任何生产读入口、hostile-plugin 判据（"投一枚写着以后都允许的说明 ⇒ 档位／`allowed_dirs` 一字不变"）；本票这一格**未做** |
| **AC#5** `unwired.go` 那面旗**一枚不少**（现跑点数：改前＝改后）。 | **满足** | **本票射程内零字节**：`git diff --numstat 5341ead7a^..5341ead7a -- internal/config`＝空、`git diff --numstat 8d49dc562e^..8d49dc562e -- internal/config`＝空（§③ 第 9 条）；`unwiredKeys` 行数 `5341ead7a^`＝5、`5341ead7a`＝5、`8d49dc562e`＝5。**钉子会红**：overlay 摘掉 `internal/config/unwired.go:132` 的 `lockedKeyDisposition["net.allowlist"]` 行 ⇒ `TestEveryLockedSectionKeyIsAccountedFor` 红于 `internal/config/unwired_test.go:372`（实测，§④.8）；改前守卫用例组现跑 `ok internal/config 0.149s` | `go test ./internal/config -run 'TestEveryLockedSectionKeyIsAccountedFor\|TestUnwired' -count=1` ⇒ `ok 0.149s` | 无（本票射程内成立）。⚠ 但 **HEAD 现读＝3 行 ≠ 票面的 5 行**：两行是被 `7bcc2ae08`（票 163 AC#5，落 `risk.shell_enabled`／`risk.shell_allowlist` 的消费者即删行）拿走的，**不是 215 干的**，见 `I6` |

---

## ② 判语怎么来的（A1875 要求的「生产调用点＋会红测试名」成对表）

> `A1875`：判"满足"必须同时给 ①生产调用/派发点 `file:line` ②会因摘掉它而变红的测试名；缺一只能判"部分"。下表是每一格我**实际给出**的成对证据，以及我**没有**给出的那一半（这就是"部分"的来源）。

| AC | 生产调用/派发点（`file:line`＋原文） | 会红的测试名（我实测红／绿） | A1875 判定 |
|---|---|---|---|
| AC#1 | `cmd/wisp/panel_resident_pump.go:144` `src.Plugins = func() panel.PluginInventorySection { return pluginsSectionFromConfig(cfg) }`；`cmd/wisp/run.go:1307` `Plugins: func() panel.PluginInventorySection {` | **红**：`Test65ResidentPumpCarriesConfigTruth`（摘 `panel_resident_pump.go:144` ⇒ 红，§④.6）。**绿**：`Test215AssemblyRootCarriesThePluginsReader` 在 overlay 摘 `run.go:1307` 时**不红**（§④.4），只有盘上摘才红（§④.5，补证） | **部分**：`run.go` 那格缺"编译视图里的牙"；技能半无清单 |
| AC#2 | `cmd/wisp/panel_resident_pump.go:252` `EngineState: panel.PluginEngineNotWired215`；`cmd/wisp/run.go:1311` 同句 | **红**：`Test65ResidentPumpCarriesConfigTruth`（摘 `252` ⇒ 红于 `panel_resident_pump_windows_test.go:141`，§④.7）；**绿**：`TestPluginInventory215EnabledRowAndNotWiredSentenceRideOnePacket`（`internal/panel/plugins_inventory_215_test.go:171`） | **部分**：可见性有牙；AC 字面"塞假项⇒红"无承载体；`run.go:1311` 无牙 |
| AC#3 | **无**（`internal/config` 无 plugins 写口） | **无** | **不满足**：①与②都缺 |
| AC#4 | **无**（生产零读入口；`manifestslot` 未接线） | **无** | **不满足**：①与②都缺（空真） |
| AC#5 | 非"调用点"型判据：旗表 `internal/config/unwired.go:118-146`＋`60-96` | **红**：`TestEveryLockedSectionKeyIsAccountedFor`（摘一行处置 ⇒ 红于 `unwired_test.go:372`，§④.8）；**绿**：`TestUnwired*` 组现跑 `ok` | **满足**（本票射程内；见 `I6` 的 HEAD 3 行） |

---

## ③ 三件之① —— 现跑 grep：`PumpSources` 里 `Plugins` 那格的**赋值点**与**清单类型的生产实现处**

> 命令原文与输出逐条照录（`A1879`：每处引用写全仓相对路径；判据贴命令＋输出）。完整日志 `.scratch/wisp/probes/215/receipt/20-citations.txt`（114 行）。

**1）赋值点（`PumpSources` 字面量里 `Plugins:` 那一格）——两处，不是一处**

```
$ git grep -n 'Plugins: func() panel.PluginInventorySection' -- 'cmd/wisp/*.go'
cmd/wisp/panel_plugins_215_test.go:87:	if !strings.Contains(string(runGo), "Plugins: func() panel.PluginInventorySection") {
cmd/wisp/panel_plugins_215_test.go:141:		Plugins: func() panel.PluginInventorySection { return pluginsSection215(mgr) },
cmd/wisp/run.go:1307:		Plugins: func() panel.PluginInventorySection {
```

- **生产赋值点 A**：`cmd/wisp/run.go:1307`（`wisp run` 流水线泵；块 `1307-1324`，读 `mgr.Config()` 的 `[plugins]` 节，见 `1308`／`1310`／`1314`）
- **生产赋值点 B**：`cmd/wisp/panel_resident_pump.go:144`（**常驻**泵，即双击 `wisp.exe` 那条发货形状；映射函数 `panel_resident_pump.go:248-265`）
- 测试件里的两处（`panel_plugins_215_test.go:87/141`）是牙与镜像，**不是**生产赋值点

**2）清单类型的生产实现处（类型、读口、装配块、快照字段）**

```
$ git grep -n 'PluginInventorySection' -- 'cmd/wisp/*.go' 'internal/panel/*.go'
cmd/wisp/panel_resident_pump.go:144:		src.Plugins = func() panel.PluginInventorySection { return pluginsSectionFromConfig(cfg) }
cmd/wisp/panel_resident_pump.go:248:func pluginsSectionFromConfig(m *config.Manager) panel.PluginInventorySection {
cmd/wisp/panel_resident_pump.go:250:	sec := panel.PluginInventorySection{
cmd/wisp/run.go:1307:		Plugins: func() panel.PluginInventorySection {
cmd/wisp/run.go:1309:			sec := panel.PluginInventorySection{
internal/panel/composer.go:223:	Plugins *PluginInventorySection `json:"plugins,omitempty"`
internal/panel/plugins_inventory_215.go:106:// PluginInventorySection is the snapshot's plugins key (ticket 215a).
internal/panel/plugins_inventory_215.go:107:type PluginInventorySection struct {
internal/panel/pump.go:296:	Plugins func() PluginInventorySection
```

```
$ git grep -n 'Plugins' -- 'internal/panel/pump.go' 'internal/panel/composer.go'
internal/panel/composer.go:210:	// Plugins is ticket 215a's carrier: the Tier-2 plugin inventory the live
internal/panel/composer.go:223:	Plugins *PluginInventorySection `json:"plugins,omitempty"`
internal/panel/pump.go:283:	// Plugins reads the Tier-2 plugin inventory the live config holds (ticket
internal/panel/pump.go:296:	Plugins func() PluginInventorySection
internal/panel/pump.go:540:	if p.src.Plugins != nil {
internal/panel/pump.go:552:		sec := p.src.Plugins()
internal/panel/pump.go:562:		snap.Plugins = &sec
```

- 行载体：`internal/panel/plugins_inventory_215.go:83`（`type PluginRow struct`）、`:107`（`PluginInventorySection`）、`:70-80`（两句 canonical 常量）
- 装配块（reader-gated＋拷贝排序＋nil→`[]`）：`internal/panel/pump.go:540-562`（`553-557` 拷贝后 `sort.Slice` 按 ID，`559-561` nil→`[]`，**泵不写回读者 slice**）
- 快照字段：`internal/panel/composer.go:223`
- 真源：`internal/config/schema.go:637-640`（`type PluginsSection struct`／`Entries map[string]PluginEntry`）、`internal/config/parse.go:93`（`buildPlugins`）

**3）两句诚实句的生产位置（三处）**

```
$ git grep -n 'PluginEngineNotWired215\|SkillsNotInConfigModel215' -- '*.go' ':!*_test.go'
cmd/wisp/panel_resident_pump.go:252:		EngineState:  panel.PluginEngineNotWired215,
cmd/wisp/panel_resident_pump.go:253:		SkillState:   panel.SkillsNotInConfigModel215,
cmd/wisp/run.go:1311:				EngineState:  panel.PluginEngineNotWired215,
cmd/wisp/run.go:1312:				SkillState:   panel.SkillsNotInConfigModel215,
internal/panel/plugins_inventory_215.go:74:	PluginEngineNotWired215 = "插件引擎未装配（票 50/51）：[plugins] 各键今天没有任何读者，enabled=true 只表示已配置、不表示在跑"
internal/panel/plugins_inventory_215.go:79:	SkillsNotInConfigModel215 = "技能维今天没有配置键，无清单可枚举；补技能清单要先动 D36 配置模型（契约面，待批），本清单不列空表"
```

**4）读者＝`mgr.Config()`（锁内深拷）** —— 票面 :54 的说法在盘上成立

```
$ git grep -n 'func (m \*Manager) Config()' -- 'internal/config/*.go'
internal/config/manager.go:127:func (m *Manager) Config() *Config {
```
`internal/config/manager.go:125-131`：`// Config returns a deep copy of the current effective config …` ＋ `m.mu.Lock()` → `return deepCopyConfig(m.cur)`。

**5）`[plugins]` 配置键的生产读者（改前／改后点数）**

```
$ git grep -n 'c.Plugins\|cfg.Plugins' -- 'cmd/*.go' 'internal/**/*.go' ':!*_test.go'
cmd/wisp/panel_resident_pump.go:144:		src.Plugins = func() panel.PluginInventorySection { return pluginsSectionFromConfig(cfg) }
cmd/wisp/panel_resident_pump.go:251:		Tier2Enabled: c.Plugins.Tier2Enabled,
cmd/wisp/panel_resident_pump.go:255:	for id, e := range c.Plugins.Entries {
cmd/wisp/run.go:1310:				Tier2Enabled: c.Plugins.Tier2Enabled,
cmd/wisp/run.go:1314:			for id, e := range c.Plugins.Entries {
internal/config/parse.go:172:	tail, err := marshalPlugins(c.Plugins)
internal/panel/pump.go:540:	if p.src.Plugins != nil {
internal/panel/pump.go:552:		sec := p.src.Plugins()
```

```
$ git grep -n 'c\.Plugins\|cfg\.Plugins\|\.Plugins\.Entries' 5341ead7a^ -- '*.go' ':!*_test.go'
5341ead7a^:internal/config/parse.go:172:	tail, err := marshalPlugins(c.Plugins)
```

⇒ **改前**：`[plugins]` 在 `internal/config` **包外零消费者**（包内只有自己的 marshal 回写 `parse.go:172` 与方向审计 `manager.go:525`／`sectionwrite.go:356`）；**改后**：两处装配根读者（`run.go:1307`、`panel_resident_pump.go:144`）。AC#1 的"≥1"与票面 :52 的"引擎未建⇒无生效位"**同时**成立——读到的是配置的**字**，不是引擎的**跑**。

**6）票面数字的复读（`I4`）**

```
HEAD schema.go lines matching 'Plugin' = 19
HEAD schema.go lines matching 'Skill'  = 0
c1fa2e1d4（立票提交 2026-09-28 23:38）schema.go lines matching 'Plugin' = 18
c1fa2e1d4^ schema.go lines matching 'Plugin' = 18
```
⇒ 票面 :9／:18 的 `grep -c Plugin internal/config/schema.go`＝**8** 在**立票提交上就是 18**、HEAD 19，**从未复现**；`Skill`＝0 复现。维度结论（插件有键／技能无键）不受影响，但"8"这个数字是坏的。

**7）写口（AC#3 的缺）**

```
$ git grep -n 'func (m \*Manager) Set' -- 'internal/config/*.go'
internal/config/allowdirs.go:109:func (m *Manager) SetAllowedDirs(dirs []string) error {
internal/config/permmode.go:70:func (m *Manager) SetPermissionMode(mode risk.Mode) error {
internal/config/settings.go:100:func (m *Manager) SetProviderBaseURL(provider, url string) ([]string, error) {
internal/config/settings.go:121:func (m *Manager) SetProviderAPIKeyRef(provider, ref string) ([]string, error) {
internal/config/settings.go:143:func (m *Manager) SetModelContextWindow(provider, model string, tokens int) ([]string, error) {
internal/config/settings.go:172:func (m *Manager) SetModelPriceIn(provider, model string, per1M float64) ([]string, error) {
internal/config/settings.go:177:func (m *Manager) SetModelPriceOut(provider, model string, per1M float64) ([]string, error) {
internal/config/settings.go:224:func (m *Manager) SetRoleChatModel(model string) ([]string, error) {
internal/config/settings_187.go:104:func (m *Manager) SetRoleChatThinking(level string) ([]string, error) {
internal/config/settings_313.go:56:func (m *Manager) SetVoiceCloudASRChain(chain []string) ([]string, error) {
internal/config/settings_313.go:63:func (m *Manager) SetVoiceCloudTTSChain(chain []string) ([]string, error) {
internal/config/settings_313.go:124:func (m *Manager) SetVoiceRealtimeProvider(provider string) ([]string, error) {
internal/config/settings_313.go:134:func (m *Manager) SetVoiceRealtimeModel(model string) ([]string, error) {
internal/config/settings_313.go:208:func (m *Manager) SetRoleChatPair(provider, model string) ([]string, error) {
internal/config/settings_320.go:36:func (m *Manager) SetProviderBilling(provider, billing string) ([]string, error) {
internal/config/settings_413.go:45:func (m *Manager) SetBallSize(size int) ([]string, error) {
internal/config/settings_65.go:46:func (m *Manager) SetBallStyle(style string) ([]string, error) {
internal/config/settings_appearance_wave.go:56:func (m *Manager) SetBallOpacityIdle(v float64) ([]string, error) {
internal/config/settings_appearance_wave.go:69:func (m *Manager) SetBallClickThrough(on bool) ([]string, error) {
internal/config/settings_appearance_wave.go:80:func (m *Manager) SetPanelWidth(px int) ([]string, error) {
internal/config/settings_appearance_wave.go:95:func (m *Manager) SetPanelFontSize(px int) ([]string, error) {
internal/config/settings_local_audio.go:38:func (m *Manager) SetVoiceAecEnabled(on bool) ([]string, error) {
internal/config/settings_local_audio.go:50:func (m *Manager) SetAudioInputDevice(name string) ([]string, error) {
```
23 枚 setter，**无一枚 plugins**；`SetPlugin`／`PluginSet` 全仓 0 命中。

**8）AC#4 的空真（读入口不存在）**

```
$ git grep -n 'internal/plugin/manifestslot' -- '*.go' ':!internal/plugin/manifestslot/**' ':!*_test.go'
internal/plugin/command.go:16://     in internal/plugin/manifestslot calls this on every call; the test leg
```
⇒ `internal/plugin/manifestslot`（含 `install.go:207`／`lifecycle.go:305/429` 的 `RegisterProvider`）**没有任何生产装配点**，即没有任何外部插件文本进入运行进程的路径。

**9）AC#5 的两把尺**

```
$ git diff --numstat 5341ead7a^..5341ead7a -- internal/config      → （空）
$ git diff --numstat 8d49dc562e^..8d49dc562e -- internal/config    → （空）
$ (PowerShell) foreach ($ref in '5341ead7a^','5341ead7a','8d49dc562e','HEAD') { "$ref " + (git show "${ref}:internal/config/unwired.go" | Select-String -Pattern '^\s+path:\s+"').Count }
5341ead7a^ unwiredKeys rows = 5
5341ead7a  unwiredKeys rows = 5
8d49dc562e unwiredKeys rows = 5
HEAD       unwiredKeys rows = 3     ← 差额来自 7bcc2ae08（票 163），非 215
```

---

## ④ 三件之② —— `go test -overlay` 摘接线那格（**树上零改动**）

> 派单口径：用 `-overlay` 摘掉接线格 ⇒ **指名哪枚红**；不红判"部分"。同时给 `git status --porcelain` 读数与原件 md5 前后同。
> 全部工件在 `.scratch/wisp/probes/215/receipt/`（untracked，⛔ 未进任何 commit）。

### ④.0 先撞上的环境事实（不藏）

现树 `cmd/wisp` **不是干净的**：写腿 `452` 的三枚在飞件＋一枚 untracked 测试件（13:44 首读见到 2 枚 `M`、13:45 起 3 枚 `M`——差异在 452 自己的文件上；我存证的起手 13:45 与收尾 13:56 两次读数一致，均非我的写面）：

```
$ git status --porcelain -- cmd/wisp
 M cmd/wisp/resident_audio_windows.go
 M cmd/wisp/resident_voice_bargein_windows.go
 M cmd/wisp/resident_voice_windows.go
?? cmd/wisp/resident_audio_device_lost_452_windows_test.go
```

⇒ raw 跑即**编译失败**（不是用例红，是包没编出来）：
```
$ go test ./cmd/wisp -run 'Test215' -count=1 -v
cmd\wisp\resident_audio_windows.go:404:21: ra.captureErr undefined (type *residentAudio has no field or method captureErr)
cmd\wisp\resident_voice_bargein_windows.go:332:6: invalid operation: cannot call aec.BindingBuilt (untyped bool constant false): untyped bool is not a function
cmd\wisp\resident_voice_bargein_windows.go:350:5: invalid operation: cannot call aec.BindingBuilt (untyped bool constant false): untyped bool is not a function
FAIL	github.com/CarlosShao/wisp/cmd/wisp [build failed]
```
⇒ 照 `62-v1` 先例用 overlay **只隔离他腿在飞件**（不改树）：三枚改动件映射到 `HEAD` blob（`git cat-file blob` 落地，`git hash-object` 与 `git rev-parse HEAD:<path>` **逐枚相等**：`050b24f6…`／`e8e9e9ef…`／`0a621827…`），未提交测试件用 overlay 空替换删除。基线与突变**共用这一层隔离**。

### ④.1 基线（隔离后，未摘接线）

```
$ go test ./cmd/wisp -run 'Test215' -count=1 -v -overlay=<overlay-isolate.json>
--- PASS: Test215AssemblyRootCarriesThePluginsReader (0.00s)
--- PASS: Test215WiredPluginsReaderSendsTheLiveValues (0.01s)
--- PASS: Test215MutationReaderRemovedSendsNoPluginsKeyAtAll (0.00s)
PASS  ok …/cmd/wisp  0.147s
```

### ④.2 **摘格**：`run.go` 的 `Plugins` 块（`1307-1324`，含注释）整块删除 ⇒ **0 枚红**

突变体 `.scratch/wisp/probes/215/receipt/mut-run-go-no-plugins.go`：从 HEAD blob 里按标记切掉整块，切后全文 `Plugins` 命中 **0**、`Out: rt.bookPanelSnapshotBlocking` 仍在（块边界没切歪）。

```
$ go test ./cmd/wisp -run 'Test215' -count=1 -v -overlay=<overlay-mut-no-plugins.json>
--- PASS: Test215AssemblyRootCarriesThePluginsReader (0.00s)   ← ← ← 不红
--- PASS: Test215WiredPluginsReaderSendsTheLiveValues (0.01s)
--- PASS: Test215MutationReaderRemovedSendsNoPluginsKeyAtAll (0.00s)
PASS  ok …/cmd/wisp  0.262s
```
**读数：摘掉接线格 ⇒ 0 枚红。** 按派单口径这一格**只能判"部分"**（原因在 ④.5）。

### ④.3 控制发：证明 overlay **确实**是编译器视图（否则 ④.2 的"不红"可能只是 overlay 没生效）

把同一块里的字段名改成 `PluginsZzz:`（其余逐字不动）⇒ 编译器在**突变体文件**上报错：
```
$ go test ./cmd/wisp -run 'Test215' -count=1 -overlay=<overlay-mut-bad-ident.json>
.\.scratch\wisp\probes\215\receipt\mut-run-go-bad-ident.go:1307:3: unknown field PluginsZzz in struct literal of type panel.PumpSources
FAIL	github.com/CarlosShao/wisp/cmd/wisp [build failed]
```
⇒ overlay 生效（错误指向突变体、行号 1307 与盘上一致）；**④.2 的 0 红是真读数，不是机制没接上**。

### ④.4 追加探针：把那一格的**内容**改坏（`Tier2Enabled: !c.Plugins.Tier2Enabled`）⇒ 仍 3/3 PASS

```
--- PASS: Test215AssemblyRootCarriesThePluginsReader
--- PASS: Test215WiredPluginsReaderSendsTheLiveValues
--- PASS: Test215MutationReaderRemovedSendsNoPluginsKeyAtAll
PASS  ok …/cmd/wisp  0.142s
```
⇒ **`run.go` 那格的映射内容没有任何测试看得见**：`Test215WiredPluginsReaderSendsTheLiveValues` 用的是**测试自己的镜像** `pluginsSection215`（`cmd/wisp/panel_plugins_215_test.go:100-117`），不是 `run.go:1308-1323`。这是"接线格只有字面钉"的第二个证据。

### ④.5 补证（**不是** overlay 口径）：那枚红**存在，但只对盘上文件红**

字面钉的判据本身就是读盘：`cmd/wisp/panel_plugins_215_test.go:82-89`（`os.ReadFile("run.go")` ＋ `strings.Contains(…, "Plugins: func() panel.PluginInventorySection")`）。⇒ overlay 换编译视图它看不见（票面 `docs/evidence/s1/215-plugins-assembly.md:34-35` 的"作废一发"记的正是这件事，我复现了它）。为了证明"牙不是没有、是盘上型"，我把生产的测试二进制编译出来、放进一个**只有 mutant `run.go`** 的沙箱目录再跑（等价于"盘上真删"的文件状态，工作树零触碰）：

```
$ (cd .scratch/wisp/probes/215/receipt/sandbox && ./wisp215.test.exe -test.run=Test215AssemblyRootCarriesThePluginsReader -test.v=true)
--- FAIL: Test215AssemblyRootCarriesThePluginsReader (0.00s)
    panel_plugins_215_test.go:88: run.go's PumpSources literal no longer carries the plugins reader: the snapshot's plugins key would silently vanish from every packet this process pushes (ticket 215a)
exit=1
$ (同一二进制，另两枚)
--- PASS: Test215WiredPluginsReaderSendsTheLiveValues (0.01s)
--- PASS: Test215MutationReaderRemovedSendsNoPluginsKeyAtAll (0.00s)
exit=0
```
**结论**：`Test215AssemblyRootCarriesThePluginsReader` 会红——**当且仅当 `run.go` 磁盘文件里没有那串字**；overlay 摘格口径下 0 红。两种说法都要写出来，不许拿补证冒充 overlay 口径。

### ④.6 另一生产赋值点（常驻泵）**有真牙**：摘 `panel_resident_pump.go:144` ⇒ 指名红

```
$ go test ./cmd/wisp -run 'Test65ResidentPumpCarriesConfigTruth' -count=1 -v -overlay=<overlay-mut-resident-pump.json>
    panel_resident_pump_windows_test.go:134: the resident packet carries no plugins key; top-level keys: [ballStyle panelFontSize pending composer generatedAt commands ballSize ballOpacityIdle ballClickThrough panelWidth results modelSwitch]
--- FAIL: Test65ResidentPumpCarriesConfigTruth (0.04s)
FAIL  …/cmd/wisp  0.254s
```
⇒ **摘掉常驻泵那一格 ⇒ `Test65ResidentPumpCarriesConfigTruth` 红**（红在 `cmd/wisp/panel_resident_pump_windows_test.go:134`，而且 reds 里点名 `plugins` 缺席）。

### ④.7 AC#2 的牙（摘掉"未装配"那句）

```
$ go test ./cmd/wisp -run 'Test65ResidentPumpCarriesConfigTruth' -count=1 -v -overlay=<overlay-mut-resident-no-engine-state.json>
    panel_resident_pump_windows_test.go:141: plugins.engineState = "", want the carrier's honest not-wired sentence constant
--- FAIL: Test65ResidentPumpCarriesConfigTruth (0.03s)
```

### ④.8 AC#5 的牙（摘掉一行处置 ⇒ 完整性钉红）

```
$ go test ./internal/config -run 'TestEveryLockedSectionKeyIsAccountedFor' -count=1 -v -overlay=<overlay-mut-unwired-row.json>
    unwired_test.go:372: locked key "net.allowlist" has no disposition in lockedKeyDisposition (ticket 83 AC#1: every locked key needs a verdict)
    unwired_test.go:389: unwiredKeys guards "net.allowlist" but no locked-section key points at it
--- FAIL: TestEveryLockedSectionKeyIsAccountedFor (0.00s)
```

### ④.9 树上零改动（派单要的两条读数）

```
run.go md5(before) = B7D7EBC0534528EFAF3BA1B3B287B327
run.go md5(after)  = B7D7EBC0534528EFAF3BA1B3B287B327      ← 同
$ git status --porcelain -- internal/panel internal/config
（空）
$ git status --porcelain -- cmd/wisp        （13:45 与 13:56 两次同读，均只有 452 的四行，见 §④.0）
```
⚠ **诚实边界**：`cmd/wisp` 的 `git status` **不是空读数**——那是写腿 452 的在飞件（起手就在，收尾仍在；我零触碰）。我能给的是：**我的写面零字节**（`internal/panel`、`internal/config` 空读数；`cmd/wisp` 无我新增/修改的文件）＋`run.go` md5 前后同。

### ④.10 探针总表

| # | 突变（全部 overlay，树上零改） | 命令靶 | 读数 | 判 |
|---|---|---|---|---|
| 0 | 无（隔离他腿在飞件） | `-run 'Test215'` | 3/3 PASS 0.147s | 基线绿 |
| 1 | 摘 `run.go:1307-1324` 整块 | `-run 'Test215'` | **3/3 PASS** 0.262s | **0 红 ⇒ 部分** |
| 2 | 字段改名 `PluginsZzz`（控制发） | 同 | build failed，指到突变体 `:1307` | overlay 确在编译视图 |
| 3 | 映射内容取反（`!c.Plugins.Tier2Enabled`） | 同 | 3/3 PASS 0.142s | 内容无牙 |
| 4 | 沙箱 cwd＋mutant run.go（补证，非 overlay） | 字面钉 | **FAIL** `panel_plugins_215_test.go:88` | 盘上型牙存在 |
| 5 | 摘 `panel_resident_pump.go:144` | `Test65ResidentPumpCarriesConfigTruth` | **FAIL** `:134` | 常驻赋值点有真牙 |
| 6 | 摘 `panel_resident_pump.go:252`（未装配句） | 同 | **FAIL** `:141` | AC#2 有牙 |
| 7 | 摘 `unwired.go:132` 处置行 | `TestEveryLockedSectionKeyIsAccountedFor` | **FAIL** `unwired_test.go:372` | AC#5 有牙 |

---

## ⑤ 三件之③ —— 票面记的"已落" vs 盘上现读（逐条对上或点名不一致）

> 尺：票面每一句"已落"都要在盘上找到对应字节，或者点名它今天读不到。`✓`＝对上；`≈`＝结论对上但**字面读数/行号已腐烂**；`≠`＝对不上。

| # | 票面原文（行号） | 盘上现读（`file:line`） | 判 | 说明 |
|---|---|---|---|---|
| P1 | `:52` "`panel.Snapshot` 新增 reader-gated 顶层键 `plugins`（载体 `internal/panel/plugins_inventory_215.go`＋`composer.go` 字段＋`pump.go` 读口与装配块）" | `internal/panel/plugins_inventory_215.go:107`；`internal/panel/composer.go:223`；`internal/panel/pump.go:296`；装配块 `internal/panel/pump.go:540-562` | ✓ | 三件套逐处命中 |
| P2 | `:52` "真源现跑核过：Plugin 8 命中/Skill 0" | `internal/config/schema.go` 现读 **Plugin 19 行／Skill 0 行**；**立票提交 `c1fa2e1d4` 上就是 18／0** | **≠（`I4`）** | 维度结论对（有 `[plugins]`、无 skills），**"8"从未复现** |
| P3 | `:52` "泵拷贝后按 ID 排序（map 顺序不上线、泵不写回读者 slice）、nil rows 渲染 `[]`" | `internal/panel/pump.go:553-557`（拷贝后 `sort.Slice`）、`559-561`（nil→`[]`） | ✓ | 与票面逐字相符 |
| P4 | `:52` "无读者基线＝无 `plugins` 键＋顶层恰 4 键" | `internal/panel/plugins_inventory_215_test.go:43-59`（无读者⇒无键、`len(generic) != 4` 即红），现跑 PASS | ✓ | 用例现跑绿（§④.10 基线） |
| P5 | `:52` "`internal/config/**` 零字节，unwired guard 表 5 行改前＝改后" | `git diff --numstat 5341ead7a^..5341ead7a -- internal/config` 空；`unwiredKeys` 行数 5→5 | ✓（射程内） | HEAD 现读 **3** 行，是 `7bcc2ae08`（票 163）拿走的 2 行（`I6`） |
| P6 | `:52` "生产装配根一行（`cmd/wisp/run.go` 加 `Plugins:` 读 config）未接——本波 `cmd/wisp/**` 被 219 腿独占，归下一腿" | `git show 5341ead7a:cmd/wisp/run.go` 里 `Plugins:` **0 命中** | ✓ | 当时的事实，下一腿 `8d49dc562e` 接了 |
| P7 | `:53` "commit `5341ead7a` 自跑复验，`internal/panel` 红 3→1（只剩 C21）" | 今天 `go test ./internal/panel -run 'TestC21DesignTokensFourWayAgree' -count=1` ⇒ `ok`（两次读数 0.068s／0.271s，第二次 `-v` 见 `13-c21-current.txt` 的 `--- PASS`） | ≈（`I5`） | 时刻读数；C21 今天已绿（他人后修 ball token），不据以翻案 |
| P8 | `:53` "界面腿给 `plugins` 键补了 TS 声明（`panel.ts:372-384` `plugins?: PluginsSection`）" | `frontend/src/lib/panel.ts:516` `plugins?: PluginsSection;`（段落在 `436-448`），落盘提交 `a3fdb50dd`（`177/0`，commit message 点名 `PluginsSection`） | ≈（`I2`） | 事实对（声明在、提交真），**行号漂了**（372-384→436-448/516） |
| P9 | `:53` "recon 程序红随界面腿落地即转绿" | `go test ./internal/panel -run 'TestComposerContractTypesMatchFrontend'` ⇒ PASS，`composer_test.go:98` 打 `Snapshot <-> PanelSnapshot: 16 JSON keys reconciled` | ✓ | 现跑绿 |
| P10 | `:54` "`cmd/wisp/run.go` PumpSources 字面量加 `Plugins: func() panel.PluginInventorySection` 读口（`run.go:1104`）" | `cmd/wisp/run.go:1307`；`git show 8d49dc562e:cmd/wisp/run.go` 里该行**恰在 1104** | ≈（`I3`） | 落盘时 1104 是真的，HEAD 漂到 1307（后续腿在 run.go 前面加了代码）；**引用今天必须写 1307** |
| P11 | `:54` "净 +37 行零删除" | `git show --numstat 8d49dc562e` ⇒ `37 0 cmd/wisp/run.go`；另 `209 0 cmd/wisp/panel_plugins_215_test.go`、`56 0 docs/evidence/s1/215-plugins-assembly.md` | ✓ | 逐数对上 |
| P12 | `:54` "读者＝`mgr.Config()`（Manager 锁内深拷）" | `internal/config/manager.go:125-131`（`// Config returns a deep copy …`＋`m.mu.Lock()`→`deepCopyConfig`） | ✓ | 并发面说法成立 |
| P13 | `:54` "`internal/config` **无** plugins 专用 reader（Options 族仅 BallStyle/ChatSwitch，现跑核过）" | **无** `PluginsOptions`（全仓 0）；但 `Manager.*Options()` 今天**7 枚**：`internal/config/settings_65.go:86`、`settings_187.go:254`、`settings_413.go:75`、`settings_appearance_wave.go:119/141/159/179` | ✓前半／≈后半（`I7`） | 前半（无 plugins reader）今天仍真；后半"仅两枚"是 11-08 那一刻的读数，被 413／appearance 波次加宽 |
| P14 | `:54` "新测试件 `cmd/wisp/panel_plugins_215_test.go` 照 187 三牙" | 209 行、三枚用例：`:82`／`:126`／`:201` | ✓ | 三牙俱在 |
| P15 | `:54` "M1 真突变＝`sed` 摘 `run.go` 接线块 → 字面钉 FAIL → `cmp` 字节还原 → 复绿；**先试 overlay 换 pristine 副本没红＝作废如实记：字面钉读盘上文件、不读 overlay 编译视图**" | 我独立复现：overlay 摘块 ⇒ **0 红**（§④.2）；沙箱盘上型 ⇒ 红于 `:88`（§④.5） | ✓ | **票面自己已经点名了这个 overlay 盲区**；本件把它做成可复跑的探针，并把"盘上型红"也补成了读数 |
| P16 | `:54` "`go test ./cmd/wisp -run 'Test215\|TestPlugins' -count=1` 3/3 PASS 0.087s" | 隔离后现跑 **3/3 PASS**（0.147s／0.262s 两次）；raw 被 452 在飞件挡在编译（§④.0） | ✓ | 结果对上；今天的 raw 不可复跑，原因具名 |
| P17 | `:54` "`go test ./internal/panel -count=1` 改前=改后恰 1 红 `TestC21DesignTokensFourWayAgree`（ball 族既存非本面）" | 今天 C21 **绿**：`go test ./internal/panel -run 'TestC21DesignTokensFourWayAgree' -count=1 -v` ⇒ `--- PASS`（`tokens_fourway_test.go:899` 打 `checks 1-4 reconciled 22 of 22 mapped pairs`，日志 `13-c21-current.txt`） | ≈（`I5`） | 时刻读数，后修 |
| P18 | `:54` "改动件 CR=0" | `cmd/wisp/run.go`、`cmd/wisp/panel_plugins_215_test.go`、`internal/panel/plugins_inventory_215.go` 的 CR 字节读数均 **0** | ✓ | 字节级核过 |
| P19 | `:54` "运行进程从此发 `plugins` 键（215a 线上面收口）" | 两条生产泵都有读者：`cmd/wisp/run.go:1307`（流水线）＋`cmd/wisp/panel_resident_pump.go:144`（常驻）；常驻那条由 `Test65ResidentPumpCarriesConfigTruth` 行为钉住 | ✓（常驻）／**静态**（流水线） | 流水线那条只有字面钉（§④.2/④.4/④.5） |
| P20 | `:52`/`:54` 证据件两枚 | `docs/evidence/s1/215-skills-plugins-inventory.md`（62 行，`5341ead7a` 落 60 行＋`:53` 补录 2 行）、`docs/evidence/s1/215-plugins-assembly.md`（56 行），`git ls-files` 均**已入库** | ✓ | 两件真在 HEAD 链 |
| P21 | 票面 `:54` 全篇只记**一处**装配根接线（`cmd/wisp/run.go`） | 盘上**两处**：`cmd/wisp/panel_resident_pump.go:144` 由 `a5f84886f`（2026-11-08 10:22，65 常驻泵腿）落，**晚于** `8d49dc562e`（09:47）35 分钟，且 `git merge-base --is-ancestor 8d49dc562e a5f84886f` ＝ 0（是后代） | **≠（`I1`）** | 票面不是错——落那一刻只有一处；**但今天照票面找会漏掉常驻那一处**，而它恰恰是唯一有行为牙的赋值点。收件判断必须看两处 |

**不一致处点名（回话第一行要用）**：
1. **`I1`（要命）**：票面只记 `cmd/wisp/run.go` 一处 `Plugins` 赋值点；盘上还有 `cmd/wisp/panel_resident_pump.go:144`（`a5f84886f`，65 腿后落）。**唯一有编译视图行为牙的正是票面没记的那一处。**
2. **`I4`**：票面 `grep -c Plugin internal/config/schema.go`＝8；立票提交上 18、HEAD 19，**8 从未复现**（`Skill`＝0 复现）。
3. **`I3`**：`run.go:1104` → HEAD `run.go:1307`（落盘时确为 1104，已核 `git show 8d49dc562e:…`）。
4. **`I2`**：`panel.ts:372-384` → HEAD `panel.ts:436-448`／声明行 `:516`。
5. **`I5`**：`internal/panel` 的"C21 恰 1 红"与"红 3→1"是时刻读数；C21 今天绿。
6. **`I6`**：`unwired guard 表 5 行` → HEAD `unwiredKeys` 3 行（`7bcc2ae08`／票 163 删 2 行，非 215）。
7. **`I7`**：`Options 族仅 BallStyle/ChatSwitch` → HEAD 7 枚（413／appearance 波次加宽；无 plugins 专用 reader 这一半仍真）。

---

## ⑥ 边界（具名，不许被读成别的）

1. **没有 commit**：本腿写面只有 `docs/evidence/s1/215-receipt.md` 与 `.scratch/wisp/probes/215/receipt/**`（均 untracked）。本仓现存的裁决件（`docs/evidence/s1/195-r1-direction-verdict.md`、`docs/evidence/s1/351-ac2-ac3-verdict.md`、`docs/evidence/s1/314-write-v1.md`）同样是 untracked，由编排者随收件批次落；本腿照此，**不自行 commit**；⛔ push 全程零次。
2. **`cmd/wisp` 归写腿 452，我只读**：没跑 `go test ./cmd/wisp` 全包（派单禁忌），只跑 `-run 'Test215'` 单组与 `-run 'Test65ResidentPumpCarriesConfigTruth'` 单枚，**且都在 overlay 隔离 452 在飞件之后**；raw 编译失败的三条错**不是我的判断对象**，只是我为什么必须隔离的事实（§④.0）。
3. **〔未跑包〕**：`go test ./cmd/wisp` 全包**未跑**（452 在飞＋派单禁忌）；`d22scan`／`gofumpt`／`go vet`／`go build` 本轮**未复跑**（不需要：本件不改产码，派单只要求摘接线与点名红；票面的门禁读数按"时刻自述"处理，不据此判 AC）。
4. **票面前置句（票面 `:46-48` 要求每票交件必须具名写）**：**界面点不动**——网页→Go 那一跳今天只有命令行路（账 `A410`）；213/214/215 三票终态都是「**Go 侧可枚举＋命令行可执行＋界面看不到**」。本件判的 5 格**全部**在这个前提下成立，⛔ 不许被读成"owner 菜单里能用了"。
5. **AC 框与 Status 零触碰**：票面 `- [ ]` 五枚全保持未勾（⛔ 派单禁改票面）。
6. **`manifestslot` 的读法**：它存在且有包内用例，但**不在生产装配上**；AC#4 判"不满足"是判"本票没交这一格"，⛔ 不是判那个包有缺陷。
7. **"8"这件事的读法**：`I4` 是**数字腐烂**，不是维度翻转——`[plugins]` 确实有配置维、技能确实没有；如果谁要引用票面这一行，请引用"Plugin 有／Skill 0"的结论，别引用 8。

---

## ⑦ 复跑清单（下一腿/编排者逐条可重放）

```powershell
# 0) 锚点与脏源（先看，别把 452 的在飞件当我改的）
git rev-parse HEAD
git status --porcelain -- cmd/wisp internal/panel internal/config

# 1) 赋值点＋生产实现处（§③ 1/2/5）
git grep -n 'Plugins: func() panel.PluginInventorySection' -- 'cmd/wisp/*.go'
git grep -n 'PluginInventorySection' -- 'cmd/wisp/*.go' 'internal/panel/*.go'
git grep -n 'c.Plugins\|cfg.Plugins' -- 'cmd/*.go' 'internal/**/*.go' ':!*_test.go'

# 2) 基线与突变（overlay 工件都在 .scratch/wisp/probes/215/receipt/）
go test ./cmd/wisp -run 'Test215' -count=1 -v -overlay=<abs>/overlay-isolate.json
go test ./cmd/wisp -run 'Test215' -count=1 -v -overlay=<abs>/overlay-mut-no-plugins.json      # 0 红
go test ./cmd/wisp -run 'Test65ResidentPumpCarriesConfigTruth' -count=1 -v -overlay=<abs>/overlay-mut-resident-pump.json  # 红 :134
go test ./internal/config -run 'TestEveryLockedSectionKeyIsAccountedFor' -count=1 -v -overlay=<abs>/overlay-mut-unwired-row.json  # 红 :372

# 3) 盘上型红（补证，非 overlay）：见 §④.5 的 sandbox 目录（mutant run.go + 已编译测试二进制）
# 4) 树零改动
git status --porcelain -- internal/panel internal/config    # 应空
Get-FileHash -Algorithm MD5 cmd\wisp\run.go                 # B7D7EBC0534528EFAF3BA1B3B287B327
```

**归档路径**（全部 untracked，可整目录留档或忽略）：

| 路径 | 是什么 |
|---|---|
| `.scratch/wisp/probes/215/receipt/00-git-status-cmdwisp-before.txt`、`21-git-status-cmdwisp-after.txt` | 起手/收尾 `git status --porcelain -- cmd/wisp` |
| `.scratch/wisp/probes/215/receipt/01-baseline-test215.txt` | raw 编译失败（452 在飞）原样 |
| `02b`／`03`／`05`／`09`／`14`／`15` | 隔离基线、摘格、改内容、常驻泵摘格、摘句、摘处置行 |
| `04-control-overlay-applies.txt` | 控制发（证明 overlay 生效） |
| `06`／`07`／`08` | 测试二进制构建、沙箱 cwd 的盘上型红、另两枚绿 |
| `10`／`11`／`12`／`13` | unwired 组、panel 载体组、composer recon、C21 现读 |
| `20-citations.txt` | §③ 全部命令原文＋输出 |
| `head/`、`mut-*.go`、`overlay-*.json`、`sandbox/` | 突变体与 overlay（可逐字重放） |

**撤销口令**：本件与 probes 全为 untracked，且**没有**任何产码/票面/台账被改；不采纳就地忽略这两处路径即可，无需回滚任何提交（`git status --porcelain` 里属于我的行＝0；`cmd/wisp` 的那四行是写腿 452 的在飞件）。
