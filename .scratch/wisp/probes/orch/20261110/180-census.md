# 180-census — AC#4 只读普查（2026-11-10，零产码、零 go 命令、AC 框一枚未碰）

- 腿＝票 180 **AC#4**（`schema.go` 里还有几枚「有 `default` 但生产零读者」的字段，逐枚 grep 现量＋逐枚归口）。
- 起手 HEAD `a2495ae15`（A1686 wave119 收口那枚）；全程只跑文本尺，未跑任何编译/测试类命令。
- 台件（只建不删，与本表同目录）：`ruler.sh`（母体抽取＋尺 A/B）、`ruler2.sh`（样本行＋段内立尺 D）、`ruler3.sh`（反控两读数＋正控）。
- 权威源：票 83 AC#1 全量表（`issues/83-…done.md:126-162`，09-21 归口登记）＋ 票 255 `cmd/wisp/config_readers_255.go` 的 `hotRowClaims`/`sectionReadSites` 名册（已含票 233/313/413/419/65/44/320 的更新）＋ `internal/config/unwired.go`（146 行版）。

## 0. A1667 前提复验：票面与前几程的断言今天漂了多少

| # | 旧断言 | 今日现量 |
|---|---|---|
| 1 | 母体 69 枚带 `default:` 字段（10-08 `180-c1`） | **70 枚**：票 65 给 `BallSection` 新增了 `Style`（`toml:"style" default:"glass"`）。`grep -c 'default:' schema.go`＝71，差的那一枚仍是 `schema.go:8` 文档注释行 |
| 2 | 10-08 尺 A 零命中 17 枚 | **14 枚**：`WarmTimeoutSec`（3）、`SampleRate`（12）、`KeepAliveInSession`（1）已离开零类 |
| 3 | 10-08「真·全树无提及只有 EchoRef／SizeMB／Days 三枚」 | **只剩两枚**：`Days` 今日 raw=9（全测试读者），`EchoRef`／`SizeMB` 仍 raw=0 |
| 4 | 票面 11-09 复验引 `panel_resident_windows.go:207` | 该读点已漂到 **`:231`**（`return cfg.Panel.Width, cfg.Panel.Height`；roster 里的 `:207` 是旧读数——roster 自己的规矩就是 token 唯一命中当判据、行号只当读数） |
| 5 | 10-02 a1「`[voice]` 段 21/21 无人读」 | **已过期**：今日 `voice.enabled`、`voice.asr.provider`、`voice.wake_word.enabled`、`voice.realtime.{enabled,provider}` 都有真生产读者 |
| 6 | 10-08 c1「19 枚归属未定」 | 本程用**段内立尺（尺 D）**＋票 255 名册把其中大半裁掉了：见 §3 判读列 |
| 7 | 10-02 a1「reload 档在出货进程没有听众」 | 名册化后 `Manager.OnReload` 已挂到 reload 回执链（本程未量 caller 枚数，**不追认不推翻**，只记句式过期风险） |

母体外备注（不进 70 枚计数）：`[panel] scale`（`schema.go:576`，无 `default` 标签）在票 255 名册 panel 行里不被提 ⇒ 疑似仍零读者；`net.allowlist`／`risk.shell_allowlist`（无 `default` 标签）由 `unwiredKeys` 守卫拦着。`[panel] height`（无 `default`）有读者（`:231`）。

## 1. 母体与判读总账（70 枚）

判读口径：**「有生产读者」＝internal/config 包外存在一条非测试生产路径按这个值做事**（boot 档／常驻进程档都算；「值已换进内存但无人按新值做事」不算）。射程只到 `internal/`＋`cmd/` 的 Go 侧，`frontend/**`／`design/**` 未读未引。

| 判读 | 枚数 |
|---|---|
| **R 有生产读者** | **33** |
| **G 守卫拦截**（写非默认值＝加载报错，不说谎；先例族＋unwired 族） | **7** |
| **Z 零生产读者（真剩余，AC#4 的交付对象）** | **30** |
| 合计 | **70** |

## 2. ★主产出：30 枚零生产读者逐枚归口表

**归口状态是本程最重要的新发现：30 枚里 25 枚的既有归口票已经 `-done`，缺口没有随那张票走掉（悬空）；只有 5 枚还挂在未结票上。** 归口悬空＝那张票结了、字段的缺口没被它带走、票池里今天没有一张 open 票认领 ⇒ **逐枚新开票/挂台账是编排者的事（本腿 ⛔ 改票）**，本表只把「该谁」与「那张票现在什么状态」写到字段级。

