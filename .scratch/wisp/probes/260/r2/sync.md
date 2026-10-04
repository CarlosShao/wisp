# 票 260 · 260-r2 证据件 —— 把一枚因编排者裁定而过期的既有真窗判据，同步到已裁的行为上

本轮唯一一格：`internal/ball/hotkey_live_test.go`（`//go:build windows && winlive`）里那枚
`requireEscBorrowed` 的"借到的必是裸 Esc"断言，在形ⓐ（确认窗借用吃 `[hotkey] cancel` 配置值）
落地之后变成对旧行为的背书。裁定＝**期望值跟着种子走**，不跟着写死常量走；
默认档零漂移那一格不许丢；⛔ 不许删断言 / 不许 `t.Log` 化 / 不许 `t.Skip`。

本枚 = 骨架 commit（先落"我读到的冲突原文逐字"＋"名册"＋门禁基线）。
§4 之后的「改后逐字／终值读数」在本轮后半程的 commit 里展开；写这一句是为了让"哪些节已满、
哪些节还没满"是**具名状态**而不是一句空格。

件的两枚 commit：骨架 `84dbda52`（读码＋名册＋基线读数）／尺同步 `e3e19e8e`
（`internal/ball/hotkey_live_test.go` ＋新档 `internal/ball/hotkey_cancel_borrow_expect_260r2_test.go`）。
本件满稿在第三枚（只带证据件）。锚点：起手 HEAD `5ad8ec27`（`dev`）。
260-r1 四枚 = `eb2c0173`／`d9bb6817`／`b1e59d63`／`ba2a8358`。

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
  期望＝**写死的位对** `wantDefaultBorrow260r2`（新档 `:38`），红句逐字点名 `{mods 0x4000, VK 0x1B}`；
  ⛔ 不拿 `escBorrowAcc()`／`vkEscape` 常量当期望（那样"默认档漂了"和"常量漂了"就分不开），
  它与产码默认支的一致性由新档 `TestBorrowWishDefaultCellIsBitIdentical260r2` 钉住；
  桌面探测仍探裸 Esc 且要求不空闲。与 260-r1 常驻尺
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

### 4.1 R1（VK 字面量格）

改前（`hotkey_live_test.go` 旧 `:107-108`，逐字）：

```go
	if got := rep.Live()[hkCancel]; got.VK != vkEscape {
		t.Fatalf("the borrowed cancel slot holds VK 0x%X, want VK_ESCAPE 0x1B: %+v", got.VK, rep.Bindings())
	}
```

改后（现 `internal/ball/hotkey_live_test.go:145-159`，逐字）：

```go
	if got := rep.Live()[hkCancel]; got != wish.Acc {
		switch {
		case wish.Refused != nil:
			t.Fatalf("the borrowed cancel slot holds %+v, want the default pair %+v: [hotkey] cancel %q is not a "+
				"parsable combination (%v), and condition ③ says the borrow takes the default and SAYS so: %+v",
				got, wish.Acc, seed, wish.Refused, rep.Bindings())
		case wish.Default:
			t.Fatalf("ticket 260 condition ① RED: [hotkey] cancel is the default (%q), so the borrowed slot must "+
				"hold the raw pair {mods 0x4000, VK 0x1B}, got %+v: %+v", seed, got, rep.Bindings())
		default:
			t.Fatalf("the borrowed cancel slot holds %+v, want the configured cancel key %q = %+v: after ticket 260 "+
				"形ⓐ the borrow registers what the config names, so this is either the pre-260 hard-coded Esc or a "+
				"key the config never named: %+v", got, seed, wish.Acc, rep.Bindings())
		}
	}
```

期望值来源（同一 helper 头部，现 `internal/ball/hotkey_live_test.go:127-128`）：

```go
	seed := b.ConfiguredHotkeys().Cancel
	wish := wantCancelBorrow260r2(seed)
```

