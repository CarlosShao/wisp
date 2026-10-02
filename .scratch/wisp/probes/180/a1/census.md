# 180-a1（10-02 第二程）—— 配置面全名册现量：哑键普查 · `unwiredKeys` 覆盖率 · 生效级别三档 · "此项不生效"的出口 · `[panel] width` 断链点

派单：只读取证腿 `180-a1`，10-02 09:5x 起手。**产码零字节改动**；票面 AC 框一枚未碰；票面正文未改（只在末尾**追加**一节）。

---

## §0 起手锚（逐字）

- `date -Iseconds` ⇒ `2026-10-02T09:57:47+08:00`
- `git log -1` ⇒ commit `00e7efeb0ca22e1ce3cbe62763d0a3088dfe4330`（"ledger(A521)：167-a1 结档…"，Author CarlosShao，Date Fri Oct 2 09:57:39 2026 +0800）
- `git rev-parse --abbrev-ref HEAD` ⇒ `dev`
- `git status --porcelain internal/config internal/panel cmd/wisp` ⇒ **空输出**（三目录起手干净）
- 全树 `git status --porcelain`：**不干净**——别人的脏件与未推临时件大量在场（`design/**` 一批 D/M、`.scratch/**` 一批 ??、`scripts/testdata/portable-tests/go` M）。本程**不还原、不提交**它们的任何东西。
- 本程读到的**所有行号均取自 HEAD `00e7efe`**（除注明）。票面锚是 `039efb47`、上一程 `180-a1` 锚是 `38fc7c0e`——**两处行号都已过期**，差集见 §8。

## §1 全名册：叶子键 → 读者 → 非测试／仅测试／无

**起手 HEAD**：`00e7efe`（本表所有 `file:line` 均在此号上读；终态复核见 §8 末条）。

### 1.1 先把"多少枚"钉住（三把独立尺）

- 段数：**18 枚 section ＋ 1 枚根表裸键 `schema_version`**。
  甲尺：`Config` struct 的段字段枚数＝18（`internal/config/schema.go:110-133`）。
  乙尺：票 198 的首建真件 `.scratch/wisp/probes/198/r1/firstrun-defaults.dump.txt`
  有 33 个表头 ＝ 18 顶层 ＋ 15 子表；它的 `key =` 行 **115** 枚。
  丙尺：`schema.go` 静态叶子走查＝**115** 枚，与首建真件逐枚 `comm` 差集只差
  `models.local_override`（空 map 不落盘）与 dump 里一行 `time`（非 schema 键）。
- 叶子键：**150 枚**＝静态 115 ＋ 动态子表模板 35
  （`llm.providers.<id>` 10 ＋ 其下 `models.<id>` 21 ＋ `plugins.<id>` 4）。
  动态子表按**模板**计数、不按实例；换口径的后果见 §6 第 1 条。
- 带 `default:` 标签的行 **69** 枚 —— 与票面 09-28 记录里 `180-c1` 报的 69 独立对上。

### 1.2 判据与三条反查尺

判定问的是"**有没有人把这个键的值读出去做事**"，不是"有没有人提到这个名字"。
射程：`cmd internal tools`（显式根，**不含** `frontend/**`、`design/**`、`testdata`），排除 `*_test.go`。
`internal/config` 包内的命中单列——那是管道（填默认、校枚举、算放宽方向、整块拷贝比较），不是消费者。

三条反查（防"名字没提但值被读走"）：
① 反射按名取值：`FieldByName|NumField()|reflect.TypeOf` 在射程内非测试码 **零命中**
   （反射只在 `internal/config/defaults.go` 里把 `default:` 标签**写进**字段）。
② 按点分路径字符串寻址：全路径字面量在 `internal/config` 之外只命中
   `cmd/wisp/config_reload.go:300`（重启档三枚键名，只用于说话）
   ＋两枚同名巧合（`internal/statemachine/events.go:21` 事件名、`internal/memory/schema.go:149` SQLite meta 键）。
③ 生产里取得 `*Config` 的入口枚数：**8 处**（`cmd/wisp/run.go:423`、`cmd/wisp/providers.go:98`、
   `cmd/wisp/models.go:184`、`cmd/wisp/panel_config_store.go:92`、`cmd/wisp/approval_always.go:155`、
   `internal/perm/store.go:159`、`cmd/wisp/firstrun.go:82`（建默认文件）、`cmd/balldebug/main.go:238`（旁路））。
   ⇒ 一切配置消费都必须经过这 8 处，这就是逐枚判定的收敛面。

### 1.3 判读图例（七档）

| 档 | 含义 |
|---|---|
| R | 有生产读者：值被读出去驱动行为，具名 file:line |
| B | 有读者但只作旁注：`internal/llm/probe_health.go:95-97` 逐字 "It exists ONLY to name mismatches; it never feeds a verdict" |
| W | 类型在场、生产零赋值：`config.Price` 被 `internal/agent/cost.go:38` 消费，但生产没有一处给它赋值 ⇒ 值恒为零价卡 |
| S | 只有旁路调试程序读者（`cmd/balldebug`），出货路径 `wisp run` 零读者 |
| G | 无人读，**已在 `unwiredKeys` 登记**：写非默认值直接加载失败（有出口） |
| L | 无人读，**只在 `lockedKeyDisposition` 登记** not built：写它不报错、也不生效（无出口） |
| D | 无人读、名册外 ⇒ 票 180 那一族的本体 |

### 1.4 枚数总账（逐枚名列 §1.5）

| 档 | 全 150 枚 | 仅静态 115 枚 |
|---|---|---|
| R | 50 | 40 |
| B | 4 | 0 |
| W | 5 | 0 |
| S | 4 | 4 |
| G | 6 | 6 |
| L | 5 | 1 |
| D | 76 | 64 |
| 合计校验 | 150 | 115 |
| 无人读／无实效（B＋W＋S＋G＋L＋D） | 100 | 75 |
| **既无人读又不在任何名册（D）** | **76** | **64** |

### 1.5 逐枚名册（150 行）

