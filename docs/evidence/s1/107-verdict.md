# 107 复验裁决表（ticket 107, round 3 / read-only verification leg）

- 票面：`.scratch/wisp/issues/107-ticket102-allowlist-judgment-is-red-on-posix.md`（Status `rejected-needs-fix`）
- 复验对象：HEAD = `ba0ea5121c2e23b4a48b0636205735e618a0ddde`（2026-11-05 10:2x +08）
- 本腿角色：只读复验（不写码、不改票、不动台账、不跑 git 写操作）
- 结论一句话：**修复 `8e10095` 确在 HEAD（祖先，内容在位），票 107 自己那条判据在两侧都判对；但 HEAD 上票面的字面 POSIX 命令跑不起来（票 23 的 WIP 把 `internal/tools` 的 POSIX 编译打断），且两包 `go vet`/`go test` 仍 rc=1（两条红都能点名到非本票的用例/文件）。因此 AC#3、AC#5 按字面口径判不通过。**

---

## 0. 裁决总判

| 项 | 判 | 依据 |
| --- | --- | --- |
| 修复是否已落地（`8e10095`） | **是** | `git merge-base --is-ancestor 8e10095 HEAD` = 真；`git log -S "treeResolvedAsNamed" -- internal/tools/paths.go` 唯一命中 `8e1009504`；HEAD 内容 `paths.go:99`（root 腿 `!res.Resolved && !treeResolvedAsNamed(...)`）、`:230-233`（target 腿 `resolvedForm` 再判 containment）、`:285` `treeOnDisk`、`:304` `treeResolvedAsNamed`、`:325` `resolvedForm` 全在位 |
| 票 102 那条红（POSIX） | **已修**（修前红已复现） | 修前 `4b4cfc95a` 快照 → docker/alpine `-test.count=2` → 2 次 `--- FAIL: TestPathCanonicalizerAccountsForRewrittenRoots`，报文与 CI 逐字同；HEAD 等价快照 → 2 次 PASS |
| 放行侧方向（不许比拒绝侧宽） | **守住了** | 探针 A（root 是链接）/ C（root 内链接出树）/ B（链接 root 未改写拼法）两侧均 PASS，且变异①把 target 腿退回 → POSIX ProbeC 立刻红 |
| HEAD 上票面字面检查 | **不全绿** | POSIX `internal/tools` 编译断（`system_bridge.go:129 undefined: platformSystemBackend`）；`go vet` Windows rc=1（`system_backend_windows.go:337/355 unsafe.Pointer`）；`go test` 两包各 1 条非本票红 |
| 是否还需要改票 107 的代码 | **不需要** | 见 §5：剩下要修的是**票 23 的 WIP 文件**，不是本票 |

---

## 1. 落地核查（命令 + 读数）

```
$ git rev-parse HEAD
ba0ea5121c2e23b4a48b0636205735e618a0ddde
$ git merge-base --is-ancestor 8e10095 HEAD && echo YES
YES
$ git log --oneline -S "treeResolvedAsNamed" -- internal/tools/paths.go
8e1009504 feat(92,AC#1-#7): 工作区选择 + composer 封套 + 前端输入框 v2；档位/工作区只发起请求
$ git log --oneline -- internal/tools/paths.go | head -5
f9c46f46e 356 AC#1码:pointerNotice加A档Deny支(...)
16901acba 票 181 AC#7 生产侧落地(...)
2b1a30710 fix(tools/paths): ticket 252 AC#2, one form for one path before the allowlist legs
8e1009504 feat(92,...)   <-- 107b 的修复随这枚进树（commit message 少署 107b，票面 19:5x 已登记）
be8b40338 fix(107,AC#1-#5): ...（第一轮修法，被验收判 AC#3 不通过）
```

内容复核（HEAD 工作树，只读 `grep -n`）：

