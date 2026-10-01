# `33-r4` — 面板宿主仪器腿（⛔ 只改 `cmd/wisp/*_test.go`；六处"看起来在测、其实看不见"的尺装牙）

- 程：`33-r4`（**写腿·仪器**）｜工单：`.scratch/wisp/issues/33-panel-host-c27.md`（票 33）
- 判据出处（逐格指回）：`docs/evidence/s1/33-panel-host-c27-v1.md` §AC#1／§AC#2／§AC#3／§AC#4／§AC#12 与 §A#24/#25/#26/#27/#30/#31
- ⛔ 本件**不勾票面任何一枚 `- [ ]`／`- [x]`**（勾框归编排者与下一枚验收腿）；⛔ 不为变绿放宽任何既有断言。
- ⛔ 产码那半边（`cmd/wisp/panel_host_windows.go`／`resident_*.go`／`internal/ball/**`／`internal/proc/**`）本程**一字未改**，归 `33-r2`。

---

## 起手锚点

| 尺 | 现量（同发取，`date` 与 `git log` 一条命令里跑） |
|---|---|
| `date "+%Y-%m-%d %H:%M:%S %z"` | `2026-10-01 11:30:24 +0800` |
| `git log -1 --format=%H` | `eed229e48cb9599c29f72fea6a1170739a5f2b41` |
| `git rev-parse --abbrev-ref HEAD` | `dev` |
| `git status --porcelain -- cmd internal` | **0 行**（起手名册＝空；终态必须等于这一枚名册） |

⚠ 起手锚的用途：收尾用 `git diff eed229e4..HEAD -- cmd internal` 复算⛔ 突变体零枚进提交；HEAD 期间会被别枚腿推进，所以本件**每条读数都带取数时刻＋当时 HEAD**。
⚠ 同刻在飞的只读腿（我不与它们抢写面，也不动它们在读的文件）：`33-a2`（读 `internal/ball`＋webview2 泵面）、`228-a3`（读 `internal/ball` 托盘面＋`internal/agent/approval`）。

---

## ① 六格逐格：改前读数／改后读数／反控读数

### 格 1 — netstat 那把尺（AC#3 的 L2 半，恒真）
- 改前（源码现读，`11:3x`，HEAD `eed229e4`）：`cmd/wisp/panel_host_windows_test.go:115` 逐字 `tcpTableOwnerPIDAll = 4`；`:116` 逐字 `tcpStateListen = 10`；`:126` 传的是那枚 `4`；`:125` 只传 `windows.AF_INET`；`:145-146` 重试耗尽 ⇒ `t.Logf(...)` ＋ `return 0`。同仓对照真值 `internal/proc/treemetrics_windows.go:57` 逐字 `tcpTableOwnerPIDAll = 5`。
- 改后读数：待回填（本行在第 100 发前必须换成读数）
- 反控读数：待回填

### 格 2 — 焦点回还（AC#4）
- 改前：`:238` 的 `prior := windows.GetForegroundWindow()` 取在 `:209` 的 `HotShow` **之后**；`:249-251` 那一跳只有 `t.Logf`；全仓对 `prevFocus` 零枚断言（v1 §A#27 的尺我自己会重跑）。
- 改后读数：待回填
- 反控读数：待回填

### 格 3 — HWND 身份（AC#1 的 "reuses window"）
- 改前：`:175`／`:237`／`:260` 三处只判 0／非 0，从未比较 hide→re-show 前后是否同一枚句柄。
- 改后读数：待回填
- 反控读数：待回填

### 格 4 — process-tree 的分母（AC#1 的 child-count）
- 改前：`:89-105` 是 `CreateToolhelp32Snapshot` 上按**进程名**全机数 `msedgewebview2.exe`，不按父 PID 走树；v1 §A#26 读数＝基线 14 枚（别人会话）／窗口期 20 枚（我们贡献约 6）／销毁后回 14。
- 改后分母读数：待回填
- 反控读数：待回填

