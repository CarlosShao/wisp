# 308-r1 · 00-anchor（落地腿锚件／commit-first 第一笔）

- 腿代号：`308-r1`（写腿，实现票 308 `AC#1` 守卫落点）
- 票：`.scratch/wisp/issues/308-spelling-invariance-test-panics-when-temp-has-no-volume.md`
- 上游：只读普查腿 `308-a1`（已交 `AC#0`，见 `probes/308/a1/90-final.md`）；守卫落点裁语＝编排者本票指令 §1 ⓐ–ⓓ

## 起手三连（`§4` 步骤 1，现量）

- `date`＝`Sun Oct 11 13:0x CST 2026`
- `git rev-parse HEAD`＝`1cd57ec1b01af41e54bdf0e6d99bd7b3530da36a`
- `git status --porcelain -- internal cmd`＝**0 行**（起手 Go 面两目录零脏；台面可信）

## 在场检查（跑读数前必做，`§3`）

- `tasklist //FI "IMAGENAME eq wisp.exe"` 命中＝`0`
- `tasklist //FI "IMAGENAME eq balldebug.exe"` 命中＝`0`
- 名册层级尺（本仓实测顶层⇄子项差 110 枚）：顶层＝`grep -cE '^--- (PASS|FAIL|SKIP)'`；子项＝`grep -cE '^    --- (PASS|FAIL|SKIP)'`；每份名册**必标层级**。

## 写面（本腿只有这些；`§3` 硬约束）

- `internal/risk/pathresolver_anchor_spelling_windows_test.go`（改守卫落点 `:103` 崩点 + ⓓ 两枚注释缺陷，注释⛔ 变可执行码）
- `internal/risk/pathresolver_junction_windows_test.go`（改守卫落点 `:152` 崩点，认内容锚）
- `internal/risk/<新增>_test.go`（本腿新增：共享 helper `driveShare` + 直接喂无卷标字面量的单元用例）
- `.scratch/wisp/probes/308/r1/**`（读数＋交付件，含 `logs/`，只建不删，⛔ `.out` 名）
- 票 308 末尾 **1 条** append-only bullet（`git diff --numstat`＝1 加 0 删）

⛔ 动产码：`pathresolver.go`／`pathresolver_windows.go`／任何非 `_test.go`（编排者裁语 ⓐ：产码 `:262` 的 `drive[:1]` 前一行逐字 `if len(drive) == 2 && drive[1] == '$' && …` 已有守卫，⛔ 在射程）。

## ⛔ 动的在飞面（共享工作树）

- ⛔ 动 `internal/tools/**`——另一枚写腿 `309-r1` 正在飞的写面。若 `go build ./...` 报在 `internal/tools/**`：具名记那发红句、稍后重试一发，⛔ 改它⛔ 修它。
- ⛔ 碰 `.gitignore`／`go.mod` 未提交改动（⛔ 我名下）；⛔ stage `docs/evidence/s1/**` 与别人 `probes/**` 未跟踪件。
- commit 一律带显式 pathspec；⛔ `git add -A`／`.`；⛔ `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`／`--no-verify`；⛔ push。

## 崩点行号（HEAD `1cd57ec` 现量，内容锚；票面行号是 11:1x 快照，会漂）

- `pathresolver_anchor_spelling_windows_test.go`：`vol := filepath.VolumeName(long)` → 下一行 `share := strings.ToLower(vol[:1]) + "$"`（在 `TestClassifySpellingInvariance` 体内，现量 `:102`/`:103`）
- `pathresolver_junction_windows_test.go`：`vol := filepath.VolumeName(long) // e.g. C:` → `share := strings.ToLower(vol[:1]) + "$"`（在 `TestPathResolverUNCAListDenied` 体内，现量 `:151`/`:152`）
- `tmpHome(t)`＝`pathresolver_junction_windows_test.go:22-35`，体内 `home := t.TempDir()` ⇒ `long` 卷标全由 temp env 决定（`§1` ⓑ：本机造不出无卷标台面，守卫牙挂输入层）。