| 位置 | 行 | 内容 |
| --- | --- | --- |
| root 腿（改写根确认） | `internal/tools/paths.go:99` | `if !res.Resolved && !treeResolvedAsNamed(res.Canonical) {` |
| 不可用根文案（票 102 语义保留） | `paths.go:100-102` | `"C26 expanded it onto %s but that tree is not confirmed on disk, authorizing nothing"` |
| target 腿（放行侧再要一次已解析形式） | `paths.go:230-233` | `rf, ok := resolvedForm(canonical); if !ok \|\| !rootsContain(p.roots, foldPath(rf)) { return false }` |
| 链接 root 的确认 | `paths.go:304-311` | `filepath.EvalSymlinks` + 与所写名字折叠相等 + `IsDir` |
| 已解析形式（允许尾部件尚不存在） | `paths.go:325-345` | `EvalSymlinks` 逐级上溯，`os.Lstat` 区分"不存在"与"存在却穿不过去" |
| 组件边界比较 | `paths.go:246-253` | `folded == r \|\| strings.HasPrefix(folded, r+pathSep)` |

---

## 2. 命令与读数（本机实测）

### 2.1 Windows（HEAD，工作树在开跑时 `git status --porcelain -- internal/tools internal/risk` 为空）

```
go test -count=2 -v ./internal/tools/   -> rc=1（日志 .scratch/107v-win-tools.log）
   === RUN 1070 = 535 名 ×2（`uniq -c` 最大重数 2）；--- PASS 770；--- FAIL 2；--- SKIP 0
   FAIL 名：Test229SpawnNamesAreRosterShaped ×2（票 229 名册普查，非本票）
go test -count=2 -v ./internal/risk/    -> rc=1（日志 .scratch/107v-win-risk.log）
   === RUN 562 = 281 名 ×2；--- PASS 340；--- FAIL 2；--- SKIP 2
   FAIL 名：TestC26RewriteAccountIsConsumedAtEverySecurityLeg ×2（新消费者，非本票文件）
   SKIP 名：TestSyncRegistryProbeLive ×2（HEAD 既有，非 -v 造成）
```

票 107 相关 6 条（Windows，两次执行都 PASS）：

```
--- PASS: TestPathCanonicalizerAccountsForRewrittenRoots        (票 102 原用例)
--- PASS: TestTicket107AllowlistJudgmentTwoShapes              (票 107 第一轮, 两形状)
--- PASS: TestTicket107AllowlistBoundaryIsComponentWise        (票 107 第一轮, 组件边界)
--- PASS: TestTicket107bProbeASymlinkedRewrittenRootAuthorizesNothing
--- PASS: TestTicket107bProbeCLinkInsideAllowedRootStaysOutside
--- PASS: TestTicket107bProbeBJudgedTreeIsTheOpenedTree
```

### 2.2 POSIX —— **HEAD 字面口径：跑不起来**

```
$ GOOS=linux CGO_ENABLED=0 go test -c -o /tmp/tools_linux.test ./internal/tools/
# github.com/CarlosShao/wisp/internal/tools [github.com/CarlosShao/wisp/internal/tools.test]
internal\tools\system_bridge.go:129:9: undefined: platformSystemBackend
BUILD_RC=1
$ GOOS=linux CGO_ENABLED=0 go vet ./internal/tools/
internal\tools\system_bridge.go:129:9: undefined: platformSystemBackend   (rc=1)
$ GOOS=linux CGO_ENABLED=0 go test -c ./internal/risk/     -> rc=0（risk 包可编）
```

根因（非本票）：`internal/tools/system_ops_windows.go` 是 `_windows.go`，而 `system_bridge.go:110` 的注释写着"`system_ops_windows.go` on Windows, `system_ops_other.go` elsewhere"，**`system_ops_other.go` 在 HEAD 上不存在**；引入者是 `119c3fa91`（票 23 WIP，`git log -1 -- internal/tools/system_bridge.go`）。

### 2.3 POSIX 补充口径（只为量票 107 自己的判据；**不在仓内、不改仓库**）

快照 `/d/tmp/107v-snap` = `git archive HEAD` + 一份 7 行的 POSIX 桩（`internal/tools/system_ops_other.go`，`//go:build !windows`，`return nil, ErrNoSystemBackend`），**该桩只补编译，不碰任何判定逻辑**。构建 rc=0，再进容器执行：

```
MSYS_NO_PATHCONV=1 docker run --rm -v "D:/tmp/107v-snap:/t" alpine sh -c \
  "/t/tools_linux.test -test.v -test.count=2 -test.run 'TestPathCanonicalizerAccountsForRewrittenRoots|TestTicket107'"
=> HEADGREEN_RC=0；=== RUN 12；--- PASS 12；--- FAIL 0
```

