# bundle-1 普查件 — 出货链上"那份前端产物由谁生成"

leg id = `bundle-1`｜只读普查｜性质：**证据件，不是判据、不是开票授权**（照 `AGENTS.md` §1.5 落 `probes/`）
本件⛔ 不改一行产码、⛔ 不新建工单、⛔ 不翻任何框；差集与缺陷形状**交编排者裁**。

---

## §0 起手锚（三条具名声明先落）

| 项 | 读数 | 尺（我现跑） |
|---|---|---|
| HEAD | `c3183379529dd7498c9b8d65d2fe58cc22c20951` | `git rev-parse HEAD` |
| 工作树脏度 | `git status --porcelain \| wc -l` ＝ **758 行**（起手即如此，非我造成） | 现跑 |
| **我没跑 `go test`／`go build`／突变** | 全件结论只来自读码＋grep＋`ls`/`wc`/`find` 计数＋`git show`／`git ls-files`。⛔ 一枚编译或测试命令都没发（同波写腿要跑整包红名册） | 本件即凭据 |
| **我读 `ci.yml` 走的是不可变快照** | `.github/workflows/ci.yml` 一律 `git show HEAD:.github/workflows/ci.yml` ＝ **910 行**，落 `logs/ci-at-HEAD.yml` 后再 grep。⛔ 没读工作副本那一份（此刻有别的腿在编辑它）。**本件里所有 ci.yml 行号都对 `HEAD` 那一份成立，不保证对工作副本成立** | `wc -l logs/ci-at-HEAD.yml` |
| **`ci.yml` 快照的锚点会漂：起手 HEAD 之后 tip 前进了** | 我起手记的 HEAD＝`c3183379…`，我的快照＝`git show c3183379:.github/workflows/ci.yml`（910 行）。**跑到本件收尾时 tip 已前进到 `f6b79ab047b419403d7d26fcfc121289647c6e02`**（同波那枚 `ci(111 r5)` 腿落的），而 `git diff -U0 c3183379 HEAD -- .github/workflows/ci.yml` ＝ **`@@ -203,0 +204,31 @@`＝纯插入 31 行**。⇒ **本件里所有 `ci-at-HEAD.yml:NNN` 行号对 `c3183379` 那一份成立**；要落在现 tip 上请**逐枚 +31**（`:540`→`:571`、`:726`→`:757`、`:732`→`:763`、`:789`→`:820`、`:795`→`:826`、`:840`→`:871`、`:848`→`:879`、`:871-872`→`:902-903`）。新文件 941 行。⛔ 我没有读工作副本，也没重取快照 | `git rev-parse HEAD`＋`git diff --stat c3183379 HEAD -- .github/workflows/ci.yml`（＝`1 file changed, 31 insertions(+)`） |
| 边界自证 | ⛔ 我没打开任何 `frontend/src/**`、没读 `design/**`。我读了 `frontend/embed.go`（**一枚 Go 文件、不在 `src/` 下**）——依据＝派单正文"以及 Go 侧内嵌资源的代码"在允许清单里，而禁句具名是 `frontend/src/**`；这层解释**由编排者追认**（见 §5 第 3 条） | — |
| 计数尺的正控（先证尺命中得了真名再判"不存在"） | 同一把尺在 CI 快照上命中 `setup-node@v4`(`:848`)、`npm run build`(`:872`)、`actions/upload-artifact@v4`(`:732`/`:795`) ⇒ 这些词面**扫得到**；于是下文"没有某一步"的负向句才有牙 | `logs/ci-hits.txt`＋`logs/ci-steps.txt` |

logs 全在 `.scratch/wisp/probes/bundle/1/logs/`（原始输出落盘，本件只引摘要）。

---

## §1 逐问结论＋尺

### A. 出货链里到底有没有"生成前端产物"这一步

**答复：ⓐ本机打包＝明确跳过（有字面承认）；ⓑCI＝有一枚 job 会生成 `frontend/dist`，但那枚 job 根本不产 `wisp.exe`，而产 `wisp.exe` 的三枚 job 一次 node 都不装。两条链今天都产不出"带包体的 exe"。**

**ⓐ 本机打包链**（`scripts/build.ps1`，SPEC-11 §2.2 的实现）：

