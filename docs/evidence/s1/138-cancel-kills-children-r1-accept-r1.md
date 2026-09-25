# 138 对抗验收 r1 — 独立重走"预留的洞"那枚裁决 ＋ 把 §6.6 那三句"简报说"回核原文

- 工单：`.scratch/wisp/issues/138-cancelling-a-task-kills-nothing-the-task-itself-spawned-single-jobscope-no-terminatejobobject.md`
- 被验件：`docs/evidence/s1/138-cancel-kills-children-r1.md`（875 行，§0..§6.8；§7 起未写＝它撞在 150 轮上限）
- 本程：**对抗验收程（非实现者）**，只读；唯一写件＝本文件
- 被验版本＝`3f6322f`；派单给的裁决问题：它说本票要防的洞"今天不存在"，成立则按"作废并登记"处理，不成立则修法还得做
- 本文件**不复算它的结论**，逐发复算它的**命令**；每格另附「本程没测什么」与「结论修正记录」

---

## 0　进场：锚点、工作树、脏树隔离

**判据**（派单 §0）：读数按 sha 取，不许把脏工作树当被验版本。

### 0.1　命令原文与读数

```
$ git rev-parse HEAD
3f6322f1b6a5d77fac4e17dd0fa7ea33c19d3d53

$ git log --oneline -14
3f6322f evidence(138 落盘代理未提交的尾段): 它把"锚点漂移发生在哪两格之间"改准了，并就地自认一次口误
2f5c7f9 probes(149): 变异尺与逐发原始读数落盘（57 枚，只建不删）
f8623c7 evidence(138 AC#4 门禁 第 6 格): 两形四数+名册两向差集+同包开跑/裁决等集；…
b417d31 fix(149 AC#2+AC#3+AC#4): corrupt 支的 offset 改填解码器自己点名的位置，并补三枚会响的用例
c4540a6 evidence(138 AC#2/AC#3 处置 第 5 格): 条件句前件判假故未落修法，…
69aefae evidence(138 AC#1 总裁决 第 4 格): 判为"给未来功能预留的洞"（三发支撑…）
d8e3957 evidence(138 AC#1③ 第 3 格): 生产 exec 路径逐名七枚、任务路径六枚 fs.* 零子进程；…
1ff7c71 evidence(138 AC#1② 第 2 格): AssignProcessToJobObject 普查…
c9237a9 evidence(138 AC#1① 第 1 格): JobScope 创建点数与调用者现量…
2fb80c2 docs(台账 A266): 第二推逐名比红名集合又是零变化；…
64858d6 docs(台账 A265): 票 147 结线五格全成立；…
de1ee7b ticket(147 结线＋144/147 改名＋149 新立): …
e9412c9 evidence(147 对抗验收 r1 第 9 格): …
1642886 evidence(147 对抗验收 r1 第 5-8 格): …

$ git status --porcelain | wc -l
79
```

**读数**：HEAD＝派单给的 `3f6322f`，分支 `dev`。工作树脏 **79 行**：`design/**`（16 枚 owner 自己的删除）、
`frontend/**`（19 枚 ` M` ＋ 15 枚未跟踪）、`docs/reports/frontend-session-log-zcode.md`、
`.zcodeignore`、`.scratch/wisp/probes/149/*.log`（5 枚未跟踪）、`docs/evidence/s1/149-three-unpinned-outlets-r1.md`（未跟踪）。
⇒ **两枚在飞写者可见**：票 149（`cmd/wisp/**`）、前端会话（`frontend/**`＋`design/**`）。本程一枚未碰。

### 0.2　被验两文件的"工作树 vs 锚点"同形自证（否则我读的是别人正在改的东西）

```
$ T=.scratch/wisp/issues/138-cancelling-a-task-kills-nothing-the-task-itself-spawned-single-jobscope-no-terminatejobobject.md
$ git diff --stat 3f6322f -- "$T" docs/evidence/s1/138-cancel-kills-children-r1.md
（空输出）
$ git show 3f6322f:"$T" | wc -l
78
$ git show 3f6322f:docs/evidence/s1/138-cancel-kills-children-r1.md | wc -l
875
```

⇒ 票面 78 行、被验件 875 行，**两枚都与锚点逐字节同**，本程之后按路径读文件＝按 sha 读，无需快照。

### 0.3　锚点到被验版本之间动过的**代码**（决定我哪些格可比）

```
$ git diff --stat 64858d6..3f6322f -- internal cmd tools scripts go.mod go.sum
 cmd/wisp/slo_report_144_windows_test.go | 268 ++++++++++++++++++++++++++++++++
 cmd/wisp/slo_windows.go                 |  60 ++++++-
 2 files changed, 325 insertions(+), 3 deletions(-)
```

**读数**：实现程锚点（`64858d6`）与我的被验版本（`3f6322f`）之间**只有票 149 的两枚 `cmd/wisp` 文件动过**；
`internal/**`、`tools/**`、`scripts/**`、`go.mod`、`go.sum` **一枚未动**。
⇒ 它的三发支撑（全部落在 `internal/**` 与 `internal/tools` 的注册表上）在我这一版**逐字可复算**，
唯一要重算的是 `cmd/wisp` 侧那两发（`agent.New` 的调用者、`run.go` 的注册段）。

### 0.4　仪器版本（现跑，不背读数）

```
$ go version
go version go1.27.1 windows/amd64

$ D:/work/base/gopath/bin/gofumpt.exe --version
v0.12.0 (go1.27.1)
```

⇒ 与票面 AC#4"宿主已装 v0.12.0"相符；本程未执行任何 `go install`。

### 0.5　本程没测什么（按"漏了它谁会先被骗"排序）

1. **没在别的分支上取数**：全部读数来自 `dev`@`3f6322f` 的跟踪文件。若有未合并分支已落 `shell.exec`，
   第 2 格支撑①的"今天无入口"当场不适用（这一条实现件 §4.5①自己写了，我没有反驳它，只是也没复算）。
2. **没做运行时观测**（跑一次真任务数子进程）：`cmd/wisp` 此刻是 149 的写者，构它取到的是中间态。
3. **没验 git 之外的凭据面**：本程不需要密钥，故未做任何凭据反扫。

### 0.6　结论修正记录（本格）

- 派单说"被验版本＝`3f6322f`"、"另一枚写码程在 `cmd/wisp/**` 跑 149"——两枚**成立**；
  但派单没提：**149 的件已经 commit 进 `3f6322f` 之前了**（`b417d31`、`2f5c7f9`），
  也就是被验版本里**含票 149 的实现**。⇒ 本程把"哪几发要重算"写进 §0.3，而不是默认它与 149 无关。

### 0.7　本程在场期间 HEAD 被推走（如实记，与 §1.3 那枚循环引用同源）

```
$ git log --oneline 3f6322f..HEAD
4d0866a evidence(139 AC#1 第 1 格): 零留痕量成读数…
c3f7224 evidence(149 第 12 格): 入库字节再跑一遍三发关键变异…
a055d9f evidence(149 第 0-11 格): 三发今天读数…
612590c docs(台账 A267): 票 138 撞顶＝落盘不是重派；它引用了三句我从未说过的"简报说"…
$ git log --format='%h %ci %s' -1 612590c
612590c 2026-09-25 23:3x +0800 docs(台账 A267): …
```

**读数**：本程进场读到 `3f6322f`，我建件时 HEAD 已是 `4d0866a`——**四枚 commit 插进来**（A267 台账一枚、
票 149 验收两枚、票 139 实现程第一格一枚）。逐枚 `git show --name-only` 数过：
路径集＝`docs/reports/pending-and-issues.md` ＋ `docs/evidence/s1/149-*.md`／`139-*.md` ＋ `.scratch/wisp/probes/{149,139}/*`，
**代码面 0 枚**——现量复算：

```
$ git diff --stat 3f6322f..HEAD -- internal cmd tools scripts go.mod go.sum
（空输出）
```

⇒ **被验版本不变，本程全部读数仍按 `3f6322f` 取**；
但这一发直接导致 §1.3 的坑：**编排者否认这三句的那格台账（A267）在我进场之后才入库**，
工作树 grep 会把它当成"盘上另有一处记载"。已在 §1.3 拆干净。

