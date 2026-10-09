# 票 296 — 296-r1（落地腿）起手锚 + 现跑行号复核

leg = 296-r1 · ticket = 296 · branch = dev · shared worktree
派单裁定 = **甲-全**（`hotReload258` 两枚 `return` 都过 `ball.ApplyHotkeyDefaults`）+ 同批改那句过期承诺 + `AC#1` 夹具先行（两形）

---

## 0. 起手闸门（逐字读数，尺 = 同一条 bash 命令的第一段）

```
$ date
Fri Oct  9 19:46:08 CST 2026

$ git log --oneline -1
7285ef83 293-a1 AC#0 （只读，零产码）：契约对「托盘勾号」＝空白 …（票 293 只读腿那一笔，本腿不引其内容）

$ git branch --show-current
dev

$ git status --porcelain -- cmd/wisp/
（空 —— 本腿开工前 cmd/wisp 无未提交改动）
```

- 工作树里 `design/**`／`frontend/**`／`.gitignore` 有别的腿未提交改动 ⇒ 本腿一律不碰、commit 必带显式 pathspec。
- 进程闸门（每发起手都要重跑，见 `AC#4` 节）：`tasklist //FI "IMAGENAME eq wisp.exe"` = `INFO: No tasks are running which match the specified criteria.`（0 枚）／`tasklist //FI "IMAGENAME eq balldebug.exe"` = 同（0 枚）。

---

## 1. 派单要求「行号自己现跑复核」——复核结果（锚 = 字面量，不是行号）

尺 = `grep -in "bridge keeps the current bindings" cmd/wisp/resident_windows.go` /
`grep -n "return ball.HotkeyConfig{}" cmd/wisp/resident_windows.go` /
`grep -n "Summon: c.Hotkey.Summon" cmd/wisp/resident_windows.go`

| 锚（字面量） | 现跑行号 | 裁定节记的行号 | 相符？ |
|---|---|---|---|
| `return ball.HotkeyConfig{}`（`hotCfg258` 内，读不到文件那支） | `resident_windows.go:188` | 未记 | — |
| `return ball.HotkeyConfig{}`（**`hotReload258` 内，读不到文件那支＝本腿第二形**） | `resident_windows.go:213` | `:213` | 相符 |
| `return ball.HotkeyConfig{Summon: c.Hotkey.Summon, …}`（**`hotReload258` 裸映射＝本腿主形**） | `resident_windows.go:215` | `:215`（票面原句 `:214`，裁定节已更正为 `:215`） | 相符 |
| `// The bridge keeps the current bindings on an empty answer`（过期承诺那一句的起行） | `resident_windows.go:208` | `:208-212` | 相符 |

桥的装配根那一行，尺 = `grep -rn "NewHotkeyReloader" cmd/wisp/resident_ball_windows.go cmd/balldebug/main.go`：

```
cmd/wisp/resident_ball_windows.go:353:		bridge := ball.NewHotkeyReloader(b, b.ConfiguredHotkeys(), hotReload)
cmd/balldebug/main.go:237:		bridge := ball.NewHotkeyReloader(b, b.ConfiguredHotkeys(), func() ball.HotkeyConfig {
```

⇒ 与票面「热加载那一半的消费者不补」的 `:353` 相符。

桥的跳过条件那把尺，裁定节记 `internal/ball/hotkey_reload.go:93` 逐字 `if r.exists && reflect.DeepEqual(r.applied, cfg) {`。
尺 = `grep -n "DeepEqual\|empty\|== \"\"" internal/ball/hotkey_reload.go` ⇒ 全文件**只命中这一行**：

```
93:	if r.exists && reflect.DeepEqual(r.applied, cfg) {
```

⇒ 复核成立：**零条 all-empty 守卫**，全空一旦与已应用集不同就真去重绑 ⇒ `:208-212` 那句注释是过期承诺（假话），本腿同批改它。

默认值合并那枚函数的生产调用点，尺 = `grep -rn "ApplyHotkeyDefaults" cmd/wisp/*.go | grep -v _test`：

```
cmd/wisp/resident_ball_windows.go:207://（注释）
cmd/wisp/resident_ball_windows.go:214://（注释）
cmd/wisp/resident_ball_windows.go:235:	cfg := ball.ApplyHotkeyDefaults(hotCfg())
cmd/wisp/resident_windows.go:193:			"empty_slots_note", "empty slots are filled from the product defaults by ball.ApplyHotkeyDefaults; "+
```

⇒ 与票面现量节相符：`hotReload258` 体内 **0 命中**；构造那一半的默认值补在消费者（`:235`）身上。

---

## 2. 本腿要动的两件事（照裁定节，不自选形）

1. **甲-全**：`hotReload258` 的**两枚** `return`（`:213` 与 `:215`）都过 `ball.ApplyHotkeyDefaults`。
   ⛔ 不动 `hotCfg258`（`:188`/`:195`）：它的默认值已由 `residentBallHotkeyChain258`（`resident_ball_windows.go:235`）补，
   且它的不可读那支带着一条 Warn + 一行 stdout —— 合并闭包 / 挪 Warn 都会改掉那句的出现时机（票面禁区第 4 条）。
2. **同批把 `:208-212` 那句注释改成说实话**：桥在「全空答案」上**不保持**当前绑定，它会真去重绑；
   本腿修完之后那枚闭包在全新装配下**不再答出四格全空**，注释按行为重写。
   ⛔ 不为此去动 `internal/ball`（乙案未裁）。

## 3. `AC#1` 夹具的形状（含一处必须具名的实现约束）

派单/票面要求「跑真装配：`startResidentBall(reg, nil, hotCfg, hotReload)` 用生产那两枚闭包的**同形**，⛔ 不许 `src, src`」。

⚠ **具名上报一条（不自选、不改判据）**：测试拿不到 `runResident` 里的**局部闭包本体**。
若照字面「同形」= 在测试里**重打一遍**那枚裸映射，则产码修没修都不影响这枚用例 ⇒ 它既不会「改前红」也不会「改后绿」，
是一枚**对这件事不敏感**的仪器（恒真/恒假尺），不满足票面 `AC#1` 的完成判据。
⇒ 本腿做法：把 `hotReload258` 的**本体**提为包级函数（`cmd/wisp/resident_windows.go`，同文件、同射程），
产码那枚闭包**调用它**，测试也**调用它** —— 这才是「跑真装配」。
先落「提取但**尚未**套 `ApplyHotkeyDefaults`」这一笔（用例必须红），再落「套上」那一笔（转绿）。
构造那一半仍按票面「同形」用 258 既有的 `loadHotkey258`（默认值由生产链 `:235` 补，与闭包本体无关），
并与热加载那一枚**分成两枚不同闭包**（⛔ `src, src`）。

两形都要交：
- 形一（机主的四行逐字）：`summon = ''` / `mute = ''` / `cancel = 'Esc'` / `panel = ''`。
- 形二（裁定节新增射程）：配置文件读不到 ⇒ `return ball.HotkeyConfig{}` 那一支。

断言：ⓐ 桥第一次 `Check()` 之后 `HotkeyReport().Live()` 仍含 summon/mute/panel（`live=3`）；ⓑ `ConfiguredHotkeys()` 读回合并后的值而非空串。