- `scripts/build.ps1:10`（`.DESCRIPTION` 逐字）：`    2. frontend: skipped (none yet, lands in S5)`
- `scripts/build.ps1:73-74`（执行体逐字；`:74` 是那一发 `Write-Host`，**整步只有一句打印**）：
  - `:73` `# --- 2. frontend ------------------------------------------------------------`
  - `:74` `Write-Host 'build.ps1: frontend step skipped (no frontend yet; embed lands in S5 per SPEC-11 §2.2)'`
- 紧接着 `:129` `    & $go.Source build -trimpath -ldflags $ldflags -o (Join-Path $outDir 'wisp.exe') ./cmd/wisp` ⇒ **`go build` 在"前端那步只打印一句话"之后直接跑**，embed 拿到的就是仓库里那枚锚文件。
- 尺的完整维度：`build.ps1` 全文只 3 行命中 `frontend|npm|node|dist|vite|bundle|embed`（`logs/build-ps1-hits.txt`＝3 行）；参数面只有 `:24 [string]$Env = 'dev'`，**SPEC-11 §2.2 写的 `--with-frontend` 旗标在这份脚本里不存在**（`grep -nE 'with-frontend' scripts/build.ps1` ＝ 零命中；正控＝同一把尺在同文件命中 `go build`/`wisp.exe`）。
- 旁证（SPEC-11 侧那两条相关落点也是空的）：`docs/specs/SPEC-11-build-deploy-containerization.md:44` 逐字 `  2. 前端产物：存在则跳过；--with-frontend 时走 docker/frontend.Dockerfile（§3.2）`；而 `docker/` 目录现量只有 `builder.Dockerfile · compose.dev.yml · compose.test.yml · mockllm.Dockerfile · model-mirror`，**`docker/frontend.Dockerfile` 不在盘上**；`grep -nEi 'frontend|npm|dist|assets/web|embed' docker/builder.Dockerfile` ＝ **0 行命中**（`logs/builder-docker-hits.txt`）。

**ⓑ CI 构建**（全部行号对 `HEAD` 快照，见 §0）：

- 六枚 job：`lint`(`:65`) · `test-core`(`:364`) · `test-windows`(`:474`) · `slo-smoke`(`:708`) · `slo-full`(`:766`) · `lint-frontend`(`:840`)（尺＝`grep -nE '^  [a-z0-9-]+:'`）。
- **会跑 `npm run build` 的只有一枚**：`ci-at-HEAD.yml:871-872` 逐字
  - `:871` `      - name: build (vite build -> frontend/dist, the bytes go:embed carries)`
  - `:872` `        run: npm run build`
  它在 `lint-frontend` job 里，该 job 的 `working-directory: frontend`(`:844`)、`actions/setup-node@v4`(`:848`)、`npm ci`(`:857`)。**这一枚 job 全程零 `go build`**（`setup-node` 在全文件只命中一次＝`:848`；尺＝`logs/ci-steps.txt` 59 行里 node 相关全在 840 之后）。
- **真正产出 `wisp.exe` 的是三枚、全都不装 node**：`:540`（`test-windows`／步名 `:535` "cgo build smoke (build.ps1 fetch-deps + mingw link + doctor)"）、`:726`（`slo-smoke`／步名 `:725` "Build wisp.exe"）、`:789`（`slo-full`／步名 `:788`）三条逐字同一句
  `        run: powershell -NoProfile -ExecutionPolicy Bypass -File scripts/build.ps1 -Env dev`
  ⇒ CI 里那枚"会生成 dist"的 job 与那三枚"产 exe"的 job **不在同一枚 job、不在同一台工作区**，`frontend/dist` 在 CI 里从来不会进入 `go build` 的 embed。
- **产物有没有以 artifact 形式被搬过去**：没有。全文件 `upload-artifact` 只两枚＝`:732`、`:795`，两处步名都是 `Upload SLO report`（`:731`、`:794`），**传的是 SLO 报表**，不是 `frontend/dist`，也没有任何 `download-artifact` 消费方。
- 同形第二处（文档侧，⛔ 我不改它）：`docs/BUILD.md:51` 逐字 `2. 前端：跳过（S5 之前无 frontend）`——票 244 §5 的 G6（`244-….md:44`）已把这一行列为**过期文档行**并裁"票 244 交件里逐枚具名登记、`docs/BUILD.md` 一字不改"。
- 结论一句：**CI 今天没有任何一条链能造出"页面字节真在里面的 wisp.exe"**；`lint-frontend` 造出的 dist 随 job 工作区蒸发。

