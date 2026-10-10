# 301-a1 锚点尺（第 1 笔，⛔ 引工作树）

跑法逐字（可重跑）：

```
date '+%Y-%m-%d %H:%M:%S %z'
git rev-parse --short HEAD
git status --porcelain -- internal cmd scripts .github
```

stdout 原样：

```
2026-10-10 13:21:15 +0800
--- HEAD ---
772ff880
--- porcelain(internal cmd scripts .github) ---
 M internal/audio/wasapi_windows.go
?? internal/audio/parse_wave_format_300_windows_test.go
--- rc=0 ---
--- branch ---
dev
```

`rc=0`

## 读数

- HEAD 短号 = `772ff880`（与票 301 现量段所引 HEAD 同枚）
- 分支 = `dev`
- porcelain 枚数 = **2 枚**（尺＝`git status --porcelain -- internal cmd scripts .github`，射程＝该四目录，量的是**工作树对 HEAD 的差异**）：
  - ` M internal/audio/wasapi_windows.go`（已跟踪、已改、未提交）
  - `?? internal/audio/parse_wave_format_300_windows_test.go`（**未跟踪**＝票 300 落地腿的在飞件）
- ⇒ 本程所有 `internal/audio` 的判据一律引 **HEAD blob**（`git show HEAD:<path>`），⛔ 引工作树；
  上面那枚未跟踪件**不构成本票表① 的分母**，它是否会被表① 拉进来属 `AC#1` 的代价面，由编排者裁。
