# 296-a1 / AC#0 — 射程普查名册：宿主把 `[hotkey]` 裸映射进 `ball` 的所有落点

尺与原始读数见 `00-anchor-and-rulers.md`（R1/R1b/R2/R3/R4/R8）。
名册口径：**R1 只跑 `cmd internal` 两棵树**；R1b 另跑一次 `.scratch` 确认变异拷贝存在并**显式排除 8 枚**
（`258/v1/mut0` 4 枚 + `258/v1/mut3` 4 枚，全是 `cmd/wisp/resident_windows.go` 的整文件变异拷贝，不参与生产路径）。
**本名册不计那 8 枚，也不勾票 296 的任何框。**

## 一、真落点 6 枚（要判"有没有过 `ApplyHotkeyDefaults`"的映射/消费点）

### 落点 A1 ＝ `cmd/wisp/resident_windows.go:181` `hotCfg258`（构造期源）——〔已过默认，但过在消费者身上〕
- **四格是否可能为空串**：**可能**，两种形：
  - `:195` 逐字 `		return ball.HotkeyConfig{Summon: h.Summon, Mute: h.Mute, Cancel: h.Cancel, Panel: h.Panel}`（`h := c.Hotkey`，`:190`）
    ⇒ 盘上有机主的三格空串实证（尺＝`sed -n '/^\[hotkey\]/,/^\[/p' C:/Users/swq/AppData/Roaming/wisp-dev/config.toml`，本腿 18:1x 现跑，逐字四行＝`summon = ''`／`mute = ''`／`cancel = 'Esc'`／`panel = ''`）；
    另一形＝`internal/config/schema.go:180-183` 只有 `Cancel` 带 `default:"Esc"` 标，`Summon`/`Mute`/`Panel` 无标 ⇒ 文件里没有 `[hotkey]` 节时也恰好是三空串（`cmd/wisp/resident_ball_windows.go:212-214` 逐字复述过这一形）。
  - `:188` 逐字 `			return ball.HotkeyConfig{}`（文件读不到时的分支）⇒ **四格全空**。
- **谁补默认**：`cmd/wisp/resident_ball_windows.go:235` 逐字 `	cfg := ball.ApplyHotkeyDefaults(hotCfg())`（在 `residentBallHotkeyChain258` 内，`:231` 定义、`:307` 被 `startResidentBall` 调）。⇒ **本枚合规**（前置条件由消费者兑现）。

### 落点 A2 ＝ `cmd/wisp/resident_windows.go:205` `hotReload258`（热加载源）——★**未过默认，本票回归点**
- **四格是否可能为空串**：**必然可能**。`:215` 逐字 `		return ball.HotkeyConfig{Summon: c.Hotkey.Summon, Mute: c.Hotkey.Mute, Cancel: c.Hotkey.Cancel, Panel: c.Hotkey.Panel}`
  ⇒ 与 A1 的 `:195` 同形裸映射，机主那份文件下回答 `summon="" mute="" cancel="Esc" panel=""`；
  `:213` 逐字 `				return ball.HotkeyConfig{}` ⇒ 文件读不到时四格全空（**这一支本腿另有一问，见文末"本腿多量到的一枚"**）。
- **谁补默认**：**没有人**。尺＝R2：`ApplyHotkeyDefaults` 在 `cmd/wisp` 的**唯一真调用**是 `resident_ball_windows.go:235`；
  `:353` 逐字 `		bridge := ball.NewHotkeyReloader(b, b.ConfiguredHotkeys(), hotReload)` 第三实参就是那枚裸闭包。
  被调方**明文要求前置条件**：`internal/ball/hotkey_reload.go:48-49` 逐字
  `// HotkeySource returns the currently effective bindings (the host's mapping` / `// of [hotkey], already defaulted via ApplyHotkeyDefaults).`
  ⇒ **宿主违反了它自己引用的前置条件**（编排者现量复核＝成立）。
