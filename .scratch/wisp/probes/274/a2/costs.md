# 274-a2 costs.md — 票 274 `AC#0`／`AC#2`／`AC#8` 备料（只读现量，⛔ 不选形）

> 腿＝`274-a2`（只读普查）。起手锚＝`.scratch/wisp/probes/274/a2/00-anchor.md`（锚 HEAD `7a367a08`／分支 `dev`／
> `git status --porcelain -- scripts .github docs internal frontend` 起手 **0 行**／全仓 porcelain **755 行**／时刻 2026-10-08 10:54 +0800）。
> 本件每一条读数后面**都逐字写了所用的尺**；凡引用别的腿的旧读数，一律标 `[读数出处]` 并写清它的锚点，
> 因为 `scripts/build.ps1`（172→313 行）与 `.github/workflows/ci.yml`（941→**992** 行）在本票开过后都动过，
> **先前各腿给的所有行号都已过期**——本件里"今天"的行号一律是我现量的。
> ⛔ 本件不裁形、不新增判据、不改任何产码/spec/票面 AC 框。

---

## §1 先查重：`AC#2` 那张三形代价表，今天缺哪几栏

### 1.1 逐件读过名册（本腿实读，非引用）

| 件 | 锚点/时刻 | 它已经答了哪几格 |
|---|---|---|
| `.scratch/wisp/probes/274/a1/census.md`（215 行） | HEAD `15699a2f`，10-07 10:39 | §1 七族语义普查（谁在门控 `Built()`／`download-artifact` 0 枚／`panel-assets` 在 ci.yml 只有 1 枚注释）；**§4 已给甲/乙/丙三行的"动哪几枚文件·要不要新 job/artifact·时长量级·谁能看见·撞的钉"**；§3 `SPEC-11 §2.2` 逐字＋"S5 不在 §2.2"；§5 六条与派单/票面对账 |
| `.scratch/wisp/probes/274/a2/shipping-chain.md`（267 行＋18 枚 logs） | HEAD `933a8341`，10-07 12:2x | §1 出货通路唯一性（当时 `build.ps1:129`）＋param 只有 `-Env`；§2 三枚产 exe 的 job 逐枚点名（当时 `:571/:757/:820`）＋`needs:` 0／`upload-artifact` 两枚（当时 `:763/:826`）；§3 run 侧 node 读数（作业侧**取不到**、机侧 v24.9.0、runner 在 E: 且与本机同一枚）；§3.4 `slo-full` 日志第 183 行那句谎话原样印出＋`test-windows` 第 2908 行 `built=false` 却 PASS；**§4 仓外真构建 26s／4 枚／602,635 字节**；§4.3 本机 dist 是陈旧件；§5.1 两处顶回 |
| `r1/impl.md`·`r2/impl.md`（写腿）＋票末「收 274-r1/r2」 | 写面 `35633445`／`91c90aa1` | 甲形**已经落进盘上**（乙-1＝第 2 步真跑 npm＋失败即红）与出处闸门（快照对拉） |
| `v1/verdict.md`·`v2/verdict.md`·`v3/verdict.md`（非实现者） | 10-07 14:5x–17:4x | AC#1/AC#4（未修码两形＋定向突变）、AC#6 门禁四数、AC#9 同上下文对拉、**AC#10 五格全射程**、AC#11 数组对拉四发；v2 `logs/10·11`＝**真跑 npm 一发绿（rc=0）／同源码第二发实测红（rc=1）**；v1 §4 抽到 CI 侧 vite 产出 **3 枚** |
| `.scratch/wisp/probes/bundle/1/census.md`（更早的同题普查，票 §1 的来源同族） | 10-06 | spec `:44`/`:84` 逐字、`docker/` 目录名册、`PLAN.md` D41 整节"没有一个字讲产物由谁生成" |

⇒ **本腿不重做上面任何一格。**下面 1.2 只列"仍缺的栏"，每条都给了本腿新取的现量或具名写"欠什么"。

### 1.2 `AC#2` 三形代价表今天仍缺的栏（逐条：缺什么／本腿补到哪／剩下的欠谁）

1. **甲形：CI 时长代价"要现量"这一栏——到今天仍无 CI 现量。** 票面 `:49` 明写"CI 时长代价要现量"。
   - 本腿补上的**组成读数**（非本机那一半）：推送态最新一发 ci run `37703959747`（headSha `cc315261`，10-07 23:44Z）的**步级时间戳**现量：
     `npm ci`＝**5 秒**（23:44:38→43）、`build (vite build -> frontend/dist, the bytes go:embed carries)`＝**4 秒**（23:44:48→52）、`actions/setup-node@v4`＝1 秒、整个 `lint-frontend`＝23:44:32→23:44:55（**23 秒**，比 `274-a1` 给的"≈24-25s"更准）。
     同发 run 里三枚**产 exe 的 job 的 build.ps1 步基线**（旧码、无 npm）＝`test-windows` 23:46:21→23:46:56 **35 秒**、`slo-smoke` 23:45:23→23:45:44 **21 秒**、`slo-full`（self-hosted）00:14:27→00:14:55 **28 秒**。
     尺＝`gh run view 37703959747 --json jobs --jq '.jobs[] | select(.name|test("lint-frontend|test-windows|slo-smoke|slo-full")) | {name, startedAt, completedAt, steps:[.steps[]|{n:.name,s:.startedAt,e:.completedAt,c:.conclusion}]}'`（只读，rc=0，未拉全量日志）。
   - **仍欠**：Windows 那三枚**加上 npm 之后的净增量**＝只能由推送后那一发给（本机/自托管的间接数：`274-a2` §4.1 冷 `npm ci` 10 秒＋`npm run build` 9 秒；`v2/logs/10` 一发端到端 `REAL_BUILD_PROCESS_EXITCODE=0`，16:47:50 起、doctor 时间戳 16:48:01 ⇒ 同机整条 ≈11 秒量级，⚠ 那一发 `node_modules EXISTS` 走的是**跳过 npm ci** 那一支）。
2. **甲形：CI 上 `npm ci` 会不会每发冷跑、有没有缓存——先前各腿没量。** 本腿现量：
   `grep -n 'cache: npm' .github/workflows/ci.yml` → **仅 1 枚＝`:935`（在 `lint-frontend` 里）**；`grep -n 'runs-on:' ci.yml` → `:66 :396 :506 :791 :849 :923`（六枚 job，其中三枚产 exe）⇒ **三枚 Windows job 里没有任何 npm 缓存**，`frontend/node_modules` 被 `node_modules/`（根 `.gitignore:17`）忽略、不会随 checkout 出现 ⇒ 甲上 CI 的形是**每发冷 `npm ci`**。
