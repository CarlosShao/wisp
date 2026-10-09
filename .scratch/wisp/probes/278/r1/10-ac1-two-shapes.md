# 278-r1 AC#1 逐条裁形状（两枚不压成一格）

取数面＝**HEAD 快照**（锚 `bf9974c2`，本腿起手时 HEAD；起手件 `00-anchor.md`）。
尺射程统一声明见同目录 `20-ac2-ac3.md` 末「尺清单」。禁 本腿未跑任何 `go build`/`go test`/`go vet`/`go list`/`go doc`。

---

## 0. 票面行号现量对拉（先报漂，再看结论）

| 票面/台账原文 | 本腿 HEAD 现量 | 漂 |
|---|---|---|
| 票 278「`resident_approval_windows.go:697` `AskOnTaskRoot`」 | `HEAD:cmd/wisp/resident_approval_windows.go:698`（`func (ra *residentApproval) AskOnTaskRoot(...)`），注释块起 `:689` | **+1** |
| 票 278「函数体内 `:700` 唯一产码调用 `askConfirmation`」 | `HEAD:cmd/wisp/resident_approval_windows.go:701`（`return ra.askConfirmation(taskCtx, c)`） | **+1** |
| 票 278「`cmd/wisp/resident_windows.go:248` 那句 `This call is the caller`」 | 那句短语落在 **`:249`**；`:248` 装的是 `// the way out - and AskOnTaskRoot / askConfirmation still had zero product` | **+1**（票面沿用了台账 `A720` 那个已被 `A7xx` 自己更正过的旧行号；注释块整体是 `:246-259`） |
| 台账 `pending-and-issues.md:11139`「尺只命中定义 `gate.go:749`」 | `HEAD:internal/agent/approval/gate.go:755`（注释 `:751-754`） | **+6** |
| 产码注释 `resident_task_source_windows.go:13`「AskOnTaskRoot (resident_approval_windows.go:**261**)」 | 同一枚符号今天在 `:698` | 产码注释里那枚行号**过期 237 行**（这是注释，本腿禁 不改，只具名） |

`git show HEAD:<file> | wc -l` 存在性尺：`resident_approval_windows.go` **989**、`resident_windows.go` **288**、
`internal/agent/approval/gate.go` **793**、`queue.go` **621**、`ui.go` **189**、`replies.go` **527**、
`resident_task_source_windows.go` **557**、`resident_approval_live_246_windows_test.go` **543**、`resident_approval_246_windows_test.go` **378**。

---

## 1. `AskOnTaskRoot`（`cmd/wisp/resident_approval_windows.go:698`）

### ① 调用形状尺（现跑，负向句锚在调用形状上）

尺 1（点号调用形状，含测试）：

```
git grep -nE "\.AskOnTaskRoot\(" HEAD -- ':(exclude).scratch' '*.go'
```

读数 **rc=0，命中 3 枚**：`cmd/wisp/resident_approval_live_246_windows_test.go:115`／`:220`／`:311`（三枚全是 `ra.AskOnTaskRoot(residentCard{...}`）。
剥掉测试后＝**0 枚产码调用者**。

尺 2（把「剥」这一步做成尺子本身，测试文件从路径面就排除）：

```
git grep -nE "[A-Za-z0-9_]\.AskOnTaskRoot\(" HEAD -- ':(exclude).scratch' '*.go' ':!*_test.go'
```

读数 **rc=1（空）** ⇒ 「零产码调用者」这枚负向句**成立**，且它是**整族**结论不是抽样。

尺 3（名字族闭合，用来堵接口派发／方法值这两形）：

```
git grep -nE "AskOnTaskRoot" HEAD -- ':(exclude).scratch'
```

读数：`.go` 面共 **8 枚**命中＝定义 `:698` ＋ 注释 `resident_approval_windows.go:689`、`resident_task_source_246_windows_test.go:7`、
`resident_task_source_windows.go:13`、`resident_windows.go:248` ＋ 尺 1 的 3 枚测试调用；其余命中全在 `docs/**` 与 `.scratch`（后者已被排除）。
**没有任何接口声明、结构体字段、方法值（`ra.AskOnTaskRoot` 不带括号）出现在这个名字上** ⇒ 本枚禁 不需要写「判不动」：
接口派发形在尺 3 下必留同名标识符，尺 3 没留。同族尺还顺带闭合了 `askConfirmation`：

```
git grep -nE "\.askConfirmation\(" HEAD -- ':(exclude).scratch' '*.go'
```

