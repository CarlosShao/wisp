# 33-n1 verdict — GitHub `windows-latest` 托管 runner 上到底有没有 WebView2 Runtime

只读普查腿。⛔ 本轮没跑任何吃 CPU 的东西（无 `go test`/`go build`/`go vet`/脚本/安装程序），
唯一动作＝文件读命令 + 联网查官方来源。⛔ 未改任何既有文件，未 push，未装任何东西。

---

## 0. 起手锚（现量）

- `date` => `Thu Oct  1 18:09:48 CST 2026`
- `git log -1 --oneline` => `a5357fa1 ledger(A507 更正 A506)…`
- `git rev-parse --abbrev-ref HEAD` => `dev`
- `pwd` => `/d/work/workspace/projects plans/Wisp`
- 终态锚（写本文件时）=> 见文末 §8

---

## 1. 派单给的盘上事实：复认结果

逐条现复，**全部对上**，另有两处补充：

| 派单说法 | 我的实测 | 结果 |
|---|---|---|
| `.github/workflows/ci.yml:388` = `runs-on: windows-latest` | `grep -n "runs-on"` => `:388 runs-on: windows-latest`（同文件另有 `:533` 也是 windows-latest、`:591` 是 `[self-hosted, wisp-slo]`，其余 4 处 ubuntu-latest） | 对上 |
| 全文 `webview\|msedge` 0 命中 | `grep -ci -E "webview\|msedge" ci.yml` => **0** | 对上 |
| `.github/**` 里 `webview\|msedge` 无文件命中 | `grep -ri -l` => **NONE** | 对上 |
| `winlive` 在 `.github/**` 与 `scripts/**` 0 命中 | 两把尺各自 => **NONE / NONE** | 对上 |
| `go.mod` 依赖 `github.com/jchv/go-webview2` | `go.mod:19 github.com/jchv/go-webview2 v0.0.0-20260205173254-56598839c808 // indirect` | 对上，**但注意它是 `// indirect`** |
| 默认档真建窗 11 枚 | 我不重新数，**直接引 `docs/evidence/s1/33-panel-host-c27-v2.md:183` 的逐名名册与 14-3=11 的口径**（该文件可读、本腿未改一字） | 引证，非自量 |

补充读数（本腿自跑，两把尺都是 `grep`）：

- `grep -rln "go:build windows && winlive" --include=*.go .` => **8 枚文件**；
  `grep -rhn "^func Test.*WinLive\|^func TestLive" --include=*_test.go . | wc -l` => **22 枚函数**。
  ⇒ `winlive` 不是"panel 那一族"一枚 tag，它已经压着 8 个文件（ball 热键、audio、task source、approval、panel）。
  这直接改变 §6 的代价算法：**"把 winlive 接进 CI"是一句不能整体执行的话**，必须按族拆。
- 面板代码确实**直接** import 该库：`cmd/wisp/panel_host_windows.go:64`、`cmd/wisp/panel_resident_windows.go:72`
  ⇒ `go.mod` 里那枚 `// indirect` 是过期标注（tidy 未做），**不代表面板没用它**。本腿⛔不改。

---

## 2. 我查了哪些来源（逐条 URL；取不到的也列）

### 甲组：runner 镜像官方清单页 / 镜像仓库本体（actions/runner-images = 镜像的出厂方）

