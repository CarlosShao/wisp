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
---

## 3. 格③（票面 AC#1③ ／ 153 的 AC#3 句②）复判：摘掉它，外部可见读数到底变没变

**它的答复**：两味变、一味不变——盘上 JSONL 里痕总行数恒 **3**，带 `"task":"` 的行 **2→0→1**（ext-A/B/C），自加的 ext-D（摘 N1）**3/2 不变**。

**本程的复算＝另起一炉**：不用它的快照、不用它的探针、换一把**键级**尺。
凭据＝`probes/155-accept/31-ext-matrix.py` ＋ `31-ext-matrix.txt` ＋ 探针原件 `30-zzaccept155-parsed-ondisk_test.go`。
- 快照＝本程 `mksnap.py` 从 `6de3d1c5` 重取的 1452 枚（`mismatched=0`）；pristine 三枚 `git hash-object` ＝`4fcd9a32a500…`／`3ea1fb8df38a…`／`2357dca86f67…`，与 `6de3d1c5` 的 blob 逐字同（也与 `b23c7f7` 同，见 §6）。
- 变异驱动＝**原样复用** `.scratch/wisp/probes/153/mut153.py`（从锚点 `ac7fb00` 取版，`git hash-object`＝`2d5f7607…`，一字未改；它 `sub()` 要求锚串**恰好一枚**否则 `REFUSED`，所以"施没施上"是仪器自证的）。
- 尺＝**`encoding/json` 解析每一枚落盘记录**，问"`task` 这枚**键在不在**、值是什么、值等不等于**这一次运行**的 `Result.TaskID`"。它和 153/155 那两把都是子串计数（`"task":"` / `task=`），子串计数分不开三形：空值、非字符串值、张冠李戴的 id。
- 用到的 helper 全是**交付里原有的**那几枚（`newHarness`/`withLogger`/`withConfig`/`buildRoundHistory`/`BudgetsFor`/`observe.InitLogWithRegistry`），本程没造新接缝。

| 发 | 施了什么 | 本程读数（盘上） | 它的读数 | 判 |
|---|---|---|---|---|
| 正控①（未变异，交付用例侧） | — | **RUN=7 / FAIL=0** | 7/0 | **复到** |
| 正控②（尺有牙） | `drop-attr` | **RUN=7 / FAIL=2**＝`TestCompressionTraceCarriesTheOwningTaskID`＋`TestCompressionTraceNeverInventsATaskID`（红名逐字同） | 2 枚红，同名 | **复到** |
| **ext-A** | 无 | records=**3** / `task` 键在=**2** / `keys=11` 与 `keys=12` 之差恰是那一枚键 / **`leg3_match=true`** | 3/2 | **复到** |
| **ext-B** | 摘 N3（痕不再写属性） | 3 / **0**（第二回独立再跑仍 3/0，uuid 每回不同＝预期噪声） | 3/0 | **复到** |
| **ext-C** | 摘 N2（loop 不再打标） | 3 / **1**，且 `leg3_match=false` 而 `leg3_taskid` 是一枚活 uuid ⇒ **归因消失、运行没消失** | 3/1 | **复到** |
| **ext-D**（它自加的那味） | 摘 N1（AC#1 那枚判据用例） | 3 / **2** ＝与 ext-A 同形 | 3/2 | **复到** |
| 句①那一形 | ext-D 同一棵树上再施 `m5` | **RUN=6 / FAIL=0**、`ok …/internal/agent` ⇒ 摘掉证人⇒逃逸回来（153 的 X2 的本程版） | 6/全绿 | **复到** |
| 还原自证 | 每发前 `restore`、收口 `git hash-object` 三枚 | 三枚＝锚点 blob、`RESTORED-IDENTICAL True×3` | 同 | **复到** |

**它 §3.1/§3.2 那枚"仪器形状"本程独立复到（这一条最该被审，因为它决定了 ext-C 能不能只靠复用仪器交卷）**：
把 **153 前一程那枚两腿到盘探针**（`zzaccept153_ondisk_test.go`，本程从锚点 `777d6cc` 取版）放进**本程**的快照并排跑——
`ext-C` 下它 **PASS**（两行、第二行照旧带 `"task"`）＝**结构上看不见 `unloop`**；`ext-B` 下它 **FAIL**（`:70 tagged record lost its task on disk`）。
⇒ 它那句"只复用已有仪器**不足以**答 ext-C"**〔成立〕**，而且它没有拿那枚 PASS 当"没变"交卷——这一处是它这一程最实的一次不放水。

