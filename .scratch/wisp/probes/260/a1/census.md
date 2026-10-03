# 票 260 AC#0 只读普查证据件（腿 `260-a1`）

判死问题：借用那一遍（`Confirming` 借那把取消键）**读不读 `[hotkey] cancel` 的配置值**。

---

## §0 锚点与口径

- 锚点 sha：`47b3876583a59a46a14bd52cdbffca1e2da8f5be`（分支 `dev`，`git log -1 --format=%H` 现跑）。
- 取数时刻：`2026-10-03 09:48 +0800`（起手 `date` 现跑）。
- 口径：**本程零 Go 命令**——没有跑过 `go build`／`go vet`／`go test`／`go run`（同机有三枚写腿在飞，任何 Go 编译都会互洗它们的读数）。全部读数只用 Read／Grep／Glob 与 `grep`／`sed`／`awk` 级文本工具现取。
- 本程自证：本文件的每一枚 commit 的 `git show --name-only` 只含 `.scratch/wisp/probes/260/a1/census.md` 这一枚路径（逐枚读数见 §6 交件名册）。
- 禁令执行：`frontend/**` 与 `design/**` 零读零写零转述（本件不引用、不转述其任何内容）；三枚冻结件与 `docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt` 一律未动；工单票面未动。
- 路径口径：下面所有引用一律报到「包/目录级路径:行」，只给文件名的引用视为不可核。
- 与台账 `A560` 的关系：`docs/reports/pending-and-issues.md:11080` 那三句是编排者自跑，本件**每一句都现跑重取**，不复述其字面；本件推翻它的地方集中在 §5。

---

## §1 问题①：借用那遍到底读不读配置值

一句话结论：**不读**（借用体注册的加速器是仓内自造的常量对，配置值在这条支路上从未被解引用）。

链条逐处（配置值 → 结构 → Options → 借用体）在下面的 §1.1–§1.4 展开，§1.5 给据以判死的那一行。

### 1.1 `[hotkey] cancel` 从 `config.toml` 被读进哪个结构

| 环节 | 落点（包/目录级路径:行） | 盘上字面 |
|---|---|---|
| TOML 键 → Go 字段 | `internal/config/schema.go:182` | `Cancel string \`toml:"cancel" default:"Esc"\`` |
| 所属结构 | `internal/config/schema.go:179` | `type HotkeySection struct {`（`:176-178` 注释自述 `[hotkey]`; hot-tier; `cancel` 在 `Confirming` 期间临时接管、会话结束归还＝D36） |
| 挂进根配置 | `internal/config/schema.go:112` | `Hotkey  HotkeySection  \`toml:"hotkey"\`` |
| 重载时的合并表 | `internal/config/manager.go:281` | `{"hotkey", &cur.Hotkey, &fresh.Hotkey, func() { cur.Hotkey = fresh.Hotkey }},` |
| 生效档位 | `internal/config/tiers.go:38` | `"hotkey":  "hot",` |

⚠ 结构口径：`config.HotkeySection`（包 `internal/config`）与 `ball.HotkeyConfig`（包 `internal/ball`，`internal/ball/hotkey_windows.go:49-54`）是**两枚不同的类型**，字段同名不共享；两者之间只做映射的地方全仓 **3 处**，逐处见 §1.2。`ApplyHotkeyDefaults`（`internal/ball/hotkey_windows.go:87-102`）是映射之后调的那一枚补空函数，`:95-97` 逐字 `if cfg.Cancel == "" { cfg.Cancel = d.Cancel }`——它**只填空、不改写非空值**，所以它不是"配置被吞掉"的现场，别把它读成吞。

### 1.2 谁把它交给 `ball.Options.Hotkeys`

全仓（tracked `*.go`，含测试）把 `config.HotkeySection` 映射成 `ball.HotkeyConfig` 的地方只有 3 处，**没有一处落在出厂主机的那一遍开机注册上**：

