# 257-r2 落地件 — 票 257 形 ⓒ 的 `cmd/wisp` 那一半（首启回执文案）

> 腿：`257-r2`；票：`.scratch/wisp/issues/257-clean-machine-provider-registry-nil-blocks-writes.md`（AC 框零触碰）。
> 前件：`257-a1/census.md`（92 行／`ed3fd270`）＋`257-a2/census.md`（257 行／`cbf4f1b2`）＋`257-r1/evidence.md`（182 行，`internal/config` 那一半，九节全满）＋票面 §8（选形 ⓒ，账 `A543`）＋`A560`。
> 本节次安排：§0–§4 随首枚实现 commit 落盘，§5 随突变 commit，§6 随尾程整包 commit，§7–§10 交件时收满。**未写满的节次此刻不在文件里**（不放过半句占位）。

## 0. 起手锚

| 项 | 读数（本腿现量） |
|---|---|
| 进场时刻 | `2026-10-05 11:04:54 +0800` |
| 进场自取 HEAD | `c6cf66e64849955bf92a4b3356c096ea81c35563`（branch `dev`）＝"ledger(A614＋A615) 票 267 五格全勾" 那枚 |
| 起手 `git status --porcelain -- cmd internal` | **0 行**＝`cmd/**`＋`internal/**` 此刻没有别人的活在树上（写面互斥闸通过）；开工后到我首次 commit 前仍是只有我这枚的形状（见 §8） |
| 起手 `git status --porcelain` 总行数 | 632＋（同机多枚在飞的正常量级；`design/**` 删除族与 `.scratch/**` 过程件占大头，与 `257-r1` 进场同形） |
| 并发腿（派单告知） | 另两枚**只读**腿在跑，不写文件；两者都在读 `internal/panel` ⇒ 本腿 ⛔ 不碰 `internal/panel/**` |
| CI 背景 | GitHub Actions 正在这台机器的 self-hosted runner 上跑 `slo-full`（D32）与 `test-windows` ⇒ 整包 `go test` 按派单留到尾程，本件尾程读数若被争用洗到，具名标〔机器争用待复跑〕 |
| 本腿写面（授权） | `cmd/wisp/firstrun.go`（回执文案）＋同包新增 `_test.go`＋`.scratch/wisp/probes/257/r2/**`；⛔ 不碰 `internal/config`、⛔ 不碰 `internal/panel` |
| 起手 `GOFLAGS= go build ./...` | rc=0（`11:08:26`，改动前基线） |

### 0.9 ★ 撞腿登记（11:23 现量，本腿停手上报前的完整事实；⛔ 不归本腿裁决）

派单给本腿的前提之一："代号 `r2` 是干净的（`257/` 里只有 `a1`／`a2`／`r1`）"＋"另有两枚只读腿在跑，它们不写文件"。**这条在 11:19–11:23 之间被盘上事实推翻**：

