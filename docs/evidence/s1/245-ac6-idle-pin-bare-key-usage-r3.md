# 245-r3 证据：AC#6② 稳态钉子改判「裸键」＋ AC#6① `wisp -h` 那句半假话 ＋ AC#4 缺的那枚钉子

- **票**：`.scratch/wisp/issues/245-the-bare-esc-cancel-default-is-a-global-hotkey-so-the-resident-process-steals-esc-from-every-other-app.md`（无 `-done`、无撤回、正文 `Status` 仍是「已立，未派」）
- **腿**：`245-r3`（产码腿；**不是**裁决腿 ⇒ 本件只是〔实现方自述＋台件〕，AC 框一枚未翻）
- **起手 HEAD**：`ff61b9dbca56`（`dev`）
- **写面（全部改动 6 个文件，逐名）**：`cmd/wisp/main.go`（只有 `usage` 字符串）／`cmd/wisp/resident_hotkey_cancel_245_test.go`（新）／`internal/ball/hotkey_idle_bare_key_245_test.go`（新）／`internal/ball/hotkey_status_test.go`／`internal/ball/hotkey_live_test.go`／`internal/ball/live_windows_test.go`
- **产码行为零改动**：`internal/**` 与 `cmd/wisp/**` 的**非测试**改面只有 `cmd/wisp/main.go` 的 `usage` 常量一个字面量（`git diff --stat` 读数见 §5）。守则未动：`registerAllWith` 的 cancel 跳过、`takeEscWithAcc`／`releaseEscWith`、`HotkeyStatus` 枚举、`DefaultHotkeys()` 的 `Cancel: "Esc"` 一字未动（⇒ 未触碰 D43/C12 冻结）。
- **为什么是这三格**：票面未闭 = AC#4（阻塞于票 228）／AC#6①（`cmd/wisp` 文案，出上一腿写面）／AC#6②（钉子精度账）／AC#8・AC#9（真机）。台账 `A1534`（2026-11-05 10:2x）逐字：「★未闭：AC#4 阻塞228／AC#6①② 出写面／AC#8/9 真机」。本腿取 **AC#6①②**（写面内、零真机）＋ AC#4 里**唯一还没钉住的那一跳**；AC#8／AC#9 是真机面，**本波不碰**（见 §6）。

---

## §0 起手读票与前提校正（我自己现读，不转抄）

| 事实 | 读数 | 尺 |
|---|---|---|
| 钢钉仍锚在「枚数」的四处 | 改前逐字：`internal/ball/hotkey_status_test.go:129` `if len(live) != 3 {`／`:181` `if len(rep.Live()) != 3 {`／`internal/ball/live_windows_test.go:74` `!got.AllLive() \|\| len(got.Live()) != 3`／`internal/ball/hotkey_live_test.go:498` `if len(rep.Live()) != 3 {` | `grep -n "len(.*Live()) != 3" internal/ball` |
| 票 260 已把「借哪一枚键」接到配置（形ⓐ），**但没有**动「稳态绑不绑」 | `internal/ball/hotkey_windows.go:458-465` 逐字 `if p.id == hkCancel {`…`continue` 仍在 | `sed -n '455,470p'` 我现读 |
| 票 258 已把 `[hotkey]` 整段接进常驻腿 | `cmd/wisp/resident_windows.go:370` 起两个闭包；`cmd/wisp/resident_ball_windows.go:317` `residentBallHotkeyChain258` | 我现读 |
| ⇒ **AC#4「阻塞于票 228」这个前提已经倒了** | 常驻腿今天**读** `c.Hotkey.Cancel`（`cmd/wisp/resident_hotkey_258_windows_test.go:105` 逐字 `Cancel: c.Hotkey.Cancel`），但**没有一枚钉子种非默认 cancel**：258 的正控 `:126` 逐字 `Cancel: "Esc"`（刻意留默认），`cmd/wisp/resident_hotkey_live_258_windows_test.go:82` 逐字 `cancel = \"Esc\"` | `grep -rn 'Cancel:' cmd/wisp/resident_hotkey_*` |

⇒ AC#4 剩下的不是「配置通路」（258 已通），而是**这一枚键没有任何钉子**：一条把用户 cancel 行悄悄丢回默认的链，今天全树绿。本腿补的就是这一枚（§3）。

