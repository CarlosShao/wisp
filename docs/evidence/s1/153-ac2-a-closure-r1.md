# 153 AC#2 部件 (a) 那两笔待补的**内容级**复核 · r1（非实现者验收程）

派单＝`.scratch/wisp/dispatches/2026-09-26-144x-accept-153-ac2a-closure.md`（存档 `1a04f3d`）。
射程＝**一格**：`docs/evidence/s1/153-trace-lies-unguarded-r1.md` 的 `### B.9`（`:600`）／`### B.10`（`:627`）
**里面写的账对不对**——157 的验收件对它自己 AC#2 那行写的是「成立（**形状**；这两笔不在本程的抽判射程）」，
⇒ 形状已被独立核过，内容没有。本程只买内容这一发。

**本程不做什么**：不翻勾（153／155／157 三枚票面一字未动）、不改任何已入库证据件、不写台账／HANDOVER／注入时间线、
**零仪器**（没跑 `go test`／`gofumpt`／`go vet`／`d22scan`；本程全是 `git cat-file`／`git show`／`git grep`／`grep -c` 级读数）。

---

## 0. 锚点与脏树分离

| 项 | 现量读数 | 命令 |
|---|---|---|
| 进场锚点（被验锚点＝我自己现读的） | `df69bcd` = `df69bcd8db71b2a85853aeb2fc4048a929499194` | `git rev-parse --short HEAD` @ 14:45 |
| 收口锚点（共享树自己往前走过） | `eee7bd9`（14:47 起读到，非我进场那枚） | 同上，每次读数当轮现取 |
| 脏树规模 | `git status --porcelain` **38** 行＝16 删／6 改／16 未跟踪 | `git status --porcelain \| wc -l` |
| 三枚被读文件的 `worktree` 与 `HEAD blob` 是否同物 | **同物**（逐枚 `git hash-object` 对 `git rev-parse HEAD:<路径>`）：153 件 `c837c7b2`、153 验收件 `ea67b2b2`、155 验收件 `1cf96570`；台账 `0a2f81d5` | 见 §附录 命令原文 |
| 派单给的锚点（`9181977`）与我读到的 | 派单写"编排者 14:4x 读到 `9181977`"；我现读 `df69bcd`（它是 `9181977` 的下一枚，＝派单存档之后那笔台账＋停车点）⇒ **按派单自己的规矩用我这枚**，且所有代码侧读数一律钉 `86b0161`／`6de3d1c5`，不读工作树冒充被验版本 | `git log --no-walk` |

**被验两节的来路（当轮现量，不是转抄）**：

```
$ git log --format='%h|%ad|%s' --date=format:'%H:%M:%S' -- docs/evidence/s1/153-trace-lies-unguarded-r1.md
ec815820|13:34:39|evidence(157 第 10 笔 落 153 证据件 附录 B)   ← B.10
7f79638e|13:33:55|evidence(157 第 9 笔 落 153 证据件 附录 B)    ← B.9
b319bab5|11:25:45|docs(票 151 结＝… 新立票 155)                 ← 157 的进场基版（改前那枚 blob）
$ 各版 wc -l：b319bab5=589 ／ 7f79638e=625 ／ ec815820=648 ／ HEAD=648
$ grep -cE '^### B\.9'：b319bab5=0 ／ 7f79638e=1 ／ ec815820=1 ／ HEAD=1 ／ 工作树=1
$ grep -cE '^### B\.10'：b319bab5=0 ／ 7f79638e=0 ／ ec815820=1 ／ HEAD=1 ／ 工作树=1
```

⇒ **本程问"原来缺不缺"一律打在 `b319bab5` 那枚 589 行 blob 上**（与 B.9 自报的"进场现量 `wc -l`＝**589**"同数），
问"现在在不在"打在 HEAD blob 与工作树上各一发。

## 0.1 链条最后一环（只补这一环，不重导）

| 环节 | 现读结果 |
|---|---|
| `153-…-accept-r1.md:831` 的 (a) 逐字档 | 现读 HEAD blob：`**附条件**——两笔待补（…件内缺那笔隐私账；…日志侧 task 与 DB 侧 task_id 不同名无人写）。最小闭合＝证据件补那两行，不动码不动判据` ⇒ **与派单所引逐字一致** |
| 同一件 `:416`（§5.3 档位行） | 现读：`**本格档位：成立，附两笔待补**（①件内缺那笔隐私账；②日志↔DB 键名差异无人写）` ⇒ §5.3 那格与 §12 的 (a) 用的是两枚措辞（"成立，附两笔待补" vs "附条件"），本程**只报读数不判**（已入库证据件不许改） |
| `:403`–`:411`／`:412`–`:414` 的行界 | 现读：第 1 笔 `:398`、第 2 笔 `:403`–`:411`（隐私账）、第 3 笔 `:412`–`:414`（键名）、`:416` 档位行 ⇒ B.9/B.10 引用的行界**对**，157 件 `:598` 那枚"更正工单 `:21` 把第 9 笔写成第 3 笔"的自报也**对** |
| 155 验收件 `:115`／`:166` | 只读核过在场（派单已建立，本程不重导） |

