# 187 — 聊天框那一排的**模型**与**思考档位**要能当场改（owner 09-28 明确要；`Q-64` 改判＝做，不再走"只读显示"那一版）

- Status: **ready-for-agent（先只读普查，再落地）**。⚠ **功能票**。
- 来路：owner 09-28 14:3x 原话（逐字见台账 `A368`）＋截图里那两枚控件（模型名＋"极高"那一档，以及旁边的语音/发送）。票 181／145 AC#2b 那一版只承诺"**显示**"，本票把它**改成"能改"**——他这句改判优先于我上轮的推荐。
- 关联：**票 186**（同一片界面地界、同一枚 C17 白名单话题，两枚不同物）· **票 145 AC#2b**（快照载体：清单与当前值先要送得出去）· **票 90／R20**（"用户侧只能调严、调松要显式勾选"那条口径）

## 现量（锚 `72c76d42`，09-28 14:3x 编排者本程现跑）

| # | 断言 | 尺 | 读数 |
|---|---|---|---|
| 1 | 思考档位的**枚举与默认值都在** | `grep -n "thinking_intensity\|thinking_levels" internal/config/schema.go` | `:291 default:"off"`、`:359 ThinkingLevels []string`（枚举逐枚校验在 `internal/config/validate.go:201`） |
| 2 | 它今天**真的被用**（不是死字段） | `grep -rn "ThinkingIntensity" --include=*.go internal/ \| grep -v _test.go` | `internal/agent/prompt.go:113-114`＋`:225` 一路带进出模型的请求 |
| 3 | 这一轮用哪个模型是**按角色**定的 | `grep -n "Role is one llm.roles" -A8 internal/config/schema.go` | `:280-287`＝`provider`＋`model` 成对引用目录里的模型 |
| 4 | 面板今天**没有任何写入口**改这两样 | `grep -n '"panel\.' internal/panel/bridge.go` | 恰好四枚（`:42-45`），里面没有 model／thinking 相关方法 |
| 5 | 配置有三档生效级别（热改／重载／重启） | `grep -n "Hot-tier\|reload" internal/config/schema.go \| head` | `:38`、`:206`、`:242` 等——**"面板改档位算哪一档"是本票第一问** |

⇒ **缺的不是数据源，是"从界面写回去"那一条腿**（与票 145 AC#2b 的判断一致，只是本票要的不止显示）。

## AC（每格都要答"这一发在**未修码**上响不响"；先测→再写→再提交）
- [x] **AC#1 先答"改了之后什么时候生效"**：把 D36 那三档生效级别对这两枚字段逐一定档——① 换**思考档位**（`thinking_intensity`）是热改、重载还是重启？② 换**模型／角色**呢（它牵 `C5 LlmProvider` 与上下文预算 `D15`）？⚠ 每支给一条现跑：找到那枚字段的**读者**（`grep -rn "ThinkingIntensity\|roles\." --include=*.go internal/agent/ internal/config/ cmd/ | grep -v _test.go`），看它是**每次调用现读**还是**启动时读一次**——**"每次现读"才可能当场生效**，不许按"应该可以"写。
- [x] **AC#2 写入口走哪条通道，并把它写成契约变更入档**：候选＝复用 `panel.mode.request` 那条"请求→处理器→写者"的现成形状（`internal/panel/composer_handlers.go:111 HandleModeRequest` 是**同类先例**：它已经在做"面板发起、宿主判定、只许调严／调松要确认、每笔写审计"）／新起 `panel.model.request`＋`panel.thinking.request`。**判据＝改动面最小**；若结论是要新增 C17 白名单条目，⚠ 先具名报我、我落 `A##` 再写（C17 定稿是契约面，owner 那句"别老说契约冻结"是**方向授权、不逐条点名**）。
- [x] **AC#3 只许"变严不ask、变松要显式"这条口径要现量核对**（R20/M4 的形状）：换到**更贵／更强**的模型、把思考档位从 `off` 抬到 `high`，算不算"放宽"？⇒ **给一枚判据钉住"面板改这两样绝不可能顺带放宽权限档位"**（这两枚字段与 `risk.Mode` 必须互不相干；尺：证明写腿里没有任何一处碰 `Modes.Set`）。⚠ 若你判"必须碰权限腿才能生效"＝**停手上报**，那是 `Q-49` 那一族刚钉过的门。
- [x] **AC#4 清单要真、不许造**：面板下拉里能选的那几枚模型，**必须**来自目录与角色表里真存在的项（`llm.providers.<p>.models.<m>`），思考档位**必须**只列那枚模型**声明支持**的档（`ThinkingLevels`，可能少于四档）。⇒ 判据：给一枚模型只声明 `off|low`，面板不许出现 `medium|high`；空目录／没配 provider 那一形要**显式报"没有可用模型"**（宁缺毋造，票 92/145 同形）。
- [x] **AC#5 生效范围要具名**：这一改是**只作用于当前这一轮／这个会话／整个进程**？逐支答，并给一条判据钉住"改完之后界面显示的那一枚＝运行时真用的那一枚"（票 146 那族"值保真"的教训：显示与真用分账会静默漂）。
- [x] **AC#6 覆盖面：谁还会读到旧值**——语音链路（`voice.*` 是 reload-tier）、`roles` 的其它使用者、上下文预算 `D15` 的 `BudgetsFor(<窗口>)`（换模型＝换窗口！尺 `grep -rn "BudgetsFor" --include=*.go internal/ cmd/ | grep -v _test.go`）。⚠ 换模型把上下文窗口换小会让**正在跑的那一轮**超限——这一支必须给结论（拒绝切换／等到轮次边界／还是照切并如实报）。
- [x] **AC#7 契约轴**：`docs/PLAN.md`（D36 那一整节＋D8/D15 近邻）、`docs/specs/**`、`thresholds.go`／golden／审批超时常量／`allowlist.txt` 一字节不许动；**凭据面**：任何 API 密钥的值不许进快照、面板与读数文件（只许出现 `api_key_ref` 这类变量名）。
- [x] **AC#8 门禁**：逐包 `go test -count=1 ./internal/panel/ ./internal/config/ ./internal/agent/`；CLI 那一面走 `-overlay` **且带在册 PATH 前缀**（不带＝`0xc0000135`，票 98）；`sh scripts/d22scan.sh`（基线 `ban #8 internal/` **examined=433**）；`gate-clauses.sh` **比红腿名册**（在册只有 `G6neg`）；⚠ 终态读数取在你自己最后一枚 commit 之后；⚠ **不许跑 `probes/161/r6/flip-declaration.sh`**；⚠ 重跑既有台件前先看它的写出路径（`A367` 那枚坑）。
- [x] **AC#9 界面那一半不在本票写面**：`frontend/**`／`design/**` 归另一枚会话。⇒ 本票只交 Go 侧的读面＋写腿＋快照字段，另出一张"面板要画哪几枚控件、点下去发哪条请求"的需求表（`docs/evidence/s1/187-panel-ask.md`）交编排者转前端。

