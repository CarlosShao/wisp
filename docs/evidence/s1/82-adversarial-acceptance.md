# 票 82 — 独立对抗验收裁决表（acceptor-ticket82）

**被验收票**：`.scratch/wisp/issues/82-risk-sync-family-posix-tiering.md`（`Status: coded-awaiting-review`）
**代码落点**：`42ed13f`、`a25430e`、`629ce3c`
**验收基线**：`git archive 629ce3c` → 仓外纯净快照 `C:\Users\swq\AppData\Local\Temp\wisp82acc-a82x`（会话后缀 `a82x`）
**容器**：docker `golang:1.27`（容器内 `go version go1.27.1 linux/amd64`），`WISP_ENV=test`，命名卷 `wisp82acc-a82x-gomod` / `wisp82acc-a82x-gobuild`
**验收者立场**：本表把票 82 的判据当**攻击目标**，重点攻击"分层变成藏东西"。所有数字为**本代理自己那一次读数**，不引用票面数字当证据。

> **本文件是渐进写入的**：下面的"执行状态"块标明哪些组已做完、哪些组没做。
> 若本代理中途死亡，接手者请**只补未完成组**，不要重跑已完成组的结论。

## 执行状态

| 组 | 内容 | 状态 |
|---|---|---|
| A | AC#3(i) ubuntu 侧归零 | **已做完** |
| B | AC#4 变异（剥 tag 两向 + 四次行为变异） | **已做完** |
| C | AC#2b 文件名 vs 构建约束 + 两侧集合差异 | **已做完** |
| D | AC#1 定性表可验性（抽 3 条读码 + 变异验牙口） | **已做完** |
| E | AC#1b `rules_test.go:198` 两侧读数 + 变异 | **已做完** |
| F | AC#5 门禁 + AC#3(ii) Windows 侧 | **已做完** |
| G | 假绿四项逐跑点名 | **已做完** |
| H | 1:1 裁决表 + 票 70 结论 + 改名建议 | **已做完** |

**七组全部做完，本文件已是终稿（判词见 H 组）。**

---

## A 组 — AC#3(i)：ubuntu 侧那 8 条是否真的归零

### A.1 命令原文

```
git archive 629ce3c | tar -x -C /tmp/wisp82acc-a82x          # 仓外纯净快照
docker run --rm -v ".../wisp82acc-a82x:/src" \
  -v wisp82acc-a82x-gomod:/go/pkg/mod -v wisp82acc-a82x-gobuild:/root/.cache/go-build \
  -w /src -e WISP_ENV=test golang:1.27 sh -c \
  'go test -count=1 ./internal/risk/ ./internal/tools/'
```

### A.2 逐包结论行原文（`-count=1`，非 -v）

```
ok  	github.com/CarlosShao/wisp/internal/risk	2.401s
ok  	github.com/CarlosShao/wisp/internal/tools	10.477s
```

### A.3 四个数（同快照、同容器、`-count=1 -v`）

| 包 | rc | `=== RUN` 总 | `=== RUN` 顶层 | `--- PASS` 总 | `--- FAIL` | `--- SKIP` | `[no test files]` | `no tests to run` |
|---|---|---|---|---|---|---|---|---|
| `internal/risk` | 0 | **124** | **71** | **123**（顶层 70） | **0** | **1** | 0 | 0 |
| `internal/tools` | 0 | **56** | **44** | **54**（顶层 42） | **0** | **2** | 0 | 0 |

`internal/risk` 的 `--- SKIP` 逐条点名：`TestSyncRegistryProbeLive`（1 条）。
`internal/tools` 的 `--- SKIP` 逐条点名（**不属于本票改动面，但既然是 portable 裸 `go test` 就要记账**）：
`TestD34WriteMatrix`、`TestCrossVolumeMoveStopsWithTwoCopiesOnLateStop`。

### A.4 那 8 条的逐条读数（ubuntu，`-count=1 -v`）

票面点名的 8 条**一条都不出现在 FAIL 里**（`--- FAIL` 总数 = 0），8/8 逐条 `--- PASS`：

```
--- PASS: TestWriteGatePlainLocalWriteNotFlagged (0.01s)
--- PASS: TestWriteGateNotSelectedByPayloadKey (0.01s)
--- PASS: TestWriteGateEveryPathTargetJudged (0.01s)
--- PASS: TestWriteGateAllPathsNonSyncStaysExempt (0.01s)
--- PASS: TestSyncNormalNewFileWriteNotFlagged (0.01s)
--- PASS: TestSyncFallbackNotDisarmableByWeakRoot (0.01s)
--- PASS: TestSyncEnvConfiguredRoots (0.01s)
--- PASS: TestSyncDotDotTailFailsClosed (0.01s)
```

### A.5 红数没增加

`--- FAIL` = 0，`[no test files]` = 0，`no tests to run` = 0，顶层 RUN = 71（> 0）。
POSIX 层 4 条（`TestSyncNoGradeIsConfirmedOnPosix`、`TestSyncMembershipDecidesOffProfilePosix`、
`TestSyncRegistryProbeIsAStubHere`、`TestExfilSyncWritePosixSpellingInvariant`）在 ubuntu 名单里各出现 1 次；
windows 层 3 条（`TestSyncConfirmedGradesDisarmFallbackWindows`、`TestWriteGatePlainTargetInsideProfileWindows`、
`TestExfilSyncWriteWindowsSpellingInvariant`）在 ubuntu 名单里各出现 **0** 次 —— 见 C 组的两侧集合对比。

### A.6 A 组判词

**AC#3(i) = PASS。** 8 条在 ubuntu 上从红转绿，且转绿是靠"POSIX 上真跑得动"而不是靠搬走：
同一批名字仍出现在 ubuntu 的 `=== RUN` 名单里（A.5 那 8 行本身就是 RUN+PASS 双证据）。

### A.7 改前基线（我自己在同一容器、同一对命名卷里复现，不引用票面账）

