# 票 111 · CI 步骤读数（第二次取证，attempt 2）

- 读取代理：`ci-read-111b`（只读；未改任何码、未 commit、未 push、未改票面）
- 读取时间（UTC / 本地 CST，逐节标注，每节之前重新 `date -u`）：起始 **2026-09-21 13:45:51Z / 21:45:51 CST**，收尾 **13:53Z / 21:53 CST**
- 仓库/分支：`D:\work\workspace\projects plans\Wisp`，本地 `dev`；工作区脏（115b/116/117/85 在写共享树）——**本次读数全部来自远端 run 日志，不来自本地工作树**
- 被读的 run：**`35606321404`**（workflow `ci`，event `push`，head `d3cc9ed8e2acd4fb04a9ffa5f09c0a10e9691a8f`，branch `dev`，created `2026-09-21T13:32:59Z`）
- run 终态：**`status=completed`、`conclusion=failure`**、`updated_at=2026-09-21T13:43:44Z`、`run_duration_ms=645000`
- 六枚 job：`lint` failure / `test-windows` failure / `test-core` success / `lint-frontend` success / `slo-smoke` success / `slo-full` success

---

## 0. 为什么换 run（不是"随便换了一枚"）

前一任 `ci-read-111` 被指定只读 run **`35605937530`**（head `65f85a6`）。它的五格全部"采不到"，机制性原因（前任已写在 `docs/evidence/s1/111-ci-step-readings.md`，本次复核一致）：

- 该 run `conclusion=cancelled`；`…/jobs` 返回 `{"total_count":0,"jobs":[]}`；`timing.billable={}`；`run_duration_ms=217000`
  ⇒ **一个 job 都没创建过、零 runner 分钟，全程只在队列里** ⇒ 没有任何日志可采。
- 砍它的是编排者的 push 节奏：`d3cc9ed` 于 `13:32:59Z` push 建了 run `35606321404`，前一枚在 `13:33:00Z` 被顶掉。
- `ci.yml:16-18` 的并发组：`group: ci-${{ github.workflow }}-${{ github.ref }}` + `cancel-in-progress: ${{ github.event_name == 'pull_request' }}`
  ⇒ push 事件下 `cancel-in-progress` 求值为 `false`，"在跑的"不被砍，但**"还在排队的"会被下一次 push 顶掉**。
  今晚最近八枚 dev run 实测（`gh api runs?per_page=8`，13:53Z 复看）：
  `35603607270` failure / `35604249114` **cancelled** / `35604909648` **cancelled** / `35605067027` **cancelled** /
  `35605359280` **cancelled** / `35605531736` failure / `35605937530` **cancelled** / `35606321404` failure
  ⇒ 六小时内 5 枚 cancelled，全部是排队期被顶掉的 push。
- **"采不到"不能读成"红"** ⇒ 票 111 AC#6/AC#9 前一任既不结案也不打回；本次换源读数见下。

**换源复核（本次自己跑的，非引用前任）** —— 命令：`git rev-parse <sha>:<path>`：

| 文件 | `65f85a6` blob | `d3cc9ed` blob | 一致 |
|---|---|---|---|
| `.github/workflows/ci.yml` | `6ab90edd9ec1eee5e6ec93d07cd6e7d6d251c25e` | `6ab90edd9ec1eee5e6ec93d07cd6e7d6d251c25e` | 是 |
| `scripts/portable-tests.sh` | `4ae5e5d8452de7e7d8d70f25a4dfbb7c198f3091` | `4ae5e5d8452de7e7d8d70f25a4dfbb7c198f3091` | 是 |
| `scripts/winsec-tests.sh` | `5ce46433192bcf89e47dcac3addc254e8e5f499b` | `5ce46433192bcf89e47dcac3addc254e8e5f499b` | 是 |
| `scripts/wisp-cli-tests.sh` | `5fd918ceff75bbe10661399e751ef8547d3572d6` | `5fd918ceff75bbe10661399e751ef8547d3572d6` | 是 |