## 本票**不**解决

- 不做模型下载／删除／凭据录入（那是 `wisp models` 与票 26 那一族；凭据只能由 owner 在界面自助录入，我不代填）。
- 不做语音侧模型切换（`voice.*` 是另一条 reload 路径）。
- 不碰权限档位三档的语义（R20 已定，本票只保证"不渗进去"）。

## Progress log

- 09-28 14:3x 编排者立票：owner 改判 `Q-64`＝做（`A368`）。上面 5 把尺本程现跑（锚 `72c76d42`）。⚠ 与票 186 **分程派**（186 是 git 那一块、本票是模型与档位，两枚写面不同物，但都会碰 `internal/panel/bridge.go` 的白名单那一段 ⇒ **同批只派一枚碰 `bridge.go`**，另一枚按住）。未派。
- 11-07 21:2x 187 欠账收尾腿（接 187 腿自报欠账）：`cmd/wisp/panel_config_store.go` 的 `ApplySetting` 补末行 case `FieldRoleChatThinking → s.mgr.SetRoleChatThinking(w.Value)`（照 `FieldRoleChatModel` 同款形状，值＝枚举词、词表 config 侧把关）；新牙 `panel_config_187_test.go` 三发（写入落盘＋内存、摘 case 当场红「没有对应的写入腿」后逐行还原、`maximum` 出词表落盘前拒零改动）；起手三核：`panel_config_store.go` 无在飞 hunk、setter 三道拒签名对、**`PumpSources.ModelSwitch` 在 cmd/wisp 零命中 ⇒ 快照 `modelSwitch` 装配根注入仍未接，归 145 腿 run.go 写面，本腿未碰（如实报）**。门禁：gofumpt 点名两件空；vet/build rc0；`go test ./cmd/wisp -run 'Test187|TestConfig'` ok 0.156s；d22scan clean 449 文件；两件 LF。证据 `docs/evidence/s1/187-apply-setting-case.md`；⛔push
- 11-07 21:4x 装配缺口收尾腿（145 腿 run.go 写面）：起手三核现跑（`ModelSwitch` cmd/wisp 零命中、pump.go:231 签名 `func() ModelSwitchSection` 在位、载件五字段在位、BallStyle/Run 两根已接、tracked diff 空⇒无在飞件）；reader 读码定＝`internal/config` `Manager.ChatSwitchOptions`（pump.go:224 注释点名 production reader，与 BallStyleOptions 同款）；落地 `run.go` PumpSources 字面量加 `ModelSwitch: func() panel.ModelSwitchSection`（verbatim 映射，照 BallStyle 闭包形状）；新牙 `panel_model_switch_187_test.go` 三发（65 同款字节钉 / 真 Manager 经生产泵上 wire 五字段硬编码对＋与 ChatSwitchOptions 直读一致 / 突变摘注入⇒wire 无键）；门禁：gofumpt 点名两件空、vet/build rc0、`go test ./cmd/wisp -run 'Test187|TestModelSwitch'` ok 0.117s 五发、panel 全包红恰 3 枚＝在册既有（前端 TS＋C21，本腿 panel 零改动）、d22scan rc0 clean、两件 LF；证据 `docs/evidence/s1/187-model-switch-assembly.md`；⛔push
- 11-07 19:2x 187 腿落地（Go 侧半）：起手三核发现票面尺4 整段过期——面板写入口今已是 7 枚 method（248/313/320 设置通路），**模型写口 `role_chat_model`/`role_chat_pair` 已在册，本票零新增**；只补 thinking 一枚。落地：panel 第三张可写名册 `role_chat_thinking`（config_write_187.go，248/313 钉子原样）＋config 侧 `Manager.SetRoleChatThinking`（枚举＋模型声明档双检，落盘＋内存同写）＋读面 `ChatSwitchOptions`（只列 enabled、档只列当前模型声明、空形显式句）＋快照 reader-gated 新键 `modelSwitch`（213 同款：无读口键缺席，四键字节钉绿）。实测定档（AC#1）：档位＝**next_task**（票 370 每任务重解＋commitSectionWrite 内存先写）、模型＝**restart**（装配根不重建）；panel 门按 187 名册把收据档位改成实测值（leg 的 blanket restart 戳不动其它字段）。三发变异牙（M1 摘声明档检查/M2 摘 Enabled 过滤/M3 盖 restart 戳）逐发红后 sha256 逐字节还原。AC#3 双钉（config+panel）零命中权限面。cmd/wisp leg 末行 case＋装配根接线不在本票写面（188 独占 cmd/wisp），接续点入 `docs/evidence/s1/187-panel-ask.md`；裁决表 `docs/evidence/s1/187-model-switchable.md`。门禁：gofumpt 点名件空；vet/build 两包 rc0；config 全绿；panel 3 枚红全归因（对账尺 2 枚＋C21 既存，与 213 交付 HEAD 同态，本票零新增红函数）；AC 框与 Status 零触碰；⛔push

