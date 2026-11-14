# 245 · AC#6② 兑现（裸键借／带修饰键常绑）— 实现腿 `245-ac6b`

> ⚠ **这不是裁决表**。按 `AGENTS.md §1.5`／`SPEC-12 §4.3`，`docs/evidence/s1/` 的判语必须出自**非实现者**；
> 本文件是实现腿的读数与台件，**票面 AC 框一枚未翻**。
>
> 派单：票 245 AC#6②，编排者 2026-11-14 裁定采纳（票面 `:195-198` 逐字：**裸键＝只在 Confirming 借；带修饰键＝常绑**）。
> 起手 HEAD：`c40e296c9f3e7262ed274140a70f764be8b67374`（2026-11-14 14:23:36 +0800）。
> 写面：`internal/ball/**`（产品码只动这一步判定及其判据）· `docs/evidence/s1/245-ac6b-*`（本件）·
> `.scratch/wisp/probes/245/ac6b/**` · 票 245 Progress log（只追加）。⛔ 未动 `cmd/wisp/**`、`internal/panel/**`、
> `internal/observe/**`、`internal/tools/**`、`scripts/**`、`frontend/**`、规格／阈值／golden。

## 0. 一句话结论

1. **判定在 `internal/ball/hotkey_windows.go:577`**：`if p.id == hkCancel && !cancelHeldWhileIdle(bindings[i])`。
   改前（`c40e296c9:hotkey_windows.go:484`）是 `if p.id == hkCancel` —— 无条件「永远待命」。
   现在：**裸键**（无 Ctrl/Alt/Shift/Win，判定函数 `:474 cancelHeldWhileIdle`）仍然`cancelIdleLine` 待命、从不注册；
   **带修饰键**的 cancel 走 `:526 registerHeldBinding` 正常注册，稳态 live。
2. **判据两条都过，且两侧正控都真红**（§2／§3）：配 `Ctrl+Alt+K` ⇒ 稳态注册集 4 枚含 cancel、`HotkeyReport().IsLive(hkCancel)`=true、
   `EscTakenOver()`=false；配 `Esc` ⇒ 稳态仍 3 枚、cancel standby、`HotkeyReport().IsLive(hkCancel)`=false、
   `EscTakenOver()`=false、借一次／还不重绑（今字未动）。
3. **归还路径必须跟上**：`ball_windows.go:1032` 由 `cancelIdleLine(b.cancelBinding)` 改为
   `cancelIdleLineWith(b.hwnd, b.cancelBinding)` —— 带修饰键的 cancel 在会话结束后要**重新绑回**（否则用户配的组合在
   第一张卡之后就被永久删掉）。裸键那一支仍是「drop、不重绑」，字节等价。

## 1. 判定现状（`file:line` 现读，不抄派单里的猜测位置）

| 事实 | 现读（`internal/ball/`） |
|---|---|
| 判定点（唯一） | `hotkey_windows.go:577` `if p.id == hkCancel && !cancelHeldWhileIdle(bindings[i])` |
| 「带修饰键才常绑」的规则 | `hotkey_windows.go:474-482 func cancelHeldWhileIdle`：解析成功且 `acc.Mods&realModifiers != 0`；空串／坏串返回 false（交给 `cancelIdleLine` 记 disabled／unparsable，与改前一致） |
| 修饰键掩码 | `hotkey_windows.go:133-140 const realModifiers = modAlt \| modControl \| modShift \| modWin`（`MOD_NOREPEAT` 不算修饰键） |
| 注册/分类本体 | `hotkey_windows.go:526-556 func registerHeldBinding`（从 `registerAllWith` 原样抽出，三个固定槽与「常绑的 cancel」共用一段，防漂移） |
| idle 取消行（裸键用） | `hotkey_windows.go:493 func cancelIdleLine`（doc 块 `:484-492`；standby 行、日志句未动） |
| 归还后的 idle 行 | `hotkey_windows.go:564-569 func cancelIdleLineWith`：裸键→`cancelIdleLine`；带修饰键→`registerHeldBinding(hotkeyRegisterer(hwnd))` |
| 归还路径调用 | `ball_windows.go:1032` `b.hotkeyReport = b.hotkeyReport.withCancel(cancelIdleLineWith(b.hwnd, b.cancelBinding))` |
| 借用路径 | **零改动**：`ball_windows.go:961 TakeEscForCancel` → `takeEscBorrow` → `takeEscWithAcc`（`escTakenOver`／`escBorrowRefused` 两旗未动） |

