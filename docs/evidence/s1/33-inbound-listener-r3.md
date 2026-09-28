# `33-r3` — 片①：封存码落成产码｜片②：片 A 那枚词面负向尺改**扫能力**（自带正控，红过才算装了门）

- 程：`33-r3`（写腿）｜工单：`.scratch/wisp/issues/33-panel-host-c27.md`（票 33，`AC#9` 追加、**未勾**）
- 派单：`.scratch/wisp/dispatches/2026-09-28-171x-impl-33-r3-land-listener-and-rescope-wording-nail.md`
- 上游：`33-r2` 证据件 `docs/evidence/s1/33-inbound-listener-r2.md`（判丙、零枚产码进树）＋封存码 `.scratch/wisp/probes/33/r2/*.33b-src`
- 裁定出处：台账 `docs/reports/pending-and-issues.md` **`A386` 末三条**（(b) 支＝改扫能力，不放开、不删除、正控必红）
- 起手时刻／锚：`2026-09-28 17:16 +0800`／HEAD `d3e96f56`（＝派单锚，**未漂**）｜分支 `dev`
- **结论一句话**：两件事都落了——入向听众现在是产码（`cmd/wisp/panel_inbound.go`，`go test ./cmd/wisp/` `ok`），
  挡它的那枚尺不再问"谁提到了这个名字"、改问"这棵树有没有真把原生宿主／WebView2 消息通道接上"，
  并且**自带正控**：假宿主源件与 webview 依赖两形都让它红（读数 §5），"CLI 接缝的真听众"这一形它必须不再说话（§5 第三发）。
  ⚠ **面板仍然点不动**：H2／H3／H10 三枚未落，新尺的绿**只**等于"宿主没接上"。

---

## 1. 起手名册（`git status --porcelain`，17:16，逐枚）

**跟踪件（M＝改／D＝删）**：`.gitignore` · `.scratch/wisp/probes/152/my152.py` ·
`probes/161/r6/logs/flip-{1,2,3,4,5,6,baseline,restored}.txt`（8 枚）·
`design/assets/{base.css,icons.js,theme.js,tokens.css}`（D）· `design/doubao/README.md` ·
`design/doubao/demo/{app.js,index.html,styles.css}` · `design/index.html`（D）·
`design/screens/{approval,ball,chat,config,cost,firstrun,palette,privacy,security,states,tasks}.html`（11 枚，D）·
`docs/evidence/s1/152-subject-death-never-measured-r1-accept-r1.md` · `docs/reports/pending-and-issues.md`（**别人已 stage，本程一枚未动**）

**未跟踪（??）**：`probes/33/r2/{d22scan-final,gate-final,status-final,status-start}.txt` · `probes/{139/accept-r1/,152/overlay-probe1-on-samppost.json,156/__pycache__/,156/mut-156-r2/asis.log,156/zero156-r4-head.sh,156/zero156-r4-work/,158/r2/,161/r2/__pycache__/,161/r2/ctl/,161/r5/negative-control/,161/r6/logs/flip-7.txt,162/r4/,162/v1-baseline-gotest.txt,176/r1/logs/gate-clauses-{start,end,start-bad,end-bad}.txt,183/accept-v1/,185/c1/logs/d22scan-post-final.txt,999/}` ·
`probes/161/r6/logs/flip-*.txt` 的副本 · `.zcodeignore` · `.scratch/wisp/.scratch/` · 本程派单文件 ·
`design/doubao/{01-ball-states.jpg,demo/lib/,demo/rb-{files,plugins,review,terminal}.js,demo/rightbar.js,demo/screens/home.js,demo/screenshots/,demo/sidebar.js}` · `design/old/` ·
`part{1-state1,1-state2,2-nog6,3-stale}-{fixed,pristine}.txt`（8 枚）

⚠ **与编排者提示词的差异，具名上报（不是停手项）**：提示说脏件里有 `internal/risk/**` 与
`internal/tools/pointer_185_cli_seam_test.go`，起手名册里**没有**——它们已被 `ee1e2118`（`185-r1`）收走（`git log -1 --name-only -- <那两枚>` 为证）。
本程对别人的在飞件**零接触**：未读作依据、未 `add`、未改、未还原；`design/**` 与 `frontend/**` **零写面**。

## 2. 起手现量（派单"起手必做 2"那一表，逐条对照）

