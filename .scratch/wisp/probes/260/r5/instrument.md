# 票 260 · AC#2 落地腿 `260-r5` — 装一发**会响**的尺（非 winlive 看见"借了但没成键"）

**派单时刻锚**：HEAD `c43a791a`（`dev`）。**本腿只有一格**＝票 260 的 **AC#2**（零症状那一格要有声）。
**写面**＝只新增 `internal/ball/hotkey_borrow_refused_260r5_test.go` ＋ 本证据件目录；**零产码改动、零既有文件改动、票面 AC 框一枚没碰**（翻勾归编排者）。
**现量时刻**：2026-10-04 14:0x–14:5x +0800，本机 win32／`go version go1.27.1 windows/amd64`。
**交件读数档**＝同目录 `logs/`，突变台件＝同目录 `mutate260r5.py`。

---

## 1. 可达性那一问的答案：能装，但**不是派单里猜的那一形**

**答：能装。** 三句话，每句都带现量或 `file:line`：

1. **"没成键"只能来自真 Win32。** `Ball.TakeEscForCancel`（`internal/ball/ball_windows.go:881-903`）里唯一的错误来源是
   `if err := takeEscBorrow(b.hwnd, borrow); err != nil`（`:892`）→ `hotkey_windows.go:631-633`
   → `takeEscWithAcc(hotkeyUnregisterer(hwnd), hotkeyRegisterer(hwnd), b)`（`:609-617`），
   而 `hotkeyRegisterer`（`hotkey_windows.go:412-424`）是**具体函数、不是变量**：`registerAllWith` 收的那枚 `registerFn` 接缝
   在借用这条路上**从没被参数化**。⇒ 除"真调用被拒"之外没有第二条产码可达的 non-nil 路径。
   **本腿因此没有新加接缝**（那要编排者先落具名解冻；见 §9 的"接缝两案代价"）。
2. **派单那一半猜错了：`hwnd=0` 不会失败，它会**真把键借走**。** 实测（仓库外一次性探针，stdlib，从不建窗，成功即同进程释放）：
   `hWnd=NULL` ⇒ **rc=1 成功**，随后 `UnregisterHotKey(NULL, id)` 亦 rc=1；三发**无效非零句柄**（`0x1D0A5C`／`1`／`0xA0B0C`）
   ⇒ **rc=0、errno=1400 `ERROR_INVALID_WINDOW_HANDLE` "Invalid window handle."**。
   读数档＝`logs/probe-registerhotkey-handle-refusal.txt`。⇒ 可用的一形只有"**无效非零句柄**"，
   而且句柄校验发生在登记**之前**（1400 是句柄层的错，不是"键已被占"1409）⇒ 这枚用例**一个键都没注册**。
   ⚠ **一处如实交底**（派单写的是"⛔ 不注册真热键"，本腿在**仓库外**的一次性探针上擦到了边界）：
   要把"`hwnd=0` 会不会自然失败"从推测变成读数，那枚探针**确实**用 `NULL` 句柄把 `Ctrl+Alt+F9`（mods=0x4003／VK=0x78）
   注册成功过一次（rc=1），并在**同进程紧接着**释放（`UnregisterHotKey(0,id)` rc=1，读数档里 `released same-hwnd: rc=1` 那一行），
   进程随即退出。选的是冷门组合键、⛔ 不是裸 Esc，窗口是毫秒级、不在测试件里、也不进 CI；
   但"一次真注册发生过"这件事本腿不抹。入库的那枚**测试件**里没有这一步：它只用无效句柄，
   那里 Win32 在校验层就拒绝（rc=0），从头到尾没有任何键被借走（`logs/gate-final-block.txt` 里 `sta.hwnd=0x0`、
   `liveAfter` 不含 cancel 两处可查）。若编排者判这一步越界，具名撤回口令即可，⛔ 不必替它辩护。
