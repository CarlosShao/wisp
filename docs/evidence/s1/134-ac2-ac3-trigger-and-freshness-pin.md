# 134 AC#2 + AC#3 取证 — `schedule` 触发落地、真 run 读数、以及反静默死钉的变异三态（2026-09-23）

被验版本：`e951dfa`（AC#4）+ 本文件所在提交（AC#2/AC#3），锚定 sha `ac6f31c`。
本票**没有 push 过任何东西**（只 `git fetch origin dev` 读了一次），下面 §1 那枚真 run 是
**别人的 commit 把我的 commit 一起带上去**的，出处见同节的祖先关系读数。

## 1. AC#2：真跑出的四项（run id + job id + step 名 + 结论）

| 项 | 值 |
|---|---|
| run id | **35810714576**（`ci`，event=push，head `58302cc`，head_branch=dev，created 2026-09-23T02:31:45Z） |
| job id | **107021435150**（job name `slo-full`，`runs-on: [self-hosted, wisp-slo]`） |
| step 名 | **`SLO full gate (six states + settle + leak)`**（第 5 步；同 job 另两步 `4 Build wisp.exe (deps cached on the runner)`、`6 Upload SLO report`） |
| 结论 | **success**（job started 02:31:49Z → completed 02:34:30Z，162 s；第 4/5/6 步逐条 conclusion=success） |

**为什么这一枚算"改动之后"的读数**（不是拿旧 run 冒充）：

```
git merge-base --is-ancestor e951dfa 58302cc   -> rc=0（我的 AC#4 commit 是它的祖先）
git branch -r --contains e951dfa               -> cnb/dev, origin/HEAD -> origin/dev, origin/dev
git rev-parse e951dfa:scripts/slo-check.ps1    -> 560186fa64ccdc0fbbabf53266f6a06e7aab3d36
git rev-parse 58302cc:scripts/slo-check.ps1    -> 560186fa64ccdc0fbbabf53266f6a06e7aab3d36   （逐字节同一枚 blob）
```

⇒ 那枚 self-hosted run 跑的 `scripts/slo-check.ps1` **就是**加了前置检查之后的版本，
且六个 job 全跑完、`slo-full` 结论为 `success` ⇒ **前置检查没有把正常的门打死**（这是 AC#2 真正要说的那句话）。

### 1a. `A103④` 三连检查（取证侧也走了一遍，10:32:48 +08 的读数）

```
(1) gh run list --branch dev --status in_progress
    -> 35810714576  2026-09-23T02:31:45Z  push  58302cc            （1 枚在飞）
(2) E:/work/base/actions-runner/_diag/Worker_*.log mtime
    -> Worker_20260923-023150-utc.log  2026-09-23_10:32:46  （now=10:32:48 ⇒ 2 秒前还在追加）
(3) 路径落在 actions-runner\_work\ 下的进程
    -> wisp.exe pid=41692 / pid=31328 / pid=6724  E:\work\base\actions-runner\_work\wisp\wisp\build\wisp.exe
```

⇒ 三条都写着"CI 正在本机取样"。所以：**这一轮之内我没有取任何时序/RSS 读数**，
只读 API；我在 10:31:49Z(=10:31:49 +08) 之前跑的并发负载实验（`go vet -a`，CPU 100%，10:27:17–10:29）
**早于**该 job 起跑，10:34:30 之后才继续跑变异测试（`sh`/`grep` 单次 <10 ms）⇒ 没有污染这枚样本。

### 1b. AC#2 尚未拿到的一格（照实报，不拿本地绿替代）

- **步级日志文本取不到**：`gh run view --job=107021435150 --log` 回
  `run 35810714576 is still in progress; logs will be available when it is complete`
  （10:37:04 时 `test-windows` 仍 in_progress）⇒ "CI 那一步里 `precheck ok` 那行到底长什么样"
  **未验证**，等 run 整体收尾后回读（票面 Progress log 记了欠账）。
- **`schedule` 触发本身零次执行**：见 §2。

## 2. `on:` 改动与 YAML 解析器复核

改动后的 `on:`（`.github/workflows/ci.yml:10-38`）与 `slo-full` 关键字段：

```yaml
on:
  pull_request:
  push:
    branches: [main, dev]
  schedule:
    - cron: '37 19 * * *'      # 19:37 UTC = 03:37 +08，编队睡觉的时段
  workflow_dispatch:
```

