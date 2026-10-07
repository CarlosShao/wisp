# 274-r2 — 写腿实施件：AC#10「本趟 npm 的出处校验」（`scripts/build.ps1` 第 2 步）

腿编号 `274-r2`；被审件 `scripts/build.ps1`（起手 265 行／md5 `644c2f6af1fa3fc7f8b3ab38c20753f5` → 终态 305 行／md5 `0b5fbbb5f43850f4e3a3545c4f4eaae2`）。
写面只有两枚：`scripts/build.ps1` 与 `.scratch/wisp/probes/274/r2/**`。⛔ 本腿零 `.go` 改动、零 CI yaml（尺＝`git show --numstat 91c90aa1` 只回 `41 1 scripts/build.ps1`）。
⛔ 一枚 AC 框都没翻；判语全在本件，翻勾归编排者。

## §0 起手锚（现量，逐字在 `logs/00-start-anchor.txt` 与 `logs/00-dist-md5-start.txt`）

- `date` ＝ `Wed Oct  7 15:25:06 CST 2026`；起手 HEAD ＝ `0bb71547`（`dev`）。
- `git status --porcelain -- cmd internal scripts .github docs frontend build` ＝ **0 行**（起手与终态同为 0，`logs/32-porcelain-end.txt`）。
- `go list -f '{{.EmbedFiles}}' ./frontend` 起手＝`[dist/.gitkeep dist/assets/index-BRKj5OIJ.css dist/assets/index-BVKlegVD.js dist/index.html]`；终态同串（`logs/31-embedfiles-end.txt`）。
- 本机 `frontend/dist` 四枚 mtime 现量（起手）：`.gitkeep` 10-07 11:57:42／`assets/index-BRKj5OIJ.css`、`assets/index-BVKlegVD.js`、`index.html` 全 **09-27 10:59:04**（陈旧件，本腿只读）。

## §1 选了哪一形：**名册快照对拉**（mtime 形弃用）

工单 AC#10（票面 `:56`）逐字给的是「本趟 npm 之前先给 dist 拍一枚名册快照，跑完后名册必须与快照不同名（**或** `index.html` 的 mtime 晚于 `frontend/src` 最新 mtime）」。取**名册形**。现量理由：

1. **mtime 形会永久放行本票要打死的那一发**。`frontend/src` 最新 mtime（起手现量）＝ **2026-09-27 11:01**，而 `frontend/dist/index.html` ＝ **2026-09-27 10:59:04**——今天这一形恰好也能红，但那是 **2 分钟的运气**：谁在构建之后碰过一枚源件而已。反过来说，**任何一次真构建都把 dist 的 mtime 推到"晚于全部 src"**，此后桩 npm 一个字节不写也永远绿——正是 AC#10 立案的那句「校存在不校出处」。本机如此，CI 的 `[self-hosted, wisp-slo]`（票面 §收 a2 那条点到它复用工作目录）更如此。
2. **mtime 不是"源件变了"的可信代理**：`frontend/src` 里躺着**生成物** `frontend/src/styles/tokens.generated.css`（起手现量＝它正是 src 里最新的一枚，出处尺＝`internal/panel/tokens_fourway_test.go:53` 与 `frontend/scripts/gen-tokens.mjs:32-33` 把它写进 `src/styles/`）。源树里有一枚会被工具重写的文件 ⇒ 「src 最新 mtime」会自己往前走，把一次合法构建判成红。
3. **checkout/拷贝会重置 mtime**：`git checkout`、`cp -r`、解包都会把 mtime 刷成落盘时刻（`cp -a`/`cp -al` 保留、`cp -r` 不保留）⇒ 同一份 dist 换一种到达方式，mtime 形给两种答案。名册形只看**文件名**，与到达方式无关。
4. **vite 的产物名是内容哈希**（`assets/index-<hash>.js`），本仓三次现量名册互不相同（票面 §4 AC#0(a)：CI `index-CZ-rxcIB.js`／本机仓外 fresh `index-vdBrT8rM.js`／`build/wisp.exe` 内嵌 `index-B8yINMF1.js`，而盘上陈旧件是第四个名字 `index-BVKlegVD.js`）⇒ "本趟真构建过"在盘上的可见痕迹就是**名字变过**。
5. **弃用 mtime 也避开墙钟**：新闸门不需要拿"这一步开始时刻"去比 mtime（那是把时钟当判据，第 3 类形状），比较的是同一次运行内的两个名册。

