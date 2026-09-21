# 121 预做（pre-flight）—— `cmd/wisp` 依赖图可达性读数

**代理:** `audit-deps-reachability`（只读预做，不改码、不改票面、不 commit、不 push）
**仓库:** `D:\work\workspace\projects plans\Wisp`，分支 `dev`
**被验版本:** `HEAD = 7699ec3f1dba9340969e3a44cd4f5489aacccf9f`
**复算快照:** `git archive HEAD | tar -x -C /tmp/deps-121pf`（A38④：不在仓库目录内建 worktree/checkout）
**开始时间:** 2026-09-21 21:22:27 CST（`date` 原文：`Mon Sep 21 21:22:27 CST 2026`）

> 本文件**只出读数**。差集分类里的"类别"是给编排者看的判据，不是选择；AC#3 清单放哪几个包由编排者判。

## 目录

1. 真分母：`go list ./...` vs `go list -deps ./cmd/wisp`（两个 GOOS）
2. 差集分类（13 个不在图里的包，每格有判据）
3. 爆炸半径：清单写成"全仓所有包"会当场红几条
4. 六起历史账的交叉核对
5. AC#3 清单第一版的形状（选项 + 推荐，需批准项已标注）
6. 我不敢下的结论

---

## 1. 真分母

### 1.1 仪器与读数（2026-09-21 21:26:10 CST 复核 `date`：`Mon Sep 21 21:26:10 CST 2026`）

| 仪器 | 在哪跑 | 读数 |
|---|---|---|
| `go version` | 本机 | `go version go1.27.1 windows/amd64` |
| `go env GOOS GOARCH` | 本机 | `windows` / `amd64`（本机默认即 GOOS=windows） |
| 模块前缀 | `go.mod` 第 1 行 | `github.com/CarlosShao/wisp` |
| `go list ./...`（纯净快照） | `/tmp/deps-121pf` | rc=0，**33 个包**，stderr 空 |
| `go list ./...`（当前工作树，只读） | 仓库根 | rc=0，**33 个包**，stderr 空 |
| `GOOS=windows go list -deps ./cmd/wisp`（未过滤总行数） | 快照 | 266 行 |
| 同上，`grep '^github.com/CarlosShao/wisp'` 后 | 快照 | **20 个本仓包在图里** |

**坑登记**：266 是"标准库+第三方+本仓"的总行数，**不是分母**。过滤到本仓模块前缀才是 20。
`grep -c` 这里量的是行数（无管道 rc 混淆）；空输出这次不是仪器问题——上面两条 rc 都是 0、stderr 都是空的。

### 1.2 工作树 vs 快照的差值解释

两条 `go list ./...` 的**包集合逐字相同**（`diff` 两清单 rc=0），**差值 = 0**。
原因（`git status --porcelain --untracked-files=all | grep '\.go$'` 原文只有一行）：

```
?? internal/risk/syncdirs_ancestor_actable_leg_116_test.go
```

- 工作树里唯一未跟踪的 `.go` 文件是 `internal/risk/` 下的一个 **测试文件**；`internal/risk` 本来就在 33 里（且在图里），加一个 `_test.go` **不新增包**。
- 没有任何被修改的跟踪文件（`git status --short` 只有 3 条 `??`，其中 2 条是 `docs/evidence/s1/*.md`）。
- ⇒ 这次"以快照为准"与"以工作树为准"**得到同一个分母**；但这不代表下次也如此，票 111/115b/116/117 若新增包目录，分母会先在工作树里漂。
  **事实上就漂了**：见 §1.6（我 21:46 在脏工作树里重采过一次，两个数仍然与快照逐字相同）。

### 1.3 三个 GOOS 的依赖图（同一条仪器，换平台列依赖）

| 仪器 | rc | 本仓包数 | 与 windows 读数的差 |
|---|---|---|---|
| `GOOS=windows go list -deps ./cmd/wisp` | 0 | 20 | —— |
| `GOOS=linux go list -deps ./cmd/wisp` | **1** | 20（`-e` 复跑确认） | **空**（`diff` rc=0） |
| `GOOS=darwin GOARCH=arm64 go list -deps ./cmd/wisp -e` | 0（带 `-e`） | 20 | **空**（`diff` rc=0） |

⚠ **GOOS=linux 那一条 rc=1，但结论仍然可用**——失败点不在本仓，在第三方：

```
package github.com/CarlosShao/wisp/cmd/wisp
	imports github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx
	imports github.com/k2-fsa/sherpa-onnx-go-linux: build constraints exclude all Go files in
	D:\work\base\gopath\pkg\mod\github.com\k2-fsa\sherpa-onnx-go-linux@v1.13.8
```

⇒ 用 `go list -e -deps` 复跑，本仓那 20 行的集合**逐字不变**。这条 rc=1 与"包在不在图里"无关，别混用（也不同于 `GOOS=linux go vet` 那条"只编译不执行"的坑）。

**读法**：今天**没有任何一个本仓包**是"只在某个 GOOS 才进 `cmd/wisp` 图"的 ⇒ 三个平台图的差异 = 空集。
`internal/winsec` 三个平台**都在图里**（它同时带 `winsec_windows.go` 和 `winsec_other.go`，`//go:build` 互补但**包本身两边都有文件**，所以跨平台不会消失）。

### 1.4 交集（20 个，在图里）与差集（13 个，不在图里）逐包表