`slo-full`：`runs-on: ['self-hosted', 'wisp-slo']`、`env` 键 `['MINGW64_ROOT','WISP_ENV']`、
`steps = 5`、`if`/`continue-on-error` **均不存在**（解析器读出来的，不是读注释）。

三行仪器复核（`python -c "import yaml..."` 那套，逐字）：

```
PARSER-OK; job count = 6
jobs = lint, test-core, test-windows, slo-smoke, slo-full, lint-frontend
all six present = True (missing: [] extra: [] )
runs-on 'ubuntu-latest' = 3 ; runs-on 'windows-latest' = 2 ; runs-on '- self-hosted\n- wisp-slo' = 1
50:  group: ci-${{ github.workflow }}-${{ github.ref }}${{ github.event_name == 'push' && format('-{0}', github.sha) || '' }}
51:  cancel-in-progress: ${{ github.event_name == 'pull_request' }}
```

`ci.yml` 里非注释的改动行只有两行（其余 33 行是注释）：

```
$ git diff -U0 -- .github/workflows/ci.yml | grep -E "^[+-]" | grep -vE "^(\+\+\+|---)" | grep -vE "^[+-][[:space:]]*#"
+  schedule:
+    - cron: '37 19 * * *'
```

⚠ **如实写清、不许当成已生效**：GitHub 的 cron **只对默认分支（`main`）上的配置求值**。
`dev` 上这条 `schedule` 现在**触发次数 = 0**，要等这枚改动进 `main` 才会真跑；
并且 scheduled run 是 best-effort（高峰期会延迟），且会跑**整枚 `ci` workflow 的 6 个 job**
（收窄它需要 `if:`，D22 mode-6 不许）。这三条限制都写进了 `ci.yml` 的注释里。

## 3. AC#3：反静默死钉 = `scripts/slo-freshness.sh` + `.github/workflows/slo-fresh.yml`

钉不能住在 `slo-full` 自己身上（那枚 job 只有被触发才跑 ⇒ 循环），也不能和 `ci.yml` 共用触发器
（那它跟门一起死）。所以它自带时钟：`slo-fresh.yml` = `schedule: '23 */6 * * *'` + `workflow_dispatch`，
`runs-on: ubuntu-latest`，`permissions: {contents: read, actions: read}`，无 `if:` / 无 `continue-on-error`
（解析器：`jobs= ['slo-fresh']`、`has if= False`、`has continue-on-error= False`）。

两道探针（都要过）：**P1 静态**——`ci.yml` 仍必须有顶层 `on:`、`push:`、`branches` 里的 main/dev、
一条 `cron:`，且 `slo-full` 那枚 job 体内不许出现 `if:` / `continue-on-error`、必须仍跑
`-Subset full`；**P2 动态**——最近的 `ci` run 里那枚 `slo-full` job 的 `created_at` 必须年轻于
`SLO_FULL_MAX_AGE_DAYS`（默认 3 天），未完成的 job 排队 >12 小时也算死。
稳定失败 token：`slo-full-trigger-missing` / `slo-full-never-triggered` / `slo-full-stale`。

### 3a. 变异三态（红→绿→红，逐发还原；P2 与 P1 各一组）

记录用的假"现在"固定为 `SLO_FRESH_NOW=2026-09-23T02:34:00Z`；真记录取自
run `35809646757` / job `107018099910`（`created_at=2026-09-23T02:15:46Z|completed|success`）。

