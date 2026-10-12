# 票 295 · 落地腿 `295-r1` · 第 1 笔：`AC#1` 修法落地（`started` 换 `atomic.Bool`）

时刻 `2026-10-12 14:45:20 +0800`。锚＝HEAD `0615c6c9`（派单锚 `73719706` 的直系后代；尺＝`git merge-base --is-ancestor 73719706 HEAD`＝yes，且 `git diff --stat 73719706 HEAD -- cmd/wisp` 输出 0 行 ⇒ 我改的这四枚字节与派单锚逐字节相同）。
定位一律用**内容锚**（ⓓ 清单的逐字原文）。文末"改后行号"＝本腿在自己这一刻现量的，⛔ 当契约、⛔ 进派单，只为本件内部对得上。

---

## 1. 形状：照同一枚 struct 里那枚邻居原子（票面 `AC#1`"照邻居做"）

`levels atomic.Uint64` 一直在那枚 struct 里（ⓓ-1 声明行的下一枚字段）；本改动把 `started` 挪进同一族，**不新造规矩**：

- ⛔ 加锁（无 `sync.Mutex`／`sync.RWMutex` 进入本文件）
- ⛔ 新造协程（无 `go func(`，D38(b) 名册零膨胀；`go build ./...` rc=0 见 `20-gates` 那笔）
- ⛔ 改默认值／D43 转移表／门的语义（`gate.Open()`/`gate.Muted()` 仍是主，这枚是派生镜像）
- ⛔ ⓓ-4 之外的任何换形状：`verdict`（`AC#0` 前件已由 `295-a1b` 复跑判不成立）、`rb.b`、`gate`／`mic`／`cancel`／`voiceEnabled`／`mutedAtBoot` **一字未动**

尺（复跑）＝`git grep -n -e 'started atomic.Bool' -e '.started.Store(' -e '.started.Load()' -- cmd/wisp` ⇒ **12 行**＝ⓓ-1 的 4 处（含声明）＋ⓓ-2 的 8 枚读点，**多一枚少一枚都算没对上**。
残形尺（应为 0）＝`git grep -n '\.started =' -- cmd/wisp` ⇒ **0 命中**；`git grep -nE '\.started([^.(]|$)' -- cmd/wisp` ⇒ **0 命中**。

---

## 2. 逐处对照（内容锚逐字 ⇄ 改后逐字）

### ⓓ-1 生产面 `cmd/wisp/resident_audio_windows.go`（4 处含声明）

| 处 | 内容锚（改前逐字） | 改后逐字 | 改后行号（现量） |
|---|---|---|---|
| 声明 | `	started bool` | `	started atomic.Bool` | `:116` |
| 写 1（`toggleMute` 体内，★票的正主＝跨线程写） | `	ra.started = gate.Open()` | `	ra.started.Store(gate.Open())` | `:187` |
| 写 2（`assembleCapture` 的 `default:` 支，boot 线程） | `		ra.started = true` | `		ra.started.Store(true)` | `:342` |
| 读 1（`posture()` 体内，生产唯一读者） | `	if !ra.started {` | `	if !ra.started.Load() {` | `:398` |

`sync/atomic` 本来就在 import 块里（邻居 `levels` 用的就是它），**未新增 import**。

### ⓓ-2 测试面（8 枚读点，条件位／参数位逐枚对）

| # | 件 | 内容锚（改前逐字） | 改后逐字 | 改后行号 |
|---|---|---|---|---|
| 1 | `resident_audio_247_windows_test.go` | `	if ra.started {`（后一行 `t.Fatal("a leg with no collector reports itself running")`） | `	if ra.started.Load() {` | `:113` |
| 2 | 同上件 | `	if ra.started {`（后一行 `t.Fatal("posture claims a running capture leg with the gate muted")`） | `	if ra.started.Load() {` | `:155` |
| 3 | 同上件（表驱动子测试，缩进 3 层） | `			if ra.started {` | `			if ra.started.Load() {` | `:256` |
| 4 | `resident_audio_247_live_windows_test.go` | `	if !ra.started {` | `	if !ra.started.Load() {` | `:148` |
| 5 | `resident_mute_290_windows_test.go` ★全场唯一一枚咬"手势路径必须写 started" | `		if !ra.started {` | `		if !ra.started.Load() {` | `:124` |
| 6 | 同上件、同一 for 轮（`\|\|` **右侧**） | `		if ra.gate.Open() \|\| ra.started {` | `		if ra.gate.Open() || ra.started.Load() {` | `:141` |
| 7 | 同上件 ★`t.Fatalf` **格式串参数位**（⛔ 条件位） | `			t.Fatalf("%s: capture stayed handed over after muting (open=%v started=%v)", gesture, ra.gate.Open(), ra.started)` | `			t.Fatalf("%s: capture stayed handed over after muting (open=%t started=%t)", gesture, ra.gate.Open(), ra.started.Load())` | `:142` |
| 8 | 同上件（`TestAC290OutcomeIsReadOffTheGateNotOffTheRequest`） | `	if ra.started {` | `	if ra.started.Load() {` | `:174` |

