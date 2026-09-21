# 105 — 独立对抗验收裁决表（acceptor-ticket105）

验收方：`acceptor-ticket105`（只读；一行生产码未改，本文件是唯一工作树写件之一，另加 `/tmp` 快照）
被验：票 105「C26 改写账只有测试在读 → 接一条真实生产消费路径 + 补祖先重解析行为用例」
HEAD 复算时点：`f6818f24ec404063fc6dde408edbe5b56acdf2ac`（分支 `dev`），本文件起笔 `date` = **2026-09-21 20:25 CST**
交件自述：票面 `## Progress log` 末条（19:5x，agent-ticket105），其 `next=` 请出本表。

## 裁决表（1 AC = 1 格）

| AC | 票面判据要点 | 裁决 | 证据标签 |
|----|--------------|------|----------|
| AC#1 | 复算"零生产读取者"，贴真实命中表，不许照抄 | **通过（复现＝成立）** | 〔独立复现〕 |
| AC#2 | 至少一条真实生产消费路径（grep 非零 + 装配根可达 + 端到端断记录内容） | **通过** | 〔独立复现〕 |
| AC#3 | 祖先重解析行为用例断可观察结果 + 变异红在行为；独立复算其可达性结论 | **不通过（FAIL）** — 可达性结论被反例驳倒，那半仍可被行为钉而未钉 | 〔独立复现〕 |
| AC#4 | 若"给人看"必须动 frontend/ 或 internal/panel/ ⇒ 停手登记 | **通过（未扩界）** — 4 枚 commit 并集只有票面+两枚测试+`bridge.go` | 〔独立复现〕 |
| AC#5 | 门禁：go test/gofmt/gofumpt/vet/GOOS=linux vet/d22scan 纯净快照 rc=0，台账不降 | **通过** — 分包四数复现，d22scan 三枚树 rc=0、`internal/` 363→370 升；两处交件数字口径差异已点名 | 〔独立复现〕 |

**接续说明（重要，避免后人误读）**：本文件上半程由 `acceptor-ticket105` 写到 AC#2 的裁决行即被模型连接中断杀掉（3317 字节，20:29 停笔）。
AC#1、AC#2 两格是**它**的独立复算，我没有重跑（只抽验 AC#2 最要紧的一刀，见下）；AC#3/AC#4/AC#5 与"总判"由接续代理 `acceptor-ticket105b` 落盘，
`date` 读数：起笔 20:31 CST、AC#3 落盘 20:4x、AC#4/#5 落盘 20:5x CST。它的快照 `/tmp/ac105-acc105` 我核对过＝`f6818f2` 纯净树（`diff -r` 与我自己解的一份**零差异**），我的复算一律跑在自己的 `/tmp/ac105b-acc105b`（探针与变异）与 `/tmp/ac105b-gates`（门禁，纯净无探针）里，**未读当前脏工作树冒充被验版本**。


---

## AC#1 — 复算"零生产读取者"

**判据**：票面 AC#1 原文命令 `grep -rn "RewrittenRoots(\|UnusableRoots(\|\.Roots(" --include=*.go . | grep -v _test.go`，贴真实命中；若其实有生产读取者 ⇒ 写"未复现"并降级本票。

**我的复算（当前 HEAD f6818f2，工作树含交件的三枚 commit）**：
- 非测试命中 **5 条**：
  - `internal/tools/bridge.go:810` `// account. Until here the only consumers of Roots()/RewrittenRoots()/` —— **注释**
  - `internal/tools/bridge.go:811` `// UnusableRoots() were tests, ...` —— **注释**
  - `internal/tools/bridge.go:864` `return b.paths.Roots(), b.paths.RewrittenRoots(), b.paths.UnusableRoots()` —— **真调用点（AC#2 新增的消费者）**
  - `internal/tools/paths.go:186` `func (p *PathCanonicalizer) UnusableRoots() []string {` —— **定义**
  - `internal/tools/paths.go:194` `func (p *PathCanonicalizer) RewrittenRoots() []string {` —— **定义**
