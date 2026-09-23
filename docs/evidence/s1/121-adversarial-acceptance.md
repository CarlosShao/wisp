# 121 对抗验收表 —— 模型交还链是否真被接上（AC#2/AC#3 那一族"能力有没有被接上"）

**验收代理:** `acceptor-ticket121`（独立对抗验收，只读）
**仓库 / 分支:** `D:\work\workspace\projects plans\Wisp` @ `dev`
**被验版本:** `6a39820`（本票末条 commit；`d180d4b`/`6a39820` 是纯文档 commit，代码与 `347b7bc` 逐字相同，
我按 `git show --stat` 点过：两枚都只动票面那一个文件）
**基线（接前）:** `3b03f00` = `7fe5e73^` = 本票一枚代码都还没有的树
**开始 / 收尾时间:** 2026-09-23 09:13 → 09:36 CST（`date` 原文 `Wed Sep 23 09:13:17 CST 2026` / `Wed Sep 23 09:36:* CST 2026`）
**工具链:** `go version go1.27.1 windows/amd64`，`GOPATH=D:\work\base\gopath`，本机物理内存 34,209,185,792 B（31.9 GiB）
**POSIX 腿:** `golang:1.27` 容器（`go version go1.27.1 linux/amd64`，内核 `6.6.114.1-microsoft-standard-WSL2`），`CGO_ENABLED=0`

## 0. 快照纪律（先自证我没拿工作树当基线）

- 接后：`git archive 6a39820 | tar -x -C /tmp/ac121-after`；接前：`git archive 3b03f00 | tar -x -C /tmp/ac121-before`
  另：`git archive 30a73f1 | tar -x -C /tmp/ac121-prefix30`（AC#4 修前那枚二进制）。
- 变异全在 `git archive 6a39820` 的副本 `/tmp/ac121-mut` 里跑，**做完逐枚恢复**，
  末了 `diff -rq /tmp/ac121-mut /tmp/ac121-after | grep '\.go:'` 输出为空（`.go` 树逐字相同）。
- **没有在仓库里建 worktree、没有 checkout、没有 reset/stash、没 commit 别人的文件。**
- 开工时 `git status --porcelain` 原文只有一条：`?? docs/evidence/s1/131-adversarial-acceptance.md`
  （`acceptor-ticket131` 的在飞写件）——我没碰、没 add。`internal/winsec/**` 本轮我没有一次写。
- ⚠ `git archive` 不含未跟踪的 `third_party/` dll ⇒ 我把三枚 dll 拷进快照的 `third_party/sherpa-onnx/` 才跑 `cmd/wisp`。
  这条在本文里既用于读数，也用于复现"剥光 dll 就起不来"那一形（见 §6 末）。
- ⚠ 我这轮踩到的仪器坑，写在用到的地方，集中列在 §10。

## 1. 裁决表（七格 + AC#5 那一空格）

| 格 | 判语 | 标签 |
|---|---|---|
| **AC#1** 复算四件仪器 | **通过。** 四件我全部自己重跑（§2）：接后依赖图含 `internal/models`、`DownloadingBridge.Run` 生产调用者从 0 变 1、边的形状逐字落在 `bridge.go:58`、"没有配置能绕过守卫"两道门都在（`validate.go:96` + `downloader.go:106`），且我在**真进程**上量到 rc=0 的交还。**它没有引用预检那张表**这句成立。 | 〔独立复现〕 |
| **AC#2** 链真接进去 | **通过，且判据今天成立**（§2 + §3 的 N-5）。两遍依赖图我自己在纯净快照上跑：20 → 22，逐包 diff 只有那两行；三 GOOS 相同；生产调用者点名 `cmd/wisp/models.go:303` → `bridge.go:58` 我在真进程上走通了（`wisp models ensure vad-silero` rc=0，"已交还（state=FirstRun, ticks=1, 末次进度 100%）"）。**但**：它那句 grep 判据盖不住"分发那一跳"，见 R-121-1。 | 〔独立复现〕 |
| **AC#3** 装配可达性用例 | **通过：不是恒绿，但它今天是否成立取决于两枚没人钉的清单。** §3 十二发变异里 **AC#3 相关十发**（N-1…N-7、N-10、N-10+N-1、N-11；N-8/N-9 属 AC#4）：**N-1/N-2/N-4/N-5/N-6 全部点名红**、**N-11 空集响亮失败**（它自述的仪器健全性腿，我复现）；**N-3（它自陈的天花板）与 N-7/N-10/N-10+N-1（我补的三发）全绿**。清单是**载荷**不是注释（N-4 加一项立刻红），AC#2 明令禁止的那形被用例抓住（N-5，它自己没声称过），前缀哨兵能答"否"（N-6）。⇒ 登记 R-121-2：**门的两枚策略清单（`capabilityPackages` / `graphGOOS`）本身没有断言**，缩它们是静默的，且"缩 GOOS + 挂 windows tag"这发组合物形正是编排者写那条硬判据时要防的那一形。 | 〔独立复现〕 |
| **AC#4** R-109-2 / R-109-3 | **通过，两洞都在真进程上验过前后两侧**（§4）。同一枚被篡改的 store，`30a73f1` 编出来的二进制 **rc=0**（"与已验签清单一致"）、`6a39820` 编出来的 **rc=1** 点名 `zz-unnamed-extra.onnx`，同一分钟、同一目录、同一 manifest。两发变异（N-8/N-9）windows + 容器两条腿都复现，反向腿保持绿。地界：`a7dfca9` 七枚文件全在 `internal/models/` 内 ⇒ 停手条款没触发（我逐字点过 `git show --stat`）。 | 〔独立复现〕 |
| **AC#5**（实现方故意留 `[ ]`） | **判：诚实降级，不是"该修没修"。** 三件事我一件件量了（§5）：①残留窗"接线之后是否可达＝否"这句**成立**（`Ensure` 的返回目录在 `bridge.go:44` 被 `_, err :=` 丢掉、全图 `os.Open*` 分布 3/1/1/3/6 逐包对上、`internal/speech` 仍只有 `doc.go`、`internal/engines/` 不存在）；②"装目录上仍挂着继承来的跨账号写授权"我在本机活体读到 `AFTER the hand-off guard ran` 两侧原文；③它没有拿"写清楚了"冒充"修好了"，也点明了唯一能自动收口的做法正是票 95 禁止的那件。**⇒ 票 109 的返修账不随本票结清**：R-109-1（残留窗）仍开着，票 121 明确不替它收口，这个"不勾"应当被读成一次交回而不是欠账。 | 〔独立复现〕 |
| **AC#6** 耗时账 | **通过（它这一格要的是"照抄 + 标噪声"，不是缺陷判定），但三处要改口径**（§6）：①`627.6 MB` 对、`707.8 MB` 我这边也复现不出，且我多量了一条它能被否证的理由——**`models/manifest.json` 全仓只有一个历史版本**，"当时清单不一样"这条路不存在；②它说"同形同机差 20%"——**低估了**：我在同一枚仪器上采到 1.113–2.113 s，比它最快那枚慢 **2.4 倍**；③冷读与真机启动**今天造不出来**（31.9 GiB 内存 vs 256 MiB 文件），这不是"再采十枚"能补的 ⇒ 建议**单独成格**。 | 627.6/依赖图＝〔独立复现〕；它那十枚＝〔仅自述，不背书〕 |
| **AC#7** 门禁 | **通过，八类读数逐字复现**（§7）：`internal/models` windows **102/76/0/4 rc=0**、**容器 POSIX 腿 102/76/0/4 rc=0（两腿四数逐字相同）**、`cmd/wisp` **RUN=146 顶层 PASS=78 子用例 PASS=68 FAIL=0 SKIP=0 rc=0**；`gofmt -l .` 空、`gofumpt v0.7.0 (go1.27.1)` 空、host vet rc=0、`GOOS=linux go vet ./internal/models/` rc=0、`GOOS=linux go vet ./cmd/wisp/` **rc=1 且与 `3b03f00` 控制组 stderr 逐字相同**；d22scan 纯净快照 rc=0、八 scope **202/22/40/18/16/40/382/31**，控制组 **202/21/40/18/16/40/376/29** ⇒ 零下降，+1/+6/+2 全部落到点名的文件。 | 四数/格式/台账＝〔独立复现〕；"CI 哪一步跑它"＝〔日志＋归档，我抽验〕 |