| 字段（toml 键） | 今日读数（raw/prod/test-only） | 读者判读 | 归口 |
|---|---|---|---|
| `app.language` | 14/7/7（prod 命中全是撞车：speech/web_search 的同名结构） | 零 | 📋 **票 39（open，GUI 配置编辑器）**（83 表 :129） |
| `app.theme` | 22/4/18（4 枚 prod 全在 config 包内） | 零 | 📋 **票 39（open）**（83 表 :129） |
| `app.single_instance` | 3/1/2（prod 命中＝manager.go 重启档比较） | 零 | 票 43 **done** ⇒ **悬空** |
| `ball.opacity_idle` | 4/2/2（prod＝validate 枚举校验） | 零 | 📋 **票 62/68（open，球视觉）**（83 表 :131） |
| `ball.click_through` | 3/0/3 | 只有测试读者 | 📋 **票 62/68（open）** |
| `ball.hide_on_fullscreen` | 1/0/1 | 只有测试读者 | 📋 **票 62/68（open）** |
| `session.settling_sec` | 1/0/1 | 只有测试读者 | 票 28 **done** ⇒ **悬空**（票 255 名册 session 行明文 "still read from nowhere"） |
| `session.conversation_idle_sec` | 1/0/1 | 只有测试读者 | 票 28 **done** ⇒ **悬空**（同上） |
| `voice.wake_word.veto_words` | 3/2/1（prod＝manager 段拷贝） | 零 | 票 26 **done** ⇒ **悬空**（票 255 名册 `voice.wake_word.veto_words` 行＝no-reader） |
| `voice.tts.provider` | scoped raw=1（全测试） | 零 | 票 61 **done** ⇒ **悬空**（云链）／票 26 done（本地 TTS） |
| `voice.tts.speed` | prod=6，唯一 config 读点 `cloudcatalog.go:119 [ep.Speed = cfg.Voice.TTS.Speed]` 住在 `ResolveCloudChain` 里——该函数**全仓非测试调用者 0** | 死读点 | 票 61 **done** ⇒ **悬空**（票 255 名册明文归票 61） |
| `voice.aec.enabled` | scoped raw=0 | 零 | 票 27 **done** ⇒ **悬空** |
| `voice.aec.echo_ref` | **raw=0（全树零提及）** | 只有测试读者 | 票 27 **done** ⇒ **悬空** |
| `voice.punctuation` | 4/2/2（prod＝manager 段拷贝） | 零 | 票 27 **done** ⇒ **悬空**（票名就是 punctuation；票 255 名册 `voice.punctuation` 行＝no-reader） |
| `voice.conversation_mode` | prod=6 全是 `statemachine.Facts.ConversationMode` 撞车——状态机那枚由**常量 `false`** 喂（`resident_voice_engines_windows.go:1241`），不是 config | 撞车 | 票 26/27 **done** ⇒ **悬空** |
| `audio.input_device` | 1/0/1 | 只有测试读者 | 票 13 **done** ⇒ **悬空**（票 255 名册 audio 行：internal/audio 走 OS 默认端点） |
| `audio.sample_rate` | 33/12/21（prod 命中全是 audio 引擎自己的 SampleRate 方法/字段撞车） | 撞车 | 票 13 **done** ⇒ **悬空**（同上：走 seam 常量） |
| `agent.loop_guard.repeat_thresholds` | prod=6＝`agent/guard.go:90-91` 读的是 **agent 侧自己的 Config**；组合根**从未把 config 值喂进去**（`grep LoopGuard` 包外 0） | 插槽在、喂线缺 | 票 10 **done** ⇒ **悬空**（83 表 :140 明文「组合根未接（票 10 已 done 的那条线尾巴）」） |
| `agent.max_rounds` | prod=7＝guard.go:101/loop.go:242 读 agent 侧结构；config 未喂 | 插槽在、喂线缺 | 票 10 **done** ⇒ **悬空** |
| `agent.token_budget` | prod=7＝guard.go:104-105 读 agent 侧结构；config 未喂 | 插槽在、喂线缺 | 票 10 **done** ⇒ **悬空** |
| `privacy.redact_paths` | prod=6＝`observe` 自己的 `Redactor.RedactPaths` 镜像字段，**由调用者填、从未从 cfg.Privacy 喂**（`logging.go:115` 的 `cfg` 是 observe.Options，不是 config.Config） | 镜像字段未喂 | 票 45 **done** ⇒ **悬空**（票 255 名册 privacy 行＝no-reader 明文） |
| `privacy.diagnostics_opt_in` | 1/0/1 | 只有测试读者 | 票 45 **done** ⇒ **悬空** |
| `privacy.retention_days` | 1/0/1 | 只有测试读者 | 票 45/票 08 **done** ⇒ **悬空** |
| `memory.l1_enabled` | 1/0/1 | 只有测试读者 | 票 29 **done** ⇒ **悬空**（票 255 名册 memory 行：L1/L3 旋钮无消费者） |
| `memory.l1_max` | 1/0/1 | 只有测试读者 | 票 29 **done** ⇒ **悬空** |
| `memory.l3_retention_days` | 1/0/1 | 只有测试读者 | 票 29 **done** ⇒ **悬空** |
| `roll.size_mb` | **raw=0**（observe 的 `RollSizeMB` 是另一结构，config 未喂） | 全树零提及 | 票 08/42 **done** ⇒ **悬空** |
| `roll.days` | 9/0/9 | 只有测试读者 | 票 08/42 **done** ⇒ **悬空** |
| `observe.level` | 419/107/312——107 枚 prod 命中**全非 config**（`risk.Level`/审批卡撞车）；`observe.level` 的唯一 config 命中是 validate 枚举校验；日志级别由 logsink 自选 | 撞车 | 票 08/66 **done** ⇒ **悬空**（票 255 名册 observe 行明文） |
| `observe.slo_sample_interval_sec` | 1/0/1 | 只有测试读者 | 票 08/66 **done** ⇒ **悬空** |