读数 **3 枚**＝测试 `resident_approval_246_windows_test.go:97`／`:161` ＋ **产码 1 枚 `resident_approval_windows.go:701`**，而 `:701` 住在 `AskOnTaskRoot` 函数体内（`:698-702`）
⇒ `askConfirmation` 是**传递不可达**（唯一上游自己零调用者）。票面这句**成立**。

### ② 编译进产物还是根本没进

- 文件：`cmd/wisp/resident_approval_windows.go` 第 1 行＝`//go:build windows`，文件名又带 `_windows.go` 后缀（双锁）。
  ⇒ **只在 `GOOS=windows` 参与编译**，且它在 `package main`（`cmd/wisp`）＝双击那个常驻二进制里 ⇒ **产码定义这一侧编进产物**（非编外）。
- 三枚调用者所在的 `cmd/wisp/resident_approval_live_246_windows_test.go` 第 1 行＝`//go:build windows && winlive`
  ⇒ 默认 `go test ./cmd/wisp` **不编它**，只有带 `-tags winlive` 才编。
- ⚠ 本腿**顶回一条台账/索引转述**：`AGENTS.md` 抄的台账 `A595` 段与 `pending-and-issues.md` 里那句「`winlive` 在 CI 无 tag 档／CI 的 `cmd/wisp` 步不跑该 tag（尺＝`grep tags|winlive scripts/wisp-cli-tests.sh .github/workflows/ci.yml` 零命中）」
  **今天在 HEAD 上已过期**：本腿尺
  `git grep -nE "winlive" HEAD -- '.github' 'scripts'` ⇒ **rc=0**，命中 `.github/workflows/ci.yml:606`
  `- name: "winlive compile gate (go vet -tags winlive, ticket 111 AC#11)"` 与 `:655` `run: go vet -tags winlive ./cmd/wisp/ ./internal/ball/`；
  落档时刻＝`git log --oneline -1 -S"winlive compile gate" -- .github/workflows/ci.yml` ⇒ **`351e5a5e` 2026-10-08（票 111 r6 AC#11）**。
  ⇒ 精确形状：**CI 里 winlive 只被 `go vet` 编译门编到、不被执行**（`ci.yml` 自己的注释还写着「do not read this step as "the winlive cases now run"」），
  所以那 3 枚用例是「CI 保证编得过、不保证跑过」的一档，比旧转述的「完全在 CI 之外」强一档、比「CI 跑过」弱一档。
- 链接器面（是否真留在最终 PE 里）：要 `go build`／`nm` 类读数才知道，本腿禁跑 ⇒ **判不动**（缺的读数＝`go build` 后的符号表；本腿不补）。

### ③ 谁在什么条件下会需要它（只给形状，不给方案）

读函数体 `:698-702` 与它上游 `askConfirmation` 体 `:642-687`：

- `AskOnTaskRoot` 只做两件事：`taskCtx, cancelTask := context.WithCancel(ra.root)`（把卡片绑到**本进程任务根**，即 D38(e) 第 3 步取消的那枚 ctx），再转 `askConfirmation`。
  注释 `:695-697` 具名写着「故意不接受调用方传的 ctx」＝一枚能脱离 root 的卡就是退出序列够不着的卡（AC#4 要关的形状）。
- `askConfirmation` 的四道前置形状：`:643-645` `c.TaskID`/`c.Reason` 空即 fail-closed 拒；`:647-650` `ra.closed` 即拒；`:667-668` `ra.gate.AdmitTextTask(c.TaskID)`（**D47 授权**，用完即撤）；`:677-682` 按 `c.Level` 分流 **L1→`gate.PendingWindow`／L2→`gate.PendingApproval`**，L0 走 `:682-687` 拒绝。
- **今天举卡走的不是它**：`resident_windows.go:260` `src := startResidentTaskSource(rt, ra)` → `resident_task_source_windows.go:269` `gate: ra.gate` → `:278` `run, code := assembleRuntime(spec)` ⇒ 卡片由**模型那一次真工具调用**经 bridge→gate→`ballCardUI.Prompt` 举起。`AskOnTaskRoot` 是**宿主自己发起**那一形（不经工具调用）的唯一现成接缝，注释 `:690-691` 具名点了两枚未来来路：**语音链**、**票 228 的 `config.toml` 接线**。
- 接上它要动的形状（枚数＝本腿读到的，非方案）：
  1. 调用者只能落在 **`package main` + GOOS=windows**——接收者 `*residentApproval` 与参数 `residentCard`（`:623`）都是非导出的，跨包禁 够不着；
  2. `residentCard` 今天**在产码里零构造点**（尺：`git grep -nE "residentCard" HEAD -- ':(exclude).scratch' '*.go' ':!*_test.go'` ⇒ **5 枚全在 `resident_approval_windows.go`**：注释 `:79`、`const residentCardTool` `:82`、类型声明 `:623`、两枚函数签名 `:642`/`:698`）⇒ 要接就得先有一枚产码构造点；
  3. `ra.root` 必须是活的、`ra.escLoad` 必须为真（否则 `:796-805 residentStatusLine()` 直接报「审批门未装配」，卡片没有可承载它的球窗口）；
  4. 宿主侧得自报一枚**已按 D47 登记**的任务标识，否则 `:643` 或 `admitCheck` 那一头把它拦在门外。
