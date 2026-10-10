# 302-a1 — 起手锚（票 302 `AC#0` 只读代价普查）

腿：`302-a1`（只读普查腿）。⛔ 任何产码。⛔ 任何 go 编译面（`go build`/`go test`/`go list`/`go vet` 一次没跑）。

## 起手（现量，落第 1 笔之前）

- `git log -1 --format='%h %ad' --date=iso-strict`
  => `6414a4bb 2026-10-10T16:10:11+08:00`
- `git status --porcelain -- scripts internal .github docs cmd | wc -l`
  => `0`（⚠ 与派单预告"工作树里躺着别人的脏件"在这五个目录上**不对上**：这四类目录此刻干净。
  脏件在别处 —— `git status --porcelain | wc -l` => `840`，含 `M .gitignore`、`D design/**`、大量 `.scratch/**` 未跟踪件。
  ⇒ **本腿一律按"对象层读＝`git show HEAD:<path>`"办事，不拿工作树当 HEAD**，这条规矩照办；分歧具名报回见 §分歧。）
- `tasklist` 现量：`wisp.exe` 枚数 => `0`；`balldebug.exe` 枚数 => `0`（`No matches found`）。
- 写面确认：仓根**没有** `probes/` 顶层目录参与本腿（`ls -d probes` 存在与否与本腿无关；先例 `301-a1` 落在
  `.scratch/wisp/probes/301/a1/`）⇒ 本腿产物落 `.scratch/wisp/probes/302/a1/`，只新建 `.md`/`.txt`。
- 根 `.gitignore` 第 1–8 行现量：`# --- Go ---` / `*.exe` / `*.exe~` / `*.dll` / `*.so` / `*.dylib` / `*.test` / `*.out`
  ⇒ **第 8 行确为全仓 `*.out`**，本腿证据件一律 `.md`/`.txt`，不用 `.out`。

## 读过的来路件（⛔ 转述，以原文为准）

- 票面：`.scratch/wisp/issues/302-page-reply-instruments-run-in-the-ci-cli-tier-and-can-never-get-a-reply.md`（整枚读完，1–35 行）。
- 起因正文：`.scratch/wisp/probes/301/orch/2026-10-10-ci-color-backfill.md` §4。
- `33-n1` 结论件：`.scratch/wisp/probes/33/n1/verdict.md`（提交于锚 `a5357fa1` 那一轮）。
- CI 原始字节：`.scratch/wisp/probes/301/orch/logs/`（`job-114162576251.log` 914906 B＝改后、
  `job-114069831344-baseline.log` 832029 B＝基线；切片 4 枚 `_cli-block`/`_cli-fail` 已在）。

## 终态锚（交件前现量，与起手同形）

- `date` ⇒ `2026-10-10 16:3x +0800`
- `git log -1 --format='%h %ad' --date=iso-strict` ⇒ 起手 `6414a4bb 2026-10-10T16:10:11+08:00` ／ 终态 `83cd66a8 2026-10-10T16:30:51+08:00`
  ⇒ **HEAD 在本腿期间被同机别的腿前进过**（起手锚＝我全部 blob 行号的读法基准，终态锚＝交件基准）。
- `git status --porcelain -- scripts internal .github docs cmd | wc -l` ⇒ 起手 `0` ／ 终态 `0`
- `tasklist` 现量 ⇒ 起手 `wisp.exe=0`、`balldebug.exe=0` ／ 终态 `wisp.exe=0`、`balldebug.exe=0`（机主全程没动 CPU）
- ⛔ 编译面：`go build`/`go test`/`go list`/`go vet`/`go run` 本腿 **0 次**；`go env` 也 **0 次**（⛔ 例外需自报，因为⛔ 例外）。
- 本腿新增件（只建⛔ 删）：`00-anchor.md`、`40-census.md`、`cm-1.txt`／`cm-2.txt`（commit 文案临时件，⛔ 入库 ⇒ pathspec 里⛔ 列它们）。
