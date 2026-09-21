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

---

---
---

# 第二验收会话记录（session B，独立复跑；与上文 session A 是**同一票的双派发**）

第二个 `acceptor-ticket75` 会话（下称 B）于 12:0x–12:5x 被独立派出发起同样的四组复现。
**本节全部数字都是 B 自己测的，凡引用上文 session A 处均明示。** 编排者要的一句话结论：
**B 独立到达与 A 一致的四组裁决；A 的三条"新发现"（4c）B 逐条亲验属实。**

## B0. 事故披露：两会话共用 `/tmp/wisp75`，A 的变异曾污染 B 的第一次取树

- B 的初始 `git archive 17efc2c | tar -x -C /tmp/wisp75`（约 12:10）第一轮读数是**修复后**形态
  （tools RC=0、risk 8 红；容器内有 sha256 打印为证）。随后 B 发现同树非 `-v`（8 红）与 `-v`（34 红）
  签名矛盾：后者打出 `<cwd>/\tmp\…` 的**修复前**特征串，且 `internal/tools/` 下多出
  `tmpTestSensitiveFileIsDenied…001` 等目录残渣——即 **A 组 2 的全树变异落在了 B 的快照上**
  （两会话被派发了同一目录名，任务书都写"禁在仓内建 worktree、在 `/tmp/wisp75` 做"）。
- B 的处置：重抽 + 与旧树 `diff -rq`（唯一源码差异 = 被换的 `pathresolver.go`），此后**每跑必带哈希钉桩**：
  修复树 `internal/risk/pathresolver.go` sha256 = `c54d52ff…`（与 `git cat-file 17efc2c:…` 逐字节一致）；
  修复前基线树（`536998b` = `27c6fe5^`）sha256 = `c51ebe85…`。Windows 侧复跑迁到私有目录 `wisp75v2`。
- **程序性建议（两会话共同署名交给编排者）**：验收派发的快照目录必须带会话后缀；
  本票同内容双发一次，纯浪费两台时区同一台机器。

## B1. 第 1 组（Linux docker `golang:1.27`，B 自己的数）

命令 = `go test ./internal/tools/... -count=1` / `go test ./internal/risk/... -count=1`
（`ci.yml:127-134` 的原样子集；两包实测无子包；容器每次 `cp -a /src/. /work/` 后在容器本地 FS 跑——
bind+9p 曾出 `compile: chdir /src` 竞态一次，换法绕开，源码哈希在跑前打印）。

| 树 | 包 | RC | 包级 | FAIL/SKIP | `-v` 台账 |
|---|---|---|---|---|---|
| 17efc2c | tools | **0** | `ok 9.428s` | 0 / 2 | RUN 56 = 42 top + 12 sub + 0 + 2 |
| 17efc2c | risk | **1** | `FAIL 2.393s` | **8** / 1 | RUN 113 = 59 top + 45 sub + 8 + 1 |
| **536998b（27c6fe5 的直接父，修复前）** | risk+tools 同一条命令 | **1** | — | **33 行** = risk 12 + tools 21（17 顶层 + 4 子项） | — |

- tools 修复前名单含票 75 声称的业务红家族全部成员（`TestFSRead*`/`TestFSList*`/`TestFSMoveSameVolume`/
  `TestFSTrashGoesToTheRecycleBin`/`TestAtomicWriteKillsMidWrite`(+3 子项)/`TestDeleteEnabledRegistersItAsL2`(+1 子项)/
  `TestSensitiveFileIsDeniedNotEscalated`/`TestToolCallRowsAreComplete`/`TestOverwriteDetectionFollowsTheCanonicalPath`/
  `TestD34WriteMatrix`/`TestEmptyAllowlistAuthorizesNothing`/`TestDeclaredRiskIsOnlyAFloor` 等）。
- **超时带名字复现（AC#3 的 B 侧独立见证）**：
  `panic: test timed out after 10m0s / running tests: TestVetoInsideTheWindowWritesNothing (4m49s)`
  （栈停 `wiring_test.go:162 → bridge.go:256`，即 C18 的 300 s 审批窗）；修复后同名两条在 B 的 `-v` 日志
  `TestL1WriteGoesThroughTheRealBlockWindow (3.01s) PASS`、`TestVetoInsideTheWindowWritesNothing (0.00s) PASS`。
- **8 ⊂ 12 集合关系（B 独有读法）**：536998b 的 risk 12 条 = 修复后 8 条**逐字同名** +
  `TestSyncWriteNegative`/`TestSyncFixtureFallbackAndMatch`/`TestSyncSuspectFallbackIsComponentBounded`/
  `TestSyncSuspectFallbackWhenUndetectable` 四条形状族 ⇒ 票 75 治掉 4 条、残留 0 新增；
  残留 8 条的归族因果链（`syncdirs.go:141 c.canonical ← pathresolver.go:84 res.Resolved ←
  pathresolver_other.go:10 resolveHandle≡false` + `syncdirs_other.go:20 registryProbe≡nil, DEFERRED(P12-macos) (ticket 55)`）
  B 独立读到同一物理位置。**票 55 家，不计 75。**
- SKIP 点名（B 侧）：tools `TestD34WriteMatrix`、`TestCrossVolumeMove…LateStop`
  （`fs_write_test.go:172/677` 打 "no volume identity for the reference path on this platform"）、
  risk `TestSyncRegistryProbeLive`（POSIX 无注册表；本机 Windows 也 SKIP，非新增遮羞）。
- 与 A/实现者读数差：B 的基线取 `536998b`（代码修复的逐字前像），A 取 `131f722`（更早、含 docs 检查点），
  实现者取 `f1033e1/131f722`；risk 12/13、tools 18/21 的差全部由 `-timeout` 闹钟落点与
  `TestResolvePerCallBudget`（见 B3）解释，族与机理三侧一致。**票面标题"19 条/47 条"是 `24a66b6` 取证数，
  B 未复现那个 sha 的读数——它不是任何人的证据档，只是个题面引用。**
- 日志：`wisp75logs/{t2,t2v,r2,r2v,prefix-base}.log`。

