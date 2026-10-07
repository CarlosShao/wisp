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
