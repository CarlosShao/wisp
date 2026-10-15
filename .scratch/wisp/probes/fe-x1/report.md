# fe-x1 只读交叉核对 —— 前端模型页 v3 接线 vs §7/§7.1 口径 vs Go 侧硬耦合

- 腿：`fe-x1`（只读，零产码）。锚 `bf25809a`（HEAD 现读）；前端件＝提交 `24509d30`（`git status -- frontend/` 现读：7 个源文件与提交一致；工作树另有一处未提交改动 `frontend/src/styles/tokens.generated.css`，与本次七件无关）。
- 读数方法：逐行现读＋`grep`；并**1:1 复刻 Go 仪器的 regex/扩展名集/`codeOnly` 做 python 纯解析**（⛔ 未跑 Go/npm 构建，派单禁）；凡"应红/应绿"结论都给出定音命令。

## 0. 结论先行

**一条不一致（唯一）**：

- **门二（`composerModeWriteRe` 的 `\bsetMode\b`）命中本次新增的两行**：`frontend/src/components/model-settings.tsx:693`（`const [mode, setMode] = useState<(typeof MODES)[number]>("内置供应商");`）与 `:765`（`<SegmentedControl onChange={setMode} …`）。
  ⇒ `TestPlantedComposerModeWriteGoesRed` 的 real-tree 半（`internal/panel/composer_test.go:141-144`：`scanComposerPermissionWrites(frontend)`＋`modeWriteMatchesInCodeFile` `:804-817`）**在 HEAD 现态即红**，不需等 320-r2；该枚在 CI portable core scope 里跑（实证：`.scratch/ci-logs/run-37158259050.log:10448`）。定音命令（本腿不能跑）：`go test ./internal/panel -run TestPlantedComposerModeWriteGoesRed`。
- 其余五条**全一致**；320-r2 要的那枚硬耦合（页上必须有 `"config.models.list"` 字面量）**已满足**。
- **给 320-r2 的判语：不能安全起跑** —— 共享编译面（`internal/panel`）现态带一枚与本票无关的红；红因不在 320-r2 写面，也不能靠放宽断言解；`frontend/**` 归机主外部 agent（本编队零写面）。改法不在本腿职责内（⛔ 本腿不给改法）。

## 1. 方法名逐字（读：一致）

- 字面量唯一发送点＝`frontend/src/lib/panel.ts:409`：`method: "config.models.list",`（大小写/点号逐字；纯 ASCII）。
- 同名的其余 7 处全部是注释（不算发送点）：`config-actions.ts:53`、`config-view.ts:88/92/98`、`model-settings.tsx:7/63`、`panel.ts:401`。
- Go 侧现读**尚无名**：`grep -rn "config\.models\.list" internal cmd --include=*.go`＝0 命中；`internal/panel/bridge.go:41-68` 现有 6 枚常量（`panel.*` 四枚＋`config.get`/`config.set`）——与 §7「320-r2 才注册」一致。
- 未注册期间点按钮的行为＝预期：方法名过不了 `ParseComposerRequest`（`bridge.go:132-134` 拒绝语「不是面板 composer 通路的能力入口」），前端 `config-feedback.ts:96` 有对应人话「暂不支持拉取模型列表」。

## 2. 信封形状（读：一致）

- `requestProviderModels`（`panel.ts:407-414`）走既有 `requestHostAnswer`（`:348-354`，唯一门 `window.wispDispatch` `:350`）：`method` `:409`＋`requestId: nextRequestId()` `:410`＋`source: "panel-composer"` `:411`＋**provider 键＝`configProvider`** `:412`。与 §7.1 给的那一行逐字段一致。
- 未新开通道：`postMessage` 调用点仍恰 2（`panel.ts:186`、`:225`，复刻读数）；`grep "fetch(\|new WebSocket\|XMLHttpRequest" frontend/src`＝0 命中。Go 侧对应钉：`bridge_test.go:161-164`（计数 2＋禁 fetch/WebSocket）、`composer_test.go:438-441`＋`:513`。

## 3. 回执解析（读：一致；容忍逐条点名）

- 标记行逐字：`config-view.ts:90`＝`MODELS_LIST_MARKER = "wisp-models-list:"`；**逐码点核过**（17 字符全 ASCII，末位 `0x3a`＝半角冒号）＝§7 口径逐字一致。
- `parseModelsListReply`（`config-view.ts:102-125`）：`indexOf` 定位标记 → 标记后第一行 → `JSON.parse` → **只读 `models[].id`**（string 且 trim 非空才收，push trim 后的值）；`provider` 仅 string 时取、否则 `""`；census §2 例里的 `ownedBy` 走「未知字段忽略」（`:98-101` 注释明示）——与 §7＋census 例一致。
- **自带容忍逐条**：⒜ 标记不钉行首（散文任意位置出现即认）；⒝ 只取标记后第一行（`afterMarker.split("\n", 1)[0]`）；⒞ 每个 id 去两侧空白、空白 id 跳过；⒟ **无**小写化/大小写归一，**无** id 去重；⒠ `models` 缺失或非数组 ⇒ `{kind:"parsed", models:[]}`（页面读成「上游返回的清单是空的」，**不是**错误态）；⒡ 不校验回执 `provider` 是否等于请求的 provider；⒢ **无 `payloadError` 分支**（config.get 的 `parseConfigReply` `:150-153` 有）——census §2 口径是「失败＝散文、无标记行」，故属口径外；若 320-r2 用 `{"payloadError":…}` 形态回执，会被读成 parsed 空表（事实登记，非不一致）。
- 失败路径：无标记/JSON 坏/非对象 ⇒ `no-payload`，散文原样交 `translateFetchReply`（`config-feedback.ts:149-158`）就地显示；无宿主 ⇒ `null` ⇒「未连接到本机服务，未拉取」。

