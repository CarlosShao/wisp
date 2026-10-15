# 模型页 v3 简报（给前端 agent）——2026-10-14 机主反馈：**StepFun 双接入模式 + 模型目录缺失**

> 来源＝机主 2026-10-14 21:2x 原话（截图状态：stepfun 行"已设置基底/密钥已保存( dpapi:… )/**暂无模型**"，默认聊天模型"当前：暂未读取"）：
> 「audio3 是通过 api 方式调用的，当前我选择的是套餐的 plan 连接方式…stepfun 有两种方式接入，一种是 plan，一种是 api 计费，audio3 是使用 api 的时候才免费；当前界面只能选择一种方式去添加 stepfun，所以我填的是 plan 的，`https://api.stepfun.com/v1` 这个才是 api 的。而且这个界面同样也根本不符合主流 harness 的模型接入方式。」

## 1. 两条功能诉求（本简报要解决的事）

1. **同一厂商要能有两行、两种计费模式**：StepFun 的"套餐(plan)"入口与"API 计费"入口是不同的基底；**语音模型 audio3 只有在 API 计费一侧才免费** ⇒ 文本腿挂 plan 行、语音腿（audio3）必须挂 API 行（`https://api.stepfun.com/v1`）。今天一行只有一个基底 ⇒ 只能二选一，audio3 无路可走。
2. **模型目录为空**：行内"暂无模型" ⇒ "默认聊天模型"选择器没有任何候选项（机主已存 base+key 仍卡死在这一步）。

## 2. Go 侧现状（编排者 21:2x 逐处现读；**这是契约面的事实，别按猜测写 UI**）

- 供应商行结构（`internal/config/schema.go:379-400`）：`protocol` / `base_url` / `api_key_ref` / **`billing`（枚举 `pay-per-token|plan`，字段已存在）** / `plan_credit_total_micro` / `compat` / `rpm` / `tpm` / **`models{}`（行内模型目录，按模型 id）**。
- 预设 12 家含 stepfun，其**默认基底就写作 `https://api.stepfun.com/v1`（＝API 口）**（`internal/config/defaults.go:18-36`）——机主填的 plan 地址是覆盖它。
- **非预设名合法**：`internal/config/validate.go:209-212` 逐字——"非预设 provider **必须显式写 protocol**" ⇒ **数据层允许 `stepfun-plan` / `stepfun-api` 这类"同名厂商、两行"**。
- 写入白名单现况：
  - 248 表七枚：`provider_base_url` / `provider_api_key_ref` / `model_context_window` / `model_price_in` / `model_price_out` / `role_chat_model` / `provider_credential`（均带 provider/model 选择列）。
  - 313 表六枚：`provider_row_create`（**选择器只吃预设名**，页面从此而来） / `role_chat_pair` / `voice_cloud_asr_chain` / `voice_cloud_tts_chain` / `voice_realtime_provider` / `voice_realtime_model`。
  - ⇒ **缺口（Go 侧、待批）**：① `billing` 没有写入字段；② `provider_row_create` 建不出"自定义名"的第二行；③ 已存在的行没有"添加模型条目"的写入路径（只有建行时可选带一枚）。**三条都是接口面变更 ⇒ 需机主点头后由编排者落地**。

## 3. 要前端页做的（在拿到 §4 点头之前，先做不依赖新写入的部分）

- 供应商行卡片：显示 `billing`（套餐 / 按量）与基底（只读即可，写入待 §4）；**"再加一行"入口**（目标：`stepfun-plan` 与 `stepfun-api` 两行并存，各自 base/key）；
- 行内**模型清单**：列已有模型（id/窗口/价格），"添加模型"入口（若 §4 批了 `model_row` 写入就直接接；未批则先做骨架 + 占位文案）；
- **语音页**：ASR/TTS 链与 realtime 两枚叶子要**显式指出挂在哪一行**（让 audio3 挂 API 行）；
- 文案沿用 v2 简报的"去内部俚语"口径；空/加载/错误三态齐。
- ⛔ 不新增任何方法名/字段名去写配置（白名单之外一律 `停手上报`）。

## 4. 待机主点头的三件（点头后编排者落地，前端再接写入）

1. `provider_billing` 写入字段（套餐 vs 按量）；
2. 建行支持**自定义名**（或"复制 stepfun 为 stepfun-api"）；
3. 模型条目写入 **或** "从 `/v1/models` 拉取列表"端口（二选一；后者能一并解决"模型真名单"的老问题）。

