# 242 落地前撞钉预检（只读腿 `242-precheck-1`）

> 状态：**骨架**（§0 起手锚已填；§1–§4 填充中）。禁令遵守：零 Go 命令、票面勾选框零接触、产码零改动、`frontend/**`/`design/**` 零读、冻结件零接触。证据来源＝grep/sed 读码＋引既有裁决表，逐枚带 file:line。

## §0 起手锚

- 锚点：`HEAD = 0781ac70`（dev），起手 `git status --porcelain -- internal cmd tools` ＝ **空**（名字册干净；工作树里已有的 `.gitignore`、`design/assets/**` 删除、`probes/152|161` 改动**非本腿所留、非本腿地界**，只记归因不判）。
- 读法：全部结论出自 grep/sed 静态读码（`cmd internal tools docs scripts .scratch` 显式根）＋三份裁决表转引（`197-ac5-selfapproval-v1.md`、`145-snapshot-growth-r3.md`、台账 `Q-74` 行）。**零动态读数**（禁跑 Go），所有"今天绿/今天红"皆注明出处。
- 票面：`.scratch/wisp/issues/242-the-grant-binding-layer-and-the-panel-facing-read-surface-both-have-zero-rulers.md`，owner 默认甲＝绑定层＋出向读面两枚零尺判据；落点甲/乙之争已挂 `Q-74`（台账 :9845），默认甲。

## §1 两枚判据的落点与形状建议（不选形，归编排者）

### §1.1 判据①（绑定层零尺）——今天的确切读数

**绑定层的生产面（铸与花各一处，无第三处）：**
- 铸：`internal/agent/approval/queue.go:162` —— `it.bind = bindDigest(corr, d.TaskID, d.Tool, d.LevelString(), q.seq, d.Args)`（在 `push` 内，吃 queue 自己加的 `q.seq`）。
- 花：`internal/agent/approval/approval.go:292-306` —— `grantStore.spend(nonce, bind)`，命中 nonce 后**无条件 delete** 再 `return equalSecret(stored, bind)`（:302-303）；比对的存储值在 `issue`（:284-288）时写入，写的是 `it.bind`（`queue.go:336`，`grantNonce` 内）。
- 摘要函数本体：`internal/agent/approval/approval.go:253`（`bindDigest(corr, taskID, tool, level, seq, args)`，域分隔注释 :250-252）。
- 全仓 `grep -rn "bindDigest" --include=*.go cmd internal tools`＝**3 枚文件**：`approval.go:250,253`、`queue.go:162`、`cmd/wisp/subagent_selfapproval_197_test.go:14,320`（后者只是注释与 t.Logf 引语，**不是断言**）。

**测试调用者计数（票面 ① 的第一问）：**
- 生产调用者：**2 处**（铸 :162、花 :292 经 `issue`/`spend`）——注意 `spend` 的调用者 `allowScoped`（`queue.go:376`）传的 `bind` 就是 `it.bind`，所以"花"路径上没有第二处独立绑定决策。
- 测试调用者：**0 枚断言**。`_test.go` 里出现 `bindDigest` 的只有被验收文件自身的两行注释/日志（`subagent_selfapproval_197_test.go:14,320`）；`internal/agent/approval` 包内 12 枚测试文件 grep `bindDigest|grantStore|spend`＝零枚断言命中（`ticket224_reply_grant_test.go:161,294` 与 `ticket97_alias_direction_test.go:72` 只是注释里的英文动词 spend）。
- 旁证（转引，非本腿跑）：`197-ac5-selfapproval-v1.md` §1.1 **D2 发**——`grantStore.spend` 整枚删掉绑定比对后，定向尺与整包 `go test ./internal/agent/approval/` 全绿（`ok 0.392s`）；§2.4 结论＝「`queue.go:44-47` 那句 bind 覆盖……今天是一条没有仪器的断言」。

