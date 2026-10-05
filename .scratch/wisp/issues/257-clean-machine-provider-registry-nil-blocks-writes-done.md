# 票 257 — 全新机器上设置页**七枚字段一枚都写不进去**：首份默认配置里服务商注册表是 nil，而写侧要求"行已经存在"；名册里又没有"可写 provider"那一枚字段

**立票时刻**：2026-10-02 13:0x +08，锚点 HEAD `cf87f7a2`（`dev`）
**来路**：只读腿 `e2e-panel-1` 的端到端跳表 `.scratch/wisp/probes/e2e-panel-1/readiness.md`（174 行／39,188 字节，八节齐、占位符 0）§3 的 **H11**；台账 `A531`。⚠ 下面五行都是**待验断言**，落地腿动手前逐条自己复认。

## 现量（腿自跑，编排者抽验了第一行）

1. `defaults.go:77-78` 对 map 的处理逐字是「leave nil」⇒ 首份默认配置里**服务商注册表是 nil**；`schema.go:419` 那一枚也没有 `default` 标签。
2. 写侧要求行已存在：`internal/config/settings.go:310-314`／`:294-308` 那两处对"要写的键所属行不存在"是**拒**（不是建行）。
3. 剩下的 `role_chat_model` 撞 `catalog.go:97-101`：model 有值而 provider 为空 ⇒ 拒。
4. 而**可写名册里没有"provider"那一枚字段**（名册七枚字段里没有一行能创建服务商条目）。
5. ⇒ 净结果＝**干净机器上，面板那七枚字段一枚都写不进去**（这与 H2 是两枚独立的死结：H2 是"文件根本没建、面板打不开"，归票 198；本格是"文件有了也写不进去"）。

## 要建什么

- [x] **AC#0 先定"哪一形对机主是诚实的"（不许直接开写）**：ⓐ 首份默认把注册表建成**一条可用的空壳行**（写侧就能改它）；ⓑ 名册里**新增一枚"选/建服务商"字段**并让写侧接受"行不存在则建"；ⓒ **不建行**，改为回执明确说"请先在配置文件里添加一节 `[providers.x]`"。三形各写：要动哪几枚文件／会不会新增 C17 面字段（ⓑ 会）／机主在界面上看到的那句话变成什么。⛔ 普查腿不许改产码，写点只准落在 `.scratch/wisp/probes/257/a1/census.md`。**〔17:4x 编排者翻勾：`257-a1` 交件 `ed3fd270`（92 行占位 0、禁 Go 遵守、四环现量没被推翻——仅一处行号差 1：case 头 `catalog.go:97-100`、真拒句 `:101-103`）；选形＝ⓒ，裁定在下面第 8 节。⚠ ⓒ 那格票面拼法被普查腿顶回：顶层 `[providers.x]` 会撞 `internal/config/parse.go:72/123` 的 `DisallowUnknownFields` ⇒ 面板整条链起不来；真路径＝`[llm.providers.<名>]`——票面原句不改，以此为准〕**
- [x] **AC#1 干净机器那一发必须真跑**：临时数据根＋**没有** `config.toml` 的起步态，建完首份配置后，逐枚试写那七枚字段，**终态＝至少 ⓐ／ⓑ／ⓒ 里被选定的那一形真兑现**（⛔ 不许用"手工先塞一份完整 config.toml"当前提——那正是今天所有测试都在做的事）。
- [x] **AC#2 失败要说清是哪一种**：拒写的三种原因（文件没建／行不存在／校验不过）**各配一句不同的话**，⛔ 不许折成一句"配置未生效"（本仓票 223 AC#4 那条定式）。
- [x] **AC#3 凭据面不许动**：⛔ 不新增任何"把值回显给页面"的方法、不落明文、凭据仍只走现成那一枚存储；对话里绝不出现凭据值。
- [x] **AC#4 越界检查**：`git diff` 出现 `frontend/**`／`design/**`／`PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt` 任一路径 ⇒ 退回。

## 禁区