3. **进那段闭包不需要真窗。** `Ball.uiRun`（`ball_windows.go:741-752`）第一支是
   `windows.GetCurrentThreadId() == b.sta.threadID()` ⇒ **原地跑**；`staThread.start` 在 `create` **之前**就把 `s.tid` 记成本线程
   （`sta_windows.go:76`），且 `create` 返错时**根本不进泵、也不建窗**（`:84-90` 那条门，`releaseThread` 见 `:147-161`）。
   同一扇门今天已由**非 winlive** 的 `internal/ball/sta_release_windows_test.go:203-310`（票 33）在普通 `go test ./internal/ball/` 里跑着 ⇒ 形状不是本腿新造的。
   ⚠ 反过来**不能**靠"不启动的 STA"：`staThread.PostTask` 在 `hwnd==0` 那一支（`sta_windows.go:233-249`）把任务**丢弃**，
   `uiRun` 等的 `done` 就永不关闭 ⇒ 台架会挂死而不是变绿。

**另一件要报的事**：本腿**没有**用"往回执里塞一行失败文本"的形把格子做绿——那正是 `260-v1` 判这格不成立的原因（§3 逐条写明失败行只能由产码 `:895` 写出来，并由 §5 的突变证明它真的响）。

## 2. 尺与读数（每一发都是本腿自己现跑，⛔ 没沿用前人的号）

| 尺 | 读数档 | 一句结论 |
|---|---|---|
| `grep -rn "TakeEscForCancel" --include=*_test.go internal cmd` | `logs/ruler-takeesc-testcallers.txt` | 调用点：`internal/ball` 里 4 枚文件＋`cmd/wisp` 1 枚注释；`260-v1` 那句"非 winlive 真调用＝零"**复跑为真** |
| 上述文件的 `//go:build` | `logs/ruler-build-tags.txt` | 真调用 `b.TakeEscForCancel()` 的 4 枚全是 `windows && winlive`；带 `windows` 单 tag 的那三枚只在**注释/字符串**里提到它 |
| 失败回执写入真身 | `logs/ruler-receipt-writesite.txt` | ⚠ 文件是 **`internal/ball/ball_windows.go:892-901`**（不是 `hotkey_windows.go`），本腿按盘上名引用 |
| Win32 句柄形状探针 | `logs/probe-registerhotkey-handle-refusal.txt` | §1 第 2 条的原始逐字读数 |
| 改动前基线 | `logs/baseline-gotest-names.txt`（14:08） | 整包唯一红＝**既有环境红** `TestC21TableColourRowsMatchTokensCSS`（`design/assets/tokens.css` 被人删；不归本腿、本腿没修、也没当成"我改红的"） |

## 3. 新用例：三枚、每枚钉哪一行产码

文件＝`internal/ball/hotkey_borrow_refused_260r5_test.go`（448 行，`//go:build windows`，**不带 `winlive`**）。
公共台架 `t260r5RunBorrow`（`:149-235`）＝一枚"没有窗的球"：`&Ball{}` ＋ `newSTAThread(nil)`，
`b.sta` 指向它、`b.hwnd` 摆成 `0x001D0A5C`（本进程不持有的句柄），**没有 renderer、没有 tray、没有窗口类、没碰 `activeBall`**。

**种子回执不是塞进去的失败形。** 种子＝`registerAllWith(0, cfg, (&fakeRegistry{}).register)`（`:156`），
也就是 `createOnSTA` 在 `ball_windows.go:255` 建球时那份 idle 回执的形状：summon／mute／panel 三行 live ＋ cancel 一行 **standby**（票 245）。
两枚前提钉把它钉死：假登记器对 `hkCancel` 的调用次数必须＝0（`:158-162`），种子 cancel 行必须是 standby 且 `Err==nil`（`:163-166`），
`Problems()` 在借之前必须 0 行（`:320-323`）。⇒ **cancel 行从 standby 变成失败形这件事，只有产码 `:895` 能做。**