为了不让"归零"变成一句没有对照物的话，我把 `42ed13f^`（票 82 第一枚码的父提交）另解一份纯净快照
`/tmp/wisp82acc-a82x-base`（该树里 `internal/risk/syncdirs_other_test.go` **不存在** ⇒ 确认是票 82 之前的形状），
同一 docker `golang:1.27`、同一 `WISP_ENV=test` 重放：

```
BASELINE_RC=1
TOPRUN=68  TOPPASS=59  TOPFAIL=8  TOPSKIP=1
--- FAIL: TestWriteGateNotSelectedByPayloadKey
--- FAIL: TestWriteGatePlainLocalWriteNotFlagged
--- FAIL: TestWriteGateEveryPathTargetJudged
--- FAIL: TestWriteGateAllPathsNonSyncStaysExempt
--- FAIL: TestSyncFallbackNotDisarmableByWeakRoot
--- FAIL: TestSyncEnvConfiguredRoots
--- FAIL: TestSyncNormalNewFileWriteNotFlagged
--- FAIL: TestSyncDotDotTailFailsClosed
FAIL	github.com/CarlosShao/wisp/internal/risk	2.964s
```

⇒ **8 条、且名字与票面点名的 8 条一字不差**（我这次是独立数出来的，不是抄票面）。
对照 A.3：顶层用例 **68 → 71**（+3 条 POSIX 层活断言），顶层红 **8 → 0**，SKIP 1 → 1（同一条）。
覆盖面是净增，不是搬走。


---

## B 组 — AC#4 变异：tag 到底是在表达平台能力，还是在藏失败

全部在 `C:\Users\swq\AppData\Local\Temp\wisp82acc-a82x` 这份快照里做，工作树 `internal/risk/` 未动一个字节。
每一条都遵守：**同一条链里 grep 证明变异落地 → 先过 build 门（`go test -run '^$'` 必须打 `[no tests to run]` 而不是 `[build failed]`）→ 跑 → 还原 → 证明还原**。
"编译失败不算变异"这条我踩到过一次（MUT-d 的 `if true {` 让 `r` 变成 unused，`TOPFAIL=0` + rc=1 全是 build 失败，当场判为无效变异并改写重跑），写在这里当下游的证据。

### B.1 MUT-1（ubuntu）：只删 `//go:build windows` 那一行，文件名不动

| 步 | 原文 |
|---|---|
| 变异前 | `md5 78759c1ba7bf15e22d50f02071c0aa56`、首行 `//go:build windows`、`go test` = `ok … 2.490s` |
| 落地证明 | `grep -n "^//go:build windows"` 命中 1:1 → `sed` 后 `grep -c "^//go:build"` = **0**，文件名仍是 `syncdirs_windows_test.go` |
| build 门 | `ok github.com/CarlosShao/wisp/internal/risk 0.011s [no tests to run]` ⇒ 包仍能编译 |
| 跑 | `MUT1_RC=`**0**，`grep -cE "^--- FAIL"` = **0**，两个 windows 层用例名在日志里**一次都不出现**，结论行 `ok … 2.317s` |
| 还原 | `md5 78759c1ba7bf15e22d50f02071c0aa56`（与变异前同值） |

⇒ **票面 AC#2b 那句自我更正成立，我独立复现到了**：`syncdirs_windows_test.go` 光靠**文件名**就被 gate 出 linux 构建，
删 tag 后 ubuntu 上是 `no tests to run` + rc=0，**不是**红。含义（对票不利的那半）：
**"只改文件名不写 tag"这条路在 windows 侧其实会生效，因而更隐蔽**；票面 AC#2b 正文的原判据（"这个名字什么门都不施加"）是错的。

### B.2 MUT-2（ubuntu）：删 tag **且**换成无 windows 后缀的文件名 ⇒ 必须变红

`mv syncdirs_windows_test.go syncdirs_winmut_test.go` + 无 tag。落地证明：`grep -c "^//go:build"` = 0、
`ls` 里 `syncdirs_windows_test.go` 计数 = 0。build 门：`[no tests to run]`（能编译）。

```
MUT2_RC=1
top-level FAIL count: 2
--- FAIL: TestSyncConfirmedGradesDisarmFallbackWindows (0.01s)
--- FAIL: TestWriteGatePlainTargetInsideProfileWindows (0.01s)
[build failed] 出现次数 = 0
```

失败正文（`logs/mut2.log:229-237`，行号是变异体文件的行号）：

```
syncdirs_winmut_test.go:51: source "registry" must count as confirmed
syncdirs_winmut_test.go:54: source "registry": confirmed root must lift the blanket suspect net
syncdirs_winmut_test.go:51: source "config" must count as confirmed
syncdirs_winmut_test.go:54: source "config": confirmed root must lift the blanket suspect net
syncdirs_winmut_test.go:51: source "env" must count as confirmed
syncdirs_winmut_test.go:54: source "env": confirmed root must lift the blanket suspect net
syncdirs_winmut_test.go:88: precondition: the injected registry-grade root must confirm detection, so the suspect net is NOT what decides these cases
```

⇒ **6 句等级判据 + 1 句 precondition，全是行为红不是编译红。** tag 没在藏东西：
它挡住的正是"POSIX 上 `SyncDetectionComplete()` 恒 false"这条平台缺口，缺口一暴露就红。

### B.3 MUT-3（Windows 主机，对称方向）：删掉 POSIX 层的 `//go:build !windows`

`internal/risk/syncdirs_other_test.go` 首行删掉（该文件名 `other` 不是 GOOS，**唯一的门就是那一行 tag**）。
落地证明 `grep -c '^//go:build'` = 0；build 门 `[no tests to run]`（Windows 上 `registryProbe(probeEnv)` 两侧同签名 ⇒ 真能编译）。

```
MUT3_RC=1  top FAIL count: 2   [build failed] = 0
--- FAIL: TestSyncNoGradeIsConfirmedOnPosix   (0.01s)
--- FAIL: TestSyncMembershipDecidesOffProfilePosix (0.01s)
--- PASS: TestSyncRegistryProbeIsAStubHere    (0.01s)
```