3. **甲形在 self-hosted 那枚 runner 上的现场代价——先前只给了"runner 在 E:"，没给工作区形状。** 本腿只读 `ls` 现量（⛔ 没在那棵树里跑任何 `git` 命令，怕碰在飞的 job）：
   `E:/work/base/actions-runner/_work/wisp/wisp/frontend/dist` ＝**只有 `.gitkeep`**（dir 与 `.gitkeep` mtime 均 `2026-10-04 08:56`）；
   `.../frontend/node_modules` ＝**不存在**；`.../build/wisp.exe` ＝**30,471,221 B，mtime `2026-10-08 08:14`**（＝那发 `slo-full` 的 build.ps1 产物）。
   ⇒ 那一枚 runner 上甲落地后的第一发会走 **冷 `npm ci` ＋ 快照＝只有锚** 这一形（＝v2 实测绿的那一形）。
4. **乙形：artifact 依赖"点名是哪两枚 job"——有一半，且行号已过期。**
   - 既有先例（本腿现量）：`grep -n 'upload-artifact\|download-artifact' .github/workflows/ci.yml` → `:814`（`slo-smoke` 的 `Upload SLO report`）／`:877`（`slo-full`）／`:893 :894`（注释）；**`download-artifact` 0 枚**。旧读数 `:763/:826` 出自 `274-a2` §2.1，**已随 ci.yml 涨到 992 行而失效**。
   - 乙形自己需要的两枚（本腿现量落点）：上传点＝`lint-frontend`（job 声明 `:922`，`npm run build` 步 `:954`）之后；下载点＝三枚产 exe job 的 build.ps1 命令行之前（`:571`／`:808`／`:871`）；`needs:` 全文件仍 **0 枚** ⇒ 三条依赖边全是新增。
   - **仍欠**：artifact 的 path/保留期/pin 形状与"下载后再 npm"这一形与**已落地的甲形**的相互作用——后者本腿给了盘上事实（见 §3 乙栏：同名即红，`v2/logs/11` 那一发 rc=1 是同形实测），⛔ 但 CI 上的那一发要推送才响。
5. **丙形：闸门落点"在哪一步"——从来没被钉过。** 本腿补的料（§3 丙栏＋§5 名册）：
   `built=true` 在 exe 侧**只看入口存在、不看字节**（`internal/panel/assets.go:56-58` 现量 `fs.Stat(tree, EntryFile)` 即置 `built`）；
   全仓"入口字节>0"这一断言只有两处：`cmd/wisp/panel_host_gate_test.go:117-119`（Go 测试，读 embed 内字节）与 `scripts/build.ps1:191-194`（读 **dist 磁盘**字节，⛔ 不读 exe）；
   `wisp panel-assets` 的 rc=1 只在 `!Built()`（`cmd/wisp/panel_assets.go:126-129`）⇒ 票面 `:49` 那句"built=true **且入口字节>0**"today 在 exe 侧**没有现成判据**，只能解析 `-manifest` 第 2 列（`panel_assets.go:105-113` 打印 `name size hash`，本腿现跑过）。
6. **"这条边以后谁能看见"这一栏**：`274-a1` §4 三行都已给（甲＝本机＋所有跑 build.ps1 的机器／乙＝只有 CI／丙＝只有执行闸门那一发）。⚠ **甲那一行今天必须重述**：甲**已在盘上**，所以它不再是"以后谁能看见"，而是"现在谁能看见"＝本机＋三枚 CI job＋那枚 self-hosted runner（与本机同一台，`274-a2` §3.3 已坐实）。
7. **不缺的格（具名指回，⛔ 不重做）**：甲/乙/丙各自的"会不会撞既有钉"＝`274-a1` §4 三段末行＋票 §3 五枚 ⛔；`release` job 存不存在＝`274-a2` §5.1-1（本腿 §5 用今天的尺复认"6 枚 job、无 release"）；"exe 从未被任何步留档"＝`274-a2` §2.2（两枚 upload 只传 SLO JSON）；"那句谎今天真印进 CI 日志"＝`274-a2` §3.4-1（⚠ 那句 `Write-Host` 在甲落地后**已从盘上消失**，见本腿 §3 甲栏，这条旧读数今后不可再引）。

---

## §2 `AC#0` 两枚数量，在"零 push"这一前提下能取到多少

### 2.1 (a) `npm run build` 今天在本机产出几枚、落在哪

**⛔ 本腿没跑 `npm`（派单禁）**，所以"今天这一发"取不到；能取到的是**盘上现状**＋**历史真跑件**＋**CI 侧既有读数**：

| 维 | 读数 | 尺 |
|---|---|---|
| 本机 `frontend/dist` 盘上 | **4 枚文件**＝1 枚跟踪锚＋3 枚未跟踪产物 | `find frontend/dist -type f -printf '%p %s %TY-%Tm-%Td %TH:%TM\n'` |
| 逐枚 | `dist/.gitkeep` 0 B（mtime 10-07 11:57）／`dist/assets/index-BRKj5OIJ.css` **49,943**（09-27 10:59）／`dist/assets/index-BVKlegVD.js` **553,469**（09-27 10:59）／`dist/index.html` **1,044**（09-27 10:59） | 同上＋`stat -c '%n size=%s mtime=%y'` |
| 逐枚哈希（前 16 位 sha256，与 exe `-manifest` 同面） | `.gitkeep e3b0c44298fc1c14`／`index-BRKj5OIJ.css f8ff337b2ce80f3b`／`index-BVKlegVD.js 793e606366463f0e`／`index.html 2078af7f2343af94` | `for f in …; do printf '%s %s ' "$f" "$(stat -c %s "$f")"; sha256sum "$f" \| cut -c1-16; done` |
| 跟踪名册 | **1 枚**＝`frontend/dist/.gitkeep`；另两枚 `index.html`/`assets/` 被 `frontend/.gitignore:12 dist/*` 与根 `.gitignore:22 frontend/dist/*` 忽略（`git check-ignore -v` 现量） | `git ls-files frontend/dist`／`git check-ignore -v frontend/dist/index.html frontend/dist/assets` |
| 落在哪 | `frontend/dist`（vite 自报路径 `dist/…`）；本机工作树那份是 **09-27 的陈旧件**（`274-a2` §4.3 已用哈希名证伪） | 历史件 |

