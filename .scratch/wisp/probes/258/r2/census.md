# 票 258 · 实现腿 `258-r2` 证据件（写码腿，非验收者）

起手时刻 `2026-10-04 09:01:48 +0800`（自取 `date`）；起手锚 HEAD `fd269de1`（`git log --oneline -1` 现取）。
本件只建不删；读数逐字；AC 框一枚未碰。

**本腿射程（派单三件事）**：
① AC#1 条件①：`[hotkey]` 节缺失时档位词必须印 "defaults"（措辞说真话），并配一枚不依赖真窗、
   对「把 defaults 改回 config」这一发有红的用例。
② AC#2 条件：两枚 winlive 尺子修形状（rebind 枚的 JSON/等号格式失配、occupied 枚的环境赌博），
   改成"读得到就判、读不到就具名失败/诚实 skip"。⛔ 本轮不跑 `-tags winlive`。
③ build 销账：全量 `go build ./...` 逐字输出＋rc 落 §3。

## §0 起手基线与绿名册

- 起手 HEAD：`fd269de1`。
- 基线整包（改码前，2026-10-04 09:05–09:14）：档 `.scratch/wisp/probes/258/r2/baseline-gotest.txt`。

| 项 | 读数 |
|---|---|
| rc | `1`（整包行逐字：`FAIL	github.com/CarlosShao/wisp/cmd/wisp	296.857s`） |
| PASS/FAIL/SKIP | `--- PASS` 行 291（含子测）／`--- FAIL` 行 1／`--- SKIP` 行 0 |
| 基线唯一红 | `TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card (41.xx s)`＝C18 窗挤压间歇族（v1 §0 同款；票 258 链零交集，本腿未修不判） |

- 基线里含 hotkey/ball/258 的绿用例名册（逐字抄自 baseline-gotest.txt，本腿动语义前实拍）：
  `Test258AssemblyRootWiresTheChainAndTheBridge`、`Test258BridgeMutationNoSrcKeepsOldBinding`、
  `Test258BridgeRebindsLiveKeysFromConfigEdit`、`Test258ConstructionChainMissingFileFallsBackAndSaysIt`、
  `Test258ConstructionChainTakesConfigValues`、`Test258OccupiedCombinationNamesTheNewValue`、
  `Test258ProvenanceWordsAreThePrintedOnes`、`Test258V1ProbeConstructionMissingFileNamesTheFallback`、
  `Test258V1ProbeSummonEditRebindsLiveBall`、`TestAC228BallHostAnswersEveryGesture`、
  `TestAC228ResidentLegIsTheBallHost`、`TestAC228ResidentLegReportsAndBooksItsBall`、
  `TestAC246EscChannelStaysUnloadedWithoutABallWindow`、`TestBallGestureWithoutPanelHostStillRecords`、
  `TestBallPanelGesturesReachThePanelThread`——全 PASS。
- ⚠ 共享树事实：起手时 `cmd/wisp/config_hotkey_source_260.go` 与 `cmd/wisp/resident_hotkey_chain_260_windows.go`
  是别的腿（票 260）的盘上未跟踪件，09:06 在场、09:17 后消失（本腿零触碰、未提交）。基线与终跑是否把它们编进了
  测试二进制取决于各自编译瞬间；两发读数内部自洽（同一颗 C18 间歇红），如实存照。
- 票 245 那族"裸 Esc 只在确认那两三秒借"的行为钉（`TestAC246EscChannelStaysUnloadedWithoutABallWindow` 等）
  基线绿、终跑绿，本腿改动（chain 措辞规则＋尺形状）不触 cancel/borrow 语义：chain 三形返回值都是
  `DefaultHotkeys()` 本尊（cancel 恒 "Esc"），`ApplyHotkeyDefaults`/`escBorrow` 零改动。

## §1 ①：档位词裁定形状（defaults 说真话）