`git diff --name-status 65f85a6 d3cc9ed` 只差两份 docs：
`A .scratch/wisp/issues/122-clear-the-staticcheck-backlog-after-85a.md`、`M docs/reports/pending-and-issues.md`。
⇒ **四个 blob 与任务书给出的前缀逐一吻合，CI 判据在两枚 run 之间字节相同，换源成立，可以继续读。**

（注：任务书说的"三个脚本"在树里的真实路径是 `scripts/portable-tests.sh` / `scripts/winsec-tests.sh` / `scripts/wisp-cli-tests.sh`；
仓库里**没有** `scripts/ci/` 目录——最初按 `scripts/ci/lint.sh` 取 blob 会 `fatal: path does not exist`，这是路径猜测错误，不是版本对不上。）

**取数形状判据（四条全中才算读到）** —— 三份日志都过：

| 日志 | job id | HTTP | `Content-Length` | 实际 body 字节 | 首行时间戳 | 期待命令的 `##[group]Run` 在场 |
|---|---|---|---|---|---|---|
| `test-windows` | `106355017881` | `HTTP/1.1 200 OK` | 407674 | 407674 | `2026-09-21T13:35:26.8733189Z Current runner version: '2.337.0'` | `##[group]Run bash scripts/winsec-tests.sh`、`…wisp-cli-tests.sh`、`…portable-tests.sh --scope=windows` 均在 |
| `test-core` | `106355018195` | `HTTP/1.1 200 OK` | 501150 | 501150 | `2026-09-21T13:35:25.0997668Z …` | `##[group]Run bash scripts/portable-tests.sh --scope=core` 在（第 338 行） |
| `lint` | `106355017673` | `HTTP/1.1 200 OK` | 52711 | 52711 | `2026-09-21T13:35:24.9205960Z …` | `##[group]Run go install honnef.co/go/tools/cmd/staticcheck@2025.1.1` + `staticcheck ./...` 在；mockllm 腿 `##[group]Run go vet ./...` 在 |

`timing.billable` 对照（这是"有 job"vs"零 job"的分水岭，不是零字节日志）：
`35606321404` → `billable={UBUNTU:{jobs:3},WINDOWS:{jobs:2}}`（5 个 job 记账，`total_ms` 私有仓显示 0 属正常），`run_duration_ms=645000`；
`35605937530` → `billable={}`、`jobs` `total_count=0`、`run_duration_ms=217000`。

---

## 1. AC#6 生死判据：`test-windows` 门禁之后的五步有没有被吃掉

**结论：修好了。五步全部拿到各自 conclusion，没有一步是 `skipped`。**

⚠ `.steps[].order` 确实全部返回 `null` ⇒ 下表按**数组位置**数步号；实测 `number` 与数组位置一一对应（pos==number）。

`test-windows`（job `106355017881`，`conclusion=failure`）非 Post 步骤逐字点名：