**历史真跑读数（本腿不重复跑，只指名引用）**：
- `274-a2` §4.1/§4.2（10-07，仓外 `git archive HEAD frontend` 导出）：`npm ci` rc=0/10s → `tsc -b` rc=0/7s → `npm run build` rc=0/9s，**26 秒、4 枚／602,635 字节、`index.html` 1,068**（⇒ 那一发是 **CRLF** 那套字节，见票「收 274-r1」§2）。
- `274-v2` `logs/10-real-npm-build.txt`（10-07 16:47，仓外 clone 里跑**甲之后的 build.ps1**）：`npm run build` 真跑（vite 8.3.0，`built in 436ms`）→ `frontend ok: 4 file(s) … index.html is 1068 bytes … .gitkeep=0 assets\index-vdBrT8rM.js=552027 assets\index-yy8KMgdf.css=49540 index.html=1068`，`REAL_BUILD_PROCESS_EXITCODE=0`。
- **CI 侧非本机读数（已有，且今天仍有效）**：票「收 274-v1」§4 从 run `37545246395` 抽到 vite 产出 **3 枚**（`index-CZ-rxcIB.js` 551.96 kB／`index-yy8KMgdf.css` 49.54 kB／`index.html` 1.04 kB）。
  ★本腿补一枚**让这格能用**的限定：`git log -1 --format='%h %ad' --date=iso -- frontend` → **`611ae8b9 2026-09-26 21:26:21 +0800`** ⇒ **页面源码自 09-26 起一字未改**，而甲/乙/丙三形都不动 `frontend/**` ⇒ **那发 CI 的 3 枚读数对"当前页面源码"仍然成立**；
  且本腿现量 `awk 'NR>=922' .github/workflows/ci.yml | grep -c 'upload-artifact'` = **0** ⇒ 那 3 枚落在 runner 工作区 `frontend/dist`、**job 结束即弃**，"落在哪"这一半今天只有"落在被丢弃的工作区里"这一形。
⇒ **仍欠的那一枚**：`AC#0(a)` 若要的是"**甲落地后** CI 那一发的产物枚数/名册"，**只有推送后的 run 才有非本机读数**（本腿不跑 npm ⇒ 本机也没有今天那一发）。

### 2.2 (b) `build/wisp.exe panel-assets` 回什么——两格分开写

**只读留证（跑 exe 之前先量，见 `00-anchor.md` §3）**：`ls -la build/wisp.exe` → **`31074357 B，mtime Oct 7 11:57`**；`stat -c` → `2026-10-07 11:57:46.007254300 +0800`。
来历：票「收 274-r1」§4-1 具名记着"那份 exe 是 `274-r1` 重跑覆盖的，起手那份另存 `/d/tmp/wisp274r1/wisp.exe.pre-274r1-backup`，撤销口令＝「还原 274 前 exe」"⇒ **它不是某个还原点，而是 10-07 11:57 那一趟构建的产物**（比 `35633445` 那笔 commit 早 5 分钟）。本腿⛔ 没动它。

**格一＝"通道通"（本腿现跑，rc 全 0）**：

```
$ ./build/wisp.exe panel-assets
panel assets embedded: 4 files, entry=index.html built=true                     rc=0
$ ./build/wisp.exe panel-assets -manifest
.gitkeep 0 e3b0c44298fc1c14
assets/index-B8yINMF1.js 551989 8f06145dac449fb6
assets/index-yy8KMgdf.css 49540 d0b198664b250973
index.html 1044 9b7856b949d63989                                                rc=0
$ ./build/wisp.exe panel-assets -check
panel assets check: entry=index.html built=true 2 asset refs resolve [./assets/index-B8yINMF1.js ./assets/index-yy8KMgdf.css]   rc=0
```
（三发各带一行 winsec INFO 抬头，非错误；尺＝逐字命令行＋`echo rc=$?`。）
⇒ **这一格只证明：embed 确实把 4 枚字节带进了那枚 exe，且入口自洽（refs 全 resolve）。**

**格二＝"带的是这份源码"——本腿判不了，⛔ 不许当反证**：
- exe 内嵌的 JS 名＝`index-B8yINMF1.js`；**盘上工作树的 dist 是 `index-BVKlegVD.js`；仓外 fresh（HEAD 源码、CRLF 那套）是 `index-vdBrT8rM.js`；CI（Linux）是 `index-CZ-rxcIB.js`** ⇒ 四个名字里没有一个两两相同，**exe 与今天盘上 dist 已经"同上下文对"失配**（票「收 274-v1」§3-3 已把它降级为**〔已成立・不可复核〕**，本腿复认那一对名字今天确实不在盘上）。
- 要证"带的是这份源码"需要**一次 fresh build 的名册＋exe 逐枚相同**（票 AC#9 的尺），而本腿被禁跑 `npm`／`go build`／`build.ps1` ⇒ **这一枚欠的是一次真跑**；其中非本机那一半**只有推送后的 CI 那发 run 才有读数**（本机能补的那一半归写腿，不属本程）。
- ⚠ 顺手一枚量级对照（⛔ 不是受控判据）：self-hosted runner 今天 08:14 由 CI 旧码产的 `wisp.exe`＝**30,471,221 B**（那一发 dist 只有锚）vs 本机这份 **31,074,357 B**（带页面）⇒ 差 **603,136 B**，与 dist 三枚产物合计 **602,635 B** 只差 501 B。两份 exe 的**源码状态不同**（`git rev-list --count origin/dev..dev` = **298**），⇒ 这只能当量级，⛔ 不许写成"exe 大小＝页面大小"的判据。

---

## §3 `AC#2` 三形的现量代价（甲／乙／丙）

**共同前提（本腿现量的今天值，⛔ 别再用旧行号）**：
`.github/workflows/ci.yml` ＝ **992 行 / 61,456 B**（mtime 10-08 09:28，最后一次动它的是 `f8810238`，注释性、⛔ 未 push）；
六枚 job＝`grep -n '^  [a-zA-Z0-9_-]*:$'` 命中 10 行里剔除 `on:` 块 4 行 → `lint:65 test-core:395 test-windows:505 slo-smoke:790 slo-full:848 lint-frontend:922`；
`runs-on` 六枚＝`:66 :396 :506(=windows-latest) :791(=windows-latest) :849(=[self-hosted, wisp-slo]) :923`；
`scripts/build.ps1`＝**313 行 / 18,140 B / md5 `43fb7e20c597192104728befb231e8ea`**（尺＝`wc -l -c`＋`md5sum`）；
`git remote`＝`origin`(GitHub)＋`cnb`(cnb.cool)，上游分支 `cnb/dev`；`git rev-list --count origin/dev..dev`＝**298**，`git merge-base --is-ancestor 35633445 origin/dev` → **NOT_on_origin**（⇒ 甲/乙/丙写进 CI 之后的颜色**今天一律无 run 侧凭据**）。

### 甲＝`build.ps1` 第 2 步真跑 `npm ci && npm run build`　**【已在盘上】**

