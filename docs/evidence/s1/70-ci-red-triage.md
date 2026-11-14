# 票 70 — CI 四步红的分诊（先读后定因；⛔ 未改产品码 / 测试 / 脚本 / 阈值 / golden / 契约）

**取证人**：agent-70-triage（只读腿）　**取证件落点**：`.scratch/wisp/probes/70/`
**判据来源**：run **43549905387**（head `096f247ed`，`dev`，`conclusion=failure`）＋ 追加的**第二发**
run **43551571039**（head `1ce39a651`＝我起手时的 HEAD，同一 job 结构，`conclusion=failure`）。
两发的 job 清单一致：`lint` ✅ `lint-frontend` ✅ `slo-smoke` ✅ `slo-full` ✅ ／
`test-windows` ❌（job `132309710360`）`test-core` ❌（job `132309710522`）。

`test-windows` 的 step 号与结论（`gh run view 43549905387 --json jobs`，逐字）：

```
 4 | failure | Windows ACL sealing gate (internal/winsec's own tests, ticket 110)
10 | failure | cmd/wisp CLI tests (needs the sherpa DLLs staged above, ticket 111 AC#4)
13 | failure | Portable windows tests (proc/secret/config/risk/ball/perm/plugin/llmrecord)
```

> 结论一句话（四步）：**step 4 / step 13 是同一个仪器缺陷的连坐（`portable-tests.sh` 的"名册宇宙"
> 把 `./cmd/wisp/` 拉进来 `go test -list`，在 DLL 暂存之前/之外必然 `0xc0000135`）——确定性红
> 19/19；step 10 是 1 枚真红（`internal/observe` 的 `QuotaLedger` 无锁并发，1/20）；step 7 是
> **另一处**真窗口（`subagent_197.go:373→374` 的"行已发布 / 句柄未挂"，2/20）。后两条都是
> **并发窗口型真缺陷**，不是慢 runner 死 flake，也不是断言该放宽。**

---

## step 4 ｜ **连坐**（步骤红是真的，但 winsec 自身 0 条断言红；`RUN=0` 是真值，不是日志假象）

**判据链（逐字，`job-test-windows.log` 行号为该 job 完整日志的行号）**

```
:280  winsec-tests.sh: platform=windows gate=internal/winsec scope=[./internal/winsec/] package=github.com/CarlosShao/wisp/internal/winsec
:281  portable-tests.sh: explicit scope [./internal/winsec/] IS the winsec tier's own scope - running it as mode=winsec against winsec_pin (ticket 251 AC#1)
      （中间 1241 行全是名册：1221 行裸 Test 名 + 10 个包各一行 `ok  <pkg>  0.0xx s`  ← 不是跑分，见下）
:1525 exit status 0xc0000135
:1526 FAIL	github.com/CarlosShao/wisp/cmd/wisp	0.138s
:1528 portable-tests.sh: go test -list failed for [./internal/winsec/... ./internal/agent/approval/ ./internal/memory/ ./internal/proc/ ./internal/audio/ ./internal/models/ ./internal/risk/ ./internal/tools/ ./internal/config/ ./cmd/wisp/] - the ledger cannot be verified against a test binary that does not build.
:1529 winsec-tests.sh: winsec result line: ok  	github.com/CarlosShao/wisp/internal/winsec	0.069s
:1530 winsec-tests.sh: four numbers (all from -v output): === RUN=0  --- PASS=0  --- FAIL=0  --- SKIP=0
:1531 winsec-tests.sh: gate is RED (rc=1) for scope=[./internal/winsec/]
```