### 格 5 — AC#12 那枚空尺
- 改前：`cmd/wisp/panel_host_gate_test.go:52-76` 全文除 `:62` 的 `t.Fatalf`（git 报错）之外**零枚 `t.Errorf`**；非锚条目与计数都只有 `t.Logf`。
- 现量（同发，`11:3x`，HEAD `eed229e4`）：`git ls-files frontend/dist` ＝ 1 枚 `frontend/dist/.gitkeep`；`git status --porcelain -- frontend` ＝ **空**；`git status --porcelain --ignored -- frontend/dist` ＝ **2 行 `!!`**（`frontend/dist/assets/`、`frontend/dist/index.html`）；`git check-ignore -v frontend/dist/index.html` ＝ `frontend/.gitignore:12:dist/*`。⇒ 工作树有产物、入库没有，且那些产物是 **ignored** 态（普通 `status` 看不见＝上一枚腿没写进表的口径细节）。
- 改后读数：待回填
- 反控读数：待回填

### 格 6 — latency 那格的两枚小账（AC#2）
- 改前：`:181` 那句 `t.Logf` 里 `"HEAD 7a41db9b"` 是硬编码字面量；票面要的 10-run P50/P95 在仓里不存在（v1 §A#20 的 `grep P50|P95`＝0 命中，我自己会重跑）。
- 改后读数：待回填
- 反控读数：待回填

---

## ② 起跑绿名册与终跑名册（逐名 diff）

待回填（起跑 `-v` 全名册已在跑，日志 `.scratch/wisp/probes/33/r4/logs/baseline-start.txt`）。

---

## ③ 门禁四数（go build ./...／go vet ./...／d22scan／gofumpt，逐条真实读数）

待回填（终态逐条现跑）。

---

## ④ 我跑了哪些尺、每条真实读数

> 本节在**前 15 发工具调用内写满**（编排者硬要求），后续只追加、不删行。

1. **起手同发四读数**（`date`＋`git log -1 --format=%H`＋`git rev-parse --abbrev-ref HEAD`＋`git status --porcelain -- cmd internal`，一条命令里取）：`2026-10-01 11:30:24 +0800`／`eed229e48cb9599c29f72fea6a1170739a5f2b41`／`dev`／**0 行**。（时刻与锚点即 §起手锚点 那一发，同发。）
2. **验收表全文已读**：`docs/evidence/s1/33-panel-host-c27-v1.md` 219 行全文（§A#24/#25 常量错位与 AF_INET6 探针读数、§A#26 分母、§A#27 焦点、§A#30/#31 AC#12、§D 逐格判语、§E 收尾三把尺）。⛔ 我没有把它的任何读数抄成本程的交付读数——本件凡引用它的地方都标了"v1 §…"出处。
3. **被测仪器源码逐行现读**（不是抄表）：`cmd/wisp/panel_host_windows_test.go`（273 行全文）、`cmd/wisp/panel_host_gate_test.go`（180 行全文）、`cmd/wisp/panel_host_windows.go`（411 行全文，⛔ 只读）、`internal/proc/treemetrics_windows.go`（对照常量真值）。
4. **`cmd/wisp` 的 run 前提（派单点名那枚）**：`PATH` 形态照 `scripts/wisp-cli-tests.sh:100-106` 抄——脚本原话是 `pwd -W` 那种带盘符的形态会 **exit status 0xc0000135**（用例根本没跑），必须用 shell 自己的路径形态：`export PATH="$dll_dir:$PATH"`，`dll_dir="$root/third_party/sherpa-onnx"`。本程命令固定为
   `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" GOFLAGS= go test ./cmd/wisp/ ...`（与 v1 §A#9 同形）。
   现量：`third_party/sherpa-onnx/` 里 `onnxruntime.dll`／`sherpa-onnx-c-api.dll`／`sherpa-onnx-cxx-api.dll` 三枚**在**（`ls` 现读，`11:3x`）。
