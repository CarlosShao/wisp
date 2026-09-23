# 票 128 AC#1 — `%APPDATA%` 缺失时数据根回落到 CWD：把**后果**量成读数

**测量会话：** `q7` · **锚定 sha：** `2620836`（dev）· **日期：** 2026-09-23（首条读数 09:37:40 +08:00）
**范围：** 仅 AC#1（四样落点在两个 CWD 下的实际路径 + 第二遍读到什么 + `icacls` 归属 + GUI 腿拒绝独立复现）。
**AC#2/AC#3/AC#4 不在本票段内**（AC#2 是语义裁定，见票面 Progress log 的 `next=`）。

## 0. 仪器与安全前提

- **仓外临时根：** `/tmp/wisp128-q7/` ⇒ Windows 实路径 `C:\Users\swq\AppData\Local\Temp\wisp128-q7\`。
  两个 CWD 是它的两个子目录（`cwd-A`、`cwd-B`），全部写点被限制在这棵树里。
- **owner 真实数据目录：未写入。** 开工前定位到两枚，且**开工前两者均为空目录**（`find -type f` 计数 = 0、`du -sb` = 0）：
  - `C:\Users\swq\AppData\Roaming\wisp`（目录 mtime 2026-09-19T14:49:20）
  - `C:\Users\swq\AppData\Roaming\wisp-dev`（目录 mtime 2026-09-20T07:06:25）
  ⇒ 基线是"零文件"，所以任何落在里面的文件都是**可判定的写入**；收尾复核见 §6。
- **不用 docker**（本机 Git Bash 下 `docker -v "C:\…"` 会静默挂空且 rc=0 = 假绿）。
- **不取时序/内存读数**：本机 `slo-full` 就跑在 self-hosted runner 上（A103），本段无需排队。

### 0.1 一处必须写明的偏离（编排者指令 vs 被量代码的形状）

派单要求"测量全程 `WISP_ENV=test` + 显式临时根"。**照做就量不到 AC#1**：
`resolveDataDir` 在 `env == "test"` 时提前 `return proc.TestDataDir()`，
根本走不到 `base, err := os.UserConfigDir(); if err != nil { base = "." }` 那三行。
所以本段主腿用 **`WISP_ENV=dev`**（落点 `wisp-dev`，与票面"已量"一节记的 `wisp-dev\logs` 同形），
`test` 只作对照腿（§5）。
隔离手段由"env=test"换成"**CWD 本身就在仓外临时根里**"——回落出来的 `wisp-dev\` 于是落在
`C:\Users\swq\AppData\Local\Temp\wisp128-q7\cwd-A\wisp-dev\`，与 owner 的
`C:\Users\swq\AppData\Roaming\wisp-dev\` 同名不同树；§6 的复核就是证这一点。

### 0.2 源码事实（读 `2620836`，非推断）

- `cmd/wisp/doctor.go:233` `resolveDataDir(env string) string`：portable 覆盖 → `env=="test"` → `proc.TestDataDir()` →
  否则 `base, err := os.UserConfigDir() // %APPDATA%`；`if err != nil { base = "." }`；
  再 `base = proc.SealableRoot(base)`；`env=="dev"` ⇒ `filepath.Join(base,"wisp-dev")`。
- `resolveDataDir` 的**四个生产调用点**：`cmd/wisp/doctor.go:101`、`cmd/wisp/doctor.go:276`（`dataDirForDisplay`）、
  `cmd/wisp/models.go:110`、`cmd/wisp/providers.go:81`、`cmd/wisp/run.go:155`。
  ⇒ run 腿的日志/`config.toml`/`memory.db` 都跟着它。
- GUI/boot 那条腿走 `internal/proc/envfork.go:236` `DefaultLayout`：
  `if err != nil { return Layout{}, fmt.Errorf("proc: user config dir: %w", err) }` ⇒ **拒绝，不回落**。
- `cmd/wisp/secret.go:123` `resolveSecretLayout` 同样**直接返回错误**（`wisp secret: user config dir: %w`），不回落。
- `SealableRoot(".")` 把 `.` 解析成**绝对路径**（`filepath.EvalSymlinks` 那半），所以日志上看到的不是 `.` 而是 CWD 的全名。

## 1. run 腿 · CWD-A（第一遍）

（待量）

## 2. run 腿 · CWD-B（第二遍，另一棵树）

（待量）

## 3. 四样落点汇总 + "第二遍读到的是不是同一棵树"

（待量）

## 4. `icacls` 归属读数（对照票 95 的私有目录纪律）

（待量）

## 5. GUI 腿在同一形下的拒绝（独立复现）

（待量）

## 6. owner 真实数据目录复核（收尾）

（待量）
