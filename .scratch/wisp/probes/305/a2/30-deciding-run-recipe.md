# 305-a2 · 30 表③：判死那一发的配方（跑的人＝编排者；本腿一发⛔ 跑过）

⚠ 本腿车道＝⛔ 任何 Go 编译面／⛔ 真窗 ⇒ 下面每一发都是**给编排者的死配方**，⛔ 读数。凡"跑到哪一格才算判死"我按格子写清，⛔ 把读码推的那层写成已判死。

## 0. 台面口径（三根轴，⛔ 任何一根漂了那两发就⛔ 可比）

| 轴 | 现量尺（我这发跑的，rc 全在 `logs/30-recipe-rulers.txt` rc=0） | 后果与硬要求 |
|---|---|---|
| ① 母仓 ↔ 干净 clone | `wc -c frontend/dist/index.html`＝**1044**；`git cat-file -s HEAD:frontend/dist/index.html`＝**rc=128**，stderr 逐字 `fatal: path 'frontend/dist/index.html' exists on disk, but not in 'HEAD'`；`git ls-tree -r --name-only 416d9d56 -- frontend` 与 `… f718e9b6 …` 都**只**含 `frontend/dist/.gitkeep` | AC13 在干净 clone 里 `--- SKIP`（具名理由逐字 `the embed resolves no entry, so there is no page content that could be covered`，出自 `panel_resident_windows_test.go` 的 `t.Skipf`）⇒ **AC13 的颜色只有母仓台面，或"clone＋把钉住的那 1044 字节放进 `frontend/dist/`"才量得到**（配方 §2 用后者，因为⛔ 在共享工作树里 detach） |
| ② 同一枚母仓里那枚未跟踪 dist 自己会漂 | `sha256sum frontend/dist/index.html`＝**`9b7856b949d639898832bbdc9a46734a17c57d8c0c4c1611fe10db45de07074f`**；mtime＝**`2026-10-10 08:51:10.506316400 +0800`；文件内含 `<title>Wisp</title>` 与唯一一枚 `id="root"`（⇒ `entryIDProbes` 只交回 1 枚，与红句"NONE of the 1 element ids"对上） | ★凡引 AC13 那一发，必须**同时**交 `wc -c`＋mtime＋sha256 三枚；⛔ 只交字节数（字节同、内容可不同）。10-01 那批 240-0-0 的绿读数用的是**另一枚字节**（R3＝已具名"⛔ 可复现"） |
| ③ ★新轴（本腿量到的，会咬人） | `git ls-files third_party/sherpa-onnx`＝**0 行**、`git ls-files build`＝**0 行**；盘上 `third_party/sherpa-onnx/{onnxruntime.dll,sherpa-onnx-c-api.dll,sherpa-onnx-cxx-api.dll}` 与 `build/*.dll` **只存在于母仓**（`build/` 里另有 33p1-*.exe 等别家产物） | ⇒ 干净 clone 里 `$PWD/third_party/sherpa-onnx` 与 `$PWD/build` 两枚路径**⛔ 存在**；照抄 `PATH="$PWD/…"` 会 `0xc0000135` 且**零 `--- FAIL`**（＝用例根本没跑，正是本项目吃过的假绿形）。⇒ **每一发的 PATH 两枚条目都写母仓绝对路径的 shell 形 `/d/work/workspace/projects plans/Wisp/third_party/sherpa-onnx` 与 `/d/…/Wisp/build`**（⛔ `D:/…` 那形），并在 run 0 里断言"具名用例真发出了 `--- PASS`/`--- FAIL`/`--- SKIP` 行" |

其他台面口径：窗口**⛔ 有 merge**（`git rev-list --count --merges 416d9d56..f718e9b6`＝**0**；`--first-parent --count`＝**963**＝全窗口枚数 ⇒ 历史线性，逐枚 detach 安全）。⛔ 仓内建 clone/worktree/checkout；clone 一律建在仓外（下面示例路径 `D:/work/workspace/bisect/305/<sha>`）。开跑前 `tasklist` 现量 `wisp.exe`／`balldebug.exe`＝0。

