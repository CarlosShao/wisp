# SPEC-03 · 配置模型、机密与环境隔离

> 追溯：D6/D36/§14.5/C15/C28/F1/D33；环境隔离为【SPEC】新增（用户要求：隔离开发与生产环境）。

## 1. 背景与问题

- D6 定了「配置真相源 = config.toml，GUI 只是它的编辑器」；D36 定了 section 树与三档生效级别，
  但未落到键的完整定义、默认值管理、热加载实现与迁移细节。
- API Key 明文存 TOML 是 F1 里的高危配置问题，必须由 C28 SecretStore 取代。
- 用户要求开发/生产环境隔离（PLAN 未覆盖）→ 本 spec §5 定义 `WISP_ENV` 环境模型。

## 2. 原则（D36 六条规则，重申为实现约束）

1. 🔒 安全节（`[risk]`/`[fs]`/`[net]`/`[plugins]`）**任何放宽**必须触发一次 L2 级重新确认，
   不得经热加载静默生效；收紧方向可热加载。
2. 未识别的键 → 报错并指出键名与所在行，不得静默忽略。
3. **默认值唯一真相源是 Go 结构体 tag**；TOML 只是序列化。
4. GUI 与 TOML 双向一致由「读写同一个结构体」保证；GUI 不得持有 TOML 里没有的状态。
5. 敏感项存引用不存值：`api_key_ref = "dpapi:<blob-id>"` 或 `"env:NAME"`；明文 `api_key` 字段禁止。
6. `config.toml` 本身进 A 档黑名单（读写都禁，F1）。

## 3. config.toml 全量 section 树（D36；【SPEC 细化】补类型与默认值）

| section | 键（类型 = 默认） | 生效级别 |
|---|---|---|
| `[app]` | `language(string)="zh-CN"` `theme("dark"\|"light"\|"auto")="dark"`（hot）`autostart(bool)=false` `single_instance(bool)=true` `portable(bool)`（只读，由 portable.txt 决定） | restart（theme 除外） |
| `[ball]` | `size(int)=56`（44–72）`position{x,y,monitor}` `opacity_idle(float)=0.35` `click_through(bool)=true` `hide_on_fullscreen(bool)=true` | hot |
| `[hotkey]` | `summon(string)` `mute(string)` `cancel(string)="Esc"`（Confirming 期间临时接管，会话结束归还）`panel(string)` | hot（须重注册） |
| `[session]` | `warm_timeout_sec(int)=90` `settling_sec(int)=3` `conversation_idle_sec(int)=30` | hot |
| `[voice]` | `enabled(bool)=true` `wake_word{enabled=false, keywords[]=[], thresholds[]=[], veto_words[]=["取消","停下","别"]}` `asr{provider="local-sherpa", model=<id>}` `tts{provider="local-sherpa", voice, speed(float)=1.0}` `punctuation(bool)=true` `conversation_mode(bool)=false` `aec{enabled(bool)=true, echo_ref="self-render"}` | reload（换模型）/ hot（阈值、否决词） |
| `[audio]` | `input_device(string)="default"` `sample_rate(int)=16000` `half_duplex(bool)=true`（**硬编码 true，只读显示，写 false 报错**）`mic_muted_default(bool)=true` | hot（须重开设备） |
| `[llm]` | `default_provider(string)` `fallback_provider(string)` `timeout_ms(int)=60000` `retry{max(int)=3, backoff_ms(int)=1000}` `temperature(float)` `providers.<name>{base_url, model, api_key_ref, context_window(int)}` | hot（api_key_ref 变更须重新解密） |
| `[agent]` | `max_rounds(int)=50` `token_budget(int)=200000` `per_tool_timeout_ms(int)` `loop_guard{repeat_thresholds=[3,5,8]}` `steering_enabled(bool)=true` | hot |
| 🔒 `[risk]` | `confirm_timeout_sec(int)=300` `l1_window_sec(int)=2` `shell_enabled(bool)=false` `allow_shell_string(bool)=false` `shell_allowlist[]=[]` `blacklist_overrides[]=[]`（B 档豁免） | 🔒 放宽需 L2 重新确认 |
| 🔒 `[fs]` | `allowed_dirs[]=[]` `reparse_point_exceptions[]=[]`（按具体路径）`delete_enabled(bool)=false` | 🔒 同上 |
| 🔒 `[net]` | `allowlist[]=[]` `block_private_ranges(bool)=true` `proxy{mode("system"\|"none"\|"manual"), url}` | 🔒 同上 |
| `[privacy]` | `redact_paths(bool)=false` `diagnostics_opt_in(bool)=false` `retention_days(int)=30`；`keep_transcript`/`keep_audio` **硬编码 false，只读显示，写入 true 报错** | hot（硬编码项不可改） |
| `[memory]` | `l1_enabled(bool)=true` `l1_max(int)=20` `l3_retention_days(int)=30` `extract_model(string)` | hot |
| `[panel]` | `enabled(bool)=true` `width(int)=640` `height(int)` `keep_alive_in_session(bool)=true` `scale(float)` | hot |
| `[cost]` | `daily_budget(micro)` `monthly_budget(micro)` `alert_threshold(float)=0.8`；达 100% 默认暂停新任务（可配 `over_budget("pause"\|"warn")="pause"`） | hot |
| 🔒 `[plugins]` | `tier2_enabled(bool)=false` `<id>.{enabled, capabilities, net_allowlist, host_api}` | 🔒 变更即重载该插件并重新确认能力 |
| `[models]` | `dir(string)` `mirror[]=[]` `verify_signature(bool)=true`（**不可关：写 false 报错而非生效**，C29）`local_override{}={model-id: 本地路径}` | hot |
| `[observe]` | `level("debug"\|"info"\|"warn"\|"error")="info"` `roll{size_mb(int)=10, days(int)=7}` `slo_sample_interval_sec(int)=5` | hot |

