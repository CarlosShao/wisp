# 274 — 出货那枚 exe 里到底有没有页面：今天**没有任何判据要求它有**，而 `build.ps1` 第 2 步那句注释两半都是过期的

立票时刻 `2026-10-07 10:4x +08`，锚点 HEAD `ca4ae2af`。来源＝机主 10-07 让我查另一棵前端分支该不该并进来，我按他拍的**乙＝先预演**跑完只读预演（`.scratch/wisp/probes/merge-dryrun/1/dryrun.md`，本票随该件同批入库）。
⚠ **预演 §3 里我写重了一句，本票第一段就是更正它**：我当时写"没人看得见空 bundle"，现读发现 `cmd/wisp/panel_host_gate_test.go:77 TestPanelBundleShapeSeparatesAnchorFromRealPage_AC12` **确实在跑**且写得很诚实——真缺口不是"没人管"，而是**那把尺两形都放行**（见 §2④）。

## 1. 现象（大白话）
同一份源码，在我这台电脑上编译出来的程序**带着页面**（面板一开就有内容）；在 CI 的干净机器上编译出来的**一个页面字节都不带**（打开面板只会得到一句"还没构建"）。而这两种结果**今天都能让所有门绿灯**——包括每次推送都在跑的那几发。⇒ 也就是说"我这边能看"和"你下载的那份能看"是两件不同的事，而我们**没有任何一枚判据把它们对上**。

## 2. 根因（四条我全部现量，逐条给尺）
① **构建脚本的第 2 步是空步，且那句注释两半都不成立**：
- `scripts/build.ps1:73` `# --- 2. frontend ---` → `:74` `Write-Host 'build.ps1: frontend step skipped (no frontend yet; embed lands in S5 per SPEC-11 §2.2)'`
- "embed 到 S5 才落"＝**过期**：`//go:embed all:dist` 早已在树里（`frontend/embed.go`，`go list -f '{{.EmbedFiles}}' ./frontend` 现回 4 枚文件，见③）。
- "没有前端"＝**反的**：前端有，CI 每晚还构建它——只是**构建它的那个 job 与造 exe 的那个 job 不是同一枚**（见②）。
- 尺＝`sed -n '73,74p' scripts/build.ps1`；`go list -f '{{.EmbedFiles}}' ./frontend`
② **CI 里"造 exe"与"造页面"从未在同一枚 job 相遇**：
- `npm run build` 全文件**只有 1 处**＝`.github/workflows/ci.yml:903`，它在 `lint-frontend`（`:871`，`runs-on: ubuntu-latest`）里，跑完就丢（job 之间不共享工作目录）。
- 真正产出 exe 的 `scripts/build.ps1` 在 CI 跑了 **3 处**＝`ci.yml:571`（test-windows 的 cgo build smoke）、`:757`（slo-smoke）、`:820`（slo-full），三处**都没有先跑 npm**。
- ⇒ CI 造出的每一枚 `wisp.exe` 都是 `built=false` 那一形。
- 尺＝`grep -n 'npm run build' .github/workflows/ci.yml`；`grep -n 'scripts/build.ps1' .github/workflows/ci.yml`
③ **干净检出时 `all:dist` 只吃到那枚占位文件，而且 `go build` 完全不报错**（这条我用仓外一把隔离尺证的，⛔ 没动本机 dist）：
- `frontend/dist` 在版本库里跟踪的只有 **1 枚**＝`.gitignore` 规则下刻意留的锚文件（尺＝`git ls-files frontend/dist`）。
- 仓外台件 `$HOME/tmp/embedprobe-274/`：只放 `dist/.gitkeep` ＋ 一枚同样写 `//go:embed all:dist` 的 `embed.go` ⇒ `go list -f '{{.EmbedFiles}}' .` 回 **`[dist/.gitkeep]`**、`go build ./...` **rc=0**。
- 对照本机现值：`go list -f '{{.EmbedFiles}}' ./frontend` ＝ `dist/.gitkeep` + `dist/assets/index-*.css` + `dist/assets/index-*.js` + `dist/index.html`（**4 枚**，来自本机 09-27 那份碰巧留下的产物）。
- ⇒ **"能不能编译"与"里面有没有页面"今天毫无关系**；带不带取决于编译机器上碰巧有什么。
④ **现有那把尺钉的是"降级要 fail-closed"，不是"出货不许降级"**（这就是为什么它全绿却不解决问题）：
- `internal/panel/assets.go:44` `fs.Sub(frontend.Dist(), "dist")` → `:54` `newAssets` 里 `:57` 只有 `index.html` 存在才置 `built=true`（`internal/panel/assets.go:29 const EntryFile = "index.html"`）→ `:75` `Resolve` 在 `:76` `!Built()` 时 `:77` 返回 `errNotBuilt`（`:34` 那句还诚实写着"run npm run build in frontend/"）。**这套语义是好形状，本票不许动它**（见 §3）。
- `cmd/wisp/panel_host_gate_test.go:77 TestPanelBundleShapeSeparatesAnchorFromRealPage_AC12`：`:115/:118/:124/:127` 管 `built=true` 那一支，`:131/:134/:137` 管 `built=false` 那一支，`:143/:147` 还专门钉了两支被混为一谈。⇒ **它两种都接受**，所以 CI 恒走 `false` 支而**整包照绿**。
- 尺＝`grep -n -E 'errNotBuilt|\.Built\(\)' --include=*.go cmd internal`；`grep -n '^func Test' cmd/wisp/panel_host_gate_test.go`
⑤ **本机反证一枚，防止我把上面说重**：`build/wisp.exe panel-assets` → `panel assets embedded: 4 files, entry=index.html built=true` ⇒ "干净机器上带不到页面"**是换机器时的必然、不是本机今天的故障**。