## 1. 二分名册（**9 个点**＝7 枚候选＋两端点；⛔ 963 枚全列）

尺＝`git log --format='%h %cd' --date=short --reverse 416d9d56..f718e9b6 -- cmd/wisp/panel_host_windows.go cmd/wisp/panel_resident_windows.go cmd/wisp/panel_resident_windows_test.go cmd/wisp/panel_host_windows_test.go cmd/wisp/panel_host_windows_live_test.go internal/panel/bridge.go internal/panel/composer.go internal/panel/composer_dispatch.go cmd/wisp/panel_inbound.go internal/panel/pump.go cmd/wisp/panel_pump_test.go`（我这发 rc=0，输出 7 行，件 `logs/30-recipe-rulers.txt`；射程＝这 11 枚路径，即 20 号件 §4 的文档链＋入向门＋时长面）。

| 序 | sha | 日期 | 为什么必须进序列 |
|---|---|---|---|
| G | `416d9d56` | 10-01 | **绿端点**（`docs/evidence/s1/33-panel-host-c27-r7.md` 的 AC13/nail2 两枚 PASS 都记在这个 HEAD 上） |
| 1 | `f7d28ef0` | 10-01 | ★能——`bringUp` 里新增 `staleCloseQueued()` 拒绝出口（建窗前 return ⇒ 步 5/6 那两次 `SetHtml` 永远⛔ 跑） |
| 2 | `7a0236b3` | 10-01 | ★能——认哨兵 `errPanelRefusedThread` ⇒ `showOnThread` **判死常驻面板线程**；并在 AC13 与 nail2 之间插进 179 行会开真窗的新用例 |
| 3 | `0d87a681` | 10-02 | 可能不能——入向门白名单扩容（`ip/bridge.go` 新增 `MethodConfigGet/Set`、`newComposerDispatchChain` 新增 `Config` 槽、`cw/panel_inbound.go` 新增 `if reply != ""`）＝这两枚用例"页面报回"那一跳的产码 |
| 4 | `480b970d` | 10-03 | 可能不能（读码推它无关，凭据＝`ip/pump.go` 是 `SnapshotPump`/`Marshal`/`Publish`，⛔ 引用 WebView/Eval/SetHtml）——列进来只为让**我这条脱罪链被跑打** |
| 5 | `4658dbb6` | 10-03 | ★能——`bringUp` 的 `WindowOptions` 字面量整块换成 `m.windowOptions()`，体内**每次建窗现调** `panelGeometrySource` → `config.LoadFile(cfgPath, nil)`（⛔ 盘读进了 STA 线程的建窗那一跳）；笔自标"半成品·未验证" |
| 6 | `cf95c799` | 10-03 | 可能不能（**仅整包序**）——`TestPanelHostRealWindowHopAndLifecycle` 的判据换形＋新增 `settleTreeReading`（`treeSettleWait = 3 * time.Second`）＝同进程真窗用例的时长/浏览器复用变了 |
| 7 | `32e74479` | 10-05 | 可能不能（**仅整包序**）——267-r2 种子迁移：`panel_pump_test.go` 的 `confirm_timeout_sec` 抬到 31、`executeOn145` 的 ctx 30s→45s |
| R | `f718e9b6` | 10-07 | **红端点**（`probes/303/orch/r3-f718e9b6-three-nails.txt` 里 nail2 已是那句 `title=""`） |

## 2. 发序（每一发都⛔ 是一条命令，⛔ 是一个台面）

**公共体**（一发＝一个仓外 clone，detach 到目标 sha，⛔ 动母仓工作树）：

