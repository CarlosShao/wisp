# 票 89 独立对抗验收（acceptor-ticket89）

**被验收对象**：`.scratch/wisp/issues/89-0600-is-decorative-on-windows-acl-for-private-data.md`（`Status: ready-for-review`，六框自勾）
**代码链**：`b994a2c`(骨架) → `3054007`(基线 icacls) → `de15a6b`(真 SD) → `57bdbb2`(接线 + AC#4/AC#5) → `0a3a445`(AC#6 门禁)
**验收环境**：本机真机，`go version go1.27.1 windows/amd64`，Windows 11，当前账户 `DESKTOP-LVS7839\swq`，
icacls = `C:\Windows\System32\icacls.exe`。本文件由**独立对抗验收代理**撰写，实现者的任何数字都不作为本代理的证据。

**快照纪律（全部在仓外 `/tmp`，仓内无 worktree、工作树未被本代理改动）**：

| 快照目录 | SHA | 用途 |
| --- | --- | --- |
| `/tmp/wisp89acc-89` | `0a3a445`（票 89 终态，`internal/winsec` 到 HEAD 无漂移，已用 `git diff --stat 0a3a445 HEAD -- internal/winsec` 核对为空） | 判据重跑 + 各项变异 |
| `/tmp/wisp89acc89-pre` | `3054007`（实现仍是 `os.Chmod` 占位） | 第 3 项"差点假绿"反向复现 |
| `/tmp/wisp89acc89-base` | `942ab5a` = `b994a2c^`（票 89 之前的生产码） | 第 1 项基线（本代理自造的仪器） |

---

## 进度总览（先落盘，边做边追加）

| 项 | 内容 | 状态 |
| --- | --- | --- |
| 2 | 纯净快照 `go test -count=2 -v ./internal/winsec/` 四数 | **已完成 → PASS** |
| 1 | 基线为假（本代理自造 icacls 仪器，票前生产码） | 进行中 |
| 3 | "差点假绿"独立证实（解析器退回跳过首行） | 未开始 |
| 4 | `PROTECTED_DACL` 变异 | 未开始 |
| 5 | AC#5 fail-closed + `sealError()` 归一化拿掉 | 未开始 |
| 6 | 覆盖面主张（`SealDir` 传播 + 出生那一刻继承） | 未开始 |
| 7 | `os.Remove` 删 junction 与票 79 A51② 的冲突 | 未开始 |
| 8 | `verifyPrivate` 白名单口径裁决 | 未开始 |

---

## 第 2 项：纯净快照判据重跑 —— **PASS（票面数字被独立复现）**

仓外纯净快照 `git archive 0a3a445 | tar -x -C /tmp/wisp89acc-89`，本代理自己的读数：

```
$ go test -count=2 -v ./internal/winsec/        # 在 /tmp/wisp89acc-89
ok  github.com/CarlosShao/wisp/internal/winsec  14.791s
RC=0

=== RUN   : 34        (全量 grep '^=== RUN' 计数)
--- PASS  : 20        (顶层；另有 14 条缩进的子测试 PASS 行，合计 34)
--- FAIL  : 0
--- SKIP  : 0
不同名测试数 : 17      (10 顶层 + 7 子测试名)
34 == 2 × 17  ✔  （每个名字恰好出现两次，无一例外，见下）
```

`=== RUN` 逐名计数（`sort | uniq -c`，全部恰好 `2`）：
`TestAC1BaselineProductionPaths` / `TestAC1BaselineForeignACEPropagatesIntoModeOnlyWrites` /
`TestAC2SealedDirCoversFilesItNeverTouched` / `TestAC3SealedWritesCarryNoForeignSID`（+ `artifact-exclusive`、
`secret-blob`、`staging-temp` 三子）/ `TestAC3ProductionDataRootIsPrivateEndToEnd` /
`TestAC4JunctionAtArtifactPositionIsNotRecursed` / `TestAC4UnlinkableLinkGivesNamedError` /
`TestAC4SealedWalkSkipsLinks` / `TestAC5FailedSealRefusesTheWrite`（+ `exclusive_artifact`、`replacement_write`、
`directory_chain`、`seal_file_and_dir_direct` 四子）/ `TestAC5FailureIsNotSwallowedByTheHappyPath`。
⇒ 票面"0 SKIP、17 条具名 PASS、无重名跑两遍冒充条数"**独立复现**。
结构上还有一条支撑：`grep -rn "t.Skip\|SkipNow\|testing.Short" internal/winsec/` = **0 命中**，本包没有任何 skip 通道。

**四种假绿逐条点名（本项）**：
1. `--- SKIP`：0 条，且包内无 `t.Skip` 调用 ⇒ 不是"skip 当 ok"。
2. `-run` 过滤：本项**没用** `-run`，全量输出计数。
3. 步骤被静默跳过：`assertPrivateACL`/`requireForeignModify` 走 `run(t,…)`，`icacls` 或 `powershell` 非零退出即 `t.Fatalf`；
   `TestAC1BaselineProductionPaths` 里 `wisp.db-wal` 缺失只会 `t.Logf` 不会 skip（本代理在第 1 项自造仪器里加了
   "一条都没测到就判失败"的 `if probed == 0 { t.Fatal }` 钉子，见下）。
4. 计数不达标 / 重名：34 = 2×17 且逐名 uniq -c 全为 2 ⇒ 票 79 那种"同一测试名跑两遍"未重现。

**残留观察（不算 FAIL，记给编排者）**：`TestAC1BaselineProductionPaths` 这个名字里的"基线"在 `0a3a445` 上
已经**不再是纯基线** —— `secret.NewStore`/`memory.Open` 在 `57bdbb2` 接了 winsec，所以那两类的 icacls 原文
测的是**修好之后**的状态（票面 AC#1 那段的原文是 `3054007` 时代录的，仍然成立）。⇒ 本代理不采信它的基线读数，
第 1 项改用票前树 `942ab5a` 上自己写的仪器重测。
