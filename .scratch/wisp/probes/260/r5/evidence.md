# 票 260 · AC#2「零症状那一格要有声」— 腿 `260-r5` 复跑件（2026-10-05 补派那一程）

**这一程是什么**：`A616` §2「写 2＝`260-r5`」在 2026-10-05 11:10 补派的那一枚。派单要求＝补一枚**非 winlive 也会执行**、射程＝**接线本身**的尺，判据＝**摘掉 `internal/ball/ball_windows.go:892-899` 那处失败回执的写入，这枚尺必须红**，且「这一发突变要真跑、红句逐字抄进证据件、还原后复跑证明字节相同」。

**★先报一条编排者需要知道的（本程最大的发现，不是代码）**：这一格的**尺昨天已经装好并已被收账**——`260-r5` 前一段（2026-10-04）落了六枚 commit（`8174c325` → `5d680d78`，名册见 §6），交付物＝`internal/ball/hotkey_borrow_refused_260r5_test.go`（`//go:build windows`，**不带 `winlive`**，三枚用例），台账 `A603` §5 那一节标题逐字「`260-r5` 收账（票 260 AC#2 那格的新仪器）」，证据件＝同目录 `instrument.md`（254 行／24,515 字节／九节）。⇒ `A616` §2 把同一格又派了一次（它 §2 末行还写着「在飞＝…／260-r5／…」）。**本程因此没有再写一枚同射程的尺**（那会造出两枚互相冒充的仪器，且按 `A603` §5「实现方判语必须归另一枚非实现者腿」，这一格的终裁仍欠 `260-v2`）。本程做的那件事＝**在今天这个锚点上把判据真跑一遍**（票面第 7 行自写「⚠ 引用前先重跑」）＋**补两发昨天没种的突变形状**（`M1b`＝摘掉块内**两枚**写入、`M4`＝只摘掉登记集镜像那一枚）＋把读数落在派单点名的这一枚 `evidence.md`（昨天的件叫 `instrument.md`，⛔ 本程没改它一字，见 §7 第 2 条）。

**现量时刻**：2026-10-05 11:12–11:20 +0800，本机 win32／`go version go1.27.1 windows/amd64`。
**写面**＝`internal/ball/**`（本程**零字节改动**，见 §4 末那把 md5 尺）＋ `.scratch/wisp/probes/260/r5/**`。⛔ 没碰 `cmd/wisp`、`internal/tools`、`internal/panel`、票面 AC 框、台账。

---

## 0. 起手锚（同发取数，逐字）

第 0 步那一发（三条命令同批）：

```
Mon Oct  5 11:12:53 CST 2026
c6cf66e6
---STATUS---
 M cmd/wisp/firstrun.go
```

- 日期换算＝2026-10-05 11:12:53 +0800（CST）。
- `git log -1 --format=%h`＝`c6cf66e6`（`dev`，即 `A616` §0 自己记的那一枚）。
- `git status --porcelain -- cmd internal`＝**1 行**：` M cmd/wisp/firstrun.go` ⇒ 那是**别的腿**（`257-r2`）的地界，本程一个字没碰过它。
- 收笔时同发复尺（11:19:35 后）＝**2 行**：多出 `?? cmd/wisp/firstrun_257_test.go`，仍是 `257-r2` 在写它自己的票 257 台件。**本程在两行之间对 `cmd/**` 零写入**（唯一碰过磁盘的产码路径＝`internal/ball/ball_windows.go` 的三次突变窗口，见 §3 末那条归因说明）。

起手先跑包、再动笔（撞钉预检定式）：那把尺的读数在 §1。

## 1. 起手名册与读到的既有断言原文

### 1.1 起手整包（改动之前，`go test -count=1 ./internal/ball/`，11:13:35 → 11:13:38，rc=1）

档＝`logs-recheck/baseline-gotest-plain.txt`（原始 stdout 全量，含 slog 那些 INFO/ERROR 行）。非 `-v` 那一发只点名红的那一枚，逐字：

```
--- FAIL: TestC21TableColourRowsMatchTokensCSS (0.00s)
    tokens_table_test.go:1468: read design/assets/tokens.css: open D:\work\workspace\projects plans\Wisp\design\assets\tokens.css: The system cannot find the path specified. - the CSS leg of this check must never skip
FAIL
FAIL	github.com/CarlosShao/wisp/internal/ball	0.312s
```

⇒ **起手红名册＝1 枚**，就是 `AGENTS.md`／派单具名的「不是本腿的红」（`design/assets/tokens.css` 被别的会话删了，早于票 260）。派单令＝具名登记、⛔ 不修不追：**本程没修它、没读它、没因为它不交尺**。

### 1.2 逐名册（同一份字节状态，11:16:56 → 11:16:58，`go test -count=1 -v ./internal/ball/`，rc=1）

`-v` 才给得出逐名，所以逐名册取**未突变**那一发；它与 §1.1 起手状态**字节相同**——凭据＝那一发前后 `internal/ball/ball_windows.go` 的 md5 都是 `2488b2ac692c99ae4e3b4831526637d7`（§3 的起手＝终值链），且本程**从未新增或修改任何 `internal/ball` 下的 `.go`**。档＝`logs-recheck/mutation-restored-go-test-v.txt`。