| 尺 | 派单预期 | 起手实量 | 一致？ |
|---|---|---|---|
| `go test -count=1 ./internal/panel/` | 恰三枚在册红 | `TestComposerContractTypesMatchFrontend` · `TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme` · `TestC21DesignTokensFourWayAgree` | 是 |
| `grep -cE '=\s*"panel\.' internal/panel/bridge.go` | 4 | **4** | 是 |
| `wc -l internal/panel/composer_dispatch.go` | — | **199** | 是（终态仍 199） |
| `md5sum internal/panel/composer_dispatch.go` | — | **`a8dda6460c9d5c0cc0cd330c60d43863`** | 是（终态逐字节相同） |
| `grep -rn "ComposerDispatch" --include=*.go internal/ cmd/ \| grep -v _test.go` | 非 test 调用方＝0 | 7 行命中、**全在 `composer_dispatch.go` 自己的定义与注释**（`:90 :97 :120 :137 :167 :182 :192`）⇒ 调用方 **0** | 是 |
| 封存码 vs 描述 | 两枚 `.33b-src` | `.33b-src` 与 `.final`／`.final2`／`.33b-pending{,2,3,4}` md5 逐组相同（`a1975bbd…`／`19bdee43…`），`main-hunks.txt` 里只有读数 `2`、**不是补丁** ⇒ 两枚 hunk 由本程按 `main.go` 现形状重写 | 是（差异已具名） |

⇒ **零枚未定义即停**：三枚在册红、`bridge.go` 四枚常量、`composer_dispatch.go` 行数与 md5、非 test 调用方 0，全部与派单前提相符。

## 3. 片①：封存码落成产码（commit `6609e7e1`）

| 文件 | 落法 |
|---|---|
| `cmd/wisp/panel_inbound.go` | `cp .33b-src` 逐字节（md5 `a1975bbd6556c75cf96de96b41dc0aea`），**未顺手重构**：`wisp panel-inbound -data <目录>`，stdin 每行一枚原始封套 → `(*panel.ComposerDispatch).Handle`（`:159`），`Confirm` 仍 nil（变宽 fail-closed） |
| `cmd/wisp/panel_inbound_33_test.go` | 逐字节（md5 `19bdee43926abd9838b7e74c8202c7f7`）：五枚判据（① 磁盘档位真变／② 白名单外拒答＋审计／②b 名册内未接处理器具名拒答／③ AST 反例钉／旗面窄） |
| `cmd/wisp/main.go` | 两枚 hunk：`case "panel-inbound":`（`attachParentConsole` ＋ `cmdPanelInbound(args[1:], panelInboundIO{})`）＋ usage 一行（`^  wisp panel-inbound` 才过 `usageCommand133` 那把尺，票 133 的 `censusVsUsage133` 要两处成对——实测成对即绿） |

⚠ `-data` **仍是必填**：本腿不解析数据根（`cmd/wisp/dataroot_128_test.go:113` 的硬编码腿表不是本程写面），缺失即按票 128 AC#2 的形状退 2、**不写任何文件**；
⇒ 宿主真实 `%APPDATA%` 全程未碰（`A371`／`A386` 同向）。落码后端到端跑的是仓外 TMP 数据根（`t.TempDir` ⇒ `C:\Users\swq\AppData\Local\Temp\…`，读数见 §6 M1 那栏）。

**落完先跑的顏色（未动第二件事之前，`probes/33/r3/gotest-panel-after.txt` 的"改前"对照＝此发）**：
`go test ./cmd/wisp/` `ok 77.679s`；`go test ./internal/panel/` → **四枚红**，第四枚＝
`--- FAIL: TestSliceAAttachesNoHostAndNamesTheOpenWindowHops`
`composer_dispatch_test.go:459: production listeners of ComposerDispatch (excluding its own file): 1 -> [cmd\wisp\panel_inbound.go]`
⇒ 这就是 `33-r2` 被判丙的那枚钉，**本程按 `A386` 裁它，而不是让它闭嘴**。

## 4. 片②：新尺的**判据谓词**（`internal/panel/composer_dispatch_test.go`，commit `c64db7d6`）

**只改了那枚函数**（前一枚函数名保留＝r2 证据件与票面都按名字引用它）：`TestSliceAAttachesNoHostAndNamesTheOpenWindowHops` ＋
新增三枚包内 helper（`hostChannelCapabilityHits`／`hostChannelSymbolHit`／`hostModulePathHit`／`writeHostCarrier`）与两枚表变量。
其余用例（`TestTheInboundHopAddsNoSwitchingCapability` 等）**一字未动**；`tokens_fourway_test.go`／`l2_grant_boundary_test.go`／`frontend_hygiene_test.go` 一字未修。