- 落点后果链（逐字，本腿自己读到）：`internal/ball/hotkey_windows.go:494-496`
  `		case b.Binding == "":` / `			b.Status = HotkeyDisabled // disabled by config: never attempted` / `			slog.Info("hotkey disabled (unset in config)", "hotkey", b.Name)`
  ⇒ 桥**照规矩**把空串读成"配置里关了"，桥没错。

### 落点 A3 ＝ `cmd/balldebug/main.go:237-242`（balldebug 的桥源）——〔已过默认 ⇒ 编排者那枚反例**成立**〕
- 本腿现读逐字（`sed -n '230,258p' cmd/balldebug/main.go`）：
  ```
  bridge := ball.NewHotkeyReloader(b, b.ConfiguredHotkeys(), func() ball.HotkeyConfig {
      h := mgr.Config().Hotkey
      return ball.ApplyHotkeyDefaults(ball.HotkeyConfig{
          Summon: h.Summon, Mute: h.Mute, Cancel: h.Cancel, Panel: h.Panel,
      })
  })
  ```
- **四格是否可能为空串**：闭包**返回值**不可能为空（`ApplyHotkeyDefaults` 的 `:95-106` 四支把空槽填成 `DefaultHotkeys()`）；
  内层字面量本身可以为空（同一份机主配置），但**在返回前已被合并**。
- **谁补默认**：`cmd/balldebug/main.go:239` 自己。⇒ **"同一个坑 balldebug 没有、只有 wisp 有"这一句我独立复核＝成立**，不是采信编排者的转述。
  注意形状差别：balldebug 把默认合并放在**源闭包体内**（正是 `hotkey_reload.go:16-20` 示例教的写法），wisp 构造期放在**消费者**、热加载期**两处都没放**。

### 落点 A4 ＝ `cmd/wisp/resident_ball_windows.go:316` `		Hotkeys:  cfg,`（`ball.New` 的实参）
- **四格是否可能为空串**：**不可能**（`cfg` 来自 `:307 cfg, provenance := residentBallHotkeyChain258(hotCfg)`，已过 `:235`）。
- **谁补默认**：A1 的消费者那一层。⇒ 合规。顺带：包内还有第二重保险 `internal/ball/ball_windows.go:153-154`
  `	if opts.Hotkeys == (HotkeyConfig{}) {` / `		opts.Hotkeys = DefaultHotkeys()`——**只在四格全空时生效**，
  ⇒ 机主那形（三空一有值）**不会被它救**，所以这重保险不能替代 `ApplyHotkeyDefaults`。

### 落点 A5 ＝ `cmd/wisp/resident_ball_windows.go:353`（wisp 侧 `NewHotkeyReloader` 的装配点）
- 与 A2 同事件的两面：这一枚是"谁把没默认的源交给桥"。`ConfiguredHotkeys()`（第二实参）取自 `internal/ball/ball_windows.go:836-844`
  ⇒ 是"上一次注册被喂进去的值"，boot 时＝合并后的默认（非空），所以桥的 `applied` 初值非空、第一次 tick 必然判成"变了" ⇒
  **开机约 1 秒后必发一次 rebind**，这与真机读数 `live=3 → live=0` 的 1.003 秒完全对得上。
- **谁补默认**：无。

### 落点 A6 ＝ `internal/ball/ball_windows.go:813 func (b *Ball) RebindHotkeys(cfg HotkeyConfig)`（被调方的注册入口）
- **四格是否可能为空串**：**可能**——它不做任何默认合并，直接 `registerAll`（`:817`）。生产里唯一的调用者是
  `internal/ball/hotkey_reload.go:96` `	rep := r.binder.RebindHotkeys(cfg)`，`cfg` 完全由宿主源决定 ⇒ 洞的"能不能补"取决于源，不取决于这里。
- **谁补默认**：设计上是宿主的义务（`hotkey_reload.go:48-49` 那句前置条件）。⇒ **甲乙两案的分工分歧就落在这一枚**（见 AC#2 件）。

## 二、非落点命中（按档位标出；本腿**不"消除"**它们）

