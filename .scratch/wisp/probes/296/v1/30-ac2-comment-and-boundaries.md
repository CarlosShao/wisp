# 296-v1 / AC#2 —— 修法那一发＋同批改写的注释：判语 **部分成立**（修法成立；注释两处读数与盘上不符，不构成回归）

HEAD＝`28b484ec`。被裁件＝`fc4aedaf`（甲-全那一发）里 `cmd/wisp/resident_windows.go` 的 39 行改动。

## 1. 修法那一发：**成立**

- 两支 `return` 都过了合并（内容锚，⛔ 不用绝对行号）：
  `return ball.ApplyHotkeyDefaults(ball.HotkeyConfig{})`（不可读支）与
  `return ball.ApplyHotkeyDefaults(ball.HotkeyConfig{Summon: c.Hotkey.Summon, Mute: c.Hotkey.Mute, Cancel: c.Hotkey.Cancel, Panel: c.Hotkey.Panel})`（裸映射支）。
  尺＝V-R1／V-R2（见 `10-ac0-roster-check.md`）⇒ `cmd/wisp` 非测试侧 `ApplyHotkeyDefaults` 真调用现为 3 处（`:235` 构造消费者＋`:49`＋`:53`）。
- "套壳真的把四格都填上"不是想当然：`internal/ball/hotkey_windows.go:93-108` 逐字四支 `if cfg.Summon/Mute/Cancel/Panel == "" { … = d.X }`
  ⇒ 注释那句 **"NO SLOT IS EVER ANSWERED EMPTY"** 对这两支**为真**（四格全覆盖，没有第五格）。
- 只套一支会漏：突变体 B／C 各红一枚（见 `20-ac1-mutation.md`）⇒ 编排者 18:2x 那条"甲射程含第二枚 `return`"被夹具钉住。
- `hotCfg258` 那一半**没被顺手改**（`:223/:230` 仍是裸返回，由 `resident_ball_windows.go:235` 的消费者补），
  它的 Warn 时机也未动（`sed -n '214,232p' cmd/wisp/resident_windows.go` 读到 `slog.Warn("ball: [hotkey] source unreadable at construction…")`＋`fmt.Printf` 逐字仍在）
  ⇒ 票面禁区"不许把两半合并成一枚闭包"未被违反。

## 2. 边界自查（四把尺，逐条）

| 问 | 尺（逐字） | 读数 | 判定 |
|---|---|---|---|
| `internal/ball` 一字未动？ | `git diff --stat ff5c193e~1 448f5a57 -- internal cmd/balldebug tools frontend design docs` | 只列 `docs/reports/HANDOVER.md`／`docs/reports/pending-and-issues.md`（编排者自己的账），**`internal/**` 与 `cmd/balldebug` 零字节** | ✔ |
| `internal/config/schema.go` 四枚默认值未动？ | `git diff ff5c193e~1 HEAD -- internal/config/schema.go` | **空**；现读 `:179-183`＝`Summon/Mute/Cancel/Panel` 四行，只有 `Cancel` 带 `default:"Esc"` 标 | ✔ |
| `none`/`off` 一族零 normalize（`Q-82` 默认＝不做）？ | `git grep -e '"none"' -e '"off"' HEAD -- cmd/wisp \| grep -v '_test\.go'` | **0 命中**；`git grep -n 'NormalizeHotkey' HEAD -- cmd/wisp` 0 命中 | ✔ |
| 名册越界？ | `git show --stat` 逐笔（`28a2ff4e`／`fc4aedaf`／`33d6cfed`／`448f5a57`） | 只含 `cmd/wisp/**`＋`.scratch/wisp/probes/296/r1/**`；三枚冻结件／golden／`thresholds.go`／`allowlist.txt` 零命中 | ✔ |

⚠ **顺带抓到的与 AC#2 无关的一枚盘上事实（归口给编排者，本腿不动）**：`33d6cfed`（AC#4 门禁那笔）里
`.scratch/wisp/probes/296/r1/logs/build.md` 是 **0 字节件**，而票面 AC#4 逐字写着"⛔ 0 字节的件＝那格没交"。
同一笔里 `logs/build-final.md`（4 行）与 `logs/build-tail.md`（1 行）载有 rc ⇒ **证据不缺**，但"每把门禁件自己落一行 `rc=N`"这一条在这一枚上没兑现。
（AC#4 已翻勾，是否追改归编排者；本腿⛔ 不翻框、⛔ 不补件。）

## 3. 注释那一发：**只说行为＝合格，但两处读数与盘上不符**