`wantCancelBorrow260r2` 的本体在新档 `internal/ball/hotkey_cancel_borrow_expect_260r2_test.go:57-66`
（默认档可跑；三支：`""` ⇒ 裸 Esc 位对、解析不了 ⇒ 位对＋带 `Refused`、其余 ⇒ `ParseAccelerator(种子)`）。

### 4.2 R2（桌面探测格）

改前（旧 `:110-115`）：探 `{modNoRepeat, vkEscape}`，`free == true` ⇒ `t.Fatal("... the borrow registered on some other key, not on VK_ESCAPE")`。

改后（现 `internal/ball/hotkey_live_test.go:160-176`，逐字）：

```go
	// The desktop, not our own bookkeeping, says the borrowed combination is ours.
	if free, err := hotkeyFreeProbe(t, b, wish.Acc); err != nil {
		t.Fatalf("the probe for the borrowed %q (%+v) could not run: %v", seed, wish.Acc, err)
	} else if free {
		t.Fatalf("ticket 245 RED: the ball says it borrowed %q (%+v), but the desktop still has that combination "+
			"free - the borrow registered on some other key: %+v", seed, wish.Acc, rep.Bindings())
	}
	if !wish.Default {
		// The second cell of 形ⓐ, on the same live ball: the key nobody configured
		// must not be the one that got taken.
		if free, err := escBorrowProbe(t, b); err != nil {
			t.Fatalf("the bare-Esc probe could not run: %v", err)
		} else if !free {
			t.Errorf("bare Esc is claimed while [hotkey] cancel says %q and the borrow holds %+v - the ball is "+
				"holding two keys: %+v", seed, wish.Acc, rep.Bindings())
		}
	}
```

`escBorrowProbe` 的签名与语义一字未动（现 `:93-99`，body 抽给新 `hotkeyFreeProbe` `:67-91`），
所以 `requireIdleRoster`（R7）和 260-r1 的真窗件 `hotkey_cancel_borrow_live_260_test.go:100` 照旧能编、照旧探裸 Esc。

### 4.3 R3/R4（措辞，不改语义）

- standby 那枚红句：`never attempted the Esc borrow` ⇒ `never attempted the cancel borrow`，并把种子打进句子
  （`[hotkey] cancel = %q`）；断言条件（live 集 3 枚且 cancel 行仍 standby）逐字保留。
- SKIP-LOUD 那枚：`another program on this desktop owns bare Esc` ⇒ `owns %+v`（＝被尝试的那枚组合）；
  `t.Logf`＋`return` 的形状是票 245 原有的，⛔ 本枚没新增任何 skip，也没把 Fatal 降成 Log。

### 4.4 R8（"还"那一格补回被形ⓐ 拆掉的牙）

改前（旧 `:454-462`）：`!EscTakenOver()` ⇒ `requireIdleRoster`（＝探裸 Esc 空闲）。
种子是 `Ctrl+Alt+V` 时，"没还得掉的那枚键"和"探的那枚键"不是同一枚 ⇒ 归还失败看不见。
改后（现 `:517-555`）：默认种子 ⇒ 早退（与改前逐位同判，⛔ 零新增探测）；非默认种子 ⇒ 追加
`hotkeyFreeProbe(t, b, wish.Acc)`，仍空闲才算还；若桌面仍占用，再分两种读数——
自己的注册集里还有 `hkCancel` ⇒ `t.Fatalf`（没还）；注册集干净但桌面仍占用 ⇒ `t.Errorf`
（卸载没落到 Win32，或跑到一半被别家抢走，"这一发分不清、两种都不是通过"）。⛔ 没有一支是 `t.Log`。

### 4.5 为什么这是"同步到已裁行为"，不是放宽断言（验收腿要攻的就是这段）

