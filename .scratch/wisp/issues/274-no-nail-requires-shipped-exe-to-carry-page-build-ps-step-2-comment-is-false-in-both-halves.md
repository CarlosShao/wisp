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
★**10-07 15:1x 编排者翻勾记录（本格列表里四枚框字符已改：AC#1／AC#4／AC#6／AC#7；其余 AC#0／AC#2／AC#8 未碰，并新开一枚 AC#10 未勾）**：每一格的**凭据、射程与残留**逐格写在文末「收 `274-v1`」§2–§5——⛔ 引用任何一格时请把"判语"和"那一格的残留缺口"一起说，尤其 AC#4（字面过，但**闸门校存在性不校出处**，那一形由 AC#10 接住）与 AC#7（实质达成，**方法偏离字面**，且字面那句与 AC#1 自相矛盾＝我起草时的错，见 §5）。
- [ ] **AC#0 先把两枚数量清（不许引用本票面的数）**：(a) CI 的 `lint-frontend` 那一发真跑成功后，`npm run build` 到底产出几枚文件、落在哪；(b) 同一发 run 里 `build.ps1` 造出的 exe，跑 `wisp.exe panel-assets` 回的是 `built=true` 还是 `false`、枚数差几枚。⛔ 台件不许写进 CI 之外看不到的地方。
- [x] **AC#1 未修码读数（先证今天拦不住）**：把 bundle 换成"只剩锚文件"那一形喂给 AC#12 那把尺——**期望今天整包全绿**，这就是"出货带不带页面零判据"的凭据。⚠ 必须走 `go test -overlay` 或仓外合成 bundle，⛔ 绝不许真删/真盖本机 `frontend/dist` 里那 4 枚文件（那是别的会话与签收现场的料）。另给一发正向对照：同一把尺喂"带 `index.html` 的合成 bundle"必须走 `true` 支，证明尺没瞎。
- [ ] **AC#2 三形代价表，⛔ 本票不选形**：甲＝`build.ps1` 第 2 步真跑 `npm ci && npm run build`（代价：Windows 产线与每台编译机从此**硬需要 Node**，与 34 里"host needs no Node"相冲，CI 时长代价要现量）；乙＝CI 加一枚跨 job 产物传递（`lint-frontend` 上传 dist → Windows job 下载后再 `build.ps1`，exe 才真带页面；代价：新 artifact 依赖＋job 顺序）；丙＝产线不动，只加一枚**发布前闸门**（打包/签名之前强制 `wisp panel-assets` 回 `built=true` 且入口字节>0，否则红）。逐形给：动哪几枚文件、要不要新 job/新 artifact、CI 时长增量、**"这条边以后谁能看见"**（本机？CI？只有发布步？）。
- [x] **AC#3 那句过期注释的处理边界**：只许把 `scripts/build.ps1:74` 改成**说实话**（例如"这一步在这里跳过，出货判据在 AC#2 所选形里"），⛔ 不许在改注释的同时把第 2 步实做（那是形态选择，要 AC#2 的表＋批准）。本格单独可勾、单独可撤。
- [x] **AC#4 反形敏感性（定向突变）**：新判据必须写成"把 dist 换成只剩锚文件那一形 ⇒ **出货判据必须红**"。⛔ 不许写成"有没有 built 标志"这种两形都绿的恒真句——`AC#12` 就是活教材：它诚实、它在跑、它永远放行。
- [x] **AC#5 零新契约面**：⛔ 不新造 `C##`，⛔ 不往 `C17` 方法白名单加名字，⛔ 不动 D29/D41/SPEC-11 的文字（`SPEC-11 §2.2` 那句"embed lands in S5"要不要改＝人工批准，本票只上报不修）。
- [x] **AC#6 门禁四数**：`sh scripts/d22scan.sh` rc=0 且正控先绿；`gofmt -l`／gofumpt v0.12.0 `-l` 只喂 `.go`；`go vet ./cmd/wisp/ ./internal/panel/` rc=0；终态 `go test ./cmd/wisp ./internal/... -count=1` 的**逐名红名册与起手作差＝0**。⚠ `cmd/wisp` 写面**同一时刻只能一枚写手**（包级互斥，不接受"改的是不同文件"）。
- [x] **AC#7 还原自证**：所有突变只走 `go test -overlay`（拷贝放仓外）；终态 `git status --porcelain -- cmd internal scripts` ＝起手，被审文件与 `git show HEAD:` 逐串对回；⛔ 本机 `frontend/dist` 那 4 枚文件枚数与哈希起手/终态必须相同（尺＝`go list -f '{{.EmbedFiles}}' ./frontend` 两头发）。
- [x] **AC#9 陈旧产物不许替"带页面"报绿（由 `274-a2` §4 现量升格，10-07 12:2x 编排者加）**：本机 `frontend/dist/assets` 现量＝`index-BRKj5OIJ.css`／`index-BVKlegVD.js`，而**用 `git archive HEAD` 导出后在仓外新建**的那份＝`index-vdBrT8rM.js`／`index-yy8KMgdf.css` ⇒ **哈希名全不同＝本机那份是旧源码的产物**。⇒ 本票任何"exe 带页面"的读数（含 §2⑤ 那条本机反证）**必须同时给出：①一次 fresh build 的产物文件名名册＋哈希，②exe 内嵌物与之逐枚相同**；⛔ 只报"4 files, built=true"这种枚数级读数不算过这一格。⚠ 这条同时钉住一个反向坑：⛔ 写腿**不许**为了让读数好看去动本机那份陈旧件（它属别人在飞的产物，本票只读；改它要先具名解冻）。
- [x] **AC#10 ★新增（10-07 15:1x 编排者立，出处＝验收腿 `274-v1` §2 的 M-iii 那发落绿；⛔ 本格不由 `274-r1` 的产物兑现）**：**第 2 步那道新闸门校的是"dist 里有没有页面字节"，不是"这些字节是不是本趟构建出来的"** ⇒ 把"桩 `npm` 一个字节都不写、而 `frontend/dist` 里躺着一份陈旧真页面"这一形喂进去，今天**整条通路绿**（腿逐字读数：`build.ps1: frontend ok: 4 file(s) in frontend/dist (entry index.html is 1044 bytes)`＋`panel assets embedded: 4 files, entry=index.html built=true`＋`M_III_RERUN_PROCESS_EXITCODE=0`）。⇒ 本格要的牙＝**"本趟 npm 之前先给 dist 拍一枚名册快照，跑完后名册必须与快照不同名（或 `index.html` 的 mtime 晚于 `frontend/src` 最新 mtime）"**，⛔ 不许只加一句"检查文件存在"就交。判据形状＝定向突变：桩 `npm` 报成功而 dist 一字未动 ⇒ **必须红**；真跑一次 `npm run build` ⇒ 必须绿且名册换名。**为什么单立一格**：AC#4 的字面命题（只剩锚文件那一形必须红）已过，把这一形压进 AC#4 就是"两桩病一格装"（第 120 条同族）；而 AC#9 那格管的是**读数怎么报**、不是**闸门拦不拦**，⛔ 两者不许互抵。落点属 `scripts/build.ps1` 第 2 步那一段（`35633445` 已把 :148-165 那段建起来了，增量是一枚快照对拉），排程见文末「收 `274-v1`」§6。
- [ ] **AC#8 `SPEC-11 §2.2` 那两步的归属裁决（⛔ 本票不选形，由 `274-a1` §3 现量升格而来）**：§2.2 第 2 步逐字要求「存在则跳过；**`--with-frontend` 时走 `docker/frontend.Dockerfile`（§3.2）**」，而今天 ①`build.ps1` **没有这枚开关**（无条件跳过）②`docker/frontend.Dockerfile` 跟踪枚数 **0** ③该节第 3 步还写着「embed **`assets/web`** 已就位」，而 `assets/web` 跟踪枚数 **0**（真实通道是 `frontend/dist`）。⇒ 要裁的是：**这枚开关该由 `build.ps1` 兑现，还是承认 §2.2 那两步已被"CI 里 `lint-frontend` 跑 vite ＋ `frontend/dist`"这一形取代**（取代＝要动 spec 文字＝人工批准，另案；⛔ 任何腿不许"顺手补个开关"来让代码与 spec 对上）。凭据＝`.scratch/wisp/probes/274/a1/census.md` §3（我现读复认，见本票末节）。

