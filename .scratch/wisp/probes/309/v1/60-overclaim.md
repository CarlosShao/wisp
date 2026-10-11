# 309-v1 · 承重诚实问题：形 Ⓑ 交付的"防御"＝**零**，那么盘上有没有哪句话把它写成了"防住了"

## 0. 先把"零行为变更"这句本身实测（⛔ 抄实现腿的数）
| 尺 | 射程 | 读数 |
|---|---|---|
| 两版 blob 剥注释后 `cmp`（`logs/ac3-comment-only.log` §3） | `bridge.go` | **IDENTICAL**（768⇄768 代码行）⇒ 产码语义一字节没动，"行为逐字节不变"**为真** |
| 两枚 236 空 id 载体（真把 `TaskID: ""` 喂进 `Execute` 的那两枚） | 三张台面：工作树／`1ad48333` 导出树／HEAD 导出树 | **三处全 PASS**（`logs/suite-worktree.log`、`logs/suite-before.log`、`logs/suite-after.log`；grep 行逐字 `--- PASS: Test236R2TaskSpawnRefusesWhenHostGaveNoTaskID` / `…TaskCancelRefusesWhenHostGaveNoCallerID`） |
| 整包名册作差（同层尺，两棵各自 `mktemp -d`） | `internal/tools` 顶层 `--- PASS` 行 | ADDED **2**（就是两枚 309 用例）／REMOVED **0**／`--- SKIP` **0**／两棵树各余 **1 枚**红＝`TestTicket252R1AlignmentAddsNoAuthorization`（**改前改后同一枚**＝导出树台面缺件，⛔ 归 309；工作树那一枚 **PASS**、整包 rc=0） |
⇒ 编排者裁语那句"行为逐字节不变 ⇒ **对空 id 的实际防御＝零**"**在盘上成立**；今天存在的东西只有两枚：`AC#1` 那枚跳线仪器＋`AC#2` 那枚枚数钉。

## 1. 逐句审（只列承重的；⛔ 形容词，每句带"盘上有没有对应的量"）