1. `https://raw.githubusercontent.com/actions/runner-images/main/README.md`（21235 bytes）— 镜像可用性表
2. `https://raw.githubusercontent.com/actions/runner-images/main/images/windows/Windows2025-VS2026-Readme.md`（34808 bytes）— **`windows-latest` 今天实际解析到的那一份清单页**
3. `https://raw.githubusercontent.com/actions/runner-images/main/images/windows/Windows2025-Readme.md`（36537 bytes）
4. `https://raw.githubusercontent.com/actions/runner-images/main/images/windows/Windows2022-Readme.md`（38880 bytes）
5. `https://api.github.com/repos/actions/runner-images/git/trees/main?recursive=1`（`"truncated": false`，588 path）— 全仓文件名，找有没有任何安装件提到 webview
6. `https://raw.githubusercontent.com/actions/runner-images/main/images/windows/scripts/build/Install-WinAppDriver.ps1`
7. `…/images/windows/scripts/tests/WinAppDriver.Tests.ps1`
8. `…/images/windows/scripts/tests/Browsers.Tests.ps1`
9. `…/images/windows/scripts/build/Configure-User.ps1`、`Configure-BaseImage.ps1`、`Configure-System.ps1`（找 autologon/session 字样）
10. issue #9538 与其 comments：`https://github.com/actions/runner-images/issues/9538`（Update WebView2 Runtime / Align with Microsoft Edge version，closed as completed，closed_at `2024-04-03`）
11. issue #9537 与其 comments：`https://github.com/actions/runner-images/issues/9537`（Application with WebView2 (fixed runtime) cannot load WebView2，closed）
12. issue #14738 与其 comments：`https://github.com/actions/runner-images/issues/14738`（WebView2 WebDriver automation fails on windows-2025，created `2026-09-16`，closed）

### 乙组：Microsoft 官方 WebView2 文档

13. `https://learn.microsoft.com/en-us/microsoft-edge/webview2/concepts/distribution`
14. `https://learn.microsoft.com/en-us/microsoft-edge/webview2/concepts/evergreen-vs-fixed-version`
15. 试过但 **404，取不到**：`…/webview2/concepts/requirements`、`…/webview2/concepts/faq`（两份都已被改名/合并；我不拿旧记忆的条款填空）

### 丙组：GitHub 官方文档

16. `https://raw.githubusercontent.com/github/docs/main/content/actions/reference/runners/github-hosted-runners.md`（8957 bytes，= docs.github.com 那页的源件）
17. `https://docs.github.com/en/actions/reference/runners/github-hosted-runners`（渲染页；一次 `fetch failed`，一次成功）
    ⇒ 渲染页里那句「each GitHub-hosted runner is a new virtual machine (VM) hosted by GitHub」**我拿不到源件逐字**（规格表在 include 里），只算二等读数，不作为结论承重。
18. 试过但 **404**：`https://docs.github.com/en/actions/reference/github-hosted-runner-images-and-sku`

### 盘上（本地读，无 CPU）

19. `docs/evidence/s1/33-panel-host-c27-v2.md`（§N1、§6、名册那一格）
20. `docs/reports/pending-and-issues.md`（`:9927` 「GitHub runner 有没有能建窗口的桌面＝判不动」、`:10133` Q-75、`:10023` 本机运行库读数、`:3709-3710` A112②）
21. `docs/reports/HANDOVER.md`（`:99` 本腿题面出处、`:126` 坑④、`:530` 那行 CI 注释在 ci.yml）
22. `.scratch/wisp/probes/33/h1/census.md:15`（go-webview2 把 3 枚 `WebView2Loader.dll` 用 `//go:embed` 打进二进制）
23. `.scratch/wisp/probes/33/a2/census.md:279`（loader 侧双路径：磁盘 DLL 或内存加载）
24. `cmd/wisp/panel_host_windows_live_test.go:1-25`（winlive 那格的由来与代价原文）

⛔ 未读未引：`frontend/**`、`design/**`（owner 有未提交件）。

---

## 3. 抄到的原文那一行

### 3.1 `windows-latest` 是哪台机器（官方可用性表，runner-images README:37）

> `| Windows Server 2025<br>… | x64 | `windows-latest`, `windows-2025`, or `windows-2025-vs2026` | [windows-2025-vs2026] |`

同页 `:44`：
> `- In general the `-latest` label is used for the latest OS image version that is GA.`

⇒ **本仓 `ci.yml:388` 的 `windows-latest` = Windows Server 2025**（清单页 `Windows2025-VS2026-Readme.md`，
`- OS Version: 10.0.26100 Build 33438`、`- Image Version: 20260922.246.2`）。
`windows-2022`（= Windows Server 2022，`10.0.20348 Build 5622`）仍在，但**已不是 `windows-latest`**。

