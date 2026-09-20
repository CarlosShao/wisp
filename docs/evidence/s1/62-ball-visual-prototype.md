# 62 — 液态玻璃悬浮球视觉原型：owner 签收证据（进行中）

**Ticket:** .scratch/wisp/issues/62-liquid-glass-ball-visuals.md
**Scope of this file:** 可运行原型 + 实测证据。**SPEC-08 §2 与 docs/PLAN.md 未改动**（契约变更须
owner 实机签收后由编排者走 D22 签字流程，见文末"待签字的 §2 措辞"）。
**Machine:** DESKTOP-LVS7839 · Windows 10 Pro 26100 · i7-8750H（12 逻辑核）· 单屏 3440x1440 ·
浅色（白）桌布 —— 即 owner 实机。

---

## 0. owner 要跑的那一条命令

```
build\balldebug.exe -stay
```

（在仓库根目录 `D:\work\workspace\projects plans\Wisp` 下执行；先
`set PATH=%CD%\third_party\sherpa-onnx;%PATH%`。）
退出：托盘图标右键 → Exit，或命令行 Ctrl+C。

## 1. 基线（改动之前）：Sleeping 态在白色桌布上不可见 —— 数值证明

测量方法（**不是**自绘合成图）：同一个屏幕区域截两次，一次悬浮球进程活着、一次把它正常关掉，
两帧逐像素相减。工具：`build\balldebug.exe -diff`（票 62 新增，父进程拉起子进程、正常收尾、
再截同区域，保证"死帧"里没有球也没有托盘残留）。

命令：

```
build\balldebug.exe -diff docs\evidence\s1\62-diff-baseline -diff-states Sleeping,Listening,Speaking -x 1600 -y 500 -diff-margin 90 -diff-sample 4s
```

图片：`62-diff-baseline\01-Sleeping-alive.png` / `-dead.png` / `-diff.png`（diff 已放大 6 倍，
所以 2/255 的差别也看得见）。数字表：`62-diff-baseline\diff-table.txt`。

| 状态 | 变化像素数（差值>=8/255） | 最大差值 | 单核 CPU | 私有工作集 | 活动定时器 |
|---|---|---|---|---|---|
| **Sleeping（现状 12px 微点 / opacity 0.35）** | **0**（>=3 的也只有 17 个） | **4/255** | 0.000% | 10.60MB | 无 |
| Listening（对照组，球是可见的） | 1471 | 115/255 | 0.008% | 12.23MB | 有 |
| Speaking（对照组） | 4824 | 102/255 | 0.027% | 12.90MB | 有 |

结论：**现状的 Sleeping 球在 owner 桌面上的可见像素变化为 0**，整个 232x232 区域里最大只有
4/255（1.6% 亮度）的差别 —— 人眼不可能分辨，这就是本票真正的 bug。对照组证明这套差分截屏
方法能看见球（Listening/Speaking 分别有 1471/4824 个像素变化 >=8/255）。

同时基线记录了 D32 现状：Sleeping 单核 CPU 0.000%、私有工作集 10.60MB、句柄 362、
GDI 10 / User 6、**零活动定时器**（timers=no 是从球对象上实测读出来的，不是查策略表）。

> 说明：`commit` 列里的 privateCommit ≈ 100MB 是 Go+cgo 进程的**提交量**（reserve 未落页），
> 与 D32 的"私有工作集"口径不同；本票一律以 private working set（任务管理器同列）为准。

## 2. 静态玻璃 Sleeping 本体（零定时器）

（进行中 —— 见文末 next=）

## 3. Summon 液态流动 + 音频驱动旋转

（进行中）

## 4. 静默边框过渡 + 靠边吸附收缩

（进行中）

## 5. 待签字的 SPEC-08 §2 措辞（owner 签收后才回填契约）

（原型完成后填写）