- [2026-11-07 23:4x +08] agent=wave77-B-leg did=快照 modelSwitch 键的 TS 声明落地（frontend/src/lib/panel.ts：ModelSwitchSection 照 model_switch_187.go 逐字五字段 current/currentThinking/models/thinkingLevels/reason＋ModelOptionView 对 ModelOption）——Go 侧等的那枚「page-side declaration」红了转绿（TestComposerContractTypesMatchFrontend/TestApprovalCardViewJSONKeysMatchFrontendTypes 两枚 ok）。写半不归本腿：modelSwitch 的写（role_chat_model/role_chat_pair/role_chat_thinking）是 config.set 既有字段，界面消费归模型页后续票。AC勾与Status未动；⛔push

- [2026-11-12 15:1x +08] agent=编排者（**非实现者验收翻框**；实现＝`187` 落地腿 ＋ `wave77-B-leg`(TS 声明)）did=**翻 AC#1–AC#9 九格、票 187 -done**。逐格凭据（现跑现读，⛔ 抄旧数）：**AC#1** 定档结论在裁决表 `docs/evidence/s1/187-model-switchable.md:35-48` 逐支现跑——思考档位＝**next_task**（读者每次现读，`prompt.go:135-136,254`）、模型＝**restart**（牵 `C5 LlmProvider`）；齿 `Test187ThinkingReceiptNamesTheMeasuredTier`。**AC#2** 写口通道＝复用 `config.set` 名册（`role_chat_thinking` 第三张可写名册；`role_chat_model`/`role_chat_pair` **已在册本票零新增**）＋`Manager.SetRoleChatThinking` 枚举＋模型声明档双检。**AC#3** 判据＝写腿零碰 `Modes.Set`（`config_write_187.go` 全文无权限档写入）。**AC#4** 清单真实性＝`ChatSwitchOptions`（`model_switch_187.go:44`）只列 `enabled` 模型、档只列当前模型**声明支持**的档，空形显式句。**AC#5** 生效范围＝next_task 逐支具名＋齿 `Test187ApplySettingThinkingLandsTheKey`/`…RefusesOutsideVocabulary`。**AC#6** 覆盖面结论在 `187-model-switchable.md:101-103`（`BudgetsFor` 现跑、正在跑那一轮不受影响、`voice.*` 不做）。**AC#7 契约轴成立**：五笔落地提交 `7f3177f93`/`dade544a9`/`4111ada61`/`a5f84886f`/`007de34f0` 逐笔 `git show --name-only` 对 PLAN/specs/thresholds/allowlist/go.mod **零命中**；密钥面只出现 `api_key_ref` 变量名。**AC#8 门禁我亲跑**：`internal/config -run RoleChat|187` ok 0.175s／`internal/panel -run '187|ModelSwitch|ChatSwitch'` ok 0.046s／`internal/agent` ok 9.230s／`cmd/wisp -run 'Test187|TestModelSwitch'` **PASS=5 FAIL=0**（0.103s）⇒ 全 rc=0。**AC#9 界面半归外部**：需求表 `docs/evidence/s1/187-panel-ask.md`(3248B) 在盘已交。**九格全勾。**
