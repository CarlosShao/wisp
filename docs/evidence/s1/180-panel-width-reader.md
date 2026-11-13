# 票 180 落地腿 `180-w1` — `[panel] width` 归属复核＋AC#4 同族再普查（2026-11-13）

- 性质：**写腿**（票面四件套：证据件＋票面 Progress log＋commit；本程**产码零字节**——AC#2 的接线已由票 255/33-r3 在盘，本程的活是**复核、判据、反向、普查**）。
- 起手 HEAD：`b5193c097`（2026-11-13，A1820 收口那枚）。共享工作树：起手时已有 3 处**他人**未提交改动
  （`.scratch/wisp/probes/291/mic/resident_audio_291_probe_windows_test.go`、
  `.scratch/wisp/probes/orch/20261112/txt/346-a2-dispatch.txt`、`frontend/fixtures/composer-states.html`）——
  全部在本程写面之外，一字未碰；本程全部门禁读数取自这棵共享树，如实记录。
- 台件（只建不删，与本件同程入库）：`.scratch/wisp/probes/180/w1/{census-ruler.ps1, summary.tsv, raw-ruler.txt}`。

## 1. AC#1 三问（全部本程现跑，起手 HEAD `b5193c097`）

### ① 这字段今天有没有读者？——**有**。票面标题那句"生产码零读者"**已不成立**（11-09 复验节判对了，本程独立复核）

字段本体（`internal/config/schema.go:586-587` 逐字）：

```go
	// Width is the panel width in px.
	Width int `toml:"width" default:"640"`
```

读者链（每跳 file:line 本程读毕原文件核对）：

1. **读**：`cmd/wisp/panel_resident_windows.go:231`（在 `panelGeometrySource` 闭包内）
   逐字 `return cfg.Panel.Width, cfg.Panel.Height`——每次建窗从盘上现读一份 config.toml；
2. **喂**：`cmd/wisp/panel_host_windows.go:354-383` `(*PanelManager).windowOptions()`——
   `gw > 0` 时 `width = gw`，`0/负值`落回宿主常量；
3. **建**：`cmd/wisp/panel_host_windows.go:378-382` `webview2.WindowOptions{... Width: uint(width) ...}`；
4. **建窗调用点**：`newPanelControl`（测试尺报 `panel_geometry_255_test.go:267` 逐字
   `the one create site is newPanelControl at 625, fed by m.windowOptions()`——全包唯一建窗点）。

另有**写侧读者**（外观波落的）：`internal/config/settings_appearance_wave.go:86` `base.Panel.Width = px`（守卫写
`SetPanelWidth`）与 `:168` `v.Width = c.Panel.Width`（读面 `PanelWidthOptions`，喂面板设置页快照）；
写入入口 `cmd/wisp/panel_config_store.go:376` `s.mgr.SetPanelWidth(v)`。

对票面原尺的复核：尺 `grep -rn "Width" --include=*.go internal/ cmd/`（排除 `_test.go`）粗命中 25 枚，
其中真读 3 枚（`:231`／settings_appearance_wave 两处），其余为 `ball_windows.go` 的 `wr.width()`（球窗，
同名不同物）、`panel_control_windows.go:391`/`panel_host_windows.go:914`（读的是**宿主自己的 options**，
其值源自 `windowOptions()`，同链下游）、config 包内 setter/读面、名册字符串与注释——**粗尺高报，段内判读为准**。

### ② 规格要求它生效吗？——**要求，且今天已兑现**

- `docs/PLAN.md:2736`（D36 三档定义，逐字）：
  `**三档生效级别**：\`hot\`（立即生效）· \`reload\`（需重载子系统，如换 ASR 模型）· \`restart\`（需重启进程）。`
- `docs/PLAN.md:2753`（D36 配置树 `[panel]` 行，逐字；**已比票面引用漂移并扩容**——现含 `font_size`(12)）：
  `| \`[panel]\` | \`enabled\` \`width\` \`height\` \`keep_alive_in_session\`(true) \`scale\` \`font_size\`(12) | \`hot\` |`
