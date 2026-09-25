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

## 3. AC#2 裁决：会响的检

（实现与用例落地后填写：判据／命令原文／**改前必红的真读数**／"承重"那一句的操作性答案。）

---

## 4. AC#3 裁决：WorkPeak 前后读数

（填写中。）

---

## 5. 本程没测什么（按"如果我漏了它，谁会先被骗"排序）

（填写中。）

---

## 6. 契约轴与门禁（AC#4 / AC#5）

（填写中。）

---

## 7. 伪授权登记（每程必填）

（填写中。）