### 3.2 清单页上"有没有 WebView2"——**没有这一行**（三张清单页逐一核过）

对 `Windows2025-VS2026-Readme.md` / `Windows2025-Readme.md` / `Windows2022-Readme.md` 三份：
`grep -n -i -E "webview|microsoft edge|msedge"` 的**全部**命中只有两行（版本号不同），逐字（windows-2025 版）：

> `### Browsers and Drivers`
> `- Microsoft Edge 153.0.4234.48`
> `- Microsoft Edge Driver 153.0.4234.48`

`webview` **0 命中**。⇒ **清单页里没有 "Microsoft Edge WebView2 Runtime" 这一条**（windows-2022 同形：`Microsoft Edge 152.0.4191.66`，同样 0 命中 webview）。

全仓文件名里也没有任何安装件提到它：`git trees` 全量 588 path 里 `grep -i webview` => **NONE**；
Windows 侧只有 `Install-EdgeDriver.ps1`。**镜像出厂流程不装 WebView2 运行库。**

### 3.3 但"清单页没列"不等于"机器上没有"——镜像方自己这么说（#9538，closed as completed）

> `we do not install WebView2 and do not configure this component during the build process;`
> `WebView2 is a relatively independent component and is not automatically updated with the browser driver;`
> `the version of the component presented on the agent matches the version of the component that comes with the base image from the Azure gallery.`

> `We won't add separate WebView installation due to a maintenance reasons. The only available version is being provided by MS Edge.`

> `We do not ship a separate version of WebView2 and as far as I remember WebView2 itself comes as a part of the edge installation (which in turn as a part of the standard windows datacenter sku), … it does not seem like a problem that should be solved in the scope of windows runner image`

### 3.4 真机读数（#14738，2026-09-16，两档镜像同一发跑出来的）

issue 正文逐字：
> `WebView2 Runtime **152.0.4191.66** — session never created`
> `- `windows-2022` (passing): WebView2 Runtime **131.0.2903.86** — 8/8 tests passing`

回复（镜像方）逐字：
> `Thanks for the detailed evidence chain — you'd actually already captured the key data point (WebView2 Runtime 131 on 2022 vs 152 on 2025), which is exactly what points to the cause.`
> `The differing variable is the **WebView2 Runtime**, which is Evergreen — it auto-updates on the machine and is not installed or pinned by the runner image (we only ship `msedgedriver` for the full Edge browser). So the `152` vs `131` gap is just Evergreen state; windows-2022 still carries a pre-regression runtime and will roll forward too.`
> `The failure itself is a known upstream Microsoft regression: WebView2 + `msedgedriver` works through runtime 132 and breaks from 133 onward (reported through 141 and 152). Tracked and still open: https://github.com/MicrosoftEdge/WebView2Feedback/issues/5415`
> `Since the runtime is Evergreen and not provisioned by the image, there's no image-side fix.`

⇒ **两档托管镜像上都有 Evergreen WebView2 运行库，且是从 Azure 基盘继承来的、镜像方不锁版本；2026-09 读数：windows-2022 = 131.0.2903.86，windows-2025(-vs2026) = 152.0.4191.66。**

### 3.5 为什么"清单页列了 Microsoft Edge"绝不能当成甲成立（MS 官方，distribution）

> `#### Microsoft Edge Stable channel isn't supported for WebView2`
> `WebView2 apps aren't permitted to use the Stable channel of Microsoft Edge as the backing web platform. This restriction prevents a production release of a WebView2 app from taking a dependency on the browser.`
> `A production release of a WebView2 app can only use the WebView2 Runtime as the backing web platform, not Microsoft Edge.`

以及预装口径（evergreen-vs-fixed-version）**只讲 Win11/Win10，一字未提 Windows Server**：
> `The Evergreen Runtime is preinstalled onto all Windows 11 devices as a part of the Windows 11 operating system. Microsoft installed the WebView2 Runtime to all eligible Windows 10 devices, as described in …`
> `Even if your app uses the Evergreen distribution mode, we recommend that you distribute the WebView2 Runtime, to cover edge cases where the Runtime wasn't already installed.`