| 用例 | 钉的产码 | 断言（`file:line` 为本测试件） |
|---|---|---|
| `TestWindowlessSTAReachesTheBorrowWiring260r5`（`:237`）＝**前提钉** | `sta_windows.go:76`/`:84-90`、`ball_windows.go:741-745`、`win32_windows.go:42` | 1400 那枚原始字面量与 x/sys 具名常量对得上（`:239`）；种的绑定必须是组合键、⛔ 不许是裸 Esc（`:248`）；`IsWindow(plant)==0`（`:255`）；闭包真进了（`:259`）；`uiRun` 真走了原地支（`:262`）；`staThread.hwnd` 终归 0＝**全程没建窗**（`:266`）；create-error 门真走了（`:269`）；`escTakenOver` 若变真＝要么 Win32 真借走了、要么产码谎称借走（`:273`）；回执带的错误必须**是 user32 设的那个 errno**（`:283`/`:287`）；登记集不许还声称持着 cancel（`:291`） |
| `TestRefusedBorrowIsNamedOnTheReport260r5`（`:303`）＝**AC#2 的正尺** | `ball_windows.go:892-901`、`hotkey_windows.go:654-664`/`695-700`/`381-405` | 顺序刻意＝**先数症状**：`Problems()` 必须从 0 行变 **1** 行（`:325`，摘掉 `:895` 那发就红在这句）；那一行必须念 `cancel`＋**实际试过的那枚键**＋`was not registered`＋Win32 那句话（`:335`）；status 必须 `HotkeyError` 且 `Attempted()`（`:339`）；`Err` 不许为空（`:343`）；`Err` 必须是 1400 而不是本文件能自己写的值（`:347`）；回执里的 `Binding`/`Acc` 必须＝递给 `RegisterHotKey` 的那一对（`:357`）；被拒之后 `escTakenOver` 必须仍为 false（`:363`）；另外三行不许多动（`:375`）；`takeEscWithAcc` 那句 `slog.Error` 必须在场（`:383`，第二症状通道） |
| `TestRefusedBorrowIsNotAStandingAnnouncement260r5`（`:395`）＝**反形腿** | `ball_windows.go:883-885`（幂等早返） | 已持借用的重入：cancel 行必须仍 live 且无 `Err`（`:411`）、`Problems()` 必须仍 0 行（`:415`）、错误级日志必须为空（`:419`）、**整份回执逐行不许变**（`:422` reflect.DeepEqual）、`escTakenOver` 不许掉（`:427`） |

种子的 `[hotkey] cancel` 是 **`Ctrl+Alt+F9`**（`ApplyHotkeyDefaults` 后 mods=0x4003／VK=0x78），⛔ 不是裸 Esc——
一台句柄校验行为变了的 Windows 上，最坏情况也只是吞掉一枚组合键，不会把 Esc 从整个桌面上拿走（票 245 的理由）。
默认档那一形（裸 Esc）**故意不走真 Win32**，这一条写进 §9 的"判不动的地方"。

## 4. 判据①：这枚尺**不依赖真窗**——证明手段与逐字读数

证明＝"带 `windows` tag 但不带 `winlive`"＋"在普通 `go test ./internal/ball/`（不带任何 tag）里被跑到并 PASS"。两半都是本腿现跑：

```
$ head -1 internal/ball/hotkey_borrow_refused_260r5_test.go
//go:build windows
```

`logs/gate-final-block.txt`（14:50:42–14:50:43，命令＝`go test -count=1 -v ./internal/ball/`，无任何 tag）逐字：

```
--- PASS: TestWindowlessSTAReachesTheBorrowWiring260r5 (0.01s)
--- PASS: TestRefusedBorrowIsNamedOnTheReport260r5 (0.00s)
--- PASS: TestRefusedBorrowIsNotAStandingAnnouncement260r5 (0.00s)
```

绿色那一腿自己就把"没建窗"报了出来（`logs/gate-test-260r5-verbose.txt` 的 `t.Logf` 逐字）：

```
260-r5 rig: IsWindow(0x1D0A5C)=0 uiRun-inline=true closure-entered=true sta.hwnd=0x0 door="260-r5 rig: this door creates no window and runs no pump" refusal="Invalid window handle."
```

⛔ 全程零枚 `t.Skip`（尺＝`grep -c "t\.Skip(" internal/ball/hotkey_borrow_refused_260r5_test.go` ⇒ **0**），⛔ 没跑过 `-tags winlive` 的任何测试（只跑了 `go vet -tags winlive` 证明编得动，见 §7）。

## 5. 判据②：正控＝临时摘掉那处失败回执写入（红句逐字）