| TOML 路径 | Go 字段 | 判 | 凭据／备注 |
|---|---|---|---|
| `schema_version` | `Config.SchemaVersion` | R | internal/config/loader.go + migrate.go：包内迁移判定就是它的用途 |
| `agent.loop_guard.repeat_thresholds` | `LoopGuard.RepeatThresholds` | D | 同上：loop.go:220 兜回 3,5,8 ⇒ 改这枚数组不改变刹车阶梯 |
| `agent.max_rounds` | `AgentSection.MaxRounds` | D | agent.Config 有这枚字段（internal/agent/loop.go:115），cmd/wisp/run.go:1010 的字面量不填 ⇒ loop.go:216-217 兜回 50 |
| `agent.per_tool_timeout_ms` | `AgentSection.PerToolTimeoutMS` | R | cmd/wisp/run.go:763 与 :1014 |
| `agent.project_instructions_enabled` | `AgentSection.ProjectInstructionsEnabled` | R | cmd/wisp/run.go:1055,1060 |
| `agent.steering_enabled` | `AgentSection.SteeringEnabled` | R | cmd/wisp/run.go:1016 |
| `agent.token_budget` | `AgentSection.TokenBudget` | D | 同上：loop.go:218-219 兜回 scaled 默认，run.go 不填 |
| `app.autostart` | `AppSection.Autostart` | D | 开机自启注册码上不存在读它的地方 |
| `app.language` | `AppSection.Language` | D | SPEC-03:26 写 restart；重启那句台词点名它（cmd/wisp/config_reload.go:300）但无 boot 读者 ⇒「重启后生效」今天兑现不了 |
| `app.single_instance` | `AppSection.SingleInstance` | D | internal/proc/singleinstance_windows.go:40 用 rt.Layout.MutexName（环境布局），不是 cfg.App.SingleInstance |
| `app.theme` | `AppSection.Theme` | D | manager.go:343-348 只比较与拷贝；全仓无人据它换肤（SPEC-03:26 写 hot） |
| `audio.half_duplex` | `AudioSection.HalfDuplex` | D | 只有 validate.go 把 false 判成加载错误（硬编码只读三枚的正解形状：出口是拒绝，不是告知） |
| `audio.input_device` | `AudioSection.InputDevice` | D | 零读者（internal/audio 不读配置） |
| `audio.mic_muted_default` | `AudioSection.MicMutedDefault` | D | internal/audio/gate.go:56 注释逐字点名 "(config AudioSection.MicMutedDefault) here at boot wiring"——只有注释没有代码（注释指名的接线一律当待验断言） |
| `audio.sample_rate` | `AudioSection.SampleRate` | D | 零读者 |
| `ball.click_through` | `BallSection.ClickThrough` | D | 零读者 |
| `ball.hide_on_fullscreen` | `BallSection.HideOnFullscreen` | D | 零读者 |
| `ball.opacity_idle` | `BallSection.OpacityIdle` | D | 零读者 |
| `ball.position.monitor` | `BallPos.Monitor` | D | 同 ball.position.x |
| `ball.position.x` | `BallPos.X` | D | internal/ball/position.go 用的是自己的 Pos 类型，不是 BallPos ⇒ 拖完不存、存了不读 |
| `ball.position.y` | `BallPos.Y` | D | 同 ball.position.x |
| `ball.size` | `BallSection.Size` | D | 球尺寸实走 internal/ball/ball_windows.go:64 的 opts.SizePx（默认与钳位 :144-151），cfg.Ball 从未参与 |
| `cost.alert_threshold` | `CostSection.AlertThreshold` | D | 只有 validate.go 校枚举 |
| `cost.daily_budget` | `CostSection.DailyBudget` | D | 注释说达 100% 暂停新任务（C23 票 44），零读者 |
| `cost.monthly_budget` | `CostSection.MonthlyBudget` | D | 同 cost.daily_budget |
| `cost.over_budget` | `CostSection.OverBudget` | D | 只有 validate.go 校枚举 pause|warn |
| `fs.allowed_dirs` | `FSSection.AllowedDirs` | R | cmd/wisp/run.go:497-498 |
| `fs.delete_enabled` | `FSSection.DeleteEnabled` | R | cmd/wisp/run.go:517 |
| `fs.reparse_point_exceptions` | `FSSection.ReparsePointExceptions` | R | cmd/wisp/run.go:501 |
| `hotkey.cancel` | `HotkeySection.Cancel` | S | 唯一读者 cmd/balldebug/main.go:238；出货路径 cmd/wisp/resident_ball_windows.go:171 用 ball.DefaultHotkeys() |
| `hotkey.mute` | `HotkeySection.Mute` | S | 唯一读者 cmd/balldebug/main.go:238；出货路径 cmd/wisp/resident_ball_windows.go:171 用 ball.DefaultHotkeys() |
| `hotkey.panel` | `HotkeySection.Panel` | S | 唯一读者 cmd/balldebug/main.go:238；出货路径 cmd/wisp/resident_ball_windows.go:171 用 ball.DefaultHotkeys() |
| `hotkey.summon` | `HotkeySection.Summon` | S | 唯一读者 cmd/balldebug/main.go:238；出货路径 cmd/wisp/resident_ball_windows.go:171 用 ball.DefaultHotkeys() |
| `llm.providers.<id>.api_key_ref` | `Provider.APIKeyRef` | R | internal/llm/resolver.go:136 / cmd/wisp/providers.go:124-125 |
| `llm.providers.<id>.base_url` | `Provider.BaseURL` | R | internal/llm/anthropic/adapter.go:146、openaichat/adapter.go:81、openairesponses/adapter.go |
| `llm.providers.<id>.billing` | `Provider.Billing` | D | 只有 validate.go 校枚举 |
| `llm.providers.<id>.compat.allow_missing_usage` | `Compat.AllowMissingUsage` | R | internal/llm/anthropic/adapter.go:265,300 与 openaichat/adapter.go:224,253（三适配器各一枚） |
| `llm.providers.<id>.compat.extra_headers` | `Compat.ExtraHeaders` | R | internal/llm/anthropic/adapter.go:158 等三适配器逐条 range |
| `llm.providers.<id>.compat.loose` | `Compat.Loose` | R | internal/llm/anthropic/adapter.go:314 等三适配器各一枚 loose() |
| `llm.providers.<id>.models.<id>.billing` | `ModelSpec.Billing` | D | 只有 validate.go 校枚举 |
| `llm.providers.<id>.models.<id>.capabilities.audio_in` | `Capabilities.AudioIn` | B | internal/llm/probe_health.go:98-110 declaredCapability（注释逐字 "it never feeds a verdict"） |
| `llm.providers.<id>.models.<id>.capabilities.audio_out` | `Capabilities.AudioOut` | D | 零读者：declaredCapability 的 switch 里只有 fc/vision/thinking/audio_in 这四支 |
| `llm.providers.<id>.models.<id>.capabilities.fc` | `Capabilities.FC` | B | internal/llm/probe_health.go:98-110 declaredCapability（注释逐字 "it never feeds a verdict"） |
| `llm.providers.<id>.models.<id>.capabilities.realtime` | `Capabilities.Realtime` | D | 零读者：declaredCapability 的 switch 里只有 fc/vision/thinking/audio_in 这四支 |
| `llm.providers.<id>.models.<id>.capabilities.stream` | `Capabilities.Stream` | D | 零读者：declaredCapability 的 switch 里只有 fc/vision/thinking/audio_in 这四支 |
| `llm.providers.<id>.models.<id>.capabilities.text` | `Capabilities.Text` | D | 零读者：declaredCapability 的 switch 里只有 fc/vision/thinking/audio_in 这四支 |
| `llm.providers.<id>.models.<id>.capabilities.thinking` | `Capabilities.Thinking` | B | internal/llm/probe_health.go:98-110 declaredCapability（注释逐字 "it never feeds a verdict"） |
| `llm.providers.<id>.models.<id>.capabilities.vision` | `Capabilities.Vision` | B | internal/llm/probe_health.go:98-110 declaredCapability（注释逐字 "it never feeds a verdict"） |
| `llm.providers.<id>.models.<id>.context_window` | `ModelSpec.ContextWindow` | R | internal/llm/resolver.go:68 → 三适配器 MaxContextWindow |
| `llm.providers.<id>.models.<id>.display` | `ModelSpec.Display` | D | 注释说 human-facing name，本包外零读者 |
| `llm.providers.<id>.models.<id>.enabled` | `ModelSpec.Enabled` | D | schema.go:369-371 注释逐字 "removes the model from discovery/selection"：resolver.go:288 的 DiscoveredModels 与 :119 取 spec 都不查它 ⇒ 注释描述的规矩码上没人执行 |
| `llm.providers.<id>.models.<id>.max_output_tokens` | `ModelSpec.MaxOutputTokens` | R | internal/llm/resolver.go:134 → anthropic/request.go:103 |
| `llm.providers.<id>.models.<id>.price.audio_in` | `Price.AudioIn` | W | 类型 config.Price 被 internal/agent/cost.go:38 消费，但生产里 Price: 赋值零枚（cmd/wisp/run.go:1010 的 agent.Config{} 不填）⇒ 值永远是零价卡 |
| `llm.providers.<id>.models.<id>.price.audio_out` | `Price.AudioOut` | W | 类型 config.Price 被 internal/agent/cost.go:38 消费，但生产里 Price: 赋值零枚（cmd/wisp/run.go:1010 的 agent.Config{} 不填）⇒ 值永远是零价卡 |
| `llm.providers.<id>.models.<id>.price.cached` | `Price.Cached` | W | 类型 config.Price 被 internal/agent/cost.go:38 消费，但生产里 Price: 赋值零枚（cmd/wisp/run.go:1010 的 agent.Config{} 不填）⇒ 值永远是零价卡 |
| `llm.providers.<id>.models.<id>.price.in` | `Price.In` | W | 类型 config.Price 被 internal/agent/cost.go:38 消费，但生产里 Price: 赋值零枚（cmd/wisp/run.go:1010 的 agent.Config{} 不填）⇒ 值永远是零价卡 |
| `llm.providers.<id>.models.<id>.price.out` | `Price.Out` | W | 类型 config.Price 被 internal/agent/cost.go:38 消费，但生产里 Price: 赋值零枚（cmd/wisp/run.go:1010 的 agent.Config{} 不填）⇒ 值永远是零价卡 |
| `llm.providers.<id>.models.<id>.quota_daily_micro` | `ModelSpec.QuotaDailyMicro` | D | 注释指向 C23 pre-dispatch hook（票 44），未落 ⇒ 零读者 |
| `llm.providers.<id>.models.<id>.quota_monthly_micro` | `ModelSpec.QuotaMonthlyMicro` | D | 同 quota_daily_micro |
| `llm.providers.<id>.models.<id>.thinking_levels` | `ModelSpec.ThinkingLevels` | D | 只有 validate.go 校枚举，无人据它限制档位 |
| `llm.providers.<id>.plan_credit_total_micro` | `Provider.PlanCreditTotalMicro` | D | 注释说 plan 模式据此核算，零读者 |
| `llm.providers.<id>.protocol` | `Provider.Protocol` | R | internal/llm/resolver.go:65 → internal/llm/provider.go:219 |
| `llm.providers.<id>.rpm` | `Provider.RPM` | R | internal/llm/resolver.go:130 → :276 |
| `llm.providers.<id>.tpm` | `Provider.TPM` | R | internal/llm/resolver.go:131 → :276 |
| `llm.retry.backoff_ms` | `RetryConfig.BackoffMS` | R | cmd/wisp/run.go:1226 |
| `llm.retry.max` | `RetryConfig.Max` | R | cmd/wisp/run.go:1227 |
| `llm.roles.chat.max_output_tokens` | `Role.MaxOutputTokens` | R | internal/llm/resolver.go:197 |
| `llm.roles.chat.model` | `Role.Model` | R | internal/llm/resolver.go:119 用 provider/model 对取 spec |
| `llm.roles.chat.provider` | `Role.Provider` | R | internal/llm/resolver.go:179 取 role → :194 建 endpoint |
| `llm.roles.chat.temperature` | `Role.Temperature` | R | internal/llm/resolver.go:195 |
| `llm.roles.chat.thinking_intensity` | `Role.ThinkingIntensity` | R | internal/llm/resolver.go:196 |
| `llm.roles.handoff.max_output_tokens` | `Role.MaxOutputTokens` | R | internal/llm/resolver.go:197 |
| `llm.roles.handoff.model` | `Role.Model` | R | internal/llm/resolver.go:119 用 provider/model 对取 spec |
| `llm.roles.handoff.provider` | `Role.Provider` | R | internal/llm/resolver.go:179 取 role → :194 建 endpoint |
| `llm.roles.handoff.temperature` | `Role.Temperature` | R | internal/llm/resolver.go:195 |
| `llm.roles.handoff.thinking_intensity` | `Role.ThinkingIntensity` | R | internal/llm/resolver.go:196 |
| `llm.roles.memory_extract.max_output_tokens` | `Role.MaxOutputTokens` | R | internal/llm/resolver.go:197 |
| `llm.roles.memory_extract.model` | `Role.Model` | R | internal/llm/resolver.go:119 用 provider/model 对取 spec |
| `llm.roles.memory_extract.provider` | `Role.Provider` | R | internal/llm/resolver.go:179 取 role → :194 建 endpoint |
| `llm.roles.memory_extract.temperature` | `Role.Temperature` | R | internal/llm/resolver.go:195 |
| `llm.roles.memory_extract.thinking_intensity` | `Role.ThinkingIntensity` | R | internal/llm/resolver.go:196 |
| `llm.roles.summarize.max_output_tokens` | `Role.MaxOutputTokens` | R | internal/llm/resolver.go:197 |
| `llm.roles.summarize.model` | `Role.Model` | R | internal/llm/resolver.go:119 用 provider/model 对取 spec |
| `llm.roles.summarize.provider` | `Role.Provider` | R | internal/llm/resolver.go:179 取 role → :194 建 endpoint |
| `llm.roles.summarize.temperature` | `Role.Temperature` | R | internal/llm/resolver.go:195 |
| `llm.roles.summarize.thinking_intensity` | `Role.ThinkingIntensity` | R | internal/llm/resolver.go:196 |
| `llm.text_chain` | `LLMSection.TextChain` | R | internal/llm/resolver.go:96 |
| `llm.timeout_ms` | `LLMSection.TimeoutMS` | R | cmd/wisp/run.go:1086 |
| `memory.extract_model` | `MemorySection.ExtractModel` | D | 零读者（注释指向 roles.memory_extract，那枚真被读） |
| `memory.l1_enabled` | `MemorySection.L1Enabled` | D | 零读者；L1 开关不在配置上 |
| `memory.l1_max` | `MemorySection.L1Max` | D | 零读者 |
| `memory.l3_retention_days` | `MemorySection.L3RetentionDays` | D | 零读者；与 privacy.retention_days 同族两枚都没人清理 |
| `models.dir` | `ModelsSection.Dir` | R | cmd/wisp/models.go:192-196 |
| `models.local_override` | `ModelsSection.LocalOverride` | R | cmd/wisp/models.go:164 |
| `models.mirror` | `ModelsSection.Mirror` | R | cmd/wisp/models.go:163 |
| `models.verify_signature` | `ModelsSection.VerifySignature` | R | internal/models/downloader.go（validate.go 另把 false 判成加载错误） |
| `net.allowlist` | `NetSection.Allowlist` | G | unwiredKeys 第 5 行（:85-90）；internal/risk/rules_network.go:74 读的是 risk.NetTarget.Allowlist，同样无生产者 |
| `net.block_private_ranges` | `NetSection.BlockPrivateRanges` | G | unwiredKeys 第 6 行（:91-96） |
| `net.proxy.mode` | `ProxyConfig.Mode` | R | cmd/wisp/run.go:1224 / cmd/wisp/providers.go:138 |
| `net.proxy.url` | `ProxyConfig.URL` | R | cmd/wisp/run.go:1225 / cmd/wisp/providers.go:138 |
| `observe.level` | `ObserveSection.Level` | D | LogConfig.Level 由 cmd/wisp/logsink.go:149 的 logSinkLevel 与 slo_windows.go:274 的字面量 "info" 填，不读 cfg.Observe.Level |
| `observe.roll.days` | `RollConfig.Days` | D | observe.LogConfig.RollDays 同上 |
| `observe.roll.size_mb` | `RollConfig.SizeMB` | D | observe.LogConfig.RollSizeMB 两处生产调用点都没填 |
| `observe.slo_sample_interval_sec` | `ObserveSection.SLOSampleIntervalSec` | D | 本包外零命中；D32 采样节奏仍走 internal/observe 自己的常数 |
| `panel.enabled` | `PanelSection.Enabled` | D | 见 §5 尺：段级令牌 .Panel 在本包外零命中 |
| `panel.height` | `PanelSection.Height` | D | 同 panel.enabled |
| `panel.keep_alive_in_session` | `PanelSection.KeepAliveInSession` | D | 同 panel.enabled |
| `panel.scale` | `PanelSection.Scale` | D | 同 panel.enabled |
| `panel.width` | `PanelSection.Width` | D | 票面那一枚；断点见 §5 |
| `plugins.Entries.<id>.capabilities` | `PluginEntry.Capabilities` | L | lockedKeyDisposition 登记 not built（unwired.go:142-146），不在 unwiredKeys ⇒ 无拦截 |
| `plugins.Entries.<id>.enabled` | `PluginEntry.Enabled` | L | lockedKeyDisposition 登记 not built（unwired.go:142-146），不在 unwiredKeys ⇒ 无拦截 |
| `plugins.Entries.<id>.host_api` | `PluginEntry.HostAPI` | L | lockedKeyDisposition 登记 not built（unwired.go:142-146），不在 unwiredKeys ⇒ 无拦截 |
| `plugins.Entries.<id>.net_allowlist` | `PluginEntry.NetAllowlist` | L | lockedKeyDisposition 登记 not built（unwired.go:142-146），不在 unwiredKeys ⇒ 无拦截 |
| `plugins.tier2_enabled` | `PluginsSection.Tier2Enabled` | L | lockedKeyDisposition 登记 not built（unwired.go:141），不在 unwiredKeys ⇒ 写它不报错也不生效 |
| `privacy.diagnostics_opt_in` | `PrivacySection.DiagnosticsOptIn` | D | 零读者 |
| `privacy.keep_audio` | `PrivacySection.KeepAudio` | D | 同 half_duplex：validate.go 把 true 判成错误 |
| `privacy.keep_transcript` | `PrivacySection.KeepTranscript` | D | 同 half_duplex |
| `privacy.redact_paths` | `PrivacySection.RedactPaths` | D | internal/observe/logging.go:50 与 diagnostics.go:43 各有 RedactPaths 字段（注释逐字 "mirrors [privacy] redact_paths"），但两处生产 InitLog（cmd/wisp/logsink.go:149、cmd/wisp/slo_windows.go:272）只填 Dir+Level ⇒ 建了但没接 |
| `privacy.retention_days` | `PrivacySection.RetentionDays` | D | 零读者 |
| `risk.allow_shell_string` | `RiskSection.AllowShellString` | G | unwiredKeys 第 2 行（:66-72） |
| `risk.blacklist_overrides` | `RiskSection.BlacklistOverrides` | G | unwiredKeys 第 4 行（:79-84） |
| `risk.confirm_timeout_sec` | `RiskSection.ConfirmTimeoutSec` | R | cmd/wisp/run.go:616 |
| `risk.l1_window_sec` | `RiskSection.L1WindowSec` | R | cmd/wisp/run.go:615 |
| `risk.permission_mode` | `RiskSection.PermissionMode` | R | cmd/wisp/approval_always.go:155 + internal/perm/store.go:159 |
| `risk.shell_allowlist` | `RiskSection.ShellAllowlist` | G | unwiredKeys 第 3 行（:73-78）；internal/risk/rules_shell.go:41 读的是 risk.Facts.ShellAllowlist，无生产者 |
| `risk.shell_enabled` | `RiskSection.ShellEnabled` | G | unwiredKeys 第 1 行（unwired.go:61-66）⇒ 写非默认值直接加载失败 |
| `session.conversation_idle_sec` | `SessionSection.ConversationIdleSec` | D | 同 session.warm_timeout_sec |
| `session.settling_sec` | `SessionSection.SettlingSec` | D | 同 session.warm_timeout_sec |
| `session.warm_timeout_sec` | `SessionSection.WarmTimeoutSec` | D | 零读者；会话三枚超时全在文档上 |
| `voice.aec.echo_ref` | `AECConfig.EchoRef` | D | 本包外零命中（尺见 §1 头部三条反查） |
| `voice.aec.enabled` | `AECConfig.Enabled` | D | 本包外零命中（尺见 §1 头部三条反查） |
| `voice.asr.model` | `ASRConfig.Model` | D | 本包外零命中（尺见 §1 头部三条反查） |
| `voice.asr.provider` | `ASRConfig.Provider` | D | 本包外零命中（尺见 §1 头部三条反查） |
| `voice.cloud_asr_chain` | `VoiceSection.CloudASRChain` | D | 本包外零命中（尺见 §1 头部三条反查） |
| `voice.cloud_tts_chain` | `VoiceSection.CloudTTSChain` | D | 本包外零命中（尺见 §1 头部三条反查） |
| `voice.conversation_mode` | `VoiceSection.ConversationMode` | D | 本包外零命中（尺见 §1 头部三条反查） |
| `voice.enabled` | `VoiceSection.Enabled` | D | 整个 [voice] 21 枚：段级令牌 .Voice 在本包外零命中（§6 第 4 条） |
| `voice.punctuation` | `VoiceSection.Punctuation` | D | 本包外零命中（尺见 §1 头部三条反查） |
| `voice.realtime.api_key_ref` | `RealtimeConfig.APIKeyRef` | D | 本包外零命中（尺见 §1 头部三条反查） |
| `voice.realtime.base_url` | `RealtimeConfig.BaseURL` | D | 本包外零命中（尺见 §1 头部三条反查） |
| `voice.realtime.enabled` | `RealtimeConfig.Enabled` | D | 本包外零命中（尺见 §1 头部三条反查） |
| `voice.realtime.model` | `RealtimeConfig.Model` | D | 本包外零命中（尺见 §1 头部三条反查） |
| `voice.realtime.provider` | `RealtimeConfig.Provider` | D | 本包外零命中（尺见 §1 头部三条反查） |
| `voice.tts.provider` | `TTSConfig.Provider` | D | 本包外零命中（尺见 §1 头部三条反查） |
| `voice.tts.speed` | `TTSConfig.Speed` | D | 本包外零命中（尺见 §1 头部三条反查） |
| `voice.tts.voice` | `TTSConfig.Voice` | D | 本包外零命中（尺见 §1 头部三条反查） |
| `voice.wake_word.enabled` | `WakeWord.Enabled` | D | 本包外零命中（尺见 §1 头部三条反查） |
| `voice.wake_word.keywords` | `WakeWord.Keywords` | D | 本包外零命中（尺见 §1 头部三条反查） |
| `voice.wake_word.thresholds` | `WakeWord.Thresholds` | D | 本包外零命中（尺见 §1 头部三条反查） |
| `voice.wake_word.veto_words` | `WakeWord.VetoWords` | D | 本包外零命中（尺见 §1 头部三条反查） |

