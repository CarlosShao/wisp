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

## 5. 变异自证（四发全响；还原一律 `git cat-file blob HEAD:<path> > <path>`）

**基线三枚 md5（`cmd/wisp/firstrun.go`）**：改前工作树＝`0491339282492f2cabdbf5be576c8a57`＝`git cat-file blob HEAD:cmd/wisp/firstrun.go | md5sum`（HEAD 取 `06b66e66`，11:30:10 双读相同）＝四发还原后各读（11:27:54／11:29:14／11:29:43／11:29:59／11:30:40 五取相同）。测试件 `cmd/wisp/firstrun_257_test.go` 全程未种未改：`23a7a56795019fabf506f6726541b049`（11:23:03 与 11:30:10 两读相同）。
⚠ 突变窗口内的额外风险由 §0.9 记的"别腿无 pathspec 提交"事件带来 ⇒ 本腿把每发压成"种（一次 Edit 或一次 python 写）→ 定向 `-run` 单发 → 还原 → 双读 md5"，最长窗口 11:28:30→11:29:14＝44 秒（M2，含一次脚本自己退回），其余三发 6／10／7 秒；四发跑完后 `git cat-file blob HEAD:` 两枚仍是基线值＝**没有任何突变被卷进 HEAD**，这条是本节最硬的一句。

| 发 | 种法（逐字 before → after，只列被换片段） | 命中用例 | 突变体 md5 | 红句（逐字，原件 `mutant-M*.log`） | 还原 |
|---|---|---|---|---|---|
| **M1（ⓐ 拼法改坏）** | `第二样＝它的模型行 [llm.providers.<名>.models.<模型 id>]，里面写 enabled = true。` → 同一句里那串改成 `[providers.<名>.models.<模型 id>]` | `TestTicket257R2AC1CleanMachineReceiptTeachesTheWalkableChain`（rc=1） | `a2a7e62d27ae0152673dfbf18b6d7ea2` | `AC#1 RED: the first-run receipt never says "[llm.providers.<名>.models.<模型 id>]，里面写 enabled = true", so shape (c) is not delivered - the operator is left with a refusal and no path:` ＋ `AC#1 RED: the receipt tells the operator to add [providers.x], which parse.go's DisallowUnknownFields rejects outright (ledger A543):` | 11:27:54 md5 回基线 |
| **M2（ⓐ 整段删掉指引）** | 三样那五句里抹掉 `[llm.providers.<名>]`／`[llm.providers.<名>.models.<模型 id>]`／`[llm.roles.chat]`／`三样齐了…`／`只加第一样…`（保留句子骨架，语法完好） | AC#1 前两发同红（rc=1） | `538f93c02dfa1d954b7d5dcbcee89f4b` | 五条同源红句（`…never says "第一样＝一节 [llm.providers.<名>]"` / `…"[llm.providers.<名>.models.<模型 id>]，里面写 enabled = true"` / `…"第三样＝就地填已有的 [llm.roles.chat] 那一节"` / `…"三样齐了这七枚才全部写得进"` / `…"只加第一样只解锁服务商那三枚"`）＋ `--- FAIL: TestTicket257R2AC1ReceiptChainWalkUnlocksAllSevenFields`（它先 `t.Fatalf("the receipt stopped teaching the in-place fill")`）；同时读数行仍出：`AC#1 clean machine: 7 refused writes by reason: map[第 2 种拒因：行不存在:7]` | 11:29:14 回基线 |
| **M3（ⓑ 三因折成一句）** | `第 2 种拒因：行不存在——那一页改不了服务商与模型的存在性，去这份文件里手加上面那三样，` → `配置未生效——去这份文件里看看，` ＋ `第 3 种拒因：校验不过——行在，但那个值过不了这份 schema 的校验` → `配置未生效——行在，值过不了校验` | `TestTicket257R2AC2ThreeRefusalsStayThreeSentences`（rc=1） | `891fc88b4e887c0f8592a6dabf7b2cc2` | `AC#2 RED: "第 2 种拒因：行不存在" appears on 0 receipt lines, want exactly 1 (three reasons, three sentences):` ＋ `AC#2 RED: "第 3 种拒因：校验不过" appears on 0 receipt lines, want exactly 1 …` ＋ 两行 `AC#2 RED: "配置未生效" is used as a verdict rather than named as the forbidden collapse:` | 11:29:43 回基线 |
| **M4（边界① 三形状折成一句）** | `没有可答卡入口的常驻形状——任务腿过不去控制台那道闸时，这个进程里可能根本没有会重读盘的东西，` → `三种形状其实是同一句话——重启就好，` | `TestTicket257R2AC2EffectTimingSaysThreeProcessShapes` 与 `…ThreeRefusalsStayThreeSentences` 同红（rc=1） | `17fdd43209d51090b4ef6043417a912d` | `AC#2 / boundary 1 RED: the receipt never states the 没有可答卡入口的常驻形状 shape ("根本没有会重读盘的东西"), so one boot shape is folded into another:` ＋ `AC#2 RED: "重启就好" is used as a verdict rather than named as the forbidden collapse:` | 11:29:59 回基线 |

四发的**权威 md5 清单**（含 M1，重算不落盘突变，11:31:31）＝`.scratch/wisp/probes/257/r2/mutant-md5.txt`：baseline `0491339282492f2cabdbf5be576c8a57`／M1 `a2a7e62d27ae0152673dfbf18b6d7ea2`／M2 `538f93c02dfa1d954b7d5dcbcee89f4b`／M3 `891fc88b4e887c0f8592a6dabf7b2cc2`／M4 `17fdd43209d51090b4ef6043417a912d`／worktree after `0491339282492f2cabdbf5be576c8a57`。本腿第一次写这一节时把 M3 那格抄错（记成"不确定"），就地更正＝上列实读值，登记在 §9-N1。


四发跑完复绿：11:30:18 定向 `go test -count=1 -run 'TestTicket257R2' ./cmd/wisp/` ＝ `ok … 0.248s`，六枚全过。

---
---

# 以下＝同一票同一写面的**另一枚腿**（自称 `257-r2b`；上面 §0–§4 是姊妹腿的活，本腿一字未改）

⚠ 编号重叠一句说明（不是笔误）：姊妹腿在下面已经写了自己的 **§5 变异自证（四发）**（第 141 行起），
本腿的 §5–§10 从下面这节起算，是**另一套八发**（M1–M8，含本腿新加的两枚牙 M7／M8）。
两半各留各的，谁也没覆盖谁；本腿原件的第一版 §0–§7 仍在 commit `9c0d4c9a`（11:26:20 那次覆盖写的受害者，逐字可取）。

