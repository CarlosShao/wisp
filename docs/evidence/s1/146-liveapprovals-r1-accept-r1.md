# 146-r1 对抗验收（r1）—— `Queue.LiveApprovals()` 的 "copies"：裁决表

- 被验对象：票 146（`internal/agent/approval/pending_read.go` 的 `cloneDecision` ＋ 新测试 ＋ 证据件 ＋ 票面 Progress log）
- 被验锚点：`7f75ff3` ＝ 本程 `git rev-parse 7f75ff3` 现量 `7f75ff3c9aed48054656b90661af274b09906b30`；
  `git cat-file -t 7f75ff3` ＝ `commit`（锚点只认自己量的这一枚）
- 本程身份：**验收程（非实现者）**。派单要求逐格出表、四发进攻齐、每格裁完即 commit
- 副本（只建不删，全部在仓外）：`D:\tmp\wisp-146acc1\snap` ＝ `git -c core.autocrlf=false -c core.eol=lf archive 7f75ff3 | tar -x`
  ＋另外三棵变体副本（下面逐处点名）。**取副本时两个 `-c` 都带了**，且复算过没被烙上 CRLF：
  `head -c 60 snap/internal/agent/approval/pending_read.go | od -c` → 只有 `\n`；
  `grep -rl $'\r' --include='*.go' snap/ | wc -l` → **0**（这条正控是必要的：派单单说这是踩过的坑）

## 0. 本程用了哪几棵副本（复算入口）

| 副本 | 内容 | 用来量什么 |
|---|---|---|
| `snap` | `7f75ff3` 纯净 ＋ 摘掉修法的变异（已还原，`diff` 证 RESTORE-CLEAN） ＋ 本程探针件 | 进攻①、`sh scripts/d22scan.sh`、`go build ./cmd/wisp` |
| `snap-before` | `630c218` 纯净 | AC#1 普查的**改前**复算（后来又往里塞了本程探针件 ⇒ 门禁读数不以它为准，见 `snap7`） |
| `snap2` | `7f75ff3` ＋ 本程的**按类型**字段普查件 ＋ 各枚变异（**全部已还原**，`diff -r snap6/internal/tools snap2/internal/tools` 逐字节一致） | 字段普查、进攻④的 M-a／M-c、变异 E／F／合体 |
| `snap3` | `7f75ff3`，**代码一字节未改** ＋ 深度探针件 | 进攻④的"一层深"复算 |
| `snap4` | `7f75ff3` ＝ 摘回 `Decision: it.Dec` ＋ 本程按票面 :42 字样插入的 ⓑ 纪律注释句 ＋ 本程自造的 ⓑ 支两枚仪器 | AC#1 反证（ⓑ 到底收不收得下 AC#2） |
| `snap6` | `7f75ff3`，**未放任何本程探针件** | AC#5 的"改后"门禁与名册、`gofumpt -l . tools/d22scan tools/mockllm`、`gofmt -l .`、`go vet` 三包 |
| `snap7` | `630c218`，**未放任何本程探针件** | AC#5 的"改前"门禁与名册 |
| `snap5` | 空目录（建早了，没用上） | — |
| `orig` `mut` | 原件副本 ＋ 每一发的读数日志（文件名逐处点名） | 还原与复核 |

本程放进副本的探针件（`zz_acceptor_*_test.go`／`qq_*_test.go`）**不进仓**：裁决者不实现（D22）。

---

## 1. 字段普查独立复算（**8 枚成立；未发现第 9 枚；票面 5 枚是漏计**）

派单单要求"自己再数一遍 `tools.Decision` 的全部引用槽位（含嵌套一层）"。本程**不采信任何一方的表**，
用一枚按 **reflect 类型**（不是按值）走全字段的普查件量（`snap2`）：

```
cd D:/tmp/wisp-146acc1/snap2
go test -count=1 -v -run TestZZAcceptorCensusTypeShape ./internal/agent/approval/
```

读数（关键行）：

```
TOP total fields=19 kinds=map[bool:2 int:2 int64:1 map:1 slice:4 string:8 struct:1]
TOP 2  Params          map[string]interface {}   kind=map
TOP 3  Args            jsontext.Value            kind=slice   (elem=uint8)
TOP 5  RulesHit        []risk.RuleID             kind=slice   (elem kind=string)
TOP 7  Paths           []string                  kind=slice
TOP 8  Capabilities    []tools.Capability        kind=slice   (elem kind=string)
TOP 13 Blacklist       tools.BlacklistNote       kind=struct
NESTED Blacklist (tools.BlacklistNote) fields=6
   Class=risk.Class(int) Absolute=[]string Unlockable=[]string AlreadyUnlocked=[]string Reason=string Unresolved=int
CONTROL json.RawMessage kind=slice elem=uint8
```

⇒ **引用槽位＝ 1 枚 map ＋ 7 枚 slice ＝ 8 枚**，逐枚是
`Params`／`Args`／`RulesHit`／`Paths`／`Capabilities`／`Blacklist.Absolute`／`Blacklist.Unlockable`／`Blacklist.AlreadyUnlocked`。

三处派单单点名要复查的疑点，逐条答：

1. **`json.RawMessage` 本质是不是 `[]byte`** ⇒ 是（正控行：`kind=slice elem=uint8`）。
   注意 `f.Type.String()` 在 go1.27 里把它打成 **`jsontext.Value`**（标准库换了本体）——
   谁按字符串认类型会数漏这一枚，本程按 `Kind` 认，故不受影响。
2. **`Capability` 是不是真的没有引用字段** ⇒ 没有：`type Capability string`（`internal/tools/capability.go:15`），
   普查件读到的 elem kind=string。**所以"元素本身还有引用字段"这一支在本结构里不存在**，
   实现件 §1.2 那句"头共用、元素不共用"与票面 :24 的"本票没量"这一格可以**结掉**。
3. **有没有藏在私有字段里的嵌套结构** ⇒ `Decision` 19 枚字段 **全 exported**、`BlacklistNote` 6 枚 **全 exported**，
   没有 `refSlots` 的 `IsExported` 那一道会漏掉的东西。

⇒ **判定：实现程把票面的 5 枚改成 8 枚，改对了**；AC#2 那两发检的覆盖面主张不必重写（覆盖面本身见 §4）。

（另记一条与普查同时抓到的尺子性质：`reflect.ValueOf([]byte(nil)).Pointer()` 与
`reflect.ValueOf([]byte{}).Pointer()` **都是 0** ⇒ probe 1 的 `pa != 0 && pa == pb` 判据对
"非 nil 但长度为 0"的切片**结构性失明**。今天 8 枚槽位在 fixture 里都 ≥1 长度，所以不咬人；记进 §7。）

---

## 2. 进攻④：深度许诺复算 —— **"一层深"这半句逐字为真；"第 9 枚不会静静进来"这半句被本程造反例推翻**（退回）

### 2.1 注释原文（被验锚点 `internal/agent/approval/pending_read.go`）

```
:32-35  The copy is ONE level deep: what a caller cannot now do is write a slot the queue
        reports; what it can still do ... is reach THROUGH Params' values into the containers
        a JSON decode left inside them (a []any or nested map[string]any is still shared).
        Treat those values as read-only too.
:36-38  ticket146_liveapprovals_backing_test.go is the instrument for the eight named slots,
        and it walks tools.Decision by reflection so a ninth reference-typed field cannot arrive quietly.
```

### 2.2 "一层深"＝**量出来为真**（代码一字节未改的 `snap3`）

```
cd D:/tmp/wisp-146acc1/snap3
go test -count=1 -v -run TestZZAcceptorDepthPromiseMatchesCode ./internal/agent/approval/
```

