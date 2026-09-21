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

**判据（票面 AC#4 逐字）**：「贴出两侧**逐字同命令**的输出与真实 exit code；runner 侧必须来自一次**真 run**，不是本地等价环境。」

**① 那条命令逐字是什么（当场核，不引旧账）**：
票 72 的判据主体是 run 35547905707 里红的那条 `test-windows` 步骤。取票面修复（`f1033e1`，2026-09-21T10:20:50+08:00，
`git merge-base --is-ancestor f1033e1 HEAD` 真）**之后、票 110 新步骤进 ci.yml 之前**的最后一次真 run =
**run `35594826728`**（headSha `b878b30f53ffe634fa299a453c08c2fb52860e4c`，2026-09-21T11:35:03Z，windows-latest）。
该 commit 的 ci.yml 第 206 行（`git show b878b30:.github/workflows/ci.yml` 逐字）：

```
bash tools/d22scan/runtests.sh ./internal/risk/ -run TestPathResolverJunctionWindows
```

步骤名 `PathResolver junction placeholder (real cases tickets 18/20)`，步骤 env `WISP_ENV: test`，
shell `C:\Program Files\Git\bin\bash.EXE --noprofile --norc -e -o pipefail`。
引用核对：票面 13 行说的 `pathresolver_junction_windows_test.go:88` **仍对得上**——HEAD 上
`func TestPathResolverJunctionWindows(t *testing.T)` 就在第 88 行（当场 grep）。
形态说明（防"拿形态不同的读数凑"）：`runtests.sh` 是包装器，它**强制** `go test -v -count=1` 并附三条致命判据
（SKIP≠过、0 PASS 0 FAIL≠过、go test rc 原样传播）；两侧跑的是**同一个包装器同一条命令行**，不是"runner 整包 vs 本机单包"那种错配。
票 72 本机半边原先自报的是 `go test ./internal/risk/... -count=2`（**整包**形态）——那与 runner 步骤**形态不同**，本格**不拿它充数**；下面两侧都按逐字同命令重采。

**② 本机侧（逐字同命令 + 同 env，Git Bash on Windows）**：

命令（逐字）：

```
WISP_ENV=test bash tools/d22scan/runtests.sh ./internal/risk/ -run TestPathResolverJunctionWindows
```

- (a) 在**当前工作树**（HEAD `e3e60b0` 起测时，采集中共享树被邻居推进到 `72bc745`，两者都含 `f1033e1`）：`RC=0`，输出逐字：

```
2026-09-21 20:45:21 INFO winsec: sealing path resolver installed resolver=risk.c26Pipeline probes_passed=2
=== RUN   TestPathResolverJunctionWindows
--- PASS: TestPathResolverJunctionWindows (0.86s)
PASS
ok  	github.com/CarlosShao/wisp/internal/risk	1.001s
runtests.sh: OK - packages=[./internal/risk/ -run TestPathResolverJunctionWindows] top-level: PASS=1 FAIL=0 SKIP=0, === RUN=1, '[no tests to run]'=0
```

- (b) 在 **`b878b30` 的纯净快照**里（`git archive b878b30 | tar -x -C /tmp/wisp-snap-b878b30`，
  与 runner 那次**同一 commit 同一命令行**）：`RC=0`，输出逐字：

```
2026-09-21 20:45:44 INFO winsec: sealing path resolver installed resolver=risk.c26Pipeline probes_passed=2
=== RUN   TestPathResolverJunctionWindows
--- PASS: TestPathResolverJunctionWindows (0.17s)
PASS
ok  	github.com/CarlosShao/wisp/internal/risk	0.216s
runtests.sh: OK - packages=[./internal/risk/ -run TestPathResolverJunctionWindows] top-level: PASS=1 FAIL=0 SKIP=0, === RUN=1, '[no tests to run]'=0
```

**③ runner 侧（真 run 日志逐字段）**：