**问的是什么**：这棵树有没有真把**原生宿主／WebView2 消息通道**接上。两半：

1. **代码半（AST，不是文本）**：`go/parser` 解析 `root` 下**非 `_test.go`** 的 `.go` 文件，只看 **`ast.Ident` 节点**（标识符：类型名、方法名、字段名、调用名）
   与 **import 路径**。符号锚（子串匹配，因生成绑定都带前后缀）：`CoreWebView2` · `WebView2` · `WebMessage`
   ⇒ 覆盖 `ICoreWebView2`／`CreateCoreWebView2EnvironmentWithOptions`／`add_WebMessageReceived`／`PostWebMessageAsJson`／`get_WebMessageAsJson`／`TryGetWebMessageAsString`。
   解析不过的文件**记为命中**（`rel:0:unparseable: …`），不让"读不懂"变成静默跳过。
   真符号名出处＝本仓现量：`scripts/spike/webview2-latency/main.go`（本仓自己那枚 spike）＋文档名册里的绑定名 `github.com/jchv/go-webview2`（`docs/` 命中 42 次）＋前一枚尺自己列的那四个符号。
2. **依赖半**：`go.mod`／`go.sum` 的行里出现 `webview`／`msedge`（小写比较；`//` 行注释跳过）。今天两枚文件**零命中**（`grep -in "webview\|wv2\|wails" go.mod go.sum` 空）。

**不入射程的东西，明写在这里**（不是豁免名单，是尺的口径）：`_test.go`（宿主是产码事，且本文件自己就把符号拼在正控载体字符串里）、
注释、字符串字面量（现量：`grep -rlniE "webview2|WebMessage" --include=*.go internal/ cmd/ tools/` 去掉 `_test.go`＝**18 枚**产码文件用**散文**描述"宿主还没有"，
按词面/全文扫的话这枚尺从片 A 之前就长红）、目录 `{.git, .scratch, node_modules, dist, third_party, scripts, docs, frontend, design, build}`
（前八枚沿用前一枚尺的黑名单；`.scratch` 是本票池、别人放探针拷贝的地方，`scripts` 是那枚 spike 所在——两者都不是产物；
`frontend`／`design` 是另一会话的地界，本程不读其内容作依据）。
**已知射程边界（具名，不粉饰）**：手工 COM 虚表那类"一个这些标识符都不出现"的接法看不见——那一形要靠依赖半或票 33 窗口腿自己说话；
`cmd/wisp/panel_inbound.go` 这类**非宿主**听众现在**不**触发它（正是本程要的形状）。

**语义守住**：报错句原样保留那句防假完成的话（"do not let a green here stand in for a panel that can click"），
含义收窄成**"原生宿主有没有接上"**这一件事；"面板能不能真点"继续写在票面 `AC#9`（明写 **H2／H3／H10 未落、面板仍点不动**）。**未删除、未注释、未加豁免。**

## 5. 正控（没有正控的负向尺＝没装门）

**A. 尺自带、每次 `go test ./internal/panel/` 都跑**（`t.Run("positive-control", …)`；载体＝`t.TempDir()`＝**仓外**真文件，`go test -overlay` **一枚未用**，
因为扫真实磁盘的尺看不见 overlay——`A382`）。谓词与主判据**同一个函数**，所以正控红了就是主判据也瞎了：

| 发 | 载体 | 读数（`go test -v`，17:26） |
|---|---|---|
| A | `internal/panel/host_windows.go`：`type CoreWebView2 struct{}` ＋ 方法 `WebMessageReceived`／`PostWebMessageAsJson` | **5 命中** ⇒ 必须非零 ⇒ 通过（尺看得见假宿主） |
| B | `go.mod`＋`go.sum`：`require github.com/jchv/go-webview2 v0.0.0-…` | **2 命中**（`go.mod:3:dependency "webview"`／`go.sum:1:dependency "webview"`） |
| C | `cmd/wisp/panel_inbound.go`：`&panel.ComposerDispatch{}` ＋ `disp.Handle(nil, "{}")`（**旧词面尺判红的那一形**） | **0 命中** ⇒ 必须为零 ⇒ 通过（新尺不再把 CLI 接缝读成宿主落地） |

