# CI runner readings — 2026-09-21（只出读数，不下结论、不勾框）

- 采集代理：`audit-runner-readings`（只读调查代理；不改生产码、不改票面、不 commit、不 push）
- 采集时刻：2026-09-21 20:3x CST（= 12:3x UTC）
- 采集时本地 HEAD：`e3fb933`（`dev`；origin/dev == cnb/dev == `34a810b`，即 `dev` 本地领先 1 枚 `e3fb933`，lag 非 0，如实报）
- 工作树状态（采集开始时 `git status --short`）：只有两枚未跟踪文件
  `?? docs/evidence/s1/104-adversarial-acceptance.md` 与 `?? docs/evidence/s1/109-adversarial-acceptance.md`
  ——票 113/111/92b 的在飞改动**此时都已进 HEAD**，树对我而言是干净的；**本文件是我唯一新增的仓内文件**。
- 四格：票 70 AC#6 / 票 72 AC#4 / 票 77 AC#4 / 票 70 AC#2。

## 格 1 · 票 70 AC#6 —— 已完成 run 的逐 job 结论表

**判据（票面 AC#6 逐字）**：「`gh run view <新 run> --json jobs` 的**逐 job 结论**贴进本票 Progress log。判据是**五个 job 全 pass**」。
（⚠ 引用腐坏点名：判据写"五个 job"，但票 77 已把 `lint-frontend` 加进流水线——本次 run 实际是 **6 个 job**。按"全 pass"精神记。）

**取数命令（逐字，20:3x 本机执行）**：

```
gh run view 35600043583 --json headSha,conclusion,url,event,displayTitle,jobs > /tmp/run35600043583.json
```

形状判据：`GH_RC=0`、首字符 `{`、`jobs` 被 python `json.load` 解析成 6 元素列表、每个元素都含 `conclusion` 与 `steps` 字段（当场断言通过）——不是"没看到 error 字样"式判成。
选 run `35600043583` 的理由：completed、conclusion 已产出、是**晚于** `35599458439` 的最新已完成 run（`35599458439` 归 `agent-ticket112b`，本代理未触碰其读数）。
headSha `d13e597b7bb728bdf2afc4b700559bc10ad42350`（= 本地 HEAD `e3fb933` 的曾祖父，祖先关系 `git merge-base --is-ancestor d13e597 HEAD` 为真）、event=push、run 整体 `conclusion=failure`。

**逐 job 逐步骤结论（gh 原文解析，全表如下）**：