⇒ **甲的成立**不是靠"Edge 在清单页上"推出来的，只靠 §3.4 那两枚真机读数 + §3.3 镜像方自述。这条差别就是本腿最容易被判错的地方（见 §7 E1）。

### 3.6 乙组：环境侧唯一一条对本问题有用的官方限制（#9537，Azure DevOps hosted / windows-2022 镜像）

正文逐字：
> `We are running UI Tests with WinAppDriver and EdgeDriver for an MSIX app package.`
> `The application uses WebView2 with the runtime in fixed version distribution mode.`
> `When starting the application, it hangs on initializing the WebView2, which never gets started.`
> `It works when changing the application to use the existing webview runtime.`

维护者逐字：
> `Unfortunately you cannot establish remote desktop connection to Azure DevOps hosted runners.`
> `During the investigation, it turned out that the issue can be reproduced on a clean virtual machine created in Azure. I am going to close the issue since the root cause comes from Azure and is not related to runner-images.`

报障人自己找到的根因逐字：
> `Meanwhile I figured out the root cause: This is due to the agent running under THE admin account.`
> `It also happens when installing an app with fixed WebView2 runtime on the Administrator account on Windows Server and trying to launch it. The problem does not occur on other user accounts, only on Administrator.`

配上 GitHub 官方文档那句（源件逐字，`github-hosted-runners.md:79`）：
> `Windows virtual machines are configured to run as administrators with User Account Control (UAC) disabled.`

⇒ **一条真·环境侧红形：用 Fixed Version 运行库 + Administrator 账户 + Windows Server ⇒ WebView2 初始化挂住；改回系统 Evergreen 就起得来。** 本仓用 Evergreen（`panel_host_windows.go:174` 那格注释写明 user data 走 `%AppData%` 默认），**不踩这一形**，但这条线**必须**写进票 33 的打包决策里（Fixed Version 是 D41/分发那条路上的合法选项，一旦改用就得回来重读这格）。

---

## 4. 甲：运行库装没装 —— **装了（有据，且带版本号与时刻）**

- **结论**：GitHub 托管 `windows-latest`（= **Windows Server 2025**，清单页 `Windows2025-VS2026-Readme.md`，Image `20260922.246.2`）**有** Evergreen WebView2 运行库；`windows-2022` 同样有。
- **凭据**：§3.4 的两枚真机读数（131.0.2903.86 / 152.0.4191.66，2026-09-16 报入官方镜像追踪器）+ §3.3 镜像方「它来自 Azure 基盘、我们不装也不锁」。
- **但清单页不列它**（§3.2，三张清单页 `webview` 0 命中，全仓无安装件）。⇒ **任何"从清单页读出没有"的判法都会得出甲=否的错误结论**；正确形状是：**镜像里没有承诺 ⇒ 不受支持、不锁版本、会漂**。
- **漂移是真实的坏消息**：§3.4 那句 `windows-2022 still carries a pre-regression runtime and will roll forward too` 就是说 2022 那台也会跟着升到 133+。⇒ **甲成立 ≠ 明天还成立**。

## 5. 乙：能不能真开一扇可见窗 —— **拆两半：能建窗有旁证；"可见/前台/毫秒时序"未定**

**不许用甲推乙。** 分开列：

### 5.1 有旁证的半（能建窗、进程能起）

