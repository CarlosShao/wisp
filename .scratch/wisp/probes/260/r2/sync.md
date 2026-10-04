# 票 260 · 260-r2 证据件 —— 把一枚因编排者裁定而过期的既有真窗判据，同步到已裁的行为上

本轮唯一一格：`internal/ball/hotkey_live_test.go`（`//go:build windows && winlive`）里那枚
`requireEscBorrowed` 的"借到的必是裸 Esc"断言，在形ⓐ（确认窗借用吃 `[hotkey] cancel` 配置值）
落地之后变成对旧行为的背书。裁定＝**期望值跟着种子走**，不跟着写死常量走；
默认档零漂移那一格不许丢；⛔ 不许删断言 / 不许 `t.Log` 化 / 不许 `t.Skip`。

本枚 = 骨架 commit（先落"我读到的冲突原文逐字"＋"名册"＋门禁基线）。
§4 之后的「改后逐字／终值读数」在本轮后半程的 commit 里展开；写这一句是为了让"哪些节已满、
哪些节还没满"是**具名状态**而不是一句空格。

锚点：起手 HEAD `5ad8ec27`（`dev`）。260-r1 四枚 = `eb2c0173`／`d9bb6817`／`b1e59d63`／`ba2a8358`。

---

## 1. 我读到的冲突原文（逐字，2026-10-04 现读；行号对得上，零漂移）

`internal/ball/hotkey_live_test.go:41-43`（种子）：

```go
func liveHotkeys() HotkeyConfig {
	return HotkeyConfig{Summon: "Ctrl+Alt+Q", Mute: "Ctrl+Alt+M", Cancel: "Ctrl+Alt+V", Panel: "Ctrl+Alt+B"}
}
```

⇒ 编排者给的 `:42` **逐字命中**（`Cancel: "Ctrl+Alt+V"` 就在 42 行；函数名行是 41）。
VK 值：`'V'` = `0x56`，mods = `MOD_ALT|MOD_CONTROL|MOD_NOREPEAT` = `0x0001|0x0002|0x4000` = `0x4003`。
（同包既有钉：`internal/ball/hotkey_cancel_borrow_260_test.go:338` 逐字
`want {mods 0x4003, VK 0x56}` —— 260-r1 已经把这一对钉在默认档，我复跑到该文件绿。）

`internal/ball/hotkey_live_test.go:105-115`（断言体，含前后两行上下文）：

```go
		t.Fatalf("after the borrow the live set holds %d entries, want 4: %+v", len(rep.Live()), rep.Bindings())
	}
	if got := rep.Live()[hkCancel]; got.VK != vkEscape {
		t.Fatalf("the borrowed cancel slot holds VK 0x%X, want VK_ESCAPE 0x1B: %+v", got.VK, rep.Bindings())
	}
	if free, err := escBorrowProbe(t, b); err != nil {
		t.Fatalf("the bare-Esc probe could not run: %v", err)
	} else if free {
		t.Fatal("ticket 245 RED: the ball says it borrowed Esc, but the desktop still has it free - " +
			"the borrow registered on some other key, not on VK_ESCAPE")
	}
```

⇒ 编排者给的 `:107-108` **逐字命中**（107 = `if got := rep.Live()[hkCancel]; got.VK != vkEscape {`，
108 = 那句 `the borrowed cancel slot holds VK 0x%X, want VK_ESCAPE 0x1B`）。

**我这一轮多读出来的一处**（编排者只点了 107-108）：同一个 helper 的**第二格** `:110-115` 也建立在
"借的是裸 Esc"上 —— 它拿 `escBorrowProbe`（只探 `modNoRepeat`+`vkEscape` 那一枚组合，见
`:57-81` 函数体）要求"桌面上裸 Esc 此刻**不**空闲"。种子是 `"Ctrl+Alt+V"` 时，形ⓐ 借走的是 `'V'`，
裸 Esc 反而**空闲** ⇒ `free == true` ⇒ 走到 `:113` 那句 `t.Fatal`。
⇒ 这枚 helper 里对旧行为的背书是**两格**（VK 字面量＋桌面探测），不是一格；只改 107 会留一枚
"改了形但还在探错键"的尺。