- 含测试总命中 **35**，非测试 **5**（口径：全仓 `--include=*.go`，`grep -v _test.go`；不带 `-v`）。

**对票面 AC#1"零生产读取者"这句话的判定**：票面 AC#1 要复算的是**建票时**的状态。交件的 19:3x 条自述在**接线前**量到 2 条命中且全是定义行、`Roots()` 分支零命中 ⇒ 那才是"零生产读取者"的原始命题，本票据此未降级。我在**接线后**HEAD 复算，非测试命中变 5，但其中 3 条是定义/注释、**唯一新增的真实读取者就是 AC#2 要造的那个**（`bridge.go:864`）。⇒ 与自述一致：**接线前该命题复现成立**；接线后它由本票 AC#2 亲口改掉。

**引用腐坏检查**：自述 AC#1 引用的 `paths.go:186/194` 定义行、`bridge.go:864` 调用点，我逐条对上，无行号漂移。

**AC#1 裁决：复现＝成立（未复现的反例不存在）**〔独立复现：命令与命中我亲手跑过〕。本票不降级，AC#2/AC#3 照验。

---

## AC#2 抽验（接续代理只攻最要紧的一刀，不重做整格）

交件说消费者落在 `internal/tools/bridge.go:864` 的 `pathAccount()`。我从**装配根**往下自己走了一遍（快照 `f6818f2`）：

`cmd/wisp/run.go:265 rt.paths = tools.NewPathCanonicalizer(...)` → `run.go:335 rt.bridge = tools.New(tools.Options{` → `run.go:337 Paths: rt.paths` → `run.go:353 Logf: rt.auditf` →
`internal/agent/loop.go:727 return l.opt.Tools.Execute(ctx, req)`（生产环路，不是测试）→ `internal/tools/bridge.go:243 func (b *Bridge) Execute` → `:271 dec.Paths = b.displayPaths(rawPaths)` →
`bridge.go:756 b.book(...)`（reject）/ `bridge.go:770 b.book(...)`（close）→ `bridge.go:820 if len(dec.Paths) > 0` → `:821 b.pathAccount()` → `:864 return b.paths.Roots(), b.paths.RewrittenRoots(), b.paths.UnusableRoots()` → `bridge.go:911 b.log()` → `b.logf` = `rt.auditf`。

- **生产调用者计数非零**：`book()` 的两个调用点都在 `Execute` 的收口上（9 处 `b.close(`/`b.reject(`，全在 `bridge.go` 生产码），每一笔工具调用必过 ⇒ 不是"只有测试构造调到"的形状。第 8 类缺陷这一刀**没有命中**。
- 行号对得上：`pathAccount()` 定义 `:860`、真读取点 `:864`、`book()` 里的发射块 `:820-825`，与票面 AC#1/AC#2 自述一致，无漂移。
- ⚠ 一处口径必须写清（不影响 AC#2 判定，但影响对 owner 的那句话）：`rt.auditf` 的实现是 `cmd/wisp/run.go:428 fmt.Fprintf(rt.stderr, "[audit] "+format+"\n", ...)`，
  `rt.stderr` 默认 `os.Stderr`（`run.go:135-136`、`main.go:99`）⇒ 这条 `PATH-ACCOUNT` 记录在生产里落的是**进程 stderr 流**，不是磁盘文件。
  票面 AC#2 要求的正是"与 `MODE-READ` 同族的那种结构化记录"，而 `MODE-READ`（`run.go:332`）走的是同一个 sink ⇒ **同族成立**；
  "盘上审计里出现这条记录"那一半是**测试自己把 logf 接到真文件**上验的（`bridge_path_account_ticket105_test.go`），交件没把它说成生产落盘，我核对过措辞。
  ⇒ 登记 `R-105-1`（要真"看得见"得有人把 stderr 收进文件或把这本账接进面板，属票 92 地界），不阻塞本票。

**AC#2 抽验结论：维持通过**〔独立复现：装配根链路我逐行走过；未重跑 AC#2 的三条用例与 M1〕。

---

## AC#3 — 祖先重解析那条腿的行为用例（本票最可能 FAIL 的一格）