**归口小结**：open 票还认领的 **5 枚**（language、theme→39；opacity_idle、click_through、hide_on_fullscreen→62/68）；**悬空 25 枚**（owner 已 done：43×1、28×2、26/27/61×6、13×2、10×3、45×3、29×3、08/42×4、66×1、61×1——按字段计 25）。25 枚悬空需要编排者裁「开新票还是挂台账」，本腿不裁（⛔ 改票；且 10-08 收口节已记过同类处置排在票 255 之后）。

## 3. G 类：7 枚「零读者但不说谎」（守卫拦截，不在零读者账里）

| 字段 | 拦截点 | 备注 |
|---|---|---|
| `risk.shell_enabled` | `unwiredKeys`（写 true＝加载报错） | capability＝shell.exec 工具（SPEC-07 §3 S3）；83 表 :142 |
| `risk.allow_shell_string` | `unwiredKeys` | 同上（SPEC-06 §10 强制 argv 规则）；83 表 :143 |
| `net.block_private_ranges` | `unwiredKeys`（写 false＝加载报错） | 真做需 D22 契约裁定（SPEC-06 §6.3 钉死）；83 表 :149 |
| `audio.half_duplex` | `validate.go:72`（写 false 报错） | 🛑★ 先例；83 表 :136 |
| `privacy.keep_transcript` | `validate.go:82` | 🛑★ 先例 |
| `privacy.keep_audio` | `validate.go:86` | 🛑★ 先例 |
| `models.verify_signature` | `validate.go:123`（写 false 报错）＋调用方**故意传字面量 true**（`cmd/wisp/models.go:170`，C29 硬编码） | 值本身设计上不消费；83 表 :159 |

（`risk.blacklist_overrides` 已于票 230 AC#4 接走、`unwiredKeys` 行已删——`unwired.go:79-83` 注释逐字在。）

## 4. R 类：33 枚有生产读者（只列名＋读点，作零读者账的对照面）