| 时刻 | 现量 | 出处 |
|---|---|---|
| 11:04:54 | 本腿进场，`ls .scratch/wisp/probes/257/`＝`a1 a2 r1`，**无 `r2` 目录** | 本腿命令实录 |
| 11:08:49–11:17:26 | 本腿做完：预检名册 → 改 `firstrun.go`（＋61 行）→ 新增 `firstrun_257_test.go` → 三发读数（`preflight-ticket198.log`／`after-text-ticket198.log`／`new-tests-first.log`／`new-tests-second.log`），md5 链：`firstrun.go`＝`0491339282492f2cabdbf5be576c8a57`、`firstrun_257_test.go`＝`23a7a56795019fabf506f6726541b049`（11:17:50 与 11:23:03 两取相同） | 本件 §1/§6 |
| 11:19:29 | commit `9c0d4c9a` 落盘：**另一枚自名 `257-r2` 的腿**提交了 `.scratch/wisp/probes/257/r2/evidence.md` 骨架（249 行），**并把我上面那两枚日志一并 commit** | `git show --name-only 9c0d4c9a` |
| 11:12:54 | 那枚件的 §0 自报进场时刻（晚于本腿 8 分钟），其 §0.1 把 ` M cmd/wisp/firstrun.go` 与两枚日志判为"我自己前一小时那半段留下的活"——**这是误认**：那 61 行＋两枚日志是本腿 11:11–11:12 产出的（本件 §2 逐字 before/after 为证） | `git show 9c0d4c9a:...evidence.md` |
| 11:14:58–11:22:42 | 同一目录继续长出**非本腿产出**：`dumpcfg.go`／`dumpcfg.log`（11:15）／`cmdwisp.test.exe` **40 MB**（11:16:54）／`msg-skeleton.txt`（11:19:02）／`snapshot-firstrun.go`＋`snapshot-firstrun_257_test.go`（11:20:12，内容＝本腿那两枚文件的快照）／`build-mutants.py`（11:21:40，读的正是 `cmd/wisp/firstrun.go`）／`mut/M1..M6`（11:22:08–11:22:42）／`baseline-cmd-wisp.log`（11:22:42＝有人在本腿写面上跑了整包） | `ls -l --time-style=+%H:%M:%S` |
| 11:21:5x | **本腿的 Write 覆盖了那枚腿在 `9c0d4c9a` 里那份 evidence.md 的工作树副本**（同一文件路径，两腿共用一枚 `257-r2` 目录）。⛔ 本腿没有删任何东西：那份骨架完整可读，取法＝`git show 9c0d4c9a:.scratch/wisp/probes/257/r2/evidence.md > <外部路径>`；本件 §0.9 就是这次覆盖的具名登记 | 本节 |

**本腿据此采取的动作（不扩权、不改别人的件）**：① 立刻把本腿的实现＋测试＋本件以**显式 pathspec** 落 commit（先让 HEAD 有一份带 md5 链的可归因产物）；② ⛔ 不 commit 别人的 `mut/`／`cmdwisp.test.exe`／`snapshot-*.go`／`dumpcfg.*`／`build-mutants.py`／`baseline-cmd-wisp.log`，也 ⛔ 不删（README 规则 8）；③ 尾程每一次突变前后**重取 md5**，若发现工作树两枚文件被外部改动，立即停手并把差异贴进本节；④ "哪一枚 `257-r2` 的件算交件"＝编排者裁，本腿不自裁。

**11:25:3x 追加：① 已经失败，本腿的 commit 被第三枚腿卷走了**。本腿 11:24 用显式 pathspec `git add` 好这十枚文件后跑
`git commit -F … --only -- <十枚显式路径>`，git 回 "no changes added to commit"（未落任何 commit）；
11:25:10 落地的 `3f0c4fff` **标题是别票的**（"probes(268-a2): 独立复核初版…"），
`git show --name-only 3f0c4fff` 却列着本腿那十枚：`cmd/wisp/firstrun.go`、`cmd/wisp/firstrun_257_test.go`、
`.scratch/wisp/probes/257/r2/{evidence,gate-d22scan,gate-pathlen,new-tests-first,new-tests-second,preflight-255-223}.log/md`
——即 `268-a2` 那一发**没带 pathspec 的 commit 把本腿暂存区整卷带走**（同机另一枚腿 `b672f853` 11:24:12 是带显式路径的，形状正常）。
本腿内容完好可核：`git cat-file blob HEAD:cmd/wisp/firstrun.go | md5sum` ＝ `0491339282492f2cabdbf5be576c8a57`
＝ 工作树 ＝ §0.9 表里那一枚；`cmd/wisp/firstrun_257_test.go` 同形（`23a7a56795019fabf506f6726541b049`）；
HEAD 版 evidence.md 含本节 §0.9（grep 命中 1）。
⇒ 后果两条，都记在本腿账上但不归本腿裁：**(a) 归因错**：本腿实现的名册里必须把 `3f0c4fff` 算作
"内容是我的、message 不是我的"那一枚；**(b) 危险**：只要还有一枚腿不带 pathspec 提交，
**本腿尾程突变的中间态就有被卷进 HEAD 的风险** ⇒ 本腿把突变窗口压到"一次 Edit＋一次定向 `-run`＋一次还原"，
每发还原后立刻重取 `md5sum` 与 `git cat-file blob HEAD:<path> | md5sum` 双读，任一漂即停手。

