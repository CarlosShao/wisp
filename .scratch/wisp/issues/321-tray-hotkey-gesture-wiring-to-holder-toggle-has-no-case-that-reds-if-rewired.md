# 321 — 托盘/热键那一跳的**执行器接线无钉**：把 `withPanelHost` 里的 `holder.toggle` 换成别的执行器，今天**没有任何用例会红**（`319-v2` 交回的具名缺口）

**Status:** 待派（小票；面＝`cmd/wisp` 的 319 宿主装配＋其测试件）
**Type:** 判据缺口（"接了但接线没人守"族；与票 121/131/133/135 同族）
**来源:** `319-v2` 的对抗验收表 `docs/evidence/s1/ticket319-319-v2-accept-v1.md` 的"toggle 语义"格（判**附条件**：语义有钉、**接线无钉**）；编排者 2026-10-15 17:4x **代码级复核** `cmd/wisp/resident_windows.go:268-270`（`withPanelHost(func(via string) bool { return holder.toggle(via) })`）。

## 为什么立案（读数不是感觉）

- 票 319 的 r3/r4 改了**二次启动**那条跳的执行器（`toggle`→`show`＋`raise` 两跳），并给两半各配钉＋变异自证；
- 但**托盘/球/热键**那一跳（装配处 `resident_windows.go:268-270`）仍走 `holder.toggle(via)`，而**没有任何用例断言"这一跳驱动的是 toggle"**——把它改成 `holder.show`（或任何别的执行器）今天不会有红。
- 反面（不许拿它冒充本格）：`TestBallPanelGesturesReachThePanelThread` 守的是"手势到得了面板线程"，**不区分执行器语义**（toggle vs show）⇒ 现有那两把尺都不接这一形。

## AC（1:1；裁决表 `docs/evidence/s1/321-*.md` 由**非实现者**出）

- [ ] **AC#1** 给"托盘/热键那一跳 → `holder.toggle`"装一枚**常驻判据**：把装配处换成 `holder.show`（`-overlay`／仓外变异）⇒ **必须红且红名点到"这一跳不是 toggle 语义"**；还原 ⇒ 绿＋`git diff` 0 字节。
- [ ] **AC#2** 门禁＋面核：`go build ./...`／gofumpt／`go vet -tags winlive ./cmd/wisp/`／`sh scripts/d22scan.sh`／长度门；改动逐枚 pathspec；⛔ 不碰 `frontend/**`／`design/**`／冻结件。
- 备注：可与票 255 的装配钉同批（同一枚文件族 `cmd/wisp` 的 assembly 测试），但**不许**用"手势到得了线程"那一枚冒充本格（见上"反面"）。

**撤销口令**：「撤 321」。

## Progress log (append-only, newest last)

- [2026-10-15 17:5x +08] 编排者立案：来源＝`319-v2` 的附条件格（319 收口时登记；319 本身五格齐、不因此票重开）。