**AC#1（真·跨卡）的构造前提——载具缺陷逐字核过：**
- `cmd/wisp/subagent_selfapproval_197_test.go:108-109`：`rt.bridge.Execute(ctx, agent.ToolRequest{ TaskID: taskID, CorrelationID: taskID, ...})` —— corr 与 taskID 同值，两次 `fs.write` 挂出的两张卡在 `q.push` 里都会被 `byID[corr]` 冲突检测改写成 `taskID#seq` 形状（`queue.go:152-153`），**但 :466 那发传给 `Native().Allow(card2.Corr, card1.Grant)` 的 `card1.Grant` 在 card2 挂出前已随 card1 结算被花/撤销**（`dropLocked`→`grants.revoke()`，`queue.go:288`），所以那一发实际量的是新鲜度，不是跨卡绑定。这与 `197-v1` §1.3.2 的 13 发读数（`child=corr1=corr2=同值`）一致——**方向修正是：corr 其实不同**（`#seq` 后缀），真正同值的是"第二发时 A 的令牌必已不活"。AC#1 要"两枚同时活着"＝两枚卡都挂出、都不答，然后互花——`startChildWrite197` 现在的写法（答完第一张才挂第二张）结构上做不到，需按票面重排：两发都点火、轮询两张卡、再互花。
- 权限侧落点确认：`allowScoped`（`queue.go:361-398`）里 `q.lookupForAllowLocked(corr)`（:231-236）**只查 `byID` 精确键、不查 alias**，alias 只进 `lookupForRefusalLocked`（:250-265）——所以"把 A 的活令牌递给 B 的 corr"今天确实只有 `spend` 里那一次 `equalSecret(stored, bind)` 在挡，D2 发证明它无尺。

### §1.2 判据②（出向读面零尺）——今天的确切读数

**PanelItem 的字段全集（票面 ② 的第一问：哪些进 Snapshot、哪些被钳）：**
- 结构体：`internal/agent/approval/ui.go:50-57` —— `CorrelationID / Tool / Level / Reason / Paths []string / Depth int`，**共 6 枚字段，无 Grant、无 Params、无 Args、无 RulesHit、无 SessionOverrideBlocked、无任何可答复字段**。头注释 :45-49 逐字「deliberately no grant field, no Allow capability」。
- 唯一构造点：`internal/agent/approval/queue.go:524-533`（`viewLocked`）——逐字段显式拷贝：Corr→CorrelationID、Dec.Tool、Dec.LevelString()、Dec.Reason、`append([]string(nil), it.Dec.Paths...)`（**Paths 是深拷**，票 146 钉过）、`q.position(it)`→Depth。**Dec.Params / Dec.Args / Dec.RulesHit / Dec.Blacklist 在此被钳掉**（结构性不出队）。
- 消费者全集（grep `PanelItem` 全仓产码）：`PanelAPI.Head/View`（`gate.go:697-699`）→ `Replies.Head/View`（`replies.go:473-489`）→ **唯一生产消费者** `cmd/wisp/approval_reply.go:374-403`（console `head`/`view` 两动词，`panelItemLine` 只读 CorrelationID/Level/Tool/Depth/Reason/len(Paths) 六枚）。
- **PanelItem 进不了 Snapshot**：快照那条路走的是**另一个类型**——`Queue.LiveApprovals()`（`pending_read.go:104-122`，返回 `LiveApproval{CorrelationID, Decision tools.Decision, Position}`）→ `cmd/wisp/panel_pump.go:58-76 liveVerdicts`（显式字段拷贝成 `panel.NativeVerdict`）→ `pump.go:98-113 NativeVerdict.CardView()` → `panel.ApprovalCardView`（`internal/panel/approval.go:39-59`，10 枚 JSON 键）。`197-v1` §1.3.1 逐字：「`approval.PanelItem` 只服务 `PanelAPI.Head()/View()`（`gate.go:697-699`），没有任何一条路把它的字段带进快照字节」。⇒ **票面 ② 问"哪些字段进 Snapshot"的准确回答是：PanelItem 一枚都不进；进 Snapshot 的是 LiveApproval/ApprovalCardView 这条平行路**。这枚区分对选形有分量（见 §1.3）。