读数（`D:\tmp\wisp-146acc1\mut\m-d-depth.txt`）：

```
:40 A-CONFIRMED: queue's Params[argv][0] was rewritten through the returned row
    ([]interface {}{"CLOBBERED", "hi"}) -> the comment's '[]any is still shared' is TRUE of the code.
:47 A-CONFIRMED: queue's Params[nested][depth] rewritten too -> 'nested map[string]any is still shared' is TRUE.
:65 B: append on the returned Paths slice changed queue len from 1 to 1 (must be equal)
```

⇒ 注释**没有夸大也没有缩小**：它说仍共享的，确实仍共享（两种形状各中一发）；它没保证的 `append`，确实进不去队列。
`cloneParamsMap` 的函数体就是 `out := make(...); for k, v := range in { out[k] = v }`，
与"只拷一层"逐字对齐。**这一半不退回。**

### 2.3 "第 9 枚引用字段不会静静进来"＝**本程造出了让它静静进来的两发**

**M-a：给 `tools.Decision` 加一枚 `cloneDecision` 不拷的 map 槽位，fixture 不填它 ⇒ 两发检全绿，而洞是活的。**

```
# snap2/internal/tools/gate.go 的 Decision 里加：Scratch map[string]string   （变异件留在 mut\gate.go.M-a）
cd D:/tmp/wisp-146acc1/snap2
go test -count=1 -v -run 'TestLiveApprovals(SharesNoReferenceSlot|InPlaceWriteCannotReachTheQueueRecord)' ./internal/agent/approval/
→ rc=0
  --- PASS: TestLiveApprovalsSharesNoReferenceSlotWithTheQueue (0.00s)
  --- PASS: TestLiveApprovalsInPlaceWriteCannotReachTheQueueRecord (0.00s)    # mut\m-a-ninth-map-no-fixture.txt
# 本程另放一枚探针，走 AC#2 禁的那个形状，只是打在 Scratch 上：
go test -count=1 -v -run TestZZAcceptorNinthSlotAliasesAndProbesMissIt ./internal/agent/approval/
→ --- FAIL  :34 AC#2-SHAPE HOLE IS LIVE: writing through Decision.Scratch rewrote the admitted record
      (queue now reads grant=CLOBBERED-BY-READER).                            # mut\m-a-live-hole.txt
```

**根因**：`refSlots`（`ticket146_liveapprovals_backing_test.go:81-100`）走的是**值**——
`case reflect.Map, reflect.Slice: if !v.IsNil()` ⇒ fixture 没填的槽位一条都不出 ⇒
名册断言仍是 8 vs 8 ⇒ 恒绿。它没有走**类型**，所以"枚数变了"这件事只能靠人记得往 fixture 里填。

**M-c：`refSlots` 的 `switch v.Kind()` 只有 `Map`／`Slice`／`Struct` 三支**（外加 `:90` 的 `IsExported` 才往下走）
⇒ **Ptr／Interface／Chan／Func／私有字段这五类一律静默跳过**，就算将来有人把它们填进 fixture 也数不到：

```
refSlots on {Ptr *[]string, Map map[string]string} = [Map]      # mut\m-a-live-hole.txt
```

⇒ 同一棵树上测试文件自己的两处同源表述也一起过期：`:18-22`"a future reference field on tools.Decision that
cloneDecision forgets is named by this case rather than discovered by an exploit"、
`:140-142`"neither a field added to Decision nor a walker that quietly stopped descending can turn probe 1 into an
always-green check"。**M-a 就是一枚"加到 Decision 上的字段"，它没有让 probe 1 红。**

**这一格的方向是收紧**，所以本程**不接受**"把注释改弱"作为唯一修法。最小闭合集合（任一即可，本程不代做）：

1. **（推荐）** 把 `refSlots` 的名册改成**按 `reflect.TypeOf` 数声明**（含 Ptr／Interface／Chan／Func／非 exported），
   名册不符即红；`cloneDecision` 与名册不一致当场红 ⇒ 注释那句从"许诺"变成"仪器"；或
2. 只把三处文字（`pending_read.go:36-38`、测试 `:18-22`、`:140-142`）收窄到真做到的那一格
   （"只认 map／slice，且只认 fixture 填过的槽位"）⇒ 文字诚实了，**但洞仍在**，须按 §7 第 1 条继续挂着。

### 2.4 顺带核掉的四处跨文件引用（这族缺陷的本名是"注释预先引用尚未产出的读数"）

| 注释句 | 现量 | 判定 |
|---|---|---|
| `:99` bound＝`DefaultMaxPending`（`queue.go`） | `internal/agent/approval/queue.go:112  DefaultMaxPending = 8` | 对得上 |
| `:99-101` 泵不随 streamed text deltas 动（`cmd/wisp/panel_pump.go`） | 同文件 `:236  // Not on streamed text deltas.` | 对得上 |
| `:101-103` 读数在 `docs/evidence/s1/146-liveapprovals-r1.md` §4 | §4 存在、A／B 两张表都有数 | 文件与节号对得上（读数真伪见 §6 AC#3） |
| `:98` "one map plus up to seven fresh backings" | `cloneDecision` 8 条赋值；`cloneBacking`／`cloneParamsMap` 对 nil 早返 | "up to" 用词准确 |

---

## 3. 进攻①：复跑改前 —— **会红，且实现件 §3.2 那 14 条明细逐字复算成立**

在 `snap`（`7f75ff3` 纯净副本）里把接线摘回票面说的那个形状（`cloneDecision` 留着不响）：

```
# internal/agent/approval/pending_read.go:117
#   -			Decision:      cloneDecision(it.Dec),
#   +			Decision:      it.Dec,
cd D:/tmp/wisp-146acc1/snap
go test -count=1 -v -run 'TestLiveApprovals(SharesNoReferenceSlot|InPlaceWriteCannotReachTheQueueRecord)' ./internal/agent/approval/
```

读数（`D:\tmp\wisp-146acc1\mut\m1-revert-all.txt`）：**`rc=1`／`^--- FAIL`＝2 行／`AC#2 RED`＝14 条**，
红句点名的槽位是**全 8 枚**：

```
--- FAIL: TestLiveApprovalsSharesNoReferenceSlotWithTheQueue (0.00s)
--- FAIL: TestLiveApprovalsInPlaceWriteCannotReachTheQueueRecord (0.00s)
:196 共用底层：[Args Blacklist.Absolute Blacklist.AlreadyUnlocked Blacklist.Unlockable Capabilities Params Paths RulesHit]
:233 队列存的 Params 里出现了调用方塞进的键 "injected-by-caller"（键 argv 现值 clobbered）
:237 :240 :243 :247 :250 :259×3 :272 :292 :295 :298
```

⇒ **检不是装饰**；"摘掉修法有没有任何外部可见读数变过"这一问的答案是"有"（rc、`--- FAIL` 行数、红明细数三处都变）。
恒真判据那一支不成立。

## 4. 进攻②：单点回退 —— **实覆 8 枚（每一枚各自会红，且红句点名被撤的那一枚）**

`sweep.sh`：每一发**只删 `cloneDecision` 里的一行**（＝只撤一枚槽位），跑同一对用例。
读数逐发留在 `D:\tmp\wisp-146acc1\mut\one-{59..66}.txt`；跑完 `cp` 还原并 `diff` 验过 RESTORE-CLEAN。