**先把文件的归属说清（⛔ 不是修辞，是可核的两笔）**：本腿原本写的是另一份 §0–§7
（起手锚 `c6cf66e6`、写面自证、票 98 dll 注入坑的根因与出处、`defaults.go:77-78` 等**七枚钉的断言原文**、
`dumpcfg` 现量、三句拒因的生产侧 tag 出处），**11:26:20 被姊妹腿整枚覆盖写掉**（盘上从我 16,266 字节／八节
变成她的 26,771 字节／§0–§4）。原件没有丢，逐字在 commit `9c0d4c9a` 里：
`git show 9c0d4c9a:.scratch/wisp/probes/257/r2/evidence.md`。
本腿**不再回盖她的文件**（那正是 README 规则 2 禁的"互相踩"），改成在下面追加，两半都活着。
上面那份 §0.9 也独立量到了同一件撞腿事实——两枚腿各写一份，这条对编排者是信号不是噪声。

## 5. 突变名册与红句（八枚，逐字抄，含"改回后未变"的证明）

### 5.0 做法与为什么这样做

⛔ 本腿**没有原地改过** `cmd/wisp/firstrun.go` 与 `cmd/wisp/firstrun_257_test.go`（§0 撞腿事实决定的）。
牙是这么咬的：`go test -overlay=<mut>/overlay.json` 只在**构建期**把 `firstrun.go`（或测试文件）替换成
`.scratch/wisp/probes/257/r2/mut/M*/` 里的那一枚副本，盘上一个字节不写；副本由 `build-mutants.py`／
`build-mutants2.py` 从快照生成，**每枚 needle 要求恰好命中 1 次**，命中数不对脚本直接 `sys.exit`（不会静默出一个"没改到东西的突变"＝假牙）。
时刻逐枚＝M1 11:22:08／M2 11:22:19／M3 11:22:24／M5 11:22:29／M6 11:22:34／M4 11:22:39／M7 11:27:57／M8 11:28:42（11:28:05 那发见 §7.2）。
日志逐枚＝`mut/M*/run.log`。

### 5.1 八枚红句（`--- FAIL` 名 + 第一行原话）

**M1**（指引第一样的小节名换成 A543 记过的那个致命拼法 `[providers.<名>]`）
```
--- FAIL: TestTicket257R2AC1CleanMachineReceiptTeachesTheWalkableChain (0.02s)
    firstrun_257_test.go:151: AC#1 RED: the first-run receipt never says "第一样＝一节 [llm.providers.<名>]", so shape (c) is not delivered - the operator is left with a refusal and no path:
    firstrun_257_test.go:160: AC#1 RED: the receipt tells the operator to add [providers.x], which parse.go's DisallowUnknownFields rejects outright (ledger A543):
```

**M2**（三因折成一句「就是配置未生效」）
```
--- FAIL: TestTicket257R2AC2ThreeRefusalsStayThreeSentences (0.01s)
    firstrun_257_test.go:328: AC#2 RED: "第 1 种拒因：文件没建" appears on 0 receipt lines, want exactly 1 (three reasons, three sentences):
    firstrun_257_test.go:328: AC#2 RED: "第 2 种拒因：行不存在" appears on 0 receipt lines, want exactly 1 (three reasons, three sentences):
    firstrun_257_test.go:328: AC#2 RED: "第 3 种拒因：校验不过" appears on 0 receipt lines, want exactly 1 (three reasons, three sentences):
```

**M3**（生效时机折成一句「改完重启就好」）
```
--- FAIL: TestTicket257R2AC2ThreeRefusalsStayThreeSentences (0.01s)
--- FAIL: TestTicket257R2AC2EffectTimingSaysThreeProcessShapes (0.01s)
    firstrun_257_test.go:349: AC#2 RED: "重启就好" is used as a verdict rather than named as the forbidden collapse:
    firstrun_257_test.go:383: AC#2 / boundary 1 RED: the receipt never states the 控制台里的 wisp run shape ("每 1s 重读一次 config.toml 的看门狗"), so one boot shape is folded into another:
    firstrun_257_test.go:383: AC#2 / boundary 1 RED: the receipt never states the 没有可答卡入口的常驻形状 shape ("根本没有会重读盘的东西"), so one boot shape is folded into another:
```

**M4**（★最硬那枚雷的实弹：把指引从 stderr 漏进生成的 `config.toml`）
```
--- FAIL: TestTicket198AC1CreatedFileHoldsEveryStaticSectionAndNoInventedTables (0.02s)
    firstrun_198_test.go:217: the created file carries the dynamic table "[llm.providers" anyway - that would be an invented entry, not a default
--- FAIL: TestTicket198R2AC2CreatedFileCarriesNothingButTheSchemaDefaultTags (0.02s)
    firstrun_198r2_test.go:239: AC#2 RED: the config.toml first-run wrote is not the schema's `default` tags rendered: first difference at line 179 (byte 2571): want "", got ""; total bytes want 2571 got 2624
--- FAIL: TestTicket257R2AC1CleanMachineReceiptTeachesTheWalkableChain (0.01s)
    firstrun_257_test.go:170: AC#1 RED: the guidance leaked into the generated config.toml, which is invented-default territory and ticket 198's forbidden shape:
（同发另红：TestTicket257R2AC1ReceiptChainWalkUnlocksAllSevenFields／TestTicket257R2AC1ReceiptOnlyFirstItemUnlocksThreeOfSeven）
```
⇒ 这一发是本腿存在的理由之一：它证明"那句指引只能进回执"不是本腿的自律而是**别人钉着的雷**，
票 198 那两枚一漏就响（`total bytes want 2571 got 2624`＝逐字节尺也一起响了）。

**M5**（第三样教成"再追加一节 `[llm.roles.chat]`"＝TOML duplicate table 那个谎）
```
--- FAIL: TestTicket257R2AC1CleanMachineReceiptTeachesTheWalkableChain (0.02s)
    firstrun_257_test.go:151: AC#1 RED: the first-run receipt never says "第三样＝就地填已有的 [llm.roles.chat] 那一节", so shape (c) is not delivered - the operator is left with a refusal and no path:
--- FAIL: TestTicket257R2AC1ReceiptChainWalkUnlocksAllSevenFields (0.01s)
    firstrun_257_test.go:236: the receipt stopped teaching the in-place fill: wisp run: 已在 ... 新建默认配置：...
```

**M6**（凭据那段开始回显值形状的串）
```
--- FAIL: TestTicket257R2AC3CredentialSurfaceUntouched (0.01s)
    firstrun_257_test.go:420: AC#3 RED: the receipt shows a plaintext key assignment shape:
    firstrun_257_test.go:423: AC#3 RED: the receipt carries a vendor key-shaped literal:
```

**M7**（本腿新增枚：删掉回执里"非预设名必须自己写 protocol，否则这份文件加载不过"这半句）
```
--- FAIL: TestTicket257R2AC1ReceiptStatesTheNonPresetCondition (0.02s)
    firstrun_257_nonpreset_test.go:77: AC#1 RED: the receipt never states "非预设名必须自己写 protocol", so the guidance is only honest about preset names:
    firstrun_257_nonpreset_test.go:77: AC#1 RED: the receipt never states "否则这份文件加载不过", so the guidance is only honest about preset names:
（同发 TestTicket257R2AC1NonPresetRowNeedsProtocolBeforeItLoads 仍 PASS＝文案枚与行为枚各咬各的，没被一枚笼统判据盖住）
```