## §2 `unwiredKeys` 名册覆盖率

### 2.1 名册今天登记了哪几枚

`unwiredKeys`（尺：`internal/config/unwired.go`，**现量 var 块在 :60-97**，派单写的 `:57-101` 已漂）
登记 **6 枚**，逐枚：

| # | path | 出处行 |
|---|---|---|
| 1 | `risk.shell_enabled` | `internal/config/unwired.go:61-66` |
| 2 | `risk.allow_shell_string` | `:66-72` |
| 3 | `risk.shell_allowlist` | `:73-78` |
| 4 | `risk.blacklist_overrides` | `:79-84` |
| 5 | `net.allowlist` | `:85-90` |
| 6 | `net.block_private_ranges` | `:91-96` |

全部落在 `risk`/`net` 两段，且都是"无人读"的键（本表 §1 判 G 的那 6 枚）。

还有第二本名册：`lockedKeyDisposition`（`internal/config/unwired.go:119-147`）**17 条**，
其中 5 条只登记 not built 而不拦（`plugins.tier2_enabled` 与四条 `plugins.<id>.*`，本表判 L）。

### 2.2 覆盖率（三个分母，三个数）

| 口径 | 分子／分母 | 覆盖率 |
|---|---|---|
| 无人读／无实效的键里，被任一名册登记者 | 11／100 | **11%** |
| 同上，但只算"会拦下来报错"的（`unwiredKeys`） | 6／100 | **6%** |
| 锁定四段（risk＋fs＋net＋plugins）里无人读的键，被登记者 | 11／11 | **100%** |
| **锁定四段之外**（14 段＋动态子表）无人读的键，被登记者 | **0／89** | **0%** |

