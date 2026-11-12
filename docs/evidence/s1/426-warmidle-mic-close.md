# 426 — Warm 行 mic close 只 1 次：**断言读得太早**（测量竞态），不是产品缺陷

- 票：`.scratch/wisp/issues/426-warm-idle-mic-close-counted-once-against-two-because-the-second-defence-never-runs.md`
- 写面：**测试**（`cmd/wisp/resident_voice_wire_r1_windows_test.go`），产码零改动
- 判定：**甲形 —— 断言落后于产码**（但**不是**「设计上只该关 1 次」，而是**断言在级联跑完之前就读了计数器**）
- 腿：warmidle（2026-11-12 20:0x–20:1x +08）

---

## AC#0 只读取证

### 1. 逐字红句 ＋ 耗时（carrier，`D:\work\soft\Git\bin\bash.exe`）

```
D:\work\soft\Git\bin\bash.exe -c "scripts/wisp-cli-staged-tests.sh '-test.run' '^TestWireR1WarmIdleDisposesOnTheOwner$'"
```

```
    resident_voice_wire_r1_windows_test.go:745: mic close counted 1 times, want at least 2 (the Warm row and the cascade's defence in depth)
--- FAIL: TestWireR1WarmIdleDisposesOnTheOwner (6.90s)
FAIL
```

日志：`D:\tmp\wisp\20261112\426-warmidle\run1.log`。与票面两次复跑（5.19s／6.81s）同句同数 ⇒ **稳定真红，确认可复现**。

### 2. ★ 关键读数：`:750` 那行 `WARM-IDLE` **从来没打印过**

红在 `:745` 的 `t.Fatalf`，`t.Logf` 在 `:750` 之后 ⇒ 原断言红时 `WARM-IDLE` 一行**根本不会出现**。
票面 AC#0.2 要求"逐字贴 `:750` 的 `t.Logf` 输出"——**在修复前的红态下该行不存在**，这是必然的（Fatalf 中止）。

因此我改用**仓外 overlay 探针**（`overlay-probe.json`，替换件 `D:\tmp\wisp\20261112\426-warmidle\probe-wire-r1_test.go`）
在**原断言之前**插入一个**owner 队列屏障**（等级联真跑完），把两发读数一次打出来：

```
    resident_voice_wire_r1_windows_test.go:756: PROBE-AFTER-BARRIER mic_closes=2 stages_at_mic_close=[Awake Disposed] disposals=1
    resident_voice_wire_r1_windows_test.go:764: WARM-IDLE stage=Disposed engines_hot=0 mic_closes=2 stages_at_mic_close=[Awake Disposed] engine_closes=1
--- PASS: TestWireR1WarmIdleDisposesOnTheOwner (17.17s)
```

日志：`D:\tmp\wisp\20261112\426-warmidle\probe1.log`

**★ 这是本票的决定性读数**：只要等到级联跑完，`mic_closes` **就是 2**，`stages_at_mic_close=[Awake Disposed]`。
⇒ **第二次关麦从来没有缺失过**，缺的是**测试读计数器的时间点**。

### 3. 三处产码：钩子实际被调用几次、分别在哪条路径

| # | 调用点 | 路径 | 该处 `Stage()` 读数 |
|---|---|---|---|
| 1 | `internal/session/scope.go:470`（`EnterWarm` 第一步，**在落 `StageWarm` 之前**，`:474`） | 回合结束 → Warm 行 | **`Awake`** |
| 2 | `internal/session/scope.go:691`（`runCascade` 末尾，**在 `beginDisposeLocked` 已落 `StageDisposed` 之后**，`:654`） | Settling 过期 → dispose 级联 | **`Disposed`** |

- 钩子本体：`cmd/wisp/resident_voice_engines_windows.go:488-504`，自增点唯一一处 `:497`。
- 级联的防御注释在 `scope.go:687-692` 逐字：「Defence in depth, deliberately AFTER the four frozen steps… a destroyed session must never leave the device armed.」⇒ **"the cascade's defence in depth" 指的就是 `runCascade` 里这一枚**，不是别人的。
- `closeEpoch`（`resident_voice_engines_windows.go:784-799`）只调 `scope.Dispose`，不自己调 `MicClose`；`rv.stop()`（`resident_voice_windows.go:789`）走的是 `plugin.DisposalScope`，与 `session.Hooks.MicClose` 无关 ⇒ **生产里 `MicClose` 只有上面两条路径，每条各一次**。

