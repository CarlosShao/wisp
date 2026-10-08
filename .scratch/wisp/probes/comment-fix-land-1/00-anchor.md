# 00-anchor — comment-fix-land-1 起手锚（13 处逐枚现取）

- 腿：comment-fix-land-1（窄写腿，只照料粘：把 01-ready-to-apply.md 的 13 处替换文本贴进代码）。
- 起手时刻：2026-10-08 19:03 +08；起手 HEAD `97cb748a`；写本件时 HEAD `a6ec48d3`（两量之间落了 next-instruments-brief-1 的 brief 提交，只动 .scratch/*.md，13 目标件零命中）。
- 料：`.scratch/wisp/probes/comment-fix-prep-1/01-ready-to-apply.md`（基点 `9e8477ed`）；复核：`.scratch/wisp/probes/comment-fix-check-1/01-check.md`（基点 `0a09a4f2`）。
- 同步漂检（现跑，rc=0）：`git diff --name-only 9e8477ed..HEAD -- <13目标件+被引件>`＝空；`0a09a4f2..HEAD -- 同`＝空；`git status --porcelain -- 同`＝空 ⇒ **漂过锚的处数＝0**。
- 逐处现量（`git show HEAD:<file> | sed -n 'N,Mp'` 与料文逐枚对照）：

| 处 | 文件:块（现量） | 块首行现量（截断示） | 与料文 |
|---|---|---|---|
| P01 | internal/panel/composer_dispatch.go:47-52 | `// WHAT THIS DOES NOT DO, stated so nobody infers` | 逐字同 |
| P02 | internal/agent/approval/doc.go:26-37 | `// Wiring still owed by ticket 12 (this package` | 逐字同（末 bullet 按复核修正后再贴） |
| P03 | cmd/wisp/panel_inbound.go:11-14 | `// line and hands the bytes to` | 逐字同 |
| P04 | internal/panel/git.go:75-81 | `// on (the WebView2 host of ticket 33, its` | 逐字同 |
| P06 | cmd/wisp/config_readers_255.go:42-45 | `// landed the tier table in internal/config/tiers.go` | 逐字同 |
| P07 | cmd/wisp/config_reload.go:11-13 | `//   - Manager.CheckAndReload - its only non-test caller` | 逐字同 |
| P08 | cmd/wisp/config_reload.go:14-17 | `//   - Manager.ConfirmLocked - zero assignments outside` | 逐字同 |
| P09 | cmd/wisp/config_reload.go:18-21 | `//   - Manager.OnRestartPending - zero assignments` | 逐字同 |
| P10 | cmd/wisp/firstrun.go:**7-14**（料文自称 8-15） | `// (NewDefaults, defaults.go:58 - the` | 真身 7-14（复核点名；:15 是分隔 `//`，不得吞） |
| P11 | cmd/wisp/resident_approval_windows.go:279-283 | `// One holder, one bind site,` | 逐字同 |
| P12 | cmd/wisp/approval_reply.go:12-16 | `// Gate.DecideFromNative / Gate.DecideFromPanel / Gate.Veto,` | 逐字同 |
| P13 | internal/risk/assessor.go:28-32 | `// Absent dependencies (PathCanonicalizer` | 逐字同 |
| P39 | cmd/wisp/run.go:331-336 | `// modeWrites is ticket 114 AC#2's gate:`（行前 1 个 tab） | 逐字同 ⇒ **不贴**（理由见下） |

- **不贴 P39**（派单硬纪律 #4）：照贴 +3 行会把 `cmd/wisp/run.go:991` 的 `[cfg := rt.cfg]` 漂到 :994，而它被 `cmd/wisp/config_readers_255.go:109` 名册逐字引、由 `TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim` 核行 ⇒ 留编排者裁。
- **修正版两处**（按复核 + 派单）：P10 按真身块 7-14 贴（照自称 8-15 会吞 :15 分隔行并在 :7 留重复行）；P02 末 bullet「which currently kills fs.write」已过期（`cmd/wisp/run.go:1008` 的 `AdmitTask: rt.admitTask` 已装订 ⇒ `internal/agent/loop.go:794` 分支不可达），改写为现行事实后再贴，不许原样贴。
- 落地纪律（派单）：一次一处；每贴完重取下一锚再贴；每处贴完跑 `cd cmd/wisp && PATH=… go test . -run 'TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim' -count=1`（不注 PATH 会 `0xc0000135`）；P01 贴完另跑 `internal/panel` 五禁词尺；`config_reload.go` 三块按 P09→P08→P07 从后往前贴。
- 写面声明：本腿只新建本目录 `*.md`，只改 13 目标件里被点名的那几处注释行；不动票面/台账/冻结件；零 push；(P39) 不动。
