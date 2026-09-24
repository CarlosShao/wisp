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
| `test-windows` step 16 `Post Cache third_party (deps.toml key)` | "本次 push 是否把 `third_party` 存回缓存"读数永不可采 | step 5 是 cache **命中**（restore-keys 命中，所以 无新 key 可存），`actions/cache` 命中路径不跑 save。**不是门禁洞**；但记账要如实：`a1fd5bf` 这一树的 cache-save 证据在这一趟是**采不到**而非"采到且为空" |
| `test-windows` step 17 `Post Run actions/setup-go@v5` | 同上，post 阶段读数永不可采 | 与 `lint` 那枚同形，**不是门禁洞** |

**这一节的总裁（step 层）**：本 run 的红**只落在 3 枚 step**：`lint/9 staticcheck`、
`test-windows/7 cmd/wisp CLI tests`、`test-windows/8 Portable windows tests`。
三枚 skipped 均非门禁，**门禁读数零丢失**；`slo-full` 的问题相反——它 step 层绿且产了日志，
但正文自陈零取样（第 4 节）。"a run being red is not a finding; a step being red is" 在这一趟
按字面成立；再加一条本趟的教训：**step 绿也不是读数，正文才是**。

---
