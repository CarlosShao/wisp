# 432 — 三枚 `Test258*` 首次真跑即红：判「甲（环境件）」，修台面

- 腿：`432-hotkey258`（写腿）
- 日期：2026-11-12
- 票：`.scratch/wisp/issues/432-...-premises-moved.md`
- 基线 HEAD：`73bcedd5b`（`dev`）
- 一句话结论：**判甲（环境件），修台面一处**。三枚要的是"Win32 把这几枚组合键交给本进程"这个前提；round 89 那 6 秒里，**另一条腿的常驻进程正占着 `Ctrl+Alt+Q`／`Ctrl+Alt+M`／`Ctrl+Alt+P`**，前提不存在，红句说的是机器不是产品。产码一字未改。

---

## 0. 先说三件会改变结论的事（都是现跑读数，不是转述）

### 0.1 ★票面「round 88 只跑了四枚」这个前提**不成立**——它跑了十一枚且全绿

票面 `## ★历史` 节逐字写「round 88 `Test258*` 只跑了**四枚**且**全 PASS**」「这三枚在 round 88 压根没跑到」。现读 round 88 的整包日志，**它跑了 11 枚 `Test258`，11 枚全 PASS**，三枚在册：

```
D:\tmp\wisp\20261112\fullclean\cmdwisp.log   (LastWriteTime 2026-11-12 19:09:19)
=== RUN   Test258BridgeMutationNoSrcKeepsOldBinding
--- PASS: Test258BridgeMutationNoSrcKeepsOldBinding (1.07s)
=== RUN   Test258OccupiedCombinationNamesTheNewValue
--- PASS: Test258OccupiedCombinationNamesTheNewValue (0.77s)
=== RUN   Test258V1ProbeSummonEditRebindsLiveBall
--- PASS: Test258V1ProbeSummonEditRebindsLiveBall (1.01s)
```

尺：`Select-String -Path ... -Pattern '=== RUN   Test258'` ⇒ **11**；`'--- PASS: Test258'` ⇒ **11**；`'--- FAIL: Test258'` ⇒ **0**。

⇒ **所以这三枚不是"第一次真跑"**：round 88（19:09）真跑过且绿，round 89（21:23）同一批码真跑过且红。
**两次之间产品侧没有任何提交**（`git log` 见 §2，最近一笔动 `cmd/wisp/resident_hotkey*` 的是 2026-11-08 的 `eec09a028`，两轮都在它之后）。
**变了的是机器**。这一点单独就把「乙 产品退化」排除掉了——同一个二进制、同一批码，19:09 绿 21:23 红，中间只有桌面状态变了。

> 这张票的"首次真跑"说法我**没有**照抄。它若成立，红就只能读成"新退化"；它不成立，红就必须先问"机器怎么了"。这一格是本票最关键的一次纠偏。

### 0.2 round 89 那 6 秒里，Q／M／P 三枚被另一个进程占着（32 条占用日志）

`D:\tmp\wisp\20261112\round89\full.log` 现数：

| 读数 | 值 |
|---|---|
| `hotkey occupied by another program, not registered` 条数 | **32** |
| 其中 `Ctrl+Alt+M` | 12 |
| 其中 `Ctrl+Alt+P` | 12 |
| 其中 `Ctrl+Alt+Q` | 7 |
| 其中 `Ctrl+Alt+V` | 1（这是 ② 自己造的 squat，属预期） |
| `hotkeys_live=3` 的球 | 5 枚 |
| `hotkeys_live=1` 的球 | 4 枚（**全是这三枚票内用例**） |

红句前后的逐字原文（`full.log` 第 4077 行起）：

```
=== RUN   Test258BridgeMutationNoSrcKeepsOldBinding
level=ERROR msg="hotkey occupied by another program, not registered" hotkey=summon binding=Ctrl+Alt+Q err="Hot key is already registered."
level=ERROR msg="hotkey occupied by another program, not registered" hotkey=mute binding=Ctrl+Alt+M err="Hot key is already registered."
level=INFO  msg="cancel hotkey left unbound while idle: ..." hotkey=cancel binding=Esc
level=ERROR msg="hotkey occupied by another program, not registered" hotkey=panel binding=Ctrl+Alt+P err="Hot key is already registered."
level=INFO  msg="ball: the resident leg created the floating ball window" hotkeys_provenance=defaults hotkeys_live=0 hotkeys="summon=Ctrl+Alt+Q occupied by another app, mute=Ctrl+Alt+M occupied by another app, ..., panel=Ctrl+Alt+P occupied by another app"
    resident_hotkey_258_windows_test.go:361: the old binding did not survive a no-rebind host (VK=0x0): the premise moved
```