**判据原文**：构造"祖先带 reparse/junction"的输入 ⇒ 断**可观察结果**（被拒 / 被记账 / 不落到错树），然后**变异**：去掉那条腿的判定 ⇒ 这条新用例必须红（红在行为，不是红在符号）。
⚠ 票面只写"真机 junction 用 `mklink /J` 在临时目录造"，交件据此造了两条用例（`internal/risk/syncdirs_ancestor_reparse_ticket105_windows_test.go`），
并在进度里请我独立复算它这条可达性结论：**"只删 `ares.Actable()`（保留重解析）在真机上触发不了——祖先串是从 C26 自己的输出切出来的，不含未展开的构造"**。

### ① 我自己重走的可达性论证（未照抄；快照 `f6818f2`）

那条腿＝`internal/risk/syncdirs.go:222-232`：`Resolve(anc, s.excs)` 的 deny 分支 + `ares.Actable()` + 空值守卫。走到它的前提（`:197-221`）：
第一次 `Resolve(raw)` 干净（`res.Actable()` 不报错 ⇒ `Rewritten=false`）、`res.Resolved=false`（叶子不存在）、`deepestExistingAncestor(canon)` 交出非空 `anc`。

它这句论证的**前半是对的**：`deepestExistingAncestor`（`:274-298`）不重新拼写，它把 `splitPathComponents(canon)` 的组件逐个 `Lstat` 之后 `TrimSuffix` 回来 ⇒ `anc` 确是 `canon` 的**字面前缀**。
但**结论错了**，错在它没看的这一步：`canon` 是 `lexCanonical`（`filepath.Abs`+`Clean`，`:206-214`）**之后**的串，而第一次 `Resolve` 的 `expandAccounted`（`:159-202`）扫的是 **RAW**。
`Clean` 会把重复分隔符 `\\` 折成 `\`，于是 **`%…%` 两个百分号的配对成员会变**：raw 里配对读到的名字含**两个**分隔符，`canon`（`anc` 从它切出来）里配对读到的名字含**一个**。
`Rewritten=false` 只要求 raw 那一次没换到值，并不要求清洗后的串也换不到。⇒ **存在**"第一次 `Actable()` 放行、第二次 `Actable()` 报错"的输入形状，那条腿不是死码。

### ② 探针（两个，都是我自己造的，跑在 `/tmp/ac105b-acc105b`）

`zz_probe_105b_test.go`（纯函数级，判据＝配对是否真的不对称）——`--- FAIL` 是**预期的"反例成立"读数**：

```
raw     = "C:\\%B\\\\C%\\notes\\new.md"
cleaned = "C:\\%B\\C%\\notes\\new.md"
expandAccounted(RAW)     -> "C:\\%B\\\\C%\\notes\\new.md" kinds=[]      (Rewritten=false)
expandAccounted(CLEANED) -> "C:\\C:\\somewhere\\else\\notes\\new.md" kinds=[env] (Rewritten=true)
Resolve(raw)     err=<nil> Rewritten=false Resolved=false
Resolve(anc)     err=<nil> Rewritten=true
second Result.Actable() err=risk: expansion moved this path onto a tree the caller did not name: ... (env)
```

`zz_probe_105b_reachable_windows_test.go`（**行为级**，走真 `Provenance.IsSyncPath`：临时目录里建真目录 `%X`、`%X\Y%`，`os.Setenv("X\\Y", "..\\escape")`，输入 `…\%X\\Y%\new.md`，`escape` 是同步树外一个**真存在**的目录——专门为了排除"`!ares.Resolved` 守卫顺手救场"这个解释）：

| 跑法 | 我的探针 | 票面 AC#3 两条用例 |
|------|----------|--------------------|
| 纯净 `f6818f2` | **PASS**：`Sync=true Root={sync-suspect/suspect-fallback}`，`Why` 末尾点名 `risk: expansion moved this path …` | 2/2 PASS（复现自述"修前 2/2"） |
| **M3＝只删 `ares.Actable()`，保留重解析**（`:226` 改 `anceCanon := ares.Canonical`） | **FAIL**：`Sync:false`、`why="write target is not under any sync root"` | **2/2 仍 PASS** |
| M2＝整段 `:222-232` 换成 `ares := Result{Resolved:true}` + `anceCanon := anc`（复现交件的变异） | FAIL（另一种形状：命中的是错根 `Ticket105bProbe` 而非 fail-closed） | **用例1 FAIL**（`Sync:false`，红在判决，原文我复现到字节）、用例2 PASS |

⇒ **交件的 M2 我复现＝成立**；但**它的可达性结论被我的反例驳倒**：`M3` 下我的探针**红了**，而那笔落进同步树（`…\%X\\Y%\new.md` 字面在注册根 `…\001` 之下）的写被判成"不在任何同步树里"＝**被放行**——正是 AC#3 声称要防的结局。

限制我照实写：这条反例的三个前提里，"目录名含 `%`"和"路径里出现 `\\`"都是模型能直接供给的；第三个前提**不是**——它要 wisp 进程的环境里存在一个**名字含分隔符**的环境变量（Windows 允许，`os.Setenv("X\\Y",…)` 实测成功）。⇒ 真机可利用性中等偏弱，但"触发不了"这句**绝对**断言不成立。
同一条 `\\`→`\` 折叠还意味着 `:223-225` 那个 deny 分支也不是不可能到达（展开后的祖先自己经过一个未登记 junction 时）——这条**我只是读码推出来的，没有造出探针**，按第三档处理。

### ③ 一个只能被符号名满足的用例算不算达成 AC#3——我的判

不算。理由三条，按分量排：
1. 本票立票的原因（R-102-3）就是"第 4 行只钉符号、没有行为用例"。交件把**重解析**那半钉成了行为（M2 红，我复现），但**`Actable` 那半今天仍然只被符号钉着**（`pathresolver_rewrite_account_test.go` 的静态"字段被读"判据），而这半**可被行为钉**——我 40 行就造出来了，不需要特权、不需要真机 junction。⇒ 立票要消灭的那个形状在交付物里**还剩一处**。
2. 票面 AC#3 的变异措辞是"去掉那条腿的**判定**"。那条腿里唯一带"判定"性质的两处就是 deny 与 `Actable()`；把整段连"重解析"一起删掉是最粗的一发，粗到 `anceCanon := anc` 必然改变字面拼写、必然红。红在最粗的一发上，不足以证明细的一发被钉住——这正是"红在符号"与"红在行为"要区分的原因。
3. 用例 2（非 exempt ⇒ fail-closed）在任何针对这条腿的变异下都拿不到红（非 exempt 祖先在**第一次** `Resolve` 就被 `reparseComponents` 拒了，够不到 `:222-232`）——它断的是**另一条腿**的行为，被本票当成第二条行为用例计数，属"断言恒真"这一类假绿的邻近形状（它不是假绿：判定本身是有意义的安全断言，但它不背 AC#3 的变异账）。交件对这一点是诚实的（它自己写了"用例 2 仍 PASS，如实记录原因"），我不据此加罪。

**AC#3 裁决：不通过（FAIL-退回本格）**〔独立复现：两个探针、M2/M3 两发变异、纯净基线全是我亲手跑的，命令与读数在下表〕。
可执行的补法（不改票 102 语义）：把 `resolveTarget` 的祖先 `Actable()` 那半做成一条**行为**用例（我的 `zz_probe_105b_reachable_windows_test.go` 可直接搬），或按我自己都不完全信服的路子——把那半明确降级成"由静态判据承载 + 票面写明它今天只在 env-name 含分隔符时可达"。

---

## AC#4 — 边界（有没有撞票 92 的地界）

**判据原文**：若 AC#2 的"给人看"这一半**必须**动 `frontend/` 或 `internal/panel/` ⇒ 停手登记交回，不许为了勾框扩界。

**先把 105 的 commit 自己找出来**（不采信票面转述）：`git log --oneline --all -- internal/tools/bridge.go` 最新一枚＝`21da8f3`；
`git log --oneline --all --grep=105` 给出整族 **4 枚**：`0d2f499`(19:41 docs AC#1) / `2afd064`(19:46 test 修前红) / `21da8f3`(19:48 feat AC#2) / `6a7dcc9`(19:56 docs 交件)。
`4693feb`、`4d43447` 我逐枚 `git log -1 --format=%s` 核过，两条都写 `test(104…)` / `fix(104…)` ⇒ 属票 104，与本票无关。

**逐枚点数（`git show --name-only`，命令原文与真实清单）**：

| commit | 文件数 | 清单 |
|--------|--------|------|
| `0d2f499` | 1 | `.scratch/wisp/issues/105-…-reader.md` |
| `2afd064` | 2 | `internal/risk/syncdirs_ancestor_reparse_ticket105_windows_test.go`(A)、`internal/tools/bridge_path_account_ticket105_test.go`(A) |
| `21da8f3` | 2 | `internal/tools/bridge.go`、`internal/tools/bridge_path_account_ticket105_test.go` |
| `6a7dcc9` | 1 | `.scratch/wisp/issues/105-…-reader.md` |

并集＝**4 个不同文件**：票面 ×1、新测试 ×2、`internal/tools/bridge.go` ×1。逐条 grep 越界词（`frontend/|internal/panel/|internal/winsec/|paths\.go|assessor|pathresolver|allowlist|PLAN|specs`）在这 4 枚的 `--name-only` 上：`NONE (零越界)`。

**⚠ 一条仪器口径，我主动写清**：我一开始想用区间 diff（`git diff --name-only 2afd064^ f6818f2 -- internal/panel frontend internal/winsec …`）来判边界，它吐出 `internal/panel/composer.go`、`internal/winsec/winsec_windows.go` 等 8 条——**那不是 105 改的**，是共树里别的票在 `2afd064^` 与 `f6818f2` 之间落进来的 commit。⇒ 判定归属只能用逐枚 `git show --name-only`（上表），区间 diff 在共享工作树上会栽成"冤案"。这条我差点报错，先登记。

**禁改语义核对**：`internal/tools/paths.go`（票 102 产生端）**不在任何一枚的清单里**；`internal/risk/syncdirs.go` 同样不在（交件承诺"不会为了勾框改 syncdirs.go"，兑现）。
"给人看"这一半走的是 `bridge.book()` 的审计行（AC#2），面板显示仍归票 92 ⇒ 与票面 `next=` 里那句"以本票的审计 sink 那一半为准、面板那一半归 92"一致。

**一处计数不严格（不改变结论）**：交件写"三枚 commit"，实际带 105 票号的是 **4 枚**（两枚纯 docs）。文件清单结论不受影响。

**AC#4 裁决：通过（未扩界）**〔独立复现：4 枚 commit 的 `--name-only` 与越界 grep 我亲手跑〕。

---

## AC#5 — 门禁（按票面原文跑）

**口径**：全部跑在 `f6818f2` 的**纯净快照** `/tmp/ac105b-gates`（`git archive f6818f2 | tar -x -C …`，容器/测试内 `ls -l /wisp/go.mod` 之类的"文件在"证明逐条附）。
交件的读数是在**当前工作树**量的（它自己标了"树＝当前工作树"）⇒ 两处不同步的数我按快照重算，并点名差异。

| 判据 | 我亲手跑的命令 | 读数 |
|------|----------------|------|
| `go test -count=2 -v` `internal/tools` | `go test -count=2 -v ./internal/tools/` | **rc=0**；`=== RUN`=**230** ＝ 不同名 115（79 顶层 + 36 子测）× 2；顶层 `--- PASS`=158、子测 PASS=72、**FAIL=0、SKIP=0**（带 `-v`） |
| `go test -count=2 -v` `internal/risk` | `go test -count=2 -v ./internal/risk/` | **rc=0**；`=== RUN`=**328** ＝ 不同名 164（96 顶层 + 68 子测）× 2；顶层 PASS=190、子测 PASS=136、FAIL=0、**SKIP=2** |
| ⚠ 我第一次跑（两包一次调用，且与 d22scan **并行**） | `go test -count=2 -v ./internal/tools/ ./internal/risk/` | `internal/tools` ok；**`internal/risk` rc=1：`--- FAIL: TestResolvePerCallBudget` 两遍**，原文 `C26 Resolve: 1667835 ns/op = 1.668 ms/op (budget 1.000 ms, 1700 samples)` → `C26 budget breach` |
| 静默重跑同一目标 | `go test -count=2 -v ./internal/risk/`（机器上无任何其他 go 进程，`tasklist` 已核） | rc=0；同一条预算读数 `0.657 ms/op`、`0.508 ms/op`（预算 1ms） |
| `gofmt -l` | 三文件（`bridge.go`、两枚 105 测试） | 空，rc=0 |
| `gofumpt -l` | `"$(go env GOPATH)/bin/gofumpt.exe" -l` 同三文件 + 整两包 | 两次都空，rc=0（该二进制本机存在：`D:\work\base\gopath\bin\gofumpt.exe`） |
| `go vet` 按包 | `go vet ./internal/tools/ ./internal/risk/` | rc=0 |
| `GOOS=linux go vet` 按包 | `GOOS=linux go vet ./internal/tools/ ./internal/risk/` | rc=0 —— **只编译不执行**，我不拿它当"POSIX 测过了" |
| `GOOS=linux go vet ./...` 全仓 | 同形 | **rc=1**，错误原文末三行：`package github.com/CarlosShao/wisp/cmd/wisp / imports …sherpa-onnx-go-linux: build constraints exclude all Go files in …`（已知洞，非本票引入） |
| `sh scripts/d22scan.sh` | 纯净快照，三枚树各跑一遍 | 三枚 **rc=0 clean**（见下方台账表） |
| POSIX **真跑**（不是交叉 vet） | `MSYS_NO_PATHCONV=1 docker run --rm -v "C:/Users/swq/AppData/Local/Temp/ac105b-gates:/wisp" -w /wisp -e CGO_ENABLED=0 golang:1.27 sh -c 'ls -l /wisp/go.mod; …go test -count=2 -v -run "…四个名字…" ./internal/tools/ ./internal/risk/'` | **rc=0**；先 `ls -l /wisp/go.mod` 与 `…ticket105_windows_test.go` 证明**不是静默空挂载**；`=== RUN`=8 ＝ 4 名 × 2，`--- PASS`=8、FAIL=0、SKIP=0；`ok internal/tools 0.064s`、`ok internal/risk … [no tests to run]`（AC#3 两枚是 `//go:build windows`，POSIX 不参与，与自述一致） |
| ⚠ 第一次 docker 尝试**没跑起来**，原文登记 | `docker run … -w /wisp …`（Git Bash 下 MSYS 把 `-w /wisp` 折成了 Windows 路径） | `docker: Error response from daemon: the working directory 'D:/work/soft/Git/wisp' is invalid, it needs to be an absolute path` ⇒ 那一次**不算** POSIX 测过，换 `MSYS_NO_PATHCONV=1` 重跑才有上面的读数 |

