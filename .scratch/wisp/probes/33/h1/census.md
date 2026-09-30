# 33-h1 — 面板宿主「窗口真开出来」落地普查（只读，零产码，零 go.mod 改动）

**腿号**：`33-h1`　**派单**：编排者 2026-09-30 23:1x 收件后的只读普查腿　**票面**：`.scratch/wisp/issues/33-panel-host-c27.md`
**开工钟点**：`2026-09-30 23:16 +0800`　**分支**：`dev`　**起手 HEAD**：`77a3b442`（与取数同发，见 §8 尺 1）
**写点**：本文件一处（`.scratch/wisp/probes/33/h1/census.md`）。票面／产码／测试／配置／冻结件**零写面**。

---

## ① 依赖：`jchv/go-webview2` 的今日现状与本机可拉性

（填写中）

## ② 宿主落点：新包 vs 装配根文件，两形各列出新增的包级依赖边

（填写中）

## ③ STA／消息循环共存：球的 `ui-sta` 循环与 WebView2 宿主在同进程会不会互抢套间

（填写中）

## ④ embed 供给：`//go:embed assets/web/dist` 的路径存在性与「embed 空目录」的行为

（填写中）

## ⑤ 「无监听端口」这条判据怎么机读（能力型，不是词面型）＋ CI 能不能真跑到

（填写中）

## ⑥ 延迟两数（冷 ≤1500ms／热 ≤200ms）的测量方法普查

（填写中）

## ⑦ 先后与互斥：与本仓在写腿的包级冲突表与建议顺序

（填写中）

## ⑧ 我跑了哪些尺，每条真实读数（原样命令）

> 所有命令均在 `D:/work/workspace/projects plans/Wisp`（分支 `dev`）由 Git Bash 跑；
> `M` = `D:/work/base/gopath/pkg/mod/github.com/jchv/go-webview2@v0.0.0-20260205173254-56598839c808`。
> **禁跑闸门遵守情况：`go build`／`go vet`／`go test`／任何带 `./...` 的命令＝零枚。** `go list` 只指名单包。