5. **起跑基线（11:34:xx 起钟，后台跑，⛔ 不与其它腿并发）**：`GOFLAGS= go test ./cmd/wisp/ -count=1 -v`（带上面那枚 PATH）→ 日志 `.scratch/wisp/probes/33/r4/logs/baseline-start.txt`。名册与逐名绿／红待收（§②）。
6. **撞钉预检①（新符号名）**：`grep -n "func TestMain|func listenSocketsForPID|func countWebviewChildren|func treePIDs|panelHostCold|panelHostHot|func percentile|func gitHeadShort|func repoRootForTest" cmd/wisp/` ⇒ 现量 **只有两枚已存在的函数**（`panel_host_gate_test.go:79 repoRootForTest`、`panel_host_windows_test.go:89 countWebviewChildren`、`:110 listenSocketsForPID`）。⇒ `cmd/wisp` **没有 `TestMain`** ⇒ 本程若要用 TestMain 汇总是安全的（不会与既有 TestMain 撞名）。
7. **撞钉预检②（行为型钉，派单第 64 条那味）**：`scripts/portable-tests.sh:195` 的 cli scope 与 `scripts/wisp-cli-tests.sh:112` 只跑包、不按用例名取数 ⇒ 新增/重命名用例不会打断 CI 取数点；`internal/panel/composer_dispatch_test.go` 那枚反转钉扫的是**产码标识符**（v1 §A#12），`_test.go` 不入射程 ⇒ 我在测试文件里新增 AST 扫描不会打红它。⚠ 这两条是**读源码得出的**，不是跑出来的；跑出来的复认在 §③ 门禁那发（终态 `go test ./internal/panel/` 红名册必须仍等于起手的 4 枚）。
8. **AC#12 的三种口径现量**（见 §①格 5 那一发）：入库＝1 枚 `.gitkeep`；普通 `git status`＝**空**（因为 `frontend/.gitignore:12` 的 `dist/*` 把整批产物标成 ignored）；`--ignored` 口径＝2 行 `!!`。⇒ 这枚尺以后必须区分"tracked／ignored-dirty／embed 内容"三条轴，⛔ 只数 `git ls-files` 的那枚旧尺看不见第二条轴。
9. **禁区自证（起手态）**：`git status --porcelain -- cmd internal`＝**0 行**；`internal/observe/thresholds.go`／`docs/SLO.md`／`docs/PLAN.md`／`docs/specs/**`／`tools/d22scan/allowlist.txt`／三枚冻结件起手均未动（终态再取一发同尺）。⛔ 本程不跑 `go mod tidy`（v1 §A#29 现量它在 HEAD 上 `exit 1` 会造出不属本程的 diff）。
10. **`frontend/**`／`design/**` 两层禁令的自证**：本程只跑过 `git ls-files`／`git status --ignored`／`git check-ignore` 三枚**条目名**尺，⛔ 未打开过任何前端或设计文件，⛔ 结论不引到 `frontend/**` 身上（AC#12 那格的判据一律从 `panel.Assets` 的 Go 侧能力问）。

---

## ⑤ 我可能写错的条目（对抗我自己）