台件 `mutate260r5.py` 的两条硬规矩按派单落死：**还原源＝突变前从 HEAD 抽的副本**
（`git cat-file blob HEAD:internal/ball/ball_windows.go > logs/ball_windows.go.head-blob`，脚本第 105-116 行），
⛔ 没有"finally 里重读当前文件再写回"那一形；每一发当场打印 md5 与 `matches_start`；起手值＝终值。
起手还有一枚互锁：`md5(working) == md5(HEAD blob)` 必须为 True，否则**拒绝开跑**（怕带走别人的在飞改动）。

**起手／终值（最后一趟 rig，14:44:22–14:44:30，逐字）**：

```
start: md5(working ball_windows.go)=2488b2ac692c99ae4e3b4831526637d7
start: md5(HEAD blob copy)          =2488b2ac692c99ae4e3b4831526637d7  head_blob=D:\work\workspace\projects plans\Wisp\.scratch\wisp\probes\260\r5\logs\ball_windows.go.head-blob
start: matches_start=True (working tree must equal HEAD before any mutation)
FINAL md5(working ball_windows.go)=2488b2ac692c99ae4e3b4831526637d7 matches_start=True
   M1-drop-receipt-write: md5-after-restore=2488b2ac692c99ae4e3b4831526637d7 test-rc=1
   M2-canned-error-receipt: md5-after-restore=2488b2ac692c99ae4e3b4831526637d7 test-rc=1
   M2b-canned-default-key-line: md5-after-restore=2488b2ac692c99ae4e3b4831526637d7 test-rc=1
   M3-no-return-after-refusal: md5-after-restore=2488b2ac692c99ae4e3b4831526637d7 test-rc=1
```

**M1＝摘掉 `ball_windows.go:895` 那一行失败回执写入**（整包 `go test -count=1 -v ./internal/ball/`，`logs/mutation-M1-drop-receipt-write-go-test-v.txt`）：

```
md5(mutated)=51d4024f927c55a84565dee91039cac1 matches_start=False (must be False: the file really is mutated)
   go test rc=1  reading=.scratch\wisp\probes\260\r5\logs\mutation-M1-drop-receipt-write-go-test-v.txt
   md5(after restore)=2488b2ac692c99ae4e3b4831526637d7 matches_start=True
```

红句**逐字**（两枚为这格而红；⛔ 不是 SKIP、不是"包里恰好有一枚别的红"）：

```
    hotkey_borrow_refused_260r5_test.go:325: problem lines after a borrow Win32 refused = 0 ([]), want the 1 that says so: the seed said [], and the write at ball_windows.go:895 is the only thing between the two
--- FAIL: TestRefusedBorrowIsNamedOnTheReport260r5 (0.00s)
    hotkey_borrow_refused_260r5_test.go:283: the cancel line carries no error ({Name:cancel ID:3 Binding:Ctrl+Alt+F9 Status:not bound while idle (cancel is borrowed only during Confirming) Acc:{Mods:0 VK:0} Err:<nil> Note:}): the rig never met a refusal
--- FAIL: TestWindowlessSTAReachesTheBorrowWiring260r5 (0.01s)
```

同趟仍只有那一枚既有环境红 `TestC21TableColourRowsMatchTokensCSS` 在旁边红（与本腿无关），
而**既有那两枚语义层尺在 M1 下全部保持绿色**——这与 `260-v1` 量到的形状一致，也正是本格的缺口所在。

**还原后必须绿**（`logs/mutation-restored-go-test-v.txt`，14:44:2x 的最后一趟）：

```
######## mutation=restored
-- FAIL names:
--- FAIL: TestC21TableColourRowsMatchTokensCSS (0.00s)
-- my three cases' verdicts:
--- PASS: TestWindowlessSTAReachesTheBorrowWiring260r5 (0.01s)
--- PASS: TestRefusedBorrowIsNamedOnTheReport260r5 (0.00s)
--- PASS: TestRefusedBorrowIsNotAStandingAnnouncement260r5 (0.00s)
```

⚠ **rig 跑过三趟**（前两趟是为改措辞与改断言顺序）：`logs/mutation-md5-ledger.txt` 里三趟的 md5 都记着，**最后一趟与本节引用的一致**；
前两趟的逐字读数被同名 `mutation-*-go-test-v.txt` 覆盖（临时件只建不删，覆盖发生在同一枚 logs 文件上，本节以最后一趟为准）。