| # | 句（逐字，带号） | 判定 | 依据（本腿自己的读数） |
|---|---|---|---|
| 1 | `bridge.go:937` "Two things keep that sentence honest, **both checkable now**" | ⚠ 半过 | 第一支的"可核"是枚数钉，钉只认字面量形（`20-ac2-nail.md` §3：形 2/3 实测 rc=0）；第二支（`scopeMarks` 那枚地板）**⛔ 任何断言**，只靠读码 |
| 2 | `bridge.go:935` "**Every other gate still runs**; only R4 is dropped" | ⚠ 未证 | 结构上真（`bridge.go:211-213`：`b.assess` 仍带 canonicalizer＋classifier，只缺 detector），但**盘上⛔ 一枚断言**为"空 id 那一发仍报得出 R2/R3"负责；仪器的两腿都只问 R4 ⇒ `Q-309-v1-other-gates-uninstrumented` |
| 3 | `bridge.go:947` "**Who rings** if this limb's shape changes is Test309…" | ✗ **overclaim** | "这一肢的形状变了就有人响"是**全称**。本腿把守卫另外两肢各摘一枚（`logs/mut-e.log`／`mut-f.log`）**两发全绿**＝形状变了⛔ 响。该句射程其实只有"空 id 那一支的挂/不挂"两向（`948-949` 那两行自己就是这么写宽的，⛔ 但 `:947` 那行比它宽） ⇒ `Q-309-v1-who-rings-overclaim` |
| 4 | `bridge.go:948-950` "it fails if a weakened guard attaches… and fails if the attach branch stops attaching… **that test is the only instrument**, and this comment states the debt, **not a guarantee**" | ✅ 老实 | 两向各自被本腿打红（形 A/C 与形 B/D），"⛔ 是保证"这句正对齐实测 |
| 5 | `bridge.go:938-941` "No production caller reaches the limb **today**… pins `agent.ToolRequest{` construction points outside internal/agent to 0" | ✅（带限定） | "today"＋backtick 锁死拼法＋具名钉，⛔ 吹成永久；`CloseTask` 那块也逐字写着 "As measured today (a present count, **not a permanent promise**)" |
| 6 | 钉头注释 `toolrequest_nail_309_test.go:3-4` "将来**任一** caller 手搭 agent.ToolRequest{ 递进空 TaskID" ⇒ "要变成'**有人新写构造点就会红**'的枚数钉" | ✗ **overclaim** | 与本票 `AC#2` 格文同一个"任一"缝：形 2（声明＋赋值）、形 3（别名）实测⛔ 响；且票 154 验收程早就在合成树上量过五形全逃（`50-ac5-recon.md` §6）⇒ 同一句在**产码注释**与**测试头注释**各出现一次，两处都缺"仅字面量拼法"这半句 |
| 7 | 仪器头注释 `assessorfor_taint_instrument_309_test.go:4-5` ""响"只能靠这枚仪器＋…枚数钉，**别处没有牙**" | ✅ 老实 | 与本腿读数一致（守卫另外两肢⛔ 仪器；`b.injected` 生产恒 false 已另记） |
| 8 | 仪器头 `:11-14` ②"对空 id 的**实际防御＝零**" | ✅ 老实 | 逐字承认零防御，⛔ 写成"会响的守卫" |
| 9 | 票尾 bullet（`b2dd7fb2`，逐字）"恒真形两向都有牙…Rig 证"加一枚即红"" | ⚠ 半过 | 前半句被四形突变兑现（真话）；后半句"加一枚即红"⛔ 限定拼法 ⇒ 与 #6 同缝，**证据件措辞**（⛔ 产码，归实现腿/编排者面） |
| 10 | 票尾 bullet "`AC#3`＝…那句**带尺**现量重述" | ✗ **overclaim** | 尺带上了，但那把尺**第一发重跑就对不上自己写的数**（`30-ac3-anchors.md` §3：注释自称 1、实测 3，自指）⇒ "带尺"的形式过了，"带一把重跑对得上的尺"没过去 |

## 2. ⛔ 一处⛔ 出现"防住了 / 会响的守卫 / covered"式措辞（全文扫过）
尺＝`grep -n "防住\|covered\|guarantee\|保证\|会响\|honest" internal/tools/bridge.go internal/tools/*_309_*_test.go`
⇒ 命中 12 枚，逐枚看过：**⛔ 一枚"防住了 / 已覆盖 / 有保证"式的正面声称**。
`covered` 那 4 枚（`:332`/`:439`/`:493`/`:1046`）＝**既有文本**、讲的是 D45 授权覆盖，与本票⛔ 相干；
`保证` 两枚＝`bridge.go:950` "not a guarantee"（否定式，✅）与 `toolrequest_nail_309_test.go:109`
"票 309 的保证＝今天非空只靠 loop.go 那枚 newTaskID() 构造点（**接线事实**）"（⛔ 否定式但自己就把口径写成了接线事实，✅）；
`honest` 三枚＝`bridge.go:937`（本票新增，见 §1 第 1 行⚠）与两枚既有文本。⇒ ⛔ 有人把零防御写成有防御；
**overclaim 全部集中在"响的射程"（哪一形会响、哪一肢会变），⛔ 在"有没有防御"**。这个区别重要：前者是措辞缺陷，后者才是欺骗。

## 3. 判给票面的口径（该由编排者落的那一句）
形 Ⓑ 之后盘上真实存在的＝**两枚跳线**（1 枚定向、1 枚认字面量），⛔ 一道防线。裁决表因此⛔ 允许任何一格读成"空 id 已有人管"；`AC#2` 只裁到"最小事"那一半，"任一 caller"那一半**⛔ 闭合**（`Q-309-v1-nail-misses-*`，且⛔ 新债、先例在票 154 族）。