---

## 1. 问一：B.9 那笔隐私账补得对不对 —— 档：**成立**（另记两枚非承重缺陷）

### 1.1 判据（可复算形式）

B.9 说"缺不缺这一行"＝**六词尺在改前那枚 blob 上全 0**；说"这一行该写什么"＝**七枚码级断言**逐枚按它自己钉的锚点复算。
⇒ 我的判据＝同尺重跑，并且**每枚 0 命中配一枚已知正控**（防死尺）。断言清单（B.9 正文 `:602`–`:625` 逐枚抽出来）：

| # | B.9 的断言 | 我的复算尺 | 读数 | 档 |
|---|---|---|---|---|
| A1 | 六词（`隐私`／`privacy`／`keep_transcript`／`私有数据`／`Q-31`／`SealDir`）在**改前件**里各 **0** 命中 | `grep -c` 打在 `git cat-file blob b319bab5:…153-trace-lies-unguarded-r1.md`（589 行版） | `0/0/0/0/0/0` 六枚全 0 | **对** |
| A1b | 同一把尺打验收件＝"7／3／6／—／2／4"（正控，证明 0 不是死尺） | 同尺打 `HEAD:153-…-accept-r1.md` | 我读到 **7／3／6／2／2／4** | **5 枚对、1 枚不符**：`私有数据` 实为 **2**（`:407` 词表自身＋`:408` 一处正文），B.9 写成"—" ⇒ **正控抄错一枚**，非承重（详见 1.3 缺陷① ） |
| A2 | 这行**现在在场**（补上了） | 同尺打 HEAD blob 与工作树 | HEAD＝`5/5/5/3/3/5`，工作树＝同数 ⇒ 六词全部从 0 变非 0 | **对** |
| A3 | 改前锚点 `86b0161` 已有**三枚**非测试行用同键名写同值：`loop.go:373`／`:741`／`:933` | `git show 86b0161:internal/agent/loop.go \| grep -n '"task"'` | 恰三枚，行号 `373`／`741`／`933` 逐字对上，值是 `taskID`／`req.TaskID`／`res.TaskID` | **对** |
| A3b | `86b0161` 真当得起"改前"这个位置 | `git rev-parse 86b0161:…loop.go` vs `b23c7f77^:…loop.go`（码侧真父）＋`merge-base --is-ancestor` | 两枚 blob 同为 `e685a6ef…`，且 `86b0161` 是 `b23c7f77` 的祖先 | **对**（`86b0161` 本身是枚前端提交 `feat(frontend 10 首启引导)` 08:59，但它载的 `loop.go` 与真父**逐字节同物** ⇒ 拿它当锚点不产生误差） |
| A4 | 到 `6de3d1c` 非测试侧带 `"task"` 的行数＝**7** | `git grep -n '"task"' 6de3d1c5 -- internal/agent ':!*_test.go'` | **7** | **对**（本程补一枚拆解：7＝4 枚真写键的行（`loop.go:373/744/936` ＋ `compress.go:238`）＋ 3 枚注释（`compress.go:49/139/232`）。B.9 写的是"行数"，字面成立；读者若把 7 当"七次写键"会读多 ⇒ 见 1.3 缺陷② ） |
| A5 | 交付版 `compress.go:236` 那句"…`[privacy] keep_transcript` does not reach it"逐字在场，且 `:236` 正是"does not reach it"那行 | `git show 6de3d1c5:internal/agent/compress.go \| sed -n '226,246p'` | `:232`–`:236` 是那五字注释块，`:236` = `// history content, so [privacy] keep_transcript does not reach it.` ⇒ 那句**跨 `:235`–`:236`**（"It is a correlation id, not" 在 `:235`），`:236` 落 "does not reach it" | **对**（同一锚点上 `compress.go` 与实现提交 `b23c7f77` 的 blob 同为 `4fcd9a32…` ⇒ "交付版"三名一致） |
| A6 | `keep_transcript` 硬编码 false＋写 true 报错：`schema.go:511`／`validate.go:82-84`／`validate_test.go:136` | `git grep -n 'KeepTranscript' 6de3d1c5 -- internal/config` ＋ 三处取行原文 | `schema.go:511 KeepTranscript bool \`toml:"keep_transcript" default:"false"\``；`validate.go:82-84` 正是 `if c.Privacy.KeepTranscript { return observe.New(observe.ClassConfig, "…writing true is rejected") }`；`validate_test.go:136` 是那条 `{"privacy.keep_transcript=true", …}` 钉 | **对**，三枚行号逐枚对上 |
| A7 | `winsec.SealDir` 在 `6de3d1c5` 非测试代码里**零调用者**（只有定义与注释），命中 6 枚全在 `winsec.go:174,190,220,222` 与 `winsec_windows.go:52,313` | `git grep -n 'SealDir' 6de3d1c5 -- '*.go' ':!*_test.go'` ＋ `':!internal/winsec/*'` ＋ 含测试的全量当正控 | 非测试 **6** 枚，行号逐枚＝所写；`func SealDir` 定义在 `winsec.go:222`（文档注释 `:220`），另 4 枚是注释；winsec 之外非测试命中 **0**；含测试全量 **37** 枚（＝尺是活的，且那 37 枚调用者全在 `internal/winsec` 自己的测试里） | **对** |
| A7b | 这件事今天在 HEAD 上还没变（B.9 未宣称，本程自加） | `git grep -c 'SealDir' HEAD -- '*.go' ':!*_test.go'` | **6**（同形） ⇒ 该断言未腐坏 | **对** |
| A8 | "它落在**票 132 那格未闭合的面上**" | `ls .scratch/wisp/issues/ \| grep '^132-'` ＋票面勾框计数＋台账 `Q-31` 现读 | 票面存在、**无 `-done` 后缀**、未勾 **5**／已勾 **0**；台账 `:3446` 记 `Q-31` 日志算私有数据＝**算、要锁** ⇒ 新建票 132，`:6871` 另有一行"给日志带上 taskID ⇒ taskID 从此落进日志文件…`Q-31`→票 132 那条路至今 `SealDir` 生产零调用者" | **对** |
| A9 | "本票交的那枚 `task` 键把一枚可归因 id 写进宿主日志文件"＋"让**成功边**那一行也进" | 取版读 `compress.go:237-240` 与落点级别 | `:237 if task := traceTaskID(ctx); task != "" { attrs = append(attrs, "task", task) }` 后 `:240 c.log().Info("agent: history compressed", …)`；`cmd/wisp/logsink.go:82 const logSinkLevel = "info"` ⇒ Info 真落盘；`internal/agent/loop.go:399` 的打标点用的是 **run() 的同一枚 `taskID`**（`run` 函数体 `:339`–`:504`，与 `:373` 同函数） | **对**，"同一枚值"这件事我按函数边界复算过（`res.TaskID` 源自 `:369 Result{TaskID: taskID}`，`req.TaskID` 源自 `:647 TaskID: taskID, CorrelationID: taskID` ⇒ 三条改前行与本票新增行写的是同一枚 run 的 task id） |