- `docs/specs/SPEC-03-config-secrets-envs.md:39`（逐字）：
  `| \`[panel]\` | \`enabled(bool)=true\` \`width(int)=640\` \`height(int)\` \`keep_alive_in_session(bool)=true\` \`scale(float)\` | hot |`
- **派单点名的 PLAN.md:1531/:1532 两行 DEFERRED 与本票无关**：本程读毕两行原文——
  `:1531` 是「**剪贴板历史**（D34）」、`:1532` 是「**\`doc.read\` 的 xlsx 解析与 OCR**（D34）」，
  全文 3619 行内 `width` 仅 ：2753（配置树）与 ：3461（SVG `stroke-width`）两处命中，
  **不存在任何"面板宽度已推迟"的规格行**。⇒ 判「接读者」支时无"规格今天不要求"的出口。

### ③ 归属判定：**接读者那一支——而且已经接上了，不是本程接的**

- 落地来源：票 **255 AC#4**（票面 11-09 复验节记提交 `908c6652c`，33-r3），判据件
  `cmd/wisp/panel_geometry_255_test.go`（457 行，四把尺）与 winlive 台件
  `cmd/wisp/panel_geometry_255_winlive_test.go`（真窗读数）、`cmd/wisp/panel_reshow_255r1_windows_test.go`。
- 票面 AC#1 已勾的 09-28 旧结论（"这棵树里根本没有面板窗口"）**当场作废两层**（11-09 复验节已判，
  本程三把尺独立复跑同样结果）：`internal/panel/pump.go:15-17` 的 "no WebView2 host" 注释仍在场但**已过期**
  （宿主在 `cmd/wisp/panel_host_windows.go`；`go.mod:9` 逐字 `github.com/jchv/go-webview2 v0.0.0-20260205173254-56598839c808`，
  且 `go.mod:24` 有 replace 到本地 fork）；`internal/panel/` 包内搜建窗得 0 是**射程错配**，不是事实。
- **今天的真实缺口**（票 255 AC#4 与 10-08 普查已登记，本程复核仍在）：`[panel] width` 生效的时态是
  **"下次建窗/下次 Show"**——`residentPanelGeometryNote`（`panel_resident_windows.go:201-203`）自己写明
  「本腿不轮询 config.toml，也没有任何东西在文件被保存那一刻去按这一下」； roster
  `cmd/wisp/config_readers_255.go:310` 同句。**"保存即生效"没有实现**，但这属于票 255 AC#4 的"重新生效"格的
  余量，不是"零读者"。

### AC#2 正向判据（本程复跑，GREEN）

票 255 的判据件就是票面 AC#2 要的"两个不同值⇒面板侧读到的尺寸随之变"的**常驻判据**：

- `TestTicket255AssemblyRootGeometrySourceReachesTheWindowOptions`：
  config.toml 写 `width=517` ⇒ 建窗几何 517；**改写 `width=733` ⇒ 再读得 733**（断言逐字
  "after [panel] width was rewritten to 733 the host still hands the create %d"——它还专门杀"快照形状"：
  若哪天有人把闭包换成装配时快照，第二次读还是 517，此尺变红）。不是"字段非零"判据。
- `TestTicket255WindowOptionsFollowTheConfigSource`：517/331、0 高度、负值、无源四形表驱动。
- `TestTicket255PanelHostBuildsItsWindowOptions`：消费端钉——全包唯一建窗点、建窗首参必须是
  `windowOptions()` 调用、`WindowOptions` 块内不许出现 int 字面量（fallback 必须走常量名）。

本程复跑读数（sherpa staging 垫 PATH 后）：**五枚测试全 PASS**（原始输出存本程会话记录；
`ok github.com/CarlosShao/wisp/cmd/wisp 0.230s`）。

### AC#3 反向判据（本程现打，与 AC#2 成对）

**不许改默认值交差**——本程对 `internal/config/schema.go` 的 diff 为**空**
（`git diff HEAD -- internal/config/schema.go` ⇒ 空输出，`default:"640"` 原样在 :587）。
改用**摘线突变**证明判据不是恒绿：

- **M1（作废，如实记）**：`-replace` 后用三参 `WriteAllText` 落盘失败 ⇒ 文件未变，测试绿——
  **这发测的是未突变代码，读数无效**，重打。
- **M2（作废，如实记）**：突变成功但替换文删掉了 `cfg` 的唯一使用 ⇒ 编译红
  （`declared and not used: cfg`）——**这是编译错不是语义红**，重打。
- **M3（有效）**：`return cfg.Panel.Width, cfg.Panel.Height` → `_ = cfg` + `return 1, 1`，逐字读数：
  `--- FAIL: TestTicket255AssemblyRootGeometrySourceReachesTheWindowOptions`
  `panel_geometry_255_test.go:152: the create block got 1x1 from a config.toml carrying width=517 height=331 - AC#4's whole claim is that these two are the same number`
  ⇒ **语义红**（值没到建窗几何，不是编译失败）。随后原字节还原：
  还原前 sha256 `C4E79861963AF52A3EA170F1331EBBC2B849D2E34BD32D91F5F11F0E91F4EF19`，
  还原后同值（三次实测一致）⇒ 突变文件零残留（`git diff HEAD -- cmd/wisp/panel_resident_windows.go` 空输出）。