三条探针在 POSIX 上仍是"真断言"（探针前提自证：`lstat modeLink=true`、`readlink` 成功、`EvalSymlinks` 看得见目标树），并且读数仍是 **不放行**：探针 A `roots=[]`、`unusable=["%WISP107B_A_ROOT%/proj": risk: path traverses a reparse point ...]`、`InAllowlist=false`。

### 2.4 修前红复现（AC#1，POSIX 真执行）

`be8b4033^` = `4b4cfc95a`（票 107 第一轮修复的前一枚）。`git archive` 到 `/d/tmp/107v-pre`，`GOOS=linux CGO_ENABLED=0 go test -c` rc=0（该提交早于票 23，无桩问题），容器内 `-test.count=2`：

```
--- FAIL: TestPathCanonicalizerAccountsForRewrittenRoots (0.01s)
    paths_rewrite_ticket102_test.go:64: InAllowlist("/tmp/TestPathCanonicalizerAccountsForRewrittenRoots4025666308/001/proj/a.txt")
      = false although the expanded root is a real tree; the fix must not turn into option (A)
--- FAIL: TestPathCanonicalizerAccountsForRewrittenRoots (0.00s)   <-- 第二次
PREFIX_RC=1
```

归因（票面要求的 `git log -S`）：`git log -S "not confirmed on disk" -- internal/tools/paths.go` → **`a66aadf4e`（票 102 fix B）**，与票面一致；这枚红**是票 102 引入、票 107 第一轮修法（`be8b403`）换错方向、第二轮（`8e10095`）修对**。

---

## 3. AC 1:1 裁决表

