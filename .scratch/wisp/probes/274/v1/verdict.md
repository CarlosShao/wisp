# 274-v1 — 非实现者对抗验收件（`scripts/build.ps1` 第 2 步真构建，形 乙-1，写面 `35633445`）

腿：`274-v1`（验收者 ≠ 实现者，D22 双角色 / `SPEC-12 §4.3` #1）。
被审实现者：`274-r1`，写面一笔 `35633445`。
⛔ 本格只裁不翻：AC 框一枚不翻、票面一字不改、不 push、不 `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`。
⛔ 本文所有读数都是**我自己现跑的**；引实现者的日志只作为"待复跑的断言"出现，不作为凭据。

---

## §0 起手锚（本腿第一条 commit 之前现量，逐字）

```
$ date
Wed Oct  7 12:28:44 CST 2026

$ git rev-parse --short HEAD
663a176d

$ git status --porcelain -- cmd internal scripts tools .github docs frontend design build
 D design/assets/base.css
 D design/assets/icons.js
 D design/assets/theme.js
 D design/assets/tokens.css
 M design/doubao/README.md
 M design/doubao/demo/app.js
 M design/doubao/demo/index.html
 M design/doubao/demo/styles.css
 D design/index.html
 D design/screens/approval.html
 D design/screens/ball.html
 D design/screens/chat.html
 D design/screens/config.html
 D design/screens/cost.html
 D design/screens/firstrun.html
 D design/screens/palette.html
 D design/screens/privacy.html
 D design/screens/security.html
 D design/screens/states.html
 D design/screens/tasks.html
?? design/doubao/01-ball-states.jpg
?? design/doubao/demo/lib/
?? design/doubao/demo/rb-files.js
?? design/doubao/demo/rb-plugins.js
?? design/doubao/demo/rb-review.js
?? design/doubao/demo/rb-terminal.js
?? design/doubao/demo/rightbar.js
?? design/doubao/demo/screens/home.js
?? design/doubao/demo/screenshots/
?? design/doubao/demo/sidebar.js
?? design/old/
（六族 `cmd internal scripts tools .github docs frontend build` 之内 **0 行**；上面每一行都属 `design/**`）

$ git log --oneline -3 -- scripts/build.ps1
35633445 274-r1: build.ps1 step 2 now builds the page and fails if it cannot (form B-1)
cc6eaa65 244-r1: build.ps1 追加 -H=windowsgui，双击不再带黑控制台
bfcb230e feat(models): C29 signed manifest (6 entries, P3 verdicts, matcha blocked-p3) + dev minisign keypair ceremony + signer tool/script + compose model-mirror (18081-18083 fixtures) + PRECHECK P3/P5

$ wc -c scripts/build.ps1 ; md5sum scripts/build.ps1 ; wc -l scripts/build.ps1
13638 scripts/build.ps1
644c2f6af1fa3fc7f8b3ab38c20753f5 *scripts/build.ps1
265 scripts/build.ps1

$ git show 35633445^:scripts/build.ps1 | wc -l
172
$ git show 35633445:scripts/build.ps1 | wc -l
265
```

锚文件：`.scratch/wisp/probes/274/v1/logs/00-anchor.txt`（上面那批的原文落盘）。

### §0.1 起手即登记的两条环境事实（⛔ 不是本腿的差，不计入任何作差）
1. **`design/**` 起手就脏**：现量＝**31 行**＝**16 枚 ` D` ＋ 4 枚 ` M` ＋ 11 枚 `??`**（尺见 `logs/00-anchor.txt` 末段：`git status --porcelain -- design | cut -c1-2 | sort | uniq -c` 逐字回 `16  D`／`4  M`／`11 ??`）。
   ⚠ **这一条与本腿简报相反**：简报写"4 枚 ` D` ＋若干 ` M`"，票面 `274-…md:145` 与台账 `A664` §5 同样写"4 枚 ` D`＋若干 ` M`"——盘上真值是 **16 枚 ` D`**。按简报纪律「盘上赢」，见 §4。
   按简报的纪律："盘上赢"——见 §4 编排者句 overturn 段。**处置不变**：不碰、不清理、不算进本票作差。
   落账差异**只可能来自**我对 `design/**` 的计数口径，而 `git status --porcelain -- <六族>` 内是 0 行，所以本票任何作差读数与 `design/**` 无关。
2. **HEAD `663a176d` 晚于写面 `35633445`**（`git merge-base --is-ancestor 35633445 HEAD` ⇒ 真），且 `scripts/build.ps1` 在 `35633445` 之后无人再动（`git log -3 -- scripts/build.ps1` 首行就是它）。⇒ 盘上那枚 265 行**就是**被审件；我审的是 `35633445:scripts/build.ps1` == 盘上文件（§6 用 md5 钉）。

### §0.2 待审清单（本腿要自己跑的尺，逐条）
- **AC#1** 未修码两形（`35633445^` 那份 172 行 `build.ps1`）：前端步静默跳过 ⇒ exe 零页面字节 ⇒ 整条通路 rc=0；＋ embed 级尺 `go list -f '{{.EmbedFiles}}'` ＋ exe 侧 `panel-assets` ＋ **正控**（同一把尺喂真页面必须回非空名册）。
- **AC#4** 反形敏感性：复跑 `274-r1` 那两发定向突变（① node/npm 取不到 ② npm 谎报成功而 dist 只剩锚），**再加本腿自造的突变 ≥1 发**。
- **AC#6** 门禁四数（独占窗口、串行）：`sh scripts/d22scan.sh` rc=0 ＋ 正控先绿；`gofmt -l`／gofumpt v0.12.0 `-l` 只喂 `.go`；`go vet ./cmd/wisp/ ./internal/panel/` rc=0；`go test ./cmd/wisp ./internal/... -count=1` RUN/PASS/FAIL/SKIP 四数＋首尾逐名红名册作差。
- **AC#7** 方法裁定：`frontend/dist` 现态 vs 实现者起手基线（逐枚 md5），裁"mv 还原"算不算 `-overlay` 字面。

---

## §1 逐格裁决

### AC#1（未修码读数：今天拦不住）＝**成立**（全部本腿自跑，⛔ 未引 `r1/logs`）
台件：仓外硬链接镜像树 `D:/tmp/wisp274v1/tree`（`cp -al` 逐条顶层目录，⛔ 不含 `.git`、⛔ 不含 `frontend/node_modules`），
`frontend/dist` 里我只 **unlink**（`rm` 硬链接＝减链接数，不触原文件）掉 3 枚陈旧件 ⇒ 该树 dist＝`.gitkeep` 一枚＝干净检出形。
仓内基线在 unlink 前后逐枚 md5 相同（`logs/01:9-19`：BEFORE/AFTER 四枚全同，`REPO DIST UNCHANGED rc=0`），`build/` 在 scratch 侧整目录 `rm -rf`（`logs/02`）⇒ 本腿从未写过仓内 `build/wisp.exe`（终态 md5 仍 `e6c8e52b…`，`logs/04` 末段）。

1. **embed 级尺，锚文件形**（`logs/01`）：
   `$ cd D:/tmp/wisp274v1/tree && go list -f '{{.EmbedFiles}}' ./frontend` ⇒ 逐字 `[dist/.gitkeep]`、`rc=0`。
2. **embed 级尺，正控**（同一把尺、本仓工作树、⛔ 未改任何字节）⇒ 逐字
   `[dist/.gitkeep dist/assets/index-BRKj5OIJ.css dist/assets/index-BVKlegVD.js dist/index.html]`、`rc=0`。尺不瞎。
3. **整条通路（未修码那份 `35633445^:scripts/build.ps1`，172 行，stage 进 scratch 后原样跑）**（`logs/03`）：
   - `logs/03:6` 逐字：`build.ps1: frontend step skipped (no frontend yet; embed lands in S5 per SPEC-11 §2.2)`
   - `logs/03:12` 逐字：`OLD_BUILD_PS1_PROCESS_EXITCODE=0`（走完 `go build ok`／DLL 并置／SHA256SUMS／`wisp doctor: PASS`）
   - `logs/03:20-21` 逐字：`wisp panel-assets: assets NOT BUILT` ⇒ `PANEL_ASSETS_EXITCODE=1`
   - `logs/03:24-25` 逐字：`wisp panel-assets: panel: embedded assets are not built (run npm run build in frontend/)` ⇒ `MANIFEST_EXITCODE=1`
   ⇒ **同一棵源码、页面字节为零、出货通路整条 rc=0**＝票面 §1 那句现象由本腿独立复现，不是转述。
4. **行尾符口径（A664 §2 要求同时报）**：我的 scratch 树是**工作树硬链接**，不是 `git archive` ⇒
   `logs/01`：`frontend/index.html` CR 数 scratch＝0／repo＝0、`cmp` 回 `IDENTICAL`；`frontend/src/App.tsx` 同样 0／IDENTICAL。
   ⛔ 因此 A664 那枚 archive-CRLF 陷阱在本腿台件里**不存在**（尺＝`tr -cd '\r' | wc -c`，⛔ 不是 `grep -c $'\r'`）。

### AC#4（反形敏感性／恒真句攻击）＝**字面成立，带一枚具名残留缺口**（详见 §2）
被审件行号一律 `35633445:scripts/build.ps1:<line>`（＝盘上 `scripts/build.ps1:<line>`，两树同 265 行、md5 见 §6）：
- `:38-41` `function Fail([string]$Message) { Write-Host "build.ps1: FATAL: $Message"` / `exit 1` ⇒ 每个失败支都**自己报名字＋非零**，不是"打印后继续"。
- `:108-111` node 取不到 → `Fail`；`:112-115` npm 取不到 → `Fail`；`:126-127`／`:136-137` `npm ci`／`npm run build` 非零 → `Fail`；
  `:148-149` dist 目录不存在 → `Fail`；`:155-158` 只剩锚文件 → `Fail`；`:159-160` 无 `index.html` → `Fail`；`:162-165` 入口 0 字节 → `Fail`。
- CI 侧真的会红：`35633445:.github/workflows/ci.yml:571`／`:757`／`:820` 逐字
  `run: powershell -NoProfile -ExecutionPolicy Bypass -File scripts/build.ps1 -Env dev`，三处**都没有** `continue-on-error`
  （尺见 §2 末），所以 `exit 1` 就是那一步红。
三发突变的原始退出码与具名红因逐字都在 `logs/04`，本腿裁论见 §2。

---

## §2 恒真句攻击（AC#4 的核心射程；五发全部本腿自跑，原始退出码逐字在 `logs/04`／`logs/05`／`logs/06`）

所有突变都在仓外镜像树 `D:/tmp/wisp274v1/tree` 里跑**当前**那份 265 行 `build.ps1`（staged＝`git show 35633445:scripts/build.ps1`，`logs/04:4`），
⛔ 仓内 `frontend/dist`／`build/wisp.exe` 一个字节没被写过（每发末尾都有 md5 回查，`logs/04` 末段、`logs/06` 末段：
`f98bfc4b…`／`db4db7a2…`／`70128a3d…`／`d41d8cd9…`／`e6c8e52b…` 与起手逐枚相同）。

| 发 | 形状 | 我看到的具名红因（逐字）／绿 | 码 |
|---|---|---|---|
| **M-i**（实现者的 (i)） | PATH 剥掉 node 与 npm | `build.ps1: FATAL: frontend step: node could not be resolved on PATH (tried node.exe, then node). …` ＝ `35633445:scripts/build.ps1:110` | `MUTATION_PROCESS_EXITCODE=1` |
| **M-iv**（**本腿自造**，实现者没碰过那一支） | **node 找得到、npm 找不到** | `build.ps1: FATAL: frontend step: npm could not be resolved on PATH (tried npm.cmd, then npm). …` ＝ `:114` | `MUTATION_PROCESS_EXITCODE=1` |
| **M-ii**（实现者的 (ii)） | 桩 `npm` 报成功而 `frontend/dist` 只剩锚文件 | `build.ps1: FATAL: frontend step: npm run build reported success but frontend/dist holds only the anchor file(s) [1 file(s), none beside frontend/dist/.gitkeep]. That is the clean-checkout shape, so the exe would ship with zero page bytes. Named cause: build produced no page artifacts.` ＝ `:157` | `MUTATION_PROCESS_EXITCODE=1` |
| **M-iii-1**（**本腿自造**） | 桩 `npm` 报成功，而 dist 里**已躺着一份真页面（09-27 陈旧件）** | **这一步绿了**：`build.ps1: frontend ok: 4 file(s) in frontend/dist (entry index.html is 1044 bytes): .gitkeep=0 assets\index-BRKj5OIJ.css=49943 assets\index-BVKlegVD.js=553469 index.html=1044`（`logs/05:19`＝`logs/06:16`）——本发整条 run 随后死在 `:192`（见下条侧发现），码 1 | `M_III_PROCESS_EXITCODE=1` |
| **M-iii-2**（同上，把 `git.exe` 也从 PATH 摘掉让通路走完） | 同上 | 整条走完：`logs/06:16` 那句 `frontend ok` 原样，`logs/06:23` `M_III_RERUN_PROCESS_EXITCODE=0`，`logs/06:28-29` `panel assets embedded: 4 files, entry=index.html built=true` `PANEL_ASSETS_EXITCODE=0`，`-manifest` 逐枚＝`.gitkeep 0`／`assets/index-BRKj5OIJ.css 49943`／`assets/index-BVKlegVD.js 553469`／`index.html 1044` ＝ **npm 这一趟一个字节都没写，exe 仍报带页面** | rc=0 **绿** |

**裁论（诚实版，不替实现者圆场）**
- AC#4 那句字面判据（票面 `:50`：「把 dist 换成只剩锚文件那一形 ⇒ 出货判据必须红」）＋「不许写成两形都绿的恒真句」——**过了**：
  三发不同形状各自报出**不同的**具名原因并 `exit 1`，七条失败支一一对得上行号；`Fail` 只有一条出口 `exit 1`（`:38-41`），没有"打印后继续"的支。
  CI 那三处调用（`35633445:.github/workflows/ci.yml:571`／`:757`／`:820`）逐字都是 `powershell -NoProfile -ExecutionPolicy Bypass -File scripts/build.ps1 -Env dev`，
  且 `grep -c '^[[:space:]]*continue-on-error:' .github/workflows/ci.yml` ＝ **0**（16 处文字命中全是注释，`logs/07` §6）⇒ `exit 1` 真会把那一步染红。
- **但本腿自造的 M-iii 落绿了**，且绿得很具体：这道新闸门校的是**"dist 里有没有页面字节"**，不是**"这些字节是不是这一趟构建出来的"**。
  ⇒ 按派单纪律（"你的突变落绿 ⇒ 那一格不敏感，要说出并拒绝接受本格"）我**不写成无条件成立**：
  本格的**字面命题成立**，**射程外残留一枚具名缺口**＝「陈旧产物冒充本次产物」这一形今天的闸门看不见。
  这枚缺口归谁：票面 **AC#9**（`:54`「陈旧产物不许替『带页面』报绿」）正是为它立的框，而 `build.ps1` 那一步**没有**任何 per-file 出处校验可以兑现 AC#9 ⇒
  翻不翻 AC#4 由编排者拍；若要它自己就有牙，最小增量是现成的：`logs/06` 那种对拉（`frontend ok` 那行的名册必须与**本趟 npm 之前**的快照不同名，或直接校 `index.html` 的 mtime 晚于 `frontend/src` 最新 mtime）。
- **侧发现一枚（不属 274，但被 274 的通路放大）**：M-iii-1 死在
  `At D:\tmp\wisp274v1\tree\scripts\build.ps1:192 char:15  +     $short = (& git rev-parse --short HEAD 2>$null)` ⇒ `NativeCommandError`＋整条 rc=1。
  因＝树里没有 `.git` 而 PATH 里有 `git.exe`，`$ErrorActionPreference='Stop'`（`:31`）把原生命令的 stderr 升成终止错误。
  **公平归因**：那一行在 `35633445^:scripts/build.ps1:99` 就有（`logs/08` (a)/(d) 与 `grep -n rev-parse` 现量），⛔ 不是 274 引入的；
  但它说明"在没有 `.git` 的目录里跑 build.ps1"今天会红成一句与前端无关的错——CI 的检出带 `.git` ⇒ 不阻塞本票，登记即可。

---

## §3 实现者（`274-r1`，`.scratch/wisp/probes/274/r1/impl.md`）那几句：逐句复跑后的裁
1. **确认**（原句逐字，`impl.md:189`）：
   `build.ps1: FATAL: frontend step: node could not be resolved on PATH (tried node.exe, then node). The page bundle is a build input of wisp.exe, so this script fails instead of skipping it. …`
   我的复跑＝M-i，`logs/04:10` 打出**同一句**并 `MUTATION_PROCESS_EXITCODE=1`。⛔ 我不引它的 `logs/09`。
2. **确认**（原句逐字，`impl.md:207`）：
   `build.ps1: FATAL: frontend step: npm run build reported success but frontend/dist holds only the anchor file(s) [1 file(s), none beside frontend/dist/.gitkeep]. That is the clean-checkout shape, so the exe would ship with zero page bytes. Named cause: build produced no page artifacts.`
   我的复跑＝M-ii，`logs/04:33` 同一句＋`MUTATION_PROCESS_EXITCODE=1`；我这发的桩 npm 是 `D:/tmp/wisp274v1/fakebin/npm.cmd`（`logs/04:30` 里 node 被解析成 `fakebin\node.exe (v24.9.0)`），与它的 `fake-npm-shim` 无继承关系。
3. **确认**（原句逐字，`impl.md:133-134`）：`**出货通路整条绿（EXITCODE=0、走到 smoke test、写了 SHA256SUMS），而交出去的那枚 exe 一个页面字节都不带…**`
   我的复跑＝`logs/03:12` `OLD_BUILD_PS1_PROCESS_EXITCODE=0`、`logs/03:9` `build.ps1: wrote D:\tmp\wisp274v1\tree\build\SHA256SUMS`、`logs/03:20-21` `assets NOT BUILT` ＋ `PANEL_ASSETS_EXITCODE=1`。
4. **确认**（原句逐字，`impl.md:319`）：`已跑的替代读数：go vet ./cmd/wisp/ ./internal/panel/ rc=0（编译面干净）＋d22scan 正控先绿再 clean。`
   我的复跑＝`logs/07`：`GOVET_RC=0`；`runtests.sh: OK - packages=[./...] top-level: PASS=35 FAIL=0 SKIP=0, === RUN=77, '[no tests to run]'=0` 之后才是真扫，`D22SCAN_RC=0`。
5. **推翻半句**（原句逐字，`impl.md:269`）：`| gofmt -l／gofumpt v0.12.0 -l | **本腿零 .go 改动 ⇒ 无 .go 可喂，未喂**…| 不适用（如实记） |`
   "按写面"它没错（`git diff --name-only 35633445^ 35633445` 只回 `scripts/build.ps1`，`logs/07:256`），⛔ 但 AC#6 那一格点名的是 `cmd/wisp`／`internal/panel` 这两枚包，喂下去**不是空**：
   `gofmt -l cmd/wisp internal/panel` ＝ **`cmd\wisp\models.go`**（`logs/07` §3），`/d/work/base/gopath/bin/gofumpt.exe --version` ＝ **`v0.12.0 (go1.27.1)`**，`gofumpt -l cmd/wisp internal/panel` ＝ **同一枚 `cmd\wisp\models.go`**（`logs/09` 末段）。
   定因＝**盘上行尾符**：`cmd/wisp/models.go` `tr -cd '\r' | wc -c` ＝ **334**（＝逐行 CRLF），全仓跟踪 `.go` 里带 CR 的 **13 枚**（`/d/tmp/wisp274v1/crlf-go-files.txt`）；
   而 `git show 35633445^:cmd/wisp/models.go` 与 `35633445:` 那份**逐字节相同**（`cmp` 回 `IDENTICAL`，`logs/08`(a)）⇒ 与 274 无关、与票 275 那条"真未格式化名册"同族。
   ⇒ 本格读数要写成"**这台机器的 bench 树上这两枚包不是 gofmt/gofumpt-clean，红因＝CRLF，不是 274**"，⛔ 不许写成"不适用"。这正是 A664 那枚行尾符变量**第二次**咬到判据仪器（第一次是 `git archive`）。

## §4 编排者的句子（票面／A664）——本腿推翻或收窄的，逐句给原句与我的尺
1. **推翻**：票面 `:145` 与台账 `A664` §5 都写 `git status --porcelain -- design` 现量＝"**4 枚 ` D` ＋若干 ` M`**"（派单简报同句照抄）。
   我的现量（`logs/00-anchor.txt` 末段，尺＝`git status --porcelain -- design | cut -c1-2 | sort | uniq -c`）＝
   **`16  D` / `4  M` / `11 ??`，共 31 行**。 ⇒ "4 枚 D"少计了 `design/index.html` 与 `design/screens/*.html` 那 11 枚删除。
   **后果**：处置不变（不碰、不算进作差），但"起手即脏的规模"被写成实际的四分之一；按简报纪律"盘上赢"，本腿所有作差一律排除 `design/**`，且六族 `cmd internal scripts tools .github docs frontend build` 现量 **0 行**（`logs/04`、`logs/06` 末段两次回查）。
2. **推翻半句**：`A664` §2 那句更正（"`git archive` 不是干净检出，它按本机 autocrlf 改字节"）**成立且我扩一条**——
   同一枚行尾符变量不只咬 `git archive`，还咬 **gofmt/gofumpt 门禁本身**（见 §3.5）。⇒ 今后任何腿报 `gofmt -l` 之前要先报 `tr -cd '\r' | wc -c`，否则"未格式化"与"CRLF"两种红分不开（这条与台账第 120 条"一个门有多种红法不许压成一格"同形）。
3. **收窄**：票面 `:108-109`／`A664` §1 的 AC#9 成立口径"exe 内嵌物 == **同一次运行**里 dist 的产物"——我复跑 exe 侧（只读，`logs/07` §7）：
   `build/wisp.exe` md5 `e6c8e52b…` → `panel assets embedded: 4 files, entry=index.html built=true`、`-manifest` ＝ `.gitkeep 0 e3b0c44298fc1c14`／`assets/index-B8yINMF1.js 551989 8f06145dac449fb6`／`assets/index-yy8KMgdf.css 49540 d0b198664b250973`／`index.html 1044 9b7856b949d63989`（与 `A664` §1 那串**逐枚相同**＝它对）。
   ⛔ 但**盘上今天的 `frontend/dist` 装的是另外两枚名字**（`index-BVKlegVD.js 553469`／`index-BRKj5OIJ.css 49943`，`logs/10`）⇒
   "同一次运行"那对**已经从盘上消失了**（dist 被还原成 09-27 基线、exe 留在 11:57 那份）。所以 AC#9 现在的凭据**只剩日志**，盘上无法重建那一对；
   凡后续再拿 `build/wisp.exe` 报"带页面"，读数主体是**它自己内嵌的那份 fresh 产物**，⛔ 不是工作树的 dist。这一格我按"已成立但不可复核"记。
4. **确认**：`A664` §4①"exe 被本票重跑覆盖＝保留新的、不还原，起手那份在 `/d/tmp/wisp274r1/wisp.exe.pre-274r1-backup`"——我现量该文件在盘：
   md5 `7032d94d36b2e060ec4a24a6a04c1f47`／31,076,405 字节，与 `impl.md` 起始那份记录一致（`logs/10`）；仓内现 exe `e6c8e52b…`＝另一枚。⇒ 那句"撤销口令＝「还原 274 前 exe」"是**可执行的**。
5. **简报里"实现者用 mv 而非 -overlay"这一条：盘上核实为真**——`/d/tmp/wisp274r1/` 里有 `dist-stale-backup`／`dist-stale-staging`／`node-modules-pre-274r1`／`fake-npm-shim` 与两份 `mut-i-no-node.ps1`／`mut-ii-stub-npm.ps1`（`logs/10` 的 `ls`），没有 `-overlay` 用到的 json 层；`impl.md:282-284` 自述三次 `mv` 动 dist＋一次动 `node_modules`。裁定见 §1-AC#7（下节）。

## AC#7（方法裁定：`mv`＋逐枚 md5 还原 算不算 `-overlay` 那句字面）＝**实质成立／字面不成立，且我认为字面那句本身写错了对象**
- 盘上事实（本腿自跑，`logs/10`）：`frontend/dist` 今天四枚＝`.gitkeep 0 d41d8cd9…`／`index-BRKjOIJ… db4db7a2… 49943`／`index-BVKlegVD.js 70128a3d… 553469`／`index.html f98bfc4b… 1044`
  与 `r1/logs/02-dist-baseline.txt` 那四枚 md5 **逐枚全等**；六族 porcelain 起手/终态都＝**0 行**。⇒ AC#7 后半句"终态 porcelain＝起手、被审文件与 `git show HEAD:` 逐串对回"＋票面 `:53` 的"dist 枚数与哈希起手/终态必须相同"**都满足**，我独立复核过。
- 但 `:53` 前半句字面是"**所有突变只走 `go test -overlay`（拷贝放仓外）**"。`274-r1` 走的是 `mv`＋还原，**字面不符**。我的裁量：
  ① 本票的突变对象是 **PowerShell 流水线的运行环境**（PATH 里有没有 node/npm、dist 目录里有没有字节），`go test -overlay` 在语义上**够不着**这一层——overlay 只替换 Go 编译器读到的文件内容，不改 `Get-ChildItem`/`Test-Path` 看到的真实目录树。⇒ 拿它作唯一合法手段，是**把 Go 侧尺的规矩套到构建脚本侧**。
  ② `:53` 的字面要求真正的**目的**在票面 `:47` 已经写明：「⛔ 绝不许真删/真盖本机 `frontend/dist` 里那 4 枚文件」，而那一格给的是**两条**合法路线——「**走 `go test -overlay` 或仓外合成 bundle**」。⇒ 两格不自洽：AC#1 允许仓外台件，AC#7 却只允许 `-overlay`。
  ③ 本腿的台件就是那条被允许的第二路线：仓外**硬链接镜像树**＋树内 `unlink`（`logs/01`：仓内四枚 md5 在我 unlink 前后逐枚相同），⛔ 全程没碰过仓内 dist——可见"不动仓内字节"并不必然需要 `-overlay`。
  ⇒ **建议**：AC#7 记「实质达成（仓内字节零改动，已由第二方复核），方法偏离字面」；要严格执行字面，就得先把 `:53` 改成与 `:47` 一致的"overlay **或**仓外台件"——那是**票面文字＝契约面**，归人工批准，⛔ 本腿不改一字。
  我**没**能核到的一点：`274-r1` 三次 `mv` 期间，别的会话若正好读 `frontend/dist`，会读到锚文件形——这段窗口是否真被并发使用过，我无法从盘上判定（记进 §5）。

### AC#6（门禁四数）＝**成立（附两条必须并存的具名读数）**（全串行为主，独占整包窗口，原始 `-v` 留仓外 `/d/tmp/wisp274v1/gotest-{start,end}.txt`）
- `sh scripts/d22scan.sh` ⇒ **rc=0**，且它的第 1 步就是那颗正控（`scripts/d22scan.sh:20-24` 逐字
  「`runtests.sh -C tools/d22scan ./... - the seeded-violation positive control. It proves the gate CAN go red`」）：
  `runtests.sh: OK - packages=[./...] top-level: PASS=35 FAIL=0 SKIP=0, === RUN=77, '[no tests to run]'=0`（`logs/07:240-241`）⇒ 门还能红。
- `go vet ./cmd/wisp/ ./internal/panel/` ⇒ **rc=0**（`logs/07` §5）。
- `gofmt -l cmd/wisp internal/panel` ⇒ 列出 **`cmd\wisp\models.go`**；`gofumpt v0.12.0 (go1.27.1)`（`/d/work/base/gopath/bin/gofumpt.exe`）`-l` 同一枚（`logs/09`）。
  红因＝该文件**盘上 CRLF**（CR 计数 334），⛔ 不是 274（`git diff --name-only 35633445^ 35633445` 只回 `scripts/build.ps1`；两份 blob `cmp` 回 IDENTICAL）。⇒ 见 §3.5。
- 整包 `go test -count=1 -v ./cmd/wisp ./internal/...` 两发（起手／终态，中间本腿只写过 `.scratch/**`）：
  **四数起手＝`RUN=1839 PASS=1234 FAIL=6 SKIP=7`（rc=1）；终态＝同一串四数（rc=1）**；
  逐名红名册起手 6 枚＝终态 6 枚（`logs/11-ac6-red-roster-start.txt` vs `logs/13-ac6-red-roster-end.txt`）：
  `TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestC21DesignTokensFourWayAgree`／`TestC21TableColourRowsMatchTokensCSS`／
  `TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestResolvePerCallBudget`；
  两发的 `diff` **只剩计时括号差异**（`(0.00s)`→`(0.01s)` 等 4 行，`logs/15`），**名字集合作差＝0**，失败包集合起手＝终态（`cmd/wisp`／`internal/ball`／`internal/panel`／`internal/risk`）。
  ⇒ 作差为 0 是靠**没有放宽任何断言**得到的：我没改过一行 Go，也没改过任何 `.go` 门禁。
  这 6 枚的因（`logs/12` 逐字）都不属本票：前 5 枚读 `design/**` 与 `frontend/**`（`design/assets/tokens.css` 现量**盘上不存在**、`git status --porcelain -- design` 里 16 枚 ` D`），
  `TestComposerContractTypesMatchFrontend` 自己写「`red as required: 1 line(s) matched, first=composer.tsx:4: <button onClick={() => bridge?.postMessage(JSON.stringify({ method: "approval.decide"...`」＝D22 禁例在页面源件上按设计红；
  `TestResolvePerCallBudget` 报「`risk: expansion moved this path onto a tree the caller did not name: C:`」＝`%TEMP%` 路径形状敏感。⇒ **一律登记为别人在飞的数，不算 274 的账**。

## §5 未做完的格（逐枚具名，⛔ 不写成通过）
1. **AC#0／AC#2／AC#8**：三格都不属本腿射程——AC#0(a)(b) 与 AC#2 的「CI 时长增量」要**推送后的 run 侧取数**（机主 10-07 说暂不推），AC#8 是 `SPEC-11 §2.2` 归属且票面写明**本票不选形**。⇒ 判不动。
2. **AC#4 的出处那一支我没有把它做成判据**：M-iii 只证了「陈旧页面今天能骗过闸门」，我没去改 `build.ps1` 加名册快照对比（那不是验收腿的写面）。⇒ 结论落 §2，修复归后续写腿。
3. **6 枚红名册到包的逐枚归属**：我给了失败包**集合**与红名**集合**，没有逐枚 `go test -run` 定位（预算帽）；`awk` 那次映射不可靠已弃用，不作数。
4. **AC#7 的并发窗口**：`274-r1` 三次 `mv` 期间是否有别的会话正好读过 `frontend/dist`，盘上不可判。⇒ 只能问编排者/别的腿。
5. **CI runner 的 PATH 里有没有 node、`actions/checkout` 落什么行尾**：仍＝〔待推送取数〕，乙-1 落地后 `test-windows`／`slo-smoke` 会不会因此变红，本腿判不动（票面 `:95` 那句代价没被关闭）。

## §6 终态自证
- **porcelain vs 起手**：起手六族＝**0 行**（`logs/00-anchor.txt`）。我这一腿**只写过** `.scratch/wisp/probes/274/v1/**`。
  ⚠ 现在同一条尺回 **` M scripts/check-path-length-budget.sh`**（`ls` 现量 mtime **14:55**，`git diff --stat` ＝ `7 +++++++`，1 file changed）＝
  **别的会话在我跑整包期间落的**，⛔ 不属本腿、也不属被审件；我不动它、不把它算进任何作差（它也不在 Go/构建胶水的被审面上）。
- **被审文件字节尺**：`scripts/build.ps1` 盘上 md5 **`644c2f6af1fa3fc7f8b3ab38c20753f5`** ＝ `git show 35633445:scripts/build.ps1 | md5sum` **同一串** ⇒ 我审的就是那一笔的落地件；
  前一份 `35633445^` 那份 md5 ＝ `9dee32a3c73317a94df456f5208cf1cc`（172 行）。相关尺具 md5：`cmd/wisp/panel_host_gate_test.go 332cf096…`／`internal/panel/assets.go f9e3b3b0…`／`frontend/embed.go d6484e65…`（⛔ 三枚都被本腿**只读**：`assets.go`/`panel_host_gate_test.go` 我一行没改，`frontend/**` 除 `embed.go` 的既有引用外没读页面代码）。
- **受保护路径起手/终态相同**：`frontend/dist` 四枚 `f98bfc4b…`／`db4db7a2…`／`70128a3d…`／`d41d8cd9…` ＝ 起手基线 ＝ `r1/logs/02-dist-baseline.txt`；`build/wisp.exe e6c8e52b…` ＝ 起手（本腿没重跑过仓内 build.ps1 到 `go build` 那一步）。
- **本腿 commit 名册**（⛔ 无 push、无 amend）：
  `5da17c58`（§0 起手锚）→ `ec1ccddd`（§1 AC#1＋AC#4 判据）→ `3f7bd77d`（§2 五发突变）→ `15ff0846`（§3/§4/AC#7）→ `955857e9`（13 份 logs）→ 本笔（AC#6 名册＋§5/§6）。
- **仓外临时件（只建不删）**：`D:/tmp/wisp274v1/`（镜像树 `tree/`、`fakebin/npm.cmd`、`bin-nodeonly/node.exe`、`gofmt/`、`gotest-start.txt` 816KB、`gotest-end.txt`、`crlf-go-files.txt`、`mutate.sh`、各 msg 文件）。

## 裁决摘要（给编排者翻勾用，本腿一枚不翻）
| 格 | 裁 | 一句凭据 |
|---|---|---|
| **AC#1** | **成立** | 未修码那份 build.ps1 整条 `EXITCODE=0`（`logs/03:12`）而同棵树 exe `assets NOT BUILT` rc=1（`logs/03:20-21`）、embed 尺回 `[dist/.gitkeep]` rc=0，正控同一把尺回 4 枚（`logs/01`） |
| **AC#4** | **字面成立＋具名残留缺口** | M-i／M-ii／M-iv 三发各自报不同具名原因且码＝1（`logs/04`）；自造 M-iii「npm 什么都没写、dist 里躺着陈旧页面」**落绿 rc=0**（`logs/06:16,23,28`）⇒ 闸门校存在不校出处，那半属 AC#9 |
| **AC#6** | **成立（附读数）** | d22scan 正控先绿（PASS=35/FAIL=0/SKIP=0）再 rc=0；`go vet` rc=0；gofmt/gofumpt v0.12.0 各列同一枚 `cmd\wisp\models.go`（CRLF，非 274）；整包四数起手＝终态 `RUN=1839 PASS=1234 FAIL=6 SKIP=7`，逐名红名册作差＝0 |
| **AC#7** | **实质成立／字面不成立**（方法裁定建议） | dist 四枚 md5 今天＝`r1/logs/02` 基线；`mv`＋md5 还原达成的是 `:47` 允许的"仓外/不动仓内字节"这一目的，但 `:53` 的字面 `-overlay` 对 PS 流水线突变够不着，且与 `:47` 自相矛盾 ⇒ 要不要改票面文字＝人工批准 |
| AC#0／AC#2／AC#8 | **判不动** | 要推送后的 run 侧取数／属 spec 归属，本票不选形 |
