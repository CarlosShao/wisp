# `274-v2` 裁决表 —— 票 274 AC#10「本趟构建出处」闸门（腿＝非实现者，攻 `274-r2` 落进仓里的那道牙）

被审件＝`scripts/build.ps1`（盘上 305 行，blob md5 `0b5fbbb5f43850f4e3a3545c4f4eaae2`，与 `91c90aa1` 逐字节相同，见 §4）。
判据原文＝`.scratch/wisp/issues/274-no-nail-requires-shipped-exe-to-carry-page-build-ps-step-2-comment-is-false-in-both-halves.md:56`（AC#10），
逐字抽出的两半：**「桩 `npm` 报成功而 dist 一字未动 ⇒ 必须红；真跑一次 `npm run build` ⇒ 必须绿且名册换名」**，外加「名册必须与快照不同名（**或** mtime 判据）」的择一形。
本腿台件全部在仓外 `D:/tmp/wisp274v2/**`；**仓内一字节没动**（§4）。**本腿 AC 框一枚没翻。**

## §0 起手锚（逐字在 `logs/00-anchor.txt`，commit `f1fa87fb`，早于任何长跑命令）

`date` = `Wed Oct  7 16:43:08 CST 2026`；`git log -1 --oneline` = `d8299ee7 A676 收 274-r2（…）`；分支 `dev`；
`git status --porcelain -- cmd internal scripts .github docs frontend build` = **0 行**（终态仍 0 行，§4）。
起手 dist 四枚 md5 ＝ `d41d8cd9…`／`db4db7a2…`／`70128a3d…`／`f98bfc4b…`，`build/wisp.exe` = `e6c8e52bb14f15a7983bbc6b093e6058`。

环境现量（⛔ 与编排者简报不符处已具名，见 §3-1）：node **v24.9.0**（`D:\work\server\node14\node.exe`）、npm **11.6.0**、go **go1.27.1**、gcc **16.2.0**（`E:\work\base\msys64\mingw64\bin`）。

台件搭法（承 `274-v1` §1 第 1 段的硬链接镜像形，⛔ 不引它的读数）：
- `real/` ＝ `git clone --depth 1 --no-hardlinks`（仓外，带自己的 `.git`）＋ `cp -al` 只读 `third_party/`＋ **`cp -a` 真拷贝 `frontend/node_modules`（167M）** ⇒ ② 用；
- `stub/`、`green/` ＝ `cp -al` 逐条顶层目录（⛔ 不含 `.git`、⛔ 不含 `build/`），`frontend/node_modules` 与 `frontend/dist` 一律 `rm` 后**真建** ⇒ ①④⑤ 用。
  `build/` 必须整目录不进镜像：`go build -o build\wisp.exe` 是截断写，硬链接会**直接盖到仓内那枚 `build/wisp.exe`**。
- ⛔ `node_modules` 全程没走硬链接（简报明令），跑完回查仓内三枚样件 md5 全 `OK`（`logs/24-final-proof.txt`）。

---

## §1 逐格裁决

### ① 复现它的突变（三形）＝**成立**

尺与原始读数（本腿自跑，⛔ 未引 `r2/logs`；`logs/20-stub-battery.txt`、`logs/22-battery2.txt`，行号如下）：