**规则**：档位词跟**合并后的集合**走，不跟"文件可不可读"走——
- 集合有任何一格被文件挪离编译默认 ⇒ "config"；
- 集合逐格等于 `ball.DefaultHotkeys()` ⇒ "defaults"，覆盖三形：文件缺失/不可读（视图全空）、
  文件在而 `[hotkey]` 节缺失（schema 只带 `cancel="Esc"` 标签，`ApplyHotkeyDefaults` 补出的就是默认表）、
  节在而逐字写默认值（终值就是默认表，"config" 会把没出过力的文件写成出处——v1 §3 点名的假话）。
- **不许静默换成另一套＝对值的约束**：三形返回值都是 `DefaultHotkeys()` 本尊，本腿只改词、零改值形状。

| 落点 | file:line |
|---|---|
| 判据本体（chain 纯函数，无窗可跑） | `cmd/wisp/resident_ball_windows.go:191-202`（`if cfg == ball.DefaultHotkeys()` 在 `:196`） |
| 装配根审计 slog 措辞随行（"bindings taken from" 在节缺失时就是那句不该印的话） | `cmd/wisp/resident_windows.go:185-188` |
| 新用例（不依赖真窗；进今天 CI） | `cmd/wisp/resident_hotkey_258_windows_test.go:189-242`＝`Test258SectionMissingTierWordIsDefaults`：前提钉（节缺失视图非全空、只有 schema cancel——旧代码正是从这里滑进 config 分支）＋ defaults 判词＋**正控**（同一目录种一格 `Ctrl+Alt+Z` ⇒ config）＋边界钉（逐字默认值节 ⇒ defaults） |
| 旧用例翻转（其第三分支钉的就是被裁掉的假话） | `cmd/wisp/resident_hotkey_258_windows_test.go:138-180`（`:158`/`:175` 两条断言现在要 "defaults"） |
| 真窗两枚的 boot 种值改植非默认 `Ctrl+Alt+Z`（"hotkeys from config" 从巧合变读数） | `cmd/wisp/resident_hotkey_258_windows_test.go:266`（Bridge，boot VK 断 'Z'）、`:373`（Occupied，boot/编辑两处 body 同改，diff 只剩 panel 一格） |
| 255 花名册证据行随链函数变长改引 | `cmd/wisp/config_readers_255.go:115`、`:123`（`:269 → :276`，token `Hotkeys:  cfg,` 现读复量在 `cmd/wisp/resident_ball_windows.go:276`） |

## §2 ②：两枚 winlive 尺子的形状修复（⛔ 本轮未跑 winlive）

**(a) rebind 枚**（`cmd/wisp/resident_hotkey_live_258_windows_test.go:120-177`，258-v1 §5 实测：
重绑行 `{"msg":"ball: hotkeys rebound after config change","summon":"Ctrl+Alt+R",...}` 就在失败附件里，
判据却 grep 等号形永不命中）：
- 判据改走 `sinkRebindLineHas258(tail, summonValue)`＝msg 短语在 **且** 冒号形 `"summon":"V"` 或等号形
  `summon=V` 至少一形在（`cmd/wisp/resident_hotkey_258_test.go:90` 起）。读不到＝维持具名 `t.Fatalf`
  （12s 有界轮询后把 sink tail 逐字贴失败件；`readSinkTail258` 在 jsonl 未落时把 "(no wisp-*.jsonl ... yet)"
  写进 tail——"读不到"本身在读数里可见，不会静默成绿）。
- 该 matcher 由**不跑真窗的 CI 用例**钉住：`cmd/wisp/resident_hotkey_258_test.go:107` 起
  `Test258SinkRebindRulerReadsBothSpellings`（JSON 形真、等号形真；反形＝旧值 Z 的同一行、无 msg 的行、
  空 tail 全假）。空 tail 反证钉死"永远满足"的换形——见 §4。
- 顺带把枚的 boot 种值从默认 Q 改成非默认 Z 并断 `hotkeys from config`＋`summon=Ctrl+Alt+Z`
  （原形在 M3-恒退defaults 腐坏下 boot 断言分辨不出——换形不红的洞就地补掉一半，另一半见 §5）。