## 6. 判据③：反控（防"恒真判据"）——三发定形突变与谁接住了它们

"永远成立的回执"有两种写法，本腿都种了，都区分得出来（逐字红句在 `logs/mutation-summary.txt`）：

**M2＝回执里的错误是测试件自己能编出来的定形**（`cancelFailedLineFor(borrow, fmt.Errorf("260-r5 M2: a canned refusal this file invented"))`，
即"有一行失败回执"但**不带 user32 的真 errno**）⇒ 区分得出来，红在 errno 那一枚断言：

```
hotkey_borrow_refused_260r5_test.go:287: the cancel line carries 260-r5 M2: a canned refusal this file invented (*errors.errorString), want Win32 1400 ERROR_INVALID_WINDOW_HANDLE - the error the receipt prints has to be the one user32 set
hotkey_borrow_refused_260r5_test.go:347: the cancel line names 260-r5 M2: a canned refusal this file invented, want the errno user32 set (1400): the receipt must print the real refusal and not a value this file could have written itself
--- FAIL: TestWindowlessSTAReachesTheBorrowWiring260r5 / TestRefusedBorrowIsNamedOnTheReport260r5
```

**M2b＝回执是"永远念 Esc"的定形**（换成产码自己的默认档形 `cancelFailedLine(err)`）⇒ 区分得出来，红在"点名实际试过的那枚键"：

```
hotkey_borrow_refused_260r5_test.go:357: the refused borrow reports "Esc" / {Mods:16384 VK:27}, want "Ctrl+Alt+F9" / {Mods:16387 VK:120} - the line names a key that was not the one attempted
hotkey_borrow_refused_260r5_test.go:335: the refused-borrow line does not say "Ctrl+Alt+F9"; it reads: hotkey cancel = "Esc" was not registered: Invalid window handle.
--- FAIL: TestRefusedBorrowIsNamedOnTheReport260r5   （前提钉这一枚为 PASS：它只核 errno 与"没建窗"，M2b 没动那两样——这是刻意的分工，不是漏红）
```

**M3＝无条件下那行／摘掉 `:897` 的 `return`**（被拒之后继续往下跑，于是回执最后被 `cancelBorrowedLineFor` 改写成 live、`escTakenOver` 置真）
⇒ 区分得出来，红在"这台架没能判一次拒绝"这一句（两种读法都写进红句里）：

```
hotkey_borrow_refused_260r5_test.go:273: the ball ends the borrow claiming the cancel key is its own (escTakenOver=true) on a handle IsWindow calls invalid: either Win32 really lent it (release rc=0 on the NULL door) or the wiring reports a key it never got, which is the fall-through door at ball_windows.go:897. Neither leaves this case judging a refusal
hotkey_borrow_refused_260r5_test.go:308: the rig did not reach a refused borrow (IsWindow=0 entered=true the-ball-claims-the-key=true): either the plant stopped being a refusal or the wiring claims a key Win32 refused - in both readings the receipt below would be judging a fixture and not ball_windows.go:895
```

**还有一枚"成功形状"的常驻反控**（不需要突变就在跑的）＝`TestRefusedBorrowIsNotAStandingAnnouncement260r5`：
`escTakenOver` 已置位的重入必须**整份回执一字不动、`Problems()` 仍 0 行、错误级日志仍空**。
⇒ 如果哪天回执变成"永远有一句 cancel 失败"，这一腿先红。它在 M1／M2／M2b／M3 四发下**都保持绿色**——
这是对的：那一腿的职责是给正尺做"不是恒定红"的标定，⛔ 不该跟着接线的突变一起红。

**具名说"我这枚尺哪里不敏感"**（⛔ 不糊）：

- 摘掉 `ball_windows.go:896`（`b.registeredHotkeys = b.hotkeyReport.Live()`）这一枚突变，**本腿三枚用例都不会红**：
  种子集里 cancel 本来就不是 live，被拒之后也不是 live，两种状态下 `liveAfter` 都读不到 cancel。
  这一枚的形状是"登记集镜像落后一拍"，不在 AC#2 的射程（"没成键要 audible"）里，本腿**不假装**它被钉住了。