⇒ 名册在它的射程内是满的；**射程外一枚都没有**。

### 2.3 射程是钉死的，而且比注释说的还窄

- 名册自己写了边界（`internal/config/unwired.go:32-37`）逐字
  "Deliberately NOT here: zero-consumer keys of unlocked sections (ball, panel, cost, voice, memory, ...)"。
- 完整性钉 `TestEveryLockedSectionKeyIsAccountedFor`（`internal/config/unwired_test.go:311`，
  **本程读到函数体为止，不是引注释**）对四段各跑一次 `collect()`：`risk`/`fs`/`net`/`plugins`。
- 我照它的 `collect()` 逻辑复算了一遍它实际收集到的路径：**16 枚**
  （risk 7 ＋ fs 3 ＋ net 5（含 `net.proxy` 表头）＋ plugins **只 1 枚**）。
  ⚠ **比票面 09-28 记录的射程又窄一枚都不止**：`PluginsSection.Entries` **没有 `toml` 标签**，
  而 `collect()` 第一句就是 `if tag == "-" || tag == "" { continue }`
  ⇒ 那五条 `plugins.<id>.*` 名册行**永远不会被这枚钉查到**（登记了，但钉不吃）。
- 钉只问"这枚锁定键在 `lockedKeyDisposition` 里有没有一句话"，**不问那句话是真是假**：
  它把 `unwired:` 前缀的行去对 `unwiredKeys`（这一步是真判据），
  但 `consumed: ...` 分支**不核实那个消费者是否在场**。⇒ 名册里写 "consumed" 而实际无人读，
  这枚钉不会红。本程逐枚读了 6 条 `consumed:` 的说明，其中 `risk.confirm_timeout_sec`／
  `risk.l1_window_sec`／`net.proxy.mode`／`net.proxy.url`／`fs.*`／`risk.permission_mode`
  我都能在码上对上具名行（见 §1 表），**这一族今天没抓到假话**。