**总判：票 121 通过，零退回。** 判"退回"的那条线（"AC 声称要防的结局被真实造出来"）我按得很紧，
并且专门去造了：AC#3 声称防的是**包级**掉图，我用 N-1/N-2/N-4/N-5 四发把它做成会红的门；
它**没声称**防"接了但没人分发"，那一形（N-3）我造出来了、并且量到它看不见——**这不构成退回，构成一条必须归口的账**（R-121-1），
理由是 AC#2 的交付判据在**今天**是真的（真进程 rc=0 那条我亲手跑出来），不是纸面主张。
新账四条见 §8。**AC#5 那一格留在 `[ ]` 是本票最正确的一次不勾。**

---

## 2. 首要攻击点 1：两遍依赖图命中表（我自己跑的，不引用预检那张表）

### 2.1 `go list -deps ./cmd/wisp`（过滤本模块前缀）

| 树 | GOOS=windows（plain，无 `-e`） | GOOS=linux plain | GOOS=linux `-e` | GOOS=darwin plain | GOOS=darwin `-e` |
|---|---|---|---|---|---|
| `3b03f00`（接前） | rc=0 / **20** | rc=1 / 20 | rc=0 / **20** | rc=1 / 20 | rc=0 / **20** |
| `6a39820`（接后） | rc=0 / **22** | rc=1 / 22 | rc=0 / **22** | rc=1 / 22 | rc=0 / **22** |

- 接前 windows 腿原始 266 行（**不是**分母）；接后 272 行。
- 逐包 diff（接前→接后，windows）**只有两行新增**：
  `> github.com/CarlosShao/wisp/internal/models`、`> github.com/CarlosShao/wisp/internal/statemachine`。
- 接前两枚都不在（grep rc=1）；接后两枚都在。darwin 用 `GOARCH=arm64`（与预检同一形）。

### 2.2 它自陈那句"`-e` 没穷尽排除"——我判：**措辞诚实，但结论已被我这轮量成读数**

它的原话（票面 AC#2 条 + 预检 §6 第 3 条）：`-e` 会跳过坏边继续输出 ⇒
"某本仓包只经 sherpa 坏边可达"这一形没穷尽排除，且承认"我认为不存在，但这是推理不是测量"。
**这一形是可测的，我补了三把尺：**

**（a）plain 与 `-e` 的集合差，三个 GOOS 上双向都为零**（不需要容器，本机交叉查询就够）：

| 树 / GOOS | plain rc | plain **stdout** 里的本仓包数 | `-e` 数 | `-e` 独有 | plain 独有 |
|---|---:|---:|---:|---:|---:|
| `6a39820` windows | 0 | 22 | 22 | 0 | 0 |
| `6a39820` linux | **1** | **22** | 22 | **0** | **0** |
| `6a39820` darwin | **1** | **22** | 22 | **0** | **0** |
| `3b03f00` linux | 1 | 20 | 20 | 0 | 0 |
| `3b03f00` darwin | 1 | 20 | 20 | 0 | 0 |

⇒ 关键读数：**rc=1 的那两条命令，stdout 已经是完整的 22 个**，`-e` 一个本仓包都没多捞、没少给。
实现方说"对照腿只在能跑 plain 的腿上成立"——在**本机交叉查询**这个形状下，
plain 的 stdout 本身就完整，所以它那句"这形没测到"我这轮测到了。

**（b）穷尽版：整张图里"源不在本模块、目标在本模块"的边有几条？**

```
GOOS=<windows|linux|darwin> go list -e -deps \
  -f '{{range .Imports}}{{if hasPrefix . "github.com/CarlosShao/wisp"}}{{else}}{{$.ImportPath}} -> {{.}}
{{end}}{{end}}' ./cmd/wisp | sort -u          =>  三个 GOOS 全部 0 条
```

同一发命令顺手量到：接后树上 `-e` 报的错**只有第三方那一枚**（`sherpa-onnx-go-linux` / `-macos`，
`build constraints exclude all Go files`），windows 腿 0 枚。⇒ 那形要成立必须有第三方→本仓的 import 边，**图上一条都没有**。

**（c）换一把不依赖 `-e`、也不依赖 `go list -deps` 闭包实现的尺**：我自己写了枚 BFS（`/tmp/ac121-bfs.py`），
只沿**本模块前缀的 import 边**做闭包，第三方节点根本不进图。