### B. Go 侧怎么拿到那份产物

**答复：`//go:embed`，编译期烘进 exe；路径来源＝Go 包目录的相对子目录 `dist`，也就是 `frontend/dist`。⛔ 没有运行时读磁盘、没有本地 HTTP、没有 CDN。**

代码锚（逐字原文，均为整行读）：

- `frontend/embed.go:19` `//go:embed all:dist` ＋ `:20` `var distFS embed.FS` ＋ `:23-25` `func Dist() embed.FS { return distFS }`。
  ⇒ **路径来源＝那枚模式串本身**：`go:embed` 的模式相对**声明它的 Go 包目录**解析，`frontend` 包在 `frontend/`，所以 `all:dist` 打的就是 `frontend/dist`。`all:` 前缀连点文件一起收（票 33 AC#12 正是拿这一点说"锚文件也会被算进去"）。
- `frontend/embed.go:4-7`（逐字，说明"唯一通路"）：`// go:embed below is the ONLY way those bytes reach wisp.exe. There is no` / `// runtime CDN and no local HTTP server (D29 - the WebView2 host answers` / `// AddWebResourceRequested from these bytes), so the binary must be` / `// self-sufficient on a machine with no node, no npm and no network.`
- 消费方：`internal/panel/assets.go:44` `	sub, err := fs.Sub(frontend.Dist(), "dist")`（函数 `BuiltinAssets()`，`:43`）；入口判据 `:56-58`＝**只看 `index.html` 在不在**（`const EntryFile = "index.html"` ＝ `:29`）；`:65` `func (a *Assets) Built() bool { return a != nil && a.built }`。
- **产码调用点枚数**（不是 rc=0、不是"包能编译"）：`panel.BuiltinAssets()` 全仓非测试调用点＝**2 枚**——`cmd/wisp/panel_resident_windows.go:238`（常驻腿真那条）＋`cmd/wisp/panel_assets.go:88`（`wisp panel-assets` 诊断命令）。其余 6 处是 `_test.go` 与 `.scratch/**` 里的旧拷贝（尺＝`grep -n 'BuiltinAssets()' *.go` 族，落 `logs/`；正控＝同一把尺命中 `assets.go:43` 的声明行）。
- ⛔ 与"镜像站哈希来源／明文密钥"那一族禁令**无关**，本问没碰那一条，也不该混（派单 §1.2 禁止清单里的 L4 那条讲的是取哈希，不是取产物）。

### C. 产物缺失时程序今天怎么表现

**答复：不是"开一个空白窗、什么都不说"。今天的形状＝fail-closed 三段：① 宿主显式打一句人看的话（`SetHtml` 一枚自造通知页）；② 该句之下的错误具名"怎么修"；③ CLI 侧 `wisp panel-assets` 以 rc=1 说 `assets NOT BUILT`。⚠ 但 `wisp doctor` 这一族今天完全不查这一维。**

- 拒绝渲染锚文件：`internal/panel/assets.go:34` 逐字
  `var errNotBuilt = errors.New("panel: embedded assets are not built (run npm run build in frontend/)")`
  调用面 `:76-78`（`Resolve` 进门就 `if !a.Built() { return nil, "", errNotBuilt }`）与 `:101-103`（`Manifest()` 同样先拒）——`:31-33` 注释逐字要求调用方**必须**把这句显示出来而不是渲染锚文件。
- 宿主那一发（**产码路径、非 mock**）：`cmd/wisp/panel_host_windows.go:428-430`
  `	if err := m.serveEntry(); err != nil {` / `		m.serveNotBuiltNoticeLocked()` / `	}`
  → `:442-451` `serveNotBuiltNoticeLocked` 用 `w.SetHtml(...)` 写出一页：`:449-450` 逐字含
  `<title>panel assets unavailable</title>`＋`<p>panel assets unavailable: the embedded bundle is not built.</p>`
  底层拒绝句在 `:456` `		return fmt.Errorf("panel host: embedded bundle not built (only the tracked placeholder is present)")`。
  ⇒ **调用点现量**：`serveNotBuiltNoticeLocked` 全仓被调 **1 次**＝`panel_host_windows.go:429`；`serveEntry()` 被调 **1 次**＝同文件 `:428`（尺＝`grep -n 'serveEntry(\|serveNotBuiltNoticeLocked'` 全 `*.go`；正控＝同一把尺命中 `:442`/`:454` 两枚声明行，命中得了真名）。