**为什么 `RUN=0` 是真值**：`scripts/winsec-tests.sh:206` 把这一发整个委派给
`scripts/portable-tests.sh`，后者在**跑任何测试之前**先做"名册复验"——
`scripts/portable-tests.sh:841-857`（`ledger_pkgs` ← ledger 里出现过的包；`universe_file`；
`:851` 一发 `go test -list '.*' "${scope[@]}" "${ledger_pkgs[@]}"`，非零即 `:856 exit 1`）。
这一发在 `:851` 就死了，`:1068` 之后的 run 块**从来没被执行**，所以
`winsec-tests.sh:208-212` 从捕获字节里数出来的四数必然是 `RUN=0 PASS=0 FAIL=0 SKIP=0`。
⇒ **不是"名册为 0"的仪器假象，是"一枚都没跑"的诚实读数**；红因不在 `internal/winsec`。

**为什么 `go test -list ./cmd/wisp/` 会死**：`exit status 0xc0000135` = Windows
`STATUS_DLL_NOT_FOUND`（进程**加载期**就死，早于任何 Go 代码）。`cmd/wisp` 的测试二进制链接
sherpa/onnxruntime，`go test -list` **不经过** `scripts/wisp-cli-staged-tests.sh` 的 DLL 暂存；
而 step 4 排在夹具之前：

- `.github/workflows/ci.yml:954`　step 4（且 `:958` 明写"it stays FIRST inside the job"）
- `.github/workflows/ci.yml:1089`　`Fetch voice fixtures`
- `.github/workflows/ci.yml:1118`　`cmd/wisp CLI tests`

**这条是什么时候被点着的（可核）**：`go test -list` 的"宇宙"那行是 2026-09-21 的老代码
（`git blame -L 851,857` ⇒ `c8c828e32e`），但 `./cmd/wisp/` **第一次**出现在 ledger 里是
**2026-11-13**：`git blame -L 831,835` ⇒ `02939b8a2`（票 111 wave 8b，`TestTicket223…`）＋
`7386f221c`（wave 10，另外 4 行）。⇒ **老机制 × 新 ledger 行 = 新红**，且从那天起每一发都红（见 附 A）。

**`ok … internal/winsec 0.069s` 为什么不能当"跑过"**：`go test -list` **自己**会在名册末尾打
`PASS` / `ok  <pkg>  <dur>` 尾巴（本机实测，见 §复核命令 R2）。`scripts/winsec-tests.sh:235` 的
GUARD 2 找的正是这个形状，于是它被**名册尾巴**满足（`winsec result line` ＝ `:1529` 那一行），
而实际上一枚测试都没跑。⇒ GUARD 2 在这里**是真阳性判据、假阳性输入**：它的"包级结果行"被
`-list` 的尾巴顶替了。这条我没有改（GUARD 属 `scripts/**`，不在写面），登记给 owner。

**修法候选（谁去修）**
1. **（推荐，主流 harness 形状）把名册复验做成逐包隔离**：`portable-tests.sh:851` 现在一发
   `go test -list` 打 `scope+ledger_pkgs`，任一包载入失败就整 tier 退。改成**逐包** `-list`，
   某包列不出来只让**它自己的 ledger 行**报 `stale/不能复验`，其余行照验、tier 照跑。
   owner：**编排者（票 71 是 `scripts/**` 的账）**，或另开一张"名册复验不是单点故障"的票。
2. **名册阶段也走 DLL 暂存**（把 `wisp-cli-staged-tests.sh` 的 shape 用到 `-list` 上）。
   治本，但等于把 cli 载具的逻辑复制到名册阶段；owner：同上。
3. ⛔ **不推荐**：把 step 4 挪到夹具之后——违反 `ci.yml:958` 的显式约束，而且 **step 13 本来就在
   夹具之后**（`:1316` vs `:1089`）也一样红（见下），所以"补夹具"根本不是这条的因。

---

## step 10 ｜ **真红**（1 枚断言红；不是连坐、不是"具名 SKIP 没 SKIP"；**派单指的两个行号不是断言**）

### (a) 派单要的 `leg_sink_nail_131_windows_test.go:358` / `:396` —— 逐字，两处都**不是断言**

源码（`cmd/wisp/leg_sink_nail_131_windows_test.go`，仓内实文件）：