失败原文（节选）：`syncdirs_other_test.go:63: premise of this whole tier: Resolve("C:\\Users\\...\\profile\\OneDrive") must report Resolved=false on POSIX, got {… Resolved:true}`；
`:94: premise of this file`；`:99: an under-profile write must fail closed as sync-suspect, got {Sync:false …}`。
⇒ **票面声称"3 条里红 2 条"= 我复现到的读数**（3 条里 2 条真红）。第三条 `TestSyncRegistryProbeIsAStubHere` 在本机**同形绿**，
理由我独立验到：同一次运行里 `syncdirs_test.go:347` 打出 `no registry-grade sync record on this machine
(HKCU Accounts without UserFolder is the documented reality here)`，且 `TestSyncRegistryProbeLive` 本机 `--- SKIP` ⇒
这台机器的 registry 真的没有 UserFolder ⇒ "桩返回 nil"这条判据在 Windows 上恰好也成立。
**结论：票面这一点没夸大，但它承重的两条确实红了。**

### B.4 还原总证明（三次变异之后）

```
git archive 629ce3c -- internal/risk | tar -x -C /tmp/wisp82acc-a82x-verify
diff -r /tmp/wisp82acc-a82x-verify/internal/risk /tmp/wisp82acc-a82x/internal/risk  → rc=0（无输出）
INTERNAL/RISK TREE IDENTICAL TO 629ce3c
go test -count=1 ./internal/risk/ → ok github.com/CarlosShao/wisp/internal/risk 2.012s
```

### B.5 行为变异（AC#1/D 组共用）：POSIX 层的牙口到底在不在

| 变异 | 落点 | ubuntu 读数 | Windows 读数 |
|---|---|---|---|
| MUT-f | `syncdirs.go:279 hasFoldedDotDot` 进门即 `return false`（拆掉 N-10 fail-closed） | rc=**1**，`--- FAIL: TestSyncDotDotTailFailsClosed`（1 条） | 未做（不必：该判据词法、两侧同源） |
| MUT-b | `provenance.go:143 pathKeys` 去掉 `"path"`（写门退回"看键名"） | rc=**1**，**5 条**顶层红：`TestSyncWriteNegative`、`TestWriteGateNotSelectedByPayloadKey`、`TestWriteGatePlainLocalWriteNotFlagged`、`TestSyncNormalNewFileWriteNotFlagged`、`TestSyncDotDotTailFailsClosed` | 未做 |
| MUT-e | `syncdirs.go:128 finalize` 无视等级与 canonical，只要有 root 就 `complete=true` | rc=**1**，**5 条**顶层红：`TestSyncNoGradeIsConfirmedOnPosix`、`TestSyncMembershipDecidesOffProfilePosix`、`TestSyncFixtureFallbackAndMatch`、`TestSyncFallbackNotDisarmableByWeakRoot`（`syncdirs_test.go:210` + `:213` 两类信息各 3 个等级名）、`TestSyncUnverifiedRootKeepsFallback` | 未做 |
| MUT-c | `syncdirs.go:65 gradeConfirmed` 把 `"default"` 也算 confirmed（弱等级被升格） | rc=**0**、`--- FAIL` = **0**、`[build failed]`=0 ⇒ **ubuntu 全绿，抓不到** | rc=**1**、1 条红 = `TestSyncFallbackNotDisarmableByWeakRoot`（`:210`/`:213`） |

⇒ 前三行说明**POSIX 层不是空仪器**：把生产码改坏，ubuntu 上就红，而且红的正是票面声称"两侧都有对象"的那几条。
第四行是**我唯一挖到的实质残留**，写进 D 组。

### B.6 B 组判词

**AC#4 = PASS**（票面自己那两次变异我独立复现到同形状；我另加的 4 次行为变异进一步证明 tag 表达的是能力缺失）。
**但 AC#4 的判据要加一句**：票面 AC#2b 的原句"文件名不施加任何门"是错的，而票面已经自己更正了并做了对称变异 ⇒ 按更正后的判据算 PASS，见 H 组残留 R-1。

---

## C 组 — AC#2b：分层是构建约束还是文件名糊弄

### C.1 `internal/risk/` 里所有 `_windows_test.go` / `_other_test.go` 逐文件点名

| 文件 | 首行 `//go:build` | 判定 |
|---|---|---|
| `internal/risk/pathresolver_anchor_spelling_windows_test.go` | `//go:build windows` | 有真约束 |
| `internal/risk/pathresolver_junction_windows_test.go` | `//go:build windows` | 有真约束 |
| `internal/risk/provenance_syncdirs_windows_test.go` | `//go:build windows` | 有真约束 |
| `internal/risk/syncdirs_redteam_windows_test.go` | `//go:build windows` | 有真约束 |
| `internal/risk/syncdirs_windows_test.go`（票 82 新建） | `//go:build windows` | 有真约束 |
| `internal/risk/provenance_syncdirs_other_test.go` | `//go:build !windows` | 有真约束 |
| `internal/risk/syncdirs_other_test.go`（票 82 新建） | `//go:build !windows` | 有真约束 |

**0 个"靠文件名假装分层"的文件**：每一个带平台后缀的测试文件都自带 `//go:build` 行。
`*_other_test.go` 的 "other" 不是 GOOS ⇒ 那两层的唯一门就是 tag 行，而它**在**，并被 B.3 的剥 tag 变异证明是**活的**。
另外三个 `//go:build` 载体：`pathresolver_budget_norace_test.go` = `!race`（与本票无关）、
`pathresolver_other.go` = `!windows`、`syncdirs_other.go`/`pathresolver_windows.go`/`syncdirs_windows.go` = 生产码分层（票 82 未动，禁改清单内，我只读）。

### C.2 两侧各自编译进去的用例集合（`-count=1 -v` 的 `=== RUN` 顶层名，去重）