### 4. `R28-2` 语义（`tri-v1-verdict.md:118` 与 `:123` 逐字）

- `:118`：`R28-2` "mic closed"是**唤醒闩**、帧仍进 ASR、设备未关（Key constraint，等 41）
- `:123`：「**票面必须把 `R28-2` 写在 `Warm` 定义旁**（"mic closed"＝**唤醒闩**，不是设备/喂食门），否则下一代会把那句当已实现。」

⇒ `MicClose` 钩子在**本仓的语义是唤醒闩那一层**（`resident_voice_engines_windows.go:336-343` 与 `resident_voice_windows.go:738-745` 都逐字写了：关的是 **FEED**，不是设备）。
**所以这条断言数的正是唤醒闩，语义没错**——它没有把"设备关闭"错当"唤醒闩"。

### 5. 票 28 的 `Warm` 定义段（`28-session-scope-warm-done.md:20` 逐字）

> `Warm`: models hot + mic CLOSED + faint warm breathing; next wake via click/hotkey → Listening

以及 `:42` AC2 逐字：「Timers: 90s Warm → Settling → dispose; 30s Conversation → Warm with mic closed FIRST」。
`:73` 记 `Warm="mic CLOSED"` 未实现／票 41 前无法实现；`:90` 记 `R28-2` 已于 2026-10-31 由票 41（`warmFeed` 真喂食门）**注销**。

⇒ **`Warm` 的定义要求"关麦 FIRST"，且"关麦"＝唤醒闩/喂食门**。断言要求的两次（Warm 一次＋级联兜底一次）**与设计一致**，不多也不少。

### 6. ★ 判定：**甲形（断言过严）——但根因是测量竞态，不是承诺过严**

> **一句话判定**：**断言过严**——`micCloseCount() >= 2` 这条承诺本身是对的（级联确实该关第二次），
> 但它在 `StageDisposed` 一可读时就读计数器，而 `StageDisposed` 由 `beginDisposeLocked` 在 `internal/session/scope.go:654` **先落**，
> 级联的 `MicClose` 在 `internal/session/scope.go:691` **后跑** ⇒ 读早了，读到的只有 Warm 行那一次。

**铁证（同一份日志，时间戳）**：

```
time=...T20:03:03.177  level=INFO msg="voice[scope]: session: SCOPE-SETTLING reason=warm.idle ..."
    resident_voice_wire_r1_windows_test.go:745: mic close counted 1 times, want at least 2 ...   ← 红在这一刻
time=...T20:03:03.228  level=INFO msg="voice[scope]: session: SCOPE-DISPOSED reason=settling.expired ..."
```

`SCOPE-DISPOSED` 是 `runCascade` 的**最后一行日志**（`scope.go:693`），它排在 `:691` 的 `MicClose` **之后**。
测试在 `.177` 就断言了，而 `.228` 才走完级联 ⇒ **测试读计数器的时刻，级联的第二次关麦还没发生**。这不是产品缺陷，是**测试与它要断言的钩子赛跑**。

**为什么不是乙形（产品缺陷）**：探针加了屏障后 `mic_closes=2`／`[Awake Disposed]`，**产码一行没改**。若真是产品缺陷，加屏障不会让第二次凭空出现。

**为什么也不是「设计上只该关 1 次」**：`scope.go:687-692` 的注释与 `tri-v1-verdict.md:118/:123` 的 `R28-2` 都说明级联兜底是**刻意设计**的第二次；票 28 `:42` 的 AC2 也只要"mic closed FIRST"，没否定第二次。所以**承诺是两次，产码也给两次**，缺的只是观测时机。

---

## AC#1 修形：只改测试（产码零改动）

### 选形理由