| 发 | 动作（落地证据） | 期望 | **实测 rc + 原文关键行** |
|---|---|---|---|
| S1 | 基线：真记录、阈值 3 天 | 绿 | **rc=0** `slo-freshness: OK - slo-full still has automatic triggers and a trigger record inside the window` |
| S2 | 喂 10 天前的记录（`created_at=2026-09-13T02:15:46Z`） | **红** | **rc=1** `FAIL slo-full-stale: no slo-full trigger record for 10 day(s) (> 3). ... D32 has had no automatic verdict since 2026-09-13T02:15:46Z.`（`age: 10 day(s) (865094 s); threshold: 3 day(s)`） |
| S3 | **临时放宽判据** `SLO_FULL_MAX_AGE_DAYS=30`（同一枚旧记录） | 绿 | **rc=0** `age: 10 day(s) (865094 s); threshold: 30 day(s)` ⇒ 红确实是这条阈值造成的，不是常数红。落地证明：`grep -n 'max_age_days=' scripts/slo-freshness.sh` → `47:max_age_days=${SLO_FULL_MAX_AGE_DAYS:-3}`（**文件本身一个字没改**，放宽只走环境变量） |
| S4 | 还原（去掉那个环境变量） | **红** | **rc=1** 同一句 `slo-full-stale ... 10 day(s) (> 3)` ⇒ 逐发还原完成 |
| S5 | **P1 结构变异**：`ci.yml` 的 /tmp 副本里删掉 `  schedule:` + `    - cron:` 两行（落地证据：`cron lines in copy: 0 ; schedule lines in copy: 0 ; real file: 1`） | **红** | **rc=1** `FAIL slo-full-trigger-missing: /tmp/wisp134/ci_nosched.yml has no 'cron:' schedule entry (ticket 134 shape B); nobody has to press a button any more` |
| S6 | 还原（指回真文件） | 绿 | **rc=0** `slo-freshness: OK ...` |
| S7 | **P1 mode-6 变异**：副本里给 `slo-full:` 之后插一行 `      if: github.event_name == 'push'`（落地证据 `grep -n -A1 "^  slo-full:" → 525:  slo-full: / 526-      if: ...`） | **红** | **rc=1** `FAIL slo-full-trigger-missing: the slo-full job now carries 'if:' (D22 mode-6: a skippable job is a job that will one day be skipped)` |
| S8 | P2 "从未触发"分支：记录 `created_at` 留空 | **红** | **rc=1** `FAIL slo-full-never-triggered: the record has no created_at: \|completed\|none\|...` |
| S9 | P2 "runner 装死"分支：`status=queued`、13 小时前 | **红** | **rc=1** `FAIL slo-full-stale: the newest slo-full job is still 'queued' after 13 hour(s) - a queued job that never starts is a silent death too.` |
| S10 | **P1 结构变异**：副本里删掉 `  push:` 一行（落地证据 `push lines: copy=0 real=1`） | **红** | **rc=1** `FAIL slo-full-trigger-missing: /tmp/wisp134/ci_nopush.yml has no 'push:' trigger; slo-full can only run when a human presses something` |
| S1' | 全部还原后再跑一遍基线 | 绿 | **rc=0** `slo-freshness: OK ...`；`grep -nE "^      if: github.event_name == 'push'$" .github/workflows/ci.yml` → **rc=1（命中 0，真文件无残留）** |

⚠ 三件照实说的：
① S5/S7/S10 的**结构变异只在 `/tmp` 的 `ci.yml` 副本上做**（脚本留了 `SLO_FRESH_CI` 这个路径注入口），
   共树里不去动别人的文件；代价是这三发测的是"副本红"，不是"工作树红"。
② `shellcheck` 本机**没装**（`command -v shellcheck` rc=1）⇒ 本地只做到 `sh -n scripts/slo-freshness.sh` rc=0；
   shellcheck 那一格等 CI 的 `Shell lint for the pin` 步（ ubuntu 镜像自带它；缺了那一步**硬红**，不静默放行）。
③ P2 的"真取 GitHub"这条腿**本地未跑**（我的 `gh` 用 keyring 鉴权，脚本要求 `GH_TOKEN`/`GITHUB_TOKEN` 非空，
   否则 exit 2 而不是假过）⇒ 那一腿的读数等 `slo-fresh.yml` 在 CI 上第一次执行，**现在说不出它的 run id**；
   本节所有 P2 读数都走的是记录注入缝（`SLO_FULL_LAST_TRIGGER`），用的值是真 API 值。

## 4. 两个阈值未动的自证（AC#5 一并）

- 本文件涉及的改动面：`.github/workflows/ci.yml`（`on:` + `slo-full` 注释/`schedule`）、
  `scripts/slo-check.ps1`（纯插入 `174 0`）、新增 `scripts/slo-freshness.sh`、新增 `.github/workflows/slo-fresh.yml`。
  `git status --porcelain` 里 `internal/**`、`cmd/**`、`tools/d22scan/**`、`allowlist.txt`、`docs/PLAN.md`、
  `docs/specs/*.md` **一个都没有**。
- 阈值原文（`internal/observe/thresholds.go`，只读）：

```
	memCapSleeping     int64 = 25 << 20
	cpuLimitSleeping     = 0.5 // % of all-core mean, 1min window
```
