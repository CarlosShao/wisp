# 274 · 三格逐格收口裁决（收口腿 `274-c1`；2026-11-14 14:0x +08）

**腿性质**：只读取证 ＋ 本件／票面只追加。零产码、零测试脚本改动、零 `ci.yml`／`build.ps1`／`docs/specs/**` 改动、零 push。
**锚点**：HEAD `069f97cac`（起手 `git status --porcelain -- cmd internal scripts tools .github docs frontend` ＝ 14 行，全部是别人的在飞件：1 枚 ` M frontend/fixtures/composer-states.html` ＋ 13 枚 `??`）。
**两发绿 run（本件全部 CI 读数的唯一来源）**：`43559883923`（head `b52365b71`）／`43561054718`（head `df6a8e035`），各自 6/6 success，逐 job ＝ `lint`／`slo-smoke`／`slo-full`／`test-windows`／`test-core`／`lint-frontend`。
**对照发（pre-甲，唯一的"改前"读数）**：`37703959747`（head `cc31526`，`cgo build smoke`／`Build wisp.exe` 两枚步 conclude=success；该 run 整发红是别的原因，与本件无关）。
**框数**：本件 ⛔ 一枚 AC 框都没翻（起手＝终态＝**9 勾／3 未勾**，尺＝`grep -cE '^- \[x\] \*\*AC#'` = 9、`grep -cE '^- \[ \] \*\*AC#'` = 3）。翻勾归编排者。

---

## 0. 三格一句话结论

| 格 | 判 | 一句话 | 下落 |
|---|---|---|---|
| **AC#0** | **能结** | 两枚数都给到：(a) `npm run build` 产出 **3 枚**（`dist/index.html` 1.04 kB ＋ `dist/assets/index-CZoXrb9N.css` 65.26 kB ＋ `dist/assets/index-zFGXlVs5.js` 640.29 kB），落 `frontend/dist/`，job 结束即丢；(b) 同一发 run 里 `build.ps1` 那份 dist ＝ **4 枚**（含 0 字节锚 `.gitkeep`），同 run 的 `AC#12` 读数 **`shape=page-bundle built=true … manifest-entries=4`**，对干净检出形（1 枚锚）**差 3 枚页面产物**。⚠ 口径差一行：字面要求的 `wisp.exe panel-assets` 那一条命令在 CI 里**没有**读数（见 §1.3） | 读数在 §1 |
| **AC#2** | **能结** | 三形四栏齐：甲已落地（写面只 `scripts/build.ps1` 1 枚；CI 增量＝前端段**实测** 13.4–24.4s（hosted）／69.7–82.5s（self-hosted）逐 job，两发 run 各一遍）；乙未选（+1 artifact ≈706,670 B、+3 条 `needs:`，且与已落地的出处闸门**正冲**）；丙未选（`release` job 在 `ci.yml` 里 **0 命中**，无栖息地） | 表在 §2 |
| **AC#8** | **归口** | 裁决对象是 `docs/specs/**` 文字（`SPEC-11 §2.2:44/:45`、`§3.2:69-84`，连带 `SPEC-01:70/:82`）＝**机主批**。本件给可直接批的提案（P1 推荐／P2 备选，逐句"现文→改文→为什么"），**票面那格维持不勾**、规格一字未改 | 提案在 §3 |

---

## 1. AC#0 —— 逐字原文 ＋ 判 ＋ 两枚读数

> **AC#0 先把两枚数量清（不许引用本票面的数）**：(a) CI 的 `lint-frontend` 那一发真跑成功后，`npm run build` 到底产出几枚文件、落在哪；(b) 同一发 run 里 `build.ps1` 造出的 exe，跑 `wisp.exe panel-assets` 回的是 `built=true` 还是 `false`、枚数差几枚。⛔ 台件不许写进 CI 之外看不到的地方。

**判：能结**（两枚数都给了；字面口径差一行已在 §1.3 具名，翻勾与否请按票内 AC#7"实质达成／方法偏离字面"的先例裁）。

### 1.1 数 (a)：`lint-frontend` 的 `npm run build` 产出

两发 run 逐字同（连哈希名都同，说明页面源码自 09-26 未变）：

