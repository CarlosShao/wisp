# 票 145 非实现者验收表 — AC#2 / AC#6 / AC#2b 三格复核（导出树全量，活树零字节）

**验收代理**：`acceptor-ticket145-r1`（独立对抗验收，未参与 `4111ada61` 的实现）
**被验 commit**：`4111ada6155b386c2ba6e0fcbe21b5f1d45f888b`（2026-11-09，`145: land census 2.14 OPTION A ...`）
**基线 commit**：`71ff0fc8786ec61b8d3b551d019575ef66c67361`（`4111ada61` 的祖先；`git merge-base --is-ancestor 71ff0fc87 4111ada61` ⇒ **YES**）
**本机 `date` 实测**：`Mon, Nov 9, 2026 5:50:45 PM`（Asia/Shanghai）；`go version go1.27.1 windows/amd64`
**活树 HEAD（收手时）**：`908c6652cbba895320649c7933716eff8b941cac`

**树与快照**：全部门禁、全部变异、全部探针只在导出树跑，活树零写入。
- 被验树 `.scratch/wisp/probes/145/acc-r1-clone/`（`git clone --shared`，HEAD `4111ada61`）
- 基线树 `.scratch/wisp/probes/145/acc-r1-base/`（HEAD `71ff0fc87`）
- 变异一律 `go test -overlay=<json>`（仓外副本 + `{"Replace":{...}}`），**被验树本身一字未改**；仓内零 worktree、零 checkout、零删除。
- ⚠ 环境注（照 AC#5 的纪律）：克隆树无 `build/`、`third_party/sherpa-onnx/`（两目录 gitignored）⇒ `cmd/wisp` 测试二进制首跑 `exit status 0xc0000135`（DLL 缺）。每次跑前 `unset GOTMPDIR TMPDIR TEMP TMP` 并把主树两目录绝对路径加进 `PATH`；**全程未设 `CGO_ENABLED=0`**。

**证据档位图例**：〔独立复现〕＝本代理自己在机器上跑出来／写探针复算；〔日志＋归档，我抽验〕＝读了断言原文并抽查；〔仅自述，不背书〕＝只有实现腿自述。

---

## 总判

**两格通过、一格不通过（具名退回 AC#2b）。**

| 格 | 判定 | 依据（一句） |
|---|---|---|
| **AC#2** | **通过，可翻** | 落地集非空＝**6 枚新键**（§1 现跑抽出）；构造点**具名、共 3 处生产 + 2 处装配 + 2 处喂线**（§3）；TS 侧"应当长这样"逐键清单在 §4，`frontend/**` 零字节 |
| **AC#6** | **通过，可翻** | 6 枚新键**逐枚有具名用例**，且每枚都被一发变异反向钉住（§2 表）；**无任何一枚是靠填常量"通"的** ⇒ AC#6 的"装饰品"判据一枚未触发 |
| **AC#2b** | **不通过，退回** | 两处**实测**缺陷（§5）：①「两模型／两档位假树逐字一致」这条判据**全仓无任何用例满足**（现存的"两模型树"第二枚恒 `enabled=false`，断言一律 `len==1`）；②AC#2b② 点名的第一个决定者 `Capabilities.Thinking` **零读者**，实测"档位维在不在"由声明词表决定（`thinking=false` ＋ 两档 ⇒ 快照照送两档） |

按规矩：**不写"通过（附条件）"**；AC#2b 附数字判 FAIL，其余逐格见下。

---

## 0. 现跑提取（⛔ 抄票面、⛔ 抄实现腿证据件）

六个新键全部由本程从生产码现跑抽出，用一枚临时反射探针（`internal/panel/acc_r1_keys_probe_test.go`，**只在导出树**）读生产结构体的 json tag：