| AC | 判 | 读数（命令/文件:行） | 备注 |
| --- | --- | --- | --- |
| **AC#1** 在能跑 POSIX 的环境里复现红，并归因到 commit | **通过** | 修前快照 `4b4cfc95a` + docker/alpine `-test.count=2` → `--- FAIL: TestPathCanonicalizerAccountsForRewrittenRoots` ×2，报文与 CI 逐字同（§2.4）；`git log -S "not confirmed on disk" -- internal/tools/paths.go` → `a66aadf4e` | HEAD 上字面复跑被票 23 的编译断点挡住，见 AC#3/AC#5；复现是在"能跑 POSIX"的等价快照里做的 |
| **AC#2** 两形状分别喂判定，定产生端/比较端 | **通过** | HEAD `paths.go:99` 是产生端条件，`:230-233` 是比较端前置；`TestTicket107AllowlistJudgmentTwoShapes` 两侧均 `shape R=true`、`shape E=true`（`roots=[...] rewritten=["%WISP107_ROOT%/proj" -> ...] unusable=[]`）、`cross=true`（POSIX 日志 `/d/tmp/107v-posix-green.log`、Windows `.scratch/107v-win-tools.log`） | 票面第一轮定的"产生端"成立；区别只在：**今天产生端已经有两侧都能用的确认能力**（`Resolved` 在 POSIX 由票 347 真落地：`internal/risk/pathresolver_other.go:39 resolveHandle` = realpath，commit `0c730cac2`），所以形状 E 也从 false 变 true |
| **AC#3** 两侧都真执行过、无 skip/无 `//go:build` 挡边 | **不通过（阻塞，非本票）** | Windows：6 条 ×2 全 PASS、SKIP=0（§2.1）。POSIX：**HEAD 上 `GOOS=linux go test -c ./internal/tools/` rc=1**（`system_bridge.go:129 undefined: platformSystemBackend`）⇒ 票面要求的"POSIX 侧被真正执行过"在 HEAD 上**无法用字面命令做到**；补充口径（HEAD + 7 行 POSIX 桩，桩只在 `/d/tmp`）12/12 PASS、rc=0（§2.3） | 两个待判归属：①这条编译断是 **票 23 WIP（`119c3fa91`）**，不是 107；②若按"HEAD + 桩"的实质口径，107 的行为两侧都绿，且 6 条用例都无 `t.Skip`/`//go:build`（`grep -n "t.Skip\|go:build" internal/tools/paths_ticket107*` 零命中） |
| **AC#4** 变异双向 | **通过** | 变异①（root 腿退回 `!res.Resolved`、target 腿退回 `rf, ok := canonical, true`）：`go test -c` rc=0 → POSIX `--- FAIL: TestTicket107bProbeCLinkInsideAllowedRootStaysOutside`，rc=1。变异②（`strings.HasPrefix(folded, r)`，去掉组件分隔）：`go test -c` rc=0 → POSIX `--- FAIL: TestTicket107AllowlistBoundaryIsComponentWise`（`.../proj-evil/b.txt` 被放行），Windows 同一条 FAIL（`C:\...\proj-evil\b.txt`）。落地行与还原：`sed -n "99p;230p;248p"` 打印改后整行；还原后 `diff -q` 对仓库 `paths.go` ⇒ `SNAP_RESTORED_IDENTICAL` | **口径变化须登记**：变异①里"root 腿退回"这一半在 HEAD 上**已不承重**（POSIX 上 `res.Resolved` 现在为真，短路了 `treeResolvedAsNamed`），所以票面期望的"ProbeA/102 用例转红"没有出现；真正把变异①抓住的是 **target 腿**（ProbeC）。"放行侧只认唯一已解析形式"这句话的牙在 target 腿上 |
| **AC#5** 门禁：gofmt/gofumpt/vet（含 GOOS=linux 按包）/`-count=2 -v` 两包/d22scan | **不通过（红逐条可点名到非本票）** | `gofmt -l internal/tools internal/risk` → 空；`gofumpt -l ...` → 空。`go vet ./internal/tools/ ./internal/risk/` → **rc=1**：`system_backend_windows.go:337:44` 与 `:355:32` `possible misuse of unsafe.Pointer`（票 23 新文件）。`GOOS=linux go vet` → **rc=1**：`undefined: platformSystemBackend`（同上）。`go test -count=2 -v` 两包 → **rc=1**：`Test229SpawnNamesAreRosterShaped`×2（票 229）、`TestC26RewriteAccountIsConsumedAtEverySecurityLeg`×2（新消费者 `cmd/wisp/panel_pump.go:283/291`、`internal/panel/git_readonly_exec.go:409/412`），SKIP `TestSyncRegistryProbeLive`×2（HEAD 既有）。`sh scripts/d22scan.sh` → **rc=0 clean**：`runtests.sh: PASS=35 FAIL=0 SKIP=0`；scope `bans#1-5 internal/=354 cmd/=56`、`ban#6 frontend/=116`、`ban#7 internal/tools/=47`、`ban#8 design/=39 frontend/=116 internal/=872 cmd/=203` | 107 自己的两条用例在两包全量 `-v` 里都是 PASS×2；`-count=2 -v` 的"逐条点名 SKIP/FAIL"读数见 §2.1。d22scan 是在工作树（此刻有他人在飞改动）上跑的一次，rc=0 |

**未核项（如实登记）**：POSIX 整包 `-count=2 -v ./internal/tools/ ./internal/risk/` 无法在 HEAD 字面上跑（编译断），所以"POSIX 整包 rc=0 / SKIP 逐条"这一格**没有读数**，只有上表 12/12 的定向读数；`internal/risk` 的 POSIX 包可编（`go test -c` rc=0），但本轮没有跑它的全量用例（与 107 判据无关，节省时间）。

---

## 4. 证据文件

| 路径 | 内容 |
| --- | --- |
| `.scratch/107v-win-tools.log` | Windows `go test -count=2 -v ./internal/tools/` 全文（1070 RUN / 2 FAIL / 0 SKIP） |
| `.scratch/107v-win-risk.log` | Windows `go test -count=2 -v ./internal/risk/` 全文（562 RUN / 2 FAIL / 2 SKIP） |
| `/d/tmp/107v-posix-green.log` | HEAD+桩快照 docker/alpine 12 RUN / 12 PASS / rc=0 |
| `/d/tmp/107v-mut1.log` | 变异① POSIX：ProbeC FAIL / rc=1 |
| `/d/tmp/107v-mut2.log` | 变异② POSIX：组件边界 FAIL / rc=1 |
| `/d/tmp/107v-prefix-red.log` | 修前 `4b4cfc95a` POSIX：票 102 用例 FAIL ×2 / rc=1 |
| `/d/tmp/107v-d22.log` | `sh scripts/d22scan.sh` 全文（rc=0 clean） |

