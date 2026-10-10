# 296-v1 / AC#1 —— "空串世界"夹具到底有没有牙：判语 **成立（有牙，逐支独立敏感）**

HEAD＝`28b484ec`。被裁件＝`cmd/wisp/resident_hotkey_296_windows_test.go`（两枚用例）
＋`cmd/wisp/resident_windows.go` 的包级 `hotkeyReloadSource296`。

## 1. 台件纪律（逐条落盘）

| 步 | 尺／命令逐字 | 读数 |
|---|---|---|
| 备份到仓外 | `cp cmd/wisp/resident_windows.go "$TEMP/wisp-296-v1/resident_windows.go.orig"` | 原件 sha256＝`c8182a684decbd54a5477e99b879249788dc900801d871752e0bbe2fba7cc968`（两份逐字相同，见 `logs-sha-before.md`） |
| 仓内不留 `.go` 拷贝 | 全程只在原地突变，未向 `.scratch/**` 解过任何 `.go` | 取 blob 一律落 `$TEMP/wisp-296-v1/`（仓外） |
| 还原 | `cp "$TEMP/wisp-296-v1/resident_windows.go.orig" cmd/wisp/resident_windows.go` | 还原后 `sha256sum` 两份逐字相同＝`c8182a684decbd54a5477e99b879249788dc900801d871752e0bbe2fba7cc968` |
| 还原证明 | `git diff --stat -- cmd/wisp` ／ `git diff --stat -- '*.go'` | **两把都空**（无任何 `.go` 处于改动状态） |
| 还原后功能证明 | `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp/ -count=1 -run 'Test296\|Test258' -v` | **rc=0**；尺名 `grep -cE '^--- PASS'`＝**13**；`grep -cE '^--- (FAIL\|SKIP)'`＝**0**（顶层尺，不含子测试；未跑整包名册） |

本腿跑过的 go 命令全集（共 5 发，全在 `cmd/wisp` 单包、全带 `-run`，⛔ 没跑整包 `-count=1`、没起任何进程）：
HEAD 定向绿 1 发（rc=0，13 枚顶层 PASS）＋突变 A/B/C 各 1 发（rc=1/1/1）＋还原后复跑 1 发（rc=0）。

## 2. 定向突变（只摘套壳，⛔ 别的一字未动）

突变形状＝把 `hotkeyReloadSource296` 里 `return` 上的 `ball.ApplyHotkeyDefaults(...)` 套壳摘掉，还原成落地前的裸返回：

- `:49` `		return ball.ApplyHotkeyDefaults(ball.HotkeyConfig{})` → `		return ball.HotkeyConfig{}`
- `:53` `	return ball.ApplyHotkeyDefaults(ball.HotkeyConfig{Summon: c.Hotkey.Summon, Mute: c.Hotkey.Mute, Cancel: c.Hotkey.Cancel, Panel: c.Hotkey.Panel})` → 去套壳

三枚突变体，各指名跑 `-run 'Test296' -v`：

| 突变体 | 摘掉哪支套壳 | rc | form-1（机主三空串） | form-2（文件读不到） |
|---|---|---|---|---|
| A | 两支全摘 | **rc=1** | **FAIL** | **FAIL** |
| B | 只摘 `:49`（不可读支） | **rc=1** | PASS | **FAIL** |
| C | 只摘 `:53`（裸映射支） | **rc=1** | **FAIL** | PASS |

⇒ 两枚用例**各自只对自己那一支敏感**，不是"一枚红带走另一枚"，也不是对整件事不敏感的恒红／恒绿仪器。
B／C 两发是"判据换成反形它也绿"那一问的正解：**反形（只修一支）当场暴露**，与票面裁定"甲射程必须含两支"同向。

### 红句逐字（突变体 A，`resident_hotkey_296_windows_test.go:176`／`:217`）

```
    resident_hotkey_296_windows_test.go:176: AC#1 form-1 (owner's three empty slots) RED: after the bridge's first Check live=0 [], want live=3 (summon, mute, panel; missing [summon mute panel]), configured={Summon: Mute: Cancel:Esc Panel:}, per-slot=summon=disabled (unset in config)="" mute=disabled (unset in config)="" cancel=not bound while idle (cancel is borrowed only during Confirming)="Esc" panel=disabled (unset in config)="". The reload source handed the bridge un-merged empty slots and the keys were unregistered - ticket 296.
```

form-2 那一行同样含 `live=0`，且 `configured={Summon: Mute: Cancel: Panel:}` **四格全空**（连 `cancel` 都 disabled）
⇒ 与 form-1 的 `Cancel:Esc` 逐字不同 ⇒ **两枚用例真的走的是两条不同分支**，与真机 `:31` 那行 `live=0` 同形。