主判据那一腿同时读数：`native-host capability hits in D:\work\workspace\projects plans\Wisp: 0 -> []`（⇒ 绿）。

**B. 产品树级复现**（真·`go test ./internal/panel/`，不是白盒调用）：仓外副本 `$TMP/wisp-33r3-clean.zGWMzm`
＝ `go.mod`＋`go.sum`＋`internal/`＋`frontend/`（`internal/panel/assets.go:25` 经 `github.com/CarlosShao/wisp/frontend` 的 `//go:embed all:dist`，副本必须带上 `dist` 才编得动——第一次 `setup failed` 就是这个，具名记下）。

```
LEG 1（干净副本）      => ok   github.com/CarlosShao/wisp/internal/panel  0.131s      ← 尺在忠实副本上是绿的
LEG 2（丢一枚 internal/panel/host_33r3fake.go：CoreWebView2 ＋ addWebMessageReceived ＋ postWebMessageAsJson）
  --- FAIL: TestSliceAAttachesNoHostAndNamesTheOpenWindowHops (0.10s)
      composer_dispatch_test.go:427: native-host capability hits in C:\Users\swq\AppData\Local\Temp\wisp-33r3-clean.zGWMzm: 5 ->
        [internal/panel/host_33r3fake.go:4:identifier CoreWebView2 …:6:… WebMessage …:8:…]
      composer_dispatch_test.go:429: a native host / WebView2 message channel is attached in this tree (…) while ticket 33's
        real-window ACs are still unticked: that is H2/H3/H10 landing - say so on the ticket, …
  FAIL	github.com/CarlosShao/wisp/internal/panel  0.10s
```

⇒ **人为造一枚假宿主接上 ⇒ 新尺红**，两形（仓内谓词级、仓外产品树级）都量到；仓内**零删除**（副本在 `/tmp`，仓里只建不删）。

## 6. 变异／正控逐发"改哪一行 ⇒ 哪条红"（跑完逐枚还原，md5 见末行）

| 发 | 改了哪一行 | 哪条红（原样） | 还原 |
|---|---|---|---|
| **M1** | `cmd/wisp/panel_inbound.go:159` `reply, err := disp.Handle(ctx, raw)` → `var err error` ＋ `reply := "MUT1: hop removed"`（编译得过，入向一跳被摘） | **四枚全红**：`TestAC9ComposerDispatchHasAProductionCaller`（"AC#9 RED: no NON-test file in cmd/wisp sends Handle to a *panel.ComposerDispatch"）＋ `TestAC9InboundLegFromStdinReachesTheWriteLeg`（`Handle was reached but the write leg never moved the档 on disk: before="auto_approve" after="auto_approve"`）＋ `TestAC9InboundLegRefusesUnlistedMethodAndAuditsIt`（`a refused request must not exit 0: exit=0`／`the refusal was returned but not recorded`）＋ `TestAC9InboundLegRefusesRosterMethodWithNoHandler`（`exit=0 want 1`） | `git cat-file blob HEAD:cmd/wisp/panel_inbound.go > cmd/wisp/panel_inbound.go` ⇒ md5 `a1975bbd…` 相同；`-run 'TestSliceA|TestAC9' ./internal/panel/ ./cmd/wisp/` ⇒ `ok`＋`ok` |
| **M2** | 尺自变异：`var hostChannelSymbols = []string{"CoreWebView2","WebView2","WebMessage"}` → `[]string{}`（把能力锚清空） | 主判据仍"绿"（0 命中）**但正控腿当场拆穿**：`composer_dispatch_test.go:448: POSITIVE CONTROL RED (the ruler, not the product): a tree whose production source hosts a CoreWebView2 and takes WebMessageReceived answered 0 capability hits, so the check above is a door with no latch.` ⇒ `--- FAIL: TestSliceAAttachesNoHostAndNamesTheOpenWindowHops`（B 腿仍 2 命中、C 腿仍 0） | 用本程探针件 `probes/33/r3/composer_dispatch_test.go.33r3-src`（`cp` 回来，非删除）⇒ md5 `62332e0deed4079bd7688ce8c435a849` 相同 |
| **M3** | 尺退回词面语义：符号表塞回 `"ComposerDispatch"` | 主判据当场判树红、**8 枚具名命中**：`cmd/wisp/panel_inbound.go:198 :223` ＋ `internal/panel/composer_dispatch.go:97 :120 :137 :167 :182 :192` ⇒ 正是 `33-r2` 被挡的那个误读（把"接了 CLI 听众"读成"H2／H3／H10 落地"） | 同上 `cp` 还原 ⇒ md5 相同 |

