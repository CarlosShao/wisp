# 票 75 — 独立验收复现（acceptor，非实现者、非编排者）

验收代理：`acceptor-ticket75`（独立复现票 75 声称的东西；不写生产码、不 commit、不 push）
日期：2026-09-21　分支：`dev`　被验收的 SHA：**`17efc2c`**
本文件是 append-only：每完成一组就追加，中途死掉也能从文件里接着干。

## 环境与被测树的取法（两组共用）

```
rm -rf /tmp/wisp75 && mkdir -p /tmp/wisp75
cd "D:\work\workspace\projects plans\Wisp" && git archive 17efc2c | tar -x -C /tmp/wisp75
```

- 快照在**仓外**（`/tmp/wisp75` = `C:/Users/swq/AppData/Local/Temp/wisp75`），**没有 git 元数据**，
  工作树里在途代理（票 70/80/81、`internal/memory/artifacts_path_invariant_test.go` 的未提交改动）
  一个字节都没进被测树。全程未在仓内建 worktree、未 checkout、未 stash。
- Linux 侧：docker `golang:1.27`（容器内 `go version go1.27.1 linux/amd64`，
  `Linux 6.6.114.1-microsoft-standard-WSL2 x86_64`）。
  容器里 `WISP_ENV=test`（与 `ci.yml` 的 `test-core` job env 一致），
  module cache 用命名卷 `wisp75-gomod`、build cache `wisp75-gobuild`。
- 命令形状：`.github/workflows/ci.yml:127-134` 的 "Portable package tests" 是
  `go test <16 个包> -count=1`。本票的战场是其中 `./internal/risk/...` 与 `./internal/tools/...`
  两项（实测这两包**无子包**，`find internal/tools internal/risk -mindepth 2 -name '*.go'` 空 ⇒
  `./internal/risk/...` ≡ `./internal/risk/`，与 CI 逐字等价）。

---

## 第 1 组：Linux（docker `golang:1.27`）全量 `-count=1`

### 1a. 逐包 `-json` 计数（用于真实 RUN/PASS/FAIL/SKIP 数）

```
docker run --rm -v <snapshot>:/src -v <logs>:/logs -v wisp75-gomod:/go/pkg/mod \
  -v wisp75-gobuild:/root/.cache/go-build -w /src -e WISP_ENV=test golang:1.27 \
  bash -c 'go test -count=1 -json ./internal/tools/ ./internal/risk/ > /logs/g1.json 2>/logs/g1.err'
GROUP1_RC=1
```

| 包 | 包级结论 | 用时 | `=== RUN` | PASS | FAIL | SKIP |
|---|---|---|---|---|---|---|
| `internal/tools` | **pass** | 10.685 s | **56** | **54** | **0** | **2** |
| `internal/risk` | fail | 2.848 s | **113** | 104 | **8** | 1 |
| 合计 | — | — | 169 | 158 | 8 | 3 |

56=54+0+2、113=104+8+1 ⇒ **RUN 与结论行逐条对得上，没有"跑了但没记账"的**。
**零 `panic: test timed out`、零 `panic:`**（g1.json 里 panic/timeout 检索为空）⇒ 票 75 声称的
600 s 超时在 `17efc2c` 上**确实不存在**（AC#3 成立）。

### 1b. 同树的 CI 逐字形态（不带 `-json`，纯文本 + 真实 exit code）

```
docker run ... bash -c 'cd /src && go test ./internal/risk/ ./internal/tools/ -count=1 2>&1 | tee /logs/g1b.txt'
GROUP1B_RC=1
grep -c -- "^--- FAIL" /logs/g1b.txt  ->  8
```

尾部原文：
```
FAIL
FAIL	github.com/CarlosShao/wisp/internal/risk	2.859s
ok  	github.com/CarlosShao/wisp/internal/tools	9.671s
FAIL
```

