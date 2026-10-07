# 273-a2 — 机器实例普查（只读腿，零产码改动）

票面：`.scratch/wisp/issues/273-shipping-process-builds-the-state-machine-without-a-sink-so-every-d43-side-effect-falls-into-a-no-op.md`（整份已读，60 行）。
本腿射程＝§1 机器实例点名 / §2 逐枚副作用名追生产者 / §3 三层条数原文 / §4 落点三事实。**不裁决、不选形、不写产码、不跑任何 `go test`。**

---

## §0 起手锚

```
$ git rev-parse --short HEAD
15699a2f

$ git status --porcelain -- cmd internal docs
(空)

$ date
Wed Oct  7 10:24:58 CST 2026
```

---

## §1 全仓"机器实例"逐枚点名

### 1.1 先把尺定下来（三把，互相补漏）

**尺 A（构造形状全量，含测试与 `.scratch`）**
```
$ grep -rn 'statemachine\.New(' --include=*.go .
```
原始末 3 行（不截断，总 8 枚命中）：
```
./internal/models/bridge_test.go:16:	machine := statemachine.New(statemachine.Options{
./internal/models/bridge_test.go:86:	machine := statemachine.New(statemachine.Options{Initial: statemachine.StateSleeping})
./internal/models/handoff_window_109_test.go:62:	return statemachine.New(statemachine.Options{Initial: statemachine.StateFirstRun})
RC=0
```
命中 8 枚＝产码 2 ＋ `_test.go` 6（`.scratch` 里那些 statemachine 变异拷贝**贡献 0 枚 New**，它们只 import）。

**尺 B（同名噪声防漏：包别名写法穷举）**——先确认 import 有没有别名，再确认 `New` 只有限定名一种写法：
```
$ grep -rhn 'wisp/internal/statemachine"' --include=*.go . | sort | uniq -c
```
原始末 3 行：
```
      1 8:	"github.com/CarlosShao/wisp/internal/statemachine"
      1 9:	"github.com/CarlosShao/wisp/internal/statemachine"
      3 3:import "github.com/CarlosShao/wisp/internal/statemachine"
```
⇒ 全仓该包 import **一律无限定别名**，所以 `statemachine.New(` 就是唯一的构造调用形状。

**尺 C（绕过 `New` 的直接复合字面量——同族写法的第二形）**
```
$ grep -rn 'statemachine\.Machine{' --include=*.go .
RC=1        (零命中)
```
⇒ 没有任何"不经 `New` 手搓 Machine"的形状；`sink` 字段是 `Machine` 的私有字段（`machine.go:52`），包外也写不到。

**尺 D（`Sink` 那枚类型的名字，防"同名噪声"给假读数）**
```
$ grep -rn 'statemachine\.Sink' --include=*.go .
RC=1        (零命中)
```
⇒ 全仓（含测试）**没有任何一处写出 `statemachine.Sink` 这个类型名**。判"有没有人接"因此**不能**靠 `grep Sink`：仓里叫 Sink 的三枚都不是它——
`grep -rn 'Sink:' --include=*.go cmd internal | grep -v _test` 的产码命中是 `cmd/wisp/providers.go:194`（`storeHealthSink`＝`llm.HealthSink`）与 `cmd/wisp/run.go:998`（`consoleSink`＝`agent.Sink`），**两枚都不是状态机的 Sink**。

**尺 E（"哪枚构造点真带了 Sink 字段"＝票 AC#1 要的尺）**——New 点向后看 4 行取 `Sink:`：
```
$ grep -rn -A4 'statemachine\.New(' --include=*.go cmd internal | grep 'Sink:'
internal/models/bridge_test.go-18-		Sink:    func(e statemachine.Effect) { effects = append(effects, e) },
RC=0
```
⇒ **8 枚构造点里只有 1 枚带 Sink，且那 1 枚在测试里**。产码 2/2 不带。这条与票面 [10-07 09:5x] 更正一/更正二**同数**（测试 6 枚、其中 1 枚真传），我复跑独立得到。

### 1.2 逐枚点名

