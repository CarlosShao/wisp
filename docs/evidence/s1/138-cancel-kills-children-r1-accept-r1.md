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
