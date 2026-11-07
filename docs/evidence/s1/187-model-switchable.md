# 187 裁决表 — 聊天框那一排的模型与思考档位要能当场改（Go 侧半）

写腿：2026-11-07 19:2x（GLM 会话，Wisp 编队写腿）。票面全文：`.scratch/wisp/issues/187-*.md`。
本票只交 Go 侧的读面＋写腿＋快照字段（AC#9）；界面那半归前端会话，需求表见
`docs/evidence/s1/187-panel-ask.md`。

## 0. 一句话

`llm.roles.chat` 的两枚「可改」里：**思考档位（thinking_intensity）＝ next_task 生效**，
写腿走 `config.set` 新增字段 `role_chat_thinking`（第三张可写名册），配置侧 setter
`Manager.SetRoleChatThinking` 落盘＋内存同写；**模型＝restart 生效**，写口
`role_chat_model`（248 名册）／`role_chat_pair`（313 名册）**已在册，本票零新增**；
读面新键 `modelSwitch`（reader-gated，213 同款范式）把可选目录＋当前档＋当前模型自己
声明的档送上线。cmd/wisp 生产 leg 的最后一行 case 接线不在本票写面（见 §6）。

## 1. 起手三核（现量 vs 票面）

| 票面断言 | 现跑读数（2026-11-07） | 结论 |
|---|---|---|
| 尺1 `thinking_intensity` 枚举与默认值在 schema.go `:291`/`:359` | `schema.go:312`（default "off"）/`:380`（ThinkingLevels）；枚举校验在 `validate.go:274-282`（票面写 `:201`） | **行号漂移，字段俱在** |
| 尺2 ThinkingIntensity 被真用：`internal/agent/prompt.go:113-114,225` | 现量 `prompt.go:135-136,254`；上游 `internal/llm/resolver.go:39,204`（RoleSettings.ThinkingIntensity） | 仍在用，行号漂移 |
| 尺3 角色按 provider+model 成对：`schema.go:280-287` | 现量 `schema.go:305-312` 附近（Role 结构） | 成立，行号漂移 |
| 尺4 面板写入口恰好四枚（`bridge.go:42-45`） | `bridge.go:42-75`＝7 枚 method：四枚 `panel.*`＋`config.get`/`config.set`（248）＋`config.models.list`（320） | **票面整段过期**：248/313/320 已开设置通路，模型写口已在册 |
| 尺5 三档生效级别 | `cmd/wisp/config_readers_255.go`＋sectionwrite.go | 成立，且 370 已把 roleSettings 改成每任务重解（票面未载） |

**新增现量（票面未载，本票方案据此定形）**：
- `panel.FieldRoleChatModel`（248 名册）＋`panel.FieldRoleChatPair`（313 名册）已可写
  `llm.roles.chat.model`／provider+model ⇒ 模型那半的写口**已经存在**，AC#2 的
  「判据＝改动面最小」落在：本票只补 thinking 一枚字段。
- `cmd/wisp/run.go:1632-1647`（票 370）：roleSettings **每任务**从
  `rt.mgr.Config()` 重解；`run.go:562-576`：endpoint/provider 链**装配一次**不重建。

## 2. AC 逐条

### AC#1 生效定档（现跑，不按「应该可以」写）

- **换思考档位＝next_task**。尺①（读者是每次现读）：
  `cmd/wisp/run.go:1632-1647` execute() 每任务 `llm.NewResolver(cfg,nil).ResolveRole(llm.RoleChat)`，
  cfg 取 `rt.mgr.Config()`（活配置）；尺②（写入同步进内存）：
  `internal/config/sectionwrite.go:305-336` commitSectionWrite 先改 `m.cur` 再落盘。
  齿：`Test187ThinkingWriteLandsTheKeyAndNamesIt` 同时断言文件与内存两处都到。
- **换模型＝restart**。尺：`cmd/wisp/run.go:562-576` 装配根一次性
  `ResolveRole`→`rt.endpoint`/`rt.provs`；全树仅 `run.go:2065-2068` 读 `rt.provs[0]`，
  无任何重建路径 ⇒ 同一次运行里模型通路不重建（248 收据的既有判断，本票复核仍成立）。
- **面板收据按实测档位说**：cmd/wisp leg 对 [llm] 一律盖 restart 戳；thinking 的实测档
  是 next_task ⇒ panel 门在 187 名册范围内把收据档位改为实测值
  （`config_handlers.go` write() 的 tier 修正，`fieldTier187` 表）。
  齿：`Test187ThinkingReceiptNamesTheMeasuredTier`（正例说 next_task／控制例保 248 字段
  restart 戳）。

### AC#2 写入口走哪条通道