| 形 | 台件 | rc | 逐字红因（截取）／绿 | 凭据 |
|---|---|---|---|---|
| **(a)** 桩 npm 报成功、一个字节不写；dist 预放陈旧真页面（仓内那 4 枚的**真拷贝**） | `stub/`，`pre-build dist snapshot = 4 file(s) [.gitkeep assets\index-BRKj5OIJ.css assets\index-BVKlegVD.js index.html]` | **1** | `FATAL: … frontend/dist holds no artifact name outside the snapshot taken before this run's npm … Named cause: unprovenanced page bytes - nothing in frontend/dist was produced by this build, so wisp.exe would embed an earlier build's artifacts.` ＋ 跑后 dist 四枚 md5 **逐枚未变** | `logs/20-stub-battery.txt:29-38`（快照行 `:29`、FATAL `:32`、rc `:33`、md5 `:35-38`） |
| **(b)** 桩 npm 只往 `frontend/node_modules` 写一枚文件，dist 一字不动（**r2 没测过这一形**） | `stub/`，同上陈旧页 | **1** | 同一句 `Named cause: unprovenanced page bytes`（`logs/22-battery2.txt:29`）；**且这一发 npm 确实干了活**：跑后 `node_modules_marker=marker-from-this-run`（`:36`）而 dist 四枚 md5 未变（`:32-35`） | `logs/22-battery2.txt:10-36` |
| **(c)** 桩 npm 把 dist 里**同名**文件重写成不同内容 | `stub/`，同上陈旧页 | **1** | 同一句具名红因（`:57`）；`dist/index.html` 的 md5 跑前 `f98bfc4b…` → 跑后 **`d8c28a0befec90aa00f0c71ab6e29fc2`**＝**字节真的变了**，名册逐字未变（`:42` vs `:63`）⇒ 盘上确实分不开「重写同名」与「什么都没写」，红得**对** | `logs/22-battery2.txt:38-64` |

**(c) 那一形给的处置走不走得通＝走得通（实测，不是读码）**：`green/` 树里先跑 (c) 同一形 ⇒ rc=1、红因逐字同上（`logs/22-battery2.txt:182`，rc 行 `:183`）；
然后**照 `:204` 红因文本逐字写的处置**做一遍——「delete or move `frontend/dist`（跟踪锚 `frontend/dist/.gitkeep` 由 git 重新生成）让前置快照回到 clean-checkout 那一形，再跑」：
`rm -rf frontend/dist` ＋ 在**仓外**那份带自己 `.git` 的树里 `git checkout -- frontend/dist` ⇒ `GIT_CHECKOUT_RC=0`、名册回到 `dist/.gitkeep` 一枚；
再跑**同一份未改的 `build.ps1`**（npm 桩这次产新名）⇒ **`CASE_D2_rerun_after_remedy_PROCESS_EXITCODE=0`**，整条通路走完（`go build ok`／DLL 并置／SHA256SUMS／`wisp doctor: PASS`），
`go list -f '{{.EmbedFiles}}' ./frontend` ＝ `[dist/.gitkeep dist/assets/index-AAAA1111.js dist/assets/index-BBBB2222.css dist/index.html]`。凭据 `logs/22-battery2.txt:167-末段`。
⇒ **这一形红得对、红因里的处置真能走通**；本腿⛔没要口子，也没有以任何形式给「同名重写算绿」背书。

### ② 补 AC#10 字面要求的**真构建那一发**＝**成立（真 vite，不是合成产物）**

选**哪条腿＋理由＋风险**（简报两选一，我取并具名）：走 **（乙）`cp -a` 真拷贝 `frontend/node_modules`**，不走 `npm ci`。
理由：①`npm ci` 要联网取 registry，本机零推送窗口里多一次外部依赖不值；②AC#10 这一格要的是「**真跑 `npm run build`** ⇒ 名册换名」，`node_modules` 从仓内**真拷贝**（不是硬链接）恰好把「装的东西＝本仓现装的东西」钉死，风险只剩一条：**这一形走的是 `node_modules EXISTS → npm ci SKIPPED` 分支，不校验 lockfile**（build.ps1 自己在 `:130` 逐字承认了这句），所以 `npm ci` 那一支我另在 ④ 用桩打。⛔ 全程没有对 `node_modules` 用 `cp -al`。