| 尺 | 读数 |
|---|---|
| 反射探针 `TestAccR1DumpCarrierJSONKeys`（`4111ada61`） | `RunSection => [tools toolsDropped? running? usage? reasoning? error? stuck? done?]`<br>`ToolRunRow => [callId name outcome resultText? argsJson? level?]`<br>`ModelOption => [provider model display?]`<br>`ModelSwitchSection => [current currentThinking models thinkingLevels reason? adapter? adapterThinkingLevels?]` |
| 同探针 `Snapshot` 顶层键（零值） | `[composer generatedAt pending results]` ⇒ **本票 6 枚新键一枚都不在顶层**，全在既有 section 内 |
| 基线树 `awk` 抽 json tag（`71ff0fc87`） | `RunSection`＝7 键（无 `done`）／`ToolRunRow`＝5 键（无 `level`）／`ModelOption`＝2 键（无 `display`）／`ModelSwitchSection`＝5 键（无 `adapter`/`adapterThinkingLevels`） |
| 两树 `Snapshot` 顶层键对比 | 均 11 枚、**逐键同名**（`pending results composer generatedAt instructions tasks commands modelSwitch ballStyle run plugins`）⇒ 与"新键全在 section 内"互证 |

⇒ **AC#6 的对象集＝6 枚**（`run.done.stop`／`tools[].level`／`modelSwitch.models[].display`／`modelSwitch.current.display`／`modelSwitch.adapter`／`modelSwitch.adapterThinkingLevels`）。

⚠ 本节与 `docs/evidence/s1/145-snapshot-field-census-r1.md` §0 的读数**不同**（那份读的是 09-25 的四字段世界，见 §6 更正）。

---

## 1. 六枚新键的载体位置（`4111ada61` 现量，行号可复跑）

| 键 | 生产定义 | 值从哪来（生产链路） | 上喂线处 |
|---|---|---|---|
| `run.done.stop` | `internal/panel/run_feed_145.go:152` `Stop string json:"stop,omitempty"`；容器 `:184` `Done *RunDoneRow json:"done,omitempty"` | `internal/agent/loop.go:1134` `Stop: stop`（`finish` 的末参）→ `cmd/wisp/run.go:2505` `c.runFeed.TaskEnded(string(e.Stop))` | `cmd/wisp/run.go:2505` |
| `tools[].level` | `internal/panel/run_feed_145.go:99` `Level string json:"level,omitempty"` | `cmd/wisp/run.go:2401` `bookedForensics(e)` 从**同一枚 booked 行**取 `RiskLevel` → `cmd/wisp/run.go:2458` `ToolEnded(..., level)` | `cmd/wisp/run.go:2458` |
| `modelSwitch.models[].display` | `internal/panel/model_switch_187.go:51` `Display string json:"display,omitempty"` | `internal/config/settings_187.go:220` `ChatModelOption{... Display: spec.Display}` ← `ModelSpec.Display` | `cmd/wisp/run.go:531`（`modelSwitchSection`）＋ `panel_resident_pump.go:179` |
| `modelSwitch.current.display` | 同上一枚（`ModelOption` 同一个字段） | `internal/config/settings_187.go:215` `out.CurrentDisplay = spec.Display`（**同一次 walk**） | `cmd/wisp/run.go:531`＋`panel_resident_pump.go:182` |
| `modelSwitch.adapter` | `internal/panel/model_switch_187.go:84` `Adapter string json:"adapter,omitempty"` | `cmd/wisp/run.go:531` `prov.Info().Protocol`（live provider 自报） | `cmd/wisp/run.go:1116` 的 reader 闭包 |
| `modelSwitch.adapterThinkingLevels` | `internal/panel/model_switch_187.go:94` | `internal/llm/thinking_levels.go:57` `ThinkingLevels(protocol)`；三家适配器 `init` **从自家映射表派生注册**（`openairesponses/adapter.go:81` ← `request.go:63` 的 `effortLevels`；`anthropic/adapter.go:106`；`openaichat/adapter.go:47` 注册 `nil`＝已知协议零档） | 同上 |

---

## 2. AC#6 逐字段表（6 枚 × 三列：具名用例 / 变异读数 / 探针读数）

**「载波跳」**＝值从它的直接来源进到线上那一跳（本票 panel 侧拥有）；**「起源跳」**＝值的源头（provider／配置目录）到中间载体那一跳。

