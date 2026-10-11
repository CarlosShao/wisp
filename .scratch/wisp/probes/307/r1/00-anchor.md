# 307-r1 · 00-anchor —— 起手锚（commit-first 第 1 笔）

## 尺与读数（逐字现跑）

- `date` = Sun Oct 11 10:09:26 CST 2026
- `git rev-parse HEAD` = cb7bb3ce33e2769e5e617479cd085451838057de
- `git branch --show-current` = dev
- `git status --porcelain -- internal cmd` = （空输出＝0 行，有效读数，非缺尺）
- `tasklist //FI "IMAGENAME eq wisp.exe"` = INFO: No tasks are running which match the specified criteria. ⇒ 0 枚
- `tasklist //FI "IMAGENAME eq balldebug.exe"` = 同上 ⇒ 0 枚
- `tasklist //FI "IMAGENAME eq mockllm.exe"` = 同上 ⇒ 0 枚

三把进程尺都 0 枚 ⇒ 台面干净，可比。

## 落点形状（编排者裁语，票面追加节 §3，逐字照做）

1. `syncdirs.go` 内容锚 `	if !s.complete && s.home != "" && isUnder(cand, s.home) {` ⇒ 把 `isUnder(cand, s.home)` 换成 `aUnder(<候选侧 pathForms>, anchorForms(s.home))`
2. `syncdirs.go` 内容锚 `		if isUnder(cand, e.canon) {` ⇒ 同族换 `aUnder(<候选侧 pathForms>, anchorForms(e.canon))`
3. ⛔ 存值一字节：`syncdirs.go:88` / `provenance.go:345`（home 两枚赋值点）、`syncdirs.go:152`/`:156`/`:160`（c.canon 三枚）、`finalize()`、`resolveTarget()` 里 `errTargetUnverified`／`ErrReparseDenied` 两句 fail-closed
4. ⛔ 新造折法：只复用 `formsOf`/`anchorForms`/`aUnder`/`aRelates`/`uncertainAnchorMiss`
5. ⛔ `risk` 包非测试面写 `filepath.Clean`/`filepath.Abs`（`307-a1` 现量＝该面这两枚 0 命中）

## 写面名册（本腿授权面）

- `internal/risk/**`（产码 1 枚文件：`syncdirs.go`；新测试夹具 1 枚 windows-tagged `_test.go`）
- `.scratch/wisp/probes/307/r1/**`

## 本文件⛔ 做到的

- ⛔ 颜色读数（本节只交锚）。改前/改后两台面名册色见后续件。