---

## §1 AC#6②：稳态钉子从「枚数」改判「裸键」

### 1.1 新判据（`internal/ball/hotkey_idle_bare_key_245_test.go`，新文件，`//go:build windows`）

- `realModifiers245 = modAlt|modControl|modShift|modWin`（**刻意不含 `MOD_NOREPEAT`**：那一位是重复抑制，`{MOD_NOREPEAT, VK_ESCAPE}` 正是生产默认那枚**裸 Esc**）
- `idleBareKeyHeld245(rep) (HotkeyBinding, bool)`：注册集里第一枚 `Status==HotkeyLive` 且 `Acc.Mods&realModifiers245 == 0` 的绑定
- `idleRosterReds245(rep) []string`：**四层同一枚判据**——① 任何槽上有一枚**无修饰键**的活绑定 ⇒ 红（逐字点名槽位与键）；② `summon`／`mute`／`panel` 三槽必须各自 live（**逐槽判，不再用枚数**）；其余状态交给既有 `Problems()`／`AllLive()`

### 1.2 四处钉子逐枚改（对照表）

| 钉子 | 改前判据 | 改后判据 | 覆盖变化 |
|---|---|---|---|
| `hotkey_status_test.go:129`（`TestRegisterAllLiveSet`） | `len(live) != 3` | `idleRosterReds245(rep)` | **收紧**：多绑一枚**裸** F9（三槽之外看不见的那枚）改前绿、改后红；多绑一枚**带修饰键**的 cancel 改前红、改后不红（AC#6② 要求的方向） |
| `hotkey_status_test.go:181`（`TestDefaultHotkeysIdlePassHoldsNoEsc`） | `len(rep.Live()) != 3` | 同上 | 同上（这一格喂的是**产品默认**，`Cancel="Esc"` 是裸键 ⇒ 判据等价且更强） |
| `live_windows_test.go:74`（`TestBallLiveLifecycle`，winlive） | `!AllLive() \|\| len(Live()) != 3` | `!AllLive() \|\| len(idleRosterReds245(got)) != 0` | 同上；`AllLive()` 原样保留 |
| `hotkey_live_test.go:498`（`requireIdleRoster`，winlive） | `Live()[hkCancel]` 在场 + `len(Live()) != 3` + 三槽逐名 + `cancel.Status==Standby` | `idleRosterReds245(rep)`（**另保留**桌面侧 `escBorrowProbe`：裸 Esc 此刻真的没人占） | 见 1.3 覆盖搬迁图 |

### 1.3 覆盖搬迁图（**谁搬走了、搬去哪、有没有洞**，明写不藏）

- **搬走的那一条**：`requireIdleRoster` 不再断言「cancel 那枚**一定**不在注册集里」，也不再数枚数 ⇒ 一枚**用户配置的带修饰键** cancel 若在稳态绑上，winlive 层不再打红。
  - **为什么这是裁定要的方向**：票面 AC#6② 逐字「用户把 `[hotkey] cancel` 配成带修饰键的组合时**应当常绑**（那不会吞别的程序）…配成裸键时仍走借用」；「判据＝稳态注册集里不存在『无修饰键』的那一枚」。枚数锚会**在 AC#4 通了之后**把这种配置打红或被实现方为过钉而悄悄不绑（票面原话）。
  - **洞不存在，因为同一事实另有两层钉**：① `internal/ball/hotkey_status_test.go:119-127` 的**尝试计数**（`len(reg.attempted) != 3` ＋ 逐枚 `id == hkCancel` 即红）——稳态**根本不许试绑** cancel；② 票 260 的 `TestConfiguredCancelStillNeverBoundWhileIdle260`（`internal/ball/hotkey_cancel_borrow_260_test.go:292`）逐字 `if n := reg.attemptsOf(hkCancel); n != 0 { t.Fatalf("ticket 245 RED: …") }`。两层都**未改一字**。
- **新长出来的那一条（改前两层都没有）**：`summon`／`mute`／`panel` 任一被配成**裸键**（例：`summon = "F9"`）⇒ 稳态就是三枚、cancel 也在 standby，**枚数锚读到全绿**，而桌面已经少了一枚 F9。改后 `idleBareKeyHeld245` 直接点名它（§1.4 控制①的逐字读数）。

