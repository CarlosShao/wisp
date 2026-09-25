# 146-liveapprovals-r1-accept-r1 — 票 146 的对抗验收裁决表（裁决者 ≠ 实现者）

- 被验对象：票 146（`Queue.LiveApprovals()` 的 "copies" 许诺）
- **被验锚点：`7f75ff3c9aed48054656b90661af274b09906b30`**（`git rev-parse 7f75ff3` 现量；`git cat-file -t` = `commit`）
- 本件作者：**验收程**（`agent=ticket146-accept-r1`），**不是**写 `cloneDecision` 的那一程
- 时刻：`2026-09-25 19:3x–20:0x +08`
- 本件**只新增这一枚文件**；未碰 `internal/**`、票面、阈值/golden、`tools/d22scan/**`、`internal/risk/**`、
  `frontend/**`、`design/**`、`internal/panel/**`、`cmd/wisp/**`（写权边界＝派单单上那一条）

## 0. 被验版本怎么取的（复算入口）

脏工作树**不当被验版本**（当时 `cmd/wisp/**` 正被票 144 落盘）。取只读纯净副本：

```
mkdir -p D:/tmp/wisp-146acc1/snap
git -c core.autocrlf=false -c core.eol=lf archive 7f75ff3 | tar -x -C D:/tmp/wisp-146acc1/snap
```

两个 `-c` 都带了，且**复算过烙没烙进 CRLF**（那一坑是真的）：

```
head -c 60 snap/internal/agent/approval/pending_read.go | od -c   → 只有 \n，无 \r
grep -rl $'\r' --include='*.go' snap/ | wc -l                     → 0
```

副本清单（**全部在仓外**，全部只建不删）：

| 目录 | 是什么 |
|---|---|
| `D:\tmp\wisp-146acc1\snap` | `7f75ff3` 纯净副本（变异过，已逐字节还原：`diff orig/pending_read.go.orig snap/…` = RESTORE-CLEAN，`diff -r snap/internal/agent/approval snap2/…` 一致） |
| `D:\tmp\wisp-146acc1\snap-before` | `630c218` 纯净副本（改前基线，AC#5 名册用） |
| `D:\tmp\wisp-146acc1\snap2` | `7f75ff3` 第二副本（字段普查/第 9 枚槽位变异用；`internal/tools/gate.go` 变异后**已还原**，变异件留在 `mut\gate.go.M-a`） |
| `D:\tmp\wisp-146acc1\snap3` | `7f75ff3` 第三副本（只放深度探针，代码一字节未改） |
| `D:\tmp\wisp-146acc1\orig` `…\mut` | 还原用的原件副本 + 每一发的读数日志（下面逐格点名） |

**本程在副本里放的三枚探针件**（`zz_acceptor_census_test.go`／`zz_acceptor_ninth_test.go`／`zz_acceptor_depth_test.go`）
是**验收仪器、不是交付物**，一枚都不进仓（D22：裁决者不实现）。

---

## 1. 字段普查独立复算（票面 5 枚 vs 实现程 8 枚）

裁决者自己数一遍，**不采信任何一方的表**。命令＝在 `snap2` 的 `internal/agent/approval/` 放一枚按**类型**
（不是按值）走 `reflect` 的普查件，跑：

```
cd D:/tmp/wisp-146acc1/snap2
go test -count=1 -v -run TestZZAcceptorCensusTypeShape ./internal/agent/approval/
```

读数（`log.txt` 未另存，全量在下表；`tools.Decision` 顶层 **19 枚字段**，kinds=`map:1 slice:4 struct:1 string:8 bool:2 int:2 int64:1`）：

| 顶层字段 | 类型（reflect 现量） | kind | 引用槽？ |
|---|---|---|---|
| `Params` | `map[string]interface {}` | map | **① ** |
| `Args` | `jsontext.Value`（＝`json.RawMessage` 在 go1.27 里的本体） | **slice** of `uint8` | **②** |
| `RulesHit` | `[]risk.RuleID`（elem kind=string） | slice | **③** |
| `Paths` | `[]string` | slice | **④** |
| `Capabilities` | `[]tools.Capability`（elem kind=string） | slice | **⑤** |
| `Blacklist` | `tools.BlacklistNote`（struct，6 枚字段全 exported） | struct | 本身不是槽，**往下数** |
| └ `Absolute` / `Unlockable` / `AlreadyUnlocked` | `[]string` ×3 | slice | **⑥⑦⑧** |
| └ `Class`(`risk.Class`=int) / `Reason`(string) / `Unresolved`(int) | — | — | 否 |
| 其余 12 枚（`Tool`/`Provider`/`Level`/`Reason`/`SessionOverrideBlocked`/`Mode`/`ModeSilenced`/`ModeKept`/`DecisionColumn`/`CorrelationID`/`TaskID`/`CallID`/`Timeout`） | string×8、int×2、bool×2、int64×1 | — | **否** |