- 它**已经在替自己说话**：注释 `:691-693` 明写「today its only callers are this package's own cases, and this file says so in the boot report rather than letting the existence of the seam read as a running pipeline」，
  且 boot 报告 `resident_windows.go:269-270` 打的 `任务来源：%s` 来自 `src.taskPosture()`（四档常数 `resident_task_source_windows.go:120-123`，无控制台时 `taskPostureAbsent`＝「无（任务入口未启用…）」）⇒ 「没接」这件事今天**有文字面在钉**，只是钉在 posture 而不是钉在 `AskOnTaskRoot` 上。

### ④ 如果不接：删 vs 留注释，各自代价一句话

- **删**：三枚 `windows && winlive` 真机用例（`:115`/`:220`/`:311`）当场编不过 ⇒ CI 那道 `go vet -tags winlive`（`ci.yml:655`）立刻红，且票 246 AC#4 的「宿主发起的卡」接缝作废、语音链落地时要重写；`askConfirmation` 会跟着变成零调用者的第二层。
- **留注释**（现状即此）：代价＝`cmd/wisp` 里长期挂一枚零产码调用者的接缝，且它会被后面的程读成「已经接上了」——**票 278 的立票理由正是这一枚误读**，而今天替它说话的是 `:691-693` 一段散文注释与 boot 报告的一枚 posture，**没有任何尺在钉**（`docs/evidence/s1/246-resident-task-source-v2.md:95` 那段末句已经说过这件事）。

### 档位（本腿不替人拍板）

**〔备用路（形状上成立），但要人拍〕**：它服务的是「宿主自发起卡」这一形，与今天已接的「工具调用发起卡」是两条不同门路；注释已自陈未接线；`residentCard` 零产码构造点是它离落地还差的那一格。禁 死路的形状证据只有「零调用者」一枚，而那一枚在票 246 的语境里是**设计选择**（`246-a1` 交件时明写「R-2：`AskOnTaskRoot` 零生产调用者 ⇒ 归票 228 后续片＋语音链」）。

---

## 2. `Gate.Replay`（`internal/agent/approval/gate.go:755`）

### ① 调用形状尺（现跑）

```
git grep -nE "\.Replay\(" HEAD -- ':(exclude).scratch' '*.go'
```
读数 **rc=0，命中 1 枚**＝`internal/agent/approval/queue_test.go:283`（`g.Replay(context.Background(), p.CorrelationID, testTask)`）⇒ 剥掉测试＝**0 枚产码调用者**。

```
git grep -nE "[A-Za-z0-9_]\.Replay\(" HEAD -- ':(exclude).scratch' '*.go' ':!*_test.go'
```
读数 **rc=1（空）** ⇒ 整族结论（非抽样）。

```
git grep -nE "\bReplay\b" HEAD -- ':(exclude).scratch'
```
读数 **20 枚**，其中 `.go` 面只 **4 枚**＝`gate.go:751`（注释）／`gate.go:755`（定义）／`queue_test.go:283`（调用）／`queue_test.go:285`（错误文案）；
其余 16 枚全在 `docs/**`。**命中树里没有 `design/**`、`frontend/**`**（本腿不读其内容，只记「零命中」这一枚事实）。
⇒ **接口派发形被闭合**：`ui.go:143-162 NativeAPI`（只有 `Allow`/`AllowSession`/`Reject`）与 `ui.go:167-171 PanelAPI`（只有 `Reject`/`Head`/`View`）**都不声明 `Replay`**，
所以它只在具体类型 `*Gate` 上可达，不存在「谁实现了这个接口」那种看不见的派发路。答案侧那枚漏斗 `Replies`（`replies.go`，非测试方法共 **18 枚**，尺：`git grep -cE "^func \(r \*Replies\)"`）里**没有 Replay 路由**；
控制台动词条 `cmd/wisp/approval_reply.go:566-584` 只有 `yes`/`session`/`no`/`veto`/`panel-no`/`panel-yes`/`always`/`head`/`view`/`help`。