**那条 `TestResolvePerCallBudget` 的 FAIL 不记在本票账上，理由我按实写**：它红的那一次我**同时在跑 `sh scripts/d22scan.sh`**（该脚本会 build 并执行 d22scan 自己的测试），CPU 被我自己占住了；这是资源类判定，编队不安静时不作数（本仓 D32 那条口径）。静默重跑两遍读数 0.657 / 0.508 ms/op，都在线内。
两遍数都留着，不挑运气那次；顺带记一条：**这条预算用例对并行负载敏感**，属既有脆弱点（票 18 AC#6 的形状），本票没改它也没加重 ⇒ 登记 `R-105-3`（不阻塞）。

**SKIP 我自己核了源码**（不信转述）：`internal/risk/syncdirs_windows_test.go:124-134`——
`roots := registryProbe(probeEnv{Home: userHomeDir()})` 之后 `if len(roots) == 0 { t.Skip("no registry-grade sync record on this machine (HKCU Accounts without UserFolder is the documented reality here)") }`。
⇒ 唯一 SKIP 名＝`TestSyncRegistryProbeLive`，两遍（`-count=2`），**带 `-v` 才报得出来**，与自述一致。不是"跳过冒充通过"：它断的是 HKCU 探针返回值的形状，本机确实没有那条记录；AC#3 的两条用例与它无关。

