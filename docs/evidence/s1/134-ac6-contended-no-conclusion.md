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

---

## 2. 后半（收紧）：新鲜度钉改按「最近一次产出有效样本的记录」计龄

### 2.1 判据物到底是什么、从哪儿读（正面回答交回项 ④）

**判据物** = GitHub Actions 里**名字恰为 `slo-full-report` 的最新一枚未过期 workflow artifact** 的 `created_at`。
读法 = `GET /repos/{owner}/{repo}/actions/artifacts?per_page=100&page=N`（`gh api`，N = 1..`SLO_FULL_ARTIFACT_PAGES`，默认 3），
jq 里三条件：`.name=="slo-full-report"` + `.expired==false` + 取 `(.created_at, .workflow_run.id, .id)`，**逐页取 max(`created_at`）**。

它为什么等价于"取到了有效样本"（三段都在同一枚 commit 里可核）：
1. `ci.yml` 的 `Upload SLO report` 那一步 `path:` **只写死一枚文件** `build/slo/slo-report.json`；
2. `scripts/slo-check.ps1` 只在**真跑完六态 + settle + leak** 之后才 `Set-Content` 那枚文件；
   争用路径开局就把上一轮遗留的 `*.json` 清掉，然后打 `0 state file(s) written, slo-report.json NOT written` 并 exit 0；
3. `actions/upload-artifact@v4`（= runner 今天真下载的 `ea165f8d…`，见 §1.2）在 `filesToUpload.length === 0` 时
   **根本不进上传那支** ⇒ 不创建 artifact。

⇒ 「artifact 在」= 「report 在」= 「数字在」。**不是 job 的 conclusion**（这一点 §2.3 的 T1 直接用一发读数钉死：
job 那侧全绿、钉照样红）。

### 2.2 `scripts/slo-freshness.sh` 改了什么

`git diff --numstat` = `198 28`（对照锚点 blob `5a4a339a…`）。三处：

- **P3 新增**（本节全部读数都打在它上面）。两个新失败令牌：`slo-full-sample-stale` / `slo-full-sample-never`。
- **P2 保留不动**（判据、阈值、12 小时 queued 那支一字未改）。留着的理由写进了文件头：
  P2 仍是**唯一**看得见"job 卡在 queued 永不启动 = 另一种静默死"的探针；P3 看不见那形（没 artifact 也没 report，
  但也没人在跑）。两枚探针各查一种死法，阈值分叉独立（`SLO_FULL_MAX_AGE_DAYS` / `SLO_FULL_SAMPLE_MAX_AGE_DAYS`，默认都 3 天）。
- **顺带纠一处腐坏的引用**：P1 注释原文写「measured: ci.yml carries **6** such lines outside slo-full」，
  2026-09-23 复测是 **9**（8 枚 `if: ${{ !cancelled() }}` + 1 枚 `if: always()`，逐行行号已写进注释）。
  只改注释，判据一字未动。
- **一处防漂移的小重构**（被自己的变异逼出来的，见 M6）：artifact 名字在脚本里原本要写两遍（jq 里一遍、消息里一遍），
  改成单一常量 `sample_artifact=slo-full-report`（**普通赋值，不是 env 接缝**）——否则"改了过滤器忘了改文案"
  这种漂移可以造出一枚自称读了 D32 样本、实际读的是别家 artifact 的钉。

新接缝（都只塑造**输入**，没有一个能关掉探针）：`SLO_FULL_LAST_SAMPLE` 注入一枚有效样本记录；
取值 `none` = 造"扫描回来一份都没有"的世界（只会更红，不会更绿，M5b 专测这条）；
`SLO_FULL_ARTIFACT_PAGES` 只改扫描页数；`can_look()` 把"查不了 ⇒ exit 2"从 P2 一处扩到两枚探针共用。

### 2.3 硬判据：红 → 绿 → 红三态（**逐字**，final script 上跑的）

世界设定（这就是 AC#6 说的那局面）：**job 侧每枚都新鲜**（注入一枚 1 小时前 `completed/success` 的 job 记录，
P2 判绿），**盘上最近一枚真 report 是 10 天前**（这一条不打诳语：P3 读的是**真 API**，
只把 `SLO_FRESH_NOW` 拨到 `2026-10-03T04:44:08Z`，也就是"最新那枚 report 已经 10 天"）。

```
### T1 RED (final script)
slo-freshness: P2 newest slo-full job record: created_at=2026-10-03T04:00:00Z status=completed conclusion=success
slo-freshness:   job url: https://github.com/CarlosShao/wisp/actions/runs/0/job/0
slo-freshness:   age: 0 day(s) (2648 s); threshold: 3 day(s)
slo-freshness: P3 newest VALID SAMPLE: artifact slo-full-report created_at=2026-09-23T04:44:08Z artifact_id=10733125568 run=35819355656
slo-freshness:   scanned 3 page(s) x 100 of the artifact listing; 70 valid-sample record(s) considered
slo-freshness:   age: 10 day(s) (864000 s); threshold: 3 day(s)
slo-freshness: FAIL slo-full-sample-stale: the last VALID SLO SAMPLE is 10 day(s) old (> 3). Jobs may well be green - since ticket 134 AC#6 a machine-contended run exits 0 without a report - so a green job is not the question. The question is when numbers last existed, and that was 2026-09-23T04:44:08Z (run 35819355656). D32 has been unverified since then.
slo-freshness: 1 probe failure(s) - ticket 134 AC#3/AC#6 nail is red
RC=1
### T2 GREEN (SLO_FULL_SAMPLE_MAX_AGE_DAYS=30)
（同上三行 P2 读数不变）
slo-freshness:   age: 10 day(s) (864000 s); threshold: 30 day(s)
slo-freshness: OK - slo-full still has automatic triggers, a trigger record inside the window, and a VALID SAMPLE inside the window
RC=0
### T3 RED again (env removed)
slo-freshness:   age: 10 day(s) (864000 s); threshold: 3 day(s)
slo-freshness: FAIL slo-full-sample-stale: ... （与 T1 同一枚令牌、同一句原文）
RC=1
```

⇒ **红→绿→红齐**，且绿那一发**只由阈值放宽造成**（同世界同数据，只改 `SLO_FULL_SAMPLE_MAX_AGE_DAYS`）——
证明红是判据咬的、不是常数红。**三发里 P2 那三行都是绿的** ⇒ 直接可视化"job 年龄这枚钉在 AC#6 之后就是装饰"。

### 2.4 其余读数（同一次会话连跑，逐字要点）

**L0 真跑、零接缝、真 API（今天的真实状态）**：

```
slo-freshness: P2 newest slo-full job record: created_at=2026-09-23T05:37:21Z status=completed conclusion=failure
slo-freshness:   job url: https://github.com/CarlosShao/wisp/actions/runs/35823168549/job/107059176883
slo-freshness: P3 newest VALID SAMPLE: artifact slo-full-report created_at=2026-09-23T04:44:08Z artifact_id=10733125568 run=35819355656
slo-freshness:   scanned 3 page(s) x 100 of the artifact listing; 70 valid-sample record(s) considered
slo-freshness: OK - ...
RC=0
```

⇒ 两件事一次看清：`conclusion=failure` 的 job（那一枚就是本代理 §0 那枚 commit 被 push 之后的 run，争用红的）
与"最新有效样 04:44:08Z"是**两个问题**；且 3 页 300 条 artifact 里 `slo-full-report` 命中 **70 枚**（去重后参与取 max）。

**M4（判据独立性）**：同一 +10d 世界，只放宽 **P2** 的阈值（`SLO_FULL_MAX_AGE_DAYS=30`）⇒ **仍 RC=1**，
令牌仍是 `slo-full-sample-stale`。⇒ P3 不是 P2 的别名，两枚旋钮各管各的。

**M5 / M5b（一份 report 都没有的世界）**：`SLO_FULL_LAST_SAMPLE=none` + job 新鲜 ⇒

```
slo-freshness: FAIL slo-full-sample-never: no uploaded slo-full-report artifact at all (SLO_FULL_LAST_SAMPLE=none (fixture: scan returned no slo-full-report artifact)). slo-full is being triggered but has produced no valid sample that this pin can see - which since AC#6 is NOT the same question as 'is the job green'. D32 (Sleeping CPU <=0.5%, private RSS <=25MB) is unverified.
RC=1
```

M5b：同世界把 `SLO_FULL_SAMPLE_MAX_AGE_DAYS` 放宽到 30 ⇒ **仍 RC=1**（`none` 不是关探针的口令）。

**M6 / M6b（过滤器是有牙的还是装饰？）**：把 /tmp 副本里的 `sample_artifact` 换成 `slo-smoke-report`，
其余一字不动，跑真 API：

```
M6  : P3 newest VALID SAMPLE: artifact slo-smoke-report  created_at=2026-09-23T05:39:19Z artifact_id=10734395325 run=35823168549 ; 100 record(s) considered
M6b : P3 newest VALID SAMPLE: artifact slo-full-report   created_at=2026-09-23T04:44:08Z artifact_id=10733125568 run=35819355656 ;  70 record(s) considered
```

⇒ 换掉一枚字符串，读到的"最新样"就换了：**过滤器真在起作用**，不是常数。更要命的是 M6 那枚读到的
`run=35823168549` 正是**本票据 commit 被 push 的那枚 run——它的 `slo-full` job 是争用红的、什么都没采样**，
而 hosted 的 `slo-smoke` 上传了报告。⇒ 如果 P3 的名字写错成 smoke，它会在全机零 D32 样本的日子里报"样很新鲜"。
这条同时解释 §2.2 里那处"单一常量"重构为什么不是洁癖。真文件复核：`grep -c slo-smoke-report scripts/slo-freshness.sh` = **0**，
变异只活在 `/tmp/134ac6/mut/`。

**M7（没有"查不了=通过"）**：`env -u GH_TOKEN -u GITHUB_TOKEN`、两个注入接缝都空 ⇒

```
slo-freshness: no GH_TOKEN/GITHUB_TOKEN - the freshness probes cannot look (this is not a pass)
RC=2
```

**M8 / M9（P1 在本格改动之后仍有牙）**：`/tmp` 里的 ci.yml 副本给 `slo-full:` 塞一行
`if: ${{ !cancelled() }}` ⇒ `RC=1 slo-full-trigger-missing: the slo-full job now carries 'if:'`；
再删掉 `- cron:` 那行 ⇒ 两条同时红（`no 'cron:' schedule entry` + `carries 'if:'`）。
⇒ 我把 `if-no-files-found: error` 改成 `warn` 之后，P1 那枚 `^[[:space:]]+if:` 锚**没有被同名前缀糊住**
（`if-no-files-found:` 不匹配 `^ *if:`，实测 M8 仍能红），真文件跑 P1 也是绿的（L0）。
还原证明：`grep -c "if: \${{ !cancelled() }}" .github/workflows/ci.yml` = 11，但其中**只有 8 枚**是步级条件行（把行首锚与行尾锚都加上的 `grep -cE` 复算 = 8，
另 3 枚是注释里出现同一串字符的行），再算上 `:291` 那枚 `if: always()` 才是 P1 说的 9 枚真条件行；
`git diff --numstat .github/workflows/ci.yml` 为空（本程改动已入库），变异文件全在 `/tmp`。

### 2.5 P3 的边界（写在这里，不留到下次才发现）

1. **出了数而数不过 ⇒ 不算"有效样本记录"**：`all_pass=false` 那枚 run `exit 1`，默认 `success()` 语义下
   `Upload` 那一步根本不跑 ⇒ 无 artifact ⇒ P3 的钟不动。**故意如此**（那种 run 自己就是红的，
   把它计入"样很新鲜"会让一枚长期红的门看起来健康）。今天有一枚实例：§0.4 的 `35806505339`。
2. **扫描上限 `SLO_FULL_ARTIFACT_PAGES * 100`**（今天 300 条 ≥ 全仓 168 条）。越界时是**部分扫描**，
   部分 max 只会**偏旧** ⇒ 误差方向是**偏红**，不会偏绿。
3. **能伪造绿的路径只剩一条，且不在钉的地界内**：谁要能让 P3 读到假"有效样本"，就得**上传一枚名字恰为
   `slo-full-report` 的 artifact**——那需要本仓一枚通过鉴权的 Actions run（编辑 ci.yml 的同一条信任边界），
   而那条边界由 P1 + `tools/d22scan` + review 看着。**删 artifact 不能骗绿**（只会更红）。P3 不下载 zip、
   不看内容：它只回答"最近一次有 report 存在是什么时候"。
4. **保留期**：artifact 默认留 90 天（P3 今天读到的那枚 `slo-full-report` id=10733125568 的
   `expires_at=2026-12-22T04:41:33Z`，`expired=false`）。3 天窗离 90 天很远；顺带一条实测：列表端点
   **不给** `file_count`/`state`（读回来是 null），所以 P3 只能用 `name` + `expired` + `created_at` 三件，
   想"顺手验一下 artifact 里有几个文件"这条路在这个端点上不存在（要验内容就得下 zip，另说）。
   `.expired==false` 那一条是防"拿一枚已过期的当有效样"，不是防 retention 到期。
5. **PR run 的 artifact 也计入**：`on.pull_request` 无分支过滤，PR 上跑出来的 `slo-full-report` 名字相同 ⇒
   会被算作"出过样"。这不是漏洞而是**放宽**（真出过样），且 P2 那侧仍只看 main/dev；本仓目前无 PR 通路在用。
6. **P3 不看 job 在哪个 run 里**：它信 artifact 的 `created_at`。若 runner 时钟漂了，
   漂向未来的 report 会显得更新 —— 但这枚 report 的时间戳是 **GitHub 服务端**写入的（artifact `created_at`），
   不是 runner 报的，所以这条不适用；真正的外部时钟是钉自己的 `now`（`date -u +%s`，ubuntu runner）。

---

## 3. 门禁（照票 134 AC#5 那一套重跑，被验版本 = `44ab500` 的树）

### 3.1 `sh scripts/d22scan.sh` 纯净快照 rc=0 + 台账八 scope 不降

```
$ rm -rf /tmp/d22-134ac6 && mkdir -p /tmp/d22-134ac6
$ git archive HEAD | tar -x -C /tmp/d22-134ac6            # HEAD = 44ab500
$ cd /tmp/d22-134ac6 && sh scripts/d22scan.sh
runtests.sh: OK - packages=[./...] top-level: PASS=21 FAIL=0 SKIP=0, === RUN=31, '[no tests to run]'=0
--- PASS: TestBuiltBinaryGoesRedEndToEnd (1.27s)   （6 条子用例全 PASS）
d22scan.sh: scan of /tmp/d22-134ac6
d22scan: examined 225 production Go files under internal/ and cmd/ of C:/Users/swq/AppData/Local/Temp/d22-134ac6
d22scan: clean - no D22 ban violations
D22-RC=0
```

| scope | 票 99 基线 | AC#5 那格（`b9b2072`） | 本次（`44ab500`） | 判定 |
|---|---|---|---|---|
| bans #1-5 `internal/` | 197 | 202 | **203** | 不降 |
| bans #1-5 `cmd/` | 20 | 22 | **22** | 持平 |
| ban #6 `frontend/` | 37 | 40 | **40** | 持平 |
| ban #7 `internal/tools/` | 17 | 18 | **18** | 持平 |
| ban #8 `design/` | 16 | 16 | **16** | 持平 |
| ban #8 `frontend/` | 37 | 40 | **40** | 持平 |
| ban #8 `internal/` | 342 | 387 | **390** | 不降 |
| ban #8 `cmd/` | 26 | 36 | **37** | 不降 |

**8 行 scope 齐全、逐格不降**。三处上涨（203 / 390 / 37）不是本票造成的：**本格零枚 `.go` 改动**
（`git diff --name-only b723978..HEAD -- '*.go'` 在本票四枚路径里为空），
涨的是同树并行的票 129/130/131 那几枚 commit。

### 3.2 两台 YAML 解析器读数（最终树，PyYAML 6.0.3）

`ci.yml`（§1.4 那张表在改动入库之后又跑了一遍，逐字相同：`job count = 6` /
`all six present = True (missing: [] extra: [])` / `runs-on ubuntu-latest=3 windows-latest=2 self-hosted+wisp-slo=1` /
`slo-full steps=5 step-level-if-or-continue=NONE job-level=NONE` / 门步骤 `unconditional run step = True` /
`with = {'name': 'slo-full-report', 'path': 'build/slo/slo-report.json', 'if-no-files-found': 'warn'}`）

`slo-fresh.yml`（独立性复测）：

```
slo-fresh.yml PARSER-OK; jobs = ['slo-fresh']
name = slo-full-must-keep-getting-triggered | runs-on = ubuntu-latest
job-level if/continue-on-error = False False
  step 1 | actions/checkout@v4 | conditional keys = []
  step 2 | slo-full freshness pin (ticket 134 AC#3) | conditional keys = []
  step 3 | Shell lint for the pin | conditional keys = []
permissions = {'contents': 'read', 'actions': 'read'}
on = {'schedule': [{'cron': '23 */6 * * *'}], 'workflow_dispatch': None}
```

⇒ 钉仍是**自带时钟 + 手动入口**、`ubuntu-latest`、`permissions` 没有加宽（`actions: read` 本来就够读 artifact 列表）、
无 `if:` / 无 `continue-on-error`、无路径过滤；它与被测的 `slo-full` 不共用任何触发器。
`slo-fresh.yml` 本程 `git diff --numstat` = `15 0`（纯注释，零删除）。

### 3.3 D32 两条阈值未动（正面回答交回项 ⑥）

- `internal/observe/thresholds.go` **锚点与 HEAD 同一枚 blob**：
  `git rev-parse b723978:internal/observe/thresholds.go HEAD:internal/observe/thresholds.go`
  ⇒ 两侧都是 `e2677b11b5a5adff8b4eb87c36d31286a274471d`；`git diff --numstat b723978..HEAD -- 该文件` 为空。
  文件里那两行判据原文（`grep -n` 于 HEAD）：
  ```
  19:	memCapSleeping     int64 = 25 << 20
  26:	cpuLimitSleeping     = 0.5 // % of all-core mean, 1min window
  ```
- `scripts/slo-check.ps1` 自锚点以来**被删掉的行**逐字全列（16 行，一类都不许漏）：
  4 行头文档 + 2 行头文档 + 6 行 `Write-Host` 文案 + 1 行 step-summary 文案 + 1 行 step-summary 表格行 +
  1 行 `Add-Content` 收尾 + **1 行 `exit 1`** ⇒ 没有任何一行是判据条件。HEAD 上仍在的原样条件：
  ```
  153:$loadNames = @('go.exe', 'gofmt.exe', 'cgo.exe', 'compile.exe', 'asm.exe', 'link.exe',
  211:$cpuBusyPct = 50
  225:if ($cpuMax -ge $cpuBusyPct) {
  381:$allPass = ($failingStates.Count -eq 0) -and $settlePass
  396:if (-not $allPass) { exit 1 }
  ```
  ⇒ "什么时候拒绝出数"（AC#4 的方向）与"出了数而不过就红"两条都在原位。
  ps1 里出现的 `0.5%` / `25MB` 六处全是**文档与日志文案**（`:29 :135 :241-242` 注释、`:256-257` 打印），
  数值本身由 `wisp slo` 判，脚本不参与。
- `docs/PLAN.md`、`docs/specs/**`、`internal/risk/**`、`rules_gateway.go`、`tools/d22scan/**`、
  `allowlist.txt`、`scripts/d22scan.sh`、任何 golden：本程 `git status --porcelain` 里从未出现，
  `git log --name-only b723978..HEAD`（我的三枚 commit）只含 §"commit 清单"里那六枚路径。

### 3.4 脚本 lint：shellcheck 这一格**这次真的跑起来了**（AC#5 欠的那格补上）

本机仍然没有 `shellcheck`（`command -v shellcheck` rc=1），但 **docker 可用**，所以按简报给的挂载口径跑
Linux 版（比 Windows 原生更贴 ubuntu runner）：

```
$ docker info --format '{{.ServerVersion}} {{.OSType}}'   ->  29.6.2 linux
$ MSYS_NO_PATHCONV=1 docker run --rm --entrypoint shellcheck -w /src \
      -v /d/work/tmp/134ac6-sc:/src:ro koalaman/shellcheck:stable --version
    ShellCheck - shell script analysis tool / version: 0.11.0
```

**假绿防护**（简报点过的那发：`docker run -v "C:\…"` 会静默挂空且 rc=0）：镜像里没有 `/bin/sh`（拿 `ls -l /src/go.mod`
进不去），所以我把"挂载真生效"换成**更强的一发**——往挂载目录里埋一枚**故意有病**的文件，让容器把**我的文件名与行内容**原样报回来：

```
In /src/mountproof.sh line 4:
  echo $1
       ^-- SC2086 (info): Double quote to prevent globbing and word splitting.
```

⇒ 内容穿过来了，挂载不是空的。

**第一发读数（要害）**：`shellcheck -s sh scripts/slo-freshness.sh` 在**我接手时的那枚文件**上 **rc=1**：

```
In /src/scripts/slo-freshness.sh line 85:
here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
              ^-- SC1007 (warning): Remove space after = if trying to assign a value ...
In /src/scripts/slo-freshness.sh line 86:
root=$(CDPATH= cd -- "$here/.." && pwd)
              ^-- SC1007 (warning): ...
```

⇒ 这两行是 **AC#3 落地的原码**（不是我写的）。也就是说 `.github/workflows/slo-fresh.yml` 的
`Shell lint for the pin` 那一步（AC#5 亲手指定为"shellcheck 那一格的兜底，缺它即硬红"）**在真装上 shellcheck 的
runner 上会直接红**——AC#5 那格当时以"本机没装"结案，欠的就是这一发。

**修法**（`scripts/slo-freshness.sh`，在本格解冻清单内）：把空的前缀赋值写成 `CDPATH=''`，
与 `CDPATH=` 语义完全相同（先在容器里拿三形 `CDPATH=''` / `CDPATH= ;` / 裸 `cd` 各跑一发，`shellcheck` 对三形都 rc=0，
我选**语义等价最强、改动最小**的第一形），并在原位写下为什么。改完：

```
$ MSYS_NO_PATHCONV=1 docker run --rm --entrypoint shellcheck -w /src -v /d/work/tmp/134ac6-sc:/src:ro \
      koalaman/shellcheck:stable -s sh /src/scripts/slo-freshness.sh
SC-RC=0
```

⚠ **过程中踩到自己造的一枚新坑，记下来防再犯**：我为此写的那段注释里有一行以
`# shellcheck 0.11.0 reads the bare ...` 开头 ⇒ **注释自己变成了一枚无法解析的 shellcheck 指令**，
`SC1073/SC1072 (error)`、rc=1。已改写措辞（`# the linter (shellcheck 0.11.0) ...`），
并全库扫过一遍：`grep -n "^\s*# shellcheck" scripts/*.sh` 除既有 directive 外 0 命中。

**其余 lint**：

| 仪器 | 对象 | 读数 |
|---|---|---|
| `sh -n` | `scripts/slo-freshness.sh` | rc=0 |
| `bash -n` | `scripts/slo-freshness.sh` | rc=0 |
| shellcheck 0.11.0 `-s sh` | `scripts/slo-freshness.sh`（HEAD 最终版） | **rc=0** |
| PSParser | `scripts/slo-check.ps1` | `PARSE-OK tokens=1842 errors=0` |
| shellcheck（信息性，非本格地界） | `d22scan.sh` 2 / `portable-tests.sh` 39 / `winsec-tests.sh` 11 / `wisp-cli-tests.sh` 5 | 共 57 处发现 |

⇒ 那 57 处**不归本票**（`slo-fresh.yml` 只 lint `slo-freshness.sh`；四枚脚本里三枚在禁改清单上），
我只把数字如实贴出来，另立不立票由编排者定。

### 3.5 测了但结论只能挂在"未验证"上的格子

见 §5 末列表；其中 shellcheck 这一格**已从"未验证"升成"已验证"**（本机 docker 内跑 Linux 版 + 上表 rc=0），
剩下的是 `slo-fresh.yml` 自身在 CI 上的第一次执行、`schedule` 那半边、以及 `GITHUB_STEP_SUMMARY` 分支。