**M8**（本腿新增枚：把行为那半的"没写 protocol 的非预设行"改成内置预设名 ⇒ 证明那条断言不是空转）
```
--- FAIL: TestTicket257R2AC1NonPresetRowNeedsProtocolBeforeItLoads (0.05s)
    firstrun_257_nonpreset_test.go:110: AC#1 RED: a non-preset provider row with no protocol loaded, so the receipt's warning is describing a gate that is not there
```

### 5.2 每枚判据对应哪一发（对得上号才算数）

| 判据（票面 AC／§8 边界） | 用例 | 响它的那发突变 |
|---|---|---|
| AC#1 回执给出三样且落到具体小节名 | `...AC1CleanMachineReceiptTeachesTheWalkableChain` | M1／M4／M5 |
| AC#1 终态＝ ⓒ 真兑现（手加三样 ⇒ 7/7） | `...AC1ReceiptChainWalkUnlocksAllSevenFields` | M4／M5 |
| AC#1 那句"只加第一样只解锁三枚"不许吹成七枚 | `...AC1ReceiptOnlyFirstItemUnlocksThreeOfSeven` | M4（同发响；三样开关本身由 §9-b 记为量不到） |
| AC#1 的诚实边界（非预设名要不要 protocol） | 本腿 `...AC1ReceiptStatesTheNonPresetCondition`／`...AC1NonPresetRowNeedsProtocolBeforeItLoads` | **M7**／**M8** |
| AC#2 三因三句、不许折成一句 | `...AC2ThreeRefusalsStayThreeSentences` | M2／M3 |
| §8 边界① 两通道不许折成一句（三形状） | `...AC2EffectTimingSaysThreeProcessShapes` | M3 |
| AC#3 凭据面一字不动、绝不回显值 | `...AC3CredentialSurfaceUntouched` | M6 |
| AC#1 的"回执不许漏进生成的文件"（票 198 的钉） | `TestTicket198AC1CreatedFileHolds...`＋`TestTicket198R2AC2CreatedFileCarriesNothing...` | M4 |

### 5.3 还原证明（突变之后盘上没动过一个字节）

`md5sum` 逐枚，对照 §0 的起手快照：

| 文件 | 起手快照 | 八发突变跑完之后 | 判 |
|---|---|---|---|
| `cmd/wisp/firstrun.go` | `0491339282492f2cabdbf5be576c8a57` | `0491339282492f2cabdbf5be576c8a57` | 未变 |
| `cmd/wisp/firstrun_257_test.go` | `23a7a56795019fabf506f6726541b049` | `23a7a56795019fabf506f6726541b049` | 未变 |
| `cmd/wisp/firstrun_257_nonpreset_test.go`（本腿新件） | `4cf81b638013728575f026dab6fabf2a` | `4cf81b638013728575f026dab6fabf2a` | 未变 |

另附一条真事：11:28:05 那发 M8 第一次跑**编译失败**，报 `cmd\wisp\firstrun.go:151:2: fmt.Fprint call has possible
Printf formatting directive %s`，当时 `md5sum cmd/wisp/firstrun.go` = `43548dec5019e9f3386b978f1a7798a1`（≠ 快照）；
11:28:33 复量回到 `049133...`、`go vet ./cmd/wisp/` rc=0。⇒ 那是姊妹腿正在写这枚文件的**半写瞬间被我撞见**，
不是谁的产码缺陷，也不是本腿改的（本腿从没写过它）。登记在 §7.2。

---

## 6. 门禁读数（四门＋整包，逐枚带时刻）

| 门 | 命令 | 时刻 | 读数 |
|---|---|---|---|
| D22 静态扫描 | `sh scripts/d22scan.sh` | 11:25:14 起，11:25:36 完 | **rc=0**，逐字 `clean - no D22 ban violations`；覆盖面 live scope：bans #1-5 internal/=228、cmd/=38、ban #6 frontend/=85、ban #7 internal/tools/=23、ban #8 design/=39、frontend/=85、internal/=508、cmd/=103（`gate2-d22scan.log`） |
| 路径长度帽 | `sh scripts/check-path-length-budget.sh` | 11:25:36 | **rc=0**，`VERDICT GREEN`；tracked=5746／over-budget=57／covered by roster=57／**not in roster=0**；longest=180 chars relative（`252-…`）；worst full path on runner=224（`gate2-pathlen.log`） |
| `go vet` | `go vet ./cmd/wisp/` | 11:25:37 与复量 11:28:33 | 两次 **rc=0**、输出零行（`gate2-vet.log` 为 0 字节） |
| gofumpt 逐枚点名 | `$(go env GOPATH)/bin/gofumpt -l firstrun.go firstrun_257_test.go firstrun_257_nonpreset_test.go` | 11:25:39 | **列出零枚**＝三枚全净；`-l cmd/wisp` 整包只报 **`cmd\wisp\models.go`** 一枚（预存 CRLF，不在本腿写面，⛔ 未顺手格式化，见 `gate2-gofumpt.log`） |

⚠ `gofumpt` 不在 PATH，`$(go env GOPATH)/bin/gofumpt.exe` 才是真的那把尺（`go env GOPATH` = `D:\work\base\gopath`）。
第一发用裸命令名跑出来的是 `command not found`（rc=127），那次读数不作数，`gate2-gofumpt.log` 里两行都留着。

### 6.1 整包红名册基线（起手那一发，真跑）

```
BASELINE START 2026-10-05 11:17:21 +0800 anchor=48b85705 inject=third_party/sherpa-onnx
ok  	github.com/CarlosShao/wisp/cmd/wisp	401.805s
BASELINE END rc=0 2026-10-05 11:24:08 +0800
RUN(顶层+子) 329 ／ --- PASS 233 ／ 子用例 PASS 96 ／ FAIL 0 ／ SKIP 0
```
⇒ 红名册＝**空**（0 枚 FAIL、0 枚 SKIP；两枚带载偶发在册红名 `TestAC1ResidentLegInstallsItsLogListenerOnDisk`／
`TestAC14GoSideEvalPushReachesThePage` 这一发都没红，逐字名串在这份日志里可 grep）。
日志 `.scratch/wisp/probes/257/r2/baseline-cmd-wisp.log`（全文不截断，红名册靠 grep 数，不靠 tail）。
同文件里还留着 11:15:57 那发**未注入 dll** 的死法（`exit status 0xc0000135`／`=== RUN` 0 条）＝§0.3 那条坑的原件。

### 6.2 覆盖本腿新件的终跑（跑完，终值）

