# 155 的三块"无人裁"——**非实现者**对抗验收 r1（复判：那三块是量出来的，还是推出来的）

- 派单：`.scratch/wisp/dispatches/2026-09-26-123x-accept-155-r1.md`（时刻 2026-09-26 12:3x）
- 被验交付物：`docs/evidence/s1/155-three-unjudged-cells-r1.md`（385 行）＋ `.scratch/wisp/probes/155/`（12 枚）
- **共同锚点＝`bb3a7a1`**（155 实现程最后一枚 commit）。本文件所有"它说 X"的 X 一律按 `git show <锚点>:<路径>` 取版，**不读脏工作树**。
- 本程身份：**验收程（非实现者）**。三格复判＋票面 4 枚框（AC#1..AC#4）逐枚裁。
- 本程写面只有两处：本文件 ＋ `.scratch/wisp/probes/155-accept/**`。零生产码、零票面、零台账、零 `-done`、零 push。
- 本程时刻：进场 `2026-09-26 12:30 +08`（`date` 现取），收口见 §11。
- **一句话总裁决（细节在下面各格）**：三块的答案**都是量出来的、不是推出来的**，本程换了三把不同的尺（编译器、JSON 键级解析、逐枚 commit 名册）各复到同一串数；
  要更正的是**四枚宣称的枚数与两句"理由"**，不是三块的答案。可翻勾性那一格判**不成立**（照它的表还不能翻勾，缺三样，最小闭合集合在 §4）。

---

## 0. 进场反查：派单给我的四条前提，逐条现量（对不上盘上的我按"不成立"写）

凭据＝`probes/155-accept/00-preflight.sh` ＋ `00-preflight.txt`、`10-ruler-caliber-and-counts.txt`。

| 派单写的前提 | 本程当轮读数 | 判定 |
|---|---|---|
| ① "它表里的 `run.go:576` 在 `5365cb22` 上是对的，在 HEAD 已漂到 `:604`" | `git show 5365cb22:cmd/wisp/run.go \| grep -n 'res := loop.Run(ctx, task)'` ⇒ **576**；同一发在 `bb3a7a1` 与 `HEAD` ⇒ **604**（该文件在两版之间从 834 行长到 862 行） | **成立**（所以我下面引任何行号都带锚点） |
| ② "`6de3d1c5` 确实存在一枚 commit，真源是 `777d6cc`/`ac7fb00`" | `cat-file -t` 三枚全＝commit；`6de3d1c5`＝`evidence(153 AC#5 第 5 格＋收口)`，且 `merge-base --is-ancestor b23c7f7 6de3d1c5` ＝ **YES**（它是生产改动那枚的**后代**，不是无关版本） | **成立**——但"号存在"不等于"号指对地方"，那一半在 §6 单独裁 |
| ③ "`git grep` 单发会吞命中并回 rc=1" | 同一发零命中命令（`git grep -n 'withTraceTask' 5365cb22 -- .`）**连跑 5 次：全 rc=1、0 命中**；同一把尺在 `b23c7f7` **连跑 5 次：全 rc=0、8 行**；`context.WithValue @5365cb22` 两跑稳定 1 枚非测试命中 | **本程未复现**。它的"零命中"宣称因此**不是**靠这条抖动站住的，是靠 §1 那把编译器尺 |
| ④ 它 §6.2 两条注入点名的三枚路径盘上反查不存在 | `dispatches/2026-09-26-122x-accept-156.md`＝**ABSENT**；`tasks/bca7d7c.output`＝**ABSENT**（`.scratch/wisp/tasks/` 这枚目录**整个不存在**）；`…155-…-continued.md`＝**ABSENT** | **成立**（本程独立复核，见 §10） |

**锚点存在性**：`bb3a7a1 829abe5 a52ded8 9136820 c10f0c4 5f6646a 4c43b68 33c6306 2a303c5 fe13c11 ff550f3 5365cb22 86b0161 b23c7f7 6de3d1c5 777d6cc ac7fb00 b319bab` **十八枚全＝commit**，无一枚假号。

