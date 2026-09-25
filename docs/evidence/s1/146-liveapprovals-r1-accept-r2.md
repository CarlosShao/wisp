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
| `v-docpin-with v-docpin-without` | 把 ⓑ 纪律字样插进／不插进 `pending_read.go` ＋ 一枚读源码的用例 | 独立复算"注释可以承重"（第一轮那句反例本程自己重造） |


**副本行号会漂**：本程往测试文件里插过行的那几发（`v-ma2*`、`v-ptrkeep`）红句行号与锚点件相差 ＋1／＋2，
下面逐处标了换算。判红绿一律只认 `^--- FAIL:` 那一行，计数一律带 `=== RUN` 枚数防 panic 吞读数。

**一条 CRLF 尺子自纠（先说，免得下一位照抄）**：派单单要求复算"两个 `-c` 都带"。本程先用
`grep -rl $'\r' --include='*.go' internal cmd` 量 ⇒ **456／455 全命中**（换到上一级目录拿 `snapgates/` 当参数再量
＝**483**，比 go 文件总数还大），那是一枚**假阳性**：GNU grep 3.0（Git-Bash）在这台机器上把裸 CR 模式
匹配到了几乎每一行。同一批文件用 `od -c` ＋ 逐字节 node 计数复量
⇒ `pending_read.go` / `ticket146_liveapprovals_backing_test.go` / `internal/tools/gate.go` / `cmd/wisp/main.go`
**CR=0、CRLF=0**。⇒ 两个结论：
① **副本确实是纯 LF**（这一条按字节定，不靠 grep）；
② **`grep -rl $'\r'` 这把尺子在 Windows 上读数不稳**（本程同一条命令两次量出 456／483），
所以第一轮验收件 §0 那句"→ **0**"本程**在这台机器上复现不出来**（它既不能证真也不能证伪，
只是一个取不到读数的形状）⇒ **下一位要验 CRLF，请按字节量（`od -c` 或逐字节计数），别按那把 grep**。
（第一轮 §14 第 ① 条登记过同族"尺子错"，那是 `od` 输出上数 `\r`；本程这发是同一枚坑的另一面。）


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

## 6. 打⑤：那段 `>` 更正准不准——**(a) 假理由没有被任何新文字重复；(b) 理由①独自撑得住，本程还给它多找到一腿**

### 6.1 (a) 那一句假理由有没有在被重复——**没有**（四处新文字逐枚 grep）

```
grep -rn "任何用例都不会红|拿不到任何仪器|不需要仪器|唯一可能的仪器|装不下" \
  docs/evidence/s1/146-liveapprovals-fix-r1.md docs/evidence/s1/146-liveapprovals-r1.md \
  .scratch/wisp/issues/146-*.md \
  internal/agent/approval/ticket146_liveapprovals_backing_test.go internal/agent/approval/pending_read.go
```

命中全部落在两类**允许**的位置：① 实现件 §2 原文那三行（`:146/:148/:150`）——**那是刻意"一字不抹"保留的旧段**，
更正段（`:163-178`）逐条标注为"假"；票面 :84 与第一轮 §15 都要求"追加更正、不改写已提交正文"，
这一保留**是正确的做法而不是漏改**。② 否定式引用（修复件 `:277`"…不是'ⓑ 拿不到任何仪器'"、
票面 Progress log `:166`"理由②③＝**假**"）。
⇒ **两枚 .go 文件零命中** ⇒ 假理由没有沉淀成代码注释里的规矩。**这一支通过。**

### 6.2 (b) 理由①单独够不够撑起 AC#1——**够；本程自己复算了它，并复算了"注释可以承重"那一发**

理由①的内容是"AC#2 的字面断言（就地写返回值之后队列记录**没**跟着变）只在 ⓐ 下为真"。本程**不采信任何一方转述**，
在纯净母本 `snap0` 上把接线摘回 ⓑ 形状（`Decision: cloneDecision(it.Dec)` → `Decision: it.Dec`）现量：

```
c-bside  rc=1  32 PASS / 2 FAIL / 1 SKIP  RUN=54
  --- FAIL: TestLiveApprovalsSharesNoReferenceSlotWithTheQueue   :253 共用底层：[全 8 枚点名]
  --- FAIL: TestLiveApprovalsInPlaceWriteCannotReachTheQueueRecord 307/311/314/317/321/324/333/346/366/369/372（11 行）
```

