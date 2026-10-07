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
- [ ] **AC#9 陈旧产物不许替"带页面"报绿（由 `274-a2` §4 现量升格，10-07 12:2x 编排者加）**：本机 `frontend/dist/assets` 现量＝`index-BRKj5OIJ.css`／`index-BVKlegVD.js`，而**用 `git archive HEAD` 导出后在仓外新建**的那份＝`index-vdBrT8rM.js`／`index-yy8KMgdf.css` ⇒ **哈希名全不同＝本机那份是旧源码的产物**。⇒ 本票任何"exe 带页面"的读数（含 §2⑤ 那条本机反证）**必须同时给出：①一次 fresh build 的产物文件名名册＋哈希，②exe 内嵌物与之逐枚相同**；⛔ 只报"4 files, built=true"这种枚数级读数不算过这一格。⚠ 这条同时钉住一个反向坑：⛔ 写腿**不许**为了让读数好看去动本机那份陈旧件（它属别人在飞的产物，本票只读；改它要先具名解冻）。
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

## 收 `274-a2`（只读普查；件 `.scratch/wisp/probes/274/a2/shipping-chain.md`，**267 行／30,884 字节**＋18 枚 `logs/`；六枚 commit `e18e32da`→`3dcfcd23`→`3d061f1e`→`7da2ab20`→`c47116da`→`44db024a`；10-07 12:2x，落账 `A659`）＋★编排者据它裁形

★**这枚腿是本票的选形前置**（我派单时写得很具体："run 侧真读数答那台机器有没有 node，⛔ 不许从 yaml 猜，也不许把'尺取不到'写成'没有'"）。它交回五节全有读数、**零未决格**，并且**两处顶回我票面的措辞——两处我都复跑坐实、都算它对**。

**1. 出货通路与三枚产 exe 的 job（我复跑对回）**：`scripts/build.ps1:129` 那一句 `go build -o build\wisp.exe ./cmd/wisp` 是**唯一**出货通路、入口只有 1 枚；`param` 块（`:22-25`）只有 `-Env dev|prod`，⛔ 没有前端开关；第 2 步（`:73-74`）**零命令**。`ci.yml` 里产 exe 的三枚 job 逐枚点名＝`test-windows`（`:505`，`windows-latest`，build.ps1 在 `:571`）／`slo-smoke`（`:739`，`windows-latest`，`:757`）／`slo-full`（`:797`，`[self-hosted, wisp-slo]`，`:820`）——**三枚全无 `setup-node`，也没有任何一步读 `frontend/dist`**；`upload-artifact` 只传 SLO 的 JSON，**exe 无人留档**。

**2. ★★dev 自己那 64 枚页面源码今天能构建（这条把本票从"要等合并"里解出来了）**：腿在**仓外**（`D:/tmp/wisp274a2/`，⛔ 没往仓内 `frontend/**` 落一个字节）用 `git archive HEAD frontend` 导出后跑：`npm ci` rc=0（10s）→ `tsc -b` rc=0（7s）→ `npm run build` rc=0（9s），**共 26s，出 4 枚产物／602,635 字节、`index.html` 1,068 字节**。⇒ `A658` §4 那句"乙形不等分支合并"从此有凭据。

**3. ★它顶回我的两处，我现量都坐实（原句不抹，就地更正）**：
- ⓐ**"`release` job 无落点"这句说轻了**：真读数＝**`ci.yml` 里根本没有 release job**——job 只有 6 枚（`lint`／`test-core`／`test-windows`／`slo-smoke`／`slo-full`／`lint-frontend`），且 `grep -ci release .github/workflows/ci.yml`＝**0**（连这个词都不出现）。⇒ 丙形那句"发布前闸门"今天**连栖息地都没有**，不是"有 job 没落点"。
- ⓑ**`build.ps1:74` 不是注释，是一句会印进 CI 日志的 `Write-Host`**（我 `sed -n '71,76p'` 现读逐字：`Write-Host 'build.ps1: frontend step skipped (no frontend yet; embed lands in S5 per SPEC-11 §2.2)'`；注释是它上面 `:73` 那段 `# --- 2. frontend ---`）。腿另在 `slo-full` 的 run 日志**第 183 行**捞到了这句原话 ⇒ **它已经在骗读日志的人**，不只是骗读代码的人。⇒ 本票 §2① 的"注释两半都过期"这个称法**当场更正**：要改的是**那句印出来的谎话本身**。