原始读数（`logs/10-real-npm-build.txt`，行号）：
- `:11` `node = D:\work\server\node14\node.exe (v24.9.0); npm = D:\work\server\node14\npm.cmd (11.6.0)`
- `:13` 快照行逐字：`build.ps1: frontend: pre-build dist snapshot = 1 file(s) [.gitkeep] (this run's npm must produce an artifact name outside it).`（跑前该树 dist 只有跟踪锚，`:4-5` `ls` 现量）
- `:17` `> tsc -b && vite build`、`:19` `vite v8.3.0 building client environment for production...`、`:24-26` 真产物 `dist/assets/index-yy8KMgdf.css`／`dist/assets/index-vdBrT8rM.js`
- `:34` 绿因逐字：`build.ps1: frontend ok: 4 file(s) in frontend/dist (entry index.html is 1068 bytes): .gitkeep=0 assets\index-vdBrT8rM.js=552027 assets\index-yy8KMgdf.css=49540 index.html=1068 || provenance: this run's npm added new artifact name(s) [assets\index-vdBrT8rM.js assets\index-yy8KMgdf.css index.html] over the pre-build snapshot 1 file(s) [.gitkeep]; gone since snapshot: []`
- `:41` **`REAL_BUILD_PROCESS_EXITCODE=0`**，随后 `:35` `go build ok (cgo linked against sherpa-onnx C API)`、`:60` `wisp doctor: PASS`
⇒ **AC#10 字面那一发补上了：真 npm、必须绿、名册确实换名**，且换出来的两枚名字与票面 AC#9 里 `git archive` 那次的现量名（`index-vdBrT8rM.js`／`index-yy8KMgdf.css`）逐字相同 ⇒ 与本机盘上那四枚陈旧名不是同一趟（`r2` 具名登记的「合成 fresh 产物」残留由这一发**消掉**）。
镜像树的 `build.ps1` 与仓内那份的关系：`IDENTICAL_MODULO_CRLF`（`logs/10` 末段 `cmp`；clone 检出走 `attr/text eol=crlf`），行数 305 同。

**同一份源码连跑两次的「真 vite 第二发」＝已实测，代价是真的**（把 `r2` §1.1 从推理升成读数）：`real/` 树里第二次跑**同一份未改的 `build.ps1`、真 `npm run build`** ⇒ `logs/11-real-second-run.txt`
`:21` 前置快照逐字 `pre-build dist snapshot = 4 file(s) [.gitkeep assets\index-vdBrT8rM.js assets\index-yy8KMgdf.css index.html]`（＝第一发产的两枚名），
`:27` `vite v8.3.0 building …` 真跑、`:32-34` 产出的还是那两枚名，
`:42` **`FATAL: … Named cause: unprovenanced page bytes … vite rewrote byte-identical content-hashed names and no disk evidence can separate that from npm writing nothing …`**，
`:43` **`SECOND_REAL_BUILD_PROCESS_EXITCODE=1`**。
⇒ ①(c) 那一形**不需要桩也成立**：真 vite 对未变源码给同名 ⇒ 红；且这一发同时是 ③(a) 那枚「同一 job 第二发」路径的**后果读数**（真命中就真红）。

### ③ 那两枚「只有读码、没有实测」的 CI 边

**(3a) 同一枚 job 内会不会跑到第二遍＝字面成立，但被本腿收窄一句**
尺＝`git show HEAD:.github/workflows/ci.yml`（941 行，落 `/d/tmp/wisp274v2/ci.yml.head`）＋逐 job 完整 step 列表（`logs-3a-steps.txt`）。
- 三处调用逐字都是 `run: powershell -NoProfile -ExecutionPolicy Bypass -File scripts/build.ps1 -Env dev`，行号 `ci.yml:571`（`test-windows`，job 起 `:505`）／`:757`（`slo-smoke`，起 `:739`）／`:820`（`slo-full`，起 `:797`）。
- **step 列表逐枚读完：每枚 job 里 `build.ps1` 只出现一次**；三枚 job 都**没有** `strategy`／`matrix`（`grep -nE 'strategy|matrix|fail-fast' ci.yml` 在 HEAD 里 **0 命中**）⇒ 同一 run 内不存在「矩阵把同一 job 在同一工作目录里跑两遍」这枚形状。
- ★**但同一枚 job 里确实存在第二发的门**，不在 yml 的 step 里，在 `scripts/slo-check.ps1:294-314`：`if (-not (Test-Path $WispExe)) { … Invoke-ExternalProgram powershell … '-File', (Join-Path $RepoRoot 'scripts\build.ps1'), '-Env','dev' … if ($build.code -ne 0) { Fail ("build.ps1 failed with exit code {0}") } }`。
  `slo-smoke`（`:757` build → 紧接 `:762`「SLO smoke gate」）与 `slo-full`（`:820` → `:825`「SLO full gate」）都是**同一 job、同一工作目录**里前后两枚 step ⇒ 一旦 `build/wisp.exe` 在那两枚 step 之间不存在（step 被改顺序／`build/` 被上一枚 step 清掉／将来有人给 gate 加 `clean`），`slo-check.ps1` 就会**在同一趟里第二次调用 `build.ps1`**，第二次的前置快照里已躺着第一次的两枚哈希名 ⇒ 命中 `:204` ⇒ 而那一步报出来的是 `slo-check.ps1:314` 那句 `build.ps1 failed with exit code 1`，**闸门具名句被包在外层，读数会难认一倍**。
  ⇒ 对编排者那句「代价只落在同机重复构建，**不落在 CI 的第一发**」的裁：**「第一发」三字成立**（happy path 每 job 只一发），但**「不落在 CI」要收窄**：CI 里存在同 job 第二发的代码路径（`slo-check.ps1:294`），只是今天它的前置条件（exe 缺失）不成立。本腿⛔没在 CI 上真跑（零推送），这一格是**读数（yml＋脚本逐字）＋形状复算**，不是 CI 实测。
  旁证一枚（不是第二发，但同名相关）：`ci.yml:887`／`:902-903` 的 `lint-frontend` 是**另一枚 job** 自己跑 `npm ci` ＋ `npm run build`，它不调 `build.ps1` ⇒ 与本闸无关，登记以免被误当第二发。