它的唯一下游：
```
git grep -nE "func \(q \*Queue\) replay|\.replay\(" HEAD -- ':(exclude).scratch' '*.go'
```
读数 **2 枚**＝`gate.go:759`（`g.q.replay(corr)`，住在 `Gate.Replay` 体内）＋ `queue.go:598`（定义）
⇒ `Queue.replay` 与 `Gate.Replay` 是**一条两格的死链**，`queue.go:69` 那枚 `replayOf` 字段（注释「一键重放」）的唯一读者也在这一支里。

### ② 编译进产物还是根本没进

- `internal/agent/approval/gate.go` 第 1 行＝`package approval`，**没有任何 `//go:build`** ⇒ **全 GOOS／全 tag 档都参与编译**，与 `AskOnTaskRoot` 那枚「windows 才有」的射程**不同**（票面禁 不许把两枚压成一格，这就是理由之一）。
- 所在包被产码 import（尺：`git grep -nE "agent/approval" HEAD -- ':(exclude).scratch' '*.go' ':!*_test.go'`）：`cmd/wisp/approval_reply.go:69`、`cmd/wisp/resident_approval_windows.go:61`、`cmd/wisp/resident_task_source_windows.go:88` ⇒ 包本身进得了产物；`Replay` 是导出类型上的导出方法，编译面无条件在内。
- **链接器是否把它剪掉＝判不动**（要 `go build` + 符号表读数，本腿禁跑）⇒ 缺的那一发读数：`go build` 产物上的 `nm`/`go tool nm` 之类；本腿禁 不补、不猜。

### ③ 谁在什么条件下会需要它

读函数体 `:755-767`：`admitCheck(taskID)`（D47，注释 `:228` 写着「is D47's guard on both routes」，其余两枚调用者 `gate.go:253`/`:503` 是 L1/L2 的正常签发路）→ `q.replay(corr)` → 给带回的 `Decision` 盖上 `taskID` → 记一行审计 `approval: replay %s -> %s tool=%s` → **`_ = ctx`**（⚠ 形状事实：签名收 ctx 但今天**整支不透传**，`q.replay` 也不接 ctx）。

**它是冻结契约的实现半件**，不是无主代码：注释 `:751-754` 具名 **C18 一键重放**，`docs/PLAN.md:1368` 的 C18 行原文含「**拒绝后任务 root ctx 不取消，可一键重放**」（`PLAN.md:3210` 第四轮定案处再写一次）。
⇒ 需要它的条件＝**任何一个能给「重放」按下键的入口**：原生侧托盘/卡片（`NativeAPI` 那一侧，按 C18/F2 只有原生侧能发「允许」）、或控制台动词（今天条上禁 没有这一枚）。
⚠ 本腿**判不动**的一发：S7 切片（`PLAN.md:1430`/`:1444`）把 C18 的完整形态排在 S7，而 `queue.go` 是票 146 AC#4 的零字节面（`docs/evidence/s1/146-liveapprovals-r1-accept-r2.md:416` 已经为 `replay` 那枚同族出口开过一票）——
「C18 的一键重放今天算不算已交付」需要人读 `PLAN.md` 与票 146/278 的边界，禁 不是本腿能裁的。本腿只给形状：**契约文字要它、产码没人调它**。

### ④ 如果不接：删 vs 留注释，各自代价一句话

- **删**：等于把 **C18 那一半实现**连同 `queue.go:598`、`queue.go:69 replayOf` 的读者一起抽走，并让 `queue_test.go:283` 当场编不过；这**触契约面（C1-C32）＝人工批准**，不是代码整理。
- **留**：代价＝它和 `docs/evidence/s1/161-doorbell-census-r1.md:266` 那族「答案侧 0 枚生产调用者」（`Native()`/`Panel()`/`Veto()`/`DecideFromNative()`/`DecideFromPanel()`/`SetLoaded()`/`Replay()`）一起挂着，单独给它写「勿删」会把一族的整体接线票（228-a5 托盘那一支）拆散。

### 档位（本腿不替人拍板）

**〔备用路（证据比 AskOnTaskRoot 更硬：有契约文字背书），但删／留要人拍〕**：零产码调用者这枚负向句成立（rc=1，整族）；
可它与 C18 的契约文字绑在一起，禁 属「死代码该删」那一档的证据不足。唯一**判不动**＝链接器是否真把它编进产物（本腿禁跑 build）。