复用 `panel.mode.request` 同款「请求→处理器→写者」形状里**已长出的设置通路**
`config.set`（248 所开，313/320 已两度追加名册）：新增字段 `role_chat_thinking`
（`config_write_187.go` 第三张名册；`config_handlers.go` knownWritableField 三表并问、
acceptedWritableFields 三表并集、takesValueFor 并列）。
**判据＝改动面最小，成立**：零新 C17 方法（bridge.go 七枚 method 一枚未增，C17 白名单
零新增 ⇒ 无需具名报契约审批）、零新 socket 名；模型半零新增（248/313 既有两字段）。
「写成契约变更入档」＝本文件 §1/§2 即入档。

### AC#3 绝不渗权限档位

- 写腿里没有任何一处碰权限写口：`grep -n "Modes.Set\|perm\." internal/panel/config_write_187.go internal/panel/model_switch_187.go internal/config/settings_187.go`＝0 命中；
  齿：`Test187ConfigWriteLegTouchesNoPermissionSurface`（config 侧）＋
  `Test187PanelFilesTouchNoPermissionSurface`（panel 侧）逐文件读字节钉死。
- `risk.Mode` 与这两枚字段无数据通路：快照新键 `modelSwitch` 只载 provider/model id
  与档位词；config.set 门禁的 locked-family 守卫（risk./fs./net.…）在 187 名册之前，
  未动一字。`[llm]` 的方向裁决＝neutral（sectionwrite.go lockedDirectionOf），
  面板写入 `Source=SourcePanel` 也不能借道（ErrPanelLoosenRefused 与本票无关但同在闸上）。
- 「换更贵模型/抬档位算不算放宽」：不算——它们不触碰 R20 的三档语义，也不触碰任何
  risk 键；R20 的「调松要卡」只对 permission_mode 这一个维度成立（本票复述口径，
  不改它）。

### AC#4 清单要真

- 下拉可选行只来自真目录：`ChatSwitchOptions()`（`internal/config/settings_187.go`）
  只列 `enabled=true` 的目录条目（resolver 对 disabled 的拒绝是运行时事实：
  `internal/llm/resolver.go:128-131`）。
  齿：`Test187ChatSwitchOptionsNeverOffersWhatTheResolverWouldRefuse`（disabled 行不得出现）。
- 档位只列当前模型**自己声明**的：读面 `ThinkingLevels` verbatim（可能少于四档、可能空）；
  写面 `requireDeclaredThinkingLevel` 拒收模型未声明的档。
  齿：`Test187ThinkingWriteRefusesALevelTheModelDoesNotDeclare`（"high" 对声明 off|low
  的模型＝拒；文件零字节动）＋`Test187ThinkingWriteOnAModelThatDeclaresNoLevelsIsEnumOnly`
  （未声明档的模型＝枚举四值全可写，不发明更严的规则）。
- 空目录／没配 provider／全 disabled：显式「没有可用模型」句（`Reason`），绝不把空列表
  妆成「没得选」。
  齿：`Test187ChatSwitchOptionsEmptyShapesSaySo`（两种空形各带各的句子）。

### AC#5 生效范围具名

- 换**模型**：整个进程（装配根不重建），重启后生效；换**档位**：从**下一任务**起
  （每任务重解），当前正在跑的那一轮不受影响。均见 AC#1 的尺。
- 「界面显示的那一枚＝运行时真用的那一枚」：线携带两个读者各自的值——
  `composer.currentModel`（boot 端点解析的模型，pump.go `src.Model()`）与
  `modelSwitch.current`（活配置里的 chat 角色）。重启前两者**允许不同**，页面拿两个
  真值可以如实说「下个任务/重启后是 X」，而不是把显示值谎报成生效值
  （146 值保真的教训：不做第二真相源，两值各出自其读者，测试各钉各的：
  `TestComposerCurrentModelTravelsOnlyFromItsReader` 与
  `TestTheModelSwitchKeyTravelsOnlyFromItsReader187`）。

### AC#6 覆盖面（谁还会读旧值）

- 上下文预算 D15：`grep -rn "BudgetsFor" --include=*.go internal/ cmd/ | grep -v _test.go`
  ＝读者是 `internal/agent/loop.go:241`（每任务 `BudgetsFor(window)`），window 来自
  per-task 重解出的 `ep.ContextWindow` ⇒ **换模型换窗口从下一任务起随新 ep 生效**；
  正在跑的那一轮持旧 assembler（构造时定死），不中途换——结论＝**轮次边界生效**，
  不拒绝、不中断、不静默超限。
- 语音链路：voice realtime 的 provider/model 是独立字段（`voice.realtime.*`），
  不读 `roles.chat`，零波及。