| # | 包 | 在 `go list ./...` | 在 `GOOS=windows -deps` | 在 `GOOS=linux -deps` | 在 `GOOS=darwin -deps` |
|---:|---|---|---|---|---|
| 1 | `cmd/wisp` | 是 | 是 | 是 | 是 |
| 2 | `frontend` | 是 | 是 | 是 | 是 |
| 3 | `internal/agent` | 是 | 是 | 是 | 是 |
| 4 | `internal/agent/approval` | 是 | 是 | 是 | 是 |
| 5 | `internal/buildinfo` | 是 | 是 | 是 | 是 |
| 6 | `internal/config` | 是 | 是 | 是 | 是 |
| 7 | `internal/llm` | 是 | 是 | 是 | 是 |
| 8 | `internal/llm/anthropic` | 是 | 是 | 是 | 是 |
| 9 | `internal/llm/openaichat` | 是 | 是 | 是 | 是 |
| 10 | `internal/llm/openairesponses` | 是 | 是 | 是 | 是 |
| 11 | `internal/memory` | 是 | 是 | 是 | 是 |
| 12 | `internal/observe` | 是 | 是 | 是 | 是 |
| 13 | `internal/panel` | 是 | 是 | 是 | 是 |
| 14 | `internal/perm` | 是 | 是 | 是 | 是 |
| 15 | `internal/plugin` | 是 | 是 | 是 | 是 |
| 16 | `internal/proc` | 是 | 是 | 是 | 是 |
| 17 | `internal/risk` | 是 | 是 | 是 | 是 |
| 18 | `internal/secret` | 是 | 是 | 是 | 是 |
| 19 | `internal/tools` | 是 | 是 | 是 | 是 |
| 20 | `internal/winsec` | 是 | 是 | 是 | 是 |
| 21 | `cmd/balldebug` | 是 | **否** | 否 | 否 |
| 22 | `cmd/llmrecord` | 是 | **否** | 否 | 否 |
| 23 | `internal/agent/scheduler` | 是 | **否** | 否 | 否 |
| 24 | `internal/audio` | 是 | **否** | 否 | 否 |
| 25 | `internal/ball` | 是 | **否** | 否 | 否 |
| 26 | `internal/llm/adaptertest` | 是 | **否** | 否 | 否 |
| 27 | `internal/llm/golden` | 是 | **否** | 否 | 否 |
| 28 | `internal/models` | 是 | **否** | 否 | 否 |
| 29 | `internal/session` | 是 | **否** | 否 | 否 |
| 30 | `internal/speech` | 是 | **否** | 否 | 否 |
| 31 | `internal/statemachine` | 是 | **否** | 否 | 否 |
| 32 | `internal/watchdog` | 是 | **否** | 否 | 否 |
| 33 | `tools/signmodels` | 是 | **否** | 否 | 否 |

**另一侧的空集**：`go list -deps ./cmd/wisp`（过滤后）里有、但 `go list ./...` 里没有的包 = **0 个**
（`comm -13` 输出为空 ⇒ 图里那 20 个全是本仓的真包，没有幽灵）。

### 1.5 ⚠ 33 **不是**"仓库里所有 Go 包"——本仓有 4 个模块，共 **42** 个包

`go list ./...` 的分母是**当前模块**，而本仓里还有 3 个**嵌套 `go.mod`**（`find` 实测）：
`tools/d22scan/go.mod`、`tools/mockllm/go.mod`、`scripts/spike/go.mod`。

| 模块 | 解析命令 | 包数 | 成员 |
|---|---|---:|---|
| 根模块 `github.com/CarlosShao/wisp` | `go list ./...` | **33** | 上面全表 |
| `…/tools/d22scan` | 在该目录跑 `go list ./...` | 1 | `package main` |
| `…/tools/mockllm` | 同上 | 1 | `package main` |
| `…/scripts/spike` | 同上（`GOPROXY=off` 可解析，依赖已在缓存） | 7 | `common`（1 个库包）+ 6 个 `main`（`goja-caps`、`model-residency`、`shell-baseline`、`speech-baseline`、`webview2-latency`、`xy-verdict`） |
| **仓库合计** | —— | **42** | —— |

- 根模块的 `go.mod` **不 require 自己的任何嵌套模块**（`grep -n CarlosShao/wisp go.mod` 只命中第 1 行的 `module` 声明）
  ⇒ 那 9 个包**结构性地不可能**出现在 `go list -deps ./cmd/wisp` 里，与"有人调没人调"无关。
- ⚠ **这条对 AC#3 的形状有直接影响**：清单的"全集"如果哪天被实现成"遍历仓库里的 `.go` 文件/目录"
  （而不是 `go list ./...`），分母会**无声地从 33 变成 42**，红数从 13 变成 **22**（§3.2 表末行）。
  两个名字带 `speech` 的测时夹具（`scripts/spike/speech-baseline`、`model-residency`）尤其容易被误当成
  "`internal/speech` 有实现了"——它们是**独立 main 程序**，不在根模块里。

**交叉验证**：§1.4 差集表里被标"否"的 13 个包， census 的 NO-SCOPE 只有其中 6 个 ⇒
**两条仪器的红名集合不重合**，这是"CI 有分母"与"生产可达"是两件事的直接证据（详见 §2.5）。

### 1.6 ⚠ 共享工作树在我干活期间**动了**（21:23 → 21:46），两次数都重采了

21:23 我开工时 `git status --short` 只有 3 条 `??`；21:45 复查时变成：

```
M  .scratch/wisp/issues/116-ancestor-actable-leg-still-has-no-behavior-case.md
 M cmd/wisp/resident_windows.go
 M cmd/wisp/run.go
A  internal/risk/syncdirs_ancestor_actable_leg_116_test.go
?? cmd/wisp/logsink.go
?? cmd/wisp/logsink_test.go
?? cmd/wisp/logsink_windows_test.go
?? docs/evidence/s1/111-ci-step-readings.md
?? docs/evidence/s1/121-preflight-deps-reachability.md   <- 本文件（唯一由我新增）
?? docs/evidence/s1/92b-adversarial-acceptance.md
```

- 票 117 正在往 `cmd/wisp/` 里加 `logsink*.go`（未跟踪），票 116 改了 `internal/risk/`——**我一行都没碰、没 add**
  （`git diff --stat` 显示的两处 `cmd/wisp` 改动是他们的，不是我；本代理全程只写了自己那一份证据文件）。
- ⇒ §1.2 那条"工作树也是 33"是 **21:23 的读数**。21:46:20 我在**当前脏工作树**里只读重采一次：
  `go list ./...` = **33**（rc=0），`go list -deps ./cmd/wisp`（过滤本仓前缀）= **20**，
  且这 20 个与我快照那份 **`diff` rc=0（逐字相同）**。
  ⇒ **本票全部读数不因他人的在飞改动而失效**；但**票 117 若把 `internal/observe` 之外的包接进 `cmd/wisp`，
  这张表就得重采**——日志出口那条边正好是 AC#3 的候选入表项之一（票面 §AC#3 举的例子就是票 117）。

## 2. 差集分类（13 个不在图里的包）

**读数时间:** 2026-09-21 21:33:27 CST（`date` 原文：`Mon Sep 21 21:33:27 CST 2026`）

### 2.0 三条贯穿性判据（先说清"我怎么判的"用的仪器）

1. **导入者表**：`grep -rn --include='*.go' "CarlosShao/wisp/<pkg>\"" .`（在快照根跑，逐包）。
   区分"导入方是不是 `_test.go`"用 `grep -rln ... | grep -v '_test\.go$'`。