## B2. 第 2 组（变异，B 的形态与 A 不同：同容器 A/B 对照）

锚点 = B 快照实测行 167-171（守卫 168-170 + 折叠 171）；改动 = `sed '168,170d'` 删守卫三行。
落地证明在**跑测试的同一条命令链**里：容器内对 `/work` 的副本 `grep -c 'not a UNC spelling'`
先 =1（CLEAN 态，`.orig` 还原后）跑一遍，再 =0（MUTANT 态，覆盖成 B 的变异树）跑一遍，
两次 `-v -count=1 -run '<票 75 新 7 条>'` **同命令同容器**：

| 态 | RC | `=== RUN` | SKIP | 结果 |
|---|---|---|---|---|
| clean | 0 | 19（7 顶层+12 子项） | 0 | 全 PASS |
| mutant | **1** | **19（同名集，编译通过 ⇒ 行为变异成立）** | 0 | **7 顶层全红**（19 行 FAIL 含子项） |

红名单：`TestResolveCanonicalIsPlatformShaped`、`TestResolveKeepsSeparatorsOutOfPosixNames`、
`TestDoubleSlashSpellingIsNotUNC`、`TestTierAAnchorsMatchNativeSpelling`、
`TestOverrideKeyKeepsPosixBackslashesDistinct`、`TestSyncAncestorWalkVerifiesNativeChain`、
`TestExfilSyncWritePosixSpellingInvariant`；失败原文即票面串（`…\tmp\…: cannot be statted`、
`CROSS-FILE BLEED (fail-open)`）。还原对拍 = 宿主机 `sha256sum` 回到 `c54d52ff…` 且 `.orig` 无残留。
**为什么 B 的变异只打红 risk 就够**：`tools/paths.go:56` 的 `Canonicalize` 委托 `risk.Resolve`，
A 组 2 已证同一变异波及 tools（+21 红 + 600 s panic）；B 不重复全树跑，机理上两把锁已扣死。
日志 `wisp75logs/{mut-clean,mut-red}.log`。**PASS。**

## B3. 第 3 组（Windows 本机，B 的私有树 `wisp75v2`，跑前哈希 `c54d52ff…`）

```
go test -v ./internal/risk/ -count=1 -run 'TestOverrideKeyKeepsPosixBackslashesDistinct$|TestExfilSyncWriteWindowsSpellingInvariant$'
RC=0；=== RUN=9（2 顶层+7 子项）、PASS=9、FAIL=0、SKIP=0、'no tests to run' 0 命中
```
前者**无 build tag** 且在 Windows 真跑（钉 Windows 侧"两拼写折一键"，非空转）；后者 7 子项逐条 PASS。
**AC#2 Windows 半边（B 的完整台账）**：`go test -v ./internal/risk/ ./internal/tools/ -count=2`
= **RC 0**、`ok risk 7.617s` / `ok tools 28.259s`、**RUN 482 = 290 top + 190 sub + 0 FAIL + 2 SKIP**
（SKIP 去重 1 个 = `TestSyncRegistryProbeLive`，与实现者日志同既有）。
**负载证词（与 A 5b 互证）**：B 更早一次同命令在 CPU 99% 时 rc=1，唯一红 = `TestResolvePerCallBudget`
（2.276/2.513 ms vs 1 ms）；单跑三次 0.843 PASS / 1.503 FAIL / 0.418 PASS ⇒ 绝对值墙钟门，负载敏感，
**非票 75 账**；B 独立同意 A 的建议（该门单独立票改相对量/隔离）。
日志 `wisp75logs/{win-newtests-v2,win-count2,win-count2b}.log`。**PASS。**

## B4. 第 4 组（CI 对账，B 只读 gh；含对 A 三条新发现的**独立证实**）

`gh run view 35558750456`：**event=push、branch=dev、ci #125、headSha=17efc2c**、completed（非 cancelled）。
jobs：test-core=failure / lint=failure / test-windows=success / slo-smoke=success / slo-full=success。
job **106207404142** 步骤级：**`Portable package tests (…)` = failure**（`ci.yml:127-134` 实测在该行号；
107-113 引用腐坏=A51⑥，与读数无关）；`Environment fork assertion` = **skipped**。
该步 `--- FAIL` = **10**（`gh run view --job … --log`，存档 `wisp75logs/ci-testcore.log`）：
8 条 = B1 的 risk 名单**逐字同名**（票 55）；
`TestSpillContainmentByDirectoryListing`（agent）与 `TestArtifactsContainmentByDirectoryListing`（memory）
= **票 81**（失败文本 `tool-output-..\..\..\..\CONTROL-escape.txt` 的 Windows 夹具钉在 POSIX 树上，B 在 CI 日志原文里读到）。
包级：`ok internal/tools 9.285s`；`FAIL risk 1.968s / agent 0.851s / memory 10.822s`，其余 17 包 ok。

**B 对 A §4c 三条的复核结果（全部证实）：**
1. **`lint=failure` 就是票 75 的账**——B 用 CI 钉的 `gofumpt v0.7.0` 在容器里对 `17efc2c` 的
   `internal/risk/provenance_syncdirs_windows_test.go`（票 75 `ea01b9f` 新建）跑 `-l` → **点名该文件**；
   对 `6ce43c7`（票 70 sweep，commit 标题自证"the one file"）的同文件 → **零输出**。
   ⇒ 票 75 引入文件时未过格式门、由票 70 清账：**违的是票面 Rules 的 A40⑤"同批该修掉"，
   不是 AC#5 的字面（没有弱化断言/tag 掉测试）**。编排者任务书里"lint 红因与票 75 无关"一句**要更正**。
2. **`Environment fork assertion` skipped 是结构性空窗**：`ci.yml:136-142` 该步无 `if: always()`，
   排在 failure 步之后 ⇒ 永不执行（两条 run 皆 skipped）。票 71 建的防"空匹配假绿"门自己成了假绿——登记。
