# 票 354 `AC#0` 只读腿 `354-a1`：残余①（复验拒删后，调用者名字下站着攻击者对象）的现读＋归因

- 日期：2026-11-08 21:1x +08。腿＝`354-a1`（票 354 排程第一节：`AC#0` 现读＋归因，⛔ 动产码字节——本腿 `.go` 字节＝0）。
- 读数锚：导出树 HEAD **`ca790c46e`**（两棵 `git clone --shared` 导出树同锚；三核③＝色取自导出树，⛔ 采信共享树半保存态）。
- 凭据：原始仪器日志两枚落盘＝`.scratch/wisp/probes/354/a1/01-ac2windows-instruments.log`（定向 6 发 rc0）＋
  `02-winsec-full-package.log`（`./internal/winsec` 整包 `ok 20.148s` rc0）；读数册＝`.scratch/wisp/probes/354/a1/readings.md`。
- ⛔ 翻框（`AC#0` 归编排者凭本件裁）；⛔ push。

## 1. AC#0 判据四件事逐件交

**① 两族读数现跑（⛔ 抄任何腿的数——本腿两族全部亲跑）**：

| 族 | 用例 | 关键读数（逐字） |
|---|---|---|
| `ancestorunlink_348_*`（Windows） | `...SwappedAncestorNeverDeletesAForeignFile342r1` | `form=symlink-dir attempts=400 cycles=67 provedReparse=67 truncated=0 deleted=0 returnedNil=76 flipErrs=11` |
| 同上（叶形） | `...LeafIdentityReCheckRefusesAForeignNameAndNeverDeletesIt342r1` | `form=leaf-symlink attempts=400 cycles=81 provedReparse=83 truncated=0 deleted=0 returnedNil=66 flipErrs=3` |
| `ancestorworld_342_348_*` 世界（消费者） | `...SteadyStateAncestorJunction...342r2` | 三发 deterministic：两列 0 damage，拒绝来自 act 之前的 C26 稳态腿（用例自陈＝348-v1"缺世界"形的确定性复现，NOT 348 pre-landing numerator） |
| 同世界（353 残留面） | `...NameTopUpBeforeTheReCheckLeavesTheTopUpObjectNamed353r1` | `entered=3/3 ... DAMAGE=0, refusals that left the caller's name standing=3` |
| 同世界（负控） | `...RemovalStillRunsWhenNobodyTopsUpTheName353r1` | 无 top-up ⇒ `name=absent`（移除照跑） |

**② 正控先跑（尺能命中非零，本仓"钝尺读 0"教训）**：分子尺＝353r1 的 shot 级 DAMAGE 行，它把
`the caller's name=present (size 64 body [an object the caller never named, standing in the same directory])`
逐字打出来——名字下站着的是攻击者 64 字节对象本体、三发全中＝非零且确定性；同尺负控（RemovalStillRuns）读 0。
敏感度（摘删除⇒正控转红）的既有记录在 `353-v1` 判语第 3 件三发去牙，本腿未重跑（归 `354-r1` `AC#3`），⛔ 单发当判据。

**③ 残留①读成读数（今天在 `ca790c46e` 上的确定性形状）**：

- **仪器面 1（叶形）**：`ancestorunlink_348_windows_test.go:169` 逐字——
  `the followed open left the bystander as "" (read err <nil>) - the bytes were refused afterwards and the NAME was never removed, which is the half this platform may claim`。
  即：复验拒掉了**字节**，但**名字从未被移除**。
- **仪器面 2（top-up 形，最直白）**：`nametopup_353_windows_test.go` 三发——复验拒绝之后，
  调用者那枚名字 `present`、读出的 body 是攻击者的对象；**而攻击者对象自己的名字 `absent`**
  ⇒ top-up 之后调用者的名字是那个外来对象的**唯一名字**，计数器 `refusals that left the caller's name standing=3/3`。
- **既有守卫行**：`ancestorwindow_342r2_windows_test.go:421` 逐字
  `and it is the attacker's object standing there. NOT covered and NOT claimed: a top-up landing ...`（`353-v1` 第 2 件已对回）。

**④ 归因（为什么残留①今天还在，且为什么它是"留下的"而不是"没做完"）**：

- **机制链**：票 353 落地的同一性复验在 `unlinkPinnedLeaf` 里（`p.root.Lstat(base)` ⇄ `os.SameFile(pinned, now)`，
  不匹配 ⇒ `return`，删除原语 `_ = p.root.Remove(base)` 逐字未动——`353-v1` 第 1 件对码）。
  复验拒绝的那一支**在 Remove 之前就返回**，所以攻击者顶进来的对象**原地留在调用者名字下**——
  这正是把"按名删会删掉调用者⛔点名的对象"关死之后的必然代价，不是回归。
- **归口**：owner 2026-10-21 拍〔乙〕（`AC#1` 已勾、原话"按推荐"）＝**把留着的那枚名字主动报出来（只加读数、不删）**。
  ⇒ 残留①今天的状态＝**已按拍定方向"可观测"（353r1 仪器三发确定性读出），待 `354-r1` 把它升格为具名分子
  `leftStandingForeignName` 的确定性仪器（`AC#2`）＋成对去牙（`AC#3`）**。本腿不越排程替 `354-r1` 动手。
- **⛔ 覆盖不到的那一支（仪器 GREEN SCOPE 逐字在场，归因必须带着它）**：top-up 在复验应答**之后**、
  `Deleteat` 的 `NtOpenFile` 绑定**之前**落地的那一窗——在平台内部，本包无钩、票 353 已逐字承认测不到。
  本腿现读同样测不到，不冒认。

## 2. 射程与纪律

- **Windows 档 only**：两族 unix 半边（`*_unix_test.go`）本机不可跑，登记〔未跑〕，⛔ 用 Windows 读数代 unix 判断。
- 本腿写面＝`docs/evidence/s1/354-a1-residual-readings-r1.md`＋`.scratch/wisp/probes/354/a1/{readings.md,01-*.log,02-*.log}`＋票尾一条 Progress bullet。
  ⛔ 动 `internal/winsec/**` 任何字节（同包单写腿留给 `354-r1`；起手现量 0 在飞）。⛔ 碰 342/348/353 三票任何原句。
- 撤销口令：删退本件＋probes 三枚＋票尾一条即回原状（全部是本腿新增，零改写）。