- 计数：**PASS 68／FAIL 1／SKIP 0**（尺＝`grep -c '^--- (PASS|FAIL|SKIP):'`）。⛔ 零枚 SKIP ⇒ 没有任何「读成通过的跳过」。
- 唯一红＝`TestC21TableColourRowsMatchTokensCSS`（同 §1.1）。
- 本票那一族的逐名（全 PASS）：
  - **本尺三枚**：`TestWindowlessSTAReachesTheBorrowWiring260r5`（#12）／`TestRefusedBorrowIsNamedOnTheReport260r5`（#13）／`TestRefusedBorrowIsNotAStandingAnnouncement260r5`（#14）
  - 260-r1/r2 那批（回执与借用取值语义层）：`TestBorrowDefaultIsBitIdentical260`（#15）／`TestBorrowFollowsConfiguredBinding260`（#16）／`TestBorrowUnparsableFallsBackLoudly260`（#17）／`TestBorrowReceiptSharesOneSourceWithRegistration260`（#18）／`TestBorrowFailureNamesTheKeyItTried260`（#19）／`TestConfiguredCancelStillNeverBoundWhileIdle260`（#20）／`TestLiveBorrowHelperPremiseIsConfigDependent260`（#21）／`TestBorrowWishFollowsSeed260r2`（#22）／`TestBorrowWishDefaultCellIsBitIdentical260r2`（#23）／`TestBorrowWishRefusedSeedIsNotSilent260r2`（#24）
  - 票 245/64 那两枚既有语义层尺：`TestCancelBorrowRoundTrip`（#27）／`TestCancelBorrowFailureIsAProblemLine`（#28）
  - 69/74 那族 C21 表尺（同包，与本票无关，只作名册）：`TestC21TableColourRowsMatchCode`／`TestC21GeometryRowsMatchCodeConstants`／`TestC21TokenConsumerReport`／`TestC21CodeTokensAreTabledOrExempt`／`TestTokenGoldenValues`／`TestNoHardcodedColorsInBallPackage`
  - 其余 41 枚（dock／liquid／position／sta／hotkey／tokens／visual／hit）逐名都在档里，本程不 Relevant 也不动。
- 定向复跑（11:15:30 → 11:15:31，`-v -run '260r5'`）＝三枚 **PASS**、`ok 0.059s`，`=== RUN` 只有三行（⇒ 那把 pattern 没顺带匹配到别人的用例）。

### 1.3 派单那把尺在今天的读数（「非 winlive 真调用 `TakeEscForCancel`」到底几枚）

`grep -rn "TakeEscForCancel" --include=*_test.go internal cmd` 命中 14 行（逐字档＋每枚文件的 `head -1` build tag＋`.TakeEscForCancel()` 真调用点名册＝`logs-recheck/ruler-takeesc-callers-20261005.txt`，11:2x 现跑）。**逐行分形**（tag 一枚一枚现读）：

| 文件 | 首行 tag | 那几行是**真调用**还是提及 |
|---|---|---|
| `internal/ball/hotkey_borrow_refused_260r5_test.go` | `//go:build windows` | **真调用 1 枚＝`:202`**；其余 5 行（`:14/:44/:145/:233/:391`）是注释 |
| `internal/ball/hotkey_cancel_borrow_260_test.go` | `//go:build windows` | `:48` **注释**（`cleanReport260` 的说明句） |
| `internal/ball/hotkey_status_test.go` | `//go:build windows` | `:257` **注释** |
| `internal/ball/hotkey_cancel_borrow_live_260_test.go` | `//go:build windows && winlive` | 真调用 `:37`／`:70` |
| `internal/ball/hotkey_live_test.go` | `//go:build windows && winlive` | 真调用 `:373` |
| `internal/ball/interaction_live_test.go` | `//go:build windows && winlive` | 真调用 `:141` |
| `internal/ball/live_windows_test.go` | `//go:build windows && winlive` | 真调用 `:126` |
| `cmd/wisp/resident_cancel_key_wording_260r3_windows_test.go` | `//go:build windows` | `:21` **注释** |

⇒ **今天的读数**：`winlive` 那四枚文件里 5 枚真调用（CI 与普通 `go test` 都不执行它们，`260-v1` 那句因此仍然为真）；**非 winlive 的真调用＝1 枚＝`hotkey_borrow_refused_260r5_test.go:202`**。`260-v1` 判 AC#2 不成立那天，这一枚＝**0**。

### 1.4 读到的既有断言原文（本程据以划射程的那两处）

`internal/ball/hotkey_status_test.go:268-283`，`TestCancelBorrowFailureIsAProblemLine` 里那两行——**测试自己把失败行塞进回执**，产码 `:892-899` 一次都没走：

```go
	rep := registerAllWith(0, fullConfig(), (&fakeRegistry{}).register).
		withCancel(cancelFailedLine(syscall.EINVAL))
```

`internal/ball/hotkey_cancel_borrow_260_test.go:47-51`，同一个形状（`line` 由调用方给，不由 `TakeEscForCancel` 生成）：