**同时刻另一枚 `257-r2` 的件在长**：`cmd/wisp/firstrun_257_nonpreset_test.go`（8,241 字节，11:24:12，未 commit、⛔ 不是本腿写的）。
本腿不碰、不删；它进不了本腿任何 commit（枚枚显式 pathspec）。


## 1. 撞钉预检名册（今天绿着的相邻用例，逐枚）

★ 派单点名的那枚钉**找到且逐枚读了断言体**：`cmd/wisp/firstrun_198_test.go:215` 在 `TestTicket198AC1CreatedFileHoldsEveryStaticSectionAndNoInventedTables` 里逐字禁 `[llm.providers`／`[models.local_override`／`[plugins.` 出现在**首建文件的内容**里（读的是 `os.ReadFile(cfgPath)`，⛔ 不是读 stderr）。⇒ ⓒ 的新文案只要落在 stderr 这条"给人看的回执"通道就不碰它；本腿新增的 AC#1 用例反过来**自己钉了一遍**"生成的文件里没有 `[llm.providers`"（`firstrun_257_test.go`，防的就是"文案漏进产物"这一形）。
同类钉逐枚读到的另三枚：`firstrun_198_test.go:237`（回执 ⛔ 不得借 cause=missing 那三串标记，逐字 `"cause=missing"`／`"config.toml 读不到：文件不存在"`／`"本次运行继续用内存里的旧配置"`）、`firstrun_198r2_test.go:363`（回执必名 7 串＋`wisp <词>` 必须真有 `cmd<词>`：`\bwisp ([a-z][a-z0-9-]*)` 拿 AST 数过）、`firstrun_198r2_test.go:414`（二跑不重复回执，含 `"新建默认配置"` 与 `"wisp secret set"`）。
⚠ 由 `:363` 那把尺得出的**硬约束，本腿照办**：新文案里出现过的 `wisp <拉丁词>` 只有 `wisp run`（`cmdRun` 存在）；`panel-inbound`／`resident` 一类**写不得**（`cmdPanel-inbound` 不存在，会当场红），所以三形状那一段用中文形状名指代，不造命令词。

**起手绿着的名册（`-v` 实跑，非静态推断）**：