⇒ `VK=0x0` 不是"注册返回了 0 号键"，是**这个槽根本没注册上**（`Live()` 里没有它，map 取零值）。`live slots=1` 同理：Z 是我们自己的、M 与 P 被别人占着、cancel 按 245 是 standby ⇒ 只剩 1。

### 0.3 占位的是**另一条腿的常驻进程**，不是"机主在用电脑"

`D:\tmp\wisp\20261112\431-ac1budget\boot1\logs\wisp-20261112-001.jsonl` 现读：

```
{"time":"2026-11-12T21:22:48.5225505+08:00", ... "hotkey":"cancel","binding":"Esc"}
{"time":"2026-11-12T21:22:49.4376617+08:00", "msg":"ball: the resident leg created the floating ball window","hotkeys_provenance":"defaults","hotkeys_live":3,
 "hotkeys":"summon=Ctrl+Alt+Q live, mute=Ctrl+Alt+M live, cancel=Esc not bound while idle ..., panel=Ctrl+Alt+P live"}
```

**21:22:48 起，431 腿的常驻进程把 Q／M／P 三枚全部 live 着**；round 89 的三枚用例在 **21:23:32–21:23:38** 跑（`full.log` 时间戳逐字）。同一个桌面、同一时刻、同一批键。

⇒ 这不是"有人在用电脑导致前台被抢"（那是另一族 `TestAC4Focus*`／`TestAC246*` 的病，本票明令不碰），而是**同机并行的另一条写腿留了一个常驻 Wisp 进程**。这一点决定了修法：不是"换冷门键躲开机主"，是"检测到占用就**大声跳过并写清原因**"——因为占位者可能是任何东西，换键只是把赌注换个地方。

---

## 1. ★特别确认：这三枚要不要独占桌面前台／热键？——**要热键，不要前台**

这是票面点名必须明确回答的一格，分两层答：

### 1.1 要「热键注册」（`RegisterHotKey`），这是硬依赖

三枚都调 `startResidentBallBounded` → `startResidentBall` → 真起一个球窗 → `registerAll` 真调 `RegisterHotKey`。
`RegisterHotKey` 是**桌面级、跨进程、与焦点无关**的：一个组合键全桌面只能有一个持有者。所以：

- **别的程序占着那枚键 ⇒ 必红**（本票 round 89 就是这个形状）；
- 这条依赖**不是**"焦点在谁那儿"的问题，**前台窗口是别人的也照样能注册**（我 §4 的对照组就是证：holder 只是持有键、不抢前台，三枚照样被逼到 skip）。

### 1.2 不要「独占桌面前台」（`SetForegroundWindow`／`SetFocus`）

三枚的断言全部读三样东西，没有一样是前台：

| 断言 | 读的是什么 | 前台依赖 |
|---|---|---|
| `:361` 旧绑定存活 | `rb.b.RegisteredHotkeys()[hkSummon258].VK`（Win32 注册集） | 无 |
| `:447` `hotkeys live 3/4` | `rb.verdict`（启动时算好的 `len(rep.Live())`） | 无 |
| `:452` 活槽＝2 | `rep.Live()`（Win32 注册集） | 无 |
| `v1probe:107` `hotkeys live 3/4` | `rb.verdict` | 无 |

尺：`git grep -n 'SetForegroundWindow\|GetForegroundWindow\|SetFocus' -- 'cmd/wisp/*.go' 'internal/ball/*.go'` 的命中**全在 `panel_*`／`*winlive*` 族**，`resident_hotkey_258_windows_test.go` 与 `resident_hotkey_v1probe_test.go` **零命中**。

### 1.3 ⇒ 直接回答票面那句

> 「本机有人在用电脑时是否必然红？」