| # | 原样命令 | 真实读数 | 钟点 |
|---|---|---|---|
| S1 | `date "+%Y-%m-%d %H:%M %z"` | `2026-09-30 23:16 +0800` | 23:16 |
| S2 | `git rev-parse --abbrev-ref HEAD` / `git log --oneline -3` | `dev` ／ `77a3b442`·`bb37fac2`·`457a824e`（起手 HEAD） | 23:17 |
| S3 | `git status --short \| head -40` | 脏项：`M .gitignore`、`M .scratch/wisp/issues/33-panel-host-c27.md`、`M cmd/wisp/{approval_reply,main,resident_sink_nail_127_windows_test,resident_windows,run}.go`、`D design/**`（16 枚）、`M docs/evidence/s1/152-*.md`／`246-*.md`、`?? .scratch/.scratch/`。**全部未动、未提交** | 23:17 |
| S4 | `ls`（仓根） | 有 `go.mod go.sum internal cmd scripts tools docs frontend third_party deps.toml`；**无 `vendor/`、无 `assets/`** | 23:17 |
| S5 | `cat go.mod` | 6 枚直接依赖（sherpa-onnx-go / go-toml v2 / x\*crypto / x\*sys v0.48.0 / x\*term / modernc.org/sqlite）＋10 枚 indirect；**`jchv` 零命中** | 23:17 |
| S6 | `ls -d vendor` | `No such file or directory`（退码 2）⇒ 本仓**不做 vendor** | 23:17 |
| S7 | `go env GOMODCACHE GOFLAGS GOPROXY GOSUMDB` | `D:\work\base\gopath\pkg\mod` ／ 空 ／ **`https://goproxy.cn,direct`** ／ `sum.golang.org` | 23:17 |
| S8 | `ls $GOMODCACHE/github.com/jchv` | `go-webview2@v0.0.0-20260205173254-56598839c808` ＋ `go-winloader@v0.0.0-20250406163304-c1995be93bd1`（**两枚都在本机模块缓存里**） | 23:17 |
| S9 | `ls $GOMODCACHE/cache/download/github.com/jchv/go-webview2/@v/`（同构跑 go-winloader） | 各含 `.info .mod .zip .ziphash .lock list` ⇒ **离线可完整解析，不需网络** | 23:20 |
| S10 | `find $M -type f` | 库共 **40 枚 `.go`** ＋ 3 枚 `.dll`（`webviewloader/sdk/{x64,x86,arm64}/WebView2Loader.dll`）＋ `LICENSE` ＋ `webviewloader/sdk/LICENSE.txt`（3-Clause BSD） | 23:21 |
| S11 | `grep -rln '"C"' $M --include=*.go` | **退码 1＝零命中**（无 cgo） | 23:17 |
| S12 | `grep -rln "go:cgo_import_dynamic\|#cgo" ...jchv` | **零命中**；唯一 C 源件＝`go-winloader@…/tinydll/tinydll.c`（不参编，被嵌成字节） | 23:22 |
| S13 | `grep -rn "^//go:build" $M --include=*.go` | `webview.go:1`／`pkg/edge/corewebview2.go:1`／`chromium*.go:1` 全 `windows`；`internal/w32` 按架构分 | 23:17 |
| S14 | `cat $M/go.mod` | `go 1.16`；require `go-winloader v0.0.0-20250406163304…` ＋ `golang.org/x/sys v0.0.0-20210218145245-beda7e5e158e` | 23:18 |
| S15 | `cat go-winloader@…/go.mod` | `go 1.14`；require `golang.org/x/sys v0.0.0-20200810151505-1b9f1253b3ed` ⇒ 两枚库要求的 x/sys 都**低于**本机 `v0.48.0`，MVS 不会抬高根 go.mod 的 x/sys | 23:22 |
| S16 | `go list -deps github.com/jchv/go-webview2/pkg/edge`（在 `scripts/spike/` 内，单包名） | 总 **89** 枚包；外部名 **10** 枚＝`golang.org/x/sys/windows`（已有）＋ **9 枚全新增**：`go-webview2/pkg/edge`·`go-webview2/internal/w32`·`go-webview2/webviewloader`·`go-winloader`·`go-winloader/internal/{loader,pe,vmem,memloader,winloader}` | 23:20 |
| S17 | `go list -deps ./internal/panel` / `./internal/ball`（单包名） | panel **151**（其中 `github.com/CarlosShao/wisp/*` **8**：winsec·secret·observe·risk·projctx·streamkey·frontend·panel）；ball **145**（wisp 内部 **5**：winsec·secret·observe·statemachine·ball）；两包外部非 wisp 依赖只有 `go-toml/v2` 家族 | 23:20 |
| S18 | `grep -c "jchv" go.sum` / `grep -n "jchv" scripts/spike/go.sum` | 根 `go.sum`＝**0**；spike `go.sum:13-16`＝**4 行**（两枚 module 各带 `h1:` 与 `/go.mod h1:`） | 23:20 |
| S19 | `cd scripts/spike && go list -m all` | 列出 `go-webview2 v0.0.0-20260205173254-56598839c808` ＋ `go-winloader v0.0.0-20250406163304-c1995be93bd1`，**未联网即成**（缓存命中） | 23:20 |
| S20 | `grep -rn "^func" $M/pkg/edge/chromium.go` | `NewChromium:47`·**`Embed:72`**·`Navigate:116`·`Show:151`·`Hide:155`·`EnvironmentCompleted:171`·`CreateCoreWebView2ControllerCompleted:186`·`MessageReceived:233`·`WebResourceRequested:281`·**`AddWebResourceRequestedFilter:292`**·`Environment:299`·`GetSettings:324`·`Focus:355` | 23:18 |
| S21 | `sed -n '72,114p' $M/pkg/edge/chromium.go` | `Embed` 内 **`:96-111` 自带一枚 `for { GetMessageW / TranslateMessage / DispatchMessageW }` 嵌套泵**，出口条件只有 `atomic.LoadUintptr(&e.inited)!=0`（`:97`）或 `GetMessageW` 返回 0（`:106`）；`inited` 在 `CreateCoreWebView2ControllerCompleted:224` 才置 | 23:18 |
| S22 | `grep -n "log.Fatal" $M/pkg/edge/*.go` | **`chromium.go:173` `log.Fatalf("Creating environment failed with %08x")`**、**`:188` `log.Fatalf("Creating controller failed with %08x")`**、`:284`（WebResourceRequested 取 request 失败）、`webview.go:141`（Eval 里 `log.Fatal`）；同步失败路 `:89/:92` 只是 `log.Printf`＋`return false` | 23:19 |
| S23 | `grep -n "func newI\|iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler" $M/pkg/edge/` | 该 handler 类型与构造函数**全小写未导出**（`corewebview2.go:238`、`chromium.go:23/:60`）；`ls $M/pkg/edge \| grep -i Environment`＝**零枚文件** | 23:22 |
| S24 | `grep -n "^func\|^type" $M/webview.go` | `New:87`·`NewWindow:92`·**`NewWithOptions:97`**·`CreateWithOptions:269`（自己 `RegisterClassExW` 类名 `"webview"`、自己 `CreateWindowExW`）·`Destroy:347`·**`Run:351`（又一枚 `GetMessageW` 死循环）**·`Terminate:381`·`Dispatch:443`（`PostThreadMessageW(w.mainthread,…)`）·`Bind:450` | 23:18 |
| S25 | `cat internal/panel/assets.go` | 现成供给缝：`BuiltinAssets:38` 走 `fs.Sub(frontend.Dist(), "dist")`；`EntryFile:31`＝`index.html`；`errNotBuilt:35`；`Built:57`；**`Resolve:78`**（返回 bytes＋Content-Type）；`Manifest:113`；`Check:158`；`contentTypeOf:196` | 23:18 |
| S26 | `ls -d frontend/dist` / `find frontend/dist -type f -printf "%s\t%p\n"` / `du -sb frontend/dist` | **存在**，4 枚文件：`.gitkeep` 0 B／`assets/index-BRKj5OIJ.css` 49943 B／`assets/index-BVKlegVD.js` 553469 B／`index.html` 1044 B；合计 **604456 B**。**只取名字与字节数，未读任何内容**（`frontend/**` 两层禁令） | 23:21 |
| S27 | `ls -d assets assets/web assets/web/dist` | **三枚全部 `No such file or directory`（退码 2）** | 23:18 |
| S28 | `ls -d internal/panel/assets internal/panel/assets/web` | 同样**不存在**（退码 2）⇒ 票面那条 `//go:embed assets/web/dist` 无论按仓根还是按包目录解释，今天都没有落点 | 23:22 |
| S29 | `sed -n '14,30p' .gitignore` | **`.gitignore:24` 逐字＝`assets/web/dist/`**；而 `:22-23` 是 `frontend/dist/*` ＋ `!frontend/dist/.gitkeep`（注释 `:18-21` 具名写"go:embed needs the directory to exist in a clean checkout (ticket 77 AC#1)"）⇒ **那枚路径被整条 ignore、连锚文件豁免都没有** | 23:23 |
| S30 | `grep -rn "all:dist" docs scripts internal cmd tools .github` | `docs/evidence/s1/156-…:762` 逐字 `frontend/embed.go:19 → //go:embed all:dist`；`docs/evidence/s1/85-preflight-staticcheck.md:256` 同；`docs/reports/frontend-handoff.md:24`／`frontend-handover-to-zcode.md:34` ⇒ **现用 embed 是 `frontend/embed.go` 的 `all:dist`，不是票面那条**（二手引仓内证据件，未读 frontend 一字） | 23:23 |
| S31 | `Read internal/ball/sta_windows.go`（163 行逐行） | `runtime.LockOSThread():52`／`CoInitializeEx(0, coinitApartmentThreaded):62`／泵 `:77-93`（`GetMessageW`→`TranslateMessage`→`DispatchMessageW`）／`releaseCOM→CoUninitialize:105`／**`PostTask:127`：`hwnd==0` 时 `:134-141` 只 `delete(tasks,id)` 后 `return`**／`runTask:146`／`quit():156-162`（`PostQuitMessage`） | 23:17 |
| S32 | `Read internal/ball/ball_windows.go`（1004 行逐行） | `New:143`→`Registry.Spawn("ui-sta","ball",nil,fn):170`→`b.sta(start(createOnSTA)):171`；`createOnSTA:182`（`ensureFactories:183`→`registerBallClass:187`→`CreateWindowExW:208`→`s.hwnd=…:219`）；**`uiRun:741-752`＝PostTask＋`<-done` 同步等**；`Close:935-968`（PostTask 内 DestroyWindow＋`sta.quit()`，随后 `<-b.sta.handle.Done()` join）；`fire:728` 契约注释 `:724-727`"callbacks must be quick and non-blocking" | 23:17 |
| S33 | `Read scripts/spike/webview2-latency/main.go`（328 行） | 用的是**高层** `webview2.NewWithOptions:149`（自建窗口），不是 `pkg/edge`；`pumpOnce:83-101`＝**手工 `PeekMessageW` 泵**；`:113-115` 具名坑："`w.Dispatch` 投递的是线程消息，只有 go-webview2 自己的 `Run()` 循环处理"；`runDriver:259` spawn 12 子进程测真冷 | 23:19 |
| S34 | `grep -n "webview\|P95\|cold\|hot" docs/evidence/s0/02-spike-report.md` ＋ `sed -n '144,152p'` | §3.4 两 run：cold P50/P95/max＝`879.7/1041.6/1116.8` 与 `1125.7/1256.4/(n=12)`；**全进程内首次 create＝1808.9 ms（run1，紧跟 12 子进程后、竞争态）／1163.7 ms（run2）**；hot show P50/P95＝`25.7/49.5` 与 `71.4/79.5`；hot 浏览器往返 p50＝3.5 ms；recreate P50/P95＝`955.7/1221.3` 与 `859.0/955.5`。§2.5 时序定义逐字见 `:82-89` | 23:19 |
| S35 | `sed -n '27,100p' internal/observe/goroutine.go` | `ResidentBaseline = 6`（`:39`）／`ResidentNames`（`:44-46`）＝ui-sta·audio-capture·hotkey-listener·db-writer·watchdog·log-flusher／**`TemporaryNames`（`:57-59`）里逐字含 `"panel-host"`（`:58`）**／`Spawn:262`，未入册名 `slog.Warn("goroutine outside the D38 roster (leak symptom)"):271` | 23:20 |
| S36 | `grep -rn "panel-host" internal cmd docs tools scripts` | 除 `goroutine.go:58`／`goroutine_test.go:255`（`CategoryTemporary`）外，产码**零调用方**；规格侧同名条目＝`docs/PLAN.md:2835`、`docs/specs/SPEC-01-architecture.md:122`（后者注明"**均在对应 DisposalScope 内**"） | 23:22 |
| S37 | `ls .github/workflows` ＋ `grep -n "runs-on:\|- name:" .github/workflows/ci.yml` | 两枚 workflow（`ci.yml`／`slo-fresh.yml`）；`ci.yml` 的 `runs-on`＝`ubuntu-latest`(`:66` lint、`:278` portable-core、`:665` frontend) ／ **`windows-latest`（`:388` test-windows、`:533` build＋SLO smoke）** ／ **`[self-hosted, wisp-slo]`（`:591` SLO full）** | 23:20 |
| S38 | `sed -n '120,200p' scripts/portable-tests.sh` | `core_pin`（`:125-148`）**含 `github.com/CarlosShao/wisp/internal/panel`（`:141`）**；`win_pin`（`:149-158`）＝proc·secret·config·perm·plugin·ball·risk·llmrecord ⇒ **`internal/panel` 不在 windows scope**；`core` 的 glob 里 `./internal/panel/...`（`:179`）、`windows` 的 glob 无 panel（`:185-189`） | 23:20 |
| S39 | `grep -rn "netstat\|GetExtendedTcpTable\|net.Listen" --include=*.go internal cmd tools scripts` | 非 test 命中 **2 处**：`internal/proc/treemetrics_windows.go:47`（`procGetExtendedTcpTable = modiphlpapi.NewProc("GetExtendedTcpTable")`）＋`tools/mockllm/server.go:243`（mock LLM 自己 listen，不在产码路径）。**全仓零 `netstat` 字样** | 23:20 |
| S40 | `grep -n "^func \|^type " internal/proc/treemetrics_windows.go` ＋ `sed -n '255,300p'` | 能跑的那把是**未导出**的 `treeTCPConnections(inTree map[uint32]bool) int`（`:269`）：`TCP_TABLE_OWNER_PIDAll`、AF_INET、**只数行数、不读 `state` 字段**（偏移只取 `tcpOwningPidOffset`），失败重试耗尽后 `return 0` | 23:20 |
| S41 | `sed -n '42,45p' internal/panel/bridge.go` ＋ `sed -n '1,60p' tools/d22scan/main.go` | 白名单四枚逐字（`:42-45`）＝mode/workspace/attachment/message，**零枚 config**（与票面 `:16` 一致）；d22scan 八禁：`1 bare-goroutine`（**每一枚 `go <anything>`**，唯一文件豁免＝`internal/observe/goroutine.go`）、`2 pathresolver-bypass`、`3 plaintext-key`、`4 wallclock-timeout`、`5 mirror-hash`、`6 panel-approval(frontend/)`、`7 internal-artifact-tool`、`8 emoji`（`internal/`＋`cmd/` 含注释与 `_test.go`） | 23:21 |
| S42 | `grep -rn "ball.New(" --include=*.go cmd internal` ＋ `grep -n "OnTrayPanel\|OnPanelHotkey" cmd/wisp/resident_ball_windows.go` | 常驻进程的球在 **`cmd/wisp/resident_ball_windows.go:114`** 建；`:127` `OnPanelHotkey` 与 `:128` `OnTrayPanel` **今天都只 `recordBallGesture(...)`**（后者串名 `"tray-open-panel"`），无宿主；另一枚 `cmd/balldebug/main.go:199` 逐字打印 `"tray: open panel (stub, ticket 33)"` | 23:21 |
| S43 | `grep -n "Width int\|Height int\|Scale float64\|KeepAliveInSession" internal/config/schema.go` | `[panel]`＝`PanelSection` 起 `:528`；`Width:532`（`default:"640"`）／`Height:534`（无默认＝0 auto）／`KeepAliveInSession:536`（`default:"true"`）／`Scale:538`（无默认＝0 跟随系统 DPI）；节注释 `:528` 具名 **hot-tier** | 23:21 |
| S44 | `sed -n '1,40p' deps.toml` ＋ `grep -c "license =" deps.toml` | `deps.toml` 头注逐字＝"pinned **native** dependencies (SPEC-11 §2.1)"，消费者＝`scripts/fetch-deps.ps1` ＋ `wisp doctor`；全文 `license =` **3 处**（都是 native 条目）。`ls third_party`＝`model-fetch·sherpa-onnx·spike-models`。**未找到任何"Go 模块许可登记册"**（尺＝`grep -rln "LICENSES\|licenses\|NOTICE" scripts tools docs/specs`＝零命中） | 23:21 |
| S45 | `sed -n '140,180p' docs/specs/SPEC-08*.md` | §5.1（`:143-154`）逐条：单例 PanelManager／**"jchv/go-webview2（MIT，纯 Go 无 cgo）"**／`AddWebResourceRequestedFilter` 从 `embed.FS` 喂、**不起 localhost**／前端无状态、每次 show 推 `panel.resync`／Runtime 缺失→无面板模式＋L2 原生降级卡＋**不得自动下载安装器**／多任务按 correlationId 分区。**五条里没有一条写"消息接收"**（与 `33-a1` §结论② 一致，本轮复核成立） | 23:21 |
| S46 | `ls "C:/Program Files (x86)/Microsoft/EdgeWebView/Application"` ＋ `find … -maxdepth 2 -name msedgewebview2.exe` | `153.0.4234.48`／`154.0.4258.37`／`SetupMetrics`；两枚版本目录**各含 `msedgewebview2.exe`**（路径逐字 `…/153.0.4234.48/msedgewebview2.exe`、`…/154.0.4258.37/msedgewebview2.exe`）⇒ **编排者 09-30 那枚读数我复认成立** | 23:21 |
| S47 | `sed -n '1,40p' internal/observe/clock.go` | 单调钟纪律（D42#9）逐字：超时/截止**必须**走单调钟，`observe.Timeout`（`NewTimeout:34`／`Elapsed:39`）是具名 sanctioned 机制；墙钟差值做超时常量＝**禁** | 23:21 |
| S48 | `grep -n "go func\|filepath.Clean\|Allowlist" tools/d22scan/main.go` ＋ `sed -n '1,60p'` | 豁免只有一条通道：`tools/d22scan/allowlist.txt`，格式 `ban-id<TAB>repo-relative path prefix<TAB>reason`，**从不通配**；`loadAllowlist:238`。**本程对该文件零字改动** | 23:22 |
| S49 | `sed -n '185,215p' .github/workflows/ci.yml` | staticcheck 步骤注明它跑 **`GOOS=linux`（34 枚 finding）与 `GOOS=windows`（78 枚 finding）** 两个平台形状 ⇒ **`_windows.go` 代码在 ubuntu CI 上确实被类型检查过（但不跑测试）**；该步 `if: ${{ !cancelled() }}` | 23:22 |
| S50 | `git diff -- .gitignore`（只读别人那枚脏项） | 别人正在加的只有 `.worktrees/` 一段（@@ -44,6 ＋15 行）⇒ **`assets/web/dist/`（`:24`）是 HEAD 里的既有规则，不是本轮新加的** | 23:23 |
| S51 | `wc -l internal/proc/treemetrics_windows.go` | **298 行**（⑨ 第 9 条引的就是这个数） | 23:26 |
| S52 | `find . -maxdepth 3 -type d -name assets` | 全树 3 层内只有 **`./design/old/assets`** 与 **`./frontend/dist/assets`** 两枚 ⇒ 票面 `assets/web/dist` 在任何既有目录下都没有落点 | 23:26 |
| S53 | `ls internal/panel \| grep -i "host\|windows"` ／ `find internal/panel -name "*_windows*.go" \| wc -l` | 前者**退码 1＝零命中**、后者＝**0 枚** ⇒ `internal/panel` 今天是**一枚零平台分叉包**（与 S38「它在 ubuntu 的 core scope 有分母」互为因果；②的落点问题就压在这一条上） | 23:26 |