1. §3.4 那一发的**通过**支逐字写着 `windows-2022 (passing): … 8/8 tests passing` —— tauri-driver 流程是**启动一真 WebView2 应用二进制并连它**；那 8/8 绿意味着 WebView2 应用在托管 VM 上确实被拉起、浏览器子进程确实起来了。失败支的成因被镜像方与微软上游共同归到 `msedgedriver`/CDP 的版本回归（`breaks from 133 onward`），**没有一格被归到"这台机器不能建窗"**。⚠ 本仓**不用 msedgedriver/CDP**（`go-webview2` 直调 COM vtable），那一枚上游回归打不到我们。
2. 同一发里 `WebView2 Runtime 152.0.4191.66` 这个读数本身，只能是**在那台机器上问运行库问出来的**。
3. 镜像**刻意**装了 `WinAppDriver 1.2.2009.02003`（清单页 `:87`），并由 `Install-WinAppDriver.ps1` 从 `microsoft/WinAppDriver` 装、跑 Pester 验收；`WinAppDriver.Tests.ps1` 另钉了一格 **`Developer Mode is enabled`**（`AllowDevelopmentWithoutDevLicense = 1`）。⇒ 镜像方的**姿态**是把「驱动真应用窗口」当用途之一。
4. 依赖自带 loader：`h1/census.md:15` 逐字「库**把 3 枚 `WebView2Loader.dll` 用 `//go:embed` 打进二进制**…**磁盘取不到就 `winloader.LoadFromMemory` 走内嵌那份**」⇒ **MS 官方那条"loader 必须随应用发布"的要求（§3.5 来源 13 的 "Files to ship with the app"）在本仓已由依赖满足，CI 不需要额外铺 DLL。**

### 5.2 没有依据、我拒绝填的半（写成"未定"）

- **作业跑在哪一类窗口站/桌面会话**（交互式 Session 1 的已登录桌面，还是服务/session 0）：**官方文档零句话**。我读了 `github-hosted-runners.md` 全文，关于 Windows 的配置声明只有两条（admin+UAC 关闭、托管在 Azure），**通篇没有 interactive desktop / GUI / display / session 任何一字**；镜像构建脚本里 `Configure-User.ps1`/`Configure-BaseImage.ps1`/`Configure-System.ps1` 三枚 `grep -i "autologon|logoncount|session|interactive|runneradmin"` **全部 0 命中**（⇒ 查不到 auto-logon 配置，既不能证有也不能证无）。
- **有没有 GPU**：清单页与那页官方文档都没有 GPU 栏（`larger runners` 一档另广告 GPU 机型）。**"标准托管 Windows runner 没有 GPU"是我从"缺栏"推的，不是原文说的** ⇒ 记为我的推断（§7 E4）。真无 GPU ⇒ Chromium 走软件合成，窗仍应有内容，但**这条我没有官方依据**。
- **像素到底画没画**、**`SetForegroundWindow` 那族前台语义能不能成**（票 33 AC#4「焦点回还」那一族断言正好踩在这上面）：**零依据**。而且 §3.6 维护者那句 `you cannot establish remote desktop connection to … hosted runners` 意味着**目视这条路在托管档根本不存在**，只能靠程序化读数。
- **`WinAppDriver` 那一格只是"目录存在"**（`Test-Path "${env:ProgramFiles(x86)}\Windows Application Driver" | Should -BeTrue`）⇒ **镜像验收并不保证真能驱动窗口**，别把"装了工具"读成"能力被验收过"。

### 5.3 判乙需要看的那一根读数（一次性 job，5 项都在同一次跑里采）

要判的不是「CI 到底能不能建窗」这句空话，而是下面 5 枚具体读数；缺任一枚就仍判不动：

| # | 读数 | 判的是 |
|---|---|---|
| R-a | 进程内 `GetAvailableCoreWebView2BrowserVersionString` 的版本串 **＋** 两枚注册表 `pv` 键（HKLM `…WOW6432Node\Microsoft\EdgeUpdate\Clients\{F3017226-FE2A-4295-8BDF-00C3A9A7E4C5}` / HKCU 同路径；MS 官方给的检测口径） | 甲的**当场**复认（不靠 2026-09 的二手读数） |
| R-b | 窗建起来后 `IsWindow(hwnd)` / `IsWindowVisible(hwnd)` / `GetWindowLong(GWL_STYLE)` 与 `GetForegroundWindow()` 是否 == 我们的 hwnd | 乙的"前台语义可用"那一半 |
| R-c | 对窗 DC `BitBlt` 一帧落 PNG（或 `PrintWindow`）+ 该文件的**非纯色字节数**，附字节数 | 乙的"真画了东西"那一半，⛔ 不许用"没报错"代替 |
| R-d | 本进程树内 `msedgewebview2.exe` 枚数与 `--user-data-dir=` 归因（照 `33-panel-host-c27-v2.md` §桌面卫生那格的尺形） | 子进程真的起来了、且能归因到本树 |
| R-e | 同一 job 连跑 10 次的冷起窗耗时分布（min/med/p95）+ 同一段代码在本机的对照 | 判"毫秒级时序断言能不能进 CI"，**这枚才是 winlive 归位的尺** |