⇒ **`internal/tools` 全绿（`ok 9.671s`，0 FAIL）**，`internal/risk` 残留 8 条，
1a/1b 两次读数**同名同数**。

### 1c. 残留 8 条红：逐条名字 + 逐条归族（**全部不是票 75 的账**）

| # | 测试名 | 失败锚点 file:line | 断言文本（首条） |
|---|---|---|---|
| 1 | `TestWriteGateNotSelectedByPayloadKey` | `provenance_test.go:469` | precondition: the injected registry-grade root must confirm detection… |
| 2 | `TestWriteGatePlainLocalWriteNotFlagged` | `provenance_test.go:541` | 同上 |
| 3 | `TestWriteGateEveryPathTargetJudged` | `provenance_test.go:608` | 同上 |
| 4 | `TestWriteGateAllPathsNonSyncStaysExempt` | `provenance_test.go:670` | 同上 |
| 5 | `TestSyncFallbackNotDisarmableByWeakRoot` | `syncdirs_test.go:118/121` ×3（registry/config/env 三个 source） | "must count as confirmed" / "confirmed root must lift the blanket suspect net" |
| 6 | `TestSyncEnvConfiguredRoots` | `syncdirs_test.go:168` | env-grade roots are confirmed evidence |
| 7 | `TestSyncNormalNewFileWriteNotFlagged` | `syncdirs_test.go:288` | sanity: registry-grade root is confirmed |
| 8 | `TestSyncDotDotTailFailsClosed` | `syncdirs_test.go:356` | precondition: registry-grade root confirmed, so only root membership may decide |

**归族结论（不是只给数字）**：8 条**全部锚在同一判据上** —— 一个"已确认（confirmed-grade 且 canonical）
的同步根"在 POSIX 上永远拿不到，于是 `SyncDetectionComplete()` 恒 false ⇒ 兜底嫌疑网关不掉 ⇒ 这些用例
的 precondition 先红。POSIX 侧的两个桩我在被测树里逐字读到：

- `internal/risk/syncdirs_other.go:20` `func registryProbe(e probeEnv) []SyncRoot { return nil }`
  —— 注释里明写 `DEFERRED(P12-macos): … (ticket 55)`。
- `internal/risk/pathresolver_other.go:10` `func resolveHandle(p string) (string, bool) { return "", false }`
  —— 注释明写 `DEFERRED(macOS/Linux)`；根没有 canonical 级证据 ⇒ `add()` 的 `res.Resolved` 门过不去。

⇒ 属**票 55 / 票 82（POSIX sync 分层）**的账，**不计入票 75**。**注意方向**：这 8 条的修法在冻结文件里
（补 realpath/lstat），把它放宽来变绿是 fail-open ⇒ 验收判**"允许残留"**，但**票 70 AC#2 的
"`test-core` 全绿"至今不成立**（见第 4 组）。

---

## 第 2 组：变异（验收的首要攻击点）—— 守卫退回无条件折叠

### 2a. 锚点定位与落地证明（先 grep，同一条 `&&` 链里改完再 grep）