| 撤掉的槽位（行号） | rc | `--- FAIL` 行数 | `AC#2 RED` 条数 | probe 1 点名 | probe 2 里响的是哪几行 |
|---|---|---|---|---|---|
| `Params`（:59） | 1 | 2 | 4 | `[Params]` | 233／237／298 |
| `Args`（:60） | 1 | 2 | 2 | `[Args]` | 240 |
| `RulesHit`（:61） | 1 | 2 | 3 | `[RulesHit]` | 243／292 |
| `Paths`（:62） | 1 | 2 | 3 | `[Paths]` | 247／272 |
| `Capabilities`（:63） | 1 | 2 | 3 | `[Capabilities]` | 250／295 |
| `Blacklist.Absolute`（:64） | 1 | 2 | 2 | `[Blacklist.Absolute]` | 259 |
| `Blacklist.Unlockable`（:65） | 1 | 2 | 2 | `[Blacklist.Unlockable]` | 259 |
| `Blacklist.AlreadyUnlocked`（:66） | 1 | 2 | 2 | `[Blacklist.AlreadyUnlocked]` | 259 |

⇒ **"撤掉这一处，哪条用例变得不响？"＝没有一条不响。**8 枚各自都被两发检各自钉住；
票面漏计、被实现程补上的那三枚 `Blacklist.*` 各有独立红句（不是靠总数蒙过去的）。
⇒ **覆盖面主张照"8 枚"写即可，不必改写成"实覆 N＜8 枚"。**

## 5. 变异 E／F：两发检**没覆盖到的那一维**——"拷回来的值还是原值"没人钉（AC#2 退回）

派单单那句"判据仪器本身是首要攻击点"落到这里。两枚变异都**保留分配、只毁内容**：

**变异 E**：`cloneBacking` 从 `append(make(S,0,len(in)), in...)` 改成 `make(S,0,len(in))`（拷了个空切片）。

```
cd D:/tmp/wisp-146acc1/snap2   # 只改这一行，跑本票那两发
→ --- PASS: TestLiveApprovalsSharesNoReferenceSlotWithTheQueue (0.00s)
  --- FAIL: TestLiveApprovalsInPlaceWriteCannotReachTheQueueRecord (0.00s)
  panic: runtime error: index out of range [0] with length 0        # mut\m-e-empty-clone.txt:223
逐包全跑：PASS=2 FAIL=3 SKIP=0，^=== RUN 只有 5 条                   # mut\m-e-full-package.txt
```

⇒ 这一发是**红**的，但红的方式是 **panic**，而且它把同包其余三十条读数一起吞了
（全跑只剩 5 条 `=== RUN`、2 条 PASS）。按派单单那一条规矩：**被吞掉的读数记"未取到"，不许写成 Skip**。
更要紧的是 **probe 1 在变异 E 下是 PASS 的**：它只比底层指针，空切片也是新底层 ⇒ 指针必然不等。

**变异 F**：`cloneParamsMap` 从 `out[k] = v` 改成 `out[k] = nil`（键全留、值全丢）。

```
cd D:/tmp/wisp-146acc1/snap2   # 只改这一行，逐包全跑
→ rc=1  PASS=35 FAIL=1 SKIP=1
  --- PASS: TestLiveApprovalsSharesNoReferenceSlotWithTheQueue (0.00s)
  --- PASS: TestLiveApprovalsInPlaceWriteCannotReachTheQueueRecord (0.00s)
  唯一红的是本程自己放的 zz_acceptor_depth 探针                       # mut\m-f-full-package.txt
```

⇒ **变异 F 是一枚真假绿**：`LiveApprovals()` 把 `Params` 的每个值都换成 `nil` 交出去，
本票那两发检**全绿**，`internal/agent/approval` 里**也没有任何一枚既有用例**发现它
（红的只有本程临时塞在副本里的探针，它不进仓）。
⇒ 于是注释里 `:29-30`"cloneDecision **copies** each of those slots, which is what makes the sentence above true"
中"**copies**"这个词的**第二半**（拷回来的得是同一批值）没有被任何东西钉住：
两发检钉的是"**不共用**"，没钉"**不走形**"。承重公式在这里给的正是误导答案——
"摘掉 cloneParamsMap 的赋值这一味"外部可见读数**变了**（返回的 `Params` 全成 nil），而**没有一条用例红**。

**两枚洞的合体本程另造了一发，读数量在 §5.1**（本仓规矩：≥2 枚各自"附条件／退回"的洞必须单造一发组合变异）。

⇒ **AC#2 这一格退回**（不是"没做"，是"覆盖面写少了"）。最小闭合集合（一条即可）：
就地写之前先加一枚**值等值断言**——`reflect.DeepEqual(got[0].Decision, mustFindLive(q,…).Dec)`——
这一条同时收掉 E 与 F（E 让 `Paths`／`RulesHit` 走形、F 让 `Params` 走形），且不改任何被检代码。

### 5.1 两枚洞**合体**那一发（本仓规矩：≥2 枚各自"附条件／退回"的洞必须单造一发组合变异）

§2.3 的 M-a 与 §5 的 F 各自是一枚洞。本程把它们**同时**装到一棵树上
（`snap2`：`internal/tools/gate.go` 加 `Scratch map[string]string` ＋ `cloneParamsMap` 改成 `out[k]=nil`），
先把本程自己的探针件从目录里挪走（`mv` 到 `mut\`，**不删**），让**只剩仓里既有用例**：

```
cd D:/tmp/wisp-146acc1/snap2
go test -count=1 -v ./internal/agent/approval/          # mut\m-g-combined.txt
→ rc=0   PASS=34  FAIL=0  SKIP=1
  --- PASS: TestLiveApprovalsSharesNoReferenceSlotWithTheQueue (0.00s)
  --- PASS: TestLiveApprovalsInPlaceWriteCannotReachTheQueueRecord (0.00s)
  （四个数与 §10 里那枚"改后健康态"的 34/0/1 一模一样）