1. **判据只多不少**（逐段 `sed` 取 helper 本体数出来的，不是估的）：
   `requireEscBorrowed` 改前＝5 支 fatal-class（4 `t.Fatalf`＋1 `t.Fatal`）＋1 支 SKIP-LOUD `t.Logf`；
   改后＝9 支 fatal-class（8 `t.Fatalf`＋1 `t.Errorf`）＋同 1 支 SKIP-LOUD。
   `requireEscReturned` 改前＝1 支 `t.Fatal`（其余转手 `requireIdleRoster`）；
   改后＝3 支 fatal-class（1 `t.Fatal`＋2 `t.Fatalf`）＋1 支 `t.Errorf`。
   ⛔ 一枚没删、⛔ 一支 Fatal 都没降级、⛔ 没新增 `t.Skip`（同文件 `t.Skip` 计数改前改后都是 `4`，
   全是票 64/245 原有的注入两支＋占用那一支＋另一枚，本枚一枚未动，见 §5.3 读数）。
2. **期望值变严了，不是变松了**：旧尺只查 `got.VK != vkEscape`（**连 mods 都不比**）；新尺比**整枚
   `Accelerator{Mods, VK}`**，且这枚值来自配置种子。旧尺会把"用户配了 Ctrl+Alt+V、实际还借裸 Esc"
   判**绿**（那正是票 260 立案的缺陷形状）；新尺判**红**（`mutation-N2` 就是这个形状，§5.4 红句逐字）。
3. **红→绿的唯一一支**，就是"借的键＝配置里那枚键"——那是编排者 A588 已裁的行为，不是新放宽的口子。
4. **默认档零漂移**：`wish.Default` 为真时期望是**写死的位对** `{0x4000, 0x1B}`（不是 `vkEscape` 常量），
   桌面探测仍是那一发裸 Esc 探测，判定链与改前逐字同形；常驻尺
   `TestCancelBorrowRoundTrip`／`TestDefaultHotkeysIdlePassHoldsNoEsc`／`TestBorrowDefaultIsBitIdentical260`
   改后仍 PASS（§5），`TestLiveCancelBorrowDefaultUnchanged260` 的判据没被搬动（真窗档未跑，具名 §6）。
5. **不是同义反复**：期望走 `ParseAccelerator(种子)`，实际走 `resolveCancelBorrow(种子)` 落进报告的
   `Live()`——两条独立路径；`mutation-N1`（期望退回常量）与 `mutation-N2`（产码退回写死）
   都能各自把尺打红（§5.4 红句逐字），说明两侧任一偷懒都会响。

---

## 5. 门禁读数

### 5.1 起手基线（本枚实测，2026-10-04 10:15-10:16 +08，树上 HEAD `5ad8ec27`，改前）

- `go test -count=1 -v ./internal/ball/` ⇒ `--- PASS 62`／`--- FAIL 1`／`--- SKIP 0`，包行 `FAIL github.com/CarlosShao/wisp/internal/ball 0.200s`（exit 1）。
  唯一红名＝`TestC21TableColourRowsMatchTokensCSS`（`tokens_table_test.go:1468`
  `read design/assets/tokens.css: ... cannot find the path specified`）＝编排者点名的**基线常红**，
  先于本腿存在，⛔ 不算本腿回归、⛔ 不修。原件 `gate-test-ball-baseline-verbose.txt`。
- `go vet ./internal/ball/` ⇒ exit 0，零字节输出（`gate-vet-ball-baseline.txt`）。
- `go vet -tags winlive ./internal/ball/` ⇒ exit 0，零字节输出（`gate-vet-ball-winelive-baseline.txt`）
  ＝ winlive 档**编译**在改前是干净的；⚠ 这只证明编得过，不证明断言跑得对（那档今天⛔ 不许执行）。

### 5.2 改后终值（本枚实测；逐份原件落盘时刻＝vet 10:25／真窗编译 10:25／d22scan 10:28-10:30 +08，整包在 10:30 首取、10:38 进树后复跑同值）

