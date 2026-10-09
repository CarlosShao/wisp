# 票 296 AC#1 — 夹具先行：两形**改前逐字红句**（产码未改）

leg = 296-r1 · 尺 = `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp/ -count=1 -run 'Test296BridgeKeeps' -v`
起手闸门（同一次命令第一段）：

```
INFO: No tasks are running which match the specified criteria.   (wisp.exe = 0)
INFO: No tasks are running which match the specified criteria.   (balldebug.exe = 0)
Fri Oct  9 20:22:04 CST 2026
```

rc=1（两枚用例全红 = 本格的完成判据）。原始读数全文＝同目录 `red-before-fix.md`（逐字 stdout，未截断）。

---

## 1. 这一笔（commit 2）到底改了什么产码

`hotReload258` 的**本体**被提为包级函数 `hotkeyReloadSource296(dataDir string)`，产码那枚闭包改成调用它，
测试也调用它。**行为零变化**（`GOFLAGS= go build ./...` rc=0；本节的红读数就是提取之后、合并之前的状态）。

为什么要提：票面 AC#1 要「生产那两枚闭包的同形」。测试碰不到 `runResident` 的局部闭包本体，
若在测试里重打一遍裸映射，则产码修没修这枚用例都不敏感 ⇒ 它给不出「改前红」。
同仓先例＝票 258 把构造那一半提为 `residentBallHotkeyChain258`。
★这一条已在 `00-anchor.md` §3 具名上报，判据本身没改（仍是甲-全：两枚 return 过合并）。

`hotCfg258`（构造那一半）**一字未动**，测试的构造闭包用 258 既有的 `loadHotkey258`（同形）；
两半在测试里是**两枚不同闭包**，⛔ 不是 `src, src`。

## 2. 形一：机主那四行逐字（`summon=''`／`mute=''`／`cancel='Esc'`／`panel=''`）

夹具本体（`cmd/wisp/resident_hotkey_296_windows_test.go`）：
`const ownerHotkeyBody296 = "[hotkey]\nsummon = \"\"\nmute = \"\"\ncancel = \"Esc\"\npanel = \"\"\n"`
＋一枚前提自证：`loadHotkey258` 读回来的必须是 `{Cancel: Esc}`，否则夹具不再是这一形。

桥第一次 `Check()` 之后，进程自己打的三行与真机那三行逐字同形：

```
level=INFO msg="hotkey disabled (unset in config)" hotkey=summon
level=INFO msg="hotkey disabled (unset in config)" hotkey=mute
level=INFO msg="hotkey disabled (unset in config)" hotkey=panel
level=INFO msg="ball: hotkeys rebound after config change" summon="" mute="" cancel=Esc panel="" live=0
```

（对照真机 `.scratch/wisp/probes/orch/2026-10-09-c1c2-resident.raw.md:27/:28/:30/:31` —— `live=3` 之后 1.003 秒的同一组读数）

**逐字红句（两行，ⓐ 与 ⓑ 都报了）**：

```
    resident_hotkey_296_windows_test.go:176: AC#1 form-1 (owner's three empty slots) RED: ConfiguredHotkeys() = {Summon: Mute: Cancel:Esc Panel:}, want the merged set {Summon:Ctrl+Alt+Q Mute:Ctrl+Alt+M Cancel:Esc Panel:Ctrl+Alt+P} - the bridge rebound onto the raw empty slots
    resident_hotkey_296_windows_test.go:176: AC#1 form-1 (owner's three empty slots) RED: after the bridge's first Check live=0 [], want live=3 (summon, mute, panel; missing [summon mute panel]), configured={Summon: Mute: Cancel:Esc Panel:}, per-slot=summon=disabled (unset in config)="" mute=disabled (unset in config)="" cancel=not bound while idle (cancel is borrowed only during Confirming)="Esc" panel=disabled (unset in config)="". The reload source handed the bridge un-merged empty slots and the keys were unregistered - ticket 296.
--- FAIL: Test296BridgeKeepsThreeLiveKeysWithOwnerEmptySlots (0.30s)
```

## 3. 形二：配置文件读不到那一支（裁定节新增的射程）

夹具：构造那一半读**有机主四行**的目录，热加载那一半读**没有 config.toml 的目录**；
前提自证＝`config.LoadFile` 对那个目录**确实报错**（`err == nil || c != nil` ⇒ `premise moved`），
所以这一形走的是 `return ball.HotkeyConfig{}` 那一条分支，不是「文件在、节空」。

**逐字红句**：

```
    resident_hotkey_296_windows_test.go:217: AC#1 form-2 (config.toml unreadable at tick time) RED: ConfiguredHotkeys() = {Summon: Mute: Cancel: Panel:}, want the merged set {Summon:Ctrl+Alt+Q Mute:Ctrl+Alt+M Cancel:Esc Panel:Ctrl+Alt+P} - the bridge rebound onto the raw empty slots
    resident_hotkey_296_windows_test.go:217: AC#1 form-2 (config.toml unreadable at tick time) RED: after the bridge's first Check live=0 [], want live=3 (summon, mute, panel; missing [summon mute panel]), configured={Summon: Mute: Cancel: Panel:}, per-slot=summon=disabled (unset in config)="" mute=disabled (unset in config)="" cancel=disabled (unset in config)="" panel=disabled (unset in config)="". The reload source handed the bridge un-merged empty slots and the keys were unregistered - ticket 296.
--- FAIL: Test296BridgeKeepsThreeLiveKeysWhenConfigUnreadable (0.04s)
```

⇒ 形二同样 `live=0`，且它杀的是**四格全空**那一支（连 `cancel` 都成 disabled）——
这一形只套 `:215` 那枚裸映射是修不到的，裁定节的「两枚 return 都要过合并」被夹具钉住。

## 4. 下一笔（commit 3）该长什么样

`hotkeyReloadSource296` 的**两枚** `return` 都包 `ball.ApplyHotkeyDefaults(...)` ⇒ 两形都转绿；
同批把 `:208-212` 那句「the bridge keeps the current bindings on an empty answer」改成按行为写的实话
（桥在全空答案上**不保持**当前绑定，它真去重绑；`hotkey_reload.go` 里唯一的跳过条件是 `reflect.DeepEqual`）。
⛔ 不动 `internal/ball`（乙案未裁）、⛔ 不动 `internal/config/schema.go`、⛔ 不动机主的 `config.toml`、
⛔ 不加 `none`/`off` 的 normalize（`Q-82`，默认动作＝不做）。