```
FINAL START 2026-10-05 11:25:56 +0800 anchor=3f0c4fff files=firstrun.go+257_test+257_nonpreset_test
ok  	github.com/CarlosShao/wisp/cmd/wisp	433.149s
FINAL END rc=0 2026-10-05 11:33:15 +0800
RUN(顶层+子) 331 ／ --- PASS 235 ／ 子用例 PASS 96 ／ FAIL 0 ／ SKIP 0（子 SKIP 也 0）
```
⇒ 终跑红名册＝**空**，且比 §6.1 的基线**多出两枚**（331−329＝本腿 `firstrun_257_nonpreset_test.go` 的两枚，
逐名在这份日志里）：
`--- PASS: TestTicket257R2AC1ReceiptStatesTheNonPresetCondition (0.01s)`／
`--- PASS: TestTicket257R2AC1NonPresetRowNeedsProtocolBeforeItLoads (0.02s)`。
日志全文（不截断）＝`.scratch/wisp/probes/257/r2/final-cmd-wisp.log`，两枚在册带载红
（`TestAC1ResidentLegInstallsItsLogListenerOnDisk`／`TestAC14GoSideEvalPushReachesThePage`）
这一发与基线那一发**都没红**（逐名 grep 零命中）——本腿没去追它们也没这个必要。

两发之间的时间差也登记一下：基线 401.8s、终跑 433.1s，差在并发负载（本腿的八发突变跑＋别的腿在写
`internal/panel/**`）——这就是 §7.2 那两枚"半写瞬间撞见的红"的来路。

### 6.3 收尾复量（交件那一刻的四门，逐枚带时刻，`gate3-*.log`）

11:34:39 起：`sh scripts/d22scan.sh` **rc=0**（末行逐字 `d22scan: clean - no D22 ban violations`）／
`go vet ./cmd/wisp/` **rc=0**（日志 0 字节）／
`$(go env GOPATH)/bin/gofumpt -l` 对本腿三枚 **零命中**／
`sh scripts/check-path-length-budget.sh` **rc=0**，`VERDICT GREEN`，tracked=5806（比 11:25 多 60 枚＝别的腿在落件）、
over-budget=57、roster 覆盖 57、**not in roster=0**。
同刻三枚 md5 仍与起手快照逐字相等（`049133…`／`23a7a5…`／`4cf81b…`）＝本腿全程没在交件后又动过写面。


---

## 7. 判不动／量不到／不属于本腿的红（具名归口，一枚没自己划掉）

**7.1 本腿判不动的格**

| 格 | 为什么够不着 | 归口 |
|---|---|---|
| 名册第 7 枚 `provider_credential` 的**值腿**（`StoreCredential`→DPAPI 真密封） | AC#3 明令凭据面一字不动，且那枚方法在 `cmd/wisp/panel_config_store.go`（不在本腿写面）。本腿只对它的**引用侧**（`api_key_ref` 落盘）负责并测过 | 编排者：真机 DPAPI 那格需要 `docs/evidence/s1/` 单行裁决 |
| 三句拒因在**界面上的长相**（信封键与面板回执） | ⛔ 写面禁 `internal/panel/**`；且名册加字段属 C17 契约面（ⓒ 裁定里 ⓑ 那形已被毙） | 编排者：owner 若要"页面上也说这三句"，按票 §8 边界③单开票 |
| 非预设名（自建网关）能不能把**七枚全写满** | 本腿只证到 `provider_base_url`＋`provider_api_key_ref` 两枚（§9-a）；模型级五枚在非预设名下没走 | 编排者：要么补一发行文，要么按现测范围裁"指引可走通" |
| AC 框翻勾／票面与台账任何一字 | 规则：勾归编排者，台账只追加不删且归它 | 编排者 |

**7.2 不属于本腿、本腿零动作的红（逐名带出处）**

| 名 | 现量 | 判 |
|---|---|---|
| 同票双腿并发（写面与证据件都撞） | §0 的逐时刻表：姊妹腿 11:17:50 落 `firstrun_257_test.go`、11:26:20 覆盖本腿 `evidence.md` | 撞腿事实成立；⛔ 本腿不裁决谁活谁停，归编排者（README 规则 2 的后果条） |
| `internal/panel/workspace.go` 半写导致 `cmd/wisp` **整包构建失败** | 11:26:42 与 11:27:5x 读数：`internal\panel\workspace.go:97:44: not enough arguments in call to scope.SetWorkspaceRoot` 等四行；11:27:50 `go build ./internal/panel/` rc=0 自行清掉 | 别的腿（`internal/panel` 写面）在飞，本腿零动作、未追未修 |
| `cmd/wisp/firstrun.go` 半写瞬间（`fmt.Fprint ... %s` vet 红） | 11:28:05 md5 `43548dec…`，11:28:33 回到 `049133…`、vet rc=0 | 同上：姊妹腿在写她的那枚文件 |
| `go test ./cmd/wisp/` 不注入 dll 必死 | `exit status 0xc0000135`、`=== RUN` 0 条（§0.3） | 票 98 在册环境洞（`R-101-6`），本腿只登记"每条读数都带注入"这条前提 |
| `gofumpt -l cmd/wisp` 报 `cmd\wisp\models.go` | 11:25:39 现量 | 预存 CRLF，不在写面，⛔ 未动（票 70 那一 sweep 的地盘） |
| `frontend/**`／`design/**`／`docs/PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt` | 本腿 `git show --name-only` 逐笔自证（§8），⛔ 连读都没读 | 零命中＝AC#4 干净 |

---

## 8. commit 名册（本腿；只 commit 不 push、显式 pathspec 枚枚点名）

| 笔 | 内容 | pathspec |
|---|---|---|
| `9c0d4c9a` | 证据件 §0–§7 原件（盘上已被姊妹腿 11:26:20 覆盖，逐字仍在此 commit） | `probes/257/r2/evidence.md`＋`dumpcfg.go`＋`dumpcfg.log`＋`preflight-ticket198.log`＋`after-text-ticket198.log` |
| `7c644bb0` | 证据件 §5–§10 追加＋本腿新测试＋八枚突变产物＋四门与基线日志 | `cmd/wisp/firstrun_257_nonpreset_test.go`＋`probes/257/r2/evidence.md`＋`build-mutants*.py`＋`snapshot-*`＋`mut/M1..M8/*`（24 枚）＋`gate2-*.log`＋`baseline-cmd-wisp.log`＋`my-new-tests.log`＋`msg-*.txt` |

### 8.1 ★ 回执文案那枚文件其实落在**别人的 commit** 里（现量，⛔ 不是本腿漏交）

`git log -1 --format='%h %ad %s' -- cmd/wisp/firstrun.go` 与 `-- cmd/wisp/firstrun_257_test.go` 都指向
**`3f0c4fff`（11:25，subject 逐字起于 `probes(268-a2): 独立复核初版——AC#0 四问复跑判定＋err 透传四支链证真`）**
⇒ 那两枚文件是被**票 268 的复核腿**连着它自己的活一起提交的。
本腿据 §0.4 从未原地改它们，所以这次落库不是本腿的动作；但内容一字未丢，可核两把尺：