### 1.4 RED-before → GREEN-after（同一发用例，逐字读数）

用例：`TestAC6IdlePinJudgesTheBareKeyNotTheCount245`（四段：① 枚数盲区控制 ② 裸 Esc 必红 ③ 用户组合键必不红 ④ 产品默认不红）。

**RED（判据还是改前形状——只判三槽在不在）**：

```
--- FAIL: TestAC6IdlePinJudgesTheBareKeyNotTheCount245 (0.00s)
    hotkey_idle_bare_key_245_test.go:114: the idle-roster nail did not name the bare summon key: [] (report [{Name:summon ID:1 Binding:F9 Status:live Acc:{Mods:16384 VK:120} …
```
（`Acc:{Mods:16384 VK:120}` ＝ `MOD_NOREPEAT|F9`，无任何真修饰键；同一发报告里 `len(Live())==3`、cancel 是 standby ⇒ **改前的枚数锚在这一发上是绿的**，用例第一段自己断言了 `len(rep.Live()) != 3` 会 `Fatal`，即「这发红不是因为枚数变了」。）

**GREEN（判据落地后）**：

```
--- PASS: TestAC6IdlePinJudgesTheBareKeyNotTheCount245 (0.00s)
ok  	github.com/CarlosShao/wisp/internal/ball	0.044s
```

### 1.5 突变两发（证明这枚尺有牙，且两向都咬）+ 逐字节还原

| 突变 | 改法 | 红句逐字 | 还原 |
|---|---|---|---|
| **M1** 判据里的裸键红去掉 | `if b, held := idleBareKeyHeld245(rep); held {` → `…; false && held {` | `hotkey_idle_bare_key_245_test.go:119: the idle-roster nail did not name the bare summon key: [] (report [... Binding:F9 Status:live Acc:{Mods:16384 VK:120} …])` ⇒ `--- FAIL` | `cp` 回原字节，`diff -q` 报 `restore-bit-identical` |
| **M2** 把任何活绑定都当裸键（丢修饰键判） | `if b.Acc.Mods&realModifiers245 == 0 {` → `if true {` | `hotkey_idle_bare_key_245_test.go:132: a bare Esc held while idle did not red with a sentence that names it: [ticket 245 RED: the idle ball holds summon = "Ctrl+Alt+Q" as a bare key ({Mods:16387 VK:81}, no Ctrl/Alt/Shift/Win) …]` ⇒ `--- FAIL`（控制②先炸：带修饰键的 summon 被判裸） | 同上，`restore-bit-identical` |

**M2 的红句里那句 `{Mods:16387 VK:81}`**＝`MOD_NOREPEAT|MOD_CONTROL|MOD_ALT` + `Q`：带修饰键的组合键**不能**被判裸——这正是 AC#6② 要的那一向。

### 1.6 这一格**没**证明什么（明写）

- **winlive 真机没跑**（`-tags winlive` 需真桌面＋真键盘，本波派单禁）：`go vet -tags winlive ./internal/ball/ ./cmd/wisp/` rc=0 只证**编得动**，不证真窗口上的读数。⇒ `requireIdleRoster`／`live_windows_test.go` 两处改动属〔改后未实测〕，交裁决腿或下一波 winlive 窗口。
- **没有新增词面型仪器**：`idleBareKeyHeld245` 问的是注册集里的**加速器位**，不问文案、不扫注释。

---

## §2 AC#6①：`wisp -h` 那句「four global hot keys」是半假话

### 2.1 改前 / 改后（逐字）

改前（`cmd/wisp/main.go:24-31`，`HEAD` 版）：

```
  wisp             GUI resident process: boots the runtime skeleton, hosts the
                   floating ball window, its tray icon and its four global hot
                   keys in this same process, shows confirmation cards on that
                   ball, and takes tasks from the operator's console when this
                   process really owns one (ticket 246 AC#7: …
```

改后：

```
  wisp             GUI resident process: boots the runtime skeleton, hosts the
                   floating ball window, its tray icon and its three always-on
                   global hot keys in this same process (summon, mute, panel;
                   the cancel slot is the fourth key, bound only while a
                   confirmation card is waiting and handed back when it leaves,
                   so an idle wisp holds no bare key against the rest of the
                   desktop), shows confirmation cards on that ball, and takes
                   tasks from the operator's console when this process really
                   owns one (ticket 246 AC#7: …
```