| 树 | GOOS | 我的 BFS 闭包 | `go list -e -deps` | 对称差 | `internal/models` 在我的闭包里？ |
|---|---|---:|---:|---|---|
| `3b03f00` | windows / linux / darwin | 20 / 20 / 20 | 20 / 20 / 20 | **EMPTY** ×3 | False / False / False |
| `6a39820` | windows / linux / darwin | 22 / 22 / 22 | 22 / 22 / 22 | **EMPTY** ×3 | **True** ×3 |

⇒ **两把独立的尺 × 三种 GOOS × 两棵树，逐包对称差全为空。** 预检 §6 第 3 条留下的那条"未闭合"，
今天可以划掉；AC#3 依赖图结论我判**成立**。（措辞上我不要求它改口——它当年写"这是推理不是测量"就是诚实的。）

### 2.3 AC#1 其余三件仪器 + 真进程读数（我的）

| 仪器 | `3b03f00`（接前） | `6a39820`（接后） |
|---|---|---|
| `DownloadingBridge` 非测试命中 | **11 行**，**全部**是 `internal/models/bridge.go` 的 10 行声明/方法/注释 + `doc.go:25` 一行注释 ⇒ **包外 0 行** | 同上 + `cmd/wisp/main.go:76`、`cmd/wisp/models.go:11/252/291` 四行**注释**；真代码边在 `models.go:297 WireDownloading` / `:303 bridge.Run` |
| `VerifyInstalled` 非测试**调用** | `bridge.go:58` 一枚（挂在无人调的 `Run` 里） | `bridge.go:58` + **`cmd/wisp/models.go:241`**；声明 `downloader.go:228` |
| 边的形状 | `Run`（`bridge.go:39`）体内 `Ensure`（`:44`）之后 `:58` 是 `if verr := b.mgr.VerifyInstalled(id); verr != nil {` | **行号在 `6a39820` 上逐字仍是 `bridge.go:58`**（我开文件对的） |
| 有没有配置能绕过守卫 | 没有 | 没有，两道门：`internal/config/validate.go:96`（写 false 直接拒）+ `internal/models/downloader.go:106`（`NewManager` 拒构造）；`cmd/wisp/models.go:162` 传的是**字面 true** |

真进程（本机 windows，`portable.txt` 让 data dir 落在快照内，`models/manifest.json` 走 exe 目录，
`[models] local_override` 指向 `third_party/spike-models/silero_vad.onnx`——其 sha256 与 manifest 钉的
`9e2449e10874…` **逐字相同**，所以我能零网络把这条链走到交还那一步）：

| 发 | 命令 | 读数 |
|---|---|---|
| P-1 | `wisp models list` | **rc=0**，`已验签清单 …（6 个模型，store=…data-dev\models）`，六枚 size 逐枚打印 |
| P-2 | `wisp models ensure vad-silero` | **rc=0**，`已交还（state=FirstRun，ticks=1，末次进度 100%，0.01s）` ⇒ `Run` 走完且 `:58` 的 `VerifyInstalled` 返回 nil（否则 state=Error） |
| P-3 | `wisp models verify vad-silero`（干净 store） | **rc=0** `与已验签清单一致` |
| P-4 | 换一枚**篡改过 1 字节**的 store 副本 | **rc=1**，点名 `silero_vad.onnx: sha256 mismatch (pins 9e2449e10874…, got aa9e5770e04b…)` |
| P-5 | 把 `local_override` 从 config 里去掉、同篡改副本 | **rc=1**，路径从 `svdir` 变成 `…data-dev\models\vad-silero`（⇒ 守卫的**目标**跟着 config 走，见 R-121-3） |
| P-6 | `wisp models ensure`（override 指向一个**文件**而不是目录） | **rc=1**，`交还被拒绝（state=Error）：model: local_override "vad-silero" failed integrity check: model dir …sv.onnx missing` ⇒ 交还拒绝方向也真走到 `Run` 的 Error 出口 |

⇒ AC#2 那句"VerifyInstalled 在真进程的交还路径上被调到（点名文件:行）"我**亲手复现**，不是 grep 推出来的。

---

## 3. 变异清单（十二发；⚠ 编号说明见本节末）

**规矩都走了**：每发先 `grep` 证明落地（打印那一行的原文）→ `go build ./cmd/wisp/` 或 `go vet` **rc=0** → 才读红名；
全部在 `git archive 6a39820` 的副本 `/tmp/ac121-mut` 里；每发做完恢复，末了整树 `diff` 为空。