### 2.4 本问最值钱那一格的答

**名册里安静＝不等于全仓没有哑键，而且差的不是几枚，是 89 枚。**
`unwiredKeys` 与 `lockedKeyDisposition` 两本名册合起来只覆盖 4/18 段；
其余 14 段（含 owner 最关心的 `[llm]` 目录里那 12 枚动态模板哑键与 `[panel]`/`[voice]`/`[memory]`/`[cost]`/
`[observe]` 整段）**没有任何一处登记过"写了不生效"**。

## §3 "生效级别"这一维在码上到底存不存在

**结论：作为"数据"不存在，作为"硬编码分支"存在，而且三档里有一档在出货进程里没有听众。**

### 3.1 三档在码上的真实现场（逐档点名）

| 档 | 在场形状 | 现场 file:line |
|---|---|---|
| 类型与常量 | `type Tier string` ＋ `TierHot`/`TierReload`/`TierRestart` 三枚常量 | `internal/config/schema.go:30-44` |
| 逐段分发（真行为） | `Manager.plan()`：锁定段走 `planLocked`，`[app]` 走 `planApp`，`[voice]` 走 `planVoice`，其余 12 段整块 hot | `internal/config/manager.go:243-297`、`:304-334`、`:338-357`、`:361-417` |
| restart 那一支 | `planApp` 比 `Language`/`Autostart`/`SingleInstance` 三枚，改了**不写进 cur**、只 `rep.Restart = append(rep.Restart, "app")` | `internal/config/manager.go:344-356` |
| reload 那一支 | `planVoice` 比模型管线键，置 `reload` ⇒ `rep.Reload = append(..., "voice")` | `internal/config/manager.go:369-379`、`:411-413` |
| 裁决的载体 | `Report{Hot, Reload, Restart []string; Locked []LockedDecision}` | `internal/config/manager.go:87-97` |
| 台词出口（restart） | `OnRestartPending` 挂在 `rt.reportRestartPending`，说出段名＋三枚键名＋为什么不生效 | `cmd/wisp/config_reload.go:115`、`:278-294`、`:299-301` |
| 台词出口（hot） | `rep.Hot` 逐段打进 stdout | `cmd/wisp/config_reload.go:170-172` |

### 3.2 三处"只写在文档里"的形状（具名）

1. **`config.Tier` 那枚类型零读者。**
   尺：`grep -rn "Tier\b|TierHot|TierReload|TierRestart" --include=*.go cmd internal tools | grep -v _test.go`
   ⇒ `internal/config` 之外命中的 `Tier` **全是别人的**：`panel.EffectiveTier`
   （`internal/panel/config_handlers.go:196-203`）、`risk.Tier`（`internal/risk/assessor.go:147`）、
   `projctx.Tier`（`internal/projctx/projctx.go:89`）、`PluginEntry` 那几处 Tier-1/Tier-2 字样。
   `schema.go:30-44` 那三枚常量**没有任何一处 `config.TierHot` 使用点** ⇒ 词汇表在场、码上没牙。
