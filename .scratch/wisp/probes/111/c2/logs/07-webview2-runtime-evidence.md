# 07 · 托管那台 Windows 上 WebView2 Runtime 到底有没有（⛔ 不靠猜、⛔ 没开任何窗）

## 1 workflow 里没有任何装它/等它的步骤

见 `06-*.md` §2：`grep -in 'webview\|edge\|evergreen\|runtime' .github/workflows/ci.yml` ⇒ 只有 2 行散文注释（`:445`、`:755`），**零安装步骤**。⇒ 有没有它，只能从**在那台机器上取到的读数**判。

## 2 仓里已经有两枚腿量过这件事（具名归口，⛔ 本腿不新立）

### 2.1 `33-n1`（普查腿，件 `.scratch/wisp/probes/33/n1/verdict.md`，锚 `a5357fa1`，10-01）
- §4 逐字结论：**"GitHub 托管 `windows-latest`（= Windows Server 2025，清单页 `Windows2025-VS2026-Readme.md`，Image `20260922.246.2`）有 Evergreen WebView2 运行库；`windows-2022` 同样有"**，凭据是官方镜像追踪器里报障人与镜像方双方确认的版本读数（`131.0.2903.86` / `152.0.4191.66`）＋镜像方"它来自 Azure 基盘、我们不装也不锁"。
- §5 明确**不许用甲推乙**，并把"可见性/前台/像素/毫秒时序"记为**未定**，列出判乙所需的 5 枚读数 **R-a…R-e**（§5.3）——**那 5 枚到今天仍没采**（本腿复认：仓里没有 R-a..R-e 的任何件；`find .scratch -maxdepth 3 -type d -iname '*winlive*'`＝0 命中，`docs/reports/**` 里也没有那五栏的到账记录）。
- §7 E1 自陈最软一环：那两枚版本读数出自 issue 正文自述、**它没独立复核**。
- §6.1 结论：⛔ 不要在 CI 加安装步骤；§6.5 逐字：**"今天真正没测的原因：不是'没运行库'"**。

### 2.2 台账 `A526`（编排者自己读完一发真推送的 CI 日志，10-02 10:51，run `36956193382`／head `ed459d09`）
逐字：`wisp: panel window is up (test-harness, cold 231.2 ms, hot path 0.0 ms)`，同形 5 发（`cold 196.6 / 288.6 / 735.1 / 976.8 ms`），另有真句柄 `hwnd 0x40232` 等，`TestAC4FocusReturnToPriorWindowGap33r5` **PASS**。
⇒ 那一格写的是：**甲（那台机器装了 WebView2 Runtime）由【旁证级】升为【已证】**；⛔ 乙不升档；`TestPanelHostLatencyPercentilesAC2` 那发 `cold=3499.615 ms（budget 1500）` 是**那台机器自己的时序**、⛔ 不许当我们的性能读数。
⇒ **本仓已有一发"CI 上面板窗真起来了"的直接读数**，且早于本腿 6 天。

## 3 本腿自己新取到的直接读数（CI 原文，⛔ 不是二手）

尺＝`02-*.md` 那套切步＋`grep -h 'machine-wide\|hot re-show measured'`，逐发：

| run | CI 原文（逐字，节选） |
|---|---|
| 304 | `panel_host_windows_test.go:713: AC#1 pre-hide tree settle: reading held after 1 extra sample(s) at our tree webview=7 (machine-wide 7)`；`:749: AC#1 re-show tree settle: ... our tree webview=6 (machine-wide 6)` |
| 303 | 同形（`our tree webview=7 (machine-wide 7)` → `=6 (machine-wide 6)`） |
| 302 | 同形 |
| 301 | 同形 |
| 300 | 同形 |
| 299 | 同形（`machine-wide 7`） |
| 294 | `panel_host_windows_test.go:591: AC#1 denominators (head 941805d): our tree webview=7 direct-children=1 tree-pids=7 \| **machine-wide msedgewebview2=7** \| same HWND 0xc015a across hide->re-show=true`；`:596: hot re-show measured on this box: 38.649 ms` |
| 305 | **无该行**（该发用例终止在 `:662` 的 `t.Fatalf`，AC#1 那几栏根本没跑到——这本身是"红句缺哪一栏"的证据） |

run 305 里另外两枚同族正读数（本腿现切）：
- `TestAC4FocusReturnToPriorWindowGap33r5` **PASS**：`panel_host_windows_test.go:944: AC#4 focus hop (head cc31526): ... after Show 0x10204 \| panel hwnd 0x10204 \| ... Hide attempted restore to 0x20196 (SetForegroundWindow 1, SetFocus 131478)` ⇒ **窗建起来了、拿到真 HWND、Show 之后前台归它**。
- `TestTicket255PanelHostBuildsItsWindowOptions`／`TestPanelThreadIsSTAAndExitsCleanly`／`TestAC13BringUp*` 等 **8 枚窗口依赖用例同发 PASS**（`03-*.md` 的表）。

## 4 "根本没运行时"那一支能不能被这些读数排除？

产码里"没运行时"的形状是**有名字的**：`cmd/wisp/panel_host_windows.go:394-396`
`if w == nil { return fmt.Errorf("panel host: WebView2 window creation returned nil (runtime missing or blocked)") }`
测试侧那一支也有名字：`panel_host_windows_test.go:641-642` 的 `t.Fatalf("bringUp on the test thread: %v (WebView2 window creation is available on this box - ...; a nil here means the host, not the environment)")`。
⇒ **CI 那发红在 `:662`，不在 `:642`**；且 `:647 IsCreated()`、`:650 hwnd != 0` 两道 `t.Fatalf` 都没响。
⇒ **两枚红句的形状属于"有运行时、窗也建了，但通道/那一跳没成"，不属于"根本没运行时"**。
更硬的正面证据是 run 294–304 那六发的 `msedgewebview2` 子进程计数（本树 7 == 机器全局 7）与 `same HWND across hide->re-show`／`hot re-show 38.649 ms`——**浏览器子进程真起来了、窗真复用了**。

## 5 仍然取不到的（⛔ 别把 §3/§4 读成"乙也结了"）

- **像素画没画、无人值守下的可见性**：`33-n1` §5.2/§5.3 的 **R-c（BitBlt/PrintWindow 落 PNG ＋非纯色字节数）本仓至今零读数**，本腿也没取到（取它要在 CI 跑一次采集 job ⇒ ⛔ 本腿无权推、无权开窗）。
- **`GetAvailableCoreWebView2BrowserVersionString` ＋两枚注册表 `pv` 键的当场版本读数（R-a）**：**零读数**。本腿只能证明"有窗有子进程"，⛔ 不能报版本号。
- **R-e（同一 job 连跑 10 次的冷起分布）**：**零读数**。今天有的只是单发 `-count=1` 的 3 126–3 941 ms（且那是**预算判据**用的，不是分布）。
- **本机对位读数（R-e 的另一半）**：本机历史是 932–983 ms 那族（`04-*.md`），⛔ 与 CI 那族不可相减。
- **run 305 那一发为什么漂到 `-1`**：本腿**没有**能区分它的读数（同 SHA `cc31526` 的 run 304 是 3414 ms 正读数）。欠的那一发＝同一枚 SHA 上多跑几发、把 E4 那一支的耗时拆开（本腿⛔不能推、⛔不能跑测试）。

rc=0
