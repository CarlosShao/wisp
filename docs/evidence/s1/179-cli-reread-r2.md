# 179-r2 交件表 — 真机 CLI 上「模型拿宿主指针调 `fs.read`」那一发跑到了，但**逐字节读回没拿到**：断点从上一程的台件写法挪到了 `R4` 命中升 `L2`

- 写手程：`179-r2`（派单 `.scratch/wisp/dispatches/2026-09-28-110x-impl-179-r2-land-the-cli-fs-read-reread-without-touching-production.md`）
- 起手时刻：2026-09-28 11:09:13 +08（`date` 现量，非抄派单号）
- 锚：分支 `dev`、起手 HEAD `5190eea6`（**与派单写的 `016a3b4a` 不符**：`016a3b4a` 是它的父一枚，`5190eea6` 就是这纸派单自己的 commit）
- 结论一句话：**AC#8 未结案**。179-r1 的修法在本读数里成立（`task.output` 真执行、`风险未分级` 那句话三区计数 0），本程把台件的取数写法改成按 `tool_call_id` 字段取，于是 `fs.read` 那一发**第一次真打在宿主自己写下的指针上**；它被 `level=L2 rules=[R4] reason="R4: 包含来自 task.output 的内容"` 拒掉，run B 退出 1、任务 `cancelled`。**零产码改动**（判"非改不可"＝停手上报，本程没到那一步——要继续走通的是**另一枚缺陷**，不是本票的修法）。
- 本程名下 commit：`391822a2`（台件＋读数 8 枚文件）；本表与票面 Progress log 追加为后续 commit。**票面 AC 框一枚没勾**（档位等 `179-v1` 非实现者裁）。

## 1. step-0 五件

| 件 | 读数 |
|---|---|
| `date` | `Mon Sep 28 11:09:13 CST 2026` |
| `git rev-parse --abbrev-ref HEAD` | `dev` |
| `git rev-parse HEAD` | `5190eea6ddd35e459a4fb8ae7bf3dad8f07238cc`（派单写 `016a3b4a`＝**不符**，已在顶格报回） |
| `git status --porcelain -- internal/ cmd/` | 空（起手与终态各现量一次，见 §6） |
| 基线三包 | `ok internal/agent 1.917s` / `ok internal/tools 15.294s` / `ok internal/risk 5.552s`（`risk` **单跑**，照 `A359`） |

尺：
```bash
date; git rev-parse --abbrev-ref HEAD; git rev-parse --short HEAD
git status --porcelain -- internal/ cmd/
go test -count=1 ./internal/agent/; go test -count=1 ./internal/tools/; go test -count=1 ./internal/risk/
```

## 2. 本程**没**测什么（免得这张表被当成全量绿灯）

- 没跑全仓 `./...`（票面 AC#7 禁）；只跑了 `./internal/agent/`、`./internal/tools/`、`./internal/risk/`（单跑）、`tools/d22scan/runtests.sh -C tools/d22scan ./...`（那是 d22scan 自己那十枚工具的用例）。
- **没跑** `cmd/wisp` 的既有跟踪用例：本程只用 `-overlay` 跑了名下那一枚续读腿（`-run TestRereadHostPointerOnCLISeam179r2`），所以 `scripts/wisp-cli-tests.sh` 那句 `PASS=33` 是**在册记录、不是本程现量**。
- 没跑 `.scratch/wisp/probes/161/r6/flip-declaration.sh`（跑一次就脏跟踪日志）。
- 没重跑 `probes/176/r1` 那发端到端（＝票 179 AC#1 那一格，仍待别人做）。
- 没动 SLO 阈值／golden／`thresholds.go`／`allowlist.txt`／`PassThroughUnclassifiedRisk` 的任何赋值点；没动 `internal/**` 一字节（§6 附尺）。
- 没读、没写、没引 `frontend/**` 与 `design/**`（别家地界）。
- 续读那一发的**逐字节读回**没测到（§4 说清卡在哪一句），所以票 177 AC#3 的端到端条件、票 175 AC#5、票 176 AC#3/AC#4 这四格本程**一枚都不算翻**。
- `fs.read` 读**别的**真文件（非宿主指针）那一形没测——本票只要指针那一发。