## 3. 关键约束（⛔ 触碰即退回）
- ⛔ **不许改 `internal/panel/assets.go` 的 fail-closed 语义**（`errNotBuilt`／`Built()`／`Resolve` 那一族）：那是"不给假页面"的安全形状，票 33 的 AC#12 正钉着它。本票只问"出货那枚该不该是 built=true"。
- ⛔ **不许为了"让产物入库"去动 `frontend/.gitignore` 里 `dist/*` 那条规则**（那是票 77 AC#1 定的规矩，不是疏漏）。
- ⛔ **不许动 `winlive`**：任何"真开面板看一眼"的读数都不在本票射程（只能登记〔仅本机可量，未批〕）。
- ⛔ **不碰页面代码**：机主 09-28 定的口径——`frontend/**` 与 `design/**` 不派、不写、不读、不转达。本票只读构建胶水与 CI。⚠ **本票取证时我读了 `frontend/embed.go` 的第 15–24 行（10 行 Go 构建胶水，为了确认那条 embed 指令）**，越界与否归我裁并已具名记进台账 `A654`，页面侧其余一个字节都没读。
- ⛔ **不许顺手把"build.ps1 真跑 npm"当成本票的修法直接落**：那是形态选择，要走 AC#2 的代价表＋人工批准（它会把 Node 变成 Windows 产线的硬依赖，与票 34 里"host needs no Node"那条约束正冲）。
- ⛔ 不许为凑绿放宽任何既有断言；不许动 SLO 阈值／golden／`thresholds.go`；`git commit` 必须带显式 pathspec。

## 4. 与票 34 的分工（⛔ 不重复劳动）
票 34（`.scratch/wisp/issues/34-frontend-scaffold.md`，未结案，`Last update 2026-09-19`）已经声称"embed pipeline"，AC#1 写的是 `docker build → dist → go build serves the scaffold in the real panel window (33)`。两票**不撞**，因为：
- 34 的前提路径**在树里不存在**：`docker/frontend.Dockerfile` 跟踪枚数 **0**（尺＝`git ls-files docker | grep -ci front`）、`assets/web` 跟踪枚数 **0**（尺＝`git ls-files | grep -c '^assets/web'`）；仓库实际用的是 `frontend/` + vite + `frontend/dist`。
- 34 的 AC#1 **要真开一扇窗**才能判，今天不许（`winlive` 未批）⇒ 它是一枚"现在无法执行"的判据。
- ⇒ **34 不动、不并、不改写**；本票只补那一格里**不需要开窗就能自动判**的那半句："交给用户的那枚 exe 里到底有没有页面字节"。真开窗口那一格仍归 34/33。