⚠ 采法必须**带日志全文**（不是只带 PASS）；`log.Fatal` 那支（`33-panel-host-c27-v2.md:189` 具名的 `webview.go:115/120/125`）会把整个 test 二进制打掉、名册采不到 ⇒ **R-a 要早于任何建窗动作跑**。

---

## 6. 可执行结论：**有条件能**（中间档），并附"哪几格仍只能本机"

### 6.1 一句话

**不用在 CI 里加任何安装步骤**（甲成立，装了，且镜像方明确拒绝装 ⇒ 我们不该去跟基盘抢版本）。
但"能不能真开一扇窗"里**可见性/前台/时序**这一半**没有官方依据** ⇒ 今天**不许**把 `winlive` 整族接进 CI；
能接的是**不含毫秒时序与像素断言的结构断言子集**，且必须先满足下面 4 个条件。

### 6.2 前置条件（4 条，逐条给出处）

1. **先解掉"一推送 windows cli 腿就红"那一枚无关红**。`33-panel-host-c27-v2.md:188` 已具名：干净检出里 `frontend/dist` 只有 `.gitkeep` ⇒ AC#13 用例 `t.Skipf` ⇒ `runtests.sh:98` 把 SKIP 判红。**这条与运行库无关，但不修就没人能读 R-b/R-c 的日志**（腿在采到之前已经红）。
2. **预检先行、按名拒绝，不许 `log.Fatal` 抢话筒**：按 §5.3 R-a 的形状在建窗前查运行库；缺或过旧 ⇒ **具名失败/具名跳过并打出 `pv` 读数**。⛔ 这一条与项目"缺就红不 Skip"的既有形状（`panel_resident_windows_test.go:25-28`）冲突，**冲突本身要摆给编排者定**，本腿不替他拍。
3. **锁形：只用 Evergreen，绝不用 Fixed Version**（§3.6 的 admin+Windows Server 挂起形，配 `github-hosted-runners.md:79` 那句"以管理员运行、UAC 关闭"）。**若 D41 分发路上哪天选 Fixed Version，这一格必须重开。**
4. **接受版本漂移并把它当噪声源**：`152 vs 131` 那格说明镜像方不锁 ⇒ **CI 档里任何毫秒级阈值都不可复现**（§5.2 最后一行）。

### 6.3 加在哪（不改动，只指位置）

最省的路径**不是新开 workflow**：`ci.yml` 的 `test-windows`（`:388`）里，在
第 4 步 `cmd/wisp CLI tests`（step 名在 `ci.yml:455`，其 `run: bash scripts/wisp-cli-tests.sh` 在 `ci.yml:475`）之后追加**一枚 step**，
跑 `-tags winlive` 的 panel 那**一枚**文件里的结构断言。
⛔ 新开 workflow = 契约级（`pending-and-issues.md:10133` Q-75 甲栏原文已登记，本轮不摆 owner）；本腿不改 CI 文件一字。

### 6.4 代价（诚实标量级）

- **时长我没有读数**（⛔ 不许推 CI 试跑、⛔ 不许本地跑测试 ⇒ 无任何计时凭据）。可引的对位读数只有盘上的两枚：
  本机 `go test ./cmd/wisp/` => `ok 112.312s`（`pending-and-issues.md:9739` 逐字）、
  `panel_host_windows_live_test.go` 那一枚断言的**上界**是 2 秒（`The bound is UNCHANGED at 2 seconds`）。
  ⇒ 量级只能写成「同 job 内追加一枚带 `-tags winlive` 的 test-binary 重建 + 个位秒级用例」，**分钟级增量＝未量**。