| 栏 | 现量 |
|---|---|
| 动哪几枚文件·枚数 | **1 枚**＝`scripts/build.ps1`。尺＝`git show --name-only --format='' <c> \| grep -v '^\.scratch'` 对 `35633445`／`91c90aa1`／`674c1920` 逐笔跑 ⇒ 三笔都只回这一枚；`.github/workflows/ci.yml` **一笔没动** |
| 增量行区间 | numstat：`35633445`＝**95/2**（真跑＋失败即红）→ `91c90aa1`＝**41/1**（出处快照对拉）→ `674c1920`＝**14/6**（数组对拉修 AC#11）；累计 172→**313** 行 |
| 今天第 2 步的射程 | `:86`（`# --- 2. frontend ---`）→ `:215`（`frontend ok` 那行打印）；`:34-38` param 仍**只有 `-Env dev\|prod`**；`Fail` 定义 `:47-50`，唯一非零出口 `:49 exit 1`（`:313 exit 0`）；`Fail` 调用点全文件 **17 枚**，其中第 2 步内 **10 枚**＝`:101 :119 :123 :136 :166 :178 :186 :189 :193 :212`（尺＝`grep -n 'Fail '`＝17 枚 vs `grep -c 'Fail'`＝18 枚（含函数定义），⚠ 两种尺写法差 1） |
| `go build` 落点 | 现量 `:270`（`274-a2` §1.1 给的 `:129` **已过期**，别再引） |
| 新 job/新 artifact | **都不需要**（`274-a1` §4 已给，本腿复认 ci.yml 名册未变） |
| CI 时长 | 见 §1.2-1：组成读数有（Linux 5s+4s；三枚 Windows build 步基线 21/35/28s），**净增量＝〔待推送〕**；npm 缓存在那三枚 job＝**0 枚**（§1.2-2） |
| Node 硬依赖 | 已成真：`:119`（node 取不到⇒红）／`:123`（npm 取不到⇒红），⛔ 无跳过后门；那句"没有 node 就整条 FAIL"已写进 `.DESCRIPTION` `:10-14` |
| 已实现的额外代价（票面没写、盘上有） | **同一份源码连跑两次必红**＝`:212`（编排者在票「收 274-r2」§3 裁"可接受"，`274-v2` 用真 npm 实测成读数：`logs/11-real-second-run.txt` 第 42-43 行 FATAL＋`SECOND_REAL_BUILD_PROCESS_EXITCODE=1`）⇒ 这条以后**任何**在 CI/本机重复起 build.ps1 的通路都会撞上（`slo-check.ps1` 那一支见 §5-4） |
| "以后谁能看见" | 已由"预测"变"现实"：本机＋三枚 CI job＋self-hosted runner 全部要 node；无 node 的干净编译机从"能编出 false 形"变"直接红"（`274-a1` §4 甲栏末行所说的那枚副作用**已经发生**） |
| 撞已钉死的东西 | ⛔ `internal/panel/assets.go` 的 fail-closed 语义——**没碰**（本腿复认 `:56-58/:75-77` 原样）；⛔ `frontend/.gitignore:12-13 dist/*`+`!dist/.gitkeep`——**没碰**；⛔ `frontend/**` 只读口径——甲只动构建胶水；`tools/d22scan` 分母——不扫 `scripts/`（§5-5）；**唯一被甲改掉的一句"谎"＝旧 `:74` 那句 `Write-Host ... frontend step skipped ...`**（票「收 274-a2」§3-ⓑ 与 `274-a2` §3.4-1 引的 CI 日志第 183 行从此不再出现 ⇒ 那两枚旧读数今后不可再当"今天"引） |

### 乙＝CI 跨 job 传产物（`lint-frontend` 上传 dist → Windows job 下载后再 `build.ps1`）

| 栏 | 现量 |
|---|---|
| 动哪几枚文件·枚数 | **1 枚**＝`.github/workflows/ci.yml`（`274-a1` §4 乙栏原话"只有 ci.yml"今天仍成立） |
| 行区间（今天） | 上传步插在 `:954`（`npm run build`）之后；下载步＋`needs:` 插在 `:571`／`:808`／`:871` 三枚 build 命令之前；job 声明行 `:922`（上传方）与 `:505/:790/:848`（三枚下载方） |
| 新 job/新 artifact | 新 job＝0；**新 artifact＝1**；`download-artifact` 全文件 **0 枚**（正控＝同尺 `upload-artifact` 命中 `:814 :877`＋注释 `:893 :894`）⇒ 属新引入依赖，按仓里先例（`:893-894` 那段注释正在 pin SHA 的先例）要 pin |
| job 顺序 | `grep -c 'needs:'`＝**0** ⇒ 六枚 job 今天全并行；乙要新增 **3 条边**（`lint-frontend → test-windows/slo-smoke/slo-full`） |
| CI 时长 | 关键路径多出 `lint-frontend` 整枚（今天 **23 秒**，本腿 §1.2-1 现量）＋三枚 job 各一次下载（仓里**无** dist 体积先例，dist 三枚产物合计 602,635 B 可作量级）；⛔ 上传/下载净秒数＝〔待推送〕 |
| 与已落地的甲的相互作用（⚠ 本腿新点出的一格代价） | 甲已在第 2 步真跑 npm，且 `:212` 那道出处闸门要求"npm 之后名册里出现快照没有的名字"。⇒ 乙那一形"先把 dist 下载进 `frontend/dist`、再跑 build.ps1"若 npm 产出**逐字节相同**的内容，vite 给同名 ⇒ **快照里已有 ⇒ 红**（同形实测＝`v2/logs/11` rc=1）。要让乙在甲之后还成立，只有两形：甲的 npm 被跳过（today⛔ 不存在跳过支，`:86-215` 全段无条件），或 CI 里 npm 产出与下载件**不同名**（跨机器名册本已互不相同：CI `index-CZ-rxcIB.js`／本机 `index-vdBrT8rM.js`／exe 内 `index-B8yINMF1.js`，票「收 274-v1」§4 已具名）。⇒ **这条边的成败只能由推送后的那一发定**，本腿只能给出"形状上会撞哪一行"。 |
| 撞已钉死的东西 | "把 dist 提交入库"这一乙变体被票 §3 第 2 枚 ⛔（`dist/*` 规则不可动）直接禁；`internal/panel/assets.go`／`frontend/**` 不需碰；`tools/d22scan` 不受影响 |
| "以后谁能看见" | 只有 CI；本机形状不变（`274-a1` §4）——⚠ 甲落地后这句要重述：**本机已经变成"必须 npm"那一形**，乙只是把 CI 三枚 job 拉成与本机同形 |

### 丙＝产线不动，只加一枚发布前闸门（强制 `built=true` 且入口字节>0）