```
:358	s := readLegSink131(t, dataDir, stdout, stderr)
:396		t.Logf("MODELS LEG RECORDED ON DISK: %q (install dir=%s, exit=%d)", hit.Msg, s.dir, code)
```

CI 日志里这两行号对应的**全部是日志输出**（`:358` 之所以报到调用者行号，是 `readLegSink131` 里有
`t.Helper()`）：

```
leg_sink_nail_131_windows_test.go:358: leg sink C:\Users\runneradmin\AppData\Local\wisp-test-gotmp\cli-20261114-031758-432\TestAC2ModelsLegBooksItsHandOffVerdictOnDisk4144639867\001\logs: 2 record(s) [INFO/wisp: persistent log sink installed INFO/models: hand-off REFUSED id=absent-in-signed-manifest-131 state=Error err=model: model id "absent-in-signed-manifest-131" not in signed manifest]
leg_sink_nail_131_windows_test.go:396: MODELS LEG RECORDED ON DISK: "models: hand-off REFUSED id=absent-in-signed-manifest-131 state=Error err=model: model id \"absent-in-signed-manifest-131\" not in signed manifest" (install dir=C:\Users\runneradmin\AppData\Local\wisp-test-gotmp\cli-20261114-031758-432\TestAC2ModelsLegBooksItsHandOffVerdictOnDisk4144639867\001\logs, exit=1)
--- PASS: TestAC2ModelsLegBooksItsHandOffVerdictOnDisk (0.02s)
```

⇒ 这两个行号所属的用例 **PASS**。派单里"该处断言红"的前提不成立（`:358` 是夹具读取调用、
`:396` 是 `t.Logf`）；这两处**没有可抄的失败断言文本**。

### (b) 这一发真正的红：`Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands`

```
portable-tests.sh: END ANCHOR (ticket 314): roster=619 top-level case(s) from `go test -list '.*' ./cmd/wisp/`; the run must reach the last one and print its verdict: TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask
--- PASS: TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask (0.37s)
FAIL
FAIL  github.com/CarlosShao/wisp/cmd/wisp	716.00s
portable-tests.sh: four numbers (all from -v output): === RUN=798  --- PASS=599  --- FAIL=1  --- SKIP=19
--- FAIL: Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands (11.16s)
    subagent_selfapproval_197_test.go:232: exit 1
```

**END ANCHOR 成立**：`619 = PASS 599 + FAIL 1 + SKIP 19`，最后一名 `TestZZ315A1…` **打印了
`--- PASS`** ⇒ 这一发**跑完了**，不是"没跑完"，也不是被前一步连坐（step 7/8/9 全 success）。

**红的内容**：`cmd/wisp/subagent_selfapproval_197_test.go:232` 是
`t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, …)`，即"**根任务退出码非 0**"，
不是某个审批断言不成立。根任务为什么 exit 1，日志里有它自己的尸检（逐字，`job-test-windows.log:7821`）：

```
level=ERROR msg="goroutine panic recovered" goroutine=agent-task-11289555-… owner=agent
  recovered="runtime error: invalid memory address or nil pointer dereference"
  … observe.(*QuotaLedger).roll(0x31b54c228900)  internal/observe/quota_44.go:256
     observe.(*QuotaLedger).Spend(…)              internal/observe/quota_44.go:284
     agent.(*QuotaGate).Decide(…)                 internal/agent/quota_44r4.go:155
     agent.(*gatedChainElement).Stream(…)         internal/agent/quota_44r4.go:177
     llm.(*ChainRunner).Stream(…) / cmd/wisp.(*gatedChainDispatch).Stream(…) run.go:1919
     agent.(*Loop).stream / run / RunAsync / observe.(*Registry).run(goroutine.go:340)
```

紧接着 `wisp run: 退出码 1（error_class=无）`、屏幕上是"任务  结束（，0 轮，0 次工具调用…）"。