```
lint-frontend  build (vite build -> frontend/dist, the bytes go:embed carries)
  2026-11-14T05:06:42.7067362Z ✓ 2454 modules transformed.
  2026-11-14T05:06:42.8247869Z dist/index.html                   1.04 kB │ gzip:   0.63 kB
  2026-11-14T05:06:42.8248681Z dist/assets/index-CZoXrb9N.css   65.26 kB │ gzip:  12.28 kB
  2026-11-14T05:06:42.8249440Z dist/assets/index-zFGXlVs5.js   640.29 kB │ gzip: 193.92 kB
  2026-11-14T05:06:42.8327867Z ✓ built in 375ms
```

- 枚数＝**3**；落点＝**`frontend/dist/`**（该 job 有 `defaults.run.working-directory: frontend`，`ci.yml:1543-1545`）。
- 第二发（run `43561054718`，job `132344591193`）逐字同：`dist/index.html 1.04 kB`／`dist/assets/index-CZoXrb9N.css 65.26 kB`／`dist/assets/index-zFGXlVs5.js 640.29 kB`，`✓ built in 442ms`。
- **它去哪了＝哪也没去**：`ci.yml` 里 `uses: actions/upload-artifact` 只有 `:1415`／`:1490` 两枚，内容都是 `build/slo/slo-report.json`；`^  needs:` **0 枚** ⇒ 该 job 的 dist 随 job 结束丢弃。
- 时长：`lint-frontend` 整 job 27s（run1）／35s（run2）；其中 `npm ci` 3s/4s、vite build 步 4s/6s。

### 1.2 数 (b)：同一发 run 里 `build.ps1` 造出的 exe 里的页面

`build.ps1` 在 CI 共 **3** 处（`ci.yml:1035` test-windows／`:1408` slo-smoke／`:1484` slo-full），三枚 job 的日志**逐字同**（run1 与 run2 也同）：

```
build.ps1: frontend: node = C:\Program Files\nodejs\node.exe (v22.23.4); npm = C:\Program Files\nodejs\npm.cmd (10.9.9)
build.ps1: frontend: node_modules ABSENT -> running npm ci (clean install pinned by package-lock.json).
build.ps1: frontend: pre-build dist snapshot = 1 file(s) [.gitkeep] (this run's npm must produce an artifact name outside it).
build.ps1: frontend ok: 4 file(s) in frontend/dist (entry index.html is 1068 bytes):
  .gitkeep=0 assets\index-CZoXrb9N.css=65264 assets\index-DmXz6VZ6.js=640338 index.html=1068
  || provenance: this run's npm added new artifact name(s) [assets\index-CZoXrb9N.css assets\index-DmXz6VZ6.js index.html]
     over the pre-build snapshot 1 file(s) [.gitkeep]; gone since snapshot: []
```

同 run、同 head 的 `built` 判定行（`test-windows` 的 `cmd/wisp CLI tests` 步，`panel_host_gate_test.go:118`）——**这是本件给 `built=true` 的那一条读数**：

```
panel_host_gate_test.go:118: AC#12 reading (head b52365b7): shape=page-bundle built=true
  entry-bytes=1068 entry-ctype="text/html; charset=utf-8" entry-err=<nil> refs=2 check-err=<nil>
  manifest-entries=4 manifest-err=<nil> | git-metadata=true tracked=1 tracked-beyond-anchor=0
  tracked-has-entry=false ignored-or-untracked=2
```

第二发（head `df6a8e03`）逐字同：`shape=page-bundle built=true entry-bytes=1068 … manifest-entries=4`。

**两枚数**：

| 问 | 读数 |
|---|---|
| `built=true` 还是 `false` | **`built=true`**（`manifest-entries=4`、`entry-bytes=1068`；`built` 的定义＝`internal/panel/assets.go:56-57` 里 `fs.Stat(tree, "index.html")` 成功即 `true`，本趟 dist 名册含 `index.html=1068`，两处对得上） |
| 枚数差几枚 | **4 − 1 ＝ 差 3 枚页面产物**（同 run 的 exe 内嵌 4 枚＝3 枚页面产物 ＋ 0 字节锚 `.gitkeep`；干净检出形＝1 枚锚 `.gitkeep`） |

**反向对照（同一把仪器在"没页面"那一形上的读数，pre-甲 run `37703959747`／head `cc31526`）**——证明这把尺不是恒真：

```
build.ps1: frontend step skipped (no frontend yet; embed lands in S5 per SPEC-11 §2.2)
panel_host_gate_test.go:109: AC#12 reading (head cc31526): shape=anchor-only built=false
  entry-bytes=0 … manifest-entries=0 … ignored-or-untracked=0
```