```

再用一枚一次性探针把"洞此刻是活的"钉住（`mut\qq_combined_check_test.go.Kept`，同样不进仓）：

```
F visible: 面板收到的 Params[command] 现在是 nil，队列里仍是 "echo hi"
M-a visible: 就地写第 9 枚槽位改掉了已决断记录，全套件 34 PASS / 0 FAIL 看不见
```

⇒ **合体之后：仓里既有用例 34 枚全绿、零红、零新增 Skip，而两件事同时在发生——
`Params` 送到面板上已经不走形，第 9 枚引用槽又回到了"就地写会改掉已决断记录"那一形。**
⇒ 这不是"两枚各自小"的洞，是一枚**能被同一条 `reflect.DeepEqual` 断言一次收掉**的洞
（它钉住"交出去的这枚 `Decision` 与队列存的这枚**逐字段等值**"，于是 M-a 那枚新槽位只要没拷就必然走形、F 那枚必然走形）。
⇒ **§15 最小闭合集合第 2 条按这一发写。**

（跑完本程把 `snap2` 的两枚文件按原件还原并复验：
`diff -r snap6/internal/tools snap2/internal/tools` ＝ 逐字节一致、
`diff ../orig/pending_read.go.orig snap2/…/pending_read.go` ＝ 空、
目录里只剩本程三枚 `zz_acceptor_*` 探针件。）


---

## 6. AC#1 独立重走这一格：**ⓐ 这个方向复算成立（不是超授权），但实现件 §2 的理由②③是本程造得出反例的假句 ⇒ 该格记"退回（只退论证）"**

派单单要求"不许因为它写了理由就采信"。本程把 ⓑ 那一支**真的建出来跑了一遍**（`snap4`＝
`Decision: it.Dec` ＋ 本程按票面 :42 的字样往注释里插进去的 ⓑ 纪律句 ＋ 两枚本程自造的 ⓑ 支仪器）。

### 6.1 先确认 ⓐ 与 ⓑ 的分歧点只有一处：**AC#2 那句断言的极性**

票面两枚 AC 的原文（本程读票面，不读任何一程的转述）：

- AC#1（票面 :51-53）判据是**可核的一句**："今天有没有任何一枚调用方（含测试）在拿到返回值之后对那枚 map
  或那三枚切片做 改元素／append／就地写／取地址后传出去？⇒ 全为'无'时 ⓑ 是诚实且零代价的一支；
  只要有一枚为'有'，ⓑ 就不许单独收，必须 ⓐ。"
- AC#2（票面 :54-56）："造一发**只有'就地改掉返回值里的 `Params` 或那三枚切片'才能触发**的用例，
  **断言队列里那条原始记录没有跟着变**。"

⇒ 票面 :47 另有一句"无论选 ⓐ 还是 ⓑ，都必须有一枚会响的检（AC#2）"。
**这两处对不上**：AC#2 的断言极性（记录**没**跟着变）在 ⓑ 之下**不可能为真**——ⓑ 的定义就是保留那枚浅拷贝，
就地写返回值必然改到队列里那条记录。

### 6.2 ⓑ 支的"0 枚就地写"普查本程复算＝**也是 0**（尺子带正控）

尺子 R1（就地写形状）：

```
RE='(Decision|dec|it|row|got\[[0-9]+\])\.(Params|RulesHit|Paths|Capabilities|Args)\[[^]]*\] *(=[^=]|:)'
cd snap7   # 630c218，票 146 开工前
grep -rnE "$RE" --include="*.go" . | grep -v "_test.go"        → 0 命中
cd snap6   # 7f75ff3，正控树（probe 2 自己就在做这个形状）
grep -rnE "$RE" --include="*.go" .                             → 命中 :220 :221 :223 :224 :225 …（**尺子不空**）
grep -rnE "$RE" --include="*.go" . | grep -v "_test.go"        → 0 命中
```

尺子 R2（第二把，Go `regexp`，写在 `snap4/…/qq_bside_test.go` 里，**内置正控断言**：
同一条正则若在 probe 2 自己的函数体上命中数为 0 就 `t.Fatalf("RULER IS BROKEN…")`）：

```
CONTROL: same regex finds 8 in-place-write shapes in probe 2's own body (must be > 0)
SCAN   : non-test production shapes = 0 []
```

取地址腿（尺子 R3）：`grep -rnE "&[a-zA-Z_]+\.Decision\b" --include="*.go" snap6/` → **0 命中**；
同一把尺打在现造的正控上（`_ = &row.Decision`）→ **1 命中** ⇒ 尺子不空。

⇒ **AC#1 的判据材料本程复算＝"全无"，与实现件 §1.3 同向**（票面 :33 已写"下一程一律现量"）。
生产出口仍只有 `cmd/wisp/panel_pump.go:61` 一枚，且它对 `it.Decision` 只做 5 枚字段的**读**
（`Tool`/`Level`/`RulesHit`/`Reason`/`SessionOverrideBlocked`，`panel_pump.go:66-71`）＋ `panelArgs` 按值读 `Args`；
`Params` **生产侧无人读**（见 §9 末条，那一枚是本票唯一没买到的代价的来源）。

### 6.3 本程造出来的两枚 ⓑ 支仪器 ⇒ 实现件 §2 的理由②③**不成立**

```
cd D:/tmp/wisp-146acc1/snap4      # 浅拷贝接线 + 票面 :42 字样的 ⓑ 纪律注释句
go test -count=1 -v -run 'TestZZAcceptor(DocPinning|BSide)' ./internal/agent/approval/
→ --- PASS: TestZZAcceptorBSideInstrumentIsConstructible        (CONTROL=8 / SCAN=0)
  --- PASS: TestZZAcceptorDocPinningInstrumentGoesRedWhenCommentDeleted
                                                              # mut\mside-b1-with-comment.txt