**仪器名反查**（它 §0 那两枚"同名不同物"的 `mut153.py`）：本程按 `git show ac7fb00:.scratch/wisp/probes/153/mut153.py` 与
`git show 777d6cc:.scratch/wisp/probes/153/accept-r1/zzaccept153_ondisk_test.go` **从锚点取版**（不是抄工作树），
`git hash-object` 分别＝`2d5f7607…`／`04cb8d41…`——后一枚与它 §12 名册宣称的 `04cb8d41b97b3772…` **逐字相同**。
它快照里那四枚探针的 `git hash-object` 本程逐枚复到（`f50cfa30…`／`04cb8d41…`／`9b7d5afc…`／`40f734ab…`，**四枚全等**；
注：它那一列的表头写"sha1"，实测**是 git blob 号**（`git hash-object`）而不是裸文件 sha1sum——下一位若拿 `sha1sum` 去比对会误报漂移，本程两把都跑了才敢这么说）。

**共享工作树脏在哪（登记，不碰不还原，也不算进任何零命中宣称）**：`git status --porcelain` 现读 **40 行**——
`.scratch/wisp/probes/152/my152.py` ` M`、`design/**` 16 枚 `D`/`M`、`docs/evidence/s1/152-…-accept-r1.md` ` M`、
`frontend/**` 与 `.scratch/wisp/probes/154/**` 多枚 `??`、另有**本程自己**的 `?? .scratch/wisp/probes/155-accept/`。
进场当刻 HEAD＝`1cde9fa`，本程跑第一发时 HEAD 已换到 `e9a1db0`（票 154 与编排者在推）——**被验版本＝`bb3a7a1`，不是工作树**。

---

## 1. 格①（票面 AC#1① ／ 153 的 AC#2 问①）复判：改之前 `Compress` 拿不拿得到 taskID？

**它的答复**：拿不到，三条通道（参数表／接收者／ctx）逐条空；正控＝同一把尺打在 `b23c7f7` 命中 **12 行**。

**本程换的尺＝编译器本身**（不是 grep）。做法：在它**没参与写过的**第二棵快照里，往 `Compress` 函数体第一行插一发引用，
让 `go build ./internal/agent/` 判"这一发在不在作用域里"。编译不过 = 量到"不存在"；同一发在改后版编得过 = 这把尺能答"在"。
凭据＝`probes/155-accept/20-compiler-ruler.py`、`20-compiler-ruler-before-5365cb22.txt`、`21-compiler-ruler-after-6de3d1c5.txt`。
快照＝本程自己的 `mksnap.py` 从 `5365cb22`（1344 枚）与 `6de3d1c5`（1452 枚）逐 blob 取版、逐枚重算 sha1 对拍，`mismatched=0`。

| 通道 | 那一发（插在 `func (c *Compressor) Compress(` 之后） | 改前 `5365cb22` | 改后 `6de3d1c5`（其 `internal/agent` 与 `b23c7f7` 逐字节同版，§6 现量） |
|---|---|---|---|
| —（仪器的正控） | `_ = ctx` | **BUILD-OK** | **BUILD-OK** |
| —（仪器的正控） | `_ = c.b` | **BUILD-OK** | **BUILD-OK** |
| **① 参数表／作用域内任何名字** | `_ = taskID` | `compress.go:129:6: undefined: taskID` | `compress.go:175:6: undefined: taskID` |
| **② 接收者字段** | `_ = c.taskID` | `c.taskID undefined (type *Compressor has no field or method taskID)` | 同一枚错 |
| **③ ctx 上有没有可读的任务身份** | `_ = traceTaskID(ctx)` | `undefined: traceTaskID` | **BUILD-OK** |

⇒ **格① 判〔成立〕，而且比它自己那版更硬**：三条通道在**编译器眼里**逐条是空的；
且第三行的两列对照顺手量到一件它没写的事——**改后版里 `taskID` 与 `c.taskID` 仍然编不过**，
所以"把 id 送进那条痕"这件事**今天唯一的通道就是 ctx**，不是"我们选了 ctx"而是"另两条通道到改后为止也还是空的"。
这一发同时否掉了 153 票面那条"要不要给 `Compress` 加参数"的备选读法（加了就得动 `compress_test.go` 的调用点，那是另一笔账）。