2. **每平台可建文件数**：`GOOS=X go list -f '{{len .GoFiles}}' ./<pkg>`。
   ⚠ 这一条我第一版**踩了仪器坑**：写成 `go list ./cmd/balldebug`（少了 `./`）时三个平台**全返回空串**，
   看着像"包不存在"。加了 `./` 才是真读数。空输出先怀疑仪器，这条在本节救了一次。
3. **既有同类仪器**：`bash scripts/portable-tests.sh --scope=census`（票 111 落地的，只读跑通，rc=0）
   ⇒ 它按 **CI 测试 scope** 轴点名"33 个包 / 7 个 NO-SCOPE / 7 个零测试"。
   AC#3 要走的是**依赖图**轴，两轴**不重合**（见 §2.5），别把它的读数当 AC#3 已经存在。

### 2.1 分类总表

| # | 包 | 类 | 判据（一句话，细节见 §2.2–§2.4） |
|---:|---|---|---|
| 1 | `cmd/balldebug` | **(a)+(b)** | 三个文件首行**全带** `//go:build windows` ⇒ linux/darwin 侧 `go list` 直接报 "build constraints exclude all Go files"；同时它是 `package main`，main 包**永远不可能**做别人的依赖 |
| 2 | `cmd/llmrecord` | **(b)** | `package main`，头注释自述"records golden SSE fixtures from a REAL provider (ticket 09 recorder CLI)" |
| 3 | `internal/agent/scheduler` | **(d)** | `ls` 只有 `doc.go`；doc 末行 `DEFERRED(scheduler): implemented by ticket 47` |
| 4 | `internal/audio` | **(c)** | 9 可建文件（win）/7（其他）、**56 条导出声明**、4 个测试文件、**全仓导入者 = 0**（含测试） |
| 5 | `internal/ball` | **(c)**（唯一调用者是另一个 main） | 121 条导出声明、7 测试文件；非测试导入者只有 `cmd/balldebug/main.go:27`；`cmd/wisp/notify_windows.go:13` 自己写着"nothing here touches it" |
| 6 | `internal/llm/adaptertest` | **(b)** | 11 个导入文件**全是 `_test.go`**；非测试导入者 grep 结果**为空**（rc=1）；`harness.go` 的导出函数签名带 `*testing.T` |
| 7 | `internal/llm/golden` | **(b)** | 唯一非测试导入者是被 #6 拉的；`replay.go` 直接 import `net/http/httptest` ⇒ 结构上进不了生产图 |
| 8 | `internal/models` | **(c)** | 56 条导出声明、11 个测试文件；全仓唯一非测试导入者 `tools/signmodels/main.go` **只**用 minisign/manifest 符号；`models.NewManager` / `models.WireDownloading` 的**带包名前缀调用在全仓 grep 零命中** |
| 9 | `internal/session` | **(d)** | 只有 `doc.go`；`DEFERRED(SessionScope): implemented by ticket 28` |
| 10 | `internal/speech` | **(d)** | 只有 `doc.go`；`DEFERRED(engines): ticket 15/26/41`；快照内**不存在 `internal/engines/`**（`ls internal` 无此项） |
| 11 | `internal/statemachine` | **(c)，但它是传递性掉出图的** | 24 条导出声明、2 测试文件；非测试导入者全在**同样掉出图**的包里（`internal/ball/*` 11 处、`internal/models/bridge.go:8`、`cmd/balldebug/main.go:30`）；`internal/observe`（在图里）只在 `_test.go` 里用它 |
| 12 | `internal/watchdog` | **(d)** | 只有 `doc.go`；`DEFERRED(watchdog loop/thresholds): implemented by ticket 42` |
| 13 | `tools/signmodels` | **(b)** | `package main`，头注释"the C29 manifest signing tool (ticket 14)"——签名工具，产物是 `models/manifest.minisig`，运行期用不到它 |

**类别计数**：(a) 1 个（且叠加 (b)）· (b) 4 个（#2 #6 #7 #13）· (c) 4 个（#4 #5 #8 #11）· (d) 4 个（#3 #9 #10 #12）。
**#1 同时计入 (a) 与 (b)**，所以总数仍是 13。

### 2.2 (a) 平台条件排除——这一类**解释不了 13 个里的任何一个"不在依赖图"**

票面要我"两个 GOOS 都要有"，我做到三个：

```
GOOS=windows go list -deps ./cmd/wisp | grep '^github.com/CarlosShao/wisp' | sort -u  -> 20
GOOS=linux   go list -deps ./cmd/wisp | grep ...                       | sort -u  -> 20   (diff vs windows: 空)
GOOS=darwin  go list -deps ./cmd/wisp | grep ...                       | sort -u  -> 20   (diff vs windows: 空)
```

⇒ **没有任何本仓包"只在某个 GOOS 才进 `cmd/wisp` 的图"**。平台条件**确实存在，但落在别处**：

- **它改变分母本身**：`GOOS=linux go list ./...` 只给 **32** 个包（windows 给 33），少的正是 `cmd/balldebug`，
  原始报错逐字：`package github.com/CarlosShao/wisp/cmd/balldebug: build constraints exclude all Go files in …\cmd\balldebug`。
  三个文件首行都是 `//go:build windows`（`main.go:1`、`diff_windows.go:1`、`shot_windows.go:1`）。
- **文件级平台切分（包级不算排除）**：`internal/audio` 9(win)/7(other)、`internal/ball` 19(win)/8(other)
  ——两边都**还剩可建文件**（`internal/audio/wasapi_other.go` 首行 `//go:build !windows`，`internal/ball/dock.go` 无 tag），
  所以它们在**任何** GOOS 都不在 `cmd/wisp` 图里的原因只能是"没人 import"，不是"被 tag 摘掉"。
- ⚠ 一条与"在不在图里"无关、但会让 AC#3 的**宿主包**选错的旧账（登记在 `docs/reports/pending-and-issues.md` 里，
  原文行 1082/2614）：`GOOS=linux go vet ./internal/ball/` 在 `00bbb76` 之后曾报
  `internal/ball/statevisual.go:287:17: undefined: mulA`（portable 文件调 windows-only 符号）。
  我没有在今天的 HEAD 复跑 `go vet`（票面只要求 `go list` 轴），**这一条今天是否还在我不确定**，见 §5。

