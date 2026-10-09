# 296-a1 / AC#2 — 甲乙两形代价表 ＋ "`none`/`off` 那一问"的整族穷举答案（仍只读）

尺名册见 `00-anchor-and-rulers.md`；落点编号沿用 `10-ac0-roster.md`（A1…A6、B1…B4）。
本件**不挑甲也不挑乙**——那是编排者的裁量；本件只把两形的射程与那一问的答案量清楚。

---

## 一、甲（改宿主：`hotReload258` 体内套 `ball.ApplyHotkeyDefaults(...)`，与 `resident_ball_windows.go:235` 同形）

**要动的文件（生产码）**
| 文件 | 动哪一处 | 射程 |
|---|---|---|
| `cmd/wisp/resident_windows.go` | 落点 A2 的闭包体（`:205-216`）——**唯一必需改动** | 2 行量级 |
| 同上（可选，见下"甲的两个亚形"） | `:208-212` 那段过期注释 | 文案 |

- **`cmd/wisp/resident_ball_windows.go` 零字节**（消费者侧不动，`:235` 那半本来就对）；
- **`cmd/balldebug/main.go` 零字节**——它本来就是甲的形状（`main.go:239` 逐字 `			return ball.ApplyHotkeyDefaults(ball.HotkeyConfig{`），**甲不会让 balldebug 变**；
- **`internal/ball/**` 零字节**、**`internal/config/schema.go` 零字节**（票 296 禁区两条都满足）；
- ⇒ **与票面 AC#4 的越界判据"名册只含 `cmd/wisp/`（甲案）＋ `probes/296/**`"完全吻合，不用改票面。**

**甲的两个亚形（差别就是本腿多量到的那一枚，见 `10-ac0-roster.md` 第三节）**
- **甲‑窄**：只把 `ApplyHotkeyDefaults` 套在 `:215` 那一支 ⇒ 修好机主那份文件（票面主形），但 `:213` 那个"文件读不到 ⇒ `return ball.HotkeyConfig{}`"**仍是 `live=0` 的入口**：桥的 `applied` 初值是合并后的默认（非空），全空 ≠ 相等 ⇒ 过 `internal/ball/hotkey_reload.go:93` 的 `reflect.DeepEqual` 判"变了"、走 `:96` `RebindHotkeys(cfg)` ⇒ 一次配置不可读就把三枚键拆掉。
- **甲‑全**：闭包整体返回值套默认（两支都过）⇒ 主形与 `:213` 那支一起修掉；代价是"文件读不到"时桥会把默认值**再绑一遍**（值不变 ⇒ `DeepEqual` 相等 ⇒ 不碰 Win32，`:93-95` 那条"unchanged file costs zero Win32 calls"仍然成立），语义变成"读不到就用默认"，与 `:188` 构造期那支的注释承诺（"falling back to the compiled defaults"）反而**一致**了。

**`[hotkey] mute = 'none'` / `'off'` 这一形将来怎么表达**
⇒ **甲既不支持、也不堵死。**现量：`'none'` 不是空串 ⇒ `ApplyHotkeyDefaults`（`internal/ball/hotkey_windows.go:95-106`）原样放过 ⇒ `ParseAccelerator`（`:131-…`，token 分支里没有 `none`/`off`）落到 `default:` 报错 ⇒ `:500` `b.Status = HotkeyUnparsable` ＋ `:502` `slog.Error("hotkey not attempted (bad binding)")`。
也就是说这一形今天**读起来是"配置写坏了"，不是"我关掉了它"**。将来若机主拍板要"关"，落点在**宿主闭包里做一次 normalize**（`none`/`off` → 一枚哨兵），射程仍只在 `cmd/wisp`，`internal/ball` 一字不动 ⇒ **甲案对"将来支持关键"是开放姿态**。

---

## 二、乙（改被调方：`NewHotkeyReloader`／`Check()` 对 src 返回值过一次 `ApplyHotkeyDefaults` ＋ ball 侧写文档＋用例）

**要动的文件（生产码）**
| 文件 | 动哪一处 | 射程 |
|---|---|---|
| `internal/ball/hotkey_reload.go` | `:48-49` 那句前置条件从"要求宿主"改成"我自己保证"；`Check()` 的 `:90 cfg := r.src()` 之后套一次默认（或 `NewHotkeyReloader:73-75` 装 src 时包一层） | 1 处逻辑 |
| `internal/ball/hotkey_windows.go` | `:54` 逐字 `// HotkeyConfig is the parsed [hotkey] section. Empty string = disabled.` 与 `:78-86` 的 "empty = unset" 必须**改写对齐**（今天这两句自己就不一致） | 文档/不变式 |
| `cmd/wisp/resident_windows.go` | `:208-212` 那段"the bridge keeps the current bindings on an empty answer"过期注释**必须一起改**（否则注释与代码再次分家，而那正是乙案声称要修的病） | 文案 |
| `cmd/balldebug/main.go` | `:239` 的 `ApplyHotkeyDefaults` 变**双重合并**（幂等、行为不变，但是死代码）；留着＝两份真相，删掉＝balldebug 也进名册 | 0 或 1 处 |