| 发 | 时刻 | 命令 | 读数 |
|---|---|---|---|
| 1 | 11:08:49 | `go test -count=1 -run 'TestTicket198' -v ./cmd/wisp/` | rc=0；**11 枚 `--- PASS`、0 FAIL、0 SKIP**，逐名：`TestTicket198AC1FirstRunCreatesConfigThenLeavesItAlone`／`…AC1DirectoryShapedConfigIsNeverOverwritten`／`…AC1CreatedFileHoldsEveryStaticSectionAndNoInventedTables`／`TestTicket198FirstRunReceiptStaysOutOfTheFourCauseSentences`／`TestTicket198FirstRunCallerIsTheRunEntryOnly`／`TestTicket198R2AC2CreatedFileCarriesNothingButTheSchemaDefaultTags`／`TestTicket198R2AC2ExportedDefaultTableStillRendersAsItsTags`／`TestTicket198R2AC5CreationFailureIsLoudAndLeavesNoHalfFile`／`TestTicket198R2AC4ReceiptNamesTheRealEntryPoints`／`TestTicket198R2J1TheAssemblyRootStillCreatesNothing`／`TestTicket198AC3CreatedFileLandsPrivate`。原件 `.scratch/wisp/probes/257/r2/preflight-ticket198.log` |
| 2 | 11:09:05 | `go test -count=1 -run 'TestTicket255\|TestTicket223' ./cmd/wisp/` | rc=0（不带 `-v`＝读不到名册，这发只算"整族绿"的粗筛；名册在发 3 补） |
| 3 | 11:19:08–11:19:5x | 同发 2 加 `-v`，重定向 `preflight-255-223.log` | rc=0；**24 枚 `--- PASS`、0 FAIL、0 SKIP**：255 族 13 枚（`SplitOnlyClaimsSectionsWithALiveReader`／`HotRowRosterCoversTheRegistry`／`RosterEvidenceLinesStillSayWhatTheyClaim`／`RosterStillMatchesTheActualReadSites`／`ReceiptOmitsPanelFromTheImmediateSentence`／`ReceiptOmitsASectionWithNoReaderAnywhere`／`ReceiptStillNamesTheLiveReadSection`／`ReceiptSentenceAssemblyIsFiltered`／`WindowOptionsFollowTheConfigSource`／`AssemblyRootGeometrySourceReachesTheWindowOptions`／`PanelHostBuildsItsWindowOptions`／`HostStillDoesNotParseConfigItself`／`PanelRosterVerdictIsTheHonestShape`／`RestartTierKeysAreBackedByATest`）＋223 族 10 枚（`RunArmsTheReloadTick`／`HandEditedFsLooseningCostsAnL2Card`／`RefusedLooseningKeepsOldValues`／`TighteningRaisesNoCard`／`ModeLooseningChangesTheRunningModeAfterAllow`／`RestartTierSaysItWillNotApply`／`FailureSentencesAreDistinct`／`PanelInboundSaysHotReloadIsDisabled`／`PermissionDeniedSitsInItsOwnSentence`／`R2FailureSentenceRouting`） |

**相邻但本腿不动的钉（读过断言体，判"会不会被我打红"）**：

| 钉 | 位置 | 我为什么不会打红它 |
|---|---|---|
| 首建文件逐字节＝tag 渲染 | `firstrun_198r2_test.go:221` | 我的改动一行都没进 `SaveFile`／`NewDefaults`；文案只进 stderr。发 1 与发 3 复跑仍绿（11:12:28 那发见 §6） |
| 回执 ⛔ 借 cause=missing 措辞 | `firstrun_198_test.go:237` | 新文案三处自查：`AC#2` 用例直接把这两串当反控告出（`firstrun_257_test.go` 的 `borrowed` 循环） |
| 二跑不重复回执 | `firstrun_198r2_test.go:414` | 三段新印品全在 `ensureFirstRunConfig` 的创建分支内（`os.Stat` 判缺才走到），第二跑根本不进这一支；复跑绿 |
| 名册七枚枚数＝7（`reflect.DeepEqual`） | `internal/panel/config_route_248_test.go` 族 | 形 ⓑ 已毙 ⇒ 我一枚字段都没加，`internal/panel` 未碰 |
| 七枚的写门读盘不读内存／拒因三 tag | `internal/config/settings_257_test.go`（257-r1 交） | 产码未动；我只**字面**引用那三句 tag（不 import 常量，理由见 §2 末），r1 那五枚仍绿＝尾程整包名册里对账 |

**"测试里先手工塞一份完整 config.toml"这一族（＝本票立案的那个洞）在本腿的处理**：⛔ 不照抄。`cmd/wisp` 里预塞文件的相邻件（`run_test.go:90`、`providers_test.go:46`、`panel_config_248_test.go:76`、`secret_test.go:774`、`run_mode101_test.go:123`、`approval_reply_201_test.go:128`、`resident_task_source_live_246_windows_test.go:336`、`leg_sink_nail_131_windows_test.go:420`）全部**手写** `[llm.providers.acme]` 起步——它们在量别的格子，不是本票的洞。我的三枚 AC#1 用例一律从 `t.TempDir()`＋**没有** `config.toml` 起步、由真入口 `runTextTask` 建文件（`§3` 逐字写了这一发怎么造），手加段是**机主动作**、发生在建完之后。