| 栏 | 现量 |
|---|---|
| 动哪几枚文件·枚数 | **1 枚**＝`.github/workflows/ci.yml`（挂既有 job）；**或 2 枚**＝再加 `scripts/build.ps1` 尾部（`:307-313` 那一支 smoke 之后），后者把闸门挪进产线，⚠ 与甲同文件 |
| 落点（今天） | 三枚 build 命令 `:571`／`:808`／`:871` **各自之后**各插一步；既有 upload 步在 `:814`／`:877` ⇒ 丙若挂在 slo 线，插入点是 `:813` 与 `:876` 之前。**票面说的"打包/签名之前"今天没有栖息地**：`grep -n '^  [a-zA-Z0-9_-]*:$'` 六枚 job 里**没有 release**，`grep -rn 'release' .github/workflows` 亦 0（`274-a2` §5.1-1 的顶回，本腿用今天的尺复认成立） |
| 新 job/新 artifact | 挂既有 job＝0／0；真做"发布前"＝**新 job**（票「收 274-a2」§5-丙明令⛔ 不许以"顺手加个 job"扩张本票） |
| 现成的牙 | `cmd/wisp/panel_assets.go:126-129`（`!Built()` ⇒ stderr＋**rc=1**）与 `:105-113`（`-manifest` 打印 `name size hash`）——本腿 §2.2 现跑通；⚠ **"入口字节>0"在 exe 侧没有任何现成判据**：`assets.go:56-58` 只 `fs.Stat`，一枚 0 字节的 `index.html` 会照样 `built=true`（唯一断言在 Go 测试 `gate_test:117-119` 与 dist 侧 `build.ps1:191-194`） |
| CI 时长 | 每 job 一条命令，秒级（三形中最小；`274-a1` §4 已给） |
| 会不会碰到已经钉死的东西 | 不改 `assets.go` 语义、不改 `.gitignore`、不开窗 ⇒ 票 §3 五枚 ⛔ 都不挡；⚠ 但丙**单独**落地时 exe 仍不带页面（只把"没页面"从静默变红牌）——`274-a1` §4 丙栏已具名；today 甲已落地 ⇒ 这一句**从"警告"变"无关"**（exe 带页面已是产线默认） |
| "以后谁能看见" | 只有执行闸门的 CI 那一发；本机不自跑（`274-a1` §4） |

### 3.4 ★那句"host needs no Node"在哪份文件哪一行（派单要我现量找，找不到就写"不存在"）

**它存在**，共四处，逐字：

| 处 | 尺 | 原文 |
|---|---|---|
| `.scratch/wisp/issues/34-frontend-scaffold.md:32`（票 274 `:36`/`:49` 引的就是这一枚） | `grep -rn "host needs no Node"` | `- Build: \`docker/frontend.Dockerfile\` (node:20-alpine, npm ci, build, export dist) →` / 下一行 ``  `assets/web/dist`; host needs no Node. CI: tsc+lint+build job.`` |
| `docs/specs/SPEC-11-build-deploy-containerization.md:83`（中文同义句，在 §3.2＝§2.2 第 2 步指过去的那节里） | `grep -n '宿主机无需' docs/specs/SPEC-11-*.md` | `- 产物落 \`assets/web/\` 供 \`//go:embed\`；宿主机无需 Node；版本随 \`package-lock.json\` 钉死。` |
| `frontend/embed.go:4-7`（活码注释；本腿只读该 4 行构建胶水，未读页面源件） | `grep -rni 'no node, no npm'` | `// self-sufficient on a machine with no node, no npm and no network.` |
| `cmd/wisp/panel_assets.go:7`（同一句的第二份拷贝） | 同尺 | `// machine with no node, no npm and no network. That is only provable if` |

⚠ **两处要把话说清的（不算顶回，算收窄）**：
1. 这四句的射程都写着 **"宿主机／运行机"**（跑起来的 exe 不许依赖 node/网络），或 **"前端在容器里构建"**；**没有一句写着"编译产线不许需要 node"**。票 274 `:36`/`:49` 把它引申成"甲形与之正冲"，那是**编排者的引申，不是原句字面**——甲落地后 `embed.go:4-7` 那句仍然成立（运行时自足），因为甲要 node 的是**构建机**。⇒ 拍形时这一栏的真实代价是"每台编译机硬需要 node"（已发生），不是"违反了某句冻结文字"。
2. 派单问的那句"我记的是 `cmd/wisp/panel_host_windows.go` 里 `:441`/`:449`"——**行号过期**（现量尺 `grep -n 'not built\|NotBuilt' cmd/wisp/panel_host_windows.go`，全文 880 行）：`:451` 调 `m.serveNotBuiltNoticeLocked()`；`:460` 函数定义；**`:468` 才是那页告示的 HTML**（逐字 `</head><body><p>panel assets unavailable: the embedded bundle is not built.</p></body></html>`）；`:474` 是错误串 `panel host: embedded bundle not built (only the tracked placeholder is present)`。⇒ 派单/旧件里的 `:441`/`:449` 今天分别落在注释与调用点，**告示正文在 `:468`**。

---

## §4 `AC#8` 归属裁决的料（`SPEC-11 §2.2` 那两步今天的真身；⛔ 本腿不裁）

### 4.1 spec 侧原文（本腿现量的行号；`wc -l docs/specs/SPEC-11-*.md`＝**191**）

