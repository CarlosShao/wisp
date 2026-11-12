# 424 — Linux 编译断裂：`memory_retention_180k.go`（无 build 约束）引用 Windows-only 的 `residentVoice` ⇒ CI `go vet ./...` 恒红、整个 Linux 构建挂

**Status:** open（2026-11-12 10:5x 编排者建；来源＝`A1764` 现跑：CI run `43139602748` 的 `go vet (module)` step 逐字 `vet: cmd/wisp/memory_retention_180k.go:73:38: undefined: residentVoice`；本机复核 `residentVoice` 的 type 定义只在 `resident_voice_windows.go:375`）
**Type:** 缺陷（跨平台构建断裂）
**Blocks:** 票 85 AC#5（`go vet ./...` rc=0）、票 122 AC#6（`lint` 步绿）、以及**任何**要求 CI 绿的 AC
**Blocked by:** nothing
**Packages:** `cmd/wisp`
**Introduce:** `e83494cf5`（票 423 wave 2 / 180-census-followup leg K）

## 现象（现跑铁证）

1. **CI**（run `43139602748` / job `130923107457` / step 11 `go vet (module)`）：**failure**，逐字
   `vet: cmd/wisp/memory_retention_180k.go:73:38: undefined: residentVoice`
2. **本机语义复核**：`type residentVoice struct` **只在** `cmd/wisp/resident_voice_windows.go:375` 定义（`_windows.go` ⇒ 只在 Windows 编）。
   而 `cmd/wisp/memory_retention_180k.go` **无 build 约束**（文件名无 `_windows` 后缀、首行仅 `package main`）。
   ⇒ **任何非 Windows GOOS，`residentVoice` 未定义 ⇒ `cmd/wisp` 编不过 ⇒ `go vet ./...` 先编译失败。**
3. **为何 `e83494cf5` 自报"vet rc0"**：那条腿在 **Windows** 上跑的 vet；跨平台缺口只在 Linux 显形。

## 四个函数、调用面（现跑）

| 函数（`memory_retention_180k.go`） | 定义行 | 调用点 | 平台 |
|---|---|---|---|
| `retentionWindowDays180K` | :40 | `run.go:705` | **跨平台** |
| `memoryL1Enabled180K` | :56 | `config_readers_255.go:257` | **跨平台** |
| `attachMemoryL1IfEnabled180K(rv *residentVoice, …)` | :73 | `resident_windows.go:759`（+ `config_wire_180k_windows_test.go`） | **Windows-only**（唯一用到 `residentVoice`） |
| `retentionBootLine180K` | :86 | `run.go:709` | **跨平台** |

## 修法（编排者拍，理由：最有价值＋主流 harness）

**只把 `attachMemoryL1IfEnabled180K` 一枚函数移进新文件 `cmd/wisp/memory_retention_180k_windows.go`**（该文件首行加 `//go:build windows`，`package main`）。
**其余 3 个跨平台函数留在 `memory_retention_180k.go` 不动。**

**为什么不是"整文件改名 `_windows.go`"**：那会把 `retentionWindowDays180K`／`memoryL1Enabled180K`／`retentionBootLine180K` 也变成 Windows-only，而 `run.go`／`config_readers_255.go` 是跨平台调用者 ⇒ **会把 Linux 另一处编挂**。**必须只切那一枚函数。**

**为什么不是"给整个文件加 `//go:build windows`"**：同上，那 3 个函数 Linux 也要用。

## AC

- [x] **AC#1 迁移函数**：`attachMemoryL1IfEnabled180K`（含它用到的私有 helper，若有）从 `memory_retention_180k.go` 移到新 `memory_retention_180k_windows.go`（`//go:build windows`）；`memory_retention_180k.go` 保留 3 个跨平台函数、删掉移走的那个。
- [x] **AC#2 Linux 编过（判据核心）**：`GOOS=linux CGO_ENABLED=1 go build ./cmd/wisp`（或 CI 等价）**必须 rc=0 且不再出现 `undefined: residentVoice`**。⚠ 若本机无 cgo/onnx 工具链导致 sherpa 依赖先挂，改用 **`GOOS=linux go vet` 的 type-check 面**或说明"本机无法完整交叉编译"并以"CI run 不再报 undefined: residentVoice"为判据。
- [x] **AC#3 Windows 不回归**：`go build ./cmd/wisp` rc=0；`go vet ./cmd/wisp` rc=0；`go test ./cmd/wisp -run 'Test180K|Test423|TestTicket255' -count=1` 与改前同绿（点名）。
- [x] **AC#4 越界**：只动 `cmd/wisp/memory_retention_180k.go` ＋ 新建 `cmd/wisp/memory_retention_180k_windows.go`；⛔ 改任何断言/阈值/契约/其它包。
- [x] **AC#5 门禁**：`gofumpt -l <两个文件>` 空、`go vet ./cmd/wisp` rc=0。

