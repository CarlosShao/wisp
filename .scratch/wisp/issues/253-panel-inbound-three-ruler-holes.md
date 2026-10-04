# 票 253 — 面板入向那三枚尺洞：可达性钉只认 `postMessage(`、`config.get`/`config.set` 不带 `panel.` 前缀掉在两把名册尺之外、旧封套被宿主库**静默 `resolve(null)` 且 Go 侧零审计**

**立票时刻**：2026-10-02 10:0x +08，锚点 HEAD `f25f9572`（`dev`）
**来路**：只读腿 `114-a2`（`.scratch/wisp/probes/114/a2/census.md`，241 行 / 57,999 字节，`8d9f4102`→`f5f9cc34`→`981f08b2`）的 R9／R10 两条＋它新增的一条缺陷；台账 `A519`。
**性质**：⚠ **第三枚算产品行为**（页面以为话发出去了、Go 侧什么都没收到、也没有任何审计行），前两枚是仪器射程。
**⛔ 本票不接管票 114 的 9 枚未勾格**——那 9 枚是"差集"（档位／附件／工作区／消息），归票 114 自己的落地腿。

## 现量（腿的读数；行号一律〔待验〕，落地腿起手先复跑）

1. **可达性钉看不见生产那扇门**。钉只认 `.postMessage(`（`internal/panel/composer_test.go:381`〔待验〕），而生产入口是 `w.Bind("wispDispatch")` → `window.external.invoke` 那条 RPC（`cmd/wisp/panel_host_windows.go:551` 一侧）。⇒ **"这条链路可达"那一格，今天是用一把看不见门的尺在判**。
2. **名册两把尺的射程之外多出一批名字**。票 248 新增的 `config.get`／`config.set` **不带 `panel.` 前缀**（`internal/panel/bridge.go:66-67`，我现读为真），因此掉在 `routeLiteralRe`（`internal/panel/composer_test.go:394`〔待验〕）与 `panelMethodRe`（`cmd/wisp/git_test.go:394`〔待验〕）之外 ⇒ 票 114 的 **AC#8 那个洞被削宽**：新增方法名可以在**不触发任何名册尺**的情况下进树。
3. **旧封套静默丢弃**。页面若还按旧形状直接 `postMessage`，宿主库会**把这次调用的 Promise 解析成 `null` 并返回**、**Go 侧零审计**（腿报位置在依赖模块 `webview.go:162-168`〔待验，且属第三方模块，本仓不修上游〕）。⇒ 用户视角的症状＝**点了按钮，什么都没发生，日志里也什么都没有**。

## 要建什么

- [ ] **AC#1 可达性那枚钉改成问能力**：判据不许再认某个 JS 调用的**词面**，要问**"页面发起的一次调用，Go 侧有没有收到并回执"**（载体沿用票 33 那族"由页面自己报回"的定式：`REPLIED` 那形）。⚠ 正控必配：造一发"只 `postMessage`、不 `wispDispatch`"的假腿 ⇒ 指名用例**必须红**。⛔ 不许把两种拼法都列进白名单了事（那是把尺调宽，不是把门修好）。
- [ ] **AC#2 名册尺的射程要覆盖"不带前缀的那批"**：`config.get`/`config.set` 这类非 `panel.` 前缀的方法名必须被某把尺数到（要么两把尺各自扩射程，要么单立一枚"全量入向名册"尺）。判据＝种一枚"新增方法名没进名册"的假腿 ⇒ 指名那一步**必须红**；⛔ 现在它今天**是绿的**（这就是缺陷本体）。
- [ ] **AC#3 静默丢弃那一支要有声**：⛔ 不许改上游库；本仓射程内能做的是**在 Go 侧给"收到但形状不认识"落一条审计/告警**（或至少让面板宿主那条腿**具名拒收并回一个错误**）。判据＝一发假腿按旧封套发一次 ⇒ ①Go 侧有可读的现场（审计行或错误文本），②页面拿到的不是无声 `null`。
- [ ] **AC#4 三枚各自结线，不许合并成一句"入向修好了"**：AC#1 是钉的射程、AC#2 是名册的射程、AC#3 是产品行为——结案时逐格给读数。

## 禁区

- ⛔ 不许动 `frontend/**`／`design/**`（**既不读也不引**）；界面侧那半（改 `postMessage` 调用形）属 owner 自己用的那枚 agent，本票**只写 Go 侧**，需要它配合的地方写进"归口"由 owner 带话。
- ⛔ 不许改 `internal/panel/l2_grant_boundary_test.go` 那批冻结件射程之外的口径；⛔ 不许放宽任何名册尺来"让它绿"。
- ⛔ 阈值／SLO／golden／`internal/observe/thresholds.go`／`tools/d22scan/allowlist.txt` 一字节不动；⛔ 不许 `t.Skip`；⛔ 不许把 SKIP 读成通过。
- Git 纪律：只 commit 不 push；显式 pathspec；⛔ `git add -A`／`.`／`-a`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；临时件只建不删。

## 归口与排程