⛔ 名册新增字段属 **C17 契约面**（ⓑ 那形）⇒ 先停手，由编排者落具名 `A##` 才许动；⛔ 三枚冻结件与上面那张冻结清单一字不动；不许为变绿放宽断言、不许 `t.Skip`、不许把 SKIP 读成通过。Git：只 commit 不 push、显式 pathspec、禁 `add -A`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`、临时件只建不删。

## 排程

写面＝`internal/config`（＋可能 `internal/panel`）⇒ ⛔ 与 `255-r1`（同写 `internal/config`）**串行**；⚠ 与票 198 的 `198-r2` 是同一条链上的前后两跳（198 先把"文件建出来"接上，本票才有干净机器可谈）⇒ **排在 `198-r2` 之后**。默认动作＝**先派 `257-a1`（只读 AC#0 三形代价）**，⛔ 未量不裁。撤销口令「257 撤」。

---

## 8. `257-a1` 收档 ⇒ AC#0 翻勾＋**选形＝ⓒ**（编排者裁定，账 `A543`，2026-10-02 17:4x，锚 `5c8fa195`）

**1. 交件核过**：`ed3fd270`（72/12）＝census **92 行／22,091 字节**、占位 0、禁 Go 遵守；四环现量没被推翻（`defaults.go:77-78`／`schema.go:419`／`settings.go:310-314`＋`:294-308` 双层关／名册七枚无建行入口），一处行号差 1 具名（`catalog.go:101-103`）。我抽验 ⓒ 的最硬坑属实：`internal/config/parse.go:72/123` 真有 `DisallowUnknownFields()`。

**2. 三形判语（读数全在 census）**：
- **ⓐ 空壳行＝毙**。要动 `defaults.go:77-78` ⇒ 撞 **4 枚钉**（票 198 两枚逐字节钉＋`firstrun_198_test.go:220` 逐字禁首建文件出现 `[llm.providers`＋`settings_248_test.go:130`）；非预设名被 `validate.go:171-177` 拒 ⇒ **首建文件自己加载不过、进程起不来**；只解锁 3/7。
- **ⓑ 名册加字段＝毙（本轮）**。C17 现量 6 枚方法名；新字段落 `SettingsView` 顶的是信封新键钉 `TestAC1SettingsKeysAreTheOnlyNewOnesOnTheSharedEnvelope`；最硬的坑＝**建 provider 天生多键写**（名＋protocol 至少两键）撞票 248 AC#11"一键一写"（票 226 地盘）。要做得单开 C17 批准＋多键写裁定，两件都是契约面，本轮不做。
- **ⓒ 指引手加＝选它**。真路径 **`[llm.providers.<名>]`**（⛔ 票面原拼法 `[providers.x]` 会撞 `DisallowUnknownFields` ⇒ 链起不来——这处是普查腿顶回我票面，记我）；静态链＝手加 provider 行＋model 子表＋`roles.chat.provider` 点名三样 ⇒ **7/7 解锁**（只加 provider 行 3/7）；产码面最窄＝只动 `cmd/wisp/firstrun.go` 文案＋三枚回执钉（7 必名串／cause 句隔离／二跑不重复），与票 198 刚结案的回执形状同族。

**3. ★ 选形＝ⓒ（五样齐）**：**文件**＝`cmd/wisp/firstrun.go`（回执文案）＋对应测试；**理由**＝ⓐ 撞钉＋自杀形、ⓑ 两件契约面、ⓒ 零产码结构改动＋owner 自己能走通＋7/7；**边界**＝① 回执必须分清**两条通道**："热加载认手改（run 腿 1s tick）"与"设置页写入要重启"并存且都对，不许混成一句（census §ⓒ 诚实边界）；② 写侧那道"行不存在拒写"**不动**——机主按指引手加后，设置页那 7 枚就落在已有行上可写了，这正是"诚实"所在：界面不替他建行，但告诉他怎么建；③ ⓑ 的两件契约面（C17 新增＋多键写）**单开一张票**留给"哪天真要从界面建行"那一天，本轮不立（owner 没要过"界面建行"，他要的是"能配上模型"）。**撤销口令「257 改形 ⓐ」**。

**4. 落地腿 `257-r1` 判据预告（暂不派——写面 `cmd/wisp` 被 `253-r2` 占着，包级互斥）**：AC#1 干净机真跑（选形 ⓒ 兑现＝回执给出 `[llm.providers.<名>]` 三样指引）／AC#2 三句不同的拒写原因各一句／AC#3 凭据面不动／AC#4 越界。写面＝`cmd/wisp`（firstrun.go 文案＋测试），⛔ 不碰 `internal/config`。

**5. 同族生产死线收档（10-03 09:5x，只读腿 `257-a2` 257 行／54,036 字节最末 `cbf4f1b2`；★死线那三行由编排者本人复量坐实，落账 `A560`）**
- ★`config.NewManager(path, res)` 的**第二参数产线三处全传 `nil`**（我复量逐字＝`cmd/wisp/run.go:409`／`cmd/wisp/panel_inbound.go:230`／`cmd/balldebug/main.go:231`；`internal/ball/hotkey_reload.go:15` 那处只是注释示例）⇒ `internal/config/loader.go:46` 的 `resolveRefs` **生产里永不执行**、`(*Manager).Resolved()`（`manager.go:123`）**非测试调用者 0 枚** ⇒ 配置层那条"引用形凭据"解引用通道**今天没接**。⇒ 对本票的直接影响：**回执⛔ 不许把那通道列成入口**（列了＝指一条走不通的路）；凭据真走的是另一条（`run.go:435-436`→`internal/llm/resolver.go:141`，每建一次 Endpoint 解一次）。另具名登记：`voice.realtime.api_key_ref` 今天**没有运行时消费者**。
- 两处留〔腿报，未复核〕：`internal/secret/migrate.go:86` 明文→DPAPI 那台机零生产调用者；**写门读盘／读面看内存**（`internal/config/settings.go:205`→`readCurrentFile :274-289` vs `cmd/wisp/panel_config_store.go:92`）＝owner 那句"我填了为什么没生效"的反向形状，⚠ 属票 248 那一族不吞进本票。
- ★**它就地打过自己一处，这条对我有反噬**：常驻进程**并非必然有两枚 Manager**——`cmd/wisp/resident_task_source_windows.go:224-231` 那道 `interactiveStdin()` 闸没过时**该进程连一个重读盘的东西都没有** ⇒ 我在 `A542`／票 248 面上用过的"手改＝热加载认／面板写＝要重启"**两格措辞本身是半谎**，登记为过期；**落地回执要按三形状写**（有控制台／无控制台／面板写入）。撤销／更正都不改上面原句，按台账追加规矩走。

---

## 9. 收 `257-v1` 终裁 ⇒ **AC#1–AC#4 全勾＋改名 `-done`**（编排者翻勾，账 `A621`，2026-10-05 12:5x +08，裁决表 `docs/evidence/s1/257-clean-machine-provider-registry-v1.md`＝262 行／39,919 字节／占位 0，四笔自落 commit `65142d9a`→`0793e35e`→`3862e0e2`→`14dc6f48`）

**1. 凭据形状（我核的是盘，不是通知）**：五发突变 M1–M5 全种在 `cmd/wisp/firstrun.go` 一枚文件、每发跑完立刻还原、`md5sum` 五取相等＝`0491339282492f2cabdbf5be576c8a57`（逐字＝`git cat-file blob HEAD:cmd/wisp/firstrun.go`，我 12:5x 复算过同一值，两把尺同发各取一次相等）；门禁四门它自跑（d22scan rc=0 clean／path-length VERDICT GREEN／`go vet ./cmd/wisp/` rc=0／gofumpt 只列预存的 `cmd/wisp/models.go`）；票面 AC 框与台账**一字未改**（两枚文件不在它任何 pathspec 里）。⇒ 这四格从此有非实现者判语，上面第 8 节第 4 行那句"落地腿 `257-r1` 判据预告（暂不派）"至此走完。

**2. 四条判语（逐格一句，边界随格一起读，⛔ 不许摘成半句）**
- **AC#1 成立**：干净机**两形都真跑过**——包内八枚 rc=0 并复现逐字 `map[第 2 种拒因：行不存在:7]`，真进程那发（它自己搭的临时数据根＋`portable.txt`）rc=2、stderr **5,674 字节**里三样指引各到一次、生成的 `config.toml` **2,571 字节且 `grep -c "llm.providers"`＝0**。⛔ 没有一枚用例拿"手工先塞一份完整 config.toml"当前提（起步态是断言出来的：先 `os.Stat` 必 `fs.ErrNotExist`＋退码仍 2＋stderr 有 `新建默认配置`，四道任一不满足就 `t.Fatalf`）。**边界两条**：① 名册第 7 枚 `provider_credential` 的**值腿**（`StoreCredential`→DPAPI 真密封）两枚测试都没走、走的是引用侧 `SetProviderAPIKeyRef`（`firstrun_257_test.go:127` 那一行逐字写着 `"Manager.SetProviderAPIKeyRef (StoreCredential's config-side leg)"`，注释自注"枚数 7、腿数 6"）；② 七枚**在设置页上的长相**不在本票写面（§8 边界③）。
- **AC#2 成立**：三因三句各带互不重叠的补救（tag 现读 `internal/config/settings.go:78-80`＋各自真身分支 `:242`／`:293`／`:367`＋`:379`；回执三句 `cmd/wisp/firstrun.go:162-171` 与三枚 tag 同串）；**把第 2、3 句折成同一句 ⇒ 指名用例红，四条红句逐字在裁决表 §4**（M4，12:36:08→12:36:17）。另一格"什么时候算用上"分三形状由 `…AC2EffectTimingSaysThreeProcessShapes` 钉着，折成"重启就好"同样被那把只许 `配置未生效` 出现一次的计数尺拦着。
- **AC#3 成立**（附一条准确性 nit）：范围内 `git diff -G"StoreCredential" --name-only 21bec8a1..HEAD` **产码零枚**，`cmd/wisp/panel_config_store.go`＋`internal/panel/config_handlers.go`＋`internal/panel/bridge.go` 三枚契约文件命中 **0**；`firstrun.go` 内函数枚数仍 **1**；**往凭据那句里回显 `api_key = "sk-…"` ⇒ 唯一红枚**（M5）。nit＝`firstrun.go:183` 那句"只说已录入还是没录入"比盘上窄（`config_handlers.go:459-461` 还会印 `引用 %s`＝引用名不是值），⛔ 不是值泄漏、⛔ 不改那 61 行——登记为残余④。
- **AC#4 成立**：八枚禁区模式除 `internal/config/**` 那 **5 枚全部逐笔归票 267**（`37f8e5c6`／`ec6a47a8`／`873c3063`，subject 逐字起于 `ticket 267 r1: band [risk].confirm_timeout_sec at load (shape a)`）外全 0；本票各腿自己的七笔 commit 枚枚 0 命中。**范围外一处落点差照旧具名**（`8e443ed7`＝`257-r1b` 动过 `internal/config/{loader,settings}.go`＋新增 `settings_257_test.go`，早于本范围锚 `21bec8a1`，`A606` 已收账）。

**3. ★三格加裁（这三条是三枚派单问题，不在原 AC 里，一并记在票上）**
- **两测试套＝两套都留、同批算交件**：决定性那发是 M2——只抹 `firstrun.go:156` 半句 ⇒ **只红 `…ReceiptStatesTheNonPresetCondition` 一枚，六枚那一套全绿**；反发 M4／M5 只红六枚那一套。真重叠只有 **4 枚断言**（`firstrun_257_test.go:425-428`／`:406-418` 与 `firstrun_257_nonpreset_test.go:145-148`／`:151-153` 同源），合并省 4 枚断言的代价是把两套唯一共用的 helper 变成跨文件依赖＝正是 `A617` 那两条事故要避开的形状。⚠ 一根脆针具名：`firstrun_257_nonpreset_test.go:72` 那枚 needle **单独使用时没有牙**（票 198 旧文案 `firstrun.go:116-117` 里有同一串），别把它当牙用。
- **正向断言不恒真（M1 把 `Fprint` 目标换成 `io.Discard`、字串一字不动 ⇒ 红 3 枚）**；★**口径更正：那一格真到达的通道是 stderr，不是 stdout**（`main.go:159-163` 逐字 `stdout: os.Stdout, stderr: os.Stderr` 交给 `runTextTask`；真进程 stdout 只有版本行 **119 字节**，三样指引全在 stderr 那 5,674 字节里）。⇒ 本票 AC#1 按"到达用户可见通道"裁＝成立，**今后引用这一格不许写 stdout**。
- **那 61 行没有把指引漏进生成的 `config.toml`，而且漏了必响**：实弹 M3 往生成文件追加 `[llm.providers.t257v1leak]` ⇒ **红 6 枚、三把尺同响**（票 198 名册反控＋票 198 逐字节尺 `want 2571 got 2599`＋本票 AC#1 泄漏反控），并顺带把 §8-2 里 ⓐ 那形"首建文件自己加载不过"复现成红句。

**4. 那枚红＋那枚命名跳＝机器争用，本票新增红 0 枚**：它整包 `RUN 331／PASS 234／FAIL 1／SKIP 0`（12:22:3x→12:29:07），与实现腿那一发**顶层名集 `diff -q` 完全相同**，只差三枚颜色（lifecycle FAIL→PASS、latency SKIP→PASS 且真带 `n=1` 样本、`TestAC14GoSideEvalPushReachesThePage` PASS→FAIL）；三枚都有低载隔离复跑顶着（单发＋`-count=3` 共 5 个冷启样本 917/980/960/757/802 ms，无一越过 1500 ms）。⇒ **⚠ 同一发整包的窗口里 `259-r2` 在 `internal/agent/approval` 落了笔**（前后 status 它逐字申报：12:21:49 那两枚文件 0 行 → 12:30:36 `M approval.go`＋`M queue.go`，109 增／13 删、读片段确为注释级）；它的判语是"变色那三枚的判据不读 `approval`"，我认这个理由，并把"若不作数需要复跑的是包内顺序能否复现 AC#14 的红、不是票 257 的六格"这句一起留在账上。

**5. 残余四格（逐枚具名归口，⛔ 不吞进本票、也不因结案而消失）**
- ① **`provider_credential` 值腿**（`StoreCredential`→DPAPI 真密封在真机上的长相）＝`docs/evidence/s1/` 欠一行真机裁决，与 `A617` §1 残余第一条同格；⛔ 不许为测试写真 blob。
- ② **设置页上的完整长相**（信封键／回执串怎么渲染）＝归**票 248 那一族**，要从界面建行＝单开票（§8 边界③，两件契约面：C17 新增＋多键写）。
- ③ **stderr 通道没有 `exec` 级用例长期钉住**（今天两枚测试都在同进程传 buffer，它那发真进程读数不可重放为门禁）＝**归票 244**（同一族"从终端起这进程到底看得见什么"），要动＝一发行文＋一枚 exec 级用例。
- ④ `firstrun.go:183` 那句"只说已录入还是没录入"比盘上窄（还念引用名）＝纯文案准确性，⛔ 不夹进别的票顺手改；低利害、可后并。

**6. ⚠⚠ 幻影用例名（记在票上，不改别人的件）**：`9c0d4c9a` 那份合并件 §3 逐字描述了四枚用例名（`TestTicket257R2AC1CleanMachineRefusesAllSevenFields`／`…FollowingTheReceiptUnlocksAllSevenFields`／`…ReceiptNamesTheThreeHandAddedThings`／`…SecondRunSaysTheGuidanceOnlyOnce`，外加 §3 末行那枚省略号前缀 `TestTicket257R2AC0`）——`257-v1` 现量**工作树 0 命中、`git grep HEAD` 0 命中**，我这轮自己逐名复跑同样 **0 命中**，盘上真名是 `…ReceiptTeachesTheWalkableChain`／`…NonPresetRowNeedsProtocolBeforeItLoads`／`…ReceiptChainWalkUnlocksAllSevenFields`／`…ReceiptOnlyFirstItemUnlocksThreeOfSeven`／`…ReceiptStatesTheNonPresetCondition`／`…AC2EffectTimingSaysThreeProcessShapes`／`…AC2ThreeRefusalsStayThreeSentences`／`…AC3CredentialSurfaceUntouched`（八枚）。⇒ 处置＝**原件一字不改**（临时件只建不删、也不改别人的件，§2 末段已具名），**谁按 `9c0d4c9a` 那份件翻 AC 框就会给四枚不存在的用例记账**——本票四格的凭据只认 `docs/evidence/s1/257-clean-machine-provider-registry-v1.md`。这条与台账里那条"注释／证据件里的测试名一律当待验断言"同族，第五次生效。

**7. 结案动作**：AC#0（`A543` 裁形 ⓒ）＋AC#1–AC#4 ＝五格全勾；`git mv` 改名 `257-clean-machine-provider-registry-nil-blocks-writes-done.md`（⛔ 票面原句零改动，第 43 行那句"暂不派——写面被 `253-r2` 占着"、第 28 行那条排程都照原样留着当历史）。撤销口令**「257 撤勾」**＝把这四枚框降回 `- [ ]` 并恢复原名。
