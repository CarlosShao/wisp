# 308-a1 · 00-anchor（本腿锚件／commit-first 第一笔）

- 腿代号：`308-a1`（只读普查腿；编排者＝用户本人）
- 票：`.scratch/wisp/issues/308-spelling-invariance-test-panics-when-temp-has-no-volume.md`
- 本腿唯一要答的格：`AC#0`（"哪一发真能崩"现量清楚，四把台面各一发）
- 起手时刻：2026-10-11（本程）

## 锚点与门禁自证（起手现量）

- 锚笔短号（`git rev-parse --short HEAD`）＝`02053534`
- 分支＝`dev`
- `git status --porcelain -- internal cmd` 枚数＝`0`（起手：Go 面两目录零脏行，本腿⛔ 碰它们）
- 共享索引整仓脏行数＝`869`（含 `.gitignore`／`design/**`／别人 `probes/**` 等⛔ 我名的改动；
  本腿⛔ 读过其内容以外的一行、⛔ 动、⛔ 提交；commit 一律带显式 pathspec `.scratch/wisp/probes/308/a1`）

## 本腿调用的门禁面（划界）

- ⛔ 调用任何 `go build`／`go vet`（另一枚验收腿 `298-v1` 正在独占 Go 门禁面）。
- ⛔ 跑 `scripts/d22scan.sh`／`scripts/check-path-length-budget.sh`（那是别腿的门禁面）。
- 只调用 `go test`（票 308 AC#0 判据要求的颜色尺）；⛔ 打 `cmd/wisp` 包（那面一次只给一枚腿）。
- 台面在场探针＝仓外 `os.MkdirTemp` 独立件（放 `C:/Users/swq/tmp/w308a1/probe/`，
  用 `go test` 跑它；⛔ 建在仓内，以免挪动 `gofmt`／`d22scan` 分母）。

## 写面（本腿只有新建）

- `.scratch/wisp/probes/308/a1/00-anchor.md`（本件）
- `.scratch/wisp/probes/308/a1/10-surface-readings.md`（四发台面读数）
- `.scratch/wisp/probes/308/a1/20-comment-defects.md`（Ⓐ 两枚注释缺陷复尺）
- `.scratch/wisp/probes/308/a1/30-codecensus.md`（Ⓑ 产码⛔ 在射程的名册表）
- `.scratch/wisp/probes/308/a1/90-final.md`（终局：四发读数表＋结论＋ⒶⒷⒸ＋欠账＋逐笔 commit）
- `.scratch/wisp/probes/308/a1/logs/**`（原始 go test 输出／探针输出，只建不删；⛔ 用 `.out` 名）

## 读过的三件（按序）

1. 票面 308（射程／六格／边界／规矩）
2. 崩点两枚测试文件（`pathresolver_anchor_spelling_windows_test.go`、`pathresolver_junction_windows_test.go`）
3. 台账最新一笔 `docs/reports/pending-and-issues.md` 的 `## A843`（§2 两枚注释缺陷）

## 崩点行号（HEAD 02053534 现量，非票面 11:1x 快照）

- `pathresolver_anchor_spelling_windows_test.go:102 vol := filepath.VolumeName(long)` →
  `:103 share := strings.ToLower(vol[:1]) + "$"`（在 `TestClassifySpellingInvariance` 体内）
- `pathresolver_junction_windows_test.go:151 vol := filepath.VolumeName(long) // e.g. C:` →
  `:152 share := strings.ToLower(vol[:1]) + "$"`（在 `TestPathResolverUNCAListDenied` 体内）
- 两枚共用 `tmpHome(t)`（`pathresolver_junction_windows_test.go:22-25`），体内 `home := t.TempDir()`
  ＋ `t.Setenv("USERPROFILE", home)` ⇒ `long` 的卷标完全由 temp env 决定。
- 票面写的 `anchor:77/:78` 与当前文件的 `:102/:103` 是同一处（行号漂移，票面自注"行号一律当 11:1x 快照"）。