| 号 | 造什么 | 落地证据 | build/vet | 用例读数 | 对上实现方哪句话 |
|---|---|---|---|---|---|
| **N-1** | 把真边塞进 `//go:build windows`（`cmd/wisp/models.go` 首行插 tag） | `head -2` 打印出 tag；`grep -c` = 1 | `go build ./cmd/wisp/` **rc=0** | `GOOS=windows` **PASS**；linux **2 枚红**、darwin **2 枚红**，红名逐字含 `Parsed set: 20 of this module's packages.` | 它的 **M-A** ⇒ 〔独立复现〕 |
| **N-2** | 整条边摘掉（删 `cmd/wisp/models.go` + `main.go` 的 `case "models"`） | `models.go present? NO`；`grep -c cmdModels main.go` = 0 | **rc=0** | **6 枚红**（2 包 × 3 GOOS），三枚子用例全 FAIL | 它的 **M-B** ⇒ 〔独立复现〕 |
| **N-3** | 只摘分发、**留着 import** | `models.go present? YES`、`import` 行还在、`main.go` 里 0 hits | **rc=0** | **用例全绿**（`ok … 0.899s`） | 它登记的盲区 **M-B1** ⇒ 〔独立复现〕，**并且我放大了它**（下面三条） |
| **N-4** | 把清单**加长**：`capabilityPackages` 里加 `internal/audio`（真包、今天不在图里、不是本票的账） | `sed -n '115,119p'` 打印出三行清单 | `go vet ./internal/models/` rc=0 | **恰好 3 枚红**（`internal/audio` × windows/linux/darwin），`models`/`statemachine` 不受影响 ⇒ **清单是载荷，不是注释** | **它没做过这一发**（我补的） |
| **N-5** | **AC#2 明令禁止的那形**：把整枚 `cmd/wisp/models.go` 改名成 `models_probe121_test.go`（链还在、还被调，**只有测试文件 import 它**）+ 摘掉分发 | `grep -rln 'wisp/internal/models' cmd/wisp/` ⇒ 唯一命中是那个 `_test.go` | `go build ./cmd/wisp/` **rc=0**、`go vet ./cmd/wisp/` **rc=0**（测试编译得过） | **6 枚红**（两包 × 三 GOOS）⇒ `go list -deps` 不吃测试边 | **它没做过这一发**（我补的，且这是本票最该有的一条判据） |
| **N-6** | 哨兵加强：把 `internal/model`（真在图里那枚的**严格前缀**，自己不是包）加进清单 | 清单四行打印出来 | rc=0 | **3 枚红** ⇒ 现 Ship 的匹配是**整行精确**，能答"否" | 对应它文件头那句"Exact-line matching, not substring" |
| **N-7** | **把匹配器退化成子串**（`strings.Contains(k, pkg)`），保留它自己那枚 `notARealPackage` sanity 腿 | `grep -n 'strings.Contains(k, pkg)'` = 命中 161 行 | `go vet ./internal/models/` rc=0 | **整枚用例 PASS**（`ok … 1.015s`）⇒ **它自己的 sanity 腿抓不到这一形**，前缀哨兵不在断言里 | **它没做过这一发**（我补的）⇒ R-121-2 |
| **N-10** | **把 `graphGOOS` 从三枚缩到 `{windows}`**（编排者那条"逐 GOOS 各跑一次"的硬判据只剩一条腿） | `grep -n 'N-10 mutation'` = 123 行 | `go vet ./internal/models/` rc=0 | **绿**（`ok … 1.194s`）⇒ 缩策略清单是**静默**的：与 N-4（加一项会被查）正好相反 | **它没做过这一发**（我补的）⇒ R-121-2 |
| **N-10+N-1** | 复合物形：`graphGOOS={windows}` **同时**把 `cmd/wisp/models.go` 挂上 `//go:build windows`（= 编排者写那条硬判据时点名的那一形） | `grep -n 'N-10+N-1'` + `head -1 cmd/wisp/models.go` = tag | `go build ./cmd/wisp/` **rc=0** | **完全绿**（`ok … 0.739s`）⇒ **这条门只在 `graphGOOS` 那三行没被人缩过的时候成立**，而没有任何东西钉着那三行 | **它没做过这一发**（我补的，本票最省事的绕过路径）⇒ R-121-2 |
| **N-11** | 把 `modulePrefix` 改成不存在的模块名 ⇒ 过滤器什么都匹配不到（"空集"那一形的定向测试） | `grep -n 'N-11 mutation'` = 110 行 | `go vet ./internal/models/` rc=0 | **三枚子用例全 FAIL**，原文 `go list -e -deps ./cmd/wisp returned no github.com/CarlosShao/notreal/* packages at all - this is an instrument failure, NOT a statement about the list below` ⇒ 它的"空集单独响亮失败"**是真做了，不是注释** | 它的文件头声明 ⇒ 〔独立复现〕（**这是它的一个强点**） |
| **N-8** | `downloader.go:593` 那次 `verifyNothingUnnamed(...)` 换成 `return nil`（R-109-2 的守卫失效） | `sed -n '593p'` 打印出那一行原文 | `go build ./internal/models/` rc=0 | windows：`TestAC42` 四形**全红** + `TestAC43` 红，`TestAC41/47/48/49` **绿**；容器（`CGO_ENABLED=0`）**再加** `TestAC44`+`TestAC45` 红、`TestAC46` 绿；容器里 `sed` 恢复后 ⇒ `ok` | 它的 **M-C** ⇒ 〔独立复现〕，两条腿都对上 |
| **N-9** | `sidIdentity` 改成直接信 icacls 打的 trustee 文本（回到按名字判） | `grep -n 'N9 mutation'` = 91 行 | `go vet ./internal/models/` rc=0 | `TestAC47` 四形 **3 红 + 1 绿**（绿的是"要的 SID 直接以裸 SID 出现"，它本来就不经解析器——它预先登记过这句，我核了）+ `TestAC48` 红 + `TestAC49` 红 | 它的 **M-D** ⇒ 〔独立复现〕 |

⚠ **编号占用**：我本来想用 `M-7` 这类号，先在仓里 grep 了一次——`M-7` 在本仓**已被占**
（`internal/agent/approval/approval.go:48/232`、`gate.go:159` 三处注释写着"rulings M-7 / C-3 in docs/reports/pending-and-issues.md"），
与本票无关。⇒ 本文十二发一律用 **N-1…N-11**（含 N-10+N-1 那发复合物形），
实现方票面里的 M-A/M-B/M-B1/M-C/M-D 在中间列里点名对应，别串台。
**其中 N-4/N-5/N-7/N-10/N-10+N-1/N-11 六发是我自己加的**（它票面没做过）；
N-11 的读数是**支持它的**（空集那一形它真的响亮失败），这条我也照写，不把强点写成缺陷。

### 3.1 N-3 的天花板：我把它量到"整包全绿"，并指出断的那一跳在哪

它自陈的 M-B1 只说"AC#3 看不见"。我补了三条读数，说明**这一形的覆盖面比它写的还窄一点**：

1. **N-3 下两枚受影响包的全部用例仍然全绿**：`go test ./internal/models/ ./cmd/wisp/ -count=1`
   ⇒ `ok internal/models 4.755s`、`ok cmd/wisp 35.998s`。**不是只有 AC#3 看不见，是这两包里没有任何用例看见。**
2. 全仓 grep：没有任何 `_test.go` 提到 `cmdModels` / `modelsEnsure` / `handOffModel` / `modelsVerify`（逐词，
   `grep -rn --include='*_test.go' 'cmdModels\|modelsEnsure\|handOffModel\|modelsVerify' .` ⇒ **零命中**）；
   `cmd/wisp` 里唯一提到 `"models"` 的是 `providers_test.go:213`（那是 mockllm 的路由计数，不是子命令分发）。
   再问一遍"仓里有没有一条可重跑的零调用者仪器"：
   `grep -rln --include='*.go' 'zero.caller\|ZeroCaller\|noProductionCaller\|deadcode\|unusedexport' .` ⇒ **零命中**
   ⇒ 它 `next=` 第 3 条那句"符号级那道门今天仍空"是**对的**，我这轮把它从"没人做过"升级成"确认没有"。