**不是"必然红"，是"取决于那个人／那个程序有没有占住 Q／M／P"**。分三种：

1. 有人在用电脑、但**没占这三枚键** ⇒ **照绿**（我 §4 的 run1／run2／run4／run5 就是：机主当时在用这台机，四发全绿）。
2. 有程序**占住其中任何一枚** ⇒ **必红**（round 89；我 §4 的 M1 也复现了）。
3. 有程序占住、而**本票的台面修复在盘上** ⇒ **必 skip 且逐字说明**（我 §4 的 mut-control-1）。

所以准确说法是：**这三枚依赖「热键注册」，不依赖「前台」；占用即红，与"有人在用电脑"本身无必然关系。**

---

## 2. 最近动过 `cmd/wisp/resident_hotkey*` 的提交（AC#0 第一件事）

`git log --oneline -10 -- 'cmd/wisp/resident_hotkey*'` 现跑：

| commit | 日期 | 一句话 |
|---|---|---|
| `eec09a028` | **2026-11-08 18:53** | 票408 r2：整包真窗/热键活体测试无限挂起第二根因修复（有界启动等待 `startResidentBallBounded`） |
| `9e84ae42e` | 2026-11-05 18:14 | 245-r3：稳态热键钉改判「裸键」＋`wisp -h` 半假话改条件句＋补 cancel 配置那枚钉 |
| `28a2ff4ea` | 2026-10-09 20:25 | 296-r1 AC#1 夹具先行（产码面只提取、零行为变化） |
| `cbece45ea` | 2026-10-04 09:20 | 258-r2 写码：档位词跟着集合说真话＋两枚 winlive 尺子修形状 |
| `d9bce7ce6` | 2026-10-03 22:01 | 258-v1 残件入库：v1 探针文件＋9 枚读数档 |
| `bd5c049f3` | 2026-10-03 22:00 | 258-v1 验收补 commit（6 枚主判据＋4 枚 winlive 尺） |

**关键读数**：最近一笔 `eec09a028` 是 **11-08**，round 88（11-12 19:09）与 round 89（11-12 21:23）**都在它之后**，两轮之间**零提交**。
⇒ 乙（产品退化）在时间线上就不成立：**同一个 HEAD、同一批码，19:09 绿、21:23 红。**

---

## 3. 甲／乙／丙 的判定，逐条给 `文件:行`

### 判：**甲（环境件）** —— 修台面一处

| 岔路 | 判 | 判据（`文件:行` ＋ 现跑读数） |
|---|---|---|
| **甲 环境件** | ✅ **就是它** | `internal/ball/hotkey_windows.go:504-509`：`reg()` 回 `ERROR_HOTKEY_ALREADY_REGISTERED` ⇒ `b.Status = HotkeyTaken`。round 89 逐字 `err="Hot key is already registered."`（32 条），占位者是 431 腿 21:22:48 起的常驻进程（§0.3）。⇒ 前提（键归我们）不存在。 |
| **乙 产品退化** | ❌ 排除 | ①时间线：§2，round 88／89 之间零提交。②同一批码 round 88 三枚全 PASS（§0.1）。③`cbece45ea` 之后 `hotkey_windows.go` 的最近改动是 2026-11-08 的 `eec09a028`（且只动测试层有界等待），两轮都在其后。④**产品行为在 round 89 是正确的**：它老老实实报了 `occupied by another app`、`hotkeys_live=0`／`1`，这正是 `hotkey_windows.go:394-398` `Problems()` 该说的话。坏的不是产品，是"拿这台机器的读数当产品的读数"这件事。 |
| **丙 断言过严** | ❌ 排除 | ①`:447`／`v1probe:107` 要 `hotkeys live 3/4` 与产码一致：`cmd/wisp/resident_ball_windows.go:526-527` 印 `len(rep.Live())`，而 idle 稳态按票 245 就是 **3**（cancel standby，`hotkey_windows.go:484-490` `registerAllWith` 对 `hkCancel` 走 `cancelIdleLine` 不注册）。②`:452` 要「活槽＝2」也与设计一致：`summon`(Z)＋`mute`(M) live、`cancel` standby、`panel` 被 squat 占 ⇒ **恰好 2**。这不是过严，是把设计算清楚了。③**反证**：round 88 与 round 89 之外的 run1／run2／run4／run5 四发（桌面干净时）**这个 2 每次都成立**（§4），若断言过严它不会在干净桌面上绿。 |