2. **没有逐键的级别数据。** `schema.go` 里出现过的 struct 标签只有 `toml:` 与 `default:` 两种
   （尺：`grep -oE '`[^`]*`' internal/config/schema.go` 取标签键名去重 ⇒ `default:` 69 枚，无 `tier:`）；
   ⇒ 级别不在字段上，只在 `plan()` 的分支里，且粒度是**段**不是键（唯一到键粒度的是 `planApp`/`planVoice` 那两处手比）。
3. **`Report` 载的是段名，不是键名**（除 `Locked.Keys`）。⇒ 一句"panel 立即生效"能覆盖
   `[panel]` 里五枚全都无人读的实话（见 §5）。

### 3.3 reload 档：引擎在、事件在、听众不在

`Manager.CheckAndReload` 在 `len(rep.Reload) > 0 && m.OnReload != nil` 时才发事件
（`internal/config/manager.go:198-200`）。
`OnReload` 的**非测试赋值点全仓一枚**：`cmd/balldebug/main.go:244`（旁路调试程序）。
`cmd/wisp/config_reload.go` 只挂了 `ConfirmLocked`（`:114`）与 `OnRestartPending`（`:115`）。
⇒ 出货进程里 `OnReload == nil`：**换 ASR/TTS 模型不会触发模型加载或卸载**，
而 SPEC-03:134 逐字要求"reload 段触发重载事件"。这一格归 §7 甲#3（是分工未到期还是漏，本程不裁）。

### 3.4 第二套词汇：面板侧的 `EffectiveTier` 是四档，不是三档

`internal/panel/config_handlers.go:196-203`：`now`／`next_task`／`restart`／`not_applied`，
台词在 `tierSentence`（`:428-441`）。
它和 D36 的 hot/reload/restart **不同名也不同数**（`now`↔hot，`next_task` 无对应，
`reload` 无对应，多出 `not_applied`），且只走**设置写入的回执**那一条路。
⇒ "生效级别"这一维今天在码上是**两套并行词汇**，一套死词汇（`config.Tier`）、
一套只服务写入回执（`EffectiveTier`），没有一套是逐键的数据。

## §4 "此项今天不生效"有没有出口 ＋ 最近可扩展点

### 4.1 四条候选出口逐条量

| 候选出口 | 能不能承载"我读到了这枚键，但它没人用" | 凭据 |
|---|---|---|
| **加载期硬错误** | **能，但只对 6 枚**：`validateUnwired` 逐枚比 `unwiredKeys.fires()`，命中就回 "config.toml: X is written but does nothing: ... " | `internal/config/unwired.go:100-111` |
| **重启档台词** | 能，但只对 3 枚（`app.language`/`app.autostart`/`app.single_instance`），且说的是"要重启"而不是"没人读" | `cmd/wisp/config_reload.go:278-294`、`:299-301` |
| **设置写回执** | **能，且已经是现成的档**：`SettingWriteResult.Tier` ＝ `EffectiveNotApplied` ⇒ `tierSentence` 回"这一项没有被应用。" | `internal/panel/config_handlers.go:196-203`、`:210-217`、`:428-441` |
| **面板快照** | **不能**：`Snapshot` 六个 JSON 键（`pending`/`results`/`composer`/`generatedAt` 在 `:58-61`，`instructions` 在 `:74`，`tasks` 在 `:91`）里没有任何一栏承载配置键状态 | `internal/panel/composer.go:57-92`（键集尺：`awk '/^type Snapshot struct/,/^}/' internal/panel/composer.go` 数 `json:"` ⇒ 6） |
| **`config.get` 的返回** | **不能**：`config.get` → `renderSettingsView` 只说可读性／服务商数／凭据状态／聊天模型，**没有逐键状态** | `internal/panel/config_handlers.go:292-303`（`read`）、`:445-462`（渲染） |

### 4.2 判语

**读侧盘上没有这句话；写侧有一模一样的档，但只有 7 枚字段能走到它面前。**

- 写侧那档是**现成的**：`EffectiveNotApplied` ＋ "这一项没有被应用。"
  但它只在 `config.set` 被**受理**之后才存在。`[panel] width` 这枚名字
  过不了受理关——白名单只有 7 枚字段（`internal/panel/config_handlers.go:57-69`），
  且在到白名单之前还有一道拒绝（`:315-326`：`unlisted-field`）。
  ⇒ 用户手改 `config.toml` 里的 `width`，**连那句话都走不到**（那条路只服务面板发起的写）。
- 手改 `config.toml` 这条路唯一会说的话，是 §5 里那句**反方向**的话："这些段已立即生效"。

### 4.3 最近的可扩展点（具名到字段；⛔ 本程不选、不裁）

按"改动半径从小到大"排，每一处都点名到字段：

1. `internal/config/unwired.go:119-147` `lockedKeyDisposition` ＋ `:60-97` `unwiredKeys`
   ——登记面**已经存在**，缺的只是把射程从 4 段扩到 18 段（纯数据行）。
   ⚠ 但扩了之后 `validateUnwired`（`:100-111`）会把它们变成**加载失败**，
   那是"拒绝写入"而不是"告知不生效"，形状与硬编码只读三枚一样重（`audio.half_duplex` 那一族）。
2. `internal/panel/config_handlers.go:142-156` `SettingsView` —— `config.get` 的返回类型。
   它已经有"状态而不是空值"的先例（`Readable` 字段、`CredentialUnknown` 那一族五档）。
3. 同文件 `:445-462` `renderSettingsView` —— 句子成形处，`:447-450` 已经把"配置不可读"当一状态说出来，
   是"此项没人读"这句话最自然的近邻落点。
4. 同文件 `:210-217` `SettingWriteResult.Tier` —— 写侧的 `not_applied` 已在，
   把它引到读侧需要 2 的新栏。
5. `internal/panel/composer.go:57-92` `Snapshot` —— **最不该先动**：那六枚键是有钉的名册，
   `:40-56` 的注释逐字记着"**FOUR, not one**……加一枚要同批改 `panel.ts` 与两枚字节钉"，
   且 Q-51（谁能写 `panel.ts`）仍开。
   ⚠ 顺带一枚过期读数：那段注释今天仍写着"四枚"，而 struct 里已经是**六枚** `json:` 标签
   （`:58-61` 四枚 ＋ `:74` `instructions`（票 200）＋ `:91` `tasks`（票 197））——
   注释在钉数量上已经过期，这正是"注释里的计数一律当待验断言"的又一例。

⛔ 以上第 2／3／5 条都落在 **C17／契约面**（快照键集与方法白名单）。
**本程不裁该不该加、加哪一栏**——按 AGENTS §0 第 2 句与 §2「未定义即停」，
这一格交编排者落批准记录。

## §5 `[panel] width` 那一枚具体走到哪一步断的

### 5.1 逐环点名（HEAD `00e7efe`）

| 环 | 现场 | 值在不在 |
|---|---|---|
| 1 文件 | `config.toml` 的 `[panel] width` | 在 |
| 2 解码 | `internal/config/parse.go` ＋ `loader.go` → `PanelSection`（`internal/config/schema.go:529-539`，`Width` 在 **`:532`** `toml:"width" default:"640"`） | 在 |
| 3 默认 | `internal/config/defaults.go:60` 反射把 `default:` 标签**写进**字段 | 在 |
| 4 内存 | `Manager.cur`（`internal/config/manager.go:106`），`Config()` 深拷贝交出去（`:115-119`） | 在 |
| 5 热加载 | `plan()` 的 `rest` 表里 `{"panel", &cur.Panel, &fresh.Panel, ...}`（`internal/config/manager.go:285`）⇒ 变了就整块拷贝并 `rep.Hot = append(rep.Hot, "panel")`（`:290-295`） | 在，**而且会说话** |
| 6 台词 | `reportReload` 把 `rep.Hot` 打进 stdout："wisp run: 配置热加载：这些段已立即生效（D36 立即档）：[panel]"（`cmd/wisp/config_reload.go:170-172`） | **它说生效了** |
| 7 消费者 | **零枚**。尺：`grep -rn "\.Panel\b" --include=*.go cmd internal tools \| grep -v _test.go` ⇒ 命中全是 `HotkeySection.Panel`（`internal/ball/hotkey_windows.go:98,455`）与审批面板方法 `g.Panel()`（`internal/agent/approval/replies.go:478,487`）；`config.PanelSection` 在本包外**一名读者都没有** | **断在这里** |
| 8 窗口实际尺寸 | `cmd/wisp/panel_host_windows.go:298-307`：`webview2.NewWithOptions(... WindowOptions{Title: panelTitle, Width: 420, Height: 260})`——两枚**字面量** | 顶替了配置 |
| 9 宿主能不能拿到配置 | 不能：`PanelManager`（`:137-167`）没有 config 字段；构造函数 `NewPanelManager(disp, assets, dataPath)`（`:175-177`）不收配置；生产调用点 `cmd/wisp/panel_resident_windows.go:204` 也只传这三样 | 缺整根管道 |
| 10 事后重设 | 无：`SetBounds\|MoveWindow\|SetWindowPos\|Resize(` 在 `cmd/wisp internal/panel` 非测试码里只命中 `panel_host_windows.go:45` 的**注释** | 无第二条路 |

### 5.2 断在哪一环（一句话）

**断在第 7 环：面板宿主从来拿不到配置对象，尺寸在建窗那一行（第 8 环）是写死的 420×260 字面量——
`Width` 的值走完了"文件→解码→内存→热加载→台词"整条路，最后一步被两枚字面量顶替。**

### 5.3 比票面更糟的一格（本程新读数，票面与 09-28 两程都没写过）

第 5→6 环不是静默：`manager.go:285` 那行会让 stdout 逐字说
**"这些段已立即生效（D36 立即档）：[panel]"**。
⇒ 手改 `[panel] width` 的用户得到的不是"没有回应"，而是**一句主动的误报**。
这与票面 §为什么值得做 那句"他没有任何一处能读到'这项今天不生效'"相比更坏一档：
那里还有一句说它生效了。同一句误报覆盖 `[panel]` **全部五枚**（§1 表里 `panel.*` 判 D 的那五行）。

### 5.4 票面 AC#1 那一问今天的答复变了

票面 09-28 两程的结论是"今天没有人决定面板窗口尺寸，因为这棵树里根本没有面板窗口"。
**该读数已过期**：`go.mod:19` 已有 `github.com/jchv/go-webview2`，
`cmd/wisp/panel_host_windows.go` 与 `cmd/wisp/panel_resident_windows.go` 在场且真的建窗。
⇒ "谁决定尺寸"的答案从**没有人**变成**`panel_host_windows.go:304-305` 那两枚字面量**；
⇒ AC#2 的"缺整条链"今天缩成两件事：把配置递进 `NewPanelManager`（第 9 环缺参数），
以及让 `:304-305` 读它。编排者 09-28 那句排程结论（AC#2 排在票 33 之后）
的前置条件**已经满足**，这一格值得重排。

## §6 我可能判错的条目

> 每条给"为什么会错"与"错了的后果"。这节先于表写满，是因为表填完把这节留空＝不算交件。

1. **"叶子键＝150 枚"是我定的口径，不是盘上唯一的口径。**
   尺：`internal/config/schema.go` 的 struct 走查（递归展开 `toml` 标签，`toml:"-"` 除外），
   动态子表（`llm.providers.<id>`、其下 `models.<id>`、`plugins.Entries.<id>`）的键**按模板算一枚**，不按实例算。
   ⇒ 静态 115 ＋ 动态模板 35 ＝ 150。
   **如果错了**：换口径（例如"按一份真实 `config.toml` 里出现的行数"）分母会变，
   于是 §2 的覆盖率百分比会变；**但 D 档名单的成员一枚不变**——判定用的是逐枚读者检索，不依赖分母。
   最可能的分歧是 `schema_version`：我把它算作一枚叶子键（根表里的裸键，不属任何 section），
   有人可能不算。它判 R（`loader.go`/`migrate.go` 消费），挪出去不改任何结论。

2. **"18 个 section"成立，但我用了两把独立的尺才敢这么说。**
   甲：`Config` struct 的 section 字段枚数＝18（`schema.go:110-133`，`Plugins` 虽然带 `toml:"-"` 但由
   `parse.go` 的尾块吸收器落盘，文件里确实有 `[plugins]`）。
   乙：票 198 的首建真件 `.scratch/wisp/probes/198/r1/firstrun-defaults.dump.txt`
   有 33 个表头 ＝ 18 个顶层段 ＋ 15 个子表。两把尺独立给出 18。
   **如果错了**：例如有人把 `position`/`roles`/`wake_word` 这类子表也叫"section"，那 18 会变 33；
   后果只是叙述口径，不影响逐枚判定。

3. **"零读者"的尺是选择器链正则，不是数据流分析。**
   覆盖面我补了三把反查才敢下判：
   ① 反射按名取值的口子：`FieldByName|NumField()|reflect.TypeOf` 在 `cmd internal tools` 的非测试码里
   **零命中**（只有 `internal/config` 自己用反射写默认值）；
   ② 按点分路径字符串寻址：全路径字面量在 `internal/config` 之外只有
   `cmd/wisp/config_reload.go:300`（重启档三枚键名）与两枚同名巧合
   （`internal/statemachine/events.go:21` 的 `"hotkey.mute"` 是事件名、`internal/memory/schema.go:149`
   的 `"schema_version"` 是 SQLite meta 键）；
   ③ 逐枚键名的非测试提及枚数（`cmd internal tools`，排除 `internal/config`）。
   **残余风险**：值经**接口方法**返回后再无人接手（例：`Loop.Budgets()` 返回派生量而非原始键值）；
   以及我按字段名匹配、**同名异主**的命中被我判成"非读者"时判错方向的风险——
   这类我逐枚读过命中行（见 §1 表里具名的 file:line），但只读了命中行本身，没读调用链全文。
   **如果错了**：某一枚我判 D 的键其实有人读 ⇒ 那一格从"哑键"降为"已接"，§2 的哑键总数减一枚。

4. **`[voice]` 21 枚全判 D，是本报告最重的一格，也是我最怕判错的一格。**
   尺：`config.Voice|cfg.Voice|VoiceSection` 在 `internal/config` 之外的非测试码里**零命中**；
   段级令牌 `.Voice` 同样零命中。
   但 `internal/audio/gate.go:56` 的注释逐字写着 `(config AudioSection.MicMutedDefault) here at boot wiring`
   ——**注释指名了一条本该存在的接线**，而按本仓纪律注释指名的东西一律当待验断言。
   **如果错了**（语音链路其实经另一条形如 `audio.Boot(wav, muted)` 的通道取值）：
   `[voice]`/`[audio]` 共 25 枚从 D 改判，本报告的哑键总数从 72 掉到 47 量级，
   而"owner 能配的模型参数里六成套了锁"这句话会弱化成"只有语音那一族没接"。
   ⇒ 这一格的终判该由**语音现量腿（票 228/241 系）**复核，不是这条只读腿能钉死的。

5. **`[hotkey]` 4 枚我判 S（仅旁路调试程序读者），不是 R 也不是 D。**
   依据：唯一生产读者是 `cmd/balldebug/main.go:238` 的 `mgr.Config().Hotkey`；
   出货路径 `cmd/wisp/resident_ball_windows.go:171` 逐字写 `Hotkeys: ball.DefaultHotkeys()`，
   而 `internal/ball/hotkey_windows.go:69` 逐字返回 `"Ctrl+Alt+Q"/"Ctrl+Alt+M"/"Esc"/"Ctrl+Alt+P"` 四枚字面量。
   `mgr.OnReload = bridge.OnReload()` 也**只在 balldebug:244** 出现。
   **如果错了**（balldebug 被算作一个真宿主）：这 4 枚改判 R，但 `wisp run` 里改 `[hotkey]` 仍然没反应——
   结论"改了就无效"不变，变的只是"有没有人读到过它"。

6. **`price.{in,out,cached,audio_in,audio_out}` 5 枚我单独判 W（类型在场·生产零赋值），是最可能判错的一档。**
   依据：`internal/agent/loop.go:130` 有 `Price config.Price` 字段、`internal/agent/cost.go:38` 收它算账，
   但**给它赋值的非测试语句零枚**（尺：`Price:` 全根零命中 rc=1；生产里 `agent.Config{` 字面量只有
   `cmd/wisp/run.go:1010` 一枚且不写 Price）。⇒ 成本账今天按零价卡算。
   **如果错了**（另有装配点把价卡塞进去而我没认出来）：
   票 248 设置页那两枚 `model_price_in/out` 的写入就是**端到端通**的，我这条"写了没人读"的指控当场作废，
   而 owner「自己配模型」这条路上最疼的一格（价卡）其实已经接上。**这一格务必让落地腿复跑。**

7. **`[panel] width` 断点判在"建窗那一行"，尺是反向的：我先证明没有别处能改尺寸。**
   `SetBounds|MoveWindow|SetWindowPos|Resize(` 在 `cmd/wisp internal/panel` 非测试码里
   只命中 `panel_host_windows.go:45` 的**注释**（讲 `(*edge.Chromium).Resize()`，不是调用）。
   **如果错了**：存在一条我找不到的 DPI/scale 重设路径，则断点后移到"那条路径读不读 `cfg.Panel.Scale`"，
   而 §5 那句"断在建窗"要改成"断在配置没递进构造函数"——两句话的修法不同（前者改一行、后者要加参数）。

8. **行号时效。** 本报告所有行号取自起手 HEAD `00e7efe`。
   票面引的 `schema.go:527-528`、派单引的 `unwired.go:57-101` 都已漂（现量见 §8）。
   **如果错了**（我读文件与落笔之间树又前进）：§0 的终态锚与 §8 的差集节就是给人抓 this 的；
   本程没有第二把终态尺之前，任何"某行存在某字段"的引用都应带 HEAD 号读。

---

## §7 判不动的地方

> 甲＝别人补得上（给命令与期望读数）；乙＝补不上，明写不做。

**甲#1 —— "这些哑键里，哪几枚是规格要求今天就必须生效的"**
我只量了码上有没有读者，没裁"该不该有"。规格文字逐段核对（D36 表＋SPEC-03 §3 表＋SPEC-08）需要另一程。
命令：`sed -n '2726,2748p' docs/PLAN.md` 与 `sed -n '24,43p' docs/specs/SPEC-03-config-secrets-envs.md`
逐段对照本表 §1 的 D 档名单，逐枚标"规格要求生效／规格未要求／规格自己写漏"。
期望读数：`[panel]`/`[ball]`/`[hotkey]`/`[session]`/`[voice]`/`[audio]`/`[memory]`/`[privacy]`/`[cost]`/
`[observe]`/`[agent]` 三段里除 §1 已判 R 的三枚 `[agent]` 键外**全部标着 `hot`**
⇒ 若照规格字面读，D 档每一枚都是"规格要求生效但码上没接"，本票的射程会从"一枚 width"扩到几十枚。
**这一格不该由只读腿裁**（它决定要不要把票 180 拆成一族票）。

**甲#2 —— "语音/音频那 25 枚到底是不是真没人读"**
见 §6 第 4 条。命令与期望读数：
`grep -rn "MicMutedDefault\|InputDevice\|SampleRate\|HalfDuplex\|VetoWords\|Thresholds" --include=*.go cmd internal tools | grep -v _test.go | grep -v '^internal/config/'`
⇒ 本程读数：**只剩 `internal/audio/gate.go:56` 一行注释**（非代码）。
语音腿要答的是"这条注释指的 boot wiring 在哪个文件、接线票号是谁"。

**甲#3 —— "`OnReload` 在 `wisp run` 里没挂，是不是票 42（watchdog）的既定分工"**
现量：`OnReload` 的非测试赋值全仓**两枚**——`cmd/balldebug/main.go:244` 与
`internal/config/manager.go:198` 的**读取**；`cmd/wisp/config_reload.go` 只挂了 `ConfirmLocked`（:114）
和 `OnRestartPending`（:115）。⇒ reload 档在出货进程里"有引擎、没听众"。
命令：读票 42 的票面与 `internal/watchdog/doc.go`，判"reload 事件没人接"是**分工未到期**还是**漏**。
期望读数：`internal/config/manager.go:198` 那三行若在生产跑，`OnReload` 必为非 nil；
`cmd/wisp` 里零枚赋值 ⇒ 它就是 nil ⇒ 换 ASR 模型不会触发模型加载/卸载。

**甲#4 —— "面板快照该不该加一栏承载'此项不生效'"**
见 §4。这属 C17／契约面（快照键集是有钉的名册），**只读腿不裁、不许自己决定加不加**。
编排者裁；裁完若要验"加了有没有人读"，参照
`internal/panel/composer.go:40-60` 注释里那条"新段在 pump 装配根没有读者就还是常量"的判据形状。

**乙#1 —— 不做：把 92 枚无消费者键逐枚写出"它该由哪个票号接"**
票面 AC#4 只要"出名单并逐枚归口"。归口要读 `.scratch/wisp/issues/**` 全池与 SPEC 切片表，
射程远超一条只读腿，且 `issues/83-...md:155` 那张归口表（票面 09-28 记录称 18 枚全在其中）
我**没有复认**——引用它之前得自己读到那行为止，而这一程我没读到。
⇒ 本程只交"名册＋覆盖率＋断点"，归口一律标为未做，不写"以后加固"。

**乙#2 —— 不做：跑任何编译/测试类尺**
`go test ./internal/config/`、`go list -deps` 全都没跑（`198-v1` 正在整包跑 `cmd/wisp`）。
⇒ 本报告的枚数全部出自文本尺与 Python 走查，**没有一把是编译器给的**。
若某枚判定最终依赖"这个字段确实被链接进来"，本程给不出这个级别的凭据。

**乙#3 —— 不做：`frontend/**` 与 `design/**` 两层的任何核对**
禁令覆盖。⇒ 若面板侧另有配置编辑界面（哪怕它已经渲染 `panel.width` 的输入框），
本程看不见、也不引它行号；"此项不生效没有出口"这句的射程因此**只到 Go 侧为止**。

## §8 我推翻票面与派单哪一句

> 起手 HEAD `00e7efe`。每条都给现量尺。

1. **派单："18 个 section"——成立，不是我猜的。**
   两把独立尺（`Config` struct 段字段枚数＝18；首建真件 33 个表头＝18 顶层＋15 子表）。
   另加一句修正：严格说是 **18 段 ＋ 1 枚根表裸键 `schema_version`**，
   有人按 19 数会被我这句绊到。
2. **派单：`internal/config/unwired.go:57-101`——行号已漂。**
   现量：`unwiredKeys` var 块在 **`:60-97`**；第二本名册 `lockedKeyDisposition` 在 **`:119-147`**。
3. **票面现量第 1 条：`internal/config/schema.go:527-528`——行号已漂。**
   现量：`PanelSection` 在 **`:529-539`**，`Width` 那句在 **`:532`**。字段与 `default:"640"` 都在，
   只是不在那一行。
4. **票面现量第 2 条的尺（`grep "\.Width"` 空输出）——仍然成立。** 本程独立复跑，仍空。
5. **票面 09-28 两程的"这棵树今天没有面板窗口"——已过期。** 见 §5.4。
   这条最重：它直接决定 AC#2 的排程前置（票 33 宿主）是不是还在等人。
6. **票面 AC#4"还有几枚"——现量答案比 09-28 的 18 枚大得多：名册外哑键 76 枚。**
   原因：`180-c1` 那一遍只数了**带 `default:` 标签**的字段（69 行 / 61 个去重字段名），
   本程按叶子键全量数（150 枚）。两个数不矛盾，分母不同。
7. **票面 09-28 记录"这 18 枚全部已在票 83 的 AC#1 全量表里登记并归口"——我未复认。**
   按派单硬边界，`issues/83-...md:155` 那一行本程没读到，**当待验断言**引用（见 §7 乙#1）。
8. **派单问 3 的预设（"三档在码上有没有真实现场"）——答案是分裂的，比"存在/不存在"都更难看。**
   级别**分发**在场（`plan()`／`planApp`／`planVoice`），但 `config.Tier` 那枚类型
   与三枚常量**零读者**，且没有任何逐键的级别数据（无 `tier:` 标签，见 §3.2）。
   另加一句派单没写的：面板侧另有一套**四档**词汇 `EffectiveTier`，与 D36 三档不同名不同数。
9. **派单问 4 的预设（"没有就明写盘上没有"）——部分推翻。**
   写侧**已经有**这一档（`EffectiveNotApplied` ⇒ "这一项没有被应用。"，票 248 落的），
   缺的是**读侧与快照侧**。照"明写盘上没有"交出去会把已有的一块现成设施说没。
10. **终态复认**：本程落笔期间未主动跑 `git log -1` 第二把尺的每一枚行号。
    终态读数与差集见本节末（下一次 commit 回执里逐条给）。
    若终态 HEAD ≠ `00e7efe`，本文件所有行号按"取自如未复认的 `00e7efe`"读。