- ⚠ 顺带把编排者转述里的一格纠掉（详见 §5 第 1 条）：**票 33 AC#13 那发"探测页盖掉真页面"的形状今天已经不在**——`:420-426` 的注释逐字写着 `AC#13's order, and it is the whole fix: prove the message channel first, / hand the page over LAST. ... running it after serveEntry meant every cold start finished on the probe page instead of the panel (ticket 33 AC#13, 33-v1 §AC#3 "供给那半")`，而现读时序确实是 `:426 firstRoundTripLocked` 在前、`:428 serveEntry` 在后。
- **`wisp doctor` 那一族**：`cmd/wisp/doctor.go`＝339 行，`grep -niE 'panel|asset|embed|dist|built|bundle|frontend'` 只 **2 行命中**（`:42` 打印 version/commit/built-at 那枚 `built=%s`，与 `:81` sherpa-onnx 对 onnxruntime 的版本），**零** 与前端包体相关。⇒ **开机自检里没有"页面产物在不在"这一维**；而 `PLAN.md` D41（`docs/PLAN.md:2900-2901`）把"启动自检（任一失败给明确错误，**不得静默降级**）"的清单写死为 `onnxruntime.dll` 版本 → 自身签名 → C29 公钥 → 数据目录可写，**这四项里没有前端产物那一枚**（原文未列＝**规格没覆盖**，不是我推断它要求）。
- 运维口径的诊断命令（有、可用）：`cmd/wisp/panel_assets.go:126-128` 逐字
  `		if !assets.Built() {` / `			fmt.Fprintln(os.Stderr, "wisp panel-assets: assets NOT BUILT")` / `			return 1`
  另有 `:114-124` 的 `-check` 支（专抓"go:embed 比 npm build 少装了东西"那一形：入口在、它引用的 hash 资产 404 ⇒ `Built()` 仍 true 而页面空白），与 `internal/panel/assets.go:137-160 Check()`。

### D. 查重（现有哪些票管到哪一格）

尺：`.scratch/wisp/issues/` 现量 **261 枚**文件；票池 grep 全部**落 `logs/`** 再看摘要，⛔ 全仓 `grep -rn`。词面族按"这句话在说什么"铺开（照 `A473` 那条教训）：`frontend/dist` · `go:embed` · `npm run build` · `Built()` · `not built` · `BuiltinAssets` · `build.ps1` · `assets/web` · `with-frontend` · `frontend.Dockerfile` · `windowsgui`。

命中枚数（去掉过滤前的原始末几行都落盘了）：`issue-body-hits.txt`＝**18 枚**票正文命中 embed/bundle 族；`buildps1-tickets.txt`＝**16 枚**提到 `build.ps1`，其中**具名"frontend step／前端"这一格的＝0 枚**（唯一非零是 `268-…`，而它那一行讲的是 stdout 钉数与禁改清单，不是打包步骤——已整行读）。