**定因（我的读法，标明推断层）**：`internal/observe/quota_44.go:221-228` 的 `QuotaLedger`
**没有任何锁**，而 `roll()`（`:247-260`）在每一次 `Spend()`/`Post()` 上都要**写**
`dayKey`/`monthKey` 并 `clear()` 三个 map；`QuotaGate.Decide`（`internal/agent/quota_44r4.go:155`）
在**每次派发前**调 `g.ledger.Spend()`。Test197 的用例形态（根代理在跑、子代理同时在跑）正是
**两个 goroutine 同时进 `Decide`** 的形状 ⇒ 读到被另一侧改了一半的 string 头（torn read），
表现就是"nil pointer dereference"。**这一层是推断**：日志只给到"panicked inside roll"，
给不到是 `l` 为 nil 还是 `l.now` 为 nil（两条构造路径都有 nil 兜底：
`NewQuotaLedger` `:232-242`、`NewQuotaGate` `internal/agent/quota_44r4.go:76-78`），
**我不编排一个我量不到的结论**。可确证的是：**最近 20 发里它红 1 发（1/20）**；而紧挨着的下一发
（`43551571039`，当前 HEAD）红的是**另一处**并发窗口——`Test221` 家族的"发布/挂句柄"窗口，
见 step 7。两条都不是"夹具写错"，是并发面缺护栏。

**修法候选（谁去修）**：`internal/observe`（票 44 家族）给 `QuotaLedger` 加互斥（`roll` 在 `Lock`
下、`Spend` 出快照），或让 `QuotaGate.Decide` 用 `Spend()` 的**只读快照**接口。**不是**改
Test197 —— 它量的是"根任务真的把活干完"，没有放宽的余地。⛔ 我没改一个字。

### (c) 19 枚 SKIP **不是**本步红因（实测）

最近 20 发里 **step 10 有 17 发 success**，其中多数四数同样是 `FAIL=0 --- SKIP=19`
（例：run `43545130245` 的 step 10 = **success**，四数 `=== RUN=798 --- PASS=600 --- FAIL=0 --- SKIP=19`）。
⇒ **同样的 19 枚 SKIP，step 10 照样绿**（判据件 `.scratch/wisp/probes/70/per-run-steps.json`）。
（cli 支路 `scripts/portable-tests.sh:1119` 调载具时**不传** ledger 的 `-skip`
模式，那些 SKIP 只被 `:1131-1134` 打印、不进 rc。这条本身是个账实缺口——ledger 那 5 行
`./cmd/wisp/` 的 row 在 cli 支路**从来没被真正过滤过**——但**它不是今天的红因**。）

---

## step 13 ｜ **连坐**（同一根因；**不是"跑完有红"**：END ANCHOR 段落**没响**，四数一行都没打）

**逐字读数**

```
:8839 portable-tests.sh: go test -list failed for [./internal/proc/ ./internal/secret/ ./internal/config/ ./internal/risk/ ./internal/ball/ ./internal/perm/ ./internal/plugin/ ./internal/plugin/manifestslot/ ./cmd/balldebug/ ./cmd/llmrecord/ ./internal/session/ ./internal/projctx/ ./internal/audio/ ./internal/agent/approval/ ./internal/memory/ ./internal/proc/ ./internal/audio/ ./internal/models/ ./internal/risk/ ./internal/tools/ ./internal/config/ ./cmd/wisp/] - the ledger cannot be verified against a test binary that does not build.
:10369 exit status 0xc0000135
:10370 FAIL	github.com/CarlosShao/wisp/cmd/wisp	0.096s
:10371 ##[error]Process completed with exit code 1.
```

