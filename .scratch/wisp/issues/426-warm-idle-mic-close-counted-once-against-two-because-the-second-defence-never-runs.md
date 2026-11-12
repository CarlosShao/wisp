# 426 — `TestWireR1WarmIdleDisposesOnTheOwner` 恒红：`mic close` 只记 1 次而断言要 ≥2（Warm 行＋级联防御），且这条断言自 2026-10-26 起可能从未真跑过

- Status: **ready-for-agent**（写面＝`cmd/wisp/**` 与 `internal/voice/**`／`internal/session/**` 中**定位到的那一处**；⚠ 先只读取证再定写面，⛔ 不许先改后查）
- 来源：round 88 整包复跑（`D:\tmp\wisp\20261112\fullclean\cmdwisp.log`，`PASS=593 / FAIL=14`）＋我两次独立复跑（`D:\tmp\wisp\20261112\puncfix\run.log` 5.19s／`D:\tmp\wisp\20261112\warmidle\a.log` 6.81s）
- 关联：票 28（`Warm` 定义与 `R28-2`）· 票 26／27（`wire-r1` 同批装配腿 `d803344ae`）· 票 381r2／390r1（同一 `resident_voice_wire` 写面）· 台账 `A1799`／`A1801`

## 这是什么（人话）

有一段代码，在**进程从"醒着"降到"温着"**（Warm）的时候要**先把麦克风关掉**，然后才落状态。
测试想证明这个"先关麦"是真的发生过，办法是数**关麦钩子跑了几次**，要求**至少 2 次**：
一次是 Warm 那一行自己关的，一次是**级联拆除时的兜底防御**再关一次。

现在实测**只记 1 次**，所以红了。

## 现量（可逐字复跑）

**逐字红句**（`cmd/wisp/resident_voice_wire_r1_windows_test.go:745`）：

> `mic close counted 1 times, want at least 2 (the Warm row and the cascade's defence in depth)`

**我两次独立复跑**（不同时间、不同 carrier 调用）：
- `2026-11-12 18:5x`：**FAIL 5.19s**
- `2026-11-12 19:4x`：**FAIL 6.81s**，同一红句

⇒ **稳定真红，⛔ 偶发。**

## ★为什么这枚特别值得查（历史遮蔽）

- 该断言由 `d803344ae`（**2026-10-26 19:59:46**，`cmd/wisp 装配腿 wire-r1：26/27/28 三张票的同一个缺口一次接完`）加入。
- 2026-11-12 15:0x 那轮整包（`230ac5-self/full.log`）里它 **FAIL 0.01s**，但红因是
  `resident_voice_wire_r1_windows_test.go:697: install the streaming ASR into …`——
  **`WISP_SPIKE_MODELS is not set` 环境遮蔽，测试根本没跑到 `:745` 那条断言**。
- 本轮环境三层打通后（见 `docs/evidence/s1/round88-summary.md`），`:745` **第一次真正被执行**。
- ⇒ **这条断言可能自 2026-10-26 起就没绿过**，只是被环境一直遮着。⛔ 但**这只是假设，必须由你取证证实或推翻**。

## 机制链（现读，供你起步；⛔ 别照抄结论）

- `micCloses` 字段：`cmd/wisp/resident_voice_engines_windows.go:313`
- 自增点：**唯一一处** `:497`，在 `MicClose` 钩子内（`:488-504`）：
  - `:489-491` 注释逐字：「Read the stage HERE: EnterWarm calls this hook before it lands the transition, so a hook that sees Awake is the proof that "the mic was closed FIRST" is an order and not a sentence.」
  - 即：**钩子在 `EnterWarm` 落状态之前被调用**。
- 读数：`micCloseCount()`（`:739-746`）、`micCloseStages()`（`:748-755`）。
- 测试断言三处：`:471`（`closes < 1`）、`:744`（`got < 2`）、`:750-751`（`t.Logf` 打 `mic_closes=` 与 `stages_at_mic_close=`）。
- 另一处 `:740` 断言 `asr.closes.Load() >= 1`（引擎释放），与本格无关。

## ⚠ 一个必须自己判的岔路（⛔ 不许我替你定）

`docs/evidence/s1/tri-v1-verdict.md:123` 逐字：

> **28：可结**——`AC3` 归 59；其余 5 格可翻，但**票面必须把 `R28-2` 写在 `Warm` 定义旁**（"mic closed"＝**唤醒闩**，不是设备/喂食门）

