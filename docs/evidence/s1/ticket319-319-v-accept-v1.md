# ticket 319 · 319-v 非实现者对抗验收（v1）

- 角色：验证腿 **319-v**，非实现者（实现者＝死腿 319-r1（编排者代提 `e6d1373e`）＋续腿 319-r2（`2ad0aa86`））。
  本腿零产码：不改任何跟踪的 Go 文件、不动 `Status:`/票名/AC 框、不碰台账/HANDOVER。
- 锚：起步 `33e429c1`；腿中 HEAD 被别家推进（`592a69cd` → `90a68597` → `a920c4d4`），
  **319 面在这四枚之间逐字节不变**（`git diff --stat 33e429c1 a920c4d4 -- <4 生产件＋3 319 测试件>` = 0 行；
  四枚生产件 sha256 与 `git show HEAD:<f>` 相等）。所有读数都是对**已落地的 319 码**的读数。
- 方法：只读 + **仓外**仪器。全部变异用 `go test -overlay <json>` 注入（Replace 指向
  `C:\Users\swq\AppData\Local\Temp\wisp319v-mut\` 的副本）；每发之后核 sha256＝HEAD、
  `git diff -- <件> | wc -c` = 0、`git status --porcelain -- cmd/wisp/ internal/proc/` = 0 行。
- 两枚探针（热路径形状 A/B、真两进程）也是 overlay 注入、**从不落进仓**（源码见附录 A）。
- 环境：`GOTMPDIR=C:\Users\swq\AppData\Local\wisp-test-gotmp\wisp319vgotmp`；
  `PATH` 前置 `third_party/sherpa-onnx`。机主实例 `wisp.exe` PID 40296 全程在场、未扰
  （每发活体捕获都带 `tasklist` 自证行）；winlive **一次一枚**、一枚一进程。
- 读数件：`.scratch/wisp/probes/319/v/`（00–05，6 件）。以下所有引文均为**我自己的**跑件逐字，
  ⛔ 没有把 319-r1/r2 的日志当读数引用。

## 1. 逐格判语（先行）

| 格 | 判语 | 一句话理由 |
|---|---|---|
| AC#0 根因（既证，本腿只核） | **成立** | census 的关窗致死链与我在活体上半段逐字重放一致（真窗 `WM_SYSCOMMAND/SC_CLOSE` ⇒ 线程 `finished`）。 |
| AC#1 三态（headless，N1 持有者半） | **成立** | 自己造的变异（摘 `isFinished()` 支）⇒ 用例红名点到"重建"与"spawn 第 2 次被问"；现树绿；还原绿＋`git diff` 0 字节。 |
| AC#1 活体（N3 两半） | **成立** | 真窗关窗后死线拒单＋无窗；同一手势经 holder ⇒ **新 HWND ≠ 旧**、`mgr.IsCreated/IsShown=true`；同变异活体红名逐字、还原绿。 |
| AC#2 活体唤窗 | **附条件成立** | 生产线在三种形状（冷启动／已有窗但隐藏／**真第二进程**）都真出窗＋真夺前台；两处具名残余：① 面板**已可见**时二次启动把它**藏起来**（divergence D1）② ASFW 的必要性未被隔离证明（divergence D2）。且仓内正式用例仍只有同进程形，真两进程那层是**本腿的仓外探针**，不是常驻钉子。 |
| AC#3 钉子与变异 | **成立** | 三钉逐钉自造变异：N1（重建＋spawn 第二次被问）、N1 拒单半（headless＋live 两处）、N2（hook 恰一次＋Reset 两发变异）、N3（活体两半）；全部"弄坏⇒红名逐字⇒还原⇒绿"。 |
| AC#4 门禁＋面核 | **成立** | 四门在 319 面全绿（见 §7 关于一次 `vet` rc=1 的别家归因）；两笔落地只落在 319 面＋票面＋probes；冻结点名零命中。 |
| 回归（改前既有 4 枚） | **成立** | N1 原始半／N2 两枚／N3 原始半照旧全绿；`-run '319' ./cmd/wisp/ ./internal/proc/ -count=1` 绿。 |

## 2. AC#1 三态（headless）——读数逐字

变异（仓外副本 + overlay）：`cmd/wisp/resident_windows.go:535`
`if h.rp == nil || h.rp.isFinished() {` → `if h.rp == nil {`；
`grep -c "h.rp.isFinished()" <mutant>` = 0；`go build -overlay ... ./cmd/wisp/` rc=0。

- **① 未修形（变异）**：`FAIL: TestPanelHolderRebuildsADeadThread319`，红名逐字
  - `panel_holder_319_windows_test.go:146: panelHolder.toggle returned false ... a dead thread must make the holder start a fresh one for this gesture`
  - `panel_holder_319_windows_test.go:149: the spawn was asked 1 time(s) for one gesture after a dead thread, want exactly 2 (the eager first thread plus one rebuild)`
  - 读数：`holder toggle=false; spawns=1; rebuilt thread queued requests=-1; same manager=false`
- **② 现树**：绿。`319 holder rebuild reading: dead thread refused its own gesture; holder toggle=true; spawns=2; rebuilt thread queued requests=1; same manager=true`
- **③ 还原**：绿（同上读数再次得到）；期间与之后 `sha256(resident_windows.go)=06eacd97…`（＝HEAD）、`git diff` 0 字节、`git status --porcelain -- cmd/wisp/ internal/proc/` 0 行。

逐字件：`probes/319/v/01-ac1-headless-three-state.txt`。

## 3. AC#1 活体（winlive，N3）——读数逐字

每发前 `tasklist` 自证：`wisp.exe 40296 …`；一次一枚。

- **原始半**（`TestPanelDeathReplay319UnfixedReading`，现树）：
  `319 unfixed live reading: hwnd=0x18300D02 closed via WM_SYSCOMMAND/SC_CLOSE; thread finished=true; next gesture RequestToggle=false; mgr.IsCreated=false; mgr hwnd=0x0` ⇒ PASS（1.08s）。
- **后半**（`TestPanelDeathThenHolderGestureRebuilds319`，现树）：
  `old hwnd=0x20C0228 … (alive after the close=false); after the close: mgr.IsCreated=false mgr hwnd=0x0 dead-thread RequestToggle=false; after the holder's rebuild: holder toggle=true new hwnd=0x18330D02 mgr.IsCreated=true mgr.IsShown=true old window alive=false` ⇒ PASS（3.14s）。
  生产日志同跑逐字：`panel host: starting a fresh panel thread for this gesture via=319-live-rebuild why="the previous panel thread had exited"`。
- **活体变异**（同 N1 变异）：
  `panel_death_replay_319_winlive_test.go:162: panelHolder.toggle refused the gesture after the panel thread had died; AC#1 says this gesture must make the holder start a fresh thread and the window must come back` ⇒ FAIL（0.85s，**在红名处返回，不是 15s 超时**）。
- **还原**：`old hwnd=0xD007E … new hwnd=0x2C370B42 mgr.IsCreated=true mgr.IsShown=true` ⇒ PASS（7.70s）。
- **拒单半自己的牙**（变异 M-A：`panel_resident_windows.go:373` 摘 `|| rp.isFinished()`）：
  活体 `panel_death_replay_319_winlive_test.go:104: a gesture was accepted by a panel thread whose pump had returned …` ⇒ FAIL；还原 ⇒ PASS（`RequestToggle=false`）。
  （M-A 也把 headless 两枚打红：`:84 a gesture handed to a dead panel thread was accepted (toggle=true show=true)`＋`:142 a gesture was accepted by the dead thread itself`。）

关键判据：**新 HWND ≠ 旧 HWND** 两次独立复现（0x20C0228→0x18330D02、0xD007E→0x2C370B42）；
旧窗 `IsWindow` = false；`mgr.IsCreated()` = true。逐字件：`probes/319/v/02-ac1-live-n3-three-state.txt`。

## 4. AC#2 活体唤窗与形状裁定

**4.1 仓内正式用例（我复跑，不是我引它的日志）**：
`TestSecondLaunchActivationRaisesThePanel319` ⇒ PASS（0.87s）。读数逐字：
`signal=SignalExistingInstance("Local\\wisp-test-319-activation-10040-… .activate") from this process; hook receipt="the panel executor took the request (second-launch); the window verdict lands on the panel thread's own log lines"; panel hwnd=0x18370D02 mgr.IsCreated=true mgr.IsShown=true window alive=true; foreground window=0x18370D02; prior window=0x144307AC; event loop still running=true`。
整发捕获里 `AllowSetForegroundWindow(ASFW_ANY) was refused` 命中 = **0**（无拒绝告警）。

**4.2 我的三枚仓外探针**（overlay 注入，源见附录 A）：

- **形状 B（窗已建、但隐藏——21:40:38 那个现场形状；仓内用例没覆盖 created=true）**：
  `same hwnd 0x64D0C12 shown again=true alive=true; foreground after=0x64D0C12 (panel=0x64D0C12, prior=0xD780E56)` ⇒ PASS。
- **形状 A（窗已可见时二次启动）**：
  `shown before=true after=false; window 0x99507AA alive=true; foreground after=0xA07E8 (panel=0x99507AA, prior=0x126C0E26) -- THE SECOND LAUNCH HID THE PANEL (toggle semantics)` ⇒ 见分歧 D1。
- **真两进程（319-r2 具名没做的那层）**：父进程先自己持前台（editor 窗），spawn 第二进程（同一测试二进制
  `-test.run=^TestProbe319vTwoProcessChild$`）；子进程自建窗、夺前台、再调 `SignalExistingInstance` 并保持存活；父进程的面板窗预先藏起。读数逐字：
  `parent pid=44140 child pid=26676; child window=0x2AB11E4; child took the foreground before signalling=true; ASFW refusal warning in the child=false; parent's panel: hwnd=0x144907AC shown-after=true alive=true; foreground after=0x144907AC (panel=0x144907AC, child=0x2AB11E4, parent prior=0x9B20F74); parent held the foreground at signal time=false; RAISED-BY-THE-SECOND-PROCESS=true` ⇒ PASS。
- **负控（把 `AllowSetForegroundWindow` 调用删掉，同一探针重跑）**：读数除 pid/handle 外**同形**，仍
  `RAISED-BY-THE-SECOND-PROCESS=true` ⇒ **这一形状里 ASFW 不是载荷**（见分歧 D2）。

**4.3 裁定：这一形够不够当 AC#2 的证据？——"附条件成立"**

- **够的部分**：生产线端到端真跑（真命名事件→`RunEventLoop`/`handleActivation`→`SetActivationHook`→
  `secondLaunchActivation319`→`holder.toggle`→真窗 Show＋`SetForegroundWindow`），三种形状都真夺到前台；
  回执句不夸大（只说"请求交出去了，窗的判定在面板线程自己的日志行上"）。
- **同进程形的边界（逐条）**：
  1. **进程边界**：仓内用例的"第二启动"是本进程调 `SignalExistingInstance`；内核对象操作（OpenEvent/SetEvent）
     与跨进程同义，但"调用者的资格"不同源。本腿的真两进程探针补了这一层（4.2 第三发）且通过——**但它不是
     常驻钉子，仓内正式形仍是同进程**。
  2. **ASFW 由"前台进程"发出的那一层**：本腿负控显示：在我能造的形状里把 ASFW 删掉结果不变 ⇒ **"ASFW 是
     现场夺前台的必要条件"未被证明**。现场形状（前台被**无关**进程如 Explorer/浏览器占着、第二实例由外壳
     拉起）仍未测得——子进程是常驻进程自己 spawn 的，Windows 的前台规则可能对它另有放行。**不能把 ASFW 那
     一跳记成已证**。
  3. **冷启动 vs 热路径**：仓内用例是冷启动（created=false）；本腿形状 B 补了"窗已建但隐藏"的 re-show 热
     路径（真出窗＋真夺前台，PASS）；形状 A 补了"窗已可见"的第三种态，读出**藏窗**（D1）。
  4. **`RunEventLoop` 残留 goroutine**：用例注释已点名。我的读数是"不污染后续用例"：`handleActivation`
     每 50ms 读 `rt.Instance.ActivateEvent()`；用例 cleanup 的 `rt.Shutdown(true)` → `Instance.Release()`
     把事件句柄置 0（`singleinstance_windows.go:92-96`），之后的 `WaitForSingleObject(0)` 报错、该轮不消费
     （`boot_windows.go:198-204`）；它不经 `Registry.Spawn`、不进 D38 名册，无名册读数能点到它；winlive 又
     是一枚一进程。⇒ 接受为"具名残留"，不是阻断项；若将来同进程连跑多枚 winlive，需重新评估。

## 5. AC#3 三钉与变异（逐钉，自己造）

- **N1（死线拒单＋holder 重建，含 spawn 第二次被问）**：见 §2 与 §3 的 M-A/N1 两发；红名逐字，还原复绿。
- **N2（hook 恰一次＋Reset；无 hook 仍消费并清）**：
  - 变异 N2a＝删 `boot_windows.go:205-212` 整段 hook 分支 ⇒ 红名：
    `activation_hook_319_windows_test.go:79: the registered activation hook ran 0 time(s) for ONE second launch, want exactly 1`；
    同跑 `TestActivationWithoutHook319StillConsumesAndClears` 仍绿（无 hook 分支的正控）。
  - 变异 N2b＝只删 `boot_windows.go:210` 那行 `ResetActivation` ⇒ 红名：
    `:85 a second round consumed a signal: ResetActivation did not clear the manual-reset event`＋
    `:88 the activation hook ran 2 time(s) across two rounds, want 1 - one launch's signal leaked into two handoffs`。
  - 还原 ⇒ 两枚全绿（`ok internal/proc 0.212s`）。
- **N3（活体两半）**：见 §3：原始半绿＋活体变异红＋后半绿＋换柄读数。逐字件：`probes/319/v/04-ac3-mutants.txt`。

## 6. 回归

- 改前既有 4 枚：`TestDeadPanelThreadRefusesEveryGesture319`（绿，读数 `RequestToggle=false RequestShow=false isFinished=true`）·
  `TestActivationHook319RunsExactlyOnceAndResets`（绿）· `TestActivationWithoutHook319StillConsumesAndClears`（绿，日志逐字
  `handled_here="booked and cleared only"`）· `TestPanelDeathReplay319UnfixedReading`（活体绿）。
- 合并：`go test -run '319' ./cmd/wisp/ ./internal/proc/ -count=1` ⇒
  `ok cmd/wisp 0.132s` ＋ `ok internal/proc 0.279s`（在 319 面上红灯 0）。
- ⛔ 未跑全包裸 `go test ./cmd/wisp/`（重包正被 320-r1 写；本腿按派单只跑 319 面回归；r2 的全包七红读数**未被本腿采用**）。

## 7. AC#4 门禁＋面核

- `go build ./...` rc=0（活树与 `git archive 33e429c1` 快照两次）。
- `gofumpt -l`：319 的 7 枚 Go 件无输出（rc=0）；本腿零 Go 写件。
- `go vet -tags winlive ./cmd/wisp/ ./internal/proc/`：
  - 13:40 活树一次 **rc=1**，但归因明确：`internal\config\settings_320.go:26:2: "reflect" imported and not used`
    —— 320-r1 当时**未提交的**半成品（`git status` = ` M internal/config/settings_320.go`，mtime 13:40），与 319 面无关；
  - 干净快照（`/tmp/wisp319v-main`）rc=0；320-r1 落地（HEAD `a920c4d4`）后活树再跑 rc=0。
- `sh scripts/d22scan.sh` rc=0：`d22scan: clean - no D22 ban violations; … bans #1-5 internal/=238 … ban #8 cmd/=132 …`。
- `sh scripts/check-path-length-budget.sh --with-self-test` rc=0：`VERDICT GREEN … not in roster=0`。
- 面核：`e6d1373e` 11 枚＝4 生产件＋3 测试件＋票面＋probes/319/r1/logs×3；`2ad0aa86` 13 枚＝2 测试件＋票面＋probes/319/r2×10；
  冻结点名（`docs/PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`frontend/**`／`design/**`）**两笔零命中**。
- 工作树里 319 面干净；脏件全属别家/机主（320-r1 的 `panel_config_*`、orchestrator 的 HANDOVER/台账、
  机主外部 agent 的 `design/**`+`frontend/**`、`go.mod` 的 indirect→direct、probes/152 与 probes/268 半成品——按派单未碰）。
- 逐字件：`probes/319/v/05-gates-face.txt`。

## 8. 分歧登记（逐处，供编排者裁定）

- **D1（AC#2，新发现，读数见 §4.2 形状 A）**：二次启动驱动的是 `toggle`；`RequestToggle` 对**已显示**的面板走
  `mgr.Hide()`（`cmd/wisp/panel_resident_windows.go:447-457`）。⇒ 面板**已可见**时双击 exe 会把面板**藏起来**，
  而不是"唤到前台"。票面字面对这一态不成立（`唤到前台`）。现场的常见态是"隐藏/不存在"（该态已证真出窗+夺前台），
  故本腿判 AC#2 附条件成立而非不成立；**要不要把 AC#2 的执行器从 toggle 改成"只展示不切换"由编排者/票面裁定**。
- **D2（AC#2）**：`AllowSetForegroundWindow(ASFW_ANY)` 在真两进程形状里**无它也一样夺到前台**（负控同形通过）⇒
  ASFW 一跳的"现场必要性"未证。若要证，需一枚前台被**无关**进程占据的形状（不能机主实例做背景，需另议仪器）。
- **D3（AC#1，既有）**：`internal/proc` 的 `RunEventLoop` 无导出停止口，winlive 用例留残留 goroutine（319-r2 已具名；
  本腿给出"不污染"的机制读数，见 §4.3-4）。
- **D4（观感，非 319 缺陷）**：活体用例日志里 `goroutine outside the D38 roster (leak symptom) goroutine=panel-sta` 的告警
  在**改前既有用例**（raw 半）也出现 ⇒ 属既有观测面行为，319 的 holder 没有新增它（重建后同样只多一条同名告警）。

## 9. 具名没做

- ⛔ 未造"前台被无关进程占据"的 ASFW 隔离形状（D2 的补证路径未建）；未把真两进程探针升级为常驻钉子（那要改仓内测试件，出本腿面）。
- ⛔ 未跑全包裸 `go test ./cmd/wisp/`（重包被 320-r1 占用/正在落地；只跑 319 面回归与 4 枚活体单发）。
- ⛔ 未验 `frontend/**` 侧任何"可见提示"（面归机主外部 agent；仓内只用了既有 `SetTrayTip`，本腿未做 UI 验收）。
- ⛔ 未模拟 WebView2 崩/`DestroyWindow` 直调等其它致死变体（AC#0 的乙/丙假说 census 已反证；不在本票验收格内）。
- ⛔ 未推（只 commit）；未碰台账/HANDOVER/AC 框/票名；未碰 probes/152、probes/268、`design/**`。
- ⛔ 全程未 kill/未关窗机主实例 `wisp.exe` 40296；未拉第二个 `--env dev`（mutex/event 名与机主同名）。

---

## 附录 A · 探针源码（overlay 注入，仓外件；本腿仪器，非仓内测试面）

### A.1 probe_319v_hotpath_winlive_test.go（形状 A/B）

```go
//go:build windows && winlive

// 319-v probe (verification leg instrument; NEVER part of the repo's test face -
// it is injected with -overlay from outside the tree and archived under
// .scratch/wisp/probes/319/v/). It measures the two AC#2 shapes the committed
// case TestSecondLaunchActivationRaisesThePanel319 does not cover, both on the
// real production path (Proc.Boot + RunEventLoop + SetActivationHook +
// secondLaunchActivation319 + panelHolder.toggle + SignalExistingInstance):
//
//   A. the panel window is already VISIBLE when the second launch arrives.
//      RequestToggle hides a shown panel, so the ticket's letter ("bring the
//      existing instance's UI to the front") and the implementation (toggle)
//      can disagree here. The reading says which one happened.
//   B. the panel window EXISTS but is hidden (created earlier, then toggled
//      away) - the field shape behind the 2026-10-14 21:40:38 entry. The
//      committed case only covers created=false (cold first window).
//
// Both cases leave RunEventLoop's goroutine behind exactly like the committed
// case names; each is run one at a time, one process per case.

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/buildinfo"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/proc"
	"golang.org/x/sys/windows"
)

func probe319vBoot(t *testing.T) *proc.Runtime {
	t.Helper()
	layout := proc.Layout{
		Env:               buildinfo.EnvProd,
		DataDir:           t.TempDir(),
		MutexName:         fmt.Sprintf(`Local\wisp-test-319v-probe-%d-%d`, os.Getpid(), time.Now().UnixNano()),
		ActivateEventName: fmt.Sprintf(`Local\wisp-test-319v-probe-%d-%d.activate`, os.Getpid(), time.Now().UnixNano()),
		MutexEnabled:      true,
	}
	rt, err := proc.Boot(buildinfo.EnvProd, proc.WithLayout(layout), proc.WithRegistry(observe.NewRegistry()))
	if err != nil {
		t.Fatalf("proc.Boot: %v", err)
	}
	t.Cleanup(func() { rt.Shutdown(true) })
	return rt
}

// probe319vFirstWindow opens the panel the production way and returns hwnd.
func probe319vFirstWindow(t *testing.T, h *panelHolder, mgr *PanelManager) uintptr {
	t.Helper()
	if !h.toggle("319v-probe-first-open") {
		t.Fatalf("the holder refused the first open")
	}
	waitPanelTrue(t, "the first panel window to be shown", func() bool { return mgr.IsShown() })
	hwnd := mgr.windowHandle()
	if hwnd == 0 {
		t.Fatalf("the first show reported HWND 0")
	}
	return hwnd
}

// probe319vLaunch runs one production second-launch hop and returns the receipt.
func probe319vLaunch(t *testing.T, rt *proc.Runtime, h *panelHolder) string {
	t.Helper()
	receipts := make(chan string, 4)
	hook := secondLaunchActivation319(h, nil)
	rt.SetActivationHook(func() string {
		s := hook()
		select {
		case receipts <- s:
		default:
		}
		return s
	})
	go rt.RunEventLoop()
	if err := proc.SignalExistingInstance(rt.Layout.ActivateEventName); err != nil {
		t.Fatalf("SignalExistingInstance: %v", err)
	}
	select {
	case s := <-receipts:
		return s
	case <-time.After(panelThreadWait):
		t.Fatalf("the activation hook never answered within %v", panelThreadWait)
		return ""
	}
}

// TestProbe319vSecondLaunchWhilePanelVisible is shape A.
func TestProbe319vSecondLaunchWhilePanelVisible(t *testing.T) {
	h, mgr := startPanelHolder319(t)
	hwnd := probe319vFirstWindow(t, h, mgr)

	// The "panel open but not foreground" analogue: a prior window of this
	// process takes the front. The panel stays visible behind it.
	prior := editorOnPanelThread(t, h.rp)
	if fg := uintptr(windows.GetForegroundWindow()); fg != uintptr(prior) {
		t.Fatalf("the prior window 0x%X did not take the foreground (0x%X); measurement undecidable", uintptr(prior), fg)
	}
	shownBefore := mgr.IsShown()

	rt := probe319vBoot(t)
	receipt := probe319vLaunch(t, rt, h)
	// The hook answers when the request is QUEUED; give the panel thread a
	// bounded moment to act on it.
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) && mgr.IsShown() == shownBefore {
		time.Sleep(10 * time.Millisecond)
	}
	shownAfter := mgr.IsShown()
	fgAfter := uintptr(windows.GetForegroundWindow())
	if shownAfter {
		t.Logf("319v shape A READING (panel already visible): shown before=%v after=%v; window 0x%X alive=%v; foreground after=0x%X (panel=0x%X, prior=0x%X); receipt=%q -- the second launch did NOT hide the panel", shownBefore, shownAfter, hwnd, isLiveWindow(hwnd), fgAfter, hwnd, uintptr(prior), receipt)
	} else {
		t.Logf("319v shape A READING (panel already visible): shown before=%v after=%v; window 0x%X alive=%v; foreground after=0x%X (panel=0x%X, prior=0x%X); receipt=%q -- THE SECOND LAUNCH HID THE PANEL (toggle semantics), which is not the ticket's letter for this shape", shownBefore, shownAfter, hwnd, isLiveWindow(hwnd), fgAfter, hwnd, uintptr(prior), receipt)
	}
}

// TestProbe319vSecondLaunchWhilePanelHidden is shape B: window exists, hidden.
func TestProbe319vSecondLaunchWhilePanelHidden(t *testing.T) {
	h, mgr := startPanelHolder319(t)
	hwnd := probe319vFirstWindow(t, h, mgr)

	// Toggle it away the same way a tray click would (the field shape).
	if !h.toggle("319v-probe-hide") {
		t.Fatalf("the holder refused the hide toggle")
	}
	waitPanelTrue(t, "the panel window to be hidden", func() bool { return !mgr.IsShown() })
	if !isLiveWindow(hwnd) {
		t.Fatalf("the window 0x%X is not alive after a hide - not the field shape (hide must not destroy)", hwnd)
	}
	prior := editorOnPanelThread(t, h.rp)
	if fg := uintptr(windows.GetForegroundWindow()); fg != uintptr(prior) {
		t.Fatalf("the prior window 0x%X did not take the foreground (0x%X); measurement undecidable", uintptr(prior), fg)
	}

	rt := probe319vBoot(t)
	receipt := probe319vLaunch(t, rt, h)
	waitPanelTrue(t, "the hidden panel to be shown by the second launch", func() bool {
		return mgr.IsShown() && mgr.windowHandle() == hwnd && uintptr(windows.GetForegroundWindow()) == hwnd
	})
	fgAfter := uintptr(windows.GetForegroundWindow())
	t.Logf("319v shape B READING (panel exists but hidden): same hwnd 0x%X shown again=%v alive=%v; foreground after=0x%X (panel=0x%X, prior=0x%X); receipt=%q", hwnd, mgr.IsShown(), isLiveWindow(hwnd), fgAfter, hwnd, uintptr(prior), receipt)
}
```

### A.2 probe_319v_twoprocess_winlive_test.go（真两进程 + ASFW 负控）

```go
//go:build windows && winlive

// 319-v probe #2 (verification leg instrument; injected with -overlay, never in
// the tree). It measures the layer 319-r2 named as not built: a REAL second
// process calls proc.SignalExistingInstance (AllowSetForegroundWindow(ASFW_ANY)
// included) while the resident's process does NOT hold the foreground, so the
// resident's SetForegroundWindow can only succeed if that grant is real.
//
// Shape: the parent (resident side) has its panel window created and hidden,
// then holds the foreground with one of its own editor windows. It spawns the
// child (this same test binary, -test.run=^TestProbe319vTwoProcessChild$). The
// child creates ITS OWN window, takes the foreground with it, prints READY + HWND,
// then calls SignalExistingInstance on the parent's per-run event name and stays
// alive holding the foreground. The parent then waits for its panel to be shown
// AND to become the foreground window. If the grant were refused, the panel
// would still be shown (Show runs regardless) but the foreground would stay on
// the child's window - that difference is what this probe reads.
//
// The owner's running instance is untouched: the event name is this run's own
// (Local\wisp-test-319v-2p-<pid>-<ns>), never the dev/prod names.

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/proc"
	"golang.org/x/sys/windows"
)

// TestProbe319vTwoProcessChild is the child half (selected by name).
func TestProbe319vTwoProcessChild(t *testing.T) {
	eventName := os.Getenv("WISP319V_CHILD_EVENT")
	if eventName == "" {
		t.Skip("319v two-process helper: only meaningful as a spawned child")
	}
	hwnd, err := createEditorWindow("wisp 319v child foreground")
	if err != nil || hwnd == 0 {
		fmt.Printf("CHILD319V WINDOWFAIL hwnd=0x%X err=%v\n", uintptr(hwnd), err)
		os.Exit(3)
	}
	pnlShowWindow.Call(uintptr(hwnd), pnlSwShow)
	pnlUpdateWindow.Call(uintptr(hwnd))
	sf, _, _ := pnlSetForeground.Call(uintptr(hwnd))
	fg := uintptr(windows.GetForegroundWindow())
	fmt.Printf("CHILD319V READY pid=%d hwnd=0x%X setforeground=%d fg=0x%X\n", os.Getpid(), uintptr(hwnd), sf, fg)
	time.Sleep(500 * time.Millisecond)
	if err := proc.SignalExistingInstance(eventName); err != nil {
		fmt.Printf("CHILD319V SIGNALFAIL err=%v\n", err)
		os.Exit(4)
	}
	fmt.Printf("CHILD319V SIGNALED event=%q (the AllowSetForegroundWindow refusal warning, if any, is on stderr above)\n", eventName)
	// Hold the foreground until the parent has read the outcome.
	time.Sleep(20 * time.Second)
	os.Exit(0)
}

// TestProbe319vSecondLaunchFromRealChildProcess is the parent half.
func TestProbe319vSecondLaunchFromRealChildProcess(t *testing.T) {
	h, mgr := startPanelHolder319(t)
	hwnd := probe319vFirstWindow(t, h, mgr)
	if !h.toggle("319v-probe-hide") {
		t.Fatalf("the holder refused the hide toggle")
	}
	waitPanelTrue(t, "the panel to be hidden", func() bool { return !mgr.IsShown() && isLiveWindow(hwnd) })
	prior := editorOnPanelThread(t, h.rp)
	if fg := uintptr(windows.GetForegroundWindow()); fg != uintptr(prior) {
		t.Fatalf("the parent's prior window 0x%X did not take the foreground (0x%X); measurement undecidable", uintptr(prior), fg)
	}

	rt := probe319vBoot(t)
	receipts := make(chan string, 4)
	hook := secondLaunchActivation319(h, nil)
	rt.SetActivationHook(func() string {
		s := hook()
		select {
		case receipts <- s:
		default:
		}
		return s
	})
	go rt.RunEventLoop()

	logPath := filepath.Join(t.TempDir(), "child-319v.log")
	logFile, err := os.Create(logPath)
	if err != nil {
		t.Fatalf("child log file: %v", err)
	}
	defer logFile.Close()
	child := exec.Command(os.Args[0], "-test.run=^TestProbe319vTwoProcessChild$", "-test.v")
	child.Env = append(os.Environ(), "WISP319V_CHILD_EVENT="+rt.Layout.ActivateEventName)
	child.Stdout = logFile
	child.Stderr = logFile
	if err := child.Start(); err != nil {
		t.Fatalf("spawning the second process: %v", err)
	}
	defer func() {
		_ = child.Process.Kill()
		_ = child.Wait()
	}()

	childHwnd := uintptr(0)
	waitPanelTrue(t, "the child process to take the foreground with its own window", func() bool {
		b, _ := os.ReadFile(logPath)
		m := regexp.MustCompile(`CHILD319V READY pid=\d+ hwnd=0x([0-9A-Fa-f]+)`).FindSubmatch(b)
		if m == nil {
			return false
		}
		fmt.Sscanf(string(m[1]), "%X", &childHwnd)
		return childHwnd != 0 && uintptr(windows.GetForegroundWindow()) == childHwnd
	})

	// The child signals after READY; wait for the receipt, then for the panel.
	var receipt string
	select {
	case receipt = <-receipts:
	case <-time.After(panelThreadWait):
		t.Fatalf("the activation hook never answered - the child's signal reached nothing")
	}
	raised := true
	deadline := time.Now().Add(panelThreadWait)
	for time.Now().Before(deadline) {
		if mgr.IsShown() && uintptr(windows.GetForegroundWindow()) == mgr.windowHandle() {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	raised = mgr.IsShown() && uintptr(windows.GetForegroundWindow()) == mgr.windowHandle()
	fgAfter := uintptr(windows.GetForegroundWindow())
	b, _ := os.ReadFile(logPath)
	refused := bytes.Contains(b, []byte("AllowSetForegroundWindow(ASFW_ANY) was refused"))
	t.Logf("319v two-process READING: parent pid=%d child pid=%d; child window=0x%X; child took the foreground before signalling=true (checked while waiting); ASFW refusal warning in the child=%v; parent's panel: hwnd=0x%X shown-after=%v alive=%v; foreground after=0x%X (panel=0x%X, child=0x%X, parent prior=0x%X); parent held the foreground at signal time=false; receipt=%q; RAISED-BY-THE-SECOND-PROCESS=%v",
		os.Getpid(), child.Process.Pid, childHwnd, refused, hwnd, mgr.IsShown(), isLiveWindow(hwnd), fgAfter, mgr.windowHandle(), childHwnd, uintptr(prior), receipt, raised)
	if !raised {
		t.Errorf("the second process's signal was consumed (%q) but the panel did not end up shown AND foreground (shown=%v fg=0x%X want 0x%X) while the child held 0x%X%s",
			receipt, mgr.IsShown(), fgAfter, mgr.windowHandle(), childHwnd, map[bool]string{true: " (ASFW was refused by the child)", false: ""}[refused])
	}
}
```

（注：两枚探针注入件的 `gofmt -l` 空、`go vet -overlay -tags winlive ./cmd/wisp/` rc=0；它们共用 `panel_319` 活体件与 A.1 的 `probe319vBoot`/`probe319vFirstWindow` helpers。）

## 附录 B · 变异清单（全部仓外副本 + overlay；靶件 sha256 收尾＝HEAD）

| # | 靶件:行 | 变异 | 用途 |
|---|---|---|---|
| N1 | `cmd/wisp/resident_windows.go:535` | `if h.rp == nil \|\| h.rp.isFinished() {` → `if h.rp == nil {` | headless＋live 重建红 |
| M-A | `cmd/wisp/panel_resident_windows.go:373` | `if rp.startUpErr() != nil \|\| rp.isFinished() {` → `if rp.startUpErr() != nil {` | 死线拒单红（headless＋live） |
| N2a | `internal/proc/boot_windows.go:205-212` | 整段 hook 分支删除 | hook 恰一次红 |
| N2b | `internal/proc/boot_windows.go:210` | 删 hook 分支内的 `ResetActivation` 一行 | Reset 红 |
| C-ASFW | `internal/proc/singleinstance_windows.go:170-173` | ASFW 调用块删除（保留 import 形状） | 真两进程负控 |