命令（逐字，日志行来自 `gh run view --job 106317137389 --log`，job=test-windows of run `35594826728`）：

```
2026-09-21T11:37:12.7954731Z ##[group]Run bash tools/d22scan/runtests.sh ./internal/risk/ -run TestPathResolverJunctionWindows
2026-09-21T11:37:12.7980179Z shell: C:\Program Files\Git\bin\bash.EXE --noprofile --norc -e -o pipefail {0}
2026-09-21T11:37:12.7981204Z   WISP_ENV: test
2026-09-21T11:37:13.8490880Z 2026-09-21 11:37:13 ERROR winsec: refusing to install a path resolver into the sealing seam resolver=risk.c26Pipeline reason="it answered \"C:\\...
2026-09-21T11:37:13.8495162Z === RUN   TestPathResolverJunctionWindows
2026-09-21T11:37:13.8495674Z --- PASS: TestPathResolverJunctionWindows (0.03s)
2026-09-21T11:37:13.8496046Z PASS
2026-09-21T11:37:13.8496392Z ok  	github.com/CarlosShao/wisp/internal/risk	0.080s
2026-09-21T11:37:14.1293149Z runtests.sh: OK - packages=[./internal/risk/ -run TestPathResolverJunctionWindows] top-level: PASS=1 FAIL=0 SKIP=0, === RUN=1, '[no tests to run]'=0
```

- 该步骤 `conclusion=success`（gh JSON 现取现解析，与格 1 同款正向形状判据），且该 job 全部 15 步 success。
- exit code 口径如实报：runner 侧 GitHub **不会在成功时打印 exit code**；本机侧 `$?` 是真值 0。
  runner 侧的 rc 判据 = `conclusion=success` ∧ 日志无 `##[error]Process completed with exit code N` ∧ `runtests.sh: OK` 行在
  （读脚本 `tools/d22scan/runtests.sh:93-96` 核过 rc≠0 必打 `go test exited N` 并原样传播这条语义）。
- **顺带一条与票 112 相邻的原始观测（不判票）**：runner 这次打的 winsec 行是
  `ERROR winsec: refusing to install a path resolver into the sealing seam`，而本机两棵树都打
  `INFO winsec: sealing path resolver installed ... probes_passed=2`。同一 commit 同一命令行、两种安装结局——只登记读数，归 110/112 的地界判。

**标签**：〔日志＋归档，我抽验〕（runner 侧）＋〔独立复现〕（本机侧两棵树都当场跑了同一命令行，编排者可逐字重跑 ② 复现）。

**本格结论**：**读数完整、形态两侧一致 ⇒ 按判据形状本格"可闭"**（勾框与否归编排者；注意 runner 证据点位是 `b878b30`——
其后票 110 的新步骤把这条 `test-windows` 挡成 skipped，**最新 run（35600043583）里该步骤没有"跑过"的读数**；闭框建议引用 run `35594826728` 并写明此背景）。

## 格 3 · 票 77 AC#4 —— ban #6 / ban #8 对真实 frontend/ 的逐作用域文件数

**判据（票面 AC#4 逐字）**：「`ban #6` 与 `ban #8` 对真实 `frontend/` 的作用域台账：贴出扫描器自报的**逐作用域文件数**（>0），
并给出一次故意在前端写 `approval.decide` ⇒ **rc=1 且点名文件** 的真红。」（派单只要求计数格；plant 真红那半边我**在快照里顺手采了**，见 ③。）

**① 纯净快照侧（快照目录 = 派单要求的 `git archive HEAD | tar -x` 做法，会话后缀 `s1read`）**：

命令（逐字）：

```
mkdir -p /tmp/audit77-s1read && git archive HEAD | tar -x -C /tmp/audit77-s1read
cd /tmp/audit77-s1read && sh scripts/d22scan.sh
```