原始读数全文＝`.scratch/wisp/probes/153-ac2/q1-b9-recheck.txt`、`q1-b9-code-anchors.txt`、
`q1-b9-config-and-sealdir.txt`、`q1-b9-q31-and-seal-route.txt`、`q1-b9-logsink-and-132.txt`、`q1-b9-scope-of-taskid.txt`、`q1-b9-loglevel.txt`。

### 1.2 关键命令原文（复算用）

```sh
# A1 改前那枚 blob 上六词全 0（正控＝同一把尺打验收件）
git cat-file blob b319bab5:docs/evidence/s1/153-trace-lies-unguarded-r1.md > /tmp/wisp153-prefix.md
for w in 隐私 privacy keep_transcript 私有数据 Q-31 SealDir; do printf '%s %s\n' "$w" "$(grep -c -- "$w" /tmp/wisp153-prefix.md)"; done
# A3/A4 键名名册（改前三枚、交付版七枚），全部锚点取版
git show 86b0161:internal/agent/loop.go | grep -n '"task"'
git grep -n '"task"' 6de3d1c5 -- internal/agent ':!*_test.go' | wc -l
# A7 那把锁有没有人调
git grep -n 'SealDir' 6de3d1c5 -- '*.go' ':!*_test.go'
git grep -n 'SealDir' 6de3d1c5 -- '*.go' ':!internal/winsec/*' ':!*_test.go'   # => 0
```

### 1.3 结论与两枚缺陷

**B.9 的隐私账：〔成立〕。** 它承重的那句（"改前件里没有这笔账"＝六词全 0，现版已在场）与它给出的**九枚码级／账级断言全部复算得到**，
且它引的验收程原文（`:403`–`:411`＋`:416`）逐字在场、HEAD blob 与工作树各 1 枚命中（路径存在且那段在）。
⇒ **153 验收件 §12 (a) 那"两笔待补"的第一笔，内容上补对了。**

记两枚**非承重**缺陷（不改档，交编排者入账）：

1. **正控抄错一枚**：B.9 正文 `:608` 把同一把尺打在验收件上的读数写成"7／3／6／—／2／4"，我复算为"7／3／6／**2**／2／4"
   （`私有数据` 在验收件里有 2 行：`:407` 词表自身、`:408` 正文）。157 那枚 commit message（`7f79638e`）里写的是"7/3/6/2/4"，**只有五枚数**对六词 ⇒ 两处抄录彼此也不一致。
   ⇒ 影响面：正控只用来证明"0 不是死尺"，其余五枚已足以证明 ⇒ 不翻档。
2. **7 这一枚的行数语义未拆**：A4 的 7 含 3 枚注释行（真写键的是 4 枚）。B.9 措辞是"带 `"task"` 的**行数**"，字面无误；
   但这枚数被"同一枚键名、同一枚值"那段当分量用，读者会读成"七次写键" ⇒ 属可读性缺陷，不属假账。

---