- **`--- FAIL` 逐枚：0 枚**（step 13 整段＝`job-test-windows.log` 行 `8799…10372`，一条都没有）。
- **四数**：`portable-tests.sh: four numbers` **一次都没打印** ⇒ 没跑到跑测试那一步（`:1068+`）。
- **END ANCHOR**：**没响**，而且它**不可能**在这里响——`scripts/portable-tests.sh:1029`
  是 `if [ "$mode" = cli ]`，END ANCHOR 只属于 cli tier。step 13 是 `mode=windows`。
  ⇒ 派单里"若锚响了就把那段抄下来"的前提在 step 13 上**不成立**；真正响的是更早的那道
  **名册宇宙闸**（`:851`），也就是 step 4 的同一个闸。
- **17 个 `ok  …/internal/xxx  0.0xx s` 不是跑分**：它们是 `go test -list` 的尾巴。本机实测
  （§复核命令 R2）：`go test -list '.*' ./internal/proc/` 的最后一行就是
  `ok  	github.com/CarlosShao/wisp/internal/proc	0.087s`。⇒ 谁把这 17 行读成"17 个包过了"，
  读到的就是名册。
- ⚠ **step 13 在夹具之后仍然红**（`.github/workflows/ci.yml:1316` vs `:1089`），
  因为 `go test -list` **不走** `scripts/wisp-cli-staged-tests.sh` 的 DLL 暂存。这条排除了
  "把夹具提前就好"的方案。

**修法候选（谁去修）**：同 step 4 的 1（逐包隔离）/ 2（名册也暂存 DLL）。owner 同 step 4。

---

## step 7 ｜ **时序性真缺陷（narrow race）**，不是"慢 runner 死 flake"，也不是"断言该放宽"

> 派单口误更正：这枚用例**不在** `internal/agent/`。`internal/agent/` 下没有
> `task_cancel_221_legs_test.go`；实际文件是 **`internal/tools/task_cancel_221_legs_test.go`**。

### ① 那个断言的上界是多少、在等什么

`internal/tools/task_cancel_221_legs_test.go:450-453`：

```go
	stopA := x.cancel(t, parent197, aID)
	if stopA.IsError {
		t.Errorf("父任务停甲号孩子被拒（拒绝类正控不许把有权的那一支也一起堵死）：%s", stopA.Text)
	}
```

**这个断言本身没有上界**：它是一次同步的 `task.cancel` 派发，回来立刻判 `IsError`。
**`30.01s` 不是它的上界**，而是它后面那段 park 循环（`:458-464`）等到的东西：
两枚 `task.spawn` 一直挂在 `heldProvider`（`internal/tools/subagent_197_test.go:60-79`）里，
直到 **`internal/tools/bridge.go:28` 的 `const DefaultToolTimeout = 30 * time.Second`** 到点、
两枚 spawn 各返回"超时"结果，`parkA/parkB` 才收口。⇒ **30.01s 是结果，不是判据**，
更不是"撞上上界所以是 flake"的证据。

### ② 同一枚用例在最近 20 发里失败过几次（给分母）

口径：`gh run list --branch dev --limit 20`（20 发全取，**分母＝20**：每一发的 `test-core` job 与
它的 `Portable package tests` 步都有结论，绿的也算样本，因为绿＝这枚用例跑过且过了）；
逐发 `gh run view <id> --json jobs` 取三个 step 的结论，再 `--log-failed` grep `--- FAIL: <name>`。

| 发 | head | test-core / step7 | step 10（cli） | 本族红（`--- FAIL`） |
|---|---|---|---|---|
| 43551571039 | `1ce39a651` | failure / **failure** | success | **`Test221ParentStopsItsOwnChildRowAndStreamSettle`**（`legs_test.go:322`） |
| 43549905387 | `096f247ed` | failure / **failure** | failure | **`Test221SubagentCannotStopSiblingOrItself`**（`legs_test.go:452`） |
| 43549913600 | `15cf0f94b` | success / success | success | — |
| 43549248016 | `bd4b9037e` | success / success | failure | —（cli 红：`Test390r1BallClickPlayingInterruptAbortsAndWalks27`） |
| 43549226873 … 43451733615（14 发） | … | success / success | success | — |
| 43442602845 | `843ad88aa` | success / success | failure | —（cli 红：`Test16r1PremortemModelSha256CorruptRefusesTheHandOff`） |
| 43496333545（整发绿） | `d49caf4fc` | success / success | success | — |