⇒ 票面 AC#6① 要的「带条件的事实句（说清稳态三枚、进入 Confirming 借第四枚）」到位；**未改**键名（`Esc` 仍是否决键，D43 #22 不动）。

### 2.2 真值来源（不是文案自证）

| 事实 | 读数 | 尺 |
|---|---|---|
| 稳态注册集＝3：`summon/mute/panel`，cancel 是 standby | `TestDefaultHotkeysIdlePassHoldsNoEsc` PASS（同一发里 `reg.bindsEsc()==false` ＋ `Cancel != "Esc"` 即红） | `go test ./internal/ball -run TestDefaultHotkeysIdlePassHoldsNoEsc -v` |
| `usage` 只被一枚门读（形状级，不是词面级）：bare 行 `^  wisp {2,}\S` | `TestAC1AC2DispatchHopGate133` PASS（改后现跑） | 见 §4 门禁 |
| 常驻腿自己报的是 `hotkeys live %d/4`，`%d = len(rep.Live())` | `cmd/wisp/resident_ball_windows.go:476` 逐字 `hotkeys live %d/4: %s`，`len(rep.Live())` 现算 | 我现读 |

### 2.3 这一格的两个欠账（明写）

1. **票面 AC#6① 要「凭据＝现跑那发 winlive 读数」，本腿没跑**（真机禁令）⇒ 那句 `hotkeys live 3/4: …` 的**本波读数缺失**，只交了确定性层（`TestDefaultHotkeysIdlePassHoldsNoEsc` PASS ＋ `usageBare133` 门 PASS）。**不许读成「winlive 已核」**。
2. 词面钉子一枚都没新增（票面禁「词面型仪器」）；`usageBare133`（`cmd/wisp/leg_dispatch_gate_133_test.go:181`）的行首形状**保住**（`TestAC1AC2DispatchHopGate133` PASS 读数在 §4）。

---

## §3 AC#4：补上「用户改的 cancel 真的进常驻腿」这枚钉子

### 3.1 新钉子（`cmd/wisp/resident_hotkey_cancel_245_test.go`，新，`//go:build windows`，非 winlive）

`Test245AC4ResidentChainReadsAConfiguredCancel`：真 `config.toml`（`[hotkey] cancel = "Ctrl+Alt+V"`，其余三槽非默认）→ `config.LoadFile` → `loadHotkey258` 视图 → **`residentBallHotkeyChain258`**（生产链）→ `ball.HotkeyConfig`。

四条断言：
1. 夹具读数：视图答 `cancel = "Ctrl+Alt+V"`，且**不等于**产品默认 ⇒ 这一发不是「用户啥也没改」；
2. 链答 `Cancel == "Ctrl+Alt+V"`（不红句逐字点名 AC#4）；
3. `provenance == hotkeyProvenanceConfig`（不是回退句）；
4. **正控**：同一夹具**只删掉**用户那行 cancel ⇒ 同一条链必须**答回产品默认表**且 `provenance == defaults`（⇒ 断言 2 不是因为链无视输入而绿）。

### 3.2 它证到哪儿 / 证不到哪儿

- **证到**：常驻腿这一跳（文件 → 视图 → 链 → 球拿到的 `HotkeyConfig.Cancel`）**跟着用户改**。这是 AC#4 里唯一没有钉子的那一跳。
- **证不到（明写）**：`resolveCancelBorrow`（借用那遍到底注册哪一枚）**在 `internal/ball` 里是包内符号**，本包够不着；那一跳已由票 260 的 `hotkey_cancel_borrow_260_test.go:292` 与 `hotkey_cancel_borrow_expect_260r2_test.go` 钉住（本腿未改、复跑绿）。**真机端到端（改配置 → 界面上那句话真的变）** 仍〔未实测〕，与票 260 收口时那句欠账同一格。

### 3.3 突变（证明这枚新钉补的正是 258 正控看不见的那一格）

改法：`cmd/wisp/resident_ball_windows.go` 的链里加一行 `cfg.Cancel = ball.DefaultHotkeys().Cancel`（**用 `-overlay` 做，不动工作树**），跑两枚用例：