**尺数合计：53 枚**（S1–S53）。**未跑清单（具名）**：`go build`／`go vet`／`go test`／`./...` 形态／`go get`／`go mod tidy`／`GetInstalledVersion()` 实调／WebView2 真拉起——前三类是派单硬禁，后两枚是只读腿无载具，全部登记在 ⑥⑨⑩。

---

## ⑨ 我可能写错的条目（对抗我自己，必交）

> 逐条摆"我打算这么说／哪一发的证据其实没那么硬／被推翻的代价"。**不许写"无"。**

1. **`cmd/wisp` 的依赖名册我没量。** 派单禁止因 `246-r2` 在写它而跑 `./...`，我把这条精神推广到"也不对它跑 `go list -deps`"（S3 现量：`cmd/wisp` 那 5 枚文件此刻是 `M` 态）。⇒ ②里凡是"落在 `cmd/wisp` 会多出哪几名"的算法，**是从库侧（S16）推的，不是从 `cmd/wisp` 测的**。若有人现在偷跑，读数可能是写到一半的脏集。
2. **"纯 Go 无 cgo"我是三把尺拼的，不是构建验的。** S11（`import "C"` 零命中）＋S12（`#cgo`／`go:cgo_import_dynamic` 零命中）＋S10（40 枚 `.go`／3 枚 `.dll`）。`go-winloader` 里**确有一枚 C 源文件** `tinydll/tinydll.c`（S12），但它没有对应的 `import "C"`，是**被嵌成字节**的样板。⇒ 结论"不需要 C 工具链"我按"无 cgo 指令"下，**没按"真的 `go build` 过"下**；`02-spike-report.md:211` 反而具名写 spike 全套"需 Go 1.27.1 ＋ **mingw64 gcc**"（那是 sherpa-onnx cgo 那一支要的，不是 go-webview2 要的）——**这两件事极容易被我读成一件事，请复核。**
3. **我把 `Embed` 的嵌套泵说成"会吃掉球的泵"，这一步有推断成分。** 直接证据只有 S21（`:96-111` 那段循环是库里的、`GetMessageW(hwnd=0)` 取本线程全部消息）。**"所以球仍能出帧"是我从 Win32 消息队列是线程级这一条推的，我没有跑过**（只读腿无载具，S22/S31 都没法证）。⇒ ③里凡是"实测会/不会卡"的句子我都写成"待量＋怎么量"，没写成结论。
4. **`log.Fatalf`＝进程死，我按 Go 语义说，未在本机验证。** S22 显示 `chromium.go:173/:188` 用 `log.Fatalf`。我的推断"异步失败即整进程退出、与票面 AC#5『no-crash』直接冲突"成立的前提是这两支真的被走到；**缺 runtime 时到底走同步支（`:89/:92` return false）还是异步支（`:173` Fatalf），S46 那种目录存在性证明不了**。⇒ 这一条在 ⑩ 里升格成待裁（我拿不定，且它决定逃生通道怎么写）。
5. **`GetInstalledVersion()` 能当预检我用的是源码读，不是现量。** `webviewloader/module.go` 的 `GetInstalledVersion` 在 S46 那台机器上必然返回非空串，我**没跑过**；而且它依赖 `WebView2Loader.dll` 可载（磁盘或内嵌字节，`module.go` 的 `loadFromMemory` 两路），**组策略禁用 WebView2 时返回什么我没查**。⇒ ③④里"预检再创建"的接法我标了"未实测"。
6. **"本机离线可拉"我只证到缓存里有 zip/mod/info（S9）＋`go list -m all` 不联网即成（S19）。** 我**没证**"并进根 `go.mod` 后 `go build` 也不联网"：那要 `GOFLAGS=-mod=mod` 改契约文件（严禁），且 `GOPROXY=https://goproxy.cn,direct`（S7）意味着**真缺件时会去镜像站**（见第 12 条）。
7. **`assets/web/dist` 我扫了两个解释口径，仍可能有第三个。** S27（仓根不存在）＋S28（`internal/panel/` 下不存在）。若那枚 embed  intends 落在**新包**目录（如 `internal/panel/host/assets/web/dist`），我这条尺没伸到，因为**目录根本不存在，无父可举**（**S52**：`find . -maxdepth 3 -type d -name assets` 现量＝只 `./design/old/assets` 与 `./frontend/dist/assets` 两枚；另 **S53**：`ls internal/panel | grep -i "host\|windows"` 退码 1、`find internal/panel -name "*_windows*.go"`＝**0 枚** ⇒ `internal/panel` 今天是**零平台分叉包**）。⇒ ④的结论我限定成"票面那条路径按仓根与按 `internal/panel` 两种口径都不存在"。
8. **我没读 `frontend/**` 一字，所以"现用 embed 是 `all:dist`"是二手。** 出处＝S30 那批仓内证据件（`156-…:762` 逐字写 `frontend/embed.go:19 → //go:embed all:dist`）。两层禁令我按派单走（不读、结论也不引到它身上）：④只报**文件在不在、多大、路径对不对**（S26），不报内容、不判它好坏。⚠ 但**"dist 里那两枚 hashed 文件名对得上 `index.html` 引用"这种话我说不了**——那是 `panel.Assets.Check()`（S25 `:158`）该干的事，不是普查该干的。
9. **`treeTCPConnections` 我读成"不能直接拿来当『无监听端口』的尺"，是基于它不读 `state` 字段（S40）。** 如果它其实通过别的入口暴露了 state，我就把话说重了。我只看了 `:269-300` 那一段＋枚 `^func` 清单（S40），**没通读 `treemetrics_windows.go` 全文（S51：`wc -l`＝298 行）**。⇒ ⑤里我把建议写成"新增/导出一枚读 state 的枚举"，没写成"这仓做不到"。
10. **"internal/panel 不在 windows CI scope"我只核了 `portable-tests.sh` 的 `win_pin`（S38）与 `ci.yml` 的 runs-on（S37）。** `ci.yml:533` 那枚 job 叫 "Build wisp.exe ＋ SLO smoke gate"，**它 `go build` 出来的 exe 里含 panel**，build 失败也会红。⇒ 所以准确说法是：**panel 的 windows-tagged 用例在 CI 零分母，但 panel 的 windows-tagged 代码在 CI 有两处分母（`build wisp.exe` 的链接、staticcheck 的 `GOOS=windows`，S49）**。我把这两件事说混＝⑤⑥会误导排程。
11. **延迟两数我用的是 09-19 的两轮读数（S34），未复跑，且时序口径是 spike 自己定的。** `:82-84` 的 cold 定义含 `SetHtml` ＋**首次浏览器往返**；本票改成"embed 供给＋`AddWebResourceRequestedFilter`"后，**往返里多了一趟 Go 侧读 embed 的分支**，这条 delta 我**没有数**。⇒ ⑥只交方法不交数。
12. **`GOPROXY=https://goproxy.cn,direct`（S7）与 d22scan ban #5 `mirror-hash`（S41）会不会打架，我判不动。** ban #5 的射程是"**hash material** mentioned together with a mirror"，而 `go.sum` 的 `h1:` 就是 hash、goproxy.cn 就是镜像——**但 `go.sum` 不在 d22scan 的 Go 产码射程里，且今天 S18 那 4 行 hash 就是 spike 从同一 proxy 拿的**。⇒ 我不下结论，摆给编排者（⑩第 5 条）。
13. **超预算风险自陈**：派单没给硬顶，但两枚先例（`33-a1` 42/35、`33-r2` ~50/35、`33-r3` ~47/30）都超。本程到 S50 为止的调用数我会随 commit 报，**不美化**；①-⑦的结论若在 100 轮内写不满，按派单闸门"先把⑧⑨⑩写满并 commit"办，**不拿放宽结论凑数**。