改前形状（可逐字重跑）：

```
$ git show c40e296c9:internal/ball/hotkey_windows.go | sed -n '483,492p'
	for i, p := range hkNames {
		if p.id == hkCancel {                      <- 无条件待命
			...
			rep.bindings = append(rep.bindings, cancelIdleLine(bindings[i]))
$ git show c40e296c9:internal/ball/ball_windows.go | grep -n 'cancelIdleLine(b.cancelBinding)'
1019: b.hotkeyReport = b.hotkeyReport.withCancel(cancelIdleLine(b.cancelBinding))
```

`summon` / `mute` / `panel` 的注册与分类**未动**（抽出函数逐字是原 `switch` 体）；`[hotkey]` schema、默认值
（`DefaultHotkeys()` / `ApplyHotkeyDefaults`）**未动**：`DefaultHotkeys().Cancel` 仍是 `"Esc"`。

## 2. 判据①（带修饰键常绑）＋ 正控

**尺**：`internal/ball/hotkey_cancel_borrow_260_test.go:311 TestConfiguredCancelIdleBindingFollowsModifierRuling245`（第 1 段）。

```
$ go test -count=1 ./internal/ball -run TestConfiguredCancelIdleBindingFollowsModifierRuling245 -v
--- PASS: TestConfiguredCancelIdleBindingFollowsModifierRuling245 (0.02s)
```

断言（全部现读通过）：

- `:317-320` `reg.attemptsOf(hkCancel) == 1` —— idle 过一遍，cancel 正好被注册一次；
- `:325-328` 行 = `live / Ctrl+Alt+K / Acc{0x4003,'K'}`（配置拼写与真正交给 Win32 的 accelerator 同源）；
- `:329-331` `rep.IsLive(hkCancel)` true；`:332-334` `len(rep.Live()) == 4`；
- `:335-337` `Problems()` 空；`:338-340` `idleBareKeyHeld245(rep)` 找不到裸键（AC#6② 的新锚）；
- `:341-343` `idleRosterReds245(rep)` 空；
- `:344-350` **宿主实际读的访问器**：`HotkeyReport().IsLive(hkCancel)` = true、`EscTakenOver()` = false（`cancelIdleSeam245`，`:397-409` 走 windowless ui-sta door，与 AC#7 rig 同一扇门）。

**正控 A（把判定改回「永远待命」）**：产品码只改 `:577` 一处，回到 `if p.id == hkCancel {`；`go build ./internal/ball` rc=0，
判据**逐字红**（原文存在 `.scratch/wisp/probes/245/ac6b/mutA-always-standby.txt`）：

```
$ go test -count=1 ./internal/ball -run TestConfiguredCancelIdleBindingFollowsModifierRuling245 -v
=== RUN   TestConfiguredCancelIdleBindingFollowsModifierRuling245
time=2026-11-14T14:29:35.467+08:00 level=INFO msg="cancel hotkey left unbound while idle: the configured key is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Ctrl+Alt+K
    hotkey_cancel_borrow_260_test.go:318: AC#6② RED: the idle pass made 0 RegisterHotKey attempts for a modifier-combination cancel (want exactly 1: a combination is held, not borrowed): [1 2 4]
--- FAIL: TestConfiguredCancelIdleBindingFollowsModifierRuling245 (0.01s)
```

**正控 A'（判定反相：`return acc.Mods&realModifiers == 0`）**：也 rc=0 可编译，同一行逐字红（
`mutB-inverted-predicate.txt`），红句仍是 `:318`、`[1 2 4]`。

## 3. 判据②（裸键一字不改）＋ 正控

**尺**：同一测试第 2 段（`:352-389`），配 `ApplyHotkeyDefaults(HotkeyConfig{Cancel: "Esc"})`：

- `:356-358` `attemptsOf(hkCancel) == 0`（idle 一次都没试）；
- `:359-361` `bindsEsc()` false（没有任何「裸 VK_ESCAPE」注册）；
- `:362-365` cancel 行 = `standby / "Esc" / 零 accelerator`；
- `:366-368` `len(rep.Live()) == 3`；
- `:369-372` `HotkeyReport().IsLive(hkCancel)` false；`:373-375` `EscTakenOver()` false；
- `:376-378` 借一次成功；`:379-381` 借只发生 1 次；`:382-386` 还之后 `attemptsOf(hkCancel)` 仍 1（**归还只 drop、不重绑**）；
  `:387-389` `UnregisterHotKey` 至少一次。

