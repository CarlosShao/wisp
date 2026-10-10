# 303-a1 `14-suspect-classification.md` — 嫌疑面七枚"动的是产码还是测试码"（尺＝`git show --name-only --format= <sha>`，母仓现量）

尺的射程＝该笔**全部**改名文件，我只留 `.go/.sh/.ps1/.ts/.tsx/.js` 那几行（`.md`／台账件略去，⛔ 影响判类）；
逐笔名册可逐字重跑（本件不落行号，只落"这枚文件在不在那一笔里"）。
判类口径：**产码**＝非 `_test.go` 的 `cmd/**`／`internal/**`；**测试码**＝`*_test.go`；**台件／探针存档**＝`.scratch/**`。

| 嫌疑笔 | 产码文件 | 测试码文件 | 类别 |
|---|---|---|---|
| `9995f9b1`（33-r11 改名 `firstRoundTripLocked`→`firstRoundTrip`） | `cmd/wisp/panel_host_windows.go` | `cmd/wisp/panel_locked_naming_33r11_windows_test.go`（新增）·`cmd/wisp/panel_pageover_33r10_windows_test.go` | 产码＋测试码（`.scratch/wisp/probes/33/r11/**` 另有 7 枚存档件与 2 枚 `.sh`） |
| `8b32060b`（255-r1 面板宽度到真窗） | `cmd/wisp/config_readers_255.go` · `cmd/wisp/panel_host_windows.go` · `cmd/wisp/panel_resident_windows.go` | `cmd/wisp/panel_geometry_255_test.go` · `cmd/wisp/panel_reshow_255r1_windows_test.go` | 产码＋测试码 |
| `70b00885`（33-r10 冷启交接补 window-free 尺） | `cmd/wisp/panel_host_windows.go` | `cmd/wisp/panel_pageover_33r10_windows_test.go` | 产码＋测试码 |
| `3a343bc7`（票 35 `AC#8` 夹具面） | — | `cmd/wisp/panel_transport_35r2_test.go` | **⛔ 产码，纯测试码**（与票面 `:15` "⛔ 产码改动"一致） |
| `2fc5f5c9`（票 35 `AC#8` 遗产，编排者代提） | — | `cmd/wisp/panel_transport_35r2_test.go` | **⛔ 产码，纯测试码**（＋`.scratch/wisp/probes/35/r4/scripts/**` 两枚台件） |
| `286a7f30`（35-r2 钩子改"先存原生出口＋再入守卫"） | `cmd/wisp/panel_host_windows.go` | `cmd/wisp/panel_transport_35r2_test.go` | 产码＋测试码 |
| `fb2fb802`（35-r1 `installPanelTransport`／绑定 `wispDispatch`／postMessage 转发 Init） | `cmd/wisp/panel_host_windows.go` | `cmd/wisp/panel_transport_35r1_test.go` | 产码＋测试码 |

⇒ 七枚里 **5 枚动过产码**（`panel_host_windows.go` 出现在其中 **5 枚**里）、**2 枚纯夹具**（`3a343bc7`／`2fc5f5c9`）。
⇒ 只有落在 `panel_host_windows.go`（宿主传输边）那一类，才可能是"门断了"；落在 `*_test.go` 的只可能是"量门的尺断了"。
⚠ 本件⛔ 归因：归因＝`AC#1` 的 bisect 读数，另件。

## 两把嫌疑面尺的复核（⛔ 抄票面枚数）

- 全程尺：`git rev-list --count cc315261..bcd0a543` ＝ **742**（clone 里现量，rc=0）
- 那把 4 枚文件的尺：`git log --oneline cc315261..bcd0a543 -- cmd/wisp/panel_host_windows.go cmd/wisp/panel_resident_windows.go cmd/wisp/panel_pageover_33r10_windows_test.go cmd/wisp/panel_transport_35r2_test.go | cat` ＝ **7 枚**，名单逐枚＝`9995f9b1 8b32060b 70b00885 3a343bc7 2fc5f5c9 286a7f30 fb2fb802`（rc=0）
  ⇒ 与票面 `:15` 的编排者那把（7 枚）**一致**；票面 `:15` 还说"腿那把 4 枚文件＝5 枚"——那是**另一枚文件集**（⛔ 同名册），我没重跑那把、⛔ 需要它。

## 一处更正（记在我名下，⛔ 改已提交的原件）

- `bisect-step.sh`（已提交于 `0edce3b6`）注释里把 PATH 形状的出处写成 `scripts/wisp-cli-tests.sh:113-118`，
  现量＝`grep -n 'PATH is handed over'` → **:101**、`export PATH=` → **:109**（该段是 **:101-109**）。
  那一笔已入库的件我⛔ 改（改了就⛔ 是"跑过的那一版"），更正记在这里。