---

## 2. 名册：同包内所有"背书借的是裸 Esc"的尺（逐枚读函数体，不是 grep 词面命中）

词面尺：`grep -n 'vkEscape\|requireEscBorrowed\|escBorrowProbe' internal/ball/*.go`
（另加 `escBorrowAcc\|TakeEscForCancel\|ReleaseEscAfterSession` 三枚，因为借用的"是谁"写在这几个符号上）。
下表每一行的"它到底断什么"一列来自函数体，不是来自名字／注释。

| # | 位置 | 符号 | 它到底断什么 | 形ⓐ 之后 | 本枚处置 |
|---|---|---|---|---|---|
| R1 | `hotkey_live_test.go:90-116`（VK 格 `:107-108`） | `requireEscBorrowed` | 借用活着时：live 集 4 枚；**且** cancel 槽那枚 `VK == vkEscape` | ⚠ 变旧（种子 `'V'` ⇒ 该槽持 `0x56`） | **同步**：期望改成"等于按种子解出的那枚组合"，默认种子那一支仍逐字要求 `{0x4000,0x1B}` |
| R2 | 同 helper `:110-115` | `requireEscBorrowed` 第二格 | 拿 `escBorrowProbe` 问 **Win32**：裸 Esc 此刻不许空闲（"真被我们借走了"） | ⚠ 变旧（非默认种子下 Esc 本就该空闲） | **同步**：探"借的那枚"（按种子解出的组合）必须不空闲；非默认种子再追一格"裸 Esc 必须空闲"＝形ⓐ 的承诺 |
| R3 | 同 helper `:96-97` | standby 分支 `t.Fatalf` | "live 集只有 3 枚且 cancel 行仍 standby" ⇒ 从没试过借用（措辞写死 "the Esc borrow"） | ✅ 断的是"试没试"，不是"哪枚键"；仅**措辞**过期 | 措辞同步（`cancel borrow`），断言语义不动 |
| R4 | 同 helper `:100-102` | `HotkeyTaken/Error` 分支 `t.Logf` | 试过但被 Win32 拒 ⇒ SKIP-LOUD 一句"别人占了**裸 Esc**" | ⚠ 措辞断的是错的键（被拒的是配置那枚） | 措辞同步成"被探的那枚组合"，仍 `return`（不是新增 skip，原样保留 SKIP-LOUD 语义） |
| R5 | `hotkey_live_test.go:45-56`＋`:57-81` | `escBorrowProbe` | 在 foreign id 上真注册 `{modNoRepeat, vkEscape}` 再立刻卸载，问桌面裸 Esc 空不空闲 | ✅ 它是"探裸 Esc"这件事本身，不是背书借用 | 保留签名（`:100` 与 `hotkey_cancel_borrow_live_260_test.go:100` 两处调用不许改），把"探任意一枚组合"抽成新 helper 供 R2 用 |
| R6 | `hotkey_live_test.go:83-89` | `requireEscBorrowed` 的 doc 注释 | 散文背书："the bare Esc is in OUR registration set" | ⚠ 散文过期 | 随同步改写 |
| R7 | `hotkey_live_test.go:423-452` | `requireIdleRoster` | 稳态：live 恰 3 枚、不含 cancel、三枚点名，**且**裸 Esc 在桌面上空闲 | ✅ 稳态从不绑 cancel（票 245 未动），断言仍真；但它**不再**能证明"借的那枚键已归还"（旧行为里借的就是它探的那枚） | 本体不动；归还那一格见 R8 |
| R8 | `hotkey_live_test.go:454-462` | `requireEscReturned` | "还"＝`!EscTakenOver()` 且转手叫 R7 | ⚠ 隐性变弱：非默认种子下它探的是**没被借过**的那枚键，归还失败看不见 | **同步**：追一格"借过的那枚组合此刻必须空闲"（默认种子＝裸 Esc ⇒ 与 R7 同一探测，逐位等价；非默认种子＝配置那枚 ⇒ 真看得见没还） |
| R9 | `hotkey_status_test.go:202-262`（VK 字面量在 `:223`）／`hotkey_cancel_borrow_260_test.go:57-108` | `TestCancelBorrowRoundTrip`／`TestBorrowDefaultIsBitIdentical260` | 走 `takeEscWith`（＝`resolveCancelBorrow("")`，`hotkey_windows.go:619-625` 注释自述"the DEFAULT-档 borrow"）借出的必是 `{vkEscape, modNoRepeat}`，且 `:60` 先把 `modNoRepeat==0x4000`／`vkEscape==0x1B` 两枚常量钉住 | ✅ **不是冲突**：这两枚读的是"默认档"那条 wrapper，260-r1 特意留下 wrapper（`hotkey_windows.go:621-622` 逐字 "this wrapper stays so the pre-260 assertions keep checking the default case instead of being edited into a new shape"）正是为了让旧断言继续量默认那一支 | ⛔ 一字不动（写面内也不动） |
| R10 | `hotkey_cancel_borrow_260_test.go:324-349`（func 在 `:332`，`t.Logf` 在 `:347`） | `TestLiveBorrowHelperPremiseIsConfigDependent260` | 260-r1 上报本冲突的那枚：默认档可跑，断种子 `"Ctrl+Alt+V"` 借出 `{mods 0x4003, VK 0x56}` 并 `t.Logf` 记下 helper premise 变成 config-dependent | ✅ 它是**上报尺**，不是背书尺 | ⛔ 不动。同步完成后它的记录仍成立（说的正是 R1/R2） |
| R11 | `interaction_live_test.go:232-239` | `injectBinding(t, "Esc")` ⇒ `t.Fatalf("the taken-over Esc did not cancel Confirming (machine %s)")` | 真窗注入**裸 Esc**，要求 Confirming 被取消 | ⚠ **同类过期**：那枚 ball 的种子是 `liveHotkeys()`（`:148`），形ⓐ 下借走的是 `Ctrl+Alt+V`，注入裸 Esc 不该取消 ⇒ 该 Fatalf 会红 | **⛔ 写面外**（本枚只许写 `hotkey_live_test.go` ＋新文件）。已具名进 §6/§7，归 260-v1／真窗窗口，⛔ 我自己不动它 |
| R12 | `interaction_live_test.go:141`／`:204`／`:217`／`:224`／`:247-248`／`:272`；`live_windows_test.go:71-73` | 注释散文 | "Esc is the cancel key while Confirming" 一类 | ⚠ 散文过期（同 `b1e59d63` 处理日志那句的口径） | ⛔ 写面外，具名 §7 |