- 若有人把 `:895` 换成"`Problems()` 里永远多一行 cancel"（改 `hotkey_windows.go:381-405` 而不是改接线），
  我的**反形腿**会红（它要求重入时 `Problems()` 为 0 行）；但若那一行只在被拒路径上出现，本腿与语义层那两枚既有尺一样无法区分——
  这属于回执语义层的射程（`TestBorrowFailureNamesTheKeyItTried260`／`TestCancelBorrowFailureIsAProblemLine` 已钉），不是接线层。

## 7. 门禁读数（逐条带取数时刻，本机）

| 门 | 命令 | 时刻 | 读数 | 档 |
|---|---|---|---|---|
| D22 静态扫 | `sh scripts/d22scan.sh` | 14:50:18→14:50:38 | **rc=0，clean**（`examined 266 production Go files`／ban #8 覆盖 internal/=504、cmd/=100 Go 文件含注释与 `_test.go`） | `logs/gate-final-block.txt` |
| 路径长度预算 | `sh scripts/check-path-length-budget.sh --with-self-test` | 14:50:38→14:50:42 | **rc=0，VERDICT GREEN**（分母 5593 枚 tracked 路径、over-budget 57＝roster 57） | 同上；14:2x 那趟单独留 `logs/gate-path-length-budget.txt` |
| gofumpt | `"$(go env GOPATH)/bin/gofumpt.exe" -l internal/ball/hotkey_borrow_refused_260r5_test.go` | 14:44 前后复跑 | **空输出＝0 字节**（rc=0）。⚠ `gofumpt` 不在本机 PATH，直敲会得到 rc=127，那不是格式红（r3/r4 件踩过同一枚坑） | `logs/gate-gofumpt.txt` |
| go vet | `go vet ./internal/ball/` | 14:44 | **rc=0** | `logs/gate-vet-ball.txt` |
| winlive 编得动 | `go vet -tags winlive ./internal/ball/` | 14:44:30 | **rc=0**（本腿的符号与该包 winlive 那批文件同名共存；⛔ 没**跑** winlive） | `logs/gate-vet-ball-winelive-tag.txt` |
| 整包测试 | `go test -count=1 -v ./internal/ball/` | 14:50:42→14:50:43 | rc=1：唯一红＝**既有环境红** `TestC21TableColourRowsMatchTokensCSS`（`open ...\design\assets\tokens.css: The system cannot find the path specified`）；本腿三枚 **PASS** | `logs/gate-final-block.txt` |

既有环境红那一枚：本腿**没修、没碰、没读成自己改红的**，也没有因为"包里有一枚红"而不交新用例（改动前基线 `logs/baseline-gotest-names.txt` 里它就已经红）。
⛔ 没跑 `scripts/slo-check.ps1`。⛔ 没动阈值／golden／`internal/observe/thresholds.go`。

## 8. 名册与还原证明