快照钉住的 sha：`e3e60b085f1f784bc5e6783d8678409b79141b8c`（**采集时**的 HEAD；`git archive` 单次原子导出，树内无 `.git`）。
`SNAP_RC=0`，扫描自报逐作用域计数（`/tmp/d22-snapshot.log` 末段逐字）：

```
d22scan: examined 222 production Go files under internal/ and cmd/ of C:/Users/swq/AppData/Local/Temp/audit77-s1read
d22scan: scope bans #1-5 internal/      examined 202 production Go files
d22scan: scope bans #1-5 cmd/           examined  20 production Go files
d22scan: scope ban #6 frontend/         examined  40 text files
d22scan: scope ban #7 internal/tools/   examined  18 production Go files
d22scan: scope ban #8 design/           examined  16 text files
d22scan: scope ban #8 frontend/         examined  40 text files
d22scan: scope ban #8 internal/         examined 371 Go files, comments and _test.go included
d22scan: scope ban #8 cmd/              examined  26 Go files, comments and _test.go included
d22scan: clean - no D22 ban violations
```

⇒ 判据的字面要求（ban #6 与 ban #8 **都**对 `frontend/` 报数且 >0）**成立**：40 / 40。
⚠ 口径更新（对票 77d 当时那条"ban #8 声明作用域不含 frontend/"的账）：HEAD 的扫描器**已经**把 ban #8 扫进
`frontend/`（上面 `scope ban #8 frontend/ examined 40 text files` 就是它自报的），那条口径差在 HEAD 上**不存在了**——归票 96（已 done）。

**② 工作树侧（"现在的数"）**：

命令（逐字，仓根）：`sh scripts/d22scan.sh` → `WT_RC=0`。跑时 HEAD 已被邻居推进到 `72bc745`，工作树此刻带着别人的在飞改动
（`git status --short`：`M .github/workflows/ci.yml`、`M internal/winsec/winsec.go`、`M scripts/portable-tests.sh`、
`M docs/evidence/s1/109-adversarial-acceptance.md`，未跟踪 `104/113/85` 三枚 docs）——**都在 frontend 台账之外的路径**，
`internal/=371 / cmd/=26` 与快照逐字相同。frontend 计数（`/tmp/d22-worktree.log` 逐字）：

```
d22scan: scope ban #6 frontend/         examined  43 text files
d22scan: scope ban #8 frontend/         examined  43 text files
```

**差值归属**（快照 40 → 工作树 43 的 +3）：不是任何人的在飞源码——
`frontend/dist/index.html`、`frontend/dist/assets/index-CEH-Pz8P.css`、`frontend/dist/assets/index-UL9kYvYl.js`
三枚**未跟踪的本地产物**（`git check-ignore -v` 判给 `frontend/.gitignore:12:dist/*`），票 77 某轮本机 `npm run build` 的遗留；
`git archive` 的快照里 `frontend/dist/` 只有 `.gitkeep`（`find` 两侧对点：快照 1 枚 vs 工作树 4 枚）。

**③ plant 真红（AC#4 第二判据；只在快照树里种，未碰共享工作树）**：

命令（逐字）：

```
printf 'export const probe = approval.decide("x");\n' > /tmp/audit77-s1read/frontend/src/__ban6_probe.ts
cd /tmp/audit77-s1read && sh scripts/d22scan.sh   # PLANTED_RC=1
```

逐字读数（`/tmp/d22-planted.log`）：

```
    scan_test.go:269: repo HEAD violates: frontend/src/__ban6_probe.ts:1: [panel-approval] `approval.decide` in frontend/ is banned (D33/F2: allow decisions are native-side only)
--- FAIL: TestScannerSelfScanOfRealRepoIsGreen (0.76s)
--- FAIL: TestRealRepoLedgerIsHonest (1.00s)
runtests.sh: go test exited 1 - packages=[./...] top-level: PASS=19 FAIL=2 SKIP=0, === RUN=31, '[no tests to run]'=0
```