| 那一格 | 已有票（到哪一行） | 状态 |
|---|---|---|
| embed 机制本身＋`Built()`＋"无 node 机器自给自足"这句要求 | **票 77**：`77-…shadcn.md:53`（骨架＋`go:embed` 打进 `wisp.exe`）、`:74`（CI 新 job `lint-frontend`＝`npm ci && typecheck && lint && build`）、`:80` `- [ ] **AC#1** frontend/ 可构建，产物被 go:embed 打进二进制，wisp.exe 在无 node 环境的机器上能拉起面板`，判据（`:81`）**剥光 PATH 里的 node/npm 后真拉起一次** | AC#1 **未勾**；AC#2 已勾。文件名无 `-done` |
| "能构建≠有页面可发"＋要一枚能区分**只匹配到占位文件**与**真有一包产物**的仪器 | **票 33 AC#12**（`33-panel-host-c27.md:37` 整行读）＋**票 248 AC#9**（`248-….md:57` 整行读）。两枚都逐字写着 `⚠ 产物由界面侧那枚 agent 产出，本编队⛔ 不写 frontend/** ⇒ 这一格在"谁把 dist 填上"落定前勾不了`，且 248 那句 `具名交 owner 带话（要带的话在 §5）` | 两枚 AC **都未勾**（同一枚归口的两面） |
| 探测页盖掉真页面（冷启动最终显示空壳页） | **票 33 AC#13**（`:40` 整行读，判据①次序重排②会响的断言③反控调换次序必须红） | **产码已落地**（`panel_host_windows.go:420-430` 现读＋注释自认 `AC#13's order, and it is the whole fix`）；⛔ 我没翻框，勾不勾归非实现者 |
| 宿主怎么取字节（`AddWebResourceRequested` 从 embed 答） | **票 33**（`dist:2 embed:2 ci:4`，`33-panel-host-c27.md:56` 仍写着**旧路径** `//go:embed assets/web/dist`） | 票面路径与活码不一致（见 §2 第 3 行） |
| 设置页要真能点出来 ⇒ 前提也是"dist 真有一包产物" | **票 248**（`dist:3 embed:2`） | 已登记，同上被"谁填 dist"卡住 |
| `-H=windowsgui`／PE 子系统／双击黑框这一族 | **票 244**（open，`244-….md` AC#0–AC#6＋§5 编排者收件七格；`:11` 引 `SPEC-11…:50`；`:16` 具名 `scripts/build.ps1:103-110` 的 `$ldflags` 六枚全 `-X`、零枚 `-H`） | **相邻但不是本问**：它管"exe 的子系统"，不管"exe 里有没有页面字节"。它 AC#1 那格还写着**⛔ 票 244 未答 AC#5 ⓐ 之前整票不许动构建链**（`:29` 逐字）⇒ **构建链那一面今天同时被 244 的门禁着** |
| 打包/分发/签名（S8） | **票 56** `56-s8-signing-distribution-naming.md:1` 标题即 `DEFERRED`，全文 `frontend/embed/bundle/dist` 命中 **0**（只在 `:11` 讲 signing/distribution） | 不管这一格 |
| **出货那一刻，谁执行 `npm run build` 把它放进 `go build` 之前** | **名义射程＝票 34**：`34-frontend-scaffold.md:32 Status: ready-for-agent`、`:42` `- [ ] \`docker build\` → dist → \`go build\` serves the scaffold in the real panel host (33).`、`:28` 讲 `assets/web/dist`、`:6` `Blocked by: 01-build-chain (assets/web embed target exists)` | ⚠ **这一格有票但票的落点是过期名册**：它写 `assets/web/dist`＋docker builder，而活码是 `frontend/dist`＋`//go:embed all:dist`；`01-build-chain` 已 `-done` 且**没带走 assets/web 那枚 embed 目标**（`assets/` 目录现量**不存在**） |
| **CI 产出的那枚 exe 必须真带包体（embed 条目枚数 > 1）这一维** | **零票**：`grep -rl 'with-frontend' .scratch/wisp/issues`＝**0**；`grep -rl 'frontend.Dockerfile'`＝**0**；没有任何票要求 `build.ps1` 的前端那步不许只是 `Write-Host`，也没有任何票要求 CI 产 exe 的 job 装 node 或消费 dist | **无票**（见 §2 末行） |

### E. 现状一致性