## 3. `cmd/wisp` 的 PATH 两形对照（逐字，两形都贴）

| 形 | 命令 | 读数 | 文件 |
|---|---|---|---|
| 带前缀 | `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test -count=1 -run XXX_NONE_PKG ./cmd/wisp/` | `rc=0`、`ok github.com/CarlosShao/wisp/cmd/wisp 0.085s [no tests to run]` | `.scratch/wisp/probes/179/r2/logs/path-with-prefix.txt`（66 字节） |
| 不带前缀 | `go test -count=1 -run XXX_NONE_PKG ./cmd/wisp/` | `rc=1`、`exit status 0xc0000135` + `FAIL github.com/CarlosShao/wisp/cmd/wisp 0.060s` | `.scratch/wisp/probes/179/r2/logs/path-no-prefix.txt`（76 字节） |

尺：
```bash
wc -c .scratch/wisp/probes/179/r2/logs/path-with-prefix.txt .scratch/wisp/probes/179/r2/logs/path-no-prefix.txt
```
⇒ 派单 §1 那条前提**核实为真**：加载期 `0xc0000135` 是 PATH 条件、不是代码形状（与票 98 的更正一致）。本程续读那一发全程带前缀跑。

## 4. 续读那一发：回执原文＋"逐字节相等"那两个数

台件：`.scratch/wisp/probes/179/r2/zz179r2_e2e_test.go`（overlay `.scratch/wisp/probes/179/r2/overlay-e2e.json`，编进 `cmd/wisp` 包走真 CLI 接缝）。状态机**只认工具行的 `tool_call_id` 字段**（`internal/agent/loop.go:715 -> toolResultMessage(c.ID, ...)`），这就是 179-r1 断掉的那一处；改完 `fs.read` 那一发真的发出去了。

跑法与逐字全文：`.scratch/wisp/probes/179/r2/logs/test-run2.txt`（25602 字节）、读数 `.scratch/wisp/probes/179/r2/logs/e2e-readings.txt`（28496 字节，`wc -c` 现量非空）。两跑（第一次误落 `cmd/wisp/logs/`、已原样搬成本程的 `logs/e2e-readings-run1-misplaced.txt`）字节数相同＝28496，读数可复现。

```bash
PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" \
  go test -count=1 -overlay=.scratch/wisp/probes/179/r2/overlay-e2e.json \
  -run TestRereadHostPointerOnCLISeam179r2 ./cmd/wisp/
```

**模型调用序列**（读数逐字）：`task.output c66c0634-9fe8-42fb-8e40-136bdc0fba73` → `fs.read C:\Users\swq\AppData\Local\Temp\TestRereadHostPointerOnCLISeam179r22423749009\001\artifacts\tool-output-agent-task-c66c0634-9fe8-42fb-8e40-136bdc0fba73.txt（整份，不带 max_bytes）`。
指针与名册登记那条**同一条路径**（台件里两向比对通过：`pointer != ToSlash(ArtifactPath) && pointer != ArtifactPath` 不成立）。

**第一发 `task.output` 的回执**（宿主给的，逐字；`isError`/层级来自同一次运行的 `[audit]` 行）：
```
[…输出已落文件：省略 17200 字符，总长 20000 字节 / 约 5000 token；注意：路径授权判定者未接线（fail-closed：C26 没接进来，这条路径是否还读得回来无法核实，按读不到处理），全文见 C:\Users\swq\...\artifacts\tool-output-agent-task-c66c0634-9fe8-42fb-8e40-136bdc0fba73.txt…]
```
```
[audit] wisp run: 风险判定 tool=task.output level=L0 rules=[] reason="无规则命中（L0 直接执行）"
[audit] tools: call kind=success task=c3d0f0f7-... tool=task.output risk=L0 decision=allow outcome=success rules_hit=[] in_allowlist_scope=true reason="无规则命中（L0 直接执行）"
[工具 task.output -> success]
```
⇒ 179-r1 那半到手的样子在这里第三次复现并留盘：**声明 L0 不再被自家门拒**。那句"路径授权判定者未接线（C26）"是票 174 AC#2b 的正交格，**不是**本发的阻断者（`task.output` 自己照样 `outcome=success`）。

