# 266-r1 交件（编排者额度警戒 · deadline-pause 干净边界停车）

**腿**：266-r1（实现者；判语归非实现者）
**停车时刻**：2026-10-15 16:2x +08（`date` 现读，见下各读数的取数时刻）
**起手锚**：`git rev-parse --short HEAD` = `b74b87d0`（现读 2026-10-15 16:15:04 +0800）
**派单存档**：`.scratch/wisp/probes/orch/20261015/txt/266r1-dispatch.txt`（sha256 `69024b4ceb93e5156b882f9fa764fd962fdc1f99722d5d3602ecd5a2c3f73e9f`）
**状态**：AC#0 读数已取并落件；**AC#1–AC#4 未做**（"未做的格别自评"，见下）。写面零改动：
`git status --porcelain -- scripts tools .github` 空（16:16 与 commit 前各现跑一次）。

## AC#0（已完成：现跑复认 §现量 1/3 ＋ 现读两枚钉今天的真身行号）

取数时刻带 2026-10-15 16:15:04 至 16:18:26 +0800，命令与读数逐条如下。

### ① §现量 1 复认——`scripts/` 不在 `tools/d22scan` 的 scope 名册里（对）

- 尺逐字：`grep -c '"scripts' tools/d22scan/main.go` ＝ **0**，`rc=1`。
- 补尺（本腿加）：`grep -cn '"\.scratch\|"\.github' tools/d22scan/main.go` ＝ **0**。
- 名册（`main.go` 头注逐字）：production scope `internal/` + `cmd/`（`:5-6`）、ban#6 `frontend/`（`:23-24`）、ban#7 `internal/tools/`（`:39`）、ban#8 `design/`+`frontend/`+`internal/`+`cmd/` 的 Go 源（`:40-43`）、ban#9 只判 `internal/`+`cmd/` 的生产 Go（`:55-56`，逐字"adding either tree is a scope change for an owner to approve"）。
- 推论（供 AC#4，⚠ 仍要前后现跑比对）：本票写面 `scripts/**` 与 `.github/workflows/ci.yml` **不在该扫描器任何一处 scope 字面里**，结构上不该带出 d22scan 红名集合的变化。

### ② §现量 2（"删钉调用 ⇒ rc=0 全绿"）**本刻未现跑**

它是 AC#2 正控的对象，本腿 deadline-pause 未做。⛔ 不引用 263-v1 的读数当凭据（票面自己要求）。

### ③ §现量 3 复认——P3 的射程与"按天钝"（对，且入口在 dev 上 inert）

- `scripts/slo-freshness.sh:137` `sample_max_age_days=${SLO_FULL_SAMPLE_MAX_AGE_DAYS:-3}`（常量，默认 3 天）。
- `:548` 打印 `age: N day(s) ... threshold: M day(s)`；`:549` 判 `-gt`；`:550` 红 token `slo-full-sample-stale`。
  判的是**最新有效 sample 产物**（`:332` `sample_artifact=slo-full-report`，单一真源），**不是钉的存在**。
- `:98-99` 逐字 "This is the probe that stays sharp when every run is a machine-contended no-conclusion run"。
  ⇒ 删钉后只要还出一份结论它就安静；真要它响得 D32 连续出不了结论（最早第 4 个整天）。
- 自动入口：`.github/workflows/slo-fresh.yml:46-52` 只有 `schedule`（cron `23 */6 * * *`）+ `workflow_dispatch`；`:41-43` 自陈 cron 只在默认分支评估（dev 上 inert）；P2/P3 需 gh+token，无则 `exit 2`（`slo-freshness.sh:235-253`）。⇒ 与 266-a1 一致：这条 dev 路线上基本不会响。

### ④ 两枚钉今天的真身（现读，**不照抄**票面 §现量 4 的旧行号）

行号取数时刻 2026-10-15 16:15:04 +0800；`git ls-files --eol scripts/slo-check.ps1` = `i/lf w/lf attr/text eol=crlf`（工作树当前 LF，属性写 CRLF ⇒ 外牙读它必须剥 `\r`）。

| 件 | file:line | 内容 |
|---|---|---|
| 能力钉定义 | `scripts/slo-check.ps1:214` | `function Test-ExitCodeInstrument {` |
| 能力钉探测 | `:223` | `Invoke-ExternalProgram -FilePath $comSpec -ArgumentList @('/c', 'exit 7')` |
| 能力钉 ok 句 | `:230` | `instrument self-check ok - external programs are started, waited for, ...` |
| 词面钉定义 | `:233` | `function Test-ScriptShapeNail {` |
| 词面钉自path判 | `:247-249` | `$PSCommandPath` 空即 `Fail-InstrumentBroken` |
| 词面钉模式 | `:251-252` | sigil 由 `[char]36` 组装（自身不匹配自己） |
| 词面钉命中判 | `:255-256` | 命中即 `Fail-InstrumentBroken`（含首命中文本） |
| 词面钉 ok 句 | `:258` | `shape nail ok - this file reads no automatic exit-code variable` |
| **调用点 1** | `:291` | `Test-ScriptShapeNail`（顶格、独立行） |
| **调用点 2** | `:292` | `Test-ExitCodeInstrument`（顶格、独立行、紧邻上一行） |
| "测量之前跑"注释 | `:286-290` | 逐字 "Both nails run before anything is measured (ticket 263)" |

