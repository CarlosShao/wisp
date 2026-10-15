# 319 — 面板线程一死，UI **永不复生**；且二次启动只能"登记"**不能把窗口唤到前台**——2026-10-14 机主实报"界面打不开"

- **Status**：**待派**（面＝`cmd/wisp` 面板宿主侧；本波该面被 `133-r1` 占，排队等面空）。
- **来源**：机主 2026-10-14 20:2x 实报 + 现场日志 `%APPDATA%\wisp-dev\logs\wisp-20261014-001.jsonl`（编排者逐字现读）。
- **现象链（逐字）**：
  1. `15:41:45` `panel thread ending`／`why: the library pump returned`／`window_opened: true` —— 面板线程退出（泵返回），进程仍活着。
  2. `17:21:53`–`20:23:11` 多次 `panel gesture could not be handed over`／`gesture: tray-open-panel`／`why: the panel thread took no request (it is starting up, its queue is full, or it already exited)` —— 每次点托盘/球都无处交付。
  3. `17:22:05`–`20:08:48` 机主反复双击 exe（11 次）：`activation requested by second launch`／`handled_here: booked and cleared only`／`bring_to_front: no path from this loop to a window; the ball, when a process has one…` —— 单实例只登记、无唤窗路径 ⇒ 双击"毫无反应"。
- **编排者临时恢复（已执行）**：`20:30`/`20:35` 掐掉僵尸实例 + 重新拉起（并借机重编含 v2 设置页的 exe：`build/wisp.exe` 20:34、`wisp.exe doctor` rc=0、嵌入式产物名 `index-CQdw3HH0.js`＋`index-MtwQJ5iZ.css` 现量在二进制内）。
- **关联**：票 228（球/托盘与干活进程＝拓扑；本票是"活着的那条进程里 UI 也会死且不可复活"）· "ticket 07"（代码自注 bring-to-front 落它名下——**请核**）· 313/145（面板面）。
- **AC（每格都要答"未修码上响不响"）**：
  - [ ] **AC#0 根因**：`15:41:45` 那条"泵返回"为什么返回（关窗？WebView2 崩？）——给可复算的取证路径；答不出就**明写"根因未定、只登记现象"**。
  - [ ] **AC#1 线程退出后必须能重启**：真机重放该形（面板线程退出）⇒ 再点托盘/球**必须重建出窗**；或给出具名裁决"今天不可重建的原因"＋**机主可自己走**的替代恢复路径。
  - [ ] **AC#2 二次启动必须把已有实例的 UI 唤到前台**（或给机主可见的提示：球闪烁/托盘气泡并写字），⛔ 不许只登记。
  - [ ] **AC#3 钉子**：AC#1/AC#2 各配常驻判据（⛔ 不许只靠日志人读）；变异自证（弄坏⇒红）。
  - [ ] **AC#4 门禁＋面核**：gofumpt／d22scan／长度门；⛔ 不碰 `frontend/**`；改动逐枚给 pathspec。
- **撤销口令**：**「撤 319」**。

## Progress log (append-only, newest last)

- [2026-10-14 20:4x +08] 编排者立案（机主实报＋日志两链逐字；临时恢复动作与构建读数一并入档）。
- [2026-10-14 22:05 +08] agent=319-a1 did=只读根因调查（锚 `ccdceff5`；件 `.scratch/wisp/probes/319/a1/census.md`＋`evt-app-1541.ps1`）。**AC#0 定案**：15:41:45 那发是**关窗**——逐字消息序 `0x112`(WM_SYSCOMMAND)→**`0x10`(WM_CLOSE，同毫秒)**→`DestroyWindow`→`0x2`(WM_DESTROY)→`PostQuitMessage(0)`⇒`Run()` 见 `WM_QUIT` 返回⇒`panel_resident_windows.go:319-320` 收尾⇒`:280`/`:373` 永久拒单；双向锚＝`panel_control_windows.go:266-275`＋`:559-573`。**乙(WebView2 崩)反证**：Application 日志 15:20–16:00 零条、Crashpad `reports/` 空、序里起头是程序化 `WM_CLOSE`；**丙反证**：无 `WM_QUERYENDSESSION/ENDSESSION`。**AC#1 判定＝今天零复活路径**（`startResidentPanel` 唯一生产调用 `resident_windows.go:198`；顺带：`RequestDispose` 注释说的"keeps the pump running"**不成立**，dispose 与按 X 同一条死路）。**AC#2 判定＝21:40:38 最新一发是在线程/窗口**都活着**时也不唤 ⇒ 与 AC#1 可分离**；"ticket 07 名下"那枚自注只活在注释里（`singleinstance_windows.go:17`/`:84`），而 `internal/ball/doc.go:24` 的 07＝动画与静态态。 next=编排者定形（甲：手势处惰性重建，面 `cmd/wisp/resident_windows.go`＋`panel_resident_windows.go`；乙：`internal/proc` 激活 hook＋`AllowSetForegroundWindow`）＋AC#3 钉子（N1 headless 持有者重建／N2 `handleActivation`／N3 winlive `SC_CLOSE` 真机重放）。⛔ 本腿零产码、未动台账/AC 框、未碰 `frontend/**`、未扰机主在跑实例。

- [2026-10-15 10:4x +08] **编排者：319-r1 死腿遗产代提**。该腿被平台 ECONNRESET 掐死（吊约 11.2h 才断，**零回执**），但**盘上活已做完**：改 4 枚（`cmd/wisp/resident_windows.go` +217／`panel_resident_windows.go` +31／`internal/proc/boot_windows.go` +92／`singleinstance_windows.go` +43）＋新 3 枚测试件（`panel_holder_319_windows_test.go` 90 行＝N1「死线程拒绝一切手势」／`activation_hook_319_windows_test.go` 119 行＝N2「hook 恰一次且 Reset」＋N2b「无 hook 仍消费并清」／`panel_death_replay_319_winlive_test.go` 74 行＝N3 真窗重放）＋`probes/319/r1/logs/` 三枚 pre 读数（n1 raw／n2 编译红／n3 winlive 未修形）。
  **代提只跑的门**（四门＋两发）：`go build ./...` rc=0；`gofumpt -l`（七枚）**空** rc=0；`go vet -tags winlive ./cmd/wisp/ ./internal/proc/` rc=0；`sh scripts/d22scan.sh` rc=0；长度门 GREEN；**headless 两包** `-run '319' ./cmd/wisp/ ./internal/proc/`＝**ok 2.371s／ok 0.535s**。
  ⚠ **本笔未跑**：AC#1 三态（未修/修后/还原）／AC#2 活体唤窗／N3 真窗／变异两态——**归 `319-v` 或续腿**（编排者按「代提只跑四门」纪律照实登记，⛔ 未宣称已验）。
