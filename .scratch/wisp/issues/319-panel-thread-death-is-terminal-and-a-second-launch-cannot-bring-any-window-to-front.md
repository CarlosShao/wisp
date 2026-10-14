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
