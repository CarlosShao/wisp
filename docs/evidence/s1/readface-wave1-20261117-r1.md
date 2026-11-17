# readface-wave1-20261117-r1 — 产品页读面波 1（写腿 agent-102 · 五项裁决表）

派单：`.scratch/wisp/probes/orch/20261117/txt/readface-wave1-dispatch.txt`（五项）。
写面：`cmd/wisp/**`（唯一写腿）、`internal/panel/**`、`internal/agent/**`、`internal/config/**`。
⛔ 未动前端（`git status` 下 `frontend/` 零改动属其他在飞腿的既有条目，本腿未触碰）；⛔ 未动 thresholds/golden/allowlist/PLAN/specs/D43/design；⛔ 未 push；⛔ 未提交 dist 字节。

## 五项落点（file:line）

### 1. config.get 全节树读面（解锁 39 AC#1）
- 生产读面：`internal/config/sections_456.go:45`（`SectionRows`，反射按 toml tag 走全 D36 树，含 `[plugins]` 动态子表；tier 词全部取自 `TierRegistry`，锁定段答 `locked`）。
- 视图/线上：`internal/panel/config_handlers.go`（`SettingsView.Sections` ＋ payload `sections[]{key,value,tier}`，omitempty＝腿没答不上线）；装配：`cmd/wisp/panel_config_store.go:147`（`configSectionRows456`，copy-not-rewalk）。
- reload 事件（写后推）：`internal/panel/config_handlers.go:579`（`ConfigReloadEventMarker = "wisp-config-reload:"`；accepted config.set 收据追加一行 JSON `{event:"config.reload",tier,written[],sections[]}`，sections＝written 键路径首段去重排序；被拒写零事件）。事件走 awaited binding reply（248/A877 裁定的同一条页面可达面），并如实标注：本树没有第二个推送跳，快照推送不在此腿接。
- tier 词表＝D36 三档词（hot/reload/restart）＋锁定标记 locked。派单写的 tier∈hot/restart/nextboot 中 "nextboot"＝restart、"reload" 既是档也是事件；本腿**不改 D36 词表**（改契约须人工批），如实按 registry 的词答。

### 2. errors[]{line,message} 读面（解锁 39 AC#4）
- 结构化行：`internal/config/loaderrors_456.go:43`（`ParseErrorsFile`，重跑读管线内容半：版本探测→迁移→严格解码→validate）；`loaderrors_456.go:54`（`Manager.ParseErrors`，读自己的 path，不动 cur）。
- 行号单源：`internal/config/loaderrors_456.go:108`（`decodeErrorRows`：StrictMissingError 每枚未知键一行＋该键行号；DecodeError 一行＋line/col；validate/read 失败 Line=0＝"无行号"）；`internal/config/parse.go`（`formatDecodeError` 改为**由同一批 rows 渲染散文**，散文逐字节同旧文——包内解码测试钉住）。
- 线上：`internal/panel/config_handlers.go` payload `errors[]{line,message}`；装配 `cmd/wisp/panel_config_store.go:147`（`configLoadErrors456`）。

### 3. 成本聚合字段（解锁 40 AC#3/AC#4）
- 字段落点：`internal/panel/cost_handlers.go:54`（`CostAggregateRow`：day/月累计 tokens+micros、日/月预算、BudgetState、BudgetAllowNewTask、PricingVersion/PricingUpdated）＋ `cost_handlers.go:92`（`CostAggregateReader` 可选面——独立接口， days-only 旧腿与测试桩不破）；`CostDayRow.CachedTokens`（带 payload `cachedKnown=false` 门：**D35 表无 cached 列，本腿不改 internal/memory**，0 不是读数，门是话）。
- 装配：`cmd/wisp/panel_cost_194.go:94`（`CostAggregate`：cost_daily 真行求月累计；预算经 `observe.BudgetFrom/Judge`——基点整数比较，无浮点碰钱）；`panel_cost_194.go:159`（`budgetStateWord`：normal/warn/paused 三态，paused＝LevelOver 且 AllowNewTask=false）；计价版本＝`observe.DefaultPriceTable()`（pricing.json 自带 version+updated）。
- 接线：`cmd/wisp/panel_inbound.go`（`newCostReader194(dir, mgr)`）。
- 对账牙：`cmd/wisp/panel_cost_reconcile_40_test.go:40`（`Test40CostReconcileAgreesToTheMicroUnit`）——真 task_log 两行（Start/FinishTaskLog）＋真 BumpCostDay 两笔 ⇒ `observe.CostMeter.Reconcile` Agree=true 且 1000/1000 micro；再+1 micro 漂移 ⇒ Agree=false 且 1001/1000 双数点名。