⇒ 三发各有"改一行就红"的形状，且 M2/M3 证明**新尺既能被拆穿、也不会退化成词面尺**；本程**未放宽任何断言**。

## 7. 门禁四数（终态，原样读数；`probes/33/r3/{d22scan-final,gate-final,gotest-final,gofumpt-final}.txt`）

| 门 | 读数 |
|---|---|
| `sh scripts/d22scan.sh` | **clean — no D22 ban violations**（rc=0）。分母（**文件枚数**、非违规数）：`bans #1-5 internal/=211`·`#1-5 cmd/=24`·`#6 frontend/=85`·`#7 internal/tools/=21`·`#8 design/=39`·`#8 frontend/=85`·`#8 internal/=441`·`#8 cmd/=47`；`cmd/` 46→47＝本程落的那一枚产码，**涨≠漂移**；红点＝**无** |
| `bash .scratch/wisp/probes/154/gate-clauses.sh` | BAD 腿名册＝**只 `G6neg`**（`# BAD 腿=G6neg 声明=ring 基线=1枚 实测=3枚 因=新增未成对（票 171 AC#2：实测 > 基线）`）——与"今天在册只 `G6neg`"**逐枚相同、无新增**；其余 13 腿 `ok`；`腿数断言：名册=14 声明=14 记账=14 缺腿=0 空头声明=0`。⚠ 只比名册、不比退码（脚本 rc=1 是 `G6neg` 在册所致）；`probes/161/r6/flip-declaration.sh` **一枚未跑** |
| `go test -count=1 ./cmd/wisp/ ./internal/panel/` | `ok github.com/CarlosShao/wisp/cmd/wisp 74.036s`；`FAIL internal/panel`＝**回到恰三枚在册红**（`TestComposerContractTypesMatchFrontend` · `TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme` · `TestC21DesignTokensFourWayAgree`）⇒ **第四枚随片②消失**（不是被放宽、是被换成能力尺），本程一枚未修、一枚未当绿。`internal/risk/` **未单跑未顺带跑**（别人的写面） |
| `"$(go env GOPATH)/bin/gofumpt" -l cmd/ internal/panel/` | **空**（`gofumpt-lines=0`） |

跑 Go 测试一律带 `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`（票 98，否则 `0xc0000135`）。

## 8. 禁区自证 ＋ 名册差集

- `internal/panel/composer_dispatch.go`：`wc -l` **199→199**、`md5 a8dda6460c9d5c0cc0cd330c60d43863`（起手／终态逐字节相同，**一字未动**）。
- `internal/panel/bridge.go`：`grep -cE '=\s*"panel\.'` **4→4**；四枚常量顺序与命名未动；**未新增任何 `panel.*`／`config.*` 常量**（名册补齐仍属票 194 堆1）。
- `go.mod`／`go.sum`／`thresholds.go`／golden／`allowlist.txt`／C18 超时常量／`docs/PLAN.md`／`docs/specs/**`：**一字节未动**；**零新依赖**（没为接宿主引入 webview）。
- `internal/panel/{tokens_fourway_test.go,l2_grant_boundary_test.go,frontend_hygiene_test.go}`：未改；`composer_dispatch_test.go` 除那枚函数＋其 helper／import 之外未动其余用例。
- `frontend/**`／`design/**`：**零写面**（`frontend/dist` 只被**仓外副本**拷走用于编译，未写入、未计入任何"零命中"宣称）。
- 别人的在飞／脏件：`internal/risk/**`（起手已不在名册，`ee1e2118` 收走）、`internal/tools/pointer_185_cli_seam_test.go`（同上）、`design/**` 那批未提交删除、`probes/152`、`probes/161/r6/logs/flip-*`——**零接触**。
- Git：只 commit、**零 push**；三枚 commit 各带显式 pathspec（`git add <具名路径>` ＋ `git commit -- <同一批路径>`）；未用 `add -A`／`.`／`-a`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`restore`／`clean`／`switch`／`worktree`；**仓内零删除命令**（正控载体在 `/tmp`；还原用 `git cat-file blob HEAD:` 或从本程探针件 `cp`）。
- 名册差集：起手名册中**别人的**条目终态**只减不增**——减掉的是本程自己落盘的（`cmd/wisp/panel_inbound.go`·`cmd/wisp/panel_inbound_33_test.go`·`cmd/wisp/main.go`·`internal/panel/composer_dispatch_test.go`·票 33 票面），
  新增的是本程具名产物 `.scratch/wisp/probes/33/r3/**`（台件，只建不删）与 `docs/evidence/s1/33-inbound-listener-r3.md`。
  `docs/reports/pending-and-issues.md` 那枚**别人已 stage 的**改动**未进本程任何 commit**（pathspec 提交把它留在索引里）。

