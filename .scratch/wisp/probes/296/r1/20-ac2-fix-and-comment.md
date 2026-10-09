# 票 296 AC#2 — 落地那一发（甲-全：两枚 return 都过合并）+ 同批改那句过期承诺

leg = 296-r1 · 尺 = `git diff --stat HEAD` ＋ `sed -n` 现读 ＋ 目标用例两发改前改后对拉

## 1. 套了几个 return

`cmd/wisp/resident_windows.go` 里 `hotkeyReloadSource296`（＝ `hotReload258` 的本体，见 `00-anchor.md` §3）
**两枚 `return` 全部包过 `ball.ApplyHotkeyDefaults`**，不是只套裸映射那一枚：

```go
func hotkeyReloadSource296(dataDir string) ball.HotkeyConfig {
	c, _, err := config.LoadFile(filepath.Join(dataDir, configFileName), nil)
	if err != nil || c == nil {
		// Ticket 296: the unreadable branch answers the compiled defaults, not
		// four empty strings - the same fallback hotCfg258 prints a Warn for at
		// construction, said per tick by the value it hands the bridge.
		return ball.ApplyHotkeyDefaults(ball.HotkeyConfig{})
	}
	// Ticket 296: the three slots the owner's file leaves empty are "unset", and
	// the bridge would read them as "disabled by config". Merge before the hop.
	return ball.ApplyHotkeyDefaults(ball.HotkeyConfig{Summon: c.Hotkey.Summon, Mute: c.Hotkey.Mute, Cancel: c.Hotkey.Cancel, Panel: c.Hotkey.Panel})
}
```

⇒ 票面 AC#2 甲案的「只套 `:215` 那一行会漏掉 `:213` 这一形」被形二的夹具堵住（改前红、改后绿）。

**没动的**（禁区逐条）：`hotCfg258`（构造那一半，含它 `:223` 那条 Warn/printf 与它的裸映射 return——默认值仍由消费者
`residentBallHotkeyChain258`（`resident_ball_windows.go`）补，出现时机一字未变）、
`internal/ball/**`（0 字节）、`internal/config/schema.go`（0 字节）、`cmd/balldebug/main.go`（0 字节）、
机主的 `config.toml`（⛔ 未读未写；夹具用的是 `t.TempDir()`）、`HotkeyDisabled` 与 `none`/`off` 一族
（`Q-82`＝默认不做，本腿零 normalize）。

## 2. 那句过期承诺改成什么了（逐字新句）

**旧句（假话，逐字删掉）**：

```go
		// The bridge keeps the current bindings on an empty answer
		// (hotkey_reload.go leaves the applied set untouched when the source
		// errors into all-empty after the first diff), and the refresh Warn
		// it prints names the failure - the AC#1 fallback stays said, per
		// tick, in the log.
```

**新句（按行为写；函数文档注释里，逐字）**：

```go
// NO SLOT IS EVER ANSWERED EMPTY (ticket 296). Both branches merge through
// ball.ApplyHotkeyDefaults, which is the precondition ball.HotkeySource states
// in its own doc comment and the shape the construction half gets from
// residentBallHotkeyChain258. The bridge does NOT keep the current bindings on
// an empty answer: Check()'s only skip is the DeepEqual diff against what is
// applied, and internal/ball has no all-empty guard, so four empty strings from
// here would re-register the set and leave summon, mute and panel dead until
// the next read - which is exactly the regression this closure caused on the
// owner's machine (an empty [hotkey] slot means "unset", not "turned off").
```

判据问的是「空答案到底保不保持当前绑定」⇒ 答：**不保持**（`internal/ball/hotkey_reload.go` 的 `Check()` 里唯一的
跳过条件是 `reflect.DeepEqual`，全文件零条 all-empty 守卫；本腿现跑尺见 `00-anchor.md` §1）。
新句写的是这个事实，并且写明这一支修完之后不再答出全空。
分支体内另留两句短注释（`// Ticket 296: ...`）说明每一支为什么在这里合并。

`runResident` 里那枚调用点的注释也同批更正（它原写「the reload view goes to the bridge as-is」，
合并落地后那句会变假）：

```go
	// purpose: the construction view is read once and merged by its consumer
	// (residentBallHotkeyChain258), the reload view is merged by this source
	// itself because its consumer - the bridge - performs no merge at all
	// (ticket 296). Hoisting is what makes the merge measurable: a re-typed copy
	// of the closure inside a test reads green no matter what production does.
```

⛔ 没有为了迁就注释去动 `internal/ball`（乙案射程更大，编排者没裁乙）。

## 3. 目标用例：改前红 → 改后绿（逐字）

尺（两发同一把）＝`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp/ -count=1 -run 'Test296BridgeKeeps' -v`
起手闸门（每发都先跑）：`tasklist //FI "IMAGENAME eq wisp.exe"` = 0、`... balldebug.exe` = 0。

- 改前（HEAD=`28a2ff4e`，产码未合并）：rc=1，形一与形二**都红**，红句逐字含 `live=0`：
  ```
  AC#1 form-1 (owner's three empty slots) RED: after the bridge's first Check live=0 [], want live=3 (summon, mute, panel; missing [summon mute panel]), ...
  AC#1 form-2 (config.toml unreadable at tick time) RED: after the bridge's first Check live=0 [], want live=3 ...
  ```
  进程自己的读数与真机那组逐字同形：`msg="ball: hotkeys rebound after config change" summon="" mute="" cancel=Esc panel="" live=0`（形一）／`summon="" mute="" cancel="" panel="" live=0`（形二）。
  全文＝`logs/red-before-fix.md`，逐行摘录＝`10-ac1-red-before-fix.md`。
- 改后（本笔）：rc=0，两枚全绿：
  ```
  --- PASS: Test296BridgeKeepsThreeLiveKeysWithOwnerEmptySlots (0.31s)
  --- PASS: Test296BridgeKeepsThreeLiveKeysWhenConfigUnreadable (0.04s)
  ok  	github.com/CarlosShao/wisp/cmd/wisp	0.444s
  ```
  全文＝`logs/green-after-fix.md`。**改后名册里 `hotkeys rebound after config change` 一句 0 命中**（尺＝
  `grep -o "hotkeys rebound after config change.*live=[0-9]*" logs/green-after-fix.md` ⇒ `No matches found`）：
  合并后的四格与建球时那四格逐字相等 ⇒ `Check()` 的 `DeepEqual` 跳过 ⇒ 零次 Win32 重绑。
  这正是票要的形状：桥不再每秒把用户刚绑上的键拆掉。

## 4. AC#3 归编排者跑

票面 `AC#3`（同一台机器、同一份 `config.toml` 一字不改，开机应见 `hotkeys_live=3` 且不再出现 `live=0`，
再按一次 `Ctrl+Alt+M` 取开门/关门读数）按派单**归编排者亲跑**；本腿⛔ 未起 `build/wisp.exe`、⛔ 未动 `build/**` 里的 exe。