**AC#2 的三条端到端用例确实跑到了**（防"跑错对象"）：按名字在 tools 日志里点名数过——
`TestRealToolCallWritesRewriteAccountIntoAudit` RUN=2 PASS=2、`TestUnusableRootIsVisibleInTheAuditRecord` RUN=2 PASS=2、`TestAccountRecordIsWrittenEvenWithNothingToReport` RUN=2 PASS=2。

**d22scan 台账（同一条脚本，三枚树）**：

| scope | `2afd064^`=`d43e277`（105 代码之前） | **`f6818f2`（被验树）** | 当前 HEAD=`f6a86db` | 交件自述（脏工作树） |
|-------|--------------------------------------|--------------------------|---------------------|------------------------|
| bans#1-5 `internal/` | 202 | 202 | 202 | 202 ✅ |
| bans#1-5 `cmd/` | 20 | 20 | 20 | 20 ✅ |
| ban#6 `frontend/` | 40 | 40 | 40 | **43** ❌ 见下 |
| ban#7 `internal/tools/` | 18 | 18 | 18 | 18 ✅ |
| ban#8 `design/` | 16 | 16 | 16 | 16 ✅ |
| ban#8 `frontend/` | 40 | 40 | 40 | **43** ❌ |
| ban#8 `internal/` | 363 | **370** | **371** | 366 ❌ |
| ban#8 `cmd/` | 26 | 26 | 26 | 26 ✅ |

