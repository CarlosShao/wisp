# 274-a1 只读普查 census（票 274：没有钉要求交出去的 exe 真带页面字节）

腿：274-a1（只读普查，禁产码/禁 go test/禁改他文件）。本件唯一产码面＝本文件 + `logs/`。

## §0 起手锚

- `git rev-parse --short HEAD` → `15699a2f`
- 分支 → `dev`
- `date` → `2026-10-07 10:28:00 +0800`
- `tasklist | grep -ci "go\.exe"` → `2`（跑前量一次，本腿全程未启 go test；同仓写腿 272-r2 读数窗口未被本腿触碰）
- 工单实际文件名（与派单文字有一字之差，见 §5-冲突①）：
  `.scratch/wisp/issues/274-no-nail-requires-the-shipped-exe-to-carry-the-page-build-ps-step-2-comment-is-false-in-both-halves.md`
  （派单里写的是 `build-ps1-step-2`，盘上是 `build-ps-step-2`）

## §1 全仓语义变体普查——「有没有任何东西要求交出去的 exe 里真有页面字节」

方法说明：一律用 Grep 工具（ripgrep 内核），下表「命令」列为可直接复跑的 `rg` 等价式；
「命中枚数」为 **count 模式全仓逐文件总数**（含 `.scratch/**` 探针/工单回声，票 274 自身命中已单独标注，引用时可减）。
本文件 `.scratch/wisp/probes/274/a1/census.md` 写在扫描之后，不污染以上读数。

### 族1 `built`／`Built(`／`errNotBuilt`／`not built`
- 命令：`rg -c 'Built\(|errNotBuilt|not built|built'`（全仓）
- 命中：**2857 枚 / 852 文件**——裸 `built` 是子串尺（built-in/rebuild/anchored 全进），**语义读数看符号级**：
  `internal/panel/assets.go` 12、`internal/panel/assets_test.go` 17、`cmd/wisp/panel_host_gate_test.go` 19、
  `cmd/wisp/panel_assets.go` 4、`.gitignore` 3、票 274 自身 9（`.scratch` 回声不计入他数）。
- 前三条命中（glob `!.scratch/**`，逐字）：
  1. `.github\workflows\ci.yml:256`: `        #    built this job it produced ZERO findings (step conclusions:`
  2. `.github\workflows\ci.yml:368`: `        # exec'd, no test binary is built, and the milliseconds-sensitive packages`
  3. `.gitignore:20`（族1+族2 交叉，本节最有信息量的原文，逐字）：
     `# Built() reports false while the bundle is only that anchor, and the host shows`
- 判读：`errNotBuilt`/`Built()` 一族只活在 **运行时报告面**（assets.go / panel_assets.go / doctor）与 **测试面**（assets_test / panel_host_gate_test）；
  **没有任何构建步或 CI 步以 `Built()` 为门**（ci.yml 里 built 的两枚命中都是注释词，非门）。正控＝非零读数本身（assets.go 12 枚在列）。

### 族2 `placeholder`／`anchor`／`gitkeep`
- 命令：`rg -ic 'placeholder|anchor|gitkeep'`（全仓，不分大小写）
- 命中：**5581 枚 / 1037 文件**（`anchor` 一词在探针/台账里是通用词，子串噪声极大；语义锚定读数看下面的 gitkeep 行）。
- 前三条命中（glob `!.scratch/**`，逐字）——三条全部来自根 `.gitignore` 的 dist 锚条款：
  1. `.gitignore:19`: `# checkout (ticket 77 AC#1), so one anchor file stays tracked; panel.Assets.`
  2. `.gitignore:20`: `# Built() reports false while the host shows` → 原文全文见族1第3条（同一行）
  3. `.gitignore:23`: `!frontend/dist/.gitkeep`
- 正控：`git ls-files | grep -i gitkeep` → 唯一跟踪件 `frontend/dist/.gitkeep`（文件名维），与编排者现量一致。
- 判读：**仓里对「dist 只有锚件」这件事的态度是写进 .gitignore 注释的说明文，不是一枚门**。`placeholder` 在本语境 0 枚语义命中（全为 CSS/表单占位符噪声，`.scratch` 外无「页面是 placeholder」式表述——该子结论射程＝上述一条尺，未另开独立尺，若验收要独立读数请按命令复跑）。