**第二发 `fs.read`（照指针那一发）的回执**（逐字）：
```
[确认 L2 fs.read] R4: 包含来自 task.output 的内容
  规则 R4
  ...
  影响路径：C:\Users\swq\...\artifacts\tool-output-agent-task-c66c0634-9fe8-42fb-8e40-136bdc0fba73.txt
  窗口 0.0s
[工具 fs.read -> error]
[audit] wisp run: 风险判定 tool=fs.read level=L2 rules=[R4] reason="R4: 包含来自 task.output 的内容"
[audit] tools: call kind=refused task=c3d0f0f7-... tool=fs.read risk=L2 decision=reject outcome=error rules_hit=[R4] in_allowlist_scope=true reason="R4: 包含来自 task.output 的内容"
[audit] tools: PATH-ACCOUNT task=c3d0f0f7-... tool=fs.read roots=1 rewritten=[] unusable=[]
[audit] tools: C25 scope closed task=c3d0f0f7-... was_open=true dropped=1 open_scopes=0 close_err=<nil>
```
`isError`＝**真**（`outcome=error`，`decision=reject`）；层级＝**L2**（由 R4 命中升上来）；`in_allowlist_scope=true`（路径本身在 `[fs] allowed_dirs` 根内，**不是**越界升档）。

**"逐字节相等"那两个数**：
| 量 | 现量 |
|---|---|
| 产物字节数（`os.ReadFile(rec.ArtifactPath)`，名册登记那条指针所指） | **20000 字节**（且与 payload `long179r2` 逐字相等） |
| `fs.read` 读回字节数 | **0 字节**——`rows["call-179r2-full"]` 与 `rows["call-179r2-head"]` 都是空（`sha256=e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`＝空串的哈希），因为那一发在桥里就被 `R4` 拒了，环路随后 `cancelled`，工具行**没有回到模型上下文** |

⇒ **这一发没拿到**，卡在那一句逐字：`R4: 包含来自 task.output 的内容`。
本程台件另配了一发"有界读回"（`max_bytes=10000`，为的是让读回的字节整段留在上下文里好做逐字比对），它同样被同一句 `R4` 拒掉，所以两个数都是 0。读数里的 `run A 退出 0 / run B 退出 1 / 确认卡 0 -> 1` 就是这一路的收尾。

## 5. 四枚分岔的逐枚读数（含正控说明）

⚠ 先说尺的形状：`e2e-readings.txt` 里**自带**一行 `分岔「X」出现次数: console=… rig-turns=… tool-rows=…`，那行本身就写着 `X`，所以"整份文件 `grep -c X`"每枚都至少 1——那是尺的自污，不是产品里的命中。真读数取下面那三区的计数（台件里对 `console`／`rig turns`／`tool-rows` 三区分别 `strings.Count`），文件级 `grep` 只用来定位行。

```bash
L=.scratch/wisp/probes/179/r2/logs/e2e-readings.txt; grep -n "出现次数" "$L"
```
| 分岔 | 读数（console / rig-turns / tool-rows） | 正控（同一把尺在别处响得出来） |
|---|---|---|
| `风险未分级` | **0 / 0 / 0** | `grep -rn "风险未分级" --include=*.go internal/` ⇒ `internal/agent/loop.go:801` 与 `internal/agent/declared_l0_risk_179_test.go:42` 各 1 行 ⇒ 尺认得这句话，本读数里真没有 |
| `…L2 级…已拒绝执行`（loop 的审批未接入那一形） | **0 / 0 / 0** | `grep -rn "L2 级" --include=*.go internal/agent/loop.go internal/tools/gate.go` ⇒ 命中（`loop.go:789`、`gate.go:126/144`）⇒ 尺认得。**但要报尺不够宽**：升 L2 这件事在本读数里真实存在，只是拼写不同——`grep -c "L2" "$L"` ⇒ **5**，逐形 `grep -o "level=L2\|risk=L2\|确认 L2 fs.read" "$L" | sort | uniq -c` ⇒ `level=L2` 2、`risk=L2` 1、`[确认 L2 fs.read]` 1。所以"0 命中"只否定"`审批通道尚未接入`那一形"，**不否定 L2 升档本身** |
| R4 命中 | **8 / 0 / 0**（＝**响**，这就是本发的阻断者） | `grep -rln "R4" --include=*.go internal/tools/` ⇒ `bridge.go`、`gate.go`、`fs_test.go` ⇒ 尺认得；`[audit] ... rules_hit=[R4] reason="R4: 包含来自 task.output 的内容"` 逐字在 §4 |
| `查不到这个任务` | **0 / 0 / 0** | `grep -rn "查不到这个任务" --include=*.go internal/tools/task.go` ⇒ `task.go:215` 1 行 ⇒ 尺认得，本读数里真没有 |