⇒ 同一枚 `AC#12`、同一台 hosted Windows runner：改前 `anchor-only/built=false/manifest-entries=0`，改后 `page-bundle/built=true/manifest-entries=4`。**这就是 AC#4 要的反形敏感性，而且这次是 CI 侧的读数，不是本机台件。**

### 1.3 口径差一行（具名，⛔ 不粉饰）

- AC#0(b) 字面写的是"跑 `wisp.exe panel-assets` 回的是 `built=true` 还是 `false`"。**CI 今天不对出炉的 exe 跑这条命令**：`panel-assets` 在 `ci.yml` 里只出现 **1 处、且是注释**（`:1587`，L2 渲染证据那段拿它当输入样例）；两枚 `upload-artifact` 只传 SLO JSON ⇒ **`build/wisp.exe` 无人留档**，外部拿不到那枚 exe 去补跑。
- ⇒ 本格给的 `built=true` 走的是**同 run 的两条读数**：`build.ps1` 的 dist 名册行（embed 的**输入**）＋ `AC#12` 的 `built` 判定行（同一棵 dist 被 `//go:embed all:dist` 收进同 run 的测试二进制）。两者与 head、与 run 都同源，但**不是** `panel-assets` 那一行。
- **若编排者要字面那一行**（可选加固，本件不实施，因 `ci.yml` 不在本腿写面）：在 `slo-smoke` 的 `Build wisp.exe`（`:1408`）之后加一发 `./build/wisp.exe panel-assets -check`，**本机现量 150 ms、exit=0**（见 §4）。它只证明"entry 与 2 枚 refs 可解析"，不比上面两条多信息，但会把 `panel-assets` 这个词印进日志。

---

## 2. AC#2 —— 逐字原文 ＋ 判 ＋ 三形四栏表

> **AC#2 三形代价表，⛔ 本票不选形**：甲＝`build.ps1` 第 2 步真跑 `npm ci && npm run build`（代价：Windows 产线与每台编译机从此**硬需要 Node**，与 34 里"host needs no Node"相冲，CI 时长代价要现量）；乙＝CI 加一枚跨 job 产物传递（`lint-frontend` 上传 dist → Windows job 下载后再 `build.ps1`，exe 才真带页面；代价：新 artifact 依赖＋job 顺序）；丙＝产线不动，只加一枚**发布前闸门**（打包/签名之前强制 `wisp panel-assets` 回 `built=true` 且入口字节>0，否则红）。逐形给：动哪几枚文件、要不要新 job/新 artifact、CI 时长增量、**"这条边以后谁能看见"**（本机？CI？只有发布步？）。

**判：能结**——四栏齐；其中"CI 时长增量"这一栏过去唯一欠的就是甲（`274-a1` §4 具名"要现量"），**今天两发绿 run 逐 job 给到了**。⛔ 两发之外的"形"栏是量级/下界，已在表下具名。

### 2.1 三形四栏表

| 形 | 动哪几枚文件 | 新 job / 新 artifact | CI 时长增量（现量） | "这条边以后谁能看见" |
|---|---|---|---|---|
| **甲**（**已落地**：`35633445`→`91c90aa1`→`674c1920`，三笔写面各只 `scripts/build.ps1` 一枚；现 370 行） | `scripts/build.ps1` **1 枚**（`scripts/**` 之外 0 枚，`ci.yml` 0 枚） | 无 / 无 | 前端段逐 job（`npm ci`＋快照＋`npm run build`）：run1＝**13.4s** test-windows·**20.5s** slo-smoke·**82.5s** slo-full；run2＝**21.7s**·**24.4s**·**69.7s**。其中 `npm run build` 本身 5.4–20.5s，`npm ci` 8.0–62.0s（这三枚 job **都没有 npm cache**：`cache: npm` 全文件只 `:1554` 一枚，在 lint-frontend） | **本机 ＋ CI**：三枚产 exe job 每发都在日志里印 `build.ps1: frontend ok: 4 file(s) …` 那一行；本机同一行也印 |
| **乙**（未选；盘上不存在） | `.github/workflows/ci.yml` **1 枚**（lint-frontend 加 1 枚 upload；三枚 Windows job 各加 1 枚 download；另加 **3 条 `needs:`**——今天 `^  needs:` ＝ **0**） | **1 枚新 artifact**（dist ≈ **706,670 B**＝65264＋640338＋1068）＋ 3 条新依赖边 | upload/download 按现量级 **1–2 s/次**（现两枚 SLO 报表 artifact 的 upload 步＝1 s/1 s（run1）、2 s/2 s（run2））；另加 lint-frontend 整 job **24–35 s** 的起步等待（test-windows 自身 1267 s/1928 s ⇒ **不落关键路径**） | **只有 CI**（artifact 页）；exe 仍无人留档 |
| **丙**（未选；今天连栖息地都没有） | 新 job／新 workflow（`release` 在 `ci.yml` 全文 **0 命中**，连这个词都不出现） | **1 枚新 job** / 无 | 新 job 起步 ≈ **30–40 s**（参照现 `lint-frontend` 整 job 24 s/35 s）＋ 闸门本身 **150–197 ms**（本机现量，见 §4） | **只有发布步**（今天不存在） |