**派单要我回答的那一问：ext-D 这一发推翻的是 153 的哪句原话？——答：哪句都没推翻。**
153 实现件 §4.2 关于 N1 的原话是：**"N1 那句的答案是'不会'：它是判据、不是输出面……（本程未为它造外部读数，可造的那一侧只是'少一枚证人'）"**。
它**自己先声明了没有读数**，所以 ext-D 不可能推翻它——**ext-D 是把一枚空缺的读数补上，顺带证实那条推理成立**。
⇒ 按派单给的二选一：这一发**不是"纠正"，也不是"多余"**，它是**补档**；而 155 自己的措辞（§3.4"实现件那条推理**成立，而且现在有读数了**"）与此同调，**没有越权把它写成纠正**。真正"被推翻"的是另一句：**153 前一程 §8 第 3 条自报"没复算 ext-A/B/C"**——那一格今天从"没复算"变成"两程各自复算过"。

**"更强的尺"这一说法要收窄一档（本程现量到的一处口径）**：它写"它量的是测试进程 stderr 上的 `task=`，本程量的是宿主日志文件里的 `"task":"` ⇒ 强一档"。
盘上：宿主那条默认 logger 是**一枚 fan-out**——`cmd/wisp/logsink.go:153` `slog.SetDefault(slog.New(teeHandler{primary: p.Handler(), mirror: slog.NewTextHandler(os.Stderr, …)}))`，
文件腿是"redactHandler 套 JSON writer"、stderr 腿是镜像。⇒ **"更强"成立**（只有文件腿穿过脱敏器与滚动写手并回到盘上被重读，stderr 腿不证明持久化），
**但不能读成"两枚互相独立的通道"**：同一条记录分两路，若脱敏器哪天改写了键名，两腿会一起改、两把尺会一起哑。
它自己在 §5 第 1/2 条把"宿主那条真链没跑""进程默认 logger 那一档未证"两笔留着，这一处**没有**冒充跑过——登记为**成立附一条口径收窄**。

**本格档位**：**〔成立〕**。句②三味各自都有到盘读数，本程换尺复算到同一串数（2/0/1 与 3/2），且"恒 3"这一条地基复到。
两程读数**不互相抵账**（本程的读数另立一栏，与它的一致只是"这把尺没有只量出它想看的"）。
---

## 4. AC#2（它 §4 那张"三块 ↔ 153 哪一枚 AC"推进表）复判：能不能只照它翻勾

**两把尺分开答**（票面那把、派单那把）：

- **按票面 155 AC#2 的字面要求**——"交件里必须有一行'这三块各自把 153 的哪一枚 AC 从'未裁'推到'可定档'"：§4 表有**三行**，逐枚点名 `AC#2 的 (b) 问①`／`AC#2 的 (c) 问②`／`AC#3 的整枚句②`，每行还带"本程之后还差谁"。⇒ **〔成立〕**。
- **按派单那把更强的尺**——"这张表读完，我能不能**只照它**翻勾而不回头读 153 的表？"⇒ **〔不成立〕**。缺三样，逐样点名（凭据都在盘上）：