| 键 | 具名用例（断言"值来自真来源"） | 变异读数（反向钉） | 判定 |
|---|---|---|---|
| `run.done.stop` | `Test145RunningUsageAndDoneTravelTheWire`（`cmd/wisp/panel_run_145_test.go:152`；:204 断言 `run.Done.Stop == string(llm.StopCancelled)`，:207 字节钉 `"stop":"cancelled"`）＋`TestRunFeedDoneRowFollowsTheEvent`（`internal/panel/run_feed_145_test.go`，字节钉 `"done":{"stop":"max_tokens"}` 与空形 `"done":{}`） | **`acc-r1-mut2`**：`TaskEnded(string(e.Stop))` → `TaskEnded("")` ⇒ 该用例 **FAIL** | **载波跳有；起源跳无**（见 §2.1） |
| `tools[].level` | `Test145ArgsJSONComesFromTheBookedForensicsRow`（`cmd/wisp/panel_run_145_test.go:217`）：**真** `memory.Open` ＋ 真 `st.InsertToolCall{RiskLevel:"L1", CorrelationID:"task-145#call-1"}` ＋ 真 join，断言字节 `"level":"L1"`；两臂突变（无 store／无 booked 行）断言键**缺席** | **`mut-level`**：`level` → `"L0"` ⇒ 该用例 **FAIL** | **两跳都在真源上** |
| `modelSwitch.models[].display` | `Test145ChatSwitchOptionsCarriesTheCatalogDisplayNames`（`internal/config/settings_145_ac2b_test.go:45`，真配置文件 → 真 `config.NewManager`，断言 `Models[0].Display=="DeepSeek Chat"`）＋`Test145ModelSwitchMappingCarriesTheDisplayNames`（`cmd/wisp/panel_model_switch_145_test.go:23`，断言线上字节 `"display":"DeepSeek Chat"`） | **`mut-display`**：`Models[]` 摘掉 `Display` ⇒ 两枚 **FAIL** | **两跳都在真源上** |
| `modelSwitch.current.display` | 同上一对用例（分别断言 `CurrentDisplay` 与 `Current.Display`） | **`mut-display2`**：`out.CurrentDisplay = ""` ⇒ 两枚 **FAIL** | **两跳都在真源上** |
| `modelSwitch.adapter` | `Test145ModelSwitchMappingCarriesTheAdapterVocabulary`（`cmd/wisp/panel_model_switch_145_test.go:69`，真 `llm.NewProvider(llm.EndpointOptions{Protocol:"openai-responses"})`，断言 `Adapter=="openai-responses"`＋字节钉） | **`mut-adapter2`**：`adapter = ""` ⇒ 该用例 **FAIL** | **两跳都在真源上** |
| `modelSwitch.adapterThinkingLevels` | 同上（断言 `[]string{"high","low","medium"}`）＋`Test145ThinkingLevelsComeFromTheAdaptersOwnTables`（`internal/llm/thinking_levels_145_test.go:20`，断言词表**就是**适配器自家映射表；`:34` 钉 `openai-chat` 已知协议零档；`:47` 返回副本；`:58` 重复注册 panic） | **`mut-adapter`**：`adapterLevels = nil` ⇒ 该用例 **FAIL** | **两跳都在真源上** |

**结论：6 枚新键一枚不落地为"装饰品"——每枚都有具名用例，且每枚都被一发只改该字段的变异反向钉红。** AC#6 的补救动作（"从落地集里拿掉"）**一枚未触发**。

### 2.1 唯一具名缺口：`run.done.stop` 的**起源跳**全仓无断言（本程实测，非推断）

本程做了三件互相独立的事来把这件事从"推断"变成"读数"：

1. **全量名册相等**〔独立复现〕：`acc-r1-mut`（把 `loop.go:1134` 的 `Stop: stop` 改成 `Stop: ""`，即只掐起源跳）下跑**整个** `./cmd/wisp`，与干净树跑同一命令的 `--- FAIL` 名册逐名 `diff` ⇒
   `clean=63 mut=63`、**`IDENTICAL: mutating the loop->Event.Stop hop changes ZERO cmd/wisp test outcomes`**。
   （日志：`probes/145/acc-r1-clean-full2.txt`、`probes/145/acc-r1-mut-full.txt`、名册 `probes/145/acc-r1-clean-cmdwisp-fails.txt`。）
