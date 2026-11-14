# 票 443 — Test52「超时杀树」仪器在夹具没报出孙进程 pid 时判红：两载体、复现、牙、普查

- 腿＝票 443 写腿（委派）；2026-11-14 16:1x–17:0x +08；起手 HEAD `d2cfa96e8`（写腿期间共享树被别的腿推进过，本件读数全部在**本腿工作树**现取，逐条标了命令与 log 路径）。
- 写面：`internal/plugin/manifestslot/command_52_test.go`、`command_52_main_test.go`、`command_138_cancel_test.go`（三枚均为 `*_test.go`）
  ＋ `docs/evidence/s1/443-instrument-abstention.md`（本件）＋ `.scratch/wisp/probes/443/**`（新）＋ 票 443 Progress log 尾一条 bullet。
- ⛔ 未动：`internal/plugin/manifestslot` 的任何非测试文件（产码 0 字节，§6）、阈值／golden／契约、`cmd/wisp/**`、`scripts/**`、`frontend/**`、其它腿的未提交改动。⛔ push。

## 一、AC#1 复现：同一个「夹具没报 pid」，本机 N 与命中率

### 1.1 机理（先写清，否则下面的数字只是噪音）

`Test52TimeoutKillsTheProcessTree` 给夹具 `hang` 模式 500ms 预算（`plugin.MinCommandTimeoutMS * 10`）。夹具是**本包的测试二进制自己**（约 90MB，`installFixtureCLI` 复制一份再 re-exec），`hang` 模式要再 spawn 一个孙进程、由孙进程写 pid 文件。预算从 spawn 起算，**早于**夹具把孙进程造出来；runner 一冷／一挤，树就在孙进程存在之前被杀，pid 文件永远不出现。

修复前那一支是：

```go
gcPID := readPIDFile(t, pidFile)
if gcPID <= 0 {
    t.Fatalf("the fixture never reported a grandchild pid (%s)", pidFile)
}
```

逐字取自 `git show d2cfa96e8:internal/plugin/manifestslot/command_52_test.go`。CI 两发红句 `command_52_test.go:748: the fixture never reported a grandchild pid` 就是这一行。
⇒ 这是**调度结果被读成缺陷**，与票 441 同一机理。

### 1.2 自然读数（不设旋钮）

| 台面 | 分母 | N | 弃权命中 | 真红 | log |
|---|---|---|---|---|---|
| 修后、本机、无旋钮（首批） | 每次 1 发 | **25** | **3/25 = 12%** | 0 | `.scratch/wisp/probes/443/03-natural-25.txt` |
| 修后、本机、无旋钮（写腿末批，改完 `clearPIDFile` 重命名后） | 每次 1 发 | **25** | **0/25** | 0 | `13-natural-25-postedit.txt` |

**合并 N=50 ⇒ 3 发落到「夹具没报 pid」这一支 = 6%**。这 3 发在修复前会**逐字**打出 CI 那句红；
修复后它们打出 §2.2 的弃权行并**由控制载体把杀树问题答掉**，所以那 3 发是绿的——绿的依据是控制运行，不是「没看见就当没事」。
3 发弃权的逐字读数（budget/elapsed/waited/polls 见行内）：

```
AC#2 ABSTAIN(443:timeout-tree:raced-no-pid) budget-ms=500 elapsed-ms=1317 waited-ms=1008 polls=81 pid-seen=no ...
AC#2 ABSTAIN(443:timeout-tree:raced-no-pid) budget-ms=500 elapsed-ms=911  waited-ms=1009 polls=84 pid-seen=no ...
AC#2 ABSTAIN(443:timeout-tree:raced-no-pid) budget-ms=500 elapsed-ms=863  waited-ms=1005 polls=82 pid-seen=no ...
```

★ 诚实交代：`elapsed-ms` 863–1317 都**超过** 500ms 预算，说明「预算到 → 杀树 → 返回」这条链本身要 0.4–0.8s（90MB 镜像的代价），这段时间正是孙进程可能来不及出生的窗口。

### 1.3 放大旋钮（票面要的「N 与命中率」，不是「跑一百次没红」）

夹具侧加了一个**只在复现时使用**的旋钮 `WISP52_443_FIXTURE_DELAY_MS`（`command_52_main_test.go` 的 `fixtureDelay443`）：让 `hang` 模式在 spawn 孙进程前先睡这么多毫秒。未设置＝0，**提交的套件里它永远是 0**，只有本证与将来复现会设它。