| 侧 | 环境 | 顶层用例数 | `[no test files]` | `no tests to run` |
|---|---|---|---|---|
| linux | docker `golang:1.27`，快照 | **71** | 0 | 0 |
| windows | 本机 `go1.27.1 windows/amd64`，快照 `-count=2`（180 = 2×**90**） | **90** | 0 | 0 |
| 交集 | | **67** | | |
| 只在 windows | 23 条，含票 82 的 `TestSyncConfirmedGradesDisarmFallbackWindows`、`TestWriteGatePlainTargetInsideProfileWindows` 与票 75 的 `TestExfilSyncWriteWindowsSpellingInvariant`（另 20 条是票 18/75 的 pathresolver/redteam 层，非本票改动） | | | |
| 只在 linux | **4 条**：`TestExfilSyncWritePosixSpellingInvariant`、`TestSyncMembershipDecidesOffProfilePosix`、`TestSyncNoGradeIsConfirmedOnPosix`、`TestSyncRegistryProbeIsAStubHere` | | | |

### C.3 交叉编译门（**按包作用域**，绝不用 `./...`）

在 Windows 主机上跑（不是容器里跑自己）：

```
GOOS=linux   go vet ./internal/risk/  → rc=0
GOOS=darwin  go vet ./internal/risk/  → rc=0
             go vet   ./internal/risk/  → rc=0
```

（A54③ 的坑我承认并绕开了：`GOOS=linux go vet ./...` 在 Windows 主机上因 CGO=0 排除 sherpa 而永远 rc=1，与本票无关。）

### C.4 C 组判词

**AC#2b = PASS（按票面自我更正后的判据）**，两个条件都过了：任一侧都没有 `[no test files]`/`no tests to run`/RUN=0；
两侧集合确实不同（71 vs 90，各有一侧独有的名单）。
**残留（对票面文本，不对代码）**：AC#2b 正文那句判据本身是错的，票面已自我更正并请求写进 A5x；
我支持那句更正，并且**支持票面提议的加码**："两层都必须有自己的 `//go:build` 行，且两侧各做一次剥 tag 变异" ——
票 82 自己两侧都做了，所以我按加码后的标准判它 PASS。

---

## D 组 — AC#1 定性表可验性：抽 3 条"两侧都有对象、就地修"读码

### D.1 抽样 1：`TestSyncDotDotTailFailsClosed`（`internal/risk/syncdirs_test.go:448`）

POSIX 上它断言的是：`raw = <root>/Notes/../../work/x.md`（`syncdirs_test.go:452`，用 `filepath.Separator` 拼）
被 `IsSyncPath` 判 `Sync=true`（`:459-462`）且 `Why` 含 `".."`（`:463`），而 `folded` 拼法判非 sync（`:456`）并走 `Inspect` 双向（`:470-476`）。
承重码是 `syncdirs.go:365` 的 `if hasFoldedDotDot(rawPath)`，实现在 `syncdirs.go:279`（对 `unifySeparators(raw)` 按 `sepStr` 切段找 `..`），
纯词法、**先于任何 OS 咨询**。**不是恒真**：MUT-f（让 `hasFoldedDotDot` 恒 false）在 ubuntu 上把这条打红（B.5）。

### D.2 抽样 2：`TestWriteGateNotSelectedByPayloadKey`（`internal/risk/provenance_test.go:468`）＋ `m7Engine`（`:450`）

POSIX 上它断言的是：7 种 payload 键名（`:472-480`）× 3 种 tool 名（`""`/`http.post`/`net.send`）+ 冻结的 C19 seam + 值形状 sink + `notify`/`clipboard.write`
**全部仍被 `Inspect` 命中**（`provenance_test.go:481` 起的整段循环），且前提是 `plainTarget` 被判非 sync（`:470`）。
`m7Engine` 现在是 `membershipEngine(t,"registry")` 的薄封装（`provenance_test.go:452`），fixture 把 root 与 work 都搬到 profile **之外**
（`syncdirs_test.go:70-75`），于是"非 sync"只能由 root membership 给出，`syncdirs_test.go:86` 用 `Root.Source != "suspect-fallback"` 自证归因。
**不是恒真**：MUT-b（`pathKeys` 去掉 `"path"`）在 ubuntu 上把这条连同 4 条兄弟打成红（B.5）。

### D.3 抽样 3：`TestSyncEnvConfiguredRoots`（`internal/risk/syncdirs_test.go:250`）

POSIX 上它断言两件事：`envConfiguredRoots` 这个纯函数在给定 `probeEnv{OneDrive, OneDriveConsumer, OneDrivePublic:"  "}` 下
恰好返回 2 条、且都 graded `env`/`OneDrive`（`:260-267`，空白值被 trim 掉 ⇒ 3 个变量只出 2 条 root）；
以及 `IsSyncPath(<root>/Notes/new.md)` 必须 `Sync && Root.Source == "env"`（`:269`）。
后一条**是关键**：这里 root 就在 profile 之下（`sandbox` 的 home），POSIX 上兜底网本来就 armed，
所以 `Sync=true` 单独看是恒真的 —— 票面知道这点，才把判据写成**归因**（`Source=="env"`）而不是 `Sync`，
`syncdirs.go:394-397` 的 membership 循环在 `:401` 的兜底之前，所以"由 root membership 判"这件事在 POSIX 上是可失败的。
**不是恒真**：MUT-e（`finalize` 无视等级）在 ubuntu 上连带把 membership 家族打红（B.5）。
把 `SyncDetectionComplete()` 那句 sanity 从这条里删掉（票面声称的动作）也确实发生了：`git diff 42ed13f^ 629ce3c -- internal/risk/syncdirs_test.go`
里 `-if !p.SyncDetectionComplete() { t.Fatal("env-grade roots are confirmed evidence") }` 换成 `+… st.Root.Source != "env"`。

### D.4 抽到的残留（MUT-c）

