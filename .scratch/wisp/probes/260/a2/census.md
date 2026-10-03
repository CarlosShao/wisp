# 票 260 AC#0 只读普查证据件（第二程，腿 `260-a2`）

判死问题：借用那一遍（`Confirming` 借那把取消键）**读不读 `[hotkey] cancel` 的配置值**——并在 258-r1 已落地的树上现量（a1 的锚 `47b38765` 早于 258-r1，其 `cmd/wisp` 侧结论与部分行号已过期，见 §6）。

路径口径：下面所有引用一律「仓库根相对全路径:行」（以 `D:\work\workspace\projects plans\Wisp\` 为根），不给裸文件名引用。

---

## §0 起手锚

- HEAD：`dffd9456d25ac11559f9f476b5afa848a9b0526b`（`git log -1 --format=%H`，分支 `dev`；提交时刻 `2026-10-03 21:07:28 +0800`）。a1 锚＝`47b38765`，本程起手比它新若干枚；产码差集见 §6 第 1 条。
- 时刻：`2026-10-03 21:0x +0800` 起、本程内现跑。
- `git status --porcelain internal cmd` 逐字（产码子集，起手）：

```
?? cmd/wisp/resident_hotkey_258_windows_test.go
?? cmd/wisp/resident_hotkey_live_258_windows_test.go
```

两枚未跟踪件＝258 落地腿的在飞测试件（前人件），具名登记、不评判；除此之外产码子集干净。终态复跑见 §8。

- 口径：**本程零 Go 命令**（`go test/build/vet/run/list` 一律未跑；本机 cmd/wisp 有测量收尾在飞，任何 Go 编译都会互洗读数）。工具＝Read/Grep/Glob 与 `grep`/`sed`/`git log|diff|show|status|add|commit`。
- 禁令执行：`frontend/**` 与 `design/**` 零读零写零转述（本件不含其任何内容）；`docs/**`／`internal/**`／`cmd/**` 零写；三枚冻结件、`internal/observe/thresholds.go`、golden、`tools/d22scan/allowlist.txt` 未动；工单票面未动。
- 与 a1 件的关系：a1 的**产码读数**本程逐处现跑复核，不复述当凭据；它判不动的格子（258-r1 桥落到借用侧没有）正是本程的主任务。a1 未覆盖的格子＝258-r1 之后的 `cmd/wisp` 接线全貌、245 裁定原文在 HEAD 的逐字对读（a1 只引了 `:21`/`:47`，本程补齐 `:26`/`:45`/`:54`/`:64`/`:72` 与 PLAN 冻结行）。

---

## §1 借用遍逐处 file:line 名册

**一句话结论：不读。** 借用体注册的加速器是包内两枚常量的字面组合；配置值在这条支路上**零环携带**——链条四枚签名（`TakeEscForCancel()`→`takeEsc(hwnd)`→`takeEscWith(unreg, reg)`→`escBorrowAcc()`）的形参里键的数量是零。258-r1 之后配置值**已经活着进了球内**（`b.cancelBinding`／`b.boundCfg`，见 §4），但借用体对这两枚字段**一眼都不看**。

### 1.1 借用链逐处（现跑 `grep -n` 于 HEAD `dffd9456`）

| 环节 | 落点 | 盘上字面（逐字） | 携带键信息？ |
|---|---|---|---|
| 借用的唯一产码触发点 | `cmd/wisp/resident_approval_windows.go:438-441` | `u.ra.cards.Record(p)` → `if p.Level == "L1" {` → `b.TakeEscForCancel()` | 无参；守卫＝仅 L1 卡（L2 队列卡不借，`cmd/wisp/resident_ball_windows.go:215` 注释逐字 "an L2 queue card does not borrow the key at all"） |
| 球侧借体 | `internal/ball/ball_windows.go:873`／`:878` | `func (b *Ball) TakeEscForCancel() {` → `if err := takeEsc(b.hwnd); err != nil {` | 只传 hwnd；**不读 `b.cancelBinding`／`b.boundCfg`**（函数体 `:873-889` 通读无一处引用这两枚字段） |
| 借的原语壳 | `internal/ball/hotkey_windows.go:543-545` | `func takeEsc(hwnd windows.HWND) error { return takeEscWith(hotkeyUnregisterer(hwnd), hotkeyRegisterer(hwnd)) }` | 只传 seam，无键 |
| 借的原语体 | `internal/ball/hotkey_windows.go:531`／`:532-533` | `func takeEscWith(unreg unregisterFn, reg registerFn) error {` → `unreg(hkCancel)` → `if err := reg(hkCancel, escBorrowAcc()); err != nil {` | 只有 id＋seam；**加速器来自 `escBorrowAcc()` 零参调用** |
| 加速器来源（判死行） | `internal/ball/hotkey_windows.go:524` | `func escBorrowAcc() Accelerator { return Accelerator{Mods: modNoRepeat, VK: vkEscape} }` | **零参、返回字面量** |
| 两枚常量的家 | `internal/ball/hotkey_windows.go:116`／`:118` | `modNoRepeat = 0x4000`／`vkEscape = 0x1B` | — |
| 打到 Win32 的那一手 | `internal/ball/hotkey_windows.go:391-392` | `r, _, err := pRegisterHotKey.Call(uintptr(hwnd), uintptr(id), uintptr(acc.Mods), uintptr(acc.VK))` | acc 自上一环 |
| 产码注册调用点口径 | `internal/ball/hotkey_windows.go:391`＋`internal/ball/win32_windows.go:42` | `pRegisterHotKey.Call(...)`／proc 声明 | `grep -rn "pRegisterHotKey" --include=*.go internal cmd | grep -v _test` 现跑＝**只这两枚**（票面现量 1 在 HEAD 仍成立；射程外两枚见 §6 第 6 条，与借用无关） |

借到手的值＝`Accelerator{Mods: 0x4000, VK: 0x1B}`＝裸 Esc＋`MOD_NOREPEAT`。id＝`hkCancel = 3`（`internal/ball/hotkey_windows.go:32`）。**判法不是"我没找到"**：`escBorrowAcc` 是借用体取加速器的唯一来源（引用者全仓只有 `:533` 注册与 `:584` 报表，`grep -n "escBorrowAcc" internal/ball/hotkey_windows.go` 现跑＝`:521` 注释／`:524` 定义／`:533`／`:584` 四命中），形参为空、返回体是常量字面，故 `HotkeyConfig.Cancel` 与它之间**不存在数据通路**；借用体内如果哪天有人想读配置，唯一现成的取值字段就是 `b.cancelBinding`／`b.boundCfg.Cancel`（§1.2），而 `TakeEscForCancel` 函数体（`internal/ball/ball_windows.go:873-889`）现读**通体没有**这两个名字。

### 1.2 配置值此刻活在借用体旁边哪两枚字段里（258-r1 之后）

| 字段 | 落点 | 写入点 | 读口 |
|---|---|---|---|
| `b.cancelBinding string` | `internal/ball/ball_windows.go:126`（注释逐字 "configured cancel binding, for B1 release"） | `:253`（开机 `= b.opts.Hotkeys.Cancel`）／`:821`（rebind `= cfg.Cancel`） | 全仓产码读它只有一处：`:904` `cancelIdleLine(b.cancelBinding)`（**归还时**，只进报表分类，不注册） |
| `b.boundCfg HotkeyConfig` | `internal/ball/ball_windows.go:124`（注释 "the set the last pass was given (rebind updates it)"） | `:254`／`:819` | `ConfiguredHotkeys()` `:840-844`（uiRun 读）——借用体不调它 |

⇒ 配置来的 cancel 串在包内有**两处被"读"**（idle 分类 `internal/ball/hotkey_windows.go:428-450` 的 `ParseAccelerator(bind)` `:435`——解出来的 `Accelerator` 被丢弃 `_`，只用 err 分类；归还分类 `internal/ball/ball_windows.go:904`），**零处被"注册"**。缺的就是 §1.1 那一环。

### 1.3 idle 那一遍（对照组，现行号复核）

- 出厂注册：`internal/ball/ball_windows.go:255` `registerAll(b.hwnd, b.opts.Hotkeys)` → `internal/ball/hotkey_windows.go:453` `registerAllWith`。
- 跳过：`internal/ball/hotkey_windows.go:457-463`，`:457` `if p.id == hkCancel {`、`:458` 注释逐字 "The one slot an idle ball must not register. takeEscWith is the"、`:462` `rep.bindings = append(rep.bindings, cancelIdleLine(bindings[i]))`、`:463` `continue`——**按 id 跳、与配置值无关**。
- 钉子：`internal/ball/hotkey_status_test.go:205-206` 逐字 `if reg.attemptsOf(hkCancel) != 0 { t.Fatalf("idle pass attempted cancel: ...") }`。
- 报表侧 `Standby` 不进 `Problems()`：`internal/ball/hotkey_windows.go:365` `case HotkeyLive, HotkeyDisabled, HotkeyStandby:`（HEAD 复核 ✓）。

### 1.4 写死的第二张脸：报表与日志

- `internal/ball/hotkey_windows.go:519` `const escBorrowBinding = "Esc"`（注释 `:513-518` 自述 "the bare Esc, **whatever the configured binding says**"）；引用＝`:535`（借失败日志 `"binding", escBorrowBinding`）＋`:582`（`cancelBorrowedLine()` 的 `Binding:`）。
- `internal/ball/hotkey_windows.go:578-586` `cancelBorrowedLine()`：`Binding: escBorrowBinding`＋`Acc: escBorrowAcc()`；`:591-596` `cancelFailedLine(err)` 由它派生；两处调用＝`internal/ball/ball_windows.go:881`（借失败）／`:886`（借成功）。
- 借成功日志：`internal/ball/hotkey_windows.go:534-535` 只在失败时打；借期报表把 cancel 行整行换成裸 Esc 字面——配置串在借期内只从 `ConfiguredHotkeys()`（`internal/ball/ball_windows.go:840-844`）还看得见。

### 1.5 判死那一句

**不读。** 判死行＝`internal/ball/hotkey_windows.go:524`。判法＝唯一加速器来源零参、返回常量字面（§1.1），且借用体函数体（`internal/ball/ball_windows.go:873-889`）与它下面那条原语链对 `b.cancelBinding`／`b.boundCfg` **零引用**——258-r1 已把配置值喂到字段旁边，借用体没伸手。

---

## §2 缺零件清单（若要吃配置）

一句话结论：**解析器、修饰键、取值字段、Win32 手四样全在包里**；缺的是把 `escBorrowAcc()` 的字面来源换成配置值所需的**一枚参数化＋一条新分支＋两层字面翻新**。逐枚：

| # | 零件 | 现状落点 | 要动什么 |
|---|---|---|---|
| P1 | `Accelerator` 类型与 `ParseAccelerator`（串→加速器） | `internal/ball/hotkey_windows.go:105-108`（类型）／`:125-158`（解析器；`:142` `acc.VK = vkEscape`、`:156` `acc.Mods |= modNoRepeat`——**每条解析结果自带 NOREPEAT**，与借用手写的 `Mods: modNoRepeat` 同形） | 不用新造；直接用 |
| P2 | 修饰键常量 | `internal/ball/hotkey_windows.go:112-116`（`modAlt/modControl/modShift/modWin/modNoRepeat`） | 不用动 |
| P3 | 取值字段 | `b.cancelBinding`／`b.boundCfg.Cancel`（§1.2），258-r1 后**已被配置与 rebind 喂活** | 借用体去读它——最短形＝`TakeEscForCancel` 在 uiRun 闭包里读 `b.cancelBinding`→`ParseAccelerator`，**导出签名不动**（外部调用者 `cmd/wisp/resident_approval_windows.go:441`、`cmd/balldebug/main.go:635` 一字不改） |
| P4 | `escBorrowAcc()` 参数化 | `internal/ball/hotkey_windows.go:524` 零参返回字面量 | 改成吃一枚 `Accelerator`（或新增带参变体、旧函数留给测试）；引用者只有 `:533`／`:584` 两处 |
| P5 | 键怎么进借用体 | `takeEscWith(unreg, reg)` `:531`／`takeEsc(hwnd)` `:543-545` 均无键参 | 两形：改签名（波及 `takeEsc`＋测试 4 处调用＝`internal/ball/hotkey_status_test.go:210`／`:259`／`:270`／`:340`）——**或**新增 `takeEscAcc` 变体不动旧函数。P3＋P4＋P5 合起来的最小形＝"借时读字段、新变体带 acc" |
| P6 | **不可解析配置的借用分支（a1 未点名的新零件）** | 今天不存在——借用从不读串，无失败分支 | 用户把 cancel 配成 `ParseAccelerator` 拒收的串（`internal/ball/hotkey_windows.go:145-155` 未知键名／无键两处报错）时借什么：回退裸 Esc？报表 `unparsable`＋不借？**必须先定**，否则形ⓐ 落地即引入"借一半"状态 |
| P7 | 报表/日志字面翻新 | `escBorrowBinding` `:519`、`:535`、`cancelBorrowedLine` `:578-586`、`cancelFailedLine` `:591-596` | 借期报表与失败日志的键名要跟着实际借的键走，否则报表说谎 |
| P8 | 注释诚实层（改形后变假的句子族） | a1 §2.4 点名 11 处（`internal/ball/hotkey_windows.go:7-15`／`:513-518`／`:521-523`／`:541-542`／`:547-552`／`:250-258`；`internal/ball/ball_windows.go:242-247`／`:863-872`／`:891-894`；`cmd/wisp/resident_ball_windows.go:203-205`（现行号，原 `:127-130`）／`cmd/wisp/resident_approval_windows.go:13`） | 逐句改成带条件的事实句（245 AC#3 同族规矩） |
| P9 | 测试锚重钉 | 4 枚裸 Esc 锚：`internal/ball/hotkey_status_test.go:223`（VK＋Mods 双判）／`internal/ball/hotkey_live_test.go:107`（只判 VK）／`internal/ball/hotkey_status_test.go:85-91` `bindsEsc()`（调用 `:140`／`:178`／`:216`）／真机探针 `internal/ball/hotkey_live_test.go:61-62`（自己拿 `spareHKID` 向 Win32 注册裸 Esc） | 245 `:64` 的顺序纪律"先改钉、再改行为"管着：ⓐ 落地必须先交钉子改动 |

**最贵的一枚＝P9（4 枚裸 Esc 锚）捆绑 P8（注释族）**：不是代码量贵，是**顺序贵**——245 `:64` 已裁"先改钉、再改行为"，且四枚钉今天还锚在枚数/VK 字面上（`grep -n "got.VK != vkEscape" internal/ball/hotkey_status_test.go internal/ball/hotkey_live_test.go` 现跑＝`:223`／`:107` 两命中），落地腿的第一枚 commit 必然是钉子 commit 而不是行为 commit。P6 是**唯一的新决策**（其余全是接线），但它小。

**不需要动的**：WM_HOTKEY 路由按 id 不按键（`internal/ball/ball_windows.go:662` `case hkCancel:` → `fire(b.opts.Events.OnCancelHotkey)`）；归还只按 id 丢（`internal/ball/hotkey_windows.go:553` `func releaseEscWith(unreg unregisterFn) { unreg(hkCancel) }`）；`registerFn` seam 形状 `:386`（acc 本来就是它的第二参）。**改加速器来源不改取消的落地路径。**

---

## §3 与票 245 裁定边界的对读

**一句话结论：撞不到"已生效裁定"的字面边界；撞到的是 245 里那格"未勾"的 AC#6② 的方向，以及它 `:26` 具名的那半代价的登记形状。** 票 260 面 `:23` 的禁令（要动 245 的边界＝人工批准）因此不被形ⓐ 的最小形触发。

### 3.1 票 245 裁定原文（本程逐字现读，HEAD 下行号）

- 裁定句 `.scratch/wisp/issues/245-the-bare-esc-cancel-default-is-a-global-hotkey-so-the-resident-process-steals-esc-from-every-other-app.md:21`：
  > **走乙：稳态不绑 cancel 那一枚，只在进入 `Confirming` 时 `takeEsc`、离开时 `releaseEsc`。**
- 三支理由 `:22-24`（不动冻结文字／复用现成机制／丙形与既有设计相反——本程复读 `:23` 逐字含 "「a live hotkey the user can see, never a dead silent one」⇒ 靠"留空"关不掉"）。
- 半代价 `:26` 逐字：
  > **⚠ 乙形留下的那半代价，具名不藏**：**确实在"等人批"那几秒里，别的程序的 Esc 仍会被抢走**（窗口只有 2–3 秒；`SPEC-06` 的 L1 定义）。这半条**今天不解决**，登记在本票末尾当残余。
- 收件状态（HEAD 现读）：`:54`（A480 收表：AC#1/#2/#3/#5 翻勾，AC#4/#6 不翻）＋`[x] AC#1`＝`:59`、`[x] AC#2`＝`:60`、`[x] AC#3`＝`:61`、`[ ] AC#4`＝`:62`、`[x] AC#5`＝`:63`、`[ ] AC#6 两半都不翻`＝`:64`、`[ ] AC#7/AC#9`＝`:66`／`:68`；票尾 `:72`「勾 4／未勾 5 ⇒ 不加 `-done`」。
- **待生效**的 AC#6② `.scratch/wisp/issues/245-...md:47` 逐字（关键半句）：
  > ⇒ 用户把 `[hotkey] cancel` 配成**带修饰键**的组合时**应当常绑**（那不会吞别的程序…）；配成**裸键**时仍走借用。
  该格**框未勾**＝待生效判形，不是盘上现成规矩。
- 顺序纪律 `:64` 内含逐字「**顺序必须"先改钉、再改行为"**」。
- 已生效裁定的实现落点（两处，HEAD 复核）：`internal/ball/hotkey_windows.go:457-463`（idle 按 id 跳）＋`internal/ball/hotkey_windows.go:553`（归还只丢不 rebind，理由注释 `:547-552`）。
- 冻结文字：`docs/PLAN.md:3082`（D43 第 22 行）逐字含「否决（**单击球 / `Esc` / KWS 否决词 / 面板拒绝**，B1）」——C12 冻结，本件零写；形ⓐ 默认档不换键名 ⇒ 无需触碰。

### 3.2 逐改动对读（撞/不撞）

| 改动 | 撞不撞 245 已生效句（`:21`） | 依据 |
|---|---|---|
| `escBorrowAcc()` 参数化＋借用体读配置（P3–P5） | **不撞** | `:21` 裁的是「稳态不绑／借—还成对」这**对机制**，点名的是 `takeEsc`／`releaseEsc` 两枚**函数**，没有点名"借的那把键必须是裸 Esc"——"裸 Esc"身份活在设计理由（`:14` 表与 `:26`）与代码注释里，不在裁定句的字面里。两个实现落点（`:457-463`／`:553`）判 id 与"还不还"，与 acc 从哪来无涉 |
| 默认档行为 | **不撞，逐位相同** | `ParseAccelerator("Esc")`＝`{Mods:0x4000, VK:0x1B}`（`:141-142` esc 分支不加修饰键＋`:156` 补 NOREPEAT）＝`escBorrowAcc()` 手写同一对；反向钉 `internal/ball/hotkey_status_test.go:190-192`（默认 `Cancel:"Esc"` 不许改）不动 ⇒ 默认档**一发出册都不变** |
| 归还那一手 | **不撞行为，撞注释** | `:553` 只 `unreg(hkCancel)`，钉在 `internal/ball/hotkey_status_test.go:239-242`（逐字 "never re-bind the configured binding"）；注释 `:547-552` 的 "the production default (a bare Esc)" 在用户配组合键那支不再覆盖真实形状＝注释过期类（P8） |
| 用户配**组合键**后的借期 | **撞 AC#6②（`:47`，未勾）的方向** | `:47` 说组合键"**应当常绑**"；形ⓐ 的最小形是"组合键也在 Confirming 借 2–3 秒"——**借但不常绑**。这不是推翻 `:21`（稳态仍不绑），但与 `:47` 的终态方向相抵；且借用户自己的组合键 2–3 秒＝把 `:26` 那半代价从"裸 Esc"扩展到"用户的键"，**残余登记要五字段更新**。⛔ 该格归 245 未勾框，260 不许替它定（票 260 面 `:23` 明写"要动它的边界＝人工批准"）——**本格只报这个相抵点，不判形** |
| 用户可见文案 | 落点在、判形不在本格 | `internal/agent/approval/approval.go:90` 逐字 `ChannelEsc: "按 Esc 键",`（枚举 `:55`）＝卡片"按 Esc 取消"现役出处；配置改键后这句与实际借到的键漂移，是 AC#1"用户看得见什么"要答的第一枚 |

---

## §4 与 258-r1 桥的关系判定

**一句话判定：链已通——通到借用侧的"取值字段"为止；借用的"注册路径"是链上唯一没接的一段。** 准确说法不是"链不到借用那一侧"，而是"链把配置值喂到了 `b.cancelBinding`／`b.boundCfg` 这两枚字段（借用体伸手可及处），借用体没伸手"。

链逐处（全部 HEAD 现跑）：

1. **取值闭包**：`cmd/wisp/resident_windows.go:175-189` `hotCfg258`（每次调用现读 `config.LoadFile`；`:184` `h := c.Hotkey`、`:188` 四字段映射）、`:198-209` `hotReload258`（每 tick 现读）。
2. **进球**：`cmd/wisp/resident_windows.go:210` `startResidentBall(rt.Registry, ra.vetoByEsc, hotCfg258, hotReload258, …)` → `cmd/wisp/resident_ball_windows.go:260` `cfg, provenance := residentBallHotkeyChain258(hotCfg)`（`:183-194`，`ApplyHotkeyDefaults` 在 `:194`）→ `:263` `ball.New(ball.Options{` → **`:269` `Hotkeys:  cfg,`**（a1 记的 `:171 DefaultHotkeys()` 已不存在）→ 开机写 `internal/ball/ball_windows.go:253-254`（`b.cancelBinding`／`b.boundCfg`）。
3. **桥**：`cmd/wisp/resident_ball_windows.go:305-324`（`if hotReload != nil && rb.b != nil`）→ `:306` `bridge := ball.NewHotkeyReloader(b, b.ConfiguredHotkeys(), hotReload)` → `:308-309` `reg.Spawn("hotkey-listener", "ball", …, func(ctx){ bridge.Run(ctx, time.Second) })`（1s poll）→ armed 日志 `:322-324`。
4. **rebind 回写**：`internal/ball/hotkey_reload.go:126` `Run` → `:81` `Check` → `:96` `rep := r.binder.RebindHotkeys(cfg)` → `internal/ball/ball_windows.go:813-825` `RebindHotkeys`（**经 `b.uiRun` post-and-wait**，派单那句核实成立）：`:819` `b.boundCfg = cfg`、`:821` `b.cancelBinding = cfg.Cancel`、`:822` `b.escTakenOver = false`（在飞借用被 rebind 丢弃＝245 残余，注释 `:807-812` 逐字 "A host that rebinds while a card is waiting therefore has to call TakeEscForCancel again … no caller does that today"）。
5. **借用侧断点**：`internal/ball/ball_windows.go:873-889` `TakeEscForCancel` 对 `:126`/`:124` 两枚字段零引用（§1.1）→ `internal/ball/hotkey_windows.go:524` 字面量。

**推论（对 AC#1 的射程含义）**：形ⓐ 若走"借时读 `b.cancelBinding`"，**258 的链零改动自动获益**——rebind 已把新值写进同一枚字段（`:821`），rebind 丢弃在飞借用后下一发卡片再借时读到的就是新配置；258-r1 落地腿不需要为此补任何线。反向也成立：**258-r1 落地本身没有（也不应该）顺手接借用**——它的面（258）与 260 的面被两份票面明文"相邻但不互相吞并"（票 260 `:31`），桥的全部效果截至 HEAD 就是"三枚非 cancel 槽跟着配置活"。

**258-r1 顺带改变的可观察面（本程新发现，归 §5/§6）**：258-r1 之后 `internal/ball/hotkey_windows.go:444-446` 那条 idle 日志的 `"binding"` 字段**已经是配置值**（`cancelIdleLine(bindings[i])` 的 `bindings[2]`＝cfg.Cancel＝配置），而句子字面仍写死 "Esc is borrowed only during Confirming"——a1 §5.6 预言的"自相矛盾句"在 HEAD **已是现在时**：用户配 `cancel="Ctrl+Alt+K"` 的出厂主机今天就会打出「cancel hotkey left unbound while idle: **Esc** is borrowed … binding=**Ctrl+Alt+K**」。

---

## §5 AC#1 两形代价预告（各一句，⛔ 本程不选形）

- **ⓐ 借用读配置**：零件九枚（§2 P1–P9）全在包内、最小形＝借时读 `b.cancelBinding`→`ParseAccelerator`→带 acc 变体（导出签名与外部调用者零改），但**落地顺序被 245 `:64` 钉死为先交 4 枚裸 Esc 锚的重钉 commit（P9）再动行为**，外加 P6 那枚"不可解析配置借什么"的新分支要编排者先定、P8 注释族 11 处翻新、以及 245 `:26` 残余登记的五字段更新（借期吞的从"裸 Esc"变"用户的键"）——默认档一发出册逐位不变（§3.2），用户改键后借期内报表/日志说新键名，但卡片"按 Esc 键"（`internal/agent/approval/approval.go:90`）会与实际键漂移、须同票处理。
- **ⓑ 维持写死＋说到明处**：产码零改（`escBorrowAcc`/借还链原样），代价集中在**两处已在打的句子改字面**——`internal/ball/hotkey_windows.go:444-446` idle 日志（258-r1 后 `binding` 字段已是配置值、句内写死 "Esc"，出厂改键用户今天就能看见自相矛盾，ⓑ 必须改它否则"说到明处"自己先说谎）与 `cmd/wisp/resident_ball_windows.go:338` verdict／回执侧补一句"cancel 配置项今日不生效、确认卡借用固定为裸 Esc"（⚠ 同时不能再留"四枚都以配置为准"这类句子——258-r1 后 summon/mute/panel **已**吃配置，假话只剩 cancel 一枚，这句要按"三枚吃配置＋cancel 借用固定"的形状说），用户看得见＝日志与回执明说"你配的 cancel 键不会成键"；另 `main.go:25` "its four global hot keys" 那句 245 AC#6① 的旧账仍开着（HEAD `cmd/wisp/main.go:25` 现读仍在），ⓑ 顺手改它最省。

---

## §6 推翻清单（票面／a1／本派单，每句都是待验断言；共 7 指）

1. **a1 §1.2/§5.3「出厂主机 `Hotkeys: ball.DefaultHotkeys()` 在 `cmd/wisp/resident_ball_windows.go:171`、`config.HotkeySection` 在出厂腿从未被读、'三枚成立'只在 balldebug」——已被 258-r1 推翻（本程主推翻）**。反证：`cmd/wisp/resident_ball_windows.go:269` 逐字 `Hotkeys:  cfg,`，cfg 来自 `residentBallHotkeyChain258`（`:183-194`）← `hotCfg258`（`cmd/wisp/resident_windows.go:175-189`，现读 config.toml）。HEAD 下 summon/mute/panel **已**跟着配置走，`cmd/wisp/config_readers_255.go:114-115`/`:123` 的名册行也已改口（"Since ticket 258 (form A) … binds from [hotkey]"）。**票 260 面 `:11` 的推论（用户改 cancel 则该键任何路径不成键）经本程复核在 HEAD 仍成立**，但成立原因收窄为"cancel 一枚"：idle 按 id 跳（§1.3）＋借用写字面（§1.1）；"只对三枚成立"那句在 HEAD 恰好**变回真话**（a1 说它连三枚都不成立——那在 a1 的锚上对，今天不对）。
2. **票面 `:10` 两枚行号（a1 已指，本程 HEAD 复核维持）**："`:521` `escBorrowAcc`" ⇒ `:521` 是注释首行、声明在 `internal/ball/hotkey_windows.go:524`；"借还成对体在 `internal/ball/ball_windows.go:875-901`" ⇒ 借体 `:873-889`、还体 `:895-907`（`875-901` 把含全仓唯一"归还读配置"的 `:904` 关在区间外）。
3. **票面 `:12` 的 gate 行号已漂**：`case <-deadline` 现 `internal/agent/approval/gate.go:325`（票写 `:319`）、`ANSWER-EXPIRED … timeout->execute` 日志现 `:338`（票写 `:332`）、`return tools.AnswerTimeout` 现 `:344`；`internal/agent/approval/queue.go:116` `DefaultL1Window = 3 * time.Second` 未漂（夹逼 `:120`/`:122` 2–3s 仍在）。**丢借用＝零症状的结论本身复核成立**：`internal/ball/hotkey_windows.go:365` `Standby` 不进 `Problems()` ✓。
4. **a1 §5.6 的将来时态已过期（现在时）**：「一旦 228/258 把 `[hotkey]` 接进 `Options.Hotkeys`，同一行会开始打印自相矛盾句」——HEAD 已是：`internal/ball/hotkey_windows.go:444-446` 的 `binding` 字段＝配置值、句内 "Esc" 写死（§4 末条）。该句从"预言"升级为"现役缺陷读数"，且是形ⓑ 的头号必改落点。
5. **本派单「起点＝票面现量 2 的四处」作为普查边界不够**：四处（`hotkey_windows.go:118`/`:458`/`:521`/`ball_windows.go:875-901`）只盖 `internal/ball` 内部；借用链的**触发点**（`cmd/wisp/resident_approval_windows.go:438-441`）与 **258 桥的取值侧**（`cmd/wisp/resident_windows.go:175-209`、`resident_ball_windows.go:260`/`:269`）是 AC#0 答案的一半，本程已补进 §1/§4。
6. **票面 `:9`「全仓产码只有一枚调用点」的量词仍以尺射程为准（a1 §5.1 指认，本程 HEAD 复核其两枚反证均在）**：`cmd/wisp/testdata/esclistener/main.go:71`/`:243`（`procRegisterHotKey`，模式 `pRegisterHotKey` 结构上漏它）与 `scripts/spike/common/winshell.go:45`/`:437`（射程外）。对借用判定无影响（那两枚都与 `hkCancel` 无关），只作量词更正。
7. **a1 §2.2 第 3 条「`takeEscWith` 加键参会波及它的五处测试调用」的枚数不对**：`grep -n "takeEscWith(" internal/ball/hotkey_status_test.go` 现跑＝**4 处**（`:210`/`:259`/`:270`/`:340`），第五处是产码包装 `takeEsc`（`internal/ball/hotkey_windows.go:543-545`）不是测试。枚数更正不改变结论（签名变更面仍是最大的一枚）。

**复认成立（不推翻，HEAD 现跑）**：`internal/ball/hotkey_windows.go:458` "The one slot an idle ball must not register" ✓；`:116`/`:118` 常量对 ✓；`internal/ball/hotkey_windows.go` 自 a1 锚**零改动**（`git diff 47b38765..HEAD --name-status` 不含该文件）⇒ a1 对它的全部行号原样有效；`internal/ball/ball_windows.go` 仅 +18 行（尾部追加 `DebugRegisterHotkeySquat`/`DebugUnregisterHotkeySquat` 两枚 seam，`git diff` 逐行核过），借还区行号未漂 ✓；`docs/PLAN.md:3082` D43 第 22 行含 "`Esc`" 逐字未动 ✓。

---

## §7 量不到的格子＋复量法

1. **借的那把键在真桌面上今天是否真成键**（RegisterHotKey 运行时事实）：尺在仓内＝`internal/ball/hotkey_live_test.go:57-71` `escBorrowProbe`（`:61-62` 拿 `spareHKID` 实注册）与 `cmd/wisp/resident_hotkey_live_258_windows_test.go`（未跟踪在飞件）；`winlive` build tag 不进 CI（`grep -n "tags" .github/workflows/ci.yml` 现跑＝0 命中），读数只有真机有。复量法＝真机 `-tags winlive` 跑上述两枚；本程禁 Go ⇒ 量不到。这格归 AC#2 的"会响的尺"，本程不预判。
2. **四枚裸 Esc 锚在形ⓐ 下的红/绿**：运行事实，本程零 Go ⇒ 量不到；静态只能给锚的判据形状（§2 P9：`:223` VK+Mods 双判、`:107` 只判 VK、`bindsEsc()` 三调用、探针自锚裸键）。复量法＝钉子 commit 后 `go test ./internal/ball/...`（由有 Go 权的腿跑）。
3. **真机开机时 provenance 三形的实际落点**（config 可读/不可读/无视图，`cmd/wisp/resident_ball_windows.go:183-194`）：要真启动常驻进程；复量法＝真机起 `wisd`/resident 后读 `:322-324` armed 日志与 `:338` verdict 的 `hotkeys from %s` 字段。
4. **用户改完配置后借期实际借到哪枚键**：今天恒裸 Esc（静态判死），但"改配置→rebind→下一发卡再借"的全链运行读数需要真机；复量法＝改 `config.toml` `[hotkey] cancel`→等 1s poll→发一张 L1 卡→读 `internal/ball/hotkey_windows.go:444-446` 日志与 `HotkeyReport()` cancel 行。
5. **面板侧取消键现役显示形状**：本程只留名不判（AC#1 落地腿现读）；`internal/panel` 此刻不在本程读数射程内（禁令只列了 frontend/design，但 253 写腿历史在飞，判它会读半成品）。
6. **行号漂移复量法（给下一程）**：`grep -n -E "escBorrowAcc|escBorrowBinding|func takeEsc|func releaseEsc|hkCancel" internal/ball/hotkey_windows.go`；`grep -n -E "TakeEscForCancel|cancelBinding|boundCfg|RebindHotkeys" internal/ball/ball_windows.go`；`grep -n -E "Hotkeys:|hotCfg258|hotReload258|NewHotkeyReloader" cmd/wisp/resident_windows.go cmd/wisp/resident_ball_windows.go`。

---

## §8 自查

- **零 Go 命令自证**：本程全部工具调用＝Read/Grep/Glob 与 `date`/`git log|diff|show|status|add|commit`/`grep`/`sed`/`ls`/`wc`；没有一枚 `go build`/`go vet`/`go test`/`go run`/`go list`。
- **写面自证**：本程只写 `.scratch/wisp/probes/260/a2/census.md`（本件）＋`.scratch/wisp/probes/260/a2/msg-01.txt`＋`.scratch/commit-msg-260a2.txt`（两枚说明文本，临时件只建不删、msg-01 不进 commit）；`git diff --name-only` 产码子集为空；commit 全部显式单文件 pathspec、零 `add -A`、零 amend/reset/rebase/stash/checkout/clean、零 push。
- **禁令自证**：`frontend/**`/`design/**` 零打开零转述；`docs/**` 只读两处（`PLAN.md:3082` 冻结行对读）；三枚冻结件、`internal/observe/thresholds.go`、golden、`tools/d22scan/allowlist.txt` 未动；`^- [ ]` 勾框零触碰（票 260 与 245 的框均未动）。
- **完整性格查**：§1–§8 无"（待填）"占位；每枚结论带 file:line 或具名"量不到"（§7）。
- **交件名册**：骨架 commit `9fb444cd`（§0＋节框）；本枚＝§1–§8 满件；终态读数与满件 sha 见交件消息（满件 commit 后如需回填名册，追加枚、不改写前枚）。
- **终态产码子集 status**（本件落笔后复跑）：

```
?? cmd/wisp/resident_hotkey_258_windows_test.go
?? cmd/wisp/resident_hotkey_live_258_windows_test.go
```

与起手一致：两枚 258 在飞测试件仍未跟踪，本程零产码写入、零产码改动。