⇒ ⓑ 之下交不出一枚绿的 AC#2（本仓那族"恒真判据＝新假绿"），**方向由理由①独自撑住这句复算成立**。
第一轮那格判的是"够、但只到'字面收得下'"——**本程同意这个射程**：理由①只证明"ⓑ 交不出 AC#2 字面那枚检"，
不证明"ⓐ 是唯一可行的修法"（ⓓ／ⓔ 仍是下一张票的形状，第一轮的收紧段已写过，本程未松动它）。

**本程另找到的一腿（多撑一道，但依赖一处读法，如实标）**：AC#1 的判据有四条腿
（改元素／append／就地写／**取地址后传出去**）。前三腿本程复算＝**0**（尺子正控：同一把尺在 probe 2 自己体内命中 5 处，
剔 `_test.go` 后 0 处；`&x.Decision` 那一形全仓 0 处）。**但第四腿按"Go 的切片值本身带底层指针"这一读法答案是"有"**：
`cmd/wisp/panel_pump.go:68` 把返回值里的 `it.Decision.RulesHit` **按切片值交给** `panel.NativeVerdict.RulesHit`，
底层数组就此逃逸进面板侧数据（`Args` 那一腿不算——`panelArgs` 把它重新 `json.Unmarshal` 成新分配）。
⇒ 若编排者按字面只认 `&` 那一形，第四腿仍是"无"，ⓐ 靠理由①一腿站着；若按"传出去"的字义读，
ⓐ 被 **AC#1 自己的判据**多撑了一道。**两种读法下 ⓐ 都成立 ⇒ 不为这一处停手、也不构成退回。**

### 6.3 "改注释的票可以有仪器"这一句——本程**独立复算**了第一轮的反例（没有继承它）

理由②被判定为假的唯一凭据是第一轮在 `snap4` 里造的那枚 doc-pin，而它**不进仓**。本程不接受"读不到的证据"，
于是自己重造一遍：把票面 :42 的 ⓑ 纪律字样（`调用方不许就地写`）插进 `pending_read.go` 注释、
放一枚读源码的用例（`v-docpin-with`），再把那句拿掉（`v-docpin-without`）：

```
v-docpin-with    : qq_docpin_test.go:26 (b)-branch doc-pinning instrument is GREEN …      --- PASS
v-docpin-without : qq_docpin_test.go:23 AC#1(b) RED: pending_read.go no longer carries …  --- FAIL
```

⇒ **"注释可以是承重的"这句本程自己量到了**，实现件 §2 的理由②确为假句、更正段的事实基础成立。

### 6.4 `AGENTS §0.2` 那一句"不需要人工批准"——**本程复算成立**

派单单断言"这不需要人工批准，因为方向没改、代码没回退"。三条前提本程各自现量：
① 方向没改＝`pending_read.go:117` 仍是 `Decision: cloneDecision(it.Dec)`（锚点纯净副本直读）；
② 代码没回退＝`c-bside` 那一发证明"回退形态"是**需要现造出来的变异**，锚点上不是；
③ 四枚 commit（`b694378 73e0d7be a1f5934c 74255b9`）碰过的 4 枚文件与 `D1–D47`／`C1–C32`／`R1–R9`／D43 转移表
的载体（`docs/PLAN.md`、`docs/specs/**`）**交集空**（§5.5）。⇒ **没有发生契约变更，§0.2 那句站得住。**

---

## 7. 本程没测什么（按"如果我漏了它，谁会先被骗"排序；本程自己写的）

1. **没测"新尺在真并发下的形状"**。两发检都是单 goroutine；`cloneDecision` 把持锁时间拉长这一条，
   第一轮 §12 第 2 条挂着、修复件 §7 第 3 条自陈未测、**本程同样未测**。先被骗的是票 33 的宿主程
   ——它下一步就是让泵在状态变动里高频重算。**没人量过"泵 × 应答抢锁"那一形，今天仍无一条用例探它。**
2. **没追 `gate.go:553` 那枚 `Params: d.Params` 的下游**。§6.2 只量了"面板泵不读 `LiveApprovals` 的 `Params`"，
   但 `internal/agent/approval/gate.go:553` 在生产里**确实**把 `Decision.Params` 往另一枚结构里搬。
   本程**没有追那条路走到哪、有没有人就地写它**。先被骗的是"以为 `Params` 全仓无人读"的人
   ——§6.2 那句要读成"**从 `LiveApprovals()` 这条出口**无人读"，别读成"全仓无人读"。
