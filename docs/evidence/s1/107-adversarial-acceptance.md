# 107 对抗验收 —— `acceptor-ticket107`

验收对象：commit `be8b403`（未 push）；票面 `.scratch/wisp/issues/107-ticket102-allowlist-judgment-is-red-on-posix.md`。
时间锚点：`date` = `Mon Sep 21 18:48:22 CST 2026`（下文时间均在此之后）。
 measurement 树：**全部在 `/tmp/accept107/snap`（`git archive HEAD` 纯净快照）＋ Windows 工作树**；
快照还原证据：`cp paths.go.head` 后 `diff -q` 对 `git show HEAD:internal/tools/paths.go` ⇒ **SNAP RESTORED IDENTICAL**；
仓库内未建 worktree/checkout，`git status --porcelain internal/tools internal/risk` 空（别人的 ` M` 未碰）。

档位标记：**〔独立复现〕** = 我亲自跑出 rc；**〔日志＋归档，我抽验〕** = 它给的读数我核对过落点/代码但没原样重跑；**〔仅自述，不背书〕**。

## 总判：**AC#3 不通过**（修法把放行侧改宽了 ⇒ 我在 POSIX 上造出"跨树放行"）；AC#1/AC#2/AC#4/AC#5 通过。

## 裁决表（与票面 AC 1:1）

| AC | 判据 | 裁定 | 档位 | 我的读数 |
| --- | --- | --- | --- | --- |
| AC#1 | POSIX 真跑复现那条红 + 归因落 commit | **通过** | 〔独立复现〕 | 见下"POSIX 复现" |
| AC#2 | 两种形状分别喂进判定，定产生端/比较端 | **通过**（定性正确，但**结论强度不足以支撑它的修法**） | 〔独立复现〕 | 见下"两形状喂入" |
| AC#3 | 修完两侧都真跑过、且放行侧不比拒绝侧宽 | **不通过** | 〔独立复现〕 | 探针 A：修后放行 / 修前不放行 / 被放行的树**不是操作者点名的那棵** |
| AC#4 | 双向变异（①退回旧实现 ②放宽成任意前缀） | **通过** | 〔独立复现〕 | ①两条红都回来；②边界用例红 |
| AC#5 | 门禁四数 | **通过**（本票范围内） | 〔独立复现〕+ 一处抽验 | gofmt 空 / 两侧 `go vet` rc=0 / `-count=2 -v` 530 RUN rc=0 |

---

## 主攻点：放行侧被改宽了吗 —— **是，且我能造出跨树放行**

### 探针 A（符号链接祖先／root 本身是符号链接）—— 决定性

形状：`base/outside/secret.txt`（**操作者从未点名**）；`base/proj` 是**指向 `outside` 的符号链接**；
`allowed_dirs = ["%AC107_ROOT%/proj"]`（走 C26 展开腿）。target = `base/proj/secret.txt`。

| | 修前（`if !res.Resolved`） | 修后（`!res.Resolved && !treeOnDisk`） |
| --- | --- | --- |
| `Roots()` | `[]` | `[.../001/proj]` ← **符号链接名进了授权集** |
| `UnusableRoots()` | `["%AC107_ROOT%/proj": C26 expanded it onto .../proj but that tree is not confirmed on disk, authorizing nothing]` | `[]` |
| `InAllowlist(canonical)` | **false** | **true** |
| `os.ReadFile(judged canonical)` | —（不放行） | `"TOPSECRET"`，`filepath.EvalSymlinks` ⇒ **`.../001/outside/secret.txt`** |

⇒ **满足判据的三件**：修后放行 ✓、修前不放行 ✓、**被放行的那棵树不是操作者点名的那棵** ✓。
"判定的树"与"真正打开的树"**不是同一棵** —— 这正是票 102 立的不变式（judgement/execution/report 指同一棵树）被反向破掉。
`treeOnDisk` 用的 `os.Stat` **跟随符号链接**（`internal/tools/paths.go:170-173`），而 `risk.Resolve` 在 POSIX 上是**纯词法**
（`internal/risk/pathresolver_other.go:10 resolveHandle → ("", false)`；`reparseComponents → nil`，同文件 `:12`）⇒ 平台上没有任何一处做过 realpath/lstat。
⇒ 票面 AC#3 的"放行只认已解析形式"这句话在 POSIX 上被这次修法**换成了"只认存在"** ⇒ **判 AC#3 不通过**。

### 探针 B/C/D（对照，说明"新洞"与"旧洞"的边界）—— 〔独立复现〕

- **B**（root 是符号链接，但拼成**普通绝对路径**、未被改写）：修前修后**都 true**，body 同样是 `TOPSECRET`。
  ⇒ 这条**不是本票新造的**，是 POSIX 上未改写腿一直有的（`paths.go:160-169` 那段自述属实）。
