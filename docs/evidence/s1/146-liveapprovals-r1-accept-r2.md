# 146-r1 对抗验收（第二轮 r2）—— 修复程补上的那两味尺子：裁决表

- 被验对象：票 146 修复程那一修（`internal/agent/approval/ticket146_liveapprovals_backing_test.go`
  ＋ 两枚 docs 件 ＋ 票面 Progress log 一段）
- 被验锚点：**`74255b9`** ＝ 本程 `git rev-parse --short HEAD` 现量、`git cat-file -t 74255b9` ＝ `commit`。
  本程另外用到的一枚 sha 全部现量：`b694378` → `b6943788775e3cd5fa08da122da45a7fe0b30738`（代码修），
  `b694378^` → 取"改前的旧尺"，`git diff --stat 7f75ff3 b694378^ -- internal/agent/approval internal/tools`
  ＝ **空** ⇒ 旧尺在两枚锚点上逐字节同一，"改前 vs 改后"这一对比不必在两棵树之间跳。
- 本程身份：**第二轮对抗验收程（非实现者）**。派单单的立场设定：简报里编排者写的每一句都是未验证断言。
  **复算结果：简报那五行变异对照表逐字成立**（下表），它自陈的那处仍瞎也成立（§4）。
- 只裁两格：**`AC#1`、`AC#2`**（＋ 修复程自己新引入的东西）。`AC#3/AC#4/AC#5` 第一轮已判成立，本程未重做。

## 0. 本程用了哪几棵副本（复算入口；全部仓外、只建不删）

根目录 `D:\tmp\wisp-146acc2\`。建法一律
`git -C "<repo>" -c core.autocrlf=false -c core.eol=lf archive <锚> | tar -x -C <目标>`（**两个 `-c` 都带**）。

| 目录 | 内容 | 用来量什么 |
|---|---|---|
| `snap` | `74255b9` 纯净（代码一字节未改）＋ 后追加一枚本程探针件 | 健康对照 34/0/1、"再往下一层"探针 |
| `snap0` | `74255b9` 纯净、**未放任何探针件** | `c-*` 那九发的母本（见 §5 的一处自纠） |
| `snappre` | `b694378^` 纯净 ＝ **旧尺** | 证明修前三发全绿、旧尺本来覆盖过什么 |
| `snapgates` | `74255b9` 纯净，跑完门禁后放入一枚禁符正控件 | 门禁四发＋d22scan 正控 |
| `v-ctrl v-f v-ma v-e v-comb` | ①那五行 | 复跑修复程的五行表 |
| `v-ma2a v-ma2b v-ma2c` | ②非空夹具三形 | 判"第 9 枚被哪一味堵住" |
| `v-ptrkeep v-ptrdrop v-ifacedrop v-chandrop v-funcdrop v-nodefault v-arraykeep v-noarray` | ③单点回退 | 五支各自撤一发 |
| `v-nilpaths v-nilpaths-noguard` | 夹具完备性守卫的正控／撤守卫 | 判那枚新守卫承不承重 |
| `pre-f pre-ma pre-e pre-comb pre-ma2b pre-nilpaths` | 旧尺 ＋ 同一批变异 | 改前对照 |
| `c-bside c-drop-{8 枚}` | 新尺 ＋ 摘修法／逐枚撤克隆 | 承重＋覆盖面（本程自量，不抄第一轮） |

**副本行号会漂**：本程往测试文件里插过行的那几发（`v-ma2*`、`v-ptrkeep`）红句行号与锚点件相差 ＋1／＋2，
下面逐处标了换算。判红绿一律只认 `^--- FAIL:` 那一行，计数一律带 `=== RUN` 枚数防 panic 吞读数。

**一条 CRLF 尺子自纠（先说，免得下一位照抄）**：派单单要求复算"两个 `-c` 都带"。本程先用
`grep -rl $'\r' --include='*.go' internal cmd | wc -l` 量 ⇒ **455／455 全命中**，那是一枚**假阳性**
（Git-Bash 的 GNU grep 在文本模式下把裸 CR 模式匹配到每一行）。同一批文件用 `od -c` ＋ 逐字节 node 计数复量
⇒ `pending_read.go` / `ticket146_liveapprovals_backing_test.go` / `internal/tools/gate.go` / `cmd/wisp/main.go`
**CR=0、CRLF=0**。⇒ 结论：**副本干净**，但**那把 grep 尺子在 Windows 上不能用**，CR 只能按字节量。
（第一轮验收件 §14 第 ① 条登记过同族"尺子错"，那是 `od` 输出上数 `\r`；本程这发是反方向的同一枚坑。）

---

## 1. 打①：修复程那五行——**原样复跑，五行全部对上**（本轮第一凭据）

母本 `snap`（健康）／`snappre`（旧尺）。命令逐发
`go test -count=1 -v ./internal/agent/approval/`（**逐包单跑**，未与 `./cmd/wisp/`、`./internal/panel/` 合跑）。
读数文件 `out-<变体名>.txt` 留在根目录。

| 变异 | 旧尺（`snappre` ＝ `b694378^`）本程量 | 新尺（`74255b9`）本程量 | 修复件 §3.6 报的 | 对上否 |
|---|---|---|---|---|
| 健康 ctrl | — | **34／0／1、RUN 54**、rc=0 | 34／0／1 | 对 |
| **F** 键留值丢 | **34／0／1、rc=0（全绿）** | **33／1／1** | 修前 34/0/1 → 修后 33/1/1 | 对 |
| **M-a** 新 map 槽、夹具不填 | **34／0／1、rc=0（全绿）** | **33／1／1** | 修前 34/0/1 → 修后 33/1/1 | 对 |
| **E** 空克隆 | **2／3／0、RUN=5、`panic: index out of range [0] with length 0`** | **31／3／1、RUN 54** | 修前红但吞 30 条 → 修后 31/3/1、RUN 54 | 对（连"吞读数"那一格都对上：RUN 从 5 回到 54） |
| **合体 M-a＋F** | **34／0／1、rc=0** | **32／2／1** | 修前与健康态一模一样 → 修后 32/2/1 | 对 |

**"修后不再一模一样"这句成立**：旧尺下 F／M-a／合体三发的四个数与 `=== RUN` 枚数**与健康态逐字相同**
（34/0/1、RUN 54、rc=0），新尺下三发各自 `rc=1` 且红句点名。⇒ **派单单要的那枚"本轮第一凭据"本程自己复算出来了。**

红句逐发点名（行号换算：F／M-a 无插行＝锚点行号）：

- F → `ticket146_liveapprovals_backing_test.go:285: AC#2 RED (值保真): … 不逐字段相等`，
  且**同一发里 probe 1 是 PASS**（本程单独 `grep` 过：`--- PASS: TestLiveApprovalsSharesNoReferenceSlotWithTheQueue`）。
  ⇒ 值保真这一味在 F 形下是**唯一的** catcher。
