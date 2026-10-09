# 296-a1 — 锚与尺的原始读数（只读普查腿，零产码）

## 锚

- `date` 的 stdout 逐字：`Fri Oct  9 18:11:52 CST 2026` ⇒ 本腿落笔钟点 `2026-10-09 18:11:52 +08`
- `git log --oneline -1` 逐字：`d432d728 A787＋§4.0ar 落账：真机窗口收口、票 296/297 立票、记我自己两处（探针活着时引整包名册 / 让他点东西没先说会看到什么）`
- `git branch --show-current`：`dev`
- 本腿身份：`296-a1`，只做票 296 的 AC#0（射程普查）与 AC#2（甲乙两形代价表 + `none` 那一问）。
  ⛔ 零产码改动、零 `go build`/`go vet`/`go test`、零翻框、零 push。

## 尺 R1 —— 宿主侧 `ball.HotkeyConfig{` 复合字面量全名册

命令逐字：

```
grep -rn "ball\.HotkeyConfig{" --include=*.go cmd internal | grep -v _test
```

原始读数（6 枚，逐字）：

```
cmd/balldebug/main.go:239:			return ball.ApplyHotkeyDefaults(ball.HotkeyConfig{
cmd/wisp/resident_windows.go:188:			return ball.HotkeyConfig{}
cmd/wisp/resident_windows.go:195:		return ball.HotkeyConfig{Summon: h.Summon, Mute: h.Mute, Cancel: h.Cancel, Panel: h.Panel}
cmd/wisp/resident_windows.go:213:			return ball.HotkeyConfig{}
cmd/wisp/resident_windows.go:215:		return ball.HotkeyConfig{Summon: c.Hotkey.Summon, Mute: c.Hotkey.Mute, Cancel: c.Hotkey.Cancel, Panel: c.Hotkey.Panel}
internal/ball/hotkey_reload.go:18://		return ball.ApplyHotkeyDefaults(ball.HotkeyConfig{
```

## 尺 R1b —— 同一条尺在 `.scratch/**` 上的读数（本件据此**排除**变异拷贝）

命令逐字：

```
grep -rn "ball\.HotkeyConfig{" --include=*.go .scratch | grep -v _test
```

原始读数＝**8 枚**，逐字：

```
.scratch/wisp/probes/258/v1/mut0/cmd/wisp/resident_windows.go:182:			return ball.HotkeyConfig{}
.scratch/wisp/probes/258/v1/mut0/cmd/wisp/resident_windows.go:188:		return ball.HotkeyConfig{Summon: h.Summon, Mute: h.Mute, Cancel: h.Cancel, Panel: h.Panel}
.scratch/wisp/probes/258/v1/mut0/cmd/wisp/resident_windows.go:206:			return ball.HotkeyConfig{}
.scratch/wisp/probes/258/v1/mut0/cmd/wisp/resident_windows.go:208:		return ball.HotkeyConfig{Summon: c.Hotkey.Summon, Mute: c.Hotkey.Mute, Cancel: c.Hotkey.Cancel, Panel: c.Hotkey.Panel}
.scratch/wisp/probes/258/v1/mut3/cmd/wisp/resident_windows.go:183:			return ball.HotkeyConfig{}
.scratch/wisp/probes/258/v1/mut3/cmd/wisp/resident_windows.go:189:		return ball.HotkeyConfig{Summon: h.Summon, Mute: h.Mute, Cancel: h.Cancel, Panel: h.Panel}
.scratch/wisp/probes/258/v1/mut3/cmd/wisp/resident_windows.go:207:			return ball.HotkeyConfig{}
.scratch/wisp/probes/258/v1/mut3/cmd/wisp/resident_windows.go:209:		return ball.HotkeyConfig{Summon: c.Hotkey.Summon, Mute: c.Hotkey.Mute, Cancel: c.Hotkey.Cancel, Panel: c.Hotkey.Panel}
```

⇒ **排除声明：本件名册排除 `.scratch/wisp/probes/**` 下的变异拷贝 8 枚**（两枚变异体 `258/v1/mut0`、`258/v1/mut3`，各 4 枚，均为 `cmd/wisp/resident_windows.go` 的整文件拷贝，行号相对生产件偏移 −7／−6）。
它们不参与生产路径，也不计入 AC#0 枚数。R1 的尺本身只跑 `cmd internal` 两棵树，所以这 8 枚是"我用第二条尺确认它们存在、然后显式剔除"，不是"尺没看见"。