3. **后续 push 已到 `a0aa44b`**：run **35560673885**（push、completed、failure）——
   test-core 步骤 `Portable package tests` 仍 failure，`--- FAIL` = **10**：
   同样 8 条 risk + `TestArtifactsContainmentByDirectoryListing` + `TestArtifactsLiteralBackslashKeepsItsAsymmetry`
   （memory，票 81 家）；`internal/agent` 本次 **ok**（spill 那条被票 81 修了）、
   `ok internal/tools 9.231s`；lint 的 gofumpt 步已绿、红转 **staticcheck**。
   ⇒ **"test-core 全绿"截至 12:5x 仍未达成，且残留无一为票 75 形状账**——
   票 70 AC#2 若以"test-core 绿"为判据，至今不成立（与 A 一致）。

## B5. B 的 AC 裁决表（独立，不抄 A）

| AC | 票面勾框 | B 裁决 | B 证据 |
|---|---|---|---|
| #1 改前基线留档 | [x] | **PASS**（B 的数）：536998b scoped `--- FAIL` 33 行（risk 12 + tools 21 行）+ 带名字 600 s panic；票面"47" 是 24a66b6 取证，B 未复现该 sha（见 B6-2） | `prefix-base.log` |
| #2 两侧门 | [x] | **PASS**：Linux tools RC=0；Windows 全包 `-count=2` RC=0（482 全对账）；预算门负载假红 1 次已复验推翻 | `t2.log`、`win-count2b.log` |
| #3 超时具名 | [x] | **PASS**：B 亲眼复现 panic 名与 4m49s 等待，修复后 3.01s/0.00s | `prefix-base.log:94`、`t2v.log` |
| #4 D22 交接 | [ ] | **等人关**：交接段落三处改动 B 用 `git diff 536998b..17efc2c -- pathresolver.go` 逐字对上（守卫/`sepStr`+`unifySeparators`/`tailExistsBelow` 重拼接，无其它）；缺的是 owner 裁定，不是代理没做完 | diff + 票面 :97-113 |
| #5 不放宽 | [x] | **PASS（附更正）**：修复范围 diff 零新增 `t.Skip`（注释反而明写 "guard, not t.Skip"）、无既有测试被 tag 走（新 tag 只在两个**新增**分层文件上）；**但票 75 新文件曾让 lint 红（B4-1），违 A40⑤ 不违 AC#5** | `git diff`、gofumpt A/B |
| #6 runner 证明 | [ ] | **等人关**：run/job/step 三级读数 + 10 条逐条来处 B 全部独立取到（见 B4）；票面要求的回执行措辞应改为"**票 75 的账在 runner 上已清零；test-core 仍因 55/81 红**" | run 35558750456 / 35560673885 |

## B6. B 视角的"未证"清单（B 独有项加粗；A 列的 1–7 B 复核后**全部维持**）

1. **`24a66b6` 上的"4 包 47 条"原始读数两个验收会话都没复现过**——它是票 70 取证转述，
   票 75 的 AC#1 实际复现的是 35/33/31 一族 scoped 数。谁要签"47"这个词，得自己去跑那个 sha。
2. **票 75 修复后 tools 侧票 18/20/73 四组 `_windows_test` 的逐文件 `-count=2`（Progress 日志的 20/12/8/8/2）
   两个会话都没逐文件复跑**；B 只以全包 `-count=2` RC=0 + 两条新仪器覆盖。若编排者要逐文件账，仍在实现者自述里。
3. **POSIX 真 UNC 形态（`//server/share` 挂载点）交回什么** 无覆盖（A#3 说得好，B 复核 `pathshape_portable_test.go` 同意）。
4. **`bOverrides` 生产接线**（A#6/B 同）：`internal/config → Gate` 无调用方，票 80 的地盘，未证。
5. **macOS 零读数**；且 B 补一句：**B 的 Linux 全是 WSL2/Docker Desktop**，与 `ubuntu-latest` runner
   的内核/FS 不同型——墙钟类断言两侧不可比（形状类结论不受影响）。
6. **staticcheck 在 `a0aa44b` 的新红**（B4-3）没人归因过：是不是票 75 的遗产 B 未查（查它需要的上下文在票 70/84）。

## B7. B 动过仓库里的什么

本文件上 B 的全部写入 = 本节 + 下方 A 片段的中继，无其它。
**实况更正（B 不给自己留面子）**：B 的第一次落笔因与 A 的并发写撞车，一度把 A 的
commit 后内容截断；B 已用"HEAD + 本节 + A 的中断期片段（`/tmp/wisp75logs/sessionA_post-commit_fragment_1303.md` 保真备份）"
重组，**A 的字节零丢失**（其组 5/裁决表/未证清单/第 6 节完整在案，组 2-4 出现一次可去重的重复，去留归 A/owner）。
工作树里其余在飞文件（`frontend/`、票 80/81/84、`.scratch`）一个字节未碰。B 的树与日志：
`C:/Users/swq/AppData/Local/Temp/{wisp75,wisp75v,wisp75v2}` + `wisp75logs/{t2,t2v,r2,r2v,win-newtests,win-newtests-v2,win-count2,win-count2b,mut-clean,mut-red,prefix-base,ci-testcore,fixed_windows_test}.log`。
快照 `wisp75` 在 B 变异后已还原并由哈希证明 = `17efc2c` blob；`wisp75v2` 从未被变异过。

---

## 附：session A 于 13:02 快照里的 commit 后增补片段（B 误伤重建时替 A 保真转存；与上文第 2-4 组的重复去留请 A/owner 定夺——第 5 组及其后为本片段独有）


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
		TestVetoInsideTheWindowWritesNothing (4m45s)
```
goroutine dump 里那条 4m45s 的**阻塞点有栈**（同一份 `-json` 日志，`FAIL internal/tools 600.108s`）：
```
goroutine 124 [select]:
  …/internal/agent/approval.(*Gate).PendingApproval(…)  /src/internal/agent/approval/gate.go:473 +0x816
  …/internal/tools.(*Bridge).route(…)                   /src/internal/tools/bridge.go:340 +0x489
  …/internal/tools.(*Bridge).Execute(…)                 /src/internal/tools/bridge.go:256 +0x7dc
  …/internal/tools_test.call(…)                         /src/internal/tools/wiring_test.go:106
  TestVetoInsideTheWindowWritesNothing                  /src/internal/tools/wiring_test.go:162
