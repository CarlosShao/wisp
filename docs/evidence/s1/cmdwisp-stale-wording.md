# cmd/wisp 过期文字清场 — 票 228 AC#8 ＋ 同段两枚同族假话 ＋ usage 托盘 Exit 假话（写腿 cmdwisp-wording，2026-11-08）

- 起手锚 `66a03a3201cdcd808158f7674a8b592b08ac2ae2`（branch dev）；产码 commit `9aaf601c8`；本表收尾 commit 见文末。
- 射程＝`cmd/wisp/**`（本波独占写面）＋ 本证据件 ＋ 票 228／245 各 Progress log 一行。⛔ 未动 `internal/**`、`frontend/**`、`scripts/**`、specs、thresholds、golden；⛔ 未翻任何 AC 框、未动任何 Status 行。
- 门禁台面＝`git clone --shared` 导出树 `D:\tmp\wisp\20261108\cmdwisp-wording\tree`（commit `9aaf601c8`；盘符根零污染落点），非共享树。

## 1. 票 228 AC#8（本票判据格）：`cmd/wisp/main.go:8` 假 godoc

**改前**（逐字，锚 `66a03a32` 上 `sed -n '8p'`）：

```
// frozen D38(e) 10-step shutdown order. The floating ball GUI is ticket 07.
```

**改后**（逐字，commit `9aaf601c8` 上 `sed -n '5,10p'`）：

```
// The no-args path (ticket 03 first, extended since) boots the runtime
// skeleton (Job Object, per-session single instance, goroutine registry - the
// init-time self-checks) and runs the resident event loop until an exit
// signal or a stop request fires, then exits through the frozen D38(e)
// 10-step shutdown order. The floating ball is not a pending ticket: the
// resident leg of this same process assembles it at startup
// (startResidentBall in resident_ball_windows.go, called from runResident).
```

**代码事实（新句每一半的出处）**：
- 球由本文件所在常驻腿装配：`cmd/wisp/main.go:81` `if len(args) == 0 {` → `:87` `runResident()`（无参数分支）→ `cmd/wisp/resident_windows.go:419` `rb := startResidentBall(rt.Registry, ...)` → 定义在 `cmd/wisp/resident_ball_windows.go:382` `func startResidentBall(...)`。四点本腿现跑 grep 复认。
- 环路等的是"exit signal 或 stop request"：`internal/proc/boot_windows.go:182` `func (rt *Runtime) RunEventLoop() string`，体内 `:184` `signal.NotifyContext(..., os.Interrupt, syscall.SIGTERM)`、`:185` `stopSignal := rt.StopSignal()`、`:194` `case <-stopSignal: return "stop request (" + rt.StopReason() + ")"`（票 228 AC#11 的 latch，注释自陈）。

⛔ 未新增扫注释票号的词面型仪器（票面判据第 ③ 条照守）。

## 2. 同段两句同族假话（同批改真；归属＝台账 `A477`「main.go 那族句子判归 228 的 AC#7/AC#8 文案面」）

### 2a. 「empty event loop」（改前原第 7 行）

**改前逐字**：`// and runs an empty event loop until an exit signal, then exits through the`
**假在哪**：环路今天不空——`boot_windows.go:186-201` 现读：每 tick `rt.handleActivation()`（激活泵），`select` 等 `sigCtx.Done()`／`stopSignal`／`50ms` tick；托盘 Exit、session-end 都经 stop-request 位进这只环路（`cmd/wisp/lifecycle_windows.go:529` `rt.RequestStop("session-end (" + why + ")")`）。
**改后逐字**：`// init-time self-checks) and runs the resident event loop until an exit` ＋ `// signal or a stop request fires, then exits through the frozen D38(e)`。

### 2b. 「wisp run still echoes its task text as a placeholder」（改前原第 9 行）

**改前逐字**：``// `wisp run` still echoes its task text as a placeholder (agent loop: ticket``
**假在哪**：agent loop 已落地且 `wisp run` 就是它的组装根——`cmd/wisp/run.go:1-27` 文件头文档链自陈 `agent.Loop(provider, bridge, AdmitTask=D47 hook)`；`:285` `return rt.execute(s.taskCtx, task)`；usage 自己（ticket 12 行）都写「run one task end to end through the agent loop」。`grep -n "placeholder\|echoes" cmd/wisp/run.go`＝零命中。
**改后逐字**：``// `wisp run` runs one task end to end through the agent loop (ticket 12);``

## 3. usage `-h` 正文假话：托盘 Exit「no stop path」（改前原第 38-39 行）

**改前逐字**（锚 `66a03a32` 上 `sed -n '37,40p'`）：