⇒ **「mic closed」在本仓的语义是"唤醒闩"，⛔ 不是"关设备"。** 你要先判清：
- 这枚测试数的 `MicClose` 钩子，是**唤醒闩**那一层，还是**设备关闭**那一层？
- 断言里的「the cascade's defence in depth」指的是**谁**的兜底？（去读 `dispose` 级联与 `MicClose` 的关系）
- **如果是"设计上只该关 1 次"⇒ 这条断言过严，属于「断言落后于产码」**（同票 425 那一族）；
- **如果是"设计上确实该关 2 次而第二次没发生"⇒ 这是产品缺陷**，要修产码。
**这两种结论的修法完全相反。⛔ 不许在没判清之前动任何一行。**

## 要做的（按序）

- [ ] **AC#0 只读取证（⛔ 本格不许改任何文件）**：
  1. 用 carrier 复现，贴逐字红句与耗时（**必须用 `D:\work\soft\Git\bin\bash.exe`**；本机 `bash` 是 WSL 转发器，一枚测试都不跑）：
     ```
     D:\work\soft\Git\bin\bash.exe -c "scripts/wisp-cli-staged-tests.sh '-test.run' '^TestWireR1WarmIdleDisposesOnTheOwner$'"
     ```
  2. 把 `:750` 那行 `WARM-IDLE` 的 `t.Logf` 输出**逐字贴出来**（含 `mic_closes=` 与 `stages_at_mic_close=`）——`stages_at_mic_close` 里若出现 `Awake` 就证明"先关麦"那一半成立，缺的是第二次。
  3. 读 `EnterWarm`／`dispose` 级联／`MicClose` 钩子三处产码，画出**钩子实际会被调用几次、分别在哪条路径上**。
  4. 读 `docs/evidence/s1/tri-v1-verdict.md:118` 与 `:123` 两条，判 `R28-2` 的语义。
  5. 读票 28 的 `Warm` 定义段。
  **产出：一句话判定「断言过严」还是「产品缺陷」，并给出支撑它的 `文件:行` 与逐字读数。**

- [ ] **AC#1 按 AC#0 的判定修（两形互斥，只做被判定那一形）**：
  - **若判「断言过严」**：改 `:744-746` 的断言，⛔ 不许放宽成"≥1 就算过"（那等于删掉这条测试）——要改成**钉住真正的承诺**（例如"Warm 那一行必须关麦一次且 `stages_at_mic_close` 含 `Awake`，级联若关则必须发生在 Warm 之后"）。**必须配突变控制**：把 `:497` 的 `micCloses++` 删掉 ⇒ 新断言必须红。
  - **若判「产品缺陷」**：修产码让级联防御真的关第二次麦。**必须配突变控制**：把修好的那行删掉 ⇒ 断言必须回到红。
  - ⚠ 无论哪形，**突变控制必须真跑**（`-overlay` 拷贝放**仓外**），⛔ 不许只写在注释里。

- [ ] **AC#2 越界检查**：`git diff --name-only` 逐行贴。出现 `docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**` ⇒ **直接退回**。

- [ ] **AC#3 门禁**：`go build ./cmd/wisp/` · `go vet ./cmd/wisp/` · `gofumpt -l` 零输出 · `tools/d22scan` 不新增违规。四条读数全贴。
  ⚠ 已知残留：`build/mut_rvw.go` 是**未跟踪、被 `.gitignore` 吃掉的残留件**（mtime 2026-11-04），会让 `go build ./...` 红。**⛔ 不是你的写面，⛔ 不许删**；用 `go build ./cmd/wisp/` 而非 `./...`。

- [ ] **AC#4 与票 28 的关系写清**：本票只处理 `:744` 这一格，⛔ 不翻票 28 任何框；若判定涉及 `R28-2` 语义，只在证据件里写明「需回填票 28 票面」，⛔ 不自行改票 28。

## 本票**不**解决

- 不处理另外两枚 `TestAC1ResidentLeg*`（那两枚是预算/偶发，另计）。
- 不碰 `Test26P*`／`Test384*`／`Test12r2*`／`Test365*` 那四枚。
- 不改 `docs/PLAN.md` 与 `docs/specs/**` 一字。

## 派单约束