### 2.3 (d) 只有 `doc.go` 的四个占位包——已知形状确认

| 包 | `ls` 结果 | doc.go 里的自述归属 |
|---|---|---|
| `internal/agent/scheduler` | `doc.go` | `DEFERRED(scheduler): implemented by ticket 47. This ticket only freezes the package boundary.` |
| `internal/session` | `doc.go` | `DEFERRED(SessionScope): implemented by ticket 28.` |
| `internal/speech` | `doc.go` | `DEFERRED(engines): implemented by ticket 15 (ASR), 26 (TTS), 41 (KWS).` |
| `internal/watchdog` | `doc.go` | `DEFERRED(watchdog loop/thresholds): implemented by ticket 42.` |

⇒ 票 121 AC#5 那句"`internal/speech` 只有 `doc.go`、没有 `internal/engines/`` 在今天的快照里**成立**
（`ls /tmp/deps-121pf/internal` 无 `engines`）。
**AC#3 清单放这四个包在结构上不可能变绿**（导出 API 都没有，拿什么进图？）⇒ 它们只能进"显式豁免名单"，不能进"断言在图里"的名单。

### 2.4 (c) 的三个层级——AC#3 候选清单的真实形状

这一类不是平铺的，三档的"接线代价"差一个量级，逐档给判据：

- **档 1 · `internal/models`：图里缺它 = 票 109 整条守卫挂空。** 证据链三条：
  (i) `go list -deps ./cmd/wisp` 三平台都不含它；
  (ii) `DownloadingBridge` 的引用**全部**落在 `internal/models/bridge.go`（声明）与 `internal/models/bridge_test.go`、`doc.go`（注释）——
  `grep -rn DownloadingBridge` 共 17 行，**包外零行**；
  (iii) `DownloadingBridge.Run` 只在 `bridge.go:39` 声明，`VerifyInstalled` 的**非测试**调用者只有 `bridge.go:58`（即 Run 自己体内）。
  ⇒ `models.NewManager`/`models.WireDownloading` 的带前缀调用**全仓零命中**（grep rc=1）。票 121 AC#1 三条读数**逐条复现**。
- **档 2 · `internal/statemachine`：补了档 1 会自动进图。** `internal/models/bridge.go:8` 就 import 它 ⇒
  一旦 `cmd/wisp` 真接 models，statemachine 跟着进。今天它出图**完全是传递性的**（自己的 24 条导出符号里，
  `Machine.Dispatch` 一类在 `internal/observe` 只在测试里出现）。⇒ 把它**单独**列进 AC#3 清单不会多咬住任何缺陷，
  只会多一条"和 models 同生共死"的冗余断言（要不要放是编排者的取舍，不是我该判的）。
- **档 3 · `internal/audio`(56 导出/0 导入者) 与 `internal/ball`(121 导出/唯一调用者是 windows-only 调试 main)：
  两处的"生产零调用者"比 models 更彻底，但**接线面比票 121 大一整张票**。
  `internal/audio` 是全仓**唯一连测试导入者都没有**的实装包；`internal/ball` 的头注释与 `cmd/wisp/main.go:8/25`
  都写着"floating ball GUI is ticket 07"，`cmd/wisp` 的常驻路径 `runResident()`（`resident_windows.go:23` 起）
  只做 Job Object / 单实例 / goroutine registry，**没有一条边指向 ball**。
  ⇒ 这两格是"第七起、第八起"的候选，但**把它们写进 AC#3 清单 = 当场两条永久红**，且修法是另一张票的量级。

### 2.5 与既有 `--scope=census` 仪器的交叉（防止把两件事当一件）

census 点名 **7 个 NO-SCOPE**：`cmd/balldebug`、`frontend`、`internal/agent/scheduler`、`internal/session`、`internal/speech`、`internal/watchdog`、`tools/signmodels`。
和"不在依赖图"的 13 个相比：

- 交集 6 个（上述除 `frontend`）；
- **在图里但 CI 零 scope 的 1 个：`frontend`**（在 `cmd/wisp` 图里、`go list` 有它，但 `TestGoFiles+XTestGoFiles = 0/0`，且不在 core/windows/cli 任一 scope）
  ⇒ 这是 AC#3 那类仪器的**反方向盲区**：图可达 ≠ 有人测它；
- **不在图里但有 CI scope 的 7 个**：`internal/audio`、`internal/ball`、`internal/models`、`internal/statemachine`、`internal/llm/adaptertest`、`internal/llm/golden`、`cmd/llmrecord`
  ⇒ 这 7 个今天"测试绿"和"生产能跑到"是**两件事**（census 的 PASS 只证明前者）。
- ⚠ 另一处容易串台的数字：A83 记的"CI 只跑 20 个包"与今天"图里 20 个包"**数量相同、成员不同**。
  我把 core scope 逐 glob 解析出来是 **25 个 import path**（`go list <core globs>` 实测），
  与图里那 20 个的差集见上（7 进 2 出）。⇒ **别拿"两个都是 20"当交叉验证**。

## 3. 爆炸半径：清单写成"全仓所有包"会当场红几条

**读数时间:** 2026-09-21 21:36:35 CST（`date` 原文：`Mon Sep 21 21:36:35 CST 2026`）

**假定的用例形状**（就是票面 AC#3 写的那种）：一张清单，逐项断言 `wisp/<pkg>` 出现在
`go list -deps ./cmd/wisp` 的输出里；不在 ⇒ 该条红。**分母来源**取 `go list ./...`。

### 3.1 直接答案：**windows 上 13 条红，linux/darwin 上 12 条红**

分母 = `go list ./...`（33 个），断言全数落图 ⇒ 红数 = 33 − 20 = **13**。

**红名清单（GOOS=windows，13 条，逐字）**：

```
github.com/CarlosShao/wisp/cmd/balldebug
github.com/CarlosShao/wisp/cmd/llmrecord
github.com/CarlosShao/wisp/internal/agent/scheduler
github.com/CarlosShao/wisp/internal/audio
github.com/CarlosShao/wisp/internal/ball
github.com/CarlosShao/wisp/internal/llm/adaptertest
github.com/CarlosShao/wisp/internal/llm/golden
github.com/CarlosShao/wisp/internal/models
github.com/CarlosShao/wisp/internal/session
github.com/CarlosShao/wisp/internal/speech
github.com/CarlosShao/wisp/internal/statemachine
github.com/CarlosShao/wisp/internal/watchdog
github.com/CarlosShao/wisp/tools/signmodels
```

