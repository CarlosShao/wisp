# 424-fix-b — Linux 编译断裂第二批（5 处 / 2 族）修证

- 腿码：`424-fix-b`
- 基线：`b528c760df66a5d603e5b6cb4acbef11ccf2319a`（进场现读）
- 性质：跨平台搬符号，函数体零语义改动；`install*`/seed 留 windows 文件。

## 1. 改了什么（逐字）

### 族 A：`configErrText`（纯函数）
- 新建 `cmd/wisp/config_err_text.go`（无 build tag，`package main`，import `internal/config`）：
  `configErrText` 全文搬入，体逐字不动，仅加一条"为何跨平台"的说明注释。
- `cmd/wisp/resident_audio_windows.go`：删除该函数及专属注释（调用点 2 处在文件内，自动可见）。
- 调用点 `logsink.go:224` / `memory_retention_180k.go:62` 未改动。

### 族 B：3 个 `applyBall*` 包装 + 3 个 applier 变量
- 新建 `cmd/wisp/ball_config_appliers.go`（无 build tag，`package main`）：
  - `var ballStyleApplier65 func(name string) bool`
  - `var ballIdleOpacityApplier func(v float64) bool`
  - `var ballClickThroughApplier func(on bool) bool`
  - `func applyBallStyle65(name string) bool` / `func applyBallIdleOpacity(v float64) bool` / `func applyBallClickThrough(on bool) bool`
  六个符号逐字搬入（nil 检查 + 调 applier 形状不变）。
- `cmd/wisp/ball_style_65_windows.go`：删 `ballStyleApplier65`、`applyBallStyle65`；`installBallStyleApplier65`、`ballStyleSeed65` 留守。
- `cmd/wisp/ball_idle_appearance_windows.go`：删 `ballIdleOpacityApplier`、`ballClickThroughApplier`、`applyBallIdleOpacity`、`applyBallClickThrough`；`installBallIdleAppearanceAppliers`、`ballIdleAppearanceSeed` 留守。
- import 清理：三个 windows 文件无未用 import 产生（各自仍用 `slog/filepath/ball/config`），零删。`go fmt ./cmd/wisp` 无输出。

## 2. 判据读数

### Linux 全量 type-check（复用第一腿 overlay `D:\tmp\wisp\20261112\424fix\overlay-gomod.json`，`Test-Path` = True）
```
GOOS=linux CGO_ENABLED=0 go build -overlay <gm> -gcflags=all=-e -o NUL ./cmd/wisp
undefined 行数 = 0      （全量错误输出为空；日志 D:\tmp\wisp\20261112\424fixb\linux-full.log）
```
5 处原 `undefined`（logsink.go:224 / memory_retention_180k.go:62 / panel_config_store.go:404/410/415）全部消失；无新族冒出。

### Windows 不回归（go 1.x `D:\work\base\go\bin\go.exe`，PATH 前置 sherpa-onnx + build）
```
go build ./cmd/wisp  → rc=0
go vet ./cmd/wisp    → rc=0
go test ./cmd/wisp -run 'Test180K|Test423|TestTicket255|Test65|Test255|Test324' -count=1 -v
  → rc=0；RUN=43 PASS=42 FAIL=0 SKIP=1（剥子测试＋去重四数；日志 D:\tmp\wisp\20261112\424fixb\win-test.log）
```

### 过程如实记
- 第一轮 Windows 构建曾报 `ball_style_65_windows.go:53` 语法错误：删除 `var ballStyleApplier65` 时误吞了 `installBallStyleApplier65` 注释首行（搬符号笔误，非语义问题），补回注释后复跑三关全绿。上表 rc 为复跑后读数。
- Linux 判据在搬完两个族后一次通过（上述语法错误只影响 Windows 侧第一轮）。

## 3. 越界自证（`git status --porcelain -- cmd/wisp docs/evidence/s1`）
本腿摸过的面：
```
 M cmd/wisp/ball_idle_appearance_windows.go
 M cmd/wisp/ball_style_65_windows.go
 M cmd/wisp/resident_audio_windows.go
?? cmd/wisp/ball_config_appliers.go        （本腿新）
?? cmd/wisp/config_err_text.go             （本腿新）
```
其余 `??`（如 `cmd/wisp/panel_resident_323*_test.go`）是别的腿既有件，本腿未提交、未触碰。internal/** frontend/** design/** 零写面。

## 4. 临时件落点
`D:\tmp\wisp\20261112\424fixb\`（linux-full.log / win-test.log），仓库内不落临时件。