```
                   task entry says so out loud and the ball keeps running;
                   Ctrl+C stops the process, and the tray's Exit item still has
                   no stop path attached to it)
```

**假在哪**：票 228 AC#11（乙-1 stop-request 位）已落地——托盘 Exit 执行者链完整：
- `cmd/wisp/resident_ball_windows.go:151` 注释自陈「withStopRequest hands the ball host the tray's 退出 executor (ticket 228」；`:158-163` `trayExit()` 挂 `Events.OnTrayExit`，回调 `recordTrayExit(h.requestStop)`；
- `cmd/wisp/resident_windows.go:846` `return rt.RequestStop`（hoisted accessor，即 requestStop 的真身）；
- `internal/proc/boot_windows.go:194-200` 收 stop-request 后正常返回 ⇒ 走满 D38(e) 十步（`228-r3c` 突变复跑腿：ⓐ-first／ⓑ-1 四枚钉两台面 `4/4/0/0` 全绿，M1/M2 反形具名红）。

**改后逐字**（`9aaf601c8` 上 `sed -n '36,40p'`）：

```
                   owns one (ticket 246 AC#7: no interactive console means the
                   task entry says so out loud and the ball keeps running;
                   Ctrl+C stops the process, and so does the tray's Exit item,
                   which asks the resident loop for the orderly D38(e) shutdown
                   instead of exiting in place (ticket 228 AC#11))
```

**钉面核对**：`leg_dispatch_gate_133_test.go:181` `usageBare133` 钉的是 `(?m)^  wisp {2,}\S` 行首形状——本次改的是它后面第 33-37 行的段落，`  wisp   ` 行首一字未碰（本腿现跑 `grep -n "wisp " cmd/wisp/main.go` 复认行首形状原样）。`grep -rn "no stop path\|Exit item still" cmd/wisp/*_test.go`＝零命中 ⇒ 无仪器钉旧句。

## 4. 票 245 AC#6①：复认「已被 245-r3 落地」，本腿零改动

派单假设这句还是「its four global hot keys」——**现读已不是**。commit `9e84ae42e`（245-r3：稳态热键钉改判裸键 AC#6②＋usage 半假话改条件句 AC#6①＋cancel 配置钉 AC#4）已把 usage 改成条件事实句。现文逐字（`9aaf601c8` 上 `sed -n '25,31p'`）：

```
  wisp             GUI resident process: boots the runtime skeleton, hosts the
                   floating ball window, its tray icon and its three always-on
                   global hot keys in this same process (summon, mute, panel;
                   the cancel slot is the fourth key, bound only while a
                   confirmation card is waiting and handed back when it leaves,
                   so an idle wisp holds no bare key against the rest of the
                   desktop), shows confirmation cards on that ball, and takes
```

**与代码事实对撞（本腿现跑，非转抄）**：
- 稳态名册不含 cancel：`internal/ball/ball_windows.go:857` 逐字 `// set is the idle set, and the idle set does not include cancel (ticket 245).`
- 借期是唯一绑定路径：`:920` 逐字 `// Ticket 245: this is the ONLY path that registers the cancel id. An idle ball`；`internal/ball/hotkey_windows.go:13-14` 逐字 `// never registers this id at all (HotkeyStandby), and takeEscWithAcc /` ＋ `// releaseEscWith are the ONLY path that binds and unbinds it.`；`:549` `escBorrowAcc`＝借期注册 `VK_ESCAPE` 加速键。
- 判语：**245-r3 的改法为真、措辞与代码事实一致** ⇒ 按「⛔ 为凑数改对的文字」零改动。245 票面 Progress log 缺 245-r3 那一笔（票面 log 止于 245-r1）——那是收表腿/编排者的账，本腿不代补判定，只在 log 里记本腿复认一行。

## 5. 扫面清单（同族候选；本腿改了 1、列出不改 5）

尺：`grep -rn "TODO\|FIXME\|XXX\|待填\|placeholder\|ticket 07\|stub" cmd/wisp/*.go | grep -v _test` ＋ 逐处读上下文。

