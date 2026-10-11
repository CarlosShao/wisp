# writeface-1 — 90 final（六节）

锚＝HEAD `934b350a` 起手；本腿共 **4 笔**（`5f2289e3`／`47e0ff78`／`a900bc10`／本笔），全程只读，写面只有 `本目录/**`。正文总量＝5 枚 `.md`、`wc -l 本目录/*.md` ⇒ 交本笔前 **254 行**（≤400 预算内）。

## ① 交了几笔 commit

| 笔 | sha | 名册枚数（尺＝`git show --name-only --format= <sha> | grep -c .`） | 目录外枚数（尺＝同尺 `grep -vc '^\.scratch/wisp/probes/orch/20261011/writeface-1/'`，应 0） |
|---|---|---|---|
| 1 锚件 | `5f2289e3` | **1** | **0**（现量已跑） |
| 2 名册＋颜色＋logs | `47e0ff78` | **39** | **0**（现量已跑） |
| 3 互斥图＋本件＋4 枚新 logs | `a900bc10` | **6**（`30-…md`／`90-…md`／`logs/atomic-bool.txt`／`golist-all.txt`／`new-since-anchor.txt`／`started-internal-audio.txt`） | **0**（尺现跑） |
| 4 本更正笔 | 本笔 | **1**（只有 `90-final.md`） | **0**（尺＝`git show --name-only --format= HEAD | grep -vc …`，交完复跑；⛔ `--amend`，按规矩**追加**） |

★共享工作树实证：笔 2 与笔 3 之间**别家又落了 1 笔**（`24929768 docs(308-a1)`），本腿两笔名册仍**目录外 0 枚**（尺＝显式 pathspec 那一把）⇒ "⛔ `add -A`、每笔带 pathspec"在真并发下有效；收尾复跑 `git status --porcelain -- internal cmd`＝**0 行**、本目录 0 字节件＝**0 枚**。