| # | 落点 | 性质 | 读数 |
|---|---|---|---|
| A | `cmd/balldebug/main.go:237-242` | **旁支调试程序**（`-config` 才走），不是出厂主机 | `bridge := ball.NewHotkeyReloader(b, b.ConfiguredHotkeys(), func() ball.HotkeyConfig{ h := mgr.Config().Hotkey; return ball.ApplyHotkeyDefaults(ball.HotkeyConfig{Summon: h.Summon, Mute: h.Mute, Cancel: h.Cancel, Panel: h.Panel}) })`（`:238`＝`h := mgr.Config().Hotkey`） |
| B | `internal/ball/hotkey_live_test.go:132`／`:160`（映射体 `hotkeysFromConfig` 在 `:499`） | 测试 | `bootCfg := ApplyHotkeyDefaults(hotkeysFromConfig(mgr))` |
| C | `internal/ball/hotkey_reload.go:16-19` | **注释里的示例，不是代码** | 该文件 `:15-20` 整段是 doc 注释（`:15` `//	m, _ := config.NewManager(path, res)`）；`internal/ball/hotkey_reload.go:49` 逐字「of [hotkey], already defaulted via ApplyHotkeyDefaults)」是 `HotkeySource` 的注释 |

出厂主机那一遍：

- `cmd/wisp/resident_ball_windows.go:165` `b, err := ball.New(ball.Options{` → `:171` 逐字 **`Hotkeys:  ball.DefaultHotkeys(),`** ⇒ 交给球的是**编译期默认值**（`internal/ball/hotkey_windows.go:69` 逐字 `return HotkeyConfig{Summon: "Ctrl+Alt+Q", Mute: "Ctrl+Alt+M", Cancel: "Esc", Panel: "Ctrl+Alt+P"}`），`config.HotkeySection` 在这条腿上从未被读。
- 这条事实有仓内的名册尺自证（不是我的推断）：`cmd/wisp/config_readers_255.go:110-113` 的 `"hotkey"` 行逐字写着唯一读者是 `cmd/balldebug/main.go:238 [h := mgr.Config().Hotkey]`，出厂球「binds ball.DefaultHotkeys() instead」（`:113`）。
- 库内的第二次兜底：`internal/ball/ball_windows.go:153-154` 逐字 `if opts.Hotkeys == (HotkeyConfig{}) {` → `opts.Hotkeys = DefaultHotkeys()`（全零值也换默认；`Options` 字段自述在 `internal/ball/ball_windows.go:67` `Hotkeys  HotkeyConfig   // zero value = DefaultHotkeys()`）。

⇒ **出厂路径上配置断了两截**：第一截在这里（`Options.Hotkeys` 压根没吃配置），第二截在借用体里（§1.3）。两截都得算进票 260 的射程；AC#1 若只修第二截，出厂主机仍然读不到（`cmd/wisp/resident_ball_windows.go:171` 不改＝配不进来）。

### 1.3 借用那一遍用的是哪一枚值（链条逐处）

| 环节 | 落点 | 盘上字面 | 携带键信息？ |
|---|---|---|---|
| 出厂调用者 | `cmd/wisp/resident_approval_windows.go:441` | `b.TakeEscForCancel()` | 无参 |
| 球侧借 | `internal/ball/ball_windows.go:873`／`:878` | `func (b *Ball) TakeEscForCancel() {` → `if err := takeEsc(b.hwnd); err != nil {` | **只有 hwnd** |
| 借的原语壳 | `internal/ball/hotkey_windows.go:543-545` | `func takeEsc(hwnd windows.HWND) error { return takeEscWith(hotkeyUnregisterer(hwnd), hotkeyRegisterer(hwnd)) }` | 只有 seam |
| 借的原语体 | `internal/ball/hotkey_windows.go:531`／`:533` | `func takeEscWith(unreg unregisterFn, reg registerFn) error {` → `if err := reg(hkCancel, escBorrowAcc()); err != nil {` | 只有 seam＋id |
| 加速器来源 | `internal/ball/hotkey_windows.go:524` | `func escBorrowAcc() Accelerator { return Accelerator{Mods: modNoRepeat, VK: vkEscape} }` | **零参、返回字面量** |
| 两枚常量 | `internal/ball/hotkey_windows.go:116`／`:118` | `modNoRepeat = 0x4000` / `vkEscape = 0x1B` | — |
| 真正打到 Win32 | `internal/ball/hotkey_windows.go:391-392` | `r, _, err := pRegisterHotKey.Call(uintptr(hwnd), uintptr(id), uintptr(acc.Mods), uintptr(acc.VK))` | acc 来自上一行 |