调用点清点（`requireEscBorrowed` 共 4 处，编排者说的"四枚 winlive 调用点"数量对上）：
`hotkey_live_test.go:311`（种子 `liveHotkeys()` ⇒ 受 R1/R2 影响）、
`live_windows_test.go:127`（`:59` 种子 `liveHotkeys()` ⇒ 受影响）、
`interaction_live_test.go:205` 与 `:257`（`:148` 种子 `liveHotkeys()` ⇒ 受影响）、
`hotkey_cancel_borrow_live_260_test.go:50`（种子 `DefaultHotkeys()` ⇒ **不受影响，必须继续绿**）。

---

## 3. 同步方案（裁定已下，不重设计）

`requireEscBorrowed(t, b)` 的签名**不改**（改签名要动 R11/R12 那几个写面外文件）。
改成：**期望值从 ball 自己的种子解出来**，解的过程用一只纯函数
`wantCancelBorrow260r2(bind string)`，放在新增的**默认档**文件
`internal/ball/hotkey_cancel_borrow_expect_260r2_test.go`（`//go:build windows`）里 ⇒
今天就能跑、就能被正控打红，而不是只在 `-tags winlive` 那档里等着撞。

两格分开写，⛔ 不合成"热键接好了"一句：

- **格①（默认档零漂移，形ⓐ 条件①）**：种子解出 `{0x4000, 0x1B}` 时——
  断言逐字要求"等于 `escBorrowAcc()`＝`{Mods 0x4000, VK 0x1B}`"，桌面探测仍探裸 Esc 且要求不空闲。
  这一支的期望值**写死位模式**（不写 `vkEscape` 常量），与 260-r1 常驻尺
  `hotkey_cancel_borrow_260_test.go:57` 同向。