裸键的**既有尺全部照旧且全绿**（本腿未放宽）：`TestRegisterAllLiveSet`（`hotkey_status_test.go:120`，`fullConfig()` 的
cancel 是 `"Esc"`）、`TestDefaultHotkeysIdlePassHoldsNoEsc`（`:181`，正控逐字仍查 `DefaultHotkeys().Cancel == "Esc"`）、
`TestCancelBorrowRoundTrip`（`:208`）、`TestConfiguredCancelStillNeverBoundWhileIdle260` 的裸键半边已并入新尺第 2 段。

**正控 B（让裸键也被常绑）**：`cancelHeldWhileIdle` 改成「能解析就 true」（`:474-477` 换成 `_, err := ParseAccelerator(bind); return err == nil`），
`go build ./internal/ball` rc=0，第 1 段仍绿、第 2 段**逐字红**（`mutC-bare-also-held.txt`）：

```
$ go test -count=1 ./internal/ball -run TestConfiguredCancelIdleBindingFollowsModifierRuling245 -v
=== RUN   TestConfiguredCancelIdleBindingFollowsModifierRuling245
    hotkey_cancel_borrow_260_test.go:357: ticket 245 RED: the idle pass attempted a BARE cancel key (1 attempts): [1 2 3 4]
--- FAIL: TestConfiguredCancelIdleBindingFollowsModifierRuling245 (0.11s)
```

⇒ 两条判据都**不是**「写死会绿」的句子：改回待命红①，放开裸键红②。

## 4. 还原与哈希（三发突变之后）

```
$ Copy-Item .scratch/wisp/probes/245/ac6b/hotkey_windows.go.pristine internal/ball/hotkey_windows.go -Force
$ (Get-FileHash internal/ball/hotkey_windows.go -Algorithm SHA256).Hash
C66CE2E336A8402BA7389C272C1177B80408ED40C9B39D484E7719B7DE765C69     <- 突变前、还原后同串
$ (Get-FileHash internal/ball/ball_windows.go -Algorithm SHA256).Hash
C9F20037771274B59B7457E088605D69098F2231505F4E51D3007848378172A1     <- 本腿产品改动后、全程未再变
```

突变件与读数全部留在 `.scratch/wisp/probes/245/ac6b/`（`*.pristine`、`mutA-*.txt`、`mutB-*.txt`、`mutC-*.txt`、`gates.txt`、
`post-restore-test.txt`）。还原后 `go test -count=1 ./internal/ball` = `ok … 1.313s`。

## 5. 门禁（还原后现跑，`gates.txt` 逐行）

| 门 | 命令 | 读数 |
|---|---|---|
| vet | `go vet ./internal/ball` | rc=0 |
| 格式 | `gofumpt -l internal/ball`（`v0.12.0 (go1.27.1)`，`D:\work\base\gopath\bin\gofumpt.exe`） | 空（0 行），rc=0 |
| D22 | `sh scripts/d22scan.sh` | `clean - no D22 ban violations`，rc=0 |
| 包测 | `go test -count=1 ./internal/ball` | `ok github.com/CarlosShao/wisp/internal/ball 1.313s` |
| winlive 编译 | `go vet -tags winlive ./internal/ball` | rc=0（**只编译，未跑**；桌面不交出去） |

## 6. 本腿改了哪些既有尺（没有一条是放宽）