**#7 单独交代（派单 `AC#2` ⓓ 那一格）**：改前把 `atomic.Bool` 整个丢进 `%v` 会打出 `{true}`／`{false}` 的结构体形状——编译过、断言不红，但失败日志开始撒谎，而这枚用例的用途恰恰是把 `started` 写进操作员可见的那一行。⇒ 实参已改 `.Load()`，两个 `%v` 同时**收成 `%t`**（更严：非布尔实参在 `go vet`／`go test` 的内建 vet 下就会响）。
那行的输出形状（本腿实测，变异体那一发里逐字读到，见 `20-anti-form.md`）＝`capture stayed handed over after muting (open=false started=true)`：两枚都是裸 `true`／`false`，**零枚花括号**。

---

## 3. 注释那两组：措辞说"派生镜像"，再补"为什么这一枚必须原子"

ⓓ-1 已具名这两组（声明上方 3 行、`toggleMute` 写点上方 4 行）。原句**一字未删**，只在各组末尾追加一段；追加段落里那句"必须原子"的落点：

- **声明组**（原首行 `	// started records that the gate handed the device to the inner source and`）
  新段落 10 行，最后一行逐字＝`	// source (ticket 295 AC#1 "照邻居做").`，它**紧贴在 `	started atomic.Bool` 那行（现量 `:116`）的上一行**。
  句内点名的三件事：写者＝`toggleMute`（球的手势回调，ui-sta 线程，逐字 `// ball's global hotkeys on the ui-sta thread` 在 `resident_ball_windows.go`）、读者＝`posture()`（boot 线程）、⛔ 仪器看得见（`-race` 不在门禁名册＝票面现量末条）。
- **写点组**（原首行 `	// started mirrors what the gate now reports. It is derived state, not a`）
  新段落 8 行，最后一行逐字＝`	// their mirror, stored atomically.`，它**紧贴在 `	ra.started.Store(gate.Open())` 那行（现量 `:187`）的上一行**。
  这一组还写死"⛔ 自持一份静音状态"：`gate.Open()/gate.Muted() stay the truth; this is only their mirror`（票面禁区①／票 293 禁区同源）。

尺（复跑这两句的位置）＝`git grep -n -e 'source (ticket 295 AC#1' -e 'their mirror, stored atomically' -- cmd/wisp/resident_audio_windows.go` ⇒ 2 命中，分别站在 `:116` 与 `:187` 的上一行。

---

## 4. 本格的凭据读数（不是"加了 atomic"本身）

- `AC#1` 只负责形状；**牙在 `AC#2`**：删写变异那一发指名红在 `:124` 那枚断言（逐字红句、还原证据、7 枚不敏感逐名＝见 `20-anti-form.md`）。
- 改后 focused 两发全绿（`21-ac2-pristine-named.txt`）：指名用例 `--- PASS: TestAC290BothMuteGesturesTurnTheGateAndBack (0.01s)`；`-run 'TestAC247|TestAC290|TestLive' ./cmd/wisp/` 名册＝**顶层 13 枚 PASS ＋ 1 枚 SKIP ＋ 0 枚 FAIL**（同一发第二把口径含子项＝16 枚 PASS；两把⛔ 相加当一枚数）。SKIP 那枚＝`TestAC247LiveMicrophoneLevelsReachTheBallSeam`＝真机件，改前两发同样 SKIP，⛔ 我这一发改出来的。
- ⛔ 放宽任何断言：`git diff` 里 4 枚件的 `-` 行 16 枚／`+` 行 35 枚（尺＝`grep -c '^-'`、`grep -c '^+'` 对 `22-ac1-diff.txt`），逐枚可核＝**删的全是被换掉的旧形状与旧 `%v`，无一枚断言被改宽**（diff 全文落 `logs/22-ac1-diff.txt`，7,395 字节）。
