# 134 AC#6 取证 — `machine-contended` 改判「本 run 无结论」+ 新鲜度钉改按「最近一次产出有效样本」计龄

被验版本：锚定 sha `b723978c0f17612c5e3210670fa6da9db8e94ff`（`git rev-parse HEAD`，2026-09-23 13:33 +08 起锚；
本程期间共树另有代理把 HEAD 从 `717d822` 推到 `a35f611` 再推到 `b723978`，四枚可写件的 blob 我按锚点记死如下）：

| 文件 | 锚点 blob | 本格是否动 |
|---|---|---|
| `scripts/slo-check.ps1` | `560186fa64ccdc0fbbabf53266f6a06e7aab3d36` | 动（前半：上色） |
| `scripts/slo-freshness.sh` | `5a4a339a9f593d48c767f1d98dd4156ba380ec35` | 动（后半：计龄依据） |
| `.github/workflows/slo-fresh.yml` | `872a0c4d313ae05626695b1abca4f010b004d1f2` | 动（只补注释：新失败令牌） |
| `.github/workflows/ci.yml` | `b2cd7632b8dd463c0d83a1e6cbeaec51857da953` | 动一处，**超出简报的「那一步」**，见 §0.5 报备 |
| `internal/observe/thresholds.go` | 未打开、未改 | 禁改（§4 给 grep 证明） |

## 0. 前置取证：今天 `machine-contended` 确实让**整枚 job** 红（简报四项齐，且逐条复核）

简报给的线索是**断言**，我独立复算：`gh api repos/CarlosShao/wisp/actions/runs/35817761098/jobs` +
`gh api repos/CarlosShao/wisp/actions/jobs/107042866358/logs`（原文落 `/tmp/134ac6/job-107042866358.log`，240 行）。

- **run id**：`35817761098`（`ci` / `dev` / `push`，`head_sha df4a60a9c6d50c819335fb2fd1889d215491fddb`，
  该 sha 的 `scripts/slo-check.ps1` blob = `560186fa64cc…` = 我锚点的同一枚 blob，逐字节相同 ⇒ 下面那段日志就是这份代码的读数）
- **job id**：`107042866358`（name `slo-full`，`status=completed`，`conclusion=failure`）
- **step 名**：第 **5** 步 `SLO full gate (six states + settle + leak)`，`conclusion=failure`；
  同 job 的 1-4 步全 `success`，第 6 步 `Upload SLO report` = `skipped`（默认 `success()` 语义，不是 `if:`）
- **结论**：整枚 job `failure`（`##[error]Process completed with exit code 1.`）

步内原文（逐字，时间戳保留）：

```
2026-09-23T04:18:38.5229166Z ##[group]Run powershell -NoProfile -ExecutionPolicy Bypass -File scripts/slo-check.ps1 -Subset full -SecondsPerState 6
2026-09-23T04:18:39.6622868Z slo-check.ps1: sampling validity precheck (ticket 134 AC#4)
2026-09-23T04:18:43.3837014Z slo-check.ps1: FAIL machine-contended - subset=full refused to sample, no numbers were produced
2026-09-23T04:18:43.3856111Z slo-check.ps1: machine-contended reason: foreign toolchain/wisp process present: go.exe pid=43684 started=2026-09-23 12:10:25 path=D:\work\base\go\bin\go.exe
2026-09-23T04:18:43.3869181Z slo-check.ps1: machine-contended reason: machine-wide cpu utilisation 61% over a 1s window (>= 50%)
2026-09-23T04:18:43.3882696Z slo-check.ps1: machine-contended: 2 reason(s), 0 state file(s) written, slo-report.json NOT written
2026-09-23T04:18:43.3898047Z slo-check.ps1: machine-contended: this is a SAMPLING VALIDITY verdict, not a performance
2026-09-23T04:18:43.3910652Z slo-check.ps1: machine-contended: result - D32 stays unverified for this run. Rerun when the
2026-09-23T04:18:43.3922197Z slo-check.ps1: machine-contended: machine is quiet (ticket 134 AC#4, shape C: refuse loudly
2026-09-23T04:18:43.3934374Z slo-check.ps1: machine-contended: instead of reporting a contaminated sample as truth).
2026-09-23T04:18:43.6918568Z ##[error]Process completed with exit code 1.
```