- **AC#3 的另一半**（"若判不要求生效须响亮登记"）不适用：规格要求生效且已生效（§1②③）。
- 真窗那一发（`panel_geometry_255_winlive_test.go` / `panel_reshow_255r8_winlive_test.go`）本程**没跑**：
  要开真 WebView2 窗、占常驻线程池，共享树上还有别腿在飞（票面 10-08 收口节同款边界）；
  该尺在盘、GREEN 读数由票 255 交付时留档，本程不补跑也不据此宣称真窗读数。

## 2. AC#4 同族再普查（主产出）

**母体**（本程现跑）：`grep -c 'default:' internal/config/schema.go` ＝ **71 行**，其中 `schema.go:8`
是讲 `default:"..."` 写法本身的文档注释 ⇒ **字段母体 70 枚**；去重名 63（本程 ruler 喂名 63）。
尺：对每枚字段名跑 `\.名\b`，射程 `internal/ cmd/`、排除 `_test.go`，逐枚全命中行落
`probes/180/w1/raw-ruler.txt`（1181 行）；段内歧义名另跑段内立尺（App\.Language、Voice\.TTS\.Provider 等）。
**判类口径**（沿 11-10 普查）：R＝`internal/config` 包外有非测试生产路径按值做事（boot/常驻/live 都算）；
G＝"零读者但不说谎"（写非默认值＝加载报错或守卫拦）；Z＝其余＝**改了白改**。

### 总账

| 判类 | 11-10 普查 | **本程现量** |
|---|---|---|
| R 有生产读者 | 33 | **48** |
| G 守卫拦截 | 7 | **8** |
| Z 零生产读者 | 30 | **7** |
| 合计 | 70 | **70** |

### Z 类 7 枚（AC#4 的交付对象；归口＝11-10 普查 §2 已逐枚写，本程只更读数与判类，处置归编排者）