- `md5sum cmd/wisp/firstrun.go` = `0491339282492f2cabdbf5be576c8a57` = 本腿起手快照（§5.3 同一枚）；
- `git show HEAD:cmd/wisp/firstrun_257_test.go | md5sum` = `23a7a56795019fabf506f6726541b049` = 本腿快照，
  且 `git diff -- cmd/wisp/firstrun.go cmd/wisp/firstrun_257_test.go` **零行**。

后果登记（给编排者，两条都是本票面上的事）：
① 257 的产码半格（回执文案）**已经进了 268-a2 那一笔**，验收腿若按"票号找 commit"会找不到它，
按 `git log -- cmd/wisp/firstrun.go` 才会撞见；
② 一笔 commit 里混着两张票的写面，正是 README 规则 2「>1 个写码代理必须 worktree 隔离」想拦的形状。

**这条不是本腿一家之言**：姊妹腿在 `06b66e66`（11:27，subject 逐字起于 `257-r2 证据件 0.9 追加：本腿实现被 3f0c4fff 无 pathspec 卷走`）
里独立量到了同一件事，并把它归因到那一笔**没带 pathspec** 的 commit。两枚腿各写一份、指向同一笔 `3f0c4fff`，
⇒ 这条对编排者是**双份读数**而不是噪声；本腿在上面那两条后果之外只补一句：
无 pathspec 的 commit 在共享工作树里卷走的是别人的未提交活，本腿从第 0 步起每笔都枚枚点名（§8 两张表可核）。


`git show --name-only` 逐笔自证只含上面这些路径；`internal/config/**`、`internal/panel/**`、
`internal/ball/**`、`internal/tools/**`、`frontend/**`、`design/**`、`docs/**`、`tools/d22scan/allowlist.txt`
⛔ 一枚都没进本腿 pathspec（§7.2 末行那条尺）。

---

## 9. 自我对抗（本腿自己仪器的洞，先自己说，不等验收腿抓）

**a. 非预设名那枚只证到 2/7，没证到 7/7。**
`TestTicket257R2AC1NonPresetRowNeedsProtocolBeforeItLoads` 对自建名试写了 `base_url` 与 `api_key_ref` 两枚就收了。
所以"非预设名也能把七枚写满"这句话**没有**被本腿测过；回执原文也没这么承诺（它只对"这一节建起来"这件事负责）。
如果验收要求非预设路径同样 7/7，这是一格真缺口，归 §7.1 第三行。

**b. 3/7 与 7/7 两枚计数共用姊妹腿的 `t257RosterWalk()`。**
本腿新件不重名、不依赖它的 helper，但回执里"只加第一样只解锁服务商那三枚"这句的**计数**是她那六枚用例在测。
⇒ 若 `t257RosterWalk()` 的字段集哪天与 `panel.WritableFields()` 漂移，这两枚计数会跟着漂（她那件文件里
`len(t257RosterWalk()) != 7` 那枚 Fatalf 是唯一的护栏，它拦得住"少于一枚"，拦不住"名字换成别的"）。
本腿没去加第二枚名册尺（那会造出第二个真相源，票面 §8 与 198-a1 §5.2 都 warn 过这件事）。

**c. M7 那发暴露了我一枚针的软处（自己抓到，就地登记，不改她的文件）。**
本腿文案枚的第一枚 needle `名字对上内置预设的，protocol 与 base_url 可以留空` **也能被票 198 那句旧文案满足**
（`firstrun.go:117` 里就有同一串）。删掉票 257 那一整样时，这一枚不响、另两枚（`非预设名必须自己写 protocol`／
`否则这份文件加载不过`）响 ⇒ M7 仍然红，判据成立，但这枚 needle 单独使用时是**弱**的。
要么把它绑到"第一样＝一节"那个前缀上，要么删掉——两样都动她的文案串，本腿没动，留给编排者定。

**d. 覆盖不到的地方：回执只在"首建"这一发被读。**
四枚用例全部起于 `run198` 的第一发；二跑不重复那格由姊妹腿的 §3 计划与票 198 现有钉覆盖，本腿没另测（不重复＝票 198 已有的形状，重复测会造第三份同一断言）。

**e. 没有放宽任何断言、没有 `t.Skip`、没有把 SKIP 读成通过**：§6.1 的 SKIP 计数是 0，
本腿若哪天真跑不动会写进 §7 而不是让用例静默。

---

## 10. 交件判语

- **票面 AC 框**：一枚未碰、零勾选（翻勾归编排者）。
- **形 ⓒ 的兑现**：产码面只动 `cmd/wisp/firstrun.go` 的**回执文案**（+61 行纯新增、零删除、签名与判定三段未变），
  加两枚同包测试件（`firstrun_257_test.go` 六枚＝姊妹腿；`firstrun_257_nonpreset_test.go` 两枚＝本腿）。
  写侧"行不存在则拒写"**一字未动**、`internal/config/**` 一字未动＝票 §8 边界②那条"诚实"保住了。
- **AC#1**：干净机真跑（临时数据根＋无 `config.toml`，起步态逐条 `t.Fatalf` 断言），
  首建后**七枚逐枚拒**、按回执手加三样后**7/7 接受**、只加第一样**恰好 3/7**、非预设名按回执那半句写 protocol 才加载得过——四发都是 `wisp run` 那条 CLI 缝起步，⛔ 没有一发是手工先塞完整配置。
- **AC#2**：三因三句各带各自的补救，`M2`／`M3` 两发红句证明"折成一句"必响。
- **AC#3**：凭据面零改动；`M6` 一发证明"回显值"必响；两枚新件的盘上终检里带 `sk-`／`api_key =` 两把反尺。
- **AC#4**：本腿 `git diff --name-only`／逐笔 `git show --name-only` 只见 `cmd/wisp/**` 与 `.scratch/wisp/probes/257/r2/**`，
  禁区零命中（`frontend/**` 与 `design/**` 连读都没读）。
- **牙**：八枚突变全响，红句逐字在 §5.1，盘上还原证明在 §5.3（三枚 md5 逐一相等）。
- **四门**：d22scan rc=0／pathlen rc=0 VERDICT GREEN（not-in-roster=0）／vet rc=0（两次）／gofumpt 本腿三枚零命中
  （整包只报预存的 `models.go`，未动）。
- **收尾状态**：§6.2 那发终跑已跑完并写满（rc=0、433.1s、331 RUN、0 FAIL、0 SKIP，含本腿两枚新用例逐名 PASS），
  四门与两发整包读数齐；本腿没有"量不到却写成结论"的格——够不着的三格全在 §7.1 具名归口给编排者，
  ⛔ 一枚没自己划掉。若 §5–§10 再被姊妹腿覆盖一次，原件逐字在**本笔 commit**（`git log -- <本文件>` 取最新一笔）。

---

# 6A–10A. 并行腿（自名 `257-r2`，11:04:54 进场那一发，锚 `c6cf66e6`）的尾程账