## 5. 完成判据（AC 框只有编排者能翻，腿一枚都不许碰）
- [ ] **AC#0 先把两枚数量清（不许引用本票面的数）**：(a) CI 的 `lint-frontend` 那一发真跑成功后，`npm run build` 到底产出几枚文件、落在哪；(b) 同一发 run 里 `build.ps1` 造出的 exe，跑 `wisp.exe panel-assets` 回的是 `built=true` 还是 `false`、枚数差几枚。⛔ 台件不许写进 CI 之外看不到的地方。
- [ ] **AC#1 未修码读数（先证今天拦不住）**：把 bundle 换成"只剩锚文件"那一形喂给 AC#12 那把尺——**期望今天整包全绿**，这就是"出货带不带页面零判据"的凭据。⚠ 必须走 `go test -overlay` 或仓外合成 bundle，⛔ 绝不许真删/真盖本机 `frontend/dist` 里那 4 枚文件（那是别的会话与签收现场的料）。另给一发正向对照：同一把尺喂"带 `index.html` 的合成 bundle"必须走 `true` 支，证明尺没瞎。
- [ ] **AC#2 三形代价表，⛔ 本票不选形**：甲＝`build.ps1` 第 2 步真跑 `npm ci && npm run build`（代价：Windows 产线与每台编译机从此**硬需要 Node**，与 34 里"host needs no Node"相冲，CI 时长代价要现量）；乙＝CI 加一枚跨 job 产物传递（`lint-frontend` 上传 dist → Windows job 下载后再 `build.ps1`，exe 才真带页面；代价：新 artifact 依赖＋job 顺序）；丙＝产线不动，只加一枚**发布前闸门**（打包/签名之前强制 `wisp panel-assets` 回 `built=true` 且入口字节>0，否则红）。逐形给：动哪几枚文件、要不要新 job/新 artifact、CI 时长增量、**"这条边以后谁能看见"**（本机？CI？只有发布步？）。
- [ ] **AC#3 那句过期注释的处理边界**：只许把 `scripts/build.ps1:74` 改成**说实话**（例如"这一步在这里跳过，出货判据在 AC#2 所选形里"），⛔ 不许在改注释的同时把第 2 步实做（那是形态选择，要 AC#2 的表＋批准）。本格单独可勾、单独可撤。
- [ ] **AC#4 反形敏感性（定向突变）**：新判据必须写成"把 dist 换成只剩锚文件那一形 ⇒ **出货判据必须红**"。⛔ 不许写成"有没有 built 标志"这种两形都绿的恒真句——`AC#12` 就是活教材：它诚实、它在跑、它永远放行。
- [ ] **AC#5 零新契约面**：⛔ 不新造 `C##`，⛔ 不往 `C17` 方法白名单加名字，⛔ 不动 D29/D41/SPEC-11 的文字（`SPEC-11 §2.2` 那句"embed lands in S5"要不要改＝人工批准，本票只上报不修）。
- [ ] **AC#6 门禁四数**：`sh scripts/d22scan.sh` rc=0 且正控先绿；`gofmt -l`／gofumpt v0.12.0 `-l` 只喂 `.go`；`go vet ./cmd/wisp/ ./internal/panel/` rc=0；终态 `go test ./cmd/wisp ./internal/... -count=1` 的**逐名红名册与起手作差＝0**。⚠ `cmd/wisp` 写面**同一时刻只能一枚写手**（包级互斥，不接受"改的是不同文件"）。
- [ ] **AC#7 还原自证**：所有突变只走 `go test -overlay`（拷贝放仓外）；终态 `git status --porcelain -- cmd internal scripts` ＝起手，被审文件与 `git show HEAD:` 逐串对回；⛔ 本机 `frontend/dist` 那 4 枚文件枚数与哈希起手/终态必须相同（尺＝`go list -f '{{.EmbedFiles}}' ./frontend` 两头发）。
- [ ] **AC#8 `SPEC-11 §2.2` 那两步的归属裁决（⛔ 本票不选形，由 `274-a1` §3 现量升格而来）**：§2.2 第 2 步逐字要求「存在则跳过；**`--with-frontend` 时走 `docker/frontend.Dockerfile`（§3.2）**」，而今天 ①`build.ps1` **没有这枚开关**（无条件跳过）②`docker/frontend.Dockerfile` 跟踪枚数 **0** ③该节第 3 步还写着「embed **`assets/web`** 已就位」，而 `assets/web` 跟踪枚数 **0**（真实通道是 `frontend/dist`）。⇒ 要裁的是：**这枚开关该由 `build.ps1` 兑现，还是承认 §2.2 那两步已被"CI 里 `lint-frontend` 跑 vite ＋ `frontend/dist`"这一形取代**（取代＝要动 spec 文字＝人工批准，另案；⛔ 任何腿不许"顺手补个开关"来让代码与 spec 对上）。凭据＝`.scratch/wisp/probes/274/a1/census.md` §3（我现读复认，见本票末节）。