**它的正控数不成立（一枚枚数错，方向不影响结论）**：它宣称"同一枚 grep 表达式打在 `b23c7f7` 命中 **12 行**"。
本程**逐字复到它那一条命令**（`git show b23c7f7:internal/agent/compress.go \| grep -n 'withTraceTask\|traceTaskID\|traceTaskKey\|ctx.Value'`）
＝**11 行**（行号 `50 134 147 148 150 152 156 159 160 161 237`），它自己的探针 `01` 第 6 节**贴的也正是这 11 行**。
本程另找了三把可能给出 12 的尺（单文件 `grep -c`＝11、加上 `context.WithValue` 仍＝11、全包 `git grep`＝15 枚行/分文件 11+2+2），**没有一把给 12**。
⇒ 判：**正控本身成立（非零＝尺是活的），但那枚"12"是错的宣称**。这是它这一程**第四枚**按枚数腐坏的数（前三枚它自己改了口：名册 7→11→12、禁改面 13→15；这是名册之外的第一枚没人抓到的）。

**两处顺手复到的（不改判，但下一位会撞）**：
1. 它 §1 末段写改前那条痕是"**九枚键**的字面量参数…，**第八枚就是最后一枚**"——同一句里自相矛盾。现量：键值对 **8 枚**（`tokens_before … history_changed`），
   那条 `Info` 里被引号包住的 token 共 **9 枚**（第 9 枚是消息串 `"agent: history compressed"`）。**句子想说的是"没有 task 位"，那部分是对的**（本程 §3 的 ext-A 第一行读数互指）。
2. 它 §1 表里"改前那版全文 `ctx` 只有 5 处（`:39/:128/:152/:197/:201`）"——本程逐字复到，**5 枚全等**。

---

## 2. 格②（票面 AC#1② ／ 153 的 AC#2 问②）复判：唯一持有者是谁？四句"不算"逐句判

**它的答复**：唯一持有者＝`Loop.run` 的参数 `taskID`（它写 `loop.go:339`）；组合根不是。同帧另有四枚副本，它逐枚给了一句"为什么不算另一个持有者"。

**先钉它引用的行号（按它的锚点取，派单坑①）**——`git show 5365cb22:internal/agent/loop.go` 现读：
`:239 comp: NewCompressor(b, opt.Summarizer, WithLogger(opt.Logger))`、`:321 func (l *Loop) RunAsync(`、`:322 id := newTaskID()`、
`:332 func (l *Loop) Run(`、`:333 return l.run(ctx, newTaskID(), input)`、`:339 func (l *Loop) run(ctx context.Context, taskID, input string) Result`、
`:396 nh, rep, err := l.comp.Compress(ctx, hist)`。⇒ **它表里那七枚行号逐枚复到，无一枚错号**（在 HEAD 上会漂，在它自己的锚点上对）。

**本程换的尺＝同一发编译器探针的第二半**（`probes/155-accept/20-compiler-ruler.py`，插在 `:396` 那一发之前）：

| 那一发 | 改前 `5365cb22` | 改后 `6de3d1c5`（插点 `:399`） |
|---|---|---|
| `_ = taskID` | **BUILD-OK** ⇒ 参数在压缩那一刻的作用域里 | **BUILD-OK** |
| `_ = l.taskID` | `l.taskID undefined (type *Loop has no field or method taskID)` | 同一枚错 |
| `_ = root.ID` | **BUILD-OK** | **BUILD-OK** |
| `_ = res.TaskID` | **BUILD-OK** | **BUILD-OK** |
| `_ = j.taskID` | **BUILD-OK** | **BUILD-OK** |

**组合根那一半（它说"那一发不带 id、cmd 侧全是事后读"）**：本程不复算它的 grep，换一枚它没用的尺——
`git grep -n 'RunAsync' 5365cb22 -- cmd` **两跑皆 rc=1、0 命中**（全仓非测试命中只有 `loop.go:320` 注释与 `:321` 定义本身）。
⇒ 宿主在锚点那版**根本不经过异步那条腿**，所以 `RunningTask.ID` 今天在生产上没有活体；这一发把它的"同步腿根本不经过它"从**断言**升成**读数**。判定：**〔成立〕**，且比它自己给的更硬。

**四句"不算"逐句判（派单要的就是这一列）**：