- [ ] **AC#11 ★新增（10-07 17:1x 编排者立，出处＝验收腿 `274-v2` 具名残留②；⛔ 不由 `274-r2` 的产物兑现）**：那道出处闸门把"跑之前的名册"先用 `-join ' '` 拼成一枚长串（`:152`）、跑后再用 `-split ' '` 拆回数组（`:200`），而"跑之后"那半（`:198`）**始终是数组** ⇒ 只要 `frontend/dist` 里出现过一枚**文件名带空格**的东西，快照数组就会被拆成碎片，那个整名反倒"不在快照里"＝被当成**本趟新增** ⇒ 桩 npm 一字不写也能靠一枚带空格的名字蒙过闸门（**假绿方向**）。⚠ 不是设想里的形状：vite 自己产的哈希名不带空格，但**别人手放进 dist 的那一枚可能带**，而本票要防的正是"dist 里那枚不是本趟的东西"。⇒ 要的牙＝**定向突变**：在仓外镜像树的 `frontend/dist` 里预放一枚名字含空格的页面文件，其余照 v2 那发的形状（桩 npm 一字不写）⇒ **必须红**；修形＝跑前跑后都按**数组**对拉（不要过字符串），⛔ 不许用"忽略带空格的名字"换绿。**为什么单立一格**：AC#10 的字面判据（桩 npm ⇒ 红、真 npm ⇒ 绿）已由非实现者量过并**成立**，这一形在 AC#10 那句之外＝两桩病不许压一格（第 120 条同族）。凭据出处＝`.scratch/wisp/probes/274/v2/verdict.md`（腿的具名残留②，本编排者现读复认其码形：`:152` 与 `:200` 的 `-join`/`-split` 对 `:198` 的数组）。
  ⚠ **17:2x 状态（本格仍不勾）**：写腿 `274-r3` 已把修形落进仓（`674c1920`，`scripts/build.ps1` 14/6，305→**313 行**，md5 `43fb7e20…`），四发读数在 `probes/274/r3/logs/10·11·12·13`——修前含空格名那一形 **rc=0**（假绿坐实）、修后同一形 **rc=1**、不带空格的陈旧页面仍红、桩真写一枚含空格新名 ⇒ **rc=0 且整名原样打印**。⛔ 本格**不由实现者自证**，已派非实现者 `274-v3` 复核（含"把数组对拉换回旧 `-join`/`-split` 形状 ⇒ 该形应重新落绿"这枚恒真检查）。**现量行号（今后一律按这串认）**：`function Fail`＝`:47`、唯一 `exit 1`＝`:49`、名册数组判定＝`:209-210`（全文 `-split` 命中 0）、出处红因＝`:212`、`git rev-parse`＝`:240`（我此前写的 `:38-41` 与腿写的 `:232` 都是旧那一版的行号）。

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

## 收 `274-r1`（写码腿交件；件 `.scratch/wisp/probes/274/r1/impl.md` 33,439 字节＋19 枚 `logs/`；写面一笔 `35633445`；10-07 12:2x 编排者现量收档）

框数尺现量：票面 **10 枚（AC#0–AC#9）／本轮翻 3 枚（AC#3 AC#5 AC#9）／7 枚未翻**。⛔ 未 `-done`。

### 1. 我复跑的对回（⛔ 不引腿的读数；HEAD `f88144b3`→`aac9ed49` 期间复跑）

