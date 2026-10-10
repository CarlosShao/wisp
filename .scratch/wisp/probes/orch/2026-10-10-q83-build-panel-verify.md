# Q-83 甲 编排者亲跑记录（出厂构建链 + 真机面板那一发）· 2026-10-10 08:4x–08:5x

来路＝机主 2026-10-10 08:3x 答「1. 需要」＝批准 `Q-83` 甲（把带页面的版本跑出来）。射程与撤销口令见台账 `A793`（口令「撤 Q-83 甲」）。⛔ 那一句批准不含推送，本轮 ⛔ 零 push。

## 0. 起手闸门（同一条命令里量，缺一不起手）

- `date` ⇒ `Sat Oct 10 08:46:09 CST 2026`
- `tasklist //FI "IMAGENAME eq wisp.exe"` ⇒ `INFO: No tasks are running which match the specified criteria.`（0 枚）
- `tasklist //FI "IMAGENAME eq balldebug.exe"` ⇒ 同上（0 枚）
- `node --version` ⇒ `v24.18.0`（bash 侧）／出厂链自己解析到的是 `D:\work\server\node14\node.exe (v24.9.0)`；`npm --version` ⇒ `11.6.0`
- `Get-CimInstance Win32_Processor | Select -Expand LoadPercentage` ⇒ **48%**（< 70% 才加派）
- ⛔ 未读机主 `config.toml` 的值，⛔ 未读同目录 `secrets/`

## 1. 起手前我先量的一件事：这台机器上的 exe 到底带不带页面字节

尺＝`grep -c -a -o '<产物名>' build/wisp.exe`（对**二进制**取词面命中）＋ PE 子系统字段（`e_lfanew+92` 那两字节，2=GUI／3=CUI）

| 读数（起手时，`build/wisp.exe` mtime 2026-10-09 17:00） | 值 |
|---|---|
| 旧产物名 `index-BVKlegVD.js` 命中 | **2** |
| 旧产物名 `index-BRKj5OIJ.css` 命中 | **2** |
| `id="root"` 命中 | **1** |
| PE 子系统 | **2（GUI）** |
| `build/SHA256SUMS` 里记的 `wisp.exe` 前缀 | `5e4891f4…` |
| 盘上那枚 `wisp.exe` 实际前缀 | `0867f1b6…` ⇒ **名册已过期，那枚 exe 不是名册对应的那一发** |

⇒ 起手时"exe 里零页面字节"这句**不成立**（页面字节在，只是**上一发的**旧产物）。⚠ 但 `frontend/dist` 是 09-27 10:59 那一发，比今天构建它的工具链旧 ⇒ "带不带页面"与"带的是不是这一发的页面"是两件事，`build.ps1` 那道换名闸门管的正是后者。

## 2. 起手前的自我保护（⚠ 这条是我自己加的，不在票面）

`frontend/dist` 在库里只跟踪 `frontend/dist/.gitkeep` 一枚（尺＝`git ls-files frontend/dist` ⇒ **1**），其余三枚是**盘上唯一的一份页面字节**；而 `vite build` 会先清空输出目录 ⇒ 万一 `npm run build` 中途死，那三枚就没了、`go build` 会退化成"只嵌 .gitkeep"、`panel.Assets.Built()` 翻 false、后续所有面板腿一起丢对象。
⇒ **先把整棵 `dist` 按字节备份到仓外**：`C:\Users\swq\AppData\Local\Temp\wisp-dist-backup-2026-10-10`（4 枚文件／600 KB，含 `.gitkeep`；备份时刻 08:4x，逐枚 sha 前缀 `2078af7f…`/`f8ff337b…`/`793e6063…`）。⛔ 未删原件，⛔ 未把备份放进仓。

## 3. 出厂链那一发（`scripts/build.ps1 -Env dev`）

件＝`.scratch/wisp/probes/orch/logs/q83-build/build-stdout.md`（4,248 字节，含 `PSRC`）。**rc=0**。逐字关键行：