| 项 | 读数 | 尺／出处 |
|---|---|---|
| `git ls-files frontend/dist` | **1 枚**＝`frontend/dist/.gitkeep`（与编排者转述一致） | `logs/git-ls-files-frontend-dist.txt` |
| 本机另有一份未跟踪产物 | `find frontend/dist -type f \| wc -l` ＝ **4**（含 `.gitkeep`；顶层可见 `index.html`＋`assets/`）⇒ 与台账里那枚"三枚未跟踪构建产物"相比**多了 1 枚**。**我只数枚数，⛔ 没读它们一个字节** | 现跑；台账旁证 `docs/reports/pending-and-issues.md:1567` 逐字 `差值 +3 已证是**三枚未跟踪的 \`frontend/dist/\` 本地产物**（\`git check-ignore\`）`、`:4427` 同形再述 |
| `.gitignore` 那条规则的**出处** | 规则本体 `.gitignore:22-23` `frontend/dist/*` ＋ `!frontend/dist/.gitkeep`；**理由写在 `:18-21`**，逐字：`# dist/ is a build product, but go:embed needs the directory to exist in a clean` / `# checkout (ticket 77 AC#1), so one anchor file stays tracked; panel.Assets.` / `# Built() reports false while the bundle is only that anchor, and the host shows` / `# an "assets not built" state instead of ever rendering it.` ⇒ **出处具名＝票 77 AC#1**；`frontend/embed.go:9-14` 是同一句话的第二处 | `logs/gitignore-hits.txt`＋整行读 `.gitignore:13-26` |
| 口径甲「dist 入库」撞什么 | ① 直接撞 `.gitignore:22`（要动那条 negation）；② 撞 d22scan 的 `frontend/` 计数分母——`tools/d22scan/scan_test.go:2021` 逐字 `{"dist/.gitkeep", false, false, "the negation keeps the go:embed anchor - THIS is what preserves CI's 40"}` ⇒ **入库会改 CI 那枚 40 的分母**；台账 `A207`（`:5875`）已经在记"扫描器把 `frontend/dist/` 算进 `frontend/` 计数"这一形；③ `tools/d22scan/gitignore.go:13` 注释具名 `go:embed needs the directory to exist in a clean checkout (ticket 77 AC#1)` | 逐枚现读 |
| 口径乙「每次构建生成」撞什么 | ① `build.ps1:74` 那枚 `Write-Host` 要换成 `npm ci && npm run build` ⇒ 脚本从此**依赖 node 工具链**，而它 `.DESCRIPTION:4` 自我声明是 `One-shot build of wisp.exe on Windows (SPEC-11 §2.2). From a clean clone,`（SPEC-11 §3.2 的解法是走 docker 导出：`产物落 \`assets/web/\` 供 \`//go:embed\`；宿主机无需 Node`）；② **SPEC-11 的落点与活码不是同一枚路径**：SPEC-11 `:45`/`:79-80`/`:83` 全指 `assets/web`，活码指 `frontend/dist`；`assets/` 目录现量不存在；③ 顺序是硬约束：embed 在 `go build` 期解析（`internal/tools/tasklist_deferred_236r3_teeth_test.go:105` 那句 `//go:embed is resolved` 讲的就是这维），所以"先生成后编译"不可反；④ 自用期 SLO／D32 那几枚门（票 263 一族）在 `slo-check.ps1` 消费 exe，多一步 npm 会改 CI 时长面 | — |
| 规格到底写清没有 | **写清了一半**：SPEC-11 `:44` 说了"存在则跳过；`--with-frontend` 时走 docker/frontend.Dockerfile"，`:84` 说了"S5 之前 frontend 目录可空，构建脚本跳过该步" ⇒ **"跳过"本身是有规格的**；而**"S5 起这一步落在哪条命令、由谁执行、CI 哪枚 job 产带包体的 exe"规格没覆盖**（`docs/PLAN.md` D41 全节 `:2882-2926` 整节读，只写 `resources\  ← 前端产物（亦可 embed.FS）` 与 `(d) CI 产物矩阵：自用期只出 windows/amd64`，**没有一个字讲产物由谁生成**） | 整节读 PLAN.md:2876-2927；SPEC-11 `logs/spec11-hits.txt`＋`logs/spec11-and-t34.txt` |

---

## §2 差集表（哪一格已有票、哪一格无票）