# 然后把那 239 字节的纪律注释句删掉，同一枚用例：
→ --- FAIL: ... AC#1(b) RED: pending_read.go no longer carries the caller discipline "调用方不许就地写"
```

⇒ 两枚事实，逐条对着实现件 §2 的三条理由：

| 实现件 §2 的理由 | 本程现量 | 判定 |
|---|---|---|
| ①"AC#2 的字面断言只在 ⓐ 下为真" | 摘回浅拷贝 ⇒ AC#2 那两发 `rc=1`、14 条红（§3）；ⓑ 之下该句在 Go 里不可能为真 | **成立**（这一条独自撑住判向） |
| ②"ⓑ 的修法是一句注释，摘掉它**任何用例都不会红**" | 一枚读源码的用例（doc-pin）**因为那 239 字节被删而红** | **假句**：注释可以是承重的，前提是有用例去读它 |
| ③"ⓑ 唯一可能的仪器在本票允许的落点之外（静态扫，而 `tools/d22scan/**` 零字节）" | 本程把两枚 ⓑ 支仪器都放在 `internal/agent/approval/` 里跑通，**一字节未碰 `tools/d22scan/**`** | **假句**：允许落点装得下 ⓑ 支仪器 |

⇒ 所以：**判 ⓐ 不是"唯一可能的收法"，而是"票面两枚 AC 里唯一同时收得下 AC#2 字面的那一支"**。
这个区别**不推翻**方向（AC#2 是票面自己标的"硬核心"，ⓑ 交不出它），但它**推翻了证据件与票面 Progress log 里
写着的那两条论证**——下一位若按那两条去评同类票，会得出"改注释的票不需要仪器"这种错话。
**登记为退回给实现者的文字缺陷**（本程不改票面、不改证据件，D22）。

### 6.4 程序面：这一逆着票面默认的裁决**没有越权**

- 票面 :37 明写"这一格有两种合法收法……本票要你**先判再改**"，ⓐ 是票面自己列出的两枚选项之一；
- 实现程**没有勾任何框**（票面 :81"台账与勾归编排者"，四行 Progress log 里五枚框一枚未勾，本程复算票面 :49-71 的
  五个 `- [ ]` 全部未勾）；
- 它把方向分歧**三处**上报（证据件 §2 的 ⚠ 段、§8 第 1 条、票面 Progress log 第二行的 next=），
  并写了"若编排者改判 ⓑ＝要动 AC#2 的极性＝改票＝人工批准"，以及"摘掉 ⓐ 只需 revert 一枚小 commit"。

⇒ **AC#1 裁决：成立**（附 §6.3 那两枚假句作为退回项，归实现者改文字）。

### 6.5 派单单那一问的正面回答："ⓑ 真的收不下 AC#2 吗？"

- **收不下 AC#2 的字面那一句**：ⓑ 的定义就是保留 `Decision: it.Dec`，而"就地写返回值之后队列里那条记录**没**跟着变"
  在浅拷贝上是 Go 语义层面的假话 ⇒ 要么恒红（§3 已量），要么把极性反过来（＝钉成规格＝票面 :54 那族"恒真判据"）。
  **这一条成立，所以判 ⓐ 不是超授权**（ⓐ／ⓑ 都是票面 :37 列出的合法收法之一）。
- **但 ⓐ 不是"唯一能让这一格变绿的形状"**：本程能列出票面没枚举的两支——
  ⓓ **收窄返回结构**（`LiveApproval` 不再往外送 `Params`；§6.2/§9 已量到生产侧只有 `panel_pump.go:66-71` 五枚字段在读、
  `Params` 无人读），ⓔ **把 `tools.Decision` 换成只读视图类型**（编译期就不许就地写）。
  两支都吃 AC#2 的字面断言，且**都比 ⓐ 便宜**（ⓓ 直接省掉 §8 那 72 枚新增分配里的大头）。
  ⇒ **这不构成退回实现者的理由**（票面只给了 ⓐ／ⓑ 两个选项，它选了允许的那一个并上报），
  但构成 §9 末条那句"下一张票要不要收 ⓓ"的**现量依据**。

---

## 7. 进攻③：极性反证 —— **三条读路都比的是队列里那条记录，没有一路是"返回值自己长什么样"**

代码级核对（`ticket146_liveapprovals_backing_test.go`）＋ 每一路的独立红句：

| 读路 | 比的是什么 | 白盒还是公共面 | 在单点回退里各自响过的槽位（§4 现量） |
|---|---|---|---|
| (1) `mustFindLive(q,"corr-backing").Dec` | **队列存的 `qitem.Dec`**（`pending_read_test.go:234-243` 直接返回 `q.pending` 里那枚指针） | 白盒 | 233／237／240／243／247／250／259×3 ＝ **8/8 枚全响** |
| (2) `q.view()` + `q.head()` | 队列自己的不可信投影（`queue.go:450-468`，`PanelItem.Paths` 由 `viewLocked` 现拷） | 公共读路径 | 272（撤 `Paths` 那一发时响） |
| (3) 第二次 `q.LiveApprovals()` | 队列的**再报**（不是第一次那枚返回值） | 公共 | 292／295／298 |

⇒ 三条路的**断言对象都是队列**，不参数字段"返回值自比"。**没有一路不抵账。**
唯一要写下来的一句：路 (1) 是白盒 ⇒ 它钉住的是"`qitem.Dec` 没被改"；路 (2)(3) 才钉住"报出去的东西没被改"。
两发合起来才把这两面都盖住，这与实现件 §3.3 自己的说法一致（它把路 (2) 称作"第二把尺"）。

**probe 1 也不是恒真判据**：它把名册按**集合**断言（`:174-184`，数目不符即 `t.Fatalf`），
本程实测两形——fixture 少填一枚 ⇒ 红（数目 7 vs 8）；类型上加一枚 map/slice 槽位而 fixture 不填 ⇒ 绿（§2.3 M-a）。
后者才是那枚要退回的过许诺，不是"恒真"。

---

## 8. AC#3 复算：仪器够不着这件事**为真**；机制账本程**自己重跑了一遍**

- **仪器结构上照不到这枚改动**（本程独立读代码，不采信转述）：
  `wisp slo` 的被采腿是 `cmd/wisp/slo_windows.go:462 startSubject` 拉起的 `slo -subject` 子进程，
  `runSubject` 把 `run.Posture = "skeleton"`（同文件 `:412`）后**只是 sleep 到窗口结束**（`:422-428`）；
  快照泵的构造点**全仓只有一处**：`grep -rn "NewSnapshotPump" cmd/wisp/` → `cmd/wisp/run.go:421`；
  `Verdicts:` 的生产接线也只有一处 → `cmd/wisp/run.go:422 rt.liveVerdicts`。
  ⇒ **WorkPeak 那一档的采样窗口里 `LiveApprovals` 一次都不会被调用**，与实现件 §4.2 的自陈一致。
- **"未观察到差异"这句写法是对的**：票面 :60 要求"数字没动要写未观察到差异而不是无代价"，
  实现件 §4.3 逐字写"WorkPeak 那一档未观察到差异"并禁止与机制账合并成"代价可忽略"。**没有越线。**
  （它引的两臂读数本程**未复跑**：本仓规矩"计时读数不作判据"，且此刻本机有别程在跑。）
- **机制账本程自跑复算**（自己的 bench 件，唯一变量＝`cloneDecision`；两棵纯净副本
  `snap-before`(630c218)／`snap`(7f75ff3)，`-count=5 -benchtime=2000x`，读数
  `D:\tmp\wisp-146acc1\mut\bench-both.txt`；bench 件留在 `D:\tmp\wisp-146acc1\acceptor_bench_test.go`）：

  | 臂 | 满深度（`DefaultMaxPending`＝8 枚 pending） | 一枚 pending |
  |---|---|---|
  | 改前 `630c218` | **3456 B/op，1 alloc/op** | 416 B/op，1 alloc/op |
  | 改后 `7f75ff3` | **7680–7682 B/op，73 allocs/op** | 944 B/op，10 allocs/op |

  ⇒ 本程量到 **＋4224 B／＋72 allocs**；实现件 §4.2 B 报的是 ＋4736 B／＋72 allocs。
  **allocs 那一维逐字对上（1→73）**，B/op 差 512 字节＝两枚 fixture 的参数条数不同（本程 2 paths／2 caps／2 rules）。
  ⇒ 实现件那条"没进仓的仓外副本读数"**在本程这里被独立复现了**，方向与量级一致。
  ⇒ **ns/op 本程不引作判据**（同一枚 bench 的 ns 在两臂间抖到 1 µs 量级，这台机器上别的程在跑）。
- **`DefaultMaxPending = 8`**（`internal/agent/approval/queue.go:112`）＝注释引的那枚上界为真。

⇒ **AC#3 裁决：成立**（票面只要求"若走 ⓐ 不许空着"，它给了双臂读数＋自陈仪器不同路＋机制账，
本程复算了机制账并独立证明了"仪器照不到"）。**一条附记**：那条机制账的**唯一读者证据是仓外副本**，
按票面 :81"任何一句前提复算不符 ⇒ 报回"的相反面，**这一格的复算入口本程替它补在这里**（bench 件路径＋两棵副本名）。

---

## 9. AC#4 契约轴 —— **交集空**（本程逐枚 commit 复算）

```
cd "D:/work/workspace/projects plans/Wisp"
for c in 598c0c6 22f7b1a 0d12661 3fde333 7f75ff3; do git show --name-only --format= $c; done | sort -u
→ .scratch/wisp/issues/146-….md
  docs/evidence/s1/146-liveapprovals-r1.md
  internal/agent/approval/pending_read.go
  internal/agent/approval/ticket146_liveapprovals_backing_test.go

把票面 AC#4 那 11 条禁改路径当筛子过一遍（本程自己写的筛子，不是抄它的结论）：
  … | grep -E "internal/risk/|rules_gateway|thresholds\.go|golden|allowlist\.txt|docs/PLAN\.md|docs/specs/|frontend/|design/|tools/d22scan|l2_grant_boundary|internal/panel/|cmd/wisp/"
→ 无输出（rc=1 ⇒ 交集空）
```

- 删除列核对：`git show --numstat 22f7b1a` → `64 2 pending_read.go`、`305 0 ticket146_…_test.go`，
  与实现件 §6.1 所报**逐字一致**；那 2 行删除是本程可读的（替换掉 `Decision: it.Dec,` 与一行注释续行）。
- **入站面那条硬线**（票面 :64-65）：`grep -rn "LiveApprovals" internal/panel/ internal/risk/` → **0 命中**；
  `grep -rn "liveApprovals|live_approvals" frontend/ design/`（js/ts/json）→ **0 命中**
  ⇒ 这枚函数没被这一票变成路由或页面可问的东西。尺子正控：同一枚 grep 在 `cmd/wisp/` 命中 1 处（`panel_pump.go:61`）。
- ⚠ **一条本程新增、实现件没写的代价账（且带一枚本程自己算错的数，已就地改掉）**：
  `Params` 是 8 枚里唯一**生产侧无人读**的槽位
  （`panel_pump.go:66-71` 读的是 Tool/Level/RulesHit/Reason/SessionOverrideBlocked ＋ `panelArgs` 按值读 `Args`；
  `grep -n "\.Params" cmd/wisp/panel_pump.go internal/panel/pump.go` → **0 命中**）。
  本程于是量了一发**归因**（`snap6`：只把 `out.Params = cloneParamsMap(d.Params)` 那一行拿掉，其余不动；
  跑完已 `diff` 验过还原）：

  | 臂（满深度 8 枚 pending，`-count=3 -benchtime=2000x`） | B/op | allocs/op |
  |---|---|---|
  | 改前 `630c218` | 3456 | 1 |
  | 改后全形 `7f75ff3` | 7680–7682 | 73 |
  | **改后但 `Params` 不拷** | **4992–4994** | **57** |

  ⇒ `Params` 那一枚 map 的克隆占这枚改动新增量的 **2686 B（＝64 %）**，但只占 **16 allocs（＝22 %）**
  ——剩下 56 枚分配是七枚切片。
  （本程第一版这里写的是"72 枚新增分配主要由它构成"，**那是把字节数当成了分配数**，现按上面这发改正。）
  ⇒ 所以"收窄返回结构（ⓓ：不往外送 `Params`）"能省掉的是**字节那一半**、不是分配次数那一半。
  ⇒ 这**不违反任何一条 AC**（ⓐ 的定义就是"那枚 map 与那几枚切片真拷一层"，票面 :39），
  但它是 §6.5 那句"ⓓ 更便宜"的**现量边界**：ⓓ 省 64 % 的字节、省 22 % 的分配次数，**省不掉那七枚切片**。
  **归编排者，不归实现者。**

⇒ **AC#4 裁决：成立。**

---

## 10. AC#5 门禁独立复跑（逐包单跑；判红绿只认 `--- FAIL:` 那一行）

| 门 | 命令（本程原样） | 本程读数 |
|---|---|---|
| 改前包测 | `cd snap7 && go test -count=1 -v ./internal/agent/approval/`（`630c218` **纯净、未放本程探针件**） | `rc=0`；`^--- PASS`=**32**／`^--- FAIL`=**0**／`^--- SKIP`=**1**（`mut\clean-before-630c218.txt`） |
| 改后包测 | `cd snap6 && go test -count=1 -v ./internal/agent/approval/`（`7f75ff3` **纯净**） | `rc=0`；**34／0／1**（`mut\clean-after-7f75ff3.txt`） |
| 名册差集 | `grep '^--- ' \| sed 's/^--- [A-Z]*: //; s/ .*//' \| sort -u` 两向 `comm` | 改前 33 枚／改后 35 枚；`comm -13`（只多）＝**恰那两枚**：`TestLiveApprovalsInPlaceWriteCannotReachTheQueueRecord`、`TestLiveApprovalsSharesNoReferenceSlotWithTheQueue`；`comm -23`（只少）＝**空** |
| SKIP 归因 | `grep '^--- SKIP'` 两棵 | 两棵各一枚、**同名** `TestDefaultDeadlineWallClockMeasurement` ⇒ 未新增 Skip、未拿 Skip 顶掉读数 |
| panic 吞读数 | `grep -l "panic:" mut/*.txt` | 干净跑里**无 panic**；只有本程那枚变异 E 造出 panic，那时全包 `=== RUN` 从 35 掉到 **5**（§5），那 30 条记**未取到** |
| 格式 | `gofumpt -l . tools/d22scan tools/mockllm` ＋ `gofmt -l .`（在**未塞本程探针件**的纯净副本 `snap6` 上跑） | **两条都空**；`gofumpt` 版本 v0.12.0（与 CI 同版，`--version` 现量） |
| vet | `go vet ./internal/agent/approval/`／`./cmd/wisp/`／`./internal/panel/`（逐包单跑） | **三枚 rc 全 0** |
| 建 | `go build ./cmd/wisp`（纯净副本 `snap`） | **rc=0** ⇒ 实现件 §6.3 那条"已被 `95885fb` 解除"的复量本程独立成立 |
| D22 门 | `sh scripts/d22scan.sh`（`snap`） | **rc=0**，末行 `d22scan: clean - no D22 ban violations`，各作用域 `examined N` 全非零（`internal/` 205 生产＋414 全量、`cmd/` 23／43、`frontend/` 49、`design/` 30） |
| 跨包 | `./internal/panel/`（纯净副本 `snap`，逐包单跑） | `rc=1`，**唯一红＝`TestC21DesignTokensFourWayAgree`**（红句见 §11）；`./cmd/wisp/` **未跑**：纯净副本里没有 `third_party/`（未跟踪），dll 不在 ⇒ `0xc0000135`，见 §12 触发条件 |