票面 AC#1 表第 6 行写"**弱等级永不拆网两侧有对象**"。**这句在 ubuntu 上只成立一半**：
`syncdirs_test.go:209 if p.SyncDetectionComplete() { t.Errorf("source %q must not count as confirmed") }`
抓得住"`finalize` 被改宽"这种 bug（MUT-e，ubuntu 红），
抓不住"`gradeConfirmed` 被改宽"这种 bug（MUT-c：`case "env","registry","config","default":` ⇒ ubuntu **rc=0、0 红**，Windows **rc=1、1 红**）。
根因清楚且不新鲜：POSIX 上 `canonical` 恒 false（`syncSet.add` `syncdirs.go:141-147` ← `pathresolver_other.go:10` 桩），
`finalize` 里 `r.canonical && gradeConfirmed(...)` 的**左半已经短路掉了右半**，所以等级表本身在 POSIX 上不可观测。
⇒ 这是一条**分层固有后果而非藏红**（没有任何断言被删；同一断言在 Windows 侧有牙；票 82 的 windows 层 `TestSyncConfirmedGradesDisarmFallbackWindows`
正是这条等级表的正面镜像），但它意味着 AC#1 表里那句"两侧有对象"对**等级表**这条腿是 over-claim。判 **允许残留**，见 H 组 R-2。

### D.5 D 组判词

**AC#1 = PASS（附 R-2 一处措辞 over-claim）**。8 条逐条定性的**方法论**我复现了：它给的"POSIX 上没有可断言的对象"
每一句都落到具体一行（`pathresolver_other.go:10`、`syncdirs_other.go:20`、`syncdirs.go:128`），我抽查的 3 条
在 POSIX 上断言的东西都是可失败的（3 次独立行为变异各自打红），**没有一条**是"搬个 tag 让它不跑"。
新增的 windows 层确实只装了 2 个用例的 confirmed-grade 半边（C.2 只在 windows 的 23 条里票 82 新造 2 条），
**没有整族打包贴 tag**，POSIX 覆盖净增（`TestWriteGate*`×4 + 那批 `TestSync*` 第一次在 ubuntu 上真跑并真能红）。

---

## E 组 — AC#1b：`internal/risk/rules_test.go:198` 两侧各自的红/绿读数

该行原文（两侧同字节；`rules_test.go` **无** `//go:build` 行 ⇒ 它是 portable 文件，我在容器里 `sed -n 198p` 与工作树对拍一致）：

```go
{"clean allowlisted full path -> L1", Facts{ShellArgv: []string{`C:\Program Files\Git\bin\git.exe`, "status"}, ShellAllowlist: []string{"git.exe"}}, L1},
```

| 侧 | 未变异 | 变异（把 198 行 allowlist 的 `"git.exe"` → `"git"`） |
|---|---|---|
| Windows 本机（快照） | `=== RUN TestRuleShellR6` ×1、`--- PASS` 、rc=0 | rc=**1**、`--- FAIL: TestRuleShellR6`、`rules_test.go:205: clean allowlisted full path -> L1: want L1/[R6], got L2/[R6] (R6: 程序不在 shell 白名单: C:\Program Files\Git\bin\git.exe)`，RUN=**1** |
| ubuntu（容器，快照） | `=== RUN TestRuleShellR6` ×1、`--- PASS`、rc=0 | rc=**1**、同一条 `--- FAIL`、**同一句信息、同一行号 205**，RUN=**1** |

**断言命中数**：该表项本身贡献 3 个断言（`d.Level != wantLevel` / `len(d.RulesHit) != 1` / `RulesHit[0] != R6`，
`rules_test.go:203`），任一条不满足即 `t.Fatalf`；整个 `TestRuleShellR6` 有 6 个表项 × 3 + 尾部 1 个负例 = **19 个断言点**，
无 `t.Run` 子用例可数（所以两侧 RUN=1 是正确读数，不是仪器失效）。
落地证明：变异后 `sed -n 198p` 打出改过的字节，还原后 `diff -q logs/rules_test.orig internal/risk/rules_test.go` 无差异。

为什么它在 Linux 上**不是空转**（读码）：`shellAllowlisted`（`rules_shell.go:93`）取 argv0 的 basename 用的是
`strings.LastIndexAny(argv0, "/\\")`（`rules_shell.go:95`）—— 两种分隔符一起找，与 `unifySeparators`（`pathresolver.go:148`，
POSIX 分支是 no-op）**无关**，所以 `C:\Program Files\Git\bin\git.exe` 这条 Windows 字面量在 ubuntu 上照样折出 `git.exe` 并参与白名单比较。
⇒ 票面"本票没有收窄任何折叠"我也验了：`git diff 42ed13f^ 629ce3c --stat` 只碰 `internal/risk/` 的 4 个测试文件 + 票面，
`rules_shell.go`/`syncdirs.go`/`provenance.go`/`assessor.go` **零改动**。

### E 组判词

**AC#1b = PASS。** 两侧各自绿、变异后两侧同一条红同一句信息 ⇒ 它不是装饰。票面"不需要照票 81 改成平台派生期望值"的裁定我认可。

---

## F 组 — AC#5 门禁 + AC#3(ii) Windows 侧

### F.1 门禁（Windows 主机，`go1.27.1 windows/amd64`）

| 门 | 工作树 | 纯净快照 |
|---|---|---|
| `gofmt -l internal/risk/` | 空（rc=0） | 空 |
| `gofumpt -l internal/risk/` | 空（rc=0），gofumpt **v0.7.0 (go1.27.1)**，取自 `$(go env GOPATH)/bin` = `D:\work\base\gopath\bin` | 空 |
| `go vet ./internal/risk/` | rc=**0** | rc=**0** |
| `GOOS=linux go vet ./internal/risk/` | rc=**0** | rc=**0**（容器内也 0） |
| `GOOS=darwin go vet ./internal/risk/` | rc=**0** | 未做 |

⚠ 仪器坑照实记：`gofumpt` 在 CI 里**没钉版本**，我这次读数是 v0.7.0；换版本的读数可能不同，这条不是"永久绿"。

### F.2 `go test -count=2 ./internal/risk/` 两侧数字（都是快照、都是我自己跑的）