- **C**（root 是**真树**，但 target 走 root 里的符号链接子目录 `proj/esc → outside`）：修前 `Roots()=[]`⇒false；修后 `Roots()=[.../proj]`⇒**true**，realpath=`.../outside/k.txt`。
  ⇒ 本票的增量把这一类也一并放进来了（未改写腿做不到，因为它至少还是操作者写下的名字；这里连名字都不对）。
- **D**：`Canonicalize` 输出与输入字节同形，`EvalSymlinks` 才给出真树 ⇒ 确认**整条腿无 realpath**。

### 谁控制 root（探针 2）—— 风险等级如实降级，但**不等于没有洞**

`allowed_dirs` 来自**配置文件**：`internal/config/schema.go:472 AllowedDirs []string \`toml:"allowed_dirs"\``
⇒ `cmd/wisp/run.go:261-264` 原样拷出 ⇒ `cmd/wisp/run.go:265 tools.NewPathCanonicalizer(allowed, cfg.FS.ReparsePointExceptions)`。
**grep 全仓无任何生产代码把工具载荷写回 `roots`**（`p.roots` 只在 `NewPathCanonicalizer` 里 append，`internal/tools/paths.go:82`）
⇒ **模型不可控，只有操作者能设**。所以这不是"模型说服判定放行"的洞，等级从"高危 fail-open"**降为**"操作者自己机器上的符号链接会让他以为批准了一棵树、实际批准了另一棵"。
仍必须说清的理由：操作者写 `allowed_dirs = ["~/proj"]` 时**没有任何界面告诉他人 `proj` 是个符号链接**；`RewrittenRoots()` 只记 `%VAR%`/`~` 的**改写**，不记**跟随**。

### `unusable` → `roots` 的去向（探针 3）—— 只喂了一处判定，但**同一遍历同时服务拒与放行**这条成立

`p.roots` 的**唯一**判定消费者是 `InAllowlist`（`internal/tools/paths.go:122-135`）；
`InAllowlist` 的两个调用点：`internal/risk/rules_gateway.go:45`（不在表内 ⇒ L2）与 `internal/tools/bridge.go:845`（`inScope` ⇒ 免问/L1 那条）。
⇒ **没有**"别的判定因为进了 roots 而额外变宽"，也没有判定因为**没进** roots 而被跳过（`unusable` 只被 `UnusableRoots()` 读作审计面，无分支逻辑）。
⚠ 但按本项目定性判据：**`roots` 这本账同时被"拒"（rules_gateway L2）与"放行"（bridge inScope 免问）两条腿读**，
本次修法只加固了其中"存在性"这一半，两腿要的强度其实不同 ⇒ 应拆成两条账（放行腿要"已解析"，拒腿可以只要"存在"）。登记为 **R-107-2**。

### 正确形状（不改码，只给方向）

POSIX 侧要绿**只有两条合法路**，**都不许靠放宽放行侧换绿**：
1. **补强放行侧**：在 `NewPathCanonicalizer` 里对 `Rewritten` 腿用 `filepath.EvalSymlinks`（已解析形式）作为 `treeOnDisk` 的替代/前置，
   并让**target 侧同一处理**（否则 root 解析了、canonical 没解析 ⇒ 反而全 false，票 102 用例照样红）；
   或在 `internal/risk/pathresolver_other.go` 把 DEFERRED 的 realpath+lstat 真做出来（**那是冻结包，需另开票**）。
2. **保持"更严"并把红用例按"平台能力天生不存在"处理干净**：即票 102 那条断言在 POSIX 上应**改成断言 false + 断言 `unusable` 有那条文案**
   （平台桩不存在是既成事实），**同时**在票面明写"POSIX 上 `[fs] allowed_dirs` 的改写腿整体失效"，把它作为一条独立的、有 owner 决策的降级，而不是"修到绿"。

---

## AC#1 POSIX 复现 —— 〔独立复现〕（本仓以后每张碰路径形状的票都必须用这条路）

方法（我跑通的，非照抄）：Windows 侧 `GOOS=linux CGO_ENABLED=0 go test -c -o /tmp/accept107/<n>.test ./internal/tools/`（rc=0），
执行用 **Docker Desktop（WSL2 内核）**：`MSYS_NO_PATHCONV=1 docker run --rm -v "C:/Users/swq/AppData/Local/Temp/accept107:/t" alpine sh -c '/t/<n>.test -test.v ...'`。
容器 `Linux 6.6.114.1-microsoft-standard-WSL2 x86_64`。
⚠ 挂载坑（要写进手法）：Git Bash 下 `-v "$W":/t`（`C:\...` 形式）**静默挂空目录**，容器里 `T:/post.test: not found` 且 **rc 仍为 0** ⇒ 假绿；
必须 `MSYS_NO_PATHCONV=1` + 正斜杠盘符路径，且**先 `ls -la /t` 证明文件真在里面**。