```
D=/d/work/workspace/bisect/305
# 临时件只建不删：每发用独立目录 $D/<sha>，⛔ rm 任何既有目录
git clone --no-local "D:/work/workspace/projects plans/Wisp" "$D/<sha>"        # 建在仓外
git -C "$D/<sha>" checkout --detach <sha>
git -C "$D/<sha>" rev-parse HEAD                                              # 记进读数表
cp "D:/work/workspace/projects plans/Wisp/frontend/dist/index.html" "$D/<sha>/frontend/dist/index.html"   # ★AC13 那几发必须放；nail2 单发⛔ 放
sha256sum "$D/<sha>/frontend/dist/index.html"                                # 应逐字等于 9b7856b9…074f
cp -r "D:/work/workspace/projects plans/Wisp/third_party/sherpa-onnx" "$D/<sha>/third_party/"              # 或下面 PATH 指母仓绝对路径（二选一，逐发自报选了哪条）
cp -r "D:/work/workspace/projects plans/Wisp/build/"*.dll "$D/<sha>/build/"
PATH="/d/work/workspace/projects plans/Wisp/third_party/sherpa-onnx:/d/work/workspace/projects plans/Wisp/build:$PATH" \
  go test -C "$D/<sha>" ./cmd/wisp/ -count=1 -timeout 900s -v -run '<具名>'
```

尺形＝任务书那把：`go test -C <clone> ./cmd/wisp/ -count=1 -timeout 900s -v -run '<具名>'`，PATH 两枚条目**必须**是 `/d/…` 形（写成 `D:/…` ⇒ `0xc0000135` 且零 `--- FAIL`）。

**具名串**（⛔ 放宽、⛔ 改名、⛔ 造 SKIP）：
- nail2 单维＝`-run 'TestAC14GoSideEvalPushReachesThePage$'`
- AC13 单维＝`-run 'TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe$'`
- 三枚成对（票面现量的那把）＝`-run 'TestAC14AwaitedBindingReplyReachesThePage$|TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe$|TestAC14GoSideEvalPushReachesThePage$'`

### 发 0 ＝台面自检（1 发，⛔ 算归因）
`<sha>`＝`f718e9b6`，三枚成对具名串，**放 dist**。
判读：nail2 必须**发出 `--- FAIL` 且红句逐字**含 `Go's Eval push did not reach the document: the page reports its title as ""`；AC13 必须**⛔ 再 SKIP**（放了 dist 就该跑；它若仍 `--- SKIP` ⇒ §0 轴① 那形⛔ 成立，AC13 就**只剩母仓一发**可量，见 §4）。
⇒ 这两格任一不符＝**台面⛔ 健康**（DLL／dist／具名三坑之一），后面所有发⛔ 可比，先把台面钉死再跑。

### 发 1 ＝绿端点复现（★最值钱的一发，先于二分）
`<sha>`＝`416d9d56`，同一枚 clone 台面，三枚成对具名串，**放 dist**。
判读（这是编排者欠着的 **R2**）：
- nail2 **绿**（`--- PASS`）⇒ 症状②的绿↔红**在同台面上成立**，进发 2 二分；
- nail2 **红**（同一句 `title=""`）⇒ ★**"有绿点的回归"这一形在 clone 台面上⛔ 复现**（10-01 那枚绿是**母仓＋旧 dist 字节**面量的），必须具名改判为"这台面⛔ 复现绿点"，并把二分降级为"只在母仓面量的候选上单发对照"（母仓⛔ detach ⇒ 只能挑 1–2 枚候选用 `git show` 造补丁面，那是另一程的事，⛔ 我裁）。
- AC13 那一格在这一发**同时**给出"今日 dist 字节能不能复现 10-01 的绿"（＝R3 的可复现替代读法，⛔ 等于 R3）。