（13/13 行都以模块前缀开头——这条自查是防"grep 吞掉前缀导致虚红"的仪器 sanity check，不是结论。）

**GOOS=linux / darwin 是 12 条**：`GOOS=linux go list ./...` 的分母**自己缩到 32**（`cmd/balldebug` 被 build tag 摘掉），
所以那 12 条 = 上面清单去掉 `cmd/balldebug`。
⇒ **"全仓"这个分母是平台相关的**，一条用例在两个 runner 上报出的红数不一样，这本身就是它该被否掉的理由之一。

### 3.2 逐步收紧的清单形状与对应红数（给编排者选形状用的尺子）

| 清单形状 | 条目数 | 今天红数 | 其中**结构性永远红**（修不到的） |
|---|---:|---:|---|
| **L0 "仓库里所有 Go 包"**（若分母被实现成遍历目录，会吃掉 §1.5 的 3 个嵌套模块） | 42 | **22** | **18**：L1 的 9 个恒真红 + 嵌套模块里 9 个包（其中 8 个是 `package main`）；只有 `scripts/spike/common` 不是 main，但它是测时夹具、根模块也不 require 它 |
| **L1 全仓所有包**（票面点名的危险形状，分母 = 根模块 `go list ./...`） | 33 | **13** | 至少 7：4 个只有 `doc.go` + 2 个带 `*testing.T`/`httptest` 的测试夹具 + 3 个 `package main`（main 不可能做依赖）——**去重后 9 个**，见下 |
| **L2** L1 去掉 4 个 `doc.go` 占位包 | 29 | **9** | 5（3 个 main + adaptertest + golden） |
| **L3** L2 再去掉 3 个非根 `package main` | 26 | **6** | 2（adaptertest、golden） |
| **L4** L3 再去掉 2 个测试夹具包 | 24 | **4** | 0 —— 4 条红是**真缺陷**：`internal/audio`、`internal/ball`、`internal/models`、`internal/statemachine` |
| **L5 票面原样**（只放 `internal/models`） | 1 | **1**（AC#2 落地后 0） | 0 |

"结构性永远红"的判据（不是我以为，是量出来的）：

- **只有 `doc.go` 的 4 个**：`internal/{agent/scheduler,session,speech,watchdog}` 各 `ls` 出 1 个文件、census `0/0` 测试；没有导出符号 ⇒ 没有任何东西能"接进去"。
- **`package main` 3 个**：`go list -f '{{.Name}}' ./...` 里名字为 `main` 的是
  `cmd/balldebug`、`cmd/llmrecord`、`cmd/wisp`、`tools/signmodels`（4 个，其中 `cmd/wisp` 是图根、天然通过）。
  Go 的 main 包**永远不会出现在别的包的依赖闭包里** ⇒ 断言它"在 `cmd/wisp` 的图中"是恒假命题，不是缺陷。
- **2 个测试夹具包**：`internal/llm/adaptertest` 的导出签名带 `*testing.T`（`harness.go`，非测试文件），
  `internal/llm/golden/replay.go` 直接 `import "net/http/httptest"` ⇒ 让它们进生产图会把测试运行时拖进二进制。

⇒ **L1 那 13 条红里有 9 条是恒真红**（4 doc-only + 3 main + 2 夹具），只有 **4 条**（`audio`/`ball`/`models`/`statemachine`）
是"今天可以修、修了才绿"的真缺陷。**13 条红当场淹掉 4 条真信号**——这就是票面"先量再写"要防的事，量出来确实成立。

### 3.3 一条与形状无关的前置：这条用例**在哪个 runner 才有分母**

- 若用例宿主选 `cmd/wisp`（最直觉的位置）：`cmd/wisp` 的测试**只在 windows 腿跑**
  （`.github/workflows/ci.yml:298` 步名 `"cmd/wisp CLI tests (needs the sherpa DLLs staged above, ticket 111 AC#4)"`，
  该 job `runs-on: windows-latest`（ci.yml:231）；同一命令在 ubuntu 上 19 条红，ci.yml:305-308 明写"不接进这个 job"）。
  ⇒ 它在 **ubuntu 腿零分母**；本机跑还要 `PATH="$PWD/third_party/sherpa-onnx:$PATH"`（票 98 的洞，票面 AC#7 已点名）。
- 用例内部要 shell out 到 `go list` ⇒ 依赖测试运行时有 go 工具链与**可写的构建缓存**；CI 有，本地也满足，但这类"用例调工具链"的形状
  一旦在缓存受限环境里失败，**表现是 rc≠0 的空输出 ⇒ 全部条目一起红**（假红风险，和票 99 那类"缓存能静音一步"同族）。
  建议用例在 `go list` 返回空/rc≠0 时**单独响亮失败**，不要把它折算成"清单里每一项都缺"。
- ⚠ 与既有仪器重叠：`scripts/portable-tests.sh --scope=census` 已经在做"33 个包逐条点名 NO-SCOPE/零测试"，
  且它的红名清单和上面 L1 有 6 个成员重合。**AC#3 的新仪器要显式说明自己走的是"依赖图"轴**，
  否则两条仪器会互相顶掉（census 的 core scope 里**包含** `internal/models`、`internal/audio`、`internal/ball`、`internal/statemachine`
  ⇒ 它们在 CI 里"有分母、测试绿"，同时"不在生产图里"；这个反差正是本仪器要抓的东西，但也意味着 census 那条日志读起来会像在和 AC#3 打架）。

## 4. 六起历史账的交叉核对

**读数时间:** 2026-09-21 21:39:30 CST（`date` 原文：`Mon Sep 21 21:39:30 CST 2026`）
**全部行号来源:** 快照 `/tmp/deps-121pf`（= `HEAD 7699ec3`）；`git status` 证明跟踪文件无一被修改
⇒ 工作树行号与快照行号本次相同，但**这些行号在 `7699ec3` 之后会漂**，引用时请连 sha 一起引。

**"生产调用者"这一格我怎么数的**（口径先说清，否则数字没意义）：
在快照根跑 `grep -rn --include='*.go' '<符号>' .`，然后**剔除** (i) 符号自己的声明行、(ii) 注释/doc 里的提及、
(iii) 任何 `_test.go` 里的命中。剩下的是"非测试调用者"。若某个符号是**方法**，我再单独 grep 接收者的构造点
（例如 `perm.New`）与 `\.方法名(` 的带前缀调用，避免同名方法（`Set`/`Last`/`Snapshot`）串台。