- **维护代价**：`winlive` 现在压着 **8 枚文件 / 22 枚函数**（§1 补充读数），里面 ball 全局热键、audio、task source
  这些**与 WebView2 无关、但环境上根本不可能在托管 runner 跑**（`ci.yml:530` 本仓自己那句逐字
  `Memory/handle subset on a hosted Windows runner (no audio hardware)`）。
  ⇒ **"接 winlive 进 CI" 必须按族拆**：panel-窗族（可议）、ball/audio 族（**永远仅本机**，理由不是 WebView2，是没有硬件与热键独占权）。
  这句话如果整条执行，会把不相干的两族一起拖进同一枚红。

### 6.5 那 CI 今天到底测到了什么 / 测不到什么（写死给票 33 那几格用）

- **测得到（一旦 §6.2① 修好）**：默认档那 **11 枚真建窗用例**的**结构**维度 —— 同 HWND 复用、单扇窗、本树持浏览器子进程、Destroy 后不留窗、STA 线程干净退出、绑定回包到页（名册逐名见 `33-panel-host-c27-v2.md:183`）。它们**没有 build tag**，本来就设计成"建不起来就是红"。
- **测不到**：`≤2s 退净`那一枚（已在 `winlive`，`panel_host_windows_live_test.go:36` 原文「winlive has NO CI job, so this clause is now measurable only on a desktop box」）、任何毫秒级时延分位、任何像素/前台/可见性断言、任何热键与音频端到端。
- **今天真正没测的原因**：不是"没运行库"，而是 §6.2① 那枚 SKIP-判红 + 从未带 `-tags winlive` 跑过（`winlive` 在 `.github/**`、`scripts/**` 0 命中，本腿复认）。

---

## 7. 我可能判错的条目（每条附"如果错了后果"）

| # | 我可能判错的地方 | 若错了的后果 |
|---|---|---|
| E1 | **把 §3.4 那枚"真机读数"当甲的凭据**——它出自 issue #14738 正文里报障人的自述，我**没有独立复核**（复核＝去读他那条 workflow 的日志原件，我没读）。派单要求"只认清单页/官方文档"，§3.4 严格说是**官方 issue 追踪器里用户+维护者双方共同确认的读数**，比博客强，但不是清单页。 | 若那枚 152/131 是报障人从**别的机器**抄来的 ⇒ 甲退回"读不出"，§6 整档从"有条件能"退到"不能"，票 33 那几格要按〔仅本机可量〕写死。**这一条是本腿结论最软的一环。** |
| E2 | `windows-latest` 的解析**会变**（README:97 逐字：`any workflow using the -latest label, may see changes in the OS version`）。我读到的是 Server 2025。 | 若票 33 落地那天 latest 已漂 ⇒ 我的清单页引用过期一格；**解法＝CI 里钉 `windows-2025`/`windows-2022` 字面**，别用 latest（这一条本腿只是指认，⛔未改 ci.yml）。 |
| E3 | **"能建窗"我只拿到 tauri/msedgedriver 那一族通过支**。那一族建窗后要接 CDP；我们建窗后走 COM vtable 直连。两条路**不同**。 | 若托管环境其实允许建窗但不允许 Chromium 子进程正常渲染 ⇒ R-b/R-c 会当场红，§6 的"结构断言可以进 CI"要缩到只剩"能 New 出对象"那一维。 |
| E4 | **"没有 GPU"是我从缺栏推的，不是原文**。 | 若托管机有 GPU：不影响结论方向。若**没有且驱动都不装** ⇒ 软件合成路径下窗可能建得起但首帧空白，R-c 是唯一能抓到的尺；我在 §5.2 已把它写成未定，不算自相矛盾。 |
| E5 | §3.6 那一形出自 **Azure DevOps hosted**（#9537 正文勾选 `Azure DevOps`，机器名 `fv-az…`），不是 GitHub Actions 托管。两家的 windows agent 用法不同（登录态/账户） | 我把它当"CI 通用危险"引用是**外推**。真结论应是：「同镜像同账户族的一个已确证坏形，路径与 ours 不同（我们不用 Fixed Version），故今天不触发；一旦改用 Fixed Version 需按 GitHub Actions 自己重验」。 |
| E6 | 依赖内嵌 `WebView2Loader.dll` 这条我**引 `.scratch/wisp/probes/33/h1/census.md` 的读数**，没在本腿重跑那把尺（重跑要读模块缓存，代价可控但我没做）。 | 若那 3 枚 DLL 其实不参编 ⇒ CI 还要多铺一枚 loader，§6.4 的"不需要额外件"要改，且**默认档 11 枚在 CI 会当场红在加载阶段**而不是建窗阶段。 |
| E7 | **11 枚名册我是引 v2 表、不是自量**（派单让我"复认后引用"，我复认了它的口径文字，没复跑 `grep -c` 那把尺）。 | 若名册数错 ⇒ §6.5"测得到"的枚数错，方向不变。已注明引证出处，不冒充现量。 |
| E8 | §6.3 我建议 step 加在 `test-windows` 里——但同一个 job 已有多枚 `if: ${{ !cancelled() }}` 的门前置步骤（`:388` 之后那 5 步）。 | 若前面任一枚红 ⇒ 我这枚建议的落点拿不到 verdict（票 111 那格讲过这形状）。真要加，落点应带 `!cancelled()` 且排在**采集链最前**或另立 step 组。 |