### 2.2 甲这一栏的三条附带读数（都是现量，不是推断）

1. **"硬需要 Node"已经从风险变成实测**：hosted Windows runner 作业 PATH 上有 `node v22.23.4` ＋ `npm 10.9.9`（两发 run 逐 job 打印）；self-hosted 那台是 `D:\work\server\node24\node.exe (v24.21.0)` ＋ `npm 11.19.0`。⇒ 甲落地后**两发都绿**，Node 这条硬依赖在 CI 三枚 job 上都有。
2. **乙与已落地的牙正冲**（这条改变了乙的代价表）：乙的原形是"下载现成 dist 再让 `build.ps1` 跳过第 2 步"，但甲落地后的 `build.ps1:211` 出处闸门要求"本趟 npm 之后名册里出现快照外的新名字"。一份预先下载的 dist 恰恰**是快照里已有的名字**⇒ 要么被闸门判红（`:211`），要么又得给 `build.ps1` 开口子（＝把 AC#10/#11 的牙拔掉）。⇒ **乙今天不是"加两根线"，还要动出处闸门语义**。
3. **整步两头发不能当甲的增量**（防一个假归因）：pre-甲那一发 `slo-full` 的 `Build wisp.exe` 整步 **28 s**，今天 **188 s /169 s**；其中前端段只占 **82.5 s/69.7 s**，余量是 go 编译/缓存态差（两次之间隔了 5 周产码变动）。⇒ 甲的增量只能按**前端段**算，⛔ 不按整步差算。

### 2.3 具名残留（不影响"能结"，但不许读成实测）

- 乙、丙的"CI 时长增量"是**组成量级**（各形未落地，无法整形实测）；甲那一栏是**整形实测**（两发 run 逐 job）。
- 乙的 upload/download 秒数是拿现两枚 SLO 报表 artifact 的步时长作参照，**payload 不同**（SLO JSON 数 KB vs dist 706,670 B）⇒ 只能说量级，不能说等值。

---

## 3. AC#8 —— 逐字原文 ＋ 判（归口）＋ 可直接批的提案

> **AC#8 `SPEC-11 §2.2` 那两步的归属裁决（⛔ 本票不选形，由 `274-a1` §3 现量升格而来）**：§2.2 第 2 步逐字要求「存在则跳过；**`--with-frontend` 时走 `docker/frontend.Dockerfile`（§3.2）**」，而今天 ①`build.ps1` **没有这枚开关**（无条件跳过）②`docker/frontend.Dockerfile` 跟踪枚数 **0** ③该节第 3 步还写着「embed **`assets/web`** 已就位」，而 `assets/web` 跟踪枚数 **0**（真实通道是 `frontend/dist`）。⇒ 要裁的是：**这枚开关该由 `build.ps1` 兑现，还是承认 §2.2 那两步已被"CI 里 `lint-frontend` 跑 vite ＋ `frontend/dist`"这一形取代**（取代＝要动 spec 文字＝人工批准，另案；⛔ 任何腿不许"顺手补个开关"来让代码与 spec 对上）。凭据＝`.scratch/wisp/probes/274/a1/census.md` §3（我现读复认，见本票末节）。

**判：归口**——对象是 `docs/specs/**` 的文字＝**机主批**（`SPEC-12 §4.1`／AGENTS.md 第 4 条三类之一）。本腿**规格一字未改**（`git status --porcelain -- docs/specs` 起手＝终态＝0 行），票面该格**维持不勾**。下面是可直接批的提案。

### 3.1 今天盘上的事实（逐条带尺，全部本腿现跑）