- M-a → `:222: 返回值里的引用槽位数 9 ([… Scratch])，期望 8 ([…])`。
- 合体 → `:222` ＋ `:285` 两行，各点一枚洞。

---

## 2. 打②：新尺的软肋那一形（**非空夹具**）——**没打穿，但"被哪一味堵住"必须改写**

派单单这一发的原话："新增一枚**不拷**的引用槽位、**且夹具把它填上非空值** ⇒ 现在会不会红？"
本程把它拆成三发（`cloneDecision` 始终**不拷**那枚第 9 枚 `Scratch map[string]string`）：

| 发 | 装什么 | 红在哪一味 | 本程读数 |
|---|---|---|---|
| `v-ma2a` | 夹具填非空、**手工名册 `queueSlotPaths` 未跟着长** | **类型普查**（锚点 `:222`，本发插一行 ⇒ 报 `:223`）census 9 vs 8 | 33／1／1，failed＝probe 1 |
| `v-ma2b` | 夹具填非空、**名册已跟着长**（9==9）、`cloneDecision` 仍不拷 | **probe 1 的 backing 指针比对**（`:255`＝锚点 `:253`）红句逐字点名 `共用底层：[Scratch]` | 33／1／1，failed＝probe 1 |
| `v-ma2c` | 名册已长、`cloneDecision` 给它分配了**一份全新的空 map**（不共用、但值丢光） | **只有 DeepEqual**（`:287`＝锚点 `:285`）；probe 1 这一发**是绿的** | 33／1／1，failed＝probe 2 |

⇒ **正面回答派单单那一问："第 9 枚会静静进来"这一格被堵住，靠的是两味、不是一味**：