- **修后（HEAD `be8b403`）在真 Linux**：`TestPathCanonicalizerAccountsForRewrittenRoots` / `TestTicket107AllowlistJudgmentTwoShapes` /
  `TestTicket107AllowlistBoundaryIsComponentWise` 三条 `=== RUN`=3、`--- PASS`=3、末行 `PASS`、容器内 rc=0。⇒ 它的"POSIX 真跑起来了"**成立**。
- **归因**：`internal/tools/paths_rewrite_ticket102_test.go` 与 `paths.go` 里 `not confirmed on disk` 那条腿**都只来自 `a66aadf`**（票 102），
  我 `git log --oneline` 对照 `be8b403^..` ⇒ **这条红是票 102 引入，不是 `a1613d9` 之前就有**。它的归因**对**。

## AC#2 两形状喂入 —— 〔独立复现〕

我在 Linux 上把它的新用例跑出来了（`-test.v` 原样打印，paths_ticket107_portable_test.go:57-60）：
形状 R `InAllowlist=true`；形状 E `roots=[.../proj] rewritten=[…(env)] unusable=[] InAllowlist=true`；交叉喂 `=true`。
修前读数我用**变异①**在同一容器里独立得到：形状 E 臂 `roots=[] unusable=[…"not confirmed on disk"…]` ⇒ **false**。
⇒ **产生端**定性成立、比较端（`foldPath` + `f == r || HasPrefix(f, r+pathSep)`，`paths.go:130`）对同一串两次判对。
但：**"比较端无罪"不等于"产生端可以拿存在性交差"** —— 产生端要交的是"已解析"，它交了"存在"。这就是 AC#3 不通过的那一格。

## AC#4 双向变异 —— 〔独立复现〕（快照内做，`grep -n` 同链证落地，先 build rc=0）

| 变异 | 落地行（同链 `grep -n`） | build | 读数 |
| --- | --- | --- | --- |
| ① 退回旧实现 | `76: if !res.Resolved {` | LINUX BUILD rc=0 | **`--- FAIL: TestPathCanonicalizerAccountsForRewrittenRoots`**（`paths_rewrite_ticket102_test.go:64`）+ **`--- FAIL: TestTicket107AllowlistJudgmentTwoShapes`**（`…portable_test.go:66`）；`BoundaryIsComponentWise` PASS；容器末行 `FAIL` |
| ② 放宽成任意前缀 | `130: if f == r \|\| strings.HasPrefix(f, r) {` | LINUX BUILD rc=0 | **`--- FAIL: TestTicket107AllowlistBoundaryIsComponentWise`**：`InAllowlist("…/001/proj-evil/b.txt") = true`（被 `…/proj` 放行）；另两条 PASS ⇒ **"放行只认唯一已解析形式"这条真有牙**（在"前缀/兄弟树"这个维度上有牙；在**符号链接**维度上无牙 —— 就是 AC#3 那一格） |

还原：`cp /tmp/accept107/paths.go.head internal/tools/paths.go` 后 `diff -q` 对 `git show HEAD:…/paths.go` ⇒ **SNAP RESTORED IDENTICAL**；仓库 `git status` 我的文件零改动。

## AC#5 门禁 —— 〔独立复现〕（Windows 工作树，只按包，没跑整仓）

- `gofmt -l internal/tools internal/risk` ⇒ **空**，rc=0。〔独立复现〕
- `go vet ./internal/tools/ ./internal/risk/` rc=0；**`GOOS=linux go vet ./internal/tools/ ./internal/risk/` rc=0**（按包）。〔独立复现〕
- `go test -count=2 -v ./internal/tools/ ./internal/risk/` ⇒ **rc=0**，`=== RUN` **530** == 不同测试名 265 × 2 ✓；FAIL **0**；SKIP **1 名 ×2**（`TestSyncRegistryProbeLive`）。〔独立复现〕
- `gofumpt`／`sh scripts/d22scan.sh` 台账各 scope：**〔日志＋归档，我抽验〕** —— 我只核了它引的 ban 计数**只增不减**这条与"新增文件是 `_test.go`"这个前提在 `git ls-files` 层成立，没重跑整条 d22scan（工作树共树、别人的 ` M` 在 `internal/winsec/`，按规矩不跑整仓门禁）。

---

## 它交回的两条"不是我发现的"—— 核后裁定