```go
// cleanReport260 is an idle report over cfg with its cancel line replaced, which
// is the shape TakeEscForCancel leaves behind on the real ball.
func cleanReport260(cfg HotkeyConfig, line HotkeyBinding) HotkeyReport {
	return registerAllWith(0, cfg, (&fakeRegistry{}).register).withCancel(line)
}
```

⇒ 这两枚的射程＝**回执/`Problems()` 语义层**（「这一行如果被写进去，`Problems()` 会不会念」）。它们在 §3 的 M1／M1b 下**今天仍然全部 PASS**（逐名读数在 §3.4）——这与 `260-v1` 量到的形状一致，也正是 AC#2 那格的缺口；本程⛔ 没动它们一字（派单令＝不许为变绿去动既有尺的射程，反之也不许顺手改严）。

## 2. 新尺的形状，以及为什么它的射程是「接线本身」

文件＝`internal/ball/hotkey_borrow_refused_260r5_test.go`（448 行，首行 `//go:build windows`，无 `winlive`）。三枚用例的分工（本程逐行读过，形状与 `instrument.md` §3 所述一致，未被后来任何腿改动——凭据＝本件 §3 的突变仍然打得响）：

1. **`TestWindowlessSTAReachesTheBorrowWiring260r5`＝前提钉**。它不判症状，它判**这台架到底有没有走到那一段**：`IsWindow(0x1D0A5C)==0`（`:255`）、闭包真进了（`:259`）、`uiRun` 真走了原地支（`:262`）、`staThread.hwnd` 终归 0＝全程没建窗（`:266`）、create-error 门真走了（`:269`）、回执带的错误**必须是 user32 自己设的 errno 1400**而不是本文件能编出来的值（`:287`）。⇒ 如果哪天接线搬家或 Windows 行为变了，这一枚先红，另外两枚的红才有意义。
2. **`TestRefusedBorrowIsNamedOnTheReport260r5`＝AC#2 的正尺**。顺序刻意＝**先数症状再数形状**：`Problems()` 借之前 0 行（`:320`）→ 借之后必须 1 行（`:325`）；那一行必须念 `cancel` ＋ **实际递给 `RegisterHotKey` 的那枚键** ＋ `was not registered` ＋ Win32 那句本地化的话（`:333-336`）；`line.Status` 必须 `HotkeyError` 且 `Attempted()`（`:338`）；被拒后 `escTakenOver` 必须仍 false（`:362`）；另外三槽不许多动（`:371-378`）；`takeEscWithAcc` 那句错误级日志必须在场（`:382`）。
3. **`TestRefusedBorrowIsNotAStandingAnnouncement260r5`＝反形腿**（幂等早返那一支）：已持借用的重入必须**整份回执逐行不动**（`:421` `reflect.DeepEqual`）、`Problems()` 仍 0 行（`:414`）、错误级日志仍空（`:418`）。⇒ 它管的是「这枚尺不是一把永远红的尺」。

**为什么射程是接线而不是语义层自证**——三条，每条都可在盘上核：

- **失败行不是测试塞进去的**。种子回执＝`registerAllWith(0, cfg, (&fakeRegistry{}).register)`（`:157`），即 `createOnSTA` 建球时那份 idle 形（三行 live ＋ cancel 一行 **standby**）。两枚前提钉把「种子不许自带失败」钉死：假登记器对 `hkCancel` 的调用次数必须 0（`:158-161`）、种子 cancel 行必须 `HotkeyStandby` 且 `Err==nil`（`:162-165`）、`Problems()` 借之前必须 0 行（`:319-322`）。⇒ **cancel 行从 standby 变成失败形，盘上只有 `ball_windows.go:895` 一处能做到**，而 §3 的 M1 正是把它摘掉。
- **失败是 Win32 自己说的，不是桩说的**。`takeEscBorrow`→`takeEscWithAcc`→真 `RegisterHotKey`；借用这条路**没有** `registerFn` 接缝（`hotkeyRegisterer` 是具体函数不是变量），所以除真调用被拒之外没有第二条 non-nil 路径。种的是**无效非零句柄** ⇒ 句柄校验发生在登记**之前**（1400 而非 1409）⇒ 这枚用例一个键都没注册（派单令「⛔ 不许真注册全局热键」满足，见 §5 第 4 条对昨天那次仓库外探针的复述与归口）。
- **走的是产码入口**。`:202` 调的是 shipped 的 `b.TakeEscForCancel()`，不是 `takeEscWith` 那类原语，也不是 `withCancel(...)`。台架只负责把 `b.hwnd` 摆成本进程不持有的句柄、把 STA 的 create 那一支引到报错门（`sta_windows.go:84-90` 那条形状今天已由票 33 的**非 winlive** `sta_release_windows_test.go` 在同一条 `go test ./internal/ball/` 里跑着，不是本票新造的口子）。

## 3. 突变名册（本程真跑）＋红句逐字＋还原凭据

台件＝`mutate260r5-recheck.py`（M1／M1b／M3）＋`mutate260r5-m4.py`（M4）。两枚都**新写**、都只把读数落到 `logs-recheck/`，⛔ 不覆写昨天 `logs/` 里被 `instrument.md` 逐字引用的那批档（原因见 §7 第 2 条）。