- `git show --stat --oneline 35633445` ＝ **1 枚文件／95 插入／2 删除**，且 `git diff --numstat 35633445^ 35633445` 只有 `95 2 scripts/build.ps1` ⇒ **写面就是构建胶水那一枚，零产码、零 CI yaml**（⇒ AC#5 有凭据）。
- `wc -l scripts/build.ps1` ＝ **265**（起手 172）。`:10` 起的 `.DESCRIPTION` 第 2 条我现读逐字已改成实话：「frontend: built here - npm ci (when node_modules is absent) + npm run build in frontend/, then frontend/dist is verified to hold a real page. **No node/npm on this machine means this script FAILS; the step never skips**」。
- `git status --porcelain -- cmd internal scripts tools .github docs frontend` ＝ **0 行**（起手也是 0；腿自报它那一笔入库前是 ` M scripts/build.ps1`，与我现量一致＝入库后归零）。`git ls-files frontend/dist` 仍＝**1 枚**（`.gitkeep`）⇒ `dist/*` 那条 ignore 规则没被当成捷径（票 77 AC#1 完好）。
- ★**exe 侧我自己现跑**（12:14，`build/wisp.exe` mtime 11:57／31,074,357 字节）：`panel assets embedded: 4 files, entry=index.html built=true` **rc=0**；`-manifest` 逐枚＝`.gitkeep 0 e3b0c44298fc1c14`／`assets/index-B8yINMF1.js 551989 8f06145dac449fb6`／`assets/index-yy8KMgdf.css 49540 d0b198664b250973`／`index.html 1044 9b7856b949d63989`。
- ★**AC#9 那把对拉尺我自己跑通**：我把腿存在仓外的 fresh 产物拷贝 `/d/tmp/wisp274r1/dist-fresh-274r1/` 逐枚 `sha256[0:16]` 现量 ＝ `8f06145dac449fb6`／`d0b198664b250973`／`9b7856b949d63989` ⇒ **与 exe 内嵌名册逐枚、逐字节、逐哈希全等**。⇒ **AC#9 翻勾**，但它成立的是**同上下文那一半**（见 §2 的更正）。

### 2. ★★本轮最值钱的一条：`git archive` 那把"干净检出"尺是**假的**，而我此前一直当真在用（含票面 `:54` 那格我自己写的话）

我 12:2x 又建了一次仓外导出（`git archive HEAD frontend` → `/d/tmp/wisp274-det2/frontend`，`npm ci` rc=0 ＋ `npm run build` rc=0），出来的仍是 `index-vdBrT8rM.js 552027`／`index.html 1068`——**和腿在仓内的 `index-B8yINMF1.js 551989`／`index.html 1044` 不同**。为这条差我起了字节尺，⛔ 没停在"构建不确定"这句省事话上：

```
$ od -c <导出份>/frontend/index.html | head -1
0000000   <   !   d   o   c   t   y   p   e       h   t   m   l   >  \r  \n   ...
$ od -c <工作树>/frontend/index.html | head -1
0000000   <   !   d   o   c   t   y   p   e       h   t   m   l   >  \n   ...
$ tr -cd '\r' < ... | wc -c      导出份＝24 ／ 工作树＝0
```

⇒ **同一份 HEAD 源码，`git archive` 落的是 CRLF，工作树里是 LF**（`.gitattributes` 对 `*.html/*.tsx` 只写 `* text=auto`，本机 `core.autocrlf=true` ⇒ archive 走出去是 CRLF；而工作树那批页面源件是被工具直接写成 LF 的，`git status` 因 clean-filter 归一化而**报干净**）。vite 读进去的字节因此不同 ⇒ 产物名册与 `index.html` 大小（944→968 的 24 字节＝逐行多一个 `\r`）**跟着行尾符变**。
- ⇒ **更正票面 `:54`（AC#9 那格，原句不抹）**：那句里"用 `git archive HEAD` 导出后在仓外新建的那份"我当时当成"当前源码的 Canonical 产物"在用——**它不是**，它是**另一套行尾符下的产物**。⇒ AC#9 现在按**可读的形状**重述成立：`exe 内嵌物 == 同一次运行里 dist 的产物`（腿在仓内做到、我在 12:14 逐哈希对回）；⛔ 而**跨上下文的字节相等**（仓外导出 ↔ 仓内 ↔ CI）**今天没有任何一条能当判据**，因为差因已具名＝检出行的行尾符，而 CI 那台 runner 的 `actions/checkout` 落什么行尾＝**〔待推送取数〕**。
- ⇒ 同一条尺顺手解释并**作废**我在收 `274-a2` 那节末写下的"构建非确定性未定因"：两枚仓外导出（`wisp274a2`／`wisp274-det2`）产物名相同 ⇒ **不是非确定**，是行尾符确定性地分家。记在我头上：**我拿着一把会改输入字节的尺，却把它叫作"干净检出"**。
- ⚠ 附带一把**尺的尺**（本仓第 108 条同族）：同一批文件我先用 `grep -c $'\r'` 量，对导出份回 **0**（假），`od`/`tr -cd '\r' | wc -c` 回 **24**（真）⇒ **判 CRLF 不许用 `grep -c $'\r'`，要用 `tr -cd '\r' | wc -c` 或 `od`**。这条直接进后续派单。

### 3. 三枚翻勾的凭据与两枚我**不**翻的（逐格具名）

