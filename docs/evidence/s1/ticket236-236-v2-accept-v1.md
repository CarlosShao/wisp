# 票 236 — 验收腿 `236-v2` 判语（AC#3＋AC#4 两格；非实现者）

- 派单 2026-10-14 12:2x +08；本腿起手 12:24（HEAD `ade60808`）。窗口逐事时刻见 `.scratch/wisp/probes/236/v2/anchor.txt`。
- 被裁：**AC#3**＝`f88d98ae`（新钉 275 行／`probes/236/r4/evidence.md` 98 行／票尾）＝实现者 `236-r4`；**AC#4**＝`0c7f3c5d`（`.github/workflows/ci.yml` junction 注释）＋`faa8950d`（票尾）＝非实现者 `236-r5`。
- ⛔ 本腿零产码改动、零翻框、零 push；写面＝本件＋`.scratch/wisp/probes/236/v2/**`＋票尾 1 条 bullet。

## 0 判语表（结论先行）

| 格 | 判语 | 翻框建议 | 留档句（给编排者） |
|---|---|---|---|
| **AC#3** 键来源钉 | **成立** | **可翻** | 跨包不变式由 `internal/agent/loop.go` 字节决定（本腿自种 3 发突变全红、还原三方同串）；范围含 `internal/agent` 实测成立；A898 零越（零产码／无熵后门／熵失败不可覆盖已写明）。★两条**注释级**小缺陷随格留档、不阻塞：①钉内熵支行号过期（写 `:1132-1137`／`:1139-1140`，被裁 blob 上真实＝`:1155-1161`／`:1162-1163`）；②`r236r4ModelWireIDs` 注释称"从两 golden 目录现读"，12 枚里 `call_w1`／`chatcmpl-fake` 不在那两个目录（分别在 `internal/tools/loop_approval_test.go:288`／`cmd/wisp/secret_test.go:1133`），且注释里那把 grep 现读得 33 枚唯一 id——断言本体仍成立（两枚都是模型侧值、都过不了锁、也不等于真铸值）。 |
| **AC#4** ci.yml 注释 | **成立** | **可翻** | 两发读数逐字对上（job id／step conclusion／`--- PASS`／`runtests.sh: OK` 行／三枚 FAIL 真名＋`want` 行）；零行为（非注释改动行 0、yaml 结构 52=52 相等）；"指对对象"成立（8.3 短因落在 `--scope=windows` 步的三枚具名 FAIL）。⛔ 未发现缺陷。 |

## 1 验收项 — AC#3（被裁 `f88d98ae`）

### 1.1 钉的正当性：跨包不变式真由 `internal/agent` 字节决定
- 逐行读 `internal/tools/taskid_source_236r4_teeth_test.go`（sha256 `47fb441d…`＝r4 声明，复量一致）。载具＝**真 loop**（`agent.New`＋`RunAsync`）＋**真 C5**（openai-chat 适配器回放 `internal/agent/testdata/golden/tool-then-text.sse`）＋**真 echo 派发**（前置：`calls==1` 且 `CallID=="call_e1"`）＋真名册/真 spill；把组合根 `run.go` 的 `MarkRoot(bg.ID)`／`Backfill(…, res.TaskID, …)` 两个消费调用**原样复做**（`MarkRoot`+`Backfill`+`Look`+`Count`）。
- 断言射程核对（`git show f88d98ae:… | grep -n`，逐枚）：铸造 `id := newTaskID()`＝`loop.go:329` ✓；填 `RunningTask.ID`＝`:330` ✓；传同一 id＝`:333` ✓；落 `Result.TaskID`＝`:382` ✓；接缝 `TaskID: taskID`＝`:676` ✓；铸形定义 `func newTaskID`＝`:1153` ✓。`run.go:1105/:1110/:1124/:1135/:1145` 在 `f88d98ae` 的 blob 上亦 ✓（HEAD 处已因 315-r1 `bfc74ae9` +29 行漂为 `:1134/:1139/:1153/:1164/:1174`——引用者须重跑）。
- "射程含 agent"对核 a4 漏因：a4 点名"按包开尺看不见跨包那一枚"；本钉**住在 `internal/tools`**（不再被 AC#3 两目录过滤）但**牙由 `internal/agent` 字节决定** ⇒ 双向成立（§1.2 实证）。
- 先例核对（引原文行）：锁三段的先例属实——`internal/session/session.go:67` `Valid()`（前缀+小写 hex）、`:94` `Mint()`；便宜向 `internal/session/grants_test.go:152`（手写字面量必须过不了锁）；贵向 `cmd/wisp/ticket224_assembly_test.go:433` 邻域（真装配铸的 id 必须过锁且 ≠ 字面量）。本钉为三段（锁 `r236r4MintedShape`／便宜向两名单／贵向真铸值+相等+接缝+名册链），谓词测试本地、未加导出名 ✓。
- ⚠ 一处措辞留档（非缺陷）:r4 交件 §1 把载具整体写成 `pointer_183_cli_seam_test.go` 两枚的形状（含"真 loop＋真 C5"）；该文件内并无 agent loop/C5（它是 `Bridge`+`TaskBackfill` 的接缝腿）。**钉内注释的归属是准确的**（只把"real components, no stand-in"归给 pointer_183，C5 接线归给 `loop_approval_test.go:60-85`）⇒ 以钉内注释为准。