**(3b) self-hosted `slo-full` 的工作目录＝取不到（具名，附我拿到的现量）**
- 简报给的三把尺全空：`ls -d /d/a`、`/c/actions-runner*`、`find /d|maxdepth 3 -name '_work'`、`find /c …` 全 **0 命中**（`logs/13-runner-probe.txt`）。本腿换一把尺拿到：`tasklist //FI "IMAGENAME eq Runner.Listener.exe"` ⇒ **`Runner.Listener.exe` PID 17844 在跑**，`(Get-Process -Id 17844).Path` ⇒ **`E:\work\base\actions-runner\bin\Runner.Listener.exe`**（`logs/14-runner-hunt.txt`）。⚠ 简报那句「`find /d -maxdepth 3`」够不到，**runner 在 E 盘**。
- 工作目录（⛔ 只读，没写一个字节、没跑 runner）＝`E:\work/base/actions-runner/_work/wisp/wisp`，带 `.git`；`frontend/dist` 现在**只有 `.gitkeep` 一枚**（md5 `d41d8cd9…`，mtime `2026-10-04 08:56:28 +0800`），而 `build/wisp.exe` mtime `2026-10-07 08:40:00`、`.git/FETCH_HEAD` `08:39:12`（`logs/16-runner-workspace.txt`、`logs/18-runner-mtimes.txt`）。
- 最新一枚 job 的时间线（`_diag/Worker_20261007-003708-utc.log`，只 grep 不改）：Checkout 步骤 `00:37:15Z→00:39:13Z exit 0`（1:57.96）、`Build wisp.exe (deps cached on the runner)` `00:39:17Z→00:40:02Z exit 0`（45.7s）、SLO 全量门 `00:40:02Z→00:42:14Z exit 0`、`"result": "succeeded"`。**这一发跑在 08:39 本地，早于 `91c90aa1` 落仓（16:0x 本地）⇒ 新闸门今天在这台 self-hosted 上一次都没真跑过。**
- ⇒ **判不动的那半**：「`actions/checkout@v4` 默认 clean 在这一形下实际会不会把 dist 清回只剩 `.gitkeep`」我**取不到机制级凭据**——盘上状态确实是闸门需要的 clean-checkout 形（前置快照＝`1 file(s) [.gitkeep]`），但我没法把它归给 checkout 的自动 clean：那枚 job 的 build 步骤 rc=0（它必然在 dist 落过页面字节），此后**没有任何更新的 job**（`Worker_20261007-003708` 是最新一枚），而那些字节现在不在；同时**同为 ignored 的 `build/` 里那枚 08:40 的 exe 活了下来** ⇒ 一句 `git clean -ffdx` 解释不了这个不对称。`_diag` 里可读的行也只有 `Cleaning runner temp folder`（`logs/17-runner-diag.txt`），拿不到 checkout 的 `git clean` 命令行。
  ⇒ 交回编排者的一句话：**「self-hosted 第一发拿到 anchor-only 快照」这一形此刻在盘上为真，但机制未证**，翻勾凭据⛔不含这一格；要钉死它需要 runner 侧一枚读数（推送后或本机手跑一次 `slo-full` 的 checkout 步）。