**三条守注释（票面现量表逐字复认，行号现核）：**
1. `internal/agent/approval/approval.go:223` —「The panel-facing surface (PanelItem) has no grant」
2. `internal/agent/approval/pending_read.go:9` —「PanelItem deliberately carries no Params and no grant」
3. `cmd/wisp/approval_reply.go:28` —「PanelItem has no grant field」
   全仓 `_test.go` grep `PanelItem`＝3 枚文件、**零枚字段名册断言**（`fakes_test.go:104` 只是 fake UI 存 views；`ticket146:338,346,349` 只钉 Paths 深拷与 Tool/Reason 值保真；`subagent_selfapproval_197_test.go:696` 是变异脚注）。
- 旁证（转引）：`197-v1` §1.4 **M1a-sweep**——`PanelItem` 加 `Grant` 且 `viewLocked` 填真令牌后**全仓 0 枚红**（cmd/wisp `ok 143.673s`、approval `ok 0.598s`）。⇒ "零仪器"读数成立，且**额外发现**：M1a 落在出向读面（PanelItem），而那发全绿的同时票 146 的 `TestLiveApprovalsSharesNoReferenceSlotWithTheQueue` 也绿——因为它扫的是 `tools.Decision` 的引用槽，不是 PanelItem。**两把现成的反射尺（146 槽位普查、90/224 权威词面普查）今天都没把 PanelItem 圈进射程**，这就是判据②要补的空。

### §1.3 新符号落点建议（两形各一条，归编排者选）

**判据①（绑定层）——新符号候选：**
- 断言函数名候选：`TestGrantSpendsOnlyOnItsOwnBinding`（甲形，`internal/agent/approval` 本包，white-box 可直呼 `q.push`+`grantStore`）。
- 需要的最小夹具：两枚不同 corr 的 `tools.Decision` 先后 `q.push`（`queue.go:141`，本包测试已用同款 `mustPush`，`pending_read_test.go:45-52`）→ 各拿一枚 `q.grantNonce(it)`（:326，**unexported，甲形白盒才够得着**）→ A 的 nonce 花 `allowScoped(B.corr, A.nonce, false)` 断言 `ErrBadGrant`；正控两枚各自花自己的断言 nil。⚠ 排序注意：`spend` 命中即删（:302），所以每张卡要嘛当"被借证方"要嘛当"正控"，两枚卡各做一发需要**四枚卡**（或一枚卡重挂一次）——形状建议里要点名。
- 拒因精度：票面要「拒因指名绑定不对」。今天 `ErrBadGrant` 的文案是三合一（"缺失/已用/与本次请求不绑定"，`ui.go:130`）——`197-v1` §2.2 的 D1 发已证明这枚文案**混用门/权限**。判据①若逐字断言"绑定不对"需要新错误值或至少 `errors.Is` 分叉，**这属于新增契约面，甲形本包内自洽可做，但不应顺手改 `ErrBadGrant` 文案**（改动会影响 `subagent_selfapproval_197_test.go:254` 的 `errors.Is` 集合——那枚在 `197-r3` 的地界里，票面已写明串行）。
- 乙形（扩 `internal/panel/l2_grant_boundary_test.go`）：**不适配判据①**。那族的射程是"入向 envelope＋route 词面"（`grantFieldWords`/`grantRouteWords`，:183-195），bindDigest 是**原生侧内部**的摘要，不出现在任何 wire 键名里——乙形那套词表扫不到它。⇒ 判据①天然只有甲形。

