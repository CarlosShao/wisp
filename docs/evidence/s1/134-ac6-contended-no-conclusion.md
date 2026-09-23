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

---

## 1. 前半（放宽）：`machine-contended` 改判「本 run 无结论」，步骤仍然无条件

### 1.1 改了什么（逐枚 diff 坐标）

| 文件 | `git diff --numstat` | 改动面 |
|---|---|---|
| `scripts/slo-check.ps1` | `86 16` | 头注释两段 + `if ($reasons.Count -gt 0) { … }` **那一枚块**（含 `exit 1`→`exit 0`） |
| `.github/workflows/ci.yml` | `41 2` | 三个 hunk 全在 `slo-full` 那一段：`@@ -520`（job 头注释一行）、`@@ -525,12`（补 AC#6 说明）、`@@ -570,28`（`Upload SLO report` 那一步的 `if-no-files-found` + 注释） |

删掉的整枚 ci.yml 只有两行，逐字：

```
-  # run is a VALIDITY failure, not a performance regression: re-run it when the
-          if-no-files-found: error
```

⇒ **其余五枚 job、触发表、concurrency 两行（`:50/:51`）一字节未动**（§3 的解析器读数第二次独立复核同一条）。

### 1.2 「取不到样 ⇒ 步骤绿但 job 仍红」这一跳为什么必须动 `Upload` 那一步（**超出简报的地界，报备**）

简报写死「④ `ci.yml` 只许 `slo-full` 那一枚 job 体内的那一步」。我按字面把方案做完后撞墙，墙在这里：

`ci.yml` 的 `Upload SLO report` 那一步是 `path: build/slo/slo-report.json` + `if-no-files-found: error`。
前半硬要求「**不写** `slo-report.json`」⇒ 门步骤 exit 0 之后，那一步照样 `##[error]No files were found …` ⇒
**整枚 job 仍旧 failure，AC#6 变成装饰**。三条出路我逐条核过：

| 出路 | 判定 |
|---|---|
| 争用时也写一枚 `slo-report.json`（里面装"无结论"） | **否**：直接违反简报硬要求，而且让"报告文件存在"这件事不再等于"取到了样"——那正是后半要钉死的病 |
| 给 `Upload` 那一步加 `if:` / `continue-on-error` 让它争用时不跑 | **否**：D22 mode-6 明禁，简报也明写「步骤保持无条件」 |
| `if-no-files-found: error` → `warn`（**1 行，作用域只这一枚 step**） | **采用** |

采用之前把「`warn` 会不会仍上传一枚空 artifact」（会 ⇒ P3 当场致盲）核到源码：
今天真 CI 那枚 job 的日志写着 `Download action repository 'actions/upload-artifact@v4' (SHA:ea165f8d65b6e75b540449e92b4886f43607fa02)`，
按这枚 SHA 取 `src/upload/upload-artifact.ts`：

```
31:  if (searchResult.filesToUpload.length === 0) {
33:    switch (inputs.ifNoFilesFound) {
34:      case NoFileOptions.warn: { core.warning(`No files were found with the provided path: …`) ; break }
40:      case NoFileOptions.error: { core.setFailed(`No files were found …`) ; break }
46:      case NoFileOptions.ignore: { core.info(`No files were found …`) ; break }
53:  } else {                                        <-- 只有这一支才会真上传
```

⇒ `warn` 与 `ignore` 都落在 `if` 里、**不进 else** ⇒ 不创建 artifact；差别只是 `warn` 留一枚 run 页面可见的告警。
选 `warn` 不选 `ignore`：本仓的口径是"不许静默"，一枚告警比一行 info 诚实。
**代价如实报**：这一行把"报告文件没写出来"从**硬失败**降成了**告警**——替它兜底的正是后半的 P3（§2），
所以这两半是一桩交易、不能只交前半，这一点我按简报的定性办，没有单独交前半。

### 1.3 前半的本地读数（三形齐；`争用下取样`已标注）

取证前按 `A103④` 三连查过：`gh run list` 有 2 枚在飞（`35823168549` / `35823141501`，都是编排者 push 的 docs commit）
+ `E:/work/base/actions-runner/_diag/Worker_*.log` mtime 距今 < 60 s + 机器上有别人的 `go.exe`。
⇒ **本节所有读数都是"编队与 CI 同机争用时"取的**，只用它判**形状**（退出码、写了哪些文件、打了哪些行），
**不拿任何一个 D32 数字当结论**。