| # | 缺什么 | 现量出处 |
|---|---|---|
| 1 | **AC#2 是六件框，它只推到两件**。153 的验收件 §12（`153-trace-lies-unguarded-r1-accept-r1.md:827`）把 AC#2 拆成 **(a) 交付物／(b) 问①／(c) 问②／(d) 问③／(e) 不许造假 ID／(f) 停手条款**，并逐件给了档：(a) 成立附两笔待补、(d) 成立、(e) 成立、(f) 成立。155 §4 只写"(a) 那两笔记账级待补仍挂在 153 名下"、"(d) 前一程已裁"，**没把 (a)(e)(f) 的档位抄进来** ⇒ 拿着 155 这张表的人**无法**给 AC#2 定档，仍必须回读 153 验收件 `:831-:836` | 本程现读两枚票面：153 票面 AC#2 原文只有"三问逐问答"＋两条 ⚠（`:26-:29`），**"六件"这套字母不在票面上** |
| 2 | **字母表的出处指错了一份文档**。155 §1 写"**票面**拆出的 6 件里那一件 (b)"——票面没拆过 6 件，拆 6 件的是 153 的**验收件** §12。指针走偏＝把下一位领到错的文件，正是本票立票的那一族错（153 续程 §13 那组断头指针同源） | 155 证据件 `:71`（"票面拆出的 6 件"）与 153 票面 `:25-:29` 并排 |
| 3 | **AC#3 那一行给了指针、没给档位**。它写"句①（前一程 §3.2 已裁、本程 §3.3 正控②与 §3.4 各复到一次）"——读者仍需回 153 验收件才知道句①被定成什么档；句②它写"整枚"，可 AC#3 那枚框今天**该不该勾**仍要两句的档并排 | 155 §4 表末行"本程之后还差谁"列 |

**最小闭合集合（三条都只加行、不动任何判据、不改 153/155 票面——本程也不改）**：
① 推进表加一列"该枚 AC 全档"，把 (a)／(e)／(f) 与句① 的档各抄一行，出处写成"153 **验收件** §12 `:831-:836`"这种带行号的粒度；
② 把"票面拆出的 6 件"改成"153 验收件 §12 拆出的 6 件"；
③ AC#3 行的"还差谁"补一句"句① 档＝前一程所判〔成立〕"。

**顺带钉住一条它守住了的**：它 §4 那句"本程一枚框都不勾、不自立档位"是**真的**——
现读：153 票面最后一次改动＝`b319bab`（编排者 11:25）、155 票面最后一次改动同样＝`b319bab`，两枚票面 155 那一程**一次都没提交过**；
153 未勾框 **2**、155 未勾框 **4**、`-done`（153/155）**0**。凭据＝`probes/155-accept/41-contract-axis.txt` §6。
⇒ **本格档位：AC#2〔成立〕（票面字面）＋〔不成立·按派单强尺〕，缺上面三样，都是加一列/改一句的量级。**
---

## 5. AC#3（契约轴零字节）复判：名册现算，别抄它表里的数

凭据＝`probes/155-accept/41-contract-axis.sh` ＋ `41-contract-axis.txt`（每支两跑记 rc）。

**它自己的宣称**：commit 枚数 7 枚（收口 §9）／名册 **12** 枚（"§8 那枚 11 又是一次写它之后又多了"）／禁改面 **15 支**逐支 0／正控在 `b23c7f7` 命中 3 枚。它这一项**自报改口三次**（名册 7→11→12、支数 13→15）。

**本程现量**：

| 那一项 | 本程的尺 | 读数 | 判 |
|---|---|---|---|
| 属于 155 的 commit | 逐枚过 `829abe51^..bb3a7a1e`，判据＝该枚 commit 的文件名册里有 `docs/evidence/s1/155-` 或 `.scratch/wisp/probes/155/` | **10 枚**（同一区间里 **3 枚是票 154 的**：`352d3d8`／`76c89b2`／`304199f`，被同一把尺筛掉） | **成立**（它 §9 数成 7 枚＝写它时只有 7 枚，之后又长了三枚——与"名册枚数"同一族过期） |
| 名册并集 | 逐枚 `git show --pretty=tformat: --name-only` 相加 `sort -u` | **13 枚路径**＝证据件 1 ＋ `probes/155/` 下 12 枚 | **它的"12"再次偏小一枚**；**结论不受影响**（第 4 枚同族数错，见下） |
| 名册越界 | 名册减掉它自己那两枚写面路径 | `grep -v` **rc=1、零枚** ⇒ 13 枚**全部**落在声明的写面内 | **成立** |
| 禁改面 | 16 支（它的 15 支 ＋ `^cmd/` ＋ `^third_party/`，去重后仍含 `^internal/`），逐支两跑记 rc | **逐支 0**（每支 rc=1/1 两跑一致） | **成立** |
| 证据件那一支 | `grep '^docs/evidence/'` 于名册 | 只有 `155-three-unjudged-cells-r1.md` 自己一枚 ⇒ 票面 AC#3 那句"票 153 与票 155 之外的**任何证据件**"零命中 | **成立** |
| 正控 | `git show --name-only b23c7f7` ⇒ `^internal/` 命中 | **3 枚**（`compress.go`／`compress_trace_test.go`／`loop.go`）＝它宣称的那 3 枚 | **成立**（那把尺能命中，上面的 0 不是死尺） |