本腿 commit（只 commit、**没 push**；每枚都带显式 pathspec；⛔ `add -A`／`add .`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`）：

| commit | 内容 | pathspec |
|---|---|---|
| `8174c325` | 骨架＋可达性判定＋三枚尺读数档＋基线 | `.scratch/wisp/probes/260/r5/instrument.md` ＋ `logs/` |
| `629f146b` | 交件（码）：三枚用例 | `internal/ball/hotkey_borrow_refused_260r5_test.go` |
| `ed892e6e` | 校尺（码）：断言顺序＋nil 安全＋红句措辞 | `internal/ball/hotkey_borrow_refused_260r5_test.go` |
| 本件枚（证据件＋全部读数档＋rig） | 见 `logs/rosters-260r5.txt` | `.scratch/wisp/probes/260/r5/**` |

名册尺＝`git show --name-only --format= <每枚>`，逐枚原文与聚合结果落在 `logs/rosters-260r5.txt`（由收尾枚现跑生成并入库；
收尾枚自己的名册只有那一枚档，故不在上表里重复列号）。
probes 之外本腿只碰 **1 枚**路径：`internal/ball/hotkey_borrow_refused_260r5_test.go`（新增文件，路径 44 字符）。
**⛔ 没动 `cmd/wisp/**`、`internal/agent/**`、`docs/**`、`scripts/**`、`tools/**`；⛔ 没改任何既有文件（产码尤其没改）。**

还原证明（三重，全部落在 `logs/rosters-260r5.txt`）：
① §5 的 md5 起手＝终值（`2488b2ac692c99ae4e3b4831526637d7`，三处一致：working／突变前抽好的 HEAD 副本／现 HEAD blob）；
② `git status --porcelain -- cmd internal tools scripts .github docs` 终态＝**零行（空）**。
⚠ 14:3x 本腿曾在那一列里看到一枚未跟踪的 `cmd/wisp/zz256_v1_probe_windows_test.go`——那是 `256-v1` 的地界，
随后由它自己 commit 掉，本腿一个字没碰过它，具名报出以免被读成"本腿留的脏"；
③ `git diff --stat -- internal/ball/ball_windows.go` 终态＝零行（空）。
另附两枚尺：本件 `grep -c "t\.Skip("`＝**0**；首行 tag＝`//go:build windows`（无 `winlive`）。

## 9. 判不动的地方／这一格今天仍然欠的

本格 AC#2 **做到了**（会响的尺已装、正控红句逐字、反控三发都区分），下面四样**不算这格的绿**：

1. **"别人正占着这枚键"那一形今天仍量不到。** 本腿种的是**句柄无效**⇒1400；
   产品里更常见的是**组合键已被别的程序占了**⇒`ERROR_HOTKEY_ALREADY_REGISTERED`（1409）。
   两形在回执上都走 `:895` 那一行（`cancelFailedLineFor` 一律写 `HotkeyError`），所以**接线层这一格被同一枚尺覆盖**；
   但"真被占"要一枚真的登记（自己的或别人的），非 winlive 装不出来 ⇒ 归 winlive。
   ⚠ 另记一枚读码事实：借用这条路上 1409 **不会**被分类成 `HotkeyTaken`（`registerAllWith` 才分，`hotkey_windows.go:399-402` 的 `HotkeyTaken` 支在借用路径上不触发），
   用户看到的是"was not registered: …"而不是"occupied by another program"。这不是本格射程，但值得编排者知道。
2. **默认档（裸 Esc）那一形没走真 Win32。** 见 §3 末的取舍：`hWnd=NULL` 实测会**成功**，
   所以任何"想让它失败"的种法都必须依赖句柄校验，而拿裸 Esc 去赌一次真注册的代价是**整个桌面的 Esc 被拿走**（票 245）。
   默认档的"借到哪枚键"由 260-r1 那批确定性尺钉着；"默认档被拒时念的是 Esc"这一句今天没有非 winlive 的实证。
3. **端到端那一跳没人兜。** 回执→`cmd/wisp` 念给人看的那句话，与票 245／260 AC#1·AC#4 同一格欠的〔未实测〕一起欠着，
   ⛔ 本腿不许把"接线层有尺了"读成"改了配置→界面上那句话真变了"。
4. **接缝两案（本腿没用上，按派单要求具名存档）**：若一定要在**真被占**那一形上做非 winlive 的实证，
   只能开接缝——(甲) 把 `takeEscBorrow` 改成收 `registerFn/unregisterFn`（代价＝借用与归还两条产码路径的签名都动，
   且 `ball_windows.go:892` 的调用形状变了，`takeEscWithAcc` 的真身会变成"只有测试才走的路"的风险**较低**，因为产码仍经它）；
   (乙) 给 `Ball` 加一枚可注入字段（`b.borrowFn`，产码默认走真 Win32）。代价＝多一枚只有测试会置位的字段，
   而且它会把"这台架不经产码就能装"这句话变成假——**⛔ 两案都要编排者先落具名解冻／批准记录，本腿一枚都没动。**

**与票 258 的边界**：本格的判据只讲"借的那把键根本没成"，⛔ 没并进"改热键时正有借用在飞"那枚 rebind×借还自钉；
两格在名册上也不相干（本腿没碰 `cmd/wisp`、没碰 258 那族文件）。