**形 A — 争用拒采样（新形状：exit 0 + NO CONCLUSION）**。命令逐字：
`pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/slo-check.ps1 -Subset smoke -SecondsPerState 2`，
起跑前 `tasklist //FI "IMAGENAME eq go.exe"` 命中 1 枚（pid 4496，**不是本会话的子进程**，路径 `D:\work\base\go\bin\go.exe`）。

```
RC=0
slo-check.ps1: clearing 4 stale report file(s) from D:\work\workspace\projects plans\Wisp\build\slo
slo-check.ps1: sampling validity precheck (ticket 134 AC#4)
slo-check.ps1: NO CONCLUSION (machine-contended) - subset=smoke refused to sample, no numbers were produced
slo-check.ps1: machine-contended reason: foreign toolchain/wisp process present: go.exe pid=4496 started=2026-09-23 13:46:07 path=D:\work\base\go\bin\go.exe
slo-check.ps1: NO CONCLUSION (machine-contended): 1 reason(s), 0 state file(s) written, slo-report.json NOT written
slo-check.ps1: NO CONCLUSION (machine-contended): this is a SAMPLING VALIDITY verdict, not a
slo-check.ps1: NO CONCLUSION (machine-contended): performance result - D32 (CPU <=0.5%, private
slo-check.ps1: NO CONCLUSION (machine-contended): RSS <=25MB) stays unverified for this run.
slo-check.ps1: NO CONCLUSION (machine-contended): Rerun when the machine is quiet (ticket 134
slo-check.ps1: NO CONCLUSION (machine-contended): AC#4, shape C: refuse loudly instead of
slo-check.ps1: NO CONCLUSION (machine-contended): reporting a contaminated sample as truth).
slo-check.ps1: NO CONCLUSION (machine-contended): AC#6 recolors the refusal, it does not lift it.
slo-check.ps1: NO CONCLUSION (machine-contended): what keeps "no conclusion" from rotting in
slo-check.ps1: NO CONCLUSION (machine-contended): silence is probe P3 of scripts/slo-freshness.sh,
slo-check.ps1: NO CONCLUSION (machine-contended): which ages on the newest uploaded SLO REPORT
slo-check.ps1: NO CONCLUSION (machine-contended): artifact - never on this job, never on this
slo-check.ps1: NO CONCLUSION (machine-contended): exit code. Enough busy days in a row and the
slo-check.ps1: NO CONCLUSION (machine-contended): nail goes red by itself (ticket 134 AC#6).
slo-check.ps1: NO CONCLUSION (machine-contended): record written to D:\work\...\build\slo\slo-no-conclusion.json (NOT a report; slo-report.json is not written on this path)
slo-check.ps1: NO CONCLUSION (machine-contended): exit 0
```

落盘只有 `build/slo/slo-no-conclusion.json` 一枚（`ls -l build/slo` 全文：4 个 `*.json` 被开局清掉，
剩下 `slo-no-conclusion.json` 872 B + 若干 09-20 的旧子目录）；`state-*.json` / `settle.json` / `leak.json` /
`slo-report.json` **零枚**。record 内容（逐字，含 AC#6 要求的两处自陈）：

```json
{
  "generated_at": "2026-09-23T05:46:10Z",
  "subset": "smoke",
  "verdict": "no-conclusion",
  "reason": "machine-contended",
  "machine": "DESKTOP-LVS7839",
  "reason_count": 1,
  "reasons": [
    "foreign toolchain/wisp process present: go.exe pid=4496 started=2026-09-23 13:46:07 path=D:\\work\\base\\go\\bin\\go.exe"
  ],
  "offenders": [ { "name": "go.exe", "pid": 4496, "path": "D:\\work\\base\\go\\bin\\go.exe",
                   "why": "foreign toolchain/wisp process present", "started": "2026-09-23 13:46:07" } ],
  "state_files_written": 0,
  "slo_report_written": false,
  "d32_evaluated": false,
  "meaning": "No number was produced, so no number may be read out of this run. Rerun on a quiet machine.",
  "nail": "scripts/slo-freshness.sh P3 ages on the newest slo-full-report artifact, not on this record"
}
```