内置 LLM provider 预设（D8）：`openai` `anthropic` `deepseek` `qwen` `zhipu` `moonshot`
`siliconflow` `openrouter` `ollama`——预设只提供 base_url/model 默认，Key 一律 `api_key_ref`。

## 4. 加载、热加载与迁移

### 4.1 加载顺序

`schema_version` 检查 →（需要时）迁移 → 解析 + 校验（未知键报错指出行号）→ 合入结构体默认值 →
`api_key_ref` 解析（DPAPI 解密或 env 读取，失败即 `Unconfigured`）。

### 4.2 热加载范围

- 可热加载：`[llm]` `[agent]` `[session]` `[ball]` `[hotkey]` `[panel]` `[cost]` `[memory]`
  `[observe]` `[models]`（阈值/开关类）。
- reload 级：`[voice]` 换模型（走重载子系统路径，触发模型加载/卸载）。
- restart 级：`[app]`（theme 除外）。
- 🔒 安全节放宽：**弹 L2 重新确认卡（原生侧），确认后生效并写日志；拒绝则保留旧值**。

### 4.3 变更检测实现

【SPEC】`watchdog` goroutine 的周期循环（1s tick）顺带检查 `config.toml` 的 mtime + size 变化
（不新增 goroutine，不加 fsnotify 依赖）；变更 → 交给 `config` 模块按 §4.2 处理。
GUI 改配置 = 写回 TOML → 同一检测路径生效（单一真相源始终是文件）。

### 4.4 schema 迁移

- `config.toml` 带 `schema_version` 键；升级时自动迁移并**备份原文件**为 `config.toml.bak-<ver>`。
- 无法迁移 → `Unconfigured` 态 + 明确指引；**不得静默用默认值覆盖**（会重置目录白名单——安全问题）。

## 5. 环境隔离模型（【SPEC】新增，实现开发/测试/生产隔离）

### 5.1 环境判定