**判据②（出向读面）——新符号候选：**
- 断言函数名候选：`TestPanelItemCarriesNoGrantOrAnswerVerb`（甲形，`internal/agent/approval` 本包）。
- 形状两半（`197-v1` §1.4 自己给的那形＋票面 AC#2 的能力半）：
  1. **词面半**：反射 `reflect.TypeOf(approval.PanelItem{})` 走 `NumField`，字段名 lowercase 后含 grant 词表任一（词表可抄 `internal/panel/l2_grant_boundary_test.go:183-188` 的 `grantFieldWords`，或收窄成 `grant|allow|approve|decision|outcome|verdict|answer`）⇒ 红。**外加字段数钉**（今天 6 枚）防"换名不换实"——这与 `internal/tools/subagent_197_test.go:833-881` 的 `SubagentDeps` 字段数钉同族（那枚钉过 M3 突变，有牙的先例）。
  2. **能力半**：`reflect.TypeOf((*PanelAPI)(nil)).Elem()` 方法集＝恰好 `{Reject, Head, View}` 三枚（`ui.go:167-171`），出现 Allow/AllowSession/grant 形方法 ⇒ 红。**已有先例可抄**：`internal/tools/ticket90_test.go:433-438`（ModeSource 恰好 1 方法）与 `cmd/wisp/subagent_selfapproval_197_test.go:541,568-572`（tools.Gate 方法集恰好 2 枚）——同族形状，写法直接移植。
- 乙形（扩冻结件射程）：技术上**能**盖词面半（`bindableKeysOf` 有递归），但 PanelItem **没有 json tag**（`ui.go:50-57` 全裸字段），`bindableKeysOf` 会 fallback 到 Go 字段名（:305-309）——能扫到，但要先过 `inboundTypeRegistry`（:1720-1727，把 PanelItem 加进去＝改冻结件本体）。⚠ 且那族 AST 半（`scanGrantBoundary`）只找 **json.Unmarshal 调用点**，PanelItem 不走 decode，AST 半对它是结构性盲区。⇒ 乙形要真盖住得动三处（registry、词表、AST 假设），不是"扩一条列表"那么小——这枚代价差可进 Q-74 的拍板材料。

## §2 撞钉清单（逐枚带 file:line＋断言原文）

> 结论先行：**判据①/② 的新符号若按 §1.3 建议的形状与命名落位，本腿逐枚排查后未发现任何一枚既有钉会被顶红**。原因有二：(a) 新断言是**新增用例**，不改任何既有断言的输入或输出；(b) 现有的同族反射尺射程都以"枚举型集合/词面"划界，`PanelItem`/`bindDigest` 都不在其集合内。但下列 7 枚是**擦边**的，逐枚给出会红/不会红的判据与理由。

### §2.1 擦边钉逐枚

1. **`internal/agent/approval/ticket146_liveapprovals_backing_test.go:190-199` `queueSlotPaths`**（8 枚槽位名册钉）
   断言：`typeCensus(reflect.TypeOf(got[0].Decision))` 必须逐项等于 `queueSlotPaths`（:221-231，逐字红句「数目不符说明 tools.Decision 增删了引用字段而 queueSlotPaths / cloneDecision 没跟上」）。
   与新符号关系：**不撞**。它扫的是 `tools.Decision`，判据②扫 `PanelItem`——两枚类型无交集。但若写腿**顺手**把 PanelItem 的反射扫也接进 `typeCensus`（不要），PanelItem 全是值型字段会走 :129-133 的"值安全"分支，不红——枚举失败的是"加了引用字段"那一形，与"加了 grant 字符串字段"不同形。⇒ **判据②的词面半必须独立写，不能复用 146 的普查函数**（它对 string 字段是放行的，`Grant string` 照不红）。
2. **`internal/tools/ticket90_test.go:443-453` `TestTicket90ModeCarriesNoAllowAuthority`**（票面禁区点名的 :431 反射钉）
   断言：`reflect.TypeOf(Decision{})` 逐字段 `strings.Contains(low, "allow"|"approve"|"grant") || f.Type == Answer` ⇒ 红；`reflect.TypeOf(agent.ToolRequest{})` 同法。
   与新符号关系：**不撞，且结构同源**——判据②的词面半就是把同一把尺的射程加到 `approval.PanelItem` 上（票 224 也撞过这枚，票面已记）。⚠ 派单措辞建议：别写"扩 ticket90 那枚的射程"（那是改别人的测试），写"按同一形状在 approval 包新写一枚"——甲形。