```
:44   2. 前端产物：存在则跳过；--with-frontend 时走 docker/frontend.Dockerfile（§3.2）      ← §2.2 第 2 步
:45   3. go build（cgo: CGO_ENABLED=1, CC=mingw32-gcc; embed assets/web 已就位）            ← §2.2 第 3 步
:69   ### 3.2 `frontend.Dockerfile`（前端构建，S5 起需要）
:74   COPY frontend/package*.json ./
:76   COPY frontend/ ./
:80   # 用法：docker buildx build --output type=local,dest=assets/web -f docker/frontend.Dockerfile .
:83   - 产物落 `assets/web/` 供 `//go:embed`；宿主机无需 Node；版本随 `package-lock.json` 钉死。
:84   - S5 之前 frontend 目录可空，构建脚本跳过该步。
:122  | `lint` | … | … 零 emoji 扫描（design/ 与 frontend/） | 阻塞 |
:125  | `frontend` | ubuntu + node 镜像 | tsc/lint/build + 硬编码色值/emoji 扫描 | 阻塞（S5 起） |
```
（尺＝`grep -n 'with-frontend\|assets/web\|frontend\.Dockerfile'`＋`grep -n 'S5\|frontend'`；`274-a1` §3 给的是 `:39-52` 窗口，今天第 2 步在 `:44`、第 3 步在 `:45`，**内容一字未改**。）

### 4.2 盘上真身（本腿现量）

| 那一行字面对应的东西 | 今天 | 尺 |
|---|---|---|
| `--with-frontend` 开关 | **不存在**（param 块 `:34-38` 只有 `-Env`；全文件 `grep -n -i 'withfrontend\|with-frontend' scripts/build.ps1` → **rc=1／0 命中**） | `grep -n -i …`；`sed -n '34,38p'` |
| `docker/frontend.Dockerfile` | **跟踪 0 枚**；`git ls-files docker`＝**11 枚**＝`builder.Dockerfile`/`compose.dev.yml`/`compose.test.yml`/`mockllm.Dockerfile`/`model-mirror/{SHA256SUMS,nginx.conf,fixtures/…6 枚}` | `git ls-files docker`；`git ls-files \| grep -ci 'docker/.*front'`＝**0** |
| `assets/web`（§2.2 第 3 步与 §3.2 的落点） | **跟踪 0 枚**，连目录都不存在（`ls -d assets` → No such file）；⚠ 但根 **`.gitignore:24` 写着 `assets/web/dist/`**——一条为不存在的路径预留的忽略规则 | `git ls-files \| grep -c '^assets/web'`＝**0**；`grep -n 'dist' .gitignore` |
| §2.2 第 2 步"存在则跳过" | ⚠ **今天的盘上不但没实现，还反向**：`:86-215` 那一段**无条件**跑 `npm run build`（唯一条件支＝`:129-140` 的 `node_modules` 存在则跳过 **`npm ci`**），而"dist 已存在且内容相同"那一形在 `:212` 是 **FATAL**。尺＝`awk 'NR>=86 && NR<=215 && /Fail /' \| wc -l`＝**10**；`v2/logs/11` 实测 rc=1 | `sed -n '129,140p;195,213p'` |
| §2.2 第 3 步"embed `assets/web`" | 真实 embed＝`frontend/embed.go:19 //go:embed all:dist` → `internal/panel/assets.go:44 fs.Sub(frontend.Dist(), "dist")` ⇒ 目标目录是 **`frontend/dist`**，⛔ 不是 `assets/web` | `grep -rn 'go:embed' cmd internal frontend`；`sed -n '44,58p' internal/panel/assets.go` |
| 权威层（`docs/PLAN.md`，3,601 行） | `grep -n 'frontend\.Dockerfile\|assets/web\|with-frontend' docs/PLAN.md` → **0 命中**；D41(a) `:2891` 只写 `resources\  ← 前端产物（亦可 embed.FS）`；`:1944` 写"前端构建如何嵌入：React 产物 → `embed.FS` → Go 二进制。**需要定：**"；`:1009` 写资源加载从内存喂 `embed.FS` | 逐条 `grep -n` |
| `docs/specs/README.md` 的优先级 | PLAN.md 定案内容最高、spec 不得与之矛盾（`AGENTS.md §4` 索引行已具名）——本腿复认 PLAN.md 从未承诺过 docker/`assets/web` 那两步 | — |

### 4.3 "该由 `build.ps1` 兑现" vs "承认那两步作废"——两形各自的**盘上后果**（⛔ 不裁决）

**形一＝由 `build.ps1` 兑现 §2.2 那两步**（开关＋docker＋`assets/web`）：
1. 要新增 `docker/frontend.Dockerfile`（现 0 枚）＋一条 docker build 通路＋把产物落到**不存在**的 `assets/web/`；
2. embed 目标因此要从 `frontend/dist` 改到 `assets/web/dist` ⇒ **必改 `frontend/embed.go:19`**（撞票 §3 第 4 枚 ⛔"不碰页面代码/`frontend/**` 只读"）与 **`internal/panel/assets.go:44` 的 `"dist"`**（撞票 §3 第 1 枚 ⛔"不许改 fail-closed 语义"，`:34/:56-58/:75-77` 那一族正是它）；两条 ⛔ 都**不在本票任何腿的授权面里**；
3. §2.2 第 2 步那半句"存在则跳过"若照字面兑现，**正面拆掉 AC#10/AC#11 刚立的牙**（`build.ps1:212` 的 FATAL 存在的唯一理由就是"跳过/报成功而没写"那一形；实测＝`v1` 的 M-iii 落绿 ⇒ 立 AC#10，`v2` 五格射程 ⇒ 勾 AC#10，`r3`/`v3` ⇒ 勾 AC#11）；
4. 与甲（已在盘上）**互斥**：甲无条件跑 npm，开关一形要么把甲包进 `--with-frontend`（则默认形又回到"跳过＝零页面字节"，正是本票立案的那句谎），要么让开关只切 docker 通道（则**两条前端通路并存**，`frontend/dist` 与 `assets/web` 谁进 exe 需新判据）；
5. 根 `.gitignore:24` 那条 `assets/web/dist/` 会被这条通路**激活**（今天它是死规则）。

**形二＝承认 §2.2 那两步已被"`frontend/` ＋ vite ＋ `frontend/dist` ＋ CI `lint-frontend`"取代（作废）**：
1. 要改的文字＝`docs/specs/SPEC-11:44`（第 2 步）、`:45`（第 3 步的 `assets/web`）、`:69-84`（§3.2 整节，含 `:80` 的用法行与 `:83` 的"宿主机无需 Node"）、`:125`（§6 表 `frontend` 行"阻塞（S5 起）"）——**全在 `docs/specs/**` 禁改区，＝人工批准**（`AGENTS.md §1.1`；票 AC#5 明写"要不要改＝人工批准，本票只上报不修"）；
2. **`docs/PLAN.md` 不需动**（0 命中已量，§4.2 末两行）⇒ 作废不触碰 D29/D41 的任何定案文字，PLAN 侧只留 `:1944` 那句"需要定"待定；
3. 作废后盘上剩一枚化石规则＝根 `.gitignore:24`；
4. **今天盘上事实上已是形二**：`build.ps1` 无条件 npm（`:86-215`）、CI 的 vite 在 `lint-frontend`（`:954`）、embed 在 `frontend/dist`（`embed.go:19`）、`docker/frontend.Dockerfile`／`assets/web` 均 0 枚、`--with-frontend` 0 命中；差的就是 spec 文字没跟上；
5. ⚠ 一枚形二**也躲不掉**的残留：`SPEC-11:45` 那句"embed **assets/web** 已就位"与 `build.ps1` 第 2 步"存在则跳过"是**同一节的两个半句**，票 §2①/`274-a1` §3 都已具名"`build.ps1` 旧 `:74` 那句引用指错了节"——那句现已随甲删掉，⇒ 作废形要动的是 spec 自己的两行，而不是脚本里的一句话。

---

## §5 哪枚既有仪器会被这三形打红（名册＋尺）

**在数"exe 字节／名册／出处"的东西，本腿现量只有这些**：