**本程自己撞出来、要回给派单的一枚仪器坑（这条会**假报违规**，方向与它那一族"数小了"相反）**：
派单 §26 给的尺 `git log --pretty=tformat: --name-only <它的十枚 commit>｜sort -u` **照字面跑不界住历史**——
多 rev 而不带 `--no-walk`/区间时，`git log` 会顺每条祖先链把**整仓历史**的名字都倒进来：本程现量＝**1730 枚路径**，
里面 `internal/**`、`frontend/**`、`design/**` 全都在 ⇒ 拿它去扫禁改面会**报出一堆不存在的违规**（票 153 续程量到的"区间尺读到 2 枚 frontend 命中、commit 尺 0"就是同一枚形状的另一种方向）。
加 `--no-walk` 后本程现量 **13 枚**，与逐枚 `git show` 并集 `diff` ＝ **IDENTICAL**。
⇒ **给编排者的一句话**：这一格**只有 commit 级尺能用**这一条被两程各自独立钉到（153 续程钉过一次、本程钉第二次），派单里那把尺建议改写为 `git log --no-walk …` 或逐枚 `git show`。

**本格档位：〔成立〕**——契约轴零字节这件事，本程用完整名册（13 枚）重扫 16 支禁改面，逐支 0、正控能命中 3 枚、名册不越界。
它那三枚改口的数**没有一处把"0 命中"撑成假话**（每改一次口都是"重跑仍全 0"，本程独立重跑证实），但**枚数宣称四次偏小**这件事要记：它不是造假，是**写它之后就过期**——而 AC#3 的判据输入恰好就是那枚名册，所以"每次收口重算"这件事必须是**规矩**，不能是**自报**。
---

## 6. AC#4（活树自证）复判：那棵 blob 精确快照树，等不等于它读数所依据的那一版

**它的宣称**：快照＝`git ls-tree -r 6de3d1c5` 逐 blob `git cat-file --batch` 落盘，`written: 1452 mismatched: 0`；三枚被测文件另用 `git hash-object` 与 `6de3d1c5` 的 blob 逐字对拍；全程未用 `git archive`；开测前先跑正控。
派单要点名判的是**另一件事**：`6de3d1c5` 这枚号**真在**（§0 已复核），但它**指没指到读数那一版**。

**第一层：树＝那枚号（本程另写一炉取版器复算，不用它的脚本）**
`probes/155-accept/11-mksnap.py`（自己实现：`ls-tree` 出 sha 清单落文件 → `cat-file --batch` 走**文件重定向** stdin → 逐枚本地重算 `blob <len>\0`＋数据的 sha1 对拍）：
`anchor=6de3d1c5 blobs=1452 / written=1452 mismatched=0` ⇒ **它那枚数被第二把实现复到**。
同一把尺另取一枚锚点做对照：`anchor=5365cb22 blobs=1344 / written=1344 mismatched=0` ——**这不是一个背下来的数**（两棵树的枚数不同、都各自对拍通过）。

**第二层：它那棵快照（盘上仍在）此刻还纯不纯**
`mksnap.py verify /d/tmp/wisp155/snap manifest-6de3d1c5.json` ⇒ `files=1452 drift=0 missing=0 extra_in_tree=4`；
`mksnap.py cmp` 两棵目录树 ⇒ `shared=1452 byte-differ=0 onlyA=0 onlyB=4`。
那 4 枚多出来的**全是 `internal/agent/` 下的 `_test.go`**：`probe153_line_test.go`(`9b7d5afc…`)／`probe153_loopleg_test.go`(`40f734ab…`)／`zzaccept153_ondisk_test.go`(`04cb8d41…`)／`zzprobe155_ext_ondisk_test.go`(`f50cfa30…`)——
`git hash-object` 四枚与它 §12 名册**逐枚全等**，前三枚又与**锚点原件**（`ac7fb00`／`777d6cc`）全等。
⇒ **它交件时那棵树没有一枚变异落地、没多一枚非测试文件、没少一枚文件**（这条是它 §3.3"还原自证"想说的东西，本程独立钉住）。