| # | 事实 | 尺与读数 |
|---|---|---|
| 1 | §2.2 第 2 步现文 | `SPEC-11:44`：`  2. 前端产物：存在则跳过；--with-frontend 时走 docker/frontend.Dockerfile（§3.2）` |
| 2 | §2.2 第 3 步现文 | `SPEC-11:45`：`  3. go build（cgo: CGO_ENABLED=1, CC=mingw32-gcc; embed assets/web 已就位）` |
| 3 | `--with-frontend` 在产线里不存在 | `Select-String scripts/build.ps1 'with-frontend'` ＝ **0 命中**；`build.ps1` 的 `param` 只有 `-Env dev|prod`（`:35-38`） |
| 4 | `docker/frontend.Dockerfile` 不存在 | `git ls-files docker` ＝ **15 枚**，名册里**没有** `frontend.Dockerfile` |
| 5 | `assets/web` 不存在 | `git ls-files | grep ^assets/web` ＝ **0**；`Test-Path assets/web` ＝ **False**；但 `.gitignore:33` 留着一枚死规则 `assets/web/dist/` |
| 6 | 权威层从未写过这两步 | `docs/PLAN.md`：`frontend.Dockerfile`／`assets/web`／`with-frontend` 各 **0 命中** |
| 7 | 真实通道 | `frontend/embed.go:19` `//go:embed all:dist` → `internal/panel/assets.go:44` `fs.Sub(frontend.Dist(), "dist")`；`built` 只看 `index.html` 在不在（`internal/panel/assets.go:56-57`） |
| 8 | 代码与 §2.2 第 2 步**反向** | `build.ps1:86-215` 今天**无条件**跑 `npm ci`＋`npm run build`，且有出处闸门（`:211` 没有新名字就 FATAL）；`§2.2` 却写"存在则跳过" ⇒ 不是"没实现"，是"实现成了反的" |

### 3.2 提案 P1（**推荐**）：承认取代，把 spec 改成今天真的那一形

一次性提交、共 6 处（5 处在 `SPEC-11`，1 处连带的 `SPEC-01` 树，另 1 处 `.gitignore` 死规则）。**⛔ 本腿不动手，等机主批。**

**P1-1 `docs/specs/SPEC-11-build-deploy-containerization.md:44`**

- 现文：`  2. 前端产物：存在则跳过；--with-frontend 时走 docker/frontend.Dockerfile（§3.2）`
- 改成：`  2. 前端产物：在 frontend/ 下跑 npm ci（node_modules 不存在时）＋ npm run build，产出 frontend/dist，并校验其中有本趟构建新增的页面名册；没有 node/npm 或构建失败＝整条构建失败（无跳过路径）`
- 为什么：这就是 `build.ps1:86-215` 今天的行为（含 `:211` 出处闸门）；"存在则跳过"若留着，会让下一读者以为可以拿一份现成 dist 蒙过去——而那一形今天**必然红**（AC#10）。

**P1-2 `SPEC-11:45`**

- 现文：`  3. go build（cgo: CGO_ENABLED=1, CC=mingw32-gcc; embed assets/web 已就位）`
- 改成：`  3. go build（cgo: CGO_ENABLED=1, CC=mingw32-gcc; embed frontend/dist 已就位，见 frontend/embed.go）`
- 为什么：`assets/web` 跟踪 0 枚、目录不存在；`//go:embed` 的字面路径是 `all:dist` 且包在 `frontend/`。

**P1-3 删掉整节 `SPEC-11:69-84`（`### 3.2 frontend.Dockerfile（前端构建，S5 起需要）`）**

- 为什么要删：该节整段承诺一枚**不存在**的文件（`docker/frontend.Dockerfile`：`git ls-files docker` 15 枚里没有），还写着"产物落 `assets/web/`"（不存在）与"S5 之前 frontend 目录可空，构建脚本跳过该步"（今天 frontend 有 64 枚源码、构建脚本**不跳过**）。**只改 `:44` 而留着 §3.2，等于把同一句谎话留在下一节**——这正是本票立案时 `build.ps1:74` 那句的成因（它当时引的就是这一节）。
- 若机主想保留"将来若要宿主机免 Node"的设想：改为一段**明确标注未实现**的备忘（`⚠ 未实现：无 docker/frontend.Dockerfile，无 assets/web 通道；今天的通道是 frontend/dist`），也行；但**不能保持现在这种"写着像已有"的时态**。