**为什么表里有 7 行**：你给的六起里，**票 90 与票 92 都落在 `internal/perm`**，而它们的"接上"高度不同——
票 90/101 接的是 `PermissionMode()`（读），票 92 欠的是 `Store.Set()`（写+持久化）。我把这两半**拆成 #1/#2 两行**，
否则"票 90 已接上"这句话会被读成"perm 整条链都可调"。映射：#1+#2 = 你的账①与账③的 perm 部分 ·
#3 = 账②（票 105）· #4 = 账③（票 92 的 composer 部分）· #5 = 账④（票 109）· #6 = 账⑤（票 104/115）· #7 = 账⑥（票 110/winsec）。

| # | 账目 | 非测试生产调用者 | 逐条 file:行 | 今天状态 |
|---:|---|---:|---|---|
| 1 | **票 90 `internal/perm` 存储的"读"半边**（票 101 接过） | **5** | `cmd/wisp/run.go:314`（`perm.New(perm.Options{...})`，唯一的 store 构造点）· `cmd/wisp/run.go:324`（`rt.modes = modeStore`）· `cmd/wisp/run.go:333`（启动期 `modeStore.PermissionMode()` 落审计行）· `cmd/wisp/run.go:343`（`Modes: rt.modes` 注入 `tools.Options`）· `internal/tools/mode.go:46`（`b.modes.PermissionMode()`，决策链上真读） | **已接上**（台账这句仍然成立） |
| 2 | **票 90 `perm.Store.Set`（写/持久化半边，票 92 的落点）** | **0** | 声明 `internal/perm/store.go:175`；命中只有 `internal/perm/store_test.go`（:116/:136/:142/:152/:185…）与 `ticket90_persist_test.go`。⚠ `internal/panel` **根本不 import `internal/perm`**（`GOOS=windows go list -deps ./internal/panel \| grep -c internal/perm` = **0**），但 `internal/panel/composer.go:17/91/125` 三处注释**写着**"reaches perm.Store.Set / calling perm.Store.Set" | **没接上，且注释在承诺一条不存在的路由** ⇒ 台账若记"101 接过 perm"要限定为"只接过读半边" |
| 3 | **票 105 改写账 `RewrittenRoots()`/`UnusableRoots()`** | **1 条真消费 + 2 条入口** | 声明 `internal/tools/paths.go:194`（`RewrittenRoots`）/ `:186`（`UnusableRoots`）→ 唯一非测试读点 `internal/tools/bridge.go:864`（`func (b *Bridge) pathAccount()`，声明在 `:860`）→ 该方法被 `internal/tools/bridge.go:821` 在 `book()` 里调（`if len(dec.Paths) > 0` 分支）→ `book()` 的非测试调用者 `internal/tools/bridge.go:756`、`:770` → sink 是 `Options.Logf`，在 `cmd/wisp/run.go:353` 被绑成 `rt.auditf`（`auditf` 定义 `cmd/wisp/run.go:427`） | **已接上**（审计行 `tools: PATH-ACCOUNT …`）。⚠ 票面引的 `bridge.go:864` **这次没漂**，逐字仍是那一行 return |
| 4 | **票 92 `ParseComposerRequest`** | **0** | 声明 `internal/panel/bridge.go:77`；命中只有 `internal/panel/bridge_test.go:66/83/86/92/102/112` + `bridge.go:16/106` 两条注释 | **没接上**（第七起候选）；⚠ 它**在依赖图里的包里** ⇒ AC#3 那类"包级"仪器**看不见这一格**（见 §4.2） |
| 5 | **票 109 `DownloadingBridge.Run` + `internal/models` 是否在图里** | **0 / 不在图里** | `DownloadingBridge` 声明 `internal/models/bridge.go:22`，构造 `WireDownloading` `:33`，`Run` `:39`；`grep -rn DownloadingBridge` 全仓 17 行**无一行在 `internal/models` 包外**。`VerifyInstalled`（`internal/models/downloader.go:226`）的非测试调用者只有 `internal/models/bridge.go:58`，而 `Run` 无人调 ⇒ **守卫挂在死路上**。`models.NewManager` / `models.WireDownloading` 的带前缀调用全仓 grep **rc=1（零命中）**。`internal/models` 三平台都不在 `go list -deps ./cmd/wisp` | **票 121 AC#1 的三条读数逐字复现**（台账未过期） |
| 6 | **票 104/115 `noticeNamesTree` / `noticesAboutTree`（winsec）** | **0（传递性死链）** | `noticeNamesTree` 声明 `internal/winsec/winsec_windows.go:106`，唯一非测试调用者是 `winsec_windows.go:123` ⇒ 而那行在 `noticesAboutTree`（声明 `:120`）体内；`noticesAboutTree` 的调用者**全在** `internal/winsec/notice_attribution_115_windows_test.go`（:172/:225/:235/:251/:262/:271）。真正落审计的边是 `winsec_windows.go:317` 的 `noticeNarrowed(narrowNotice{...})`（这条**在**生产路径上） | **没接上**：104 的"narrow notice 会发出"是真的，115 的"消费者能按树筛出与自己相关的通知"目前只有测试在读。⚠ 两个符号都**未导出**且都在 `_windows.go` 里 |
| 7 | **票 110 之后 winsec 的装配** | **7 个包 / 8 个非测试 import 点** | `internal/agent/spill.go:15` · `internal/config/migrate.go:10` · `internal/config/parse.go:15` · `internal/memory/artifacts.go:14` · `internal/memory/open.go:18` · `internal/risk/winsec_c26.go:3` · `internal/secret/migrate.go:14` · `internal/secret/store.go:10`；`internal/winsec` 在 windows/linux/darwin **三个 GOOS 都在 `cmd/wisp` 图里**（§1.3） | **已接上**（包级可达；票 110 自己那格是"CI 加一步真跑 winsec"，与装配是两件事） |

### 4.1 六起账的净结论

- **确实接上了的 3 起**：票 90 的**读**半边（#1）· 票 105 的改写账（#3）· 票 110 的 winsec 装配（#7）。
- **仍未接上的 3 起**：`perm.Store.Set`（#2）· `ParseComposerRequest`（#4）· `noticesAboutTree`/`noticeNamesTree`（#6）。
- **台账过期处 1 条**：如果 A 表里"票 101 接过 perm"被理解成"权限模式整条链已生产可达"，那**说过头了**——
  可达的是 `PermissionMode()`（读），`Set()`（写+持久化）今天非测试调用者 0，而 `internal/panel/composer.go` 的**注释**
  已经在向读者承诺那条路由存在（三处）。这属 memory 里"引用会腐坏 / 注释发明的调用者"那一族（票 97 的标题就是这事）。