- 写面＝`internal/panel/**`＋`cmd/wisp/**` ⇒ **两枚包此刻都有验收腿在飞**（`198-v1` 正在整包跑 `cmd/wisp`）⇒ **本票按住**，起跑判据：`198-v1` 退出且 `git status --porcelain internal/panel cmd/wisp` 为空。
- 与票 114（差集九格）、票 242（`bindDigest` 绑定层与出向读面零仪器）**互不吞并**：那两票各自的名字与格子不许并进本票。
- 需要界面侧配合的只有一件：**把面板里的调用统一走 `wispDispatch`**（不带 `panel.` 前缀的那两名除外，它们就是 `config.get`/`config.set`）。这一句要 owner 带给界面侧 agent，⛔ 本编队永不转达。

## 8. 收只读预检腿 `253-p4`＋编排者裁定（2026-10-03 09:2x；台账 `A558`；票面原话一字未改，就地打旧）

交件凭据（盘上尺复量）：`.scratch/wisp/probes/253/p4/precheck.md` **141 行／22,709 字节**；三枚 commit `9051d51c`（锚＋骨架）→ `fae59004`（§1–§3）→ `f79c3a0b`（§4–§6）；只读、零 Go 命令、未 push、AC 框未碰。

- ★ **AC#2 我改裁＝走形ⓐ，⛔ 不解冻任何枚、⛔ 不摆给 owner**。这条**推翻我在 `A555` 里写的"AC#2 落地必撞 want-4、须具名解冻"**。p4 现量的根据：既有尺的分母里**不含"另一枚测试尺"**——`internal/panel/git_test.go:464` 的 `gitToolNamesUnder` 显式跳 `_test.go`、`whitelistMethodsFromSource` 只点名单文件 `bridge.go`、`scanRendererHostDoors` 只走 `frontend/src`、连冻结件 `internal/panel/l2_grant_boundary_test.go:369` 的产码枚举也跳 `_test.go` ⇒ **另立一枚"全量入向名册"尺（不改 `bridge.go`、不往 C17 四名里加真方法）属纯仪器面**。
- 形ⓑ（往 C17 加一枚真 `panel.` 方法）才红在 `internal/panel/git_test.go:385`（`equalStrings`，文案自述 "the four methods this ticket may not extend"）与 **`:517`**（`len(real) != 4`）；要变绿就得放宽名册尺＝撞本票禁区＋属契约面（`SPEC-12 §4.1` 人工批准）。**而且形ⓑ不对靶**：`config.get`/`config.set` 的要害正是"不带 `panel.` 前缀"（现读：`internal/panel/bridge.go:66-67` 的 `MethodConfigGet`/`MethodConfigSet`），加 `panel.` 方法不解决这个洞。
- ⛔ **AC#2 落地腿一条硬约束（我从 p4 读数推出来的）**：那发"新增方法名没进名册 ⇒ 指名那一步必须红"的正控**只能种不带 `panel.` 前缀的名字**（如 `config.foo`）。种 `panel.x` 会同时打红 `:385`/`:517`＝误伤别人的锚，会被误判成我方修法的错。
- **两处行号纠我自己**：`A555` 里我写的第二锚 **`:509` 错，真身 `:517`**（`sed` 现量）；`panelMethodRe` 在 `internal/panel/git_test.go:394`（不在 `cmd/wisp`，与 `253-r3` 同判、复认）；`routeLiteralRe` 真身 `internal/panel/composer_test.go:394`、闭集 `composerRouteLiterals()` 在 `:409-419`（5 名）、红点 `:521`。
- **AC#1 的增量改窄（推翻本票排程节里我给的宽义口径）**：`cmd/wisp` **今天确实已有"页面发起、Go 侧回执"的行为用例**——`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`（`cmd/wisp/panel_resident_windows_test.go:314`）与 `TestAC14AwaitedBindingReplyReachesThePage`（`:812`，两半互证：`countRequests==3` 且页面 `REPLIED`）。⇒ 真洞只在窄义：**没有任何词面/AST 尺断言产码里 `w.Bind("wispDispatch")` 存在**。所以 AC#1＝"把可达性裁决从词面尺迁到能力形＋补一发只 `postMessage` 不 `wispDispatch` 的正控"，⛔ 不是从零造一台仪器。可复用零件逐名（p4 现量）：`startPanelForTest:59-67`／`showAndWait:104-118`／`waitPanelTrue:85-95`／`evalOnPanelThread:166-189`／`recordingModeHandler`（`panel_host_windows_test.go:227-243`，`all()` 在 `:206-210`）／`awaitReport:215-229`／`reportJSEnv:250-253`／`countRequests:837-845`／`panelThreadWait:53`。三条复用红线＝自开 requestId 桶、别再 `SetHtml`（撞 AC#13 那枚次序钉 `:314`）、handler 只答 mode 路。
- **排程更新**：AC#1／AC#3 写面＝`cmd/wisp`（`255-r2` 在飞）⇒ 按住到它交件；AC#2 写面＝`internal/panel` 新增测试，但 `internal/panel` **编译带上了正被 `242-v1` 临时改写的 `internal/agent/approval`**（现量 1 枚 import）⇒ 同机跑测试互洗读数，按住到 `242-v1` 交件。**三格都无需 owner 介入。**