| # | 仪器 | 它到底数什么 | 三形动了会怎样 |
|---|---|---|---|
| 1 | `scripts/build.ps1:292-305`（第 5 步 SHA256SUMS） | 4 枚产物（`wisp.exe`＋3 枚 DLL）的**存在性＋哈希**，`:297` 缺任一枚即 `Fail`。**不数 exe 字节大小** | 三形都不动这段；⚠ 但它保证了 **rc=0 ⇒ exe 必存在**（这条用于第 4 行的推理） |
| 2 | `scripts/build.ps1:176-215`（第 2 步 dist 存在性＋出处闸门，10 枚 `Fail` 之一在 `:212`） | **名册对拉**：跑前 `$preRoster`（`:155-158`，数组）、跑后 `$postRoster`（`:208`，数组），要求出现快照里没有的名字（`:209-212`） | 甲已实现；**乙会撞**（§3 乙栏末行：artifact 先落 dist ⇒ 同名 ⇒ `:212` 红，实测同形 rc=1 见 `v2/logs/11:42-43`）；丙不撞 |
| 3 | `cmd/wisp/panel_host_gate_test.go:77-155`（`TestPanelBundleShapeSeparatesAnchorFromRealPage_AC12`） | 三支：① embed 侧 `built/entry-bytes/refs/manifest`（`:113-138`）② git 名册 `ls-files frontend/dist`（`:90`）与 `status --porcelain --ignored -- frontend/dist`（`:91`）③ 出处轴 `:141-149` | **甲把 CI 里它的"形"从 `anchor-only / built=false` 翻成 `page-bundle / built=true`**（旧形凭据＝`274-a2` §3.4-2 抽的 `test-windows` 日志第 2908 行；甲上 CI 后**那行不再会复现**）；两支都放行（票 §2④），⇒ **颜色不变、读数换形**；⚠ `:91` 那把尺把被忽略的**目录**折成一行（本腿现量：`git status --porcelain --ignored -- frontend/dist`＝**2 行**，而 `find … -type f`＝**4 枚**），所以 `ignored-or-untracked` 天生不是文件枚数 |
| 4 | `scripts/slo-check.ps1:294-316` | exe 缺失时**第二次**起 `scripts/build.ps1`（`:296` 那句 Write-Host＋`:305-316`），路径 `$WispExe` 默认＝`build\wisp.exe`（`:86`），与 build.ps1 的 outDir（`:259-270`）**同一枚** | 本腿补一枚**盘上收口**（v2 §推翻-2 把这条记为"仍是推断"）：`build.ps1` rc=0 必已过 `:295-297` 的四枚存在性 ⇒ **同一次运行里 exe 必在** ⇒ `:294` 那一支不可达；而 CI 两枚调用它的步（`:811` smoke／`:874` full）**都没有 `if:`**（现量 `grep -n 'if: ' ci.yml` 命中止于 `:784`，`:790-880` 区间 0 命中）⇒ build 步红则 gate 步不跑。**结论：CI 里"第二发 build"今天不可达；本机手工顺序仍未验**（⛔ 不许当凭据） |
| 5 | `tools/d22scan`（`scripts/d22scan.sh` 54 行驱动） | 八禁令＋零 emoji；emoji 段作用域现量 `main.go:581-584`＝**`design/` `frontend/`(everyFile) `internal/`(goOnly) `cmd/`(goOnly)** ⇒ **不扫 `scripts/`、不扫 `.github/`**；`:1136-1143` 明确把 **gitignored 路径剔出分母**（注释自陈 A207 那枚"跑了 build 的机器读 43、CI 读 40"的漂移就是这么补的） | 甲往 `frontend/dist` 写页面字节 ⇒ **被忽略、不进分母**（今天的形状已修）；动 `build.ps1`/`ci.yml` ⇒ **不在作用域**，⛔ 不会挪 d22scan 的任何数 |
| 6 | `scripts/check-path-length-budget.sh`（584 行；CI 落点 `ci.yml:136`＋命令 `:166`；`HAT_NAME_LIMIT=100` `:149`／`DEBT_CAP=180` `:164`） | 数**跟踪路径**的长度名册 | 三形都不新增跟踪路径（dist 产物被 ignore、artifact 在仓外）⇒ 不动；本腿两件长度现量 **40 / 36 字符**（尺＝`printf … \| wc -c`）⇒ 远低于帽 |
| 7 | `scripts/slo-freshness.sh` | 只看 `slo-full-report` **artifact 存不存在**（`:25-45` 注释自陈"artifact exists == report existed == numbers were produced"） | 三形都不改 upload 形状 ⇒ 不动；⚠ 但甲若让 `slo-full` 的 build 步在 runner 上转红，`Upload` 步不达 ⇒ 那枚 freshness 钟**停走**（间接红，值得记进拍形风险栏） |
| 8 | `scripts/wisp-cli-tests.sh:76` | 只在注释里指认"缓存链＋build.ps1"，**不执行** build.ps1 | 不动 |

**"同一枚数、两种尺写法会差 1"名册（本腿现量，逐字尺在括号里）**：
- `needs:` **0** ／ `needs` **5**（`grep -c 'needs:'` vs `grep -c 'needs'`，后者命中 `:433 :541 :573 :612 :688` 的注释与步名）
- build.ps1 在 ci.yml 的提及 **10** ／ 真 `run:` 调用 **3**（`grep -c 'build\.ps1'` vs `grep -c 'run: powershell.*build\.ps1'` ⇒ `:571 :808 :871`）
- `npm` 行 **11** ／ `npm run build` **1**（`grep -c 'npm'` vs `grep -c 'npm run build'` ⇒ `:954`）
- job 数 **6** ／ `^  name:$` 命中 **10**（含 `on:` 块 `:11 :12 :36 :38`）
- `Fail` 调用 **17** ／ `Fail` 字样 **18**（含 `:47` 函数定义）
- dist 枚数 **1（跟踪）／4（磁盘文件）／2（porcelain --ignored，目录折叠）**
- `upload-artifact` 真步 **2**（`:814 :877`）／ 字样 **4**（＋`:893 :894` 注释）

---

## §6 死格名册（本程三资源全禁：⛔ 真窗／⛔ 推送／⛔ 合并页面分支）