正控（先验尺子再用量程）：

```
CONTROL json.RawMessage kind=slice elem=uint8      ← 确认 RawMessage 本质是切片，不是"别的什么"
TOP total fields=19 kinds=map[bool:2 int:2 int64:1 map:1 slice:4 string:8 struct:1]
NESTED Blacklist (tools.BlacklistNote) fields=6    ← 六枚全部 exported，没有"藏在私有字段里"的死角
```

⇒ **判定：8 枚（1 map ＋ 7 slice）成立，没有第 9 枚。**
实现程比票面（5 枚）**多找出的 3 枚**＝`Blacklist.Absolute`/`Unlockable`/`AlreadyUnlocked`，位置与类型逐枚对上。
`Capabilities` 那行票面写"元素本票没量"⇒ 本程量了：**elem kind=string**，元素不含引用，
所以 `cloneBacking` 逐元素拷值即足够，"头共用"这一维是真闭上了。
（另有一枚正控值得记：`reflect.ValueOf([]byte(nil)).Pointer()` 与 `reflect.ValueOf([]byte{}).Pointer()` **都返回 0**
⇒ probe 1 的 `pa != 0 && pa == pb` 判据对"空但非 nil"的切片是**瞎的**。今天 8 枚槽位在 fixture 里都是长度 ≥1，
所以不咬人；这条记进 §7。）

---

## 2. 进攻④：深度许诺 vs 代码 —— **不一致的那一枚是"第 9 枚不会静静进来"那句**（退回）

### 2.1 注释原文（`snap/internal/agent/approval/pending_read.go`，逐行号）

```
:32-35  The copy is ONE level deep: what a caller cannot now do is write a slot the queue
        reports; what it can still do, and what no Go type here can prevent, is reach
        THROUGH Params' values into the containers a JSON decode left inside them
        (a []any or nested map[string]any is still shared). Treat those values as read-only too.
:36-38  ticket146_liveapprovals_backing_test.go is the instrument for the eight named slots,
        and it walks tools.Decision by reflection so a ninth reference-typed field
        cannot arrive quietly.
```

### 2.2 "一层深"那一半：**逐字成立**（量出来的，不是读出来的）

`snap3`（代码零改动）放一枚深度探针，命令：

```
cd D:/tmp/wisp-146acc1/snap3
go test -count=1 -v -run TestZZAcceptorDepthPromiseMatchesCode ./internal/agent/approval/
```

读数（`D:\tmp\wisp-146acc1\mut\m-d-depth.txt`）：

```
:40 A-CONFIRMED: queue's Params[argv][0] was rewritten through the returned row
    ([]interface {}{"CLOBBERED", "hi"}) -> the comment's '[]any is still shared' is TRUE of the code.
:47 A-CONFIRMED: queue's Params[nested][depth] rewritten too -> 'nested map[string]any is
    still shared' is TRUE of the code.
:65 B: append on the returned Paths slice changed queue len from 1 to 1 (must be equal)
```

⇒ 注释说"仍能穿过去"＝**真的仍能**（两种形状各量一发）；注释没说"append 会漏"＝**append 确实不漏**（B 行队列长度仍 1）。
`cloneParamsMap` 只做 `out[k]=v`，与"一层"的措辞**逐字一致**。这一半**不退回**。

### 2.3 "第 9 枚引用字段不会静静进来"那一半：**证伪**（本程造出来了两发）

**M-a：新增一枚 map 槽位、fixture 未填 ⇒ 两发检全绿，而洞是活的。**

```
# snap2/internal/tools/gate.go 的 Decision 里加一枚 cloneDecision 不拷的引用字段：
	Scratch map[string]string
cd D:/tmp/wisp-146acc1/snap2
go test -count=1 -v -run 'TestLiveApprovals(SharesNoReferenceSlot|InPlaceWriteCannotReachTheQueueRecord)' ./internal/agent/approval/
→ rc=0
  --- PASS: TestLiveApprovalsSharesNoReferenceSlotWithTheQueue (0.00s)
  --- PASS: TestLiveApprovalsInPlaceWriteCannotReachTheQueueRecord (0.00s)          # mut\m-a-ninth-map-no-fixture.txt
# 同一棵树上本程另放一枚探针，走的就是 AC#2 禁的那个形状，只是打在 Scratch 上：
go test -count=1 -v -run TestZZAcceptorNinthSlotAliasesAndProbesMissIt ./internal/agent/approval/
→ :34 AC#2-SHAPE HOLE IS LIVE: writing through Decision.Scratch rewrote the admitted record
      (queue now reads grant=CLOBBERED-BY-READER).                                   # mut\m-a-live-hole.txt
```