### ④ 既有八支是否仍各自红＝**成立（8/8 各自报名字，没有互相打瞎）**

`logs/22-battery2.txt` 与 `logs/20-stub-battery.txt`（行号＝FATAL 行）：

| 支 | 台件形状 | rc | 具名红因逐字（关键段） |
|---|---|---|---|
| `:123` npm 取不到（node 用 `node.cmd` 桩满足，PATH 里**摘掉** `D:\work\server\node14`） | `logs/22:165` | 1 | `FATAL: frontend step: npm could not be resolved on PATH (tried npm.cmd, then npm). … so it fails instead of skipping.` |
| `:136` `npm ci` 非零（**r2 没跑这支；node_modules ABSENT**，桩对 `ci` 退 7） | `logs/22:155` | 1 | `FATAL: frontend step: 'npm ci' failed with exit code 7 … the page bundle cannot be built without its dependencies.` |
| `:160` `npm run build` 非零（桩退 3） | `logs/22:133` | 1 | `FATAL: frontend step: 'npm run build' failed with exit code 3 in D:\tmp\wisp274v2\stub\frontend` |
| `:172` 无 dist | `logs/20:147` | 1 | `FATAL: … D:\tmp\wisp274v2\stub\frontend\dist does not exist. Named cause: dist directory missing, so wisp.exe would embed no page.` |
| `:180` 只剩锚文件 | `logs/20`（IV_180） | 1 | `FATAL: … holds only the anchor file(s) [1 file(s), none beside frontend/dist/.gitkeep] … Named cause: build produced no page artifacts.` |
| `:183` 无 `index.html`（产了 `assets/index-ZZZZ0000.js`） | `logs/22:83` | 1 | `FATAL: … has 1 artifact file(s) but no index.html … Named cause: missing entry file.` |
| `:187` 入口 0 字节 | `logs/22:108` | 1 | `FATAL: …\dist\index.html is 0 bytes - a present-but-empty entry file is not a page. Named cause: empty entry file.` |
| `:204` 新那支 | `logs/22:29/:57`、`logs/20:32` | 1 | `Named cause: unprovenanced page bytes …`（三形见 ①） |

顺带两支非 AC#10 射程但同一 `Fail` 出口的：`:101` 无 `package.json`／`:119` node 取不到（后者＝`274-v1` 的 M-iv，本腿⛔没重跑，只登记）。
**`:135` 与 `:159` 两支的 `Fail` 走的是同一枚 `:47-50`（`Write-Host "build.ps1: FATAL: $Message"` ＋ `exit 1`）**——盘上 `Fail` 定义只有一处出口，没有第二条退出路，也没有「打印后继续」。

⚠ **自报一枚我自己的台件缺陷（不作废任何判语，但必须写在这里）**：`logs/20-stub-battery.txt` 里**凡是要 npm 写文件的支**（`I_b`／`I_c`／`IV_183`／`IV_187`／`IV_160`／`IV_136`／`IV_123`）第一次跑的是**废发**——我用 `printf` 生成 `.cmd`，bash 把 `%1 `／`%2 ` 读成宽度转换符，桩文件被写坏、cmd 退 **255**，于是那几发的红是「npm 非零 ⇒ `:160`／`:136`」而不是我要的形状。凭据：`logs/20:60/88/110/129` 的 `exit code 255` 与跑后 md5 逐枚未变。**①(b)(c) 与 ④ 的这几行全部由 `logs/22-battery2.txt`（heredoc＋`awk '{printf "%s\r\n",$0}'` 重写桩）重跑覆盖**，`logs/22` 里桩正文逐字贴在每发前面。`I_a`／`IV_172`／`IV_180` 三发在 `logs/20` 里就是有效的（桩退 0 未受影响）。

### ⑤ 恒真检查＝**成立（而且比预期更硬）**