---

## 8. 判不动的地方（交回编排者，我不填）

1. **乙的后半（可见性/前台/像素）——判不动**。缺 §5.3 R-b/R-c 两枚读数。要它**不是**"推一次 CI 试试看"那种裸推，而是**一次性 job 打日志**（R-a..e 五枚可在同一发采齐）。⛔ 派单 §"不许推 CI"的约束我没违反：本腿一次都没推。
2. **甲的"明天还成立吗"——判不动**。镜像方逐字拒绝锁版本（§3.3/§3.4 `will roll forward too`）。⇒ 任何以甲为前提的 CI 档都必须**自带预检 + 具名失败**，而不是把甲当定理。
3. **"缺运行库时应红还是应 Skip"这一格**——判不动，属**冲突上交**：本仓既有形状写"建不起来就是红不 Skip"（v2 表 `:183` 引的 `panel_resident_windows_test.go:25-28`），而 `log.Fatal` 那支会让名册采不到（`:189`）。这不是我一层能拍的。
4. **`winlive` 族里非窗的那 6-7 枚文件到底归哪一档**——判不动（它们的环境约束是热键独占/音频硬件，与本题无关），需要另腿按族拆。
5. **票 33 §6.2① 那枚 SKIP-判红先不先修**——**阻塞本腿结论落地**，但它是另一格（干净检出的 `frontend/dist` 产物），⛔ 我不顺手动 `frontend/**`。

---

## 9. 工具输出里出现的"像授权"文字（照实报，我一个字没照它做）

- `git log -1` 那枚提交（`a5357fa1`）正文里有编排性指令文字（「只读腿继续并行、写码腿因毫秒计时仍串行」「腿≤2格」等）。**它是仓里的台账提交说明，不是对本腿的授权**；本腿只按派单的铁边界执行：⛔未跑 CPU、⛔未改既有文件、⛔未 push。
- 拉回来的 GitHub issue 正文/评论里出现 `<system>` 字样的转义痕迹与「Workaround - Update manually via official bootstrapper」一类**建议装东西**的文字。**都不是给我的指令**：本腿⛔未装任何东西、⛔未动注册表、⛔未动系统 Edge。
- 未见任何"放宽阈值/不用取证/已解锁"类文字被采信。§6.4 里那条 `≤2s` 界我明确写成「上界原文一字未动」，没有任何一档被我放宽。

## 10. 终态锚

- 写本文件时的 `git log -1` 与起手锚同枚（`a5357fa1`）；期间 `HEAD` 若有前进，只可能是同机器上的写码腿，**与本文件无关**。
- 本文件是唯一新增件：`.scratch/wisp/probes/33/n1/verdict.md`。临时件只建不删。