## 2. 改动逐处（before／after 逐字＋为什么落这一支）

**唯一产码改动＝`cmd/wisp/firstrun.go` 的 `ensureFirstRunConfig` 尾部追加三段印品＋一段说明注释；`git diff --numstat`＝`61  0  cmd/wisp/firstrun.go`（加 61 删 0）。新增测试件＝`cmd/wisp/firstrun_257_test.go`（472 行，本腿唯一新增文件）。**

| # | 落点（新行号） | before（逐字） | after（逐字，只列印品体） | 为什么落这一支 |
|---|---|---|---|---|
| 1 | `firstrun.go:123-150`（注释块） | 模型句之后直接 `return true, nil` | 28 行注释：三形代价的账（`A543`）、拼法更正的出处（`parse.go` 的 `DisallowUnknownFields`，现量 `:72`＋`:123`）、ⓐ 的 `roles.chat` 就地填形状（257-r1b 的 M3 红句）、`A560` 的"引用解引用没接线"、以及"只走 stderr、不落进文件"的两枚 198 钉 | 落这一支的理由＝票 §8-3 ⓒ 的**边界①②**都写在回执上；这段注释是本腿唯一允许的"新立论"位置，产码形状零变化 |
| 2 | `firstrun.go:151-160` | —— | `fmt.Fprint(stderr, "wisp run: 上面那句模型只是第一样。设置页那七枚字段（服务商的 base_url、api_key_ref、凭据，模型的 context_window、price.in、price.out，还有聊天模型）今天都不建行，只改已有的行；要在这一页配上模型，得在这份文件里手加三样，缺一不可：" + "wisp run: 第一样＝一节 [llm.providers.<名>]，就是服务商那一行（名字对上内置预设的，protocol 与 base_url 可以留空；非预设名必须自己写 protocol，否则这份文件加载不过）。第二样＝它的模型行 [llm.providers.<名>.models.<模型 id>]，里面写 enabled = true。第三样＝就地填已有的 [llm.roles.chat] 那一节，把 provider 与 model 两枚一起点上名（别再追加一节同名的，那在 TOML 里是 duplicate table，文件直接加载不过）。界面不会替你建这一行，它只会告诉你去哪一节建；三样齐了这七枚才全部写得进，只加第一样只解锁服务商那三枚。` | **AC#1 的终态格**： ⓒ 兑现＝把"静态链 7/7"教到能照着走通（257-a1 §ⓒ-2 量的三样；`3/7` 与 `7/7` 两个数都进文案，并被两枚用例各钉一边，见 §3 发 2／发 3）。第二样用 `guidanceModelRow` 的同一串拼法、第三样用"就地填"而非"追加"，是与 `internal/config` 那半句逐字对齐的结果 |
| 3 | `firstrun.go:162-170` | —— | `fmt.Fprintf(stderr, "wisp run: 写不进去的时候有三种原因，各是一句不同的话，不会合成一句「配置未生效」：\n" + "  第 1 种拒因：文件没建——这一种刚才那一发已经替你办完，%s 现在是真的文件；首启之前没有任何旧配置可言。\n" + "  第 2 种拒因：行不存在——那一页改不了服务商与模型的存在性，去这份文件里手加上面那三样，加完才写得进；这一条说的不是你的值不对。\n" + "  第 3 种拒因：校验不过——行在，但那个值过不了这份 schema 的校验（引用形没写对前缀、点名的模型不在目录里，都算这一种）；要改的是值，文件一个字节都没动。\n", cfgPath)` | **AC#2**：三句各配各的补救，且 `Fprintf` 的第二参数把"第 1 种已办完"落到**这台机的真路径**上（不是模板话）。三句的 tag 与 `internal/config/settings.go:78-80` 那三枚常量**逐字同串**（`refusalFileMissing`／`refusalRowMissing`／`refusalInvalid`），拒写现场与首启回执说的是同三种事 |
| 4 | `firstrun.go:172-184` | —— | `fmt.Fprint(stderr, "wisp run: 改完什么时候才算用上，按进程形状分三种说法，不是一句「重启就好」：\n" + "  在控制台里跑 wisp run——这个进程带着每 1s 重读一次 config.toml 的看门狗，[llm] 属可热加载档，手改的值一秒内就换进这台进程的内存；但模型通路是启动时建一次的，热加载不会替它换脑，真正发请求还是按启动时那一份。\n" + "  没有可答卡入口的常驻形状——任务腿过不去控制台那道闸时，这个进程里可能根本没有会重读盘的东西，手改与面板写在两个方向上都只能等下一次启动。\n" + "  设置页那一页——它那条腿自己明说不带轮询，写入回执固定说要重启进程；页面上的读数在重启之前也不会跟着你手改的文件走。\n" + "wisp run: 凭据这一格只有引用会进这份文件：改 api_key_ref 换的是名字不是密钥本身，把那份引用再解一次是新建端点时才做的事，所以换过 key 的引用同样要重启才算用上；这一页任何时候都不回显密钥的值，只说已录入还是没录入。`) | **票 §8 边界①＋`A560` 的反噬**：三形状各一句（有控制台／无控制台／面板写入），⛔ 不再用"手改＝热加载认／面板写＝要重启"那两格冒充全部；凭据那段只写真接了的那条路（`run.go:435-436` → `internal/llm/resolver.go:141` 每建一次端点解一次），⛔ 不把 `config.NewManager` 的 `resolveRefs` 列成入口（`A560`：生产三处第二参数全 nil） |

