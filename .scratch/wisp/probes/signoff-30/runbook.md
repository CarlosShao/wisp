# 10:30 桌面签收 runbook（编排者自用，09-30 09:3x 现跑验证过第 0 步）

⛔ 全程三条硬规矩：**不 push**（触发＝签收结束）、**不跑 `go test`／整包**（抢机器）、**不派任何会自己开窗口的腿**。
`scripts/dev/ball-cycle.ps1` **不带 `-OutDir` 一律不许跑**（`scripts/dev/ball-cycle.ps1:19` 默认写进 `docs/evidence/s1/ball-states/`，那 20 张旧规格图是不可再生对照件）。

## 0 起跑门（10:05 与 10:28 各跑一遍，四条都要过）

```
date
tasklist //FI "IMAGENAME eq balldebug.exe"
tasklist //FI "IMAGENAME eq wisp.exe"
gh run list --status in_progress
ls -l build/balldebug.exe && find internal/ball cmd/balldebug -name '*.go' -newer build/balldebug.exe
```

- 09:30 现量读数（已跑通）：两枚进程 **零枚**、`gh run list --status in_progress` **空**、`find -newer` **空**（二进制＝今天代码）。
- ⛔ `gh run list` **非空 ⇒ 延后开始时刻**，不要边跑 CI 边演示（`slo-full` 历史最长 2h12m，跑在本机 self-hosted runner 上）。
- 缩放必须 **100%**（历史数字全在 100% 下测）。

## 1 静止态：先浅色壁纸（对照表第 1、2 件）

```
PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" ./build/balldebug.exe -stay
```
- 看：玻璃质感、边光、**完全不抖不闪**；`Ctrl+C` 收（会打 `exit: signal`）。
- ⚠ `-stay` 期间**不要**在旁边跑任何测量，读数会被洗。

## 2 20 种模样挨个轮演（第 17 件；⛔ 不用 `-tour`）

```
PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" ./build/balldebug.exe -cycle-ms 6000
```
- 09:29 已用 1 秒档干跑验证：**20 行 `state=` 按 D43 顺序全部出现**（FirstRun→…→Stuck），末行 `balldebug: OK`，`rc=0`，句柄 376→390（阈值 <600）。
- 6 秒档＝20×6＝**120 秒**。⛔ **不带 `-shots`**，不然会在盘上多出 20 张图。
- 他要单独盯某一态时：`-state Warm -hold`（会顺手打 `timers=` 并写 `-status` 文件）。

## 3 贴边缩回（第 6 件）

先用真拖（他用手，最像日常）；需要我复现时：
```
./build/balldebug.exe -stay -dock left
```
`-dock` 取值 `left|right|top|bottom|none`，**打错字会 `os.Exit(2)`**，不会假装贴在中间。
⚠ 这一件**当场只报现象、不翻勾**：票 62 AC#5 已由 `62-v1` 判**不成立**（契约 `DockOverlapFrac=0.42` 对实测 35/44≈**79.5%**；点得到 `r=21` < 可见 `26`）。

## 4 "音量"驱动液体（第 4、5 件）⚠ 不是真声音

```
./build/balldebug.exe -state Speaking -level 0.9 -hold      # 大
# Ctrl+C，再来一发小的：
./build/balldebug.exe -state Speaking -level 0.2 -hold      # 小
```
- 电平稳态来源逐字＝`synthetic audio envelope 0..1 pushed to the ball at ~30fps`（`cmd/balldebug/main.go:106`），`feedLevels` 是 `syllabified synthetic envelope`（`:396`）。
- ⛔ **不许说"语音通了"**：`SetAudioLevel` 全仓非测试生产者只有 `cmd/balldebug`；"真声音→这个数"那一截**没有生产者**（`62-v1` AC#4 判语＝附条件）。
- 他要"静音收敛进边框"那一形：`-state Listening -level 0`（0＝不喂，`-- :106` 那句 `0 = feed nothing`）。

## 5 深色壁纸那一发差分（他换完壁纸之后才跑）

```
./build/balldebug.exe -diff .scratch/wisp/probes/signoff-30/dark -diff-states Sleeping
```
- ⛔ **`-diff` 的目录必须是 `probes` 下的临时目录**；`balldebug` 自己要求显式给目录，所以不会覆盖档案——**但 `ball-cycle.ps1` 会**，那条路今天不走。
- `62-v1` 的账：浅色半侧她已自跑复现（`px>=8=2101`、框 40×45、teardown clean）；**深色半侧她今天不可满足**（没动 owner 桌面设置）⇒ 明早这一发是**全仓第一枚深色桌布差分**。

## 6 交完要落的两件账（不在演示期间做）

1. 台账 **A461**：签收结果逐格（他说过的原话要逐字），并把"翻／不翻"按他看到的写。
2. 推送触发到达 ⇒ 先 `git status` 看别人地界，再 `git push`，终判据＝**远程 tip 逐字等于本地 HEAD**；现量 **14 枚未推**（`origin/dev=d5a59d66`）。