`app.autostart`（`autostart_windows.go:62,65`，票 43）· `ball.size`（`ball_size_413_windows.go:48`，票 413，boot 档）· `ball.style`（`ball_style_65_windows.go:50`＋`config_cmd_65.go:150`，票 65，boot 档）· `hotkey.cancel`（`resident_windows.go:73` 映射进 `ball.HotkeyConfig`，票 258/64，常驻档）· `session.warm_timeout_sec`（`resident_voice_windows.go:1823`，boot 档）· `voice.enabled`（`resident_audio_windows.go:257`＋`resident_voice_windows.go:1315`）· `voice.wake_word.enabled`（`resident_voice_windows.go:1864`）· `voice.asr.provider`（`resident_voice_windows.go:1815`）· `voice.realtime.enabled`/`provider`（`panel_config_store.go:108-110`，票 313，**live**）· `audio.mic_muted_default`（`resident_audio_windows.go:256`，票 247，boot 档）· `llm.roles.*.temperature`/`thinking_intensity`（`prompt.go:252,254`，喂线＝每任务 `ResolveRole`，`run.go:1842`，票 370）· `llm.timeout_ms`（`run.go:2021`）· `llm.retry.max`/`backoff_ms`（`run.go:2404,2403`）· `compat.loose`/`allow_missing_usage`（三家 adapter）· `providers.<n>.billing`/`models.<id>.billing`（`quota_44r4.go:115,121`，票 44）· `models.<id>.enabled`（`resolver.go:128`，I-G 局部读）· `agent.steering_enabled`（`run.go:1906`→`loop.go:318`，票 233 后 live）· `agent.project_instructions_enabled`（`run.go:1990,1995`）· `risk.confirm_timeout_sec`/`l1_window_sec`（`run.go:989,988`＋`resident_approval_windows.go:557,556`）· `risk.permission_mode`（perm store／bridge，票 90）· `fs.delete_enabled`（`run.go:761`→`fs_write.go:596`）· `net.proxy.mode`（`providers.go:138`＋`panel_models_320.go:110`＋`run.go:2402`）· `panel.enabled`（`panel_resident_windows.go:258`，33-r3 `908c6652c`）· `panel.width`/`height`（`:231`，票 255 AC#4）· `panel.keep_alive_in_session`（`:276`，33-r3）· `cost.alert_threshold`/`over_budget`（`costbook_44.go:184`＋`quota_44r4.go:136`，票 44）· `plugins.tier2_enabled`（`run.go:1209`＋`panel_resident_pump.go:208`→`plugin/tier2_gate.go:154`）。

## 5. 反向两读数

- **反控①（零读者字段的 raw 计数，含测试）**：30 枚里 raw=0（全树连测试都不提）＝ **2 枚**（`echo_ref`、`size_mb`）；「只有测试读者」＝ **11 枚**（click_through、hide_on_fullscreen、settling_sec、conversation_idle_sec、input_device、diagnostics_opt_in、retention_days、l1_enabled、l1_max、l3_retention_days、slo_sample_interval_sec）；其余 17 枚 prod>0 但逐枚判读为撞车／死读点／未喂线（§2 判读列）。**措辞纪律照 10-08 收口节：⛔ 不许把 30 枚写成「全树没人提」。**
- **反控②（正控，段内立尺）**：26 对段内 pattern 现跑全部非 0（`ruler3.sh` 输出逐行：`Panel.Enabled=1`、`Panel.Width=2`、`Ball.Style=9`、`Cost.AlertThreshold=2`、`Plugins.Tier2Enabled=2`…），其中 `Voice.Enabled` 那一对因尺脚本转义多了一个 `\.` 得 0、以尺 D 的 8 枚命中（`resident_audio_windows.go:257`／`resident_voice_windows.go:1315`）为准——**尺的缺陷记在此，不改读数**。
- **同名撞车警示**（沿用 10-08 纪律）：`Enabled`（21→39 枚命中，无一枚属 `PanelSection`/`VoiceSection` 之外的 config 段——段内立尺才把它们分开）、`Level`（419 枚）、`Mode`、`Size`、`Provider`、`Days`（9 枚全是测试）单跑字段名尺必然高报，本表全部以段内立尺或票 255 名册为准。

## 6. 本程没做（明写）

- **零产码、零 go 命令**：全部读数出自文本尺（grep/awk/sed），没有一把是编译器给的；未跑票 180 AC#5/AC#6 的任何格。
- **AC 框一枚未勾**（勾归编排者/落地腿）；票面一字未动；`schema.go` 一字节未动（尺：本程只读）。
- **25 枚悬空字段的处置未裁**（开新票/挂台账/并格＝改票面，本腿 ⛔）；票 83/255 名册里的过期行号（如 roster panel 行的 `:207`）**未修**（在别人的写面）。
- `frontend/**`／`design/**` 未读未引 ⇒ 「零读者」射程只到 Go 侧 `internal/`＋`cmd/`。
- 未核 25 枚悬空字段在 `docs/reports/pending-and-issues.md` 台账里是否已有 A## 账（读台账超本程射程，未做）。