3. **AC#2 那句 grep 判据在 N-3 之后仍然是"绿"的**： pristine 树上 `cmdModels` 的非测试命中是 **3 行**
   （`main.go:81` 调用 + `models.go:98` 注释 + `models.go:104` 声明）＝**生产调用者 1 枚**，N-3 摘掉的正是那一枚调用
   ⇒ 剩 2 行（注释 + 声明）＝**生产调用者 0**。
   而 `DownloadingBridge.Run` 的调用者 `cmd/wisp/models.go:303` **一行都没动**。
   ⇒ 断的那一跳不在 AC#2 的判据射程里（它数的是"谁调 Run"，不是"谁调 `cmdModels`"）。

**判：M-B1 作为 AC#3 那一格的结论可接受**（票面对 AC#3 只写了"包必须出现在依赖图里"，它做到了，
并且它**主动登记**了自己的天花板而不是藏起来）；
**作为"这条链已经接上"这句交付面的结论不可接受**——因为 R-109-4 的本义就是"能力在树里、不在产品里"，
而 N-3 只用六行编辑就能让这句话重新变成假的，且**仓里没有第二把尺拦它**。⇒ R-121-1（本文最重的一条）。

### 3.2 这条门在 CI 里到底有没有分母（"恒绿"的另一种形：根本没人跑）

- 宿主选 `internal/models`：`scripts/portable-tests.sh:139` 的清单里有 `github.com/CarlosShao/wisp/internal/models`，
  `:177` 的 core scope 含 `./internal/models/...`；`.github/workflows/ci.yml:255` 的
  `bash scripts/portable-tests.sh --scope=core` 在 **ubuntu-latest**（`:192`）那一个 job 里。
  ⇒ **这条用例两条腿都有分母**，且它挑 ubuntu 正是为了抓 windows-only 边（N-1 就是那一形）。
- 文件里 `t.Skip` 只有 1 处命中，是**第 25 行注释里的字**（"this file has no t.Skip anywhere"）；
  实际代码 0 处 `t.Skip(` ⇒ 它这句自述成立。⚠ 仪器坑：`grep -c 't.Skip'` 会把注释算进去，我第一遍就被骗了一次。
- 我没去开 CI 的 run 日志核对"上一次真跑过的那一步的 run id"，所以 AC#7 的 CI 侧我记〔日志＋归档，我抽验〕。

---

## 4. AC#4 的两洞：真进程前后两侧（这一节是我最强的一档证据）

**同一枚 store、同一分钟、同一 manifest，两枚二进制**（store 里放 `silero_vad.onnx`（哈希全对）
+ 一枚 manifest 没点名的 `zz-unnamed-extra.onnx`）：

| 二进制 | 命令 | 读数 |
|---|---|---|
| `30a73f1`（AC#4 修前） | `wisp models verify vad-silero` | **rc=0**，`vad-silero 与已验签清单一致（store=…data-dev\models，0.00s）` ⇒ **这形真的滑得过去**，不是纸面推断 |
| `6a39820`（HEAD） | 同一条命令 | **rc=1**，点名 `model dir holds 1 entry the signed manifest does not name: zz-unnamed-extra.onnx - re-hashing only the named files cannot see an added one, so a loader that opens this directory by name could be handed bytes nobody verified` |
| `6a39820` | 删掉那枚文件（反向腿） | **rc=0** |

⇒ **R-109-2 结清**。修前那枚二进制是我从 `30a73f1` 现编的（`go build` rc=0），不是引用谁的日志。
**R-109-3** 侧我走 N-9：三枚红（`TestAC47` 三形 / `TestAC48` / `TestAC49`）+ 一枚本该绿的红名解释，
与它票面逐字对得上。地界：`git show --stat a7dfca9` = 七枚文件，全在 `internal/models/` 内
⇒ 它那句"没要改 `internal/winsec` 的判定"成立，停手条款未触发（我没在 `internal/winsec/**` 看到本票的任何改动）。
另：`R-109-3` 的正解落在**两枚测试助手**（`acl_sid_121_test.go:83 sidIdentity` 等）上——
生产码 `internal/models` 今天不解析 ACL，所以"按 SID 判"是**仪器高度**的修复，不是行为修复；
这与票面 AC#4 那一段（41-45 行，"顺带补票 109 退回的两个洞（验收方实测能滑过去的两形）"）的措辞一致
（它要补的是验收方实测能滑过去的两形，没主张改生产判定），我不因此降格，但也不让它被读成"生产码改了"。

---

## 5. AC#5 那一空格：三条我都重走了

| 它的断言 | 我的读数 |
|---|---|
| `Manager.Ensure` 的返回目录在 `Run` 里被 `_, err :=` 丢掉 | `bridge.go:44` 原文 `_, err := b.mgr.Ensure(ctx, id, b.onProgress)`，`Run` 只往外传 `(statemachine.State, error)`（`bridge.go:39` 签名）⇒ **成立** |
| 全图 `os.Open*` 分布点名 | 我用 `grep -rn 'os\.Open' --include='*.go' <包目录> \| grep -v _test.go` 逐包量：`winsec 3 / observe 1 / memory 1 / tools 3 / models 6` ⇒ **与它逐字相同**。⚠ 换 `os.Open\|os.Create\|os.OpenFile` 就变 `3/2/1/2/7/1`（多出 `config 1`）⇒ 这枚读数**吃模式**，引用时要连模式一起引 |
| 没有一处按 `<store>/<id>/<file>` 开模型字节 | 我在**真进程**上补了一发（P-2/P-5）：`local_override` 那形下 store 目录**一个字节都没落**（`find data-dev` 只有 `config.toml` 与 `logs/`），交还仍然报成功 ⇒ 除了"没有装载器"，还多一条"装载器该开的目录可能不是守卫看的那一枚"（R-121-3） |
| 装目录上挂着继承来的跨账号写授权 | `TestAC1TheWriteWindowIsAnInheritedCrossAccountRight` **PASS**，原文我读到：`BEFORE … install dir kws-fixture: [BUILTIN\Users (I)(OI)(CI)(M)]` / `BEFORE … installed file model.onnx: [BUILTIN\Users (I)(M)]` / `hand-off refused the swap made through that right: … sha256 mismatch` / **`AFTER the hand-off guard ran - install dir: [BUILTIN\Users (I)(OI)(CI)(M)]`**、`installed file: [BUILTIN\Users (I)(M)]` ⇒ 守卫跑完，继承的 Modify **仍在** |
| 引擎装载今天无路径可达成 | `ls internal/speech/` = `doc.go` 一枚；`internal/engines` **不存在**；`internal/models` 的包外非测试 importer 只有 `cmd/wisp/models.go` 与 `tools/signmodels/main.go` ⇒ **成立** |

