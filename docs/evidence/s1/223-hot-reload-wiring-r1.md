# 票 223 R1 实现证据件 —— CheckAndReload / ConfirmLocked / OnRestartPending 三枚断口接线

- 工单：`.scratch/wisp/issues/223-checkandreload-has-zero-production-callers-hot-reload-never-runs.md`
- 腿：`223-r1`（写码腿：产码＋产测试＋产证据件）
- 起手锚点：`git rev-parse --short HEAD` = **`7ffa9520`**，`date` = **`2026-09-29 14:34:53 +0800`**，分支 `dev`
- 台件目录：`.scratch/wisp/probes/223/r1/`（原始输出全量落盘，每发带各自 `date`）
- 起手普查件（本件的地形图，逐条现跑复认过其中五处更正）：`.scratch/wisp/probes/223/c1/census.md`（49,708 字节，锚点 `5a251ece`）

## 起手读数（本腿 14:3x 现跑，不复用他人数字）

| 尺 | 结果 |
|---|---|
| `grep -rn "CheckAndReload" --include=*.go . \| grep -v _test.go \| grep -v /probes/` | 生产调用者 **1 枚**＝`cmd/balldebug/main.go:243`；其余命中为注释（`internal/ball/hotkey_reload.go:21/29/58`、`internal/config/allowdirs.go:151`、`loader.go:60`、`manager.go:17/62/113/117`、`permmode.go:57`、`writeguard.go:109`） |
| `cmd/wisp` 里 `CheckAndReload` 命中 | **0** |
| `grep -rn "ConfirmLocked = " --include=*.go cmd/ internal/ \| grep -v _test.go` | **0 命中**（exit 1）＝断口二成立 |
| `OnRestartPending` 全仓命中 | `manager.go:33/36/146/147` ＋ `manager_test.go:82/94` ⇒ 生产赋值点 **0**＝断口三成立 |
| `OnReload` 生产赋值点 | 仅 `cmd/balldebug/main.go:244` |
| `internal/watchdog/` | 只有 `doc.go` |
| `grep -rn DEFERRED internal/config/` | `internal/config/doc.go:14` 一枚 `DEFERRED(schema/hot-reload/migration)`（是否可摘＝见末节，本腿**不摘**） |

## AC#1 生产里真有人在轮询（CheckAndReload 的 wisp run 调用者）

**未判。**

## AC#2 三档生效级别各有读数（立即／重启后＋明确告知）

**未判。**

## AC#3 放宽必带 L2 复确认（D33 正控，走生产路径；收紧不许弹卡）

**未判。**

## AC#4 不生效与读不到是四句话（缺失／语法错／权限不够／热加载被禁用）

**未判。**

## AC#5 整包终态读数＋逐名比红名册

**未判。**（起手绿名册已存 `.scratch/wisp/probes/223/r1-start-green-roster.txt`：`./internal/config ./internal/perm` 83 枚具名用例、0 FAIL、0 SKIP）

## AC#6 断口二一起接上（ConfirmLocked 生产赋值点 ≥1，且非锁内同步等待）

**未判。**

## AC#7 断口三一起接上（OnRestartPending 生产赋值点＋可观测出口）

**未判。**

## 票面写的 vs 我读到的（不一致单列）

**未判。**

## 没做完／判不了

**未判。**