2. **同轮阳性对照**〔独立复现〕：同一棵被验树上别的变异都真红——`acc-r1-ctl`（改 `res.Status`）⇒ **15 FAIL**；`acc-r1-mut4`（`Kind: EvDone`→`EvWarn`）⇒ **5 FAIL**；`acc-r1-mut3`（填 `"BOGUS_XYZ"`）⇒ 机制活。⇒ **"零差集"不是 overlay 没生效**。
3. **可断言性证明**〔独立复现〕：本程在导出树写了一枚 6 行的临时探针 `TestAccR1TerminalEvDoneCarriesTheProvidersOwnStopWord`（复用仓里既有的 golden harness：`newHarness(t,"text-reply")` ＋ `h.run(...)`，断言 `h.sink.LastOf(EvDone).Stop == llm.StopEndTurn`）：
   - 干净树 ⇒ **PASS**（载体诚实）；
   - `acc-r1-mut` overlay 下 ⇒ **FAIL**：`terminal EvDone.Stop = "", want "end_turn" (the replayed wire's own end_turn)`。

   ⇒ 这条缺口**不是"测不了"，是"没写"**。写起来就是在既有 harness 上加 6 行。

**这不是票 145 的缺陷，是全仓形状**〔独立复现〕：本程另做两发标定，掐**既有已验收字段**的起源跳——`cal-text`（`EvDone.Text` → `"CALIBRATION_BOGUS"`）与 `cal-usage`（`EvDone` 的 token 数 → 0）——`./internal/agent` 均 **ok**。
⇒ **仓里的 agent 侧用例断言的是"载体忠实搬运事件说的值"，从不_"事件说的值来自 provider"_。** 这条形状是既有验收（票 30／370 那一族）就已接受的，本票与它同形、未加深。
**登记为 `R-145-1`（不阻断本格翻勾，见 §7 处置建议）**。

---

## 3. AC#2 的构造点（谁、几处——AC#2 原文要求"必须点名"）

**生产构造点＝3 处，全部在 `internal/panel/pump.go`**〔独立复现，尺：`grep -rn "NewSnapshot(" --include=*.go .`，排除 `_test.go`〕：

| # | 位置 | 作用 |
|---|---|---|
| 1 | `internal/panel/pump.go:409` `snap := NewSnapshot(cards, results, composer, now())` | **唯一的真值装配点**：卡片、结果、composer、注入时钟 |
| 2 | `internal/panel/pump.go:479` `snap.Run = p.src.Run()` | `run` 段（含 `done.stop`、`tools[].level`）**reader-gated**：有 reader 才上键 |
| 3 | `internal/panel/pump.go:460` `snap.ModelSwitch = &sec`（section 由 `:453` 的 `p.src.ModelSwitch()` 取） | `modelSwitch` 段（含 4 枚新键）**reader-gated**，且 `:456/:459` 把 `nil` 归成 `[]`（picker 永不遇 `null`） |

**装配根（reader 从哪来）＝2 处**：
- `cmd/wisp/run.go:1104` `Run: rt.runFeed.Section`（管线形状）
- `cmd/wisp/run.go:1116` `ModelSwitch: func() panel.ModelSwitchSection { return modelSwitchSection(mgr, chatChainProvider(rt)) }`

**喂线（值怎么进 feed）＝2 处**：`cmd/wisp/run.go:2505`（`done.stop`）、`cmd/wisp/run.go:2458`（`tools[].level`）。

**测试构造点＝13 处**（`composer_test.go:223/:251`、`display_scrub_216_test.go:47/:271`、`instructions_200_test.go:188/:294`、`rewrite_account_102r2_test.go:34/:52/:69`、`panel_clamp_33r14_test.go:87`、`projctx_test.go:357/:378`）——本程只用来佐证"生产侧不止零处"，不作判据。

⚠ **这一节直接推翻 `145-snapshot-field-census-r1.md` §4.1 的"`panel.NewSnapshot` 非测试调用者＝0"**。见 §6。

**AC#2 的 TS 侧要求（"不写 `frontend/**`，只在证据件里给出应当长这样的逐键清单"）**：满足，清单在 §4；`frontend/**` 零字节（§8 有尺）。

---

## 4. TS 侧"应当长这样"逐键清单（本程现量 `frontend/src/lib/panel.ts` 后写的对齐件，⛔ 本程未写 `frontend/**`）

现量（`4111ada61` 导出树，`frontend/src/lib/panel.ts` 共 761 行）：4 个相关 interface 各缺键，**且没有一个尺会因此变红**（见 §4.1）。