| # | 字段（toml 键） | default | 粗尺命中 | 真读数 | 本程判读依据 |
|---|---|---|---|---|---|
| 1 | `app.language` | zh-CN | 12 | **0** | 唯一 config 侧命中＝manager.go:379 重启档比较＋tiers/名册字符串；speech/web_search 的 `.Language` 是别家字段（cloudasr.go:515、web_search.go:130） |
| 2 | `app.theme` | dark | 6 | **0** | 全部在 config 包内：manager.go:378/383 比较/拷贝、validate.go:51 枚举、tiers.go:59、roster 字符串；roster :323 自己写"扫描零命中：nothing outside internal/config reads cfg.App.THEME" |
| 3 | `ball.hide_on_fullscreen` | true | **0** | **0** | 全树（含测试）仅 schema 定义＋boundary_test 赋值＋341 注释；`internal/ball` 无任何 fullscreen/HMONITOR 判定（本程 grep 复核 8 命中全是 monitorFromWindow 别用）——**功能半边不在，想接都没对象** |
| 4 | `voice.wake_word.veto_words` | 取消,停下,别 | 2 | **0** | 仅 manager.go:434/441 段拷贝/比较；KWS 侧 `kwsVetoWords` 空表（票 423 已记） |
| 5 | `voice.tts.provider` | local-sherpa | 2（段内） | **0** | 仅 manager.go:423 段拷贝；云 TTS speak 半边未实现（423） |
| 6 | `voice.tts.speed` | 1.0 | 13 | **0**（死读点） | 唯一 config 读点 `internal/speech/cloudcatalog.go:119 [ep.Speed = cfg.Voice.TTS.Speed]` 住在 `ResolveCloudChain` 里，该函数全仓非测试调用者 0（roster :333 同判，本程复核） |
| 7 | `observe.slo_sample_interval_sec` | 5 | **0** | **0** | 全树（含测试）仅 schema 定义＋boundary_test 赋值；SPEC-03 §4.3 冻结 watchdog 1s 双职责共享节拍（423 记 DEFERRED 待裁） |

### R 类 48 枚与 G 类 8 枚（全表；粗尺命中＝summary.tsv，真读数＝段内判读）