逐发原始件：`.scratch/wisp/probes/70/per-run-steps.json`、`history/<id>.log`。

⇒ **分母 20：`Portable package tests` 18 success / **2 failure**（10%）。两发命中的是**同一族、
同一句拒绝文本、两枚不同用例**——这正是"窗口"而不是"某枚用例写错了"的指纹。
同一枚用例（`:452`）在窗口内 **1/20**；本族 **2/20**。

### ③ 本地三发读数

```
go version go1.27.1 windows/amd64      (CI 的 test-core 是 go1.27.4 linux/amd64)
=== RUN   Test221SubagentCannotStopSiblingOrItself
--- PASS: Test221SubagentCannotStopSiblingOrItself (0.01s)
--- PASS: Test221SubagentCannotStopSiblingOrItself (0.00s)
--- PASS: Test221SubagentCannotStopSiblingOrItself (0.00s)
ok  	github.com/CarlosShao/wisp/internal/tools	0.533s
```
（`go test -count=3 -run '^Test221SubagentCannotStopSiblingOrItself$' ./internal/tools -v`；
原始输出 `.scratch/wisp/probes/70/local-test221-count3.txt`）⇒ **本机 3/3 绿**，窗口很窄，
本机（12 核、空载）打不开。**这不构成"CI 上是 flake"的证据**，只说明窗口窄。

### ④ 定因（`文件:行`）

CI 两句拒绝文本 **自己就自证了窗口**：

```
:452 没能停掉 1d300582-…：任务 1d300582-… 没有在跑（kind=subagent，state=Thinking），没有可停的句柄
:322 没能停掉 4a599e49-…：任务 4a599e49-… 没有在跑（kind=subagent，state=Thinking），没有可停的句柄
```

`state=Thinking` ＝ 这一行**正在跑**，而 `internal/tools/task.go:526-529` 的措辞是
"没有在跑 … 没有可停的句柄"——**括号里的两个事实互相矛盾**，因为它是 `cancel == nil` 那一支
（`task.go:519-529`：`known==true` 但 `r.cancel[taskID]==nil`）。而"行已发布、句柄未挂"的窗口
就在 `internal/tools/subagent_197.go`：

```
:337-345  opt.AdmitTask = func(taskID string) func() { t.d.Roster.PublishSubagent(taskID, parentID, label); … }
          ← 子 goroutine 侧：第一次模型调用之前**只发布行**
:371-374  childCtx, cancelChild := context.WithCancel(context.WithoutCancel(ctx))
          bg := child.RunAsync(childCtx, prompt)
          t.d.Roster.PublishSubagent(bg.ID, parentID, label)
          t.d.Roster.AttachCancel(bg.ID, cancelChild)      ← 句柄在**发布之后**才挂上
```

测试的放行点 `harness221.spawnHeld`（`:132-158`）等的是"第一次模型调用已经开始"
（`started` 通道），而那一刻**只保证行已经发布**（子侧 `AdmitTask`），**不保证句柄已挂**
（父侧 `:374` 可能还没执行到）。负载一高，父 goroutine 在 `RunAsync` 返回后被抢占一下，
窗口就打开 ⇒ 测试拿到的行在飞、句柄为空 ⇒ 真拒绝被当成误判。
（顺带：`DetachCancel`（`:466`）不是嫌疑人——孩子还挂在 provider 里，没 join。）

**定性**：这是**产品码里的发布顺序缺陷**，测试的判据（父任务停自己的孩子必须成功）是对的、
不该放宽；慢 runner 只是把窗口放大。按派单四选一 ⇒ **时序 flake 与真缺陷是同一件事的两面，
我判"真缺陷（窗口型）"**，不是"断言太紧"。