```
=== RUN   Test258ConstructionChainTakesConfigValues
--- PASS: Test258ConstructionChainTakesConfigValues (0.01s)
=== RUN   Test245AC4ResidentChainReadsAConfiguredCancel
    resident_hotkey_cancel_245_test.go:55: the resident chain handed the ball cancel = "Esc", want the user's
    planted value: a hand edit of [hotkey] cancel must reach the slot ticket 245 borrows (AC#4)
--- FAIL: Test245AC4ResidentChainReadsAConfiguredCancel (0.02s)
```

⇒ **258 那一枚照旧绿、只有新钉红**：这正是「链把用户 cancel 悄悄丢回默认，全树绿」的那个形状，改前它**没有钉子**。

---

## §4 门禁（本腿现跑，逐条读数）

| 门 | 命令 | 读数 |
|---|---|---|
| build | `go build ./cmd/... ./internal/... ./tools/...` | **rc=0** |
| `internal/ball` | `go test ./internal/ball -count=1 -timeout 420s` | **`ok … 0.528s`**（0 红；含 §1.4 那发） |
| 新增/相关用例 | `go test ./cmd/wisp -count=1 -timeout 420s -run "^(Test245AC4ResidentChainReadsAConfiguredCancel\|Test258ConstructionChainTakesConfigValues\|TestAC1AC2DispatchHopGate133)$" -v` | **3 PASS／0 FAIL**，`ok … 0.565s` |
| `cmd/wisp` 整包（**异常项，见下**） | `go test ./cmd/wisp -count=1 -timeout 900s` | **FAIL／45 红／695.026s**；**同一枚包**在 `-overlay` 把 `cmd/wisp/main.go` 换回 **HEAD 版**后跑**同一份 45 枚红名**＝**28.249s／45 红**，`diff` 名册 ⇒ **NAMES IDENTICAL (45/45)**（⇒ 本腿**新增红 = 0**） |
| gofumpt | `gofumpt -l cmd/wisp internal/ball`（`/d/work/base/gopath/bin/gofumpt.exe`） | **空** |
| d22scan | `sh scripts/d22scan.sh` | **rc=0**，`runtests.sh: OK … PASS=35 FAIL=0 SKIP=0`，`d22scan: clean`（cmd/ Go 208→**209**，恰是新测试文件那枚） |
| winlive 编译 | `go vet -tags winlive ./internal/ball/ ./cmd/wisp/` | **rc=0**（只证编得动，**不**证真机读数） |

**`cmd/wisp -timeout 420s` 这一发为什么要换成 900s（如实登记，不藏）**：
- 本机第一次按 `-timeout 420s` 跑整包，**在 `TestAC1ResidentLegInstallsItsLogListenerOnDisk` 卡死**（`cmd/wisp/secret_argv_windows_test.go:178 buildWispForTest` → `exec.Cmd.Run` 等子 `go build`；台账里这枚是既有「带载偶发」红，`pending-and-issues.md:12056`／`:23256` 记「HEAD 上单跑 3/3 也红」）。
- 整包正常时长本身就是 **695 s**（> 420 s）⇒ 420 s 这条门在本机**结构性跑不完**。改用 900 s 拿到完整读数，并用**同包 HEAD 基线**（overlay 只换 `main.go`）逐名对照，得到 45/45 同名 ⇒ 这 45 枚**全部是存量/环境**，与本腿无关。
- 这 45 枚的成因我抽读了三类：① **语音族**（`Test15e*`／`Test41a2r1*`／`Test26P*`／`Test392*`／`Test381r2*`／`Test390r1*`／`Test394R1*`／`TestWireR1*`／`Test382WR2*`／`Test26R3*`／`Test16r1*`／`Test12r2*`／`Test395r1*`）逐字红因＝`语音链未构造：VAD 没起来（speech: onnxruntime 1.17.1 is older than the 1.28 the linked sherpa-onnx C API requires; a second onnxruntime.dll is shadowing the one on PATH …）`（**DLL 影子**，环境）；② `TestAC4FocusReturnToPriorWindowGap33r5`＝本会话**没有前台桌面焦点**（逐字 `foreground … panel hwnd … 0x53c10a8`／`the panel did not take the foreground on Show`）；③ 两枚**他票既有红**：`TestTicket256ResidentGateOptionsFieldSetIsTheFiveItClaims`（`Options` 多了 `Answer` 字段，属票 256 名册过期）、`TestAC247HandingTheLevelToTheSeamIsNotVisibility`（`prototypeVisuals` 默认被别的票改过）。
- 另一发：`sh scripts/d22scan.sh` 第一次与整包测试**并发**跑时读到 `PASS=33 FAIL=2`，**独占重跑** `PASS=35 FAIL=0 rc=0` ⇒ 记「带载偶发」，读数以独占那发为准。