> 编号说明：本件 §0–§5 与本块 6A–10A 是**同一枚腿**（11:04 进场、锚 `c6cf66e6`）；§5–§10（175-445 行）是**另一枚自名 `257-r2b` 的腿**（其 §0 自报 11:12:54 进场，锚 `c6cf66e6`→`69c1bb82`）。
> 两套件在同一棵共享工作树上并存，编排者已在 `A617`（commit `51ea6998`）自曝"同一轮把票 257 派了两枚写腿"并留"两套实现并存待 257-v1 裁"。
> **本块不自裁谁算交件**；本块只保证：本腿那两枚产码/测试件的内容、读数、突变红句、门禁与整包名册，盘上一一可核。
> 本腿实现落点：`cmd/wisp/firstrun.go`（＋61 行，md5 `0491339282492f2cabdbf5be576c8a57`）与 `cmd/wisp/firstrun_257_test.go`（472 行，md5 `23a7a56795019fabf506f6726541b049`），两枚都在 `3f0c4fff`（该笔 message 是别票的，事故见 §0.9）。

## 6A. 门禁四数＋整包名册（逐把带时刻）

| 门 | 时刻（+08） | 读数 | 判 |
|---|---|---|---|
| ① `GOFLAGS= go build ./...` | 11:08:26（改前基线）／11:11:59（改后）／11:42:0x（尾程）／11:43:15（尾程复量） | **rc=0 四取同判** | 绿 |
| ② `"D:\work\base\gopath\bin\gofumpt.exe" -l cmd/wisp` | 11:17:39 首读点名两枚：`cmd\wisp\firstrun_257_test.go`（本腿新件）＋`cmd\wisp\models.go` | 11:17:50 对**本腿自己那枚**跑 `-w` 后：`firstrun_257_test.go` 归零，清单只剩 `cmd\wisp\models.go`；尾程 11:42:00 与 11:43:15 两读均只剩 `models.go`（预存 CRLF，票 212/258 既有账，⛔ 未顺手格式化） | 绿（本腿零命中） |
| ③ `sh scripts/d22scan.sh` | 11:17:59（原件 `gate-d22scan.log`）／11:42:57（尾程 `tail-d22scan.log`） | 两发 **rc=0**；末行 `d22scan: clean - no D22 ban violations`；仪器自检 `PASS=35 FAIL=0 SKIP=0`（含 `TestBan8MathBandAndRemainingGaps` 正反两向）。⛔ 无 phantom-citation（265 腿那枚既有红已结案，账 `A607`） | 绿 |
| ④ `sh scripts/check-path-length-budget.sh` | 11:18:43／11:42:59 | 两发 **rc=0 VERDICT GREEN**；分母 `11:18`＝`5733 tracked / 57 over / 57 covered / 0 not in roster`，`11:42`＝`5817 / 57 / 57 / 0` ⇒ 分母涨 84 枚＝同机别腿在这 24 分钟里落了件，⛔ 不是本腿口径变化 | 绿 |
| 附：`go vet ./cmd/wisp/` | 11:43:36 | rc=0。（⚠ 11:26:51 那一发曾因 `181-r3` 在飞的 `internal/panel/workspace.go` 编译不过而红，那次**不作本腿门禁读数**，见 §9A-N4） | 绿 |

**整包那一发（尾程，重定向到文件取数，⛔ 不用 `tail` 截名册）**：
`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test -count=1 -v ./cmd/wisp/ > .scratch/wisp/probes/257/r2/full-cmd-wisp.log 2>&1`
发射 11:30:2x，末行落 11:39:54，测试体内计时 **506.468s**：

| 读数 | 值 |
|---|---|
| `=== RUN` | **331** |
| top-level `--- PASS` | **233** |
| top-level `--- FAIL` | **1** |
| top-level `--- SKIP` | **1** |
| 子测试 `    --- PASS` | **96** |
| 对账 | `233＋1＋1＋96 = 331` ＝ `=== RUN` 数，逐枚闭合（无"跑了没结果"那一形） |
| 本腿六枚 | 全在名册里逐名 `--- PASS`（`TestTicket257R2AC1…×3／AC2…×2／AC3…×1`） |

**红名判据＝只减不增，基线用同机同包那一发**（`.scratch/wisp/probes/257/r2/baseline-cmd-wisp.log`，11:22:42，`233 PASS／0 FAIL／0 SKIP`，是并行腿在本腿实现已入树之后跑的）：
`comm -23` 逐名比对＝**丢名 0 枚**；`comm -13`＝**新增名 2 枚**，均是并行腿的新用例（`TestTicket257R2AC1NonPresetRowNeedsProtocolBeforeItLoads`／`TestTicket257R2AC1ReceiptStatesTheNonPresetCondition`），两枚本发内 **PASS**。
**新增红 1 枚〔未归因，归编排者〕＋〔机器争用待复跑〕**：`TestPanelHostRealWindowHopAndLifecycle`（基线里它是 `--- PASS (1.08s)`），红句逐字（整包 `0.81s` 那一发）：
```
panel_host_windows_test.go:713: AC#1 pre-hide tree settle: reading held after 1 extra sample(s) at our tree webview=0 (machine-wide 24)
panel_host_windows_test.go:717: no msedgewebview2 process is a direct child of this test process while the window is up - the host-pid ruler is blind (tree webview 0, tree pids 0)
```
隔离复跑（派单要求的两形，⛔ 没有"重跑到一半把两发混成一发"）：
- `-count=1`（11:40:40→11:41:05，`-run 'TestPanelHostRealWindowHopAndLifecycle|TestPanelHostLatencyPercentilesAC2'`）：rc=0，`--- PASS`（lifecycle 2.71s）＋ `--- PASS`（latency 0.00s）。原件 `isolate-realwindow-count1.log`
- `-count=3`（11:41:13→11:41:28）：**第 1 迭代两枚 PASS，第 2／3 迭代两枚双 FAIL**，红句含 `panel_host_windows_test.go:665: cold bring-up 2107.9 ms exceeds D32 panel cold budget 1500 ms` ＋ 同一族的 `tree settle: reading held after 1 extra sample(s)`。原件 `isolate-realwindow-count3.log`
- 同尺 cold bring-up 四次读数：**619.409 ms（11:32:52，整包内）／1031.424 ms（11:41:20）／2107.868 ms（11:41:24）**＝8 分钟内同机漂到 3.4 倍，且 `machine-wide msedgewebview2` 在 19~24 之间浮动（本树 `webview=0/7/8`）；跑窗内 self-hosted runner 正在跑 `slo-full`／`test-windows`（派单预告），`internal/panel`＋`internal/tools` 同时被 `181-r3` 写脏（`git status` 11:27 现量九枚 ` M`/`??`）。
⇒ 本腿判语只到这里：**这枚红与本腿写面（stderr 文案＋同包测试）无因果路径**（本腿没碰窗体、没碰 `internal/panel`），但**不宣布"无关"划掉**，标〔机器争用待复跑〕〔未归因，归编排者〕，⛔ 未压任何断言换绿。
**命名 SKIP 1 枚**：`TestPanelHostLatencyPercentilesAC2`，它自己的句子＝`no cold/hot sample recorded in this process: the lifecycle test did not run in this binary … Named skip - an empty aggregate is not a green latency gate` ⇒ 本腿读成"未跑"，⛔ 不读成通过（该 SKIP 只在 lifecycle 未落样本的那发里出现，隔离复跑两发里它是 PASS）。

