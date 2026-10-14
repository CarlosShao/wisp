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