## 尺 R2 —— `ApplyHotkeyDefaults` 的命中分布

命令逐字：

```
grep -rn "ApplyHotkeyDefaults" --include=*.go cmd internal | grep -v _test
```

原始读数（13 枚，逐字）：

```
cmd/balldebug/main.go:239:			return ball.ApplyHotkeyDefaults(ball.HotkeyConfig{
cmd/wisp/resident_ball_windows.go:207:// maps the host's [hotkey] view through ball.ApplyHotkeyDefaults and names the
cmd/wisp/resident_ball_windows.go:214:// (258-a1 census Q1a), and ApplyHotkeyDefaults merges exactly
cmd/wisp/resident_ball_windows.go:235:	cfg := ball.ApplyHotkeyDefaults(hotCfg())
cmd/wisp/resident_windows.go:193:			"empty_slots_note", "empty slots are filled from the product defaults by ball.ApplyHotkeyDefaults; "+
internal/ball/hotkey_reload.go:18://		return ball.ApplyHotkeyDefaults(ball.HotkeyConfig{
internal/ball/hotkey_reload.go:49:// of [hotkey], already defaulted via ApplyHotkeyDefaults).
internal/ball/hotkey_reload.go:70:// NewHotkeyReloader wires a ball to a hotkey source. The applied set starts
internal/ball/hotkey_reload.go:73:func NewHotkeyReloader(binder HotkeyBinder, current HotkeyConfig, src HotkeySource) *HotkeyReloader {
internal/ball/hotkey_windows.go:78:// ApplyHotkeyDefaults fills the empty fields of a config-sourced binding set
internal/ball/hotkey_windows.go:93:func ApplyHotkeyDefaults(cfg HotkeyConfig) HotkeyConfig {
internal/ball/hotkey_windows.go:546:// spelling, and ConfiguredHotkeys() keeps the whole configured set.
internal/ball/hotkey_windows.go:575://     direction ApplyHotkeyDefaults documents for this slot ("a cancel key that
```

⇒ **`cmd/wisp` 里 `ApplyHotkeyDefaults` 的"真调用"只有 1 处＝`resident_ball_windows.go:235`**（:207/:214 是注释、`resident_windows.go:193` 是日志文案字段）。
`hotReload258` 所经的 `resident_ball_windows.go:353` 那一行**不在命中里** ⇒ 编排者现量"链路里没有任何 `ApplyHotkeyDefaults`"复核**成立**。

## 尺 R3 —— 包内不带 `ball.` 前缀的 `HotkeyConfig{` 字面量（防漏）

命令逐字：

```
grep -rn "HotkeyConfig{" --include=*.go internal/ball | grep -v _test | grep -v "ball\.HotkeyConfig{"
```

原始读数（2 枚，逐字）：

```
internal/ball/ball_windows.go:153:	if opts.Hotkeys == (HotkeyConfig{}) {
internal/ball/hotkey_windows.go:75:	return HotkeyConfig{Summon: "Ctrl+Alt+Q", Mute: "Ctrl+Alt+M", Cancel: "Esc", Panel: "Ctrl+Alt+P"}
```

## 尺 R4 —— 消费点名册（谁把 config 交给 ball）

命令逐字：

```
grep -rn "NewHotkeyReloader\|ApplyHotkeyDefaults\|ConfiguredHotkeys(" --include=*.go cmd internal | grep -v _test
```

原始读数中"宿主消费点"4 枚（其余为 ball 侧定义/注释，见 R2）：

```
cmd/balldebug/main.go:237:		bridge := ball.NewHotkeyReloader(b, b.ConfiguredHotkeys(), func() ball.HotkeyConfig {
cmd/wisp/resident_ball_windows.go:235:	cfg := ball.ApplyHotkeyDefaults(hotCfg())
cmd/wisp/resident_ball_windows.go:353:		bridge := ball.NewHotkeyReloader(b, b.ConfiguredHotkeys(), hotReload)
internal/ball/ball_windows.go:840:func (b *Ball) ConfiguredHotkeys() HotkeyConfig {
```

（逐字全文见 `raw-scratch-copies.txt`，本目录内。）