---

## §5 越界自查（只给名册，不给结论句）

```
$ git status --porcelain cmd/wisp internal/ball
 M cmd/wisp/main.go
 M internal/ball/hotkey_live_test.go
 M internal/ball/hotkey_status_test.go
 M internal/ball/live_windows_test.go
?? cmd/wisp/resident_hotkey_cancel_245_test.go
?? internal/ball/hotkey_idle_bare_key_245_test.go
$ git diff --stat HEAD -- cmd/wisp internal/ball
 cmd/wisp/main.go                    | 18 +++++++++++-------
 internal/ball/hotkey_live_test.go   | 29 ++++++++++++++---------------
 internal/ball/hotkey_status_test.go | 20 +++++++++++++-------
 internal/ball/live_windows_test.go  | 21 ++++++++++++---------
```

- 禁区逐名：`docs/PLAN.md` **0**、`docs/specs/**` **0**、`internal/observe/thresholds.go` **0**、`golden` **0**、`internal/risk/**` **0**、`frontend/**` **0**、`design/**` **0**、`tools/d22scan/allowlist.txt` **0**、三枚冻结件（`internal/panel/tokens_fourway_test.go`／`l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）**0**。
- 未新增 `DEFERRED(D-xx)` 代码标记（票 245 既有裁定：现况表 13 行／代码 0 枚，加一枚会打破 1:1 双向核对）。
- 未新增状态、未动 `HotkeyStatus` 枚举、未加 `internal/agent/approval → internal/ball` 依赖边（本腿没碰 approval 包）。

---

## §6 残差与归口（一格都不藏）

1. **AC#4**：本腿补了「配置 → 链 → 球拿到的那枚 cancel」这一跳（§3），**只能交裁决腿判**「够不够翻、要不要补真机端到端」。**票面那句「阻塞于票 228」的前提今天已倒**（§0 第三行：258 已把四槽接进常驻腿），请裁决腿按现况重判，别按 09-30 的阻塞句。
2. **AC#6① 的 winlive 读数**欠一发（§2.3）——本波禁真机；`hotkeys live 3/4` 那句要由下一发 winlive 窗口现跑。
3. **AC#6② 的 winlive 两处改动**（`requireIdleRoster`／`live_windows_test.go`）〔改后未实测〕：只过了 `go vet -tags winlive`（编得动）。**这不是「我已验」**。
4. **与票 260 的边界（点名，不改）**：本腿**没有**改「稳态绑不绑」的行为（`registerAllWith` 的 cancel 跳过一字未动），260 的 `TestConfiguredCancelStillNeverBoundWhileIdle260` 与 245 的尝试计数钉都还在、还绿。⇒ AC#6② 那句「配置成组合键时应当常绑」的**行为半边今天仍不成立**（260 已裁「稳态不绑」）；本腿只落地裁定的**判据精度**半边，行为半边要动就必须先有一枚人工批准（铁律：改裁定＝要 owner 一句话）。
5. **AC#7**（借被拒的记账）已由 `245-r2`（`02799be92`）落地、本腿复跑绿，未动。**AC#8／AC#9 属真机**（第二个真进程收 Esc／借还台件），本波**未碰**。
6. **归我复跑的既有物**：`cmd/wisp` 45 枚存量红（§4）其中 `TestTicket256ResidentGateOptionsFieldSetIsTheFiveItClaims`／`TestAC247HandingTheLevelToTheSeamIsNotVisibility` 两枚是**他票名册过期**（不是环境），本腿**不越界改**，只在此具名；归属于票 256 与票 247 的下一片。