| 台面 | N | 弃权命中 | 控制载体 | 结果 | log |
|---|---|---|---|---|---|
| `WISP52_443_FIXTURE_DELAY_MS=1500` | **5** | **5/5 = 100%** | 5/5 都拿到控制 pid 且判「树已死」 | 5/5 PASS | `11-amplified-5-delay1500.txt` |
| `WISP52_443_FIXTURE_DELAY_MS=3000` | **20** | **20/20 = 100%** | 20/20 同上 | 20/20 PASS | `04-amplified-20-delay3000.txt` |

**1500 是复现值，3000 已踩到控制预算的边**（诚实边界，写进了夹具注释）：控制预算是 `MinCommandTimeoutMS*80 = 4000ms`，3000ms 延迟＋冷启动 1s 就贴着它；实测 `10-post-edit-abstain.txt` 第 1 发打出
`AC#2 RED(443:timeout-tree:control-no-pid)`——**红得对**（那一发连控制都没能报出 pid，仪器无证人），所以旋钮不是「越大越好」，1500 才是可复现的复现值。

### 1.4 旧码直接复现（上一条被中断的腿的 `prefix` 台面，本腿逐条复核内容）

`.scratch/wisp/probes/443/logs/prefix-delay700-clean-01..05.log` 是**上一条腿**（被机主重启打断的那条）留下的 5 发：`WISP52_443_FIXTURE_DELAY_MS=700`，跑的是**旧断言还在**的测试件。
本腿逐条打开核对（不是照抄文件名）：**5/5 都是 FAIL 且红句就是 CI 那句**，逐字：

```
--- FAIL: Test52TimeoutKillsTheProcessTree (0.96s)
    command_52_test.go:765: the fixture never reported a grandchild pid (C:\Users\swq\AppData\Local\wisp52-fixture-tests\grandchild.pid)
```

行号 **765** 而不是 HEAD 的 748：那条腿当时已在同一文件上半途加了东西，行号跟着移了——**断言是旧那一句**，所以这是「旧文件 + 700ms 旋钮 ⇒ 5/5 复现 CI 红句」的直接读数。
同目录 `logs/baseline-0*.log`（5 发、无旋钮）全是 `ok ...`，即修前**不设旋钮时本机也绿**——这正说明「跑绿了」不能当复现，旋钮与自然命中率要一起报。

### 1.5 与 CI 那一发对齐

票面读数：run `43567331249`（head `e1cc304e0`）红、run `43568581591`（head `1f0d04714`，只改 probes 与证据件）红、`a46352fd3` 绿 ⇒ **3 发红 2**，红句同 §1.1。
本机复现的是**同一支**（pid 文件缺席 → 旧码判红），命中率 6%（自然）／100%（1500ms 旋钮）。

## 二、AC#2 弃权 vs 真红：两个问句、两个载体、两个 token

### 2.1 拆法（一枚文件）

| 问句 | 载体 | 判法 | token |
|---|---|---|---|
| 这一发**看见**孙进程 pid 了吗？ | raced 运行后的 1s 轮询（`awaitPID443`） | 看见 → 继续判树；没看见 → **具名弃权**，不红 | `ABSTAIN(443:timeout-tree:raced-no-pid)` |
| 超时**还杀得掉这棵树**吗？ | raced 看见 pid 就用它判；没看见则由 **8x 预算（4000ms）的控制运行**判 | 树活 ⇒ **真红**；控制也报不出 pid ⇒ **真红**（夹具坏了，无证人） | `RED(443:timeout-tree:tree-survived)`、`RED(443:timeout-tree:control-no-pid)` |

另：`clearPIDFile` 在 raced 与 control 两次运行前都清 pid 文件。理由不是洁癖——`stableExeDir` 是**固定目录**（`%LOCALAPPDATA%\wisp52-fixture-tests`，跨 run 复用），上一发留下的 pid 是**死进程的证词**：`stillAlive` 会答「已死」，测试就会在**本发孙进程根本没被报出**的情况下变绿（假绿）；反过来一个被系统复用的 pid 会让它假红。

### 2.2 弃权行逐字（带数字、可 grep）

```
AC#2 ABSTAIN(443:timeout-tree:raced-no-pid) budget-ms=500 elapsed-ms=949 waited-ms=1006 polls=85 pid-seen=no pid-file="C:\\Users\\swq\\AppData\\Local\\wisp52-fixture-tests\\grandchild.pid" last-poll-error="open ...: The system cannot find the file specified." - the budget expired before the fixture reported the grandchild it was supposed to be killed with, so THIS run cannot attest the tree kill from the race; the carrier is the 4000ms control run below, not a green from absence
```