⇒ 借用到手的值是 **`Accelerator{Mods: 0x4000, VK: 0x1B}`**＝裸 Esc＋`MOD_NOREPEAT`。这条链上**没有一环带键**：`TakeEscForCancel()`→`takeEsc(hwnd)`→`takeEscWith(unreg, reg)`→`escBorrowAcc()` 四枚签名里键的数量是零。

报告的写法也一起写死：`internal/ball/hotkey_windows.go:519` `const escBorrowBinding = "Esc"` 被 `:535`（失败日志的 `"binding", escBorrowBinding`）与 `:582`（`cancelBorrowedLine()` 的 `Binding:` 字段）引用，`:584` 再取 `Acc: escBorrowAcc()` ⇒ 借期内报表上 cancel 那一行**连字面都换成裸 Esc**；配置串在此期间只从 `ConfiguredHotkeys()`（`internal/ball/ball_windows.go:840-844`，读 `b.boundCfg`）还看得见。

### 1.4 idle 那一遍（对照组：它读了配置，但只分类、不注册）

- 出厂注册：`internal/ball/ball_windows.go:255` `b.hotkeyReport = registerAll(b.hwnd, b.opts.Hotkeys)` → `internal/ball/hotkey_windows.go:416-418` → `registerAllWith`（`:453`）。
- 跳过那一枚：`internal/ball/hotkey_windows.go:457-463`，其中 `:457` `if p.id == hkCancel {`、`:458` 注释逐字「The one slot an idle ball must not register. takeEscWith is the」、`:462` `rep.bindings = append(rep.bindings, cancelIdleLine(bindings[i]))`、`:463` `continue`。
- `cancelIdleLine`（`internal/ball/hotkey_windows.go:428-450`）**确实吃配置串**：`:435` `if _, err := ParseAccelerator(bind); err != nil {`——但它把解出来的 `Accelerator` **丢弃**（`_`），只用 `err` 分类，且这个分支里 `reg` 从未被调用。测试把这枚"不注册"钉住了：`internal/ball/hotkey_status_test.go:205-207` 逐字 `if reg.attemptsOf(hkCancel) != 0 { t.Fatalf("idle pass attempted cancel: ...") }`。
- 还的那一遍也吃配置：`internal/ball/ball_windows.go:904` `b.hotkeyReport = b.hotkeyReport.withCancel(cancelIdleLine(b.cancelBinding))`，`b.cancelBinding` 的来源是 `internal/ball/ball_windows.go:253`（开机 `= b.opts.Hotkeys.Cancel`）与 `:821`（rebind `= cfg.Cancel`）——**仅用于报表分类**（standby／disabled／unparsable 三选一），不产生任何 `RegisterHotKey`。

⇒ 设计对照的结论：**包里有两处会"读"配置来的 cancel 串（idle 分类、归还分类），零处会把它"注册"出去。** 缺的就是 §1.3 那一环。

### 1.5 判死那一句

**不读。**

据以判死的那一行：`internal/ball/hotkey_windows.go:524` —— `func escBorrowAcc() Accelerator { return Accelerator{Mods: modNoRepeat, VK: vkEscape} }`。

判法（不是"我没找到"，是"它不可能读"）：这枚函数是借用体取加速器的**唯一**来源（`internal/ball/hotkey_windows.go:533` 与 `:584` 两处引用，见 §2.3），**形参为空**、返回体是两枚包级常量的字面组合，因此 `HotkeyConfig.Cancel` 与它之间不存在任何数据通路；把链条往上追到球外面（§1.3 那四枚签名）也没有任何一环能塞进一枚键。改形的最小形状必然是**动签名**（见 §2.2），不是"把常量换掉"那么轻。


---

## §2 问题②：若要吃配置，需要哪些零件

一句话结论：零件在包里**已经全齐**（`Accelerator`／`ParseAccelerator`／`boundCfg`／`cancelBinding` 四枚都在），缺的只是把借用体那两个自造常量的来源换成配置值；下面逐枚给 `file:line`。

### 2.1 `Accelerator`（修饰键＋虚拟键）从哪来