⇒ "各 scope 不降"（A64②）**在纯净树复算成立**：`internal/` 363→370（本票新增 2 枚测试）升、其余各 scope 持平，rc=0 clean。
**两个数字对不上，我点名**：`frontend/`=43 与 `internal/`=366 都只出现在**脏工作树**上（43 里含票 92b 当时未提交的 frontend 文件）。它 AC#5 里确实标了"树＝当前工作树"，所以这是**口径披露不足**而不是谎报数字；但"不降"这个结论必须建立在同一棵树的两枚点上才有意义，交件没给出那个基线（我给：`d43e277` 那枚）。

**AC#5 裁决：通过**〔独立复现：表中每一行都是我在 `f6818f2` 纯净快照上亲手跑的，含 docker 里的 POSIX 真跑；第一次 docker 尝试因 MSYS 参数折路径而**未跑成**，错误原文已登记，未拿它充当读数〕。

---

## 总判

**FAIL-退回（只退 AC#3 这一格，其余四格通过）。**

一句话理由：AC#3 声称要防的结局——"祖先被 C26 移走之后，判定拿移动后的树当证据"——我用一条 40 行行为用例造出来了，而**只删 `ares.Actable()` 这一发变异在本票交付的两条用例下不红**（我的探针红：`Sync:false`、`Why="write target is not under any sync root"`），也就是说"那半今天只能钉符号、钉不住行为"这句自述**不成立**，交件据此少写了一格行为用例；按票 108 的先例，这种格子不许写成"通过附条件"。