**一句话理由**：`internal/ball/hotkey_windows.go:504-509` 把"别人占着键"判成 `HotkeyTaken`，而 round 89 的 32 条占用日志证明当时 Q／M／P 真被别人占着（431 腿 21:22:48 起的常驻进程），所以三枚红在**前提缺失**，不在产品、也不在断言。

### 为什么**不**选"换冷门键"这条甲形修法

票面给了甲的两条修法：①换冷门键 ②检测到占用就 skip 并说明。我选 **②**，理由：

- 换键是**把赌注换个地方**：round 89 占位的是同机另一条腿的 Wisp 常驻进程，它按 `DefaultHotkeys()`（Q／M／P）绑——换到 F9 只是"这次没撞上"。下一发可能撞 F9，且撞上时红句还是同一句"前提动了"，**没人知道是机器**。
- ②把"机器状态"和"产品状态"**分开**：占用 ⇒ skip 并点名谁占着哪枚；其它任何失败（没尝试、别的 Win32 错、键活着但值不对）⇒ **照旧红**。这才是主流 harness 的做法（前提不成立时 skip-loud，断言失败时红）。
- 且 ①会**削弱判据**：三枚里 ② 那枚的整个立意就是"被占的组合键要被如实报出来"（票 258 AC#3），把它的键换成冷门键等于让它**更难**碰到被占形状。

---

## 4. 突变控制（真跑，四发读数）

对照组台件：`D:\tmp\wisp\20261112\432-hotkey258\holder\`（**仓外**，只用标准库，`holder.exe -keys Q,M,P` 在锁定线程上 `RegisterHotKey` 三枚并 `GetMessageW` 泵消息）。
它**只持键、不抢前台**——这本身就是 §1.2 那个"不要前台"论断的对照组。

| # | 台面状态 | 三枚读数 | 说明 |
|---|---|---|---|
| run1 | 修复前，桌面干净 | PASS 1.19s／PASS 0.81s／PASS 0.74s | 桌面干净时本来就绿 |
| run2 | 修复前，桌面干净 | PASS 1.19s／PASS 0.84s／PASS 0.78s | 复跑稳定 |
| run4 | **修复后**，桌面干净 | PASS 1.17s／PASS 0.78s／PASS 1.05s | 修复不改变干净桌面上的行为 |
| **mut-control-1** | **修复后 ＋ holder 持 Q/M/P** | **SKIP 1.31s／SKIP 0.74s／SKIP 0.92s** | 逐字见下 |
| **mut-M1-red** | **修复删掉 ＋ holder 持 Q/M/P** | **FAIL 1.11s／FAIL 0.69s／FAIL 0.79s** | 逐字见下 |
| mut-M1-restored | 修复恢复 ＋ holder 持 Q/M/P | SKIP 0.68s／SKIP 0.73s／SKIP 0.74s | 恢复即回 skip |
| run5 | 修复后，holder 已杀 | PASS 1.24s／PASS 0.75s／PASS 0.75s | 终态绿 |

### 4.1 突变 M1 的读法（删掉修复 ⇒ 必须回到红）

**M1 定义**：把三处 `skipIfOccupied258(...)` **调用**注掉（helper 留着，证明红不是编译不过），holder 保持持有 Q／M／P。

红句逐字（`mut-M1-red.log`）——**与 round 89 的三句同形**：

```
resident_hotkey_258_windows_test.go:410: the old binding did not survive a no-rebind host (VK=0x0): the premise moved
--- FAIL: Test258BridgeMutationNoSrcKeepsOldBinding (1.11s)

resident_hotkey_258_windows_test.go:506: the boot verdict lost the live-count shape: "the floating ball window is up in this process (tray icon added, hotkeys from config, hotkeys live 1/4: summon=Ctrl+Alt+Z live, mute=Ctrl+Alt+M occupied by another app, cancel=Esc not bound while idle (cancel is borrowed only during Confirming), panel=Ctrl+Alt+P occupied by another app)"
resident_hotkey_258_windows_test.go:511: live slots after the rebind = 1, want 2 (summon+mute; cancel is standby, panel is taken)
--- FAIL: Test258OccupiedCombinationNamesTheNewValue (0.69s)