尺：把**仓外** `green/` 树里的一份**副本** `scripts/build_mut.ps1` 的 `:203` 判据中和（`if ($freshNames.Count -le 0) {` → `if ($false) {  # 274-v2 NEUTRALISATION of the provenance gate`），⛔ 仓内 `scripts/build.ps1` 一字节没动（终态 md5 回查 `0b5fbbb5…`，§4）。
`diff` 原文（`logs/21-green-cases.txt:55-59`）逐字只有两行：`< if ($freshNames.Count -le 0) {` / `> if ($false) { …`。
跑 ①(a) 那一形（陈旧真页面＋桩 npm 一字不写）⇒
- `logs/21:74` 逐字：`build.ps1: frontend ok: 4 file(s) in frontend/dist (entry index.html is 1044 bytes): .gitkeep=0 assets\index-BRKj5OIJ.css=49943 assets\index-BVKlegVD.js=553469 index.html=1044 || provenance: this run's npm added new artifact name(s) [] over the pre-build snapshot …`（**新增名＝空集合，但通路照走**）
- `logs/21:81` **`CASE_L_tautology_neutralise_203_PROCESS_EXITCODE=0`**（一路走完 `go build ok`／SHA256SUMS／`wisp doctor: PASS`）
⇒ **中和 `:203` 之后 (a) 由红转绿**，证明「名字集合换了没换」这一条**确实由 `:203-204` 在管**，不是别的东西顺手拦住的；同时这一发把 `274-v1` 的 M-iii-2（陈旧件冒充、整条绿）在本腿台件里**独立复现了一遍**，牙与非牙的差别只在那两行。

---

## §2 对 AC#10 本身的裁（给编排者翻勾用，本腿一枚不翻）

**成立**。两半字面判据都由本腿自跑兑现：桩 npm 一字不写⇒红（①a，含 `Named cause: unprovenanced page bytes`）；真跑 `npm run build`⇒绿且名册换名（②，vite 真产物 `index-vdBrT8rM.js`／`index-yy8KMgdf.css`，整条 `rc=0`）。
`r2` 具名登记的**唯一实质残留**（「②用的是合成 fresh 产物」）由 ② 那一发**消掉**。
形选择的合法性：AC#10 原文给的是「名册不同名 **或** mtime」二选一（`274:56` 逐字），取名册形**是字面许可的一支**，⛔ 不是绕开；`r2` §1 那五条弃 mtime 的理由里「生成物 `frontend/src/styles/tokens.generated.css` 会自己往前走」这一条本腿⛔没独立核（读码继承，不当凭据）。
残余（不否决本格，逐枚具名）：见 §5。

---

## §3 ⓐ 我推翻／收窄了编排者（或 `r2`）的哪几句（⛔ 以原文与盘上为准）