**P1-4 `SPEC-11:125`（§6 表 `frontend` 行）**

- 现文：`| \`frontend\` | ubuntu + node 镜像 | tsc/lint/build + 硬编码色值/emoji 扫描 | 阻塞（S5 起） |`
- 改成：`| \`lint-frontend\` | ubuntu-latest + actions/setup-node（node 24） | tsc/lint/build（vite → frontend/dist）+ 渲染证据 + token 漂移 | 阻塞（已在跑） |`
- 为什么：job 真名是 `lint-frontend`（`ci.yml:1535`），没有 node 镜像（`runs-on: ubuntu-latest` ＋ `actions/setup-node@v4`，`:1549`），"S5 起"过期（该 job 自票 77 起就在跑，今天 6/6 绿里的一枚）。

**P1-5 `docs/specs/SPEC-01-architecture.md:70` 与 `:82`（连带，同一批批）**

- `:70` 现文：`├── assets/web/embed.go                # //go:embed frontend 产物（…D29）`
  改成：`├── frontend/embed.go                   # //go:embed all:dist（面板产物；D29）`
- `:82` 现文：`│   ├── frontend.Dockerfile            # Node 20 构建镜像（输出 dist）`
  改成：删掉该行（`docker/` 下没有这枚文件）。
- 为什么：P1-1/2/3 一改，`SPEC-01` 的目录树就成了**唯一还在说有 `assets/web/embed.go` 与 `frontend.Dockerfile`** 的地方；同一批批掉，省一次两处对不上的窗口。

**P1-6 `.gitignore:33`（死规则，低利害，可同批或自办）**

- 现文：`assets/web/dist/`
- 改成：删掉这一行。
- 为什么：该路径不存在，规则是给一枚从未出生的通道预留的；留着会让下一位以为 `assets/web` 是真实落点。

### 3.3 提案 P2（**备选，不推荐**）：兑现 `--with-frontend` 那枚开关

要动的东西（逐项，全部可量化）：

1. 新建 `docker/frontend.Dockerfile`（跟踪 0 → 1 枚），内容就是 `SPEC-11:71-81` 那段；
2. `scripts/build.ps1` 加 `param([switch]$WithFrontend)`，并把第 2 步改成"默认跳过／带开关才构建"——**这一步直接拆掉 AC#10/#11 刚立的牙**（跳过＝不跑 npm＝没有名册对拉，出处闸门失效）；
3. embed 路径三处改：`frontend/embed.go:19`（`all:dist` → `all:../assets/web` 之类，超出包边界，实际要搬家）、`internal/panel/assets.go:44` 的 `"dist"`、以及 `.gitignore:33` 从死规则变成活规则；
4. 票 §3 的两条 ⛔ 正面相撞：⛔ 不改 `internal/panel/assets.go` 的 fail-closed 语义、⛔ 不动 `frontend/.gitignore` 的 `dist/*`（票 77 AC#1）；
5. CI 侧：三枚 Windows job 本来就没有 Docker（hosted runner 有，但 self-hosted 那台未必有 docker desktop；今天 `slo-full` 用的就是 `powershell` 直跑）。

⇒ P2 买到的唯一东西是"宿主机无需 Node"，而 `PLAN.md` **从未写过**这条；换来的代价是"第 2 步可跳过"＋"闸门失效"＋"embed 通道搬家"。**不选它的理由不是"动作大"，是它买不到价值**（而且与甲已落地的实测绿相反）。撤销口令（若机主最后选 P2）：`「274 AC#8 选 P2」`——重新开一枚票，本票不转。

### 3.4 批准后的执行（谁做、几步、撤销口令）

- 批 P1 后：一枚写腿只改 4 枚文件（`SPEC-11`／`SPEC-01`／`.gitignore` ＋ 票面 AC#8 翻勾），跑 `sh scripts/d22scan.sh`（`docs/**` 在扫描面内）＋ 一次 `git status --porcelain -- docs/specs` 作差；**不改产码、不改 CI**。
- 撤销口令：`「274 P1 撤回」`＝把上面 6 处按 `git show <改前>:<路径>` 逐字节退回（改动只在一枚 commit 里，可整笔 revert）。
- ⚠ 翻勾动作仍需**非实现者**复认（`SPEC-12 §4.3`）。

---

## 4. 复跑命令（每条尺，逐条可重跑）