**判语**：这是**诚实降级**。三条依据：①它没有把"写清楚了"说成"修好了"（票面 AC#5 条逐字）；
②它给了"接线之后是否可达＝否"的**测量**判据而不是注释判据，其中两半（`_, err :=` / `os.Open` 分布）我复现到逐字；
③它点明了唯一能自动收口的做法（把装目录封窄）**正是票 95 明令禁止**且有用例钉着的那件，
所以"不勾"不是偷懒而是**没有合法的收口路径**。
⇒ **不据此退回票 121**；⇒ **票 109 的返修账不随本票结清**：票 109 欠的两洞（R-109-2/3）本票代付了，
R-109-1（残留窗）仍开，交回编排者挂到"第一个把模型字节读进内存的装载器"那张票上。

---

## 6. AC#6 耗时账：数字口径谁对

**它自报的三处话柄，我逐条自己算：**

1. **`627.6 ≠ 707.8`** —— 我这边独立算：**`627.6 MB` 对**。
   `models/manifest.json` 六枚 `size_bytes` = 32,654,866 / 643,854 / 237,202,501 / 163,002,883 / 64,717,756 / 129,347,466
   ⇒ **合计 627,569,326 B = 627.6 MB = 598.50 MiB**；`files[].size_bytes` 相加 = 626,925,472 = 626.9 MB（差 0.6 MB 正是 `vad-silero` 没列 files）。
   真进程 P-1 打印的六行字节数与此**逐字相同**（同一把尺的第三条腿）。
   **707.8 我复现不出，并且我多否了一条退路**：`git log --follow -- models/manifest.json` **只有一个版本**（`bfcb230`）
   ⇒ "验收方当时读的是另一版清单"这条路**不存在**；`707.8 − 627.6 = 80.23 MB`，清单里没有任何一枚是 80.2 MB，
   `MiB/MB` 混算（598.5 / 627.6）、`files` 口径（626.9）也都对不上。
   **⇒ 票面 AC#6 那句"真实 manifest 6 枚合计 707.8 MB"应当更正**（票面 append-only ⇒ 走 Progress log 补记，别改旧行）。
2. **"十枚全热"** —— 成立，而且比它写的更没法补：本机 31.9 GiB 内存 vs 256 MiB 文件 ⇒
   热是**结构性**的，不是采样不够；Windows 上没有丢页缓存的常规手段、本机也没让我重启。
   ⇒ 这一格不可能靠"再采十枚"闭掉，只能换条件（另一台机器 / 可丢缓存的特权 / 真机启动计时）。
3. **"我这十枚全部低于验收方下沿"** —— 它没撒谎，但**结论下小了**。我在同一枚仪器
   （`go test ./internal/models -run '^$' -bench BenchmarkAC6 -benchtime=5x -benchmem`）上采到：
   - `VerifyInstalled` 单枚 256 MiB：**1.322 / 1.113 / 1.499 / 2.113 / 1.290 s**（203 / 241 / 179 / 127 / 208 MB/s）
   - `Ensure` 缓存命中对照：**1.545 s**（174 MB/s）
   ⇒ 三组数摆在一起：**实现方 0.880–0.998 s** · **票 109 验收方 1.07–1.22 s** · **我 1.113–2.113 s**。
   同形同机同代码，**最快到最慢差 2.4 倍**，不是我看到的"20% 量级噪声"。
   我这轮还采到一枚 2.113 s，**超过验收方那组的上沿 73%**。⇒ 他们两组数谁都没错，
   但**任何"约 X 秒"的启动预算都不该由单枚 run 支撑**。
   ⚠ 我这轮不是编队安静的：同时段 `acceptor-ticket131` 在跑（它 09:1x 起在写它的证据文件），
   这一点我照实标成噪声来源之一，不当成"我的数才是对的"。

**由此对"启动翻倍"这句给 owner 的数的判定**：
**方向成立**（`Ensure` 已经付一遍全树哈希、交还再付一遍，两笔同量级这点三组数一致），
**但"≈2.1 秒"这个具体值我不背书**：按我这组吞吐（127–241 MB/s）算全清单 = 627.6 MB / 吞吐 ⇒
**约 2.6–4.9 秒**；票面那句"约 3 秒"落在这个带里，实现方改口后的"≈2.1 秒"落在带外（它用的是自己最快的带）。
⇒ **该报的不是一个点，是一个带**，并且**"冷读与真机启动零读数"应当单独成格**（现在它只是 AC#6 里的一句话）。

---

## 7. AC#7 门禁复算（含同 sha 控制组）