**合格面（先给）**：⛔ 没有新造任何能力承诺，⛔ 没有指名一枚不存在的测试当凭据——逐条点名核对：
`cmd/wisp/resident_hotkey_296_windows_test.go`（存在）、`residentBallHotkeyChain258`（存在）、
`ball.HotkeySource` 的前置条件注释（`internal/ball/hotkey_reload.go:48-49` 存在）、
`panelGeometrySource`（`cmd/wisp/panel_resident_windows.go:198` 存在）、
`hotCfg258` 那句 Warn（存在，见第 1 节）、"its consumer - the bridge - performs no merge at all"
（**真**：`internal/ball/ball_windows.go:813-824` `RebindHotkeys` 直接 `unregisterAll`+`registerAll(cfg)`，无任何默认合并；
唯一的 `opts.Hotkeys = DefaultHotkeys()` 在 `ball_windows.go:154` 的 **New() 构造支**，不在重绑路径 ⇒ 那句"不合并"成立）。
被删掉的那句过期承诺（"The bridge keeps the current bindings on an empty answer…"）**整段消失**，
`git show fc4aedaf -- cmd/wisp/resident_windows.go` 逐字可见它换成了一段行为描述。

**不符面 1（具名）**：新注释写 `Check()'s only skip is the DeepEqual diff against what is applied`。
盘上 `internal/ball/hotkey_reload.go:81-95` 的 `Check()`（`:81 func (r *HotkeyReloader) Check()`）有**两条**不重绑的出口：
`:87-88` `	if r.src == nil {` / `		return false, r.Report()`（**这条恰恰就是"保持当前绑定"**）与
`:93` `	if r.exists && reflect.DeepEqual(r.applied, cfg) {`。
⇒ "only" 一词把 `src == nil` 那一支算漏了；同一包里 `Test258BridgeMutationNoSrcKeepsOldBinding`（本腿定向跑到 **PASS**）
钉的正是"nil src ⇒ 保持旧绑定"。尺度上这句仍然**没有为"空答复"许诺任何东西**（票面要判的就是这一条），
但它对 `internal/ball` 的描述**过宽**，日后有人把热加载源置 `nil` 时会读到与盘上相反的行为。

**不符面 2（具名）**：`… would re-register the set and leave summon, mute and panel dead until the next read`。
`until the next read` 把"恢复时机"许诺过头：连续两次都答四格全空时，第二次与 `applied` 逐字相等 ⇒ `DeepEqual` 跳过 ⇒
**不会在下一次读就活回来**，要等一次"答得跟已应用集不同"的读。真机凭据同向：票面现量 `:31` 那行 `live=0` 之后，
机主那台机器上的三枚键在整个进程生命周期里没回来过（编排者 `AC#3` 复跑窗口里 `live=0` 与 `rebound` 两把尺各 0 命中，是因为**修好了**、不是因为自动恢复）。

**次要措辞（具名，不判缺陷）**：末句括注 `(an empty [hotkey] slot means "unset", not "turned off")` 是在讲**配置文件格式的语义**，
而仓里 `internal/config/schema.go:177` 逐字 `// Empty string = binding unset (feature key disabled). ` 与
`internal/ball/hotkey_windows.go:54`（`Empty string = disabled.`）**同时**存在两种说法，且 `none`/`off` 那一问已裁为**契约空白（`Q-82`，默认动作＝不做）**。
⇒ 这句话是 `ApplyHotkeyDefaults` 已写死的保守方向的复述（方向对），但**以格式语义的口吻**写出来，
会被读成"`[hotkey]` 的空串语义已定案"。建议（⛔ 本腿不改产码）降格为函数口吻：
"the merge treats an empty slot as unset (`ApplyHotkeyDefaults`'s documented direction); whether a user can turn a key off is `Q-82`, not built"。

## 4. 判语

**`AC#2` 部分成立**：修法那一发（两支都套、只动 `cmd/wisp`、四枚边界尺全零越界、`hotCfg258` 未被动）**成立**；
注释那一发**没有新造承诺**（点名核对全通过、过期那句确已删除），但**两处对 `internal/ball` 的行为读数写歪**
（`Check()` 的"only skip"漏了 `src == nil`；"until the next read"把恢复条件说轻了）＋一处措辞越过 `Q-82` 的待定案线。
⇒ 三处都**不构成回归、不必回退产码**，属注释精度问题，处置归编排者（要不要再补一笔 `cmd/wisp` 文案修订，由他裁）。
