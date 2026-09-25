# 146-r1 — `Queue.LiveApprovals()` 的 "copies"：裁决表（实现程自证件）

- 工单：`.scratch/wisp/issues/146-liveapprovals-says-copies-but-the-value-shares-a-map-and-three-slices-with-the-queue.md`
- 写手：`agent=ticket146`（实现程，**同一程写码又自证，所以本件不是验收件**；对抗验收按 `issues/README` 规则 6 由另一程做）
- 时刻：`2026-09-25 19:0x–19:5x +08`（`11:0x–11:5x Z`）
- 锚点：开工 HEAD = `630c218`，分支 `dev`，工作树含别人的未提交改动（见 §6 AC#4 的 `git diff --name-only` 清单）
- 临时件（只建不删）：`D:\tmp\wisp-146-agent-a\`——`before.txt`（改前全包 `-v` 输出）、`roster-before.txt`（改前点名册）、
  以及下面各节点名的改后读数与变异副本输出。**本件任何一处不抄凭据值。**

> **读法警告**：本件里所有数字都是**我在 `630c218` 之上的工作树现量**的，命令原文逐条附在下面。
> 票面 18:2x 的三个数与我的读数有两处不符，逐条点名在 §1.4，**没有一处是我抄票面抄来的**。

---

## 1. 现量（AC#1 的判据材料）

### 1.1 `LiveApprovals(` 的调用者枚数

命令（仓根，Git Bash）：

```
grep -rn "LiveApprovals(" --include="*.go" .
grep -ro  "LiveApprovals(" --include="*.go" . | sort | uniq -c
```

读数（`630c218` 之上的工作树）：

```
./cmd/wisp/panel_pump.go:61:                          items := rt.gate.Queue().LiveApprovals()
./internal/agent/approval/pending_read.go:42:         func (q *Queue) LiveApprovals() []LiveApproval {
./internal/agent/approval/pending_read_test.go:       8 处（:121 :133 :141 :158 :160 :170 :201 :226）
```

| 口径 | 票面 :25-26 说 | 我量到 |
|---|---|---|
| `.go` 里 `LiveApprovals(` 总处数 | 11 | **10** |
| 定义 | 1 | 1（`pending_read.go:42`） |
| 生产调用 | 1（`cmd/wisp/panel_pump.go:61`） | 1（同一枚，行号也一致） |
| 本包测试 | 9 | **8** |

⇒ **差的那 1 枚在测试侧**（票面 9、我 8），生产腿两处一致。票面 :30 已自declare"下一程一律现量，别抄我的数"，
所以这是**票面过期**、不是形状之争：**没有第六处调用者**，生产出口仍只有 `cmd/wisp` 一枚。

### 1.2 被拷贝的结构到底有几枚引用类型字段

命令：

```
sed -n '15,53p' internal/tools/gate.go            # Decision 的字段表
grep -n "type BlacklistNote" -A 16 internal/tools/mode.go
grep -n "type Capability" -A 2  internal/tools/capability.go
grep -n "type Kind"      -A 2  internal/tools/registry.go
```

现量的字段表（`tools.Decision`，`internal/tools/gate.go:15-53`）：

| # | 字段 | 类型 | 引用？ | 票面 :18-23 的表里有吗 |
|---|---|---|---|---|
| 1 | `Tool` | `string` | 否 | — |
| 2 | `Provider` | `Kind`＝`string`（`registry.go:15`） | 否 | — |
| 3 | **`Params`** | **`map[string]any`**（`gate.go:18`） | **是** | ✅ |
| 4 | **`Args`** | **`json.RawMessage`＝`[]byte`**（`gate.go:19`） | **是** | ✅ |
| 5 | `Level` | `risk.Level`（整型） | 否 | — |
| 6 | **`RulesHit`** | **`[]risk.RuleID`**（`gate.go:21`） | **是** | ✅ |
| 7 | `Reason` | `string` | 否 | — |
| 8 | **`Paths`** | **`[]string`**（`gate.go:23`） | **是** | ✅ |
| 9 | **`Capabilities`** | **`[]Capability`**，元素 `Capability＝string`（`capability.go:15`） | **是**（头共用、元素不共用） | ✅（票面写"元素本票没量"——我量了：**是 string**） |
| 10-13 | `SessionOverrideBlocked`/`Mode`/`ModeSilenced`/`ModeKept` | `bool`/`risk.Mode`/`bool`/`string` | 否 | — |
| 14 | **`Blacklist`** | **`BlacklistNote`**（`gate.go:43`）→ 结构体内含 **`Absolute []string` / `Unlockable []string` / `AlreadyUnlocked []string`**（`mode.go:66-81`） | **是 ×3** | ❌ **票面漏计** |
| 15-18 | `DecisionColumn`/`CorrelationID`/`TaskID`/`CallID`/`Timeout` | `string` ×4 / `time.Duration` | 否 | — |

⇒ **引用槽位是 8 枚（1 枚 map + 7 枚 slice），票面点了 5 枚。** 漏的三枚全在 `Blacklist BlacklistNote` 里面。
这**不推翻票面的形状判断**（共用底层是真的、`Params` 那枚"比切片更硬"也是真的），只是把它**量小了**：
按 ⓑ 只改票面那句注释，会把一次漏计**钉成文字**——注释仍会漏说 `Blacklist` 那三枚切片与队列共用底层。
⇒ 已按"复算不符"上报（票面 :81），并且这条直接进入 §2 的修法判定。

`Params` 不是死字段（这条会影响"风险是不是活的"）：

```
grep -rn "\.Params = \|Params:" --include="*.go" internal/tools/ | grep -v _test
→ internal/tools/bridge.go:267:  dec.Params = params
```

`params` 是 `decodeArgs(req.Args)` 每次调用现解的一份 JSON `map[string]any`（`bridge.go:262`），
随后原样进 `qitem.Dec`（`queue.go:156` `Dec: d` 按值存）。⇒ **每一张 L2 卡的 `Params` 都带着内容进队列**，
所以"就地写它会改掉队列里那条记录"不是理论形状，是随时可发生的形状。

### 1.3 AC#1 那句可核判据的现量：有没有调用方就地写过返回值里的 map／切片

判据（票面 :48-50 原文）：今天有没有**任何**一枚调用方（含测试）在拿到返回值之后对那枚 map 或那几枚切片做
**改元素／append／就地写／取地址后传出去**。

命令（三条，逐条读，不合并）：

```
# (a) 生产出口对返回值做了什么
grep -rn "RulesHit\|Paths\|Params\|Args\|Capabilities" cmd/wisp/panel_pump.go
# (b) 全仓切片/map 就地写形状（剔测试）
grep -rn "RulesHit =\|RulesHit\[\|Paths\[\|Params\[\|Capabilities\[\|Args\[" --include="*.go" . | grep -v "_test.go"
# (c) 有没有人取返回值或其字段的地址
grep -rn "&it\.Decision\|&items\[\|&got\[" --include="*.go" .
```

读数：

- (a) `cmd/wisp/panel_pump.go` 对 `it.Decision` 只有 6 处**读**：`Tool`/`Args`（经 `panelArgs` 读成字符串）/`Level`/
  `RulesHit`/`Reason`/`SessionOverrideBlocked`，全部是"取出来塞进 `panel.NativeVerdict`"，**没有一处 `[…] =` 形状**；
  `panelArgs(it.Decision)` 是**按值**传参（`panel_pump.go:96`），`&` 只出现在同函数里一枚本地 map（`:102` 的 `&params`）。
- (b) 全仓（剔 `_test.go`）切片就地写形状 **4 处**：`internal/agent/approval/gate.go:232`（自己 `append` 到新底层数组，
  构造新 verdict）、`internal/panel/approval.go:92,93`（`view.RulesHit` 是 `[]string`，**与 `risk.RuleID` 不同型、
  且是整枚字段重指不是改元素**）、`internal/tools/bridge.go:275`（`dec.RulesHit = verdict.RulesHit`，在 `LiveApprovals`
  **上游**，是构造 verdict 不是消费返回值）。⇒ **0 处**落在"拿到 `LiveApprovals()` 返回值之后就地写它"。
  （票面 :26-27 量的那条命令数出 3 处；同一条命令我今天数出 **4 处**——多出的那枚是 `approval.go:93`，
  形状与票面结论同向，不影响判定。）
- (c) **0 处**。
- 本包 8 处测试调用（`pending_read_test.go`）逐枚读过：`assertRow` 是**按值**收 `LiveApproval`（`:54`），
  其余全是 `len()`／下标读／比较，**没有写形状**。

⇒ **判定材料：全无。**

### 1.4 票面 vs 现量，逐条点名（票面 :81 要求"复算不符 ⇒ 写进证据件并报回"）

| 项 | 票面说 | 我量到 | 原因／后果 |
|---|---|---|---|
| `LiveApprovals(` 总处数 | 11（1+1+9） | **10**（1+1+8） | 测试侧少 1 枚（票面在 18:2x 那棵树上量的；现树 `pending_read_test.go` 是 8 行命中）。**不影响判定**：生产出口仍只有 1 枚。 |
| `RulesHit =`/`RulesHit[`（剔测试）全仓 | 3 处 | **4 处** | 多出 `internal/panel/approval.go:93`（另一枚同型重指）。**不影响判定**：4 处全在返回值消费点之外。 |
| `Decision` 的引用字段 | 5 枚（1 map + 4 slice 形） | **8 枚槽位**（1 map + 7 slice） | 票面漏计 `Blacklist BlacklistNote` 里的三枚 `[]string`（`gate.go:43` → `mode.go:66-81`）。**影响修法范围**：见 §2 判据 ③。 |
| `pending_read.go` 注释行的位置 | `:19-20`，"copies" 在行首 | 同（现树 `:19` 就是含 `copies` 那行） | 一致。 |
| 生产调用行号 | `cmd/wisp/panel_pump.go:61` | 同 | 一致。 |

---

## 2. AC#1 裁决：**ⓐ 真拷一层**（与票面默认方向 ⓑ 相反，理由与张力全部摆在下面）

**结论一句**：AC#1 的普查（§1.3）答"全无"，按票面 :50 的字面**ⓑ 是被允许的**；但同一张票的 AC#2 把"会响的检"的
断言极性写死了（"**断言队列里那条原始记录没有跟着变**"），而那一句**在 ⓑ 之下永假**，本票 AC#2 又被票面自己标为
"本票的硬核心" ⇒ **两枚 AC 只有 ⓐ 同时收得下**，所以判 ⓐ。

三条可核理由：

1. **AC#2 的字面断言只在 ⓐ 下为真。** ⓑ 的定义就是"保留 `Decision: it.Dec` 那枚浅拷贝"，
   那么在浅拷贝之上写"就地改掉返回值 ⇒ 队列里那条原始记录**没有**跟着变"是一条**必红**的断言。
   把它反向写成"记录**一定**跟着变"确实能绿，但那等于**把危险形状钉成规格**，而且踩本仓那族有名字的缺陷：
   票面 :54 的"恒真判据是一类新假绿"——一条除了把实现改成 ⓐ 之外永远不会红的用例，就是装饰。
2. **"把修法摘掉，这一发是不是从此打不红"对 ⓑ 恒答"仍不红"。** 票面 :52-53 的承重定义是按**可执行**给的：
   摘掉修法 ⇒ 用例失红。ⓑ 的修法是**一句注释**，注释不在编译与执行路径上，摘掉它**任何用例都不会红**
   ⇒ 按票面自己的作废条款（:53"答'仍不红'＝这枚检是装饰，本票作废"），ⓑ 这支**结构上**交不出一枚承重的检。
3. **ⓑ 唯一可能的仪器在本票允许的落点之外。** 一条"调用方不许就地写"的纪律，能配的非运行时仪器只有**静态扫**，
   而本票 AC#4 把 `tools/d22scan/**` 划成**零字节**，`internal/panel/l2_grant_boundary_test.go` 同理零字节
   ⇒ 在本票允许改的文件里，ⓑ 拿不到任何仪器。**ⓐ 不需要新纪律：它让注释说的那句字本身成立。**

成本与代价（票面 :36-38 要求"没量过就不许写代价可忽略"）：

- 成本是**每次泵动 × 每枚 pending 项**多一张 map（键数＝该次调用的 JSON 参数数）＋ 7 枚底层数组；
  队列上限 `DefaultMaxPending = 8`（`queue.go:112`），泵**不随 streamed delta 动**（`panel_pump.go:236-240`）。
- 这条成本的**档**属于 `PLAN.md` 的 D32 资源预算那一族 ⇒ 按票面 AC#3 现量，读数见 §4。
  **我没有把"代价可忽略"当结论写**；§4 给的是量到的数，量不到的那半标"未测、需安静窗口"。

⚠ **这一格要编排者复核的点**：票面默认 ⓑ、我判 ⓐ，**方向相反**。我判的依据全在上面三条，且都在票内文字里；
不是我偏好深拷贝。若编排者改判 ⓑ，需要同时改 AC#2 的断言极性（那就是改票＝人工批准），并按理由 3 另派一票给
`tools/d22scan`。**代码与检我已按 ⓐ 落地，摘掉 ⓐ 只需 revert 一枚小 commit（见 §3）。**

---

## 3. AC#2 裁决：会响的检 —— **两发都量过"摘掉修法必红"**

落点：`internal/agent/approval/ticket146_liveapprovals_backing_test.go`（`package approval` 白盒，
与 `pending_read_test.go` 同一棵尺、复用它的 `mustPush`/`mustFindLive`）。

### 3.1 两发是什么

| 发 | 名字 | 判据 |
|---|---|---|
| probe 1 | `TestLiveApprovalsSharesNoReferenceSlotWithTheQueue` | 反射走 `tools.Decision`，点名返回值里**哪几枚引用槽位仍与队列存储项共用底层**（`reflect.Value.Pointer()` 比数据指针/hmap 指针）。槽位**名册本身按集合断言**（8 枚，见 §1.2）⇒ walker 失效或 `Decision` 加/减引用字段都会红，不做恒真判据。 |
| probe 2 | `TestLiveApprovalsInPlaceWriteCannotReachTheQueueRecord` | AC#2 的字面：**只通过返回值**就地写全部 8 枚槽位（map 改键＋插键、`copy()` 进 `Args`、下标写 6 枚切片），然后从**三条路**读回队列那条记录：白盒 `qitem.Dec`、队列自己的不可信投影 `view()`/`head()`、第二次 `LiveApprovals()`，要求原样。 |

### 3.2 "承重"那一句，按票面 :52-53 的操作定义答

**问**：把 AC#1 选定的修法（ⓐ 的 `cloneDecision`）摘掉，这一发是不是从此打不红？
**答：会红。且这条不是推理，是量过的。** 我先写检、后装修法，中间那一发的真读数：

命令：

```
go test -count=1 -v -run 'TestLiveApprovals(SharesNoReferenceSlot|InPlaceWriteCannotReachTheQueueRecord)' ./internal/agent/approval/
```

改前（`pending_read.go` 仍是 `Decision: it.Dec`）：`rc=1`，两行 `--- FAIL:`，**14 条** `AC#2 RED` 明细，
probe 1 逐字点名全 8 枚：

```
--- FAIL: TestLiveApprovalsSharesNoReferenceSlotWithTheQueue (0.00s)
--- FAIL: TestLiveApprovalsInPlaceWriteCannotReachTheQueueRecord (0.00s)
ticket146_liveapprovals_backing_test.go:196: AC#2 RED: ... 共用底层：[Args Blacklist.Absolute
  Blacklist.AlreadyUnlocked Blacklist.Unlockable Capabilities Params Paths RulesHit]
ticket146_liveapprovals_backing_test.go:233/237: 队列存的 Params 被塞键 / argv 被就地改成 "clobbered"
ticket146_liveapprovals_backing_test.go:240: 队列存的 Args 前 8 字节被就地改成 "XXXXXXXX"
ticket146_liveapprovals_backing_test.go:243: 队列存的 RulesHit[0]="R9"（期望 R2）
ticket146_liveapprovals_backing_test.go:247: 队列存的 Paths[0]="C:/clobbered"
ticket146_liveapprovals_backing_test.go:250: 队列存的 Capabilities[0]="notify"（期望 shell）
ticket146_liveapprovals_backing_test.go:259: Blacklist.Absolute / Unlockable / AlreadyUnlocked 三枚均被就地改
ticket146_liveapprovals_backing_test.go:272: PanelItem.Paths=[C:/clobbered]      ← 投影也被污染（见 §3.4）
ticket146_liveapprovals_backing_test.go:292/295/298: 第二次读到的三枚都是上一枚读者改过的值
```

日志：`D:\tmp\wisp-146-agent-a\ac2-before-fix.txt`（只建不删）。装上 ⓐ 之后同一对用例 PASS
（`D:\tmp\wisp-146-agent-a\after.txt`，`rc=0`）。

⇒ 检**不是装饰**：它对"修法在不在"是敏感的，且敏感面是 8/8 全覆盖，不是随手挑的一枚字段。

### 3.3 三条读路的用意（防止这枚检被读薄）

- 白盒 `qitem.Dec` 那一发是票面要的那句本体。
- `view()`/`head()` 那一发是**第二把尺**：`queue.go:viewLocked` 早就给 `PanelItem.Paths` 做了 `append([]string(nil), …)`，
  所以"投影会拷"这件事**并不能保护队列的记录**——改前那一发照样红，因为它读的是**已被改掉的源**。
  这一发同时钉住另一件事：将来谁把 `viewLocked` 那枚拷贝摘掉，卡片会直接把调用方的改写画给批准者看。
- 第二次 `LiveApprovals()` 那一发对应泵的**真实用法**（`cmd/wisp/panel_pump.go:61` 每次状态变动重读一次）：
  一次污染会**跨帧**留在面板上，不是一次性的。

### 3.4 这一格没买到的东西（写清楚，别让它被读成"全隔离"）

拷贝是**一层**深。`Params` 的 value 里那层容器（`decodeArgs` 留下的 `[]any`／嵌套 `map[string]any`，
`bridge.go:262`）**仍与队列共用**——调用方把 `Params["argv"].([]any)[0]` 改掉，队列的记录还是会动，
而本票这两发**不探这一层**。原因：真要做穿到底，得递归克隆任意 `any`（面对 cycle／chan／func 都得另写一套），
代价与本票 AC#1 判出来的那一点收益不成比例；今天也没有任何调用方走这条路（§1.3 普查）。
⇒ 这一条同时写在 `pending_read.go` 的头部注释里（"Treat those values as read-only too"），
不藏在证据件里。**如果编排者认为这一层也要钉，那是一枚新票的体量，不是本票剩下的半格。**

---

## 4. AC#3 裁决：WorkPeak 前后读数 —— **未观察到差异，且这台仪器结构上看不见这枚改动**

### 4.1 测量前提（先确认安静，再取数）

```
gh api repos/<owner>/<repo>/actions/runners --jq '.runners[] | {status,busy}'
→ {"busy":false,"status":"online"}          # wisp-selfhosted-01，取数前确认，取数后再确认一次仍 busy:false
gh run list --limit 8 → 全部 status=completed（最近一枚 10:26Z，取数时刻约 11:15Z）
```

⇒ **runner 空着，这一格不是脏数。** 本机另有两件事同时发生，都记在这里：
① `cmd/wisp/slo_windows.go` 当时**被别人改着且不能编译**（`git status` = ` M`，3 枚未用 import），
   所以我**不在工作树上建 `wisp`**，改从两棵 `git archive` 的**纯净快照**建（见 4.2）；那一枚 WIP 我一字节未动。
② `internal/panel` 的 `TestC21DesignTokensFourWayAgree` 当时红着，根因是 owner 未提交的 `design/` 移动
   （`design/assets/tokens.css` 在工作树里不存在）——与本票无关，登记在 §6.3。

### 4.2 命令原文与读数

**A. 票面点名的那台仪器（`wisp slo`，WorkPeak 档，同形发法，双臂各 3 次交替）**

```
git archive 630c218 | tar -x -C /d/tmp/wisp-146-agent-a/tree-before      # 改前那棵树
git archive 22f7b1a | tar -x -C /d/tmp/wisp-146-agent-a/tree-after       # 改后那棵树（只差本票那一枚 commit）
go build -o wisp-before.exe ./cmd/wisp   （在 tree-before 内）
go build -o wisp-after.exe  ./cmd/wisp   （在 tree-after  内）
for i in 1 2 3; do for arm in before after; do
  WISP_TEST_DATA_DIR=…/data/$arm-$i ./wisp-$arm.exe slo -state WorkPeak -seconds 5 > run-$arm-$i.json
done; done
```

| arm | `mem_median_bytes`（3 次） | 均值 | arm 内极差 | `tree_private_bytes` 判据行 |
|---|---|---|---|---|
| before | 4112384 / 4075520 / 4214784 | **4134229** | 139264 | 3.9MB / 3.9MB / 4.0MB，`pass=true`，**`gate=false`** |
| after | 4038656 / 4067328 / 4153344 | **4086443** | 114688 | 3.9MB / 3.9MB / 4.0MB，`pass=true`，**`gate=false`** |

⇒ **均值差 −47787 字节（−1.16%），小于任一臂的组内极差（115 KB／139 KB）** ⇒ 按票面 :56 的要求写：
**"未观察到差异"**，**不写**"无代价"。CPU 两臂都是 0，`handles_max` 140 vs 140/146，`threads_max` 12 vs 12/13。
**阈值与 golden 一字节未动**（`memCapWorkPeak = 700<<20` 那一行我没读过就没资格碰，见 AC#4 清单）。

⚠ **更要紧的一句：这一档结构上照不到这枚改动。** WorkPeak 的 subject 是
"the REAL runtime skeleton … boot the real process, apply settle, then **sit still**"
（`cmd/wisp/slo_windows.go:344`，`posture:"skeleton"`，`goroutines_max=1`），
而 `LiveApprovals()` 的唯一生产调用者是 `agentRuntime.liveVerdicts`，它**只被 `wisp run` 接线**
（`cmd/wisp/run.go:422` `Verdicts: rt.liveVerdicts`）。⇒ 被采样那 5 秒里**这枚函数一次都没被调用**，
所以两臂同形不是"代价可忽略"的证据，只是"仪器与改动不在同一条路上"的证据。**这句话必须和上面那张表一起读。**

**B. 这枚改动的机制账（分配／耗时），在同一对纯净快照上做，唯一变量是 `cloneDecision`**

```
cp /d/tmp/wisp-146-agent-a/bench146_liveapprovals_test.go <tree-{before,after}>/internal/agent/approval/bench146_test.go
go test -count=5 -run '^$' -bench BenchmarkLiveApprovalsDepth8 -benchtime=2000x ./internal/agent/approval/
```

| arm | B/op | allocs/op | ns/op（5 次） |
|---|---|---|---|
| before | **3456** | **1** | 6330/2698/3037/5225/3354 |
| after | **8192** | **73** | 15998/22098/24382/17002/15150 |

⇒ 满深度（8 枚 pending，`DefaultMaxPending`）一次 `LiveApprovals()` 多 **4736 字节 / 多 72 次分配**
（8 项 × 9：1 map + 7 slice + map 的桶），约 **+11～18 µs**。这份账**跑在仓外副本上**，
`bench146_liveapprovals_test.go` **没有进仓**（本票只交两枚 `.go` 改动，见 §6.1）。
分配是**暂态**的（快照建完即可回收），留存集不随泵动增长；上限是"队列深度 × 一份"，不是"泵次数 × 一份"。

### 4.3 这一格的结论怎么写才算诚实

- 机制上：**有代价，量到了**（+4736 B / +72 allocs / 一次满深度调用）。
- 档位上：**WorkPeak 那一档未观察到差异**（读数在组内噪声之下），而且这台仪器今天**根本不经过**这枚函数。
-  ⇒ **不许**把这两条合并成"代价可忽略"。合并它需要的证据（一个真会调用 `LiveApprovals` 的 WorkPeak 姿势）
  本票没有，也不该由我在证据件里现造。
- 如果编排者要那一档的真前后对比，缺的不是时间而是**仪器**：`wisp slo` 得有一个"带 pending 审批的 run 姿势"
  才照得到这枚函数。那是登记项，不是本票 AC#3 能顺手收的半格。


---

## 5. 本程没测什么（按"如果我漏了它，谁会先被骗"排序）

（填写中。）

---

## 6. 契约轴与门禁（AC#4 / AC#5）

（填写中。）

---

## 7. 伪授权登记（每程必填）

（填写中。）