---

## 1　派单第 3 节：把 §6.6 那三句"简报说"回核原文（本程重点）

**判据**（派单 §3）：那三句按字面 `grep` 票面与派单；§6.6 附近它引用的每条编号／文件名核存在性；
给出结论（甲 自造引用／乙 工具输出混入／丙 改写票面当"简报说"），**核不到就写"未定性"并说明缺哪一枚读数**。

### 1.1　第一发：票面 138 里 `go build ./...` 出现几处、怎么写的

```
$ git grep -n -F 'go build ./...' 3f6322f -- '.scratch/wisp/issues/138-cancelling-a-task-kills-nothing-the-task-itself-spawned-single-jobscope-no-terminatejobobject.md'
3f6322f:.scratch/wisp/issues/138-cancelling-a-task-kills-nothing-the-task-itself-spawned-single-jobscope-no-terminatejobobject.md:56:      先证变异落地（`grep -n` 到你改那一行的原文 ＋ `go build ./...` rc=0）再读数。
rc=0（1 枚命中）
```

⇒ **派单那句成立**：票面里 `go build ./...` 只有第 56 行一枚，且**写作 rc=0**（AC#3"先证变异落地"那半句）。
把它引成"简报说 `go build ./...` rc≠0"是**符号都反的**（AC#3 要的是"改了码之后还能构干净"）。

### 1.2　第二发：票面 138 里有没有任何 CI / workflow / runner 字样（比派单更强的一发）

```
$ git grep -n -E 'CI|workflow|runner|gh run' 3f6322f -- '.scratch/wisp/issues/138-cancelling-a-task-kills-nothing-the-task-itself-spawned-single-jobscope-no-terminatejobobject.md'
rc=1（零命中）
```

⇒ **票面整份文件零枚 CI／workflow／runner／`gh run` 命中。**
所以"CI 那两枚红的成因是别人的代码"、"gofmt／go vet／d22scan 三步在 CI 是绿的"这两句
**不可能**是"改写票面某句"（假设丙）的产物——票面根本没提 CI。派单那句"没有任何一处说 CI 那两枚红的成因"
**成立，而且成立得比它写的更彻底**。

### 1.3　第三发：那三句在**整棵锚点树**里的出处（＋一枚我自己差点踩的循环引用）

```
$ git grep -n -F '三步在 CI 是绿的' 3f6322f
3f6322f:docs/evidence/s1/138-cancel-kills-children-r1.md:848
$ git grep -n -F '成因是别人的代码' 3f6322f
3f6322f:docs/evidence/s1/138-cancel-kills-children-r1.md:816
$ git grep -n -F 'CI 那两枚红' 3f6322f
3f6322f:docs/evidence/s1/138-cancel-kills-children-r1.md:816
3f6322f:docs/evidence/s1/138-cancel-kills-children-r1.md:846
$ git grep -n -F '别人的代码' 3f6322f | grep -v 138-cancel-kills-children-r1.md
（空输出）
$ git grep -n -F 'Build all packages' 3f6322f | grep -v 138-cancel-kills-children-r1.md
（空输出）
```

⇒ **在被验版本 `3f6322f` 里，三句只有一枚出处：被验件自己。**

⚠ 但**工作树里多一枚**，而那枚**不能用**：

```
$ grep -rn -F '成因是别人的代码' . --include='*.md' | grep -v 138-cancel-kills-children-r1.md
./docs/reports/pending-and-issues.md:6567:  **③ 本条真正要留的：它 §6.6 摆了一张……那三句"简报说…"我回核不到出处。**
$ git log --oneline -S'成因是别人的代码' -- docs/reports/pending-and-issues.md
612590c docs(台账 A267): 票 138 撞顶＝落盘不是重派；它引用了三句我从未说过的"简报说"…
$ git grep -c 'A267' 3f6322f -- docs/reports/pending-and-issues.md
rc=1（A267 不在被验版本里）
```

⇒ 那枚命中＝**编排者自己 A267③ 那一格**，它落在 `612590c`（`3f6322f` 的**后代**），
内容是"**我核不到这三句的出处**"——**它证的正是本案要判的那件事**。
⇒ 拿它当"盘上有记载"就是**循环引用**：只有 A267 之前的树才配当 witness。本程第一版草稿把这两枚出处混在一起了，
已在 §1.10 记一笔并改在这里。
⇒ 结论不变、且更硬：**盘上没有任何一份"三句被说过"的原文。**

⇒ **读数**：被验版本里三句**只有被验件自己这一枚出处**；台账那枚是后代的否认格、不作 witness。
⇒ 盘上没有任何一份"这三句被说过"的原文。

### 1.4　第四发（决定本格能不能定案的一发）：**派单正文在盘上不存在**

```
$ ls -la .scratch/wisp/dispatches/
total 4
-rw-r--r-- 1 swq 197609 3587 Sep 24 16:03 README.md
$ sed -n '3,10p' .scratch/wisp/dispatches/README.md   （取该目录自己的规矩）
> **为什么有这个目录**：`A162` 那笔账的根因——本仓此前**没有派单归档制度**……
> **从这里起，凡授予范围／解冻／例外处置的派单，正文里的授权段落在派完之后落进本目录。**
> ⚠ **规矩只有三条**：①**不回填旧的**（回填＝造新史料，此前派单一律按"不可核"处理）；②**只建不删、只追加**；
```

⇒ 归档制度只管**"授予范围／解冻／例外处置"**那一类派单；票 138 是**普通实现派单**，按制度**不入盘**。
⇒ 于是"派单里到底有没有这三句"**在盘上没有 witness**——
**不是"我核到它没有"，是"这一类问题本仓今天无法核"**。
⚠ 顺带：`git grep -n '138' docs/reports/HANDOVER.md` 的 6 枚命中全在 09-21～09-24 的旧停车点，
HANDOVER 最新一节仍是 `4.0f`（09-21 20:3x）⇒ 停车点里也没有派单正文。

### 1.5　第五发：§6.6 那张表的**事实层**我独立复算（这一层与归属层无关，先钉住）