四格通过的部分我不打折：AC#1 复现成立（前一格由 `acceptor-ticket105` 独立复算，我未重跑）；AC#2 的读者从装配根一路走到 `book()`、生产调用点非零、`loop.go:727 → bridge.go:243 → :756/:770 → :820-825 → :864` 逐级对上，三条端到端用例按名字点名复跑（各 2 遍、全绿）；AC#4 零越界；AC#5 五个数加 POSIX 真跑全部复现（含两条交件数字口径差异的点名）。

### 四种假绿逐条点名（AC#3/#4/#5 这三格）

- **跳过冒充通过**：唯一 SKIP＝`TestSyncRegistryProbeLive`（2 遍、带 `-v`、条件源码我逐字核过），不覆盖本票任何一条判据。AC#3 的 POSIX 不参赛是**构建标签**决定的（junction 是 Windows 产物），交件已写明，不算假绿。
- **断言恒真**：AC#3 用例 1 在 M2 下会红 ⇒ 非恒真；用例 2 对这条腿**永远拿不到红**（非 exempt 祖先在第一次 `Resolve` 就被拒），它断的是另一条腿——不算假绿，但**不该按本票的"第二条行为用例"计数**。
- **跑错对象/错的文件**：交件 AC#5 的 d22scan 与门禁读数取自**当前工作树**（共树有票 92b/111/113 在飞），`frontend/=43`、`internal/=366` 两枚数在纯净树上复算不出来（40、370）。它标了口径，故不判假绿，但"不降"的结论我换了基线（`d43e277` → `f6818f2`）才立得住。
- **门禁压根没跑**：`GOOS=linux go vet` 全仓 rc=1（sherpa-onnx 构建约束）是**已知洞**，交件没拿它冒充 POSIX 测试；POSIX 真跑两边都做了（它一次、我一次），读数一致。
- ⚠ 我自己的假绿也要点名：第一次 `-count=2` 我给 internal/risk 量到 rc=1（预算用例 1.668ms/op），原因是**我并行跑了 d22scan**。静默重跑 rc=0（0.657 / 0.508 ms/op）。这条红不记在本票账上，两遍读数都留着。