**第三层（派单那一问的正身）：`6de3d1c5` 等不等于读数所依据的那版**

| 那一发 | 本程现量 |
|---|---|
| 三枚被测文件在两枚锚点上的 blob | `compress 4fcd9a32a500`／`loop 3ea1fb8df38a`／`compress_trace_test 2357dca86f67`——**`b23c7f7`、`6de3d1c5`、`bb3a7a1` 三枚锚点上逐字相同** |
| `b23c7f7..6de3d1c5` 在 `internal/agent` 上动了什么 | `git diff --name-only … -- internal/agent` ⇒ **0 枚**（整目录没动） |
| `internal/agent` 的**依赖闭包**（`go list -deps` 与 `go list -test -deps`，`GOPROXY=off` 现取） | `winsec/secret/observe/risk/config/llm(+golden/openaichat)/plugin/memory/agent`；**闭包里没有 `internal/tools`**（`grep -c` ＝ 0，两向都是 0） |
| 闭包里每一包在该区间的改动枚数 | 逐包 `git diff --name-only b23c7f7 6de3d1c5 -- <pkg>` ⇒ **全部 0** |
| 该区间在 `internal/**` 上唯一动过的一枚 | `internal/tools/bridge.go`（16 进 3 删，其中 4 行非注释＝票 151 的 `CloseTask` 日志）——**闭包外** |
| 祖先关系 | `merge-base --is-ancestor b23c7f7 6de3d1c5` ＝ **YES** |

⇒ **判〔成立〕**：那棵树**就是**它读数所依据的那一版——不是"号能解析"这种弱判据，而是"**`internal/agent` 整目录与其测试二进制的依赖闭包在两枚锚点之间逐包 0 改动、三枚被测文件同 blob**"。
它没指错版本，**但它也没证过这一层**：§3.3 那行只证了"树＝`6de3d1c5`"（第一层），"`6de3d1c5`＝交付那一版"这一半本程补上（第三层）。
⇒ 正确的口径写法应该是"`b23c7f7` 的交付码，取自其后代 `6de3d1c5`，闭包内逐字节同版"——它写成"快照＝`6de3d1c5`（blob 精确树）"少了后半句。这不是读数错，是**射程声明不够**（下一位若拿同一棵快照去跑 `cmd/wisp`，就会踩到那枚 `bridge.go` 差异）。

**派单坑#4 的两条本程自己复现**：
① **`git archive｜tar -x` 在本仓确实不是纯净树**——现量：`.scratch/wisp/probes/139/ac1-readers.log` **blob 1643 字节 vs archive 落盘 1679 字节**（`.gitattributes:1 = * text=auto`、`core.autocrlf=true`、`.log` 无 `eol=lf` 规则；同一枚文件**工作树里是 1643**，脏的是 archive 这条路，不是工作树）。凭据＝`probes/155-accept/13-archive-pitfall-demo.txt`。
② `third_party/**` tracked＝**0**（`git ls-files third_party \| wc -l`）⇒ blob 树里没有那三枚 dll；本程与它一样**一行 `cmd/wisp` 都没跑**，所以这条对本程**不构成"仪器没跑到"的第三种形状**（它 §3.3 已经自己把这条挡在门外了，本程复核一致）。
③ `-overlay`／`-cover*`：本程**两把都没用**（文件物理落盘）；`-race` **一行没跑**——所以本程**没有**任何一格的答案建在那枚可能 `0xc0000374` 的跑上；判红绿一律只认 `--- FAIL:` 行，判"跑到没"只认 `=== RUN` 枚数（见 §3 那两列）。

**本格档位：〔成立〕，附一条射程声明要补**。