**还原源＝突变之前从 HEAD 抽的那一份**（`git cat-file blob HEAD:internal/ball/ball_windows.go` → `logs-recheck/ball_windows.go.head-blob`），⛔ 没有「finally 里重读当前文件再写回」那一形；起手互锁＝`md5(working)==md5(HEAD blob)` 必须为真，否则**拒绝开跑**（怕带走别人的在飞改动）。

### 3.1 突变名册与逐字红句

**M1＝摘掉 `:895` 那一枚失败回执写入**（票面判据原形；档＝`logs-recheck/mutation-M1-drop-receipt-write-go-test-v.txt`，窗口 11:16:46 → 11:16:49，rc=1）。红句逐字（按档内出现顺序，两枚为接线而红；整段即 `logs-recheck/mutation-M1-drop-receipt-write-go-test-v.txt` 的原样）：

```
    hotkey_borrow_refused_260r5_test.go:283: the cancel line carries no error ({Name:cancel ID:3 Binding:Ctrl+Alt+F9 Status:not bound while idle (cancel is borrowed only during Confirming) Acc:{Mods:0 VK:0} Err:<nil> Note:}): the rig never met a refusal
--- FAIL: TestWindowlessSTAReachesTheBorrowWiring260r5 (0.02s)
    hotkey_borrow_refused_260r5_test.go:325: problem lines after a borrow Win32 refused = 0 ([]), want the 1 that says so: the seed said [], and the write at ball_windows.go:895 is the only thing between the two
--- FAIL: TestRefusedBorrowIsNamedOnTheReport260r5 (0.00s)
```

这一发 M1 那一趟的**整包 FAIL 名册**＝`TestWindowlessSTAReachesTheBorrowWiring260r5`／`TestRefusedBorrowIsNamedOnTheReport260r5`／`TestC21TableColourRowsMatchTokensCSS`（第三枚＝既有环境红），SKIP 计数 0。⇒ **判据成立**：摘掉那处写入，本尺红，而且红句自己就在念「`Problems()` 里那一行没了」。

**M1b＝摘掉块内**两枚**写入**（`:895` 报告行＋`:896` 登记集镜像，留 `if err != nil` 与 `return`）——这一形是派单那句「`:892-899` 那处失败回执的写入」的**较强读法**，昨天没种过。档＝`logs-recheck/mutation-M1b-drop-both-writes-go-test-v.txt`，窗口 11:16:49 → 11:16:53，rc=1。红句逐字（两枚，与 M1 同两处断言）：

```
    hotkey_borrow_refused_260r5_test.go:283: the cancel line carries no error ({Name:cancel ID:3 Binding:Ctrl+Alt+F9 Status:not bound while idle (cancel is borrowed only during Confirming) Acc:{Mods:0 VK:0} Err:<nil> Note:}): the rig never met a refusal
--- FAIL: TestWindowlessSTAReachesTheBorrowWiring260r5 (0.02s)
    hotkey_borrow_refused_260r5_test.go:325: problem lines after a borrow Win32 refused = 0 ([]), want the 1 that says so: the seed said [], and the write at ball_windows.go:895 is the only thing between the two
--- FAIL: TestRefusedBorrowIsNamedOnTheReport260r5 (0.00s)
```

⚠ **不许把 M1b 读成「覆盖比 M1 多一格」**：M1b 的红**仍然全部来自 `:895`**；`:896` 那一枚自己单独摘掉时本尺看不见，那是量到的、不是推测的，见下面 M4。

**M3＝摘掉 `:897` 的 `return`**（被拒之后继续往下跑，于是回执末了被 `cancelBorrowedLineFor` 改写成 live）——昨天种过同一形，今天复跑。档＝`logs-recheck/mutation-M3-no-return-after-refusal-go-test-v.txt`，窗口 11:16:53 → 11:16:56，rc=1。红句逐字：

```
    hotkey_borrow_refused_260r5_test.go:273: the ball ends the borrow claiming the cancel key is its own (escTakenOver=true) on a handle IsWindow calls invalid: either Win32 really lent it (release rc=0 on the NULL door) or the wiring reports a key it never got, which is the fall-through door at ball_windows.go:897. Neither leaves this case judging a refusal
--- FAIL: TestWindowlessSTAReachesTheBorrowWiring260r5 (0.02s)
    hotkey_borrow_refused_260r5_test.go:308: the rig did not reach a refused borrow (IsWindow=0 entered=true the-ball-claims-the-key=true): either the plant stopped being a refusal or the wiring claims a key Win32 refused - in both readings the receipt below would be judging a fixture and not ball_windows.go:895
--- FAIL: TestRefusedBorrowIsNamedOnTheReport260r5 (0.00s)
```

**M4＝只摘掉 `:896` 那枚登记集镜像写入**（本程新增的**灵敏度边界探针**，答的是「这枚尺哪里看不见」）。档＝`logs-recheck/mutation-M4-drop-mirror-only-go-test-v.txt`，窗口 11:19:33 → 11:19:35，**rc=0**：