| # | 那一格 | 有票？ |
|---|---|---|
| 1 | `//go:embed` 通路存在且 `Built()` 能区分锚/真产物 | **有**＝票 77 AC#1（未勾）；产码在场 |
| 2 | 宿主在缺产物时**说一句人看的话**而非渲染锚 | **有**＝票 33 AC#12 的判据面＋票 248 AC#9；产码已落地（`panel_host_windows.go:428-429`＋`:442-451`） |
| 3 | 冷启动探测页盖掉真页面 | **有**＝票 33 AC#13，**已落地**；⚠ 票面 `:56` 的旧路径 `assets/web/dist` 与活码 `frontend/dist` 不一致（登记形状，⛔ 我不改票面） |
| 4 | 一枚能区分"embed 只到锚"与"真有一包产物"的**会响的仪器** | **有**＝票 33 AC#12／票 248 AC#9 都写了判据形状（问能力、⛔ 不问构建退出码、不问文案），**两枚都未勾**，且都具名卡在"谁把 dist 填上" |
| 5 | **出货那一刻谁执行 `npm run build`，把它钉在 `go build` 之前**（本机链＝`build.ps1` 第 2 步；CI 链＝产 exe 的那三枚 job） | **名义上有、实质上没有**＝票 34 `:42` 那一格写的是 `docker build → dist → go build`，落点 `assets/web/dist`＋`docker/frontend.Dockerfile`（该文件**不在盘上**，`with-frontend`／`frontend.Dockerfile` 词面**全票池零命中**）；它的 blocker `01-build-chain` 已结案且没带走 assets/web 目标 ⇒ **实质射程未覆盖活码路径** |
| 6 | **CI 产出的 `wisp.exe` 真带页面字节**（embed 条目枚数 > 1、artifact 有人消费） | **⛔ 无票＝本件的差集那枚** |
| 7 | 启动自检（`wisp doctor`）里"页面产物在不在"那一维 | **⛔ 无票**；且 `PLAN.md:2900-2901` 的 D41 自检四项**不含这一维**＝规格没覆盖（不是我推断规格要求它） |

⇒ **给编排者的一句话**：**没票那一格＝第 6 行**（"CI/本机产出的 exe 里真得有包体"这条**出货判据**），第 5 行是它的近亲（名义有票但名册过期），两行是否并成一枚由你裁；本件⛔ 不新开票。

---

## §3 要在干净机器上真开得出内容，缺的是哪一步（⛔ 不提"改前端"，前端不归我们）

按依赖序，每步都只说**链上缺的那一手**，⛔ 不指定实现：

1. **出货命令里那一枚"先生成"的手**：`scripts/build.ps1:73-74` 那步今天只 `Write-Host`。要在它和 `:129 go build` 之间真存在一次"把 `frontend/dist` 填上"的执行（脚本本地 `npm`、或 SPEC-11 §3.2 的 docker 导出、或票 34 `:42` 名义上那形）——**选哪一形归编排者/owner，规格只到"存在则跳过"这一句**（SPEC-11 `:44`）。⚠ 顺带要把 SPEC-11 的 `assets/web` 落点与活码 `frontend/dist` 的对齐**具名定一枚**（两边现在各说一套：SPEC-11 `:45/:79-80/:83` vs `frontend/embed.go:19`）。
2. **CI 侧那一枚"同一工作区"的手**：产 exe 的三枚 job（`:540`、`:726`、`:789`）里至少一枚要在 `build.ps1` 之前拿到真 dist（装 node＋npm ci＋build，或 download 一枚 artifact）。今天 `upload-artifact` 只有 SLO 报表两发（`:732`/`:795`），**没有 dist 这一枚 artifact，也就没有消费方**。
3. **一条会响的判据**：把"exe 里真带包体"钉成一问能力的事实（embed 条目枚数／入口真存在，判据形状票 33 AC#12 与票 248 AC#9 已经写好，⛔ 不许问构建退出码、不许问文案），现成载体＝`wisp panel-assets -manifest`／`-check`（`panel_assets.go:105-124`）。没有这一条，第 1/2 步漏掉时**盘上证据形状是全绿**（票 244 AC#6 那枚教训的原话）。
4. **开机自检那一维**（差集第 7 行）：`wisp doctor` 今天不答"包体在不在"（`doctor.go` 339 行零命中）。这一格**规格没覆盖**，要不要补由编排者裁，⛔ 我不替规格发明要求。
5. ⚠ **顺序门禁具名**：票 244 `:29` 逐字写着"**本格未答 ⓐ 之前，票 244 整票不许动构建链**"。上面第 1/2 步都要动构建链 ⇒ **是否与票 244 那枚门禁冲突，必须由编排者现裁**，任何实现腿不许自己排这个队。

---

## §4 判不动（只能真机量、而 `winlive` 未批）