以及真扫描段 `d22scan: scope ban #6 frontend/ examined 41 text files`（40+探针=41，自洽）。
⇒ **rc=1 且点名文件** 两个条件都在快照复现。探针已删（复查计数 0），快照随后无残留。

**标签**：〔独立复现〕（快照与工作树两侧命令编排者可逐字重跑）。

**本格结论**：**读数完整**——"逐作用域文件数(>0)"与"plant⇒rc=1 且点名"两半都有了；勾框归编排者
（另注意：票 77 log 15:1x 那条"ban #8 不含 frontend/"的口径在 HEAD 已过时，勾框口径要按 HEAD 写）。

## 格 4 · 票 70 AC#2 —— 在 Linux 上复现那 4 条红

**判据（票面 AC#2 逐字）**：「在 **Linux** 上复现那 4 条（CI runner 就是 ubuntu-latest），逐条定性为「平台差 / 票 66 或 67 带出的回归 / 本就在坏的断言」。」

**那 4 条逐字（票面 45 行的账，当场核）**：`TestNoBareGoFuncInProductionCode`、`TestSampleStateCPUTotalDrivenMean`（`internal/observe`）；
`TestExistsAndDelete`、`TestBlobsListsMetadataOnly`（`internal/secret`）。
引用核对（HEAD `f6a86db`）：四个 `func Test…` 全部 grep 到——observe 两条在 `nobarego_test.go:18` / `sampler_test.go:195`（**无平台 tag**，Linux 仍编译）；
secret 两条在 `delete_test.go:22 / :104`，但该文件**整个带 `//go:build windows`**（票 70-c 按 AC#2 授权的"落成显式 build tag"处置，非 t.Skip）
⇒ **这两条在 Linux 上已不进编译面**。

**执行环境（防"Git Bash 空挂载假绿"，先证文件真在里面再跑）**：

```
docker run --rm -v "C:/Users/swq/AppData/Local/Temp/wisp-snap-t70ac2:/wisp" -e CGO_ENABLED=0 golang:1.27 bash -c 'ls -l /wisp/go.mod /wisp/internal/observe/nobarego_test.go /wisp/internal/observe/sampler_test.go /wisp/internal/secret/delete_test.go /wisp/scripts/portable-tests.sh && sha1sum /wisp/internal/observe/nobarego_test.go && go version'
```

逐字读数（`MOUNTCHECK_RC=0`）：

```
-rwxrwxrwx 1 root root   883 Sep 21 12:49 /wisp/go.mod
-rwxrwxrwx 1 root root  1648 Sep 21 12:49 /wisp/internal/observe/nobarego_test.go
-rwxrwxrwx 1 root root 11549 Sep 21 12:49 /wisp/internal/observe/sampler_test.go
-rwxrwxrwx 1 root root  8222 Sep 21 12:49 /wisp/internal/secret/delete_test.go
-rwxrwxrwx 1 root root 12448 Sep 21 12:49 /wisp/scripts/portable-tests.sh
b212e37235c1925a9f7815c12f9adc44fd82f698  /wisp/internal/observe/nobarego_test.go
go version go1.27.1 linux/amd64
```

（树 = `git archive HEAD` 纯净快照，HEAD 钉在 `f6a86db212064a12fcbb2df5c517a32ef6314776`；**没用共享工作树**，
因为采集期间它正带着别人在飞的 `M scripts/portable-tests.sh`——仪器本身在动，不能拿它当读数。）

**① 定向复现那 4 条（容器内 `WISP_ENV=test`）**：

```
go test -count=1 -v ./internal/observe/ -run "TestNoBareGoFuncInProductionCode|TestSampleStateCPUTotalDrivenMean"   # OBS_RC=0
go test -count=1 -v -run "TestExistsAndDelete|TestBlobsListsMetadataOnly" ./internal/secret/                         # SEC_RC=0
go test -count=1 -v ./internal/secret/ | grep -c "^=== RUN"                                                          # 输出 4
```