1. **推翻（环境事实）**：简报「本机 **node v24.18.0** 在 PATH」⇒ 盘上 `node = D:\work\server\node14\node.exe (**v24.9.0**)`、npm `11.6.0`（`logs/10-real-npm-build.txt:11`）。不改判据，但今后谁写读数请写 v24.9.0。
2. **推翻（尺不够长）**：简报 ③(b) 给的三把尺（`/d/a`、`/c/actions-runner*`、`find /d -maxdepth 3 -name '_work'`）**全空**，runner 实际在 **`E:\work\base\actions-runner`**，工作目录 `E:\work\base\actions-runner/_work/wisp/wisp`（`logs/14`、`logs/15`）。⇒ 那把尺要改成按 `Runner.Listener.exe` 的进程路径反查。
3. **收窄（`r2` §1.1 第三段 ＋ 简报 ③a 的隐含口径）**：「CI 侧不受影响／代价只落在同机重复构建，不落在 **CI**」——**同一枚 job 内第二发是有代码路径的**：`scripts/slo-check.ps1:294-314` 在 `build/wisp.exe` 缺失时会再跑一遍 `build.ps1`，而它就在 `slo-smoke`／`slo-full` 的 build step 之后一枚 step。happy path 不命中（exe 刚被上一枚 step 建出来），所以「第一发不受影响」成立；但**「不存在第二发」不成立**，且真命中时报错文本会被 `slo-check.ps1:314` 那句包住一层。凭据＝§1③a。
4. **收窄（护栏的射程）**：简报「绝不真删／真盖仓内 `frontend/dist` **那 4 枚文件**」——盘上 `frontend/.gitignore:12` 写的是 `dist/*`，跟踪的只有 `.gitkeep` 一枚（`git ls-files frontend/dist` ＝ 1 枚，`logs/setup-real.txt:114-116`）⇒ 那 3 枚是**未跟踪产物**，`git archive`／clone 根本不带。这句话对「别动本机字节」仍然有效，但对「CI 里那一形长什么样」的推理必须按「跟踪面只有锚」来算（`r2` §1.1 已经这么算了，我核过是对的）。
5. **登记（不推翻）**：简报「`gofmt -l`/`gofumpt -l` 会列出 `cmd\wisp\models.go` 等——带 CR 的共 13 枚」⇒ 本腿实测**逐字对上**：`tr -cd '\r' | wc -c` 逐枚数＝**13 枚，产码 5 枚**（`cmd/wisp/models.go`、`internal/agent/approval/pending_read.go`、`internal/agent/tools.go`、`internal/risk/provenance.go`、`internal/tools/bridge.go`）＋ `.scratch/wisp/probes/257/r2/mut/**` 8 枚（`logs/23-gates.txt`）。⛔ 一枚没格式化。另：`gofumpt -l $(cat 全部跟踪 .go)` 会 `Argument list too long`，换 `xargs` 后 `cmd/`+`internal/` 命中 **0 枚**（`logs/gofumpt.txt`）；`gofmt -l` 28 枚里 **23 枚是 `.scratch/**` 别人写的突变件**、5 枚是上面那 5 枚产码。
6. **口径照编排者改的办**：冷跑整包红名册按**名字集合**作差。本腿 `GOCACHE=D:/tmp/gocache-274v2` 冷跑 ⇒ `GOTEST_COLD_RC=1`，具名红 **5 枚**：`TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestC21DesignTokensFourWayAgree`／`TestC21TableColourRowsMatchTokensCSS`／`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`＝**已知 6 枚去掉 `TestResolvePerCallBudget`**（`grep -c` ＝ **0**，这一发它没出现）⇒ **作差 0 枚新增**，谁的账也不记。失败的三枚包（`cmd/wisp`／`internal/ball`／`internal/panel`）与这 5 枚一一对得上。红因抽样逐字（`logs/gotest-cold.txt:49`）：`tokens_table_test.go:1468: read design/assets/tokens.css: open D:\work\workspace\projects plans\Wisp\design\assets\tokens.css: The system cannot find the path specified.`——`design/` 这枚路径今天在本机工作树里读不到，属别人在飞的数，⛔ 不是我造成的（本腿没写任何产码）、⛔ 我不去动它。
7. **一枚要编排者认领的环境事实（⛔ 不是本腿的账，也⛔不是本腿写的）**：上条那句红因不是偶发——`git status --porcelain -- design` 现在报 **31 行 ` D design/…`**（跟踪件在工作树里缺失，`ls design` 只剩 `doubao`／`old`）。起手锚的六族尺面**没含 `design`**（简报给的就是那六族），所以本腿**无法自证这 31 枚删除早于我**；能自证的是本腿全部写面只有 `.scratch/wisp/probes/274/v2/**`（每笔文件枚数见 ⓒ），对 `design/` 没有任何写/删动作。⇒ 建议台账具名挂一枚：工作树 `design/` 缺失＝那 4 枚 C21/主题红的现因。

---

## §4 ⓑ 终态自证（`logs/24-final-proof.txt`，`date` ＝ `Wed Oct  7 17:02:50 CST 2026`）