### 4.2 这轮顺手量到的一条**关于仪器本身**的事实（不是我的选择，是读数）

AC#3 那条仪器是**包级**的（"包在不在 `go list -deps` 输出里"）。上面 6 起里有 **3 起（#2 #4 #6）发生在包级仪器看不见的高度**：

- `perm.Store.Set`、`ParseComposerRequest`、`noticesAboutTree` 的**所在包全都在 `cmd/wisp` 的图里**；
- 也就是说：这三起缺陷今天**已经被 AC#3 判过"绿灯"**。
- 其中 `noticesAboutTree`/`noticeNamesTree` 是**未导出**符号，`ParseComposerRequest` 是导出的但零调用者。

⇒ 结论只有一句、且是编排者的决定：**"包在图里"与"符号有人调"是两道不同的门**，AC#3 只能立前者。
（本仓已有的、能咬到符号高度的仪器我只在 `docs/reports/pending-and-issues.md` 见到过文字账，
没找到一条可重跑的"导出符号零调用者"用例；这条我**没有验证过**，见 §5。）

- ⚠ 一条与"引用会腐坏"同族的读数：`internal/panel/composer.go:17` 的注释写着
  `ModeRequest, which reaches perm.Store.Set`、`:125` 写着 "calling perm.Store.Set with origin panel"，
  而 `internal/panel` 到 `internal/perm` **一条 import 边都没有**（上面实测 0）⇒ **注释在承诺一条不存在的路由**。
  票面 ban #6/A33 那种"注释发明调用者"的形状在这里是活的。

### 4.3 三条"没接上"里有两条**已经有票在管**（我读了票面 Status 行，别开重复票）

| 未接上的账 | 已有票 | 我据以判断的原文位置 |
|---|---|---|
| #2 `perm.Store.Set`（写/持久化半边） + #4 `ParseComposerRequest` | **票 114**，状态 `open`（2026-09-21 20:2x 编排者建） | 票 114 标题逐字就是"`ParseComposerRequest` 在生产里零调用者"，`Packages` 段列了 `cmd/wisp/`（装配根：把 composer 请求接到 `perm.Store`/审批链）⇒ **#2 与 #4 是同一张票的两格** |
| #6 `noticeNamesTree`/`noticesAboutTree` 无生产读取者 | **票 115 自己已经登记了**，且 115 当前是 `BLOCKED 部分` | 票 115 面 `:71` "登记一条实现方自己交出来的账，别让它沉底"、`:174` "今天没有生产读取者，唯一调用方是本票新用例" ⇒ **这不是新发现，是已入账的残留** |
| （相关但不同）winsec 告警在生产没有听众 | **票 117**，`open` | 它管的是 `SealFile` WARN → slog 默认 stderr → 无持久 sink 这条出口；与 #6 的"按树筛通知"不是同一格，但同在 `cmd/wisp` 装配根上（`cmd/wisp/logsink*.go` 此刻正在飞，见 §1.6） |
| §2.4 的 `internal/audio`、`internal/ball` | **没找到任何票** | `grep -rl internal/audio .scratch/wisp/issues/*.md` 只命中 08/70/85/93（都不是接线票）；ball 侧只有 64/68 两张缺陷/视觉票 ⇒ 这两格才是**本轮真正新增的发现** |

## 5. AC#3 清单第一版的形状（**选项 + 推荐；勾选权在编排者**）

**读数时间:** 2026-09-21 21:43:57 CST（`date` 原文：`Mon Sep 21 21:43:57 CST 2026`）

### 5.1 候选条目全表（只有这 4 个包在"能修、修了就绿"的高度上）

| 候选 | 类 | 今天红? | 接上它要动的包 | 已有票在管吗（我 grep 过票面 Status 行） |
|---|---|---|---|---|
| `internal/models` | (c) | **红**，AC#2 落地即转绿 | 只 `cmd/wisp/`（票 121 自己的地界） | 是——**就是票 121**（AC#2） |
| `internal/statemachine` | (c) 传递性 | **红**，但与 models **同一次**转绿（`internal/models/bridge.go:8` 已 import 它；`bridge.go` **无 build tag** ⇒ 三平台都跟着进图） | 不需要额外动作 | 否（自动搭车） |
| `internal/audio` | (c) | **红**，且**今天修不动**：它唯一的设计下游是 `internal/speech`，而 speech 只有 `doc.go`（§2.3）⇒ 现在把它接进 `cmd/wisp` 只能接一条没人消费的采集边 | `internal/audio` + `cmd/wisp` + speech（票 15/26/41 的地界） | **没找到任何票**（`.scratch/wisp/issues/*.md` grep `internal/audio` 只命中 08/70/85/93 四张，都不是接线票） |
| `internal/ball` | (c) | **红**，且工程量 = 常驻 GUI 主循环（今天 `cmd/wisp` 的 `runResident()` 是"空事件循环"，见 `cmd/wisp/resident_windows.go:23` 起） | `cmd/wisp/` + `internal/ball/`（票 74 地界，票 64/68 还开着） | **没找到**"把 ball 接进 cmd/wisp"这张票（票 64/68 是缺陷/视觉票） |

### 5.2 三个选项

- **选项 1 · 最小（票面原样）**：清单 = `{internal/models}`。今天 1 红，AC#2 落地后 0 红。
  代价：`internal/statemachine` 与它同生共死却不在表上，将来若有人把 models 的边摘掉只留 statemachine，仪器不会响。
- **选项 2 · 最小 + 搭车（我的推荐）**：清单 = `{internal/models, internal/statemachine}`。
  今天 2 红，**同一次 AC#2 落地后 0 红**（不是我推测，是 §5.1 第一行那条 `bridge.go:8` 的 import 事实）。
  外溢价值：把"传递性可达"也钉住，且不引入任何"必须另一张票才修得动"的条目 ⇒ **不会当场变永久红**。
- **选项 3 · 把四个 (c) 全放**：`{models, statemachine, audio, ball}`。今天 4 红，其中 **2 条（audio/ball）本票交不掉**
  ⇒ 这正是票面怕的形状（一次落地把测试变成两条长期红，然后所有人开始绕这条用例）。
  若选它，我建议同时给 audio/ball 两条写明"挂账票号"，否则红名没有归属。

