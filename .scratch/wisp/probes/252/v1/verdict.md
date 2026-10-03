# 票 252 · 非实现者验收腿 v1（二程，`252-v1`）—— AC#1／AC#6 结清 + 五发复认 + 恒真两问 + 调用者名册

## §0 起手锚（本腿自取，全部现跑）

- `date` = `Sat Oct  3 15:55:13 CST 2026`（起手）；收尾另有 `date` 复取，见 §9。
- `git log -1 --format=%H` = `6c7a12b656188405e7ab40d6bb7c47066ff55dd6`（分支 `dev`）。
- 起手 `git status --porcelain internal/tools internal/risk` = **0 行**（逐字空输出）。
- 起手包级名册（`go test -count=1 -v`，本腿自跑，日志 `start-tools-test.log`／`start-risk-test.log`）：
  - `./internal/tools/`：**262 PASS／0 FAIL／0 SKIP**（含子测试行计数，`ok … 18.775s`）；
  - `./internal/risk/`：**198 PASS／0 FAIL／0 SKIP**（`ok … 4.888s`）。
  - 名册文件 `roster-start-tools.txt`（262 行）／`roster-start-risk.txt`（198 行），收尾逐名 `comm` 比对见 §9。
- ⛔ 未跑全仓 `go test`／`go build ./...`／`go vet ./...`（派单禁令）；`cmd/wisp` 一枚没跑。
- 前人件已读：票面全文 + `.scratch/wisp/probes/252/r2/impl.md`（284 行）+ 其 §3 突变件目录盘点。
- 生产码锚点（本腿现跑 `grep -n`＋`md5sum`，`paths.go` 真身＝`find cmd internal tools scripts -name paths.go`
  唯一命中 `internal/tools/paths.go`）：
  - `internal/tools/paths.go` md5 = `7a86da7aeba8420639491cd252a47a0c` ＝ **与 `252-r2` §3.1 留档的 md5 逐字节相等**
    ⇒ 自 `70a935ce` 起该文件零漂移（`git log --oneline 70a935ce..HEAD -- internal/tools internal/risk` = 空，复证）。
  - `InAllowlist` 函数体：**起 `paths.go:189`，讫 `paths.go:219`**（`func (p *PathCanonicalizer) InAllowlist(canonical string) bool {` 在 `:189`）。
  - 词法腿（LEG 1）：`:196` `if !rootsContain(p.roots, f) {`（`f := foldPath(canonical)` 在 `:190`），拒绝支 `:197-198`。
  - `resolvedForm` 腿（LEG 2）：`:208` `rf, ok := resolvedForm(canonical)` ＋ `:209` `if !ok || !rootsContain(p.roots, foldPath(rf)) {`，拒绝支 `:210-211`。
  - 工作区腿（LEG 3，不在本硬禁分母）：`:215` `if p.workspace != "" && !rootsContain([]string{p.workspace}, f) {`。
  - 同形修法（`2b1a3071`）落点：`Canonicalize` 内 `:145-148`（`if res.Resolved { return res.Canonical }; return sameFormOfUnresolved(res.Canonical)`），
    `sameFormOfUnresolved` 函数体 `:169-184`。

## §1 AC#1 探针读数（先失配的是哪一腿）

（探针＝`TestTicket252V1NonImplLegsAndReach`，经 `go test -overlay` 挂载，树内零新文件；读数待填本节。）

## §2 AC#6 本机可达性（撞／不撞）

（待填。）

## §3 五发突变复认（重跑的两发：红句逐字＋还原证明）

（待填。）

## §4 恒真两问的答案

（待填。）

## §5 生产调用者名册与"行为到底变没变"的判语

（待填。）

## §6 既有钉射程检查

（待填。）

## §7 推翻清单

（待填。）

## §8 判不动／量不到的格子

（待填。）

## §9 门禁与收尾

（待填。）
