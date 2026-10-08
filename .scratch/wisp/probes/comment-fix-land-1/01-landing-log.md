# 01-landing-log — comment-fix-land-1 落地记录（12 处已贴；P39 不贴）

- 腿：comment-fix-land-1；窗口 2026-10-08 19:03–19:2x +08。起手锚＝00-anchor.md（commit `a1b19873`）。
- 漂检：13 目标件对 `9e8477ed`（料基点）与 `0a09a4f2`（复核基点）**零 diff**、工作树零脏 ⇒ **漂过锚的处数＝0**（以为不漂、实测也没漂）；12 处贴前复取锚（`git show HEAD:`＋盘上双读）逐处全等。
- 贴法：一次一处；每贴完立刻重取下一锚再贴；每处贴完跑 `cd cmd/wisp && PATH="$PWD/../../third_party/sherpa-onnx:$PWD/../../third_party/onnxruntime:$PATH" go test . -run 'TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim' -count=1`（act 尺按文件另加跑）。`config_reload.go` 三块按 P09→P08→P07 从后往前贴。
- 逐处账（锚 → 已贴/跳过 → 贴后测试 rc）：

| 处 | 锚（现量全等料文） | 动作 | 贴后测试 rc |
|---|---|---|---|
| P01 | internal/panel/composer_dispatch.go:47-52 | 已贴（6→8 行） | 禁词尺 rc=0；dispatcher 尺 rc=0；roster rc=0 |
| P02 | internal/agent/approval/doc.go:26-37 | 已贴（按复核修正末 bullet：`loop.go:794` 分支仍在、`run.go:1008` AdmitTask 已装订 ⇒ 改为「不再被走到、剩死码清理」；12→12 行） | roster rc=0 |
| P03 | cmd/wisp/panel_inbound.go:11-14 | 已贴（4→5 行） | roster rc=0；AC9 rc=0 |
| P04 | internal/panel/git.go:75-81 | 已贴（7→10 行） | panel-git 尺 rc=0；roster rc=0 |
| P06 | cmd/wisp/config_readers_255.go:42-45 | 已贴（4→5 行） | roster rc=0 |
| P07 | cmd/wisp/config_reload.go:11-13 | 已贴（PLAN.md 2715→2721；3→4 行） | roster rc=0；223 族 rc=0 |
| P08 | cmd/wisp/config_reload.go:14-17 | 已贴（按复核改引 `manager.go:172/:180`；4→6 行） | 同上 |
| P09 | cmd/wisp/config_reload.go:18-21 | 已贴（4→5 行；照料文原样） | 同上 |
| P10 | cmd/wisp/firstrun.go:**7-14**（真身块；料文自称 8-15） | 已贴（真身 8→10 行；:15 分隔 `//` 原位未吞；loader.go 238→252） | roster rc=0；257 族 rc=0 |
| P11 | cmd/wisp/resident_approval_windows.go:279-283 | 已贴（5→6 行） | roster rc=0；265 族 rc=0 |
| P12 | cmd/wisp/approval_reply.go:12-16 | 已贴（5→7 行；引号形状票面逐字行＝ASCII 直引号，现取料文核过） | roster rc=0；201 尺 rc=0 但「no tests to run」（201 族判据不在 cmd/wisp 包，不算跑到） |
| P13 | internal/risk/assessor.go:28-32 | 已贴（5→7 行；C19 横幅内只动描述句） | roster rc=0；risk 整包 rc=0（3.56s） |
| P39 | cmd/wisp/run.go:331-336 | **未贴**（见下） | — |

- **P39 报回**：我没动 P39、原因是照贴 +3 行会把 `cmd/wisp/run.go:991` 的 `[cfg := rt.cfg]` 漂到 :994，而它被 `cmd/wisp/config_readers_255.go:109`（P06 +1 后该名册行自身移到 :110，串内容未动）名册逐字引、由 `TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim` 核行 ⇒ 留编排者裁。
- commit 清单：`a1b19873`（00-anchor）→ `d7a3d62c`（P01-P03）→ `d51dde33`（P04,P06）→ `8dbae00c`（P07-P09）→ `aec801ef`（P10,P11）→ `aa6c1881`（P12,P13）→ 本件一笔。
- 终检读数：`git diff 97cb748a..HEAD`（10 个受改件）**非注释行改动＝0**（逐行 `//` 起）；`cmd/wisp/run.go` 零 diff；受改件工作树 porcelain 空；10 件合计 +65/−46。
- 具名报回（与料文/复核件的偏差，均未临场改，除复核点名两处）：
  1. P02 修正版 bullet 占 3 行 ⇒ 该块 12→12（料文原版为 12→11）。
  2. P02 引 `resident_approval_windows.go:369`：按「P02 先贴、P11 后贴」执行 ⇒ P02 落笔时 :369 为真；P11 +1 行后该引用漂到 :370（复核 §五① 预告过；注释面、无仪器）。未改（超出复核点名范围）。
  3. P09 照料文原样贴（含 AC#2 归名，复核 §四② 记轻差；未获修正令）。
  4. 派单 PATH 里的 `third_party/onnxruntime` 目录不存在（现量 ls rc=2）；DLL 全在 `third_party/sherpa-onnx/` ⇒ 测试实跑、无 0xc0000135。
- 写面声明：本目录 2 枚 `.md`＋10 个产码文件的注释行；`run.go`/票面/台账/冻结件零改动；零 push、零翻框、零 `-done`。