- `build.ps1: toolchain: go version go1.27.1 windows/amd64; CC=E:\work\base\msys64\mingw64\bin\gcc.exe (16.2.0)`
- `fetch-deps: cache hit - third_party/sherpa-onnx matches deps.toml (sherpa-onnx 1.13.8)`
- `node_modules EXISTS -> taking the npm-ci-SKIPPED branch`（具名：这一支不核 `package-lock.json`）
- 换名闸门（票 274 `AC#10` 那道）：
  - 改前快照＝`4 file(s) [.gitkeep assets\index-BRKj5OIJ.css assets\index-BVKlegVD.js index.html]`
  - 改后＝`assets\index-B8yINMF1.js=551989`／`assets\index-yy8KMgdf.css=49540`／`index.html=1044`／`.gitkeep=0`
  - ⇒ `provenance: this run's npm added new artifact name(s) [assets\index-B8yINMF1.js assets\index-yy8KMgdf.css] …; gone since snapshot: [assets\index-BRKj5OIJ.css assets\index-BVKlegVD.js]` ⇒ **闸门放行，且它给的就是"这一发的页面"的凭据**
- `go build ok (cgo linked against sherpa-onnx C API)`／`DLLs colocated`／`signed model manifest colocated`／`wrote build\SHA256SUMS`
- 冒烟：`wisp doctor: PASS`（`version=0.0.0-dev commit=29081a13 built=2026-10-10T00:51:10Z`；gcc／sherpa 1.13.8／onnxruntime 1.28.2／deps.toml 三枚 pin／data dir writable 全 PASS）

## 4. 跑后三把尺（票面没写、我自己要的）

| 尺 | 读数 |
|---|---|
| `find frontend/dist -type f \| sort` | `.gitkeep`／`assets/index-B8yINMF1.js`／`assets/index-yy8KMgdf.css`／`index.html`（4 枚） |
| 新 exe 里**新**产物名命中（`grep -c -a -o`） | `index-B8yINMF1.js` ⇒ **2**；`index-yy8KMgdf.css` ⇒ **2** |
| 新 exe 里**旧**产物名命中 | `index-BVKlegVD.js` ⇒ **0**（⛔ 不是"两发混在一枚里"） |
| `id="root"` 命中 | **1** |
| PE 子系统 | **2（GUI）** ⇒ 双击不弹黑框 |
| `(cd build && sha256sum -c SHA256SUMS)` | `wisp.exe: OK`／`onnxruntime.dll: OK`／`sherpa-onnx-c-api.dll: OK`／`sherpa-onnx-cxx-api.dll: OK` ⇒ 名册与产物重新对上（起手时那枚 `0867f1b6 ≠ 5e4891f4` 的过期已消） |
| `git status --short -- frontend cmd internal tools scripts docs AGENTS.md` | **0 枚** ⇒ 出厂链没把任何**被跟踪**的源码/页面源件弄脏（`dist` 三枚本就被 `frontend/.gitignore:12` 忽略） |
| `ls -la build/wisp.exe` | 31,270,806 字节，mtime `Oct 10 08:51` |

⚠ **记我一处预测错**：我起手前判"同一份源码再跑一次会撞那道换名闸门 ⇒ 必红"（依据＝`frontend/src` 最后一次提交是 09-26 21:26、`dist` 是 09-27 10:59）。**实测红没发生**：闸门绿，因为产物名**真的换了**（旧那三枚不是当前工具链产的形状）。⇒ 定式：**"必红"这种预测也要先跑尺再说**，我这次差点因此不去跑（差点把机主批的那一发省掉）。

## 5. 真机那一发：面板到底看不看得见（件 `.scratch/wisp/probes/orch/logs/q83-live/boot.md`）

- 08:52:57 起 `build/wisp.exe`（PATH 只在同一条命令里赋值一次；stdout 重定向进件）⇒ `wisp.exe` **PID 31272**，`CommandLine` 逐字 `"D:\work\workspace\projects plans\Wisp\build\wisp.exe"`（`Get-CimInstance` 现量，是**我自己**起的）
- 起手读数（顺带把票 296 的修法在**出厂那一枚 exe** 上又验了一遍）：`hotkeys_live=3`、`summon/mute/panel=""`＋`cancel=Esc`、`live=0` 那行 **0 命中**
- 08:53:24 合成 `Ctrl+Alt+P`（`keybd_event`，脚本 `C:\Users\swq\AppData\Local\Temp\wisp-q83-panel-tap.ps1`，⛔ 不进仓）⇒ 日志两行：
  - `panel: INBOUND-DISPATCH … method="wispProbeRT" … err=… 不是面板 composer 通路的能力入口`
  - `wisp: panel window is up (panel-hotkey, cold -1.0 ms, hot path 0.0 ms)`