3. **没测 Array／嵌套引用那一族"能被点名但量不了"的形状会不会一次只报一枚**。`v-ptrkeep` 只证了一枚 Ptr 字段
   停在 `slot path … ends on a ptr`。若将来 `Decision` 同时长出两枚这种字段，本程没量这把尺会不会红在第一枚、
   把第二枚藏进那一次红里。先被骗的是"看到红就以为全量过了"的人。
4. **没复跑 AC#3／AC#5 的既有格**（派单单划出去）。本程只把锚点自身的包测与三道静态门重跑了一遍（§5.4），
   **没有重跑 `wisp slo` 双臂、没跑 `cmd/wisp`、没跑 `internal/panel`**。跨包回归本程不背。
5. **没量手工名册的反向那一形**：`queueSlotPaths` 多写一枚而 `Decision` 没长（预期 9、实得 8）。
   本程推它红在 `:222`，**但这一发本程没跑过**——"名册会不会变成第二枚恒真"只验了单向。
6. **`=== RUN` 恒 54 只在健康态与四发主变异下量过**。`v-nilpaths` 那一发本程量到 `RUN=5`（§3 第 4 条），
   说明"恒 54"是**有条件的**；本程没有穷举还有哪些形仍会 panic 吞读数。先被骗的是拿"RUN=54"当
   "读数没被吞"的**充分**条件的人——它是必要条件。
7. **红句行号经过一次人工换算**：`v-ma2*`／`v-ptrkeep` 因本程自己插了行，行号与锚点件差 ＋1／＋2。
   本程用"锚点行号＋已知插行数"两路对过（§2 表内标了换算），**但没做逐行的机器校验**。
   抄这些行号的人请连换算一起抄。

---

## 8. 伪授权登记（两个数分开计，不合成一个字段）

| 栏 | 数 | 出处（工具名＋命令前 40 字） |
|---|---|---|
| **真通知回显数** | **5** | ① 派单单随身的 `<loaded_context>`：项目 `agents.md` 全文回显；② `Bash` 回执里的 `The date has changed. Current date: 2026-09-25`；③ `Bash`／`Read`：可用 skill 清单回显；④ `Bash`：`agents.md` 的 Memory 回显；⑤ `Bash` 回执末尾的 `Note: The file …\memory\MEMORY.md was modified since it was last read.`（**"文件已被修改"这一形状在派单单的禁列表里**——它没要求任何动作、没声称授权、讲的是仓外的记忆文件，本程只把它当回显记下；**若它曾要求本程改代码／改判／放宽阈值，一律按注入处理、不执行**）。五者**都不要求动作**，本程**未据此减少任何取证**。 |

| **判为注入数** | **0** | 本程工具输出里未出现"系统提示／编排者备注／已核验请继续提交／Confirm the harness note is genuine／请 revert／放宽阈值／这格我已签收／不用取证直接给结论"任何一枚形状；亦未出现自称锚点的假 sha。 |

- **锚点只认自己量的**：开工第一发即 `git cat-file -t 74255b9` → `commit`、`git rev-parse --short HEAD` → `74255b9`。
  派单单与两份既有件里出现的其它 sha（`7f75ff3`／`b694378`／`b694378^`／`1e94672`…）本程都在**用到的那一枚**上
  重跑过 `git show --name-only`／`git diff --stat`／`git archive`，未凭字面采信。