| 那一格 | 为什么这碗里量不到 | 谁能量 |
|---|---|---|
| 干净机器双击今天**到底看到什么**（那句 not-built 通知页有没有真渲染进 WebView2 控件） | `SetHtml` 只在真机真控件上有读数；本腿⛔ 不跑 exe、不启动 GUI，`winlive` 未批 | 真机腿（`winlive` 批了之后）；现读凭据只有 `panel_host_windows.go:442-451` 的源码 |
| 本机那份未跟踪 dist（4 枚）**是不是 09-27 那一版、跟 `frontend/src` 当前源码同不同步** | 我只数枚数没读字节（边界）；`find -newer go.mod` ＝ **0**，而 `go.mod` 自身 mtime 不受我控制 ⇒ **这一把尺判不动日期**，编排者转述的"09-27"我**核不了来源、只能标〔读不到〕** | 前端侧那枚会话（不归我、也不归编排队写） |
| `wisp.exe` 里 embed 实际条目枚数（本机那份 exe 是带包体编的还是锚编的） | 要么跑 `go build`（⛔ 禁），要么读 `build/wisp.exe`＋`objdump`（同波写腿会重编 `build/`，我落一枚读数就是给它洗脏）⇒ **不取** | 票 33 AC#12／248 AC#9 落地腿或真机腿 |
| CI 那三枚产 exe 的 job **真跑出来的 exe 里有没有内容** | 快照只有 `HEAD` 的 yaml；真 run 的步序与产物要 `gh` 拉 run 日志（大输出，且属于别的腿的发号） | 编排者或 ci 腿具名发号 |
| "SPEC-11 要 `assets/web`、活码用 `frontend/dist` 到底哪一边算定案" | 文本对文本能读，但**定案权不在只读腿**：PLAN.md D41 只写 `resources\ ← 前端产物（亦可 embed.FS）`，没钉目录 | 编排者裁（必要时上 owner 清单） |

---

## §5 具名报回：与我（子代理）先前听说的转述不符之处

1. **"很可能开出一个空白窗口、什么都不说"——今天不成立（就代码形状而言）。** `cmd/wisp/panel_host_windows.go:428-429` 在 `serveEntry()` 报错时立刻 `serveNotBuiltNoticeLocked()`，`:449-450` 写出的那页含 `<p>panel assets unavailable: the embedded bundle is not built.</p>`；`internal/panel/assets.go:76-78` 从进门就拒渲染锚文件。**真机有没有真把这一页画出来＝§4 第一行，量不到。** 所以：缺陷形状要登记的话，登记的应是**"链上没有生成步骤"（§2 第 5/6 行）**，不是"UI 沉默"。
2. **票 33 AC#13 那发"探测页盖掉真页面"我看是已经落地的**（`:420-426` 注释自认 `the whole fix`＋现读时序 `:426` 在前、`:428` 在后）。若编排者手上还把它当"未修的缺陷形状"在排程，这一条要更正。⛔ 我没翻任何框。
3. **边界解释**：我读了 `frontend/embed.go` 并数了 `frontend/dist` 的文件枚数。依据＝派单允许清单里的"Go 侧内嵌资源的代码"与本件 E 问具名要求的 `git ls-files frontend/dist`；禁句具名是 `frontend/src/**`（**我零打开**）与 `design/**`（**我零读**）。⛔ 我没读 dist 里任何一字节内容、没转达前端语义。若这一层解释超出你给的边界，请具名纠正，我按纠正重做。
4. **本机 dist 枚数与台账数字差 1**：`docs/reports/pending-and-issues.md:1567`/`:4427` 记"三枚未跟踪 `frontend/dist/` 产物"，我今天 `find frontend/dist -type f` ＝ **4**（含 `.gitkeep`）。差异可能只是 `.gitkeep` 计不计入或前端侧新落了一枚，**不在我射程内定性**，只登记。
5. **SPEC-11 的 `--with-frontend` 与 `docker/frontend.Dockerfile` 是"规格有、盘上没有、票池零命中"三件同向**（`docker/` 现量 5 枚、`grep -rl` 两词面＝0）。⇒ 编排者转述里"前端产物由谁生成"这一问，**规格侧确实给过一个答案（docker 导出到 `assets/web`），但那个答案既没实现、也已被活码路径 `frontend/dist` 取代**；这一格要不要算"规格过期"归你裁，本件只报现状。