1. **产码是对的** ⇒ 不能改产码（改了就成乙形，且会把一个正确的设计改坏）。
2. 票面明令 **⛔ 不许放宽成"≥1 就算过"**（那等于删掉这条测试）。我的改法**没有放宽**：
   - 仍然要求 `closes >= 2`（原承诺一字不改）；
   - 新增 `len(stages) >= 2`；
   - 新增**两条顺序钉子**：`stages[0] == "Awake"`、`stages[last] == "Disposed"`。
3. 顺序钉子把票面担心的"级联若关则必须发生在 Warm 之后"变成**可执行**的：`[Awake Disposed]` 才过；`[Warm Disposed]`（关麦晚于落状态）或 `[Awake Awake]`（第二次是假的）都红。
4. 屏障用 `rv.queue.post`：级联跑在 owner goroutine（`voiceOwnerClock` 把定时器回调 marshal 到它，`resident_voice_engines_windows.go:207-220`），所以**排在 settle fire 之后 post 的 job 必然排在级联之后**——这是**顺序保证**，不是 sleep 猜时长。
5. 不引入墙钟时间差做超时判断（D42#9 的禁项）：`time.After(20*time.Second)` 只作**测试防挂死**上限，与生产超时无关，且 `voice15eWaitFor` 本来就是同一形状（`resident_voice_15e_windows_test.go:269`）。

### 修后绿发

```
D:\work\soft\Git\bin\bash.exe -c "scripts/wisp-cli-staged-tests.sh '-test.run' '^TestWireR1WarmIdleDisposesOnTheOwner$'"
```

```
    resident_voice_wire_r1_windows_test.go:786: WARM-IDLE stage=Disposed engines_hot=0 mic_closes=2 stages_at_mic_close=[Awake Disposed] engine_closes=1 disposals=1
--- PASS: TestWireR1WarmIdleDisposesOnTheOwner (6.59s)
PASS
```

日志：`D:\tmp\wisp\20261112\426-warmidle\green1.log`

---

## ★ 突变控制（三发，全部经 `go test -overlay` 真跑，拷贝全在仓外 `D:\tmp\`）

### M1（票面指定的那一发）：删掉 `:497` 的 `r.micCloses++`

替换件 `D:\tmp\wisp\20261112\426-warmidle\m1-engines.go`（只删这一行），overlay `overlay-m1.json` 指向 `cmd/wisp/resident_voice_engines_windows.go`。

```
    resident_voice_wire_r1_windows_test.go:767: mic close counted 0 times, want at least 2 (the Warm row and the cascade's defence in depth)
--- FAIL: TestWireR1WarmIdleDisposesOnTheOwner (6.52s)
```

⇒ **红（0 次）**。日志：`mut-m1-red.log`。**恢复（真文件未动）⇒ 绿 6.59s**（见上）。

### M2：把 `EnterWarm` 的 `mute()` 挪到落 `StageWarm` **之后**（顺序反了）

替换件 `m2-scope.go`，overlay `overlay-m2.json` 指向 `internal/session/scope.go`。

```
    resident_voice_wire_r1_windows_test.go:778: the first mic close saw stage "Warm", want Awake: the Warm row must close the mic BEFORE the transition lands
--- FAIL: TestWireR1WarmIdleDisposesOnTheOwner (6.93s)
```

⇒ **红**，且**次数仍是 2**（没触发 `closes < 2`）——证明**新增的顺序钉子自己有牙**，不是靠"数不够"顺带变红的。
日志：`mut-m2-red.log`。

### M3：只删掉新加的**屏障**（其余修法保留）

替换件 `m3-wire_test.go`，overlay `overlay-m3.json` 指向测试文件本身。

```
    resident_voice_wire_r1_windows_test.go:758: mic close counted 1 times, want at least 2 (the Warm row and the cascade's defence in depth)
--- FAIL: TestWireR1WarmIdleDisposesOnTheOwner (7.19s)
```

⇒ **红（1 次）**，逐字复现原红句 —— 证明**屏障本身是承重的**，不是装饰。日志：`mut-m3-red.log`。