- `go test -count=1 -v ./internal/ball/` ⇒ `--- PASS 65`／`--- FAIL 1`／`--- SKIP 0`，包行
  `FAIL github.com/CarlosShao/wisp/internal/ball 0.182s`（exit 1）。
  **红名集合与基线同名**＝只有 `TestC21TableColourRowsMatchTokensCSS`（还是 `tokens.css` 不在盘上那一因，
  与本票无关、⛔ 本枚不修）。`65 = 62 基线 ＋ 本枚新增 3 枚`，三枚逐名 PASS：
  `TestBorrowWishFollowsSeed260r2`／`TestBorrowWishDefaultCellIsBitIdentical260r2`／
  `TestBorrowWishRefusedSeedIsNotSilent260r2`（`gate-test-ball-final-verbose.txt`）。
  同档里"默认档零漂移"的既有常驻尺逐枚仍 PASS：`TestCancelBorrowRoundTrip`／
  `TestDefaultHotkeysIdlePassHoldsNoEsc`／`TestRegisterAllLiveSet`／`TestBorrowDefaultIsBitIdentical260`／
  `TestBorrowFollowsConfiguredBinding260`／`TestBorrowUnparsableFallsBackLoudly260`／
  `TestBorrowReceiptSharesOneSourceWithRegistration260`／`TestBorrowFailureNamesTheKeyItTried260`／
  `TestConfiguredCancelStillNeverBoundWhileIdle260`／`TestLiveBorrowHelperPremiseIsConfigDependent260`。
- `go vet ./internal/ball/` ⇒ exit 0，零字节（`gate-vet-ball-final.txt`）。
- `go vet -tags winlive ./internal/ball/` ⇒ exit 0，零字节（`gate-vet-ball-winelive-final.txt`）
  ＝同步后的真窗尺**编译得过**；⚠ 这只算编译读数，⛔ 本枚一次都没执行真窗档。
- `gofmt -l internal/ball/` ⇒ 空；`$(go env GOPATH)/bin/gofumpt.exe -l internal/ball/` ⇒ 空。
- 进树后复跑 `go test -count=1 -v ./internal/ball/`（10:38）⇒ 读数与首取**逐字同值**（65／1／0，红名同名）。
- 同文件 `t.Skip` 计数：改前 `4`／改后 `4`（⛔ 本枚没新增任何 skip；SKIP-LOUD 那支 `t.Logf` 是票 245 原有的形状）。

### 5.3 D22 门（`sh scripts/d22scan.sh`）——撞上一枚别家在飞的门，如实记

- 本枚跑 ⇒ **exit 1**，且**停在第一步正控**（脚本是 `set -eu`，第二步真扫描因此根本没跑）。
  唯一 FAIL＝`TestSelfTestFlagIsWiredInTheBuiltBinary`：`d22scan -self-test: 4 direction(s) failed, 36/40 passed`，
  四支全在 **ban #9 phantom-citation** 的 silent 方向（原件 `gate-d22scan-final.txt`）。
- 归因（起手就取过数）：`git status --short tools/d22scan/` ＝ ` M main.go`／` M scan_test.go`／` M selftestsamples.go`
  ＝别家在飞的写面，本枚⛔ 未碰该目录一字（§5.4 名册可复算）；票面禁区也写死 `tools/d22scan/allowlist.txt` 不许动。
  ⇒ 这一红**不是本枚造成的**，也⛔ 不由本枚去"修绿"。
- 同一次跑里读到一枚有用的绿：`TestScannerSelfScanOfRealRepoIsGreen` `--- PASS (1.08s)`
  ＝那把（在飞的）扫描器真走过本仓工作树并报 clean，覆盖本枚改动的两枚文件。
- 补一发**HEAD 版**扫描器以拿到可引用的独立读数：用 `git cat-file blob HEAD:tools/d22scan/<f>` 取
  `main.go/gitignore.go/selftest.go/selftestsamples.go/allowlist.txt/go.mod` 到 `$TEMP`（⛔ 不在仓内建 checkout／worktree），
  `go build` 后 `-root "D:/work/workspace/projects plans/Wisp"` ⇒ **exit 0**，
  `clean - no D22 ban violations`；`ban #8 internal/ 501 Go files, comments and _test.go included`
  （260-r1 终值 500 ⇒ ＋1 就是本枚新档），`ban #8 cmd/ 98`（比 r1 的 97 多 1＝别家在飞的 `cmd/wisp` 文件，与本枚无关）。
  原件 `gate-d22scan-headbinary-scan.txt`。

