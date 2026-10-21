# 票 124 落地腿 124-r1 —— internal/config 一批：把用例侧的未解析根接上 proc.SealableRoot

本件是**实现者**交件，不是裁决。裁决表由验收方另出（`docs/evidence/s1/124-*.md` 既有七份不属本腿）。

## 0. 一句话

把 `internal/config` 用例侧 6 处「把 `t.TempDir()` 的答案直接交给密封底线」的**递根点**
改成先过 `proc.SealableRoot`（复用本包**早已存在**的那枚一行委托 `sealableTempDir124`，不新增走法），
软链 TMPDIR 形状下本包 red 名 **33 -> 0**（only-in-link 口径），误伤 0。

## 1. 开工前两把尺（我自己现跑，未抄派单）

(1) 写面脏态：`git status --porcelain -- internal cmd` = **空**（写面内无别人未提交改动）。全树脏项全在 `.scratch/**` 等非写面。

(2) 派单基线尺复跑 = **20**，与派单一致；但读法要写清，我用三把尺同数才算：

```
$ git grep -n 'winsec\.' HEAD -- '*_test.go' | grep -E 'TempDir|os\.TempDir' | wc -l
20
$ grep -c '' <list>   => 20
$ awk 'END{print NR}' <list> => 20
$ tail -c 1 <list> | od -An -tx1 => 0a      (末行有换行，不是缺行被少计)
```

**一处钝尺登记**：同一发我先用 PowerShell `Select-String` 读出 **21**，是假读数——
PowerShell 把 `git grep` 输出数组化后末元素被重复计一次。改用 Git Bash 逐字同形得 20，
再用 `grep -c ''` 与 `awk END{NR}` 两条独立尺交叉确认为 20。**派单的 20 成立。**

## 2. 票面前提已腐坏（这是本腿最重要的读数，请裁决方先看这一节）

那 20 处**全部**已是 `internal/winsec/SealableTempDirForTest124(t)` 的调用点（19 处调用 + 1 行注释），
**没有一处还在把裸 `t.TempDir()` 交给 winsec**。两条硬读数：

```
$ git grep -n -E 'winsec\.[A-Za-z0-9_]*\([^)]*TempDir' -- '*_test.go'
(0 命中，rc=1)
$ git grep -n 'winsec\.' -- '*_test.go' | grep -E 't\.TempDir|os\.TempDir'
(0 命中，rc=1)
```

=>  派单那把尺量的是「已转换后」的形态，今天**恒等于已修状态**，不再能测「还剩多少没转换」。
**但票的终判据并未因此成立**：按 AC#2b 跑法量软链形，HEAD 上仍有 only-in-link red **108 枚 / 37 枚文件**：

| 包 | only-in-link red | 定义在几枚文件 |
|---|---|---|
| `internal/config` | 33 | 8 |
| `internal/tools` | 38 | 16 |
| `cmd/wisp` | 37 | 13 |

其中自家块内就打印 `not provably resolved` 的 74 枚（config 32 / tools 17 / cmd-wisp 25），
余 34 枚是间接（经 `secret.NewStore` / `memory.Open` / config manager 递根），含票 123 那族 300 s 腿三枚（按账豁免）。
**本腿只做 `internal/config` 这一包**（33 枚 / 8 文件，最内聚、与别票零交集）；
`internal/tools` 与 `cmd/wisp` 逐名清单交回编排者排后续批（`cmd/wisp` 是共享编译面，按 AGENTS.md 每波只许一枚写腿）。

## 3. 逐枚文件：原来交了什么、现在过哪一跳

