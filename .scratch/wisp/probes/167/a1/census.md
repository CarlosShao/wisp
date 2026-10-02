# 167-a1 — 面板四枚小出口＋崩溃自救：Go 侧真源取证（只读腿）

> 本文件由只读取证腿 `167-a1` 产出。**零产品码改动**；唯一可写文件即本文件。
> 票面：`.scratch/wisp/issues/167-four-small-outputs-the-panel-needs-occupancy-queue-stop-draft-plus-crash-self-rescue.md`（未改动，AC 框一枚未碰）。
> 搜索根一律显式：`cmd internal tools docs scripts .scratch`。`frontend/**` 与 `design/**` 未读、未引用、结论里不出现其行号。
> 未跑 `go test` / `go build`（两条写腿在飞）。允许动作只有 `grep` / `sed` / `ls` / `git log` / `go list`。
> 骨架先落（本轮只有 §0 与各区标题），§1–§6 由后续 commit 逐节补齐。

## 0. 起手锚（逐字读数）

`date -Iseconds`：

```
2026-10-02T09:29:29+08:00
```

分支（`git rev-parse --abbrev-ref HEAD`）：

```
dev
```

`git log -1 --format='%H%n%ad%n%an%n%s'`：

```
f5f9cc3433385c13e24cb2a7d8054d8210f99cde
Fri Oct 2 09:27:56 2026 +0800
CarlosShao
census(114-a2)：§1 调用者名册＋§2 逐格三态交齐——ParseComposerRequest 生产调用者=1(composer_dispatch.go:155)且 Handle 有两枚产码听众；未勾实测 9 枚不是 8 枚；248-r1 一格都没顺手做掉，反而把 AC#8 的覆盖面削薄(config.* 不带 panel. 前缀、两把尺都看不见)
```

`git status --porcelain internal/panel cmd/wisp`（起手读数，逐字）：

```
 M cmd/wisp/run.go
?? cmd/wisp/firstrun.go
?? cmd/wisp/firstrun_198_test.go
?? cmd/wisp/firstrun_acl_198_windows_test.go
```

起手读数解读（仅记录事实，非结论）：`cmd/wisp` 侧有在飞的写腿产物（`run.go` 已改 + 三枚 `firstrun*198*` 新件），与派单说明的 `198-r1` 占 `cmd/wisp` 相符；`internal/panel` 干净——本腿读到的 `internal/panel` 内容即 HEAD 的内容。

## 1. 名册（五枚出口 × 四问）

（后续 commit）

## 2. 逐枚"缺哪一环"

（后续 commit）

## 3. 最小落地面

（后续 commit）

## 4. 我可能判错的条目

（后续 commit）

## 5. 判不动的地方

（后续 commit）

## 6. 我推翻票面/派单哪一句

（后续 commit）
