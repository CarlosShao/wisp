# 票 296 `AC#3` 的编排者真机跑法（写于 2026-10-09 20:5x，⛔ 现在**不跑**）

**为什么现在不跑**：`296-r1` 正在跑 `cmd/wisp` 整包名册（`logs/prefix-full-*.md` 20:47/20:49 还在长）。★**起 `build\wisp.exe` 会把那把尺搅出假红**（本仓今日实测两次：PID 9084 那枚探针顶出 `TestAC246DevLegIgnoresTheTestTaskInjection`）。⇒ **顺序**：腿交回 → 我复跑门禁并裁 `AC#1`/`AC#2`/`AC#4` → **再**跑本节 → 派 `296-v1`。

## 起手闸门（缺一不许起进程）

1. `tasklist //FI "IMAGENAME eq wisp.exe"` ⇒ **0 枚**；`tasklist //FI "IMAGENAME eq balldebug.exe"` ⇒ **0 枚**。
2. 无其它在飞写腿（`git log -1` 之后盘上不再有新的 `cmd/wisp` 写入）。
3. ⛔ **不改机主那份 config**（票面逐字"⛔ 一个字都不改它"）。

## 前提我 20:5x 已现量（这条是 `AC#1` 夹具与真机同形的凭据）

尺＝按节抽 `[hotkey]` 的**键名＋值长度**（⛔ 不打印值、不碰同目录的 `secrets/`）⇒ 文件 `C:\Users\swq\AppData\Roaming\wisp-dev\config.toml`：
`summon length=2`／`mute length=2`／`cancel length=5`／`panel length=2` ⇒ 即 `""`／`""`／`"Esc"`／`""`。
⇒ 与夹具那四行逐字一致；⚠ 注意布局是 **`wisp-dev`**（dev 档）不是 `wisp`（那一档今天只有 `logs/`，没有 config）。

## 读数判据（三形都要，逐形带命令）

- **ⓐ 改后不再拆键**：`GOFLAGS= go build -o build/wisp.exe ./cmd/wisp` → `./build/wisp.exe`（常驻腿，无终端；日志落 `%APPDATA%\wisp-dev\logs\`）→ 起手应见 `hotkeys_live=3 hotkeys="summon=Ctrl+Alt+Q live, mute=Ctrl+Alt+M live, panel=…"`；**存活 ≥5 s**（＞一个 1 s 轮询拍）后 **`live=0` 那行必须 0 命中**（尺＝对日志文件 `grep -c "hotkeys rebound after config change.*live=0"` ⇒ **0**）。⚠ 只看启动那一行不够——回归恰恰是"起手 3、一秒后 0"。
- **ⓑ 按一次 `Ctrl+Alt+M`** ⇒ stdout/日志应出现票 290 那两句之一（`resident_ball_windows.go:536` 的 `wisp: ball mute-hotkey: …`，开门/关门各一形）。合成输入的试法＝PowerShell `Add-Type` + `keybd_event`（VK 序：`LCONTROL`+`LMENU`+`M`）。⚠ **若 `keybd_event` 打不出这一句**：那不等于"键没修好"——`RegisterHotKey` 对注入键的待遇是平台行为，**这一支我按"仪器造不出的世界"具名归并到机主在场那次合并窗口**（他按一次即可），⛔ 不许拿"合成输入没反应"当"修法无效"，也不许为它放宽任何断言。
- **ⓒ 收口干净**：托盘「退出」优先（那一次顺带把票 247 欠的 `levels_delivered=` 收尾句一并要掉＝合并窗口的第二形）；拿不到就走具名 PID：先 `Get-CimInstance Win32_Process -Filter "Name='wisp.exe'"` 看 `CommandLine` 确认是**我自己**刚起的那枚，再 `Stop-Process -Id <pid>`。⛔ 不许 `taskkill /IM`／不许杀任何不是我起的进程（`292-r1` 那条裁定同样适用）。
- 跑完 `tasklist` 再核一次 ⇒ `wisp.exe` **0 枚**（残留会污染下一枚整包尺）。

## 落账

三形各一行读数（带命令逐字＋命中数）＋ `AC#3` 翻不翻的判语，一起写进票 296 的裁定节与本票 Progress log（⛔ 派单里"归编排者跑"的那格，腿没跑＝不算它的欠账）。