---

## ⑩ 判不动的地方（逐条摆甲／乙／不做＋现量再写，交编排者裁）

> 未定义即停（D22 闸门③／`AGENTS.md` §0.1）。**每条都写"现量到哪一步"，缺的都是裁，不是猜。**

**J1（⇒③）WebView2 该投在哪枚线程：`ui-sta` 共用 vs 独立 STA 线程。**
现量：票面 `:30` Key constraints 逐字"WebView2 created on the shared `ui-sta` STA thread (D38a)"；S31 `sta_windows.go:52/:62/:77` 已占死一枚 STA；S35 `panel-host` 在 **`TemporaryNames`（`goroutine.go:58`）**、S36 `SPEC-01:122` 注明"**均在对应 DisposalScope 内**"，而常驻 6 枚名册（`ResidentNames:44-46`）里没有它 ⇒ **"另开一枚常驻 STA 线程"这个名字在 D38b 名册里今天不存在**，硬用 `panel-host` 就是把"临时·在处置范围内"的名字拿去当常驻。
　甲＝按票面走 `ui-sta`（名册零变更，代价＝S21 那枚嵌套泵压在球身上）。
　乙＝新线程（要 D38b 名册变更＝**人工批准**，我不碰）。
　不做＝我自行选乙并起名（＝改契约，禁）。
**请裁：甲，还是先量一遍甲的代价再谈乙。**