```
⇒ **4m45s 已经越过 C18 的 300 s 上界**，卡在 `PendingApproval` 的 select 里 ⇒
这条**独立佐证了 A55②/票 84 的判据**（"correlation_id 无对应待审批项 ⇒ 无界阻塞"），
不是转述。改前树 `131f722` 那一跑我也量到同族越界：
`panic: test timed out after 15m0s / running tests: TestLateVetoRendersTheApprovalLayersAppliedStepsReport (4m41s)`（见 5a）。
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

1. **`lint=failure` 就是票 75 的账（**这条不是新发现**：账本 **A52②** 已经登记过同一事实；
   是**本次派发给验收代理的简报**里那句"红因 `gofmt (gofumpt)` 步骤，与票 75 无关"口径错了）。
   我把这条独立复算了一遍：CI 日志 `gofumpt would reformat:` 只点了**一个文件**：
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

---

## 第 5 组（补测）：AC#1 的"改前基线"我自己量了一遍 + 争抢修正

### 5a. 改前树 `131f722`（票 72 落地后、票 75 修复前）在同一 Linux 容器里

```
git archive 131f722 | tar -x -C /tmp/wisp75pre        # 又一棵仓外纯净快照
docker run … -e WISP_ENV=test golang:1.27 bash -c 'cd /src && go test ./internal/risk/ ./internal/tools/ -count=1 -timeout 15m'
PRE_FIX_RC=1
```
| 包 | 改前 `131f722` | 修复后 `17efc2c`（第 1 组） |
|---|---|---|
| `internal/risk` | **13** 条 `--- FAIL`，5.753 s | **8** 条，2.848 s |
| `internal/tools` | **18** 条 `--- FAIL` + `panic: test timed out after 15m0s`，包 **900.299 s** | **0** 条，`ok 9.671s` |
| `--- FAIL` 合计 | **31** | **8** |

**AC#3 的具名等待我自己复现到了**（改前树，逐条）：
```
--- FAIL: TestL1WriteGoesThroughTheRealBlockWindow (300.02s)
--- FAIL: TestVetoInsideTheWindowWritesNothing (300.28s)
--- FAIL: TestLoopPassesDeclaredL1WriteThroughTheGate (10.56s)
panic: test timed out after 15m0s
```
两条各烧满 **300 s** = C18 的"超时一律判拒绝"，第三条被闹钟切断 ⇒
**"哪条测试、哪一次等待"这条解释成立**，不是"它变快了"。
⇒ AC#1/AC#3 **PASS**（票面说"47-ish"，实现者在同 sha 量到 35/18/12；我量到 31/18/13，
差的那 1 条 risk 是 `TestResolvePerCallBudget`，见 5b）。

### 5b. 3b 那条 Windows `TestResolvePerCallBudget` 红：**确认为 CPU 争抢，已推翻自己的读数**

机器安静后（`docker ps` 无 golang 容器）复跑，同一快照 `/tmp/wisp75w`：
```
go test ./internal/risk/ -run TestResolvePerCallBudget -count=2 -v
=== RUN × 2   →  0.474 ms/op、0.532 ms/op（预算 1 ms）    --- PASS × 2    ok 4.210s    rc=0
```
Linux 安静复跑（`go test ./internal/risk/ -count=1 -v`）：
```
LINUX_RC=1   顶层 --- PASS = 59   --- FAIL = 8（名字与第 1 组逐字相同）   --- SKIP = 1（TestSyncRegistryProbeLive）
C26 Resolve: 790 ns/op = 0.001 ms/op（budget 1 ms，1543290 samples）  ->  PASS
```
⇒ 3b 的 2 条红**不是缺陷**：同一个门在争抢时 1.6 ms/op、安静时 0.47–0.53 ms/op。
**这条也说明票 75 证据档里"全量并跑 13 条 vs 单跑 12 条"的分歧根因就是它**，
与形状无关；`TestResolvePerCallBudget` 是**绝对值 wall-clock 预算**，在共享 runner/本机有邻居时
天然不稳（我第 2 组变异读数里那条 2.483 ms/op 同理，已剥离）。**建议单独立票改成相对量或隔离**（不是票 75 的账）。

---

## 裁决表：票 75 每条 AC → 我的裁决 → 证据档

| AC | 票面要求（一句） | 我的裁决 | 证据（本文件节）+ 数字 |
|---|---|---|---|
| **AC#1** 改前用 CI 逐字命令在 docker 复现并记数 | **PASS（我自己量到，非转述）** | 5a：`131f722` 树 Linux `go test ./internal/risk/ ./internal/tools/ -count=1` → `--- FAIL` **31**（risk 13 / tools 18 + `panic: test timed out after 15m0s`，包 900.299 s），rc=1。命令形状取自 `ci.yml:127-134`（我 grep 复核行号，107-113 是错的） |
| **AC#2** tools 在 `golang:1.27` 绿 **且** Windows `-count=2` 守住绿 | **PASS** | Linux 侧 1a/1b/5b：`ok 9.671s`、RUN 56 / PASS 54 / **FAIL 0** / SKIP 2；改后 `17efc2c` 与 CI run `35558750456`（`ok internal/tools 9.285s`）、最新 push `35560673885`（`ok 9.231s`）三处一致。Windows 侧 3b/5b：tools `-count=2` **`ok 43.819s`、0 FAIL、0 SKIP**；risk 的 2 条红是 `TestResolvePerCallBudget` 争抢，安静复跑 `-count=2` **rc=0、PASS×2（0.474/0.532 ms/op）** |
| **AC#3** 600 s 超时具名解释 | **PASS** | 5a：改前树 `TestL1WriteGoesThroughTheRealBlockWindow (300.02s)` + `TestVetoInsideTheWindowWritesNothing (300.28s)`（各满 300 s = C18 一律判拒绝）+ 第三条被 15m 闹钟切；2c：变异树上同一条 300.04 s 复现、`panic: test timed out after 10m0s` 挂在 `TestVetoInsideTheWindowWritesNothing`；修复后整包 9.671 s |
| **AC#4**（D22 门，编排者已用 R17 关闭） | **维持关闭，且证据支持其事实前提** | 2d/2f：我**逐字读到**契约读数的物理后果——不改这个冻结文件，POSIX 上 `Resolve` 交回的就是 `/src/internal/risk/\tmp\…`（OS 开不了、也不再绝对）；变异（只回退守卫）⇒ 两包合计红由 **8 → 56**（risk 8→35、tools 0→21）。"实现不符合冻结契约"这一判定站得住。**次序偏离（先改后提案）票面已自记，我不重复扣分** |
| **AC#5** 不弱化断言、不为变绿 tag 掉测试 | **PASS（附一条更正）** | 4e：`17efc2c` 上票 75 四个仪器文件 `t.Skip(` **= 0**（3 处命中全是注释）、既有测试**零**被新 tag 走、无阈值改动。更正的那条：**票 75 自己新增的 `provenance_syncdirs_windows_test.go` 让 `lint` 在 `17efc2c` 红**（gofumpt 唯一被点名的文件，我用 v0.7.0 复现），由票 70 的 `6ce43c7` 清账 ⇒ 属"同批该修掉"的疏漏，**不是 AC#5 的字面违规**，但票面 4c-1 那句话要改 |
| **AC#6**（仍开着）runner 可见证据 | **证据充分可关，但措辞要改；"test-core 绿"没达成** | 4a/4b/4c-3：run `35558750456`（job `106207404142`）步骤 `Portable package tests` = **failure**、`--- FAIL` **10** = 我在 Linux 逐字复现的 8 条 risk（票 55/82）+ 2 条 Windows 形状夹具（我在 docker 里单独复现，票 81）；`internal/tools` 19+超时 → `ok 9.285s`。最新 push `35560673885` 同样 10 条（agent 那条已被票 81 修、memory 多 1 条）。**回执那行该写成"票 75 的账已清零，test-core 仍因票 55/81/82 红"**，不能写成"全绿" |

## 我认为票 75 至今**无人证明**的部分（明写"未证"）

1. **POSIX 上"根外目标不得被标"这一半从未被执行过**：
   `internal/risk/provenance_syncdirs_other_test.go` 的反向腿是**休眠断言**
   （`SyncDetectionComplete()` 为假时只 `t.Log`）。我在 Linux 全量跑里看到它跑了但没有判据 ⇒
   **"POSIX 上同步写门既不漏判也不误判"未证**。票面自己写了这点（不糊票 55），我确认它仍然空。
2. **`TestExfilSyncWriteWindowsSpellingInvariant` 的"根外 3 种拼写不升嫌疑"腿只在 Windows 成立**：
   第 3 组我确实在 Windows 上跑绿了它，但**同一性质在 Linux 上没有任何在判的断言**（见 1）。
3. **`normalizeLocalUNC` 的 `//` 前缀分支在 POSIX 上是死代码这件事，只用一条测试钉住**
   （`TestDoubleSlashSpellingIsNotUNC`，我在变异组看到它会红 ⇒ 仪器承重）。
   但"**POSIX 上真的 UNC/`//` 开头的路径**（如 `//server/share` 形式的 SMB 挂载点）交回什么"
   没有任何用例覆盖 ⇒ 未证（低风险，但别说成"覆盖了 UNC 全分支"）。
4. **`internal/risk` 的 POSIX 全绿**未证（8 条红 = 票 55/82）；**`test-core` 全绿**未证（今天三个 push 全红）。
   票 75 若被理解成"修完 test-core 就绿"，那是**超出证据的说法**。
5. **真实 macOS 上的形状**未证：本轮全部 POSIX 测量在 Linux（WSL2 内核 + `golang:1.27` 容器）；
   `pathresolver_other.go`/`syncdirs_other.go` 的 `!windows` 分支在 darwin 上没有一次真机/真容器读数。
6. **`bOverrides` 无生产调用方**这条票面自记的新缺陷候选：我只复核了它在测试里被驱动
   （第 2 组变异把它打红，说明"跨文件豁免串用"的钉是活的），
   **`internal/config` 的 `BlacklistOverrides` 是否真被喂进 `Gate`** 我未查、票 75 也未证 ⇒ 归票 80 那条"死线"。
7. **CI 与本机的一致性是"命令逐字 + 容器同版本"级别的**，不是"同一台 runner"级别：
   我的 Linux 是 Docker Desktop/WSL2（`Linux 6.6.114.1-microsoft-standard-WSL2`），
   与 `ubuntu-latest` runner 的 FS/调度不同 ⇒ **wall-clock 类断言（5b）在两侧不可比**，
   这一点我按"数字难看也照实报"处理，没有拿它当豁免。

## 我动过仓库里的什么

只写了本文件：`docs/evidence/s1/75-independent-verification.md`（新文件）。
未 commit、未 push、未 checkout、未 stash、未 `git worktree`；四个在飞代理的文件一个没碰
（全部测量在 `C:/Users/swq/AppData/Local/Temp/wisp75{,w,pre}` 三棵 `git archive` 快照里做，
变异过的那棵已用双 oracle 证明还原，见 2e）。
日志与快照留在 `/tmp/wisp75-logs/`：`g1.json`/`g1b.txt`/`g2.json`/`g3.txt`/`g4b.txt`/`g5.txt`/
`g6_pre.txt`/`g7_win.txt`/`g8_linux_risk_v.txt`/`ci_testcore.log`/`pathresolver.go.orig`。

---

## 第 6 节：编排者在我跑到一半时已结票（`dc91999`，12:43:17）—— 三件要回补的事

1. **`dc91999` 提的是我的半成品。** 那个 commit 把本文件**当时的 359 行版本**（内容止于 4e）
   连同票面改名 `-done` + "六框全绿"一起落了地。我在这之后写的
   **第 5 组（AC#1 改前基线我自己量：`131f722` 上 31 条红、tools 18 条 + 15m 超时、
   两条各满 300 s 的用例具名）+ 第 6 节 + 裁决表 + "无人证明"清单**
   今天**只在磁盘上，不在任何 commit 里**。⇒ **需要编排者再 commit 一次本文件**（我不 commit，硬规矩）。
   渐进写的用意正在这里：别让后半段跟着代理一起死。
2. **我的六条裁决与"六框全绿"不冲突**（AC#1/2/3/5 = PASS；AC#4 = 维持关闭，且我的变异组给出它
   依赖的事实前提；AC#6 = 按"本票的账清零"可关，**但不能被读成"test-core 绿"**——今天最新的
   `35560673885` 仍 10 条红）。两处要改口径：
   (a) 派发简报里"lint 红与票 75 无关"这句**与账本 A52② 自相矛盾**，A52② 才是对的
   （红在票 75 新加的 `provenance_syncdirs_windows_test.go`），我的独立复算支持 A52②；
   (b) **新发现（账本里没有）**：`Environment fork assertion` 步骤在 `35558750456` 与
   `35560673885` 两次 run 里都是 **`skipped`**——它排在 `Portable package tests` 之后且**没有
   `if: always()`** ⇒ 只要 test-core 那一步红，票 71 AC#3 专门建的"防 `-run` 空匹配"门**从不执行**。
   这不是票 75 的账，但它属"步骤被静默跳过"这一类假绿，建议并给票 70 的 CI 面或单开一票。
3. **票 84 的判据我这边有一份独立栈证据**（见 2c 的 goroutine dump：`PendingApproval`
   @ `internal/agent/approval/gate.go:473`，同一用例已跑到 **4m45s**）+ 改前树的 **4m41s**（5a）。
   与 `verify75` 的 288 s 观察**不同机器、不同树、同一函数** ⇒ 票 84 的 AC#1"两侧都量"
   现在 Linux 侧有两份，**Windows 侧仍 0 份**。
   顺带把 AC#3 的解释再拧准一格：**600 s 不全是"300 s × 2 的有界窗口"**，
   其中至少一次是 `PendingApproval` 的**无界**阻塞 ⇒ 票 75 面上"各吃 300 s"这句话
   对它自己没影响（形状根因不变），但**别拿它当"等待都有上界"的证据**。

---

# session A 最终追加（13:1x，覆盖式声明）

A 的最后几段（**第 2 组补 / 第 5 组（A 补测）/ 裁决表（A）/ 未证 / 我动过（A）/ 第 6 节（A）**）
是 A 的**最终读数**，与前面 B 替 A 转存的 13:0x 片段内容重叠：**重叠处以本段为准**。
本段多了：goroutine 栈证据（`PendingApproval` @ gate.go:473，4m45s 越过 300 s 上界）、
`8 → 56` 红数口径、`lint` 那条对 A52② 的归因让位、以及"写战 + 撞目录"两条程序性建议。
**B0–B7 一字未改；A 不 commit，请编排者落一次 commit。**

## 第 2 组补（A 追加）：变异日志里的 goroutine 栈 —— 票 84 判据的独立证据

2c 的 `panic: test timed out after 10m0s` 挂着的用例，时长与阻塞点都在同一份 `-json` 日志里：
```
panic: test timed out after 10m0s
	running tests:
		TestVetoInsideTheWindowWritesNothing (4m45s)
…
goroutine 124 [select]:
  …/internal/agent/approval.(*Gate).PendingApproval(…)  /src/internal/agent/approval/gate.go:473 +0x816
  …/internal/tools.(*Bridge).route(…)                   /src/internal/tools/bridge.go:340 +0x489
  …/internal/tools.(*Bridge).Execute(…)                 /src/internal/tools/bridge.go:256 +0x7dc
  …/internal/tools_test.call(…)                         /src/internal/tools/wiring_test.go:106
  TestVetoInsideTheWindowWritesNothing                  /src/internal/tools/wiring_test.go:162
FAIL	github.com/CarlosShao/wisp/internal/tools	600.108s
```
⇒ **4m45s 已越过 C18 的 300 s 上界**，且停在 `PendingApproval` 的 select 上。
这条**独立佐证 A55②/票 84 的判据**（"`correlation_id` 无对应待审批项 ⇒ 无界阻塞"），不是转述。
⚠ 与 B1 的措辞分歧：B 把同一现场写成"栈停 `wiring_test.go:162 → bridge.go:256`，即 C18 的 300 s 审批窗"。
**栈再往上一层才是真相：窗口有 300 s 上界，`PendingApproval` 没有。**
⇒ 顺带把 AC#3 的解释拧准一格：600 s **不全是"300 s × 2 的有界等待"**，其中至少一次是无界阻塞 ⇒
**别拿票 75 的这段归因当"审批等待都有上界"的证据**（对票 75 自己的形状结论无影响）。

---

## 第 5 组（A 补测）：AC#1 的"改前基线"我自己量了一遍 + 争抢修正

### 5a. 改前树 `131f722`（票 72 落地后、票 75 修复前）在同一 Linux 容器里

```
git archive 131f722 | tar -x -C /tmp/wisp75pre        # 又一棵仓外纯净快照
docker run … -e WISP_ENV=test golang:1.27 bash -c 'cd /src && go test ./internal/risk/ ./internal/tools/ -count=1 -timeout 15m'
PRE_FIX_RC=1
```
| 包 | 改前 `131f722` | 修复后 `17efc2c`（第 1 组） |
|---|---|---|
| `internal/risk` | **13** 条 `--- FAIL`，5.753 s | **8** 条，2.848 s |
| `internal/tools` | **18** 条 `--- FAIL` + `panic: test timed out after 15m0s`，包 **900.299 s** | **0** 条，`ok 9.671s` |
| `--- FAIL` 合计 | **31** | **8** |

AC#3 的具名等待我自己复现到了（改前树，逐条）：
```
--- FAIL: TestL1WriteGoesThroughTheRealBlockWindow (300.02s)
--- FAIL: TestVetoInsideTheWindowWritesNothing (300.28s)
--- FAIL: TestLoopPassesDeclaredL1WriteThroughTheGate (10.56s)
panic: test timed out after 15m0s
	running tests:
		TestLateVetoRendersTheApprovalLayersAppliedStepsReport (4m41s)
```
两条各烧满 300 s = C18 一律判拒绝；**被 15 分钟闹钟切掉的那条已跑到 4m41s**（同 2c 补：越过上界）。
⇒ AC#1/AC#3 **PASS**。（票面标题的"47"是 `24a66b6` 的取证数，**我和 B 都没复现过那个 sha**。）

### 5b. 3b 那条 Windows `TestResolvePerCallBudget` 红：**确认为 CPU 争抢，我自己推翻自己的读数**

机器安静后（`docker ps` 无 golang 容器）复跑，同一快照 `/tmp/wisp75w`：
```
go test ./internal/risk/ -run TestResolvePerCallBudget -count=2 -v
=== RUN × 2  →  0.474 ms/op、0.532 ms/op（预算 1 ms）   --- PASS × 2   ok 4.210s   rc=0
```
Linux 安静复跑（`go test ./internal/risk/ -count=1 -v`）：
```
LINUX_RC=1   顶层 --- PASS = 59   --- FAIL = 8（与第 1 组逐字同名）   --- SKIP = 1（TestSyncRegistryProbeLive）
C26 Resolve: 790 ns/op = 0.001 ms/op（budget 1 ms，1543290 samples） ->  PASS
```
⇒ 3b 的 2 条红**不是缺陷**：同一个门争抢时 1.6 ms/op、安静时 0.47–0.53 ms/op。
**票 75 证据档里"全量并跑 13 条 vs 单跑 12 条"的分歧根因就是它**，与形状无关；
第 2 组变异读数里那条 2.483 ms/op 同理，已从证据里剥离。
⇒ **建议单独立票把这条绝对值墙钟门改成相对量或隔离**（不是票 75 的账）。B3 的"负载证词"
（2.276/2.513 ms 红，单跑三次 0.843 PASS / 1.503 FAIL / 0.418 PASS）与此**同结论、独立取得**。

---

## 裁决表（A）：票 75 每条 AC → 我的裁决 → 证据

| AC | 票面要求（一句） | 我的裁决 | 证据（本文件节）+ 数字 |
|---|---|---|---|
| **AC#1** 改前用 CI 逐字命令在 docker 复现并记数 | **PASS（我自己量到，非转述）** | 5a：`131f722` 树 Linux `go test ./internal/risk/ ./internal/tools/ -count=1` → `--- FAIL` **31**（risk 13 / tools 18 + `panic: test timed out after 15m0s`，包 900.299 s），rc=1。命令形状取自 `ci.yml:127-134`（`grep -n` 复核；"107-113"是腐坏引用） |
| **AC#2** tools 在 `golang:1.27` 绿 **且** Windows `-count=2` 守住绿 | **PASS** | Linux：1a/1b/5b `ok 9.671s`、RUN 56 / PASS 54 / **FAIL 0** / SKIP 2；与 run `35558750456`（`ok internal/tools 9.285s`）、`35560673885`（`ok 9.231s`）三处一致。Windows：tools `-count=2` **`ok 43.819s`、0 FAIL、0 SKIP**；risk 的 2 条红是墙钟门争抢，安静复跑 `-count=2` **rc=0**（B3 亦独立给出全包 `-count=2` **RC=0 / RUN 482**） |
| **AC#3** 600 s 超时具名解释 | **PASS（并把口径拧准）** | 5a 改前树两条 300.02/300.28 s + 15m 闹钟；2c+「第 2 组补」变异树 300.04 s 与 `TestVetoInsideTheWindowWritesNothing (4m45s)` 卡在 `PendingApproval`（**越过 300 s 上界**）⇒ 修复后整包 9.671 s。**"各吃 300 s"这句对一半**：另一次是无界阻塞（票 84 的账） |
| **AC#4**（D22 门，编排者已用 R17 关闭） | **维持关闭，且证据支持其事实前提** | 2d/2f：我逐字读到契约读数的物理后果——不改这个冻结文件，POSIX 上 `Resolve` 交回的就是 `/src/internal/risk/\tmp\…`（OS 开不了、也不再绝对）；只回退守卫 ⇒ 两包合计红 **8 → 56**（risk 8→35、tools 0→21）+ 10 分钟超时。"实现不符合已冻结契约"站得住；**先改后提案的次序偏离票面已自记，我不重复扣分** |
| **AC#5** 不弱化断言、不为变绿 tag 掉测试 | **PASS（附一条更正）** | 4e：`17efc2c` 上票 75 四个仪器文件 `t.Skip(` **= 0**（3 处命中全是注释）、既有测试**零**被新 tag 走、无阈值/黄金文件改动。更正的那条：**票 75 新增的 `provenance_syncdirs_windows_test.go` 让 `lint` 在 `17efc2c` 红**（gofumpt 唯一被点名的文件；A52② 已登记，我和 B 各自独立复算）⇒ 违的是票面 Rules 的 **A40⑤"同批该修掉"**，**不是 AC#5 的字面** |
| **AC#6** runner 可见证据 | **按"本票的账清零"可关；**不得读成"test-core 绿"** | 4a/4b/4c-3：run `35558750456` / job `106207404142` / 步骤 `Portable package tests` = **failure**、`--- FAIL` **10** = 我逐字复现的 8 条 risk（票 55/82）+ 2 条 Windows 形状夹具（我在 docker 里单独复现，票 81）；`internal/tools` 由 19+超时 → `ok 9.285s`。最新 push `35560673885` 仍 10 条、无一属票 75。**A58② 已把这条线推到 47→10→8**，与本表一致 |

## 我认为票 75 至今**无人证明**的部分（明写"未证"）

1. **POSIX 上"根外目标不得被标"这一半从未被执行过**：`provenance_syncdirs_other_test.go` 的反向腿是
   **休眠断言**（`SyncDetectionComplete()` 为假时只 `t.Log`）。我在 Linux 全量跑里看到它跑了但没有判据
   ⇒ **"POSIX 上同步写门既不漏判也不误判"未证**。（票面自己写了不糊票 55，我确认它仍然空。）
2. **`TestExfilSyncWriteWindowsSpellingInvariant` 的"根外 3 种拼写不升嫌疑"只在 Windows 成立**：
   第 3 组我跑绿了它，但同一性质在 Linux 上没有任何在判的断言（见 1）。
3. **POSIX 上真的 `//server/share` 形态（SMB 挂载点类）交回什么**：无覆盖。
   `//` 分支在 POSIX 上是死代码这一点只由 `TestDoubleSlashSpellingIsNotUNC` 一条钉住
   （第 2 组看到它**会红** ⇒ 仪器承重），但别把 UNC 全分支说成已覆盖。
4. **`internal/risk` 的 POSIX 全绿**未证（8 条红 = 票 55/82）；**`test-core` 全绿**未证。
   把票 75 理解成"修完 test-core 就绿"是**超出证据的说法**。
5. **真实 macOS 上的形状**未证：本轮全部 POSIX 读数在 Linux（Docker Desktop / WSL2 内核
   `6.6.114.1-microsoft-standard-WSL2`），`//go:build !windows` 分支在 darwin 上**零真机/零真容器读数**。
   ⇒ 与 `ubuntu-latest` runner 不同型：**墙钟类断言两侧不可比**（形状类结论不受影响）。
6. **`bOverrides` 是否真接进生产**未证：我只复核它在测试里被驱动（第 2 组变异把它打红 ⇒
   "跨文件豁免串用"的钉是活的）；`internal/config` 的 `BlacklistOverrides` → `Gate` 的接线**我没查**，
   票 75 也没证 ⇒ 归票 80 那条"死线"。
7. **票面 Progress 里"逐文件 Windows `-count=2`"的四组数（20/12/8/8/2）两个验收会话都没逐文件复跑**；
   我们只做到"全包 `-count=2` rc=0 + 两条新仪器在 Windows 真跑"。要逐文件账，它仍在实现者自述里。
8. **`24a66b6` 上"4 包 47 条"这个题面数字没人复现过**（实现者取的是 `f1033e1`/`131f722`，
   我取 `131f722` = 31，B 取 `536998b` = 33）。**谁要签"47"这个词，得自己去跑那个 sha。**

## 我动过仓库里的什么（A）

只写了本文件：`docs/evidence/s1/75-independent-verification.md`（append-only；
被第二个验收会话 B 也在同一文件里追加了 B0–B7，见"第 6 节"）。
**未 commit、未 push、未 checkout、未 stash、未建 worktree**；在飞代理（票 70/80/81/84 + `frontend/`）
的文件一个没碰。全部测量在三棵 `git archive` 快照里做：
`/tmp/wisp75`（`17efc2c`，变异过、已用双 oracle 证明还原，见 2e）、
`/tmp/wisp75w`（`17efc2c` 第二次独立归档，Windows 侧用它）、
`/tmp/wisp75pre`（`131f722` 改前基线）。
日志在 `/tmp/wisp75-logs/`：`g1.json` `g1b.txt` `g2.json` `g3.txt` `g4b.txt` `g5.txt`
`g6_pre.txt` `g7_win.txt` `g8_linux_risk_v.txt` `ci_testcore.log` `pathresolver.go.orig`。

---

## 第 6 节（A）：编排者在我跑到一半时已结票，且本文件被写战争过

1. **`dc91999`（12:43:17）提的是我的半成品**：那次 commit 把当时 359 行的本文件连同
   票面改名 `-done` + "六框全绿"一起落地了，我之后的读数只在磁盘上。
   我此前又发现 `1a26d21` 提的仍是**缺后半段**的版本 ⇒ **需要编排者再 commit 一次本文件**（我不 commit，硬规矩）。
2. **本文件现在有 A/B 两个验收会话的记录，且发生过一次真实的写战**：
   B 的追加把 A 的 1d/1e–第 4 组与 A 后半段**挤出过工作副本**（我从 `dc91999` + 自己的日志重建了
   缺失部分并保留 B 的全部字句；若两侧都继续"整文件重写"，下一轮还会丢）。
   ⇒ 程序性建议（A/B 共同署名）：**同一票的两个验收会话必须写两个文件**
   （`75-independent-verification-A.md` / `-B.md`），或一律禁止覆写、只允许 `>>` 追加。
3. **`git archive … | tar -x -C /tmp/wisp75` 这个目录名被两个会话共用**（任务书原文就这么写）。
   时间线对得上的一半是我的责任：**我在 12:16–12:33 变异过 `/tmp/wisp75`，同一目录当时也是 B 的取树点**
   ⇒ A 的变异落进了 B 的第一棵快照（B0 已披露，B 用 sha256 钉桩自证后续读数未受影响）。
   ⇒ 建议登记：**验收派发的快照目录必须带会话后缀**；并明确"两份独立复现"独立到什么程度
   （读数各自跑 ⇒ 独立；取树目录曾共用 ⇒ 不独立，要靠哈希钉桩补信任）。
   A55 的 `verify75` 用哪个目录、有没有落进同一污染窗口，**没人写过** ⇒ 签"两份"之前问一句。
4. **我的六条裁决与"六框全绿"不冲突**（AC#1/2/3/5 = PASS、AC#4 = 维持关闭、
   AC#6 = 按"本票的账清零"可关），但两处口径要改：
   (a) 派发简报里"lint 红因与票 75 无关"与 **A52②** 自相矛盾，**A52② 才是对的**（我和 B 各自独立复算支持）；
   (b) **账本里还没有的一条新发现**：`Environment fork assertion` 步骤在 `35558750456` 与
   `35560673885` 两次 run 里都是 **`skipped`**——它排在 `Portable package tests` 之后且**没有 `if: always()`**
   ⇒ 只要 test-core 那一步红，票 71 AC#3 专门建的"防 `-run` 空匹配"门**从不执行**。
   不是票 75 的账，但属"步骤被静默跳过"这一类假绿，建议并给票 70 的 CI 面或单开一票。