**(b) occupied 枚**（`cmd/wisp/resident_hotkey_live_258_windows_test.go:193-237`，v1 §7-R7：赌
`Ctrl+Alt+U` 被第三方占、无兜底 ⇒ 键空着时假绿/假红全看机器心情）：
- 前提改**自建**：`squatHotkeyForRuler258`（`:264-291`）在 `bootResidentLeg` 之前于**本测试进程**
  的锁定 OS 线程（`runtime.LockOSThread`＋`user32!RegisterHotKey`，hWnd=0、id `liveSquatID258`）注册该组合：
  注册成功 ⇒ 占用者是本进程（`t.Cleanup` 同线程 `UnregisterHotKey`）；答
  `ERROR_HOTKEY_ALREADY_REGISTERED` ⇒ 第三方在场，前提由它成立、无需释放；**其余任何答 ⇒
  `t.Fatalf` 具名"建不起前提不判"**（三形各有名，`owner` 描述进失败/日志句）。
- 判词随之加牙：Problems 行（`occupied by another program`＋`panel = "Ctrl+Alt+U"` 逐字）之外，
  新钉 `hotkeys live 2/4`（前提成立时算术上没有第三枚 live 的可能；v1 那次 `3/4` 只是 Logf 软记，现改硬断言）。
- 诚实 skip 只剩一枚且沿用家族形：主机无球窗（`ballAbsentClaim`）⇒ `SKIP-LOUD`；winlive tag CI 不跑
  （258-a2 §6(c) 尺），不会把 SKIP 喂给 runtests.sh。
- 前提构建器本体：`cmd/wisp/resident_hotkey_live_258_windows_test.go:264-291`。

**线程形状**：两枚改的全是测试件；产码侧 rebind 仍走球自己 ui-sta 的 `uiRun` post-and-wait，
`startResidentBall` 桥装配与 D38(e) join 次序零改动——票 33 十一裁射程未触，无上报项。

## §3 门禁读数

| 尺 | 逐字末行 | rc |
|---|---|---|
| 基线整包 `PATH=... go test -count=1 -v ./cmd/wisp/`（09:05–09:14，改码前） | `FAIL	github.com/CarlosShao/wisp/cmd/wisp	296.857s` | `1`（红名册＝§0 唯一 C18 间歇） |
| 终跑整包（同尺，改码后；档 `family-final-gotest.txt`） | `FAIL	github.com/CarlosShao/wisp/cmd/wisp	243.333s`；`--- PASS` 行 **293**（含子测；顶层 `^--- PASS` 204）／`--- FAIL` **1**／`--- SKIP` **0** | `1` |
| 终跑唯一红复量 | `TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card (41.05s)` 整包红；**solo `-run ...$` 复跑＝`--- PASS (1.29s)`、rc=0**（档 `c18-solo-rerun.txt`）＝挤压间歇形，基线同法同红——非 258-r2 回归（判据只涉 C18 always-branch 审批窗，与热键链零交集） | — |
| 258 家族 11 枚终跑（9 旧＋`Test258SectionMissingTierWordIsDefaults`＋`Test258SinkRebindRulerReadsBothSpellings`） | 全 `--- PASS`（逐字行在 family-final-gotest.txt，09:17 一发 `-run Test258` 亦 rc=0） | `0` |
| `go vet -tags winlive ./cmd/wisp/`（只编译，**未跑任何 winlive 用例**） | 无输出；逐字记法 `rc=0`（档 `vet-readings.txt`，09:19:20） | `0` |
| `go vet ./cmd/wisp/` | 无输出 `rc=0` | `0` |
| **`go build ./...`（仓根全量，交付件③）** | 逐字：`## 258-r2 build 销账：go build ./... 于仓根`／`09:19:25+0800`／（**零行输出**）／`rc=0`（档 `build-full.txt`） | `0` |
| `scripts/d22scan.sh` | `d22scan: clean - no D22 ban violations; ... ban #8 cmd/=97`；runtests.sh 全仓行 `top-level: PASS=34 FAIL=0 SKIP=0`；`tools/d22scan ok 17.522s` | `0` |
| gofumpt | PATH 上命中一枚 gofumpt（`command -v` 分支），`-l` 对本腿六枚文件输出空 ⇒ 全已合规 | `0` |

