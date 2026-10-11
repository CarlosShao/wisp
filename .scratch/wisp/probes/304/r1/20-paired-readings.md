# `304-r1` 20 — `AC#2` 成对两发（两枚 spawn 点各发一遍，含 `Q-304-v2-A`）

每发都落两件事：**那个坏世界真被造出来了**的证据行 ＋ 残留读数。⛔ "0 枚"当读数（`A850` §2#3 那一形：帽设得比用例还长 ⇒ 量的是一件没发生的事）。
尺＝`probes/orch/20261011/mockllm-procs.ps1`（`Get-CimInstance Win32_Process Name='mockllm.exe'`，归属守卫＝命令行含 `*Wisp\internal\llm\testdata\golden*`）＋本腿新造的 `watch-kill.ps1`（**在同一个 PowerShell 进程里**先观测到"我这一发的子代活着＋父 PID 就是我起的"再 `Stop-Process`）＋`cmdline-dump.ps1`（按 `t.TempDir()` 首段＝测试名归属，票面 `:11-:14` 那把尺）。台件＝`paired.sh fwd`／`paired.sh rev`，逐字读数＝`logs/paired-fwd.txt`（68 行）／`logs/paired-rev.txt`（76 行）。

## 正控＝落地之后那枚树（`logs/paired-fwd.txt`，18:05–18:06，射程＝工作树）
| 发 | spawn 点／形 | 坏世界证据行（同行落件） | 之后 `count_repo_golden` | 处置 |
|---|---|---|---|---|
| 1 | **SP2** `openaichat/mockllm_integ_test.go:73`，(b) `go test -timeout 3s` | `bad-world-evidence: panic: test timed out after 3s` ＋ `world_constructed t=3x0.5s … count_total=1 count_repo_golden=1 pids=[23872]`（`go test` 还在飞） | **0**（`VERDICT-b-SP2: world_constructed=1 ; count_total=0 count_repo_golden=0 pids=[]`） | `clean: nothing-to-clean`（＝没东西可停） |
| 2 | **SP2**，(c) 外力 `TerminateProcess` | `watch-kill: world_constructed: child pid=28320 ppid=22852 cmdline_has_repo_golden=YES` ＋ `parent_present_before_stop=True` ＋ `terminate_issued pid=22852` ＋ `parent_present_after_stop=False` | **0**（`VERDICT-c-SP2 … target_pid=22852 ; count_repo_golden=0 pids=[]`） | `clean: nothing-to-clean` |
| 3 | **SP1** `adaptertest/mockllm.go:72`（经 `-run ^TestMockllmSharedFaultSuite$` 走到，＝`Q-304-v2-A` 欠的那枚 helper），(b) | `bad-world-evidence: panic: test timed out after 3s` ＋ `world_constructed … pids=[26912]` | **0**（`VERDICT-b-SP1: world_constructed=1 ; count_repo_golden=0`） | `clean: nothing-to-clean` |
| 4 | **SP1**，(c) | `watch-kill: world_constructed: child pid=20436 ppid=31584 …` ＋ `terminate_issued pid=31584` ＋ `parent_present_after_stop=False` | **0**（`VERDICT-c-SP1 … count_repo_golden=0 pids=[]`） | `clean: nothing-to-clean` |
| — | 收尾整台面 | — | `=== end dir=fwd 18:06:10 final: count_total=0 count_repo_golden=0 pids=[]` | — |

⚠ 归属行在正控里是**空**的：`attribute` 那把尺（`cmdline-dump.ps1`）在四发之后都只落一行 `count_repo_golden=0` ⇒ **上面空着的那几栏＝0 枚**（⛔ 把"没打印"读成"尺没打中"）。

## 反形＝把新加的 job 那一步摘掉（`logs/paired-rev.txt`，18:01–18:02，射程＝同一枚工作树的两枚文件回到 `git show b9f37892:` 的原形）
| 发 | 形 | 具名残留 PID ＋父 PID 现量 | 归属（测试名在 `t.TempDir()` 首段） | `count_repo_golden` | 按三守卫停回 |
|---|---|---|---|---|---|
| 1 | **SP2** (b) | `pid=25952 ppid=22740 REPO_GOLDEN PARENT_GONE` | `tempdir_testname=TestMockllmGoldenByteIdenticalEvents3140925201` | **1** | `killed pid=25952 ppid=22740 guards=name=mockllm.exe+repo-golden-path+parent-gone` |
| 2 | **SP2** (c) | `pid=7876 ppid=22432 … PARENT_GONE`（父＝我起的 `paired.test.exe` pid=22432，`parent_present_after_stop=False`） | `…3292549456` | **1** | `killed pid=7876 ppid=22432 guards=…parent-gone` |
| 3 | **SP1** (b) | `pid=10184 ppid=6744 … PARENT_GONE`（`logs/paired-rev.txt:57-58`） | `tempdir_testname=TestMockllmSharedFaultSuite2868402985` | **1** | `killed pid=10184 ppid=6744 guards=name=mockllm.exe+repo-golden-path+parent-gone` |
| 4 | **SP1** (c) | `pid=29264 ppid=32160 … PARENT_GONE` | `…TestMockllmSharedFaultSuite1603513437` | **1** | `killed pid=29264 ppid=32160 guards=…parent-gone` |
| — | 收尾 | — | — | `final: count_repo_golden=0` | 