3. **`internal/tools/grant_test.go:508-516` `TestTicket224GrantSeamCarriesNoCallerSuppliedAuthority`**
   断言：`reflect.TypeOf(agent.ToolRequest{})` 逐字段含 `session|grant` ⇒ 红（「lets the CALLER name the session or hand in a grant」）。
   与新符号关系：**不撞**。射程是 `agent.ToolRequest`。但它是**第三个同族**——票面禁区那句「不许把 SessionID/grant 塞进 tools.Decision 或 agent.ToolRequest」的三处反射钉（ticket90:443-453、grant_test:508-516、ticket224_setter_scope_test.go:37-67）合起来就是"判据②该长成什么样"的族谱；新钉落 approval 包不进它们任何一个的 grep 根（`internal/tools`），零交集。
4. **`cmd/wisp/subagent_selfapproval_197_test.go:557-594` `Test197NoAllowDoorIsReachableFromASubagentsAssembly`**
   断言：`allowDoorMethods197 = ["Allow","Native","DecideFromNative","DecideFromPanel","GrantNonce"]`（:541）在四枚装配根（`tools.SubagentDeps/TaskDeps/Options`、`agent.Options`）的类型图里零命中（:563）；`tools.Gate` 方法集恰好 `{PendingWindow, PendingApproval}`（:569-572）。
   与新符号关系：**不撞，但有两个必记的耦合**。①它的射程是**装配根的类型图**，`PanelItem`/`Queue` 都不是那四枚根的字段类型（`agent.Options` 无 approval 字段，grep 过），新钉不改变它的输入。②**判据①的测试若放在 `cmd/wisp`（不推荐）并需要调用 `q.grantNonce`，那枚 unexported 方法名恰在 `allowDoorMethods197` 词表里**——不会红（它只扫类型图，不扫测试文件），但写腿容易误读。⇒ 又一条判据①该落 approval 包本包的理由：`grantNonce` 本来就够不着，甲形白盒直呼即可。
5. **`internal/panel/l2_grant_boundary_test.go:1595-1709` `TestGrantVocabularyIsNotSatisfiedByTheRealEnvelopes`**（冻结件）
   断言族：`grantRouteWordWitnesses` 11 行双向对账（:1652-1678）；`grantFieldWords` 对三枚种植 envelope 必红（:1616-1628）。
   与新符号关系：**不撞（判据②不新增 wire 键）**。这是票面最怕的那枚（A471 那笔账），逐字核对结论：那族全部读数围绕 `ComposerRequest/ModeRequest/AttachmentPayload/AttachmentRef`（`inboundTypeRegistry`，:1720-1727）与 `knownComposerMethod`（`bridge.go:146-152`）——**全是入向**；`PanelItem` 在 `internal/agent/approval`，`internal/panel` 不 import approval（grep 产码零命中），AST 半也只 parse 本包。甲形落新钉与它零交集；乙形要动它本体（见 §1.3 代价）。
6. **`internal/panel/pump_test.go:123,291` 四键钉**（`"composer,generatedAt,pending,results"`）
   与新符号关系：**不撞**。判据②不给 Snapshot 加键（Snapshot 的 pending 走 `ApprovalCardView`，判据②钉的是 PanelItem——两个类型）。⚠ 但派单里要防一个**反向误判**：写腿若把判据②错误地写成"扫 `panel.ApprovalCardView` 的字段"（它才是进 Snapshot 的），会撞 `TestApprovalCardViewJSONKeysMatchFrontendTypes`（`approval_test.go:105-136`）——**那枚本来就红**（`197-v1` §1.5 与 `145-growth-r3` :117 同名在册：前端接口与 Go 结构体键不对账的常红），红上加红会把"新增判据红"与"既有常红"混在一口锅里，突变读数就废了。判据②的反射对象**必须是 `approval.PanelItem`，不是 `panel.ApprovalCardView`**——这一句是本预检最想说的一条。
