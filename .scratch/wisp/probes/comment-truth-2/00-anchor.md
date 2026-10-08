# comment-truth-2 / 00 — 起手锚

## ① 现量（命令原文与输出，逐字）

- `date '+%Y-%m-%d %H:%M:%S%z'` -> `2026-10-08 16:14:09+0800`
- `git log -1 --format='%H %ad %s'` -> HEAD = `7a452d0aca1bf5157a97fb5329ca6eb38b9f0694`，
  `Thu Oct 8 16:10:31 2026 +0800`，subject 起手 `A721 落账（节头 16:0x＝追加前再跑一次 date＝16:06:41 现量…）`。
  ⇒ 派单里那句「HEAD 约为 `7a452d0a`」**与盘上一致**，无冲突。
- `git status --porcelain | head -20`（rc=0）＝只登记别人的在飞，⛔ 一概不动：
  ` M .gitignore` ／ ` M .scratch/wisp/probes/152/my152.py` ／
  ` M .scratch/wisp/probes/161/r6/logs/flip-1..6.txt`、`flip-baseline.txt`、`flip-restored.txt` ／
  ` M .scratch/wisp/probes/242/r3/logs/probe-routed.txt` ／ ` M .scratch/wisp/probes/268/v1/evidence.md` ／
  ` D design/assets/base.css`、`icons.js`、`theme.js`、`tokens.css` ／ ` M design/doubao/README.md`、
  `demo/app.js`、`demo/index.html`、`demo/styles.css`。
  ⇒ 按禁区：`design/**` 不读；取数一律 `git show HEAD:<path>`。

## ② 工作名册（整族 15 枚，非抽样）

来源：`git show HEAD:.scratch/wisp/probes/stale-claim-1/01-prod-comment-roster.md`（下称「表一」）。
表一 `:5` 逐字：`分档计数（含 P39）：**已过期 15 ／ 判不动 9 ／ 仍成立 15**。`
表一 `:6` 逐字：`（甲节标"14 枚"是 P39 补入前的数；P39 是**半枚过期**，计入已过期一侧，故甲节实际 15 枚。）`
⇒ 本程按**整族 15 枚**做，枚数＝15，不是抽样。

⚠ 一枚**名册内部冲突，具名报回、不自改**：表一标题行 `:1` 逐字 `# stale-claim-1 / 01 — 表一：产码注释名册（HEAD `6547fd30`）`，
而我这一程的 HEAD 是 `7a452d0a`。⛔ 这不是"分类不成立"，是**锚会漂**——所以 15 枚的 ⓐ 锚全部在 HEAD 现取（下表）。
表一 `:22`（P01 段）里那条"同一枚断言的行号在台账/票面里被引成 `composer_dispatch.go:48`…漂 2 行"就是同一颗雷的先例。

| 号 | 表一给的锚（HEAD `6547fd30` 时代） | HEAD `7a452d0a` 现量锚 | 本程档位（复跑后填） |
|---|---|---|---|
| P01 | `internal/panel/composer_dispatch.go:50` | 见 `01-triage.md` | |
| P02 | `internal/agent/approval/doc.go:26` | 见 `01-triage.md` | |
| P03 | `cmd/wisp/panel_inbound.go:11` | 见 `01-triage.md` | |
| P04 | `internal/panel/git.go:76` | 见 `01-triage.md` | |
| P06 | `cmd/wisp/config_readers_255.go:43-45` | 见 `01-triage.md` | |
| P07 | `cmd/wisp/config_reload.go:12` | 见 `01-triage.md` | |
| P08 | `cmd/wisp/config_reload.go:14` | 见 `01-triage.md` | |
| P09 | `cmd/wisp/config_reload.go:18` | 见 `01-triage.md` | |
| P10 | `cmd/wisp/firstrun.go:11` | 见 `01-triage.md` | |
| P11 | `cmd/wisp/resident_approval_windows.go:282` | 见 `01-triage.md` | |
| P12 | `cmd/wisp/approval_reply.go:11-16` | 见 `01-triage.md` | |
| P13 | `internal/risk/assessor.go:28-29` | 见 `01-triage.md` | |
| P14 | `cmd/wisp/resident_task_source_windows.go:13-14` | 见 `01-triage.md` | |
| P15 | `cmd/wisp/resident_windows.go:248` | 见 `01-triage.md` | |
| P39 | `cmd/wisp/run.go:333-334`（半枚过期） | 见 `01-triage.md` | |

## ③ 纪律（本程自钉）

- ⛔ 零 Go 命令（`go build/vet/test/list/doc/env` 一律不跑；此刻写腿 `253-r1` 在 `cmd/wisp` 上跑测试）。
- ⛔ 零翻框、不加 `-done`、不改任何既有文件、不新建判据、不提修法（只分档）。
- ⛔ 不 push；只 `git add -- <显式路径> && git commit -F <消息> -- <显式路径>`。
- 写面＝只新建本目录下的 `.md`；⛔ `.sh/.ps1/.txt/.out`。
- 范围尺不接 `2>/dev/null`，跑完必检 rc；枚数尺写明射程文件。
- `frontend/**` 与 `design/**` 不读不引。

## ④ 已落账部分（避免重复发现）

台账 `docs/reports/pending-and-issues.md` 的 `A720`＝编排者已用**自己的尺**复认了表一里的 **5 枚**。
具体哪 5 枚与它的原话，见 `02-attribution.md` §0（现量摘录，⛔ 不在此转述）。