```ts
// :283 —— 现为 { provider: string; model: string }
export interface ModelOptionView {
  provider: string;
  model: string;
  /** 人面向名，来自 ModelSpec.Display；未声明则缺席（omitempty），
      页面回退显示 model id，不许自造名。 */
  display?: string;
}

// :294 —— 现为 { current, currentThinking, models, thinkingLevels, reason? }
export interface ModelSwitchSection {
  current: ModelOptionView;
  currentThinking: string;
  models: ModelOptionView[];
  thinkingLevels: string[];
  reason?: string;
  /** live provider 自报的协议词（Info().Protocol）；无 provider 链时缺席。 */
  adapter?: string;
  /** 该适配器实际接受的档位词表；与 thinkingLevels 是两问：
      thinkingLevels = 当前模型声明的档，adapterThinkingLevels = 这家协议收的档。
      页面若要"只显示真能发出去的档"，须自行取交集——Go 侧今天不算这个交集。 */
  adapterThinkingLevels?: string[];
}

// :339 —— 现为 { callId, name, outcome, resultText? }
export interface ToolRunRowView {
  callId: string;
  name: string;
  outcome: string;
  resultText?: string;
  /** sanitize 后的入参 JSON，来自 booked forensics 行；无 booked 行则缺席。 */
  argsJson?: string;
  /** 该次调用被门按 L0/L1/L2 哪档放行/拒绝，与 argsJson 同一枚 booked 行；无行则缺席。 */
  level?: string;
}

// :365 —— 现为 { tools, toolsDropped?, reasoning?, error?, stuck? }
export interface RunDoneRowView {
  /** provider 自己的终态词（end_turn / max_tokens / cancelled / error…）；
      loop 自有终态（C22 刹车、预算拒、buildRequest 失败）没有 provider 词 ⇒ 缺席，
      不是空串。 */
  stop?: string;
}

export interface ToolRunningRowView { callId: string; name: string; }

export interface RunUsageRowView {
  tokensIn: number;
  tokensOut: number;
  cachedTokens?: number;
  micros: number;
  currency: string;   // 与 micros 成对，永不同时缺席
}

export interface RunSection {
  tools: ToolRunRowView[];
  toolsDropped?: number;
  reasoning?: string;
  error?: RunErrorRowView | null;
  stuck?: RunStuckRowView | null;
  done?: RunDoneRowView;
  running?: ToolRunningRowView[];   // 非本票落地集（wave4 已落），panel.ts 同样缺
  usage?: RunUsageRowView;          // 同上
}
```

**归属**：本票落地集＝前 6 枚（`display`×2、`adapter`、`adapterThinkingLevels`、`argsJson`、`level`、`done`）；`running`/`usage` 是 wave4 的键，一并列出只为让前端会话一次对齐完，**不计入本票 AC#6 的对象集**。

### 4.1 ⚠ 具名缺口二：这 6 枚新键**没有任何一把尺在盯**

- 合同尺 `TestComposerContractTypesMatchFrontend`（`internal/panel/composer_test.go:48`）的 `pairs` **只配 6 个顶层类型**（`Snapshot`↔`PanelSnapshot`／`ComposerState`／`ModeView`／`WorkspaceView`／`AttachmentRef`／`ResultChunk`）。6 枚新键全在**嵌套** section 内 ⇒ **这把尺对它们完全瞎**。
- 现量佐证〔独立复现〕：`frontend/src/lib/panel.ts` 里 `display`／`adapter`／`adapterThinkingLevels`／`argsJson`／`running`／`usage` 作为**声明键**命中数**均为 0**（`grep -c "\b<key>[?]*:"`），而 `internal/panel` 全量套件在 `4111ada61` 仍是 **PASS 212 / FAIL 2**（两枚在册红与本节无关）。⇒ **Go 侧加了 6 枚键、TS 侧一枚没跟上，零枚用例变红。**
- 这正好回答 AC#4③ 那句"要用一枚用例证明 Go 加了字段而 `panel.ts` 没跟上时，谁会响"：**顶层会响，嵌套不会**。
- **登记为 `R-145-2`**（不阻断本票两格翻勾：AC#2 只要求"给出清单"，不要求造尺；但它是"下一次同类漂移无人拦"的洞）。

---

## 5. AC#2b：不通过（两处实测缺陷，逐条具名）

