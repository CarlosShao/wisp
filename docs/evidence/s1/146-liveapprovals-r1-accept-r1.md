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
| `snap` | `7f75ff3` 纯净 ＋ 摘掉修法的变异（已还原，`diff` 证 RESTORE-CLEAN） ＋ 本程探针件 | 进攻①、门禁复跑、名册差集 |
| `snap-before` | `630c218` 纯净（票 146 开工前） | AC#5 改前基线、AC#1 普查的改前复算 |
| `snap2` | `7f75ff3` ＋ 本程的**按类型**字段普查件 ＋ 各枚变异（全部已还原） | 字段普查、进攻④、变异 E／F |
| `snap3` | `7f75ff3`，**代码一字节未改** ＋ 深度探针件 | 进攻④的"一层深"复算 |
| `snap4` | `7f75ff3` ＝ 摘回 `Decision: it.Dec` ＋ 本程自造的 ⓑ 支仪器 | AC#1 反证（ⓑ 到底收不收得下 AC#2） |
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

**两枚洞的合体已经单造过一发**：M-a 同时满足"新槽位＋fixture 未填＋`cloneDecision` 未拷"，
结果是 §2.3 里那五行全绿＋洞活着的组合形；本节的 F 是另一形（洞不在"共用"，在"走形"）。

⇒ **AC#2 这一格退回**（不是"没做"，是"覆盖面写少了"）。最小闭合集合（一条即可）：
就地写之前先加一枚**值等值断言**——`reflect.DeepEqual(got[0].Decision, mustFindLive(q,…).Dec)`——
这一条同时收掉 E 与 F（E 让 `Paths`／`RulesHit` 走形、F 让 `Params` 走形），且不改任何被检代码。

<!-- 本件下面各节尚未取证；next= 见文件末尾 -->