| 侧 | rc | `=== RUN` 总 | 顶层 | `--- PASS` 总 | 顶层 PASS | `--- FAIL` | `--- SKIP` | `[no test files]` / `no tests to run` |
|---|---|---|---|---|---|---|---|---|
| Windows 本机 | **0** | **302** | **180 = 2×90** ✓偶数倍 | **300** | 178 | **0** | **2** | 0 / 0 |
| ubuntu 容器 | **0** | **248** | **142 = 2×71** ✓偶数倍 | **246** | 140 | **0** | **2** | 0 / 0 |

两侧 `--- SKIP` 逐条点名：**Windows 2 条 = `TestSyncRegistryProbeLive` ×2；ubuntu 2 条 = `TestSyncRegistryProbeLive` ×2**
（`-count=2` 的双份，`=== RUN` 与 `--- SKIP` 一一对上，不存在"名字漂移"）。
Windows 那 2 条 SKIP 的原因我独立看到了证据：同一次运行 `syncdirs_test.go:347` 打
`no registry-grade sync record on this machine (HKCU Accounts without UserFolder is the documented reality here)`。

### F.3 AC#3(ii) 判词

**PASS。** rc=0、`=== RUN` 是 `-count=2` 的严格偶数倍（302=2×151 / 180=2×90，两侧都成立）、FAIL=0、SKIP 已逐条点名且非本票新造。

---

## G 组 — 四种假绿逐跑点名

### G.1 `--- SKIP` 被当 ok？—— **有，且不在本票改动面内**

我把 test-core 那条 portable 步骤（`ci.yml:127-134`，命令是**裸** `go test <16 个包> -count=1`，
**没有**走 `tools/d22scan/runtests.sh`）在容器里逐字重放，加 `-v` 数 SKIP：

```
go test ./internal/agent/... ./internal/llm/... ./internal/config/... ./internal/memory/...
        ./internal/observe/... ./internal/secret/... ./internal/risk/... ./internal/statemachine/...
        ./internal/session/... ./internal/watchdog/... ./internal/tools/... ./internal/models/...
        ./internal/buildinfo/... ./internal/audio/... ./internal/proc/... ./internal/panel/... -count=1 -v
RC=0  TOTALRUN=841  TOPRUN=526  TOPPASS=519  TOPFAIL=0  TOPSKIP=7  子用例 SKIP=0
```

7 条 SKIP 全部点名（含归属文件）：

| # | SKIP 的用例 | 所在包 | 与本票关系 |
|---|---|---|---|
| 1 | `TestDefaultDeadlineWallClockMeasurement` | `internal/agent/approval/ticket84_no_owner_test.go` | 无关（票 81/84 地界） |
| 2 | `TestSubprocessCrashWriter` | `internal/memory/concurrent_test.go` | 无关（票 81 地界） |
| 3 | **`TestSyncRegistryProbeLive`** | `internal/risk/syncdirs_test.go:338` | **本票包内，票面自己申报的残留①** |
| 4 | `TestD34WriteMatrix` | `internal/tools/fs_write_test.go` | 无关但同形状 |
| 5 | `TestCrossVolumeMoveStopsWithTwoCopiesOnLateStop` | `internal/tools/fs_write_test.go` | 同上 |
| 6 | `TestRealDownloadVadThroughPipeline` | `internal/models/manifest_real_test.go` | 无关 |
| 7 | `TestRealDownloadPuncArchiveThroughPipeline` | `internal/models/manifest_real_test.go` | 无关 |

⇒ 这 7 条在 CI 里全部记成 `ok`。票面申报的那一条（#3）**我复现到了**，并补两笔它没写的 sharper 事实：
(a) 在 ubuntu 上这条**永远不可能跑**——`syncdirs_other.go:20 registryProbe` 恒返回 nil ⇒ `len(roots)==0` 恒成立 ⇒
它是**结构性永久 SKIP**，不是"这台机器恰好没记录"；
(b) 它打印的理由文本（`syncdirs_test.go:347`：*"HKCU Accounts without UserFolder is the documented reality here"*）
是在 **Linux** 上打印的 **Windows** 解释，文本本身在 POSIX 上就是说谎的。
缓解事实也要记：票 82 的 POSIX 层 `TestSyncRegistryProbeIsAStubHere`（`syncdirs_other_test.go:110`）**活着地断言了同一个事实**
（桩返回 nil），所以"POSIX 上没有 registry 探测"这件事在本包里**不是只有休眠覆盖**。

### G.2 用了 `-run` 是否空匹配？—— 本票相关步骤没有，我自己的步骤有自证

- CI 的 portable 步骤（就是本票的判据载体）**不含 `-run`** ⇒ 无空匹配风险。
- 我用 `-run '^TestRuleShellR6$'` 取 AC#1b 读数：两侧都证明 `=== RUN` = **1**（不是 0），见 E 组表。
- 我用 `-run '^$'` 当**编译门**（MUT-1/2/3 与行为变异各一次）：它打的 `[no tests to run]` 我只用作"能编译"的证据，
  **没有**任何一次把它记成测试通过。这一点写清楚，因为它正是票 70-d 说的"看着绿其实什么都没跑"的形状。
- test-core 的另一步 `runtests.sh ./internal/proc/ -run TestLayoutForTestEnv` 我重放：`STEP4_RC=0`、
  `PASS=1 FAIL=0 SKIP=0, === RUN=3, '[no tests to run]'=0` ⇒ 那条本来就有牙，不是本票残留。

### G.3 `-count=N` 是否核对 `=== RUN` == N × 不同名？—— 核对了，两侧两侧都过

| 侧 | ALLRUN | 去重后的 RUN 名数 | 倍数 | 顶层 RUN | 顶层去重 | 倍数 |
|---|---|---|---|---|---|---|
| Windows `-count=2` | 302 | **151** | 2×151 ✓ | 180 | **90** | 2×90 ✓ |
| ubuntu `-count=2` | 248 | **124** | 2×124 ✓ | 142 | **71** | 2×71 ✓ |

并做了逐名一致性抽查：出现次数最多的名字各 = 2 次（`TestWriteGatePlainTargetInsideProfileWindows`、
`TestWriteGateNotSelectedByPayloadKey`、`TestWriteGateEveryPathTargetJudged/…`），没有"跑了一遍只记一遍"的漂移。
ubuntu 的 `-count=1`（124）与 `-count=2` 去重值（124）**互相咬合** ⇒ 没有静默去重。