- **格②（非默认种子跟着配置）**：种子是别的组合时——
  断言"live 的 cancel 槽 == 按种子解出的那枚"，桌面探测换成探**那枚组合**且要求不空闲；
  再追一格"裸 Esc 此刻必须空闲"（形ⓐ 的意义＝搬走 Esc 的用户不再连带丢掉 Esc）。
- **正控**（本仓死规矩：负向尺必配"种 X 必响"）：默认档那枚新用例种 `"Esc"` 与 `"Ctrl+Alt+Y"` 两枚
  种子，断言两格期望**各自**等于自己的位字面量、且彼此不等 ⇒ 期望值真的跟着种子在动；
  再种一枚垃圾值，断期望退回默认＋带上"被拒"那句话（P6 同向）。真窗那一支的正控按第 5 条处理，
  具名〔未实测〕。

⛔ 不删断言、⛔ 不 `t.Fatalf`→`t.Log`、⛔ 不加 `t.Skip`（`tools/d22scan/runtests.sh:98` 把 SKIP 判红）。

---

## 4. 改前 / 改后逐字

（本节在实现 commit 里展开；§1 的改前逐字已在本枚落盘。）

---

## 5. 门禁读数

起手基线（本枚实测，2026-10-04 10:15 +08，树上 HEAD `5ad8ec27`，改前）：

- `go test -count=1 -v ./internal/ball/` ⇒ `--- PASS 62`／`--- FAIL 1`／`--- SKIP 0`，包行 `FAIL github.com/CarlosShao/wisp/internal/ball 0.200s`（exit 1）。
  唯一红名＝`TestC21TableColourRowsMatchTokensCSS`（`tokens_table_test.go:1468`
  `read design/assets/tokens.css: ... cannot find the path specified`）＝编排者点名的**基线常红**，
  先于本腿存在，⛔ 不算本腿回归、⛔ 不修。原件 `gate-test-ball-baseline-verbose.txt`。
- `go vet ./internal/ball/` ⇒ exit 0，零字节输出（`gate-vet-ball-baseline.txt`）。
- `go vet -tags winlive ./internal/ball/` ⇒ exit 0，零字节输出（`gate-vet-ball-winelive-baseline.txt`）
  ＝ winlive 档**编译**在改前是干净的；⚠ 这只证明编得过，不证明断言跑得对（那档今天⛔ 不许执行）。

---

## 6. 欠账（具名，⛔ 不算已量）

- 真窗／winlive 全档今天**没跑**（`//go:build winlive` 会在机主桌面真注册全局热键、真开球窗；
  本轮没有这个词）。所以：R1/R2/R8 的同步后行为、以及"改前那四枚调用点会红"这条，
  **都只是读码＋默认档解值推出来的预测**，⛔ 不许写成"已红"。
- 归 260-v1／真窗窗口：`go test -tags winlive ./internal/ball/ -run "TestLiveHotkey|TestBallLive|TestLiveConfirming|260" -v`
  跑一次，取 R1/R2/R8 同步后的逐枚读数（默认种子那一支必须与形ⓐ 之前逐位同判）。
- R11（`interaction_live_test.go:232-238` 注入裸 Esc 那格）落在本枚**写面外**，未同步、未跑。

---

## 7. 判不动的地方

- **J-1**：R11/R12 的同步需要写 `internal/ball/interaction_live_test.go`／`live_windows_test.go`，
  本枚写面只有 `hotkey_live_test.go` ＋带 `260r2` 的新文件 ⇒ 交回编排者派下一枚。
- **J-2**：`requireEscBorrowed` 这个名字本身（含 `escBorrowProbe`）在形ⓐ 之后是过期词面，
  改名要动 4 个调用文件（3 个在写面外）⇒ 不改名，只在注释里把它说清；这条与 `b1e59d63`
  那句"idle 日志不许再说 Esc is borrowed"同性质，属**词面**账，不属断言账。
- **J-3**：改前那四枚调用点到底几枚红、红句长什么样——今天量不到（不许跑真窗）。
  写"会红"就是预测，不是读数。