1. **穷尽普查**只在"改了 `tools.Decision` 而没改手工名册"时承重（`v-ma2a`）。它的性质是
   **逼一次同步评审**（名册与 `cloneDecision` 必须一起长），不是"数出少拷了一枚"。
2. 真"少拷一枚"那一形（`v-ma2b`）被堵住靠的是 **probe 1 的底层指针比对**——而**这一味旧尺本来就有**：
   本程把同一发装到旧尺上（`pre-ma2b`），旧尺红在 `:198`（旧文件的同一处 `t.Errorf`），
   **不是**红在新普查上。⇒ **诚实的写法是**：非空夹具那一形**不是第二轮新买到的东西**；
   第二轮新买到的是"nil 夹具那一形（`v-ma`）与未知 kind 那一形（§3 的 `default`）"。
3. **DeepEqual 在非空夹具那一形下确实不承重**（`v-ma2b` 里两边是同一枚 map、`DeepEqual` 必然真），
   它承重的是 `v-ma2c`／F 那一族"分配了但没带上值"的形状。⇒ 修复件 §2.2 那句
   "（M-a 单独的'新槽位'由 §2.1 的类型普查收）"方向对，但它 §5.1／验收件那句
   "同一条 `reflect.DeepEqual` 断言一次收掉 E、F 与 M-a"**在本程的 `v-ma2b` 形上收不掉**——
   **那一格要按"三味分管三形"改写**，不是谁的洞，是文字归属。

⇒ **打②没有造出真洞**（三发全红），但**造出了一处必须改写的归属句**。记入账内、不构成本格退回理由。

---

## 3. 打③：`declaredRefSlots` 那五支逐枚各撤一发——**没有一支撤掉后变不响；但承重的不是那五支，是 `default`**

问法按派单单："撤掉这一支，哪条用例变得不响？"每一发都**同时**给 `tools.Decision` 装一枚该 kind 的字段
（否则那支今天根本不触发），跑逐包单跑：

| 撤掉的那支 | 装的字段 | 读数 | 红在哪 | 撤掉后响不响 |
|---|---|---|---|---|
| （不撤，只装 Ptr 字段＋名册＋夹具）`v-ptrkeep` | `Scratch *[]string` | 33／1／1 | probe 1 `:238: slot path "Scratch" ends on a ptr, not a map or slice` | — |
| Ptr `v-ptrdrop` | `Scratch *[]string` | 33／1／1 | probe 1 `:218: 引用槽位普查(Scratch): 未预期的 reflect.Kind ptr` ＝ **default 支** | **仍响** |
| Interface `v-ifacedrop` | `Scratch any` | 33／1／1 | 同上，kind 报 `interface` | **仍响** |
| Chan `v-chandrop` | `Scratch chan int` | 33／1／1 | 同上，`chan` | **仍响** |
| Func `v-funcdrop` | `Scratch func()` | 33／1／1 | 同上，`func` | **仍响** |
| Array `v-noarray` | `Scratch [2]map[string]string` | 33／1／1 | 同上，`array`（不撤时红在 `:222` 名册 9 vs 8，即 `v-arraykeep`） | **仍响** |
| **未知 kind（＝撤掉 `default`）** `v-nodefault` | `Scratch *[]string` ＋ 同时撤 Ptr 支 | **34／0／1、rc=0、全绿** | —— | **不响！洞是活的** |

`v-nodefault` 那一发本程另放了一枚**只进副本、不进仓**的探针（`zz_acceptor_liveness_test.go`）把"活"钉住：

```
DEEPEQUAL: green (aliased *[]string compares equal)
HOLE IS LIVE: writing through the returned *[]string rewrote the admitted record
              (queue now reads "CLOBBERED-BY-READER")
```

⇒ **本程的读法与结论**（这条不是退回，是把话说准）：

1. **"fail-closed 穷尽"这句按现量成立，不必降级成"实覆 N 支"。** 承重的那一支是
   `default: t.Fatalf` —— 撤掉它才出现"34/0/1 全绿而洞活着"；那五支（Ptr／Interface／Chan／Func／Array）
   各自撤掉时**都仍响**，只是红的位置从"名册数目不符"挪到"未预期的 kind"。
   按派单单"别要求新加那支先响"的交代，**这一格本程不判退回**。