- 0 字节件复跑尺＝`find 本目录 -type f -size 0` ⇒ **0 枚**（曾有 3 枚 grep 零命中件是 0 字节，已逐枚补"命令＋rc＋为何 0"使其⛔ 为空；见 `logs/zero-byte-check*.txt`）。
- ⛔ 证据件命名 `.out`（根 `.gitignore` 第 8 行全仓 `*.out`＝静默不收）；本腿件名一律 `.md`／`.txt`／blob 副本 `.go`／`.sh`／`.yml`。
- ⛔ `push`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`／`--no-verify`／`add -A`／`add .` 命中 0；每笔带显式 pathspec＝本目录。

## ② 每格答了什么（一句一条）

1. **逐腿文件名册**（`10-…`）：`296-r2`＝**1 枚**且是**产码文件** `cmd/wisp/resident_windows.go` 的 `hotkeyReloadSource296` 文档块（三处锚 `:38`/`:40-41`/`:42`）；`295-r1`＝**4 枚**（`resident_audio_windows.go` 4 处含声明＋`resident_mute_290_windows_test.go` 4 读点＋`resident_audio_247_windows_test.go` 3＋`resident_audio_247_live_windows_test.go` 1）；`294-r1`＝**1–2 枚新增测试件**＋1 枚 `hkMute…` 本地镜像常量（只**读** `resident_windows.go` 的 AST）；`309-r1`＝**2 枚**（`internal/tools/bridge.go` 注释＋1 枚新增 `internal/tools/*_test.go`）；`298` 残余＝**0 枚**。
2. **颜色面**（`20-…`）：`296-r2` 定向 `-run 'Test296|Test258'`；`295-r1` 要 `./cmd/wisp/` ＋ `./internal/audio/`；`294-r1` 要 `./cmd/wisp/` 整包；`309-r1` 要 `./internal/tools/`；⇒ **三枚 cmd/wisp 腿共用同一枚包名册**，309 与它们⛔ 同包。机器需求归口 6 条 M-1..M-6（真桌面会话／机主真麦／`-race`／CI-only／历史台面严格形）。
3. **互斥图**（`30-…`）：10 个配对里**⛔ 一对是"同一枚文件的写-写"**；真边只有两条软边＝`go build ./...` 整模块编译（所有腿共享）＋`./cmd/wisp/` 名册归因（296/295/294 三枚互相）。波次＝**A：`296-r2` ∥ `309-r1`** ⇒ **B：`295-r1`** ⇒ **C：`294-r1`→`294-v1`**；若派单钉死 ①②③ 三条，B/C 可并。
4. **顶回**（§⑥）：**5 条**具名（其中第 1/3/4 条＝编排者预先点名那两枚雷＋那条他自己写的串行约束；"⛔ 硬造"这一条没触发——真找得到）。

**复跑对拉结论**（要本腿复跑的承重尺，逐条；✓＝对上，✗＝打旧）：
- 295-a1b ⓓ 的 `.started`＝**11** ⇒ ✓ 复现（分布 3 产／8 试）。
- 309 裁定节的 `ToolRequest{` 60＝59 测试／1 产码 ⇒ ✓ 复现，并补一枚它没写的维度：同尺去 pathspec＝全仓 tracked **252**。
- 票 295 现量"全仓已有 **5** 处在用 `atomic.Bool`" ⇒ ✓（尺＝`git grep -nE "atomic\.Bool" HEAD -- internal cmd`＝5）。
- 票 294 裁定节"测试直调 **3**" ⇒ ✗ 现量 **5**；"手搓 `&residentBall{` **13**" ⇒ ✗ 现量 **16**；产码 `:336` ⇒ ✗ **`:342`**；定义 `:472` ⇒ ✗ **`:478`**。四条同一成因＝**票 293 的落地腿 `resident_tray_mute_293_windows_test.go` 已进 HEAD**（那枚件带 2 枚直调＋2 枚手搓），⛔ 改原句，只追加。
- `startResidentBall` "测试真调用 10" ⇒ ✓；"整族 12" ⇒ **射程差**：本腿尺 `-- cmd/wisp` 给 **14**＝测试内 12（含注释 2）＋测试外 2（定义＋生产调用）⇒ 两把都对，今后名册必写含⛔ 含定义与生产点。

## ③ 欠读数（⛔ 本腿代答，逐条写"该由谁用什么尺答"）

| # | 欠的是什么 | 该由谁、用什么尺 |
|---|---|---|
| U-1 | 三枚 cmd/wisp 腿同台面并发时，逐名名册**实际**互相搅出几枚假红/假绿（本腿只证"文件面∅"，⛔ 证"颜色面⛔ 会串"） | 编排者在派单里指定"每波独占 `./cmd/wisp/` 色"，落地腿 `go test ./cmd/wisp/ -count=1 -v` 两发取交集＋逐名作差 |
| U-2 | `-run 'Test296\|Test258'` 在本机桌面会话下真发 `--- PASS`（⛔ 退成 SKIP） | `296-r2` 起手首发断言具名用例有名（票 309「边界」⏂ 那条轴同形） |
| U-3 | `gofumpt` 残留 1 枚（`panel_transport_35r2_test.go`）的 HEAD 现势 | 下一枚碰 `cmd/wisp` 的门禁腿并排两把（工作树⇄blob）复现；本腿⛔ 跑 gofumpt，现引＝`298-v1`/`298-r1` 的旧读数 |
| U-4 | `294-r1` 会不会选择"把 `AC#3` 断言加进 `resident_mute_290_windows_test.go`"（那一选＝与 `295-r1` **当场**文件互斥） | ⛔ 静态读不出来＝落地腿的决策；编排者在派单钉一句"写面限自家新件"即可消解 |
| U-5 | 票 295 `AC#4` 第二包 `./internal/audio/` 的必要性（`residentAudio.started` 那发⛔ 触到 `internal/audio` 任何符号：尺＝`git grep -nE "\bstarted\b" HEAD -- internal/audio`＝**14 命中，全是该包自己的 `started`**，与 `cmd/wisp` 那枚同名⛔ 同物） | 编排者裁：要不要砍掉第二包那一发以省机时；⛔ 由写腿自删票面判据 |
| U-6 | `309-a1` 自报的 **7** 枚欠读数（本腿⛔ 重述⛔ 代答） | 归口逐枚写在台账 `A844` §3（编排者已结 U-3、结一半 U-4）；件＝`probes/309/a1/90-final.md` |
| U-7 | 票 294 `AC#2` 那枚"cmd/wisp 本地镜像常量"的最终命名与放置（裁定只说"缺它"） | `294-r1` 落地时具名；先例逐字 `	hkSummon258 uint32 = 1` |

已**自结**一枚（用本腿允许的那把尺，⛔ 算欠）：`.scratch/**` 下的 `.go` blob 副本会⛔ 进 `go build ./...` 的包枚举＝尺 `go list ./...` ⇒ **35 枚包，含 `.scratch` 者 0 枚**（点开头目录⛔ 被通配展开）⇒ 编排者 A791 那句"blob 解到仓内会被 `go build ./...` 扫到"对 **`.scratch/**` 这一形⛔ 成立**（真风险面是"非点开头目录"，如仓根）。

## ④ 本腿跑过的编译面命令（自报）

| 命令 | 条数 | 事实 |
|---|---|---|
| `go build` | **0** | ⛔ |
| `go vet` | **0** | ⛔ |
| `go test` | **0** | ⛔（`298-v1` 当时在飞用 `-overlay` 取色；本腿⛔ 碰） |
| `gofmt` / `gofumpt` | **0 / 0** | ⛔ |
| `sh scripts/d22scan.sh` / `scripts/lenbudget.sh` | **0 / 0** | ⛔ |
| `go list` | **1** | `go list ./...`（→ `logs/golist-all.txt`，35 行，rc=0；为结 §③ 那枚自读数） |
| `go env` | **0** | 未跑（⛔ 需要 `$(go env GOPATH)` 那一把，那是门禁面） |
| 其余 | — | `git`（`show`/`grep`/`ls-tree`/`status`/`log`/`diff --stat`/`merge-base --is-ancestor`/`rev-parse`/`cat-file`/`add`/`commit`）、`sed -n`、`grep`/`cut`/`sort`/`uniq`/`wc`/`find`/`ls`/`date`、`mkdir`、`printf`（只写自家 `logs/` 三枚说明件） |

## ⑤ 本腿⛔ 做、但射程边上有想法的事（一句话一条，⛔ 动手）

1. 票 294 现量/裁定那五条数字（3/13/:336/:472/整族 12）已随 293 落地整体过期 ⇒ 下一份 `294-r1` 派单该由编排者**现跑后重写落点表**，⛔ 让腿带着过期号上路。
2. 裁定节点名的 `cmd/wisp/resident_mute_294_test.go` 这枚**无 `_windows` 后缀**的文件名对 `AC#1`（纯 AST）合法，对 `AC#2`/`AC#3`（要用 `residentAudio`/`toggleMute`/球窗符号）**会**在 Linux 台面上打死编译 ⇒ 建议拆成两枚件（AST 一枚＋`_windows_test.go` 一枚），本腿只标风险⛔ 改裁定。
3. 票 309 的枚数钉若按"全仓"射程写，分母今天＝252 且**任何一枚腿**新交一份含 `ToolRequest{` 的探针件都会挪它 ⇒ 建议钉到 `internal cmd` ＋排 `_test.go`（现量 60/59/1 那把）。
4. `298` 的残余里没有"把 `-race` 请进门禁"那一枚（那归票 295 `AC#5` ＋人工批准）⇒ 别把它算进 298 的欠账。
5. 本腿自己给 `.scratch/**` 加了 **5 枚 `.go`＋1 枚 `.sh`** tracked 副本（现量 `git ls-files '*.go'`＝**992**）⇒ 今后凡引"22 枚 `.scratch` 故意脏拷贝"那类名册，要写明含⛔ 含探针目录的 blob 副本（本腿副本逐字＝HEAD blob，⛔ 新增脏枚）。
6. `resident_hotkey_296_windows_test.go:22` 那行注释在**转述** `hotkeyReloadSource296` 的来历；`AC#5` 改产码注释时⛔ 顺带扫它一眼——本腿尺（`grep -nE "only skip|dead until|turned off|DeepEqual"`）在它身上 0 命中，所以**今天⛔ 同步义务**，但那是"一行的事"，下一位改注释的人值得知道。

## ⑥ 顶回编排者（具名）

1. **雷#1 的前提⛔ 成立**：编排者写"票 296 `AC#5` 的注释住在 `cmd/wisp` 的**测试件**里"。现量＝三处锚全在**产码文件** `cmd/wisp/resident_windows.go` 的 `:38`／`:40-41`／`:42`；第二把尺（反证）＝`git show HEAD:cmd/wisp/resident_hotkey_296_windows_test.go | grep -nE "only skip|dead until|turned off|DeepEqual"` ⇒ **0 命中**。⇒ 连带：`296-r2` 与票 298 那两枚测试件**⛔ 同名册交集**，两件事⛔ 是同一枚文件。
2. **`296-r2` × `295-r1` 的串行⛔ 是文件必然**：写面交集＝∅（一枚文件名单见 §②-1）。把它们绑在一起的只有 `go build ./...` 与 `./cmd/wisp/` 名册归因——那是**颜色归因**约束，⛔ 是"碰同一枚文件"。
3. **`295-r1 先于 294-r1`⛔ 是文件必然**（编排者要我当面顶的那一条）：`294-r1` 的**入库**写面⛔ 含 `resident_audio_windows.go`——它的 `AC#1` 只 **parse**、`AC#3` 的反形是**瞬态变异**（该形已有两枚仓内代偿法：`-overlay`〔`298-v1`〕、仓外备份＋`sha256sum` 还原证明〔`296-v1`〕），而它断言的 `executed` 与 `295` 换形状的 `started` **零耦合**（尺＝`git grep -n "executed" HEAD -- cmd/wisp`＝**19 命中**；那枚布尔的**唯一断言**逐字 `	if executed {`（`resident_mute_290_windows_test.go:197`），生产侧只有 `resident_ball_windows.go` 的 `:89` 签名／`:605` 绑定／`:606` 判，其余全是注释与一句 `Printf` 文案；`executed`⇄`started` 两个名字**⛔ 同现于一枚行**）。这条串行的真出处＝**票 294 排程节那句"与票 295（同写 `cmd/wisp/resident_audio_windows.go`）不同批"**——那句写于裁定把 `AC#3` 落点改成"新测试件"**之前**，本票自己的裁定节已把它作废，编排者引的是作废那半句。
   ⚠ 本腿⛔ 主张"所以随便并发"：并发要派单钉三条（见 `30-…` 表第 3 行①②③）；⛔ 愿意钉 ⇒ 维持串行是**安全**的，只是⛔ 该把它说成文件必然。
4. **雷#2（"照抄 258＝恒绿假钉"）**：**技术核〔已证〕**（258 那把尺逐字匹配 `call.Fun.(*ast.Ident)`，生产那枚是方法调用＝`*ast.SelectorExpr`），**但那枚钉的颜色写反了**：票面 `AC#1` ⓐ 要求计数 `== 1` ⇒ 恒 0 命中的尺会**永远红**，⛔ 永远绿；真正"恒绿"的形是 ⓑ/ⓒ 被写成"对找到的每枚调用检查 X"（空集通过＝票 309 范本注释点名的"⛔ 交一枚扫了 0 个文件然后报绿"）。⇒ 派单该钉的是"三件事共用同一枚计数、且 `==1` 单独成断言"；另外⛔ 非要新造 SelectorExpr 层——同包现成范本逐字 `		if sel, ok := d.Call.Fun.(*ast.SelectorExpr); ok && sel.Sel.Name == "stop" {`（`resident_ball_228_test.go`，正是票 294 现量点名的三枚 AST 先例之一；`git grep -n "ast.SelectorExpr" HEAD -- cmd/wisp`＝**41 命中／13 枚件**）。
5. **一条给编排者省时的台面更新**（⛔ 顶、但是他要的）：票 309「排程更新」按住 `309-r1` 的两条前提（`298-v1` 取验收色、`308-a1` 读 `internal/risk` 台面色）——本腿收尾现量 `git log --oneline 934b350a..HEAD`（⛔ 本腿的笔）＝ `7d163392` ＋ `dd638f98` ⇒ **两枚都已交回**，且 `git diff --stat 934b350a HEAD -- internal cmd` 为空、`git status --porcelain -- internal cmd` 复跑 **0 行** ⇒ **`309-r1` 的排面冲突已消**，可与 `296-r2` 同波（波 A）。⛔ 由本腿改票面那句，记为过期读数。