**为什么测试里的期望全是字面量**：`257-r1b` 的 M1 量过同包常量的盲区（改常量值＝断言跟着改，永不红）。`internal/config` 那三枚 tag 是**未导出**的，`cmd/wisp` 想引用也引用不到；本腿把三句当字面抄进 `firstrun_257_test.go:53-57`，于是"两半文案漂移"与"折叠成一句"这类真缺陷才咬得动（§5 的 M2/M3 就是它咬的）。这条同时是本件的一处**已知脆性**，具名进 §7-N1。

## 3. AC#1 干净机真跑读数

**这一发怎么造出来的（⛔ 没有"测试先塞一份 config.toml"这一步）**：
`t257CleanMachine(t)`（`cmd/wisp/firstrun_257_test.go:74-95`）＝
1. `dir := t.TempDir()`；
2. `os.Stat(filepath.Join(dir, "config.toml"))` 必须是 `fs.ErrNotExist`，否则 `t.Fatalf("AC#1 premise broke: the clean machine already holds ...")` —— **"没有 config.toml"这一步是断言出来的，不是假设出来的**；
3. 调 `run198(t, dir)`＝`runTextTask(runSpec{argv, stdout, stderr, dataDir: dir, notify: 空})`，也就是**真入口**（`main.go cmdRun -> runTextTask -> ensureFirstRunConfig`，`run.go:245`）；
4. 断言退码仍＝2（票 198 的"建完仍按未配置失败"不许被我放宽）＋文件现在真在＋stderr 里有 `"新建默认配置"`，否则 `t.Fatalf` 说"没有回执可读"。