## R-105-x 新账（建议，不代立案——票面与台账不归我写）

| 编号 | 内容 | 地界 | 阻塞结案？ |
|------|------|------|------------|
| **R-105-2** | `syncdirs.go:226-229` 祖先 `Actable()` 那半**没有行为用例**；我的 `zz_probe_105b_reachable_windows_test.go`（在 `/tmp/ac105b-acc105b/internal/risk/`）可直接搬进仓，配 M3 变异自证 | `internal/risk/`（只加测试，不改判定） | **阻塞**（就是 AC#3 退回的修复项） |
| R-105-1 | `PATH-ACCOUNT` 审计行在生产只到 **stderr**（`cmd/wisp/run.go:428`），"人能看见哪条路被改写过"仍未兑现；要么面板读这本账（票 92），要么审计流收进文件 | `internal/panel/`（票 92）或 `cmd/wisp`+`internal/observe` | 不阻塞（本票 AC 只要求"与 MODE-READ 同族"） |
| R-105-3 | `TestResolvePerCallBudget`（票 18 AC#6）对并行 CPU 负载敏感，1ms 预算在编队不安静时自己会红——本代理已踩一次 | `internal/risk/` 门禁形状 | 不阻塞；建议与既有"资源类判定"账查重后并入 |
| R-105-4 | `Resolve(anc)` 的 deny 分支（`syncdirs.go:223-225`）我**推断**也可经同一"清洗前后 `%` 配对差"到达（展开后的祖先自己经过未登记 junction），但**没造出探针** ⇒ 只是第三档观察，不足以立案 | `internal/risk/` | 不阻塞；要先有探针 |
| R-105-5 | 共树下"各 scope 不降"的对比必须锚定**同一棵纯净树的两枚 sha**；在脏工作树量的台账数不能与被验版本比 | 纪律/编排 | 不阻塞 |

## 下一张该派什么

1. **立刻派一枚小写码票（或把本票退回同一实现方）**：只做 R-105-2——在 `internal/risk/` 加祖先 `Actable()` 的行为用例（我的探针是现成的），并自带 M3 变异读数。**不许改 `syncdirs.go`**（票 102 语义）。这一枚完本票才可结案。
   替代路线（**需 owner 拍板、属改判据**）：承认那半由静态判据（`pathresolver_rewrite_account_test.go`）承载，把 AC#3 判据文字改成"重解析整段"，并在票面写明"Actable 那半只在 env 名含分隔符时可达"。我不推荐：那等于把本票立票要消灭的形状重新写回判据。
2. 面板那半（R-105-1）仍归票 92，别塞回本票。
3. 资源类门禁（R-105-3）等编队安静的一轮再单独测，本轮不派。