```
$ gh run view 36149256584 -R CarlosShao/wisp --json headSha,conclusion -q '.headSha, .conclusion'
64858d6838ced46fbe7bcce38dc4f7bb163d2f9c
failure

$ gh run view 36149256584 -R CarlosShao/wisp --json jobs -q '.jobs[] | .name + " => " + (.conclusion//"-") + " | " + ([.steps[] | select(.conclusion=="failure") | .name] | join(" ;; "))'
lint         => failure | staticcheck
slo-full     => success |
test-core    => failure | Portable package tests (core scope; the list and its guards live in scripts/portable-tests.sh)
test-windows => failure | cmd/wisp CLI tests (needs the sherpa DLLs staged above, ticket 111 AC#4) ;; Portable windows tests (proc/secret/config/risk/ball/perm/plugin/llmrecord)
lint-frontend=> success |
slo-smoke    => success |

$ gh run view 36149256584 -R CarlosShao/wisp --json jobs -q '.jobs[] | .name as $j | .steps[] | select(.name | test("Build all")) | $j + " | " + .name'
（空输出）

$ gh run view 36149256584 -R CarlosShao/wisp --json jobs -q '.jobs[] | select(.name=="lint") | .steps[] | .name + " => " + (.conclusion // "-")'
D22 scanner positive control (tools/d22scan tests, seeded red) => success
D22 seven-ban + emoji scan (tools/d22scan) => success
gofmt (gofumpt) => success
go vet (module) => success
go vet (tools/d22scan module) => success
staticcheck => failure
mockllm module vet => success

$ git show 3f6322f:.github/workflows/ci.yml | grep -n 'go build'
295:  # cgo build smoke (clean-clone + build.ps1, S0 done criterion) + the
349:        # not sit behind `cgo build smoke`, which has been red before) and it is
365:        # red before (`cgo build smoke`, run 35551819606) plus a red `PathResolver
395:      - name: "cgo build smoke (build.ps1 fetch-deps + mingw link + doctor)"
549:      # build.ps1 -> "go build ok" + doctor all PASS + rc=0.
```

**读数**：与我逐字相同的＝被验件 §6.6。三处事实层**我独立复算成立**：
① 锚点那 run 确实 `failure`，红在 **3 枚 job／4 枚步骤**，步骤名逐字同；
② 这份 workflow **没有任何一步跑 `go build ./...`**（`grep 'go build'` 五枚命中里 4 枚是注释、
一枚是步骤名 `cgo build smoke`，**无 `run:` 行**）；**也没有名为 "Build all packages" 的步骤**（空输出）；
③ `lint` 那 job 里 `gofmt (gofumpt)`／`go vet (module)`／`go vet (tools/d22scan module)`／`D22 seven-ban + emoji scan` **全是 success**，
唯一红的是 `staticcheck`。
⇒ 它那格的"现量"列**没有造假**，`staticcheck` 那枚"简报没点的"也是真的。

### 1.6　第六发：§6.6 附近它引用的每条编号／文件名的存在性（逐名）

| 它引用什么 | 我的现量命令 | 存在？ |
|---|---|---|
| `docs/specs/SPEC-12-roadmap-governance.md:61`＝DEFERRED/RESERVED 类型语义 | `git grep -n '类型语义' 3f6322f -- docs/specs/SPEC-12-roadmap-governance.md` ⇒ `:61` | **存在，行号逐字对** |
| SPEC-01 §6 层级句 `SPEC-01-architecture.md:196` | `git grep -n '会话 scope' 3f6322f -- docs/specs/SPEC-01-architecture.md` ⇒ `:196 - 层级：会话 scope（C31）⊃ 任务 scope ⊃ 工具 scope ⊃ 插件 scope。` | **存在，行号逐字对** |
| `internal/plugin/disposal.go:155 NewDisposalScope` | `git show 3f6322f:internal/plugin/disposal.go \| sed -n '150,160p'` ⇒ `func NewDisposalScope(…)` 正在 `:155` | **存在** |
| `internal/config/unwired.go:63` 那条守卫 | `git grep -n 'no shell.exec tool is registered' 3f6322f -- internal/config/unwired.go` ⇒ `:63` | **存在**（票面写的 `internal/tools/unwired.go` 才是打偏的那枚） |
| `internal/tools/unwired.go`（票面点名） | `git ls-files internal/tools/unwired.go` ⇒ 空 | **不存在**——它已如实登记为票面缺陷 |
| `internal/proc/crossvet_test.go:52 TestCrossVetForLinux` | `git show … \| sed -n '50,54p'` ⇒ `func TestCrossVetForLinux` 在 `:52` | **存在** |
| `scripts/d22scan.sh` | `git ls-files scripts/d22scan.sh` ⇒ 在 | **存在** |
| run id `36149256584` | §1.5 `gh run view` ⇒ 该 run 存在且 `headSha`＝`64858d6` | **存在** |
| `docs/evidence/s0/03-adversarial-acceptance.md`、`.scratch/wisp/issues/03-…-done.md`、票 21／50／51 | `git ls-files`／`ls .scratch/wisp/issues/` | **全部存在** |
| 票面 AC#4 那枚 ⚠（"Linux 交叉 vet 对 `cmd/wisp` 与 `cmd/balldebug` 本来就 rc=1"） | `git show 3f6322f:<票面> \| sed -n '62,63p'` | **存在**，且它在 §6.5 明确**拒绝**把这条读成 CI 那一发的成因 |

⇒ **引用层（可盘核的那一半）零伪造**：它能落到盘上的每一枚编号／路径／行号我逐名核过，行号级都对得上。
唯一不存在的那枚（`internal/tools/unwired.go`）是**票面的错**、不是它的错，且它当场登记了。

### 1.7　第七发：它 §6.6 贴的那条命令与它贴的输出**不同形**（本程自己造出的一发，派单没要求）

```
$ grep -n 'name:' .github/workflows/ci.yml | grep -i build
395:      - name: "cgo build smoke (build.ps1 fetch-deps + mingw link + doctor)"
496:      - name: Build wisp.exe
559:      - name: Build wisp.exe (deps cached on the runner)
642:      - name: build (vite build -> frontend/dist, the bytes go:embed carries)
$ grep -n 'name:' .github/workflows/ci.yml | grep -ic build
4
```

⇒ 它 §6.6 贴的**同一条命令**只给了 3 行输出，**盘上实为 4 行**：漏的第 4 枚是 `ci.yml:642`
`build (vite build -> frontend/dist, …)`，属 `lint-frontend` 那枚 job（前端产物构建，与 Go 构建无关）。
⇒ **影响判定**：不改它那条结论（"没有 `go build ./...` 那一步"照样成立，第 4 枚也不是 Go 构建），
但**"贴出的命令与贴出的输出必须同形"是本仓自己的纪律**——这一格是输出被截了一枚。
⚠ 同时：**派单（我这单的上级）写"三枚带 build 字样的步骤名"，实为 4 枚命中／3 枚 Go 形状**——
这句按"数法未指明"记入 §1.9 修正记录。

### 1.8　本格裁决

**事实层：成立。** §6.6 的"现量"列我逐发独立复算（§1.5），无造假；
`staticcheck` 在锚点是红的、workflow 里没有 `go build ./...` 那一步、没有 "Build all packages" 那枚步骤，三条**都是真的**。

**归属层：未定性。** 三种假设的现有强度（我手上的读数只到这里）：

| 假设 | 我这一程能给的证据 | 还缺哪一枚读数 |
|---|---|---|
| (丙) 改写票面当"简报说" | **只能解释一小半**：票面 AC#4 那枚 ⚠（"本来就 rc=1、与被审对象无关"）形状相近，但它在**本机交叉 vet** 上、不写 CI、不写"两枚"；AC#3 的 `go build ./...` 写作 **rc=0**，符号相反；票面 AC#4 的 gofmt/vet/d22scan 是**要它自己跑的活**，不是"CI 绿"的断言。且票面 **CI 字样 0 命中**（§1.2） | 不缺，这一支**证据不足** |
| (甲) 自造引用 | 盘上零 witness（§1.3）＋它能核的引用一枚没造假（§1.6）⇒ "惯于编造"的形状**没有旁证** | 该程转录里这三句**首次出现的位置**（它的上下文，只有平台侧有） |
| (乙) 工具输出混入（第 8 代那一类：假凭据） | 该程**没写**这三句出现在哪一枚工具结果里；本仓 §5.3 的登记口径要求"原文＋计数＋出处（工具名＋命令前 40 字）"，**这一栏整段缺失** | 同上：只有转录能判 |
| **(丁) 编排者确实说过近似的话**（派单没列这一支） | 台账 `A265④` **自己写着**"凡是'某读数的成因是 X'这种因果句……我上一轮把它当成事实写进了派单"，同一天 `A265⑤`／`A264` 都在报 CI 步级颜色；`HANDOVER.md:924` 另有一条在盘上的近似事实："cmd/wisp 的 sherpa cgo 使 `CGO_ENABLED=0 go build ./...` 在 cmd 失败" | **派单正文本身**（盘上无 witness，§1.4） |

⇒ **本格不判"代理编造引用"**，理由不是替它说话，是**这一类 dispute 在本仓今天不可判**：
派单不入盘（`dispatches/README` 规矩①的适用范围是"授予范围／解冻／例外"那一类），
于是"简报说 X"这句话**永远只能靠编排者自述**——而 `A265④` 已经记过两次"编排者自述的前提是错的"。
⇒ **可落地的处置只有一枚**：`A267③` 那张表的归属结论**不能由任何一方单方定**，
要么补制度（把"含事实断言的派单正文"也归档），要么把这枚 dispute 摆成 owner 的 `Q##`；
把它写成"代理自造引用"是一枚**没有凭据的定罪**，写成"注入"是**另一枚**（第 8 代那一形有硬指纹：
假 sha／假 token 串，这三句**没有任何凭据形状的字符串**，与那一代不同形）。

**档位：本程读数成立；对被验件 §6.6 那一格＝附条件入账。**
条件（**由编排者落笔，本程不改它一字**）：那张表的"简报的说法"列改标〔出处未定·盘上无 witness〕，
"判定"列保留事实层的"成立／对不上"。

### 1.9　本程没测什么（按"漏了它谁会先被骗"排序）