**J2（⇒③/④）缺 runtime 时库走同步支还是异步支（决定 AC#5『no-crash』怎么写）。**
现量：S22 两条路都在库里（同步 `:89/:92` `return false`；异步 `:173/:188` `log.Fatalf`）；S46 本机装了运行时 ⇒ **本票在这台机器上永远走不到那支**。
　甲＝本票照 AC#5 做 fixture（改名/遮 loader），并**先在 fixture 上量到底命中哪一支**，若是 `log.Fatalf` 那支 ⇒ 具名上报"必须绕开 `Embed`、自写异步创建路径"。
　乙＝现在就判定它必是 `log.Fatalf`，直接要求自写 COM 创建（＝新增约 100 行 unsafe 产码，且 S23 证明库未导出可复用的 handler）。
　不做＝把"runtime 缺失也不崩"当已成立勾上去。
**请裁：甲（我倾向），但甲要真载具，不属本普查射程。**

**J3（⇒⑥）1808.9 ms 那一形（进程内首次 create、竞争态）算不算本票的冷启口径。**
现量：S34 cold P95 = 1041.6/1256.4 ms（均 <1500 且 <2000）；但同表"全进程内首次 create"**1808.9 ms＞1500**（仍 <2000）；而本票落地形态恰恰就是"在已经跑着球＋音频＋DB 的常驻进程里首次 create"。`AGENTS.md` §2 与 `SPEC-12 §4.2` 把"WebView2 冷拉起 >2s → 重评 L2 卡（P11，S0）"列为**未定义即停**；`33-a1` §④ 又具名报"`grep -n "P11" docs/reports/pending-and-issues.md`＝0 命中，台账无对应 `A##`"。
　甲＝把 1808.9 认作本票冷启口径的候选真值，落地腿**必须**在常驻进程内测一发，超 1500 就具名上报（不放宽 D32、不改 `thresholds.go`，S 见 ⑥ 的硬约束）。
　乙＝沿用 spike 的 1041.6/1256.4 作口径，把 1808.9 记成"竞争态离群"。
　不做＝我顺手把冷启阈值改成 2000（＝改契约，禁）。
**请裁：口径选甲还是乙；并请先给 P11 立一枚 `A##`（台账现量 0 命中这条我没复核到，转引 `33-a1`）。**