**修法候选（谁去修）**：`internal/tools`（票 221 / 197 家族）——把"发布行"和"挂句柄"做成**一次**
发布，且**必须在子侧第一次模型调用之前**完成。最省事且语义最正的一形：
`PublishSubagentWithCancel(taskID, parentID, label, cancelChild)`，在 `subagent_197.go:337-345` 的
`AdmitTask` 里（那里已经拿到 `taskID`，且早于第一次模型调用）调用一次；父侧 `:373-374`
收敛为同一方法的幂等覆盖（或删掉）。这样"行可见 ⇒ 句柄已挂"成为不变量，
`task.go:526` 那句"没有在跑"也就永远不会用在仍在飞的行走上。
⛔ 我没有动 `internal/tools/**` 一个字。

---

## 附 A：这条红的"年龄"（可核）

- `scripts/portable-tests.sh:851`（名册宇宙）＝ `c8c828e32e`（2026-09-21）。
- ledger 里第一条 `./cmd/wisp/` 行 `:831` ＝ `02939b8a2`（**2026-11-13**，票 111 wave 8b），
  `:832-835` ＝ `7386f221c`（同日 wave 10）。
- ⇒ step 4 / step 13 的红**从 2026-11-13 起每一发都在**：最近 20 发里**19 发红**（唯一整发绿的是
  `43496333545`），这 19 发的 `test-windows` 第 4 步与第 13 步**全部 failure**（`per-run-steps.json`），
  日志里各是 `go test -list failed for` ＋ `exit status 0xc0000135`（第二发 `43551571039`＝当前 HEAD
  也一样）。**这不是偶发，是确定性红。**

## 附 B：探针件清单（`.scratch/wisp/probes/70/`）与复核命令

**入库的**：`job-test-windows.log`（2.5MB，本文件所有 `:NNNN` 行号都指它）、`job-test-core.log`（3.1MB）、
`per-run-steps.json`（20 发逐 step 结论）、`run-list-20.json`、`run-jobs.json`、`local-test221-count3.txt`、
`scan-20runs.ps1`、`history/summary.json`。
**只落盘不入库**：`history/<id>.log` 19 份 `--log-failed` 原始 dump（约 30MB，R5 一条命令即可重取）。

```bash
# R1 取两发逐 job / 逐 step 结论
gh run view 43549905387 --json headSha,conclusion,jobs
gh run view 43551571039 --json headSha,conclusion,jobs

# R2 「go test -list 会打 ok 尾巴」的本地实测（名册 ≠ 跑分）
go test -list '.*' ./internal/proc/ | tail -3      # 末行形如: ok  …/internal/proc  0.087s

# R3 step 4 / 13 的红因（一行）
gh run view --job 132309710360 --log | grep -n 'go test -list failed for\|0xc0000135'

# R4 step 10 的 END ANCHOR 与四数
gh run view --job 132309710360 --log | grep -n 'END ANCHOR (ticket 314)\|four numbers' | tail -4

# R5 本族在最近 20 发的命中（分母口径见 §step7②）
gh run view <id> --log-failed | grep -n 'Task221\|task_cancel_221_legs_test.go:'
```

## 附 C：我这个腿**没做**的事（边界）

- ⛔ 没改任何产品码 / 测试 / 脚本 / 阈值 / golden / 契约；本文件之外只落了证据与探针件。
- ⛔ 没跑整包 `./cmd/wisp`；本地只跑了 `internal/tools` 的**单枚**用例（`-run`，3 发）。
- ⛔ 没 push。
- 没做的事里唯一"本来想要"的：`QuotaLedger` 那条 nil 到底是 `l` 还是 `l.now`、以及
  `Test197` 能否在真机复跑——两者都要跑 `./cmd/wisp` 的单枚用例（要 DLL 暂存 + 夹具），
  超出本腿写面与派单的"⛔ 跑整包"边界，留给 owner。