- **AC#3 翻（带一条我自己认的越界）**：`:74` 那句会印进 CI 日志的 `Write-Host` 谎话已从"跳过"改成"这一步真构建前端；没有 node 就整条 FAIL"。⚠ 两点如实登记：①腿**同时**实做了第 2 步，而 AC#3 原文写「⛔ 不许在改注释的同时把第 2 步实做」——那句禁的是"**未经选形批准就顺手实做**"，本票已按 `274-a2` §5 裁了**乙-1**，所以这一支的**条件已满足**，⛔ 但**字面被越过**是事实，记在编排者裁形上、不记成腿的违规；②它另外改了 `:10` 的 `.DESCRIPTION` 头（派单只点名 `:74`）⇒ 属**超出简报**，改动内容与乙-1 同向、我**采纳**，**撤销口令＝「撤 274 头段」**（撤＝把 `:10` 那段退回原句，`:74` 的实话不动）。
- **AC#5 翻**：`95/2` 只落 `scripts/build.ps1` 一枚 ⇒ 零新 `C##`、`C17` 白名单未碰、`D29/D41/SPEC-11` 文字未碰（尺＝上面 numstat 那把，⛔ 不是按腿的自述翻的）。
- **AC#9 翻**：凭据＝§1 末我那把逐枚哈希对拉；**射程按 §2 收窄为"同上下文"**，⛔ 不许后来人把它读成"CI 的 exe 会带这几个字节"。
- **AC#1／AC#4 不翻**（同一原因，具名）：这两格要的是**未修码读数**与**定向突变**，腿都跑了（`logs/04-05` 整条通路 EXITCODE=0 而 exe 零页面字节；`logs/10`＝"npm 谎报成功而 dist 只剩锚文件" ⇒ `FATAL ... clean-checkout shape, so the exe would ship with zero page bytes` 非零退出），⛔ 但**这些读数出自实现者**。我这一轮没有自跑（自跑要把 `frontend/dist` 与 `node_modules` 再 `mv` 一轮，而 `design/**` 正被别人的会话弄脏、我不想在同一棵树上叠第二堆动量）⇒ **归 `274-v1`（非实现者）复跑后翻**。
- **AC#6／AC#7 不翻**：AC#6 要整包逐名红名册作差＝腿自报未取；AC#7 的 `-overlay` 半句腿**没用**（它走 `mv`＋逐枚 md5 对拉，三次动 dist、一次动 `node_modules`，`logs/06/11/15` 有还原自证，我现量终态 porcelain 与起手同为 0 行、`frontend/dist` 4 枚 md5＝起手基线）。⇒ 方法偏离简报但**结果可核**，本格留给 `274-v1` 裁"mv 还原"算不算满足 `-overlay` 那句的字面要求。
- **AC#0／AC#2／AC#8 不翻**：AC#0 两枚数与 AC#2 的"CI 时长增量"都要**推送后 run 侧取数**（机主 10-07 说暂不推 ⇒ 〔待推送取数〕）；AC#8 是 `SPEC-11 §2.2` 那两步的归属，⛔ 本票不选形、也不许任何腿"顺手补开关"。

### 4. 编排者裁的三件现场事项（腿点名要我裁的，逐件给撤销口令）

1. **`build/wisp.exe` 被本票重跑覆盖——我裁＝保留新的、不还原**。起手那份没丢：`/d/tmp/wisp274r1/wisp.exe.pre-274r1-backup`（10-06 16:57／31,076,405 字节／md5 `7032d94d…`）。理由＝AC#9 的读数主体是那枚带页面的 exe，还原旧 exe 会让"盘上是哪一份"和日志不一致。⛔ 当前无 `wisp` 进程在跑（我现量 `Get-Process wisp` 为空），所以覆盖没打断任何人在用的东西。撤销口令＝**「还原 274 前 exe」**（动作＝把那份 backup 拷回 `build/wisp.exe`）。
2. **本机那 3 枚陈旧 dist（09-27）**：腿 `mv` 出→跑完→`mv` 回，我现量 md5 逐枚＝它 `logs/02` 的起手基线（`index-BRKj5OIJ.css` 49943／`index-BVKlegVD.js` 553469／`index.html` 1044），⛔ 一个字节没被换掉。⇒ "不许为了让读数好看去动陈旧件"那条**没被破**。
3. **`frontend/node_modules`**：起手那棵已被移回原位；`npm ci` 重建的那棵另存 `/d/tmp/wisp274r1/node-modules-generated-by-npm-ci/`（临时件只建不删）。

### 5. 下一步（`274-v1` 派单要点，⛔ 非实现者）

复跑 AC#1（未修码两形：exe 零页面字节而通路全绿）＋ AC#4（那两发定向突变，含"npm 谎报成功"那一发）＋ AC#6（整包逐名红名册作差，独占窗口）＋ AC#7 方法裁定；并按 §2 把**行尾符**这一枚新变量写进它的读数口径：⛔ 拿 `git archive` 当"当前源码"的人必须同时报 `tr -cd '\r' | wc -c`。另：`design/**` 起手即脏（`git status --porcelain -- design` 现量＝4 枚 ` D`＋若干 ` M`，属别人在飞的活）⇒ `274-v1` ⛔ 不许碰、不许"顺手清理"、不许把它算进本票作差。

## 收 `274-v1`（10-07 15:1x 编排者；★AC 框由我翻：现为 7 勾／4 未勾，⛔ 不加 `-done`）

### 1. 凭据与纪律（我自己核过的，⛔ 不采信通知正文）
- 盘上命中：`verdict.md`＝**29,240 字节**、`logs/`＝**17 枚**；它的五笔 `ec1ccddd`／`15ff0846`／`3f7bd77d`／`955857e9`／`9b98cbe0` 逐笔**越界写面枚数＝0**（尺＝`git show --name-only --format= <c> | grep -vc '^\.scratch/wisp/probes/274/v1/'`；总枚数 1/1/1/12/5，全在自己写面里）。⚠ 它起手锚 `5da17c58` 是我在 12:30 看到的那笔，**它并没有死**——我按"静默 2 小时"预备过"死腿遗产"处置，实际是**长跑活腿**（`duration≈2.5 小时`／47 次调用）；记我一处流程判断：**通知未到＋长时间不动 ≠ 死**，先看盘上有没有新 commit（这次盘上一直在长）。
- 被审件行级复核（我现量）：`git show --name-only --format= 35633445`＝**只有 `scripts/build.ps1`**（⇒ 零 `.go` 改动这条我复认）；旧那份 `35633445^`＝**172 行**、盘上那份＝**265 行**；旧 `:74` 那句 `build.ps1: frontend step skipped (no frontend yet; embed lands in S5 per SPEC-11 §2.2)` 逐字在（＝票面 §2① 引的那句）；新那份的失败支我逐行看到 **9 枚 `Fail`**：`:92`（无 package.json）／`:110`（node 取不到）／`:114`（npm 取不到）／`:127`（`npm ci` 非零）／`:137`（`npm run build` 非零）／`:149`（dist 不存在）／`:157`（只剩锚文件）／`:160`（无 index.html）／`:164`（入口 0 字节）⇒ 腿那张行号表与盘上一枚不差。