| # | 位置（`file:line`） | ⓐ 落在哪个函数／入口 | 那条链今天真有人调吗（调用点枚数） | ⓑ 填了 `Options` 哪几个字段 | ⓒ 活多久 | ⓓ 有没有 Sink |
|---|---|---|---|---|---|---|
| 1 | `cmd/wisp/models.go:303` | `(*modelStore).handOffModel`（`cmd/wisp`，出货二进制里那台） | 生产链 3 跳，逐跳 **1 枚**真调用点：`handOffModel` ← `models.go:293`（`modelsEnsure`）；`modelsEnsure` ← `models.go:127`（`cmdModels` 的 `case "ensure"`）；`cmdModels` ← `cmd/wisp/main.go:109`（`case "models"`，`func main`）。⇒ 只有跑 `wisp models ensure <id>` 这一条 CLI 腿会造它 | 只 `Initial: StateFirstRun`（`Sink`/`Timeouts` 未填） | **一次 CLI 调用**：`:304` `defer machine.Close()`，函数返回即死；不是进程生命周期 | **无** ⇒ `machine.go:64→:67` 缺省 no-op |
| 2 | `cmd/balldebug/main.go:188` | `func main()`（`cmd/balldebug`，开发期调试器） | `main` 即入口，无上游；下游机器被 `gesture`(`:591`)／`dispatch`(`:617`)／`syncMachineToBall`(`:630`) 使用，这三个由 `ball.Events` 回调（`:194-208`）驱动。⚠ **`balldebug` 不是出货件**：`scripts/build.ps1:129` 只 `go build ... ./cmd/wisp`；`balldebug` 仅 `scripts/dev/ball-cycle.ps1:30` 构建 | 只 `Initial: StateSleeping` | **一次 balldebug 进程**（dev 手跑，随进程退出） | **无** ⇒ 同样 no-op |
| 3 | `internal/ball/hotkey_live_test.go:382` | `TestLiveMuteHotkeyEndToEnd` | 测试；`winlive` 门控用例 | 只 `Initial` | 一次用例 | **无** |
| 4 | `internal/ball/interaction_live_test.go:54` | `TestLiveClickSummonsAndDragDoesNot` | 测试；`winlive` 门控用例 | 只 `Initial` | 一次用例 | **无** |
| 5 | `internal/ball/interaction_live_test.go:134` | `TestLiveConfirmingCancelAndEscReturned` | 测试；`winlive` 门控用例 | 只 `Initial` | 一次用例 | **无** |
| 6 | `internal/models/bridge_test.go:16`（跨 `:16-19`） | `newWalkRig` 测试夹具 | 测试夹具（`grep -c newWalkRig` 见下） | **`Initial` + `Sink`**（唯一一枚） | 一次用例 | **有** ⇒ 全仓唯一真收件人，把 `Effect` append 进切片 |
| 7 | `internal/models/bridge_test.go:86` | 同文件另一枚（`Initial: StateSleeping` 的负控形状） | 测试 | 只 `Initial` | 一次用例 | **无** |
| 8 | `internal/models/handoff_window_109_test.go:62` | `newWalkMachine()` 夹具 | 测试夹具 | 只 `Initial` | 一次用例 | **无** |

**小结（我这把尺的数，不是引 a1 的）**：机器实例 **8 枚**＝产码 **2** ＋测试 **6**；带 Sink 的 **1** 枚（测试 #6）；产码 **0** 枚带 Sink。
⇒ 出货进程（`wisp.exe`）里**只有 `wisp models ensure` 这一条腿会造出一台机器**；`wisp run`（`cmd/wisp/run.go`）与常驻腿（`runResident`，`resident_windows.go:33`／`resident_other.go:22`）**一台都没有**——常驻腿只把 `statemachine.State*` 当常量用（`resident_ball_windows.go:275` 是 `ball.Options{Initial: ...}`，不是 Machine；`resident_approval_windows.go:956` 是 `b.SetState(...)`）。

**"某符号零调用者"锚在调用形状**：上表的调用点枚数一律锚成 `标识符(` 或 `x.method(` 形状（`grep -rn 'cmdModels'`／`'modelsEnsure'`／`'handOffModel'` 的命中里，注释行与声明行我都逐条分辨后剔除；例如 `handOffModel` 的 7 枚文本命中里只有 `models.go:293` 是调用，其余 6 枚是声明与注释）。