⇒ 一句 `tools.Decision` 加字段就能让"就地写返回值改掉已决断记录"这条路**重新打开，而票 146 的两发检全绿**。
原因不神秘：probe 1 走的是**值**（`sharedDecision()` 那枚 fixture 填了什么就数什么），不是**类型**。
fixture 没填 ⇒ `refSlots` 一条不出 ⇒ 集合断言仍是 8 vs 8 ⇒ 恒绿。

**M-c：`refSlots` 对 Ptr／Interface 种类根本不看 ⇒ 就算 fixture 填了也照不到。**

```
refSlots on {Ptr *[]string, Map map[string]string} = [Map]
CONFIRMED BLIND SPOT: refSlots names only [Map] …
```

（`ticket146_liveapprovals_backing_test.go:81-100` 的 `switch v.Kind()` 只有 `Map`/`Slice`/`Struct` 三支，
外加 `:90` 的 `if !f.IsExported() { continue }` ⇒ Ptr／Interface／Chan／Func／私有字段五类**静默跳过**。）

⇒ **`pending_read.go:36-38` 与测试文件自己的两处同源表述（`:18-22`"a future reference field … is named by this case
rather than discovered by an exploit"、`:140-142`"neither a field added to Decision nor a walker that quietly stopped
descending can turn probe 1 into an always-green check"）都是**过许诺**。第二处尤其直接：M-a 就是"一枚加到
Decision 上的字段"，它**没有**让 probe 1 变红。
**这一格方向按派单单是"收紧"，所以不给"把注释改弱"的修法。** 最小闭合集合（交回实现者，任选其一即可，本程不代做）：

1. 把普查改成**按类型**数槽位（`reflect.TypeOf(tools.Decision{})` 递归导字段，凡 kind ∈ {Map, Slice, Ptr, Interface, Chan, Func}
   或非 exported 引用型 ⇒ 必须出现在名册里），fixture 只需为"名册存在性"负责，不必为"记忆"负责；
   `cloneDecision` 与名册一不一致当场红；或
2. 至少把三处文字**收窄到真做到的那格**："`refSlots` 只认 map／slice，且只认 fixture 填过的槽位 ⇒ 第 9 枚
   map/slice 字段若 fixture 未填、或任何指针／接口型字段，仍会静静进来"。⇒ 文字诚实了，但**洞仍在**，
   那一洞要另记一档（本程把它记在 §7 第 1 条）。

**2.4 顺带核掉的另外两处跨文件许诺（这一族叫"注释预先引用尚未产出的读数"）**

| 注释句 | 现量 | 判定 |
|---|---|---|
| `:99` "The bound is the queue's own (`DefaultMaxPending`, `queue.go`)" | `internal/agent/approval/queue.go:112  DefaultMaxPending = 8` | **对得上** |
| `:99-101` "the pump does not run on streamed text deltas (`cmd/wisp/panel_pump.go`)" | `cmd/wisp/panel_pump.go:236  // Not on streamed text deltas.` | **对得上** |
| `:101-103` "its readings … are recorded in `docs/evidence/s1/146-liveapprovals-r1.md` §4" | §4 存在、A/B 两张表都有数 | **对得上**（读数真伪见 §6 AC#3） |
| `:98` "one map plus up to seven fresh backings per pending item" | `cloneDecision` 8 条赋值，`cloneBacking`/`cloneParamsMap` 对 nil 早返 | **"up to" 用词准确** |

⇒ **§2 净结论：深度边界那半写得比代码好；覆盖面那半写得比代码好。两处都属"注释许诺 > 仪器做到"，
但只有覆盖面那半本程造出了活洞 ⇒ 该格退回。**

---

## 3. 进攻①：复跑改前 —— **会红，读数与实现程所报逐字对得上**

命令（`snap`＝`7f75ff3` 纯净副本，只把接线摘回票面说的那个形状；`cloneDecision` 留着不响）：