```
--- PASS: TestWindowlessSTAReachesTheBorrowWiring260r5 (0.01s)
--- PASS: TestRefusedBorrowIsNamedOnTheReport260r5 (0.00s)
--- PASS: TestRefusedBorrowIsNotAStandingAnnouncement260r5 (0.00s)
```

⇒ **本尺在这一枚上不敏感**，原因可解释：种子集里 cancel 本来就不 live，被拒之后也不该 live，两种状态下 `liveAfter` 都读不到 `hkCancel`。这一枚的形状是「登记集镜像落后一拍」，不在 AC#2 的射程（「没成键要 audible」）里，本程⛔ 不假装它被钉住了（→ §5 第 3 条）。⚠ **这一发的读法口径要说清**：M4 用的是定向 `-run 260r5`（只跑本尺三枚）而不是整包，因为它问的问题是「本尺红不红」；「整包在 M4 下也只红那一枚既有环境红」这句**本程没测**，不许从这一发推出来。

### 3.2 md5 与还原凭据（起手值＝终值，逐字取自 `logs-recheck/mutation-md5-ledger.txt`）

```
start: md5(working ball_windows.go)=2488b2ac692c99ae4e3b4831526637d7
start: md5(HEAD blob copy)          =2488b2ac692c99ae4e3b4831526637d7
start: matches_start=True (working tree must equal HEAD before any mutation)
   M1-drop-receipt-write: md5-after-restore=2488b2ac692c99ae4e3b4831526637d7 test-rc=1
   M1b-drop-both-writes: md5-after-restore=2488b2ac692c99ae4e3b4831526637d7 test-rc=1
   M3-no-return-after-refusal: md5-after-restore=2488b2ac692c99ae4e3b4831526637d7 test-rc=1
FINAL md5(working ball_windows.go)=2488b2ac692c99ae4e3b4831526637d7 matches_start=True
```

M4 那一趟（同档末段）：`start md5=2488b2ac692c99ae4e3b4831526637d7` → mutated `3cfc5ee9d780e796aa20abae9f4b2791` → `FINAL md5=2488b2ac692c99ae4e3b4831526637d7 matches_start=True`。

突变态各枚 md5（可重放核对）：M1＝`51d4024f927c55a84565dee91039cac1`（**与昨天 `instrument.md` §5 记的那一枚逐字相同** ⇒ 突变片段与文件都未漂）、M1b＝`94811a15a6641db2b82ac31e6696216b`、M3＝`41ccbf322f6792307aa8e2858a6bf104`、M4＝`3cfc5ee9d780e796aa20abae9f4b2791`。

三重还原证明（全部现跑）：
1. `md5(internal/ball/ball_windows.go)` 起手＝终值＝`2488b2ac692c99ae4e3b4831526637d7`，且与 `git cat-file blob HEAD:…` 现读同值（11:18:24 那一发也复核过）；
2. `git status --porcelain -- cmd internal tools` 终态＝**2 行，全在 `cmd/wisp`（`257-r2` 的地界）**，`internal/` **零行**（档＝`logs-recheck/gate-readings.txt` 末段）；
3. `git diff --stat -- internal/ball/ball_windows.go`＝**空**（同因：第 2 条里 `internal/` 那一列为空已含它）。

### 3.3 本程真造出的突变窗口（给别的腿做归因用，逐枚带时刻）

`cmd/wisp` 编译时把 `internal/ball` 当源码一起编。⇒ 我这三个窗口里如果 `257-r2` 或任何跑 `cmd/wisp` 测试的腿正好在跑，它的红名册可能被我的突变洗脏。**本程没发生过 `cmd/wisp` 的跑（我没跑），但别人那侧我看不见**，所以逐枚报出：

| 突变态 | 窗口（起 → 止） | 盘上处于突变态的时长 |
|---|---|---|
| M1 | 11:16:46 → 11:16:49 | ≈3 s |
| M1b | 11:16:49 → 11:16:53 | ≈4 s |
| M3 | 11:16:53 → 11:16:56 | ≈3 s |
| 未突变对照 | 11:16:56 → 11:16:58 | 0（未突变） |
| M4 | 11:19:33 → 11:19:35 | ≈2 s |

合计 **≈12 秒**，全部集中在 11:16:46–11:19:35 这一分钟段内。

### 3.4 既有那两枚语义层尺在 M1 下的逐名（本程现跑，`logs-recheck/mutation-M1-drop-receipt-write-go-test-v.txt`）

```
--- PASS: TestBorrowFailureNamesTheKeyItTried260 (0.00s)
--- PASS: TestCancelBorrowFailureIsAProblemLine (0.00s)
--- PASS: TestRefusedBorrowIsNotAStandingAnnouncement260r5 (0.00s)
```

（未突变那一发同样三枚 PASS，见 `logs-recheck/mutation-restored-go-test-v.txt`。`TestRefusedBorrowIsNotAStandingAnnouncement260r5` 在 M1／M1b／M3 三发下都保持绿——它的职责是给正尺做「不是恒定红」的标定，⛔ 不该跟着接线的突变一起红，这一条与 `instrument.md` §6 的分工一致。）