## 禁区

- ⛔ 改 `internal/**`、`frontend/**`、`design/**`
- ⛔ 改任何 `_test.go` 的断言（测试可跟着改 build tag，但不许改断言）
- ⛔ `thresholds.go`／golden／`C1–C32`／`D1–D47`

## Progress

- [2026-11-12 10:5x +08] agent=编排者 建票 + 现跑取证（见上）。派腿修复。

## Progress（续）

- [2026-11-12 11:0x +08] agent=编排者 用 424-fix 的 overlay 技法（`GOOS=linux CGO_ENABLED=0 go build -overlay go.mod→sherpa桩 -gcflags=all=-e`）拿到**权威全量错误清单**。`go vet` 只报每包第一枚，`-gcflags=all=-e` 才全量；结果 **5 处 / 2 族**：
  ```
  cmd/wisp/logsink.go:224:80              undefined: configErrText
  cmd/wisp/memory_retention_180k.go:62:28 undefined: configErrText
  cmd/wisp/panel_config_store.go:404:6    undefined: applyBallStyle65
  cmd/wisp/panel_config_store.go:410:6    undefined: applyBallIdleOpacity
  cmd/wisp/panel_config_store.go:415:6    undefined: applyBallClickThrough
  ```
  ⇒ **424-fix 只修了第一枚 residentVoice；这 5 处（configErrText×2 + applyBall×3）是同一族的第二批**，须第二腿修完 Linux 才真绿。
- [2026-11-12 11:0x] agent=编排者 判每枚性质（决定搬法）：
  - `configErrText(err error, c *config.Config) string` = **纯函数**（`resident_audio_windows.go:475`，只 err.Error()/nil 判，无 win32）⇒ **搬进跨平台文件**。
  - `applyBallStyle65`/`applyBallIdleOpacity`/`applyBallClickThrough` = **薄包装**（各自 `if <applier> == nil { return false }; return <applier>(…)`），applier 变量 `ballStyleApplier65`/`ballIdleOpacityApplier`/`ballClickThroughApplier` 现声明在 windows 文件 ⇒ **把 3 个 applier 变量声明 + 3 个包装函数搬进跨平台文件**（Linux 下 nil 返 false），`installBallStyleApplier65`/`installBallIdleAppearanceAppliers`（真赋值，接 `*ball.Ball`）**留 windows 文件**。
- [2026-11-12 11:0x] agent=编排者 **派第二腿 `424-fix-b`（glm high）**修这 5 处。判据：`GOOS=linux go build -overlay … -gcflags=all=-e ./cmd/wisp` **零 undefined 行** ＋ Windows build/vet/test 不回归。

## Progress（424-review 验收）

- [2026-11-12 11:2x +08] agent=424-review 验收：CI run 43146479546 step11 go vet=success；尺3 linux undefined=0、win build/vet rc=0；翻 AC#1-7。判词见 `docs/evidence/s1/424-verdict.md`。

## AC（续）

- [x] **AC#6 第二批 5 处清零**：`configErrText` 搬跨平台；3 个 `applyBall*` 包装 + 3 个 applier 变量声明搬跨平台；`install*` 留 windows。判据 `GOOS=linux go build -overlay … -gcflags=all=-e ./cmd/wisp` 零 `undefined:` 行。
- [x] **AC#7 两平台不回归**：Windows `go build ./cmd/wisp` rc=0、`go vet ./cmd/wisp` rc=0、`go test ./cmd/wisp -run 'Test180K|Test423|TestTicket255|Test65|Test255'` 与改前同绿。