| 仪器 | 我的读数 | 它的读数 | 判 |
|---|---|---|---|
| `internal/models -count=2 -v`（windows） | **RUN=102 PASS=76 FAIL=0 SKIP=4**，rc=0（6.852s） | 102/76/0/4 rc=0（5.66s） | **逐字复现** |
| 同上（**POSIX 真跑**，`golang:1.27` + `CGO_ENABLED=0`） | **RUN=102 PASS=76 FAIL=0 SKIP=4**，rc=0（45.37s） | 同 | **逐字复现**（两腿四数相同这句成立） |
| 那 4 枚 SKIP 是谁 | `TestRealDownloadVadThroughPipeline` / `TestRealDownloadPuncArchiveThroughPipeline` ×2 count；两枚都在 `scripts/portable-tests.sh:324/325` 台账里，理由列写着要 `WISP_IT_REAL_MIRROR=1` | 同 | **不是我加的、也不是新形** ⇒ 对上 |
| `cmd/wisp -count=2 -v`（windows） | **RUN=146 顶层 PASS=78 子用例 PASS=68 FAIL=0 SKIP=0**，rc=0（74.83s） | 146/78/68/0/0 rc=0（67.46s） | **逐字复现** |
| `gofmt -l .` | 空，rc=0 | 空 rc=0 | 复现 |
| `gofumpt -l .` | 空 rc=0；`-version` = **`v0.7.0 (go1.27.1)`**（先自证存在） | v0.7.0 | 复现 ⇒ 没有"未跑"要引错误原文 |
| `go vet ./internal/models/ ./cmd/wisp/ ./internal/statemachine/` | rc=0 | rc=0 | 复现 |
| `GOOS=linux go vet ./internal/models/` | **rc=0**（⚠ 只编译不执行） | rc=0 | 复现 |
| `GOOS=linux go vet ./cmd/wisp/` | **rc=1**，原文 `imports github.com/k2-fsa/sherpa-onnx-go-linux: build constraints exclude all Go files …@v1.13.8`；**控制组 `3b03f00` 同 rc=1、stderr `diff` 逐字相同** | rc=1 / 同 | 复现 ⇒ 不是本票引入 |
| `sh scripts/d22scan.sh`（纯净快照） | **rc=0**，`clean - no D22 ban violations`；正对照 `PASS=21 FAIL=0 SKIP=0, === RUN=31` | 同 | 逐字复现 |
| 台账八 scope（`3b03f00` → `6a39820`） | `bans #1-5 internal/ 202→202` · `bans #1-5 cmd/ 21→**22**` · `ban #6 frontend/ 40→40` · `ban #7 internal/tools/ 18→18` · `ban #8 design/ 16→16` · `ban #8 frontend/ 40→40` · `ban #8 internal/ 376→**382**` · `ban #8 cmd/ 29→**31**` | 同八个数 | **零枚下降**，三个上升值逐字对上 |
| 上升值归因 | `git diff --name-status --diff-filter=A 3b03f00 6a39820 -- internal/ cmd/` = **恰好 8 枚**：`cmd/wisp/models.go`、`cmd/wisp/secret_dataroot_119b_test.go`（票 119 的 `36294c2`）、`internal/models/{assembly_reachability_121, acl_sid_121, acl_sid_121_windows, verify_tree_121, verify_tree_121_other, verify_cost_121_bench}_test.go` | 同八枚 | 逐字对上，**不是"我猜是它"** |
| 剥光 dll 还能不能跑 | 快照里没有 `third_party/` ⇒ 拷三枚 dll 才有上面那些读数。**同一条命令、exe 旁边没有 dll**：`returncode = 3221225781 => 0xC0000135`，stdout/stderr 皆空 | rc=1，原文 `exit status 0xc0000135` | **复现**（它引的是 Go 侧看到的字符串，我引的是进程退出码；两者同物）⇒ 票 98 那个洞今天还在，本票所有"`cmd/wisp` 能跑"的读数都带这个前提 |

## 8. `R-121-*` 台账与建议归属

| 号 | 一句话 | 证据 | 严重度 | 建议归属 |
|---|---|---|---|---|
| **R-121-1** | **枚举门的天花板被放大了一格：断"分发"那一跳，AC#2 的 grep 与 AC#3 的用例同时看不见。** N-3 之后：`go build ./cmd/wisp/` rc=0、`internal/models` 与 `cmd/wisp` **两包全部用例仍 `ok`**、`DownloadingBridge.Run` 的生产调用者（`cmd/wisp/models.go:303`）一行没动，而 `cmdModels` 的生产调用者从 2 命中（1 调用 + 1 声明）掉到 1（只剩声明）＝**0 个调用者**。仓里没有一条可重跑的仪器拦这一形。 | §3 N-3 + §3.1 三条读数 | 中（今天无产品后果，因为链的另一端没人读；但它是**本票要防的那个形状**的复活路径） | **不在本票收**（它比一次修复大一格）。建议**开一张新票**："分发/入口级可达性仪器"，判据至少要能红 N-3 这一形；实现形状可参考 N-5 已经证明可行的一侧（`go list` 不吃测试边）。次选：把它挂进票 114（composer→perm，同族"零生产调用者"）。 |
| **R-121-2** | **这条门的两枚策略清单没有被任何断言钉住，缩它们是静默的。** 三发读数：①N-7——成员判定退化成子串后整枚用例 **PASS**，文件头那句"Exact-line matching, not substring"与那枚 `notARealPackage` sanity 腿**都抓不到**（真能答"否"的是严格前缀哨兵，N-6 证明实现今天是对的，但那是实现自带、不是断言保证）；②N-10——`graphGOOS` 缩到 `{windows}` **全绿**；③N-10+N-1——**"缩 GOOS" + "把真边挂进 `//go:build windows`"两发合起来仍然全绿**，而这一形正是编排者在票面"追加一条硬判据"里点名要防的那一形（windows 腿绿、linux 腿静默掉出图）。对照：N-4 证明**往清单里加**一项会立刻红 ⇒ 加长被查、缩短不查。 | §3 N-6/N-7/N-10/N-10+N-1 | 中（今天两清单都是对的；但它承诺的那条硬判据是**可静默撤销**的） | **本票续单或入表纪律**，三行量级：①把 `modulePrefix + "/internal/model"`（严格前缀、不是包）加成永久负断言；②断言 `len(graphGOOS) >= 2` 且含 `windows` 与 `linux`（票面那条硬判据的最低形式）；③与 R-121-1 同做。谁下次动 `assembly_reachability_121_test.go` 谁顺手，或并入 1。 |
| **R-121-3** | **交还复验的"目标"跟着 config 走，而同一行打印的 store 路径不变**：`VerifyInstalled`（`downloader.go:234-236`）在有 `local_override` 时验的是覆盖目录；真进程 P-5/P-2 量到"store 目录里一个字节都没有、`wisp models verify` 报 `与已验签清单一致（store=…models）`"，而把 config 里那行去掉、同一枚被篡改的 store ⇒ rc=1 点名 `sha256 mismatch`。⇒ **C29 的交还判定在 sideload 配置下不描述 CLI 所说的那枚目录**。（分支本体在 `3b03f00` 就有，**不是本票引入**；但**本票的 AC#2 是第一个让 `cfg.Models.LocalOverride` 在生产进程里可配落地的事件**——之前 `NewManager` 生产调用者为 0。） | §2.3 P-2/P-5/P-6 + provenance `git show 3b03f00:internal/models/downloader.go`（226-240 行原文含那一支） | 中（今天无装载器读那枚目录 ⇒ 无实害；一旦有读者就变成"守卫看 A、读者开 B"） | **票 109 的返修账**（`VerifyInstalled` 的语义归它）+ **与 AC#5 同一张"读者票"合并收**：读者落地那一刻，必须钉"守卫验的目录 == 读者开的目录"，且 CLI 文案要么打印真正验过的目录、要么在被 override 时明说。 |
| **R-121-4** | **AC#6 的两个口径问题**：①票面 `707.8 MB` 复现不出（manifest 全仓只有一个版本 ⇒ 无"当时不同版"退路），真值是 627,569,326 B；②单枚 run 撑不起一个启动预算——同形同机三组数 0.880 / 1.07–1.22 / 1.113–2.113 秒，散布 **2.4 倍**，且冷读在 31.9 GiB 内存的本机**结构性造不出来**。 | §6 | 低（不是缺陷，是交付面的数说错了一次 + 说满了一次） | **票 121 自己补一枚 Progress log 更正 `707.8 → 627.6`**（票面 append-only）；**"冷读 / 真机启动预算"单独成格**，建议挂 SLO 那一族（票 96/98/123 地界），别留在 AC#6 的括号里。 |