- **测试面**：乙要新增 `internal/ball` 那包的用例（票面 AC#2 的原文就是这么写的："把'空＝没配'这条不变式在 ball 侧写成文档＋用例"）⇒
  **票面 AC#4 的门禁名册要跟着扩**（现在只列了 `go test ./cmd/wisp/`；乙改了 `internal/ball` 就得再加 `go test ./internal/ball/`），
  且 `git show --stat` 的"名册只含 `cmd/wisp/`＋`probes/296/**`"这一条越界判据**在乙案下会红**⇒ 需要编排者先改票面再派 `296-r1`（本腿不动票面框）。
- **契约面风险（具名）**：`internal/ball` 的 `HotkeySource` 语义（`hotkey_reload.go:48-49`）与 `HotkeyConfig` 的"Empty string = disabled"（`hotkey_windows.go:54`）都是**被别处引用的说法**——
  `cmd/wisp/resident_ball_windows.go:212-214`、`cmd/wisp/config_readers_255.go:129` 都在引用这条链路的语义；乙动 ball 侧 ⇒ 那两处宿主注释也可能要复核。射程比甲大一个包。
- 乙的亚形也要裁：是"src 返回后套 `ApplyHotkeyDefaults`"（把 `:213` 那支变成默认键），还是"src 返回全空 ⇒ 保持现状不 rebind"（把 `:208-212` 那句过期承诺**实现成真的**）。
  两形对 `:213` 的处置**不同**，且只有后者才让 `Rebinds()`（`hotkey_reload.go:167-173`）那枚"不许每 tick 空转"的接缝继续成立。

**`[hotkey] mute = 'none'` / `'off'` 这一形将来怎么表达**
⇒ **乙把它钉得更死。**乙案的本质是"把'空＝没配、由 ball 补默认'升级成包内不变式"，于是：
- `''` ⇒ 永远被 ball 侧填回默认（用户彻底无法表达"关"）；
- `'none'`/`'off'` ⇒ 仍然 `HotkeyUnparsable`（同甲）；
- 将来若要支持"关"，**必须改 `internal/ball` 的公开契约**（`HotkeySource` 的含义、`HotkeyConfig:54` 那句注释、`HotkeyDisabled` 的产生条件三处一起动），比甲多跨一个包、且碰到的是被别处引用的语义。

---

## 三、★本格最要紧那一问的整族穷举答案

> **问：仓里今天有没有任何一条生产路径，能让用户真正把一枚热键关掉（而不是被重新填上默认值）？**

### 尺（整族穷举，非抽样；范围 `internal/ball`＋`cmd/wisp`＋`internal/config`，⛔ `_test.go`、⛔ `.scratch/**`）

```
R5  grep -rn '"none"\|"off"\|"disabled"\|"false"\|"no"' --include=*.go internal/ball cmd/wisp internal/config | grep -v _test
R6  grep -rn "HotkeyDisabled\|disabled by config\|Unregister\|== \"\"" --include=*.go internal/ball cmd/wisp internal/config | grep -v _test
R11 grep -rn "Hotkey\.\(Summon\|Mute\|Cancel\|Panel\) *=\|Hotkeys\.\(Summon\|Mute\|Cancel\|Panel\) *=" --include=*.go cmd internal | grep -v _test
R12 grep -rniE "\"(none|off|disable|disabled|nokey|no-key|unbind|unbound)\"" --include=*.go internal/ball cmd/wisp internal/config | grep -v _test
R13 grep -rn "HotkeyDisabled\|HotkeyStandby\|HotkeyUnparsable" --include=*.go cmd internal | grep -v _test
R14 grep -rn "Hotkey\|hotkey" --include=*.go internal/config/loader.go internal/config/defaults.go internal/config/writeguard.go internal/config/tiers.go
R9  grep -rn "Hotkey" --include=*.go internal/config | grep -v _test
```

### 读数要点（逐字）
1. **热键族里 `"none"`/`"off"` 的全部命中只有 1 枚，且它是注释**：
   `internal/ball/hotkey_windows.go:84` 逐字 `// through here; a config with summon = "none"/"off" is the host's business`
   ＋ `:85` 逐字 `// (it arrives as "" and re-enables the default, which is the safe direction: `
   ＋ `:86` 逐字 `// a live hotkey the user can see, never a dead silent one).`
   ⇒ 这三行就是编排者点名的那一句，**实际行号是 `:84-85`（派单写作 `:83-84`，票 296 正文写作 `:78-85` 段），实质一致**。
   R5/R12 里其余 `"none"`/`"off"` 命中全部**属于别的族**，逐条排除：`internal/ball/dock_windows.go:336 case sameName(s, "none"):`（停靠边 `top|bottom|none`）、
   `internal/config/schema.go:70 ProxyNone = "none"` 与 `:496-497`（`[net]` 代理枚举）、`schema.go:85 ThinkingOff = "off"` 与 `:291 default:"off"`（思考强度）、
   `internal/config/manager.go:478`（`[net]` 松紧判定）、`cmd/wisp/config_reload.go:141`（定时任务那枚 "disabled" 的措辞）。