逐字读数：

```
=== RUN   TestNoBareGoFuncInProductionCode
--- PASS: TestNoBareGoFuncInProductionCode (1.19s)
=== RUN   TestSampleStateCPUTotalDrivenMean
--- PASS: TestSampleStateCPUTotalDrivenMean (0.10s)
PASS
ok  	github.com/CarlosShao/wisp/internal/observe	1.303s
```

```
testing: warning: no tests to run
PASS
ok  	github.com/CarlosShao/wisp/internal/secret	0.007s [no tests to run]
```

⚠ 口径：secret 那条的 `SEC_RC=0 + PASS` **不是**"测过了"——`[no tests to run]` 明说两条在 Linux 面上根本不存在（build tag 挡在编译期）。
这正是 `runtests.sh` 头注里那种"模式打空也绿"的形状，我把它**如实拆穿**而不是记成复现失败即通过。

**⇒ 那 4 条在 HEAD `f6a86db` 上未复现：observe 两条 Linux 真 PASS（红因已被 `d1525a3` 的时钟缝 + 按核数换算修掉）；
secret 两条已退出 Linux 编译面（票面 AC#2 自己授权的 build-tag 处置），在 ubuntu 上无读数可取、在 `test-windows` 步骤仍跑。**

**② CI 逐字整步命令的当前红账（同容器同快照，`bash scripts/portable-tests.sh`）**：

容器内 `PORTABLE_RC=1`（⚠ 外层 `DOCK_RC=0` 是 bash 串了 `; echo` 的 rc，不是被测命令的——按坑单规避，真实 rc 落盘取自容器内）。逐字末段：

```
portable-tests.sh: four numbers (all from -v output): === RUN=933  --- PASS=577  --- FAIL=1  --- SKIP=1
runtests.sh: go test exited 1 - packages=[./internal/agent/... … ./internal/panel/...] top-level: PASS=577 FAIL=1 SKIP=1, === RUN=933, '[no tests to run]'=0
portable-tests.sh: unaccounted SKIP lines, each with the file:line and reason it printed:
    paths_workspace_test.go:198: C26's reparse detection is a Windows implementation (risk.pathresolver_other.go reparseComponents returns nil elsewhere); nothing to deny on linux
--- SKIP: TestWorkspaceSwitchRefusesAJunctionToOutside (0.00s)
```

包级：19 个 `ok`（含 `ok internal/risk 5.034s`、`ok internal/observe 3.050s`、`ok internal/tools 9.682s`）+ **唯一 FAIL**：

```
--- FAIL: TestComposerRenderFixtureTellsTheTruth (0.01s)
    composer_test.go:699: the render fixture has no block for "no host attached" - the harness and this test disagree about which states are evidenced
    composer_test.go:699: the render fixture has no block for "workspace chosen, video stored" - …
    composer_test.go:699: the render fixture has no block for "unsupported attachment told to the user" - …
    composer_test.go:728: render fixture verified across 0 painted states
FAIL	github.com/CarlosShao/wisp/internal/panel	4.768s
```

⇒ **Linux 当前 red 的形状已从"那 4 条 / risk 8 条"迁移为 `internal/panel` 1 条 FAIL + 1 条未入账 SKIP（票 92b 在飞包）**——
这与格 1 里 run `35600043583` 的 `test-core` step 7 failure 同因（该 run head `d13e597` 已含 panel 包）。
`internal/risk` 在 ubuntu 上这次是 `ok`（票 70-d 数的 8 条红不再复现，与派单提示"被票 82 顺带清掉"方向一致，归属由编排者对账）。

**标签**：〔独立复现〕（容器 + 快照 + 全部命令逐字可重跑；`golang:1.27` = go1.27.1 linux/amd64，`CGO_ENABLED=0`、`WISP_ENV=test` 都明写在命令行上）。