**发 1＝`TestTicket257R2AC1CleanMachineReceiptTeachesTheWalkableChain`（11:17:16 PASS）**
- 回执含 `[llm.providers.<名>`／`第一样＝一节 [llm.providers.<名>]`／`[llm.providers.<名>.models.<模型 id>]，里面写 enabled = true`／`第三样＝就地填已有的 [llm.roles.chat] 那一节`／`三样齐了这七枚才全部写得进`／`只加第一样只解锁服务商那三枚`（逐串断言，缺一串即红）；
- 反控：回执 ⛔ 不含 `[providers.`（`A543`：那串会撞 `DisallowUnknownFields`，面板整条链起不来）；
- 通道分离：新建的 `config.toml` 里 ⛔ 不含 `[llm.providers` ⇒ 指引只在"给人看的回执"里；
- `config.NewManager(cfgPath, nil)`（生产形状，`A560`）加载过，且 `m.Config().LLM.Providers == nil`（票面前提在**我这枚树上、我这个入口建出来的文件**上重认一次）；
- 逐枚试写七枚（名册七枚＝`provider_base_url`／`provider_api_key_ref`／`model_context_window`／`model_price_in`／`model_price_out`／`role_chat_model`／`provider_credential`，第 7 枚按 257-r1 的口径走**引用腿** `SetProviderAPIKeyRef`，⛔ 不碰 `StoreCredential`，AC#3）：
  **读数（逐字，`-v` 里的 `t.Logf`）**：`AC#1 clean machine: 7 refused writes by reason: map[第 2 种拒因：行不存在:7]` —— 七枚全拒、七枚**各带且只带一枚** tag（`seen != 1` 即红），且七句都点名 `[llm.providers.`（ⓒ 的"告诉你去哪一节建"），折叠串 `"配置未生效："`／`"写入失败，配置未生效"` 逐枚断言不在场。

**发 2＝`TestTicket257R2AC1ReceiptChainWalkUnlocksAllSevenFields`（11:17:16 PASS）**：把回执教的三样**照教的样子**动文件——`strings.Replace` 就地填已有的 `[llm.roles.chat]`（前置断言：文件里必须找得到 `"[llm.roles.chat]\nprovider = ''\nmodel = ''"`，找不到即红＝"教错了形状"）＋追加 `[llm.providers.deepseek]`／两条 `models.<id>` 且各带 `enabled = true` ⇒ `NewManager` 加载过 ⇒ **七枚 7/7 全部接受且各报键路径**（`len(written)==0` 判红，所以"接受了却没说落了哪"也拦得住）；文件终态 ⛔ 无 `\napi_key =` 明文键（AC#3）。
⚠ 本腿在这发上**先红后绿过一次**，两次跑都留了盘（`new-tests-first.log`／`new-tests-second.log`）：首跑读数 `AC#1 RED: the taught chain unlocked 5 of 7 fields, want 7`，红因是**我的测试自己**把 `api_key_ref` 与 role 的 `model` 先填成了后面要写的同一个值 ⇒ `writeOneKey` 老实回"文件里已经是这个值"、键路径为空。这不是产码缺陷，是"接受了却没键路径"与"没变化"两种形状在我仪器里没分开；改法是让手加段与试写段用**不同**的占位值（`env:T257_R2_HANDADD_ENV_NAME` vs `env:T257_R2_ENV_NAME_NEVER_SET`；`model = deepseek-chat` vs 写 `deepseek-reasoner`），11:17:16 复跑 7/7。

**发 3＝`TestTicket257R2AC1ReceiptOnlyFirstItemUnlocksThreeOfSeven`（11:17:16 PASS）**：只加第一样 ⇒ 恰好 3 枚接受，且逐名核对这 3 枚就是服务商那三枚（多一枚也算红）——这是把回执里"只加第一样只解锁服务商那三枚"这句数量话**变成可失败断言**的那一发（`257-a1` §1-7／`257-a2` §2.3 末两条独立复认过的 3/7，本腿在真跑里量到同一数）。

**终态判语**：形 ⓒ 在 `cmd/wisp` 这一面**真兑现**＝干净机起步 → 首份配置由真入口建 → 七枚逐枚拒、拒句只讲"行不存在"并点名 `[llm.providers.<名>]` → 照回执教的三样手加 → 7/7 写得进；两个数量断言（3/7 与 7/7）都在盘上跑过，不靠静态推。

## 4. AC#2 三句不同的话：逐字＋各自触发条件

**先在产码里找到那三支今天各走哪条路（本腿自己 grep＋逐枚读，不抄派单行号）**：

