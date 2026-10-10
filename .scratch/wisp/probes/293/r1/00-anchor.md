# 票 293 · 写腿 `293-r1` · `00` 起手锚 + 射程声明

agent=293-r1 ｜ 时刻＝`2026-10-10 09:43:32 +0800`（`date` stdout，非手打）

## 射程（本波交付的格）

- 交付格＝**`AC#2` / `AC#3` / `AC#4`**。`AC#0`/`AC#1` 已由 `293-a1` + 编排者裁完（票面 `[x]`）；
  `AC#5`（肉眼那枚勾）**归机主在场的合并真机窗口**，本腿不碰、不翻框。
- 写面（合法名册）＝**`cmd/wisp/**`** ＋ **`.scratch/wisp/probes/293/r1/**`**。
  `internal/ball` 按甲形裁＝**零改动**（本腿名册里它应当 **0 字节 diff**）。
- ⛔ 不碰：`cmd/wisp/panel_host_windows.go`、`internal/panel/**`（票 299 在飞射程）、
  `frontend/**`、`design/**`、`PLAN.md`、`docs/specs/**`、`internal/observe/thresholds.go`、
  golden、`tools/d22scan/allowlist.txt`、三枚冻结件（`internal/panel/tokens_fourway_test.go`·
  `internal/panel/l2_grant_boundary_test.go`·`internal/perm/ticket90_persist_test.go`）。
- ⛔ 默认值三枚一字不动；⛔ D43 表／球侧 `Muted` 态／`Q-81` 不碰；⛔ 裸 `go func(`；⛔ 新协程 0；
  ⛔ 新包级依赖边 0；⛔ 新造第二/第三枚真相源。
- ⛔ 票面任何 `- [ ]` 框不翻（只追加 Progress log 一条）。⛔ 零 push；commit 必带显式 pathspec。

## HEAD（以我自己这一发为准）

```
6ef14788 票 296 · 验收腿 296-v1 交回三格判语（AC#0 成立／AC#1 成立-夹具有牙／AC#2 部分成立）
```
分支＝`dev`。`git status --porcelain -- cmd/wisp internal/ball` 起手＝**空**（干净）。

## 内容锚复跑（本腿 09:43 现跑，⛔ 不引用他人行号快照）

尺一＝`git grep -n 'SetTrayChecks' HEAD -- internal cmd`：
```
HEAD:internal/ball/ball_windows.go:952:// SetTrayChecks updates the mute / pause-wake checkmarks.
HEAD:internal/ball/ball_windows.go:953:func (b *Ball) SetTrayChecks(muted, pausedWake bool) {
```
⇒ **命中两行＝注释 + 定义；生产调用者 0 枚**（与编排者 09:4x、`293-a1`、票面 16:4x 三发一致）。

尺二＝`git grep -n 'trayMuted' HEAD -- internal cmd`：
```
HEAD:cmd/wisp/resident_windows.go:334:  （注释里的那枚名字，不是驱动点）
HEAD:internal/ball/ball_windows.go:128:	trayMuted     bool        （声明）
HEAD:internal/ball/ball_windows.go:674:	sel := showMenu(b.hwnd, b.trayMuted, b.trayPauseWake)  （每次右键现建时读）
HEAD:internal/ball/ball_windows.go:955:	b.trayMuted = muted   （唯一的写点，在 SetTrayChecks 体内）
```
⇒ 写点只在无人调的 setter 体内 ⇒ `trayMuted` 恒零值 ⇒ 勾永远不带。**成立**。

尺三＝后置 setter 一族（甲形要照的形状）：`git grep -n 'attachMuteGate\|currentMuteGate\|muteMux' HEAD -- cmd internal`
⇒ `resident_ball_windows.go:173/:178`（`muteMux sync.Mutex`）、`:469/:472/:476-477`（`attachMuteGate`）、
`:481/:483/:487-488`（`currentMuteGate` 带锁读回）、`:516`（`fn := rb.currentMuteGate()`）；
产码挂载点＝`resident_windows.go:336` `rb.attachMuteGate(raudio.toggleMute)`；
测试面夹具＝`resident_mute_290_windows_test.go:93/:206/:221/:235`。

尺四＝⛔ 禁用的死字段：`git grep -n 'mutedAtBoot' HEAD -- cmd internal` ⇒
`resident_audio_windows.go:101` 声明 + `:256` 一处写，**零读** ⇒ 写读比 1:0 **成立** ⇒ 勾不接它。

## 同批要改写的那段注释（内容锚，⛔ 不认行号）

三行逐字（本腿 09:4x 于工作树量到 `resident_windows.go:333-335`；票面日志写过 `:322-324`；两个号都不作尺）：
```
	// booked), and no mirror of the gate into the tray's own checkmark - the
	// ball-side trayMuted display flag stays undriven here rather than becoming a
	// second authority nobody reconciles.
```
⇒ 前提已变（`OnTrayMute → muteGesture → gate.SetMuted`），本腿**必须**把它改成主从关系的实话，
⛔ 旧句不许与新代码并存；紧邻的挂载点 `rb.attachMuteGate(...)` 那一行 ⛔ 不动。