## 4. 门禁读数（四门，逐条带时刻；档＝`logs-recheck/gate-readings.txt`）

| 门 | 命令 | 时刻 | 读数 |
|---|---|---|---|
| D22 静态扫 | `sh scripts/d22scan.sh` | 11:18:16 前起（末行 `(end 11:18:16)`） | **rc=0，clean**（`examined 266 production Go files`；ban #8 覆盖＝`internal/` 508 枚、`cmd/` 102 枚 Go 文件含注释与 `_test.go`，`design/` 39、`frontend/` 85 文本件） |
| 路径长度预算 | `sh scripts/check-path-length-budget.sh --with-self-test` | 11:18:16 → 11:18:23 | **rc=0，VERDICT GREEN**（分母 5733 枚 tracked 路径；over-budget 57＝roster 57、not in roster 0；longest 180 相对＝`252-…`；runner 最差全路径 224）；三发正控 1/3–3/3 全 ok |
| go vet | `go vet ./internal/ball/` | 11:18:23 → 11:18:24 | **rc=0，零输出** |
| gofumpt | `"$(go env GOPATH)/bin/gofumpt.exe" -l internal/ball/hotkey_borrow_refused_260r5_test.go`，`go env GOPATH`＝`D:\work\base\gopath` | 11:18:24 → 11:18:24 | **空输出，rc=0**（逐枚点名那一枚尺；⛔ 没用 `~/GOPATH/...` 那把瞎尺） |
| 整包 | `go test -count=1 -v ./internal/ball/`（未突变那一发） | 11:16:56 → 11:16:58 | rc=1，唯一红＝`TestC21TableColourRowsMatchTokensCSS`（既有环境红），本尺三枚 PASS，SKIP 0 |

- ⛔ 本程没跑 `winlive` 的任何测试（连 `go vet -tags winlive` 都没跑——那是昨天 `instrument.md` §7 跑过的，本程没新增需要它编得动的符号）。⛔ 没跑 `scripts/slo-check.ps1`。⛔ 没动阈值／golden／`internal/observe/thresholds.go`／`tools/d22scan/allowlist.txt`。⛔ 没放宽任何断言、没加 `t.Skip`、没把任何 SKIP 读成通过。
- **emoji 那一支门的射程要说清（免得下一枚腿把这行读成"证据件过了 emoji 门"）**：仪器扫的是 `internal/`＋`cmd/` 的 Go 文件（含注释与 `_test.go`，上面那行 508／102 计数就是它自己报的）＋`design/`／`frontend/` 文本件；`.scratch/**` 下的 markdown 与 `.py` 台件**不在射程内**。本件因此沿用仓内文档的既有符号（`⛔`／`★`／`⇒`，与 `instrument.md`、票面、`AGENTS.md` 同形）。⛔ 本程没写任何 Go 源码 ⇒ Go 那一侧零新增符号（写面名册见 §6），emoji 门那一句对**真产生 Go 的形态**是「不适用」，不是「过了」。
- 路径长度那一门里那两行 `warning: in the working copy of '…'，LF will be replaced by CRLF` 是**该脚本正控自己在 `/tmp` bench 目录里种的台件**打的，不是本仓的脏面；bench 目录按规则 8 留在盘上不删（`/tmp/tmp.HL2CGJQD87`）。

## 5. 判不动的格／量不到的格（逐条具名归口，⛔ 不自己划掉）

1. **「别的程序正占着这枚键」那一形（`ERROR_HOTKEY_ALREADY_REGISTERED` 1409）今天仍量不到**。非 winlive 要种它就得真登记一次。接线层这一格由同一枚尺覆盖（`cancelFailedLineFor` 对任何非 nil err 都写 `HotkeyError`），但「真被占」归 `winlive`。⇒ **归编排者＋机主那一个词**（`A616` §4 待人项里那一条）。
2. **默认档（裸 Esc）被拒时念的是 `Esc` 这一句没有非 winlive 实证**。种子的句柄校验发生在登记之前，所以任何「想让它失败」的种法都依赖那一层；拿裸 Esc 去赌一次真注册的代价是整个桌面的 Esc 被拿走（票 245 的理由）。⇒ 本尺刻意用组合键 `Ctrl+Alt+F9`，这一格**没做**，⛔ 不许从 AC#2 的绿推出来。
3. **`ball_windows.go:896` 单独摘掉＝本尺不敏感**（今天实测 rc=0，见 §3.1 M4）。形状＝「登记集镜像落后一拍」，不在 AC#2 射程。要钉它得另一枚尺（射程＝`RegisteredHotkeys()` 与回执的一致性），**本票没顺手做**；要不要单开一格归编排者判。
4. **接缝两案**（把 `takeEscBorrow` 改成收 `registerFn/unregisterFn`，或给 `Ball` 加可注入借用函数）——两案都要**编排者先落具名解冻／批准记录**。本程⛔ 一枚产码字节都没动，`hotkeyRegisterer` 仍是具体函数。昨天那次**仓库外一次性探针真注册过 `Ctrl+Alt+F9`（rc=1，随后同进程释放）**已由 `260-r5` 主动交底、`A603` §4 具名上报；本程**没有重复那一发**（判据不需要它，且派单令＝⛔ 不许真注册全局热键），昨天那格的处理仍挂在机主裁决上。
5. **端到端那一跳没人兜**：回执 → `cmd/wisp` 念给人看的那句话。与票 245／260 AC#1·AC#4 同一格欠的〔未实测〕一起欠着。本程⛔ 没把「接线层有尺了」写成「改了配置→界面上那句话真变了」。
6. **`cmd/wisp` 侧今天仍然零枚非 winlive 真调用**（§1.3 那张表最后一行＝注释提及）。若 AC#2 的射程被判定要包括宿主侧，那一格**本程动不了**：派单第 33 行与 `A616` §2 都写明 `cmd/wisp` 此刻是 `257-r2` 的地界 ⇒ **停手上报，不自行跨包**。
7. **AC#2 这一格本腿不许自己翻勾**。凭据纪律＝`A603` §5 那句「实现方判语必须归另一枚非实现者腿」，`260-v2` 此刻按在 `A616` §2 的队列里（理由＝跨包突变互斥）。票面第 20 行那个框**一字未动**（尺＝`git status --porcelain -- .scratch/wisp/issues`＝空，见 §8）。