### 1.2 牙：本腿自种 3 发（全 overlay；3/3 为 r4 没种过的形；红句逐字见 `logs/02-mutations.txt`）
| 发 | 突变（`internal/agent/loop.go`） | rc | 红名（逐句见 logs） |
|---|---|---|---|
| m1 | `:330` 填点换 uuid 形常量 `00000000-0000-4000-8000-000000000000`（模型可给的一枚） | 1 | **4 枚**：`:203` 相等／`:223` 接缝／`:247` 正文读回／`:255` 名册 2 行。**形状锁 `:189` 不红**（uuid 形过锁）⇒ 贵向相等断言是这发唯一咬点，与文件头"便宜向是枚举不是结构"的天花板声明一致 |
| m2 | `:1165` 铸形排版换位（36 字符・4 连字符，第 1 个落 4 位） | 1 | **恰 1 枚**：`:189` 形状锁 |
| m3 | `:382` `Result.TaskID` 换 `"call_e1"` | 1 | **4 枚**：`:203`／`:247`／`:255`／`:267`；接缝 `:223` 不红（接缝仍走真 `taskID`） |
- **交叉读数**：仓里唯一那枚 `internal/agent` 形状尺（`compress_trace_test.go:526`：36 字符/4 连字符）在 m1、m2 下**都 rc=0 PASS**（harness 走同步铸点 `:340`）⇒ 本钉的牙**是新增的**，不是复述旧红；m2 尤其证明本锁（钉死 8/13/18/23 位）比旧尺严。
- **还原证**（overlay 结构性不写盘）：`sha256sum internal/agent/loop.go` ＝ 仓外 backup ＝ `git show HEAD:` blob **三方同串** `cb7f9da7a6d277f341d3b40caee59eecec3bea4e4faee31777b1ad47c3763189`；`git hash-object`＝HEAD blob `0d53277d…`；三发后未突变重跑 Test236 整族 **rc=0**（9 枚）。

### 1.3 正控＋零放宽＋race
- 正控（未突变）：`go test ./internal/tools/ -count=1 -v` ⇒ **rc=0**，顶层 **PASS=220 / FAIL=0 / SKIP=0**（16.768s）；**名册逐名抄下**（`logs/01-baseline.txt`）。与 r4 的 220/0/0 同口径同数。
- 零放宽：`git show f88d98ae --name-status` ⇒ 只 3 件（票尾 M+1 行／evidence A／新测试文件 A 275 行）；**没有任何既有测试或产码文件被改** ⇒ "放宽既有断言"结构上不存在；新文件无 `t.Skip`、唯一导出符号＝用例名。
- race：`go test ./internal/tools/ -count=1 -run 'Test236' -race -v` ⇒ **rc=0**，9 枚顶层 PASS（`logs/03-race.txt`）。
- ⚠ 仪器环境形（本腿自己撞上，非钉缺陷）：GOTMPDIR 写 posix 形 `C:/tmp/…` 时 go1.27 把它当测试临时目录根 ⇒ `mklink /J` 目标带正斜杠 ⇒ 5 枚 junction 用例红（rc=1, 215/5/0）；改反斜杠形 `C:\tmp\…` 后同批全绿。已留档，后程复现请用反斜杠形。

### 1.4 A898 射程核对
- ⛔ 没新造工具面：提交零产码，注册面（`BuiltinTaskEntries` 等）零改动；钉内只加测试本地谓词。
- ⛔ 没动熵源后门：`grep -rn "randRead" --include=*.go .` **零命中**；未在产码加注入缝。
- "覆盖不到熵失败"**已写明**（钉文件头 + r4 evidence §5.1：`rand.Read` 无注入缝、退化支与正常形同 36 字符形状；⛔ 不为测试开生产后门）。
- 其余声明齐：同步铸点 `:340` 不在射程；`cmd/wisp/run.go` 未被真正执行（进程内复做）；模型给 uuid 形值时形状锁排除不了（本腿 m1 实测这一格，且证明贵向仍咬得住）。

