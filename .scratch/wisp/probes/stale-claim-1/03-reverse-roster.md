# stale-claim-1 / 03 — 表三：反向名册（"已经接线／已在生产／唯一调用者"这一族）

同一把尺（`00` §1.2 调用形状），同一个 ref `HEAD 6547fd30`，同一套三层（包内／包外／到得了入口）。
正向断言过期时**同样骗人，只是方向相反**：它骗的不是"别派这活"，而是"**别派那活 / 这格已经安全了**"。

射程：产码非测试文件里，本程判为**正向状态断言**的 14 枚。
分档计数：**仍成立 10（含一枚"仍成立但脆"）／ 已过期 1 ／ 判不动 3**。

## 甲 · 仍成立（9 枚）——这批可以放心当现状读

| # | HEAD 位置（逐字） | 断言 | 尺 -> 输出摘要 | 三层 |
|---|---|---|---|---|
| R01 | `HEAD:cmd/wisp/approval_always.go:158` `// bookWaitingState is ticket 201 AC#6's production caller: every time this run` | `(*agentRuntime).bookWaitingState` 有生产调用者 | `git grep -nE "bookWaitingState" HEAD -- ':(exclude).scratch' '*.go'` -> 定义 `:171`；非测试调用点 **2**：`cmd/wisp/run.go:1361 u.run.bookWaitingState("ui-prompt")`、`run.go:1407 u.run.bookWaitingState("ui-" + string(e.Kind))` | 包内(main) 有 ／ 入口 `runTextTask`（`run.go:181`）→ console UI 回调 **到得了** |
| R02 | `HEAD:cmd/wisp/config_readers_255.go:235` `if tier, ok := config.TierOf(name); ok { // TierOf's first production caller` | `config.TierOf` 的第一枚（＝唯一一枚）生产调用者 | `TierOf(` 非测试命中：本行（调用）+ `internal/config/tiers.go:93`（定义）。"first"＝"only" **成立** | 包外(cmd/wisp -> internal/config) 有 ／ 入口：本行在 `hotRowsFor`（`:234`）体内，`hotRowsFor` 被同文件 `:266`、`:292` 调；更上游未追 -> `04` B-06 |
| R03 | `HEAD:cmd/wisp/firstrun.go:27-28` `// runTextTask - the single production caller chain main.go cmdRun ->` `// runTextTask.` | **`single`**（唯一入口链） | `ensureFirstRunConfig`（定义 `firstrun.go:72`）非测试调用点 **恰 1** = `run.go:245`；入口链 3 枚全中：`main.go:91 os.Exit(cmdRun(args[1:]))`、`main.go:151 func cmdRun`、`main.go:159 return runTextTask(runSpec{`。**"single" 这把尺我量得动，量到的是 1。** | 包内 有 ／ 入口 **到得了 `main`** ／ ⚠ 函数体归属由 `awk 'NR>=176 && NR<=260 && /^func /'` 钉住 = 176–260 行里唯一的 `func` 是 `:181 runTextTask`，故 `:245` 确在其体内 |
| R04 | `HEAD:cmd/wisp/approval_reply.go:9-11` `// structural fact and it is the whole subject of this file: the ASKING side is` `// live in production - internal/tools/bridge.go calls Gate.PendingWindow for an` `// L1 verdict and Gate.PendingApproval for an L2 one, and cmd/wisp's own mode` | `PendingWindow`／`PendingApproval` 在产发热路径上活着 | `git grep -nE "\.PendingWindow\(\|\.PendingApproval\("` -> 非测试：`internal/tools/bridge.go:443 a, why := b.gate.PendingWindow(ctx, *dec)`、`bridge.go:458 a, why := b.gate.PendingApproval(ctx, *dec)`（正是注释点名的那两处），另有 `cmd/wisp/run.go:844`、`cmd/wisp/approval_always.go:115`、`cmd/wisp/config_reload.go:279`、`cmd/wisp/resident_approval_windows.go:678/:680` | 包外 有 ／ 入口 两条腿（`wisp run` 与常驻）都到 ／ **这枚是表一 P12 的正向半边：它没说错** |
| R05 | `HEAD:internal/tools/task.go:685-686` `// had and nobody called (ticket 221's whole complaint). This tool is its first` `// production caller, and it is the ONLY one: no second stop path gets invented.` | **`ONLY one`**（排他断言，最容易骗人的一种） | 尺：`\.Cancel\(` 非测试全谱剥同形异符号（`observe.Root`/context 根 15 枚）-> `TaskRoster.Cancel` 调用点 **恰 1** = `internal/tools/task.go:740 stopped, why := t.d.Roster.Cancel(target)`；定义 `task.go:447` | 包内 有 ／ 注册链到得了入口：`run.go:544 tools.BuiltinTaskEntries(tools.TaskDeps{Roster: rt.tasks, Paths: rt.paths})` -> `task.go:600 BuiltinTaskEntries` -> `taskCancel.Execute` ／ **"ONLY" 成立〔已证〕** |
| R06 | `HEAD:internal/agent/approval/approval.go:535-536` `// What the bind argument is on a routed call: the only production caller of this` `// function is Queue.allowScoped, which passes the item's own it.bind - the same` | `spend` 的唯一生产调用者 | `allowScoped` 非测试：调用形状 **恰 1** = `internal/agent/approval/gate.go:670 if err := g.q.allowScoped(corr, grant, true); err != nil`；其余 11 枚命中全是注释／定义（`approval.go:359/366/373/415/441/470/496/536`、`queue.go:50/58/343`） | 包内 有 ／ **排他成立〔已证〕** |
| R07 | `HEAD:internal/config/manager.go:23` `// ConfirmLocked (ticket 223). That sentence used to be a lie for ConfirmLocked:` | 该 hook 过去是死线、现在不是（正向） | 同 P08 尺：赋值 `cmd/wisp/config_reload.go:114`、读侧 `manager.go:172` | 包外 有 ／ 入口 `run.go:813` **到得了** ／ 仍成立 |
| R08 | `HEAD:internal/tools/mode.go:53` `// BlacklistNote is risk.Gate's own reading of one call's paths, kept beside the` ＋ `:53(AC#5 段) // AC#5 asked for risk.Gate to stop being a function with zero production call` | `risk.Gate` 已脱离"零生产调用者" | `git grep -nE "risk\.Gate\("` -> 非测试调用点 **恰 1** = `internal/tools/mode.go:98 d := risk.Gate(c, overrides)` | 包外(tools -> risk) 有 ／ 正向半边成立；⚠ **`overrides` 能不能非空**＝表一 P29 判不动那枚，两者不许压成一枚 |
| R09 | `HEAD:cmd/wisp/config_reload.go:116` `// AC#1's production caller. Owner names the subsystem that spawned it; the` | 该 tick 的生产调用者已存在 | 尺：`startConfigReload` 非测试 2 枚（定义 `:105` + 调用 `run.go:813`），见 Q09 | 包内 有 ／ 入口 `runTextTask` **到得了** ／ 仍成立（常驻腿不到，见 Q09） |