2. **但那五支今天给的不是"测到"，是"拒绝测"。** `v-ptrkeep` 那一发（把该做的都做了：名册长、夹具填）
   红句是 `slot path "Scratch" ends on a ptr` —— 也就是说**一旦 `Decision` 真出现一枚 Ptr／Interface／
   Chan／Func／数组内嵌引用槽，这把尺只能停下来要求重设计，量不出共用与否**。
   这句是"fail-closed"的本义、不是缺陷，但**它使 `pending_read.go:36-38` 那句
   "a ninth reference-typed field cannot arrive quietly" 的确切射程＝"第九枚会被**点名**，
   其中 map／slice 两 kind 能被**测量**"**。⇒ 文字要按这个射程写，别写成"九种 kind 都能测"。
   （与第一轮同一族：许诺的射程不得大于仪器的射程。）
3. **夹具完备性守卫（`:245`）是真的有牙**，本程造了正控：把夹具里 `Paths` 那一行拿掉
   （`v-nilpaths`）⇒ 红在 `:244: fixture 未填引用槽 "Paths"（nil 或空）：backing 指针对它是失明的`；
   再把守卫本身撤掉、并撤掉 `Paths` 那枚克隆（`v-nilpaths-noguard`）⇒ **probe 1 变绿（瞎）**，
   只剩 probe 2 以 `panic: index out of range` 红。⇒ 守卫独立承重，**不是装饰**。
4. **一枚残留（不进退回，进账）**：`v-nilpaths` 那一发尽管守卫红，**同一次运行仍 `RUN=5`**
   —— 探针 2 无条件 `row.Decision.Paths[0] = …`，夹具一 nil 就 panic、**吞掉 49 条读数**。
   修复件 §3.4 那句"读数不再被吞"只对 **E 那一形**成立，**对"nil 夹具"那一形不成立**。
   按本仓规矩这仍是红（不是假绿、也不是拿 Skip 顶数），但**"RUN 恒 54"这句被写宽了**。

---

## 4. 打④：`reflect.DeepEqual` 会不会是枚恒真断言——**不是；本程自造两发"只有它红"的变异钉住**

派单单要查的四点，逐点现量（`tools.Decision` 全文读自锚点 `internal/tools/gate.go:15-53`）：

- **两侧是不是同一对象的两个别名？** 不是。左：`row.Decision` ＝ `LiveApproval` 结构体里那枚
  由 `cloneDecision` 造出的**新结构＋八枚新底层**；右：`mustFindLive(q,…).Dec` ＝ 队列 `q.pending`
  里那枚 `*qitem` 的字段（`pending_read_test.go:234`，白盒直读 `q.pending`）。两枚不同对象。
  **反向证据本程也量了**：把接线摘回浅拷贝（`c-bside`）⇒ 两侧真的同源 ⇒ DeepEqual **绿**，
  红的是 probe 1 的 `:253`（八枚全点名）＋ probe 2 的 11 行就地写断言。
  ⇒ 它是**值等值尺、不是共用尺**，与注释里那句"不共用底层却把值拷丢"的分工一致。
- **有没有时间戳／指针／不可比字段？** 19 枚字段：`Timeout` 是 `time.Duration`（int64 时长，**不是墙上时钟**）；
  **无 func、无 chan、无裸指针、无 interface 字段**；引用型只有 1 map ＋ 7 slice；
  `Blacklist` 是值结构 `BlacklistNote`（6 枚字段全 exported）。⇒ `DeepEqual` 全程走值语义、不 panic、
  不会因不可比类型退化成指针对比。**这一支不成立。**
- **"只有 DeepEqual 能抓、其它检都抓不到"那一发——本程自己造出来两发**：
  ① `v-f`（`out[k]=nil`，键留值丢）⇒ **33／1／1、唯一红＝`:285`，probe 1 同发 PASS**；
  ② `v-ma2c`（第 9 枚槽位给一份**全新的空 map**）⇒ **33／1／1、唯一红＝`:287`，probe 1 绿**。
  两发的红都是 `t.Fatalf` ⇒ 用例在那一行就终止，**后面的就地写断言根本没跑** ⇒
  "只有它红"是被结构钉住的，不是靠行号猜。**⇒ DeepEqual 独立承重，不是装饰。**