## 9. 收 AC#2 落地腿 `253-r5`＋编排者裁定（10-03 10:0x，台账 `A562`；票面原话一字未改，就地追加）

- **交件为真（我核过盘）**：唯一产码写面＝`internal/panel/inbound_roster_253_test.go`（新增 **834 行／35,189 字节**，commit `87bc6aca`）；`internal/panel/bridge.go` **一字未动**（我现跑 `git diff --numstat HEAD -- internal/panel/bridge.go`＝空，md5 `63e69a79422b577cca6a320c52915261` 与腿报逐字相同）；证据件 `.scratch/wisp/probes/253/r5/impl.md` 315 行／34,461 字节，commit 链 `32af81ac→87bc6aca→9099a006→ccd9543b→036eb133`，⛔ 未 push。**终态四数**（腿跑）＝`internal/panel` **111 PASS／4 FAIL／0 SKIP**，那 4 枚 FAIL 与起手基线**逐名相同**（＝在册常红的页面契约族，⛔ 不算凭据、未修、未放宽）。
- **这把尺的形状（比票面要求的更严，我认下）**：一枚**写死在测试里的六名封闭入向名册**（四名 `panel.*`＋`config.get`／`config.set`，取字面而⛔ 不引 Go 常量）对**五个独立分母双向对账**＝①`bridge.go` 包级字符串常量（标识符以 `Method` 开头**或**值呈点分形状 ⇒ 不认命名空间，`config.peek`／`settings.read`／`x.y` 全在射程）②真问一遍运行中的 `knownComposerMethod`（不复刻判据）③`bridge.go:146` guard 的 case 标签（AST）④`internal/panel/composer_dispatch.go:175` router 的 case 标签（AST，并要求 `default:` 必须调 `rosterMismatch`）⑤全包 ANSWERED 字面。**正控实测**（真树种一枚无前缀常量 `const ConfigPeek253 = "config.peek"`）⇒ 指名红在 `inbound_roster_253_test.go:437`，红句逐字含"a new inbound route entered the tree without being registered anywhere"。**恒真自查**＝把分母①换成反形（空集也算过＋只查单向）后**同一棵种了假腿的树全绿**，它自己的 `t.Logf` 还写着 `7 names declared` ⇒ **反形不可作凭据**（这条正是我那条"判据换成反形还绿＝不敏感"的规矩，本腿自己执行了）。
- ★**我的裁定（一件重要的下游约束）**：它**故意**让"只把名字补进测试名册、Go 不答它"这一式**继续红**（红在三支答侧：running guard／guard AST／router AST，判词逐字 "a rostered inbound name Go does not answer is a page-visible dead door"）。⇒ 我认这枚更严形＝**票 253 AC#2 的结案口径**：名字进名册就必须同时进 guard 与 router，**"挂个门牌没人接"以后不许算登记**。⚠ **这条与 `Q-76` 直接咬合**：若 owner 批了第五枚入向方法名（面板"停止"），**落地那一次五个分母必须同批动**，⛔ 不许只补名册、⛔ 不许删这把尺换绿。副本读数（完整登记⇒`missing=(none) extra=(none)`）由 `TestRegistrationSilencesTheRosterRuler` 在 TempDir 里量出＝**真树没动过 guard**（那是契约面，本腿不许做，做得对）。
- **它新量到一枚 `253-p4` 未记的射程事实**：冻结尺的 `answered` 只从 `knownComposerMethod` 导出（`internal/panel/l2_grant_boundary_test.go:777` 显式排除其它函数）⇒ **router 的 case 列表此前不在任何尺的分母里**，分母④补的就是这一格。
- **want-4 两枚锚一行未动**（我复认其形状）：`internal/panel/git_test.go:385`（∈`TestGitDimensionHasNoModelCallableTool`）与 `:517`（∈`TestPlantedGitToolShapesGoRed`）起手／种形／终态三次读数均 PASS；⛔ 也没让 `bridge.go` 真多出第五枚 `panel.*`。另纠一枚易误读：`grep -c '"panel\.' bridge.go` 给 5，那 1 枚是 `:35` 的文档注释、尺的正则匹配不到 ⇒ 尺数到的仍是 4。
- ⛔ **AC#2 这格我不翻**（实现者自述不算凭据）：待非实现者验收腿 `253-v1`，排在 `220-r1` 退出之后（那枚腿正写 `internal/panel`，同时跑测试＝互洗）。⚠ 一处半量欠账归我裁：`cmd/wisp/config_reload.go:86` 那枚 `config.reload` **算不算入向名**——`cmd/wisp` 有 `255-r2` 在飞，等它交件我现跑再裁。
- ★**AC#1／AC#3 的排程更新**：`internal/panel` 此刻空了，但这两格写面＝**`cmd/wisp`**（`255-r2` 在飞）⇒ 照旧按住；`220-r1` 已被我派进 `internal/panel`＋`internal/agent/approval`（票 220 的 join 那一格），⛔ 别把它当成"我这轮的 253 腿"。