## 乙 · 已过期 1 枚（R10）＋ 仍成立但脆 1 枚（R11）——正向里"过期方向相反"的那两枚

### R10 `HEAD:cmd/wisp/resident_approval_windows.go:13-15`
> `// (ledger A480 ②/A481 ③): TakeEscForCancel/ReleaseEscAfterSession had zero`
> `// callers in any product process, and cmd/wisp/resident_ball_windows.go said so`
> `// about itself in three places. So D43's veto row had no executor in the process`

- 断言形状：**负向**（"had zero in any product process"），但它的**语义是进度**——同一句在讲"这枚由本文件补掉了"。
- 尺：`git grep -nE "TakeEscForCancel|ReleaseEscAfterSession" HEAD -- ':(exclude).scratch' ':(exclude)*_test.go'`
- 输出摘要（非测试调用形状）：`cmd/wisp/resident_approval_windows.go:877 b.TakeEscForCancel()`、`cmd/wisp/resident_approval_windows.go:955 b.ReleaseEscAfterSession()`、`cmd/balldebug/main.go:635 b.TakeEscForCancel()`、`cmd/balldebug/main.go:638 b.ReleaseEscAfterSession()`。定义：`internal/ball/ball_windows.go:881`、`:911`。
- **分档：已过期（作为现状读）〔已证〕。现在到得了它的是 `HEAD:cmd/wisp/resident_approval_windows.go:877` 与 `:955`**（同文件内，`resident_windows.go:217 startResidentBall(...)` 那条常驻入口链）。
- ⚠ 害法：这枚过期时，读它的人会以为**"球侧的取消键借还仍然没有执行者"**，从而把票 245/246/260 那一族已落的活重开。它与表一 P20（`no caller does that today`，同一件事的另一半）**指向同一枚未裁读数**。

### R11 `HEAD:cmd/wisp/run.go:532-535`
> `// task.cancel's production landing point (ticket 221 AC#2: TaskRoster's Cancel`
> `// had zero production callers until the tool that reads it got registered -`
> `// count it with a grep for calls of that method outside _test, which answers`
> `// exactly one hit now, internal/tools/task.go's taskCancel.Execute; the`
> `// description of the measure is spelled out here without the call syntax on`
> `// purpose, so this comment can never be miscounted as a caller).`

