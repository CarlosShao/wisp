# CI 读数审计 `ci-read-a1fd5bf-r1`：run 36003984868 逐步账

角色：只读 CI 日志审计程（READ-ONLY）。仓库 `D:\work\workspace\projects plans\Wisp`，默认分支 `dev`。
本程零构建、零 `go test` / `go vet` / `gofmt` / `wisp` / `docker`、零 push；工具只用 `gh`（read 子命令）、
`Read` / `Grep` / `Glob`、`git log/show/diff/blame/ls-tree/merge-base`。唯一可写件＝本文件。

取数时刻：`2026-09-24`（读数窗口见第 5 节）。所有 sha 已过 `git cat-file -t`，见 §0。

---

## 0. 锚点与被审对象

| 项 | 值 | `git cat-file -t` |
|---|---|---|
| 本 run | `36003984868`（工作流 `ci`，run number 266，event `push`，branch `dev`） | n/a |
| 被检树（run `head_sha`） | `a1fd5bf1795196df3913686b804664c787833700` | commit |
| 派单给的本仓 tip | `5889559`（多一枚 docs） | commit |
| 本程取数时的实际 tip | `c1e122b`（`5889559` 之后另一枚只读程连交 3 枚 `136 AC#15 census r2` 分节 commit：`77eac9f`、`1061238`、`c1e122b`，逐枚 `--name-only` 核过只碰它自己那枚文件） | commit |
| 35 枚窗口之前的 tip（基数） | `30e19ef`（上一枚 push run `35996231784` 的 head_sha） | commit |
| 35 枚窗口的第一枚 | `5a946d3`（＝ `a1fd5bf~34`，含自身共 35 枚） | commit |
| 票 85a 内联基线锚 | `d3cc9ed` | commit |
| AGENTS.md 生成锚 | `4e66817` | commit |

**35 枚窗口的口径已核死**（不是猜的）：`git diff --name-only 30e19ef..a1fd5bf -- '*.go'` 只有三枚文件，
全在 `internal/observe`：`sampler.go`、`sampler_settle_coverage_136_test.go`、`sampler_settle_gate_136_test.go`。
非 `docs/`、非 `.scratch/` 的改动集合与它逐字相同（零枚额外）。窗口内 35 枚 commit 的 subject 里只有
`52191ce`（`feat(136 AC#14 Gate)`）动过可执行码，另两枚（`f9bc512`、`c28d4e8`）自称只改注释，与三枚文件的
diff 事实不冲突。于是后面任何"这枚红是不是今天 35 枚造的"的判定都以这个集合为尺子。

**共树观察（如实记，不美化）**

- 派单里那句"另一枚程在写 `docs/evidence/s1/136-ac15-denominator-census-r1.md`"**文件名对不上**：本程取数期间
  索引里真实出现过的、必须永不暂存的那枚是 **`docs/evidence/s1/136-ac15-target-census-r2.md`**（另一枚只读程）。
  本程全程未 `git add` 它、未把它写进任何 commit。
- 本程第一次 commit 尝试（只带我那枚 pathspec）**失败并 rc=1**，报错是
  `error: pathspec 'docs/evidence/s1/ci-read-a1fd5bf-r1.md' did not match any file(s) known to git`
  ——原因是我那枚文件当时还是 untracked，`git commit -- <path>` 对 untracked 路径不收。**那次失败没有产生任何 commit，
  也没有吞掉别人索引里的东西**（失败后 `git log -1` 仍是那一程自己的 commit，`git show --name-only` 只列它自己那枚文件）。
  教训值得单记：**在共享工作树里发不带 `git add` 的 pathspec commit，第一次必然踩这枚坑**；
  正确形状是先 `git add -- <自己的那一枚路径>` 再 `git commit -- <同一枚路径>`，两步都不许出现 `-A`／`.`。
- 工作树里 16 枚 `design/**` 的删除与未跟踪的 `design/old/`、`design/doubao/` 是 owner 那侧的活，
  本程一枚未碰、一枚未暂存（见 [[frontend-delegated-to-external-agent]] 09-24 扩展条）。
- `gh run view --log` 与 `gh api .../jobs/<id>/logs` 同一条日志行微秒尾数不同这件事（票 134 面记过）在本程复现：
  本程所有逐字引用一律取自 `gh run view --log`，取法在此写明。

---

## 1. 逐步账（step ledger）

**run 终态：`failure`**（`status=completed`，`createdAt 2026-09-24T13:11:38Z`，`completedAt 13:20:31Z`）。
6 枚 job、**69 枚 step**：`success` 64 / `failure` 3 / `skipped` 3。
"产日志"列的判据：该 step 有正文行（`gh run view --log` 能取到非 `##[group]` 的输出）。

### 1.1 `lint`（job `107647320781`，`ubuntu-latest`，job conclusion **failure**）

| # | step 名 | status | conclusion | 产日志 |
|---|---|---|---|---|
| 1 | Set up job | completed | success | 是 |
| 2 | Run actions/checkout@v4 | completed | success | 是 |
| 3 | Run actions/setup-go@v5 | completed | success | 是 |
| 4 | D22 scanner positive control (tools/d22scan tests, seeded red) | completed | success | 是 |
| 5 | D22 seven-ban + emoji scan (tools/d22scan) | completed | success | 是 |
| 6 | gofmt (gofumpt) | completed | success | 是 |
| 7 | go vet (module) | completed | success | 是 |
| 8 | go vet (tools/d22scan module) | completed | success | 是 |
| 9 | **staticcheck** | completed | **failure** | 是（本 job 唯一红，见第 2 节） |
| 10 | mockllm module vet | completed | success | 是 |
| 19 | **Post Run actions/setup-go@v5** | completed | **skipped** | **否，读数永久采不到** |
| 20 | Post Run actions/checkout@v4 | completed | success | 是 |
| 21 | Complete job | completed | success | 是 |

正面记账：`mockllm module vet`（step 10）在 step 9 红的前提下**仍交出自己的 conclusion=success**，于是
`ci.yml` 里票 85a/111 记的 **R-4**（非逐字复述：staticcheck 崩掉时依赖步被记成 `skipped`）在 run 级别已闭合，
两枚步在同一枚 run 里各自带颜色，正是那段注释写的结案条件。

### 1.2 `test-windows`（job `107647320284`，`windows-latest`，job conclusion **failure**）