```
# snap/internal/agent/approval/pending_read.go:117
#   - 			Decision:      cloneDecision(it.Dec),
#   + 			Decision:      it.Dec,
cd D:/tmp/wisp-146acc1/snap
go test -count=1 -v -run 'TestLiveApprovals(SharesNoReferenceSlot|InPlaceWriteCannotReachTheQueueRecord)' ./internal/agent/approval/
```

读数（`D:\tmp\wisp-146acc1\mut\m1-revert-all.txt`）：**`rc=1`／`^--- FAIL`=2 行／`AC#2 RED`=14 条**

```
--- FAIL: TestLiveApprovalsSharesNoReferenceSlotWithTheQueue (0.00s)
--- FAIL: TestLiveApprovalsInPlaceWriteCannotReachTheQueueRecord (0.00s)
:196 共用底层：[Args Blacklist.Absolute Blacklist.AlreadyUnlocked Blacklist.Unlockable Capabilities Params Paths RulesHit]
:233 队列存的 Params 里出现了调用方塞进的键 "injected-by-caller"（键 argv 现值 clobbered）
:237 队列存的 Params["argv"] 被就地改成 "clobbered"          :240 队列存的 Args 前 8 字节被就地改成 "XXXXXXXX"
:243 队列存的 RulesHit[0]="R9"，期望仍是 "R2"                :247 队列存的 Paths[0]="C:/clobbered"
:250 队列存的 Capabilities[0]="notify"                       :259 Blacklist.Absolute/Unlockable/AlreadyUnlocked ×3
:272 PanelItem.Paths=[C:/clobbered]
:292 第二次读到 RulesHit[0]="R9"  :295 第二次读到 Capabilities[0]="notify"  :298 第二次读到的 Params 里带着别人塞的键
```

⇒ **检不是装饰。红句点名 8/8 枚，14 条明细、两行 FAIL 与实现程 §3.2 所报完全一致。**
（"摘掉修法有没有任何外部可见读数变过"这一问：有——`rc` 0→1、`--- FAIL` 0→2、`AC#2 RED` 0→14。承重成立。）

## 4. 进攻②：单点回退 —— 覆盖面主张按"8 枚"**逐枚复算成立，实覆 8 枚**

命令（`D:\tmp\wisp-146acc1\sweep.sh`，每一发只删 `cloneDecision` 里的一行，跑同一对用例，读数逐发留在
`D:\tmp\wisp-146acc1\mut\one-{59..66}.txt`；跑完 `cp` 还原并 `diff` 验过 RESTORE-CLEAN）：

| 摘掉的那一行（`pending_read.go`） | rc | `--- FAIL` 行数 | `AC#2 RED` 条数 | probe 1 点名 | probe 2 是否也红 |
|---|---|---|---|---|---|
| `:59 out.Params = cloneParamsMap(d.Params)` | 1 | 2 | 4 | `[Params]` | 是（`:233 :237 :298`） |
| `:60 out.Args = cloneBacking(d.Args)` | 1 | 2 | 2 | `[Args]` | 是（`:240`） |
| `:61 out.RulesHit = …` | 1 | 2 | 3 | `[RulesHit]` | 是 |
| `:62 out.Paths = …` | 1 | 2 | 3 | `[Paths]` | 是（含 `:272` 投影那一发） |
| `:63 out.Capabilities = …` | 1 | 2 | 3 | `[Capabilities]` | 是 |
| `:64 out.Blacklist.Absolute = …` | 1 | 2 | 2 | `[Blacklist.Absolute]` | 是（`:259`） |
| `:65 out.Blacklist.Unlockable = …` | 1 | 2 | 2 | `[Blacklist.Unlockable]` | 是（`:259`） |
| `:66 out.Blacklist.AlreadyUnlocked = …` | 1 | 2 | 2 | `[Blacklist.AlreadyUnlocked]` | 是（`:259`） |

⇒ **"撤掉这一处，哪条用例变得不响？"＝没有一条。8 枚里任意一枚单独回退，两发检都各自红，且红句逐字点名被撤的那一枚。**
票面漏计、被实现程补上的那三枚（`Blacklist.*`）**不是靠"总条数变多"蒙上的**，是各自有独立红句。
⇒ **覆盖面主张可以从"8 枚"照写，不必改写**；要改写的是 §2.3 那句"第 9 枚不会静静进来"。
（两枚洞的合体已单造一发：`M-a` 同时满足"新槽位 ＋ fixture 未填 ＋ cloneDecision 未拷"，结果 5 枚读数全绿 ⇒ 见 §7 第 1 条。）

<!-- 以下各节尚未取证 -->
