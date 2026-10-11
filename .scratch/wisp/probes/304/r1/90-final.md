# `304-r1` 90 — 一页结论（读数，⛔ 判语）

**本腿＝实现者**：⛔ 翻任何框、⛔ 宣布任何 `AC` 成立、⛔ push。判语归下一程非实现者腿（`304-v1`）。

## 交了什么
- `AC#1` 落地形＝甲（编排者已裁，本腿⛔ 重新裁形）：两枚 mockllm spawn 点**同批**接上在库的 `internal/proc.JobScope`（`KILL_ON_JOB_CLOSE`），windows 走 `StartInJob`、非 windows ＝**具名的 no-op**（⛔ 假装受保护）。文件:行 ⇄ 逐枚说明＝**`10-impl.md`**。
- `AC#2` 成对两发＝两枚 spawn 点各发 (b)/(c)，含编排者欠的 `Q-304-v2-A`（helper 那枚的三形实测）：正控 4 发全 `count_repo_golden=0` **且**每发同行落"坏世界真被造出来"的证据行（`panic: test timed out after 3s` × 2／`parent_present_after_stop=False` × 2），反形 4 发逐枚落**具名 PID＋父 PID 现量已没**（`25952/22740`、`7876/22432`、`10184/6744`、`29264/32160`，四枚都带 `REPO_GOLDEN PARENT_GONE` 行＋`t.TempDir()` 首段＝测试名的归属行），还原＝`cmp`×4 YES ＋ `sha256sum` 双证 ⇄ `logs/landed/`。逐字＝**`20-paired-readings.md`** ＋ `logs/paired-fwd.txt`(68)／`logs/paired-rev.txt`(76)。
- `AC#3` 八把尺、每把件内一行 `rc=N`：`build_rc=0`、`vet_rc=0`、`d22scan_rc=0`（`clean - no D22 ban violations`，分母 `internal/` 533 Go files 含注释与 `_test.go`）、`pathlen_rc=0`（GREEN，`tracked=10828 over-budget=57 roster=57 not-in-roster=0 worst=228`）、格式名册**两把并排**（工作树 ⇄ `git show HEAD:` blob 取到仓外再量）＝**两把都 0 枚**、`bash scripts/portable-tests.sh --scope=census` 的 totals 行**逐字未变**＝`packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7 unclaimed-with-tests=0`（跑前⇄跑后作差 `IDENTICAL`，且⇄票面 `AC#3` 写的现量基线逐字相等）。逐字＝**`30-gates.md`** ＋ `logs/g-*.txt`／`logs/over-reach.txt`。
- 顺带一把正向对照（⛔ 判语用）：落地树上两枚 spawn 点各跑一发**正常结束**形 (a)＝`rc=0`、5 枚顶层 PASS／6 枚 PASS、`panic=0`、`count_repo_golden=0`（尺＝`-Quick`，`logs/p1-shape-a.txt`，17:38）；第二症状 `temp_test_dirs` 在我这轮从 24 走到 33（同因递增，⇄ 票面 `:16` 的 21 与 `A850` 的 23；⛔ 归因、只报数）。

## commit（只 commit、⛔ push；每笔都带显式 pathspec，`git add` 与 `git commit` 各一遍）
`b9f37892`（锚件）→ `e875f00f`（4 枚 `internal/llm/**` ＋首批件）→ 本件所在两笔见 `logs/over-reach.txt`（逐笔 `git show --name-only --format=`）。
起手索引现量＝**0 枚**（`git -c core.quotePath=false diff --cached --name-only | grep -c ''`）⇒ 别家路径 0 枚、我⛔ 清过任何人的东西、⛔ `reset`/`stash`/`checkout .`/`clean`。
⚠ 共享工作树窗口要具名：为交反形，`internal/llm/adaptertest/mockllm.go` 与 `internal/llm/openaichat/mockllm_integ_test.go` 在 **17:50–18:05** 呈"⛔ 带 job"的原形（用 `git show b9f37892:<path>` 重定向取回、随后逐字节还原），这段时间在飞的其他腿看到的是无 job 的那两枚文件。