1. **修正常量不等于那把尺就有牙**。我把 class 改成 5、state 改成 2 之后，仍必须做**真开一枚 `127.0.0.1` 监听 socket** 的正控；如果正控不红（尺看不见自己进程刚开的端口），那说明我的行 layout 或归属口径还是错的，这一格判"仍无牙"并上报，⛔ 不许拿"常量改对了"交差。
2. **行 layout 是我自己算的，不是抄来的**：`MIB_TCPROW_OWNER_PID`＝6 枚 uint32＝24 字节、state 在偏移 0、owningPid 在偏移 20；`MIB_TCP6ROW_OWNER_PID`＝48 字节（localAddr16＋localPort4＋remoteAddr16＋remotePort4＋state40＋pid44）。v1 §B#14 在 `AF_INET6` 上取到过"表里有 LISTEN 行、own-pid 却对上 0 枚"，所以**IPv6 那一半我很可能又取错**；我只在本程自己的正控真跑通之后才主张任何 IPv6 结论，跑不通就具名写"这一维我没量到"。
3. **重试耗尽静默 `return 0` 的修法可能造出新假红**。我把它改成 `(count, error)` 后，调用方 `t.Fatalf` 会在**机器负载高**时红——那是"量不到"不是"违规"。这一枚红要有自证（错误句必须写清是工具失败），⛔ 不许反过来把它又降级成 `t.Logf`。
4. **焦点那一跳可能被我修成"两个都恒真"**。产品 `Show` 是在自己内部记 `prevFocus`（`panel_host_windows.go:269`，取在 `ShowWindow` 之前），所以只要测试的取样时刻也在 Show 之前，`afterHide == beforeShow` 就**几乎必然成立**——这正是"把断言写在产码已经做的那一步之后"的形状。⇒ 我必须补两枚才谈得上有牙：断言**面板 Show 后确实拿到前台**（否则"回还"这个问题根本不成立），以及**断言回还目标不是刚被藏起来的那扇面板窗**（today 的病形）。这两枚哪一枚在真机上是恒绿，我要具名说。
5. **焦点断言可能受 Windows 前台锁影响而 flake**。`SetForegroundWindow` 对非前台进程可能被拒（v1 §A#27 那句注释就是这意思）。⇒ 我不为它加 sleep 重试把它拧绿；如果它 flake，我把每一发的三个句柄逐名落表，交编排者判是不是产码那一跳真不稳。
6. **HWND 身份断言可能被 recreate 路径合法打断**。`Destroy` 之后 `Show` 会造**新**窗口（`:319-330` 明写），所以"同一枚 HWND"只在 **hide→re-show** 这一对上成立。⇒ 断言只放在 hide→HotShow 那一跳；dispose 之后另说，⛔ 不许把它写成"整个生命周期 HWND 不变"。
7. **process-tree 的分母我可能数窄了**。WebView2 的子进程可能不是测试进程的直接子级（loader/服务重挂父）；按 `ParentPID == 本 pid` 直接走会数到 0，那把"我们的树"读成"没人起窗"＝新的假绿。⇒ 我从 `os.Getpid()` 往下**递归**收集整棵子树，并具名报"直接子级枚数／全子树枚数／其中 msedgewebview2 枚数"三个数；如果全子树也是 0，这一格我判"改后仍量不到"，不硬说它有牙。
8. **AC#12 那枚尺的两个危险形状**：① 把"工作树有产物"读成"AC#12 已满足"（那是 33-v1 已判过的假绿，入库口径仍只有 `.gitkeep`）；② 把"入库只有占位"判成一枚**红**，替界面侧那枚 agent 背锅——归口未落是票面 `:37` 的原话，⛔ 本程不许拿红替编排者做那个决定。⇒ 我的判据只钉两条不变式（Built() 与 fail-closed 必须自洽；Built() 为真时入口引用必须全能解析），归口那一格写成**具名 skip 或具名读数**，不写成绿。
9. **反控可能需要临时改产码**：本程唯一允许的形态是**仓外副本**里改；若我不得不在树内做，改完立刻 `git cat-file blob HEAD:<path> > <path>` 并取还原前后 md5 相同。⛔ 绝不用 `checkout`／`restore`／`stash`，⛔ 突变体零枚进提交（编排者用 `git diff eed229e4..HEAD -- cmd internal` 复算）。
10. **10-run P50/P95 是自相关样本**：同进程 `-count=10` 的第 2 到第 10 发会吃到第一发留下的 WebView2 缓存/常驻子进程，把冷读低。⇒ 分位数落表时逐发带时刻与冷/热原文，并具名写"同进程连续 10 发"这一口径，⛔ 不写成"冷启真实分布"；P11 那道 2s 重评线我按同口径说明未触发。
11. **⛔ 不许把 `internal/panel` 起手那 4 枚红算成本程的**（`tokens_fourway_test.go:441` 落在**冻结件**上，红因＝工作树那 16 枚 ` D` 的 `design/**`）。⇒ 终跑名册按**逐名**比，不比退码；本程只要求 `cmd/wisp` 侧不新增红，新增的红只允许是"我给尺装牙装出来的、且红因具名"那几枚。
12. **时序 flake `TestAC228ExitRequestDuringBootStillLeavesThroughD38E` 可能在本程起跑/终跑里换态**。派单点名它起手即在。⇒ 偶发红用 `-count=5` 隔离复认并具名登记，⛔ 不改它、不放宽它、也不把它算成本程新增红（除非 `-count=5` 里它枚枚都红，那我如实上报）。
13. **硬编码 HEAD 的修法本身可能引入一枚新尺错**：`git rev-parse --short HEAD` 在 `go test` 的 cwd（`cmd/wisp`）里能取到仓根，但在**仓外副本**里 archive 出来的树**没有 `.git`** ⇒ 那枚 helper 若 `t.Fatalf` 会在副本里炸。⇒ 我把它写成"取不到就具名报 `HEAD-unknown`"，⛔ 不取用假值。
14. **d22scan 的 emoji 扫描（ban #8）可能打到我的新测试文件**：`cmd/` 在 ban #8 的走树里（v1 §A#34 现量 `ban #8 cmd/=78` 是分母），注释豁免、**字符串不豁免**。⇒ 新写的字符串里我不用 `✓`(U+2713)、`≤`(U+2264)、任何 U+1F000–U+1FAFF；终态 `sh scripts/d22scan.sh` 必须仍 exit 0。
15. **票面进度行**：往 `.scratch/wisp/issues/33-panel-host-c27.md` 的 Progress log 追加时**只用 Edit/Write**，⛔ 不带反引号的 heredoc（记忆里那族"静默吃字"事故），⛔ 不改上面任何一格、⛔ 不碰 AC 复选框。