**J4（⇒②）"装配根 `cmd/wisp` 是唯一接缝"与"正向依赖边一律不开"能不能推广到宿主。**
现量：这两条是本仓既有裁定（派单 `:②` 具名提醒），但**我没有在 `PLAN.md`／`SPEC-*` 里读到它们对"原生窗口宿主"说过一个字**；`docs/evidence/s1/114-ac1-status-table.md:203` 反而把"票 114 与票 33 的地界"列成**未定需 owner 拍**。⇒ 我这一问是**我自己的推广**，不是仓里本来就有的规矩。
　甲＝宿主落**新包**（如 `internal/panel/host`），`cmd/wisp` 只注入（沿用 238 的形状，但这是我的推广）。
　乙＝宿主落 `cmd/wisp/panel_host_windows.go`（沿用 197 的形状，同样是推广）。
　不做＝我选一个并写成"本仓规矩如此"。
**请裁：甲/乙（②里两形的依赖增量我都算好了，见 ②）。**

**J5（⇒①）从 `goproxy.cn` 拿 go.sum hash，撞不撞 ban #5『mirror-hash』。**
现量：S7 那枚 GOPROXY；S18 spike 已用同 proxy 落了 4 行 hash；S41 ban #5 的原文射程＝"**hash material** mentioned together with a mirror (C29: hashes come from the signed manifest only)"，且扫描射程是 `internal/`＋`cmd/` 的 Go 产码。`deps.toml:22-26` 具名示范过镜像写法（`mirror_prefix = "https://ghfast.top/"`）并靠 **sha256 钉**兜底。
　甲＝按现状加两名 module（hash 由工具链自取），我**不**动手写任何 hash。
　乙＝先请人工确认 go.sum 是否需要与 deps.toml 同一套"官方源 hash ＋ 镜像取件"规矩。
　不做＝我为了绕 ban #5 去手写/改写任何 hash（禁，且本票不改 go.sum）。
**请裁：甲是否即可。**

