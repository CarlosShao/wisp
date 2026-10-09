# livewin-v1 — 起手锚（非实现者验收腿）

任务：攻 2026-10-09 真机窗口这批凭据够不够格，逐格给票 247（AC#2/AC#4/AC#6 未跑的归属）、
票 290（AC#3 两形）、票 291（AC#2 四形）的判语。不改产码、不翻框、不 push、零 go 命令。

## 起手锚

- 取锚时刻：`2026-10-09 18:12 +0800`（`date` 原文读数 `2026-10-09 18:12 +0800`）
- `git log --oneline -1` = `d432d728 A787＋§4.0ar 落账：真机窗口收口、票 296/297 立票、记我自己两处（探针活着时引整包名册 / 让他点东西没先说会看到什么）`
- 分支：`git branch --show-current` = `dev`
- 工作树非干净（`git status --short` 首 30 行里含 `.gitignore`、`design/**` 多枚 D/M 等他腿在飞的件）⇒ 本腿**只**用显式 pathspec 提交自己 `.scratch/wisp/probes/livewin-v1/` 下的 `.md`。

## 尺 A — 四枚凭据的现量行数（`wc -l`，18:12 同刻）

| 文件 | 行数 |
|---|---|
| `.scratch/wisp/probes/orch/2026-10-09-live-window-findings.md` | 101 |
| `.scratch/wisp/probes/orch/2026-10-09-a1-run1-quiet.raw.md` | 33 |
| `.scratch/wisp/probes/orch/2026-10-09-a1-run2-voice.raw.md` | 33 |
| `.scratch/wisp/probes/orch/2026-10-09-c1c2-resident.raw.md` | 80 |

findings 101 行与派单给的"101 行"一致（零处不符，第一枚对得上）。

## 尺 B — 凭据的入库出处（`git log --oneline -1 -- <file>`，逐枚）

| 文件 | 最后一次动它的 commit |
|---|---|
| `2026-10-09-live-window-findings.md` | `d432d728`（即 HEAD；"作者自己的具名更正"就在这一笔里进的盘） |
| `2026-10-09-a1-run1-quiet.raw.md` | `73db600a` |
| `2026-10-09-a1-run2-voice.raw.md` | `73db600a` |
| `2026-10-09-c1c2-resident.raw.md` | `73db600a` |

- 派单说三枚原始件"同目录，commit `73db600a`"⇒ 三枚全部对得上（`git cat-file -t 73db600a` = `commit`）。
- ⚠ 具名记下这一枚**时序事实**：三枚原始件落在 `73db600a`，而 findings 的**当前版**落在 HEAD `d432d728`，
  即 findings 在原始件入库**之后又被改过**（含那条作者自更正）。本腿引 findings 一律按 HEAD 版行号现跑，
  并另跑一把 `git log --oneline -- <findings>` 看它被改过几笔。

## 尺 C — 判据文件的现量

- `cmd/wisp/resident_audio_247_live_windows_test.go` = `252` 行（派单点名 `:194` 那一行存在性待本腿自证，见 10-rulers）。

## 没跑完的尺 / 资源禁令

- 本腿按派单**禁一切 go 命令**（build/vet/test）。凡"只能靠跑才能裁"的格，本腿不跑，具名写欠并记账到编排者。