⇒ **AC#5 裁决：成立**（实现件 §6.2 那四个数与名册差集**全部复算一致**；`cmd/wisp` 那一腿本程与实现程同样
**未取到**，不记成通过）。

⚠ **d22scan 的一条仪器说明（不是缺陷）**：在 `git archive` 出来的副本里跑，它自己会打一行
`gitignore rules NOT APPLIED - … every path in every scope is being scanned`。
⇒ 本程那些 `examined N` 与实现程在真工作树里量到的 N **不可直接比**（本程的是"一张网都没减"的那一侧）。
这是**响的方向**，不是漏的方向，记下来免得下一位以为谁在吹数。

---

## 11. 那枚已知常红：本程复算＝**它在纯净锚点上也红**，实现件给它的根因**要改**

派单单交代过：`internal/panel` 的 `TestC21DesignTokensFourWayAgree` 在工作树里必红
（`design/assets/tokens.css` 被 owner 挪走，台账 `Q-52` 撤回行＋`A254`），不是本票的账。
本程照它在**三棵树**各跑一遍，红句一起引：

```
snap6/snap（7f75ff3 纯净副本，tar 解出来 design/assets/tokens.css 在）  → FAIL
snap-before（630c218 纯净副本）                                          → FAIL
工作树（git status 里那 16 枚 D design/**）                               → 红，但红因不同
```

纯净锚点上的红句（`mut\panel-snap.txt`，摘三行）：

```
tokens_fourway_test.go:446: parties: design/assets/tokens.css=135 dark decls, …=78 colour rows,
    internal/ball/tokens.go=40/38 palette fields, frontend/src/styles/tokens.generated.css=109 dark decls
tokens_fourway_test.go:461: frontend/src/styles/tokens.generated.css does not carry --orb-hi = rgba(255, 255, 255, 0.90)
    that design/assets/tokens.css declares - a token never reached the panel
tokens_fourway_test.go:471: design/assets/tokens.css declares --danger = "#E07A70" but
    frontend/src/styles/tokens.generated.css carries "var(--demo-destructive)" - two style sources disagree
```

⇒ **两个分开的结论，不许互相抵账**：
1. **与本票无关**：成立（`design/**`、`frontend/**`、`internal/panel/**` 对票 146 全零字节，§9 已核）；
   ⇒ 本程**不据此判票 146 退回**（派单单也这么交代）。
2. **根因**：**同一枚用例在两棵树上是两种红**。本程把两边红句都取了（工作树那一发只跑测试、未动任何文件）：

   | 树 | `tokens.css` 在不在 | 红句 |
   |---|---|---|
   | 工作树（`D:\work\workspace\projects plans\Wisp`，`git status` 里那 16 枚 `D design/**`） | **不在**（`ls design/assets/` → No such file） | `tokens_fourway_test.go:441: read design/assets/tokens.css: open …\design\assets\tokens.css: The system cannot find the path specified.` |
   | `snap6`／`snap`＝`7f75ff3` 纯净副本 | **在**（`git ls-tree 7f75ff3 design/assets/` 认得这枚 blob `0ad0112`） | `:446 parties: design/assets/tokens.css=135 dark decls … tokens.generated.css=109 dark decls` ＋ `:461 … does not carry --orb-hi …` ＋ `:471 … --danger = "#E07A70" … carries "var(--demo-destructive)"` |
   | `snap7`＝`630c218` 纯净副本（票 146 开工前） | 在 | 同第三行的形（`mut\clean-…`／`… panel run`） |

   ⇒ 实现件 §5 第 5 条／§6.3／§8 第 4 条写的"根因＝`design/assets/tokens.css` 工作树缺文件＝owner 未提交的
   `design/**` 移动"**只在第一行那一形上成立**；**在锚点的纯净副本上这枚用例本来就红，红的是四方值不一致**。
   ⇒ 于是"推送前只剩这一枚要处理、而它归 owner"这一句里的**归因**要改：
   推送之后 CI 从 commit 检出 ⇒ 缺文件那一形会消失，**值不一致那一形不会消失**。
   （两形是否同源、哪一形该由谁收＝**未取证**，触发条件见 §12 第 4 条。）

