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
| 6 `gofmt (gofumpt)` | success | `gofumpt -l` 输出为空（该步的判据就是"输出非空即 exit 1"）⇒ **格式债为零枚** |
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

⇒ 这枚门禁**真的执行完了**：`toolchain-crash-lines=0`，票 85a 那枚"export data version 4 is greater than
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
⇒ 这一步今天**既没变好也没变坏**，本窗口 35 枚对它的贡献是零。

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
   台账同一节判它是"**真隐患不是风格**"。⇒ **两处分歧未并案，本程不裁**，只把它登记成"下一步该有人裁"。
   （票 77 AC#1 拥有那枚注释，票 122 裁这条，`ci.yml` 已写明它 KEEPS being reported，不设 `//nolint`。）
2. `internal/observe/thresholds.go:42` U1000 —— 落在**"一字节都不许动"的阈值文件**里。
   本程核过：该文件 `30e19ef..a1fd5bf` **零 diff**，所以这一趟**没人动它**；但这条 finding 的修法
   天然要碰那枚文件，⇒ **今后任何程要清这枚，必须先走人工批准**，不能当成一次普通 lint 清理。
3. `internal/observe/goroutine_test.go:151` SA4000（逻辑或两侧的表达式逐字相同）——
   在**D38 命名 goroutine 台账**的守卫测试里，两侧同表达式意味着这一判据**有一半是死的**。这不是风格，是**测试判别力**问题。
4. `internal/llm/anthropic/adapter.go:207` SA4006（赋给 `stop` 的值从未被读）—— 生产码，流式停止条件那一侧。
5. 三枚 `runtime.GOROOT` SA1019（`internal/llm/adaptertest/mockllm.go:58`、
   `internal/llm/openaichat/mockllm_integ_test.go:52`、`internal/proc/crossvet_test.go:59`）——
   Go 1.24 起弃用，**随工具链升级只会更响**，且第 3 枚在票 78 AC#3 那把跨-vet 尺上。

---