| 因 | 今天的真身（现量） | 现场原话挂在谁身上 |
|---|---|---|
| 文件没建 | `internal/config/loader.go:77`（`readConfigFile` 的缺文件分支，句尾接 `refusalFileMissing`＋"运行一次 wisp run 写出首份配置"）；`cmd/wisp` 侧的旧形状是 `run.go:417` 的 `wisp run: 配置未就绪（Unconfigured）：%v`（票 198 已把"干净机建文件"接在 `run.go:245`，所以首启后这一因在 `wisp run` 上不再发生） | 257-r1 那半；本腿的**回执**替它说话（第 1 句把真路径 `%s` 印出来，宣布"这一种已办完"） |
| 行不存在 | `internal/config/settings.go:367`／`:379`（`requireCatalogEntry`／`unknownProviderErr`，接 `refusalRowMissing`＋`guidanceModelRow`／`guidanceProviderRow`）＋`:235-243` 的 `requireChatProvider`（role 那枚在"config 没点名 provider"时改判第 2 种） | 拒写现场；本腿的**回执第 2 句**预先说清"这一条说的不是你的值不对" |
| 校验不过 | `internal/config/settings.go:293`（`writeOneKey` 预写 `validate()` 门，接 `refusalInvalid`＋"要改的是值；文件一个字节都没动"） | 拒写现场；本腿的**回执第 3 句**给同一形状 |

**三句逐字（首启回执里的版本，11:12:28 那发的真 stderr，见 `after-text-ticket198.log`）**：
1. `  第 1 种拒因：文件没建——这一种刚才那一发已经替你办完，<绝对路径>\config.toml 现在是真的文件；首启之前没有任何旧配置可言。`
2. `  第 2 种拒因：行不存在——那一页改不了服务商与模型的存在性，去这份文件里手加上面那三样，加完才写得进；这一条说的不是你的值不对。`
3. `  第 3 种拒因：校验不过——行在，但那个值过不了这份 schema 的校验（引用形没写对前缀、点名的模型不在目录里，都算这一种）；要改的是值，文件一个字节都没动。`

**触发条件各不相同**：① 只在该数据根从没建过配置时出现（并宣布自己刚被解决）；② 文件在、要写的行不在（provider 行／model 行／roles.chat 的 provider 三形）；③ 行在、值过不了同一份 `validate()`（引用形缺前缀、点名的模型不在目录）。

**仪器（`TestTicket257R2AC2ThreeRefusalsStayThreeSentences`，11:17:16 PASS）**的牙齿，不是"含不含"而是：
- 三枚 tag 各**恰好出现在一枚行上**（`len(hits) != 1` 即红＝既拦折叠、也拦一句里塞两因）；
- 每枚 tag 那一行还必须带**它自己的补救**：第 1 种＝`替你办完`＋`没有任何旧配置可言`；第 2 种＝`改不了服务商与模型的存在性`＋`手加上面那三样`；第 3 种＝`过不了这份 schema 的校验`＋`文件一个字节都没动`（三组补救串互不重叠 ⇒ "两行话其实是一行话"会红，见 §5-M3）；
- 折叠判据正向化：`配置未生效`／`重启就好` 在整个回执里**只能出现一次**，且只能出现在禁句那行（"不会合成一句…"／"不是一句…"），出现在别的任何行都判红；
- 借句判据：`cause=missing`／`本次运行继续用内存里的旧配置` 不得出现在回执（票 223/198 的归因隔离，本腿自己复钉一遍）。

**第二发（`TestTicket257R2AC2EffectTimingSaysThreeProcessShapes`，11:17:16 PASS）**＝票 §8 边界①在 AC#2 上的延伸：**"什么时候算用上"也不许折成一句**——三形状各一枚（`每 1s 重读一次 config.toml 的看门狗`／`根本没有会重读盘的东西`／`自己明说不带轮询`），另钉两枚限定句（`热加载不会替它换脑`＝内存与装配期一次建好的模型通路不是一回事；`把那份引用再解一次是新建端点时才做的事`＝`A560` 那条没接线的通道不许被说成"改了就生效"）。