## 6. commit 名册

### 6.1 前一程那六枚（2026-10-04，本程逐枚 `git show --name-only` 现核）


| # | commit | 内容 | 名册（probes 之外的路径） |
|---|---|---|---|
| 1 | `8174c325` | 骨架＋可达性判定＋尺读数档 | 只 `probes/260/r5/**` |
| 2 | `629f146b` | 交件（码）：三枚用例 | `internal/ball/hotkey_borrow_refused_260r5_test.go` |
| 3 | `ed892e6e` | 校尺（码）：断言顺序＋措辞 | `internal/ball/hotkey_borrow_refused_260r5_test.go` |
| 4 | `6ab2efdd` | 交件（证据）：`instrument.md` 九节 | 只 `probes/260/r5/**` |
| 5 | `5ad1a7b8` | 收尾：名册与还原证明 | 只 `probes/260/r5/**` |
| 6 | `5d680d78` | 交底（证据）：仓库外探针那次真注册 | 只 `probes/260/r5/**` |

**本程那一枚（第 7 枚）**：写它这一节的时候那一笔还没落地，所以**不预填号**；名册（枚枚点名，全部显式 pathspec）＝

- `.scratch/wisp/probes/260/r5/evidence.md`（本件）
- `.scratch/wisp/probes/260/r5/mutate260r5-recheck.py`（M1／M1b／M3 台件）
- `.scratch/wisp/probes/260/r5/mutate260r5-m4.py`（M4 灵敏度边界台件）
- `.scratch/wisp/probes/260/r5/logs-recheck/` 十枚档：`ball_windows.go.head-blob`／`baseline-gotest-plain.txt`／`gate-readings.txt`／`mutation-md5-ledger.txt`／`mutation-M1-drop-receipt-write-go-test-v.txt`／`mutation-M1b-drop-both-writes-go-test-v.txt`／`mutation-M3-no-return-after-refusal-go-test-v.txt`／`mutation-M4-drop-mirror-only-go-test-v.txt`／`mutation-restored-go-test-v.txt`