四个数字都可核：`budget-ms`（这一发给的预算）、`elapsed-ms`（调用实际耗时）、`waited-ms`/`polls`（轮询等多久、几次）、`pid-seen=no`（pid 是否出现过）＋ `last-poll-error`（缺席与「写了但读不动」是两种发现，分开印）。
`last-poll-error` 把「文件不存在」与「文件在但不是 pid」分开，这是 `readPIDFileErr` 相对旧 `readPIDFile` 的唯一增量（旧包装原样留给票 138 的调用点）。

### 2.3 控制载体为什么不是「放宽」

控制运行**不继承 raced 那一发的抽签**：同一夹具、同一 runtime、同一 kill 路径、预算 8 倍。它要同时满足三件事才绿：
①调用确实以 timeout 结束（`IsError && Truncated && contains "timed out"`，否则 `control-no-timeout` 真红）；
②它自己报出了孙进程 pid（否则 `control-no-pid` 真红）；
③那个 pid 在预算后**已经死了**（否则 `tree-survived` 真红）。
⇒ 它只替换「这一发能不能在 500ms 里把孙进程造出来」这一件运气，**不替换**「树死没死」这个断言。

树活的那一支在返回前先 `taskkill /T /F` 收尸再 `t.Fatalf`，所以红完不会把挂着的孙进程留给后续用例。

### 2.4 为什么不用 `t.Skip`

与票 441 同一条硬边界：`tools/d22scan/runtests.sh` 把 core 档顶层 `--- SKIP` 判**致命**（SKIP 不是 pass），而 `scripts/**` 不在本票写面。
⇒ 采用票面同一句话的后半句「也不许静默绿」的实现：**具名行＋数字＋可 grep token**，弃权不红，但绿必须由控制载体挣来。

### 2.5 普查顺带发现并修掉的兄弟（同句、同夹具目录）

`command_138_cancel_test.go:76` 是**同一句** `the fixture never reported a grandchild pid`，且同样读 `stableExeDir` 的固定目录、同样没有清旧文件。
它的硬红**保留**（它不是判错：预算 120s、等 10s、**没有杀树在飞**，10s 还报不出 pid 就是夹具坏了），只补了 `clearPIDFile` 一行，堵掉「陈旧 pid ⇒ 假绿／假红」那一面。
改后 `Test138TaskCancelKillsTheProcessTree` PASS（`12-post-edit-138.txt`）。

## 三、AC#3 牙还在：仓外导出树里的突变，真红点名，还原 SHA256 同串

### 3.1 导出树与突变形状（为什么不是「把 kill 整个拿掉」）

导出树 `C:\Users\swq\AppData\Local\Temp\wisp443-export`（`robocopy` 白名单拷贝：`go.mod`/`go.sum`/`cmd`/`internal`/`tools`/`forks`/各包 `testdata`，**仓外**、不含 `.git`/`.scratch`）。先跑一遍确认绿：`ok ... 1.571s`。

**先试过的错突变（记录下来，因为它解释了下一条为什么这样选）**：把 `Kill()` 的 `TerminateJobObject` 与 `KILL_ON_JOB_CLOSE` 一起摘掉 ⇒ **测试挂死 600s**，最后是 `go test` 自己的 10 分钟超时 panic（还有 2 枚 `wisp52-fixture` 残留进程被本腿清掉）。
原因在 `command.go:256` 的 `out := proc.Wait()`：直子进程命里就是 `block()`，**只有 kill 才让它退出**；kill 整枚没了，`Wait` 就永远不返回。⇒ 那个突变根本不是「杀树失效」，是「调用不返回」，拿它当 AC#3 会得到 panic 而不是点名红。

**采用的突变（历史缺陷形状）**：Job Object 的树杀**降级成单进程杀**——
`command_job_windows.go` 里 `childJob` 加 `pid` 字段，`Kill()` 由 `TerminateJobObject(h,1)` 改成 `OpenProcess(PROCESS_TERMINATE,pid)+TerminateProcess`，同时把 `JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE` 置 0（否则 `Close()` 仍会把树带走，突变看不见）。
这正是一句测试注释里写的历史缺陷：`a plain Process.Kill would leave behind`（只杀直子进程，孙进程活）。

### 3.2 三个 SHA256 与两侧读数

| 状态 | `command_job_windows.go` SHA256 | 读数 |
|---|---|---|
| pristine | `B597F1029314B0F85D2BDB6446DE8AFF3875CB944E9BF66D438660765A4BD9D0` | `ok ... 1.571s`（导出树基线绿） |
| mutated（树杀→单进程杀） | `7EFD1619DA1AE9249FFB642C0636C40E5FDC7C3E3E7CF9FCDA395EA1040DF3C9` | **FAIL**，见下逐字 |
| restored（拷回 pristine） | `B597F1029314B0F85D2BDB6446DE8AFF3875CB944E9BF66D438660765A4BD9D0`（**同串**） | `ok ... 2.234s`，`--- PASS` |