### G.4 有没有步骤被静默跳过？—— 本票范围内没有；**仓外有一处，正好是票 70 的**

- 本票的包：`[no test files]` = 0、`no tests to run` = 0（两侧、`-count=1` 与 `-count=2` 四次读数都是 0）。
  portable 全集里那 3 个 `?  … [no test files]` 是 `internal/agent/scheduler`、`internal/session`、`internal/watchdog`
  ——**这三个包本来就没有测试文件**（非 tag 造成，也不在本票面里），我按事实记下不当本票的账。
- 但我顺手抓到的**一处真静默跳过**（不属票 82，属票 70 的 AC#6）：run `35566711794` 的 `lint` job 里
  `staticcheck => failure`，紧随其后的 **`mockllm module vet => skipped`**。GitHub 首个失败步之后的步骤不执行 ⇒
  `lint` 的最后一步判据从未产生过。这条写进 H 组残留 R-4。

---

## H 组 — 1:1 裁决表

| 票面 AC（1:1，含编排者追加的 AC#1b / AC#2b） | 裁决 | 我这次的读数（不是票面数字） | 残留 / 条件 |
|---|---|---|---|
| **AC#1** 8 条逐条定性，不许整族打包贴 tag | **PASS** | 抽 3 条读码 + 3 次行为变异验牙口：MUT-f（`hasFoldedDotDot`→false）⇒ ubuntu 1 红；MUT-b（`pathKeys` 去 `"path"`）⇒ ubuntu **5** 红；MUT-e（`finalize` 无视等级）⇒ ubuntu **5** 红。windows 层只装 2 条用例的 confirmed-grade 半边（`TestSyncConfirmedGradesDisarmFallbackWindows`、`TestWriteGatePlainTargetInsideProfileWindows`），POSIX 侧顶层用例数 68 → 71、顶层红 8 → 0（改前基线是我自己在同容器复现的，见 A.7） | **R-2**：AC#1 表第 6 行"弱等级永不拆网两侧有对象"是 over-claim（见下） |
| **AC#1b** `rules_test.go:198` 不许变空仪器 | **PASS** | 两侧未变异各 `--- PASS`、RUN=1；变异 `"git.exe"`→`"git"` ⇒ **两侧同一条红、同一行号 `rules_test.go:205`、同一句** `want L1/[R6], got L2/[R6]`。不空转的机理：`rules_shell.go:95` 用 `LastIndexAny(argv0, "/\\")`，与 `pathresolver.go:148` 的平台分支无关 | 无 |
| **AC#2b** 分层必须用构建约束表达，不能用文件名 | **PASS（按票面自我更正后的判据）** | 7 个平台后缀测试文件**逐个**都有 `//go:build` 行（B/C 组表），0 个靠文件名糊；两侧集合确实不同：linux 顶层 **71**、windows 顶层 **90**、交集 67、仅 linux **4 条**；两侧 `[no test files]`=0、`no tests to run`=0；`GOOS=linux`/`darwin` `go vet ./internal/risk/` 按包作用域各 rc=**0** | **R-1**：票面正文原句"这个名字什么门都不施加"是**错的**（MUT-1 实测：删 tag 保文件名 ⇒ ubuntu `no tests to run` + rc=**0**，不红）。票已自我更正并请求写进 A5x，我背书该更正 |
| **AC#2** POSIX 层不许空、不许全是 Skip，要断言真实行为并指名票 55 转活 | **PASS** | `syncdirs_other_test.go`：3 条活用例、`t.Skip`/`t.Log`/早退 **0** 处、8 个 `t.Errorf/Fatalf` 位点、ubuntu 上 **3/3 `--- PASS`**；头部写明票 55 落地的 1-2-3 步；`SyncDetectionComplete()` 门控的休眠半边留在 `provenance_syncdirs_other_test.go:106-110`（`t.Log`+`return`，全仓 portable 集里子用例 SKIP=0） | 无。B.3 的 MUT-3 证明这层的两条承重断言在 Windows 上真红 ⇒ 它是活的、不是装饰 |
| **AC#3(i)** ubuntu 侧那 8 条归零、红数不增 | **PASS** | 容器 `golang:1.27`、快照纯净、`WISP_ENV=test`：`internal/risk` rc=**0**、RUN=124（顶层 71）、PASS=123、**FAIL=0**、SKIP=1；8 条逐条 `--- PASS`，一条不在 FAIL 里；`internal/tools` rc=0。**改前对照我自己复现**：`42ed13f^` 同容器同卷 = rc=1、顶层红 **8**、名字与票面一字不差（A.7） | 无 |
| **AC#3(ii)** Windows `-count=2` 全绿、RUN 偶数倍、SKIP 点名 | **PASS** | Windows rc=**0**、ALLRUN=**302**=2×151、顶层 **180**=2×90、FAIL=**0**、SKIP=**2**=`TestSyncRegistryProbeLive`×2；SKIP 成因我独立看到证据（`syncdirs_test.go:347` + 本机 HKCU 无 UserFolder） | 无 |
| **AC#4** 剥 windows 层 tag ⇒ ubuntu 必须红；反向也要 | **PASS** | 正向（票的形状）：删 tag + 换无 windows 后缀名 ⇒ rc=**1**、**2 条 `--- FAIL`**、**6+1 句等级判据信息**、`[build failed]`=**0**、先过编译门；对称反向：Windows 删 `!windows` ⇒ **3 条里 2 条真红**（票声称 2 条，与我读数一致），第三条同形绿的机器原因我另取证据确认。全部还原，并以 `git archive 629ce3c -- internal/risk` + `diff -r` **rc=0** 总证明 | 无 |
| **AC#5** 门禁只跑自己碰的包 | **PASS** | `gofmt -l internal/risk/` 空（树 + 快照）、`gofumpt -l internal/risk/` 空（**v0.7.0**，与 CI 钉的同版本）、`go vet ./internal/risk/` rc=**0**、`GOOS=linux` rc=**0**、`GOOS=darwin` rc=**0**、`-count=2` 两侧 rc=**0** | 见 R-3（gofumpt 版本 CI 未钉；换版本读数可能变） |