## §4 变异自证

**①真跑（go test 可跑的那发）**：把 `cmd/wisp/resident_ball_windows.go:199` 的
`return cfg, hotkeyProvenanceDefaults` 换回 `hotkeyProvenanceConfig`（＝翻勾节点名的那一发腐坏形），
`-run "Test258SectionMissingTierWordIsDefaults|Test258ConstructionChain"` 实测红 **4 条**，逐字：
```
resident_hotkey_258_windows_test.go:158: missing-file provenance = "config", want "defaults": AC#1 requires the fallback to be named, not silent
resident_hotkey_258_windows_test.go:175: section-missing provenance = "config", want "defaults": the [hotkey] table is absent, the set is the compiled defaults
resident_hotkey_258_windows_test.go:212: AC#1 condition-1 RED: section-missing provenance = "config", want "defaults" (the merged set IS ball.DefaultHotkeys(); "config" names a file that contributed no binding)
resident_hotkey_258_windows_test.go:240: a verbatim-defaults [hotkey] section answered "config": the set equals ball.DefaultHotkeys(), so only defaults is true of it
```
（rc=1；档 `mut-word-config-run.txt`。）还原自备份 `resident_ball_windows.go.r2backup`（cp 覆回＋grep 复量
`:199 defaults`/`:201 config` 在位）⇒ `-run` 复跑 `ok ... 0.122s` 绿。

**②两枚尺子（真窗禁跑 ⇒ 静态反证，按派单不喊"有牙"）**：
- rebind 枚：把 `rebound = sinkRebindLineHas258(sinkTail, newSummon258)` 换成"永远满足"形
  `rebound = true` ⇒ 循环第一转即退出、Fatalf 永不触达——**真窗下桥被摘（M1 形）也会绿**，
  因为读盘那一环被跳过、判据不再消费任何进程产出。这条"换形不红"正是派单禁用的形状；现在承重落在
  matcher 身上，而 matcher 的反形（空 tail、非重绑 msg、旧值行）在 §2 的 CI 钉里**实测为假**——
  也就是说"永远满足"换形若被装入，`Test258SinkRebindRulerReadsBothSpellings` 当场红（第三条反证断言
  直接判空 tail）。残余不敏感面如实具名：matcher 是 raw substring 读数，一台**每 tick 无 diff 也谎写
  重绑行**的腐坏进程会骗过它——那需要真窗下的行为读数，本腿量不到（§5）。
- occupied 枚：把 Problems/live 断言删成软 Logf（＝旧枚的"赌环境"形）⇒ 键空着的真机上
  `panel=Ctrl+Alt+U live`、无 Problems 行，枚**静默绿**（v1 §5 亲量的正是这个假象的镜像：那次红得
  冤枉）。现形里前提由 RegisterHotKey 三形之一**当场建立或具名拒判**，"读不到"不存在第三种静默出口。
  残余不敏感面：若产码谎报 occupied 行而 Win32 实际注册成功——live 2/4 断言由 verdict 算式（`len(rep.Live())`）
  兜住，但"谎报+计数也假"这一整层要真窗复跑才能证伪，同记 §5。

## §5 判不动的地方（欠账具名）

1. **修完的两枚 winlive 尺子有没有牙＝本轮量不到**（⛔ 派单禁跑 `-tags winlive`，真窗会在 owner 在用
   的桌面真开球窗真抢全局热键）。已证部分＝编译过（vet rc=0）＋判据形状静态可读（§2）＋matcher 反形
   CI 实测（§4）。**欠账**＝squat 形 occupied 枚首发、修形 rebind 枚首发、以及两枚各一发 M1-摘桥树突变
   （⚠ 方法照 v1 §8-3：overlay 不穿透 `buildWispForTest` 的子 `go build`，必须真树突变＋md5 还原；
   本腿备份件在 `.scratch/wisp/probes/258/r2/`，突变腿自己再核）。归属＝**258-v2 之后与 owner 同意的窗口**，
   本腿不自行绕跑。