## 3. 票面点名的另外三问

- **是不是两枚不同闭包？** 是。`hotCfg := func() ball.HotkeyConfig { return loadHotkey258(t, dir) }`
  与 `hotReload := func() ball.HotkeyConfig { return hotkeyReloadSource296(dir) }` 是**两枚不同实参**（form-2 里分别是 `bootDir`／`reloadDir` 两个 TempDir），
  ⛔ 没有 258 那枚 `startResidentBall(..., src, src)` 复用形状。尺＝`grep -n 'startResidentBall(' cmd/wisp/resident_hotkey_296_windows_test.go` ⇒ 两行都是 `(…, hotCfg, hotReload)`。
- **夹具那四行是不是逐字机主形状？** 是。`ownerHotkeyBody296`＝
  `"[hotkey]\nsummon = \"\"\nmute = \"\"\ncancel = \"Esc\"\npanel = \"\"\n"`
  ⇒ 逐行对机主 `summon = ''`／`mute = ''`／`cancel = 'Esc'`／`panel = ''`，键序与值一致。
  并带**前提钉子**：`if raw := loadHotkey258(t, dir); raw != (ball.HotkeyConfig{Cancel: "Esc"}) { t.Fatalf("premise moved…") }`
  ⇒ 若哪天 schema 在解析期就把槽填了，这枚用例会**自己停下来**而不是悄悄恒绿（票面 258 那格"源可全空"未被放宽）。
- **断言是不是票面要的两件？** 是：ⓐ `live=3`（`missingLive296` 逐名点名 summon/mute/panel）ⓑ `ConfiguredHotkeys()` 等于 `ball.DefaultHotkeys()`
  （ⓑ 用 `t.Errorf` 先说、ⓐ 用 `t.Fatalf` 收口 ⇒ 一发红同时报两半）。
  另有 `hotkeyBridgeArmed` 与 boot `live!=3` 两处 **SKIP-LOUD**（无 Win32 的宿主）⇒ 本机两枚用例都是 `--- PASS`（0.04s／0.04s），未触发跳过。

## 4. ★具名评估那一处偏离（提为包级 `hotkeyReloadSource296`）：**维持编排者的"不算越界"，本腿不推翻**

理由（本腿自己的凭据，不是转述）：
- 产码与测试**共用同一枚本体**＝`git grep -n 'hotkeyReloadSource296' HEAD -- cmd/wisp`：定义 `:43`、产码调用 `resident_windows.go:249`
  `hotReload258 := func() ball.HotkeyConfig { return hotkeyReloadSource296(rt.Layout.DataDir) }`、测试调用 `:156/:201` ⇒ 测试驱动的是**装配根真正跑的那段体**。
- 提取那一笔**行为零变化**：`git show 28a2ff4e -- cmd/wisp/resident_windows.go | grep -E '^[+-]' | grep -c 'ApplyHotkeyDefaults'` ⇒ **0**
  （合并套壳是下一笔 `fc4aedaf` 才加的，那笔我逐字读过 diff）。
- ⛔ 在测试里重打一遍裸映射闭包会造出**恒绿仪器**：本腿的突变台件正是"只改产码本体"，测试随即红——若测试自带一份拷贝，同样的突变**不会红**。
  ⇒ 票面"用生产那两枚闭包的同形"这句在 Go 里不可达（函数体内局部闭包不可从包外引用），**属票面射程写窄**。

⚠ **一枚本腿看到的、编排者那句裁定没覆盖的代价（具名报回，处置归编排者）**：夹具的敏感度**绑在"合并发生在源本体"这一处**。
若日后有人用**等价修法**把合并挪到调用点（`resident_ball_windows.go:353` 的 `NewHotkeyReloader` 第三个实参外面套一层），
产码行为正确、而这枚夹具会**红**（它直接调 `hotkeyReloadSource296`，不经过调用点）。那属**假红**而非回归。
⇒ 建议（不是本腿改）：要么在 296 件里写明"合并必须在源本体，调用点等价修法不采纳"，要么补一枚走 `startResidentBall` 装配根的接线钉
（现成的形状＝同包 `Test258AssemblyRootWiresTheChainAndTheBridge`，本腿读到它 PASS，但它钉的是 258 那半的接线，没钉热加载半的合并位置）。

## 5. 判语

**`AC#1` 成立**：夹具**有牙**——三枚突变体把两枚 `return` 的每一支分别钉住（A 双红／B 只红 form-2／C 只红 form-1），
红句逐字含 `live=0`；两枚闭包不同形、四行空串逐字＝机主形状、并带"前提挪位就停手"的钉子。
落地腿报的"改前两形都红、改后全绿"经本腿**独立复现**（不是采信）。