| 格 | 按字面能不能跑 | 它具体欠哪一枚读数 |
|---|---|---|
| `AC#0(b)` 后半句"同一发 run 里 `build.ps1` 造出的 exe 回 `built=true/false`、枚数差几枚" | **按字面跑不了**（判据自带"同一发 run"） | 推送后那一发的：①该 job `Build wisp.exe`/`cgo build smoke` 步日志里的 `frontend ok: N file(s)` ②同发 exe 的 `panel-assets -manifest` 逐枚名册 ③与同发 dist 名册的逐枚差。今天最新 ci run＝`37703959747 @ cc315261`，`35633445` **不在 origin/dev**（`git merge-base --is-ancestor` → NOT_on_origin） |
| `AC#0(a)` 的"今天这一发"本机枚数 | **本腿取不到**（⛔ 禁跑 npm） | 一次 `npm run build` 的名册＋字节；⚠ 已有替身读数（§2.1 表末两行：CI 3 枚＋两发仓外真跑），**缺的是"当前工作树、当前源码"那一发** |
| `AC#2` 甲栏"CI 时长代价要现量"（票面字面） | **按字面跑不了**（要推送） | Windows 三枚 job 加 npm 后的**步级时长增量**（组成侧本腿已给：Linux `npm ci` 5s／build 4s；三枚基线 35/21/28s） |
| `AC#2` 甲栏"hosted Windows 作业 PATH 有没有 node" | **取不到**（要推送或在 CI 现加一步去问；`274-a2` §3.1 已判"零枚步骤试过"） | 一发 `node --version` 的**作业侧**输出；机侧 v24.9.0／`D:\work\server\node14`（目录名与版本名实不符）与 runner `externals/node20+node24` 都**不可外推**（`ci.yml:855-858` 注释自己具名过"交互式 PATH ≠ 作业 PATH"） |
| `AC#2` 乙栏"artifact 上传/下载时长＋`needs` 边后果" | **按字面跑不了** | 上传/下载秒数、三枚 job 起跑延迟、以及 §3 乙栏那枚"artifact 先落 dist 再 npm 会不会红"的 **CI 实测** |
| `AC#2` 丙栏"发布前"那一半 | **没有栖息地**（不是要资源，是盘上没有） | `grep` 现量 ci.yml 6 枚 job、无 `release`；票「收 274-a2」§5-丙已把它改成"要么新造 job（超射程）要么认 build.ps1 的非零退出本身是牙" |
| `AC#8` | **不欠读数**（本腿 §4 已把盘上事实交齐） | 欠的是**裁决本身**；若裁成"改 `SPEC-11:44/:45/:69-84/:125`"，那一步落进 `docs/specs/**` 禁改区＝**人工批准**，⛔ 任何腿不许顺手补开关 |
| 票 34 的 `AC#1`（真开一扇窗） | 与 274 无关的**另一枚死格**（票 §4 已具名"它是一枚现在无法执行的判据"） | 真窗读数，⛔ 不在本票射程（票 §3 第 3 枚 ⛔） |

★**顺手替机主解掉一枚假依赖**：派单写"要让 exe 真带页面字节，得先有…才能拍形"，语境里隐含"页面分支可能挡路"。盘上事实：**本票三格没有任何一枚需要合并 `D:/wt/fe` 那支**——
`git log -1 -- frontend`＝**`611ae8b9 2026-09-26 21:26`**（dev 的页面源码自 09-26 起未变、且**已跟踪**），`git ls-files frontend/src \| wc -l`＝**64**（`274-a2` §0 与票 §2② 口径），且 dev 这一支**今天可构建**（`274-a2` §4：`npm ci`/`tsc -b`/`npm run build` 三档 rc=0，26 秒，4 枚）。
⇒ "先不合"不挡这张表；**甲形也已经不需要它**（甲已在盘上且自带出处闸门）。

---

## §7 本腿自报（用了什么、没用什么）

- **用了**：`date`／`git log`／`git rev-parse`／`git status --porcelain`／`git ls-files`／`git check-ignore`／`git show --numstat --name-only`／`git merge-base --is-ancestor`／`git rev-list`／`ls`／`stat`／`find`／`sed`／`awk`／`grep`／`wc`／`md5sum`／`sha256sum`／`printf … \| wc -c`／`gh run list --json`（1 次）／`gh run view --json`（1 次，只抽步级时间戳，⛔ 未拉全量日志、未落盘大文件）。
- **⛔ 一次都没跑**：`go build`／`go vet`／`go test`／`npm`／`npx`／`vite`／`powershell scripts/build.ps1`。
- **⛔ 本腿没跑 `go env`／`go list`**（派单允许但要自报——本腿选择不跑，因为 `33-v4` 正独占同机 Go/仪器读数，`go list` 也会写构建缓存）。⇒ 代价：`go list -f '{{.EmbedFiles}}' ./frontend` 那一格本腿**未复量**，用的是磁盘 `find`＋票 `:54`/`274-a1` §5/`274-v1` §2 的历史读数，已在 §2.1/§4 具名标注为"历史件"。
- **跑了 exe**：`build/wisp.exe panel-assets`（默认／`-manifest`／`-check` 三发，均 rc=0）；跑前已 `ls -la`＋`stat` 留证，⛔ 未动该文件。
- **只读看了别的工作树**：`E:/work/base/actions-runner/_work/wisp/wisp/{frontend/dist,frontend/node_modules,build}` 三个 `ls`（只列名/字节/mtime，⛔ 没在那棵树里跑任何 `git` 子命令，怕动到在飞 job 的 index）。⛔ 未进 `D:/wt/fe`（本腿一次都没访问）。
- **只读边界**：⛔ 未读 `frontend/src/**`、⛔ 未读 `design/**`、未读 `frontend/dist` 内**内容**（只取文件名/字节数/sha256 摘要），`frontend/embed.go` 只取注释里那 4 行构建胶水（`:4-7`）与 `:19` 那枚 directive。`docs/specs/**`／`docs/PLAN.md` 只读，⛔ 一字未改。
- **写面**：只有 `.scratch/wisp/probes/274/a2/00-anchor.md`＋本件＋票 274 末追加的一节。commit 逐笔带显式 pathspec；⛔ 无 `add -A`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；⛔ 零 push。

---

## §8 终态自证

| 尺 | 读数 |
|---|---|
| `git status --porcelain -- scripts .github docs internal frontend \| wc -l` | 起手 **0**／中途一次瞬时 **1**（未复现，且**不是本腿写面**——本腿全程只写 `.scratch/wisp/probes/274/a2/`）／终态 **0** |
| `git status --porcelain -- cmd internal scripts tools .github docs frontend \| wc -l` | 终态 **0 行** |
| `git status --porcelain -- .scratch/wisp/probes/274 \| wc -l` | 锚 commit 后 **0 行**（本件与票追加各自成笔） |
| `ls -la build/wisp.exe` | 终态仍 **31,074,357 B / Oct 7 11:57**（＝起手锚那行，⛔ 未被覆盖） |
| `find frontend/dist -type f \| wc -l` | 终态 **4 枚**（同起手；⛔ 本腿没动 dist 一个字节） |
| `wc -l -c scripts/build.ps1`＋`md5sum` | 终态 **313 / 18,140 / `43fb7e20c597192104728befb231e8ea`**＝起手值（⛔ 本腿没动它） |
| `wc -l -c .github/workflows/ci.yml` | 终态 **992 / 61,456**＝起手值（⛔ 本腿没动它） |
| 票 274 AC 框 | 起手 `grep -cE '^[[:space:]]*- \[ \]'`＝**3**、`grep -cE '^[[:space:]]*- \[x\]'`＝**9**；本腿追加一节后复量＝**3 / 9**（⛔ 一枚没翻，见票末追加节里的同数自证） |
| 本件 | `wc -l -c` 终值随交件回报给出；0 字节＝没交，非 0 |

rc=0