### 发 2–4 ＝二分 nail2（同一台面，⛔ 混发）
线性 9 点，取中：中＝**序 4 `480b970d`**（nail2 单维具名串，⛔ 放 dist 也⛔ 关它——建议照放，保持与发 1 逐字同台面）。
红⇒往低半（中＝序 1 `f7d28ef0`）；绿⇒往高半（中＝序 6 `cf95c799`）。**最多再加 2 发**（log₂8＝3）。
★**单调性闸门（这条⛔ 省）**：一旦出现"低半里有红、高半里有绿"的非单调色（这几枚候选机制互不相干，很可能非单调），**⛔ 信二分**，改**逐枚全跑**序 1–7＝7 发 nail2 单维，把 7 枚的颜色**逐枚列成名册**交回来。⛔ 报"二分定位到 X"。

### 发 5 ＝AC13 那一维（只在放了 dist 的 clone 上量）
把发 2–4 的具名串换成 AC13 单维，在**二分得到的同一枚 sha**＋**发 1 的 sha**各一发，两发⛔ 用不同 dist 字节（sha256 必须逐字相同并写进读数表）。
判读：AC13 与 nail2 **同枚翻色** ⇒ "部分同源"多一枚凭据；AC13 翻色的 sha ≠ nail2 翻色的 sha ⇒ **两枚症状各有各的笔**，⛔ 合并这一条要重裁（回来找我，⛔ 落地腿自己并）。

### 发 6 ＝跨维度负控的**前置读数**（票面 `AC#3` 那枚负控在这一程能拿到的那一半）
在 nail2 翻色的那枚 sha 上，用三枚成对具名串再发一次（同台面）。
判读：如果 nail1 绿、nail2 红、AC13 也红 ⇒ 与"⛔ 同源"相容；如果那一枚 sha 上 **AC13 绿而 nail2 红** ⇒ 症状①⛔ 是症状②的上游（部分同源里"共享上游"那一半被否），这一格是落地腿⛔ 能替的。

## 3. 每一发必须交回的字段（缺一格＝那发⛔ 可比）

`sha` ＋ `git -C <clone> rev-parse HEAD` ＋ clone 绝对路径 ＋ 具名 `-run` 串逐字 ＋ `rc` ＋ `--- PASS/FAIL/SKIP` 行**逐字**（红句逐字引，⛔ 只报枚数） ＋ `wc -c`/`sha256sum`/mtime 三枚（放了 dist 的发） ＋ PATH 两枚条目逐字（`/d/…` 形） ＋ `tasklist` 现量 `wisp.exe`／`balldebug.exe` 枚数。每把尺自己落一行 `rc=N`。

## 4. 判死到哪一格就⛔ 停（⛔ 我越界）

- 判死＝"哪一笔把颜色从绿推到红"（发 1–5 的产物）。**⛔** 判死＝"为什么"——`Eval` 落地那一层的三形（脚本被丢／落到另一份文档实例／文档在两发间被重建）我这程仍⛔ 分不开（＝票 §6 的 **R1/R4**，只有编译面能分），且**第四形我读码新添**：见 20 号件 §4 —— `post()` 在建窗后一律走 `w.Dispatch`（"Run() is that queue's reader"），而首扇窗的 `bringUp` 跑在 `drainTasks()` 里、`w.Run()`（`panel_resident_windows.go:299`）**尚未接管** ⇒〔读码推到〕最后一次 `SetHtml` 之后这条链**⛔ 手工泵点**（链内唯一泵点在 `firstRoundTrip` 的循环里，锚 `panel_host_windows.go:897`）。这一形⛔ 要新机制，一发就能分：**在具名串里加一枚同文件既有的 `TestAC13ColdStartPageOverEndsOnEntryContentNotTheProbe`（`panel_pageover_33r10_windows_test.go`，无头 `docSink` 驱动）** 对照——它若绿而 AC13 红，就是"落地/泵"那一族，⛔ 是"顺序"那一族。
- ⛔ 依赖面（`ExecuteScript(…, 0,)` 在第三方模块）、⛔ 泵形状（STA＋`Run()`）、⛔ `Q-84`、⛔ 放宽任何判据、⛔ 造 SKIP。
