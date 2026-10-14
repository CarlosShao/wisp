# 票 313 验收腿 `313-v1` 判语（非实现者；面＝internal/panel＋internal/config＋cmd/wisp）

- **锚**：`bfb62488`（2026-10-14 10:07 +08 终跑）；被裁 `89b3b550`／`2f8964b9`＋`9e3013e4`；页面侧 `9e9964f8` 只读。
- **载具**：`.scratch/wisp/probes/313/v1/`（overlay 探针＝真 Manager/Store/ComposerDispatch；keys 对拉脚本；logs 两枚；锚与门禁 `00-anchor-and-gates.md`）。⛔ 未只引实现者的件：载荷、突变、失败名单全部本腿现跑。

## 0. 结论先行

| # | 格 | 判语 |
|---|---|---|
| 1 | 读面载荷 三方对拉 | **成立**（键集 11/7/5 三方全等）＋**1 枚形状分歧具名**（零行 `providers:null`，见 §1.3） |
| 2 | A881 ⓐ/ⓑ/ⓒ | **ⓐ 成立**；**ⓑ 成立**；**ⓒ 成立**——两发突变证牙，牙的位置具名（§2） |
| 3 | 六枚 setter 拒绝面/原子性 | **成立**：十发拒绝文件 sha256 不变＋合法控制移动 hash（§3） |
| 4 | cmd/wisp 端到端 | **成立**；行为突变甲＝钉红（有牙）／乙＝三包零红（**一处覆盖缺口具名**，§4） |
| 5 | 密钥面 | **成立**：并集拒绝文案现读；canary 四表面零命中（亲跑）；248 哨兵钉 PASS（§5） |
| 6 | 失败集归因 | **0 新红**；5 枚旧红全在名册三族，逐枚归因（§6） |
| 7 | 翻框建议 | AC#1/#2/#3 的 **Go 半＝成立**；页面半不判；`AC#3` 机制＝**既有且未被动**（§7） |

## 1. 读面载荷：Go 现跑 ↔ 规格 ↔ 页面

1.1 现跑（探针四发，真装配）：空配置与"三步写＋云链后"两枚载荷逐字入 `logs/10-run-and-keys.txt`。
1.2 机械对拉（`compare-keys-313v1.py`）：顶层 `["chatModel","chatModelKnown","chatProvider","configured","configuredKnown","credential","presets","presetsKnown","providers","readable","voice"]`、行七枚、voice 五枚——**go-vs-spec、go-vs-page、spec-vs-page 三方向差集全空**；页面读键机械提取也无"页面读 Go 不发"的键。写方向同拉：页面发的 7 枚字段（3 步＋云链/实时 4 枚）⊆ 并集 13 枚，page-only=0。
1.3 ★**具名分歧（形状、非键）**：零 provider 行时载荷逐字 `"providers":null`（Go nil 切片），而 `313-panel-ask.md` §1 的 `providers[]` 未标 null、`frontend/src/lib/config-view.ts:80` 声明 `providers: ConfigViewProvider[]`（非空数组），`components/model-settings.tsx:111` `view.providers.find(...)`／`:211` `.some(...)` **无 null 守卫** ⇒ 全新配置（AC#1 的目标态）下页面渲染会 TypeError。修法二选一（Go 发 `[]` 或页面守卫），**本腿不修**（页面归机主前端会话）；判语＝Go 半附此条条件。
1.4 AC#1 的数据源亲验：空配置 `configuredKnown:true/configured:false`＋presets 12 枚含 `stepfun`（字母序）⇒ 高亮条件成立；用户写过后 `configured:true`⇒"以用户为准"成立。`voice.cloudASRChain/cloudTTSChain` 未配置为 `null`——页面接口本来 `string[]|null`，无双拉缺口。

## 2. A881 三条件（独立现读）

- **ⓐ 248 钉零动＝成立**：`config_route_248_test.go` 最后落笔＝`0d87a681`（248 本体，313 两笔零动）；`git diff 0d87a681 HEAD -- internal/panel/config_handlers.go` 的全部 hunk 均在表块之外（表块＝`// The writable fields` 注释→`WritableFields()` 收尾共 50 行，两版**逐字节相同**，sha256 双侧 `de47e7f8…5db55`）；248 冻结钉 `TestAC1FieldTableCarriesNoApprovalVocabulary` 现跑 PASS。
- **ⓑ 单一 accessor＝成立**：六枚 313 字段字面量在仓内非测试 Go 代码里**只**出现在 `config_write_313.go:43-53`（常量）；并集只在 `acceptedWritableFields()`（`config_write_313.go:104`）算一次，唯二用点＝`config_handlers.go:364/376` 两处拒绝文案；无第三处名册。`cmd/wisp` 的 `switch w.Field` 是逐 case 派发、不是名册：门收了而 switch 缺 case ⇒ 默认支 `没有对应的写入腿`（fail-closed，非 fail-open）。
- **ⓒ 并集等值钉＝成立（牙有两处，具名）**：钉＝`Test313WritableUnionEqualsTheTwoRosters`（`cmd/wisp/panel_config_313_test.go:221`）。突变 B1（`acceptedWritableFields` 截末员）⇒ 钉红，句逐字：`the door reports the union [...voice_realtime_model]; the two rosters say [...voice_realtime_model voice_realtime_provider] (A881: union == table1 + table2)`。突变 B2（`fields313` 删 `FieldVoiceRealtimeModel`）⇒ 该钉**仍绿**（两侧同源），但 `internal/panel` `Test313AdditionalFieldsReachTheLegAndThe248PinIsUntouched` 红（`field "voice_realtime_model" was refused by the door: …不在可写名册内…`）＋`cmd/wisp` 端到端红。⇒ 判语：等值钉对"报告文本≠两表并集"类分歧有牙；"表成员被删"类的牙长在 roster 用例＋端到端用例上，两者合起来无盲区。

