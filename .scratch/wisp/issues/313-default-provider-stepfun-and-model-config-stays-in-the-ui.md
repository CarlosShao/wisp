# 票 313 — 默认供应商＝**stepfun**（文本）＋模型与密钥**全部在界面上配**（⛔ 不走 CLI）＋语音模型也走 API（audio3 免费但仅 API）

**立票**：2026-10-13 08:4x 编排者（来路＝机主原话，逐字两段：
① 「我只想用界面去配置模型，不是用你这种方式，因为这个 stepfun 吧稍微有点复杂，当前因为 stepfun 新模型比较火爆，所以官方准备修改套餐了，所以文本模型自然用 stepfun 的套餐，但是语音模型因为 audio3 是免费的，但是只能接入 api 使用，所以语音模型又得用 api 方式，没招啊」
② 「所以最好的方式还是等界面出来了，我自己去界面上配置模型吧，但可以把 stepfun 作为我们默认的模型来使用，api key 该配还得配」）。
**性质**：**机主的产品决定**（不是缺陷）：① 默认供应商＝stepfun；② **配置面必须是界面**（他明确拒绝 CLI 录 key 那条路）；③ 语音模型也走 API 配置（那是 stepfun 的 audio3，免费但只在 API 侧开放）。

## 现量（2026-10-13 现跑，改前重跑）

1. `internal/config/defaults.go:28` 逐字 `"stepfun": {Protocol: ProtocolOpenAIChat, BaseURL: "https://api.stepfun.com/v1"}` ⇒ **预设里已经有 stepfun**（2026-09-19 那次批准的扩展），缺的只是"默认选中"与界面入口。
2. ★**`provider` 与 `model` 必须成对写**（`internal/config/catalog.go:95-113` 三条分支：都空＝该功能自降级；只写一枚＝`set both or neither` 报错）；且 `provider` 必须已在该 catalog 里（`:110`）、model 必须存在于该 provider 下（`:113`）⇒ **"把 stepfun 设为默认"＝至少要把 provider＋model 一起落进某处**，⛔ 只写 provider。
3. 设置写入那条路 Go 侧已落（票 248 系列）；**页面侧的设置页在那棵前端树**（`dsh/feat/frontend-p0-v2`，⛔ dev 自造页面——机主 10-12 源头指令）；`api_key_ref` 机制=D36 规则 5（`dpapi:<blob-id>` 或 `env:NAME`，`internal/config/schema.go:387`）。
4. ⛔ **前置**：面板得先渲染得出来（票 299，`299-r1` 在飞）＋ 出厂链得能把"那棵树的页面产物"嵌进 exe（顺序＝先 299 后这条，见账 `A869` §3）。

## 要建什么（三格，等前置通了再逐格细化）

- [ ] **AC#1 默认选中**：设置页（与首次引导，若有）里 **stepfun 为默认供应商**；⛔ 不当硬编码写死——用户改过就以用户的为准（落点与判据由落地腿按 catalog 现量裁，⛔ 绕"成对写"那条）。
- [ ] **AC#2 界面里能配"两个模型"**：文本角色（`llm.roles.chat` 等）与**语音链**（ASR/TTS/实时那几段的 provider＋model）都能在界面上改；语音那条要能填 stepfun（audio3）的 API 形配置。⛔ 不改 `config.toml` schema 语义。
- [ ] **AC#3 界面里录 key**：走既有 `api_key_ref`（DPAPI blob）机制。⚠ 判据要含：**密钥值⛔ 回显、⛔ 进日志、⛔ 进模型可读文本**（先例＝票 63 那条"key 从不进 argv／日志／对话"）+ 页面上只显示"已配置/未配置"与 blob 名。
- 顺带（⛔ 本票不承诺）：票 312（扫描本地已有配置）落地后，本票的界面要把扫到的候选**预填**进来。

## 禁区

- ⛔ 任何把明文密钥写进 `config.toml` 的形（D36 规则 5 逐字禁）；⛔ 把 key 塞 argv／日志／模型文本（票 63／117／131 那一族既有钉）。
- ⛔ 动 `frontend/index.html` 的 CSP `<meta>`；⛔ 在 dev 那棵树里改页面（页面归 `D:/wt/fe`）；⛔ 动 `go.mod`。
- ⛔ 为"让默认值生效"去放宽任何既有断言（`internal/config` 那几枚成对/可知性用例照旧）。

## 排程

前置＝票 299（面板渲染）⇒ 出厂链嵌新页面那条路 ⇒ 本票。当前**按住**（仅立票，⛔ 派腿）。
- 2026-10-13 313-a1 只读普查交付（HEAD 0e5253fe）：`.scratch/wisp/probes/313/a1/` 四表——`a-get.md`（config.get 现读：provider 行来自 config.toml 已声明行、全新首建=零行、stepfun 要出行须手写该节；模型清单只有计数，"面板无字段列模型"判真）·`b-set.md`（白名单 7 枚全在 [llm]、语音段零枚；新字段候选＋校验逐字（票面 :95-113 已漂 :99-115）＋"不能建行"结构缺口）·`c-default.md`（role 未配=回退 text_chain 首元素；"默认选中"零承载物；候选①界面预选／②首建写入／③tag 或运行时兜底＋各自代价；stepfun 模型名仓内零枚）·`d-page.md`（页面侧现状＋15 行待办规格＋AC#1-3 判据草案：Go 可自测/必须页面/必须真机各归位）。⛔ 未动产码。