**J6（⇒④）票面那条 `//go:embed assets/web/dist` 与仓里活的 `frontend/embed.go` `all:dist` 是什么关系。**
现量：S27/S28（路径两种口径都不存在）＋S29（`.gitignore:24` 逐字 `assets/web/dist/`，**整条 ignore、无锚文件豁免**）＋S30（`156-…:762` 具名 `frontend/embed.go:19 //go:embed all:dist`）＋S25（`internal/panel.Assets.Resolve` 已备好供宿主用的读接口）。Go 侧语义（`pkg.go.dev/embed`，原句）："Matches for empty directories are ignored" ＋ "each pattern in a //go:embed line must match at least one file or **non-empty directory**. If any patterns are invalid or have invalid matches, **the build will fail**"。
　甲＝宿主直接复用 `panel.Assets`（S25），票面那句按"已过期的路径写法"处理，**票面我不改一字**，由编排者裁是否留一行更正记录。
　乙＝真要新立 `assets/web/dist` ⇒ 先要动 `.gitignore:24`（那是别人正在弄脏的文件之一，S3/S50）并加锚文件，且**清检出时会硬失败**（上一条 Go 语义）。
　不做＝我改票面或改 `.gitignore`。
**请裁：甲/乙。这是"构建期会不会硬失败"的总开关。**