- `git status --porcelain -- cmd internal scripts .github docs frontend build` ＝ **0 行**（＝起手）。
- 仓内 dist 四枚 md5 终态＝起手，`diff` 空（`DIST_MD5_DIFF_EMPTY`）；`build/wisp.exe` 仍 `e6c8e52bb14f15a7983bbc6b093e6058`（本腿没盖它）。
- **`git show 91c90aa1:scripts/build.ps1 | md5sum` ＝ `0b5fbbb5f43850f4e3a3545c4f4eaae2` ＝ 盘上 `md5sum < scripts/build.ps1`**（逐字节对回；`i/lf w/lf attr/text eol=crlf`，`git show --numstat 91c90aa1` ＝ `41 1 scripts/build.ps1`，仍只动这一枚文件）。
- 仓内 `frontend/node_modules` 三枚样件 `md5sum -c` 全 **OK**（⛔ 硬链接没用过）；`third_party/sherpa-onnx/*.dll` 链接数 15（我 `cp -al` 出去的镜像还挂着，只读，未发生任何写入）。
- 门禁必要那几把：`sh scripts/d22scan.sh` rc=**0**（`d22scan: clean - no D22 ban violations`）、`go vet ./cmd/wisp/ ./internal/panel/` rc=**0**、gofmt/gofumpt 见 §3-5、整包 `go test -count=1`（`GOCACHE=D:/tmp/gocache-274v2` 冷跑）**只跑了这一发**。

---

## §5 没做完／射程外（逐枚具名，⛔ 不写成通过）

1. ~~同一份源码连跑两次的真 vite 第二发~~ **已补量**：`logs/11-real-second-run.txt:42-43`，真 vite 第二发 ⇒ **红在 `:204`，rc=1**。见 §1② 末段。本腿交表时**没有未量的 AC#10 字面判据**。
2. **③(b) 的机制**（checkout@v4 默认 clean 到底清不清 ignored 的 `frontend/dist`）＝**取不到**：runner 与工作目录找到了、盘上状态与 job 时间线拿到了，但归因不成立（不对称：同为 ignored 的 `build/` 活下来了），`_diag` 里没有可读的 `git clean` 命令行。见 §1③b 末段。
3. **名册串的接缝（读码观察，⛔ 未实测）**：`:152` 用 `-join ' '` 把名册拍成一枚字符串，`:200` 再 `-split ' '` 拆回来 ⇒ **一旦 `frontend/dist` 里出现文件名带空格**（真实 vite 不产，但用户可往 dist 里放任何东西），`a b.js` 会被拆成两枚"名字"，可能**凭空造出一枚「快照里没有的名字」⇒ 假绿**，也可能反过来把两枚真名字并成一枚。这条属"名册形"的实现接缝，不是 AC#10 要的判据本身；本腿没打这一形（预算），**登记为残留，不当已通过**。
4. `Fail` 之外的两支（`:101` 无 `package.json`、`:119` node 取不到）本腿没重跑：前者不在 AC#10 射程，后者已由 `274-v1` M-i 量过（⛔ 我未复跑，只登记行号 `scripts/build.ps1:119`）。
5. **CI 真跑**：⛔ 零推送，全部 CI 侧结论是 yml＋脚本逐字读数与形状复算，不是 run 侧观测。
6. **`:130` 那句「node_modules EXISTS ⇒ 不校验 lockfile」** 的整条后果（装了不等于对）不属本票射程，② 那一发就在这条分支上跑出来的（已具名）。

---

## 裁决摘要（一句话版）

**AC#10 这道闸门有牙**：三形反形全部红在 `:204` 那枚具名因上（含 r2 没测过的「npm 只写 node_modules」与「同名重写不同字节」），红因给的处置实测能走通到 `rc=0`；AC#10 字面要求的**真 `npm run build`** 一发由本腿补上（真 vite、真换名、整条 `rc=0`）；把 `:203` 中和之后同一发立刻转绿（rc=0）＝牙确实在 `:203-204`；既有八支各自报各自的名字，没被打瞎。**两处未决 ＋一枚接缝残留**：self-hosted 那枚 clean 的**机制**取不到（盘上状态与 job 时间线量到了，归因不成立）；名册串的 `-join ' '`/`-split ' '` 对**含空格文件名**可能造出假绿（读码观察，未实测，真实 vite 不产这种名）。两枚都⛔不否决本格，但翻勾凭据里不该写它们。另：AC#10 的**全部**字面判据（含真 npm 一发、真 vite 第二发、三形反形、恒真中和、八支各红）都已由本腿实测覆盖，⛔ 无未量项。