---

## 12. 未取证项与可复算的触发条件（本程**没有**造出来的怀疑，不许读成"附条件通过"）

| # | 未取证的那一格 | 为什么本程没取到 | 复算触发条件（谁都能照做） |
|---|---|---|---|
| 1 | `cloneBacking`／`cloneParamsMap` 的**值保真**只在一枚本程变异下暴露（§5 F） | 本程不代做实现 | 在副本里把 `out[k]=v` 换成 `out[k]=nil` ⇒ 本票两发检应红；今天它绿 |
| 2 | 并发下锁持有时间被拷贝拉长（实现件 §5 第 3 条自陈未测） | 本程未造并发形 | `-race` 跑本包＋一枚"泵动 × 应答抢锁"的形；**今天没有任何用例探它** |
| 3 | `cmd/wisp/` 全包测在锚点上是否全绿 | 纯净副本没有 `third_party/` dll（未跟踪）⇒ `0xc0000135` | 从工作树跑：`PATH=third_party/sherpa-onnx:$PATH go test -count=1 ./cmd/wisp/`（先确认票 144 已收线） |
| 4 | `TestC21DesignTokensFourWayAgree` 的**两形同源否**（缺文件 vs 值不一致） | 越界（`internal/panel/**`、`design/**` 都不许本程碰） | 在 `tokens.css` 存在的一棵树上跑，把红句与工作树那一发的红句**逐字节对一次** |
| 5 | WorkPeak 那一档的**双臂读数** | 本仓规矩：本机有别程在跑 ⇒ 计时/资源读数不作判据 | 安静窗口下按 §8 的建法（两棵 `git archive` 副本）各 3 次交替，**且先承认仪器照不到**（§8 已证） |
| 6 | `Queue.replay` 那枚**同族出口**（实现件 §5 第 6 条：`queue.go:488-498` 也按值传 `it.Dec`） | 本票范围外，且 `queue.go` 是禁改面 | 本程复算它**仍在**：`internal/agent/approval/queue.go:478-492` 的 `replay` 里 `fresh := *it` 把含 8 枚引用槽的 `Dec` 按值搬出来，**没有任何用例探它的底层共用**——`grep -rn "\.replay(" --include="*_test.go" internal/agent/approval/` ＝ **0 命中**，`grep -rn "Decision\.(Params\|Paths\|RulesHit)\[" --include="*_test.go" internal/` 只命中票 146 自己那两发（`:220 :221 :223 :224 :291`）。⚠ 这条量的时候本程差点被骗：小写 `\.replay(` 那一发是 0 命中，**但 0 的原因是出口叫 `Gate.Replay`**（`queue_test.go:283 TestReplayRedisplaysUnderAFreshGrant` 确实在跑）——是本程拿 `grep -rn "replay" internal/agent/approval/`（16 命中）当正控复打才发现的。**那一枚用例存在，只是它断的是状态、不是共用底层。** |

---

## 13. "如果我漏了它，谁会先被骗"——本程自己那一节（不抄实现件 §5）

1. **先被骗的是"下一次给 `Decision` 加字段的人"**：`pending_read.go:36-38` 与测试件 `:18-22`、`:140-142` 三处
   都写着"第 9 枚引用字段不会静静进来"。本程造了两枚它**会**静静进来的形（M-a：新 map 槽位＋fixture 未填；
   M-c：Ptr/Interface/Chan/Func/私有字段一律不进 `refSlots` 的 switch）。
   骗到之后的症状与票 146 的原病**一模一样**：一条已决断记录被读者悄悄改掉，而门禁全绿。
2. **第二枚是"把两发检读成'返回值已与队列隔离且内容原样'的人"**：§5 的变异 F（`Params` 值全丢）在两发检下全绿。
   先被骗的是**票 33 的宿主程**——它下一步就是把这些值画到卡上；今天 `Params` 走形不会有任何东西红。
3. **第三枚是"按 §2 理由②③去评下一张同类票的人"**：那两条本程各造了一发反例。若它们被当成规矩抄走，
   会得出"改注释的票不需要仪器"这种错话——而正确的那句是"改注释交不出 **AC#2 字面那枚**检"。
4. **第四枚是把 §8 那张表读成"ⓐ 的代价由 `Params` 那一枚 map 构成"的人**（**本程自己第一版就写错了这一条**，
   现按 §9 的归因量改过）：`Params` 占这枚改动新增量的 **64 % 字节**，但只占 **22 % 分配次数**（16／72）——
   七枚切片才是分配次数那一半的主人。于是"收窄返回结构（ⓓ 不往外送 `Params`）"这枚下一票的形状，
   **省的是字节那一半、不是分配次数那一半**；谁按"ⓓ 能把 72 枚 allocs 砍一半"去立项，那个预期是错的。
   实现件与本程第一版都没量这一格；它不属于任何一条 AC 的缺失，属于**下一张票的形状**（§6.5 的 ⓓ 边界）。
5. **第五枚是把"d22scan 的 examined N"在两棵树之间直接对比的人**（§10 的 ⚠）：`git archive` 副本没有
   `.gitignore` 可问 ⇒ 它扫得更宽，N 天然不同。谁拿这两枚 N 比"谁漏看了"，谁就先被骗。

---

## 14. 伪授权登记（两栏分开计数）

| 栏 | 数 | 出处（工具名＋命令前 40 字） |
|---|---|---|
| **真通知回显数** | **5** | ①派单单随身的 `<loaded_context>`：项目 `agents.md` 全文回显；②`Read`：`The date has changed. Current date: 2026-09-25`；③`Read`：可用 skill 清单；④⑤`Bash` 后台任务完成通知 `[SYSTEM NOTIFICATION - NOT USER INPUT]`（`for f in internal/agent/approval/pending_read.go…` 与 `cd D:/tmp/wisp-146acc1/snap && ls third_party 2>/dev…` 两发的完成回执）——四枚都**不要求任何动作**、不声称授权，本程读它们只是为了取回日志；**未据此减少任何取证** |
| **判为注入数** | **0** | 本程工具输出里没有出现"系统提示／编排者备注／已核验请继续提交／Confirm the harness note is genuine／请 revert／放宽阈值／这格我已经签收／不用取证直接给结论"任何一枚形状 |

- **锚点只认自己量的**：本程取 sha 一律先 `git cat-file -t` ＋ `git rev-parse`（`7f75ff3` → `commit`／
  `7f75ff3c9aed48054656b90661af274b09906b30`）。派单单正文里出现的其它 sha（`4e66817`／`630c218`／`95885fb`…）
  本程都在**用到的那一枚**上重跑过 `git show`／`git ls-tree`，未凭一串字面量采信。