⇒ 简报那四项**成立**，动码许可生效。

### 0.1 「每次推送几乎必红」这句断言的实测口径

`gh api repos/CarlosShao/wisp/actions/runs?per_page=30` 里最近 25 枚 `ci` run 的 `slo-full` job：
**13 枚 success / 12 枚 failure**；只看 09-23 当天（19 枚）：**11 success / 8 failure**。
⇒ 「每推必红」的字面意义**不成立**（不到一半）。但把窗口收到争用激烈的最近 7 枚 push（04:14Z 之后）：
**6 枚 failure / 1 枚 success**。

09-23 这 8 枚 failure 的**失败步逐枚点到同一格**
（`gh api repos/CarlosShao/wisp/actions/jobs/<id>` → `[.steps[]|select(.conclusion=="failure")|.name]`）：

```
job=107056192824 -> SLO full gate (six states + settle + leak)
job=107049793053 -> SLO full gate (six states + settle + leak)
job=107046632390 -> SLO full gate (six states + settle + leak)
job=107042866358 -> SLO full gate (six states + settle + leak)
job=107042519143 -> SLO full gate (six states + settle + leak)
job=107042238141 -> SLO full gate (six states + settle + leak)
job=107038116636 -> SLO full gate (six states + settle + leak)
job=107008457499 -> SLO full gate (six states + settle + leak)
```

逐枚 `gh api .../actions/jobs/<id>/logs | grep -c machine-contended`：
`12 / 8 / 7 / 8 / 13 / 11 / 7 / 0` ⇒ **前 7 枚全是争用拒采样**（没有一枚是 D32 数字不过），
最后那枚 `107008457499`（run `35806505339`，01:29Z）争用行 **0** —— 它跑的是 AC#4 之前的旧 blob
`fcebe28828892d6f498bf2f4b46446729fb2924b`，六态真取了样、`state Sleeping exit=2 pass=False`、
`report written … (all_pass=False)` ⇒ **那是"出了数且数不过"的一枚，正是本格不许改判的那一形**（见 §0.4）。
**判**：断言在「编队活跃的窗口里几乎必红」这个意义上成立（最近 7 枚 push 6 红、红的因全是争用），
在「每推必红」的字面意义上不成立（当天 11:8）；本格要治的是前者（红海吞信号），后者的账不改判据。

### 0.4 一条必须留住的形：出了数而数不过 ⇒ 仍然判红

`35806505339` 那一枚的步内尾部（逐字）：

```
2026-09-23T01:30:54.5872947Z slo-check.ps1: state Sleeping exit=2 pass=False
2026-09-23T01:32:11.7565447Z slo-check.ps1: state WorkPeak exit=0 pass=True
2026-09-23T01:32:23.9196902Z slo-check.ps1: settle exit=0 pass=True
2026-09-23T01:32:35.2737353Z slo-check.ps1: report written to E:\work\base\actions-runner\_work\wisp\wisp\build\slo\slo-report.json (all_pass=False)
2026-09-23T01:32:35.4402151Z ##[error]Process completed with exit code 1.
```

⇒ 本格的改动**只落在 `if ($reasons.Count -gt 0)` 那一块**（争用分支）；上面这条路（写了 report、
`all_pass=False` ⇒ `exit 1`）一个字节不动，§1 的 `git diff --numstat` 与 YAML 解析器读数一起证这一点。

### 0.2 顺带复核简报另外两条断言（一条不成立，如实报）