1. **没拿到派单原文**（盘上无 witness）。⇒ 谁先被骗：**读 `A267③` 就把"三句都不存在"当已结案事实的下一程**——
   它会把一枚未定性的归属当成已定性的"代理造假"，从此在被验件的表上打折。
2. **没读该程的转录**（平台侧才有）。⇒ (乙) 既不能立也不能废；我只能报"它没按本仓 §5.3 的口径写出处"。
3. **没核 §6.6 之外那 6 处简报引用的真实性**（`:6`／`:13`／`:74`／`:538`／`:539`／`:619`／`:637`／`:756`／`:864`）。
   其中 `:6` `:13` `:74`（锚点 sha、落点 `internal/agent/**`）与 `:619`（"这活的本体是杀掉被取消任务自己起的子进程"）
   我**顺手核过＝与票面 §1 同形**；其余 5 处同属派单侧、同样无 witness，**没逐条展开**。
4. **没去查那三枚红在 CI 里的日志正文**（`staticcheck` 红在哪一行、portable 红在哪一枚用例）。
   它与"这三句是谁说的"无关，实现件 §6.8④ 自己也标了没做。
5. **没测 (丁) 假设的下限**：没去比对编排者**当天推送给 CI 的那两推的 commit message** 里有没有近似句子——
   那枚读数存在，但只能证明"编排者说过近似的话"，仍证不了"派单里说过"，我没做。

### 1.10　结论修正记录（本格）

0. **推翻本程自己第一版草稿一处（未 commit 前自纠）**：§1.3 初稿写"三句在盘上有两枚出处：被验件＋台账 `A267③`"，
   并把那枚台账命中算进了 `git grep … 3f6322f` 的输出里。**实际：A267 是 `612590c`，是 `3f6322f` 的后代，
   在被验版本里 `git grep -c 'A267' 3f6322f` ⇒ rc=1（零枚）**；那枚命中只存在于工作树。
   ⇒ 形状＝**"把后代的转述当成被验版本里的记载"**，与本仓"行号带版本"那条同源，只是这次搬的是**台账条目**。
   已在 §1.3 按两发分开重贴命令与输出。
1. **推翻派单一处（小）**：派单写"三枚带 build 字样的步骤名你自己贴出来"。现量**同一条命令给 4 枚**
   （多 `ci.yml:642` 的 `build (vite build -> frontend/dist …)`）。派单那句只在"只数 Go 形状的构建"下成立。
2. **加强派单一处**：派单说"票面没有任何一处说 CI 那两枚红的成因是别人的代码"。
   现量更强：**票面整份零枚 `CI`／`workflow`／`runner`／`gh run` 字样**（§1.2）。
3. **推翻派单隐含的一处**：派单把可判的假设列了 (甲)(乙)(丙) 三支。现量**必须补第 (丁) 支**
   （编排者当日有"把因果句当事实写进派单"的两次自记，`A265④`），否则本格的"未定性"会被读成"只剩甲乙"。
4. **没推翻任何一处自己的**：本格 §1.5 那三发我第一次就复算到与被验件逐字同，没有中途改口。

---

## 2　三发支撑独立重走（"预留的洞"那枚裁决的承重点）

**判据**（派单 §2）：三发逐发独立重走，**不复算结论、复算命令**；
⚠ **只要有一发不成立**（例如某条生产路径今天其实能 exec），"预留的洞"这个裁决就塌，本票要改回"已实现功能里的缺口"。

### 2.1　支撑① —— "任务侧起不了子进程"：我自己从 `run.go` 走到名册，再换一把更宽的尺

**(a) 组合根只有一处注册工具，且那处只注册 `fs.*`**（149 已把 `cmd/wisp` 动过两枚文件，所以这一发我按 `3f6322f` 重取）：

```
$ git grep -n '\.Register(' 3f6322f -- '*.go' ':!*_test.go'
3f6322f:cmd/wisp/run.go:345:            if err := reg.Register(e); err != nil {
3f6322f:tools/mockllm/main.go:51:       srv.Register(mux)
$ git grep -n 'Builtin[A-Za-z]*Entries' 3f6322f -- '*.go' ':!*_test.go'
3f6322f:cmd/wisp/run.go:341:    for _, e := range tools.BuiltinFSEntries(tools.FSDeps{
3f6322f:internal/tools/fs.go:323:func BuiltinFSEntries(d FSDeps) []Entry {
3f6322f:internal/tools/fs_write.go:766:func BuiltinFSWriteEntries(d FSDeps) []Entry {
```

⇒ 生产码里 `Registry.Register` 的调用者**只有 `run.go:345` 那一枚**（`tools/mockllm` 那枚是 HTTP mux，不是工具注册表）。

**(b) 环路拿到的 Tools 就是那枚桥**：

```
$ git show 3f6322f:cmd/wisp/run.go | sed -n '547,550p'
        loop, err := agent.New(agent.Options{
                Provider: prov,
                Tools:    rt.bridge,
                Sink:     consoleSink{out: rt.stdout, stream: rt.stream, publish: rt.publishPanelSnapshot},
$ git grep -n -A28 'type Options struct' 3f6322f -- internal/agent/loop.go
internal/agent/loop.go:149:   // Tools is the C4 surface ...
internal/agent/loop.go-150-  Tools ToolProvider
```

⇒ `agent.Options` 的全部字段里**没有任何 Job／scope／进程属主入口**（只有 `Tools`、`Registry *observe.Registry`、
`AdmitTask`、`Journal`、`Sink`、`Summarizer`、`Control`、`Logger`、`Config`）。

**(c) 名册逐名**（我不引它的表，自己从 `Name()` 抽）：

```
$ git grep -n -A2 'func .*Name() string' 3f6322f -- 'internal/tools/*.go' ':!*_test.go'
fs.go:122        func (fsRead) Name() string  { return "fs.read" }
fs.go:188        func (fsList) Name() string  { return "fs.list" }
fs_write.go:233  func (fsWrite) Name() string { return "fs.write" }
fs_write.go:368  func (fsTrash) Name() string { return "fs.trash" }
fs_write.go:439  func (fsMove) Name() string  { return "fs.move" }
fs_write.go:582  func (fsDelete) Name() string{ return "fs.delete" }
```

⇒ **六枚，全是 `fs.*`**，`internal/tools` 整包没有第七枚工具实现。

**(d) `fs.trash` 那一枚"碰外部世界"的：进程内 COM，不起进程**；非 Windows 侧是**拒绝**不是静默 no-op：

```
$ git grep -n 'procSHFileOperationW = ' 3f6322f -- internal/tools/recycle_windows.go
3f6322f:internal/tools/recycle_windows.go:58: procSHFileOperationW = windows.NewLazySystemDLL("shell32.dll").NewProc("SHFileOperationW")
$ git show 3f6322f:internal/tools/platform_other.go | sed -n '29,33p'
// shellTrash refuses. It returns an error before touching the disk.
func shellTrash(canonical string) (trashDetail, error) {
        return trashDetail{}, errors.New(
                "本平台无 Shell 回收站 API（Wisp 的 trash 只走回收站，绝不做删除），已拒绝执行：" + canonical)
```

**(e) 我这把尺比它的宽**：它 §3.1 用 `git grep -l '"os/exec"'`（只能抓 import 了 `os/exec` 的码）。
我改普查**所有 spawn 形状**（`exec.Command`／`StartProcess`／`CreateProcess`／`ShellExecute`／`WinExec`／
`plugin.Open`／`os.ForkExec`／`forkExec`），仍然**非测试文件全出**：

```
$ git grep -nE 'CreateProcess|ShellExecute|WinExec|StartProcess|exec\.Command|plugin\.Open|forkExec|os\.ForkExec' 3f6322f -- '*.go' ':!*_test.go'
cmd/balldebug/diff_windows.go:409            cmd := exec.Command(exe, childArgs...)
cmd/wisp/doctor.go:316                       out, err := exec.Command(cc, "--version").Output()
cmd/wisp/slo_windows.go:490                  cmd := exec.Command(exe, args...)
internal/llm/adaptertest/mockllm.go:65,72    exec.Command(goBin, "build", …) / exec.Command(exe, "-addr", …)
scripts/spike/webview2-latency/main.go:283,294
tools/d22scan/gitignore.go:373               cmd := exec.CommandContext(ctx, "git", args...)
```