- **它对什么瞎？** 本程在**健康锚点**上放了一枚探针（`zz_acceptor_deeper_test.go`，不进仓）现量：

  ```
  DeepEqual before the write=true, after the write=true ; queue's Params[argv] now=[CLOBBERED-THROUGH-NESTED hi]
  ```

  ⇒ 通过返回值往 `Params` 的 value 里那层 `[]any` 写，**队列里那条已决断记录真的被改掉了**，
  而 DeepEqual **写前写后都是 true**。⇒ **修复件 §7 第 1 条那句自陈（"DeepEqual 对再往下一层的共用仍然看不见"）
  逐字准确，既没夸大也没缩小。**本程没有把它读成"深拷贝已买到"。**

---

## 5. 常规必查（放水两问／名册差集／门禁独立复跑／覆盖面自量）

### 5.1 那 37 行删除——逐行过，**没有任何一条会红的断言被删或被改弱**

`git show --numstat b694378` → **111 / 37**，只碰一枚文件（复算对上）。把删除行里**带判定形状**的
（`t.Error|t.Fatal|if |panic`）单筛出来，**只有 3 行**，全部落在被替换的旧尺函数体内部：

```
-		if !v.IsNil() {        (旧 refSlots：nil 槽不进名册 —— 正是第一轮 M-a 的根因)
-			if !f.IsExported() {  (旧 refSlots：私有字段静默跳过 —— 正是第一轮 M-c 的根因)
-			if prefix == "" {     (点号路径拼接，无判定作用)
```

⇒ **一条 `t.Errorf`／`t.Fatalf` 都没被删。**保留下来的三条名册判定（数目不符 `:221`、逐枚比对 `:227`、
共用底层 `:248`）极性一字未动，只是把 `slotPaths(值走)` 换成 `typeCensus(类型走)`。

### 5.2 放水问②"helper 是不是原有那枚"——**被换掉那把尺覆盖过的东西没有掉下来**（本程逐枚复量）

派单单那一问：第一轮量到"实覆 8 枚"，换成新尺之后**那 8 枚逐枚各撤一发还响不响**。本程不抄第一轮 §4，
在纯净母本 `snap0` 上重做九发（**每发只删 `cloneDecision` 里的一行**）：

| 撤掉的那枚克隆 | 读数 | probe 1 点名 | probe 2 里响的行（锚点行号） |
|---|---|---|---|
| 全部摘回 `it.Dec`（＝ⓑ 形状／修法整体不装） | 32／2／1 | `[Args Blacklist.Absolute Blacklist.AlreadyUnlocked Blacklist.Unlockable Capabilities Params Paths RulesHit]` 全 8 枚 | 307／311／314／317／321／324／333／346／366／369／372 ＝ **11 行** |
| `Params` | 32／2／1 | `[Params]` | 307／311／372 |
| `Args` | 32／2／1 | `[Args]` | 314 |
| `RulesHit` | 32／2／1 | `[RulesHit]` | 317／366 |
| `Paths` | 32／2／1 | `[Paths]` | 321／346 |
| `Capabilities` | 32／2／1 | `[Capabilities]` | 324／369 |
| `Blacklist.Absolute` | 32／2／1 | `[Blacklist.Absolute]` | 333 |
| `Blacklist.Unlockable` | 32／2／1 | `[Blacklist.Unlockable]` | 333 |
| `Blacklist.AlreadyUnlocked` | 32／2／1 | `[Blacklist.AlreadyUnlocked]` | 333 |

⇒ **"撤掉哪一枚，有没有一条用例变得不响？"＝没有一枚不响，且红句逐字点名被撤的那一枚**
（`共用底层：[Blacklist.Unlockable]` 这类是本程逐发 grep 出来的，不是从总数推的）。
⇒ **覆盖面主张照"8 枚"写即可，不必写成"实覆 N＜8"；换尺没造成覆盖回落。**
⇒ 同一批发在旧尺上也复算过一枚关键的：`pre-ma2b`（非空夹具＋名册已长＋少拷）旧尺红在 `:198`
⇒ **旧尺本来覆盖过的那一形，新尺仍覆盖，靠的是同一味 backing 比对。**