- 写面：先只读取证（AC#0 零写），判定后再定；**产码与测试二选一，⛔ 两边都改**（否则突变控制证明不了任何事）。
- ⛔ `git add -A`／`add .`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；`add`／`commit` 两处显式 pathspec。
- 只 commit，⛔ push（推送由编排者做）。
- ⛔ 翻任何 AC 框（翻框由编排者做，非实现者）。
- 临时件只建不删，落 `D:\tmp\wisp\<YYYYMMDD>\426-<腿码>\`。

## Progress log

### 2026-11-12 20:0x–20:1x +08 — 腿 `426-warmidle`（⛔ 未翻任何 AC 框）

**判定＝甲形（断言过严），根因是测量竞态，不是产品缺陷。**
一句话：`micCloseCount() >= 2` 这条承诺本身是对的（级联确实该关第二次），但测试在 `StageDisposed` 一可读时就读计数器，
而 `StageDisposed` 由 `beginDisposeLocked` 在 `internal/session/scope.go:654` **先落**、级联的 `MicClose` 在 `internal/session/scope.go:691` **后跑** ⇒ 读早了，只读到 Warm 行那一次。

- **红发（carrier）**：`resident_voice_wire_r1_windows_test.go:745: mic close counted 1 times, want at least 2 …` / **FAIL 6.90s**，与票面 5.19s／6.81s 同句同数。
- **★决定性读数**：仓外 overlay 探针在断言前插 owner 队列屏障（等完级联）⇒ **`mic_closes=2`、`stages_at_mic_close=[Awake Disposed]`、`disposals=1`、PASS 17.17s`**，**产码一行未改** ⇒ 第二次关麦从未缺失。
- **铁证（同一份日志时间戳）**：`:745` 红在 `T20:03:03.177`，`SCOPE-DISPOSED`（`scope.go:693`，排在 `:691` 的 `MicClose` 之后）在 `T20:03:03.228` ⇒ 断言时级联还没跑完。
- **三处产码**：`EnterWarm` 在落 `StageWarm`（`:474`）**之前**调 `MicClose`（`scope.go:470`，看到 `Awake`）；`runCascade` 末尾再调一次（`scope.go:691`，看到 `Disposed`）；`closeEpoch`／`rv.stop()` 不走 `session.Hooks.MicClose` ⇒ **生产里恰好两条路径、各一次**。
- **`R28-2` 语义**：`tri-v1-verdict.md:118`＋`:123` 逐字"mic closed＝唤醒闩，不是设备/喂食门"；`MicClose` 钩子在本仓关的是 **FEED**（`resident_voice_engines_windows.go:336-343`、`resident_voice_windows.go:738-745`）⇒ **断言数的正是唤醒闩，语义没错**。票 28 `:20`／`:42` 的 `Warm` 定义要求"关麦 FIRST"，两次的承诺与设计一致。
- **修形：只改测试**（产码零改动）。**未放宽**成"≥1 就算过"：仍要求 `closes >= 2`，并新增 `len(stages) >= 2` ＋ 两条顺序钉子 `stages[0]=="Awake"`、`stages[last]=="Disposed"`。
- **突变控制三发（全经 `go test -overlay`，拷贝全在仓外 `D:\tmp\wisp\20261112\426-warmidle\`）**：
  - **M1** 删 `:497` 的 `micCloses++` ⇒ **红** `counted 0 times` 6.52s；恢复 ⇒ **绿** 6.59s。
  - **M2** 把 `EnterWarm` 的 `mute()` 挪到落状态**之后** ⇒ **红** `the first mic close saw stage "Warm", want Awake` 6.93s，**次数仍 2**（证明顺序钉子有独立牙）。
  - **M3** 只删新加的屏障 ⇒ **红** `counted 1 times` 7.19s，逐字复现原红句（证明屏障承重）。
- **门禁**：`go build ./cmd/wisp/` exit=0 · `go vet ./cmd/wisp/` exit=0 · `gofumpt -l cmd/wisp/` 零输出 · `scripts/d22scan.sh` exit=0（`clean - no D22 ban violations`，`runtests.sh: OK PASS=45 FAIL=0 SKIP=0`）。
- **越界**：只改 `cmd/wisp/resident_voice_wire_r1_windows_test.go`（+40／−4）＋ 本证据件 ＋ 本 Progress 条。未删 `build/mut_rvw.go`，未动 `PLAN.md`／`specs/**`／`thresholds.go`／golden／`allowlist.txt`／`design/**`。
- **证据件**：`docs/evidence/s1/426-warmidle-mic-close.md`
- **AC#4 需编排者回填**：建议把 `tri-v1-verdict.md:123` 那句"必须写在 `Warm` 定义旁"落到 `28-session-scope-warm-done.md:20` 旁（本腿按禁令**未自行改票 28**）。
