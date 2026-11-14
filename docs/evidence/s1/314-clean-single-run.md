# 票 314 — 独占干净全量：**单发不复现「整进程死亡／挂起」**，并发才是机制

> 裁决人：编排者（非实现者）。时间 `2026-11-14 09:26–09:50 +08`。机器：本机 Windows。
> 本件回答票 314 的核心前提：「`go test ./cmd/wisp` 在窗口区会把整个进程打死」。

## 1. 方法（可逐字重跑）

```bash
# 前提：同一时刻没有任何别的 wisp.test 在跑（脚本自己先量了）
sh scripts/wisp-cli-staged-tests.sh -test.timeout=25m
```
- 脚本 `D:\tmp\wisp\20261114\314-clean-exclusive.sh`，日志 `D:\tmp\wisp\20261114\314-clean-exclusive.txt`（7126 行）。
- 开跑前实测：`competitor wisp.test processes before: 0`，`LogonUI: 0`。
- **并发铁律**：本编队此前犯过一次错——编排者自己的锁屏测量与一条腿的全量**同时在跑**，导致两条读数作废（见 `A1860`）。本次是**唯一**一发。

## 2. 读数（逐字）

```
CLEAN-FULL-RC=1
CLEAN-FULL-ELAPSED=1442s
PASS=786  FAIL=1  SKIP=12
末尾：--- PASS: TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask (0.34s)
      FAIL  github.com/CarlosShao/wisp/cmd/wisp  1439.00s
```

### 关键三点
1. **没有信号死亡**：`rc=1`（测试红），**不是** `rc=139`／`128+n`。
2. **没有挂起**：1442s＝24 分钟跑完并自然退出，**没有** `panic: test timed out`，没有需要外部 kill。
3. **跑到了最后一枚**：终点锚 `TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask` 有 `PASS` 裁决行（该锚由 `832e04489` 引入，见 `314-end-anchor.md`）。

## 3. 唯一一枚红（与 314 无关）

`TestAC4FocusReturnToPriorWindowGap33r5`——真窗前台时序类，既有已知红，**不是**进程死亡类。

## 4. 12 枚 SKIP 全部具名（票 71 AC#3 纪律）

| 用例 | 具名理由（`文件:行`） |
|---|---|
| `TestAC247LiveMicrophoneLevelsReachTheBallSeam` | `resident_audio_247_live_windows_test.go:130` — 需要真麦克风，两把 ruler 都读 0，`WISP_LIVE_MIC=1` 才跑 |
| `Test41a2r1DrainTimeoutVerdictOptIn` | `resident_voice_41a2r1_windows_test.go:251` — opt-in：烧 2×20s 预算 |
| `Test41a2r1DrainTimeoutArmChild` | `:324` — 上面那枚的子进程入口，靠 re-exec 才跑 |
| `Test26PEcho…`／`Test26PWithoutAEC…`／`Test26PMargin…`／`Test26PBargein…`／`Test26PFeedOnline…` | `resident_voice_bargein_26p_*_windows_test.go:109/364/343/521` — **AEC3 绑定不在本构建**（`-tags webrtcap` ＋ deps.toml CGO 旗） |
| `Test392Keepalive…`／`Test392AecTuning…`／`Test392Dispose…` | `resident_voice_bargein_392_windows_test.go:404/554/647` — 同上，AEC3 不在本构建 |
| `Test16r1PremortemMicPermissionDeniedLiveProbe` | `voice_e2e_demo_16r1_windows_test.go:461` — live probe，需 16-r1 pre-mortem rig 翻过 OS 同意存储 |

## 5. 裁定

- **票 314 的「整进程死亡（rc=139）」在单发干净条件下不复现。**
- **此前两处异常读数（`rc=139` @370s、90 分钟挂起）都出现在「两发整包并发」的窗口里**（编排者自己的测量与一条腿的全量相撞）⇒ **机制是并发，不是产品在窗口区自杀**。
- 这与诊断腿 `314-hang-rootcause.md` 的结论一致：真机制＝**竞态 ＋ 嵌套变异 `go test` 抬长整包**（`approval_always_219ac2_test.go:379` 在整包里又跑一遍整包，单枚 `Test219AC2Mutation` 39.70s），于是外部时钟落在哪一枚**随负载漂**。
- 已落地的两道闸：①**终点锚**（`832e04489`）让「没跑完」响亮变红并点名；②**并发纪律**写进 goal：`go test ./cmd/wisp` 整包**同一时刻只许一发**。

## 6. 仍未验的一格（不糊）

- **候选 (a)「锁屏／空闲桌面对真窗测试的影响」仍未测**：编排者尝试自己 `rundll32 user32.dll,LockWorkStation`，实测 `LogonUI` 前后均为 **0**、console 仍 `Active`（后台进程里静默失败）。
- 待机主离开时按一次 `Win+L`，编排者轮询检测到锁屏后自动跑同一发，再补读数。