```
WISP_ENV ∈ { prod(默认) | dev | test }
```

- 安装产物默认 `prod`；`go run`/本地调试默认 `dev`（由 `buildinfo` 里的构建期默认值 +
  环境变量覆盖）；CI 与测试进程显式 `WISP_ENV=test`。
- CLI 同样受 `WISP_ENV` 影响（`wisp run` 在 dev 环境打 dev 数据目录）。

### 5.2 每环境资源分叉

| 项 | prod | dev | test |
|---|---|---|---|
| 数据目录 | `%APPDATA%\wisp\` | `%APPDATA%\wisp-dev\` | `WISP_TEST_DATA_DIR` 或 `%TEMP%\wisp-test-<pid>\` |
| 单实例互斥 | `Local\wisp-single-instance` | `Local\wisp-dev-single-instance` | 不注册（测试可并行） |
| 默认日志级别 | info | debug | 可注入 |
| 自动更新检查 | 开（默认手动检查，§6） | **关** | 关 |
| LLM 默认端点 | 用户配置 | `http://127.0.0.1:18080/v1`（compose mock-llm，SPEC-11） | 显式注入 |
| 模型镜像默认 | `[models] mirror` | `http://127.0.0.1:18081`（compose model-mirror） | 显式注入 |
| KWS 常驻 | opt-in 同 prod | 默认关（开发机免打扰） | 按测试注入 |
| 诊断包导出 | 允许 | 允许 | 关（CI 无意义） |

- dev 环境首次启动生成 `config.toml`（dev 默认值 + 指向 mock 的 provider 示例，注释说明），
  与 prod 配置完全互不可见——调试不会污染自用的画像/记忆/取证/费用账本。
- 便携模式优先于环境分叉目录（portable.txt 存在 → exe 同级 `data\`；dev 仍是 `data-dev\`）。
- 环境标识必须**在悬浮球 tooltip 与面板标题可见**（`Wisp · dev`），防调试时误以为在操作真实数据。
- SLO 验收一律在 `prod` 构建口径下测（`WISP_ENV` 只分叉数据与端点，不改变运行时行为参数）。

### 5.3 SecretStore 在各环境的行为（C28）

- 存取：Windows DPAPI（`CryptProtectData`，CurrentUser scope）；blob 落数据目录 `secrets\`，
  每 ref 一个文件；日志/诊断包只留 Key 后 4 位。
- `env:` 引用：CI/测试专用（`api_key_ref = "env:WISP_TEST_LLM_KEY"`，值可以是任意占位串——mock 不验）。
- 首次启动检测到明文 `api_key` → 自动迁移到 DPAPI、从 TOML 删除、写日志告知（备份原文件）。
- 便携模式见 SPEC-02 §6（P13 建议：仅 `env:`）。

## 6. 测试决策

- 校验测试：每个 section 一组「合法/未知键/类型错/安全节放宽」用例；未知键断言报错含行号。
- 热加载测试：写文件触发 → hot 段即时生效、reload 段触发重载事件、🔒 放宽断言弹确认且拒绝后保留旧值。
- 迁移测试：构造 v-1 旧配置 → 迁移后值保留 + 备份存在；构造不可迁移文件 → Unconfigured 且原文件未被改。
- 环境分叉测试：同一进程内以三种 `WISP_ENV` 各启动一次（测试进程内模拟），断言数据目录/互斥名/端点分叉；
  dev/prod 目录互不可见。
- SecretStore：明文迁移 round-trip；DPAPI blob 换用户不可解密（用降权进程模拟）→ 明确错误非静默。

## 7. 不做什么

- 不做配置的远程下发/多配置层叠加（每环境一个 config.toml，无 overlay）。
- 不做 GUI 里的「高级模式」绕过 TOML 直接持久化状态（GUI 只是编辑器，D6）。
- 不做 `keep_audio=true` 的任何实现路径（硬编码 false，D16）。