⇒ 照票 179 AC#8 的分岔规矩：**冒出来的是第三枚（R4）并连带第二枚（L2 升档）**，不是第一枚。这条正是 AC#8 写死的"不许写成本票已结案，要具名报回并指回票 177/162"那一支。本程照实报：
- 指回**票 177 shape A（Q-61甲 精确豁免）／票 175-r2 桥级盖章**：`task.output` 里 `box.set(rec.ArtifactPath)` 那条豁免（`internal/tools/task.go:253-255`、`internal/tools/bridge.go:566-578` 的 `mark`/`hostPathBox`）**在这条路上没能护住续读**——模型这一发唯一自带的字符串就是宿主自己写下的那条路径，却仍判 `rules_hit=[R4]`。同一次运行里 `C25 scope closed ... dropped=1` 是一条要人看的线索（本程**没有**去动它，判据与产码一字节没改）。
- 指回**票 162 那一族**：`窗口 0.0s` → 立刻拒 → 任务 `cancelled`（`run B 退出 1`），这条审批/取消路径的形状与本票正交。

## 6. 有没有动过产码

**没有。**
```bash
git status --porcelain -- internal/ cmd/          # 起手空、每次提交后各现量一次＝空
git show --stat --name-only 391822a2              # 8 枚文件全在 .scratch/wisp/probes/179/r2/**
```
本程名下只写过三类文件：`.scratch/wisp/probes/179/r2/**`（台件＋读数）、本表、票 179 的 Progress log 追加行。
一处要自报的越界**形状**（不是越界改码）：第一跑的 `runtime.Caller(0)` 在 `-overlay` 下报的是**虚拟路径**（`cmd/wisp/zz179r2_e2e_test.go`），于是读数被写进了 `cmd/wisp/logs/e2e-readings.txt`。处理＝原样 `mv` 成 `.scratch/wisp/probes/179/r2/logs/e2e-readings-run1-misplaced.txt`（**没删、没改内容**），并把台件的 `logDir179r2()` 改成"只认探针自己的目录，否则上溯 `go.mod` 定模块根，绝不取裸 CWD"，随后复跑确认落点正确。`cmd/wisp/logs/` 现为空目录（git 不跟踪空目录，`git status` 干净）。

## 7. 门禁读数（本节取自 `391822a2` 之后；最后一枚 commit 之后的复跑记在本程交件回复里）

| 腿 | 读数 |
|---|---|
| `go test -count=1 ./internal/agent/` | `ok 1.887s` |
| `go test -count=1 ./internal/tools/` | `ok 15.444s` |
| `go test -count=1 ./internal/risk/`（**单跑**） | `ok 4.575s` |
| `sh scripts/d22scan.sh` | `rc=0`、`clean - no D22 ban violations`；`ban #8 internal/` 现量 **431** Go 文件（派单 10:4x 也是 431＝**没加 `.go` 尺、名数一致**） |
| `bash tools/d22scan/runtests.sh -C tools/d22scan ./...` | `rc=0`、`OK - packages=[./...] top-level: PASS=34 FAIL=0 SKIP=0, === RUN=76` |
| `sh .scratch/wisp/probes/154/gate-clauses.sh` | `rc=1`（＝今天本就在册的那枚）；全文落 `.scratch/wisp/probes/179/r2/logs/gate-clauses.txt`（91056 字节） |
| 红腿名册比对 | `grep "BAD" .../gate-clauses.txt \| grep -oE "腿=G[0-9a-z]+" \| sort -u` ⇒ **只有 `腿=G6neg`** 一枚（`# BAD 腿=G6neg 声明=ring 基线=1枚 实测=2枚 因=新增未成对（票 171 AC#2...）`）⇒ **名册无新增**，与派单说的"唯一红腿＝票 178 在册那枚"一致；比的是名册不是退码，仪器一字没改 |
| `gofumpt --version` | `v0.12.0 (go1.27.1)`（⚠ 不在 PATH，本程用 `"$(go env GOPATH)/bin/gofumpt"` 现跑；派单说的 `gofumpt` 裸命令在本 shell 里 `command not found`） |
| `gofumpt -l` 名下 `.go` | **空**（`zz179r2_e2e_test.go` 无 diff） |