突变那一发的完整红路径（`05-mut-singlekill-raced.txt`）——**它先弃权、再由控制载体点名**：

```
AC#2 ABSTAIN(443:timeout-tree:raced-no-pid) budget-ms=500 elapsed-ms=1033 waited-ms=1004 polls=90 pid-seen=no ...
AC#2 control run: grandchild pid 58852 reached the test (budget-ms=4000 elapsed-ms=4311 waited-ms=1 polls=1)
AC#2 RED(443:timeout-tree:tree-survived) control run: grandchild pid 58852 is still running after the budget expired - the kill did not reach the tree
--- FAIL: Test52TimeoutKillsTheProcessTree (7.82s)
```

★ 这一发同时证了两件事：**杀树失效仍然真红**（点名 pid 58852），且**它是从弃权那一支红出来的**——即「弃权」没有把缺陷洗绿，控制载体把问题接住了。还原后同一发回到 PASS（`06-restored-green.txt`）。

## 四、AC#4 同类普查（先普查，它决定本腿改几处）

尺（可复跑，产物在 `.scratch/wisp/probes/443/`）：

```
git grep -n -I -E 'never reported|never appeared|did not appear|never showed|never arrived|never came back|did not report|never wrote|never produced' -- '*.go' ':(exclude).scratch'
```

现量＝**62 行**（其中 `*_test.go` 50 行；`.scratch` 台的探针不算，CI 不编译它们）。逐条按「这一行会不会因**调度**变成红」分类，**只有两枚属本族**：

| `file:line` | 句子 | 档 | 会不会弃权变红 | 本腿 |
|---|---|---|---|---|
| `internal/plugin/manifestslot/command_52_test.go:748`（旧） | `the fixture never reported a grandchild pid` | windows 也跑（无 build tag），test-windows | **票面这枚**，3 发红 2 | **已修**：两载体（§2） |
| `internal/plugin/manifestslot/command_138_cancel_test.go:76` | `the fixture never reported a grandchild pid` | 同上 | **不会**（120s 预算／10s 等待／无杀树在飞 ⇒ 缺席＝夹具真坏），但**缺清旧文件**（陈旧 pid ⇒ 假绿/假红） | **补 `clearPIDFile` 一行，硬红保留** |
| `internal/proc/singleinstance_windows_test.go:147` | `%q did not appear within %v` | windows | 不会：等的是子进程**确定的**一行 stdout，无 kill 竞争 | 不改 |
| `internal/watchdog/watchdog_test.go:603`、`config_poll_test.go:177` | `the watchdog goroutine never appeared in the registry` | 全平台 | 不会：2ms tick 对 5s 预算（2500x），缺席＝注册没发生 | 不改 |
| `internal/agent/spill_path_invariant_test.go:135` | `the artifact did not appear at the asserted path` | 全平台 | 不会：`Prepare` 同步写完再读，无等待 | 不改 |
| `cmd/wisp/panel_assets_299_windows_test.go:208` | `the page never reported at all` | windows | 不会：**缺席就是被测缺陷**（白窗），红是结论不是弃权 | 不改 |
| `cmd/wisp/resident_hotkey_live_258_windows_test.go:161` | `the process has not rebound the live keys`（12s 有界轮询） | windows live 档 | **形似**：活进程＋有界轮询，重载慢时可能到点没记录 | ⛔ 不改（不在本票写面／live 档另票）→ 见下「交回」 |
| `cmd/wisp/subagent_blocked_197_test.go:149` | `this run never produced a subagent roster row` | windows | 不会：driver 内同步计数，`id==""` 是生成失败而非调度 | 不改 |
| `cmd/wisp/slo_windows.go:527`、`cmd/balldebug/diff_windows.go:460` 等 | 同上句式 | 产码 | **不是自检**，是产码错误串 | — |
| `internal/command/roster_test.go:381`、`internal/config/writeguard_226_test.go:127` 等 | `never appeared in %d rows` / `did not report the keys` | 全平台 | 不会：表驱动覆盖检查，无等待窗口 | — |

★ **交回**：`resident_hotkey_live_258` 那枚（活进程、12s 有界轮询、到点即 `LIVE RED`）与票 441 §4 里另外 6 枚「到达数」臂同类：它的缺席**可能是**调度。它跑在 live 档（需真实桌面），本腿不碰（写面外，且要另设计确定性载体）⇒ 建议下一票按 §2.1 的两问句处理。

