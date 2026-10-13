# 313-panel-ask — 给设置页的接口规格（票 313 Go 侧，`313-r1` 交）

**这份是谁给谁的**：`313-r1` 落地腿 → 编排者转交机主外部 agent（页面在 `D:/wt/fe`，分支 `dsh/feat/frontend-p0-v2`）。**Go 侧已落**（`internal/panel` ＋ `internal/config`）；`cmd/wisp` 的装配那一跳还没落（见 §5）。行号/形状以本件为准，`census d-page.md` 那 15 行是**改前的料**。
> ⚠ **更正（编排者 2026-10-13 深夜）**：前端工作树**已并入 dev**（`36b44ee4`，10-13 09:39；该分支 tip `16c2f038` 已在 dev 内、**零未合并提交**，落后 dev 2562 笔＝纯残留）⇒ **页面工作直接在主仓 `dev` 上继续，不再走 `D:/wt/fe` 分支循环；该 worktree 退役**。原句「页面在 `D:/wt/fe`，分支 `dsh/feat/frontend-p0-v2`」为 10-13 合并前快照，已过期（⛔ 不照它做）。

**⛔ 页面侧不许做的事照旧**：⛔ 动 CSP `<meta>`、⛔ 新方法名（C17）、⛔ 解析回执里的散文、⛔ emoji、⛔ 把密钥值放 `configValue`。

## 1. 读：`config.get` 回执 = 老两句散文 ＋ 一行结构化载荷

调用形不变：`window.wispDispatch(JSON.stringify({method:"config.get", requestId:"pc-…", source:"panel-composer"}))`。
回的是一个**多行字符串**：第 1 段是票 248 冻结的散文（逐字未变："配置可读。服务商 N 家；凭据状态 X。" / "当前聊天模型 Y。" / 每行的 `- 名：引用 …，目录模型 N 个，密钥 …`；不可读时为"配置当前不可读，…"那一句），**再一行**以标记开头：

```
wisp-config-view:{"readable":true,"credential":"no_ref_declared", …}
```

- 标记常量＝`panel.SettingsPayloadMarker` = **`wisp-config-view:`**（Go 侧导出，页面复制这个字面量即可）。取回执里以它开头的那一行，`strings[i+len(marker):]` 走 `JSON.parse`。
- **⛔ 别解析散文**；散文只是给人看的。载荷是 JSON 对象，键**冻结**（`config_313_test.go` 逐键钉住）：

| 键 | 型 | 含义 |
|---|---|---|
| `readable` | bool | 配置读到了没有。`false` 时其余键是空壳，屏上照旧显示"不可读（这不是未配置）" |
| `credential` | string | 五值之一：`unknown`/`config_unreadable`/`no_ref_declared`/`all_recorded`/`partly_missing` |
| `presetsKnown` | bool | **腿答了没有**。`false` ⇒ ⛔ 别把 `presets:null` 当成"没有预设" |
| `presets` | string[] \| null | 内置服务商预设名（**含 `stepfun`**，字母序）＝**服务商下拉的数据源** |
| `configuredKnown` / `configured` | bool / bool | "用户配过没有"。`configured` 真＝服务商行非空或 `roles.chat.provider` 非空 |
| `chatProvider` / `chatModel` / `chatModelKnown` | string/string/bool | 当前聊天选中的**两半** |
| `providers[]` | obj[] | 每行：`name`、`apiKeyRef`（引用名，随机 id，**不是值**）、`hasBaseURL`、`modelCount`、**`models`（string[] \| null，模型名清单＝模型下拉的数据源）**、`credentialSet`、`credentialError` |
| `voice` | obj | `cloudASRChain`、`cloudTTSChain`（string[] \| null）、`realtimeEnabled`、`realtimeProvider`、`realtimeModel` |

- **⛔ 编模型名**：仓里 `stepfun` 名下**零枚**真实模型 id（`census c-default.md` C4 现量）。`models` 里只会有用户自己写进去的 id。⇒ 模型下拉**必须带"手填模型 id"退路**；真 id 要预置，得先有"发现→写入"那条口（`wisp providers discover stepfun`，现在只打印不落盘）。
- 键名/值里没有任何凭据值（`apiKeyRef` 是 `dpapi:<id>` / `env:NAME`，可安全显示）；`credentialError` 是状态字符串，可显示。

## 2. 默认高亮 stepfun（`A877` 裁定③）

- 规则：**`configuredKnown === true && configured === false`** 时，模型/服务商控件**默认高亮 `stepfun`**；否则一律按 `chatProvider`/`presets` 走。
- 只读＋默认高亮：**⛔ 不点保存不写任何东西**；"用户改过以用户为准"由这一条天然成立。
- ⛔ 不写首建文件（票 198 钉：首建文件禁出现 `[llm.providers`）、⛔ 不给 `default` 标签/运行时兜底（catalog 缺行即拒载）。Go 侧只交"预设清单＋是否已配置＋模型清单"这三样**读**。

## 3. 写：`config.set` 可用字段（白名单，门先于腿）

写一发＝**一枚字段**：`{method:"config.set", requestId, source:"panel-composer", configField:<名>, configProvider?, configModel?, configValue?}`。