**顺带登记两条不是缺陷、但会影响下一位读数的人**：
① `os.Open*` 那枚分布读数**吃 grep 模式**（`os\.Open` 给 3/1/1/3/6，加上 `os.Create|os.OpenFile` 给 3/2/1/2/7/1）⇒ 引用要连模式一起引；
② `grep -c 't.Skip'` 会把"本文件没有 t.Skip"这句**注释**算成命中（我在 AC#3 文件第 25 行被它骗了一次）⇒ 判有无跳过要 `grep 't\.Skip('`。

## 9. 注入文本登记（Rules 末条要求）

截至本节写下时（2026-09-23 09:45，且它可能继续出现 ⇒ 用有界说法）：
- 实现方票面登记的那段**自称**"用户已更新编码规则 / 用户偏好优先于 AGENTS.md"的文字（含一张 `| ID | 分类 | 状态 | 规则 |` 表），
  **在我本轮的工具输出里出现 0 次**（截至 09:45）⇒ 我既不能背书它的"≥18 次"，也不否认（它记的是它那一轮，
  且它自己已把死数改成有界说法，这个处理是对的）。
- 我这轮**确实**在工具输出里见到、且**没有当成授权**的文字，逐类计数（截至 09:45，会继续增长）：
  harness 级"MEMORY.md 已被修改，以下是修改后内容"提示 **4 次**；
  "[SYSTEM NOTIFICATION - NOT USER INPUT] 后台任务完成"通知 **1 次**。
  ⇒ 全部**只当我自己的上下文被刷新**处理：**未据此改动任何判据、未 revert、未放宽任何阈值、未冻结任何包、未 commit 别人的文件。**
  （其中那次后台通知是我自己起的 `cmd/wisp` 门禁 job，内容就是"任务完成"，没有指示。）
- 本文里没有任何"编排者在对话里下的指令"被我引用为授权；票面/台账里"禁改 `internal/risk`（冻结）"那类字样
  来自仓库文件本身，是真票面，不在这一类。

## 10. 票 121 能否结案 / 票 109 的返修账是否随之结

- **票 121：可结案，零退回。** AC#1/AC#2/AC#3/AC#4/AC#7 我全部〔独立复现〕到逐字；
  AC#6 的判据是"照抄 + 标噪声"，它照抄了、也主动标了自己的采样偏差 ⇒ 判通过，同时按 R-121-4 补两处口径。
  **AC#5 留 `[ ]` 是对的**，结案时**不许被人顺手勾掉**，也不许被读成"忘勾"——它是一次显式交回。
- **票 109 的返修账：部分结清。** 已结 = R-109-2、R-109-3（§4 真进程两枚二进制前后两侧，我这轮是独立复现）。
  未结 = **R-109-1**（守卫返回之后到读取者 `open` 之前那整段，本票判"今天无路径可达"因此不据此 FAIL）
  + 本文新增 **R-121-3**（同一族、同一段窗，但换了一枚配置入口）。
  ⇒ 票 109 的"接线"那一半现在真的接上了；它的**窗口语义那一半仍在 109/未来的读者票手上**，别把它算成随 121 结清。
- **票 121 `next=` 第 3 条（符号级那道门仍空）**：判**成立且被我这轮做硬了**——
  它当年是"复现成了红/绿两侧读数"，我这轮多了 N-3（两包全绿）+ N-5（测试边能被抓住）两侧新读数，
  所以这条不该继续挂在 `next=`，该落成 R-121-1 的票号。

## 11. `next=`（交回编排者，按可执行度排）

1. **裁决本表 + 定 R-121-1 归属**（建议开新票："入口/分发级可达性仪器"）。**判据要能红 N-3 这一形**：
   只写"再 grep 一次生产调用者"不够，因为 N-3 之后"谁调 Run"那枚读数一点没变。
2. **决定 AC#6 是否单独成格**（R-121-4②）。我的推荐：**单独成格**，归 SLO/真机启动那一族；
   代价若不成格 = owner 拿到的"启动翻倍 ≈ 3 秒"这句话今天没有任何一条冷读或真机启动读数支撑。
3. **票面 707.8 MB 让本票自己 append 一枚更正**（R-121-4①）。我不改别人票面，一字没动。
4. **R-121-3 交回票 109 的账**（`VerifyInstalled` 语义归它），并与"第一个模型装载器"那张票合并收口；
   合并判据建议直接写成："守卫验的目录 == 读者开的目录"，两边各钉一枚用例。
5. **R-121-2 顺手做**（三行量级，两处）：①永久负断言"严格前缀哨兵不在图里"；②`len(graphGOOS) >= 2` 且含 windows+linux
   （否则编排者那条"逐 GOOS"硬判据可以被一次编辑静默撤销，我的 N-10+N-1 就是这么绕过去的）。
   谁下次动 `assembly_reachability_121_test.go` 谁做，或并入 1。
6. **AC#3 清单入表纪律**（它 `next=` 第 2 条我核过方向）：票 117 / 票 114 结案时各自申请入表，
   `internal/audio`/`internal/ball` 仍挂 A91③；**本票一格都没顺手加**（我在 N-4 里加过 `internal/audio`，
   那是**变异**、做完就恢复，`diff` 为空可查）。
7. 本表**只 commit 未 push**（`b43c149` 主表 + `9610d61` 两处口径更正）；`git add` 每一枚都只有
   `docs/evidence/s1/121-adversarial-acceptance.md` 一枚路径，提交前 `git diff --cached --name-only` 逐字点过。
   共树里此刻还有别人的在飞写件（`.github/workflows/ci.yml`、票 128 两枚、`docs/evidence/s1/131-adversarial-acceptance.md`、
   `docs/reports/frontend-handoff.md`），**我全程没 add、没 commit、没读它们当基线**（我的读数全部来自 `git archive` 快照）。