```
JOB lint           => failure (jobId=106334066361)
   step  1 [success  ] Set up job
   step  2 [success  ] Run actions/checkout@v4
   step  3 [success  ] Run actions/setup-go@v5
   step  4 [success  ] D22 scanner positive control (tools/d22scan tests, seeded red)
   step  5 [success  ] D22 seven-ban + emoji scan (tools/d22scan)
   step  6 [success  ] gofmt (gofumpt)
   step  7 [success  ] go vet (module)
   step  8 [success  ] go vet (tools/d22scan module)
   step  9 [failure  ] staticcheck
   step 10 [skipped  ] mockllm module vet
   step 19 [skipped  ] Post Run actions/setup-go@v5
   step 20 [success  ] Post Run actions/checkout@v4
   step 21 [success  ] Complete job
JOB lint-frontend  => success (jobId=106334066228)
   step  1 [success  ] Set up job
   step  2 [success  ] Run actions/checkout@v4
   step  3 [success  ] Run actions/setup-go@v4(→setup-node)  ※原文为 "Run actions/setup-node@v4"
   step  4 [success  ] npm ci (lockfile is the only source of deps)
   step  5 [success  ] typecheck (tsc -b, noEmit per tsconfig.*.json)
   step  6 [success  ] lint (oxlint)
   step  7 [success  ] token drift guard (generated theme must equal the C21 table)
   step  8 [success  ] build (vite build -> frontend/dist, the bytes go:embed carries)
   step  9 [success  ] L2 card renders the real risk fields (AC#3 render evidence)
   step 17 [success  ] Post Run actions/setup-node@v4
   step 18 [success  ] Post Run actions/checkout@v4
   step 19 [success  ] Complete job
JOB slo-full       => success (jobId=106334066098)
   step  1 [success  ] Set up job
   step  2 [success  ] Run actions/checkout@v4
   step  3 [success  ] Run actions/setup-go@v5
   step  4 [success  ] Build wisp.exe (deps cached on the runner)
   step  5 [success  ] SLO full gate (six states + settle + leak)
   step  6 [success  ] Upload SLO report
   step 11 [success  ] Post Run actions/setup-go@v5
   step 12 [success  ] Post Run actions/checkout@v4
   step 13 [success  ] Complete job
JOB slo-smoke      => success (jobId=106334066499)
   step  1 [success  ] Set up job
   step  2 [success  ] Run actions/checkout@v4
   step  3 [success  ] Run actions/setup-go@v5
   step  4 [success  ] Cache third_party (deps.toml key)
   step  5 [success  ] Build wisp.exe
   step  6 [success  ] SLO smoke gate
   step  7 [success  ] Upload SLO report
   step 12 [success  ] Post Cache third_party (deps.toml key)
   step 13 [success  ] Post Run actions/setup-go@v5
   step 14 [success  ] Post Run actions/checkout@v4
   step 15 [success  ] Complete job
JOB test-core      => failure (jobId=106334066356)
   step  1 [success  ] Set up job
   step  2 [success  ] Run actions/checkout@v4
   step  3 [success  ] Run actions/setup-go@v5
   step  4 [success  ] Start compose test services (mock-llm on 18080)
   step  5 [success  ] Probe mock-llm
   step  6 [success  ] Environment fork assertion (WISP_ENV=test data dir)
   step  7 [failure  ] Portable package tests (agent/llm/config/memory/observe/secret/risk/statemachine/...)
   step  8 [success  ] Stop compose services
   step 15 [skipped  ] Post Run actions/setup-go@v5
   step 16 [success  ] Post Run actions/checkout@v4
   step 17 [success  ] Complete job
JOB test-windows   => failure (jobId=106334066496)
   step  1 [success  ] Set up job
   step  2 [success  ] Run actions/checkout@v4
   step  3 [success  ] Run actions/setup-go@v5
   step  4 [failure  ] Windows ACL sealing gate (internal/winsec's own tests, ticket 110)
   step  5 [skipped  ] Cache third_party (deps.toml key)
   step  6 [skipped  ] cgo build smoke (build.ps1 fetch-deps + mingw link + doctor)
   step  7 [skipped  ] Portable windows tests (proc/secret/config/risk)
   step  8 [skipped  ] PathResolver junction placeholder (real cases tickets 18/20)
   step 15 [skipped  ] Post Run actions/setup-go@v5
   step 16 [success  ] Post Run actions/checkout@v4
   step 17 [success  ] Complete job
```

注：上表 step 3/`setup-node` 一行的 `※` 是我加的说明，其余名称与结论逐字来自 gh JSON（原始落盘 `/tmp/readings-run-35600043583.txt`，本文件采集时该 JSON 在 `/tmp/run35600043583.json`）。

**读数即事实**：3 success（`lint-frontend`/`slo-full`/`slo-smoke`）+ 3 failure（`lint`/`test-core`/`test-windows`）。
顺带三条**步级**新事实（只报读数不判票）：`go vet (module)` 首次在这条链上 success（票 78 那两颗哑弹的 CI 侧首次绿证，就本 run 而言）；`lint` 现在红的正是票 70-d 预警过的 `staticcheck`；
`test-windows` 的 `PathResolver …` 步骤这次是 **skipped**（被票 110 新加的 step 4 `Windows ACL sealing gate` failure 挡在前面）——⇒ 本格**不能**充当票 72 AC#4 的 runner 侧读数，见格 2。

**标签**：〔日志＋归档，我抽验〕——JSON 由 gh 从 GitHub API 现取现解析，run url 在
`https://github.com/CarlosShao/wisp/actions/runs/35600043583`，编排者可 `gh run view 35600043583 --json jobs` 复跑抽验。

**本格结论**：**读数已采到；格不能闭**——AC#6 的判据是"全 pass"，本 run 6 个 job 里 3 个 failure。


## 格 2 · 票 72 AC#4 —— 本机 + runner 两侧逐字同命令

## 格 3 · 票 77 AC#4 —— ban #6 / ban #8 对真实 frontend/ 的逐作用域文件数

## 格 4 · 票 70 AC#2 —— 在 Linux 上复现那 4 条红

## 我没能采到的