7. **`cmd/wisp/leg_dispatch_gate_133_test.go:178-211` `TestAC1AC2DispatchHopGate133`**（票 33 那族）
   断言：本包（cmd/wisp）腿普查 ≥ `minLegs133=11`（:148）；`legCovers133` 四行 registry（:169-174）每行 test 名必须是真 `func TestXxx`；孤儿 sink 从 main 不可达 ⇒ 红。
   与新符号关系：**不撞，条件是判据落在 approval 包**。若两枚新钉按甲形落 `internal/agent/approval`，cmd/wisp 的腿普查数与 registry 全不变。若落 cmd/wisp（乙选择之一）：普查**下限**是 `>=` 语义（:187），**加**腿不红——但 `coverageReds133` 会把新测试名与 `legCovers133` 对账，新 test 名若恰好撞 registry 四行里任何一行（不会，名字带 242）或让 minLegs 语义漂移才需要注意。⇒ 结论：落 cmd/wisp 也不红，但落 approval 包更干净（不进这枚 gate 的普查宇宙）。

### §2.2 信封钉的射程问题（任务 ④）

**票面写的 `TestAC11SettingsKeysAreTheOnlyNewOnesOnTheSharedEnvelope` 实际名字是 `TestAC1SettingsKeysAreTheOnlyNewOnesOnTheSharedEnvelope`**（`cmd/wisp/panel_config_248_test.go:369-379`，AC1 不是 AC11——派单转述时号记岔了一枚）。
- 断言原文（:371-378）：`jsonKeys248(reflect.TypeOf(panel.ComposerRequest{}))` 必须含 `configField/configProvider/configModel/configValue` 四键、不含 `credentialValueKey`。
- **射程**：`panel.ComposerRequest`（`internal/panel/bridge.go:91-121`）的字段名枚数——它是**入向 envelope**（页面→Go 的 decode 目标）。
- **盖不盖到 PanelItem 的新读面**：**不盖**。三重不盖：①`PanelItem` 不是 `ComposerRequest`，无嵌入关系（ComposerRequest 嵌的是 `AttachmentPayload`）；②`jsonKeys248` 只递归 anonymous struct 字段（:394-402），不追任何外部包类型；③方向也不对——ComposerRequest 是入向，PanelItem 是出向，判据②即便某天给 PanelItem 加 json tag，也进不了这枚钉的反射根。⇒ 同族的真正"出向侧对账尺"只有 `TestApprovalCardViewJSONKeysMatchFrontendTypes`（`approval_test.go:105`，对 `ApprovalCardView/ResultChunk/Snapshot` 三枚，**今天在册常红**）——判据②若只做"新增 json 键"那一形确实有这把尺（票面 :47 那句对），但它红的读数**不可用**，且它对 PanelItem 同样零射程。

## §3 我可能判错的条目