## 4. 三个写入接线（读：一致，逐字）

- 计费：`model-settings.tsx:424-428` → `performWrite({ field: "provider_billing", provider: row.name, value })`；值域字面量 `"plan" | "pay-per-token"`（按钮 `:545-550` 套餐/按量）。逐字对 §7。
- 自定义建行：`:737-741` → `field: "provider_row_create_custom"`、`value: `${protocol},${customBase.trim()}``；协议词表恰三枚（Picker `:804-808`＝`openai-chat`/`openai-responses`/`anthropic`）。逐字对 §7；预设名本地即拦（`:713-715`）。
- 加模型：`:437`（手动）与 `:457`（拉取后逐枚，一枚一次）→ `field: "model_row_create"`＋`provider`＋`model`，**无 value**（`requestConfigWrite` 只在 `value !== undefined` 时上 `configValue`，`panel.ts:384`）。逐字对 §7。
- 其余既有字段同形未动：`provider_base_url`（`:491`）、`role_chat_pair`（`:907`）、`provider_row_create`（`:732`，预设名＋model 选填）。
- 前缀/信封键：全树 `"panel.` 字面量仅 5 处（四枚白名单＋`panel.approval.request` `:188`，均既有）；新代码**零 `panel.` 前缀、零新信封键**（只用 `configField/configProvider/configModel/configValue/credentialValue`，`panel.ts:376-398`）。

## 5. 零 emoji（逐文件：全零命中）

- 扫的禁带区＝仪器 `tools/d22scan/main.go:186` 的 `emojiRe` 同集（`1F000-1FAFF`／`2200-22FF`／`2600-27BF`／`2B00-2BFF`／`FE0F`／`1F1E6-1F1FF`）。
- 本次七件：`model-settings.tsx`、`voice-settings.tsx`、`settings-v2-preview.tsx`、`config-actions.ts`、`config-feedback.ts`、`config-view.ts`、`panel.ts` —— **7/7 零命中**（注释也在射程内，一并扫过）；附带 `tokens.generated.css`（未提交件）也零命中。
- 另用 `TestFrontendHasNoEmoji`（`frontend_hygiene_test.go:246-279`）的文件集复核（`frontend/src`＋`frontend/scripts`＋`index.html`，102 件）：**零命中**。
- 两张截图落 `docs/reports/313-model-page-v2-shots/*.png`——不在 `frontend/` 面（d22scan emoji scope `main.go:582` 只含 design/frontend/internal/cmd），不触发该面。
- ban#6 同扫：`approval\.decide`（`main.go:171`）全树 0 命中；composer 门二形状之外的 `setMode` 命中见 §0/§6。

## 6. 会给 320-r2 的风险（逐枚点名）

- `TestEveryWhitelistedInboundNameIsAlsoSentByThePage`（`bridge_test.go:184`；helper `namesAbsentFromPage:230-238` 要求 `"<name>"` 出现在 panel.ts）⇒ 页上已有 `"config.models.list"`（`:409`）⇒ **绿**（320-r2 注册新常量后此枚不会红）。
- `TestTheRendererHoldsExactlyOneDoorToTheHost`（`composer_test.go:572-594`）复刻读数：2 调用点／0 运行时点号拼装／0 计算式 `sendRequest`／0 未知 `panel.*` ⇒ **绿**。
- `TestFrontendComposerRequestsMatchTheEnvelope`（`bridge_test.go:122-165`）⇒ 四名＋source＋requestId＋计数 2 ⇒ **绿**。
- `TestTheBoundDispatchDoorIsReachableFromThePage`（`bound_dispatch_door_114_test.go:438`）⇒ 新增代码零 `window.*` 读（`window.wispDispatch` 仍只在 `panel.ts:350`）⇒ **绿**。
- `frontend_hygiene_test.go` 四枚（stateless／颜色字面量／emoji／approval.decide）复刻：新件零命中 ⇒ **绿**。
- **`TestPlantedComposerModeWriteGoesRed`（`composer_test.go:141-144`）⇒ 红**，逐字点名 `frontend/src/components/model-settings.tsx:693`、`:765`（判据链：`scanComposerPermissionWrites` `:771-800` 走 `rendererSourceFile` 扩集 `:429-436` → `modeWriteMatchesInCodeFile` `:804-817` → `codeOnly` `:824-842` → `composerModeWriteRe` `:92-95` 含 `\bsetMode\b`；`:693/:765` 无注释前缀，`codeOnly` 不剥）。
- **判语**：**320-r2 不能安全起跑**（前置＝上述门二命中消解；消解后由编排者复跑该枚转绿再发）。

## 具名没做

- 没跑任何 Go/npm 命令（无 `go test` / `go build` / `npm`）；上述"应红/应绿"均为 regex 复刻读数，定音命令已在 §0 给出。
- 没改任何文件（含 `frontend/**`、`design/**`）、没碰票面/台账/HANDOVER、没动 `probes/orch/20261015/` 的派单档。
- 没核 320-r2 的 Go 侧注册四点（`bridge.go`/`composer_dispatch.go`/`config_handlers.go`/`inbound_roster_253_test.go`）——不属本腿面。
- 没跑 d22scan 二进制（只按其 `main.go` 的 `emojiRe`/ban#6 规则复刻扫）。

## 你需要做什么

1. 处理 §0 那一条门二命中（红在 HEAD 现态、CI portable core 会跑；改法归机主前端 agent，本编队 `frontend/**` 零写面）。
2. 在该命中消解前别发 `320-r2`：发出去收件时 `go test ./internal/panel` 必带这枚红，「零新增红」收口不成立。
3. 若要把本读数用作台账输入：署名 `fe-x1`＋本报告路径；本腿未动任何票面/台账。