| # | step 名 | status | conclusion | 产日志 |
|---|---|---|---|---|
| 1 | Set up job | completed | success | 是 |
| 2 | Run actions/checkout@v4 | completed | success | 是 |
| 3 | Run actions/setup-go@v5 | completed | success | 是 |
| 4 | Windows ACL sealing gate (internal/winsec's own tests, ticket 110) | completed | success | 是（`PASS=58 FAIL=0 SKIP=0`，`=== RUN=101`） |
| 5 | Cache third_party (deps.toml key) | completed | success | 是（cache 命中） |
| 6 | cgo build smoke (build.ps1 fetch-deps + mingw link + doctor) | completed | success | 是 |
| 7 | **cmd/wisp CLI tests** (ticket 111 AC#4) | completed | **failure** | 是（`PASS=50 FAIL=4 SKIP=0`，`=== RUN=101`） |
| 8 | **Portable windows tests** (proc/secret/config/risk/ball/perm/plugin/llmrecord) | completed | **failure** | 是（`PASS=266 FAIL=12 SKIP=1`，`=== RUN=409`） |
| 9 | PathResolver junction placeholder (real cases tickets 18/20) | completed | success | 是（`PASS=1 FAIL=0 SKIP=0`） |
| 16 | **Post Cache third_party (deps.toml key)** | completed | **skipped** | **否，读数永久采不到** |
| 17 | **Post Run actions/setup-go@v5** | completed | **skipped** | **否，读数永久采不到** |
| 18 | Post Run actions/checkout@v4 | completed | success | 是 |
| 19 | Complete job | completed | success | 是 |

### 1.3 `test-core`（job `107647320486`，`ubuntu-latest`，conclusion **success**）

11 枚 step 全 `completed/success`，**零枚 skipped**，全产日志：
1 Set up job / 2 checkout / 3 setup-go / 4 Start compose test services (mock-llm on 18080) /
5 Probe mock-llm / 6 Environment fork assertion / 7 Portable package tests (`--scope=core`) /
8 Stop compose services / 15 Post setup-go / 16 Post checkout / 17 Complete job。

### 1.4 `slo-smoke`（job `107647320562`，`windows-latest`，conclusion **success**）

11 枚 step 全 `completed/success`，**零枚 skipped**：1 Set up / 2 checkout / 3 setup-go /
4 Cache third_party / 5 Build wisp.exe / 6 **SLO smoke gate** / 7 Upload SLO report /
12 Post Cache / 13 Post setup-go / 14 Post checkout / 15 Complete job。

### 1.5 `slo-full`（job `107647320611`，`[self-hosted, wisp-slo]`，runner `wisp-selfhosted-01`，conclusion **success**）

9 枚 step 全 `completed/success`，**零枚 skipped**：1 Set up / 2 checkout / 3 setup-go /
4 Build wisp.exe (deps cached on the runner) / 5 **SLO full gate (six states + settle + leak)** /
6 Upload SLO report / 11 Post setup-go / 12 Post checkout / 13 Complete job。

**注意：本 job 的 "success" 是假绿形状**：step 5 零枚取样、`slo-report.json` 未写、`Upload SLO report`
以 `##[warning]No files were found` 收尾，整枚 run 名下**没有 `slo-full-report` 这枚 artifact**。
逐字取证与判读全部在**第 4 节**。这里先在账上钉住：9/9 step 交出了颜色、零枚 skipped，
所以这道"没做事却绿"的洞**不是 skipped 吃掉的，是 step 自己以 success 交出来的**——
读日志的人不看正文就会签它。

### 1.6 `lint-frontend`（job `107647320580`，`ubuntu-latest`，conclusion **success**）

12 枚 step 全 `completed/success`，**零枚 skipped**：1 Set up / 2 checkout / 3 setup-node / 4 npm ci /
5 typecheck (tsc -b) / 6 lint (oxlint) / 7 token drift guard / 8 build (vite build) /
9 L2 card renders the real risk fields / 17 Post setup-node / 18 Post checkout / 19 Complete job。

### 1.7 skipped 逐枚点名与射程判定（本仓规矩：skipped＝那一步的读数**永久**采不到，不是"这次没采"）

三枚 skipped 全部是 Actions 的 post 步，**没有一枚是门禁步**，本 run 的门禁读数一个都没丢：

| 哪枚 job 的哪一步 | 后果（可采 / 永不可采） | 判定 |
|---|---|---|
| `lint` step 19 `Post Run actions/setup-go@v5` | "该 job 里 Go 缓存被 post 阶段写回过"这个读数永不可采 | 形状解释：`setup-go` 未启 `cache`（`slo-full` 才显式 `cache: false`），无东西可存，所以 GitHub 不派 post 日志。**不是门禁洞** |
| `test-windows` step 16 `Post Cache third_party (deps.toml key)` | "本次 push 是否把 `third_party` 存回缓存"读数永不可采 | step 5 是 cache **命中**（restore-keys 命中，所以无新 key 可存），`actions/cache` 命中路径不跑 save。**不是门禁洞**；但记账要如实：`a1fd5bf` 这一树的 cache-save 证据在这一趟是**采不到**而非"采到且为空" |
| `test-windows` step 17 `Post Run actions/setup-go@v5` | 同上，post 阶段读数永不可采 | 与 `lint` 那枚同形，**不是门禁洞** |

**这一节的总裁（step 层）**：本 run 的红**只落在 3 枚 step**：`lint/9 staticcheck`、
`test-windows/7 cmd/wisp CLI tests`、`test-windows/8 Portable windows tests`。
三枚 skipped 均非门禁，**门禁读数零丢失**；`slo-full` 的问题相反——它 step 层绿且产了日志，
但正文自陈零取样（第 4 节）。"a run being red is not a finding; a step being red is" 在这一趟
按字面成立；再加一条本趟的教训：**step 绿也不是读数，正文才是**。

---

## 2. `lint` 红因归因：唯一失败的是 `staticcheck`，44 枚 finding，逐枚归到 file:line

### 2.1 是哪一条命令红的（读 `.github/workflows/ci.yml` 的 `lint` job 定义＋该步正文，不是猜）

`lint` job（`ubuntu-latest`）按 `ci.yml` 定义共 10 条实质步（`ci.yml:65-217`）：
positive control（`bash tools/d22scan/runtests.sh -C tools/d22scan ./...`）、真扫描（`sh scripts/d22scan.sh`）、
`gofmt (gofumpt)`（`go install mvdan.cc/gofumpt@latest` 后 `gofumpt -l . tools/d22scan tools/mockllm`，
非空即 exit 1）、`go vet (module)`、`go vet (tools/d22scan module)`、`staticcheck`、`mockllm module vet`。

**红的只有 `staticcheck`（step 9）**。逐条排除都有正文：

| 步 | 结论 | 排除依据（该步正文） |
|---|---|---|
| 4 positive control | success | `tools/d22scan` 自带测试全跑，含 `TestBuiltBinaryGoesRedEndToEnd` 那一族 |
| 5 d22scan 禁令＋emoji 扫描 | success | `scripts/d22scan.sh` 走的是被扫面自证形状，正文里 `scan_test.go:250: real repo production Go files in scope: 225`，**分母非零**，所以这一步是真跑了且真干净，不是"匹配到零个文件"的假绿 |
| 6 `gofmt (gofumpt)` | success | `gofumpt -l` 输出为空（该步的判据就是"输出非空即 exit 1"），所以**格式债为零枚** |
| 7 / 8 `go vet`（两枚 module） | success | `go vet ./...` rc=0 |
| 9 **`staticcheck`** | **failure** | 见 §2.2 的自报分母 |
| 10 `mockllm module vet` | success | 在 step 9 红之下仍答了（`if: ${{ !cancelled() }}` 生效），R-4 结案 |

`staticcheck` 步自己的分母行（逐字，取法 `gh run view --log`）：

```
--- module .: exit=1 packages=32 findings=41 toolchain-crash-lines=0
--- module tools/d22scan: exit=1 packages=1 findings=2 toolchain-crash-lines=0
--- module tools/mockllm: exit=1 packages=1 findings=1 toolchain-crash-lines=0
staticcheck self-report: version=staticcheck 2026.2.1 (0.8.1) modules=3 packages=34 findings=44 toolchain-crash-lines=0 (step exit is 1)
```

所以这枚门禁**真的执行完了**：`toolchain-crash-lines=0`，票 85a 那枚"export data version 4 is greater than
maximum supported version 2、121 枚 run 零条 finding"的病形状在这一趟**没有复发**。
44 枚＝41＋2＋1，跨 3 枚 module、34 枚 package。这一步按 `ci.yml` 的票 85a 注释就是 **EXPECTED RED**，
且明写"lint turning green is NOT a criterion of this ticket"。

### 2.2 44 枚 finding 全量点名＋逐枚归因

归因方法（三把独立的尺子，任何一把都不够）：
(1) 该 finding 的 `file:line` 用 `git blame -L n,n a1fd5bf` 取最后改它的 commit；
(2) 那枚 commit 是否 `30e19ef`（35 枚窗口之前的 tip）的祖先，用 `git merge-base --is-ancestor` 逐枚判；
(3) 窗口内被改过的只有 3 枚 Go 文件这一事实，加"unused 类 finding 会因调用方被删而新造"这一枚**特定风险**，
对落在 `internal/observe` 的 3 枚 finding 另做**两棵树对照**（见 §2.3）。

| # | file:line | code | blame | 日期 | 分类 |
|---|---|---|---|---|---|
| 1 | `cmd/wisp/dataroot_128_test.go:60` | ST1012 | `4e5d240` | 09-23 | 存量 |
| 2 | `cmd/wisp/leg_sink_gate_131_test.go:225` | U1000 | `6f702ae` | 09-23 | 存量 |
| 3 | `cmd/wisp/run.go:110` | U1000 | `cd011b8` | 09-20 | 存量 |
| 4 | `cmd/wisp/run.go:234` | U1000 | `cd011b8` | 09-20 | 存量 |
| 5 | `cmd/wisp/run.go:306` | S1011 | `cd011b8` | 09-20 | 存量 |
| 6 | `cmd/wisp/run_test.go:70` | U1000 | `cd011b8` | 09-20 | 存量 |
| 7 | `frontend/embed.go:4` | SA9009 | `9318bb8` | 09-21 | 存量 |
| 8 | `internal/agent/approval/approval.go:316` | U1000 | `35200c7` | 09-20 | 存量 |
| 9 | `internal/agent/approval/fakes_test.go:104` | U1000 | `35200c7` | 09-20 | 存量 |
| 10 | `internal/agent/approval/fakes_test.go:394` | U1000 | `35200c7` | 09-20 | 存量 |
| 11 | `internal/agent/approval/queue.go:472` | U1000 | `35200c7` | 09-20 | 存量 |
| 12 | `internal/agent/harness_test.go:69` | U1000 | `ae23da6` | 09-20 | 存量 |
| 13 | `internal/agent/testtools_test.go:31` | U1000 | `1e5d96f` | 09-20 | 存量 |
| 14 | `internal/agent/testtools_test.go:49` | U1000 | `1e5d96f` | 09-20 | 存量 |
| 15 | `internal/audio/audio.go:156` | U1000 | `2aeae00` | 09-19 | 存量 |
| 16 | `internal/audio/device.go:45` | U1000 | `77538f1` | 09-19 | 存量 |
| 17 | `internal/audio/device.go:51` | U1000 | `77538f1` | 09-19 | 存量 |
| 18 | `internal/audio/device.go:126` | U1000 | `77538f1` | 09-19 | 存量 |
| 19 | `internal/ball/tokens_table_test.go:1109` | S1011 | `4ca2c6c` | 09-21 | 存量 |
| 20 | `internal/ball/tokens_table_test.go:1115` | S1011 | `cf8406e` | 09-21 | 存量 |
| 21 | `internal/config/parse.go:233` | U1000 | `fe126e2` | 09-19 | 存量 |
| 22 | `internal/llm/adaptertest/mockllm.go:58` | SA1019 | `5ddedf7` | 09-20 | 存量 |
| 23 | `internal/llm/adaptertest/registry_test.go:222` | U1000 | `5ddedf7` | 09-20 | 存量 |
| 24 | `internal/llm/anthropic/adapter.go:207` | SA4006 | `5ddedf7` | 09-20 | 存量 |
| 25 | `internal/llm/openaichat/adapter_test.go:39` | U1000 | `3b7bf24` | 09-19 | 存量 |
| 26 | `internal/llm/openaichat/mockllm_integ_test.go:52` | SA1019 | `3b7bf24` | 09-19 | 存量 |
| 27 | `internal/llm/ratelimit.go:200` | U1000 | `fe87e12` | 09-19 | 存量 |
| 28 | `internal/memory/artifacts_path_invariant_test.go:745` | U1000 | `d9224af` | 09-21 | 存量 |
| 29 | `internal/memory/privacy.go:262` | U1000 | `f27869e` | 09-19 | 存量 |
| 30 | `internal/models/archive.go:31` | U1000 | `92f4861` | 09-19 | 存量 |
| 31 | `internal/observe/goroutine_test.go:151` | SA4000 | `21a8738` | 09-19 | 存量 |
| 32 | `internal/observe/sampler_settle_gate_136_test.go:62` | U1000 | `aef82f5` | 09-24 19:08 | 存量（**今天但不在本窗口**，见 §2.3） |
| 33 | `internal/observe/thresholds.go:42` | U1000 | `739bb15` | 09-20 | 存量 |
| 34 | `internal/panel/attachments_test.go:33` | U1000 | `f4bf0fa` | 09-21 | 存量 |
| 35 | `internal/perm/ticket90_persist_test.go:59` | U1000 | `582b9a9` | 09-21 | 存量 |
| 36 | `internal/proc/crossvet_test.go:59` | SA1019 | `d1ee0a0` | 09-21 | 存量 |
| 37 | `internal/tools/fs_staging.go:122` | U1000 | `a0072b0` | 09-21 | 存量 |
| 38 | `internal/tools/fs_staging.go:130` | U1000 | `a0072b0` | 09-21 | 存量 |
| 39 | `internal/tools/helpers_test.go:24` | U1000 | `64c5fea` | 09-20 | 存量 |
| 40 | `internal/tools/mode.go:120` | U1000 | `1d4f289` | 09-21 | 存量 |
| 41 | `internal/tools/platform_other.go:39` | U1000 | `0986d63` | 09-20 | 存量 |
| 42 | `tools/d22scan/main.go:96`（日志里以 module 相对路径打印为 `main.go:96`） | U1000 | `bcf44d6` | 09-21 | 存量 |
| 43 | `tools/d22scan/main.go:103`（同上，打印为 `main.go:103`） | U1000 | `64d083d` | 09-20 | 存量 |
| 44 | `tools/mockllm/chat.go:219`（打印为 `chat.go:219`） | U1000 | `0839407` | 09-19 | 存量 |

**形状分布**：U1000（未用的字段／函数／类型／常量）37 枚、S1011（循环可写成 append）3 枚、
SA1019（弃用 API，三枚全是 `runtime.GOROOT`）3 枚、SA9009 1 枚、SA4006 1 枚、SA4000 1 枚、ST1012 1 枚。

### 2.3 新增＝零枚。三把尺子的第二、第三把（防"unused 类因删调用方而新造"）

第一把已过：44/44 的 blame commit 全是 `30e19ef` 的祖先（`git merge-base --is-ancestor` 逐枚 YES），
窗口内 35 枚没有一枚落在这些行上。

但**光看 blame 会漏**：`U1000 unused` 只要最后一枚调用方被删掉就会新造出来，
而被删的调用方可以在**别的文件**里。窗口内改过的 3 枚文件全在 package `observe`，
所以唯一可能的漏点是"本窗口删掉了 observe 包内某枚符号的引用"。逐枚两树对照（`git grep` 在
`30e19ef` 与 `a1fd5bf` 各数一次命中）：

| 符号 | 基数树 `30e19ef` | tip 树 `a1fd5bf` | 判 |
|---|---|---|---|
| `gateFailed`（表里第 32 枚） | 2 处命中：`:61` 注释＋`:62` 声明体，**零枚调用方** | 逐字相同 | 基数树上就已 unused，**存量** |
| `captureBufferLimitConversation`（第 33 枚） | 2 处命中：`:39` 注释＋`:42` 常量值，零枚调用方 | 逐字相同 | **存量**；另 `git diff 30e19ef..a1fd5bf -- internal/observe/thresholds.go` **空输出**，于是 本窗口一字节未碰那枚"一字节不许动"的文件 |
| `goroutine_test.go:151`（第 31 枚，SA4000 行内形状） | 该文件不在窗口 diff 名单内 | 同 | **存量** |

**第三把尺（外部对照，与本机仪器无关）**：台账 `docs/reports/pending-and-issues.md`（09-24 16:2x/16:4x 那节，
"A168/A169 编队"附近）已记 `lint` 红＝**44 条 staticcheck 积压**（`##[error]` 计数＝44），
落点名包括 `cmd/wisp/run.go:110/234/306`、`internal/agent/approval/queue.go:472`、
`internal/audio/device.go:45/51/126`、`frontend/embed.go:4` —— 与本趟 CI 自报 `findings=44` **数与名都对得上**。
所以这一步今天**既没变好也没变坏**，本窗口 35 枚对它的贡献是零。

### 2.4 一处**读数口径**必须点破（别把三个数混着用）

- 本趟 CI 原生数：**44**，口径＝该步自己的 `findings=` 自报（对每行做 `\((SA|ST|S[0-9]|QF|U)[0-9]+\)$` 计数），跨 3 module。
- 台账 09-24 的数：**44**，口径＝日志里 `##[error]` 行数。两个数在这一趟相等，但**它们是两把不同的尺**，
  换一趟（比如某步多出非 finding 的 error 行）就会分叉，今后引用要连口径一起引。
- `ci.yml:176` 内联那句"**37 lines across the 3 modules at d3cc9ed, linux shape**"**本程未能复算，也不当作废**：
  它是票 85a 在 `git archive HEAD` 快照上**本机**量的数，与 CI 原生尺不同源。本程只对它做了一件弱得多的核：
  44 枚里 blame 出在 `d3cc9ed` 之后的只有 **3 枚**（`4e5d240` 09-23、`6f702ae` 09-23、`aef82f5` 09-24），
  所以"37 到 44"那 7 枚的差额**不能全归到这三枚行上**——差异至少有一部分来自两把尺不同源。
  所以建议（不动手，只登记）：`ci.yml` 那行内联数该改成指向 CI 自报行而不是本机快照数。

### 2.5 五枚"不是纯风格"的，单拎出来给票 122（本程不修）

按本仓纪律，staticcheck 的账归票 122，且"不许为了变绿放宽断言"。这里只做**形状分类**，不开药：

1. `frontend/embed.go:4` SA9009 ineffectual `// go:embed` directive ——
   `ci.yml:171-174` 判它是**已知假阳性**（正文注释恰好以 `// go:embed` 开头，真指令在 `embed.go:19`）；
   台账同一节判它是"**真隐患不是风格**"。**两处分歧未并案，本程不裁**，只把它登记成"下一步该有人裁"。
   （票 77 AC#1 拥有那枚注释，票 122 裁这条，`ci.yml` 已写明它 KEEPS being reported，不设 `//nolint`。）
2. `internal/observe/thresholds.go:42` U1000 —— 落在**"一字节都不许动"的阈值文件**里。
   本程核过：该文件 `30e19ef..a1fd5bf` **零 diff**，所以这一趟**没人动它**；但这条 finding 的修法
   天然要碰那枚文件，所以**今后任何程要清这枚，必须先走人工批准**，不能当成一次普通 lint 清理。
3. `internal/observe/goroutine_test.go:151` SA4000（逻辑或两侧的表达式逐字相同）——
   在**D38 命名 goroutine 台账**的守卫测试里，两侧同表达式意味着这一判据**有一半是死的**。这不是风格，是**测试判别力**问题。
4. `internal/llm/anthropic/adapter.go:207` SA4006（赋给 `stop` 的值从未被读）—— 生产码，流式停止条件那一侧。
5. 三枚 `runtime.GOROOT` SA1019（`internal/llm/adaptertest/mockllm.go:58`、
   `internal/llm/openaichat/mockllm_integ_test.go:52`、`internal/proc/crossvet_test.go:59`）——
   Go 1.24 起弃用，**随工具链升级只会更响**，且第 3 枚在票 78 AC#3 那把跨-vet 尺上。

### 2.6 补记（本程自纠，追加不删）

`§2` 那枚 commit（`7cc5050`）**落盘时正文里还留着 5 枚 U+21D2 箭头连接符**，与本审计"输出不得含箭头/数学符号"
的纪律相违，也压在 ban #8 的码点段边界上（`docs/` 不在 ban #8 射程内，所以这不是 CI 违规，是本程违规）。
本程随后用词替换掉它们（"所以"／就地删），**不 amend、不改写已提交的那一枚**，修正随 §3 那枚 commit 一起落（同一枚文件、同一枚 pathspec）。
教训单记：**符号自查要在 commit 之前跑，不是之后**——本程在 §1 自查过一次，§2 新增正文后没再跑就交了。

---

## 3. `test-windows` 红因归因到测试名：两个口径都给，逐名对台账

### 3.1 本 run 的原始读数（job `107647320284`，两枚红步各自的 `runtests.sh` 自报行，逐字）

```
step 7 cmd/wisp CLI tests        : go test exited 1 - packages=[./cmd/wisp/ ...] top-level: PASS=50 FAIL=4 SKIP=0, === RUN=101, '[no tests to run]'=0
step 8 Portable windows tests    : go test exited 1 - packages=[./internal/proc/ ./internal/secret/ ./internal/config/ ./internal/risk/ ./internal/ball/ ./internal/perm/ ./internal/plugin/ ./cmd/llmrecord/ ...] top-level: PASS=266 FAIL=12 SKIP=1, === RUN=409, '[no tests to run]'=0
```

其余三枚测试步都是绿的：step 4 winsec `PASS=58 FAIL=0 SKIP=0／=== RUN=101`、step 6 cgo build smoke success、
step 9 PathResolver junction `PASS=1 FAIL=0 SKIP=0／=== RUN=1`。

### 3.2 失败测试名清单（16 枚，去重）

**step 7 `cmd/wisp CLI tests`（4 枚）**

| 测试名 | 落点 `file:line` | 耗时 |
|---|---|---|
| `TestTicket101ManualSwitchSurvivesRestart` | `cmd/wisp/run_mode101_test.go:207` | 3.80s |
| `TestTicket101UntouchedConfigRestartsAtDefault` | `cmd/wisp/run_mode101_test.go:328` | 6.63s |
| `TestTicket101SessionGrantDoesNotCrossRestart` | `cmd/wisp/run_mode101_test.go:389` | 4.71s |
| `TestComposedGateBlocksAWriteForTwoSeconds` | `cmd/wisp/run_test.go:330` | **301.08s** |

（第 4 枚那枚 **300.0x s** 按本仓已定的判据先怀疑 **C18 审批超时（300s 常量）**，不是性能回归——票 123 那一族，
台账 1390 行记的同一形状：`TestComposedGateBlocksAWriteForTwoSeconds` 耗时 301.06 秒＝正好撞满那个超时。）

**step 8 `Portable windows tests`（12 枚）**

| 测试名 | 落点 `file:line` |
|---|---|
| `TestClassifyAnchorSpellingIsNotVerdict` | `internal/risk/pathresolver_anchor_spelling_windows_test.go:29` |
| `TestCanonicalInputGainsNoSecondForm` | `internal/risk/pathresolver_anchor_spelling_windows_test.go:125` |
| `TestAListWinsWhereBothTablesHit` | `internal/risk/pathresolver_anchor_spelling_windows_test.go:204` |
| `TestPathResolverShortNameAListDenied` | `internal/risk/pathresolver_junction_windows_test.go:123` |
| `TestPathResolverUNCAListDenied` | `internal/risk/pathresolver_junction_windows_test.go:148` |
| `TestPathResolverExtendedLengthPrefixAListDenied` | `internal/risk/pathresolver_junction_windows_test.go:174` |
| `TestBListDefaultDenyAndOverride` | `internal/risk/pathresolver_junction_windows_test.go:272` |
| `TestSyncFixtureFallbackAndMatch` | `internal/risk/syncdirs_test.go:152` |
| `TestSyncFallbackNotDisarmableByWeakRoot` | `internal/risk/syncdirs_test.go:198` |
| `TestSyncUnverifiedRootKeepsFallback` | `internal/risk/syncdirs_test.go:224` |
| `TestSyncSuspectFallbackIsComponentBounded` | `internal/risk/syncdirs_test.go:344` |
| `TestSyncSuspectFallbackWhenUndetectable` | `internal/risk/syncdirs_test.go:390` |

### 3.3 两个口径的数（把口径写死，不许只引数字）

| 口径 | 本 run 36003984868 | 台账锚 run 35840958334（`4e9adcc`，09-23 09:06Z） |
|---|---|---|
| **A：`--- FAIL:` 的顶层测试名去重枚数** | **16** | **17**（本程独立复算，与台账逐字对上） |
| **B：含子测试路径的 `--- FAIL:` 行总枚数** | **16** | **21**（本程独立复算，与台账逐字对上） |
| A 与 B 之差＝嵌套子测试红行数 | **0**（本趟一枚嵌套子测试都没红） | **4**（全是 `TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/` 下的 `runTextTask`、`cmdModels`、`cmdProviders`、`cmdDoctor`；同族第 5 条腿 `resolveSecretLayout` 当时就是绿的） |
| 旁证：两步 `runtests.sh` 的顶层 FAIL 计数相加 | 4 加 12 ＝ **16** | 5 加 12 ＝ **17** |

所以本趟的诚实答案是 **16 与 16**，两数相等**只因这一趟没有任何嵌套子测试红**。
"17 与 21 是两个不同的诚实答案"那课在本趟的对应物是 **16 与 16**——
两个数相等本身是个事实，不是口径可以省掉的许可；换个有嵌套红的树它立刻分叉，所以口径必须继续一起引。

台账原文（`docs/reports/pending-and-issues.md` 09-24 那节，门禁读数第 2 条）写的是
"`test-windows` 红＝17 枚具名用例（`--- FAIL` 的名字去重＝17；含子项路径是 21 枚，两个数别混）"，
本程用 `gh run view --log` 对 run `35840958334`（job `107115610835`）逐行 `grep -c` 复算，
**17／21 两数复现成功**，不是〔仅自述〕。

### 3.4 逐名 存量／新增 判定：**新增 0 枚，存量 16 枚，另有一枚存量转绿**

台账锚那 17 枚与本趟 16 枚的差集，本程自己算了（不引台账的结论）：

- 17 减 16 ＝ 唯一少掉的那枚是 **`TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128`**（`cmd/wisp/dataroot_128_test.go`，5 条腿）。
- **它不是被跳过，是真变绿**（本仓规矩：判"不再红"先分清变绿还是被跳过）。日志正文三条都在：
  `=== RUN` 母项 1 枚＋5 条子腿全在，随后 `--- PASS` 母项 (0.04s) ＋ 5 条 `--- PASS` 子腿全在，
  且母项旁边还有真子进程诊断行（`dataroot_128_windows_test.go:69: AC#2 real process, leg "run" ... rc=2`、
  `leg "doctor" ... rc=1` 等）。所以 **变绿，非 SKIP，非没跑**。
- 归因到 commit：修好它的是 **`c2fa2e9`**（`test(wisp/128 AC#4): 进程内那五条腿钉住 WISP_ENV —— 先栽 runner 的 test 再 pin 回`，
  09-24 17:23），**在本窗口之外**（`4e9adcc..30e19ef` 之间，属今天更早的一批推送）。
  角色按 commit message 前缀 `test(...)` 读＝**实现程**；翻勾那枚是 `d56b6f5`（`docs(128,AC#4 翻勾 + AC#5 追加)`）。
  所以本审计**不认领这枚绿**，也不许它被算进 35 枚的成绩。
- 余下 16 枚逐名都在这 35 枚之外的存量集合里（16 枚的落点文件全部**不在**窗口 diff 名单——
  窗口只碰了 `internal/observe` 的 3 枚 Go 文件，而 16 枚落在 `cmd/wisp/run_mode101_test.go`、
  `cmd/wisp/run_test.go`、`internal/risk/pathresolver_*.go`、`internal/risk/syncdirs_test.go`）。
  台账对这一族各归各的票也已在场：`TestPathResolver*` 三枚＝`R-115-2` 的 8.3 短名族，
  `TestTicket101*` 三枚＝票 101 接线那一族，`TestComposedGateBlocksAWriteForTwoSeconds`＝票 123 那一族，
  `TestAC2*128`＝`128 AC#4/AC#5`（已清）。

所以 **`test-windows` 这一趟没有任何新增红；有且仅有 1 枚存量红被（别人、更早那批）修真了。**

### 3.5 一枚本趟新暴露、但**归因也是存量**的第二个红因（step 8 单独点名）

step 8 的自报是 `FAIL=12 SKIP=1`。那枚 `SKIP=1` 是 **`TestSyncRedTeamRealOneDrive`**
（`internal/risk/syncdirs_redteam_windows_test.go:205`，`:208` 处 `t.Skip("no profile home")`）。
判它要紧，因为 `tools/d22scan/runtests.sh` 的规矩是"**SKIP 不是 pass**"
（该脚本 `:24` 注释、`:88` 计数 `^--- SKIP`、`:99` 见非零即判死），所以：

- **step 8 今天有两枚彼此独立的红因**：12 枚 FAIL，加 1 枚**未登记的 SKIP**。
- 那枚 SKIP **不在** `scripts/portable-tests.sh` 的豁免账本里（`grep -n TestSyncRedTeamRealOneDrive scripts/portable-tests.sh`
  零命中；该步正文自报"11 ledger entries, 7 accounted on this platform"，7 条被点名的都不含它）。
- 归因：`:205`/`:208` blame 到 **`c091ee2b`（09-20 13:06）**，且同一枚 SKIP 在锚 run `35840958334` 的
  step 8 正文里**已经在场**（那趟同样是 `FAIL=12 SKIP=1`，八数与今天**逐字相同**：`PASS=266 FAIL=12 SKIP=1／=== RUN=409`）
  所以 **存量**，与这 35 枚无关。
- **但这条台账里没有**：台账那两条读数只数了 `--- FAIL` 的名字，没给 SKIP 记账。
  后果具名：**把那 12 枚 FAIL 全修完，step 8 仍然是红的**，因为那枚未登记的 SKIP 会自己把它判死。
  建议（本程不动手）：给 `TestSyncRedTeamRealOneDrive` 在 `scripts/portable-tests.sh` 补一行带理由的账，
  **或者**让它在 windows runner 上真跑；两条路都由实现程走，谁都不许把 `runtests.sh` 的 SKIP 判死放宽。

### 3.6 顺带纠一枚过期内联声明（`ci.yml`，非红、不改）

`ci.yml:466-471` 写着 `PathResolver junction placeholder` 这一步"DO NOT READ THIS WIRING AS A FIX OR A MASK:
the step **FAILS** on windows-latest today"，引的是 run `35551819606`。
本程核：run `35840958334`（09-23）与本 run（09-24）该步**两趟都是 success**，正文 `PASS=1 FAIL=0 SKIP=0／=== RUN=1`。
那句话是**过期声明**（引的是更早的 run），不是今天的事实；照 AGENTS.md 的规矩这里只登记不改文件——
"注释里说它红"和"CI 读数说它绿"同时在场时，**以 step 级读数为权威**。

---



## 4. 今天的 `internal/observe` 门禁改动有没有造成什么？

### 4.1 先确认改的是什么（读码，不读注释）

- 翻勾处：`internal/observe/sampler.go:556`，`const settleCoverageRowGates = true`，
  由 **`52191ce`**（`feat(136 AC#14 Gate): flip settleCoverageRowGates true + rewrite 2 coverage legs ...`，
  commit 时间 `2026-09-24 20:12:18 +0800`，＝ UTC 12:12:18）落地。
  `git show 52191ce -- internal/observe/sampler.go` 的该常量的增减行**只有那一枚布尔 false 变 true**，与票面"可执行码零行"之外的承诺一致。
- 它能走到的最远处（逐字引自 `sampler.go:546-555` 的注释自述，再用码面对一遍）：
  `buildSettleVerdicts`（`:566`）造一枚 `{Metric: "sampling", Pass: covered, Gate: settleCoverageRowGates}`，
  其中 `covered := rep.SampleErrors == 0 && len(rep.Samples) > 0`（`:567`）；
  `foldSettlePass`（`:592`）里 `if v.Gate && !v.Pass { pass = false }`，所以该行能否决 `SettleReport.Pass`；
  `rep.Pass = foldSettlePass(memOK && backInTime && releaseOK, rep.Verdicts)`（`:538`）。
- 到 CI 颜色的那一跳，我在 `scripts/slo-check.ps1` 里核到了，不是听注释的：
  `:381` `$allPass = ($failingStates.Count -eq 0) -and $settlePass`，`:396` `if (-not $allPass) { exit 1 }`。
**所以 settle 的 pass 位真的一票否决整枚 job 的颜色**，这条链在码面上是通的。

**方向判据（重要，决定"有没有造成什么"能怎么答）**：这枚翻勾**只会收紧不会放松**——
它只可能把 `pass=True` 变成 `pass=False`（把步弄红），**不可能**把红的弄绿。
所以"今天两步都是绿的"这一事实，配上这个单调性，等于说：**今天没有任何一步是被它弄红的。**

### 4.2 `slo-smoke`：真跑了取样，且**它今天的绿在翻勾之后比昨天更值钱**

`SLO smoke gate`（step 6）正文逐字（取法 `gh run view --log`，job `107647320562`）：

```
slo-check.ps1: sampling validity precheck (ticket 134 AC#4)
slo-check.ps1: precheck ok - no foreign toolchain/runner process, machine-wide cpu max 3%
slo-check.ps1: sampling state Sleeping for 4s
slo-check.ps1: state Sleeping exit=0 pass=True
slo-check.ps1: sampling state Warm for 4s
slo-check.ps1: state Warm exit=0 pass=True
slo-check.ps1: settle check (10s window)
slo-check.ps1: settle exit=0 pass=True
slo-check.ps1: leak fixture self-test (100MB, must FAIL)
slo-check.ps1: leak exit=1 flipped_to_fail=True
slo-check.ps1: report written to D:\a\wisp\wisp\build\slo\slo-report.json (all_pass=True)
```

判：**不是零样本 fail-closed，是真的采到了**。而且这一点**今天不需要相信它的措辞**——

- 翻勾之后 `covered` 要求 `len(rep.Samples) > 0 && rep.SampleErrors == 0`，且 `Gate=true`。
  零样本那一支（`:575` `"fail-closed disclosure: sample_errors=0 but 0 valid samples: this window measured nothing"`）
  会把 `Pass` 判 false，经 `foldSettlePass` 否决 `rep.Pass`，经 `slo-check.ps1:381/:396` 把这一步**直接弄红**。
- 它今天**是绿的**（`settle exit=0 pass=True`，步 conclusion=success，
  artifact `slo-smoke-report` **7179 字节**确实上传了，`Upload SLO report` 无 warning）。
- 所以"它真采到了样本"这一条，现在是**由门禁反推出来的结论**，不是由脚本自报的措辞接受的读数。
**所以这是今天这枚翻勾在 CI 上唯一一次真正咬合，而它咬合的方向是收紧**：同一枚绿，判据强度变了。
- 旁证：precheck 自报 `machine-wide cpu max 3%`（宿主不忙），两步 `state ... pass=True` 之间有真实墙钟推进
  （13:12:54 到 13:13:07 到 13:13:20 到 13:13:32，与 `-SecondsPerState 4` 加启动开销相容），
  `leak exit=1 flipped_to_fail=True` 说明泄漏自检那枚"必须红"的反向对照也真响了。

### 4.3 `slo-full`：**这道绿灯下什么都没做**——本节的靶子

`SLO full gate (six states + settle + leak)`（step 5，job `107647320611`，self-hosted `wisp-selfhosted-01`）
正文逐字：

```
slo-check.ps1: sampling validity precheck (ticket 134 AC#4)
slo-check.ps1: NO CONCLUSION (machine-contended) - subset=full refused to sample, no numbers were produced
slo-check.ps1: machine-contended reason: machine-wide cpu utilisation 99% over a 1s window (>= 50%)
slo-check.ps1: NO CONCLUSION (machine-contended): 1 reason(s), 0 state file(s) written, slo-report.json NOT written
slo-check.ps1: NO CONCLUSION (machine-contended): performance result - D32 (CPU <=0.5%, private
slo-check.ps1: NO CONCLUSION (machine-contended): RSS <=25MB) stays unverified for this run.
slo-check.ps1: NO CONCLUSION (machine-contended): record written to E:\work\base\actions-runner\_work\wisp\wisp\build\slo\slo-no-conclusion.json (NOT a report; slo-report.json is not written on this path)
slo-check.ps1: NO CONCLUSION (machine-contended): exit 0
```

三条独立物证钉住"它没做事却绿了"：

1. `0 state file(s) written, slo-report.json NOT written`，且它**自己声明** `D32 ... stays unverified for this run`。
2. `Upload SLO report` 步（step 6）正文：`##[warning]No files were found with the provided path: build/slo/slo-report.json. No artifacts will be uploaded.`
   ——该步 conclusion 仍是 `success`（因为 `ci.yml:597` 把 `if-no-files-found` 从 `error` 改成了 `warn`）。
3. 整枚 run 名下的 artifact 清单**只有一枚**：`slo-smoke-report`（7179 字节，13:13:45Z）。
   **没有 `slo-full-report`。** 全仓最新一枚 `slo-full-report` 是 **2026-09-24T10:18:08Z**，
   属 run `35985947228`（`7cf8075`），比这趟早 2 小时 53 分。所以上一枚真样品之后，
   `35996231784`（`30e19ef`，11:59Z）与本 run 连续两趟 D32 腿**都没产出数**。

**这不是 skipped，也不是 bug，是票 134 AC#6 的既定形状**（owner Q-36 拍板"都按推荐"：
有效性拒绝不再判红，改由 `scripts/slo-freshness.sh` 的探针 P3 在"最新上传的 report artifact"上计时来逼账）。
但**按本仓自己的纪律，它必须被点名成一条 finding**：

> `slo-full` 的 step 5 这一趟 conclusion=success，却**一枚样本都没取**、
> **一个 D32 数都没评**、**一份 report 都没写**。它是**"passed while doing nothing"那一形**，
> 而且它 step 层交出了颜色、产出了日志、零枚 skipped——
> 所以任何只看 `steps[].conclusion` 的读数器都会把它记成"D32 今天过了"。**今天 D32 没过，今天 D32 没被测量。**

同时给一枚**公平的反向记账**，免得这条被读成"AC#6 失败"：那 1 秒窗口 99% 的机器忙判定**大概率就是被这一批 CI 自己抢的**
（同 run 的 self-hosted `Build wisp.exe` 步在 13:13:26 到 13:17:33 之间跑了 4 分 07 秒，紧接着 13:17:33 取样），
加上本机同时在飞另一枚只读程。，形状是"自我污染被 precheck 正确拒掉"，不是"测量仪坏了"。
本仓记忆条目 `wisp-ci-selfhosted-topology` 写的"每次 push 自启抢 CPU"这一条，第一次拿到 **step 正文级**的实证。

### 4.4 那枚翻勾"打到了 CI 颜色"吗？**在它真正为谁而翻的那条腿上，到今天为止一次都没打到**

这是本节最该被读到的一句，我把判据摆全：

- 翻勾的动机（引 `sampler.go:546-550` 的自述，非逐字引用）：让覆盖行能否决 `slo -settle` 的 pass 位，
  并"打到 slo-smoke / slo-full 的 CI 颜色"。
- **本 run 是史上第一枚带上这枚 true 的 `ci` run**。依据：`52191ce` 落在 UTC 12:12:18；
  前一枚 dev run `35996231784` 的头是 `30e19ef`（11:59:10Z，**在翻勾之前**，`52191ce` 不在它树里）；
  `gh run list --limit 3` 现量证明**本 run 之后没有更晚的 `ci` run**。
- 而第一枚带上它的 run，**恰恰在 `slo-full` 那条腿上被 precheck 挡在取样之前**（§4.3）。
  门禁的否决逻辑住在 `buildSettleVerdicts`／`foldSettlePass`，而 AC#6 的拒绝路径在
  `slo-check.ps1` 的 **precheck**（"refused to sample, no numbers were produced"），
  发生在 settle 之前，所以**翻勾连被求值的机会都没有**。
- 结论两条，一分为二：
  1. **`slo-smoke`（hosted）这一腿：咬合了，且是有效咬合**（§4.2 的反推）。
     但它跑的是 `-Subset smoke`、只有 `Sleeping`＋`Warm` 两态，**不是 D32 那两个数（CPU 与私有 RSS）的评委**。
  2. **`slo-full`（D32 合并门禁）这一腿：翻勾至今 CI 颜色层面零次被求值。**
     码注释里那句 backing（"Nine clean settle runs backed it (sample_errors 0, exit 0, row gate=true)"）
     指向 `docs/evidence/s1/136-ac14b-impl.md` 第 4 节与 `136-ac14b-r2-acceptance.md` 第 5 节，
     **那些是本机/自建自采读数，不是 CI step 级颜色**，所以按本仓"说不出 run id 加 step 就当门不存在"的规矩，
     **这枚门禁目前只在一台 host 上存在过，除 smoke 那条窄腿外还没在 CI 上存在过。**
- **所以最小建议（本程不动手）**：**在机器安静时手动重跑一次带这枚 true 的 `slo-full`**
  （`workflow_dispatch`，或趁 fleet 睡觉时推一枚 docs），让它真取一次样并产出 `slo-full-report`。
  在那之前，任何"`slo-full` 绿＝D32 达标"的表述都不许写进任何交付面。
  这活该**编排者**做（要挑时间窗、要 push），不该实现程程做。
- 另注：别把它读成"要撤勾"：撤销口令是"撤 136 Gate 批准"（本程未使用、也不建议——它是收紧方向，
  今天的证据没有任何一条指向放水）。

### 4.5 本节的直接答案

**今天这 35 枚（其中唯一动可执行码的是那枚翻勾）没有把任何一步弄红；
它把 `slo-smoke` 那枚 settle 绿的证据强度提高了；
而它在 `slo-full` 上造成的唯一后果，是让我们撞见"那一步本来就什么都没做"。**
本节唯一硬 finding 属于 `slo-full`，且**它不是翻勾造成的，是 AC#6 的既定形状加上一次真实的机器争用**——
翻勾只是把那盏灯照得更亮了一点。

---

## 5. 历史对照（一个数，外加两枚必须分开的分母）

### 5.1 我实际读了多少、从什么时候起（先声明窗口，再报数）

仪器：`gh run list --workflow ci --limit 100`（**不带 branch 过滤**，即全分支全事件）。

| 项 | 值 |
|---|---|
| 实际读到的 run 数 | **100 枚**（`--limit 100` 就是这个工具的窗口上限，我没有翻第二页，所以这是"最近 100 枚"而不是"全部历史"） |
| 窗口 | **2026-09-21T12:31:37Z**（最旧，run `35600043583`）到 **2026-09-24T13:11:38Z**（最新＝本审计对象 `36003984868`），跨度约 72 小时 |
| 窗口内出现的分支 | `["dev"]` 一枚，**没有第二枚分支**（`main` 在本仓不存在，`master` 那支死种子从不触发本 workflow，与记忆条目 `wisp-ci-selfhosted-topology` 一致） |
| 结论分布 | `failure` 92 枚 ＋ `cancelled` 8 枚 ＝ 100 枚 |
| **`conclusion=success` 的枚数** | **0** |

所以 **台账那句"`ci` 最近 100 枚 0 枚成功"今天仍然成立**，但要带两条限定才能继续引用：
其一，台账那次是**截至 2026-09-23**的 100 枚，我今天这 100 枚窗口**整体前移了约 72 小时**，
**是不同的一批 100 枚**（我今天这窗口把 09-21 那批挤了出去一部分、又收进 09-24 那 21 枚）；
其二，8 枚 `cancelled` 不是红也不是绿，它们**没有任何 step 级读数可采**（与 skipped 同性质），
所以"0 枚 success"里含 8 枚"根本没判"的格子——这条台账原文没区分，今后引这数要连它一起说。

### 5.2 run 级与 step 级是两枚不同的分母（这一节存在的理由）

同一个窗口里，下面两句**同时为真**：

- "`ci` 工作流最近 100 枚 **0 枚成功**。"（分母＝run）
- "`slo-full` 在 09-24 一天 21 枚 dev run 里，门禁步 **21 枚全绿**。"（分母＝step）

我把 09-24 全天 21 枚 dev run 逐枚 `gh api runs/<id>/jobs` 与 `runs/<id>/artifacts` 拉了一遍，现量：

| 量 | 09-24 现量（21 枚 run 全查，零外推） |
|---|---|
| `ci` run 级 conclusion | 21 枚 **全 failure** |
| `slo-full` **job** 级 conclusion | success 19 ／ failure 2 |
| `slo-full` **step 5（SLO full gate）** conclusion | **success 21 ／ failure 0** |
| 那 2 枚 job failure 的红因 | **都是 step 6 `Upload SLO report`**（`##[error]Failed to FinalizeArtifact: Unable to make request: ECONNRESET`），**不是 D32 没过** |
| 真取了样且 `all_pass=True`（判据：日志有 `report written ... (all_pass=True)`；含 9 枚成功上传 artifact ＋ 2 枚 report 写成了但上传被 ECONNRESET 吞了） | **11 枚** |
| **零取样、走 `NO CONCLUSION (machine-contended)` 而 step 5 仍绿** | **10 枚** |

**"今天这枚 step 绿"与那枚 claim 的关系**（claim＝"slo-full 自改造以来已绿过两枚"）：

1. **数量级不对**：光 09-24 一天，`slo-full` 的门禁步就绿了 **21 枚**、job 绿了 **19 枚**，
   不是"两枚"。更早的 134 AC#6 证据面（`docs/evidence/s1/134-ac6-contended-no-conclusion.md:812-826`）
   在 09-23 就已经记了"带 AC#6 代码的三枚 run 3/3 success，它之前那批 8/9 failure"。
   所以那句"绿过两枚"要么是在说**别的分母**（比如"绿过两枚**且真取到样**"），要么是当时的一次局部读数。
   本程**不去替它圆**，只把三个分母的现量摆出来：run 0／job 19／step 21，同日真取样 11、拒采 10。
2. **今天这一枚与那枚 claim 的口径相容性**：本 run 的 step 5 绿是**那 10 枚"没测任何东西"的绿之一**（§4.3 三条物证），
   所以它**不能**被算进"真过 D32 的绿"。真过 D32 且留了 artifact 的最新一枚是
   run `35985947228`（`7cf8075`，report 建于 **2026-09-24T10:18:08Z**，比本审计对象早 **2 小时 53 分**）。
3. **一个必须点破的仪器坑（我这次差点踩）**：**"没有 `slo-full-report` artifact"不等于"没测"**。
   有两枚 run（`35948947994`、`35946437702`）日志明写 `report written ... (all_pass=True)`，
   是 `Upload SLO report` 撞了 `ECONNRESET` 才没留 artifact；另外 gate 步真红时上传步会被整枚跳过，也不留 artifact。
   所以今后**用 artifact 数当"真样品种数"的代理**必须叠加"读 step 5 正文"这一道，
   单用 artifact 会把 11 枚数成 9 枚，而**探针 P3 恰恰就是用 artifact 计时的**（票 134 AC#6 的设计）。
   这不是假想风险：那两枚 run 的 D32 结果对 P3 而言**已经丢了**。

### 5.3 本节的一句话

**run 级 0／100 与 step 级 21／21 是同一天同一批 run 的两个真话**；
`slo-full` 今天的绿灯与"改造后绿过两枚"那说法**不相容地低估了一个量级**，
而它与"今天真过了 D32"那说法**完全不相容**——今天它一枚样本都没取。

---