| `configField` | 选择器 | 值 | 落到哪 |
|---|---|---|---|
| `provider_base_url` | `configProvider` | URL | `llm.providers.<p>.base_url` |
| `provider_api_key_ref` | `configProvider` | `dpapi:<id>` / `env:NAME` | `llm.providers.<p>.api_key_ref` |
| `model_context_window` | `+configModel` | 整数 | `…models.<m>.context_window` |
| `model_price_in` / `model_price_out` | `+configModel` | 数字 | `…models.<m>.price.in/out` |
| `role_chat_model` | — | 模型 id | `llm.roles.chat.model`（只写 model 半） |
| **`provider_row_create`** | `configProvider`（**必须**）；`configModel` 可选 | — | 建 `[llm.providers.<p>]`（带 model 时同笔建 `…models.<m>`，`enabled=true`）。**只收内置预设名**；非预设名会被拒（无 protocol 字段，写出来下次加载不过）；**行已存在＝拒**（此步只建不改） |
| **`role_chat_pair`** | `configProvider` ＋ `configModel` | — | **一笔写 `llm.roles.chat.provider` 与 `.model` 两键**（`set both or neither`：`validateRoles` 两种半态都拒，两支单键 setter 从空态根本到不了合法态） |
| **`voice_cloud_asr_chain`** / **`voice_cloud_tts_chain`** | — | 链串，元素 `provider/model`，多元素**逗号分隔** | `voice.cloud_asr_chain` / `voice.cloud_tts_chain`。每元素过 `checkChainElement`：provider 必须在 registry、model 必须在该 provider 目录里，**`local-sherpa` 一律拒**（它是本地兜底、永不进云链）。空串会被拒（见 §5 缺口①） |
| **`voice_realtime_provider`** / **`voice_realtime_model`** | — | 名 / id | `voice.realtime.provider` / `.model`。β 双字段形（`A877` ①）。`enabled` 不可写（见 §5 缺口②） |
| `provider_credential` | `configProvider` | 值**只许**放 `credentialValue` 键 | 值 → DPAPI blob，`config.toml` 只出现 `dpapi:<id>` 引用 |

- **⛔ 值不许放 `configValue`、不许进共享封套其它键**（共享封套只声明 `configField/configProvider/configModel/configValue` 四枚，页面/消息路都能看见）。
- 拒因三句照原文（`internal/config/settings.go:77-81` 三种拒因：文件没建／行不存在／校验不过）；**成对拦在页面**（`set both or neither`）。
- **回执即事实**：写入回执带"实际落盘的键路径 ＋ 生效级别"，现量所有设置写入＝**要重启**（`panel_config_store.go:228`，理由逐字见该文件头）。⛔ 页面⛔ 自编"已生效"。
- 推荐流程（文本模型）：`provider_row_create`(provider=`stepfun`, model=手填 id) → `provider_credential`(录 key) → `role_chat_pair`(provider=`stepfun`, model=手填 id)。语音：先建行，再 `voice_cloud_asr_chain`/`voice_cloud_tts_chain`。

## 4. 密钥录入（`AC#3`，机制未变）

`config.set` + `configField:"provider_credential"` + `configProvider:<名>`；值**只能**在 `credentialValue`。输入框提交后立刻清空、不存 state、屏上任何位置不渲染值/尾号；状态只显示「已录入/未录入」+ 引用名（`dpapi:<id>` / `env:NAME`）。回执句（`renderCredentialReceipt`）可直接展示。

## 5. 本腿**具名没做**／留给裁定（⛔ 页面别自己造）

1. **清空云链没有界面口**：链字段的值为空会被门按"空值不写成默认"拒（沿用既有纪律，⛔ 未新造 clear 哨兵）。要清空 ⇒ 需编排者给一枚显式字段或归后续票。
2. **`voice.realtime.enabled` 不可写**：`A877` ① 的实时段名册只有 provider/model（＋随需要的 api_key_ref/base_url），**没有 enabled** ⇒ 页面配得成对，但点不亮实时头。
3. **`cmd/wisp` 装配那一跳未落**（`ReadSettings` 填字段 ＋ switch 六 case），本腿**未动**该面（该面在飞件归 `294-r1`）；精确落点包＝`.scratch/wisp/probes/313/r1/10-cmdwisp-landing-package.md`。在它落地前，读载荷恒为零值（"腿没答"），写字段会被腿按"没有对应的写入腿"拒——**这是实话，不是 bug**。
4. **`role_chat_pair` 是本票唯一"一字段写两键"**：`A877` ① 为实时段守住了 one-key-per-call；`roles.chat` 两种半态都被 `validateRoles` 拒 ⇒ 那里 β 形结构上不存在。此形**具名上报待裁**。
5. 首启引导：产品 `App` 今天**没挂** `FirstRunScreen` ⇒ `AC#1` 的落点就是设置页。

## 6. `census d-page.md` 那 15 行的去向（对账）

1 「模型」类 waiting 文案 → 换控件区：**本件 §1/§2** 给数据源，⛔ 不动其它 11 类。2 读面调用形 → **§1**（读标记行，⛔ 别解析散文）。3 服务商下拉＝预设清单 → **§1 `presets`**。4 模型下拉＝该 provider 的 models 键名 → **§1 `providers[].models`** ＋**手填退路**。5 文本写路径 → **§3 `role_chat_pair`**（provider 半已有口）。6 语音写路径 → **§3 链两枚／realtime 两枚**。7 密钥录入 → **§4**。8 密钥不回显 → **§4**。9 回执即事实＋重启 → **§3**。10 成对写拦截 → **§3**（成对拦在页面，拒因照原文）。11 行存在性 → **§3 `provider_row_create`**（显式建行，唯一口）。12 快照读法 → **不属本票**（另一枚票）。13 规范（⛔ emoji／⛔ 新方法名）→ 照旧。14 读不到就说"未读到" → **§1 `*Known`**。15 stepfun 默认高亮 → **§2**。