### 5.4 写面名册与 AC#3 越界面（只算本枚的 commit）

- `internal/ball/hotkey_live_test.go` ⇒ 改前对 HEAD 的 `git diff --numstat` ＝ `125  32`（＋125／−32），
  原件 `diff-hotkey-live-test.txt`；已随 `e3e19e8e` 进树。
- `internal/ball/hotkey_cancel_borrow_expect_260r2_test.go` ⇒ **新档**（`//go:build windows`，默认档跑）。
- `.scratch/wisp/probes/260/r2/**` ⇒ 证据件（含两份突变驱动、四份突变读数、三份门禁读数、diff 原件）。
- ⛔ 未碰产码：`internal/ball/hotkey_windows.go` md5 起手＝终值＝`be79cbda486a334d74fe9d278558acba`
  且 `git diff --numstat -- internal/ball/hotkey_windows.go` 为空；`internal/ball/ball_windows.go` 零改动。
- ⛔ 未碰：`internal/agent/approval/**`（212-r3 在飞）、`cmd/wisp/**`（256-r1 在飞）、
  三枚冻结件、`.scratch/wisp/issues/**` 票面复选框（一枚没动）、台账。
- AC#3 名单命中数（对本枚两枚 commit `84dbda52`／`e3e19e8e` 逐枚 `git show --name-only`）＝ **0**；
  票面／台账／三枚冻结件／两个在飞包（`internal/agent/approval`、`cmd/wisp`）命中数也＝ **0**。
  命令与原始输出＝`gate-ac3-scope.txt`（里面还具名写了"起手..HEAD 整段区间含别家 commit，
  同一 git user 用 `--author` 分不开"这一口径坑）。
- 起手 `git status` 里那批 ` M`／` D`（`design/**` 16 枚删除、`tools/d22scan/**`、`cmd/wisp/**`、
  `internal/agent/approval/pending_read.go`、`.gitignore`、别家 probes）**一律不是本枚写的**；
  本枚 add 只带显式 pathspec，共树别人的增量留在原地不动。

### 5.5 突变＝本枚尺的正控（红句逐字）

驱动与还原：`md5-before.txt`／`md5-restored-final.txt` 两枚文件的 md5 逐字相同
（产码 `be79cbda…`、本枚新档 `cb3997a5…`），产码 `git diff --numstat` 空 ⇒ 突变只进盘一次、原样写回。
⚠ **第一版驱动 `mutate260r2.py` 的 M1/M2 被它们自己的读数否掉**：两支都只印 `package exit=1`、
零枚 `--- FAIL`＝改动让局部变量没人用，`go test` 死在**编译**，那不叫红。第二版 `mutate260r2b.py`
保持所有变量仍被使用、只改"值"，四支全打到断言。旧驱动与旧读数按"临时件只建不删"留在
`mutation-M1M2M3.txt`。四支都跑在**默认档**；真窗那几枚的"种 X 必响"取不到读数，欠账具名 §6。

- **N1＝期望退回常量**（`wantCancelBorrow260r2` 无视种子只报默认位对）⇒ `TestBorrowWishFollowsSeed260r2` 红：
  ```
  hotkey_cancel_borrow_expect_260r2_test.go:89: the winlive fixture seed at hotkey_live_test.go:42: seed "Ctrl+Alt+V" expects the borrowed slot to read {Mods:16387 VK:86}, got {Mods:16384 VK:27}
  hotkey_cancel_borrow_expect_260r2_test.go:117: the expectation does not follow the seed: "Esc" and "Ctrl+Alt+Y" resolve to the same accelerator
  ```