| 零件 | 落点 | 字面／作用 |
|---|---|---|
| 类型 | `internal/ball/hotkey_windows.go:105-108` | `type Accelerator struct { Mods uint32; VK uint32 }` |
| 解析器（串→加速器） | `internal/ball/hotkey_windows.go:125-158` | `func ParseAccelerator(s string) (Accelerator, error)`；`:131-151` 逐段吃 `ctrl/alt/shift/win/esc/space`＋`parseVK`；`:156` `acc.Mods |= modNoRepeat`（**每条解析出来的组合都自带 NOREPEAT**，与借用体现在手写的 `Mods: modNoRepeat` 同形） |
| 修饰键常量 | `internal/ball/hotkey_windows.go:112-116` | `modAlt 0x0001` / `modControl 0x0002` / `modShift 0x0004` / `modWin 0x0008` / `modNoRepeat 0x4000` |
| 虚拟键常量 | `internal/ball/hotkey_windows.go:118-119`＋`parseVK :161-226` | `vkEscape = 0x1B`、`vkSpace = 0x20`，其余走 `parseVK`（F1-F12／方向键／单字符） |
| 注入口（seam） | `internal/ball/hotkey_windows.go:386`／`:389-401` | `type registerFn func(id uint32, acc Accelerator) error`；`hotkeyRegisterer(hwnd)` 就是它到 Win32 的那一手 |
| 已经解出来、又丢掉的那一枚 | `internal/ball/hotkey_windows.go:435` | `if _, err := ParseAccelerator(bind); err != nil {` —— **`_` 里躺着的就是"配置来的 cancel 加速器"**；同一行注释 `:425-427` 逐字承认「It parses the binding even though the parse result is thrown away」 |

⇒ "吃配置"不需要新造类型、不需要新解析器、不需要新 Win32 手：**解出来的加速器已经在 `cancelIdleLine` 的手上，只是被扔了。**

### 2.2 `mods` 怎么带进借用体

现状：§1.3 那条链上的四枚签名（`TakeEscForCancel()`／`takeEsc(hwnd)`／`takeEscWith(unreg, reg)`／`escBorrowAcc()`）**零枚携带键**，所以 mods 今天进不去。要进去，零件按"最少改动"排：

1. **键的来源字段（仓里已有，借用体没去读）**：
   - `internal/ball/ball_windows.go:124` `boundCfg HotkeyConfig // the set the last pass was given (rebind updates it)`，写在 `:254`（开机 `= b.opts.Hotkeys`）与 `:819`（rebind `= cfg`），读口 `ConfiguredHotkeys()` `:840-844`；
   - `internal/ball/ball_windows.go:126` `cancelBinding string // configured cancel binding, for B1 release`，写在 `:253` 与 `:821`。
   ⇒ 最短的一条是：**`TakeEscForCancel` 在 UI 线程里读 `b.cancelBinding`／`b.boundCfg.Cancel`，导出签名不动**（外部调用者 `cmd/wisp/resident_approval_windows.go:441`、`cmd/balldebug/main.go:635` 一字不改）。
2. **`escBorrowAcc()`（`internal/ball/hotkey_windows.go:524`）参数化**：从"零参返回字面量"变成"吃一枚 `Accelerator` 或吃配置串"。它是借用体唯一的加速器来源（引用者见 §2.3）。
3. **`takeEscWith` 加一枚键参数（`internal/ball/hotkey_windows.go:531`）**：这一动会波及它的**五处测试调用**（`internal/ball/hotkey_status_test.go:210`／`:235`／`:259`／`:270`／`:340`，形状逐枚见 §2.3）＋`takeEsc(hwnd)`（`:543-545`）。这是签名变更面最大的一枚，也是唯一能让"注入口测试"继续成立的做法——现 seam 的形状就是 `registerFn`（`:386`），传 acc 不动 seam。
4. **报表行的拼法**（行为不变也得改，否则报表说谎）：`escBorrowBinding`（`:519`）与 `cancelBorrowedLine()`（`:578-586`，字段 `Binding: escBorrowBinding`＋`Acc: escBorrowAcc()`）；`cancelFailedLine(err)`（`:591-596`）由 `cancelBorrowedLine()` 派生，跟着走。
5. **不需要动的**：WM_HOTKEY 的路由按 **id** 不按键——`internal/ball/ball_windows.go:662-663` 逐字 `case hkCancel:` → `b.fire(b.opts.Events.OnCancelHotkey)`，出厂接法在 `cmd/wisp/resident_ball_windows.go:177` `OnCancelHotkey: func() { recordCancelHotkey(onCancelEsc) }`。换掉加速器**不改**取消的落地路径；归还那一手也只按 id 丢（见 §2.4 第 4 条）。