## 6. 排程（编排者自己记，不许腿替我改）
本票**按住**，排在 `cmd/wisp` 写面空出之后：此刻 `272-r2` 是本波唯一跑 Go 突变的腿，它的终态门禁要跑整包红名册，任何写腿落脏文件都会洗它的读数。
AC#0/AC#1 可以**先派只读普查腿**（不写产码、不跑整包）——它与票 273 的落点普查 `273-a1` 同族，但**不许在同一枚包里同时跑整包测试**。

## 收 `274-a1`（只读普查；件 `.scratch/wisp/probes/274/a1/census.md`，215 行／29,577 字节／占位 0；三枚 commit，起手锚 `15699a2f`；10-07 11:1x，落账 `A656`）

★**两枚对我不利的更正，我都现量对回了，且都让本票的射程更准（⛔ 票面 §2① 原句不抹）：**

1. **§2① 那句"注释两半都过期"说轻了。** `docs/specs/SPEC-11…:39 §2.2「构建顺序（一键流程）」` 我现读逐字，它第 2 步写的是「**前端产物：存在则跳过；`--with-frontend` 时走 `docker/frontend.Dockerfile`（§3.2）**」——
   - 该节**通篇没有"S5"这个字**（`S5` 在 §3.2 标题与 §6 里），⇒ `build.ps1:74` 那句 `per SPEC-11 §2.2` 是**引用指错了节**，不只是内容过期。
   - 更要紧的：**`build.ps1` 今天没有 `--with-frontend` 这个开关**，它是**无条件**跳过；而 §2.2 第 3 步还写着「embed `assets/web` 已就位」，`assets/web` 跟踪枚数＝**0**（我早前量过）。⇒ 形状应从"一句过期注释"升格为 **「SPEC-11 §2.2 的第 2、3 步今天都没实现，且注释把责任推给了一节没说过这话的 spec」**。⛔ spec 文字我不改（改它＝契约面，另案）。
   - ⇒ 新增 **AC#8（本票不选形）**：把"`--with-frontend` 这枚开关该不该由 `build.ps1` 兑现、还是承认 spec 那一步已被 `frontend/dist` ＋ CI `lint-frontend` 取代"摆成一支待裁项；⛔ 不许任何腿"顺手补个开关"来让 §2.2 变得对上。
2. **"没有任何判据要求它带页面"这句要配一枚例外——牙其实存在，只是没人咬。** `cmd/wisp/panel_assets.go:125-129`（`default` 支）现读：`if !assets.Built() { fmt.Fprintln(os.Stderr, "wisp panel-assets: assets NOT BUILT"); return 1 }` ⇒ **`wisp panel-assets` 不带子命令时，未构建就直接 rc=1**。腿的尺＋我的复跑一致：`ci.yml` 里 `panel-assets` **唯一命中在 `:912` 的一句注释**，`scripts/**` 与 `tools/**` 里 `setup-node|npm ci|npm run build` **0 命中**（正控＝同尺在 `.github/**` 命中 4 枚）、`download-artifact` 跟踪语料 **0 枚**（正控＝同尺 `upload` 4 枚）。⇒ **本票真正的缺口收窄成一句：现成的牙（rc=1）没有任何脚本或 CI 步骤去咬。** 这把丙形（发布前闸门）从"要新造判据"降成"要接一根线"，代价量级也因此有了现成数（腿给的：`lint-frontend` ≈24-25s、`slo-smoke` ≈2m、`test-windows` ≈10m）。
3. **腿替我把 CI 侧坐实了**：三发 `dev` run（`37545246395`／`37406757402`／`37406422380`）的步级 conclusion 证明 **`npm run build` 真跑过且全 success**，且**没有一发在 `build.ps1` 之前有 node 步** ⇒ §2②那把尺不再是"yaml 里写着"而是"run 里出现过"。另：`needs:` 全 `ci.yml` **0 枚** ⇒ 乙形要新加依赖边；`release` job 无落点 ⇒ 丙形的"发布前"**今天没有栖息地**（这两条直接进 AC#2 的代价表）。
4. **⛔ 我不采纳它的一条**：它报"工单盘上真名比派单少个『1』"。尺＝`grep -rl 'build-ps1-step-2' .scratch docs` ＝ **0 命中**，仓内所有引用与本文件名一致 ⇒ 那处拼写差异**盘上查无实据**，记为"腿的自报未经我复认"，不影响它其余读数。
5. **它自报未做完的三格我不代填**：AC#0(a) `npm run build` 在 CI 侧产出枚数（要全量 run 日志）、更早 run 扩样、`placeholder` 独立尺。⇒ 归后续验收腿逐枚处置。**本票九枚框（AC#0–AC#8，其中 AC#8 由本次收表新增）一枚未翻。**