### §1.1 名册形的代价（我自己认，⛔ 不藏）

**同一份源码连跑两次 build.ps1 会红。** 第二次 npm 真跑了，vite `emptyOutDir` 删旧写新，但**内容哈希名与快照完全相同** ⇒ `freshNames` 为空 ⇒ 报 `Named cause: unprovenanced page bytes`。这不是 bug 而是本格的语义：盘上证据无法区分「npm 重写了逐字节相同的名字」与「npm 一个字节没写」，而 AC#9 已裁死「陈旧产物不许替『带页面』报绿」。红因里逐字写了这一支的处置（把 `frontend/dist` 删掉/移走，让前置快照回到 clean-checkout 那一形，再跑）。
读数＝`logs/09-rerun-ident.txt`（桩 npm 按快照里的同名重写 ⇒ **rc=1**，同一句具名红因），这一发是我自造的**代价发**，不是突变发。
CI 侧不受影响：三处调用（`ci.yml:571`/`:757`/`:820`）都在 fresh checkout 上跑，`dist/*` 被 ignore、跟踪的只有 `.gitkeep` ⇒ 前置快照必然是「1 file(s) [.gitkeep]」，真构建必然换名。⚠ 这一句是**推理不是读数**——推送后的 run 侧取数不在本腿射程（机主 10-07 说暂不推）。

## §2 代码落在哪（`git diff -U0` ＝ **4 枚 hunk／+41 -1**，只动 `scripts/build.ps1`）

| 段落 | 终态行号 | 内容 |
|---|---|---|
| `.DESCRIPTION` 第 2 条 | `:17-24` | 补"存在不等于本趟产出"，并写明为什么用名册不用 mtime（第 1 点的现量口径） |
| 名册快照 | `:146-155` | 本趟 `npm run build` **之前**拍：`$preDistExists`／`$preDistFiles`／`$preRoster`（相对 `frontend/dist` 的路径，含"没有 dist 目录"那一形），并打印 `pre-build dist snapshot = …` |
| **出处闸门** | `:195-206` | `$postRoster` vs `$preNames` ⇒ `$freshNames`；空则 `Fail`（`:203-205`） |
| `frontend ok` 行 | `:207-208` | 保留原前缀，尾部并置「本次新增名」与「前置快照」，绿的时候一眼看出换了名 |

具名红因逐字（`scripts/build.ps1:204`，走既有 `:38-41` 那一枚 `Fail`→`exit 1` 出口，无新出口）：

```
build.ps1: FATAL: frontend step: npm run build reported success but frontend/dist holds no artifact name outside the snapshot taken before this run's npm (Env=dev, this run's repo root <root>; snapshot was <n> file(s) [<roster>], dist now holds [<roster>]). Named cause: unprovenanced page bytes - nothing in frontend/dist was produced by this build, so wisp.exe would embed an earlier build's artifacts. A stubbed or no-op npm that writes nothing lands here on purpose. If this was a real build of unchanged sources, vite rewrote byte-identical content-hashed names and no disk evidence can separate that from npm writing nothing: delete or move frontend/dist (the tracked anchor frontend/dist/.gitkeep regenerates from git) so the pre-run snapshot is the clean-checkout shape, then re-run.
```