## AC#1–AC#4 状态（未做；按停车纪律不自评、不勾）

- **AC#1 未做**：外牙一行未落地。本腿走得最远的一步＝形已选（见下 `next=`），未写件。
- **AC#2 未做**：正控两发未跑，"删钉必响逐字"不存在（外牙不存在）；还原证明这一半可给且已给（上"写面零改动"＋`git status --porcelain -- scripts tools .github` 空）。
- **AC#3 未做**：无产码 ⇒ 四条偷懒形状自查无可对，未自评。
- **AC#4 未做**：`sh scripts/d22scan.sh` 与票 212 ban #9 读数**前后都没取**；唯一已取＝路径门基线（下）。

## 卫生门读数（本刻已取的真身，供 AC#4 当前基线）

- `sh scripts/check-path-length-budget.sh --with-self-test`：起 16:18:20 / 止 16:18:26 +0800，**rc=0**；
  自控 3/3 PASSED；`denominator read: 13212 tracked paths`；`over-budget=60  covered by roster=60  not in roster=0`；`VERDICT GREEN`；roster=60。
  ⚠ 票面写的 "5594 tracked / 57" 是 2026-10-04 的旧读数；**本刻真身 = 13212 / 60**（引用前先现跑）。
- `tools/d22scan`：未取（见 AC#4）；`ban #9` 命中数：未取。

## next=（下一腿精确剩余；形仍按编排者裁定 ①/②，③ 未具名解冻）

1. **形＝①-A（本腿选定，未落地）**：新尺 `scripts/slo-nail-census.sh`（POSIX sh，读 `scripts/slo-check.ps1`）＋ `.github/workflows/ci.yml` lint 作业**新增一步**（`if: ${{ !cancelled() }}`，落在作业尾部下一个空位，先例＝"Portable tests carrier self-test" 步；⛔ 不动 `on:`/`paths:`/并发块、⛔ 不动任何既有步）。
   - 尺要判四件（都在被检文件之外，具名 `file:line` 落在外牙自己身上）：(a) 两枚函数定义在（`:214`/`:233` 顶格）；(b) 每个函数体含自己的 ok 句 token（`instrument self-check ok` / `shape nail ok`）；(c) 两枚调用点各**恰好一条顶格、未注释**的独立行（今形 `:291`/`:292`）；(d) 读不到文件/缺任一条 ⇒ `exit 2` 具名"无法判定"，⛔ 不静默绿（先例：`Fail-InstrumentBroken`、票 134 AC#4）。
   - 实现要点（已量）：剥 `\r`（属性 CRLF）、去 `<#...#>` 与行注释后再匹配（防"注释里放同形字面"假绿）、根目录从脚本自身位置推（可任意 cwd 运行）。
   - AC#1 的答句（草）：声明（N=2、哪两枚）在**外牙自己的件＋ci.yml 新步**里，"删检查"（删 `slo-check.ps1` 里的调用）不动声明 ⇒ 外牙响；**具名绕过**＝一枚 commit 同时删掉外牙脚本与其 CI 步，今天没有任何东西会发现（`slo-freshness.sh` P1 只锚 slo-full 作业体与触发块）⇒ 别把它包装成"完整牙"。另：静态尺不能证"真跑"，只钉形状；"真跑"的运行时半只能由 ② 给（见 2）。
2. **② 的现读障碍（供裁定）**：②-字面（外层断 slo-check.ps1 stdout 上的两行 ok）在 GitHub Actions 里要求改动**既有** gate 步（管道到文件再断），而本票写面只许**新增**步 ⇒ 字面 ② 不可落地。**②-报告字段变体可落地**：改 `scripts/slo-check.ps1`（在写面内）令两枚钉把"已跑"记账进 `slo-report.json`，再在 `slo-smoke` 作业尾部新增一步读该报告并要求字段在场（缺字段/缺报告＝具名红）。代价：本机正控需跑整条 smoke（build + 采样），而 precheck 会把同波在飞腿的 go.exe 判成 contention 拒采 ⇒ 正控难现跑，须在 hosted runner 上取读数。
3. **AC#2 正控（未跑）**：仓外拷贝 `scripts/slo-check.ps1` ＋ sha256 还原证；变异册至少：删两枚调用（必响，红句逐字）、注释掉两枚调用、只删一枚、删一枚函数定义、删一枚 ok 句、件缺失（⇒ exit 2）；外加未变异件跑绿。
4. **AC#3**：四条偷懒形状（动阈值／放宽 `Set-StrictMode`／把某步改成跳过／`try/catch` 吞）逐条点名"我没做哪一种"。
5. **AC#4**：`sh scripts/d22scan.sh` 与 ban #9 命中数在"落地前/落地后"各现跑一遍，红名集合逐名比对带取数时刻（基线：路径门 13212/60 GREEN，见上）。
