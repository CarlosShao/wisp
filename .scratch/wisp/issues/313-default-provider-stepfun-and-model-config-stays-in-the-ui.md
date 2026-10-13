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
- 2026-10-13 313-r1（界面配模型·Go 侧）交件（HEAD 39f6341d）：按 `A877` 四裁定落三格 Go 半——甲读面＝`SettingsView` 追加字段（`Presets(+Known)`／`Configured(+Known)`／`ChatProvider`／`ProviderSettings.Models`／`Voice`）＋`renderSettingsPayload` 一行 JSON 载荷（标记 `SettingsPayloadMarker`，挂 `config.get` 回执**散文之后**，票 248 两句话逐字未动）；乙写面＝`internal/panel/config_write_313.go` 第二张可写表六枚（`voice_cloud_asr_chain`/`voice_cloud_tts_chain`/`voice_realtime_provider`/`voice_realtime_model`/`provider_row_create`/`role_chat_pair`，⛔ 未动票 248 七枚那张被钉住的表）＋`internal/config/settings_313.go` 六枚 setter（显式建行＝`A877` ②，只收预设名、行已存在即拒、崩溃原子；云链每元素过 `checkChainElement`）；丙＝只做读面半（预设清单＋是否已配置），⛔ 未写首建文件/⛔ 未加 default 标签。⛔ 新 C17 方法、⛔ 编模型名、⛔ 动 `cmd/wisp`（精确落点包＝`probes/313/r1/10-cmdwisp-landing-package.md`）、⛔ 改任何既有断言；页面接口规格＝`docs/evidence/s1/313-panel-ask.md`。新增测试 11 枚全绿（panel 6／config 5）；`go build ./...` rc=0、gofumpt 两文件已洗、`internal/panel` 新增红 0 枚（现树 6 枚失败 ⊂ 纯净树基线 8 枚）。未 push。

- [2026-10-13 21:29 +08] agent=编排者 did=313 收尾链两枚裁定落地（台账 A881）：① `role_chat_pair`＝成立（三条件：两键枚举／原子／回执双报）② 第二张可写表＋两表并查＝成立（三条件：248 钉零动／单一 accessor／并集等值钉）next=落地腿待 294-v1 收口后放行（落点包 `probes/313/r1/10-cmdwisp-landing-package.md`；页面规格转交＝机主动作）

- [2026-10-13 23:44 +08] agent=313-r2 did=`cmd/wisp` 精确落点包落地（`panel_config_store.go`：`ReadSettings` 六味＋`modelNames` 排序切片；`ApplySetting` 六 case＋`splitChainValue`；三禁区一字未动）＋新增 `cmd/wisp/panel_config_313_test.go`（空配置→建 stepfun 行→两云链→realtime 两枚→chat pair 端到端，六枚 case 各报自身键路径；A881 ⓒ 并集等值钉——r1 未放本腿补，现读 PASS）。A881 三条件现读齐：ⓐ 第一张表块 `0d87a681`↔HEAD 字节相同＋248 钉未动（本腿 diff 0）；ⓑ 并集单一 accessor；ⓒ 钉在。门禁：build rc=0／gofumpt 空／panel+config 改前=改后失败集（既有 5 红）／cmd/wisp rc=1 具名 255 两枚——本腿落地使 voice 读者行需补 `cmd/wisp/panel_config_store.go`、`config_readers_255.go:103` 的 llm 证据行需从 `:96` 重裁到 `:98`（修面越界，转 313 收尾链/255-r2，push 前必须清）／d22scan clean／长度门 RED=`.scratch/wisp/issues/173-…-done.md` 超长路径（2a8ab9c4 引入，非本腿）。⛔ 未 push。next=编排者收件；255 名册重裁后清账再推。