它自己报名字（`Named cause: unprovenanced page bytes`）并说清是哪一趟（`Env=$Env` ＋ `this run's repo root $RepoRoot` ＋ `before this run's npm`），写法照 `:180/:183/:187` 那三句的 `… Named cause: …`。
⛔ 不是"检查文件存在"那种恒真句：存在性四支（`:172`/`:180`/`:183`/`:187`）在它之前各自独立，新闸门只比**名字**。

## §3 定向突变（必红）——就是 v1 §2 的 M-iii-2 那一发，本腿打死它

台件形状照抄 v1：**仓外**硬链接镜像树（`/d/tmp/wisp274r2/tree-07-mutation`，`cp -al` 逐枚顶层目录，⛔ 不含 `.git`；`frontend` 走 `cp -a` 真拷贝并排除 `node_modules`；`build/` 建空目录而非链接，免得树里写的 exe 与仓内 `build/wisp.exe` 共享 inode），PATH 里**摘掉 `git.exe`**（v1 同法，`:192` 那条已知副作用本腿不修），桩 `npm.cmd` = `D:/tmp/wisp274r2/fakebin/npm.cmd`（报成功、**一个字节不写**），`frontend/dist` 预先躺着**从仓内真 dist 拷来的陈旧真页面**（`.gitkeep`＋`index-BRKj5OIJ.css`＋`index-BVKlegVD.js`＋`index.html` 1044 字节 ⇒ 与 v1 逐字读数同一枚名册）。

- **rc ＝ 1**（`CASE_07-mutation_RC=1`，`logs/07-mutation.txt`）。
- 打出的具名红因逐字（本发实值，非模板）：

```
build.ps1: FATAL: frontend step: npm run build reported success but frontend/dist holds no artifact name outside the snapshot taken before this run's npm (Env=dev, this run's repo root D:\tmp\wisp274r2\tree-07-mutation; snapshot was 4 file(s) [.gitkeep assets\index-BRKj5OIJ.css assets\index-BVKlegVD.js index.html], dist now holds [.gitkeep assets\index-BRKj5OIJ.css assets\index-BVKlegVD.js index.html]). Named cause: unprovenanced page bytes - …
```

- 同一发里可见 `frontend: pre-build dist snapshot = 4 file(s) […]`（快照真拍了）与桩 npm 的 `reporting success while writing not one byte`。
- v1 那句绿读数（`frontend ok: 4 file(s) … built=true` ＋ `EXITCODE=0`）**今天取不到了**：通路在 `:204` 就 `exit 1`，`go build` 那一发 exe 根本没被造出来。⇒ M-iii-2 落绿这一格被本发打死。

## §4 正控（必绿）：本趟真产出了新名册

`logs/10-positive.txt` ＋ `logs/10-positive-driver.txt`（`CASE_10-positive_RC` 那一行落在 driver 的 stdout 里——powershell＋GUI 子系统 exe 把日志句柄拖到了文件末尾之后，本腿不重跑，逐字引 driver：`--- 10-positive rc=0 ---`）。
形状：同一把尺、同一枚陈旧 dist 起手，桩 npm 这次**删旧写新**（合成 fresh 名册 `index.html`/`.gitkeep`/`assets/index-FRESHr2.js`/`assets/index-FRESHr2.css`）。

- **rc ＝ 0**，整条通路走完（`fetch-deps: cache hit` → 前端 → `go build ok (cgo linked against sherpa-onnx C API)` → DLL/manifest → SHA256SUMS → `wisp doctor: PASS` → `build.ps1: done. …`）。
- `frontend ok:` 那行逐字（一眼看出换名）：

```
build.ps1: frontend ok: 4 file(s) in frontend/dist (entry index.html is 1048 bytes): .gitkeep=0 assets\index-FRESHr2.css=552 assets\index-FRESHr2.js=551 index.html=1048 || provenance: this run's npm added new artifact name(s) [assets\index-FRESHr2.css assets\index-FRESHr2.js] over the pre-build snapshot 4 file(s) [.gitkeep assets\index-BRKj5OIJ.css assets\index-BVKlegVD.js index.html]; gone since snapshot: [assets\index-BRKj5OIJ.css assets\index-BVKlegVD.js]
```

