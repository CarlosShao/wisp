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
| 声明 | `	started bool` | `	started atomic.Bool` | `:105` |
| 写 1（`toggleMute` 体内，★票的正主＝跨线程写） | `	ra.started = gate.Open()` | `	ra.started.Store(gate.Open())` | `:168` |
| 写 2（`assembleCapture` 的 `default:` 支，boot 线程） | `		ra.started = true` | `		ra.started.Store(true)` | `:323` |
| 读 1（`posture()` 体内，生产唯一读者） | `	if !ra.started {` | `	if !ra.started.Load() {` | `:379` |

★这四枚行号＝票面"现量"那四枚（`:105`／`:168`／`:323`／`:379`）**一字未动**——不是巧合，是第 5 节那条硬约束的结果（本件必须净零增行，否则票 255 的回执行证会漂）。

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

## 3. 注释那两组：措辞说"派生镜像"，再补"为什么这一枚必须原子"——而且**回到原行数**

ⓓ-1 已具名这两组（声明上方 3 行、`toggleMute` 写点上方 4 行）。第 5 节那条约束（票 255 的行证按**行号**引本件）⇒ 这两组的新措辞必须**各自还是 3 行／4 行**，⛔ 追加段落。两组都保留了原句的全部claim（"记录门把设备交给了内层 source、采集真在跑"／"posture() 读它，所以建了门却没开过麦克风的腿不能打印 running"／"镜像、不是第二真相源"／"保持最新才不会有后来的读句谎称设备"），并在组末补上为什么这一枚必须原子：

- **声明组**（现量 `:102-104`，紧贴 `	started atomic.Bool`（`:105`）的上一行是它自己的第三行）
  末行逐字＝`	// never opened a mic cannot print "running". ui-sta writes, boot reads: 295.`
  ⇒ "必须原子"那一枚落点＝**`:104`**，站在 `:105` 的上一行；句里点名的三件事＝写者线程（`ui-sta`）／读者（boot 的 `posture()`）／票号 295。
- **写点组**（现量 `:164-167`，紧贴 `	ra.started.Store(gate.Open())`（`:168`））
  末行逐字＝`	// It crosses threads (this body is ui-sta, boot reads it), so it is a Store.`
  ⇒ "必须原子"那一枚落点＝**`:167`**，站在 `:168` 的上一行；同组第一行仍写死"derived state, not a second truth source"（票面禁区①／票 293 禁区同源：`gate.Open()`/`gate.Muted()` 是主，这枚只是镜像）。

尺（复跑这两句的位置）＝`git grep -n -e 'ui-sta writes, boot reads: 295' -e 'so it is a Store' -- cmd/wisp/resident_audio_windows.go` ⇒ 2 命中＝`:104` 与 `:167`。

---

## 4. 本格的凭据读数（不是"加了 atomic"本身）

- `AC#1` 只负责形状；**牙在 `AC#2`**：删写变异那一发指名红在 `:124` 那枚分支守的断言上（`t.Fatalf` 打在 `:125`，逐字红句、还原证据、7 枚不敏感逐名＝见 `20-anti-form.md`）。
- 改后 focused 两发全绿（`21-ac2-pristine-named.txt`）：指名用例 `--- PASS: TestAC290BothMuteGesturesTurnTheGateAndBack (0.01s)`；`-run 'TestAC247|TestAC290|TestLive' ./cmd/wisp/` 名册＝**顶层 13 枚 PASS ＋ 1 枚 SKIP ＋ 0 枚 FAIL**（同一发第二把口径含子项＝16 枚 PASS；两把⛔ 相加当一枚数）。SKIP 那枚＝`TestAC247LiveMicrophoneLevelsReachTheBallSeam`＝真机件，改前两发同样 SKIP，⛔ 我这一发改出来的。
- ⛔ 放宽任何断言：`git diff` 里 4 枚件的 `-` 行 16 枚／`+` 行 35 枚（尺＝`grep -c '^-'`、`grep -c '^+'` 对 `22-ac1-diff.txt`），逐枚可核＝**删的全是被换掉的旧形状与旧 `%v`，无一枚断言被改宽**（diff 全文落 `logs/22-ac1-diff.txt`，7,395 字节）。

---

## 5. ★本格踩到的一枚真回归（已复绿）：本件的**行数**是票 255 回执的引据

`d5ef9cb9`（AC#1 第一版）把两组注释各追加了一段（净 ＋18 行）⇒ 本件 `:256`／`:257` 两行下移 ⇒
`cmd/wisp/config_readers_255.go` 的 `hotRowClaims` 里 5 枚按行号引本件的行证失配 ⇒
`TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim` 由改前两发的 PASS 翻成 FAIL＝**本腿造成的新增红 1 枚**。
本件底部自己的注释早就写着这条规矩（逐字）：`// of this file on purpose - ticket 255's config receipt cites lines inside this file by number, so appending here is what keeps those rows from drifting.`
处置＝**两组注释压回原行数**（本件净零增行，`git diff --numstat 0615c6c9` ⇒ `10	10`），⛔ 动 `config_readers_255.go`／`config_receipt_255_test.go`（⛔ 本票写面），⛔ 让回执迁就漂移（＝放宽断言）。
复绿读数＝`logs/41-driftfix-verify.txt`：那枚用例逐字 `--- PASS: TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim (0.00s)`，整发 `ok github.com/CarlosShao/wisp/cmd/wisp 8.766s`／`DRIFTFIX_RC=0`（33 枚顶层 PASS＋1 枚 SKIP）。
⇒ 这一枚已进 `90-final.md` 的欠账名册当**已销**条（成因、尺、复绿三件齐），验收腿要现读那三件。