1. **"PanelItem 不进 Snapshot"的强表述**：本腿依据 `197-v1` §1.3.1＋今日 grep（无第二条 PanelItem→panel 包的路）。若编排者读过某个我没 grep 到的旁路（例如未来某分支已经把 `Queue.view` 接进 pump），§1.2 的表述要收窄。我 grep 的根与词：`PanelItem` 全仓产码 17 行命中、逐行读过（清单在 §1.2），`grep -rn "CarlosShao/wisp/internal/agent/approval" internal/panel/*.go` 产码零命中。判错概率低但表述是全称否定，特此标注。
2. **AC#1"造不出两枚同名卡"的票面读数**：票面 :14 说"两张卡共用同一个 correlation id"（照 `197-v1` §1.3.2 的 13 发读数）。本腿细读代码后判**票面那句要修正**：`q.push` 的冲突检测（`queue.go:152-153`）会把第二张卡改成 `taskID#seq`，corr 其实不同；真正让 :466 那发测不到绑定的是**时序**——card1 结算时 `revoke()` 已把令牌烧了（`queue.go:288`），第二发拿到的是已空的 grantStore。若我读错了（比如 corr 在 :153 的 `d.CorrelationID` 为空时另有分支），AC#1 的构造建议要重写。判据不受影响（今天照样零尺），但"前置＝先把 :109 改成真 corr"那半句**可能不必要**——票面 :19 的前置条款本腿持保留意见，交编排者。
3. **`minLegs133` 的 `>=` 语义**：我按 :187 的代码判"加腿不红"。若 `enumerateLegs133` 的腿普查对"测试文件里的断言函数"另有分类（我读过 :306 的 `tests map[string][]*decl133` 但没逐行追它怎么分类非 leg 的 test 函数），落 cmd/wisp 的新钉可能被卷进 coverage 对账。甲形（approval 包）不受此不确定影响。
4. **`ErrBadGrant` 文案是否该在判据①里逐字断言**：票面 AC#1 要"拒因指名绑定不对"，本腿判"今天的三合一文案满足不了逐字要求，得新错误值或 `errors.Is` 分叉"（§1.3）。但"新错误值"也可能被编排者判为过度设计（拒因枚举是 D37/错误模型的契约面）——若判"现文案够用"，判据①断言就写 `errors.Is(err, approval.ErrBadGrant)`，精度损失与 `197-v1` D1 发的发现同款（门/权限文案混用），这条取舍我没资格替票定。
5. **`197-r3` 串行范围**：票面 :39 说与 `197-r3` 串行（同动 `subagent_selfapproval_197_test.go`）。本腿判甲形两枚新钉**不需要动那枚文件**（判据①在 approval 包自造卡，不借 197 载具）——若编排者想复用 197 的载具改 :109 那枚常量来造两枚活卡，就撞进串行窗，别派。

## §4 量不到的格子

1. **所有"今天绿/今天红"的动态读数**：本腿禁跑 Go，全部转引（`197-v1` §1.1 D2 发、§1.4 M1a-sweep、§1.5 在册四红）。编排者派单前按票面 :40 的规矩自己现跑 `go test ./internal/agent/approval/ ./cmd/wisp/ -count=1` 抄当日用例名——本腿给不了当日 PASS 名册，这是预检与派单的分工。
2. **`internal/panel` 4 枚在册红的当日状态**：`TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourWayAgree`。本腿只核到"它们在 `197-v1` §1.5 与 `145-growth-r3` :117 两份读数里同名在册"＋成因（前端键名册不对账／`design/assets` 被删），**今日 10-0x 是否仍在册没量**（§0 名册只覆盖 internal/cmd/tools 的 git 状态，不是测试态）。`internal/ball` 1 枚同理。
3. **frontend 那半的 `panel.ts` 键名册现状**：`TestApprovalCardViewJSONKeysMatchFrontendTypes` 红的具体差集（Go 多了哪些键、ts 多了哪些）本腿没读 `frontend/src/lib/panel.ts`（禁令）。判据②不受影响（对象是 PanelItem），但若 AC#3 判"甲形不够、要连出向 JSON 键一起钉"，那一格的现状数字要界面那支的腿来量。
4. **乙形的确切改动面**：§1.3 说乙形要动三处（registry/词表/AST 假设），这是静态推演；真正动冻结件的 diff 大小、以及 `TestPlantedGrantWiringGoesRedInASnapshot`（:2168）里六枚 plant 对新射程的反应，只有写了才知道——本腿不预支那个数字。
5. **判据①两形（同一块石头值不值两格）**：票面 AC#3 把这刀给验收腿。本腿只补一枚事实：判据①与判据②共享的是**"无尺"这个状态**（D2 发、M1a 发各证一半），实现上无共享代码——两枚新钉、两个包（甲形下都在 approval 包）、两套夹具，"一块石头"的判法在验收腿，不在本预检。