- exe 侧配套读数（`tree-10-positive/build/wisp.exe panel-assets`，⛔ 不是仓内那枚 exe）：`panel assets embedded: 4 files, entry=index.html built=true`／`PANEL_ASSETS_RC=0`；`-manifest` 逐枚＝`.gitkeep 0 e3b0c44298fc1c14`／`assets/index-FRESHr2.css 552 b91e45239550f53a`／`assets/index-FRESHr2.js 551 ea632664defd00d4`／`index.html 1048 765fd902d054107c` ⇒ **内嵌名册＝本趟 dist 名册**（同一次运行内对拉，第 125 条口径；⛔ 不许跨机器/跨上下文当相等）。
- ⚠ 这一发用的是**合成产物**（AC#10 简报明许可"仓外合成的 fresh 文件"）。它证明的是"闸门对'名册换名'放行、通路走完"，⛔ **不证明**真 vite 构建在本机的换名行为——真 vite 名册读数仍归 AC#0(b)/票 34 那格，本腿没跑真 `npm run build`（树里没 `node_modules`，`npm ci` 会把 Node 依赖再拉一遍，不在本格射程）。

## §5 既有支没被打瞎（各一发，全部 rc=1，各自**不同**的具名红因）

| 发 | 形状 | 命中支（终态行号，起手号在括号里） | 具名红因（尾句逐字） | rc |
|---|---|---|---|---|
| `01-anchor` | dist 只剩 `.gitkeep`，桩 npm 不写 | `:180`（原 `:157`） | `Named cause: build produced no page artifacts.` | 1 |
| `02-noentry` | dist 有 1 枚产物但无 `index.html` | `:183`（原 `:160`） | `Named cause: missing entry file.` | 1 |
| `03-emptyentry` | `index.html` 0 字节 | `:187`（原 `:164`） | `Named cause: empty entry file.` | 1 |
| `04-nodist` | 整个 `frontend/dist` 不存在 | `:172`（原 `:149`） | `Named cause: dist directory missing, so wisp.exe would embed no page.` | 1 |
| `05-failbuild` | 桩 `npm run build` 退 1 | `:160`（原 `:137`） | `'npm run build' failed with exit code 1 in …` | 1 |
| `06-nonpm` | node 找得到、npm 找不到（PATH 里只放 `node.exe`） | `:123`（原 `:114`） | `npm could not be resolved on PATH (tried npm.cmd, then npm). …` | 1 |
| `09-rerun-ident` | 桩 npm 按快照同名重写（§1.1 代价发） | **`:204` 新支** | `Named cause: unprovenanced page bytes` | 1 |

七支各自报**不同的**具名原因，`Fail` 仍只有 `exit 1` 一条出口 ⇒ 新闸门没有把谁吞成恒真或恒红。
⚠ 一枚台件副产物具名登记：`logs/08-positive.txt` 是我第一次正控——桩 npm 自己的 `-f` 格式化炸了（`Error formatting a string`）⇒ 它退 1 ⇒ build.ps1 正确报 `:160` 那一支。那是**我的台件缺陷**（已修在 `fakebin/fake-npm.ps1`），不是被审件的病；这一发顺手多证了一次"`npm run build` 非零必须红"，`08` 保留不删。

## §6 门禁四数（本腿自跑，⛔ 未引 v1 的 logs）