2. **恒真问的"读装配根"那一发**（v1 §1 M0/M3 默认档绿＝6 枚主判据自带闭包注入不读装配根）——翻勾节
   已裁"不动判据、留给后续程"，本腿未动；§2(a) 顺带把 winlive rebind 枚对 M3 形变敏感（boot 断
   `hotkeys from config`＋非默认种值），**默认档六枚的覆盖洞原样在册**。
3. **票 260 射程**（配置里的 cancel 组合键从来没注册过／借键不吃配置）：本腿零触碰——chain 三形的 cancel
   恒 `Esc` 不新增注册；起手时盘上那两枚 260 未跟踪 Go 文件非本腿所写、未随本腿提交（§0 存照）。
4. **C18 间歇红**（`TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card`）：基线与终跑同发倒下、solo 均
   回绿，属既有挤压族（v1 R2 归因换名后的同判决），非本腿回归；修它不在 258-r2 射程，不动。

## §6 越界与禁区自查

- `git show cbece45e --name-only`＝**六枚全在 cmd/wisp/**：`config_readers_255.go`、
  `resident_ball_windows.go`、`resident_hotkey_258_test.go`、`resident_hotkey_258_windows_test.go`、
  `resident_hotkey_live_258_windows_test.go`、`resident_windows.go`。
- 禁区逐名过：`docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／
  `tools/d22scan/allowlist.txt`／`frontend/**`／`design/**`／三枚冻结件／`internal/**`／`tools/d22scan/**`
  ＝**零出现**（本腿 diff 内无任何仓根外路径；写面纪律包级互斥：212-r2 的三目录未碰）。
- 票面 `- [ ]`/`- [x]` AC 框：未编辑过票文件（`.scratch/wisp/issues/258-*.md` 不在任何 commit pathspec 里）。
- 线程模型：产码改动仅 chain 纯函数＋slog 措辞＋测试件；`uiRun`/桥装配/join 次序/`Run()` 泵零改动。
- d22scan：clean（§3 行）；无新增 `go func(`、无 PathResolver 外 Clean/Abs、无墙钟超时、字符串无 U+2713/U+2264。
- Git：三枚 commit 均显式 pathspec、无 `-A`/`.`、无 amend/reset/rebase/stash/checkout ./clean；未 push。

## §7 Progress log

- 09:0x 起手：date/锚现取；票面翻勾节＋258-v1 判决书＋落点五件现读；基线整包后台跑（09:05 编译快照）。
- 09:0x 证据件骨架 commit `4d8763b5`（七节表位）。
- 09:1x ①落地（chain 集合词规则＋无窗用例＋真窗两枚改植 Z＋255 花名册改引 :276＋装配根 slog 措辞）；
  ②落地（sinkRebindLineHas258＋CI 钉＋squat 自建前提＋live 2/4 硬断言）；vet 双 tag rc=0；gofumpt 空。
- 09:1x 变异①实跑红 4 条（§4 逐字）→ 还原绿；`-run Test258` 11 枚 rc=0。
- 09:19 门禁：`go build ./...` rc=0（档 build-full.txt）；`go vet -tags winlive` rc=0；d22scan clean。
- 09:2x 终跑整包 293/1/0（唯一红＝基线同款 C18 间歇，solo 复跑 PASS）；代码腿 commit `cbece45e`（六枚文件）。
- 09:3x 本件终态＋读数档随证据 commit `c74ab57e`（九枚档，显式 pathspec，warning 仅 autocrlf 换行提示）。
  本腿 commit 链：骨架 `4d8763b5` → 代码 `cbece45e`（六枚 cmd/wisp 文件）→ 证据终态 `c74ab57e`；
  最后一枚 census 收尾行 commit 只增这一行，链尾＝见 git log。未 push。