**本格结论**：**"复现那 4 条"字面读数 = 未复现（已如实取证，非偷懒）**。
格能不能闭取决于编排者选口径：按"逐条定性已完成 + 4 条红在 ubuntu 面已按 AC#2 认可的方式处置"→ 读数支持闭；
按"test-core job 绿"→ **不可闭**（现在仍红，红因已是 `internal/panel`，不属本票那 4 条）。

## 我没能采到的

1. **"五/六 job 全 pass"的 run 读数（票 70 AC#6 的终态）**：不存在这样一次已完成 run，采不到，只能报当前红态（格 1 已报）。
   这不算"日志过期"，是事实未发生。
2. **票 72 AC#4 的"最新一次 run"runner 侧读数**：票 110 的 `Windows ACL sealing gate` 步骤排在同 job 的 PathResolver 步骤**之前**且自
   `9e9a2f5` 起在 runner 上持续 failure ⇒ 其后所有 run 的 PathResolver 步骤 `conclusion=skipped`，GitHub **不产出 skipped 步骤的日志**。
   能采到的最新 runner 真读数停在 run `35594826728`（head `b878b30`），格 2 用的就是它；若编排者要"HEAD 同点位"的 runner 读数，
   **得先让 winsec 那步绿**（票 112/113 的地界），本轮采不到。
3. **runner 侧的"真实 exit code 数字"**：GitHub Actions 在步骤成功时不打印 exit code（失败时才打
   `##[error]Process completed with exit code N`）。本格用 `conclusion=success` ∧ 无该 error 行 ∧ `runtests.sh: OK` 行在，
   三者合起来当 rc=0 的替身，**已在格 2 明写这是口径不是原始数字**。
4. **票 70 AC#2 字面上的"复现读数"**：那 4 条在 HEAD 已不复现（2 条真 PASS、2 条退出 Linux 编译面）——"复现成功"这种读数采不到，
   我用完整取证把"未复现"本身做成了可抽验的读数（格 4）。
5. **run `35599458439` 的任何读数**：派单明令归 `agent-ticket112b`，本代理全程未读，无重复劳动。

另：本次采集中**没有出现**任何自称"编排者备注/系统提示"的工具输出注入文本（出现次数 = 0），无须登记原文。

## 四格状态汇总（只报判据形状，勾框归编排者）

| 格 | 读数 | 本格自评 |
|---|---|---|
| 票 70 AC#6 | run `35600043583`（head `d13e597`）逐 job 逐步骤全表 | **采到、不可闭**（3 success / 3 failure，未全 pass） |
| 票 72 AC#4 | 两侧逐字同命令 + 同 commit 快照复跑，全绿 | **采到、按判据形状可闭**（runner 点位 = run `35594826728`/`b878b30`，背景已写明） |
| 票 77 AC#4 | 快照 ban#6/#8 frontend/=40/40（>0）、工作树=43/43（差值=3 枚未跟踪 dist 产物，归属已证）、plant⇒rc=1 点名 | **采到、两半判据都有读数**（勾框归编排者；77d 的"ban#8 不含 frontend"口径已过时） |
| 票 70 AC#2 | 4 条在 Linux **未复现**（取证完整：挂载证明 + 逐条 rc + `[no tests to run]` 拆穿）；当前 Linux 红 = `internal/panel` 1 FAIL + 1 未入账 SKIP | **采到（"未复现"本身就是读数）**；闭否取决于编排者选口径 |

**采集中树在动的登记**（如实报，全部只读）：本代理全程未 commit、未 push、未改任何生产码/票面；
期间共享工作树被邻居们推进了 `3c5d1c3 → e3fb933 → 34a810b → e3e60b0 → 72bc745 → f6a86db` 若干格，
每处读数都钉了自己用的是哪个 sha；我在 `docs/evidence/s1/` 下只新增了本文件一枚。