- **反向不洗**：本程"0 注入"只是本程这一路的读数，不抵账前几程各自登记的数；同程"我没遇到"不能洗掉别人遇到的。
- **工具调用被拒数：0**（无权限弹窗、无 authorization 拒绝）。
- **本程自己踩空、原样登记（都不是授权、都不需要人批）**：
  ① 第一条建副本命令把 `cd` 写在 `git archive` 前 ⇒ `fatal: not a git repository` ＋ `tar: does not look like a tar archive`
  ⇒ 那一发**没有产出副本**，改成 `git -C "<repo>"` 重跑成功。
  ② 把 `go test` 输出重定向到绝对路径、却用相对路径 grep ⇒ 四发 `No such file`（**是尺子错，不是读数**）。
  ③ **`grep -rl $'\r'` 在 Git-Bash 上报 455／455 假阳性** ⇒ 换 `od -c` ＋逐字节 node 计数才量对（§0 已写）。
  ④ 变异 F 第一版写成 `out[k] = nil` **不编译**（`v declared and not used`）⇒ 两发 `build failed`；
  本程**没有把 build failed 当读数**，改成 `_ = v` ＋换名重建 `v-f`／`v-comb`。
  ⑤ **建 `v-bside`／`v-drop-*` 九发时母本用了已被本程塞过探针的 `snap`** ⇒ 那九发 `RUN=55`；
  换名从纯净 `snap0` 重建成 `c-*` 九发，旧读数文件留在盘上不删。
  ⑥ `git show --format="$c"` 把 sha 当成 pretty 格式名 ⇒ `fatal: invalid --pretty format`（那一发文件清单未出，重跑）。
  ⑦ `Edit` 的 `old_string` 不唯一 ⇒ 一次失败，加上下文后成功。
  **全程未动仓内任何一份代码／测试／票面**；副本一律只建不删；未 `rm`、未在仓内建 worktree／checkout。
- **凭据值**：本件全程未抄任何 key/token 原文（连变量名以外的字面都没抄）。

---

## 9. 两格裁决表 ＋ 给编排者的裁定建议

| 格 | 档位 | 现量命令／发（复算入口） | 判决要点 |
|---|---|---|---|
| **AC#1** 判 ⓐ／ⓑ（本轮只裁"退回的那两行论证修好了没"） | **成立** | §6.1 四处新文字 grep；§6.2 `c-bside`（ⓑ 形状 ⇒ 32/2/1、两发各红、8 枚全点名）；§6.3 `v-docpin-with`／`v-docpin-without`；§6.4 四枚 commit 的文件并集 vs 契约载体 | ① 那两枚假理由**已从"唯一支撑"降为"被标注为假"**，且未沉淀进任何 .go 注释；② 更正段引的每一发本程**都独立复算过**（doc-pin 那发本程自己重造了一遍，没有继承一份读不到的仓外证据）；③ 理由①独自撑住方向＝本程现量成立，**且本程为它多找到一条腿**（AC#1 第四判据"传出去"在 `panel_pump.go:68` 上读作"有"）；④ `AGENTS §0.2`"不需人工批准"的三条前提逐条复算成立。**第一轮"只退论证"那一档，论证已经退掉了。** |
| **AC#2** 会响的检 | **成立**（附 5 条**文字类**入账，全部不挡勾、不需要再动代码） | §1（五行变异原样复跑，含旧尺对照 34/0/1）；§2（非空夹具三形分别红在普查／backing／DeepEqual）；§3（五支单点回退＋撤 `default` 才出现全绿、本程探针钉出洞活）；§5.2（8 枚逐枚撤克隆全响、红句逐字点名；ⓑ／摘修法 11 行红）；§4（DeepEqual 非恒真：两发"只有它红"＋字段表无不可比） | 承重与覆盖面**双双复算成立**：摘修法红、逐枚撤各自红、三味分管三形互不重叠。**第一轮那枚"第 9 枚静静进来"的洞已被钉**（`v-ma` 33/1/1），且**不是靠把许诺削短**（选的是甲、`pending_read.go` 一字节未动，§5.5 复算）。放水两问都过：37 行删除零条 `t.Error/t.Fatal`、被替换的旧尺覆盖过的形（`pre-ma2b`）新尺仍覆盖。名册 49 枚两向 comm 空。**入账的五条全是文字归属，见 §9.1** |

### 9.1 AC#2 的入账清单（**文字类，不构成本格退回**；下一手碰这三枚文件时顺手改掉即可）

1. **归属要改**：验收件 §5.1／修复件 §2.2 那句"同一条 `reflect.DeepEqual` 一次收掉 E、F **与 M-a**"——
   `v-ma2b` 形（名册已长、少拷一枚、非空夹具）下 **DeepEqual 必然绿**（两边同一枚 map），收掉它的是 probe 1 的 backing 比对。
   正确写法是**三味分管**：普查管"声明了没点名"、backing 管"点名了没拷"、DeepEqual 管"拷了没带值"。