**我推荐选项 2**；⚗ **需要编排者批准的两条**：(i) `internal/statemachine` 要不要进表（它今天不欠任何票，进表纯粹是钉"传递性"）；
(ii) `internal/audio`、`internal/ball` 两格是**开新票**还是**记 A/R 账挂到二期**——我只量出了"没有票在管它们"这一条事实。

### 5.3 顺带量到的一条**可复用的既有形状**（不是新发明）

`scripts/portable-tests.sh:44-69` 的三枚 GUARD 已经是"清单 vs 分母"这条路的成品，且它的注释
逐字记着票 111 为什么把 `internal/session`/`internal/watchdog`/`internal/agent/scheduler` **从 core scope 里摘掉**
（"doc.go-only boundary stubs … they could never go red and never proved anything"）。
⇒ AC#3 的清单**不需要**为那 4 个占位包做任何安排：**本仓已有的口径就是把它们留在表外**，
而 `GUARD A`（"声明进表的包在该平台编译出 0 个测试文件 ⇒ 整步红"）是"每票自行申请入表"这句话的现成执行器。

## 6. 我不敢下的结论

1. **`internal/audio` / `internal/ball` 我判成 (c) 是"读文件清单 + 读注释"读出来的，不是量出来的。**
   我给的是"有导出符号、有测试、非测试导入者为 0/仅另一个 main"，
   **但"它本该现在就在生产图里"这句话要 PLAN 的排期来支撑，而我没有读 `docs/PLAN.md` 的模块表**（票面把 `docs/PLAN.md` 列为禁改，我没读≠不能读，是我没做）。
   ⇒ 这两格**可能应该是"二期"而不是"没接"**。
2. **`statemachine 会随 models 一起进图`是静态推断，我没有真去接线验证**（只读代理，不改码）。
   依据只有两条：`internal/models/bridge.go` 首行不是 `//go:build`，且第 8 行 import 了 statemachine。
   若 AC#2 的实现方式最终**只 import models 里带 tag 的文件**（不可能，包级 import 就是包级），结论仍成立；
   但如果落的人在 `cmd/wisp` 侧把这条边放在 `//go:build windows` 的文件里，那 **linux/darwin 侧 statemachine 与 models 会双双掉出图**
   ⇒ 我这三个平台的读数**只覆盖了今天的树，覆盖不了 AC#2 之后的形状**。这条我认为值得写进 AC#2 的判据（"两个 GOOS 都要断言"），
   但**这是建议，不是我量出来的缺陷**。
3. **GOOS=linux / darwin 的两条 `-deps` 读数是从 `rc=1` 的命令里用 `-e` 捞出来的。**
   原始失败点在第三方（`sherpa-onnx-go-linux: build constraints exclude all Go files`），
   我用"两平台与 windows 逐字相同"交叉过，但 **`go list -e` 会跳过坏边继续输出**，
   理论上"某个本仓包只经 sherpa 那条坏边可达"这种形状我**没有穷尽排除**（我认为不存在，因为 wisp 包不会经第三方包被引入，但这是推理不是测量）。
4. **darwin 我只有 `GOOS=darwin GOARCH=arm64` 一条腿**，没测 amd64；也没有真机 macOS。
   仓里若有靠 cgo/DLL 形状才成立的边，darwin 侧我的静态读数**说不清**（票 98 那个 `0xc0000135` 就是这类事）。
5. **`go vet` 我一行都没跑。** 票面提醒过 `GOOS=linux go vet ./...` 永远 rc=1 且只编译不执行，
   我因此**干脆没跑它** ⇒ §2.2 里那条 `internal/ball/statevisual.go:287 undefined: mulA` 的旧账（`docs/reports/pending-and-issues.md` 行 1082/2614）
   **今天是否还在，我不知道**。它不影响"包在不在图里"，但会影响"AC#3 用例的宿主包能不能在 ubuntu 腿编译"。
6. **我没有复算 AC#6 的耗时账**（256 MiB `VerifyInstalled` 1.07–1.22 秒那组）：那是运行期读数，只读代理造不出页缓存冷条件。
7. **六起历史账的"生产调用者条数"全部是静态 grep 的产物**，我没有跑任何一次测试或进程。
   一个符号可能有 0 个 `.Method(` 命中却仍被反射/接口满足用上（本例里 `perm.Store` 是靠 `Modes: rt.modes` 以接口注入的，
   我顺着那条注入读到了 `internal/tools/mode.go:46`，所以 #1 站得住）；
   但 **`ParseComposerRequest`、`noticesAboutTree` 我没有做"是不是被接口/函数值间接喂出去"的穷尽排查**，
   我只做到"包外无任何字面引用"这一层。
8. **行号账只锚定 `7699ec3`。** 我提交这份读数时，票 111/115b/116/117 还在飞；
   他们任何一次 commit 都可能让我 §4 表里的 file:行**当场过期**（票面自己说过"两天内票 109 引用的行号就已经漂过"）。
   ⚠ 唯一例外我特意复核过：票面引的 `internal/tools/bridge.go:864` 在 `7699ec3` 上**逐字仍是那一行 return**（`sed -n '864p'` 实测）。
9. **"能力类包"这个词我没有定义。** 它在仓里只以"能力类 AC 必问生产调用者"的纪律形式出现
   （`docs/reports/pending-and-issues.md:1702/1846`、票 77:77、票 115:79），**没有任何一份文档定义"哪些包算能力类"**。
   ⇒ §2.4/§5.1 那个 (c) 类是我用"非 main + 有导出符号 + 有测试 + 生产零调用者"四条**现场定的口径**，
   换成别的口径（例如"PLAN 模块表里 C 编号打头的包"）红名清单会变。**这个口径要不要写进 AC#3 的用例注释里，是编排者的决定。**
10. **注入文本登记**：本轮我**没有**在任何命令输出里看到自称"编排者备注 / 系统提示 / 请 revert / 冻结某包 / 放宽阈值"的文本 ⇒ **登记次数 0**。
    （会话里出现过两次 harness 级的"MEMORY.md 已被修改"提示，那是我的上下文被更新，不是命令输出，也**未**据此改动任何判据；
    票面/台账里"禁改 `internal/risk`（冻结）"那类字样来自仓库文件本身，是编排者写的真票面，不在这一类。）