- **N2＝产码退回写死**（`resolveCancelBorrow` 的配置支强制 `VK: vkEscape`）⇒ 260-r1 三支＋本枚一支全红：
  ```
  hotkey_cancel_borrow_260_test.go:127: modifiers plus a letter: cancel = "Ctrl+Alt+K" resolves to {Mods:16387 VK:27}, want {Mods:16387 VK:75}
  hotkey_cancel_borrow_260_test.go:244: cancel = "Ctrl+Alt+K": the receipt prints "Ctrl+Alt+K" = {Mods:16387 VK:75} but registers {Mods:16387 VK:27} - the page and the key disagree
  hotkey_cancel_borrow_260_test.go:338: the borrow of "Ctrl+Alt+V" resolved to VK_ESCAPE again - condition ①'s default-only rule was over-applied
  hotkey_cancel_borrow_expect_260r2_test.go:106: the winlive fixture seed at hotkey_live_test.go:42: RegisterHotKey was handed {Mods:16387 VK:27} while the ruler expects {Mods:16387 VK:86}
  ```
  ⇒ 这一支就是"旧尺会判绿的那个缺陷形状"，新尺（期望⇄注册两条独立路径）判红。
- **N3＝解析失败被静默吞掉**（去掉 `Refused`）⇒ `TestBorrowWishRefusedSeedIsNotSilent260r2` 红：
  ```
  hotkey_cancel_borrow_expect_260r2_test.go:162: seed "Ctrl+Alt+NotAKey+" parsed cleanly; the ruler would then expect the user's key while production borrows {Mods:16384 VK:27}
  ```
- **N4＝默认位对被搬动**（`wantDefaultBorrow260r2` 的 VK 改 0x1C）⇒ 三支红，其中点名条件①那一支：
  ```
  hotkey_cancel_borrow_expect_260r2_test.go:127: the raw default pair moved to {Mods:16384 VK:28}; condition ①'s ruler is measuring a value nobody registers
  hotkey_cancel_borrow_expect_260r2_test.go:89: unset in config: seed "" expects the borrowed slot to read {Mods:16384 VK:27}, got {Mods:16384 VK:28}
  hotkey_cancel_borrow_expect_260r2_test.go:93: the shipped default spelling: seed "Esc" Default = false, want true
  ```
  原件：`mutation-N1-run.txt`／`mutation-N2-run.txt`／`mutation-N3-run.txt`／`mutation-N4-run.txt`／
  `mutation-N1N2N3N4-summary.txt`。

---

## 6. 欠账（具名，⛔ 不算已量）

1. **真窗／winlive 全档今天一次没跑**（`//go:build winlive` 会在机主桌面真注册全局热键、真开球窗；
   本轮没有这个词，任务书也写死⛔ 不许跑）。因此下面这些都还是**读码＋默认档解值推出来的预测**：
   - 「改前那四枚调用点里，三枚（种子 `liveHotkeys()`＝`Ctrl+Alt+V`）会撞 R1/R2」＝**预测，不是读数**；
     第四枚 `hotkey_cancel_borrow_live_260_test.go:50`（种子 `DefaultHotkeys()`）按同步后逻辑应当
     与形ⓐ 之前逐位同判，这条同样**没跑过**。
   - 同步后的尺本身（R1/R2/R8 的新分支）在真窗上有没有牙、红句长什么样＝**未量**。
     只有它的**期望算式**在默认档有读数（§5.2／§5.5）。
   - 取数口令（归 260-v1／真窗窗口，⛔ 别在编队里并发跑，两把真窗尺会互相抢热键）：
     `go test -tags winlive ./internal/ball/ -run "TestLiveHotkey|TestLiveCancelBorrow|TestBallLive|TestLiveConfirming" -v`
     要盯的逐枚判据：`hotkey_live_test.go:311`（`TestLiveHotkeyOccupiedVsNotAttempted` 的借还段）、
     `live_windows_test.go:126-129`、`interaction_live_test.go:205`／`:257`、
     `hotkey_cancel_borrow_live_260_test.go:50`（默认档那一支，必须与改前同判）。