1. `sh scripts/d22scan.sh` ＝ **rc=0**（`D22SCAN_RC=0`）。**正控先绿**才认它：`logs/20-d22scan-end.txt:239` ＝ `runtests.sh: OK - packages=[./...] top-level: PASS=35 FAIL=0 SKIP=0, === RUN=77, '[no tests to run]'=0`，其后才是 `d22scan: clean - no D22 ban violations`。
2. `go vet ./cmd/wisp/ ./internal/panel/` ＝ **rc=0**（`GOVET_RC=0`）。
3. `gofmt -l cmd/wisp internal/panel` ＝ `cmd\wisp\models.go`；`/d/work/base/gopath/bin/gofumpt.exe`（`v0.12.0 (go1.27.1)`）`-l` ＝ **同一枚**。⚠ 红因＝**盘上 CRLF**（`cmd/wisp/models.go` `tr -cd '\r' | wc -c` ＝ 334；`internal/panel/assets.go` 同尺＝0），本机工作树 13 枚跟踪 `.go` 带 CR（`core.autocrlf=true` 的幻影）⇒ **不是未格式化，本腿一个字都没"格式化"**，与 274 无关（v1 §3-5 同判）。
4. 整包 `go test -count=1 ./cmd/wisp ./internal/...`（`GOCACHE=D:/tmp/gocache-274r2`，两发同一缓存）：
   - 起手（15:25，`logs/gotest-start.txt`）＝ **rc=1、6 枚红**：`TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestC21DesignTokensFourWayAgree`／`TestC21TableColourRowsMatchTokensCSS`／`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestResolvePerCallBudget`。
   - 终态（`logs/gotest-end.txt`）＝ **rc=1、5 枚红**（前 5 枚同名）。
   - ⚠ **逐名作差 ≠ 0**：差 1 枚，方向＝**红→绿**，掉的是 `internal/risk` 的 `TestResolvePerCallBudget`（起手 2.67s 红、终态无该行＝绿）。本腿**零 `.go` 改动**（`git show --numstat 91c90aa1` 只回 `scripts/build.ps1`），`internal/risk` 起手到终态 porcelain 恒为 0 行 ⇒ 归因不在本腿。我补了一把尺：`go test -count=2 -run TestResolvePerCallBudget ./internal/risk` ＝ **ok**（两次都过）；出处 `internal/risk/pathresolver_budget_norace_test.go:21` 是一枚**每-call 预算／时序敏感**测试，起手那一发我同时在跑镜像树构建与整包测试 ⇒ **最可能的红因是负载抖动**。⇒ 如实写成「差 1 枚，红→绿，本腿未取到"起手那一发为什么会红"的直接证据」，⛔ 不许按"应该相等"抹平，也不许为它动任何断言。

## §7 还原自证（仓内 `frontend/dist` 只读，一个字节没动）

- 起手 md5（`logs/00-dist-md5-start.txt`）与终态 md5（`logs/30-dist-md5-end.txt`）`diff` ＝ **空**（`DIST_MD5_START_EQ_END(rc=0)`），四枚同名同串：`.gitkeep d41d8cd9…`／`index-BRKj5OIJ.css db4db7a2…`／`index-BVKlegVD.js 70128a3d…`／`index.html f98bfc4b…`。中途也回查过一次（`DIST_MD5_MID_IDENTICAL`）。
- `go list -f '{{.EmbedFiles}}' ./frontend` 起手＝终态＝4 枚同名（`logs/31-embedfiles-end.txt`）。
- `git status --porcelain -- cmd internal scripts .github docs frontend build` 起手＝终态＝**0 行**（终态那份在 `logs/32-porcelain-end.txt`；`scripts/build.ps1` 已入库故不显示）。
- ⛔ 任何 `rm -rf`/`Remove-Item` 只发生在我自己的仓外树的 `frontend/dist`（`/d/tmp/wisp274r2/tree-*/frontend/dist`）与桩 npm 的合成写入；仓内 dist／仓内 `build/wisp.exe`／`frontend/node_modules` 一个字节没碰。镜像树与桩件**只建不删**，全部留在 `/d/tmp/wisp274r2/`（`tree-01`…`tree-10`、`fakebin/`、`logs/`、两份 gotest 原始件 12KB/11KB 已入库）。
- 树的 `third_party/` 是硬链接——本腿先跑了 `07` 一发确认 `fetch-deps` 走的是 **cache hit**（只读校验，不写 `third_party`）才继续开其余树；这条风险已具名消掉。

