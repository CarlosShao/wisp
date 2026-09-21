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
| G | 假绿四项逐跑点名 | 未做（事实已在 F/G 段采集，结论未写） |
| H | 1:1 裁决表 + 票 70 结论 + 改名建议 | 未做 |

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
被 `IsSyncPath` 判 `Sync=true` 且 `Why` 含 `".."`（`:460`），而 `folded` 拼法判非 sync（`:458`）并走 `Inspect` 双向（`:470-476`）。
承重码是 `syncdirs.go:363` 的 `if hasFoldedDotDot(rawPath)`，实现在 `syncdirs.go:278`（对 `unifySeparators(raw)` 按 `sepStr` 切段找 `..`），
纯词法、**先于任何 OS 咨询**。**不是恒真**：MUT-f（让 `hasFoldedDotDot` 恒 false）在 ubuntu 上把这条打红（B.5）。

### D.2 抽样 2：`TestWriteGateNotSelectedByPayloadKey`（`internal/risk/provenance_test.go:463`）＋ `m7Engine`（`:449`）

POSIX 上它断言的是：7 种 payload 键名 × 3 种 tool 名（`""`/`http.post`/`net.send`）+ 冻结的 C19 seam + 值形状 sink + `notify`/`clipboard.write`
**全部仍被 `Inspect` 命中**（`provenance_test.go:481-518`），且前提是 `plainTarget` 被判非 sync（`:466`）。
`m7Engine` 现在是 `membershipEngine(t,"registry")` 的薄封装（`provenance_test.go:450`），fixture 把 root 与 work 都搬到 profile **之外**
（`syncdirs_test.go:71-75`），于是"非 sync"只能由 root membership 给出，`syncdirs_test.go:86` 用 `Root.Source != "suspect-fallback"` 自证归因。
**不是恒真**：MUT-b（`pathKeys` 去掉 `"path"`）在 ubuntu 上把这条连同 4 条兄弟打成红（B.5）。

### D.3 抽样 3：`TestSyncEnvConfiguredRoots`（`internal/risk/syncdirs_test.go:250`）

POSIX 上它断言两件事：`envConfiguredRoots` 这个纯函数在给定 `probeEnv{OneDrive, OneDriveConsumer, OneDrivePublic:"  "}` 下
恰好返回 2 条、且都 graded `env`/`OneDrive`（`:260-267`，空白值被 trim 掉 ⇒ 3 个变量只出 2 条 root）；
以及 `IsSyncPath(<root>/Notes/new.md)` 必须 `Sync && Root.Source == "env"`（`:269`）。
后一条**是关键**：这里 root 就在 profile 之下（`sandbox` 的 home），POSIX 上兜底网本来就 armed，
所以 `Sync=true` 单独看是恒真的 —— 票面知道这点，才把判据写成**归因**（`Source=="env"`）而不是 `Sync`，
`syncdirs.go:394` 的 membership 循环在 `:400` 的兜底之前，所以"由 root membership 判"这件事在 POSIX 上是可失败的。
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