### 族3 `EmbedFiles`／`go:embed`
- 命令：`rg -c 'EmbedFiles|go:embed'`（全仓）
- 命中：**308 枚 / 142 文件**；非回声侧关键分布：`.github\workflows\ci.yml` 2、`frontend\embed.go` 2、
  `internal\panel\assets.go` 3、`cmd\wisp\panel_host_gate_test.go` 1、`docs\specs\SPEC-11-build-deploy-containerization.md` 1、`tools\d22scan\scan_test.go` 2。
- 前三条命中（glob `!.scratch/**`，逐字）：
  1. `.gitignore:18`: `# dist/ is a build product, but go:embed needs the directory to exist in a clean`
  2. `.github\workflows\ci.yml:289`: `        # "// go:embed"; the real directive is embed.go:19) and KEEPS being`
  3. `.github\workflows\ci.yml:902`: `      - name: build (vite build -> frontend/dist, the bytes go:embed carries)`
- 正控：`frontend/embed.go` 2 枚在列（真 directive 在 `embed.go:19`，ci.yml:289 注释亲口指认）。
- 判读：**没有任何东西检查 `EmbedFiles` 非空**——命中全在注释/测试自述/文档；ci.yml:902 那步只存在于 lint-frontend job（自产自检，不喂别的 job，见族4/§2）。

### 族4 `upload-artifact`／`download-artifact`
- 命令：`rg -c 'upload-artifact|download-artifact'`（全仓，glob `!.scratch/**` 的 count 为 10 枚/4 文件）
- 命中：**10 枚 / 4 文件**：`.github\workflows\ci.yml` 4、`scripts\slo-freshness.sh` 1、`docs\evidence\s1\134-ac6-contended-no-conclusion.md` 4、`docs\reports\pending-and-issues.md` 1。
- 前四条命中（glob `!.scratch/**`，逐字）：
  1. `scripts\slo-freshness.sh:34`: `#                                 written" when it refuses). upload-artifact`
  2. `.github\workflows\ci.yml:763`: `        uses: actions/upload-artifact@v4`
  3. `.github\workflows\ci.yml:826`: `        uses: actions/upload-artifact@v4`
  4. `.github\workflows\ci.yml:842`: `          #      actions/upload-artifact@v4 -> ea165f8d65b6e75b540449e92b4886f43607fa02,`
- `download-artifact` 精确枚数在补测格（见末列）——前 6 条窗口里全为 upload 侧；若计数为 0 则正控＝同尺 upload 侧 4 枚（ci.yml）已证尺可用。
- 判读：**跨 job 传物的口子在 slo 线（slo-smoke/slo-full，763/826 附近），页面字节没有走这个口子的任何一步**——`download-artifact` 在 ci.yml 命中里不出现（待补测格钉死）。

