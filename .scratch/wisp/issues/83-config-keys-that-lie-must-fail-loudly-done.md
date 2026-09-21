# 83 — 配置里"看着能用、其实没接线"的键要**响亮地失败**（票 80 裁决 (C)；先例 `SPEC-03:42 verify_signature`）

**Status:** **done**（编排者对抗验收 2026-09-21 14:3x）—— 五框全 PASS，裁决表
`docs/evidence/s1/83-adversarial-acceptance.md`。我亲自复跑：包内门禁逐字对上
（`go test -count=2 -v ./internal/config/` rc=0，`=== RUN 194 = 2 × 97`、`PASS 106`、`FAIL 0`、`SKIP 0`），
并**自己下了一刀它没下过的变异**（把 `validate.go:30` 的 `validateUnwired(c),` 注释掉 ⇒
`TestUnwiredSecurityKeysFailLoudly` 等 **4 条命名用例同时红**，还原后 `git diff --quiet` 干净）
⇒ 判据真的咬得住，且"诚实配置不受影响""报错早于任何写盘"这两条不是装饰。
五个交回问题的裁决（原文在裁决表末段）：`[net]` 两键**算扩面**、`[plugins]` 五键**登记不扩面**（归票 50）、
两处冻结文件的限定语**我不代落**（只有 owner 能改 ⇒ 已转成 Q-27 的批语）、墙钟脆弱性**建票 86**、
`ban #6` 与 `frontend/` 的豁免翻转**由我在建树同批做**（不许建目录的代理顺手改豁免文件）。
起点 HEAD `4b05be0`，本轮 commit：`6028977`(认领) → `a95ee3a`(守卫码+用例) → `4dec91b`(AC#1 表) →
`e4a0fac`(L3/L4) → `5ca30a7`(**它自己修掉自己带进 CI 的 emoji 哑弹**)

**Type:** 安全可用性/诚实性（一个说谎的配置键）——**不是**新能力
**Blocks:** nothing · **Blocked by:** nothing（`internal/config` 此刻无人写；票 80 已交回且零 Go 改动）
**Packages:** `internal/config/`（校验与加载路径）+ 新建的用例。**禁改**：`docs/PLAN.md`、`docs/specs/*.md`
（冻结面，那两行限定语由编排者落）、`internal/risk/assessor.go`、`internal/risk/pathresolver*.go`、
`rules_gateway.go`、`internal/agent`+`internal/memory` 的测试（票 81）、`internal/risk` 的 sync 家族测试（票 82）。

## 依据：票 80 的穷尽清单（读/写逐条在票面 L2，证据 `docs/evidence/s1/80-ac1-ac2-verdict.md`）

`blacklist_overrides` 的真相**比建票时写的更空**：
- 写：`internal/config/schema.go:457`（TOML 解码）；读：**只有** `manager.go:356` 的方向审计。**没有任何消费者**。
- `risk.Gate` 在生产里**零调用点**。真实判定链是 `cmd/wisp/run.go:262 tools.New` →
  `internal/tools/bridge.go:159 NewSensitiveClassifier()` → `internal/paths.go:146 risk.Classify()`，
  而 `Classify`（`blacklist.go:69`）**签名里没有 override 形参** ⇒ "接线"不是灌一个 map，是要先在判定链上
  **开出能携带豁免集合的位置**。
- 该键唯一残存效果（🔒 变更检测，`manager.go:239`）也是哑的：`ConfirmLocked` 生产从未赋值。

**票 80 的 AC#2 裁决（编排者已签）**：契约**从未承诺**这个键在决策时生效——
`SPEC-06:67` 把豁免定义成**运行时事件**（"豁免 = 一次 L2 强确认 + 写日志"）、`PLAN.md:2396-2399` 明写"**必须由人点**"、
`assessor.go:144-146` 的冻结注释把 exemption flow 推给票 21。⇒ 接静态预授权 = **新造放行侧能力**（选项 A/B，需 D22，
且本来就是票 21 的欠账）。本票走**选项 (C)**：让键**响亮地失败**。

## 为什么 (C) 而不是 (D) 维持现状

现状是"配置项说谎"：用户写下一行他认为在放开某个文件的键，程序默默吃掉它。
先例就在同一份契约里：`SPEC-03:42` 的 `verify_signature` 写 `false` 时报错而非静默生效。
(C) 是四个选项里**唯一收紧侧**的那个，且**不需要**动安全语义、不需要 D22 批准新能力。

## AC（1:1，裁决表 `docs/evidence/s1/83-*.md`）

- [x] **AC#1** **先把同类键扫全**（不要只修一个键就交）：`grep` 出 `internal/config` 里每个被解析、
  但在 `internal/config` 之外**零消费者**的键。票 80 已经点了三个候选：`[risk] shell_enabled`、
  `allow_shell_string`、`shell_allowlist`（`grep "\.Risk\."` 全仓只命中 `L1WindowSec`/`ConfirmTimeoutSec`），
  以及 🔒 段的 `ConfirmLocked`（D36 规则 1，生产未赋值）。
  **交付物 = 一张表**：键 → 解析处 → 消费者（"无"也要给 grep 证据）→ 处置（响亮失败 / 已有计划接（指名票号）/ 保留但文档说明）。
  ⚠ "grep 零命中"在这里**不足以**定罪：要顺带查是否存在**间接消费**（反射、按名取值的 map、TOML 原样透传）。
- [x] **AC#2** 校验落地：对表里判为"响亮失败"的键，**加载时**报错并指名键名 + 一句"该能力尚未实现/由票 N 实现"，
  **错误发生在任何副作用之前**（照 `verify_signature` 先例的形状）。
  判据用例必须包含：**不写这个键的用户不受影响**（默认配置仍加载成功）。
- [x] **AC#3** 双向变异：(i) 把校验去掉 ⇒ 用例红；(ii) 把校验做宽（例如对所有未知键报错）⇒
  必须有用例**因此变红**（防止"用一个大棒假装修好了"）。锚点=承载行为的那一行，同链 grep 自证，还原后 `diff -q`/`git diff --quiet`。
- [x] **AC#4** 明确**不做**：不接 `risk.Gate`、不给 `Classify` 开 override 形参、不动审批队列。
  在票面写一行"这些属于票 21 / 需 D22"，并把它链接到 A51⑤ 与票 80 的裁决。
- [x] **AC#5** 门禁（只跑自己碰的包）：`gofmt -l` 空、`gofumpt -l <files>` 空、`go vet ./internal/config/` rc=0、
  `GOOS=linux go vet ./internal/config/` rc=0、`go test -count=2 ./internal/config/` rc=0，
  逐跑点名 `--- SKIP`/`--- FAIL`。**并且**：如果两包并跑出现墙钟抖动（票 80 量到 `TestResolvePerCallBudget`
  在 `config`+`risk` 并跑时红、单跑绿），**如实登记为外项**，不许顺手调那个 1ms 预算。

## 每条 AC 的落点（裁决表就是票面自己，证据文件 `docs/evidence/s1/83-*.md` 本票**没建**——理由见下）

| AC | 落点 | 达成 / 未达成 |
|---|---|---|
| AC#1 | L2 的键表 + 七条仪器 I-A~I-G | **达成（有声明的边界）**：Go 侧全量；**非 Go 侧（IPC/脚本按字符串取键）未扫**，L2 末段明写 |
| AC#2 | `internal/config/unwired.go`（6 行守卫表）+ `validate.go:30` 挂载 + `unwired_test.go` 六个子例 & "默认配置照常加载" 四个子例 & 副作用次序用例 | **达成** |
| AC#3 | L3（去掉校验 ⇒ 11 条红；做宽成"连有消费者的键一起拒" ⇒ 6 条红，含设计锚点；两次都有 grep 自证与还原证明） | **达成** |
| AC#4 | L4 上半（三件不做的事 + 归票 21 / 需 D22 的链接） | **达成** |
| AC#5 | L4 的表（本票包五项全 rc=0；两包并跑红 2 条 = `TestResolvePerCallBudget`，**登记为外项、未动 1ms**） | **达成（包内）**；⚠ 本轮**单跑也红**，见 L4 的抖动登记 |
| 未闭的框 | **无一个框未闭**，但两处**没做到位的地方要说清**：① 没有单独的 `docs/evidence/s1/83-*.md`（AC 抬头要求"裁决表 `docs/evidence/s1/83-*.md`"）——票面 L2/L3/L4 已逐条 1:1 承载，且本仓规矩是"不主动新建报告 md"，要不要另立一份交编排者定；② `[plugins]` 五键**只登记未守卫**（L2 表 + L4 交回项 ②） | 见"需要编排者拍板" |

## Rules（本仓固定）

15 次工具调用内交第一枚 checkpoint commit；每次 commit 同步 Status + 勾框 + 末行 `next=`；
`git commit -q -F - -- <显式路径>` + **带引号** heredoc；禁 `git add -A`/`.`；commit 前核对 `git diff --cached --name-only`；
禁 `--amend`/`reset`/`rebase`/`stash`/`checkout .`（共树 A34）；**不 push**；不在仓内建 worktree（A38④）；
票面 Progress log append-only，**要改的那行先读再替换**；四种假绿逐跑点名。
**spec 那两行限定语（`SPEC-03:34` 键表加"未实现 ⇒ 加载报错"）你不要写** —— 把拟好的文字放进票面，编排者落。

## Progress log（append-only）

### L0 — 认领 + checkpoint（2026-09-21）

Status `open` ⇒ `claimed`。本轮次序：AC#1 全量同类键扫描（含间接消费排查）⇒ AC#2 校验落地 ⇒ AC#3 双向变异 ⇒ AC#5 门禁。
零 Go 改动，纯 checkpoint。
`next=` 把 AC#1 的键表全量写进本票（含间接消费排查证据），再跑双向变异与并跑门禁。

### L2 — AC#1 全量同类键表（含间接消费排查）（2026-09-21）

**先给仪器，再给表。** 六条检索式全部在本仓 HEAD 上跑过，命中数写在括号里：

- **I-A 主仪器**：`grep -rn --include=*.go -E '\.(App|Ball|Hotkey|Session|Voice|Audio|LLM|Agent|Risk|FS|Net|Privacy|Memory|Panel|Cost|Plugins|Models|Observe)\.' internal cmd tools | grep -v '^internal/config/'`
  ⇒ 去重后**只有 16 条键路径**被 `internal/config` 之外读到（含测试文件）：
  `.Agent.PerToolTimeoutMS .Agent.SteeringEnabled .FS.AllowedDirs .FS.DeleteEnabled .FS.ReparsePointExceptions
  .LLM.Providers .LLM.Retry.BackoffMS .LLM.Retry.Max .LLM.Roles .LLM.Roles.Chat .LLM.TextChain .LLM.TimeoutMS
  .Net.Proxy.Mode .Net.Proxy.URL .Risk.ConfirmTimeoutSec .Risk.L1WindowSec`。
- **I-B 整段被取走**（段名后不接点，防 `cfg.Risk` 整段传下去）：非注释命中只有
  `internal/ball/hotkey_live_test.go:368` 的 `m.Config().Hotkey`（**测试**），以及
  `internal/ball/hotkey_reload.go:15-19` 的**文档示例注释**；生产装配 `cmd/wisp/run.go` 未装这个 reloader。
- **I-C 谁 import 了本包**：`grep -rl '"github.com/CarlosShao/wisp/internal/config"' --include=*.go`
  ⇒ 非测试 8 个文件（`cmd/balldebug/main.go`、`cmd/wisp/{providers,run}.go`、`internal/agent/{cost,loop}.go`、
  `internal/llm/{discover,probe_health,resolver}.go`），**逐个打开核对过**。整只 `*Config` 被传下去的只有
  `llm.NewResolver(cfg,…)` 与 `chainOptions(cfg)` 两处 ⇒ 其字段读取逐行落到下面的 [llm] 表。
- **I-D 间接消费①（反射按名取值）**：`grep -rn "FieldByName\|FieldByIndex" --include=*.go internal cmd` ⇒ **1 处**，
  且它是 `parse.go:77` 对 `reflect.StructOf` 生成的解码镜像做**下标**访问，不是拿 TOML 键名去查字段。
  `defaults.go` / `manager.go` / `parse.go` 的 reflect 只在 **Config 结构**上走默认值/深拷贝/差集，
  **没有任何一处把配置里的字符串当字段名用**。
- **I-E 间接消费②（按名取值的 map / 原样透传）**：`grep -rn "map\[string\]any" internal/config` ⇒
  只有 `migrate.go:104-161`（v1→v2 重写 `[llm]` 子树）与 `parse.go:30/92/181`（`[plugins]` 动态子表吸收器）。
  **两条都把结果重新塞回同一个强类型结构、再走同一个 `validate()`** ⇒ 没有"TOML 原样交给下游"的旁路。
- **I-F 间接消费③（第二个 TOML 解析者）**：`grep -rn "config\.toml" --include=*.go internal cmd tools | grep -v '^internal/config/' | grep -v _test.go`
  ⇒ 23 命中，逐条打开：全是注释与路径常量（`risk/blacklist.go:219-220` 把 `%APPDATA%\wisp\config.toml` 钉成 A 档、
  `observe/diagnostics.go:147` 把它再脱敏导出、`secret/migrate.go:26` 的备份名）。
  **不存在第二处 `os.ReadFile` + 自解析**。
- **I-G 局部变量读**（I-A 的盲区）：`internal/llm/resolver.go:119-134` 用 `p.` / `spec.` 读 Provider/ModelSpec 字段，
  I-A 数不到 ⇒ 已逐行核对并写进 [llm] 行。

**仪器的闭包性（为什么"零命中"在这里够用）**：一个值要被消费，只有三条路 ——
(a) 有人读一次 Config 字段（I-A + I-B + I-G 覆盖），(b) 有人另开一个 TOML 解析者绕过结构（I-F 排除），
(c) 有人按名反射/按名查 map 取值（I-D + I-E 排除）。三条都堵 ⇒ 下表"无"字是有证明的，不是 grep 一次就定罪。

**处置图例**：✔=有消费者 · 🛑=本票"响亮失败" · 🛑★=落地前先例（已有 guard，非本票新增） ·
📋=已有计划接（指名票号） · 📄=保留但文档说明 · ⬜=结构容器（无键语义）

| 键（config.toml 路径） | 解析处 | 消费者（无 ⇒ 指回仪器） | 处置 |
|---|---|---|---|
| `schema_version` | `schema.go:108` | `loader.go:44 peekSchemaVersion`（选迁移路径；键的消费者**就是**加载管线本身） | ✔ |
| `app.language` / `app.theme` / `app.autostart` / `app.single_instance` | `schema.go:138,140,142,144` | 无（I-A：`.App.` 在包外零命中） | 📋 票 39（GUI 配置编辑器）/ 票 43（电源生命周期） |
| `app.portable` | `schema.go:149`（`toml:"-"`） | **文件里根本不许出现**：写了就是 unknown-key 错（`schema.go:145-149` 注释钉死） | 📄 无罪 |
| `ball.size` / `ball.position.{x,y,monitor}` / `ball.opacity_idle` / `ball.click_through` / `ball.hide_on_fullscreen` | `schema.go:165-173`、`155-159` | 无（I-A） | 📋 票 62/65/68（球体视觉）+ 票 39 |
| `hotkey.summon` / `mute` / `cancel` / `panel` | `schema.go:180-183` | 无生产喂入；**消费侧插槽已存在**：`ball/hotkey_reload.go` + `ball/hotkey_windows.go:384` | 📋 票 64（未完成）。⚠ 与六键**同形**（半条线），但非 🔒 段、不承载放宽语义 ⇒ 不进守卫 |
| `session.warm_timeout_sec` / `settling_sec` / `conversation_idle_sec` | `schema.go:188-190` | 无（I-A） | 📋 票 28（session scope warm） |
| `voice.*`（`enabled`、`wake_word.*`、`asr.*`、`tts.*`、`punctuation`、`conversation_mode`、`aec.*`、`realtime.{enabled,provider,model,base_url}`、`cloud_asr_chain`、`cloud_tts_chain`） | `schema.go:244-266` | 无（I-A）；`voice.realtime.api_key_ref` 例外：`loader.go:121 resolveRefs` 在包内解析进 `Resolved.RealtimeKey`，但**生产装配未取用** | 📋 票 15/26/27/41/59/60/61 |
| `audio.input_device` / `sample_rate` / `mic_muted_default` | `schema.go:271,272,277` | 无（I-A）；`mic_muted_default` 只被 `audio/gate.go:56` 的**注释**指名 | 📋 票 13 尾巴 / 票 26 |
| `audio.half_duplex=false` | `schema.go:275` | `validate.go:67` 已经拒绝 | 🛑★ 先例 |
| `llm.text_chain` / `llm.timeout_ms` / `llm.retry.{max,backoff_ms}` / `llm.roles.*` | `schema.go:411,415,306-307,297-302` | `llm/resolver.go:95-97`、`resolver.go:195-201`；`cmd/wisp/run.go:335`、`run.go:450-451` | ✔ |
| `llm.providers.<n>.{protocol,base_url,api_key_ref,compat.{loose,allow_missing_usage,extra_headers},rpm,tpm}` / `llm.providers.<n>.models.<id>.{capabilities,context_window,max_output_tokens}` | `schema.go:379-401,348-372` | `llm/resolver.go:119-134`（I-G）、三家 adapter 的 `request.go`、`cmd/wisp/providers.go` | ✔ |
| `llm.providers.<n>.{billing,plan_credit_total_micro}` / `models.<id>.{display,enabled,thinking_levels,billing}` / `price.{audio_in,audio_out}` / `quota_daily_micro` / `quota_monthly_micro` | `schema.go:391-392,350,371,359,361,341-342,367-368` | 无（`price.{in,out,cached}` 由 `agent/cost.go:56-61` 消费，其余零命中；`billing`/`thinking_levels` 只被枚举校验） | 📋 票 11（probe/目录选择）/ 票 44（C23 配额三层）/ 票 39（GUI） |
| `agent.{max_rounds,token_budget,loop_guard.repeat_thresholds}` | `schema.go:431,433,425` | 无（I-A）；消费侧插槽存在但从未被生产喂值：`agent/guard.go:28-35`（注释自称来自 config） | 📋 组合根未接（票 10 已 done 的那条线尾巴） |
| `agent.per_tool_timeout_ms` / `agent.steering_enabled` | `schema.go:436,440` | `cmd/wisp/run.go:272,325` / `run.go:327` | ✔ |
| **`risk.shell_enabled`** | `schema.go:451` | 无：I-A 只命中 `L1WindowSec`/`ConfirmTimeoutSec`；`internal/tools` 里**没有注册任何 shell.exec 工具** | 🛑 本票守卫（SPEC-07 §3 的 S3 工具） |
| **`risk.allow_shell_string`** | `schema.go:453` | 无：`risk.Facts.ShellString` **零生产写入者**（`grep -rn "ShellString\|ShellArgv\|ShellAllowlist" --include=*.go \| grep -v _test.go` 的全部命中是 `config/manager.go` 与 `risk` 的定义/读取处；生产唯一的 Facts 装配点 `tools/bridge.go:550-558` 只填 `Declared`/`Paths`） | 🛑 本票守卫 |
| **`risk.shell_allowlist`** | `schema.go:455` | 无：同上，`Facts.ShellAllowlist` 零生产写入者 ⇒ `rules_shell.go:41` 的 R6 白名单分支永不执行 | 🛑 本票守卫 |
| **`risk.blacklist_overrides`** | `schema.go:457` | 无（票 80 L2 的穷尽清单：写只有解码、读只有 `manager.go:356` 方向审计；`risk.Gate` 生产零调用点，`Classify` 签名里没有 override 形参） | 🛑 本票守卫（⇒ 票 21） |
| `risk.confirm_timeout_sec` / `risk.l1_window_sec` | `schema.go:446,449` | `cmd/wisp/run.go:258-259` | ✔ |
| `fs.allowed_dirs` / `fs.reparse_point_exceptions` / `fs.delete_enabled` | `schema.go:463,466,468` | `cmd/wisp/run.go:223-227,243` | ✔ |
| **`net.allowlist`** | `schema.go:483` | 无：`risk.NetTarget.Allowlist` 零生产写入者，且今天没有任何 web/open 工具（`internal/tools/` 无 `net/http` 工具，票 22 未开工）⇒ `rules_network.go:74` 的域名门永不执行 | 🛑 本票守卫（⇒ 票 22） |
| **`net.block_private_ranges=false`** | `schema.go:486` | 无：`rules_network.go:64` 的 `isPrivateHost` 是**无条件**判 L2，`Facts` 里根本没有承载它的字段 ⇒ 写 `false` 连"放宽"这个动作都表达不出来 | 🛑 本票守卫（要真做需 D22：SPEC-06 §6.3 钉着） |
| `net.proxy.mode` / `net.proxy.url` | `schema.go:475,477` | `cmd/wisp/run.go:448-449 chainOptions` | ✔ |
| `privacy.redact_paths` | `schema.go:496` | 无（I-A）；消费侧插槽 `observe/logging.go:50-51`、`redact.go:138`（注释自称 mirrors `[privacy] redact_paths`），生产装配未喂 ⇒ 半条线 | 📋 票 45（diagnostics guards） |
| `privacy.diagnostics_opt_in` / `retention_days` | `schema.go:498,500` | 无（I-A） | 📋 票 45 / 票 08 |
| `privacy.keep_transcript` / `keep_audio` | `schema.go:502,504` | `validate.go:78,83` 已拒绝 | 🛑★ 先例 |
| `memory.{l1_enabled,l1_max,l3_retention_days,extract_model}` | `schema.go:509-512` | 无（I-A） | 📋 票 29（memory L1/L2） |
| `panel.{enabled,width,height,keep_alive_in_session,scale}` | `schema.go:517-525` | 无（I-A） | 📋 票 33/34/36 |
| `cost.{daily_budget,monthly_budget,alert_threshold,over_budget}` | `schema.go:531-535` | 无（I-A；`alert_threshold`/`over_budget` 只被枚举校验） | 📋 票 44（C23） |
| `plugins.tier2_enabled` / `plugins.<id>.{enabled,capabilities,net_allowlist,host_api}` | `schema.go:560,543-549`、`parse.go:92 buildPlugins` | 无：`internal/plugin` **不 import 本包**（I-C），`grep -rn "\.Plugins\."` 包外零命中 | 📄 **同类、但本票不扩面**，理由见 L4，交编排者拍 |
| `models.dir` / `models.mirror` | `schema.go:569,572` | 无（I-A） | 📋 票 14 尾巴 / 票 56 |
| `models.verify_signature=false` | `schema.go:575` | `validate.go:92` 已拒绝（本票守卫的形状就是抄它） | 🛑★ 先例 |
| `models.local_override` | `schema.go:577` | 无（I-A）；消费侧插槽 `models/downloader.go:65,160`，生产装配未喂 ⇒ 半条线 | 📋 票 14/56 |
| `observe.level` / `roll.{size_mb,days}` / `slo_sample_interval_sec` | `schema.go:589,582-583,593` | 无（I-A）；`observe/logging.go:80` 读的是 `observe.Options.Level`，装配点未从 config 喂 | 📋 票 08/42/66 |
| （非 TOML 键）`Manager.ConfirmLocked` | `manager.go:28` 声明 | **零生产赋值**：`grep -rn "ConfirmLocked" --include=*.go .` ⇒ 非测试 3 命中全在 `manager.go` 自身（声明 25-28 + 调用 239），其余 4 处是 `manager_test.go` | 📄 它不是 TOML 键 ⇒ 无法"加载报错"；`nil` 分支已是 **fail-closed 拒绝 + `slog.Warn`**（`manager.go:239-248`），即它**不说谎**，只是没人接。接线归票 21 段 2 / 票 37 |

**本票没扫到的（明写扫不全的地方）**：
`.go` 之外的东西**不在仪器覆盖内** —— 我没有验证"某个外部脚本/前端 WebView 通过 IPC 把 config 键名当字符串读出来"这类消费
（`design/`、`scripts/`、`third_party/` 未纳入 I-A~I-F）。已做的补查：`grep -rn "blacklist_overrides|shell_enabled|shell_allowlist|allow_shell_string|block_private_ranges" --include=*.toml --include=*.md --include=*.json .`
⇒ 除 `docs/`、`.scratch/` 的文档与本票链外**零命中**（`deps.toml` 不含这些键）。
如果编排者知道有 GUI/IPC 侧按字符串取键的通道，本表结论要按那条通道复核。

`next=` AC#3 双向变异的实测数字（写在 L3），AC#5 的并跑抖动登记（L4）。

### L3 — AC#3 双向变异（2026-09-21，两次都在同一条 `&&` 链里先 grep 自证再跑）

承载行为的锚点行 = `internal/config/validate.go:30` 的 `validateUnwired(c),`（守卫唯一入口）
与 `internal/config/unwired.go` 的守卫表。

**(i) 把校验去掉 ⇒ 用例红。**
`sed -i 's/^\t\tvalidateUnwired(c),$/\t\t\/\/ MUTATION(i): guard removed/' internal/config/validate.go`
⇒ 同一链内 grep 自证：`grep -c 'MUTATION(i)' validate.go` = **1**，`grep -c 'validateUnwired(c),' validate.go` = **0**。
跑 `go test -count=1 -v -run 'TestUnwired|TestEveryLocked|TestRoundTrip|TestMigrate|TestValidate' ./internal/config/`
⇒ **rc=1**，`=== RUN` 69、`--- PASS` 21、`--- FAIL` **4 顶层 + 7 子例 = 11 条**：
`TestUnwiredSecurityKeysFailLoudly`（6 个子例全红：六个键各自"必须失败却通过了"）
+ `TestUnwiredGuardLeavesHonestConfigsAlone/tightening_a_guarded_key_past_its_default_is_still_rejected`
+ `TestUnwiredGuardFiresBeforeAnyFileWrite` + `TestUnwiredKeysStillRoundTripAtTheByteLevel`。
⇒ 用例真锚在校验上，不是靠别的东西绿。

**(ii) 把校验做宽（大棒）⇒ 必须有用例因此变红。**
注入的"大棒"= 往守卫表里加一行，把**有消费者的** `fs.allowed_dirs`（`cmd/wisp/run.go:223` 在喂 C26 canonicalizer）
也折进拒绝集合（`fires: len(c.FS.AllowedDirs) > 0`）。
⇒ 同一链内 grep 自证：`grep -c 'MUTATION(ii)' unwired.go` = **2**（注入行 + 其 `lands`）、
守卫表里 `path: "fs.allowed_dirs"` = **1** 行。
跑同一串 ⇒ **rc=1**，`=== RUN` 69、`--- PASS` 19、`--- FAIL` **6 顶层 + 0 子例**：
`TestUnwiredGuardIsScopedNotABigStick`（**设计好的那根反大棒锚点，如预期变红**：一份合法的
`[fs] allowed_dirs` + `[risk] l1_window_sec` + `[net].proxy` + `[plugins]` + 30 多个"零消费者但未守卫"的键的配置被拒了）
+ `TestEveryLockedSectionKeyIsAccountedFor`（守卫表与处置表必须逐键对齐）
+ `TestRoundTripLoadMarshalLoad` + `TestRoundTripBytes` + `TestMigrateV1Fixture` + `TestMigrateNoVersionKeyAssumedV1`。
⇒ 大棒式修法（"把所有没接线的键一律拒 / 把所有未知键一律拒"）**会被这棵树打红**，不是"改宽一点也算修好"。

**还原证明**：变异 (i) 后 `git diff --quiet -- internal/config/validate.go` **rc=0**；
变异 (ii) 后 `git diff --name-only -- internal/config/` 只剩 `unwired_test.go`
（那是我当轮**故意要提交**的锚点加强，不是变异残留），`grep -c MUTATION internal/config/*.go` **16 个文件全 0**；
随后 `go test -count=2 ./internal/config/` rc=0。

`next=` AC#5 门禁 + 并跑抖动登记 + 本票明确不做的那件事（L4）。

### L4 — AC#4 明确不做 + AC#5 门禁 + 墙钟抖动登记（2026-09-21）

**AC#4：本票刻意不做的事**（写了就等于新造一条没人审过的**放行侧**通道）：
1. **不接 `risk.Gate`**：不把 `cfg.Risk.BlacklistOverrides` 灌进 `bOverrides`，也不把生产判定链从
   `risk.Classify` 换到 `risk.Gate`（换过去还要在 `rules_gateway.go:84 case TierB:` 上开豁免位 ⇒ 冻结面）。
2. **不给 `Classify` 开 override 形参**（`internal/risk/blacklist.go:69` 的签名不动）。
3. **不动审批队列 / 不动 `Manager.ConfirmLocked` 的 nil=fail-closed 语义**，不新增任何"L2 强确认之后把文件记进某张 map"的容器。
   ⇒ 这三件事属于 **ticket 21**（`assessor.go:144-146` 冻结注释原文 "exemption flow lands with the approval queue, ticket 21"），
   且按票 80 的 **AC#2 裁决 / A53②** 的 (A)(B) 两分支，任何静态预授权形态**都要 D22 人工批准**（要改 `SPEC-06 §4.1` 的
   "豁免 = 一次 L2 强确认 + 写日志" 与 `PLAN.md:2396-2399` 的"必须由人点"）。本票只走 (C)。

**AC#5 门禁（只跑本票碰到的包 `internal/config/`）**

| 命令 | rc | 真实输出 / 计数 |
|---|---|---|
| `gofmt -l internal/config` | **0** | 空输出（16 个文件无一待格式化） |
| `gofumpt -l . tools/d22scan tools/mockllm` | **0** | **空输出**（全仓口径，含别人在飞的文件；`gofumpt` 不在 PATH，用 `$(go env GOPATH)/bin/gofumpt`） |
| `go vet ./internal/config/` | **0** | 无输出 |
| `GOOS=linux go vet ./internal/config/` | **0** | 无输出（Linux 编译哑弹已查） |
| `go test -count=2 -v ./internal/config/` | **0** | `ok github.com/CarlosShao/wisp/internal/config 3.050s`；`=== RUN` **194** = count=1 的 97 的**整 2 倍**（N 倍核对过）；`--- PASS` **106**（=53×2）、`    --- PASS` **88**（=44×2）、`--- SKIP` **0**、`--- FAIL` **0** |
| `go test -count=2 -v ./internal/config/ ./internal/risk/`（两包并跑） | **1** | `ok internal/config 7.050s` ／ `FAIL internal/risk 14.870s`；总 `=== RUN` **492**、`--- PASS` 278 + `    --- PASS` 210、`--- SKIP` **2**（同一测试 `TestSyncRegistryProbeLive` × count=2）、`--- FAIL` **2**（`TestResolvePerCallBudget` × count=2） |

**四种假绿逐跑点名**：`--- SKIP` 只有并跑那次的 `TestSyncRegistryProbeLive` ×2（**不是本票的包**，票 82 家族既有条件跳过）；
本票包内 **0 条 SKIP、0 条 FAIL**；没用 `-run` 收窄做交付判据（`-run` 只用在变异与复核上，
且每次都数了 `=== RUN` 条数、没有出现 `no tests to run`）；`-count=2` 按上面两行做了 N 倍核对；
没有步骤被静默跳过。

**⚠ 墙钟抖动登记为外项（不调那个 1ms 预算，一票没动它）**：
`internal/risk/pathresolver_budget_norace_test.go:34,37` 的 `TestResolvePerCallBudget`
在两包并跑时 **1.945 / 2.420 ms/op**（预算 1ms，852/447 samples）⇒ 两条 `--- FAIL`；
本轮**连单跑都红**：`go test -count=2 -run TestResolvePerCallBudget ./internal/risk/` 给出 **1.356 / 1.194 ms/op**（739/1352 samples）。
⇒ 比票 80 L3 记的数字更糟，但形状一致（这台机器此刻同时有 agent-ticket70-d / 77 / 81 在编译与跑测试）。
**本票不改它、不放宽它、也不把它算进自己的门禁结论**：它是票 18 AC#6 的墙钟性能契约，
处置建议仍按票 80 交回的那两条（`testing.Short()` 门控，或改成同机基线的倍数），**单独一票**。
⇒ 结论：`internal/config` 侧全绿；`internal/risk` 侧的红**与本票改动无关**（本票零文件落在 `internal/risk/`）。

**交回编排者的两行键表限定语（本票不写 `docs/`，原文照抄可用）**
- `docs/PLAN.md:2732`（🔒 `[risk]` 行，接在 `blacklist_overrides[]` 之后）：
  > ｜**其中 `shell_enabled` / `allow_shell_string` / `shell_allowlist` / `blacklist_overrides` 四项今天没有消费者（B 档豁免与 shell 工具都尚未实现），写入非默认值即加载报错，不静默接受**（票 83；先例 `verify_signature`）
- `docs/specs/SPEC-03-config-secrets-envs.md:34`（🔒 `[risk]` 行，同位置）：
  > （`shell_enabled` `allow_shell_string` `shell_allowlist` `blacklist_overrides` **未实现 ⇒ 写非默认值加载报错，指名由谁实现**，票 83）
- `docs/specs/SPEC-03-config-secrets-envs.md:36`（🔒 `[net]` 行，本票扫描新增的两键）：
  > （`allowlist` **未实现 ⇒ 写入即加载报错**，票 22 落地时移除；`block_private_ranges` **写 false 报错而非生效**：R5 的私有网段判定是硬编码的，SPEC-06 §6.3，关掉它需 D22）

`next=` 交回编排者：① `[net]` 两键是否算本票扩面（票 80 只点了 `[risk]` 三候选）；② `[plugins]` 五键要不要照同一判据一起守卫；
③ `SPEC-03:34/36` 与 `PLAN.md:2732` 的限定语请由编排者落；④ `TestResolvePerCallBudget` 单独立票。
写码代理这边**没有未闭的框**，等对抗验收；验收判据 = L2 表逐键 + L3 两个方向的变异数 + L4 门禁表。

### L5 — 我自己造的一枚 CI 哑弹：`a95ee3a` 在 Go 注释里放了 U+1F512（D22 ban #8 零 emoji）（2026-09-21）

票 70 的 `564be23` 抓到并点名：`internal/config/unwired.go:13` 的注释写着 `locked (🔒) section`。
它**没替我修**（简报禁改区点名 `internal/config`，且 `a95ee3a` 是几分钟前的活人产物）⇒ 修法在我这边。
本轮改掉（改写成 ASCII "the four marked rows of SPEC-03 sec 3"）后逐跑复核：

- **独立仪器（不信扫描器的嘴）**：`find internal cmd -name '*.go'` = **329** 个文件，
  用 ban #8 的同一字符类 `[\x{1F000}-\x{1FAFF}\x{2600}-\x{27BF}\x{2B00}-\x{2BFF}\x{FE0F}\x{1F1E6}-\x{1F1FF}]`
  全量扫 ⇒ **命中 0 处**；我这 5 个 Go 文件单独扫也 ⇒ **0 处**。
- **门禁五项复跑**（改注释之后）：`gofmt -l internal/config` 空 rc=0、
  `gofumpt -l . tools/d22scan tools/mockllm` 空 rc=0、`go vet ./internal/config/` rc=0、
  `GOOS=linux go vet ./internal/config/` rc=0、`go test -count=2 ./internal/config/` rc=0
  （`=== RUN` 194 / `--- PASS` 106 / `    --- PASS` 88 / `--- SKIP` **0** / `--- FAIL` **0**）。
- `sh scripts/d22scan.sh` 现在 **exit 2**，但**红点已不在 ban #8**：
  `scope ban #8 internal/ examined 305 Go files`、`cmd/ examined 24 Go files` 均无 finding；
  exit 2 的原因是 `frontend/` 目录被票 77 建出来了、而 ban #6 还登记着 "absent-but-exempt"
  ⇒ **与本票无关，交编排者/票 70 处置**（它自己的输出文字已经写明"flip it to live:true in the same commit that creates the tree"）。
- ⚠ 顺带提醒：**票面 markdown 里继续用 🔒/⚠ 是安全的**（ban #8 走 `.go` 与声明的文本树，`docs/`、`.scratch/` 不在 `internal/`+`cmd/` 口径内），
  但**任何 Go 注释里都不许**——这张票的教训。

`next=` 交回编排者（L4 末尾的四项不变），另加第 ⑤ 项：**`sh scripts/d22scan.sh` 的 ban #6 "absent-but-exempt" 需要与 `frontend/` 同 commit 翻成 `live:true`**（票 77 刚建了树）。





### L1 — AC#2 校验先落地（承载体在，表随后补全）（2026-09-21）

新增 `internal/config/unwired.go`：`unwiredKeys` 六行守卫表 + `validateUnwired()`，挂在 `validate()` 里
（次序 `validateModels` 之后、`validateCost` 之前），照 `verify_signature` 的形状：
**只在非默认值上开火**（四个 `[risk]` 键 + `net.allowlist` + `net.block_private_ranges`），
错误文案 = 键路径 + "哪个消费者不存在"（点名到符号）+ "由谁落地"（票号 / SPEC 行 / 需 D22）。

副作用次序已由用例证明：`applyMigrations` 的 dry-check 走同一个 `validate()`，
所以一张写着 `shell_allowlist` 的 **v1** 文件被拒时 `config.toml` 逐字节未变、`.bak-1` 未生成、
临时文件未出现（`TestUnwiredGuardFiresBeforeAnyFileWrite`）。

改动面（5 个文件，全在 `internal/config/`）：
`unwired.go`(新) `unwired_test.go`(新) `validate.go`(挂 + 一处注释) `boundary_test.go`(fullConfig 去掉 6 行非默认值)
`migrate_test.go`(v1 fixture 的 `[risk]` 穿透断言从 `shell_allowlist` 换成 `l1_window_sec`，
因为前者现在按定义不可加载)。

用例（`unwired_test.go`）：`TestUnwiredSecurityKeysFailLoudly`(6 子例) / `TestUnwiredGuardLeavesHonestConfigsAlone`(4 子例)
/ `TestUnwiredGuardFiresBeforeAnyFileWrite` / `TestUnwiredGuardIsScopedNotABigStick`(AC#3(ii) 的锚)
/ `TestUnwiredKeysStillRoundTripAtTheByteLevel` / `TestEveryLockedSectionKeyIsAccountedFor`(AC#1 完整性的可执行版)。

门禁实测（`gofumpt` 在 `$(go env GOPATH)/bin/gofumpt`，不在 PATH）：
`gofmt -l internal/config` 空 rc=0；`gofumpt -l . tools/d22scan tools/mockllm` **空 rc=0**；
`go vet ./internal/config/` rc=0；`GOOS=linux go vet ./internal/config/` rc=0；
`go test -count=2 -v ./internal/config/` rc=0，`=== RUN` **194**（= count=1 的 97 的整 2 倍，N 倍核对通过）、
`--- PASS` **106**、`    --- PASS` **88**、`--- SKIP` **0**、`--- FAIL` **0**，`ok ... 0.782s`。
本轮**没有**并跑 `internal/risk` ⇒ 未复现 `TestResolvePerCallBudget` 的墙钟抖动（AC#5 里单独量一次并跑）。

`next=` 把 AC#1 的键表全量写进本票（含间接消费排查证据），再跑双向变异与并跑门禁。