### 5.1 缺陷①：判据点名的"两模型／两档位假树逐字一致"这条，全仓**无任何用例**满足

AC#2b 逐字要求：**"拿一份两模型／两档位的假配置树喂进去，快照里读到的清单与那棵树逐字一致（与票 92 那格同形），⚠ 不许用'字段非空'充当判据"**。

现量〔独立复现〕：
- 仓里确有两棵"两模型"假树：`internal/config/settings_145_ac2b_test.go` 的 `settings145DisplayBase` 与 `internal/config/settings_187_test.go` 的 `settings187Base`。**但两棵树的第二枚模型都是 `enabled=false`**，而两处断言一律 `len(opts.Models) != 1` ⇒ 断言的是**过滤方向**（未启用的不许进），**不是复数方向**（两枚启用的都要在）。
- 全仓 `len(...Models)` 断言普查：`!= 1`／`!= 0`／`!= 2`；其中两处 `!= 2`（`settings_257_test.go:234`、`firstrun_257_test.go:251`）数的是**某 provider 的 model map**，不是 picker 清单。⇒ **没有任何用例断言过"picker 清单里有两行"。**
- **变异读数**：`mut-onerow`（`settings_187.go:220` 追加 `if len(out.Models) >= 1 { return out }`，即"只肯给第一枚"）⇒ `./internal/config` **ok**、`./cmd/wisp` **ok**。**全仓零红。**
- **阳性对照**：`mut-zerorow`（清空 `Models`）⇒ `Test145ChatSwitchOptionsCarriesTheCatalogDisplayNames` ＋ `Test187ChatSwitchOptionsNeverOffersWhatTheResolverWouldRefuse` **双红** ⇒ overlay 确实生效，`mut-onerow` 的绿是**真空洞**，不是仪器没接上。
- **载体本身是对的**〔独立复现〕：本程写的临时探针 `TestAccR1PluralPickerListAgreesWithATwoModelTree`（2 provider／4 模型／**3 枚 enabled**，逐行断言 `panel.ModelOption{Provider,Model,Display}` 逐字、`Current` 带名、`CurrentThinking=="low"`、`ThinkingLevels==[off,low]`（**当前模型自己的**）、线上第 2 行 `"Kimi K2"`，适配器半枚走真 `llm.NewProvider`）：
  - 干净树 ⇒ **PASS**；
  - `mut-onerow` overlay 下 ⇒ **FAIL**：`offered rows = [{Provider:deepseek Model:deepseek-chat Display:DeepSeek Chat}], want the tree's THREE enabled models`。

  ⇒ **生产码满足 AC#2b 的判据，仓里没有用例钉住它。** 判据是"用例断言"，所以这一格**不能翻**——探针在导出树里，不是仓里的用例。

### 5.2 缺陷②：AC#2b② 点名的第一个决定者 `Capabilities.Thinking` **零读者**

AC#2b② 逐字：**"`Capabilities.Thinking` 决定'有没有这一维' ＋ 该适配器实际接受的档位词表决定'有哪几档'，两问分开答"**。

现量〔独立复现〕：
- `Capabilities.Thinking` 全仓非测试读者**只有 3 处**，无一处与 picker／档位有关：`internal/llm/probe_health.go:105`（`case ProbeThinking: return c.Thinking`，探测面）、`cmd/wisp/run.go:2169`（从 flag 构造结构体）、`tools/mockllm/control.go`。`internal/config/settings_187.go` 的 `ChatSwitchOptions` walk（`:195-225`）**从不读 `Capabilities`**；`internal/config/validate.go:243-259` 的 `validateModelSpec` **也不做** `Capabilities.Thinking` 与 `ThinkingLevels` 的交叉校验（只逐枚验枚举）。
- **两问实测**〔独立复现〕：
  - `TestAccR1ThinkingDimensionGate`（`internal/config`，两棵树）：`capabilities.thinking=false` ＋ `thinking_levels=["off","low"]` ⇒ `ThinkingLevels=[off low]`（len=2）；`capabilities.thinking=true` ＋ 无档位表 ⇒ `ThinkingLevels=[]`（len=0）。⇒ **"有没有这一维"由声明词表决定，位从来没被问过。**
  - `TestAccR1TierListIgnoresBothNamedDeciders`（`cmd/wisp`，真 `config.NewManager` ＋ 真 `llm.NewProvider`）：树声明 `capabilities.thinking=false`、`thinking_levels=["off","high"]`、协议 `openai-chat`（注册为**已知协议、零档**）⇒ 读数
    ```
    adapter protocol      = "openai-chat"
    adapter vocabulary    = [] (len=0)
    offered tier list     = [off high] (len=2)
    READING A: the dimension EXISTS (2 tiers) although capabilities.thinking=false
    READING B: offered tiers the adapter does NOT accept = [off high]
    ```
    ⇒ **两问都被绕过**：维在不在不看位；"有哪几档"也不看适配器词表——两个清单**并列**送到线上（`thinkingLevels` 与 `adapterThinkingLevels`），**交集留给页面算**。