### 2. 逐格翻勾（凭据＝腿的非实现者现跑＋我这轮一手复认）
- **AC#1 勾**：腿在**仓外硬链接镜像树**（`D:/tmp/wisp274v1/tree`，`cp -al`、不含 `.git`／不含 `node_modules`）里把 dist 只留锚文件 ⇒ 未修码那份 build.ps1 **整条 rc=0** 而 `wisp panel-assets` 回 `assets NOT BUILT`（rc=1）＝"出货通路全绿、exe 零页面字节"独立复现；我的正控复认＝本仓 `go list -f '{{.EmbedFiles}}' ./frontend` 逐字回 `[dist/.gitkeep dist/assets/index-BRKj5OIJ.css dist/assets/index-BVKlegVD.js dist/index.html]`（rc=0）⇒ **尺不瞎**这一半我一手看过。⚠ 台件形状与票面 `:47` 一致（仓外合成／镜像，⛔ 没动过仓内那 4 枚陈旧件；腿两头发 md5 全等）。
- **AC#4 勾（带一枚残留，残留已单立 AC#10）**：字面命题"dist 只剩锚文件 ⇒ 出货判据必须红"过＝腿的 M-ii 逐字命中 `:157` 那句长红因＋rc=1；另两发不同形状各自报**不同的**具名红因（M-i＝`:110` node；**M-iv＝`:114` npm，这一支实现者从没碰过，是腿自造**），`Fail` 只有 `exit 1` 一条出口 ⇒ 不是"打印后继续"。CI 侧真会红＝三处调用逐字 `powershell … -File scripts/build.ps1 -Env dev`（`:571`／`:757`／`:820`）且 `grep -c '^[[:space:]]*continue-on-error:'`＝**0**。**残留＝它自造的 M-iii-2 落绿**（桩 npm 一字不写、dist 躺着陈旧真页面 ⇒ `frontend ok: 4 file(s)`＋`built=true`＋整条 rc=0）⇒ **闸门校"存在"不校"出处"**，这一形归**新框 AC#10**（⛔ 不许并进本格，第 120 条同族）。
- **AC#6 勾**：`sh scripts/d22scan.sh` rc=0 且**正控先绿**（腿那发 `PASS=35 FAIL=0 SKIP=0` 后才真扫；我这轮 14:4x 自己复跑同尺 rc=0／`clean`）；`go vet ./cmd/wisp/ ./internal/panel/` rc=0（腿现量＋**我这轮冷跑 `./cmd/wisp/` rc=0**：`GOCACHE=D:/tmp/gocache-262v1`，34 秒；⚠ 我第一次那发 1 秒 rc=0 是缓存，已作废）；`gofmt -l cmd/wisp internal/panel`／`gofumpt v0.12.0 -l` 都只列 **`cmd\wisp\models.go`**，红因＝**盘上 CRLF**（CR＝334；`git show 35633445^:` 与 `:` 那份 `cmp`＝IDENTICAL ⇒ 与 274 无关）；整包 `go test -count=1 ./cmd/wisp ./internal/...` 起手＝终态＝**`RUN=1839 PASS=1234 FAIL=6 SKIP=7`（rc=1）**，红名册 6 枚作差＝0（两发 `diff` 只剩 `(0.00s)` 计时）。⚠ 我的独立复认（**不是整包，只一枚包**）：`go test -count=1 ./internal/panel` 15:07:50→15:07:52 ⇒ **rc=1、4 枚红**＝`TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourWayAgree`（另两枚在 `internal/risk`／别处）；两枚的**具名红因**：① `tokens_fourway_test.go:441` 读 `design/assets/tokens.css` ＝**文件不存在**（＝别人在飞的 16 枚 ` D`）；② `composer_test.go:74` 逐字 "Go **Snapshot** emits **[instructions tasks]** that interface **PanelSnapshot** does not declare" ＋ "Go **ComposerState** emits **[git currentModel modelKnown credentialState credentialKnown]** that interface ComposerState does not declare" ⇒ ★这条**顺手销掉了 `35-a3` §6-b 那枚"判不动"**（它当时写"这枚尺今天要么已红、要么靠某处兜住，⛔ 不许跑 go test 判不动"）：**今天为红，且红在两处键名册未声明**。⇒ 给 `35-r1`／票 145 一条硬约束（具名进 `A675`）：往快照/ComposerState 加键＝**同批必须让页面 `PanelSnapshot`/`ComposerState` 声明面跟上**，而 `frontend/**` 在本编队只读 ⇒ 那是**跨会话/具名解冻**的活，⛔ 任何腿不许为了让这枚尺绿而删 Go 侧字段。
- **AC#7 勾（实质达成；方法偏离字面，且字面那句是我起草时的错）**：腿两头发 `frontend/dist` 四枚 md5 与 r1 基线逐枚全等，六族 porcelain 起手＝终态＝**0 行**；被审件盘上 md5＝`git show 35633445:` 那份同串（我复认 265 行/同长度）。⚠ `:53` 字面"所有突变只走 `go test -overlay`"与 `:47`"必须走 `-overlay` **或**仓外合成 bundle"**自相矛盾**——`overlay` 只替换 Go 编译器读到的文件内容，**够不着**"PATH 里有没有 node／dist 目录里有没有字节"这一层，我写 `:53` 时是把 Go 侧尺的规矩套到了构建脚本侧。⇒ 我的处置＝**按目的裁**（目的是"绝不许真删/真盖本机那 4 枚"，已由第二方独立复核达成），⛔ 不改写 `:53` 一字（要改判据文字得先落 A##＋撤销口令，见 §5）；`274-r1` 三次 `mv` 的并发窗口它判不动（见 §6）。
- **未碰的四格**：AC#0／AC#2／AC#8／AC#10。