⇒ 命中集**与它 §3.1 那张七枚表同一集合**（`internal/proc/jobscope_windows.go` 那枚在这把尺上**根本不出现在**，
因为它只 import `os/exec` 当参数类型、不 exec——这一点它自己写对了）。
⇒ **`internal/agent/**`、`internal/tools/**`、`internal/plugin/**`、`internal/panel/**`、`internal/llm/**`（生产半边）
在两种尺下都是零枚 spawn。**
`internal/llm/adaptertest` 只被测试 import：

```
$ git grep -l 'llm/adaptertest' 3f6322f -- '*.go' | wc -l        → 13
$ git grep -l 'llm/adaptertest' 3f6322f -- '*.go' | grep -c '_test.go'  → 13
```

**(f) 三枚 C4 槽不会"被问到就吐出一段可执行的东西"**（它 §3.4 只贴了 `providerSlots` 的字面，我补上射程内的消费者）：

```
$ git grep -n 'Slot(\|ErrSlotNotLanded\|SlotErr' 3f6322f -- '*.go' ':!*_test.go'
internal/tools/registry.go:29/33/35/40/42/44/47/48/52/161/169   ← 全部在 registry.go 自己那一棵
```

⇒ 生产码里除 `registry.go` 自身**没有任何消费者问到那三枚槽**，槽的答案是 `SlotErr`（`errors.Is` 得 `ErrSlotNotLanded`），
不是一个空 Provider。

**⇒ 支撑①：成立。** 而且比我进来前预期的更硬：不是"我没找到入口"，是**两把不同颗粒度的尺给出同一个空集**。

### 2.2　支撑② —— `internal/agent` 的依赖闭包里没有 `internal/proc`

```
$ go list -deps ./internal/agent | grep -c 'internal/proc'
0

$ go list -deps ./internal/agent | grep 'CarlosShao/wisp' | sort
github.com/CarlosShao/wisp/internal/agent
github.com/CarlosShao/wisp/internal/config
github.com/CarlosShao/wisp/internal/llm
github.com/CarlosShao/wisp/internal/memory
github.com/CarlosShao/wisp/internal/observe
github.com/CarlosShao/wisp/internal/plugin
github.com/CarlosShao/wisp/internal/risk
github.com/CarlosShao/wisp/internal/secret
github.com/CarlosShao/wisp/internal/winsec
```

**⇒ 支撑②成立**（我那发 `grep -c` 给 0，与它逐字同）。
⚠ 但第二发**给它这句加了一条它没写的限定**，我把它说清楚免得下游读歪：
**闭包里有 `internal/plugin`**（C11 的容器类型 `DisposalScope` 就在里面）。
⇒ 所以"任务级归属今天缺什么"的准确答案是：
**缺的不是 import 通路、不是容器类型，是"有人在新建这枚容器"以及"往这枚容器上挂一件 OS 级杀进程的动作"**。
它 §4.1② 那句"还缺一条从组合根把 `internal/proc` 递进 `agent.Options` 的通路"是对的；
但**读成"agent 连容器都拿不到"会偏**——容器拿得到，`JobScope` 拿不到。

### 2.3　支撑③ —— `NewDisposalScope` 生产码调用者 0 枚

```
$ git grep -n 'NewDisposalScope' 3f6322f -- '*.go'
internal/memory/retention_test.go:193   scope := plugin.NewDisposalScope("test-retention", nil,
internal/plugin/disposal.go:128         // not usable; use NewDisposalScope.            <- 注释
internal/plugin/disposal.go:152         // NewDisposalScope creates a scope derived…    <- 注释
internal/plugin/disposal.go:155         func NewDisposalScope(…)                        <- 声明
internal/plugin/disposal_test.go:17     s := NewDisposalScope(name, …)
internal/plugin/disposal_test.go:168    s := NewDisposalScope("task-tail", …)
internal/risk/provenance_test.go:307    scope := plugin.NewDisposalScope("session-1", …)
```

⇒ **七枚命中＝1 声明＋2 注释＋4 测试**，与它逐字同。**生产调用者 0 枚：成立。**

**我把它往下多问了一层**（它没问：除 `NewDisposalScope` 之外还有没有别的生产载体）：

```
$ git grep -ln 'wisp/internal/plugin' 3f6322f -- '*.go' | grep -v '_test.go'
internal/memory/retention.go
$ git grep -n 'StartRetentionJob' 3f6322f -- '*.go' | grep -v '_test.go'
internal/memory/retention.go:98:func (s *Store) StartRetentionJob(scope *plugin.DisposalScope, cfg RetentionConfig) {
internal/memory/retention.go:100:       panic("memory: StartRetentionJob requires a DisposalScope")
internal/memory/retention.go:234:// command); production scheduling goes through StartRetentionJob.
$ git grep -n 'CloseTask' 3f6322f -- '*.go'
internal/tools/bridge.go:642:// CloseTask closes one task's taint scope. …
internal/tools/bridge.go:646:func (b *Bridge) CloseTask(taskID string) {
```

⇒ **唯一 import `internal/plugin` 的生产文件是 `internal/memory/retention.go`**，
而它那道 `StartRetentionJob(scope *plugin.DisposalScope)` 在**生产码里零枚调用者**（只有 `retention_test.go:198`）。
⇒ **`Bridge.CloseTask`（把 per-task taint scope 关掉的唯一门）在全仓零枚调用者——连测试都没有**；
`OpenTask` 却在派发热路径上被调（`bridge.go:559`）。
⇒ **支撑③成立**："任务级归属"这仓里确实**只有契约文本、没有生产载体**。

### 2.4　本程造出的一发（它没做、也不该由它做的）：**同一枚"无载体"另有一处今天是实伤**

§2.3 那枚 `CloseTask` 零调用者，不是"给未来预留"的形状——`OpenTask` **今天就在生产路径上被调**：

```
$ git grep -n 'OpenScope\|CloseScope' 3f6322f -- '*.go' ':!*_test.go'
cmd/wisp/panel_assets.go:232   prov.OpenScope(taintSourceScopeID)
internal/tools/bridge.go:638           b.prov.OpenScope(taskID)     <- 由 OpenTask 调，OpenTask 由 bridge.go:559 调
internal/tools/bridge.go:655           b.prov.CloseScope(taskID)    <- 只有 CloseTask 调它，而 CloseTask 无人调
internal/risk/provenance.go:387  logf("risk/C25: scope %q not open (CloseScope raced or composition gap); taint stored — it cannot leak into other scopes", scopeID)
```

⇒ 读数：**per-task 污点 scope 今天只开不关**（关它的那道门唯一的接线说明是
`bridge.go:642-643` 那句"the composition root defers this on the task's DisposalScope"——而那枚 DisposalScope 生产零枚）。
⇒ **这不是本票的对象**（本票管"进程还活着吗"，这管"污点记录还在不在内存里"），
但它是**同一枚缺失载体的第二个消费者**，且**今天可达**。
⇒ 处置：按本仓既有裁定（**勾交付物＋把洞归口另一张票＋原话不改**），我**不**把它塞进 138 的裁决里，
只登记并指名归属：`internal/risk/provenance.go` 的 C25 族（票 19/25/26 那一线，`provenance.go:56` 自己点名"tickets 12/19/21/22/26"）。
⚠ 本程**没有量它的后果边界**（见 §2.6）。

### 2.5　本格裁决

**三发支撑全部成立**（各自独立复算，命令原文在上）⇒
**"今天没有任何一条生产路径能在一轮任务里 exec 出子进程"这句话成立**，
于是票面 AC#2 的条件句前件**确实为假**、"给未来功能预留的洞"这枚裁决**不塌**。
派单给的那道"只要一发不成立就翻案"**没有被触发**。