快照位置（均不在仓库内、仓库未建 worktree/checkout）：`/d/tmp/107v-snap`（HEAD+POSIX 桩，已还原为 HEAD 内容，`diff -q` IDENTICAL）、`/d/tmp/107v-pre`（`4b4cfc95a`）。

---

## 5. STOP 条：还需改的代码不在票 107

1. **`internal/tools/system_ops_other.go` 缺失**（`//go:build !windows`）。
   - 证据：`system_bridge.go:110` 注释自称"`system_ops_windows.go` on Windows, `system_ops_other.go` elsewhere"，但 `ls internal/tools/system_ops_other.go` 不存在；`GOOS=linux go vet/test -c` 均 rc=1，唯一的未定义符号是 `platformSystemBackend`（`system_bridge.go:129`）。
   - 正确形状（本腿不执行，只写落点）：新增
     `//go:build !windows` / `package tools` / `func platformSystemBackend() (systemFuncs, error) { return nil, ErrNoSystemBackend }`（`ErrNoSystemBackend` 定义在 `system_bridge.go:157`）。
   - 归属：**票 23**（`119c3fa91`，WIP 提交）。不是票 107 的代码，但它是 AC#3/AC#5 POSIX 口径在 HEAD 上跑不动的直接原因。
2. **`internal/tools/system_backend_windows.go:337 / :355` `unsafe.Pointer` 用法**：`go vet` 两条。归属同上（票 23）。
3. 两条非本票红：`Test229SpawnNamesAreRosterShaped`（票 229 名册普查，`spawn_roster_census_229_test.go:192` 自报"名册外有登记 2 枚、未复审 1 枚"）、`TestC26RewriteAccountIsConsumedAtEverySecurityLeg`（新消费者 `cmd/wisp/panel_pump.go:283/291`、`internal/panel/git_readonly_exec.go:409/412` 未读改写账）。
4. **不要**为了变绿去改票 107 的判定方向（放行侧仍必须比拒绝侧窄）：本轮所有变异读数都说明，target 腿（`paths.go:230-233`）是当前唯一承重的那一发，删它就是重新打开验收探针 C。

---

## 6. 口径与诚实登记

1. **票 107 的 POSIX 历史前提已经变了**：票面写的"POSIX 上 `resolveHandle` 是 `("", false)` 桩"在 HEAD 上**不再成立**——票 347（`0c730cac2`）把 POSIX 的 realpath+lstat 落地了（`internal/risk/pathresolver_other.go:39`）。因此：
   - 今天 POSIX 上 `res.Resolved` 为真 ⇒ `treeResolvedAsNamed` 那条 fallback 在 POSIX 上**基本是冷路径**（只有 `Resolved=false` 且 `Rewritten=true` 才走）；
   - 票 107 最初那条红（`TestPathCanonicalizerAccountsForRewrittenRoots`）在 HEAD 上**早已不会红**，红修好了有两层原因（8e10095 的放行侧 + 347 的解析能力），无法把今天的绿单独记在 8e10095 名下；本轮用"修前 `4b4cfc95a` 红 / HEAD 绿"和"变异打回 → 红"两条读数把它钉住。
2. **工作树在他人在写**：本腿开跑时（10:22）`git status --porcelain -- internal/tools internal/risk` 为空；10:3x 复核时同路径出现 ` M internal/tools/bridge.go`、` M internal/tools/system.go`、` M internal/tools/system_23_test.go`、` M internal/tools/tool.go`、`?? internal/tools/ticket23_image_channel_test.go`（票 23 的在飞改动，**不是本腿做的**）。§2.1 的 Windows 读数是干净树时跑的；本腿全程未改仓库任何代码/票/台账，只新建本文件。
3. **POSIX 桩的范围**：`/d/tmp/107v-snap` 里的 `system_ops_other.go` 只补"平台后端"这一处编译面，`return nil, ErrNoSystemBackend` 是 fail-closed 方向，不参与 `InAllowlist`/`resolvedForm`/`treeResolvedAsNamed` 的任何判定；因此它的读数对票 107 的放行侧结论有效，但**不能替代**"HEAD 字面 POSIX 口径"那一格（那一格判不通过）。
4. 本腿未 push、未 commit、未跑任何 git 写操作；`git status` 里 `docs/evidence/s1/` 新增的只有本文件。