| 副本 | 它给的理由 | 本程的尺 | 判 |
|---|---|---|---|
| `Root.ID`（`loop.go:345`） | "`NewRootFrom` 内部只 `context.WithCancel`，没有 `WithValue` ⇒ id 没进 ctx" | `git grep -n 'WithValue' 5365cb22 -- internal/observe` **两跑 rc=1、0 命中**；`Root struct{ID;Ctx;Cancel;…}` 现读；§1 那发 `_ = traceTaskID(ctx)` 在改前编不过 | **不是打发**：理由独立复到。但"不算另一个持有者"只在"**能从 `Compress` 里摸到**"这个口径下成立——它确实是同一帧的第二枚持有者（本程 `_ = root.ID` 在 `:396` 之前编得过） |
| `l.current`（`loop.go:353`） | "只有 3 处**读点**：`:348` 注释、`:523`" | `git grep -nE '\.current\b' 5365cb22 -- internal/agent/loop.go` ＝ **3 枚**：`:348` 注释、`:523` 读、`:987` **写**；字段声明 `:189 current *observe.Root` 是第 4 枚提及 | **数与列不吻合**：宣称"3 处读点"却只列 2 枚，且那 3 枚里含一枚写点。**结论仍成立**——`Compressor` 三枚字段（`b/sum/lg`）、六枚方法签名无一枚接 `Loop`/`Root`（现读），够不着 |
| `Result.TaskID`（`loop.go:369`） | "那是**出口**，且**晚于压缩那一发**" | `:369` vs `:396`：同一枚函数体、同一条 `for` 之前 ⇒ **早 27 行** | **打发，而且打错**：赋值时刻早于压缩。"出口"那半句对（它是返回值的载体），"**晚于压缩**"那半句盘上不成立 |
| `RunningTask.ID`（`loop.go:323`） | "它在 `reg.Spawn` 的闭包外，同步腿根本不经过它" | 上面那发 `RunAsync @cmd` ＝ 0 命中两跑 | **不是打发，但它没给凭据**：结论对，理由本程替它钉成了读数 |

**它漏报的两枚同帧持有者**（要登记的就是这一笔）：
① 日志那一发 `j := newTaskJournal(l.opt.Journal, taskID, taskID)`（`loop.go:356`；`taskJournal` 有字段 `taskID`，`journal.go:52/:61` 现读）——本程 `_ = j.taskID` 编译**过**；
② 宿主审批回调 `l.opt.AdmitTask(taskID)`（`loop.go:362`；`Options.AdmitTask func(taskID string) (revoke func())` 现读，`cmd/wisp/run.go:558 AdmitTask: rt.gate.AdmitTextTask`）——id 在这一刻**已经出了包**。

⇒ **对它 §2 收尾那句"精确读法"的判定**：前半句"**造 id 的通道唯一**（包私有 `newTaskID`、全仓 4 枚命中、`cmd` 侧 0）"**〔成立〕**；
后半句"**压缩发生那一刻手里有 id 的只有 `Loop.run` 那一枚参数**"**〔不成立〕**——同一刻至少 `root.ID`／`l.current.ID`／`j.taskID`／`res.TaskID` 四枚编得过，外加一份已递给 gate 的拷贝。
它自己 §2 表里那句"**没有任何一枚住在 `Compressor` 够得着的地方**"才是站得住的落点，本程逐枚复到、判**成立**。

**本格档位**：**问② 的答复〔成立〕**（持有者是 `Loop.run` 的参数、组合根不是，两半都有独立读数）；
**"唯一持有者"这个措辞〔说过头〕**，应收窄成"唯一**在压缩那一发够得着**的持有者"。这一处**不改 AC#2 的定档**——问②要的是"id 从哪来、要不要动契约面"，那两问本程都复到同一答案（不动契约面：`Compress` 签名与 `Result`/`RunningTask`/`Root` 三枚结构体一字未变即把归因闭上，§1 的编译器尺就是这条的反证）；
但措辞必须登记，因为"唯一"两个字一旦被人照抄，下一个读者会以为这帧里只有一份 id，从而在 Warm-window hook 落地时（`DEFERRED(D28-1)`，`loop.go:395` 注释现读）**漏掉另四份拷贝里任何一份该不该带标**这个问题。