**档位：成立（附两处点名，不构成退回）**：
1. 被验件 §3.7③ 那句 **"`internal/panel`、`internal/speech` 至今只有 `doc.go`" 是事实错**——
   `internal/panel` 有 **9 枚生产 `.go`**（`approval.go`／`assets.go`／`attachments.go`／`bridge.go`／
   `composer.go`／`composer_handlers.go`／`doc.go`／`pump.go`／`workspace.go`），只有 `internal/speech` 是 `doc.go` 一枚。

   ```
   $ git ls-tree -r --name-only 3f6322f -- internal/panel | cat        （20 枚，其中 9 枚非 _test.go）
   $ git ls-tree -r --name-only 3f6322f -- internal/speech | cat       （1 枚：doc.go）
   ```

   ⚠ **它那句错话在哪、要紧在哪**：它出现在 §3.7「本格没测什么」第③条的理由里
   （"因为那条组合根本程没找到"）。**结论没受影响，但理由是错的**：
   面板那棵**今天确实不是第二枚工具组合根**，正确的凭据是**它不 import `internal/agent` 也不 import `internal/tools`**：

   ```
   $ git grep -n 'wisp/internal/agent\|wisp/internal/tools' 3f6322f -- internal/panel
   （空输出，rc=1）
   ```

   ⇒ 下游若照它那句"panel 只有 doc.go"去排"第二枚组合根"，会在**panel 已经有 20 枚文件的今天**排错面。
2. §2.2 那枚限定：`internal/plugin` **在** agent 的闭包里，所以"缺通路"要说准是缺 **`internal/proc` 的通路**、不是缺容器。

### 2.6　本程没测什么（按"漏了它谁会先被骗"排序）

1. **没有运行时观测**：没跑一次真任务数子进程数。`cmd/wisp` 今天仍是在飞写者（149 已交、139 在写），
   构它取到的是中间态。⇒ §2.1 那个"0"与它一样，是**静态调用点**读数。谁把它升格成"运行时也没有"，谁就在替本程说谎。
2. **没排第三方自发子进程**：`modernc.org/libc`、sherpa 那族 cgo/native 绑定里有没有 `CreateProcess`。
   它 §3.7② 同处自陈没排；**我也没排**，所以"任务不会起子进程"严格说只覆盖**本仓自己写的码**。
3. **没量 §2.4 那枚实伤的后果**：provenance map 会不会无上界长大、`TreeProcessCount` 类仪器看不看得见它，
   本程一枚都没测。⇒ 若有人据此开票，得先有尺，别抄我的"零调用者"当危害。
4. **没验槽位在真被 `ProviderFor` 问到时的行为**：我只证明了生产码无人问（§2.1f）。
   真要落 D46，那枚 SlotErr 会不会被某个上层当成"空插件"咽掉，是落地当天的判据，不是今天的。
5. **没在别的分支取数**（`git branch -a` 都没扫）。它与它 §4.5① 同一枚限制。

### 2.7　结论修正记录（本格）

1. **推翻被验件一处（事实层）**：`internal/panel` 绝非"只有 `doc.go`"（§2.5 点名 1 的读数为凭）。
   受影响的只是它 §3.7③ 的**理由**，不影响它 §3.5 的**结论**——我用另一条凭据（panel 不 import agent/tools）
   把它那句"没找到第二枚组合根"的**结论救回来了**。**这一处该记在它表上，但记的是"理由错、结论对"，不是"结论错"。**
2. **加强被验件一处**：它 §3.1 那把尺（`os/exec` importer 普查）覆盖面窄于本案需要；
   我换成 spawn 形状普查后**命中集不变** ⇒ 支撑①比我进来前预期的更硬。这一条是**给它加分**，不是挑刺。
3. **收紧被验件一处措辞**：`internal/plugin` 在 `internal/agent` 闭包里（§2.2），
   所以 §4.1②"连整机 Job 都拿不到"要说准：**拿不到的是 `JobScope`，不是 scope 容器**。
4. **本程自纠一处**：我第一次跑 spawn 普查时把 `':!*_test.go'` 与 `-- '*.go'` 的顺序写反了一次，
   输出仍含测试文件；发现后重跑（§2.1e 贴的是重跑那发）。**没有拿那一发下过任何结论。**

---

## 3　派单第 4 节的两判 —— 那笔 `DEFERRED` 登记能不能落地 ／ 本票该"作废并登记"还是"改判保留"

**判据**（派单 §4）：ⓐ 登记是否只能摆成 owner 的 `Q##`，说清为什么不能由 agent 单方落；
ⓑ 作废并登记 vs 改判保留，按本仓规矩各适用哪一条，**别混**。

### 3.1　第一发：`docs/specs/**` 到底被什么冻着（把两枚不同来源的冻结分开）