⚠ **一个真发现（不是本格的改动，登记）**：`Get-OwnProcessTree` 沿父子链**双向**长，所以它会把我这个 shell 的
**共同祖先下的全部进程**都认成"自己人"。实测：我先 `go build`（pid 31420）再跑门 ⇒ `precheck ok`（该枚 `go.exe` 被判成自家树）；
换成**别人**的 `go.exe`（pid 4496）才拿到 contended。本机手动复跑争用形时，负载必须来自**另一棵进程树**，
否则预检会照放。这个性质 AC#4 就有（本格一字未动），但它会让"我本地造个负载试一下"这条路在 CI 之外的形状上骗人 ⇒ 记在这里。
CI 里对应的面是 pwsh 的祖先链 = `Runner.Worker`/`Runner.Listener` 那棵树，别人（另一枚 run、编队）进不来。

**形 B — 安静放行（证明没被改成"永远无结论"）**。同一枚命令，`go.exe` 那两枚退了之后：

```
slo-check.ps1: precheck ok - no foreign toolchain/runner process, machine-wide cpu max 26%
slo-check.ps1: state Sleeping exit=0 pass=True
slo-check.ps1: state Warm exit=0 pass=True
slo-check.ps1: settle exit=0 pass=True
slo-check.ps1: leak exit=1 flipped_to_fail=True
slo-check.ps1: report written to D:\work\...\build\slo\slo-report.json (all_pass=True)
RC=0
```

**形 C — 出了数而数不过 ⇒ 仍旧 exit 1 判红（AC#6 不许动的那一形，本地抓到一发）**。
同一条命令，负载起在别处时 `wisp slo` 的采样器自己失败了：

```
slo-check.ps1: sampling state Sleeping for 2s
wisp slo: out-of-tree sampling failed: resource: observe: baseline tree read: resource: proc: external system process snapshot: NtQuerySystemInformation: buffer never sufficient (last 1120008 bytes)
slo-check.ps1: state Sleeping exit=2 pass=False
slo-check.ps1: state Warm exit=0 pass=True
slo-check.ps1: settle exit=0 pass=True
slo-check.ps1: leak exit=1 flipped_to_fail=True
slo-check.ps1: report written to D:\work\...\build\slo\slo-report.json (all_pass=False)
RC=1
```

⇒ `all_pass=False` 那条路的 `exit 1` 活得好好的（`slo-check.ps1` 尾部那三行我一个字没动）。
`exit=2` 那枚失败是 registry A14 家族的已知仪器脆性（票 66 记过 `buffer never sufficient`），
**不归本格管**，但它正合用：它证明"红"没有被我在这一格里顺手洗掉。
CI 侧的同形读数另有 §0.4 那枚 `35806505339`（`all_pass=False` + `exit 1`，跑的是旧 blob）。

### 1.4 前半的 mode-6 合规读数（解析器，PyYAML 6.0.3）

```
PARSER-OK; job count = 6
jobs = lint, test-core, test-windows, slo-smoke, slo-full, lint-frontend
all six present = True (missing: [] extra: [] )
runs-on: self-hosted+wisp-slo = 1 | ubuntu-latest = 3 | windows-latest = 2
job slo-full       steps=5  step-level-if-or-continue=NONE  job-level=NONE
  slo-full step 4 | SLO full gate (six states + settle + leak) | extra keys = []
    unconditional run step = True | run = powershell -NoProfile -ExecutionPolicy Bypass -File scripts/slo-check.ps1 -Subset full -SecondsPerState 6
  slo-full step 5 | Upload SLO report | extra keys = []
    with = {'name': 'slo-full-report', 'path': 'build/slo/slo-report.json', 'if-no-files-found': 'warn'}
on: = {'pull_request': None, 'push': ['branches'], 'schedule': [{'cron': '37 19 * * *'}], 'workflow_dispatch': None}
push.branches = ['main', 'dev'] | schedule = [{'cron': '37 19 * * *'}] | pull_request present = True | workflow_dispatch present = True
50:  group: ci-${{ github.workflow }}-${{ github.ref }}${{ github.event_name == 'push' && format('-{0}', github.sha) || '' }}
51:  cancel-in-progress: ${{ github.event_name == 'pull_request' }}
```

要害一行：`if-no-files-found` 是**动作的输入**（在 `with:` 里），不是 workflow 的条件键；
step 级 `if` / `continue-on-error` 在 `slo-full` 五枚 step 上仍是 `NONE`，`extra keys = []`。
全文件 `grep -cE '^[[:space:]]+continue-on-error:'` = **0**；`^[[:space:]]+if:` = 9 行（`:178 :215 :291 :351 :389 :399 :421 :457 :473`），
**逐行都在 `slo-full` 之外**（`slo-full:` 的 job 键现在在 `:537`），与钉的 P1 探针同源。