| # | 位置 | 档 | 逐字要点 |
|---|---|---|---|
| B1 | `internal/ball/hotkey_reload.go:16-20` | 〔注释＝用法示例〕 | `//	r := ball.NewHotkeyReloader(b, b.HotkeyConfig(), func() ball.HotkeyConfig {` / `//		return ball.ApplyHotkeyDefaults(ball.HotkeyConfig{` ——**正例**，wisp 的 A2 与它不同形 |
| B2 | `cmd/wisp/resident_windows.go:193` | 〔日志文案〕 | `"empty_slots_note", "empty slots are filled from the product defaults by ball.ApplyHotkeyDefaults; "` ——这句话**只对 A1 那半成立**；同一进程里 A2 那半没有兑现它，属**文案与实际链路不一致**（本腿只报不修） |
| B3 | `cmd/wisp/resident_ball_windows.go:207`／`:214`／`:224` | 〔注释体＝编排者说的"第二枚"〕 | `// maps the host's [hotkey] view through ball.ApplyHotkeyDefaults and names the` / `// (258-a1 census Q1a), and ApplyHotkeyDefaults merges exactly` ——`residentBallHotkeyChain258` 的 doc 注释与三档说明，**不是代码路径**；按派单口径标〔注释/日志文案〕档，⛔ 本腿不动它 |
| B4 | `internal/ball/hotkey_windows.go:75`／`:93` | 〔默认表／被调函数本体〕 | `return HotkeyConfig{Summon: "Ctrl+Alt+Q", Mute: "Ctrl+Alt+M", Cancel: "Esc", Panel: "Ctrl+Alt+P"}` / `func ApplyHotkeyDefaults(cfg HotkeyConfig) HotkeyConfig {` ——映射落点之外的定义处，尺顺带命中 |

## 三、本腿多量到的一枚（编排者现量里没有，具名报回）

**`cmd/wisp/resident_windows.go:213` 那个"文件读不到 ⇒ 返回四格全空"的分支，其注释承诺的行为在桥里不存在。**
- `:208-212` 逐字：`//				// The bridge keeps the current bindings on an empty answer` /
  `//				// (hotkey_reload.go leaves the applied set untouched when the source` /
  `//				// errors into all-empty after the first diff), and the refresh Warn`
- 现读 `internal/ball/hotkey_reload.go:81-105` 的 `Check()`：唯一的跳过条件是 `:93`
  `	if r.exists && reflect.DeepEqual(r.applied, cfg) {` ⇒ **"全空"不是跳过条件**，全空照样走到 `:96` 的 `RebindHotkeys(cfg)`。
  全文（`:1-174`）尺 `grep -n "HotkeyConfig{}" internal/ball/*.go` 只命中 `ball_windows.go:153` 与 `hotkey_windows.go:75`，
  **`hotkey_reload.go` 里 0 处 all-empty 守卫** ⇒ 那句注释是**过期承诺**。
- 实际后果（同一枚洞的第二形）：boot 后若 `config.toml` 变得不可读（被移动、被写坏、权限变化），`hotReload258` 回答四格全空，
  桥与 `applied`（非空默认）判为"变了" ⇒ **照样 `live=0`**。这一形与票 296 主形**同链不同支**：主形修好了（甲在 `:215` 那一支套默认）**不覆盖 `:213`**。
- 因此给编排者的裁法约束（**不是本腿挑案**）：甲若只把 `ApplyHotkeyDefaults` 套在 `:215` 那一支，`:213` 那一支仍是 `live=0` 的来源；
  要么甲套在闭包**整个返回值**上（含错误分支），要么乙在 ball 侧给 `Check()` 加"全空＝源没回答，保持现状"的守卫。两种都得在 AC#1 夹具里跑得到。

## 四、票 296「现量」一节的独立复核结果（⛔ 不是照抄）