⛔ `add -A`／`add .`；⛔ `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；⛔ 没 push（推送由编排者做，本程只往本地加笔）。核号尺（任何腿可自取）＝`git log --oneline -- .scratch/wisp/probes/260/r5`。

### 6.2 补录（第 8 枚那一笔现跑的逐字名册）

上一节那条「不预填号」的承诺在这一节兑现：本件第 7 枚的号与逐枚 `git show --name-only` 名册由第 8 枚追加，⛔ 两处读数都来自现跑不是回忆。

## 7. 自我对抗（本程真抓出来的东西，逐条）

1. **差点拿昨天的读数交今天的件**。起手我确实准备直接引 `logs/mutation-M1…txt`；抓到它的不是洁癖而是派单那句「这一发突变要真跑」＋票面第 7 行「引用前先重跑」。复跑之后多出来两样昨天没有的东西：M1b 这一形（派单原话是 `:892-899`，昨天只摘 `:895`）与 M4 这一枚**边界**（昨天在 §6 里写「不会红」是靠读码判断，本程把它跑成了 rc=0 的实读数）。
2. **差一点毁掉前一程的证据**。`mutate260r5.py` 的 `run_tests()` 写死 `logs/mutation-<tag>-go-test-v.txt`，而 `instrument.md` §5/§6 逐字引用着那几个文件——直接重跑=**覆写别人已入库的凭据**（规则 8「临时件只建不删」的精神＋台账 `A603` 那条「引用它的读数一旦失效，那张裁决表就掉回仅自述」）。处置＝新写两枚台件、日志落 `logs-recheck/`；尺＝`git status --porcelain -- .scratch/wisp/probes/260/r5/logs` 现跑＝**空**（见 §3 与本节，`logs/` 一字未动）。
3. **本件初稿把 §3.4 那段插到了 §8 里面**（章节序号错位：一份九节件里出现「§8 之下挂着 3.4」）。抓它的尺＝落笔后把 § 标题按行序拉了一遍，⛔ 不是靠通读。处置＝把那段整体移回 §3 之下（§3.3 与 §4 之间），本件现在序号单调。**这一条留在账上**，因为它正是「骨架名册」那把尺存在的理由：先有编号清单，错位的段落才会被看见。
4. **M1b 不是多出来的牙**——它的红仍来自 `:895`。如果我把它写成「覆盖了 `:896`」，那就是给一票多报一格灵敏度；M4 正是为了让这句话有一个反方向的实读数。
5. **M4 用的是定向 `-run`，不是整包**。「本尺不敏感」这个结论站得住；「整包在 M4 下只红一枚环境红」这个结论**站不住**，本程没测，已在 §3.1 写明不许推。
6. **起手那发 `go test` 不带 `-v`，只点名 1 枚红**。逐名册因此来自 11:16:56 那一发；我用 md5 链（起手＝终值）证明这两发读的是**同一份字节**，⛔ 不是「我改过之后再测」。这一条若不写，名册的取数时刻就成了装饰。
7. **窗口污染的可能**：`cmd/wisp` 把 `internal/ball` 当源码编。我在 §3.3 把五个窗口的起止时刻逐枚报出，而不是写「几秒钟，无所谓」——`A603` §5 记过一条同类互斥（两枚带突变的腿并发＝互相洗读数，跨包比同包更隐蔽）。本程是唯一在飞的那一侧看不见的人，所以给料不给结论。
8. **`cmd internal tools` 那一列两行脏**（`M cmd/wisp/firstrun.go`／`?? cmd/wisp/firstrun_257_test.go`）⛔ 被我写成「干净」。前一程在 `instrument.md` §8 里踩过同族形状（那一枚是 `zz256_v1_probe_windows_test.go`），本程同样具名报成「别人的地界」（§0、§3.2）。
9. **有没有可能这一格其实还欠码？** 我把三种「还欠」的可能各自证了一次：(甲) 射程若要求覆盖 `:896` ⇒ 那是 §5 第 3 条那枚**新格**，不是 AC#2 那句话；(乙) 若要求 1409 那一形 ⇒ 非 winlive 装不出来（§5 第 1 条），欠的是机主的词不是我的码；(丙) 若要求宿主侧 ⇒ 跨包，派单令停手（§5 第 6 条）。⇒ 本程不新写第二枚同射程尺；写一枚的话，两把尺会互冒充，而 AC#2 只需要一把会响的。
10. **⛔ 用 SKIP 或 tag 把麻烦洗掉**：本尺 `t.Skip` 计数 0（`instrument.md` §4 那把尺本程复核＝首行 `//go:build windows`、无 `winlive`；本程整包读数里 SKIP＝0）。我⛔ 没给这枚尺加 tag 让它躲开 M4 那种不好看的答案。

## 8. 交件判语

- **票 260 AC#2 的判据，在今天这个锚点上成立**：摘掉 `internal/ball/ball_windows.go:895` 那处失败回执写入（M1）或摘掉块内两枚写入（M1b）⇒ `internal/ball` 有两枚**为接线层而红**的用例（红句逐字在 §3.1），而既有那两枚语义层尺 `TestCancelBorrowFailureIsAProblemLine`／`TestBorrowFailureNamesTheKeyItTried260` **保持 PASS**（§3.4）；还原后本尺三枚复跑绿、`ball_windows.go` md5 起手＝终值、`internal/` 零脏面。

- **本程零产码改动**（写面 6 枚路径全在 `.scratch/wisp/probes/260/r5/**` 之下；`internal/ball/**` 字节未变，凭据＝§3.2 三重还原）。四门全 rc=0（§4）。
- **本程没做的事，具名**：⛔ 没翻票面任何框（尺＝`git status --porcelain -- .scratch/wisp/issues`＝空，§0 那一发 `cmd internal` 两行脏全在 `cmd/wisp`）；⛔ 没改台账；⛔ 没修 `TestC21TableColourRowsMatchTokensCSS`（派单令＝具名登记、不修不追，§1.1）；⛔ 没动 `cmd/wisp`／`internal/tools`／`internal/panel`／禁区六枚；⛔ 没跑 winlive；⛔ 没 push。
- **交件给编排者的三句判断**（都不许由本腿自裁）：
  1. **这一格派重了**：`A603` §5 已在 2026-10-04 收账，`A616` §2 于 2026-10-05 又写「写 2＝`260-r5`」并把 260-r5 同时列为「在飞」。若那是**故意重派**（例如要第二枚腿复测），请具名说明它要补什么；若不是，则本件 §1–§4 就是那一次复测的凭据，⛔ 不必再有第三程。
  2. **AC#2 的终裁仍欠 `260-v2`**（非实现者腿），本件是实现方读数，⛔ 不许被读成裁决。
  3. **两处「不是本票的红但会误导下一枚腿」**：`TestC21TableColourRowsMatchTokensCSS`（缺 `design/assets/tokens.css`）与 §5 第 3 条那枚登记集镜像的灵敏度边界。后者若判要钉，请单开一格并点名写面。