---

## ⑥ 判不动的地方（逐条甲／乙／不做＋现量）

1. **格 2 的"回还那一跳在产码里到底做没做"**。现量（我起手自己读的行）：`panel_host_windows.go:269` 确实记 `prevFocus`，`:311-313` 确实 `SetForegroundWindow(prev)`＋`SetFocus(prev)` ⇒ 与 v1 §D 那句"机制代码存在、证据为零"同向，所以真断言**不一定红**。⚠ 但 `:274-276` 只在 `!shown` 时才 `SetForegroundWindow(hwnd)`，而 `:269` 的记值发生在 `ShowWindow` 之前——这两行合起来意味着"面板没拿到前台时记的 prev 是别人、拿到了之后记的 prev 是自己"，**这条链路谁对谁错只有真机能定**。甲＝我把三枚句柄逐名断言，红就报红因在产码（归 `33-r2`）；乙＝要我先改产码——⛔ 本程写面禁；不做＝⛔ 不许把断言降级成 `t.Logf` 换绿。
2. **`AF_INET6` 那一族要不要进断言射程**。现量：v1 §A#24 探针在 class=3 的 `AF_INET6` 表里读到 3 枚 `state=LISTEN` 行，但 own-pid 归属对不上（它自己 §B#14 承认可能是 offset 错）；本仓 `internal/proc/treemetrics_windows.go:277` 的既有采样器**也只传 `AF_INET`**。⇒ 这一格是**射程决策**（要不要让 AC#3 覆盖 IPv6），⛔ 不是我这一发该替编排者做的扩范围。甲＝我只量、只报（自己开 `[::1]` 监听的正控能不能被同一套 layout 数到）；乙＝编排者点头我才把 IPv6 并进断言；不做＝⛔ 不拿"今天这机上 webview2 没开 IPv6 端口"当"不需要覆盖"的理由。
3. **AC#13（冷启探测页盖掉真页面）**。现量：`bringUp` 的 `:221 serveEntry()` 之后 `:227 firstRoundTripLocked(...)`，而后者在 `:374-375` 又 `SetHtml` 一枚自造探针页。⇒ 这是**产码时序**，票面 `:39` 那格的判据（重排次序／过滤器＋一枚会响的"最终文档含真入口"断言）落点在 `panel_host_windows.go`，本程⛔ 不许改。甲＝本程只把"显示的是探针页还是入口页"这条**量出来**（读 embed 入口字节的特征 vs 探针页特征，从 Go 侧字符串比，不读前端文件），交 `33-r2`；乙＝要我现在钉成断言 ⇒ 会造出一枚"本程自己判自己产码缺跳"的红，归属混淆；不做＝⛔ 不留空尺。
4. **`PanelManager` 生产调用者枚数**（dispose 半、也是 owner 那句"我要能点"的那半）。现量：v1 §A#14 全仓 `grep -rn "PanelManager"`＝定义文件＋测试文件两枚 ⇒ **0 枚非 test 调用者**。这一格归 `33-r2`（接常驻）与票 248，本程⛔ 动不了；我能做的是把"dispose 在产品里不可达"从散文变成**一枚会随接线自动转红的具名 skip**。
5. **AC#12 的勾与归口**（谁把 `frontend/dist` 填上）。现量：入库 1 枚 `.gitkeep`、工作树 2 行 `!!`（§④#8）；票面 `:37` 逐字"这一格在'谁把 dist 填上'落定前勾不了"。⇒ 本程**只装牙、不判归口**，⛔ 不造页面产物、⛔ 不写 `frontend/**`。
6. **`winlive` 档在 CI 零岗位**。现量：本程所有真机读数（窗口、焦点、netstat、时延分位）只在开机这台机器成立；`scripts/portable-tests.sh` 的 cli scope 里没有 `-tags winlive` 这一档（v1 §A#19/#32 同判，编排者 `A489`/`Q-75` 第四次具名）。⇒ 每条读数带时刻＋HEAD，单开 workflow＝契约级，我不重开这一裁。
7. **`go mod tidy` 那味**（webview2 记成 `// indirect`）。现量：v1 §A#29 在仓外副本 `exit 1`，仓内 CI 无 tidy 门。⇒ 属 `go.mod` 禁区，本程⛔ 不跑不修；只在名册里说明" tidy 若被人跑一次会造出不属本程的 diff"。
8. **staticcheck 本机版与 CI 钉版不同**（本机 `2025.1.1`，CI `2026.2.1`，`ci.yml:181-200` 自陈旧版在 go1.27 export data 上导入期即崩）。⇒ 门禁四数里这一门**不复认**，交 CI；本件标〔未复认〕，⛔ 不跑它拿假结论。
9. **起跑名册的"没跑"与"绿"之分**。现量：`cmd/wisp` 缺 sherpa DLL 时是 `exit status 0xc0000135`＋`0.0xxs`＋**零枚 `--- FAIL`**（`scripts/wisp-cli-tests.sh:100-106` 逐字记录过这个坑）。⇒ §② 的名册只在日志里**真有 `--- PASS`／`--- FAIL` 行**时才算量到；若起跑那发是 `0xc0000135` 形，我具名写"这一维我没量到"，⛔ 不写"包全绿"。
10. **反控的载体归口**：格 1／2／5 的反控都可能需要"让产码或依赖变坏"。现量：派单允许只在仓外副本做，仓内做则须 md5 还原＋双尺复算。⇒ 我优先仓外（`git archive HEAD` 抽副本，与 v1 同一手）；副本没有 `.git`，凡我新写的 git 尺在副本里都走 §⑤#13 那枚 `HEAD-unknown` 分支，⛔ 不在仓内建 worktree／checkout。