| 步号(=数组位置) | `name` 逐字 | conclusion |
|---|---|---|
| 1 | Set up job | success |
| 2 | Run actions/checkout@v4 | success |
| 3 | Run actions/setup-go@v5 | success |
| 4 | Windows ACL sealing gate (internal/winsec's own tests, ticket 110) | **failure（允许红）** |
| 5 | Cache third_party (deps.toml key) | **success（有 conclusion，非 skipped）** |
| 6 | cgo build smoke (build.ps1 fetch-deps + mingw link + doctor) | **success（有 conclusion，非 skipped）** |
| 7 | cmd/wisp CLI tests (needs the sherpa DLLs staged above, ticket 111 AC#4) | **failure（有 conclusion，非 skipped）** |
| 8 | Portable windows tests (proc/secret/config/risk/ball/perm/plugin/llmrecord) | **failure（有 conclusion，非 skipped）** |
| 9 | PathResolver junction placeholder (real cases tickets 18/20) | **success（有 conclusion，非 skipped）** |

（10–15 号位是 Post 步骤：`Post Cache third_party` 与 `Post Run actions/setup-go@v5` 为 `skipped`，`Post checkout`/`Complete job` 为 success —— 属清理步骤，不参与 AC#6 判据。）

关键机制证据：`ci.yml:247/285/295/317/353/369` 上每一步都带 `if: ${{ !cancelled() }}`
⇒ 第 4 步红掉之后，第 5/6/7/8/9 步仍然被调度并各自产出日志（上面每一行都能在 `win.log` 里找到对应的 `##[group]Run …` 块）。
这正是 AC#6 要修的东西，**本枚 run 上它是真的修好了**，不再是"一步红吃掉后续五步"。

第 4 步（门禁本体）逐字读数：`winsec-tests.sh: four numbers (all from -v output): === RUN=82  --- PASS=36  --- FAIL=6  --- SKIP=0`，
6 条顶层 FAIL：`TestAC1SealFileReportsTheInheritedGrantItCleared`、`TestAC2InheritedNoticeHasANoiseBound`（含子用例 `sealing_only_children_reports_each_of_them_once`）、
`TestSealReportsThePrincipalsItCleared`、`TestNoticeAttributionSurvivesAn83AliasOfItsOwnTree`、`TestNoticeAttributionKeepsTwoTreesApart`、`TestSealNarrowsAndNamesThePrincipalItRemovedBySID`。
⇒ 与"ACL 门禁允许红"一致，不计入 AC#6 失败。

---

## 2. AC#9：winsec 那一步在 ubuntu 腿真跑了吗

**结论：真跑了，且是真执行不是只编译。**

- 步号：**`test-core` job `106355018195` 第 7 步**，`name` 逐字 = `Portable package tests (core scope; the list and its guards live in scripts/portable-tests.sh)`，**conclusion = success**（该 job 整体 success）。
- 该步日志第 338 行起为 `##[group]Run bash scripts/portable-tests.sh --scope=core`，脚本自己打印的 scope 逐字（第 344 行）末尾含 **`./internal/winsec/`** ⇒ winsec 确实在 core 分母里。
- 真跑出来的一行（core.log 第 3821 行，逐字）：
  `2026-09-21T13:36:37.8446582Z ok  	github.com/CarlosShao/wisp/internal/winsec	0.019s`
  紧跟（第 3848 行）：`portable-tests.sh:   ok (own line)  github.com/CarlosShao/wisp/internal/winsec`
- 不是"只编译不执行"的旁证：同一包体里有真执行断言输出，例如
  `=== RUN   TestAC3POSIXSealDoesNotFoldABackslashIntoASeparator` → `--- PASS: … (0.00s)`，
  以及它打印的 `/tmp/TestAC3POSIX…` 路径与 `placement_symlink_113_other_test.go:336/:359` 的逐字 t.Log。`GOOS=linux go vet` 那类不会产 `ok … 0.NNNs`。
- 历史基线复核：此前整份 `test-core` 日志里 `internal/winsec` 出现 **0 次**；本枚 `grep -c "internal/winsec" core.log` = **4**（scope 行 / `ok … 0.019s` 行 / `runtests.sh: OK` 包列表行 / 结果表行）。
- ⚠ 覆盖度提醒（不算 AC#9 失败，但要写清）：core 腿跑的是 winsec 的 **POSIX 半边**（linux 构建标签下的用例）；Windows 半边仍只在 `test-windows` 第 4 步那 82 个 RUN 里，两腿不互相顶替。

---

## 3. AC#4：`cmd/wisp CLI tests` 步级结论

**runner 上的答案与实现方本机不同 —— 这是第二个答案，照实贴，不当 flake 抹掉。**

- 步号 7（`test-windows`），conclusion = **failure**，日志末行 `##[error]Process completed with exit code 1.`
- 逐字四数（win.log 第 1889 行）：`portable-tests.sh: four numbers (all from -v output): === RUN=67  --- PASS=29  --- FAIL=4  --- SKIP=0`
- 逐字收尾（1888/1891 行）：`runtests.sh: go test exited 1 - packages=[./cmd/wisp/ -count=1 -skip …]`、`portable-tests.sh: strict runner exited 1 for scope=[./cmd/wisp/]`
- 包行：`FAIL	github.com/CarlosShao/wisp/cmd/wisp	328.972s`
- 4 条失败逐字点名 + 失败原因（都是审批窗口/超时，不是缺 DLL）：
  1. `--- FAIL: TestTicket101ManualSwitchSurvivesRestart (3.10s)` — `run_mode101_test.go:309: the silenced L1 write should have executed, got: 审批超时（1 秒未确认），C18 一律判拒绝，已自动拒绝`；`:320: the silenced write never landed on disk`
  2. `--- FAIL: TestTicket101UntouchedConfigRestartsAtDefault (4.13s)` — `run_mode101_test.go:366: cold start 1/2/3: an unvetoed L1 window means EXECUTE, got: 审批超时（1 秒未确认）…`（三次冷启动全一致 ⇒ 非随机）
  3. `--- FAIL: TestTicket101SessionGrantDoesNotCrossRestart (4.06s)` — `run_mode101_test.go:482: control write should NOT error under auto_approve (审批超时（1 秒未确认）…)`
  4. `--- FAIL: TestComposedGateBlocksAWriteForTwoSeconds (301.06s)` — `run_test.go:378: an unvetoed L1 window means EXECUTE, got: 审批超时（300 秒未确认），C18 一律判拒绝，已自动拒绝`
- 分母对照：本机期望 `PASS=33 / SKIP=0 / rc=0`；runner `PASS=29 + FAIL=4 = 33` ⇒ **顶层用例总数对得上，差的是这 4 条在 runner 上判红，rc=1**。
- 排除"mingw/cache 拖累"这一解释：第 5 步 `Cache third_party` = **success**、第 6 步 `cgo build smoke (fetch-deps + mingw link + doctor)` = **success**，第 9 步（同 runner、依赖 risk 包）= success；失败集中在**审批超时**这一类断言，其中 4 条 `1 秒未确认`、1 条 `300 秒未确认`（跑到测试自身上限 301.06s），形状像 L1 窗口/审批确认在 hosted runner 上没被满足，且**在同一枚 run 里可复现（3/3 冷启动）**。
- ⇒ AC#4 **不结案**：步级结论是 failure，实现方本机数与 runner 数不一致，且不一致的原因是功能性断言而非环境缺件。

---

## 4. AC#2/#3 的分母与"未记账 skip"

- **core 腿（第 7 步）结果表：逐行数 = 25 行**（`grep -c "portable-tests.sh:   ok (own line)\|FAIL (own line)" core.log` = 25），全 25 行逐字点名：
  `cmd/llmrecord`、`internal/agent`、`internal/agent/approval`、`internal/audio`、`internal/ball`、`internal/buildinfo`、`internal/config`、`internal/llm`、`internal/llm/adaptertest`、`internal/llm/anthropic`、`internal/llm/golden`、`internal/llm/openaichat`、`internal/llm/openairesponses`、`internal/memory`、`internal/models`、`internal/observe`、`internal/panel`、`internal/perm`、`internal/plugin`、`internal/proc`、`internal/risk`、`internal/secret`、`internal/statemachine`、`internal/tools`、`internal/winsec`
  ⇒ **25 行，符合期待**（末行 `internal/winsec` 即 AC#9 的新增分母）。
  四数行逐字：`=== RUN=1036  --- PASS=661  --- FAIL=0  --- SKIP=0`；job `test-core` 整体 **success**。
- **windows 腿（第 8 步）结果表：逐行数 = 8 行**（win.log 第 2964–2971 行）：
  `ok  cmd/llmrecord`、`ok  internal/ball`、`ok  internal/config`、`ok  internal/perm`、`ok  internal/plugin`、`ok  internal/proc`、
  **`FAIL (own line) github.com/CarlosShao/wisp/internal/risk`**、`ok  internal/secret`
  ⇒ **8 行，符合期待**；该步四数行：`=== RUN=405  --- PASS=262  --- FAIL=12  --- SKIP=1`。
- **未记账 skip 的分布（这格的答案是"两腿不对称"，不是"归零"）**：
  - ubuntu（core）腿：**真的归零** —— 四数行 `--- SKIP=0`，且整份 core.log 里**没有** `unaccounted SKIP lines` 这一段（`grep -c "unaccounted SKIP" core.log` = 0）。
  - windows 腿：**仍有 1 条未记账 skip**，脚本按新形状把 file:line + 原因一起打了出来（win.log 第 2962–2963 行）逐字：
    `portable-tests.sh: unaccounted SKIP lines, each with the file:line and reason it printed:`
    `   syncdirs_redteam_windows_test.go:220: no live sync root on this machine (detected roots: [])`
    `-- SKIP: TestSyncRedTeamRealOneDrive (0.00s)`
    ⇒ 这条是 hosted runner 上没有真 OneDrive 根的**环境性 skip**；机制上它现在**可见、可点名**（这正是票要的形状），但"windows 腿归零"这件事本身没成立。登记给编排者：若要它归零，得给该用例一个永不 skip 的替身层或明写豁免判据。

---

## 5. lint 两格（交给票 85/85a 用）

⚠ **版本口径必须明写**：本枚 run 读的 `.github/workflows/ci.yml` blob = `6ab90edd…`，与 `65f85a6` 逐字节相同，
**是票 85a 之前的门禁定义**（85a 仍在写、未进 `d3cc9ed` 树：`git ls-tree d3cc9ed .scratch/wisp/issues` 只有 `85-lint-tools-never-produced-a-verdict.md` 与新登记的 `122-clear-the-staticcheck-backlog-after-85a.md`，没有 85a 本体）。
⇒ 下面两格是**旧门禁**的读数，不要当成 85a 之后的状态。

**① `staticcheck`（`lint` job 第 9 步）= failure，崩溃串照旧。**
- 步名 `staticcheck`，命令逐字 `go install honnef.co/go/tools/cmd/staticcheck@2025.1.1` + `staticcheck ./...`；末行 `##[error]Process completed with exit code 1.`
- 5 行崩溃输出逐字（lint.log 第 528–532 行，全部同一形状）：
  - `-: internal error in importing "internal/byteorder" (cannot decode "internal/byteorder", export data version 4 is greater than maximum supported version 2); please report an issue (compile)`
  - `-: internal error in importing "internal/cpu" (cannot decode "internal/cpu", export data version 4 is greater than maximum supported version 2); please report an issue (compile)`
  - `-: internal error in importing "internal/goarch" (cannot decode "internal/goarch", export data version 4 is greater than maximum supported version 2); please report an issue (compile)`
  - `-: internal error in importing "math/bits" (cannot decode "math/bits", export data version 4 is greater than maximum supported version 2); please report an issue (compile)`
  - `-: internal error in importing "unicode/utf8" (cannot decode "unicode/utf8", export data version 4 is greater than maximum supported version 2); please report an issue (compile)`
  ⇒ 基线 `export data version 4 is greater than maximum supported version 2` **命中，一字不差**；崩溃在**导入标准库预编译导出数据**阶段，`staticcheck ./...` 因此**一条 finding 都没产出**（工具没裁决，不是"裁决为干净"）。根因形状仍是 staticcheck `2025.1.1` 的 export-data 读取器跟不上 runner 上那版 Go 的 `4`。
- 前置步骤（供 85 参考，均 success）：第 4 步 D22 positive control、第 5 步 D22 seven-ban + emoji scan、第 6 步 gofmt (gofumpt)、第 7 步 go vet (module)、第 8 步 go vet (tools/d22scan module)。

**② `mockllm module vet`（`lint` job 第 10 步）= success ⇒ R-4（被 staticcheck 连带 skip）在这枚 run 上已经不成立。**
- 步名逐字 `mockllm module vet`，`ci.yml:111` 带 `if: ${{ !cancelled() }}`、`ci.yml:113` `working-directory: tools/mockllm`。
- API 结论：`{"number":10,"name":"mockllm module vet","conclusion":"success"}`（**不是 skipped**）。
- 日志在场证据：第 543 行 `##[group]Run go vet ./...` 紧随 staticcheck 的 `##[error]Process completed with exit code 1.`（第 533 行）之后启动，命令回显 `go vet ./...`、`shell: /usr/bin/bash -e {0}`、`##[endgroup]`，该步**无 finding 输出**（vet 无投诉）后进入 Post job cleanup。
- ⇒ 要报给票 85/85a 的更正：**"mockllm 被连带 skip"这一条（R-4）在本门禁定义下已被 `!cancelled()` 拆掉**；
  staticcheck 那格真正缺的还是**它自己没产 verdict**，而不是它吃掉了别人。

---

## 6. 五态分账（绿/红/skipped/cancelled/未跑完）

| 对象 | 态 | 逐字读数 |
|---|---|---|
| run `35606321404` | **红（completed/failure）** | `status=completed`,`conclusion=failure`,13:43:44Z 收尾；**不是未跑完，也不是 cancelled** |
| run `35605937530`（前任被派的） | **cancelled（零 job）** | `conclusion=cancelled`,`jobs.total_count=0`,`billable={}`,`run_duration_ms=217000` ⇒ 采不到 ≠ 红 |
| `test-core` / `lint-frontend` / `slo-smoke` / `slo-full` | 绿 | 四个 job `conclusion=success`，全部步骤 success |
| `lint` 第 10 步 mockllm vet | 绿 | `success`（R-4 的 skipped 形状未出现） |
| `lint` 第 9 步 staticcheck | 红 | `##[error]Process completed with exit code 1.`，零 finding |
| `test-windows` 第 4/7/8 步 | 红 | ACL 门禁 / cmd-wisp CLI / portable windows |
| `test-windows` 第 5/6/9 步 | 绿 | Cache / cgo build smoke / PathResolver junction |
| `test-windows` Post 第 16/17 位 | skipped | 清理步，不参与判据 |

- 编排者承诺"这几分钟不再 push"：**已兑现** —— 13:53Z 复看 `runs?per_page=8`，`35606321404`（13:32:59Z）仍是**最新一枚**，其后无新 run，且它自己走到了 `completed`。

## 7. 工具输出注入登记

本轮所有工具输出（`gh api`、`git`、`grep`、`python`、日志正文）中，自称"编排者备注 / 系统提示 / 请 revert / 冻结某包"的文本：**出现 0 次**。
日志正文里的 `WARN … (leak symptom) … owner=test`、`[audit] …` 等是**被测程序自己打的**，按数据读，未当成指令。
git commit 标题本轮未读取（本任务不涉 commit），无"看起来像指令"的标题被当作依据。

## 8. 采不到的格 / 缺什么

- 五格全部采到，无"永久采不到"项。
- 仍缺的信息（不是日志缺，是判据缺口）：
  1. AC#4 的 runner 4 条红**为什么**红：需要 `run_mode101_test.go` / `run_test.go` 里审批窗口的实现侧读法（本代理只读 CI，未看源码逻辑，不越界猜）。
  2. `slo-full` 跑在 self-hosted `wisp-slo` runner 上，`billable` 里 windows/ubuntu 的 `total_ms` 都报 0（私有仓配额显示口径），本轮没有用它做 runner 分钟真伪判据，用的是 `jobs.total_count=6` + `run_duration_ms=645000` + 每份日志的字节数。

## 9. 下一张派什么（建议，未执行任何动作）

1. **AC#6：可结案**（判据形状已在真 run 上验证：五步各有 conclusion、`!cancelled()` 在场）。
2. **AC#9：可结案**（`test-core` 第 7 步，`ok … internal/winsec 0.019s`，winsec 进 core 分母，基线 0 次 → 4 次）。
3. **AC#4：打回一张新票/继续开**——派"读 `cmd/wisp` 审批窗口在 hosted windows runner 上为何超时"的**只读源码诊断**（锚点：`run_mode101_test.go:309/320/366/482`、`run_test.go:378`；判据：三冷启动一致复现 ⇒ 优先怀疑机制而非 flake）。
4. **票 85/85a：两条更正要带进去**——① staticcheck 崩串基线**照旧命中**，本 run 是**85a 之前**的定义；② **R-4 已死**（mockllm 腿 `!cancelled()` 且 `success`），85a 不必再为它设计。
5. **新的小格**：windows 腿仍有 1 条未记账 skip（`TestSyncRedTeamRealOneDrive`，`syncdirs_redteam_windows_test.go:220`）——若要"windows 腿 SKIP=0"成立，需一个永不 skip 的替身层；否则把它写成明面豁免判据。
6. **派单节奏**：push 建 run 后到该 run `completed` 之前别再 push（本轮 6 小时内 8 枚里 5 枚 cancelled，全属排队期被顶掉 ⇒ 读数一旦落在 cancelled run 上就是永久采不到）。