- roles 其它使用者：memory_extract/handoff/summarize 各读各的 Role 字段
  （resolver.go:186-227 逐 role 取），chat 的改动不外溢。

### AC#7 契约轴

`docs/PLAN.md`/`docs/specs/**`/thresholds/golden/审批超时常量/`allowlist.txt` 零触碰
（git 读数见收件报）。凭据面：`ChatModelOption`/`ModelSwitchSection` 结构性只载
provider 与 model id 字符串＋档位词，无任何凭据字段；`api_key_ref` 不上新键。

### AC#8 门禁读数

见 §4（终态读数取自最后一枚 commit 之前的最终工作树状态；commit 后未再改码）。

### AC#9 界面需求表

`docs/evidence/s1/187-panel-ask.md`（本文件交付时一并落盘）。

## 3. 变异牙（三发，逐发红后 sha256 逐字节还原）

| 发 | 突变 | 红 |
|---|---|---|
| M1 | settings_187.go：摘掉「模型声明档」检查＋Check 置 nil | `Test187ThinkingWriteLandsTheKeyAndNamesIt`＋`Test187ThinkingWriteRefusesALevelTheModelDoesNotDeclare` |
| M2 | settings_187.go：摘掉 Enabled 过滤 | `Test187ChatSwitchOptionsNeverOffersWhatTheResolverWouldRefuse` 等 3 枚 |
| M3 | config_handlers.go：tier 修正改为盖 restart | `Test187ThinkingReceiptNamesTheMeasuredTier` |

还原核对：`sha256sum` 四文件与突变前一致（settings_187.go＝949347eb…、config_write_187.go＝ec865c68…、config_handlers.go＝a7b9c923…、pump.go＝be37283a…）。

## 4. 门禁读数（终态）

- `gofumpt -l`（10 枚点名件）＝空。
- `go vet ./internal/panel ./internal/config` rc0；`go build` 两包 rc0。
- `go test ./internal/config -count=1`＝**ok**（含 7 枚新 187 用例）。
- `go test ./internal/panel -count=1`＝3 枚红函数，**全部与 213 交付时 HEAD 同态**：
  - `TestApprovalCardViewJSONKeysMatchFrontendTypes`：报文 `Go Snapshot emits [commands modelSwitch] that interface PanelSnapshot does not declare`——213 交付时该尺已因 `commands` 红（213 交付说明：等界面腿写 panel.ts，Q-51）；本票 `modelSwitch` 加入同一名单，红测试函数数不变（2 枚对账尺）。
  - `TestComposerContractTypesMatchFrontend`：同上＋`files` 键既存红（190）。
  - `TestC21DesignTokensFourWayAgree`：既存红（ball tokens，HEAD 同态，213 交付时已归因）。
  - 四键字节钉（reader-gated 无读口形）保持绿：`TestTheModelSwitchKeyTravelsOnlyFromItsReader187` 断言无读口 pump 恰发 4 键。
- `sh scripts/d22scan.sh`：scan 本体 1 finding＝本票注释引用 `docs/evidence/s1/187-model-switchable.md` 在文件落盘前先存在过（phantom citation，ticket 212 ban #9）——本文件落盘后复跑见收件报读数；d22scan 自身包另有 2 枚测试红＝同一 finding 的两面，文件落盘后消失。

## 5. 红绿归因（第三核）

3 枚红测试函数与 213 交付 commit（`0e68772a0`）时的 HEAD 同态归因：对账尺两枚本就
预定红（等前端 panel.ts 声明，Q-51），TestC21 既存红（ball tokens）。本票零新增红
测试函数。

## 6. 诚实边界（没做到的）

1. **cmd/wisp 生产 leg 的最后一行接线不在本票写面**：`cmd/wisp/panel_config_store.go`
   的 `ApplySetting` switch 需要一行
   `case panel.FieldRoleChatThinking: written, err = s.mgr.SetRoleChatThinking(w.Value)`
   才能把门禁放行的请求落到 setter。本波 cmd/wisp 是 188 的独占写面（写面互斥铁律），
   本票 ⛔ 碰 ⇒ 当前树上门禁会认 `role_chat_thinking` 并路由，leg 会按既有 default 分支
   拒绝（「没有对应的写入腿，未落盘」，不静默）。接续点已写进
   `docs/evidence/s1/187-panel-ask.md` §4，交编排者派 188 或下一枚 cmd/wisp 腿。
2. 装配根把 `ChatSwitchOptions` 接上 `PumpSources.ModelSwitch` 的那几行同在 cmd/wisp
   （run.go 装配区），同上不碰；读面与 reader-gated 泵路径已备好并测毕。
3. 前端 panel.ts 的 `modelSwitch` 声明归界面腿（Q-51）；两枚对账尺因此保持预定红。