### 2.3 `escBorrowAcc`／`escBorrowBinding` 的引用者（逐枚，含测试）

**直接引用（产码，包 `internal/ball`）**
- `internal/ball/hotkey_windows.go:524` 定义 `escBorrowAcc`（注释 `:521-523` 自述「VK_ESCAPE with nothing but the no-repeat flag. Same spelling ParseAccelerator("Esc") produces」）。
- `internal/ball/hotkey_windows.go:533` 借用体：`if err := reg(hkCancel, escBorrowAcc()); err != nil {`。
- `internal/ball/hotkey_windows.go:584` 报表行：`Acc: escBorrowAcc(),`。
- `internal/ball/hotkey_windows.go:519` 定义 `const escBorrowBinding = "Esc"` → 被 `:535`（失败日志 `"binding", escBorrowBinding`）与 `:582`（`Binding: escBorrowBinding`）引用。
- 派生者：`cancelBorrowedLine()`（`:578-586`）→ `cancelFailedLine(err)`（`:591-596`）→ 两处调用在 `internal/ball/ball_windows.go:881`／`:886`。

**间接引用（按"裸 Esc 的字面"断言的测试）**
- `internal/ball/hotkey_live_test.go:107` `if got := rep.Live()[hkCancel]; got.VK != vkEscape {` —— **只判 VK、不判 Mods**；红句 `the borrowed cancel slot holds VK 0x%X, want VK_ESCAPE 0x1B`（`:108`）。
- `internal/ball/hotkey_status_test.go:223` `if got := borrowed.Live()[hkCancel]; got.VK != vkEscape || got.Mods != modNoRepeat {` —— VK＋Mods 双判，`t.Errorf` 字面「want the bare VK_ESCAPE with MOD_NOREPEAT」（`:224`）。
- `internal/ball/hotkey_status_test.go:85-91` helper `bindsEsc()`（`:87` `if acc.VK == vkEscape && acc.Mods == modNoRepeat {`），**四处调用**：`:140`（idle 遍必须 false）、`:178`（生产默认那遍必须 false）、`:216`（借之后必须 true）。
- `internal/ball/hotkey_live_test.go:61-62` 真机探针 `escBorrowProbe` 自己拿 `spareHKID` 向 Win32 注册裸 Esc（`:61-62` `pRegisterHotKey.Call(uintptr(b.hwnd), uintptr(spareHKID), uintptr(modNoRepeat), uintptr(vkEscape))`）——**探针也锚在同一枚裸键上**，配置一改它问的就不是借走的那把键。
- `internal/ball/hotkey_live_test.go:41-43` `liveHotkeys()` 把 cancel 定成 `Ctrl+Alt+V`（注释 `:32-40` 逐字解释为什么真机用例躲开裸键）⇒ 这批用例的"配置值"与"借到的值"**今天本来就是两枚不同的键**，改形ⓐ 之后两者的关系会翻转。
- 无关的那两枚同名常量引用（不动）：`internal/ball/hotkey_test.go:24-25` `{"Esc", 0, vkEscape, false}` / `{"Escape", 0, vkEscape, false}` —— 它们测的是 `ParseAccelerator` 本身，不测借用。

### 2.4 借还两处按"裸 Esc"写死的字面（逐处行号＋字面）

票面给的范围 `internal/ball/ball_windows.go:875-901` 比实际**窄**：借体是 `:873-889`、还体是 `:895-907`（`875-901` 把还体的 `:902-907` 切掉了，而那 6 行里正好有全仓唯一一处"归还时读配置"的 `:904`）。逐处：