四发的父 PID 现量都⛔ 在（`PARENT_GONE`／`parent_alive=False` 两种写法同一把尺，逐字分别在 `logs/roster-rev-*.txt` 与 `logs/cmdline-rev-*.txt`），且 `-CleanIds` 那把尺**只有三条守卫全满足才停**（名字 `mockllm.exe` ＋ 命令行含本仓 golden ＋ 父 PID 现量不存在），父还在就打 `skip pid=… reason=parent-still-alive`——本轮四发⛔ 出现 skip（尺＝`grep -c "skip pid" logs/paired-rev.txt` 现量 **0**，⚠ 命中 0 ⇒ `grep -c` rc＝1，按 `|| true` 取）。

## 两发的坏世界⛔ 相同＝成对才成立
正控 4 发与反形 4 发用**同一个** `paired.sh`（同顶、同 `-timeout 3s` 帽、同归属守卫），差别只有那一步 job：`grep -c StartChildGuarded internal/llm/adaptertest/mockllm.go internal/llm/openaichat/mockllm_integ_test.go` 现量＝反形 `0 0`（件里逐字：`reversed_marker(0=job step absent): …:0 …:0`）→ 落地版 `2 2`（`logs/restore-verify.txt` 的 `guard_marker(2=job step back)`）。

## 还原双证（`logs/restore-verify.txt`）
`cmp_identical …=YES` × 4 枚（`mockllm.go`／`mockllm_integ_test.go`／`childjob_windows.go`／`childjob_other.go`）＝**0 枚 NO**（尺＝`grep -c "=NO" logs/restore-verify.txt` 现量 `0`，⚠ 命中 0 时 `grep -c` rc＝1，我是按 `|| true` 取的）。
`sha256sum` 双证（`logs/landed/*.go` ⇄ 工作树同路径，四枚逐字相等）：
`cfdc71f5…` mockllm.go ⇔ `cfdc71f5…`；`3c5d6034…` mockllm_integ_test.go ⇔ 同；`1df2b14f…` childjob_windows.go ⇔ 同；`7f7c9e70…` childjob_other.go ⇔ 同（全号在件里）。

## 具名自报两处仪器史（⛔ 藏）
1. **反形第一发（旧仪器）里 (c)-SP1 那一发⛔ 算执行过**：观测与 `Stop-Process` 分两次 powershell 调用，之间隔了 ≈2.5 s 进程启动，而 `TestMockllmSharedFaultSuite` 编成 exe 后约 2.5–3 s 就跑完了 ⇒ `Stop-Process : Cannot find a process with the process identifier 22456`＝父已经自己退、`t.Cleanup` 已经跑过 ⇒ 那一发量的是"没发生的世界"（正是一枚假 0）。定式修正＝`watch-kill.ps1`（观测＋终止在同一进程里）。那一轮的 4 枚残留号（`28504`／`24960`／`10564`／`16644`，逐字在会话转录里）被同名件 `logs/paired-rev.txt` **覆写**了＝我的临时件纪律缺陷 ⛔ 删、只建不删我没守住（覆写⛔ 删但仍丢了第一发原件）⇒ 具名＝**`R-304-r1-1`**；正控旧仪器那发（`logs/paired-fwd.txt` 同名覆写）同理，号＝**`R-304-r1-2`**。判语⛔ 归我：这两枚都是"仪器/纪律缺陷"，⛔ 读数本身。
2. **⚠ `go build` 那一枚子进程⛔ 在 job 里**（边界②：job 只端显式成员，`go build` 一塞就等于把编译器子树带进去）⇒ 外力杀发生在 build 那一瞬的残留形状本轮零读数＝**`Q-304-r1-2`**（与上一条编号不同，此处 `Q-` 是欠读数、`R-` 是缺陷自报）。