**4. ★★最要紧的一条新发现：本机 `frontend/dist` 是陈旧件，它会替"带页面"报假绿（我自己复跑坐实）**
- 尺＝把本机 `frontend/dist/assets` 的文件名与腿在仓外**用 HEAD 现建**的那份逐枚对：本机＝`index-BRKj5OIJ.css`／`index-BVKlegVD.js`，HEAD 新建＝`index-vdBrT8rM.js`／`index-yy8KMgdf.css` ⇒ **哈希名全不同＝内容不同**。
- ⇒ **直接后果**：本票 §2⑤ 那条"本机反证"（`build/wisp.exe panel-assets` → `panel assets embedded: 4 files, entry=index.html built=true`）**结构上成立、字节上是旧账**——它证明的是"dist 非空时 embed 真的带进去"，⛔ **不证明"当前源码能带进当前 exe"**。凡后续拿"本机 exe 带不带页面"作读数，**必须与一次 fresh build 的哈希逐枚相同**才算。⇒ 新增 **AC#9**（框数 **9→10**，一枚未翻）。

**5. ★我据此裁形（编排者裁，⛔ 写腿不许自行改选形）**：
- **落形＝乙-1**：在 `scripts/build.ps1` 的第 2 步**真构建前端产物**（`npm ci` ＋ `npm run build`，26s 有凭据），并且——⛔ **静默跳过就是本票要修的病**——**node/npm 取不到或构建失败时，具名报错＋非零退出**（不是 `Write-Host` 一句略过、⛔ 不许 `continue-on-error`）。
- **丙并入乙，不新造 job**：ⓐ说清楚了"没有 release job"，所以"发布前闸门"要么新造一枚 job（超出本票射程），要么就认：**build.ps1 那一枚非零退出本身就是牙**，与 `cmd/wisp/panel_assets.go:125-129` 现成的 rc=1 串成一条线。⇒ ⛔ 任何腿不许以"顺手加个 job"的名义扩张本票。
- **甲随批**：`=74` 那句改成实话（说清"这一步现在构建前端；机器上没有 node 就整条失败"）。
- **AC#8 维持不选形**：`--with-frontend` 那枚开关属 spec 文字面（改 `SPEC-11 §2.2`＝人工批准，另案）；乙-1 落地后 §2.2 第 2 步"存在则跳过"这半句自然对上，`--with-frontend` 那一支变**可选**。
- ⚠ **已知代价，写进派单**：`test-windows`／`slo-smoke` 这两枚 hosted job 的**作业 PATH 里有没有 node＝取不到**（腿 §3：零步试过；机侧另有 node v24.9.0／npm 11.6.0 与 runner 自带 `externals/node20+node24`，但 `ci.yml:804-806` 自己具名写过"交互式 PATH ≠ 作业 PATH"那枚坑）。⇒ **乙-1 落地后这两枚可能变红，红因＝构建前端那一步跑不起来**。按票 111 AC#2 已有的裁例处置＝**红名逐条登记、不撤步骤、不改产码**；⛔ 更不许用"检测不到 node 就跳过"来消红——那正是本票立案的那句谎话。

**6. 排程（写腿按住，具名理由）**：`274-r1`（唯一写面＝`scripts/build.ps1` ＋ 可能 `.github/workflows/ci.yml`）**按住等 `272-v1` 交完**——两枚都要跑整包 Go 读数（本票 AC#5 要求逐名红名册作差），⛔ 同批抢＝洗对方的数。`273-v1`（只读、零 `go test`）不受此限。⛔ 只 commit 不 push；⛔ 不动 `frontend/**` 一个字节（机主放开的只是**只读**）；⛔ 不动 `frontend/.gitignore` 的 `dist/*`（票 77 AC#1）；⛔ 不给 CI 加 `-tags winlive`。