**J7（⇒⑤）「无监听端口」这条判据要不要进 CI，以及进了谁的scope。**
现量：S39 全仓零 `netstat`；S40 唯一能枚举 TCP 的是 `internal/proc` 里**未导出、不读 state** 的一把；S37/S38/S49 给出三种分母；S38 `internal/panel` 不在 win_pin。⇒ **"该判据只能在真机成立"的风险是实的**，而本仓已有三枚同形坑（派单 `⑤` 具名）。
　甲＝判据写成能力型（本进程 PID ＋ 树内 PID 的 TCP 表里 `state==LISTEN` 计数＝0），**并把"它跑在哪台机器"随读数具名**，跑不到的那格**登进台账**（新增 `A##` 由编排者落）。
　乙＝同时把 `internal/panel/` 追加进 `win_pin`（＝改门禁分母，`A374` 里编排者已自决"暂不加"，加要单独落 `A##`）。
　不做＝我自行往 `portable-tests.sh`/`ci.yml` 里加 scope（本程零写面）。
**请裁：甲必做；乙要不要现在做。**

**J8（⇒⑦）票 248（`internal/panel` 设置路由）与本票落地腿的先后。**
现量：S41 白名单四枚零 config（票面 `:16` 已具名"点设置录 key"归票 248）；owner 09-30 那句"我要点击设置自己配置"要的是**两枚都落地**。两枚同包 ⇒ 串行。
　甲＝先 33 宿主（窗口能开），再 248（能点设置）。
　乙＝先 248（白名单＋路由），再 33（一接上就能用）。
　不做＝并行开两枚写腿打同一包。
**请裁：甲/乙（我给的顺序建议见 ⑦，理由是"宿主是白名单的唯一真听众"，但那是我的理由不是仓里的规矩）。**