| 文件 | 原来直接交给 winsec 的是什么 | 现在过哪一跳 | 覆盖 red |
|---|---|---|---|
| `internal/config/manager_test.go` | `newTestManager` 的 `dir := t.TempDir()`（`SaveFile` 封它所在目录）；`TestManagerMissingFileKeepsCurrentAndErrors` 同形一处 | 两处 => `sealableTempDir124(t)` | 多数 writeguard/settings 用例的根 |
| `internal/config/settings_257_test.go` | `cleanMachineManager` 的 `t.TempDir()`；`...AC1HandAddedRows...`；`...AC2MissingFile...`（`filepath.Join(t.TempDir(), ...)`）；`...AC2MissingRow...`；`...AC2InvalidValue...` 共 5 处 | 5 处 => `sealableTempDir124(t)` | 257 族 3 枚 |
| `internal/config/enabled_261_test.go` | `writeEnabledFile` 的 `filepath.Join(t.TempDir(), "config.toml")` | => `sealableTempDir124(t)` | 261 族 1 枚 |
| `internal/config/writeguard_226_test.go` | `TestGuardedWriteFailureRollsBackMemory` 的 `t.TempDir()` | => `sealableTempDir124(t)` | 226 族 |

**未改动、按名登记为「有意不改」**：`internal/config/c26_seam_posix_125_test.go`（2 处）——
票 125 的 seam 探针，**故意**看未解析拼写并在宿主临时根自带链接时 `t.Skipf`；
给它们解根会把「拒绝测量」换成「测错形状」，该文件的既有注释已写明这条边界。

**没有新增 helper**：`sealableTempDir124` 是票 124 前几批已入树的一行委托（`tempdir_resolved_124_test.go`，
`return proc.SealableRoot(t.TempDir())`），本腿零新增走法 =>  票 125 `R-125-2` 那本副本账不增行。

## 4. 判定分支未动（票面 AC#3）

```
$ git diff --numstat -- internal cmd
1	1	internal/config/enabled_261_test.go
6	2	internal/config/manager_test.go
5	5	internal/config/settings_257_test.go
1	1	internal/config/writeguard_226_test.go
```

删除侧**全文只有 6 行 `t.TempDir()` 那一个表达式**（其余 3 行为注释增补），
新增侧除同名替换外只有 `manager_test.go` 的 4 行注释。**断言、阈值、golden、判定分支 0 hunk。**
禁改列（`internal/winsec/**`、`internal/proc/**`、`internal/risk/**`、`tools/d22scan/**`、
`allowlist.txt`、`thresholds.go`、`go.mod`/`go.sum`、`frontend/**`、`docs/PLAN.md`、`docs/specs/**`、
`docs/reports/**`）**一律零 hunk**。

## 5. 三读数（逐字贴，全部 `golang:1.27` / `go1.27.1 linux/amd64`，每形一枚新容器，`GOPROXY=off`）

跑法：纯净快照 `git archive HEAD | tar -x`（挂载自证 `go.mod` 1026 字节）；软链形
`ln -s /realpriv /varlink` + `TMPDIR=/varlink/w124tmp`，`exit 97/98` 硬断言；普通形 `exit 99`（`/varlink` 不许存在）。

```
形状            读数        RUN  PASS  FAIL SKIP  mech
PRE-LINK      (改前基线)    186   150    35    1    35
MUT-LINK      (变异)       186   143    42    1    41
RESTORED-LINK (还原)       194   191     2    1     0
POST-LINK     (原样)       194   192     1    1     0
PRE-PLAIN     (改前普通形)  196   194     2    0     0
POST-PLAIN    (改后普通形)  196   196     0    0     0
```

- **变异红**：`MUT-LINK FAIL=42`，比基线 35 **多 7 枚**（6 枚新增 + 1 枚既有 flake 也翻红），
  且基线 35 名**逐名全部回归**（差集一侧为空）。
- **还原绿**：`RESTORED-LINK` 与 `POST-LINK` 同形同数，mech 归 0。
- **原样绿**：普通形 `POST-PLAIN 196/196/0/0`，改前 `PRE-PLAIN 196/194/2/0`。

> 本件因 160 行帽拆三份：本件 = sec 0..5（读数与文件清单）；`124-r1-impl-part2.md` = sec 6（变异自证与逐名红册）；`124-r1-impl-part3.md` = sec 7..10（残余红归因、门禁、交回项）。