## 5. 归属

页面＝**前端 agent**；§4 的写入面＝**编排者**（批准后实现，含 313 的字段门与契约测试）。

---

## 6. ★机主批准（2026-10-14 22:5x，逐字落账 A968）——本简报的 §4 三件已获批

机主原话逐字：「**模型列表要允许拉，三个写入能力可以做**」＋方向逐字：「你现在不要因为 stepfun 去妥协，你最主要的是把各处的页面功能项什么的都给尽善尽美……按照票上来」。

⇒ **§4 待批三件全部获批**：① `provider_billing` 写入；② 建自定义名供应商行（`stepfun-plan`/`stepfun-api` 两行并存）；③ 模型条目写入 **＋** "从 `/v1/models` 拉取列表"端口。Go 侧落地票＝`.scratch/wisp/issues/320-model-page-v3-go-side-*.md`（设计腿先出形状）。
⇒ **给前端 agent 的口径（机主方向）**：页面按"主流 harness 的完整形态"做——供应商列表（多行/计费模式/模型清单/拉取）＋模型目录＋角色映射；**不做**"为了能录 key 就应付"的简版。Go 侧新字段/新方法落地后再接写入（本简报 §3 的只读展示可先行）。

---

## 7. ★接口面定稿（编排者 2026-10-15 12:0x 裁定，落账 A972；依据＝设计腿 `320-a1` 交件 `.scratch/wisp/probes/320/a1/census.md`，编排者已逐条复跑核过）

- **三个写入字段**（都骑**既有** `config.set` 方法，⛔ 不新增方法名）：`provider_billing`（值 `plan` / `pay-per-token`）· `provider_row_create_custom`（值 `<protocol>,<base_url>`，可带首枚模型 id）· `model_row_create`（provider+model，无值）。⇒ 落地后页面的"计费切换 / 再加一行 / 添加模型"三个动作就能接真写入，**前端不必为它们改方法名**（写入仍走既有 `requestConfigWrite` 形状，只是 `field` 取三个新名）。
- **"拉取模型列表"端口＝新方法名 `config.models.list`**（⛔ 不叫 `models.list`：那个名字在 `SPEC-08 §5.2` 已属本地签署模型仓 C29/D26，语义相撞；⛔ 不带 `panel.` 前缀：`panel.*` 四枚白名单钉死）。请求载荷复用既有键 `configProvider`（⛔ 不发明新信封键）；回执＝散文＋**一枚标记行** `wisp-models-list:` 后跟 JSON（照 `wisp-config-view:` 的先例，`config_handlers.go:546-555`）；失败只给 D37 分类文案，**不回显上游原文/密钥**。

### 7.1 给前端 agent 的**一行活**（请机主转交；这是 Go 半落地的前置）

`frontend/src/lib/panel.ts` 加**一枚**导出发送函数（形状照该文件既有 `requestConfigView`／`requestConfigWrite` 那几处，用现成 `requestHostAnswer`）：

```ts
/** Ask the host to fetch one provider's model catalog (config.models.list). */
export function requestProviderModels(provider: string): Promise<string | null> {
  return requestHostAnswer({
    method: "config.models.list",
    requestId: nextRequestId(),
    source: "panel-composer",
    configProvider: provider,
  });
}
```

三条硬约束（仓里既有测试钉死的，不是我加的口径）：① 方法名**必须写成字面量**（`sendRequest`/`requestHostAnswer` 的第一参数不许算出来）；② **不许**新开 `postMessage` 调用点（恰 2 处）／不许 `fetch(`／不许 `new WebSocket`；③ 用 `config.` 前缀，⛔ 别写成 `panel.`（四枚白名单钉死）。

页面侧（模型页 v3 顺手做）：行内模型区加**"拉取模型列表"**按钮；Go 侧是**同步**调用（预算 ≤10s）⇒ 按钮要有**加载态**；拉到后列 id 供勾选，逐枚走 `model_row_create` 写入（§7 第一行）。

**时序**：前端这一行**先落**即可（按现读，`config.*` 字面量不进 `panel.*` 两道扫描面，落它不会让任何测试变红；此时点按钮会收到"未知方法"的分类拒绝，属预期）→ Go 半（`320-r2`：方法名注册四点＋拉取实现）随其后，两侧合起来 `TestEveryWhitelistedInboundNameIsAlsoSentByThePage` 才绿。