- 结论：AC#2b② 的**形状实现了**（两个数据都在线上、各有真源），**但 AC 写的判据不是它的判据**。这是**规格与代码的分歧**，不是"缺一枚测试"⇒ 本格**不能翻**，须编排者裁：要么改 AC#2b② 的文字（"维在不在由声明词表决定；适配器词表并列送出，交集由页面算"——**改 AC 文字属人工批准**），要么让代码照 AC 写（`Capabilities.Thinking=false` ⇒ 不送该维；`thinkingLevels` 与 `adapterThinkingLevels` 取交集后送）。**本程不选形**。

---

## 6. `145-snapshot-field-census-r1.md` §4.1 的过期更正（本程最大的一处更正）

| 事实 | 尺 | 读数 |
|---|---|---|
| 普查件锚点 | `git log -1 --format=%ad 5d8ffdaa2` | **2026-09-25 14:25:39 +0800** |
| 票 200 给面板载体装上生产者 | `git log -1 --format=%ad a90677e92` | **2026-09-28 20:12:50 +0800** |
| 票 145 `run` 段落地 | `git log -1 --format=%ad 079a3f529` | **2026-11-07 21:01:09 +0800** |
| 两条是否在普查之前 | `git merge-base --is-ancestor a90677e92 5d8ffdaa2` / `… 079a3f529 5d8ffdaa2` | **NO / NO** |
| 今天 `NewSnapshot` 的非测试调用者 | `grep -rn "NewSnapshot(" --include=*.go . \| grep -v _test.go` | **3 处**（`composer.go:201` 定义、`pump.go:320`、`pump.go:409`）⇒ 生产调用者 **2 处** |

⇒ 普查 §4.1 那句"`panel.NewSnapshot` 非测试调用者＝**0**"在**写下那天是对的，今天已过期**：泵现在有生产者、且真的在推——
`cmd/wisp/panel_pump.go:511` `rt.pump.Publish()` → `internal/panel/pump.go:524 Publish()` → `:528 p.Snapshot()`；驱动点 `cmd/wisp/run.go:1170`（`rt.ui.publish = rt.publishPanelSnapshot`）、`cmd/wisp/resident_present_windows.go:197`、`cmd/wisp/resident_windows.go:546`；常驻腿另有 `cmd/wisp/panel_resident_pump.go:150`。
**⇒ 09-26 那句"没有地方能给它真值"的前提已被 09-28 与 11-07 两枚 commit 推翻，本票 AC#2 的落地集因此从"空"变成"非空"。这条必须写进证据，不能继续引普查的旧读数。**

---

## 7. 门禁读数（全部在导出树 `acc-r1-clone`，锚 `4111ada61`；环境照 §开头）