## 五、AC#5 门禁四数（本机现跑，逐条命令）

| 门 | 读数 | 命令 |
|---|---|---|
| `go vet` | **rc=0** | `go vet ./internal/plugin/manifestslot` |
| `gofumpt -l`（本包） | **空**（0 文件） | `gofumpt -l internal/plugin/manifestslot/`（`go install mvdan.cc/gofumpt@v0.12.0`，`v0.12.0 (go1.27.1)`，与 CI 同 pin） |
| `gofumpt -l`（CI 全量 `internal cmd tools scripts frontend`） | **1 枚非本腿**：`internal/models/acl_sid_121_test.go`（HEAD 已提交、工作树干净；`gofumpt -d` 显示整文件 CRLF→LF 重写，不是本腿产物） | 同上，范围改全量 |
| d22scan | **clean**：`d22scan: clean - no D22 ban violations`，产码扫 482 枚（internal 404 / cmd 78），正控 `PASS=45 FAIL=0 SKIP=0`，rc=0 | `sh scripts/d22scan.sh`（Git sh）→ `07-gate-d22scan.txt` |
| path-length | **VERDICT RED，原因不在本腿**：`over-budget=70 covered=69 not in roster=1`，缺的是 `.scratch/wisp/issues/274-no-nail-requires-shipped-exe-to-carry-page-build-ps-step-2-comment-is-false-in-both-halves-done.md`（相对 123 字符、文件名 **102** 字符 > 100 帽；HEAD 已提交、工作树干净，**不是本腿写的**，也不在本票写面） | `sh scripts/check-path-length-budget.sh` → `08-gate-pathlen.txt` |
| 产码 0 字节 | `git diff --name-only -- internal/plugin/manifestslot/ ':!*_test.go'` ＝ **空** | 同上 |

定向测试（单条、串行跑，⛔ 未跑 `cmd/wisp`）：

```
go test -run 'Test52TimeoutKillsTheProcessTree$' -count=1 ./internal/plugin/manifestslot   # PASS
go test -run 'Test138TaskCancelKillsTheProcessTree$' -count=1 ./internal/plugin/manifestslot # PASS
```

## 六、边界自证与撤销

- **产码 0 字节**：本腿只改三枚 `*_test.go`；`command_job_windows.go` 等产码在仓库里**一字未动**（`git status --porcelain -- internal/plugin/manifestslot/` 只有那三枚 ` M`）。
- **突变只发生在仓外导出树**：仓库内该文件 SHA256 未变；导出树已还原成 pristine 同串（§3.2）。
- **未碰**：阈值／golden／`thresholds.go`／契约／`go.mod`／`go.sum`／`.github/**`／`scripts/**`／`frontend/**`／其它票的台账与票面正文（本票只追加一条 Progress bullet）。
- **本件的射程**：证「这一枚仪器的弃权不再判红、弃权带数字、杀树失效仍真红、还原同串、同族普查给定点名表」。
  ⛔ 证「CI 从此不会红」（还有 §4 表里那些臂与别的红源）；⛔ 把本机自然 N=50 的 6% 读成 CI 的实际频率。
- **撤销口令**：「票 443 撤」＝删本件＋回退本腿 commit（新 commit 反改，不 reset/amend）；`WISP52_443_FIXTURE_DELAY_MS` 不设即为 0，删旋钮等价于去掉 `fixtureDelay443` 与两处 env 传递。
- ⛔ `add -A`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`；commit 用显式 pathspec；⛔ push。

## 七、落地复跑（commit `d78a0ae5c` 的字节）

落地后工作树对本腿路径为干净（`git status --porcelain -- internal/plugin/manifestslot/ docs/evidence/s1/443-instrument-abstention.md .scratch/wisp/probes/443` ＝ 空），在该字节上复跑三条（串行，未跑 `cmd/wisp`）：

| 跑法 | 读数 |
|---|---|
| `go test -run 'Test52TimeoutKillsTheProcessTree$' -count=1`（无旋钮） | **PASS**（`raced run: grandchild pid 44696 … is gone`） |
| 同上 + `WISP52_443_FIXTURE_DELAY_MS=1500` | **PASS**，先 `ABSTAIN(443:timeout-tree:raced-no-pid) … elapsed-ms=755 polls=96 pid-seen=no`，再由 `control run: grandchild pid 59440 … is gone` 接住 |
| `go test -run 'Test138TaskCancelKillsTheProcessTree$' -count=1` | **PASS**（补了 `clearPIDFile` 之后） |

原始 log：`.scratch/wisp/probes/443/14-post-commit-rerun.txt`。