## 2 验收项 — AC#4（被裁 `0c7f3c5d`＋`faa8950d`）

### 2.1 读数真伪（gh 现读；逐字见 `logs/04-gh-ac4.txt`）
- run 328＝`38484004521`（head `2c764c0480…`＝dev@2c764c04，failure，2026-10-13T23:24:03Z）；job `115593146575`（test-windows）。
- run 327＝`38332849280`（同 head，failure，01:35:26Z）；job `115086275360`。
- **step 11「PathResolver junction placeholder (real cases tickets 18/20)」两发均 conclusion=success**；逐字正文：run328 `--- PASS: TestPathResolverJunctionWindows (0.02s)` ＋ `runtests.sh: OK - packages=[./internal/risk/ -run TestPathResolverJunctionWindows] top-level: PASS=1 FAIL=0 SKIP=0, === RUN=1, '[no tests to run]'=0`；run327 同形 `(0.03s)`。注释引文逐字对上（注释引到 `=== RUN=1` 截断）。
- 三枚短名 FAIL 真名（两发逐字，且各在 `--scope=windows` 步内；run328 该步 group 头在 log line 4412）：`TestPathResolverShortNameAListDenied`／`TestPathResolverUNCAListDenied`／`TestPathResolverExtendedLengthPrefixAListDenied`，正文均含 `want "C:\\Users\\RUNNER~1\\…"`；该步收尾 `portable-tests.sh: strict runner exited 1 …`＋`##[error]Process completed with exit code 1.` ⇒ **"指对对象"成立**（8.3 短因没消失、只是不砸在 junction 步上）。
- 附：junction 用例在 step10 `--scope=windows` 里也 PASS（0.08s/0.04s）——r5 的附加读数，本腿复核一致；注释未写这条（不构成缺陷）。

### 2.2 零行为核（本腿自跑）
- `git diff 0c7f3c5d^ 0c7f3c5d -- .github/workflows/ci.yml`：全部 `+/-` 行 `lstrip` 后以 `#` 起 ⇒ **非注释改动行＝0**；无步骤增删／`if:`／`continue-on-error`（`if:` 与 `run:` 是纯上下文行）。
- `yaml.safe_load` 两版 rc=0；**结构相等**；steps 52＝52；jobs 名集相同；剥注释后字节逐字相同。
- 新注释区块 `grep -n "today\|FAILS"` ⇒ **rc=1（零残留）**；写有**读数＋日期＋run 号**（`2c764c04`／run 328/327 job 号／2026-10-13 时间戳），⛔ 无"今天为真"句式。
- `runtests.sh` 的非零传播声明核：`tools/d22scan/runtests.sh` 尾 `if [ "$rc" -ne 0 ]; then … exit "$rc"` ✓。
- 改动文件集：`0c7f3c5d`＝ci.yml＋probes 6 件；`faa8950d`＝票尾 +1 行；⛔ 无其他文件、⛔ 未碰他腿脏件。

## 3 工作树与具名没做
- 本腿在 `internal cmd scripts .github` 零改动（他腿在飞 2 枚 `cmd/wisp/config_readers_255.go`／`cmd/wisp/run.go` 非本腿，未读未碰）。
- ⛔ 具名没做：①未复跑 `internal/agent` 整包基线（只跑交叉用例 rc=0）；②未复算 r4 原两发（改用本腿三发，覆盖其两形且更强）；③未验证 8.3 短名"机制怎么修"（归票 252/115）；④未跑 d22scan／长度门（本腿零产码、零 Go 新文件，跑它们不增信息）。⑤未复跑 CI（run 38508311376 在飞）。

## 交件与临时件
- 交件 6 件（逐件 ≤200 行）：本件＋`probes/236/v2/anchor.txt`＋`probes/236/v2/logs/01..04`；票尾 1 条 bullet。
- raw 大日志（留盘未提交，`.scratch/wisp/probes/236/v2/raw/`）：`baseline-tools-v2.txt`（154,056 B，sha256 `31826a2d…`）、`run328-job.log`（920,292 B，`5d68cb65…`）、`run327-job.log`（920,412 B，`c85e37f2…`）、`race-236.txt`（1,658 B，`9d8ea2f1…`）、`baseline-tools-v1-forward-slash-gotmpdir.txt`、`mut-m1/m2/m3.txt`、`xread-m1/m2.txt`。