```bash
# 两发 run 的逐 job 结论与时长（JSON 已存 .scratch/wisp/probes/274/c1/logs/）
gh run view 43559883923 --json headSha,conclusion,jobs
gh run view 43561054718 --json headSha,conclusion,jobs
gh run view 37703959747 --json headSha,conclusion,jobs          # pre-甲 对照

# 四分 job 的完整日志（⚠ 不用 --log-failed，它会截断测试名）
gh run view 43559883923 --job 132340969677 --log   # lint-frontend
gh run view 43559883923 --job 132340969645 --log   # test-windows（含 AC#12 行）
gh run view 43559883923 --job 132340969630 --log   # slo-smoke
gh run view 43559883923 --job 132340969634 --log   # slo-full
gh run view 43561054718 --job 132344591193 --log   # 第二发 lint-frontend
gh run view 43561054718 --job 132344590973 --log   # 第二发 test-windows
gh run view 37703959747 --job 113073784763 --log   # pre-甲 test-windows（skipped 行＋anchor-only 行）

# 结构尺
grep -n 'npm run build' .github/workflows/ci.yml        # 1 枚 = :1578（lint-frontend 内）
grep -n 'File scripts/build.ps1' .github/workflows/ci.yml  # :1035 :1408 :1484 三处
grep -c '^  needs:' .github/workflows/ci.yml            # 0
grep -n 'uses: actions/upload-artifact' .github/workflows/ci.yml   # :1415 :1490，均 SLO JSON
grep -ci release .github/workflows/ci.yml               # 0
git ls-files docker | wc -l                             # 15，无 frontend.Dockerfile
git ls-files | grep -c '^assets/web'                    # 0
```

**本机读数（只读，未改任何字节）**：

```
build\wisp.exe panel-assets            → rc=0，197 ms（进程起停含在内），
  panel assets embedded: 4 files, entry=index.html built=true
build\wisp.exe panel-assets -check     → rc=0，150 ms，
  panel assets check: entry=index.html built=true 2 asset refs resolve
  [./assets/index-DpRdJ-7r.js ./assets/index-Dxt2y5oW.css]
```

⚠ 那枚本机 exe 的 mtime＝`2026-11-02 17:16:22`／32,966,840 B，内嵌名册（`index-DpRdJ-7r.js`／`index-Dxt2y5oW.css`）与今天 CI 的（`index-DmXz6VZ6.js`／`index-CZoXrb9N.css`）**逐枚不同** ⇒ 又一次印证票内那条口径：**跨上下文的名册不许当判据**，本机这行只用来给丙形的闸门耗时，不用于 AC#0。

---

## 5. 纪律自证（收件三把尺：提交真在／件真在／自己复跑）

- **写面**：本件 ＋ 票 274 末尾追加一节 ＋ `.scratch/wisp/probes/274/c1/**`（读数件）。`cmd/ internal/ scripts/ tools/ .github/ frontend/ design/ docs/specs/**` **零写入**；`git status --porcelain -- docs/specs` 起手＝终态＝**0 行**。
- **零测试**：⛔ 未跑 `go test`／`go build`／`go vet`（`cmd/wisp` 面全程没碰）；四条 `go list -f '{{.EmbedFiles}}' ./frontend` 级别的只读尺是 `go list`，不编译不落盘。
- **零 push**、⛔ 无 `add -A`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`；commit 走显式 pathspec。
- **框数**：起手 9/12（未勾 `AC#0`／`AC#2`／`AC#8`），本件交件后**仍是 9/12**——⛔ 一枚没翻（翻勾归编排者）。
- **六族 porcelain 两头（如实记一处别人的在飞变动）**：起手＝**14 行**（1 ` M` `frontend/fixtures/composer-states.html` ＋ 13 `??`）；本腿交件前现量＝**16 行**（3 ` M` ＋ 13 `??`），多出的两枚 ` M` 是 `cmd/wisp/resident_ball_live_228_windows_test.go` 与 `internal/ball/hotkey_refused_borrow_release_245_test.go`，出现在本腿**只读**期间，属别的腿在飞，本腿不改、不算进作差。
- **未采信任何腿的自报**：本件全部数字来自 `gh` 现拉＋本机现跑；引用旧件（`274-ac0-counts.md`／`near-done-triage.md`）只用于指出"哪些栏过去判过"，未拿它们的数当本件读数。
- **不采信通知**：本腿没有别人交来的件要收。