| 字段（toml 键） | default | 粗尺 | 真读 | 读点（票） |
|---|---|---|---|---|
| `app.autostart` | false | 13 | 2 | autostart_windows.go:62,65（43） |
| `app.single_instance` | true | 3 | 1 | single_instance_boot_180p.go:69 |
| `ball.size` | 56 | 68⚠ | 1 | ball_size_413_windows.go:48（413） |
| `ball.style` | glass | 53⚠ | 3 | ball_style_65_windows.go:50＋config_cmd_65.go:150,235（65） |
| `ball.opacity_idle` | 0.35 | 10 | 3 | ball_idle_appearance_windows.go:56＋settings_appearance_wave.go:61,128 |
| `ball.click_through` | true | 7 | 3 | 同上 ：56＋settings_appearance_wave.go:71,149 |
| `hotkey.cancel` | Esc | 124⚠ | 3 | resident_windows.go:73 → ball.ApplyHotkeyDefaults（258/64） |
| `session.warm_timeout_sec` | 90 | 3 | 1 | resident_voice_windows.go:1874 |
| `session.settling_sec` | 3 | 2 | 1 | 同上 ：1879（423） |
| `session.conversation_idle_sec` | 30 | 2 | 1 | 同上 ：1878（423） |
| `voice.wake_word.enabled` | false | 84⚠ | 1 | resident_voice_windows.go:1927 |
| `voice.asr.provider` | local-sherpa | 117⚠ | 1 | 同上 ：1866 |
| `voice.aec.enabled` | true | 84⚠ | 1 | 同上 ：1884（180K） |
| `voice.realtime.enabled` | false | 84⚠ | 2 | panel_config_store.go:117（313, live）＋catalog.go:126 校验 |
| `voice.realtime.provider` | "" | 117⚠ | 1 | panel_config_store.go:118（313, live） |
| `voice.punctuation` | true | 6 | 1 | resident_voice_windows.go:1880（423） |
| `voice.conversation_mode` | false | 12⚠ | 1 | 同上 ：1883（423-K）；statemachine 12 命中是 Facts 同名字段 |
| `audio.input_device` | default | 7 | 2 | resident_audio_windows.go:307＋panel_config_store.go:103（423-K） |
| `audio.mic_muted_default` | true | 8 | 2 | resident_audio_windows.go:291,314（247） |
| `llm.roles.*.temperature` | 0.7 | 9 | 2 | prompt.go:252＋resolver.go:203（370） |
| `llm.roles.*.thinking_intensity` | off | 16 | 3 | prompt.go:254＋resolver.go:204＋anthropic/openairesponses request.go（370） |
| `llm.retry.max` | 3 | 27⚠ | 2 | run.go:2520＋llm/retry.go:79,192 |
| `llm.retry.backoff_ms` | 1000 | 1 | 1 | run.go:2519 |
| `compat.loose` | false | 39 | 8 | 三家 adapter（anthropic:315 / openaichat:274 / openairesponses:275）＋resolver.go:69 |
| `compat.allow_missing_usage` | false | 7 | 6 | 同上三家＋resolver.go:70 |
| `providers.<n>.billing` | pay-per-token | 25⚠ | 3 | quota_44r4.go:115＋validate.go:274＋settings_320.go:48（44/320） |
| `models.<id>.billing` | pay-per-token | 25⚠ | 2 | quota_44r4.go:121＋validate.go:311（44） |
| `models.<id>.enabled` | true | 84⚠ | 3 | resolver.go:128,298＋catalog.go:126 系（187） |
| `llm.timeout_ms` | 60000 | 10 | 1 | run.go:2137 |
| `agent.loop_guard.repeat_thresholds` | 3,5,8 | 8 | 2 | run.go:2022 → loop.go:245 → guard.go:90（423） |
| `agent.max_rounds` | 50 | 8 | 2 | run.go:2020 → loop/guard（423） |
| `agent.steering_enabled` | true | 2 | 2 | run.go:2011＋loop.go:318 |
| `agent.project_instructions_enabled` | true | 2 | 2 | run.go:2106,2111（200） |
| `risk.confirm_timeout_sec` | 300 | 3 | 2 | run.go:1058＋resident_approval_windows.go:557（256） |
| `risk.l1_window_sec` | 2 | 4 | 2 | run.go:1057＋resident :556；manager.go:467 方向审计 |
| `risk.shell_enabled` | false | 3 | 1 | run.go:864 → ShellDeps（163） |
| `risk.permission_mode` | ask_every_step | 31⚠ | 3 | permmode.go:40＋validate.go:219＋tools/bridge.go（90） |
| `fs.delete_enabled` | false | 9 | 1 | run.go:780 → fs_write.go:596 |
| `net.proxy.mode` | system | 118⚠ | 3 | providers.go:138＋panel_models_320.go:110＋run.go:2517 |
| `privacy.redact_paths` | false | 11 | 3 | logsink.go:233（boot）＋diagnostics_cmd.go:157,166 |
| `privacy.retention_days` | 30 | 3 | 2 | memory_retention_180k.go:45-46（180K） |
| `memory.l1_enabled` | true | 3 | 1 | memory_retention_180k.go:62-65（180K） |
| `memory.l1_max` | 20 | 2 | 1 | run.go:691（180K） |
| `memory.l3_retention_days` | 30 | 3 | 2 | memory_retention_180k.go:42-43（180K） |
| `panel.enabled` | true | 84⚠ | 1 | panel_resident_windows.go:258（33-r3） |
| `panel.width` | 640 | 25⚠ | 3 | panel_resident_windows.go:231（255）＋settings_appearance_wave.go:86,168＋panel_config_store.go:376 |
| `panel.keep_alive_in_session` | true | 1 | 1 | panel_resident_windows.go:276（33-r3） |
| `panel.font_size` | 12 | 4 | 2 | run.go:1255／panel_resident_pump.go:206 → PanelFontSizeOptions＋setter（外观波） |
| `cost.alert_threshold` | 0.8 | 4 | 2 | costbook_44.go:184＋quota_44r4.go:136（44） |
| `cost.over_budget` | pause | 4 | 2 | 同上两处 |
| `plugins.tier2_enabled` | false | 14 | 4 | run.go:1300＋panel_resident_pump.go:251＋tier2_gate.go:154,226 |
| `observe.level` | info | 128⚠ | 2 | logsink.go:226＋validate.go:245；128 命中里 100+ 是 risk.Level/审批卡撞车 |
| `observe.roll.size_mb` | 10 | 2 | 1 | logsink.go:231（180K） |
| `observe.roll.days` | 7 | 8⚠ | 1 | logsink.go:232（180K）；observe/costmeter 的 `.days` 是别家 |
| `voice.aec.echo_ref` | self-render | 1 | 0 | **G**：validate.go:148 词汇表守卫（写错值＝加载报错，423-K 落） |
| `audio.sample_rate` | 16000 | 19⚠ | 1 | **G**：validate.go:120 必须 16000（423-K）；19 命中是引擎 SampleRate() 撞车 |
| `audio.half_duplex` | true | 1 | 0 | **G**：validate.go:116（写 false 报错，83 先例） |
| `risk.allow_shell_string` | false | 3 | 0 | **G**：unwired.go:74-77 |
| `net.block_private_ranges` | true | 3 | 0 | **G**：unwired.go:91-94 |
| `privacy.keep_transcript` | false | 1 | 0 | **G**：validate.go:159 |
| `privacy.keep_audio` | false | 1 | 0 | **G**：validate.go:163 |
| `models.verify_signature` | true | 2 | 1 | **G**：validate.go:173＋C29 字面量 true（值设计上不消费） |