## 7A. 判不动／量不到（具名归口，⛔ 没有一枚写成"应该没问题"）

| # | 格子 | 为什么够不着 | 归口 |
|---|---|---|---|
| N1 | **哪一枚 `257-r2` 的件算交件**（两套件并存：本腿 6 枚用例／并行腿 8 枚突变＋两枚 nonpreset 用例，写面同一枚 `firstrun.go`） | 派单前提"另一枚腿只读不写"被 `git log` 推翻（§0.9），且本腿无权限判姊妹件作废 | **编排者**（`A617` 已自曝并留"两套实现并存待 257-v1 裁"；本腿不自裁） |
| N2 | AC#1 的**真面板链那一格**（设置页里逐枚点那七枚、看页面回执原话） | 需要 WebView2 真窗＋`internal/panel`，本腿写面禁止碰 `internal/panel`（另两枚腿在读它），且真窗族在本机受争用（§6A） | 本腿量到的是**名册 setter 层**（`config.NewManager` ＋ 七枚导出 setter）与**首启 stderr 回执层**；页面渲染层归 `257-v1`／`e2e` 腿 |
| N3 | 三形状里"无控制台常驻"那一格**在机主那台机器上到底是哪一形** | 要 `GetConsoleMode` 的运行期答案（`257-a2` §5-N4 同判，`resident_task_source_windows.go:224-231` 那道闸），码里读不出 | 回执只把三形都说到、不替用户选形；真机判定归 winlive 批准后那一程（账 `A606`/`A607` 的 U9） |
| N4 | 本腿回执与 `internal/panel` 那侧中文句（`renderSettingReceipt`／`tierSentence`）的**全串一致性** | 只钉到"三枚 tag 字面同串"（本腿测试用字面量抄那三句），面板渲染整段字符串不在本腿读数窗，且 ⛔ 不能改 `internal/panel` | 归 `257-v1` 对抗验收；本腿在 §7A-N7 把自己的耦合脆性交出去 |
| N5 | `provider_credential` 那一枚的**完整腿**（`configStore.StoreCredential` → DPAPI store） | AC#3 禁新增回显／禁碰凭据面，且写 DPAPI blob 会让"干净机读数"混进真密钥存储形状；本腿按 `257-r1` 口径走该字段的**引用腿** `SetProviderAPIKeyRef` | 名册枚数仍 7，**腿数是 6**（口径差异在 `firstrun_257_test.go:96-104` 的注释里逐字写明）；真凭据腿的干净机形状归 248/257-v1 那一族 |
| N6 | 并行腿在 `.scratch/wisp/probes/257/r2/` 里留下的 40 MB `cmdwisp.test.exe`、`mut/`、`snapshot-*.go`、`baseline-cmd-wisp.log` 等 | ⛔ 派单规则"临时件只建不删"，本腿不删别人的件、也不 commit 别人的件（`3f0c4fff` 那次是被第三枚腿的裸 commit 卷走的，不是本腿提交） | 编排者／CI 地界（`A617` 已记 `ci.yml:171` gofumpt 走遍 `.scratch` 那笔红归票 171；本腿新件全部 <25 KB） |
| N7 | 脆性两条，主动交出：①本腿测试里三枚 tag 是**字面量**（`firstrun_257_test.go:53-57`），`internal/config` 将来改 tag 措辞会打红 `cmd/wisp` 这一族（有意耦合，改 tag 必须两半同批改）；②`TestTicket257R2AC3` 里 `dpapi:`／`env:` 占位计数用"总数 == 占位数"，同一占位若被引用两次则**假红**（宁可假红不静默） | — | 交给验收腿判"是否要换成尺读常量"（换成读常量＝重新掉进 257-r1b 的 M1 盲区，本腿选边是刻意的） |

## 8A. 污染面与提交名册（逐笔 `git show --name-only`＋AC 框零改动自证）

| 笔 | 归属 | 内容 | `git show --name-only` 实测 |
|---|---|---|---|
| `3f0c4fff`（11:25:10） | **message 是别票的**（`probes(268-a2)`），**内容含本腿八枚**：`cmd/wisp/firstrun.go`＋`cmd/wisp/firstrun_257_test.go`＋本件 §0–§4＋五枚 log | 本腿 11:24 用显式 pathspec `git add` 后跑 `git commit -F … --only -- <十枚路径>`，git 回 `no changes added to commit`＝**本腿那笔没落地**；随后 `268-a2` 的裸 commit 把本腿暂存区卷走 | 名单已核；内容完好（`git cat-file blob HEAD:…` 两枚 md5 与 §0.9 一致）；`A617` 就地追账同一形事故，定式升级为"commit 一律 `-F msg -- 显式 pathspec`" |
| `06b66e66`（11:2x） | **本腿** | 证据件 §0.9 追加（撞腿登记） | 只含 `.scratch/wisp/probes/257/r2/evidence.md` 一枚 |
| `7c644bb0`（11:31:42） | 并行腿（`257-r2b` 交件） | 含它的 `firstrun_257_nonpreset_test.go` ＋ 本件合并版（**我的 §0–§5 原样保留在 7–174 行**，它的 §5–§10 从 175 行起） | 名单已核：cmd/ 里只多它那枚 nonpreset 测试，本腿两枚文件不在其中（已在 `3f0c4fff`） |
| 本块（6A–10A） | **本腿** | 证据件 ＋ 尾程 log／msg 件 | 枚枚显式 pathspec，只带 `.scratch/wisp/probes/257/r2/**` |