1. 「`slo-freshness.sh` 的 P2 今天查的是最近一枚 `slo-full` **job** 的年龄 > 3 天即红」——**成立**：
   `scripts/slo-freshness.sh:144-154` 逐枚 run 取 `select(.name=="slo-full")` 的 `created_at`，
   `:178` `[ "$age_days" -gt "$max_age_days" ]`，`:47` `max_age_days=${SLO_FULL_MAX_AGE_DAYS:-3}`。
   ⇒ 判据物是 **job 记录**，与「这次有没有取到样」无关。job 一改成 exit 0 就永远新鲜，这条我复核后同意。
2. 台账 `Q-36` 行写「`slo-fresh.yml`（`slo-full-stale`，**10 天**旧记录就红）」——**不成立**：真默认是 **3 天**
   （`scripts/slo-freshness.sh:47`），「10 天」只是 AC#3 变异时喂进去的那枚记录年龄（票面 AC#3 那格写的就是「喂 10 天前的记录 rc=1」）。
   ⇒ 我按 **3 天** 办事，没有为此改任何代码；这条只登记，不回改台账（append-only）。

### 0.3 「有效样本」判据物的选型：读到的东西必须是 API 上有的

原判据物（job 年龄）不够用了，换什么？候选与实测：

| 候选判据物 | API 是否取得到 | 判定 |
|---|---|---|
| job 的 `conclusion` | 是 | **排除**：简报明禁，且 exit 0 之后恒绿 |
| job 步内日志文本（`report written to …`） | 是（`/actions/jobs/{id}/logs`） | 排除：要逐枚下载日志（今天 25 枚 × 数 MB），且本质仍是"从 job 反推" |
| `GITHUB_STEP_SUMMARY` | 否（无稳定 REST 读取面） | 排除 |
| **workflow artifact `slo-full-report` 的存在与 `created_at`** | 是（`/repos/{o}/{r}/actions/artifacts`） | **采用** |

采用理由（全部实测于 13:2x-13:3x，`gh api`）：这份 artifact 的 `path:` 就是 `build/slo/slo-report.json` 本身——
**脚本只在真取到样时才写那枚文件**（争用路径 `0 state file(s) written, slo-report.json NOT written`，且进预检前先把上一轮
遗留的 `*.json` 清掉，`scripts/slo-check.ps1:104-112`）⇒ artifact 在不在 = 那枚 report 在不在 = 这次是不是真出了数。

实测两枚对照：

```
# 出了样的那枚（run 35810714576 / job 107021435150 / step 5 success）
$ gh api repos/CarlosShao/wisp/actions/runs/35810714576/artifacts
{"created_at":"2026-09-23T02:34:22Z","id":10729737755,"name":"slo-full-report","size_in_bytes":21711}
{"created_at":"2026-09-23T02:34:49Z","id":10729637897,"name":"slo-smoke-report","size_in_bytes":6988}

# 争用那枚（run 35817761098 / job 107042866358 / step 5 failure）
$ gh api repos/CarlosShao/wisp/actions/runs/35817761098/artifacts
{"created_at":"2026-09-23T04:21:29Z","id":10732580634,"name":"slo-smoke-report","size_in_bytes":7113}
（没有 slo-full-report 这一行）
```

仓库级一枚（P3 走的就是这条，单次调用可分页）：`GET /repos/{o}/{r}/actions/artifacts?per_page=100` →
`total_count=168`，其中 `name=="slo-full-report"` **43 枚**，`expired=false`，最新一枚
`created_at=2026-09-23T04:44:08Z / size_in_bytes=21474 / workflow_run.id=35819355656`。

⚠ 一条**必须写进脚本**的实测事实：这枚列表**不是按 `created_at` 排序的**（第 1 页前 9 行是
`05:25 / 04:53 / 04:44 / 04:38 / 04:21 / 03:43 / 03:56 / 04:18 / 03:45`，`id` 与时间也不单调：
`10729737755@02:34:22` 与 `10729637897@02:34:49` 反过来）。⇒ 「取第 1 行」是错的，P3 必须**逐页取 max(`created_at`)**。
这条同时是「钉会不会自己变成装饰」的一个真实坑：如果 P3 只 `head -1`，它读到的"最新样"可能是任意一枚旧样。