| 文件 | 改动 | 为什么／强度 |
|---|---|---|
| `hotkey_cancel_borrow_260_test.go:283-390` | `TestConfiguredCancelStillNeverBoundWhileIdle260` → `TestConfiguredCancelIdleBindingFollowsModifierRuling245` | 旧名的「configured cancel 永远不在 idle 注册」正是裁定**收窄掉**的字面。旧尺的裸键半边（借一次、还不重绑、unregister）**逐条保留**到新尺第 2 段；带修饰键半边按裁定**换向**（0 次→1 次、standby→live、3→4），并**新增** `HotkeyReport()`/`EscTakenOver()` 两个访问器读数。断言只换向／新增，未删 |
| `hotkey_refused_borrow_release_245_test.go` | AC#7 rig 的 planted binding `Ctrl+Alt+F8` → `"Esc"`（裸键）；闭包先问 `IsWindow`，不是窗口就**一步都不走** | 裁定后「借被拒」只存在于裸键这一支；rig 的断言（standby 归位、3 枚、Problems 空、幂等）**一字未动**。安全性从「用组合当替身」改成「先证明不是窗口再动手」，对裸键更强 |
| `hotkey_borrow_refused_260r5_test.go` | seed 前提：cancel 0 次/standby → 1 次/live（组合常绑）；同样加 `IsWindow` 门 | 260-r5 判的是「借被拒要记账」，refusal 断言（错误行、1400、两旗、邻居不受影响）**一字未动**；坏的只是「seed 里 cancel 必须 standby」这条与裁定冲突的前提 |
| `hotkey_live_test.go`（winlive） | `liveHotkeys()` 的 cancel 是 `Ctrl+Alt+V` ⇒ 计数 3→4（`:235`／`AllLive` 语境）、`requireEscBorrowed` 的 4 条不变、`requireEscReturned` 对「非默认且可解析」的 seed 改为「claimed 且是本进程自己的」 | 归还在裁定后**必须**把组合绑回；旧尺「组合必须在桌面 free」会在正确行为上红。默认（裸 Esc）那一支仍走 `requireIdleRoster` 的 free 读数，未加新探针 |
| `interaction_live_test.go` / `live_windows_test.go`（winlive） | 只改注释里「3→4」「还回后 Esc 不属于我们」的**描述** | 没有改任何断言表达式 |
| `hotkey_idle_bare_key_245_test.go` / `hotkey_status_test.go` | 注释里「idle 从不注册 cancel」改成「裸键不注册；组合注册」 | 判据函数 `idleRosterReds245` **未动**（它本来就按裸键判） |

## 7. 残余／未做（都点名，不粉饰）

1. **`cmd/wisp/main.go:28-34` 的 usage 现在对「带修饰键的 cancel」是半句假话**：现文逐字
   `the cancel slot is the fourth key, bound only while a confirmation card is waiting and handed back when it leaves,
   so an idle wisp holds no bare key against the rest of the desktop`。
   出厂默认（裸 Esc）下仍真，但用户配 `Ctrl+Alt+X` 后 cancel **稳态就在册**。`cmd/wisp/**` 是别腿的写面（本腿禁改）⇒
   **交编排者派一条文案腿**：改成「裸键才只在卡片等人批时借；带修饰键的组合常驻」。同理，宿主若有「hotkeys live 3/4」的
   打印，组合配置下会读 4/4（数字本身是对的行为，不是缺陷）。
2. **winlive 层未跑**（只 `go vet -tags winlive` rc=0 编译过）。若编排者要真桌面读数，建议干净桌面单发
   `go test -tags winlive ./internal/ball -run 'TestLiveCancelBorrowFollowsConfig260|TestLiveHotkeyRebindEndToEnd'`；
   本腿不交桌面，理由：本轮门禁不含 winlive，且组合常绑会让真实窗口短暂持有该组合。
3. **归还时的新失败面**：组合配置下 `ReleaseEscAfterSession` 会真的再 `RegisterHotKey` 一次；若此刻被别的程序抢走，
   report 会写一条**真**的 refusal（attempted/error，`Problems()` 出声），而不是 standby。这是诚实读数，不是旧 AC#7 那种
   「从没绑过却一直记着失败」；已在 `ball_windows.go:998-1035` 的注释里写明（`:1004-1018` 是这一支的理由）。
4. **AC#7 rig 的 fixture 改动**属于改写上一轮证据的输入（不是断言）；原文件未动，改动原因写在上表，供非实现者复核。

## 8. 撤销口令

一句话复原：**`245 AC6② 回永远待命`** —— 把 `internal/ball/hotkey_windows.go:577` 的条件改回 `if p.id == hkCancel {`，
并把 `ball_windows.go:1032` 改回 `cancelIdleLine(b.cancelBinding)`，即回到 `c40e296c9` 的判定；本件与测试改动一并按仓规
追加提交回退（⛔ 不 `--amend`／`reset`／`checkout .`）。