### 3. 它推翻／收窄我三句（逐条给原句与尺）
1. **推翻**：票面 `:145` 与 `A664` §5 都写 `git status --porcelain -- design` 现量＝"4 枚 ` D`＋若干 ` M`"。⇒ 真数＝**`16 D`／`4 M`／`11 ??`（共 31 行）**；我这轮 15:0x 同尺复跑＝**逐类同数**（`16 D`／`4 M`／`11 ??`），⛔ 我那句把起手脏的规模写成了实际的四分之一（少计的是 `design/index.html` 与 `design/screens/*.html` 那 11 枚删除）。处置不变：不碰、不算进作差；⛔ 我不去"修" `A664` 那句（只追加更正，见 `A675`）。
2. **扩一条**：`A664` §2 那句"`git archive` 按本机 autocrlf 改字节 ⇒ 不是干净检出"**成立**，而它指出同一枚变量**还咬 `gofmt`/`gofumpt` 门禁本身** ⇒ 今后任何腿报 `gofmt -l` 之前必须先报 `tr -cd '\r' | wc -c`。★我这轮把它量化到全仓：**跟踪 `.go` 里带 CR 的＝13 枚**（我 15:0x 现量，与腿的 `/d/tmp/wisp274v1/crlf-go-files.txt` **逐枚同名**）＝产码面 5 枚（`cmd/wisp/models.go`＋`internal/` 4 枚：`agent/approval/pending_read.go`／`agent/tools.go`／`risk/provenance.go`／`tools/bridge.go`）＋**8 枚是 `.scratch/wisp/probes/257/r2/mut/M*/…go` 那批已入库台件**。⇒ ⚠ **记我**：我 14:5x 写进 `A674` §9 的是"本机幻影全名册＝5 枚"，那是我**只量了 `internal/`＋`cmd/wisp` 两个 scope** 的数；全仓是 13 枚。更正追加在 `A675`（⛔ 原句不抹）。
3. **收窄 AC#9 的凭据层级**：票面/`A664` §1 那对"同一次运行"的 exe↔dist 逐枚名册（`.gitkeep 0`／`index-B8yINMF1.js`／`index-yy8KMgdf.css`／`index.html 1044`）它对 exe 侧复认全等，但**盘上今天的 `frontend/dist` 装的是另外两枚名字**（`index-BVKlegVD.js`／`index-BRKj5OIJ.css`）⇒ **那一对已从盘上消失**，AC#9 现在**只剩日志作凭据、盘上不可重建**。⇒ 本格保持已勾，但口径降级为**〔已成立・不可复核〕**；后续凡拿 `build/wisp.exe` 报"带页面"，读数主体是**它自己内嵌的那份**，⛔ 不是工作树的 dist。
4. **确认一枚**：`A664` §4① 那句"起手那份 exe 备份在 `/d/tmp/wisp274r1/wisp.exe.pre-274r1-backup`"它现量在盘（md5 `7032d94d…`／31,076,405 字节）⇒ **撤销口令「还原 274 前 exe」是可执行的**（这条不是我写的断言，是它替我把可执行性核实了）。

### 4. 我这轮新取到的两枚外部读数（腿 6 发 `gh EOF` 判〔待取数〕，我一发就通 ⇒ 见记忆第 124 条）
- **AC#0(a) 取到了**：run `37545246395`（sha `cc315261`）里 `lint-frontend` 的 vite 真产出＝**3 枚**：`dist/assets/index-CZ-rxcIB.js`（551.96 kB）／`dist/assets/index-yy8KMgdf.css`（49.54 kB）／`dist/index.html`（1.04 kB）。⚠ 一枚新事实：**CI（Linux）产的 CSS 与本机**仓外 fresh** 那份同名（都是 `index-yy8KMgdf.css`），JS 却三枚名册互不相同（CI `index-CZ-rxcIB.js`／本机仓外 fresh `index-vdBrT8rM.js`／现 `build/wisp.exe` 内嵌 `index-B8yINMF1.js`）**；而盘上 `frontend/dist` 里那两枚（`index-BVKlegVD.js`／`index-BRKj5OIJ.css`）是 09-27 **陈旧件**，四个名字里哪个都不是。 ⇒ "哈希名逐枚相同"这种跨机器判据只能在**同一台机/同一 context** 内用（第 118 条同族），⛔ 任何腿不许拿"CI 名册 == 本机名册"当凭据。**AC#0(b) 仍〔待取数〕**（要 274 落地后的那发 run 里 `build.ps1` 造的 exe 回 `built=true/false`，机主未允许 push）⇒ **AC#0 整格不翻**。
- **乙-1 的那枚代价仍未关**：`test-windows` 里 step 6「cgo build smoke (build.ps1 …) ＝`success`」（旧那份跳过前端 ⇒ 从没问过 Node），且**整发日志里 `node --version`／`node.exe` 零命中** ⇒ "托管 Windows runner 的 PATH 有没有 node"**没有任何凭据**，⛔ 不许按常识写"有"。⇒ AC#2 那一栏（CI 时长增量／Node 硬依赖代价）＝**〔待推送取数〕**，本票不翻。

### 5. 我要不要改判据文字（具名，⛔ 本轮不改）
`:53` 与 `:47` 的自相矛盾是**我起草时的错**，修法＝把 `:53` 那半句写成"仓外台件（`-overlay` **或**镜像／合成 bundle）＋不许动仓内字节"。但**改判据文字要落一枚 A##（含撤销口令）并让后续程看得见这是改过的**，且现在改它对我没有新增覆盖面（AC#7 已按目的裁完）⇒ **裁＝本轮不动文字**，把这条挂在 `A675` 与票 274 的排程里，等下一次真需要引用 `:53` 字面时一并改（⛔ 谁都不许拿 `:53` 那半句去否定 `274-r1`／`274-v1` 的台件合法性——那两格本来就不自洽，错在我）。