```
$ T=.scratch/wisp/issues/138-cancelling-a-task-kills-nothing-the-task-itself-spawned-single-jobscope-no-terminatejobobject.md
$ git show 3f6322f:"$T" | sed -n '9,12p'
**Packages:** `internal/proc/**`（`jobscope_windows.go`）· `internal/agent/loop.go`（取消路径）·
              `cmd/wisp/**` 若有接线点。**⚠ 冻结面照旧禁改**：`internal/risk/**` · `rules_gateway.go` ·
              `tools/d22scan/**` · `allowlist.txt` · D32 阈值（CPU≤0.5%／RSS≤25MB）· `thresholds.go` ·
              任何 golden · `frontend/**`（owner 已交外部 agent）· `docs/PLAN.md` · `docs/specs/*.md`。
$ git show 3f6322f:docs/specs/SPEC-12-roadmap-governance.md | sed -n '36,39p'
### 4.1 契约变更流程
改 C1–C32 或 D1–D47 = **人工批准**；同步更新 PLAN.md、`docs/DECISIONS.md`、受影响切片卡。
agent 单方面改契约 = 跑歪模式 #1，对抗验收判失败。
```

**读数（两枚冻结不是一枚）**：
1. **路径冻结：成立、且就写在本票票面**（`docs/specs/*.md` 在那行"冻结面照旧禁改"里）。
   ⇒ 对**实现程**与对**本验收程**这都是硬边界，不需要引任何别处的条文。
2. **"§4.1 契约变更"这一枚：不覆盖这笔登记。** `SPEC-12 §4.1:38` 的射程是**由 `C1–C32`／`D1–D47` 那两份枚举报决定的**
   （这个判法不是我发明的，是本仓 `A181②` 已经这么裁过一次并把方法写进了台账）。
   `SPEC-12 §5` 的 DEFERRED/RESERVED 表行**不在那两份名单里** ⇒
   **"往 §5 加一行"按 §4.1 的字面不构成契约变更**。
3. 反过来，同一份 SPEC-12 又把登记列为**每片完成的强制动作**：

```
$ git show 3f6322f:docs/specs/SPEC-12-roadmap-governance.md | grep -n '登记表更新'
57:5. DEFERRED/RESERVED/REJECTED 登记表更新（有新增推迟必须五字段齐全）
```

⇒ 那一节是 `### 4.3 每切片完成后的强制动作`（`:51`），第 5 项在 `:57`。

⇒ **所以这件事不是"没人有权做"，是"做的人不能是想要它的那一枚"**（实现者给自己开推迟条＝裁决者≠实现者，`AGENTS §0.3`）。
⇒ ⚠ 本票内部还有一处**自相矛盾**要登记：**票面 `:12` 冻结 `docs/specs/*.md`，票面 `:44-46` 又命令"必须去 `SPEC-12 §5` 登记"**。
两枚出自同一张票的要求**今天不可能同时满足**——与它 §3.4 抓到并登记的那枚 `internal/tools/unwired.go` 路径打偏同类，
是**票面缺陷、不是实现缺陷**；解的人只能是编排者（改 AC#1 那句的落笔者，或给票面加一枚 `>` 具名豁免）。

### 3.2　第二发：`AGENTS §1.1` 那条"1:1 双向"到底挡不挡这笔登记（我把两个方向都量了）

权威文本（不在 AGENTS，AGENTS 自称薄索引）：

```
$ git show 3f6322f:.scratch/wisp/issues/README.md | sed -n '188p'
- `DEFERRED(D-xx)` code markers must map 1:1 to SPEC-12 §5 registry (bidirectional check).
```

**方向 A（标记 → 表行）：今天已经违着两条。**

```
$ git grep -n 'DEFERRED(D28-1)\|DEFERRED(D11-3)' 3f6322f -- 'internal/**'
3f6322f:internal/agent/compress.go:26:  // DEFERRED(D28-1): move the call into the Warm window hook; keep this
3f6322f:internal/agent/loop.go:394:               // DEFERRED(D28-1): the Warm-window hook owns this call; the
3f6322f:internal/agent/control.go:15:     // this ticket (DEFERRED(D11-3), see the ticket report).
$ git show 3f6322f:docs/specs/SPEC-12-roadmap-governance.md | grep -n 'D28\|D11'
91:| REJECTED | 意图分类四级流水线 | D11 过度设计已废 | — | D11 两层替代 | 无前置分类器（特性非缺失） |
```

⇒ `D28-1`／`D11-3` 两枚**代码标记在 `SPEC-12 §5` 里没有对应行**（`:91` 那枚是 D11 的 REJECTED 行，不管 `D11-3`）。
⇒ 这一枚**在飞的票 139 今天也独立量到了**（`4d0866a` 的 commit 标题："D28-1 在 SPEC-12 §5 里根本没有那一格"）。
⇒ 历史上它被判过一次：`docs/evidence/s1/10-adversarial-acceptance.md:106`——
"`DEFERRED(D28-1)` … 符合契约明文豁免 … **不算违规**"。⇒ 也就是说**这条 1:1 今天没有任何仪器在守**
（我扫过：`git grep -rn 'DEFERRED' 3f6322f -- scripts tools .github` 只命中 `scripts/portable-tests.sh:50/:95` 的**注释**，
没有一把尺）。

**方向 B（表行 → 标记）：这条要求根本不成立。**

```
$ git show 3f6322f:docs/specs/SPEC-12-roadmap-governance.md | sed -n '66,92p' | grep -c '^|'
27
（同一段的"项"列逐名：macOS 平台层 / 代码签名包管理器分发 / 插件 SDK+文档+示例 / 社区插件 registry /
  i18n非中文 / 无障碍 / AEC barge-in / 快捷键路径语音否决（B1）/ 剪贴板历史（D34）/ `doc.read` xlsx/OCR（D34）/
  `system.eject`（D34）/ 完整错误文案体系 / 竞品对比文档 / MCP client / Codex guardian 双轴评分 /
  跨会话持久授权档 / `available_decisions` / embedding 语义记忆 / Linux 支持 / 精确 taint tracking /
  周期性调度器（cron）/ 邮件日历IM 核心内置 / Coding IDE 能力 / 多 Agent 协作 / 意图分类四级流水线 / `fs.delete` 默认提供 …）
```

⇒ **27 枚数据行全是散文键**，一枚都不带 `DEFERRED(D-xx)` 代码标记（个别在"项"列里带 `(D34)` 这种**出处标注**，
那不是代码标记）。⇒ 权威文本那句 `DEFERRED(D-xx)` **限定的就是标记那一侧的形状**，
**不是**"每条登记都得配一枚代码标记"。

**⇒ 这一发的裁决：被验件 §4.2 那格给的理由有一枚是读严了的。**
它写"本程单侧改表会造出一枚**没有代码标记对应**、也没有验收表的登记"——
前半句按盘上的表与 `issues/README:188` 的字面**不成立**；
**成立的后半句只有一条**：`docs/specs/*.md` 被票面 `:12` 冻着，且**它自己是那笔登记的受益人**（裁决者≠实现者）。
⇒ 这不影响它的**行为**（不落笔是对的），影响的是**下一程引这条理由时的射程**——
别让人以为"配不上代码标记就不能登记"，那会让 27 枚现存行全都变成假想违规。

### 3.3　第三发：类型取 `DEFERRED` 还是 `RESERVED`（我复核它取 DEFERRED 的前提）

它取 DEFERRED 的前提是"shell.exec **有计划**（SPEC-07 §3，S3）／D46 是已采纳项、有票 50"。逐条核：

```
$ git grep -n 'shell.exec' 3f6322f -- docs/specs/SPEC-07-tools-and-plugins.md
3f6322f:docs/specs/SPEC-07-tools-and-plugins.md:71:| `shell.exec` | 任意命令 | L2（白名单命中 → L1） | shell | S3 | **默认禁用**；argv 向量强制 |
$ git grep -n '^## ' 3f6322f -- docs/specs/SPEC-07-tools-and-plugins.md | sed -n '3p'
32:## 3. D34 内置工具权威表（唯一来源；落 `docs/TOOLS.md`）        <- 它引的"SPEC-07 §3"就是这一节，对得上
$ git show 3f6322f:.scratch/wisp/issues/50-tier1-manifest-plugins.md | sed -n '3p'
**Status:** ready-for-agent
$ git show 3f6322f:docs/evidence/s1/138-cancel-kills-children-r1.md | sed -n '554p' | awk -F'|' '{print "登记行的列数="NF-2}'
登记行的列数=6
$ git show 3f6322f:docs/specs/SPEC-12-roadmap-governance.md | sed -n '64p'
| 类型 | 项 | 为什么现在不做（依据） | 完成判据 | 前置 | 当前残缺表现 |
```

⇒ **DEFERRED 取对了**：`RESERVED` 的语义那句是"**无实现计划**"，而 `shell.exec` 在 SPEC-07 §3 权威表里有行、
D46 命令插件有 `ready-for-agent` 的票 50 ⇒ "无计划"不成立。
⇒ 它那行**六列与表头逐列对齐**（`fields=6`），"五字段缺一不可"在两种数法下都满足——这一条我不改它。
⇒ 唯一要补的：它的"依据"列里写的是一串**〔本程现量〕**，而那批现量今天才被我这一格复算完。
⇒ **落笔时"依据"列应改指本验收件**（否则登记行的凭据是一枚未验收的自述——正是 `A265③` 那族"把自述当读数"的形状）。

### 3.4　第四发（本程造出的一发，派单没要求）：`ready-for-human` 这枚状态值**在本仓不存在**

票面 `:44-46` 命令"把本票降级为 `ready-for-human`"。我去核池子的状态词表：

```
$ git show 3f6322f:.scratch/wisp/issues/README.md | sed -n '9,15p'
| Status | Meaning |
|---|---|
| `ready-for-agent` | Frontier ticket, unclaimed |
| `in-progress` | Claimed; must have ≥1 Progress-log entry |
| `blocked` | Waiting on external decision (ticket notes which) |
| `review` | Work done, awaiting adversarial acceptance (SPEC-10 §8) |
| `done` | All boxes checked; **rename file with `-done` suffix** + title `(DONE ✅)` |
$ git grep -n 'Rules (prevent' 3f6322f -- .scratch/wisp/issues/README.md
3f6322f:.scratch/wisp/issues/README.md:17:**Rules (prevent abandoned/in-progress tech debt):**
$ git show 3f6322f:.scratch/wisp/issues/README.md | sed -n '30,33p'
4. Completion: check all acceptance boxes → `Status: done` → rename file `NN-slug.md` →
   `NN-slug-done.md` → update this index → commit+push.
5. If a decision is missing, set `Status: blocked` with the open question in the log
   (D22 闸门③: undefined = stop and ask, never assume).
$ git grep -rln 'ready-for-human' 3f6322f
3f6322f:.scratch/wisp/issues/138-cancelling-a-task-kills-nothing-the-task-itself-spawned-single-jobscope-no-terminatejobobject.md
3f6322f:docs/evidence/s1/138-cancel-kills-children-r1.md
```

⇒ 池子的词表是**五枚**：`ready-for-agent` / `in-progress` / `blocked` / `review` / `done`。
⇒ **`ready-for-human` 全仓只出现在本票与它的件里——是一枚池外状态值。**
⇒ **`ready-for-human` 全仓只出现在本票与它的件里——是一枚池外状态值。**
⇒ 于是"降级为 ready-for-human"这句**要按池内语言翻译才能执行**，而翻译结果唯一：
**本票的处境正是规则 5 那句 "a decision is missing"**（缺的决定＝"owner 已批准的一枚票要不要改判为今天不做"）
⇒ **`Status: blocked` ＋ 待决问题进日志 ＝ 一枚 `Q##`**。
⇒ 顺带一句实现程的判断没走歪：它 §4.2 那格写"Status 属同一族的票面框，未改"——**它一枚字都没动是对的**，
虽然它引的那句"票面框由编排者按非实现者验收表来定"是派单侧、盘上无 witness（见本件 §1.4）。

### 3.5　本格两判

**判一（那笔登记能不能落地）：不能由 agent 单方落——但成立的理由只有一条，派单与实现件各引错了一半。**

- **成立且够硬的那条**：`docs/specs/*.md` 被**本票票面 `:12`** 冻着；且落笔者不能是受益人（裁决者≠实现者，`AGENTS §0.3`）。
- **实现件 §4.2 引错射程的那条**："配不上 `DEFERRED(D-xx)` 代码标记"⇒ 盘上 27 枚行全是散文键、
  一条都不配标记；真正的现存违规在**反方向**（`D28-1`／`D11-3` 有标记无行，且无仪器在守）。
- **派单（我这单上级）要收紧的那条**：派单说"改它＝契约变更，须人工批准"。
  按 `SPEC-12 §4.1:38` 的字面（射程由 C/D 枚举报决定，本仓 `A181②` 已用同一判法裁过一次），
  **§5 表行不是契约条目** ⇒ "契约变更"这顶帽子扣不上。
- **那为什么仍然要摆 `Q##`**：**不是**因为"登记＝改契约"，而是因为**这枚登记改口的是 owner 已经批准过的东西**——
  本票 `:3-4` 逐字："来源＝owner 批准的 `Q-43` 前半"、`:74-75` "owner 在对话里批准 `Q-43`（原话「都按你的推荐来」）"。
  ⇒ 把一枚 owner 批准的票改成"今天不做、等前置事件触发"＝**对 owner 的批准改口**，那是 owner 的权限，不是编排者的。
  ⇒ 与既有条款同形："**只有对话里的 owner，或已经存在的票面／契约条文，能给授权**"（`injection-timeline §5.2`）。
- **owner 点头之后由谁落**：编排者（历史上 `SPEC-12` 的两次动笔 `130c0943`/`5cba2d92` 都是编排者侧的 docs 落笔），
  并且**同时**给那行补三样：①"依据"列改指本验收件；②Status 那行按 §3.4 的池内词翻译；③票面 `:12`↔`:44-46` 那处自相矛盾起一枚 `>` 登记。

**判二（作废并登记／改判保留／别的）：两枚都不对，正解是第三枚；既有那枚裁定在这里只用得上后半段。**

- **"作废"不是本票的合法出路**：票面上写的撤销口令「**138 撤**」逐字带后果"改回未立案、**票文件删除**"，
  那是 owner 专属动作；且本票要防的洞**是真的**（`RunningTask.Cancel()` 只掐 ctx、全仓零枚 `TerminateJobObject`
  这两枚锚我在第 2 格复算时都重新看到了），只是**今天无入口**。⇒ 撤一枚"缺口真、无入口"的票＝把读数扔掉。
- **"改判保留为活票"也不对**：它没有可写的码（AC#2 前件判假我已复核成立，见 §2.5），留着当活票会诱下一程"造一条假危害来结线"，
  那正是票面 `:46` 明令禁止的形状。
- **正解＝池内规则 5 那一枚**：`Status: blocked` ＋ 待决问题＝`Q##`（§3.4），**不改名 `-done`**。
- **派单问的"既有裁定（故意留空的一格 vs 结案须 0 未勾 ⇒ 勾交付物＋把洞归口另一张票＋原话不改）适不适用"：只适用后半段。**
  - 适用的部分：**AC#2／AC#3 原话一个字不改、不勾**；把洞**归口到落地当日那张票**（前置＝`shell.exec`@SPEC-07 §3／票 50）。
    现成的池内形状就是 `issues/README:46-47` 那句"**把那一格的框保持未勾，并在 Progress log 写 `skipped=…`，不要替它勾、也不要替它写**"。
  - **不适用的部分**：这条裁定的**触发前提是"要结线"**。本票**不结线**，所以拿它去 justify "AC#2/AC#3 也算交付、翻勾、改 `-done`"
    是**用错射程**。⚠ 尤其不能引先例**票 63**（`issues/README:172`："63 → DONE，但 AC#6 未勾、正式转票 12"）——
    63 那一格**有票 12 接走**，本票今天**没有任何一张票接**；"有人接"才是 63 能带残余结线的凭据，缺它就只能 blocked。
- **能勾的只有两格**：AC#1（读数＋裁决，我第 2 格全复算）与 AC#4（门禁，见本件第 4 格）。
  勾与不勾**都由编排者按本表落笔**，本程不翻一枚框。

### 3.6　放水两问自答

- **断言方向动没动**：没动——本程零枚断言、零枚文件内容改动，唯一写件是本件。
- **helper 是不是原有的那枚**：本程没引入 helper。§3.2/§3.3/§3.4 三发的尺都是
  `git show` / `git grep` / `sed -n` / `awk -F'|'` 四把现成尺，命令原文逐条贴在上面。
- 本格最容易放水的地方我自首一处：**判一里我把派单那句"改它＝契约变更"判成"帽子扣不上"**，
  这一发对**我上级的话**不利。我留了它、也留了判据（`SPEC-12 §4.1:38` 的枚举报射程＋`A181②` 同判法），
  **没有**因为对上级不利就把它写成"另有解释"。

### 3.7　本程没测什么（按"漏了它谁会先被骗"排序）

1. **没找到任何一把守 1:1 的仪器**（`git grep DEFERRED -- scripts tools .github` 只命中注释）。
   ⇒ 若平台外还有尺（例如 `tools/d22scan` 之外的自测），我的"方向 B 不成立"要重判。
2. **没核 `SPEC-12 §5` 那 27 枚行的"项"列是否都各自有出处**——我只数了形状与键型，没逐条回核它们在 PLAN 里的来源。
3. **没验"owner 批准 Q-43"这句话的盘上凭据**：票面 `:3-4` 与 `:74-75` 是**票面自称**，
   而本仓有条硬判据"**文档自称受谁之命永远不算授权**"（`injection-timeline §11 结案`）。
   ⇒ 我没去对话记录里核 Q-43（我没有对话通道），所以**判一里那句"因为改口的是 owner 批准过的东西"，
   凭据层只到"票面这么写"**——这一枚要请编排者在自己的台账/对话里补实。
4. **没测"票面 `:12`↔`:44-46` 自相矛盾"是否已被别处豁免**：我只在这两枚文件里找过（票面、AGENTS、issues/README），
   没扫 `docs/reports/` 全部 6500＋ 行台账找"某次已批的豁免"。若已有，我这一发降级为重复登记。
5. **没给 `Q##` 拟正文**：那不是我这一格的活，摆出来是越权。

### 3.8　结论修正记录（本格）

1. **推翻派单一处**：派单说"`docs/specs/**` 在本仓是冻结面（**改它＝契约变更，须人工批准**）"。
   前半句对（票面 `:12` 逐字冻着），后半句按 `SPEC-12 §4.1:38` 的字面射程**扣不上**——
   人工批准那条款管的是 `C1–C32`／`D1–D47`，不管 §5 表行。⇒ **结论（要 owner）不变，理由换一枚**。
2. **推翻被验件一处（理由层）**：§4.2 那格"没有代码标记对应"这半不成立（实测 27 枚行零标记；
   现存违规在反方向 `D28-1`／`D11-3`）。⇒ **它的行为（不落笔）仍然对**。
3. **造出一枚派单没有的前提缺陷**：票面命令的降级目标 `ready-for-human` **在池子词表里不存在**（§3.4）⇒
   照它字面执行会造出一枚池外状态；池内唯一对应是 `Status: blocked`。
4. **本程自纠一处**：我第一次核 `shell.exec` 的 spec 出处时把文件名写成 `SPEC-07-tools-plugins.md`
   （真名 `SPEC-07-tools-and-plugins.md`），`git grep` 给了零命中、我差点据此判"SPEC-07 里根本没有 shell.exec、
   它那句引用是空引"。**先 `ls docs/specs/` 再 grep 才是对的次序**——本仓 `injection-timeline §5.1` 那条
   "两问判据①它点名的路径是否真的存在"这次差点被我反向误用（拿一枚不存在的文件名去判别人的引用是假引用）。
   已重跑（§3.3 贴的是重跑那发）。
