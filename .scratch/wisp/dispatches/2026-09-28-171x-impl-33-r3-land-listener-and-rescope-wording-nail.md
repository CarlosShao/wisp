# 派单 33-r3（写腿）— 把封存的那枚入向听众落成产码，并按"扫能力"改写片 A 那枚词面负向尺

- 时间：2026-09-28 17:1x +08　编排者锚（**不等于你的锚**）：`d3e96f56`
- 上游：`33-r2` 的证据件 `docs/evidence/s1/33-inbound-listener-r2.md`＋封存码 `.scratch/wisp/probes/33/r2/panel_inbound.go.33b-src`／`panel_inbound_33_test.go.33b-src`（**零枚产码进过树**）· 台账裁定在 **`A386`**（务必先读这一节的最后三条）· 派单原文 `.scratch/wisp/dispatches/2026-09-28-164x-impl-33-r2-inbound-production-listener.md`
- 工单：票 33（`.scratch/wisp/issues/33-panel-host-c27.md`，`AC#9` 已追加、**未勾**）

## 你要做的两件事（缺一不可，顺序不能反）

**① 把封存码落成产码**：`cmd/wisp/panel_inbound.go`（`wisp panel-inbound -data <目录>`，stdin 每行一枚原始封套 → `(*ComposerDispatch).Handle`）＋判据件 `cmd/wisp/panel_inbound_33_test.go`。**照 `.33b-src` 那份落**，不要顺手重构；落完先跑一遍颜色再动第二件事。⚠ `-data` **必须仍是必填**（fail-closed，别为了让端到端"好看"去写宿主真实 `%APPDATA%`，那会撞 `dataroot_128_test.go:113` 那把尺和 `A371` 的旧账）。

**② 改写片 A 那枚词面负向尺（本程的硬骨头）**：`internal/panel/composer_dispatch_test.go:440-464` 现在是全仓 `WalkDir`＋**纯文本**匹配 `"ComposerDispatch"`、命中非测试文件即红、唯二豁免是它自己那枚文件，**且不以"接没接原生宿主"为条件**。它挡的就是①。
- **裁定＝改扫能力，不许删除、不许加豁免名单、不许注释掉**（那等于把门拆了）。
- **新尺要问的是**："这棵树有没有把**原生宿主／WebView2 消息通道**接上"——用**能力特征**判，例如 `go.mod`／`go.sum` 里出现 webview 依赖、生产代码里出现消息处理器／窗口宿主那一族符号（自己现量取真符号名，别照抄我这句）。
- **必须自带正控**（没有正控的负向尺＝没装门）：**人为造一枚假宿主接上 ⇒ 新尺要红**；正控跑法与读数进证据件。⚠ 正控别用 `go test -overlay`——**扫真实文件的尺看不见 overlay**（`A382`）；载体建在**仓外副本**，仓内零删除。
- **旧语义要保住**：那枚钉的报错原话是"别让这里的绿冒充一个能点的面板"。⇒ 新尺红/绿的含义**只说"原生宿主有没有接上"**，而**"面板能不能真点"这件事继续写在票面 `AC#9`**里：明写**到这一步面板仍然点不动，H2／H3／H10 未落**，不许因为产码落了就把 `AC#9` 报成"界面通了"。

## 写面（超出即判失败）

`cmd/wisp/**`（**除 `cmd/wisp/run.go`**，别人的串行队列）＋ `internal/panel/composer_dispatch_test.go` 的**指定那格**（`440-464` 那枚函数，**其余用例不许动**）＋票 33 票面（追加，不改原文）＋`docs/evidence/s1/33-inbound-listener-r3.md`＋`.scratch/wisp/probes/33/r3/**`。

## 禁区（一条破即判失败）

- `internal/panel/composer_dispatch.go` **一字不动**（起手 md5 自己取，终态必须相同）；`bridge.go:42-45` 四枚常量的顺序与命名不动；**不新增任何 `panel.*`／`config.*` 常量**（名册补齐属票 194 堆1，排你之后）。
- `internal/panel/tokens_fourway_test.go`／`l2_grant_boundary_test.go`／`frontend_hygiene_test.go` 一字不修；`frontend/**`／`design/**` **零写面**（读可以，写不行，也不许计入任何"零命中"宣称）。
- `go.mod`／`go.sum`／`thresholds.go`／golden／`allowlist.txt`／C18 超时常量／`docs/PLAN.md`／`docs/specs/**` 一字节不动（**不许为了接宿主引入 webview 依赖**——那是票 33 余下几格的射程，不是本程）。
- 别人的在飞／脏件不碰：`internal/risk/**`（含 `hostpath_185.go`）、`internal/tools/pointer_185_cli_seam_test.go`、`design/**` 那批未提交删除、`probes/152`、`probes/161/r6/logs/flip-*`。
- Git 纪律：只 commit **绝不 push**；每次 commit 带**显式 pathspec**；禁 `git add -A`／`.`；禁 `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`restore`／`clean`；**仓内零删除命令**（临时件只建不删；恢复跟踪件用 `git cat-file blob HEAD:<path> > <path>`）；**每完成一格 commit 一次**。

## 起手必做

1. `date`＋`git log --oneline -1`＋`git status --porcelain`：起手名册逐枚抄进证据件第一行。
2. 现量并记入表：`go test -count=1 ./internal/panel/`（起手应＝**恰三枚在册红**：`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourWayAgree`）；`grep -cE '=\s*"panel\.' internal/panel/bridge.go`（＝**4**）；`wc -l`＋md5 取 `composer_dispatch.go`；`grep -rn "ComposerDispatch" --include=*.go internal/ cmd/ | grep -v _test.go`（非 test 命中＝**0**）。
3. **未定义即停**：上面任何一条与派单不符（例如红不是那三枚、或封存码与描述不一致），**停手具名报回**，不要自创形状、不要放宽任何断言来变绿。

## 门禁与终态

`sh scripts/d22scan.sh`（只看 clean／红点逐名；⚠ `ban #8 internal/=N`、`cmd/=N` 那两格是**文件枚数**、不是违规数，本程新增产码会让 `cmd/` 涨，涨≠漂移）；`bash .scratch/wisp/probes/154/gate-clauses.sh`（只比 **BAD 腿名册**＝今天在册**只 `G6neg`**；**绝不跑 `probes/161/r6/flip-declaration.sh`**）；`go test -count=1 ./cmd/wisp/ ./internal/panel/`（后者必须**回到恰三枚在册红**，第四枚红＝你把尺改坏了）；`"$(go env GOPATH)/bin/gofumpt" -l cmd/ internal/panel/` 空。跑 CLI 相关测试带 `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`（否则 `0xc0000135`）。终态 `git status --porcelain` **等于起手名册**（逐枚具名差集为空）。

## 回报（≤350 字，六栏）

① 新尺的**判据谓词**是什么（扫哪几个真符号／哪几个文件）＋**正控那一发红在哪条用例**；② 落码前后 `./internal/panel/` 与 `./cmd/wisp/` 的颜色逐名；③ 变异／正控逐发"改哪一行 ⇒ 哪条红"；④ 门禁四数原样读数；⑤ 你被拒过的每次工具调用（一栏，别省）；⑥ **离"界面点一下后端真收到"还差哪几跳**（H2／H3／H10 那几枚，具名）。
票面：`AC#9` 写"产码已落＋新尺改扫能力＋**面板仍点不动**"，**别勾**（勾要非实现者表）。预算：硬顶 **30**，第 5 枚内落盘第一枚 commit；超预算要具名说被哪一步吃掉、**不许用超预算作为放宽断言的理由**。