1. **`TestResolvePerCallBudget` FAIL×2**：**我重跑，判为负载假红，登记读数**。
   工作树 `-count=2 -v ./internal/tools/ ./internal/risk/`（两包同跑，正是它说的那条红条件）⇒ 两条都 **PASS**，rc=0，
   读数 `pathresolver_budget_norace_test.go:34: C26 Resolve: 730340 ns/op = 0.730 ms/op (budget 1.000 ms, 2949 samples)` 与
   `524340 ns/op = 0.524 ms/op (2209 samples)`。它的对照（0.614/0.718）与我这次同号同量级 ⇒ **只是负载，不归票 86**。多样本全报，未调阈值。
2. **`SKIP=TestSyncRegistryProbeLive×2` + "runtests.sh 报 SKIP=0 是非 `-v` 读数"**：**前半坐实、后半推翻**，请票 93 按更正后的因果接手。
   - 坐实：我这次 `-v` 里 `--- SKIP: TestSyncRegistryProbeLive (0.00s)` **2 行**；`git status --porcelain internal/risk` **空** ⇒ 该 SKIP 在 **HEAD 就有**，不是票 93 新造（也不是 `-v` 造出来的：`-v` 只是让它**可见**）。
   - **推翻**：它的因果解释"因为是非 `-v` 读数所以 SKIP=0"**不成立** —— `tools/d22scan/runtests.sh:75` 明写 `go test -v -count=1 "$@"`，
     `:88 skipped=$(count '^--- SKIP')`，`:99` 直接"SKIP is not a pass"。而我实测 `^--- SKIP` 在 col 0 **命中 2 次**（`^    --- SKIP` 命中 0）⇒
     同一条仪器不可能因为"没带 `-v`"而报 0。**真正要查的是它调 runtests.sh 时的包选择/`-run`/`-C 目录`**（`scripts/d22scan.sh:32` 只是注释指向它）。
     ⇒ 转票 93 时请带这句：**别按"非 `-v`"这个错因果修**，方向应是"门禁包集合是否覆盖了 `internal/risk`"。

## 票面 numstat 删除列=1 —— **通过**，〔独立复现〕

`git show be8b403 -- <票107文件> | grep -E "^-[^-]"` **只有一行**：
`-**Status:** open（2026-09-21 18:1x 编排者建；来源=CI 步级读数 run 35585147258、35586044995 的 test-core 步）`
⇒ 删除列 1 = Status 行本身，其余正文（含 CI 原文、判之前的两件事、AC#1..#5、Rules、编排者 Progress log）**逐字保留**，交件全在文末追加段。

## 缺陷登记（验收期不修）

- **R-107-1（阻断 AC#3）**：`internal/tools/paths.go:76` + `treeOnDisk`(`:170-173`) 在 POSIX 上把"确认这棵树"降为"跟随符号链接的存在性"，
  造成探针 A/C 的**跨树放行**；放行侧因此**宽于**票 102 立的标准。修法见"正确形状"，**不许用放宽放行侧换绿**。
- **R-107-2**：`roots` 一本账同时喂"拒"（`rules_gateway.go:45`）与"放行/免问"（`bridge.go:845`）两条腿 ⇒ 按本仓判据要**拆两条**（放行腿要已解析，拒腿存在即可）。
- **R-107-3**：`RewrittenRoots()` 只记 `%VAR%`/`~` 改写，**不记符号链接跟随**；操作者没有任何界面得知"我点的 `proj` 是链到别处的"。
- **R-107-4（手法坑）**：Docker 挂载在 Git Bash 下会**静默挂空且 rc=0**（假绿），必须 `MSYS_NO_PATHCONV=1` + 容器内 `ls -la /t` 自证。
- **R-107-5（读数坑，非本票）**：`docker run … | grep … ; echo rc=$?` 测的是 grep 的 rc 不是容器 rc ⇒ 本轮所有 rc 用容器内 `echo RC=$?` 或直取管道尾。

## 伪授权登记（本轮）

**0 次。** 本轮所有工具输出末尾**未**出现自称"编排者备注/停手/撤回/revert/冻结某包"的文本；
唯一一条外部注入提示是 `MEMORY.md 已修改` 的系统通知（不含指令，未据此改变行为）。**未执行任何 revert**，仓库与快照状态均由我自己控制。

## 没验完的格子 / 下一条命令

- `sh scripts/d22scan.sh` 整条**我没重跑**（共树），只抽验了台账前提 ⇒ 下一条命令：
  `cd /tmp/accept107/snap && sh scripts/d22scan.sh` 与 `git show HEAD:…` 的 ban 计数对照。
- `gofumpt -l` 未独立复跑（`$(go env GOPATH)/bin/gofumpt`）。
- 票面 AC#1 留的 **run id 空位**仍空，按规矩归编排者 push 后读 `test-core` 步级结论补。