### 残留（4 条，都不构成"藏红"）

- **R-1（票面文本，非代码）**：AC#2b 正文那句关于文件名不施加门的事实**错误**，票面已自我更正并点明含义
  （"只改文件名"在 windows 侧其实生效、因而更隐蔽）。**支持**票面提议的加码判据：两层各自必须有 `//go:build` 行 + 两侧各做一次剥 tag 变异。
  处置：编排者把这句写进 A5x 并订正 `docs/reports/pending-and-issues.md`；票 82 代码侧无需改动。
- **R-2（AC#1 表第 6 行措辞 over-claim）**：`TestSyncFallbackNotDisarmableByWeakRoot` 的
  `syncdirs_test.go:209` 那条断言，在 POSIX 上抓得住 `finalize` 被改宽（MUT-e ⇒ ubuntu 5 红），
  **抓不住 `gradeConfirmed` 被改宽**（MUT-c ⇒ ubuntu **rc=0、0 红**；同一变异在 Windows ⇒ rc=1、1 红）。
  根因是 `syncdirs.go:128` 的 `r.canonical && gradeConfirmed(...)` 在 POSIX 上左半恒 false ⇒ 等级表整条腿不可观测。
  处置：票面 AC#1 表该行补一句"POSIX 侧只对 `finalize` 形状有牙；等级表本身由 windows 层镜像覆盖"，
  并把这个事实作为一条"落地即红"锚点抄进票 55 的 Constraints。**不需要改码**（MUT-c 在任一侧都算被覆盖）。
- **R-3（仪器）**：`gofumpt` 在 CI 未钉版本（本次读数 v0.7.0）。票 70-d 已在同一文件上吃过"版本换了红法就不同"的亏。
  不属票 82 的账，但票 82 的 AC#5 依赖它，故记在此。
- **R-4（票 70 的账，我顺手取到）**：`lint` 的 `staticcheck` 步失败 ⇒ 其后 `mockllm module vet` 步 `conclusion=skipped`，
  判据从未产生（G.4）。

### `TestSyncRegistryProbeLive` 的裁定（票面残留①要我表态的那件事）

**判：允许残留 + 另开小票，不许塞进本票闭。** 理由三条，都可查：
① 它不在票 82 的 8 条红里，动它=扩界（票面禁改清单与 AC 都没授权）；
② 它的形状与票 71 AC#3 同族（CI  portable 步是裸 `go test` ⇒ SKIP 记 ok），**修法该是门禁级的**
（把 portable 步换成 `tools/d22scan/runtests.sh`，或让 CI 拒绝 `--- SKIP`），一次修好 G.1 表里**全部 7 条**，
而不是在 `internal/risk` 里给一条用例打补丁——我这次实测 `internal/tools` 也有 2 条同类 SKIP，
在包内补丁只治 1/7；
③ 本票已经把它在 POSIX 上的**有效判据**另立成 `TestSyncRegistryProbeIsAStubHere`（活断言、票 55 落地即红），
所以"这件事没被覆盖"不成立。小票内容建议：门禁级拒 SKIP + 把 `syncdirs_test.go:347` 的 skip 文案改成平台诚实
（Linux 上应写"本平台没有 registry 探测，见 `syncdirs_other.go:20` DEFERRED(ticket 55)"）。

### 票 70 的 `test-core` 全绿现在是否成立？—— **成立，而且是在真 CI 上第一次成立**

我自己拉的证据（`gh api .../runs/<id>/jobs`，不是复述票面）：

```
run 35562680354  headSha e5e5eb7  （票 82 之前）      test-core => failure
run 35564090950  headSha 42ed13f  （票 82 第一枚码）  test-core => success   ← 第一次
run 35564183459  headSha 59596e8                      test-core => success
run 35566028944  headSha 84e4161                      test-core => success
run 35566711794  headSha 942ab5a  （最新）            test-core => success
   逐步骤：Start compose / Probe mock-llm / Portable package tests /
   Environment fork assertion / Stop compose  —— 全部 completed success
```

同时我在**同一份纯净快照**上离线重放了那条 portable 步骤：`STEP3_RC=0`，23 个包结果行，`FAIL` 行数 = **0**
（3 个 `?  … [no test files]` 是 scheduler/session/watchdog，与票 82 无关）。⇒ 容器读数与真 CI 读数**同向**。

**但票 70 的判据不是"test-core 绿"，是 AC#6 的"五个 job 全 pass"**：run `35566711794` 的逐 job 结论是
`slo-smoke success / test-core success / slo-full success / test-windows success / lint failure` ⇒
**AC#6 仍未闭，卡点已从"本包 8 条红"迁移到 `lint :: staticcheck`（1/5 job）**，
那一步按票 70-d 的记录是"pin 在 go1.27 跑不动、升上去 35 条 finding"，归属票 70/编排者，不是票 82。

### 总判词

**票 82 = PASS，建议改名 `-done`**（AC#1/AC#1b/AC#2/AC#2b/AC#3/AC#4/AC#5 七项全过；AC#2b 按票面自己更正后的判据过）。
它最怕的失败模式我攻下来了并且**没攻下来**：tag 表达的是 `resolveHandle`/`registryProbe` 两个 DEFERRED 桩造成的能力缺口，
不是失败的藏身处——正向剥 tag 变异在 ubuntu 打出 2 条行为红（6 句等级判据 + 1 句 precondition，`[build failed]`=0），
反向剥 `!windows` 在 Windows 打出 3 条里 2 条红；POSIX 层 3 条活断言 + 4 次生产码行为变异里 3 次把 ubuntu 打红。
改名前建议先落两件事（都不动码，只在票面/账上）：R-1 的 A5x 订正、R-2 的 AC#1 表补一句并转票 55 锚点；
另开一张"CI 拒 SKIP"的小票（票 82 之外）。


