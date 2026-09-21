# A18 特征化：真 `taskkill /F` 打断暂存写盘之后，今天到底留下什么

**时间：** 2026-09-21 09:1x–09:2x（本地）· **执行者：** 票 20 接续代理（第 5 框同批顺手做的 A18 特征化）
**判据来源：** `docs/reports/pending-and-issues.md` **A18**（票 20 AC#3 的另一半未证）
**新文件：** `internal/tools/bridge_a18_kill_windows_test.go`（`//go:build windows`，零 `t.Skip`）

## 0. 这条为什么原来是假的（不是我又发现一遍，是 registry 已经写着的）

A18 的原文判断：`TestAtomicWriteKillsMidWrite` 用的是**进程内** `Hooks.Kill`
（`internal/tools/fs_write.go:47-53` 返回 error ⇒ **会**跑 Go 的清理路径），
所以 `fs_write_test.go:119-130` 那条"不留暂存文件"的断言在真实故障下**不可观测**。
我不重做这个判断，只补它缺的那一半：**真的把进程杀掉，然后看盘上有什么。**

## 1. 怎么做到"真"

- 父测试用 `exec.Command(os.Args[0], "-test.run=^TestA18…$", "-test.timeout=90s")` 重启**自己这个测试二进制**；
  子进程在同一条用例的开头按环境变量 `WISP_A18_CHILD_DIR` 走子分支（所以不需要 `TestMain`，本仓 `internal/` 零 `TestMain`）。
- 子进程通过**真桥**（`fsDepsBridge` + AnswerAllow 门）调 `fs.write`，`WriteChunk: 4`，
  在 `Hooks.AtStep("write:8")` 处**落一个信号文件后永久阻塞** ⇒ 被杀的那一刻它确实已经在写暂存区
  （不是"还没开始写"的廉价场景）。
- 父进程轮询到信号文件后跑 **`taskkill /F /PID <pid>`**（不是 `Process.Kill`，A18 的判据原文要的就是真 taskkill）。
  taskkill 跑不动 ⇒ `t.Fatalf` 写明需要 `System32\taskkill.exe`；**没有 t.Skip，也没有降级回进程内 Kill**。
- 子进程自己带两道反空跑：绕过阻塞点就直接 `os.Exit(5)` 打印"wrote without blocking"，
  父进程读到这句话就判失败；父进程若没杀成，`defer` 补一刀 + `Wait` 收尸，不留孤儿。

## 2. 真实结果（本机 2026-09-21，`-count=1` 与 `-count=2` 各跑过）

| 断言 | 结果 |
|---|---|
| 覆写分支：目标 `existing.txt` 在真 kill 之后 | **逐字仍是 OLD-BYTES**（D31 在真故障下成立：暂存建在目标目录 + 只有一次 `os.Rename`，`fs_write.go:277/325`） |
| 新建分支：目标 `fresh.txt` | **不存在**（既不是半截、也不是延迟出现） |
| 每次 kill 的残留 | **恰好 1 个 `.wisp-tmp-*`**，两次 kill = 2 个；实测每个 **4 字节**（= 被杀时已进暂存区的量，`t.Logf` 记着） |
| 残留会不会被"下次启动"扫掉 | **不会。** 全新 `PathCanonicalizer` + 全新 registry/桥在同目录做了一次**成功**写盘（合法路径全程走通），两个 `.wisp-tmp-*` 仍原地不动 |

第四条是**特征化**不是修复：全仓 `grep` 只有 `fs_write.go:39/277/497` 三处引用 `tempPrefix`，
**没有任何读取者** ⇒ 今天没有清扫器可称"下次启动"。用例把"仍在"钉成断言，
谁加了清扫器这条就会红，逼他连同 registry A18 判据②一起改（`t.Fatalf` 文本里写着这句）。
**我没有把"没扫"改成"扫掉"** —— 按 A18 原文，那要么是行为设计、要么该写进 SPEC，归 owner 判。

## 3. 变异 / 反空跑证据

- 反空跑①：若 `AtStep` 的阻塞没生效，子进程打印 `A18 child wrote without blocking` 并 `os.Exit(5)`，
  父进程读到即判失败（实测未触发 ⇒ 阻塞真发生过、kill 真发生在写盘中途）。
- 反空跑②：残留数量断言写死 `len == 2`，任何"顺手清理"都会让它红（见 §2 第四行）。
- 阳性对照：`fresh.txt` 那条同时断言 `Stat` 的 error **必须是** `IsNotExist`，
  把"权限报错"与"确实不存在"分开，防止用 err≠nil 蒙混成"没落地"。

## 4. 门禁数字

`go test -count=1 -v -run TestA18 ./internal/tools/` → 1 顶层 + 3 子项 `--- PASS`，0 SKIP（09:2x）；
本文件与第 5 框那条一起落地后复跑：`go test -count=2 ./internal/tools/` `ok 25.192s`、
`go test -race -count=1 ./internal/tools/` `ok 17.709s`（race 下子进程照常起、照常被打死）。

## 5. 没证明但看起来成立的事

- **没证明"崩溃时的 4 字节会不会被卷缓存丢"**：我只断言文件存在与目标未受影响；
  暂存文件在真 kill 后的**内容**是否等于已 Write 的量，取决于 Windows 缓存管理器的收尾，
  用例把大小放进 `t.Logf` 而**不**断言——想钉它得先决定要不要 `Sync` 语义变化，那是行为设计。
- **没证明"下次真正的 wisp 进程启动"**：`cmd/wisp` 的启动路径里没有任何 `.wisp-tmp-*` 处理，
  我按硬规矩没去 `cmd/wisp` 加（也不许改），所以"下一次启动"这里只能用"新建一套桥"来代表。
- **crash 发生在 `rename` 之后但 `Remove` 之前的分支未测**：那一支今天由 `Hooks.Kill` 覆盖不到，
  真 kill 的时机我固定在写盘边界，跨边界的时间点没有可重复的注入手段（需要调试特权）。