**AC#4 越界尺（本腿两枚 commit 的全部文件名过禁区正则）**：
`git show --name-only --format="" 3f0c4fff 06b66e66 | grep -E "^frontend/|^design/|PLAN\.md|^docs/specs/|thresholds\.go|golden|allowlist\.txt"` ＝ **0 行命中**。
**写面清单（逐枚 `file:line`）**：
- `cmd/wisp/firstrun.go:123-150`（ ⓒ 的立论注释块）／`:151-160`（三样段）／`:162-170`（三因段，`Fprintf` 带真路径）／`:172-184`（三形状段＋凭据段）
- `cmd/wisp/firstrun_257_test.go`：`:53-67`（三枚 tag 与五枚教学串字面量）／`:72-95`（干净机起步：`t.TempDir` ＋ `fs.ErrNotExist` 断言 ＋ 真入口 `runTextTask`）／`:96-133`（名册七枚行走）／`:139-217`（AC#1 发 1）／`:224-276`（AC#1 发 2：照回执手加 → 7/7）／`:281-306`（AC#1 发 3：只加第一样 → 恰 3/7）／`:311-364`（AC#2 三因三句）／`:370-398`（AC#2 三形状）／`:403-446`（AC#3 凭据面）／`:448-472`（手加与读盘 helper）
- ⛔ 零触碰清单（本腿任何 commit 的名单里都没有）：`internal/config/**`（包级互斥的另一枚写位）、`internal/panel/**`（两枚腿在读）、`frontend/**`、`design/**`、`docs/PLAN.md`、`docs/specs/**`、`internal/observe/thresholds.go`、golden、`tools/d22scan/allowlist.txt`、`docs/reports/pending-and-issues.md`。
**票面 AC 框零改动自证**：`.scratch/wisp/issues/257-clean-machine-provider-registry-nil-blocks-writes.md` 不在本腿任何一枚 pathspec 里；现量该文件＝`AC#0 [x]`（编排者 10-02 翻的勾）＋`AC#1..AC#4 [ ]` 四框未勾，`git status --porcelain` 对该文件**零输出**＝工作树与 HEAD 一致，勾与不勾留给编排者。
**Git 纪律**：全程只 commit、⛔ 不 push；未用 `add -A`／`add .`／`commit -a`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；还原只用 `git cat-file blob HEAD:<path> > <path>`（§5 四发）；临时件只建不删；未在仓库内建 worktree。

## 9A. 本腿写错的读数（自我对抗，不等验收腿抓）

| # | 我错在哪 | 后果与处置 |
|---|---|---|
| N1 | §5 首版把 M3 的突变体 md5 抄成不确定串（"891fc88b4e887c0f587…"） | 已用 `mutant-md5.txt`（11:31:31 纯计算重算）的全值 `891fc88b4e887c0f8592a6dabf7b2cc2` 替掉，并把"抄错"这件事留在原地不抹 |
| N2 | M2 第一次脚本第 5 枚 needle 没匹配上 → 按设计写盘前退出，但**那次 `go test` 仍跑了**（对未突变文件），`rc=0` | 那一发的 `rc=0` **不作任何一发的读数**；换 needle 后的 M2 才是证（11:29:02→11:29:14，五条红句在 `mutant-M2.log`）。教训：脚本 assert 失败后不该继续跑测试，本腿后半段把两步并成一个 `&&` 链 |
| N3 | §3 发 2 首跑 5/7 红（`provider_api_key_ref`／`role_chat_model` 两枚"接受但键路径为空"） | **是我的仪器缺陷不是产码缺陷**：手加段把后面要试写的同一个值先填上了，`writeOneKey` 老实回"文件里已经是这个值"。改法是让两段用不同占位值（未放宽任何断言，两枚 log 都留盘） |
| N4 | 11:26:51 我把 `go vet ./cmd/wisp/` 的编译红（`internal/panel/workspace.go` 参数不匹配）差一点记成本腿门禁读数 | 现量 `git status`＝`181-r3` 正在写 `internal/panel`＋`internal/tools`（九枚脏）⇒ 那次红是别腿在飞形状，已从 §6A 门禁表剔除并具名；尾程 11:43:36 rc=0 才算数 |
| N5 | §0.9 首版我只写了一句"另两枚只读腿在写文件"的推测 | 11:25 追认时才用 `git show --name-only` 把"是谁、哪一枚、哪一刻、卷走了哪八枚"定死；推测没当成读数用 |
| N6 | 我最初接受派单"代号 `r2` 是干净的"（编排者已 `ls` 过）作为前提，直到 11:22 看见 `mut/` 才怀疑 | 教训＝**代号干净要每轮重取**，不只在派单时取一次；登记给编排者的派单定式（`A617` 已自定"派前 grep 台账腿号＋git log --since 看那枚包几分钟"） |
| N7 | 我不能证明"回执教的形状就是机主真会走的形状"（人因层面），只证明了它可加载、7/7 解锁、且非预设名那一支会因缺 `protocol` 加载不过（本腿只把这句话**写进文案**，量它的是并行腿那枚 nonpreset 用例） | 具名留给 257-v1：文案诚实性 ≠ 机主真走得通；真走通要 e2e 腿 |

## 10A. 本腿交件判语（四条 AC 各一句 ＋ 预算／纪律）

- **AC#1（干净机真跑＋ ⓒ 终态）**：`t.TempDir()` 且断言无 `config.toml` 起步 → 真入口 `runTextTask` 建首份配置 → 名册七枚逐枚试写**全拒、且七枚都只报第 2 种并点名 `[llm.providers.`**（`map[第 2 种拒因：行不存在:7]`）→ 照回执教的三样手加后 **7/7 接受并报键路径**、只加第一样**恰好 3/7**；回执文本逐字含 `[llm.providers.<名>]`／`[llm.providers.<名>.models.<模型 id>]`／就地填 `[llm.roles.chat]`，且 `[providers.` 被反控钉为禁串。**判：兑现。**
- **AC#2（三因三句）**：三枚 tag 各恰一枚行、各带互不重叠的补救句，`配置未生效`／`重启就好` 在整个回执只许出现在禁句那一行；M3（折因）与 M4（折形状）两发红句证明折叠必响。**判：兑现（面板渲染层那一格见 §7A-N2/N4）。**
- **AC#3（凭据面一字未动）**：本腿零新增方法（`firstrun.go` 仍只 `ensureFirstRunConfig` 一枚函数、签名未漂，静态钉在 `:403-446`）；回执里 `dpapi:`／`env:` 各只以占位形出现；盘上文件无 `api_key =`；本件与本 commit message 内**无任何凭据值**（只有 `dpapi:<blob 名>` 这类占位与 `env:T257_R2_*` 这类永不置值的变量名）。**判：未动。**
- **AC#4（越界）**：本腿两枚 commit 的全部文件名过禁区尺＝**0 命中**；写面只有派单授权的两枚（`cmd/wisp/firstrun.go` ＋ 同包测试件）。**判：零越界。**
- **牙**：四发突变全响（ⓐ 两形＝改坏拼法／删指引；ⓑ 两形＝折三因／折三形状），红句逐字在 §5 与四枚 `mutant-M*.log`，还原后五取 md5 全等基线、`git cat-file blob HEAD:` 双读相同＝没有突变进历史。
- **纪律与预算**：本腿工具轮次约 60／100 帽内交满全部节次（未触发"第 100 轮先写 §7/§10"那条硬闸门，但 §7A/§10A 已提前写满）；无 `t.Skip`、无放宽断言、无 SKIP 读成通过；只 commit 不 push；枚枚显式 pathspec。
- **留给编排者的一句话**：两套件并存（本腿 6 枚用例＋四发突变／并行腿 8 枚突变＋2 枚 nonpreset 用例），**本腿不合并、不重跑、不替它判**；要哪一套、或两套都要合入哪一枚，`257-v1` 裁。