⚠ **一处本程自己的操作失误，原样登记**：`v-bside`／`v-drop-*` 那九发我第一次建母本时用了 `snap`，
而 `snap` 那时已经被本程塞进一枚探针件（`zz_acceptor_deeper_test.go`）⇒ 那九发读数是 `RUN=55／PASS=33`。
本程**没有去清那批目录**（临时件只建不删），而是换名从纯净 `snap0` 重建成 `c-bside`／`c-drop-*`，
上表是**重跑后的干净读数**（`RUN=54／PASS=32`）。旧读数文件 `out-v-bside.txt`、`out-v-drop-*.txt` 一并留在根目录。

### 5.3 名册差集——**49 枚两向 comm 皆空，复算成立**（并说明为什么本程另一把尺报 35）

```
grep -oE '\-\-\- (FAIL|PASS|SKIP): [A-Za-z0-9_/]+' <两向> | sed 's/^--- [A-Z]*: //' | sort -u
→ before=49  after=49 ；comm -23 空 ；comm -13 空
```

⇒ 修复件 §5 那句"49 枚同名、未增删用例函数"**逐字成立**（本程先按锚定行首 `^--- ` 量到 35，
那是不含子测试的顶层数；49＝含子测试名。**两把尺都自洽，别当成谁吹数**）。
SKIP 同一枚 `TestDefaultDeadlineWallClockMeasurement`，改前就在，未新增 Skip、未拿 Skip 顶读数。

### 5.4 门禁独立复跑（一律逐包单跑；**未**把 `./cmd/wisp/` 与 `./internal/panel/` 放进同一次调用）

| 门 | 命令（本程原样，全在纯净副本 `snapgates`） | 读数 |
|---|---|---|
| 包测 | `go test -count=1 -v ./internal/agent/approval/` | **rc=0；34／0／1、RUN 54**（锚点自身，未塞探针） |
| 格式 | `gofumpt -l . tools/d22scan tools/mockllm`（本机 `D:\work\base\gopath\bin\gofumpt.exe` **v0.12.0**＝CI 同版，`--version` 现量） | **空**，rc=0 |
| 静态 | `go vet ./internal/agent/approval/` | **rc=0** |
| D22 | `sh scripts/d22scan.sh` | **rc=0 clean**，各作用域 `examined N` 全非零：`internal/ 205`、`cmd/ 23`、`frontend/ 52`、`internal/tools/ 18`、`design/ 30`、ban#8 `internal/ 412`、`cmd/ 43` |
| **D22 正控**（"0 命中"先打一遍那把尺） | 往 `snapgates/internal/agent/approval/` 放一枚含 **U+2713** 的 Go 字符串，重跑 | **命中**：`zz_acceptor_ctrl.go:4: [emoji] ban #8 glyph in scope internal/ is banned (D23): non-comment text, string literals included` ⇒ 那把尺有牙，上面的 clean 是真读数 |
| 已知红 | 未复跑（不归本票）：`TestC21DesignTokensFourWayAgree` 两形红已登记 `A253④/A254/A257④`；`staticcheck` 48 条历史积压 | 引红连红句：本程**未**引它作任何判据 |

⛔ **计时读数一条未取**（多程在跑、`slo-full` 本机 self-hosted）。`cmd/wisp`／`internal/panel` 本程未跑
（别人的落点，且 AC#4 名单内）。

### 5.5 契约轴（AC#4 复算，只看本程被要求裁的两格的改动面）

```
for c in b694378 73e0d7be a1f5934c 74255b9; do git show --name-only --format= $c; done | sort -u
→ internal/agent/approval/ticket146_liveapprovals_backing_test.go
  docs/evidence/s1/146-liveapprovals-fix-r1.md
  docs/evidence/s1/146-liveapprovals-r1.md
  .scratch/wisp/issues/146-…-with-the-queue.md
```

过一遍 AC#4 那 11 条禁改路径 ⇒ **交集空**（`internal/risk/`、`rules_gateway`、`thresholds.go`、golden、
`allowlist.txt`、`docs/PLAN.md`、`docs/specs/`、`frontend/`、`design/`、`tools/d22scan`、`l2_grant_boundary`、
`internal/panel/`、`cmd/wisp/`、`scripts/slo-check` 全 0 命中）。
`pending_read.go`／`queue.go` 本程复算**确实一字节未动**（修复件 §2.1 那句"判为不需要动它"成立）。

---