### 族5 `setup-node`／`npm ci`／`npm run build`（`.github/**`、`scripts/**`、`tools/**`）
- 命令：`rg -c 'setup-node|npm ci|npm run build'`（全仓 332 枚/79 文件，**绝大部分是 `dist` 混入？否——本尺三词精确；大头在 `.scratch` 探针回声与 `docs/reports/frontend-*`**，具名分布见下）
- 作用域读数：
  - `.github`：4 枚，全部集中在 lint-frontend job：`ci.yml:879` `      - uses: actions/setup-node@v4`；`ci.yml:887` `      - name: npm ci (lockfile is the only source of deps)`；`ci.yml:888` `        run: npm ci`；`ci.yml:903` `        run: npm run build`（903 行与编排者现量一致；902 是步名，见族3）。
  - `scripts`：**0 枚**（该作用域读数为零）。正控＝同尺在 `.github` 4 枚命中 + 同文件 `scripts/build.ps1` 用另一把尺（族5 内容尺加宽 `frontend|dist` 后命中 `build.ps1:10/73/74`，见族7原文）——即「scripts/ 里 npm/node 三词真没有，不是尺瞎」。
  - `tools`：**0 枚**（根 count 逐文件清单中无 `tools\` 条目）。正控同上。
- 前三条命中（path=`.github`，逐字）：ci.yml:879 / 887-888 / 903（上面已逐字）。
- 判读：**node 世界只活在 lint-frontend 一个 job 里**；跑 build.ps1 的 job（含 :571 那个 `cgo build smoke` 步）前面没有任何 node/npm 步——:566 步名与 :571 命令行之间无 setup-node（:879 远在其后，属另一 job）。

### 族6 `panel-assets`／`built=`
- 命令：`rg -c 'panel-assets|built='`（全仓）
- 命中：**1227 枚 / 321 文件**（`panel-assets` 同时是子命令名、工单 143 文件名词、与大量探针回声；语义主件：`cmd\wisp\panel_assets.go` **18**、`cmd\wisp\main.go` 2、`cmd\wisp\doctor.go` 1、`cmd\wisp\panel_pump.go` 2、`cmd\wisp\panel_assets_143_test.go` 2、`cmd\wisp\panel_host_gate_test.go` 10、`.github\workflows\ci.yml` **1**、`internal\panel\assets.go` 2、`internal\panel\assets_test.go` 1）。
- 上下文命中样例（glob `!.scratch/**`，`-C 3`，逐字——注意该窗口多数行为上下文行，非匹配行本身）：
  - `scripts\check-path-length-budget.sh:368`: `    R[".scratch/wisp/issues/143-wisp-panel-assets-l2-cannot-produce-an-r4-card-the-taint-detector-is-not-wired-on-that-path-done.md"] = "hist: ticket filed before README rule 9 (100-char hat, added 2026-10-04); closed, name doubles as the -done anti-double-claim key"`
  - （同窗 :365/:366/:367 为同名映射行，同为长行原文，未截断。）
- **「是谁在调 panel-assets」**：调用点名册（谁分发子命令、谁写 `built=` 行）在 `cmd/wisp/main.go` 与 `cmd/wisp/panel_assets.go` 内，两件的逐行读文归 §4 批（本腿纪律：先交首笔 commit）。
- **「有没有脚本/CI 步看它的退出码或 built= 行」**：ci.yml 对 `panel-assets` 只有 **1 枚** 命中——该行原文（是步名还是注释、有没有 `continue-on-error`/退出码检查）列入补测格，§2/§4 钉。当前可下的半格判读：族1/族5 显示没有任何构建步以 `Built()`/npm 为门，故此 1 枚若为注释或纯展示步则全仓无「看门人」。
- 正控：`cmd\wisp\panel_assets.go` 18 枚在列。

### 族7 `dist`（谁生成、谁消费）
- 命令：`rg -c '\bdist\b|dist/'`（全仓）
- 命中：**1441 枚 / 431 文件**（含 `pending-and-issues.md` 台账 47、`frontend-watchdog.md` 34、探针回声；语义主件：`tools\d22scan\scan_test.go` 1274（镜像真 .gitignore 的夹具注释，见逐字第3条）、`frontend\vite.config.ts` 2、`frontend\.gitignore` 3、`.gitignore` 4、`docs\specs\SPEC-11-...md` 2、`.github\workflows\ci.yml` 1、`cmd\wisp\panel_host_windows.go` 1、`cmd\wisp\panel_host_gate_test.go` 8、`internal\panel\assets.go` 5、`frontend\embed.go` 4）。
- 前三条命中（glob `!.scratch/**`，逐字）：
  1. `tools\d22scan\scan_test.go:1274`: `// The .gitignore contents mirror the real ones (root .gitignore `frontend/dist/*``
  2. `tools\d22scan\scan_test.go:1275`: `// + `!frontend/dist/.gitkeep`, frontend/.gitignore `dist/*` + `!dist/.gitkeep`)``
  3. `tools\d22scan\scan_test.go:1277`: `// negation matters for the counts: `frontend/dist/.gitkeep` is one of the 40 files`
- **谁生成**：`vite build`（`frontend/vite.config.ts` outDir；CI 中唯一执行者是 lint-frontend 的 `npm run build`，ci.yml:902 步名自述「the bytes go:embed carries」）。
  **谁消费**：`frontend/embed.go` 的 go:embed 指令（embed.go:19）→ `internal/panel/assets.go` 的 `Assets.Built()`。
  **本机生成链**：`scripts/build.ps1:10/73/74` 明文跳过 frontend 步（三行逐字见族5 作用域读数）——干净克隆里 dist 只有 `.gitkeep`，embed 带 0 字节页面，两种形状全仓无门区分。
- 正控：上面三条 + `frontend\dist\.gitkeep` 跟踪件（git ls-files 尺）。

### 补测格（本节尚未钉死的读数，后续 commit 补）
- [x] `download-artifact` 精确计数：`rg -c 'download-artifact'` → **7 枚 / 3 文件，全部在 `.scratch`**（本 census.md 自指 5、`probes/bundle/1/census.md` 1、`probes/orchestrator/preflight-1/colors.md` 1）⇒ **跟踪件与 CI/scripts/docs 语料中 0 枚**。正控＝同尺 `upload-artifact` 4 枚（ci.yml :763/:826/:842 等）。
- [x] ci.yml 唯一 `panel-assets` 命中＝**注释**，`:912` 逐字：
  `        # which is verbatim ` + '`wisp.exe panel-assets -l2 fs.delete -irreversible delete`' + `（反引号原文如此）`
  ——它是 lint-frontend job 里解释 fixture 来历的说明行，**没有任何 CI 步执行 `wisp panel-assets`、更没有任何步检查其退出码**。
- [x] `built=` 输出行原文（`cmd/wisp/panel_assets.go`，本腿已整读该件）：
  - `:123`（-check 支）：`fmt.Printf("panel assets check: entry=%s built=%t %d asset refs resolve [%s]\n", panel.EntryFile, assets.Built(), len(served), strings.Join(served, " "))`
  - `:135`（默认支）：`fmt.Printf("panel assets embedded: %d files, entry=%s built=%t\n", len(lines), panel.EntryFile, assets.Built())`——票⑤那句 `built=true` 读数即此格式。
  - 关键形状：`:126-129` 默认支在 `!assets.Built()` 时 `fmt.Fprintln(os.Stderr, "wisp panel-assets: assets NOT BUILT")` 并 **`return 1`** ⇒ 命令本身有牙，**但全仓没有任何脚本/CI 步咬这根牙**（族6 ci.yml 仅注释 1 枚）。

## §2 CI 侧真名册（gh 只读；`gh run list/view --json`，RC=0，无大日志落盘，logs/ 本腿未用）

最近 3 发 `dev` 的 ci run（`gh run list --branch dev --limit 6`：另有两发 `slo-fresh` workflow success，与本题无关不展开）：

| run | 触发/时刻(UTC) | headSha | 整发 | lint-frontend | test-windows | slo-smoke | slo-full |
|---|---|---|---|---|---|---|---|
| 37545246395 | schedule 10-06 23:12 | cc315261 | failure | **success**(23:12:29→23:12:54, ≈25s) | failure(cgo build smoke 步 success) | success(≈2m06s) | success(次日 00:37→00:42, ≈5m22s) |
| 37406757402 | push 10-06 02:58 | cc315261 | failure | **success**(≈24s) | failure(同) | success | success |
| 37406422380 | push 10-06 02:54 | b948bcb8 | failure | **success**(≈24s) | failure(同) | success | success |

具名回答（判据＝step 的 conclusion 字段，非「yaml 里有」）：
1. **`npm run build` 真跑过且成功过吗——跑过，三发全部 success**。三发的 lint-frontend job 步名册里都有这四枚步且 conclusion=success：`Run actions/setup-node@v4` / `npm ci (lockfile is the only source of deps)` / **`build (vite build -> frontend/dist, the bytes go:embed carries)`** / `Post Run actions/setup-node@v4`。
2. **有没有任何一发 run 里 `build.ps1` 之前存在 node 步骤——这三发里没有**。test-windows 按 `npm|node|artifact|build` 正则过滤后只剩 `cgo build smoke (build.ps1 fetch-deps + mingw link + doctor)`（success）一枚匹配；slo-smoke/slo-full 只剩 `Build wisp.exe`。且三发的 slo 线与 test-windows 都没有任何 `Download *` / setup-node 步。
3. AC#0 (a) 格「lint-frontend 的 npm run build 产出几枚、落在哪」：**CI 侧取不到读数**——步名册证明它成功，但 vite 产物随 job 工作区丢弃（该 job 无 upload-artifact 步，见族4），要读数须 `gh run view --log` 全文（本腿未拉，属写腿/裁决腿的活）。
4. AC#0 (b) 格「build.ps1 造的 exe 跑 panel-assets 回什么」：**没有任何 run 执行过 `wisp panel-assets`**（§1 补测格钉死：ci.yml 仅注释 1 枚）⇒ 该格今天无 CI 读数；由形状可推 `built=false`（windows job 干净检出 dist 只有 .gitkeep），但本腿如实登记为「推断，非读数」。
5. 附注：本腿 grep 窗口（head 6）亲见 `run: ... scripts/build.ps1` 于 `ci.yml:571` 与 `:757` 两处；票面说第三处在 `:820`，第三枚的完整名册列入 §4 批的穷举尺复量（不算与本腿读数冲突，只是窗口截断）。

## §3 `SPEC-11 §2.2` 到底说了什么（逐字），那句话今天还成立吗

件：`docs/specs/SPEC-11-build-deploy-containerization.md`。§2.2 全文（:39-52）逐字：

> ### 2.2 构建顺序（一键流程）
>
> ```
> scripts/build.ps1 [-Env dev|prod] :
>   1. fetch-deps.ps1        # 校验/补齐 third_party/（有缓存则秒过）
>   2. 前端产物：存在则跳过；--with-frontend 时走 docker/frontend.Dockerfile（§3.2）
>   3. go build（cgo: CGO_ENABLED=1, CC=mingw32-gcc; embed assets/web 已就位）
>   4. 产物：wisp.exe + onnxruntime.dll + sherpa-onnx c dll 同目录（§7.1）
>   5. 输出 SHA256SUMS
> ```
>
> - CLI 与 GUI 同一二进制：无参 = GUI（`-H=windowsgui`）；`wisp run` 子命令
>   `AttachConsole(ATTACH_PARENT_PROCESS)` 输出。【SPEC】
> - 构建期注入：版本号、commit、`buildinfo` 里的 C29 minisign 公钥。

**「S5」不在 §2.2 里**——spec 的 S5 字样在相邻两处：§3.2 标题 `### 3.2 frontend.Dockerfile（前端构建，S5 起需要）`（:69）与其末条 `S5 之前 frontend directory 可空，构建脚本跳过该步`（:84，原文为「S5 之前 frontend 目录可空，构建脚本跳过该步」）；§6 job 表 `frontend` 行门禁级别 `阻塞（S5 起）`（:125）。

回答「今天还成立吗」：
- **不成立，且是三个方向的过期**：
  ① `build.ps1:74` 引「SPEC-11 §2.2」说“embed lands in S5”——**§2.2 本身没有一个字提 S5**，引用指错了节；§2.2 第 2 步的语义也不是“到 S5 才做”，而是“前端产物**存在则跳过**；要产就走 `docker/frontend.Dockerfile`（--with-frontend）”。
  ② §2.2/§3.2 承诺的 embed 通道是 **docker → `assets/web/`**（:45 「embed assets/web 已就位」、:83 「产物落 assets/web/ 供 //go:embed」）——**该通道在树里 0 枚**（票 §4 尺：`docker/frontend.Dockerfile` 与 `assets/web` 均未跟踪；本腿另验 `docker/` 下仅 `builder.Dockerfile` 在族1读数里出现）。今天真正的 embed 通道是 `frontend/embed.go` 的 `//go:embed all:dist` → `frontend/dist`，spec 文字从未更新到这个形状。
  ③ “embed 到 S5 才落”半句：embed 指令**早已在树里**（embed.go，ci.yml:289 注释自指 embed.go:19），所以「还没落」为过期；而“到 S5 才把页面喂进 exe”那半句对应的机制（frontend.Dockerfile 路径）根本没建——**票②的实证（exe 产线从未与 node 产物同 job）与 spec 的 S5 承诺是两条都没合拢的边**。
- ⛔ 本腿不改 spec 一字（AC#5：`SPEC-11 §2.2` 那句要不要改＝人工批准）；此处只回答“那句话指什么、今天对不对得上”。

## §4 喂页面给 exe 的合法落点普查（甲/乙/丙，只描形不选形）
状态：**未开始**。

## §5 边界与自陈
状态：**未开始**。