- **这枚是整格普查里最好的一枚，所以它进的是"已过期"档、进的是好事。**
- 尺（照它教的办法复跑）：非测试里 `TaskRoster.Cancel` 调用形状 = `internal/tools/task.go:740`，**恰 1** ⇒ **`exactly one hit now` 今天仍成立**。
- **分档：仍成立（但脆）**，理由见下一条。
- 我把它单拎出来**不是**因为它错了，而是因为：它说 `until the tool that reads it got registered` —— 而"registered"这一枚**它没给尺**（注册表在 `BuiltinTaskEntries`，`run.go:544`）。我把注册链补量了（`task.go:600 BuiltinTaskEntries` -> `run.go:544`），结论仍成立。**记它"脆"是因为它把"exactly one"钉在了一个会随 HEAD 动的计数上** —— 下一次谁再补一枚 stop path，这句就得同时改 `run.go:532-535` 和 `task.go:686`（两枚排他断言，一处一动必漂另一处）。
- ⇒ **给后续程的一句话事实（不是修法）**：本仓**排他式正向断言（"ONLY one"／"exactly one hit"）在两枚不同文件里钉同一枚事实**：`internal/tools/task.go:686` 与 `cmd/wisp/run.go:534`。它们今天**同真**，一动就分家。

## 丙 · 判不动（3 枚）

| # | HEAD 位置（逐字） | 为什么判不动 | 缺的读数 |
|---|---|---|---|
| R12 | `HEAD:cmd/wisp/run.go:341-344` `// pump assembles the packet the panel is rendered from, out of the live` `// objects this struct already holds (ticket 35's data half). It is the first` `// thing in this repository that calls panel.NewSnapshot /` `// panel.NewComposerState from a running process rather than from a test -` | "first … from a running process rather than from a test"是**排他＋历史双重断言**，需要全仓调用点分类（测试／非测试）＋"running process"的归属判断。我这把尺能给"非测试调用点集合"，⛔ 不能证明"first"。 | `git grep -nE "NewSnapshot\(|NewComposerState\(" HEAD -- ':(exclude).scratch' '*.go'` 的全谱＋逐枚分类；见 `04` B-09 |
| R13 | `HEAD:tools/d22scan/gitignore.go:297` `// moves in between (ticket 142 AC#1). runGitIndex is its only production caller.` | **"its" 的先行词没在句里**（`:294` 讲的是另一枚函数的形状）。`runGitIndex` **自身**的非测试调用点是 **恰 1**（`gitignore.go:510 g.idx = runGitIndex(g.root)`；定义 `:261`）⇒ 若"its"指 `runGitIndex`，仍成立；若指 `:294-297` 那块描述的那枚 helper，**我没读出 helper 的名字**。 | 读 `gitignore.go:286-300` 整块定先行词；见 `04` B-10 |
| R14 | `HEAD:cmd/wisp/approval_always.go:45` `// tools.NewPathCanonicalizer mid-run), and it is reported as such.` ＋ `HEAD:internal/tools/bridge.go:87` `// which is the fail-closed direction - an unwired composition can never run` ＋ `HEAD:internal/tools/mode.go:14` `//     (nil = the strictest mode, so an unwired composition cannot accidentally` | 这三枚是**行为分支断言**（"没接进来会怎样"），不是状态断言。按尺它们**永远不过期**（它们描述的是 nil 分支的语义），⛔ 也就⛔ 不能拿来当"当前接没接"的证据。 | 无 —— 本程把它们**显式剔出名册**，只列在这里防止下一程把它们误计为"已过期"或"仍成立"。 |

---

## 表三·补 · 正向断言为什么比负向更难验（本程实测到的三件事）

1. **`first` / `only` / `single` 是三把不同的尺**：`R03` 的 `single` 我用"调用点计数=1"验掉了；`R02` 的 `first` 靠"当前只有 1 枚"推出来；`R12` 的 `first` **验不了**（它是历史序，不是当前集）。
2. **正向断言会随时间自己变真，但不会自己变假**——除了新增调用者（把 `only` 打成假）与删掉调用者（把 `is live` 打成假）。⇒ 只有 `only/single/exactly one` 这一型会过期；纯存在型（"有生产调用者"）**只会因为删码而过期**，本程的尺对此有盲区（`git log -S` 才能量，我没跑）。
3. **排他断言成对出现时必漂其一**（`R11` 那枚实例）：`internal/tools/task.go:686` 与 `cmd/wisp/run.go:534`。