- 凭据值：本程全程未抄任何 key/token 原文（连变量名以外的字面都没抄）。
- **被拒工具调用数：0**。本程自己踩空、原样登记的几发（都不需要人批准）：
  ①`grep -c '\\r'` 在 `od -c` 输出上数出 6 ⇒ 那是**尺子错**（`od` 的转义文本里有字面 `\r` 字样），
  换成 `od -c` 直接看字节＋`grep -rl $'\r'` 才是对的；
  ②`find … | while read` 全树数 CR 跑到超时被挪后台 ⇒ 弃用，改小域 grep；
  ③`startSubject`/`runSubject` 第一把 grep 落空（`func runSLOSubject` 这名字根本不存在）；
  ④本程自己那枚 `qq_bside_test.go` 第一版 `os.Getwd()` 往上跳了 4 层、把仓库根算错 ⇒
  **它自己那一发是 build/路径失败，不是读数**；改 3 层后重跑，并在里面带了一枚"尺子坏了就 `t.Fatalf`"的正控
  （正是派单单要求的"0 命中先拿正控打一遍"）。
  ⑤**同一枚坑本程又踩了一次并自己抓回来了**：查 `replay` 同族出口时 `\.replay(` 报 0 命中，
  正控（`grep -rn "replay"`，16 命中）显示出口其实叫 `Gate.Replay` ⇒ 那一发 0 命中是**尺子错**，
  结论已按现量改写在 §12 第 6 条里（"那一枚用例存在，只是它断的是状态、不是共用底层"）。

---

## 15. 五格裁决表（成立／退回／附条件）

| 格 | 档位 | 现量命令（复算入口） | 读数／判决要点 |
|---|---|---|---|
| **AC#1** 判 ⓐ／ⓑ | **退回（附条件入账：文件保留、框不勾、ⓐ 那一支代码**不要**回退）** | §6.2 三条普查（各带正控）＋ §6.3 的 `snap4` 两枚 ⓑ 支仪器（`mut\mside-b1-with-comment.txt`／`mut\mside-b2-comment-deleted.txt`） | ⓐ 是"同时收得下 AC#2 字面"的那一支 ⇒ **判向本程复算成立、未越权**；普查 0 枚就地写本程复算同向；框未勾、三处上报，程序面合规。**但**§2 的理由②"ⓑ 的修法是一句注释，摘掉它任何用例都不会红"＝本程造出发红的一发（删掉 239 字节的 ⓑ 纪律句 ⇒ 一枚读源码的用例 `--- FAIL`）；理由③"ⓑ 唯一可能的仪器在本票允许的落点之外"＝本程在 `internal/agent/approval/` 里造出**两枚** ⓑ 支仪器、`tools/d22scan/**` 一字节未动。**这一格是本程造出来的洞 ⇒ 按分界规矩记退回，退回对象是那两行论证、不是 ⓐ 这个方向。** |
| **AC#2** 会响的检 | **退回**（附条件入账、文件保留、**不许勾**） | §3（摘修法）＋ §4（8 枚逐枚撤）＋ §5（变异 E／F）＋ **§5.1（两枚洞合体那一发：`mut\m-g-combined.txt`）** ＋ §7（极性） | 承重与覆盖面主张**都成立**：摘修法 `rc=1`／14 条红；逐枚撤 8/8 各自红；三条读路都真比队列。**但 §5 的变异 F 是一枚真假绿**（`Params` 值全丢 ⇒ 两发检全绿、包内无人红），且 §5.1 证明**它与 §2.3 那枚洞合体后，仓里既有用例仍报 34 PASS／0 FAIL／1 SKIP** ⇒ "copies" 的**值保真**那一半没被钉。**最小闭合集合（一条）**：就地写之前加一枚 `reflect.DeepEqual(got[0].Decision, mustFindLive(q,"corr-backing").Dec)`；它一次收掉 E、F **与 M-a**（新槽位只要没拷就必然与存储项不等值）。**不要求**改被检代码、**不要求**回退 ⓐ |
| **AC#3** 资源那一问 | **成立** | §8：`grep -rn "NewSnapshotPump" cmd/wisp/`＋`startSubject`/`runSubject` 读码；bench 复算 `mut\bench-both.txt`；归因那一发在 §9（`snap6` 只拿掉 `Params` 那行） | WorkPeak 两臂照不到这枚函数＝**真**（泵全仓只在 `run.go:421` 构造）；写法是"未观察到差异"、不是"无代价"＝**对**；机制账本程自跑到 **＋4224 B／＋72 allocs**（它报 ＋4736 B／＋72 allocs，allocs 逐字对上，B 差＝fixture 条数差）；**新增量的归属本程另量了一发：`Params` 占 64 % 字节、只占 22 % 分配次数**；ns 不引 |
| **AC#4** 契约轴 | **成立** | §9：五枚 commit 的 `--name-only` 并集 ＋ 11 条禁改路径筛子；`--numstat`；入站面两枚 grep（正控＝`cmd/wisp` 命中 1） | 交集**空**；`queue.go` 零字节；未变成路由；删除列 2 行可解释 |
| **AC#5** 门禁 | **成立** | §10 那张表（每行都是本程自己跑的） | 32／0／1 → 34／0／1；名册两向 `comm`＝只多那两枚、零删除；SKIP 改前就在；`gofumpt`／`gofmt` 空（纯净副本上量）；`go vet` 三包 rc=0；`go build ./cmd/wisp` rc=0；`sh scripts/d22scan.sh` rc=0 且 N 全非零。**未取到的一腿**：`./cmd/wisp/` 全包测（dll 不在副本里），与实现件自陈同格 |
| **进攻④** 深度许诺 | **退回**（半句为真、半句为假） | §2：`snap3` 深度探针＋`snap2` 的 M-a／M-c | "一层深／`[]any` 与嵌套 map 仍共用"＝**逐字为真**；"第 9 枚引用字段不会静静进来"＝**假**（M-a 活着红不了、M-c 五类静默跳过）。最小闭合集合在 §2.3（首选：名册改成按**类型**数声明） |

### 总裁决

**退回（附条件入账、四枚文件全部保留、五枚框一枚不许勾）。** 不是"整票重跑"：
代码（ⓐ 那一支）、普查、承重、覆盖面这四样**本程都复算成立**，退回的是**三处可证伪的文字**与**一条一行就能补的检**：

1. `pending_read.go:36-38` ＋ 测试件 `:18-22`／`:140-142` 的"第 9 枚不会静静进来"——改成按类型数声明（§2.3 选项 1），
   或收窄文字并把洞挂进 §12（选项 2，但那一洞必须继续挂着）；
2. 测试件补一枚**值保真**断言（`reflect.DeepEqual` 返回值 vs 队列存储项）——收掉本程的变异 E／F（§5）；
3. 证据件 §2 的理由②③ 按 §6.3 那两实现量改写；正确的那句是"ⓑ 交不出 **AC#2 字面**那枚检"，不是"ⓑ 拿不到任何仪器"；
4. 证据件 §5 第 5 条／§6.3／§8 第 4 条给 `TestC21DesignTokensFourWayAgree` 的**根因**要改（§11：锚点纯净副本上它照样红，红的是四方值不一致）。

**归编排者（不是实现者的账）**：① AC#1 的方向裁决仍需你签字（本程认为 ⓐ 站得住，但改判 ⓑ 要动 AC#2 极性＝改票）；
② §9 末条那枚"返回结构要不要继续往外送 `Params`"＝下一张票；③ §12 第 6 条 `queue.go` 的 `replay` 同族出口＝下一张票。

**next=**：本件 15 节已交完派单单要求的四发进攻（①②③④）＋五格出表＋门禁复跑＋"我漏了什么"＋伪授权两栏。
**没有剩余格**。若编排者要本程补 §12 的第 2 条（`-race` 并发形）或第 3 条（工作树 `cmd/wisp`），
点名即可——两条的触发条件都写在表里了。