⚠＝粗尺同名高报（Enabled 84／Provider 117／Mode 118／Level 128／Size 68／Style 53／Max 27／Days 8 等），
真读数一律以段内立尺或票 255 名册为准——**这是 10-08 普查"票面单跑尺 A 会乐观 3 倍"的同款纪律**。

### 与 11-10 普查（票 423 的来路）对账：**Z 30 → 7，差 23 枚全部有票可指**

| # | 字段 | 11-10 判 | 本程判 | 谁接的 |
|---|---|---|---|---|
| 1 | `session.settling_sec` | Z（悬空） | R | 票 423 收口腿 |
| 2 | `session.conversation_idle_sec` | Z | R | 票 423 |
| 3 | `voice.punctuation` | Z | R | 票 423 |
| 4 | `agent.loop_guard.repeat_thresholds` | Z | R | 票 423 |
| 5 | `agent.max_rounds` | Z | R | 票 423 |
| 6 | `agent.token_budget` | Z | R | 票 423（无 default 标签，不进本程 70 母体，仍记对账） |
| 7 | `observe.level` | Z | R | 票 423-K（logsink 自喂） |
| 8 | `observe.roll.size_mb` | Z（raw=0） | R | 票 423-K |
| 9 | `observe.roll.days` | Z | R | 票 423-K |
| 10 | `privacy.redact_paths` | Z | R | 票 423-K |
| 11 | `privacy.retention_days` | Z | R | 票 423-K |
| 12 | `memory.l1_enabled` | Z | R | 票 423-K |
| 13 | `memory.l1_max` | Z | R | 票 423-K |
| 14 | `memory.l3_retention_days` | Z | R | 票 423-K |
| 15 | `voice.conversation_mode` | Z（撞车） | R | 票 423-K |
| 16 | `voice.aec.enabled` | Z | R | 票 423-K |
| 17 | `audio.input_device` | Z | R | 票 423-K |
| 18 | `audio.mic_muted_default` | 11-10 已记 R | R | （一致） |
| 19 | `voice.aec.echo_ref` | Z | **G** | 票 423-K（词汇表守卫——从"零读者"改判"不说谎"） |
| 20 | `audio.sample_rate` | Z（撞车） | **G** | 票 423-K（必须 16000 守卫） |
| 21 | `app.single_instance` | Z（悬空） | R | **新发现**：`single_instance_boot_180p.go:69` 真读（11-10 表漏了这枚读点） |
| 22 | `privacy.diagnostics_opt_in` | Z（423 记 DEFERRED） | **R（新发现）** | `wisp diagnostics preview/export` 的总闸 `diagnosticsGate180P`（diagnostics_cmd.go:137 逐字 `if !c.Privacy.DiagnosticsOptIn {`→拒绝并指名）——423 的"BuildDiagnosticsBundle 零生产调用点"理由**只对 GUI/自动采集那半成立，CLI 半边有真读者** |
| 23 | `ball.opacity_idle`／`click_through` | Z（挂票 62/68） | R | 外观波（settings_appearance_wave.go setter+读面＋ball_idle_appearance_windows.go:56 boot 喂线）——11-10 表只算了 validate，漏了喂线那半 |