## 本腿欠的读数（逐枚编号，⛔ 我自己销）
- **`Q-304-r1-1`＝ubuntu/linux 那一支的实测**：本轮零读数。落地形在非 windows 那支＝`return cmd.Start()`（`childjob_other.go:21-23`），票面 `:25` 甲形代价里那句"ubuntu 那一支仍会漏"我这轮⛔ 推翻了它、也⛔ 修了它；要实测需要 linux/WSL 台面（本腿只有本机 windows）。
- **`Q-304-r1-2`＝`go build` 那一枚子进程的残留形状**：零读数。`mockllm.go:65`／`mockllm_integ_test.go:65` 的 `go build` 用 `CombinedOutput()` 同步跑、边界②（⛔ 把编译器子树塞进同一枚 job）⇒ 我⛔ 把它接进 job；外力恰好杀在 build 那一瞬会漏什么，本轮⛔ 造过、⛔ 量过。
- **`Q-304-r1-3`＝census 的"改前那一发"由我自己跑**：见 `30-gates.md` §3 第 1 条（唯一呈改前原形的窗口花在同一条窗口的 `AC#2` 反形上）。
- **`Q-304-r1-4`＝托管 runner 上 job 的射程边界**（票面 `AC#0` 必答 (I)／台账 `Q-304-a1-B`）：⛔ 我能交，需 CI 色＝归编排者按在推送那批。本腿只能交"读码级"那半：每枚 job 只有一个显式成员（`StartInJob` 一枚）、`Assign` 我一次⛔ 调。
- **`Q-304-r1-5`＝走 helper 的其余测试文件⛔ 逐枚重跑**：`304-a1` 表①数到 17 枚文件走 `adaptertest.StartMockllm`，本腿只跑了 1 枚用例穿过 SP1（`TestMockllmSharedFaultSuite`）；逐枚跑＝⛔ 整包（在飞腿占了 Go 编译面，且派单 §4 禁大输出）。
- **`Q-304-r1-6`＝`-timeout 3s` 之外档位的 (b) 形**：帽是按 `logs/p1-shape-a.txt` 的实测时长（4.924s／3.798s）选的单一档位；不同帽下"panic 落在子代存活期内"的普适性⛔ 证。
- **`AC#4`／`AC#4b`（起手闸门与终态锚的名册更正）⛔ 本腿射程**：写面⛔ 包含模板件与 `docs/reports/**`，派单 §3 交付件里也⛔ 有它 ⇒ 归编排者/下一腿；`A850` 已注明"丙的另一半与甲⛔ 互相替代"。
- **仪器史自报（⛔ 读数缺陷、是纪律缺陷）**：`R-304-r1-1`／`R-304-r1-2`＝我用同名件覆写了自己前两发的原始日志（旧仪器：观测与终止分两次 powershell 调用 ⇒ 反形 (c)-SP1 那一发 `Stop-Process: Cannot find a process with the process identifier 22456`＝父已自行跑完退场、量的是没发生的世界）。修正＝`watch-kill.ps1`（同一进程内观测＋终止），随后两向各重打一遍＝本文所引读数。前两发的号（残留 `28504/24960/10564/16644`、正控子代 `32044`）只活在会话转录里，**盘上已⛔ 存在**。

## 派单 ⇄ 票面比对（⛔ 静默择一）
逐条对过：落点两枚（票面 `:15` ⇄ 派单 §1）、⛔ 新造 job 封装（票面 `AC#1` ⑤＋`A850` §4 ⇄ 派单 §1）、边界①②（`A850` §4 ⇄ 派单 §2.5/§2.6）、`AC#3` 那把 census 基线（票面 ⇄ 派单 §3 逐字同）。**冲突＝0 处**。唯一⛔ 一致的是**行号快照**：派单 §1 引 `mockllm_integ_test.go` 的"现有 `t.Cleanup`"在原 `:106`，我落地之后那行在 `:113`（同一枚 `t.Cleanup`，语义⛔ 动）；票面 `:7` 自己立的"引用前先重跑"就是这个意思 ⇒ 只作过期提示上报，⛔ 改任何原句。