| # | 落点 | 逐字 | 判 |
|---|---|---|---|
| 1 | `cmd/wisp/resident_sink_nail_127_windows_test.go:30-31`（注释） | `//   - No window. Ticket 07 has not landed, so this leg boots, prints and parks` ＋ `//     in the empty event loop (MainWindowHandle = 0).` | **假话（票 07 已落地）**。但属票 243 普查的「ticket 07 注释族」归属面（该票 Census 逐处判归属），本腿不跨票改，**列清待裁**。⚠ 该注释描述的是 127 那族测试的台面形状（测试进程不开窗），改动前需先核测试语义，不是纯文字活。 |
| 2 | `cmd/wisp/doctor.go:45`（运行时字符串） | `" (placeholder until C29 lands; hardcoded into buildinfo per SPEC-11 §7.3)"` | **判不动，不改**。`internal/buildinfo/buildinfo.go:25-31` 自陈：key 是票 14 生成的 dev keypair（真身不在我的读数里、被脱敏），production 轮换在 S8；"C29（ModelManifest 签名清单）算不算 landed"需契约考古＋models 链读数，超出本腿文字格射程。 |
| 3 | `cmd/wisp/notify_windows.go:12`（注释） | `// behind it (the resident process's own tray icon is ticket 07/62's, in` | 合法历史出处（过去时归属引来源），不改。 |
| 4 | `cmd/wisp/models.go:19-23`（注释） | `// ... Until ticket 15e this paragraph was honest about the gap ... 15e closed it` | 已自更新（自带作废声明），合法，不改。 |
| 5 | `cmd/wisp/console_other.go:9`、`cmd/wisp/console_windows.go:27` | `// ... not ticket 07) ...` ／ `// ... not ticket 07): such processes` | 合法历史出处（过去时归属），不改。 |
| 6 | `cmd/wisp/panel_host_windows.go:519`（运行时字符串） | `"panel host: embedded bundle not built (only the tracked placeholder is present)"` | 真：`frontend/dist` 只跟踪 `.gitkeep`，bundle 未构建时这句如实，不改。 |

`TODO`／`FIXME`／`XXX`／`待填` 在 `cmd/wisp` 产码＝零命中（同尺）。

## 6. 门禁读数（全部在导出树 `D:\tmp\wisp\20261108\cmdwisp-wording\tree`＝commit `9aaf601c8`）

| 门 | 读数 |
|---|---|
| LF | `tr -dc "\r" < cmd/wisp/main.go | wc -c`＝**0**；`git ls-files --eol cmd/wisp/main.go`＝`i/lf w/lf attr/text eol=lf` |
| gofumpt（绝对路径 `/d/work/base/gopath/bin/gofumpt.exe`） | `-l cmd/wisp/main.go` 空、`-l cmd/wisp/` 全空，rc=0 |
| `go vet ./cmd/wisp/` | rc=0 |
| `go build ./...` | **第一发 rc=1＝假红，与本腿无关**：本腿为 PATH 拷共享树 `build/` 目录时把**别人在飞的未跟踪 scratch `build/mut_rvw.go`**（引用 `voiceTtsFiles`／`residentBall` 等未定义符号）带进了导出树；共享树 `git ls-files build/`＝零枚跟踪 `.go` ⇒ 纯 committed 树无此文件。剥掉外来 `.go`（保留 3 枚 dll）后重跑 **rc=0** |
| `go test ./cmd/wisp -count=1 -run 'TestNothing|TestMain' -timeout 5m` | `ok github.com/CarlosShao/wisp/cmd/wisp 0.207s [no tests to run]`，rc=0（只跑快的，未跑整包 579s） |
| `sh scripts/d22scan.sh` | `clean - no D22 ban violations`，rc=0（ban#8: internal/=965, cmd/=242, frontend/=121, design/=14） |
| 环境 | 每发前 `unset GOTMPDIR TMPDIR TEMP TMP`；`export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"`；⛔ 未设 `CGO_ENABLED=0` |

## 7. 尺（可逐字重跑）

```
sed -n '5,14p' cmd/wisp/main.go
grep -n "runResident\|len(args)" cmd/wisp/main.go
grep -rn "func startResidentBall" cmd/wisp/
grep -n "Exit item" cmd/wisp/main.go
grep -rn "TODO\|FIXME\|XXX\|待填\|placeholder\|ticket 07\|stub" cmd/wisp/*.go | grep -v _test
grep -n "idle set does not include cancel\|ONLY path that registers the cancel id" internal/ball/ball_windows.go
```

## 8. 诚实记账

- 派单两格里**只做成一格**：228 AC#8 落地；245 AC#6① **已被 245-r3（`9e84ae42e`）先行落地**，本腿复认为真、零改动（§4）——非「没做到」，是「已有人做对」。
- 超出 AC#8 判据字面的两处改动（§2a／§2b／§3）都在 `cmd/wisp/main.go` 这一个文件、都是「注释/文案说的和代码做的相反」同族假话、归属都落在台账 `A477` 把 main.go 那族句子判归 228 文案面的裁定内；每句改前改后逐字＋代码事实均在本表。
- 本腿 commit 两枚（产码 `9aaf601c8`；证据件＋票面 log 一枚），**只 commit 不 push**。