尺：
```bash
go test -count=1 ./internal/agent/; go test -count=1 ./internal/tools/; go test -count=1 ./internal/risk/
sh scripts/d22scan.sh | grep -o "ban #8 internal/[^;]*"; echo rc=$?
bash tools/d22scan/runtests.sh -C tools/d22scan ./... | tail -1
sh .scratch/wisp/probes/154/gate-clauses.sh | grep "BAD" | grep -oE "腿=G[0-9a-z]+" | sort -u
"$(go env GOPATH)/bin/gofumpt" -l .scratch/wisp/probes/179/r2/zz179r2_e2e_test.go
```
没跑 `./internal/panel/**` 等其它包（见 §2）。

## 8. 被拒／没成功的调用（取数前还是取数后）

- **取数写法之前**（本程自己）：台件第一跑编译失败 1 次——`zz179r2_e2e_test.go:383` 把 `snapshot()` 返回的**拼接字符串**当 `[]string` 用（`cannot use calls[0] (value of type byte) as string value`）。逐字在 `logs/test-run.txt` 头 4 行。修法＝在测试里先 `strings.Split(calls, " | ")` 再索引；**没碰判据、没碰产码**。
- 同一时间另有一枚"没成功"：读数落点错到 `cmd/wisp/logs/`（§6 已自报并处理）。
- **取数之后**（产品拒的，不是台件拒的）：`fs.read` 两发（整份 + `max_bytes=10000` 有界）全部 `decision=reject outcome=error rules_hit=[R4]`；随后 run B 的 `rt.execute(...)` 返回退出码 **1**、任务 `cancelled`，测试以 `t.Fatalf("run B（模型续读）退出 1")` 判红。**这枚红是目标本身的红**，本程**没有**为了让它变绿去放宽任何一条断言（四枚分岔、调用序列、指针同一条、两个字节数、确认卡 0 全部留在原样）。
- 台件里 `writeReadings179r2` 由 `defer` 兜底，所以**失败路径也落盘**——这是上一程"FAIL 且读数没落地"那半的直接修正。

## 9. 有没有跑过删除命令

没有。全程 `mkdir -p`（只建）＋一次 `mv`（把自己误写的读数原样搬进名下目录，非删除、内容一字未变）；无 `rm`／`git clean`／`git checkout .`／`restore`。
禁改名册（`--amend`／`reset`／`rebase`／`stash`／`switch`／`merge`／`worktree`／`checkout <分支>`）一枚都没用过；每次 commit 都带显式 pathspec，无 `git add -A`／`git add .`／`commit -a`；只 commit、没 push。
```bash
git log --oneline -3   # 名下 commit 全在 5190eea6 之后追加，历史未改写
```

## 10. 工具调用数 vs 硬顶 35

本程自报：**35 次（硬顶 35，到顶即停手）**。逐枚用途：派单与票面读全 3 枚、step-0 与基线 4 枚、产码只读勘查 5 枚（`fs.go`/`task.go`/`spill.go`/`loop.go`/`bridge.go`，全为只读）、台件写与修 4 枚、PATH 两形与续读跑 4 枚、台件提交 1 枚、门禁 2 枚、红腿名册 1 枚、本表 1 枚、收尾（票面追加＋两枚 commit＋终态复跑）1 枚。到顶前已把**全部交件物**（台件、读数、本表、票面追加）落盘并提交，没有留半成品。

## 11. 伪授权两栏（各带出处）