- 窗**真建出来了**：`list_windows` ⇒ `title="Wisp panel"`，`bounds 640x260 @ (286,286)`，`processId 31272`
- ⚠ **屏上是全白**：`get_window_state` 截图＝纯白客户区；无障碍树 22 枚节点里只有 `region Wisp - Web content` 套了七层空 `region`，**零文本、零控件**
- 08:55:40 第二发 `Ctrl+Alt+P` ⇒ 窗**从名册里消失**（toggle 关掉了），⚠ 但**日志零行**（关窗这一支今天不出声，具名记一笔）
- 08:56:42 `Stop-Process -Id 31272 -Force`（先 `Get-CimInstance` 核过 `CommandLine`）⇒ 事后 `tasklist` `wisp.exe`／`balldebug.exe` 各 **0 枚**
- ⛔ 托盘「退出」那一形**没拿到**：`list_windows` 的 `app` 过滤参数在本机不生效（两次都回全量名册），我没去盲点通知区域（机主托盘图标一堆，误点会开出他别的应用）⇒ 票 247 `AC#6` 的 `levels_delivered=` 收尾行**仍欠**，归并到下一次机主在场的合并窗口（与票 293 `AC#5` 肉眼那枚勾、票 296 `AC#3` ⓒ、"设备被占"那一形同一次）

## 6. 结论（这一发买到什么、没买到什么）

1. **`Q-83` 甲在"构建面"这一格已闭合**：出厂链 rc=0，`frontend/dist` 换成了这一发的产物，`build/wisp.exe` 里嵌的就是它（新名 2 命中／旧名 0 命中／`id="root"` 1 命中），`SHA256SUMS` 重新对上，GUI 子系统，`doctor` PASS。⛔ 零 push。
2. **但"用户看得见面板"仍然不成立**，而且**原因不在 dist、不在数据供给**：屏上那层白是**票 33 `AC#13` 早就写死的那一发探测页**——`cmd/wisp/panel_host_windows.go` 的冷启动时序是"先把 embed 入口 `SetHtml` 进控件，随后 `firstRoundTripLocked` 再发一发**自造的探测页**"，最终留在屏上的就是那张探测页。★**我这次拿到的正是它点名的那枚凭据**：页面主动打回来的方法名逐字是 `wispProbeRT`，而这串字符在 `frontend/src/**` 与本轮 `dist` 产物里 **0 命中**（尺＝`git grep -n 'wispProbeRT' HEAD -- frontend` ⇒ 零命中；`grep -c -a -o 'wispProbeRT' frontend/dist/assets/index-B8yINMF1.js` ⇒ **0**），它只在 `cmd/wisp/panel_host_windows.go` 的 `w.Bind("wispProbeRT", …)` 与那枚自造页里存在 ⇒ **屏幕上跑的不是面板页面，是探测页**。
3. ⇒ 排程后果：**机主要"看得见的面板"，下一枚该动的是票 33 `AC#13`，不是再刷一遍 dist**。`AC#13` 那格 ⛔ 不翻（本轮读数恰恰证明它仍不成立），我只把凭据追加进票 33 的 Progress log。
4. ⚠ 顺带一枚观测：面板**关窗那一支零日志**（开窗有 `panel window is up`，关窗没声）——不在本轮射程，具名转给票 33。

## 7. 盘上现场（⛔ 一条都没删）

- 备份：`C:\Users\swq\AppData\Local\Temp\wisp-dist-backup-2026-10-10`（改前那三枚 09-27 产物）
- 合成按键脚本：`C:\Users\swq\AppData\Local\Temp\wisp-q83-panel-tap.ps1`
- 出厂链 stdout：`.scratch/wisp/probes/orch/logs/q83-build/build-stdout.md`
- 真机日志：`.scratch/wisp/probes/orch/logs/q83-live/boot.md`
- ⚠ 起手时就存在的**别人的**未提交现场（本轮 ⛔ 未动、⛔ 未入库、也⛔ 不是我造的）：`design/**` 16 枚 `D` ＋ 15 枚 `M`、`.gitignore` 一枚 `M`（内容＝新增 `.worktrees/` 忽略段）、`.scratch/wisp/probes/152/my152.py` 一枚 `M`、仓根一枚 0 字节文件名就叫 `-` 的旧件（mtime 10-03 10:20）。⇒ 引用"工作树干净不干净"之前先跑 `git status --short -- frontend cmd internal tools scripts docs` 那一把窄尺（本轮＝0 枚）。