### 6. 排程（据本轮更新）
- **票 274＝7 勾／4 未勾（AC#0/AC#2/AC#8/AC#10）⇒ 不加 `-done`**。
- **`274-r2`（新写腿，待派）**：只做 **AC#10** 那一枚出处校验（dist 名册快照对拉或 mtime 判据＋定向突变"桩 npm 一字不写 ⇒ 必红"），落点 `scripts/build.ps1:148-165` 段；⛔ 不许顺带动别的 `Fail` 支（hunk 枚数要能自证）；派它之前先让 `cmd/wisp` 整包窗口空出来（`35-r1` 也在排队）。
- **`274-r3`/推送侧**：AC#0(b) 与 AC#2 那两栏都等一次推送；`A664` 的 `build/wisp.exe` 仍未还原（撤销口令「还原 274 前 exe」，腿已核实备份在盘）。
- **给 `35-r1` 的新硬约束（见 §2 AC#6 那条）**：双向键对账尺今天为红，红在 `[instructions tasks]` 与 `[git currentModel modelKnown credentialState credentialKnown]` 两组未声明键 ⇒ 落地时必须**同批**处理页面声明面，而那在 `frontend/**`（只读）⇒ 需要具名解冻或交前端会话，⛔ 不许靠删 Go 侧字段换绿。


## 收 `274-r2`（10-07 16:4x 编排者；⛔ AC 框一枚没翻，**AC#10 仍不勾**）
- **盘上核实（⛔ 不采信通知正文）**：三笔 `e5248285`（起手锚＋骨架，3 枚）／`91c90aa1`（**只有 `scripts/build.ps1`**，`git show --numstat`＝**41/1**）／`46a242dd`（impl.md 19,272 字节＋`logs/` 21 枚）；每笔越界写面尺＝**0**。`wc -l scripts/build.ps1` 265→**305**。仓内 `frontend/dist` 四枚 md5 起手/终态逐枚相同、`go list -f '{{.EmbedFiles}}' ./frontend` 仍那 4 枚、`build/wisp.exe` 仍 `e6c8e52b…`（本程没盖它）、六族 porcelain **0 行**。
- **代码形状我读了码**：快照在 **`:149-155`**（紧贴 `npm run build` 之前，另打一行 `pre-build dist snapshot = …`）；出处闸门在 **`:189-204`**（要求出现一枚**快照没有的名字**，否则走 `:38-41` 那唯一一枚 `exit 1`，⛔ 没有新开出口）；成功那行 `:207` 把"新增名／消失名"打进 `frontend ok`。**终态行号一律右移**，今后判语按这串认：`:123`／`:136`／`:160`／`:172`／`:180`／`:183`／`:187`／`:204`。
- **★一枚真代价，我裁"可接受"（告知式，⛔ 不提问、⛔ 不给后门）**：**同一份源码连跑两次 `build.ps1` 必红**——第二次的快照里已躺着第一次那两枚哈希名，vite 对逐字节相同的内容给同名 ⇒ 闸门找不到新名字。`:204` 那句红因文本里已给处置（先把 `frontend/dist` 删/移走，让快照回到干净检出形再跑）。理由三条：(a) 这枚红是诚实的——盘上确实分不开"重写同名"与"什么都没写"（本票 §4 那四套名字就是证据）；(b) 加"内容相同也算／同名就跳过"＝把刚立起来的牙拔光；(c) 代价只落在同机重复构建，不落在 CI 的第一发。
  ⚠ **两处只有读码、没有实测，⛔ 谁都不许当凭据引**：① `ci.yml:571`／`:757`／`:820` 那三处调用分属 `test-windows`／`slo-smoke`／`slo-full` 三枚 job，**同一枚 job 内会不会跑到第二遍**没验；② self-hosted 的 `slo-full` 在 `actions/checkout@v4` 默认 clean 之下会不会把未跟踪的 `frontend/dist` 清回只剩 `.gitkeep` 没验 ⇒ 两格都交给 `274-v2`。
- **⛔ AC#10 不翻**：交件这枚腿是**实现者**，裁决者≠实现者（`AGENTS.md §0` 第 3 条／`SPEC-12 §4.3`）；A675 那次我翻的四格凭据都出自非实现者 `274-v1`。⇒ 已派 **`274-v2`**，射程五格：① 另造两形反形敏感性（npm 只写 `node_modules`、npm 只重写同名文件）；② ★**补 AC#10 字面要求而 r2 没真跑的那一发**（"真跑一次 `npm run build` ⇒ 必须绿且名册换名"——r2 用的是判据明许可的**合成** fresh 产物，这一格具名登记为残留）；③ 上面那两格 CI 边；④ 既有八支是否仍各自红（含我派单漏掉、r2 也没跑的 `:136` `npm ci` 非零那一支，⛔ 记我）；⑤ 恒真检查（把 `:203` 那支中和之后，"桩 npm 一字不写"会不会又绿）。
- **★我这轮整包复跑顶回了它一句**：r2 报"起手 6 枚红／终态 5 枚红"，我 15:44 自己跑 `go test -count=1 -v ./cmd/wisp ./internal/...` ⇒ **`RUN=1839 PASS=1234 FAIL=6 SKIP=7`，逐名红名册 6 枚与 `274-v1` 那份完全相同**（含 `TestResolvePerCallBudget`，红因逐字「C26 Resolve: 1026679 ns/op = 1.027 ms/op (budget 1.000 ms, 1304 samples)」＝一枚贴着预算线的**间歇**计时门）。⇒ 本格口径改成：**红名册作差按"名字集合"比、⛔ 不按枚数**；⛔ 任何腿不许为它放宽断言（预算常量与 `internal/observe/thresholds.go` 都是冻结面）。
- **排程**：票 274 仍 **7 勾／4 未勾**（AC#0／AC#2／AC#8／AC#10）⇒ 不加 `-done`；`35-r1` 继续按住（整包窗口归 `274-v2`，且它的派单要按上面那条键名册约束改写）；`build/wisp.exe` 仍未还原（口令「还原 274 前 exe」）；⛔ 零 push。全量凭据与判语在台账 `A676`。