| # | 落点 | 字面 | 按裸 Esc 写死？ |
|---|---|---|---|
| 1 | `internal/ball/ball_windows.go:878` | `if err := takeEsc(b.hwnd); err != nil {` | **是（间接）**：调用只传 hwnd，键由 callee 的字面量决定（`internal/ball/hotkey_windows.go:524`） |
| 2 | `internal/ball/ball_windows.go:881` | `b.hotkeyReport = b.hotkeyReport.withCancel(cancelFailedLine(err))` | **是**：`cancelFailedLine` → `cancelBorrowedLine()` → `Binding: escBorrowBinding`／`Acc: escBorrowAcc()`（`internal/ball/hotkey_windows.go:582-584`） |
| 3 | `internal/ball/ball_windows.go:886` | `b.hotkeyReport = b.hotkeyReport.withCancel(cancelBorrowedLine())` | **是**（同上） |
| 4 | `internal/ball/ball_windows.go:900` | `releaseEsc(b.hwnd)` | **否**：`internal/ball/hotkey_windows.go:553` 逐字 `func releaseEscWith(unreg unregisterFn) { unreg(hkCancel) }` —— 只按 id 丢，与借的是哪枚键无关 ⇒ 改形ⓐ 不需要动归还的第一性形状 |
| 5 | `internal/ball/ball_windows.go:904` | `b.hotkeyReport = b.hotkeyReport.withCancel(cancelIdleLine(b.cancelBinding))` | **否（这处吃配置）**：`b.cancelBinding` 来自 `:253`／`:821`；只用于报表分类，不注册 |
| 6 | `internal/ball/ball_windows.go:875`／`:885`／`:897`／`:901` | `if b.escTakenOver {`／`b.escTakenOver = true`／`if !b.escTakenOver {`／`b.escTakenOver = false` | 否（布尔位，字段名带 esc 但值与键无关；字段声明 `internal/ball/ball_windows.go:125`） |

**注释／自述层按裸 Esc 写死、改形后会变假的句子**（不是行为，但票 245 AC#3 那族"改成带条件的事实句"的规矩会管到它们）：
`internal/ball/hotkey_windows.go:7-15`（B1 头条，`:11-12` 「the production cancel binding IS the bare "Esc", DefaultHotkeys below」）、`:513-518`（`escBorrowBinding` 的"whatever the configured binding says"，**这句是明示的写死自述**）、`:521-523`、`:541-542`（「It is the ONLY path that registers the cancel id」）、`:547-552`（release 那段「re-registering the configured binding on the way out would put the production default (a bare Esc) straight back」）、`:250-258`（`HotkeyStandby` 文档）；`internal/ball/ball_windows.go:863-872`（`TakeEscForCancel` doc「temporarily binds Esc as the cancel key」）、`:891-894`（`ReleaseEscAfterSession` doc）、`:242-247`（开机注释）；`cmd/wisp/resident_ball_windows.go:127-130`（「the fourth slot, cancel, is left standby by internal/ball because its production binding is a bare Esc」）、`:132-140`（ticket 246 那段）；`cmd/wisp/resident_approval_windows.go:13`。
另有**面向用户的字面**：`internal/agent/approval/approval.go:90` `ChannelEsc: "按 Esc 键",`（通道枚举声明 `:55` `ChannelEsc Channel = "esc" // global Esc hotkey`）——改形ⓐ 之后卡片上这句会不会跟着配置变，属"用户看得见什么"那一格，AC#1 必答（本格只报落点）。


---

## §3 问题③：与票 245 裁定的边界

一句话结论：把借用体改成吃配置值**不碰**票 245 裁的那条「稳态不绑、只在确认那两三秒借那把键」——那条裁定的落点在 `registerAllWith` 与 `releaseEscWith`，不在加速器的来源上；碰得到的只有"借期是否仍为裸键"这一维，逐行给在 §3.2。

### 3.1 票 245 裁的是哪一句（原文落点）

### 3.2 改加速器来源会碰哪一行／不碰哪一行

---

## §4 量不到的格子

本程不跑 Go，因此凡是"要跑起来才知道"的读数一律不进本件；具体哪几格量不到、为什么、以及要什么尺才量得到，逐条列在 §4.1–§4.3。全称否定的那几处（§4.0）给尺的字面命令与正控。

### 4.0 全称否定清单（尺的字面＋正控）

### 4.1 量不到：需要真机 Win32 读数的三格

### 4.2 量不到：需要跑测试才能证的判据形状

### 4.3 量不到：票面引用了本仓不存在的形状

---

## §5 我推翻前人（含编排者 `A560`）哪几句

逐句指认：哪一句、原文在哪、我的反证 `file:line`、我为什么认为它过期或不成立。

---

## §6 交件名册

本件按节分枚提交，每枚 `git show --name-only` 只含本文件一枚路径；下面逐枚记 sha、时刻、覆盖的节。末枚记行数与字节数。