## §8 没做完／射程外（逐枚具名，⛔ 不写成通过）

1. **真 `npm run build` 在本闸门下的换名行为**＝未取数（本腿用合成 fresh 名册，§4 已具名）。它要 `npm ci` 拉一遍依赖，且落点属 AC#0(b)／票 34。
2. **CI 那一形（fresh checkout ⇒ 前置快照＝`.gitkeep` 一枚）必然绿**＝**推理不是读数**；run 侧取数要推送，机主 10-07 说暂不推 ⇒ 本腿没推、没取。
3. **AC#6 第 4 数的红名册作差＝1 枚（红→绿）**，直接证据未取到（§6-4）。
4. `:192` 的 `git rev-parse` 在无 `.git` 目录里会被 `$ErrorActionPreference='Stop'` 升成终止错误——⛔ 不属本票，本腿按简报未顺手修。
5. 本腿**没跑** `npm ci` 那一支（`:140`，原 `:127`）：树里放了 `node_modules` 标记目录让它跳过。它不是简报点名的五支之一，但仍是"闸门未覆盖的一支"，如实登记。

## §9 与本腿简报不一致之处（⛔ 以盘上原文为准，逐条具名）

1. 简报给的**行号全部对得上盘上起手那份 265 行／md5 `644c2f6a…`**（`:104-168` 是前端段、五支在 `:114/:137/:149/:157/:160/:164`，`Fail` 定义 `:38-41` 且只有这一枚 `exit 1` 出口——我逐条 `git show 91c90aa1^:scripts/build.ps1 | grep -n 'Fail '` 复认）。⚠ 唯一要提请注意的是**本腿落码后这些行号一律右移**：`Fail` 现在在 `:123`（npm 取不到）／`:160`（`npm run build` 非零）／`:172`（无 dist）／`:180`（只剩锚文件）／`:183`（无 `index.html`）／`:187`（入口 0 字节）／`:204`（**新支·unprovenanced**）。后续任何引用（含编排者翻勾时的判语）请认终态行号，⛔ 别拿简报那组去对盘。
2. 简报让我只读 v1 verdict 的 §2/§3——照办（没整读，也**没自量**那份的字节数，简报说的 29KB 我按票面 `§收 v1 §1` 的 29,240 字节引用，⛔ 不算本腿读数）。本腿 §3 那发打死的就是它 §2 表里 M-iii-2 那一格，形状逐条复刻（同一枚陈旧 dist 名册、同样摘掉 `git.exe` 的 PATH、仓外镜像树）。
3. 简报「五支既有支」里 `npm 取不到（`:114`）`——本腿实测那一支在起手那份是 `:114` ✓，但**同一支在 v1 §2 表里被记为 M-iv**（腿自造），本腿复跑的是同一形状（PATH 只给 `node.exe`），红因逐字与 v1 那行一致 ⇒ 无冲突，只是"哪一发出自谁"要分清。
4. 简报说 `md5sum frontend/dist/.gitkeep frontend/dist/assets/* frontend/dist/index.html` 起手能存四枚——盘上确实四枚；但**`.gitkeep` 的 mtime 是 10-07 11:57（不是 09-27）**，即 r1 那次 `mv` 出/回留下的痕迹。它不影响名册形（名册不比时间），但它是"mtime 形会被搬运痕迹干扰"的第 4 枚现量旁证，已写进 §1。
5. 简报把 `d22scan.sh` 的第 1 步正控写作 `PASS=35 FAIL=0 SKIP=0`——盘上本腿这一发**逐字同串**（`logs/20-d22scan-end.txt:239`）✓。
6. 其余（写面范围、commit-first、不 push、pathspec、证据件不许叫 `.out`）均照办；无冲突。