## 3. 六枚 setter：拒绝面与原子性（探针，10 拒＋1 控制）

逐发拒绝（door 半态 ×2、leg 半对 ×1、幽灵模型/元素 ×2、`local-sherpa` ×1、非预设建行 ×2（door+leg）、已存在行 ×1、空链元素 ×1）**全部拒绝且 config.toml sha256 前后不变**（`e8436df1…`）；控制发（合法云链）移动 hash 至 `1b8cb8bd…` ⇒ 尺有灵敏度。拒绝句三族归位：`第 2 种拒因：行不存在` / `set both or neither` / `不是内置预设名…没有 protocol 字段`，逐字入 `logs/10-run-and-keys.txt`。

## 4. cmd/wisp 端到端复跑与行为突变

- 复跑：三包 `-run 'Test313'` 全 ok（含 r2 端到端与并集钉）；`TestAC2CredentialSentinelAppearsInNoArtifact` 亲跑 PASS；`TestTicket255` 全 PASS（r2 记的 255 名册两处欠账已由 255-r7 `dcaff922` 清，非 313 悬账）。
- **突变甲（case 换脚）**：`FieldVoiceCloudASRChain` 改调 `SetVoiceCloudTTSChain` ⇒ `Test313EmptyConfigToStepfunRowChainAndPair` 红，句逐字：`receipt does not name "voice.cloud_asr_chain": "设置已写入。 实际落盘的键路径：voice.cloud_tts_chain。…"` ⇒ 钉有牙。
- ★**突变乙（`splitChainValue` 改坏＝按 `;` 切）**：三包 `-run 'Test313'` **全绿**（仓内零钉覆盖多元素逗号链——r2 端到端只写单元素；config 包直接传切片）⇒ **具名覆盖缺口**：多元素链值无仓内钉。本腿 overlay 探针（两元素链）当场红：`element "stepfun/step-2,stepfun/step-2-mini" references unknown model "step-2,stepfun/step-2-mini" of provider "stepfun"`。建议落地腿补一发两元素端到端（修面不在本腿）。

## 5. 密钥面（AC#3 Go 半）

- 拒绝文案现读（并集全 13 枚）：`…不在可写名册内；可写的只有 model_context_window/model_price_in/model_price_out/provider_api_key_ref/provider_base_url/provider_credential/provider_row_create/role_chat_model/role_chat_pair/voice_cloud_asr_chain/voice_cloud_tts_chain/voice_realtime_model/voice_realtime_provider…`。
- 亲跑 canary（`credentialValue` 键）：回执只含 `dpapi:974b…` 引用与状态；磁盘 `api_key_ref = 'dpapi:…'`；回执/磁盘/审计/log（248 钉含 slog 与快照/台账/数据根）**零命中**；载荷只带 `apiKeyRef` 引用。`credentialValue` 只被 `cmd/wisp` 的只写封套绑定，`ComposerRequest` 不声明该键（现读）⇒ `AC#3`＝**既有且未被动**（r2 diff 三枚 hunk 全在新增面，`StoreCredential` 路零改动）。

## 6. 失败集归因

- `internal/panel` 整包现跑：5 红——`TestApprovalCardViewJSONKeysMatchFrontendTypes`＋`TestComposerContractTypesMatchFrontend`（contract 尺；逐字＝`Go Snapshot emits [instructions tasks]`／`Go ComposerState emits [git currentModel modelKnown credentialState credentialKnown]`；`credentialState` 系 `0d87a681`（248）引入、`instructions` 系 `a90677e9`（200），均 313 之前）／`TestC21DesignTokensFourWayAgree`（`read design/assets/tokens.css … cannot find the path`，机主移走）／35r8 两枚（故意红）。**5 枚全在派单已知三族名册内，新红 0 枚**。
- `internal/config` 整包：全绿。`cmd/wisp`：整包不跑（314 死区），`-run` 小射程见 §4/§5。
- 页面侧落笔核：`git log -- frontend/` 最新＝`9e9964f8`；两枚 contract 尺读的 `PanelSnapshot`/`ComposerState` 接口在 `9e9964f8` 前后**逐字节相同**（该笔只加 settings 门）；被裁四笔 `git diff -- frontend/` 为空 ⇒ 旧红与四笔无因果。

## 7. 翻框建议（⛔ 不翻框，只写表）

| AC | Go 半判语 | 页面半（不判，归机主前端会话） |
|---|---|---|
| AC#1 默认选中 | **成立**（读面给 `presets`含 stepfun＋`configuredKnown/configured`；"改过以用户为准"由载荷语义天然成立） | 默认高亮渲染、不点不写 |
| AC#2 两个模型 | **成立**（文本＝`role_chat_pair`；语音＝两链＋realtime 两枚；写后读回亲验；⛔ 未改 schema 语义） | 下拉/手填/表单交互 |
| AC#3 界面录 key | **成立＝既有且未被动**（DPAPI 路＋不回显四面亲验） | 输入框清空/状态显示 |

**不翻框的两个具名挂账**：①§1.3 `providers:null` 形状分歧需裁（Go 发 `[]` 或页面守卫）；②§4 多元素链覆盖缺口需补钉。二者都不是"AC 不成立"，属收尾清单。

## 8. 具名没做的格

整包 `./cmd/wisp/`（314 死区）· `-race`（未要求）· 页面侧渲染与视觉判语 · 未 push、未动 `Status:`/票名/AC 框。