改前（`/tmp/wisp75/internal/risk/pathresolver.go`，实测行号）：
```
150:		return strings.ReplaceAll(p, "/", `\`)      <- 另一处，在 UNC 分支里，与本变异无关
167:func normalizeLocalUNC(p string) string {
168:	if !strings.HasPrefix(p, `\\`) && !strings.HasPrefix(p, `//`) {   <- 票 75 的守卫
169:		return p // not a UNC spelling: this function has no business here
170:	}
171:	u := strings.ReplaceAll(p, "/", `\`)          <- 目标折叠行
```
改动 = **删掉 168/169/170 三行**（即把 `normalizeLocalUNC` 恢复成票 75 之前的"对一切路径无条件折 `\`"）。
删后同链再 grep：
```
grep -n 'ReplaceAll(p, ' internal/risk/pathresolver.go   ->  150: … 168: u := strings.ReplaceAll(p, "/", `\`)
grep -n 'not a UNC spelling' internal/risk/pathresolver.go -> （无匹配，rc=1）
```
函数体现在是 `func normalizeLocalUNC(p string) string {` 的**第一条语句就是 ReplaceAll**。

### 2b. 这是**行为**变异，不是编译失败

```
go vet ./internal/risk/   ->  VET_RC=0，/logs/g2_vet.txt 为空（0 行 "cannot"）
```

### 2c. Linux（同一容器、同一命令、同一 module/build cache）跑变异树

```
docker run --rm -v <snapshot>:/src … -e WISP_ENV=test golang:1.27 \
  bash -c 'cd /src && go test -count=1 -json ./internal/tools/ ./internal/risk/'
GROUP2_JSON_RC=1        墙钟：12:17 -> 12:28（>10 分钟，被 go test 闹钟切断）
```

| 包 | 基线（第 1 组） | **变异树** | 增量 |
|---|---|---|---|
| `internal/tools` | result=pass 10.685 s，RUN 56 / PASS 54 / **FAIL 0** / SKIP 2 | result=**fail** **600.144 s**，RUN 53 / PASS 30 / **FAIL 21** / SKIP 1 | **+21 条红 + 超时** |
| `internal/risk` | result=fail 2.848 s，RUN 113 / PASS 104 / **FAIL 8** / SKIP 1 | result=fail **8.395 s**，RUN 113 / PASS 77 / **FAIL 35** / SKIP 1 | **+27 条红** |

变异树里 `internal/tools` 的 panic 原文（`-json` 输出，挂在
`TestVetoInsideTheWindowWritesNothing` 上）：
```
panic: test timed out after 10m0s
	running tests:
		TestVetoInsideTheWindowWritesNothing (…)
```
⇒ **AC#3 的 600 s 超时不是"变快了"，是那条 300 s 等待被去掉了**：变异树上
`TestL1WriteGoesThroughTheRealBlockWindow (300.04s)` 打出
`审批超时（300 秒未确认），C18 一律判拒绝 … RiskLevel:L2`（`wiring_test.go:124`），
`TestLoopPassesDeclaredL1WriteThroughTheGate (10.15s)`，第二条 300 s 用例走到闹钟被切。
**基线树上这两条都是 PASS、整包 9.671 s。**

### 2d. 票 75 新加的用例**真的变红**（逐条，这是关键证据）

"基线绿 → 变异红"的名字集合差（脚本对 `g1.json`/`g2.json` 取差，**基线红→变异绿 = 0 条**，
即变异只加红、不洗牌）：

`internal/risk` 新增 27 条，其中**票 75 自己写的仪器**（`internal/risk/pathshape_portable_test.go`、
`internal/risk/provenance_syncdirs_other_test.go`）红的是：
```
+ TestResolveCanonicalIsPlatformShaped
+ TestResolveKeepsSeparatorsOutOfPosixNames
+ TestDoubleSlashSpellingIsNotUNC
+ TestTierAAnchorsMatchNativeSpelling            (+ 10 个子项)
+ TestOverrideKeyKeepsPosixBackslashesDistinct
+ TestSyncAncestorWalkVerifiesNativeChain
+ TestExfilSyncWritePosixSpellingInvariant      (+ 2 个子项)
+ TestFourChannelExfilSuite / fs.write_into_sync_dir
```
打出来的失败文本就是票面的串（**没有任何 OS 能打开、也不再绝对**）：
```
pathshape_portable_test.go:47: canonical "/src/internal/risk/\tmp\TestResolveCanonicalIsPlatformShaped4231820040\001\profile\notes.txt"
                             carries the foreign separator "\\": C26 must hand back the platform shape, not the Windows one
pathshape_portable_test.go:54: the canonical C26 handed back cannot be statted (lstat …: no such file or directory)
tools/pathshape_portable_test.go:47: fs.read/fs.write open what Canonicalize returned, and that path is not there:
                             lstat /src/internal/tools/\tmp\…\payload.txt: no such file or directory
risk/pathshape_portable_test.go:246: CROSS-FILE BLEED (fail-open): the override recorded for "…" unlocked the different file "…"
```

`internal/tools` 新增 21 条，含 `TestCanonicalizeReturnsAPathTheOSCanOpen`、
`TestCanonicalizeAgreesWithTheOSName` 两条票 75 仪器 + 19 条业务红（fs.read/fs.list/move/trash/delete/
D34 写矩阵/原子写/污点/审批环）。

**变异是有选择性的，不是大锤**：票 75 另外两条仪器
`TestFoldPathKeepsPosixBackslashesDistinct`、`TestDirOfBaseOfCutOnThePlatformSeparator`
在本变异下**仍然绿**（它们钉的是 `paths.go` 的 fold 与 `dirOf/baseOf`，不是 `normalizeLocalUNC`）
⇒ 变红的确确实实是这**一处**守卫造成的。

**一条要如实剥离的读数**：`TestResolvePerCallBudget` 也在变异新增红里
（`pathresolver_budget_norace_test.go:37`，2.483 ms/op vs 1 ms 预算）。它是 wall-clock 门，
同容器里 `internal/tools` 正在烧 600 s，且实现者证据里这条在全量并跑时本就飘
（`75-rootcause-linux-path-shape.md`：全量 13 条 vs 单跑 12 条）。
⇒ **不把它记成形状变异的信号**，也不据此加戏。

### 2e. 还原与对拍（快照里没有 git，用双 oracle）

```
cp /tmp/wisp75-logs/pathresolver.go.orig internal/risk/pathresolver.go
diff -q internal/risk/pathresolver.go /tmp/wisp75-logs/pathresolver.go.orig   ->  IDENTICAL（无输出，rc=0）
diff -q internal/risk/pathresolver.go /tmp/wisp75w/internal/risk/pathresolver.go ->  IDENTICAL
     ^ 第二个 oracle：为第 3 组另起的 `git archive 17efc2c` 独立归档，与 .orig 无关
grep -c 'not a UNC spelling' internal/risk/pathresolver.go  ->  1
diff -r --brief /tmp/wisp75 /tmp/wisp75w                   ->  空（rc=0）
     ^ **整棵被测树与 17efc2c 的独立归档逐字节相同**：变异没留残渣，也证明第 1/2 组跑的是同一棵树
```

### 2f. 第 2 组结论

`17efc2c` 的守卫是**承重的**：把它退回票 75 之前的无条件折叠，Linux 上
`internal/tools` 由 **0 FAIL → 21 FAIL + 10 分钟超时 panic**、`internal/risk` 由 **8 → 35**，
且票 75 新写的 8 组（含子项 20 条）形状用例**逐条转红并打印出票面描述的那条串**。
⇒ 票 75 的修复**不是靠测试松口的绿**。**PASS**。

---

## 第 3 组：Windows 本机（票 75 新增的两条用例）

测量树：`/tmp/wisp75w`（**同 SHA `17efc2c` 的第二次独立 `git archive` 快照**，不是被变异过的那棵；
宿主 `go version go1.27.1 windows/amd64`，Windows 原生，非交叉编译）。

```
go test ./internal/risk/ -run 'TestOverrideKeyKeepsPosixBackslashesDistinct|TestExfilSyncWriteWindowsSpellingInvariant' -v -count=1
GROUP3_RC=0
=== RUN   = 9      --- PASS（含子项）= 9      --- FAIL = 0      --- SKIP = 0
无 "testing: warning: no tests to run"（专门 grep 过）
ok  	github.com/CarlosShao/wisp/internal/risk	0.057s
```

- `TestOverrideKeyKeepsPosixBackslashesDistinct`：**没有 build tag**
  （`head -5 internal/risk/pathshape_portable_test.go` 第一行就是 `package risk`），
  Windows 上 **`=== RUN` 有、`--- PASS (0.00s)` 有、零 SKIP** ⇒ 满足"若 SKIP 即验收不通过"这条硬规矩。
- `TestExfilSyncWriteWindowsSpellingInvariant`：`//go:build windows`
  （`provenance_syncdirs_windows_test.go` 第 1 行），7 个子项**逐条列出并逐条 PASS**
  （4 条 sync 根内拼写 + 3 条根外拼写），0.02 s。

⇒ Windows 侧两条新仪器**真跑、真绿、无 SKIP、非空匹配**。**PASS**。

### 3b. AC#2 的另一半（Windows 全包 `-count=2`）——**这台机器上不复现 rc=0**

```
cd /tmp/wisp75w && go test ./internal/risk/ ./internal/tools/ -count=2
WIN_COUNT2_RC=1     --- FAIL = 2（同一个名字的两轮）  --- SKIP = 0
FAIL	github.com/CarlosShao/wisp/internal/risk	39.102s
ok  	github.com/CarlosShao/wisp/internal/tools	43.819s
```
两条红都是 **`TestResolvePerCallBudget`**（`pathresolver_budget_norace_test.go:37`）：
`1.604 ms/op` 与 `1.596 ms/op`，预算 1 ms，样本 1172/714。
**这是 wall-clock 门，不是我这台机器的形状问题，而且这次读数被污染了**：
读数期间 `docker ps` 显示**另一个代理的 `golang:1.27` 容器在跑**（`Up 2 minutes`），
我自己的 AC#1 预修复基线容器同时在烧 tools 包；risk 包整包用时从实现者报的 4.853 s
变成我的 39.102 s（≈8 倍），量级只可能来自争抢。
⇒ **这条读数不作指控**，等机器安静后复跑（见 3c 与本节末尾）。
`internal/tools` 在 Windows `-count=2` 下 **ok 43.819s**，零 FAIL 零 SKIP ⇒ AC#2 的 tools 半边成立。

---

## 第 4 组：票面对账 + 真 CI run 取证

### 4a. run `35558750456`（headSha `17efc2c1d950975f712c8c159221be096de91c1a`，workflow `ci`，run #125，status **completed**，非 cancelled）

`gh run view 35558750456 --json jobs` → job 级：
`test-windows=success` / `test-core=failure` / `slo-full=success` / `lint=failure` / `slo-smoke=success`。

`gh run view 35558750456 --json jobs --jq '.jobs[]|select(.name=="test-core")|{conclusion,steps:[…]}'`
→ **步骤级**（票面 A44③/A49④ 要的正是这个）：
```
"Portable package tests (agent/llm/config/memory/observe/secret/risk/statemachine/...)" : failure
"Environment fork assertion (WISP_ENV=test data dir)"                                  : skipped   <-- 见 4d
"Probe mock-llm" / "Start compose test services" / "Stop compose services"             : success
```
`test-core` job id = **106207404142**；步骤内 `##[error]Process completed with exit code 1.`。

我另外确认了引用不腐坏：`ci.yml` 里那条 16 包命令的 `go test ./internal/agent/... ` 起于
**127-134 行**（`grep -n` 实测），CI 日志的 `##[group]Run go test ./internal/agent/... \` 与之逐字一致；
只读代理报告写的 `107-113` 是错的（A51⑥ 已登记，我这一步复核支持"127-134"）。

### 4b. 那 10 条 `--- FAIL` 的来处：**我的 Linux 复跑能逐条解释**

CI 日志里 `grep -c -- "--- FAIL"` = **10**（编排者给的数字复核一致），名单与归属：

| CI 里的红 | 我在 Linux 复跑 | 归属 |
|---|---|---|
| `TestWriteGateNotSelectedByPayloadKey` `…PlainLocalWriteNotFlagged` `…EveryPathTargetJudged` `…AllPathsNonSyncStaysExempt` | **同名同锚点复现**（`provenance_test.go:469/541/608/670`，日志文本逐字相同） | 票 55/82 |
| `TestSyncFallbackNotDisarmableByWeakRoot` `TestSyncEnvConfiguredRoots` `TestSyncNormalNewFileWriteNotFlagged` `TestSyncDotDotTailFailsClosed` | **同名同锚点复现**（`syncdirs_test.go:118/121/168/288/356`） | 票 55/82 |
| `TestSpillContainmentByDirectoryListing`（`internal/agent`） | 我在 docker 里 `-v -run` 复现：失败文本是 `tool-output-..\\..\\..\\..\\CONTROL-escape.txt` **反斜杠夹具** | 票 81 |
| `TestArtifactsContainmentByDirectoryListing`（`internal/memory`） | 同上复现：`data/artifacts/nested\canary-separator.txt` | 票 81 |

CI 该步的包级结论行：`ok internal/tools 9.285s`（我：`ok 9.671s`）、
`FAIL internal/risk 1.968s`（我：2.859 s，同名 8 条）、`FAIL internal/agent`、`FAIL internal/memory`；
**`test timed out` 0 次、`panic:` 0 次** ⇒ 票 75 的 19 条 + 超时在 runner 上确实消失。

### 4c. 三条**新发现**（票面/编排者的话与实测不符，照实报）

1. **`lint=failure` 就是票 75 的账，不是"与票 75 无关"。**
   CI 日志 `gofumpt would reformat:` 只点了**一个文件**：
   `internal/risk/provenance_syncdirs_windows_test.go` —— 它正是票 75 在 **`ea01b9f`** 里新增的文件。
   我在这台机器上用 `gofumpt v0.7.0` 对 `17efc2c` 快照跑 CI 的同一条命令
   `gofumpt -l . tools/d22scan tools/mockllm` → 输出**逐字只有这一行**（其余文件全清）；
   再对 `6ce43c7`（票 70 的格式化 sweep）的同文件跑 → **零输出**。
   `git log -- <file>` = `ea01b9f`（票 75 建）+ `6ce43c7`（票 70 修）。
   ⇒ 归因：**票 75 引入了这次 lint 红，票 70 的 commit 替它清的账**；这正是票面 Rules 里
   A40⑤ 警告过的模式（"加宽的门抓到既有违规就得同批修掉"）。
   对签字的影响：**不改判 AC#5**（既没弱化断言也没 tag 掉测试，见 4e），但 AC#6 那句
   "红因与票 75 无关"要更正。
2. **`Environment fork assertion` 步骤是 `skipped`，不是 success。**
   它在 `Portable package tests` 之后且**没有 `if: always()`** ⇒ 只要 test-core 的那一步红，
   这道"防 `-run` 空匹配"的门**根本不执行**。这条门是票 71 AC#3 专门为堵"空匹配假绿"建的，
   今天两次 run（`35558750456`、`35560673885`）里它都是 skipped。
   ⇒ 建议登记（不是票 75 的账，但票 75 的 run 取证把它照出来了）。
3. **AC#6 指的"下一次 dev push"已经有 3 次了**，不止编排者读的那一次：
   `830621b`→run 35560137482、`294f35d`→35560544295、`a0aa44b`→**35560673885**（04:21Z，最新，completed/failure）。
   最新一次的 `test-core` 步骤仍 `failure`，`--- FAIL` = **10**：**同样那 8 条 risk** +
   `TestArtifactsContainmentByDirectoryListing` + **`TestArtifactsLiteralBackslashKeepsItsAsymmetry`**（memory，票 81 家族），
   `internal/agent` 这次 `ok`（spill 那条已被修），`internal/tools` **`ok 9.231s`**；
   `lint` 的 gofumpt 步骤已绿、改由 **`staticcheck`** 失败。
   ⇒ **"test-core 全绿"到今天为止仍未达成**，且残留没有一条是票 75 的形状账；
   票 70 AC#2 若以"test-core 绿"为判据，**至今不成立**。

### 4d. 假绿四项点名（第 4 组）

- **步骤被静默跳过**：`Environment fork assertion = skipped`（见 4c-2），已点名。
- `--- SKIP` 当 ok：CI 那 8 条是真 `--- FAIL`，不是 SKIP；我复现时同名同锚点。
- `-run` 空匹配：我跑票 81 那两条时用了 `-run`，**专门 grep 过 `no tests to run` 与 `testing: warning`**，
  且 `=== RUN` 有 2 条、两条都有失败文本 ⇒ 非空匹配。
- `-count=N` 没核对：3b 的 Windows `-count=2` 里 `TestResolvePerCallBudget` 出现**两次**
  （8.98 s 与 1.40 s 各一轮），证明确实跑了 2 轮而不是被 `-count` 骗过。

### 4e. AC#5 的独立测量（不采信实现者的自述）

对 `17efc2c` 这棵快照树本身取数（不是对 diff）：
```
git grep -c "t\.Skip(" 17efc2c -- internal/risk/pathshape_portable_test.go \
  internal/tools/pathshape_portable_test.go internal/risk/provenance_syncdirs_other_test.go \
  internal/risk/provenance_syncdirs_windows_test.go     ->  rc=1，无输出 = **0 处 t.Skip**
```
（`grep "t.Skip"` 命中的 3 行**全是注释**："Written as a guard rather than a t.Skip"、
"guard, not t.Skip: see …" ⇒ 平台差异写成 early-return 守卫，SKIP 台账不增格。）
`ed74595` 曾加过 2 处真 `t.Skip(`，`493b4e7` 又把它们改成守卫 ⇒ **净 0**，AC#5 的方向对。
build tag 只出现在两个**新增**文件上（`provenance_syncdirs_other_test.go` = `//go:build !windows`、
`provenance_syncdirs_windows_test.go` = `//go:build windows`），**没有一条既有测试被 tag 走**；
`pathshape_portable_test.go` 两个文件**无 tag**（第 3 组已在 Windows 真跑证明）。

### 1d. 假绿四项点名检查（第 1 组）

- `--- SKIP` 当 ok？**3 条 SKIP，逐条点名**：`internal/tools` 的 `TestD34WriteMatrix`
  （`fs_write_test.go:151`，"this machine has a single volume"）与
  `TestCrossVolumeMoveStopsWithTwoCopiesOnLateStop`（`fs_write_test.go:140`，
  "no volume identity for the reference path on this platform"）；
  `internal/risk` 的 `TestSyncRegistryProbeLive`（P12 registry 探针本体，POSIX 无注册表）。
  三条都是**卷标号/注册表这类平台 API 限制**，不是路径形状被藏起来。**我自己看到的反证**：
  第 2 组的变异树里 `TestCrossVolumeMoveStopsWithTwoCopiesOnLateStop` 与 `TestSyncRegistryProbeLive`
  **照样 SKIP**（变异树 tools SKIP=1 / risk SKIP=1），⇒ 它们不是这次修复造出来的；
  同一条 `TestD34WriteMatrix` 在变异树上从 SKIP 变成 **FAIL**，说明这三格 SKIP 没有藏形状问题。
  ⇒ 不构成假绿。
- `-run` 空匹配：**本组没用 `-run`**，整包全跑。
- `-count=N` 没核对 N 倍：本组是 `-count=1`，且 RUN 数与去重后的测试名 + 子项数量自洽（1a 表）。
- 步骤被静默跳过：`go test` 的两行包级结论都在（`ok`/`FAIL`），无 "no test files"。

### 1e. 第 1 组结论

票 75 的**主断言成立**：`17efc2c` 上 Linux 的 `internal/tools` **19 FAIL + 超时 → 0 FAIL、`ok 9.671s`**；
`internal/risk` 的红从基线的 12/13 条降到 **8 条**，且这 8 条**逐条**是票 55 家族的 precondition 红，
没有一条是形状红。**PASS**。