| 编排者的话 | 本腿现跑读数 | 判定 |
|---|---|---|
| 机主 `[hotkey]` 四行＝`summon=''`/`mute=''`/`cancel='Esc'`/`panel=''` | `sed -n '/^\[hotkey\]/,/^\[/p'` 现跑逐字相符 | ✔ 成立 |
| `resident_windows.go:181` 构造源裸四格，消费者 `resident_ball_windows.go:235` 有 `ApplyHotkeyDefaults` | 行号逐字相符（`:181`、`:195`、`:235`） | ✔ 成立 |
| `resident_windows.go:205-214` 热加载源裸四格 | 闭包起点 `:205` ✔；**裸字面量实际在 `:215`**（`:213` 是全空分支） | ⚠ 行号偏 1，实质成立 |
| `resident_ball_windows.go:353` 直接交给 `NewHotkeyReloader`，链路无 `ApplyHotkeyDefaults` | `:353` 逐字相符；R2 命中分布确认 `cmd/wisp` 真调用只有 `:235` | ✔ 成立 |
| `hotkey_reload.go:48-49` 明文前置条件 | 逐字相符 | ✔ 成立 |
| `hotkey_windows.go:494-496` 空串＝"配置里关了" | 逐字相符（`:494`/`:495`/`:496`） | ✔ 成立 |
| `ApplyHotkeyDefaults` 注释 `hotkey_windows.go:78-92`（票内写 `:78-85`）预言了 "straight in" 那一形 | 注释整段＝`:78-92`；票 296 正文 `:78-85` 少算一行；逐字含 `// host that maps config.Hotkey straight in would silently end up with NO` / `// summon key at all on a fresh install` | ⚠ 行号偏，实质成立 |
| balldebug `:237-241` 那枚 src **有** `ApplyHotkeyDefaults` ⇒ 反例 | `:237` 是 `NewHotkeyReloader`、`:239` 是 `ball.ApplyHotkeyDefaults(ball.HotkeyConfig{`（闭包到 `:242` 结束） | ✔ 成立（本腿自读，未采信转述） |
| 引入点＝`5e8748b3`（2026-10-03） | `git log -S hotReload258 -- cmd/wisp/resident_windows.go --oneline` | ✔ 成立（唯一一条命中） |
| 真机读数 `:18` `hotkeys_live=3` / `:27/:28/:30` 三行 disabled / `:31` `live=0`，差 1.003 秒 | `grep -n` 现跑：`:18` 与 `:31` 逐字命中（`17:27:04.936` → `17:27:05.939`＝**1.003 秒**）；三行 `hotkey disabled (unset in config)` 在 `:27`(summon) `:28`(mute) `:30`(panel)，`:29` 是 cancel 的 standby 行 | ✔ 成立 |
| 258 用例夹具四格全显式赋值、且 `src, src` 同源；`:153` 把"源可全空"钉成期望 | 夹具实际在 `resident_hotkey_258_windows_test.go:272-275`（票写 `:266-281`，那是"含函数名的整块"）；`rb := startResidentBall(observe.NewRegistry(), nil, src, src)` 实际在 **`:278`**（票写 `:283`）；`:153` 那句 `	if c := loadHotkey258(t, dir); c != (ball.HotkeyConfig{}) {` **行号逐字相符** ⇒ 成立 | ⚠ 两处行号偏（约 −5 行），实质成立 |
| 尺 `grep -rn 'mute = ..' --include=*_test.go cmd/wisp` ⇒ 没有任何一枚桥用例写空串 | 本腿跑 `grep -rn "summon = ..|mute = ..|panel = ..|cancel = .." --include=*_test.go cmd/wisp` ⇒ 命中全是**非空显式组合**（`Ctrl+Alt+Z/X/V/R/Q/M/P`、`Esc`），**0 枚空串夹具** | ✔ 成立 |

⇒ **不符之处只有"行号"三类（`:215` 而非 `:214`；注释段 `:78-92` 而非 `:78-85`；测试 `:272-275`/`:278` 而非 `:266-281`/`:283`），实质判断零处推翻**；
另**多量到**第三节那一枚（`:213` 全空分支的过期注释承诺 + 独立的 `live=0` 触发形）。