resident_hotkey_v1probe_test.go:115: 258-v1 PROBE: boot verdict lost the live-count shape: "... hotkeys live 1/4: summon=Ctrl+Alt+Z live, mute=Ctrl+Alt+M occupied by another app, ... panel=Ctrl+Alt+P occupied by another app)"
--- FAIL: Test258V1ProbeSummonEditRebindsLiveBall (0.79s)
```

⇒ **突变有牙**：把修好的那处拿掉，三枚**精确回到** round 89 的红（`:361`→`:410`、`:453`→`:511`、`v1probe:108`→`:115`，行号偏移是我加注释造成的，红句文本一字未变）。

### 4.2 修复生效的读法（skip 必须大声）

`mut-control-1.log` 逐字：

```
resident_hotkey_258_windows_test.go:397: SKIP-LOUD: the no-rebind control (it requires the boot Ctrl+Alt+Q to be registered) needs Win32 to hand THIS process its own combinations, but another program on this machine already owns summon=Ctrl+Alt+Q (RegisterHotKey answered ERROR_HOTKEY_ALREADY_REGISTERED). The premise this case judges is missing, not the product: a second Wisp resident, or any other app holding the key, produces exactly this reading. Re-run where those combinations are free (or stop the holder) to get a judging run.
--- SKIP: Test258BridgeMutationNoSrcKeepsOldBinding (1.31s)

resident_hotkey_258_windows_test.go:446: SKIP-LOUD: ... already owns mute=Ctrl+Alt+M, panel=Ctrl+Alt+P (...)
--- SKIP: Test258OccupiedCombinationNamesTheNewValue (0.74s)

resident_hotkey_v1probe_test.go:71: SKIP-LOUD: ... already owns mute=Ctrl+Alt+M, panel=Ctrl+Alt+P (...)
--- SKIP: Test258V1ProbeSummonEditRebindsLiveBall (0.92s)
```

⇒ **零静默跳过**：三条 skip 都点名了**谁**（哪个槽）**哪枚键**、**为什么**（`ERROR_HOTKEY_ALREADY_REGISTERED`）、以及**怎么拿回一发能判的运行**。

---

## 5. 改了什么（只改台面一处，产码零改）

`git diff --stat`：

```
 cmd/wisp/resident_hotkey_258_windows_test.go | 62 +++++++++++++++++++++++++++-
 cmd/wisp/resident_hotkey_v1probe_test.go     |  7 ++++
 2 files changed, 67 insertions(+), 2 deletions(-)