**Z 类剩余 7 枚的归口状态**：`app.language`／`app.theme` → 📋 票 39（open，GUI 配置编辑器）；
`ball.hide_on_fullscreen` → 📋 票 62/68（open）——本程新量出**功能半边不在**（internal/ball 无 fullscreen 判定），
接线前先要有"谁判全屏"；`voice.wake_word.veto_words` → 票 26 done 悬空；`voice.tts.provider`/`voice.tts.speed`
→ 票 61/26 done 悬空；`observe.slo_sample_interval_sec` → 423 记 DEFERRED 待裁。
**开新票还是挂台账＝编排者裁，本腿不裁**（票 423 的处置格原样有效，本程只把 25 缩到 7、
并给每枚配了今天可复跑的读数）。

## 3. 门禁五条读数（commit 前工作树现跑；最后一枚 commit 之后另有一发终态复跑，读数随腿回话记录）

1. `go build ./cmd/wisp/` ⇒ **rc=0**（无输出）。
2. `go vet ./cmd/wisp/` ⇒ **rc=0**（无输出）。
3. `go test -count=1 ./internal/config/ ./internal/panel/` ⇒
   `ok github.com/CarlosShao/wisp/internal/config 2.187s` ＋
   `ok github.com/CarlosShao/wisp/internal/panel 15.533s`，**rc=0 全绿**
   （注：10-08 普查时 internal/panel 那 2 枚红——frontend 主色与缺失的 design/assets/tokens.css——
   本机今天都不红，读数如实记录）。
4. `gofumpt -l internal cmd tools scripts frontend` ⇒ **空输出**（gofumpt v0.12.0，rc=0）。
5. `sh scripts/d22scan.sh` ⇒ 末行逐字
   `d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=403, bans #1-5 cmd/=77, ... ban #8 internal/=1004, ban #8 cmd/=276`（rc=0）。
   ⚠ `./cmd/wisp/` 的包测试按票 98 只在判据/突变时定向跑，不进整包红账。

## 4. 本程没做（明写，不写成"以后加固"）

- **没改任何产码**：AC#2 的接线在盘（票 255），本程的"落地"＝复核＋判据复跑＋反向突变＋普查；
  `git diff HEAD -- internal/config/schema.go cmd/wisp/panel_resident_windows.go` 双空读数在 §1/§3。
- **没跑真窗 winlive**（§1 末）；没有真窗像素读数，不据 AC#2 宣称"屏幕上量到了"。
- **没做 7 枚 Z 类的处置**（开票/挂台账/删字段＝改票面/契约，归编排者；`ball.hide_on_fullscreen` 若要接，
  先要有"全屏判定"的生产对象，那是票 62/68 的地界）。
- **没读未引 `frontend/**`／`design/**`**（起手树里 `frontend/fixtures/composer-states.html` 有他人未提交改动，
  本程零触碰）；"零读者"与普查射程只到 Go 侧 `internal/`＋`cmd/`。
- 突变探针打歪两发（M1 落盘失败／M2 编译红）**原样记录未洗**；有效那一发 M3 的还原 sha256 三方一致。

## 5. 一句结论

票 180 的标题今天只剩历史价值：`[panel] width` 有规格（D36 hot）、有读者（`:231`→`windowOptions()`→建窗）、
有守卫写侧（`SetPanelWidth` 拒非正数）、有能杀"快照退化"的常驻判据（摘线即红）。
同族哑键从 11-10 的 30 枚降到 **7 枚**，其中 2 枚全树零提及（`hide_on_fullscreen`／`slo_sample_interval_sec`）、
2 枚死读点/半边缺（`tts.speed`／`tts.provider`）、1 枚 KWS 空表（`veto_words`）、2 枚整段无包外读者
（`app.language`／`app.theme`，挂票 39）。AC#4 的下一程只剩这 7 枚的处置。