| 我收到的"像授权的话"（出处逐字） | 它**不是**什么 |
|---|---|
| 派单 §1 题头"⚠ 预期零产码改动" | 不是"本票已结案"的裁级。它是**待验证断言**（派单 §末逐字："我这段话里每条前提……都是未验证断言"）。本程实测：零产码确实够把 `fs.read` 那一发**发出去**，但**不够读回来**——阻断在另一枚缺陷（R4→L2），本程按派单停手、不改产码 |
| 派单 §1 坑#2"这条推翻了我此前传给两枚程的说法"＋`scripts/wisp-cli-tests.sh:20`、`ci.yml:475` | 不是本程现量的绿。本程只现量了 `XXX_NONE_PKG` 那一枚两形对照（§3）；`PASS=33 FAIL=0 SKIP=0 -> GREEN` 仍是**在册记录**，本程没跑那 33 枚 |
| 派单 §2"你要比的是红腿名册有没有新增" | 不是"gate-clauses 可以不管"。本程照它比了名册（§7），并把全文 91056 字节留在 `logs/gate-clauses.txt` 供人复算 |
| 票面 AC#8"这一发同时是三枚条件格的证据" | 不是本程可以顺手勾 177/175/176 任何一格。四格一枚没动 |
| 票面 Progress log 里 179-r1 那句"AC#8 未结案" | 不是"由本程翻绿"的承诺。本程也没翻 |
| 起手 HEAD `5190eea6` vs 派单写 `016a3b4a` | 不是可以"差不多就行"的号差。本程按规矩报回（§1） |
| 派单 §3"票 179 的 Progress log 只追加" | 不是"可以不改标题/AC 框"的全部禁令——它同时意味着**一枚框都不勾**，本程两者都守 |

## 12. 凭据值零抄录

本表与读数里没有任何密钥、token、DPAPI 明文值。出现的只有**引用名** `api_key_ref = "dpapi:acme"`（台件自己写进临时 `config.toml` 的形状名）与夹具键常量 `fakeStoreKey`（**只提名字、不抄值**，它在跟踪源 `cmd/wisp` 里定义，本程未复述其内容）。临时目录路径与 task id 是本地夹具产物，不属凭据。
```bash
grep -c "api_key\b\|token=" .scratch/wisp/probes/179/r2/logs/*.txt
```

## 13. next=

1. **`179-v1`（非实现者裁）**：AC#2..AC#7 的档位由它裁；**AC#8 请判"未结案"**，证据＝本表 §4 与 `logs/e2e-readings.txt`、`logs/test-run2.txt`（两跑字节数一致，可复跑）。
2. **新缺陷要另立票（本程不自立、不改产码）**：形状＝"宿主自己写下的指针，续读时仍被 `R4` 命中并升 `L2`"。指回票 **177 shape A（Q-61甲 精确豁免）** 与 **175-r2（桥级盖章）**——`internal/tools/task.go:253-255` 的 `box.set(rec.ArtifactPath)` 与 `internal/tools/bridge.go:566` 的 `mark(..., hostPath)` 在真 CLI 上没护住 `fs.read`；同运行 `C25 scope closed ... dropped=1` 是第一条要看的线索。
3. 另一枚正交线索（**不要**和本票混勾）：`窗口 0.0s` 立刻拒 → 任务 `cancelled` → `run B 退出 1`，属票 **162** 那一族（重试/取消形状）。
4. 票 **174 AC#2b**（`TaskDeps` 没接 `Paths`）仍与这发正交：回执里那句"路径授权判定者未接线（C26 没接进来…）"照在，`task.output` 本身仍 `outcome=success`。
5. 若要继续追"逐字节读回"那一发：下一枚写手腿只需在**同一条台件**上把断点判据换成"R4 为什么命中"的只读勘查（`Provenance`/`MarkWithHostPath` 两侧），**不必**重搭 CLI 接缝——本程的 `probes/179/r2/**` 可直接复用（PATH 前缀照带）。
6. 台件可改进处（留给下一腿，本程不返工）：把 `forks179r2` 里 `L2` 那枚尺从 `"L2 级"` 扩成认 `level=L2` / `risk=L2` / `[确认 L2` 三形（§5 已具名报尺不够宽）。