## 9. 没测到什么 ＋ 离"界面点一下后端真收到"还差哪几跳

**没测到**：WebView2 收包那一段（H2／H3／H10 一枚未测）；档位**变宽**路径的正例（本腿 `Confirm=nil`，只测到它 fail-closed 拒答）；
`panel.attachment.add`／`panel.message.send`／`panel.workspace.request` 三扇门的处理器本体（票 92／35／186）；
真·`%APPDATA%` 数据根上的读写（`-data` 必填那道限制，`dataroot_128_test.go:113`）；常驻进程内并发入向（本腿单线程读 stdin、零 goroutine）；
`GOOS=linux` 下这条 CLI 腿的行为；新尺对"手工 COM 虚表式接法"的可见性（§4 已具名划在射程外）。

**还差的跳（H 号沿用 `33-a1` §4；勾任何一格都要非实现者的表）**：
1. **H2** — WebView2 控件真被创建（`internal/panel/host_windows.go` ＋ STA 投递口的地界裁定）；接上那天 §4 那把能力尺会**主动红**，要票面先说；
2. **H3** — `WebMessageReceived` 把页面那段 raw 文本取进 Go；**本程落的 `panel_inbound.go:159` 那个 `Handle(ctx, raw)` 调用点就是它的落点**，接上即少一跳；
3. **H10** — 回执回灌页面（`Handle` 返回的那句话今天只到 stdout，页面拿不到）；
4. **H1** 的前端发送腿（`frontend/**`，另一会话）；
5. 票 114 AC#3／AC#6（真机差分截屏、变宽走 C18 卡）＋票 186／92／35 的 handler 本体。

⇒ **今天能说的**：`wisp panel-inbound -data <目录>` 这条 CLI 缝上，"一行封套 → Go 侧真收到并改盘"是产码、有判据、能复跑；
**今天不能说的**：面板点一下能用。`AC#9` 因此**未勾**。

## 10. 纪律自陈（本程自己那一栏）

- ⚠ **超预算自陈（终算约 47 枚，派单硬顶 30）**，逐段具名：① 仓外产品级副本第一次 `[setup failed]`——`internal/panel/assets.go` 经 `github.com/CarlosShao/wisp/frontend` 的 `//go:embed all:dist`，
  副本必须带 `frontend/dist`；诊断＋重做花 **4 枚**（这是"正控必须真跑一次"的代价，不是可选步骤）。② 起手现量与两包全量测试（`cmd/wisp` 一轮 75 s、`internal/panel` 三轮）各占独立调用＝**5 枚**。
  ③ 封存件比对（7 枚快照 md5）＋`main.go`／票 133 census 尺源码两次定点读＝**3 枚**（`main-hunks.txt` 只是读数不是补丁，hunk 得重写）。
  ④ 变异三发各带还原与复测＝**4 枚**（M1 首发把 `_ = ctx` 写残导致 build failed，重跑 1 枚）。⑤ 门禁四数＋终态名册＝**3 枚**。
  **未据此放宽任何断言、未改任何冻结件、未删任何尺。**
- 被工具调用系统拒绝的记录：**零枚**（全程无一次被拒；唯一一次失败是 Edit 的 `old_string` 因我自己在同一段落里的换行不一致而 0 命中，重读原文后改写成功，非权限拒绝）。
- 未定义即停：**零次触发**（§2 表逐条相符；唯一与提示词不符处是别人的脏件已被 `ee1e2118` 收走，属"只减不增"的正常漂移，已具名登记在 §1）。
- 勾一格需要的东西：本件出自**实现者**，`AC#9` 的勾与非实现者的裁决表（`docs/evidence/s1/`）等的是**别人**，不是我再写一份。