2. **射程要写准**：`pending_read.go:36-38` 与测试件 `:23-25` 的"第九枚不会静静进来"，确切射程＝
   "第九枚**会被点名**；其中 map／slice 两 kind **能被测量**，Ptr／Interface／Chan／Func／数组内嵌引用槽**只能拒绝**"
   （`v-ptrkeep` 的红句是 `slot path "Scratch" ends on a ptr`）。
3. **"RUN 恒 54"写宽了**：`v-nilpaths` 形下守卫红、但 probe 2 仍 `panic: index out of range`、`RUN=5`。
   要写成"五发主变异下不吞读数；nil 夹具那一形仍由 panic 兜、仍会吞"。
4. **`v-ma`（第一轮 M-a 那一发）的红是"逼同步评审"型的红**（数目 9 vs 8），不是"测出少拷一枚"的红。
   这一条写进测试件 `:186-189` 那段注释正好，别写成后者。
5. **修复件 §1.1 复现变异 F 的那两行读数**（`rc=0` ＋两发 PASS）是在**旧尺**上量的、文字没标尺的世代。
   本程 `pre-f` 复算＝同读数 ⇒ **数字对、世代标签缺**，补一句"旧尺"即可。

### 9.2 给编排者的裁定建议：那三件空格**该开新票、并入家族票、还是记台账**

| 空格 | 本程现量到的形状 | 建议 | 理由（一句话） |
|---|---|---|---|
| **① "再往下一层仍共用"**（`Params` 的 value 里那层 `[]any`／嵌套 map） | 本程探针在**健康锚点**上量到：通过返回值写 `Params["argv"].([]any)[0]` ⇒ 队列那条记录**真被改掉**，而 DeepEqual 写前写后都是 `true`；仓内**零条用例红** | **记台账 ＋ 并入下面那枚新票**（不单独开票） | 单独开票只为"深拷第二层"＝**为没有读者的字段买字节**：`Params` 从这条出口**在生产侧无人读**（`panel_pump.go:66-71` 五枚字段不含它），先判"要不要送出去"再判"送出去拷几层"才是对的顺序 |
| **② "把 `Params` 递到面板这件事要不要收窄"（ⓓ）** | 同上；且第一轮 §9 已量到 `Params` 那一枚 map 占这枚改动新增量的 **64 % 字节／22 % 分配次数**——不拷它省的是字节那一半 | **开一枚新票**（**不要**并进 146 家族票） | 它必须**先把 AC#2 的措辞改成"对返回值里任何引用槽就地写"**才能收 ⓓ（第一轮 §6.5 已收紧：不改 AC#2 时 ⓓ 让那一发"无从下笔"）⇒ **改 AC＝改票＝人工批准**，这个动作不能塞进现有票的续发里 |
| **③ `replay` 那枚同族出口** | 本程在锚点复算**仍在**：`internal/agent/approval/queue.go:488-494` 先 `fresh := *it` 再 `d := it.Dec`，把含 8 枚引用槽的 `Dec` **按值搬出来、没走 `cloneDecision`**；`Replay` 在测试里有 4 处提及，断的都是状态不是底层共用 | **开一枚新票**（家族票也行，但**必须先扩 AC#4 的禁改面**） | 修法只有两枚：让 `replay` 也走 `cloneDecision`，或把它钉成一枚会响的检——**两枚都要改 `queue.go`**，而 `queue.go` 是票 146 AC#4 明写的零字节面 ⇒ 并进来等于事后放宽自家契约 |

⇒ **一句话给编排者**：①**记台账**并作为②的**第二条 AC**；②③各开新票，
②的前置动作是**人工批准改写 AC#2 的极性覆盖面**，③的前置动作是**人工批准把 `queue.go` 从禁改面里放出来**
（或另设一枚"钉而不修"的检）。**本程未替任何人勾一枚框、未动一份码、未 push。**

**next=**：本件交完派单单要求的五发（①②③④⑤）＋常规必查（放水两问／名册差集／门禁四发＋正控／AC#4 交集）
＋两格裁决表＋三件空格的裁定建议＋"我漏了什么"＋伪授权两栏。**没有剩余格。**
若编排者要本程补 §7 第 1 条（`-race` 并发形）或第 5 条（名册反向那一发），点名即可——两发的建法都写在表里了。