| 突变 | 改了什么 | 读数 | 结论 |
|---|---|---|---|
| M1 | 删 `micCloses++` | **红** `counted 0 times` 6.52s | 计数器是承重的 |
| M2 | `mute()` 挪到落状态之后 | **红** `saw stage "Warm"` 6.93s（次数仍 2） | 顺序钉子有独立牙 |
| M3 | 只删屏障 | **红** `counted 1 times` 7.19s | 屏障是承重的 |
| 恢复 | 真文件原样 | **绿** 6.59s | — |

---

## AC#2 越界检查

`git diff --name-only`（逐行，工作树全部改动，含**别的腿**的在飞件）：

```
.scratch/wisp/probes/291/mic/resident_audio_291_probe_windows_test.go
.scratch/wisp/probes/orch/20261112/txt/346-a2-dispatch.txt
cmd/wisp/resident_voice_wire_r1_windows_test.go
cmd/wisp/settle_trace_12r2_windows_test.go
frontend/fixtures/composer-states.html
```

**★ 我这一腿只改 1 个文件**：

```
cmd/wisp/resident_voice_wire_r1_windows_test.go | 44 ++++++++++++++++++++++---
1 file changed, 40 insertions(+), 4 deletions(-)
```

其余 4 个是**进腿前就已在工作树的别人的在飞件**（mtime 证明：`291_probe` = 11-09 22:01、`346-a2-dispatch` = 11-12 13:39、`settle_trace_12r2` = 11-12 20:10（票 427 的并发腿）、`composer-states.html` = 11-11 16:59）。我一枚都没碰。

禁改项逐条核：`docs/PLAN.md` **无** · `docs/specs/**` **无** · `internal/observe/thresholds.go` **无** · golden **无** · `tools/d22scan/allowlist.txt` **无** · `frontend/**`（`composer-states.html` 是**别人**的，非我改）· `design/**` **无** · `build/mut_rvw.go` **未删**（未跟踪残留，不是我的写面）。

---

## AC#3 门禁（四条读数）

| 门禁 | 命令 | 读数 |
|---|---|---|
| build | `go build ./cmd/wisp/` | **exit=0**，零输出 |
| vet | `go vet ./cmd/wisp/` | **exit=0**，零输出 |
| gofumpt | `D:\work\base\gopath\bin\gofumpt.exe -l cmd/wisp/` | **零输出**，exit=0 |
| d22scan | `D:\work\soft\Git\bin\bash.exe -c "scripts/d22scan.sh"` | **exit=0**，`d22scan: clean - no D22 ban violations`；`runtests.sh: OK - PASS=45 FAIL=0 SKIP=0` |

- 用 `go build ./cmd/wisp/` 而非 `./...`（票面明令；`build/mut_rvw.go` 未跟踪残留会让 `./...` 红，不是我的写面）。
- `scripts/d22scan.sh` 是唯一受支持的入口（`tools/d22scan` 是**独立 module**，`go run ./tools/d22scan` 从根模块跑会失败——我实测撞到 `main module does not contain package`，改用 wrapper 后通过）。

---

## AC#4 与票 28 的关系

- 本票**只处理 `:744` 这一格**，⛔ **未翻票 28 任何框**，⛔ 未改票 28 票面。
- **需回填票 28 票面**：`R28-2` 的语义（"mic closed"＝唤醒闩，不是设备/喂食门）**本次取证完全支持**，且票 28 `:90` 已记它由票 41 注销。**建议编排者**把 `tri-v1-verdict.md:123` 那句"必须写在 `Warm` 定义旁"真正落到 `28-…-done.md:20` 旁（本腿按票面禁令**不自行改票 28**）。

---

## 一句话结论

> **甲形（断言过严）**：承诺两次是对的，产码也给两次（`scope.go:470` 一次 `Awake`、`scope.go:691` 一次 `Disposed`），
> 但断言在 `StageDisposed` 一可读（`scope.go:654`，**早于**级联的 `scope.go:691`）时就读计数器 ⇒ 读早了，只读到 1 次。
> 修法是**只改测试**：加 owner 队列屏障等完级联，并把"两次"钉成"`[Awake … Disposed]`"这个顺序形状。