```

1. `cmd/wisp/resident_hotkey_258_windows_test.go:71-76` —— 补 `hkMute258 uint32 = 2` 常量（原来只有 summon/panel 两枚；`hkMute` 在 `internal/ball/hotkey_windows.go:37`）。
2. `cmd/wisp/resident_hotkey_258_windows_test.go:78-119` —— 新增 `skipIfOccupied258`，**只**对 `ball.HotkeyTaken` 生效。
3. `cmd/wisp/resident_hotkey_258_windows_test.go:397`／`:446` 与 `cmd/wisp/resident_hotkey_v1probe_test.go:71` —— 三处调用，都放在**原有断言之前、原有 FATAL 之后**（`hotkeyBridgeArmed` 那两枚 fatal 仍在 skip 之前，所以接线退化照样红，不会被 skip 吞掉）。

**为什么这个 skip 不是"放宽断言"**：

- 它只在 `Status == ball.HotkeyTaken` 时触发，即 `hotkey_windows.go:506` 那条"Win32 自己答了 1409"的分支；
- `HotkeyDisabled`（没尝试）、`HotkeyUnparsable`、`HotkeyError`（别的 Win32 错）、`HotkeyLive` 但键值不对 —— **一律不 skip，照旧落到原断言上红**；
- 所以**产品退化躲不过它**：本票的 M1 已经证明"红得掉"，而 §3 的乙排除是独立的时间线证据。

---

## 6. 门禁四读数

| 门禁 | 命令 | 读数 |
|---|---|---|
| build | `go build ./cmd/wisp/` | **rc=0** |
| vet | `go vet ./cmd/wisp/` | **rc=0** |
| gofumpt | `D:\work\base\gopath\bin\gofumpt.exe -l <两枚改动文件>` | **零输出**，rc=0 |
| d22scan | `scripts/d22scan.sh`（**不是** `go run ./tools/d22scan`） | **rc=0**，`d22scan: clean - no D22 ban violations`；`runtests.sh: OK - packages=[./...] top-level: PASS=45 FAIL=0 SKIP=0` |

⚠ `build/mut_rvw.go` **未删**（票面明令）。

---

## 7. 越界检查（AC#2）

`git diff --name-only` 全量（含别人在飞的改动，本腿只认领前两行）：

```
.scratch/wisp/issues/426-...md                          <- 别人
.scratch/wisp/issues/427-...md                          <- 别人
.scratch/wisp/issues/428-...md                          <- 别人
.scratch/wisp/issues/429-...md                          <- 别人
.scratch/wisp/issues/430-...md                          <- 别人
.scratch/wisp/probes/291/mic/resident_audio_291_probe_windows_test.go   <- 别人
.scratch/wisp/probes/orch/20261112/txt/346-a2-dispatch.txt              <- 别人
cmd/wisp/resident_hotkey_258_windows_test.go            <- ★本腿
cmd/wisp/resident_hotkey_v1probe_test.go                <- ★本腿
cmd/wisp/resident_sink_nail_127_windows_test.go         <- 别人（431 腿，21:58 还在改）
frontend/fixtures/composer-states.html                  <- 别人
```

**本腿只碰上面标 ★ 的两枚。** 禁改清单逐条对照：`docs/PLAN.md` 未现、`docs/specs/**` 未现、`internal/observe/thresholds.go` 未现、golden 未现、`tools/d22scan/allowlist.txt` 未现、`frontend/**` 与 `design/**` 在**本腿名下**未现（`frontend/fixtures/composer-states.html` 是别人在飞的件，我没碰）。

### 7.1 一次并发事故，如实记

21:58:04，431 腿正在改 `cmd/wisp/resident_sink_nail_127_windows_test.go`（共享编译面），中途我这一发 carrier 撞上它的半成品，`go test -c` 报：

```
cmd\wisp\resident_sink_nail_127_windows_test.go:672:30: undefined: waitForFact127
cmd\wisp\resident_sink_nail_127_windows_test.go:675:24: undefined: waitForFact127
```

这不是我的改动引起的（我的两枚文件当时已 gofumpt/vet/build 全绿），是**同包别人写到一半**。等它补上 `waitForFact127` 后我复跑即绿（run4）。按项目铁律"`cmd/wisp/**` 每波只许一枚写腿"，这次撞车是排程问题，**我只记录、不裁定**，留给编排者。

---

## 8. 本票**没有**做的事（免得被读成做了）

1. **没改任何产码**：`internal/ball/**`、`cmd/wisp/resident_ball_windows.go` 一字未动。
2. **没碰** `TestAC4Focus*`／`TestAC246*`／`TestAC1ResidentLeg*`（票面明令）。
3. **没改** `thresholds.go`／SLO／golden。
4. **没翻**票 258／246／432 任何 AC 框；**没 push**。
5. **没"修好 round 89"**——round 89 的红是**机器状态**，我没有、也不该去让那台机器在那一刻变绿；我做的是让下一次遇到同样机器时**说清原因而不是谎报产品坏了**。

---

## 9. 给下一程的两句实话

1. **round 88 的"只跑四枚"是错的**，而它已经写进票面 `## ★历史` 并被当成"别当退化"的依据。**依据本身不成立**，虽然结论（不是退化）恰好对——我给出了真正的依据（两轮之间零提交＋同码 19:09 绿）。**建议编排者更正票面那一节**，否则下一个人还会照着错的前提推理。
2. **同机并行跑整包 + 跑真窗/真热键用例，本身就会互相污染**：round 89 的红不是孤例，32 条占用日志说明当时有 3 枚键被长期占着。真窗族要么独占，要么就该像本票这样**把占用变成一句可读的 skip**。