### 4. 隐私路由（解锁 40 AC#1）
- 名册：`internal/panel/bridge.go:141`（`privacy.delete-one` / `privacy.purge` / `privacy.export` 三常量；无 `panel.` 前缀、无 grantRouteWords 词表成员；理由写在常量注）；guard 单行 case 同步扩。
- 派发表：`internal/panel/composer_dispatch.go:174`（`PrivacyRequestHandler` 三方法）＋ dispatch 三 case＋`Privacy` socket（nil＝具名拒）。
- 选择子：`internal/panel/bridge.go`（`ComposerRequest.PrivacyDomain` / `PrivacyItemID`，omitempty）。
- 处理器：`cmd/wisp/panel_privacy_40.go:89`（delete-one）/ `HandlePrivacyPurge` / `HandlePrivacyExport`（域词向 `memory.PrivacyDomains()` **问词表不抄词表**；per-read `memory.Open`；purge 必带域、无"全清"参；export 单行 JSON payload 走 `PrivacyExportPayloadMarker`，count 从 payload 自身 records 数读出——store 信封 Count 未填是 memory 层既成事实，本腿不虚报）；接线 `cmd/wisp/panel_inbound.go:411`。
- 仪器随动：`internal/panel/inbound_roster_253_test.go`（8→11、非前缀 4→7、roster 表、两枚 guardLine 植物行——该文件自己写明的"diff 必须亲手挪"摩擦）；`internal/panel/bridge_test.go`（`rosterPageGapLedger194` 三条注册，页字面量未落地、本腿不动 frontend——both-ways 审计照钉）。
- 三牙（计数对 DB）：`cmd/wisp/panel_privacy_40_test.go`——delete-one 3→2 且收据点名 id；purge 3→0 且收据点名 store 自报 3；export payload records=2 且收据"2 条记录"；另有 wiring nail（Production chain Privacy 实型断言）与拒绝两枚（未知域点名五域词表、无 id 不删库）。

### 5. BatchView 聚合卡字段（解锁 37 AC#5）
- 卡字段：`internal/panel/approval.go:81`（`ApprovalCardView` 增五枚**平铺 omitempty** 字段 batchTotal/batchDirs/batchFirst/batchDeferred/batchSummary——刻意不立新嵌套 struct，契约尺不需要新页侧 interface；非批次卡 JSON 与旧形逐字节同）。
- 载体：`internal/panel/pump.go`（`NativeVerdict.Batch`＋`NativeBatchView`；`:126 CardView` 映射，copy 不再算）。
- 数据源：`internal/agent/approval/cardsummary_37.go:39`（`CardSummary`：display 字段任意档 ≥3 目标都答，R4/R3 污染不出；**Aggregate 的确认合并契约原样未动**——L2 仍永不聚合成便宜确认，R7 升级守卫 `gate.go:317-346` 未动一字）。
- 白名单随动：`internal/agent/approval/pending_read.go:148`（panelVerdict 携 Paths，cloneBacking）；`ticket148_panel_whitelist_test.go`（Paths dropped→allowed，census 注释同步——该尺自己写明的"同轮评审做掉"）。
- 接线：`cmd/wisp/panel_pump.go:101`（nativeVerdictsFromGate 喂 CardSummary）。ledger：`internal/panel/composer_test.go`（`composerTSGapLedger["ApprovalCardView"]` 五键注册，判词写明"真线上事实、页半归前端腿"）。
- 牙：`internal/agent/approval/cardsummary_37_test.go`（10 项目标 total/first5/deferred5/dirs2；60 目标升级后 CardSummary 在而 Aggregate nil；2 项目标/污染全 nil）；`internal/panel/readface_wave1_test.go`（带 Batch 卡线上五键、无批次卡零 batch 键）；`cmd/wisp/readface_wave1_test.go:Test37VerdictBatchFieldsReachThePacket`（真 gate＋真队列 10 路径判 ⇒ 包里卡 batchTotal=10）。