| 尺 | 读数 |
|---|---|
| `gofumpt.exe -l internal/panel internal/agent internal/config internal/llm cmd/wisp` | **空**，rc=0 |
| `go vet`（同五包） | **rc=0** |
| `go build`（同五包） | **rc=0** |
| `go test ./internal/panel` | **PASS 212 / FAIL 2** —— `TestComposerRenderFixtureTellsTheTruth`（克隆树缺 fixture，环境性，wave3/ac2r1 同读数）、`TestC21DesignTokensFourWayAgree`（在册）。基线树 `acc-r1-base` 同尺：**PASS 210 / FAIL 2**，**逐名同两枚** ⇒ **零新红** |
| `go test ./cmd/wisp`（干净树全量，`-count=1`） | **63 枚具名 FAIL**，全部落在既有真窗／语音／名册族（`Test12r2…`／`Test15e…`／`Test16r1…`／`Test26P/R3…`／`Test313…`／`Test381r2…`／`Test382WR2…`／`Test390r1…`／`Test392…`／`Test394R1…`／`Test395r1…`／`Test41a2r1…`／`TestAC1…`／`TestAC228…`／`TestAC246…`／`TestAC247…`／`TestAC2RealProcess…`／`TestAC3Early…`／`TestAC4Focus…`／`TestLockedSuffix…`／`TestSecret…`／`TestTicket198/224/255×2/256…`／`TestWireR1…`×6）。**`Test145*`／`Test187*` 零 FAIL**；基线树 `71ff0fc87` 同族名册 ⇒ **本票 commit 未新增 cmd/wisp 红** |
| 定向复跑（被验的 6 枚新键相关用例） | `Test145*`（8 枚）＋`Test187*`（5 枚）**全 PASS** |
| `sh scripts/d22scan.sh` | **rc=0 clean** |
| `sh scripts/check-path-length-budget.sh` | **GREEN**，rc=0（63 条超预算路径在册） |

---

## 8. 边界与如实登记

- **活树零写入**：本程所有跑动、变异、探针都在 `.scratch/wisp/probes/145/` 的两棵克隆树与仓外 overlay 副本里。本程只改两处：本证据件 ＋ 票 145 的三枚 AC 框（翻两枚）＋ 票面 Progress log 一行。
- **禁面零字节**〔独立复现〕：`git show --name-status 4111ada61` 共 **20 文件**；`grep -cE "^(frontend/|forks/|internal/winsec/|docs/PLAN.md|docs/specs/|thresholds|.*golden)"` ⇒ **0**。
- **被验树＝活树（对全部被验路径）**〔独立复现〕：`4111ada61..908c6652c` 共 **12** 枚 commit，其中触及 `internal/panel` 的只有 `908c6652c`（33-r3，两枚**测试件**＋`panel_host/resident_windows`）；对 6 枚新键的定义行与读者行逐枚 `git diff` ⇒ **空**（`run_feed_145.go`／`model_switch_187.go`／`composer.go`／`pump.go`／`loop.go`／`settings_187.go`／`schema.go`／`validate.go`／`thinking_levels.go`／`run.go`／`panel_pump.go`／`panel.ts` 全部无 diff）。
- **临时件只建不删**：`probes/145/` 下新增 5 枚探针件（`acc_r1_keys_probe_test.go`／`acc_r1_plural_probe_test.go`／`acc_r1_stop_probe_test.go`／`acc_r1_capability_probe_test.go`／`acc_r1_tier_probe_test.go`，**全部只在克隆树、未跟踪、永不入仓**）＋ 4 枚日志（`acc-r1-clean-cmdwisp-fails.txt`／`acc-r1-clean-full2.txt`／`acc-r1-mut-full.txt`／`acc-r1-probes-final.txt`／`acc-r1-probes-mutated-final.txt`）。零删除。
- **本程没做的**：①没有选形——§5.2 的"改 AC 文字 vs 改代码"、§2.1 的"补起源跳用例"、§4.1 的"给嵌套 section 造尺"三项**一律交编排者**；②没有动 `frontend/**`、`design/**`、thresholds/golden、`docs/PLAN.md`、`docs/specs/**`；③没有重跑票 145 之外的门（如 `-race`、`scripts/wisp-cli-tests.sh` 全跑）。
- **两条登记（不阻断翻勾，建议另立）**：
  - **`R-145-1`**：agent 侧"事件的值来自 provider"这一跳全仓无用例（`cal-text`／`cal-usage` 两发标定同色）；`run.done.stop` 的起源跳是它的一个实例。**可断言、已证明**（本程探针 6 行即红），建议单独开一枚"给 golden harness 加终态词断言"的小票，别挂在 145 上。
  - **`R-145-2`**：合同尺只配顶层类型，嵌套 section（`RunSection`／`ModelSwitchSection`／`ToolRunRow`／`ModelOption`）**无尺**⇒ Go 加了 6 枚键、`panel.ts` 一枚没跟上，零红。建议把 `pairs` 扩成"含嵌套 section"的树形对账。