2. **R11 没同步**：`interaction_live_test.go:232-239` 注入裸 Esc 并要求取消 Confirming ——
   那枚 ball 的种子是 `liveHotkeys()`，形ⓐ 下借走的是 `Ctrl+Alt+V`，这一支按读码会红。
   它在**本枚写面外**（任务书写死只许 `hotkey_live_test.go` ＋带 `260r2` 的新档）⇒ 未动、未跑，交回编排者。
3. **R12 词面账没清**：`interaction_live_test.go:141`／`:204`／`:217`／`:224`／`:247-248`／`:272`
   与 `live_windows_test.go:71-73` 那六句"Esc 就是此刻的取消键"一类散文（同 `b1e59d63` 处理 idle 日志那句的口径），
   同样在写面外，未动。另有 `A595` 那批"十枚人看文案"（编排者已具名解冻、派的是别的腿）⛔ 不是本枚的格。
4. **helper 名字本身**：`requireEscBorrowed`／`escBorrowProbe` 在形ⓐ 之后是过期词面（J-2），
   本枚只把注释与红句改到实话，⛔ 没改名（改名要动 4 个调用文件，其中 3 个在写面外）。
5. **桌面探测的固有歧义没被消掉**：`hotkeyFreeProbe` 分不清"这枚组合是我们占的"与"别家占的"——
   这是票 245 那把尺天生就有的性质，本枚没假装解决；R8 的处置是把两种读数**都判不通过**
   （自己注册集仍留着＝`t.Fatalf`；注册集干净而桌面仍占＝`t.Errorf`），⛔ 没有一支降级成 `t.Log`。

---

## 7. 判不动的地方

- **J-1（写面外）**：R11/R12 的同步要写 `internal/ball/interaction_live_test.go`／`live_windows_test.go`，
  本枚写面只有 `hotkey_live_test.go` ＋新档 ⇒ 交回编排者派下一枚；⛔ 本枚没碰。
- **J-2（改名代价）**：见 §6.4。符号改名＝跨 4 枚文件的产码级变更（用例名进 CI 口径），不是尺能自己决定的。
- **J-3（预测≠读数）**：改前那几枚真窗调用点到底几枚红、红句长什么样，今天量不到（不许跑真窗）。
  §5.2 的 `go vet -tags winlive`＝0 只证明**编得过**，⛔ 不证明断言跑得对。
- **J-4（门被别人卡住）**：`sh scripts/d22scan.sh` 的**第一步正控**今天红（ban #9 self-test 4/40 支失败），
  起因是别家在飞的 `tools/d22scan/{main.go,scan_test.go,selftestsamples.go}` 三枚未提交增量；
  脚本是 `set -eu`，所以第二步真扫描**根本没跑**。本枚⛔ 不动 `tools/d22scan/**`（票面禁区＋别人的写面）
  ⇒ 只能报"停在正控"＋用 HEAD 版扫描器另取一发独立读数（§5.3）。**这道门今天没法给本枚出终值**，
  归下一枚能碰 `tools/d22scan` 的腿或编排者。
- **J-5（R8 是补牙不是裁定原文）**：A588 裁的是"期望值跟着种子走"，`requireEscReturned` 那条追加
  是形ⓐ 顺手拆掉的覆盖（归还失败看不见）——本枚按"同步到已裁行为"补回，默认种子那一支**零新增探测**、
  与改前逐位同判。如果编排者认为这一格该由别的腿做，回退范围＝`requireEscReturned` 一个函数体（§4.4）。
- **J-6（`ConfiguredHotkeys()` 当种子的来源）**：同步尺从 `b.ConfiguredHotkeys().Cancel` 取种子
  （＝产码 `b.cancelBinding` 的同一枚值，两处赋值点 `ball_windows.go:253-254`／`:819-821` 现读一致）。
  真窗上如果哪天这两处不再同值，尺会拿"配置说的"和"注册用的"两条腿互相指认——那正是**该红**的形状，
  但本枚没法在没有真窗的情况下证明它红得了。