## 红绿读数（自己复跑）

- `go test ./internal/config/ -count=1` → **ok**（含新牙 TestSectionRowsCoversEveryD36Section / TestSectionRowsTierWordsComeFromTheRegistry / TestSectionRowsCarriesWrittenValues / TestSectionRowsRefusesNilConfig / TestSectionRowsWalksPluginEntries / TestParseErrorsNamesUnknownKeysWithLines / TestParseErrorsNamesTypeErrorsWithLines / TestParseErrorsCleanFileAnswersNil / TestManagerParseErrorsReadsItsOwnPath）。
- `go test ./internal/agent/approval/ -count=1` → **ok**（含新牙 3 枚＋148 双尺照绿）。
- `go test ./internal/panel/ -count=1` → **ok**（29.2s 全包，含 253 closed-roster（11 名/7 非前缀）、114 镜像门（三条 ledger 注册 both-ways）、415/422 契约尺（ApprovalCardView 五键注册判定）、147/148 等全部存量仪器）。
- `go vet ./cmd/wisp/ ./internal/panel/ ./internal/config/ ./internal/agent/approval/` → clean。
- 载具 `scripts/wisp-cli-staged-tests.sh`（Git bash，跑前 `Get-Process wisp.test` 无在跑）：单跑新牙三批全绿——`-test.run 'Test40'` 8 枚 PASS（对账/三态/链上聚合/privacy 五枚）；`-test.run 'Test37VerdictBatchFieldsReachThePacket|Test39ConfigStoreReadFaceSectionsErrors'` 2 枚 PASS。
- ⚠ 整包 `go test ./cmd/wisp` 未跑（裸跑已知 0xc0000135 onnxruntime 影绑＋票 450 exit 139 既定项，⛔ 裸跑 A1888）；按 A1882/A1893 口径**载具单跑**本波新牙如实记录如上，整包留给跑窗。

## 没做完 / 卡点（如实）

1. **整包 cmd/wisp 未复跑**（只跑了载具单跑的新牙＋vet）：整包 exit 139 在 WebView2 区段是票 450 既定项，编排者跑窗复量。
2. **reload 事件走收据行**（`wisp-config-reload:` marker），不是快照新顶层键：加顶层键要前端 panel.ts 同步（本腿⛔动前端）＋pump 字节钉，按既有 judged-gap/两向尺口径留给页腿接的时候收口。收据是页面已可达的返回面（A877 裁定同源）。
3. **CACHED 无持久源**：cost_daily 无 cached 列（D35 表，internal/memory 不在本腿写面），payload 以 `cachedKnown=false` 门如实说"未聚合"，页渲染"未聚合"不得把 0 当读数；持久化需 D35 表变更＝人工批准面。
4. **export 信封 Count 未填**是 internal/memory 既有事实（不在写面），处理器从 payload records 自数，收据不虚报。
5. **privacy 三枚页字面量未落**（panel.ts）——已在 `rosterPageGapLedger194` 注册三条（理由在注），页腿落地同 commit 删除。

## 附：契约红线自检

- ⛔ C18/D36 语义零改动（tier 词表照 registry 原词输出）。
- ⛔ C17 白名单反向动作零（八→十一枚属**派单点名的新增**，理由如上具名；无删减）。
- ⛔ 未放宽任何断言：所有仪器随动均为其文件自述的既定程序（253 计数挪动、148 白名单挪列、415 gap 注册、114 ledger 注册）。
- goroutine：新测试仅用既有 gate/UI 面，无新生产协程（无需 observe.Registry.Spawn）。