2. **代码里没有任何一处给热键槽赋值或归一化**：R11 **0 命中**。整个仓里对 `HotkeyConfig` 唯一的改写只有两枚，方向都是"补默认"：
   `internal/ball/hotkey_windows.go:95-106`（四支 `if cfg.X == "" { cfg.X = d.X }`）与 `internal/ball/ball_windows.go:153-154`
   `	if opts.Hotkeys == (HotkeyConfig{}) {` / `		opts.Hotkeys = DefaultHotkeys()`（**只救四格全空**，机主那形三空一有值救不到）。
3. **`HotkeyDisabled` 这一枚状态是真的写出来的、且真的会被打**：
   产生点只有 2 枚——`internal/ball/hotkey_windows.go:455`（cancel 槽，`cancelIdleLine`）与 `:495`（其余三槽），
   判据都是"传进来的 `Binding` 是空串"（`:454 case bind == "":` / `:494 case b.Binding == "":`）。R13 全量：`cmd/wisp` 与 `cmd/balldebug` **都不产、不读** `HotkeyDisabled`，只有 `internal/ball` 内部在产与在格式化（`:273`、`:360`、`:388`）。
4. **`config` 层对 `[hotkey]` 只做三件事**（R9/R14）：`schema.go:179-184` 定四格（只有 `Cancel` 带 `default:"Esc"`）、
   `tiers.go:38 "hotkey": "hot"`、`manager.go:281` 热档整节替换。
   **loader/defaults/writeguard 里没有任何针对 `[hotkey]` 的空串/哨兵处理**（R14 只命中 `writeguard.go:14` 的一句注释与 `tiers.go:38`）。
5. **而 `config` 层的文档恰恰许诺了"空＝关"**：`internal/config/schema.go:177` 逐字
   `// Empty string = binding unset (feature key disabled). ` ——与 `internal/ball/hotkey_windows.go:54`
   `// HotkeyConfig is the parsed [hotkey] section. Empty string = disabled.` 同调，
   **但两句都被 `ApplyHotkeyDefaults:95-106` 在链路上就地推翻了。**

### 分三档回答

- **〔已实现并接上〕＝ 没有。**不存在任何一条**由用户表达、被代码理解、并被如实上报**的"关掉一枚键"的路径。
  理由就是 1–5：热键族里没有任何赋值/归一化点（R11 0 命中），所有对空串的处理方向都是"填回默认"。
- **〔写了但没人调用〕＝ 半档成立，且要说准**：`HotkeyDisabled` 这套**下游形状是写好的、也是活的**（`:494-496` 每次都真的执行），
  缺的是**上游**——生产里今天唯一能把空串喂给它的路径，**恰好就是本票要修的 `hotReload258` 那枚裸闭包（落点 A2）**。
  ⇒ 也就是说：**"关得掉"这个能力现在是靠 bug 在兑现的**；票 296 一落地（甲或乙都堵住 A2），`HotkeyDisabled` 在生产里就**再无入口**（`:213` 那支除外，它同样被甲‑全／乙覆盖）。
- **〔契约/注释许过但代码里没这个形状〕＝ 这就是 `none`/`off` 那一问的落点**。
  `"none"`/`"off"` 只活在 `hotkey_windows.go:84` 那一行注释里，而且那行注释把它们**明确推给宿主**（"is the host's business"），
  宿主侧（`cmd/wisp`、`cmd/balldebug`）**没有任何一处实现过这件事**（R5/R12 在 `cmd/wisp` 的热键族命中＝0）；
  同时 `schema.go:177` 与 `hotkey_windows.go:54` 两句文档还在说"空＝关"，代码却是"空＝默认"。
  ⇒ **文档互相矛盾、代码无对应形状 ⇒ 这是契约空白，不是实现空白。**

### 对排程的直接结论（不替人拍板）

**要摆机主一句话。**本票 AC#2 落地的**任何一形**（甲窄／甲全／乙）都只决定"空串怎么被读"，
都不新增"关得掉"；而"用户到底能不能主动关掉一枚全局热键"今天**没有定案**：
`PLAN.md`/`SPEC-08` 那侧许诺的是"有一枚一键静音热键"，`schema.go:177` 那侧许诺的是"空＝关"，两者在代码里由 `ApplyHotkeyDefaults` 统一读成"空＝默认"。
⇒ 建议编排者把这一问上机主清单时带上三栏：①维持现状（永远关不掉，"关"用别的入口表达）②`none`/`off` 折成"哨兵＝关"（甲案下只动 `cmd/wisp`）③给 ball 加第四态或让 `HotkeyDisabled` 有可寻址的表达（乙案面，跨包、动 `:54` 那句文档）。
⛔ 本腿不挑任何一栏，也不勾框。