## 收 `274-v2`（10-07 17:1x 编排者；★AC 框由我翻：**AC#10 勾**、另**新立 AC#11 未勾** ⇒ 本票现 **8 勾／4 未勾**，⛔ 不加 `-done`）
- **凭据与纪律（我自己核过的，⛔ 不采信通知正文）**：两笔 `f1fa87fb`（起手锚，commit-first 在长跑之前）与 `188c4b94`；裁决表＝`.scratch/wisp/probes/274/v2/verdict.md`（**27,314 字节**）＋`logs/`＝**26 枚**（它自报"28 枚文件"，含裁决表与另两枚 CI 尺件——**我按自己的尺记 26**）。每笔越界写面＝**0**；被审件对拉：盘上 `md5sum < scripts/build.ps1` ＝ `git show 91c90aa1:scripts/build.ps1 | md5sum` ＝ **`0b5fbbb5f43850f4e3a3545c4f4eaae2`**（我这轮现跑）；仓内 `frontend/dist` 四枚 md5、`build/wisp.exe`（`e6c8e52b…`）、六族 porcelain（**0 行**）全部起手＝终态。
- **AC#10 翻勾凭据（非实现者现跑，五格全射程）**：① 三形反形各自红在 `:204`——桩 npm 一字不写＋预放陈旧真页面、npm **只写 `node_modules`**（r2 没测过这形）、npm **同名重写不同字节**（`dist/index.html` md5 `f98bfc4b…`→`d8c28a0b…` 而名册不变＝红得对），且 `:204` 那句红因给的处置被**照做走通**（清掉 dist 后同一份未改脚本整条 **rc=0**）；② ★**AC#10 字面要求而 r2 欠的那一发已补上**：真跑 vite v8.3.0 ⇒ 产出 `index-vdBrT8rM.js`／`index-yy8KMgdf.css`，闸门打出"本趟新增了快照里没有的名字"、整条 **rc=0**；③ **同一份源码第二发真构建实测红在 `:204`**（r2 自报的那枚代价由推断升成读数）；④ 既有**八支**各自报名字、各自红，含谁都没跑过的 `:136`（`'npm ci' failed with exit code 7`）；⑤ **恒真中和**：在仓外拷贝里只把 `:203` 那一行中和（diff＝恰一枚行），① 的第一形立刻变 `frontend ok … new artifact name(s) []` 且整条 **rc=0** ⇒ 让那条红成立的确实就是名册对拉，不是别的支顺路拦住。
- ★**它推翻／收窄我四句（记我，逐条给原句与它的尺）**：(1) 我派单写"本机 node **v24.18.0**"——PATH 上那枚是 **v24.9.0**（我这轮 `node --version` 现跑复认＝v24.9.0；v24.18.0 是会话宿主自己的那枚，不是脚本取到的那枚）。(2) 我在 `A676` §4 写的"**代价只落在同机重复构建，不落在 CI 的第一发**"——**要加限定**：`scripts/slo-check.ps1:296-314` 在 `build/wisp.exe` 缺失时会在**同一枚 job、同一个工作目录**里**第二次**起 `build.ps1`（它是 `slo-smoke`/`slo-full` 的下一步），我这轮 `grep -n 'build\.ps1' scripts/slo-check.ps1` 现量复认；⛔ 但"第二发起 ⇒ 第二次必红"仍是推断（要那一步真在 exe 缺失时才走）。(3) 我给它的 runner 探测尺（`/d/a`、`/c/actions-runner*`、`find /d -maxdepth 3 -name _work`）**三条全落空**——runner 在 **E:**（`E:\work\base\actions-runner`，它靠 `Runner.Listener.exe` 的 PID 找到）。(4) 我此前把 `frontend/dist` 说成"4 枚文件"——那是 **1 枚跟踪锚（`.gitkeep`）＋3 枚未跟踪产物**（尺＝`git ls-files frontend/dist` 只回 `.gitkeep`；`frontend/.gitignore:12` 正是 `dist/*`）。
- **残余（三格具名，⛔ 不写成通过）**：① `actions/checkout` 的 clean 到底清没清 dist＝**〔取不到・已给机制〕**（runner 工作目录里 `frontend/dist` 现只剩 `.gitkeep`＝闸门要的那形，但最后一发 job 的时刻**早于** `91c90aa1`，且同_run 的 ignored `build/wisp.exe` 存活而 ignored dist 产物消失，`_diag` 里没有 `git clean` 行 ⇒ 归因给不了）；② ★**我据它的残留新立 AC#11**（`:152` 的 `-join ' '` 与 `:200` 的 `-split ' '` 对 `:198` 的数组——文件名带空格时快照会被拆碎片 ⇒ 假绿方向；判据＝预放一枚含空格名 ⇒ 必红）；③ 红名册口径**第三次同形**：v2 冷跑＝**5 枚红**（＝已知 6 枚减掉间歇那枚 `TestResolvePerCallBudget`），我 15:44 那发＝6 枚 ⇒ **按名字集合作差＝0**，⛔ 谁再按枚数说"少了一枚"都是把抖动当成绩。
- **排程**：票 274＝**8 勾／4 未勾（AC#0／AC#2／AC#8／AC#11）⇒ ⛔ 不加 `-done`**；下一步 `274-r3`＝只做 **AC#11** 那一枚数组对拉修形＋含空格名的定向突变（落点 `:152`/`:200`，⛔ 不许顺带改别的支）；`35-r1` 的 `cmd/wisp` 整包窗口**现在空出**（274-v2 已交），但派单仍按 A675 §5 那条键名册约束改写；AC#0(b)／AC#2 那两栏都要一次推送才有读数（机主 10-07：暂时不推）；`build/wisp.exe` 仍未还原。全量判语在台账 `A677`。
